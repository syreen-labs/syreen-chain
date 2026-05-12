package keeper_test

import (
	"fmt"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/require"

	"syreen/x/dex/keeper"
	"syreen/x/dex/types"
)

func testAddr5() string {
	return sdk.AccAddress([]byte("fifth__addr_padding_")).String()
}

// helper: create a pool and fund the module account for swap outputs.
func setupPoolForCopyTrade(t *testing.T, k keeper.Keeper, ctx sdk.Context, bk *mockBankKeeper) (uint64, string) {
	t.Helper()
	creator := testAddr()
	bk.fundAccount(creator, sdk.NewCoins(
		sdk.NewInt64Coin("usyreen", 100_000_000),
		sdk.NewInt64Coin("uusdc", 100_000_000),
	))
	poolID, err := k.CreatePool(ctx, creator, "usyreen", "uusdc", math.NewInt(100_000_000), math.NewInt(100_000_000))
	require.NoError(t, err)

	// Fund module account so swaps can send output tokens
	bk.balances[types.ModuleName] = bk.balances[types.ModuleName].Add(
		sdk.NewCoin("usyreen", math.NewInt(100_000_000)),
		sdk.NewCoin("uusdc", math.NewInt(100_000_000)),
	)
	return poolID, creator
}

// ---------------------------------------------------------------------------
// TestFollowTrader
// ---------------------------------------------------------------------------

func TestFollowTrader(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	follower := testAddr()
	trader := testAddr2()

	err := k.FollowTrader(ctx, follower, trader, math.NewInt(500_000), math.NewInt(10_000_000), 5000)
	require.NoError(t, err)

	// Verify settings stored
	settings, exists := k.GetCopySettings(ctx, follower, trader)
	require.True(t, exists)
	require.Equal(t, follower, settings.Follower)
	require.Equal(t, trader, settings.Trader)
	require.Equal(t, math.NewInt(500_000), settings.MaxPerTrade)
	require.Equal(t, math.NewInt(10_000_000), settings.TotalBudget)
	require.True(t, settings.Spent.IsZero())
	require.Equal(t, int64(5000), settings.CopyRatio)
	require.True(t, settings.Active)
	require.Equal(t, int64(100), settings.CreatedAt) // block height 100

	// Verify follower count updated
	stats, found := k.GetTraderStats(ctx, trader)
	require.True(t, found)
	require.Equal(t, int64(1), stats.FollowerCount)

	// Verify follower appears in GetFollowers
	followers := k.GetFollowers(ctx, trader)
	require.Len(t, followers, 1)
	require.Equal(t, follower, followers[0])

	// Verify GetFollowing returns the setting
	following := k.GetFollowing(ctx, follower)
	require.Len(t, following, 1)
	require.Equal(t, trader, following[0].Trader)
}

// ---------------------------------------------------------------------------
// TestFollowTraderSelf
// ---------------------------------------------------------------------------

func TestFollowTraderSelf(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	addr := testAddr()

	err := k.FollowTrader(ctx, addr, addr, math.NewInt(500_000), math.NewInt(10_000_000), 5000)
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot follow yourself")
}

// ---------------------------------------------------------------------------
// TestFollowTraderDuplicate
// ---------------------------------------------------------------------------

func TestFollowTraderDuplicate(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	follower := testAddr()
	trader := testAddr2()

	err := k.FollowTrader(ctx, follower, trader, math.NewInt(500_000), math.NewInt(10_000_000), 5000)
	require.NoError(t, err)

	err = k.FollowTrader(ctx, follower, trader, math.NewInt(500_000), math.NewInt(10_000_000), 5000)
	require.Error(t, err)
	require.Contains(t, err.Error(), "already following")
}

// ---------------------------------------------------------------------------
// TestFollowTraderMaxFollows
// ---------------------------------------------------------------------------

