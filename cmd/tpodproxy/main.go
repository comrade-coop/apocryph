package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"

	tpk8s "github.com/comrade-coop/apocryph/pkg/kubernetes"
	"github.com/spf13/cobra"
)

const SpecialHeaderName = "X-Apocryph-Expected"

// const SpecialCookieName = "X-Apocryph-Expected" // Could support it as a cookie too

const VerificationServiceURL = "http://verification-service.kube-system.svc.cluster.local:8080"

var ValidServiceNames = regexp.MustCompile("^[a-z][a-z0-9-]{0,62}$")

var currentBackingService string                     // e.g. tpod-XXX
var backingServiceSuffix string                      // e.g. .NS.svc.cluster.local
var extraAttestationDetails []tpk8s.AttestationValue // TODO: wire up
var serveAddress = ":9999"

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	go func() {
		<-interruptChan
		cancel()
	}()

	if err := proxyCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

var proxyCmd = &cobra.Command{
	Use:   "tpodproxy [backing service] [backing service suffix] [extra attestation info]",
	Short: "Start a facade routing requests to the right version of an application and serving attestation information",
	Long: "Start a facade routing requests to the right version of an application and serving attestation information.\n" +
		"Example: tpodproxy svc1 .namespace.cluster.local:80 docker.io/org/image@sha256:1234...",
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		currentBackingService = args[0]
		backingServiceSuffix = args[1]
		extraAttestationDetailsJson := args[2]
		err := json.Unmarshal([]byte(extraAttestationDetailsJson), &extraAttestationDetails)
		if err != nil {
			return err
		}
		if !ValidServiceNames.MatchString(currentBackingService) {
			return fmt.Errorf("Invalid value (%s) for backing service, must be a valid service name", currentBackingService)
		}

		mux := http.NewServeMux()
		mux.HandleFunc("GET /.well-known/network.apocryph.attest", attestHandler)
		mux.HandleFunc("/", wildcardHandler)
		s := &http.Server{
			Addr:    serveAddress,
			Handler: mux,
		}
		go func() {
			<-cmd.Context().Done()
			err := s.Shutdown(context.TODO())
			cmd.PrintErr(err)
		}()

		cmd.PrintErr("Server is listening on ", serveAddress)
		err = s.ListenAndServe()
		cmd.PrintErr(err)

		return err
	},
}

func init() {
	proxyCmd.Flags().StringVar(&serveAddress, "address", "", "port to serve on")
}

type attestation struct { // From constellation/verify/server/server.go
	Data []byte `json:"data"`
}
type AttestationResult struct {
	AttestationData  []byte                   `json:"attestation"`
	AttestationError string                   `json:"error"`
	ExtraData        []tpk8s.AttestationValue `json:"extra"`
	Header           string                   `json:"header"`
}

func attestHandler(w http.ResponseWriter, r *http.Request) {
	result := AttestationResult{}

	// TODO: Have the extra attestation data be the user data parameter for validation
	attestationJson, attestationErr := http.Get(VerificationServiceURL + "/?nonce=" + url.QueryEscape(r.URL.Query().Get("nonce")))
	if attestationErr != nil {
		result.AttestationError = attestationErr.Error()
	} else {
		err := json.NewDecoder(attestationJson.Body).Decode(&result.AttestationData)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
	}

	result.ExtraData = extraAttestationDetails
	result.Header = fmt.Sprintf("%s: %s", SpecialHeaderName, currentBackingService)

	// Write the info string to the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}

func wildcardHandler(w http.ResponseWriter, r *http.Request) {
	expectedService := r.Header.Get(SpecialHeaderName)
	if expectedService == "" {
		expectedService = currentBackingService
	}
	if !ValidServiceNames.MatchString(expectedService) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid value for " + SpecialHeaderName))
		return
	}
	proxyRequest := r.Clone(r.Context())
	proxyRequest.URL.Scheme = "http"
	proxyRequest.URL.Host = expectedService + backingServiceSuffix
	proxyRequest.RequestURI = ""
	// TODO: Figure out if we are doing anything about proxyRequest.Host / proxyRequest.Header["Host"]
	proxyResponse, err := http.DefaultClient.Do(proxyRequest)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	if r.Method == http.MethodOptions {
		proxyResponse.Header.Add("Access-Control-Request-Headers", SpecialHeaderName)
	}
	for header, values := range proxyResponse.Header { // HACK: ...there's got to be a simpler way...
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}
	w.WriteHeader(proxyResponse.StatusCode)
	io.Copy(w, proxyResponse.Body)
	for header, values := range proxyResponse.Trailer {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}
}
