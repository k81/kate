package config

import (
	"time"

	"github.com/k81/kate/redismgr"
	"github.com/k81/kate/utils"
	"gopkg.in/ini.v1"
)

// Redis is the redis config instance
var Redis = &RedisConfig{&redismgr.RedisConfig{}}

// RedisConfig defines the redis config
type RedisConfig struct {
	*redismgr.RedisConfig
}

// SectionName implements the `Config.SectionName()` method
func (conf *RedisConfig) SectionName() string {
	return "redis"
}

// Load implements the `Config.Load()` method
func (conf *RedisConfig) Load(section *ini.Section) error {
	addrs := section.Key("addrs").MustString("127.0.0.1:6379")
	conf.Addrs = utils.Split(addrs, ",")
	conf.MaxIdle = section.Key("max_idle").MustInt(10)
	conf.MaxActive = section.Key("max_active").MustInt(50)
	conf.ConnectTimeout = section.Key("connect_timeout").MustDuration(time.Second)
	conf.ReadTimeout = section.Key("read_timeout").MustDuration(500 * time.Millisecond)
	conf.WriteTimeout = section.Key("write_timeout").MustDuration(500 * time.Millisecond)
	conf.IdleTimeout = section.Key("idle_timeout").MustDuration(30 * time.Second)
	conf.Wait = section.Key("wait").MustBool(true)
	return nil
}
