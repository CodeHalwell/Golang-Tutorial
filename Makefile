.PHONY: help fmt lint vet test test-race test-coverage build clean install-tools run-example

# Default target
help:
	@echo "Go Mastery Repository - Available targets:"
	@echo ""
	@echo "  make fmt              - Format all Go code"
	@echo "  make lint             - Run linters on all code"
	@echo "  make vet              - Run go vet"
	@echo "  make test             - Run all tests"
	@echo "  make test-race        - Run all tests with race detector"
	@echo "  make test-coverage    - Run tests with coverage report"
	@echo "  make build            - Build all examples"
	@echo "  make clean            - Remove build artifacts"
	@echo "  make install-tools    - Install required development tools"
	@echo "  make run-example      - Run graceful shutdown example (EXAMPLE=path/to/example)"
	@echo ""

# Format code according to Go standards
fmt:
	@echo "Formatting Go code..."
	gofmt -s -w .
	@echo "✓ Formatting complete"

# Run linters
lint:
	@echo "Running linters..."
	golangci-lint run ./...
	@echo "✓ Linting complete"

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...
	@echo "✓ Vet complete"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...
	@echo "✓ Tests complete"

# Run tests with race detector
test-race:
	@echo "Running tests with race detector..."
	go test -race -v ./...
	@echo "✓ Race detector tests complete"

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage analysis..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

# Build all examples
build:
	@echo "Building examples..."
	go build -o /tmp/go-mastery-examples ./...
	@echo "✓ Build complete"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	go clean ./...
	rm -f coverage.out coverage.html
	@echo "✓ Clean complete"

# Install required development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "✓ Tools installed"

# Run specific example
run-example:
	@if [ -z "$(EXAMPLE)" ]; then \
		echo "Usage: make run-example EXAMPLE=03-concurrency/05-graceful-shutdown"; \
		exit 1; \
	fi
	@echo "Running example: $(EXAMPLE)"
	go run ./$(EXAMPLE)/graceful_shutdown.go

# Development workflow shortcuts
dev-setup: install-tools
	@echo "Development environment ready!"

# Pre-commit checks (run before committing)
pre-commit: fmt lint vet test
	@echo "✓ All pre-commit checks passed"

# Module management
mod-download:
	go mod download
	@echo "✓ Dependencies downloaded"

mod-tidy:
	go mod tidy
	@echo "✓ go.mod tidied"

mod-verify:
	go mod verify
	@echo "✓ go.mod verified"

# Profiling targets
profile-cpu:
	@echo "Running CPU profile on graceful_shutdown example..."
	go run -cpuprofile=cpu.prof ./03-concurrency/05-graceful-shutdown/graceful_shutdown.go
	go tool pprof cpu.prof
	@echo "Note: Type 'web' in pprof to generate visualization"

profile-mem:
	@echo "Running memory profile on graceful_shutdown example..."
	go run -memprofile=mem.prof ./03-concurrency/05-graceful-shutdown/graceful_shutdown.go
	go tool pprof mem.prof

# Escape analysis for optimization
escape-analysis:
	@echo "Analyzing escape analysis..."
	go build -gcflags="-m=2" ./05-internals/04-escape-analysis 2>&1 | head -20

# Benchmark targets
benchmark-graceful-shutdown:
	@echo "Benchmarking graceful shutdown..."
	go test -bench=BenchmarkShutdown -benchmem ./03-concurrency/05-graceful-shutdown

benchmark-concurrency:
	@echo "Benchmarking concurrency patterns..."
	go test -bench=. -benchmem ./03-concurrency/...

benchmark-all:
	@echo "Running all benchmarks..."
	go test -bench=. -benchmem ./...

# Documentation
docs:
	@echo "Available documentation:"
	@echo "  - README.md: Main curriculum guide"
	@echo "  - 03-concurrency/05-graceful-shutdown/README.md: Graceful shutdown pattern"
	@find . -name "README.md" -type f | grep -v ".git" | sort

# Information targets
info-go-version:
	@echo "Go Version:"
	@go version

info-env:
	@echo "Environment Information:"
	@echo "GOROOT: $$(go env GOROOT)"
	@echo "GOPATH: $$(go env GOPATH)"
	@echo "GOMAXPROCS: $$(go env GOMAXPROCS)"

# Testing utilities
test-verbose:
	@echo "Running tests with verbose output..."
	go test -v -cover ./...

test-short:
	@echo "Running short tests only..."
	go test -short ./...

test-single:
	@if [ -z "$(TEST)" ]; then \
		echo "Usage: make test-single TEST=TestBasicGracefulShutdown"; \
		exit 1; \
	fi
	go test -v -run $(TEST) ./...

# Continuous integration helpers
ci-check: mod-verify lint vet test-race
	@echo "✓ CI checks passed"

# Generate code (if using Wire or protobuf)
generate:
	@echo "Running go generate..."
	go generate ./...
	@echo "✓ Code generation complete"

# Update dependencies
update-deps:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy
	@echo "✓ Dependencies updated"

# Create exercises directory structure (can be extended per module)
exercises:
	@mkdir -p exercises/{01-basics,02-intermediate,03-concurrency,04-architecture,05-internals,06-capstone}
	@echo "✓ Exercise directories created"

.PHONY: help fmt lint vet test test-race test-coverage build clean install-tools run-example \
        dev-setup pre-commit mod-download mod-tidy mod-verify profile-cpu profile-mem \
        escape-analysis benchmark-graceful-shutdown benchmark-concurrency benchmark-all \
        docs info-go-version info-env test-verbose test-short test-single ci-check \
        generate update-deps exercises
