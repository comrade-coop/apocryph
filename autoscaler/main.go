package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/comrade-coop/apocryph/pkg/autoscaler"
	pbcon "github.com/comrade-coop/apocryph/pkg/proto/protoconnect"
	tpraft "github.com/comrade-coop/apocryph/pkg/raft"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	mux := http.NewServeMux()
	self, err := tpraft.NewPeer(fmt.Sprintf("/ip4/0.0.0.0/tcp/%v", autoscaler.RAFT_P2P_PORT))
	if err != nil {
		fmt.Println("Failed creating p2p node")
		return
	}

	fmt.Printf("PEER ID: %v\n", self.ID())

	server, err := autoscaler.NewAutoSalerServer("http://eth-rpc.eth.svc.cluster.local:8545")
	if err != nil {
		fmt.Println("Failed creating AutoScaler Server")
		return
	}

	server.MainLoop = autoscaler.SetAppGatewayExample
	path, handler := pbcon.NewAutoscalerServiceHandler(server)
	mux.Handle(path, handler)
	log.Println("Autoscaler RPC Server Started")
	http.ListenAndServe(
		":8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)

}
