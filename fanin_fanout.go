package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Fan-in and Fan-out are important concurrency patterns in Go:
// - Fan-out: Distribute work across multiple goroutines
// - Fan-in: Collect results from multiple goroutines into a single channel

// Example 1: Basic Fan-out pattern
func basicFanOut() {
	fmt.Println("=== Basic Fan-out Example ===")

	// Input data
	jobs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Channel for jobs
	jobChan := make(chan int, len(jobs))

	// Send jobs to channel
	for _, job := range jobs {
		jobChan <- job
	}
	close(jobChan)

	// Fan-out: Start multiple workers
	const numWorkers = 3
	var wg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for job := range jobChan {
				// Process job
				result := job * job
				fmt.Printf("Worker %d processed job %d -> %d\n", workerID, job, result)
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("Fan-out example completed!\n")
}

// Example 2: Basic Fan-in pattern
func basicFanIn() {
	fmt.Println("=== Basic Fan-in Example ===")

	// Create multiple producer channels
	ch1 := make(chan string)
	ch2 := make(chan string)
	ch3 := make(chan string)

	// Start producers
	go func() {
		defer close(ch1)
		for i := 1; i <= 3; i++ {
			ch1 <- fmt.Sprintf("Producer1-Item%d", i)
			time.Sleep(300 * time.Millisecond)
		}
	}()

	go func() {
		defer close(ch2)
		for i := 1; i <= 3; i++ {
			ch2 <- fmt.Sprintf("Producer2-Item%d", i)
			time.Sleep(400 * time.Millisecond)
		}
	}()

	go func() {
		defer close(ch3)
		for i := 1; i <= 3; i++ {
			ch3 <- fmt.Sprintf("Producer3-Item%d", i)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Fan-in: Merge multiple channels into one
	merged := fanIn(ch1, ch2, ch3)

	// Consume from merged channel
	for item := range merged {
		fmt.Printf("Received: %s\n", item)
	}

	fmt.Println("Fan-in example completed!\n")
}

// fanIn merges multiple channels into a single channel
func fanIn(channels ...<-chan string) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup

	// Start a goroutine for each input channel
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan string) {
			defer wg.Done()
			for item := range c {
				out <- item
			}
		}(ch)
	}

	// Close output channel when all input channels are done
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// Example 3: Fan-out/Fan-in with result collection
type Job struct {
	ID   int
	Data int
}

type Result struct {
	JobID  int
	Result int
	Worker int
}

func fanOutFanInWithResults() {
	fmt.Println("=== Fan-out/Fan-in with Results Example ===")

	// Create jobs
	jobs := make(chan Job, 10)
	for i := 1; i <= 10; i++ {
		jobs <- Job{ID: i, Data: rand.Intn(100)}
	}
	close(jobs)

	// Fan-out: Start workers
	const numWorkers = 3
	resultChannels := make([]<-chan Result, numWorkers)

	for i := 0; i < numWorkers; i++ {
		resultChannels[i] = worker(i+1, jobs)
	}

	// Fan-in: Merge results
	results := fanInResults(resultChannels...)

	// Collect and display results
	var allResults []Result
	for result := range results {
		allResults = append(allResults, result)
		fmt.Printf("Job %d processed by Worker %d: %d -> %d\n",
			result.JobID, result.Worker, result.JobID, result.Result)
	}

	fmt.Printf("Processed %d jobs total\n", len(allResults))
	fmt.Println("Fan-out/Fan-in with results completed!\n")
}

// worker processes jobs and returns results
func worker(id int, jobs <-chan Job) <-chan Result {
	results := make(chan Result)

	go func() {
		defer close(results)
		for job := range jobs {
			// Simulate processing time
			time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond)

			// Process job (square the data)
			result := Result{
				JobID:  job.ID,
				Result: job.Data * job.Data,
				Worker: id,
			}

			results <- result
		}
	}()

	return results
}

// fanInResults merges multiple result channels
func fanInResults(channels ...<-chan Result) <-chan Result {
	out := make(chan Result)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan Result) {
			defer wg.Done()
			for result := range c {
				out <- result
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// Example 4: Pipeline with Fan-out/Fan-in stages
func pipelineWithFanOutFanIn() {
	fmt.Println("=== Pipeline with Fan-out/Fan-in Example ===")

	// Stage 1: Generate numbers
	numbers := generateNumbers(1, 20)

	// Stage 2: Fan-out for parallel processing (multiply by 2)
	const stage2Workers = 3
	stage2Channels := make([]<-chan int, stage2Workers)
	for i := 0; i < stage2Workers; i++ {
		stage2Channels[i] = multiplyBy2(numbers)
	}

	// Fan-in stage 2 results
	stage2Results := fanInNumbers(stage2Channels...)

	// Stage 3: Fan-out for parallel processing (add 10)
	const stage3Workers = 2
	stage3Channels := make([]<-chan int, stage3Workers)
	for i := 0; i < stage3Workers; i++ {
		stage3Channels[i] = add10(stage2Results)
	}

	// Fan-in final results
	finalResults := fanInNumbers(stage3Channels...)

	// Collect results
	var results []int
	for result := range finalResults {
		results = append(results, result)
	}

	fmt.Printf("Pipeline processed %d numbers\n", len(results))
	fmt.Printf("Sample results: %v\n", results[:min(5, len(results))])
	fmt.Println("Pipeline with fan-out/fan-in completed!\n")
}

func generateNumbers(start, count int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := start; i < start+count; i++ {
			out <- i
			time.Sleep(10 * time.Millisecond)
		}
	}()
	return out
}

func multiplyBy2(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * 2
			time.Sleep(50 * time.Millisecond) // Simulate processing
		}
	}()
	return out
}

func add10(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n + 10
			time.Sleep(30 * time.Millisecond) // Simulate processing
		}
	}()
	return out
}

func fanInNumbers(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for n := range c {
				out <- n
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// Example 5: Load balancer pattern (round-robin fan-out)
func loadBalancerExample() {
	fmt.Println("=== Load Balancer Example ===")

	// Create worker channels
	const numWorkers = 3
	workers := make([]chan string, numWorkers)
	var wg sync.WaitGroup

	for i := range workers {
		workers[i] = make(chan string, 1) // Buffered to prevent blocking
		wg.Add(1)

		// Start worker
		go func(id int, jobs <-chan string) {
			defer wg.Done()
			for job := range jobs {
				fmt.Printf("Worker %d processing: %s\n", id+1, job)
				time.Sleep(time.Duration(rand.Intn(200)+100) * time.Millisecond)
			}
		}(i, workers[i])
	}

	// Load balancer function - simple round-robin
	loadBalancer := func(jobs []string, workers []chan string) {
		for i, job := range jobs {
			workerIndex := i % len(workers)
			workers[workerIndex] <- job
		}
	}

	// Create jobs
	jobs := []string{}
	for i := 1; i <= 10; i++ {
		jobs = append(jobs, fmt.Sprintf("Job-%d", i))
	}

	// Start load balancer
	loadBalancer(jobs, workers)

	// Close worker channels
	for _, worker := range workers {
		close(worker)
	}

	wg.Wait() // Wait for all workers to finish
	fmt.Println("Load balancer example completed!\n")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	fmt.Println("Go Concurrency Examples: Fan-in/Fan-out Patterns")
	fmt.Println("=================================================")

	rand.Seed(time.Now().UnixNano())

	basicFanOut()
	basicFanIn()
	fanOutFanInWithResults()
	pipelineWithFanOutFanIn()
	loadBalancerExample()

	fmt.Println("Fan-in/Fan-out examples completed!")
}
