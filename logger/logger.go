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
	level   Level
	outputs []io.Writer
	mu      sync.RWMutex
}

// New creates a new logger instance
func New(opts ...Option) *Logger {
	l := &Logger{
		level:   InfoLevel, // default level
		outputs: make([]io.Writer, 0),
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
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
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
func Default(opts ...Option) *Logger {
	options := append([]Option{
		WithLevel(InfoLevel),
		WithConsole(),
	}, opts...)

	return New(options...)
}

// SetDefault sets the default logger instance
func SetDefault(l *Logger) {
	defaultLogger = l
}

func init() {
	SetDefault(Default())
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

func (l *Logger) Log(level Level, msg string, fields ...zap.Field) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Logger.Log(level, msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Log(DebugLevel, msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Log(InfoLevel, msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Log(WarnLevel, msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.Log(ErrorLevel, msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.Log(PanicLevel, msg, fields...)
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	return &Logger{
		Logger:  l.Logger.With(fields...),
		level:   l.level,
		outputs: l.outputs,
	}
}
