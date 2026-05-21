package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"syreen/x/flashloan/types"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: types.ModuleName, Short: "Flashloan transaction subcommands",
		DisableFlagParsing: true, SuggestionsMinimumDistance: 2, RunE: client.ValidateCmd,
	}
	cmd.AddCommand(CmdFlashLoan(), CmdCreateFlashPool(), CmdFundFlashPool(), CmdWithdrawFlashPool())
	return cmd
}

func CmdFlashLoan() *cobra.Command {
	cmd := &cobra.Command{
		Use: "flash-loan [denom] [amount]", Short: "Execute a flash loan", Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			amount, ok := math.NewIntFromString(args[1])
			if !ok { return fmt.Errorf("invalid amount: %s", args[1]) }
			msg := &types.MsgFlashLoan{Sender: clientCtx.GetFromAddress().String(), Denom: args[0], Amount: amount}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdCreateFlashPool() *cobra.Command {
	cmd := &cobra.Command{
		Use: "create-pool [denom] [fee-rate]", Short: "Create a flash loan pool (authority only)", Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			feeRate, err := math.LegacyNewDecFromStr(args[1])
			if err != nil { return fmt.Errorf("invalid fee-rate: %w", err) }
			msg := &types.MsgCreateFlashPool{Authority: clientCtx.GetFromAddress().String(), Denom: args[0], FeeRate: feeRate}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdFundFlashPool() *cobra.Command {
	cmd := &cobra.Command{
		Use: "fund-pool [denom] [amount]", Short: "Fund a flash loan pool", Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			amount, ok := math.NewIntFromString(args[1])
			if !ok { return fmt.Errorf("invalid amount: %s", args[1]) }
			msg := &types.MsgFundFlashPool{Sender: clientCtx.GetFromAddress().String(), Denom: args[0], Amount: amount}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdWithdrawFlashPool() *cobra.Command {
	cmd := &cobra.Command{
		Use: "withdraw-pool [denom] [shares]", Short: "Withdraw from a flash loan pool", Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			shares, ok := math.NewIntFromString(args[1])
			if !ok { return fmt.Errorf("invalid shares: %s", args[1]) }
			msg := &types.MsgWithdrawFlashPool{Sender: clientCtx.GetFromAddress().String(), Denom: args[0], Shares: shares}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
