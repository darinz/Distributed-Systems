# Lab 2A: Primary-Backup Replication with ViewService

## Overview

This lab implements a primary-backup replication system using a ViewService to manage server roles and handle failures. You'll build a fault-tolerant key-value store where one server acts as the primary (handling all client requests) and another as the backup (replicating state for failover).

## Architecture Components

### System Players

- **ViewServer**: Central coordinator that manages server roles and view transitions
- **Primary**: Active server handling all client requests
- **Backup**: Standby server replicating primary's state
- **Reserve Servers**: Additional servers waiting to become backup
- **Client**: Application making requests to the system

### System Flow

```
Client → Primary → ViewServer
         ↓
       Backup
```

## Core Concepts

### What is a View?

A **View** represents the current configuration of the system:
- Contains the identity of the primary and backup servers
- Has an increasing view number (sequence)
- Servers contact the ViewServer to learn current view
- View transitions only occur when the current primary acknowledges the current view (prevents split-brain)
- ViewServer detects failures through heartbeat monitoring

**Important**: The ViewServer is a single point of failure, but this is acceptable for this lab.

## Go Implementation Framework

### Basic Types and Interfaces

```go
// NodeID represents a unique server identifier
type NodeID string

// View represents the current system configuration
type View struct {
    ViewNum int64   `json:"view_num"`
    Primary NodeID  `json:"primary"`
    Backup  NodeID  `json:"backup"`
}

// ViewService manages server roles and view transitions
type ViewService interface {
    Ping(args *PingArgs, reply *PingReply) error
    Get(args *GetArgs, reply *GetReply) error
}

// Server represents a primary-backup server
type Server interface {
    Ping(args *PingArgs, reply *PingReply) error
    Get(args *GetArgs, reply *GetReply) error
    Put(args *PutArgs, reply *PutReply) error
    Append(args *AppendArgs, reply *AppendReply) error
}
```

### Message Types

```go
// Ping message for heartbeat and view synchronization
type PingArgs struct {
    ServerID NodeID `json:"server_id"`
    ViewNum  int64  `json:"view_num"`
}

type PingReply struct {
    View View `json:"view"`
}

// Get current view
type GetArgs struct{}

type GetReply struct {
    View View `json:"view"`
}

// Client operations
type PutArgs struct {
    Key   string `json:"key"`
    Value string `json:"value"`
    View  View   `json:"view"`
}

type PutReply struct {
    Err Err `json:"err"`
}

type GetArgs struct {
    Key  string `json:"key"`
    View View   `json:"view"`
}

type GetReply struct {
    Value string `json:"value"`
    Err   Err    `json:"err"`
}

type AppendArgs struct {
    Key   string `json:"key"`
    Value string `json:"value"`
    View  View   `json:"view"`
}

type AppendReply struct {
    Err Err `json:"err"`
}
```
## Ping and Heartbeat System

### Heartbeat Messages

Ping RPCs serve dual purposes:
1. **Liveness confirmation**: Server is alive and responding
2. **View synchronization**: Server reports the most recent view it knows

The ViewServer uses this information to:
- Detect server failures
- Synchronize view state across servers
- Make view transition decisions

### Failure Detection

**Critical Rule**: Do NOT store timestamps on the ViewServer! This must be deterministic for search tests.

Instead, use a **ping counter approach**:

```go
type ViewServer struct {
    mu           sync.RWMutex
    currentView  View
    servers      map[NodeID]*ServerState
    pingTimeout  int // Number of ping cycles before considering server dead
}

type ServerState struct {
    lastPingCycle int64
    isAlive       bool
    viewNum       int64
}

// Ping cycle counter (incremented periodically)
var pingCycle int64

func (vs *ViewServer) Ping(args *PingArgs, reply *PingReply) error {
    vs.mu.Lock()
    defer vs.mu.Unlock()
    
    // Update server state
    if serverState, exists := vs.servers[args.ServerID]; exists {
        serverState.lastPingCycle = pingCycle
        serverState.isAlive = true
        serverState.viewNum = args.ViewNum
    } else {
        vs.servers[args.ServerID] = &ServerState{
            lastPingCycle: pingCycle,
            isAlive:       true,
            viewNum:       args.ViewNum,
        }
    }
    
    // Check for view transitions
    vs.checkViewTransitions()
    
    // Return current view
    reply.View = vs.currentView
    return nil
}

// Called periodically to detect failures
func (vs *ViewServer) checkFailures() {
    vs.mu.Lock()
    defer vs.mu.Unlock()
    
    for serverID, state := range vs.servers {
        if pingCycle-state.lastPingCycle > int64(vs.pingTimeout) {
            state.isAlive = false
            log.Printf("Server %s marked as dead", serverID)
        }
    }
    
    vs.checkViewTransitions()
}
```

