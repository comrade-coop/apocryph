package raft

import (
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb/v2"
	"github.com/libp2p/go-libp2p"
	raftp2p "github.com/libp2p/go-libp2p-raft"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
)

const (
	SNAPSHOT_RETAIN = 3
	RAFT_TIMEOUT    = 10 * time.Second
)

type TestState struct {
	Value int
}

type RaftNode struct {
	Raft        *raft.Raft
	Consensus   *raftp2p.Consensus
	Actor       *raftp2p.Actor
	Transport   *raft.NetworkTransport
	Config      raft.Config
	Snapshots   raft.SnapshotStore
	StableStore raft.StableStore
	LogStore    raft.LogStore
}

func NewRaftNode(host host.Host, peers []*peer.AddrInfo, raftDir string) (*RaftNode, error) {
	for _, peerPtr := range peers {
		peer := *peerPtr
		// skip adding self
		if peer.Addrs[0].String() != host.Addrs()[0].String() {
			host.Peerstore().AddAddrs(peer.ID, peer.Addrs, peerstore.PermanentAddrTTL)
		}
	}

	transport, err := raftp2p.NewLibp2pTransport(host, time.Minute)
	if err != nil {
		return nil, fmt.Errorf("Failed creating transport:%s", err)
	}

	config := raft.DefaultConfig()
	config.LogOutput = io.Discard
	config.Logger = nil
	config.LocalID = raft.ServerID(host.ID().String())

	raftNode := RaftNode{
		Transport: transport,
	}

	if raftDir == "" {
		raftNode.LogStore = raft.NewInmemStore()
		raftNode.StableStore = raft.NewInmemStore()
		raftNode.Snapshots = raft.NewInmemSnapshotStore()
	} else {
		boltDB, err := raftboltdb.New(raftboltdb.Options{
			Path: filepath.Join(raftDir, "raft.db"),
		})
		if err != nil {
			return nil, fmt.Errorf("new bolt store: %s", err)
		}
		raftNode.LogStore = boltDB
		raftNode.StableStore = boltDB
		raftNode.Snapshots, err = raft.NewFileSnapshotStore(filepath.Join(raftDir, "raft_snapshots"), SNAPSHOT_RETAIN, nil)
		if err != nil {
			return nil, fmt.Errorf("Failed creating file snapshot store")
		}
	}
	servers := make([]raft.Server, 0)
	for _, peerPtr := range peers {
		peer := *peerPtr
		servers = append(servers, raft.Server{
			Suffrage: raft.Voter,
			ID:       raft.ServerID(peer.ID.String()),
			Address:  raft.ServerAddress(peer.ID.String()),
		})
	}
	serversCfg := raft.Configuration{Servers: servers}

	// initializes a server's storage with the given cluster configuration,
	err = raft.BootstrapCluster(config, raftNode.LogStore, raftNode.StableStore, raftNode.Snapshots, raftNode.Transport, serversCfg.Clone())
	if err != nil {
		return nil, fmt.Errorf("Failed bootstraping cluster: %v", err)
	}
	raftNode.Config = *config
	return &raftNode, nil
}

func (node *RaftNode) StartNode(c *raftp2p.Consensus) error {
	node.Consensus = c
	instance, err := raft.NewRaft(&node.Config, node.Consensus.FSM(), node.LogStore, node.StableStore, node.Snapshots, node.Transport)
	if err != nil {
		return fmt.Errorf("Failed creating Raft instance:%v", err)
	}
	node.Raft = instance
	node.Actor = raftp2p.NewActor(node.Raft)
	node.Consensus.SetActor(node.Actor)
	return nil
}

func NewPeer(multiaddr string) (host.Host, error) {
	h, err := libp2p.New(libp2p.ListenAddrStrings(multiaddr))
	if err != nil {
		return nil, err
	}
	return h, nil
}
