package config

import "gopkg.in/ini.v1"

// Profiling is the profiling config instance
var Profiling = &ProfilingConfig{}

// ProfilingConfig defines the profiling config
type ProfilingConfig struct {
	Enabled bool
	Port    int
}

// SectionName implements the `Config.SectionName()` method
func (conf *ProfilingConfig) SectionName() string {
	return "profiling"
}

// Load implements the `Config.Load()` method
func (conf *ProfilingConfig) Load(section *ini.Section) error {
	conf.Enabled = section.Key("enabled").MustBool(true)
	conf.Port = section.Key("port").MustInt(18000)
	return nil
}
