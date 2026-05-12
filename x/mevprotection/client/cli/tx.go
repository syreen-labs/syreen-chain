package cli

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"syreen/x/mevprotection/types"
)

// GetTxCmd returns the transaction commands for the mevprotection module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "MEV protection transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdCommitTx(),
		CmdRevealTx(),
	)

	return cmd
}

func CmdCommitTx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit-tx [tx-hash-hex] [encrypted-tx-hex]",
		Short: "Submit a commit-phase encrypted transaction",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txHash, err := hex.DecodeString(args[0])
			if err != nil {
				return fmt.Errorf("invalid tx-hash hex: %w", err)
			}

			encryptedTx, err := hex.DecodeString(args[1])
			if err != nil {
				return fmt.Errorf("invalid encrypted-tx hex: %w", err)
			}

			msg := &types.MsgCommitTx{
				Sender:      clientCtx.GetFromAddress().String(),
				TxHash:      txHash,
				EncryptedTx: encryptedTx,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdRevealTx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reveal-tx [committed-hash-hex] [tx-body-hex] [nonce-hex]",
		Short: "Reveal a previously committed transaction",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			commitHash, err := hex.DecodeString(args[0])
			if err != nil {
				return fmt.Errorf("invalid committed-hash hex: %w", err)
			}

			txBody, err := hex.DecodeString(args[1])
			if err != nil {
				return fmt.Errorf("invalid tx-body hex: %w", err)
			}

			nonce, err := hex.DecodeString(args[2])
			if err != nil {
				return fmt.Errorf("invalid nonce hex: %w", err)
			}

			msg := &types.MsgRevealTx{
				Sender:     clientCtx.GetFromAddress().String(),
				CommitHash: commitHash,
				TxBody:     txBody,
				Nonce:      nonce,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
