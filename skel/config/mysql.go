package config

import (
	"time"

	"gopkg.in/ini.v1"
)

// MySQL is the mysql config instance
var MySQL = &MySQLConfig{}

// MySQLConfig defines the mysql config
type MySQLConfig struct {
	DataSource      string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	DebugSQL        bool
}

// SectionName implements the `Config.SectionName()` method
func (conf *MySQLConfig) SectionName() string {
	return "mysql"
}

// Load implements the `Config.Load()` method
func (conf *MySQLConfig) Load(section *ini.Section) error {
	conf.DataSource = section.Key("data_source").String()
	conf.MaxIdleConns = section.Key("max_idle_conns").MustInt(20)
	conf.MaxOpenConns = section.Key("max_open_conns").MustInt(60)
	conf.ConnMaxLifetime = section.Key("conn_max_lifetime").MustDuration(60 * time.Second)
	conf.DebugSQL = section.Key("debug_sql").MustBool(false)
	return nil
}
