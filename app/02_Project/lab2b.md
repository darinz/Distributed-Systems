# Lab 2B: Complete Primary-Backup System with Lamport Clocks

## Overview

This lab completes the primary-backup replication system by implementing state transfer mechanisms, comprehensive message handling, and logical clocks (Lamport and Vector clocks) for causal ordering in distributed systems.

## Progress Check

- **Lab 2A Status**: ViewService implementation should be complete
- **Lab 2B Focus**: Complete primary-backup system with state transfer and logical clocks
- **Complexity**: This part is significantly more challenging than ViewService

## Ping Timeout System

### PingCheckTimeout Explanation

The ping timeout system uses a deterministic approach to detect server failures:

```
Time:     0ms    100ms   200ms
S1:       ping    ping    ping
S2:       ping    [miss]  ping
Result:   alive   alive   alive (S2 missed one cycle but recovered)
```

**Key Points:**
- Servers are considered alive if they ping within the timeout window
- Multiple missed pings in a row mark a server as dead
- The system is deterministic (no timestamps) for search test compatibility

### Implementation

```go
type PingTimeoutManager struct {
    mu           sync.RWMutex
    servers      map[NodeID]*ServerPingState
    pingTimeout  int64 // Number of ping cycles before marking as dead
    currentCycle int64
}

type ServerPingState struct {
    lastPingCycle int64
    isAlive       bool
    missedPings   int64
}

func (ptm *PingTimeoutManager) UpdatePing(serverID NodeID) {
    ptm.mu.Lock()
    defer ptm.mu.Unlock()
    
    if state, exists := ptm.servers[serverID]; exists {
        state.lastPingCycle = ptm.currentCycle
        state.missedPings = 0
        state.isAlive = true
    } else {
        ptm.servers[serverID] = &ServerPingState{
            lastPingCycle: ptm.currentCycle,
            isAlive:       true,
            missedPings:   0,
        }
    }
}

func (ptm *PingTimeoutManager) CheckTimeouts() {
    ptm.mu.Lock()
    defer ptm.mu.Unlock()
    
    ptm.currentCycle++
    
    for serverID, state := range ptm.servers {
        if ptm.currentCycle-state.lastPingCycle > ptm.pingTimeout {
            state.missedPings++
            if state.missedPings >= ptm.pingTimeout {
                state.isAlive = false
                log.Printf("Server %s marked as dead after %d missed pings", 
                    serverID, state.missedPings)
            }
        }
    }
}
```

## Complete Primary-Backup System

### System Architecture

```
Client → Primary → ViewService
         ↓
       Backup
```

### Message Design

#### Core Message Types

```go
// Base message with logical clock
type BaseMessage struct {
    Type      string    `json:"type"`
    From      NodeID    `json:"from"`
    To        NodeID    `json:"to"`
    View      View      `json:"view"`
    Timestamp int64     `json:"timestamp"` // Lamport timestamp
    VectorClock []int64 `json:"vector_clock,omitempty"` // Vector clock
}

// Client operations
type PutArgs struct {
    BaseMessage
    Key   string `json:"key"`
    Value string `json:"value"`
    Seq   int64  `json:"seq"` // Client sequence number
}

type PutReply struct {
    BaseMessage
    Err Err `json:"err"`
}

type GetArgs struct {
    BaseMessage
    Key string `json:"key"`
    Seq int64  `json:"seq"`
}

type GetReply struct {
    BaseMessage
    Value string `json:"value"`
    Err   Err    `json:"err"`
}

type AppendArgs struct {
    BaseMessage
    Key   string `json:"key"`
    Value string `json:"value"`
    Seq   int64  `json:"seq"`
}

type AppendReply struct {
    BaseMessage
    Err Err `json:"err"`
}

// State transfer messages
type StateTransferArgs struct {
    BaseMessage
    AppState *AMOApplication `json:"app_state"`
    ViewNum  int64          `json:"view_num"`
}

type StateTransferReply struct {
    BaseMessage
    Err Err `json:"err"`
}

// Backup operation forwarding
type BackupOpArgs struct {
    BaseMessage
    OpType string      `json:"op_type"`
    OpData interface{} `json:"op_data"`
    Seq    int64       `json:"seq"`
}

type BackupOpReply struct {
    BaseMessage
    Result interface{} `json:"result"`
    Err    Err         `json:"err"`
}
```

