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
// Risk Score Constants
// ---------------------------------------------------------------------------

const (
	riskScorePrefix   = "risk_score/"
	riskHistoryPrefix = "risk_history/"
	maxHistoryEntries = 100 // keep last 100 blocks of risk history per pool

	// H6: cap the number of pools whose risk score we recompute per block.
	// On a chain with thousands of pools, scoring all of them every block is
	// O(N) work that can stall consensus. Pools are processed in round-robin
	// order using a stored cursor.
	MaxPoolsPerBlock = 500
)

// Risk levels
const (
	RiskLevelLow     = "LOW"
	RiskLevelMedium  = "MEDIUM"
	RiskLevelHigh    = "HIGH"
	RiskLevelExtreme = "EXTREME"
)

// Pool grades
const (
	PoolGradeA = "A"
	PoolGradeB = "B"
	PoolGradeC = "C"
	PoolGradeD = "D"
	PoolGradeF = "F"
)

// ---------------------------------------------------------------------------
// Trade Risk Assessment
// ---------------------------------------------------------------------------

// TradeRiskScore represents a pre-trade risk assessment.
type TradeRiskScore struct {
	PoolID           uint64         `json:"pool_id"`
	InputDenom       string         `json:"input_denom"`
	InputAmount      math.Int       `json:"input_amount"`
	PriceImpactScore int            `json:"price_impact_score"` // 0-100
	SlippageScore    int            `json:"slippage_score"`     // 0-100
	PoolHealthScore  int            `json:"pool_health_score"`  // 0-100
	TimingScore      int            `json:"timing_score"`       // 0-100
	OverallScore     int            `json:"overall_score"`      // 0-100 weighted average
	RiskLevel        string         `json:"risk_level"`         // LOW, MEDIUM, HIGH, EXTREME
	PriceImpactPct   math.LegacyDec `json:"price_impact_pct"`
	ExpectedOutput   math.Int       `json:"expected_output"`
	OutputDenom      string         `json:"output_denom"`
	Warnings         []string       `json:"warnings"`
}

// ---------------------------------------------------------------------------
// Pool Risk Score
// ---------------------------------------------------------------------------

// PoolRiskScore represents the overall risk assessment for a pool.
type PoolRiskScore struct {
	PoolID                uint64         `json:"pool_id"`
	LiquidityDepthScore   int            `json:"liquidity_depth_score"`   // 0-100 (100 = deep liquidity)
	VolatilityScore       int            `json:"volatility_score"`        // 0-100 (100 = very volatile = risky)
	VolumeConcentration   int            `json:"volume_concentration"`    // 0-100 (100 = one whale dominating)
	ReserveBalanceScore   int            `json:"reserve_balance_score"`   // 0-100 (100 = perfectly balanced)
	ImpermanentLossRisk   int            `json:"impermanent_loss_risk"`   // 0-100 (100 = high IL risk)
	OverallScore          int            `json:"overall_score"`           // 0-100 (0 = healthy, 100 = risky)
	PoolGrade             string         `json:"pool_grade"`              // A, B, C, D, F
	BlockHeight           int64          `json:"block_height"`
	IsHalted              bool           `json:"is_halted"`
	HaltReason            string         `json:"halt_reason,omitempty"`
	PriceChangeSinceStart math.LegacyDec `json:"price_change_since_start"`
}

// ---------------------------------------------------------------------------
// Risk History Entry
// ---------------------------------------------------------------------------

// RiskHistoryEntry stores a snapshot of a pool's risk score at a given block.
type RiskHistoryEntry struct {
	BlockHeight   int64  `json:"block_height"`
	OverallScore  int    `json:"overall_score"`
	PoolGrade     string `json:"pool_grade"`
	VolatilityScore int  `json:"volatility_score"`
	LiquidityScore  int  `json:"liquidity_score"`
	IsHalted      bool   `json:"is_halted"`
}

// RiskHistory holds a list of risk history entries for a pool.
type RiskHistory struct {
	PoolID  uint64             `json:"pool_id"`
	Entries []RiskHistoryEntry `json:"entries"`
}

// ---------------------------------------------------------------------------
// Store Keys
// ---------------------------------------------------------------------------

func riskScoreKey(poolID uint64) []byte {
	return append([]byte(riskScorePrefix), types.PoolKey(poolID)...)
}

func riskHistoryKey(poolID uint64) []byte {
	return append([]byte(riskHistoryPrefix), types.PoolKey(poolID)...)
}

