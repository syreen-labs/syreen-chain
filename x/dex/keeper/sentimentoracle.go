package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/dex/types"
)

// ---------------------------------------------------------------------------
// Sentiment Oracle Primitive — aggregates multiple on-chain signals into a
// composite score that other modules can consume.
// ---------------------------------------------------------------------------

const (
	// OracleUpdateInterval is how often oracle signals are recalculated (in blocks).
	// Piggybacks on the existing 5-block sentiment update cycle.
	OracleUpdateInterval = 5

	// OracleVolumeWindow is the number of blocks to analyze for volume spikes.
	OracleVolumeWindow = 200 // same as SentimentWindow

	// OracleLiquidityWindow is the number of blocks for liquidity change analysis.
	OracleLiquidityWindow = 100

	// OraclePriceMomentumWindow is the number of blocks for price momentum.
	OraclePriceMomentumWindow = 50

	// WhaleAlertThresholdBps is the minimum swap size as basis points of pool reserves
	// to qualify as a whale alert (500 bps = 5%).
	WhaleAlertThresholdBps = 500

	// MaxWhaleAlerts is the maximum number of whale alerts stored per pool.
	MaxWhaleAlerts = 100
)

// ---------------------------------------------------------------------------
// Data Structures
// ---------------------------------------------------------------------------

// OracleSignal represents an individual oracle signal score.
type OracleSignal struct {
	Name      string `json:"name"`
	Score     int64  `json:"score"`      // -100 to +100
	Weight    int64  `json:"weight"`     // weight in composite
	UpdatedAt int64  `json:"updated_at"` // block height
}

// OracleComposite holds the full composite oracle state for a pool.
type OracleComposite struct {
	PoolID         uint64         `json:"pool_id"`
	CompositeScore int64          `json:"composite_score"` // -100 to +100
	Trend          string         `json:"trend"`           // "bullish", "bearish", "neutral"
	Signals        []OracleSignal `json:"signals"`
	WhaleAlerts    []WhaleAlert   `json:"whale_alerts"`
	UpdatedAt      int64          `json:"updated_at"`
}

// WhaleAlert represents a large trade that exceeded the whale threshold.
type WhaleAlert struct {
	Trader      string   `json:"trader"`
	PoolID      uint64   `json:"pool_id"`
	Amount      math.Int `json:"amount"`
	Side        string   `json:"side"`         // "buy" or "sell"
	BlockHeight int64    `json:"block_height"`
}

// ---------------------------------------------------------------------------
// Signal Weights
// ---------------------------------------------------------------------------

var oracleSignalWeights = map[string]int64{
	"volume_spike":       25,
	"whale_movement":     20,
	"liquidity_change":   20,
	"orderbook_imbalance": 15,
	"price_momentum":     20,
}

// ---------------------------------------------------------------------------
// KV Store — OracleComposite
// ---------------------------------------------------------------------------

func (k Keeper) GetOracleComposite(ctx context.Context, poolID uint64) OracleComposite {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.OracleCompositeKey(poolID))
	if err != nil || bz == nil {
		return OracleComposite{
			PoolID:         poolID,
			CompositeScore: 0,
			Trend:          "neutral",
			Signals:        []OracleSignal{},
			WhaleAlerts:    []WhaleAlert{},
		}
	}
	var composite OracleComposite
	if err := json.Unmarshal(bz, &composite); err != nil {
		return OracleComposite{
			PoolID:         poolID,
			CompositeScore: 0,
			Trend:          "neutral",
			Signals:        []OracleSignal{},
			WhaleAlerts:    []WhaleAlert{},
		}
	}
	return composite
}

func (k Keeper) SetOracleComposite(ctx context.Context, composite OracleComposite) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(composite)
	kvStore.Set(types.OracleCompositeKey(composite.PoolID), bz)
}

// ---------------------------------------------------------------------------
// KV Store — WhaleAlerts
// ---------------------------------------------------------------------------

func (k Keeper) getWhaleAlertKey(poolID uint64, height int64, index uint64) []byte {
	return types.WhaleAlertKey(poolID, height, index)
}

