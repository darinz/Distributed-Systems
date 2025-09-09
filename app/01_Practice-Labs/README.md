# Distributed Systems Practice Labs

## Overview

This repository contains hands-on practice labs for learning distributed systems concepts through practical implementation. These labs are designed to build your understanding of fundamental distributed systems principles by implementing real-world systems using Go.

**Prerequisites**: Go version 1.25.1 or later

## Learning Objectives

By completing these practice labs, you will gain practical experience with:

- **Distributed Systems Fundamentals**: RPC, message passing, and network communication
- **Fault Tolerance**: Handling server failures and network partitions
- **Consensus Algorithms**: Understanding agreement protocols
- **Replication**: Primary-backup and state machine replication
- **Sharding**: Data partitioning and load balancing
- **Persistence**: Data durability and recovery
- **MapReduce**: Large-scale data processing

## Lab Progression

The labs are designed to be completed in order, with each building upon concepts from previous labs:

```
Lab 1: MapReduce
├── Master-worker architecture
├── Task distribution and coordination
├── Fault tolerance and recovery
└── Large-scale data processing

Lab 2: Primary-Backup Replication
├── ViewService for failure detection
├── Primary-backup state machine
├── State transfer and synchronization
└── Client failover mechanisms

Lab 3: Paxos Consensus
├── Basic Paxos algorithm
├── Multi-Paxos for replicated logs
├── Leader election and optimization
└── Log synchronization

Lab 4: Sharded Key-Value Store
├── ShardMaster for configuration management
├── ShardStore with dynamic reconfiguration
├── Shard movement and load balancing
└── Cross-shard operations

Lab 5: Persistent Storage
├── Atomic file operations
├── Crash recovery and durability
├── Disk layout optimization
└── Performance tuning
```

## Lab Descriptions

### Lab 1: MapReduce
**Objective**: Implement a distributed MapReduce system for large-scale data processing

**Key Concepts**:
- Master-worker architecture
- Task distribution and coordination
- Fault tolerance and recovery
- Large-scale data processing

**Learning Outcomes**:
- Understand distributed computing paradigms
- Implement fault-tolerant task scheduling
- Handle worker failures and recovery
- Process large datasets efficiently

**Files**:
- `lab1.md`: Detailed lab instructions
- `src/mapreduce/`: MapReduce implementation
- `src/main/`: Example applications (word count, etc.)

### Lab 2: Primary-Backup Replication
**Objective**: Build a fault-tolerant key-value store using primary-backup replication

**Key Concepts**:
- ViewService for failure detection
- Primary-backup state machine replication
- State transfer and synchronization
- Client failover mechanisms

**Learning Outcomes**:
- Understand replication strategies
- Implement failure detection algorithms
- Master state synchronization techniques
- Handle split-brain scenarios

**Files**:
- `lab2a.md`: ViewService and primary-backup setup
- `lab2b.md`: State transfer and logical clocks
- `src/viewservice/`: ViewService implementation
- `src/pbservice/`: Primary-backup service

### Lab 3: Paxos Consensus
**Objective**: Implement the Paxos consensus algorithm for fault-tolerant agreement

**Key Concepts**:
- Basic Paxos algorithm (Prepare/Accept phases)
- Multi-Paxos for replicated logs
- Leader election and optimization
- Log synchronization and garbage collection

**Learning Outcomes**:
- Master consensus algorithms
- Understand leader election
- Implement replicated state machines
- Handle network partitions and failures

**Files**:
- `lab3a.md`: Basic Paxos implementation
- `lab3b.md`: Multi-Paxos and replicated logs
- `src/paxos/`: Paxos consensus implementation

### Lab 4: Sharded Key-Value Store
**Objective**: Build a scalable, sharded key-value store with dynamic reconfiguration

**Key Concepts**:
- ShardMaster for configuration management
- ShardStore with dynamic reconfiguration
- Shard movement and load balancing
- Cross-shard operations

**Learning Outcomes**:
- Understand sharding and partitioning
- Implement dynamic reconfiguration
- Master load balancing techniques
- Handle cross-shard operations

**Files**:
- `lab4a.md`: ShardMaster implementation
- `lab4b.md`: ShardStore with reconfiguration
- `src/shardmaster/`: ShardMaster service
- `src/shardkv/`: Sharded key-value store

