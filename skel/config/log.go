package config

import (
	"fmt"
	"path"

	"github.com/k81/kate/app"
	"gopkg.in/ini.v1"
)

// Log is the log config instance
var Log = &LogConfig{}

// LogConfig defines the Log config
type LogConfig struct {
	Level   string
	LogFile string
	ErrFile string
}

// SectionName implements the `Config.SectionName()` method
func (conf *LogConfig) SectionName() string {
	return "log"
}

// Load implements the `Config.Load()` method
func (conf *LogConfig) Load(section *ini.Section) error {
	var (
		defaultLogFile = path.Join(app.GetHomeDir(), "log", fmt.Sprintf("%s.log", app.GetName()))
		defaultErrFile = path.Join(app.GetHomeDir(), "log", fmt.Sprintf("%s.log.wf", app.GetName()))
	)

	conf.Level = section.Key("level").MustString("DEBUG")
	conf.LogFile = section.Key("log_file").MustString(defaultLogFile)
	conf.ErrFile = section.Key("err_file").MustString(defaultErrFile)

	return nil
}
