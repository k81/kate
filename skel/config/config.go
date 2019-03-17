package config

import (
	"context"

	"github.com/k81/kate/configer"
	"github.com/k81/log"
	"gopkg.in/ini.v1"
)

var (
	mctx = log.WithContext(context.Background(), "module", "config")
)

// Config defines the config interface
type Config interface {
	// SectionName return the section name
	SectionName() string

	// Load load the config in the section specified in `SectionName()`
	Load(*ini.Section) error
}

var Configer configer.Configer = &iniConfiger{}