func (k Keeper) RecordWhaleAlert(ctx context.Context, trader string, poolID uint64, amount math.Int, side string) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	height := sdkCtx.BlockHeight()

	// Check for reversal: if the same trader made an opposite trade in the last 5 blocks,
	// skip recording this alert to prevent whale manipulation via buy-then-sell patterns.
	recentAlerts := k.GetWhaleAlerts(ctx, poolID, 100)
	for _, recent := range recentAlerts {
		if recent.BlockHeight < height-5 {
			break // only check last 5 blocks
		}
		if recent.Trader == trader && recent.Side != side {
			k.Logger(ctx).Info("WHALE ALERT REVERSAL DETECTED — skipping",
				"pool_id", poolID,
				"trader", trader,
				"new_side", side,
				"previous_side", recent.Side,
				"block", height,
			)
			return // skip recording this alert
		}
	}

	alert := WhaleAlert{
		Trader:      trader,
		PoolID:      poolID,
		Amount:      amount,
		Side:        side,
		BlockHeight: height,
	}

	// Find the next index for this pool+height
	kvStore := k.storeService.OpenKVStore(ctx)
	var index uint64
	for {
		key := types.WhaleAlertKey(poolID, height, index)
		existing, err := kvStore.Get(key)
		if err != nil || existing == nil {
			break
		}
		index++
		if index > 100 {
			break // safety cap
		}
	}

	bz, _ := json.Marshal(alert)
	kvStore.Set(types.WhaleAlertKey(poolID, height, index), bz)

	k.Logger(ctx).Info("WHALE ALERT",
		"pool_id", poolID,
		"trader", trader,
		"amount", amount.String(),
		"side", side,
		"block", height,
	)
}

// GetWhaleAlerts returns the most recent whale alerts for a pool, up to limit.
func (k Keeper) GetWhaleAlerts(ctx context.Context, poolID uint64, limit int) []WhaleAlert {
	if limit <= 0 {
		limit = 50
	}
	if limit > MaxWhaleAlerts {
		limit = MaxWhaleAlerts
	}

	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := types.WhaleAlertPoolPrefix(poolID)
	endKey := prefixEndBytes(prefix)

	// Use reverse iterator to get most recent alerts first
	iter, err := kvStore.ReverseIterator(prefix, endKey)
	if err != nil {
		return nil
	}
	defer iter.Close()

	var alerts []WhaleAlert
	for ; iter.Valid() && len(alerts) < limit; iter.Next() {
		var alert WhaleAlert
		if err := json.Unmarshal(iter.Value(), &alert); err != nil {
			continue
		}
		alerts = append(alerts, alert)
	}
	return alerts
}

// ---------------------------------------------------------------------------
// Signal Computation
// ---------------------------------------------------------------------------

// calcVolumeSpike computes the volume spike signal (-100 to +100).
// Compares current block's volume to the 24h (200 block) average.
// Positive = volume above average (bullish activity), negative = below.
func (k Keeper) calcVolumeSpike(blockData []BlockTradingData, currentHeight int64) int64 {
	if len(blockData) == 0 {
		return 0
	}

	// Total volume across all blocks
	totalVolume := math.ZeroInt()
	for _, bd := range blockData {
		totalVolume = totalVolume.Add(bd.BuyVolume).Add(bd.SellVolume)
	}

	if totalVolume.IsZero() {
		return 0
	}

	avgPerBlock := totalVolume.QuoRaw(int64(len(blockData)))
	if avgPerBlock.IsZero() {
		return 0
	}

	// Get the last 5 blocks volume (recent window)
	var recentVolume = math.ZeroInt()
	recentCount := 0
	for i := len(blockData) - 1; i >= 0 && recentCount < 5; i-- {
		bd := blockData[i]
		recentVolume = recentVolume.Add(bd.BuyVolume).Add(bd.SellVolume)
		recentCount++
	}

	if recentCount == 0 {
		return 0
	}

	recentAvg := recentVolume.QuoRaw(int64(recentCount))

	// Ratio = (recentAvg - avgPerBlock) / avgPerBlock * 100, clamped to -100..+100
	if avgPerBlock.IsZero() {
		return 0
	}

	diff := recentAvg.Sub(avgPerBlock)
	score := diff.MulRaw(100).Quo(avgPerBlock).Int64()
	return clampOracleScore(score)
}

// calcWhaleMovement computes whale movement signal (-100 to +100).
// Positive = net whale buying, negative = net whale selling.
func (k Keeper) calcWhaleMovement(blockData []BlockTradingData) int64 {
	totalWhaleBuy := math.ZeroInt()
	totalWhaleSell := math.ZeroInt()

	for _, bd := range blockData {
		totalWhaleBuy = totalWhaleBuy.Add(bd.WhaleBuyVolume)
		totalWhaleSell = totalWhaleSell.Add(bd.WhaleSellVolume)
	}

	total := totalWhaleBuy.Add(totalWhaleSell)
	if total.IsZero() {
		return 0
	}

	// Score: (buyVol - sellVol) / total * 100
	diff := totalWhaleBuy.Sub(totalWhaleSell)
	score := diff.MulRaw(100).Quo(total).Int64()
	return clampOracleScore(score)
}

