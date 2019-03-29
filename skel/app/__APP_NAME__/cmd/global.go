package cmd

import (
	"context"

	"github.com/k81/log"
)

type globalFlags struct {
	Debug      bool
	ConfigFile string
}

var (
	GlobalFlags = &globalFlags{}
	mctx        = log.WithContext(context.Background(), "module", "cmd")
)
