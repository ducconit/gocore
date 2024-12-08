package logger

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

// Level represents the severity level of a log message
type Level = zapcore.Level

// These constants define log levels using zap's levels
const (
	DebugLevel  = zapcore.DebugLevel  // -1
	InfoLevel   = zapcore.InfoLevel   // 0
	WarnLevel   = zapcore.WarnLevel   // 1
	ErrorLevel  = zapcore.ErrorLevel  // 2
	DPanicLevel = zapcore.DPanicLevel // 3
	PanicLevel  = zapcore.PanicLevel  // 4
	FatalLevel  = zapcore.FatalLevel  // 5
)

// ParseLevel converts a level string to a Level value
func ParseLevel(levelStr string) Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn", "warning":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "dpanic":
		return DPanicLevel
	case "panic":
		return PanicLevel
	case "fatal":
		return FatalLevel
	default:
		return InfoLevel // default to info
	}
}
