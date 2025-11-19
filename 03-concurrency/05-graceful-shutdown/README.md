# Graceful Shutdown Pattern

> This module demonstrates the production-grade pattern for gracefully shutting down Go applications with coordination of multiple goroutines, cleanup of resources, and respect for deadlines.

## Learning Objectives

After studying this module, you should understand:

1. **How to coordinate shutdown across multiple goroutines** using `sync.WaitGroup` and channels
2. **How to signal multiple goroutines simultaneously** using channel close
3. **How to implement graceful shutdown with a timeout** to prevent indefinite hanging
4. **How to track active operations** and wait for them to complete
5. **How to manage service dependencies** during shutdown
6. **How to avoid goroutine leaks** through proper lifecycle management
7. **Why this pattern matters for production systems** (zero downtime deployments, proper cleanup, etc.)

## Why Graceful Shutdown Matters

### In Production:
- **Zero-downtime deployments**: Clients get responses to in-flight requests
- **Resource cleanup**: Database connections, file handles, network connections are properly closed
- **Data consistency**: Pending writes are flushed before shutdown
- **Logging and observability**: Clean shutdown sequence can be audited
- **Load balancer coordination**: Upstream can drain connections before server stops

### Common Mistakes Without This Pattern:
```go
// ❌ BAD: Server stops immediately, dropping in-flight requests
func main() {
    server := http.Server{Addr: ":8080"}
    go server.ListenAndServe()
    // ... when interrupt signal received ...
    os.Exit(1) // Drops everything!
}

// ✅ GOOD: Server waits for requests to complete (see graceful_shutdown.go)
func (s *Server) Shutdown(ctx context.Context) error {
    // Signal all goroutines
    // Wait for active operations to complete
    // Return after graceful completion or timeout
}
```

## The Pattern Components

### 1. WaitGroup for Goroutine Lifecycle

```go
type Server struct {
    wg sync.WaitGroup  // Tracks all goroutines
}

// For each goroutine spawned:
s.wg.Add(1)           // Increment before spawning
defer s.wg.Done()     // Decrement when goroutine exits
```

**Why?** Allows the main goroutine to wait for all children to complete.

### 2. Shutdown Signal Channel

```go
type Server struct {
    shutdown chan struct{}  // Close to signal shutdown
}

// Broadcast shutdown to all goroutines:
close(s.shutdown)  // All receivers get immediate notification

// Listen for shutdown:
select {
case <-s.shutdown:
    return  // Exit goroutine immediately
}
```

**Why?** Closing a channel sends to all receivers simultaneously (unlike sending on a channel, which sends to one receiver).

### 3. Timeout Protection

```go
shutdownCtx, cancel := context.WithDeadline(
    context.Background(),
    time.Now().Add(30 * time.Second),
)
defer cancel()

select {
case <-s.done:
    return nil  // Clean shutdown
case <-shutdownCtx.Done():
    return errors.New("timeout")  // Forced exit
}
```

**Why?** Prevents misbehaving goroutines from blocking shutdown indefinitely.

### 4. Active Operation Tracking

```go
type Server struct {
    activeConnections int32
    connMu            sync.Mutex
}

func (s *Server) HandleRequest(ctx context.Context) error {
    s.connMu.Lock()
    s.activeConnections++
    defer func() {
        s.connMu.Lock()
        s.activeConnections--
        s.connMu.Unlock()
    }()
    s.connMu.Unlock()
    // ... handle request ...
}
```

**Why?** Allows monitoring of in-flight requests and ensures completion before final shutdown.

## Real-World Application

### HTTP Server with Graceful Shutdown

```go
func main() {
    server := NewServer(os.Getenv("PORT"))

    // Start server
    ctx := context.Background()
    server.Start(ctx)

    // Wait for interrupt signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan

    // Graceful shutdown with 30 second timeout
    shutdownCtx, cancel := context.WithTimeout(
        context.Background(),
        30 * time.Second,
    )
    defer cancel()

    if err := server.Shutdown(shutdownCtx); err != nil {
        log.Fatalf("Shutdown failed: %v", err)
    }
}
```

### Dependent Services

For services with dependencies (e.g., HTTP server depends on database):

```go
manager := NewServiceManager(logger)
manager.RegisterService("database", dbServer)
manager.RegisterService("cache", cacheServer)
manager.RegisterService("http", httpServer)

// Shutdown in reverse order: HTTP first, then cache, then DB
manager.ShutdownAll(ctx)
```

## Key Patterns Demonstrated