// ---------------------------------------------------------------------------
// Store Getters/Setters
// ---------------------------------------------------------------------------

func (k Keeper) GetPoolRiskScore(ctx context.Context, poolID uint64) (PoolRiskScore, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(riskScoreKey(poolID))
	if err != nil || bz == nil {
		return PoolRiskScore{}, false
	}
	var score PoolRiskScore
	if err := json.Unmarshal(bz, &score); err != nil {
		return PoolRiskScore{}, false
	}
	return score, true
}

func (k Keeper) SetPoolRiskScore(ctx context.Context, score PoolRiskScore) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(score)
	kvStore.Set(riskScoreKey(score.PoolID), bz)
}

func (k Keeper) GetRiskHistory(ctx context.Context, poolID uint64) RiskHistory {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(riskHistoryKey(poolID))
	if err != nil || bz == nil {
		return RiskHistory{PoolID: poolID, Entries: []RiskHistoryEntry{}}
	}
	var history RiskHistory
	if err := json.Unmarshal(bz, &history); err != nil {
		return RiskHistory{PoolID: poolID, Entries: []RiskHistoryEntry{}}
	}
	return history
}

func (k Keeper) SetRiskHistory(ctx context.Context, history RiskHistory) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(history)
	kvStore.Set(riskHistoryKey(history.PoolID), bz)
}

// ---------------------------------------------------------------------------
// Trade Risk Scoring — called before trade execution (read-only)
// ---------------------------------------------------------------------------

// CalculateTradeRisk computes a comprehensive risk score for a proposed trade.
func (k Keeper) CalculateTradeRisk(ctx context.Context, poolID uint64, inputDenom string, inputAmount math.Int) (TradeRiskScore, error) {
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return TradeRiskScore{}, types.ErrPoolNotFound
	}

	// Determine reserves
	var reserveIn, reserveOut math.Int
	var denomOut string
	if inputDenom == pool.DenomA {
		reserveIn = pool.ReserveA
		reserveOut = pool.ReserveB
		denomOut = pool.DenomB
	} else if inputDenom == pool.DenomB {
		reserveIn = pool.ReserveB
		reserveOut = pool.ReserveA
		denomOut = pool.DenomA
	} else {
		return TradeRiskScore{}, types.ErrInvalidDenom
	}

	if reserveIn.IsZero() || reserveOut.IsZero() {
		return TradeRiskScore{}, types.ErrInsufficientLiquidity
	}

	// Calculate expected output using constant-product formula
	feeBps := pool.SwapFee.MulInt64(10000).TruncateInt64()
	amountInAfterFee := inputAmount.MulRaw(10000 - feeBps).QuoRaw(10000)
	expectedOutput := reserveOut.Mul(amountInAfterFee).Quo(reserveIn.Add(amountInAfterFee))

	// 1. Price Impact Score
	// priceImpact = 1 - (outputAmount * reserveIn) / (inputAmount * reserveOut)
	priceImpactPct := math.LegacyZeroDec()
	if !inputAmount.IsZero() && reserveOut.IsPositive() {
		spotNum := math.LegacyNewDecFromInt(expectedOutput).Mul(math.LegacyNewDecFromInt(reserveIn))
		spotDen := math.LegacyNewDecFromInt(inputAmount).Mul(math.LegacyNewDecFromInt(reserveOut))
		if spotDen.IsPositive() {
			ratio := spotNum.Quo(spotDen)
			priceImpactPct = math.LegacyOneDec().Sub(ratio)
			if priceImpactPct.IsNegative() {
				priceImpactPct = math.LegacyZeroDec()
			}
		}
	}
	priceImpactScore := priceImpactToScore(priceImpactPct)

	// 2. Slippage Risk Score — order size relative to pool liquidity
	// ratio = inputAmount / reserveIn
	sizeRatio := math.LegacyNewDecFromInt(inputAmount).Quo(math.LegacyNewDecFromInt(reserveIn))
	slippageScore := sizeRatioToScore(sizeRatio)

	// 3. Pool Health Score (inverse — 100 = unhealthy pool for the trade)
	poolHealthScore := k.calculatePoolBalanceScore(pool)
	// Invert: high balance score = healthy = low trade risk
	poolHealthRisk := 100 - poolHealthScore

	// 4. Timing Score — based on recent volatility from risk state
	timingScore := k.calculateTimingScore(ctx, poolID)

	// Weighted overall score
	// Price impact: 40%, Slippage: 25%, Pool health: 20%, Timing: 15%
	overall := (priceImpactScore*40 + slippageScore*25 + poolHealthRisk*20 + timingScore*15) / 100
	if overall > 100 {
		overall = 100
	}
	if overall < 0 {
		overall = 0
	}

	riskLevel := scoreToRiskLevel(overall)

	// Build warnings
	warnings := buildTradeWarnings(priceImpactPct, sizeRatio, poolHealthRisk, timingScore)

	return TradeRiskScore{
		PoolID:           poolID,
		InputDenom:       inputDenom,
		InputAmount:      inputAmount,
		PriceImpactScore: priceImpactScore,
		SlippageScore:    slippageScore,
		PoolHealthScore:  poolHealthRisk,
		TimingScore:      timingScore,
		OverallScore:     overall,
		RiskLevel:        riskLevel,
		PriceImpactPct:   priceImpactPct,
		ExpectedOutput:   expectedOutput,
		OutputDenom:      denomOut,
		Warnings:         warnings,
	}, nil
}

