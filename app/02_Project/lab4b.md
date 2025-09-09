# Lab 4B: Sharded Key-Value Store

## Overview

This lab implements a sharded key-value store that uses the ShardMaster from Lab 4A to manage shard assignments and configurations. The system consists of multiple Paxos replica groups, each responsible for a subset of shards, with dynamic reconfiguration and shard movement capabilities.

## Architecture Overview

### System Components

```
ShardMaster (Paxos Group)
    ↓ (Configuration)
ShardStoreClient → ShardStoreServer (Paxos Replica Group)
                      ↓
                  Paxos Subnode
```

### Key Components

- **ShardMaster**: Manages shard-to-group mappings (from Lab 4A)
- **ShardStoreServer**: Handles key-value operations for assigned shards
- **ShardStoreClient**: Routes requests to appropriate replica groups
- **Paxos Subnode**: Provides consensus within each replica group

## ShardStore Architecture

### Core Types and Interfaces

```go
// ShardStore command types
type ShardStoreCommand interface {
    GetType() string
    GetConfigNum() int64
}

// Single key operations
type SingleKeyCommand struct {
    Op       string `json:"op"`        // "Get", "Put", "Append"
    Key      string `json:"key"`
    Value    string `json:"value,omitempty"`
    ConfigNum int64 `json:"config_num"`
    ClientID NodeID `json:"client_id"`
    SeqNum   int64  `json:"seq_num"`
}

// Shard movement commands
type ShardMove struct {
    ConfigNum int64                    `json:"config_num"`
    Shard     int                      `json:"shard"`
    Data      map[string]string        `json:"data"`
    ClientSeq map[NodeID]int64         `json:"client_seq"`
}

type ShardMoveAck struct {
    ConfigNum int64 `json:"config_num"`
    Shard     int   `json:"shard"`
}

// Configuration change command
type NewConfig struct {
    ConfigNum int64     `json:"config_num"`
    Config    Configuration `json:"config"`
}

// ShardStore request/reply
type ShardStoreRequest struct {
    Command ShardStoreCommand `json:"command"`
}

type ShardStoreReply struct {
    Err   Err    `json:"err"`
    Value string `json:"value,omitempty"`
}
```

### ShardStore Server Implementation

```go
type ShardStoreServer struct {
    mu              sync.RWMutex
    me              NodeID
    shardMasters    []NodeID
    paxos           *MultiPaxosNode
    config          Configuration
    configNum       int64
    shards          map[int]map[string]string  // shard -> key-value store
    clientSeq       map[NodeID]int64           // client -> last seq num
    shardMoves      map[int]bool               // shard -> moved
    shardAcks       map[int]bool               // shard -> acked
    done            chan struct{}
}

func NewShardStoreServer(me NodeID, shardMasters []NodeID) *ShardStoreServer {
    ss := &ShardStoreServer{
        me:           me,
        shardMasters: shardMasters,
        shards:       make(map[int]map[string]string),
        clientSeq:    make(map[NodeID]int64),
        shardMoves:   make(map[int]bool),
        shardAcks:    make(map[int]bool),
        done:         make(chan struct{}),
    }
    
    // Initialize Paxos
    ss.paxos = NewMultiPaxosNode(me, []NodeID{}) // Will be set during init
    
    return ss
}

func (ss *ShardStoreServer) Start() error {
    go ss.paxos.Start()
    go ss.configPuller()
    return nil
}

func (ss *ShardStoreServer) Stop() {
    close(ss.done)
    ss.paxos.Stop()
}
```

## Command Processing Framework

### Generic Command Processing

