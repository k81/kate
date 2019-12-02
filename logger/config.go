package logger

import (
	"path"
	"sort"
	"time"

	"github.com/k81/kate/app"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	// Level is the minimum enabled logging level. Note that this is a dynamic
	// level, so calling Config.Level.SetLevel will atomically change the logger
	// level of all loggers descended from this config.
	Level zap.AtomicLevel `json:"level,omitempty" yaml:"level,omitempty"`
	// Development puts the logger in development mode, which changes the
	// behavior of DPanicLevel and takes stacktraces more liberally.
	Development bool `json:"development,omitempty" yaml:"development,omitempty"`
	// DisableCaller stops annotating logs with the calling function's file
	// name and line number. By default, all logs are annotated.
	DisableCaller bool `json:"disableCaller,omitempty" yaml:"disableCaller,omitempty"`
	// DisableStacktrace completely disables automatic stacktrace capturing. By
	// default, stacktraces are captured for WarnLevel and above logs in
	// development and ErrorLevel and above in production.
	DisableStacktrace bool `json:"disableStacktrace,omitempty" yaml:"disableStacktrace,omitempty"`
	// Sampling sets a sampling policy. A nil SamplingConfig disables sampling.
	Sampling *zap.SamplingConfig `json:"sampling,omitempty" yaml:"sampling,omitempty"`
	// OutputPaths is a list of URLs or file paths to write logging output to.
	// See Open for details.
	OutputPaths map[zapcore.Level]string `json:"outputPaths,omitempty" yaml:"outputPaths,omitempty"`
	// ErrorOutputPaths is a list of URLs to write internal logger errors to.
	// The default is standard error.
	//
	// Note that this setting only affects internal errors; for sample code that
	// sends error-level logs to a different location from info- and debug-level
	// logs, see the package-level AdvancedConfiguration example.
	ErrorOutputPaths []string `json:"errorOutputPaths,omitempty" yaml:"errorOutputPaths,omitempty"`
	// InitialFields is a collection of fields to add to the root logger.
	InitialFields map[string]interface{} `json:"initialFields,omitempty" yaml:"initialFields,omitempty"`
}

func (cfg *Config) Build() (*zap.Logger, error) {
	cores := make([]zapcore.Core, 0, len(cfg.OutputPaths))
	closeSinks := make([]func(), 0, len(cfg.OutputPaths))
	for level, filePath := range cfg.OutputPaths {
		var (
			encoder      = NewEncoder()
			levelEnabler = cfg.getLevelEnabler(level)
		)

		if !path.IsAbs(filePath) {
			filePath = path.Join(app.GetHomeDir(), "log", filePath)
		}
		sink, closeSink, err := zap.Open(filePath)
		if err != nil {
			cfg.cleanUp(closeSinks)
			return nil, err
		}
		closeSinks = append(closeSinks, closeSink)

		core := zapcore.NewCore(encoder, sink, levelEnabler)
		cores = append(cores, core)
	}

	errSink, _, err := zap.Open(cfg.ErrorOutputPaths...)
	if err != nil {
		cfg.cleanUp(closeSinks)
		return nil, err
	}

	options := cfg.buildOptions(errSink)
	return zap.New(zapcore.NewTee(cores...), options...), nil
}

func (cfg *Config) buildOptions(errSink zapcore.WriteSyncer) []zap.Option {
	opts := []zap.Option{zap.ErrorOutput(errSink)}

	if cfg.Development {
		opts = append(opts, zap.Development())
	}

	if !cfg.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}

	stackLevel := zap.ErrorLevel
	if cfg.Development {
		stackLevel = zap.WarnLevel
	}
	if !cfg.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}

	if cfg.Sampling != nil {
		opts = append(opts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewSampler(core, time.Second, int(cfg.Sampling.Initial), int(cfg.Sampling.Thereafter))
		}))
	}

	if len(cfg.InitialFields) > 0 {
		fs := make([]zap.Field, 0, len(cfg.InitialFields))
		keys := make([]string, 0, len(cfg.InitialFields))
		for k := range cfg.InitialFields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fs = append(fs, zap.Any(k, cfg.InitialFields[k]))
		}
		opts = append(opts, zap.Fields(fs...))
	}

	return opts
}

func (cfg *Config) getLevelEnabler(level zapcore.Level) zapcore.LevelEnabler {
	return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == level
	})
}

func (cfg *Config) cleanUp(closeFuncs []func()) {
	for _, f := range closeFuncs {
		f()
	}
}
