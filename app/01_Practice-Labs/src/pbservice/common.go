// Package pbservice implements a primary/backup key-value service with fault tolerance.
// This service provides a distributed key-value store where one server acts as the
// primary (handling all client requests) and another acts as a backup (replicating
// the primary's state). The service ensures at-most-once semantics for operations
// and maintains consistency between primary and backup servers.
package pbservice

import "hash/fnv"

// Error constants for the primary/backup service
const (
	// OK indicates successful operation
	OK = "OK"
	// ErrNoKey indicates the requested key does not exist
	ErrNoKey = "ErrNoKey"
	// ErrWrongServer indicates the request was sent to the wrong server
	ErrWrongServer = "ErrWrongServer"
	// ErrBackup indicates a backup operation failed
	ErrBackup = "ErrBackup"
)

// Err represents an error type for the primary/backup service
type Err string

// PutArgs contains the arguments for a Put or PutHash operation
type PutArgs struct {
	Key    string // The key to store the value under
	Value  string // The value to store
	DoHash bool   // If true, perform PutHash operation (hash previous value + new value)
	Id     uint64 // Unique operation ID for duplicate detection
	Client string // Client identifier for duplicate detection
}

// PutReply contains the response for a Put or PutHash operation
type PutReply struct {
	Err           Err    // Error status of the operation
	PreviousValue string // Previous value (for PutHash operations)
}

// GetArgs contains the arguments for a Get operation
type GetArgs struct {
	Key    string // The key to retrieve
	Id     uint64 // Unique operation ID for duplicate detection
	Client string // Client identifier for duplicate detection
}

// GetReply contains the response for a Get operation
type GetReply struct {
	Err   Err    // Error status of the operation
	Value string // The retrieved value (empty if key doesn't exist)
}

// PutEntry represents a logged Put operation for duplicate detection
type PutEntry struct {
	Reply  PutReply // The reply that was sent for this operation
	Client string   // Client that performed the operation
}

// GetEntry represents a logged Get operation for duplicate detection
type GetEntry struct {
	Reply  GetReply // The reply that was sent for this operation
	Client string   // Client that performed the operation
}

// ForwardPutArgs contains arguments for forwarding a Put operation from primary to backup
type ForwardPutArgs struct {
	Args  *PutArgs  // Original Put arguments
	Reply *PutReply // Reply to be sent back to client
}

// ForwardGetArgs contains arguments for forwarding a Get operation from primary to backup
type ForwardGetArgs struct {
	Args  *GetArgs  // Original Get arguments
	Reply *GetReply // Reply to be sent back to client
}

// ForwardStateArgs contains arguments for forwarding complete state from primary to backup
type ForwardStateArgs struct {
	Primary string                // Address of the primary server
	Store   map[string]string     // Complete key-value store
	Gets    map[uint64]GetEntry   // Complete log of Get operations
	Puts    map[uint64]PutEntry   // Complete log of Put operations
}

// ForwardReply contains the response for forwarded operations
type ForwardReply struct {
	Err Err // Error status of the forwarded operation
}

// hash computes a 32-bit hash of the input string using FNV-1a algorithm.
// This function is used by PutHash to generate deterministic hash values.
func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
