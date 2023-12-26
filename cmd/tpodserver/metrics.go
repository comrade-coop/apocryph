// SPDX-License-Identifier: GPL-3.0

package main

import (
	"fmt"
	"io"
	"os"

	"github.com/comrade-coop/apocryph/pkg/prometheus"
	pb "github.com/comrade-coop/apocryph/pkg/proto"
	"github.com/comrade-coop/apocryph/pkg/resource"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
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
		namespaceConsumptions := resource.NamespaceConsumptions{}

		err := prometheus.GetPrometheusClient(prometheusUrl).FetchResourceMetrics(namespaceConsumptions)
		if err != nil {
			return err
		}

		pricingTables, err := openPricingTables()
		if err != nil {
			return err
		}

		if len(pricingTables) > 0 {
			for _, ptm := range pricingTables {
				namespaceConsumptions.Display(cmd.OutOrStdout(), ptm)
				fmt.Fprintf(cmd.OutOrStdout(), "totals: %v\n", namespaceConsumptions.Price(ptm))
			}
		} else {
			namespaceConsumptions.Display(cmd.OutOrStdout(), nil)
		}

		return err
	},
}

func openPricingTables() (map[common.Address]resource.PricingTableMap, error) {
	if pricingFileContents != "" {
		pricingTables := &pb.PricingTables{}
		err := pb.Unmarshal(pricingFileFormat, []byte(pricingFileContents), pricingTables)
		if err != nil {
			return nil, err
		}

		return resource.ConvertPricingTables(pricingTables.Tables), nil
	}

	file, err := os.Open(pricingFile)
	if err != nil {
		return nil, err
	}

	pricingTableContents, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	pricingTables := &pb.PricingTables{}
	err = pb.Unmarshal(pricingFileFormat, pricingTableContents, pricingTables)
	if err != nil {
		return nil, err
	}

	return resource.ConvertPricingTables(pricingTables.Tables), nil
}

func init() {
	metricsCmd.AddCommand(getMetricsCmd)

	getMetricsCmd.Flags().StringVar(&prometheusUrl, "prometheus", "", "address at which the prometheus API can be accessed")

	AddConfig("pricing.table.filename", &pricingFile, "", "absolute path to file containing pricing information")
	AddConfig("pricing.table.contents", &pricingFileContents, "", "absolute path to file containing pricing information")
	AddConfig("pricing.table.format", &pricingFileFormat, "", fmt.Sprintf("pricing file format. one of %v", pb.FormatNames))
}
