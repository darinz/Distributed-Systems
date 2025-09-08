package shardmaster

import (
	"fmt"
	"net/rpc"
	"time"
)

// Clerk represents a client for the shard master service.
// It provides methods to interact with the shard master and handles
// automatic retry logic and server failover.
type Clerk struct {
	servers []string // List of shard master server addresses
}

// MakeClerk creates a new client for the shard master service.
// servers is a list of shard master server addresses that the client can connect to.
func MakeClerk(servers []string) *Clerk {
	return &Clerk{
		servers: servers,
	}
}

// call sends an RPC to the specified server and waits for a reply.
// It returns true if the server responded successfully, false otherwise.
// The reply argument should be a pointer to a reply structure.
// This function handles connection establishment, RPC call, and cleanup.
func call(srv string, rpcname string, args interface{}, reply interface{}) bool {
	c, err := rpc.Dial("unix", srv)
	if err != nil {
		return false
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}

// Query retrieves a configuration from the shard master.
// num specifies the configuration number to retrieve (-1 for the latest configuration).
// Returns the requested configuration.
func (ck *Clerk) Query(num int) Config {
	for {
		// Try each known server until one responds
		for _, srv := range ck.servers {
			args := &QueryArgs{Num: num}
			var reply QueryReply
			ok := call(srv, "ShardMaster.Query", args, &reply)
			if ok {
				return reply.Config
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// Join requests that a replica group join the system.
// gid is the unique group identifier, and servers is the list of server addresses in the group.
// The shard master will create a new configuration that includes this group.
func (ck *Clerk) Join(gid int64, servers []string) {
	for {
		// Try each known server until one responds
		for _, srv := range ck.servers {
			args := &JoinArgs{GID: gid, Servers: servers}
			var reply JoinReply
			ok := call(srv, "ShardMaster.Join", args, &reply)
			if ok {
				return
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// Leave requests that a replica group leave the system.
// gid is the group identifier of the group that wants to leave.
// The shard master will create a new configuration that excludes this group
// and redistributes its shards to the remaining groups.
func (ck *Clerk) Leave(gid int64) {
	for {
		// Try each known server until one responds
		for _, srv := range ck.servers {
			args := &LeaveArgs{GID: gid}
			var reply LeaveReply
			ok := call(srv, "ShardMaster.Leave", args, &reply)
			if ok {
				return
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// Move requests that a specific shard be moved to a specific replica group.
// shard is the shard number to move, and gid is the group ID to assign it to.
// This is primarily used for testing and fine-tuning load balance.
func (ck *Clerk) Move(shard int, gid int64) {
	for {
		// Try each known server until one responds
		for _, srv := range ck.servers {
			args := &MoveArgs{Shard: shard, GID: gid}
			var reply MoveReply
			ok := call(srv, "ShardMaster.Move", args, &reply)
			if ok {
				return
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}
