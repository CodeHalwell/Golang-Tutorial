// Package gateway implements the core API gateway components
package gateway

import (
	"errors"
	"sync"
	"time"
)

// RateLimiter implements the token bucket rate limiting algorithm.
//
// The token bucket algorithm works as follows:
// 1. Start with a capacity of N tokens
// 2. Tokens refill at rate R tokens/second
// 3. Each request costs 1 token
// 4. Request is allowed if tokens >= 1, then tokens--
// 5. If tokens < 1, request is denied
//
// This provides smooth rate limiting with burst capacity.
type RateLimiter struct {
	capacity  float64       // Maximum tokens in bucket
	tokens    float64       // Current tokens available
	fillRate  float64       // Tokens per second (refill rate)
	lastFill  time.Time     // When bucket was last refilled
	mu        sync.Mutex
}

// NewRateLimiter creates a new token bucket rate limiter.
//
// Parameters:
//   - capacity: maximum tokens (allows bursts up to this)
//   - tokensPerSecond: refill rate (steady-state throughput)
//
// Example: NewRateLimiter(100, 10) = 100 token burst, 10 req/sec average
func NewRateLimiter(capacity int, tokensPerSecond float64) *RateLimiter {
	return &RateLimiter{
		capacity:  float64(capacity),
		tokens:    float64(capacity), // Start full
		fillRate:  tokensPerSecond,
		lastFill:  time.Now(),
	}
}

// Allow checks if a request for n tokens is allowed.
//
// Returns true if n tokens are available (and removes them).
// Returns false if fewer than n tokens are available.
//
// This method is thread-safe via mutex protection.
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

// TryAllow attempts to get tokens with a timeout.
//
// Useful in async contexts where you might want to wait briefly
// for tokens to become available.
func (rl *RateLimiter) TryAllow(n int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)

	for {
		if rl.Allow(n) {
			return true
		}

		if time.Now().After(deadline) {
			return false
		}

		// Wait a bit before retrying
		time.Sleep(1 * time.Millisecond)
	}
}

// GetTokens returns the current number of available tokens.
// Useful for monitoring and debugging.
func (rl *RateLimiter) GetTokens() float64 {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.refill()
	return rl.tokens
}

// SetRate updates the refill rate.
// Useful for dynamic rate limit adjustment.
func (rl *RateLimiter) SetRate(tokensPerSecond float64) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.fillRate = tokensPerSecond
}

// Reset resets the bucket to full capacity.
func (rl *RateLimiter) Reset() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.tokens = rl.capacity
	rl.lastFill = time.Now()
}

// Private helper methods

// refill adds tokens based on time elapsed since last fill.
// Must be called with lock held.
func (rl *RateLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(rl.lastFill).Seconds()

	// Add tokens for elapsed time
	tokensToAdd := elapsed * rl.fillRate

	rl.tokens = rl.tokens + tokensToAdd

	// Cap at capacity
	if rl.tokens > rl.capacity {
		rl.tokens = rl.capacity
	}

	rl.lastFill = now
}

// UserRateLimiter manages rate limits per user.
// Each user gets their own token bucket.
type UserRateLimiter struct {
	defaultCapacity    int
	defaultRate        float64
	limiters           map[string]*RateLimiter
	mu                 sync.RWMutex
	cleanupInterval    time.Duration
	lastCleanupTime    time.Time
	inactivityDuration time.Duration
}

// NewUserRateLimiter creates a user-aware rate limiter.
//
// Parameters:
//   - capacity: default tokens per user
//   - tokensPerSecond: default refill rate per user
//   - inactivityDuration: remove users inactive for this long
func NewUserRateLimiter(capacity int, tokensPerSecond float64, inactivityDuration time.Duration) *UserRateLimiter {
	return &UserRateLimiter{
		defaultCapacity:    capacity,
		defaultRate:        tokensPerSecond,
		limiters:           make(map[string]*RateLimiter),
		cleanupInterval:    5 * time.Minute,
		inactivityDuration: inactivityDuration,
	}
}

// Allow checks if a user can make a request.
func (url *UserRateLimiter) Allow(userID string, tokens int) bool {
	url.mu.Lock()

	limiter, exists := url.limiters[userID]
	if !exists {
		// Create new limiter for this user
		limiter = NewRateLimiter(url.defaultCapacity, url.defaultRate)
		url.limiters[userID] = limiter
	}

	url.mu.Unlock()

	// Check with the limiter (unlocked for concurrency)
	return limiter.Allow(tokens)
}

