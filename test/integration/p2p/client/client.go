package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	ipfs_utils "github.com/comrade-coop/trusted-pods/pkg/ipfs-utils"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:5000", "the address to connect to")
)

func main() {
	node, err := ipfs_utils.ConnectToLocalNode()
	if err != nil {
		println("CLIENT: could not connect to local node")
		return
	}
	cid, err := ipfs_utils.AddFile(node, "./client-pod")
	if err != nil {
		println("CLIENT: could not add pod package to ipfs")
		return
	}
	providerID := os.Getenv("PROVIDER_ID")

	println("CLIENT: pod package cid:", cid)
	cmd := exec.Command("ls", "-al", "client-pod")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("CLIENT: could not list directory files:", err)
		fmt.Println("CLIENT: Command output:", string(output))
		return
	}
	// Redirect all gRPC client requests to 127.0.0.1:5000 to IPFS,
	// which will in turn route these requests to the provider ID.
	// The provider ID will be listening on that protocol
	// And will further route all requests of the same protocol to his own gRPC server. running on a DIFFRENT port"
	err = ipfs_utils.ForwardConnection(node, pb.ProvisionPod, "/ip4/127.0.0.1/tcp/5000", providerID)
	if err != nil {
		println("Could not forward connection")
		return
	}
	fmt.Println("CLIENT: Pod Package:", string(output))
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("CLIENT: did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSampleClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SendPod(ctx, &pb.SampleProvisionPodRequest{Cid: cid})
	if err != nil {
		log.Fatalf("CLIENT: could not send pod: %v", err)
	}
	log.Printf("CLIENT: Pod Endpoint: %s", r.GetEndpoint())

}
