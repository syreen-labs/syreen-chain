package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/dex/types"
)

// ---------------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------------

const (
	// KV store prefixes
	copyFollowPrefix       = "copytrade/follow/"       // follower/trader -> CopySettings
	copyFollowerIndexPrefix = "copytrade/follower_idx/" // trader/follower -> exists
	copyTraderStatsPrefix  = "copytrade/trader_stats/"  // trader -> TraderStats
	copyTradeHistoryPrefix = "copytrade/trade_history/" // trader -> []TradeRecord (last 100)
	copyPendingTradesPrefix = "copytrade/pending/"      // block-scoped pending trades
	copyTradeLogPrefix     = "copytrade/log/"           // follower -> []CopyTradeLog (last 100)

	// Limits
	MaxFollowersPerTrader = 1000
	MaxFollowsPerUser     = 20
	MaxTradeHistory       = 100
	MaxCopyTradeLog       = 100

	// Protection defaults
	DefaultMaxSlippageBps    = int64(500)  // 5% baseline slippage tolerance
	DefaultCooldownBlocks    = int64(2)    // minimum blocks between copy trades for same follower+trader
	// M10: -20% was weaponizable. A malicious trader could intentionally
	// take a single ~20% loss to mass-unfollow every follower in one shot
	// (a free griefing primitive). Raise the bar to -50%.
	AutoUnfollowPnLThreshold = -5000      // -50% in basis points

	// MaxFollowerSlippageBps is the maximum allowed price deviation between
	// the pre-trade price and the follower's execution price. If the leader's
	// trade moved the price beyond this threshold, the copy trade is skipped.
	// This mitigates front-running where the leader gets a better price.
	MaxFollowerSlippageBps = int64(200) // 2%

	// MinCopyTradeAmount is the minimum amount (in base units) for a copy trade.
	// Copy trades below this threshold are skipped to prevent gas-griefing
	// followers with micro-trades.
	MinCopyTradeAmount = int64(1000000) // 1 token (assuming 6 decimal places)

	// M10: copy-trade slippage tolerance scales with the relative size of
	// the copy trade vs the pool — larger relative trades face larger price
	// impact, so they need a wider tolerance to land at all. Capped to keep
	// it from becoming a free pass.
	CopyMaxAdditionalSlippageBps = int64(2000) // up to +20% extra

	// M10: oracle/TWAP-style sanity check. If the copy trade's expected
	// minOut would imply a price more than this far from the current pool
	// spot price, skip the trade entirely (do NOT unfollow). This blocks
	// sandwich/manipulation copy-baits.
	CopyTwapSanityBps = int64(500) // 5%
)

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// CopySettings stores the follow relationship and copy parameters.
type CopySettings struct {
	Follower    string   `json:"follower"`
	Trader      string   `json:"trader"`
	MaxPerTrade math.Int `json:"max_per_trade"`  // max amount per copy trade
	TotalBudget math.Int `json:"total_budget"`   // total budget for copying
	Spent       math.Int `json:"spent"`          // how much of budget has been spent
	CopyRatio   int64    `json:"copy_ratio_bps"` // basis points: 5000 = 50%
	Active      bool     `json:"active"`
	CreatedAt   int64    `json:"created_at"`     // block height
	LastCopyAt  int64    `json:"last_copy_at"`   // block height of last copy execution
}

// TraderStats tracks per-trader performance metrics.
type TraderStats struct {
	Trader       string         `json:"trader"`
	TotalTrades  int64          `json:"total_trades"`
	WinCount     int64          `json:"win_count"`
	LossCount    int64          `json:"loss_count"`
	TotalVolume  math.Int       `json:"total_volume"`   // sum of all trade input amounts
	TotalPnL     math.Int       `json:"total_pnl"`      // cumulative PnL in base units
	AvgTradeSize math.Int       `json:"avg_trade_size"`
	WinRate      math.LegacyDec `json:"win_rate"`       // 0.0 to 1.0
	FollowerCount int64         `json:"follower_count"`
	UpdatedAt    int64          `json:"updated_at"`
}

// TradeRecord stores a single trade for history/PnL calculation.
type TradeRecord struct {
	PoolID     uint64   `json:"pool_id"`
	TokenIn    sdk.Coin `json:"token_in"`
	TokenOut   sdk.Coin `json:"token_out"`
	BlockHeight int64   `json:"block_height"`
	PnL        math.Int `json:"pnl"` // positive = profit, negative = loss (estimated)
}

