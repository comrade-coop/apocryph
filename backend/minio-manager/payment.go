package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/comrade-coop/apocryph/backend/abi"
	"github.com/comrade-coop/apocryph/backend/prometheus"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const PERMISSION_WITHDRAW uint16 = 4
const PERMISSION_NO_LIMIT uint16 = 8

var RequiredAllowance = &big.Int{}
var MaximumOverdraft = &big.Int{}

var ChannelDiscriminator = [32]byte{}

func init() {
	discriminatorString := []byte("storage.apocryph.io")
	for i := range discriminatorString {
		ChannelDiscriminator[i] = discriminatorString[i]
	}

	_, _ = RequiredAllowance.SetString("1000000", 10) // 1e6
	_, _ = MaximumOverdraft.SetString("10000000", 10) // 10e6
}

func ExpectedPaymentContractAddress(aappAddress common.Address) common.Address {
	return crypto.CreateAddress(aappAddress, 0)
}

type PaymentManager struct {
	minio          *minio.Client
	minioAdmin     *madmin.AdminClient
	minioMetrics   *madmin.MetricsClient
	ethereum       *ethclient.Client
	tokenAddress   common.Address
	tokenErc20     *abi.IERC20
	paymentAddress common.Address
	payment        *abi.SimplePayment
	aappVersion    uint64
	transactOpts   *bind.TransactOpts
	withdrawTo     common.Address
	prometheus     *prometheus.PrometheusAPI
}

func NewPaymentManager(minioAddress string, minioCreds *credentials.Credentials, ethereumAddress string, tokenAddress common.Address, transactOpts *bind.TransactOpts, withdrawTo common.Address, aappVersion uint64, prometheusClient *prometheus.PrometheusAPI) (*PaymentManager, error) {
	minioClient, err := minio.New(minioAddress, &minio.Options{
		Creds: minioCreds,
	})
	if err != nil {
		return nil, err
	}
	minioAdmin, err := madmin.NewWithOptions(minioAddress, &madmin.Options{
		Creds: minioCreds,
	})
	if err != nil {
		return nil, err
	}
	minioMetrics, err := madmin.NewMetricsClientWithOptions(minioAddress, &madmin.Options{
		Creds: minioCreds,
	})
	if err != nil {
		return nil, err
	}

	ethereumClient, err := ethclient.Dial(ethereumAddress)
	if err != nil {
		return nil, err
	}

	tokenErc20, err := abi.NewIERC20(tokenAddress, ethereumClient)
	if err != nil {
		return nil, err
	}

	return &PaymentManager{
		minio:          minioClient,
		minioAdmin:     minioAdmin,
		minioMetrics:   minioMetrics,
		tokenAddress:   tokenAddress,
		tokenErc20:     tokenErc20,
		paymentAddress: ExpectedPaymentContractAddress(transactOpts.From),
		payment:        nil,
		ethereum:       ethereumClient,
		transactOpts:   transactOpts,
		prometheus:     prometheusClient,
		withdrawTo:     withdrawTo,
	}, nil
}

func (p *PaymentManager) ensureSimplePaymentDeployed(ctx context.Context) (err error) {
	if p.payment != nil {
		return nil
	}

	accountNonce, err := p.ethereum.NonceAt(ctx, p.transactOpts.From, nil)
	if err != nil {
		return
	}
	if accountNonce == 0 { // Fresh account, deploy contract
		var resultAddress common.Address
		resultAddress, _, _, err = abi.DeploySimplePayment(p.transactOpts, p.ethereum, p.tokenAddress)
		if err != nil {
			return fmt.Errorf("Error deploying simple payment contract: %e", err)
		}
		if p.paymentAddress != resultAddress {
			return fmt.Errorf("Invalid payment address for newly-created payment contract; expected %s, got %s", p.paymentAddress, resultAddress)
		}
	}

	paymentContract, err := abi.NewSimplePayment(p.paymentAddress, p.ethereum)
	if err != nil {
		return
	}

	owner, err := paymentContract.Owner(&bind.CallOpts{})
	if err != nil {
		return fmt.Errorf("Failed getting payment contract owner: %e", err)
	}

	if owner != p.transactOpts.From {
		return fmt.Errorf("Invalid payment contract owner; expected %s, got %s", p.transactOpts.From, owner)
	}

	log.Printf("Found payment contract at %s, with correct owner", p.paymentAddress)

	p.payment = paymentContract

	return
}

func (p *PaymentManager) Run(ctx context.Context) error {
	go func() {
		for {
			err := p.ensureSimplePaymentDeployed(ctx)
			if err != nil {
				log.Printf("Error while ensuring payment contract is deployed: %v\n", err)
				log.Printf("Retrying in 30s; deployer address is %s", p.transactOpts.From)
				time.Sleep(30 * time.Second)
				continue
			}
			break
		}

		err := p.reconcilationLoop(ctx)
		if err != nil {
			log.Printf("Error while reconciling payments: %v\n", err)
		}
		time.Sleep(5 * time.Second) // HACK: For dev. purposes, a quick first loop lets us wait on Prometheus starting

		for {
			err := p.reconcilationLoop(ctx)
			if err != nil {
				log.Printf("Error while reconciling payments: %v\n", err)
			}

			time.Sleep(5 * time.Minute)
		}
	}()
	return nil
}

