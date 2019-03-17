package main

import (
	"context"
	"fmt"

	"github.com/k81/kate/app"
	"github.com/k81/kate/redismgr"
	"github.com/k81/kate/utils"
	"github.com/k81/log"

	"__PROJECT_DIR__/config"
	"__PROJECT_DIR__/httpsrv"
	"__PROJECT_DIR__/model"
	"__PROJECT_DIR__/profiling"
)

var (
	mctx = log.WithContext(context.Background(), "module", "main")
)

func main() {
	app.Setup(config.Configer)

	defer func() {
		if r := recover(); r != nil {
			log.Fatal(mctx, "panic", "error", r, "stack", utils.GetPanicStack())
		}

		log.Info(mctx, "shutting down ...")

		redismgr.Uninit()
		app.Cleanup()
		log.Info(mctx, fmt.Sprintf("%s stopped", app.GetName()))
	}()

	if config.Profiling.Enabled {
		profiling.Start(config.Profiling.Port)
	}

	redismgr.Init(config.Redis.RedisConfig)
	model.Init()

	log.Info(mctx, fmt.Sprintf("%s started", app.GetName()))
	httpsrv.ListenAndServe()
}
