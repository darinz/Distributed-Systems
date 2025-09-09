# Lab 1: Building Distributed Systems from Scratch with Go

## Overview

This lab introduces the fundamental concepts and patterns for building distributed systems from scratch using the Go programming language. You'll learn how to design and implement a highly available, scalable, fault-tolerant, and transactional key-value store that serves as the foundation for modern cloud computing systems.

## Prerequisites

- Go 1.25.1 or later
- Understanding of basic Go concepts (goroutines, channels, interfaces)
- Familiarity with distributed systems concepts (consensus, replication, fault tolerance)

## Core Architecture

### Framework Design

The distributed systems framework is built around several key abstractions:

```go
// Core interfaces for the distributed system
type Node interface {
    Start() error
    Stop() error
    Send(to NodeID, msg Message) error
    SetTimer(duration time.Duration, timer Timer) error
}

type Message interface {
    GetType() string
    GetFrom() NodeID
    GetTo() NodeID
    Serialize() ([]byte, error)
    Deserialize(data []byte) error
}

type Timer interface {
    GetID() string
    GetData() interface{}
}
```

### Node Implementation

A Node represents a single machine or process in the distributed system. Here's how to implement one:

```go
type BaseNode struct {
    id       NodeID
    peers    map[NodeID]Node
    handlers map[string]MessageHandler
    timers   map[string]*time.Timer
    mu       sync.RWMutex
    done     chan struct{}
}

type MessageHandler func(msg Message) error

// Key methods to implement:
func (n *BaseNode) Start() error {
    // Initialize node state
    // Start message processing loop
    // Set up periodic tasks
}

func (n *BaseNode) Send(to NodeID, msg Message) error {
    // Implement reliable message delivery
    // Handle network failures and retries
}

func (n *BaseNode) SetTimer(duration time.Duration, timer Timer) error {
    // Set up timeout mechanisms
    // Handle timer expiration
}
```

**Critical Design Principles:**
- **Deterministic**: Same input should always produce same output
- **Idempotent**: Operations should be safe to retry
- **Fault-tolerant**: Handle node failures gracefully
- **Thread-safe**: Use proper synchronization primitives

### Client/Server Patterns

#### Client Implementation

```go
type Client interface {
    SendCommand(cmd Command) (Result, error)
    HasResult() bool
    GetResult() (Result, error)
}

type KVClient struct {
    servers    []NodeID
    currentSeq int64
    pending    map[int64]chan Result
    mu         sync.RWMutex
}

func (c *KVClient) SendCommand(cmd Command) (Result, error) {
    c.mu.Lock()
    seq := c.currentSeq
    c.currentSeq++
    resultChan := make(chan Result, 1)
    c.pending[seq] = resultChan
    c.mu.Unlock()
    
    // Send command with sequence number
    msg := &CommandMessage{
        Command: cmd,
        Seq:     seq,
        ClientID: c.id,
    }
    
    // Try servers in round-robin fashion
    for _, server := range c.servers {
        if err := c.sendToServer(server, msg); err == nil {
            break
        }
    }
    
    // Wait for result with timeout
    select {
    case result := <-resultChan:
        return result, nil
    case <-time.After(5 * time.Second):
        return nil, errors.New("command timeout")
    }
}
```

#### Server Implementation

```go
type Server struct {
    node      Node
    app       Application
    handlers  map[string]MessageHandler
    log       []LogEntry
    state     StateMachine
    mu        sync.RWMutex
}

type Application interface {
    Execute(cmd Command) (Result, error)
    GetState() interface{}
    ApplyLogEntry(entry LogEntry) error
}

func (s *Server) handleCommand(msg *CommandMessage) error {
    // Validate command
    // Apply to state machine
    // Send response back to client
    result, err := s.app.Execute(msg.Command)
    if err != nil {
        return err
    }
    
    response := &ResultMessage{
        Result: result,
        Seq:    msg.Seq,
        ClientID: msg.ClientID,
    }
    
    return s.node.Send(msg.ClientID, response)
}
```

### Message Handling and Serialization

#### Message Types

```go
// Base message structure
type BaseMessage struct {
    Type     string    `json:"type"`
    From     NodeID    `json:"from"`
    To       NodeID    `json:"to"`
    Timestamp time.Time `json:"timestamp"`
}

// Specific message types
type CommandMessage struct {
    BaseMessage
    Command  Command `json:"command"`
    Seq      int64   `json:"seq"`
    ClientID NodeID  `json:"client_id"`
}

type ResultMessage struct {
    BaseMessage
    Result   Result `json:"result"`
    Seq      int64  `json:"seq"`
    ClientID NodeID `json:"client_id"`
}

type HeartbeatMessage struct {
    BaseMessage
    Term     int64 `json:"term"`
    LeaderID NodeID `json:"leader_id"`
}
```

#### Serialization Best Practices

