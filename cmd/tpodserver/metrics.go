package main

import (
	"fmt"
	"io"
	"log"
	"math/big"
	"os"

	"github.com/comrade-coop/trusted-pods/pkg/contracts"
	"github.com/comrade-coop/trusted-pods/pkg/crypto"
	"github.com/comrade-coop/trusted-pods/pkg/prometheus"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/comrade-coop/trusted-pods/pkg/resource"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Operations related to monitoring and pricing pod execution",
}

var prometheusUrl string
var pricingFile string
var pricingFileFormat string
var pricingFileContents string

var getMetricsCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the metrics stored in prometheus",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resourceMeasurements := resource.ResourceMeasurementsMap{}

		err := prometheus.NewPrometheusAPI(prometheusUrl).FetchResourceMetrics(resourceMeasurements)
		if err != nil {
			return err
		}

		pricingTable, err := openPricingTable()
		if err != nil {
			return err
		}
		resourceMeasurements.Display(cmd.OutOrStdout(), pricingTable)
		if pricingTable != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "totals: %v\n", resourceMeasurements.Price(pricingTable))
		}

		return err
	},
}

func openPricingTable() (*pb.PricingTable, error) {
	if pricingFileContents != "" {
		pricingTable := &pb.PricingTable{}
		err := protojson.Unmarshal([]byte(pricingFileContents), pricingTable)
		if err != nil {
			return nil, err
		}

		return pricingTable, nil
	}

	if pricingFile == "" {
		return nil, nil
	}

	file, err := os.Open(pricingFile)
	if err != nil {
		return nil, err
	}

	pricingTableContents, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	pricingTable := &pb.PricingTable{}
	err = pb.Unmarshal(pricingFileFormat, pricingTableContents, pricingTable)
	if err != nil {
		return nil, err
	}

	return pricingTable, nil
}

var uploadMetricsCmd = &cobra.Command{
	Use:   "uploadMetric",
	Short: "Upload metrics",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		clientAddress := common.HexToAddress(ClientAddress)
		tokenAddress := common.HexToAddress(TokenContractAddress)
		paymentAddress := common.HexToAddress(PaymentContractAddress)
		contracts.VerifyContractAddress(PaymentContractAddress, allowedContractCodeHashes)
		units := big.NewInt(Units)
		// derive the Account from the private key
		ks := crypto.CreateKeystore(keystorePath)

		// get an ethclient
		client, err := contracts.Connect(chainRpc)
		if err != nil {
			return err
		}
		_, providerAuth, err := contracts.DeriveAccountConfigs(ProviderKey, Password, exportPassword, client, ks)
		if err != nil {
			return err
		}

		// get a contract instance
		payment, err := contracts.GetContractInstance(client, paymentAddress)

		_, err = contracts.UploadMetrics(providerAuth, payment, clientAddress, big.NewInt(podID), tokenAddress, units)
		if err != nil {
			return err
		}
		log.Println("Uploaded Metrics Successfully")
		return nil
	},
}

func init() {
	metricsCmd.AddCommand(getMetricsCmd)
	metricsCmd.AddCommand(uploadMetricsCmd)

	getMetricsCmd.Flags().StringVar(&prometheusUrl, "prometheus", "", "address at which the prometheus API can be accessed")

	uploadMetricsCmd.Flags().StringVar(&ClientAddress, "cAddr", "", "client public address")
	uploadMetricsCmd.Flags().StringVar(&TokenContractAddress, "tokenAddr", "", "token contract address")
	uploadMetricsCmd.Flags().Int64Var(&Units, "units", 5, "final units of execution")
	uploadMetricsCmd.Flags().Int64Var(&podID, "id", 1, "pod ID")

	AddConfig("pricing.table.filename", &pricingFile, "", "absolute path to file containing pricing information")
	AddConfig("pricing.table.contents", &pricingFileContents, "", "absolute path to file containing pricing information")
	AddConfig("pricing.table.format", &pricingFileFormat, "", fmt.Sprintf("pricing file format. one of %v", pb.UnmarshalFormatNames))
}
