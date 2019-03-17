package redismgr

// Redis connection manager
// Modules use redis connections should be stateless
// that each time fetch a connection from redismgr, then release it when finished.
// Stateless make configuration update at runtime easy.

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/k81/log"
	"github.com/k81/retry"
	"github.com/sony/gobreaker"
)

var (
	mctx    = log.WithContext(context.Background(), "module", "redismgr")
	manager *RedisConnectionManager
)

// RedisConfig defines the redis config
type RedisConfig struct {
	Addrs          []string
	MaxIdle        int
	MaxActive      int
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	Wait           bool
}

// RedisConnectionManager is the redis connection manager
type RedisConnectionManager struct {
	conf     *RedisConfig
	pools    []*redis.Pool
	breakers []*gobreaker.CircuitBreaker
}

// uninit do the clean up
func (mgr *RedisConnectionManager) uninit() {
	for _, pool := range mgr.pools {
		// nolint:errcheck
		pool.Close()
	}
}

func getDialFunc(addr string, connectTimeout, readTimeout, writeTimeout time.Duration) func() (redis.Conn, error) {
	return func() (redis.Conn, error) {
		c, err := redis.DialTimeout("tcp", addr, connectTimeout, readTimeout, writeTimeout)
		if err != nil {
			log.Error(mctx, "dail to server", "redis_server", addr, "error", err)
		}
		return c, err
	}
}

// newRedisConnMgr create a RedisConnectionManager instance
func newRedisConnMgr(conf *RedisConfig) *RedisConnectionManager {
	mgr := &RedisConnectionManager{
		conf:     conf,
		pools:    make([]*redis.Pool, len(conf.Addrs)),
		breakers: make([]*gobreaker.CircuitBreaker, len(conf.Addrs)),
	}

	for idx, addr := range conf.Addrs {
		log.Info(mctx, "creating redis pool", "redis_server", addr)

		mgr.breakers[idx] = gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:        fmt.Sprint("circuit breaker redis-", addr),
			MaxRequests: 10,
			Interval:    5 * time.Second,
			Timeout:     10 * time.Second,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				return counts.TotalFailures >= uint32(20)
			},
			OnStateChange: func(name string, from, to gobreaker.State) {
				log.Info(mctx, "circuit breaker state changed", "name", name, "from", from, "to", to)
			},
		})

		mgr.pools[idx] = &redis.Pool{
			MaxIdle:     conf.MaxIdle,
			MaxActive:   conf.MaxActive,
			IdleTimeout: conf.IdleTimeout,
			Wait:        conf.Wait,
			Dial:        getDialFunc(addr, conf.ConnectTimeout, conf.ReadTimeout, conf.WriteTimeout),
		}
	}

	return mgr
}

// Init initialize the global RedisConnectionManager instance
func Init(conf *RedisConfig) {
	manager = newRedisConnMgr(conf)
}

// Uninit do the clean up for the global RedisConnectionManager instance
func Uninit() {
	manager.uninit()
}

// GetConn return a redis.Conn instance
func GetConn() (c redis.Conn) {
	c, _ = getConnWithCircuitBreaker()
	return
}

// getConnWithCircuitBreaker return a redis.Conn instance and the coresponding breaker
func getConnWithCircuitBreaker() (redis.Conn, *gobreaker.CircuitBreaker) {
	idx := rand.Intn(len(manager.pools))
	return manager.pools[idx].Get(), manager.breakers[idx]
}

// TryRedisCmd retry a redis command with respect of the circuit breaker status
func TryRedisCmd(ctx context.Context, strategy retry.ResettableStrategy, cmd string, args ...interface{}) (reply interface{}, err error) {
	if strategy == nil {
		strategy = defaultRetryStrategy()
	}

	retry.DoWithReset(ctx, strategy, func() bool {
		c, breaker := getConnWithCircuitBreaker()
		defer func() {
			// nolint:errcheck
			c.Close()
		}()

		if reply, err = breaker.Execute(func() (interface{}, error) { return c.Do(cmd, args...) }); err != nil {
			argstr := fmt.Sprintln(args...)
			if len(argstr) > 0 {
				argstr = argstr[0 : len(argstr)-1]
			}
			log.Error(ctx, "failed to try redis cmd, retrying with another connection",
				"cmd", cmd,
				"args", argstr,
				"error", err,
			)
			return false
		}
		return true
	})
	return
}

// TryRedisScript run a lua script with the respect of the circuit breaker status
func TryRedisScript(ctx context.Context, strategy retry.ResettableStrategy, script *redis.Script, keysAndArgs ...interface{}) (reply interface{}, err error) {
	if strategy == nil {
		strategy = defaultRetryStrategy()
	}

	retry.DoWithReset(ctx, strategy, func() bool {
		c, breaker := getConnWithCircuitBreaker()
		defer func() {
			// nolint:errcheck
			c.Close()
		}()

		if reply, err = breaker.Execute(func() (interface{}, error) { return script.Do(c, keysAndArgs...) }); err != nil {
			log.Error(ctx, "failed to try redis script, retrying with another connection",
				"script", script,
				"error", err,
			)
			return false
		}
		return true
	})
	return
}

// defaultRetryStrategy return the default retry strategy
func defaultRetryStrategy() retry.ResettableStrategy {
	s := &retry.AllResettable{
		&retry.MaximumTimeStrategy{
			Duration: time.Second,
		},
		&retry.ExponentialBackoffStrategy{
			InitialDelay: time.Millisecond * 30,
			MaxDelay:     time.Millisecond * 500,
		},
		&retry.CountStrategy{
			Tries: 2,
		},
	}
	return s
}
