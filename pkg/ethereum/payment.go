// SPDX-License-Identifier: GPL-3.0

package ethereum

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/comrade-coop/trusted-pods/pkg/abi"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/comrade-coop/trusted-pods/pkg/resource"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type PaymentChannelValidator struct {
	ethClient     *ethclient.Client
	transactOpts  *bind.TransactOpts
	ChainID       *big.Int
	pricingTables map[common.Address]resource.PricingTableMap
	minFunds      *big.Int
}

func NewPaymentChannelValidator(ethClient *ethclient.Client, pricingTables map[common.Address]resource.PricingTableMap, transactOpts *bind.TransactOpts) (*PaymentChannelValidator, error) {
	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	return &PaymentChannelValidator{
		ethClient:     ethClient,
		transactOpts:  transactOpts,
		ChainID:       chainID,
		pricingTables: pricingTables,
		minFunds:      big.NewInt(1),
	}, nil
}

func (v *PaymentChannelValidator) Parse(channel *pb.PaymentChannel) (*PaymentChannel, error) {
	if (&big.Int{}).SetBytes(channel.ChainID).Cmp(v.ChainID) != 0 {
		return nil, fmt.Errorf("Invalid payment channel chain id (expected %s)", common.BigToHash(v.ChainID).Hex())
	}
	provider := v.transactOpts.From
	if common.BytesToAddress(channel.ProviderAddress).Cmp(provider) != 0 {
		return nil, fmt.Errorf("Invalid provider address in payment channel (expected %s)", provider)
	}
	paymentContract := common.BytesToAddress(channel.ContractAddress)
	pricingTable, ok := v.pricingTables[paymentContract]
	if !ok {
		allowedContracts := make([]common.Address, 0, len(v.pricingTables))
		for a := range v.pricingTables {
			allowedContracts = append(allowedContracts, a)
		}
		return nil, fmt.Errorf("Invalid payment contract address (expected one of %v)", allowedContracts)
	}
	payment, err := abi.NewPayment(paymentContract, v.ethClient)
	if err != nil {
		return nil, err
	}
	p := &PaymentChannel{
		TransactOpts: v.transactOpts,
		Payment:      payment,
		Publisher:    common.BytesToAddress(channel.PublisherAddress),
		PodID:        common.BytesToHash(channel.PodID),
		PricingTable: pricingTable,
	}
	available, err := p.Available()
	if err != nil {
		return nil, err
	}
	if available.Cmp(v.minFunds) < 0 {
		return nil, errors.New("Insufficient funds in payment channel")
	}
	return p, nil
}

type PaymentChannel struct {
	Payment      *abi.Payment
	TransactOpts *bind.TransactOpts
	Publisher    common.Address
	PodID        common.Hash
	PricingTable resource.PricingTableMap
}

func (p *PaymentChannel) Available() (*big.Int, error) {
	return p.Payment.Available(&bind.CallOpts{Pending: false}, p.Publisher, p.TransactOpts.From, p.PodID)
}

func (p *PaymentChannel) Withdrawn() (*big.Int, error) {
	return p.Payment.Withdrawn(&bind.CallOpts{Pending: false}, p.Publisher, p.TransactOpts.From, p.PodID)
}

func (p *PaymentChannel) WithdrawIfOverMargin(transferAddress common.Address, amount *big.Int, tolerance *big.Int) (*types.Transaction, error) {
	withdrawn, err := p.Withdrawn()
	if err != nil {
		return nil, err
	}
	claimableAmount := (&big.Int{}).Sub(amount, withdrawn)
	if claimableAmount.Cmp(tolerance) < 0 {
		return nil, nil
	}
	return p.WithdrawUpTo(transferAddress, amount) // TODO: Handle running out of funds!!
}

func (p *PaymentChannel) WithdrawUpTo(transferAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return p.Payment.WithdrawUpTo(p.TransactOpts, p.Publisher, p.PodID, amount, transferAddress)
}
