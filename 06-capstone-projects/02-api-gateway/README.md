# Capstone Project 2: Production-Grade API Gateway with Rate Limiting

## Project Overview

Build a **distributed API gateway** that sits between clients and backend services. This gateway implements:
- HTTP request routing and forwarding
- Token bucket rate limiting (per-user, per-endpoint)
- Circuit breaker pattern for backend resilience
- Request/response middleware pipeline
- Distributed tracing (OpenTelemetry)
- Comprehensive logging and metrics
- Graceful shutdown under load

**Difficulty:** Advanced | **Estimated LOC:** 2500-3500 | **Duration:** 3-4 weeks

---

## Learning Objectives

After completing this project, you will understand:

1. **HTTP Server Design** - Building scalable HTTP services
2. **Request Pipeline Patterns** - Middleware chains and context propagation
3. **Rate Limiting Algorithms** - Token bucket, sliding window, leaky bucket
4. **Circuit Breaker Pattern** - Preventing cascading failures
5. **Distributed Tracing** - Understanding request flow across systems
6. **Advanced Concurrency** - Coordinating multiple services
7. **Observability** - Metrics, logging, tracing
8. **Production Reliability** - Error handling, timeouts, graceful degradation

---

## Architecture Overview

```
┌──────────────┐
│ HTTP Clients │
└──────┬───────┘
       │ (HTTP/REST)
       ▼
┌────────────────────────────────────┐
│     API Gateway                    │
│  ┌────────────────────────────────┐│
│  │ Request Pipeline               ││
│  │ ┌──────────────────────────────┐││
│  │ │ 1. Authentication/Auth       │││
│  │ │ 2. Request Logging           │││
│  │ │ 3. Rate Limiting             │││
│  │ │ 4. Circuit Breaker Check     │││
│  │ │ 5. Tracing Start             │││
│  │ └──────────────────────────────┘││
│  └────────────────────────────────┘│
│  ┌────────────────────────────────┐│
│  │ Backend Routing                ││
│  │ GET    /users    → User Service││
│  │ GET    /orders   → Order Svc   ││
│  │ POST   /products → Product Svc ││
│  └────────────────────────────────┘│
│  ┌────────────────────────────────┐│
│  │ Response Pipeline              ││
│  │ ┌──────────────────────────────┐││
│  │ │ 1. Response Logging          │││
│  │ │ 2. Tracing End               │││
│  │ │ 3. Metrics Export            │││
│  │ │ 4. Error Handling            │││
│  │ └──────────────────────────────┘││
│  └────────────────────────────────┘│
└───────────────────────┬──────────────┘
                        │ (HTTP)
         ┌──────────────┼──────────────┐
         ▼              ▼              ▼
    ┌─────────┐   ┌─────────┐   ┌─────────┐
    │ User    │   │ Order   │   │Product  │
    │ Service │   │ Service │   │ Service │
    └─────────┘   └─────────┘   └─────────┘
```

---

## Core Components

### 1. Rate Limiter (`gateway/rate_limiter.go`)

**Token Bucket Algorithm:**

```go
type RateLimiter struct {
    capacity  float64       // Max tokens
    tokens    float64       // Current tokens
    fillRate  float64       // Tokens per second
    lastFill  time.Time     // Last refill time
    mu        sync.Mutex
}

// Allow(n) checks if n tokens are available
// Returns true if request can proceed, false otherwise
func (rl *RateLimiter) Allow(n int) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    rl.refill()
    if rl.tokens >= float64(n) {
        rl.tokens -= float64(n)
        return true
    }
    return false
}
```

**Per-User Rate Limiting:**
```go
type UserRateLimiter struct {
    limiters map[string]*RateLimiter  // One limiter per user
    mu       sync.RWMutex
}

func (url *UserRateLimiter) Allow(userID string, tokens int) bool {
    // Create limiter if doesn't exist
    // Check if user has quota
}
```

### 2. Circuit Breaker (`gateway/circuit_breaker.go`)

**States:**
- **CLOSED**: Requests allowed, monitoring for errors
- **OPEN**: Requests blocked, preventing cascading failures
- **HALF_OPEN**: Limited requests allowed to test recovery

