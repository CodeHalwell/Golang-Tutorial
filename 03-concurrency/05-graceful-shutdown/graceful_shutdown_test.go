package shutdown

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

// TestBasicGracefulShutdown verifies that a server can start and shutdown cleanly
func TestBasicGracefulShutdown(t *testing.T) {
	logger := &SimpleLogger{}
	server := NewServer(logger)

	// Start the server
	ctx := context.Background()
	if err := server.Start(ctx); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Give it time to run
	time.Sleep(100 * time.Millisecond)

	// Shutdown should complete without error
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}

	// Verify all goroutines have completed
	if err := server.WaitForCompletion(1 * time.Second); err != nil {
		t.Fatalf("Goroutines did not complete: %v", err)
	}
}

// TestShutdownTimeout verifies that shutdown respects the timeout
func TestShutdownTimeout(t *testing.T) {
	logger := &SimpleLogger{}
	server := NewServer(logger)

	// Start the server
	ctx := context.Background()
	if err := server.Start(ctx); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Create a context that times out immediately
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// Shutdown should timeout
	start := time.Now()
	err := server.Shutdown(shutdownCtx)
	elapsed := time.Since(start)

	// We expect a timeout error
	if err == nil {
		t.Error("Expected shutdown to timeout")
	}

	// The timeout should be close to the deadline
	if elapsed < 500*time.Millisecond {
		t.Logf("Shutdown timed out quickly: %v (this is ok)", elapsed)
	}
}

// TestHandleRequest demonstrates request tracking during shutdown
func TestHandleRequest(t *testing.T) {
	logger := &SimpleLogger{}
	server := NewServer(logger)

	// Start the server
	ctx := context.Background()
	if err := server.Start(ctx); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Spawn multiple concurrent requests
	const numRequests = 10
	var completedCount atomic.Int32

	for i := 0; i < numRequests; i++ {
		go func(id int) {
			requestCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			if err := server.HandleRequest(requestCtx, "req-"+string(rune(id))); err != nil {
				t.Logf("Request failed: %v", err)
			} else {
				completedCount.Add(1)
			}
		}(i)
	}

	// Wait for requests to start
	time.Sleep(50 * time.Millisecond)

	// Check active connections
	activeCount := server.ActiveConnectionCount()
	if activeCount == 0 {
		t.Errorf("Expected active connections, got %d", activeCount)
	}

	// Shutdown should wait for requests to complete
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}

	// Verify all goroutines completed
	if err := server.WaitForCompletion(1 * time.Second); err != nil {
		t.Fatalf("Goroutines did not complete: %v", err)
	}

	// After shutdown, active connections should be zero
	finalCount := server.ActiveConnectionCount()
	if finalCount != 0 {
		t.Errorf("Expected no active connections after shutdown, got %d", finalCount)
	}
}

// TestIsShuttingDown verifies the shutdown signal detection
func TestIsShuttingDown(t *testing.T) {
	logger := &SimpleLogger{}
	server := NewServer(logger)

	// Start the server
	ctx := context.Background()
	if err := server.Start(ctx); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Server should not be shutting down initially
	if server.IsShuttingDown() {
		t.Error("Server reports shutting down before shutdown called")
	}

	// Trigger shutdown
	go func() {
		time.Sleep(100 * time.Millisecond)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)
	}()

	// Wait and check if shutdown is detected
	time.Sleep(150 * time.Millisecond)
	if !server.IsShuttingDown() {
		t.Error("Server does not report shutting down after shutdown called")
	}

	// Wait for completion
	server.WaitForCompletion(1 * time.Second)
}

// TestServiceManager demonstrates coordinating multiple services
func TestServiceManager(t *testing.T) {
	logger := &SimpleLogger{}
	manager := NewServiceManager(logger)

	// Create and register multiple services
	service1 := NewServer(logger)
	service2 := NewServer(logger)
	service3 := NewServer(logger)

	manager.RegisterService("service1", service1)
	manager.RegisterService("service2", service2)
	manager.RegisterService("service3", service3)

	// Verify duplicate registration fails
	if err := manager.RegisterService("service1", service1); err == nil {
		t.Error("Expected duplicate registration to fail")
	}

	// Start all services
	ctx := context.Background()
	if err := manager.StartAll(ctx); err != nil {
		t.Fatalf("Failed to start all services: %v", err)
	}

	// Give services time to run
	time.Sleep(100 * time.Millisecond)

	// Shutdown all services
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := manager.ShutdownAll(shutdownCtx); err != nil {
		t.Fatalf("Failed to shutdown all services: %v", err)
	}

	// Verify all services completed
	for _, service := range []*Server{service1, service2, service3} {
		if err := service.WaitForCompletion(1 * time.Second); err != nil {
			t.Fatalf("Service did not complete: %v", err)
		}
	}
}

// BenchmarkHandleRequest measures the overhead of request tracking
func BenchmarkHandleRequest(b *testing.B) {
	logger := &SimpleLogger{}
	server := NewServer(logger)

	// Start the server
	ctx := context.Background()
	server.Start(ctx)

	// Benchmark request handling
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func() {
			requestCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			server.HandleRequest(requestCtx, "bench-req")
		}()
	}

	// Cleanup
	b.StopTimer()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	server.Shutdown(shutdownCtx)
}

// BenchmarkShutdown measures the cost of shutdown coordination
func BenchmarkShutdown(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		logger := &SimpleLogger{}
		server := NewServer(logger)
		ctx := context.Background()
		server.Start(ctx)
		time.Sleep(10 * time.Millisecond) // Let it run a bit

		b.StartTimer()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		server.Shutdown(shutdownCtx)
		cancel()
		server.WaitForCompletion(1 * time.Second)
	}
}

// TestConcurrentRequests ensures no race conditions under load
func TestConcurrentRequests(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrency test in short mode")
	}

	logger := &SimpleLogger{}
	server := NewServer(logger)

	ctx := context.Background()
	if err := server.Start(ctx); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Spawn 100 concurrent requests
	const numRequests = 100
	done := make(chan struct{}, numRequests)

	for i := 0; i < numRequests; i++ {
		go func(id int) {
			requestCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			server.HandleRequest(requestCtx, "req-concurrent")
			done <- struct{}{}
		}(i)
	}

	// Wait for half the requests to complete, then shutdown
	for i := 0; i < numRequests/2; i++ {
		<-done
	}

	// Shutdown should wait for remaining requests
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}

	// Verify clean shutdown
	if err := server.WaitForCompletion(1 * time.Second); err != nil {
		t.Fatalf("Shutdown was not clean: %v", err)
	}
}

// TestContextCancellation ensures context cancellation propagates correctly
func TestContextCancellation(t *testing.T) {
	logger := &SimpleLogger{}
	server := NewServer(logger)

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	if err := server.Start(ctx); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// Cancel the context
	cancel()

	// Server should complete after context cancellation
	if err := server.WaitForCompletion(2 * time.Second); err != nil {
		t.Fatalf("Server did not complete after context cancellation: %v", err)
	}
}
