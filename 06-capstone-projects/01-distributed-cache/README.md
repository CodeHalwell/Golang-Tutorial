# Capstone Project 1: Distributed In-Memory Cache with Replication

## Project Overview

Build a **production-grade distributed cache system** similar to Redis but simpler. This project integrates:
- Core concurrency patterns (goroutines, channels, mutexes)
- gRPC for inter-node communication
- Graceful shutdown and lifecycle management
- Health checking and failure detection
- Metrics collection (hit rate, latency)
- Comprehensive testing

**Difficulty:** Intermediate | **Estimated LOC:** 1500-2000 | **Duration:** 2-3 weeks

---

## Learning Objectives

After completing this project, you will understand:

1. **Multi-node architecture** - How distributed systems communicate
2. **Replication protocols** - Maintaining consistency across nodes
3. **Health checking** - Detecting and handling node failures
4. **gRPC communication** - Modern RPC for microservices
5. **Concurrent data structures** - Thread-safe cache implementation
6. **Metrics and observability** - Monitoring system behavior
7. **Production considerations** - Error handling, logging, cleanup

---

## Architecture Overview

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │ gRPC
       ▼
┌─────────────────────────────────────┐
│      Primary Cache Node             │
│  ┌───────────────────────────────┐  │
│  │  In-Memory Store              │  │
│  │  {key: value, ttl}            │  │
│  └───────────────────────────────┘  │
└───────┬───────────────┬──────────────┘
        │               │ Replication (async)
        │               ▼
        │         ┌──────────────┐
        │         │ Replica Node │
        │         │ (In-Memory)  │
        │         └──────────────┘
        ▼
┌──────────────┐
│ Health Check │ (Every 5s)
│ Heartbeat    │
└──────────────┘
```

---

## Core Components

### 1. Cache Store (`cache/store.go`)

```go
type Store struct {
    data sync.Map           // Concurrent key-value store
    mu   sync.RWMutex       // Protects TTL data
    ttls map[string]time.Time
}

type CacheEntry struct {
    Value     []byte
    ExpiresAt time.Time
}
```

**Operations:**
- `Get(key)` - Retrieve value
- `Set(key, value, ttl)` - Store with time-to-live
- `Delete(key)` - Remove key
- `Clear()` - Remove all
- `Len()` - Get size
- `Stats()` - Return metrics

**Key Design Decisions:**
- Use `sync.Map` for lock-free reads
- RWMutex for TTL expiration handling
- Periodic cleanup goroutine for expired keys
- Return copies of values (prevent external mutations)

### 2. Server (`server/server.go`)

Implement gRPC server with:

```go
type CacheServer struct {
    store *cache.Store
    logger Logger
    replicaClient ReplicaClient  // For pushing updates
    metrics Metrics
}

// RPC Methods (defined in .proto)
func (s *CacheServer) Get(ctx context.Context, req *GetRequest) (*GetResponse, error)
func (s *CacheServer) Set(ctx context.Context, req *SetRequest) (*SetResponse, error)
func (s *CacheServer) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error)
```

### 3. Replication (`server/replication.go`)

Asynchronously replicate writes to replica nodes:

```go
type ReplicationQueue struct {
    queue chan *ReplicationEvent
    clients map[string]ReplicaClient
    mu sync.RWMutex
}

type ReplicationEvent struct {
    Operation Operation // SET, DELETE, CLEAR
    Key       string
    Value     []byte
    TTL       time.Duration
    Timestamp time.Time
}
```

**Characteristics:**
- Asynchronous writes (don't block client)
- Batching for efficiency
- Retry logic for failures
- Ordered delivery per key (using channels)

### 4. Health Checker (`server/health.go`)

```go
type HealthChecker struct {
    primary ReplicaClient
    check   *health.Check
    mu      sync.Mutex
}

