package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// Advanced concurrency patterns in Go including atomic operations,
// rate limiting, graceful shutdown, pub/sub, and resilience patterns.

// Example 1: Atomic Operations - Lock-free counter
type AtomicCounter struct {
	value int64
}

func (c *AtomicCounter) Increment() {
	atomic.AddInt64(&c.value, 1)
}

func (c *AtomicCounter) Decrement() {
	atomic.AddInt64(&c.value, -1)
}

func (c *AtomicCounter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

func (c *AtomicCounter) CompareAndSwap(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&c.value, old, new)
}

func atomicOperationsExample() {
	fmt.Println("=== Atomic Operations Example ===")

	counter := &AtomicCounter{}
	var wg sync.WaitGroup
	
	// Start multiple goroutines incrementing the counter
	const numGoroutines = 100
	const incrementsPerGoroutine = 1000

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			for j := 0; j < incrementsPerGoroutine; j++ {
				counter.Increment()
			}
			
			if id == 0 {
				fmt.Printf("Goroutine %d: Current value = %d\n", id, counter.Value())
			}
		}(i)
	}

	wg.Wait()
	
	expectedValue := int64(numGoroutines * incrementsPerGoroutine)
	actualValue := counter.Value()
	
	fmt.Printf("Expected: %d, Actual: %d, Match: %t\n", expectedValue, actualValue, expectedValue == actualValue)
	
	// Demonstrate compare-and-swap
	fmt.Println("Testing compare-and-swap...")
	oldVal := actualValue
	newVal := oldVal + 42
	success := counter.CompareAndSwap(oldVal, newVal)
	fmt.Printf("CAS(%d -> %d): %t, Final value: %d\n", oldVal, newVal, success, counter.Value())
	fmt.Println()
}

// Example 2: Token Bucket Rate Limiter
type TokenBucket struct {
	tokens    int64
	capacity  int64
	refillRate int64
	lastRefill int64
	mu        sync.Mutex
}

