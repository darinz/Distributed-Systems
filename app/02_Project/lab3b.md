# Lab 3B: Multi-Paxos with Replicated Log

## Overview

This lab extends basic Paxos to implement Multi-Paxos, which maintains a replicated log for state machine replication. Instead of agreeing on a single value, Multi-Paxos allows a set of distributed nodes to agree on a sequence of values (commands) that form a replicated log, enabling all servers to execute client requests in the same order.

## From Basic Paxos to Multi-Paxos

### Basic Paxos Limitations

- **Single Value**: Basic Paxos only agrees on one value
- **Contention**: Multiple proposers can cause conflicts with higher sequence numbers
- **Inefficiency**: Each consensus requires 2 RTTs (prepare + accept phases)

### Multi-Paxos Solution

- **Replicated Log**: Maintain a log where each slot contains a Paxos instance
- **Leader Election**: Elect a single leader to reduce contention
- **Optimized Accept**: Leader can skip prepare phase for subsequent proposals
- **State Machine Replication**: Execute commands in log order

## Replicated Log Implementation

### Log Structure

```go
// Log entry representing a single slot
type LogEntry struct {
    Slot    int64       `json:"slot"`
    Value   PaxosValue  `json:"value"`
    Decided bool        `json:"decided"`
    Executed bool       `json:"executed"`
}

// Replicated log
type ReplicatedLog struct {
    mu      sync.RWMutex
    entries map[int64]*LogEntry
    nextSlot int64
}

// Paxos instance for each log slot
type PaxosInstance struct {
    mu              sync.RWMutex
    slot            int64
    state           InstanceState
    proposalNumber  ProposalNumber
    acceptedValue   PaxosValue
    decidedValue    PaxosValue
    decided         bool
}

type InstanceState int

const (
    StatePending InstanceState = iota
    StatePrepared
    StateAccepted
    StateDecided
)
```

### Log Management

```go
func (rl *ReplicatedLog) GetNextEmptySlot() int64 {
    rl.mu.RLock()
    defer rl.mu.RUnlock()
    
    return rl.nextSlot
}

func (rl *ReplicatedLog) SetValue(slot int64, value PaxosValue) {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    if entry, exists := rl.entries[slot]; exists {
        entry.Value = value
        entry.Decided = true
    } else {
        rl.entries[slot] = &LogEntry{
            Slot:    slot,
            Value:   value,
            Decided: true,
        }
    }
    
    // Update next slot if necessary
    if slot >= rl.nextSlot {
        rl.nextSlot = slot + 1
    }
}

func (rl *ReplicatedLog) GetEntry(slot int64) (*LogEntry, bool) {
    rl.mu.RLock()
    defer rl.mu.RUnlock()
    
    entry, exists := rl.entries[slot]
    return entry, exists
}

func (rl *ReplicatedLog) GetEntries(from, to int64) []*LogEntry {
    rl.mu.RLock()
    defer rl.mu.RUnlock()
    
    var entries []*LogEntry
    for slot := from; slot < to; slot++ {
        if entry, exists := rl.entries[slot]; exists {
            entries = append(entries, entry)
        }
    }
    return entries
}
```

## Multi-Paxos Protocol Implementation

### Multi-Paxos Node

```go
type MultiPaxosNode struct {
    mu              sync.RWMutex
    me              NodeID
    peers           []NodeID
    log             *ReplicatedLog
    instances       map[int64]*PaxosInstance
    isLeader        bool
    ballotNumber    ProposalNumber
    leaderBallot    ProposalNumber
    lastHeartbeat   time.Time
    heartbeatTimer  *time.Timer
    electionTimer   *time.Timer
    done            chan struct{}
}

func (mpn *MultiPaxosNode) Start() error {
    mpn.log = NewReplicatedLog()
    mpn.instances = make(map[int64]*PaxosInstance)
    mpn.done = make(chan struct{})
    
    // Start heartbeat timer
    mpn.startHeartbeatTimer()
    
    // Start election timer
    mpn.startElectionTimer()
    
    return nil
}

func (mpn *MultiPaxosNode) ProposeCommand(command PaxosValue) error {
    mpn.mu.Lock()
    defer mpn.mu.Unlock()
    
    // Only leader can propose
    if !mpn.isLeader {
        return errors.New("not the leader")
    }
    
    // Get next empty slot
    slot := mpn.log.GetNextEmptySlot()
    
    // Create or get Paxos instance for this slot
    instance := mpn.getOrCreateInstance(slot)
    
    // Start Paxos for this slot
    go mpn.startPaxosForSlot(instance, command)
    
    return nil
}
```

