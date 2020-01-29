package config

import (
	"time"

	"gopkg.in/ini.v1"
)

// GRPC is the grpc config instance
var GRPC = &GRPCConfig{}

// GRPCConfig defines the grpc config
type GRPCConfig struct {
	Addr       string
	LogFile    string
	LogSampler LogSamplerConfig
}

// SectionName implements the `Config.SectionName()` method
func (conf *GRPCConfig) SectionName() string {
	return "grpc"
}

// Load implements the `Config.Load()` method
func (conf *GRPCConfig) Load(section *ini.Section) error {
	conf.Addr = section.Key("addr").MustString(":9090")
	conf.LogFile = section.Key("log_file").MustString("__APP_NAME__.access")
	conf.LogSampler.Enabled = section.Key("log_sampler_enabled").MustBool(false)
	conf.LogSampler.Tick = section.Key("log_sampler_tick").MustDuration(time.Second)
	conf.LogSampler.First = section.Key("log_sampler_first").MustInt(100)
	conf.LogSampler.ThereAfter = section.Key("log_sampler_thereafter").MustInt(10000)
	return nil
}
