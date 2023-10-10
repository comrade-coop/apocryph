// Package podmanagement provides functionality for managing pods (create pods, upload package to ipfs).
package podmanagement

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/comrade-coop/trusted-pods/pkg/crypto"
	"github.com/spf13/viper"
)

type pod struct {
	Name            string
	manifest        string // could not be changed after execution
	instances       uint   // number of instances currently executing across diffrent providers
	cid             string // final package manifest cid
	PackageSavePath string
	key             crypto.KeyNoncePair
}

func CreatePod(name string, psw string) (pod, error) {
	// TODO read salt from config file
	salt := []byte("salt")
	noncesalt := []byte("nonce")
	key, err := crypto.CreateKeyNoncePair([]byte(psw), salt, noncesalt)
	if err != nil {
		return pod{}, err
	}
	return pod{
		Name:      name,
		manifest:  "",
		instances: 0,
		cid:       "",
		key:       *key,
	}, nil
}

func (p *pod) AssignManifest(path string) error {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		log.Println("manifest file does not exist")
		return err
	}
	p.manifest = path
	return nil
}

func (p *pod) GetKeyNoncePair(psw string) crypto.KeyNoncePair {
	// TODO read salt from config file
	salt := []byte("salt")
	noncesalt := []byte("nonce")
	return crypto.KeyNoncePair{
		Key:   crypto.DeriveKey([]byte(psw), salt, crypto.AES_KEY_SIZE),
		Nonce: crypto.DeriveKey([]byte(psw), noncesalt, crypto.NONCE_SIZE),
	}
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
		log.Println("No image names found in the manifest file.")
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

func (p *pod) EncryptPodPackage(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		// Check if the entry is a regular file (not a directory)
		if file.Type().IsRegular() {
			// Get the file name
			name := file.Name()
			log.Println("Encrypting file:", name)
			filepath := fmt.Sprintf("%v/%v", path, name)
			// Read the file contents as a byte slice
			data, err := os.ReadFile(filepath)
			if err != nil {
				return err
			}
			cipheredtext, _, err := crypto.AESEncryptWith(data, p.key.Key, p.key.Nonce)
			if err != nil {
				return err
			}
			cipheredfilepath := fmt.Sprintf("%v.enc", filepath)
			destination, err := os.Create(cipheredfilepath)
			if err != nil {
				return err
			}
			_, err = destination.Write(cipheredtext)
			if err != nil {
				return err
			}
			os.Remove(filepath)
		}
	}
	return nil
}

func DecryptPodPackage(path string, key crypto.KeyNoncePair) error {

	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		// Check if the entry is a regular file (not a directory)
		// this could be concurrent
		if file.Type().IsRegular() {
			name := file.Name()
			log.Println("Decrypting file:", name)
			filepath := fmt.Sprintf("%v/%v", path, name)
			data, err := os.ReadFile(filepath)
			if err != nil {
				return err
			}

			data, err = crypto.AESDecryptWith(data, key.Key, key.Nonce)

			originalfilepath := strings.Replace(filepath, ".enc", "", 1)
			destination, err := os.Create(originalfilepath)
			if err != nil {
				return err
			}
			_, err = destination.Write(data)
			if err != nil {
				return err
			}
			os.Remove(filepath)

		}
	}
	return nil
}

func (p *pod) UploadPackage(provider PackageUploader, packagePath ...string) (string, error) {
	if len(packagePath) > 0 {
		return provider.uploadPackage(p, packagePath[0])
	}
	return provider.uploadPackage(p)
}