## Leader Election (Phase 1)

### Ballot Number Management

```go
type BallotNumber struct {
    Number int64
    NodeID NodeID
}

func (bn BallotNumber) Compare(other BallotNumber) int {
    if bn.Number < other.Number {
        return -1
    } else if bn.Number > other.Number {
        return 1
    }
    
    // Same number, compare by NodeID
    if bn.NodeID < other.NodeID {
        return -1
    } else if bn.NodeID > other.NodeID {
        return 1
    }
    return 0
}

func (mpn *MultiPaxosNode) generateBallotNumber() BallotNumber {
    mpn.mu.Lock()
    defer mpn.mu.Unlock()
    
    mpn.ballotNumber++
    return BallotNumber{
        Number: mpn.ballotNumber,
        NodeID: mpn.me,
    }
}
```

### Leader Election Implementation

```go
func (mpn *MultiPaxosNode) startElectionTimer() {
    mpn.electionTimer = time.NewTimer(mpn.getElectionTimeout())
    
    go func() {
        for {
            select {
            case <-mpn.electionTimer.C:
                mpn.attemptLeadership()
                mpn.electionTimer.Reset(mpn.getElectionTimeout())
            case <-mpn.done:
                return
            }
        }
    }()
}

func (mpn *MultiPaxosNode) attemptLeadership() {
    mpn.mu.Lock()
    defer mpn.mu.Unlock()
    
    if mpn.isLeader {
        return // Already leader
    }
    
    // Generate new ballot number
    ballot := mpn.generateBallotNumber()
    
    // Send prepare requests (Phase 1)
    go mpn.sendPrepareRequests(ballot)
}

func (mpn *MultiPaxosNode) sendPrepareRequests(ballot BallotNumber) {
    prepareReq := &PrepareRequest{
        PaxosMessage: PaxosMessage{
            Type:          "prepare",
            From:          mpn.me,
            ProposalNumber: ProposalNumber(ballot.Number),
        },
        Ballot: ballot,
    }
    
    var replies []*PrepareReply
    var mu sync.Mutex
    
    // Send to all peers
    for _, peer := range mpn.peers {
        go func(peer NodeID) {
            var reply PrepareReply
            err := mpn.call(peer, "MultiPaxosNode.Prepare", prepareReq, &reply)
            
            mu.Lock()
            defer mu.Unlock()
            
            if err == nil {
                replies = append(replies, &reply)
            }
        }(peer)
    }
    
    // Wait for majority
    time.Sleep(100 * time.Millisecond)
    
    mu.Lock()
    if len(replies) > len(mpn.peers)/2 {
        mpn.handlePrepareMajority(ballot, replies)
    }
    mu.Unlock()
}

func (mpn *MultiPaxosNode) handlePrepareMajority(ballot BallotNumber, replies []*PrepareReply) {
    mpn.mu.Lock()
    defer mpn.mu.Unlock()
    
    // Update log with accepted values from replies
    for _, reply := range replies {
        if reply.Accepted {
            for slot, value := range reply.AcceptedValues {
                mpn.log.SetValue(slot, value)
            }
        }
    }
    
    // Become leader
    mpn.isLeader = true
    mpn.leaderBallot = ballot
    mpn.lastHeartbeat = time.Now()
    
    log.Printf("Node %s: Became leader with ballot %d", mpn.me, ballot.Number)
    
    // Start sending heartbeats
    mpn.startHeartbeatTimer()
}
```

