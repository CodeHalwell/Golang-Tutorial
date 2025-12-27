# Chapter 1: Goroutines Basics

## Learning Objectives

1. Understand what goroutines are and how they differ from threads
2. Create and run goroutines
3. Understand the M:N threading model
4. Manage goroutine lifetime
5. Avoid race conditions with synchronization
6. Use WaitGroups for coordination
7. Handle goroutine failures and panics
8. Understand goroutine scheduling
9. Debug common goroutine issues
10. Best practices for goroutine management

---

## What Are Goroutines?

### Definition

A **goroutine** is a lightweight thread managed by the Go runtime.

```
┌─────────────────────────────────────────┐
│            Your Go Program              │
│                                         │
│  Main Goroutine    Other Goroutines    │
│  (implicit)        (explicit)          │
│      │                 │               │
│      └─────────────────┘               │
│                │                       │
│       Runtime Scheduler                │
│           (M:N Model)                  │
│                │                       │
│       ┌────────┴────────┐              │
│       │                 │              │
│    OS Thread 1   OS Thread 2...       │
│       │                 │              │
└───────┼─────────────────┼──────────────┘
```

### Key Differences from Threads

| Feature | Threads | Goroutines |
|---------|---------|-----------|
| Memory | ~1-2 MB | ~2 KB |
| Creation Cost | Expensive | Cheap |
| Quantity | Thousands | Millions |
| Scheduling | OS Kernel | Go Runtime |
| Context Switch | Slow | Fast |
| Easy to Create | Hard | Trivial |

---

## Creating Goroutines

### The `go` Keyword

```go
import (
    "fmt"
    "time"
)

func sayHello() {
    fmt.Println("Hello from goroutine!")
}

func main() {
    // Create goroutine
    go sayHello()

    // Main goroutine continues immediately
    fmt.Println("Main continues")

    // But if main exits, goroutine dies!
    time.Sleep(1 * time.Second)  // Give goroutine time to run
}

// Output (might be in any order):
// Main continues
// Hello from goroutine!
```

**Important**: When `main()` returns, ALL goroutines are terminated!

### Goroutines with Arguments

```go
func greet(name string) {
    fmt.Printf("Hello, %s!\n", name)
}

go greet("Alice")
go greet("Bob")
go greet("Charlie")
```

### Anonymous Functions in Goroutines

```go
go func() {
    fmt.Println("Anonymous goroutine")
}()

// With closure
name := "Alice"
go func() {
    fmt.Printf("Hello, %s\n", name)
}()
```

**Gotcha**: Variables captured by reference!

```go
// ✗ Wrong: Loop variable captured by reference
for i := 0; i < 3; i++ {
    go func() {
        fmt.Println(i)  // All print 3!
    }()
}

// ✓ Correct: Pass as argument
for i := 0; i < 3; i++ {
    go func(id int) {
        fmt.Println(id)  // Prints 0, 1, 2
    }(i)
}
```

---

## Basic Synchronization

### WaitGroup

```go
import (
    "fmt"
    "sync"
)

var wg sync.WaitGroup

func worker(id int) {
    defer wg.Done()  // Signal completion
    fmt.Printf("Worker %d doing work\n", id)
}

func main() {
    wg.Add(3)  // Expect 3 goroutines

    for i := 1; i <= 3; i++ {
        go worker(i)
    }

    wg.Wait()  // Block until all Done()
    fmt.Println("All workers done")
}
```

**WaitGroup Pattern**:
1. `Add(n)` - Register n goroutines
2. `Done()` - Mark one goroutine complete
3. `Wait()` - Block until all Done()

### Channels for Signaling

```go
import "fmt"

func main() {
    done := make(chan bool)

    go func() {
        fmt.Println("Working...")
        done <- true  // Signal completion
    }()

    <-done  // Wait for signal
    fmt.Println("Done")
}
```

---

## Avoiding Race Conditions

### Race Detector

```bash
go run -race main.go      # Detect races at runtime
go test -race ./...       # Test with race detection
```

### Shared Memory Problem

```go
var counter = 0

func main() {
    for i := 0; i < 100; i++ {
        go func() {
            counter++  // ✗ RACE CONDITION!
        }()
    }
    // counter might not be 100
}
```

### Solutions

#### Solution 1: Mutex Lock

```go
import "sync"

var (
    counter int
    mu      sync.Mutex
)

func increment() {
    mu.Lock()
    defer mu.Unlock()
    counter++
}

func main() {
    for i := 0; i < 100; i++ {
        go increment()
    }
}
```

#### Solution 2: Channels

```go
func main() {
    counter := make(chan int)

    go func() {
        count := 0
        for increment := range counter {
            if increment == 1 {
                count++
            }
        }
    }()

    for i := 0; i < 100; i++ {
        counter <- 1
    }
    close(counter)
}
```

#### Solution 3: Atomic Operations

```go
import "sync/atomic"

var counter int64

func increment() {
    atomic.AddInt64(&counter, 1)
}
```

---

## Goroutine Lifetime

### Goroutine States

```
Created → Ready → Running ↔ Blocked → Terminated
           ↓                    ↑
           └────────────────────┘
```

### Goroutine Never Runs

```go
// ✗ Might not execute
go func() {
    fmt.Println("I might not run")
}()
// Main exits immediately
```

### Waiting for Goroutines

```go
import (
    "fmt"
    "sync"
    "time"
)

func main() {
    var wg sync.WaitGroup

    wg.Add(1)
    go func() {
        defer wg.Done()
        time.Sleep(1 * time.Second)
        fmt.Println("Done")
    }()

    wg.Wait()  // Block until goroutine finishes
}
```

---

## Error Handling in Goroutines

### Panics Don't Affect Other Goroutines