```go
func (ss *ShardStoreServer) process(command ShardStoreCommand, replicated bool) {
    switch cmd := command.(type) {
    case *SingleKeyCommand:
        ss.processSingleKeyCommand(cmd, replicated)
    case *ShardMove:
        ss.processShardMove(cmd, replicated)
    case *ShardMoveAck:
        ss.processShardMoveAck(cmd, replicated)
    case *NewConfig:
        ss.processNewConfig(cmd, replicated)
    default:
        log.Printf("Unknown command type: %T", command)
    }
}

func (ss *ShardStoreServer) processSingleKeyCommand(cmd *SingleKeyCommand, replicated bool) {
    if !replicated {
        // Propose command through Paxos
        ss.paxos.ProposeCommand(cmd)
        return
    }
    
    // Command is replicated, now execute it
    ss.executeSingleKeyCommand(cmd)
}

func (ss *ShardStoreServer) processShardMove(cmd *ShardMove, replicated bool) {
    if !replicated {
        // Propose command through Paxos
        ss.paxos.ProposeCommand(cmd)
        return
    }
    
    // Command is replicated, now process shard move
    ss.executeShardMove(cmd)
}

func (ss *ShardStoreServer) processShardMoveAck(cmd *ShardMoveAck, replicated bool) {
    if !replicated {
        // Propose command through Paxos
        ss.paxos.ProposeCommand(cmd)
        return
    }
    
    // Command is replicated, now process ack
    ss.executeShardMoveAck(cmd)
}

func (ss *ShardStoreServer) processNewConfig(cmd *NewConfig, replicated bool) {
    if !replicated {
        // Propose command through Paxos
        ss.paxos.ProposeCommand(cmd)
        return
    }
    
    // Command is replicated, now process config change
    ss.executeNewConfig(cmd)
}
```

### Message Handlers

```go
func (ss *ShardStoreServer) HandleShardStoreRequest(req *ShardStoreRequest, reply *ShardStoreReply) error {
    // Validate request
    if req.Command == nil {
        reply.Err = ErrInvalidCommand
        return nil
    }
    
    // Check if we're in the correct configuration
    if !ss.isConfigActive() {
        reply.Err = ErrConfigNotActive
        return nil
    }
    
    // Process command (not replicated yet)
    ss.process(req.Command, false)
    
    // Note: Reply will be sent after Paxos decision
    return nil
}

func (ss *ShardStoreServer) HandlePaxosDecision(decision PaxosDecision) {
    // Process command (now replicated)
    ss.process(decision.Command, true)
}
```

## Single Key Operations

### AMO Application Implementation

```go
type AMOApplication struct {
    mu        sync.RWMutex
    data      map[string]string  // key-value store
    clientSeq map[NodeID]int64   // client -> last sequence number
}

func NewAMOApplication() *AMOApplication {
    return &AMOApplication{
        data:      make(map[string]string),
        clientSeq: make(map[NodeID]int64),
    }
}

func (amo *AMOApplication) Execute(cmd *SingleKeyCommand) (string, Err) {
    amo.mu.Lock()
    defer amo.mu.Unlock()
    
    // Check for duplicate request
    if lastSeq, exists := amo.clientSeq[cmd.ClientID]; exists && cmd.SeqNum <= lastSeq {
        // Duplicate request, return cached result
        return amo.getCachedResult(cmd), OK
    }
    
    // Execute command
    var result string
    var err Err
    
    switch cmd.Op {
    case "Get":
        result, err = amo.get(cmd.Key)
    case "Put":
        err = amo.put(cmd.Key, cmd.Value)
    case "Append":
        result, err = amo.append(cmd.Key, cmd.Value)
    default:
        err = ErrInvalidOperation
    }
    
    // Update client sequence number
    amo.clientSeq[cmd.ClientID] = cmd.SeqNum
    
    return result, err
}

func (amo *AMOApplication) get(key string) (string, Err) {
    if value, exists := amo.data[key]; exists {
        return value, OK
    }
    return "", ErrNoKey
}

func (amo *AMOApplication) put(key, value string) Err {
    amo.data[key] = value
    return OK
}

func (amo *AMOApplication) append(key, value string) (string, Err) {
    existing := amo.data[key]
    amo.data[key] = existing + value
    return existing, OK
}
```

### Single Key Command Execution