### Prepare Handler

```go
func (mpn *MultiPaxosNode) Prepare(req *PrepareRequest, reply *PrepareReply) error {
    mpn.mu.Lock()
    defer mpn.mu.Unlock()
    
    ballot := req.Ballot
    
    // Check if ballot number is higher than what we've seen
    if ballot.Compare(mpn.leaderBallot) > 0 {
        mpn.leaderBallot = ballot
        
        // Collect accepted values from our log
        acceptedValues := make(map[int64]PaxosValue)
        for slot, entry := range mpn.log.entries {
            if entry.Decided {
                acceptedValues[slot] = entry.Value
            }
        }
        
        reply.Accepted = true
        reply.AcceptedValues = acceptedValues
        
        log.Printf("Node %s: Accepted prepare for ballot %d", mpn.me, ballot.Number)
    } else {
        reply.Accepted = false
        log.Printf("Node %s: Rejected prepare for ballot %d (current: %d)", 
            mpn.me, ballot.Number, mpn.leaderBallot.Number)
    }
    
    return nil
}
```

## Optimized Accept Phase (Phase 2)

### Leader Accept Implementation

```go
func (mpn *MultiPaxosNode) startPaxosForSlot(instance *PaxosInstance, value PaxosValue) {
    mpn.mu.RLock()
    ballot := mpn.leaderBallot
    mpn.mu.RUnlock()
    
    // Skip prepare phase - go directly to accept
    acceptReq := &AcceptRequest{
        PaxosMessage: PaxosMessage{
            Type:          "accept",
            From:          mpn.me,
            InstanceID:    InstanceID(instance.slot),
            ProposalNumber: ProposalNumber(ballot.Number),
        },
        Ballot: ballot,
        Slot:   instance.slot,
        Value:  value,
    }
    
    // Send accept requests to all peers
    var replies []*AcceptReply
    var mu sync.Mutex
    
    for _, peer := range mpn.peers {
        go func(peer NodeID) {
            var reply AcceptReply
            err := mpn.call(peer, "MultiPaxosNode.Accept", acceptReq, &reply)
            
            mu.Lock()
            defer mu.Unlock()
            
            if err == nil {
                replies = append(replies, &reply)
            }
        }(peer)
    }
    
    // Wait for majority
    time.Sleep(100 * time.Millisecond)
    
    mu.Lock()
    if len(replies) > len(mpn.peers)/2 {
        mpn.handleAcceptMajority(instance, value, replies)
    }
    mu.Unlock()
}

func (mpn *MultiPaxosNode) handleAcceptMajority(instance *PaxosInstance, value PaxosValue, replies []*AcceptReply) {
    mpn.mu.Lock()
    defer mpn.mu.Unlock()
    
    // Count accepted replies
    acceptedCount := 0
    for _, reply := range replies {
        if reply.Accepted {
            acceptedCount++
        }
    }
    
    if acceptedCount > len(mpn.peers)/2 {
        // Value is decided
        instance.decided = true
        instance.decidedValue = value
        mpn.log.SetValue(instance.slot, value)
        
        // Notify learners
        go mpn.notifyLearners(instance.slot, value)
        
        log.Printf("Node %s: Decided value for slot %d", mpn.me, instance.slot)
    }
}
```

### Accept Handler

```go
func (mpn *MultiPaxosNode) Accept(req *AcceptRequest, reply *AcceptReply) error {
    mpn.mu.Lock()
    defer mpn.mu.Unlock()
    
    ballot := req.Ballot
    slot := req.Slot
    value := req.Value
    
    // Check if ballot number matches current leader
    if ballot.Compare(mpn.leaderBallot) == 0 {
        // Accept the value
        mpn.log.SetValue(slot, value)
        
        reply.Accepted = true
        
        log.Printf("Node %s: Accepted value for slot %d", mpn.me, slot)
    } else {
        reply.Accepted = false
        log.Printf("Node %s: Rejected accept for slot %d (ballot mismatch)", mpn.me, slot)
    }
    
    return nil
}
```

