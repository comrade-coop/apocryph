package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"

	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/spf13/cobra"
)

var disallowedVolumeNameCharacters = regexp.MustCompile("[\\.-_/\\\\0-9\n ]+$|^[0-9\\.-_/\\\\\n ]+|[\\.-_/\\\\\n ]")
var volumeSizeFlag = regexp.MustCompile("^ *([0-9]+) *([KMG]i?B) *$")
var initImageName string
var initVolumes []string
var initHostname string

var initPodCmd = &cobra.Command{
	Use:    "init [manifest.yaml]",
	Short:  "Parse a pod from a local manifest",
	Args:   cobra.MaximumNArgs(1),
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		manifestFile := "manifest.yaml"
		if len(args) >= 1 {
			manifestFile = args[0]
		}

		imageParts := strings.Split(initImageName, "/")

		initHostname = strings.Replace(initHostname, "XXXX", strconv.FormatInt(rand.Int63(), 16), 1)

		pod := &pb.Pod{
			Containers: []*pb.Container{
				{
					Name: imageParts[len(imageParts)-1],
					Image: &pb.Container_Image{
						Url: initImageName,
					},
					Ports: []*pb.Container_Port{
						{
							Name: "http",
							ContainerPort: 80,
							ExposedPort: &pb.Container_Port_HostHttpHost{
								HostHttpHost: initHostname,
							},
						},
					},
				},
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
								case "KB": amountBytes *= 1024
								case "MB": amountBytes *= 1024 * 1024
								case "GB": amountBytes *= 1024 * 1024 * 1024
								case "KiB": amountBytes *= 1000
								case "MiB": amountBytes *= 1000 * 1000
								case "GiB": amountBytes *= 1000 * 1000 * 1000
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
				volume.Configuration = &pb.Volume_Secret{
					Secret: &pb.Volume_SecretConfig{
						File: secretPath,
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
				Name: volume.Name,
				MountPath: mountPath,
				ReadOnly: readOnly,
			})

			pod.Volumes = append(pod.Volumes, volume)
		}

		err := pb.MarshalFile(manifestFile, manifestFormat, pod)
		if err != nil {
			return fmt.Errorf("Failed writing the manifest file: %w", err)
		}

		fmt.Fprintf(cmd.ErrOrStderr(), "Initialized pod configuration in %s\n", manifestFile)

		return err
	},
}

func init() {
	podCmd.AddCommand(initPodCmd)

	initPodCmd.Flags().StringVarP(&initImageName, "image", "t", "docker.io/library/hello-world", "Image to initialize the pod with")
	initPodCmd.Flags().StringArrayVarP(&initVolumes, "volume", "v", []string{}, "Volumes to initialize the pod with. Specify similar to Docker volumes; e.g. secret.txt:/path/in/container.txt for secrets and /path/in/container:rw,5GB for volume mounts")
	initPodCmd.Flags().StringVar(&initHostname, "http", "pod.XXXX.hostname.example", "Hostname to initialize the pod with (XXXX is replaced with a random string)")
}