func NewTokenBucket(capacity, refillRate int64) *TokenBucket {
	return &TokenBucket{
		tokens:     capacity,
		capacity:   capacity,
		refillRate: refillRate,
		lastRefill: time.Now().UnixNano(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	
	now := time.Now().UnixNano()
	elapsed := now - tb.lastRefill
	tokensToAdd := (elapsed / int64(time.Second)) * tb.refillRate
	
	tb.tokens = min(tb.capacity, tb.tokens+tokensToAdd)
	tb.lastRefill = now
	
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	
	return false
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func rateLimiterExample() {
	fmt.Println("=== Rate Limiter (Token Bucket) Example ===")
	
	// Create a rate limiter: 5 tokens capacity, refill 2 tokens per second
	limiter := NewTokenBucket(5, 2)
	
	// Simulate requests
	for i := 0; i < 20; i++ {
		if limiter.Allow() {
			fmt.Printf("Request %d: ALLOWED\n", i+1)
		} else {
			fmt.Printf("Request %d: RATE LIMITED\n", i+1)
		}
		
		// Small delay between requests
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Println()
}

// Example 3: Graceful Shutdown Pattern
type Server struct {
	name     string
	shutdown chan struct{}
	done     chan struct{}
}

func NewServer(name string) *Server {
	return &Server{
		name:     name,
		shutdown: make(chan struct{}),
		done:     make(chan struct{}),
	}
}

func (s *Server) Start() {
	fmt.Printf("Server %s starting...\n", s.name)
	
	go func() {
		defer close(s.done)
		
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		
		for {
			select {
			case <-s.shutdown:
				fmt.Printf("Server %s received shutdown signal\n", s.name)
				return
			case <-ticker.C:
				fmt.Printf("Server %s: processing...\n", s.name)
			}
		}
	}()
}

func (s *Server) Stop() {
	fmt.Printf("Server %s stopping...\n", s.name)
	close(s.shutdown)
	<-s.done
	fmt.Printf("Server %s stopped gracefully\n", s.name)
}

func gracefulShutdownExample() {
	fmt.Println("=== Graceful Shutdown Pattern Example ===")
	
	// Start multiple servers
	servers := []*Server{
		NewServer("API"),
		NewServer("Database"),
		NewServer("Cache"),
	}
	
	// Start all servers
	for _, server := range servers {
		server.Start()
	}
	
	// Let them run for a while
	time.Sleep(2 * time.Second)
	
	// Graceful shutdown all servers
	var wg sync.WaitGroup
	for _, server := range servers {
		wg.Add(1)
		go func(s *Server) {
			defer wg.Done()
			s.Stop()
		}(server)
	}
	
	wg.Wait()
	fmt.Println("All servers shut down gracefully\n")
}

// Example 4: Advanced Timeout Patterns
func advancedTimeoutExample() {
	fmt.Println("=== Advanced Timeout Patterns Example ===")
	
	// Pattern 1: Context with timeout
	fmt.Println("--- Context Timeout Pattern ---")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	
	result := make(chan string, 1)
	go func() {
		// Simulate work that might take a long time
		time.Sleep(1500 * time.Millisecond)
		result <- "work completed"
	}()
	
	select {
	case res := <-result:
		fmt.Printf("Result: %s\n", res)
	case <-ctx.Done():
		fmt.Printf("Operation timed out: %v\n", ctx.Err())
	}
	
	// Pattern 2: Multiple operations with timeout
	fmt.Println("--- Multiple Operations Timeout ---")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel2()
	
	operations := []string{"op1", "op2", "op3"}
	results := make(chan string, len(operations))
	
	for _, op := range operations {
		go func(operation string) {
			// Simulate varying operation durations
			duration := time.Duration(rand.Intn(1000)+500) * time.Millisecond
			
			select {
			case <-time.After(duration):
				results <- fmt.Sprintf("%s completed in %v", operation, duration)
			case <-ctx2.Done():
				results <- fmt.Sprintf("%s cancelled: %v", operation, ctx2.Err())
				return
			}
		}(op)
	}
	
	for i := 0; i < len(operations); i++ {
		select {
		case result := <-results:
			fmt.Printf("  %s\n", result)
		case <-ctx2.Done():
			fmt.Printf("  Timeout waiting for results: %v\n", ctx2.Err())
			break
		}
	}
	fmt.Println()
}

// Example 5: Pub/Sub (Publish-Subscribe) Pattern
type PubSub struct {
	mu          sync.RWMutex
	subscribers map[string][]chan string
	closed      bool
}

func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: make(map[string][]chan string),
	}
}

func (ps *PubSub) Subscribe(topic string) <-chan string {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	
	if ps.closed {
		return nil
	}
	
	ch := make(chan string, 10) // Buffered to prevent blocking
	ps.subscribers[topic] = append(ps.subscribers[topic], ch)
	return ch
}

func (ps *PubSub) Publish(topic, message string) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	
	if ps.closed {
		return
	}
	
	for _, ch := range ps.subscribers[topic] {
		select {
		case ch <- message:
		default:
			// Skip slow subscribers
		}
	}
}

func (ps *PubSub) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	
	if ps.closed {
		return
	}
	
	ps.closed = true
	for _, channels := range ps.subscribers {
		for _, ch := range channels {
			close(ch)
		}
	}
	ps.subscribers = nil
}

func pubSubExample() {
	fmt.Println("=== Pub/Sub Pattern Example ===")
	
	pubsub := NewPubSub()
	defer pubsub.Close()
	
	// Create subscribers for different topics
	newsSub := pubsub.Subscribe("news")
	sportsSub := pubsub.Subscribe("sports")
	newsSub2 := pubsub.Subscribe("news") // Multiple subscribers to same topic
	
	var wg sync.WaitGroup
	
	// Start subscribers
	wg.Add(3)
	
	// News subscriber 1
	go func() {
		defer wg.Done()
		for msg := range newsSub {
			fmt.Printf("[News Subscriber 1] Received: %s\n", msg)
		}
	}()
	
	// Sports subscriber
	go func() {
		defer wg.Done()
		for msg := range sportsSub {
			fmt.Printf("[Sports Subscriber] Received: %s\n", msg)
		}
	}()
	
	// News subscriber 2
	go func() {
		defer wg.Done()
		for msg := range newsSub2 {
			fmt.Printf("[News Subscriber 2] Received: %s\n", msg)
		}
	}()
	
	// Publisher
	go func() {
		messages := []struct {
			topic string
			msg   string
		}{
			{"news", "Breaking: Important news update"},
			{"sports", "Football match results"},
			{"news", "Weather forecast update"},
			{"sports", "Basketball championship"},
			{"news", "Technology breakthrough"},
		}
		
		for _, m := range messages {
			pubsub.Publish(m.topic, m.msg)
			time.Sleep(300 * time.Millisecond)
		}
	}()
	
	// Let it run for a while then close
	time.Sleep(2 * time.Second)
	pubsub.Close()
	wg.Wait()
	fmt.Println()
}

