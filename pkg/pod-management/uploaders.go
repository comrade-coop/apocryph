package podmanagement

import (
	"fmt"
	"log"
	"os/exec"

	ipfs_utils "github.com/comrade-coop/trusted-pods/pkg/ipfs-utils"
	"github.com/ipfs/kubo/client/rpc"
)

type PackageUploader interface {
	// uploadPackage uploads a package associated with a pod to IPFS and returns its CID
	uploadPackage(p *pod, path ...string) (string, error)
	// uploadImages(p *pod, ks crypto.KeyService) []string
	// uploadManifest(p *pod, ks crypto.KeyService) []string
}

type IpfsUploader struct {
	Node *rpc.HttpApi
}

func CreateIpfsUploader() (*IpfsUploader, error) {
	node, err := ipfs_utils.ConnectToLocalNode()
	if err != nil {
		log.Println("could not connect to local IPFS node")
		return nil, err
	}
	return &IpfsUploader{Node: node}, nil
}

// Uploads a package associated with a pod to IPFS
// The package includes a manifest file and Docker images.
// It performs the following steps:
// 1. Copies the manifest file to the package directory.
// 2. Saves Docker images locally in the package directory.
// 3. TODO Uses the key service (ks) to encrypt the package files
// 4. Uploads the package directory to IPFS using the specified IPFS node
func (provider IpfsUploader) uploadPackage(p *pod, path ...string) (string, error) {
	var packagePath string
	if len(path) > 0 && path[0] != "" {
		packagePath = path[0]
	} else {
		packagePath = fmt.Sprintf("./%v/", p.Name)
	}

	ensureFolder(packagePath)

	cmd := exec.Command("cp", p.manifest, packagePath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Println("could not copy manifest file to package path:", err)
		log.Println("Command output:", string(output))
		return "", err
	}

	err = p.SaveImagesLocally(packagePath)
	if err != nil {
		return "", err
	}

	err = p.EncryptPodPackage(packagePath)
	if err != nil {
		fmt.Printf("could not encrypt pod package %v", err)
		return "", err
	}

	cid, err := ipfs_utils.AddFile(provider.Node, packagePath)
	if err != nil {
		fmt.Printf("could not add package to IPFS%v", err)
		return "", err
	}
	p.cid = cid
	return cid, nil
}

// TODO uses imgcrypt + ipdr to upload the images
type IpdrUploader struct{}

func (provider *IpdrUploader) uploadPackage(p *pod) []string {
	// imgaes := p.extractImageNames()
	// use the key service to encrypt the images using imgcrypt
	// upload the encrypted images to ipdr
	return []string{"QfGhXqshaysSU"}
}