```go
// Use JSON for simplicity, consider Protocol Buffers for production
func (m *BaseMessage) Serialize() ([]byte, error) {
    return json.Marshal(m)
}

func (m *BaseMessage) Deserialize(data []byte) error {
    return json.Unmarshal(data, m)
}

// Message factory for type-safe deserialization
func DeserializeMessage(data []byte) (Message, error) {
    var base BaseMessage
    if err := json.Unmarshal(data, &base); err != nil {
        return nil, err
    }
    
    switch base.Type {
    case "command":
        var msg CommandMessage
        json.Unmarshal(data, &msg)
        return &msg, nil
    case "result":
        var msg ResultMessage
        json.Unmarshal(data, &msg)
        return &msg, nil
    default:
        return nil, fmt.Errorf("unknown message type: %s", base.Type)
    }
}
```

### Concurrency Patterns in Go

#### Goroutine Management

```go
type NodeManager struct {
    nodes    map[NodeID]*BaseNode
    network  Network
    done     chan struct{}
    wg       sync.WaitGroup
}

func (nm *NodeManager) Start() error {
    for _, node := range nm.nodes {
        nm.wg.Add(1)
        go func(n *BaseNode) {
            defer nm.wg.Done()
            n.Start()
        }(node)
    }
    return nil
}

func (nm *NodeManager) Stop() error {
    close(nm.done)
    nm.wg.Wait()
    return nil
}
```

#### Channel-based Communication

```go
type MessageRouter struct {
    incoming chan Message
    outgoing chan Message
    handlers map[string]chan Message
    mu       sync.RWMutex
}

func (mr *MessageRouter) Route() {
    for {
        select {
        case msg := <-mr.incoming:
            mr.mu.RLock()
            handler := mr.handlers[msg.GetType()]
            mr.mu.RUnlock()
            
            select {
            case handler <- msg:
            case <-time.After(100 * time.Millisecond):
                // Handle message routing timeout
            }
        case <-mr.done:
            return
        }
    }
}
```

#### Synchronization Primitives

```go
// Use sync.RWMutex for read-heavy workloads
type StateManager struct {
    state interface{}
    mu    sync.RWMutex
}

func (sm *StateManager) GetState() interface{} {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    return sm.state
}

func (sm *StateManager) UpdateState(newState interface{}) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    sm.state = newState
}

// Use sync.Cond for condition variables
type ResultWaiter struct {
    result Result
    ready  bool
    cond   *sync.Cond
    mu     sync.Mutex
}

func (rw *ResultWaiter) WaitForResult() Result {
    rw.mu.Lock()
    defer rw.mu.Unlock()
    
    for !rw.ready {
        rw.cond.Wait()
    }
    return rw.result
}

func (rw *ResultWaiter) SetResult(result Result) {
    rw.mu.Lock()
    defer rw.mu.Unlock()
    rw.result = result
    rw.ready = true
    rw.cond.Signal()
}
```

### Timer and Timeout Management

```go
type TimerManager struct {
    timers map[string]*time.Timer
    mu     sync.RWMutex
}

func (tm *TimerManager) SetTimer(id string, duration time.Duration, callback func()) {
    tm.mu.Lock()
    defer tm.mu.Unlock()
    
    // Cancel existing timer if any
    if timer, exists := tm.timers[id]; exists {
        timer.Stop()
    }
    
    tm.timers[id] = time.AfterFunc(duration, callback)
}

func (tm *TimerManager) CancelTimer(id string) {
    tm.mu.Lock()
    defer tm.mu.Unlock()
    
    if timer, exists := tm.timers[id]; exists {
        timer.Stop()
        delete(tm.timers, id)
    }
}

// Example usage in Raft leader election
func (r *RaftNode) startElectionTimer() {
    timeout := time.Duration(rand.Intn(150)+150) * time.Millisecond
    r.timerManager.SetTimer("election", timeout, func() {
        r.startElection()
    })
}
```

### Testing Strategy

#### Unit Testing

```go
func TestKVClient(t *testing.T) {
    // Create mock servers
    servers := []NodeID{"server1", "server2", "server3"}
    client := NewKVClient(servers)
    
    // Test successful command
    cmd := &PutCommand{Key: "test", Value: "value"}
    result, err := client.SendCommand(cmd)
    assert.NoError(t, err)
    assert.Equal(t, "OK", result.Status)
    
    // Test timeout scenario
    // ... test implementation
}
```

#### Integration Testing

```go
func TestDistributedSystem(t *testing.T) {
    // Create network simulation
    network := NewNetworkSimulator()
    
    // Create nodes
    nodes := make(map[NodeID]Node)
    for i := 0; i < 3; i++ {
        nodeID := NodeID(fmt.Sprintf("node%d", i))
        nodes[nodeID] = NewRaftNode(nodeID, network)
    }
    
    // Start all nodes
    for _, node := range nodes {
        go node.Start()
    }
    
    // Run test scenarios
    testLeaderElection(t, nodes)
    testLogReplication(t, nodes)
    testFaultTolerance(t, nodes)
    
    // Cleanup
    for _, node := range nodes {
        node.Stop()
    }
}
```