// PendingCopyTrade is queued for execution in EndBlock.
type PendingCopyTrade struct {
	Trader   string   `json:"trader"`
	PoolID   uint64   `json:"pool_id"`
	TokenIn  sdk.Coin `json:"token_in"`
	TokenOut sdk.Coin `json:"token_out"`
	Height   int64    `json:"height"`
}

// CopyTradeLog records copy trade execution results for a follower.
type CopyTradeLog struct {
	Follower    string   `json:"follower"`
	Trader      string   `json:"trader"`
	PoolID      uint64   `json:"pool_id"`
	TokenIn     sdk.Coin `json:"token_in"`
	TokenOut    sdk.Coin `json:"token_out"`
	Success     bool     `json:"success"`
	FailReason  string   `json:"fail_reason,omitempty"`
	BlockHeight int64    `json:"block_height"`
}

// ---------------------------------------------------------------------------
// Follow / Unfollow
// ---------------------------------------------------------------------------

// FollowTrader creates a copy-trade follow relationship.
func (k Keeper) FollowTrader(ctx context.Context, follower, trader string, maxPerTrade, totalBudget math.Int, copyRatioBps int64) error {
	if follower == trader {
		return fmt.Errorf("cannot follow yourself")
	}

	// Validate addresses
	if _, err := sdk.AccAddressFromBech32(follower); err != nil {
		return types.ErrInvalidSender
	}
	if _, err := sdk.AccAddressFromBech32(trader); err != nil {
		return fmt.Errorf("invalid trader address")
	}

	// Validate parameters
	if !maxPerTrade.IsPositive() {
		return fmt.Errorf("max_per_trade must be positive")
	}
	if !totalBudget.IsPositive() {
		return fmt.Errorf("total_budget must be positive")
	}
	if copyRatioBps <= 0 || copyRatioBps > 10000 {
		return fmt.Errorf("copy_ratio must be between 1 and 10000 basis points")
	}

	// Check if already following
	if _, exists := k.GetCopySettings(ctx, follower, trader); exists {
		return fmt.Errorf("already following this trader")
	}

	// Check max follows per user
	following := k.GetFollowing(ctx, follower)
	if len(following) >= MaxFollowsPerUser {
		return fmt.Errorf("maximum follows reached (%d)", MaxFollowsPerUser)
	}

	// Check max followers per trader
	followers := k.GetFollowers(ctx, trader)
	if len(followers) >= MaxFollowersPerTrader {
		return fmt.Errorf("trader has maximum followers (%d)", MaxFollowersPerTrader)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	settings := CopySettings{
		Follower:    follower,
		Trader:      trader,
		MaxPerTrade: maxPerTrade,
		TotalBudget: totalBudget,
		Spent:       math.ZeroInt(),
		CopyRatio:   copyRatioBps,
		Active:      true,
		CreatedAt:   sdkCtx.BlockHeight(),
		LastCopyAt:  0,
	}

	k.SetCopySettings(ctx, settings)
	k.SetFollowerIndex(ctx, trader, follower)

	// Update trader follower count
	stats := k.GetOrCreateTraderStats(ctx, trader)
	stats.FollowerCount++
	k.SetTraderStats(ctx, stats)

	k.Logger(ctx).Info("copy trade: follow",
		"follower", follower,
		"trader", trader,
		"copy_ratio_bps", copyRatioBps,
		"max_per_trade", maxPerTrade,
		"total_budget", totalBudget,
	)

	return nil
}

// UnfollowTrader removes a copy-trade follow relationship.
func (k Keeper) UnfollowTrader(ctx context.Context, follower, trader string) error {
	if _, err := sdk.AccAddressFromBech32(follower); err != nil {
		return types.ErrInvalidSender
	}

	_, exists := k.GetCopySettings(ctx, follower, trader)
	if !exists {
		return fmt.Errorf("not following this trader")
	}

	k.DeleteCopySettings(ctx, follower, trader)
	k.DeleteFollowerIndex(ctx, trader, follower)

	// Update trader follower count
	stats := k.GetOrCreateTraderStats(ctx, trader)
	if stats.FollowerCount > 0 {
		stats.FollowerCount--
	}
	k.SetTraderStats(ctx, stats)

	k.Logger(ctx).Info("copy trade: unfollow",
		"follower", follower,
		"trader", trader,
	)

	return nil
}

// UpdateCopySettings updates copy parameters for an existing follow.
func (k Keeper) UpdateCopySettings(ctx context.Context, follower, trader string, maxPerTrade, totalBudget math.Int, copyRatioBps int64) error {
	if _, err := sdk.AccAddressFromBech32(follower); err != nil {
		return types.ErrInvalidSender
	}

	settings, exists := k.GetCopySettings(ctx, follower, trader)
	if !exists {
		return fmt.Errorf("not following this trader")
	}

	if maxPerTrade.IsPositive() {
		settings.MaxPerTrade = maxPerTrade
	}
	if totalBudget.IsPositive() {
		settings.TotalBudget = totalBudget
	}
	if copyRatioBps > 0 && copyRatioBps <= 10000 {
		settings.CopyRatio = copyRatioBps
	}

	k.SetCopySettings(ctx, settings)

	k.Logger(ctx).Info("copy trade: update settings",
		"follower", follower,
		"trader", trader,
	)

	return nil
}

// ---------------------------------------------------------------------------
// Trade Recording (called from Swap)
// ---------------------------------------------------------------------------

// copyTradeCtxKey is used to mark context during copy trade execution to prevent recursive amplification.
type copyTradeCtxKey struct{}

// RecordTradeForCopyTrading records a trade so followers can copy it in EndBlock.
//
// L1: `reserveIncrease` is the actual amount that was added to the input-side
// reserve after referral fee deductions (i.e., tokenIn.Amount minus any
// referral payouts). It is used to reconstruct the pre-swap reserves
// accurately so the PnL estimate is not skewed by the referral split. Pass
// tokenIn.Amount when no referral fees applied.
func (k Keeper) RecordTradeForCopyTrading(ctx context.Context, trader string, poolID uint64, tokenIn, tokenOut sdk.Coin, reserveIncrease math.Int) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Skip recording if this swap was triggered by copy trading to prevent recursive amplification
	if sdkCtx.Value(copyTradeCtxKey{}) != nil {
		return
	}
	height := sdkCtx.BlockHeight()

	// Add pending trade for EndBlock processing
	pending := PendingCopyTrade{
		Trader:   trader,
		PoolID:   poolID,
		TokenIn:  tokenIn,
		TokenOut: tokenOut,
		Height:   height,
	}
	k.addPendingCopyTrade(ctx, pending)

	// Update trader stats
	stats := k.GetOrCreateTraderStats(ctx, trader)
	stats.TotalTrades++
	stats.TotalVolume = stats.TotalVolume.Add(tokenIn.Amount)

	// PnL estimation: normalize output to input denomination using pre-swap spot price.
	// Since RecordTradeForCopyTrading is called after pool reserves are updated,
	// we reconstruct the pre-swap reserves to get an accurate conversion rate.
	//
	// L1: use reserveIncrease (the actual delta that hit the reserve after
	// referral fees were peeled off) instead of tokenIn.Amount, otherwise
	// the reconstructed pre-swap reserve is off by the referral payout and
	// PnL is misattributed.
	if reserveIncrease.IsNil() {
		reserveIncrease = tokenIn.Amount
	}
	var pnl math.Int
	pool, poolFound := k.GetPool(ctx, poolID)
	if poolFound {
		// Reconstruct pre-swap reserves
		var preReserveIn, preReserveOut math.Int
		if tokenIn.Denom == pool.DenomA && tokenOut.Denom == pool.DenomB {
			preReserveIn = pool.ReserveA.Sub(reserveIncrease)
			preReserveOut = pool.ReserveB.Add(tokenOut.Amount)
		} else if tokenIn.Denom == pool.DenomB && tokenOut.Denom == pool.DenomA {
			preReserveIn = pool.ReserveB.Sub(reserveIncrease)
			preReserveOut = pool.ReserveA.Add(tokenOut.Amount)
		} else {
			preReserveIn = math.ZeroInt()
			preReserveOut = math.ZeroInt()
		}

		if preReserveOut.IsPositive() && preReserveIn.IsPositive() {
			// Convert output to input denomination using pre-swap spot price
			// outputValue = tokenOut.Amount * preReserveIn / preReserveOut
			outputValue := tokenOut.Amount.Mul(preReserveIn).Quo(preReserveOut)
			pnl = outputValue.Sub(tokenIn.Amount)
		} else {
			pnl = tokenOut.Amount.Sub(tokenIn.Amount)
		}
	} else {
		// Fallback: raw amount comparison (inaccurate across different denoms)
		pnl = tokenOut.Amount.Sub(tokenIn.Amount)
	}
	stats.TotalPnL = stats.TotalPnL.Add(pnl)
	if pnl.IsPositive() {
		stats.WinCount++
	} else {
		stats.LossCount++
	}

	// Update averages
	if stats.TotalTrades > 0 {
		stats.AvgTradeSize = stats.TotalVolume.Quo(math.NewInt(stats.TotalTrades))
		total := stats.WinCount + stats.LossCount
		if total > 0 {
			stats.WinRate = math.LegacyNewDec(stats.WinCount).Quo(math.LegacyNewDec(total))
		}
	}
	stats.UpdatedAt = height

	k.SetTraderStats(ctx, stats)

	// Store trade record (keep last 100)
	record := TradeRecord{
		PoolID:      poolID,
		TokenIn:     tokenIn,
		TokenOut:    tokenOut,
		BlockHeight: height,
		PnL:         pnl,
	}
	k.appendTradeRecord(ctx, trader, record)

	// Update on-chain track record for copy trading leaderboard.
	// Calculate PnL percentage: pnl / tokenIn.Amount * 100
	pnlPercent := math.LegacyZeroDec()
	if tokenIn.Amount.IsPositive() {
		pnlPercent = math.LegacyNewDecFromInt(pnl).Quo(math.LegacyNewDecFromInt(tokenIn.Amount)).MulInt64(100)
	}
	k.UpdateTraderRecord(ctx, trader, poolID, tokenIn.Amount, tokenOut.Amount, pnlPercent)
}