func (h *HealthChecker) Check(interval time.Duration) {
    ticker := time.NewTicker(interval)
    for range ticker.C {
        // Send heartbeat
        // Mark node up/down
        // Trigger failover if needed
    }
}
```

### 5. Metrics (`metrics/metrics.go`)

Collect and expose:

```go
type Metrics struct {
    hits       int64  // Successful gets
    misses     int64  // Failed gets
    sets       int64  // Set operations
    deletes    int64  // Delete operations
    errors     int64  // Operation errors
    mu         sync.Mutex
}

// Stats returned
type CacheStats struct {
    Hits      int64
    Misses    int64
    HitRate   float64
    Sets      int64
    Deletes   int64
    Errors    int64
    Size      int64
    AvgLatencyMs float64
}
```

---

## Protocol Buffer Definition (`proto/cache.proto`)

```proto
syntax = "proto3";

package cache;

option go_package = "github.com/yourorg/cache/generated";

service Cache {
  rpc Get(GetRequest) returns (GetResponse);
  rpc Set(SetRequest) returns (SetResponse);
  rpc Delete(DeleteRequest) returns (DeleteResponse);
  rpc Stats(StatsRequest) returns (StatsResponse);
  rpc Replicate(ReplicationEvent) returns (ReplicationResponse);
}

message GetRequest {
  string key = 1;
}

message GetResponse {
  bytes value = 1;
  bool found = 2;
  string error = 3;
}

message SetRequest {
  string key = 1;
  bytes value = 2;
  int64 ttl_seconds = 3;
}

message SetResponse {
  bool success = 1;
  string error = 2;
}

message DeleteRequest {
  string key = 1;
}

message DeleteResponse {
  bool success = 1;
}

message StatsRequest {
  // Empty
}

message StatsResponse {
  int64 hits = 1;
  int64 misses = 2;
  int64 sets = 3;
  int64 size = 4;
  double hit_rate = 5;
}

message ReplicationEvent {
  enum Operation {
    SET = 0;
    DELETE = 1;
  }

  Operation operation = 1;
  string key = 2;
  bytes value = 3;
  int64 ttl_seconds = 4;
  int64 timestamp_ms = 5;
}

message ReplicationResponse {
  bool success = 1;
}
```

---

## Implementation Steps

### Step 1: Setup and Project Structure
```bash
mkdir -p 06-capstone-projects/01-distributed-cache/{cache,server,client,proto,generated}

cd 06-capstone-projects/01-distributed-cache
go mod init github.com/yourorg/distributed-cache
go get google.golang.org/grpc
go get google.golang.org/protobuf
```

### Step 2: Implement Core Cache Store

Start with `cache/store.go`:
- Thread-safe storage using `sync.Map`
- TTL management with `time.Time`
- Cleanup goroutine for expired entries

```go
package cache

import (
    "sync"
    "time"
)

type Store struct {
    data sync.Map
    mu   sync.RWMutex
    ttls map[string]time.Time
}

func NewStore() *Store {
    return &Store{
        ttls: make(map[string]time.Time),
    }
}
```

### Step 3: Generate Protocol Buffer Code

```bash
# Compile .proto files to Go code
protoc --go_out=. --go-grpc_out=. proto/cache.proto
```

### Step 4: Implement gRPC Server

Create `server/server.go` with:
- Get/Set/Delete handlers
- Graceful shutdown (using patterns from Module 3)
- Context propagation for deadlines

### Step 5: Add Replication

Implement asynchronous replication to replicas

### Step 6: Implement Health Checking

Periodic heartbeat to detect failures

### Step 7: Add Metrics Collection

Track hits, misses, latency

### Step 8: Create Client Library

Simple client for applications to use

---

## Testing Strategy

### Unit Tests (`*_test.go`)

```go
// Test concurrent access
func TestConcurrentAccess(t *testing.T) {
    store := cache.NewStore()

    // Spawn 100 goroutines writing/reading
    // Verify no race conditions
}