### State Management

#### PBClient State

```go
type PBClient struct {
    mu          sync.RWMutex
    me          NodeID
    viewServer  NodeID
    currentView View
    seq         int64
    pending     map[int64]chan interface{}
    lamportClock int64
    vectorClock  []int64
}

func (c *PBClient) Put(key, value string) error {
    c.mu.Lock()
    seq := c.seq
    c.seq++
    c.lamportClock++
    
    // Create result channel
    resultChan := make(chan interface{}, 1)
    c.pending[seq] = resultChan
    c.mu.Unlock()
    
    // Get current view
    view, err := c.getCurrentView()
    if err != nil {
        return err
    }
    
    // Create message with logical clock
    args := &PutArgs{
        BaseMessage: BaseMessage{
            Type:      "put",
            From:      c.me,
            To:        view.Primary,
            View:      view,
            Timestamp: c.lamportClock,
            VectorClock: c.vectorClock,
        },
        Key:   key,
        Value: value,
        Seq:   seq,
    }
    
    // Send to primary
    var reply PutReply
    err = c.call(view.Primary, "PBServer.Put", args, &reply)
    if err != nil {
        return err
    }
    
    // Update logical clock
    c.updateLogicalClock(reply.Timestamp, reply.VectorClock)
    
    return reply.Err
}

func (c *PBClient) updateLogicalClock(timestamp int64, vectorClock []int64) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    // Update Lamport clock
    if timestamp > c.lamportClock {
        c.lamportClock = timestamp
    }
    c.lamportClock++
    
    // Update Vector clock
    if len(vectorClock) > len(c.vectorClock) {
        c.vectorClock = make([]int64, len(vectorClock))
    }
    for i, v := range vectorClock {
        if i < len(c.vectorClock) && v > c.vectorClock[i] {
            c.vectorClock[i] = v
        }
    }
}
```

#### PBServer State

```go
type PBServer struct {
    mu                    sync.RWMutex
    me                    NodeID
    viewServer            NodeID
    currentView           View
    app                   *AMOApplication
    stateTransferInProgress bool
    lamportClock          int64
    vectorClock           []int64
    processedSeqs         map[int64]bool // For deduplication
}

func (s *PBServer) Put(args *PutArgs, reply *PutReply) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    // Update logical clock
    s.updateLogicalClock(args.Timestamp, args.VectorClock)
    
    // Check if we're the primary
    if s.currentView.Primary != s.me {
        reply.Err = ErrWrongServer
        s.setReplyLogicalClock(reply)
        return nil
    }
    
    // Check view number
    if args.View.ViewNum != s.currentView.ViewNum {
        reply.Err = ErrWrongView
        s.setReplyLogicalClock(reply)
        return nil
    }
    
    // Check for duplicate request
    if s.processedSeqs[args.Seq] {
        reply.Err = OK // Already processed
        s.setReplyLogicalClock(reply)
        return nil
    }
    
    // Don't process during state transfer
    if s.stateTransferInProgress {
        reply.Err = ErrStateTransfer
        s.setReplyLogicalClock(reply)
        return nil
    }
    
    // Forward to backup if exists
    if s.currentView.Backup != "" {
        err := s.forwardToBackup(args)
        if err != nil {
            reply.Err = ErrBackupFailed
            s.setReplyLogicalClock(reply)
            return nil
        }
    }
    
    // Execute operation
    err := s.app.Put(args.Key, args.Value)
    if err != nil {
        reply.Err = err
    } else {
        reply.Err = OK
        s.processedSeqs[args.Seq] = true
    }
    
    s.setReplyLogicalClock(reply)
    return nil
}

func (s *PBServer) updateLogicalClock(timestamp int64, vectorClock []int64) {
    // Update Lamport clock
    if timestamp > s.lamportClock {
        s.lamportClock = timestamp
    }
    s.lamportClock++
    
    // Update Vector clock
    if len(vectorClock) > len(s.vectorClock) {
        s.vectorClock = make([]int64, len(vectorClock))
    }
    for i, v := range vectorClock {
        if i < len(s.vectorClock) && v > s.vectorClock[i] {
            s.vectorClock[i] = v
        }
    }
}

func (s *PBServer) setReplyLogicalClock(reply interface{}) {
    switch r := reply.(type) {
    case *PutReply:
        r.Timestamp = s.lamportClock
        r.VectorClock = make([]int64, len(s.vectorClock))
        copy(r.VectorClock, s.vectorClock)
    case *GetReply:
        r.Timestamp = s.lamportClock
        r.VectorClock = make([]int64, len(s.vectorClock))
        copy(r.VectorClock, s.vectorClock)
    case *AppendReply:
        r.Timestamp = s.lamportClock
        r.VectorClock = make([]int64, len(s.vectorClock))
        copy(r.VectorClock, s.vectorClock)
    }
}
```

