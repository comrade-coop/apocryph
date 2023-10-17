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
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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

		var config *rest.Config
		if kubeConfig == "-" {
			config, err = rest.InClusterConfig()
		} else {
			config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		}
		if err != nil {
			return err
		}

		scheme, err := tpk8s.GetScheme()
		if err != nil {
			return err
		}

		client, err := k8client.New(config, k8client.Options{
			Scheme: scheme,
			DryRun: &dryRun,
		})
		if err != nil {
			return err
		}

	Loop:
		for {
			select {
			case <-cmd.Context().Done():
				break Loop
			case <-time.After(time.Minute * 1):
				resourceMeasurements := resource.ResourceMeasurementsMap{}
				err := pro.FetchResourceMetrics(resourceMeasurements)
				if err != nil {
					return err
				}

				amountsOwed := resourceMeasurements.Price(pricingTable)

				req, err := labels.NewRequirement(tpk8s.LabelTrustedPodsPaymentChannel, selection.Exists, []string{})
				if err != nil {
					return err
				}

				var namespaces *corev1.NamespaceList
				client.List(cmd.Context(), namespaces, &k8client.ListOptions{ // TODO: Watch?
					LabelSelector: labels.NewSelector().Add(*req),
				})

				for _, n := range namespaces.Items {
					if amountOwed, ok := amountsOwed[n.Name]; ok {
						contractAddress := n.ObjectMeta.Labels[tpk8s.LabelTrustedPodsPaymentChannel]
						// TODO: Watch amounts left in contract!
						// TODO: Only update contract if the difference is "too large"
						setAmountOwed(contractAddress, types.NewU128(*amountOwed))
					}
				}
			}
		}

		return nil
	},
}

func setAmountOwed(contractAddressString string, amount types.U128) error {
	_, contractAddress, err := tptypes.NewAccountIDFromSS58(contractAddressString)
	if err != nil {
		return err
	}

	api, err := tpgsrpc.NewSubstrateAPI(chainRpc)
	if err != nil {
		return err
	}

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

	if !dryRun {
		hash, err := api.RPC.Contracts.CallContractNoWait(*contractAddress, from, inputData, types.NewU128(*big.NewInt(0)))
		if err != nil {
			return err
		}

		fmt.Printf("Claimed funds from contract %s: %d tx %s\n", tptypes.AccountIDToSS58(0, contractAddress), amount, hash.Hex())
	} else {
		fmt.Printf("Claimed funds from contract %s: %d \n", tptypes.AccountIDToSS58(0, contractAddress), amount)
	}

	return nil
}

func init() {
	monitorCmd.Flags().BoolVarP(&dryRun, "dry-run", "z", false, "Dry run mode; modify nothing.")
	monitorCmd.Flags().StringVar(&kubeConfig, "kubeconfig", "-", "absolute path to the kubeconfig file (- to the first of in-cluster config and ~/.kube/config)")
	monitorCmd.Flags().StringVar(&prometheusUrl, "prometheus", "", "address at which the prometheus API can be accessed")
	monitorCmd.Flags().StringVar(&chainRpc, "rpc", "", "Link to the Substrate RPC.")
}