var BucketLabel string = "bucket"

func (p *PaymentManager) reconcilationLoop(ctx context.Context) (err error) {
	bucketByteMinutes, err := p.prometheus.FetchBucketTotalByteMinutes()
	if err != nil {
		return
	}

	priceGbMonth := big.NewInt(int64(0.0025e6))
	minPayment := big.NewInt(int64(0.1e6))

	errs := []error{}

	for bucketId, byteMinutes := range bucketByteMinutes {
		exists, existsErr := p.minio.BucketExists(ctx, bucketId)
		if existsErr == nil && !exists {
			continue
		}
		if !common.IsHexAddress(bucketId) {
			continue
		}

		bucketOwnerAddress := common.HexToAddress(bucketId)

		totalToPay := &big.Int{}
		totalToPay = totalToPay.Mul(byteMinutes, priceGbMonth)
		// totalToPay :: byte * minute * $ / GiB / month
		totalToPay = totalToPay.Div(totalToPay, big.NewInt(30*24*60))
		// totalToPay :: byte * $ / GiB
		totalToPay = totalToPay.Div(totalToPay, big.NewInt(1000*1000))
		// totalToPay :: $

		totalPaid, err := p.getTotalPaidAmount(ctx, bucketOwnerAddress)
		if err != nil {
			errs = append(errs, fmt.Errorf("Failed getting total paid amount: %w", err))
			continue
		}

		paymentAmount := &big.Int{}
		paymentAmount = paymentAmount.Sub(totalToPay, totalPaid)
		fmt.Printf("For bucket %s: total bill so far is %s, paid is %s, difference: %s | min: %s\n", bucketId, totalToPay, totalPaid, paymentAmount, minPayment)

		authorizedAmount, err := p.getAuthorizedAmount(ctx, bucketOwnerAddress)
		if err != nil {
			errs = append(errs, fmt.Errorf("Failed getting authorized amount: %w", err))
			continue
		}

		if authorizedAmount.Cmp(paymentAmount) < 0 {
			overdraft := &big.Int{}
			overdraft.Sub(paymentAmount, authorizedAmount)

			fmt.Printf("For bucket %s: authorized is %s, overdraft is %s | max: %s\n", bucketId, authorizedAmount, overdraft, MaximumOverdraft)

			if overdraft.Cmp(MaximumOverdraft) > 0 {
				fmt.Printf("For bucket %s: deleting bucket that is over overdraft\n", bucketId)
				err = p.removeBucket(ctx, bucketId)
				if err != nil {
					errs = append(errs, fmt.Errorf("Failed deleting bucket %s: %w", bucketId, err))
					continue
				}
				continue
			} else {
				paymentAmount = authorizedAmount // Only withdraw funds we are allowed to
			}
		}
		if paymentAmount.Cmp(minPayment) > 0 && p.withdrawTo != (common.Address{}) {
			tx, err := p.tokenErc20.TransferFrom(p.transactOpts, bucketOwnerAddress, p.withdrawTo, paymentAmount)
			if err != nil {
				errs = append(errs, fmt.Errorf("Failed processing transfer for bucket %s: %w", bucketId, err))
				continue
			}
			fmt.Printf("For bucket %s: processing payment of %s, tx %s\n", bucketId, paymentAmount, tx.Hash())
		}
	}

	log.Printf("Finished reconciling bucket payments\n")

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (p *PaymentManager) getTotalPaidAmount(_ context.Context, bucketId common.Address) (*big.Int, error) {
	totalPaid, err := p.payment.TotalPaid(&bind.CallOpts{
		BlockNumber: nil,
	}, bucketId, p.aappVersion)
	if err != nil {
		return nil, err
	}

	return totalPaid, nil
}

func (p *PaymentManager) getAuthorizedAmount(ctx context.Context, bucketId common.Address) (*big.Int, error) {
	currentBlockNumber, err := p.ethereum.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}

	allowance, err := p.tokenErc20.Allowance(&bind.CallOpts{
		BlockNumber: big.NewInt(int64(currentBlockNumber)),
	}, bucketId, p.paymentAddress)
	if err != nil {
		return nil, err
	}

	balance, err := p.tokenErc20.BalanceOf(&bind.CallOpts{
		BlockNumber: big.NewInt(int64(currentBlockNumber)),
	}, bucketId)
	if err != nil {
		return nil, err
	}

	var authorizedFor *big.Int

	if allowance.Cmp(balance) < 0 {
		authorizedFor = allowance
	} else {
		authorizedFor = balance
	}

	return authorizedFor, nil
}

func (p *PaymentManager) IsAuthorized(ctx context.Context, bucketId common.Address) (bool, error) {
	authorizedFor, err := p.getAuthorizedAmount(ctx, bucketId)
	if err != nil {
		return false, err
	}

	return authorizedFor.Cmp(RequiredAllowance) < 0, nil
}

func (p *PaymentManager) removeBucket(ctx context.Context, bucketId string) error {

	err := p.minio.RemoveBucket(ctx, bucketId)
	if err != nil {
		err2 := p.minio.RemoveBucketWithOptions(ctx, bucketId, minio.RemoveBucketOptions{
			ForceDelete: true,
		})
		if err2 != nil {
			return errors.Join(err, err2)
		}
	}

	return nil
}
