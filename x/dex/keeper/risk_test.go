package keeper_test

import (
	"testing"

	"cosmossdk.io/math"

	"github.com/stretchr/testify/require"

	"syreen/x/dex/types"
)

func TestRiskEngine_InitializesState(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(1000000),
		ReserveB: math.NewInt(1000000),
	}
	k.SetPool(ctx, pool)

	k.RunRiskChecks(ctx)

	state, found := k.GetRiskState(ctx, 1)
	require.True(t, found)
	require.Equal(t, uint64(1), state.PoolID)
	require.Equal(t, math.NewInt(1000000), state.PrevReserveA)
}

func TestRiskEngine_PoolNotHaltedNormally(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(1000000),
		ReserveB: math.NewInt(1000000),
	}
	k.SetPool(ctx, pool)

	k.RunRiskChecks(ctx)

	halted, _ := k.IsPoolHalted(ctx, 1)
	require.False(t, halted)
}

func TestRiskEngine_HaltsOnPriceCrash(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	// Initialize with balanced pool
	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(1000000),
		ReserveB: math.NewInt(1000000),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	// Simulate 35% price crash (reserve ratio changes dramatically)
	pool.ReserveA = math.NewInt(1500000) // more tokenA
	pool.ReserveB = math.NewInt(600000)  // much less tokenB
	// Price was 1.0, now is 0.4 — a 60% drop
	k.SetPool(ctx, pool)

	k.RunRiskChecks(ctx)

	halted, reason := k.IsPoolHalted(ctx, 1)
	require.True(t, halted)
	require.Contains(t, reason, "price dropped")
}

func TestRiskEngine_HaltsOnPriceSpike(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(1000000),
		ReserveB: math.NewInt(1000000),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	// Simulate 60% price spike
	pool.ReserveA = math.NewInt(600000)
	pool.ReserveB = math.NewInt(1600000)
	k.SetPool(ctx, pool)

	k.RunRiskChecks(ctx)

	halted, reason := k.IsPoolHalted(ctx, 1)
	require.True(t, halted)
	require.Contains(t, reason, "price spiked")
}

func TestRiskEngine_HaltsOnReserveDrain(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(1000000),
		ReserveB: math.NewInt(1000000),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	// Drain 85% of reserve_a
	pool.ReserveA = math.NewInt(100000)
	k.SetPool(ctx, pool)

	k.RunRiskChecks(ctx)

	halted, reason := k.IsPoolHalted(ctx, 1)
	require.True(t, halted)
	require.Contains(t, reason, "reserve_a drained")
}

func TestRiskEngine_HaltExpires(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(1000000),
		ReserveB: math.NewInt(1000000),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	// Trigger halt
	pool.ReserveA = math.NewInt(100000)
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	halted, _ := k.IsPoolHalted(ctx, 1)
	require.True(t, halted)

	// Advance past halt duration (100 blocks)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 101)
	k.RunRiskChecks(ctx)

	halted, _ = k.IsPoolHalted(ctx, 1)
	require.False(t, halted)
}

func TestRiskEngine_SwapRejectedWhenHalted(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(1000000),
		ReserveB: math.NewInt(1000000),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	// Trigger halt
	pool.ReserveA = math.NewInt(100000)
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	// Try to swap — should be rejected
	err := k.CheckSwapRisk(ctx, 1, math.NewInt(1000))
	require.Error(t, err)
	require.Contains(t, err.Error(), "halted")
}

func TestRiskEngine_VolumeLimit(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(1000000000000),
		ReserveB: math.NewInt(1000000000000),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	// First swap within limit — should pass
	err := k.CheckSwapRisk(ctx, 1, math.NewInt(100000000000))
	require.NoError(t, err)

	// Second swap that exceeds per-block limit — should fail
	err = k.CheckSwapRisk(ctx, 1, math.NewInt(500000000000))
	require.Error(t, err)
	require.Contains(t, err.Error(), "volume limit")
}

func TestRiskEngine_NormalTradeDoesNotHalt(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	pool := types.Pool{
		ID:       1,
		DenomA:   "tokenA",
		DenomB:   "tokenB",
		ReserveA: math.NewInt(1000000),
		ReserveB: math.NewInt(1000000),
		SwapFee:  math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	// Small 5% price move — should NOT halt
	pool.ReserveA = math.NewInt(1050000)
	pool.ReserveB = math.NewInt(952381)
	k.SetPool(ctx, pool)
	k.RunRiskChecks(ctx)

	halted, _ := k.IsPoolHalted(ctx, 1)
	require.False(t, halted)

	// Swap risk check should pass
	err := k.CheckSwapRisk(ctx, 1, math.NewInt(1000))
	require.NoError(t, err)
}
