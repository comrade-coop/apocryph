// Package podmanagement provides functionality for managing pods (create pods, upload package to ipfs).
package podmanagement

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/comrade-coop/trusted-pods/crypto"
	ipfs_utils "github.com/comrade-coop/trusted-pods/ipfs-utils"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/spf13/viper"
)

type pod struct {
	Name          string
	manifest      string // could not be changed after execution
	biddingOpened bool   // open for bids or not
	instances     uint   // number of instances currently executing across diffrent providers
	cid           string // final package manifest cid
}

func CreatePod(name string) pod {
	return pod{
		Name:          name,
		manifest:      "",
		biddingOpened: false,
		instances:     0,
		cid:           "",
	}
}

func (p *pod) AssignManifest(path string) error {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("manifest file does not exist")
		return err
	}
	p.manifest = path
	return nil
}

// TODO
func sendPodRequest() {}

// Extracts Docker image names from a pod's configuration manifest.
// It reads the manifest file specified in the 'pod' instance and parses it to find image names.
// The manifest file is expected to be in YAML format.
func (p *pod) ExtractImageNames() []string {
	vp := viper.New()
	vp.SetConfigFile(p.manifest)
	err := vp.ReadInConfig()
	if err != nil {
		fmt.Printf("Could not read manifest file %v\n", err)
		return nil
	}
	containers := vp.Get("containers")
	yaml := fmt.Sprintf("%v", containers)

	pattern := `image:([^ \] ]+)`

	regex := regexp.MustCompile(pattern)

	matches := regex.FindAllStringSubmatch(yaml, -1)
	var images []string
	if len(matches) > 0 {
		for _, match := range matches {
			imageName := match[1]
			images = append(images, imageName)
		}
	} else {
		fmt.Println("No image names found in the manifest file.")
	}
	return images
}

// SaveImagesLocally saves Docker images associated with a pod locally to the specified directory.
// It takes a path parameter as the target directory where the images will be saved.
// The function first extracts image names associated with the pod using the ExtractImageNames method.
// Then, it creates the target directory if it doesn't already exist.
// For each image, it generates a filename based on the image name, replacing slashes with underscores,
// and saves the Docker image as a .tar file in the target directory.
// this method requires the images been available to be pulled locally
func (p *pod) SaveImagesLocally(path string) error {
	images := p.ExtractImageNames()
	ensureFolder(path)
	for _, v := range images {
		name := fmt.Sprintf("%v.tar", v)
		fileName := strings.Replace(name, "/", "_", -1)
		tarPath := fmt.Sprintf("%v/%v", path, fileName)
		args := []string{"save", "-o", tarPath, v}
		cmd := exec.Command("docker", args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("could not save image:", v, err)
			fmt.Println("Command output:", string(output))
			return err
		}
	}
	return nil
}

// Create the path folder if it doesn't exist
func ensureFolder(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (p *pod) EncryptSavedImages(path string) {}

func (p *pod) UploadPackage(provider PackageUploader, ks crypto.KeyService) string {
	return provider.uploadPackage(p, ks)
}

type PackageUploader interface {
	// uploadPackage uploads a package associated with a pod to IPFS and returns its CID
	uploadPackage(p *pod, ks crypto.KeyService, path ...string) string
	// uploadImages(p *pod, ks crypto.KeyService) []string
	// uploadManifest(p *pod, ks crypto.KeyService) []string
}

type IpfsUploader struct {
	Node *rpc.HttpApi
}

func CreateIpfsUploader() (*IpfsUploader, error) {
	node, err := ipfs_utils.ConnectToLocalNode()
	if err != nil {
		println("could not connect to local IPFS node")
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
func (provider IpfsUploader) uploadPackage(p *pod, ks crypto.KeyService, path ...string) string {
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
		fmt.Println("could not copy manifest file to package path:", err)
		fmt.Println("Command output:", string(output))
		return ""
	}
	err = p.SaveImagesLocally(packagePath)
	if err != nil {
		return ""
	}
	cid, err := ipfs_utils.AddFile(provider.Node, packagePath)
	if err != nil {
		fmt.Printf("could not add package to IPFS%v", err)
		return ""
	}
	p.cid = cid
	return cid
}

// TODO uses imgcrypt + ipdr to upload the images
type IpdrUploader struct{}

func (provider *IpdrUploader) uploadPackage(p *pod, ks crypto.KeyService) []string {
	// imgaes := p.extractImageNames()
	// use the key service to encrypt the images using imgcrypt
	// upload the encrypted images to ipdr
	return []string{"QfGhXqshaysSU"}
}
