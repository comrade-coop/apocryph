package prometheus

type PrometheusAPI struct {
	baseUrl string
}

func GetPrometheusClient(baseUrl string) *PrometheusAPI {
	if baseUrl == "" {
		baseUrl = "http://127.0.0.1:9090/"
	}
	return &PrometheusAPI{baseUrl: baseUrl}
}
