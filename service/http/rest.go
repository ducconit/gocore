package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/ducconit/gocore/logger"
	"github.com/ducconit/gocore/utils"
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

// StartDaemon starts the HTTP service in daemon mode (non-blocking)
func (s *HTTPService) StartDaemon(ctx context.Context) error {
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

// Start starts the HTTP service and blocks until interrupted
func (s *HTTPService) Start() error {
	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the service in daemon mode
	if err := s.StartDaemon(ctx); err != nil {
		return err
	}

	// Create signal channel to listen for interrupt signals
	utils.WaitOSSignalHandler(func() {
		s.logger.Info("Received signal, shutting down")
		if err := s.Stop(context.Background()); err != nil {
			s.logger.Error("Error shutting down service", zap.Error(err))
		}
	}, os.Interrupt, syscall.SIGTERM)

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