## Log Synchronization

### Heartbeat with Log Sync

```go
func (mpn *MultiPaxosNode) startHeartbeatTimer() {
    mpn.heartbeatTimer = time.NewTimer(50 * time.Millisecond)
    
    go func() {
        for {
            select {
            case <-mpn.heartbeatTimer.C:
                if mpn.isLeader {
                    mpn.sendHeartbeats()
                }
                mpn.heartbeatTimer.Reset(50 * time.Millisecond)
            case <-mpn.done:
                return
            }
        }
    }()
}

func (mpn *MultiPaxosNode) sendHeartbeats() {
    mpn.mu.RLock()
    ballot := mpn.leaderBallot
    logEntries := mpn.log.GetEntries(0, mpn.log.nextSlot)
    mpn.mu.RUnlock()
    
    heartbeat := &HeartbeatMessage{
        From:       mpn.me,
        Ballot:     ballot,
        LogEntries: logEntries,
        Timestamp:  time.Now(),
    }
    
    for _, peer := range mpn.peers {
        go func(peer NodeID) {
            var reply HeartbeatReply
            mpn.call(peer, "MultiPaxosNode.Heartbeat", heartbeat, &reply)
        }(peer)
    }
}

func (mpn *MultiPaxosNode) Heartbeat(req *HeartbeatMessage, reply *HeartbeatReply) error {
    mpn.mu.Lock()
    defer mpn.mu.Unlock()
    
    // Update leader ballot if higher
    if req.Ballot.Compare(mpn.leaderBallot) >= 0 {
        mpn.leaderBallot = req.Ballot
        mpn.isLeader = false // We're not the leader
        mpn.lastHeartbeat = time.Now()
        
        // Sync log entries
        for _, entry := range req.LogEntries {
            if entry.Decided {
                mpn.log.SetValue(entry.Slot, entry.Value)
            }
        }
        
        reply.Accepted = true
    } else {
        reply.Accepted = false
    }
    
    return nil
}
```

## Garbage Collection

### Garbage Collection Implementation

```go
type GarbageCollector struct {
    mu              sync.RWMutex
    executedSlots   map[NodeID]int64
    minExecutedSlot int64
}

func (gc *GarbageCollector) UpdateExecutedSlot(nodeID NodeID, slot int64) {
    gc.mu.Lock()
    defer gc.mu.Unlock()
    
    gc.executedSlots[nodeID] = slot
    gc.updateMinExecutedSlot()
}

func (gc *GarbageCollector) updateMinExecutedSlot() {
    if len(gc.executedSlots) == 0 {
        return
    }
    
    minSlot := int64(^uint64(0) >> 1) // Max int64
    for _, slot := range gc.executedSlots {
        if slot < minSlot {
            minSlot = slot
        }
    }
    
    gc.minExecutedSlot = minSlot
}

func (gc *GarbageCollector) GetMinExecutedSlot() int64 {
    gc.mu.RLock()
    defer gc.mu.RUnlock()
    return gc.minExecutedSlot
}

func (mpn *MultiPaxosNode) GarbageCollect() {
    mpn.mu.Lock()
    defer mpn.mu.Unlock()
    
    minExecuted := mpn.gc.GetMinExecutedSlot()
    
    // Remove log entries that are no longer needed
    for slot := range mpn.log.entries {
        if slot < minExecuted {
            delete(mpn.log.entries, slot)
            delete(mpn.instances, slot)
        }
    }
    
    log.Printf("Node %s: Garbage collected up to slot %d", mpn.me, minExecuted)
}
```

## Hole Handling

### Handling Gaps in Log