// ---------------------------------------------------------------------------
// Pool Risk Scoring — called each block in BeginBlock
// ---------------------------------------------------------------------------

// CalculatePoolRiskScore computes the risk score for a pool and stores it.
func (k Keeper) CalculatePoolRiskScore(ctx context.Context, pool types.Pool) PoolRiskScore {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// 1. Liquidity Depth Score (100 = deep, 0 = shallow)
	liquidityScore := k.calculateLiquidityDepthScore(pool)

	// 2. Volatility Score (100 = very volatile, 0 = stable)
	volatilityScore := k.calculateVolatilityScore(ctx, pool.ID)

	// 3. Volume Concentration Score (100 = one whale dominating, 0 = distributed)
	volumeConcentration := k.calculateVolumeConcentrationScore(ctx, pool.ID)

	// 4. Reserve Balance Score (100 = perfectly balanced, 0 = completely skewed)
	reserveBalanceScore := k.calculatePoolBalanceScore(pool)

	// 5. Impermanent Loss Risk (100 = high IL risk, 0 = no IL risk)
	ilRisk := k.calculateImpermanentLossRisk(ctx, pool.ID)

	// Overall pool risk: weighted combination
	// Lower = healthier. We invert liquidity and balance scores since they are "positive" metrics.
	// Volatility: 30%, IL risk: 25%, Volume concentration: 20%, Liquidity (inverted): 15%, Balance (inverted): 10%
	overallRisk := (volatilityScore*30 + ilRisk*25 + volumeConcentration*20 + (100-liquidityScore)*15 + (100-reserveBalanceScore)*10) / 100
	if overallRisk > 100 {
		overallRisk = 100
	}
	if overallRisk < 0 {
		overallRisk = 0
	}

	poolGrade := scoreToPoolGrade(overallRisk)

	halted, haltReason := k.IsPoolHalted(ctx, pool.ID)

	// Price change since risk state initialization
	priceChange := math.LegacyZeroDec()
	state, found := k.GetRiskState(ctx, pool.ID)
	if found && !state.PrevPrice.IsZero() && !pool.ReserveA.IsZero() {
		currentPrice := math.LegacyNewDecFromInt(pool.ReserveB).Quo(math.LegacyNewDecFromInt(pool.ReserveA))
		priceChange = currentPrice.Sub(state.PrevPrice).Quo(state.PrevPrice)
	}

	score := PoolRiskScore{
		PoolID:                pool.ID,
		LiquidityDepthScore:   liquidityScore,
		VolatilityScore:       volatilityScore,
		VolumeConcentration:   volumeConcentration,
		ReserveBalanceScore:   reserveBalanceScore,
		ImpermanentLossRisk:   ilRisk,
		OverallScore:          overallRisk,
		PoolGrade:             poolGrade,
		BlockHeight:           sdkCtx.BlockHeight(),
		IsHalted:              halted,
		HaltReason:            haltReason,
		PriceChangeSinceStart: priceChange,
	}

	// Store the score
	k.SetPoolRiskScore(ctx, score)

	// Append to history
	k.appendRiskHistory(ctx, pool.ID, score)

	return score
}

// getRiskScoreCursor returns the next pool ID at which to resume risk score
// processing. Zero means "start from the beginning of the pool list".
func (k Keeper) getRiskScoreCursor(ctx context.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.RiskScoreCursorKey))
	if err != nil || bz == nil || len(bz) != 8 {
		return 0
	}
	return binaryBigEndianUint64(bz)
}