```go
func (ss *ShardStoreServer) executeSingleKeyCommand(cmd *SingleKeyCommand) {
    ss.mu.Lock()
    defer ss.mu.Unlock()
    
    // Check if we're responsible for this key's shard
    shard := keyToShard(cmd.Key)
    if !ss.isResponsibleForShard(shard) {
        log.Printf("Not responsible for shard %d (key: %s)", shard, cmd.Key)
        return
    }
    
    // Check for duplicate request
    if lastSeq, exists := ss.clientSeq[cmd.ClientID]; exists && cmd.SeqNum <= lastSeq {
        // Duplicate request, ignore
        return
    }
    
    // Execute command
    var result string
    var err Err
    
    switch cmd.Op {
    case "Get":
        result, err = ss.get(cmd.Key, shard)
    case "Put":
        err = ss.put(cmd.Key, cmd.Value, shard)
    case "Append":
        result, err = ss.append(cmd.Key, cmd.Value, shard)
    default:
        err = ErrInvalidOperation
    }
    
    // Update client sequence number
    ss.clientSeq[cmd.ClientID] = cmd.SeqNum
    
    // Send reply to client
    ss.sendReply(cmd.ClientID, result, err)
}

func (ss *ShardStoreServer) get(key string, shard int) (string, Err) {
    if shardData, exists := ss.shards[shard]; exists {
        if value, exists := shardData[key]; exists {
            return value, OK
        }
    }
    return "", ErrNoKey
}

func (ss *ShardStoreServer) put(key, value string, shard int) Err {
    if ss.shards[shard] == nil {
        ss.shards[shard] = make(map[string]string)
    }
    ss.shards[shard][key] = value
    return OK
}

func (ss *ShardStoreServer) append(key, value string, shard int) (string, Err) {
    if ss.shards[shard] == nil {
        ss.shards[shard] = make(map[string]string)
    }
    
    existing := ss.shards[shard][key]
    ss.shards[shard][key] = existing + value
    return existing, OK
}
```

## Configuration Management

### Configuration Pulling

```go
func (ss *ShardStoreServer) configPuller() {
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            ss.pullConfig()
        case <-ss.done:
            return
        }
    }
}

func (ss *ShardStoreServer) pullConfig() {
    // Query ShardMaster for current configuration
    queryCmd := &QueryCommand{ConfigNum: -1}
    
    for _, shardMaster := range ss.shardMasters {
        var reply ConfigReply
        err := ss.call(shardMaster, "ShardMaster.Query", queryCmd, &reply)
        if err == nil && reply.Err == OK {
            if reply.Config.ConfigNum > ss.configNum {
                // New configuration available
                newConfigCmd := &NewConfig{
                    ConfigNum: reply.Config.ConfigNum,
                    Config:    reply.Config,
                }
                ss.process(newConfigCmd, false)
            }
            break
        }
    }
}
```

### Configuration Change Processing

```go
func (ss *ShardStoreServer) executeNewConfig(cmd *NewConfig) {
    ss.mu.Lock()
    defer ss.mu.Unlock()
    
    oldConfig := ss.config
    newConfig := cmd.Config
    
    // Update configuration
    ss.config = newConfig
    ss.configNum = newConfig.ConfigNum
    
    // Determine shard movements needed
    shardsToSend := ss.getShardsToSend(oldConfig, newConfig)
    shardsToReceive := ss.getShardsToReceive(oldConfig, newConfig)
    
    // Send shards to new owners
    for _, shard := range shardsToSend {
        ss.sendShard(shard, newConfig)
    }
    
    // Wait for shards to be received
    ss.waitForShardMoves(shardsToReceive)
    
    log.Printf("Configuration changed to %d", newConfig.ConfigNum)
}

func (ss *ShardStoreServer) getShardsToSend(oldConfig, newConfig Configuration) []int {
    var shards []int
    
    for shard := 0; shard < NShards; shard++ {
        oldGroup := oldConfig.Shards[shard]
        newGroup := newConfig.Shards[shard]
        
        if oldGroup != newGroup && ss.isInGroup(oldGroup, newConfig) {
            shards = append(shards, shard)
        }
    }
    
    return shards
}

func (ss *ShardStoreServer) getShardsToReceive(oldConfig, newConfig Configuration) []int {
    var shards []int
    
    for shard := 0; shard < NShards; shard++ {
        oldGroup := oldConfig.Shards[shard]
        newGroup := newConfig.Shards[shard]
        
        if oldGroup != newGroup && ss.isInGroup(newGroup, newConfig) {
            shards = append(shards, shard)
        }
    }
    
    return shards
}
```

## Shard Movement

### Sending Shards

