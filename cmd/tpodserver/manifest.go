package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	ipfs_utils "github.com/comrade-coop/trusted-pods/pkg/ipfs-utils"
	tpk8s "github.com/comrade-coop/trusted-pods/pkg/kubernetes"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/ipfs/boxo/coreiface/path"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
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

type provisionPodServer struct {
	pb.UnimplementedProvisionPodServiceServer
	ipfs *rpc.HttpApi
}

func transformError(err error) (*pb.ProvisionPodResponse, error) {
	return &pb.ProvisionPodResponse{
		Error: err.Error(),
	}, nil
}

func (s *provisionPodServer) ProvisionPod(ctx context.Context, request *pb.ProvisionPodRequest) (*pb.ProvisionPodResponse, error) {
	cid, err := cid.Cast(request.PodManifestCid)
	if err != nil {
		return transformError(err)
	}

	node, err := s.ipfs.Unixfs().Get(ctx, path.IpfsPath(cid))
	if err != nil {
		return transformError(err)
	}
	file, ok := node.(files.File)
	if !ok {
		return transformError(errors.New("Supplied CID not a file"))
	}
	manifestBytes, err := io.ReadAll(file)
	if err != nil {
		return transformError(err)
	}
	pod := &pb.Pod{}
	err = proto.Unmarshal(manifestBytes, pod)
	if err != nil {
		return transformError(err)
	}

	cl, err := getNamespacedClient(ctx)
	if err != nil {
		return nil, err
	}
	response := &pb.ProvisionPodResponse{}
	err = tpk8s.ApplyPodRequest(ctx, cl, pod, response)
	if err != nil {
		return transformError(err)
	}

	return response, nil
}

var serveManifestCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve a service listening for incomming manifests",
	RunE: func(cmd *cobra.Command, args []string) error {
		ipfs, err := ipfs_utils.ConnectToLocalNode()
		if err != nil {
			return err
		}

		listener, err := ipfs_utils.NewP2pApi(ipfs).Listen(pb.ProvisionPod)
		if err != nil {
			return err
		}

		s := grpc.NewServer()
		pb.RegisterProvisionPodServiceServer(s, &provisionPodServer{
			ipfs: ipfs,
		})

		go s.Serve(listener)
		defer s.Stop()

		<-cmd.Context().Done()

		s.GracefulStop()

		return nil
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
	manifestCmd.AddCommand(serveManifestCmd)

	formatNames := make([]string, 0, len(formats))
	for name := range formats {
		formatNames = append(formatNames, name)
	}
	applyManifestCmd.Flags().StringVar(&manifestFormat, "format", "pb", fmt.Sprintf("Manifest format. One of %v", formatNames))
	applyManifestCmd.Flags().BoolVarP(&dryRun, "dry-run", "z", false, "Dry run mode; modify nothing.")
	serveManifestCmd.Flags().BoolVarP(&dryRun, "dry-run", "z", false, "Dry run mode; modify nothing.")

	defaultKubeConfig := "-"
	if home := homedir.HomeDir(); home != "" {
		defaultKubeConfig = filepath.Join(home, ".kube", "config")
	}
	applyManifestCmd.Flags().StringVar(&kubeConfig, "kubeconfig", defaultKubeConfig, "absolute path to the kubeconfig file (- to use in-cluster config)")
	serveManifestCmd.Flags().StringVar(&kubeConfig, "kubeconfig", defaultKubeConfig, "absolute path to the kubeconfig file (- to use in-cluster config)")
}
