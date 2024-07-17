package main

import (
	"fmt"
	"time"

	"github.com/comrade-coop/apocryph/pkg/raft"
	"github.com/libp2p/go-libp2p/core/host"
)

func main() {

	peer1, _ := raft.NewPeer("/ip4/127.0.0.1/tcp/9997")
	fmt.Printf("Node ID:%s \n", peer1.ID().String())
	peer2, _ := raft.NewPeer("/ip4/127.0.0.1/tcp/9998")
	fmt.Printf("Node ID:%s \n", peer2.ID().String())
	peer3, _ := raft.NewPeer("/ip4/127.0.0.1/tcp/9999")
	fmt.Printf("Node ID:%s \n", peer3.ID().String())

	defer peer1.Close()
	defer peer2.Close()
	defer peer3.Close()

	peers := []*host.Host{&peer1, &peer2, &peer3}
	node1, err := raft.NewRaftNode(peer1, peers, "")
	if err != nil {
		fmt.Printf("Error:%s", err)
		return
	}
	node2, err := raft.NewRaftNode(peer2, peers, "")
	if err != nil {
		fmt.Printf("Error:%s", err)
		return
	}
	node3, err := raft.NewRaftNode(peer3, peers, "")
	if err != nil {
		fmt.Printf("Error:%s", err)
		return
	}
	// create the kv stores
	store1, _ := raft.NewKVStore(node1)
	store2, _ := raft.NewKVStore(node2)
	store3, _ := raft.NewKVStore(node3)

	fmt.Println("Waiting for leader election")
	// Provide some time for leader election
	time.Sleep(5 * time.Second)

	go func() {
		for range node1.Raft.LeaderCh() {
			fmt.Println("Leadership changed")
		}
	}()

	domain := "www.webapp.com"
	fmt.Printf("Trying to set the key:%v with node1\n", domain)
	err = store1.Set(domain, "/ip4/127.0.0.1/tcp/9997")
	if err != nil {
		if _, ok := err.(*raft.NotLeaderError); ok {
			fmt.Println("Node1 not a leader: Trying to Set with node2")
			err = store2.Set(domain, "/ip4/127.0.0.1/tcp/9998")
			if _, ok := err.(*raft.NotLeaderError); ok {
				fmt.Println("Node2 not a leader: Trying to Set with node3")
				err = store3.Set(domain, "/ip4/127.0.0.1/tcp/9999")
				if err != nil {
					panic(err)
				}
			} else {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	fmt.Println("Waiting for updates...")
	time.Sleep(5 * time.Second)

	value, _ := store1.Get(domain)
	value2, _ := store2.Get(domain)
	value3, _ := store3.Get(domain)

	fmt.Printf("Store1 domain value:%v\n", value)
	fmt.Printf("Store2 domain value:%v\n", value2)
	fmt.Printf("Store3 domain value:%v\n", value3)

	fmt.Printf("Deleting key:%v from the store...\n", domain)
	err = store1.Delete(domain)
	if err != nil {
		if _, ok := err.(*raft.NotLeaderError); ok {
			fmt.Println("Node1 not a leader: Trying to Delete with node2")
			err = store2.Delete(domain)
			if _, ok := err.(*raft.NotLeaderError); ok {
				fmt.Println("Node2 not a leader: Trying to Delete with node3")
				err = store3.Delete(domain)
				if err != nil {
					panic(err)
				}
			} else {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	fmt.Println("Waiting for updates...")
	time.Sleep(5 * time.Second)

	// Final states
	finalState1, err := node1.Consensus.GetCurrentState()
	if err != nil {
		fmt.Println(err)
		return
	}
	finalState2, err := node2.Consensus.GetCurrentState()
	if err != nil {
		fmt.Println(err)
		return
	}
	finalState3, err := node3.Consensus.GetCurrentState()
	if err != nil {
		fmt.Println(err)
		return
	}
	finalRaftState1 := finalState1.(*raft.KVStore)
	finalRaftState2 := finalState2.(*raft.KVStore)
	finalRaftState3 := finalState3.(*raft.KVStore)

	fmt.Printf("Raft1 final state: %v\n", finalRaftState1)
	fmt.Printf("Raft2 final state: %v\n", finalRaftState2)
	fmt.Printf("Raft3 final state: %v\n", finalRaftState3)

	node1.Raft.Shutdown().Error()
	node2.Raft.Shutdown().Error()
	node3.Raft.Shutdown().Error()

}
