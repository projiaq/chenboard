package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Init 初始化日志
func Init(level, format string) error {
	var cfg zap.Config

	if format == "json" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	// 设置日志级别
	switch level {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	var err error
	log, err = cfg.Build()
	if err != nil {
		return err
	}

	return nil
}

// Get 获取日志实例
func Get() *zap.Logger {
	if log == nil {
		log, _ = zap.NewProduction()
	}
	return log
}

// Debug 调试日志
func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

// Info 信息日志
func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

// Warn 警告日志
func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

// Error 错误日志
func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

// Fatal 致命错误日志
func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
}

// Sync 同步日志
func Sync() error {
	if log != nil {
		return log.Sync()
	}
	return nil
}
