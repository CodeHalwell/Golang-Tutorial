// Package shutdown demonstrates production-grade graceful shutdown patterns.
// This is the style guide for the entire repository.
package shutdown

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Server represents a service that needs graceful shutdown coordination.
// This pattern can be applied to HTTP servers, gRPC services, databases, etc.
type Server struct {
	// shutdown coordinates the shutdown process
	shutdown chan struct{}

	// done signals when all goroutines have finished
	done chan struct{}

	// wg tracks all active goroutines
	wg sync.WaitGroup

	// logger is for observability (could be zap, logrus, etc.)
	logger Logger

	// activeConnections tracks in-flight requests
	activeConnections int32
	connMu            sync.Mutex
}

// Logger interface allows dependency injection of logging
type Logger interface {
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// SimpleLogger is a basic logger implementation for examples
type SimpleLogger struct{}

func (l *SimpleLogger) Infof(format string, args ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", args...)
}

func (l *SimpleLogger) Warnf(format string, args ...interface{}) {
	fmt.Printf("[WARN] "+format+"\n", args...)
}

func (l *SimpleLogger) Errorf(format string, args ...interface{}) {
	fmt.Printf("[ERROR] "+format+"\n", args...)
}

// NewServer creates a new server with proper initialization
func NewServer(logger Logger) *Server {
	if logger == nil {
		logger = &SimpleLogger{}
	}

	return &Server{
		shutdown: make(chan struct{}),
		done:     make(chan struct{}),
		logger:   logger,
	}
}

// Start begins the server's background work.
// This demonstrates the fundamental pattern of goroutine lifecycle management.
func (s *Server) Start(ctx context.Context) error {
	s.logger.Infof("Server starting")

	// Start the main work loop
	s.wg.Add(1)
	go s.workLoop(ctx)

	// Start background maintenance tasks
	s.wg.Add(1)
	go s.maintenanceLoop(ctx)

	// Close done channel when all goroutines finish
	go func() {
		s.wg.Wait()
		close(s.done)
		s.logger.Infof("All goroutines have finished")
	}()

	return nil
}

// workLoop simulates the main server work
func (s *Server) workLoop(ctx context.Context) {
	defer s.wg.Done()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-s.shutdown:
			// Shutdown signal received, stop this loop
			s.logger.Infof("workLoop: shutdown signal received")
			return

		case <-ctx.Done():
			// Context cancelled (e.g., deadline exceeded)
			s.logger.Infof("workLoop: context cancelled")
			return

		case <-ticker.C:
			// Regular work
			s.logger.Infof("workLoop: doing work")
		}
	}
}

// maintenanceLoop simulates background maintenance tasks
func (s *Server) maintenanceLoop(ctx context.Context) {
	defer s.wg.Done()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.shutdown:
			s.logger.Infof("maintenanceLoop: shutdown signal received")
			return

		case <-ctx.Done():
			s.logger.Infof("maintenanceLoop: context cancelled")
			return

		case <-ticker.C:
			s.logger.Infof("maintenanceLoop: doing maintenance")
		}
	}
}

// Shutdown initiates graceful shutdown with a timeout.
// This is the production pattern for coordinating multiple goroutines.
//
// The timeout prevents hung goroutines from blocking shutdown indefinitely.
// A typical timeout is 30 seconds to 1 minute depending on workload.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Infof("Shutdown: initiating graceful shutdown")

	// Signal all goroutines to stop
	close(s.shutdown)

	// Create a deadline to wait for graceful completion
	// Use the caller's context to respect cancellation and deadlines
	shutdownDeadline := time.Now().Add(30 * time.Second)
	ctxDeadline, ok := ctx.Deadline()
	if ok && ctxDeadline.Before(shutdownDeadline) {
		shutdownDeadline = ctxDeadline
	}

	shutdownCtx, cancel := context.WithDeadline(
		ctx, // Use caller's context, not Background
		shutdownDeadline,
	)
	defer cancel()

	// Wait for completion with timeout
	select {
	case <-s.done:
		s.logger.Infof("Shutdown: completed gracefully")
		return nil

	case <-shutdownCtx.Done():
		s.logger.Warnf("Shutdown: timeout waiting for graceful completion")
		return fmt.Errorf("graceful shutdown timeout exceeded")
	}
}