```go
type CircuitBreaker struct {
    state       State  // CLOSED, OPEN, HALF_OPEN
    failures    int
    threshold   int      // Open after N failures
    timeout     time.Duration
    lastFail    time.Time
    successCount int
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    switch cb.state {
    case CLOSED:
        // Normal operation
        if err := fn(); err != nil {
            cb.recordFailure()
            return err
        }
        return nil

    case OPEN:
        // Check if timeout has passed
        if time.Since(cb.lastFail) > cb.timeout {
            cb.setState(HALF_OPEN)
            return cb.Call(fn)
        }
        return ErrCircuitOpen

    case HALF_OPEN:
        // Test recovery
        if err := fn(); err != nil {
            cb.setState(OPEN)
            return err
        }
        cb.successCount++
        if cb.successCount >= 3 {
            cb.reset()
        }
        return nil
    }
    return nil
}
```

### 3. Request Router (`gateway/router.go`)

```go
type Router struct {
    routes map[string]Route  // Pattern → Backend service
    mu     sync.RWMutex
}

type Route struct {
    Pattern     string  // "/users/:id"
    Backend     string  // "http://user-service:8080"
    RateLimit   int     // Requests per second
    Timeout     time.Duration
    CircuitBreaker *CircuitBreaker
}

func (r *Router) Route(method, path string) (*Route, map[string]string, error) {
    // Match request to route pattern
    // Extract parameters (e.g., :id)
    // Return route and params
}
```

### 4. Middleware Pipeline (`gateway/middleware.go`)

```go
type Handler func(context.Context, *http.Request) (*http.Response, error)

type Middleware func(Handler) Handler

func AuthMiddleware(next Handler) Handler {
    return func(ctx context.Context, req *http.Request) (*http.Response, error) {
        // Check authentication
        token := req.Header.Get("Authorization")
        if token == "" {
            return ErrorResponse(401, "Unauthorized"), nil
        }
        return next(ctx, req)
    }
}

func LoggingMiddleware(next Handler) Handler {
    return func(ctx context.Context, req *http.Request) (*http.Response, error) {
        start := time.Now()
        resp, err := next(ctx, req)
        duration := time.Since(start)
        log.Printf("%s %s took %v", req.Method, req.URL.Path, duration)
        return resp, err
    }
}
```

### 5. Distributed Tracing (`gateway/tracing.go`)

```go
type Tracer struct {
    exporter trace.SpanExporter
}

func (t *Tracer) StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
    return t.tracer.Start(ctx, name)
}

// Usage in middleware:
ctx, span := tracer.StartSpan(ctx, "http.request")
defer span.End()

// Automatic context propagation to downstream services
```

### 6. Metrics Collection (`gateway/metrics.go`)

```go
type Metrics struct {
    RequestsTotal      prometheus.Counter
    RequestDuration    prometheus.Histogram
    BackendErrors      prometheus.Counter
    RateLimitExceeded  prometheus.Counter
    CircuitBreakerOpen prometheus.Gauge
}
```

---

## Implementation Steps

### Step 1: Project Setup

```bash
mkdir -p 06-capstone-projects/02-api-gateway/{gateway,middleware,backends,cmd}
cd 06-capstone-projects/02-api-gateway
go mod init github.com/yourorg/api-gateway

go get github.com/gorilla/mux
go get go.opentelemetry.io/otel
go get github.com/prometheus/client_golang
```

### Step 2: Implement Rate Limiter

Create `gateway/rate_limiter.go`:
- Token bucket algorithm
- Per-user rate limiting
- Endpoint-specific limits

```go
// Example usage:
limiter := NewRateLimiter(10, 1.0)  // 10 tokens, refill at 1/sec
if limiter.Allow(1) {
    // Request allowed
}
```

### Step 3: Implement Circuit Breaker

Create `gateway/circuit_breaker.go`:
- Three states: CLOSED, OPEN, HALF_OPEN
- Automatic state transitions
- Failure tracking

### Step 4: Setup Request Router

Create `gateway/router.go`:
- Route pattern matching
- Parameter extraction
- Route-specific configuration

### Step 5: Build Middleware Pipeline

Create `gateway/middleware.go`:
- Authentication/Authorization
- Request logging
- Rate limit checking
- Tracing integration

### Step 6: Implement HTTP Server

Create `gateway/server.go`:
- HTTP request handling
- Middleware chain execution
- Backend request forwarding

### Step 7: Add Distributed Tracing

Integrate OpenTelemetry:
```bash
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/exporters/jaeger
```

### Step 8: Metrics and Monitoring

Integrate Prometheus metrics

### Step 9: Comprehensive Testing

- Unit tests for rate limiter, circuit breaker
- Integration tests with mock backends
- Load tests with concurrent requests

---

## Testing Strategy

### Unit Tests