// ---------------------------------------------------------------------------
// EndBlock: Process Pending Copy Trades
// ---------------------------------------------------------------------------

// ProcessCopyTrades executes pending copy trades. Called from EndBlock.
func (k Keeper) ProcessCopyTrades(ctx context.Context) {
	pending := k.getPendingCopyTrades(ctx)
	if len(pending) == 0 {
		return
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	height := sdkCtx.BlockHeight()

	// H4-FIX: Cap total copy trade executions per block to prevent block timeout
	const maxCopyTradesPerBlock = 100
	totalExecutions := 0

	for _, trade := range pending {
		followers := k.GetFollowers(ctx, trade.Trader)
		for _, followerAddr := range followers {
			// H4-FIX: Stop processing if we've hit the per-block cap
			if totalExecutions >= maxCopyTradesPerBlock {
				break
			}
			settings, exists := k.GetCopySettings(ctx, followerAddr, trade.Trader)
			if !exists || !settings.Active {
				continue
			}

			// Cooldown check
			if settings.LastCopyAt > 0 && (height-settings.LastCopyAt) < DefaultCooldownBlocks {
				continue
			}

			// Budget check
			remaining := settings.TotalBudget.Sub(settings.Spent)
			if !remaining.IsPositive() {
				// Budget exhausted, deactivate
				settings.Active = false
				k.SetCopySettings(ctx, settings)
				continue
			}

			// Calculate copy amount based on ratio
			copyAmount := trade.TokenIn.Amount.MulRaw(settings.CopyRatio).QuoRaw(10000)
			if copyAmount.IsZero() {
				continue
			}

			// Minimum copy trade size check — prevent gas-griefing with micro-trades
			if copyAmount.LT(math.NewInt(MinCopyTradeAmount)) {
				k.appendCopyTradeLog(ctx, followerAddr, CopyTradeLog{
					Follower:    followerAddr,
					Trader:      trade.Trader,
					PoolID:      trade.PoolID,
					TokenIn:     sdk.NewCoin(trade.TokenIn.Denom, copyAmount),
					Success:     false,
					FailReason:  fmt.Sprintf("copy amount %s below minimum %d", copyAmount, MinCopyTradeAmount),
					BlockHeight: height,
				})
				continue
			}

			// Cap at max per trade
			if copyAmount.GT(settings.MaxPerTrade) {
				copyAmount = settings.MaxPerTrade
			}

			// Cap at remaining budget
			if copyAmount.GT(remaining) {
				copyAmount = remaining
			}

			// Front-running protection: check if the current pool spot price
			// (which reflects the leader's trade impact) has deviated too far.
			// We use the pool's current spot price (reserveOut/reserveIn) as the
			// pre-trade price for the FOLLOWER's slippage check, then compare
			// the follower's expected execution price against it.
			pool, poolFound := k.GetPool(ctx, trade.PoolID)
			if poolFound {
				// Get the pool's current spot price (this is what the follower faces)
				var currentSpotPrice math.LegacyDec
				var spotErr error
				if trade.TokenIn.Denom == pool.DenomA {
					currentSpotPrice, spotErr = k.GetSpotPrice(ctx, trade.PoolID, pool.DenomA, pool.DenomB)
				} else {
					currentSpotPrice, spotErr = k.GetSpotPrice(ctx, trade.PoolID, pool.DenomB, pool.DenomA)
				}
				// Also compute the leader's effective price to detect front-running
				if spotErr == nil && currentSpotPrice.IsPositive() && trade.TokenIn.Amount.IsPositive() && trade.TokenOut.Amount.IsPositive() {
					leaderPrice := math.LegacyNewDecFromInt(trade.TokenOut.Amount).Quo(math.LegacyNewDecFromInt(trade.TokenIn.Amount))
					// If spot price moved against the follower by more than MaxFollowerSlippageBps
					// compared to the leader's achieved price, skip
					priceDiff := leaderPrice.Sub(currentSpotPrice).Quo(leaderPrice)
					maxSlippage := math.LegacyNewDec(MaxFollowerSlippageBps).QuoInt64(10000)
					if priceDiff.GT(maxSlippage) {
						k.appendCopyTradeLog(ctx, followerAddr, CopyTradeLog{
							Follower:    followerAddr,
							Trader:      trade.Trader,
							PoolID:      trade.PoolID,
							TokenIn:     sdk.NewCoin(trade.TokenIn.Denom, copyAmount),
							Success:     false,
							FailReason:  fmt.Sprintf("leader front-running: spot price %s vs leader price %s, diff %s > max %s", currentSpotPrice, leaderPrice, priceDiff, maxSlippage),
							BlockHeight: height,
						})
						continue
					}
				}
			}

			// Calculate min output with slippage protection.
			// Baseline DefaultMaxSlippageBps + a pool-weighted bonus (M10):
			// the larger the copy trade relative to the input-side reserve,
			// the more price impact we expect, so widen the tolerance.
			expectedOut := math.ZeroInt()
			if trade.TokenIn.Amount.IsPositive() {
				expectedOut = trade.TokenOut.Amount.Mul(copyAmount).Quo(trade.TokenIn.Amount)
			}

			extraBps := int64(0)
			if !poolFound {
				pool, poolFound = k.GetPool(ctx, trade.PoolID)
			}
			if poolFound {
				var reserveIn math.Int
				if trade.TokenIn.Denom == pool.DenomA {
					reserveIn = pool.ReserveA
				} else if trade.TokenIn.Denom == pool.DenomB {
					reserveIn = pool.ReserveB
				} else {
					reserveIn = math.ZeroInt()
				}
				if reserveIn.IsPositive() {
					// extra = copyAmount / reserveIn, expressed in bps and capped.
					ratioBps := copyAmount.MulRaw(10000).Quo(reserveIn).Int64()
					if ratioBps < 0 {
						ratioBps = 0
					}
					if ratioBps > CopyMaxAdditionalSlippageBps {
						ratioBps = CopyMaxAdditionalSlippageBps
					}
					extraBps = ratioBps
				}
			}

			effectiveSlippage := DefaultMaxSlippageBps + extraBps
			if effectiveSlippage >= 10000 {
				effectiveSlippage = 9999
			}
			minOut := expectedOut.MulRaw(10000 - effectiveSlippage).QuoRaw(10000)

			// M10: TWAP / spot-price sanity check. If the implied minOut
			// price differs from the current pool spot price by more than
			// CopyTwapSanityBps, skip this copy trade — but DO NOT unfollow
			// (keep the relationship intact, just refuse this single trade).
			if poolFound && expectedOut.IsPositive() && copyAmount.IsPositive() {
				// implied price = expectedOut / copyAmount (denomOut per denomIn)
				impliedPrice := math.LegacyNewDecFromInt(expectedOut).Quo(math.LegacyNewDecFromInt(copyAmount))
				var spotPrice math.LegacyDec
				var spotErr error
				if trade.TokenIn.Denom == pool.DenomA {
					spotPrice, spotErr = k.GetSpotPrice(ctx, trade.PoolID, pool.DenomA, pool.DenomB)
				} else {
					spotPrice, spotErr = k.GetSpotPrice(ctx, trade.PoolID, pool.DenomB, pool.DenomA)
				}
				if spotErr == nil && spotPrice.IsPositive() {
					// diff = |implied - spot| / spot
					diff := impliedPrice.Sub(spotPrice).Abs().Quo(spotPrice)
					maxDiff := math.LegacyNewDec(CopyTwapSanityBps).QuoInt64(10000)
					if diff.GT(maxDiff) {
						k.appendCopyTradeLog(ctx, followerAddr, CopyTradeLog{
							Follower:    followerAddr,
							Trader:      trade.Trader,
							PoolID:      trade.PoolID,
							TokenIn:     sdk.NewCoin(trade.TokenIn.Denom, copyAmount),
							Success:     false,
							FailReason:  fmt.Sprintf("twap sanity check failed: diff=%s > %s", diff, maxDiff),
							BlockHeight: height,
						})
						continue
					}
				}
			}

			// Execute the copy swap with copy-trade context flag to prevent recursive recording
			copyTokenIn := sdk.NewCoin(trade.TokenIn.Denom, copyAmount)
			copyCtx := sdkCtx.WithValue(copyTradeCtxKey{}, true)
			tokenOut, err := k.Swap(copyCtx, followerAddr, trade.PoolID, copyTokenIn, minOut)

			logEntry := CopyTradeLog{
				Follower:    followerAddr,
				Trader:      trade.Trader,
				PoolID:      trade.PoolID,
				TokenIn:     copyTokenIn,
				BlockHeight: height,
			}

			// H4-FIX: Count execution attempt
			totalExecutions++

			if err != nil {
				logEntry.Success = false
				logEntry.FailReason = err.Error()
				logEntry.TokenOut = sdk.Coin{}
				// H3-FIX: No penalty on failure — prevents budget drain attacks.
				// Failed trades cost nothing since no actual swap occurred.
			} else {
				logEntry.Success = true
				logEntry.TokenOut = tokenOut

				// Update spent
				settings.LastCopyAt = height
				settings.Spent = settings.Spent.Add(copyAmount)
				k.SetCopySettings(ctx, settings)
			}

			k.appendCopyTradeLog(ctx, followerAddr, logEntry)
		}
		// H4-FIX: Break outer loop too if cap hit
		if totalExecutions >= maxCopyTradesPerBlock {
			break
		}
	}

	// Clear pending trades
	k.clearPendingCopyTrades(ctx)

	// Auto-unfollow check: if trader PnL dropped below threshold
	k.checkAutoUnfollow(ctx)
}

// checkAutoUnfollow unfollows traders whose PnL has dropped below -20%.
func (k Keeper) checkAutoUnfollow(ctx context.Context) {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(copyTraderStatsPrefix)
	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return
	}

	// Collect traders to unfollow — do NOT mutate KVStore during iteration
	var toUnfollow []TraderStats
	for ; iter.Valid(); iter.Next() {
		var stats TraderStats
		if err := json.Unmarshal(iter.Value(), &stats); err != nil {
			continue
		}
		if !stats.TotalVolume.IsPositive() {
			continue
		}
		pnlBps := stats.TotalPnL.MulRaw(10000).Quo(stats.TotalVolume).Int64()
		if pnlBps < int64(AutoUnfollowPnLThreshold) {
			toUnfollow = append(toUnfollow, stats)
		}
	}
	iter.Close()

	// Cache all followers before any mutation to avoid iterator invalidation
	type traderFollowers struct {
		stats     TraderStats
		followers []string
	}
	cached := make([]traderFollowers, 0, len(toUnfollow))
	for _, stats := range toUnfollow {
		cached = append(cached, traderFollowers{
			stats:     stats,
			followers: k.GetFollowers(ctx, stats.Trader),
		})
	}

	// Apply modifications after all reads are done
	for _, tf := range cached {
		for _, f := range tf.followers {
			k.DeleteCopySettings(ctx, f, tf.stats.Trader)
			k.DeleteFollowerIndex(ctx, tf.stats.Trader, f)
		}
		tf.stats.FollowerCount = 0
		k.SetTraderStats(ctx, tf.stats)
	}
}

