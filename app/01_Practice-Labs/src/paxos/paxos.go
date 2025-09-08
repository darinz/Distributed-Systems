// Package paxos implements the Paxos consensus algorithm for distributed systems.
// This library provides a fault-tolerant consensus mechanism that allows a set of
// distributed processes to agree on a sequence of values even in the presence of
// network failures, message loss, and process crashes.
//
// The Paxos algorithm ensures that:
// - A majority of peers can reach agreement on a value
// - Once a value is decided, all peers will eventually agree on that value
// - The system can tolerate minority failures and network partitions
// - Multiple instances can be agreed upon concurrently
//
// Key Features:
// - Manages a sequence of agreed-upon values across multiple instances
// - Handles network failures, partitions, and message loss
// - Supports concurrent agreement on multiple instances
// - Implements memory management through instance forgetting
// - Provides both in-memory and persistent storage options
//
// The application interface:
//   px = paxos.Make(peers []string, me int, rpcs *rpc.Server, saveToDisk bool, dir string, restart bool)
//   px.Start(seq int, v interface{}) -- start agreement on new instance
//   px.Status(seq int) (decided bool, v interface{}) -- get info about an instance
//   px.Done(seq int) -- ok to forget all instances <= seq
//   px.Max() int -- highest instance seq known, or -1
//   px.Min() int -- instances before this seq have been forgotten
package paxos

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// GCInterval defines the interval at which the garbage collector runs
// to clean up forgotten Paxos instances and free memory.
const GCInterval = 500 * time.Millisecond

// Debug controls whether debug output is printed
const Debug = false

// DPrintf prints debug output if Debug is enabled
func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

// Paxos represents a Paxos peer that participates in consensus decisions.
// Each peer maintains its own state and communicates with other peers
// to reach agreement on a sequence of values.
type Paxos struct {
	mu         sync.Mutex        // Mutex for protecting shared state
	l          net.Listener      // Network listener for RPC connections
	dead       bool              // Flag indicating if peer should shut down
	unreliable bool              // Flag for unreliable network simulation (testing)
	rpcCount   int               // Counter for RPC calls (testing)
	peers      []string          // List of all peer addresses
	me         int               // Index of this peer in the peers array

	// Paxos state
	instances map[int]PaxosInstance // Map of sequence number to instance state
	done      map[int]int           // Map of peer index to highest Done() value
	nseq      int                   // Highest sequence number seen
	nmajority int                   // Number of peers required for majority

	// Persistence
	dir        string // Directory for persistent storage
	saveToDisk bool   // Whether to save state to disk
}

// PaxosInstance represents the state of a single Paxos agreement instance.
// Each instance tracks the proposal numbers and values for the three phases
// of the Paxos algorithm: prepare, accept, and decide.
type PaxosInstance struct {
	N_p     int         // Highest prepare request number seen
	N_a     int         // Highest accept request number seen
	V_a     interface{} // Value from the highest accept request
	Decided bool        // Whether this instance has been decided
}

// PrepareArgs contains arguments for the prepare phase of Paxos
type PrepareArgs struct {
	Seq int // Sequence number of the instance
	N   int // Proposal number
}

// PrepareReply contains the response to a prepare request
type PrepareReply struct {
	N      int         // Highest accept number seen
	V      interface{} // Value from the highest accept
	Reject bool        // Whether the prepare was rejected
	Done   int         // Highest Done() value from this peer
}

// AcceptArgs contains arguments for the accept phase of Paxos
type AcceptArgs struct {
	Seq int         // Sequence number of the instance
	N   int         // Proposal number
	V   interface{} // Proposed value
}

// AcceptReply contains the response to an accept request
type AcceptReply struct {
	Reject bool // Whether the accept was rejected
	Done   int  // Highest Done() value from this peer
}

// DecideArgs contains arguments for the decide phase of Paxos
type DecideArgs struct {
	Seq int         // Sequence number of the instance
	N   int         // Proposal number
	V   interface{} // Decided value
}

