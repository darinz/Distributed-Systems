# Lab 4C: Distributed Transactions with Two-Phase Commit

## Overview

This lab implements distributed transactions across multiple keys that may be located in different replica groups using Two-Phase Commit (2PC). The system supports MultiGet, MultiPut, and Swap operations that can span multiple shards while maintaining ACID properties.

## Why Two-Phase Commit?

### Problem with Paxos Alone

- **Paxos Limitations**: Each replica group runs Paxos independently
- **Cross-Shard Operations**: Transactions may involve keys in different replica groups
- **Consistency Requirements**: Need to ensure all-or-nothing execution across groups
- **Lock Coordination**: Must prevent concurrent modifications to the same keys

### Two-Phase Commit Solution

- **Coordinator Selection**: One replica group acts as transaction coordinator
- **Prepare Phase**: Acquire locks on all involved keys
- **Commit Phase**: Execute transaction if all locks acquired successfully
- **Abort Handling**: Release locks and abort if any participant fails

## Transaction Architecture

### Core Types and Interfaces

```go
// Transaction interface
type Transaction interface {
    GetType() string
    GetReadSet() []string
    GetWriteSet() []string
    GetKeySet() []string
    GetConfigNum() int64
    GetAttemptNum() int64
}

// MultiGet transaction
type MultiGet struct {
    Keys       []string `json:"keys"`
    ConfigNum  int64    `json:"config_num"`
    AttemptNum int64    `json:"attempt_num"`
    ClientID   NodeID   `json:"client_id"`
    SeqNum     int64    `json:"seq_num"`
}

func (mg *MultiGet) GetType() string { return "MultiGet" }
func (mg *MultiGet) GetReadSet() []string { return mg.Keys }
func (mg *MultiGet) GetWriteSet() []string { return []string{} }
func (mg *MultiGet) GetKeySet() []string { return mg.Keys }
func (mg *MultiGet) GetConfigNum() int64 { return mg.ConfigNum }
func (mg *MultiGet) GetAttemptNum() int64 { return mg.AttemptNum }

// MultiPut transaction
type MultiPut struct {
    Values     map[string]string `json:"values"`
    ConfigNum  int64             `json:"config_num"`
    AttemptNum int64             `json:"attempt_num"`
    ClientID   NodeID            `json:"client_id"`
    SeqNum     int64             `json:"seq_num"`
}

func (mp *MultiPut) GetType() string { return "MultiPut" }
func (mp *MultiPut) GetReadSet() []string { return []string{} }
func (mp *MultiPut) GetWriteSet() []string {
    keys := make([]string, 0, len(mp.Values))
    for key := range mp.Values {
        keys = append(keys, key)
    }
    return keys
}
func (mp *MultiPut) GetKeySet() []string { return mp.GetWriteSet() }
func (mp *MultiPut) GetConfigNum() int64 { return mp.ConfigNum }
func (mp *MultiPut) GetAttemptNum() int64 { return mp.AttemptNum }

// Swap transaction
type Swap struct {
    Key1       string `json:"key1"`
    Key2       string `json:"key2"`
    ConfigNum  int64  `json:"config_num"`
    AttemptNum int64  `json:"attempt_num"`
    ClientID   NodeID `json:"client_id"`
    SeqNum     int64  `json:"seq_num"`
}

func (s *Swap) GetType() string { return "Swap" }
func (s *Swap) GetReadSet() []string { return []string{s.Key1, s.Key2} }
func (s *Swap) GetWriteSet() []string { return []string{s.Key1, s.Key2} }
func (s *Swap) GetKeySet() []string { return []string{s.Key1, s.Key2} }
func (s *Swap) GetConfigNum() int64 { return s.ConfigNum }
func (s *Swap) GetAttemptNum() int64 { return s.AttemptNum }

// 2PC message types
type Prepare struct {
    Transaction Transaction `json:"transaction"`
    ConfigNum   int64       `json:"config_num"`
    AttemptNum  int64       `json:"attempt_num"`
}

type PrepareOK struct {
    Transaction Transaction     `json:"transaction"`
    Values      map[string]string `json:"values"`
    ConfigNum   int64           `json:"config_num"`
    AttemptNum  int64           `json:"attempt_num"`
}

type PrepareAbort struct {
    Transaction Transaction `json:"transaction"`
    ConfigNum   int64       `json:"config_num"`
    AttemptNum  int64       `json:"attempt_num"`
    Reason      string      `json:"reason"`
}

type Commit struct {
    Transaction Transaction `json:"transaction"`
    ConfigNum   int64       `json:"config_num"`
    AttemptNum  int64       `json:"attempt_num"`
}

type CommitOK struct {
    Transaction Transaction `json:"transaction"`
    ConfigNum   int64       `json:"config_num"`
    AttemptNum  int64       `json:"attempt_num"`
}

type Abort struct {
    Transaction Transaction `json:"transaction"`
    ConfigNum   int64       `json:"config_num"`
    AttemptNum  int64       `json:"attempt_num"`
    Reason      string      `json:"reason"`
}
```

