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

func setupPoolForSentiment(t *testing.T, k keeper.Keeper, ctx sdk.Context) types.Pool {
	t.Helper()
	pool := types.Pool{
		ID:          1,
		DenomA:      "usyreen",
		DenomB:      "uusdc",
		ReserveA:    math.NewInt(10_000_000),
		ReserveB:    math.NewInt(10_000_000),
		TotalShares: math.NewInt(10_000_000),
		SwapFee:     math.LegacyNewDecWithPrec(3, 3),
		Creator:     testAddr(),
		CreatedAt:   ctx.BlockHeight(),
	}
	k.SetPool(ctx, pool)
	k.SetNextPoolID(ctx, 2)
	return pool
}

// ---------------------------------------------------------------------------
// TestSentiment_BlockTradingDataCRUD
// ---------------------------------------------------------------------------

func TestSentiment_BlockTradingDataCRUD(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	// Initially no data
	_, found := k.GetBlockTradingData(ctx, 1, 100)
	require.False(t, found)

	// Set data
	data := keeper.NewBlockTradingData(1, 100)
	data.BuyVolume = math.NewInt(5000)
	data.SellVolume = math.NewInt(3000)
	data.BuyCount = 10
	data.SellCount = 5
	k.SetBlockTradingData(ctx, data)

	// Retrieve
	got, found := k.GetBlockTradingData(ctx, 1, 100)
	require.True(t, found)
	require.Equal(t, math.NewInt(5000), got.BuyVolume)
	require.Equal(t, math.NewInt(3000), got.SellVolume)
	require.Equal(t, uint64(10), got.BuyCount)
	require.Equal(t, uint64(5), got.SellCount)

	// Delete
	k.DeleteBlockTradingData(ctx, 1, 100)
	_, found = k.GetBlockTradingData(ctx, 1, 100)
	require.False(t, found)
}

// ---------------------------------------------------------------------------
// TestSentiment_BlockTradingDataRange
// ---------------------------------------------------------------------------

func TestSentiment_BlockTradingDataRange(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	// Store data for blocks 90-99
	for h := int64(90); h < 100; h++ {
		data := keeper.NewBlockTradingData(1, h)
		data.BuyVolume = math.NewInt(100)
		data.SellVolume = math.NewInt(50)
		data.BuyCount = 1
		k.SetBlockTradingData(ctx, data)
	}

	// Query range 90-99
	results := k.GetBlockTradingDataRange(ctx, 1, 90, 99)
	require.Len(t, results, 10)

	// Query partial range
	results = k.GetBlockTradingDataRange(ctx, 1, 95, 99)
	require.Len(t, results, 5)

	// Query empty range
	results = k.GetBlockTradingDataRange(ctx, 1, 200, 210)
	require.Len(t, results, 0)
}

// ---------------------------------------------------------------------------
// TestSentiment_RecordSwap
// ---------------------------------------------------------------------------

func TestSentiment_RecordSwap(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSentiment(t, k, ctx)

	// Record a buy (buying DenomA = usyreen by selling DenomB = uusdc)
	k.RecordSwapForSentiment(ctx, pool.ID, math.NewInt(500_000), true, pool)

	data, found := k.GetBlockTradingData(ctx, pool.ID, ctx.BlockHeight())
	require.True(t, found)
	require.Equal(t, math.NewInt(500_000), data.BuyVolume)
	require.True(t, data.SellVolume.IsZero())
	require.Equal(t, uint64(1), data.BuyCount)
	require.Equal(t, uint64(0), data.SellCount)

	// Record a sell
	k.RecordSwapForSentiment(ctx, pool.ID, math.NewInt(300_000), false, pool)

	data, _ = k.GetBlockTradingData(ctx, pool.ID, ctx.BlockHeight())
	require.Equal(t, math.NewInt(500_000), data.BuyVolume)
	require.Equal(t, math.NewInt(300_000), data.SellVolume)
	require.Equal(t, uint64(1), data.BuyCount)
	require.Equal(t, uint64(1), data.SellCount)
}

// ---------------------------------------------------------------------------
// TestSentiment_RecordWhaleSwap
// ---------------------------------------------------------------------------

