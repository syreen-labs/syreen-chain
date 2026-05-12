package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/dex/types"
)

// GetTxCmd returns the transaction commands for the dex module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "DEX transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdCreatePool(),
		CmdAddLiquidity(),
		CmdRemoveLiquidity(),
		CmdSwap(),
		CmdPlaceOrder(),
		CmdCancelOrder(),
		CmdModifyOrder(),
		CmdCreateReferralCode(),
		CmdRegisterReferral(),
		CmdClaimReferralRewards(),
		CmdFollowTrader(),
		CmdUnfollowTrader(),
		CmdMultiHopSwap(),
		CmdSetPoolFeeConfig(),
	)

	return cmd
}

func CmdCreatePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pool [denom-a] [amount-a] [denom-b] [amount-b]",
		Short: "Create a new liquidity pool with initial reserves",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			denomA := args[0]
			amountA, ok := math.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("invalid amount-a: %s", args[1])
			}

			denomB := args[2]
			amountB, ok := math.NewIntFromString(args[3])
			if !ok {
				return fmt.Errorf("invalid amount-b: %s", args[3])
			}

			msg := &types.MsgCreatePool{
				Sender:  clientCtx.GetFromAddress().String(),
				DenomA:  denomA,
				DenomB:  denomB,
				AmountA: amountA,
				AmountB: amountB,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdAddLiquidity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-liquidity [pool-id] [amount-a] [amount-b] [min-shares]",
		Short: "Add liquidity to an existing pool",
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

			amountA, ok := math.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("invalid amount-a: %s", args[1])
			}

			amountB, ok := math.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid amount-b: %s", args[2])
			}

			minShares, ok := math.NewIntFromString(args[3])
			if !ok {
				return fmt.Errorf("invalid min-shares: %s", args[3])
			}

			msg := &types.MsgAddLiquidity{
				Sender:       clientCtx.GetFromAddress().String(),
				PoolID:       poolID,
				AmountA:      amountA,
				AmountB:      amountB,
				MinSharesOut: minShares,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdRemoveLiquidity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-liquidity [pool-id] [shares] [min-a] [min-b]",
		Short: "Remove liquidity from a pool by burning LP shares",
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

			shares, ok := math.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("invalid shares: %s", args[1])
			}

			minA, ok := math.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid min-a: %s", args[2])
			}

			minB, ok := math.NewIntFromString(args[3])
			if !ok {
				return fmt.Errorf("invalid min-b: %s", args[3])
			}

			msg := &types.MsgRemoveLiquidity{
				Sender:        clientCtx.GetFromAddress().String(),
				PoolID:        poolID,
				SharesIn:      shares,
				MinAmountAOut: minA,
				MinAmountBOut: minB,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdSwap() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap [pool-id] [token-in-denom] [token-in-amount] [min-out]",
		Short: "Swap tokens through a pool",
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

			denom := args[1]
			amount, ok := math.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid token-in-amount: %s", args[2])
			}

			minOut, ok := math.NewIntFromString(args[3])
			if !ok {
				return fmt.Errorf("invalid min-out: %s", args[3])
			}

			msg := &types.MsgSwap{
				Sender:      clientCtx.GetFromAddress().String(),
				PoolID:      poolID,
				TokenIn:     sdk.NewCoin(denom, amount),
				MinTokenOut: minOut,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdPlaceOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "place-order [pool-id] [side] [order-type] [price] [quantity] [time-in-force]",
		Short: "Place a limit/market/stop-loss/take-profit order",
		Long: `Place an order on the order book.

Side: buy, sell
Order types: limit, market, stop_loss, take_profit, stop_loss_limit, take_profit_limit
Time in force: GTC (good-till-cancelled), IOC (immediate-or-cancel), FOK (fill-or-kill)
Price: required for limit orders, use "0" for market orders
Use --trigger-price flag for conditional orders (stop_loss, take_profit, stop_loss_limit, take_profit_limit)`,
		Args: cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid pool-id: %w", err)
			}

			side := args[1]
			orderType := args[2]

			priceStr := args[3]

			quantityStr := args[4]

			timeInForce := args[5]

			triggerPriceStr, _ := cmd.Flags().GetString("trigger-price")

			msg := &types.MsgPlaceOrder{
				Creator:      clientCtx.GetFromAddress().String(),
				PoolID:       poolID,
				Side:         side,
				OrderType:    orderType,
				Price:        priceStr,
				Quantity:     quantityStr,
				TimeInForce:  timeInForce,
				TriggerPrice: triggerPriceStr,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("trigger-price", "", "Trigger price for conditional orders (stop-loss/take-profit)")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdCancelOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-order [order-id]",
		Short: "Cancel an open order",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			orderID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid order-id: %w", err)
			}

			msg := &types.MsgCancelOrder{
				Creator: clientCtx.GetFromAddress().String(),
				OrderID: orderID,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdModifyOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify-order [order-id] [new-price] [new-quantity]",
		Short: "Modify an existing order's price and/or quantity",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			orderID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid order-id: %w", err)
			}

			newPriceStr := args[1]

			newQuantity, ok := math.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid new-quantity: %s", args[2])
			}

			msg := &types.MsgModifyOrder{
				Creator:     clientCtx.GetFromAddress().String(),
				OrderID:     orderID,
				NewPrice:    newPriceStr,
				NewQuantity: newQuantity,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdCreateReferralCode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-referral-code [optional-custom-code]",
		Short: "Create a referral code for your account (auto-generated if not provided)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgCreateReferralCode{
				Creator: clientCtx.GetFromAddress().String(),
			}
			if len(args) > 0 {
				msg.CustomCode = args[0]
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdRegisterReferral() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-referral [referral-code]",
		Short: "Register under a referral code for fee discounts",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgRegisterReferral{
				User:         clientCtx.GetFromAddress().String(),
				ReferralCode: args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdClaimReferralRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-referral-rewards",
		Short: "Claim accumulated referral fee rewards",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgClaimReferralRewards{
				Referrer: clientCtx.GetFromAddress().String(),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdFollowTrader() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "follow-trader [trader-address] [max-per-trade] [total-budget] [copy-ratio-bps]",
		Short: "Follow a trader for copy trading",
		Long:  "Copy ratio is in basis points (1-10000). E.g., 5000 = 50% of trader's position size.",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Validate but pass as strings
			if _, ok := math.NewIntFromString(args[1]); !ok {
				return fmt.Errorf("invalid max-per-trade: %s", args[1])
			}
			if _, ok := math.NewIntFromString(args[2]); !ok {
				return fmt.Errorf("invalid total-budget: %s", args[2])
			}

			copyRatio, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid copy-ratio: %w", err)
			}

			msg := &types.MsgFollowTrader{
				Follower:    clientCtx.GetFromAddress().String(),
				Trader:      args[0],
				MaxPerTrade: args[1],
				TotalBudget: args[2],
				CopyRatio:   copyRatio,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdMultiHopSwap() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multi-hop-swap [route] [token-in-denom] [token-in-amount] [min-out]",
		Short: "Swap tokens through multiple pools atomically",
		Long: `Execute an atomic multi-hop swap through a sequence of pools.
The route is a comma-separated list of pool IDs (e.g., "1,3,5").
Maximum 4 hops allowed.

Example:
  syreend tx dex multi-hop-swap 1,3,5 usyr 1000000 900000 --from mykey`,
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Parse route (comma-separated pool IDs)
			routeStrs := splitRoute(args[0])
			if len(routeStrs) == 0 {
				return fmt.Errorf("route must contain at least one pool ID")
			}
			if len(routeStrs) > 4 {
				return fmt.Errorf("route cannot exceed 4 hops")
			}

			route := make([]uint64, len(routeStrs))
			for i, s := range routeStrs {
				poolID, err := strconv.ParseUint(s, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid pool ID in route at position %d: %s", i, s)
				}
				route[i] = poolID
			}

			denom := args[1]
			if _, ok := math.NewIntFromString(args[2]); !ok {
				return fmt.Errorf("invalid token-in-amount: %s", args[2])
			}
			if _, ok := math.NewIntFromString(args[3]); !ok {
				return fmt.Errorf("invalid min-out: %s", args[3])
			}

			msg := &types.MsgMultiHopSwap{
				Sender:        clientCtx.GetFromAddress().String(),
				Route:         route,
				TokenInDenom:  denom,
				TokenInAmount: args[2],
				MinTokenOut:   args[3],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// splitRoute splits a comma-separated string of pool IDs into a slice.
func splitRoute(s string) []string {
	var parts []string
	current := ""
	for _, c := range s {
		if c == ',' {
			if current != "" {
				parts = append(parts, current)
			}
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

func CmdUnfollowTrader() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unfollow-trader [trader-address]",
		Short: "Stop copy trading a trader",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgUnfollowTrader{
				Follower: clientCtx.GetFromAddress().String(),
				Trader:   args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdSetPoolFeeConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-pool-fee-config [pool-id] [enabled] [base-fee] [max-fee] [volatility-multiplier]",
		Short: "Configure dynamic fee settings for a pool (pool creator only)",
		Long: `Set dynamic fee configuration for a pool. Only the pool creator can call this.

enabled: true/false - whether to use volatility-based dynamic fees
base-fee: minimum fee (decimal, e.g. "0.003" for 0.3%)
max-fee: maximum fee cap (decimal, e.g. "0.05" for 5%)
volatility-multiplier: how much volatility affects fees (decimal, e.g. "1.0")

When enabled, the effective fee = base_fee + (volatility * multiplier), capped at max_fee.
When disabled, the regular pool swap fee is used.`,
		Args: cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid pool-id: %w", err)
			}

			enabled := args[1] == "true"

			// Validate decimal format but pass as strings
			if _, err := math.LegacyNewDecFromStr(args[2]); err != nil {
				return fmt.Errorf("invalid base-fee: %w", err)
			}
			if _, err := math.LegacyNewDecFromStr(args[3]); err != nil {
				return fmt.Errorf("invalid max-fee: %w", err)
			}
			if _, err := math.LegacyNewDecFromStr(args[4]); err != nil {
				return fmt.Errorf("invalid volatility-multiplier: %w", err)
			}

			msg := &types.MsgSetPoolFeeConfig{
				Sender:               clientCtx.GetFromAddress().String(),
				PoolID:               poolID,
				VolatilityFeeEnabled: enabled,
				BaseFee:              args[2],
				MaxFee:               args[3],
				VolatilityMultiplier: args[4],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
