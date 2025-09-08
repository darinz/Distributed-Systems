package pbservice

import (
	"bytes"
	"crypto/rand"
	"distributed-systems/app/01_Practice-Labs/src/viewservice"
	"encoding/binary"
	"fmt"
	"net/rpc"
	"time"
)

// Clerk represents a client for the primary/backup key-value service.
// It maintains a connection to the view service and handles client operations
// with automatic retry logic and primary server discovery.
type Clerk struct {
	vs     *viewservice.Clerk // View service client for discovering primary/backup
	vshost string             // Address of the view service
	me     string             // Client identifier
	view   viewservice.View   // Current view (primary/backup information)
}

// MakeClerk creates a new client for the primary/backup service.
// vshost is the address of the view service, and me is the client identifier.
func MakeClerk(vshost string, me string) *Clerk {
	ck := &Clerk{
		vs:     viewservice.MakeClerk(me, vshost),
		vshost: vshost,
		me:     me,
		view:   viewservice.View{Viewnum: 0, Primary: "", Backup: ""},
	}
	return ck
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

// uuid generates a 64-bit unique identifier for operation deduplication.
// It uses cryptographically secure random number generation to ensure uniqueness.
func (ck *Clerk) uuid() uint64 {
	// Generate 8 random bytes
	rbytes := make([]byte, 8)
	n, err := rand.Read(rbytes)
	for n < 8 || err != nil {
		n, err = rand.Read(rbytes)
	}

	// Convert bytes to uint64 using little-endian encoding
	var randid uint64
	binary.Read(bytes.NewReader(rbytes), binary.LittleEndian, &randid)

	return randid
}

// vsping pings the view service to get the current view information.
// It updates the client's view with the latest primary/backup assignments.
func (ck *Clerk) vsping() {
	args := viewservice.PingArgs{Me: ck.me, Viewnum: ck.view.Viewnum}
	var reply viewservice.PingReply

	// Ping the view service to get current view
	if ok := call(ck.vshost, "ViewServer.Ping", &args, &reply); !ok {
		return
	}
	ck.view = reply.View
}

// Get retrieves the value associated with the given key from the primary server.
// If the key has never been set, it returns an empty string.
// Get() will keep retrying until it either gets the value or the primary
// confirms the key doesn't exist. It handles primary server failures by
// automatically discovering the new primary through the view service.
func (ck *Clerk) Get(key string) string {
	DPrintf("Client Get(%s)\n", key)

	args := GetArgs{Key: key, Id: ck.uuid(), Client: ck.me}
	var reply GetReply

	// If this is the first request, ping for current view
	if ck.view.Viewnum == 0 {
		ck.vsping()
	}

	// Retry until the RPC is successful
	ok := call(ck.view.Primary, "PBServer.Get", &args, &reply)
	for !ok || (reply.Err != OK && reply.Err != ErrNoKey) {
		// Sleep for a tick to avoid overwhelming the system
		time.Sleep(viewservice.PingInterval)
		// Check if the view changed (new primary)
		ck.vsping()
		ok = call(ck.view.Primary, "PBServer.Get", &args, &reply)
	}

	return reply.Value
}

// PutExt performs a Put or PutHash operation on the primary server.
// It will keep retrying until the operation succeeds, handling primary
// server failures by automatically discovering the new primary.
// dohash determines whether to perform a PutHash operation (hash previous + new value).
// Returns the previous value (for PutHash operations).
func (ck *Clerk) PutExt(key string, value string, dohash bool) string {
	args := PutArgs{Key: key, Value: value, DoHash: dohash, Id: ck.uuid(), Client: ck.me}
	var reply PutReply

	// If this is the first request, ping for current view
	if ck.view.Viewnum == 0 {
		ck.vsping()
	}

	// Retry until the RPC is successful
	ok := call(ck.view.Primary, "PBServer.Put", &args, &reply)
	for !ok || reply.Err != OK {
		// Sleep for a tick to avoid overwhelming the system
		time.Sleep(viewservice.PingInterval)
		// Check if the view changed (new primary)
		ck.vsping()
		ok = call(ck.view.Primary, "PBServer.Put", &args, &reply)
	}

	return reply.PreviousValue
}

// Put stores a key-value pair in the primary/backup service.
// It will keep retrying until the operation succeeds.
func (ck *Clerk) Put(key string, value string) {
	DPrintf("Client Put(%s, %s)\n", key, value)
	ck.PutExt(key, value, false)
}

// PutHash performs a hash operation on the key's value.
// It concatenates the previous value with the new value and stores the hash.
// Returns the previous value before the hash operation.
func (ck *Clerk) PutHash(key string, value string) string {
	DPrintf("Client PutHash(%s, %s)\n", key, value)
	return ck.PutExt(key, value, true)
}