### Transaction State Management

```go
type TransactionState struct {
    Transaction Transaction
    ConfigNum   int64
    AttemptNum  int64
    Status      string // "PREPARING", "COMMITTING", "COMMITTED", "ABORTED"
    Locks       map[string]bool // key -> locked
    Participants []int64 // group IDs involved
    PrepareOKs  map[int64]bool // group ID -> received prepare OK
    CommitOKs   map[int64]bool // group ID -> received commit OK
    Values      map[string]string // key -> value (for reads)
}

type TransactionManager struct {
    mu           sync.RWMutex
    transactions map[string]*TransactionState // transaction ID -> state
    locks        map[string]string // key -> transaction ID
    me           NodeID
    config       Configuration
}

func NewTransactionManager(me NodeID) *TransactionManager {
    return &TransactionManager{
        transactions: make(map[string]*TransactionState),
        locks:        make(map[string]string),
        me:           me,
    }
}

func (tm *TransactionManager) generateTransactionID(txn Transaction) string {
    return fmt.Sprintf("%s_%d_%d_%s", txn.GetType(), txn.GetConfigNum(), txn.GetAttemptNum(), txn.GetClientID())
}
```

## Coordinator Logic

### Coordinator Implementation

```go
type TransactionCoordinator struct {
    mu           sync.RWMutex
    me           NodeID
    config       Configuration
    txnManager   *TransactionManager
    shardMasters []NodeID
    done         chan struct{}
}

func NewTransactionCoordinator(me NodeID, shardMasters []NodeID) *TransactionCoordinator {
    return &TransactionCoordinator{
        me:           me,
        txnManager:   NewTransactionManager(me),
        shardMasters: shardMasters,
        done:         make(chan struct{}),
    }
}

func (tc *TransactionCoordinator) Start() error {
    go tc.configPuller()
    return nil
}

func (tc *TransactionCoordinator) Stop() {
    close(tc.done)
}
```

### Transaction Execution

```go
func (tc *TransactionCoordinator) ExecuteTransaction(txn Transaction) error {
    tc.mu.Lock()
    defer tc.mu.Unlock()
    
    // Determine coordinator (group with highest ID)
    coordinatorGroup := tc.determineCoordinator(txn)
    if coordinatorGroup != tc.getMyGroupID() {
        // Forward to coordinator
        return tc.forwardToCoordinator(txn, coordinatorGroup)
    }
    
    // Create transaction state
    txnID := tc.txnManager.generateTransactionID(txn)
    state := &TransactionState{
        Transaction:  txn,
        ConfigNum:    txn.GetConfigNum(),
        AttemptNum:   txn.GetAttemptNum(),
        Status:       "PREPARING",
        Locks:        make(map[string]bool),
        Participants: tc.getParticipantGroups(txn),
        PrepareOKs:   make(map[int64]bool),
        CommitOKs:    make(map[int64]bool),
        Values:       make(map[string]string),
    }
    
    tc.txnManager.transactions[txnID] = state
    
    // Start 2PC
    go tc.runTwoPhaseCommit(txnID)
    
    return nil
}

func (tc *TransactionCoordinator) determineCoordinator(txn Transaction) int64 {
    groups := tc.getParticipantGroups(txn)
    if len(groups) == 0 {
        return 0
    }
    
    maxGroup := groups[0]
    for _, group := range groups {
        if group > maxGroup {
            maxGroup = group
        }
    }
    return maxGroup
}

func (tc *TransactionCoordinator) getParticipantGroups(txn Transaction) []int64 {
    groups := make(map[int64]bool)
    
    for _, key := range txn.GetKeySet() {
        shard := keyToShard(key)
        group := tc.config.Shards[shard]
        if group != 0 {
            groups[group] = true
        }
    }
    
    result := make([]int64, 0, len(groups))
    for group := range groups {
        result = append(result, group)
    }
    
    return result
}
```

