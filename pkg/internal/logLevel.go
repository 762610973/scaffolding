package internal

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

var AtomicLevel = zap.NewAtomicLevel()

func SetLogLevel(l string) {
	var level zapcore.Level
	switch strings.ToLower(l) {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "panic":
		level = zapcore.PanicLevel
	default:
		level = zapcore.InfoLevel
	}

	AtomicLevel.SetLevel(level)
}
