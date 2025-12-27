package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ===== EXAMPLE 1: Basic Goroutine =====

func example1_BasicGoroutine() {
	fmt.Println("\n=== EXAMPLE 1: Basic Goroutine ===")

	go func() {
		fmt.Println("Hello from goroutine!")
	}()

	fmt.Println("Main continues immediately")

	// Must wait or main exits
	time.Sleep(100 * time.Millisecond)
}

// ===== EXAMPLE 2: Multiple Goroutines =====

func example2_MultipleGoroutines() {
	fmt.Println("\n=== EXAMPLE 2: Multiple Goroutines ===")

	for i := 1; i <= 3; i++ {
		go func(id int) {
			fmt.Printf("Goroutine %d executing\n", id)
			time.Sleep(50 * time.Millisecond)
			fmt.Printf("Goroutine %d done\n", id)
		}(i)
	}

	time.Sleep(200 * time.Millisecond)
	fmt.Println("All goroutines complete")
}

// ===== EXAMPLE 3: WaitGroup Synchronization =====

func example3_WaitGroup() {
	fmt.Println("\n=== EXAMPLE 3: WaitGroup Synchronization ===")

	var wg sync.WaitGroup

	// Register 5 goroutines
	wg.Add(5)

	for i := 1; i <= 5; i++ {
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Worker %d: starting\n", id)
			time.Sleep(50 * time.Millisecond)
			fmt.Printf("Worker %d: done\n", id)
		}(i)
	}

	fmt.Println("Waiting for workers...")
	wg.Wait()
	fmt.Println("All workers complete")
}

// ===== EXAMPLE 4: Race Condition Problem =====

func example4_RaceCondition() {
	fmt.Println("\n=== EXAMPLE 4: Race Condition Problem ===")

	counter := 0
	var wg sync.WaitGroup

	// Create 10 goroutines, each incrementing 100 times
	wg.Add(10)

	for g := 0; g < 10; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				counter++ // ✗ RACE CONDITION!
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Expected: 1000, Got: %d (likely less due to race)\n", counter)
	fmt.Println("Run with: go run -race main.go")
}

// ===== EXAMPLE 5: Mutex Solution =====

type SafeCounter struct {
	mu    sync.Mutex
	value int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *SafeCounter) Read() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func example5_MutexSolution() {
	fmt.Println("\n=== EXAMPLE 5: Mutex Solution ===")

	counter := &SafeCounter{}
	var wg sync.WaitGroup

	// Create 10 goroutines, each incrementing 100 times
	wg.Add(10)

	for g := 0; g < 10; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				counter.Increment()
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Thread-safe counter: %d (guaranteed correct)\n", counter.Read())
}

// ===== EXAMPLE 6: Atomic Operations =====

func example6_AtomicOperations() {
	fmt.Println("\n=== EXAMPLE 6: Atomic Operations ===")

	var counter int64
	var wg sync.WaitGroup

	wg.Add(10)

	for g := 0; g < 10; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				atomic.AddInt64(&counter, 1)
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Atomic counter: %d (lock-free)\n", atomic.LoadInt64(&counter))
}

// ===== EXAMPLE 7: Channel Communication =====

func example7_Channels() {
	fmt.Println("\n=== EXAMPLE 7: Channel Communication ===")

	results := make(chan string)

	go func() {
		fmt.Println("Goroutine: computing...")
		time.Sleep(100 * time.Millisecond)
		results <- "computation complete"
	}()

	fmt.Println("Main: waiting for result...")
	result := <-results
	fmt.Printf("Main: received '%s'\n", result)
}

// ===== EXAMPLE 8: Multiple Channels =====

func example8_MultipleChannels() {
	fmt.Println("\n=== EXAMPLE 8: Multiple Channels ===")

	results := make(chan int, 3)
	var wg sync.WaitGroup

	wg.Add(3)

	// Goroutine 1
	go func() {
		defer wg.Done()
		results <- 10
	}()

	// Goroutine 2
	go func() {
		defer wg.Done()
		results <- 20
	}()

	// Goroutine 3
	go func() {
		defer wg.Done()
		results <- 30
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect all results
	sum := 0
	for result := range results {
		fmt.Printf("Received: %d\n", result)
		sum += result
	}
	fmt.Printf("Total: %d\n", sum)
}

// ===== EXAMPLE 9: Loop Variable Capture Bug =====

func example9_LoopVariableBug() {
	fmt.Println("\n=== EXAMPLE 9: Loop Variable Capture Bug ===")

	var wg sync.WaitGroup

	fmt.Println("✗ Wrong way (all print 3):")
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			fmt.Printf("  i = %d\n", i) // All will be 3!
		}()
	}
	wg.Wait()

	fmt.Println("✓ Correct way (print 0, 1, 2):")
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func(id int) {
			defer wg.Done()
			fmt.Printf("  i = %d\n", id) // Correct!
		}(i)
	}
	wg.Wait()
}

