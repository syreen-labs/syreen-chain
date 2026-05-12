package keeper

import (
	"context"
	"encoding/json"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/dex/types"
)

// ---------------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------------

const (
	signalStatePrefix   = "signal_state/"
	signalHistoryPrefix = "signal_history/"

	// Default price history depth (configurable per deployment)
	DefaultPriceHistoryDepth = 100

	// Moving average windows
	SMAWindowShort  = 10
	SMAWindowMedium = 20
	SMAWindowLong   = 50

	// RSI period
	RSIPeriod = 14

	// Signal history retention (how many block snapshots to keep)
	SignalHistoryRetention = 200
)

// Signal rating thresholds (score 0-100)
const (
	StrongSellThreshold = 20
	SellThreshold       = 40
	NeutralLow          = 40
	NeutralHigh         = 60
	BuyThreshold        = 60
	StrongBuyThreshold  = 80
)

// Signal rating strings
const (
	SignalStrongBuy  = "STRONG_BUY"
	SignalBuy        = "BUY"
	SignalNeutral    = "NEUTRAL"
	SignalSell       = "SELL"
	SignalStrongSell = "STRONG_SELL"
)

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// PricePoint stores a price at a specific block height.
type PricePoint struct {
	Height int64              `json:"height"`
	Price  sdkmath.LegacyDec `json:"price"`
	Volume sdkmath.Int        `json:"volume"`
}

// MovingAverages holds SMA and EMA values for different windows.
type MovingAverages struct {
	SMA10 sdkmath.LegacyDec `json:"sma_10"`
	SMA20 sdkmath.LegacyDec `json:"sma_20"`
	SMA50 sdkmath.LegacyDec `json:"sma_50"`
	EMA10 sdkmath.LegacyDec `json:"ema_10"`
	EMA20 sdkmath.LegacyDec `json:"ema_20"`
	EMA50 sdkmath.LegacyDec `json:"ema_50"`
}

// MomentumIndicators holds momentum-related signals.
type MomentumIndicators struct {
	RSI            sdkmath.LegacyDec `json:"rsi"`
	PriceROC       sdkmath.LegacyDec `json:"price_roc"`        // % change over last N blocks
	VolumeMomentum sdkmath.LegacyDec `json:"volume_momentum"`  // current volume / avg volume
}

// PoolSignalState stores the full signal computation state for a pool.
type PoolSignalState struct {
	PoolID       uint64       `json:"pool_id"`
	PriceHistory []PricePoint `json:"price_history"`
	LastHeight   int64        `json:"last_height"`
	// EMA state: carry forward for incremental computation
	EMA10 sdkmath.LegacyDec `json:"ema_10"`
	EMA20 sdkmath.LegacyDec `json:"ema_20"`
	EMA50 sdkmath.LegacyDec `json:"ema_50"`
}

// PoolSignals is the public-facing signal output for a pool.
type PoolSignals struct {
	PoolID          uint64             `json:"pool_id"`
	Height          int64              `json:"height"`
	CurrentPrice    sdkmath.LegacyDec  `json:"current_price"`
	MovingAverages  MovingAverages     `json:"moving_averages"`
	Momentum        MomentumIndicators `json:"momentum"`
	VolatilityScore sdkmath.LegacyDec  `json:"volatility_score"`
	CompositeScore  int64              `json:"composite_score"` // 0-100
	Signal          string             `json:"signal"`          // STRONG_BUY, BUY, NEUTRAL, SELL, STRONG_SELL
}

// SignalHistoryEntry stores a snapshot of signals at a given height.
type SignalHistoryEntry struct {
	Height         int64              `json:"height"`
	Price          sdkmath.LegacyDec  `json:"price"`
	CompositeScore int64              `json:"composite_score"`
	Signal         string             `json:"signal"`
	RSI            sdkmath.LegacyDec  `json:"rsi"`
	Volatility     sdkmath.LegacyDec  `json:"volatility"`
}

// SignalHistory holds the signal history for a pool.
type SignalHistory struct {
	PoolID  uint64               `json:"pool_id"`
	Entries []SignalHistoryEntry  `json:"entries"`
}

