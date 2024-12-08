package errors

import (
	"fmt"
	"runtime"
)

// Error represents a custom error with stack trace and metadata
type Error struct {
	Message    string
	Code       string
	Err        error
	StackTrace string
	Metadata   map[string]any
}

// New creates a new Error instance
func New(message string) *Error {
	return &Error{
		Message:    message,
		StackTrace: getStackTrace(),
		Metadata:   make(map[string]any),
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, message string) *Error {
	if err == nil {
		return nil
	}

	return &Error{
		Message:    message,
		Err:        err,
		StackTrace: getStackTrace(),
		Metadata:   make(map[string]any),
	}
}

// WithCode adds an error code to the error
func (e *Error) WithCode(code string) *Error {
	e.Code = code
	return e
}

// WithMetadata adds metadata to the error
func (e *Error) WithMetadata(key string, value any) *Error {
	e.Metadata[key] = value
	return e
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *Error) Unwrap() error {
	return e.Err
}

func getStackTrace() string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	var trace string
	for {
		frame, more := frames.Next()
		trace += fmt.Sprintf("\n%s:%d", frame.File, frame.Line)
		if !more {
			break
		}
	}
	return trace
}