func TestFollowTraderMaxFollows(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	follower := testAddr()

	// Follow 20 unique traders (the max)
	for i := 0; i < keeper.MaxFollowsPerUser; i++ {
		traderAddr := sdk.AccAddress([]byte(fmt.Sprintf("trader_%02d_padding__", i))).String()
		err := k.FollowTrader(ctx, follower, traderAddr, math.NewInt(500_000), math.NewInt(10_000_000), 5000)
		require.NoError(t, err, "follow %d should succeed", i)
	}

	// 21st follow should fail
	extraTrader := sdk.AccAddress([]byte("trader_extra_paddin_")).String()
	err := k.FollowTrader(ctx, follower, extraTrader, math.NewInt(500_000), math.NewInt(10_000_000), 5000)
	require.Error(t, err)
	require.Contains(t, err.Error(), "maximum follows reached")

	// Verify exactly 20 follows
	following := k.GetFollowing(ctx, follower)
	require.Len(t, following, keeper.MaxFollowsPerUser)
}

// ---------------------------------------------------------------------------
// TestUnfollowTrader
// ---------------------------------------------------------------------------

func TestUnfollowTrader(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	follower := testAddr()
	trader := testAddr2()

	err := k.FollowTrader(ctx, follower, trader, math.NewInt(500_000), math.NewInt(10_000_000), 5000)
	require.NoError(t, err)

	// Verify follow exists
	_, exists := k.GetCopySettings(ctx, follower, trader)
	require.True(t, exists)

	// Unfollow
	err = k.UnfollowTrader(ctx, follower, trader)
	require.NoError(t, err)

	// Verify removed
	_, exists = k.GetCopySettings(ctx, follower, trader)
	require.False(t, exists)

	// Verify follower count decremented
	stats, found := k.GetTraderStats(ctx, trader)
	require.True(t, found)
	require.Equal(t, int64(0), stats.FollowerCount)

	// Verify no longer in followers list
	followers := k.GetFollowers(ctx, trader)
	require.Len(t, followers, 0)
}

// ---------------------------------------------------------------------------
// TestUnfollowTraderNotFollowing
// ---------------------------------------------------------------------------

func TestUnfollowTraderNotFollowing(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	follower := testAddr()
	trader := testAddr2()

	err := k.UnfollowTrader(ctx, follower, trader)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not following")
}

// ---------------------------------------------------------------------------
// TestUpdateCopySettings
// ---------------------------------------------------------------------------

func TestUpdateCopySettings(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	follower := testAddr()
	trader := testAddr2()

	err := k.FollowTrader(ctx, follower, trader, math.NewInt(500_000), math.NewInt(10_000_000), 5000)
	require.NoError(t, err)

	// Update settings
	err = k.UpdateCopySettings(ctx, follower, trader, math.NewInt(1_000_000), math.NewInt(20_000_000), 7500)
	require.NoError(t, err)

	// Verify changes
	settings, exists := k.GetCopySettings(ctx, follower, trader)
	require.True(t, exists)
	require.Equal(t, math.NewInt(1_000_000), settings.MaxPerTrade)
	require.Equal(t, math.NewInt(20_000_000), settings.TotalBudget)
	require.Equal(t, int64(7500), settings.CopyRatio)

	// Verify spent was not reset
	require.True(t, settings.Spent.IsZero())
}

// ---------------------------------------------------------------------------
// TestCopyTradeExecution
// ---------------------------------------------------------------------------

func TestCopyTradeExecution(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	poolID, _ := setupPoolForCopyTrade(t, k, ctx, bk)

	trader := testAddr2()
	follower := testAddr3()

	// Follower follows trader with 50% copy ratio
	err := k.FollowTrader(ctx, follower, trader, math.NewInt(1_000_000), math.NewInt(10_000_000), 5000)
	require.NoError(t, err)

	// Fund trader for the swap
	bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1_000_000)))

	// Fund follower for the copy trade (50% of 1M = 500k)
	bk.fundAccount(follower, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1_000_000)))

	// Record follower balance before
	followerUusdcBefore := bk.getBalance(follower, "uusdc")

	// Trader swaps (this automatically calls RecordTradeForCopyTrading)
	tokenOut, err := k.Swap(ctx, trader, poolID, sdk.NewInt64Coin("usyreen", 1_000_000), math.ZeroInt())
	require.NoError(t, err)
	require.True(t, tokenOut.Amount.IsPositive())

	// Advance block height past cooldown (context starts at 100, set last_copy_at=0 so no cooldown issue)
	// Process copy trades (EndBlock equivalent)
	k.ProcessCopyTrades(ctx)

	// Verify follower received uusdc tokens
	followerUusdcAfter := bk.getBalance(follower, "uusdc")
	require.True(t, followerUusdcAfter.GT(followerUusdcBefore), "follower should have received uusdc from copy trade")

	// Verify follower spent some usyreen
	followerUsyreenAfter := bk.getBalance(follower, "usyreen")
	require.True(t, followerUsyreenAfter.LT(math.NewInt(1_000_000)), "follower should have spent usyreen")

	// Verify settings updated (spent, lastCopyAt)
	settings, exists := k.GetCopySettings(ctx, follower, trader)
	require.True(t, exists)
	require.True(t, settings.Spent.IsPositive())
	require.Equal(t, int64(100), settings.LastCopyAt)
}

