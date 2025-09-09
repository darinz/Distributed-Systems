# Lab 3A: Paxos Consensus Algorithm

## Overview

This lab implements the Paxos consensus algorithm, one of the most fundamental algorithms in distributed systems. Paxos allows a set of distributed nodes to agree on a single value even in the presence of failures, ensuring that all servers execute client requests in the same order (linearizability).

## Motivation

### Consensus Requirements

- **Consensus**: Decide on a single value among multiple distributed nodes
- **Fault Tolerance**: Continue operation if a majority of nodes are alive
- **Asynchronous Operation**: Work without synchronized clocks
- **Linearizability**: All servers must execute all client requests in the same order
- **State Machine Replication**: Each node is a replica of a state machine (key/value store)

### Why Paxos?

Paxos provides a mathematically proven solution to the consensus problem that:
- Guarantees safety (no two different values are chosen)
- Ensures liveness (a value will eventually be chosen if a majority is alive)
- Works in asynchronous networks with message delays and failures

## Paxos Algorithm Overview

### Core Components

1. **Proposers**: Initiate consensus by proposing values
2. **Acceptors**: Vote on proposed values
3. **Learners**: Learn the chosen value

### Two-Phase Protocol

**Phase 1 (Prepare)**: Proposer asks acceptors to promise not to accept proposals with lower numbers
**Phase 2 (Accept)**: Proposer asks acceptors to accept a specific proposal

## Go Implementation Framework

### Basic Types and Interfaces

```go
// Proposal number type
type ProposalNumber int64

// Paxos value type
type PaxosValue interface {
    Serialize() ([]byte, error)
    Deserialize(data []byte) error
}

// Paxos instance identifier
type InstanceID int64

// Paxos participant roles
type PaxosRole int

const (
    Proposer PaxosRole = iota
    Acceptor
    Learner
)

// Paxos node interface
type PaxosNode interface {
    Start() error
    Stop() error
    Propose(instance InstanceID, value PaxosValue) error
    Learn(instance InstanceID, value PaxosValue) error
}
```

### Message Types

```go
// Base Paxos message
type PaxosMessage struct {
    Type       string        `json:"type"`
    From       NodeID        `json:"from"`
    To         NodeID        `json:"to"`
    InstanceID InstanceID    `json:"instance_id"`
    ProposalNumber ProposalNumber `json:"proposal_number"`
}

// Prepare request (Phase 1)
type PrepareRequest struct {
    PaxosMessage
}

type PrepareReply struct {
    PaxosMessage
    Accepted bool        `json:"accepted"`
    Value    PaxosValue  `json:"value,omitempty"`
    Err      error       `json:"err,omitempty"`
}

// Accept request (Phase 2)
type AcceptRequest struct {
    PaxosMessage
    Value PaxosValue `json:"value"`
}

type AcceptReply struct {
    PaxosMessage
    Accepted bool  `json:"accepted"`
    Err      error `json:"err,omitempty"`
}

// Decision notification
type DecisionNotification struct {
    PaxosMessage
    Value PaxosValue `json:"value"`
}
```

## Proposer Protocol Implementation

### Proposer State

```go
type PaxosProposer struct {
    mu              sync.RWMutex
    me              NodeID
    acceptors       []NodeID
    currentProposal ProposalNumber
    instanceStates  map[InstanceID]*ProposalState
    done            chan struct{}
}

type ProposalState struct {
    instanceID      InstanceID
    proposalNumber  ProposalNumber
    value           PaxosValue
    phase           ProposalPhase
    prepareReplies  map[NodeID]*PrepareReply
    acceptReplies   map[NodeID]*AcceptReply
    decided         bool
    decidedValue    PaxosValue
}

type ProposalPhase int

const (
    PhasePrepare ProposalPhase = iota
    PhaseAccept
    PhaseDecided
)
```

### Proposer Implementation

