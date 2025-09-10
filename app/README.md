# Distributed Systems Applications

[![Go Version](https://img.shields.io/badge/Go-1.25.1-blue.svg)](https://golang.org)
[![Status](https://img.shields.io/badge/Status-Active%20Development-green.svg)](https://github.com/yourusername/Distributed-Systems)

> **Master distributed systems through hands-on implementation of real-world algorithms and architectures**

## Why Distributed Systems Matter

Modern computing is inherently distributed. From web applications serving millions of users to cloud infrastructure spanning continents, distributed systems power the digital world. However, building reliable distributed systems presents unique challenges that don't exist in single-machine applications:

- **Partial Failures**: Networks fail, servers crash, and disks corrupt—but the system must continue operating
- **Concurrency**: Multiple processes execute simultaneously, creating complex interaction patterns
- **Consistency**: Ensuring data remains correct across multiple replicas under concurrent access
- **Scalability**: Growing from single machines to clusters of thousands of nodes
- **Latency**: Network communication introduces delays that affect system behavior

This repository provides a structured path to master these challenges through practical implementation. Rather than just reading about distributed algorithms, you'll build them, test them under failure conditions, and understand their real-world implications.

## Learning Philosophy

**Theory without practice is incomplete.** These labs bridge the gap between academic understanding and production engineering by:

1. **Implementing Core Algorithms**: Build Paxos consensus, MapReduce processing, and sharded storage from scratch
2. **Testing Under Failure**: Simulate network partitions, server crashes, and message reordering
3. **Measuring Performance**: Profile bottlenecks, optimize throughput, and understand trade-offs
4. **Engineering for Production**: Handle edge cases, implement proper error handling, and design for observability

Each lab builds upon previous concepts, creating a comprehensive understanding of how distributed systems work in practice.

## Quick Start

### Prerequisites
- **Go 1.25.1+** [[memory:8397771]] - Required for all examples and tests
- **4GB+ RAM** (8GB recommended for large datasets and concurrent testing)
- **Basic Go knowledge** - Goroutines, channels, interfaces, and error handling
- **Understanding of concurrency** - Mutexes, atomic operations, and race conditions
- **Network fundamentals** - TCP/IP, client-server communication, and RPC concepts

### System Requirements
- **Operating System**: Linux, macOS, or Windows with WSL2
- **CPU**: Multi-core processor recommended for concurrent testing
- **Storage**: 2GB free space for source code and test data
- **Network**: Internet connection for downloading dependencies

### 30-Second Demo
```bash
# 1. Verify Go installation and environment
go version  # Should show go1.25.1 or later
go env GOPATH GOROOT  # Check Go environment

# 2. Run MapReduce word count on King James Bible
cd app/01_Practice-Labs/src/main
go run wc.go master kjv12.txt sequential
# This processes ~4MB of text using distributed computing principles

# 3. Test sharded key-value store with fault tolerance
cd ../shardmaster && go test -v
cd ../shardkv && go test -race -v
# These tests verify consensus, replication, and sharding
```

### What You'll See
The demo showcases three fundamental distributed systems concepts:
1. **MapReduce**: Parallel data processing across multiple workers
2. **Consensus**: Agreement protocols for distributed decision-making  
3. **Sharding**: Horizontal partitioning for scalability

## Learning Path

### Phase 1: Fundamentals (40-50 hours)
Start with the **Practice Labs** for guided, hands-on learning. Each lab builds essential distributed systems skills:

| Lab | Topic | Duration | Complexity | Key Skills Developed |
|-----|-------|----------|------------|---------------------|
| **Lab 1** | MapReduce | 4-6 hours | Beginner | Distributed computing, fault tolerance, task scheduling |
| **Lab 2** | Primary-Backup | 6-8 hours | Intermediate | Replication, failure detection, state synchronization |
| **Lab 3** | Paxos Consensus | 8-12 hours | Advanced | Consensus algorithms, leader election, log replication |
| **Lab 4** | Sharded KV Store | 10-15 hours | Expert | Sharding, load balancing, dynamic reconfiguration |
| **Lab 5** | Persistent Storage | 6-10 hours | Advanced | Durability, crash recovery, performance optimization |

### Learning Progression Strategy

**Week 1-2: Foundation Building**
- Complete Labs 1-2 to understand basic distributed computing patterns
- Focus on understanding failure modes and recovery mechanisms
- Practice with Go concurrency primitives (goroutines, channels, mutexes)

**Week 3-4: Consensus and Agreement**
- Tackle Lab 3 to master consensus algorithms
- Study the theoretical foundations in `../reference/`
- Implement and test Paxos under various failure scenarios

**Week 5-6: Scalability and Persistence**
- Complete Labs 4-5 to understand large-scale system design
- Learn about sharding strategies and performance optimization
- Practice with real-world data persistence challenges

### Phase 2: Advanced Projects (Self-Directed)
Extend your knowledge with production-ready implementations:

- **RPC Framework** — Build reliable remote procedure calls with retries, timeouts, and circuit breakers
- **Consensus Library** — Implement Raft or PBFT with proper leader election and log compaction
- **Distributed Cache** — Create a fault-tolerant caching layer with consistent hashing
- **Service Mesh** — Build service discovery, load balancing, and traffic management
- **Observability Stack** — Add comprehensive metrics, structured logging, and distributed tracing
- **Distributed Database** — Implement ACID transactions across multiple shards
- **Message Queue** — Build reliable message delivery with at-least-once semantics

## Project Structure

```
app/
├── 01_Practice-Labs/          # Start here - Guided exercises
│   ├── lab1.md               # MapReduce implementation
│   ├── lab2a.md              # ViewService & failure detection
│   ├── lab2b.md              # Primary-backup replication
│   ├── lab3a.md              # Basic Paxos consensus
│   ├── lab3b.md              # Multi-Paxos & replicated logs
│   ├── lab4a.md              # ShardMaster configuration
│   ├── lab4b.md              # Sharded key-value store
│   ├── lab5.md               # Persistent storage & recovery
│   └── src/                  # Complete implementations
│       ├── mapreduce/        # Lab 1: Distributed processing
│       ├── viewservice/      # Lab 2A: Failure detection
│       ├── pbservice/        # Lab 2B: Replication
│       ├── paxos/            # Lab 3: Consensus
│       ├── shardmaster/      # Lab 4A: Configuration mgmt
│       ├── shardkv/          # Lab 4B: Sharded storage
│       ├── diskv/            # Lab 5: Persistent storage
│       └── main/             # Example applications
└── 02_Practice-Project/      # Advanced projects (coming soon)
```

## Detailed Lab Overview

### Lab 1: MapReduce - Distributed Data Processing

**Architecture**: Master-Worker pattern with fault-tolerant task scheduling

**What You'll Build**:
- A master process that coordinates work across multiple workers
- Worker processes that execute map and reduce functions
- Fault tolerance mechanisms for handling worker failures
- Task reassignment and recovery protocols

**Key Learning Outcomes**:
- Understand the master-worker architectural pattern
- Implement distributed task scheduling and coordination
- Handle partial failures in distributed computations
- Process large datasets using parallel computation

**Technical Challenges**:
- Coordinating work across multiple processes
- Detecting and recovering from worker failures
- Ensuring all tasks complete despite failures
- Optimizing data transfer between map and reduce phases

**Testing Commands**:
```bash
cd app/01_Practice-Labs/src/mapreduce
go test -v                    # Basic functionality
go test -race -v              # Concurrency testing
go test -run TestFailure -v   # Failure scenarios
```

### Lab 2: Primary-Backup Replication - Fault-Tolerant Storage

**Architecture**: Primary-backup replication with ViewService for failure detection

**What You'll Build**:
- ViewService that monitors server health and maintains cluster view
- Primary-backup key-value store with automatic failover
- State synchronization between primary and backup servers
- Client-side failover mechanisms

**Key Learning Outcomes**:
- Implement failure detection and cluster membership protocols
- Design primary-backup replication with consistency guarantees
- Handle split-brain scenarios and network partitions
- Build client libraries that handle server failures gracefully

**Technical Challenges**:
- Detecting server failures accurately without false positives
- Synchronizing state between primary and backup servers
- Handling concurrent client requests during failover
- Ensuring consistency during network partitions

**Testing Commands**:
```bash
cd app/01_Practice-Labs/src/viewservice && go test -v
cd ../pbservice && go test -race -v
go test -run TestPartition -v  # Network partition testing
```

### Lab 3: Paxos Consensus - Distributed Agreement

**Architecture**: Paxos consensus algorithm for fault-tolerant agreement

**What You'll Build**:
- Basic Paxos implementation for single-value consensus
- Multi-Paxos for replicated log consensus
- Leader election and optimization mechanisms
- Log synchronization and garbage collection

**Key Learning Outcomes**:
- Master the Paxos consensus algorithm and its variants
- Understand leader election and log replication
- Implement replicated state machines
- Handle network partitions and Byzantine failures

**Technical Challenges**:
- Implementing the complex Paxos protocol correctly
- Optimizing for common-case performance
- Handling log holes and synchronization issues
- Managing memory usage with log compaction

**Testing Commands**:
```bash
cd app/01_Practice-Labs/src/paxos
go test -race -v              # Always use race detection
go test -run TestUnreliable -v # Network failure testing
go test -run TestPartition -v  # Partition tolerance testing
```

### Lab 4: Sharded Key-Value Store - Scalable Distributed Storage

**Architecture**: Sharded storage with dynamic reconfiguration

**What You'll Build**:
- ShardMaster for configuration management and load balancing
- ShardKV servers that handle individual shards
- Dynamic shard movement and reconfiguration
- Cross-shard operations and consistency guarantees

**Key Learning Outcomes**:
- Design and implement horizontal partitioning (sharding)
- Build dynamic reconfiguration systems
- Implement load balancing and shard migration
- Handle cross-shard operations and consistency

**Technical Challenges**:
- Designing efficient sharding and rebalancing algorithms
- Ensuring consistency during shard movement
- Handling concurrent reconfigurations
- Optimizing for load distribution and performance

**Testing Commands**:
```bash
cd app/01_Practice-Labs/src/shardmaster && go test -race -v
cd ../shardkv && go test -race -v
go test -run TestConcurrent -v    # Concurrent reconfiguration
go test -run TestUnreliable -v    # Network failure testing
```

### Lab 5: Persistent Storage - Durability and Crash Recovery

**Architecture**: Persistent storage with atomic operations and crash recovery

**What You'll Build**:
- Atomic file operations for data durability
- Crash recovery procedures for system restart
- Disk layout optimization for performance
- Integration with sharded key-value store

**Key Learning Outcomes**:
- Implement persistent storage with ACID properties
- Design crash recovery and durability mechanisms
- Optimize disk I/O and storage layout
- Integrate persistence with distributed systems

**Technical Challenges**:
- Ensuring atomicity of file operations
- Designing efficient crash recovery procedures
- Optimizing disk layout for performance
- Handling concurrent access to persistent storage

**Testing Commands**:
```bash
cd app/01_Practice-Labs/src/diskv
go test -run Test4 -v         # Lab 4 compatibility
go test -race -v              # Full Lab 5 suite
go test -run TestCrash -v     # Crash recovery testing
```

## Development Workflow and Methodology

### 1. Theoretical Foundation
Before implementing any distributed system, understand the underlying theory:

```bash
# Study the concepts before coding
cd ../reference
# Read relevant sections for each lab:
# - 07-1_Paxos.md for Lab 3
# - 04-1_Primary-Backup.md for Lab 2
# - 12-1_GFS.md for Lab 1 (MapReduce)
```

**Key Reading Strategy**:
- Start with the problem statement and motivation
- Understand the algorithm's correctness properties
- Study the failure scenarios and edge cases
- Review the performance characteristics and trade-offs

### 2. Incremental Implementation Strategy

**Phase 1: Basic Functionality**
- Implement the simplest working version first
- Focus on correctness over performance
- Write comprehensive unit tests
- Use `go test -race` for every commit

**Phase 2: Feature Addition**
- Add features one at a time
- Test thoroughly after each change
- Maintain backward compatibility
- Document API changes

**Phase 3: Failure Handling**
- Add failure detection and recovery
- Implement timeout and retry mechanisms
- Test under various failure scenarios
- Measure performance impact

**Phase 4: Optimization**
- Profile for bottlenecks
- Optimize critical code paths
- Add performance benchmarks
- Document performance characteristics

### 3. Testing Methodology

**Unit Testing**:
```bash
# Test individual components
go test -v ./src/mapreduce/...
go test -race -v ./src/paxos/...
```

**Integration Testing**:
```bash
# Test component interactions
go test -run TestIntegration -v
go test -run TestEndToEnd -v
```

**Failure Testing**:
```bash
# Test under failure conditions
go test -run TestFailure -v
go test -run TestPartition -v
go test -run TestUnreliable -v
```

**Performance Testing**:
```bash
# Benchmark critical operations
go test -bench=. -benchmem
go test -bench=BenchmarkConsensus -benchmem
```

### 4. Performance Analysis and Optimization

**Profiling Workflow**:
```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof
(pprof) top10
(pprof) list functionName

# Memory profiling
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
(pprof) top10
(pprof) list functionName

# Goroutine analysis
go test -blockprofile=block.prof -bench=.
go tool pprof block.prof
```

**Optimization Strategies**:
- Identify hot code paths through profiling
- Reduce memory allocations in critical sections
- Optimize network communication patterns
- Implement efficient data structures
- Use connection pooling and caching

### 5. Code Quality and Best Practices

**Go-Specific Guidelines**:
- Follow Go naming conventions and idioms
- Use interfaces for testability and modularity
- Implement proper error handling with context
- Use channels for communication, mutexes for synchronization
- Write self-documenting code with clear variable names

**Distributed Systems Best Practices**:
- Design for failure from the beginning
- Implement idempotent operations where possible
- Use timeouts and circuit breakers
- Log important state changes and decisions
- Implement health checks and monitoring

## Comprehensive Testing Strategy

### Test Categories and Commands

**Concurrency Testing** (Critical for distributed systems):
```bash
# Always use race detection for concurrent code
go test -race ./...
go test -race -v ./src/paxos/...
go test -race -count=10 ./src/shardkv/...  # Run multiple times
```

**Coverage Analysis**:
```bash
# Generate coverage reports
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out  # View in browser
go test -covermode=atomic -race ./...  # Race-aware coverage
```

**Performance Benchmarking**:
```bash
# Benchmark critical operations
go test -bench=. -benchmem ./...
go test -bench=BenchmarkConsensus -benchmem ./src/paxos/
go test -bench=BenchmarkShardOperation -benchmem ./src/shardkv/
```

**Failure Scenario Testing**:
```bash
# Test under various failure conditions
go test -run TestFailure -v ./...
go test -run TestPartition -v ./...
go test -run TestUnreliable -v ./...
go test -run TestCrash -v ./...
```

**Load and Stress Testing**:
```bash
# Test under high load
go test -run TestLoad -v ./...
go test -run TestStress -v ./...
go test -run TestConcurrent -v ./...
```

### Testing Best Practices

**Concurrency Testing**:
- **Always use `-race`** for concurrent code - this catches data races that can cause subtle bugs
- **Run tests multiple times** with `-count=N` to catch flaky tests
- **Test with different goroutine counts** to ensure scalability
- **Use `-timeout`** to prevent hanging tests

**Coverage Requirements**:
- **Aim for 90%+ coverage** on critical consensus and replication code
- **Focus on error paths** - they're often less tested but more important
- **Test edge cases** like empty inputs, boundary conditions, and error states
- **Use branch coverage** (`-covermode=atomic`) for concurrent code

**Failure Testing**:
- **Test network partitions** - simulate network splits and reconnections
- **Test server crashes** - verify recovery and state consistency
- **Test message reordering** - ensure algorithms work with out-of-order delivery
- **Test Byzantine failures** - verify behavior under malicious nodes

**Performance Testing**:
- **Benchmark critical paths** - consensus, replication, and shard operations
- **Measure memory usage** - prevent memory leaks in long-running systems
- **Test scalability** - verify performance with increasing node counts
- **Profile under load** - identify bottlenecks under realistic conditions

**Property-Based Testing**:
- **Use property-based testing** for complex algorithms like Paxos
- **Test invariants** - properties that should always hold
- **Generate random inputs** - test with diverse scenarios
- **Verify correctness properties** - safety and liveness guarantees

## Troubleshooting and Debugging

### Common Issues and Solutions

#### Build and Environment Problems

**Go Version Issues**:
```bash
# Check Go version
go version  # Should be 1.25.1+
go env GOROOT GOPATH  # Verify environment

# Fix module issues
go mod tidy
go mod download
go clean -modcache  # If modules are corrupted
```

**Import Path Issues**:
```bash
# Check import paths
go list -m all  # List all modules
go mod why <package>  # Understand why a package is needed
```

#### Concurrency and Race Conditions

**Race Condition Detection**:
```bash
# Always test with race detection
go test -race ./...
go test -race -count=10 ./...  # Run multiple times

# Common race condition patterns to avoid:
# - Shared mutable state without synchronization
# - Data races in concurrent map access
# - Race conditions in consensus algorithms
```

**Deadlock Detection**:
```bash
# Use Go's deadlock detector
go test -race -timeout=30s ./...
# Look for "fatal error: all goroutines are asleep - deadlock!"
```

**Goroutine Leaks**:
```bash
# Profile goroutines
go test -blockprofile=block.prof ./...
go tool pprof block.prof
(pprof) top
(pprof) list functionName
```

#### Performance Issues

**CPU Profiling**:
```bash
# Profile CPU usage
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof
(pprof) top10
(pprof) list functionName
(pprof) web  # Generate SVG graph
```

**Memory Profiling**:
```bash
# Profile memory usage
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
(pprof) top10
(pprof) list functionName
(pprof) web
```

**Common Performance Issues**:
- **Excessive memory allocations** in hot paths
- **Inefficient data structures** for the use case
- **Network round trips** that could be batched
- **Lock contention** in concurrent code

#### Network and Communication Issues

**Connection Problems**:
```bash
# Check port availability
netstat -an | grep :<port>
lsof -i :<port>

# Test network connectivity
telnet localhost <port>
```

**RPC Issues**:
- **Timeout configuration** - adjust for network latency
- **Connection pooling** - reuse connections when possible
- **Retry logic** - implement exponential backoff
- **Circuit breakers** - prevent cascading failures

**Common Network Issues**:
- **Firewall blocking ports** - check system firewall settings
- **Port conflicts** - ensure ports are available
- **Network partitions** - test with simulated network failures
- **Message ordering** - handle out-of-order message delivery

#### Algorithm-Specific Issues

**Paxos Consensus Problems**:
- **Split votes** - ensure proper proposal numbering
- **Log holes** - handle missing log entries correctly
- **Leader election** - implement proper leader selection
- **Garbage collection** - manage log compaction safely

**Sharding Issues**:
- **Shard movement** - ensure atomic shard transfers
- **Load balancing** - distribute load evenly across shards
- **Configuration changes** - handle dynamic reconfiguration
- **Cross-shard operations** - maintain consistency across shards

**Replication Problems**:
- **Split-brain scenarios** - prevent multiple primaries
- **State synchronization** - ensure consistent state across replicas
- **Failure detection** - avoid false positives in failure detection
- **Client failover** - handle client reconnection gracefully

### Debugging Tools and Techniques

**Go Debugging Tools**:
```bash
# Race detection
go test -race ./...

# Deadlock detection
go test -race -timeout=30s ./...

# Goroutine analysis
go test -blockprofile=block.prof ./...
go tool pprof block.prof

# Memory analysis
go test -memprofile=mem.prof ./...
go tool pprof mem.prof

# CPU profiling
go test -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof
```

**Logging and Tracing**:
```bash
# Enable verbose logging
go test -v ./...

# Use structured logging
log.Printf("Consensus: proposal=%d, value=%v", proposal, value)

# Add debug prints for complex algorithms
fmt.Printf("DEBUG: Paxos state: %+v\n", state)
```

**System-Level Debugging**:
```bash
# Monitor system resources
top -p $(pgrep -f "go test")
iostat -x 1  # Monitor I/O
netstat -i   # Monitor network interfaces

# Check for resource limits
ulimit -a
```

### Getting Help and Support

**Self-Help Checklist**:
1. **Read lab documentation** thoroughly - most issues are covered
2. **Run tests with verbose output**: `go test -v`
3. **Use debugging tools**: `go tool pprof`, `go test -race`
4. **Review error messages** carefully - they often contain the solution
5. **Check system resources** - memory, CPU, disk space
6. **Verify environment setup** - Go version, PATH, modules

**When to Seek Help**:
- Tests pass locally but fail in CI/CD
- Performance issues that profiling doesn't reveal
- Complex race conditions that are hard to reproduce
- Algorithm correctness questions
- Integration issues between components

**Providing Good Bug Reports**:
- **Include Go version**: `go version`
- **Include test output**: `go test -v -race`
- **Include system info**: OS, architecture, memory
- **Describe steps to reproduce**
- **Include relevant code snippets**
- **Explain expected vs actual behavior**

## Learning Resources

### Go Programming
- [Go Documentation](https://golang.org/doc/) - Official Go docs
- [Effective Go](https://go.dev/doc/effective_go) - Best practices
- [Go Concurrency Patterns](https://golang.org/doc/effective_go.html#concurrency) - Concurrency guide
- [A Tour of Go](https://go.dev/tour/welcome/1) - Interactive tutorial

### Distributed Systems Theory
- **Concepts**: `../reference/` - Theoretical foundations
- **Case Studies**: `../research/` - Real-world systems
- **Papers**: `../research/papers/` - Academic research

### External Resources
- [MIT 6.824](https://pdos.csail.mit.edu/6.824/) - Distributed Systems course
- [Raft Consensus](https://raft.github.io/) - Raft algorithm visualization
- [Paxos Made Simple](https://lamport.azurewebsites.net/pubs/paxos-simple.pdf) - Classic paper

## Learning Objectives

By completing these labs, you'll master:

- **Distributed Computing** — Master-worker patterns, task distribution
- **Fault Tolerance** — Failure detection, recovery, and prevention
- **Consensus Algorithms** — Paxos, Raft, and agreement protocols
- **Replication** — Primary-backup, state machine replication
- **Sharding** — Data partitioning, load balancing, reconfiguration
- **Persistence** — Durability, crash recovery, performance
- **Testing** — Concurrency testing, failure injection, benchmarking

## Contributing

We welcome contributions! Here's how to get started:

1. **Fork** the repository
2. **Create** a feature branch
3. **Implement** your changes with tests
4. **Run** all tests: `go test -race ./...`
5. **Submit** a pull request

### Contribution Guidelines
- Follow Go coding conventions
- Write comprehensive tests
- Document public APIs
- Include performance benchmarks
- Update documentation as needed

---

**Ready to build distributed systems?** Start with [Lab 1: MapReduce](01_Practice-Labs/lab1.md) and work your way through the guided exercises. Each lab builds upon the previous one, creating a comprehensive understanding of distributed systems engineering.

**Need help?** Check the [troubleshooting section](#-troubleshooting) or review the detailed lab documentation in `01_Practice-Labs/`.