### Two-Phase Commit Implementation

```go
func (tc *TransactionCoordinator) runTwoPhaseCommit(txnID string) {
    tc.mu.RLock()
    state := tc.txnManager.transactions[txnID]
    tc.mu.RUnlock()
    
    if state == nil {
        return
    }
    
    // Phase 1: Prepare
    if !tc.preparePhase(txnID) {
        tc.abortTransaction(txnID, "Prepare phase failed")
        return
    }
    
    // Phase 2: Commit
    tc.commitPhase(txnID)
}

func (tc *TransactionCoordinator) preparePhase(txnID string) bool {
    tc.mu.RLock()
    state := tc.txnManager.transactions[txnID]
    tc.mu.RUnlock()
    
    if state == nil {
        return false
    }
    
    // Send prepare messages to all participants
    prepare := &Prepare{
        Transaction: state.Transaction,
        ConfigNum:   state.ConfigNum,
        AttemptNum:  state.AttemptNum,
    }
    
    for _, groupID := range state.Participants {
        go tc.sendPrepare(prepare, groupID)
    }
    
    // Wait for prepare responses
    timeout := time.After(5 * time.Second)
    for {
        select {
        case <-timeout:
            return false
        default:
            tc.mu.RLock()
            state := tc.txnManager.transactions[txnID]
            tc.mu.RUnlock()
            
            if state == nil {
                return false
            }
            
            // Check if all participants responded
            if len(state.PrepareOKs) == len(state.Participants) {
                return true
            }
            
            // Check if any participant aborted
            if state.Status == "ABORTED" {
                return false
            }
            
            time.Sleep(10 * time.Millisecond)
        }
    }
}

func (tc *TransactionCoordinator) commitPhase(txnID string) {
    tc.mu.RLock()
    state := tc.txnManager.transactions[txnID]
    tc.mu.RUnlock()
    
    if state == nil {
        return
    }
    
    // Update status
    tc.mu.Lock()
    state.Status = "COMMITTING"
    tc.mu.Unlock()
    
    // Send commit messages to all participants
    commit := &Commit{
        Transaction: state.Transaction,
        ConfigNum:   state.ConfigNum,
        AttemptNum:  state.AttemptNum,
    }
    
    for _, groupID := range state.Participants {
        go tc.sendCommit(commit, groupID)
    }
    
    // Wait for commit responses
    timeout := time.After(5 * time.Second)
    for {
        select {
        case <-timeout:
            log.Printf("Commit phase timeout for transaction %s", txnID)
            return
        default:
            tc.mu.RLock()
            state := tc.txnManager.transactions[txnID]
            tc.mu.RUnlock()
            
            if state == nil {
                return
            }
            
            // Check if all participants committed
            if len(state.CommitOKs) == len(state.Participants) {
                tc.mu.Lock()
                state.Status = "COMMITTED"
                tc.mu.Unlock()
                
                // Send response to client
                tc.sendTransactionResponse(txnID)
                return
            }
            
            time.Sleep(10 * time.Millisecond)
        }
    }
}
```

## Participant Logic

### Participant Implementation

