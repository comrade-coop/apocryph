package main

import (
	"github.com/comrade-coop/trusted-pods/pkg/crypto"
	ipfs_utils "github.com/comrade-coop/trusted-pods/pkg/ipfs-utils"
	podmanagement "github.com/comrade-coop/trusted-pods/pkg/pod-management"
)

func main() {
	// create a pod named mypod
	p := podmanagement.CreatePod("mypod")
	// assign the pod manifest
	p.AssignManifest("./manifest.yaml")
	// choose an upload option
	provider, _ := podmanagement.CreateIpfsUploader()
	// create a key service for encrypting the package
	ks := crypto.MockKeyService{}
	// upload the package to ipfs
	cid := p.UploadPackage(provider, ks)
	if cid != "" {
		// retreive the package from ipfs
		ipfs_utils.RetreiveFile(provider.Node, cid, "/tmp/package/")
	} else {
		println("failed uploading pod package manifest")
	}
}
