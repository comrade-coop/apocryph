package main

import (
	"log"
	"os"

	ipfs_utils "github.com/comrade-coop/trusted-pods/pkg/ipfs-utils"
	podmanagement "github.com/comrade-coop/trusted-pods/pkg/pod-management"
)

func main() {
	// create a pod named mypod
	p, _ := podmanagement.CreatePod("mypod", "mypassword")
	// assign the pod manifest
	p.AssignManifest("./package-files/manifest.yaml")
	// choose an upload option
	provider, _ := podmanagement.CreateIpfsUploader()
	// create a key service for encrypting the package
	// upload the package to ipfs
	cid, _ := p.UploadPackage(provider, os.Args[1])
	if cid != "" {
		// retreive the package from ipfs
		ipfs_utils.RetreiveFile(provider.Node, cid, "/tmp/package/")
	} else {
		println("failed uploading pod package manifest")
	}
	err := podmanagement.DecryptPodPackage("/tmp/package/", p.GetKeyNoncePair("mypassword"))
	if err != nil {
		log.Fatalf("%v", err)
	}
}
