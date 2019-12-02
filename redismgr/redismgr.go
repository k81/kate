package redismgr

// Redis connection manager
// Modules use redis connections should be stateless
// that each time fetch a connection from redismgr, then release it when finished.
// Stateless make configuration update at runtime easy.

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

var manager *RedisConnectionManager
var logger *zap.Logger

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
			logger.Error("dail to server", zap.String("redis_server", addr), zap.Error(err))
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
		logger.Info("creating redis pool", zap.String("redis_server", addr))

		mgr.breakers[idx] = gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:        fmt.Sprint("circuit breaker redis-", addr),
			MaxRequests: 10,
			Interval:    5 * time.Second,
			Timeout:     10 * time.Second,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				return counts.TotalFailures >= uint32(20)
			},
			OnStateChange: func(name string, from, to gobreaker.State) {
				logger.Info("circuit breaker state changed",
					zap.String("name", name),
					zap.Any("from", from),
					zap.Any("to", to),
				)
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
func Init(conf *RedisConfig, l *zap.Logger) {
	manager = newRedisConnMgr(conf)
	logger = l
}

// Uninit do the clean up for the global RedisConnectionManager instance
func Uninit() {
	if manager != nil {
		manager.uninit()
	}
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

// Do a redis command with respect of the circuit breaker status
func Do(cmd string, args ...interface{}) (reply interface{}, err error) {
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
		logger.Error("failed to try redis cmd, retrying with another connection",
			zap.String("cmd", cmd),
			zap.String("args", argstr),
			zap.Error(err),
		)
		return nil, err
	}
	return reply, nil
}

// DoScript run a lua script with the respect of the circuit breaker status
func DoScript(script *redis.Script, keysAndArgs ...interface{}) (reply interface{}, err error) {
	c, breaker := getConnWithCircuitBreaker()
	defer func() {
		// nolint:errcheck
		c.Close()
	}()

	if reply, err = breaker.Execute(func() (interface{}, error) { return script.Do(c, keysAndArgs...) }); err != nil {
		logger.Error("failed to try redis script, retrying with another connection",
			zap.String("script", script.Hash()),
			zap.Error(err),
		)
		return nil, err
	}
	return reply, nil
}