```go
func (ss *ShardStoreServer) sendShard(shard int, config Configuration) {
    newGroup := config.Shards[shard]
    newGroupServers := config.Groups[newGroup]
    
    // Prepare shard data
    shardData := make(map[string]string)
    if data, exists := ss.shards[shard]; exists {
        for k, v := range data {
            shardData[k] = v
        }
    }
    
    // Prepare client sequence data
    clientSeq := make(map[NodeID]int64)
    for client, seq := range ss.clientSeq {
        clientSeq[client] = seq
    }
    
    // Create shard move command
    shardMove := &ShardMove{
        ConfigNum: config.ConfigNum,
        Shard:     shard,
        Data:      shardData,
        ClientSeq: clientSeq,
    }
    
    // Send to all servers in new group
    for _, server := range newGroupServers {
        go func(srv NodeID) {
            var reply ShardStoreReply
            err := ss.call(srv, "ShardStoreServer.HandleShardMove", shardMove, &reply)
            if err == nil && reply.Err == OK {
                ss.markShardSent(shard)
            }
        }(server)
    }
}

func (ss *ShardStoreServer) HandleShardMove(cmd *ShardMove, reply *ShardStoreReply) error {
    // Process shard move command
    ss.process(cmd, false)
    
    // Send acknowledgment
    ack := &ShardMoveAck{
        ConfigNum: cmd.ConfigNum,
        Shard:     cmd.Shard,
    }
    ss.process(ack, false)
    
    reply.Err = OK
    return nil
}
```

### Receiving Shards

```go
func (ss *ShardStoreServer) executeShardMove(cmd *ShardMove) {
    ss.mu.Lock()
    defer ss.mu.Unlock()
    
    // Check if we should receive this shard
    if !ss.shouldReceiveShard(cmd.Shard, cmd.ConfigNum) {
        return
    }
    
    // Update shard data
    if ss.shards[cmd.Shard] == nil {
        ss.shards[cmd.Shard] = make(map[string]string)
    }
    
    for key, value := range cmd.Data {
        ss.shards[cmd.Shard][key] = value
    }
    
    // Update client sequence numbers
    for client, seq := range cmd.ClientSeq {
        if currentSeq, exists := ss.clientSeq[client]; !exists || seq > currentSeq {
            ss.clientSeq[client] = seq
        }
    }
    
    log.Printf("Received shard %d with %d keys", cmd.Shard, len(cmd.Data))
}

func (ss *ShardStoreServer) executeShardMoveAck(cmd *ShardMoveAck) {
    ss.mu.Lock()
    defer ss.mu.Unlock()
    
    ss.shardAcks[cmd.Shard] = true
    log.Printf("Received ack for shard %d", cmd.Shard)
}
```

## ShardStore Client

### Client Implementation

```go
type ShardStoreClient struct {
    mu           sync.RWMutex
    me           NodeID
    shardMasters []NodeID
    config       Configuration
    configNum    int64
    seq          int64
}

func NewShardStoreClient(me NodeID, shardMasters []NodeID) *ShardStoreClient {
    return &ShardStoreClient{
        me:           me,
        shardMasters: shardMasters,
        seq:          0,
    }
}

func (sc *ShardStoreClient) Get(key string) (string, error) {
    return sc.executeCommand("Get", key, "")
}

func (sc *ShardStoreClient) Put(key, value string) error {
    _, err := sc.executeCommand("Put", key, value)
    return err
}

func (sc *ShardStoreClient) Append(key, value string) (string, error) {
    return sc.executeCommand("Append", key, value)
}

func (sc *ShardStoreClient) executeCommand(op, key, value string) (string, error) {
    // Get current configuration
    config, err := sc.getCurrentConfig()
    if err != nil {
        return "", err
    }
    
    // Determine which group handles this key
    shard := keyToShard(key)
    group := config.Shards[shard]
    
    if group == 0 {
        return "", errors.New("no group assigned to shard")
    }
    
    // Get servers for this group
    servers := config.Groups[group]
    if len(servers) == 0 {
        return "", errors.New("no servers in group")
    }
    
    // Create command
    sc.mu.Lock()
    sc.seq++
    seq := sc.seq
    sc.mu.Unlock()
    
    cmd := &SingleKeyCommand{
        Op:        op,
        Key:       key,
        Value:     value,
        ConfigNum: config.ConfigNum,
        ClientID:  sc.me,
        SeqNum:    seq,
    }
    
    // Send to all servers in group
    for _, server := range servers {
        var reply ShardStoreReply
        err := sc.call(server, "ShardStoreServer.HandleShardStoreRequest", 
            &ShardStoreRequest{Command: cmd}, &reply)
        if err == nil {
            if reply.Err == OK {
                return reply.Value, nil
            } else if reply.Err == ErrWrongGroup {
                // Configuration changed, retry
                return sc.executeCommand(op, key, value)
            }
        }
    }
    
    return "", errors.New("all servers failed")
}

func (sc *ShardStoreClient) getCurrentConfig() (Configuration, error) {
    sc.mu.RLock()
    config := sc.config
    configNum := sc.configNum
    sc.mu.RUnlock()
    
    // Query ShardMaster for latest configuration
    queryCmd := &QueryCommand{ConfigNum: -1}
    
    for _, shardMaster := range sc.shardMasters {
        var reply ConfigReply
        err := sc.call(shardMaster, "ShardMaster.Query", queryCmd, &reply)
        if err == nil && reply.Err == OK {
            sc.mu.Lock()
            sc.config = reply.Config
            sc.configNum = reply.Config.ConfigNum
            sc.mu.Unlock()
            return reply.Config, nil
        }
    }
    
    return config, errors.New("failed to get configuration")
}
```

