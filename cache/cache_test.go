package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOptions(t *testing.T) {
	opts := NewOptions()
	assert.Equal(t, DefaultExpiration, opts.DefaultExpiration)
	assert.Equal(t, DefaultCleanupInterval, opts.CleanupInterval)
	assert.Equal(t, DefaultMaxEntries, opts.MaxEntries)
	assert.Equal(t, "localhost:6379", opts.RedisOptions.Addr)
	assert.Equal(t, []string{"localhost:11211"}, opts.MemcachedAddrs)
}

func TestOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *Options
		wantErr bool
	}{
		{
			name:    "valid options",
			opts:    NewOptions(),
			wantErr: false,
		},
		{
			name: "negative default expiration",
			opts: &Options{
				DefaultExpiration: -1,
			},
			wantErr: true,
		},
		{
			name: "negative cleanup interval",
			opts: &Options{
				CleanupInterval: -1,
			},
			wantErr: true,
		},
		{
			name: "negative max entries",
			opts: &Options{
				MaxEntries: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMemoryCache(t *testing.T) {
	ctx := context.Background()
	cache, err := NewMemoryCache(nil)
	require.NoError(t, err)

	t.Run("basic operations", func(t *testing.T) {
		// Set
		err := cache.Set(ctx, "key1", "value1", time.Minute)
		require.NoError(t, err)

		// Get
		value, err := cache.Get(ctx, "key1")
		require.NoError(t, err)
		assert.Equal(t, "value1", value)

		// Delete
		err = cache.Delete(ctx, "key1")
		require.NoError(t, err)

		// Get after delete
		_, err = cache.Get(ctx, "key1")
		assert.Error(t, err)
	})

	t.Run("expiration", func(t *testing.T) {
		err := cache.Set(ctx, "key2", "value2", time.Millisecond)
		require.NoError(t, err)

		time.Sleep(2 * time.Millisecond)

		_, err = cache.Get(ctx, "key2")
		assert.Error(t, err)
	})

	t.Run("multi operations", func(t *testing.T) {
		items := map[string]any{
			"mkey1": "mvalue1",
			"mkey2": "mvalue2",
		}

		// SetMulti
		err := cache.SetMulti(ctx, items, time.Minute)
		require.NoError(t, err)

		// GetMulti
		values, err := cache.GetMulti(ctx, []string{"mkey1", "mkey2", "nonexistent"})
		require.NoError(t, err)
		assert.Equal(t, "mvalue1", values["mkey1"])
		assert.Equal(t, "mvalue2", values["mkey2"])
		assert.NotContains(t, values, "nonexistent")

		// DeleteMulti
		err = cache.DeleteMulti(ctx, []string{"mkey1", "mkey2"})
		require.NoError(t, err)

		values, err = cache.GetMulti(ctx, []string{"mkey1", "mkey2"})
		require.NoError(t, err)
		assert.Empty(t, values)
	})

	t.Run("clear", func(t *testing.T) {
		err := cache.Set(ctx, "key3", "value3", time.Minute)
		require.NoError(t, err)

		err = cache.Clear(ctx)
		require.NoError(t, err)

		_, err = cache.Get(ctx, "key3")
		assert.Error(t, err)
	})

	t.Run("key prefix", func(t *testing.T) {
		opts := NewOptions()
		opts.KeyPrefix = "test"
		cache, err := NewMemoryCache(opts)
		require.NoError(t, err)

		err = cache.Set(ctx, "key4", "value4", time.Minute)
		require.NoError(t, err)

		value, err := cache.Get(ctx, "key4")
		require.NoError(t, err)
		assert.Equal(t, "value4", value)
	})
}

func TestRedisCache(t *testing.T) {
	t.Skip("Requires Redis server")
	// Similar tests as memory cache
}

func TestMemcachedCache(t *testing.T) {
	t.Skip("Requires Memcached server")
	// Similar tests as memory cache
}