// ---------------------------------------------------------------------------
// Store keys
// ---------------------------------------------------------------------------

func signalStateKey(poolID uint64) []byte {
	return append([]byte(signalStatePrefix), types.PoolKey(poolID)...)
}

func signalHistoryKey(poolID uint64) []byte {
	return append([]byte(signalHistoryPrefix), types.PoolKey(poolID)...)
}

// ---------------------------------------------------------------------------
// State CRUD
// ---------------------------------------------------------------------------

func (k Keeper) GetSignalState(ctx context.Context, poolID uint64) (PoolSignalState, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(signalStateKey(poolID))
	if err != nil || bz == nil {
		return PoolSignalState{}, false
	}
	var state PoolSignalState
	if err := json.Unmarshal(bz, &state); err != nil {
		return PoolSignalState{}, false
	}
	return state, true
}

func (k Keeper) SetSignalState(ctx context.Context, state PoolSignalState) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(state)
	kvStore.Set(signalStateKey(state.PoolID), bz)
}

func (k Keeper) GetSignalHistory(ctx context.Context, poolID uint64) (SignalHistory, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(signalHistoryKey(poolID))
	if err != nil || bz == nil {
		return SignalHistory{}, false
	}
	var history SignalHistory
	if err := json.Unmarshal(bz, &history); err != nil {
		return SignalHistory{}, false
	}
	return history, true
}

func (k Keeper) SetSignalHistory(ctx context.Context, history SignalHistory) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(history)
	kvStore.Set(signalHistoryKey(history.PoolID), bz)
}

// ---------------------------------------------------------------------------
// BeginBlock hook: UpdateAllSignals
// ---------------------------------------------------------------------------

// UpdateAllSignals recalculates trading signals for all pools. Called in BeginBlock.
func (k Keeper) UpdateAllSignals(ctx context.Context) {
	pools := k.GetAllPools(ctx)
	for _, pool := range pools {
		k.UpdatePoolSignals(ctx, pool)
	}
}

// UpdatePoolSignals updates signals for a single pool.
func (k Keeper) UpdatePoolSignals(ctx context.Context, pool types.Pool) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	height := sdkCtx.BlockHeight()

	// Calculate current price: price of DenomA in terms of DenomB
	if pool.ReserveA.IsZero() || pool.ReserveB.IsZero() {
		return
	}
	currentPrice := sdkmath.LegacyNewDecFromInt(pool.ReserveB).Quo(sdkmath.LegacyNewDecFromInt(pool.ReserveA))

	// Get or initialize signal state
	state, found := k.GetSignalState(ctx, pool.ID)
	if !found {
		state = PoolSignalState{
			PoolID:       pool.ID,
			PriceHistory: []PricePoint{},
			LastHeight:   height,
			EMA10:        currentPrice,
			EMA20:        currentPrice,
			EMA50:        currentPrice,
		}
	}

	// Skip if we already processed this height
	if found && state.LastHeight >= height {
		return
	}

	// Get current volume from risk state (if available)
	volume := sdkmath.ZeroInt()
	riskState, riskFound := k.GetRiskState(ctx, pool.ID)
	if riskFound {
		volume = riskState.VolumeThisBlock
	}

	// Append new price point
	state.PriceHistory = append(state.PriceHistory, PricePoint{
		Height: height,
		Price:  currentPrice,
		Volume: volume,
	})

	// Trim to max depth
	if len(state.PriceHistory) > DefaultPriceHistoryDepth {
		state.PriceHistory = state.PriceHistory[len(state.PriceHistory)-DefaultPriceHistoryDepth:]
	}

	// Update EMAs incrementally
	state.EMA10 = calcEMA(currentPrice, state.EMA10, SMAWindowShort)
	state.EMA20 = calcEMA(currentPrice, state.EMA20, SMAWindowMedium)
	state.EMA50 = calcEMA(currentPrice, state.EMA50, SMAWindowLong)

	state.LastHeight = height
	k.SetSignalState(ctx, state)

	// Compute full signals and store history snapshot
	signals := k.ComputePoolSignals(ctx, pool.ID)
	if signals != nil {
		k.appendSignalHistory(ctx, pool.ID, *signals)
	}
}