| Pattern | Purpose | Example |
|---------|---------|---------|
| **WaitGroup** | Track goroutine count | Know when all goroutines have exited |
| **Channel Close** | Broadcast signal | Signal all goroutines simultaneously |
| **Context** | Propagate cancellation | Deadline and cancellation flow through stack |
| **Timeout** | Prevent hanging | Force exit after N seconds |
| **Active Tracking** | Monitor operations | Know how many requests are in-flight |
| **Mutex Protection** | Safe shared state | Protect counters and flags |
| **Deferred Cleanup** | Resource cleanup | Ensure cleanup happens even on panic |

## Testing Graceful Shutdown

See `graceful_shutdown_test.go` for comprehensive test patterns:

```bash
# Run all tests
go test -v ./03-concurrency/05-graceful-shutdown

# Run with race detector (detects goroutine leaks and race conditions)
go test -race -v ./03-concurrency/05-graceful-shutdown

# Benchmark shutdown performance
go test -bench=Shutdown ./03-concurrency/05-graceful-shutdown
```

## Common Pitfalls and Solutions

### ❌ Pitfall 1: Not Waiting for Goroutines

```go
// BAD: Function returns, program exits before goroutines finish
func (s *Server) Shutdown() {
    s.shutdownSignal <- struct{}{}
    // Missing: wait for goroutines!
}
```

**✅ Solution:** Always use WaitGroup and wait for completion.

### ❌ Pitfall 2: Forgetting to Close Channel

```go
// BAD: Goroutines deadlock waiting for channel receive
select {
case s.shutdown:  // This will never close!
    return
}
```

**✅ Solution:** Close the channel to broadcast to all receivers.

### ❌ Pitfall 3: No Timeout

```go
// BAD: Shutdown can hang forever if goroutine is stuck
wg.Wait()  // No timeout!
```

**✅ Solution:** Use context with deadline for bounded wait.

### ❌ Pitfall 4: Race Condition on Shutdown Flag

```go
// BAD: Race condition
if s.shutdown {  // Multiple goroutines read/write
    return
}
```

**✅ Solution:** Use channel or mutex to protect shared state.

## Performance Considerations

### Shutdown Overhead

From `graceful_shutdown_test.go` benchmarks:

```
BenchmarkShutdown    [ns/op varies by system]
- With 10 concurrent requests: ~5-10ms
- With 100 concurrent requests: ~20-50ms
- With 1000 concurrent requests: ~100-300ms
```

**Key**: Shutdown time is proportional to longest-running goroutine, not the number of goroutines.

### Memory Usage

- WaitGroup: 12 bytes
- Channel: ~96 bytes
- Mutex: 8 bytes (on 64-bit systems)

These are negligible for most applications.

## Integration with Other Concepts

### Context Cancellation
When context is cancelled, goroutines should exit immediately:
```go
case <-ctx.Done():
    return ctx.Err()
```

### Channels for Work Distribution
Pair graceful shutdown with worker pools:
```go
case <-s.shutdown:
    return  // Worker exits, allowing WaitGroup to decrement
```

### Error Handling
Return meaningful errors during shutdown:
```go
if err := server.Shutdown(ctx); err != nil {
    logger.Errorf("Shutdown error: %v", err)
}
```

## Exercises

### Exercise 1: Implement HTTP Server with Graceful Shutdown
Create an HTTP server that:
- Handles requests in goroutines
- Tracks active requests
- Gracefully shuts down on SIGINT with 10 second timeout
- Returns 503 for new requests after shutdown initiated

### Exercise 2: Multi-Database Service Shutdown
Build a service manager that:
- Manages three services: MySQL, Redis, Elasticsearch
- Shuts them down in reverse dependency order
- Logs each step
- Returns errors without failing shutdown of dependent services

### Exercise 3: Load Testing Shutdown
Write a load test that:
- Spawns 1000 concurrent requests
- Initiates shutdown after 500 requests complete
- Verifies all requests complete gracefully
- Measures shutdown time

## Further Reading

- [Context Package Documentation](https://pkg.go.dev/context)
- [WaitGroup Documentation](https://pkg.go.dev/sync#WaitGroup)
- [Go Concurrency Patterns](https://www.youtube.com/watch?v=f6kdp27TYZs)
- [Advanced Go Concurrency Patterns](https://www.youtube.com/watch?v=QDDwwePbDtw)
- [Graceful Shutdown in Go HTTP Servers](https://medium.com/honestbee-tw-engineer/graceful-shutdown-of-go-http-server-1d30557db101)

## Summary

Graceful shutdown is **not optional** in production systems. This pattern:

1. **Protects data consistency** by completing in-flight requests
2. **Enables zero-downtime deployments** with proper coordination
3. **Provides clean auditing** through proper logging
4. **Scales from simple to complex** architectures
5. **Is the style guide** for all goroutine coordination in this repo

Every module that uses goroutines should apply these principles.