#### Property-based Testing

```go
func TestConsistencyProperties(t *testing.T) {
    quick.Check(func(commands []Command) bool {
        // Run commands on system
        // Verify consistency properties
        return verifyLinearizability(commands)
    }, nil)
}
```

### Debugging and Monitoring

#### Logging Strategy

```go
import "log/slog"

type Logger struct {
    *slog.Logger
    nodeID NodeID
}

func NewLogger(nodeID NodeID) *Logger {
    return &Logger{
        Logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
        nodeID: nodeID,
    }
}

func (l *Logger) LogMessage(msg Message, action string) {
    l.Info("message_handled",
        "node_id", l.nodeID,
        "message_type", msg.GetType(),
        "action", action,
        "from", msg.GetFrom(),
        "to", msg.GetTo(),
    )
}
```

#### Metrics Collection

```go
type Metrics struct {
    MessagesSent     int64
    MessagesReceived int64
    CommandsProcessed int64
    Errors          int64
    mu              sync.RWMutex
}

func (m *Metrics) IncrementMessagesSent() {
    atomic.AddInt64(&m.MessagesSent, 1)
}

func (m *Metrics) GetStats() map[string]int64 {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    return map[string]int64{
        "messages_sent":     atomic.LoadInt64(&m.MessagesSent),
        "messages_received": atomic.LoadInt64(&m.MessagesReceived),
        "commands_processed": atomic.LoadInt64(&m.CommandsProcessed),
        "errors":           atomic.LoadInt64(&m.Errors),
    }
}
```

### RPC Semantics and Reliability

#### At-Least-Once Semantics

```go
func (c *Client) sendWithRetry(msg Message, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        if err := c.send(msg); err == nil {
            return nil
        }
        time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
    }
    return errors.New("max retries exceeded")
}
```

#### At-Most-Once Semantics

```go
type RequestTracker struct {
    processed map[int64]bool
    mu        sync.RWMutex
}

func (rt *RequestTracker) IsProcessed(seq int64) bool {
    rt.mu.RLock()
    defer rt.mu.RUnlock()
    return rt.processed[seq]
}

func (rt *RequestTracker) MarkProcessed(seq int64) {
    rt.mu.Lock()
    defer rt.mu.Unlock()
    rt.processed[seq] = true
}
```

#### Exactly-Once Semantics

```go
func (s *Server) handleCommandWithDeduplication(msg *CommandMessage) error {
    if s.requestTracker.IsProcessed(msg.Seq) {
        // Return cached result
        return s.sendCachedResult(msg)
    }
    
    // Process command
    result, err := s.app.Execute(msg.Command)
    if err != nil {
        return err
    }
    
    // Cache result and mark as processed
    s.cacheResult(msg.Seq, result)
    s.requestTracker.MarkProcessed(msg.Seq)
    
    return s.sendResult(msg, result)
}
```

### Best Practices and Common Pitfalls

#### Do's:
- Use context.Context for cancellation and timeouts
- Implement proper error handling and retry logic
- Use structured logging with correlation IDs
- Design for failure - assume network partitions and node failures
- Implement health checks and monitoring
- Use connection pooling for network communication
- Implement backpressure mechanisms
- Test with network partitions and message drops

#### Don'ts:
- Don't ignore goroutine leaks - always clean up resources
- Don't use global variables for state management
- Don't block on network calls without timeouts
- Don't assume messages will be delivered in order
- Don't ignore clock skew in distributed systems
- Don't forget to handle split-brain scenarios
- Don't use blocking operations in message handlers

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with race detection
go test -race ./...

# Run tests with coverage
go test -cover ./...

# Run specific test with verbose output
go test -v -run TestKVStore ./...

# Run benchmarks
go test -bench=. ./...

# Run tests with timeout
go test -timeout 30s ./...
```

### Project Structure

```
distributed-systems/
├── cmd/                    # Main applications
│   ├── server/            # Server binary
│   └── client/            # Client binary
├── internal/              # Private application code
│   ├── node/              # Node implementations
│   ├── network/           # Network simulation
│   ├── consensus/         # Consensus algorithms
│   └── storage/           # Storage engines
├── pkg/                   # Public library code
│   ├── types/             # Common types
│   ├── messages/          # Message definitions
│   └── utils/             # Utility functions
├── tests/                 # Integration tests
├── docs/                  # Documentation
└── scripts/               # Build and test scripts
```

### Next Steps

1. Implement the basic Node interface
2. Create a simple key-value store application
3. Add client-server communication
4. Implement basic consensus (Raft)
5. Add fault tolerance mechanisms
6. Implement transactions
7. Add comprehensive testing
8. Performance optimization and monitoring

This foundation will prepare you for building the complete distributed key-value store in subsequent labs.