// ---------------------------------------------------------------------------
// Signal Computation
// ---------------------------------------------------------------------------

// ComputePoolSignals computes all trading signals for a pool from stored state.
func (k Keeper) ComputePoolSignals(ctx context.Context, poolID uint64) *PoolSignals {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	state, found := k.GetSignalState(ctx, poolID)
	if !found || len(state.PriceHistory) == 0 {
		return nil
	}

	currentPrice := state.PriceHistory[len(state.PriceHistory)-1].Price
	prices := extractPrices(state.PriceHistory)

	// Moving averages
	ma := MovingAverages{
		SMA10: calcSMA(prices, SMAWindowShort),
		SMA20: calcSMA(prices, SMAWindowMedium),
		SMA50: calcSMA(prices, SMAWindowLong),
		EMA10: state.EMA10,
		EMA20: state.EMA20,
		EMA50: state.EMA50,
	}

	// Momentum indicators
	rsi := calcRSI(prices, RSIPeriod)
	priceROC := calcPriceROC(prices, RSIPeriod)
	volumeMomentum := calcVolumeMomentum(state.PriceHistory)

	momentum := MomentumIndicators{
		RSI:            rsi,
		PriceROC:       priceROC,
		VolumeMomentum: volumeMomentum,
	}

	// Volatility
	volatility := calcVolatility(prices, SMAWindowMedium)

	// Composite score
	score := calcCompositeScore(currentPrice, ma, momentum, volatility)
	signal := scoreToSignal(score)

	return &PoolSignals{
		PoolID:          poolID,
		Height:          sdkCtx.BlockHeight(),
		CurrentPrice:    currentPrice,
		MovingAverages:  ma,
		Momentum:        momentum,
		VolatilityScore: volatility,
		CompositeScore:  score,
		Signal:          signal,
	}
}

// ---------------------------------------------------------------------------
// Signal History
// ---------------------------------------------------------------------------

func (k Keeper) appendSignalHistory(ctx context.Context, poolID uint64, signals PoolSignals) {
	history, found := k.GetSignalHistory(ctx, poolID)
	if !found {
		history = SignalHistory{
			PoolID:  poolID,
			Entries: []SignalHistoryEntry{},
		}
	}

	entry := SignalHistoryEntry{
		Height:         signals.Height,
		Price:          signals.CurrentPrice,
		CompositeScore: signals.CompositeScore,
		Signal:         signals.Signal,
		RSI:            signals.Momentum.RSI,
		Volatility:     signals.VolatilityScore,
	}

	history.Entries = append(history.Entries, entry)

	// Trim to retention limit
	if len(history.Entries) > SignalHistoryRetention {
		history.Entries = history.Entries[len(history.Entries)-SignalHistoryRetention:]
	}

	k.SetSignalHistory(ctx, history)
}

// ---------------------------------------------------------------------------
// Math helpers — all pure functions, no allocations beyond LegacyDec
// ---------------------------------------------------------------------------

// extractPrices extracts just the price decimals from price points.
func extractPrices(points []PricePoint) []sdkmath.LegacyDec {
	prices := make([]sdkmath.LegacyDec, len(points))
	for i, p := range points {
		prices[i] = p.Price
	}
	return prices
}

// calcSMA computes Simple Moving Average over the last `window` values.
// If fewer values than window, uses all available.
func calcSMA(prices []sdkmath.LegacyDec, window int) sdkmath.LegacyDec {
	if len(prices) == 0 {
		return sdkmath.LegacyZeroDec()
	}

	n := window
	if n > len(prices) {
		n = len(prices)
	}

	sum := sdkmath.LegacyZeroDec()
	start := len(prices) - n
	for i := start; i < len(prices); i++ {
		sum = sum.Add(prices[i])
	}
	return sum.Quo(sdkmath.LegacyNewDec(int64(n)))
}

// calcEMA computes Exponential Moving Average incrementally.
// EMA = price * k + prevEMA * (1 - k), where k = 2 / (window + 1)
func calcEMA(currentPrice, prevEMA sdkmath.LegacyDec, window int) sdkmath.LegacyDec {
	// k = 2 / (window + 1)
	k := sdkmath.LegacyNewDec(2).Quo(sdkmath.LegacyNewDec(int64(window + 1)))
	oneMinusK := sdkmath.LegacyOneDec().Sub(k)
	return currentPrice.Mul(k).Add(prevEMA.Mul(oneMinusK))
}

