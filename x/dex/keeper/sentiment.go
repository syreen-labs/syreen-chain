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
// Sentiment Oracle — derives market sentiment from on-chain trading patterns
// ---------------------------------------------------------------------------

const (
	// SentimentWindow is the number of past blocks to analyze for sentiment.
	SentimentWindow = 200

	// WhaleThresholdBps is the minimum trade size as basis points of pool reserves
	// to qualify as a whale trade (100 bps = 1%).
	WhaleThresholdBps = 100
)

// MarketMood describes the overall market mood for a pool.
type MarketMood string

const (
	MoodBullish         MarketMood = "BULLISH"
	MoodSlightlyBullish MarketMood = "SLIGHTLY_BULLISH"
	MoodNeutral         MarketMood = "NEUTRAL"
	MoodSlightlyBearish MarketMood = "SLIGHTLY_BEARISH"
	MoodBearish         MarketMood = "BEARISH"
)

// AlertType classifies a sentiment alert.
type AlertType string

const (
	AlertFearSpike          AlertType = "fear_spike"
	AlertGreedSpike         AlertType = "greed_spike"
	AlertWhaleAccumulation  AlertType = "whale_accumulation"
	AlertLiquidityExodus    AlertType = "liquidity_exodus"
)

// ---------------------------------------------------------------------------
// BlockTradingData — per-pool per-block counters stored in KV store
// ---------------------------------------------------------------------------

// BlockTradingData holds trading activity counters for a single block for a single pool.
type BlockTradingData struct {
	PoolID          uint64   `json:"pool_id"`
	Height          int64    `json:"height"`
	BuyVolume       math.Int `json:"buy_volume"`        // volume of DenomA bought (DenomB sold)
	SellVolume      math.Int `json:"sell_volume"`       // volume of DenomA sold (DenomB bought)
	BuyCount        uint64   `json:"buy_count"`
	SellCount       uint64   `json:"sell_count"`
	LPAddVolume     math.Int `json:"lp_add_volume"`     // total value added as liquidity
	LPRemoveVolume  math.Int `json:"lp_remove_volume"`  // total value removed from liquidity
	LPAddCount      uint64   `json:"lp_add_count"`
	LPRemoveCount   uint64   `json:"lp_remove_count"`
	WhaleBuyCount   uint64   `json:"whale_buy_count"`   // large buys (>1% of pool)
	WhaleSellCount  uint64   `json:"whale_sell_count"`  // large sells (>1% of pool)
	WhaleBuyVolume  math.Int `json:"whale_buy_volume"`
	WhaleSellVolume math.Int `json:"whale_sell_volume"`
}

// NewBlockTradingData creates a zero-valued BlockTradingData for a pool at a height.
func NewBlockTradingData(poolID uint64, height int64) BlockTradingData {
	return BlockTradingData{
		PoolID:          poolID,
		Height:          height,
		BuyVolume:       math.ZeroInt(),
		SellVolume:      math.ZeroInt(),
		LPAddVolume:     math.ZeroInt(),
		LPRemoveVolume:  math.ZeroInt(),
		WhaleBuyVolume:  math.ZeroInt(),
		WhaleSellVolume: math.ZeroInt(),
	}
}

// ---------------------------------------------------------------------------
// SentimentIndicators — individual signal scores
// ---------------------------------------------------------------------------

// SentimentIndicators holds the individual indicator scores (each 0-100).
type SentimentIndicators struct {
	BuySellRatio   int64 `json:"buy_sell_ratio"`    // >50 = more buying
	TradeSizeTrend int64 `json:"trade_size_trend"`  // >50 = sizes increasing
	IntentSentiment int64 `json:"intent_sentiment"` // >50 = more bullish intents
	LiquidityFlow  int64 `json:"liquidity_flow"`    // >50 = net inflow
	WhaleActivity  int64 `json:"whale_activity"`    // >50 = whales buying
	VolumeTrend    int64 `json:"volume_trend"`      // >50 = volume increasing
}

// ---------------------------------------------------------------------------
// SentimentState — the current computed sentiment for a pool
// ---------------------------------------------------------------------------

