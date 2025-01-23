// SPDX-License-Identifier: GPL-3.0

package prometheus

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
)

const resourceQuotaSecondsQuery = "sum by (bucket)(sum_over_time(minio_bucket_usage_total_bytes[150000000m]))" // TODO: 250000m might be too short?

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

func (api *PrometheusAPI) FetchBucketTotalByteMinutes() (result map[string]*big.Int, err error) {
	result = map[string]*big.Int{}
	queryUrl, err := url.JoinPath(api.baseUrl, "/api/v1/query")
	if err != nil {
		return
	}

	queryUrl = queryUrl + "?" + url.Values{
		"query": []string{resourceQuotaSecondsQuery},
	}.Encode()

	resp, err := http.Get(queryUrl)
	if err != nil {
		return
	}

	response := QueryResponse{}
	json.NewDecoder(resp.Body).Decode(&response)

	if response.Status != QueryStatusSuccess {
		err = errors.New(fmt.Sprintf("Bad response status: %s", response.Status))
		return
	}

	for _, qr := range response.Data.Result {
		bucket := qr.Metric["bucket"]
		valueString := qr.Value[1].(string)
		value := &big.Int{}
		value, ok := value.SetString(valueString, 10)
		if !ok {
			err = fmt.Errorf("Not a valid number: %s", valueString)
			return
		}
		if result[bucket] != nil {
			result[bucket] = result[bucket].Add(result[bucket], value)
		} else {
			result[bucket] = value
		}
	}

	return
}
