package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	ipfs_utils "github.com/comrade-coop/trusted-pods/pkg/ipfs-utils"
	podmanagement "github.com/comrade-coop/trusted-pods/pkg/pod-management"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 6000, "The server port")
)

type server struct {
	pb.UnimplementedSampleServer
}

func (s *server) SendPod(ctx context.Context, in *pb.SampleProvisionPodRequest) (*pb.SampleProvisionPodReply, error) {
	println("PROVIDER: pod package cid:", in.Cid)
	// cleanup unpinned files
	exec.Command("ipfs", "repo", "gc")
	provider, _ := podmanagement.CreateIpfsUploader()
	ipfs_utils.RetreiveFile(provider.Node, in.Cid, "/tmp/pod-package")
	cmd := exec.Command("ls", "-al", "/tmp/pod-package")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("PROVIDER: could not list directory files:", err)
		fmt.Println("PROVIDER: Command output:", string(output))
		return nil, err
	}
	fmt.Println("PROVIDER: Pod Package:", string(output))
	return &pb.SampleProvisionPodReply{Endpoint: "http://provider.com/mypod"}, nil
}

func main() {
	// Create a channel to receive the interrupt signal
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-interruptChan
		os.Exit(0)
	}()

	node, err := ipfs_utils.ConnectToLocalNode()
	if err != nil {
		println("PROVIDER: could not connect to local node")
		return
	}

	// route all ipfs p2p connections of the provios-pod protocol to the grpc server
	err = ipfs_utils.CreateIpfsService(node, pb.ProvisionPod, "/ip4/127.0.0.1/tcp/6000")
	if err != nil {
		println("Could not create IPFS service")
		return
	}
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("PROVIDER: failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSampleServer(s, &server{})
	log.Printf("PROVIDER: server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("PROVIDER: failed to serve: %v", err)
	}
}
