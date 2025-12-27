package main

import (
	"context"
	"fmt"
	"time"
)

// === CONTEXT PACKAGE EXAMPLES ===
//
// The context package provides:
// 1. Cancellation signals (stop all operations)
// 2. Deadlines (operations must complete by time T)
// 3. Timeouts (operations must complete in duration D)
// 4. Values (passing request-scoped data)
//
// Context is the idiomatic way to handle timeouts and cancellation in Go.

// Example 1: Background context (root context)
func example1BasicContext() {
	fmt.Println("=== Example 1: Background Context ===")

	// Background creates an empty context (root of context tree)
	ctx := context.Background()

	fmt.Printf("Context type: %T\n", ctx)
	fmt.Printf("Deadline: %v\n", ctx.Deadline())
	fmt.Printf("Done: %v\n", ctx.Done())
	fmt.Println()
}

// Example 2: Timeout context
func example2TimeoutContext() {
	fmt.Println("=== Example 2: Timeout Context ===")

	// Create context that times out after 100ms
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()  // Always cancel to clean up

	// Simulate work
	select {
	case <-time.After(50 * time.Millisecond):
		fmt.Println("Work completed successfully")
	case <-ctx.Done():
		fmt.Printf("Context cancelled: %v\n", ctx.Err())
	}

	// Try again with work taking longer than timeout
	ctx, cancel = context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	select {
	case <-time.After(100 * time.Millisecond):
		fmt.Println("Work completed")
	case <-ctx.Done():
		fmt.Printf("Context timed out: %v\n", ctx.Err())
	}
	fmt.Println()
}

// Example 3: Deadline context
func example3DeadlineContext() {
	fmt.Println("=== Example 3: Deadline Context ===")

	// Set absolute deadline (10 seconds from now)
	deadline := time.Now().Add(10 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// Check if context has deadline
	if d, ok := ctx.Deadline(); ok {
		fmt.Printf("Deadline: %v\n", d)
		fmt.Printf("Time until deadline: %v\n", time.Until(d))
	}
	fmt.Println()
}

// Example 4: Cancellation context
func example4CancellationContext() {
	fmt.Println("=== Example 4: Cancellation Context ===")

	// Create context that can be manually cancelled
	ctx, cancel := context.WithCancel(context.Background())

	// Start a worker goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Worker: Received cancellation signal")
				return
			default:
				fmt.Println("Worker: Doing work...")
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// Let it work for a bit
	time.Sleep(250 * time.Millisecond)

	// Cancel it
	fmt.Println("Main: Cancelling context")
	cancel()

	// Give goroutine time to clean up
	time.Sleep(100 * time.Millisecond)
	fmt.Println()
}

// Example 5: Context values
func example5ContextValues() {
	fmt.Println("=== Example 5: Context Values ===")

	// Create context with values
	ctx := context.Background()
	ctx = context.WithValue(ctx, "user_id", "12345")
	ctx = context.WithValue(ctx, "request_id", "req-abc-123")

	// Retrieve values
	userID := ctx.Value("user_id")
	requestID := ctx.Value("request_id")
	missing := ctx.Value("missing_key")

	fmt.Printf("User ID: %v\n", userID)
	fmt.Printf("Request ID: %v\n", requestID)
	fmt.Printf("Missing key: %v\n", missing)
	fmt.Println()
}

// Example 6: Propagating context through functions
func example6ContextPropagation() {
	fmt.Println("=== Example 6: Context Propagation ===")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Pass context to function
	result := doWork(ctx, "process A")
	fmt.Printf("Result: %v\n", result)

	// Try with shorter timeout
	ctx, cancel = context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	result = doWork(ctx, "process B")
	fmt.Printf("Result: %v\n", result)
	fmt.Println()
}

func doWork(ctx context.Context, name string) string {
	select {
	case <-time.After(100 * time.Millisecond):
		return fmt.Sprintf("%s completed", name)
	case <-ctx.Done():
		return fmt.Sprintf("%s cancelled: %v", name, ctx.Err())
	}
}

// Example 7: Combining contexts (nesting)
func example7ContextNesting() {
	fmt.Println("=== Example 7: Context Nesting ===")

	// Parent context with 500ms timeout
	parentCtx, parentCancel := context.WithTimeout(
		context.Background(),
		500*time.Millisecond,
	)
	defer parentCancel()

	// Child context with shorter timeout (inherits cancellation)
	childCtx, childCancel := context.WithTimeout(parentCtx, 200*time.Millisecond)
	defer childCancel()

	// Child context uses the shorter timeout
	fmt.Println("Child context will timeout after 200ms (shorter of two)")

	select {
	case <-time.After(300 * time.Millisecond):
		fmt.Println("Work completed")
	case <-childCtx.Done():
		fmt.Printf("Child context timed out: %v\n", childCtx.Err())
	}
	fmt.Println()
}

// Example 8: Concurrent goroutines with shared context
func example8ConcurrentWorkers() {
	fmt.Println("=== Example 8: Concurrent Workers ===")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	numWorkers := 3

	// Start multiple workers
	for i := 0; i < numWorkers; i++ {
		go worker(ctx, i)
	}

	// Wait for cleanup
	time.Sleep(500 * time.Millisecond)
	fmt.Println()
}

func worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d: Shutting down\n", id)
			return
		case <-time.After(100 * time.Millisecond):
			fmt.Printf("Worker %d: Working...\n", id)
		}
	}
}

