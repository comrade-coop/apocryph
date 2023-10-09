package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

type QueryStatus string

const (
	QueryStatusSuccess QueryStatus = "success"
)

type QueryResponse struct {
	Status QueryStatus `json:"status"`
	Data   QueryData   `json:"data"`
}

type QueryData struct {
	ResultType string        `json:"resultType"`
	Result     []QueryResult `json:"result"`
}

type QueryResult struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
}

type ResourceUsage map[string]map[string]*big.Float

func (r ResourceUsage) Add(namespace string, resource string, value *big.Float) {
	if r[namespace] == nil {
		r[namespace] = make(map[string]*big.Float, 1)
	}

	if r[namespace][resource] == nil {
		r[namespace][resource] = value
	} else {
		r[namespace][resource].Add(r[namespace][resource], value)
	}
}

func (r ResourceUsage) Display(writer io.Writer) {
	fmt.Fprint(writer, "Resources Used:\n")

	for namespace, resources := range r {
		fmt.Fprintf(writer, " - Namespace: %s\n", namespace)
		for resource, value := range resources {
			fmt.Fprintf(writer, "   - %s : %s\n", resource, value.Text('f', 2))
		}
	}
}

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Operations related to monitoring and pricing pod execution",
}

var prometheusUrl string

const resourceQuotaSecondsQuery = "sum by (namespace, resource)(sum_over_time(kube_pod_container_resource_requests[250000m]))"

var getMetricsCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the metrics stored in prometheus",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {

		queryUrl, err := url.JoinPath(prometheusUrl, "/api/v1/query")
		if err != nil {
			return err
		}

		queryUrl = queryUrl + "?" + url.Values{
			"query": []string{resourceQuotaSecondsQuery},
		}.Encode()

		resp, err := http.Get(queryUrl)
		if err != nil {
			return err
		}

		response := QueryResponse{}
		json.NewDecoder(resp.Body).Decode(&response)

		if response.Status != QueryStatusSuccess {
			return errors.New(fmt.Sprintf("Bad response status: %s", response.Status))
		}

		resourcesUsed := ResourceUsage{}

		for _, qr := range response.Data.Result {
			namespace := qr.Metric["namespace"]
			resource := qr.Metric["resource"] + "QS" // QS -- quota-second
			valueString := qr.Value[1].(string)
			value, _, err := big.ParseFloat(valueString, 10, 5, big.ToNearestEven)
			if err != nil {
				return err
			}

			resourcesUsed.Add(namespace, resource, value)
		}

		resourcesUsed.Display(cmd.OutOrStdout())

		return err
	},
}

func init() {
	metricsCmd.AddCommand(getMetricsCmd)

	getMetricsCmd.Flags().StringVar(&prometheusUrl, "prometheus", "", "address at which the prometheus API can be accessed")
}
