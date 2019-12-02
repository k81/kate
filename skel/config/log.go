package config

import (
	"gopkg.in/ini.v1"
)

// Log is the log config instance
var Log = &LogConfig{}

// LogConfig defines the Log config
type LogConfig struct {
	Level string
}

// SectionName implements the `Config.SectionName()` method
func (conf *LogConfig) SectionName() string {
	return "log"
}

// Load implements the `Config.Load()` method
func (conf *LogConfig) Load(section *ini.Section) error {
	conf.Level = section.Key("level").MustString("DEBUG")

	return nil
}
