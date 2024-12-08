# Logger Package

The logger package provides a flexible and extensible logging system for Go applications.

## Features

- Multiple Log Levels (Debug, Info, Warn, Error, Fatal)
- Structured Logging
- Output Formatting
- Log Rotation
- Custom Writers Support
- Context-aware Logging

## Usage

### Basic Usage

```go
import "github.com/ducconit/gocore/logger"

// Create a new logger
log := logger.New()

// Log messages
log.Info("Server starting...")
log.Debug("Debug message")
log.Error("Error occurred", logger.Fields{
    "error": err,
    "code":  500,
})
```

### With Configuration

```go
log := logger.New(
    logger.WithLevel(logger.InfoLevel),
    logger.WithFormat(logger.JSONFormat),
    logger.WithOutput(os.Stdout),
)
```

### Structured Logging

```go
log.Info("User logged in",
    logger.Fields{
        "user_id":    123,
        "ip":         "192.168.1.1",
        "user_agent": "Mozilla/5.0...",
    },
)
```

## Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| WithLevel | Set logging level | InfoLevel |
| WithFormat | Set output format (Text/JSON) | TextFormat |
| WithOutput | Set output writer | os.Stdout |
| WithTimeFormat | Set time format | RFC3339 |
| WithCaller | Include caller information | false |

## Log Levels

- Debug: Detailed information for debugging
- Info: General operational information
- Warn: Warning messages for potentially harmful situations
- Error: Error messages for serious problems
- Fatal: Critical errors that stop the program

## Examples

### File Logging

```go
file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
if err != nil {
    log.Fatal("Could not open log file", logger.Fields{"error": err})
}

log := logger.New(
    logger.WithOutput(file),
    logger.WithFormat(logger.JSONFormat),
)
```

### Custom Formatter

```go
log := logger.New(
    logger.WithFormatter(func(entry *logger.Entry) []byte {
        // Custom formatting logic
        return []byte(fmt.Sprintf("[%s] %s: %s\n",
            entry.Time.Format(time.RFC3339),
            entry.Level,
            entry.Message,
        ))
    }),
)
```

### Context-aware Logging

```go
ctx := context.Background()
ctx = logger.WithContext(ctx, logger.Fields{
    "request_id": "123",
    "user_id":    "456",
})

log.FromContext(ctx).Info("Processing request")
```

## Best Practices

1. Use appropriate log levels
2. Include relevant context in structured fields
3. Configure log rotation for production
4. Use context-aware logging for request tracing
5. Include error details in error logs