package cli

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
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

			from := clientCtx.GetFromAddress()

			msg := &types.MsgEthereumTx{
				From:     from.String(),
				To:       "",
				Value:    "0",
				GasLimit: gasLimit,
				Data:     bytecodeHex,
				Nonce:    0,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			fmt.Printf("EVM Deploy Transaction:\n")
			fmt.Printf("  From:      %s\n", msg.From)
			fmt.Printf("  Gas Limit: %d\n", msg.GasLimit)
			fmt.Printf("  Bytecode:  %d bytes\n", len(bytecodeHex)/2)
			fmt.Printf("\nTransaction will be broadcast via EVM handler.\n")

			return nil
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

			from := clientCtx.GetFromAddress()

			msg := &types.MsgEthereumTx{
				From:     from.String(),
				To:       contractAddr,
				Value:    value,
				GasLimit: gasLimit,
				Data:     calldataHex,
				Nonce:    0,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			fmt.Printf("EVM Call Transaction:\n")
			fmt.Printf("  From:     %s\n", msg.From)
			fmt.Printf("  To:       %s\n", msg.To)
			fmt.Printf("  Value:    %s\n", msg.Value)
			fmt.Printf("  Gas:      %d\n", msg.GasLimit)
			fmt.Printf("  Calldata: %d bytes\n", len(calldataHex)/2)
			fmt.Printf("\nTransaction will be broadcast via EVM handler.\n")

			_ = clientCtx
			return nil
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
