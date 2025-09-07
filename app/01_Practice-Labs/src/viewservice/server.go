package viewservice

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"sync"
	"time"
)

// ViewServer implements the view service for the primary/backup system.
// It maintains the current view state and manages view transitions based on
// server health and acknowledgment status.
type ViewServer struct {
	mu   sync.Mutex // Protects all fields below
	l    net.Listener
	dead bool
	me   string

	// View service state
	ticks   int            // Number of ticks since server initialization
	pings   map[string]int // Last tick when we heard from each server
	current View           // Current view state
	ack     bool           // Whether current view has been acknowledged by primary
}

// getIdleServer finds an available server that can serve as backup.
// It removes dead servers from the pings map as it searches.
// Returns an empty string if no suitable server is found.
func (vs *ViewServer) getIdleServer() string {
	for server, lastTick := range vs.pings {
		// Remove dead servers from tracking
		if vs.ticks-lastTick > DeadPings {
			delete(vs.pings, server)
			continue
		}
		
		// Check if this server is not already primary or backup
		if server != vs.current.Primary && server != vs.current.Backup {
			return server
		}
	}
	return ""
}

// Ping handles Ping RPC requests from servers.
// It updates server health tracking and manages view transitions based on:
// - Server health status (alive/dead)
// - View acknowledgment status
// - Server role (primary/backup/idle)
func (vs *ViewServer) Ping(args *PingArgs, reply *PingReply) error {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	server := args.Me

	// Handle first server initialization
	if vs.current.Viewnum == 0 {
		vs.current.Viewnum = 1
		vs.current.Primary = server
		vs.ack = false
		reply.View = vs.current
		return nil
	}

	// Update server health tracking
	vs.pings[server] = vs.ticks

	// Process based on server role and view number
	switch server {
	case vs.current.Primary:
		vs.handlePrimaryPing(args)
	case vs.current.Backup:
		vs.handleBackupPing(args)
	default:
		vs.handleIdleServerPing(args)
	}

	reply.View = vs.current
	return nil
}

// handlePrimaryPing processes ping from the current primary server.
func (vs *ViewServer) handlePrimaryPing(args *PingArgs) {
	switch args.Viewnum {
	case vs.current.Viewnum:
		// Primary acknowledges current view
		vs.ack = true
	case 0:
		// Primary crashed and restarted
		if !vs.ack {
			log.Fatal("primary crashed before acknowledging view")
		}
		vs.transitionToNewView(vs.current.Backup, vs.getIdleServer())
	}
}

// handleBackupPing processes ping from the current backup server.
func (vs *ViewServer) handleBackupPing(args *PingArgs) {
	if args.Viewnum == 0 {
		// Backup crashed and restarted
		if !vs.ack {
			// If the view hasn't been acknowledged yet, we can't safely transition
			// The backup will need to wait until the primary acknowledges the view
			return
		}
		vs.transitionToNewView(vs.current.Primary, vs.getIdleServer())
	}
}

// handleIdleServerPing processes ping from an idle server.
func (vs *ViewServer) handleIdleServerPing(args *PingArgs) {
	// Promote idle server to backup if no backup exists and view is acknowledged
	if vs.current.Backup == "" && vs.ack {
		vs.transitionToNewView(vs.current.Primary, args.Me)
	}
}

// transitionToNewView creates a new view with the specified primary and backup.
func (vs *ViewServer) transitionToNewView(primary, backup string) {
	vs.current.Viewnum++
	vs.current.Primary = primary
	vs.current.Backup = backup
	vs.ack = false
}

// Get handles Get RPC requests from clients.
// Returns the current view without affecting server state or health tracking.
func (vs *ViewServer) Get(args *GetArgs, reply *GetReply) error {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	
	reply.View = vs.current
	return nil
}

// tick is called periodically to check server health and manage view transitions.
// It detects failed servers and promotes backups or idle servers as needed.
// Only proceeds with view changes if the current view has been acknowledged.
func (vs *ViewServer) tick() {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	vs.ticks++

	// No view changes until we have at least one server and current view is acknowledged
	if vs.current.Viewnum == 0 || !vs.ack {
		return
	}

	// Check primary server health
	if vs.isServerDead(vs.current.Primary) {
		vs.handlePrimaryFailure()
		return
	}

	// Check backup server health
	if vs.current.Backup != "" && vs.isServerDead(vs.current.Backup) {
		vs.handleBackupFailure()
	}
}

// isServerDead checks if a server should be considered dead based on ping history.
func (vs *ViewServer) isServerDead(server string) bool {
	lastTick, exists := vs.pings[server]
	return !exists || (vs.ticks-lastTick > DeadPings)
}

// handlePrimaryFailure handles the case where the primary server has failed.
func (vs *ViewServer) handlePrimaryFailure() {
	delete(vs.pings, vs.current.Primary)
	vs.transitionToNewView(vs.current.Backup, vs.getIdleServer())
}

// handleBackupFailure handles the case where the backup server has failed.
func (vs *ViewServer) handleBackupFailure() {
	delete(vs.pings, vs.current.Backup)
	vs.transitionToNewView(vs.current.Primary, vs.getIdleServer())
}

// Kill shuts down the view server gracefully.
// This method is used for testing and should not be modified.
func (vs *ViewServer) Kill() {
	vs.dead = true
	vs.l.Close()
}

// StartServer creates and starts a new ViewServer instance.
// It initializes the server state, sets up RPC handling, and starts
// background goroutines for connection handling and periodic health checks.
func StartServer(me string) *ViewServer {
	vs := &ViewServer{
		me:      me,
		ticks:   0,
		pings:   make(map[string]int),
		current: View{Viewnum: 0, Primary: "", Backup: ""},
		ack:     false,
	}

	// Set up RPC server
	rpcs := rpc.NewServer()
	rpcs.Register(vs)

	// Set up network listener
	os.Remove(vs.me) // Clean up any existing unix socket
	l, err := net.Listen("unix", vs.me)
	if err != nil {
		log.Fatal("listen error: ", err)
	}
	vs.l = l

	// Start connection handling goroutine
	go vs.acceptConnections(rpcs)

	// Start periodic health check goroutine
	go vs.healthCheckLoop()

	return vs
}

// acceptConnections handles incoming RPC connections in a separate goroutine.
func (vs *ViewServer) acceptConnections(rpcs *rpc.Server) {
	for !vs.dead {
		conn, err := vs.l.Accept()
		if err != nil {
			if !vs.dead {
				fmt.Printf("ViewServer(%v) accept: %v\n", vs.me, err.Error())
				vs.Kill()
			}
			return
		}
		
		if !vs.dead {
			go rpcs.ServeConn(conn)
		} else {
			conn.Close()
		}
	}
}

// healthCheckLoop runs the periodic health check in a separate goroutine.
func (vs *ViewServer) healthCheckLoop() {
	for !vs.dead {
		vs.tick()
		time.Sleep(PingInterval)
	}
}
