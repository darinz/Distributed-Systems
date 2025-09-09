# Lab 4A: ShardMaster - Distributed Sharding Service

## Overview

This lab implements a ShardMaster service that manages sharding and configuration for a distributed key-value store. The ShardMaster uses Paxos (from Lab 3) to maintain consistent configuration state across multiple replicas, enabling dynamic load balancing and system reconfiguration.

## Lab 3 Dependencies

- **Lab 4 relies heavily on Lab 3**: ShardMaster uses Paxos for consensus
- **Multi-Paxos Implementation**: Even if Lab 3 isn't perfect, you can still work on Lab 4
- **Priority**: Finish Lab 3 first for best results

## Lab 4 Overview

### Goals

- **Reliability**: Paxos increases fault tolerance
- **Performance**: Sharding increases throughput and scalability
- **Multi-key Transactions**: Part 3 adds multi-key transaction support using two-phase commit
- **Linearizability**: Maintain consistency across all operations
- **Dynamic Load Balancing**: Similar to Amazon's DynamoDB or Google's Spanner

### Architecture

```
Client → ShardMaster → Paxos Group
         ↓
    Configuration
         ↓
    Shard Mapping
         ↓
    Replica Groups
```

## Sharding Concepts

### What is Sharding?

**Sharding** divides the keyspace into multiple groups called **shards**:

- **Parallel Processing**: Different shards can run in parallel without communication
- **Performance Scaling**: Performance increases proportionally to the number of shards
- **Fault Isolation**: Failures in one shard don't affect others
- **Load Distribution**: Spread load across multiple replica groups

### Terminology

- **ShardMaster**: Service that manages shard-to-group mappings using Paxos
- **Configuration**: Specifies which shards map to which replica groups (similar to View in Lab 2)
- **Paxos Replica Group**: Group of servers performing Paxos agreement (from Lab 3)
- **Shard**: Subset of key/value pairs (e.g., keys starting with "a-g")
- **Configuration Number**: Monotonically increasing version number

## ShardMaster Architecture

### Core Components

```go
// Configuration represents the current shard-to-group mapping
type Configuration struct {
    ConfigNum int64                    `json:"config_num"`
    Groups    map[int64][]NodeID       `json:"groups"`    // groupID -> servers
    Shards    [NShards]int64           `json:"shards"`    // shard -> groupID
}

// ShardMaster command types
type ShardMasterCommand interface {
    GetType() string
    GetConfigNum() int64
}

// Join command: Add a new replica group
type JoinCommand struct {
    GroupID int64    `json:"group_id"`
    Servers []NodeID `json:"servers"`
}

// Leave command: Remove a replica group
type LeaveCommand struct {
    GroupID int64 `json:"group_id"`
}

// Move command: Move a shard to a specific group
type MoveCommand struct {
    Shard   int   `json:"shard"`
    GroupID int64 `json:"group_id"`
}

// Query command: Get configuration
type QueryCommand struct {
    ConfigNum int64 `json:"config_num"`
}

// ShardMaster reply types
type ShardMasterReply interface {
    GetErr() Err
}

type OKReply struct {
    Err Err `json:"err"`
}

type ConfigReply struct {
    Err    Err          `json:"err"`
    Config Configuration `json:"config"`
}
```

### ShardMaster Implementation

```go
type ShardMaster struct {
    mu           sync.RWMutex
    me           NodeID
    peers        []NodeID
    paxos        *MultiPaxosNode
    configs      []Configuration
    currentConfig int64
    done         chan struct{}
}

func NewShardMaster(peers []NodeID, me NodeID) *ShardMaster {
    sm := &ShardMaster{
        me:      me,
        peers:   peers,
        configs: make([]Configuration, 1),
        done:    make(chan struct{}),
    }
    
    // Initialize with empty configuration
    sm.configs[0] = Configuration{
        ConfigNum: 0,
        Groups:    make(map[int64][]NodeID),
        Shards:    [NShards]int64{},
    }
    
    // Initialize Paxos
    sm.paxos = NewMultiPaxosNode(me, peers)
    
    return sm
}

func (sm *ShardMaster) Start() error {
    go sm.paxos.Start()
    go sm.configPuller()
    return nil
}

func (sm *ShardMaster) Stop() {
    close(sm.done)
    sm.paxos.Stop()
}
```

## Configuration Management

### Configuration Operations

