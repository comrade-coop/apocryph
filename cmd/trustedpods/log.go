package main

import (
	"fmt"
	"io"

	tpipfs "github.com/comrade-coop/trusted-pods/pkg/ipfs"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

var containerName string

var logPodCmd = &cobra.Command{
	Use:   "log",
	Short: "get pod conatiner logs",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ipfs, ipfsMultiaddr, err := tpipfs.GetIpfsClient(ipfsApi)
		if err != nil {
			return err
		}

		request := &pb.PodLogRequest{
			ContainerName: containerName,
			Credentials:   &pb.Credentials{},
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

		stream, err := client.GetPodLogs(cmd.Context(), request)
		if err != nil {
			return err
		}

		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), protojson.Format(response))
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	podCmd.AddCommand(logPodCmd)

	logPodCmd.Flags().StringVar(&ipfsApi, "ipfs", "", "multiaddr where the ipfs/kubo api can be accessed (leave blank to use the daemon running in IPFS_PATH)")

	logPodCmd.Flags().StringVar(&providerPeer, "provider", "", "provider peer id")
	logPodCmd.Flags().StringVar(&containerName, "container", "", "pod namespace (returned from deploy pod)")
}
