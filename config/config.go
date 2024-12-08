package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// Config interface defines methods for configuration management
type Config interface {
	// Core methods
	Get(key string) any
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetFloat64(key string) float64
	GetStringMap(key string) map[string]any
	GetStringSlice(key string) []string
	GetIntSlice(key string) []int
	IsSet(key string) bool
	AllSettings() map[string]any
	AllKeys() []string

	// Extended methods
	Set(key string, value any)
	SetDefault(key string, value any)
	LoadFromFile(path string, options ...Option) error
	LoadFromDB(db any, tableName string) error
	Reload() error
	Watch(key string, callback func(any))

	// Additional Viper Get methods
	GetDuration(key string) time.Duration
	GetTime(key string) time.Time
	GetUint(key string) uint
	GetUint32(key string) uint32
	GetUint64(key string) uint64
	GetInt32(key string) int32
	GetInt64(key string) int64
	GetSizeInBytes(key string) uint
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string

	// extract
	Unmarshal(rawVal any, opts ...viper.DecoderConfigOption) error
	UnmarshalKey(key string, rawVal any, opts ...viper.DecoderConfigOption) error
}

type Option func(*viperConfig)

// WithConfigType sets the config type (yaml, json, etc)
func WithConfigType(configType string) Option {
	return func(c *viperConfig) {
		c.SetConfigType(configType)
	}
}

// WithEnvPrefix sets the environment variables prefix
func WithEnvPrefix(prefix string) Option {
	return func(c *viperConfig) {
		c.SetEnvPrefix(prefix)
	}
}

// WithEnvKeyReplacer sets the environment key replacer
func WithEnvKeyReplacer(oldNew ...string) Option {
	return func(c *viperConfig) {
		c.SetEnvKeyReplacer(strings.NewReplacer(oldNew...))
	}
}

type viperConfig struct {
	*viper.Viper
	watchMu   sync.RWMutex
	watches   map[string][]func(any)
	lastState map[string]any
}

// NewConfig creates a new configuration instance
func NewConfig() Config {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	return &viperConfig{
		Viper:     v,
		watches:   make(map[string][]func(any)),
		lastState: make(map[string]any),
	}
}

func (c *viperConfig) LoadFromFile(path string, options ...Option) error {
	// Apply options
	for _, opt := range options {
		opt(c)
	}

	// If config type not set, infer from extension
	if c.ConfigFileUsed() == "" {
		ext := filepath.Ext(path)
		c.SetConfigType(strings.TrimPrefix(ext, "."))
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("config file not found: %s", path)
	}

	// Read config file
	c.SetConfigFile(path)
	if err := c.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Watch for changes if fsnotify is enabled
	c.WatchConfig()
	c.OnConfigChange(func(e fsnotify.Event) {
		if err := c.Reload(); err != nil {
			// Log error but don't fail
			fmt.Printf("failed to reload config: %v\n", err)
		}
	})

	// Update last state
	c.updateLastState()

	return nil
}

func (c *viperConfig) LoadFromDB(db any, tableName string) error {
	var data map[string]any
	var err error

	switch v := db.(type) {
	case *sql.DB:
		data, err = loadFromSQL(v, tableName)
	case *gorm.DB:
		data, err = loadFromGorm(v, tableName)
	default:
		return fmt.Errorf("unsupported database type: %T", db)
	}

	if err != nil {
		return err
	}

	// Set all values from database
	for key, value := range data {
		c.Set(key, value)
	}

	// Update last state
	c.updateLastState()

	return nil
}

func (c *viperConfig) Reload() error {
	if err := c.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to reload config: %w", err)
	}

	// Check for changes and notify watchers
	c.watchMu.RLock()
	defer c.watchMu.RUnlock()

	for key, callbacks := range c.watches {
		oldValue := c.lastState[key]
		newValue := c.Get(key)
		if !reflect.DeepEqual(oldValue, newValue) {
			for _, callback := range callbacks {
				callback(newValue)
			}
		}
	}

	// Update last state
	c.updateLastState()

	return nil
}

func (c *viperConfig) updateLastState() {
	c.watchMu.Lock()
	defer c.watchMu.Unlock()

	for key := range c.watches {
		c.lastState[key] = c.Get(key)
	}
}

func (c *viperConfig) Watch(key string, callback func(any)) {
	c.watchMu.Lock()
	defer c.watchMu.Unlock()

	c.watches[key] = append(c.watches[key], callback)
	callback(c.Get(key))
}

func (c *viperConfig) Set(key string, value any) {
	c.Viper.Set(key, value)
}

func (c *viperConfig) SetDefault(key string, value any) {
	c.Viper.SetDefault(key, value)
}

func loadFromSQL(db *sql.DB, tableName string) (map[string]any, error) {
	// Create table if not exists
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			key_name VARCHAR(255) PRIMARY KEY,
			value TEXT NOT NULL
		)
	`, tableName)

	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	// Query all config entries
	query = fmt.Sprintf("SELECT key_name, value FROM %s", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query config: %w", err)
	}
	defer rows.Close()

	result := make(map[string]any)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Try to unmarshal JSON value
		var jsonValue any
		if err := json.Unmarshal([]byte(value), &jsonValue); err == nil {
			result[key] = jsonValue
		} else {
			result[key] = value
		}
	}

	return result, nil
}

func loadFromGorm(db *gorm.DB, tableName string) (map[string]any, error) {
	// Create table if not exists
	type ConfigEntry struct {
		KeyName string `gorm:"column:key_name;primaryKey"`
		Value   string `gorm:"column:value"`
	}

	if err := db.Table(tableName).AutoMigrate(&ConfigEntry{}); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	// Query all config entries
	var entries []ConfigEntry
	if err := db.Table(tableName).Find(&entries).Error; err != nil {
		return nil, fmt.Errorf("failed to query config: %w", err)
	}

	result := make(map[string]any)
	for _, entry := range entries {
		// Try to unmarshal JSON value
		var jsonValue any
		if err := json.Unmarshal([]byte(entry.Value), &jsonValue); err == nil {
			result[entry.KeyName] = jsonValue
		} else {
			result[entry.KeyName] = entry.Value
		}
	}

	return result, nil
}
