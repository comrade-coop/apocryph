package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path"

	tpipfs "github.com/comrade-coop/trusted-pods/pkg/ipfs"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
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
	_, pod, err := ProvisionPod(client, publisherAddress, "./manifest-guestbook.json")
	if err != nil {
		log.Printf("Could not provision pod: %v", err)
		return
	}

	pod.Containers[0].Ports[0].ContainerPort = 99
	pod.Containers[0].Name = "js-redis"
	request := &pb.UpdatePodRequest{Credentials: Credentials, Pod: pod}
	UpdatePod(client, request)

	// GetPodLogs(client, &pb.PodLogRequest{ContainerName: "anvil"}) // empty credentials should fail
	// GetPodLogs(client, &pb.PodLogRequest{ContainerName: response.Addresses[0].ContainerName, Credentials: Credentials}) // the pod quickly scaled down?

	log.Println("Press Enter to Delete Namespace")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')

	DeletePod(client, &pb.DeletePodRequest{Credentials: Credentials})
}

func ProvisionPod(client pb.ProvisionPodServiceClient, publisherAddress []byte, podPath string) (*pb.ProvisionPodResponse, *pb.Pod, error) {

	pod := &pb.Pod{}
	err := pb.UnmarshalFile(podPath, "json", pod)
	if err != nil {
		return nil, nil, err
	}

	keys := []*pb.Key{}

	err = tpipfs.TransformSecrets(pod,
		tpipfs.ReadSecrets(path.Dir(podPath)),
		tpipfs.EncryptSecrets(&keys),
		tpipfs.UploadSecrets(context.Background(), ipfs),
	)
	if err != nil {
		return nil, nil, err
	}

	podCid, err := tpipfs.AddProtobufFile(ipfs, pod)
	if err != nil {
		return nil, nil, err
	}

	request := &pb.ProvisionPodRequest{
		PodManifestCid: podCid.Bytes(),
		Payment:        &pb.PaymentChannel{PublisherAddress: publisherAddress},
		Keys:           keys,
	}

	response, err := client.ProvisionPod(context.Background(), request)
	if err != nil {
		return nil, nil, err
	}
	log.Printf("pod provision response: %v", response)

	return response, pod, nil
}

func DeletePod(client pb.ProvisionPodServiceClient, request *pb.DeletePodRequest) {
	response, err := client.DeletePod(context.Background(), request)
	if err != nil {
		log.Printf("rpc delete method failed: %v", err)
		return
	}
	log.Printf("Pod Deletion response: %v", response)
}

func UpdatePod(client pb.ProvisionPodServiceClient, request *pb.UpdatePodRequest) {
	response, err := client.UpdatePod(context.Background(), request)
	if err != nil {
		log.Printf("rpc update method failed: %v", err)
		return
	}
	log.Printf("Pod Update response: %v", response)

}

func GetPodLogs(client pb.ProvisionPodServiceClient, request *pb.PodLogRequest) {
	response, err := client.GetPodLogs(context.Background(), request)
	if err != nil {
		log.Printf("rpc log method failed: %v", err)
		return
	}
	log.Printf("Pod Log response: %v", response)

}
