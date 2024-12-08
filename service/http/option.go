package http

import (
	"net/http"

	"github.com/ducconit/gocore/logger"
)

// HTTPOption represents an option for configuring HTTPService
type HTTPOption func(*HTTPService)

// WithAddr sets the server address
func WithAddr(addr string) HTTPOption {
	return func(s *HTTPService) {
		s.addr = addr
	}
}

// WithHandler sets the HTTP handler
func WithHandler(handler http.Handler) HTTPOption {
	return func(s *HTTPService) {
		s.handler = handler
	}
}

// WithLogger sets the logger
func WithLogger(l *logger.Logger) HTTPOption {
	return func(s *HTTPService) {
		s.logger = l
	}
}