### Ping Cycle Management

```go
// Start ping cycle timer
func (vs *ViewServer) startPingCycle() {
    ticker := time.NewTicker(100 * time.Millisecond)
    go func() {
        for {
            select {
            case <-ticker.C:
                atomic.AddInt64(&pingCycle, 1)
                vs.checkFailures()
            case <-vs.done:
                ticker.Stop()
                return
            }
        }
    }()
}
```
## View Transition Logic

### View Transition Rules

View transitions follow a strict sequence:

1. **Initial State**: `STARTUP_VIEWNUM = 0` with `{null, null}`
2. **First Server**: When server A pings, transition to `INITIAL_VIEWNUM = 1` with `{primary=A, null}`
3. **Add Backup**: When server B pings and primary A has acknowledged view 1, transition to view 2 with `{primary=A, backup=B}`

### Critical Rule: Primary Acknowledgment

**Only move to a new view (i+1) if the primary of view (i) has acknowledged view (i)!**

This prevents split-brain scenarios where multiple servers think they're primary.

### View Transition Implementation

```go
const (
    STARTUP_VIEWNUM  = 0
    INITIAL_VIEWNUM  = 1
)

func (vs *ViewServer) checkViewTransitions() {
    currentView := vs.currentView
    
    switch {
    case currentView.ViewNum == STARTUP_VIEWNUM:
        // Initial state: {null, null}
        vs.handleStartupView()
        
    case currentView.ViewNum == INITIAL_VIEWNUM:
        // View 1: {primary=A, null}
        vs.handleInitialView()
        
    default:
        // View 2+: {primary=A, backup=B}
        vs.handleNormalView()
    }
}

func (vs *ViewServer) handleStartupView() {
    // Find first alive server to become primary
    for serverID, state := range vs.servers {
        if state.isAlive {
            vs.currentView = View{
                ViewNum: INITIAL_VIEWNUM,
                Primary: serverID,
                Backup:  "",
            }
            log.Printf("Transition to view %d: primary=%s", INITIAL_VIEWNUM, serverID)
            return
        }
    }
}

func (vs *ViewServer) handleInitialView() {
    primaryID := vs.currentView.Primary
    primaryState, exists := vs.servers[primaryID]
    
    // Check if primary has acknowledged this view
    if !exists || !primaryState.isAlive || primaryState.viewNum < INITIAL_VIEWNUM {
        // Primary failed or hasn't acknowledged - can't transition
        return
    }
    
    // Primary has acknowledged, look for backup
    for serverID, state := range vs.servers {
        if state.isAlive && serverID != primaryID {
            vs.currentView = View{
                ViewNum: INITIAL_VIEWNUM + 1,
                Primary: primaryID,
                Backup:  serverID,
            }
            log.Printf("Transition to view %d: primary=%s, backup=%s", 
                INITIAL_VIEWNUM+1, primaryID, serverID)
            return
        }
    }
}

func (vs *ViewServer) handleNormalView() {
    primaryID := vs.currentView.Primary
    backupID := vs.currentView.Backup
    
    primaryState, primaryExists := vs.servers[primaryID]
    backupState, backupExists := vs.servers[backupID]
    
    // Check if primary has acknowledged current view
    if !primaryExists || !primaryState.isAlive || 
       primaryState.viewNum < vs.currentView.ViewNum {
        // Primary failed - promote backup
        if backupExists && backupState.isAlive {
            vs.promoteBackupToPrimary(backupID)
        }
        return
    }
    
    // Check if backup failed
    if backupID != "" && (!backupExists || !backupState.isAlive) {
        vs.findNewBackup(primaryID)
    }
}

func (vs *ViewServer) promoteBackupToPrimary(backupID NodeID) {
    // Find new backup if available
    var newBackup NodeID
    for serverID, state := range vs.servers {
        if state.isAlive && serverID != backupID {
            newBackup = serverID
            break
        }
    }
    
    vs.currentView = View{
        ViewNum: vs.currentView.ViewNum + 1,
        Primary: backupID,
        Backup:  newBackup,
    }
    log.Printf("Promoted backup %s to primary in view %d", backupID, vs.currentView.ViewNum)
}

func (vs *ViewServer) findNewBackup(primaryID NodeID) {
    for serverID, state := range vs.servers {
        if state.isAlive && serverID != primaryID {
            vs.currentView = View{
                ViewNum: vs.currentView.ViewNum + 1,
                Primary: primaryID,
                Backup:  serverID,
            }
            log.Printf("Found new backup %s in view %d", serverID, vs.currentView.ViewNum)
            return
        }
    }
    
    // No backup available
    vs.currentView = View{
        ViewNum: vs.currentView.ViewNum + 1,
        Primary: primaryID,
        Backup:  "",
    }
    log.Printf("No backup available in view %d", vs.currentView.ViewNum)
}
```

