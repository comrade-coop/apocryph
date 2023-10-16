package main

import (
	"fmt"
	"io"
	"os"

	"github.com/comrade-coop/trusted-pods/pkg/prometheus"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/comrade-coop/trusted-pods/pkg/resource"
	"github.com/spf13/cobra"
)

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Operations related to monitoring and pricing pod execution",
}

var prometheusUrl string
var pricingFile string
var pricingFileFormat string

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

func init() {
	metricsCmd.AddCommand(getMetricsCmd)

	getMetricsCmd.Flags().StringVar(&prometheusUrl, "prometheus", "", "address at which the prometheus API can be accessed")
	getMetricsCmd.Flags().StringVar(&pricingFile, "pricing", "", "file containing pricing information")
	getMetricsCmd.Flags().StringVar(&pricingFileFormat, "pricing-format", "json", fmt.Sprintf("pricing file format. one of %v", pb.UnmarshalFormatNames))
}
