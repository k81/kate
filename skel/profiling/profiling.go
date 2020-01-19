package profiling

import (
	"fmt"
	"net/http"
	"time"

	// register the pprof handler
	_ "net/http/pprof"

	"go.uber.org/zap"
)

var addr string
var logger *zap.Logger

// Start start the http pprof service
func Start(port int, l *zap.Logger) {
	logger = l

	go loop(port)
}

func loop(port int) {
	defer func() {
		if r := recover(); r != nil {
			logger.Fatal("got panic", zap.Any("error", r), zap.Stack("stack"))
		}
	}()
	// delay to avoid listen addr conflict with parent process
	time.Sleep(5 * time.Second)

	addr = fmt.Sprint("0.0.0.0:", port)

	logger.Info("profiling started listening", zap.String("addr", addr))

	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Error("serve http profiling", zap.String("addr", addr), zap.Error(err))
	}
}
