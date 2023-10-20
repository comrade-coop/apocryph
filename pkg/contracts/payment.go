package contracts

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/comrade-coop/trusted-pods/pkg/abi"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/exp/slices"
)


type PaymentChannelValidator struct {
	ethClient *ethclient.Client
	transactOpts *bind.TransactOpts
	ChainID *big.Int
	allowedContracts []common.Address
	tokenAddress common.Address
	minFunds *big.Int
}

func NewPaymentChannelValidator(ethClient *ethclient.Client, allowedContractAddresses []string, transactOpts *bind.TransactOpts, tokenAddress []byte) (*PaymentChannelValidator, error) {
	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	allowedContracts := make([]common.Address, len(allowedContractAddresses))
	for i, a := range allowedContractAddresses {
		allowedContracts[i] = common.HexToAddress(a)
	}
	return &PaymentChannelValidator{
		ethClient: ethClient,
		transactOpts: transactOpts,
		ChainID: chainID,
		allowedContracts: allowedContracts,
		tokenAddress: common.BytesToAddress(tokenAddress),
		minFunds: big.NewInt(1),
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
	if !slices.Contains(v.allowedContracts, paymentContract) {
		return nil, fmt.Errorf("Invaid payment contract address (expected one of %v)", v.allowedContracts)
	}
	token := common.BytesToAddress(channel.TokenAddress)
	if v.tokenAddress.Cmp(token) != 0 {
		return nil, fmt.Errorf("Wrong token address (expected %v)", v.tokenAddress)
	}
	payment, err := abi.NewPayment(paymentContract, v.ethClient)
	if err != nil {
		return nil, err
	}
	p := &PaymentChannel {
		TransactOpts: v.transactOpts,
		Payment: payment,
		Publisher: common.BytesToAddress(channel.PublisherAddress),
		Token: token,
		PodID: common.BytesToHash(channel.PodID),
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
	Payment *abi.Payment
	TransactOpts *bind.TransactOpts
	Publisher common.Address
	PodID common.Hash
	Token common.Address
}

func (p *PaymentChannel) Available() (*big.Int, error) {
	return p.Payment.Available(&bind.CallOpts{Pending: false}, p.Publisher, p.TransactOpts.From, p.PodID, p.Token)
}

func (p *PaymentChannel) Withdrawn() (*big.Int, error) {
	return p.Payment.Withdrawn(&bind.CallOpts{Pending: false}, p.Publisher, p.TransactOpts.From, p.PodID, p.Token)
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
	return p.WithdrawUpTo(transferAddress, amount)
}

func (p *PaymentChannel) WithdrawUpTo(transferAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return p.Payment.WithdrawUpTo(p.TransactOpts, p.Publisher, p.PodID, p.Token, amount, transferAddress)
}


