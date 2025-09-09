# Distributed Systems Applications

[![Go Version](https://img.shields.io/badge/Go-1.25.1-blue.svg)](https://golang.org)
[![Status](https://img.shields.io/badge/Status-Active%20Development-green.svg)](https://github.com/yourusername/Distributed-Systems)

> **Hands-on programming exercises and implementations of distributed systems concepts**

This directory contains practical implementations and guided exercises that bring distributed systems theory to life. Build real systems, understand failure modes, and master the engineering challenges of distributed computing.

## Quick Start

### Prerequisites
- **Go 1.25.1+** [[memory:8397771]]
- **4GB+ RAM** (8GB recommended for large datasets)
- **Basic Go knowledge** (goroutines, channels, interfaces)

### 30-Second Demo
```bash
# 1. Verify Go installation
go version  # Should show go1.25.1 or later

# 2. Run MapReduce word count
cd app/01_Practice-Labs/src/main
go run wc.go master kjv12.txt sequential

# 3. Test sharded key-value store
cd ../shardmaster && go test
cd ../shardkv && go test
```

## Learning Path

### Phase 1: Fundamentals
Start with the **Practice Labs** for guided, hands-on learning:

| Lab | Topic | Duration | Key Skills |
|-----|-------|----------|------------|
| **Lab 1** | MapReduce | 4-6 hours | Distributed computing, fault tolerance |
| **Lab 2** | Primary-Backup | 6-8 hours | Replication, failure detection |
| **Lab 3** | Paxos Consensus | 8-12 hours | Consensus algorithms, leader election |
| **Lab 4** | Sharded KV Store | 10-15 hours | Sharding, load balancing, reconfiguration |
| **Lab 5** | Persistent Storage | 6-10 hours | Durability, crash recovery, performance |

### Phase 2: Advanced Projects
Extend your knowledge with self-directed implementations:

- **RPC Framework** — Build reliable remote procedure calls
- **Consensus Library** — Implement Raft or PBFT
- **Distributed Cache** — Create a fault-tolerant caching layer
- **Service Mesh** — Build service discovery and load balancing
- **Observability Stack** — Add metrics, logging, and tracing

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

## Lab Overview

### Lab 1: MapReduce
**Goal**: Build a distributed data processing system
- Master-worker architecture
- Task distribution and fault tolerance
- Large-scale data processing

```bash
cd app/01_Practice-Labs/src/mapreduce
go test
```

### Lab 2: Primary-Backup Replication
**Goal**: Create a fault-tolerant key-value store
- ViewService for failure detection
- State machine replication
- Client failover mechanisms

```bash
cd app/01_Practice-Labs/src/viewservice && go test
cd ../pbservice && go test
```

### Lab 3: Paxos Consensus
**Goal**: Implement distributed consensus
- Basic Paxos algorithm
- Multi-Paxos for replicated logs
- Leader election and optimization

```bash
cd app/01_Practice-Labs/src/paxos
go test -race  # Always use race detection
```

### Lab 4: Sharded Key-Value Store
**Goal**: Build a scalable, sharded system
- Dynamic configuration management
- Shard movement and load balancing
- Cross-shard operations

```bash
cd app/01_Practice-Labs/src/shardmaster && go test
cd ../shardkv && go test -race
```

### Lab 5: Persistent Storage
**Goal**: Add durability and crash recovery
- Atomic file operations
- Crash recovery procedures
- Performance optimization

```bash
cd app/01_Practice-Labs/src/diskv
go test -run Test4  # Lab 4 compatibility
go test             # Full Lab 5 suite
```

## Development Workflow

### 1. Read Theory First
```bash
# Study the concepts before coding
cd ../reference
# Read relevant sections for each lab
```

### 2. Implement Incrementally
- Start with the simplest working version
- Add features one at a time
- Test thoroughly after each change
- Use `go test -race` for concurrency bugs

### 3. Test Failure Scenarios
- Simulate network partitions
- Test server crashes and recovery
- Verify correctness under failures
- Measure performance impact

### 4. Benchmark and Optimize
```bash
# Run performance benchmarks
go test -bench=. -benchmem

# Profile for bottlenecks
go test -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

## Testing Strategy

### Essential Test Commands
```bash
# Run all tests with race detection
go test -race ./...

# Test specific lab
go test -race ./src/mapreduce/...

# Run with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...

# Profile memory usage
go test -memprofile=mem.prof ./...
go tool pprof mem.prof
```

### Testing Best Practices
- **Always use `-race`** for concurrent code
- **Aim for 90%+ coverage** on critical paths
- **Test failure scenarios** extensively
- **Benchmark performance** regularly
- **Use property-based testing** for complex algorithms

## Troubleshooting

### Common Issues

#### Build Problems
```bash
# Wrong Go version
go version  # Should be 1.25.1+
go mod tidy # Fix dependency issues
```

#### Race Conditions
```bash
# Always test with race detection
go test -race ./...
# Fix any reported race conditions
```

#### Performance Issues
```bash
# Profile CPU usage
go test -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof

# Profile memory usage
go test -memprofile=mem.prof ./...
go tool pprof mem.prof
```

#### Network Issues
- Check firewall settings
- Verify port availability
- Adjust timeout values for slow networks
- Test with different network conditions

### Getting Help
1. **Check lab documentation** first
2. **Run tests with verbose output**: `go test -v`
3. **Use debugging tools**: `go tool pprof`, `go test -race`
4. **Review error messages** carefully
5. **Search existing issues** in the repository

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