```go
func (sm *ShardMaster) getCurrentConfig() Configuration {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    
    return sm.configs[sm.currentConfig]
}

func (sm *ShardMaster) createNewConfig() Configuration {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    
    // Create new configuration based on current one
    newConfig := Configuration{
        ConfigNum: sm.currentConfig + 1,
        Groups:    make(map[int64][]NodeID),
        Shards:    [NShards]int64{},
    }
    
    // Copy current groups
    for groupID, servers := range sm.configs[sm.currentConfig].Groups {
        newConfig.Groups[groupID] = make([]NodeID, len(servers))
        copy(newConfig.Groups[groupID], servers)
    }
    
    // Copy current shard assignments
    copy(newConfig.Shards[:], sm.configs[sm.currentConfig].Shards[:])
    
    return newConfig
}

func (sm *ShardMaster) addConfig(config Configuration) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    sm.configs = append(sm.configs, config)
    sm.currentConfig = config.ConfigNum
}
```

### Configuration Persistence

```go
func (sm *ShardMaster) configPuller() {
    for {
        select {
        case <-sm.done:
            return
        default:
            sm.pullConfigs()
            time.Sleep(100 * time.Millisecond)
        }
    }
}

func (sm *ShardMaster) pullConfigs() {
    // Pull configurations from Paxos log
    for {
        decided, value := sm.paxos.GetDecision(InstanceID(sm.currentConfig + 1))
        if !decided {
            break
        }
        
        if command, ok := value.(ShardMasterCommand); ok {
            sm.applyCommand(command)
        }
    }
}
```

## Join Operation

### Join Implementation

```go
func (sm *ShardMaster) Join(args *JoinCommand, reply *OKReply) error {
    // Validate input
    if args.GroupID <= 0 {
        reply.Err = ErrInvalidGroupID
        return nil
    }
    
    if len(args.Servers) == 0 {
        reply.Err = ErrEmptyServers
        return nil
    }
    
    // Check if group already exists
    currentConfig := sm.getCurrentConfig()
    if _, exists := currentConfig.Groups[args.GroupID]; exists {
        reply.Err = ErrGroupExists
        return nil
    }
    
    // Propose Join command through Paxos
    err := sm.proposeCommand(args)
    if err != nil {
        reply.Err = ErrPaxosFailed
        return nil
    }
    
    reply.Err = OK
    return nil
}

func (sm *ShardMaster) applyJoinCommand(cmd *JoinCommand) {
    // Create new configuration
    newConfig := sm.createNewConfig()
    
    // Add new group
    newConfig.Groups[cmd.GroupID] = make([]NodeID, len(cmd.Servers))
    copy(newConfig.Groups[cmd.GroupID], cmd.Servers)
    
    // Redistribute shards
    sm.redistributeShards(&newConfig)
    
    // Add new configuration
    sm.addConfig(newConfig)
    
    log.Printf("Applied Join: Group %d added with %d servers", 
        cmd.GroupID, len(cmd.Servers))
}
```

### Shard Redistribution Algorithm

```go
func (sm *ShardMaster) redistributeShards(config *Configuration) {
    groupIDs := make([]int64, 0, len(config.Groups))
    for groupID := range config.Groups {
        groupIDs = append(groupIDs, groupID)
    }
    
    if len(groupIDs) == 0 {
        // No groups, assign all shards to 0 (null)
        for i := range config.Shards {
            config.Shards[i] = 0
        }
        return
    }
    
    // Calculate shards per group
    shardsPerGroup := NShards / len(groupIDs)
    extraShards := NShards % len(groupIDs)
    
    // Distribute shards
    shardIndex := 0
    for i, groupID := range groupIDs {
        shardsForThisGroup := shardsPerGroup
        if i < extraShards {
            shardsForThisGroup++
        }
        
        for j := 0; j < shardsForThisGroup; j++ {
            config.Shards[shardIndex] = groupID
            shardIndex++
        }
    }
    
    log.Printf("Redistributed %d shards among %d groups", NShards, len(groupIDs))
}
```

## Leave Operation

### Leave Implementation