// ---------------------------------------------------------------------------
// Leaderboard
// ---------------------------------------------------------------------------

// GetLeaderboard returns top traders sorted by PnL.
func (k Keeper) GetLeaderboard(ctx context.Context, limit int) []TraderStats {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(copyTraderStatsPrefix)
	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var allStats []TraderStats
	for ; iter.Valid(); iter.Next() {
		var stats TraderStats
		if err := json.Unmarshal(iter.Value(), &stats); err != nil {
			continue
		}
		if stats.TotalTrades > 0 {
			allStats = append(allStats, stats)
		}
	}

	// Sort by PnL descending, with address-based tiebreaker for determinism
	sort.SliceStable(allStats, func(i, j int) bool {
		if allStats[i].TotalPnL.Equal(allStats[j].TotalPnL) {
			return allStats[i].Trader < allStats[j].Trader
		}
		return allStats[i].TotalPnL.GT(allStats[j].TotalPnL)
	})

	if limit > 0 && len(allStats) > limit {
		allStats = allStats[:limit]
	}

	return allStats
}

// ---------------------------------------------------------------------------
// KV Store: Copy Settings
// ---------------------------------------------------------------------------

func copySettingsKey(follower, trader string) []byte {
	return []byte(copyFollowPrefix + follower + "/" + trader)
}

