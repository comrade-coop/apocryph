package main

import (
	"fmt"
	"io"
	"os"

	tpk8s "github.com/comrade-coop/trusted-pods/pkg/kubernetes"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var manifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Operations related to with raw pod manifests",
}

var manifestFormat string
var paymentContract string
var kubeConfig string
var dryRun bool

var applyManifestCmd = &cobra.Command{
	Use:   "apply <file>",
	Short: "Apply a manifest from a file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := os.Open(args[0])
		if err != nil {
			return err
		}

		manifestContents, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		pod := &pb.Pod{}
		err = pb.Unmarshal(manifestFormat, manifestContents, pod)
		if err != nil {
			return err
		}

		cl, err := tpk8s.GetClient(kubeConfig, dryRun)
		if err != nil {
			return err
		}

		response := &pb.ProvisionPodResponse{}
		err = tpk8s.RunInNamespaceOrRevert(cmd.Context(), cl, tpk8s.NewTrustedPodsNamespace(paymentContract), func(cl client.Client) error {
			return tpk8s.ApplyPodRequest(cmd.Context(), cl, pod, response)
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

	applyManifestCmd.Flags().StringVar(&manifestFormat, "format", "pb", fmt.Sprintf("Manifest format. One of %v", pb.UnmarshalFormatNames))
	applyManifestCmd.Flags().BoolVarP(&dryRun, "dry-run", "z", false, "Dry run mode; modify nothing.")
	applyManifestCmd.Flags().StringVar(&paymentContract, "payment", "", "Payment contract address.")
	applyManifestCmd.Flags().StringVar(&kubeConfig, "kubeconfig", "-", "absolute path to the kubeconfig file (- to the first of in-cluster config and ~/.kube/config)")
}
