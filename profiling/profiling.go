package profiling

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/k81/log"
)

var (
	logger *log.Logger
	mctx   = context.Background()
	addr   string
)

func Enabled() bool {
	return GetConfig().Enabled
}

func Start() {
	logger = log.With("module", "profiling")
	go loop()
}

func loop() {
	var err error

	// delay to avoid listen addr conflict with parent process
	time.Sleep(10 * time.Second)

	addr = fmt.Sprint("0.0.0.0:", GetConfig().Port)

	logger.Info(mctx, "starting", "addr", addr)

	if err = http.ListenAndServe(addr, nil); err != nil {
		logger.Error(mctx, "serve http profiling", "addr", addr, "error", err)
	}
}
