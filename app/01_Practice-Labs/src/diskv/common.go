package diskv

import "crypto/rand"
import "math/big"

//
// Sharded key/value server.
// Lots of replica groups, each running op-at-a-time paxos.
// Shardmaster decides which group serves each shard.
// Shardmaster may change shard assignment from time to time.
//
// You will have to modify these definitions.
//
// Shared types, error codes, and helpers for the persistent sharded key/value service.

const (
	OK            = "OK"
	ErrNoKey      = "ErrNoKey"
	ErrWrongGroup = "ErrWrongGroup"
	ErrNoConfig   = "ErrNoConfig"
	Nil           = ""
)

type Err string

type PutAppendArgs struct {
	Key   string
	Value string
	Op    string // "Put" or "Append"
	// You'll have to add definitions here.
	// Field names must start with capital letters,
	// otherwise RPC will break.
	Seq    int64 // per-client monotonic sequence for at-most-once
	Client int64 // unique client identifier
}

type PutAppendReply struct {
	Err Err
}

type GetArgs struct {
	Key string
	// You'll have to add definitions here.
	Seq    int64 // per-client monotonic sequence for at-most-once
	Client int64 // unique client identifier
}

type GetReply struct {
	Err   Err
	Value string
}

//
// Generates a 64-bit UUID - random number
//
func uuid() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}
