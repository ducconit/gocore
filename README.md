# GoCore

A lightweight and flexible core library for Go applications by DNT.

## Features

### Configuration Management
- Flexible configuration loading from multiple sources
- Type-safe configuration access
- Default value support
- Environment variable support
- Comprehensive configuration methods (Get, GetString, GetInt, etc.)

### Service Management
- HTTP service management with graceful start/stop
- Health check support
- Flexible service configuration
- Logging integration
- Signal handling (SIGTERM, Interrupt)

### Utilities
- OS signal handling
- Graceful shutdown support
- Cross-platform compatibility

## Installation

```bash
go get github.com/ducconit/gocore
```

## Usage Examples

### Configuration Management
```go
import "github.com/ducconit/gocore/config"

// Create new configuration
cfg := config.NewConfig()

// Load configuration from file
err := cfg.LoadFromFile("config.yaml")
if err != nil {
    log.Fatal(err)
}

// Access configuration values
port := cfg.GetInt("server.port")
host := cfg.GetString("server.host")
```

### HTTP Service
```go
import "github.com/ducconit/gocore/service"

// Create new HTTP service
svc := service.NewHTTPService("api",
    service.WithAddress(":8080"),
    service.WithHandler(handler),
)

// Start service
err := svc.Start(context.Background())
if err != nil {
    log.Fatal(err)
}

// Graceful shutdown
svc.Stop(context.Background())
```

### Signal Handling
```go
import "github.com/ducconit/gocore/utils"

// Register interrupt handler
utils.RegisterSignalInterruptHandler(func() {
    // Cleanup code here
})
```

## Requirements
- Go 1.23.2 or higher

## License
MIT License

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
