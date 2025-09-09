# Distributed Systems Project: Building a Scalable Key-Value Store

## Overview

This project implements a comprehensive distributed systems curriculum focused on building a highly available, scalable, fault-tolerant, and transactional key-value store from scratch using Go. The project is designed to provide hands-on experience with fundamental distributed systems concepts through a series of progressively complex labs.

## Learning Objectives

By completing this project, you will gain practical experience with:

- **Distributed Systems Fundamentals**: Nodes, clients, servers, messages, timers, and RPC semantics
- **Consensus Algorithms**: Paxos and Multi-Paxos for fault-tolerant consensus
- **Replication**: Primary-backup replication with deterministic failure detection
- **Logical Clocks**: Lamport and Vector clocks for ordering events
- **Sharding**: Dynamic sharding with load balancing and reconfiguration
- **Distributed Transactions**: Two-Phase Commit (2PC) for cross-shard transactions
- **Fault Tolerance**: Handling server failures, network partitions, and message loss
- **Concurrency**: Go-specific concurrency patterns and synchronization primitives

## Project Architecture

The project is organized into four main labs, each building upon the previous:

```
Lab 1: Distributed Systems Framework
├── Node abstraction and message passing
├── Client-server interactions
├── Timer management
└── Testing framework

Lab 2: Primary-Backup Replication
├── ViewService for failure detection
├── Primary-backup state machine
├── State transfer and synchronization
└── Logical clocks (Lamport & Vector)

Lab 3: Paxos Consensus
├── Basic Paxos algorithm
├── Multi-Paxos for replicated logs
├── Leader election and optimization
└── Log synchronization and garbage collection

Lab 4: Sharded Key-Value Store
├── ShardMaster for configuration management
├── ShardStore with dynamic reconfiguration
├── Shard movement and load balancing
└── Distributed transactions with 2PC
```

## Technical Requirements

### System Requirements
- **Go Version**: 1.25.1 or above
- **Operating System**: Linux, macOS, or Windows
- **Memory**: Minimum 4GB RAM (8GB recommended)
- **Storage**: 1GB free space

### Dependencies
- Go standard library (no external dependencies required)
- Testing framework: `testing` package
- Logging: `slog` package for structured logging
- Concurrency: `sync`, `atomic`, and `context` packages

### Development Tools
- **IDE**: VS Code, GoLand, or any Go-compatible editor
- **Testing**: `go test` with race detection
- **Profiling**: `go tool pprof` for performance analysis
- **Linting**: `golangci-lint` for code quality

## Lab Descriptions

### Lab 1: Distributed Systems Framework
**Objective**: Build the foundational framework for distributed systems

**Key Concepts**:
- Node abstraction and message passing
- Client-server RPC interactions
- Timer management and event handling
- Application state management
- Testing and debugging tools

**Deliverables**:
- `Node` interface and implementation
- `Message` types and serialization
- `Client` and `Server` abstractions
- `Application` interface for state machines
- Comprehensive test suite

**Learning Outcomes**:
- Understand distributed systems architecture
- Master Go concurrency patterns
- Implement robust testing strategies
- Handle network communication and failures

### Lab 2: Primary-Backup Replication
**Objective**: Implement fault-tolerant replication using primary-backup pattern

**Key Concepts**:
- ViewService for failure detection
- Primary-backup state machine replication
- Deterministic failure detection
- State transfer and synchronization
- Logical clocks (Lamport and Vector)

**Deliverables**:
- `ViewService` for managing server views
- `PBServer` with primary-backup logic
- `PBClient` with automatic failover
- State transfer mechanisms
- Logical clock implementations

**Learning Outcomes**:
- Understand replication strategies
- Implement failure detection algorithms
- Master state synchronization techniques
- Handle split-brain scenarios

### Lab 3: Paxos Consensus
**Objective**: Implement the Paxos consensus algorithm for fault-tolerant agreement

**Key Concepts**:
- Basic Paxos algorithm (Prepare/Accept phases)
- Multi-Paxos for replicated logs
- Leader election and optimization
- Log synchronization and garbage collection
- Hole handling in Paxos logs

**Deliverables**:
- `PaxosNode` with basic consensus
- `MultiPaxosNode` for replicated logs
- Leader election mechanisms
- Log synchronization algorithms
- Garbage collection strategies

