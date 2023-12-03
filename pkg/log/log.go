package log

import (
	"os"
	"scaffolding/pkg/config"
	"scaffolding/pkg/internal"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

const (
	_jsonFormat    = "json"
	_consoleFormat = "console"
)

func InitLog() {
	Logger = zap.New(logCore())
	if config.Conf.Zap.ShowLine {
		Logger = Logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
	}
}

func logCore() zapcore.Core {
	writer := logWriter()
	if config.Conf.Zap.LogInFile && config.Conf.Zap.LogInConsole {
		return zapcore.NewTee(
			zapcore.NewCore(encoder(), zapcore.AddSync(os.Stdout), internal.AtomicLevel),
			zapcore.NewCore(encoder(), writer, internal.AtomicLevel))
	} else if config.Conf.Zap.LogInFile {
		return zapcore.NewCore(encoder(), writer, internal.AtomicLevel)
	} else {
		// 默认写入控制台
		return zapcore.NewCore(encoder(), zapcore.AddSync(os.Stdout), internal.AtomicLevel)
	}
}

func encoder() zapcore.Encoder {
	switch strings.ToLower(config.Conf.Zap.Format) {
	case _consoleFormat:
		return zapcore.NewConsoleEncoder(encoderConfig())
	case _jsonFormat:
		return zapcore.NewJSONEncoder(encoderConfig())
	default:
		return zapcore.NewJSONEncoder(encoderConfig())
	}
}

func encoderConfig() zapcore.EncoderConfig {
	zap.NewDevelopmentEncoderConfig()
	cfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "ts",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "fn",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    nil,
		EncodeTime:     zapcore.TimeEncoderOfLayout(time.DateTime),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	if strings.ToLower(config.Conf.Zap.Format) == _consoleFormat {
		cfg.ConsoleSeparator = " "
	}
	switch config.Conf.Zap.EncodeLevel {
	case "LowercaseLevelEncoder": // 小写编码器(默认)
		cfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	case "LowercaseColorLevelEncoder": // 小写编码器带颜色
		cfg.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case "CapitalLevelEncoder": // 大写编码器
		cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	case "CapitalColorLevelEncoder": // 大写编码器带颜色
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		cfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	}

	return cfg
}

func logWriter() zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.Conf.Lumberjack.FileName,
		MaxSize:    config.Conf.Lumberjack.MaxSize,
		MaxAge:     config.Conf.Lumberjack.MaxAge,
		MaxBackups: config.Conf.Lumberjack.MaxBackups,
		LocalTime:  true,
		Compress:   config.Conf.Lumberjack.Compress,
	})
}

func Debug(message string, fields ...zapcore.Field) {
	Logger.Debug(message, fields...)
}

func Info(message string, fields ...zapcore.Field) {
	Logger.Info(message, fields...)
}

func Warn(message string, fields ...zapcore.Field) {
	Logger.Warn(message, fields...)
}

func Error(message string, fields ...zapcore.Field) {
	Logger.Error(message, fields...)
}

func Panic(message string, fields ...zapcore.Field) {
	Logger.Panic(message, fields...)
}
