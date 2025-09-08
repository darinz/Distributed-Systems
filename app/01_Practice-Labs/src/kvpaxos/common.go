// Package kvpaxos implements a fault-tolerant key-value storage system using the Paxos consensus algorithm.
// This system provides sequential consistency and at-most-once semantics for all operations.
// Multiple replicas maintain identical state through Paxos agreement on operation ordering.
package kvpaxos

import "hash/fnv"

// Error constants for the key-value service
const (
	// OK indicates successful operation
	OK = "OK"
	// ErrNoKey indicates the requested key does not exist
	ErrNoKey = "ErrNoKey"
)

// Err represents an error type for the key-value service
type Err string

// PutArgs contains the arguments for a Put or PutHash operation
type PutArgs struct {
	Key    string // The key to store the value under
	Value  string // The value to store
	DoHash bool   // If true, perform PutHash operation (hash previous value + new value)
	Client uint64 // Client identifier for duplicate detection
	OpId   uint64 // Unique operation ID for duplicate detection
}

// PutReply contains the response for a Put or PutHash operation
type PutReply struct {
	Err           Err    // Error status of the operation
	PreviousValue string // Previous value (for PutHash operations)
}

// GetArgs contains the arguments for a Get operation
type GetArgs struct {
	Key    string // The key to retrieve
	Client uint64 // Client identifier for duplicate detection
	OpId   uint64 // Unique operation ID for duplicate detection
}

// GetReply contains the response for a Get operation
type GetReply struct {
	Err   Err    // Error status of the operation
	Value string // The retrieved value (empty if key doesn't exist)
}

// hash computes a 32-bit hash of the input string using FNV-1a algorithm.
// This function is used by PutHash to generate deterministic hash values.
func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
