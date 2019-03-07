package config

import (
	"strings"
	"time"

	"github.com/k81/log"
)

// 本地配置
type LocalConfig struct {
	EtcdAddrs     []string      `yaml:"etcd_addrs"`
	EtcdPrefix    string        `yaml:"etcd_prefix"`
	WatchEnabled  bool          `yaml:"watch_enabled"`
	WatchInterval time.Duration `yaml:"watch_interval"`
	UseCacheOnly  bool          `yaml:"use_cache_only"`
}

var (
	Local *LocalConfig = &LocalConfig{}
)

/*************************************************
Description: 从配置文件加载本地配置
Input:
Output:
Return:
Others:
*************************************************/
func (l *LocalConfig) init() {
	var err error

	if err = LoadConfig(configFilePath, Local); err != nil {
		log.Fatal(mctx, "load config file", "config_file", configFilePath, "error", err)
	}

	if Local.WatchInterval <= 0 {
		Local.WatchInterval = 15
	}

	Local.WatchInterval *= time.Second

	log.Info(mctx, "local config",
		"config_file", configFilePath,
		"etcd_addrs", strings.Join(Local.EtcdAddrs, ","),
		"etcd_prefix", Local.EtcdPrefix,
		"watch_enabled", Local.WatchEnabled,
		"watch_interval", Local.WatchInterval,
	)
}