// calcLiquidityChange computes liquidity change signal (-100 to +100).
// Uses the last OracleLiquidityWindow blocks.
// Positive = net LP additions, negative = net LP removals.
func (k Keeper) calcLiquidityChange(blockData []BlockTradingData, currentHeight int64) int64 {
	// Filter to last OracleLiquidityWindow blocks
	cutoff := currentHeight - OracleLiquidityWindow
	totalAdd := math.ZeroInt()
	totalRemove := math.ZeroInt()

	for _, bd := range blockData {
		if bd.Height < cutoff {
			continue
		}
		totalAdd = totalAdd.Add(bd.LPAddVolume)
		totalRemove = totalRemove.Add(bd.LPRemoveVolume)
	}

	total := totalAdd.Add(totalRemove)
	if total.IsZero() {
		return 0
	}

	diff := totalAdd.Sub(totalRemove)
	score := diff.MulRaw(100).Quo(total).Int64()
	return clampOracleScore(score)
}

// calcOrderBookImbalance computes order book imbalance signal (-100 to +100).
// Positive = more bid depth (bullish), negative = more ask depth (bearish).
func (k Keeper) calcOrderBookImbalance(ctx context.Context, poolID uint64) int64 {
	ob := k.GetOrderBook(ctx, poolID)

	bidDepth := math.ZeroInt()
	for _, level := range ob.Bids {
		bidDepth = bidDepth.Add(level.Quantity)
	}

	askDepth := math.ZeroInt()
	for _, level := range ob.Asks {
		askDepth = askDepth.Add(level.Quantity)
	}

	total := bidDepth.Add(askDepth)
	if total.IsZero() {
		return 0
	}

	diff := bidDepth.Sub(askDepth)
	score := diff.MulRaw(100).Quo(total).Int64()
	return clampOracleScore(score)
}

// calcPriceMomentum computes price momentum signal (-100 to +100).
// Compares first-half vs second-half average prices over the momentum window.
// Positive = price increasing, negative = price decreasing.
func (k Keeper) calcPriceMomentum(ctx context.Context, poolID uint64, currentHeight int64) int64 {
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return 0
	}

	// Use buy/sell volume ratio as a proxy for price direction.
	// If there's more buying (DenomB sold for DenomA), price of DenomA goes up.
	fromHeight := currentHeight - OraclePriceMomentumWindow
	if fromHeight < 1 {
		fromHeight = 1
	}

	blockData := k.GetBlockTradingDataRange(ctx, poolID, fromHeight, currentHeight-1)
	if len(blockData) < 2 {
		return 0
	}

	_ = pool // pool available if needed for spot price comparison

	// Split into two halves and compare buy pressure
	mid := len(blockData) / 2
	firstBuy := math.ZeroInt()
	firstSell := math.ZeroInt()
	secondBuy := math.ZeroInt()
	secondSell := math.ZeroInt()

	for i, bd := range blockData {
		if i < mid {
			firstBuy = firstBuy.Add(bd.BuyVolume)
			firstSell = firstSell.Add(bd.SellVolume)
		} else {
			secondBuy = secondBuy.Add(bd.BuyVolume)
			secondSell = secondSell.Add(bd.SellVolume)
		}
	}

	// Calculate buy ratio for each half
	firstTotal := firstBuy.Add(firstSell)
	secondTotal := secondBuy.Add(secondSell)

	if firstTotal.IsZero() && secondTotal.IsZero() {
		return 0
	}

	var firstRatio, secondRatio int64
	if firstTotal.IsPositive() {
		firstRatio = firstBuy.MulRaw(100).Quo(firstTotal).Int64()
	} else {
		firstRatio = 50
	}
	if secondTotal.IsPositive() {
		secondRatio = secondBuy.MulRaw(100).Quo(secondTotal).Int64()
	} else {
		secondRatio = 50
	}

	// Momentum = (secondRatio - firstRatio) * 2, clamped
	score := (secondRatio - firstRatio) * 2
	return clampOracleScore(score)
}

// ---------------------------------------------------------------------------
// Oracle Update — called every 5 blocks alongside UpdateSentiment
// ---------------------------------------------------------------------------

