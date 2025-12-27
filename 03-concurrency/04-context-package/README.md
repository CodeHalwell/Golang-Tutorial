# Chapter 4: The Context Package - Timeouts, Cancellation, and Values

## Learning Objectives

1. Understand when and why to use context
2. Implement timeouts for long-running operations
3. Coordinate cancellation across goroutines
4. Pass request-scoped values safely
5. Propagate context through function calls
6. Avoid common context mistakes
7. Use context in real-world scenarios

---

## What is Context?

The `context` package provides:

1. **Cancellation**: Stop all operations (and propagate to child goroutines)
2. **Deadlines**: Operations must complete by time T
3. **Timeouts**: Operations must complete in duration D
4. **Values**: Pass request-scoped data safely

**Key Rule**: Always pass context as the first parameter of a function.

---

## Core Types

### context.Context interface

```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key interface{}) interface{}
}
```

**Methods:**
- `Deadline()`: Get the deadline, if any
- `Done()`: Channel that closes when context is cancelled
- `Err()`: Why context was cancelled (nil, context.Canceled, context.DeadlineExceeded)
- `Value(key)`: Get a value stored in context

---

## Creating Contexts

### 1. Background Context (Root)

```go
ctx := context.Background()
```

- Empty context, no deadline or values
- Used as root of context tree
- Never cancelled

### 2. Todo Context (Placeholder)

```go
ctx := context.TODO()
```

- Like Background, but signals "context not yet determined"
- Use when you're not sure what context to use
- Prefer Background in most cases

### 3. Timeout Context

```go
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel()  // CRITICAL: Always defer cancel()
```

- Context that expires after duration
- `cancel()` stops the context immediately
- Always defer cancel() to clean up

### 4. Deadline Context

```go
deadline := time.Now().Add(1 * time.Hour)
ctx, cancel := context.WithDeadline(parent, deadline)
defer cancel()
```

- Context expires at absolute time
- Use when you have a specific deadline
- Otherwise prefer WithTimeout

### 5. Cancellation Context

```go
ctx, cancel := context.WithCancel(parent)
defer cancel()
```

- Manual cancellation (no timeout)
- cancel() stops the context immediately
- Useful for explicit shutdown

### 6. Value Context

```go
ctx = context.WithValue(ctx, "user_id", "12345")
```

- Stores request-scoped values
- Immutable (creates new context)
- Don't overuse

---

## Using Channels from Context

### Wait for Context Cancellation

```go
select {
case <-ctx.Done():
    return ctx.Err()
case <-time.After(5 * time.Second):
    // Operation completed
}
```

### Check Error Before Returning

```go
if err := ctx.Err(); err != nil {
    // Already cancelled
    return err
}
```

---

## Practical Examples

### Example 1: HTTP Request with Timeout

```go
// Create context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Make request respecting timeout
req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
resp, err := http.DefaultClient.Do(req)
```

### Example 2: Database Query with Timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

// Database query respects timeout
result := db.QueryContext(ctx, "SELECT * FROM users")
```

### Example 3: Graceful Shutdown

```go
func gracefulShutdown(server *http.Server) {
    // Create shutdown context with 30-second timeout
    ctx, cancel := context.WithTimeout(
        context.Background(),
        30 * time.Second,
    )
    defer cancel()

    // Shutdown respects context deadline
    server.Shutdown(ctx)
}
```

### Example 4: Worker Pool with Cancellation

```go
func workers(ctx context.Context, numWorkers int) {
    var wg sync.WaitGroup

    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()

            for {
                select {
                case <-ctx.Done():
                    // Stop when context cancelled
                    return
                case work := <-taskChan:
                    processWork(work)
                }
            }
        }(i)
    }

    wg.Wait()
}
```

### Example 5: Timeout with Early Exit

```go
ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
defer cancel()

// Check before expensive operation
if err := ctx.Err(); err != nil {
    return fmt.Errorf("already cancelled: %w", err)
}

// Do expensive work
result := expensiveComputation(ctx)
```

---

## Context Values

### Good: Request-Scoped Data

```go
// Store request metadata
ctx = context.WithValue(ctx, "request_id", reqID)
ctx = context.WithValue(ctx, "user_id", userID)

// Retrieve in handlers
requestID := ctx.Value("request_id")
```

### Bad: Overusing Values

```go
// ✗ Don't: using context like a config object
ctx = context.WithValue(ctx, "database", db)
ctx = context.WithValue(ctx, "cache", cache)
ctx = context.WithValue(ctx, "logger", logger)

// This is what dependency injection is for!
// Pass these as function parameters instead
```

### Best: Type-Safe Value Keys

```go
type key string