// SentimentState holds the full sentiment analysis for a pool.
type SentimentState struct {
	PoolID          uint64              `json:"pool_id"`
	FearGreedIndex  int64               `json:"fear_greed_index"`  // 0-100
	FearGreedLabel  string              `json:"fear_greed_label"`  // "Extreme Fear", "Fear", etc.
	Mood            MarketMood          `json:"mood"`
	Indicators      SentimentIndicators `json:"indicators"`
	UpdatedAtHeight int64               `json:"updated_at_height"`
	BlocksAnalyzed  int64               `json:"blocks_analyzed"`   // how many blocks of data used
}

// ---------------------------------------------------------------------------
// SentimentAlert — a detected sentiment shift event
// ---------------------------------------------------------------------------

// SentimentAlert represents a detected sentiment event.
type SentimentAlert struct {
	PoolID    uint64    `json:"pool_id"`
	AlertType AlertType `json:"alert_type"`
	Message   string    `json:"message"`
	Severity  string    `json:"severity"` // "low", "medium", "high"
	Height    int64     `json:"height"`
	Value     int64     `json:"value"`    // the metric value that triggered the alert
}

// SentimentAlerts holds the list of active alerts for a pool.
type SentimentAlerts struct {
	PoolID uint64           `json:"pool_id"`
	Alerts []SentimentAlert `json:"alerts"`
}

// ---------------------------------------------------------------------------
// KV Store — BlockTradingData
// ---------------------------------------------------------------------------

func (k Keeper) GetBlockTradingData(ctx context.Context, poolID uint64, height int64) (BlockTradingData, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.SentimentBlockDataKey(poolID, height))
	if err != nil || bz == nil {
		return BlockTradingData{}, false
	}
	var data BlockTradingData
	if err := json.Unmarshal(bz, &data); err != nil {
		return BlockTradingData{}, false
	}
	return data, true
}

func (k Keeper) SetBlockTradingData(ctx context.Context, data BlockTradingData) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(data)
	kvStore.Set(types.SentimentBlockDataKey(data.PoolID, data.Height), bz)
}

// DeleteBlockTradingData removes old block data to keep only the last SentimentWindow blocks.
func (k Keeper) DeleteBlockTradingData(ctx context.Context, poolID uint64, height int64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	kvStore.Delete(types.SentimentBlockDataKey(poolID, height))
}

// GetBlockTradingDataRange returns block trading data for a pool over a range of heights.
// Uses a prefix iterator instead of individual Get calls per height, so only blocks
// that actually have stored data are visited (D-13 fix).
func (k Keeper) GetBlockTradingDataRange(ctx context.Context, poolID uint64, fromHeight, toHeight int64) []BlockTradingData {
	kvStore := k.storeService.OpenKVStore(ctx)

	// Build start and end keys for the range within this pool's block data.
	startKey := types.SentimentBlockDataKey(poolID, fromHeight)
	// End key is exclusive: use toHeight+1 so toHeight is included.
	endKey := types.SentimentBlockDataKey(poolID, toHeight+1)

	iter, err := kvStore.Iterator(startKey, endKey)
	if err != nil {
		return nil
	}
	defer iter.Close()

	const maxResults = 1000
	var results []BlockTradingData
	for ; iter.Valid() && len(results) < maxResults; iter.Next() {
		var data BlockTradingData
		if err := json.Unmarshal(iter.Value(), &data); err != nil {
			continue
		}
		results = append(results, data)
	}
	return results
}

// ---------------------------------------------------------------------------
// KV Store — SentimentState
// ---------------------------------------------------------------------------

func (k Keeper) GetSentimentState(ctx context.Context, poolID uint64) (SentimentState, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.SentimentStateKey(poolID))
	if err != nil || bz == nil {
		return SentimentState{}, false
	}
	var state SentimentState
	if err := json.Unmarshal(bz, &state); err != nil {
		return SentimentState{}, false
	}
	return state, true
}

func (k Keeper) SetSentimentState(ctx context.Context, state SentimentState) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(state)
	kvStore.Set(types.SentimentStateKey(state.PoolID), bz)
}

// GetFearGreedIndex returns the current fear/greed index for a pool (for cross-module use).
// Returns (index, label, ok) where ok is false if no sentiment state exists yet.
func (k Keeper) GetFearGreedIndex(ctx context.Context, poolID uint64) (int64, string, bool) {
	state, ok := k.GetSentimentState(ctx, poolID)
	if !ok {
		return 0, "", false
	}
	return state.FearGreedIndex, state.FearGreedLabel, true
}