## State Transfer Implementation

### State Transfer Process

```go
func (s *PBServer) initiateStateTransfer(backupID NodeID) {
    s.mu.Lock()
    s.stateTransferInProgress = true
    s.mu.Unlock()
    
    // Create state transfer message
    args := &StateTransferArgs{
        BaseMessage: BaseMessage{
            Type:      "state_transfer",
            From:      s.me,
            To:        backupID,
            View:      s.currentView,
            Timestamp: s.lamportClock,
            VectorClock: s.vectorClock,
        },
        AppState: s.app,
        ViewNum:  s.currentView.ViewNum,
    }
    
    var reply StateTransferReply
    err := s.call(backupID, "PBServer.StateTransfer", args, &reply)
    
    s.mu.Lock()
    s.stateTransferInProgress = false
    s.mu.Unlock()
    
    if err != nil || reply.Err != OK {
        log.Printf("State transfer failed: %v", err)
    } else {
        log.Printf("State transfer completed successfully")
    }
}

func (s *PBServer) StateTransfer(args *StateTransferArgs, reply *StateTransferReply) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    // Update logical clock
    s.updateLogicalClock(args.Timestamp, args.VectorClock)
    
    // Check if we're the backup
    if s.currentView.Backup != s.me {
        reply.Err = ErrWrongServer
        s.setReplyLogicalClock(reply)
        return nil
    }
    
    // Check view number
    if args.ViewNum != s.currentView.ViewNum {
        reply.Err = ErrWrongView
        s.setReplyLogicalClock(reply)
        return nil
    }
    
    // Replace application state
    s.app = args.AppState
    
    reply.Err = OK
    s.setReplyLogicalClock(reply)
    return nil
}
```

### State Transfer Rules

1. **Include all data**: Send entire AMOApplication state
2. **Block operations**: Primary drops requests during state transfer
3. **Handle duplicates**: Backup can receive duplicate state transfer messages
4. **View consistency**: Only overwrite state once per view change
5. **Ping behavior**: Primary pings reflect old view during state transfer
6. **View transition**: Only move to new view after state transfer completes

## Logical Clocks Implementation

### Lamport Clocks

```go
type LamportClock struct {
    timestamp int64
    mu        sync.Mutex
}

func (lc *LamportClock) Tick() int64 {
    lc.mu.Lock()
    defer lc.mu.Unlock()
    lc.timestamp++
    return lc.timestamp
}

func (lc *LamportClock) Update(receivedTimestamp int64) int64 {
    lc.mu.Lock()
    defer lc.mu.Unlock()
    if receivedTimestamp > lc.timestamp {
        lc.timestamp = receivedTimestamp
    }
    lc.timestamp++
    return lc.timestamp
}

func (lc *LamportClock) GetTimestamp() int64 {
    lc.mu.Lock()
    defer lc.mu.Unlock()
    return lc.timestamp
}
```

### Vector Clocks