### Failure Scenarios

1. **Primary fails**: Backup becomes new primary, try to find new backup
2. **Backup fails**: Find new backup or set to null
3. **Primary fails with no backup**: Do nothing, wait for recovery
4. **Both fail**: Do nothing, system is down
## Example Call Flow

```
Time  Server 1    ViewServer    Server 2    Server 3
----  --------    ----------    --------    --------
0     Ping(0)     View 1 {S1}   
1                 View 1 {S1}   Ping(0)     
2     Ping(1)     View 2 {S1,S2}
3                 View 2 {S1,S2} Ping(1)    
4     Ping(2)     View 2 {S1,S2}
5                 View 2 {S1,S2} Ping(2)    
6     [CRASH]     
7                 View 3 {S2,null} Ping(2)
8                 View 3 {S2,null} Ping(3)
9                 View 3 {S2,null}          Ping(0)
10                View 4 {S2,S3}
```

**Key Points:**
- S1 must acknowledge view 1 before view 2 can be created
- S1 must acknowledge view 2 before view 3 can be created
- State transfer happens between ping acknowledgments

## Primary-Backup Implementation

### Client Interaction

```go
type PBClient struct {
    viewServer NodeID
    currentView View
    mu         sync.RWMutex
}

func (c *PBClient) Get(key string) (string, error) {
    // Get current view from ViewServer
    view, err := c.getCurrentView()
    if err != nil {
        return "", err
    }
    
    // Send request to primary
    args := &GetArgs{
        Key:  key,
        View: view,
    }
    
    var reply GetReply
    err = c.call(view.Primary, "Server.Get", args, &reply)
    if err != nil {
        return "", err
    }
    
    return reply.Value, reply.Err
}

func (c *PBClient) getCurrentView() (View, error) {
    c.mu.RLock()
    view := c.currentView
    c.mu.RUnlock()
    
    // Get latest view from ViewServer
    var reply GetReply
    err := c.call(c.viewServer, "ViewServer.Get", &GetArgs{}, &reply)
    if err != nil {
        return view, err
    }
    
    c.mu.Lock()
    c.currentView = reply.View
    c.mu.Unlock()
    
    return reply.View, nil
}
```

### Primary Server Implementation

```go
type PBServer struct {
    mu          sync.RWMutex
    me          NodeID
    viewServer  NodeID
    currentView View
    app         *AMOApplication
    stateTransferInProgress bool
}

func (s *PBServer) Get(args *GetArgs, reply *GetReply) error {
    s.mu.RLock()
    currentView := s.currentView
    stateTransfer := s.stateTransferInProgress
    s.mu.RUnlock()
    
    // Check if we're the primary
    if currentView.Primary != s.me {
        reply.Err = ErrWrongServer
        return nil
    }
    
    // Check view number
    if args.View.ViewNum != currentView.ViewNum {
        reply.Err = ErrWrongView
        return nil
    }
    
    // Don't process requests during state transfer
    if stateTransfer {
        reply.Err = ErrStateTransfer
        return nil
    }
    
    // Forward to backup if exists
    if currentView.Backup != "" {
        err := s.forwardToBackup(args)
        if err != nil {
            reply.Err = ErrBackupFailed
            return nil
        }
    }
    
    // Execute operation
    value, err := s.app.Get(args.Key)
    if err != nil {
        reply.Err = err
        return nil
    }
    
    reply.Value = value
    reply.Err = OK
    return nil
}

func (s *PBServer) forwardToBackup(args interface{}) error {
    s.mu.RLock()
    backupID := s.currentView.Backup
    s.mu.RUnlock()
    
    if backupID == "" {
        return nil // No backup to forward to
    }
    
    // Forward request to backup
    var reply BackupReply
    err := s.call(backupID, "Server.BackupOp", args, &reply)
    if err != nil {
        return err
    }
    
    if reply.Err != OK {
        return errors.New("backup operation failed")
    }
    
    return nil
}
```