```go
type TransactionParticipant struct {
    mu           sync.RWMutex
    me           NodeID
    config       Configuration
    txnManager   *TransactionManager
    shardMasters []NodeID
    done         chan struct{}
}

func NewTransactionParticipant(me NodeID, shardMasters []NodeID) *TransactionParticipant {
    return &TransactionParticipant{
        me:           me,
        txnManager:   NewTransactionManager(me),
        shardMasters: shardMasters,
        done:         make(chan struct{}),
    }
}

func (tp *TransactionParticipant) Start() error {
    go tp.configPuller()
    return nil
}

func (tp *TransactionParticipant) Stop() {
    close(tp.done)
}
```

### Prepare Phase Handling

```go
func (tp *TransactionParticipant) HandlePrepare(prepare *Prepare, reply *PrepareOK) error {
    tp.mu.Lock()
    defer tp.mu.Unlock()
    
    txn := prepare.Transaction
    txnID := tp.txnManager.generateTransactionID(txn)
    
    // Validate configuration
    if prepare.ConfigNum != tp.config.ConfigNum {
        reply.Err = ErrConfigMismatch
        return nil
    }
    
    // Check if keys are locked
    for _, key := range txn.GetKeySet() {
        if lockedTxnID, exists := tp.txnManager.locks[key]; exists && lockedTxnID != txnID {
            // Key is locked by another transaction
            reply.Err = ErrKeyLocked
            return nil
        }
    }
    
    // Lock keys
    for _, key := range txn.GetKeySet() {
        tp.txnManager.locks[key] = txnID
    }
    
    // Create transaction state
    state := &TransactionState{
        Transaction: txn,
        ConfigNum:   prepare.ConfigNum,
        AttemptNum:  prepare.AttemptNum,
        Status:      "PREPARING",
        Locks:       make(map[string]bool),
        Values:      make(map[string]string),
    }
    
    for _, key := range txn.GetKeySet() {
        state.Locks[key] = true
    }
    
    tp.txnManager.transactions[txnID] = state
    
    // Read values for read operations
    if txn.GetType() == "MultiGet" {
        for _, key := range txn.GetReadSet() {
            if value, exists := tp.getKeyValue(key); exists {
                state.Values[key] = value
            }
        }
    }
    
    // Send prepare OK
    reply.Transaction = txn
    reply.Values = state.Values
    reply.ConfigNum = prepare.ConfigNum
    reply.AttemptNum = prepare.AttemptNum
    reply.Err = OK
    
    return nil
}
```

### Commit Phase Handling

```go
func (tp *TransactionParticipant) HandleCommit(commit *Commit, reply *CommitOK) error {
    tp.mu.Lock()
    defer tp.mu.Unlock()
    
    txn := commit.Transaction
    txnID := tp.txnManager.generateTransactionID(txn)
    
    // Validate transaction
    if commit.ConfigNum != tp.config.ConfigNum {
        reply.Err = ErrConfigMismatch
        return nil
    }
    
    state := tp.txnManager.transactions[txnID]
    if state == nil {
        reply.Err = ErrTransactionNotFound
        return nil
    }
    
    if state.AttemptNum != commit.AttemptNum {
        reply.Err = ErrAttemptMismatch
        return nil
    }
    
    // Execute transaction
    switch txn.GetType() {
    case "MultiGet":
        // Values already read in prepare phase
        break
    case "MultiPut":
        mp := txn.(*MultiPut)
        for key, value := range mp.Values {
            tp.setKeyValue(key, value)
        }
    case "Swap":
        s := txn.(*Swap)
        val1, _ := tp.getKeyValue(s.Key1)
        val2, _ := tp.getKeyValue(s.Key2)
        tp.setKeyValue(s.Key1, val2)
        tp.setKeyValue(s.Key2, val1)
    }
    
    // Update client sequence number
    tp.updateClientSequence(txn.GetClientID(), txn.GetSeqNum())
    
    // Unlock keys
    for key := range state.Locks {
        delete(tp.txnManager.locks, key)
    }
    
    // Update state
    state.Status = "COMMITTED"
    
    // Send commit OK
    reply.Transaction = txn
    reply.ConfigNum = commit.ConfigNum
    reply.AttemptNum = commit.AttemptNum
    reply.Err = OK
    
    return nil
}
```

### Abort Handling

