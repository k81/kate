package profiling

import (
	"context"
	"fmt"
	"net/http"
	"time"

	// register the pprof handler
	_ "net/http/pprof"

	"github.com/k81/kate/utils"
	"github.com/k81/log"
)

var (
	mctx = log.WithContext(context.Background(), "module", "profiling")
	addr string
)

// Start start the http pprof server
func Start(port int) {
	go loop(port)
}

func loop(port int) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal(mctx, "got panic", "error", r, "stack", utils.GetPanicStack())
		}
	}()
	// delay to avoid listen addr conflict with parent process
	time.Sleep(5 * time.Second)

	addr = fmt.Sprint("0.0.0.0:", port)

	log.Info(mctx, "starting", "addr", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Error(mctx, "serve http profiling", "addr", addr, "error", err)
	}
}