// ---------------------------------------------------------------------------
// TestCopyTradeBudgetLimit
// ---------------------------------------------------------------------------

func TestCopyTradeBudgetLimit(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	poolID, _ := setupPoolForCopyTrade(t, k, ctx, bk)

	trader := testAddr2()
	follower := testAddr3()

	// Follow with a very small budget (100 tokens) at 100% ratio
	err := k.FollowTrader(ctx, follower, trader, math.NewInt(1_000_000), math.NewInt(100), 10000)
	require.NoError(t, err)

	// Fund accounts
	bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1_000_000)))
	bk.fundAccount(follower, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1_000_000)))

	// Trader swaps 1M (copy amount would be 1M but budget is only 100)
	_, err = k.Swap(ctx, trader, poolID, sdk.NewInt64Coin("usyreen", 1_000_000), math.ZeroInt())
	require.NoError(t, err)

	k.ProcessCopyTrades(ctx)

	// After first copy trade, budget should be spent (100 tokens)
	settings, exists := k.GetCopySettings(ctx, follower, trader)
	require.True(t, exists)
	require.Equal(t, math.NewInt(100), settings.Spent)

	// Advance context to new block to bypass cooldown
	ctx = ctx.WithBlockHeader(cmtproto.Header{Height: 103})

	// Second swap: budget should be exhausted
	bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1_000_000)))
	_, err = k.Swap(ctx, trader, poolID, sdk.NewInt64Coin("usyreen", 1_000_000), math.ZeroInt())
	require.NoError(t, err)

	k.ProcessCopyTrades(ctx)

	// Settings should now be deactivated
	settings, exists = k.GetCopySettings(ctx, follower, trader)
	require.True(t, exists)
	require.False(t, settings.Active, "copy should be deactivated after budget exhausted")
}

// ---------------------------------------------------------------------------
// TestCopyTradeCooldown
// ---------------------------------------------------------------------------

func TestCopyTradeCooldown(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	poolID, _ := setupPoolForCopyTrade(t, k, ctx, bk)

	trader := testAddr2()
	follower := testAddr3()

	err := k.FollowTrader(ctx, follower, trader, math.NewInt(1_000_000), math.NewInt(50_000_000), 5000)
	require.NoError(t, err)

	// Fund accounts generously
	bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10_000_000)))
	bk.fundAccount(follower, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10_000_000)))

	// First trade at block 100
	_, err = k.Swap(ctx, trader, poolID, sdk.NewInt64Coin("usyreen", 1_000_000), math.ZeroInt())
	require.NoError(t, err)
	k.ProcessCopyTrades(ctx)

	// Verify first copy executed
	settings, _ := k.GetCopySettings(ctx, follower, trader)
	firstSpent := settings.Spent
	require.True(t, firstSpent.IsPositive(), "first copy should have executed")

	// Second trade at block 101 (only 1 block later, cooldown is 2)
	ctx = ctx.WithBlockHeader(cmtproto.Header{Height: 101})
	bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1_000_000)))
	_, err = k.Swap(ctx, trader, poolID, sdk.NewInt64Coin("usyreen", 1_000_000), math.ZeroInt())
	require.NoError(t, err)
	k.ProcessCopyTrades(ctx)

	// Spent should NOT have increased (cooldown prevents copy)
	settings, _ = k.GetCopySettings(ctx, follower, trader)
	require.Equal(t, firstSpent, settings.Spent, "copy should be blocked by cooldown")

	// Third trade at block 102 (2 blocks later, cooldown satisfied)
	ctx = ctx.WithBlockHeader(cmtproto.Header{Height: 102})
	bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1_000_000)))
	_, err = k.Swap(ctx, trader, poolID, sdk.NewInt64Coin("usyreen", 1_000_000), math.ZeroInt())
	require.NoError(t, err)
	k.ProcessCopyTrades(ctx)

	// Now spent should have increased
	settings, _ = k.GetCopySettings(ctx, follower, trader)
	require.True(t, settings.Spent.GT(firstSpent), "copy should execute after cooldown expires")
}

