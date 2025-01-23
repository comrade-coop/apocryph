package main

import (
	"context"
	"log"
	"math/big"
	"os"
	"os/signal"

	"github.com/comrade-coop/apocryph/backend/prometheus"
	"github.com/comrade-coop/apocryph/backend/swarm"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	_ "github.com/joho/godotenv/autoload"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/cobra"
)

var identityServeAddress string
var minioAddress string
var minioAccessKey string
var minioSecretKey string
var privateKeyString string
var disableReplication bool
var serfAddress string
var hostname string
var disablePayments bool
var prometheusAddress string
var ethereumAddress string
var chainIdString string
var paymentContractAddress string
var withdrawAddress string

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

		if minioAccessKey == "" {
			minioAccessKey = os.Getenv("ACCESS_KEY")
		}
		if minioSecretKey == "" {
			minioSecretKey = os.Getenv("SECRET_KEY")
		}
		if privateKeyString == "" {
			privateKeyString = os.Getenv("PRIVATE_KEY")
		}
		
		privateKey, err := crypto.HexToECDSA(privateKeyString)
		if err != nil {
			return
		}
	
		replicationSigner, err := NewTokenSigner(privateKey)
		if err != nil {
			return
		}
		minioCreds := credentials.NewStaticV4(minioAccessKey, minioSecretKey, "")
		
		if !disableReplication {
			swarm, err := swarm.NewSwarm(serfAddress, hostname)
			if err != nil {
				return err
			}
			
			replication, err := NewReplicationManager(minioAddress, minioCreds, swarm, replicationSigner)
			if err != nil {
				return err
			}
			err = replication.Run(cmd.Context())
			if err != nil {
				return err
			}
		}
		
		if !disablePayments {
			prometheusClient, err := prometheus.GetPrometheusClient(prometheusAddress)
			if err != nil {
				return err
			}
			paymentAddress := common.HexToAddress(paymentContractAddress)
			withdrawTo := common.HexToAddress(withdrawAddress)
			chainId := &big.Int{}
			chainId, ok := chainId.SetString(chainIdString, 10)
			if !ok {
				return err
			}
			transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
			if err != nil {
				return err
			}
			
			payment, err := NewPaymentManager(minioAddress, minioCreds, ethereumAddress, paymentAddress, transactOpts, withdrawTo, prometheusClient)
			if err != nil {
				return err
			}
			err = payment.Run(cmd.Context())
			if err != nil {
				return err
			}
		}
		
		err = RunIdentityServer(cmd.Context(), identityServeAddress, replicationSigner.GetPublicAddress(), minioCreds)
		if err != nil {
			return err
		}
		return
	},
}

func init() {
	backendCmd.Flags().StringVar(&identityServeAddress, "bind", ":8593", "Bind address to serve the minio identity plugin on")
	backendCmd.Flags().StringVar(&minioAddress, "minio", "localhost:9000", "Address to query minio on")
	
	backendCmd.Flags().StringVar(&minioAccessKey, "minio-access", "", "Access key for Minio (defaults to $ACCESS_KEY from .env)")
	backendCmd.Flags().StringVar(&minioSecretKey, "minio-secret", "", "Secret key for Minio (defaults to $SECRET_KEY from .env)")
	backendCmd.Flags().StringVar(&privateKeyString, "private-key", "", "Private key to use for replication token signing (defaults to $PRIVATE_KEY from .env)")
	
	backendCmd.Flags().BoolVar(&disablePayments, "disable-payments", false, "Disable payments")
	backendCmd.Flags().StringVar(&prometheusAddress, "prometheus", "http://localhost:9090", "Address to query prometheus on")
	backendCmd.Flags().StringVar(&ethereumAddress, "ethereum", "http://localhost:8545", "Address to query ethereum on")
	backendCmd.Flags().StringVar(&paymentContractAddress, "payment-contract", "0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9", "Address of the payment contract")
	backendCmd.Flags().StringVar(&withdrawAddress, "withdraw-address", "", "Address to withdraw to")
	backendCmd.Flags().StringVar(&chainIdString, "chain-id", "31337", "Ethereum Chain ID")
	
	backendCmd.Flags().BoolVar(&disableReplication, "disable-replication", false, "Disable replication")
	backendCmd.Flags().StringVar(&serfAddress, "serf", "localhost:7373", "Address to query serf on")
	backendCmd.Flags().StringVar(&hostname, "hostname", "localhost", "Hostname & local serf node name")
}
