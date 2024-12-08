package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ducconit/gocore/logger"
	"github.com/stretchr/testify/assert"
)

type mockHandler struct {
	handleFunc func(w http.ResponseWriter, r *http.Request)
}

func (h *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.handleFunc != nil {
		h.handleFunc(w, r)
	}
}

const testPort = ":8089" // Use a test port instead of 80

func TestNewHTTPService(t *testing.T) {
	t.Run("default_configuration", func(t *testing.T) {
		svc := NewHTTPService("test")
		assert.Equal(t, "test", svc.Name())
		assert.Equal(t, ":3000", svc.addr)
		assert.NotNil(t, svc.logger)
		assert.NotNil(t, svc.stopChan)
		assert.False(t, svc.started)
	})

	t.Run("with_options", func(t *testing.T) {
		handler := &mockHandler{}
		customLogger := logger.NewLogger()
		svc := NewHTTPService("test",
			WithAddress(":8080"),
			WithHandler(handler),
			WithLogger(customLogger),
		)

		assert.Equal(t, ":8080", svc.addr)
		assert.Equal(t, handler, svc.handler)
		assert.Equal(t, customLogger, svc.logger)
	})
}

func TestHTTPService_Start(t *testing.T) {
	t.Run("successful_start", func(t *testing.T) {
		handler := &mockHandler{
			handleFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, "ok")
			},
		}

		svc := NewHTTPService("test",
			WithAddress(testPort),
			WithHandler(handler),
		)

		ctx := context.Background()
		err := svc.Start(ctx)
		assert.NoError(t, err)
		assert.True(t, svc.IsStarted())

		// Wait for server to start
		time.Sleep(100 * time.Millisecond)

		// Make test request
		addr := fmt.Sprintf("http://localhost%s", svc.server.Addr)
		resp, err := http.Get(addr)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, "ok", string(body))
		resp.Body.Close()

		// Cleanup
		svc.Stop(ctx)
	})

	t.Run("already_started", func(t *testing.T) {
		svc := NewHTTPService("test", WithAddress(testPort))
		ctx := context.Background()

		// Start first time
		err := svc.Start(ctx)
		assert.NoError(t, err)

		// Try to start again
		err = svc.Start(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already started")

		// Cleanup
		svc.Stop(ctx)
	})
}

func TestHTTPService_Stop(t *testing.T) {
	t.Run("successful_stop", func(t *testing.T) {
		svc := NewHTTPService("test", WithAddress(testPort))
		ctx := context.Background()

		// Start service
		err := svc.Start(ctx)
		assert.NoError(t, err)
		assert.True(t, svc.IsStarted())

		// Stop service
		err = svc.Stop(ctx)
		assert.NoError(t, err)

		// Verify server is no longer accepting connections
		_, err = http.Get(fmt.Sprintf("http://localhost%s", svc.server.Addr))
		assert.Error(t, err)
	})
}

func TestHTTPService_Health(t *testing.T) {
	t.Run("healthy_service", func(t *testing.T) {
		svc := NewHTTPService("test", WithAddress(testPort))
		ctx := context.Background()

		// Start service
		err := svc.Start(ctx)
		assert.NoError(t, err)

		// Check health
		err = svc.Health(ctx)
		assert.NoError(t, err)

		// Cleanup
		svc.Stop(ctx)
	})

	t.Run("not_started_service", func(t *testing.T) {
		svc := NewHTTPService("test")
		ctx := context.Background()

		// Check health without starting
		err := svc.Health(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "is not running")
	})
}

func TestHTTPService_SetHandler(t *testing.T) {
	svc := NewHTTPService("test")

	handler := &mockHandler{
		handleFunc: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	}

	svc.SetHandler(handler)
	assert.Equal(t, handler, svc.handler)
	assert.Equal(t, handler, svc.server.Handler)
}

func TestHTTPService_Server(t *testing.T) {
	svc := NewHTTPService("test")
	server := svc.Server()
	assert.NotNil(t, server)
	assert.Equal(t, svc.server, server)
}

func TestHTTPService_Integration(t *testing.T) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "test response")
	}))
	defer ts.Close()

	// Create service with test server's handler
	svc := NewHTTPService("test",
		WithAddress(testPort),
		WithHandler(ts.Config.Handler),
	)

	ctx := context.Background()

	// Start service
	err := svc.Start(ctx)
	assert.NoError(t, err)

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Make request
	resp, err := http.Get(fmt.Sprintf("http://localhost%s", svc.server.Addr))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "test response", string(body))
	resp.Body.Close()

	// Stop service
	err = svc.Stop(ctx)
	assert.NoError(t, err)
}