```go
func (tp *TransactionParticipant) HandleAbort(abort *Abort, reply *OKReply) error {
    tp.mu.Lock()
    defer tp.mu.Unlock()
    
    txn := abort.Transaction
    txnID := tp.txnManager.generateTransactionID(txn)
    
    state := tp.txnManager.transactions[txnID]
    if state == nil {
        reply.Err = OK // Already aborted or not found
        return nil
    }
    
    // Unlock keys
    for key := range state.Locks {
        delete(tp.txnManager.locks, key)
    }
    
    // Update state
    state.Status = "ABORTED"
    
    reply.Err = OK
    return nil
}

func (tp *TransactionParticipant) abortTransaction(txnID string, reason string) {
    tp.mu.Lock()
    defer tp.mu.Unlock()
    
    state := tp.txnManager.transactions[txnID]
    if state == nil {
        return
    }
    
    // Unlock keys
    for key := range state.Locks {
        delete(tp.txnManager.locks, key)
    }
    
    // Update state
    state.Status = "ABORTED"
    
    log.Printf("Transaction %s aborted: %s", txnID, reason)
}
```

## Lock Management

### Lock Implementation

```go
type LockManager struct {
    mu    sync.RWMutex
    locks map[string]string // key -> transaction ID
}

func NewLockManager() *LockManager {
    return &LockManager{
        locks: make(map[string]string),
    }
}

func (lm *LockManager) TryLock(key, txnID string) bool {
    lm.mu.Lock()
    defer lm.mu.Unlock()
    
    if existingTxnID, exists := lm.locks[key]; exists && existingTxnID != txnID {
        return false // Key is locked by another transaction
    }
    
    lm.locks[key] = txnID
    return true
}

func (lm *LockManager) Unlock(key string) {
    lm.mu.Lock()
    defer lm.mu.Unlock()
    
    delete(lm.locks, key)
}

func (lm *LockManager) UnlockTransaction(txnID string) {
    lm.mu.Lock()
    defer lm.mu.Unlock()
    
    for key, lockedTxnID := range lm.locks {
        if lockedTxnID == txnID {
            delete(lm.locks, key)
        }
    }
}

func (lm *LockManager) IsLocked(key string) bool {
    lm.mu.RLock()
    defer lm.mu.RUnlock()
    
    _, exists := lm.locks[key]
    return exists
}

func (lm *LockManager) GetLockedBy(key string) string {
    lm.mu.RLock()
    defer lm.mu.RUnlock()
    
    return lm.locks[key]
}
```

## Transaction Client

### Client Implementation

```go
type TransactionClient struct {
    mu           sync.RWMutex
    me           NodeID
    shardMasters []NodeID
    config       Configuration
    seq          int64
}

func NewTransactionClient(me NodeID, shardMasters []NodeID) *TransactionClient {
    return &TransactionClient{
        me:           me,
        shardMasters: shardMasters,
        seq:          0,
    }
}

func (tc *TransactionClient) MultiGet(keys []string) (map[string]string, error) {
    txn := &MultiGet{
        Keys:       keys,
        ConfigNum:  tc.config.ConfigNum,
        AttemptNum: 1,
        ClientID:   tc.me,
        SeqNum:     tc.getNextSeq(),
    }
    
    return tc.executeTransaction(txn)
}

func (tc *TransactionClient) MultiPut(values map[string]string) error {
    txn := &MultiPut{
        Values:     values,
        ConfigNum:  tc.config.ConfigNum,
        AttemptNum: 1,
        ClientID:   tc.me,
        SeqNum:     tc.getNextSeq(),
    }
    
    _, err := tc.executeTransaction(txn)
    return err
}

func (tc *TransactionClient) Swap(key1, key2 string) error {
    txn := &Swap{
        Key1:       key1,
        Key2:       key2,
        ConfigNum:  tc.config.ConfigNum,
        AttemptNum: 1,
        ClientID:   tc.me,
        SeqNum:     tc.getNextSeq(),
    }
    
    _, err := tc.executeTransaction(txn)
    return err
}

func (tc *TransactionClient) executeTransaction(txn Transaction) (map[string]string, error) {
    // Get current configuration
    config, err := tc.getCurrentConfig()
    if err != nil {
        return nil, err
    }
    
    // Determine coordinator
    coordinatorGroup := tc.determineCoordinator(txn, config)
    if coordinatorGroup == 0 {
        return nil, errors.New("no coordinator found")
    }
    
    // Get coordinator servers
    coordinatorServers := config.Groups[coordinatorGroup]
    if len(coordinatorServers) == 0 {
        return nil, errors.New("no coordinator servers")
    }
    
    // Send transaction to coordinator
    for _, server := range coordinatorServers {
        var reply TransactionReply
        err := tc.call(server, "TransactionCoordinator.ExecuteTransaction", txn, &reply)
        if err == nil {
            if reply.Err == OK {
                return reply.Values, nil
            } else if reply.Err == ErrConfigMismatch {
                // Retry with new configuration
                return tc.executeTransaction(txn)
            }
        }
    }
    
    return nil, errors.New("all coordinator servers failed")
}

func (tc *TransactionClient) determineCoordinator(txn Transaction, config Configuration) int64 {
    groups := make(map[int64]bool)
    
    for _, key := range txn.GetKeySet() {
        shard := keyToShard(key)
        group := config.Shards[shard]
        if group != 0 {
            groups[group] = true
        }
    }
    
    maxGroup := int64(0)
    for group := range groups {
        if group > maxGroup {
            maxGroup = group
        }
    }
    
    return maxGroup
}
```

