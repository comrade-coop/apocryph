// SPDX-License-Identifier: GPL-3.0

package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"connectrpc.com/connect"
	tpipfs "github.com/comrade-coop/trusted-pods/pkg/ipfs"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	pbcon "github.com/comrade-coop/trusted-pods/pkg/proto/protoconnect"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/libp2p/go-libp2p/core/peer"
)

func main() {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-interruptChan
		os.Exit(0)
	}()

	if os.Args[1] == "provider" {
		err := mainProvider()
		if err != nil {
			fmt.Printf("provider: error %v\n", err)
			os.Exit(1)
		}
	} else if os.Args[1] == "publisher" {
		err := mainPublisher()
		if err != nil {
			fmt.Printf("publisher: error %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Unknown role: %s\n", os.Args[1])
		os.Exit(1)
	}
}

type server struct {
	pbcon.UnimplementedProvisionPodServiceHandler
	node *rpc.HttpApi
}

func (s *server) ProvisionPod(ctx context.Context, in *connect.Request[pb.ProvisionPodRequest]) (*connect.Response[pb.ProvisionPodResponse], error) {
	return connect.NewResponse(&pb.ProvisionPodResponse{Error: fmt.Sprint(in.Msg.Pod.Replicas.Max)}), nil
}

func mainProvider() error {
	node, multiaddr, err := tpipfs.GetIpfsClient("")
	if err != nil {
		return err
	}

	lis, err := tpipfs.NewP2pApi(node, multiaddr).Listen(pb.ProvisionPod)
	if err != nil {
		return err
	}
	defer lis.Close()
	fmt.Printf("provider: server listening at %v\n", lis.Addr())

	mux := http.NewServeMux()

	mux.Handle(pbcon.NewProvisionPodServiceHandler(&server{node: node}))

	return (&http.Server{Handler: mux}).Serve(lis)
}

func mainPublisher() error {
	node, multiaddr, err := tpipfs.GetIpfsClient("")
	if err != nil {
		return err
	}

	providerPeerId, err := peer.Decode(os.Args[2])
	if err != nil {
		return err
	}

	num := rand.Uint32()

	request := &pb.ProvisionPodRequest{
		Pod: &pb.Pod{Replicas: &pb.Replicas{Max: num}},
	}

	addr, err := tpipfs.NewP2pApi(node, multiaddr).ConnectTo(pb.ProvisionPod, providerPeerId)
	if err != nil {
		return err
	}
	defer addr.Close()

	// Set up a connection to the server.
	c := pbcon.NewProvisionPodServiceClient(http.DefaultClient, (&url.URL{Scheme: "http", Host: addr.String()}).String())

	response, err := c.ProvisionPod(context.Background(), connect.NewRequest(request))
	if err != nil {
		return err
	}
	if fmt.Sprint(num) != response.Msg.Error {
		return fmt.Errorf("publisher: TEST FAILED, expected %v, received %v", fmt.Sprint(num), response.Msg.Error)
	}
	fmt.Printf("publisher: test passed successfully\n")
	return nil
}