## Testing Strategy

### Unit Testing

```go
func TestShardStoreServer(t *testing.T) {
    server := NewShardStoreServer("server1", []NodeID{"sm1", "sm2"})
    go server.Start()
    defer server.Stop()
    
    // Test single key operations
    cmd := &SingleKeyCommand{
        Op:        "Put",
        Key:       "test",
        Value:     "value",
        ConfigNum: 1,
        ClientID:  "client1",
        SeqNum:    1,
    }
    
    server.process(cmd, true) // Simulate replicated command
    
    // Verify data was stored
    value, err := server.get("test", 0)
    assert.NoError(t, err)
    assert.Equal(t, "value", value)
}

func TestShardMovement(t *testing.T) {
    // Create two servers
    server1 := NewShardStoreServer("server1", []NodeID{"sm1"})
    server2 := NewShardStoreServer("server2", []NodeID{"sm1"})
    
    go server1.Start()
    go server2.Start()
    defer server1.Stop()
    defer server2.Stop()
    
    // Add some data to server1
    cmd := &SingleKeyCommand{
        Op:        "Put",
        Key:       "test",
        Value:     "value",
        ConfigNum: 1,
        ClientID:  "client1",
        SeqNum:    1,
    }
    server1.process(cmd, true)
    
    // Simulate shard movement
    shardMove := &ShardMove{
        ConfigNum: 2,
        Shard:     0,
        Data:      map[string]string{"test": "value"},
        ClientSeq: map[NodeID]int64{"client1": 1},
    }
    
    server2.process(shardMove, true)
    
    // Verify server2 has the data
    value, err := server2.get("test", 0)
    assert.NoError(t, err)
    assert.Equal(t, "value", value)
}
```

### Integration Testing

```go
func TestShardStoreIntegration(t *testing.T) {
    // Create ShardMaster
    sm := NewShardMaster([]NodeID{"sm1", "sm2", "sm3"}, "sm1")
    go sm.Start()
    defer sm.Stop()
    
    // Create ShardStore servers
    servers := make([]*ShardStoreServer, 3)
    for i := 0; i < 3; i++ {
        serverID := NodeID(fmt.Sprintf("server%d", i))
        servers[i] = NewShardStoreServer(serverID, []NodeID{"sm1", "sm2", "sm3"})
        go servers[i].Start()
    }
    
    defer func() {
        for _, server := range servers {
            server.Stop()
        }
    }()
    
    // Create client
    client := NewShardStoreClient("client1", []NodeID{"sm1", "sm2", "sm3"})
    
    // Add replica groups
    join1 := &JoinCommand{
        GroupID: 1,
        Servers: []NodeID{"server0", "server1", "server2"},
    }
    var reply OKReply
    sm.Join(join1, &reply)
    
    // Wait for configuration to propagate
    time.Sleep(200 * time.Millisecond)
    
    // Test operations
    err := client.Put("key1", "value1")
    assert.NoError(t, err)
    
    value, err := client.Get("key1")
    assert.NoError(t, err)
    assert.Equal(t, "value1", value)
    
    result, err := client.Append("key1", "extra")
    assert.NoError(t, err)
    assert.Equal(t, "value1", result)
    
    value, err = client.Get("key1")
    assert.NoError(t, err)
    assert.Equal(t, "value1extra", value)
}
```

