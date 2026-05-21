package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"syreen/x/clmm/types"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "CLMM (concentrated liquidity) transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdCreateCLPool(),
		CmdCreatePosition(),
		CmdAddLiquidity(),
		CmdRemoveLiquidity(),
		CmdCollectFees(),
		CmdCLSwap(),
	)

	return cmd
}

func CmdCreateCLPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pool [denom-a] [denom-b] [tick-spacing] [fee-rate] [initial-price]",
		Short: "Create a new concentrated liquidity pool",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			tickSpacing, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid tick-spacing: %w", err)
			}
			feeRate, err := math.LegacyNewDecFromStr(args[3])
			if err != nil {
				return fmt.Errorf("invalid fee-rate: %w", err)
			}
			initialPrice, err := math.LegacyNewDecFromStr(args[4])
			if err != nil {
				return fmt.Errorf("invalid initial-price: %w", err)
			}
			msg := &types.MsgCreateCLPool{
				Sender:       clientCtx.GetFromAddress().String(),
				DenomA:       args[0],
				DenomB:       args[1],
				TickSpacing:  tickSpacing,
				FeeRate:      feeRate,
				InitialPrice: initialPrice,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdCreatePosition() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-position [pool-id] [tick-lower] [tick-upper] [amount0-desired] [amount1-desired] [amount0-min] [amount1-min]",
		Short: "Create a new LP position in a price range",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid pool-id: %w", err)
			}
			tickLower, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid tick-lower: %w", err)
			}
			tickUpper, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid tick-upper: %w", err)
			}
			amount0Desired, ok := math.NewIntFromString(args[3])
			if !ok {
				return fmt.Errorf("invalid amount0-desired: %s", args[3])
			}
			amount1Desired, ok := math.NewIntFromString(args[4])
			if !ok {
				return fmt.Errorf("invalid amount1-desired: %s", args[4])
			}
			amount0Min, ok := math.NewIntFromString(args[5])
			if !ok {
				return fmt.Errorf("invalid amount0-min: %s", args[5])
			}
			amount1Min, ok := math.NewIntFromString(args[6])
			if !ok {
				return fmt.Errorf("invalid amount1-min: %s", args[6])
			}
			msg := &types.MsgCreatePosition{
				Sender:         clientCtx.GetFromAddress().String(),
				PoolID:         poolID,
				TickLower:      tickLower,
				TickUpper:      tickUpper,
				Amount0Desired: amount0Desired,
				Amount1Desired: amount1Desired,
				Amount0Min:     amount0Min,
				Amount1Min:     amount1Min,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdAddLiquidity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-liquidity [position-id] [amount0-desired] [amount1-desired]",
		Short: "Add more liquidity to an existing position",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			positionID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid position-id: %w", err)
			}
			amount0, ok := math.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("invalid amount0-desired: %s", args[1])
			}
			amount1, ok := math.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid amount1-desired: %s", args[2])
			}
			msg := &types.MsgAddLiquidity{
				Sender:         clientCtx.GetFromAddress().String(),
				PositionID:     positionID,
				Amount0Desired: amount0,
				Amount1Desired: amount1,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdRemoveLiquidity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-liquidity [position-id] [liquidity-amount]",
		Short: "Remove liquidity from a position",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			positionID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid position-id: %w", err)
			}
			liquidityAmount, err := math.LegacyNewDecFromStr(args[1])
			if err != nil {
				return fmt.Errorf("invalid liquidity-amount: %w", err)
			}
			msg := &types.MsgRemoveLiquidity{
				Sender:          clientCtx.GetFromAddress().String(),
				PositionID:      positionID,
				LiquidityAmount: liquidityAmount,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdCollectFees() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collect-fees [position-id]",
		Short: "Collect accrued fees from a position",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			positionID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid position-id: %w", err)
			}
			msg := &types.MsgCollectFees{
				Sender:     clientCtx.GetFromAddress().String(),
				PositionID: positionID,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdCLSwap() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap [pool-id] [denom-in] [amount-in] [min-amount-out]",
		Short: "Swap tokens through a concentrated liquidity pool",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid pool-id: %w", err)
			}
			amountIn, ok := math.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid amount-in: %s", args[2])
			}
			minAmountOut, ok := math.NewIntFromString(args[3])
			if !ok {
				return fmt.Errorf("invalid min-amount-out: %s", args[3])
			}
			msg := &types.MsgCLSwap{
				Sender:       clientCtx.GetFromAddress().String(),
				PoolID:       poolID,
				DenomIn:      args[1],
				AmountIn:     amountIn,
				MinAmountOut: minAmountOut,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
