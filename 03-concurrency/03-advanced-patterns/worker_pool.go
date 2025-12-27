// Package advanced_patterns demonstrates production-grade concurrency patterns
package advanced_patterns

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// WorkerPool implements a pool of workers that process tasks from a queue.
// This pattern is essential for managing resource usage and throughput.
type WorkerPool struct {
	// taskChan receives work to be processed
	taskChan chan Task

	// resultChan sends back completed work
	resultChan chan Result

	// numWorkers is the fixed number of workers
	numWorkers int

	// wg tracks all active goroutines
	wg sync.WaitGroup

	// ctx controls cancellation
	ctx context.Context
	cancel context.CancelFunc

	// metrics tracks pool performance
	metrics *PoolMetrics
}

// Task represents a unit of work to be processed
type Task struct {
	ID    int
	Data  []byte
	Delay time.Duration // Simulates work duration
}

// Result represents the outcome of a processed task
type Result struct {
	TaskID    int
	Output    []byte
	Duration  time.Duration
	Error     error
	WorkerID  int
	Timestamp time.Time
}

// PoolMetrics tracks worker pool statistics
type PoolMetrics struct {
	TasksProcessed  int64
	TasksFailed     int64
	TotalDuration   time.Duration
	AvgTaskDuration time.Duration
	mu              sync.Mutex
}

// NewWorkerPool creates a new worker pool with the specified number of workers.
//
// Example usage:
//
//	pool := NewWorkerPool(5)
//	pool.Submit(Task{ID: 1, Data: []byte("data")})
//	result := <- pool.Results()
//	pool.Close()
func NewWorkerPool(numWorkers int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &WorkerPool{
		taskChan:   make(chan Task, numWorkers*2), // Buffered channel for throughput
		resultChan: make(chan Result, numWorkers*2),
		numWorkers: numWorkers,
		ctx:        ctx,
		cancel:     cancel,
		metrics:    &PoolMetrics{},
	}

	// Start worker goroutines
	pool.wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go pool.worker(i)
	}

	return pool
}

// Submit adds a task to the work queue.
// Non-blocking; returns immediately if queue is not full.
func (p *WorkerPool) Submit(task Task) error {
	select {
	case p.taskChan <- task:
		return nil
	case <-p.ctx.Done():
		return fmt.Errorf("worker pool is closed")
	default:
		return fmt.Errorf("task queue is full")
	}
}

// Results returns the channel where results are sent.
func (p *WorkerPool) Results() <-chan Result {
	return p.resultChan
}

// Close gracefully shuts down the worker pool.
// Waits for all in-flight tasks to complete.
func (p *WorkerPool) Close() error {
	p.cancel()

	// Wait for all workers to finish
	p.wg.Wait()

	// Close result channel
	close(p.resultChan)

	return nil
}

// Metrics returns a copy of the pool metrics.
func (p *WorkerPool) Metrics() PoolMetrics {
	p.metrics.mu.Lock()
	defer p.metrics.mu.Unlock()
	return *p.metrics
}

// Private helper methods

// worker processes tasks from the queue.
// Each worker runs in its own goroutine.
func (p *WorkerPool) worker(workerID int) {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			// Pool is shutting down
			return

		case task, ok := <-p.taskChan:
			if !ok {
				// Task channel closed
				return
			}

			// Process the task
			result := p.processTask(task, workerID)

			// Send result (non-blocking, with context awareness)
			select {
			case p.resultChan <- result:
			case <-p.ctx.Done():
				return
			}
		}
	}
}

// processTask executes a single task.
func (p *WorkerPool) processTask(task Task, workerID int) Result {
	start := time.Now()
	result := Result{
		TaskID:    task.ID,
		WorkerID:  workerID,
		Timestamp: start,
	}

	// Simulate work
	time.Sleep(task.Delay)

	// Process the data (example: uppercase conversion)
	result.Output = make([]byte, len(task.Data))
	for i, b := range task.Data {
		if b >= 'a' && b <= 'z' {
			result.Output[i] = b - 32 // Convert to uppercase
		} else {
			result.Output[i] = b
		}
	}

	result.Duration = time.Since(start)

	// Update metrics
	p.metrics.recordCompletion(result.Duration)

	return result
}

// recordCompletion updates metrics after task completion
func (m *PoolMetrics) recordCompletion(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.TasksProcessed++
	m.TotalDuration += duration

	// Calculate moving average
	if m.AvgTaskDuration == 0 {
		m.AvgTaskDuration = duration
	} else {
		m.AvgTaskDuration = (m.AvgTaskDuration + duration) / 2
	}
}

// BoundedWorkerPool is a more advanced variant that limits concurrent execution.
//
// Use this when you have:
// - Limited resources (connections, file handles, memory)
// - Need strict control over concurrency level
// - Want fairness in task processing
type BoundedWorkerPool struct {
	semaphore chan struct{}
	pool      *WorkerPool
	maxActive int
}

// NewBoundedWorkerPool creates a worker pool with bounded resource usage.
func NewBoundedWorkerPool(numWorkers int, maxConcurrent int) *BoundedWorkerPool {
	return &BoundedWorkerPool{
		semaphore: make(chan struct{}, maxConcurrent),
		pool:      NewWorkerPool(numWorkers),
		maxActive: maxConcurrent,
	}
}

// Submit acquires a resource before submitting the task.
func (bwp *BoundedWorkerPool) Submit(task Task) error {
	// Try to acquire a resource
	select {
	case bwp.semaphore <- struct{}{}:
		// Got a resource, submit the task
		if err := bwp.pool.Submit(task); err != nil {
			// Release the resource if submission failed
			<-bwp.semaphore
			return err
		}

		// Wrap result to release resource on completion
		go func() {
			<-bwp.pool.Results()
			<-bwp.semaphore // Release resource
		}()

		return nil

	default:
		return fmt.Errorf("all resources in use")
	}
}

// Results returns the underlying pool results channel
func (bwp *BoundedWorkerPool) Results() <-chan Result {
	return bwp.pool.Results()
}

// Close shuts down the bounded pool
func (bwp *BoundedWorkerPool) Close() error {
	return bwp.pool.Close()
}

// RateLimitedWorkerPool adds rate limiting to task submission.
//
// Use this when you need to:
// - Respect rate limits (API quotas, database connections)
// - Smooth out traffic spikes
// - Provide backpressure to callers
type RateLimitedWorkerPool struct {
	pool       *WorkerPool
	rateLimiter <-chan time.Time
	ticker     *time.Ticker
}

// NewRateLimitedWorkerPool creates a worker pool with rate limiting.
// tasksPerSecond controls the maximum task submission rate.
func NewRateLimitedWorkerPool(numWorkers int, tasksPerSecond int) *RateLimitedWorkerPool {
	ticker := time.NewTicker(time.Second / time.Duration(tasksPerSecond))

	return &RateLimitedWorkerPool{
		pool:        NewWorkerPool(numWorkers),
		rateLimiter: ticker.C,
		ticker:      ticker,
	}
}

// Submit respects the rate limit before submitting.
func (rwp *RateLimitedWorkerPool) Submit(task Task) error {
	// Wait for rate limit to allow submission
	<-rwp.rateLimiter

	// Submit the task
	return rwp.pool.Submit(task)
}

// Results returns the underlying pool results channel
func (rwp *RateLimitedWorkerPool) Results() <-chan Result {
	return rwp.pool.Results()
}

// Close shuts down the rate-limited pool
func (rwp *RateLimitedWorkerPool) Close() error {
	rwp.ticker.Stop()
	return rwp.pool.Close()
}
