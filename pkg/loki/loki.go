// SPDX-License-Identifier: GPL-3.0

package loki

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"connectrpc.com/connect"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/encoding/protojson"
)

const tailPath = "/loki/api/v1/tail"

type LokiConfig struct {
	Host  string
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

type TailData struct {
	Streams []QueryResult `json:"streams"`
}

func GetLogs(namespace, containerName, limit string, url string) ([]*pb.LogEntry, error) {

	requestURL := fmt.Sprintf("%s/api/v1/query_range", url)
	query := fmt.Sprintf("{container=\"%s\",namespace=\"%s\"}", containerName, namespace)
	params := map[string]string{
		"query": query,
		"limit": limit,
	}

	req, err := http.NewRequest("GET", requestURL, nil)
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
		logEntries[i] = &pb.LogEntry{} // inittialize the logEntry or it will cause null pointer dereference (i miss rust)
		if err := protojson.Unmarshal([]byte(line[1]), logEntries[i]); err != nil {
			fmt.Printf("Error unmarshalling entry:%v", err)
			return nil, err
		}

		nanosecondsUnixEpoch, _ := strconv.ParseUint(line[0], 10, 64)
		logEntries[i].NanosecondsUnixEpoch = nanosecondsUnixEpoch
	}
	return logEntries, nil
}

func GetStreamedEntries(namespace, containerName string, srv *connect.ServerStream[pb.PodLogResponse], lokiHost string) error {

	query := fmt.Sprintf("{container=\"%s\",namespace=\"%s\"}", containerName, namespace)
	requestURL := fmt.Sprintf("ws://%s/loki/api/v1/tail?query=%s", lokiHost, query)
	baseURL := &url.URL{
		Scheme: "ws",
		Host:   "loki.loki.svc.cluster.local:3100",
		Path:   tailPath,
	}
	queryValues := baseURL.Query()

	queryValues.Set("query", fmt.Sprintf("{container=\"%s\",namespace=\"%s\"}", containerName, namespace))
	baseURL.RawQuery = queryValues.Encode()
	requestURL = baseURL.String()

	c, _, err := websocket.DefaultDialer.Dial(requestURL, nil)
	if err != nil {
		fmt.Println("Failed Dialing Server:", err)
		return err
	}
	defer c.Close()
	for {
		var logs TailData
		err := c.ReadJSON(&logs)
		if err != nil {
			fmt.Printf("Failed unmarshalling Json: %v \n", err)
			return err
		}
		for _, result := range logs.Streams {
			entries, err := parseEntries(result.Values)
			if err != nil {
				fmt.Println("Failed Parsing Entries:", err)
				return err
			}
			for _, entry := range entries {
				response := pb.PodLogResponse{}
				response.LogEntry = entry
				if err := srv.Send(&response); err != nil {
					fmt.Printf("Error generating response: %v", err)
					return err
				}
			}
		}
	}
}

//
// type LogEntry struct {
// 	NanosecondsUnixEpoch int64  `json:"nanoseconds_unix_epoch"`
// 	Log                  string `json:"log"`
// 	Stream               string `json:"stream"`
// 	Time                 string `json:"time"`
// }
