# Config Package

The config package provides a flexible configuration management system with support for multiple sources and formats.

## Features

- Multiple Configuration Sources
  - Files (YAML, JSON, TOML)
  - Environment Variables
  - Command Line Flags
- Dynamic Configuration Updates
- Type-safe Access
- Default Values
- Configuration Validation
- Nested Configuration Support

## Usage

### Basic Usage

```go
import "github.com/ducconit/gocore/config"

// Create new configuration
cfg := config.New()

// Load from file
err := cfg.LoadFile("config.yaml")
if err != nil {
    log.Fatal(err)
}

// Access values
port := cfg.GetInt("server.port")
host := cfg.GetString("server.host")
debug := cfg.GetBool("app.debug")
```

### Environment Variables

```go
// Load from environment
cfg.LoadEnv()

// Access with fallback
dbHost := cfg.GetString("DB_HOST", "localhost")
```

### Multiple Sources

```go
cfg := config.New(
    config.WithFile("config.yaml"),
    config.WithEnv(),
    config.WithFlags(),
)
```

## Configuration Methods

### Getters

```go
// Basic getters
str := cfg.GetString("key")
num := cfg.GetInt("key")
float := cfg.GetFloat64("key")
bool := cfg.GetBool("key")
duration := cfg.GetDuration("key")

// With default values
str := cfg.GetString("key", "default")
num := cfg.GetInt("key", 8080)

// Typed getters
var serverConfig ServerConfig
cfg.Get("server", &serverConfig)
```

### Setters

```go
cfg.Set("key", "value")
cfg.SetDefault("server.port", 8080)
```

## Configuration Structure

### YAML Example

```yaml
app:
  name: MyApp
  debug: true
  
server:
  host: localhost
  port: 8080
  
database:
  host: localhost
  port: 5432
  name: mydb
  user: user
  password: password
```

### Environment Variables

```bash
APP_NAME=MyApp
SERVER_PORT=8080
DB_HOST=localhost
```

## Examples

### Custom Configuration Type

```go
type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Name     string `yaml:"name"`
    User     string `yaml:"user"`
    Password string `yaml:"password"`
}

var dbConfig DatabaseConfig
cfg.Get("database", &dbConfig)
```

### Watch Configuration Changes

```go
cfg.OnChange(func(e config.ChangeEvent) {
    log.Printf("Config changed: %s = %v", e.Key, e.NewValue)
    
    // Reload specific components
    if e.Key == "database.host" {
        reconnectDatabase()
    }
})
```

### Validation

```go
type ServerConfig struct {
    Host string `yaml:"host" validate:"required"`
    Port int    `yaml:"port" validate:"required,min=1024,max=65535"`
}

var serverConfig ServerConfig
if err := cfg.GetValidated("server", &serverConfig); err != nil {
    log.Fatal(err)
}
```

## Best Practices

1. Use structured configuration
2. Implement validation for critical values
3. Provide sensible defaults
4. Use environment variables for sensitive data
5. Document configuration options
6. Handle configuration errors gracefully

## Configuration Hierarchy

The package follows this configuration hierarchy (highest to lowest priority):

1. Command Line Flags
2. Environment Variables
3. Configuration Files
4. Default Values

## Security Considerations

1. Never commit sensitive data in configuration files
2. Use environment variables for secrets
3. Implement proper access controls
4. Validate configuration values
5. Use secure storage for sensitive data