func TestSentiment_RecordWhaleSwap(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSentiment(t, k, ctx)

	// A whale trade is >1% of pool reserves
	// Pool has 10M each side, so 1% = 100K
	// Record a 150K buy — should be flagged as whale
	k.RecordSwapForSentiment(ctx, pool.ID, math.NewInt(150_000), true, pool)

	data, found := k.GetBlockTradingData(ctx, pool.ID, ctx.BlockHeight())
	require.True(t, found)
	require.Equal(t, uint64(1), data.WhaleBuyCount)
	require.Equal(t, math.NewInt(150_000), data.WhaleBuyVolume)
	require.Equal(t, uint64(0), data.WhaleSellCount)

	// Record a small trade (50K) — not a whale
	k.RecordSwapForSentiment(ctx, pool.ID, math.NewInt(50_000), true, pool)
	data, _ = k.GetBlockTradingData(ctx, pool.ID, ctx.BlockHeight())
	require.Equal(t, uint64(1), data.WhaleBuyCount) // still 1 whale
	require.Equal(t, uint64(2), data.BuyCount)       // 2 total buys

	// Record a whale sell
	k.RecordSwapForSentiment(ctx, pool.ID, math.NewInt(200_000), false, pool)
	data, _ = k.GetBlockTradingData(ctx, pool.ID, ctx.BlockHeight())
	require.Equal(t, uint64(1), data.WhaleSellCount)
	require.Equal(t, math.NewInt(200_000), data.WhaleSellVolume)
}

// ---------------------------------------------------------------------------
// TestSentiment_RecordLiquidity
// ---------------------------------------------------------------------------

func TestSentiment_RecordLiquidity(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	_ = setupPoolForSentiment(t, k, ctx)

	// Record LP add
	k.RecordLiquidityForSentiment(ctx, 1, math.NewInt(2_000_000), true)

	data, found := k.GetBlockTradingData(ctx, 1, ctx.BlockHeight())
	require.True(t, found)
	require.Equal(t, math.NewInt(2_000_000), data.LPAddVolume)
	require.Equal(t, uint64(1), data.LPAddCount)
	require.True(t, data.LPRemoveVolume.IsZero())

	// Record LP remove
	k.RecordLiquidityForSentiment(ctx, 1, math.NewInt(500_000), false)
	data, _ = k.GetBlockTradingData(ctx, 1, ctx.BlockHeight())
	require.Equal(t, math.NewInt(500_000), data.LPRemoveVolume)
	require.Equal(t, uint64(1), data.LPRemoveCount)
}

// ---------------------------------------------------------------------------
// TestSentiment_NeutralWithNoActivity
// ---------------------------------------------------------------------------

func TestSentiment_NeutralWithNoActivity(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	_ = setupPoolForSentiment(t, k, ctx)

	// Run sentiment update with no trading data
	k.UpdateSentiment(ctx)

	state, found := k.GetSentimentState(ctx, 1)
	require.True(t, found)
	require.Equal(t, int64(50), state.FearGreedIndex)
	require.Equal(t, "Neutral", state.FearGreedLabel)
	require.Equal(t, keeper.MoodNeutral, state.Mood)
	require.Equal(t, int64(50), state.Indicators.BuySellRatio)
	require.Equal(t, int64(50), state.Indicators.LiquidityFlow)
	require.Equal(t, int64(50), state.Indicators.WhaleActivity)
	require.Equal(t, int64(50), state.Indicators.VolumeTrend)
}

// ---------------------------------------------------------------------------
// TestSentiment_BullishOnBuyingPressure
// ---------------------------------------------------------------------------

func TestSentiment_BullishOnBuyingPressure(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSentiment(t, k, ctx)

	// Simulate heavy buying over many blocks
	for h := int64(50); h < 100; h++ {
		blockCtx := ctx.WithBlockHeight(h)
		data := keeper.NewBlockTradingData(pool.ID, h)
		data.BuyVolume = math.NewInt(800_000)
		data.SellVolume = math.NewInt(100_000)
		data.BuyCount = 8
		data.SellCount = 1
		// Add whale buys
		data.WhaleBuyCount = 2
		data.WhaleBuyVolume = math.NewInt(300_000)
		data.WhaleSellCount = 0
		data.WhaleSellVolume = math.ZeroInt()
		// LP additions
		data.LPAddVolume = math.NewInt(500_000)
		data.LPRemoveVolume = math.NewInt(50_000)
		data.LPAddCount = 2
		k.SetBlockTradingData(blockCtx, data)
	}

	// Run sentiment at block 100
	k.UpdateSentiment(ctx)

	state, found := k.GetSentimentState(ctx, pool.ID)
	require.True(t, found)

	// Should be bullish — heavy buying, whale accumulation, LP inflows
	require.Greater(t, state.FearGreedIndex, int64(60), "Fear & Greed should be >60 (greedy) with heavy buying")
	require.Contains(t, []keeper.MarketMood{keeper.MoodSlightlyBullish, keeper.MoodBullish}, state.Mood)
	require.Greater(t, state.Indicators.BuySellRatio, int64(50))
	require.Greater(t, state.Indicators.WhaleActivity, int64(50))
	require.Greater(t, state.Indicators.LiquidityFlow, int64(50))
}

