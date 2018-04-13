package redismgr

// Redis connection manager
// Modules use redis connections should be stateless
// that each time fetch a conection from redismgr, then release it when finished.
// Stateless make configuration update at runtime easy.

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	"context"

	"github.com/k81/kate/log"
	"github.com/k81/kate/retry"

	"github.com/garyburd/redigo/redis"
	"github.com/sony/gobreaker"
)

var (
	mctx = log.SetContext(context.Background(), "module", "redismgr")
	gMgr atomic.Value
)

// redis连接池管理器
type RedisMgr struct {
	addrs    []string
	conf     *PoolsConfig
	pools    []*redis.Pool
	breakers []*gobreaker.CircuitBreaker
}

/*************************************************
Description: 清理redis连接池管理器，内部使用
Input:
Output:
Return:
Others:
*************************************************/
func (mgr *RedisMgr) uninit() {
	for _, pool := range mgr.pools {
		pool.Close()
	}
}

/*************************************************
Description: 根据接入地址、连接池参数创建一个redis连接池管理器
Input:
Output:
Return:
Others:
*************************************************/
func newRedisConnMgr(addrs []string, pc *PoolsConfig) *RedisMgr {
	mgr := &RedisMgr{
		addrs:    make([]string, len(addrs)),
		conf:     pc,
		pools:    make([]*redis.Pool, len(addrs)),
		breakers: make([]*gobreaker.CircuitBreaker, len(addrs)),
	}
	copy(mgr.addrs, addrs)

	for idx, addr := range mgr.addrs {
		log.Info(mctx, "creating redis pool", "redis_server", addr)

		mgr.breakers[idx] = gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:        fmt.Sprint("circuit breaker redis-", addr),
			MaxRequests: 10,
			Interval:    time.Duration(5 * time.Second),
			Timeout:     time.Duration(10 * time.Second),
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				return counts.TotalFailures >= uint32(20)
			},
			OnStateChange: func(name string, from, to gobreaker.State) {
				log.Info(mctx, "circuit breaker state changed", "name", name, "from", from, "to", to)
			},
		})

		mgr.pools[idx] = &redis.Pool{
			MaxIdle:     mgr.conf.MaxIdle,
			MaxActive:   mgr.conf.MaxActive,
			IdleTimeout: mgr.conf.IdleTimeout,
			Wait:        mgr.conf.Wait,
			Dial: func() (redis.Conn, error) {
				c, err := redis.DialTimeout("tcp", addr, mgr.conf.ConnectTimeout, mgr.conf.ReadTimeout, mgr.conf.WriteTimeout)
				if err != nil {
					log.Error(mctx, "dail to server", "redis_server", addr, "error", err)
				}
				return c, err
			},
		}
	}

	return mgr
}

/*************************************************
Description: 初始化redis连接池管理器
Input:
Output:
Return:
Others:
*************************************************/
func Init() {
	mgr := newRedisConnMgr(GetAddrs(), GetPoolsConfig())
	gMgr.Store(mgr)
}

/*************************************************
Description: 清理redis连接池管理器
Input:
Output:
Return:
Others:
*************************************************/
func Uninit() {
	mgr := gMgr.Load().(*RedisMgr)
	mgr.uninit()
}

/*************************************************
Description: 连接池参数更新回调，按照新的参数初始化新的连接池，并销毁旧的连接池
Input:
	v		 连接池的新配置参数
Output:
Return:
Others:
*************************************************/
func OnPoolsConfigUpdate(v interface{}) {
	poolsConf := v.(*PoolsConfig)
	oldMgr := gMgr.Load().(*RedisMgr)
	newMgr := newRedisConnMgr(oldMgr.addrs, poolsConf)
	gMgr.Store(newMgr)
	oldMgr.uninit()
}

/*************************************************
Description: 随机获取一个redis连接，可能是不同的接入地址
Input:
Output:
Return:		 redis连接
Others:
*************************************************/
func GetConn() (c redis.Conn) {
	c, _ = getConnWithCircuitBreaker()
	return
}

func getConnWithCircuitBreaker() (redis.Conn, *gobreaker.CircuitBreaker) {
	mgr := gMgr.Load().(*RedisMgr)
	idx := rand.Intn(len(mgr.pools))
	return mgr.pools[idx].Get(), mgr.breakers[idx]
}

/*************************************************
Description: 执行redis命令，并按照指定的retry策略重试
Input:
	strategy	retry策略
	cmd			redis命令
	args		redis命令的参数
Output:
Return:
	成功时	 redis结果，nil
	失败时	 nil，error
Others:
*************************************************/
func TryRedisCmd(ctx context.Context, strategy retry.ResettableStrategy, cmd string, args ...interface{}) (reply interface{}, err error) {
	if strategy == nil {
		strategy = defaultRetryStrategy()
	}

	retry.DoWithReset(ctx, strategy, func() bool {
		c, breaker := getConnWithCircuitBreaker()
		defer func() {
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

/*************************************************
Description: 执行redis的lua脚本，并按照指定的策略重试
Input:
	strategy	retry策略
	script		redis的lua脚本对象
	keysAndArgs	lua脚本的keys和args参数
Output:
Return:
	成功时	 redis结果，nil
	失败时	 nil，error
Others:
*************************************************/
func TryRedisScript(ctx context.Context, strategy retry.ResettableStrategy, script *redis.Script, keysAndArgs ...interface{}) (reply interface{}, err error) {
	if strategy == nil {
		strategy = defaultRetryStrategy()
	}

	retry.DoWithReset(ctx, strategy, func() bool {
		c, breaker := getConnWithCircuitBreaker()
		defer func() {
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
