package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"syreen/x/lending/types"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: types.ModuleName, Short: "Lending transaction subcommands",
		DisableFlagParsing: true, SuggestionsMinimumDistance: 2, RunE: client.ValidateCmd,
	}
	cmd.AddCommand(CmdDeposit(), CmdWithdraw(), CmdBorrow(), CmdRepay(), CmdLiquidate(), CmdCreateLendingPool())
	return cmd
}

func CmdDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use: "deposit [pool-id] [amount]", Short: "Deposit tokens into a lending pool", Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid pool-id: %w", err) }
			amount, ok := math.NewIntFromString(args[1])
			if !ok { return fmt.Errorf("invalid amount: %s", args[1]) }
			msg := &types.MsgDeposit{Sender: clientCtx.GetFromAddress().String(), PoolID: poolID, Amount: amount}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdWithdraw() *cobra.Command {
	cmd := &cobra.Command{
		Use: "withdraw [pool-id] [amount]", Short: "Withdraw tokens from a lending pool", Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid pool-id: %w", err) }
			amount, ok := math.NewIntFromString(args[1])
			if !ok { return fmt.Errorf("invalid amount: %s", args[1]) }
			msg := &types.MsgWithdraw{Sender: clientCtx.GetFromAddress().String(), PoolID: poolID, Amount: amount}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdBorrow() *cobra.Command {
	cmd := &cobra.Command{
		Use: "borrow [borrow-pool-id] [amount] [collateral-pool-id] [collateral-amount]", Short: "Borrow tokens using collateral", Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			bpid, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid borrow-pool-id: %w", err) }
			amount, ok := math.NewIntFromString(args[1])
			if !ok { return fmt.Errorf("invalid amount: %s", args[1]) }
			cpid, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil { return fmt.Errorf("invalid collateral-pool-id: %w", err) }
			camt, ok := math.NewIntFromString(args[3])
			if !ok { return fmt.Errorf("invalid collateral-amount: %s", args[3]) }
			msg := &types.MsgBorrow{Sender: clientCtx.GetFromAddress().String(), BorrowPoolID: bpid, Amount: amount, CollateralPoolID: cpid, CollateralAmount: camt}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdRepay() *cobra.Command {
	cmd := &cobra.Command{
		Use: "repay [borrow-id] [amount]", Short: "Repay a borrow position", Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			bid, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid borrow-id: %w", err) }
			amount, ok := math.NewIntFromString(args[1])
			if !ok { return fmt.Errorf("invalid amount: %s", args[1]) }
			msg := &types.MsgRepay{Sender: clientCtx.GetFromAddress().String(), BorrowID: bid, Amount: amount}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdLiquidate() *cobra.Command {
	cmd := &cobra.Command{
		Use: "liquidate [borrow-id]", Short: "Liquidate an unhealthy borrow", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			bid, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil { return fmt.Errorf("invalid borrow-id: %w", err) }
			msg := &types.MsgLiquidate{Liquidator: clientCtx.GetFromAddress().String(), BorrowID: bid}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdCreateLendingPool() *cobra.Command {
	cmd := &cobra.Command{
		Use: "create-pool [denom] [collateral-factor] [dex-pool-id] [price-denom]", Short: "Create a lending pool (governance only)", Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil { return err }
			cf, err := math.LegacyNewDecFromStr(args[1])
			if err != nil { return fmt.Errorf("invalid collateral-factor: %w", err) }
			dpid, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil { return fmt.Errorf("invalid dex-pool-id: %w", err) }
			msg := &types.MsgCreateLendingPool{Authority: clientCtx.GetFromAddress().String(), Denom: args[0], CollateralFactor: cf, DexPoolID: dpid, PriceDenom: args[3]}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
