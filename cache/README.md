# Cache Package

The cache package provides a flexible caching system with support for multiple cache backends.

## Features

- Multiple Cache Backends (Memory, Redis)
- TTL Support
- Automatic Key Expiration
- Thread-safe Operations
- Cache Tags Support
- Bulk Operations

## Usage

### Memory Cache

```go
import "github.com/ducconit/gocore/cache"

// Create memory cache
c := cache.NewMemoryCache()

// Set value with TTL
c.Set("key", "value", 5*time.Minute)

// Get value
val, err := c.Get("key")
if err != nil {
    log.Printf("Cache miss: %v", err)
}
```

### Redis Cache

```go
// Create Redis cache
c := cache.NewRedisCache(
    cache.WithRedisAddr("localhost:6379"),
    cache.WithRedisPassword("password"),
    cache.WithRedisDB(0),
)

// Set value
c.Set("key", "value", 0) // 0 for no expiration
```

## Cache Interface

```go
type Cache interface {
    Get(key string) (interface{}, error)
    Set(key string, value interface{}, ttl time.Duration) error
    Delete(key string) error
    Clear() error
    Has(key string) bool
    GetMultiple(keys []string) (map[string]interface{}, error)
    SetMultiple(values map[string]interface{}, ttl time.Duration) error
    DeleteMultiple(keys []string) error
}
```

## Options

### Redis Options

| Option | Description | Default |
|--------|-------------|---------|
| WithRedisAddr | Redis server address | "localhost:6379" |
| WithRedisPassword | Redis password | "" |
| WithRedisDB | Redis database number | 0 |
| WithRedisPoolSize | Connection pool size | 10 |

### Memory Cache Options

| Option | Description | Default |
|--------|-------------|---------|
| WithCapacity | Maximum items in cache | 1000 |
| WithEvictionPolicy | Cache eviction policy | LRU |
| WithCleanupInterval | Cleanup interval | 10m |

## Examples

### Cache with Tags

```go
// Set values with tags
c.TaggedSet("user:1", user1, []string{"users", "active"}, 1*time.Hour)
c.TaggedSet("user:2", user2, []string{"users", "inactive"}, 1*time.Hour)

// Invalidate by tag
c.InvalidateTag("active")
```

### Bulk Operations

```go
// Set multiple values
values := map[string]interface{}{
    "key1": "value1",
    "key2": "value2",
}
c.SetMultiple(values, 5*time.Minute)

// Get multiple values
keys := []string{"key1", "key2"}
results, err := c.GetMultiple(keys)
```

### Cache Middleware

```go
func CacheMiddleware(c *cache.Cache) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        key := ctx.Request.URL.String()
        if value, err := c.Get(key); err == nil {
            ctx.JSON(200, value)
            ctx.Abort()
            return
        }
        ctx.Next()
        // Cache response
        if ctx.Writer.Status() == 200 {
            c.Set(key, ctx.Keys["response"], 5*time.Minute)
        }
    }
}
```

## Best Practices

1. Choose appropriate TTL values
2. Use bulk operations for better performance
3. Implement proper error handling
4. Monitor cache usage and hit rates
5. Use tags for related cache invalidation
6. Implement fallback mechanisms
