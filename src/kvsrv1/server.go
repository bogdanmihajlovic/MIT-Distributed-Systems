package kvsrv

import (
	"log"
	"sync"

	"6.5840/kvsrv1/rpc"
	"6.5840/labrpc"
	tester "6.5840/tester1"
)

const Debug = false

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

type KVServer struct {
	mu sync.Mutex

	values   map[string]string
	versions map[string]rpc.Tversion
}

func MakeKVServer() *KVServer {
	kv := &KVServer{}
	// Your code here.
	kv.values = make(map[string]string)
	kv.versions = make(map[string]rpc.Tversion)
	return kv
}

// Get returns the value and version for args.Key, if args.Key
// exists. Otherwise, Get returns ErrNoKey.
func (kv *KVServer) Get(args *rpc.GetArgs, reply *rpc.GetReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	value, exists := kv.values[args.Key]
	if exists {
		reply.Err = rpc.OK
		reply.Value = value
		reply.Version = kv.versions[args.Key]
		return
	}
	reply.Err = rpc.ErrNoKey
}

// Update the value for a key if args.Version matches the version of
// the key on the server. If versions don't match, return ErrVersion.
// If the key doesn't exist, Put installs the value if the
// args.Version is 0, and returns ErrNoKey otherwise.
func (kv *KVServer) Put(args *rpc.PutArgs, reply *rpc.PutReply) {

	kv.mu.Lock()
	defer kv.mu.Unlock()
	// check for version
	version, exists := kv.versions[args.Key]
	if exists {
		if version == args.Version {
			kv.values[args.Key] = args.Value
			kv.versions[args.Key] = version + 1
			reply.Err = rpc.OK
		} else {
			reply.Err = rpc.ErrVersion
		}
	} else {
		if args.Version == 0 {
			kv.values[args.Key] = args.Value
			kv.versions[args.Key] = 1
			reply.Err = rpc.OK
		} else {
			reply.Err = rpc.ErrNoKey
		}
	}

}

// You can ignore Kill() for this lab
func (kv *KVServer) Kill() {
}

// You can ignore all arguments; they are for replicated KVservers
func StartKVServer(ends []*labrpc.ClientEnd, gid tester.Tgid, srv int, persister *tester.Persister) []tester.IService {
	kv := MakeKVServer()
	return []tester.IService{kv}
}