### Lab 5: Persistent Storage
**Objective**: Add persistence and crash recovery to the sharded key-value store

**Key Concepts**:
- Atomic file operations
- Crash recovery and durability
- Disk layout optimization
- Performance tuning

**Learning Outcomes**:
- Understand persistence mechanisms
- Implement crash recovery
- Master disk layout optimization
- Handle performance tuning

**Files**:
- `lab5.md`: Persistence implementation
- `src/diskv/`: Persistent disk storage

## Getting Started

### Prerequisites

#### System Requirements
- **Go Version**: 1.25.1 or later
- **Operating System**: Linux, macOS, or Windows
- **Memory**: Minimum 4GB RAM (8GB recommended)
- **Storage**: 2GB free space

#### Development Environment
- **IDE**: VS Code, GoLand, or any Go-compatible editor
- **Terminal**: Command-line interface for running tests
- **Git**: Version control for tracking changes

### Go Installation

#### macOS

**Option 1: Official Installer (Recommended)**
1. Download the macOS installer from [go.dev/dl](https://go.dev/dl/)
2. Select the appropriate architecture (Intel or Apple Silicon)
3. Run the `.pkg` installer
4. Verify installation:
   ```bash
   go version
   ```

**Option 2: Homebrew**
```bash
brew install go@1.25
```

**Option 3: Version Manager (g)**
```bash
# Install g
curl -sSL https://git.io/g-install | sh -s
# Install Go 1.25.1
g install 1.25.1
```

#### Linux

**Option 1: Official Binary (Recommended)**
```bash
# Download and extract
wget https://go.dev/dl/go1.25.1.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.1.linux-amd64.tar.gz

# Add to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH=$PATH:/usr/local/go/bin
```

**Option 2: Package Manager**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang

# CentOS/RHEL/Fedora
sudo dnf install golang
# or
sudo yum install golang
```

**Option 3: Version Manager (gvm)**
```bash
# Install gvm
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
# Install Go 1.25.1
gvm install go1.25.1
gvm use go1.25.1 --default
```

#### Windows

**Option 1: Official Installer (Recommended)**
1. Download the Windows installer from [go.dev/dl](https://go.dev/dl/)
2. Run the `.msi` installer
3. Follow the installation wizard
4. Verify installation:
   ```cmd
   go version
   ```

**Option 2: Chocolatey**
```cmd
choco install golang
```

**Option 3: Scoop**
```cmd
scoop install go
```

### Environment Setup

**For Modern Go Development (Recommended):**
```bash
# Initialize Go modules in your project
cd /path/to/your/project
go mod init distributed-systems-practice

# Verify Go environment
go env GOPATH
go env GOROOT
```

**For Legacy GOPATH Setup (if needed):**
```bash
# Set GOPATH (not recommended for new projects)
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

## Running the Labs

### Lab 1: MapReduce
```bash
# Navigate to MapReduce directory
cd app/01_Practice-Labs/src/mapreduce

# Run tests
go test

# Run example applications
cd ../main
go run wc.go master kjv12.txt sequential
go run wc.go master kjv12.txt distributed
```

### Lab 2: Primary-Backup Replication
```bash
# Navigate to ViewService directory
cd app/01_Practice-Labs/src/viewservice

# Run ViewService tests
go test

# Navigate to Primary-Backup directory
cd ../pbservice

# Run Primary-Backup tests
go test
```

### Lab 3: Paxos Consensus
```bash
# Navigate to Paxos directory
cd app/01_Practice-Labs/src/paxos

# Run Paxos tests
go test

# Run with race detection
go test -race
```

### Lab 4: Sharded Key-Value Store
```bash
# Navigate to ShardMaster directory
cd app/01_Practice-Labs/src/shardmaster

# Run ShardMaster tests
go test

# Navigate to ShardKV directory
cd ../shardkv

# Run ShardKV tests
go test

# Run with race detection
go test -race
```

### Lab 5: Persistent Storage
```bash
# Navigate to DiskV directory
cd app/01_Practice-Labs/src/diskv

# Run Lab 4 compatible tests
go test -run Test4

# Run full Lab 5 test suite
go test

# Run with race detection
go test -race
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
go test ./src/mapreduce/...
go test ./src/viewservice/...
go test ./src/pbservice/...
go test ./src/paxos/...
go test ./src/shardmaster/...
go test ./src/shardkv/...
go test ./src/diskv/...

# Run with coverage
go test -cover ./...

# Run with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./...

# Run specific test
go test -run TestSpecificFunction ./...
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

### Lab-Specific Troubleshooting

#### Lab 1: MapReduce
- **Worker Failures**: Ensure proper task reassignment
- **Memory Issues**: Monitor memory usage for large datasets
- **Network Timeouts**: Adjust timeout values for slow networks

#### Lab 2: Primary-Backup
- **Split-Brain**: Implement proper failure detection
- **State Transfer**: Ensure atomic state transfer
- **Client Failover**: Handle client reconnection properly

#### Lab 3: Paxos
- **Leader Election**: Ensure proper leader selection
- **Log Synchronization**: Handle log holes and gaps
- **Garbage Collection**: Implement proper log cleanup

#### Lab 4: Sharding
- **Configuration Changes**: Handle dynamic reconfiguration
- **Shard Movement**: Ensure atomic shard transfers
- **Load Balancing**: Distribute load evenly

#### Lab 5: Persistence
- **Crash Recovery**: Implement proper recovery procedures
- **Atomic Operations**: Ensure file operations are atomic
- **Performance**: Optimize disk I/O operations

## Learning Resources

### Documentation
- [Go Documentation](https://golang.org/doc/)
- [Go Concurrency Patterns](https://golang.org/doc/effective_go.html#concurrency)
- [Go Testing](https://golang.org/doc/tutorial/add-a-test)
- [Go Profiling](https://golang.org/doc/diagnostics.html)

### Distributed Systems
- [Paxos Made Simple](https://lamport.azurewebsites.net/pubs/paxos-simple.pdf)
- [Raft Consensus Algorithm](https://raft.github.io/)
- [MapReduce Paper](https://static.googleusercontent.com/media/research.google.com/en//archive/mapreduce-osdi04.pdf)
- [Distributed Systems Concepts](https://en.wikipedia.org/wiki/Distributed_computing)

### Tools and Libraries
- [Go Testing Package](https://pkg.go.dev/testing)
- [Go Sync Package](https://pkg.go.dev/sync)
- [Go Context Package](https://pkg.go.dev/context)
- [Go Log Package](https://pkg.go.dev/log)

## Project Structure

```
app/01_Practice-Labs/
├── README.md                 # This file
├── lab1.md                   # Lab 1: MapReduce
├── lab2a.md                  # Lab 2A: ViewService
├── lab2b.md                  # Lab 2B: Primary-Backup
├── lab3a.md                  # Lab 3A: Basic Paxos
├── lab3b.md                  # Lab 3B: Multi-Paxos
├── lab4a.md                  # Lab 4A: ShardMaster
├── lab4b.md                  # Lab 4B: ShardStore
├── lab5.md                   # Lab 5: Persistence
├── src/                      # Source code directory
│   ├── mapreduce/            # Lab 1: MapReduce implementation
│   ├── viewservice/          # Lab 2A: ViewService implementation
│   ├── pbservice/            # Lab 2B: Primary-Backup implementation
│   ├── paxos/                # Lab 3: Paxos implementation
│   ├── shardmaster/          # Lab 4A: ShardMaster implementation
│   ├── shardkv/              # Lab 4B: ShardStore implementation
│   ├── diskv/                # Lab 5: Persistent storage
│   └── main/                 # Example applications
└── img/                      # Images and diagrams
```

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

## Support

If you encounter any issues or have questions:

1. **Check** the troubleshooting section above
2. **Search** existing issues in the repository
3. **Create** a new issue with detailed information
4. **Join** the community discussions
5. **Read** the lab documentation and materials

## Acknowledgments

- **MIT 6.824**: Distributed Systems course for inspiration
- **Go Team**: For the excellent Go programming language
- **Distributed Systems Community**: For research papers and best practices
- **Open Source Contributors**: For tools and libraries used in this project

---

**Happy Learning!**

These practice labs will give you hands-on experience with building real distributed systems. Take your time, understand each concept thoroughly, and don't hesitate to experiment and explore different approaches.