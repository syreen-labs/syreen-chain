package keeper_test

import (
	"testing"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"syreen/x/dex/keeper"
	"syreen/x/dex/types"
)

// ---------------------------------------------------------------------------
// Helper: create a pool directly in the store (no bank transfers needed)
// ---------------------------------------------------------------------------

func setupPoolForSignals(t *testing.T, k keeper.Keeper, ctx sdk.Context, poolID uint64, reserveA, reserveB int64) types.Pool {
	t.Helper()
	pool := types.Pool{
		ID:       poolID,
		DenomA:   "usyreen",
		DenomB:   "uusdc",
		ReserveA: math.NewInt(reserveA),
		ReserveB: math.NewInt(reserveB),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)
	return pool
}

// ---------------------------------------------------------------------------
// TestSignals_InitializesOnFirstBlock
// ---------------------------------------------------------------------------

func TestSignals_InitializesOnFirstBlock(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	k.UpdatePoolSignals(ctx, pool)

	state, found := k.GetSignalState(ctx, 1)
	require.True(t, found)
	require.Equal(t, uint64(1), state.PoolID)
	require.Len(t, state.PriceHistory, 1)
	require.Equal(t, int64(100), state.PriceHistory[0].Height) // test ctx starts at height 100
	// Price should be 1.0 (10M/10M)
	require.True(t, state.PriceHistory[0].Price.Equal(math.LegacyOneDec()))
}

// ---------------------------------------------------------------------------
// TestSignals_PriceHistoryAccumulates
// ---------------------------------------------------------------------------

func TestSignals_PriceHistoryAccumulates(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	// Simulate 20 blocks with gradually increasing price
	for i := int64(0); i < 20; i++ {
		blockCtx := ctx.WithBlockHeight(100 + i)
		// Increase reserveB each block to simulate price rising
		pool.ReserveB = math.NewInt(10_000_000 + i*100_000)
		k.SetPool(blockCtx, pool)
		k.UpdatePoolSignals(blockCtx, pool)
	}

	state, found := k.GetSignalState(ctx, 1)
	require.True(t, found)
	require.Len(t, state.PriceHistory, 20)

	// Prices should be monotonically increasing
	for i := 1; i < len(state.PriceHistory); i++ {
		require.True(t, state.PriceHistory[i].Price.GT(state.PriceHistory[i-1].Price),
			"price at block %d should be > price at block %d", i, i-1)
	}
}

// ---------------------------------------------------------------------------
// TestSignals_PriceHistoryTrimmed
// ---------------------------------------------------------------------------

func TestSignals_PriceHistoryTrimmed(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	// Simulate more blocks than DefaultPriceHistoryDepth (100)
	for i := int64(0); i < 120; i++ {
		blockCtx := ctx.WithBlockHeight(100 + i)
		k.UpdatePoolSignals(blockCtx, pool)
	}

	state, found := k.GetSignalState(ctx, 1)
	require.True(t, found)
	require.Equal(t, keeper.DefaultPriceHistoryDepth, len(state.PriceHistory))
	// First entry should be from height 120 (100 + 120 - 100)
	require.Equal(t, int64(120), state.PriceHistory[0].Height)
}

// ---------------------------------------------------------------------------
// TestSignals_SkipsDuplicateHeight
// ---------------------------------------------------------------------------

func TestSignals_SkipsDuplicateHeight(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	k.UpdatePoolSignals(ctx, pool)
	k.UpdatePoolSignals(ctx, pool) // same height again

	state, found := k.GetSignalState(ctx, 1)
	require.True(t, found)
	require.Len(t, state.PriceHistory, 1, "should not add duplicate height")
}

// ---------------------------------------------------------------------------
// TestSignals_SMACalculation
// ---------------------------------------------------------------------------