```go
func (mpn *MultiPaxosNode) FillHoles() {
    mpn.mu.RLock()
    nextSlot := mpn.log.nextSlot
    mpn.mu.RUnlock()
    
    // Check for holes in the log
    for slot := int64(0); slot < nextSlot; slot++ {
        if entry, exists := mpn.log.GetEntry(slot); !exists || !entry.Decided {
            // Found a hole, try to fill it
            go mpn.fillHole(slot)
        }
    }
}

func (mpn *MultiPaxosNode) fillHole(slot int64) {
    // Try to learn the value for this slot from other nodes
    learnReq := &LearnRequest{
        Slot: slot,
        From: mpn.me,
    }
    
    var replies []*LearnReply
    var mu sync.Mutex
    
    for _, peer := range mpn.peers {
        go func(peer NodeID) {
            var reply LearnReply
            err := mpn.call(peer, "MultiPaxosNode.Learn", learnReq, &reply)
            
            mu.Lock()
            defer mu.Unlock()
            
            if err == nil && reply.HasValue {
                replies = append(replies, &reply)
            }
        }(peer)
    }
    
    // Wait for replies
    time.Sleep(100 * time.Millisecond)
    
    mu.Lock()
    if len(replies) > 0 {
        // Use the first reply (all should be the same)
        reply := replies[0]
        mpn.log.SetValue(slot, reply.Value)
        log.Printf("Node %s: Filled hole at slot %d", mpn.me, slot)
    }
    mu.Unlock()
}
```

## Testing Strategy

### Unit Testing

```go
func TestReplicatedLog(t *testing.T) {
    log := NewReplicatedLog()
    
    // Test setting values
    value1 := &TestValue{Data: "value1"}
    log.SetValue(0, value1)
    
    entry, exists := log.GetEntry(0)
    assert.True(t, exists)
    assert.Equal(t, value1, entry.Value)
    assert.True(t, entry.Decided)
    
    // Test next slot
    assert.Equal(t, int64(1), log.GetNextEmptySlot())
}

func TestMultiPaxosNode(t *testing.T) {
    node := NewMultiPaxosNode("node1", []NodeID{"node2", "node3"})
    go node.Start()
    defer node.Stop()
    
    // Test command proposal
    command := &TestCommand{Op: "put", Key: "test", Value: "value"}
    err := node.ProposeCommand(command)
    assert.NoError(t, err)
    
    // Wait for decision
    time.Sleep(200 * time.Millisecond)
    
    // Verify command was decided
    entry, exists := node.log.GetEntry(0)
    assert.True(t, exists)
    assert.True(t, entry.Decided)
}
```

### Integration Testing

```go
func TestMultiPaxosConsensus(t *testing.T) {
    // Create 3 Multi-Paxos nodes
    nodes := make([]*MultiPaxosNode, 3)
    for i := 0; i < 3; i++ {
        nodeID := NodeID(fmt.Sprintf("node%d", i))
        peers := make([]NodeID, 0, 2)
        for j := 0; j < 3; j++ {
            if i != j {
                peers = append(peers, NodeID(fmt.Sprintf("node%d", j)))
            }
        }
        
        nodes[i] = NewMultiPaxosNode(nodeID, peers)
        go nodes[i].Start()
    }
    
    // Wait for leader election
    time.Sleep(300 * time.Millisecond)
    
    // Find leader
    var leader *MultiPaxosNode
    for _, node := range nodes {
        if node.isLeader {
            leader = node
            break
        }
    }
    assert.NotNil(t, leader, "No leader elected")
    
    // Propose commands
    commands := []*TestCommand{
        {Op: "put", Key: "a", Value: "1"},
        {Op: "put", Key: "b", Value: "2"},
        {Op: "append", Key: "a", Value: "3"},
    }
    
    for _, cmd := range commands {
        err := leader.ProposeCommand(cmd)
        assert.NoError(t, err)
    }
    
    // Wait for consensus
    time.Sleep(500 * time.Millisecond)
    
    // Verify all nodes have the same log
    for i, node := range nodes {
        for j, expectedCmd := range commands {
            entry, exists := node.log.GetEntry(int64(j))
            assert.True(t, exists, "Node %d missing entry %d", i, j)
            assert.True(t, entry.Decided, "Node %d entry %d not decided", i, j)
            assert.Equal(t, expectedCmd, entry.Value, "Node %d entry %d mismatch", i, j)
        }
    }
}
```