### Backup Server Implementation

```go
func (s *PBServer) BackupOp(args interface{}, reply *BackupReply) error {
    s.mu.RLock()
    currentView := s.currentView
    s.mu.RUnlock()
    
    // Check if we're the backup
    if currentView.Backup != s.me {
        reply.Err = ErrWrongServer
        return nil
    }
    
    // Check view number
    var viewNum int64
    switch a := args.(type) {
    case *GetArgs:
        viewNum = a.View.ViewNum
    case *PutArgs:
        viewNum = a.View.ViewNum
    case *AppendArgs:
        viewNum = a.View.ViewNum
    default:
        reply.Err = ErrInvalidArgs
        return nil
    }
    
    if viewNum != currentView.ViewNum {
        reply.Err = ErrWrongView
        return nil
    }
    
    // Execute operation on backup
    switch a := args.(type) {
    case *GetArgs:
        value, err := s.app.Get(a.Key)
        if err != nil {
            reply.Err = err
        } else {
            reply.Value = value
            reply.Err = OK
        }
    case *PutArgs:
        err := s.app.Put(a.Key, a.Value)
        reply.Err = err
    case *AppendArgs:
        err := s.app.Append(a.Key, a.Value)
        reply.Err = err
    }
    
    return nil
}
```

### State Transfer Implementation

```go
type StateTransferArgs struct {
    View View
    App  *AMOApplication
}

type StateTransferReply struct {
    Err Err
}

func (s *PBServer) StateTransfer(args *StateTransferArgs, reply *StateTransferReply) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    // Check if we're the backup
    if s.currentView.Backup != s.me {
        reply.Err = ErrWrongServer
        return nil
    }
    
    // Check view number
    if args.View.ViewNum != s.currentView.ViewNum {
        reply.Err = ErrWrongView
        return nil
    }
    
    // Replace application state
    s.app = args.App
    
    reply.Err = OK
    return nil
}

func (s *PBServer) initiateStateTransfer(backupID NodeID) {
    s.mu.Lock()
    s.stateTransferInProgress = true
    s.mu.Unlock()
    
    // Send entire application state to backup
    args := &StateTransferArgs{
        View: s.currentView,
        App:  s.app,
    }
    
    var reply StateTransferReply
    err := s.call(backupID, "Server.StateTransfer", args, &reply)
    
    s.mu.Lock()
    s.stateTransferInProgress = false
    s.mu.Unlock()
    
    if err != nil || reply.Err != OK {
        log.Printf("State transfer failed: %v", err)
    } else {
        log.Printf("State transfer completed successfully")
    }
}
```

## Critical Rules

1. **Primary in view i+1 must have been backup or primary in view i**
2. **Primary must wait for backup to accept/execute each op before doing op and replying to client**
3. **Backup must accept forwarded requests only if view is correct**
4. **Non-primary must reject client requests**
5. **Every operation must be before or after state transfer**

## Testing Strategy

```go
func TestViewService(t *testing.T) {
    // Create ViewServer
    vs := NewViewServer()
    go vs.Start()
    defer vs.Stop()
    
    // Create servers
    servers := make(map[NodeID]*PBServer)
    for i := 0; i < 3; i++ {
        serverID := NodeID(fmt.Sprintf("server%d", i))
        servers[serverID] = NewPBServer(serverID, vs.GetNodeID())
        go servers[serverID].Start()
    }
    
    // Test view transitions
    testInitialView(t, vs, servers)
    testBackupAddition(t, vs, servers)
    testPrimaryFailure(t, vs, servers)
    testBackupFailure(t, vs, servers)
}

func testInitialView(t *testing.T, vs *ViewServer, servers map[NodeID]*PBServer) {
    // First server should become primary
    time.Sleep(200 * time.Millisecond)
    
    view := vs.GetCurrentView()
    assert.Equal(t, int64(1), view.ViewNum)
    assert.NotEmpty(t, view.Primary)
    assert.Empty(t, view.Backup)
}
```