// setRiskScoreCursor stores the next pool ID to resume from.
func (k Keeper) setRiskScoreCursor(ctx context.Context, poolID uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	putUint64BigEndian(bz, poolID)
	kvStore.Set([]byte(types.RiskScoreCursorKey), bz)
}

// UpdateAllPoolRiskScores calculates and stores risk scores for up to
// MaxPoolsPerBlock pools per call, advancing a stored cursor each block so
// that every pool is eventually processed (round-robin). Called in BeginBlock.
//
// H6: previously this iterated every pool unconditionally, which is O(N) per
// block and an obvious DoS vector once the chain has many pools.
func (k Keeper) UpdateAllPoolRiskScores(ctx context.Context) {
	pools := k.GetAllPools(ctx)
	if len(pools) == 0 {
		return
	}

	// Sort by pool ID for stable round-robin ordering.
	sort.Slice(pools, func(i, j int) bool { return pools[i].ID < pools[j].ID })

	cursor := k.getRiskScoreCursor(ctx)

	// Find the starting index >= cursor. If the cursor falls past the end
	// (e.g., pools were deleted), wrap around to the beginning.
	startIdx := 0
	for i, p := range pools {
		if p.ID >= cursor {
			startIdx = i
			break
		}
		if i == len(pools)-1 {
			startIdx = 0
		}
	}

	processed := 0
	idx := startIdx
	for processed < MaxPoolsPerBlock && processed < len(pools) {
		k.CalculatePoolRiskScore(ctx, pools[idx])
		processed++
		idx = (idx + 1) % len(pools)
	}

	// Persist next cursor as the pool ID we would resume from.
	k.setRiskScoreCursor(ctx, pools[idx].ID)
}

// binaryBigEndianUint64 / putUint64BigEndian — small local helpers to avoid
// importing encoding/binary in this file (already imported elsewhere in the
// package, but we keep this self-contained for clarity).
func binaryBigEndianUint64(b []byte) uint64 {
	_ = b[7]
	return uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
}

func putUint64BigEndian(b []byte, v uint64) {
	_ = b[7]
	b[0] = byte(v >> 56)
	b[1] = byte(v >> 48)
	b[2] = byte(v >> 40)
	b[3] = byte(v >> 32)
	b[4] = byte(v >> 24)
	b[5] = byte(v >> 16)
	b[6] = byte(v >> 8)
	b[7] = byte(v)
}

// ---------------------------------------------------------------------------
// Component Score Calculations
// ---------------------------------------------------------------------------

// calculateLiquidityDepthScore returns 0-100 where 100 = deep liquidity.
// Uses geometric mean of reserves compared to thresholds.
func (k Keeper) calculateLiquidityDepthScore(pool types.Pool) int {
	if pool.ReserveA.IsZero() || pool.ReserveB.IsZero() {
		return 0
	}

	// Geometric mean of reserves (using sqrt(A*B) which is already the LP shares concept)
	// Compare against thresholds:
	// < 1M = 10, < 10M = 30, < 100M = 50, < 1B = 70, < 10B = 85, >= 10B = 100
	geoMean := isqrt(pool.ReserveA.Mul(pool.ReserveB))

	thresholds := []struct {
		limit math.Int
		score int
	}{
		{math.NewInt(1_000_000), 10},
		{math.NewInt(10_000_000), 30},
		{math.NewInt(100_000_000), 50},
		{math.NewInt(1_000_000_000), 70},
		{math.NewInt(10_000_000_000), 85},
	}

	for _, t := range thresholds {
		if geoMean.LT(t.limit) {
			return t.score
		}
	}
	return 100
}

// calculateVolatilityScore returns 0-100 where 100 = high volatility.
// Based on recent price change from the risk state.
func (k Keeper) calculateVolatilityScore(ctx context.Context, poolID uint64) int {
	state, found := k.GetRiskState(ctx, poolID)
	if !found {
		return 0
	}

	// Use absolute price change percentage
	absPriceChange := state.PriceChangePct.Abs()

	// Map price changes to volatility scores:
	// 0% = 0, 1% = 10, 5% = 30, 10% = 50, 20% = 75, 30%+ = 100
	pctInt := absPriceChange.MulInt64(100).TruncateInt64() // convert to integer percentage

	switch {
	case pctInt >= 30:
		return 100
	case pctInt >= 20:
		return 75
	case pctInt >= 10:
		return 50
	case pctInt >= 5:
		return 30
	case pctInt >= 1:
		return 10
	default:
		return 0
	}
}

