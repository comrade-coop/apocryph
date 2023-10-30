package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	tpk8s "github.com/comrade-coop/trusted-pods/pkg/kubernetes"
	"github.com/comrade-coop/trusted-pods/pkg/provider"
	"github.com/ipfs/kubo/client/rpc"
)

func main() {
	if len(os.Args) < 2 {
		println("Usage: server <listening-Address>")
		return
	}
	serverAddress := os.Args[1]
	ipfsApi, err := rpc.NewLocalApi()
	if err != nil {
		log.Printf("Failed to Connect to local ipfs node: %v", err)
		return
	}

	k8cl, err := tpk8s.GetClient("", false)
	if err != nil {
		log.Printf("Failed to configure k8s client: %v", err)
		return
	}
	// skip kubeConfig
	s, err := provider.NewTPodServer(ipfsApi, false, k8cl, "", nil)
	if err != nil {
		log.Printf("Failed to create grpc server: %v", err)
		return
	}
	// listen on regular address instead of p2p
	listener, err := provider.GetListener(serverAddress)
	if err != nil {
		log.Printf("Failed to get Listening Address: %v", err)
		return
	}

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-interruptChan
		s.GracefulStop()
		os.Exit(0)
	}()

	log.Printf("tpodserver: server listening at %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("PROVIDER: failed to serve: %v", err)
	}

}
