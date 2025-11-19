# Go Mastery: Zero to Top 1% Developer

> A comprehensive, production-grade curriculum designed to take you from "Hello World" to understanding Go's runtime internals, building distributed systems, and joining the top 1% of Go developers.

**Author's Note:** This repository is structured like a professional onboarding curriculum at a top-tier tech company. Each module builds deliberately on the previous one, emphasizing not just how to write Go code, but how to write *idiomatic* Go that scales to production systems.

---

## 📋 Table of Contents

1. [Philosophy](#philosophy)
2. [Repository Structure](#repository-structure)
3. [The Curriculum](#the-curriculum)
4. [Capstone Projects](#capstone-projects)
5. [Learning Path](#learning-path)
6. [How to Use This Repo](#how-to-use-this-repo)
7. [Prerequisites](#prerequisites)

---

## Philosophy

This curriculum is built on three principles:

1. **Idiomatic Go**: We follow the [Go Proverbs](https://go-proverbs.github.io/) and the philosophy behind them. Code that works isn't enough; code must be *clear*, *simple*, and *maintainable*.

2. **From Theory to Practice**: Each concept is introduced with theory, progresses through progressively complex examples, and culminates in real-world applications in the capstone projects.

3. **Production-Grade Mindset**: From the first module, we think about error handling, logging, testing, and observability. This isn't a toy codebase; it's built like real systems at Google, Uber, and Shopify.

---

## Repository Structure

```
go-mastery-zero-to-one-percent/
├── README.md (this file)
├── go.mod
├── go.sum
├── Makefile
├── .golangci.yml
│
├── 01-basics/
│   ├── 01-hello-world/
│   │   ├── main.go
│   │   ├── README.md (learning objectives)
│   │   └── exercises/
│   ├── 02-variables-types/
│   │   ├── main.go
│   │   ├── types_test.go
│   │   ├── README.md
│   │   └── exercises/
│   ├── 03-control-flow/
│   │   ├── main.go
│   │   ├── control_test.go
│   │   ├── README.md
│   │   └── exercises/
│   ├── 04-functions/
│   │   ├── main.go
│   │   ├── functions_test.go
│   │   ├── README.md
│   │   └── exercises/
│   └── 05-structs-interfaces/
│       ├── main.go
│       ├── structs_test.go
│       ├── README.md
│       └── exercises/
│
├── 02-intermediate/
│   ├── 01-error-handling/
│   │   ├── main.go
│   │   ├── errors_test.go
│   │   ├── README.md (sentinel errors, error wrapping, custom errors)
│   │   └── exercises/
│   ├── 02-generics/
│   │   ├── main.go
│   │   ├── generics_test.go
│   │   ├── README.md (Go 1.18+ generics, constraints)
│   │   └── exercises/
│   ├── 03-testing/
│   │   ├── main.go
│   │   ├── testing_test.go
│   │   ├── benchmarks_test.go
│   │   ├── README.md (table-driven tests, fixtures, mocking)
│   │   └── exercises/
│   ├── 04-file-operations/
│   │   ├── main.go
│   │   ├── files_test.go
│   │   ├── README.md
│   │   └── exercises/
│   └── 05-standard-library-deep-dive/
│       ├── main.go
│       ├── stdlib_test.go
│       ├── README.md (strings, bytes, io, fmt)
│       └── exercises/
│
├── 03-concurrency/
│   ├── 01-goroutines-basics/
│   │   ├── main.go
│   │   ├── goroutines_test.go
│   │   ├── README.md (M:N threading, scheduler basics)
│   │   └── exercises/
│   ├── 02-channels/
│   │   ├── main.go
│   │   ├── channels_test.go
│   │   ├── README.md (buffered, unbuffered, direction)
│   │   └── exercises/
│   ├── 03-advanced-patterns/
│   │   ├── worker_pool.go
│   │   ├── fan_out_fan_in.go
│   │   ├── pipelines.go
│   │   ├── patterns_test.go
│   │   ├── README.md (worker pools, fan-out/fan-in, pipelines)
│   │   └── exercises/
│   ├── 04-context-package/
│   │   ├── main.go
│   │   ├── context_test.go
│   │   ├── README.md (cancellation, timeouts, deadlines)
│   │   └── exercises/
│   └── 05-graceful-shutdown/
│       ├── graceful_shutdown.go
│       ├── graceful_shutdown_test.go
│       ├── README.md
│       └── exercises/
│
├── 04-architecture/
│   ├── 01-project-layout/
│   │   ├── main.go
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── pkg/
│   │   ├── README.md (standard Go project structure)
│   │   └── Makefile
│   ├── 02-interfaces-design/
│   │   ├── main.go
│   │   ├── interfaces_test.go
│   │   ├── README.md (interface nil vs typed nil, composition)
│   │   └── exercises/
│   ├── 03-dependency-injection/
│   │   ├── manual_di.go
│   │   ├── wire_di.go (using Wire)
│   │   ├── di_test.go
│   │   ├── README.md
│   │   └── exercises/
│   ├── 04-clean-architecture/
│   │   ├── domain/
│   │   ├── usecase/
│   │   ├── delivery/
│   │   ├── repository/
│   │   ├── main.go
│   │   ├── README.md (hexagonal, clean architecture in Go)
│   │   └── exercises/
│   └── 05-grpc-protobuf/
│       ├── proto/
│       │   └── service.proto
│       ├── generated/
│       ├── server.go
│       ├── client.go
│       ├── grpc_test.go
│       ├── README.md
│       └── Makefile
│
├── 05-internals/
│   ├── 01-scheduler/
│   │   ├── main.go
│   │   ├── README.md (G/M/P model, GOMAXPROCS, runtime tracing)
│   │   └── exercises/
│   ├── 02-memory-management/
│   │   ├── main.go
│   │   ├── memory_test.go
│   │   ├── README.md (stack vs heap, allocation patterns)
│   │   └── exercises/
│   ├── 03-garbage-collector/
│   │   ├── main.go
│   │   ├── gc_test.go
│   │   ├── README.md (GC tuning, GOGC, concurrent marking)
│   │   └── exercises/
│   ├── 04-escape-analysis/
│   │   ├── main.go
│   │   ├── escape_test.go
│   │   ├── README.md (escape analysis, heap allocations)
│   │   └── exercises/
│   └── 05-profiling/
│       ├── main.go
│       ├── profiling_test.go
│       ├── README.md (pprof, trace, benchmarking)
│       └── exercises/
│
└── 06-capstone-projects/
    ├── 01-distributed-cache/
    │   ├── main.go
    │   ├── cache/
    │   ├── server/
    │   ├── client/
    │   ├── go.mod
    │   ├── README.md (objectives and requirements)
    │   └── Makefile
    ├── 02-api-gateway/
    │   ├── main.go
    │   ├── gateway/
    │   ├── middleware/
    │   ├── go.mod
    │   ├── README.md
    │   └── Makefile
    └── 03-distributed-kv-store/
        ├── main.go
        ├── kv/
        ├── raft/
        ├── rpc/
        ├── go.mod
        ├── README.md
        └── Makefile
```

---

## The Curriculum

### Module 1: Basics (Chapters 01-05)

**Learning Objectives:**
- Write your first Go program and understand the compilation model
- Master Go's type system and zero values
- Control program flow with idiomatic Go constructs
- Write functions with proper error handling from day one
- Design structs and interfaces for composable code

**Key Concepts:**
- Package and module system
- Type declarations (int, string, bool, arrays, slices, maps)
- Functions and multiple return values
- Structs, methods, and receivers
- Interface design and composition
- Defer, panic, and recover

**Deliverables:**
- 20+ working examples
- Table-driven test suite for each chapter
- 40+ exercises with solutions

---

### Module 2: Intermediate (Chapters 01-05)

**Learning Objectives:**
- Master error handling beyond "if err != nil"
- Leverage Go 1.18+ generics effectively
- Write production-grade tests with comprehensive coverage
- Manipulate files and data efficiently
- Understand and use stdlib effectively

**Key Concepts:**
- Sentinel errors vs custom error types vs error wrapping
- Generics constraints and type parameters
- Table-driven tests, mocking patterns, and fixtures
- File I/O, buffering, and streaming
- strings, bytes, io, fmt packages mastery
- Reflection for advanced use cases

**Deliverables:**
- Real error handling patterns used in production systems
- Generic data structures (Stack, Queue, Tree)
- 100+ test cases demonstrating best practices
- File processing utilities with streaming support

---

### Module 3: Concurrency (Chapters 01-05)

**Learning Objectives:**
- Understand Go's M:N threading model and scheduler
- Write safe concurrent code with channels
- Implement advanced patterns (worker pools, pipelines, fan-out/fan-in)
- Master the context package for cancellation and timeouts
- Implement graceful shutdown in real applications
- Use sync/atomic for lock-free programming

**Key Concepts - Advanced Level:**
- Goroutines and their lifecycle
- Buffered and unbuffered channels
- Channel direction and select statements
- Worker pool pattern for resource management
- Fan-out/fan-in for parallel work distribution
- Pipelines for data processing chains
- Context for cancellation, deadlines, and values
- Graceful shutdown with WaitGroups and Context
- sync/atomic for atomic operations
- Memory ordering and visibility guarantees
- False sharing in CPU caches (struct field layout optimization)

**Deliverables:**
- Complete worker pool implementation with metrics
- Multi-stage pipeline processing system
- Graceful shutdown example (provided as style guide)
- Context-based cancellation patterns
- Atomic counter implementations with benchmarks

---

### Module 4: Architecture (Chapters 01-05)

**Learning Objectives:**
- Structure projects for scalability and maintainability
- Design interfaces that enable testing and flexibility
- Implement dependency injection patterns
- Apply clean architecture principles to Go projects
- Build gRPC services with Protocol Buffers

**Key Concepts - Advanced Level:**
- Standard Go project layout (cmd/, internal/, pkg/)
- Interface nil vs typed nil (critical distinction)
- Struct composition vs inheritance
- Dependency injection patterns (manual, Wire library)
- Hexagonal architecture in Go
- Repository pattern for data access
- Use cases for business logic
- gRPC and Protocol Buffer code generation
- Middleware patterns for cross-cutting concerns

**Deliverables:**
- Complete project template following best practices
- DI framework comparison and implementation
- Clean architecture example with all layers
- gRPC service with client and server
- Middleware pipeline for logging, metrics, auth

---

### Module 5: Internals (Chapters 01-05)

**Learning Objectives:**
- Understand the Go runtime at a deep level
- Optimize programs using scheduler insights
- Tune garbage collection for latency-sensitive applications
- Analyze escape analysis to reduce allocations
- Profile and optimize real applications

**Key Concepts - Expert Level:**
- **Scheduler (G/M/P model):**
  - Goroutines (G): execution units
  - Machine threads (M): OS threads
  - Processors (P): execution contexts
  - Run queues and work stealing
  - GOMAXPROCS impact
  - Preemption and fairness

- **Memory Management:**
  - Stack allocation and stack frames
  - Heap allocation and GC roots
  - Pointer semantics vs value semantics
  - Allocation patterns and performance

- **Garbage Collector:**
  - Generational hypothesis in Go
  - Tri-color marking algorithm
  - Concurrent marking and scanning
  - Write barriers
  - GOGC tuning (default 100)
  - GC pause times and throughput

- **Escape Analysis:**
  - Stack vs heap allocations
  - Escape directives
  - Optimization through structure design
  - Cost of allocations in hot paths

- **Profiling & Observability:**
  - pprof: CPU, memory, goroutine, block profiling
  - Execution tracer (go tool trace)
  - Benchmark-driven optimization
  - Metric collection and analysis

**Deliverables:**
- Scheduler behavior demonstrations with tracing
- GC tuning guide with GOGC examples
- Escape analysis examples with assembly output
- Comprehensive profiling guide with tools
- Performance optimization checklist

---

### Module 6: Capstone Projects

Three projects of increasing complexity that integrate all previous concepts.

#### **Project 1: Distributed Cache (Beginner Capstone)**

**Objective:** Build a multi-node in-memory cache with replication

**Requirements:**
- [ ] In-memory key-value store with TTL support
- [ ] Multi-node communication (gRPC)
- [ ] Replication protocol
- [ ] Health checking and failure detection
- [ ] Metrics (hit rate, latency)
- [ ] Graceful shutdown
- [ ] Tests with 80%+ coverage

**Concepts Applied:**
- goroutines and channels for async replication
- context for timeouts and cancellation
- gRPC for inter-node communication
- Clean architecture with interfaces
- Comprehensive testing patterns

**Estimated Lines of Code:** 1500-2000

---

#### **Project 2: Rate-Limited API Gateway (Intermediate Capstone)**

**Objective:** Build a production-grade API gateway with multiple rate-limiting strategies

**Requirements:**
- [ ] HTTP server with routing
- [ ] Token bucket rate limiter (per-user, per-endpoint)
- [ ] Circuit breaker pattern for backend resilience
- [ ] Request/response middleware pipeline
- [ ] Distributed tracing integration (OpenTelemetry)
- [ ] Graceful shutdown under load
- [ ] Load testing and benchmarks
- [ ] Comprehensive logging and observability

**Concepts Applied:**
- Advanced concurrency patterns (worker pools, pipelines)
- Context package for deadline propagation
- Interface design for pluggable components
- Metrics collection and monitoring
- Performance profiling and optimization
- Error handling at scale

**Estimated Lines of Code:** 2500-3500

---

#### **Project 3: Distributed Key-Value Store with Raft Consensus (Advanced Capstone)**

**Objective:** Build a production-grade distributed system with consensus

**Requirements:**
- [ ] Raft consensus implementation or using etcd/raft library
- [ ] Persistent state machine
- [ ] Multi-node cluster coordination
- [ ] gRPC-based inter-node communication
- [ ] Client API (Get, Set, Delete)
- [ ] Leader election and failure recovery
- [ ] Snapshot support for efficient recovery
- [ ] Network partition handling
- [ ] Comprehensive testing (unit, integration, chaos)
- [ ] Performance benchmarks
- [ ] Monitoring and observability
- [ ] Production-grade error handling

**Concepts Applied:**
- Deep understanding of distributed systems
- Advanced concurrency patterns
- Context package for sophisticated cancellation
- gRPC and Protocol Buffers
- State machine design
- Failure detection and recovery
- Profiling under distributed load
- Clean architecture at scale

**Estimated Lines of Code:** 4000-6000

---

## Learning Path

### Recommended Timeline

**Total Duration:** 3-6 months (intense study)
- **Module 1 (Basics):** 2-3 weeks
- **Module 2 (Intermediate):** 3-4 weeks
- **Module 3 (Concurrency):** 4-5 weeks
- **Module 4 (Architecture):** 3-4 weeks
- **Module 5 (Internals):** 4-6 weeks
- **Module 6 (Capstone Projects):** 6-10 weeks

### Three Learning Strategies

#### **Strategy A: Linear (Recommended for Beginners)**
Follow modules 1-6 sequentially. Modules build on each other.

#### **Strategy B: Parallel (For Intermediate Developers)**
If you're already familiar with programming:
- Week 1-2: Skim Module 1, focus on Go-specific concepts
- Week 3+: Do Modules 2-4 in parallel with Module 3
- Week 8+: Deep dive Module 5 while starting capstone projects

#### **Strategy C: Fast-Track (For Experienced Systems Programmers)**
- Week 1: Modules 1-2 (skim)
- Week 2-3: Focus on Module 3 (concurrency is key)
- Week 4+: Module 5 (internals) + capstone projects simultaneously

---

## How to Use This Repo

### 1. Setup

```bash
git clone https://github.com/your-org/go-mastery-zero-to-one-percent.git
cd go-mastery-zero-to-one-percent
go mod download
```

### 2. Run Examples

Each chapter has a `main.go`:

```bash
cd 01-basics/01-hello-world
go run main.go
```

### 3. Run Tests

All examples include comprehensive tests:

```bash
cd 02-intermediate/03-testing
go test -v
```

### 4. Run Benchmarks

Performance is critical:

```bash
go test -bench=. -benchmem ./03-concurrency/...
```

### 5. Format and Lint

Use provided tools:

```bash
make fmt
make lint
make vet
```

### 6. Check Escape Analysis

Understand memory allocations:

```bash
go build -gcflags="-m" ./05-internals/04-escape-analysis
```

### 7. Profile Code

Learn profiling from Module 5:

```bash
go run -cpuprofile=cpu.prof ./05-internals/05-profiling/main.go
go tool pprof cpu.prof
```

---

## Prerequisites

### Required Knowledge
- Basic programming fundamentals (variables, loops, functions)
- Understanding of how compilers work (basic level)
- Familiarity with command line

### Software Requirements
- **Go 1.21+** (We use recent features like slices.All)
- **Git** for version control
- **Make** for running build commands
- **Docker** (for later capstone projects)
- **golangci-lint** for linting
- **Protocol Buffer compiler** (for gRPC module)

### Installation

```bash
# macOS
brew install go protobuf golangci-lint

# Ubuntu/Debian
sudo apt-get install golang-go protobuf-compiler
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Generate protobuf code
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

---

## Key Distinctions: What Makes a Top 1% Go Developer

### 1. Interface Nil vs Typed Nil
Understanding why:
```go
var i interface{} = (*MyType)(nil)
if i != nil { // TRUE! interface{} is not nil
    // This will execute
}
```

### 2. False Sharing and CPU Cache Layout
Knowing that field ordering in structs affects CPU cache efficiency:
```go
// Bad: false sharing
type Counter struct {
    count1 int64
    padding [7]int64 // not needed with proper layout
    count2 int64
}

// Good: separate cache lines
type Counter struct {
    count1 int64
    count2 int64
    _      [8]int64 // padding on last field
}
```

### 3. Escape Analysis Impact
Understanding when allocations happen:
```go
// Escapes to heap
func bad() *int {
    x := 10
    return &x // escapes!
}

// Stays on stack
func good() {
    x := 10
    // x never escapes
}
```

### 4. Goroutine Lifecycle and Scheduler
Knowing that goroutines aren't "free" and understanding GOMAXPROCS:
```go
runtime.GOMAXPROCS(runtime.NumCPU())
// Sets P count equal to CPUs for CPU-bound work
```

### 5. Context Propagation
Properly canceling work in distributed systems:
```go
// Cancels ALL downstream work
ctx, cancel := context.WithCancel(context.Background())
defer cancel() // cleanup
```

### 6. Graceful Shutdown Patterns
Coordinating shutdown of multiple goroutines:
- See `03-concurrency/05-graceful-shutdown/` for production pattern

---

## Repository Standards

### Code Style
- Follow [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- Use `gofmt` for formatting
- Run `golangci-lint` before committing
- Comment all exported functions and types

### Testing Standards
- Minimum 80% code coverage
- Table-driven tests for complex logic
- Benchmarks for performance-critical paths
- Integration tests for multi-component scenarios

### Documentation
- Every module has a detailed README.md
- Code examples are documented with inline comments
- Each capstone project has architecture diagram and docs

---

## Troubleshooting

### "Cannot find package"
Ensure you're in the correct module directory and have run `go mod download`

### "Escape analysis seems wrong"
Run with `-m=2` for more detailed output:
```bash
go build -gcflags="-m=2" ./05-internals/04-escape-analysis
```

### "Tests are slow"
Check for goroutine leaks and excessive allocations:
```bash
go test -race -v ./...
```

---

## Contributing

Found an issue or want to improve examples?
- Open an issue with your Go version (`go version`)
- Submit pull requests with improvements
- Share your solutions to exercises

---

## License

MIT License - See LICENSE file

---

## Acknowledgments

This curriculum is inspired by:
- Google's internal Go training
- Uber's Go Style Guide
- Dave Cheney's Go talks
- Rob Pike's Go Proverbs
- Concurrency in Go by Katherine Cox-Buday

---

**Ready to begin? Start with `01-basics/01-hello-world` and progress methodically. The path is designed, the curriculum is comprehensive, and every concept builds toward the goal: becoming a top 1% Go developer.**