func (k Keeper) GetCopySettings(ctx context.Context, follower, trader string) (CopySettings, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(copySettingsKey(follower, trader))
	if err != nil || bz == nil {
		return CopySettings{}, false
	}
	var s CopySettings
	if err := json.Unmarshal(bz, &s); err != nil {
		return CopySettings{}, false
	}
	return s, true
}

func (k Keeper) SetCopySettings(ctx context.Context, s CopySettings) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(s)
	kvStore.Set(copySettingsKey(s.Follower, s.Trader), bz)
}

func (k Keeper) DeleteCopySettings(ctx context.Context, follower, trader string) {
	kvStore := k.storeService.OpenKVStore(ctx)
	kvStore.Delete(copySettingsKey(follower, trader))
}

// ---------------------------------------------------------------------------
// KV Store: Follower Index (trader -> list of followers)
// ---------------------------------------------------------------------------

func followerIndexKey(trader, follower string) []byte {
	return []byte(copyFollowerIndexPrefix + trader + "/" + follower)
}

func followerIndexPrefix(trader string) []byte {
	return []byte(copyFollowerIndexPrefix + trader + "/")
}

func (k Keeper) SetFollowerIndex(ctx context.Context, trader, follower string) {
	kvStore := k.storeService.OpenKVStore(ctx)
	kvStore.Set(followerIndexKey(trader, follower), []byte("1"))
}

