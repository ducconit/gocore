package logger

import "go.uber.org/zap"

var (
	defaultLogger *Logger
)

func Info(msg string, fields ...zap.Field) {
	if defaultLogger != nil {
		defaultLogger.Info(msg, fields...)
	}
}

func Debug(msg string, fields ...zap.Field) {
	if defaultLogger != nil {
		defaultLogger.Debug(msg, fields...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	if defaultLogger != nil {
		defaultLogger.Warn(msg, fields...)
	}
}

func Error(msg string, fields ...zap.Field) {
	if defaultLogger != nil {
		defaultLogger.Error(msg, fields...)
	}
}

func Panic(msg string, fields ...zap.Field) {
	if defaultLogger != nil {
		defaultLogger.Panic(msg, fields...)
	}
}

func Instance() *Logger {
	return defaultLogger
}
