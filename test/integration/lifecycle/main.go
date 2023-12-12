// SPDX-License-Identifier: GPL-3.0

package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"connectrpc.com/connect"
	"github.com/comrade-coop/trusted-pods/pkg/ethereum"
	tpeth "github.com/comrade-coop/trusted-pods/pkg/ethereum"
	tpipfs "github.com/comrade-coop/trusted-pods/pkg/ipfs"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	pbcon "github.com/comrade-coop/trusted-pods/pkg/proto/protoconnect"
	"github.com/comrade-coop/trusted-pods/pkg/publisher"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ipfs/kubo/client/rpc"
)

var ipfs *rpc.HttpApi

func main() {
	if len(os.Args) < 2 {
		println("Usage: server <Minikube-IP> <GRPC-PORT> <IPFS-PORT>")
		return
	}
	minikubeIp := os.Args[1]
	grpcPort := os.Args[2]
	ipfsPort := os.Args[3]
	serverAddress := fmt.Sprintf("%v:%v", minikubeIp, grpcPort)
	ipfsAddress := fmt.Sprintf("/ip4/%v/tcp/%v", minikubeIp, ipfsPort)

	ethClient, err := ethereum.GetClient("http://127.0.0.1:8545")
	if err != nil {
		log.Fatalf("could not get eth client %v", err)
	}
	auth, sign, err := tpeth.GetAccountAndSigner("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", ethClient)
	if err != nil {
		log.Fatalf("failed creating Account %v", err)
	}
	publisherAddress := auth.From.Bytes()

	ipfs, _, err = tpipfs.GetIpfsClient(ipfsAddress)
	if err != nil {
		log.Fatalf("failed to retreived Ipfs Client: %v", err)
	}

	deployment := &pb.Deployment{Payment: &pb.PaymentChannelConfig{}}
	deployment.Payment.PodID = common.BytesToHash([]byte("123456")).Bytes()
	deployment.Payment.PublisherAddress = publisherAddress

	interceptor := pbcon.NewAuthInterceptorClient(deployment, pbcon.ProvisionPodServiceProvisionPodProcedure, 10, sign)

	client := pbcon.NewProvisionPodServiceClient(
		http.DefaultClient,
		serverAddress,
		connect.WithInterceptors(interceptor),
	)

	reader := bufio.NewReader(os.Stdin)

	// Create Pod
	containerName := ProvisionPod(client, publisherAddress, "./manifest-guestbook.json")

	// Get Pod logs
	go GetPodLogs(client, &pb.PodLogRequest{ContainerName: containerName, PublisherAddress: publisherAddress})
	// Update Pod
	log.Printf("Press Enter to Update Pod \n\n")
	_, _ = reader.ReadString('\n')
	// containerName = UpdatePod(client, Credentials)
	UpdatePod(client, publisherAddress)

	// retreived logs for the updated container
	go GetPodLogs(client, &pb.PodLogRequest{ContainerName: containerName, PublisherAddress: publisherAddress})

	// Delete Pod
	log.Printf("Press Enter to Delete Namespace\n\n")
	_, _ = reader.ReadString('\n')
	DeletePod(client, &pb.DeletePodRequest{PublisherAddress: publisherAddress})
}

func ProvisionPod(client pbcon.ProvisionPodServiceClient, publisherAddress []byte, podPath string) string {

	podFile, _, pod, deployment, err := publisher.ReadPodAndDeployment([]string{"./logger-manifest.json"}, "", "")

	err = publisher.UploadSecrets(context.Background(), ipfs, filepath.Dir(podFile), pod, deployment)
	if err != nil {
		log.Fatalf("failed uploading Manifest: %v", err)
	}
	err = publisher.UploadImages(context.Background(), ipfs, pod, deployment)
	if err != nil {
		log.Fatalf("failed uploading Manifest: %v", err)
	}
	request := &pb.ProvisionPodRequest{PublisherAddress: publisherAddress}
	request.Pod = publisher.LinkUploadsFromDeployment(pod, deployment)
	request.Payment = &pb.PaymentChannel{PublisherAddress: publisherAddress, PodID: []byte("123456")}

	response, err := client.ProvisionPod(context.Background(), connect.NewRequest(request))
	if err != nil {
		log.Fatalf("Provision Pod failed: %v", err)
	}
	log.Printf("pod provision response: %v", response)
	return response.Msg.Addresses[0].ContainerName
}

func DeletePod(client pbcon.ProvisionPodServiceClient, request *pb.DeletePodRequest) {
	response, err := client.DeletePod(context.Background(), connect.NewRequest(request))
	if err != nil {
		log.Printf("rpc delete method failed: %v", err)
		return
	}
	log.Printf("Pod Deletion response: %v", response)
}

func UpdatePod(client pbcon.ProvisionPodServiceClient, publisherAddress []byte) string {

	podFile, _, pod, deployment, err := publisher.ReadPodAndDeployment([]string{"./updated-logger.json"}, "", "")

	err = publisher.UploadSecrets(context.Background(), ipfs, filepath.Dir(podFile), pod, deployment)
	if err != nil {
		log.Fatalf("failed uploading Manifest: %v", err)
	}
	err = publisher.UploadImages(context.Background(), ipfs, pod, deployment)
	if err != nil {
		log.Fatalf("failed uploading Manifest: %v", err)
	}
	request := &pb.UpdatePodRequest{PublisherAddress: publisherAddress}
	request.Pod = publisher.LinkUploadsFromDeployment(pod, deployment)
	response, err := client.UpdatePod(context.Background(), connect.NewRequest(request))
	if err != nil {
		log.Printf("rpc update method failed: %v", err)
		return ""
	}
	log.Printf("Pod Update response: %v", response)
	return response.Msg.Addresses[0].ContainerName
}

func GetPodLogs(client pbcon.ProvisionPodServiceClient, request *pb.PodLogRequest) {
	stream, err := client.GetPodLogs(context.Background(), connect.NewRequest(request))
	if err != nil {
		log.Printf("Could not get logs stream: %v", err)
		return
	}
	for stream.Receive() {
		resp := stream.Msg()
		valStr := fmt.Sprintf("%s:%s", resp.LogEntry.Time, resp.LogEntry.Log)
		log.Println(valStr)
	}
	err = stream.Err()
	if err != nil {
		log.Printf("Failed reading Stream: %v", err)
		return
	}
	log.Println("Log Stream Ended")
	return
}
