package pbservice

import (
	"distributed-systems/app/01_Practice-Labs/src/viewservice"
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
)

// Debug controls whether debug output is printed
const Debug = false

// DPrintf prints debug output if Debug is enabled
func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		n, err = fmt.Printf(format, a...)
	}
	return
}

// PBServer represents a primary/backup key-value server.
// It can act as either a primary (handling client requests) or backup
// (replicating primary state) depending on the current view.
type PBServer struct {
	l          net.Listener        // Network listener for RPC connections
	dead       bool                // Flag indicating if server should shut down (for testing)
	unreliable bool                // Flag for unreliable network simulation (for testing)
	me         string              // Server's network address/identifier
	vs         *viewservice.Clerk  // View service client for discovering current view
	done       sync.WaitGroup      // WaitGroup for graceful shutdown
	finish     chan interface{}    // Channel to signal completion

	// Server state
	vshost string               // Address of the view service
	view   viewservice.View     // Current view (primary/backup assignments)
	store  map[string]string    // Key-value store
	gets   map[uint64]GetEntry  // Log of Get operations for duplicate detection
	puts   map[uint64]PutEntry  // Log of Put operations for duplicate detection
	mu     sync.Mutex           // Mutex for protecting shared state
}

// Put handles Put and PutHash operations from clients.
// Only the primary server can execute Put operations.
// It ensures at-most-once semantics by checking for duplicate operation IDs.
// For PutHash operations, it computes a hash of the previous value + new value.
// It forwards the operation to the backup server before committing locally.
func (pb *PBServer) Put(args *PutArgs, reply *PutReply) error {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	DPrintf("Server Put(%s, %s)\n", args.Key, args.Value)

	// Only the primary can execute Put operations
	if pb.me != pb.view.Primary {
		reply.PreviousValue = ""
		reply.Err = ErrWrongServer
		return nil
	}

	// Check for duplicate operation ID from the same client
	if entry, exists := pb.puts[args.Id]; exists {
		if args.Client == entry.Client {
			reply.Err = entry.Reply.Err
			reply.PreviousValue = entry.Reply.PreviousValue
			return nil
		}
	}

	// Handle PutHash operation
	if args.DoHash {
		// Get previous value or use empty string if key doesn't exist
		if value, exists := pb.store[args.Key]; exists {
			reply.PreviousValue = value
		} else {
			reply.PreviousValue = ""
		}
		// Compute hash of previous value + new value
		args.Value = strconv.Itoa(int(hash(reply.PreviousValue + args.Value)))
	}
	reply.Err = OK

	// Forward the operation to backup server if one exists
	if pb.view.Backup != "" {
		fargs := ForwardPutArgs{Args: args, Reply: reply}
		var freply ForwardReply

		ok := call(pb.view.Backup, "PBServer.ForwardPut", &fargs, &freply)
		if !ok || freply.Err != OK {
			reply.PreviousValue = ""
			reply.Err = ErrBackup
			return nil
		}
	}

	// Commit the operation locally and log it
	pb.puts[args.Id] = PutEntry{Reply: *reply, Client: args.Client}
	pb.store[args.Key] = args.Value

	return nil
}

// ForwardPut handles forwarded Put operations from the primary server.
// Only the backup server can accept forwarded operations.
// It replicates the operation on the backup's local store and logs it.
func (pb *PBServer) ForwardPut(args *ForwardPutArgs, reply *ForwardReply) error {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	DPrintf("Server ForwardPut(%s, %s)\n", args.Args.Key, args.Args.Value)

	// Only the backup server can accept forwarded operations
	if pb.me != pb.view.Backup {
		reply.Err = ErrWrongServer
		return nil
	}

	// Replicate the operation on the backup's store and log it
	pb.puts[args.Args.Id] = PutEntry{Reply: *args.Reply, Client: args.Args.Client}
	pb.store[args.Args.Key] = args.Args.Value

	reply.Err = OK
	return nil
}


// Get handles Get operations from clients.
// Only the primary server can execute Get operations.
// It ensures at-most-once semantics by checking for duplicate operation IDs.
// It forwards the operation to the backup server to maintain consistency.
func (pb *PBServer) Get(args *GetArgs, reply *GetReply) error {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	DPrintf("Server Get(%s)\n", args.Key)

	// Only the primary can execute Get operations
	if pb.me != pb.view.Primary {
		reply.Value = ""
		reply.Err = ErrWrongServer
		return nil
	}

	// Check for duplicate operation ID from the same client
	if entry, exists := pb.gets[args.Id]; exists {
		if args.Client == entry.Client {
			reply.Err = entry.Reply.Err
			reply.Value = entry.Reply.Value
			return nil
		}
	}

	// Retrieve the value from the store
	if value, exists := pb.store[args.Key]; exists {
		reply.Err = OK
		reply.Value = value
	} else {
		reply.Err = ErrNoKey
		reply.Value = ""
	}

	// Forward the operation to backup server if one exists
	if pb.view.Backup != "" {
		fargs := ForwardGetArgs{Args: args, Reply: reply}
		var freply ForwardReply

		ok := call(pb.view.Backup, "PBServer.ForwardGet", &fargs, &freply)
		if !ok || freply.Err != OK {
			reply.Err = ErrBackup
			reply.Value = ""
			return nil
		}
	}

	// Log the operation for duplicate detection
	pb.gets[args.Id] = GetEntry{Reply: *reply, Client: args.Client}

	return nil
}

