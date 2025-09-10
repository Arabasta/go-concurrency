package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Mutexes provide mutual exclusion to protect shared resources from concurrent access.
// Go provides sync.Mutex for exclusive access and sync.RWMutex for read-write scenarios.

// Example 1: Basic Mutex usage - protecting a counter
type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func basicMutexExample() {
	fmt.Println("=== Basic Mutex Example ===")

	counter := &Counter{}
	var wg sync.WaitGroup

	// Start 10 goroutines that increment the counter
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < 1000; j++ {
				counter.Increment()
			}
			fmt.Printf("Goroutine %d finished\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Printf("Final counter value: %d (should be 10000)\n\n", counter.Value())
}

// Example 2: RWMutex for read-heavy workloads
type SafeMap struct {
	mu   sync.RWMutex
	data map[string]int
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		data: make(map[string]int),
	}
}

func (sm *SafeMap) Set(key string, value int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[key] = value
}

func (sm *SafeMap) Get(key string) (int, bool) {
	sm.mu.RLock() // Read lock allows multiple concurrent readers
	defer sm.mu.RUnlock()
	value, ok := sm.data[key]
	return value, ok
}

func (sm *SafeMap) Keys() []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	keys := make([]string, 0, len(sm.data))
	for k := range sm.data {
		keys = append(keys, k)
	}
	return keys
}

func rwMutexExample() {
	fmt.Println("=== RWMutex Example ===")

	safeMap := NewSafeMap()
	var wg sync.WaitGroup

	// Writer goroutines (fewer, slower)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(writerID int) {
			defer wg.Done()

			for j := 0; j < 5; j++ {
				key := fmt.Sprintf("writer%d_key%d", writerID, j)
				safeMap.Set(key, writerID*10+j)
				fmt.Printf("Writer %d set %s\n", writerID, key)
				time.Sleep(100 * time.Millisecond)
			}
		}(i)
	}

	// Reader goroutines (many, faster)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(readerID int) {
			defer wg.Done()

			for j := 0; j < 10; j++ {
				keys := safeMap.Keys()
				if len(keys) > 0 {
					// Read a random key
					key := keys[rand.Intn(len(keys))]
					if value, ok := safeMap.Get(key); ok {
						fmt.Printf("Reader %d read %s: %d\n", readerID, key, value)
					}
				}
				time.Sleep(50 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("RWMutex example completed!\n")
}

// Example 3: Mutex with condition variables (using sync.Cond)
type Buffer struct {
	mu       sync.Mutex
	items    []int
	maxSize  int
	notEmpty *sync.Cond
	notFull  *sync.Cond
}

func NewBuffer(maxSize int) *Buffer {
	b := &Buffer{
		items:   make([]int, 0, maxSize),
		maxSize: maxSize,
	}
	b.notEmpty = sync.NewCond(&b.mu)
	b.notFull = sync.NewCond(&b.mu)
	return b
}

func (b *Buffer) Put(item int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Wait until buffer is not full
	for len(b.items) == b.maxSize {
		b.notFull.Wait()
	}

	b.items = append(b.items, item)
	fmt.Printf("Put item %d, buffer size: %d\n", item, len(b.items))

	// Signal that buffer is not empty
	b.notEmpty.Signal()
}

func (b *Buffer) Get() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Wait until buffer is not empty
	for len(b.items) == 0 {
		b.notEmpty.Wait()
	}

	item := b.items[0]
	b.items = b.items[1:]
	fmt.Printf("Got item %d, buffer size: %d\n", item, len(b.items))

	// Signal that buffer is not full
	b.notFull.Signal()

	return item
}

func conditionVariableExample() {
	fmt.Println("=== Condition Variable Example ===")

	buffer := NewBuffer(3)
	var wg sync.WaitGroup

	// Producer goroutines
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(producerID int) {
			defer wg.Done()

			for j := 0; j < 5; j++ {
				item := producerID*10 + j
				buffer.Put(item)
				time.Sleep(200 * time.Millisecond)
			}
		}(i)
	}

	// Consumer goroutines
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(consumerID int) {
			defer wg.Done()

			for j := 0; j < 5; j++ {
				item := buffer.Get()
				fmt.Printf("Consumer %d consumed: %d\n", consumerID, item)
				time.Sleep(300 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("Condition variable example completed!\n")
}

// Example 4: Avoiding deadlocks with ordered locking
type Account struct {
	id      int
	mu      sync.Mutex
	balance int
}

func (a *Account) Balance() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.balance
}

// Transfer money between accounts with ordered locking to avoid deadlocks
func Transfer(from, to *Account, amount int) bool {
	// Always lock accounts in order of their IDs to prevent deadlock
	first, second := from, to
	if from.id > to.id {
		first, second = to, from
	}

	first.mu.Lock()
	defer first.mu.Unlock()

	second.mu.Lock()
	defer second.mu.Unlock()

	if from.balance >= amount {
		from.balance -= amount
		to.balance += amount
		fmt.Printf("Transferred %d from account %d to account %d\n", amount, from.id, to.id)
		return true
	}

	fmt.Printf("Transfer failed: insufficient funds in account %d\n", from.id)
	return false
}

func deadlockAvoidanceExample() {
	fmt.Println("=== Deadlock Avoidance Example ===")

	account1 := &Account{id: 1, balance: 1000}
	account2 := &Account{id: 2, balance: 1000}
	account3 := &Account{id: 3, balance: 1000}

	var wg sync.WaitGroup

	// Multiple concurrent transfers
	transfers := []struct {
		from   *Account
		to     *Account
		amount int
	}{
		{account1, account2, 100},
		{account2, account3, 150},
		{account3, account1, 200},
		{account2, account1, 50},
		{account1, account3, 75},
		{account3, account2, 125},
	}

	for _, t := range transfers {
		wg.Add(1)
		go func(from, to *Account, amount int) {
			defer wg.Done()
			Transfer(from, to, amount)
			time.Sleep(100 * time.Millisecond)
		}(t.from, t.to, t.amount)
	}

	wg.Wait()

	fmt.Printf("Final balances: Account1=%d, Account2=%d, Account3=%d\n",
		account1.Balance(), account2.Balance(), account3.Balance())
	fmt.Println("Deadlock avoidance example completed!\n")
}

func main() {
	fmt.Println("Go Concurrency Examples: Mutexes")
	fmt.Println("=================================")

	rand.Seed(time.Now().UnixNano())

	basicMutexExample()
	rwMutexExample()
	conditionVariableExample()
	deadlockAvoidanceExample()

	fmt.Println("Mutex examples completed!")
}