// ---------------------------------------------------------------------------
// TestSentiment_BearishOnSellingPressure
// ---------------------------------------------------------------------------

func TestSentiment_BearishOnSellingPressure(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSentiment(t, k, ctx)

	// Simulate heavy selling with LP exits
	for h := int64(50); h < 100; h++ {
		data := keeper.NewBlockTradingData(pool.ID, h)
		data.BuyVolume = math.NewInt(50_000)
		data.SellVolume = math.NewInt(900_000)
		data.BuyCount = 1
		data.SellCount = 9
		data.WhaleSellCount = 3
		data.WhaleSellVolume = math.NewInt(500_000)
		data.WhaleBuyCount = 0
		data.WhaleBuyVolume = math.ZeroInt()
		// LP exodus
		data.LPAddVolume = math.NewInt(10_000)
		data.LPRemoveVolume = math.NewInt(800_000)
		data.LPRemoveCount = 5
		k.SetBlockTradingData(ctx.WithBlockHeight(h), data)
	}

	k.UpdateSentiment(ctx)

	state, found := k.GetSentimentState(ctx, pool.ID)
	require.True(t, found)

	require.Less(t, state.FearGreedIndex, int64(40), "Fear & Greed should be <40 (fearful) with heavy selling")
	require.Contains(t, []keeper.MarketMood{keeper.MoodSlightlyBearish, keeper.MoodBearish}, state.Mood)
	require.Less(t, state.Indicators.BuySellRatio, int64(50))
	require.Less(t, state.Indicators.WhaleActivity, int64(50))
	require.Less(t, state.Indicators.LiquidityFlow, int64(50))
}

// ---------------------------------------------------------------------------
// TestSentiment_FearGreedLabels
// ---------------------------------------------------------------------------

func TestSentiment_FearGreedLabels(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSentiment(t, k, ctx)

	// To fully control all indicators, we set buy/sell, whale, and LP data together.
	// All indicators are driven by the same bias direction.
	tests := []struct {
		name        string
		buyPct      int64 // buy volume percentage
		whaleBuyPct int64 // whale buy pct
		lpAddPct    int64 // LP add pct
		expectLabel string
		expectMood  keeper.MarketMood
	}{
		{"Extreme Fear", 2, 0, 0, "Extreme Fear", keeper.MoodBearish},
		{"Fear", 20, 15, 15, "Fear", keeper.MoodSlightlyBearish},
		{"Neutral", 50, 50, 50, "Neutral", keeper.MoodNeutral},
		{"Greed", 80, 85, 85, "Greed", keeper.MoodSlightlyBullish},
		{"Extreme Greed", 98, 100, 100, "Extreme Greed", keeper.MoodBullish},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for h := int64(1); h < 100; h++ {
				data := keeper.NewBlockTradingData(pool.ID, h)
				data.BuyVolume = math.NewInt(tc.buyPct * 1000)
				data.SellVolume = math.NewInt((100 - tc.buyPct) * 1000)
				data.BuyCount = 1
				data.SellCount = 1
				data.WhaleBuyVolume = math.NewInt(tc.whaleBuyPct * 100)
				data.WhaleSellVolume = math.NewInt((100 - tc.whaleBuyPct) * 100)
				if tc.whaleBuyPct > 0 {
					data.WhaleBuyCount = 1
				}
				if tc.whaleBuyPct < 100 {
					data.WhaleSellCount = 1
				}
				data.LPAddVolume = math.NewInt(tc.lpAddPct * 100)
				data.LPRemoveVolume = math.NewInt((100 - tc.lpAddPct) * 100)
				k.SetBlockTradingData(ctx.WithBlockHeight(h), data)
			}

			k.UpdateSentiment(ctx)

			state, found := k.GetSentimentState(ctx, pool.ID)
			require.True(t, found)
			require.Equal(t, tc.expectLabel, state.FearGreedLabel)
			require.Equal(t, tc.expectMood, state.Mood)
		})
	}
}

// ---------------------------------------------------------------------------
// TestSentiment_VolumeTrendIncreasing
// ---------------------------------------------------------------------------

func TestSentiment_VolumeTrendIncreasing(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSentiment(t, k, ctx)

	// First half: low volume; second half: high volume
	midHeight := int64(50)
	for h := int64(1); h < 100; h++ {
		data := keeper.NewBlockTradingData(pool.ID, h)
		if h <= midHeight {
			data.BuyVolume = math.NewInt(100)
			data.SellVolume = math.NewInt(100)
		} else {
			data.BuyVolume = math.NewInt(10_000)
			data.SellVolume = math.NewInt(10_000)
		}
		data.BuyCount = 1
		data.SellCount = 1
		k.SetBlockTradingData(ctx.WithBlockHeight(h), data)
	}

	k.UpdateSentiment(ctx)

	state, found := k.GetSentimentState(ctx, pool.ID)
	require.True(t, found)
	require.Greater(t, state.Indicators.VolumeTrend, int64(50), "Volume trend should be >50 when volume is increasing")
}