// DecideReply contains the response to a decide request
type DecideReply struct {
	Reject bool // Whether the decide was rejected
	Done   int  // Highest Done() value from this peer
}

// call sends an RPC to the specified server and waits for a reply.
// It returns true if the server responded successfully, false otherwise.
// The reply argument should be a pointer to a reply structure.
// This function handles connection establishment, RPC call, and cleanup.
func call(srv string, name string, args interface{}, reply interface{}) bool {
	c, err := rpc.Dial("unix", srv)
	if err != nil {
		err1 := err.(*net.OpError)
		if err1.Err != syscall.ENOENT && err1.Err != syscall.ECONNREFUSED {
			fmt.Printf("paxos Dial() failed: %v\n", err1)
		}
		return false
	}
	defer c.Close()

	err = c.Call(name, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}

// updatePaxos updates the Paxos instance state, optionally persisting to disk
func (px *Paxos) updatePaxos(seq int, instance PaxosInstance) {
	if px.saveToDisk {
		px.fileUpdatePaxos(seq, instance)
	}
	px.instances[seq] = instance
}

// fileUpdatePaxos saves a Paxos instance to disk using atomic file operations
func (px *Paxos) fileUpdatePaxos(seq int, op PaxosInstance) error {
	DPrintf("Saving paxos sequence %d: %v on server %d", seq, op, px.me)
	fullname := px.dir + "/paxos-" + strconv.Itoa(seq)
	tempname := px.dir + "/temp-" + strconv.Itoa(seq)

	// Encode the instance to bytes
	w := new(bytes.Buffer)
	e := gob.NewEncoder(w)
	if err := e.Encode(op); err != nil {
		return err
	}

	// Write to temporary file first, then rename atomically
	if err := ioutil.WriteFile(tempname, w.Bytes(), 0666); err != nil {
		return err
	}
	if err := os.Rename(tempname, fullname); err != nil {
		return err
	}
	return nil
}

// fileRetrievePaxos loads all Paxos instances from disk
func (px *Paxos) fileRetrievePaxos() map[int]PaxosInstance {
	m := map[int]PaxosInstance{}
	d := px.dir
	files, err := ioutil.ReadDir(d)
	if err != nil {
		log.Fatalf("fileRetrievePaxos could not read %v: %v", d, err)
	}
	
	for _, fi := range files {
		n1 := fi.Name()
		if len(n1) >= 6 && n1[0:6] == "paxos-" {
			key, err := strconv.Atoi(n1[6:])
			if err != nil {
				log.Fatalf("fileRetrievePaxos bad file name %v: %v", n1, err)
			}
			fullname := px.dir + "/" + n1
			content, err := ioutil.ReadFile(fullname)
			if err != nil {
				log.Fatalf("fileRetrievePaxos fileGet failed for %v: %v", key, err)
			}
			buf := bytes.NewBuffer(content)
			g := gob.NewDecoder(buf)
			var inst PaxosInstance
			if err := g.Decode(&inst); err != nil {
				log.Fatalf("fileRetrievePaxos decode failed for %v: %v", key, err)
			}
			DPrintf("Retrieved paxos sequence %d: %v on server %d", key, inst, px.me)
			m[key] = inst
		}
	}
	return m
}

// propose implements the proposer side of the Paxos algorithm.
// It attempts to get a majority of peers to agree on the proposed value
// for the given sequence number.
func (px *Paxos) propose(seq int, value interface{}) {
	// Initialize proposal number to ensure uniqueness across peers
	// Each peer uses a different starting point based on its index
	initProposalNum := (px.me + seq) % len(px.peers)
	n := initProposalNum

	// Continue until agreement is reached or peer is killed
	for !px.dead {
		// Check if this instance has already been decided
		if inst, ok := px.instances[seq]; ok && inst.Decided {
			return
		}

		// Phase 1: Prepare
		nseen := -1
		nprepareOK := 0
		servPrepareOK := make([]int, 0)
		
		// Send prepare(n) to all servers including self
		for peer := 0; peer < len(px.peers); peer++ {
			args := PrepareArgs{Seq: seq, N: n}
			var reply PrepareReply

			if px.send(peer, "Paxos.Prepare", args, &reply) && !reply.Reject {
				nprepareOK++
				servPrepareOK = append(servPrepareOK, peer)
				// Choose value with highest proposal number seen
				if reply.N > nseen {
					nseen = reply.N
					value = reply.V
				}
			}
		}

		// If we didn't get majority approval, try with higher proposal number
		if nprepareOK < px.nmajority {
			if n < nseen {
				n = (nseen/len(px.peers))*len(px.peers) + initProposalNum
			}
			n += len(px.peers)
			continue
		}

		// Phase 2: Accept
		nacceptOK := 0
		servAcceptOK := make([]int, 0)
		for peer := 0; peer < len(px.peers); peer++ {
			args := AcceptArgs{Seq: seq, N: n, V: value}
			var reply AcceptReply
			if px.send(peer, "Paxos.Accept", args, &reply) && !reply.Reject {
				nacceptOK++
				servAcceptOK = append(servAcceptOK, peer)
			}
		}

		// If we didn't get majority approval, try with higher proposal number
		if nacceptOK < px.nmajority {
			if n < nseen {
				n = (nseen/len(px.peers))*len(px.peers) + initProposalNum
			}
			n += len(px.peers)
			continue
		}

		// Phase 3: Decide - send decided value to all peers
		for peer := 0; peer < len(px.peers); peer++ {
			args := DecideArgs{Seq: seq, N: n, V: value}
			var reply DecideReply
			if px.send(peer, "Paxos.Decided", args, &reply) && !reply.Reject {
				px.done[peer] = reply.Done
			}
		}

		// Agreement reached, exit the loop
		break
	}
}

// send makes either a local procedure call (LPC) to self or an RPC to other peers.
// This abstraction allows the proposer to treat local and remote calls uniformly.
func (px *Paxos) send(peer int, pc string, args interface{}, reply interface{}) bool {
	// Make RPC call to remote peer
	if peer != px.me {
		return call(px.peers[peer], pc, args, reply)
	}

	// Make local procedure call to self
	pc = strings.TrimPrefix(pc, "Paxos.") // Remove "Paxos." prefix for local calls
	switch pc {
	case "Prepare":
		px.Prepare(args.(PrepareArgs), reply.(*PrepareReply))
	case "Accept":
		px.Accept(args.(AcceptArgs), reply.(*AcceptReply))
	case "Decided":
		px.Decided(args.(DecideArgs), reply.(*DecideReply))
	default:
		return false
	}

	return true
}

// Prepare handles the prepare phase of the Paxos algorithm.
// It responds to prepare requests from proposers, promising not to accept
// proposals with lower numbers and returning the highest accepted value.
func (px *Paxos) Prepare(args PrepareArgs, reply *PrepareReply) error {
	px.mu.Lock()
	defer px.mu.Unlock()

	// Get or create instance for this sequence number
	var pi PaxosInstance
	if pii, ok := px.instances[args.Seq]; ok {
		pi = pii
	} else {
		pi = PaxosInstance{N_p: -1, N_a: -1, V_a: nil, Decided: false}
	}

	// Acceptor's prepare(n) handler:
	// If n > n_p, promise not to accept proposals with lower numbers
	if args.N > pi.N_p {
		// Update n_p to the new proposal number
		instance := PaxosInstance{
			N_p:     args.N,
			N_a:     pi.N_a,
			V_a:     pi.V_a,
			Decided: pi.Decided,
		}
		px.updatePaxos(args.Seq, instance)

		// Reply with prepare_ok(n_a, v_a)
		reply.N = pi.N_a
		reply.V = pi.V_a
		reply.Reject = false
	} else {
		// Reject prepare request with lower proposal number
		reply.Reject = true
	}

	// Piggyback Done value
	reply.Done = px.done[px.me]
	return nil
}

// Accept handles the accept phase of the Paxos algorithm.
// It responds to accept requests from proposers, accepting values
// with proposal numbers greater than or equal to the highest prepare seen.
func (px *Paxos) Accept(args AcceptArgs, reply *AcceptReply) error {
	px.mu.Lock()
	defer px.mu.Unlock()

	// Get instance for this sequence number
	pi := px.instances[args.Seq]

	// Acceptor's accept(n, v) handler:
	// If n >= n_p, accept the proposal
	if args.N >= pi.N_p {
		// Update n_p, n_a, and v_a
		instance := PaxosInstance{
			N_p:     args.N,
			N_a:     args.N,
			V_a:     args.V,
			Decided: pi.Decided,
		}
		px.updatePaxos(args.Seq, instance)
		reply.Reject = false
	} else {
		// Reject accept request with lower proposal number
		reply.Reject = true
	}

	// Piggyback Done value
	reply.Done = px.done[px.me]
	return nil
}

// Decided handles the decide phase of the Paxos algorithm.
// It marks the instance as decided with the agreed-upon value.
func (px *Paxos) Decided(args DecideArgs, reply *DecideReply) error {
	px.mu.Lock()
	defer px.mu.Unlock()

	// Mark the instance as decided with the agreed value
	instance := PaxosInstance{
		N_p:     args.N,
		N_a:     args.N,
		V_a:     args.V,
		Decided: true,
	}
	px.updatePaxos(args.Seq, instance)

	// Piggyback the Done value
	reply.Done = px.done[px.me]
	reply.Reject = false
	return nil
}

// Start begins agreement on a new Paxos instance with the given sequence number and value.
// It returns immediately without waiting for agreement to complete.
// The application should call Status() to check if/when agreement is reached.
func (px *Paxos) Start(seq int, v interface{}) {
	// Start the proposer in a separate goroutine
	go px.propose(seq, v)
}

// Done indicates that the application on this machine is done with all instances <= seq.
// This allows Paxos to forget information about old instances to free memory.
// See the comments for Min() for more explanation.
func (px *Paxos) Done(seq int) {
	px.mu.Lock()
	defer px.mu.Unlock()
	px.done[px.me] = seq
}

// Max returns the highest instance sequence number known to this peer.
// Returns -1 if no instances have been seen yet.
func (px *Paxos) Max() int {
	px.mu.Lock()
	defer px.mu.Unlock()
	return px.nseq
}

// Min returns one more than the minimum among z_i, where z_i is the highest number
// ever passed to Done() on peer i. A peer's z_i is -1 if it has never called Done().
//
// Paxos is required to have forgotten all information about any instances it knows
// that are < Min(). The point is to free up memory in long-running Paxos-based servers.
//
// Paxos peers need to exchange their highest Done() arguments in order to implement Min().
// These exchanges can be piggybacked on ordinary Paxos agreement protocol messages,
// so it is OK if one peer's Min does not reflect another peer's Done() until after
// the next instance is agreed to.
//
// The fact that Min() is defined as a minimum over *all* Paxos peers means that Min()
// cannot increase until all peers have been heard from. So if a peer is dead or
// unreachable, other peers' Min()s will not increase even if all reachable peers call Done.
// The reason for this is that when the unreachable peer comes back to life, it will need
// to catch up on instances that it missed -- the other peers therefore cannot forget these instances.
func (px *Paxos) Min() int {
	px.mu.Lock()
	defer px.mu.Unlock()

	min := math.MaxInt64
	for _, seq := range px.done {
		if seq < min {
			min = seq
		}
	}
	return min + 1
}

// Status checks whether this peer thinks an instance has been decided,
// and if so, what the agreed value is. Status() should just inspect the local
// peer state; it should not contact other Paxos peers.
func (px *Paxos) Status(seq int) (bool, interface{}) {
	px.mu.Lock()
	defer px.mu.Unlock()

	// Check if instance exists and is decided
	if pi, ok := px.instances[seq]; ok {
		return pi.Decided, pi.V_a
	}

	// Instance for the sequence number does not exist yet
	if seq > px.nseq {
		px.nseq = seq
	}

	return false, -1
}

// Kill tells the peer to shut itself down.
// This method is used for testing and graceful shutdown.
func (px *Paxos) Kill() {
	px.dead = true
	if px.l != nil {
		px.l.Close()
	}
}

// Make creates a new Paxos peer that participates in consensus decisions.
// peers contains the addresses of all Paxos peers (including this one).
// me is the index of this peer in the peers array.
// rpcs is the RPC server to register with (nil to create a new one).
// saveToDisk determines whether to persist state to disk.
// dir is the directory for persistent storage.
// restart indicates whether this is a restart from a crash.
func Make(peers []string, me int, rpcs *rpc.Server, saveToDisk bool, dir string, restart bool) *Paxos {
	px := &Paxos{
		peers:      peers,
		me:         me,
		dir:        dir,
		saveToDisk: saveToDisk,
	}

	// Register PaxosInstance for gob encoding
	gob.Register(PaxosInstance{})

	// Initialize Paxos state
	px.instances = make(map[int]PaxosInstance)
	px.done = make(map[int]int)
	for peer := 0; peer < len(px.peers); peer++ {
		// A peer's z_i is -1 if it has never called Done()
		px.done[peer] = -1
	}
	px.nseq = -1
	px.nmajority = len(px.peers)/2 + 1

	// Load state from disk if restarting
	if saveToDisk && restart {
		px.instances = px.fileRetrievePaxos()
	}

	// Set up RPC server
	if rpcs != nil {
		// Caller will create socket and handle connections
		rpcs.Register(px)
	} else {
		rpcs = rpc.NewServer()
		rpcs.Register(px)

		// Prepare to receive connections from clients
		// Change "unix" to "tcp" to use over a network
		os.Remove(peers[me]) // Only needed for "unix"
		l, err := net.Listen("unix", peers[me])
		if err != nil {
			log.Fatal("listen error: ", err)
		}
		px.l = l

		// Create a thread to accept RPC connections
		go func() {
			for !px.dead {
				conn, err := px.l.Accept()
				if err == nil && !px.dead {
					if px.unreliable && (rand.Int63()%1000) < 100 {
						// Discard the request (simulate network failure)
						conn.Close()
					} else if px.unreliable && (rand.Int63()%1000) < 200 {
						// Process the request but force discard of reply
						c1 := conn.(*net.UnixConn)
						f, _ := c1.File()
						err := syscall.Shutdown(int(f.Fd()), syscall.SHUT_WR)
						if err != nil {
							fmt.Printf("shutdown: %v\n", err)
						}
						px.rpcCount++
						go rpcs.ServeConn(conn)
					} else {
						// Normal processing
						px.rpcCount++
						go rpcs.ServeConn(conn)
					}
				} else if err == nil {
					conn.Close()
				}
				if err != nil && !px.dead {
					fmt.Printf("Paxos(%v) accept: %v\n", me, err.Error())
				}
			}
		}()
	}

	// Start garbage collector to clean up forgotten instances
	go func() {
		for !px.dead {
			px.mu.Lock()
			min := px.Min()
			for ninst := range px.instances {
				if ninst < min {
					delete(px.instances, ninst)
					if px.saveToDisk {
						fullname := px.dir + "/paxos-" + strconv.Itoa(ninst)
						os.Remove(fullname)
					}
				}
			}
			px.mu.Unlock()
			time.Sleep(GCInterval)
		}
	}()

	return px
}