// UpdateOracleSignals recomputes oracle signals for a pool.
func (k Keeper) UpdateOracleSignals(ctx context.Context, poolID uint64) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentHeight := sdkCtx.BlockHeight()

	// Gather block trading data
	fromHeight := currentHeight - OracleVolumeWindow
	if fromHeight < 1 {
		fromHeight = 1
	}
	blockData := k.GetBlockTradingDataRange(ctx, poolID, fromHeight, currentHeight-1)

	// Compute each signal
	signals := []OracleSignal{
		{
			Name:      "volume_spike",
			Score:     k.calcVolumeSpike(blockData, currentHeight),
			Weight:    oracleSignalWeights["volume_spike"],
			UpdatedAt: currentHeight,
		},
		{
			Name:      "whale_movement",
			Score:     k.calcWhaleMovement(blockData),
			Weight:    oracleSignalWeights["whale_movement"],
			UpdatedAt: currentHeight,
		},
		{
			Name:      "liquidity_change",
			Score:     k.calcLiquidityChange(blockData, currentHeight),
			Weight:    oracleSignalWeights["liquidity_change"],
			UpdatedAt: currentHeight,
		},
		{
			Name:      "orderbook_imbalance",
			Score:     k.calcOrderBookImbalance(ctx, poolID),
			Weight:    oracleSignalWeights["orderbook_imbalance"],
			UpdatedAt: currentHeight,
		},
		{
			Name:      "price_momentum",
			Score:     k.calcPriceMomentum(ctx, poolID, currentHeight),
			Weight:    oracleSignalWeights["price_momentum"],
			UpdatedAt: currentHeight,
		},
	}

	// Compute weighted composite score
	var weightedSum, totalWeight int64
	for _, s := range signals {
		weightedSum += s.Score * s.Weight
		totalWeight += s.Weight
	}

	compositeScore := int64(0)
	if totalWeight > 0 {
		compositeScore = weightedSum / totalWeight
	}
	compositeScore = clampOracleScore(compositeScore)

	// Determine trend
	trend := "neutral"
	if compositeScore > 25 {
		trend = "bullish"
	} else if compositeScore < -25 {
		trend = "bearish"
	}

	// Get recent whale alerts
	whaleAlerts := k.GetWhaleAlerts(ctx, poolID, 10)
	if whaleAlerts == nil {
		whaleAlerts = []WhaleAlert{}
	}

	composite := OracleComposite{
		PoolID:         poolID,
		CompositeScore: compositeScore,
		Trend:          trend,
		Signals:        signals,
		WhaleAlerts:    whaleAlerts,
		UpdatedAt:      currentHeight,
	}
	k.SetOracleComposite(ctx, composite)
}

// UpdateAllOracleSignals updates oracle signals for all pools.
// Called in BeginBlock alongside UpdateSentiment.
func (k Keeper) UpdateAllOracleSignals(ctx context.Context) {
	pools := k.GetAllPools(ctx)
	for _, pool := range pools {
		k.UpdateOracleSignals(ctx, pool.ID)
		// Prune old whale alerts to prevent unbounded state growth
		k.PruneOldWhaleAlerts(ctx, pool.ID)
	}
}

// GetOracleSignalByName returns a specific signal by name for a pool.
func (k Keeper) GetOracleSignalByName(ctx context.Context, poolID uint64, signalName string) (OracleSignal, bool) {
	composite := k.GetOracleComposite(ctx, poolID)
	for _, s := range composite.Signals {
		if s.Name == signalName {
			return s, true
		}
	}
	return OracleSignal{}, false
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// clampOracleScore clamps a value to -100..+100.
func clampOracleScore(v int64) int64 {
	if v < -100 {
		return -100
	}
	if v > 100 {
		return 100
	}
	return v
}

// CheckAndRecordWhaleAlert checks if a swap qualifies as a whale alert
// and records it. Called from Swap when tokenIn.Amount > 5% of reserveIn.
func (k Keeper) CheckAndRecordWhaleAlert(ctx context.Context, sender string, poolID uint64, tokenIn sdk.Coin, pool types.Pool) {
	var reserveIn math.Int
	var side string

	if tokenIn.Denom == pool.DenomA {
		reserveIn = pool.ReserveA
		side = "sell" // selling DenomA for DenomB
	} else if tokenIn.Denom == pool.DenomB {
		reserveIn = pool.ReserveB
		side = "buy" // selling DenomB for DenomA (buying DenomA)
	} else {
		return
	}

	// Check if amount > 5% of reserve
	threshold := reserveIn.MulRaw(WhaleAlertThresholdBps).QuoRaw(10000)
	if tokenIn.Amount.GTE(threshold) {
		k.RecordWhaleAlert(ctx, sender, poolID, tokenIn.Amount, side)
	}
}

// PruneOldWhaleAlerts removes whale alerts older than SentimentWindow blocks.
func (k Keeper) PruneOldWhaleAlerts(ctx context.Context, poolID uint64) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cutoff := sdkCtx.BlockHeight() - SentimentWindow
	if cutoff < 1 {
		return
	}

	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := types.WhaleAlertPoolPrefix(poolID)
	// Iterate from beginning to cutoff height
	endKey := types.WhaleAlertKey(poolID, cutoff, 0)

	iter, err := kvStore.Iterator(prefix, endKey)
	if err != nil {
		return
	}
	defer iter.Close()

	var keysToDelete [][]byte
	for ; iter.Valid(); iter.Next() {
		keysToDelete = append(keysToDelete, append([]byte{}, iter.Key()...))
	}

	for _, key := range keysToDelete {
		kvStore.Delete(key)
	}
}

// suppress unused import
var _ = fmt.Sprintf