// ---------------------------------------------------------------------------
// TestSentiment_VolumeTrendDecreasing
// ---------------------------------------------------------------------------

func TestSentiment_VolumeTrendDecreasing(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSentiment(t, k, ctx)

	// First half: high volume; second half: low volume
	midHeight := int64(50)
	for h := int64(1); h < 100; h++ {
		data := keeper.NewBlockTradingData(pool.ID, h)
		if h <= midHeight {
			data.BuyVolume = math.NewInt(10_000)
			data.SellVolume = math.NewInt(10_000)
		} else {
			data.BuyVolume = math.NewInt(100)
			data.SellVolume = math.NewInt(100)
		}
		data.BuyCount = 1
		data.SellCount = 1
		k.SetBlockTradingData(ctx.WithBlockHeight(h), data)
	}

	k.UpdateSentiment(ctx)

	state, found := k.GetSentimentState(ctx, pool.ID)
	require.True(t, found)
	require.Less(t, state.Indicators.VolumeTrend, int64(50), "Volume trend should be <50 when volume is decreasing")
}

// ---------------------------------------------------------------------------
// TestSentiment_TradeSizeTrendIncreasing
// ---------------------------------------------------------------------------

func TestSentiment_TradeSizeTrendIncreasing(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSentiment(t, k, ctx)

	midHeight := int64(50)
	for h := int64(1); h < 100; h++ {
		data := keeper.NewBlockTradingData(pool.ID, h)
		if h <= midHeight {
			// Small trades: 10 trades of 100 each = avg 100
			data.BuyVolume = math.NewInt(1000)
			data.BuyCount = 10
			data.SellVolume = math.ZeroInt()
		} else {
			// Large trades: 2 trades of 5000 each = avg 5000
			data.BuyVolume = math.NewInt(10000)
			data.BuyCount = 2
			data.SellVolume = math.ZeroInt()
		}
		k.SetBlockTradingData(ctx.WithBlockHeight(h), data)
	}

	k.UpdateSentiment(ctx)

	state, found := k.GetSentimentState(ctx, pool.ID)
	require.True(t, found)
	require.Greater(t, state.Indicators.TradeSizeTrend, int64(50), "Trade size trend should be >50 when avg trade sizes are increasing")
}

// ---------------------------------------------------------------------------
// TestSentiment_AlertFearSpike
// ---------------------------------------------------------------------------

func TestSentiment_AlertFearSpike(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSentiment(t, k, ctx)

	// Create extreme fear conditions: overwhelming sell pressure, LP exodus
	for h := int64(1); h < 100; h++ {
		data := keeper.NewBlockTradingData(pool.ID, h)
		data.BuyVolume = math.NewInt(1_000)
		data.SellVolume = math.NewInt(99_000)
		data.BuyCount = 1
		data.SellCount = 20
		data.WhaleSellCount = 5
		data.WhaleSellVolume = math.NewInt(50_000)
		data.WhaleBuyCount = 0
		data.WhaleBuyVolume = math.ZeroInt()
		data.LPRemoveVolume = math.NewInt(100_000)
		data.LPAddVolume = math.NewInt(1_000)
		data.LPRemoveCount = 3
		k.SetBlockTradingData(ctx.WithBlockHeight(h), data)
	}

	k.UpdateSentiment(ctx)

	state, found := k.GetSentimentState(ctx, pool.ID)
	require.True(t, found)
	require.LessOrEqual(t, state.FearGreedIndex, int64(20), "Should be extreme fear")

	alerts := k.GetSentimentAlerts(ctx, pool.ID)
	require.NotEmpty(t, alerts.Alerts, "Should have alerts during extreme fear")

	// Check for fear spike alert
	foundFear := false
	for _, alert := range alerts.Alerts {
		if alert.AlertType == keeper.AlertFearSpike {
			foundFear = true
			require.Equal(t, "high", alert.Severity)
		}
	}
	require.True(t, foundFear, "Should have a fear_spike alert")
}

// ---------------------------------------------------------------------------
// TestSentiment_AlertGreedSpike
// ---------------------------------------------------------------------------

