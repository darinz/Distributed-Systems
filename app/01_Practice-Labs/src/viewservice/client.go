package viewservice

import (
	"fmt"
	"net/rpc"
)

// Clerk provides a client interface to the view service.
// It maintains connection information and provides methods to interact
// with the view service for both servers and clients.
type Clerk struct {
	me     string // Client's network address (host:port)
	server string // View service's network address (host:port)
}

// MakeClerk creates a new Clerk instance for communicating with the view service.
// The 'me' parameter identifies this client, while 'server' specifies the
// view service's network address.
func MakeClerk(me string, server string) *Clerk {
	return &Clerk{
		me:     me,
		server: server,
	}
}

// call sends an RPC to the specified server and waits for a reply.
// It handles connection establishment, RPC execution, and cleanup.
// Returns true if the RPC succeeded, false otherwise.
// The reply argument must be a pointer to the expected reply structure.
//
// This function is used by all RPC calls in the viewservice package.
// It automatically handles timeouts and connection errors.
func call(srv string, rpcname string, args interface{}, reply interface{}) bool {
	// Establish connection to the server
	client, err := rpc.Dial("unix", srv)
	if err != nil {
		return false
	}
	defer client.Close()

	// Make the RPC call
	err = client.Call(rpcname, args, reply)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

// Ping sends a ping to the view service to indicate this server is alive
// and to learn about the current view. The viewnum parameter indicates
// the server's current view number (0 for server restart/crash).
// Returns the current view from the view service or an error if the RPC fails.
func (ck *Clerk) Ping(viewnum uint) (View, error) {
	args := &PingArgs{
		Me:      ck.me,
		Viewnum: viewnum,
	}
	var reply PingReply

	ok := call(ck.server, "ViewServer.Ping", args, &reply)
	if !ok {
		return View{}, fmt.Errorf("Ping(%v) failed", viewnum)
	}

	return reply.View, nil
}

// Get retrieves the current view from the view service without
// affecting server state or health tracking.
// Returns the current view and a boolean indicating success.
func (ck *Clerk) Get() (View, bool) {
	args := &GetArgs{}
	var reply GetReply
	
	ok := call(ck.server, "ViewServer.Get", args, &reply)
	if !ok {
		return View{}, false
	}
	
	return reply.View, true
}

// Primary returns the network address of the current primary server.
// This is a convenience method that calls Get() and extracts the primary.
// Returns an empty string if the RPC fails.
func (ck *Clerk) Primary() string {
	view, ok := ck.Get()
	if !ok {
		return ""
	}
	return view.Primary
}