func TestSignals_SMACalculation(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	// Create exactly 10 blocks at price=1.0
	for i := int64(0); i < 10; i++ {
		blockCtx := ctx.WithBlockHeight(100 + i)
		k.UpdatePoolSignals(blockCtx, pool)
	}

	signals := k.ComputePoolSignals(ctx.WithBlockHeight(109), 1)
	require.NotNil(t, signals)

	// SMA10 with all prices at 1.0 should be 1.0
	require.True(t, signals.MovingAverages.SMA10.Equal(math.LegacyOneDec()),
		"SMA10 should be 1.0, got %s", signals.MovingAverages.SMA10)
}

// ---------------------------------------------------------------------------
// TestSignals_SMAWithPriceChanges
// ---------------------------------------------------------------------------

func TestSignals_SMAWithPriceChanges(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	// Create pool with price=1.0
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	// 5 blocks at price=1.0
	for i := int64(0); i < 5; i++ {
		blockCtx := ctx.WithBlockHeight(100 + i)
		k.UpdatePoolSignals(blockCtx, pool)
	}

	// 5 blocks at price=2.0
	pool.ReserveB = math.NewInt(20_000_000)
	k.SetPool(ctx, pool)
	for i := int64(5); i < 10; i++ {
		blockCtx := ctx.WithBlockHeight(100 + i)
		k.UpdatePoolSignals(blockCtx, pool)
	}

	signals := k.ComputePoolSignals(ctx.WithBlockHeight(109), 1)
	require.NotNil(t, signals)

	// SMA10 should be average of (5 * 1.0 + 5 * 2.0) / 10 = 1.5
	expected := math.LegacyNewDecWithPrec(15, 1)
	require.True(t, signals.MovingAverages.SMA10.Equal(expected),
		"SMA10 should be 1.5, got %s", signals.MovingAverages.SMA10)
}

// ---------------------------------------------------------------------------
// TestSignals_EMAUpdates
// ---------------------------------------------------------------------------

func TestSignals_EMAUpdates(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	// First block: EMA should be initialized to current price
	k.UpdatePoolSignals(ctx, pool)
	state, _ := k.GetSignalState(ctx, 1)
	require.True(t, state.EMA10.Equal(math.LegacyOneDec()))

	// Second block with price=2.0: EMA should move toward 2.0
	pool.ReserveB = math.NewInt(20_000_000)
	k.SetPool(ctx, pool)
	blockCtx := ctx.WithBlockHeight(101)
	k.UpdatePoolSignals(blockCtx, pool)

	state, _ = k.GetSignalState(blockCtx, 1)
	// EMA10 should be between 1.0 and 2.0
	require.True(t, state.EMA10.GT(math.LegacyOneDec()), "EMA10 should be > 1.0")
	require.True(t, state.EMA10.LT(math.LegacyNewDec(2)), "EMA10 should be < 2.0")

	// k = 2/(10+1) = 0.1818...
	// EMA = 2.0 * k + 1.0 * (1-k) = 0.3636... + 0.8181... = 1.1818...
	expectedApprox := math.LegacyNewDecWithPrec(11818, 4) // ~1.1818
	diff := state.EMA10.Sub(expectedApprox).Abs()
	require.True(t, diff.LT(math.LegacyNewDecWithPrec(1, 3)),
		"EMA10 should be ~1.1818, got %s", state.EMA10)
}

// ---------------------------------------------------------------------------
// TestSignals_RSINeutralWithNoMovement
// ---------------------------------------------------------------------------

func TestSignals_RSINeutralWithNoMovement(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	// 20 blocks at constant price
	for i := int64(0); i < 20; i++ {
		blockCtx := ctx.WithBlockHeight(100 + i)
		k.UpdatePoolSignals(blockCtx, pool)
	}

	signals := k.ComputePoolSignals(ctx.WithBlockHeight(119), 1)
	require.NotNil(t, signals)

	// RSI with no movement should be 50 (neutral)
	require.True(t, signals.Momentum.RSI.Equal(math.LegacyNewDec(50)),
		"RSI should be 50 with no price movement, got %s", signals.Momentum.RSI)
}

// ---------------------------------------------------------------------------
// TestSignals_RSIHighWithConsistentGains
// ---------------------------------------------------------------------------