func TestSentiment_AlertGreedSpike(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSentiment(t, k, ctx)

	// Create extreme greed conditions
	for h := int64(1); h < 100; h++ {
		data := keeper.NewBlockTradingData(pool.ID, h)
		data.BuyVolume = math.NewInt(99_000)
		data.SellVolume = math.NewInt(1_000)
		data.BuyCount = 20
		data.SellCount = 1
		data.WhaleBuyCount = 5
		data.WhaleBuyVolume = math.NewInt(50_000)
		data.WhaleSellCount = 0
		data.WhaleSellVolume = math.ZeroInt()
		data.LPAddVolume = math.NewInt(100_000)
		data.LPRemoveVolume = math.NewInt(1_000)
		data.LPAddCount = 3
		k.SetBlockTradingData(ctx.WithBlockHeight(h), data)
	}

	k.UpdateSentiment(ctx)

	state, found := k.GetSentimentState(ctx, pool.ID)
	require.True(t, found)
	require.GreaterOrEqual(t, state.FearGreedIndex, int64(80), "Should be extreme greed")

	alerts := k.GetSentimentAlerts(ctx, pool.ID)
	foundGreed := false
	for _, alert := range alerts.Alerts {
		if alert.AlertType == keeper.AlertGreedSpike {
			foundGreed = true
			require.Equal(t, "high", alert.Severity)
		}
	}
	require.True(t, foundGreed, "Should have a greed_spike alert")
}

// ---------------------------------------------------------------------------
// TestSentiment_AlertWhaleAccumulation
// ---------------------------------------------------------------------------

func TestSentiment_AlertWhaleAccumulation(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSentiment(t, k, ctx)

	// Whales buying heavily, no whale sells
	for h := int64(1); h < 100; h++ {
		data := keeper.NewBlockTradingData(pool.ID, h)
		data.BuyVolume = math.NewInt(500_000)
		data.SellVolume = math.NewInt(500_000)
		data.BuyCount = 5
		data.SellCount = 5
		data.WhaleBuyCount = 3
		data.WhaleBuyVolume = math.NewInt(300_000)
		data.WhaleSellCount = 0
		data.WhaleSellVolume = math.ZeroInt()
		k.SetBlockTradingData(ctx.WithBlockHeight(h), data)
	}

	k.UpdateSentiment(ctx)

	alerts := k.GetSentimentAlerts(ctx, pool.ID)
	foundWhale := false
	for _, alert := range alerts.Alerts {
		if alert.AlertType == keeper.AlertWhaleAccumulation {
			foundWhale = true
			require.Equal(t, "medium", alert.Severity)
		}
	}
	require.True(t, foundWhale, "Should have a whale_accumulation alert")
}

// ---------------------------------------------------------------------------
// TestSentiment_AlertLiquidityExodus
// ---------------------------------------------------------------------------

func TestSentiment_AlertLiquidityExodus(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSentiment(t, k, ctx)

	// Heavy LP removals, minimal additions
	for h := int64(1); h < 100; h++ {
		data := keeper.NewBlockTradingData(pool.ID, h)
		data.BuyVolume = math.NewInt(500_000)
		data.SellVolume = math.NewInt(500_000)
		data.BuyCount = 5
		data.SellCount = 5
		data.LPAddVolume = math.NewInt(1_000)
		data.LPRemoveVolume = math.NewInt(500_000)
		data.LPRemoveCount = 5
		k.SetBlockTradingData(ctx.WithBlockHeight(h), data)
	}

	k.UpdateSentiment(ctx)

	alerts := k.GetSentimentAlerts(ctx, pool.ID)
	foundExodus := false
	for _, alert := range alerts.Alerts {
		if alert.AlertType == keeper.AlertLiquidityExodus {
			foundExodus = true
			require.Equal(t, "high", alert.Severity)
		}
	}
	require.True(t, foundExodus, "Should have a liquidity_exodus alert")
}

// ---------------------------------------------------------------------------
// TestSentiment_NoAlertsInNeutralMarket
// ---------------------------------------------------------------------------

func TestSentiment_NoAlertsInNeutralMarket(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSentiment(t, k, ctx)

	// Balanced trading — should produce no alerts
	for h := int64(1); h < 100; h++ {
		data := keeper.NewBlockTradingData(pool.ID, h)
		data.BuyVolume = math.NewInt(500_000)
		data.SellVolume = math.NewInt(500_000)
		data.BuyCount = 5
		data.SellCount = 5
		data.LPAddVolume = math.NewInt(100_000)
		data.LPRemoveVolume = math.NewInt(100_000)
		data.WhaleBuyCount = 1
		data.WhaleSellCount = 1
		data.WhaleBuyVolume = math.NewInt(100_000)
		data.WhaleSellVolume = math.NewInt(100_000)
		k.SetBlockTradingData(ctx.WithBlockHeight(h), data)
	}

	k.UpdateSentiment(ctx)

	state, found := k.GetSentimentState(ctx, pool.ID)
	require.True(t, found)
	// Allow +/- 1 for rounding
	require.InDelta(t, 50, state.FearGreedIndex, 1, "Fear & Greed should be ~50 in balanced market")
	require.Equal(t, "Neutral", state.FearGreedLabel)

	alerts := k.GetSentimentAlerts(ctx, pool.ID)
	require.Empty(t, alerts.Alerts, "No alerts expected in a balanced market")
}

