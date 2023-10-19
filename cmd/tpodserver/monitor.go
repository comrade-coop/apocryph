package main

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/comrade-coop/trusted-pods/pkg/abi"
	"github.com/comrade-coop/trusted-pods/pkg/contracts"
	"github.com/comrade-coop/trusted-pods/pkg/prometheus"
	"github.com/comrade-coop/trusted-pods/pkg/resource"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
)

var clientAddress common.Address
var tokenAddress common.Address
var paymentAddress common.Address

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor the metrics stored in prometheus",
	RunE: func(cmd *cobra.Command, args []string) error {
		pro := prometheus.NewPrometheusAPI(prometheusUrl)

		err := contracts.VerifyContractAddress(PaymentContractAddress, allowedContractCodeHashes)
		if err != nil {
			return err
		}
		err = setUp()
		if err != nil {
			return err
		}

		clientAddress = common.HexToAddress(ClientAddress)
		tokenAddress = common.HexToAddress(TokenContractAddress)
		paymentAddress = common.HexToAddress(PaymentContractAddress)

		pricingTable, err := openPricingTable()
		if err != nil {
			return err
		}
		if pricingTable == nil {
			return errors.New("Pricing table is required for the monitor")
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

			for _, n := range namespaces.Items {
				if amountOwed, ok := amountsOwed[n.Name]; ok {
					// TODO: Watch amounts left in contract!
					// TODO: Only update contract if the difference is "too large"
					err = setAmountOwed(ethClient, Instance, amountOwed)
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

func setAmountOwed(client *ethclient.Client, instance *abi.Payment, amount *big.Int) error {
	_, err := instance.WithdrawUpTo(ProviderAuth, clientAddress, common.HexToHash(podID), tokenAddress, amount, common.Address{})
	if err != nil {
		return err
	}
	return nil
}

func init() {
	monitorCmd.Flags().BoolVarP(&dryRun, "dry-run", "z", false, "Dry run mode; modify nothing.")
	monitorCmd.Flags().StringVar(&kubeConfig, "kubeconfig", "-", "absolute path to the kubeconfig file (- to the first of in-cluster config and ~/.kube/config)")
	monitorCmd.Flags().StringVar(&prometheusUrl, "prometheus", "", "address at which the prometheus API can be accessed")
	monitorCmd.Flags().StringVar(&podID, "id", "00", "pod ID")
	monitorCmd.Flags().StringVar(&ClientAddress, "cAddr", "", "client public address")
	monitorCmd.Flags().StringVar(&TokenContractAddress, "tokenAddr", "", "token contract address")

}
