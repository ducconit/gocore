package config

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.json")
	configContent := []byte(`{
		"database": {
			"host": "localhost",
			"port": 5432,
			"name": "testdb"
		},
		"server": {
			"port": 8080,
			"debug": true
		},
		"features": ["auth", "api", "websocket"]
	}`)
	err := os.WriteFile(configFile, configContent, 0644)
	assert.NoError(t, err)

	cfg := NewConfig()
	err = cfg.LoadFromFile(configFile, WithConfigType("json"))
	assert.NoError(t, err)

	t.Run("get_values", func(t *testing.T) {
		assert.Equal(t, "localhost", cfg.GetString("database.host"))
		assert.Equal(t, 5432, cfg.GetInt("database.port"))
		assert.Equal(t, true, cfg.GetBool("server.debug"))
		assert.Equal(t, []string{"auth", "api", "websocket"}, cfg.GetStringSlice("features"))
	})
}

func TestDatabaseConfig(t *testing.T) {
	// Create SQLite database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := sql.Open("sqlite3", dbPath)
	assert.NoError(t, err)
	defer db.Close()

	// Create table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS config (
			key_name VARCHAR(255) PRIMARY KEY,
			value TEXT NOT NULL
		)
	`)
	assert.NoError(t, err)

	// Create config instance
	cfg := NewConfig()

	// Insert test data
	_, err = db.Exec("INSERT INTO config (key_name, value) VALUES (?, ?)", "test.key", `"test value"`)
	assert.NoError(t, err)

	// Load from database
	err = cfg.LoadFromDB(db, "config")
	assert.NoError(t, err)

	// Test values
	assert.Equal(t, "test value", cfg.GetString("test.key"))
}

func TestExtendedGetMethods(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.json")
	configContent := []byte(`{
		"duration": "5s",
		"time": "2024-12-08T15:56:00+07:00",
		"numbers": {
			"uint": 42,
			"int32": 123,
			"int64": 456,
			"float64": 3.14
		},
		"maps": {
			"stringMap": {
				"key1": "value1",
				"key2": "value2"
			},
			"sliceMap": {
				"key1": ["value1", "value2"],
				"key2": ["value3", "value4"]
			}
		}
	}`)
	err := os.WriteFile(configFile, configContent, 0644)
	assert.NoError(t, err)

	cfg := NewConfig()
	err = cfg.LoadFromFile(configFile, WithConfigType("json"))
	assert.NoError(t, err)

	t.Run("duration_and_time", func(t *testing.T) {
		assert.Equal(t, 5*time.Second, cfg.GetDuration("duration"))
		expectedTime, _ := time.Parse(time.RFC3339, "2024-12-08T15:56:00+07:00")
		assert.Equal(t, expectedTime, cfg.GetTime("time"))
	})

	t.Run("number_types", func(t *testing.T) {
		assert.Equal(t, uint(42), cfg.GetUint("numbers.uint"))
		assert.Equal(t, int32(123), cfg.GetInt32("numbers.int32"))
		assert.Equal(t, int64(456), cfg.GetInt64("numbers.int64"))
		assert.Equal(t, float64(3.14), cfg.GetFloat64("numbers.float64"))
	})

	t.Run("map_types", func(t *testing.T) {
		stringMap := cfg.GetStringMapString("maps.stringMap")
		assert.Equal(t, "value1", stringMap["key1"])
		assert.Equal(t, "value2", stringMap["key2"])

		sliceMap := cfg.GetStringMapStringSlice("maps.sliceMap")
		assert.Equal(t, []string{"value1", "value2"}, sliceMap["key1"])
		assert.Equal(t, []string{"value3", "value4"}, sliceMap["key2"])
	})
}

func TestExtendedMethods(t *testing.T) {
	cfg := NewConfig()

	t.Run("set_values", func(t *testing.T) {
		cfg.Set("new.key", "value")
		assert.Equal(t, "value", cfg.GetString("new.key"))
	})

	t.Run("defaults", func(t *testing.T) {
		cfg.SetDefault("missing.key", "default")
		assert.Equal(t, "default", cfg.GetString("missing.key"))
	})

	t.Run("watch", func(t *testing.T) {
		watchKey := "watch.key"
		watchCalled := false

		// Set initial value
		cfg.Set(watchKey, "test")

		// Add watcher
		cfg.Watch(watchKey, func(value any) {
			watchCalled = true
			assert.Equal(t, "test", value)
		})

		assert.True(t, watchCalled)
	})
}
