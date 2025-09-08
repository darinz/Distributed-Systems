package shardmaster

import (
	"bytes"
	crand "crypto/rand"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"sync"
	"syscall"
	"time"

	"distributed-systems/app/01_Practice-Labs/src/paxos"
)

// Debug controls whether debug output is printed
const Debug = 0

// DPrintf prints debug output if Debug is enabled
func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug > 1 {
		fmt.Printf(format, a...)
	} else if Debug > 0 {
		log.Printf(format, a...)
	}
	return
}

// ShardMaster represents a fault-tolerant shard master server.
// It manages configurations that describe which replica groups are responsible
// for which shards, and handles dynamic reconfiguration as groups join and leave.
type ShardMaster struct {
	mu         sync.Mutex        // Mutex for protecting shared state
	l          net.Listener      // Network listener for RPC connections
	me         int               // Server identifier
	dead       bool              // Flag indicating if server should shut down (for testing)
	unreliable bool              // Flag for unreliable network simulation (for testing)
	px         *paxos.Paxos      // Paxos peer for consensus

	configs  []Config // Array of configurations indexed by configuration number
	seqTried int      // Next sequence number to try for Paxos operations
	seqDone  int      // Last sequence number for which Done() was called
}

// Operation type constants for Paxos operations
const (
	Join  = "Join"  // Join operation: add a new replica group
	Leave = "Leave" // Leave operation: remove a replica group
	Move  = "Move"  // Move operation: move a shard to a specific group
	Query = "Query" // Query operation: retrieve a configuration
)

// Op represents an operation that can be agreed upon through Paxos.
// It contains all the information needed to execute a shard master operation.
type Op struct {
	Name    string   // Operation type (one of the constants above)
	OpId    uint64   // Unique operation ID for duplicate detection
	GID     int64    // Group ID (for Join, Leave, Move operations)
	Shard   int      // Shard number (for Move operations)
	Num     int      // Configuration number (for Query operations)
	Servers []string // Server addresses (for Join operations)
}

// uuid generates a 64-bit unique identifier for operation identification.
// It uses cryptographically secure random number generation to ensure uniqueness.
func (sm *ShardMaster) uuid() uint64 {
	// Generate 8 random bytes
	rbytes := make([]byte, 8)
	n, err := crand.Read(rbytes)
	for n < 8 || err != nil {
		n, err = crand.Read(rbytes)
	}

	// Convert bytes to uint64 using little-endian encoding
	var randid uint64
	binary.Read(bytes.NewReader(rbytes), binary.LittleEndian, &randid)

	return randid
}

// waitForPaxos waits for a Paxos instance to be decided.
// It uses exponential backoff to avoid busy waiting.
// Returns true if the instance was decided, false if it was already done.
func (sm *ShardMaster) waitForPaxos(seq int) bool {
	for to := 10 * time.Millisecond; ; {
		if decided, _ := sm.px.Status(seq); decided {
			return true
		}
		time.Sleep(to)
		if to < 10*time.Second {
			to *= 2
		}

		DPrintf("Waiting for paxos (seq %d) on server %d.", seq, sm.me)
	}
}

// decide attempts to assign a Paxos sequence number to the given operation.
// It tries different sequence numbers until it successfully gets agreement on the operation.
// Returns the sequence number where the operation was agreed upon.
func (sm *ShardMaster) decide(op Op) int {
	for {
		DPrintf("Waiting for lock (to decide) on server %d.", sm.me)

		sm.mu.Lock()
		seq := sm.seqTried
		sm.seqTried++
		sm.mu.Unlock()

		sm.px.Start(seq, op)
		sm.waitForPaxos(seq)

		if decided, val := sm.px.Status(seq); decided {
			if op.OpId == val.(Op).OpId {
				return seq
			}
		}
	}
}

