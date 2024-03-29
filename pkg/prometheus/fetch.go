// SPDX-License-Identifier: GPL-3.0

package prometheus

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"

	"github.com/comrade-coop/apocryph/pkg/resource"
)

const resourceQuotaSecondsQuery = "sum by (namespace, resource)(sum_over_time(kube_pod_container_resource_requests[250000m]))" // TODO: 250000m might be too short?

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

var sixty = big.NewFloat(60.0)

func (api *PrometheusAPI) FetchResourceMetrics(namespaceConsumptions resource.NamespaceConsumptions) error {
	queryUrl, err := url.JoinPath(api.baseUrl, "/api/v1/query")
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

	for _, qr := range response.Data.Result {
		namespace := qr.Metric["namespace"]
		resourceName := qr.Metric["resource"]
		valueString := qr.Value[1].(string)
		value, _, err := big.ParseFloat(valueString, 10, 6, big.ToNearestEven)
		if err != nil {
			return err
		}
		value = value.Mul(value, sixty)

		namespaceConsumptions.Add(namespace, resource.GetResource(resourceName, resource.ResourceKindReservation), value)
	}

	return nil
}
