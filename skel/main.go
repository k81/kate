package main

import (
	"context"
	"fmt"

	"github.com/facebookgo/grace/gracehttp"

	"github.com/k81/kate/app"
	"github.com/k81/kate/log"
	"github.com/k81/kate/utils"

	"__PROJECT_DIR__/httpsrv"
	"__PROJECT_DIR__/repo"
)

var (
	mctx = log.SetContext(context.Background(), "module", "main")
)

func main() {
	app.Setup()

	defer func() {
		if r := recover(); r != nil {
			log.Fatal(mctx, "panic", "error", r, "stack", utils.GetPanicStack())
		}

		log.Info(mctx, "shutting down ...")

		app.Cleanup()
		log.Info(mctx, fmt.Sprintf("%s stopped", app.GetName()))
	}()

	repo.Init()

	log.Info(mctx, fmt.Sprintf("%s started", app.GetName()), "listen_addr", httpsrv.GetListenAddr())

	if err := gracehttp.Serve(httpsrv.GetServer()); err != nil {
		log.Error(mctx, "serve error", "error", err)
	}
}
