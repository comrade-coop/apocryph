package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/comrade-coop/apocryph/backend/abi"
	"github.com/comrade-coop/apocryph/backend/prometheus"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const PERMISSION_WITHDRAW uint16 = 4
const PERMISSION_NO_LIMIT uint16 = 8

var RequiredAllowance = &big.Int{}

var ChannelDiscriminator = [32]byte{}

func init() {
	discriminatorString := []byte("storage.apocryph.io")
	for i := range discriminatorString {
		ChannelDiscriminator[i] = discriminatorString[i]
	}

	_, _ = RequiredAllowance.SetString("10000000000000000000", 10) // 10e18
}

type paymentChannelWatch struct {
	totalPaid    atomic.Value
	allowance    atomic.Value
	waitingForTx atomic.Value
}

type PaymentManager struct {
	minio        *minio.Client
	minioAdmin   *madmin.AdminClient
	minioMetrics *madmin.MetricsClient
	ethereum     *ethclient.Client
	tokenErc20   *abi.IERC20
	transactOpts *bind.TransactOpts
	withdrawTo   common.Address
	prometheus   *prometheus.PrometheusAPI

	watches map[common.Address]*paymentChannelWatch
}

func NewPaymentManager(minioAddress string, minioCreds *credentials.Credentials, ethereumAddress string, tokenAddress common.Address, transactOpts *bind.TransactOpts, withdrawTo common.Address, prometheusClient *prometheus.PrometheusAPI) (*PaymentManager, error) {
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
		minio:        minioClient,
		minioAdmin:   minioAdmin,
		minioMetrics: minioMetrics,
		tokenErc20:   tokenErc20,
		ethereum:     ethereumClient,
		transactOpts: transactOpts,
		prometheus:   prometheusClient,
		withdrawTo:   withdrawTo,
		watches:      map[common.Address]*paymentChannelWatch{},
	}, nil
}

