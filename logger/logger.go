package logger

import (
	"io"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger represents a logger instance
type Logger struct {
	*zap.Logger
	level      Level
	outputs    []io.Writer
	timeFormat string
	mu         sync.RWMutex
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// NewLogger creates a new logger instance
func NewLogger(opts ...Option) *Logger {
	l := &Logger{
		level:      InfoLevel, // default level
		outputs:    make([]io.Writer, 0),
		timeFormat: "2006-01-02 15:04:05",
	}

	// Apply options
	for _, opt := range opts {
		opt(l)
	}

	// Create encoder config
	encConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create cores
	var cores []zapcore.Core
	if len(l.outputs) > 0 {
		enc := zapcore.NewJSONEncoder(encConfig)
		for _, output := range l.outputs {
			core := zapcore.NewCore(enc, zapcore.AddSync(output), l.level)
			cores = append(cores, core)
		}
	}

	// Create logger
	l.Logger = zap.New(zapcore.NewTee(cores...))
	return l
}

// Default returns the default logger instance
func Default() *Logger {
	once.Do(func() {
		defaultLogger = NewLogger(
			WithLevel(InfoLevel),
			WithConsole(),
		)
	})
	return defaultLogger
}

// SetDefault sets the default logger instance
func SetDefault(l *Logger) {
	defaultLogger = l
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetLevelString sets the minimum log level using a string
func (l *Logger) SetLevelString(levelStr string) {
	l.SetLevel(ParseLevel(levelStr))
}

// GetLevel returns the current log level
func (l *Logger) GetLevel() Level {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.level
}

// AddOutput adds a new output writer
func (l *Logger) AddOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.outputs = append(l.outputs, w)

	// Update zap logger with new output
	encConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var cores []zapcore.Core
	enc := zapcore.NewJSONEncoder(encConfig)
	for _, output := range l.outputs {
		core := zapcore.NewCore(enc, zapcore.AddSync(output), l.level)
		cores = append(cores, core)
	}

	l.Logger = zap.New(zapcore.NewTee(cores...))
}

// ClearOutputs removes all output writers
func (l *Logger) ClearOutputs() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.outputs = make([]io.Writer, 0)

	// Reset zap logger
	l.Logger = zap.New(zapcore.NewTee())
}

// With creates a child logger with the given fields
func (l *Logger) With(fields ...zapcore.Field) *Logger {
	return &Logger{
		Logger:     l.Logger.With(fields...),
		level:      l.level,
		outputs:    l.outputs,
		timeFormat: l.timeFormat,
	}
}
