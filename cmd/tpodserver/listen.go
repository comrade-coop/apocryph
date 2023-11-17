package main

import (
	"context"
	"net"
	"net/http"

	"github.com/comrade-coop/trusted-pods/pkg/ethereum"
	tpipfs "github.com/comrade-coop/trusted-pods/pkg/ipfs"
	tpk8s "github.com/comrade-coop/trusted-pods/pkg/kubernetes"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/comrade-coop/trusted-pods/pkg/provider"
	"github.com/spf13/cobra"
)

var ipfsApi string
var serveAddress string
var localOciRegistry string

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Start a service listening for incomming execution requests",
	RunE: func(cmd *cobra.Command, args []string) error {
		ipfs, ipfsMultiaddr, err := tpipfs.GetIpfsClient(ipfsApi)
		if err != nil {
			return err
		}

		ethClient, err := ethereum.GetClient(ethereumRpc)
		if err != nil {
			return err
		}

		providerAuth, _, err := ethereum.GetAccountAndSigner(providerKey, ethClient)
		if err != nil {
			return err
		}

		pricingTables, err := openPricingTables()
		if err != nil {
			return err
		}

		validator, err := ethereum.NewPaymentChannelValidator(ethClient, pricingTables, providerAuth)

		var listener net.Listener
		if serveAddress == "" {
			listener, err = tpipfs.NewP2pApi(ipfs, ipfsMultiaddr).Listen(pb.ProvisionPod)
		} else {
			listener, err = net.Listen("tcp", serveAddress)
		}
		if err != nil {
			return err
		}

		k8cl, err := tpk8s.GetClient(kubeConfig, dryRun)
		if err != nil {
			return err
		}

		mux := http.NewServeMux()
		mux.Handle(provider.NewTPodServerHandler(ipfs, dryRun, k8cl, localOciRegistry, validator, "loki.loki.svc.cluster.local:3100"))
		server := &http.Server{Handler: mux}

		go server.Serve(listener)

		defer server.Close()

		<-cmd.Context().Done()

		server.Shutdown(context.TODO()) // cmd.Context() is already done.. :/

		return nil
	},
}

func init() {
	listenCmd.Flags().StringVar(&serveAddress, "address", "", "port to serve on (leave blank to automatically pick a port and register a listener for it in ipfs)")

	listenCmd.Flags().BoolVarP(&dryRun, "dry-run", "z", false, "Dry run mode; modify nothing.")
	listenCmd.Flags().StringVar(&kubeConfig, "kubeconfig", "", "absolute path to the kubeconfig file (leave blank for the first of in-cluster config and ~/.kube/config)")
	listenCmd.Flags().StringVar(&ipfsApi, "ipfs", "", "multiaddr where the ipfs/kubo api can be accessed (leave blank to use the daemon running in IPFS_PATH)")
	listenCmd.Flags().StringVar(&localOciRegistry, "oci-registry", "", "OCI registry used to resolve IPDR images")
	listenCmd.Flags().StringVar(&ethereumRpc, "ethereum-rpc", "http://127.0.0.1:8545", "client public address")
	listenCmd.Flags().StringVar(&providerKey, "ethereum-key", "", "provider account string (private key | http[s]://clef#account | /keystore#account | account (in default keystore))")
}
