# Errors Package

The errors package provides enhanced error handling capabilities with error wrapping, stack traces, and error types.

## Features

- Error Wrapping
- Stack Traces
- Error Types
- Error Codes
- HTTP Error Integration
- Error Context
- Localization Support

## Usage

### Basic Usage

```go
import "github.com/ducconit/gocore/errors"

// Create new error
err := errors.New("something went wrong")

// Create error with code
err := errors.NewWithCode(404, "not found")

// Wrap error
err = errors.Wrap(err, "failed to process request")
```

### Error Types

```go
// Predefined error types
var (
    ErrNotFound     = errors.NewType("not_found", 404)
    ErrUnauthorized = errors.NewType("unauthorized", 401)
    ErrBadRequest   = errors.NewType("bad_request", 400)
)

// Create error from type
err := ErrNotFound.New("user not found")
```

## Error Interface

```go
type Error interface {
    error
    Code() int
    Message() string
    Cause() error
    Stack() string
    Context() map[string]interface{}
}
```

## Examples

### Error with Context

```go
err := errors.New("database error").
    WithContext("table", "users").
    WithContext("operation", "insert")
```

### HTTP Error Handling

```go
func ErrorHandler(c *gin.Context) {
    c.Next()
    
    if len(c.Errors) > 0 {
        err := c.Errors.Last()
        if e, ok := err.Err.(errors.Error); ok {
            c.JSON(e.Code(), gin.H{
                "error": e.Message(),
                "code":  e.Code(),
            })
            return
        }
        c.JSON(500, gin.H{"error": err.Error()})
    }
}
```

### Stack Trace

```go
err := errors.New("something went wrong")
fmt.Println(err.Stack())
```

### Error Chain

```go
err1 := errors.New("original error")
err2 := errors.Wrap(err1, "wrapped once")
err3 := errors.Wrap(err2, "wrapped twice")

// Print full error chain
fmt.Println(err3.Error())
```

## Best Practices

1. Use appropriate error types
2. Include relevant context
3. Wrap errors with meaningful messages
4. Handle errors at appropriate levels
5. Use stack traces for debugging
6. Implement proper error logging

## Error Types

### Standard Error Types

```go
var (
    ErrNotFound     = errors.NewType("not_found", 404)
    ErrUnauthorized = errors.NewType("unauthorized", 401)
    ErrBadRequest   = errors.NewType("bad_request", 400)
    ErrInternal     = errors.NewType("internal", 500)
    ErrValidation   = errors.NewType("validation", 422)
)
```

### Custom Error Types

```go
// Define custom error type
var ErrDatabaseConnection = errors.NewType("database_connection", 503)

// Use custom error type
err := ErrDatabaseConnection.New("failed to connect to database")
```

## Error Context

### Adding Context

```go
err := errors.New("validation failed").
    WithContext("field", "email").
    WithContext("value", "invalid@email").
    WithContext("rule", "email_format")
```

### Retrieving Context

```go
if e, ok := err.(errors.Error); ok {
    ctx := e.Context()
    field := ctx["field"].(string)
    value := ctx["value"].(string)
}
```

## Error Codes

### Standard HTTP Status Codes

```go
const (
    StatusBadRequest          = 400
    StatusUnauthorized        = 401
    StatusForbidden          = 403
    StatusNotFound           = 404
    StatusMethodNotAllowed   = 405
    StatusConflict           = 409
    StatusUnprocessableEntity = 422
    StatusInternalError      = 500
    StatusServiceUnavailable = 503
)
```

## Localization

### Error Messages

```go
// Register translations
errors.RegisterTranslation("en", map[string]string{
    "not_found": "Resource not found: %s",
    "unauthorized": "Unauthorized access",
})

// Create localized error
err := errors.NewLocalized("not_found", "en", "user")
```

## Integration

### Gin Integration

```go
func ErrorMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            if e, ok := err.Err.(errors.Error); ok {
                c.JSON(e.Code(), gin.H{
                    "error": e.Message(),
                    "code": e.Code(),
                    "context": e.Context(),
                })
                return
            }
            c.JSON(500, gin.H{"error": err.Error()})
        }
    }
}
