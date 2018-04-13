package profiling

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/k81/kate/log"
)

var (
	mctx    = log.SetContext(context.Background(), "module", "profiling")
	started bool
	addr    string
)

func Enabled() bool {
	return GetConfig().Enabled
}

func Start() {
	go loop()
}

func loop() {
	var err error

	// delay to avoid listen addr conflict with parent process
	time.Sleep(10 * time.Second)

	addr = fmt.Sprint("0.0.0.0:", GetConfig().Port)

	log.Info(mctx, "starting", "addr", addr)

	if err = http.ListenAndServe(addr, nil); err != nil {
		log.Error(mctx, "serve http profiling", "addr", addr, "error", err)
	}
}
