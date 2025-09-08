package kvpaxos

import (
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	"syscall"
	"time"

	"distributed-systems/app/01_Practice-Labs/src/paxos"
)

// Debug controls whether debug output is printed
const Debug = 0

// DPrintf prints debug output if Debug is enabled
func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug > 0 {
		log.Printf(format, a...)
	}
	return
}

// Nop represents a no-operation used to fill holes in the Paxos log
var Nop = Op{Client: 0, OpId: 0, Put: 0, Key: "", Value: ""}

// Op represents an operation that can be agreed upon through Paxos.
// It contains all the information needed to execute a key-value operation.
type Op struct {
	Client uint64 // Client identifier for duplicate detection
	OpId   uint64 // Unique operation ID for duplicate detection
	Put    int    // Operation type: 0 = Get, 1 = Put, 2 = PutHash
	Key    string // The key for the operation
	Value  string // The value for Put operations (empty for Gets)
}

// OpResult represents the result of executing an operation.
// It stores the operation, its result, and the Paxos sequence number.
type OpResult struct {
	Op     Op     // The operation that was executed
	Result string // Result value (for Get) or previous value (for PutHash)
	Seq    int    // Paxos sequence number where this operation was agreed upon
}

// KVPaxos represents a Paxos-based key-value server.
// It maintains a replicated key-value store using Paxos consensus
// to ensure all replicas stay synchronized.
type KVPaxos struct {
	mu         sync.Mutex        // Mutex for protecting shared state
	l          net.Listener      // Network listener for RPC connections
	me         int               // Server identifier
	dead       bool              // Flag indicating if server should shut down (for testing)
	unreliable bool              // Flag for unreliable network simulation (for testing)
	px         *paxos.Paxos      // Paxos peer for consensus

	// Key-value store state
	store    map[string]string           // Replicated key-value store
	opLog    map[uint64]map[int]OpResult // Operation log per client per sequence number
	idLog    map[uint64]bool             // Set of completed operation IDs for duplicate detection
	seqTried int                         // Hint for next sequence number to try
	seqDone  int                         // Last sequence number for which Done() was called
}

// waitForPaxos waits for a Paxos instance to be decided.
// It uses exponential backoff to avoid busy waiting.
// Returns true if the instance was decided, false if it was already done.
func (kv *KVPaxos) waitForPaxos(seq int) bool {
	to := 10 * time.Millisecond
	for {
		decided, _ := kv.px.Status(seq)
		if decided {
			return true
		}
		time.Sleep(to)
		if to < 10*time.Second {
			to *= 2
		} else {
			DPrintf("Waiting for Paxos seq(%d) on server %d", seq, kv.me)
		}
		if seq < kv.seqDone {
			DPrintf("Seq(%d) is already done, through %d. Stop waiting on server %d.", seq, kv.seqDone, kv.me)
			return false
		}
	}
}

// doPut executes a Put or PutHash operation on the key-value store.
// This is a non-locking internal helper function.
// For PutHash operations, it computes the hash of (previous value + new value).
// Returns the previous value (for PutHash) or empty string (for Put).
func (kv *KVPaxos) doPut(op Op, seq int, forceLog bool) string {
	if op.Put > 1 { // PutHash operation
		prevValue, ok := kv.store[op.Key]
		if !ok {
			prevValue = ""
		}
		// Compute hash of previous value + new value
		value := strconv.Itoa(int(hash(prevValue + op.Value)))
		kv.store[op.Key] = value
		kv.logResult(OpResult{Op: op, Result: prevValue, Seq: seq})
		return prevValue
	} else { // Regular Put operation
		kv.store[op.Key] = op.Value
		kv.logResult(OpResult{Op: op, Result: "", Seq: seq})
		return ""
	}
}

