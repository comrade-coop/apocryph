package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	keyservice "github.com/comrade-coop/trusted-pods/pkg/crypto"
	ipfs_utils "github.com/comrade-coop/trusted-pods/pkg/ipfs-utils"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	iface "github.com/ipfs/boxo/coreiface"
	"github.com/ipfs/boxo/files"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var podCmd = &cobra.Command{
	Use:   "pod",
	Short: "Operations related to with raw pod manifests",
}

var manifestFormat string
var providerPeer string
var ipfsApi string
var paymentContract string

func TransformSecrets(ctx context.Context, ipfs iface.CoreAPI, basepath string, pod *pb.Pod, keys *[]*pb.Key) error {
	for _, volume := range pod.Volumes {
		if volume.Type == pb.Volume_VOLUME_SECRET {
			secretConfig := volume.GetSecret()
			if secretConfig.File != "" {
				secretPath := secretConfig.File
				if !path.IsAbs(secretPath) {
					secretPath = path.Join(basepath, secretPath)
				}
				secretFile, err := os.Open(secretPath)
				if err != nil {
					return err
				}
				secretBytes, err := io.ReadAll(secretFile)
				if err != nil {
					return err
				}
				keyData, err := keyservice.CreateRandomKey()
				if err != nil {
					return err
				}

				encryptedData, nonce, err := keyservice.AESEncryptWith(secretBytes, keyData)
				if len(nonce) != keyservice.NONCE_SIZE {
					return errors.New("Wrong nonce size")
				}
				encryptedSecretBytes := append(nonce, encryptedData...)

				encryptedSecretPath, err := ipfs.Unixfs().Add(ctx, files.NewBytesFile(encryptedSecretBytes))
				if err != nil {
					return err
				}

				secretConfig.Cid = encryptedSecretPath.Cid().Bytes()
				secretConfig.KeyIdx = int32(len(*keys))
				secretConfig.File = ""
				*keys = append(*keys, &pb.Key{
					Data: keyData,
				})
			}
		}
	}
	return nil
}

var deployPodCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a pod from a local manifest",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ipfs, ipfsMultiaddr, err := ipfs_utils.GetIpfsClient(ipfsApi)
		if err != nil {
			return err
		}

		podPath := args[0]

		file, err := os.Open(podPath)
		if err != nil {
			return err
		}

		podManifestCotents, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		pod := &pb.Pod{}
		err = pb.Unmarshal(manifestFormat, podManifestCotents, pod)
		if err != nil {
			return err
		}

		keys := []*pb.Key{}

		err = TransformSecrets(cmd.Context(), ipfs, path.Dir(podPath), pod, &keys)
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

		request := &pb.ProvisionPodRequest{
			PodManifestCid: podManifestPath.Cid().Bytes(),
			Keys:           keys,
			ClientAddress:  ClientAddress,
			TokenAddress:   TokenContractAddress,
			PodID:          []byte{1},
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

	deployPodCmd.Flags().StringVar(&manifestFormat, "format", "pb", fmt.Sprintf("Manifest format. One of %v", pb.UnmarshalFormatNames))
	deployPodCmd.Flags().StringVar(&ipfsApi, "ipfs", "-", "multiaddr where the ipfs/kubo api can be accessed (- to use the daemon running in IPFS_PATH)")

}