// GetUserLimit returns the current token count for a user.
func (url *UserRateLimiter) GetUserLimit(userID string) (float64, error) {
	url.mu.RLock()
	defer url.mu.RUnlock()

	limiter, exists := url.limiters[userID]
	if !exists {
		return 0, errors.New("user not found")
	}

	return limiter.GetTokens(), nil
}

// SetUserLimit sets a custom rate limit for a specific user.
func (url *UserRateLimiter) SetUserLimit(userID string, capacity int, tokensPerSecond float64) {
	url.mu.Lock()
	defer url.mu.Unlock()

	url.limiters[userID] = NewRateLimiter(capacity, tokensPerSecond)
}

// Cleanup removes inactive user limiters.
// Call periodically to prevent memory leaks.
func (url *UserRateLimiter) Cleanup() {
	url.mu.Lock()
	defer url.mu.Unlock()

	now := time.Now()
	if now.Sub(url.lastCleanupTime) < url.cleanupInterval {
		return
	}

	// Remove inactive limiters
	for userID, limiter := range url.limiters {
		if now.Sub(limiter.lastFill) > url.inactivityDuration {
			delete(url.limiters, userID)
		}
	}

	url.lastCleanupTime = now
}

// EndpointRateLimiter manages rate limits per endpoint.
// Different endpoints can have different rate limits.
type EndpointRateLimiter struct {
	limits map[string]*RateLimiter
	mu     sync.RWMutex
}

// NewEndpointRateLimiter creates an endpoint-aware rate limiter.
func NewEndpointRateLimiter() *EndpointRateLimiter {
	return &EndpointRateLimiter{
		limits: make(map[string]*RateLimiter),
	}
}

// RegisterEndpoint registers a rate limit for an endpoint.
//
// Example: RegisterEndpoint("/api/users", 1000, 100)
// = 1000 request burst, 100 requests per second
func (erl *EndpointRateLimiter) RegisterEndpoint(endpoint string, capacity int, tokensPerSecond float64) {
	erl.mu.Lock()
	defer erl.mu.Unlock()

	erl.limits[endpoint] = NewRateLimiter(capacity, tokensPerSecond)
}

// Allow checks if a request to an endpoint is allowed.
func (erl *EndpointRateLimiter) Allow(endpoint string) bool {
	erl.mu.RLock()
	limiter, exists := erl.limits[endpoint]
	erl.mu.RUnlock()

	if !exists {
		// No limit registered for this endpoint, allow
		return true
	}

	return limiter.Allow(1)
}

// UpdateRate dynamically updates the rate for an endpoint.
func (erl *EndpointRateLimiter) UpdateRate(endpoint string, tokensPerSecond float64) error {
	erl.mu.RLock()
	defer erl.mu.RUnlock()

	limiter, exists := erl.limits[endpoint]
	if !exists {
		return errors.New("endpoint not found")
	}

	limiter.SetRate(tokensPerSecond)
	return nil
}

// SlidingWindowRateLimiter is an alternative to token bucket.
// It counts requests in a sliding time window.
type SlidingWindowRateLimiter struct {
	maxRequests int
	window      time.Duration
	requests    []time.Time
	mu          sync.Mutex
}

// NewSlidingWindowRateLimiter creates a sliding window rate limiter.
//
// Parameters:
//   - maxRequests: maximum requests allowed
//   - window: the time window to count requests in
//
// Example: NewSlidingWindowRateLimiter(100, 1*time.Second)
// = 100 requests per second
func NewSlidingWindowRateLimiter(maxRequests int, window time.Duration) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		maxRequests: maxRequests,
		window:      window,
		requests:    make([]time.Time, 0, maxRequests),
	}
}

// Allow checks if a request is allowed under the sliding window.
func (swrl *SlidingWindowRateLimiter) Allow() bool {
	swrl.mu.Lock()
	defer swrl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-swrl.window)

	// Remove old requests outside the window
	newRequests := make([]time.Time, 0)
	for _, req := range swrl.requests {
		if req.After(cutoff) {
			newRequests = append(newRequests, req)
		}
	}
	swrl.requests = newRequests

	// Check if we can accept this request
	if len(swrl.requests) < swrl.maxRequests {
		swrl.requests = append(swrl.requests, now)
		return true
	}

	return false
}
