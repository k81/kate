package config

import (
	"fmt"

	"gopkg.in/ini.v1"
)

// Config defines the config interface
type Config interface {
	// SectionName return the section name
	SectionName() string

	// Load load the config in the section specified in `SectionName()`
	Load(*ini.Section) error
}

// Load load all configs
func Load(file string) error {
	var (
		iniFile *ini.File
		err     error
	)

	if iniFile, err = ini.Load(file); err != nil {
		return fmt.Errorf("load config: %v", err)
	}

	configs := []Config{
		Main,
		Log,
		Profiling,
		DB,
		Redis,
		HTTP,
	}

	for _, config := range configs {
		section := iniFile.Section(config.SectionName())
		if err = config.Load(section); err != nil {
			return fmt.Errorf("load config: section=%v, error=%v", config.SectionName(), err)
		}
	}
	return nil
}