// Test TTL expiration
func TestTTLExpiration(t *testing.T) {
    // Set key with 100ms TTL
    // Wait 150ms
    // Verify key is gone
}

// Test metrics
func TestMetricsCollection(t *testing.T) {
    // Perform operations
    // Verify correct hit/miss rates
}
```

### Integration Tests

```go
// Test multi-node replication
func TestReplicationAcrossNodes(t *testing.T) {
    // Start primary and replica
    // Write to primary
    // Read from replica
    // Verify consistency
}

// Test failover
func TestFailoverOnReplicaDown(t *testing.T) {
    // Start primary + replica
    // Stop replica
    // Verify health checker detects it
    // Verify operations still work
}
```

### Load Tests

```bash
# Use Apache Bench or similar
ab -n 10000 -c 100 http://localhost:8080/cache/key

# Or custom Go benchmark
BenchmarkCacheGet-8  1000000  1234 ns/op
BenchmarkCacheSet-8  500000   2567 ns/op
```

---

## Requirements Checklist

- [ ] **Core Cache**: In-memory storage with TTL
- [ ] **Concurrency**: Thread-safe operations
- [ ] **gRPC API**: Get, Set, Delete, Stats operations
- [ ] **Replication**: Async replication to replicas
- [ ] **Health Checking**: Periodic heartbeat
- [ ] **Metrics**: Hit rate, latency, operation counts
- [ ] **Graceful Shutdown**: Clean lifecycle management
- [ ] **Error Handling**: Proper error messages and recovery
- [ ] **Logging**: Structured logging of operations
- [ ] **Documentation**: Comments on public APIs
- [ ] **Tests**: 80%+ code coverage
- [ ] **Benchmarks**: Performance analysis

---

## Performance Targets

| Operation | Target Latency |
|-----------|---|
| Get (cache hit) | <1ms |
| Set | <5ms |
| Delete | <1ms |
| Replication (async) | <50ms per batch |

---

## Extensions

After completing the basic project, consider:

1. **Persistence**: Add RocksDB for disk persistence
2. **Clustering**: Add automatic node discovery
3. **Eviction Policies**: LRU, LFU when capacity exceeded
4. **Monitoring**: Export Prometheus metrics
5. **Admin Console**: Web UI to view/manage cache
6. **Transactions**: Multi-key atomic operations
7. **Pub/Sub**: Topic subscription for invalidations
8. **Distributed Locking**: Implement distributed mutex

---

## Real-World Patterns Applied

| Pattern | Location | Purpose |
|---------|----------|---------|
| **Graceful Shutdown** | `server/server.go` | Clean lifecycle |
| **WaitGroup** | `server/replication.go` | Track goroutines |
| **Context** | RPC handlers | Deadline propagation |
| **Channels** | `replication_queue.go` | Goroutine communication |
| **RWMutex** | `cache/store.go` | Concurrent reads |
| **Health checking** | `server/health.go` | Failure detection |
| **Metrics** | `metrics/` | Observability |
| **Circuit breaker** | `client/` | Resilience |

---

## Success Criteria

✅ **Excellent (95-100%)**
- All requirements met
- Comprehensive tests with 85%+ coverage
- Clean, well-documented code
- Handles edge cases
- Performance targets achieved
- Graceful degradation on failures

✅ **Good (85-94%)**
- All core requirements met
- Solid test coverage (75%+)
- Clear code and documentation
- Handles most error cases
- Reasonable performance

✅ **Acceptable (75-84%)**
- Basic requirements met
- Moderate test coverage (60%+)
- Code is readable
- Some error handling
- Works on happy path

---

## Next Steps

After completing this project:
1. **Review capstone project 2** (API Gateway) for advanced middleware patterns
2. **Study Module 5** (Internals) to optimize performance
3. **Add monitoring** with Prometheus and Grafana
4. **Deploy** using Docker and Kubernetes

This project is a foundation for understanding distributed systems architecture!