func (k Keeper) DeleteFollowerIndex(ctx context.Context, trader, follower string) {
	kvStore := k.storeService.OpenKVStore(ctx)
	kvStore.Delete(followerIndexKey(trader, follower))
}

// GetFollowers returns all follower addresses for a trader.
func (k Keeper) GetFollowers(ctx context.Context, trader string) []string {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := followerIndexPrefix(trader)
	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var followers []string
	for ; iter.Valid(); iter.Next() {
		key := string(iter.Key())
		// Key format: copytrade/follower_idx/trader/follower
		// Extract follower from after the last "/"
		prefixStr := copyFollowerIndexPrefix + trader + "/"
		if len(key) > len(prefixStr) {
			follower := key[len(prefixStr):]
			followers = append(followers, follower)
		}
	}
	sort.Strings(followers)
	return followers
}

// GetFollowing returns all traders a user is following.
func (k Keeper) GetFollowing(ctx context.Context, follower string) []CopySettings {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(copyFollowPrefix + follower + "/")
	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var result []CopySettings
	for ; iter.Valid(); iter.Next() {
		var s CopySettings
		if err := json.Unmarshal(iter.Value(), &s); err != nil {
			continue
		}
		result = append(result, s)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Follower != result[j].Follower {
			return result[i].Follower < result[j].Follower
		}
		return result[i].Trader < result[j].Trader
	})
	return result
}

