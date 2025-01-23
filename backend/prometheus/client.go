// SPDX-License-Identifier: GPL-3.0

package prometheus

type PrometheusAPI struct {
	baseUrl string
}

func GetPrometheusClient(baseUrl string) (*PrometheusAPI, error) {
	return &PrometheusAPI{baseUrl: baseUrl}, nil
}