```go
type VectorClock struct {
    clock []int64
    mu    sync.Mutex
}

func NewVectorClock(nodeCount int) *VectorClock {
    return &VectorClock{
        clock: make([]int64, nodeCount),
    }
}

func (vc *VectorClock) Tick(nodeID int) int64 {
    vc.mu.Lock()
    defer vc.mu.Unlock()
    if nodeID < len(vc.clock) {
        vc.clock[nodeID]++
        return vc.clock[nodeID]
    }
    return 0
}

func (vc *VectorClock) Update(receivedClock []int64, nodeID int) {
    vc.mu.Lock()
    defer vc.mu.Unlock()
    
    // Ensure clock is large enough
    if len(receivedClock) > len(vc.clock) {
        newClock := make([]int64, len(receivedClock))
        copy(newClock, vc.clock)
        vc.clock = newClock
    }
    
    // Update clock: max(our_clock[i], received_clock[i]) for all i
    for i := range vc.clock {
        if i < len(receivedClock) && receivedClock[i] > vc.clock[i] {
            vc.clock[i] = receivedClock[i]
        }
    }
    
    // Increment our own clock
    if nodeID < len(vc.clock) {
        vc.clock[nodeID]++
    }
}

func (vc *VectorClock) GetClock() []int64 {
    vc.mu.Lock()
    defer vc.mu.Unlock()
    result := make([]int64, len(vc.clock))
    copy(result, vc.clock)
    return result
}

func (vc *VectorClock) HappensBefore(other []int64) bool {
    vc.mu.Lock()
    defer vc.mu.Unlock()
    
    if len(vc.clock) != len(other) {
        return false
    }
    
    // Check if vc.clock < other (happens before)
    for i := range vc.clock {
        if vc.clock[i] > other[i] {
            return false
        }
    }
    
    // Check if they're not equal
    for i := range vc.clock {
        if vc.clock[i] < other[i] {
            return true
        }
    }
    
    return false
}
```

### Logical Clock Usage Example

```go
// Lamport Clock Example
func (s *PBServer) handleMessage(msg Message) {
    // Update Lamport clock on receive
    s.lamportClock = s.lamportClock.Update(msg.GetTimestamp())
    
    // Process message
    s.processMessage(msg)
    
    // Increment clock on send
    reply := s.createReply(msg)
    reply.SetTimestamp(s.lamportClock.Tick())
    
    s.sendReply(reply)
}

// Vector Clock Example
func (s *PBServer) handleMessageWithVectorClock(msg Message) {
    // Update Vector clock on receive
    s.vectorClock.Update(msg.GetVectorClock(), s.nodeID)
    
    // Process message
    s.processMessage(msg)
    
    // Increment our clock on send
    reply := s.createReply(msg)
    reply.SetVectorClock(s.vectorClock.GetClock())
    
    s.sendReply(reply)
}
```

## Testing Strategy

### Unit Testing

```go
func TestLamportClock(t *testing.T) {
    clock := &LamportClock{}
    
    // Test initial tick
    assert.Equal(t, int64(1), clock.Tick())
    assert.Equal(t, int64(2), clock.Tick())
    
    // Test update
    assert.Equal(t, int64(5), clock.Update(4))
    assert.Equal(t, int64(6), clock.Tick())
}

func TestVectorClock(t *testing.T) {
    clock := NewVectorClock(3)
    
    // Test tick
    assert.Equal(t, int64(1), clock.Tick(0))
    assert.Equal(t, int64(1), clock.Tick(1))
    
    // Test update
    received := []int64{2, 1, 0}
    clock.Update(received, 0)
    assert.Equal(t, int64(3), clock.GetClock()[0])
    assert.Equal(t, int64(1), clock.GetClock()[1])
}
```

### Integration Testing

