package logger

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
	}{
		{
			name: "default logger",
			opts: []Option{},
		},
		{
			name: "with console",
			opts: []Option{WithConsole()},
		},
		{
			name: "with level",
			opts: []Option{WithLevel(DebugLevel)},
		},
		{
			name: "with level string",
			opts: []Option{WithLevelString("debug")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := New(tt.opts...)
			assert.NotNil(t, logger)
			assert.NotNil(t, logger.Logger)
		})
	}
}

func newTestLogger(buf *bytes.Buffer) *Logger {
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

	// Create core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encConfig),
		zapcore.AddSync(buf),
		zapcore.DebugLevel,
	)

	// Create logger
	l := &Logger{
		Logger:  zap.New(core),
		level:   DebugLevel,
		outputs: []io.Writer{buf},
	}

	return l
}

func TestLogger_Levels(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	tests := []struct {
		level   Level
		logFunc func(string, ...zapcore.Field)
	}{
		{DebugLevel, logger.Debug},
		{InfoLevel, logger.Info},
		{WarnLevel, logger.Warn},
		{ErrorLevel, logger.Error},
		{DPanicLevel, logger.DPanic},
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			buf.Reset()
			msg := "test message"
			tt.logFunc(msg)
			output := buf.String()
			t.Logf("Output: %s", output)
			assert.Contains(t, output, msg)
			assert.Contains(t, strings.ToLower(output), tt.level.String())
		})
	}
}

func TestLogger_WithFields(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	// Create child logger with fields
	childLogger := logger.With(
		zap.String("module", "test"),
		zap.Int("version", 1),
	)

	// Log message with child logger
	msg := "test with fields"
	childLogger.Info(msg)
	output := buf.String()

	t.Logf("Output: %s", output)
	assert.Contains(t, output, msg)
	assert.Contains(t, output, "module")
	assert.Contains(t, output, "test")
	assert.Contains(t, output, "version")
	assert.Contains(t, output, "1")
}

func TestLogger_FileOutput(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "logger_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create logger with file output
	logFile := filepath.Join(tmpDir, "test.log")
	logger := New(WithFile(logFile))

	// Log some messages
	msg := "test file output"
	logger.Info(msg)

	// Wait for file write
	time.Sleep(100 * time.Millisecond)

	// Read log file
	content, err := os.ReadFile(logFile)
	assert.NoError(t, err)
	t.Logf("File content: %s", string(content))
	assert.Contains(t, string(content), msg)
}

func TestLogger_MultipleOutputs(t *testing.T) {
	// Create buffers for testing
	var buf1, buf2 bytes.Buffer

	// Create encoder config
	encConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}

	// Create cores
	core1 := zapcore.NewCore(
		zapcore.NewJSONEncoder(encConfig),
		zapcore.AddSync(&buf1),
		zapcore.DebugLevel,
	)
	core2 := zapcore.NewCore(
		zapcore.NewJSONEncoder(encConfig),
		zapcore.AddSync(&buf2),
		zapcore.DebugLevel,
	)

	// Create logger
	logger := &Logger{
		Logger:  zap.New(zapcore.NewTee(core1, core2)),
		level:   DebugLevel,
		outputs: []io.Writer{&buf1, &buf2},
	}

	// Log message
	msg := "test multiple outputs"
	logger.Info(msg)

	t.Logf("Buffer 1: %s", buf1.String())
	t.Logf("Buffer 2: %s", buf2.String())

	// Verify message appears in both outputs
	assert.Contains(t, buf1.String(), msg)
	assert.Contains(t, buf2.String(), msg)
}

func TestLogger_LevelParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
	}{
		{"debug", DebugLevel},
		{"DEBUG", DebugLevel},
		{"info", InfoLevel},
		{"warn", WarnLevel},
		{"warning", WarnLevel},
		{"error", ErrorLevel},
		{"dpanic", DPanicLevel},
		{"panic", PanicLevel},
		{"fatal", FatalLevel},
		{"invalid", InfoLevel}, // default to info
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			level := ParseLevel(tt.input)
			assert.Equal(t, tt.expected, level)
		})
	}
}

func TestLogger_DynamicOutputs(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	// Log message
	msg := "test dynamic output"
	logger.Info(msg)
	assert.Contains(t, buf.String(), msg)

	// Clear outputs and create new buffer
	logger.ClearOutputs()
	var newBuf bytes.Buffer
	logger.AddOutput(&newBuf)

	// Log new message
	newMsg := "new message"
	logger.Info(newMsg)
	assert.Contains(t, newBuf.String(), newMsg)
}

func TestLogger_TimeFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	// Log with timestamp
	msg := "test time format"
	logger.Info(msg)
	output := buf.String()
	t.Logf("Output: %s", output)

	// Verify output contains timestamp
	assert.Contains(t, output, "time")
	assert.Contains(t, output, time.Now().Format("2006"))
}
