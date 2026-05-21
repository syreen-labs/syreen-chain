package cli

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	"syreen/x/evm/types"
)

// GetTxCmd returns the transaction commands for the EVM module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "EVM transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdDeployContract(),
		CmdCallContract(),
	)

	return cmd
}

// CmdDeployContract deploys a Solidity contract from bytecode
func CmdDeployContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy [bytecode-hex-or-file] [gas-limit]",
		Short: "Deploy an EVM smart contract from bytecode",
		Long: `Deploy an EVM smart contract. Provide the contract bytecode as hex
or a path to a file containing hex bytecode.

Example:
  syreend tx evm deploy 608060405234801561001057600080fd5b5060... 3000000 --from mykey --fees 100000usyreen
  syreend tx evm deploy ./contract.bin 3000000 --from mykey --fees 100000usyreen`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			bytecodeHex := args[0]

			// Check if it's a file path
			if _, statErr := os.Stat(bytecodeHex); statErr == nil {
				data, readErr := os.ReadFile(bytecodeHex)
				if readErr != nil {
					return fmt.Errorf("failed to read bytecode file: %w", readErr)
				}
				bytecodeHex = strings.TrimSpace(string(data))
			}

			// Strip 0x prefix if present
			bytecodeHex = stripHexPrefix(bytecodeHex)

			// Validate hex
			if _, err := hex.DecodeString(bytecodeHex); err != nil {
				return fmt.Errorf("invalid bytecode hex: %w", err)
			}

			gasLimit, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid gas limit: %w", err)
			}

			nonce, err := queryEVMNonce(clientCtx, clientCtx.GetFromAddress().String())
			if err != nil {
				return fmt.Errorf("failed to query EVM nonce: %w", err)
			}

			msg := &types.MsgEthereumTx{
				From:     clientCtx.GetFromAddress().String(),
				To:       "",
				Value:    "0",
				GasLimit: gasLimit,
				Data:     bytecodeHex,
				Nonce:    nonce,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdCallContract calls an EVM contract
func CmdCallContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "call [contract-hex-address] [calldata-hex] [gas-limit] [value]",
		Short: "Call an EVM smart contract",
		Long: `Call a deployed EVM smart contract with encoded calldata.

Example:
  syreend tx evm call 0x1234...abcd a9059cbb000000... 100000 0 --from mykey --fees 100000usyreen`,
		Args: cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			contractAddr := args[0]
			calldataHex := stripHexPrefix(args[1])

			if _, err := hex.DecodeString(calldataHex); err != nil {
				return fmt.Errorf("invalid calldata hex: %w", err)
			}

			gasLimit, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid gas limit: %w", err)
			}

			value := "0"
			if len(args) == 4 {
				value = args[3]
			}

			nonce, err := queryEVMNonce(clientCtx, clientCtx.GetFromAddress().String())
			if err != nil {
				return fmt.Errorf("failed to query EVM nonce: %w", err)
			}

			msg := &types.MsgEthereumTx{
				From:     clientCtx.GetFromAddress().String(),
				To:       contractAddr,
				Value:    value,
				GasLimit: gasLimit,
				Data:     calldataHex,
				Nonce:    nonce,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func stripHexPrefix(s string) string {
	if len(s) >= 2 && (s[:2] == "0x" || s[:2] == "0X") {
		return s[2:]
	}
	return s
}

// queryEVMNonce queries the EVM nonce for a bech32 address from the store
func queryEVMNonce(clientCtx client.Context, bech32Addr string) (uint64, error) {
	addr := common.BytesToAddress(clientCtx.GetFromAddress().Bytes())
	key := append([]byte{types.PrefixNonce}, addr.Bytes()...)
	res, _, err := clientCtx.QueryStore(key, types.StoreKey)
	if err != nil {
		return 0, nil // default to 0 if query fails
	}
	if len(res) == 0 {
		return 0, nil
	}
	return binary.BigEndian.Uint64(res), nil
}
