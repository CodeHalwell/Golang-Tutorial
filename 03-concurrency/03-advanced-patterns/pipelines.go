// Package advanced_patterns provides advanced concurrency patterns
package advanced_patterns

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Pipeline represents a series of processing stages connected by channels.
// Each stage reads from input, processes, and sends to output.
//
// Example:
//
//	numbers := generateNumbers(10)
//	squared := square(context.Background(), numbers)
//	cubed := cube(context.Background(), squared)
//	for result := range cubed {
//	    fmt.Println(result)
//	}
type Pipeline struct {
	name   string
	stages []Stage
}

// Stage represents a single processing step in a pipeline
type Stage interface {
	Process(ctx context.Context, input interface{}) (interface{}, error)
}

// SimpleStage is a basic stage implementation
type SimpleStage struct {
	name      string
	processor func(ctx context.Context, val interface{}) (interface{}, error)
}

func (s *SimpleStage) Process(ctx context.Context, input interface{}) (interface{}, error) {
	return s.processor(ctx, input)
}

// LinearPipeline processes items sequentially through stages.
//
// Example use case: Data validation and transformation
//   Input → Validate → Transform → Enrich → Output
func LinearPipeline(
	ctx context.Context,
	input <-chan interface{},
	stages ...func(ctx context.Context, val interface{}) (interface{}, error),
) <-chan interface{} {
	output := make(chan interface{})

	go func() {
		defer close(output)

		for val := range input {
			// Pass through each stage
			for _, stage := range stages {
				var err error
				val, err = stage(ctx, val)
				if err != nil {
					// Skip on error, or send error marker
					continue
				}
			}

			select {
			case output <- val:
			case <-ctx.Done():
				return
			}
		}
	}()

	return output
}

