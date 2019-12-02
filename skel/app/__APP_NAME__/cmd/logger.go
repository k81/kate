package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogFile struct {
	Level zapcore.Level
	Path  string
}

type LoggerConfig struct {
	// Level is the minimum enabled logging level. Note that this is a dynamic
	// level, so calling Config.Level.SetLevel will atomically change the log
	// level of all loggers descended from this config.
	Level zap.AtomicLevel `json:"level" yaml:"level"`
	// Development puts the logger in development mode, which changes the
	// behavior of DPanicLevel and takes stacktraces more liberally.
	Development bool `json:"development" yaml:"development"`
	// DisableCaller stops annotating logs with the calling function's file
	// name and line number. By default, all logs are annotated.
	DisableCaller bool `json:"disableCaller" yaml:"disableCaller"`
	// DisableStacktrace completely disables automatic stacktrace capturing. By
	// default, stacktraces are captured for WarnLevel and above logs in
	// development and ErrorLevel and above in production.
	DisableStacktrace bool `json:"disableStacktrace" yaml:"disableStacktrace"`
	// Sampling sets a sampling policy. A nil SamplingConfig disables sampling.
	Sampling *SamplingConfig `json:"sampling" yaml:"sampling"`
	// Encoding sets the logger's encoding. Valid values are "json" and
	// "console", as well as any third-party encodings registered via
	// RegisterEncoder.
	Encoding string `json:"encoding" yaml:"encoding"`
	// EncoderConfig sets options for the chosen encoder. See
	// zapcore.EncoderConfig for details.
	EncoderConfig zapcore.EncoderConfig `json:"encoderConfig" yaml:"encoderConfig"`
	// OutputPaths is a list of URLs or file paths to write logging output to.
	// See Open for details.
	OutputPaths []LogFile `json:"outputPaths" yaml:"outputPaths"`
	// ErrorOutputPaths is a list of URLs to write internal logger errors to.
	// The default is standard error.
	//
	// Note that this setting only affects internal errors; for sample code that
	// sends error-level logs to a different location from info- and debug-level
	// logs, see the package-level AdvancedConfiguration example.
	ErrorOutputPaths []string `json:"errorOutputPaths" yaml:"errorOutputPaths"`
	// InitialFields is a collection of fields to add to the root logger.
	InitialFields map[string]interface{} `json:"initialFields" yaml:"initialFields"`
}

func (cfg *LoggerConfig) Load(filePath string) (err error) {
	var content []byte
	if content, err = ioutil.ReadFile(filePath); err != nil {
		return fmt.Errorf("read file: file=%v, error=%w", filePath, err)
	}

	if err = json.Unmarshal(content, cfg); err != nil {
		return fmt.Errorf("unmarshal file: file=%v, error=%w", filePath, err)
	}
	return nil
}

func (cfg *LoggerConfig) Build() (*zap.Logger, error) {
	cores := make([]zapcore.Core, 0, len(cfg.OutputPaths))
	closeSinks := make([]func(), 0, len(cfg.OutputPaths))
	for _, file := range cfg.OutputPaths {
		var (
			encoder      = zapcore.NewJSONEncoder(cfg.EncoderConfig)
			levelEnabler = cfg.getLevelEnabler(file.Level)
		)

		sink, closeSink, err := zap.Open(file.Path)
		if err != nil {
			cfg.cleanUp(closeSinks)
			return nil, fmt.Errorf("failed to open output: file=%s, error=%w", file.Path, err)
		}
		closeSinks = append(closeSinks, closeSink)

		core := zapcore.NewCore(encoder, sink, levelEnabler)
		cores = append(cores, core)
	}

	errSink, _, err := zap.Open(cfg.ErrorOutputPaths...)
	if err != nil {
		cfg.cleanUp(closeSinks)
		return nil, fmt.Errorf("failed to open error output: files=%v, error=%w", cfg.ErrorOutputPaths, err)
	}

	options := cfg.buildOptions(errSink)
	return zap.New(zapcore.NewTee(cores...), options...), nil
}

func (cfg *LoggerConfig) buildOptions(errSink zapcore.WriteSyncer) []zap.Option {
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

func (cfg *LoggerConfig) getLevelEnabler(level zapcore.Level) zapcore.LevelEnabler {
	return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == level
	})
}

func (cfg *LoggerConfig) cleanUp(closeFuncs []func()) {
	for _, f := range closeFuncs {
		f()
	}
}