// doGet executes all operations up to and including the given sequence number.
// This is a non-locking internal function. Callers are responsible for locking.
// It ensures the server is caught up with the Paxos log and executes the requested operation.
// Returns the result of the operation (value for Get, previous value for PutHash).
func (kv *KVPaxos) doGet(req Op, seqEnd int) string {
	if req != Nop {
		DPrintf("DoGet(%d) with Done = %d on server %d", seqEnd, kv.seqDone, kv.me)
	}

	// If we've already processed this sequence, look up the result in the log
	if seqEnd <= kv.seqDone {
		record, ok := kv.checkLog(req)
		if ok {
			return record.Result
		}
		return ""
	}

	// Jump-start decisions for sequences we need to catch up on
	// This is essential to avoid deadlock when servers fall behind
	for seq := kv.seqDone + 1; seq < seqEnd; seq++ {
		decided, _ := kv.px.Status(seq)
		if !decided {
			// Submit a no-op to force agreement on this sequence
			kv.px.Start(seq, Nop)
		}
	}

	// Execute all operations from seqDone+1 to seqEnd-1
	for seq := kv.seqDone + 1; seq < seqEnd; seq++ {
		decided, val := kv.px.Status(seq)
		if !decided {
			kv.waitForPaxos(seq)
			decided, val = kv.px.Status(seq)
		}
		op := val.(Op)
		var res string
		if op.Put > 0 {
			res = kv.doPut(op, seq, false)
		} else {
			res = kv.store[op.Key]
			kv.logResult(OpResult{Op: op, Result: res, Seq: seq})
		}
		DPrintf("Server %d has caught up with Seq(%d): Get/Put(%s) ID = %d: %s", kv.me, seq, op.Key, op.OpId, res)
	}

	// Execute the requested operation at seqEnd
	decided, val := kv.px.Status(seqEnd)
	for !decided {
		kv.px.Start(seqEnd, Nop)
		kv.waitForPaxos(seqEnd)
		decided, val = kv.px.Status(seqEnd)
	}

	op := val.(Op)
	var value string
	if op.Put > 0 {
		value = kv.doPut(op, seqEnd, true)
	} else {
		if v, ok := kv.store[op.Key]; ok {
			value = v
		} else {
			value = ""
		}
		kv.logResult(OpResult{Op: op, Result: value, Seq: seqEnd})
	}

	// Tell Paxos we are finished with this operation and all previous ones
	kv.px.Done(seqEnd)
	kv.seqDone = seqEnd
	return value
}

// logResult logs the result of a Put or PutHash operation for duplicate detection.
// Get operations are not logged as they don't need duplicate detection.
func (kv *KVPaxos) logResult(op OpResult) {
	if op.Op.Put == 0 {
		return // Don't need to log Get operations
	}
	if kv.opLog[op.Op.Client] == nil {
		kv.opLog[op.Op.Client] = make(map[int]OpResult)
	}
	kv.opLog[op.Op.Client][op.Seq] = op
}

// checkLog checks if an operation has already been executed by looking in the operation log.
// Returns the operation result and true if found, false otherwise.
func (kv *KVPaxos) checkLog(op Op) (OpResult, bool) {
	records, ok := kv.opLog[op.Client]
	if !ok {
		return OpResult{Op: Nop, Result: "", Seq: -1}, false
	}
	for _, record := range records {
		if record.Op.OpId == op.OpId {
			return record, true
		}
	}
	return OpResult{Op: Nop, Result: "", Seq: -1}, false
}

// clearLog clears old operation log entries to prevent memory growth.
// It keeps operations at and ahead of the given sequence, plus 8 operations behind it.
func (kv *KVPaxos) clearLog(client uint64, seq int) {
	records, ok := kv.opLog[client]
	if !ok {
		return
	}
	prev := make(map[int]bool)
	prevMin := -1
	for _, record := range records {
		if record.Seq < prevMin {
			delete(records, record.Seq)
		} else if record.Seq < seq {
			prev[record.Seq] = true
			if len(prev) > 6 {
				delete(records, prevMin)
				prevMin = seq
				for key := range prev {
					if prevMin < key {
						prevMin = key
					}
				}
				delete(prev, prevMin)
			}
		}
	}
	kv.opLog[client] = records
}

// decideSeq attempts to assign a Paxos sequence number to the given operation.
// It tries different sequence numbers until it successfully gets agreement on the operation.
// Returns the sequence number where the operation was agreed upon.
func (kv *KVPaxos) decideSeq(op Op) int {
	for {
		seq := kv.seqTried

		// Check if this sequence is already decided
		if decided, _ := kv.px.Status(seq); decided {
			kv.seqTried++
			continue
		}

		// Check if the operation has already been processed
		kv.doGet(Nop, seq-1)
		record, ok := kv.checkLog(op)
		if ok {
			return record.Seq
		}

		DPrintf("OpID=%d not found in log on server %d. Attempting Paxos seq=%d", op.OpId, kv.me, seq)

		kv.seqTried++

		// Try to get agreement on this sequence with our operation
		kv.px.Start(seq, op)
		kv.waitForPaxos(seq)
		decided, val := kv.px.Status(seq)
		if decided && val != nil && (op.OpId == val.(Op).OpId) {
			return seq
		}
	}
}