```go
func (p *PaxosProposer) Propose(instanceID InstanceID, value PaxosValue) error {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    // Check if already decided
    if state, exists := p.instanceStates[instanceID]; exists && state.decided {
        return nil // Already decided
    }
    
    // Create or update proposal state
    proposalNumber := p.generateProposalNumber()
    state := &ProposalState{
        instanceID:     instanceID,
        proposalNumber: proposalNumber,
        value:          value,
        phase:          PhasePrepare,
        prepareReplies: make(map[NodeID]*PrepareReply),
        acceptReplies:  make(map[NodeID]*AcceptReply),
    }
    p.instanceStates[instanceID] = state
    
    // Start Phase 1: Prepare
    go p.startPreparePhase(state)
    
    return nil
}

func (p *PaxosProposer) startPreparePhase(state *ProposalState) {
    // Send prepare requests to all acceptors
    prepareReq := &PrepareRequest{
        PaxosMessage: PaxosMessage{
            Type:          "prepare",
            From:          p.me,
            InstanceID:    state.instanceID,
            ProposalNumber: state.proposalNumber,
        },
    }
    
    for _, acceptor := range p.acceptors {
        go p.sendPrepareRequest(acceptor, prepareReq, state)
    }
}

func (p *PaxosProposer) sendPrepareRequest(acceptor NodeID, req *PrepareRequest, state *ProposalState) {
    var reply PrepareReply
    err := p.call(acceptor, "PaxosAcceptor.Prepare", req, &reply)
    
    p.mu.Lock()
    defer p.mu.Unlock()
    
    if err != nil {
        log.Printf("Prepare request failed to %s: %v", acceptor, err)
        return
    }
    
    state.prepareReplies[acceptor] = &reply
    
    // Check if we have majority
    if len(state.prepareReplies) > len(p.acceptors)/2 {
        p.handlePrepareMajority(state)
    }
}

func (p *PaxosProposer) handlePrepareMajority(state *ProposalState) {
    if state.phase != PhasePrepare {
        return // Already moved to next phase
    }
    
    // Find highest numbered proposal among replies
    var highestValue PaxosValue
    var highestNumber ProposalNumber = -1
    
    for _, reply := range state.prepareReplies {
        if reply.Accepted && reply.Value != nil {
            // This acceptor has accepted a proposal
            if reply.ProposalNumber > highestNumber {
                highestNumber = reply.ProposalNumber
                highestValue = reply.Value
            }
        }
    }
    
    // Choose value: highest numbered proposal or our own
    chosenValue := state.value
    if highestValue != nil {
        chosenValue = highestValue
    }
    
    // Start Phase 2: Accept
    state.phase = PhaseAccept
    state.value = chosenValue
    
    go p.startAcceptPhase(state)
}

func (p *PaxosProposer) startAcceptPhase(state *ProposalState) {
    // Send accept requests to all acceptors
    acceptReq := &AcceptRequest{
        PaxosMessage: PaxosMessage{
            Type:          "accept",
            From:          p.me,
            InstanceID:    state.instanceID,
            ProposalNumber: state.proposalNumber,
        },
        Value: state.value,
    }
    
    for _, acceptor := range p.acceptors {
        go p.sendAcceptRequest(acceptor, acceptReq, state)
    }
}

func (p *PaxosProposer) sendAcceptRequest(acceptor NodeID, req *AcceptRequest, state *ProposalState) {
    var reply AcceptReply
    err := p.call(acceptor, "PaxosAcceptor.Accept", req, &reply)
    
    p.mu.Lock()
    defer p.mu.Unlock()
    
    if err != nil {
        log.Printf("Accept request failed to %s: %v", acceptor, err)
        return
    }
    
    state.acceptReplies[acceptor] = &reply
    
    // Check if we have majority
    if len(state.acceptReplies) > len(p.acceptors)/2 {
        p.handleAcceptMajority(state)
    }
}

func (p *PaxosProposer) handleAcceptMajority(state *ProposalState) {
    if state.phase != PhaseAccept {
        return // Already moved to next phase
    }
    
    // Count accepted replies
    acceptedCount := 0
    for _, reply := range state.acceptReplies {
        if reply.Accepted {
            acceptedCount++
        }
    }
    
    if acceptedCount > len(p.acceptors)/2 {
        // Value is chosen! Notify learners
        state.phase = PhaseDecided
        state.decided = true
        state.decidedValue = state.value
        
        go p.notifyLearners(state)
    }
}

func (p *PaxosProposer) notifyLearners(state *ProposalState) {
    decision := &DecisionNotification{
        PaxosMessage: PaxosMessage{
            Type:       "decision",
            From:       p.me,
            InstanceID: state.instanceID,
        },
        Value: state.decidedValue,
    }
    
    // Notify all nodes (including self)
    for _, node := range p.acceptors {
        go p.call(node, "PaxosLearner.Learn", decision, &struct{}{})
    }
}

func (p *PaxosProposer) generateProposalNumber() ProposalNumber {
    p.currentProposal++
    return p.currentProposal
}
```

## Acceptor Protocol Implementation

