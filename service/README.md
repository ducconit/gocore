# Service Package

The service package provides a robust framework for managing different types of services in your Go applications.

## Features

- HTTP Service Management
- WebSocket Support (Pusher API v7 Compatible)
- Graceful Start/Stop
- Health Check Support
- Middleware Integration
- Logging Integration

## Components

### HTTP Service

The HTTP service provides a foundation for building REST APIs and web applications.

```go
import "github.com/ducconit/gocore/service/http"

// Create a new HTTP service
srv := http.NewHTTPService(
    http.WithAddr(":8080"),
    http.WithLogger(logger),
    http.WithMiddleware(middleware.Cors()),
)

// Start the service
if err := srv.Start(); err != nil {
    log.Fatal(err)
}
```

### WebSocket Service

WebSocket service with Pusher API v7 compatibility for real-time applications.

```go
import "github.com/ducconit/gocore/service/websocket"

// Create a new WebSocket service
ws := websocket.NewPusherService(
    websocket.WithAddr(":6001"),
    websocket.WithAppKey("your-app-key"),
    websocket.WithSecret("your-secret"),
)

// Start the service
if err := ws.Start(); err != nil {
    log.Fatal(err)
}
```

## Service Options

### HTTP Service Options

| Option | Description | Default |
|--------|-------------|---------|
| WithAddr | Set the server address | ":8080" |
| WithLogger | Set custom logger | nil |
| WithMiddleware | Add middleware | nil |
| WithTLS | Enable TLS | false |
| WithCertFile | Set TLS cert file | "" |
| WithKeyFile | Set TLS key file | "" |

### WebSocket Service Options

| Option | Description | Default |
|--------|-------------|---------|
| WithAddr | Set the WebSocket server address | ":6001" |
| WithAppKey | Set Pusher app key | Required |
| WithSecret | Set Pusher secret | Required |
| WithSSL | Enable SSL for WebSocket | false |
| WithCapacity | Set channel capacity | 100 |

## Examples

### Basic HTTP Service

```go
package main

import (
    "github.com/ducconit/gocore/service"
    "github.com/ducconit/gocore/service/http"
)

func main() {
    // Create HTTP service
    srv := http.NewHTTPService(
        http.WithAddr(":8080"),
    )

    // Add routes
    srv.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // Start service
    if err := srv.Start(); err != nil {
        log.Fatal(err)
    }
}
```

### WebSocket Service with Pusher Compatibility

```go
package main

import (
    "github.com/ducconit/gocore/service/websocket"
)

func main() {
    // Create WebSocket service
    ws := websocket.NewPusherService(
        websocket.WithAddr(":6001"),
        websocket.WithAppKey("app-key"),
        websocket.WithSecret("secret"),
    )

    // Add event handlers
    ws.OnConnection(func(client *websocket.Client) {
        log.Printf("New client connected: %s", client.ID)
    })

    // Start service
    if err := ws.Start(); err != nil {
        log.Fatal(err)
    }
}
```