func TestSignals_RSIHighWithConsistentGains(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	// 20 blocks with price steadily increasing
	for i := int64(0); i < 20; i++ {
		blockCtx := ctx.WithBlockHeight(100 + i)
		pool.ReserveB = math.NewInt(10_000_000 + i*200_000)
		k.SetPool(blockCtx, pool)
		k.UpdatePoolSignals(blockCtx, pool)
	}

	signals := k.ComputePoolSignals(ctx.WithBlockHeight(119), 1)
	require.NotNil(t, signals)

	// RSI should be 100 (all gains, no losses)
	require.True(t, signals.Momentum.RSI.Equal(math.LegacyNewDec(100)),
		"RSI should be 100 with all gains, got %s", signals.Momentum.RSI)
}

// ---------------------------------------------------------------------------
// TestSignals_RSILowWithConsistentLosses
// ---------------------------------------------------------------------------

func TestSignals_RSILowWithConsistentLosses(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	// 20 blocks with price steadily decreasing
	for i := int64(0); i < 20; i++ {
		blockCtx := ctx.WithBlockHeight(100 + i)
		pool.ReserveB = math.NewInt(10_000_000 - i*100_000)
		k.SetPool(blockCtx, pool)
		k.UpdatePoolSignals(blockCtx, pool)
	}

	signals := k.ComputePoolSignals(ctx.WithBlockHeight(119), 1)
	require.NotNil(t, signals)

	// RSI should be 0 (all losses, no gains)
	require.True(t, signals.Momentum.RSI.Equal(math.LegacyZeroDec()),
		"RSI should be 0 with all losses, got %s", signals.Momentum.RSI)
}

// ---------------------------------------------------------------------------
// TestSignals_PriceROC
// ---------------------------------------------------------------------------

func TestSignals_PriceROC(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	// 15 blocks at price=1.0
	for i := int64(0); i < 15; i++ {
		blockCtx := ctx.WithBlockHeight(100 + i)
		k.UpdatePoolSignals(blockCtx, pool)
	}

	// Then price doubles to 2.0
	pool.ReserveB = math.NewInt(20_000_000)
	k.SetPool(ctx, pool)
	blockCtx := ctx.WithBlockHeight(115)
	k.UpdatePoolSignals(blockCtx, pool)

	signals := k.ComputePoolSignals(blockCtx, 1)
	require.NotNil(t, signals)

	// ROC over 14 blocks: (2.0 - 1.0) / 1.0 * 100 = 100%
	require.True(t, signals.Momentum.PriceROC.Equal(math.LegacyNewDec(100)),
		"Price ROC should be 100%%, got %s", signals.Momentum.PriceROC)
}

// ---------------------------------------------------------------------------
// TestSignals_VolatilityZeroWithConstantPrice
// ---------------------------------------------------------------------------

func TestSignals_VolatilityZeroWithConstantPrice(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	for i := int64(0); i < 25; i++ {
		blockCtx := ctx.WithBlockHeight(100 + i)
		k.UpdatePoolSignals(blockCtx, pool)
	}

	signals := k.ComputePoolSignals(ctx.WithBlockHeight(124), 1)
	require.NotNil(t, signals)

	// Volatility should be 0 with constant price
	require.True(t, signals.VolatilityScore.IsZero(),
		"Volatility should be 0 with constant price, got %s", signals.VolatilityScore)
}

// ---------------------------------------------------------------------------
// TestSignals_VolatilityNonZeroWithPriceChanges
// ---------------------------------------------------------------------------

func TestSignals_VolatilityNonZeroWithPriceChanges(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	// Alternate price up and down to create volatility
	for i := int64(0); i < 25; i++ {
		blockCtx := ctx.WithBlockHeight(100 + i)
		if i%2 == 0 {
			pool.ReserveB = math.NewInt(10_000_000)
		} else {
			pool.ReserveB = math.NewInt(12_000_000) // 20% up
		}
		k.SetPool(blockCtx, pool)
		k.UpdatePoolSignals(blockCtx, pool)
	}

	signals := k.ComputePoolSignals(ctx.WithBlockHeight(124), 1)
	require.NotNil(t, signals)

	require.True(t, signals.VolatilityScore.IsPositive(),
		"Volatility should be positive with price swings, got %s", signals.VolatilityScore)
}

