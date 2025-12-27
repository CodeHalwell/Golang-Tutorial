// Package gateway provides API gateway components including circuit breaker
package gateway

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// CircuitBreakerState represents the state of a circuit breaker
type CircuitBreakerState int

const (
	// StateClosed: Normal operation, requests pass through
	StateClosed CircuitBreakerState = iota
	// StateOpen: Failure detected, requests blocked to prevent cascading failures
	StateOpen
	// StateHalfOpen: Testing recovery, limited requests allowed
	StateHalfOpen
)

func (s CircuitBreakerState) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

// ErrCircuitOpen is returned when a circuit breaker is open
var ErrCircuitOpen = errors.New("circuit breaker is open")

// ErrCircuitBreakerTimeout is returned when trying too many times in HALF_OPEN
var ErrCircuitBreakerTimeout = errors.New("circuit breaker timeout")

// CircuitBreaker implements the circuit breaker pattern.
//
// The circuit breaker prevents cascading failures by:
// 1. Monitoring for failures in a service
// 2. Blocking requests when failure rate exceeds threshold
// 3. Periodically testing if service has recovered
//
// States:
// - CLOSED (normal): Requests pass through, failures monitored
// - OPEN (failing): Requests blocked immediately without calling service
// - HALF_OPEN (recovering): Limited requests allowed to test recovery
type CircuitBreaker struct {
	// Configuration
	failureThreshold      int           // Number of failures before opening
	successThreshold      int           // Number of successes in HALF_OPEN before closing
	timeout               time.Duration // How long to wait in OPEN state
	halfOpenMaxRequests   int           // Max requests allowed in HALF_OPEN

	// State
	state          CircuitBreakerState
	failureCount   int
	successCount   int
	lastFailTime   time.Time
	lastStateChange time.Time

	// Metrics
	metrics *CircuitBreakerMetrics

	// Synchronization
	mu sync.RWMutex
}

// CircuitBreakerMetrics tracks circuit breaker statistics
type CircuitBreakerMetrics struct {
	TotalRequests      int64
	SuccessfulRequests int64
	FailedRequests     int64
	RejectedRequests   int64 // Rejected while OPEN
	StateChanges       int64
	CurrentState       CircuitBreakerState
	LastStateChange    time.Time
	mu                 sync.RWMutex
}

// NewCircuitBreaker creates a new circuit breaker.
//
// Parameters:
//   - failureThreshold: Open after N failures
//   - successThreshold: Close after N successes in HALF_OPEN
//   - timeout: Wait this long in OPEN state before trying HALF_OPEN
//
// Example:
//
//	cb := NewCircuitBreaker(5, 2, 30*time.Second)
//	result, err := cb.Call(func() error {
//	    return callBackendService()
//	})
func NewCircuitBreaker(
	failureThreshold int,
	successThreshold int,
	timeout time.Duration,
) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold:    failureThreshold,
		successThreshold:    successThreshold,
		timeout:             timeout,
		halfOpenMaxRequests: 3,
		state:               StateClosed,
		lastStateChange:     time.Now(),
		metrics: &CircuitBreakerMetrics{
			CurrentState:    StateClosed,
			LastStateChange: time.Now(),
		},
	}
}

// Call executes a function through the circuit breaker.
//
// Returns:
// - nil if function succeeds
// - ErrCircuitOpen if breaker is OPEN
// - Original error if function fails
func (cb *CircuitBreaker) Call(fn func() error) error {
	cb.mu.Lock()

	// Check current state
	switch cb.state {
	case StateClosed:
		// Normal operation
		cb.mu.Unlock()
		return cb.callInClosedState(fn)

	case StateOpen:
		// Check if timeout has passed to transition to HALF_OPEN
		if time.Since(cb.lastFailTime) > cb.timeout {
			cb.setState(StateHalfOpen)
			cb.mu.Unlock()
			return cb.callInHalfOpenState(fn)
		}

		// Still in timeout, reject request
		cb.recordRejected()
		cb.mu.Unlock()
		return ErrCircuitOpen

	case StateHalfOpen:
		cb.mu.Unlock()
		return cb.callInHalfOpenState(fn)
	}

	cb.mu.Unlock()
	return fmt.Errorf("unknown circuit breaker state: %v", cb.state)
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreaker) State() CircuitBreakerState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Reset resets the circuit breaker to CLOSED state
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.setState(StateClosed)
	cb.failureCount = 0
	cb.successCount = 0
	cb.lastFailTime = time.Time{}
}

// Metrics returns a snapshot of the circuit breaker metrics
func (cb *CircuitBreaker) Metrics() CircuitBreakerMetrics {
	cb.metrics.mu.RLock()
	defer cb.metrics.mu.RUnlock()
	return *cb.metrics
}

// Private helper methods

// callInClosedState handles execution in CLOSED state
func (cb *CircuitBreaker) callInClosedState(fn func() error) error {
	err := fn()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.recordFailure()

		// Check if we should open the circuit
		if cb.failureCount >= cb.failureThreshold {
			cb.setState(StateOpen)
			return err
		}

		return err
	}

	// Success in CLOSED state - reset failure count
	cb.failureCount = 0
	cb.recordSuccess()
	return nil
}