## Best Practices and Common Pitfalls

### Do's

1. **Always check view numbers**: Every message should include and validate view numbers
2. **Use proper synchronization**: Protect shared state with mutexes
3. **Handle state transfer carefully**: Block operations during state transfer
4. **Implement proper error handling**: Return meaningful error codes
5. **Use deterministic failure detection**: Don't rely on timestamps
6. **Test with network partitions**: Simulate message drops and delays
7. **Log important state changes**: Include view transitions and failures
8. **Validate server roles**: Always check if you're primary/backup before processing

### Don'ts

1. **Don't store timestamps on ViewServer**: Use ping counters instead
2. **Don't process requests during state transfer**: Block until transfer completes
3. **Don't ignore view number mismatches**: Always validate view numbers
4. **Don't forget to forward to backup**: Primary must wait for backup acknowledgment
5. **Don't allow split-brain**: Only transition views when primary acknowledges
6. **Don't use blocking operations in handlers**: Use channels and goroutines
7. **Don't forget to handle network failures**: Implement retry logic
8. **Don't assume message ordering**: Messages may arrive out of order

### Common Implementation Mistakes

```go
// WRONG: Storing timestamps
type ViewServer struct {
    lastPingTime map[NodeID]time.Time // DON'T DO THIS!
}

// CORRECT: Using ping counters
type ViewServer struct {
    lastPingCycle map[NodeID]int64 // DO THIS!
}

// WRONG: Not checking view numbers
func (s *PBServer) Get(args *GetArgs, reply *GetReply) error {
    // Missing view number check!
    value, err := s.app.Get(args.Key)
    reply.Value = value
    return err
}

// CORRECT: Always validate view numbers
func (s *PBServer) Get(args *GetArgs, reply *GetReply) error {
    if args.View.ViewNum != s.currentView.ViewNum {
        reply.Err = ErrWrongView
        return nil
    }
    // ... rest of implementation
}

// WRONG: Not waiting for backup
func (s *PBServer) Put(args *PutArgs, reply *PutReply) error {
    err := s.app.Put(args.Key, args.Value)
    reply.Err = err
    return nil // Missing backup forwarding!
}

// CORRECT: Forward to backup first
func (s *PBServer) Put(args *PutArgs, reply *PutReply) error {
    // Forward to backup first
    if err := s.forwardToBackup(args); err != nil {
        reply.Err = ErrBackupFailed
        return nil
    }
    
    // Then execute on primary
    err := s.app.Put(args.Key, args.Value)
    reply.Err = err
    return nil
}
```

### Performance Considerations

1. **Minimize lock contention**: Use RWMutex for read-heavy operations
2. **Batch operations**: Group related operations together
3. **Use connection pooling**: Reuse network connections
4. **Implement backpressure**: Handle overload gracefully
5. **Monitor resource usage**: Track memory and CPU usage

### Debugging Tips

1. **Add structured logging**: Include view numbers and server IDs
2. **Use consistent log levels**: Debug, Info, Warn, Error
3. **Include correlation IDs**: Track requests across servers
4. **Monitor view transitions**: Log all view changes
5. **Test with delays**: Add artificial delays to find race conditions

### Running Tests

```bash
# Run ViewService tests
go test -v -run TestViewService ./...

# Run with race detection
go test -race -run TestViewService ./...

# Run specific test
go test -v -run TestPrimaryFailure ./...

# Run with coverage
go test -cover -run TestViewService ./...

# Run benchmarks
go test -bench=. -run TestViewService ./...
```

### Next Steps

1. Implement the ViewService with proper view transitions
2. Build the primary-backup replication system
3. Add state transfer mechanisms
4. Implement comprehensive testing
5. Add monitoring and logging
6. Performance optimization

This foundation will prepare you for implementing the complete primary-backup system with fault tolerance and consistency guarantees.
```