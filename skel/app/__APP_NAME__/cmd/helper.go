package cmd

import (
	"path"

	"github.com/k81/kate/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"__PACKAGE_NAME__/config"
)

func initLogger() *zap.Logger {
	enc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	core := zapcore.NewTee(
		log.MustNewCore(zapcore.DebugLevel, path.Join(config.Main.LogDir, "debug.log"), enc),
		log.MustNewCore(zapcore.InfoLevel, path.Join(config.Main.LogDir, "info.log"), enc),
		log.MustNewCore(zapcore.WarnLevel, path.Join(config.Main.LogDir, "warn.log"), enc),
		log.MustNewCore(zapcore.ErrorLevel, path.Join(config.Main.LogDir, "error.log"), enc),
		log.MustNewCore(zapcore.FatalLevel, path.Join(config.Main.LogDir, "fatal.log"), enc),
	)

	opts := []zap.Option{
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCaller(),
	}

	logger := zap.New(core, opts...)
	zap.ReplaceGlobals(logger)

	return logger
}

func initDevLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}