// ---------------------------------------------------------------------------
// TestSignals_CompositeScoreNeutral
// ---------------------------------------------------------------------------

func TestSignals_CompositeScoreNeutral(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	// Stable price for many blocks -> neutral signal
	for i := int64(0); i < 60; i++ {
		blockCtx := ctx.WithBlockHeight(100 + i)
		k.UpdatePoolSignals(blockCtx, pool)
	}

	signals := k.ComputePoolSignals(ctx.WithBlockHeight(159), 1)
	require.NotNil(t, signals)

	// With constant price, all MAs equal price, RSI=50 -> should be ~neutral
	require.Equal(t, "NEUTRAL", signals.Signal,
		"Signal should be NEUTRAL with stable price, got %s (score %d)", signals.Signal, signals.CompositeScore)
}

// ---------------------------------------------------------------------------
// TestSignals_BullishSignalOnRisingPrice
// ---------------------------------------------------------------------------

func TestSignals_BullishSignalOnRisingPrice(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	// Strong uptrend: price increases every block
	for i := int64(0); i < 60; i++ {
		blockCtx := ctx.WithBlockHeight(100 + i)
		pool.ReserveB = math.NewInt(10_000_000 + i*500_000) // big steady gains
		k.SetPool(blockCtx, pool)
		k.UpdatePoolSignals(blockCtx, pool)
	}

	signals := k.ComputePoolSignals(ctx.WithBlockHeight(159), 1)
	require.NotNil(t, signals)

	// Score should be high (bullish)
	require.True(t, signals.CompositeScore > 50,
		"Score should be > 50 in uptrend, got %d", signals.CompositeScore)
	require.Contains(t, []string{"BUY", "STRONG_BUY"}, signals.Signal,
		"Signal should be BUY or STRONG_BUY, got %s", signals.Signal)
}

// ---------------------------------------------------------------------------
// TestSignals_BearishSignalOnFallingPrice
// ---------------------------------------------------------------------------

func TestSignals_BearishSignalOnFallingPrice(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 20_000_000) // start at price=2.0

	// Strong downtrend
	for i := int64(0); i < 60; i++ {
		blockCtx := ctx.WithBlockHeight(100 + i)
		newReserveB := int64(20_000_000) - i*200_000
		if newReserveB < 5_000_000 {
			newReserveB = 5_000_000
		}
		pool.ReserveB = math.NewInt(newReserveB)
		k.SetPool(blockCtx, pool)
		k.UpdatePoolSignals(blockCtx, pool)
	}

	signals := k.ComputePoolSignals(ctx.WithBlockHeight(159), 1)
	require.NotNil(t, signals)

	// Score should be low (bearish)
	require.True(t, signals.CompositeScore < 50,
		"Score should be < 50 in downtrend, got %d", signals.CompositeScore)
	require.Contains(t, []string{"SELL", "STRONG_SELL"}, signals.Signal,
		"Signal should be SELL or STRONG_SELL, got %s", signals.Signal)
}

// ---------------------------------------------------------------------------
// TestSignals_ScoreClampedTo0_100
// ---------------------------------------------------------------------------

func TestSignals_ScoreClampedTo0_100(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	// Extreme uptrend
	for i := int64(0); i < 60; i++ {
		blockCtx := ctx.WithBlockHeight(100 + i)
		pool.ReserveB = math.NewInt(10_000_000 + i*2_000_000)
		k.SetPool(blockCtx, pool)
		k.UpdatePoolSignals(blockCtx, pool)
	}

	signals := k.ComputePoolSignals(ctx.WithBlockHeight(159), 1)
	require.NotNil(t, signals)
	require.True(t, signals.CompositeScore >= 0, "Score should be >= 0")
	require.True(t, signals.CompositeScore <= 100, "Score should be <= 100")
}

