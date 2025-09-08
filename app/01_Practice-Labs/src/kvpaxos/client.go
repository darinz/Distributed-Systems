package kvpaxos

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net/rpc"
	"time"
)

// Clerk represents a client for the Paxos-based key-value service.
// It maintains a connection to multiple server replicas and handles
// client operations with automatic retry logic and server failover.
type Clerk struct {
	servers []string // List of available server addresses
	me      uint64   // Stable client identifier for all RPC calls
}

// MakeClerk creates a new client for the Paxos-based key-value service.
// servers is a list of server addresses that the client can connect to.
func MakeClerk(servers []string) *Clerk {
	ck := &Clerk{
		servers: servers,
		me:      0, // Will be set by uuid()
	}
	ck.me = ck.uuid()
	return ck
}

// uuid generates a 64-bit unique identifier for client identification.
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

// Get retrieves the value associated with the given key.
// Returns an empty string if the key does not exist.
// Keeps trying different servers until one responds successfully.
func (ck *Clerk) Get(key string) string {
	DPrintf("Client Get(%s)\n", key)
	args := GetArgs{Key: key, Client: ck.me, OpId: ck.uuid()}
	var reply GetReply

	// Try different servers until one responds
	for i := 0; true; i++ {
		ok := call(ck.servers[i%len(ck.servers)], "KVPaxos.Get", args, &reply)
		if ok {
			return reply.Value
		}
		time.Sleep(100 * time.Millisecond)
	}
	// This should never be reached, but Go requires a return statement
	return ""
}

// PutExt performs a Put or PutHash operation on the key-value store.
// It keeps trying different servers until one responds successfully.
// dohash determines whether to perform a PutHash operation (hash previous + new value).
// Returns the previous value (for PutHash operations).
func (ck *Clerk) PutExt(key string, value string, dohash bool) string {
	args := PutArgs{Key: key, Value: value, DoHash: dohash, Client: ck.me, OpId: ck.uuid()}
	var reply PutReply

	for i := 0; true; i++ {
		ok := call(ck.servers[i%len(ck.servers)], "KVPaxos.Put", args, &reply)
		if ok {
			return reply.PreviousValue
		}
		time.Sleep(100 * time.Millisecond)
	}
	// This should never be reached, but Go requires a return statement
	return ""
}

// Put stores a key-value pair in the key-value store.
// It will keep trying until the operation succeeds.
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