// ---------------------------------------------------------------------------
// KV Store: Trader Stats
// ---------------------------------------------------------------------------

func traderStatsKey(trader string) []byte {
	return []byte(copyTraderStatsPrefix + trader)
}

func (k Keeper) GetTraderStats(ctx context.Context, trader string) (TraderStats, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(traderStatsKey(trader))
	if err != nil || bz == nil {
		return TraderStats{}, false
	}
	var s TraderStats
	if err := json.Unmarshal(bz, &s); err != nil {
		return TraderStats{}, false
	}
	return s, true
}

func (k Keeper) GetOrCreateTraderStats(ctx context.Context, trader string) TraderStats {
	stats, exists := k.GetTraderStats(ctx, trader)
	if !exists {
		stats = TraderStats{
			Trader:      trader,
			TotalVolume: math.ZeroInt(),
			TotalPnL:    math.ZeroInt(),
			AvgTradeSize: math.ZeroInt(),
			WinRate:     math.LegacyZeroDec(),
		}
	}
	return stats
}

func (k Keeper) SetTraderStats(ctx context.Context, s TraderStats) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(s)
	kvStore.Set(traderStatsKey(s.Trader), bz)
}

// ---------------------------------------------------------------------------
// KV Store: Trade History (per trader, last 100)
// ---------------------------------------------------------------------------