// ===== EXAMPLE 10: Goroutine Panic Recovery =====

func example10_PanicRecovery() {
	fmt.Println("\n=== EXAMPLE 10: Goroutine Panic Recovery ===")

	var wg sync.WaitGroup

	// Safe goroutine wrapper with panic recovery
	safeGo := func(fn func(), id int) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Goroutine %d recovered: %v\n", id, r)
				}
			}()
			fn()
		}()
	}

	safeGo(func() {
		fmt.Println("Goroutine 1: working")
		time.Sleep(50 * time.Millisecond)
		fmt.Println("Goroutine 1: done")
	}, 1)

	safeGo(func() {
		panic("Error in goroutine 2")
	}, 2)

	wg.Wait()
	fmt.Println("All goroutines complete (even with panic)")
}

// ===== EXAMPLE 11: Worker Pool Pattern =====

func example11_WorkerPool() {
	fmt.Println("\n=== EXAMPLE 11: Worker Pool Pattern ===")

	numWorkers := 3
	jobs := make(chan int, 10)
	results := make(chan int, 10)
	var wg sync.WaitGroup

	// Start workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobs {
				fmt.Printf("Worker %d: processing %d\n", workerID, job)
				time.Sleep(50 * time.Millisecond)
				results <- job * 2
			}
		}(w)
	}

	// Send jobs
	go func() {
		for j := 1; j <= 5; j++ {
			fmt.Printf("Sending job %d\n", j)
			jobs <- j
		}
		close(jobs)
	}()

	// Wait for workers
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	fmt.Println("Results:")
	for result := range results {
		fmt.Printf("  %d\n", result)
	}
}

// ===== EXAMPLE 12: Fan-Out/Fan-In Pattern =====

func example12_FanOut() {
	fmt.Println("\n=== EXAMPLE 12: Fan-Out/Fan-In Pattern ===")

	// Generate work
	work := func() <-chan int {
		out := make(chan int)
		go func() {
			for i := 1; i <= 5; i++ {
				out <- i
			}
			close(out)
		}()
		return out
	}()

	// Fan out: broadcast each value to multiple workers
	// Each worker gets ALL values from the input
	fanOut := func(in <-chan int, numWorkers int) []<-chan int {
		outputs := make([]<-chan int, numWorkers)
		channels := make([]chan int, numWorkers)
		
		for i := 0; i < numWorkers; i++ {
			ch := make(chan int, 10) // Larger buffer to reduce blocking risk
			channels[i] = ch
			outputs[i] = ch
		}

		go func() {
			for val := range in {
				// Broadcast to all workers with timeout protection
				for i := 0; i < numWorkers; i++ {
					select {
					case channels[i] <- val:
						// Successfully sent
					case <-time.After(100 * time.Millisecond):
						// Worker is too slow, skip this value for this worker
						fmt.Printf("Warning: Worker %d is slow, skipping value\n", i)
					}
				}
			}
			// Close all output channels
			for i := 0; i < numWorkers; i++ {
				close(channels[i])
			}
		}()

		return outputs
	}

	// Process function for each worker
	process := func(id int, in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				result := n * (id + 1) // Different processing per worker
				fmt.Printf("Worker %d: %d -> %d\n", id, n, result)
				out <- result
			}
		}()
		return out
	}

	// Fan out work to 2 workers
	workerInputs := fanOut(work, 2)
	
	// Process with each worker
	ch1 := process(1, workerInputs[0])
	ch2 := process(2, workerInputs[1])

	// Fan in: merge results concurrently using sync.WaitGroup
	merged := make(chan int)
	var wg sync.WaitGroup
	
	wg.Add(2)
	
	// Collect from ch1
	go func() {
		defer wg.Done()
		for x := range ch1 {
			merged <- x
		}
	}()
	
	// Collect from ch2
	go func() {
		defer wg.Done()
		for x := range ch2 {
			merged <- x
		}
	}()
	
	// Close merged when all inputs are done
	go func() {
		wg.Wait()
		close(merged)
	}()

	// Collect all results
	for result := range merged {
		fmt.Printf("Got result: %d\n", result)
	}
}

// ===== EXAMPLE 13: Timeout Pattern =====

func example13_Timeout() {
	fmt.Println("\n=== EXAMPLE 13: Timeout Pattern ===")

	slowTask := func() <-chan string {
		out := make(chan string)
		go func() {
			fmt.Println("Task: starting (will take 2 seconds)")
			time.Sleep(2 * time.Second)
			out <- "Task complete"
		}()
		return out
	}

	// Attempt 1: Timeout before completion
	fmt.Println("Attempt 1 (timeout = 100ms):")
	select {
	case result := <-slowTask():
		fmt.Printf("Result: %s\n", result)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("Timeout!")
	}

	// Attempt 2: Wait long enough
	fmt.Println("\nAttempt 2 (timeout = 3s):")
	select {
	case result := <-slowTask():
		fmt.Printf("Result: %s\n", result)
	case <-time.After(3 * time.Second):
		fmt.Println("Timeout!")
	}
}

