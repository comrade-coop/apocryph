package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	ipfs_utils "github.com/comrade-coop/trusted-pods/pkg/ipfs-utils"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	tptypes "github.com/comrade-coop/trusted-pods/pkg/substrate/types"
	"github.com/ipfs/boxo/files"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var podCmd = &cobra.Command{
	Use:   "pod",
	Short: "Operations related to with raw pod manifests",
}

var formats = map[string]func(b []byte, m protoreflect.ProtoMessage) error{
	"json": protojson.Unmarshal,
	"pb":   proto.Unmarshal,
	"text": prototext.Unmarshal,
}

var manifestFormat string
var providerPeer string
var ipfsApi string
var paymentContract string

var deployPodCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a pod from a local manifest",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ipfs, ipfsMultiaddr, err := ipfs_utils.GetIpfsClient(ipfsApi)
		if err != nil {
			return err
		}

		file, err := os.Open(args[0])
		if err != nil {
			return err
		}

		podManifestCotents, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		Unmarshal := formats[manifestFormat]
		if Unmarshal == nil {
			return errors.New("Unknown format: " + manifestFormat)
		}

		pod := &pb.Pod{}
		err = Unmarshal(podManifestCotents, pod)
		if err != nil {
			return err
		}

		podManifestBytes, err := proto.Marshal(pod)
		if err != nil {
			return err
		}

		podManifestPath, err := ipfs.Unixfs().Add(cmd.Context(), files.NewBytesFile(podManifestBytes))
		if err != nil {
			return err
		}

		_, paymentContractAddress, err := tptypes.NewAccountIDFromSS58(paymentContract)
		if err != nil {
			return err
		}

		request := &pb.ProvisionPodRequest{
			PodManifestCid:         podManifestPath.Cid().Bytes(),
			Keys:                   []*pb.Key{},
			PaymentContractAddress: paymentContractAddress.ToBytes(),
		}

		providerPeerId, err := peer.Decode(providerPeer)
		if err != nil {
			return err
		}

		addr, err := ipfs_utils.NewP2pApi(ipfs, ipfsMultiaddr).ConnectTo(pb.ProvisionPod, providerPeerId)
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

		response, err := client.ProvisionPod(cmd.Context(), request)
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
	podCmd.AddCommand(deployPodCmd)

	formatNames := make([]string, 0, len(formats))
	for name := range formats {
		formatNames = append(formatNames, name)
	}
	deployPodCmd.Flags().StringVar(&manifestFormat, "format", "pb", fmt.Sprintf("Manifest format. One of %v", formatNames))
	deployPodCmd.Flags().StringVar(&providerPeer, "provider", "", "P2p identity of the provider to deploy to")
	deployPodCmd.Flags().StringVar(&paymentContract, "payment", "", "Payment contract address.")

	deployPodCmd.Flags().StringVar(&ipfsApi, "ipfs", "-", "multiaddr where the ipfs/kubo api can be accessed (- to use the daemon running in IPFS_PATH)")

}
