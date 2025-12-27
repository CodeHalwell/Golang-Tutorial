package gateway

import (
	"errors"
	"testing"
	"time"
)

// TestCircuitBreakerClosedState verifies normal operation in CLOSED state
func TestCircuitBreakerClosedState(t *testing.T) {
	cb := NewCircuitBreaker(3, 2, 1*time.Second)

	// Successful calls should pass through
	if err := cb.Call(func() error { return nil }); err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}

	if cb.State() != StateClosed {
		t.Errorf("Expected CLOSED, got %v", cb.State())
	}
}

// TestCircuitBreakerOpensOnThreshold verifies circuit opens after threshold failures
func TestCircuitBreakerOpensOnThreshold(t *testing.T) {
	cb := NewCircuitBreaker(3, 2, 1*time.Second)

	failError := errors.New("service failed")

	// Fail exactly threshold times
	for i := 0; i < 3; i++ {
		if err := cb.Call(func() error { return failError }); err == nil || err == ErrCircuitOpen {
			t.Errorf("Iteration %d: Expected failError, got %v", i, err)
		}
	}

	// Should now be OPEN
	if cb.State() != StateOpen {
		t.Errorf("Expected OPEN after threshold, got %v", cb.State())
	}

	// Further calls should be rejected immediately
	if err := cb.Call(func() error { return nil }); err != ErrCircuitOpen {
		t.Errorf("Expected ErrCircuitOpen, got %v", err)
	}
}

// TestCircuitBreakerHalfOpenTransition verifies transition to HALF_OPEN
func TestCircuitBreakerHalfOpenTransition(t *testing.T) {
	timeout := 100 * time.Millisecond
	cb := NewCircuitBreaker(1, 1, timeout)

	// Open the circuit
	cb.Call(func() error { return errors.New("fail") })

	if cb.State() != StateOpen {
		t.Errorf("Expected OPEN, got %v", cb.State())
	}

	// Immediately should still be OPEN
	if err := cb.Call(func() error { return nil }); err != ErrCircuitOpen {
		t.Errorf("Expected ErrCircuitOpen, got %v", err)
	}

	// Wait for timeout
	time.Sleep(timeout + 50*time.Millisecond)

	// Next call should transition to HALF_OPEN
	err := cb.Call(func() error { return nil })
	if err != nil {
		t.Errorf("Expected success in HALF_OPEN, got error: %v", err)
	}

	if cb.State() != StateHalfOpen {
		t.Errorf("Expected HALF_OPEN after timeout, got %v", cb.State())
	}
}

// TestCircuitBreakerClosesAfterRecovery verifies circuit closes after recovery
func TestCircuitBreakerClosesAfterRecovery(t *testing.T) {
	cb := NewCircuitBreaker(1, 1, 100*time.Millisecond)

	// Open the circuit
	cb.Call(func() error { return errors.New("fail") })

	// Wait for transition window
	time.Sleep(150 * time.Millisecond)

	// Successful call in HALF_OPEN should close it
	if err := cb.Call(func() error { return nil }); err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}

	if cb.State() != StateClosed {
		t.Errorf("Expected CLOSED after recovery, got %v", cb.State())
	}

	// Metrics should show state change
	metrics := cb.Metrics()
	if metrics.StateChanges < 2 {
		t.Errorf("Expected at least 2 state changes, got %d", metrics.StateChanges)
	}
}

// TestCircuitBreakerReopensOnFailureInHalfOpen verifies reopening on failure during recovery
func TestCircuitBreakerReopensOnFailureInHalfOpen(t *testing.T) {
	cb := NewCircuitBreaker(1, 2, 100*time.Millisecond)

	// Open the circuit
	cb.Call(func() error { return errors.New("fail") })

	// Wait for transition
	time.Sleep(150 * time.Millisecond)

	// Fail in HALF_OPEN should reopen
	cb.Call(func() error { return errors.New("fail again") })

	if cb.State() != StateOpen {
		t.Errorf("Expected OPEN after failure in HALF_OPEN, got %v", cb.State())
	}
}

