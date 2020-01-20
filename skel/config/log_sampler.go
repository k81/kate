package config

import "time"

type LogSamplerConfig struct {
	Tick       time.Duration
	First      int
	ThereAfter int
}
