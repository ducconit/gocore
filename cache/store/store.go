package store

import (
	"context"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
)

// Store represents a cache store interface
type Store interface {
	// Get retrieves a value from cache by key
	Get(ctx context.Context, key any) (any, error)

	// Set stores a value in cache
	Set(ctx context.Context, key any, value any, options ...store.Option) error

	// Delete removes a value from cache
	Delete(ctx context.Context, key any) error

	// Clear removes all values from cache
	Clear(ctx context.Context) error

	// GetType returns the store type
	GetType() string
}

// CacheStore wraps gocache store implementation
type CacheStore struct {
	cache *cache.Cache[any]
}

// NewStore creates a new cache store
func NewStore(store store.StoreInterface) *CacheStore {
	return &CacheStore{
		cache: cache.New[any](store),
	}
}

// Get retrieves a value from cache by key
func (s *CacheStore) Get(ctx context.Context, key any) (any, error) {
	return s.cache.Get(ctx, key)
}

// Set stores a value in cache
func (s *CacheStore) Set(ctx context.Context, key any, value any, options ...store.Option) error {
	return s.cache.Set(ctx, key, value, options...)
}

// Delete removes a value from cache
func (s *CacheStore) Delete(ctx context.Context, key any) error {
	return s.cache.Delete(ctx, key)
}

// Clear removes all values from cache
func (s *CacheStore) Clear(ctx context.Context) error {
	return s.cache.Clear(ctx)
}

// GetType returns the store type
func (s *CacheStore) GetType() string {
	return s.cache.GetType()
}

// WithExpiration returns a store option with expiration
func WithExpiration(expiration time.Duration) store.Option {
	return store.WithExpiration(expiration)
}

// WithTags returns a store option with tags
func WithTags(tags ...string) store.Option {
	return store.WithTags(tags)
}

// WithCost returns a store option with cost
func WithCost(cost int64) store.Option {
	return store.WithCost(cost)
}
