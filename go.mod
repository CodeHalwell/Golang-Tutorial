module github.com/yourusername/go-mastery-zero-to-one-percent

go 1.21

require (
	// Concurrency and async patterns
	golang.org/x/sync v0.5.0

	// gRPC and Protocol Buffers (for Module 4: Architecture)
	google.golang.org/grpc v1.59.0
	google.golang.org/protobuf v1.31.0

	// Dependency Injection
	github.com/google/wire v0.5.0

	// Linting and code quality
	github.com/golangci/golangci-lint v1.54.2

	// Structured logging (example for production patterns)
	go.uber.org/zap v1.26.0

	// Testing utilities
	github.com/stretchr/testify v1.8.4

	// HTTP and REST patterns
	github.com/gorilla/mux v1.8.1

	// Distributed tracing (OpenTelemetry for Module 2: Architecture)
	go.opentelemetry.io/otel v1.19.0
	go.opentelemetry.io/otel/exporters/jaeger/otlpgrpc v1.19.0
	go.opentelemetry.io/otel/sdk v1.19.0

	// Metrics and observability
	github.com/prometheus/client_golang v1.17.0

	// Configuration management
	github.com/spf13/viper v1.17.0

	// CLI utilities
	github.com/spf13/cobra v1.7.0

	// Rate limiting utilities (for capstone projects)
	golang.org/x/time v0.5.0

	// UUID generation
	github.com/google/uuid v1.5.0

	// Test mocking
	github.com/golang/mock v1.6.0

	// Raft consensus (for capstone project 3)
	github.com/etcd-io/etcd/client/v3 v3.5.10
)

require (
	// Indirect dependencies - included for transparency
	github.com/golang/protobuf v1.5.3 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231211222908-948df04f6b5c // indirect
)

// Development dependencies (optional, can be installed separately)
// For code generation, profiling, and advanced tooling
