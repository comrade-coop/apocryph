package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	tpk8s "github.com/comrade-coop/trusted-pods/pkg/kubernetes"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var manifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Operations related to with raw pod manifests",
}

var formats = map[string]func(b []byte, m protoreflect.ProtoMessage) error{
	"json": protojson.Unmarshal,
	"pb":   proto.Unmarshal,
	"text": prototext.Unmarshal,
}

var manifestFormat string
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

		Unmarshal := formats[manifestFormat]
		if Unmarshal == nil {
			return errors.New("Unknown format: " + manifestFormat)
		}

		podManifest := &pb.Pod{}
		err = Unmarshal(manifestContents, podManifest)
		if err != nil {
			return err
		}

		cl, err := getNamespacedClient(cmd.Context())
		if err != nil {
			return err
		}

		response := &pb.ProvisionPodResponse{}
		err = tpk8s.ApplyPodRequest(cmd.Context(), cl, podManifest, response)
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

func getNamespacedClient(ctx context.Context) (client.Client, error) {
	var config *rest.Config
	var err error
	if kubeConfig == "-" {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
	}
	if err != nil {
		return nil, err
	}

	return tpk8s.GetNamespacedClient(ctx, config, dryRun)
}

func init() {
	manifestCmd.AddCommand(applyManifestCmd)

	formatNames := make([]string, 0, len(formats))
	for name := range formats {
		formatNames = append(formatNames, name)
	}
	applyManifestCmd.Flags().StringVar(&manifestFormat, "format", "pb", fmt.Sprintf("Manifest format. One of %v", formatNames))
	applyManifestCmd.Flags().BoolVarP(&dryRun, "dry-run", "z", false, "Dry run mode; modify nothing.")

	defaultKubeConfig := "-"
	if home := homedir.HomeDir(); home != "" {
		defaultKubeConfig = filepath.Join(home, ".kube", "config")
	}
	applyManifestCmd.Flags().StringVar(&kubeConfig, "kubeconfig", defaultKubeConfig, "absolute path to the kubeconfig file (- to use in-cluster config)")
}