// Get handles Get RPC requests from clients.
// It ensures the operation is agreed upon through Paxos and returns the current value.
func (kv *KVPaxos) Get(args *GetArgs, reply *GetReply) error {
	DPrintf("Server Get(%s), ID=%d, to server %d\n", args.Key, args.OpId, kv.me)
	kv.mu.Lock()
	defer kv.mu.Unlock()

	op := Op{Client: args.Client, OpId: args.OpId, Put: 0, Key: args.Key, Value: ""}

	// First, check if this operation has already been processed
	kv.doGet(Nop, kv.seqTried-1)

	record, ok := kv.checkLog(op)
	if ok && record.Op.OpId == args.OpId {
		DPrintf("Server Get(%s) to server %d found in log\n", args.Key, kv.me)
		reply.Value = record.Result
		reply.Err = ""
		return nil
	}

	// Try to get agreement on a Paxos sequence for this operation
	seq := kv.decideSeq(op)

	DPrintf("Server Get(%s) on server %d decided for Seq(%d)", args.Key, kv.me, seq)

	// Execute the operation and get the result
	reply.Value = kv.doGet(op, seq)
	reply.Err = ""

	// Clean up old log entries for this client
	kv.clearLog(args.Client, seq)

	DPrintf("Server Get(%s), ID = %d, on server %d, Seq(%d) returns value %s", args.Key, args.OpId, kv.me, seq, reply.Value)

	return nil
}

// Put handles Put and PutHash RPC requests from clients.
// It ensures the operation is agreed upon through Paxos and executes it atomically.
func (kv *KVPaxos) Put(args *PutArgs, reply *PutReply) error {
	DPrintf("Server Put(%s), ID=%d, to server %d\n", args.Key, args.OpId, kv.me)
	kv.mu.Lock()
	defer kv.mu.Unlock()

	// Determine operation type: 1 = Put, 2 = PutHash
	put := 1
	if args.DoHash {
		put = 2
	}
	op := Op{Client: args.Client, OpId: args.OpId, Put: put, Key: args.Key, Value: args.Value}

	// First, check if this operation has already been processed
	kv.doGet(Nop, kv.seqTried-1)
	record, ok := kv.checkLog(op)
	if ok && record.Op.OpId == args.OpId {
		DPrintf("Server Put(%s) to server %d found in log.\n", args.Key, kv.me)
		reply.PreviousValue = record.Result
		reply.Err = ""
		return nil
	}

	// Check if we've seen this operation ID before (old request)
	if kv.idLog[args.OpId] {
		reply.PreviousValue = ""
		return nil
	}

	// Try to get agreement on a Paxos sequence for this operation
	seq := kv.decideSeq(op)

	DPrintf("Server Put(%s) on server %d decided on Seq(%d)", args.Key, kv.me, seq)

	// Execute the operation and get the result
	if args.DoHash {
		reply.PreviousValue = kv.doGet(op, seq)
	} else {
		reply.PreviousValue = ""
	}

	reply.Err = ""
	kv.clearLog(args.Client, seq)
	kv.idLog[args.OpId] = true

	DPrintf("Server Put(%s, %s) ID=%d on server %d, Seq(%d) returns value %s", args.Key, args.Value, args.OpId, kv.me, seq, reply.PreviousValue)
	return nil
}

// Kill tells the server to shut itself down.
// This method is used for testing and graceful shutdown.
func (kv *KVPaxos) Kill() {
	DPrintf("Kill(%d): die\n", kv.me)
	kv.dead = true
	kv.l.Close()
	kv.px.Kill()
}

// StartServer creates and starts a new KVPaxos server.
// servers contains the ports of all servers that will cooperate via Paxos.
// me is the index of the current server in the servers array.
func StartServer(servers []string, me int) *KVPaxos {
	// Register Op struct for RPC marshalling/unmarshalling
	gob.Register(Op{})

	kv := &KVPaxos{
		me:         me,
		seqTried:   0,
		seqDone:    -1,
		store:      make(map[string]string),
		opLog:      make(map[uint64]map[int]OpResult),
		idLog:      make(map[uint64]bool),
	}

	rpcs := rpc.NewServer()
	rpcs.Register(kv)

	kv.px = paxos.Make(servers, me, rpcs, false, "", false)

	// Set up Unix socket listener
	os.Remove(servers[me])
	l, err := net.Listen("unix", servers[me])
	if err != nil {
		log.Fatal("listen error: ", err)
	}
	kv.l = l

	// Start RPC server goroutine to handle incoming connections
	go func() {
		for !kv.dead {
			conn, err := kv.l.Accept()
			if err == nil && !kv.dead {
				if kv.unreliable && (rand.Int63()%1000) < 100 {
					// Discard the request (simulate network failure)
					conn.Close()
				} else if kv.unreliable && (rand.Int63()%1000) < 200 {
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
			if err != nil && !kv.dead {
				fmt.Printf("KVPaxos(%v) accept: %v\n", me, err.Error())
				DPrintf("RPC handler error: Killing server %d.\n", me)
				kv.Kill()
			}
		}
	}()

	return kv
}