```go
func TestStateTransfer(t *testing.T) {
    // Create ViewServer
    vs := NewViewServer()
    go vs.Start()
    defer vs.Stop()
    
    // Create servers
    primary := NewPBServer("primary", vs.GetNodeID())
    backup := NewPBServer("backup", vs.GetNodeID())
    
    go primary.Start()
    go backup.Start()
    
    // Wait for view establishment
    time.Sleep(200 * time.Millisecond)
    
    // Test state transfer
    testStateTransfer(t, primary, backup)
}

func testStateTransfer(t *testing.T, primary, backup *PBServer) {
    // Add some data to primary
    primary.Put("key1", "value1")
    primary.Put("key2", "value2")
    
    // Trigger state transfer
    primary.initiateStateTransfer("backup")
    
    // Verify backup has the data
    value, err := backup.Get("key1")
    assert.NoError(t, err)
    assert.Equal(t, "value1", value)
}
```

## Best Practices and Common Pitfalls

### Do's

1. **Always update logical clocks**: On send and receive
2. **Include clocks in all messages**: For proper causal ordering
3. **Handle state transfer carefully**: Block operations during transfer
4. **Validate view numbers**: Check view consistency on every message
5. **Implement deduplication**: Use sequence numbers to avoid duplicate processing
6. **Use proper synchronization**: Protect shared state with mutexes
7. **Test with network partitions**: Simulate message drops and delays
8. **Log important events**: Include logical clock values in logs

### Don'ts

1. **Don't forget to update clocks**: Always increment on send, update on receive
2. **Don't process during state transfer**: Block until transfer completes
3. **Don't ignore view number mismatches**: Always validate view numbers
4. **Don't forget to forward to backup**: Primary must wait for backup acknowledgment
5. **Don't use blocking operations in handlers**: Use channels and goroutines
6. **Don't assume message ordering**: Messages may arrive out of order
7. **Don't forget to handle duplicates**: Use sequence numbers for deduplication
8. **Don't ignore clock synchronization**: Keep clocks consistent across nodes

### Common Implementation Mistakes

```go
// WRONG: Not updating logical clock
func (s *PBServer) Put(args *PutArgs, reply *PutReply) error {
    // Missing clock update!
    err := s.app.Put(args.Key, args.Value)
    reply.Err = err
    return nil
}

// CORRECT: Always update logical clock
func (s *PBServer) Put(args *PutArgs, reply *PutReply) error {
    s.updateLogicalClock(args.Timestamp, args.VectorClock)
    err := s.app.Put(args.Key, args.Value)
    reply.Err = err
    s.setReplyLogicalClock(reply)
    return nil
}

// WRONG: Not handling state transfer
func (s *PBServer) Get(args *GetArgs, reply *GetReply) error {
    // Missing state transfer check!
    value, err := s.app.Get(args.Key)
    reply.Value = value
    reply.Err = err
    return nil
}

// CORRECT: Check state transfer status
func (s *PBServer) Get(args *GetArgs, reply *GetReply) error {
    if s.stateTransferInProgress {
        reply.Err = ErrStateTransfer
        return nil
    }
    // ... rest of implementation
}
```

### Performance Considerations

1. **Minimize lock contention**: Use RWMutex for read-heavy operations
2. **Batch operations**: Group related operations together
3. **Use connection pooling**: Reuse network connections
4. **Implement backpressure**: Handle overload gracefully
5. **Monitor resource usage**: Track memory and CPU usage

### Debugging Tips

1. **Add structured logging**: Include logical clock values
2. **Use consistent log levels**: Debug, Info, Warn, Error
3. **Include correlation IDs**: Track requests across servers
4. **Monitor state transfers**: Log all state transfer events
5. **Test with delays**: Add artificial delays to find race conditions

### Running Tests

```bash
# Run all Lab 2B tests
go test -v -run TestLab2B ./...

# Run with race detection
go test -race -run TestLab2B ./...

# Run specific test
go test -v -run TestStateTransfer ./...

# Run with coverage
go test -cover -run TestLab2B ./...

# Run benchmarks
go test -bench=. -run TestLab2B ./...
```

### Next Steps

1. Implement complete primary-backup system with state transfer
2. Add Lamport and Vector clock support
3. Implement comprehensive message handling
4. Add deduplication and error handling
5. Implement comprehensive testing
6. Add monitoring and logging
7. Performance optimization

This foundation will prepare you for implementing a complete, fault-tolerant primary-backup system with logical clocks and proper causal ordering.
