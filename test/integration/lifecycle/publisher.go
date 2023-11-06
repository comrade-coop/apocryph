package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/comrade-coop/trusted-pods/pkg/publisher"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ipfs/kubo/client/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ipfs *rpc.HttpApi

func main() {
	if len(os.Args) < 2 {
		println("Usage: server <Server-Address>")
		return
	}
	serverAddress := os.Args[1]
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Could not dial server: %v", serverAddress)
		return
	}

	defer conn.Close()

	ks := keystore.NewKeyStore("./keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := ks.NewAccount("123")
	if err != nil {
		fmt.Printf("could not create account %v", err)
		return
	}
	publisherAddress := []byte(acc.Address.Hex())

	signature, err := pb.SignPayload(publisherAddress, acc, "123", ks)
	if err != nil {
		fmt.Printf("could not sign message: %v", err)
	}
	Credentials := &pb.Credentials{PublisherAddress: publisherAddress, Signature: signature}
	ipfs, err = rpc.NewLocalApi()
	if err != nil {
		log.Printf("Failed to Connect to local ipfs node: %v", err)
		return
	}

	client := pb.NewProvisionPodServiceClient(conn)
	ProvisionPod(client, publisherAddress, "./manifest-guestbook.json")
	GetPodLogs(client, &pb.PodLogRequest{ContainerName: "php-redis", Credentials: Credentials})
	UpdatePod(client, Credentials)
	log.Println("Press Enter to Delete Namespace")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')

	DeletePod(client, &pb.DeletePodRequest{Credentials: Credentials})
}

func ProvisionPod(client pb.ProvisionPodServiceClient, publisherAddress []byte, podPath string) {

	request, _, err := publisher.UploadManifest(context.Background(), "./manifest-guestbook.json", "json", "", false)
	if err != nil {
		log.Fatalf("failed uploading Manifest: %v", err)
	}
	request.Payment = &pb.PaymentChannel{PublisherAddress: publisherAddress}
	response, err := client.ProvisionPod(context.Background(), request)
	if err != nil {
		log.Fatalf("Provision Pod failed: %v", err)
	}
	log.Printf("pod provision response: %v", response)
}

func DeletePod(client pb.ProvisionPodServiceClient, request *pb.DeletePodRequest) {
	response, err := client.DeletePod(context.Background(), request)
	if err != nil {
		log.Printf("rpc delete method failed: %v", err)
		return
	}
	log.Printf("Pod Deletion response: %v", response)
}

func UpdatePod(client pb.ProvisionPodServiceClient, credentials *pb.Credentials) {

	request, pod, err := publisher.UploadManifest(context.Background(), "./updated-guestbook.json", "json", "", false)
	response, err := client.UpdatePod(context.Background(), &pb.UpdatePodRequest{Pod: pod, Credentials: credentials, Keys: request.Keys})
	if err != nil {
		log.Printf("rpc update method failed: %v", err)
		return
	}
	log.Printf("Pod Update response: %v", response)

}

func GetPodLogs(client pb.ProvisionPodServiceClient, request *pb.PodLogRequest) {
	stream, err := client.GetPodLogs(context.Background(), request)
	if err != nil {
		log.Printf("Could not get logs stream: %v", err)
		return
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Println("Log Stream Ended")
			return
		} else if err == nil {
			valStr := fmt.Sprintf("%s:%s", resp.LogEntry.Time, resp.LogEntry.Log)
			log.Println(valStr)
		}

		if err != nil {
			log.Printf("Failed reading Stream: %v", err)
			return
		}

	}
}