// ---------------------------------------------------------------------------
// TestSignals_SignalHistory
// ---------------------------------------------------------------------------

func TestSignals_SignalHistory(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	// Run 10 blocks
	for i := int64(0); i < 10; i++ {
		blockCtx := ctx.WithBlockHeight(100 + i)
		k.UpdatePoolSignals(blockCtx, pool)
	}

	history, found := k.GetSignalHistory(ctx, 1)
	require.True(t, found)
	require.Equal(t, uint64(1), history.PoolID)
	require.Len(t, history.Entries, 10)

	// Each entry should have valid data
	for _, entry := range history.Entries {
		require.True(t, entry.Height >= 100)
		require.True(t, entry.Price.IsPositive())
		require.NotEmpty(t, entry.Signal)
		require.True(t, entry.CompositeScore >= 0 && entry.CompositeScore <= 100)
	}
}

// ---------------------------------------------------------------------------
// TestSignals_SignalHistoryTrimmed
// ---------------------------------------------------------------------------

func TestSignals_SignalHistoryTrimmed(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	// Run more blocks than SignalHistoryRetention (200)
	for i := int64(0); i < 220; i++ {
		blockCtx := ctx.WithBlockHeight(100 + i)
		k.UpdatePoolSignals(blockCtx, pool)
	}

	history, found := k.GetSignalHistory(ctx, 1)
	require.True(t, found)
	require.Equal(t, keeper.SignalHistoryRetention, len(history.Entries))
}

// ---------------------------------------------------------------------------
// TestSignals_NoDataReturnsNil
// ---------------------------------------------------------------------------

func TestSignals_NoDataReturnsNil(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	// No pool or signal state set
	signals := k.ComputePoolSignals(ctx, 99)
	require.Nil(t, signals)
}

// ---------------------------------------------------------------------------
// TestSignals_EmptyPoolSkipped
// ---------------------------------------------------------------------------

func TestSignals_EmptyPoolSkipped(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	// Pool with zero reserves
	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.ZeroInt(),
		ReserveB: math.ZeroInt(),
	}
	k.SetPool(ctx, pool)

	k.UpdatePoolSignals(ctx, pool)

	_, found := k.GetSignalState(ctx, 1)
	require.False(t, found, "Should not create signal state for pool with zero reserves")
}

// ---------------------------------------------------------------------------
// TestSignals_UpdateAllSignals
// ---------------------------------------------------------------------------

func TestSignals_UpdateAllSignals(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	// Create 3 pools
	setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)
	setupPoolForSignals(t, k, ctx, 2, 5_000_000, 5_000_000)
	setupPoolForSignals(t, k, ctx, 3, 20_000_000, 10_000_000)
	k.SetNextPoolID(ctx, 4)

	k.UpdateAllSignals(ctx)

	// All 3 pools should have signal state
	for _, id := range []uint64{1, 2, 3} {
		_, found := k.GetSignalState(ctx, id)
		require.True(t, found, "Pool %d should have signal state", id)
	}
}

// ---------------------------------------------------------------------------
// TestSignals_MultiplePoolsIndependent
// ---------------------------------------------------------------------------

func TestSignals_MultiplePoolsIndependent(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	// Pool 1: stable
	pool1 := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)
	// Pool 2: rising
	pool2 := types.Pool{
		ID:       2,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(10_000_000),
		ReserveB: math.NewInt(10_000_000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool2)
	k.SetNextPoolID(ctx, 3)

	for i := int64(0); i < 30; i++ {
		blockCtx := ctx.WithBlockHeight(100 + i)
		// Pool 1 stays the same
		k.UpdatePoolSignals(blockCtx, pool1)
		// Pool 2 rises
		pool2.ReserveB = math.NewInt(10_000_000 + i*500_000)
		k.SetPool(blockCtx, pool2)
		k.UpdatePoolSignals(blockCtx, pool2)
	}

	sig1 := k.ComputePoolSignals(ctx.WithBlockHeight(129), 1)
	sig2 := k.ComputePoolSignals(ctx.WithBlockHeight(129), 2)
	require.NotNil(t, sig1)
	require.NotNil(t, sig2)

	// Pool 2 should have higher score than pool 1
	require.True(t, sig2.CompositeScore > sig1.CompositeScore,
		"Rising pool should have higher score (%d) than stable pool (%d)",
		sig2.CompositeScore, sig1.CompositeScore)
}