// ---------------------------------------------------------------------------
// TestCopyTradeMaxPerTrade
// ---------------------------------------------------------------------------

func TestCopyTradeMaxPerTrade(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	poolID, _ := setupPoolForCopyTrade(t, k, ctx, bk)

	trader := testAddr2()
	follower := testAddr3()

	// Follow with 100% ratio but max 200 per trade
	err := k.FollowTrader(ctx, follower, trader, math.NewInt(200), math.NewInt(50_000_000), 10000)
	require.NoError(t, err)

	// Fund accounts
	bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1_000_000)))
	bk.fundAccount(follower, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1_000_000)))

	// Trader swaps 1M (copy would be 1M at 100% but capped at 200)
	_, err = k.Swap(ctx, trader, poolID, sdk.NewInt64Coin("usyreen", 1_000_000), math.ZeroInt())
	require.NoError(t, err)
	k.ProcessCopyTrades(ctx)

	// Verify spent is exactly 200 (the max per trade cap)
	settings, exists := k.GetCopySettings(ctx, follower, trader)
	require.True(t, exists)
	require.Equal(t, math.NewInt(200), settings.Spent, "copy amount should be capped at max_per_trade")
}

// ---------------------------------------------------------------------------
// TestGetLeaderboard
// ---------------------------------------------------------------------------

func TestGetLeaderboard(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	poolID, _ := setupPoolForCopyTrade(t, k, ctx, bk)

	// Three traders swap different amounts to create different PnL
	trader1 := testAddr2()
	trader2 := testAddr3()
	trader3 := testAddr4()

	// Trader 1: large swap (more negative PnL due to price impact)
	bk.fundAccount(trader1, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10_000_000)))
	_, err := k.Swap(ctx, trader1, poolID, sdk.NewInt64Coin("usyreen", 10_000_000), math.ZeroInt())
	require.NoError(t, err)

	// Trader 2: small swap (less negative PnL)
	bk.fundAccount(trader2, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100_000)))
	_, err = k.Swap(ctx, trader2, poolID, sdk.NewInt64Coin("usyreen", 100_000), math.ZeroInt())
	require.NoError(t, err)

	// Trader 3: medium swap
	bk.fundAccount(trader3, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1_000_000)))
	_, err = k.Swap(ctx, trader3, poolID, sdk.NewInt64Coin("usyreen", 1_000_000), math.ZeroInt())
	require.NoError(t, err)

	// Get leaderboard
	leaderboard := k.GetLeaderboard(ctx, 10)
	require.Len(t, leaderboard, 3)

	// Should be sorted by PnL descending (least negative first)
	for i := 0; i < len(leaderboard)-1; i++ {
		require.True(t, leaderboard[i].TotalPnL.GTE(leaderboard[i+1].TotalPnL),
			"leaderboard should be sorted by PnL descending: index %d (%s) >= index %d (%s)",
			i, leaderboard[i].TotalPnL, i+1, leaderboard[i+1].TotalPnL)
	}

	// Test limit
	limited := k.GetLeaderboard(ctx, 2)
	require.Len(t, limited, 2)
}

// ---------------------------------------------------------------------------
// TestCopyTradeLog
// ---------------------------------------------------------------------------

