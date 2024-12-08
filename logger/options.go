package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Option represents a logger option
type Option func(*Logger)

// WithLevel sets the minimum log level
func WithLevel(level Level) Option {
	return func(l *Logger) {
		l.level = level
	}
}

// WithLevelString sets the minimum log level using a string
func WithLevelString(levelStr string) Option {
	return func(l *Logger) {
		l.level = ParseLevel(levelStr)
	}
}

// WithConsole adds console output
func WithConsole() Option {
	return func(l *Logger) {
		l.outputs = append(l.outputs, os.Stdout)
	}
}

// WithFile adds file output
func WithFile(filename string) Option {
	return func(l *Logger) {
		// Create directory if not exists
		dir := filepath.Dir(filename)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creating log directory: %v\n", err)
			return
		}

		// Open file in append mode
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error opening log file: %v\n", err)
			return
		}

		l.outputs = append(l.outputs, f)
	}
}

// WithTimeFormat sets the time format for log messages
func WithTimeFormat(format string) Option {
	return func(l *Logger) {
		l.timeFormat = format
	}
}

// WithOutput adds a custom output writer
func WithOutput(w io.Writer) Option {
	return func(l *Logger) {
		l.outputs = append(l.outputs, w)
	}
}