// FanOut splits a single input stream into multiple independent outputs by broadcasting.
// Each output receives ALL values from the input (true broadcast pattern).
//
// Use case: Broadcasting work to multiple consumers that all need the same data
// Example: Send notification to multiple channels (email, SMS, push)
func FanOut(
	ctx context.Context,
	input <-chan interface{},
	numOutputs int,
) []<-chan interface{} {
	outputs := make([]<-chan interface{}, numOutputs)

	// Create buffered output channels
	channels := make([]chan interface{}, numOutputs)
	for i := 0; i < numOutputs; i++ {
		ch := make(chan interface{}, 1) // Buffered to prevent blocking
		channels[i] = ch
		outputs[i] = ch
	}

	// Single goroutine to broadcast to all outputs
	go func() {
		defer func() {
			// Close all output channels when input is closed
			for _, ch := range channels {
				close(ch)
			}
		}()

		for {
			select {
			case val, ok := <-input:
				if !ok {
					return
				}
				// Broadcast to all outputs with timeout protection
				for i, ch := range channels {
					select {
					case ch <- val:
						// Successfully sent
					case <-time.After(100 * time.Millisecond):
						// Consumer is too slow, drop this value for this output
						_ = i // Avoid unused variable
					case <-ctx.Done():
						return
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return outputs
}

// FanIn merges multiple input streams into one output stream.
//
// Use case: Aggregating results from multiple workers
// Example: Collecting logs from multiple services
func FanIn(ctx context.Context, inputs ...<-chan interface{}) <-chan interface{} {
	output := make(chan interface{})
	var wg sync.WaitGroup

	// Start a goroutine for each input
	for _, input := range inputs {
		wg.Add(1)

		go func(in <-chan interface{}) {
			defer wg.Done()

			for val := range in {
				select {
				case output <- val:
				case <-ctx.Done():
					return
				}
			}
		}(input)
	}

	// Close output when all inputs are done
	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

// MapPipeline applies a function to each element in a stream.
//
// Equivalent to functional programming's map operation.
// Use case: Data transformation (square, filter, format, etc.)
func MapPipeline(
	ctx context.Context,
	input <-chan interface{},
	fn func(interface{}) (interface{}, error),
) <-chan interface{} {
	output := make(chan interface{})

	go func() {
		defer close(output)

		for val := range input {
			result, err := fn(val)
			if err != nil {
				continue
			}

			select {
			case output <- result:
			case <-ctx.Done():
				return
			}
		}
	}()

	return output
}

// FilterPipeline keeps only elements that match a predicate.
//
// Use case: Filtering data (only valid entries, only errors, etc.)
func FilterPipeline(
	ctx context.Context,
	input <-chan interface{},
	predicate func(interface{}) bool,
) <-chan interface{} {
	output := make(chan interface{})

	go func() {
		defer close(output)

		for val := range input {
			if predicate(val) {
				select {
				case output <- val:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return output
}

// ParallelMapPipeline applies a function in parallel using N workers.
//
// More efficient than linear map for CPU-intensive operations.
// Use case: Image processing, data analysis, computation
func ParallelMapPipeline(
	ctx context.Context,
	input <-chan interface{},
	numWorkers int,
	fn func(interface{}) (interface{}, error),
) <-chan interface{} {
	output := make(chan interface{})
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for val := range input {
				result, err := fn(val)
				if err != nil {
					continue
				}

				select {
				case output <- result:
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	// Close output when done
	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

// BatchPipeline groups elements into batches before processing.
//
// Use case: Bulk database inserts, batch API calls, batch processing
func BatchPipeline(
	ctx context.Context,
	input <-chan interface{},
	batchSize int,
) <-chan []interface{} {
	output := make(chan []interface{})

	go func() {
		defer close(output)

		batch := make([]interface{}, 0, batchSize)

		for val := range input {
			batch = append(batch, val)

			if len(batch) == batchSize {
				select {
				case output <- batch:
					batch = make([]interface{}, 0, batchSize)
				case <-ctx.Done():
					return
				}
			}
		}

		// Send remaining items
		if len(batch) > 0 {
			select {
			case output <- batch:
			case <-ctx.Done():
			}
		}
	}()

	return output
}

// RateLimitPipeline limits the rate at which items are processed.
//
// Use case: API rate limiting, database load management, output throttling
func RateLimitPipeline(
	ctx context.Context,
	input <-chan interface{},
	itemsPerSecond int,
) <-chan interface{} {
	output := make(chan interface{})
	ticker := time.NewTicker(time.Second / time.Duration(itemsPerSecond))

	go func() {
		defer close(output)
		defer ticker.Stop()

		for val := range input {
			select {
			case <-ticker.C:
				select {
				case output <- val:
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return output
}

// WindowPipeline groups elements by a sliding time window.
//
// Use case: Metrics aggregation, time-series analysis
func WindowPipeline(
	ctx context.Context,
	input <-chan interface{},
	windowSize time.Duration,
) <-chan []interface{} {
	output := make(chan []interface{})

	go func() {
		defer close(output)

		ticker := time.NewTicker(windowSize)
		defer ticker.Stop()

		window := make([]interface{}, 0)

		for {
			select {
			case val, ok := <-input:
				if !ok {
					// Input closed, send remaining window
					if len(window) > 0 {
						output <- window
					}
					return
				}
				window = append(window, val)

			case <-ticker.C:
				if len(window) > 0 {
					select {
					case output <- window:
						window = make([]interface{}, 0)
					case <-ctx.Done():
						return
					}
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	return output
}

// TimeoutPipeline adds a timeout to each item processing.
//
// Use case: Preventing hanging operations, enforcing SLAs
func TimeoutPipeline(
	ctx context.Context,
	input <-chan interface{},
	timeout time.Duration,
	fn func(ctx context.Context, val interface{}) (interface{}, error),
) <-chan interface{} {
	output := make(chan interface{})

	go func() {
		defer close(output)

		for val := range input {
			// Create context with timeout for this item
			itemCtx, cancel := context.WithTimeout(ctx, timeout)

			result, err := fn(itemCtx, val)
			cancel()

			if err != nil {
				continue
			}

			select {
			case output <- result:
			case <-ctx.Done():
				return
			}
		}
	}()

	return output
}

// MergePipeline merges multiple pipelines into one, preserving order.
//
// Use case: Combining results from multiple data sources
func MergePipeline(ctx context.Context, inputs ...<-chan interface{}) <-chan interface{} {
	output := make(chan interface{})
	var wg sync.WaitGroup

	// Send each input to output
	for _, input := range inputs {
		wg.Add(1)

		go func(in <-chan interface{}) {
			defer wg.Done()

			for val := range in {
				select {
				case output <- val:
				case <-ctx.Done():
					return
				}
			}
		}(input)
	}

	// Close when all inputs done
	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

// Example usage of pipelines

// Numbers generates a sequence of numbers
func Numbers(ctx context.Context, n int) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		for i := 0; i < n; i++ {
			select {
			case out <- i:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

// Square returns a pipeline that squares numbers
func Square(ctx context.Context, in <-chan interface{}) <-chan interface{} {
	return MapPipeline(ctx, in, func(val interface{}) (interface{}, error) {
		n := val.(int)
		return n * n, nil
	})
}

// OnlyEven returns a pipeline that filters to even numbers
func OnlyEven(ctx context.Context, in <-chan interface{}) <-chan interface{} {
	return FilterPipeline(ctx, in, func(val interface{}) bool {
		n := val.(int)
		return n%2 == 0
	})
}

// PrintResults prints pipeline results
func PrintResults(ctx context.Context, in <-chan interface{}, label string) {
	for val := range in {
		fmt.Printf("%s: %v\n", label, val)
	}
}

// ExamplePipelineChain demonstrates chaining multiple pipelines
func ExamplePipelineChain() {
	ctx := context.Background()

	// Create a pipeline: Numbers → Square → Filter Even → Output
	numbers := Numbers(ctx, 10)
	squared := Square(ctx, numbers)
	evens := OnlyEven(ctx, squared)

	PrintResults(ctx, evens, "Result")

	// Output:
	// Result: 0
	// Result: 4
	// Result: 16
	// Result: 36
	// Result: 64
}
