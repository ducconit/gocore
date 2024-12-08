package service

import (
	"context"
)

// Service interface defines the basic methods that all services should implement
type Service interface {
	// Start initializes and starts the service
	Start(ctx context.Context) error

	// Stop gracefully shuts down the service
	Stop(ctx context.Context) error

	// Health checks if the service is healthy
	Health(ctx context.Context) error

	// Name returns the service name
	Name() string
}