```go
func (sm *ShardMaster) Leave(args *LeaveCommand, reply *OKReply) error {
    // Validate input
    if args.GroupID <= 0 {
        reply.Err = ErrInvalidGroupID
        return nil
    }
    
    // Check if group exists
    currentConfig := sm.getCurrentConfig()
    if _, exists := currentConfig.Groups[args.GroupID]; !exists {
        reply.Err = ErrGroupNotFound
        return nil
    }
    
    // Check if this is the last group
    if len(currentConfig.Groups) <= 1 {
        reply.Err = ErrLastGroup
        return nil
    }
    
    // Propose Leave command through Paxos
    err := sm.proposeCommand(args)
    if err != nil {
        reply.Err = ErrPaxosFailed
        return nil
    }
    
    reply.Err = OK
    return nil
}

func (sm *ShardMaster) applyLeaveCommand(cmd *LeaveCommand) {
    // Create new configuration
    newConfig := sm.createNewConfig()
    
    // Remove group
    delete(newConfig.Groups, cmd.GroupID)
    
    // Redistribute shards
    sm.redistributeShards(&newConfig)
    
    // Add new configuration
    sm.addConfig(newConfig)
    
    log.Printf("Applied Leave: Group %d removed", cmd.GroupID)
}
```

## Move Operation

### Move Implementation

```go
func (sm *ShardMaster) Move(args *MoveCommand, reply *OKReply) error {
    // Validate input
    if args.Shard < 0 || args.Shard >= NShards {
        reply.Err = ErrInvalidShard
        return nil
    }
    
    if args.GroupID <= 0 {
        reply.Err = ErrInvalidGroupID
        return nil
    }
    
    // Check if group exists
    currentConfig := sm.getCurrentConfig()
    if _, exists := currentConfig.Groups[args.GroupID]; !exists {
        reply.Err = ErrGroupNotFound
        return nil
    }
    
    // Propose Move command through Paxos
    err := sm.proposeCommand(args)
    if err != nil {
        reply.Err = ErrPaxosFailed
        return nil
    }
    
    reply.Err = OK
    return nil
}

func (sm *ShardMaster) applyMoveCommand(cmd *MoveCommand) {
    // Create new configuration
    newConfig := sm.createNewConfig()
    
    // Move shard to specified group
    newConfig.Shards[cmd.Shard] = cmd.GroupID
    
    // Add new configuration
    sm.addConfig(newConfig)
    
    log.Printf("Applied Move: Shard %d moved to group %d", cmd.Shard, cmd.GroupID)
}
```

## Query Operation

### Query Implementation

```go
func (sm *ShardMaster) Query(args *QueryCommand, reply *ConfigReply) error {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    
    configNum := args.ConfigNum
    
    // Handle special cases
    if configNum == -1 || configNum >= int64(len(sm.configs)) {
        // Return latest configuration
        configNum = sm.currentConfig
    }
    
    // Return requested configuration
    if configNum >= 0 && configNum < int64(len(sm.configs)) {
        reply.Config = sm.configs[configNum]
        reply.Err = OK
    } else {
        reply.Err = ErrConfigNotFound
    }
    
    return nil
}
```

## Paxos Integration

### Command Proposing

```go
func (sm *ShardMaster) proposeCommand(cmd ShardMasterCommand) error {
    // Get next instance number
    instanceID := InstanceID(sm.currentConfig + 1)
    
    // Propose command through Paxos
    err := sm.paxos.ProposeCommand(cmd)
    if err != nil {
        return err
    }
    
    // Wait for decision
    for {
        decided, value := sm.paxos.GetDecision(instanceID)
        if decided {
            if command, ok := value.(ShardMasterCommand); ok {
                sm.applyCommand(command)
            }
            break
        }
        time.Sleep(10 * time.Millisecond)
    }
    
    return nil
}

func (sm *ShardMaster) applyCommand(cmd ShardMasterCommand) {
    switch c := cmd.(type) {
    case *JoinCommand:
        sm.applyJoinCommand(c)
    case *LeaveCommand:
        sm.applyLeaveCommand(c)
    case *MoveCommand:
        sm.applyMoveCommand(c)
    default:
        log.Printf("Unknown command type: %T", cmd)
    }
}
```

## Advanced Shard Redistribution

### Load Balancing Algorithm

