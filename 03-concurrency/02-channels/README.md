# Chapter 2: Channels

## Learning Objectives

1. Understand channel fundamentals (send, receive, close)
2. Distinguish buffered vs unbuffered channels
3. Use directional channels (send-only, receive-only)
4. Master the select statement for multiplexing
5. Understand closing channels and signaling completion
6. Implement common channel patterns
7. Handle timeouts and deadlocks
8. Use range to iterate channels
9. Combine channels with goroutines effectively
10. Debug channel-related issues

---

## Channel Basics

### What Are Channels?

Channels enable safe communication between goroutines:

```
┌──────────────┐     ch        ┌──────────────┐
│ Goroutine 1  │ ──send──→  ── receive ──→ │ Goroutine 2 │
└──────────────┘     ch        └──────────────┘
```

### Creating Channels

```go
// Unbuffered (synchronous)
ch := make(chan int)

// Buffered (asynchronous)
ch := make(chan int, 3)  // Capacity 3

// Close a channel
close(ch)  // Panic if send, safe on receive
```

---

## Unbuffered Channels

### How They Work

```go
ch := make(chan int)  // No buffer

// Send blocks until receiver ready
ch <- 42

// Receive blocks until sender sends
value := <-ch
```

**Key Point**: Both sender and receiver must be ready.

```go
func main() {
    ch := make(chan int)

    // This would deadlock! No goroutine to receive
    // ch <- 42
    // fmt.Println(<-ch)

    // Correct: use goroutine
    go func() {
        ch <- 42
    }()

    fmt.Println(<-ch)  // Prints 42
}
```

### Synchronization with Unbuffered

```go
import (
    "fmt"
    "sync"
)

func main() {
    done := make(chan bool)
    var wg sync.WaitGroup

    wg.Add(1)
    go func() {
        defer wg.Done()
        fmt.Println("Work starting")
        // Do work
        fmt.Println("Work done")
        done <- true  // Signal completion
    }()

    <-done  // Wait for signal
    wg.Wait()
}
```

---

## Buffered Channels

### How They Work

```go
ch := make(chan int, 3)  // Buffer size 3

// Can send 3 times without blocking
ch <- 1
ch <- 2
ch <- 3

// 4th send blocks until room
// ch <- 4  // Would block!

// Receiving frees buffer
fmt.Println(<-ch)  // Now room for 4th send
ch <- 4
```

### When to Use Buffered

```go
// ✓ Good: Buffer for known number of goroutines
results := make(chan int, 10)
for i := 0; i < 10; i++ {
    go func(id int) {
        results <- id * 2
    }(i)
}

// Collect without blocking producers
for i := 0; i < 10; i++ {
    fmt.Println(<-results)
}

// ✗ Avoid: Unbounded buffer
// results := make(chan int)  // No buffer causes blocking
```

---

## Directional Channels

### Send-Only Channels

```go
func sender(ch chan<- int) {
    ch <- 1
    ch <- 2
    ch <- 3
    close(ch)
}

// Compile error:
// value := <-ch  // Can't receive on send-only
```

### Receive-Only Channels

```go
func receiver(ch <-chan int) {
    for value := range ch {
        fmt.Println(value)
    }
}

// Compile error:
// ch <- 42  // Can't send on receive-only
```

### Type Conversion

```go
ch := make(chan int)

sendCh := (chan<- int)(ch)     // Convert to send-only
recvCh := (<-chan int)(ch)     // Convert to receive-only

// Original is still bidirectional
```

---

## The Select Statement

### Basic Select

```go
import "fmt"

func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)

    go func() {
        ch1 <- "one"
    }()

    go func() {
        ch2 <- "two"
    }()

    select {
    case msg1 := <-ch1:
        fmt.Println("Received:", msg1)
    case msg2 := <-ch2:
        fmt.Println("Received:", msg2)
    }
}
```

**Key Point**: Waits for FIRST case that's ready.

### Default Case

```go
select {
case msg := <-ch:
    fmt.Println(msg)
default:
    fmt.Println("No message available")
}
```

This doesn't block - default executes immediately if others can't.

### Timeout Pattern

```go
import "time"

select {
case result := <-resultCh:
    fmt.Println(result)
case <-time.After(5 * time.Second):
    fmt.Println("Timeout!")
}
```

### Infinite Loop with Select

```go
for {
    select {
    case msg := <-ch1:
        fmt.Println("From ch1:", msg)
    case msg := <-ch2:
        fmt.Println("From ch2:", msg)
    case <-done:
        return  // Exit loop
    }
}
```

---

## Closing Channels

### When to Close

```go
// ✓ Good: Sender closes
func producer(out chan<- int) {
    for i := 0; i < 5; i++ {
        out <- i
    }
    close(out)  // Signal completion
}

// ✗ Bad: Receiver closes (can panic)
func consumer(in <-chan int) {
    close(in)  // Panic if sender sends!
}
```

### Detecting Close

```go
ch := make(chan int)

go func() {
    ch <- 1
    ch <- 2
    close(ch)
}()

// Method 1: Check ok
for {
    value, ok := <-ch
    if !ok {
        fmt.Println("Channel closed")
        break
    }
    fmt.Println(value)
}

// Method 2: Use range (cleaner)
for value := range ch {
    fmt.Println(value)
}
// Loop exits automatically when closed
```

