# GoCore

A lightweight and flexible core library for Go applications by DNT.

## Features

### Configuration Management
- Multiple configuration sources (YAML, JSON, ENV)
- Type-safe configuration access
- Dynamic configuration updates
- Configuration validation
- Environment variable support

### Caching System
- Multiple cache backends (Memory, Redis)
- TTL support
- Cache tags
- Bulk operations
- Automatic key expiration

### Message Queue
- Multiple queue backends
- Message patterns (Pub/Sub, Work Queue, RPC)
- Dead letter queue
- Message retry
- Priority queue

### Logging System
- Multiple log levels
- Structured logging
- Output formatting
- Log rotation
- Context-aware logging

### Error Handling
- Error wrapping
- Stack traces
- Error types
- HTTP error integration
- Error context
- Localization support

## Installation

```bash
go get github.com/ducconit/gocore
```

## Quick Start

### Configuration

```go
import "github.com/ducconit/gocore/config"

func main() {
    // Create configuration
    cfg := config.New(
        config.WithFile("config.yaml"),
        config.WithEnv(),
    )

    // Access configuration
    port := cfg.GetInt("server.port", 8080)
    host := cfg.GetString("server.host", "localhost")
}
```

### Caching

```go
import "github.com/ducconit/gocore/cache"

func main() {
    // Create cache
    c := cache.NewRedisCache(
        cache.WithRedisAddr("localhost:6379"),
    )

    // Use cache
    c.Set("key", "value", 5*time.Minute)
    val, err := c.Get("key")
}
```

### Logging

```go
import "github.com/ducconit/gocore/logger"

func main() {
    // Create logger
    log := logger.New(
        logger.WithLevel(logger.InfoLevel),
        logger.WithOutput("path/to/file.log"),
    )

    // Log messages
    log.Info("Server starting", zap.String("port", "3000"))
}
```

### Error Handling

```go
import "github.com/ducconit/gocore/errors"

func main() {
    // Create error
    err := errors.NewWithCode(404, "user not found").
        WithContext("user_id", 123)

    // Handle error
    if e, ok := err.(errors.Error); ok {
        fmt.Printf("Code: %d, Message: %s\n", e.Code(), e.Message())
    }
}
```

## Documentation

Detailed documentation for each package can be found in their respective directories:

- [Configuration Package](config/README.md)
- [Cache Package](cache/README.md)
- [Queue Package](queue/README.md)
- [Logger Package](logger/README.md)
- [Errors Package](errors/README.md)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License

## Support

For support, please open an issue in the GitHub repository.

## Contributors

<!-- readme: contributors -start -->
<table>
	<tbody>
		<tr>
            <td align="center">
                <a href="https://github.com/ducconit">
                    <img src="https://avatars.githubusercontent.com/u/72369814?v=4" width="100;" alt="ducconit"/>
                    <br />
                    <sub><b>Duke</b></sub>
                </a>
            </td>
		</tr>
	<tbody>
</table>
<!-- readme: contributors -end -->