// loadBalance implements load balancing for join/leave operations.
// It redistributes shards evenly among replica groups to balance load.
// Returns a new configuration with balanced shard assignments.
func (sm *ShardMaster) loadBalance(config Config) Config {
	// the number of shards per server; since shards might not distribute exactly evenly, we need two numbers
	q1 := NShards / len(config.Groups)
	q2 := q1 + 1

	// the number of groups with q2 or q1 shards, respectively
	// mathematically, n1 * q1 + n2 * q2 = NShards; n1 + n2 = len(Groups)
	n2 := NShards % len(config.Groups)
	n1 := len(config.Groups) - n2

	DPrintf("Load-balancer assertion: %d * %d + %d * %d = %d", n1, q1, n2, q2, NShards)

	// first pass: we count shards per group;
	// and make a list of un-assigned shards or shards assigned to too-heavy blocks
	assign := make(map[int64]int)
	unassign := []int{}

	for gid, _ := range config.Groups {
		assign[gid] = 0
	}

	for shard, gid := range config.Shards {
		if count, ok := assign[gid]; !ok {
			unassign = append(unassign, shard)
		} else if count >= q2 {
			unassign = append(unassign, shard)
		} else if n2 == 0 && count >= q1 {
			unassign = append(unassign, shard)
		} else if assign[gid]++; count+1 >= q2 {
			// if this caused the group to fill; mark off the count of numbers of groups
			n2 -= 1
		}
	}

	DPrintf("Unassigned: %d", unassign)

	// now we have a list of shards to re-assign

	// second pass: add the unassigned shards to groups without enough shards
	// we had better reach the end of our unassigned list and the end of our groups list at the same time
	shard_idx := 0
	for gid, count := range assign {
		// check if group is full
		if count >= q2 || (count >= q1 && n2 == 0) {
			continue
		}

		var q int
		if n2 > 0 {
			q = q2
			n2--
		} else {
			q = q1
		}


		// reassign from the list until the group is full
		for i := 0; i < q - count; i++ {
			if shard_idx >= len(unassign) {
				DPrintf("Error: Load-balancer had too few unassigned shards.")
				return config
			}
			config.Shards[unassign[shard_idx]] = gid
			shard_idx += 1
		}
	}

	if shard_idx < len(unassign) {
		DPrintf("Error: Load-balancer did not reassign all shards.")
		// the remaining shards get distributed over whichever gid comes up
		// (iterating over a map is not guaranteed to be consistent for the same
		// set of <key/value>s so this may distribute to more than one group)
		for ; shard_idx < len(unassign); shard_idx++ {
			for gid, _ := range config.Groups {
				config.Shards[unassign[shard_idx]] = gid
				break;
			}
		}
	}

	return config
}

// duplicateLast creates a copy of the last configuration with an incremented number.
// This is used when creating new configurations for Join, Leave, and Move operations.
func (sm *ShardMaster) duplicateLast() Config {
	prev := sm.configs[len(sm.configs)-1]

	// Copy over the previous configuration's shards
	var shards [NShards]int64
	for idx, elem := range prev.Shards {
		shards[idx] = elem
	}

	// Copy over the previous configuration's groups
	groups := make(map[int64][]string)
	for idx, elem := range prev.Groups {
		groups[idx] = elem
	}

	// Return the configuration with an incremented number
	return Config{Num: prev.Num + 1, Shards: shards, Groups: groups}
}

// doOps executes all operations up to and including the given sequence number.
// Since configurations are immutable, we don't need to cache versions for Query requests,
// but we do return the number of the last configuration.
func (sm *ShardMaster) doOps(seqEnd int) int {
	DPrintf("Waiting for lock (to do) on server %d.", sm.me)

	sm.mu.Lock()
	defer sm.mu.Unlock()

	DPrintf("Doing ops from %d to %d on server %d.\n", sm.seqDone, seqEnd, sm.me)

	// if we've already done this op, just return
	if seqEnd <= sm.seqDone {
		return len(sm.configs) - 1
	}

	// jump-start all Paxos sequences between our Done() one and this one
	for seq := sm.seqDone + 1; seq <= seqEnd; seq++ {
		if decided, _ := sm.px.Status(seq); !decided {
			sm.px.Start(seq, Op{Query, sm.uuid(), 0, 0, 0, nil})
		}
	}

	// now do all the ops in order; this loop is blocking on Paxos decision
	for seq := sm.seqDone + 1; seq <= seqEnd; seq++ {
		decided, val := sm.px.Status(seq)
		if !decided {
			sm.waitForPaxos(seq)
			decided, val = sm.px.Status(seq)
		}

		// Execute Join, Leave, and Move operations (Query operations don't modify state)
		switch op := val.(Op); op.Name {
		case Join:
			config := sm.duplicateLast()
			config.Groups[op.GID] = op.Servers
			config = sm.loadBalance(config)
			sm.configs = append(sm.configs, config)
		case Leave:
			config := sm.duplicateLast()
			delete(config.Groups, op.GID)
			config = sm.loadBalance(config)
			sm.configs = append(sm.configs, config)
		case Move:
			config := sm.duplicateLast()
			config.Shards[op.Shard] = op.GID
			sm.configs = append(sm.configs, config)
		}
	}

	sm.seqDone = seqEnd
	sm.px.Done(sm.seqDone)

	return len(sm.configs) - 1
}

