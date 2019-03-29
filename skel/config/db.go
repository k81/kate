package config

import (
	"time"

	"gopkg.in/ini.v1"
)

// DB is the mysql config instance
var DB = &DBConfig{}

// DBConfig defines the mysql config
type DBConfig struct {
	DataSource      string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

// SectionName implements the `Config.SectionName()` method
func (conf *DBConfig) SectionName() string {
	return "mysql"
}

// Load implements the `Config.Load()` method
func (conf *DBConfig) Load(section *ini.Section) error {
	conf.DataSource = section.Key("data_source").String()
	conf.MaxIdleConns = section.Key("max_idle_conns").MustInt(20)
	conf.MaxOpenConns = section.Key("max_open_conns").MustInt(60)
	conf.ConnMaxLifetime = section.Key("conn_max_lifetime").MustDuration(60 * time.Second)
	return nil
}
