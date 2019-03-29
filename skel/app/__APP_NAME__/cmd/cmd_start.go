package cmd

import (
	"fmt"
	"os"

	"github.com/k81/kate/app"
	"github.com/k81/kate/redismgr"
	"github.com/k81/kate/utils"
	"github.com/k81/log"
	"github.com/spf13/cobra"

	"__PROJECT_DIR__/config"
	"__PROJECT_DIR__/httpsrv"
	"__PROJECT_DIR__/model"
	"__PROJECT_DIR__/profiling"
)

func NewStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start server",
		Run:   startCmdFunc,
	}
	return cmd
}

func startCmdFunc(cmd *cobra.Command, args []string) {
	var (
		logAppender *log.FileAppender
		errAppender *log.FileAppender
		err         error
	)

	// load config
	if err = config.Load(GlobalFlags.ConfigFile); err != nil {
		log.Fatal(mctx, "load config failed", "file", GlobalFlags.ConfigFile, "error", err)
		return
	}

	// create log file
	if logAppender, err = log.NewFileAppender(log.LevelMask, config.Log.LogFile, log.PipeKVFormatter); err != nil {
		os.Exit(1)
	}
	logAppender.DisableLock()

	// create err file
	if errAppender, err = log.NewFileAppender(log.LevelError|log.LevelFatal, config.Log.ErrFile, log.PipeKVFormatter); err != nil {
		os.Exit(1)
	}
	errAppender.DisableLock()

	// set logger
	log.SetLogger(log.NewLogger(logAppender, errAppender))

	// update pid
	app.UpdatePIDFile(config.Main.PIDFile)

	defer func() {
		if r := recover(); r != nil {
			log.Fatal(mctx, "panic", "error", r, "stack", utils.GetPanicStack())
		}

		log.Info(mctx, "shutting down ...")

		redismgr.Uninit()
		app.RemovePIDFile()
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
