package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"cosmossdk.io/math"

	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/dex/types"
)

// GetQueryCmd returns the query commands for the dex module.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the dex module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryParams(),
		CmdQueryPool(),
		CmdQueryPools(),
		CmdQueryQuote(),
		CmdQuerySignals(),
		CmdQuerySignalHistory(),
		CmdQueryRiskScore(),
		CmdQueryTradeRisk(),
		CmdQueryRiskHistory(),
		CmdQueryOptimalRoute(),
		CmdQueryRoutes(),
		CmdQuerySentiment(),
		CmdQuerySentimentHistory(),
		CmdQuerySentimentAlerts(),
		CmdQueryOracle(),
		CmdQueryWhaleAlerts(),
		CmdQueryTimeWeightedPower(),
		CmdQueryTraderRecord(),
		CmdQueryTraderLeaderboard(),
		CmdQueryIBCOrders(),
		CmdQueryPoolFee(),
	)

	return cmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query dex module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryStore([]byte("params"), types.StoreKey)
			if err != nil {
				return fmt.Errorf("failed to query params: %w", err)
			}

			if len(res) == 0 {
				params := types.DefaultParams()
				bz, _ := json.MarshalIndent(params, "", "  ")
				return clientCtx.PrintBytes(bz)
			}

			var params types.Params
			if err := json.Unmarshal(res, &params); err != nil {
				return fmt.Errorf("failed to unmarshal params: %w", err)
			}

			bz, _ := json.MarshalIndent(params, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool [pool-id]",
		Short: "Query a pool by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid pool-id: %w", err)
			}

			key := types.PoolKey(poolID)
			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil {
				return fmt.Errorf("failed to query pool: %w", err)
			}

			if len(res) == 0 {
				return fmt.Errorf("pool %d not found", poolID)
			}

			var pool types.Pool
			if err := json.Unmarshal(res, &pool); err != nil {
				return clientCtx.PrintBytes(res)
			}

			bz, _ := json.MarshalIndent(pool, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryPools() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pools",
		Short: "Query all pools",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// Iterate pool IDs 1..100 (practical upper bound for CLI)
			var pools []types.Pool
			for id := uint64(1); id <= 100; id++ {
				key := types.PoolKey(id)
				res, _, err := clientCtx.QueryStore(key, types.StoreKey)
				if err != nil || len(res) == 0 {
					break
				}
				var pool types.Pool
				if err := json.Unmarshal(res, &pool); err != nil {
					continue
				}
				pools = append(pools, pool)
			}

			if len(pools) == 0 {
				bz, _ := json.MarshalIndent([]types.Pool{}, "", "  ")
				return clientCtx.PrintBytes(bz)
			}

			bz, _ := json.MarshalIndent(pools, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryQuote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quote [pool-id] [token-in-denom] [token-in-amount]",
		Short: "Get a swap quote (read-only)",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
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

			key := types.PoolKey(poolID)
			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil {
				return fmt.Errorf("failed to query pool: %w", err)
			}

			if len(res) == 0 {
				return fmt.Errorf("pool %d not found", poolID)
			}

			var pool types.Pool
			if err := json.Unmarshal(res, &pool); err != nil {
				return fmt.Errorf("failed to unmarshal pool: %w", err)
			}

			tokenIn := sdk.NewCoin(denom, amount)

			var reserveIn, reserveOut math.Int
			var denomOut string

			if tokenIn.Denom == pool.DenomA {
				reserveIn = pool.ReserveA
				reserveOut = pool.ReserveB
				denomOut = pool.DenomB
			} else if tokenIn.Denom == pool.DenomB {
				reserveIn = pool.ReserveB
				reserveOut = pool.ReserveA
				denomOut = pool.DenomA
			} else {
				return fmt.Errorf("denom %s is not in pool %d", denom, poolID)
			}

			feeBps := pool.SwapFee.MulInt64(10000).TruncateInt64()
			amountInAfterFee := tokenIn.Amount.MulRaw(10000 - feeBps).QuoRaw(10000)
			outputAmount := reserveOut.Mul(amountInAfterFee).Quo(reserveIn.Add(amountInAfterFee))

			result := map[string]interface{}{
				"pool_id":      poolID,
				"token_in":     tokenIn.String(),
				"expected_out": sdk.NewCoin(denomOut, outputAmount).String(),
				"reserve_in":   reserveIn.String(),
				"reserve_out":  reserveOut.String(),
				"swap_fee":     pool.SwapFee.String(),
			}

			bz, _ := json.MarshalIndent(result, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// ---------------------------------------------------------------------------
// Signals CLI Commands
// ---------------------------------------------------------------------------

func CmdQuerySignals() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signals [pool-id]",
		Short: "Query trading signals for a pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			_ = args
			return clientCtx.PrintString("Use the gRPC/REST endpoint: /syreen/dex/v1/signals\n")
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQuerySignalHistory() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signal-history [pool-id]",
		Short: "Query signal history for a pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			_ = args
			return clientCtx.PrintString("Use the gRPC/REST endpoint: /syreen/dex/v1/signal_history\n")
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// ---------------------------------------------------------------------------
// Risk CLI Commands
// ---------------------------------------------------------------------------

func CmdQueryRiskScore() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "risk-score [pool-id]",
		Short: "Query the risk score for a pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			_ = args
			return clientCtx.PrintString("Use the gRPC/REST endpoint: /syreen/dex/v1/risk_score\n")
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryTradeRisk() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trade-risk [pool-id] [input-denom] [amount]",
		Short: "Assess risk for a potential trade",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			_ = args
			return clientCtx.PrintString("Use the gRPC/REST endpoint: /syreen/dex/v1/trade_risk\n")
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryRiskHistory() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "risk-history [pool-id]",
		Short: "Query risk score history for a pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			_ = args
			return clientCtx.PrintString("Use the gRPC/REST endpoint: /syreen/dex/v1/risk_history\n")
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// ---------------------------------------------------------------------------
// Router CLI Commands
// ---------------------------------------------------------------------------

func CmdQueryOptimalRoute() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "optimal-route [input-denom] [output-denom] [amount]",
		Short: "Find optimal swap route between two denoms",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			_ = args
			return clientCtx.PrintString("Use the gRPC/REST endpoint: /syreen/dex/v1/optimal_route\n")
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryRoutes() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "routes [input-denom] [output-denom]",
		Short: "List all available swap routes between two denoms",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			_ = args
			return clientCtx.PrintString("Use the gRPC/REST endpoint: /syreen/dex/v1/routes\n")
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// ---------------------------------------------------------------------------
// Sentiment Oracle CLI Commands
// ---------------------------------------------------------------------------

func CmdQuerySentiment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sentiment [pool-id]",
		Short: "Query sentiment oracle data for a pool (Fear & Greed index, mood, indicators)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid pool-id: %w", err)
			}
			key := types.SentimentStateKey(poolID)
			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil || len(res) == 0 {
				result := map[string]interface{}{
					"pool_id":          poolID,
					"fear_greed_index": 50,
					"fear_greed_label": "Neutral",
					"mood":             "NEUTRAL",
					"message":          "No sentiment data yet",
				}
				bz, _ := json.MarshalIndent(result, "", "  ")
				return clientCtx.PrintBytes(bz)
			}
			var parsed interface{}
			if err := json.Unmarshal(res, &parsed); err != nil {
				return clientCtx.PrintBytes(res)
			}
			bz, _ := json.MarshalIndent(parsed, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQuerySentimentHistory() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sentiment-history [pool-id]",
		Short: "Query sentiment trading data history for a pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			_ = args
			return clientCtx.PrintString("Use the gRPC/REST endpoint: /syreen/dex/v1/sentiment/{pool_id}/history\n")
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQuerySentimentAlerts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sentiment-alerts [pool-id]",
		Short: "Query active sentiment alerts for a pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid pool-id: %w", err)
			}
			key := types.SentimentAlertKey(poolID)
			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil || len(res) == 0 {
				result := map[string]interface{}{
					"pool_id": poolID,
					"alerts":  []interface{}{},
				}
				bz, _ := json.MarshalIndent(result, "", "  ")
				return clientCtx.PrintBytes(bz)
			}
			var parsed interface{}
			if err := json.Unmarshal(res, &parsed); err != nil {
				return clientCtx.PrintBytes(res)
			}
			bz, _ := json.MarshalIndent(parsed, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// ---------------------------------------------------------------------------
// Sentiment Oracle Primitive CLI Commands
// ---------------------------------------------------------------------------

func CmdQueryOracle() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oracle [pool-id]",
		Short: "Query the sentiment oracle composite score for a pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid pool-id: %w", err)
			}
			key := types.OracleCompositeKey(poolID)
			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil || len(res) == 0 {
				result := map[string]interface{}{
					"pool_id":         poolID,
					"composite_score": 0,
					"trend":           "neutral",
					"signals":         []interface{}{},
					"whale_alerts":    []interface{}{},
					"message":         "No oracle data yet",
				}
				bz, _ := json.MarshalIndent(result, "", "  ")
				return clientCtx.PrintBytes(bz)
			}
			var parsed interface{}
			if err := json.Unmarshal(res, &parsed); err != nil {
				return clientCtx.PrintBytes(res)
			}
			bz, _ := json.MarshalIndent(parsed, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryWhaleAlerts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whale-alerts [pool-id]",
		Short: "Query recent whale alerts for a pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			_ = args
			return clientCtx.PrintString("Use the gRPC/REST endpoint: /syreen/dex/v1/whale_alerts\n")
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// ---------------------------------------------------------------------------
// IBC Cross-Chain Order CLI Commands
// ---------------------------------------------------------------------------

func CmdQueryIBCOrders() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ibc-orders [source-chain]",
		Short: "Query cross-chain IBC orders by source chain",
		Long:  "Query all cross-chain orders submitted via IBC from a specific source chain. Optionally filter by sender address using --sender flag.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			sourceChain := args[0]
			sender, _ := cmd.Flags().GetString("sender")

			// Scan IBC order keys from the store using prefix iteration.
			// Key format: ibc_order/<sourceChain>/<sender>/<orderID_8bytes>
			prefix := "ibc_order/" + sourceChain + "/"
			if sender != "" {
				prefix = "ibc_order/" + sourceChain + "/" + sender + "/"
			}

			// Use ABCI query to fetch keys with prefix
			abciRes, err := clientCtx.QueryABCI(
				abci.RequestQuery{
					Path:   fmt.Sprintf("store/%s/subspace", types.StoreKey),
					Data:   []byte(prefix),
					Prove:  false,
				},
			)
			if err != nil {
				// Fallback: direct message
				result := map[string]interface{}{
					"source_chain": sourceChain,
					"message":      "Use the gRPC/REST endpoint: /syreen/dex/v1/ibc_orders",
					"sender":       sender,
				}
				bz, _ := json.MarshalIndent(result, "", "  ")
				return clientCtx.PrintBytes(bz)
			}

			if len(abciRes.Value) == 0 {
				result := map[string]interface{}{
					"source_chain": sourceChain,
					"orders":       []interface{}{},
					"count":        0,
				}
				bz, _ := json.MarshalIndent(result, "", "  ")
				return clientCtx.PrintBytes(bz)
			}

			// Try to parse the response as JSON
			var parsed interface{}
			if err := json.Unmarshal(abciRes.Value, &parsed); err != nil {
				return clientCtx.PrintBytes(abciRes.Value)
			}
			bz, _ := json.MarshalIndent(parsed, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	cmd.Flags().String("sender", "", "filter by sender address on source chain")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryTimeWeightedPower() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "time-weighted-power [address]",
		Short: "Query time-weighted staking power for an address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			key := []byte("twp/" + args[0])
			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil || len(res) == 0 {
				result := map[string]interface{}{
					"address":             args[0],
					"time_weighted_power": "0",
				}
				bz, _ := json.MarshalIndent(result, "", "  ")
				return clientCtx.PrintBytes(bz)
			}
			var parsed interface{}
			if err := json.Unmarshal(res, &parsed); err != nil {
				return clientCtx.PrintBytes(res)
			}
			bz, _ := json.MarshalIndent(parsed, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryTraderRecord() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trader-record [address]",
		Short: "Query on-chain track record for a trader",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			key := []byte("trader_record/" + args[0])
			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil || len(res) == 0 {
				result := map[string]interface{}{
					"address":      args[0],
					"total_trades": 0,
					"total_volume": "0",
				}
				bz, _ := json.MarshalIndent(result, "", "  ")
				return clientCtx.PrintBytes(bz)
			}
			var parsed interface{}
			if err := json.Unmarshal(res, &parsed); err != nil {
				return clientCtx.PrintBytes(res)
			}
			bz, _ := json.MarshalIndent(parsed, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryTraderLeaderboard() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "leaderboard",
		Short: "Query the trader leaderboard (sorted by PnL, volume, win_rate, trades, or followers)",
		Long:  "Returns top traders sorted by the specified criterion. Default: PnL, limit: 10.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			sortBy, _ := cmd.Flags().GetString("sort-by")
			limit, _ := cmd.Flags().GetInt("limit")
			result := map[string]interface{}{
				"sort_by": sortBy,
				"limit":   limit,
				"traders": []interface{}{},
				"message": "Use the gRPC/REST endpoint: /syreen/dex/v1/trader_leaderboard?sort_by=pnl&limit=10",
			}
			bz, _ := json.MarshalIndent(result, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}
	cmd.Flags().String("sort-by", "pnl", "Sort criterion: pnl, volume, win_rate, trades, followers")
	cmd.Flags().Int("limit", 10, "Number of traders to return (max 100)")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryPoolFee() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool-fee [pool-id]",
		Short: "Query current fee information for a pool (including dynamic fee status and volatility)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			poolID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid pool-id: %w", err)
			}

			key := types.PoolKey(poolID)
			res, _, err := clientCtx.QueryStore(key, types.StoreKey)
			if err != nil {
				return fmt.Errorf("failed to query pool: %w", err)
			}

			if len(res) == 0 {
				return fmt.Errorf("pool %d not found", poolID)
			}

			var pool types.Pool
			if err := json.Unmarshal(res, &pool); err != nil {
				return fmt.Errorf("failed to unmarshal pool: %w", err)
			}

			result := map[string]interface{}{
				"pool_id":                poolID,
				"swap_fee":               pool.SwapFee.String(),
				"volatility_fee_enabled": pool.VolatilityFeeEnabled,
			}

			if pool.VolatilityFeeEnabled {
				result["base_fee"] = pool.BaseFee.String()
				result["max_fee"] = pool.MaxFee.String()
				result["volatility_multiplier"] = pool.VolatilityMultiplier.String()
				result["note"] = "Effective fee is calculated dynamically based on recent price volatility"
			} else {
				result["effective_fee"] = pool.SwapFee.String()
				result["note"] = "Dynamic fees disabled; using static swap fee"
			}

			bz, _ := json.MarshalIndent(result, "", "  ")
			return clientCtx.PrintBytes(bz)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
