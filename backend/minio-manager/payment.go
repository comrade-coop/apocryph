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

const PERMISSION_WITHDRAW uint16 = 4;
const PERMISSION_NO_LIMIT uint16 = 8;
var RequiredReservation = &big.Int{}

var ChannelDiscriminator = [32]byte{}
func init() {
	discriminatorString := []byte("storage.apocryph.io")
	for i := range discriminatorString {
		ChannelDiscriminator[i] = discriminatorString[i]
	}
	
	_, _ = RequiredReservation.SetString("10000000000000000000", 10) // 10e18
}

type paymentChannelWatch struct {
	channelId      [32]byte
	totalPaid      atomic.Int64
	waitingForTx   atomic.Value
}

type PaymentManager struct {
	minio        *minio.Client
	minioAdmin   *madmin.AdminClient
	minioMetrics *madmin.MetricsClient
	paymentV2    *abi.PaymentV2
	ethereum     *ethclient.Client
	tokenErc20   *abi.IERC20
	transactOpts *bind.TransactOpts
	withdrawTo   common.Address
	prometheus   *prometheus.PrometheusAPI
	
	watches      map[common.Address]*paymentChannelWatch
}

func NewPaymentManager(minioAddress string, minioCreds *credentials.Credentials, ethereumAddress string, paymentAddress common.Address, transactOpts *bind.TransactOpts, withdrawTo common.Address, prometheusClient *prometheus.PrometheusAPI) (*PaymentManager, error) {
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
	
	paymentV2, err := abi.NewPaymentV2(paymentAddress, ethereumClient)
	if err != nil {
		return nil, err
	}
	tokenAddress, err := paymentV2.Token(&bind.CallOpts{})
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
		paymentV2:    paymentV2,
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
	minPayment := big.NewInt(int64(0.004e18))
	
	for bucketId, byteMinutes := range bucketByteMinutes {
		if !common.IsHexAddress(bucketId) {
			continue
		}
		
		totalToPay := &big.Int{}
		totalToPay = totalToPay.Mul(byteMinutes, priceGbMin)
		totalToPay = totalToPay.Div(totalToPay, big.NewInt(1024 * 1024))
		
		watch, err := p.getWatch(ctx, common.HexToAddress(bucketId))
		if err != nil {
			return err
		}
		
		if watch.waitingForTx.Load() == (common.Hash{}) {
			totalPaid := big.NewInt(watch.totalPaid.Load())
			paymentAmount := &big.Int{}
			paymentAmount = paymentAmount.Sub(totalToPay, totalPaid)
			fmt.Printf("For bucket %s: total bill so far is %s, paid is %s, difference: %s\n", bucketId, totalToPay, totalPaid, paymentAmount)
			if paymentAmount.Cmp(minPayment) > 0 && p.withdrawTo != (common.Address{}) {
				tx, err := p.paymentV2.Withdraw(p.transactOpts, watch.channelId, p.withdrawTo, paymentAmount)
				if err != nil {
					return err
				}
				watch.waitingForTx.Store(tx.Hash())
			}
		} else {
			fmt.Printf("For bucket %s: total bill so far is %s, paid is <loading>\n", bucketId, totalToPay)
		}
		
		// TODO: Handle unlocking
		// authorized, err := p.IsAuthorized()
	}
	

	log.Printf("Finished reconciling bucket payments\n")
	return nil
}

func (p *PaymentManager) IsAuthorized(ctx context.Context, bucketId string) (bool, error) {
	bucketOwnerAddress := common.HexToAddress(bucketId)
	channelId, err := p.paymentV2.GetChannelId(&bind.CallOpts{}, bucketOwnerAddress, ChannelDiscriminator)
	if err != nil {
		return false, err
	}
	authorization, err := p.paymentV2.ChannelAuthorizations(&bind.CallOpts{}, channelId, p.transactOpts.From)
	if err != nil {
		return false, err
	}
	
	hasWithdrawPermission := authorization.Permissions & PERMISSION_WITHDRAW == PERMISSION_WITHDRAW
	hasSufficientReservation := authorization.Reservation.Cmp(RequiredReservation) >= 0
	
	return hasWithdrawPermission && hasSufficientReservation, nil
}


func (p *PaymentManager) getWatch(ctx context.Context, bucketId common.Address) (watch *paymentChannelWatch, err error) {
	ownAddress := p.transactOpts.From
	
	if existing, ok := p.watches[bucketId]; ok {
		watch = existing
		return
	}
	
	channelId, err := p.paymentV2.GetChannelId(&bind.CallOpts{}, bucketId, ChannelDiscriminator)
	if err != nil {
		return
	}
	startingBlockNumber, err := p.ethereum.BlockNumber(ctx)
	if err != nil {
		return
	}
	iterator, err := p.paymentV2.FilterWithdraw(
		&bind.FilterOpts{Start: 0, End: &startingBlockNumber},
		[][32]byte{channelId},
		[]common.Address{ownAddress},
	)
	if err != nil {
		return
	}
	events := make(chan *abi.PaymentV2Withdraw)
	subscription, err := p.paymentV2.WatchWithdraw(
		&bind.WatchOpts{Start: &startingBlockNumber, Context: ctx},
		events,
		[][32]byte{channelId},
		[]common.Address{ownAddress},
	)
	if err != nil {
		return
	}
	watch = &paymentChannelWatch{
		channelId: channelId,
	}
	watch.waitingForTx.Store(common.MaxHash) // Block payments until we've synced to current time
	p.watches[bucketId] = watch
	go func() {
		for iterator.Next() {
			event := iterator.Event
			watch.totalPaid.Add(event.Amount.Int64()) // TODO: Ensure this doesn't overflow!
		}
		watch.waitingForTx.CompareAndSwap(common.MaxHash, common.Hash{}) // Done with initial sync
		for {
			select {
			case event := <-events:
				if event.Recipient == ownAddress && event.ChannelId == channelId {
					// TODO: ensure that Add happens before CompareAndSwap here
					watch.totalPaid.Add(event.Amount.Int64())
					watch.waitingForTx.CompareAndSwap(event.Raw.TxHash, common.Hash{})
				}
			case err := <-subscription.Err():
				log.Println("Error in subscription:", err)
				// Will be recreated next reconcilation - TODO use thread-safe map for this
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
