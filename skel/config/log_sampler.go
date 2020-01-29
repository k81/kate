package config

import "time"

type LogSamplerConfig struct {
	Enabled    bool
	Tick       time.Duration
	First      int
	ThereAfter int
}
