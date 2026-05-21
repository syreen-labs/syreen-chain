package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"syreen/x/vault/types"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Vault transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		CmdCreateVault(),
		CmdDepositVault(),
		CmdWithdrawVault(),
		CmdCompoundVault(),
		CmdUpdateVaultStrategy(),
	)
	return cmd
}

func CmdCreateVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-vault [name] [deposit-denom] [strategy-type] [target-pool-ids] [performance-fee]",
		Short: "Create a new yield vault",
		Long:  "target-pool-ids is a comma-separated list (e.g. 1,2,3). performance-fee is a decimal (e.g. 0.10 for 10%)",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			var poolIDs []uint64
			for _, s := range strings.Split(args[3], ",") {
				s = strings.TrimSpace(s)
				if s == "" {
					continue
				}
				id, err := strconv.ParseUint(s, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid pool ID: %s", s)
				}
				poolIDs = append(poolIDs, id)
			}
			perfFee, err := math.LegacyNewDecFromStr(args[4])
			if err != nil {
				return fmt.Errorf("invalid performance-fee: %w", err)
			}
			msg := &types.MsgCreateVault{
				Creator:        clientCtx.GetFromAddress().String(),
				Name:           args[0],
				DepositDenom:   args[1],
				StrategyType:   args[2],
				TargetPoolIDs:  poolIDs,
				PerformanceFee: perfFee,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdDepositVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit [vault-id] [amount]",
		Short: "Deposit tokens into a vault",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			vaultID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid vault-id: %w", err)
			}
			amount, ok := math.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("invalid amount: %s", args[1])
			}
			msg := &types.MsgDepositVault{
				Sender:  clientCtx.GetFromAddress().String(),
				VaultID: vaultID,
				Amount:  amount,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdWithdrawVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw [vault-id] [shares]",
		Short: "Withdraw from a vault by burning shares",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			vaultID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid vault-id: %w", err)
			}
			shares, ok := math.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("invalid shares: %s", args[1])
			}
			msg := &types.MsgWithdrawVault{
				Sender:  clientCtx.GetFromAddress().String(),
				VaultID: vaultID,
				Shares:  shares,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdCompoundVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compound [vault-id]",
		Short: "Trigger manual compounding for a vault",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			vaultID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid vault-id: %w", err)
			}
			msg := &types.MsgCompoundVault{
				Sender:  clientCtx.GetFromAddress().String(),
				VaultID: vaultID,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdUpdateVaultStrategy() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-strategy [vault-id] [strategy-type] [target-pool-ids]",
		Short: "Update vault strategy parameters",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			vaultID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid vault-id: %w", err)
			}
			var poolIDs []uint64
			for _, s := range strings.Split(args[2], ",") {
				s = strings.TrimSpace(s)
				if s == "" {
					continue
				}
				id, err := strconv.ParseUint(s, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid pool ID: %s", s)
				}
				poolIDs = append(poolIDs, id)
			}
			msg := &types.MsgUpdateVaultStrategy{
				Creator:       clientCtx.GetFromAddress().String(),
				VaultID:       vaultID,
				StrategyType:  args[1],
				TargetPoolIDs: poolIDs,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
