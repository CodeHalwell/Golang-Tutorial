package cache

import (
	"testing"
	"time"
)

// TestBasicGetSet verifies basic get and set operations.
func TestBasicGetSet(t *testing.T) {
	store := NewStore()
	defer store.Close()

	// Set a value
	err := store.Set("key1", []byte("value1"), 0)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get the value
	val, err := store.Get("key1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if string(val) != "value1" {
		t.Errorf("Expected 'value1', got '%s'", string(val))
	}
}

// TestGetNonExistent verifies that getting a non-existent key returns ErrNotFound.
func TestGetNonExistent(t *testing.T) {
	store := NewStore()
	defer store.Close()

	_, err := store.Get("nonexistent")
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// TestDelete verifies deletion of keys.
func TestDelete(t *testing.T) {
	store := NewStore()
	defer store.Close()

	// Set a value
	store.Set("key1", []byte("value1"), 0)

	// Delete it
	store.Delete("key1")

	// Try to get it
	_, err := store.Get("key1")
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound after delete, got %v", err)
	}
}

// TestTTLExpiration verifies that TTL expiration works correctly.
func TestTTLExpiration(t *testing.T) {
	store := NewStore()
	defer store.Close()

	// Set with 100ms TTL
	err := store.Set("expiring", []byte("value"), 100*time.Millisecond)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Should exist immediately
	_, err = store.Get("expiring")
	if err != nil {
		t.Errorf("Value should exist immediately after Set")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should be expired now
	_, err = store.Get("expiring")
	if err != ErrNotFound {
		t.Errorf("Value should be expired after TTL, got %v", err)
	}
}

// TestExists verifies the Exists method.
func TestExists(t *testing.T) {
	store := NewStore()
	defer store.Close()

	// Non-existent key
	if store.Exists("nonexistent") {
		t.Error("Exists should return false for non-existent key")
	}

	// After setting
	store.Set("key1", []byte("value"), 0)
	if !store.Exists("key1") {
		t.Error("Exists should return true after Set")
	}

	// After deleting
	store.Delete("key1")
	if store.Exists("key1") {
		t.Error("Exists should return false after Delete")
	}
}

// TestClear verifies the Clear method.
func TestClear(t *testing.T) {
	store := NewStore()
	defer store.Close()

	// Add multiple keys
	store.Set("key1", []byte("value1"), 0)
	store.Set("key2", []byte("value2"), 0)
	store.Set("key3", []byte("value3"), 0)

	if store.Len() != 3 {
		t.Errorf("Expected 3 keys, got %d", store.Len())
	}

	// Clear all
	store.Clear()

	if store.Len() != 0 {
		t.Errorf("Expected 0 keys after Clear, got %d", store.Len())
	}
}

// TestLen verifies the Len method.
func TestLen(t *testing.T) {
	store := NewStore()
	defer store.Close()

	if store.Len() != 0 {
		t.Errorf("New store should be empty, got %d keys", store.Len())
	}

	store.Set("key1", []byte("value1"), 0)
	if store.Len() != 1 {
		t.Errorf("Expected 1 key, got %d", store.Len())
	}

	store.Set("key2", []byte("value2"), 0)
	if store.Len() != 2 {
		t.Errorf("Expected 2 keys, got %d", store.Len())
	}
}

// TestValueIsolation verifies that returned values are isolated from the store.
func TestValueIsolation(t *testing.T) {
	store := NewStore()
	defer store.Close()

	original := []byte("value")
	store.Set("key", original, 0)

	// Get the value
	retrieved, _ := store.Get("key")

	// Modify the retrieved value
	retrieved[0] = 'X'

	// Get again and verify it's unchanged
	retrieved2, _ := store.Get("key")
	if string(retrieved2) != "value" {
		t.Errorf("Store value was modified externally: %s", string(retrieved2))
	}
}

// TestEmptyKey verifies that empty keys are rejected.
func TestEmptyKey(t *testing.T) {
	store := NewStore()
	defer store.Close()

	err := store.Set("", []byte("value"), 0)
	if err == nil {
		t.Error("Set with empty key should fail")
	}
}

// TestNilValue verifies that nil values are rejected.
func TestNilValue(t *testing.T) {
	store := NewStore()
	defer store.Close()

	err := store.Set("key", nil, 0)
	if err == nil {
		t.Error("Set with nil value should fail")
	}
}

// TestConcurrentAccess verifies thread-safety under concurrent access.
func TestConcurrentAccess(t *testing.T) {
	store := NewStore()
	defer store.Close()

	const numGoroutines = 100
	const operationsPerGoroutine = 100

	done := make(chan struct{}, numGoroutines)

	// Spawn concurrent writers
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < operationsPerGoroutine; j++ {
				key := "key"
				store.Set(key, []byte{byte(id), byte(j)}, 0)
			}
			done <- struct{}{}
		}(i)
	}

	// Spawn concurrent readers
	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < operationsPerGoroutine; j++ {
				store.Get("key")
			}
			done <- struct{}{}
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 2*numGoroutines; i++ {
		<-done
	}

	// Should complete without panic or race condition
}

// TestMetrics verifies that metrics are tracked correctly.
func TestMetrics(t *testing.T) {
	store := NewStore()
	defer store.Close()

	// Perform operations
	store.Set("key1", []byte("value1"), 0)
	store.Set("key2", []byte("value2"), 0)
	store.Get("key1")       // Hit
	store.Get("key1")       // Hit
	store.Get("nonexistent") // Miss
	store.Delete("key1")

	// Check metrics
	stats := store.Stats()

	if stats.Sets != 2 {
		t.Errorf("Expected 2 Sets, got %d", stats.Sets)
	}

	if stats.Gets != 3 {
		t.Errorf("Expected 3 Gets, got %d", stats.Gets)
	}

	if stats.Hits != 2 {
		t.Errorf("Expected 2 Hits, got %d", stats.Hits)
	}

	if stats.Misses != 1 {
		t.Errorf("Expected 1 Miss, got %d", stats.Misses)
	}

	if stats.Deletes != 1 {
		t.Errorf("Expected 1 Delete, got %d", stats.Deletes)
	}
}

// BenchmarkGet benchmarks Get operations.
func BenchmarkGet(b *testing.B) {
	store := NewStore()
	defer store.Close()

	store.Set("key", []byte("value"), 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Get("key")
	}
}

// BenchmarkSet benchmarks Set operations.
func BenchmarkSet(b *testing.B) {
	store := NewStore()
	defer store.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Set("key", []byte("value"), 0)
	}
}

// BenchmarkConcurrentGet benchmarks concurrent Get operations.
func BenchmarkConcurrentGet(b *testing.B) {
	store := NewStore()
	defer store.Close()

	store.Set("key", []byte("value"), 0)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			store.Get("key")
		}
	})
}

// BenchmarkConcurrentSet benchmarks concurrent Set operations.
func BenchmarkConcurrentSet(b *testing.B) {
	store := NewStore()
	defer store.Close()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			store.Set("key", []byte("value"), 0)
		}
	})
}