// ===== EXAMPLE 14: Context Cancellation =====

func example14_Cancellation() {
	fmt.Println("\n=== EXAMPLE 14: Cancellation Pattern ===")

	done := make(chan struct{})
	var wg sync.WaitGroup

	// Worker that can be cancelled
	worker := func(id int) {
		defer wg.Done()
		for {
			select {
			case <-done:
				fmt.Printf("Worker %d: cancelled\n", id)
				return
			default:
				fmt.Printf("Worker %d: working\n", id)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}

	wg.Add(3)
	go worker(1)
	go worker(2)
	go worker(3)

	time.Sleep(250 * time.Millisecond)

	fmt.Println("Cancelling all workers...")
	close(done)

	wg.Wait()
	fmt.Println("All workers stopped")
}

// ===== EXAMPLE 15: Error Handling in Goroutines =====

func example15_ErrorHandling() {
	fmt.Println("\n=== EXAMPLE 15: Error Handling in Goroutines ===")

	type Result struct {
		ID  int
		Err error
	}

	results := make(chan Result, 3)
	var wg sync.WaitGroup

	// Worker that might fail
	worker := func(id int) {
		defer wg.Done()
		if id == 2 {
			results <- Result{ID: id, Err: fmt.Errorf("worker %d failed", id)}
		} else {
			fmt.Printf("Worker %d: success\n", id)
			results <- Result{ID: id, Err: nil}
		}
	}

	wg.Add(3)
	for i := 1; i <= 3; i++ {
		go worker(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	for result := range results {
		if result.Err != nil {
			fmt.Printf("Error: %v\n", result.Err)
		} else {
			fmt.Printf("Worker %d completed successfully\n", result.ID)
		}
	}
}

// ===== EXAMPLE 16: Producer-Consumer Pattern =====

func example16_ProducerConsumer() {
	fmt.Println("\n=== EXAMPLE 16: Producer-Consumer Pattern ===")

	items := make(chan int, 5)
	var wg sync.WaitGroup

	// Producer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 5; i++ {
			fmt.Printf("Producing %d\n", i)
			items <- i
			time.Sleep(50 * time.Millisecond)
		}
		close(items)
	}()

	// Consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for item := range items {
			fmt.Printf("Consuming %d\n", item)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	wg.Wait()
	fmt.Println("Done")
}

// ===== EXAMPLE 17: Concurrent Map Access =====

type SafeMap struct {
	mu   sync.RWMutex
	data map[string]int
}

func (sm *SafeMap) Set(key string, value int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[key] = value
}

func (sm *SafeMap) Get(key string) int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.data[key]
}

func example17_ConcurrentMap() {
	fmt.Println("\n=== EXAMPLE 17: Concurrent Map Access ===")

	m := &SafeMap{data: make(map[string]int)}
	var wg sync.WaitGroup

	// Multiple writers
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			m.Set(fmt.Sprintf("key%d", id), id*10)
			fmt.Printf("Writer %d: set value %d\n", id, id*10)
		}(i)
	}

	// Multiple readers
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			val := m.Get(fmt.Sprintf("key%d", id))
			fmt.Printf("Reader %d: got value %d\n", id, val)
		}(i)
	}

	wg.Wait()
}

// ===== MAIN FUNCTION =====

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║   Chapter 1: Goroutines Basics - Comprehensive Examples    ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")

	example1_BasicGoroutine()
	example2_MultipleGoroutines()
	example3_WaitGroup()
	example4_RaceCondition()
	example5_MutexSolution()
	example6_AtomicOperations()
	example7_Channels()
	example8_MultipleChannels()
	example9_LoopVariableBug()
	example10_PanicRecovery()
	example11_WorkerPool()
	example12_FanOut()
	example13_Timeout()
	example14_Cancellation()
	example15_ErrorHandling()
	example16_ProducerConsumer()
	example17_ConcurrentMap()

	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║   Key Takeaways:                                           ║")
	fmt.Println("╠════════════════════════════════════════════════════════════╣")
	fmt.Println("║ 1. Goroutines are lightweight and easy to create           ║")
	fmt.Println("║ 2. Use WaitGroup for synchronization                       ║")
	fmt.Println("║ 3. Mutex protects shared data                              ║")
	fmt.Println("║ 4. Channels enable safe communication                      ║")
	fmt.Println("║ 5. Close channels from sender side                         ║")
	fmt.Println("║ 6. Pass loop variables as arguments                        ║")
	fmt.Println("║ 7. Use timeout pattern for deadlock protection             ║")
	fmt.Println("║ 8. Always wait for goroutines in main                      ║")
	fmt.Println("║ 9. Worker pools limit concurrent execution                 ║")
	fmt.Println("║ 10. Run tests with -race flag                              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
}