// ---------------------------------------------------------------------------
// KV Store — SentimentAlerts
// ---------------------------------------------------------------------------

func (k Keeper) GetSentimentAlerts(ctx context.Context, poolID uint64) SentimentAlerts {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.SentimentAlertKey(poolID))
	if err != nil || bz == nil {
		return SentimentAlerts{PoolID: poolID, Alerts: []SentimentAlert{}}
	}
	var alerts SentimentAlerts
	if err := json.Unmarshal(bz, &alerts); err != nil {
		return SentimentAlerts{PoolID: poolID, Alerts: []SentimentAlert{}}
	}
	return alerts
}

func (k Keeper) SetSentimentAlerts(ctx context.Context, alerts SentimentAlerts) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(alerts)
	kvStore.Set(types.SentimentAlertKey(alerts.PoolID), bz)
}

// ---------------------------------------------------------------------------
// Recording — called by Swap, AddLiquidity, RemoveLiquidity
// ---------------------------------------------------------------------------

// RecordSwapForSentiment records a swap event for sentiment tracking.
// isBuy: true if the trader is buying DenomA (selling DenomB), false if selling DenomA.
func (k Keeper) RecordSwapForSentiment(ctx context.Context, poolID uint64, amount math.Int, isBuy bool, pool types.Pool) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	height := sdkCtx.BlockHeight()

	data, found := k.GetBlockTradingData(ctx, poolID, height)
	if !found {
		data = NewBlockTradingData(poolID, height)
	}

	// Determine if this is a whale trade (>1% of smaller reserve)
	smallerReserve := pool.ReserveA
	if pool.ReserveB.LT(pool.ReserveA) {
		smallerReserve = pool.ReserveB
	}
	threshold := smallerReserve.MulRaw(WhaleThresholdBps).QuoRaw(10000)
	isWhale := amount.GTE(threshold)

	if isBuy {
		data.BuyVolume = data.BuyVolume.Add(amount)
		data.BuyCount++
		if isWhale {
			data.WhaleBuyCount++
			data.WhaleBuyVolume = data.WhaleBuyVolume.Add(amount)
		}
	} else {
		data.SellVolume = data.SellVolume.Add(amount)
		data.SellCount++
		if isWhale {
			data.WhaleSellCount++
			data.WhaleSellVolume = data.WhaleSellVolume.Add(amount)
		}
	}

	k.SetBlockTradingData(ctx, data)
}

// RecordLiquidityForSentiment records a liquidity add/remove event.
func (k Keeper) RecordLiquidityForSentiment(ctx context.Context, poolID uint64, amount math.Int, isAdd bool) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	height := sdkCtx.BlockHeight()

	data, found := k.GetBlockTradingData(ctx, poolID, height)
	if !found {
		data = NewBlockTradingData(poolID, height)
	}

	if isAdd {
		data.LPAddVolume = data.LPAddVolume.Add(amount)
		data.LPAddCount++
	} else {
		data.LPRemoveVolume = data.LPRemoveVolume.Add(amount)
		data.LPRemoveCount++
	}

	k.SetBlockTradingData(ctx, data)
}

// ---------------------------------------------------------------------------
// Computation — called in BeginBlock
// ---------------------------------------------------------------------------

// UpdateSentiment computes sentiment for all pools. Called in BeginBlock.
func (k Keeper) UpdateSentiment(ctx context.Context) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentHeight := sdkCtx.BlockHeight()

	pools := k.GetAllPools(ctx)
	for _, pool := range pools {
		k.computePoolSentiment(ctx, pool, currentHeight)
	}
}

