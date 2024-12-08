package http

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ducconit/gocore/logger"
	"go.uber.org/zap"
)

// HTTPService represents a base HTTP service
type HTTPService struct {
	name     string
	server   *http.Server
	handler  http.Handler
	addr     string
	mu       sync.RWMutex
	started  bool
	stopChan chan struct{}
	logger   *logger.Logger
}

// HTTPOption represents an option for configuring HTTPService
type HTTPOption func(*HTTPService)

// WithAddress sets the server address
func WithAddress(addr string) HTTPOption {
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

// NewHTTPService creates a new HTTP service
func NewHTTPService(name string, opts ...HTTPOption) *HTTPService {
	s := &HTTPService{
		name:     name,
		addr:     ":3000",             // default address
		stopChan: make(chan struct{}), // signal channel with lower capacity
		logger:   logger.Default(),    // default console logger
	}

	// Apply options
	for _, opt := range opts {
		opt(s)
	}

	// Create server
	s.server = &http.Server{
		Addr:              s.addr,
		Handler:           s.handler,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	return s
}

// Start implements Service.Start
func (s *HTTPService) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.started {
		s.mu.Unlock()
		return fmt.Errorf("service %s already started", s.name)
	}
	s.started = true
	s.mu.Unlock()

	// Start server in a goroutine
	go func() {
		s.logger.Info("Starting HTTP service",
			zap.String("service", s.name),
			zap.String("address", s.addr))

		if err := s.server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				s.logger.Error("HTTP service error",
					zap.String("service", s.name),
					zap.Error(err))
			} else {
				s.logger.Info("HTTP service stopped",
					zap.String("service", s.name))
			}
		}
	}()

	// Wait for context cancellation or stop signal
	go func() {
		select {
		case <-ctx.Done():
			s.Stop(context.Background())
		case <-s.stopChan: // signal channel
			return
		}
	}()

	return nil
}

// Stop implements Service.Stop
func (s *HTTPService) Stop(ctx context.Context) error {
	s.mu.Lock()
	if !s.started {
		s.mu.Unlock()
		return fmt.Errorf("service %s not started", s.name)
	}
	s.started = false
	s.mu.Unlock()

	// Signal stop
	close(s.stopChan)

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	s.logger.Info("Stopping HTTP service",
		zap.String("service", s.name))

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("error shutting down service %s: %v", s.name, err)
	}

	return nil
}

// Health implements Service.Health
func (s *HTTPService) Health(ctx context.Context) error {
	s.mu.RLock()
	started := s.started
	s.mu.RUnlock()

	if !started {
		return fmt.Errorf("service %s is not running", s.name)
	}
	return nil
}

// Name implements Service.Name
func (s *HTTPService) Name() string {
	return s.name
}

// Server returns the underlying http.Server
func (s *HTTPService) Server() *http.Server {
	return s.server
}

// IsStarted returns whether the service is started
func (s *HTTPService) IsStarted() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.started
}

// SetHandler sets the HTTP handler
func (s *HTTPService) SetHandler(handler http.Handler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handler = handler
	s.server.Handler = handler
}