```go
func (sm *ShardMaster) redistributeShardsBalanced(config *Configuration) {
    groupIDs := make([]int64, 0, len(config.Groups))
    for groupID := range config.Groups {
        groupIDs = append(groupIDs, groupID)
    }
    
    if len(groupIDs) == 0 {
        // No groups, assign all shards to 0 (null)
        for i := range config.Shards {
            config.Shards[i] = 0
        }
        return
    }
    
    // Calculate target shards per group
    targetShardsPerGroup := float64(NShards) / float64(len(groupIDs))
    
    // Distribute shards to minimize imbalance
    shardCounts := make(map[int64]int)
    for _, groupID := range groupIDs {
        shardCounts[groupID] = 0
    }
    
    // Assign shards one by one to the group with least shards
    for shard := 0; shard < NShards; shard++ {
        // Find group with minimum shards
        minGroup := groupIDs[0]
        minCount := shardCounts[minGroup]
        
        for _, groupID := range groupIDs {
            if shardCounts[groupID] < minCount {
                minGroup = groupID
                minCount = shardCounts[groupID]
            }
        }
        
        // Assign shard to this group
        config.Shards[shard] = minGroup
        shardCounts[minGroup]++
    }
    
    log.Printf("Redistributed %d shards among %d groups (target: %.2f per group)", 
        NShards, len(groupIDs), targetShardsPerGroup)
}
```

## Testing Strategy

### Unit Testing

```go
func TestShardMasterJoin(t *testing.T) {
    sm := NewShardMaster([]NodeID{"peer1", "peer2", "peer3"}, "sm1")
    go sm.Start()
    defer sm.Stop()
    
    // Test Join
    joinCmd := &JoinCommand{
        GroupID: 1,
        Servers: []NodeID{"server1", "server2", "server3"},
    }
    
    var reply OKReply
    err := sm.Join(joinCmd, &reply)
    assert.NoError(t, err)
    assert.Equal(t, OK, reply.Err)
    
    // Verify configuration
    queryCmd := &QueryCommand{ConfigNum: -1}
    var configReply ConfigReply
    err = sm.Query(queryCmd, &configReply)
    assert.NoError(t, err)
    assert.Equal(t, OK, configReply.Err)
    assert.Equal(t, int64(1), configReply.Config.ConfigNum)
    assert.Contains(t, configReply.Config.Groups, int64(1))
}

func TestShardMasterLeave(t *testing.T) {
    sm := NewShardMaster([]NodeID{"peer1", "peer2", "peer3"}, "sm1")
    go sm.Start()
    defer sm.Stop()
    
    // First add a group
    joinCmd := &JoinCommand{
        GroupID: 1,
        Servers: []NodeID{"server1", "server2"},
    }
    var reply OKReply
    sm.Join(joinCmd, &reply)
    
    // Then remove it
    leaveCmd := &LeaveCommand{GroupID: 1}
    err := sm.Leave(leaveCmd, &reply)
    assert.NoError(t, err)
    assert.Equal(t, OK, reply.Err)
    
    // Verify group is removed
    queryCmd := &QueryCommand{ConfigNum: -1}
    var configReply ConfigReply
    sm.Query(queryCmd, &configReply)
    assert.NotContains(t, configReply.Config.Groups, int64(1))
}
```

### Integration Testing

```go
func TestShardMasterIntegration(t *testing.T) {
    // Create multiple ShardMaster replicas
    peers := []NodeID{"sm1", "sm2", "sm3"}
    sms := make([]*ShardMaster, len(peers))
    
    for i, peer := range peers {
        sms[i] = NewShardMaster(peers, peer)
        go sms[i].Start()
    }
    
    defer func() {
        for _, sm := range sms {
            sm.Stop()
        }
    }()
    
    // Wait for Paxos to stabilize
    time.Sleep(200 * time.Millisecond)
    
    // Test sequence of operations
    testOperations(t, sms[0])
    
    // Verify all replicas have same configuration
    verifyConsistency(t, sms)
}

func testOperations(t *testing.T, sm *ShardMaster) {
    // Join group 1
    join1 := &JoinCommand{
        GroupID: 1,
        Servers: []NodeID{"s1", "s2", "s3"},
    }
    var reply OKReply
    err := sm.Join(join1, &reply)
    assert.NoError(t, err)
    
    // Join group 2
    join2 := &JoinCommand{
        GroupID: 2,
        Servers: []NodeID{"s4", "s5", "s6"},
    }
    err = sm.Join(join2, &reply)
    assert.NoError(t, err)
    
    // Move shard 0 to group 2
    move := &MoveCommand{
        Shard:   0,
        GroupID: 2,
    }
    err = sm.Move(move, &reply)
    assert.NoError(t, err)
    
    // Leave group 1
    leave := &LeaveCommand{GroupID: 1}
    err = sm.Leave(leave, &reply)
    assert.NoError(t, err)
}

func verifyConsistency(t *testing.T, sms []*ShardMaster) {
    // Get latest configuration from all replicas
    configs := make([]Configuration, len(sms))
    for i, sm := range sms {
        query := &QueryCommand{ConfigNum: -1}
        var reply ConfigReply
        err := sm.Query(query, &reply)
        assert.NoError(t, err)
        configs[i] = reply.Config
    }
    
    // Verify all configurations are identical
    for i := 1; i < len(configs); i++ {
        assert.Equal(t, configs[0].ConfigNum, configs[i].ConfigNum)
        assert.Equal(t, configs[0].Groups, configs[i].Groups)
        assert.Equal(t, configs[0].Shards, configs[i].Shards)
    }
}
```

