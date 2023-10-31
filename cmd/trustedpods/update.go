package main

import (
	"fmt"
	"os"
	"path"

	tpipfs "github.com/comrade-coop/trusted-pods/pkg/ipfs"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

var namespace string

var updatePodCmd = &cobra.Command{
	Use:   "update",
	Short: "update a pod from a local manifest",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ipfs, ipfsMultiaddr, err := tpipfs.GetIpfsClient(ipfsApi)
		if err != nil {
			return err
		}
		podPath := os.Args[0]

		pod := &pb.Pod{}
		err = pb.UnmarshalFile(podPath, manifestFormat, pod)
		if err != nil {
			return err
		}

		keys := []*pb.Key{}

		err = tpipfs.TransformSecrets(pod,
			tpipfs.ReadSecrets(path.Dir(podPath)),
			tpipfs.EncryptSecrets(&keys),
			tpipfs.UploadSecrets(cmd.Context(), ipfs),
		)
		if err != nil {
			return err
		}

		request := &pb.UpdatePodRequest{
			Pod:       pod,
			Namespace: namespace,
		}

		providerPeerId, err := peer.Decode(providerPeer)
		if err != nil {
			return err
		}

		addr, err := tpipfs.NewP2pApi(ipfs, ipfsMultiaddr).ConnectTo(pb.ProvisionPod, providerPeerId)
		if err != nil {
			return err
		}

		defer addr.Close()

		conn, err := grpc.Dial(addr.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}

		defer conn.Close()

		client := pb.NewProvisionPodServiceClient(conn)

		response, err := client.UpdatePod(cmd.Context(), request)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintln(cmd.OutOrStdout(), protojson.Format(response))
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	podCmd.AddCommand(updatePodCmd)

	updatePodCmd.Flags().StringVar(&manifestFormat, "format", "", fmt.Sprintf("Manifest format. One of %v", pb.UnmarshalFormatNames))

	updatePodCmd.Flags().StringVar(&ipfsApi, "ipfs", "", "multiaddr where the ipfs/kubo api can be accessed (leave blank to use the daemon running in IPFS_PATH)")

	updatePodCmd.Flags().StringVar(&providerPeer, "provider", "", "provider peer id")
	updatePodCmd.Flags().StringVar(&namespace, "namespace", "", "pod namespace (returned from deploy pod)")
}
