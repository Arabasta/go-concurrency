package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// This file demonstrates various synchronization primitives in Go beyond mutexes,
// including channels, select statements, sync.Once, sync.Pool, and context.

// Example 1: Channel-based synchronization
func channelSyncExample() {
	fmt.Println("=== Channel Synchronization Example ===")

	// Unbuffered channel for synchronization
	done := make(chan bool)

	go func() {
		fmt.Println("Worker starting...")
		time.Sleep(2 * time.Second)
		fmt.Println("Worker done!")
		done <- true // Signal completion
	}()

	fmt.Println("Waiting for worker...")
	<-done // Wait for signal
	fmt.Println("Worker completed, main continues\n")
}

// Example 2: Select statement for non-blocking operations
func selectStatementExample() {
	fmt.Println("=== Select Statement Example ===")

	ch1 := make(chan string)
	ch2 := make(chan string)

	// Goroutine 1
	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "message from ch1"
	}()

	// Goroutine 2
	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "message from ch2"
	}()

	// Select with timeout
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("Received:", msg1)
		case msg2 := <-ch2:
			fmt.Println("Received:", msg2)
		case <-time.After(3 * time.Second):
			fmt.Println("Timeout!")
		}
	}

	// Non-blocking select with default
	select {
	case msg := <-ch1:
		fmt.Println("Got message:", msg)
	default:
		fmt.Println("No message available")
	}

	fmt.Println()
}

// Example 3: sync.Once for one-time initialization
var (
	instance *Singleton
	once     sync.Once
)

type Singleton struct {
	value string
}

func GetInstance() *Singleton {
	once.Do(func() {
		fmt.Println("Creating singleton instance...")
		instance = &Singleton{value: "I'm a singleton!"}
		time.Sleep(100 * time.Millisecond) // Simulate expensive initialization
	})
	return instance
}

func syncOnceExample() {
	fmt.Println("=== sync.Once Example ===")

	var wg sync.WaitGroup

	// Multiple goroutines trying to get the singleton
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			s := GetInstance()
			fmt.Printf("Goroutine %d got singleton: %s\n", id, s.value)
		}(i)
	}

	wg.Wait()
	fmt.Println("sync.Once example completed!\n")
}

// Example 4: sync.Pool for object reuse
type ExpensiveObject struct {
	data []byte
}

func NewExpensiveObject() *ExpensiveObject {
	fmt.Println("Creating expensive object...")
	return &ExpensiveObject{
		data: make([]byte, 1024), // Simulate expensive allocation
	}
}

var objectPool = sync.Pool{
	New: func() interface{} {
		return NewExpensiveObject()
	},
}

func syncPoolExample() {
	fmt.Println("=== sync.Pool Example ===")

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Get object from pool
			obj := objectPool.Get().(*ExpensiveObject)

			// Use the object
			fmt.Printf("Worker %d using object with %d bytes\n", id, len(obj.data))
			time.Sleep(100 * time.Millisecond)

			// Reset and return to pool
			obj.data = obj.data[:0] // Clear data
			objectPool.Put(obj)

			fmt.Printf("Worker %d returned object to pool\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Println("sync.Pool example completed!\n")
}

// Example 5: Context for cancellation and timeouts
func contextExample() {
	fmt.Println("=== Context Example ===")

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result := make(chan string, 1)

	go func() {
		// Simulate long-running operation
		select {
		case <-time.After(3 * time.Second):
			result <- "work completed"
		case <-ctx.Done():
			fmt.Println("Work cancelled:", ctx.Err())
			return
		}
	}()

	select {
	case res := <-result:
		fmt.Println("Result:", res)
	case <-ctx.Done():
		fmt.Println("Main: operation timed out")
	}

	// Context with cancellation
	ctx2, cancel2 := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx2.Done():
				fmt.Println("Background worker stopped")
				return
			default:
				fmt.Println("Background worker running...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	time.Sleep(1500 * time.Millisecond)
	cancel2()                          // Cancel the context
	time.Sleep(100 * time.Millisecond) // Give time for cleanup

	fmt.Println("Context example completed!\n")
}

// Example 6: Pipeline pattern with channels
func pipelineExample() {
	fmt.Println("=== Pipeline Pattern Example ===")

	// Stage 1: Generate numbers
	numbers := func() <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for i := 1; i <= 5; i++ {
				out <- i
			}
		}()
		return out
	}

	// Stage 2: Square numbers
	square := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				out <- n * n
			}
		}()
		return out
	}

	// Stage 3: Double numbers
	double := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				out <- n * 2
			}
		}()
		return out
	}

	// Build pipeline
	pipeline := double(square(numbers()))

	// Consume results
	for result := range pipeline {
		fmt.Printf("Pipeline result: %d\n", result)
	}

	fmt.Println("Pipeline example completed!\n")
}

// Example 7: Worker pool pattern
func workerPoolExample() {
	fmt.Println("=== Worker Pool Example ===")

	const numWorkers = 3
	const numJobs = 10

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// Start workers
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				fmt.Printf("Worker %d processing job %d\n", id, job)
				time.Sleep(500 * time.Millisecond) // Simulate work
				results <- job * 2                 // Send result
			}
		}(w)
	}

	// Send jobs
	go func() {
		for j := 1; j <= numJobs; j++ {
			jobs <- j
		}
		close(jobs)
	}()

	// Close results channel when all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	for result := range results {
		fmt.Printf("Result: %d\n", result)
	}

	fmt.Println("Worker pool example completed!\n")
}

func main() {
	fmt.Println("Go Concurrency Examples: Synchronization Primitives")
	fmt.Println("===================================================")

	channelSyncExample()
	selectStatementExample()
	syncOnceExample()
	syncPoolExample()
	contextExample()
	pipelineExample()
	workerPoolExample()

	fmt.Println("Synchronization examples completed!")
}
