package main

import (
	"context"
	"log"
	"os"

	"github.com/comrade-coop/apocryph/pkg/ipcr"
	"github.com/containerd/containerd"
	"github.com/containerd/nerdctl/pkg/api/types"
	img "github.com/containerd/nerdctl/pkg/cmd/image"
	"github.com/ipfs/kubo/client/rpc"
)

var ipfs *rpc.HttpApi

const (
	IPFS_ADDRESS = "/ip4/127.0.0.1/tcp/5001"
	IMAGE_NAME   = "hello-world"
	PASSWORD     = "dummy_password"
)

func main() {
	// Initialize containerd client
	client, err := containerd.New("/run/containerd/containerd.sock", containerd.WithDefaultNamespace("k8s.io"))
	if err != nil {
		panic(err)
	}
	defer client.Close()
	log.Println("Encrypting Image ...")
	pubKey, prvKey, err := ipcr.EncryptImage(context.Background(), client, IMAGE_NAME, PASSWORD)
	if err != nil {
		panic(err)
	}
	log.Println("Image Encrypted")
	// removing the original image
	err = img.Remove(context.Background(), client, []string{IMAGE_NAME}, types.ImageRemoveOptions{Stdout: os.Stdout})
	if err != nil {
		panic(err)
	}

	printImages(client)
	// push the encrypted image to ipfs
	cid, err := ipcr.PushImage(context.Background(), client, IPFS_ADDRESS, IMAGE_NAME+":encrypted")
	if err != nil {
		log.Panic(err)
	}

	// removing encrypted image
	err = img.Remove(context.Background(), client, []string{IMAGE_NAME + ":encrypted"}, types.ImageRemoveOptions{Stdout: os.Stdout})
	if err != nil {
		log.Panic(err)
	}
	printImages(client)

	// pulling the ecnrypted image from ipfs
	log.Println("Pulling Encrypted Image")
	err = ipcr.PullImage(context.Background(), client, IPFS_ADDRESS, cid, IMAGE_NAME+":encrypted")
	if err != nil {
		log.Panic(err)
	}
	printImages(client)

	// decrypting pulled Image
	err = ipcr.DecryptImage(context.Background(), client, PASSWORD, IMAGE_NAME+":encrypted", pubKey, prvKey)
	if err != nil {
		log.Panic(err)
	}
	printImages(client)

}

func printImages(client *containerd.Client) {
	images, _ := img.List(context.Background(), client, nil, nil)
	log.Println("Current List of images:")
	for idx, img := range images {
		log.Printf("idx:%v %+v\n", idx, img.Name)
	}
}
