package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"github.com/spf13/cobra"

	pb "github.com/comrade-coop/apocryph/pkg/proto"
	pbcon "github.com/comrade-coop/apocryph/pkg/proto/protoconnect"
)

var nodesAddrs []string
var raftPath string
var appManifest string
var baseUrl string

var autoscalerCmd = &cobra.Command{
	Use:   "autoscale",
	Short: "Configure the deployed autonomous autoscaler",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := pbcon.NewAutoscalerServiceClient(
			http.DefaultClient,
			baseUrl)
		log.Printf("%v\n", nodesAddrs)
		request := connect.NewRequest(&pb.ConnectClusterRequest{NodeGateway: baseUrl, Servers: nodesAddrs, Timeout: 10})
		response, err := client.ConnectCluster(context.Background(), request)
		if err != nil {
			return fmt.Errorf("Error connecting to cluster:%v", err)
		}
		log.Printf("Received response: %v", response)
		return nil

	},
}

func init() {
	// this could be removed if we store the list of providers directly in the pod configuration
	autoscalerCmd.Flags().StringSliceVarP(&nodesAddrs, "providers", "p", []string{}, "List of All providers cluster ip addresses that the autoscaler will use to redeploy your application in case of a failure")
	autoscalerCmd.Flags().StringVar(&raftPath, "path", "", "Optional: path where raft will save it's state (Default is In Memory)")
	autoscalerCmd.Flags().StringVar(&appManifest, "manifest", "", "path to application manifest to autoscale")
	autoscalerCmd.Flags().StringVar(&baseUrl, "url", "", "apocryph node gateway running autoscaler instance")

	rootCmd.AddCommand(autoscalerCmd)
}
