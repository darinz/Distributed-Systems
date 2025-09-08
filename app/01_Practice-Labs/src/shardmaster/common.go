// Package shardmaster implements a fault-tolerant shard master service that assigns
// shards to replication groups in a distributed key-value storage system.
// The shard master manages configurations that describe which replica groups
// are responsible for which shards, and handles dynamic reconfiguration as
// groups join and leave the system.
package shardmaster

// NShards defines the total number of shards in the system.
// Each shard contains a subset of the key-value pairs, and shards are
// distributed among replica groups to balance load and improve performance.
const NShards = 10

// Config represents a configuration that describes the current state of the system.
// Each configuration is numbered and contains information about which replica
// groups exist and which shards each group is responsible for.
type Config struct {
	Num    int                // Configuration number (starts at 0)
	Shards [NShards]int64     // Array mapping shard index to group ID (GID)
	Groups map[int64][]string // Map from group ID to list of server addresses
}

// JoinArgs contains the arguments for a Join RPC request.
// Join is used when a new replica group wants to join the system.
type JoinArgs struct {
	GID     int64    // Unique replica group identifier (must be > 0)
	Servers []string // List of server addresses in the replica group
}

// JoinReply contains the response for a Join RPC request.
// Currently empty as Join operations don't return data.
type JoinReply struct {
}

// LeaveArgs contains the arguments for a Leave RPC request.
// Leave is used when a replica group wants to leave the system.
type LeaveArgs struct {
	GID int64 // Group ID of the replica group that wants to leave
}

// LeaveReply contains the response for a Leave RPC request.
// Currently empty as Leave operations don't return data.
type LeaveReply struct {
}

// MoveArgs contains the arguments for a Move RPC request.
// Move is used to manually assign a specific shard to a specific group.
type MoveArgs struct {
	Shard int   // Shard number to move (0 to NShards-1)
	GID   int64 // Group ID to assign the shard to
}

// MoveReply contains the response for a Move RPC request.
// Currently empty as Move operations don't return data.
type MoveReply struct {
}

// QueryArgs contains the arguments for a Query RPC request.
// Query is used to retrieve configuration information.
type QueryArgs struct {
	Num int // Configuration number to query (-1 for latest configuration)
}

// QueryReply contains the response for a Query RPC request.
// Returns the requested configuration.
type QueryReply struct {
	Config Config // The requested configuration
}