### Acceptor State

```go
type PaxosAcceptor struct {
    mu              sync.RWMutex
    me              NodeID
    highestPrepare  ProposalNumber
    highestAccept   ProposalNumber
    acceptedValue   PaxosValue
    instanceStates  map[InstanceID]*AcceptorState
}

type AcceptorState struct {
    instanceID     InstanceID
    highestPrepare ProposalNumber
    highestAccept  ProposalNumber
    acceptedValue  PaxosValue
}
```

### Acceptor Implementation

```go
func (a *PaxosAcceptor) Prepare(req *PrepareRequest, reply *PrepareReply) error {
    a.mu.Lock()
    defer a.mu.Unlock()
    
    instanceID := req.InstanceID
    proposalNumber := req.ProposalNumber
    
    // Get or create state for this instance
    state, exists := a.instanceStates[instanceID]
    if !exists {
        state = &AcceptorState{
            instanceID:     instanceID,
            highestPrepare: -1,
            highestAccept:  -1,
        }
        a.instanceStates[instanceID] = state
    }
    
    // P1a: An acceptor can accept a proposal numbered n iff it has not responded
    // to a prepare request having number greater than n
    if proposalNumber > state.highestPrepare {
        state.highestPrepare = proposalNumber
        
        reply.Accepted = true
        reply.ProposalNumber = state.highestAccept
        reply.Value = state.acceptedValue
        
        log.Printf("Acceptor %s: Prepared proposal %d for instance %d", 
            a.me, proposalNumber, instanceID)
    } else {
        reply.Accepted = false
        log.Printf("Acceptor %s: Rejected prepare for proposal %d (highest: %d)", 
            a.me, proposalNumber, state.highestPrepare)
    }
    
    return nil
}

func (a *PaxosAcceptor) Accept(req *AcceptRequest, reply *AcceptReply) error {
    a.mu.Lock()
    defer a.mu.Unlock()
    
    instanceID := req.InstanceID
    proposalNumber := req.ProposalNumber
    value := req.Value
    
    // Get state for this instance
    state, exists := a.instanceStates[instanceID]
    if !exists {
        reply.Accepted = false
        reply.Err = errors.New("no state for instance")
        return nil
    }
    
    // P1a: An acceptor can accept a proposal numbered n iff it has not responded
    // to a prepare request having number greater than n
    if proposalNumber >= state.highestPrepare {
        state.highestAccept = proposalNumber
        state.acceptedValue = value
        
        reply.Accepted = true
        
        log.Printf("Acceptor %s: Accepted proposal %d with value %v for instance %d", 
            a.me, proposalNumber, value, instanceID)
    } else {
        reply.Accepted = false
        log.Printf("Acceptor %s: Rejected accept for proposal %d (highest prepare: %d)", 
            a.me, proposalNumber, state.highestPrepare)
    }
    
    return nil
}
```

## Unique Proposal Number Generation

### Problem: Reusing Proposal Numbers

Without proper coordination, multiple proposers might use the same proposal number, leading to conflicts.

### Solution: Unique Proposal Numbers

```go
type ProposalNumberGenerator struct {
    mu        sync.Mutex
    base      ProposalNumber
    increment ProposalNumber
    nodeID    NodeID
}

func NewProposalNumberGenerator(nodeID NodeID, totalNodes int) *ProposalNumberGenerator {
    return &ProposalNumberGenerator{
        base:      ProposalNumber(nodeID),
        increment: ProposalNumber(totalNodes),
        nodeID:    nodeID,
    }
}

func (png *ProposalNumberGenerator) Next() ProposalNumber {
    png.mu.Lock()
    defer png.mu.Unlock()
    
    png.base += png.increment
    return png.base
}

// Alternative: Use timestamp + nodeID
func (png *ProposalNumberGenerator) NextWithTimestamp() ProposalNumber {
    png.mu.Lock()
    defer png.mu.Unlock()
    
    timestamp := time.Now().UnixNano()
    return ProposalNumber(timestamp*1000 + int64(png.nodeID))
}
```

## Leader Election

### Simple Retry Strategy