// calculateVolumeConcentrationScore returns 0-100 where 100 = one whale dominating.
// Based on the ratio of current block volume to max volume per block.
func (k Keeper) calculateVolumeConcentrationScore(ctx context.Context, poolID uint64) int {
	state, found := k.GetRiskState(ctx, poolID)
	if !found {
		return 0
	}

	if state.VolumeThisBlock.IsZero() || maxVolumePerBlock.IsZero() {
		return 0
	}

	// ratio = volumeThisBlock / maxVolumePerBlock * 100
	ratio := math.LegacyNewDecFromInt(state.VolumeThisBlock).Quo(math.LegacyNewDecFromInt(maxVolumePerBlock)).MulInt64(100)
	score := int(ratio.TruncateInt64())

	if score > 100 {
		return 100
	}
	if score < 0 {
		return 0
	}
	return score
}

// calculatePoolBalanceScore returns 0-100 where 100 = perfectly balanced.
// Measures how close the reserves are to equal value.
func (k Keeper) calculatePoolBalanceScore(pool types.Pool) int {
	if pool.ReserveA.IsZero() || pool.ReserveB.IsZero() {
		return 0
	}

	// Ratio of smaller reserve to larger reserve
	// Perfect balance = 1.0 (score 100), extreme skew = 0.0 (score 0)
	var ratio math.LegacyDec
	if pool.ReserveA.GT(pool.ReserveB) {
		ratio = math.LegacyNewDecFromInt(pool.ReserveB).Quo(math.LegacyNewDecFromInt(pool.ReserveA))
	} else {
		ratio = math.LegacyNewDecFromInt(pool.ReserveA).Quo(math.LegacyNewDecFromInt(pool.ReserveB))
	}

	// ratio is between 0 and 1; map directly to 0-100
	score := int(ratio.MulInt64(100).TruncateInt64())
	if score > 100 {
		return 100
	}
	if score < 0 {
		return 0
	}
	return score
}

// calculateImpermanentLossRisk returns 0-100 where 100 = high IL risk.
// Based on price divergence from the initial price in the risk state.
func (k Keeper) calculateImpermanentLossRisk(ctx context.Context, poolID uint64) int {
	state, found := k.GetRiskState(ctx, poolID)
	if !found {
		return 0
	}

	pool, found := k.GetPool(ctx, poolID)
	if !found || pool.ReserveA.IsZero() || state.PrevPrice.IsZero() {
		return 0
	}

	currentPrice := math.LegacyNewDecFromInt(pool.ReserveB).Quo(math.LegacyNewDecFromInt(pool.ReserveA))

	// Price ratio = currentPrice / prevPrice
	priceRatio := currentPrice.Quo(state.PrevPrice)

	// IL formula (simplified): IL = 2 * sqrt(priceRatio) / (1 + priceRatio) - 1
	// For practical scoring, we use the absolute price divergence as a proxy:
	// divergence = |priceRatio - 1|
	// Map: 0% = 0, 5% = 10, 10% = 25, 20% = 45, 30% = 65, 50%+ = 85, 100%+ = 100
	divergence := priceRatio.Sub(math.LegacyOneDec()).Abs()
	divergencePct := divergence.MulInt64(100).TruncateInt64()

	switch {
	case divergencePct >= 100:
		return 100
	case divergencePct >= 50:
		return 85
	case divergencePct >= 30:
		return 65
	case divergencePct >= 20:
		return 45
	case divergencePct >= 10:
		return 25
	case divergencePct >= 5:
		return 10
	default:
		return 0
	}
}

// calculateTimingScore returns 0-100 where 100 = bad timing (volatile).
// Uses recent volatility and halt status.
func (k Keeper) calculateTimingScore(ctx context.Context, poolID uint64) int {
	// Check if halted
	halted, _ := k.IsPoolHalted(ctx, poolID)
	if halted {
		return 100
	}

	// Use volatility as timing indicator
	volatility := k.calculateVolatilityScore(ctx, poolID)

	// Also check volume concentration — high volume = more competitive
	volumeConc := k.calculateVolumeConcentrationScore(ctx, poolID)

	// Timing = 70% volatility + 30% volume concentration
	timing := (volatility*70 + volumeConc*30) / 100
	if timing > 100 {
		return 100
	}
	return timing
}

// ---------------------------------------------------------------------------
// Score Mapping Helpers
// ---------------------------------------------------------------------------