// ForwardGet handles forwarded Get operations from the primary server.
// Only the backup server can accept forwarded operations.
// The backup only needs to log the reply for consistency, not perform the actual Get.
func (pb *PBServer) ForwardGet(args *ForwardGetArgs, reply *ForwardReply) error {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	DPrintf("Server ForwardGet(%s)\n", args.Args.Key)

	// Only the backup server can accept forwarded operations
	if pb.me != pb.view.Backup {
		reply.Err = ErrWrongServer
		return nil
	}

	// The backup only needs to log the reply for consistency
	pb.gets[args.Args.Id] = GetEntry{Reply: *args.Reply, Client: args.Args.Client}

	reply.Err = OK
	return nil
}

// tick pings the view service periodically to discover view changes.
// When the view changes and this server becomes primary with a new backup,
// it forwards its complete state to the new backup server.
func (pb *PBServer) tick() {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	DPrintf("-\n")

	args := viewservice.PingArgs{Me: pb.me, Viewnum: pb.view.Viewnum}
	var reply viewservice.PingReply

	// Ping the view service to get current view
	if ok := call(pb.vshost, "ViewServer.Ping", &args, &reply); !ok {
		return
	}

	// If we are the primary and the view has changed, forward our state to the new backup
	if reply.View.Primary == pb.me && pb.view.Viewnum != reply.View.Viewnum {
		if reply.View.Backup != "" {
			fargs := ForwardStateArgs{
				Primary: pb.me,
				Store:   pb.store,
				Gets:    pb.gets,
				Puts:    pb.puts,
			}
			var freply ForwardReply

			ok := call(reply.View.Backup, "PBServer.ForwardState", &fargs, &freply)
			if !ok || freply.Err != OK {
				return
			}
		}
	}

	pb.view = reply.View
}

// ForwardState handles complete state transfer from the primary server.
// Only the backup server can accept state transfers.
// It replaces the backup's local state with the primary's state.
func (pb *PBServer) ForwardState(args *ForwardStateArgs, reply *ForwardReply) error {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	DPrintf("Server ForwardState()\n")

	// Only the backup can accept state and only from the current primary
	if pb.me != pb.view.Backup {
		reply.Err = ErrWrongServer
		return nil
	}

	// Replace backup's state with primary's state
	pb.store = args.Store
	pb.gets = args.Gets
	pb.puts = args.Puts

	reply.Err = OK
	return nil
}

// Kill tells the server to shut itself down.
// This method is used for testing and graceful shutdown.
func (pb *PBServer) Kill() {
	pb.dead = true
	pb.l.Close()
}

// StartServer creates and starts a new primary/backup server.
// vshost is the address of the view service, and me is the server's identifier.
// The server will start listening for RPC connections and begin pinging the view service.
func StartServer(vshost string, me string) *PBServer {
	pb := &PBServer{
		me:     me,
		vs:     viewservice.MakeClerk(me, vshost),
		finish: make(chan interface{}),
		vshost: vshost,
		view:   viewservice.View{Viewnum: 0, Primary: "", Backup: ""},
		store:  make(map[string]string),
		gets:   make(map[uint64]GetEntry),
		puts:   make(map[uint64]PutEntry),
	}

	rpcs := rpc.NewServer()
	rpcs.Register(pb)

	// Remove any existing socket file and create a new listener
	os.Remove(pb.me)
	l, err := net.Listen("unix", pb.me)
	if err != nil {
		log.Fatal("listen error: ", err)
	}
	pb.l = l

	// Start the RPC server goroutine
	go func() {
		for !pb.dead {
			conn, err := pb.l.Accept()
			if err == nil && !pb.dead {
				if pb.unreliable && (rand.Int63()%1000) < 100 {
					// Discard the request (simulate network failure)
					conn.Close()
				} else if pb.unreliable && (rand.Int63()%1000) < 200 {
					// Process the request but force discard of reply
					c1 := conn.(*net.UnixConn)
					f, _ := c1.File()
					err := syscall.Shutdown(int(f.Fd()), syscall.SHUT_WR)
					if err != nil {
						fmt.Printf("shutdown: %v\n", err)
					}
					pb.done.Add(1)
					go func() {
						rpcs.ServeConn(conn)
						pb.done.Done()
					}()
				} else {
					// Normal processing
					pb.done.Add(1)
					go func() {
						rpcs.ServeConn(conn)
						pb.done.Done()
					}()
				}
			} else if err == nil {
				conn.Close()
			}
			if err != nil && !pb.dead {
				fmt.Printf("PBServer(%v) accept: %v\n", me, err.Error())
				pb.Kill()
			}
		}
		DPrintf("%s: wait until all requests are done\n", pb.me)
		pb.done.Wait()
		close(pb.finish)
	}()

	// Start the tick goroutine for pinging the view service
	pb.done.Add(1)
	go func() {
		for !pb.dead {
			pb.tick()
			time.Sleep(viewservice.PingInterval)
		}
		pb.done.Done()
	}()

	return pb
}
