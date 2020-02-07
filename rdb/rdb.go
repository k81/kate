package rdb

import (
	"time"

	"github.com/go-redis/redis"
)

const (
	// RouteModeMasterOnly only route read-only commands to master node
	RouteModeMasterOnly = "master_only"
	// RouteModeMasterSlaveRandom route read-only commands to both master and slave, using random policy
	RouteModeMasterSlaveRandom = "master_slave_random"
	// RouteModeMasterSlaveLatency route read-only commands to both master and slave, using least latency policy
	RouteModeMasterSlaveLatency = "master_slave_latency"
)

// Client is the client interface for redis db
type Client interface {
	redis.Cmdable
	Do(args ...interface{}) *redis.Cmd
	Process(cmd redis.Cmder) error
	Close() error
}

// Config defines the redis config
type Config struct {
	Addrs              []string
	ClusterEnabled     bool
	ReadOnly           bool
	RouteMode          string
	MaxRedirects       int
	MaxRetries         int
	MinRetryBackoff    time.Duration
	MaxRetryBackoff    time.Duration
	ConnectTimeout     time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	PoolSize           int
	MinIdleConns       int
	MaxConnAge         time.Duration
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration
}

var rdb Client

// Init initialize the redis cluster instance
func Init(conf *Config) {
	if conf.ClusterEnabled {
		rdb = newClusterClient(conf)
	} else {
		rdb = newClient(conf)
	}
}

func newClient(conf *Config) *redis.Client {
	opt := &redis.Options{
		Addr:               conf.Addrs[0],
		MaxRetries:         conf.MaxRetries,
		MinRetryBackoff:    conf.MinRetryBackoff,
		MaxRetryBackoff:    conf.MaxRetryBackoff,
		DialTimeout:        conf.ConnectTimeout,
		ReadTimeout:        conf.ReadTimeout,
		WriteTimeout:       conf.WriteTimeout,
		PoolSize:           conf.PoolSize,
		MinIdleConns:       conf.MinIdleConns,
		MaxConnAge:         conf.MaxConnAge,
		PoolTimeout:        conf.PoolTimeout,
		IdleTimeout:        conf.IdleTimeout,
		IdleCheckFrequency: conf.IdleCheckFrequency,
	}

	return redis.NewClient(opt)
}

func newClusterClient(conf *Config) *redis.ClusterClient {
	opt := &redis.ClusterOptions{
		Addrs:              conf.Addrs,
		MaxRedirects:       conf.MaxRedirects,
		MaxRetries:         conf.MaxRetries,
		MinRetryBackoff:    conf.MinRetryBackoff,
		MaxRetryBackoff:    conf.MaxRetryBackoff,
		DialTimeout:        conf.ConnectTimeout,
		ReadTimeout:        conf.ReadTimeout,
		WriteTimeout:       conf.WriteTimeout,
		PoolSize:           conf.PoolSize,
		MinIdleConns:       conf.MinIdleConns,
		MaxConnAge:         conf.MaxConnAge,
		PoolTimeout:        conf.PoolTimeout,
		IdleTimeout:        conf.IdleTimeout,
		IdleCheckFrequency: conf.IdleCheckFrequency,
	}

	switch conf.RouteMode {
	case RouteModeMasterOnly:
		opt.ReadOnly = false
	case RouteModeMasterSlaveRandom:
		opt.RouteRandomly = true
	case RouteModeMasterSlaveLatency:
		opt.RouteByLatency = true
	}

	return redis.NewClusterClient(opt)
}

// Uninit do the clean up for the global RedisConnectionManager instance
func Uninit() {
	if rdb != nil {
		// nolint: errcheck
		_ = rdb.Close()
	}
}

// Get() return the rdb client instance
func Get() Client {
	return rdb
}