// priceImpactToScore maps price impact percentage to a 0-100 score.
func priceImpactToScore(impact math.LegacyDec) int {
	// impact is a decimal (e.g., 0.01 = 1%)
	impactPct := impact.MulInt64(100).TruncateInt64() // integer percentage

	switch {
	case impactPct >= 20:
		return 100
	case impactPct >= 10:
		return 80
	case impactPct >= 5:
		return 60
	case impactPct >= 2:
		return 40
	case impactPct >= 1:
		return 20
	default:
		// Sub-1% — use fractional mapping
		// impact * 10000 gives basis points
		bps := impact.MulInt64(10000).TruncateInt64()
		if bps >= 50 { // 0.5%
			return 15
		}
		if bps >= 10 { // 0.1%
			return 5
		}
		return 0
	}
}

// sizeRatioToScore maps order-size-to-reserve ratio to a 0-100 score.
func sizeRatioToScore(ratio math.LegacyDec) int {
	// ratio = inputAmount / reserveIn
	pct := ratio.MulInt64(100).TruncateInt64()

	switch {
	case pct >= 50:
		return 100
	case pct >= 25:
		return 80
	case pct >= 10:
		return 60
	case pct >= 5:
		return 40
	case pct >= 1:
		return 20
	default:
		return 5
	}
}

// scoreToRiskLevel converts a 0-100 risk score to a level string.
func scoreToRiskLevel(score int) string {
	switch {
	case score >= 75:
		return RiskLevelExtreme
	case score >= 50:
		return RiskLevelHigh
	case score >= 25:
		return RiskLevelMedium
	default:
		return RiskLevelLow
	}
}

// scoreToPoolGrade converts a 0-100 pool risk score to a letter grade.
// Lower risk = better grade.
func scoreToPoolGrade(score int) string {
	switch {
	case score <= 15:
		return PoolGradeA
	case score <= 30:
		return PoolGradeB
	case score <= 50:
		return PoolGradeC
	case score <= 70:
		return PoolGradeD
	default:
		return PoolGradeF
	}
}

// ---------------------------------------------------------------------------
// Trade Warnings
// ---------------------------------------------------------------------------

func buildTradeWarnings(priceImpact math.LegacyDec, sizeRatio math.LegacyDec, poolHealthRisk int, timingScore int) []string {
	var warnings []string

	// Price impact warnings
	impactPct := priceImpact.MulInt64(100)
	if impactPct.GTE(math.LegacyNewDec(5)) {
		warnings = append(warnings, fmt.Sprintf("High price impact: %s%%", impactPct.TruncateInt()))
	} else if impactPct.GTE(math.LegacyNewDec(2)) {
		warnings = append(warnings, fmt.Sprintf("Moderate price impact: %s%%", impactPct.TruncateInt()))
	}

	// Size warnings
	sizePct := sizeRatio.MulInt64(100)
	if sizePct.GTE(math.LegacyNewDec(10)) {
		warnings = append(warnings, fmt.Sprintf("Large trade relative to pool size: %s%% of reserves", sizePct.TruncateInt()))
	}

	// Pool health warnings
	if poolHealthRisk >= 60 {
		warnings = append(warnings, "Low liquidity warning")
	}

	// Timing warnings
	if timingScore >= 50 {
		warnings = append(warnings, "Pool is currently volatile")
	}

	return warnings
}

// ---------------------------------------------------------------------------
// History Management
// ---------------------------------------------------------------------------

func (k Keeper) appendRiskHistory(ctx context.Context, poolID uint64, score PoolRiskScore) {
	history := k.GetRiskHistory(ctx, poolID)

	entry := RiskHistoryEntry{
		BlockHeight:     score.BlockHeight,
		OverallScore:    score.OverallScore,
		PoolGrade:       score.PoolGrade,
		VolatilityScore: score.VolatilityScore,
		LiquidityScore:  score.LiquidityDepthScore,
		IsHalted:        score.IsHalted,
	}

	history.Entries = append(history.Entries, entry)

	// Trim to max entries — keep only the most recent
	if len(history.Entries) > maxHistoryEntries {
		history.Entries = history.Entries[len(history.Entries)-maxHistoryEntries:]
	}

	// Ensure entries are sorted by block height
	sort.Slice(history.Entries, func(i, j int) bool {
		if history.Entries[i].BlockHeight == history.Entries[j].BlockHeight {
			return history.Entries[i].OverallScore < history.Entries[j].OverallScore
		}
		return history.Entries[i].BlockHeight < history.Entries[j].BlockHeight
	})

	k.SetRiskHistory(ctx, history)
}
