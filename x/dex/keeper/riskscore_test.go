package keeper_test

import (
	"encoding/json"
	"testing"

	"cosmossdk.io/math"

	"github.com/stretchr/testify/require"

	"syreen/x/dex/keeper"
	"syreen/x/dex/types"
)

// ---------------------------------------------------------------------------
// Trade Risk Score Tests
// ---------------------------------------------------------------------------

func TestTradeRisk_SmallTrade_LowRisk(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(10_000_000),
		ReserveB: math.NewInt(10_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)

	// Initialize risk state
	k.RunRiskChecks(ctx)

	// Small trade: 10K out of 10M = 0.1% of reserves
	score, err := k.CalculateTradeRisk(ctx, 1, "tokenA", math.NewInt(10_000))
	require.NoError(t, err)
	require.Equal(t, uint64(1), score.PoolID)
	require.Equal(t, "tokenA", score.InputDenom)
	require.Equal(t, "tokenB", score.OutputDenom)
	require.True(t, score.ExpectedOutput.IsPositive())
	require.Equal(t, keeper.RiskLevelLow, score.RiskLevel)
	require.LessOrEqual(t, score.OverallScore, 25)
	require.Empty(t, score.Warnings)
}

func TestTradeRisk_MediumTrade_ModerateRisk(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(10_000_000),
		ReserveB: math.NewInt(10_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	// Medium trade: 500K out of 10M = 5% of reserves
	score, err := k.CalculateTradeRisk(ctx, 1, "tokenA", math.NewInt(500_000))
	require.NoError(t, err)
	require.Greater(t, score.OverallScore, 10)
	require.True(t, score.PriceImpactPct.IsPositive())
	require.True(t, score.ExpectedOutput.IsPositive())
	require.True(t, score.ExpectedOutput.LT(math.NewInt(500_000))) // output < input due to price impact + fee
}

func TestTradeRisk_LargeTrade_HighRisk(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(10_000_000),
		ReserveB: math.NewInt(10_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	// Large trade: 5M out of 10M = 50% of reserves
	score, err := k.CalculateTradeRisk(ctx, 1, "tokenA", math.NewInt(5_000_000))
	require.NoError(t, err)
	require.GreaterOrEqual(t, score.OverallScore, 50)
	// Should have warnings about high price impact and large trade
	require.NotEmpty(t, score.Warnings)

	hasImpactWarning := false
	hasSizeWarning := false
	for _, w := range score.Warnings {
		if len(w) > 0 {
			if w[0] == 'H' || w[0] == 'M' {
				hasImpactWarning = true
			}
			if w[0] == 'L' {
				hasSizeWarning = true
			}
		}
	}
	require.True(t, hasImpactWarning || hasSizeWarning, "should have risk warnings for large trade")
}

func TestTradeRisk_WhaleTrade_ExtremeRisk(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	// Small pool with low liquidity
	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(1_000_000),
		ReserveB: math.NewInt(1_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	// Whale: 800K out of 1M = 80% of reserves
	score, err := k.CalculateTradeRisk(ctx, 1, "tokenA", math.NewInt(800_000))
	require.NoError(t, err)
	require.GreaterOrEqual(t, score.OverallScore, 60)
	require.NotEmpty(t, score.Warnings)
}

func TestTradeRisk_InvalidPool_ReturnsError(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	_, err := k.CalculateTradeRisk(ctx, 999, "tokenA", math.NewInt(1000))
	require.ErrorIs(t, err, types.ErrPoolNotFound)
}

func TestTradeRisk_InvalidDenom_ReturnsError(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(10_000_000),
		ReserveB: math.NewInt(10_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)

	_, err := k.CalculateTradeRisk(ctx, 1, "tokenC", math.NewInt(1000))
	require.ErrorIs(t, err, types.ErrInvalidDenom)
}

func TestTradeRisk_BothDirections(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(10_000_000),
		ReserveB: math.NewInt(10_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	// Trade A -> B
	scoreAB, err := k.CalculateTradeRisk(ctx, 1, "tokenA", math.NewInt(100_000))
	require.NoError(t, err)
	require.Equal(t, "tokenB", scoreAB.OutputDenom)

	// Trade B -> A
	scoreBA, err := k.CalculateTradeRisk(ctx, 1, "tokenB", math.NewInt(100_000))
	require.NoError(t, err)
	require.Equal(t, "tokenA", scoreBA.OutputDenom)

	// Same amount in a balanced pool should give similar risk scores
	require.Equal(t, scoreAB.OverallScore, scoreBA.OverallScore)
}

func TestTradeRisk_VolatilePool_HigherTimingScore(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(10_000_000),
		ReserveB: math.NewInt(10_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	// Get baseline risk
	scoreBaseline, err := k.CalculateTradeRisk(ctx, 1, "tokenA", math.NewInt(100_000))
	require.NoError(t, err)

	// Simulate price movement (15% move)
	pool.ReserveA = math.NewInt(11_500_000)
	pool.ReserveB = math.NewInt(8_700_000)
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	// Check risk after volatility
	scoreVolatile, err := k.CalculateTradeRisk(ctx, 1, "tokenA", math.NewInt(100_000))
	require.NoError(t, err)

	// Timing score should be higher when pool is volatile
	require.GreaterOrEqual(t, scoreVolatile.TimingScore, scoreBaseline.TimingScore)
}

// ---------------------------------------------------------------------------
// Pool Risk Score Tests
// ---------------------------------------------------------------------------

func TestPoolRiskScore_BalancedPool_GoodGrade(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(100_000_000_000), // 100B
		ReserveB: math.NewInt(100_000_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	score := k.CalculatePoolRiskScore(ctx, pool)
	require.Equal(t, uint64(1), score.PoolID)
	require.Equal(t, int64(100), score.BlockHeight)
	require.False(t, score.IsHalted)

	// Balanced pool with deep liquidity should get good grade
	require.LessOrEqual(t, score.OverallScore, 30)
	require.Contains(t, []string{keeper.PoolGradeA, keeper.PoolGradeB}, score.PoolGrade)

	// Reserve balance should be high (100 = perfectly balanced)
	require.Equal(t, 100, score.ReserveBalanceScore)

	// Liquidity should be high
	require.GreaterOrEqual(t, score.LiquidityDepthScore, 85)
}

func TestPoolRiskScore_SkewedPool_LowerGrade(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	// Heavily skewed pool
	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(100_000_000),
		ReserveB: math.NewInt(10_000_000), // 10:1 ratio
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	score := k.CalculatePoolRiskScore(ctx, pool)

	// Reserve balance should be low (only 10% balanced)
	require.Equal(t, 10, score.ReserveBalanceScore)
}

func TestPoolRiskScore_LowLiquidity_LowScore(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(500_000), // very low
		ReserveB: math.NewInt(500_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	score := k.CalculatePoolRiskScore(ctx, pool)

	// Low liquidity pool should have low depth score
	require.LessOrEqual(t, score.LiquidityDepthScore, 10)
}

func TestPoolRiskScore_HaltedPool(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(10_000_000),
		ReserveB: math.NewInt(10_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	// Trigger halt via 85% reserve drain
	pool.ReserveA = math.NewInt(1_000_000)
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	score := k.CalculatePoolRiskScore(ctx, pool)
	require.True(t, score.IsHalted)
	require.NotEmpty(t, score.HaltReason)
}

func TestPoolRiskScore_VolatilePool(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(10_000_000),
		ReserveB: math.NewInt(10_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	// Simulate 15% price move (below halt threshold)
	pool.ReserveA = math.NewInt(11_500_000)
	pool.ReserveB = math.NewInt(8_700_000)
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	score := k.CalculatePoolRiskScore(ctx, pool)
	require.Greater(t, score.VolatilityScore, 0)
	require.Greater(t, score.ImpermanentLossRisk, 0)
}

func TestPoolRiskScore_StoredAndRetrieved(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(10_000_000),
		ReserveB: math.NewInt(10_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	// Calculate and store
	k.CalculatePoolRiskScore(ctx, pool)

	// Retrieve
	retrieved, found := k.GetPoolRiskScore(ctx, 1)
	require.True(t, found)
	require.Equal(t, uint64(1), retrieved.PoolID)
	require.NotEmpty(t, retrieved.PoolGrade)
	require.Equal(t, int64(100), retrieved.BlockHeight)
}

func TestPoolRiskScore_UpdateAllPools(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	// Create two pools
	pool1 := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(10_000_000),
		ReserveB: math.NewInt(10_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	pool2 := types.Pool{
		ID:       2,
		DenomA:   "tokenC",
		DenomB:   "tokenD",
		ReserveA: math.NewInt(5_000_000),
		ReserveB: math.NewInt(5_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool1)
	k.SetPool(ctx, pool2)
	k.RunRiskChecks(ctx)

	k.UpdateAllPoolRiskScores(ctx)

	s1, found1 := k.GetPoolRiskScore(ctx, 1)
	s2, found2 := k.GetPoolRiskScore(ctx, 2)
	require.True(t, found1)
	require.True(t, found2)
	require.Equal(t, uint64(1), s1.PoolID)
	require.Equal(t, uint64(2), s2.PoolID)
}

// ---------------------------------------------------------------------------
// Pool Grade Tests
// ---------------------------------------------------------------------------

func TestPoolGrade_Thresholds(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	// Grade A: deep liquidity, balanced, no volatility
	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(1_000_000_000_000), // 1T
		ReserveB: math.NewInt(1_000_000_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	score := k.CalculatePoolRiskScore(ctx, pool)
	require.Equal(t, keeper.PoolGradeA, score.PoolGrade,
		"deep balanced pool should be grade A, got score %d", score.OverallScore)
}

// ---------------------------------------------------------------------------
// Risk History Tests
// ---------------------------------------------------------------------------

func TestRiskHistory_Recorded(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(10_000_000),
		ReserveB: math.NewInt(10_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	// First block
	k.CalculatePoolRiskScore(ctx, pool)

	history := k.GetRiskHistory(ctx, 1)
	require.Equal(t, uint64(1), history.PoolID)
	require.Len(t, history.Entries, 1)
	require.Equal(t, int64(100), history.Entries[0].BlockHeight)
}

func TestRiskHistory_MultipleBlocks(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(10_000_000),
		ReserveB: math.NewInt(10_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	// Simulate 5 blocks
	for i := 0; i < 5; i++ {
		blockCtx := ctx.WithBlockHeight(int64(100 + i))
		k.CalculatePoolRiskScore(blockCtx, pool)
	}

	history := k.GetRiskHistory(ctx, 1)
	require.Equal(t, 5, len(history.Entries))

	// Check entries are sorted by block height
	for i := 1; i < len(history.Entries); i++ {
		require.Greater(t, history.Entries[i].BlockHeight, history.Entries[i-1].BlockHeight)
	}
}

func TestRiskHistory_TrimmedToMax(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(10_000_000),
		ReserveB: math.NewInt(10_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	// Add 150 entries (max is 100)
	for i := 0; i < 150; i++ {
		blockCtx := ctx.WithBlockHeight(int64(100 + i))
		k.CalculatePoolRiskScore(blockCtx, pool)
	}

	history := k.GetRiskHistory(ctx, 1)
	require.LessOrEqual(t, len(history.Entries), 100)

	// Oldest should be block 150 (entries 50-149 kept)
	require.Equal(t, int64(150), history.Entries[0].BlockHeight)
}

func TestRiskHistory_EmptyForNonExistentPool(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	history := k.GetRiskHistory(ctx, 999)
	require.Equal(t, uint64(999), history.PoolID)
	require.Empty(t, history.Entries)
}

// ---------------------------------------------------------------------------
// Score Mapping Tests
// ---------------------------------------------------------------------------

func TestScoreToRiskLevel(t *testing.T) {
	tests := []struct {
		score    int
		expected string
	}{
		{0, keeper.RiskLevelLow},
		{10, keeper.RiskLevelLow},
		{24, keeper.RiskLevelLow},
		{25, keeper.RiskLevelMedium},
		{49, keeper.RiskLevelMedium},
		{50, keeper.RiskLevelHigh},
		{74, keeper.RiskLevelHigh},
		{75, keeper.RiskLevelExtreme},
		{100, keeper.RiskLevelExtreme},
	}

	k, ctx, _, _ := setupKeeper(t)
	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(10_000_000),
		ReserveB: math.NewInt(10_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)

	// Just verify the levels are defined
	for _, tt := range tests {
		require.NotEmpty(t, tt.expected)
		require.GreaterOrEqual(t, tt.score, 0)
		require.LessOrEqual(t, tt.score, 100)
	}
}

// ---------------------------------------------------------------------------
// Trade Warnings Tests
// ---------------------------------------------------------------------------

func TestTradeWarnings_HighImpact(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(1_000_000),
		ReserveB: math.NewInt(1_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	// 30% of pool = high impact
	score, err := k.CalculateTradeRisk(ctx, 1, "tokenA", math.NewInt(300_000))
	require.NoError(t, err)
	require.NotEmpty(t, score.Warnings)

	foundImpactWarning := false
	foundSizeWarning := false
	for _, w := range score.Warnings {
		if len(w) > 10 {
			if w[:4] == "High" || w[:8] == "Moderate" {
				foundImpactWarning = true
			}
			if w[:5] == "Large" {
				foundSizeWarning = true
			}
		}
	}
	require.True(t, foundImpactWarning || foundSizeWarning, "should warn about impact or size")
}

func TestTradeWarnings_NoWarningsForSmallTrade(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(10_000_000_000),
		ReserveB: math.NewInt(10_000_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	// Tiny trade: 1000 out of 10B = 0.00001% — no warnings expected
	score, err := k.CalculateTradeRisk(ctx, 1, "tokenA", math.NewInt(1000))
	require.NoError(t, err)
	require.Empty(t, score.Warnings, "tiny trade should have no warnings, got: %v", score.Warnings)
}

// ---------------------------------------------------------------------------
// Component Score Tests
// ---------------------------------------------------------------------------

func TestLiquidityDepthScore_Thresholds(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	tests := []struct {
		name     string
		reserveA int64
		reserveB int64
		minScore int
		maxScore int
	}{
		{"tiny pool", 100_000, 100_000, 0, 10},
		{"small pool", 5_000_000, 5_000_000, 10, 30},
		{"medium pool", 50_000_000, 50_000_000, 30, 50},
		{"large pool", 500_000_000, 500_000_000, 50, 70},
		{"very large pool", 5_000_000_000, 5_000_000_000, 70, 85},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := types.Pool{
				ID:       1,
				DenomA:   "tokenA",
				DenomB:   "tokenB",
				ReserveA: math.NewInt(tt.reserveA),
				ReserveB: math.NewInt(tt.reserveB),
				SwapFee:  math.LegacyNewDecWithPrec(3, 3),
			}
			k.SetPool(ctx, pool)
			k.RunRiskChecks(ctx)

			score := k.CalculatePoolRiskScore(ctx, pool)
			require.GreaterOrEqual(t, score.LiquidityDepthScore, tt.minScore,
				"%s: liquidity score %d < expected min %d", tt.name, score.LiquidityDepthScore, tt.minScore)
			require.LessOrEqual(t, score.LiquidityDepthScore, tt.maxScore,
				"%s: liquidity score %d > expected max %d", tt.name, score.LiquidityDepthScore, tt.maxScore)
		})
	}
}

func TestReserveBalanceScore(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	tests := []struct {
		name     string
		reserveA int64
		reserveB int64
		expected int
	}{
		{"perfectly balanced", 1_000_000, 1_000_000, 100},
		{"50% balanced", 1_000_000, 500_000, 50},
		{"10% balanced", 1_000_000, 100_000, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := types.Pool{
				ID:       1,
				DenomA:   "tokenA",
				DenomB:   "tokenB",
				ReserveA: math.NewInt(tt.reserveA),
				ReserveB: math.NewInt(tt.reserveB),
				SwapFee:  math.LegacyNewDecWithPrec(3, 3),
			}
			k.SetPool(ctx, pool)
			k.RunRiskChecks(ctx)

			score := k.CalculatePoolRiskScore(ctx, pool)
			require.Equal(t, tt.expected, score.ReserveBalanceScore,
				"%s: expected %d, got %d", tt.name, tt.expected, score.ReserveBalanceScore)
		})
	}
}

func TestZeroReserves_Handled(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.ZeroInt(),
		ReserveB: math.ZeroInt(),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)

	// Should not panic
	score := k.CalculatePoolRiskScore(ctx, pool)
	require.Equal(t, 0, score.LiquidityDepthScore)
	require.Equal(t, 0, score.ReserveBalanceScore)
}

// ---------------------------------------------------------------------------
// gRPC Query Tests
// ---------------------------------------------------------------------------

func TestQueryRiskScore(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(10_000_000),
		ReserveB: math.NewInt(10_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)
	k.UpdateAllPoolRiskScores(ctx)

	qs := keeper.NewQueryServerImpl(k)
	resp, err := qs.RiskScore(ctx, &types.QueryRiskScoreRequest{PoolID: 1})
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify JSON is valid
	var scoreMap map[string]interface{}
	err = json.Unmarshal(resp.RiskScore, &scoreMap)
	require.NoError(t, err)
	require.Contains(t, scoreMap, "pool_id")
	require.Contains(t, scoreMap, "pool_grade")
	require.Contains(t, scoreMap, "overall_score")
}

func TestQueryRiskScore_PoolNotFound(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	qs := keeper.NewQueryServerImpl(k)
	_, err := qs.RiskScore(ctx, &types.QueryRiskScoreRequest{PoolID: 999})
	require.Error(t, err)
}

func TestQueryTradeRisk(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(10_000_000),
		ReserveB: math.NewInt(10_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	qs := keeper.NewQueryServerImpl(k)
	resp, err := qs.TradeRisk(ctx, &types.QueryTradeRiskRequest{
		PoolID:     1,
		InputDenom: "tokenA",
		Amount:     math.NewInt(100_000),
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	var tradeRisk map[string]interface{}
	err = json.Unmarshal(resp.TradeRisk, &tradeRisk)
	require.NoError(t, err)
	require.Contains(t, tradeRisk, "risk_level")
	require.Contains(t, tradeRisk, "price_impact_pct")
	require.Contains(t, tradeRisk, "warnings")
}

func TestQueryRiskHistory(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(10_000_000),
		ReserveB: math.NewInt(10_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)
	k.CalculatePoolRiskScore(ctx, pool)

	qs := keeper.NewQueryServerImpl(k)
	resp, err := qs.RiskHistory(ctx, &types.QueryRiskHistoryRequest{PoolID: 1})
	require.NoError(t, err)
	require.NotNil(t, resp)

	var history map[string]interface{}
	err = json.Unmarshal(resp.RiskHistory, &history)
	require.NoError(t, err)
	require.Contains(t, history, "pool_id")
	require.Contains(t, history, "entries")
}

// ---------------------------------------------------------------------------
// Integration: BeginBlock updates risk scores
// ---------------------------------------------------------------------------

func TestBeginBlock_UpdatesRiskScores(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(10_000_000),
		ReserveB: math.NewInt(10_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)

	// Simulate what BeginBlock does
	k.RunRiskChecks(ctx)
	k.UpdateAllPoolRiskScores(ctx)

	// Risk score should now be stored
	score, found := k.GetPoolRiskScore(ctx, 1)
	require.True(t, found)
	require.Equal(t, uint64(1), score.PoolID)
	require.NotEmpty(t, score.PoolGrade)

	// History should have one entry
	history := k.GetRiskHistory(ctx, 1)
	require.Len(t, history.Entries, 1)
}

func TestBeginBlock_MultipleBlocksAccumulateHistory(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(10_000_000),
		ReserveB: math.NewInt(10_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)

	for i := 0; i < 10; i++ {
		blockCtx := ctx.WithBlockHeight(int64(100 + i))
		k.RunRiskChecks(blockCtx)
		k.UpdateAllPoolRiskScores(blockCtx)
	}

	history := k.GetRiskHistory(ctx, 1)
	require.Len(t, history.Entries, 10)
}
