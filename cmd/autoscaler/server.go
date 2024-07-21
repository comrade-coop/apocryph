package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"connectrpc.com/connect"
	pb "github.com/comrade-coop/apocryph/pkg/proto"
	pbcon "github.com/comrade-coop/apocryph/pkg/proto/protoconnect"
	tpraft "github.com/comrade-coop/apocryph/pkg/raft"
	"github.com/hashicorp/raft"
	"github.com/libp2p/go-libp2p/core/host"
)

const RAFT_P2P_PORT = 9999
const RAFT_PATH = ""

type AutoScalerServer struct {
	pbcon.UnimplementedAutoscalerServiceHandler
	node        *tpraft.RaftNode
	store       *tpraft.KVStore
	peers       []string
	self        host.Host
	nodeGateway string
	mainLoop    func(*AutoScalerServer)
}

func (s *AutoScalerServer) ConnectCluster(c context.Context, req *connect.Request[pb.ConnectClusterRequest]) (*connect.Response[pb.ConnectClusterResponse], error) {
	log.Println("Forming a Raft Cluster with the following providers:", req.Msg.Servers)

	var peers []*host.Host
	for _, addr := range req.Msg.Servers {
		addr = fmt.Sprintf("/ip4/%v/tcp/%v", addr, RAFT_P2P_PORT)
		log.Printf("Adding Peer:%v\n", addr)
		peer, err := tpraft.NewPeer(addr)
		if err != nil {
			return connect.NewResponse(&pb.ConnectClusterResponse{
				Success: false,
				Error:   fmt.Sprintf("Failed creating peer %v: %v", addr, err),
			}), nil
		}
		peers = append(peers, &peer)
		log.Printf("Peer %s ID:%s \n", addr, peer.ID().String())
	}

	node, err := tpraft.NewRaftNode(s.self, peers, RAFT_PATH)
	if err != nil {
		log.Println("Error:Could not Creat Raft Node")
		return connect.NewResponse(&pb.ConnectClusterResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed Creating Raft Node %v\n", err),
		}), nil

	}

	// create the KVStore
	store, err := tpraft.NewKVStore(node)

	s.node = node
	s.store = store
	s.peers = req.Msg.Servers
	s.nodeGateway = req.Msg.NodeGateway

	err = s.waitLeaderElection(req.Msg.Timeout)
	if err != nil {
		response := &pb.ConnectClusterResponse{Success: false, Error: err.Error()}
		return connect.NewResponse(response), nil
	}

	response := &pb.ConnectClusterResponse{Success: true}
	return connect.NewResponse(response), nil
}

func (s *AutoScalerServer) waitLeaderElection(timeout uint32) error {

	log.Printf("Waiting for leader election with %v seoncds timout ...", timeout)
	time.Sleep(5 * time.Second)

	go s.watchNewStates()

	// using Raft.leaderCh() wont help because it does not count
	// first leader election as a leadership change, therefore from the docs,
	// this is the way to detect the new leader
	obsCh := make(chan raft.Observation, 1)
	observer := raft.NewObserver(obsCh, false, nil)
	s.node.Raft.RegisterObserver(observer)
	defer s.node.Raft.DeregisterObserver(observer)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	timeoutCh := time.After(time.Duration(timeout) * time.Second)
	for {
		select {
		case obs := <-obsCh:
			switch obs.Data.(type) {
			case raft.RaftState:
				if leaderAddr, _ := s.node.Raft.LeaderWithID(); leaderAddr != "" {
					go s.mainLoop(s)
					return nil
				}
			}
		case <-ticker.C:
			if leaderAddr, _ := s.node.Raft.LeaderWithID(); leaderAddr != "" {
				go s.mainLoop(s)
				return nil
			}
		case <-timeoutCh:
			log.Println("timed out waiting for leader")
			return fmt.Errorf("Timed out waiting for leadership election")
		}
	}
}

func (s *AutoScalerServer) watchNewStates() {
	for range s.node.Consensus.Subscribe() {
		newState, err := s.node.Consensus.GetCurrentState()
		if err != nil {
			log.Printf("Failed getting current state %v\n", err)
		}
		log.Printf("State Changed, New State: %v\n", newState)
	}
}

// example of main loop setting the value of a test domain with the current node
// gateway every 5 seconds
func setAppGatewayExample(s *AutoScalerServer) {
	log.Println("Starting Main Loop:")
	for {
		if s.node.Raft.State() == raft.Leader {
			log.Println("Setting the domain value to the current apocryph node gateway")
			err := s.store.Set("www.test.com", s.nodeGateway)
			// leadership could be changed right before setting the value
			if _, ok := err.(*tpraft.NotLeaderError); ok {
				newLeaderAddr, newLeaderID := s.node.Raft.LeaderWithID()
				log.Printf("Leadership changed to %v:%v", newLeaderAddr, newLeaderID)
			} else {
				log.Printf("Failed setting key:%v: %v\n", "www.test.com", err)
			}
		}
		time.Sleep(5 * time.Second)
	}
}