// Join handles Join RPC requests from clients.
// It adds a new replica group to the system and redistributes shards to balance load.
func (sm *ShardMaster) Join(args *JoinArgs, reply *JoinReply) error {
	DPrintf("Join %d on server %d.\n", args.GID, sm.me)

	op := Op{Name: Join, OpId: sm.uuid(), GID: args.GID, Shard: 0, Num: 0, Servers: args.Servers}
	seq := sm.decide(op)
	sm.doOps(seq)

	return nil
}

// Leave handles Leave RPC requests from clients.
// It removes a replica group from the system and redistributes its shards to remaining groups.
func (sm *ShardMaster) Leave(args *LeaveArgs, reply *LeaveReply) error {
	DPrintf("Leave %d on server %d.\n", args.GID, sm.me)

	op := Op{Name: Leave, OpId: sm.uuid(), GID: args.GID, Shard: 0, Num: 0, Servers: nil}
	seq := sm.decide(op)
	sm.doOps(seq)

	return nil
}

// Move handles Move RPC requests from clients.
// It moves a specific shard to a specific replica group.
// This is primarily used for testing and fine-tuning load balance.
func (sm *ShardMaster) Move(args *MoveArgs, reply *MoveReply) error {
	DPrintf("Move %d on server %d.\n", args.GID, sm.me)

	op := Op{Name: Move, OpId: sm.uuid(), GID: args.GID, Shard: args.Shard, Num: 0, Servers: nil}
	seq := sm.decide(op)
	sm.doOps(seq)

	return nil
}

// Query handles Query RPC requests from clients.
// It retrieves a configuration by number (-1 for the latest configuration).
// Query(-1) may return fresher data than the Paxos log would suggest, but it won't return stale data.
// To return exact values, we would need a per-op cache, and we don't have a good way to garbage-collect that cache.
func (sm *ShardMaster) Query(args *QueryArgs, reply *QueryReply) error {
	DPrintf("Query %d on server %d.\n", args.Num, sm.me)

	op := Op{Name: Query, OpId: sm.uuid(), GID: 0, Shard: 0, Num: args.Num, Servers: nil}
	seq := sm.decide(op)
	con := sm.doOps(seq)

	DPrintf("Waiting for lock (to query) on server %d.", sm.me)

	sm.mu.Lock()
	if 0 <= args.Num && args.Num < len(sm.configs) {
		reply.Config = sm.configs[args.Num]
	} else {
		reply.Config = sm.configs[con]
	}
	sm.mu.Unlock()

	return nil
}

// Kill tells the server to shut itself down.
// This method is used for testing and graceful shutdown.
func (sm *ShardMaster) Kill() {
	sm.dead = true
	sm.l.Close()
	sm.px.Kill()
}

// StartServer creates and starts a new ShardMaster server.
// servers contains the ports of all servers that will cooperate via Paxos.
// me is the index of the current server in the servers array.
func StartServer(servers []string, me int) *ShardMaster {
	// Register Op struct for RPC marshalling/unmarshalling
	gob.Register(Op{})

	sm := &ShardMaster{
		me:         me,
		seqTried:   0,
		seqDone:    -1,
		configs:    make([]Config, 1),
	}

	// Initialize the first configuration (configuration 0)
	sm.configs[0].Groups = make(map[int64][]string)

	rpcs := rpc.NewServer()
	rpcs.Register(sm)

	sm.px = paxos.Make(servers, me, rpcs, false, "", false)

	// Set up Unix socket listener
	os.Remove(servers[me])
	l, err := net.Listen("unix", servers[me])
	if err != nil {
		log.Fatal("listen error: ", err)
	}
	sm.l = l

	// Start RPC server goroutine to handle incoming connections
	go func() {
		for !sm.dead {
			conn, err := sm.l.Accept()
			if err == nil && !sm.dead {
				if sm.unreliable && (rand.Int63()%1000) < 100 {
					// Discard the request (simulate network failure)
					conn.Close()
				} else if sm.unreliable && (rand.Int63()%1000) < 200 {
					// Process the request but force discard of reply
					c1 := conn.(*net.UnixConn)
					f, _ := c1.File()
					err := syscall.Shutdown(int(f.Fd()), syscall.SHUT_WR)
					if err != nil {
						fmt.Printf("shutdown: %v\n", err)
					}
					go rpcs.ServeConn(conn)
				} else {
					// Normal processing
					go rpcs.ServeConn(conn)
				}
			} else if err == nil {
				conn.Close()
			}
			if err != nil && !sm.dead {
				fmt.Printf("ShardMaster(%v) accept: %v\n", me, err.Error())
				sm.Kill()
			}
		}
	}()

	return sm
}