```go
type PaxosLeader struct {
    mu              sync.RWMutex
    me              NodeID
    proposer        *PaxosProposer
    isLeader        bool
    leaderTimeout   time.Duration
    lastLeaderTime  time.Time
    retryInterval   time.Duration
    done            chan struct{}
}

func (l *PaxosLeader) Start() error {
    go l.leaderElectionLoop()
    return nil
}

func (l *PaxosLeader) leaderElectionLoop() {
    ticker := time.NewTicker(l.retryInterval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            l.checkLeadership()
        case <-l.done:
            return
        }
    }
}

func (l *PaxosLeader) checkLeadership() {
    l.mu.Lock()
    defer l.mu.Unlock()
    
    if l.isLeader {
        // Check if we're still the leader
        if time.Since(l.lastLeaderTime) > l.leaderTimeout {
            l.isLeader = false
            log.Printf("Node %s: Lost leadership due to timeout", l.me)
        }
    } else {
        // Try to become leader
        l.attemptLeadership()
    }
}

func (l *PaxosLeader) attemptLeadership() {
    // Try to propose a leadership value
    leadershipValue := &LeadershipValue{
        LeaderID: l.me,
        Timestamp: time.Now().UnixNano(),
    }
    
    err := l.proposer.Propose(0, leadershipValue) // Use instance 0 for leadership
    if err == nil {
        l.isLeader = true
        l.lastLeaderTime = time.Now()
        log.Printf("Node %s: Became leader", l.me)
    }
}
```

## Complete Paxos Implementation

### Paxos Node

```go
type PaxosNode struct {
    mu        sync.RWMutex
    me        NodeID
    peers     []NodeID
    proposer  *PaxosProposer
    acceptor  *PaxosAcceptor
    learner   *PaxosLearner
    leader    *PaxosLeader
    instances map[InstanceID]*PaxosInstance
    done      chan struct{}
}

type PaxosInstance struct {
    instanceID InstanceID
    decided    bool
    value      PaxosValue
    mu         sync.RWMutex
}

func (pn *PaxosNode) Start() error {
    // Initialize components
    pn.proposer = NewPaxosProposer(pn.me, pn.peers)
    pn.acceptor = NewPaxosAcceptor(pn.me)
    pn.learner = NewPaxosLearner(pn.me)
    pn.leader = NewPaxosLeader(pn.me, pn.proposer)
    
    // Start components
    go pn.proposer.Start()
    go pn.acceptor.Start()
    go pn.learner.Start()
    go pn.leader.Start()
    
    return nil
}

func (pn *PaxosNode) Propose(instanceID InstanceID, value PaxosValue) error {
    pn.mu.RLock()
    defer pn.mu.RUnlock()
    
    // Check if already decided
    if instance, exists := pn.instances[instanceID]; exists && instance.decided {
        return nil
    }
    
    // Only leader can propose
    if !pn.leader.IsLeader() {
        return errors.New("not the leader")
    }
    
    return pn.proposer.Propose(instanceID, value)
}

func (pn *PaxosNode) GetDecision(instanceID InstanceID) (PaxosValue, bool) {
    pn.mu.RLock()
    defer pn.mu.RUnlock()
    
    if instance, exists := pn.instances[instanceID]; exists && instance.decided {
        return instance.value, true
    }
    
    return nil, false
}
```

## Testing Strategy

### Unit Testing

```go
func TestPaxosProposer(t *testing.T) {
    proposer := NewPaxosProposer("proposer1", []NodeID{"acceptor1", "acceptor2", "acceptor3"})
    
    // Test proposal
    value := &TestValue{Data: "test"}
    err := proposer.Propose(1, value)
    assert.NoError(t, err)
    
    // Wait for decision
    time.Sleep(100 * time.Millisecond)
    
    // Verify decision
    decided, exists := proposer.GetDecision(1)
    assert.True(t, exists)
    assert.Equal(t, value, decided)
}

func TestPaxosAcceptor(t *testing.T) {
    acceptor := NewPaxosAcceptor("acceptor1")
    
    // Test prepare
    prepareReq := &PrepareRequest{
        PaxosMessage: PaxosMessage{
            InstanceID:    1,
            ProposalNumber: 1,
        },
    }
    
    var reply PrepareReply
    err := acceptor.Prepare(prepareReq, &reply)
    assert.NoError(t, err)
    assert.True(t, reply.Accepted)
}
```

### Integration Testing

