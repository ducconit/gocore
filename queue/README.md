# Queue Package

The queue package provides a robust message queue system with support for multiple backends and message patterns.

## Features

- Multiple Queue Backends
  - In-Memory Queue
  - Redis Queue
  - RabbitMQ Support
- Message Patterns
  - Publish/Subscribe
  - Work Queue
  - RPC Pattern
- Message Persistence
- Dead Letter Queue
- Message Retry
- Priority Queue

## Usage

### Basic Queue

```go
import "github.com/ducconit/gocore/queue"

// Create new queue
q := queue.New(
    queue.WithType(queue.TypeMemory),
    queue.WithCapacity(1000),
)

// Publish message
err := q.Publish("topic", []byte("message"))

// Subscribe to messages
sub := q.Subscribe("topic")
for msg := range sub.Messages() {
    fmt.Printf("Received: %s\n", string(msg.Data))
}
```

### Redis Queue

```go
// Create Redis queue
q := queue.New(
    queue.WithType(queue.TypeRedis),
    queue.WithRedisAddr("localhost:6379"),
)
```

### Work Queue Pattern

```go
// Create worker
worker := q.Worker("queue-name")
worker.Process(func(msg *queue.Message) error {
    // Process message
    return nil
})
```

## Queue Interface

```go
type Queue interface {
    Publish(topic string, data []byte) error
    Subscribe(topic string) Subscription
    Worker(queue string) Worker
    Close() error
}

type Subscription interface {
    Messages() <-chan *Message
    Unsubscribe() error
}

type Worker interface {
    Process(handler HandlerFunc) error
    Stop() error
}
```

## Options

### General Options

| Option | Description | Default |
|--------|-------------|---------|
| WithType | Queue backend type | TypeMemory |
| WithCapacity | Queue capacity | 1000 |
| WithRetryLimit | Max retry attempts | 3 |
| WithRetryDelay | Delay between retries | 5s |

### Redis Options

| Option | Description | Default |
|--------|-------------|---------|
| WithRedisAddr | Redis address | "localhost:6379" |
| WithRedisPassword | Redis password | "" |
| WithRedisDB | Redis database | 0 |

## Examples

### Priority Queue

```go
// Create priority queue
q := queue.New(
    queue.WithPriority(),
    queue.WithPriorityLevels(3),
)

// Publish with priority
q.PublishWithPriority("topic", []byte("high"), queue.PriorityHigh)
q.PublishWithPriority("topic", []byte("low"), queue.PriorityLow)
```

### Dead Letter Queue

```go
// Create queue with DLQ
q := queue.New(
    queue.WithDeadLetter("dlq"),
    queue.WithRetryLimit(3),
)

// Process messages with DLQ handling
worker := q.Worker("queue")
worker.Process(func(msg *queue.Message) error {
    if err := processMessage(msg); err != nil {
        return err // Message will go to DLQ after retry limit
    }
    return nil
})
```

### RPC Pattern

```go
// Server
rpcServer := q.RPCServer("calc")
rpcServer.Handle("add", func(req []byte) ([]byte, error) {
    var nums struct{ A, B int }
    json.Unmarshal(req, &nums)
    result := nums.A + nums.B
    return json.Marshal(result)
})

// Client
rpcClient := q.RPCClient("calc")
response, err := rpcClient.Call("add", request)
```

## Best Practices

1. Choose appropriate queue type for your use case
2. Implement proper error handling
3. Use dead letter queues for failed messages
4. Monitor queue size and performance
5. Implement proper retry strategies
6. Use appropriate message serialization

## Message Patterns

### Publish/Subscribe
- One-to-many message distribution
- Topics based routing
- Non-persistent messages

### Work Queue
- Task distribution among workers
- Message persistence
- At-least-once delivery

### RPC
- Request/Response pattern
- Synchronous communication
- Timeout handling

## Error Handling

```go
// Error types
var (
    ErrQueueFull    = errors.New("queue is full")
    ErrQueueClosed  = errors.New("queue is closed")
    ErrInvalidTopic = errors.New("invalid topic")
)

// Handle errors
if err := q.Publish("topic", data); err != nil {
    switch err {
    case queue.ErrQueueFull:
        // Handle full queue
    case queue.ErrQueueClosed:
        // Handle closed queue
    default:
        // Handle other errors
    }
}
```
