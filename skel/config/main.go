package config

import (
	"fmt"
	"path"

	"github.com/k81/kate/app"
	"gopkg.in/ini.v1"
)

// Main is the log config instance
var Main = &MainConfig{}

// MainConfig defines the Main config
type MainConfig struct {
	PIDFile string
}

// SectionName implements the `Config.SectionName()` method
func (conf *MainConfig) SectionName() string {
	return "main"
}

// Load implements the `Config.Load()` method
func (conf *MainConfig) Load(section *ini.Section) error {
	defaultPIDFile := path.Join(app.GetHomeDir(), "run", fmt.Sprintf("%s.pid", app.GetName()))
	conf.PIDFile = section.Key("pid_file").MustString(defaultPIDFile)
	return nil
}
