package log

import (
	"context"
	"os"
)

const (
	keyLogContext = int(0xf7f7f7f7)
)

var (
	mctx                    = SetContext(context.Background(), "module", "log")
	RootContext *LogContext = NewLogContext()
	logger      Logger      = NewStdLogger(KVFormatter)
	level                   = TraceLevel
)

func SetLogger(l Logger) {
	logger = l
}

func GetLogger() Logger {
	return logger
}

func SetLevelByName(lvlName LevelName) {
	level = lvlName.ToLevel()
}

func GetLevel() Level {
	return level
}

func Enabled(lvl Level) bool {
	return lvl >= level
}

func getLogContext(ctx context.Context) *LogContext {
	if logctx, ok := ctx.Value(keyLogContext).(*LogContext); ok {
		return logctx
	}
	return RootContext
}

// 为当前context绑定日志上下文变量
func SetContext(ctx context.Context, keyvals ...interface{}) context.Context {
	logCtx := getLogContext(ctx).With(keyvals...)
	return context.WithValue(ctx, keyLogContext, logCtx)
}

func Trace(ctx context.Context, msg string, keyvals ...interface{}) {
	if Enabled(TraceLevel) {
		logger.Log(getLogContext(ctx).newEntry(TraceLevel, msg, keyvals))
	}
}

func Debug(ctx context.Context, msg string, keyvals ...interface{}) {
	if Enabled(DebugLevel) {
		logger.Log(getLogContext(ctx).newEntry(DebugLevel, msg, keyvals))
	}
}

func Info(ctx context.Context, msg string, keyvals ...interface{}) {
	if Enabled(InfoLevel) {
		logger.Log(getLogContext(ctx).newEntry(InfoLevel, msg, keyvals))
	}
}

func Warning(ctx context.Context, msg string, keyvals ...interface{}) {
	if Enabled(WarningLevel) {
		logger.Log(getLogContext(ctx).newEntry(WarningLevel, msg, keyvals))
	}
}

func Error(ctx context.Context, msg string, keyvals ...interface{}) {
	if Enabled(ErrorLevel) {
		logger.Log(getLogContext(ctx).newEntry(ErrorLevel, msg, keyvals))
	}
}

func Fatal(ctx context.Context, msg string, keyvals ...interface{}) {
	if Enabled(FatalLevel) {
		logger.Log(getLogContext(ctx).newEntry(FatalLevel, msg, keyvals))
	}
	os.Exit(-1)
}
