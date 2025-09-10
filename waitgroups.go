package main

import (
	"fmt"
	"sync"
	"time"
)

// WaitGroups are used to wait for a collection of goroutines to finish executing.
// They provide a way to coordinate and synchronize goroutines without using channels.

// Example 1: Basic WaitGroup usage
func basicWaitGroupExample() {
	fmt.Println("=== Basic WaitGroup Example ===")

	var wg sync.WaitGroup

	// Start 3 goroutines
	for i := 1; i <= 3; i++ {
		wg.Add(1) // Increment the WaitGroup counter

		go func(id int) {
			defer wg.Done() // Decrement the counter when goroutine completes

			fmt.Printf("Worker %d starting\n", id)
			time.Sleep(time.Second) // Simulate work
			fmt.Printf("Worker %d done\n", id)
		}(i)
	}

	wg.Wait() // Block until all goroutines complete
	fmt.Println("All workers completed!\n")
}

// Example 2: WaitGroup with error handling pattern
func waitGroupWithErrorHandling() {
	fmt.Println("=== WaitGroup with Error Handling ===")

	var wg sync.WaitGroup
	errChan := make(chan error, 3) // Buffered channel for errors

	tasks := []string{"task1", "task2", "task3"}

	for _, task := range tasks {
		wg.Add(1)

		go func(taskName string) {
			defer wg.Done()

			fmt.Printf("Processing %s\n", taskName)

			// Simulate potential error
			if taskName == "task2" {
				errChan <- fmt.Errorf("error in %s", taskName)
				return
			}

			time.Sleep(500 * time.Millisecond)
			fmt.Printf("Completed %s\n", taskName)
		}(task)
	}

	// Close error channel when all goroutines complete
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Collect any errors
	for err := range errChan {
		fmt.Printf("Error occurred: %v\n", err)
	}

	fmt.Println("All tasks processed (with error handling)\n")
}

// Example 3: Bounded concurrency with WaitGroup
func boundedConcurrencyExample() {
	fmt.Println("=== Bounded Concurrency Example ===")

	const maxWorkers = 2
	semaphore := make(chan struct{}, maxWorkers) // Limit concurrent goroutines
	var wg sync.WaitGroup

	jobs := []int{1, 2, 3, 4, 5, 6}

	for _, job := range jobs {
		wg.Add(1)

		go func(jobID int) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }() // Release semaphore

			fmt.Printf("Worker processing job %d\n", jobID)
			time.Sleep(time.Second) // Simulate work
			fmt.Printf("Worker finished job %d\n", jobID)
		}(job)
	}

	wg.Wait()
	fmt.Println("All jobs completed with bounded concurrency!\n")
}

// Example 4: Nested WaitGroups (producer-consumer pattern)
func nestedWaitGroupsExample() {
	fmt.Println("=== Nested WaitGroups Example ===")

	var producerWG sync.WaitGroup
	var consumerWG sync.WaitGroup

	dataChan := make(chan int, 10)

	// Start producers
	producerWG.Add(2)
	for i := 1; i <= 2; i++ {
		go func(producerID int) {
			defer producerWG.Done()

			for j := 1; j <= 3; j++ {
				value := producerID*10 + j
				dataChan <- value
				fmt.Printf("Producer %d produced: %d\n", producerID, value)
				time.Sleep(200 * time.Millisecond)
			}
		}(i)
	}

	// Start consumers
	consumerWG.Add(2)
	for i := 1; i <= 2; i++ {
		go func(consumerID int) {
			defer consumerWG.Done()

			for data := range dataChan {
				fmt.Printf("Consumer %d consumed: %d\n", consumerID, data)
				time.Sleep(300 * time.Millisecond)
			}
		}(i)
	}

	// Close channel when all producers are done
	go func() {
		producerWG.Wait()
		close(dataChan)
	}()

	consumerWG.Wait()
	fmt.Println("Producer-consumer example completed!\n")
}

func main() {
	fmt.Println("Go Concurrency Examples: WaitGroups")
	fmt.Println("====================================")

	basicWaitGroupExample()
	waitGroupWithErrorHandling()
	boundedConcurrencyExample()
	nestedWaitGroupsExample()

	fmt.Println("WaitGroup examples completed!")
}
