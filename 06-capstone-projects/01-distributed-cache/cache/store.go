// Package cache provides a thread-safe in-memory cache with TTL support.
package cache

import (
	"errors"
	"sync"
	"time"
)

// ErrNotFound is returned when a key doesn't exist in the cache.
var ErrNotFound = errors.New("key not found")

// Entry represents a single cached value with optional expiration.
type Entry struct {
	Value     []byte
	ExpiresAt time.Time
}

// IsExpired returns true if the entry has passed its TTL.
func (e *Entry) IsExpired() bool {
	return !e.ExpiresAt.IsZero() && time.Now().After(e.ExpiresAt)
}

// Store is a thread-safe in-memory cache with TTL support.
// It uses sync.Map for lock-free reads and a mutex for TTL management.
type Store struct {
	// data holds the actual cached values (thread-safe reads)
	data sync.Map // map[string]*Entry

	// mu protects the TTL map and expiration tracking
	mu sync.RWMutex

	// metrics tracks cache operations
	metrics *StoreMetrics

	// cleanupInterval determines how often to clean expired keys
	cleanupInterval time.Duration

	// stopCleanup signals the cleanup goroutine to stop
	stopCleanup chan struct{}

	// done signals when cleanup goroutine has finished
	done chan struct{}
}

// StoreMetrics tracks cache statistics.
type StoreMetrics struct {
	Gets       int64
	Sets       int64
	Deletes    int64
	Hits       int64
	Misses     int64
	Errors     int64
	TotalSize  int64
	AvgLatency time.Duration
	mu         sync.RWMutex
}

// NewStore creates a new cache store with default cleanup interval of 1 minute.
func NewStore() *Store {
	return NewStoreWithCleanupInterval(1 * time.Minute)
}

// NewStoreWithCleanupInterval creates a new cache store with custom cleanup interval.
func NewStoreWithCleanupInterval(interval time.Duration) *Store {
	s := &Store{
		metrics:         &StoreMetrics{},
		cleanupInterval: interval,
		stopCleanup:     make(chan struct{}),
		done:            make(chan struct{}),
	}

	// Start the cleanup goroutine
	go s.cleanupLoop()

	return s
}

// Get retrieves a value from the cache.
// Returns ErrNotFound if key doesn't exist or has expired.
func (s *Store) Get(key string) ([]byte, error) {
	start := time.Now()
	defer func() {
		s.recordLatency(time.Since(start))
	}()

	s.metrics.recordGet()

	// Try to get from the data store
	val, ok := s.data.Load(key)
	if !ok {
		s.metrics.recordMiss()
		return nil, ErrNotFound
	}

	entry := val.(*Entry)

	// Check if expired
	if entry.IsExpired() {
		// Remove expired entry
		s.data.Delete(key)
		s.metrics.recordMiss()
		return nil, ErrNotFound
	}

	// Return a copy to prevent external modifications
	result := make([]byte, len(entry.Value))
	copy(result, entry.Value)

	s.metrics.recordHit()
	return result, nil
}

// Set stores a value in the cache with optional TTL.
// ttl of 0 means the key never expires.
func (s *Store) Set(key string, value []byte, ttl time.Duration) error {
	start := time.Now()
	defer func() {
		s.recordLatency(time.Since(start))
	}()

	if len(key) == 0 {
		s.metrics.recordError()
		return errors.New("key cannot be empty")
	}

	if value == nil {
		s.metrics.recordError()
		return errors.New("value cannot be nil")
	}

	s.metrics.recordSet()

	// Create a copy of the value to prevent external modifications
	valueCopy := make([]byte, len(value))
	copy(valueCopy, value)

	// Calculate expiration time
	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	// Store the entry
	entry := &Entry{
		Value:     valueCopy,
		ExpiresAt: expiresAt,
	}

	s.data.Store(key, entry)
	return nil
}

// Delete removes a key from the cache.
func (s *Store) Delete(key string) {
	start := time.Now()
	defer func() {
		s.recordLatency(time.Since(start))
	}()

	s.metrics.recordDelete()
	s.data.Delete(key)
}

// Clear removes all keys from the cache.
func (s *Store) Clear() {
	s.data.Range(func(key, value interface{}) bool {
		s.data.Delete(key)
		return true
	})

	s.metrics.recordDelete()
}

// Len returns the approximate number of keys in the cache.
// Note: This may include expired keys that haven't been cleaned up yet.
func (s *Store) Len() int {
	count := 0
	s.data.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

// Exists checks if a key exists and hasn't expired.
func (s *Store) Exists(key string) bool {
	val, ok := s.data.Load(key)
	if !ok {
		return false
	}

	entry := val.(*Entry)
	if entry.IsExpired() {
		s.data.Delete(key)
		return false
	}

	return true
}

// Stats returns cache statistics.
func (s *Store) Stats() *StoreMetrics {
	return s.metrics.Copy()
}

// Close gracefully shuts down the store and its cleanup goroutine.
func (s *Store) Close() {
	close(s.stopCleanup)
	<-s.done
}

// Private helper methods

// cleanupLoop periodically removes expired entries.
func (s *Store) cleanupLoop() {
	defer close(s.done)

	ticker := time.NewTicker(s.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCleanup:
			return

		case <-ticker.C:
			s.cleanupExpired()
		}
	}
}

// cleanupExpired removes all expired entries from the cache.
// It performs a full scan of all entries to check for expiration.
func (s *Store) cleanupExpired() {
	// Full scan of all entries to check for expiration
	s.data.Range(func(key, value interface{}) bool {
		entry := value.(*Entry)
		if entry.IsExpired() {
			s.data.Delete(key)
		}
		return true
	})
}

// recordLatency records operation latency for metrics.
func (s *Store) recordLatency(d time.Duration) {
	s.metrics.mu.Lock()
	defer s.metrics.mu.Unlock()

	// Simple moving average
	if s.metrics.AvgLatency == 0 {
		s.metrics.AvgLatency = d
	} else {
		s.metrics.AvgLatency = (s.metrics.AvgLatency + d) / 2
	}
}

// StoreMetrics helper methods

// recordGet increments the Get counter.
func (m *StoreMetrics) recordGet() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Gets++
}

// recordSet increments the Set counter.
func (m *StoreMetrics) recordSet() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Sets++
}

// recordDelete increments the Delete counter.
func (m *StoreMetrics) recordDelete() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Deletes++
}

// recordHit increments the Hit counter.
func (m *StoreMetrics) recordHit() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Hits++
}

// recordMiss increments the Miss counter.
func (m *StoreMetrics) recordMiss() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Misses++
}

// recordError increments the Error counter.
func (m *StoreMetrics) recordError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Errors++
}

// Copy returns a copy of the metrics for safe access.
func (m *StoreMetrics) Copy() *StoreMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return &StoreMetrics{
		Gets:       m.Gets,
		Sets:       m.Sets,
		Deletes:    m.Deletes,
		Hits:       m.Hits,
		Misses:     m.Misses,
		Errors:     m.Errors,
		TotalSize:  m.TotalSize,
		AvgLatency: m.AvgLatency,
	}
}