// Example 6: Circuit Breaker Pattern
type CircuitBreaker struct {
	mu            sync.Mutex
	maxFailures   int
	resetTimeout  time.Duration
	failures      int
	lastFailTime  time.Time
	state         string // "CLOSED", "OPEN", "HALF_OPEN"
}

func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        "CLOSED",
	}
}

func (cb *CircuitBreaker) Call(fn func() error) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	// Check if we should move from OPEN to HALF_OPEN
	if cb.state == "OPEN" {
		if time.Since(cb.lastFailTime) > cb.resetTimeout {
			cb.state = "HALF_OPEN"
			cb.failures = 0
		} else {
			return fmt.Errorf("circuit breaker is OPEN")
		}
	}
	
	// Execute the function
	err := fn()
	
	if err != nil {
		cb.failures++
		cb.lastFailTime = time.Now()
		
		if cb.failures >= cb.maxFailures {
			cb.state = "OPEN"
		}
		return err
	}
	
	// Success - reset if we were in HALF_OPEN state
	if cb.state == "HALF_OPEN" {
		cb.state = "CLOSED"
	}
	cb.failures = 0
	
	return nil
}

func (cb *CircuitBreaker) State() string {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

func circuitBreakerExample() {
	fmt.Println("=== Circuit Breaker Pattern Example ===")
	
	cb := NewCircuitBreaker(3, 2*time.Second)
	
	// Simulate a flaky service
	callCount := 0
	flakyService := func() error {
		callCount++
		// Fail for the first 5 calls, then succeed
		if callCount <= 5 {
			return fmt.Errorf("service unavailable")
		}
		return nil
	}
	
	// Test the circuit breaker
	for i := 1; i <= 10; i++ {
		err := cb.Call(flakyService)
		state := cb.State()
		
		if err != nil {
			fmt.Printf("Call %d: ERROR (%v) - State: %s\n", i, err, state)
		} else {
			fmt.Printf("Call %d: SUCCESS - State: %s\n", i, state)
		}
		
		time.Sleep(300 * time.Millisecond)
		
		// After circuit opens, wait for reset
		if i == 6 && state == "OPEN" {
			fmt.Println("Waiting for circuit breaker reset...")
			time.Sleep(2 * time.Second)
		}
	}
	fmt.Println()
}

// Example 7: Retry Pattern with Exponential Backoff
type RetryConfig struct {
	MaxRetries  int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	Multiplier  float64
}

func RetryWithBackoff(config RetryConfig, fn func() error) error {
	var lastErr error
	
	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		if attempt > 0 {
			delay := time.Duration(float64(config.BaseDelay) * 
				   float64(attempt) * config.Multiplier)
			if delay > config.MaxDelay {
				delay = config.MaxDelay
			}
			
			fmt.Printf("  Retry attempt %d after %v delay\n", attempt, delay)
			time.Sleep(delay)
		}
		
		err := fn()
		if err == nil {
			if attempt > 0 {
				fmt.Printf("  Success after %d retries\n", attempt)
			}
			return nil
		}
		
		lastErr = err
		fmt.Printf("  Attempt %d failed: %v\n", attempt+1, err)
	}
	
	return fmt.Errorf("failed after %d attempts: %v", config.MaxRetries+1, lastErr)
}

func retryPatternExample() {
	fmt.Println("=== Retry Pattern with Exponential Backoff Example ===")
	
	config := RetryConfig{
		MaxRetries: 4,
		BaseDelay:  100 * time.Millisecond,
		MaxDelay:   2 * time.Second,
		Multiplier: 2.0,
	}
	
	// Simulate a service that fails a few times then succeeds
	attempts := 0
	unreliableService := func() error {
		attempts++
		if attempts <= 3 {
			return fmt.Errorf("temporary failure %d", attempts)
		}
		return nil
	}
	
	fmt.Println("Testing retry pattern:")
	err := RetryWithBackoff(config, unreliableService)
	
	if err != nil {
		fmt.Printf("Final result: FAILED - %v\n", err)
	} else {
		fmt.Printf("Final result: SUCCESS\n")
	}
	fmt.Println()
}