const (
    requestIDKey key = "request_id"
    userIDKey    key = "user_id"
)

// Store
ctx = context.WithValue(ctx, requestIDKey, "abc-123")

// Retrieve
id := ctx.Value(requestIDKey).(string)
```

---

## Common Patterns

### Pattern 1: Timeout with Fallback

```go
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()

select {
case result := <-fastPath(ctx):
    return result
case <-ctx.Done():
    // Timeout, use fallback
    return cachedResult()
}
```

### Pattern 2: Deadline Propagation

```go
func service1(ctx context.Context) error {
    return service2(ctx)  // Context deadline propagates
}

func service2(ctx context.Context) error {
    return service3(ctx)  // Service3 respects deadline from service1
}
```

### Pattern 3: Multiple Timeouts

```go
// Outer timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// Inner timeout (shorter)
innerCtx, innerCancel := context.WithTimeout(ctx, 5*time.Second)
defer innerCancel()

// innerCtx will timeout at 5 seconds OR when ctx times out at 10 seconds
// (whichever comes first)
```

### Pattern 4: Cancel on First Error

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// If any goroutine encounters error, cancel all
go func() {
    if err := doWork(); err != nil {
        cancel()  // Stop all other work
    }
}()
```

---

## Gotchas and Common Mistakes

### ✗ Mistake 1: Not Deferring Cancel

```go
// BAD: Resource leak
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// ... forgot to cancel
```

**Fix**: Always defer cancel()

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()  // ✓ Guaranteed cleanup
```

### ✗ Mistake 2: Passing Nil Context

```go
// BAD: Panics
func handler(ctx context.Context) {
    handler(nil)  // ✗ Crashes
}
```

**Fix**: Use context.Background()

```go
handler(context.Background())  // ✓ Correct
```

### ✗ Mistake 3: Using context.TODO() Incorrectly

```go
// TODO should be temporary
ctx := context.TODO()
// ... but forgot to replace it with proper context
```

**Fix**: Replace with Background or specific context type

```go
ctx := context.Background()  // ✓ or WithTimeout, etc.
```

### ✗ Mistake 4: Not Checking Context

```go
// BAD: Ignores context
func process(ctx context.Context) {
    // Never checks ctx.Done()
    time.Sleep(10 * time.Second)  // Ignores timeout!
}
```

**Fix**: Respect the context

```go
func process(ctx context.Context) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-time.After(10 * time.Second):
        // Process completed
        return nil
    }
}
```

### ✗ Mistake 5: Storing non-request-scoped data in context

```go
// BAD: Global config in context
ctx = context.WithValue(ctx, "database", db)
```

**Fix**: Pass as function parameter

```go
func handler(ctx context.Context, db *Database) {
    // Use db directly
}
```

---

## Best Practices

1. **Always pass context as first parameter**
   ```go
   func MyFunction(ctx context.Context, arg1 string, arg2 int)
   ```

2. **Always defer cancel() when creating cancelable contexts**
   ```go
   ctx, cancel := context.WithTimeout(...)
   defer cancel()
   ```

3. **Never pass nil as context**
   ```go
   // Bad: function(nil)
   // Good: function(context.Background())
   ```

4. **Check context before expensive operations**
   ```go
   if err := ctx.Err(); err != nil {
       return err
   }
   ```

5. **Use context.Background() as root**
   ```go
   ctx := context.Background()  // ✓
   ctx := context.TODO()         // Only if unsure
   ```

6. **Only store request-scoped values**
   ```go
   // ✓ Good
   ctx = context.WithValue(ctx, "request_id", id)

   // ✗ Bad
   ctx = context.WithValue(ctx, "database", db)
   ```

7. **Use type-safe value keys**
   ```go
   type contextKey string
   const userIDKey contextKey = "user_id"
   ```

8. **Let context propagate naturally**
   ```go
   func outer(ctx context.Context) {
       inner(ctx)  // Pass to inner functions
   }
   ```

---

## Exercises

### Exercise 1: HTTP Client with Timeout
Create an HTTP client that respects context timeouts.

### Exercise 2: Worker Pool Shutdown
Implement a worker pool that shuts down when context is cancelled.

### Exercise 3: Database Batch Processing
Process database records with per-record timeout.

### Exercise 4: Request Tracking
Add request IDs to context and log them in handlers.

---

## Summary

The context package is **essential** for Go development:

- **Timeouts**: Prevent hanging operations
- **Cancellation**: Stop all work immediately
- **Propagation**: Deadlines flow through call stack
- **Values**: Pass request metadata safely

Every Go developer must master context.

---

## Next: Chapter 5 - Graceful Shutdown

We'll combine context, goroutines, and channels to implement production-grade graceful shutdown.
