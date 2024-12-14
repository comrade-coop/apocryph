package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"

	swie "github.com/spruceid/siwe-go"
	"github.com/spf13/cobra"
)

var AuthDomain string = "localhost:5173" // 's3.apocryph.io'
var serveAddress string

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	go func() {
		<-interruptChan
		cancel()
	}()

	if err := backendCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

var backendCmd = &cobra.Command{
	Use:   "tpodstoragebackend",
	RunE: func(cmd *cobra.Command, args []string) error {
		mux := http.NewServeMux()
		mux.HandleFunc("POST /", authenticateHandler)
		s := &http.Server{
			Addr:    serveAddress,
			Handler: mux,
		}
		go func() {
			<-cmd.Context().Done()
			err := s.Shutdown(context.TODO())
			cmd.PrintErr(err)
		}()

		cmd.PrintErrln("Server is listening on ", serveAddress)
		err := s.ListenAndServe()
		cmd.PrintErrln(err)

		return err
	},
}

func init() {
	backendCmd.Flags().StringVar(&serveAddress, "address", ":8593", "port to serve on")
}

type AuthenticationFailure struct {
	Reason string `json:"reason"`
}

type AuthenticationResult struct {
	User string `json:"user"`
	MaxValiditySeconds int `json:"maxValiditySeconds"`
	Claims map[string]interface{} `json:"claims"`
}

type Token struct {
	Message string `json:"message"`
	Signature string `json:"signature"`
}


func authenticateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokenString := r.URL.Query().Get("token")

	token := &Token{}
	err := json.Unmarshal([]byte(tokenString), token)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(AuthenticationFailure {
			Reason: err.Error(),
		})
		return
	}
	message, err := swie.ParseMessage(token.Message)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(AuthenticationFailure {
			Reason: err.Error(),
		})
		return
	}
	println(message.String())

	/*if token.Message.GetResources() != crypto.Keccak256Hash(...) {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(AuthenticationFailure {
			Reason: "Token is for the wrong Bucket!",
		})
		return
	}*/ // TODO

	_, err = message.Verify(token.Signature, &AuthDomain, nil, nil)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(AuthenticationFailure {
			Reason: err.Error(),
		})
		return
	}

	address := message.GetAddress()
	println(hex.EncodeToString(address[:]))

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(AuthenticationResult{
		User: address.Hex(),
		MaxValiditySeconds: 3600, // token.ExpirationTime.Unix() - time.Now().Unix()
		Claims: map[string]interface{}{
			"preferred_username": hex.EncodeToString(address[:]),
		},
	})
}