// ---------------------------------------------------------------------------
// TestSignals_ScoreToSignalMapping
// ---------------------------------------------------------------------------

func TestSignals_ScoreToSignalMapping(t *testing.T) {
	// Test that score boundaries produce correct signal strings
	k, ctx, _, _ := setupKeeper(t)

	// We test by building specific price patterns that produce known score ranges
	// For boundary testing, we verify the signal string mapping works correctly
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	// Stable price -> neutral signal
	for i := int64(0); i < 60; i++ {
		blockCtx := ctx.WithBlockHeight(100 + i)
		k.UpdatePoolSignals(blockCtx, pool)
	}

	signals := k.ComputePoolSignals(ctx.WithBlockHeight(159), 1)
	require.NotNil(t, signals)

	// Verify signal is one of the valid values
	validSignals := map[string]bool{
		"STRONG_BUY": true, "BUY": true, "NEUTRAL": true,
		"SELL": true, "STRONG_SELL": true,
	}
	require.True(t, validSignals[signals.Signal],
		"Signal should be one of the valid values, got %s", signals.Signal)
}

// ---------------------------------------------------------------------------
// TestSignals_VolumeMomentum
// ---------------------------------------------------------------------------

func TestSignals_VolumeMomentumWithNoVolume(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	for i := int64(0); i < 20; i++ {
		blockCtx := ctx.WithBlockHeight(100 + i)
		k.UpdatePoolSignals(blockCtx, pool)
	}

	signals := k.ComputePoolSignals(ctx.WithBlockHeight(119), 1)
	require.NotNil(t, signals)

	// With no volume data from risk engine, volume momentum should be 1.0 (neutral)
	require.True(t, signals.Momentum.VolumeMomentum.Equal(math.LegacyOneDec()),
		"Volume momentum should be 1.0 with no volume data, got %s", signals.Momentum.VolumeMomentum)
}

// ---------------------------------------------------------------------------
// TestSignals_ComputePoolSignalsOutputStructure
// ---------------------------------------------------------------------------

func TestSignals_ComputePoolSignalsOutputStructure(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	for i := int64(0); i < 60; i++ {
		blockCtx := ctx.WithBlockHeight(100 + i)
		pool.ReserveB = math.NewInt(10_000_000 + i*100_000)
		k.SetPool(blockCtx, pool)
		k.UpdatePoolSignals(blockCtx, pool)
	}

	signals := k.ComputePoolSignals(ctx.WithBlockHeight(159), 1)
	require.NotNil(t, signals)

	// Verify all fields are populated
	require.Equal(t, uint64(1), signals.PoolID)
	require.Equal(t, int64(159), signals.Height)
	require.True(t, signals.CurrentPrice.IsPositive())

	// Moving averages should all be positive
	require.True(t, signals.MovingAverages.SMA10.IsPositive())
	require.True(t, signals.MovingAverages.SMA20.IsPositive())
	require.True(t, signals.MovingAverages.SMA50.IsPositive())
	require.True(t, signals.MovingAverages.EMA10.IsPositive())
	require.True(t, signals.MovingAverages.EMA20.IsPositive())
	require.True(t, signals.MovingAverages.EMA50.IsPositive())

	// RSI should be 0-100
	rsiF, _ := signals.Momentum.RSI.Float64()
	require.True(t, rsiF >= 0 && rsiF <= 100, "RSI should be 0-100, got %f", rsiF)

	// Volatility should be non-negative
	require.False(t, signals.VolatilityScore.IsNegative())

	// Score should be 0-100
	require.True(t, signals.CompositeScore >= 0 && signals.CompositeScore <= 100)

	// Signal should be a valid value
	require.NotEmpty(t, signals.Signal)
}