// computePoolSentiment calculates sentiment for a single pool.
func (k Keeper) computePoolSentiment(ctx context.Context, pool types.Pool, currentHeight int64) {
	// Determine the analysis window
	fromHeight := currentHeight - SentimentWindow
	if fromHeight < 1 {
		fromHeight = 1
	}

	// Gather block data
	blockData := k.GetBlockTradingDataRange(ctx, pool.ID, fromHeight, currentHeight-1)

	// Prune old data beyond the window
	pruneHeight := fromHeight - 1
	if pruneHeight > 0 {
		k.DeleteBlockTradingData(ctx, pool.ID, pruneHeight)
	}

	// Aggregate statistics
	var (
		totalBuyVolume       = math.ZeroInt()
		totalSellVolume      = math.ZeroInt()
		totalBuyCount        uint64
		totalSellCount       uint64
		totalLPAddVolume     = math.ZeroInt()
		totalLPRemoveVolume  = math.ZeroInt()
		totalWhaleBuyCount   uint64
		totalWhaleSellCount  uint64
		totalWhaleBuyVolume  = math.ZeroInt()
		totalWhaleSellVolume = math.ZeroInt()
	)

	// For volume trend: split window into halves
	midHeight := fromHeight + (currentHeight-1-fromHeight)/2
	var firstHalfVolume, secondHalfVolume = math.ZeroInt(), math.ZeroInt()

	// For trade size trend: split into halves
	var firstHalfTradeCount, secondHalfTradeCount uint64
	var firstHalfTradeVol, secondHalfTradeVol = math.ZeroInt(), math.ZeroInt()

	for _, bd := range blockData {
		totalBuyVolume = totalBuyVolume.Add(bd.BuyVolume)
		totalSellVolume = totalSellVolume.Add(bd.SellVolume)
		totalBuyCount += bd.BuyCount
		totalSellCount += bd.SellCount
		totalLPAddVolume = totalLPAddVolume.Add(bd.LPAddVolume)
		totalLPRemoveVolume = totalLPRemoveVolume.Add(bd.LPRemoveVolume)
		totalWhaleBuyCount += bd.WhaleBuyCount
		totalWhaleSellCount += bd.WhaleSellCount
		totalWhaleBuyVolume = totalWhaleBuyVolume.Add(bd.WhaleBuyVolume)
		totalWhaleSellVolume = totalWhaleSellVolume.Add(bd.WhaleSellVolume)

		blockVol := bd.BuyVolume.Add(bd.SellVolume)
		blockTradeCount := bd.BuyCount + bd.SellCount

		if bd.Height <= midHeight {
			firstHalfVolume = firstHalfVolume.Add(blockVol)
			firstHalfTradeCount += blockTradeCount
			firstHalfTradeVol = firstHalfTradeVol.Add(blockVol)
		} else {
			secondHalfVolume = secondHalfVolume.Add(blockVol)
			secondHalfTradeCount += blockTradeCount
			secondHalfTradeVol = secondHalfTradeVol.Add(blockVol)
		}
	}

	blocksAnalyzed := int64(len(blockData))

	// Calculate individual indicators (each 0-100)
	indicators := SentimentIndicators{
		BuySellRatio:    calcBuySellRatio(totalBuyVolume, totalSellVolume),
		TradeSizeTrend:  calcTradeSizeTrend(firstHalfTradeVol, firstHalfTradeCount, secondHalfTradeVol, secondHalfTradeCount),
		IntentSentiment: 50, // default neutral — updated below if intent data available
		LiquidityFlow:   calcLiquidityFlow(totalLPAddVolume, totalLPRemoveVolume),
		WhaleActivity:   calcWhaleActivity(totalWhaleBuyVolume, totalWhaleSellVolume, totalWhaleBuyCount, totalWhaleSellCount),
		VolumeTrend:     calcVolumeTrend(firstHalfVolume, secondHalfVolume),
	}

	// Composite Fear & Greed Index (weighted average of indicators)
	// Weights: BuySellRatio=25%, TradeSizeTrend=10%, IntentSentiment=15%,
	//          LiquidityFlow=20%, WhaleActivity=20%, VolumeTrend=10%
	fearGreedIndex := (indicators.BuySellRatio*25 +
		indicators.TradeSizeTrend*10 +
		indicators.IntentSentiment*15 +
		indicators.LiquidityFlow*20 +
		indicators.WhaleActivity*20 +
		indicators.VolumeTrend*10) / 100

	// Clamp to 0-100
	if fearGreedIndex < 0 {
		fearGreedIndex = 0
	}
	if fearGreedIndex > 100 {
		fearGreedIndex = 100
	}

	label := fearGreedLabel(fearGreedIndex)
	mood := deriveMood(fearGreedIndex)

	state := SentimentState{
		PoolID:          pool.ID,
		FearGreedIndex:  fearGreedIndex,
		FearGreedLabel:  label,
		Mood:            mood,
		Indicators:      indicators,
		UpdatedAtHeight: currentHeight,
		BlocksAnalyzed:  blocksAnalyzed,
	}
	k.SetSentimentState(ctx, state)

	// Detect alerts
	k.detectSentimentAlerts(ctx, pool.ID, state, totalBuyVolume, totalSellVolume,
		totalWhaleBuyVolume, totalWhaleSellVolume, totalWhaleBuyCount, totalWhaleSellCount,
		totalLPAddVolume, totalLPRemoveVolume, currentHeight)
}

