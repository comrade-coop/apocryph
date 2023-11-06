package main

import (
	"fmt"
	"path"

	"github.com/bufbuild/protoyaml-go"
	tpipfs "github.com/comrade-coop/trusted-pods/pkg/ipfs"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/comrade-coop/trusted-pods/pkg/publisher"
	"github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

var podCmd = &cobra.Command{
	Use:   "pod",
	Short: "Operations related to with raw pod manifests",
}

var manifestFormat string
var providerPeer string
var ipfsApi string
var paymentContract string
var noIpdr bool
var readSecrets bool

var parsePodCmd = &cobra.Command{
	Use:    "parse",
	Short:  "Parse a pod from a local manifest",
	Args:   cobra.ExactArgs(1),
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		manifestFile := args[0]

		pod := &pb.Pod{}
		err := pb.UnmarshalFile(manifestFile, manifestFormat, pod)
		if err != nil {
			return fmt.Errorf("Failed reading the manifest file: %w", err)
		}

		if readSecrets {
			err = tpipfs.TransformSecrets(pod,
				tpipfs.ReadSecrets(path.Dir(manifestFile)),
			)
			if err != nil {
				return fmt.Errorf("Failed reading from the filesystem: %w", err)
			}
		}

		resultBytes, err := protoyaml.MarshalOptions{}.Marshal(pod)
		if err != nil {
			return fmt.Errorf("Failed reading the manifest file: %w", err)
		}

		_, err = cmd.OutOrStdout().Write(resultBytes)

		return err
	},
}

var uploadPodCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a pod from a local manifest",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		partialRequest, _, err := publisher.UploadManifest(cmd.Context(), args[0], manifestFormat, ipfsApi, noIpdr)
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.ErrOrStderr(), "Pass the following parameter instead of a pod file to deploy the pod\n")

		_, err = fmt.Fprintln(cmd.OutOrStdout(), protojson.MarshalOptions{}.Format(partialRequest))
		if err != nil {
			return err
		}

		return nil
	},
}

var deployPodCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a pod from a local manifest",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ipfs, ipfsMultiaddr, err := tpipfs.GetIpfsClient(ipfsApi)
		if err != nil {
			return fmt.Errorf("Failed connectig to IPFS: %w", err)
		}

		request := &pb.ProvisionPodRequest{}

		err = protojson.Unmarshal([]byte(args[0]), request)
		if err != nil {
			request, _, err = publisher.UploadManifest(cmd.Context(), args[0], manifestFormat, ipfsApi, noIpdr)
			if err != nil {
				return err
			}
		}

		var podIdBytes common.Hash
		if podId != "" {
			podIdBytes = common.HexToHash(podId)
		} else {
			podIdBytes = common.BytesToHash(request.PodManifestCid)
		}

		request.Payment, err = createPaymentChannel(podIdBytes)
		if err != nil {
			return err
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

		fmt.Fprintf(cmd.ErrOrStderr(), "Dialed provider over IPFS p2p\n")

		defer conn.Close()

		client := pb.NewProvisionPodServiceClient(conn)

		fmt.Fprintf(cmd.ErrOrStderr(), "Sending request...\n")

		response, err := client.ProvisionPod(cmd.Context(), request)
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.ErrOrStderr(), "Success! Received result...\n")

		_, err = fmt.Fprintln(cmd.OutOrStdout(), protojson.Format(response))
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	podCmd.AddCommand(deployPodCmd)
	podCmd.AddCommand(uploadPodCmd)
	podCmd.AddCommand(parsePodCmd)

	podCmd.PersistentFlags().StringVar(&manifestFormat, "format", "", fmt.Sprintf("Manifest format. One of %v (leave empty to auto-detect)", pb.UnmarshalFormatNames))

	podCmd.PersistentFlags().StringVar(&ipfsApi, "ipfs", "", "multiaddr where the ipfs/kubo api can be accessed (leave blank to use the daemon running in IPFS_PATH)")

	deployPodCmd.Flags().StringVar(&ethereumRpc, "ethereum-rpc", "http://127.0.0.1:8545", "ethereum rpc node")
	deployPodCmd.Flags().StringVar(&publisherKey, "ethereum-key", "", "account string (private key | http[s]://clef#account | /keystore#account | account (in default keystore))")

	deployPodCmd.Flags().StringVar(&providerPeer, "provider", "", "provider peer id")
	deployPodCmd.Flags().StringVar(&providerEthAddress, "provider-eth", "", "provider public address")

	deployPodCmd.Flags().StringVar(&paymentContractAddress, "payment-contract", "", "payment contract address")
	deployPodCmd.Flags().StringVar(&podId, "pod-id", "", "pod id (empty to pick one automatically)")
	deployPodCmd.Flags().StringVar(&funds, "funds", "5000000000000000000", "intial funds")
	deployPodCmd.Flags().Int64Var(&unlockTime, "unlock-time", 5*60, "time for unlocking tokens (in seconds)")
	podCmd.PersistentFlags().BoolVar(&noIpdr, "no-ipdr", false, "disable ipdr")

	parsePodCmd.Flags().BoolVar(&readSecrets, "with-secrets", false, "include secrets in the parsed output")

}
