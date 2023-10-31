package loki

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

type LokiConfig struct {
	Url   string
	Limit string
}

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
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

func GetLogs(containerName string, limit string, url string) ([]*pb.LogEntry, error) {

	query := fmt.Sprintf("{container=\"%s\"}", containerName)
	params := map[string]string{
		"query": query,
		"limit": limit,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("failed sending request to loki: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body:%v", err)
		return nil, err
	}

	response := QueryResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("Error unmarshalling body:%v", err)
		return nil, err
	}
	if response.Status != QueryStatusSuccess {
		return nil, errors.New(fmt.Sprintf("Bad response status: %s", response.Status))
	}

	var logEntries []*pb.LogEntry
	for _, qr := range response.Data.Result {
		logEntries, _ = parseEntries(qr.Values)
	}

	return logEntries, nil
}

func parseEntries(lines [][]string) ([]*pb.LogEntry, error) {
	logEntries := make([]*pb.LogEntry, len(lines))
	for i, line := range lines {
		logEntries[i] = &pb.LogEntry{} // inittialize the logEntry or it will cause null pointer derefrence (i miss rust)
		if err := protojson.Unmarshal([]byte(line[1]), logEntries[i]); err != nil {
			fmt.Printf("Error unmarshalling entry:%v", err)
			return nil, err
		}

		nanosecondsUnixEpoch, _ := strconv.ParseUint(line[0], 10, 64)
		logEntries[i].NanosecondsUnixEpoch = nanosecondsUnixEpoch
	}
	return logEntries, nil
}

// type LogEntry struct {
// 	NanosecondsUnixEpoch int64  `json:"nanoseconds_unix_epoch"`
// 	Log                  string `json:"log"`
// 	Stream               string `json:"stream"`
// 	Time                 string `json:"time"`
// }
