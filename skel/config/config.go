package config

import (
	"context"
	"os"

	"github.com/k81/log"
	"gopkg.in/ini.v1"
)

var (
	mctx = log.WithContext(context.Background(), "module", "config")
	cfg  *ini.File
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
	var err error
	if cfg, err = ini.Load(file); err != nil {
		return err
	}

	configs := []Config{
		Log,
		Profiling,
		MySQL,
		Redis,
		HTTP,
	}

	for _, config := range configs {
		section := cfg.Section(config.SectionName())
		if err = config.Load(section); err != nil {
			log.Fatal(mctx, "load section failed",
				"section", loader.SectionName(),
				"error", err)
			os.Exit(1)
		}
	}
	log.Info(mctx, "config loaded successfully")
}