**Learning Outcomes**:
- Master consensus algorithms
- Understand leader election
- Implement replicated state machines
- Handle network partitions and failures

### Lab 4: Sharded Key-Value Store
**Objective**: Build a complete sharded key-value store with distributed transactions

**Key Concepts**:
- ShardMaster for configuration management
- ShardStore with dynamic reconfiguration
- Shard movement and load balancing
- Distributed transactions with 2PC
- Cross-shard operations

**Deliverables**:
- `ShardMaster` for shard management
- `ShardStore` with Paxos replication
- `ShardStoreClient` with automatic routing
- Transaction coordinator and participants
- Two-Phase Commit implementation

**Learning Outcomes**:
- Understand sharding and partitioning
- Implement dynamic reconfiguration
- Master distributed transactions
- Handle cross-shard operations

## Project Structure

```
app/02_Project/
├── README.md                 # This file
├── go.mod                    # Go module definition
├── lab1.md                   # Lab 1: Distributed Systems Framework
├── lab2a.md                  # Lab 2A: ViewService and Primary-Backup
├── lab2b.md                  # Lab 2B: State Transfer and Logical Clocks
├── lab3a.md                  # Lab 3A: Basic Paxos
├── lab3b.md                  # Lab 3B: Multi-Paxos
├── lab4a.md                  # Lab 4A: ShardMaster
├── lab4b.md                  # Lab 4B: ShardStore
├── lab4c.md                  # Lab 4C: Distributed Transactions
├── src/                      # Source code directory
│   ├── common/               # Common types and utilities
│   ├── lab1/                 # Lab 1 implementation
│   ├── lab2/                 # Lab 2 implementation
│   ├── lab3/                 # Lab 3 implementation
│   └── lab4/                 # Lab 4 implementation
├── tests/                    # Test files
│   ├── unit/                 # Unit tests
│   ├── integration/          # Integration tests
│   └── benchmarks/           # Performance benchmarks
└── docs/                     # Additional documentation
    ├── architecture.md       # System architecture
    ├── api.md               # API documentation
    └── troubleshooting.md   # Troubleshooting guide
```

## Getting Started

### Prerequisites
1. Install Go 1.25.1 or above
2. Set up your development environment
3. Clone the repository
4. Navigate to the project directory

### Building the Project
```bash
# Initialize Go module
go mod init distributed-systems-project

# Build all components
go build ./...

# Run tests
go test ./...

# Run with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./...
```

### Running Individual Labs
```bash
# Lab 1: Distributed Systems Framework
go run src/lab1/main.go

# Lab 2: Primary-Backup Replication
go run src/lab2/main.go

# Lab 3: Paxos Consensus
go run src/lab3/main.go

# Lab 4: Sharded Key-Value Store
go run src/lab4/main.go
```

## Testing Strategy

### Unit Testing
- **Coverage**: Aim for 90%+ code coverage
- **Race Detection**: Always run tests with `-race` flag
- **Property Testing**: Use property-based testing for complex algorithms
- **Mocking**: Mock external dependencies and network calls

### Integration Testing
- **End-to-End**: Test complete workflows
- **Failure Scenarios**: Test network failures, server crashes, and partitions
- **Performance**: Measure latency, throughput, and resource usage
- **Scalability**: Test with increasing numbers of nodes and clients

### Testing Commands
```bash
# Run all tests
go test ./...

# Run specific lab tests
go test ./src/lab1/...
go test ./src/lab2/...
go test ./src/lab3/...
go test ./src/lab4/...

# Run with coverage
go test -cover ./...

# Run with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./...

# Run specific test
go test -run TestSpecificFunction ./...
```

## Performance Considerations

### Optimization Strategies
- **Concurrency**: Use goroutines and channels effectively
- **Memory Management**: Minimize allocations and use object pools
- **Network Efficiency**: Batch operations and use compression
- **Caching**: Implement intelligent caching strategies
- **Load Balancing**: Distribute load evenly across nodes

### Monitoring and Profiling
```bash
# CPU profiling
go test -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof ./...
go tool pprof mem.prof

# Race detection
go test -race ./...

# Benchmarking
go test -bench=. -benchmem ./...
```

## Troubleshooting

