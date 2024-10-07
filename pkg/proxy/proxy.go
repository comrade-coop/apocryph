package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"connectrpc.com/connect"
	"github.com/comrade-coop/apocryph/pkg/constants"
	pb "github.com/comrade-coop/apocryph/pkg/proto"
	pbcon "github.com/comrade-coop/apocryph/pkg/proto/protoconnect"
)

// TargetServerURL is the URL of the server to forward the request to
const TargetServerURL = "tpodserver.trustedpods.svc.cluster.local:8080"

func main() {
	http.HandleFunc("/", handler)
	log.Printf("Server is listening on port 9999")
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Get the environment variable value
	namespace := os.Getenv(constants.NamespaceKey)
	if namespace == "" {
		http.Error(w, "Environment variable not set", http.StatusInternalServerError)
		return
	}

	r.Header.Set("X-Namespace", namespace)
	client := pbcon.NewProvisionPodServiceClient(
		http.DefaultClient,
		(&url.URL{Scheme: "http", Host: TargetServerURL}).String())
	// Forward the request to the target server
	request := pb.PodInfoRequest{Namespace: namespace}
	info, err := forwardRequest(client, &request)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error forwarding request: %v", err), http.StatusInternalServerError)
		return
	}

	// Write the info string to the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(info))
}

func forwardRequest(client pbcon.ProvisionPodServiceClient, request *pb.PodInfoRequest) (string, error) {
	response, err := client.GetPodInfos(context.Background(), connect.NewRequest(request))
	if err != nil {
		return "", err
	}
	log.Printf("Retrived info: %v\n", response.Msg.Info)
	return response.Msg.Info, nil
}
