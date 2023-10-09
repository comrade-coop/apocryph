package prometheus

type PrometheusAPI struct {
	baseUrl string
}

func NewPrometheusAPI(baseUrl string) PrometheusAPI {
	return PrometheusAPI{baseUrl: baseUrl}
}