```go
go func() {
    panic("Error in goroutine")
}()

go func() {
    fmt.Println("I still run!")
}()

time.Sleep(100 * time.Millisecond)
```

**Output**:
```
I still run!
panic: Error in goroutine
```

### Recovering from Panics

```go
func safeGo(fn func()) {
    go func() {
        defer func() {
            if r := recover(); r != nil {
                fmt.Printf("Recovered: %v\n", r)
            }
        }()
        fn()
    }()
}

safeGo(func() {
    panic("Handled!")
})
```

### Error Channels

```go
func worker(id int, errors chan<- error) {
    defer wg.Done()

    if id == 2 {
        errors <- fmt.Errorf("worker %d failed", id)
        return
    }

    errors <- nil
}

func main() {
    errors := make(chan error, 3)
    wg.Add(3)

    for i := 1; i <= 3; i++ {
        go worker(i, errors)
    }

    wg.Wait()
    close(errors)

    for err := range errors {
        if err != nil {
            fmt.Println("Error:", err)
        }
    }
}
```

---

## Goroutine Limits

### Resource Exhaustion

```go
// ✗ Don't do this - creates millions of goroutines
for i := 0; i < 1000000; i++ {
    go someFunction()
}

// Stack overflow or memory exhaustion
```

### Worker Pool Pattern

```go
import (
    "fmt"
    "sync"
)

func worker(jobs <-chan int, results chan<- int) {
    for job := range jobs {
        fmt.Printf("Processing %d\n", job)
        results <- job * 2
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)

    // Create only 3 workers
    for w := 1; w <= 3; w++ {
        go worker(jobs, results)
    }

    // Send 10 jobs
    for j := 1; j <= 10; j++ {
        jobs <- j
    }
    close(jobs)

    // Collect results
    for a := 1; a <= 10; a++ {
        <-results
    }
}
```

---

## Common Goroutine Issues

### Issue 1: Goroutine Leak

```go
// ✗ Bad: Goroutine never exits
func leak() {
    for {
        // Infinite loop, no exit condition
    }
}

// ✓ Good: Goroutine can exit
func worker(done <-chan struct{}) {
    for {
        select {
        case <-done:
            return
        default:
            // Do work
        }
    }
}
```

### Issue 2: Deadlock

```go
// ✗ Bad: Deadlock
ch := make(chan int)  // Unbuffered
ch <- 1  // Block, no one reading!

// ✓ Good: Use buffered channel
ch := make(chan int, 1)
ch <- 1
val := <-ch
```

### Issue 3: Lost Goroutine

```go
// ✗ Bad: Goroutine might not complete
go func() {
    fmt.Println("Hello")
}()
// Main might exit before goroutine runs

// ✓ Good: Wait for completion
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    fmt.Println("Hello")
}()
wg.Wait()
```

---

## Best Practices

### 1. Always Wait for Goroutines

```go
// ✓ Good
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // Work
}()
wg.Wait()
```

### 2. Close Channels from Sender

```go
// ✓ Good
go func(out chan<- int) {
    for i := 0; i < 10; i++ {
        out <- i
    }
    close(out)  // Sender closes
}(out)

// Receive until closed
for val := range out {
    fmt.Println(val)
}
```

### 3. Use Context for Cancellation

```go
import "context"

go func(ctx context.Context) {
    select {
    case <-ctx.Done():
        return
    default:
        // Do work
    }
}(ctx)
```

### 4. Guard with Select and Timeout

```go
import (
    "time"
)

select {
case result := <-resultChan:
    fmt.Println(result)
case <-time.After(5 * time.Second):
    fmt.Println("Timeout")
}
```

### 5. Use sync.WaitGroup for Coordination

```go
var wg sync.WaitGroup
wg.Add(numWorkers)

for i := 0; i < numWorkers; i++ {
    go func() {
        defer wg.Done()
        // Work
    }()
}

wg.Wait()
```

---

## Practical Examples

### Example 1: Concurrent Downloads

```go
import (
    "fmt"
    "sync"
    "time"
)

func download(url string, results chan<- string) {
    time.Sleep(1 * time.Second)  // Simulate download
    results <- fmt.Sprintf("Downloaded %s", url)
}

func main() {
    results := make(chan string, 3)

    go download("url1", results)
    go download("url2", results)
    go download("url3", results)

    for i := 0; i < 3; i++ {
        fmt.Println(<-results)
    }
}
```

### Example 2: Map-Reduce Pattern

```go
func mapReduce() {
    nums := []int{1, 2, 3, 4, 5}
    results := make(chan int, len(nums))

    // Map: Process each number
    for _, n := range nums {
        go func(num int) {
            results <- num * 2
        }(n)
    }

    // Reduce: Collect results
    sum := 0
    for i := 0; i < len(nums); i++ {
        sum += <-results
    }

    fmt.Println("Sum:", sum)  // 30
}
```

### Example 3: Timeout Pattern

```go
import (
    "fmt"
    "time"
)

func slowTask() <-chan string {
    out := make(chan string)
    go func() {
        time.Sleep(2 * time.Second)
        out <- "Done"
    }()
    return out
}

func main() {
    select {
    case result := <-slowTask():
        fmt.Println(result)
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout!")
    }
}
```

---

## Summary

- **Goroutines**: Lightweight, cheap, managed by runtime
- **go keyword**: Creates goroutine (non-blocking)
- **WaitGroup**: Coordinate multiple goroutines
- **Channels**: Communicate between goroutines
- **Race detector**: Find data races at runtime
- **Synchronization**: Mutex, Atomic, Channels
- **Gotchas**: Loop variable closure, main exit, leaks
- **Patterns**: Worker pool, fan-out/fan-in, pipeline
- **Error handling**: Error channels, panics, recovery

Next: Module 3 Chapter 2 - Channels.