```go
func TestPaxosConsensus(t *testing.T) {
    // Create 3 Paxos nodes
    nodes := make([]*PaxosNode, 3)
    for i := 0; i < 3; i++ {
        nodeID := NodeID(fmt.Sprintf("node%d", i))
        peers := make([]NodeID, 0, 2)
        for j := 0; j < 3; j++ {
            if i != j {
                peers = append(peers, NodeID(fmt.Sprintf("node%d", j)))
            }
        }
        
        nodes[i] = NewPaxosNode(nodeID, peers)
        go nodes[i].Start()
    }
    
    // Wait for leader election
    time.Sleep(200 * time.Millisecond)
    
    // Find leader
    var leader *PaxosNode
    for _, node := range nodes {
        if node.leader.IsLeader() {
            leader = node
            break
        }
    }
    assert.NotNil(t, leader, "No leader elected")
    
    // Propose value
    value := &TestValue{Data: "consensus_test"}
    err := leader.Propose(1, value)
    assert.NoError(t, err)
    
    // Wait for consensus
    time.Sleep(500 * time.Millisecond)
    
    // Verify all nodes learned the same value
    for _, node := range nodes {
        decided, exists := node.GetDecision(1)
        assert.True(t, exists, "Node %s did not learn decision", node.me)
        assert.Equal(t, value, decided, "Node %s learned different value", node.me)
    }
}
```

## Best Practices and Common Pitfalls

### Do's

1. **Always check proposal numbers**: Ensure proper ordering
2. **Handle majority correctly**: Count only accepted replies
3. **Persist acceptor state**: Survive crashes and restarts
4. **Use unique proposal numbers**: Avoid conflicts between proposers
5. **Implement proper timeouts**: Handle network delays and failures
6. **Log important events**: Include proposal numbers and instance IDs
7. **Test with failures**: Simulate node crashes and network partitions
8. **Handle concurrent proposals**: Multiple proposers can run simultaneously

### Don'ts

1. **Don't ignore proposal number ordering**: Always check against highest seen
2. **Don't forget to persist state**: Acceptor state must survive crashes
3. **Don't assume synchronous execution**: Handle asynchronous message delivery
4. **Don't ignore majority requirements**: Need majority for both phases
5. **Don't forget to handle timeouts**: Proposals can fail and need retry
6. **Don't assume single proposer**: Multiple proposers can compete
7. **Don't ignore message ordering**: Messages may arrive out of order
8. **Don't forget to handle duplicate messages**: Implement idempotency

### Common Implementation Mistakes

```go
// WRONG: Not checking proposal number ordering
func (a *PaxosAcceptor) Prepare(req *PrepareRequest, reply *PrepareReply) error {
    // Missing proposal number check!
    reply.Accepted = true
    return nil
}

// CORRECT: Always check proposal number ordering
func (a *PaxosAcceptor) Prepare(req *PrepareRequest, reply *PrepareReply) error {
    if req.ProposalNumber > a.highestPrepare {
        a.highestPrepare = req.ProposalNumber
        reply.Accepted = true
    } else {
        reply.Accepted = false
    }
    return nil
}

// WRONG: Not handling majority correctly
func (p *PaxosProposer) handlePrepareMajority(state *ProposalState) {
    if len(state.prepareReplies) > 0 { // Wrong! Need majority
        p.startAcceptPhase(state)
    }
}

// CORRECT: Check for majority
func (p *PaxosProposer) handlePrepareMajority(state *ProposalState) {
    if len(state.prepareReplies) > len(p.acceptors)/2 { // Correct!
        p.startAcceptPhase(state)
    }
}
```

### Performance Considerations

1. **Minimize message overhead**: Batch operations when possible
2. **Use efficient serialization**: Choose appropriate encoding
3. **Implement connection pooling**: Reuse network connections
4. **Optimize proposal number generation**: Avoid conflicts
5. **Monitor consensus latency**: Track time to reach decisions

### Debugging Tips

1. **Add structured logging**: Include proposal numbers and instance IDs
2. **Use consistent log levels**: Debug, Info, Warn, Error
3. **Include correlation IDs**: Track proposals across nodes
4. **Monitor consensus progress**: Log phase transitions
5. **Test with delays**: Add artificial delays to find race conditions

### Running Tests

```bash
# Run Paxos tests
go test -v -run TestPaxos ./...

# Run with race detection
go test -race -run TestPaxos ./...

# Run specific test
go test -v -run TestPaxosConsensus ./...

# Run with coverage
go test -cover -run TestPaxos ./...

# Run benchmarks
go test -bench=. -run TestPaxos ./...
```

### Next Steps

1. Implement basic Paxos proposer and acceptor
2. Add unique proposal number generation
3. Implement leader election
4. Add state persistence
5. Implement comprehensive testing
6. Add monitoring and logging
7. Performance optimization

This foundation will prepare you for implementing a complete, fault-tolerant consensus system using the Paxos algorithm.
