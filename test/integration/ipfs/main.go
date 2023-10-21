package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"

	"github.com/comrade-coop/trusted-pods/pkg/ipfs"
	tpipfs "github.com/comrade-coop/trusted-pods/pkg/ipfs"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/libp2p/go-libp2p/core/peer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	pb.UnimplementedProvisionPodServiceServer
	node *rpc.HttpApi
}

func (s *server) ProvisionPod(ctx context.Context, in *pb.ProvisionPodRequest) (*pb.ProvisionPodResponse, error) {
	pod := &pb.Pod{}
	tpipfs.GetProtobufFile(s.node, cid.MustParse(in.PodManifestCid), pod)

	return &pb.ProvisionPodResponse{Error: fmt.Sprint(pod.Replicas.Max)}, nil
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

	s := grpc.NewServer()
	pb.RegisterProvisionPodServiceServer(s, &server{node: node})

	return s.Serve(lis)
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

	cid, err := ipfs.AddProtobufFile(node, &pb.Pod{Replicas: &pb.Replicas{Max: num}})
	if err != nil {
		return err
	}

	request := &pb.ProvisionPodRequest{
		PodManifestCid: cid.Bytes(),
	}

	addr, err := tpipfs.NewP2pApi(node, multiaddr).ConnectTo(pb.ProvisionPod, providerPeerId)
	if err != nil {
		return err
	}
	defer addr.Close()

	// Set up a connection to the server.
	conn, err := grpc.Dial(addr.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	c := pb.NewProvisionPodServiceClient(conn)

	response, err := c.ProvisionPod(context.Background(), request)
	if err != nil {
		return err
	}
	if fmt.Sprint(num) != response.Error {
		return fmt.Errorf("publisher: TEST FAILED, expected %v, received %v", fmt.Sprint(num), response.Error)
	}
	fmt.Printf("publisher: test passed successfully\n")
	return nil
}
