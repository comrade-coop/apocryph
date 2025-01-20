package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"

	"github.com/comrade-coop/apocryph/backend/swarm"
	_ "github.com/joho/godotenv/autoload"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/cobra"
)

var minioAddress string
var serfAddress string
var hostname string
var identityServeAddress string
var dnsServeAddress string
var privateKey string
var minioAccessKey string
var minioSecretKey string

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	go func() {
		<-interruptChan
		cancel()
	}()

	log.Printf("%v", os.Args)

	if err := backendCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

var backendCmd = &cobra.Command{
	Use: "apocryph-s3-backend",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		swarm, err := swarm.NewSwarm(serfAddress, hostname)
		if err != nil {
			return
		}

		replicationSigner, err := NewTokenSigner(privateKey)
		if err != nil {
			return
		}

		if minioAccessKey == "" {
			minioAccessKey = os.Getenv("ACCESS_KEY")
		}
		if minioSecretKey == "" {
			minioSecretKey = os.Getenv("SECRET_KEY")
		}
		if privateKey == "" {
			privateKey = os.Getenv("PRIVATE_KEY")
		}
		minioCreds := credentials.NewStaticV4(minioAccessKey, minioSecretKey, "")

		replication, err := NewReplicationManager(minioAddress, minioCreds, swarm, replicationSigner)
		if err != nil {
			return
		}

		err = errors.Join(
			replication.Run(cmd.Context()),
			RunIdentityServer(cmd.Context(), identityServeAddress, replicationSigner.GetPublicAddress(), minioCreds),
		)
		return
	},
}

func init() {
	backendCmd.Flags().StringVar(&identityServeAddress, "bind", ":8593", "Bind address to serve the minio identity plugin on")
	backendCmd.Flags().StringVar(&minioAddress, "minio", "localhost:9000", "Address to query minio on")
	backendCmd.Flags().StringVar(&minioAccessKey, "minio-access", "", "Access key for Minio (defaults to $ACCESS_KEY from .env)")
	backendCmd.Flags().StringVar(&minioSecretKey, "minio-secret", "", "Secret key for Minio (defaults to $SECRET_KEY from .env)")
	backendCmd.Flags().StringVar(&serfAddress, "serf", "localhost:7373", "Address to query serf on")
	backendCmd.Flags().StringVar(&hostname, "hostname", "localhost", "Hostname & local serf node name")
	backendCmd.Flags().StringVar(&privateKey, "private-key", "", "Private key to use for replication token signing (defaults to $PRIVATE_KEY from .env)")
}