// HandleRequest demonstrates how to track active operations during shutdown.
// This pattern ensures we know when all in-flight requests have completed.
func (s *Server) HandleRequest(ctx context.Context, requestID string) error {
	// Add to active connections
	s.connMu.Lock()
	s.activeConnections++
	connCount := s.activeConnections
	s.connMu.Unlock()

	s.logger.Infof("HandleRequest %s: started (active: %d)", requestID, connCount)

	// Add goroutine to wait group
	s.wg.Add(1)
	defer s.wg.Done()

	// Decrement active connections on completion
	defer func() {
		s.connMu.Lock()
		s.activeConnections--
		s.connMu.Unlock()
		s.logger.Infof("HandleRequest %s: completed", requestID)
	}()

	// Listen for shutdown signal
	select {
	case <-s.shutdown:
		return fmt.Errorf("request %s: shutdown received", requestID)

	case <-ctx.Done():
		return ctx.Err()

	case <-time.After(2 * time.Second):
		// Request completed successfully
		s.logger.Infof("HandleRequest %s: completed successfully", requestID)
		return nil
	}
}

// ActiveConnectionCount returns current number of active requests
func (s *Server) ActiveConnectionCount() int32 {
	s.connMu.Lock()
	defer s.connMu.Unlock()
	return s.activeConnections
}

// IsShuttingDown returns true if shutdown has been initiated
func (s *Server) IsShuttingDown() bool {
	select {
	case <-s.shutdown:
		return true
	default:
		return false
	}
}

// WaitForCompletion blocks until all goroutines have finished.
// Useful for testing and ensuring clean shutdown.
func (s *Server) WaitForCompletion(timeout time.Duration) error {
	select {
	case <-s.done:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("timeout waiting for completion after %v", timeout)
	}
}

// AdvancedGracefulShutdownPattern demonstrates coordinating shutdown
// across multiple services with dependencies.
type ServiceManager struct {
	services map[string]*Server
	logger   Logger
	mu       sync.RWMutex
}

// NewServiceManager creates a new service manager
func NewServiceManager(logger Logger) *ServiceManager {
	return &ServiceManager{
		services: make(map[string]*Server),
		logger:   logger,
	}
}

// RegisterService registers a service for lifecycle management
func (sm *ServiceManager) RegisterService(name string, server *Server) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.services[name]; exists {
		return fmt.Errorf("service %s already registered", name)
	}

	sm.services[name] = server
	return nil
}

// StartAll starts all registered services
func (sm *ServiceManager) StartAll(ctx context.Context) error {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sm.logger.Infof("Starting all services")

	for name, service := range sm.services {
		if err := service.Start(ctx); err != nil {
			sm.logger.Errorf("Failed to start service %s: %v", name, err)
			return err
		}
	}

	return nil
}

// ShutdownAll gracefully shuts down all registered services in reverse order.
// This demonstrates dependency-aware shutdown coordination.
func (sm *ServiceManager) ShutdownAll(ctx context.Context) error {
	sm.mu.RLock()
	services := make([]*Server, 0, len(sm.services))
	serviceNames := make([]string, 0, len(sm.services))

	// Collect services in reverse registration order (LIFO)
	for name, service := range sm.services {
		services = append(services, service)
		serviceNames = append(serviceNames, name)
	}
	sm.mu.RUnlock()

	sm.logger.Infof("Shutting down %d services", len(services))

	// Shutdown in reverse order: last registered first
	for i := len(services) - 1; i >= 0; i-- {
		service := services[i]
		name := serviceNames[i]

		sm.logger.Infof("Shutting down service: %s", name)
		if err := service.Shutdown(ctx); err != nil {
			sm.logger.Warnf("Error shutting down service %s: %v", name, err)
		}
	}

	return nil
}

// Key Principles Demonstrated Here:
//
// 1. **WaitGroup for Goroutine Coordination**: Used to track all active goroutines
//    and wait for their completion. Always call Add() before spawning and Done()
//    when complete.
//
// 2. **Channel Signaling for Shutdown**: The shutdown channel broadcasts the
//    shutdown signal to all goroutines simultaneously. Closing a channel is the
//    idiomatic way to signal multiple goroutines.
//
// 3. **Context Propagation**: Context flows through all functions, allowing
//    cancellation to propagate through the call stack. Use context.WithCancel,
//    context.WithTimeout for different scenarios.
//
// 4. **Timeout Protection**: Graceful shutdown has a bounded timeout (30s default).
//    Without this, hung goroutines could block indefinitely.
//
// 5. **Active Operation Tracking**: Counting active operations allows the server
//    to know when all in-flight requests have completed. Use sync.Mutex for
//    synchronization of shared state.
//
// 6. **Error Handling**: Different shutdown modes (immediate vs graceful) can be
//    implemented by checking the shutdown channel.
//
// 7. **Observability**: Logging at each state transition makes debugging easier
//    and helps understand the shutdown sequence in production.
//
// 8. **LIFO Service Shutdown**: Services with dependencies are shut down in
//    reverse order of registration (last in, first out).
