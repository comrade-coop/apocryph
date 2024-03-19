// SPDX-License-Identifier: GPL-3.0

package publisher

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	pb "github.com/comrade-coop/apocryph/pkg/proto"
	"github.com/mitchellh/go-homedir"
)

var DefaultDeploymentPath = "~/.apocryph/deployment"
var DefaultPodFile = "manifest.yaml"

func GenerateDeploymentFilename(podFile string, deploymentFormat string) (deploymentFile string, relPodFile string, err error) {
	if deploymentFormat == "" {
		deploymentFormat = "yaml"
	}

	absPodFile, err := filepath.Abs(podFile)
	if err != nil {
		return
	}

	podFileHash := sha256.Sum256([]byte(absPodFile))
	podFileHashHex := hex.EncodeToString(podFileHash[:])
	deploymentFilename := fmt.Sprintf("%s.%s", podFileHashHex, deploymentFormat)
	deploymentRoot, err := homedir.Expand(DefaultDeploymentPath)
	if err != nil {
		return
	}

	err = os.MkdirAll(deploymentRoot, 0755)
	if err != nil {
		return
	}

	deploymentFile = filepath.Join(deploymentRoot, deploymentFilename)

	absDeploymentFile, err := filepath.Abs(deploymentFile)
	if err != nil {
		return
	}
	relPodFile, err = filepath.Rel(filepath.Dir(absDeploymentFile), absPodFile)

	return
}

func ReadPodAndDeployment(args []string, manifestFormat string, deploymentFormat string) (podFile string, deploymentFile string, pod *pb.Pod, deployment *pb.Deployment, err error) {
	deployment = &pb.Deployment{}
	readDeployment := false

	switch len(args) {
	case 0:
		podFile = DefaultPodFile
	case 1:
		podFile = args[0]
	case 2:
		podFile = args[0]
		deploymentFile = args[1]
	default:
		err = fmt.Errorf("Wrong number of arguments passed to ReadPodAndDeployment")
		return
	}

	// Get the name of the deployment file if it was not passed in the args
	if deploymentFile == "" {
		deploymentFile, deployment.PodManifestFile, err = GenerateDeploymentFilename(podFile, deploymentFormat)
		if err != nil {
			err = fmt.Errorf("Failed resolving deployment file path: %w", err)
			return
		}
	}

	if !readDeployment {
		err = pb.UnmarshalFile(deploymentFile, deploymentFormat, deployment)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			err = fmt.Errorf("Failed reading deployment file %s: %w", deploymentFile, err)
			return
		}
	}

	pod = &pb.Pod{}
	err = pb.UnmarshalFile(podFile, manifestFormat, pod)
	if err != nil {
		err = fmt.Errorf("Failed reading manifest file %s: %w", podFile, err)
	}

	return
}

func SaveDeployment(deploymentFile string, deploymentFormat string, deployment *pb.Deployment) error {
	err := pb.MarshalFile(deploymentFile, deploymentFormat, deployment)
	if err != nil {
		return fmt.Errorf("Failed saving deployment file: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Stored deployment data in %s\n", deploymentFile)
	// if deployment == nil {
	// 	err := os.Remove(deploymentFile)
	// 	if err != nil && !errors.Is(err, os.ErrNotExist) {
	// 		return fmt.Errorf("Failed to remove deployment file: %w", err)
	// 	}
	// 	fmt.Fprintf(os.Stderr, "Removed deployment data from %s\n", deploymentFile)
	// }
	return nil
}