// ---------------------------------------------------------------------------
// TestSentiment_SentimentStatePersistsAcrossBlocks
// ---------------------------------------------------------------------------

func TestSentiment_SentimentStatePersistsAcrossBlocks(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	_ = setupPoolForSentiment(t, k, ctx)

	// Set some data and update
	data := keeper.NewBlockTradingData(1, 99)
	data.BuyVolume = math.NewInt(700_000)
	data.SellVolume = math.NewInt(300_000)
	data.BuyCount = 7
	data.SellCount = 3
	k.SetBlockTradingData(ctx.WithBlockHeight(99), data)

	k.UpdateSentiment(ctx)

	state1, _ := k.GetSentimentState(ctx, 1)
	require.Greater(t, state1.FearGreedIndex, int64(50))

	// Advance block and update again — state should still be accessible
	ctx2 := ctx.WithBlockHeight(101)
	k.UpdateSentiment(ctx2)

	state2, found := k.GetSentimentState(ctx2, 1)
	require.True(t, found)
	require.Equal(t, int64(101), state2.UpdatedAtHeight)
}

// ---------------------------------------------------------------------------
// TestSentiment_MultiplePoolsIndependent
// ---------------------------------------------------------------------------

func TestSentiment_MultiplePoolsIndependent(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	// Create two pools
	pool1 := types.Pool{
		ID: 1, DenomA: "usyreen", DenomB: "uusdc",
		ReserveA: math.NewInt(10_000_000), ReserveB: math.NewInt(10_000_000),
		TotalShares: math.NewInt(10_000_000), SwapFee: math.LegacyNewDecWithPrec(3, 3),
	}
	pool2 := types.Pool{
		ID: 2, DenomA: "uatom", DenomB: "uusdc",
		ReserveA: math.NewInt(5_000_000), ReserveB: math.NewInt(5_000_000),
		TotalShares: math.NewInt(5_000_000), SwapFee: math.LegacyNewDecWithPrec(3, 3),
	}
	k.SetPool(ctx, pool1)
	k.SetPool(ctx, pool2)
	k.SetNextPoolID(ctx, 3)

	// Pool 1: heavy buying (bullish)
	for h := int64(1); h < 100; h++ {
		data := keeper.NewBlockTradingData(1, h)
		data.BuyVolume = math.NewInt(900_000)
		data.SellVolume = math.NewInt(100_000)
		data.BuyCount = 9
		data.SellCount = 1
		k.SetBlockTradingData(ctx.WithBlockHeight(h), data)
	}

	// Pool 2: heavy selling (bearish)
	for h := int64(1); h < 100; h++ {
		data := keeper.NewBlockTradingData(2, h)
		data.BuyVolume = math.NewInt(100_000)
		data.SellVolume = math.NewInt(900_000)
		data.BuyCount = 1
		data.SellCount = 9
		k.SetBlockTradingData(ctx.WithBlockHeight(h), data)
	}

	k.UpdateSentiment(ctx)

	state1, _ := k.GetSentimentState(ctx, 1)
	state2, _ := k.GetSentimentState(ctx, 2)

	require.Greater(t, state1.FearGreedIndex, int64(50), "Pool 1 should be bullish")
	require.Less(t, state2.FearGreedIndex, int64(50), "Pool 2 should be bearish")

	// Moods should differ
	require.NotEqual(t, state1.Mood, state2.Mood)
}

// ---------------------------------------------------------------------------
// TestSentiment_GetSentimentHistory
// ---------------------------------------------------------------------------

func TestSentiment_GetSentimentHistory(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	_ = setupPoolForSentiment(t, k, ctx)

	// Store 50 blocks of data
	for h := int64(50); h < 100; h++ {
		data := keeper.NewBlockTradingData(1, h)
		data.BuyVolume = math.NewInt(1000)
		k.SetBlockTradingData(ctx.WithBlockHeight(h), data)
	}

	history := k.GetSentimentHistory(ctx, 1, 50)
	require.Len(t, history, 50)

	// Limit to 10
	history = k.GetSentimentHistory(ctx, 1, 10)
	require.Len(t, history, 10)
}

// ---------------------------------------------------------------------------
// TestSentiment_QueryEndpoints
// ---------------------------------------------------------------------------