// TestCircuitBreakerMetrics verifies metric collection
func TestCircuitBreakerMetrics(t *testing.T) {
	cb := NewCircuitBreaker(5, 2, 1*time.Second)

	// Make some calls
	cb.Call(func() error { return nil })
	cb.Call(func() error { return nil })
	cb.Call(func() error { return errors.New("fail") })

	metrics := cb.Metrics()

	if metrics.SuccessfulRequests != 2 {
		t.Errorf("Expected 2 successful requests, got %d", metrics.SuccessfulRequests)
	}

	if metrics.FailedRequests != 1 {
		t.Errorf("Expected 1 failed request, got %d", metrics.FailedRequests)
	}

	if metrics.TotalRequests != 3 {
		t.Errorf("Expected 3 total requests, got %d", metrics.TotalRequests)
	}
}

// TestCircuitBreakerReset verifies reset functionality
func TestCircuitBreakerReset(t *testing.T) {
	cb := NewCircuitBreaker(1, 1, 1*time.Second)

	// Open the circuit
	cb.Call(func() error { return errors.New("fail") })

	if cb.State() != StateOpen {
		t.Errorf("Expected OPEN, got %v", cb.State())
	}

	// Reset
	cb.Reset()

	if cb.State() != StateClosed {
		t.Errorf("Expected CLOSED after reset, got %v", cb.State())
	}

	// Should now accept calls
	if err := cb.Call(func() error { return nil }); err != nil {
		t.Errorf("Expected success after reset, got error: %v", err)
	}
}

// TestMultiCircuitBreaker verifies managing multiple endpoints
func TestMultiCircuitBreaker(t *testing.T) {
	config := CircuitBreakerConfig{
		FailureThreshold:    2,
		SuccessThreshold:    1,
		Timeout:             100 * time.Millisecond,
		HalfOpenMaxRequests: 3,
	}
	mcb := NewMultiCircuitBreaker(config)

	// Open circuit for endpoint A
	mcb.Call("endpoint-a", func() error { return errors.New("fail") })
	mcb.Call("endpoint-a", func() error { return errors.New("fail") })

	cbA := mcb.GetOrCreate("endpoint-a")
	if cbA.State() != StateOpen {
		t.Errorf("Expected endpoint-a to be OPEN")
	}

	// Endpoint B should still be CLOSED
	cbB := mcb.GetOrCreate("endpoint-b")
	if cbB.State() != StateClosed {
		t.Errorf("Expected endpoint-b to be CLOSED")
	}

	// Make call on endpoint B
	if err := mcb.Call("endpoint-b", func() error { return nil }); err != nil {
		t.Errorf("Expected success on endpoint-b, got error: %v", err)
	}

	// Get metrics
	metrics := mcb.GetMetrics()
	if len(metrics) < 2 {
		t.Errorf("Expected metrics for at least 2 endpoints, got %d", len(metrics))
	}
}

// TestCircuitBreakerConcurrency verifies thread safety
func TestCircuitBreakerConcurrency(t *testing.T) {
	cb := NewCircuitBreaker(5, 2, 1*time.Second)

	// Run concurrent calls
	const numGoroutines = 10
	const callsPerGoroutine = 10

	done := make(chan struct{}, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < callsPerGoroutine; j++ {
				// Some calls succeed, some fail
				err := errors.New("fail")
				if (id+j)%3 == 0 {
					err = nil
				}

				cb.Call(func() error { return err })
			}
			done <- struct{}{}
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Should complete without panic or race condition
	metrics := cb.Metrics()
	if metrics.TotalRequests == 0 {
		t.Error("Expected some metrics, got 0")
	}
}

// BenchmarkCircuitBreakerClosed benchmarks call throughput in CLOSED state
func BenchmarkCircuitBreakerClosed(b *testing.B) {
	cb := NewCircuitBreaker(100, 2, 1*time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cb.Call(func() error { return nil })
	}
}

// BenchmarkCircuitBreakerOpen benchmarks rejection in OPEN state
func BenchmarkCircuitBreakerOpen(b *testing.B) {
	cb := NewCircuitBreaker(1, 2, 1*time.Hour)

	// Open the circuit
	cb.Call(func() error { return errors.New("fail") })

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cb.Call(func() error { return nil })
	}
}

// BenchmarkMultiCircuitBreaker benchmarks performance with multiple endpoints
func BenchmarkMultiCircuitBreaker(b *testing.B) {
	config := CircuitBreakerConfig{
		FailureThreshold:    5,
		SuccessThreshold:    2,
		Timeout:             1 * time.Second,
		HalfOpenMaxRequests: 3,
	}
	mcb := NewMultiCircuitBreaker(config)

	endpoints := []string{"service-a", "service-b", "service-c"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		endpoint := endpoints[i%len(endpoints)]
		mcb.Call(endpoint, func() error { return nil })
	}
}
