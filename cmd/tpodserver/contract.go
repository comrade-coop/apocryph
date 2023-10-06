package main

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	tpgsrpc "github.com/comrade-coop/trusted-pods/pkg/substrate"
	tptypes "github.com/comrade-coop/trusted-pods/pkg/substrate/types"
	"github.com/spf13/cobra"
)

var contractCmd = &cobra.Command{
	Use:   "contract",
	Short: "Operations related to managing contracts",
}

var chainRpc string
var substrateKey string
var allowedContractCodeHashes []string

const (
	GetSelector   tptypes.ContractSelector = 0x2f865bd9
	ClaimSelector tptypes.ContractSelector = 0xb388803f
)

func ValidateCodeHash(codeHash types.Hash) error {
	codeHashHex := codeHash.Hex()

	for _, v := range allowedContractCodeHashes {

		if v == codeHashHex {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Contract code hash (%s) not in the list of allowed code hashes (in config)", codeHashHex))
}

var checkContractCmd = &cobra.Command{
	Use:   "check <contract>",
	Short: "Check whether a payment contract has a permissible hash",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		_, contractAddress, err := tptypes.NewAccountIDFromSS58(args[0])
		if err != nil {
			return err
		}

		api, err := tpgsrpc.NewSubstrateAPI(chainRpc)
		if err != nil {
			return err
		}

		block, err := api.RPC.Chain.GetBlockHashLatest() // GetFinalizedHead
		if err != nil {
			return err
		}

		contractHash, err := api.RPC.Contracts.GetCodeHash(*contractAddress, &block)
		if err != nil {
			return err
		}

		err = ValidateCodeHash(contractHash)
		if err != nil {
			return err
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), "Contract code is valid")
		}

		inputData, err := tptypes.EncodeInputData(GetSelector)
		if err != nil {
			return err
		}

		execResult, err := api.RPC.Contracts.QueryContract(*contractAddress, *contractAddress, inputData, types.NewU128(*big.NewInt(0)), &block)
		if err != nil {
			return err
		}
		var returnValue tptypes.Result[types.UCompact, tptypes.LangError]
		err = execResult.DecodeResult(&returnValue)
		if err != nil {
			return err
		}

		if returnValue.IsError {
			return fmt.Errorf("Contract execution error: %v", returnValue.Error)
		}

		funds := big.Int(returnValue.Value)

		fmt.Fprintf(cmd.OutOrStdout(), "Available funds in contract: %v\n", &funds)
		return nil
	},
}

var claimContractCmd = &cobra.Command{
	Use:   "claim <contract> <amount>",
	Short: "Claim funds from a payment contract",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		_, contractAddress, err := tptypes.NewAccountIDFromSS58(args[0])
		if err != nil {
			return err
		}

		amount, ok := (&big.Int{}).SetString(args[1], 10)
		if !ok {
			return errors.New("Could not parse amount")
		}

		api, err := tpgsrpc.NewSubstrateAPI(chainRpc)
		if err != nil {
			return err
		}

		properties, err := api.RPC.System.Properties()
		if err != nil {
			return err
		}

		from, err := signature.KeyringPairFromSecret(substrateKey, uint16(properties.AsSS58Format))
		if err != nil {
			return err
		}

		inputData, err := tptypes.EncodeInputData(ClaimSelector, types.UCompact(*amount))
		if err != nil {
			return err
		}

		subscription, err := api.RPC.Contracts.CallContract(*contractAddress, from, inputData, types.NewU128(*big.NewInt(0)))

		defer subscription.Unsubscribe()
		timeoutChan := time.After(60 * time.Second)
	Loop:
		for {
			select {
			case status := <-subscription.Chan():
				if status.IsInBlock {
					break Loop
				}
			case <-timeoutChan:
				return errors.New("Timed out waiting for claim transaction to be accepted")
			}
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Successfully claimed funds from contract: %d\n", amount)

		return nil
	},
}

func init() {
	contractCmd.AddCommand(checkContractCmd)
	contractCmd.AddCommand(claimContractCmd)

	contractCmd.PersistentFlags().StringVar(&chainRpc, "rpc", "", "Link to the Substrate RPC.")

	AddConfig("contract.substrateKey", &substrateKey, []string{"0x0b7694d1042fb3ff479212be20ec1d46507cf46773a91fac4952e755e9add1a8"}, "Key for signing Substrate transactions.")
	AddConfig("contract.allowedPaymentContractCodeHashes", &allowedContractCodeHashes, []string{"0xf0cc3415b17e718eeb8cd0e9dea22bb856df14414937ac6fdb666b3148b795f5"}, "Allowed hashes for contracts.")
}
