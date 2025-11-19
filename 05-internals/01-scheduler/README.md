# Chapter 1: Go Scheduler and Concurrency Model

## Learning Objectives

1. Understand the M:N threading model (Goroutines : OS Threads)
2. Understand G/M/P scheduling concepts
3. Learn about work stealing and scheduling decisions
4. Understand GOMAXPROCS and its impact
5. Understand preemption and fairness
6. Use scheduler insights to optimize programs
7. Understand goroutine starvation prevention

---

## The G/M/P Model

Go's scheduler is based on three key entities:

### G - Goroutine

- User-level lightweight thread
- Managed by Go runtime (not OS)
- Very cheap to create (~2KB memory)
- Thousands or millions can exist
- Scheduled independently

### M - Machine (OS Thread)

- Operating system thread
- Expensive (~1MB memory)
- Limited by OS resources (typically 1,000-10,000)
- Actually executes the code
- Managed by Go runtime

### P - Processor (Logical Processor)

- Execution context for M
- Holds run queue of Goroutines
- Number = GOMAXPROCS (default = CPU count)
- Scheduler assigns G→M→P

### The Relationship

```
Goroutines (thousands)    M (dozens-hundreds)         Processors (= CPUs)
        G1                    M1 (running)                 P0
        G2         ──→        M2 (idle)           ──→    P1
        G3                    M3 (blocked on I/O)         P2
        G4                                                 P3
        ...
```

**Key Insight**: One M can be running, others blocked. When M blocks, another M picks up the P.

---

## Scheduling Mechanics

### 1. Local Run Queue

Each P has a local run queue of Goroutines:

```go
type P struct {
    runq [256]G  // Local queue (FIFO)
    ...
}
```

**When a Goroutine is ready**:
1. If same P, append to local queue (fast)
2. If different P, append to global queue (slower)

### 2. Global Run Queue

Shared queue for work from other Ps:

```go
var globalRunq []*G  // Shared, requires lock
```

### 3. Work Stealing

When a P's local queue is empty:

1. Check global queue for work
2. Check other Ps' run queues
3. "Steal" half of another P's queue
4. Block M until work available

**This is why Go can efficiently use many Goroutines**.

### 4. Scheduling Loop

```go
for {
    // Check local queue
    g := p.runq.pop()
    if g == nil {
        // Check global queue
        g = globalRunq.pop()
    }
    if g == nil {
        // Try to steal from other Ps
        g = stealWork()
    }
    if g == nil {
        // No work, block M
        park()
        continue
    }

    // Execute goroutine
    execute(g)
}
```

---

## GOMAXPROCS

Controls the number of Ps (and thus parallelism):

```go
import "runtime"

// Get current GOMAXPROCS
numProcs := runtime.GOMAXPROCS(-1)

// Set GOMAXPROCS to 4
runtime.GOMAXPROCS(4)

// Use all CPUs (default in Go 1.5+)
runtime.GOMAXPROCS(runtime.NumCPU())
```

### CPU-Bound vs I/O-Bound

**CPU-Bound Work**:
```go
// Set to number of CPUs
runtime.GOMAXPROCS(runtime.NumCPU())
```

**I/O-Bound Work**:
```go
// Can use more Ps than CPUs (because many will be blocked on I/O)
runtime.GOMAXPROCS(runtime.NumCPU() * 2)
```

---

## Preemption and Fairness

### Preemption (Go 1.14+)

Goroutines can be preempted:

```go
// Goroutines now get ~10ms time slice
for {
    // After ~10ms, this goroutine can be preempted
    doWork()
}
```

**Before Go 1.14**: Only on function calls (could starve others).

**After Go 1.14**: Can preempt in loops (fair scheduling).

---

## How to Optimize Using Scheduler Knowledge

### 1. Avoid Blocking Ms

**Bad**: Holding locks while doing CPU-bound work
```go
// ✗ Bad: Lock holds M for too long
func processData() {
    mu.Lock()
    for i := 0; i < 1000000; i++ {
        // Long-running work while holding lock
        expensiveComputation(i)
    }
    mu.Unlock()
}
```

**Good**: Do work outside lock
```go
// ✓ Good: Work outside lock
results := make([]Result, 0)
for i := 0; i < 1000000; i++ {
    results = append(results, expensiveComputation(i))
}
mu.Lock()
for _, r := range results {
    saveResult(r)
}
mu.Unlock()
```