## Best Practices and Common Pitfalls

### Do's

1. **Use Paxos for all operations**: Ensure consistency within replica groups
2. **Check configuration validity**: Verify operations are sent to correct groups
3. **Handle duplicate requests**: Use client sequence numbers for at-most-once semantics
4. **Implement proper shard movement**: Send both data and client sequence mappings
5. **Wait for configuration stability**: Don't process operations during reconfiguration
6. **Use efficient shard distribution**: Minimize data movement during reconfiguration
7. **Test with multiple configurations**: Verify system works across config changes
8. **Handle network failures**: Implement retry logic for failed operations

### Don'ts

1. **Don't bypass Paxos**: All operations must go through consensus
2. **Don't ignore configuration changes**: Always check if config is active
3. **Don't forget client sequence tracking**: Essential for at-most-once semantics
4. **Don't process operations during reconfiguration**: Wait for stability
5. **Don't ignore shard movement**: Properly handle data transfer between groups
6. **Don't assume single replica**: Test with multiple servers per group
7. **Don't forget error handling**: Handle wrong group and timeout errors
8. **Don't skip validation**: Always validate commands before processing

### Common Implementation Mistakes

```go
// WRONG: Not checking configuration
func (ss *ShardStoreServer) executeSingleKeyCommand(cmd *SingleKeyCommand) {
    // Missing configuration check!
    ss.shards[0][cmd.Key] = cmd.Value
}

// CORRECT: Always check configuration
func (ss *ShardStoreServer) executeSingleKeyCommand(cmd *SingleKeyCommand) {
    if !ss.isConfigActive() {
        return
    }
    
    shard := keyToShard(cmd.Key)
    if !ss.isResponsibleForShard(shard) {
        return
    }
    
    // ... rest of implementation
}

// WRONG: Not handling duplicate requests
func (ss *ShardStoreServer) executeSingleKeyCommand(cmd *SingleKeyCommand) {
    // Missing duplicate check!
    ss.shards[0][cmd.Key] = cmd.Value
}

// CORRECT: Check for duplicates
func (ss *ShardStoreServer) executeSingleKeyCommand(cmd *SingleKeyCommand) {
    if lastSeq, exists := ss.clientSeq[cmd.ClientID]; exists && cmd.SeqNum <= lastSeq {
        return // Duplicate request
    }
    
    // ... rest of implementation
    ss.clientSeq[cmd.ClientID] = cmd.SeqNum
}
```

### Performance Considerations

1. **Minimize shard movements**: Use efficient redistribution algorithms
2. **Batch operations**: Group related operations when possible
3. **Use efficient data structures**: Optimize key-value store operations
4. **Implement caching**: Cache frequently accessed configurations
5. **Monitor consensus latency**: Track Paxos performance

### Debugging Tips

1. **Add structured logging**: Include configuration numbers and shard IDs
2. **Use consistent log levels**: Debug, Info, Warn, Error
3. **Include correlation IDs**: Track operations across components
4. **Monitor shard movements**: Log all shard transfer operations
5. **Test with delays**: Add artificial delays to find race conditions

### Running Tests

```bash
# Run ShardStore tests
go test -v -run TestShardStore ./...

# Run with race detection
go test -race -run TestShardStore ./...

# Run specific test
go test -v -run TestShardStoreIntegration ./...

# Run with coverage
go test -cover -run TestShardStore ./...

# Run benchmarks
go test -bench=. -run TestShardStore ./...
```

### Next Steps

1. Implement basic ShardStore server with Paxos integration
2. Add single key operations with AMO semantics
3. Implement configuration management and pulling
4. Add shard movement and reconfiguration
5. Implement ShardStore client with routing
6. Add comprehensive testing
7. Performance optimization

This foundation will prepare you for implementing a complete, fault-tolerant sharded key-value store that can dynamically handle configuration changes and shard movements while maintaining consistency and at-most-once semantics.
