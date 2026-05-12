package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/compute/types"
)

// GetTxCmd returns the transaction commands for the compute module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Compute (smart contracts) transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdStoreCode(),
		CmdInstantiateContract(),
		CmdExecuteContract(),
		CmdMigrateContract(),
	)

	return cmd
}

func CmdStoreCode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "store-code [code-file]",
		Short: "Upload contract code from a WASM file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			wasmBytes, err := os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("failed to read wasm file: %w", err)
			}

			msg := &types.MsgStoreCode{
				Sender:       clientCtx.GetFromAddress().String(),
				WASMByteCode: wasmBytes,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdInstantiateContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instantiate [code-id] [init-msg-json] [label] [optional:admin]",
		Short: "Instantiate a contract from stored code",
		Long: `Instantiate a contract from stored code.
Use --amount flag to send funds to the contract on instantiation.`,
		Args: cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			codeID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid code-id: %w", err)
			}

			initMsg := json.RawMessage(args[1])
			if !json.Valid(initMsg) {
				return fmt.Errorf("invalid init-msg JSON: %s", args[1])
			}

			label := args[2]

			admin := ""
			if len(args) == 4 {
				admin = args[3]
			}

			amountStr, _ := cmd.Flags().GetString("amount")
			var funds sdk.Coins
			if amountStr != "" {
				funds, err = sdk.ParseCoinsNormalized(amountStr)
				if err != nil {
					return fmt.Errorf("invalid amount: %w", err)
				}
			}

			msg := &types.MsgInstantiateContract{
				Sender: clientCtx.GetFromAddress().String(),
				Admin:  admin,
				CodeID: codeID,
				Label:  label,
				Msg:    initMsg,
				Funds:  funds,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("amount", "", "Coins to send to the contract on instantiation")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdExecuteContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute [contract-addr] [msg-json]",
		Short: "Execute a contract message",
		Long: `Execute a message on a smart contract.
Use --amount flag to send funds with the execution.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			contractAddr := args[0]
			execMsg := json.RawMessage(args[1])
			if !json.Valid(execMsg) {
				return fmt.Errorf("invalid msg JSON: %s", args[1])
			}

			amountStr, _ := cmd.Flags().GetString("amount")
			var funds sdk.Coins
			if amountStr != "" {
				funds, err = sdk.ParseCoinsNormalized(amountStr)
				if err != nil {
					return fmt.Errorf("invalid amount: %w", err)
				}
			}

			msg := &types.MsgExecuteContract{
				Sender:   clientCtx.GetFromAddress().String(),
				Contract: contractAddr,
				Msg:      execMsg,
				Funds:    funds,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("amount", "", "Coins to send to the contract with execution")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdMigrateContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate [contract-addr] [new-code-id] [migrate-msg-json]",
		Short: "Migrate a contract to a new code version",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			contractAddr := args[0]

			newCodeID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid new-code-id: %w", err)
			}

			migrateMsg := json.RawMessage(args[2])
			if !json.Valid(migrateMsg) {
				return fmt.Errorf("invalid migrate-msg JSON: %s", args[2])
			}

			msg := &types.MsgMigrateContract{
				Sender:   clientCtx.GetFromAddress().String(),
				Contract: contractAddr,
				CodeID:   newCodeID,
				Msg:      migrateMsg,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