// callInHalfOpenState handles execution in HALF_OPEN state
func (cb *CircuitBreaker) callInHalfOpenState(fn func() error) error {
	// Check if we've exceeded max requests in HALF_OPEN
	cb.mu.Lock()
	if cb.successCount >= cb.halfOpenMaxRequests {
		cb.mu.Unlock()
		return ErrCircuitBreakerTimeout
	}
	cb.mu.Unlock()

	// Try the call
	err := fn()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		// Failed in HALF_OPEN, go back to OPEN
		cb.setState(StateOpen)
		cb.failureCount = 0
		cb.successCount = 0
		cb.recordFailure()
		return err
	}

	// Success in HALF_OPEN
	cb.successCount++
	cb.recordSuccess()

	// Check if we should close the circuit
	if cb.successCount >= cb.successThreshold {
		cb.setState(StateClosed)
		cb.failureCount = 0
		cb.successCount = 0
	}

	return nil
}

// setState changes the circuit breaker state and records metrics
func (cb *CircuitBreaker) setState(newState CircuitBreakerState) {
	if cb.state != newState {
		cb.state = newState
		cb.lastStateChange = time.Now()

		// Record metrics
		cb.metrics.mu.Lock()
		cb.metrics.CurrentState = newState
		cb.metrics.LastStateChange = cb.lastStateChange
		cb.metrics.StateChanges++
		cb.metrics.mu.Unlock()

		// Reset counters for new state
		if newState == StateClosed {
			cb.failureCount = 0
			cb.successCount = 0
		}
	}
}

// recordFailure increments failure count
func (cb *CircuitBreaker) recordFailure() {
	cb.lastFailTime = time.Now()

	cb.metrics.mu.Lock()
	cb.metrics.FailedRequests++
	cb.metrics.TotalRequests++
	cb.metrics.mu.Unlock()
}

// recordSuccess increments success count
func (cb *CircuitBreaker) recordSuccess() {
	cb.metrics.mu.Lock()
	cb.metrics.SuccessfulRequests++
	cb.metrics.TotalRequests++
	cb.metrics.mu.Unlock()
}

// recordRejected increments rejected count
func (cb *CircuitBreaker) recordRejected() {
	cb.metrics.mu.Lock()
	cb.metrics.RejectedRequests++
	cb.metrics.TotalRequests++
	cb.metrics.mu.Unlock()
}

// MultiCircuitBreaker manages multiple circuit breakers, one per endpoint/service
type MultiCircuitBreaker struct {
	breakers map[string]*CircuitBreaker
	config   CircuitBreakerConfig
	mu       sync.RWMutex
}

// CircuitBreakerConfig holds default configuration for circuit breakers
type CircuitBreakerConfig struct {
	FailureThreshold    int
	SuccessThreshold    int
	Timeout             time.Duration
	HalfOpenMaxRequests int
}

// NewMultiCircuitBreaker creates a manager for multiple circuit breakers
func NewMultiCircuitBreaker(config CircuitBreakerConfig) *MultiCircuitBreaker {
	return &MultiCircuitBreaker{
		breakers: make(map[string]*CircuitBreaker),
		config:   config,
	}
}

// GetOrCreate gets an existing circuit breaker or creates a new one for an endpoint
func (mcb *MultiCircuitBreaker) GetOrCreate(endpoint string) *CircuitBreaker {
	mcb.mu.Lock()
	defer mcb.mu.Unlock()

	if cb, exists := mcb.breakers[endpoint]; exists {
		return cb
	}

	// Create new circuit breaker for this endpoint
	cb := NewCircuitBreaker(
		mcb.config.FailureThreshold,
		mcb.config.SuccessThreshold,
		mcb.config.Timeout,
	)
	cb.halfOpenMaxRequests = mcb.config.HalfOpenMaxRequests

	mcb.breakers[endpoint] = cb
	return cb
}

// Call executes a function through the appropriate circuit breaker
func (mcb *MultiCircuitBreaker) Call(endpoint string, fn func() error) error {
	cb := mcb.GetOrCreate(endpoint)
	return cb.Call(fn)
}

// GetMetrics returns metrics for all circuit breakers
func (mcb *MultiCircuitBreaker) GetMetrics() map[string]CircuitBreakerMetrics {
	mcb.mu.RLock()
	defer mcb.mu.RUnlock()

	metrics := make(map[string]CircuitBreakerMetrics)
	for endpoint, cb := range mcb.breakers {
		metrics[endpoint] = cb.Metrics()
	}
	return metrics
}

// ResetAll resets all circuit breakers
func (mcb *MultiCircuitBreaker) ResetAll() {
	mcb.mu.RLock()
	defer mcb.mu.RUnlock()

	for _, cb := range mcb.breakers {
		cb.Reset()
	}
}
