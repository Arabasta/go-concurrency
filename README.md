# Go Concurrency Examples
todo: 
A comprehensive collection of Go concurrency patterns, notes, and examples demonstrating various synchronization primitives and concurrent programming techniques.

## Overview

This repository contains educational examples showcasing different aspects of Go's concurrency model, including goroutines, channels, synchronization primitives, and common concurrency patterns.

## Examples

### 1. WaitGroups (`waitgroups.go`)
Demonstrates the use of `sync.WaitGroup` for coordinating goroutines:
- Basic WaitGroup usage
- Error handling patterns with WaitGroups
- Bounded concurrency with semaphores
- Producer-consumer patterns with nested WaitGroups

```bash
go run waitgroups.go
```

### 2. Mutexes (`mutex.go`)
Shows mutual exclusion patterns using `sync.Mutex` and `sync.RWMutex`:
- Basic mutex for protecting shared resources
- RWMutex for read-heavy workloads
- Condition variables with `sync.Cond`
- Deadlock avoidance techniques

```bash
go run mutex.go
```

### 3. Synchronization Primitives (`sync.go`)
Covers various synchronization mechanisms beyond mutexes:
- Channel-based synchronization
- Select statements for non-blocking operations
- `sync.Once` for one-time initialization
- `sync.Pool` for object reuse
- Context for cancellation and timeouts
- Pipeline patterns
- Worker pool patterns

```bash
go run sync.go
```

### 4. Fan-in/Fan-out Patterns (`fanin_fanout.go`)
Demonstrates advanced concurrency patterns for distributed processing:
- Basic fan-out (work distribution)
- Basic fan-in (result collection)
- Combined fan-out/fan-in with result processing
- Pipeline with multiple fan-out/fan-in stages
- Load balancer pattern

```bash
go run fanin_fanout.go
```

### 5. Advanced Concurrency Patterns (`advanced.go`)
Covers production-ready concurrency patterns and resilience techniques:
- Atomic operations with `sync/atomic` for lock-free programming
- Rate limiting with token bucket algorithm
- Graceful shutdown patterns for services
- Advanced timeout handling with context
- Pub/Sub (publish-subscribe) pattern for event-driven systems
- Circuit breaker pattern for fault tolerance
- Retry patterns with exponential backoff
- Semaphore pattern for advanced resource control
- Request-Response pattern for synchronous communication

```bash
go run advanced.go
```

## Getting Started

1. Clone this repository:
```bash
git clone https://github.com/Arabasta/go-concurrency.git
cd go-concurrency
```

2. Run any example:
```bash
go run waitgroups.go
go run mutex.go
go run sync.go
go run fanin_fanout.go
go run advanced.go
```

3. Or get usage information:
```bash
go run main.go
```

## Key Concepts Covered

- **Goroutines**: Lightweight threads managed by the Go runtime
- **Channels**: Type-safe communication between goroutines
- **WaitGroups**: Synchronization primitive for waiting on goroutine completion
- **Mutexes**: Mutual exclusion locks for protecting shared resources
- **Select Statements**: Non-blocking channel operations
- **Context**: Cancellation and timeout handling
- **Fan-out/Fan-in**: Patterns for distributing and collecting work
- **Pipeline Processing**: Chaining processing stages with channels
- **Atomic Operations**: Lock-free programming with sync/atomic
- **Rate Limiting**: Token bucket algorithm for controlling request rates
- **Circuit Breaker**: Fault tolerance pattern for resilient systems
- **Pub/Sub**: Event-driven communication pattern
- **Graceful Shutdown**: Proper service termination handling
- **Retry Patterns**: Exponential backoff for handling transient failures
- **Semaphores**: Advanced resource control and limiting
- **Request-Response**: Synchronous communication over channels

## Best Practices Demonstrated

1. **Proper resource cleanup** using `defer` statements
2. **Error handling** in concurrent contexts
3. **Bounded concurrency** to prevent resource exhaustion
4. **Deadlock avoidance** through ordered locking
5. **Context-based cancellation** for long-running operations
6. **Channel direction enforcement** for type safety

## Learning Path

1. Start with `waitgroups.go` to understand basic goroutine coordination
2. Move to `mutex.go` to learn about protecting shared resources
3. Explore `sync.go` for advanced synchronization patterns
4. Study `fanin_fanout.go` for complex distributed processing patterns
5. Master `advanced.go` for production-ready resilience and performance patterns

Each file is self-contained and can be run independently to observe the behavior of different concurrency patterns.

## Requirements

- Go 1.16 or later

## License

This project is open source and available under the MIT License.