## Best Practices and Common Pitfalls

### Do's

1. **Use Paxos for all configuration changes**: Ensure consistency
2. **Implement proper shard redistribution**: Minimize shard movements
3. **Validate all inputs**: Check group IDs, shard numbers, etc.
4. **Handle edge cases**: Empty groups, last group leaving, etc.
5. **Use efficient algorithms**: Minimize shard movements during redistribution
6. **Test with multiple replicas**: Verify consistency across replicas
7. **Log important operations**: Track configuration changes
8. **Handle Paxos failures**: Retry on Paxos errors

### Don'ts

1. **Don't bypass Paxos**: All configuration changes must go through Paxos
2. **Don't ignore validation**: Always validate input parameters
3. **Don't forget edge cases**: Handle empty groups and last group scenarios
4. **Don't assume single replica**: Test with multiple ShardMaster replicas
5. **Don't ignore shard redistribution**: Implement proper load balancing
6. **Don't forget error handling**: Handle Paxos failures gracefully
7. **Don't ignore consistency**: Ensure all replicas have same configuration
8. **Don't skip testing**: Test all operations thoroughly

### Common Implementation Mistakes

```go
// WRONG: Not using Paxos for configuration changes
func (sm *ShardMaster) Join(args *JoinCommand, reply *OKReply) error {
    // Directly modifying configuration without Paxos!
    sm.configs[sm.currentConfig].Groups[args.GroupID] = args.Servers
    reply.Err = OK
    return nil
}

// CORRECT: Use Paxos for all configuration changes
func (sm *ShardMaster) Join(args *JoinCommand, reply *OKReply) error {
    err := sm.proposeCommand(args)
    if err != nil {
        reply.Err = ErrPaxosFailed
        return nil
    }
    reply.Err = OK
    return nil
}

// WRONG: Not validating input
func (sm *ShardMaster) Move(args *MoveCommand, reply *OKReply) error {
    // Missing validation!
    sm.proposeCommand(args)
    reply.Err = OK
    return nil
}

// CORRECT: Always validate input
func (sm *ShardMaster) Move(args *MoveCommand, reply *OKReply) error {
    if args.Shard < 0 || args.Shard >= NShards {
        reply.Err = ErrInvalidShard
        return nil
    }
    if args.GroupID <= 0 {
        reply.Err = ErrInvalidGroupID
        return nil
    }
    // ... rest of implementation
}
```

### Performance Considerations

1. **Minimize shard movements**: Use efficient redistribution algorithms
2. **Batch operations**: Group related configuration changes
3. **Use efficient data structures**: Optimize group and shard lookups
4. **Implement caching**: Cache frequently accessed configurations
5. **Monitor Paxos performance**: Track consensus latency

### Debugging Tips

1. **Add structured logging**: Include configuration numbers and group IDs
2. **Use consistent log levels**: Debug, Info, Warn, Error
3. **Include correlation IDs**: Track operations across replicas
4. **Monitor configuration changes**: Log all Join/Leave/Move operations
5. **Test with delays**: Add artificial delays to find race conditions

### Running Tests

```bash
# Run ShardMaster tests
go test -v -run TestShardMaster ./...

# Run with race detection
go test -race -run TestShardMaster ./...

# Run specific test
go test -v -run TestShardMasterIntegration ./...

# Run with coverage
go test -cover -run TestShardMaster ./...

# Run benchmarks
go test -bench=. -run TestShardMaster ./...
```

### Next Steps

1. Implement basic ShardMaster with Paxos integration
2. Add Join, Leave, Move, and Query operations
3. Implement efficient shard redistribution algorithms
4. Add comprehensive testing
5. Performance optimization
6. Integration with sharded key-value store

This foundation will prepare you for implementing a complete, fault-tolerant sharding service that can dynamically manage replica groups and shard assignments.