func (p *PaymentManager) Run(ctx context.Context) error {
	go func() {
		err := p.reconcilationLoop(ctx)
		if err != nil {
			log.Printf("Error while reconciling payments: %v\n", err)
		}
		time.Sleep(5 * time.Second) // HACK: For dev. purposes, a quick first loop lets us load all watches
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

	priceGbMin := big.NewInt(int64(0.000004e18) * 60)
	minPayment := big.NewInt(int64(0.1e18))

	for bucketId, byteMinutes := range bucketByteMinutes {
		if !common.IsHexAddress(bucketId) {
			continue
		}

		bucketOwnerAddress := common.HexToAddress(bucketId)

		totalToPay := &big.Int{}
		totalToPay = totalToPay.Mul(byteMinutes, priceGbMin)
		totalToPay = totalToPay.Div(totalToPay, big.NewInt(1024*1024))

		watch, err := p.getWatch(ctx, bucketOwnerAddress)
		if err != nil {
			return err
		}

		authorized, err := p.IsAuthorized(ctx, bucketId)
		if err != nil {
			return err
		}
		if !authorized {
			fmt.Printf("For bucket %s: we are not authorized (has: %s, required: %s)!\n", bucketId, watch.allowance.Load(), RequiredAllowance)
			continue // TODO: should remove bucket
		}

		if watch.waitingForTx.Load() == (common.Hash{}) {
			totalPaid := watch.totalPaid.Load().(*big.Int)
			paymentAmount := &big.Int{}
			paymentAmount = paymentAmount.Sub(totalToPay, totalPaid)
			fmt.Printf("For bucket %s: total bill so far is %s, paid is %s, difference: %s | min: %s\n", bucketId, totalToPay, totalPaid, paymentAmount, minPayment)
			if paymentAmount.Cmp(minPayment) > 0 && p.withdrawTo != (common.Address{}) {
				tx, err := p.tokenErc20.TransferFrom(p.transactOpts, bucketOwnerAddress, p.withdrawTo, paymentAmount)
				if err != nil {
					return err // TODO: Single error blocks all payments to latter buckets
				}
				watch.waitingForTx.Store(tx.Hash())
				fmt.Printf("For bucket %s: processing payment of %s, tx %s\n", bucketId, paymentAmount, tx.Hash())
			}
		} else {
			fmt.Printf("For bucket %s: total bill so far is %s, paid is <loading>\n", bucketId, totalToPay)
		}
	}

	log.Printf("Finished reconciling bucket payments\n")
	return nil
}

func (p *PaymentManager) IsAuthorized(ctx context.Context, bucketId string) (bool, error) {
	bucketOwnerAddress := common.HexToAddress(bucketId)

	watch, err := p.getWatch(ctx, bucketOwnerAddress)
	if err != nil {
		return false, err
	}

	if watch.allowance.Load().(*big.Int).Cmp(RequiredAllowance) < 0 {
		return false, nil
	}

	return true, nil
}

func (p *PaymentManager) getWatch(ctx context.Context, bucketId common.Address) (watch *paymentChannelWatch, err error) {
	ownAddress := p.transactOpts.From

	if existing, ok := p.watches[bucketId]; ok {
		watch = existing
		return
	}

	startingBlockNumber, err := p.ethereum.BlockNumber(ctx)
	if err != nil {
		return
	}

	allowance, err := p.tokenErc20.Allowance(&bind.CallOpts{
		BlockNumber: big.NewInt(int64(startingBlockNumber)),
	}, bucketId, p.transactOpts.From)
	if err != nil {
		return
	}

	iterator, err := p.tokenErc20.FilterTransfer(
		&bind.FilterOpts{Start: 0, End: &startingBlockNumber},
		[]common.Address{bucketId},
		[]common.Address{p.withdrawTo}, // NOTE: Multiple managers with the same withdraw address _will_ get confused about payments
	)
	if err != nil {
		return
	}

	events := make(chan *abi.IERC20Transfer)
	subscription, err := p.tokenErc20.WatchTransfer(
		&bind.WatchOpts{Start: &startingBlockNumber, Context: ctx},
		events,
		[]common.Address{bucketId},
		[]common.Address{p.withdrawTo},
	)
	if err != nil {
		return
	}

	allowanceEvents := make(chan *abi.IERC20Approval)
	allowanceSubscription, err := p.tokenErc20.WatchApproval(
		&bind.WatchOpts{Start: &startingBlockNumber, Context: ctx},
		allowanceEvents,
		[]common.Address{bucketId},
		[]common.Address{ownAddress},
	)
	if err != nil {
		return
	}
	watch = &paymentChannelWatch{}
	watch.allowance.Store(allowance)         // Block payments until we've synced to current time
	watch.waitingForTx.Store(common.MaxHash) // Block payments until we've synced to current time
	p.watches[bucketId] = watch

	totalPaid := &big.Int{}

	go func() {
		for iterator.Next() {
			event := iterator.Event
			newTotalPaid := &big.Int{}
			newTotalPaid.Add(totalPaid, event.Value)
			totalPaid = newTotalPaid
			watch.totalPaid.Store(totalPaid)
		}
		watch.waitingForTx.CompareAndSwap(common.MaxHash, common.Hash{}) // Done with initial sync
		for {
			select {
			case event := <-events:
				if event.To == p.withdrawTo && event.From == bucketId {
					newTotalPaid := &big.Int{}
					newTotalPaid.Add(totalPaid, event.Value)
					totalPaid = newTotalPaid
					watch.totalPaid.Store(totalPaid)
					watch.waitingForTx.CompareAndSwap(event.Raw.TxHash, common.Hash{})
				}
			case event := <-allowanceEvents:
				if event.Spender == ownAddress && event.Owner == bucketId {
					watch.allowance.Store(event.Value)
				}
			case err := <-subscription.Err():
				log.Println("Error in subscription:", err)
				// Will be recreated next reconcilation - TODO use thread-safe map for this
				p.watches[bucketId] = nil
				return
			case err := <-allowanceSubscription.Err():
				log.Println("Error in allowance subscription:", err)
				// TODO: same as above
				p.watches[bucketId] = nil
				return
			case <-ctx.Done():
				p.watches[bucketId] = nil
				return
			}
		}
	}()

	return
}
