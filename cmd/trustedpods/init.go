// SPDX-License-Identifier: GPL-3.0

package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	pb "github.com/comrade-coop/apocryph/pkg/proto"
	"github.com/comrade-coop/apocryph/pkg/publisher"
	"github.com/spf13/cobra"
)

var disallowedVolumeNameCharacters = regexp.MustCompile("[\\.-_/\\\\0-9\n ]+$|^[0-9\\.-_/\\\\\n ]+|[\\.-_/\\\\\n ]")
var volumeSizeFlag = regexp.MustCompile("^ *([0-9]+) *([KMG]i?B) *$")
var initImageName string
var initVolumes []string
var initPort string

const initDefaultHostname = "pod.XXXX.hostname.example"
const initDefaultPort = uint64(80)

var initPodCmd = &cobra.Command{
	Use:   "init [manifest.yaml]",
	Short: "Create a new pod manifest",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		manifestFile := publisher.DefaultPodFile
		if len(args) >= 1 {
			manifestFile = args[0]
		}
		manifestFileAbs, err := filepath.Abs(manifestFile)
		if err != nil {
			return err
		}

		imageParts := strings.Split(initImageName, "/")
		containerName, _, _ := strings.Cut(imageParts[len(imageParts)-1], ":")

		initPortParts := strings.SplitN(initPort, ":", 2)
		port := initDefaultPort
		hostname := initDefaultHostname
		if len(initPortParts) == 1 {
			potentialPort, err := strconv.ParseUint(initPortParts[0], 10, 64)
			if err == nil {
				port = potentialPort
			} else {
				hostname = initPortParts[0]
			}
		} else if len(initPortParts) == 2 {
			port, err = strconv.ParseUint(initPortParts[0], 10, 64)
			if err != nil {
				return fmt.Errorf("Invalid port specification %s: %w", initPort, err)
			}
			hostname = initPortParts[1]
		}
		hostname = strings.Replace(hostname, "XXXX", strconv.FormatInt(rand.Int63(), 16), 1)

		pod := &pb.Pod{
			Containers: []*pb.Container{
				{
					Name: containerName,
					Image: &pb.Container_Image{
						Url: initImageName,
					},
					Ports: []*pb.Container_Port{
						{
							Name:          "http",
							ContainerPort: port,
							ExposedPort: &pb.Container_Port_HostHttpHost{
								HostHttpHost: hostname,
							},
						},
					},
					ResourceRequests: []*pb.Resource{
						{Resource: "cpu", Quantity: &pb.Resource_AmountMillis{AmountMillis: 100}},
						{Resource: "memory", Quantity: &pb.Resource_Amount{Amount: 100000000}},
					},
				},
			},
			Replicas: &pb.Replicas{
				Min: 0,
				Max: 1,
			},
		}

		for _, volumeSpec := range initVolumes {
			volumeSpecParts := strings.Split(volumeSpec, ":")

			readOnly := false
			volumeSize := uint64(5 * 1000 * 1000 * 1000)
			secretPath := ""
			mountPath := ""

			if len(volumeSpecParts) > 0 {
				potentialSecret := volumeSpecParts[0]
				info, err := os.Stat(potentialSecret)
				if err == nil && !info.IsDir() {
					secretPath = potentialSecret
					volumeSpecParts = volumeSpecParts[1:]
				}
			}

			if len(volumeSpecParts) > 0 {
				potentialFlags := volumeSpecParts[len(volumeSpecParts)-1]
				if !strings.Contains(potentialFlags, "/") {
					potentialFlagsParts := strings.Split(potentialFlags, ",")
					for _, flag := range potentialFlagsParts {
						if flag == "ro" {
							readOnly = true
						} else if flag == "rw" {
							readOnly = false
						} else if matches := volumeSizeFlag.FindStringSubmatch(flag); matches != nil {
							amountBytes, err := strconv.ParseInt(matches[0], 10, 64)
							if err != nil {
								return fmt.Errorf("Invalid volume size: %w", err)
							}
							switch matches[1] {
							case "KB":
								amountBytes *= 1024
							case "MB":
								amountBytes *= 1024 * 1024
							case "GB":
								amountBytes *= 1024 * 1024 * 1024
							case "KiB":
								amountBytes *= 1000
							case "MiB":
								amountBytes *= 1000 * 1000
							case "GiB":
								amountBytes *= 1000 * 1000 * 1000
							}
							volumeSize = uint64(amountBytes)
						}
					}
					volumeSpecParts = volumeSpecParts[:len(volumeSpecParts)-1]
				}
			}

			if len(volumeSpecParts) > 0 {
				mountPath = volumeSpecParts[0]
				volumeSpecParts = volumeSpecParts[1:]
			}

			if len(volumeSpecParts) > 0 {
				return fmt.Errorf("Failed to parse volume specification: %s", volumeSpec)
			}

			volume := &pb.Volume{
				AccessMode: pb.Volume_VOLUME_RW_MANY, // TODO: pb.Volume_VOLUME_RW_ONE, pb.Volume_VOLUME_RO_MANY
			}

			if secretPath != "" {
				secretPathAbs, err := filepath.Abs(secretPath)
				if err != nil {
					return err
				}
				secretPathRel, err := filepath.Rel(manifestFileAbs, secretPathAbs)
				if err != nil {
					return err
				}
				volume.Configuration = &pb.Volume_Secret{
					Secret: &pb.Volume_SecretConfig{
						File: secretPathRel,
					},
				}
				if mountPath == "" {
					mountPath = "/" + secretPath
				}
			} else {
				volume.Configuration = &pb.Volume_Filesystem{
					Filesystem: &pb.Volume_FilesystemConfig{
						ResourceRequests: []*pb.Resource{
							{
								Resource: "storage",
								Quantity: &pb.Resource_Amount{
									Amount: volumeSize,
								},
							},
						},
					},
				}
			}

			volume.Name = disallowedVolumeNameCharacters.ReplaceAllString(mountPath, "")

			pod.Containers[0].Volumes = append(pod.Containers[0].Volumes, &pb.Container_VolumeMount{
				Name:      volume.Name,
				MountPath: mountPath,
				ReadOnly:  readOnly,
			})

			pod.Volumes = append(pod.Volumes, volume)
		}

		err = pb.MarshalFile(manifestFile, manifestFormat, pod)
		if err != nil {
			return fmt.Errorf("Failed writing the manifest file: %w", err)
		}

		fmt.Fprintf(cmd.ErrOrStderr(), "Initialized pod configuration in %s\n", manifestFile)

		return err
	},
}

func init() {
	podCmd.AddCommand(initPodCmd)

	initPodCmd.Flags().StringVarP(&initImageName, "image", "t", "docker.io/library/hello-world:latest", "Image to initialize the pod with")
	initPodCmd.Flags().StringArrayVarP(&initVolumes, "volume", "v", []string{}, "Volumes to initialize the pod with. Specify similar to Docker volumes; e.g. secret.txt:/path/in/container.txt for secrets and /path/in/container:rw,5GB for volume mounts")
	initPodCmd.Flags().StringVarP(&initPort, "http", "p", fmt.Sprintf("%d:%s", initDefaultPort, initDefaultHostname), "Port and hostname to initialize the pod with (XXXX is replaced with a random string)")
}
