package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/comrade-coop/trusted-pods/pkg/ethereum"
	tpeth "github.com/comrade-coop/trusted-pods/pkg/ethereum"
	tpipfs "github.com/comrade-coop/trusted-pods/pkg/ipfs"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/comrade-coop/trusted-pods/pkg/publisher"
	"github.com/ipfs/kubo/client/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	token := pb.NewToken("123456", pb.CreatePod, 10, publisherAddress)
	interceptor := &pb.AuthInterceptorClient{Token: token, Sign: sign, ExpirationOffset: 2 * time.Second}

	conn, err := grpc.Dial(
		serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptor.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(interceptor.StreamClientInterceptor()),
	)
	if err != nil {
		log.Fatalf("Could not dial server: %v", serverAddress)
	}

	defer conn.Close()

	client := pb.NewProvisionPodServiceClient(conn)
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

func ProvisionPod(client pb.ProvisionPodServiceClient, publisherAddress []byte, podPath string) string {

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

	response, err := client.ProvisionPod(context.Background(), request)
	if err != nil {
		log.Fatalf("Provision Pod failed: %v", err)
	}
	log.Printf("pod provision response: %v", response)
	return response.Addresses[0].ContainerName
}

func DeletePod(client pb.ProvisionPodServiceClient, request *pb.DeletePodRequest) {
	response, err := client.DeletePod(context.Background(), request)
	if err != nil {
		log.Printf("rpc delete method failed: %v", err)
		return
	}
	log.Printf("Pod Deletion response: %v", response)
}

func UpdatePod(client pb.ProvisionPodServiceClient, publisherAddress []byte) string {

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
	response, err := client.UpdatePod(context.Background(), request)
	if err != nil {
		log.Printf("rpc update method failed: %v", err)
		return ""
	}
	log.Printf("Pod Update response: %v", response)
	return response.Addresses[0].ContainerName
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
