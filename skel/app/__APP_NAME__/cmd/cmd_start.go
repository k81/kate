package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/k81/kate/app"
	"github.com/k81/kate/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"__PROJECT_DIR__/config"
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

func initLog() *zap.Logger {
	loggerCfgPath := path.Join(app.GetHomeDir(), "conf", "logger.json")
	data, err := ioutil.ReadFile(loggerCfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read logger config: %v\n", err)
		os.Exit(1)
	}

	cfg := &logger.Config{}
	if err = json.Unmarshal(data, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "failed to unmarshal logger config: %v\n", err)
		os.Exit(1)
	}

	logger, err := cfg.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "create logger failed: %v\n", err)
		os.Exit(1)
	}
	return logger
}

func startCmdFunc(cmd *cobra.Command, args []string) {
	logger := initLog()

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

		logger.Info("server shutting down ...")

		app.RemovePIDFile()
		logger.Info("server stopped")
	}()

	if config.Profiling.Enabled {
		profiling.Start(config.Profiling.Port, logger)
	}

	logger.Info("server starting")

	// TODO: server start here

	logger.Info("server started", zap.String("k1", "v1"), zap.Any("k2", config.Profiling))
}
