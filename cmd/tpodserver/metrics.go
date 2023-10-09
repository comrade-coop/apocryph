package main

import (
	"errors"
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

		if pricingFile != "" {
			file, err := os.Open(pricingFile)
			if err != nil {
				return err
			}

			pricingTableContents, err := io.ReadAll(file)
			if err != nil {
				return err
			}

			Unmarshal := formats[pricingFileFormat]
			if Unmarshal == nil {
				return errors.New("Unknown format: " + pricingFileFormat)
			}

			pricingTable := &pb.PricingTable{}

			err = Unmarshal(pricingTableContents, pricingTable)
			if err != nil {
				return err
			}

			resourceMeasurements.Display(cmd.OutOrStdout(), pricingTable)
			fmt.Fprintf(cmd.OutOrStdout(), "totals: %v\n", resourceMeasurements.Price(pricingTable))
		} else {
			resourceMeasurements.Display(cmd.OutOrStdout(), nil)
		}

		return err
	},
}

func init() {
	metricsCmd.AddCommand(getMetricsCmd)

	getMetricsCmd.Flags().StringVar(&prometheusUrl, "prometheus", "", "address at which the prometheus API can be accessed")
	getMetricsCmd.Flags().StringVar(&pricingFile, "pricing", "", "file containing pricing information")

	formatNames := make([]string, 0, len(formats))
	for name := range formats {
		formatNames = append(formatNames, name)
	}
	getMetricsCmd.Flags().StringVar(&pricingFileFormat, "pricing-format", "json", fmt.Sprintf("pricing file format. one of %v", formatNames))
}