func TestCopyTradeLog(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	poolID, _ := setupPoolForCopyTrade(t, k, ctx, bk)

	trader := testAddr2()
	follower := testAddr3()

	err := k.FollowTrader(ctx, follower, trader, math.NewInt(1_000_000), math.NewInt(50_000_000), 5000)
	require.NoError(t, err)

	// Fund accounts
	bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1_000_000)))
	bk.fundAccount(follower, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1_000_000)))

	// Trader swaps
	_, err = k.Swap(ctx, trader, poolID, sdk.NewInt64Coin("usyreen", 1_000_000), math.ZeroInt())
	require.NoError(t, err)
	k.ProcessCopyTrades(ctx)

	// Check log entries
	logs := k.GetCopyTradeLog(ctx, follower)
	require.Len(t, logs, 1)
	require.Equal(t, follower, logs[0].Follower)
	require.Equal(t, trader, logs[0].Trader)
	require.Equal(t, poolID, logs[0].PoolID)
	require.True(t, logs[0].Success)
	require.Empty(t, logs[0].FailReason)
	require.Equal(t, "usyreen", logs[0].TokenIn.Denom)
	require.True(t, logs[0].TokenIn.Amount.IsPositive())
	require.Equal(t, "uusdc", logs[0].TokenOut.Denom)
	require.True(t, logs[0].TokenOut.Amount.IsPositive())
	require.Equal(t, int64(100), logs[0].BlockHeight)
}

// ---------------------------------------------------------------------------
// TestAutoUnfollow
// ---------------------------------------------------------------------------

func TestAutoUnfollow(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)

	// Create a heavily imbalanced pool so a single big swap drives PnL well
	// past -50% (the M10-tightened auto-unfollow threshold).
	creator := testAddr()
	bk.fundAccount(creator, sdk.NewCoins(
		sdk.NewInt64Coin("usyreen", 100_000_000),
		sdk.NewInt64Coin("uusdc", 10_000_000),
	))
	poolID, err := k.CreatePool(ctx, creator, "usyreen", "uusdc", math.NewInt(100_000_000), math.NewInt(10_000_000))
	require.NoError(t, err)

	bk.balances[types.ModuleName] = bk.balances[types.ModuleName].Add(
		sdk.NewCoin("usyreen", math.NewInt(100_000_000)),
		sdk.NewCoin("uusdc", math.NewInt(10_000_000)),
	)

	trader := testAddr2()
	follower := testAddr3()
	follower2 := testAddr4()

	// Two followers follow the trader
	err = k.FollowTrader(ctx, follower, trader, math.NewInt(1_000_000), math.NewInt(50_000_000), 5000)
	require.NoError(t, err)
	err = k.FollowTrader(ctx, follower2, trader, math.NewInt(1_000_000), math.NewInt(50_000_000), 5000)
	require.NoError(t, err)

	// Verify 2 followers
	followers := k.GetFollowers(ctx, trader)
	require.Len(t, followers, 2)

	// Swap 100M usyreen into a 100M/10M pool. Output is ~5M uusdc, so PnL
	// (in input units, normalized via the pool spot price) is deeply
	// negative — well past -50% in bps terms.
	bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100_000_000)))
	_, err = k.Swap(ctx, trader, poolID, sdk.NewInt64Coin("usyreen", 100_000_000), math.ZeroInt())
	require.NoError(t, err)

	// Check trader stats PnL is deeply negative
	stats, found := k.GetTraderStats(ctx, trader)
	require.True(t, found)
	require.True(t, stats.TotalPnL.IsNegative(), "PnL should be negative after large swap")

	// PnL in bps: (pnl / volume) * 10000
	pnlBps := stats.TotalPnL.MulRaw(10000).Quo(stats.TotalVolume).Int64()
	require.True(t, pnlBps < int64(keeper.AutoUnfollowPnLThreshold),
		"PnL bps (%d) should be below threshold (%d)", pnlBps, keeper.AutoUnfollowPnLThreshold)

	// Process copy trades triggers checkAutoUnfollow
	// Fund followers so copy trades don't fail due to insufficient funds
	bk.fundAccount(follower, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 50_000_000)))
	bk.fundAccount(follower2, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 50_000_000)))
	k.ProcessCopyTrades(ctx)

	// After auto-unfollow, all followers should be removed
	followers = k.GetFollowers(ctx, trader)
	require.Len(t, followers, 0, "all followers should be auto-unfollowed")

	// Copy settings should be deleted
	_, exists := k.GetCopySettings(ctx, follower, trader)
	require.False(t, exists, "follower copy settings should be deleted")
	_, exists = k.GetCopySettings(ctx, follower2, trader)
	require.False(t, exists, "follower2 copy settings should be deleted")

	// Trader stats follower count should be 0
	stats, _ = k.GetTraderStats(ctx, trader)
	require.Equal(t, int64(0), stats.FollowerCount)
}
