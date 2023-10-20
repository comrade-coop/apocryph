package main

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/comrade-coop/trusted-pods/pkg/contracts"
	tpk8s "github.com/comrade-coop/trusted-pods/pkg/kubernetes"
	"github.com/comrade-coop/trusted-pods/pkg/prometheus"
	"github.com/comrade-coop/trusted-pods/pkg/resource"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
)

var withdrawAddressString string
var withdrawTime int64
var withdrawTolerance string

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor the metrics stored in prometheus",
	RunE: func(cmd *cobra.Command, args []string) error {
		pro := prometheus.NewPrometheusAPI(prometheusUrl)

		ethClient, err := contracts.Connect(ethereumRpc)
		if err != nil {
			return err
		}

		pricingTable, err := openPricingTable()
		if err != nil {
			return err
		}

		providerAuth, err := contracts.GetAccount(providerKey, ethClient)
		if err != nil {
			return err
		}

		paymentChannelValidator, err := contracts.NewPaymentChannelValidator(ethClient, allowedContractAddresses, providerAuth, pricingTable.TokenAddress)
		if err != nil {
			return err
		}

		client, err := tpk8s.GetClient(kubeConfig, dryRun)
		if err != nil {
			return err
		}

		withdrawAddress := common.HexToAddress(withdrawAddressString)

		tolerance, _ := (&big.Int{}).SetString(withdrawTolerance, 10)
		if tolerance == nil {
			return errors.New("Invalid tolerance value")
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
			err = client.List(cmd.Context(), namespaces, tpk8s.TrustedPodsNamespaceFilter)
			if err != nil {
				return err
			}

			for _, n := range namespaces.Items {
				if amountOwed, ok := amountsOwed[n.Name]; ok {
					err := func() error {
						paymentChannelProto, err := tpk8s.TrustedPodsNamespaceGetChannel(&n)
						if err != nil {
							return err
						}
						if paymentChannelProto == nil {
							return nil
						}

						paymentChannel, err := paymentChannelValidator.Parse(paymentChannelProto)
						if err != nil {
							// TODO: Stop namespaces that no longer parse (e.g. because they are under minFunds)!
							return err
						}

						tx, err := paymentChannel.WithdrawIfOverMargin(withdrawAddress, amountOwed, tolerance)
						if err != nil {
							return err
						}
						if tx != nil {
							fmt.Printf("namespace %s: Uploaded metrics for %d\n", n.Name, amountOwed)
						}

						return nil
					}()
					if err != nil {
						fmt.Printf("Error while processing namespace %s: %v\n", n.Name, err)
					}
				}
			}
			select {
			case <-cmd.Context().Done():
				break Loop
			case <-time.After(time.Second * time.Duration(withdrawTime)):
				continue
			}
		}

		return nil
	},
}

func init() {
	monitorCmd.Flags().BoolVarP(&dryRun, "dry-run", "z", false, "Dry run mode; modify nothing.")
	monitorCmd.Flags().StringVar(&kubeConfig, "kubeconfig", "-", "absolute path to the kubeconfig file (- to the first of in-cluster config and ~/.kube/config)")
	monitorCmd.Flags().StringVar(&prometheusUrl, "prometheus", "", "address at which the prometheus API can be accessed")
	monitorCmd.Flags().StringVar(&ethereumRpc, "ethereum-rpc", "http://127.0.0.1:8545", "client public address")
	monitorCmd.Flags().StringVar(&providerKey, "ethereum-key", "", "provider account string (private key | http[s]://clef#account | /keystore#account | account (in default keystore))")

	AddConfig("withdraw.address", &withdrawAddressString, "", "ethereum address to withdraw funds to")
	AddConfig("withdraw.tolerace", &withdrawTolerance, "10", "tolerance for withdrawing from address")
	AddConfig("withdraw.time", &withdrawTime, 100, "time in seconds between sucessive billing checks")

}