---

## Common Channel Patterns

### Pattern 1: Pipeline

```go
// Generate
generate := func() <-chan int {
    out := make(chan int)
    go func() {
        for i := 1; i <= 3; i++ {
            out <- i
        }
        close(out)
    }()
    return out
}

// Transform
square := func(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

// Consume
for n := range square(generate()) {
    fmt.Println(n)  // 1, 4, 9
}
```

### Pattern 2: Fan-Out/Fan-In

```go
// Fan-out: distribute work
func distribute(in <-chan int, workers int) []<-chan int {
    chs := make([]<-chan int, workers)
    for i := 0; i < workers; i++ {
        ch := make(chan int)
        go func(c chan<- int) {
            for n := range in {
                c <- n * 2
            }
            close(c)
        }(ch)
        chs[i] = (<-chan int)(ch)
    }
    return chs
}

// Fan-in: collect results
func merge(chs ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    wg.Add(len(chs))

    for _, ch := range chs {
        go func(c <-chan int) {
            for n := range c {
                out <- n
            }
            wg.Done()
        }(ch)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

### Pattern 3: Quit Channel

```go
func worker(jobs <-chan int, quit <-chan struct{}) {
    for {
        select {
        case job := <-jobs:
            fmt.Println("Processing", job)
        case <-quit:
            fmt.Println("Quitting")
            return
        }
    }
}

func main() {
    jobs := make(chan int)
    quit := make(chan struct{})

    go worker(jobs, quit)

    jobs <- 1
    jobs <- 2

    close(quit)  // Signal stop
}
```

---

## Common Mistakes

### Mistake 1: Sending on Closed Channel

```go
// ✗ Panic!
ch := make(chan int)
close(ch)
ch <- 42  // Panic
```

### Mistake 2: Receiving from Never-Closed Channel

```go
// ✗ Deadlock
ch := make(chan int)
fmt.Println(<-ch)  // Blocks forever, no sender
```

### Mistake 3: Multiple Closes

```go
// ✗ Panic!
ch := make(chan int)
close(ch)
close(ch)  // Panic on second close
```

### Mistake 4: All Goroutines Blocked

```go
// ✗ Deadlock
func main() {
    ch := make(chan int)
    ch <- 42  // Sends but no receiver!
}
```

---

## Best Practices

### 1. Sender Closes Channel

```go
func producer(out chan<- int) {
    defer close(out)
    for i := 1; i <= 5; i++ {
        out <- i
    }
}
```

### 2. Use Receive-Only in Parameters

```go
func consume(in <-chan int) {
    for value := range in {
        // Process
    }
}

func produce(out chan<- int) {
    defer close(out)
    // Send values
}
```

### 3. Use Range for Simple Iteration

```go
// ✓ Good
for value := range ch {
    process(value)
}

// Less common
for {
    value, ok := <-ch
    if !ok {
        break
    }
    process(value)
}
```

### 4. Use Select for Multiple Operations

```go
// ✓ Good
select {
case value := <-ch1:
    // Handle ch1
case value := <-ch2:
    // Handle ch2
case <-timeout:
    // Handle timeout
}
```

### 5. Buffer When Knowing Capacity

```go
// ✓ Good
results := make(chan int, numWorkers)
for i := 0; i < numWorkers; i++ {
    go worker(results)
}

// ✓ Also good - iterate results later
for i := 0; i < numWorkers; i++ {
    <-results
}
```

---

## Practical Examples

### Example 1: Worker Pool with Result Collection

```go
func work(id int) int {
    return id * 10
}

func main() {
    results := make(chan int, 5)

    for i := 1; i <= 5; i++ {
        go func(id int) {
            results <- work(id)
        }(i)
    }

    for i := 0; i < 5; i++ {
        fmt.Println(<-results)
    }
}
```

### Example 2: Multiplexing Multiple Sources

```go
func getData(source string, out chan<- string) {
    time.Sleep(1 * time.Second)
    out <- fmt.Sprintf("Data from %s", source)
}

func main() {
    ch := make(chan string)

    go getData("A", ch)
    go getData("B", ch)
    go getData("C", ch)

    for i := 0; i < 3; i++ {
        fmt.Println(<-ch)
    }
}
```

### Example 3: Rate Limiting

```go
import "time"

func main() {
    ticker := time.NewTicker(100 * time.Millisecond)
    requests := []int{1, 2, 3, 4, 5}

    for _, req := range requests {
        <-ticker.C
        fmt.Println("Processing", req)
    }

    ticker.Stop()
}
```

---

## Summary

- **Channels**: Safe communication between goroutines
- **Unbuffered**: Synchronous, both parties must be ready
- **Buffered**: Asynchronous, with capacity
- **Directional**: Send-only (chan<-) and receive-only (<-chan)
- **Select**: Multiplex operations on multiple channels
- **Range**: Iterate until channel closes
- **Close**: Sender closes, never on receive-only
- **Patterns**: Pipeline, fan-out/fan-in, quit channel
- **Gotchas**: Close on closed, multiple closes, deadlocks
- **Best Practice**: Always close from sender

Next: Module 3 Chapters 3-5 (Advanced Patterns, Context, Graceful Shutdown) - already complete!
