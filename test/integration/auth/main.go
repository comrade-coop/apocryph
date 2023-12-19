// SPDX-License-Identifier: GPL-3.0

package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"connectrpc.com/connect"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	pbcon "github.com/comrade-coop/trusted-pods/pkg/proto/protoconnect"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	key, err := crypto.GenerateKey()
	if err != nil {
		fmt.Println(err)
		return
	}
	ethAddr := crypto.PubkeyToAddress(key.PublicKey)
	sign := func(data []byte) ([]byte, error) { // Copied from "github.com/comrade-coop/trusted-pods/pkg/ethereum".GetAccountAndSigner
		hash := crypto.Keccak256(data)
		sig, err := crypto.Sign(hash, key)
		if err != nil {
			return nil, err
		}
		return sig, nil
	}
	
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-interruptChan
		os.Exit(0)
	}()
	
	lisChan := make(chan net.Listener, 1)
	
	go func() {
		lis, err := net.Listen("tcp", "localhost:0")
		if err != nil {
			fmt.Println(err)
			close(lisChan)
			return
		}

		lisChan<-lis

		mux := http.NewServeMux()
		mux.Handle(pbcon.NewProvisionPodServiceHandler(&server{}, connect.WithInterceptors(
			pbcon.NewAuthInterceptor(nil),
		)))

		err = (&http.Server{Handler: mux}).Serve(lis)
		if err != nil {
			fmt.Println(err)
			lis.Close()
			return
		}
	}()

	addr := (<-lisChan).Addr()
	
	request := &pb.ProvisionPodRequest{
		Payment: &pb.PaymentChannel{
			PublisherAddress: ethAddr.Bytes(),
		},
	}
	
	client := pbcon.NewProvisionPodServiceClient(
		http.DefaultClient,
		(&url.URL{Scheme: "http", Host: addr.String()}).String(),
		connect.WithInterceptors(pbcon.NewAuthInterceptorClient(
			&pb.Deployment{
				Payment: &pb.PaymentChannelConfig{
					PodID: []byte{0, 1, 2, 3},
					PublisherAddress: ethAddr.Bytes(),
				},
			},
			10,
			sign,
		)),
	)
	
	response, err := client.ProvisionPod(context.Background(), connect.NewRequest(request))
	if err != nil {
		fmt.Println(err)
		return
	}
	if response.Msg.Error == "" {
		fmt.Printf("TEST FAILED, expected != '', received '%v'", response.Msg.Error)
		return
	}
	if strings.Count(response.Msg.Error, "-") != 1 {
		fmt.Printf("TEST FAILED, expected only one '-' in string, '%v'", response.Msg.Error)
		return
	}
	fmt.Printf("Test passed successfully (%v)\n", response.Msg.Error)
}

type server struct {
	pbcon.UnimplementedProvisionPodServiceHandler
}

func (s *server) ProvisionPod(ctx context.Context, in *connect.Request[pb.ProvisionPodRequest]) (*connect.Response[pb.ProvisionPodResponse], error) {
	return connect.NewResponse(&pb.ProvisionPodResponse{Error: pbcon.GetNamespace(in)}), nil
}