## Best Practices and Common Pitfalls

### Do's

1. **Use consistent ballot numbers**: Ensure proper leader election
2. **Handle log holes**: Fill gaps in the replicated log
3. **Implement garbage collection**: Prevent memory growth
4. **Use heartbeats**: Keep followers informed of leader status
5. **Optimize accept phase**: Skip prepare phase for leaders
6. **Handle concurrent proposals**: Multiple commands can be proposed
7. **Sync logs efficiently**: Use heartbeats for log synchronization
8. **Test with failures**: Simulate leader crashes and network partitions

### Don'ts

1. **Don't ignore ballot numbers**: Always check ballot validity
2. **Don't forget to handle holes**: Log gaps must be filled
3. **Don't skip garbage collection**: Memory will grow indefinitely
4. **Don't assume single leader**: Multiple nodes can attempt leadership
5. **Don't ignore log synchronization**: Followers must stay in sync
6. **Don't forget to handle timeouts**: Leaders can fail
7. **Don't assume ordered execution**: Commands may be decided out of order
8. **Don't ignore majority requirements**: Need majority for both phases

### Common Implementation Mistakes

```go
// WRONG: Not checking ballot numbers
func (mpn *MultiPaxosNode) Accept(req *AcceptRequest, reply *AcceptReply) error {
    // Missing ballot check!
    mpn.log.SetValue(req.Slot, req.Value)
    reply.Accepted = true
    return nil
}

// CORRECT: Always check ballot numbers
func (mpn *MultiPaxosNode) Accept(req *AcceptRequest, reply *AcceptReply) error {
    if req.Ballot.Compare(mpn.leaderBallot) == 0 {
        mpn.log.SetValue(req.Slot, req.Value)
        reply.Accepted = true
    } else {
        reply.Accepted = false
    }
    return nil
}

// WRONG: Not handling log holes
func (mpn *MultiPaxosNode) ExecuteCommands() {
    for slot := int64(0); slot < mpn.log.nextSlot; slot++ {
        // Missing hole check!
        entry, _ := mpn.log.GetEntry(slot)
        mpn.executeCommand(entry.Value)
    }
}

// CORRECT: Handle holes in log
func (mpn *MultiPaxosNode) ExecuteCommands() {
    for slot := int64(0); slot < mpn.log.nextSlot; slot++ {
        if entry, exists := mpn.log.GetEntry(slot); exists && entry.Decided {
            mpn.executeCommand(entry.Value)
        } else {
            // Found hole, try to fill it
            mpn.fillHole(slot)
            break
        }
    }
}
```

### Performance Considerations

1. **Minimize message overhead**: Use heartbeats for log sync
2. **Optimize accept phase**: Skip prepare phase for leaders
3. **Implement efficient garbage collection**: Regular cleanup
4. **Use connection pooling**: Reuse network connections
5. **Monitor log growth**: Track memory usage

### Debugging Tips

1. **Add structured logging**: Include ballot numbers and slot numbers
2. **Use consistent log levels**: Debug, Info, Warn, Error
3. **Include correlation IDs**: Track commands across nodes
4. **Monitor leader election**: Log leadership changes
5. **Test with delays**: Add artificial delays to find race conditions

### Running Tests

```bash
# Run Multi-Paxos tests
go test -v -run TestMultiPaxos ./...

# Run with race detection
go test -race -run TestMultiPaxos ./...

# Run specific test
go test -v -run TestMultiPaxosConsensus ./...

# Run with coverage
go test -cover -run TestMultiPaxos ./...

# Run benchmarks
go test -bench=. -run TestMultiPaxos ./...
```

### Next Steps

1. Implement basic Multi-Paxos with replicated log
2. Add leader election with ballot numbers
3. Implement optimized accept phase
4. Add log synchronization and hole handling
5. Implement garbage collection
6. Add comprehensive testing
7. Performance optimization

This foundation will prepare you for implementing a complete, fault-tolerant replicated state machine using Multi-Paxos.
