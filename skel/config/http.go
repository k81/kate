package config

import (
	"time"

	"gopkg.in/ini.v1"
)

// HTTP is the http config instance
var HTTP = &HTTPConfig{}

// HTTPConfig defines the HTTP config
type HTTPConfig struct {
	Addr           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
	MaxBodyBytes   int64
}

// SectionName implements the `Config.SectionName()` method
func (conf *HTTPConfig) SectionName() string {
	return "http"
}

// Load implements the `Config.Load()` method
func (conf *HTTPConfig) Load(section *ini.Section) error {
	conf.Addr = section.Key("addr").MustString(":8080")
	conf.ReadTimeout = section.Key("read_timeout").MustDuration(2000 * time.Millisecond)
	conf.WriteTimeout = section.Key("write_timeout").MustDuration(0)
	conf.MaxHeaderBytes = section.Key("max_header_bytes").MustInt(1048576)
	conf.MaxBodyBytes = section.Key("max_body_bytes").MustInt64(1073741824)
	return nil
}