// Example 9: Timeout for HTTP-like requests
func example9SimulatedHTTPRequest() {
	fmt.Println("=== Example 9: Simulated HTTP Request ===")

	// Simulate making requests with timeouts
	makeRequest("http://fast-service", 50*time.Millisecond, 200*time.Millisecond)
	makeRequest("http://slow-service", 150*time.Millisecond, 100*time.Millisecond)
}

func makeRequest(url string, duration time.Duration, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fmt.Printf("Requesting %s (timeout: %v)\n", url, timeout)

	select {
	case <-time.After(duration):
		fmt.Printf("  ✓ Response received after %v\n", duration)
	case <-ctx.Done():
		fmt.Printf("  ✗ Request failed: %v\n", ctx.Err())
	}
}

// Example 10: Checking context before expensive operations
func example10ContextAwareness() {
	fmt.Println("=== Example 10: Context Awareness ===")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	fmt.Println("Before checking context...")
	time.Sleep(150 * time.Millisecond)

	// Check if context is already cancelled before doing expensive work
	if err := ctx.Err(); err != nil {
		fmt.Printf("Context error: %v (skipping expensive operation)\n", err)
		return
	}

	fmt.Println("Doing expensive operation...")
}

// Example 11: Context with values for request tracking
func example11RequestTracking() {
	fmt.Println("=== Example 11: Request Tracking ===")

	// Create context with request metadata
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "req-123")
	ctx = context.WithValue(ctx, "user_id", "user-456")
	ctx = context.WithValue(ctx, "trace_id", "trace-789")

	// Pass to handler
	handleRequest(ctx)
}

func handleRequest(ctx context.Context) {
	requestID := ctx.Value("request_id")
	userID := ctx.Value("user_id")
	traceID := ctx.Value("trace_id")

	fmt.Printf("Handling request:\n")
	fmt.Printf("  Request ID: %v\n", requestID)
	fmt.Printf("  User ID: %v\n", userID)
	fmt.Printf("  Trace ID: %v\n", traceID)
}

// Example 12: Best practices
func example12BestPractices() {
	fmt.Println("=== Example 12: Best Practices ===")

	// ✓ GOOD: Use context.Background() as root
	ctx := context.Background()

	// ✓ GOOD: Always defer cancel()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// ✓ GOOD: Pass context as first parameter
	func(ctx context.Context) {
		fmt.Println("Context passed as first parameter")
	}(ctx)

	// ✓ GOOD: Don't pass nil context
	// func(ctx context.Context) { /* use ctx */ }(nil)  // ✗ Bad

	// ✓ GOOD: Store values with package-specific keys
	type key string
	const userIDKey key = "user_id"
	ctx = context.WithValue(ctx, userIDKey, "123")
	_ = ctx.Value(userIDKey)

	fmt.Println("Best practices demonstrated")
	fmt.Println()
}

func main() {
	example1BasicContext()
	example2TimeoutContext()
	example3DeadlineContext()
	example4CancellationContext()
	example5ContextValues()
	example6ContextPropagation()
	example7ContextNesting()
	example8ConcurrentWorkers()
	example9SimulatedHTTPRequest()
	example10ContextAwareness()
	example11RequestTracking()
	example12BestPractices()

	fmt.Println("=== KEY TAKEAWAYS ===")
	fmt.Println("1. context.Background() creates root context")
	fmt.Println("2. WithTimeout() for deadline-based cancellation")
	fmt.Println("3. WithCancel() for manual cancellation")
	fmt.Println("4. WithValue() for request-scoped data")
	fmt.Println("5. Always defer cancel() to clean up")
	fmt.Println("6. Pass context as first parameter")
	fmt.Println("7. Never pass nil context")
	fmt.Println("8. Use context for timeouts, not time.Sleep()")
	fmt.Println("9. Contexts propagate to child goroutines")
	fmt.Println("10. Check context.Err() before expensive operations")
}
