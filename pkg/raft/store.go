package raft

import (
	"fmt"
	"sync"

	"github.com/hashicorp/raft"
	consensus "github.com/libp2p/go-libp2p-consensus"
	raftp2p "github.com/libp2p/go-libp2p-raft"
)

// Data is exported because it needs to be serialized
type KVStore struct {
	Data map[string]string
	mu   sync.Mutex
	node *RaftNode
}

// all fields must be public otherwise raft will dismiss them in (de)serilization
type KVOperation struct {
	Op    string
	Key   string
	Value string
}

func NewKVStore(node *RaftNode) (*KVStore, error) {
	store := KVStore{Data: make(map[string]string), node: node}
	// c := raftp2p.NewConsensus(&KVStore{})
	c := raftp2p.NewOpLog(&store, &KVOperation{})
	err := node.StartNode(c)
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (s *KVOperation) ApplyTo(state consensus.State) (consensus.State, error) {
	// get the underlying store
	store := state.(*KVStore)
	switch s.Op {
	case "set":
		err := store.applySet(s.Key, s.Value)
		if err != nil {
			return store, err
		}
		return store, nil
	case "delete":
		err := store.applyDelete(s.Key)
		if err != nil {
			return store, err
		}
		return store, nil
	default:
		return s, fmt.Errorf("unrecognized command op: %s", s.Op)
	}
}

// Get returns the value for the given key.
func (s *KVStore) Get(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Data[key], nil
}

func (s *KVStore) Set(key, value string) error {
	if s.node.Raft.State() != raft.Leader {
		return &NotLeaderError{}
	}

	op := KVOperation{
		Op:    "set",
		Key:   key,
		Value: value,
	}
	_, err := s.node.Consensus.CommitOp(&op)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes the given key.
func (s *KVStore) Delete(key string) error {
	if s.node.Raft.State() != raft.Leader {
		return &NotLeaderError{}
	}
	op := KVOperation{
		Op:  "delete",
		Key: key,
	}
	_, err := s.node.Consensus.CommitOp(&op)
	if err != nil {
		return err
	}
	return nil
}

func (store *KVStore) applySet(key, value string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.Data[key] = value
	return nil
}

func (store *KVStore) applyDelete(key string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	delete(store.Data, key)
	return nil
}

type NotLeaderError struct{}

// Implement the Error method for NotLeaderError
func (e *NotLeaderError) Error() string {
	return "Not leader"
}