func TestSentiment_QuerySentiment(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSentiment(t, k, ctx)

	// Populate data
	for h := int64(1); h < 100; h++ {
		data := keeper.NewBlockTradingData(pool.ID, h)
		data.BuyVolume = math.NewInt(600_000)
		data.SellVolume = math.NewInt(400_000)
		data.BuyCount = 6
		data.SellCount = 4
		k.SetBlockTradingData(ctx.WithBlockHeight(h), data)
	}
	k.UpdateSentiment(ctx)

	qs := keeper.NewQueryServerImpl(k)

	// Query sentiment
	resp, err := qs.Sentiment(ctx, &types.QuerySentimentRequest{PoolID: pool.ID})
	require.NoError(t, err)
	require.Equal(t, pool.ID, resp.PoolID)
	require.Greater(t, resp.FearGreedIndex, int64(50))
	require.NotEmpty(t, resp.FearGreedLabel)
	require.NotEmpty(t, resp.Mood)
}

func TestSentiment_QuerySentimentHistory(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	_ = setupPoolForSentiment(t, k, ctx)

	for h := int64(50); h < 100; h++ {
		data := keeper.NewBlockTradingData(1, h)
		data.BuyVolume = math.NewInt(1000)
		k.SetBlockTradingData(ctx.WithBlockHeight(h), data)
	}

	qs := keeper.NewQueryServerImpl(k)
	resp, err := qs.SentimentHistory(ctx, &types.QuerySentimentHistoryRequest{PoolID: 1, Limit: 50})
	require.NoError(t, err)
	require.Equal(t, uint64(1), resp.PoolID)
	require.NotNil(t, resp.History)
}

func TestSentiment_QuerySentimentAlerts(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSentiment(t, k, ctx)

	// Create extreme fear conditions: heavy selling, whale sells, LP exodus
	for h := int64(1); h < 100; h++ {
		data := keeper.NewBlockTradingData(pool.ID, h)
		data.BuyVolume = math.NewInt(1_000)
		data.SellVolume = math.NewInt(99_000)
		data.BuyCount = 1
		data.SellCount = 20
		data.WhaleSellCount = 5
		data.WhaleSellVolume = math.NewInt(50_000)
		data.WhaleBuyCount = 0
		data.WhaleBuyVolume = math.ZeroInt()
		data.LPRemoveVolume = math.NewInt(100_000)
		data.LPAddVolume = math.NewInt(1_000)
		data.LPRemoveCount = 3
		k.SetBlockTradingData(ctx.WithBlockHeight(h), data)
	}
	k.UpdateSentiment(ctx)

	qs := keeper.NewQueryServerImpl(k)
	resp, err := qs.SentimentAlerts(ctx, &types.QuerySentimentAlertsRequest{PoolID: pool.ID})
	require.NoError(t, err)
	require.Equal(t, pool.ID, resp.PoolID)
	require.NotNil(t, resp.Alerts)
}

func TestSentiment_QuerySentimentNoPool(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	qs := keeper.NewQueryServerImpl(k)

	// Query for non-existent pool returns neutral defaults
	resp, err := qs.Sentiment(ctx, &types.QuerySentimentRequest{PoolID: 999})
	require.NoError(t, err)
	require.Equal(t, int64(50), resp.FearGreedIndex)
	require.Equal(t, "Neutral", resp.FearGreedLabel)
}

func TestSentiment_QuerySentimentZeroPoolID(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	qs := keeper.NewQueryServerImpl(k)

	_, err := qs.Sentiment(ctx, &types.QuerySentimentRequest{PoolID: 0})
	require.Error(t, err)
}

// ---------------------------------------------------------------------------
// TestSentiment_IntegrationWithSwap — verify Swap() records sentiment data
// ---------------------------------------------------------------------------

func TestSentiment_IntegrationWithSwap(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()
	trader := testAddr2()

	// Create pool via the keeper (full flow)
	bk.fundAccount(creator, sdk.NewCoins(
		sdk.NewInt64Coin("usyreen", 10_000_000),
		sdk.NewInt64Coin("uusdc", 10_000_000),
	))
	poolID, err := k.CreatePool(ctx, creator, "usyreen", "uusdc", math.NewInt(10_000_000), math.NewInt(10_000_000))
	require.NoError(t, err)

	// Fund module for swap outputs
	bk.balances[types.ModuleName] = bk.balances[types.ModuleName].Add(
		sdk.NewCoin("usyreen", math.NewInt(10_000_000)),
		sdk.NewCoin("uusdc", math.NewInt(10_000_000)),
	)

	// Perform a swap (buying uusdc by selling usyreen)
	bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1_000_000)))
	_, err = k.Swap(ctx, trader, poolID, sdk.NewInt64Coin("usyreen", 1_000_000), math.ZeroInt())
	require.NoError(t, err)

	// Check that swap was recorded for sentiment
	data, found := k.GetBlockTradingData(ctx, poolID, ctx.BlockHeight())
	require.True(t, found)

	// Selling usyreen (DenomA) = selling DenomA = not a "buy" (buy = buying DenomA)
	// tokenIn.Denom == pool.DenomA means isBuy=false (selling DenomA)
	require.True(t, data.SellVolume.IsPositive(), "Should record sell volume for selling DenomA")
	require.Equal(t, uint64(1), data.SellCount)
}