func tradeHistoryKey(trader string) []byte {
	return []byte(copyTradeHistoryPrefix + trader)
}

func (k Keeper) getTradeHistory(ctx context.Context, trader string) []TradeRecord {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(tradeHistoryKey(trader))
	if err != nil || bz == nil {
		return nil
	}
	var records []TradeRecord
	if err := json.Unmarshal(bz, &records); err != nil {
		return nil
	}
	return records
}

func (k Keeper) appendTradeRecord(ctx context.Context, trader string, record TradeRecord) {
	records := k.getTradeHistory(ctx, trader)
	records = append(records, record)
	if len(records) > MaxTradeHistory {
		records = records[len(records)-MaxTradeHistory:]
	}
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(records)
	kvStore.Set(tradeHistoryKey(trader), bz)
}

// ---------------------------------------------------------------------------
// KV Store: Pending Copy Trades (per block)
// ---------------------------------------------------------------------------

func (k Keeper) addPendingCopyTrade(ctx context.Context, trade PendingCopyTrade) {
	kvStore := k.storeService.OpenKVStore(ctx)
	key := []byte(copyPendingTradesPrefix)

	var trades []PendingCopyTrade
	bz, err := kvStore.Get(key)
	if err == nil && bz != nil {
		json.Unmarshal(bz, &trades)
	}
	trades = append(trades, trade)
	bz, _ = json.Marshal(trades)
	kvStore.Set(key, bz)
}

func (k Keeper) getPendingCopyTrades(ctx context.Context) []PendingCopyTrade {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(copyPendingTradesPrefix))
	if err != nil || bz == nil {
		return nil
	}
	var trades []PendingCopyTrade
	if err := json.Unmarshal(bz, &trades); err != nil {
		return nil
	}
	return trades
}

func (k Keeper) clearPendingCopyTrades(ctx context.Context) {
	kvStore := k.storeService.OpenKVStore(ctx)
	kvStore.Delete([]byte(copyPendingTradesPrefix))
}

// ---------------------------------------------------------------------------
// KV Store: Copy Trade Log (per follower, last 100)
// ---------------------------------------------------------------------------

func copyTradeLogKey(follower string) []byte {
	return []byte(copyTradeLogPrefix + follower)
}

func (k Keeper) GetCopyTradeLog(ctx context.Context, follower string) []CopyTradeLog {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(copyTradeLogKey(follower))
	if err != nil || bz == nil {
		return nil
	}
	var logs []CopyTradeLog
	if err := json.Unmarshal(bz, &logs); err != nil {
		return nil
	}
	return logs
}

func (k Keeper) appendCopyTradeLog(ctx context.Context, follower string, entry CopyTradeLog) {
	logs := k.GetCopyTradeLog(ctx, follower)
	logs = append(logs, entry)
	if len(logs) > MaxCopyTradeLog {
		logs = logs[len(logs)-MaxCopyTradeLog:]
	}
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(logs)
	kvStore.Set(copyTradeLogKey(follower), bz)
}
