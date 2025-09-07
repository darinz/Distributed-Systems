// Package viewservice implements a non-replicated view service for a simple
// primary/backup key-value system. This service maintains a sequence of numbered
// views, each containing a primary server and optionally a backup server.
//
// The view service ensures that:
//   - The primary in a view is always either the primary or backup of the previous view
//   - At most one primary is active at a time through acknowledgment mechanisms
//   - Failed servers are detected and replaced automatically
//   - New views are only created after the current primary acknowledges the current view
//
// This design provides fault tolerance for the primary/backup key-value service
// while maintaining consistency guarantees.
package viewservice

import "time"

// View represents a numbered view in the primary/backup system.
// Each view contains a view number and the network addresses of the
// primary and backup servers for that view.
type View struct {
	Viewnum uint   // Sequential view number, starting from 1
	Primary string // Network address (host:port) of the primary server
	Backup  string // Network address (host:port) of the backup server, empty if none
}

// PingInterval defines how frequently servers should send Ping RPCs
// to the view service to indicate they are alive.
const PingInterval = time.Millisecond * 100

// DeadPings defines the number of consecutive missed Ping RPCs
// before the view service considers a server dead.
const DeadPings = 5

// PingArgs contains the arguments for the Ping RPC.
// Servers use this to:
//   - Inform the view service they are alive
//   - Report their current view number
//   - Learn about the latest view from the view service
type PingArgs struct {
	Me      string // Server's network address (host:port)
	Viewnum uint   // Server's current view number (0 indicates server restart/crash)
}

// PingReply contains the response from the Ping RPC.
// The view service returns the current view to the requesting server.
type PingReply struct {
	View View // Current view from the view service's perspective
}

// GetArgs contains the arguments for the Get RPC.
// This is used by clients to fetch the current view without
// volunteering to be a server.
type GetArgs struct {
	// No arguments needed for Get RPC
}

// GetReply contains the response from the Get RPC.
// Returns the current view without affecting server state.
type GetReply struct {
	View View // Current view from the view service's perspective
}