// Example 8: Semaphore Pattern for Resource Control
type Semaphore struct {
	permits chan struct{}
}

func NewSemaphore(maxConcurrency int) *Semaphore {
	return &Semaphore{
		permits: make(chan struct{}, maxConcurrency),
	}
}

func (s *Semaphore) Acquire() {
	s.permits <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.permits
}

func (s *Semaphore) TryAcquire() bool {
	select {
	case s.permits <- struct{}{}:
		return true
	default:
		return false
	}
}

func semaphoreExample() {
	fmt.Println("=== Semaphore Pattern Example ===")
	
	// Create a semaphore that allows maximum 3 concurrent operations
	sem := NewSemaphore(3)
	
	var wg sync.WaitGroup
	
	// Start 10 goroutines, but only 3 can run concurrently
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			fmt.Printf("Task %d waiting for permit...\n", id)
			sem.Acquire()
			fmt.Printf("Task %d acquired permit, working...\n", id)
			
			// Simulate work
			time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond)
			
			fmt.Printf("Task %d completed, releasing permit\n", id)
			sem.Release()
		}(i)
	}
	
	wg.Wait()
	fmt.Println("All tasks completed")
	
	// Demonstrate TryAcquire
	fmt.Println("\nTesting non-blocking acquire:")
	for i := 1; i <= 5; i++ {
		if sem.TryAcquire() {
			fmt.Printf("  Successfully acquired permit %d\n", i)
		} else {
			fmt.Printf("  Failed to acquire permit %d (would block)\n", i)
		}
	}
	
	// Release the permits we acquired
	for i := 1; i <= 3; i++ {
		sem.Release()
	}
	fmt.Println()
}

// Example 9: Request-Response Pattern
type Request struct {
	ID       int
	Data     string
	Response chan string
}

type RequestHandler struct {
	requests chan Request
	workers  int
}

func NewRequestHandler(workers int) *RequestHandler {
	rh := &RequestHandler{
		requests: make(chan Request, 100),
		workers:  workers,
	}
	
	// Start worker goroutines
	for i := 0; i < workers; i++ {
		go rh.worker(i)
	}
	
	return rh
}

func (rh *RequestHandler) worker(id int) {
	for req := range rh.requests {
		// Process the request
		result := fmt.Sprintf("Worker %d processed request %d: %s -> PROCESSED", 
			id, req.ID, req.Data)
		
		// Simulate processing time
		time.Sleep(time.Duration(rand.Intn(200)+100) * time.Millisecond)
		
		// Send response back
		req.Response <- result
		close(req.Response)
	}
}

func (rh *RequestHandler) HandleRequest(id int, data string) string {
	responseChan := make(chan string, 1)
	
	req := Request{
		ID:       id,
		Data:     data,
		Response: responseChan,
	}
	
	// Send request (non-blocking due to buffered channel)
	select {
	case rh.requests <- req:
		// Wait for response
		return <-responseChan
	default:
		return fmt.Sprintf("Request %d rejected: handler busy", id)
	}
}

func (rh *RequestHandler) Close() {
	close(rh.requests)
}

func requestResponseExample() {
	fmt.Println("=== Request-Response Pattern Example ===")
	
	handler := NewRequestHandler(3) // 3 workers
	defer handler.Close()
	
	var wg sync.WaitGroup
	
	// Send multiple requests concurrently
	for i := 1; i <= 8; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			data := fmt.Sprintf("data_%d", id)
			response := handler.HandleRequest(id, data)
			fmt.Printf("Request %d response: %s\n", id, response)
		}(i)
	}
	
	wg.Wait()
	fmt.Println()
}

func main() {
	fmt.Println("Go Advanced Concurrency Patterns")
	fmt.Println("================================")
	fmt.Println()

	atomicOperationsExample()
	rateLimiterExample()
	gracefulShutdownExample()
	advancedTimeoutExample()
	pubSubExample()
	circuitBreakerExample()
	retryPatternExample()
	semaphoreExample()
	requestResponseExample()

	fmt.Println("All advanced concurrency examples completed!")
}