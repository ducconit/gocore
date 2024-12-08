package queue

import (
	"context"
	"time"
)

// Message represents a queue message
type Message struct {
	ID        string
	Body      []byte
	Metadata  map[string]string
	Timestamp time.Time
}

// Queue interface defines methods for queue operations
type Queue interface {
	// Push adds a message to the queue
	Push(ctx context.Context, msg *Message) error
	
	// Pop retrieves and removes a message from the queue
	Pop(ctx context.Context) (*Message, error)
	
	// Peek retrieves but does not remove a message from the queue
	Peek(ctx context.Context) (*Message, error)
	
	// Length returns the number of messages in the queue
	Length(ctx context.Context) (int64, error)
	
	// Clear removes all messages from the queue
	Clear(ctx context.Context) error
}

// Consumer interface defines methods for message consumption
type Consumer interface {
	// Start starts consuming messages
	Start(ctx context.Context) error
	
	// Stop stops consuming messages
	Stop(ctx context.Context) error
	
	// OnMessage is called when a message is received
	OnMessage(handler func(ctx context.Context, msg *Message) error)
}

// Producer interface defines methods for message production
type Producer interface {
	// Start starts the producer
	Start(ctx context.Context) error
	
	// Stop stops the producer
	Stop(ctx context.Context) error
	
	// Send sends a message
	Send(ctx context.Context, msg *Message) error
}

// Options represents queue configuration options
type Options struct {
	// MaxSize is the maximum number of messages in the queue
	MaxSize int64
	
	// BatchSize is the number of messages to process in a batch
	BatchSize int
	
	// PollInterval is the interval between polls
	PollInterval time.Duration
	
	// RetryCount is the number of times to retry failed operations
	RetryCount int
	
	// RetryDelay is the delay between retries
	RetryDelay time.Duration
}

// NewOptions creates default queue options
func NewOptions() *Options {
	return &Options{
		MaxSize:      10000,
		BatchSize:    100,
		PollInterval: time.Second,
		RetryCount:   3,
		RetryDelay:   time.Second,
	}
}
