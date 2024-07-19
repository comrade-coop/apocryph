package main

import (
	"context"
	"fmt"
	"log"
	"time"

	tpraft "github.com/comrade-coop/apocryph/pkg/raft"
	"github.com/hashicorp/raft"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/spf13/cobra"
)

var peersMultiAddr []string
var port uint16
var raftPath string
var appDomain string
var gateway string

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the autonomous autoscaler",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("Forming a Raft Cluster with the following peers:", peersMultiAddr)

		self, _ := tpraft.NewPeer(fmt.Sprintf("/ip4/0.0.0.0/tcp/%v", port))
		var peers []*host.Host
		for _, addr := range peersMultiAddr {
			peer, _ := tpraft.NewPeer(addr)
			peers = append(peers, &peer)
			log.Printf("Peer %s ID:%s \n", addr, peer.ID().String())
		}

		node, err := tpraft.NewRaftNode(self, peers, raftPath)
		if err != nil {
			return fmt.Errorf("Failed creating Raft node: %s", err)
		}

		// create the KVStore
		store, err := tpraft.NewKVStore(node)

		log.Println("Waiting for leader election")
		time.Sleep(5 * time.Second)

		setKey := func() {
			if node.Raft.State() == raft.Leader {
				log.Println("Setting the domain value to the current apocryph node gateway")
				err := store.Set(appDomain, fmt.Sprintf("%v", gateway))
				// leadership could be changed right before setting the value
				if _, ok := err.(*tpraft.NotLeaderError); ok {
					newLeaderAddr, newLeaderID := node.Raft.LeaderWithID()
					log.Printf("Leadership changed to %v:%v", newLeaderAddr, newLeaderID)
				} else {
					log.Printf("Failed setting key:%v: %v\n", appDomain, err)
				}
			}
		}

		// print the new state whenever it changes
		printNewState := func() {
			for range node.Consensus.Subscribe() {
				newState, err := node.Consensus.GetCurrentState()
				if err != nil {
					log.Printf("Failed getting current state %v", err)
				}
				log.Printf("State Changed, New State: %v", newState)
			}
		}

		// using Raft.leaderCh() wont help because it does not count
		// first leader election as a leadership change, therefore from the docs,
		// this is the way to detect the new leader
		obsCh := make(chan raft.Observation, 1)
		observer := raft.NewObserver(obsCh, false, nil)
		node.Raft.RegisterObserver(observer)
		defer node.Raft.DeregisterObserver(observer)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		ticker := time.NewTicker(time.Second / 2)
		defer ticker.Stop()
		for {
			select {
			case obs := <-obsCh:
				switch obs.Data.(type) {
				case raft.RaftState:
					if leaderAddr, _ := node.Raft.LeaderWithID(); leaderAddr != "" {
						setKey()
						printNewState()
					}
				}
			case <-ticker.C:
				if leaderAddr, _ := node.Raft.LeaderWithID(); leaderAddr != "" {
					setKey()
					printNewState()
				}
			case <-ctx.Done():
				fmt.Println("timed out waiting for Leader")
				return nil
			}
		}
	},
}

func init() {
	startCmd.Flags().StringSliceVarP(&peersMultiAddr, "peers", "p", []string{}, "List of peers multiaddresses that the autoscaler will use to redeploy your application in case of a failure")
	startCmd.Flags().Uint16Var(&port, "port", 9999, "port number for this node")
	startCmd.Flags().StringVar(&raftPath, "path", "", "path where raft will save it's state (Default is In Memory)")
	startCmd.Flags().StringVar(&appDomain, "doamin", "www.apocryph.com", "the application domain name")
	startCmd.Flags().StringVar(&gateway, "gateway", "127.0.0.1", "Apocryph node Gateway")
}