```go
// Test rate limiter
func TestRateLimiter(t *testing.T) {
    rl := NewRateLimiter(10, 1.0)

    // Should allow up to capacity
    for i := 0; i < 10; i++ {
        if !rl.Allow(1) {
            t.Error("Should allow within capacity")
        }
    }

    // Should block when exhausted
    if rl.Allow(1) {
        t.Error("Should block when exhausted")
    }
}

// Test circuit breaker
func TestCircuitBreaker(t *testing.T) {
    cb := NewCircuitBreaker(3, 1*time.Second)

    // Simulate 3 failures
    for i := 0; i < 3; i++ {
        cb.Call(func() error { return errors.New("failed") })
    }

    // Should be OPEN now
    if cb.State() != OPEN {
        t.Error("Circuit should be OPEN after threshold failures")
    }
}
```

### Integration Tests

```go
// Test with mock backend
func TestGatewayWithMockBackend(t *testing.T) {
    // Start mock backend server
    mockServer := httptest.NewServer(
        http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(200)
            w.Write([]byte("OK"))
        }),
    )
    defer mockServer.Close()

    // Create gateway with route to mock
    gateway := NewGateway()
    gateway.RegisterRoute("/api/:service", mockServer.URL)

    // Test request
    req, _ := http.NewRequest("GET", "/api/test", nil)
    rec := httptest.NewRecorder()
    gateway.ServeHTTP(rec, req)

    if rec.Code != 200 {
        t.Errorf("Expected 200, got %d", rec.Code)
    }
}
```

### Load Tests

```bash
# Apache Bench
ab -n 100000 -c 1000 http://localhost:8080/api/users

# Custom Go load test
go test -bench=BenchmarkGateway -benchtime=10s
```

---

## Requirements Checklist

- [ ] **HTTP Server**: Accept and route requests
- [ ] **Rate Limiting**: Per-user, per-endpoint limits
- [ ] **Circuit Breaker**: Prevent cascading failures
- [ ] **Request Routing**: Dynamic route matching
- [ ] **Middleware Pipeline**: Composable request processing
- [ ] **Distributed Tracing**: OpenTelemetry integration
- [ ] **Metrics**: Prometheus metrics collection
- [ ] **Error Handling**: Proper error responses
- [ ] **Graceful Shutdown**: Clean lifecycle management
- [ ] **Logging**: Structured logging (zap or similar)
- [ ] **Tests**: 80%+ code coverage
- [ ] **Documentation**: API and architecture docs

---

## Performance Targets

| Metric | Target |
|--------|--------|
| Request latency (p99) | <100ms |
| Throughput | >10k req/s |
| Memory per request | <1MB |
| Circuit breaker trigger | <100ms |
| Rate limiter overhead | <1µs |

---

## Extensions

After completing the basic project:

1. **Load Balancing**: Round-robin, least-connections, weighted
2. **Caching**: Response caching with TTL
3. **Request Transformation**: Modify requests/responses
4. **Authentication**: JWT, OAuth2, API keys
5. **API Versioning**: Support multiple API versions
6. **WebSocket Support**: Upgrade to WebSocket connections
7. **GraphQL Gateway**: Support GraphQL queries
8. **Request Validation**: JSON schema validation

---

## Real-World Patterns Applied

| Pattern | Purpose |
|---------|---------|
| **Graceful Shutdown** | Clean resource cleanup |
| **Circuit Breaker** | Failure resilience |
| **Rate Limiting** | Resource protection |
| **Middleware** | Cross-cutting concerns |
| **Context** | Request-scoped data |
| **Distributed Tracing** | End-to-end visibility |
| **Metrics** | Observability |
| **Health Checks** | Service monitoring |

---

## Success Criteria

✅ **Excellent (95-100%)**
- All requirements implemented
- Comprehensive test coverage (85%+)
- Handles edge cases
- Performance targets met
- Production-ready error handling
- Clean, well-documented code

✅ **Good (85-94%)**
- Core requirements met
- Solid test coverage (75%+)
- Good performance
- Clear error handling

✅ **Acceptable (75-84%)**
- Basic requirements met
- Moderate coverage (60%+)
- Works on happy path

---

## Real-World Applications

This API Gateway pattern is used by:
- Kong (API Gateway platform)
- AWS API Gateway
- Google Cloud Endpoints
- Nginx (with modules)
- Traefik (cloud-native gateway)

---

## Next Steps

After completing this project:
1. Deploy using Docker and Kubernetes
2. Add monitoring dashboard (Grafana)
3. Implement auto-scaling
4. Add advanced authentication (OAuth2, OIDC)
5. Implement request/response transformation

This project integrates advanced Go patterns and prepares you for building production microservices!