### Common Issues

#### Build Issues
- **Go Version**: Ensure you're using Go 1.25.1 or above
- **Module Issues**: Run `go mod tidy` to resolve dependencies
- **Import Paths**: Check that import paths are correct

#### Runtime Issues
- **Race Conditions**: Use `go test -race` to detect race conditions
- **Deadlocks**: Use `go tool pprof` to analyze goroutine stacks
- **Memory Leaks**: Monitor memory usage with `go tool pprof`

#### Network Issues
- **Connection Failures**: Check network configuration and firewall settings
- **Timeout Issues**: Adjust timeout values based on network conditions
- **Port Conflicts**: Ensure ports are available and not in use

### Debugging Tools
```bash
# Race detection
go test -race ./...

# Profiling
go test -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof

# Goroutine analysis
go test -blockprofile=block.prof ./...
go tool pprof block.prof

# Memory analysis
go test -memprofile=mem.prof ./...
go tool pprof mem.prof
```

## Best Practices

### Code Quality
- **Go Conventions**: Follow Go coding standards and conventions
- **Error Handling**: Always handle errors explicitly
- **Documentation**: Document all public APIs and complex algorithms
- **Testing**: Write comprehensive tests for all functionality
- **Logging**: Use structured logging with appropriate levels

### Concurrency
- **Goroutines**: Use goroutines for concurrent operations
- **Channels**: Prefer channels over shared memory for communication
- **Synchronization**: Use appropriate synchronization primitives
- **Context**: Use context for cancellation and timeouts
- **Atomic Operations**: Use atomic operations for simple shared state

### Performance
- **Profiling**: Regularly profile your code for performance issues
- **Benchmarking**: Write benchmarks for critical code paths
- **Memory Management**: Be mindful of memory allocations
- **Network Efficiency**: Minimize network round trips
- **Caching**: Implement appropriate caching strategies

## Contributing

### Development Workflow
1. **Fork** the repository
2. **Create** a feature branch
3. **Implement** your changes
4. **Write** comprehensive tests
5. **Run** all tests and ensure they pass
6. **Submit** a pull request

### Code Review Process
- **Automated Tests**: All tests must pass
- **Code Quality**: Code must follow Go conventions
- **Documentation**: All public APIs must be documented
- **Performance**: No significant performance regressions
- **Security**: No security vulnerabilities

### Pull Request Guidelines
- **Clear Description**: Describe what the PR does and why
- **Test Coverage**: Include tests for new functionality
- **Documentation**: Update documentation as needed
- **Breaking Changes**: Clearly mark any breaking changes
- **Performance Impact**: Document any performance implications

## Resources

### Documentation
- [Go Documentation](https://golang.org/doc/)
- [Go Concurrency Patterns](https://golang.org/doc/effective_go.html#concurrency)
- [Go Testing](https://golang.org/doc/tutorial/add-a-test)
- [Go Profiling](https://golang.org/doc/diagnostics.html)

### Distributed Systems
- [Paxos Made Simple](https://lamport.azurewebsites.net/pubs/paxos-simple.pdf)
- [Raft Consensus Algorithm](https://raft.github.io/)
- [Two-Phase Commit](https://en.wikipedia.org/wiki/Two-phase_commit_protocol)
- [Distributed Systems Concepts](https://en.wikipedia.org/wiki/Distributed_computing)

### Tools and Libraries
- [Go Testing Package](https://pkg.go.dev/testing)
- [Go Sync Package](https://pkg.go.dev/sync)
- [Go Context Package](https://pkg.go.dev/context)
- [Go Log Package](https://pkg.go.dev/log)

## Acknowledgments

- **MIT 6.824**: Distributed Systems course for inspiration
- **Go Team**: For the excellent Go programming language
- **Distributed Systems Community**: For research papers and best practices
- **Open Source Contributors**: For tools and libraries used in this project

## Support

If you encounter any issues or have questions:

1. **Check** the troubleshooting section above
2. **Search** existing issues in the repository
3. **Create** a new issue with detailed information
4. **Join** the community discussions
5. **Read** the documentation and lab materials

---

**Happy Coding!**

This project will give you hands-on experience with building real distributed systems. Take your time, understand each concept thoroughly, and don't hesitate to experiment and explore different approaches.
