package log

import (
	"fmt"
	"path"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/k81/kate/app"
)

// MustNewCore create an zapcore.Core instance, exit if error occurred
func MustNewCore(level zapcore.Level, location string, enc zapcore.Encoder) zapcore.Core {
	if !path.IsAbs(location) {
		location = path.Join(app.GetHomeDir(), "log", location)
	}

	writer, err := NewWriter(location)
	if err != nil {
		panic(fmt.Errorf("failed to create file sink: %v, %v", location, err))
	}

	levelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == level
	})

	return zapcore.NewCore(enc, writer, levelEnabler)
}