// calcRSI computes the Relative Strength Index over the last `period` price changes.
// RSI = 100 - 100 / (1 + RS), where RS = avgGain / avgLoss
func calcRSI(prices []sdkmath.LegacyDec, period int) sdkmath.LegacyDec {
	if len(prices) < 2 {
		return sdkmath.LegacyNewDec(50) // neutral when insufficient data
	}

	n := period
	if n > len(prices)-1 {
		n = len(prices) - 1
	}

	avgGain := sdkmath.LegacyZeroDec()
	avgLoss := sdkmath.LegacyZeroDec()

	start := len(prices) - n - 1
	if start < 0 {
		start = 0
	}

	count := 0
	for i := start + 1; i < len(prices); i++ {
		change := prices[i].Sub(prices[i-1])
		if change.IsPositive() {
			avgGain = avgGain.Add(change)
		} else if change.IsNegative() {
			avgLoss = avgLoss.Add(change.Abs())
		}
		count++
	}

	if count == 0 {
		return sdkmath.LegacyNewDec(50)
	}

	periodDec := sdkmath.LegacyNewDec(int64(count))
	avgGain = avgGain.Quo(periodDec)
	avgLoss = avgLoss.Quo(periodDec)

	if avgLoss.IsZero() {
		if avgGain.IsZero() {
			return sdkmath.LegacyNewDec(50) // no movement
		}
		return sdkmath.LegacyNewDec(100) // all gains, no losses
	}

	rs := avgGain.Quo(avgLoss)
	rsi := sdkmath.LegacyNewDec(100).Sub(
		sdkmath.LegacyNewDec(100).Quo(sdkmath.LegacyOneDec().Add(rs)),
	)

	return rsi
}

// calcPriceROC computes rate of change: (current - past) / past * 100
func calcPriceROC(prices []sdkmath.LegacyDec, lookback int) sdkmath.LegacyDec {
	if len(prices) < 2 {
		return sdkmath.LegacyZeroDec()
	}

	n := lookback
	if n >= len(prices) {
		n = len(prices) - 1
	}

	current := prices[len(prices)-1]
	past := prices[len(prices)-1-n]

	if past.IsZero() {
		return sdkmath.LegacyZeroDec()
	}

	return current.Sub(past).Quo(past).MulInt64(100)
}

// calcVolumeMomentum returns current block volume / average volume.
// > 1.0 means above average, < 1.0 means below average.
func calcVolumeMomentum(points []PricePoint) sdkmath.LegacyDec {
	if len(points) == 0 {
		return sdkmath.LegacyOneDec()
	}

	totalVolume := sdkmath.LegacyZeroDec()
	for _, p := range points {
		totalVolume = totalVolume.Add(sdkmath.LegacyNewDecFromInt(p.Volume))
	}

	avgVolume := totalVolume.Quo(sdkmath.LegacyNewDec(int64(len(points))))
	if avgVolume.IsZero() {
		return sdkmath.LegacyOneDec()
	}

	currentVolume := sdkmath.LegacyNewDecFromInt(points[len(points)-1].Volume)
	return currentVolume.Quo(avgVolume)
}