// ---------------------------------------------------------------------------
// Indicator Calculations (pure arithmetic, no external calls)
// ---------------------------------------------------------------------------

// calcBuySellRatio returns 0-100 where >50 means more buying.
func calcBuySellRatio(buyVol, sellVol math.Int) int64 {
	total := buyVol.Add(sellVol)
	if total.IsZero() {
		return 50 // neutral if no activity
	}
	// ratio = buyVol / total * 100
	ratio := buyVol.MulRaw(100).Quo(total)
	return clampScore(ratio.Int64())
}

// calcTradeSizeTrend returns 0-100 where >50 means average trade sizes are increasing.
func calcTradeSizeTrend(firstHalfVol math.Int, firstHalfCount uint64, secondHalfVol math.Int, secondHalfCount uint64) int64 {
	if firstHalfCount == 0 && secondHalfCount == 0 {
		return 50
	}

	var avgFirst, avgSecond int64

	if firstHalfCount > 0 {
		avgFirst = firstHalfVol.QuoRaw(int64(firstHalfCount)).Int64()
	}
	if secondHalfCount > 0 {
		avgSecond = secondHalfVol.QuoRaw(int64(secondHalfCount)).Int64()
	}

	if avgFirst == 0 && avgSecond == 0 {
		return 50
	}

	if avgFirst == 0 {
		return 75 // sizes appeared from nothing — bullish signal
	}

	// Calculate percentage change and map to 0-100 score
	// +100% change -> score 100, -100% change -> score 0, 0% change -> score 50
	changePct := ((avgSecond - avgFirst) * 50) / avgFirst
	return clampScore(50 + changePct)
}

// calcLiquidityFlow returns 0-100 where >50 means net LP inflow.
func calcLiquidityFlow(addVol, removeVol math.Int) int64 {
	total := addVol.Add(removeVol)
	if total.IsZero() {
		return 50 // neutral
	}
	// ratio = addVol / total * 100
	ratio := addVol.MulRaw(100).Quo(total)
	return clampScore(ratio.Int64())
}

// calcWhaleActivity returns 0-100 where >50 means whales are net buying.
func calcWhaleActivity(whaleBuyVol, whaleSellVol math.Int, whaleBuyCount, whaleSellCount uint64) int64 {
	totalVol := whaleBuyVol.Add(whaleSellVol)
	if totalVol.IsZero() {
		return 50 // neutral
	}
	ratio := whaleBuyVol.MulRaw(100).Quo(totalVol)
	return clampScore(ratio.Int64())
}

// calcVolumeTrend returns 0-100 where >50 means volume is increasing.
func calcVolumeTrend(firstHalfVol, secondHalfVol math.Int) int64 {
	if firstHalfVol.IsZero() && secondHalfVol.IsZero() {
		return 50
	}
	if firstHalfVol.IsZero() {
		return 75 // volume appeared from nothing
	}

	total := firstHalfVol.Add(secondHalfVol)
	// secondHalf fraction of total, mapped to 0-100
	// If second half = first half, score = 50
	// If second half >> first half, score -> 100
	ratio := secondHalfVol.MulRaw(100).Quo(total)
	return clampScore(ratio.Int64())
}

// ---------------------------------------------------------------------------
// Labels and Mood
// ---------------------------------------------------------------------------

func fearGreedLabel(index int64) string {
	switch {
	case index <= 20:
		return "Extreme Fear"
	case index <= 40:
		return "Fear"
	case index <= 60:
		return "Neutral"
	case index <= 80:
		return "Greed"
	default:
		return "Extreme Greed"
	}
}

