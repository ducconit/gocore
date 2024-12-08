package cache

import (
	"context"
	"errors"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/ducconit/gocore/cache/store"
	cacheStore "github.com/eko/gocache/lib/v4/store"
	goCacheStore "github.com/eko/gocache/store/go_cache/v4"
	memcacheStore "github.com/eko/gocache/store/memcache/v4"
	redisStore "github.com/eko/gocache/store/redis/v4"
	goCache "github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
)

var (
	// ErrInvalidOptions is returned when the options are invalid
	ErrInvalidOptions = errors.New("invalid cache options")

	// DefaultExpiration is the default expiration time
	DefaultExpiration = 5 * time.Minute

	// DefaultCleanupInterval is the default cleanup interval
	DefaultCleanupInterval = 10 * time.Minute

	// DefaultMaxEntries is the default maximum number of entries
	DefaultMaxEntries = 10000
)

// Cache interface defines methods for caching operations
type Cache interface {
	// Get retrieves a value from cache
	Get(ctx context.Context, key string) (any, error)

	// Set stores a value in cache
	Set(ctx context.Context, key string, value any, expiration time.Duration) error

	// Delete removes a value from cache
	Delete(ctx context.Context, key string) error

	// Clear removes all values from cache
	Clear(ctx context.Context) error

	// GetMulti retrieves multiple values from cache
	GetMulti(ctx context.Context, keys []string) (map[string]any, error)

	// SetMulti stores multiple values in cache
	SetMulti(ctx context.Context, items map[string]any, expiration time.Duration) error

	// DeleteMulti removes multiple values from cache
	DeleteMulti(ctx context.Context, keys []string) error

	// GetStore returns the underlying store
	GetStore() store.Store
}

// Options represents cache configuration options
type Options struct {
	// DefaultExpiration is the default expiration time for cache entries
	DefaultExpiration time.Duration

	// CleanupInterval is the interval for cleanup of expired entries
	CleanupInterval time.Duration

	// MaxEntries is the maximum number of items in the cache
	MaxEntries int

	// OnEvicted is called when an entry is evicted from the cache
	OnEvicted func(key string, value any)

	// Redis options. default addr is localhost:6379
	RedisOptions *redis.Options

	// Memcached addresses. Default is localhost:11211
	MemcachedAddrs []string

	// KeyPrefix is the prefix added to all keys
	KeyPrefix string
}

// Validate validates the options
func (o *Options) Validate() error {
	if o.DefaultExpiration < 0 {
		return errors.New("default expiration must be >= 0")
	}
	if o.CleanupInterval < 0 {
		return errors.New("cleanup interval must be >= 0")
	}
	if o.MaxEntries < 0 {
		return errors.New("max entries must be >= 0")
	}
	return nil
}

// NewOptions creates default cache options
func NewOptions() *Options {
	return &Options{
		DefaultExpiration: DefaultExpiration,
		CleanupInterval:   DefaultCleanupInterval,
		MaxEntries:        DefaultMaxEntries,
		RedisOptions: &redis.Options{
			Addr: "localhost:6379",
		},
		MemcachedAddrs: []string{"localhost:11211"},
		KeyPrefix:      "",
	}
}

// cacheImpl implements Cache interface
type cacheImpl struct {
	store  store.Store
	prefix string
	opts   *Options
}

// NewMemoryCache creates a new memory cache instance
func NewMemoryCache(opts *Options) (Cache, error) {
	if opts == nil {
		opts = NewOptions()
	}

	if err := opts.Validate(); err != nil {
		return nil, ErrInvalidOptions
	}

	client := goCache.New(opts.DefaultExpiration, opts.CleanupInterval)
	if opts.OnEvicted != nil {
		client.OnEvicted(opts.OnEvicted)
	}
	goCacheStore := goCacheStore.NewGoCache(client)

	return &cacheImpl{
		store:  store.NewStore(goCacheStore),
		prefix: opts.KeyPrefix,
		opts:   opts,
	}, nil
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(opts *Options) (Cache, error) {
	if opts == nil {
		opts = NewOptions()
	}

	if err := opts.Validate(); err != nil {
		return nil, ErrInvalidOptions
	}

	redisClient := redis.NewClient(opts.RedisOptions)
	// Test connection
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	redisStore := redisStore.NewRedis(redisClient)

	return &cacheImpl{
		store:  store.NewStore(redisStore),
		prefix: opts.KeyPrefix,
		opts:   opts,
	}, nil
}

// NewMemcachedCache creates a new Memcached cache instance
func NewMemcachedCache(opts *Options) (Cache, error) {
	if opts == nil {
		opts = NewOptions()
	}

	if err := opts.Validate(); err != nil {
		return nil, ErrInvalidOptions
	}

	if len(opts.MemcachedAddrs) == 0 {
		return nil, errors.New("memcached addresses are required")
	}

	memcacheClient := memcache.New(opts.MemcachedAddrs...)
	if err := memcacheClient.Ping(); err != nil {
		return nil, err
	}

	memcacheStore := memcacheStore.NewMemcache(memcacheClient)

	return &cacheImpl{
		store:  store.NewStore(memcacheStore),
		prefix: opts.KeyPrefix,
		opts:   opts,
	}, nil
}

func (c *cacheImpl) buildKey(key string) string {
	if c.prefix == "" {
		return key
	}
	return c.prefix + ":" + key
}

// Get retrieves a value from cache
func (c *cacheImpl) Get(ctx context.Context, key string) (any, error) {
	return c.store.Get(ctx, c.buildKey(key))
}

// Set stores a value in cache
func (c *cacheImpl) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	if expiration == 0 {
		expiration = c.opts.DefaultExpiration
	}
	return c.store.Set(ctx, c.buildKey(key), value, store.WithExpiration(expiration))
}

// Delete removes a value from cache
func (c *cacheImpl) Delete(ctx context.Context, key string) error {
	return c.store.Delete(ctx, c.buildKey(key))
}

// Clear removes all values from cache
func (c *cacheImpl) Clear(ctx context.Context) error {
	return c.store.Clear(ctx)
}

// GetMulti retrieves multiple values from cache
func (c *cacheImpl) GetMulti(ctx context.Context, keys []string) (map[string]any, error) {
	result := make(map[string]any)
	for _, key := range keys {
		value, err := c.Get(ctx, key)
		if err != nil {
			var notFoundError *cacheStore.NotFound
			if !errors.As(err, &notFoundError) {
				return nil, err
			}
			continue
		}
		result[key] = value
	}
	return result, nil
}

// SetMulti stores multiple values in cache
func (c *cacheImpl) SetMulti(ctx context.Context, items map[string]any, expiration time.Duration) error {
	if expiration == 0 {
		expiration = c.opts.DefaultExpiration
	}
	for key, value := range items {
		if err := c.Set(ctx, key, value, expiration); err != nil {
			return err
		}
	}
	return nil
}

// DeleteMulti removes multiple values from cache
func (c *cacheImpl) DeleteMulti(ctx context.Context, keys []string) error {
	for _, key := range keys {
		if err := c.Delete(ctx, key); err != nil {
			return err
		}
	}
	return nil
}

// GetStore returns the underlying store
func (c *cacheImpl) GetStore() store.Store {
	return c.store
}