// ---------------------------------------------------------------------------
// TestSentiment_IntegrationWithLiquidity — verify AddLiquidity/RemoveLiquidity record data
// ---------------------------------------------------------------------------

func TestSentiment_IntegrationWithAddLiquidity(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()
	provider := testAddr2()

	bk.fundAccount(creator, sdk.NewCoins(
		sdk.NewInt64Coin("usyreen", 10_000_000),
		sdk.NewInt64Coin("uusdc", 10_000_000),
	))
	poolID, err := k.CreatePool(ctx, creator, "usyreen", "uusdc", math.NewInt(10_000_000), math.NewInt(10_000_000))
	require.NoError(t, err)

	// Add liquidity
	bk.fundAccount(provider, sdk.NewCoins(
		sdk.NewInt64Coin("usyreen", 2_000_000),
		sdk.NewInt64Coin("uusdc", 2_000_000),
	))
	_, err = k.AddLiquidity(ctx, provider, poolID, math.NewInt(2_000_000), math.NewInt(2_000_000), math.ZeroInt())
	require.NoError(t, err)

	// Check sentiment data recorded
	data, found := k.GetBlockTradingData(ctx, poolID, ctx.BlockHeight())
	require.True(t, found)
	require.True(t, data.LPAddVolume.IsPositive(), "Should record LP add volume")
	require.Equal(t, uint64(1), data.LPAddCount)

	// CreatePool also records an LP add — check total
	// CreatePool adds liquidity too, so we may have 2 LP add events
	require.GreaterOrEqual(t, data.LPAddCount, uint64(1))
}

// ---------------------------------------------------------------------------
// TestSentiment_ClampScore
// ---------------------------------------------------------------------------

func TestSentiment_ClampScore(t *testing.T) {
	// Test via the indicator functions by setting up extreme data
	k, ctx, _, _ := setupKeeper(t)
	pool := setupPoolForSentiment(t, k, ctx)

	// All buys, no sells — should clamp at 100
	data := keeper.NewBlockTradingData(pool.ID, 99)
	data.BuyVolume = math.NewInt(1_000_000)
	data.SellVolume = math.ZeroInt()
	data.BuyCount = 10
	k.SetBlockTradingData(ctx.WithBlockHeight(99), data)

	k.UpdateSentiment(ctx)

	state, _ := k.GetSentimentState(ctx, pool.ID)
	require.Equal(t, int64(100), state.Indicators.BuySellRatio, "100% buys should give ratio of 100")
}

// ---------------------------------------------------------------------------
// TestSentiment_SentimentAlertsCRUD
// ---------------------------------------------------------------------------

func TestSentiment_SentimentAlertsCRUD(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	// Default empty
	alerts := k.GetSentimentAlerts(ctx, 1)
	require.Equal(t, uint64(1), alerts.PoolID)
	require.Empty(t, alerts.Alerts)

	// Set some alerts
	alerts.Alerts = []keeper.SentimentAlert{
		{PoolID: 1, AlertType: keeper.AlertFearSpike, Message: "test", Severity: "high", Height: 100, Value: 10},
	}
	k.SetSentimentAlerts(ctx, alerts)

	got := k.GetSentimentAlerts(ctx, 1)
	require.Len(t, got.Alerts, 1)
	require.Equal(t, keeper.AlertFearSpike, got.Alerts[0].AlertType)
}

// ---------------------------------------------------------------------------
// TestSentiment_SentimentStateCRUD
// ---------------------------------------------------------------------------

func TestSentiment_SentimentStateCRUD(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	_, found := k.GetSentimentState(ctx, 1)
	require.False(t, found)

	state := keeper.SentimentState{
		PoolID:          1,
		FearGreedIndex:  75,
		FearGreedLabel:  "Greed",
		Mood:            keeper.MoodSlightlyBullish,
		UpdatedAtHeight: 100,
		BlocksAnalyzed:  50,
	}
	k.SetSentimentState(ctx, state)

	got, found := k.GetSentimentState(ctx, 1)
	require.True(t, found)
	require.Equal(t, int64(75), got.FearGreedIndex)
	require.Equal(t, "Greed", got.FearGreedLabel)
	require.Equal(t, keeper.MoodSlightlyBullish, got.Mood)
}
