// SPDX-License-Identifier: GPL-3.0

package main

import (
	"fmt"
	"path/filepath"

	tpipfs "github.com/comrade-coop/apocryph/pkg/ipfs"
	tpk8s "github.com/comrade-coop/apocryph/pkg/kubernetes"
	pb "github.com/comrade-coop/apocryph/pkg/proto"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var manifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Operations related to with raw pod manifests",
}

var manifestFormat string
var kubeConfig string
var dryRun bool

var applyManifestCmd = &cobra.Command{
	Use:   "apply <file>",
	Short: "Apply a manifest from a file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		manifestPath := args[0]

		pod := &pb.Pod{}
		err := pb.UnmarshalFile(manifestFormat, manifestPath, pod)
		if err != nil {
			return err
		}

		images := make(map[string]string)
		for _, c := range pod.Containers {
			images[c.Name] = c.Image.Url
		}

		secrets := make(map[string][]byte)
		for _, v := range pod.Volumes {
			if v.Type == pb.Volume_VOLUME_SECRET {
				secret := v.GetSecret()
				contents, err := tpipfs.ReadSecret(filepath.Dir(manifestPath), secret)
				if err != nil {
					return err
				}
				secrets[v.Name] = contents
			}
		}

		cl, err := tpk8s.GetClient(kubeConfig, dryRun)
		if err != nil {
			return err
		}

		response := &pb.ProvisionPodResponse{}
		namespace := tpk8s.NewTrustedPodsNamespace("tpods-xx", pod, nil)
		err = tpk8s.RunInNamespaceOrRevert(cmd.Context(), cl, namespace, dryRun, func(cl client.Client) error {
			return tpk8s.ApplyPodRequest(cmd.Context(), cl, namespace.ObjectMeta.Name, false, pod, nil, images, secrets, response, "")
		})
		if err != nil {
			return err
		}

		result, err := protojson.Marshal(response)
		if err != nil {
			return err
		}
		_, err = cmd.OutOrStdout().Write(result)
		return err
	},
}

func init() {
	manifestCmd.AddCommand(applyManifestCmd)

	applyManifestCmd.Flags().StringVar(&manifestFormat, "format", "", fmt.Sprintf("Manifest format. One of %v", pb.FormatNames))

	applyManifestCmd.Flags().BoolVarP(&dryRun, "dry-run", "z", false, "Dry run mode; modify nothing.")
	applyManifestCmd.Flags().StringVar(&kubeConfig, "kubeconfig", "", "absolute path to the kubeconfig file (leave blank for the first of in-cluster config and ~/.kube/config)")

}