// calcVolatility computes the standard deviation of price percentage changes.
func calcVolatility(prices []sdkmath.LegacyDec, window int) sdkmath.LegacyDec {
	if len(prices) < 2 {
		return sdkmath.LegacyZeroDec()
	}

	n := window
	if n > len(prices)-1 {
		n = len(prices) - 1
	}

	// Calculate percentage changes
	changes := make([]sdkmath.LegacyDec, 0, n)
	start := len(prices) - n - 1
	if start < 0 {
		start = 0
	}
	for i := start + 1; i < len(prices); i++ {
		if prices[i-1].IsZero() {
			continue
		}
		pctChange := prices[i].Sub(prices[i-1]).Quo(prices[i-1])
		changes = append(changes, pctChange)
	}

	if len(changes) == 0 {
		return sdkmath.LegacyZeroDec()
	}

	// Mean of changes
	mean := sdkmath.LegacyZeroDec()
	for _, c := range changes {
		mean = mean.Add(c)
	}
	mean = mean.Quo(sdkmath.LegacyNewDec(int64(len(changes))))

	// Variance = sum((x - mean)^2) / n
	variance := sdkmath.LegacyZeroDec()
	for _, c := range changes {
		diff := c.Sub(mean)
		variance = variance.Add(diff.Mul(diff))
	}
	variance = variance.Quo(sdkmath.LegacyNewDec(int64(len(changes))))

	// Standard deviation = sqrt(variance)
	// Use Newton's method for decimal square root
	return decSqrt(variance)
}

// decSqrt computes the square root of a LegacyDec using pure Newton's method.
// No float64 conversion — fully deterministic across all architectures.
func decSqrt(x sdkmath.LegacyDec) sdkmath.LegacyDec {
	if !x.IsPositive() {
		return sdkmath.LegacyZeroDec()
	}
	// Newton's method: guess = x/2, iterate guess = (guess + x/guess) / 2
	two := sdkmath.LegacyNewDec(2)
	guess := x.Quo(two)
	epsilon := sdkmath.LegacyNewDecWithPrec(1, 18)
	for i := 0; i < 50; i++ {
		if guess.IsZero() {
			return sdkmath.LegacyZeroDec()
		}
		next := guess.Add(x.Quo(guess)).Quo(two)
		if next.Sub(guess).Abs().LT(epsilon) {
			return next
		}
		guess = next
	}
	return guess
}

// ---------------------------------------------------------------------------
// Composite Score Calculation
// ---------------------------------------------------------------------------

