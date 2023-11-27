package main

import (
	"fmt"
	"io"

	"github.com/comrade-coop/trusted-pods/pkg/ethereum"
	tpipfs "github.com/comrade-coop/trusted-pods/pkg/ipfs"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/comrade-coop/trusted-pods/pkg/publisher"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

var containerName string

var logPodCmd = &cobra.Command{
	Use:     fmt.Sprintf("log [%s] [deployment.yaml]", publisher.DefaultPodFile),
	Short:   "get pod conatiner logs",
	Args:    cobra.ExactArgs(1),
	GroupID: "main",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, _, pod, deployment, err := publisher.ReadPodAndDeployment(args, manifestFormat, deploymentFormat)
		if err != nil {
			return err
		}

		ipfs, ipfsMultiaddr, err := tpipfs.GetIpfsClient(ipfsApi)
		if err != nil {
			return err
		}

		ethClient, err := ethereum.GetClient(ethereumRpc)
		if err != nil {
			return err
		}

		if publisherKey == "" {
			publisherKey = common.BytesToAddress(deployment.Payment.PublisherAddress).String()
		}

		publisherAuth, sign, err := ethereum.GetAccountAndSigner(publisherKey, ethClient)
		if err != nil {
			return fmt.Errorf("Could not get ethereum account: %w", err)
		}

		publisherEthAddress := publisherAuth.From.Bytes()

		interceptor := pb.NewAuthInterceptor(deployment, pb.GetPodLogs, expirationOffset, sign)

		conn, err := publisher.ConnectToProvider(tpipfs.NewP2pApi(ipfs, ipfsMultiaddr), deployment, &interceptor, nil)
		if err != nil {
			return err
		}
		defer conn.Close()

		client := pb.NewProvisionPodServiceClient(conn)

		if containerName == "" {
			if len(pod.Containers) != 1 {
				return fmt.Errorf("Specifying a container name is required for pods with more than one container")
			}
			containerName = pod.Containers[0].Name
		}

		request := &pb.PodLogRequest{
			ContainerName:    containerName,
			PublisherAddress: publisherEthAddress,
		}

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

	logPodCmd.Flags().AddFlagSet(deploymentFlags)
	logPodCmd.Flags().AddFlagSet(syncFlags)
}
