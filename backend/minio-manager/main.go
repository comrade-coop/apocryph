package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/comrade-coop/apocryph/backend/prometheus"
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

var disablePayments bool
var prometheusAddress string
var ethereumAddress string
var chainIdString string
var tokenContractAddress string
var withdrawAddress string

var replicateSites []string
var externalUrl string

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
		siweDomain := os.Getenv("GLOBAL_HOST_APP")
		if siweDomain == "" {
			siweDomain = "s3.apocryph.io"
		}
		siweDomainMatch := os.Getenv("SIWE_DOMAIN")
		if siweDomainMatch == "" {
			siweDomainMatch = siweDomain
		}
		aappVersion := os.Getenv("AAPP_VERSION")
		if aappVersion == "" {
			aappVersion = "v0000"
		}
		aappVersionInt, _ := strconv.ParseUint(aappVersion[1:], 10, 64)

		privateKey, err := crypto.HexToECDSA(privateKeyString)
		if err != nil {
			return
		}

		replicationTokenSigner, err := NewTokenSigner(privateKey, siweDomain)
		if err != nil {
			return
		}

		log.Println("Public key for storage system (VITE_STORAGE_SYSTEM): ", replicationTokenSigner.GetPublicAddress())

		minioCreds := credentials.NewStaticV4(minioAccessKey, minioSecretKey, "")

		ownUrl, err := url.Parse(externalUrl)
		if err != nil {
			return
		}

		for _, site := range replicateSites {
			siteUrl, err := url.Parse(site)
			if err != nil {
				return err
			}
			go func() {
				for range 3 {
					// Sleep right away, since bucket replication will need the token signer, and thus identity server running
					time.Sleep(time.Second * 10)
					err = ConfigureAllBucketsReplication(cmd.Context(), ownUrl, siteUrl, replicationTokenSigner)
					if err != nil {
						log.Println("Error configuring all-bucket replication", err)
					} else {
						log.Println("Succefully configured all-bucket replication!")
						return
					}
				}
				log.Fatalln("Failed configuring all-bucket replication!")
			}()
		}

		var payment *PaymentManager = nil

		if !disablePayments {
			prometheusClient, err := prometheus.GetPrometheusClient(prometheusAddress)
			if err != nil {
				return err
			}
			tokenAddress := common.HexToAddress(tokenContractAddress)
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

			payment, err = NewPaymentManager(minioAddress, minioCreds, ethereumAddress, tokenAddress, transactOpts, withdrawTo, aappVersionInt, prometheusClient)
			if err != nil {
				return err
			}
			err = payment.Run(cmd.Context())
			if err != nil {
				return err
			}
		}

		err = RunIdentityServer(cmd.Context(), identityServeAddress, siweDomainMatch, replicationTokenSigner.GetPublicAddress(), minioAddress, minioCreds, payment)
		if err != nil {
			return err
		}
		return
	},
}

var getPublicAddressCmd = &cobra.Command{
	Use: "get-public-address",
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		if privateKeyString == "" {
			privateKeyString = os.Getenv("PRIVATE_KEY")
		}

		privateKey, err := crypto.HexToECDSA(privateKeyString)
		if err != nil {
			return
		}

		fmt.Println(crypto.PubkeyToAddress(privateKey.PublicKey))

		return
	},
}

var getPaymentAddressCmd = &cobra.Command{
	Use: "get-payment-address",
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		if privateKeyString == "" {
			privateKeyString = os.Getenv("PRIVATE_KEY")
		}

		privateKey, err := crypto.HexToECDSA(privateKeyString)
		if err != nil {
			return
		}

		fmt.Println(ExpectedPaymentContractAddress(crypto.PubkeyToAddress(privateKey.PublicKey)))

		return
	},
}

func init() {
	backendCmd.AddCommand(getPublicAddressCmd)
	backendCmd.AddCommand(getPaymentAddressCmd)

	backendCmd.Flags().StringVar(&identityServeAddress, "bind", ":8593", "Bind address to serve the minio identity plugin on")
	backendCmd.Flags().StringVar(&minioAddress, "minio", "localhost:9000", "Address to query minio on")

	backendCmd.Flags().StringVar(&minioAccessKey, "minio-access", "", "Access key for Minio (defaults to $ACCESS_KEY from .env)")
	backendCmd.Flags().StringVar(&minioSecretKey, "minio-secret", "", "Secret key for Minio (defaults to $SECRET_KEY from .env)")
	backendCmd.Flags().StringVar(&privateKeyString, "private-key", "", "Private key to use for replication token signing (defaults to $PRIVATE_KEY from .env)")

	backendCmd.Flags().BoolVar(&disablePayments, "disable-payments", false, "Disable payments")
	backendCmd.Flags().StringVar(&prometheusAddress, "prometheus", "http://localhost:9090", "Address to query prometheus on")
	backendCmd.Flags().StringVar(&ethereumAddress, "ethereum", "http://localhost:8545", "Address to query ethereum on")
	backendCmd.Flags().StringVar(&tokenContractAddress, "token-contract", "", "Address of the token contract")
	backendCmd.Flags().StringVar(&withdrawAddress, "withdraw-address", "", "Address to withdraw to")
	backendCmd.Flags().StringVar(&chainIdString, "chain-id", "31337", "Ethereum Chain ID")

	backendCmd.Flags().StringSliceVar(&replicateSites, "replicate-sites", []string{}, "Replicate to a given site running with the same private key (expects comma-separated http(s):// URLs)")
	backendCmd.Flags().StringVar(&externalUrl, "external-url", "http://localhost", "URL remotes sites can use to reach us")
}