### 2. Use Buffered Channels Wisely

**Bad**: Unbuffered channels cause blocking
```go
// ✗ Unbuffered: Sender and receiver must synchronize
work := make(chan Task)
go func() {
    task := <-work  // Blocks until sender ready
}()
work <- task  // Blocks until receiver ready
```

**Good**: Buffered channels reduce blocking
```go
// ✓ Buffered: Sender doesn't block if buffer has space
work := make(chan Task, 100)
go func() {
    task := <-work
}()
work <- task  // Doesn't block (buffer has space)
```

### 3. Parallelize I/O

```go
// Goroutines block on I/O, allowing other work to proceed
for i := 0; i < 100; i++ {
    go func(id int) {
        resp, _ := http.Get(url)  // M blocks on network I/O
        process(resp)
    }(i)
}
```

### 4. Understand Goroutine Starvation

```go
// ✗ Starvation: CPU-bound goroutine might starve others
for {
    // This goroutine never yields, might starve others
    doWork()
}

// ✓ Better: Yield occasionally
for {
    doWork()
    // Preemption point (or explicit yield)
}
```

---

## Debugging and Profiling Scheduler

### 1. Runtime Tracing

```bash
# Generate execution trace
GODEBUG=gctrace=1 ./program

# View with go tool trace
go tool trace trace.out
```

### 2. Pprof CPU Profile

```bash
# Profile goroutine scheduling
go test -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof
```

### 3. Debug with Goroutine Dump

```go
runtime.Stack(buf, true)  // Dump all goroutines
```

### 4. Check Runtime Stats

```go
var m runtime.MemStats
runtime.ReadMemStats(&m)
fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
```

---

## Common Scheduler Behaviors

### Behavior 1: Goroutine May Not Start Immediately

```go
go func() {
    fmt.Println("Hello")
}()
// Goroutine might not have run yet!
time.Sleep(1 * time.Millisecond)  // Give it time
```

### Behavior 2: Order of Execution is Unpredictable

```go
for i := 0; i < 10; i++ {
    go fmt.Println(i)
}
// Goroutines might print in any order
```

### Behavior 3: Goroutines Run Concurrently, not Parallel

```go
// Concurrent (interleaved) on single P
// Parallel only if GOMAXPROCS > 1
for i := 0; i < 10; i++ {
    go fmt.Println(i)
}
```

---

## Practical Examples

### Example 1: Measure Scheduler Impact

```go
start := time.Now()

// Many goroutines doing CPU-bound work
const numGoroutines = 100
var wg sync.WaitGroup
wg.Add(numGoroutines)

for i := 0; i < numGoroutines; i++ {
    go func() {
        defer wg.Done()
        // CPU-intensive work
        sum := 0
        for j := 0; j < 10000000; j++ {
            sum += j
        }
    }()
}

wg.Wait()
fmt.Printf("Time: %v\n", time.Since(start))
```

### Example 2: Optimize with GOMAXPROCS

```go
// Single threaded
runtime.GOMAXPROCS(1)
bench1 := benchmark()

// Multi-threaded
runtime.GOMAXPROCS(runtime.NumCPU())
bench2 := benchmark()

fmt.Printf("Single: %v\nMulti: %v\n", bench1, bench2)
```

### Example 3: Work Distribution

```go
// Create worker pool matching GOMAXPROCS
numWorkers := runtime.GOMAXPROCS(-1)
work := make(chan Task, numWorkers)

for i := 0; i < numWorkers; i++ {
    go worker(work)
}

// Submit work
for task := range taskList {
    work <- task
}
```

---

## Key Takeaways

1. **M:N threading**: Thousands of goroutines on dozens of threads
2. **Work stealing**: Keeps all Ps busy
3. **Preemption**: Fair scheduling (Go 1.14+)
4. **GOMAXPROCS**: Control parallelism level
5. **No guarantee of order**: Don't assume execution order
6. **Blocking is hidden**: Network I/O allows other work
7. **Starvation prevention**: Modern Go handles most cases

---

## Next: Memory Management and Escape Analysis

Understanding how goroutines allocate memory and how to optimize with escape analysis.