// calcCompositeScore combines all indicators into a 0-100 score.
// 50 = neutral. > 50 = bullish. < 50 = bearish.
// Uses only math.LegacyDec — fully deterministic across all architectures.
func calcCompositeScore(
	currentPrice sdkmath.LegacyDec,
	ma MovingAverages,
	momentum MomentumIndicators,
	volatility sdkmath.LegacyDec,
) int64 {
	// Start at neutral (50)
	score := sdkmath.LegacyNewDec(50)

	// --- MA Signals (weight: 30%) ---
	maScore := sdkmath.LegacyZeroDec()

	// Price vs SMA crossovers
	if !ma.SMA10.IsZero() {
		if currentPrice.GT(ma.SMA10) {
			maScore = maScore.Add(sdkmath.LegacyNewDec(2))
		} else if currentPrice.LT(ma.SMA10) {
			maScore = maScore.Sub(sdkmath.LegacyNewDec(2))
		}
	}
	if !ma.SMA20.IsZero() {
		if currentPrice.GT(ma.SMA20) {
			maScore = maScore.Add(sdkmath.LegacyNewDec(3))
		} else if currentPrice.LT(ma.SMA20) {
			maScore = maScore.Sub(sdkmath.LegacyNewDec(3))
		}
	}
	if !ma.SMA50.IsZero() {
		if currentPrice.GT(ma.SMA50) {
			maScore = maScore.Add(sdkmath.LegacyNewDec(5))
		} else if currentPrice.LT(ma.SMA50) {
			maScore = maScore.Sub(sdkmath.LegacyNewDec(5))
		}
	}

	// EMA crossovers
	if !ma.EMA10.IsZero() && !ma.EMA20.IsZero() {
		if ma.EMA10.GT(ma.EMA20) {
			maScore = maScore.Add(sdkmath.LegacyNewDec(3))
		} else if ma.EMA10.LT(ma.EMA20) {
			maScore = maScore.Sub(sdkmath.LegacyNewDec(3))
		}
	}
	if !ma.EMA20.IsZero() && !ma.EMA50.IsZero() {
		if ma.EMA20.GT(ma.EMA50) {
			maScore = maScore.Add(sdkmath.LegacyNewDec(2))
		} else if ma.EMA20.LT(ma.EMA50) {
			maScore = maScore.Sub(sdkmath.LegacyNewDec(2))
		}
	}

	score = score.Add(maScore)

	// --- RSI Signal (weight: 25%) ---
	rsiVal := momentum.RSI
	thirty := sdkmath.LegacyNewDec(30)
	seventy := sdkmath.LegacyNewDec(70)
	fifty := sdkmath.LegacyNewDec(50)
	twelvePointFive := sdkmath.LegacyNewDecWithPrec(125, 1) // 12.5
	pointOneTwoFive := sdkmath.LegacyNewDecWithPrec(125, 3) // 0.125

	if rsiVal.LT(thirty) {
		// Oversold -> bullish signal
		score = score.Add(twelvePointFive.Mul(thirty.Sub(rsiVal)).Quo(thirty))
	} else if rsiVal.GT(seventy) {
		// Overbought -> bearish signal
		score = score.Sub(twelvePointFive.Mul(rsiVal.Sub(seventy)).Quo(thirty))
	}
	if rsiVal.GTE(thirty) && rsiVal.LTE(seventy) {
		score = score.Add(rsiVal.Sub(fifty).Mul(pointOneTwoFive))
	}

	// --- Price ROC Signal (weight: 20%) ---
	rocVal := momentum.PriceROC
	pointFive := sdkmath.LegacyNewDecWithPrec(5, 1) // 0.5
	ten := sdkmath.LegacyNewDec(10)
	rocContribution := rocVal.Mul(pointFive)
	if rocContribution.GT(ten) {
		rocContribution = ten
	} else if rocContribution.LT(ten.Neg()) {
		rocContribution = ten.Neg()
	}
	score = score.Add(rocContribution)

	// --- Volume Momentum Signal (weight: 10%) ---
	volMom := momentum.VolumeMomentum
	onePointFive := sdkmath.LegacyNewDecWithPrec(15, 1) // 1.5
	five := sdkmath.LegacyNewDec(5)
	if volMom.GT(onePointFive) {
		if rocVal.IsPositive() {
			score = score.Add(five)
		} else if rocVal.IsNegative() {
			score = score.Sub(five)
		}
	}

	// --- Volatility Penalty (weight: 15%) ---
	pointZeroFive := sdkmath.LegacyNewDecWithPrec(5, 2) // 0.05
	if volatility.GT(pointZeroFive) {
		pull := volatility.MulInt64(50)
		if pull.GT(ten) {
			pull = ten
		}
		if score.GT(fifty) {
			score = score.Sub(pull)
		} else if score.LT(fifty) {
			score = score.Add(pull)
		}
	}

	// Clamp to 0-100
	zero := sdkmath.LegacyZeroDec()
	hundred := sdkmath.LegacyNewDec(100)
	if score.LT(zero) {
		score = zero
	}
	if score.GT(hundred) {
		score = hundred
	}

	return score.TruncateInt64()
}

// scoreToSignal converts a numeric score to a signal string.
func scoreToSignal(score int64) string {
	switch {
	case score >= StrongBuyThreshold:
		return SignalStrongBuy
	case score >= BuyThreshold:
		return SignalBuy
	case score > NeutralLow && score < NeutralHigh:
		return SignalNeutral
	case score <= StrongSellThreshold:
		return SignalStrongSell
	case score <= SellThreshold:
		return SignalSell
	default:
		return SignalNeutral
	}
}

// GetSignalVolatility returns the current volatility score for a pool (for cross-module use).
// Returns zero if no signal data exists yet.
func (k Keeper) GetSignalVolatility(ctx context.Context, poolID uint64) sdkmath.LegacyDec {
	signals := k.ComputePoolSignals(ctx, poolID)
	if signals == nil {
		return sdkmath.LegacyZeroDec()
	}
	return signals.VolatilityScore
}

// GetSignalForPool returns (signal string, composite score 0-100, RSI) for cross-module use.
func (k Keeper) GetSignalForPool(ctx context.Context, poolID uint64) (string, int64, sdkmath.LegacyDec) {
	signals := k.ComputePoolSignals(ctx, poolID)
	if signals == nil {
		return "NEUTRAL", 50, sdkmath.LegacyNewDec(50)
	}
	return signals.Signal, signals.CompositeScore, signals.Momentum.RSI
}
