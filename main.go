package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Go Concurrency Examples")
		fmt.Println("=======================")
		fmt.Println()
		fmt.Println("Usage: go run <example>.go")
		fmt.Println()
		fmt.Println("Available examples:")
		fmt.Println("  waitgroups.go    - WaitGroup synchronization examples")
		fmt.Println("  mutex.go         - Mutex and RWMutex examples")
		fmt.Println("  sync.go          - General synchronization primitives")
		fmt.Println("  fanin_fanout.go  - Fan-in/Fan-out concurrency patterns")
		fmt.Println()
		fmt.Println("Example:")
		fmt.Println("  go run waitgroups.go")
		return
	}

	fmt.Println("Please run individual example files directly:")
	fmt.Println("go run", os.Args[1])
}
