package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/k81/kate/app"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"__PROJECT_DIR__/config"
	"__PROJECT_DIR__/httpsrv"
	"__PROJECT_DIR__/profiling"
)

var logger, _ = zap.NewDevelopment()

func NewStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start server",
		Run:   startCmdFunc,
	}
	return cmd
}

func initLog() {
	var (
		cfg     = LoggerConfig{}
		cfgPath = path.Join(app.GetHomeDir(), "conf", "logger.json")
		err     error
	)

	if err = cfg.Load(cfgPath); err != nil {
		fmt.Fprintf(os.Stderr, "failed to load logger config: error=%v", err)
		os.Exit(1)
	}

	if logger, err = cfg.Build(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build logger: error=%v", err)
		os.Exit(1)
	}
}

func startCmdFunc(cmd *cobra.Command, args []string) {
	initLog()

	// load config
	if err := config.Load(GlobalFlags.ConfigFile); err != nil {
		logger.Fatal("load config failed", zap.String("file", GlobalFlags.ConfigFile), zap.Error(err))
	}

	// update pid
	app.UpdatePIDFile(config.Main.PIDFile)

	defer func() {
		if r := recover(); r != nil {
			logger.Fatal("panic", zap.Any("error", r), zap.Stack("stack"))
		}

		logger.Info("shutting down ...")

		app.RemovePIDFile()
		logger.Info(fmt.Sprintf("%s stopped", app.GetName()))
	}()

	if config.Profiling.Enabled {
		profiling.Start(config.Profiling.Port, logger)
	}

	logger.Info(fmt.Sprintf("%s started", app.GetName()))
	httpsrv.ListenAndServe()
}
