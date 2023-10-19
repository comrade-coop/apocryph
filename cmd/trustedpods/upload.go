package main

import (
	"log"

	podmanagement "github.com/comrade-coop/trusted-pods/pkg/pod-management"
	"github.com/spf13/cobra"
)

var podName string
var psw string
var manifest string
var packagePath string

// createCmd represents the create command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a pod package to IPFS",
	RunE: func(cmd *cobra.Command, args []string) error {

		pod, err := podmanagement.CreatePod(podName, psw)
		if err != nil {
			return err
		}

		err = pod.AssignManifest(manifest)
		if err != nil {
			return err
		}

		provider, err := podmanagement.CreateIpfsUploader()
		if err != nil {
			return err
		}

		cid, err := pod.UploadPackage(provider, packagePath)
		if err != nil {
			return err
		}

		log.Println("pod package cid:", cid)
		return nil
	},
}

func init() {
	podCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().StringVar(&podName, "name", "my pod", "pod name")
	uploadCmd.Flags().StringVar(&psw, "password", "123", "pod password")
	uploadCmd.Flags().StringVar(&manifest, "manifest", "./manifest.yaml", "manifest file path")
	uploadCmd.Flags().StringVar(&packagePath, "packagePath", "pod-package", "final pod package path")

}