## Testing Strategy

### Unit Testing

```go
func TestTransactionCoordinator(t *testing.T) {
    coordinator := NewTransactionCoordinator("coord1", []NodeID{"sm1"})
    go coordinator.Start()
    defer coordinator.Stop()
    
    // Test MultiGet transaction
    txn := &MultiGet{
        Keys:       []string{"key1", "key2"},
        ConfigNum:  1,
        AttemptNum: 1,
        ClientID:   "client1",
        SeqNum:     1,
    }
    
    err := coordinator.ExecuteTransaction(txn)
    assert.NoError(t, err)
    
    // Verify transaction state
    txnID := coordinator.txnManager.generateTransactionID(txn)
    state := coordinator.txnManager.transactions[txnID]
    assert.NotNil(t, state)
    assert.Equal(t, "PREPARING", state.Status)
}

func TestTransactionParticipant(t *testing.T) {
    participant := NewTransactionParticipant("part1", []NodeID{"sm1"})
    go participant.Start()
    defer participant.Stop()
    
    // Test prepare handling
    txn := &MultiGet{
        Keys:       []string{"key1"},
        ConfigNum:  1,
        AttemptNum: 1,
        ClientID:   "client1",
        SeqNum:     1,
    }
    
    prepare := &Prepare{
        Transaction: txn,
        ConfigNum:   1,
        AttemptNum:  1,
    }
    
    var reply PrepareOK
    err := participant.HandlePrepare(prepare, &reply)
    assert.NoError(t, err)
    assert.Equal(t, OK, reply.Err)
    
    // Verify keys are locked
    assert.True(t, participant.txnManager.locks["key1"] != "")
}
```

### Integration Testing

```go
func TestDistributedTransaction(t *testing.T) {
    // Create ShardMaster
    sm := NewShardMaster([]NodeID{"sm1", "sm2", "sm3"}, "sm1")
    go sm.Start()
    defer sm.Stop()
    
    // Create replica groups
    join1 := &JoinCommand{
        GroupID: 1,
        Servers: []NodeID{"server1", "server2", "server3"},
    }
    join2 := &JoinCommand{
        GroupID: 2,
        Servers: []NodeID{"server4", "server5", "server6"},
    }
    
    var reply OKReply
    sm.Join(join1, &reply)
    sm.Join(join2, &reply)
    
    // Wait for configuration to propagate
    time.Sleep(200 * time.Millisecond)
    
    // Create transaction client
    client := NewTransactionClient("client1", []NodeID{"sm1", "sm2", "sm3"})
    
    // Test MultiGet across groups
    values, err := client.MultiGet([]string{"key1", "key2"})
    assert.NoError(t, err)
    assert.NotNil(t, values)
    
    // Test MultiPut across groups
    values = map[string]string{
        "key1": "value1",
        "key2": "value2",
    }
    err = client.MultiPut(values)
    assert.NoError(t, err)
    
    // Test Swap across groups
    err = client.Swap("key1", "key2")
    assert.NoError(t, err)
}
```