func deriveMood(fearGreedIndex int64) MarketMood {
	switch {
	case fearGreedIndex <= 20:
		return MoodBearish
	case fearGreedIndex <= 40:
		return MoodSlightlyBearish
	case fearGreedIndex <= 60:
		return MoodNeutral
	case fearGreedIndex <= 80:
		return MoodSlightlyBullish
	default:
		return MoodBullish
	}
}

// ---------------------------------------------------------------------------
// Alert Detection
// ---------------------------------------------------------------------------

func (k Keeper) detectSentimentAlerts(
	ctx context.Context,
	poolID uint64,
	state SentimentState,
	totalBuyVol, totalSellVol math.Int,
	whaleBuyVol, whaleSellVol math.Int,
	whaleBuyCount, whaleSellCount uint64,
	lpAddVol, lpRemoveVol math.Int,
	height int64,
) {
	var alerts []SentimentAlert

	// Fear spike: Fear & Greed below 20 with significant selling
	if state.FearGreedIndex <= 20 && totalSellVol.GT(totalBuyVol) {
		alerts = append(alerts, SentimentAlert{
			PoolID:    poolID,
			AlertType: AlertFearSpike,
			Message:   fmt.Sprintf("Extreme fear detected — sell volume exceeds buy volume, F&G index at %d", state.FearGreedIndex),
			Severity:  "high",
			Height:    height,
			Value:     state.FearGreedIndex,
		})
	}

	// Greed spike: Fear & Greed above 80 with significant buying
	if state.FearGreedIndex >= 80 && totalBuyVol.GT(totalSellVol) {
		alerts = append(alerts, SentimentAlert{
			PoolID:    poolID,
			AlertType: AlertGreedSpike,
			Message:   fmt.Sprintf("Extreme greed detected — buy volume exceeds sell volume, F&G index at %d", state.FearGreedIndex),
			Severity:  "high",
			Height:    height,
			Value:     state.FearGreedIndex,
		})
	}

	// Whale accumulation: whales buying significantly more than selling
	if whaleBuyCount > 0 && whaleBuyCount >= whaleSellCount*3+1 {
		alerts = append(alerts, SentimentAlert{
			PoolID:    poolID,
			AlertType: AlertWhaleAccumulation,
			Message:   fmt.Sprintf("Whale accumulation pattern — %d whale buys vs %d whale sells", whaleBuyCount, whaleSellCount),
			Severity:  "medium",
			Height:    height,
			Value:     int64(whaleBuyCount),
		})
	}

	// Liquidity exodus: removal volume significantly exceeds additions
	if !lpRemoveVol.IsZero() && lpRemoveVol.GT(lpAddVol.MulRaw(3)) {
		alerts = append(alerts, SentimentAlert{
			PoolID:    poolID,
			AlertType: AlertLiquidityExodus,
			Message:   fmt.Sprintf("Liquidity exodus — LP removals (%s) far exceed additions (%s)", lpRemoveVol, lpAddVol),
			Severity:  "high",
			Height:    height,
			Value:     state.Indicators.LiquidityFlow,
		})
	}

	sentimentAlerts := SentimentAlerts{
		PoolID: poolID,
		Alerts: alerts,
	}
	k.SetSentimentAlerts(ctx, sentimentAlerts)

	// Log high-severity alerts
	for _, alert := range alerts {
		if alert.Severity == "high" {
			k.Logger(ctx).Info("SENTIMENT ALERT",
				"pool_id", poolID,
				"type", alert.AlertType,
				"message", alert.Message,
			)
		}
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// clampScore clamps a value to the 0-100 range.
func clampScore(v int64) int64 {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}

// GetSentimentHistory returns sentiment block data history for a pool over the last N blocks.
func (k Keeper) GetSentimentHistory(ctx context.Context, poolID uint64, limit int64) []BlockTradingData {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentHeight := sdkCtx.BlockHeight()

	if limit <= 0 || limit > SentimentWindow {
		limit = SentimentWindow
	}

	fromHeight := currentHeight - limit
	if fromHeight < 1 {
		fromHeight = 1
	}

	return k.GetBlockTradingDataRange(ctx, poolID, fromHeight, currentHeight)
}
