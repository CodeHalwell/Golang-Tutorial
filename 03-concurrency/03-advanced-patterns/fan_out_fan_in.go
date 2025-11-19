package advanced_patterns

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// FanOutFanIn demonstrates the fan-out/fan-in concurrency pattern.
//
// Fan-Out: One input is replicated to multiple parallel workers.
// Fan-In: Multiple worker outputs are combined into one result stream.
//
// Use cases:
// - Parallel processing of independent items
// - Fetching from multiple sources concurrently
// - Sending notifications to multiple channels
// - Map/reduce style data processing

// FanOutFanInResult holds the result from a worker
type FanOutFanInResult struct {
	Index  int           // Which worker produced this
	Value  interface{}
	Error  error
	Time   time.Time
}

// SimpleFanOutFanIn is the basic pattern: split work, process in parallel, collect results
//
// Example:
//
//	input := []string{"a", "b", "c", "d"}
//	results := SimpleFanOutFanIn(ctx, input, 2, processItem)
//	for r := range results {
//	    fmt.Printf("Worker %d: %v\n", r.Index, r.Value)
//	}
func SimpleFanOutFanIn(
	ctx context.Context,
	items []interface{},
	numWorkers int,
	worker func(ctx context.Context, item interface{}) (interface{}, error),
) <-chan FanOutFanInResult {
	// Channel for distributing work
	work := make(chan interface{}, len(items))
	results := make(chan FanOutFanInResult, len(items))

	// Send items to work queue
	go func() {
		defer close(work)
		for _, item := range items {
			select {
			case work <- item:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()

			for item := range work {
				result := FanOutFanInResult{
					Index: workerID,
					Time:  time.Now(),
				}

				result.Value, result.Error = worker(ctx, item)

				select {
				case results <- result:
				case <-ctx.Done():
					return
				}
			}
		}(i)
	}

	// Close results when all workers done
	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

// DynamicFanOut sends the same item to multiple parallel processors.
//
// Use case: Broadcasting an event to multiple handlers
//
// Example:
//
//	event := "user_registered"
//	event_handlers := []func{"send_email", "log_event", "update_analytics"}
//	results := DynamicFanOut(ctx, event, event_handlers)
func DynamicFanOut(
	ctx context.Context,
	item interface{},
	handlers []func(ctx context.Context, item interface{}) error,
) []<-chan error {
	resultChans := make([]<-chan error, len(handlers))

	for i, handler := range handlers {
		resultChan := make(chan error, 1)
		resultChans[i] = resultChan

		go func(h func(context.Context, interface{}) error, ch chan<- error) {
			defer close(ch)

			select {
			case ch <- h(ctx, item):
			case <-ctx.Done():
			}
		}(handler, resultChan)
	}

	return resultChans
}

// FanOutFanInWithOrdering preserves the order of results matching input order.
//
// Useful when you need results in the same order as input items.
//
// Example:
//
//	input := []int{1, 2, 3, 4, 5}
//	results := FanOutFanInWithOrdering(ctx, input, 2, squareAndDouble)
//	for i, r := range results {
//	    fmt.Printf("Item %d: %v\n", i, r)
//	}
func FanOutFanInWithOrdering(
	ctx context.Context,
	items []interface{},
	numWorkers int,
	worker func(ctx context.Context, item interface{}) (interface{}, error),
) []interface{} {
	// Results slice to preserve order
	results := make([]interface{}, len(items))
	resultsChan := make(chan FanOutFanInResult, len(items))

	// Channel for distributing work with indices
	work := make(chan struct {
		index int
		item  interface{}
	}, len(items))

	// Send work with indices
	go func() {
		defer close(work)
		for i, item := range items {
			select {
			case work <- struct {
				index int
				item  interface{}
			}{i, item}:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for w := range work {
				value, err := worker(ctx, w.item)

				select {
				case resultsChan <- FanOutFanInResult{
					Index: w.index,
					Value: value,
					Error: err,
					Time:  time.Now(),
				}:
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	// Collect results in order
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Wait for all results and place in order
	for res := range resultsChan {
		results[res.Index] = res.Value
	}

	return results
}

// BatchFanOutFanIn processes items in batches with fan-out/fan-in.
//
// Use case: Database batch inserts, batch API calls
//
// Example:
//
//	items := []string{...} (1000 items)
//	batchSize := 100
//	numWorkers := 4
//	results := BatchFanOutFanIn(ctx, items, batchSize, numWorkers, processBatch)
func BatchFanOutFanIn(
	ctx context.Context,
	items []interface{},
	batchSize int,
	numWorkers int,
	batchWorker func(ctx context.Context, batch []interface{}) (interface{}, error),
) <-chan FanOutFanInResult {
	// Create batches
	var batches [][]interface{}
	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}
		batches = append(batches, items[i:end])
	}

	// Process batches as items
	batchItems := make([]interface{}, len(batches))
	for i, batch := range batches {
		batchItems[i] = batch
	}

	return SimpleFanOutFanIn(
		ctx,
		batchItems,
		numWorkers,
		func(ctx context.Context, item interface{}) (interface{}, error) {
			batch := item.([]interface{})
			return batchWorker(ctx, batch)
		},
	)
}

// NestedFanOutFanIn applies fan-out/fan-in at multiple levels.
//
// Use case: Multi-stage processing pipeline with parallelism at each stage
//
// Example: Image processing
//
//	Level 1: Download images from 4 sources in parallel (4 workers)
//	Level 2: Process downloaded images (resize, filter, etc) in parallel (8 workers)
//	Level 3: Upload processed images (3 workers)
func NestedFanOutFanIn(
	ctx context.Context,
	input []interface{},
	stages []struct {
		numWorkers int
		worker     func(ctx context.Context, item interface{}) (interface{}, error)
	},
) <-chan FanOutFanInResult {
	current := input

	for _, stage := range stages {
		// Convert results to items for next stage
		var nextInput []interface{}
		results := SimpleFanOutFanIn(ctx, current, stage.numWorkers, stage.worker)

		for res := range results {
			if res.Error == nil {
				nextInput = append(nextInput, res.Value)
			}
		}

		current = nextInput
	}

	// Final stage: collect results
	results := make(chan FanOutFanInResult, len(current))
	for i, item := range current {
		results <- FanOutFanInResult{
			Index: i,
			Value: item,
			Time:  time.Now(),
		}
	}
	close(results)

	return results
}

// ErrorHandlingFanOutFanIn demonstrates proper error handling across workers.
//
// Collects both successful results and errors.
type ErrorHandlingResult struct {
	SuccessCount int
	ErrorCount   int
	Results      []FanOutFanInResult
	FirstError   error
}

func ErrorHandlingFanOutFanIn(
	ctx context.Context,
	items []interface{},
	numWorkers int,
	worker func(ctx context.Context, item interface{}) (interface{}, error),
) *ErrorHandlingResult {
	results := SimpleFanOutFanIn(ctx, items, numWorkers, worker)

	var result ErrorHandlingResult
	result.Results = make([]FanOutFanInResult, 0)

	for res := range results {
		result.Results = append(result.Results, res)

		if res.Error != nil {
			result.ErrorCount++
			if result.FirstError == nil {
				result.FirstError = res.Error
			}
		} else {
			result.SuccessCount++
		}
	}

	return &result
}

// RateLimitedFanOutFanIn applies rate limiting to the fan-out phase.
//
// Use case: API rate limits, database connection limits
//
// Example:
//
//	items := []url{...} (1000 URLs)
//	maxConcurrent := 10  // Only 10 concurrent requests
//	results := RateLimitedFanOutFanIn(ctx, items, maxConcurrent, fetchURL)
func RateLimitedFanOutFanIn(
	ctx context.Context,
	items []interface{},
	maxConcurrent int,
	worker func(ctx context.Context, item interface{}) (interface{}, error),
) <-chan FanOutFanInResult {
	resultsChan := make(chan FanOutFanInResult, len(items))
	workChan := make(chan interface{}, len(items))

	// Send all work
	go func() {
		defer close(workChan)
		for _, item := range items {
			select {
			case workChan <- item:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Process with rate limit
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxConcurrent)

	go func() {
		for item := range workChan {
			// Acquire semaphore
			select {
			case semaphore <- struct{}{}:
			case <-ctx.Done():
				return
			}

			wg.Add(1)

			go func(i interface{}) {
				defer wg.Done()
				defer func() { <-semaphore }() // Release semaphore

				value, err := worker(ctx, i)

				select {
				case resultsChan <- FanOutFanInResult{
					Value: value,
					Error: err,
					Time:  time.Now(),
				}:
				case <-ctx.Done():
				}
			}(item)
		}
	}()

	// Close results when done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	return resultsChan
}

// TimeoutFanOutFanIn applies a timeout to each worker task.
//
// Example:
//
//	items := []url{...}
//	timeout := 5 * time.Second
//	results := TimeoutFanOutFanIn(ctx, items, 4, timeout, fetchWithTimeout)
func TimeoutFanOutFanIn(
	ctx context.Context,
	items []interface{},
	numWorkers int,
	timeout time.Duration,
	worker func(ctx context.Context, item interface{}) (interface{}, error),
) <-chan FanOutFanInResult {
	return SimpleFanOutFanIn(
		ctx,
		items,
		numWorkers,
		func(ctx context.Context, item interface{}) (interface{}, error) {
			// Create context with timeout for this task
			taskCtx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			return worker(taskCtx, item)
		},
	)
}

// Example demonstration

// ExampleFanOutFanIn shows a simple fan-out/fan-in example
func ExampleFanOutFanIn() {
	ctx := context.Background()

	// Input items
	items := []interface{}{1, 2, 3, 4, 5}

	// Worker function that squares numbers
	worker := func(ctx context.Context, item interface{}) (interface{}, error) {
		n := item.(int)
		return n * n, nil
	}

	// Process with 2 workers
	results := SimpleFanOutFanIn(ctx, items, 2, worker)

	fmt.Println("Results:")
	for res := range results {
		fmt.Printf("Worker %d: %v\n", res.Index, res.Value)
	}
}