## Best Practices and Common Pitfalls

### Do's

1. **Use proper coordinator selection**: Choose group with highest ID
2. **Validate configuration numbers**: Ensure all participants have same config
3. **Implement proper locking**: Lock all keys before prepare phase
4. **Handle timeouts**: Set reasonable timeouts for 2PC phases
5. **Implement retry logic**: Retry transactions on configuration changes
6. **Use attempt numbers**: Prevent duplicate transaction processing
7. **Handle aborts properly**: Release locks and clean up state
8. **Test with multiple groups**: Verify cross-group transactions work

### Don'ts

1. **Don't bypass 2PC**: All transactions must go through prepare/commit phases
2. **Don't ignore configuration mismatches**: Always validate config numbers
3. **Don't forget to unlock keys**: Always release locks on abort/commit
4. **Don't process duplicate transactions**: Use attempt numbers for deduplication
5. **Don't ignore timeouts**: Handle 2PC phase timeouts gracefully
6. **Don't assume single coordinator**: Handle coordinator failures
7. **Don't forget error handling**: Handle all failure scenarios
8. **Don't skip validation**: Always validate transaction parameters

### Common Implementation Mistakes

```go
// WRONG: Not validating configuration
func (tp *TransactionParticipant) HandlePrepare(prepare *Prepare, reply *PrepareOK) error {
    // Missing configuration validation!
    txn := prepare.Transaction
    // ... rest of implementation
}

// CORRECT: Always validate configuration
func (tp *TransactionParticipant) HandlePrepare(prepare *Prepare, reply *PrepareOK) error {
    if prepare.ConfigNum != tp.config.ConfigNum {
        reply.Err = ErrConfigMismatch
        return nil
    }
    // ... rest of implementation
}

// WRONG: Not handling duplicate transactions
func (tp *TransactionParticipant) HandleCommit(commit *Commit, reply *CommitOK) error {
    // Missing attempt number validation!
    txn := commit.Transaction
    // ... rest of implementation
}

// CORRECT: Check attempt numbers
func (tp *TransactionParticipant) HandleCommit(commit *Commit, reply *CommitOK) error {
    state := tp.txnManager.transactions[txnID]
    if state.AttemptNum != commit.AttemptNum {
        reply.Err = ErrAttemptMismatch
        return nil
    }
    // ... rest of implementation
}
```

### Performance Considerations

1. **Minimize lock duration**: Lock keys only during transaction execution
2. **Use efficient coordinator selection**: Choose coordinator based on group ID
3. **Implement proper timeouts**: Set reasonable timeouts for 2PC phases
4. **Batch operations**: Group related operations when possible
5. **Monitor transaction latency**: Track 2PC performance

### Debugging Tips

1. **Add structured logging**: Include transaction IDs and attempt numbers
2. **Use consistent log levels**: Debug, Info, Warn, Error
3. **Include correlation IDs**: Track transactions across components
4. **Monitor lock contention**: Log lock acquisition and release
5. **Test with delays**: Add artificial delays to find race conditions

### Running Tests

```bash
# Run transaction tests
go test -v -run TestTransaction ./...

# Run with race detection
go test -race -run TestTransaction ./...

# Run specific test
go test -v -run TestDistributedTransaction ./...

# Run with coverage
go test -cover -run TestTransaction ./...

# Run benchmarks
go test -bench=. -run TestTransaction ./...
```

### Next Steps

1. Implement basic transaction types (MultiGet, MultiPut, Swap)
2. Add coordinator logic with 2PC
3. Implement participant logic with locking
4. Add abort handling and retry logic
5. Implement transaction client
6. Add comprehensive testing
7. Performance optimization

This foundation will prepare you for implementing a complete, fault-tolerant distributed transaction system that can handle cross-shard operations while maintaining ACID properties through Two-Phase Commit.