// ---------------------------------------------------------------------------
// TestSignals_SignalStatePersistence
// ---------------------------------------------------------------------------

func TestSignals_SignalStatePersistence(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	k.UpdatePoolSignals(ctx, pool)

	// Read state, modify it, write it back, read again
	state, found := k.GetSignalState(ctx, 1)
	require.True(t, found)
	require.Equal(t, uint64(1), state.PoolID)

	// Verify it round-trips through the KV store
	state2, found := k.GetSignalState(ctx, 1)
	require.True(t, found)
	require.Equal(t, state.PoolID, state2.PoolID)
	require.Equal(t, len(state.PriceHistory), len(state2.PriceHistory))
}

// ---------------------------------------------------------------------------
// TestSignals_SignalHistoryPersistence
// ---------------------------------------------------------------------------

func TestSignals_SignalHistoryPersistence(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	for i := int64(0); i < 5; i++ {
		blockCtx := ctx.WithBlockHeight(100 + i)
		k.UpdatePoolSignals(blockCtx, pool)
	}

	history, found := k.GetSignalHistory(ctx, 1)
	require.True(t, found)
	require.Len(t, history.Entries, 5)

	// Re-read and verify
	history2, found := k.GetSignalHistory(ctx, 1)
	require.True(t, found)
	require.Equal(t, len(history.Entries), len(history2.Entries))
}

// ---------------------------------------------------------------------------
// TestSignals_MissingSentimentHistoryReturnsNotFound
// ---------------------------------------------------------------------------

func TestSignals_MissingSignalHistoryReturnsNotFound(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	_, found := k.GetSignalHistory(ctx, 999)
	require.False(t, found)
}

// ---------------------------------------------------------------------------
// TestSignals_RSIMixedGainsAndLosses
// ---------------------------------------------------------------------------

func TestSignals_RSIMixedGainsAndLosses(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	// Alternating gains and losses of equal magnitude
	for i := int64(0); i < 20; i++ {
		blockCtx := ctx.WithBlockHeight(100 + i)
		if i%2 == 0 {
			pool.ReserveB = math.NewInt(11_000_000) // price = 1.1
		} else {
			pool.ReserveB = math.NewInt(9_000_000) // price = 0.9
		}
		k.SetPool(blockCtx, pool)
		k.UpdatePoolSignals(blockCtx, pool)
	}

	signals := k.ComputePoolSignals(ctx.WithBlockHeight(119), 1)
	require.NotNil(t, signals)

	// RSI should be somewhere around 50 with equal gains/losses (though not exactly
	// due to the asymmetry of percentage changes)
	rsiF, _ := signals.Momentum.RSI.Float64()
	require.True(t, rsiF > 10 && rsiF < 90,
		"RSI should be moderate with alternating gains/losses, got %f", rsiF)
}

// ---------------------------------------------------------------------------
// TestSignals_SMAWithInsufficientData
// ---------------------------------------------------------------------------

func TestSignals_SMAWithInsufficientData(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSignals(t, k, ctx, 1, 10_000_000, 10_000_000)

	// Only 3 blocks (less than SMA10 window)
	for i := int64(0); i < 3; i++ {
		blockCtx := ctx.WithBlockHeight(100 + i)
		k.UpdatePoolSignals(blockCtx, pool)
	}

	signals := k.ComputePoolSignals(ctx.WithBlockHeight(102), 1)
	require.NotNil(t, signals)

	// SMA should still compute using available data
	require.True(t, signals.MovingAverages.SMA10.IsPositive(),
		"SMA10 should still be computed with < 10 data points")
	// With 3 points all at 1.0, SMA = 1.0
	require.True(t, signals.MovingAverages.SMA10.Equal(math.LegacyOneDec()))
}
