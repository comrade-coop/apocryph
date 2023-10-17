package main

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	tpk8s "github.com/comrade-coop/trusted-pods/pkg/kubernetes"
	"github.com/comrade-coop/trusted-pods/pkg/prometheus"
	"github.com/comrade-coop/trusted-pods/pkg/resource"
	tpgsrpc "github.com/comrade-coop/trusted-pods/pkg/substrate"
	tptypes "github.com/comrade-coop/trusted-pods/pkg/substrate/types"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	k8client "sigs.k8s.io/controller-runtime/pkg/client"
)

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor the metrics stored in prometheus",
	RunE: func(cmd *cobra.Command, args []string) error {
		pro := prometheus.NewPrometheusAPI(prometheusUrl)

		pricingTable, err := openPricingTable()
		if err != nil {
			return err
		}
		if pricingTable == nil {
			return errors.New("Pricing table is required for the monitor")
		}

		client, err := tpk8s.GetClient(kubeConfig, dryRun)
		if err != nil {
			return err
		}

		var api *tpgsrpc.SubstrateAPI
		if !dryRun {
			api, err = tpgsrpc.NewSubstrateAPI(chainRpc)
			if err != nil {
				return err
			}
		}

	Loop:
		for {
			resourceMeasurements := resource.ResourceMeasurementsMap{}
			err := pro.FetchResourceMetrics(resourceMeasurements)
			if err != nil {
				return err
			}

			amountsOwed := resourceMeasurements.Price(pricingTable)

			namespaces := &corev1.NamespaceList{}
			err = client.List(cmd.Context(), namespaces, k8client.HasLabels{tpk8s.LabelTrustedPodsPaymentChannel})
			if err != nil {
				return err
			}

			for _, n := range namespaces.Items {
				if amountOwed, ok := amountsOwed[n.Name]; ok {
					contractAddress := n.ObjectMeta.Labels[tpk8s.LabelTrustedPodsPaymentChannel]
					// TODO: Watch amounts left in contract!
					// TODO: Only update contract if the difference is "too large"
					err = setAmountOwed(api, contractAddress, types.NewU128(*amountOwed))
					if err != nil {
						fmt.Printf("Error while processing namespace %s: %v", n.Name, err)
					}
				}
			}
			select {
			case <-cmd.Context().Done():
				break Loop
			case <-time.After(time.Minute * 1):
				continue
			}
		}

		return nil
	},
}

func setAmountOwed(api *tpgsrpc.SubstrateAPI, contractAddressString string, amount types.U128) error {
	_, contractAddress, err := tptypes.NewAccountIDFromSS58(contractAddressString)
	if err != nil {
		return err
	}

	if api != nil {
		properties, err := api.RPC.System.Properties()
		if err != nil {
			return err
		}

		from, err := signature.KeyringPairFromSecret(substrateKey, uint16(properties.AsSS58Format))
		if err != nil {
			return err
		}

		inputData, err := tptypes.EncodeInputData(ClaimSelector, types.UCompact(*amount.Int))
		if err != nil {
			return err
		}

		hash, err := api.RPC.Contracts.CallContractNoWait(*contractAddress, from, inputData, types.NewU128(*big.NewInt(0)))
		if err != nil {
			return err
		}

		fmt.Printf("Claimed funds from contract %s: %d tx %s\n", tptypes.AccountIDToSS58(42, contractAddress), amount, hash.Hex())
	} else {
		fmt.Printf("Claimed funds from contract %s: %d (dry-run: no tx) \n", tptypes.AccountIDToSS58(42, contractAddress), amount)
	}

	return nil
}

func init() {
	monitorCmd.Flags().BoolVarP(&dryRun, "dry-run", "z", false, "Dry run mode; modify nothing.")
	monitorCmd.Flags().StringVar(&kubeConfig, "kubeconfig", "-", "absolute path to the kubeconfig file (- to the first of in-cluster config and ~/.kube/config)")
	monitorCmd.Flags().StringVar(&prometheusUrl, "prometheus", "", "address at which the prometheus API can be accessed")
	monitorCmd.Flags().StringVar(&chainRpc, "rpc", "", "Link to the Substrate RPC.")
}
