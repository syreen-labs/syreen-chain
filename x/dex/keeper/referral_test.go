package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"syreen/x/dex/keeper"
	"syreen/x/dex/types"
)

func testAddr3() string {
	return sdk.AccAddress([]byte("third__addr_padding_")).String()
}

func testAddr4() string {
	return sdk.AccAddress([]byte("fourth_addr_padding_")).String()
}

// ---------------------------------------------------------------------------
// TestCreateReferralCode
// ---------------------------------------------------------------------------

func TestCreateReferralCode(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	creator := testAddr()

	code, err := k.CreateReferralCode(ctx, creator, "")
	require.NoError(t, err)
	require.NotEmpty(t, code)
	require.Len(t, code, 8) // 4 bytes hex = 8 chars

	// Verify code maps to address
	addr, found := k.GetReferrerByCode(ctx, code)
	require.True(t, found)
	require.Equal(t, creator, addr)

	// Verify address maps to code
	gotCode, found := k.GetReferralCodeByAddress(ctx, creator)
	require.True(t, found)
	require.Equal(t, code, gotCode)

	// Verify stats initialized
	stats, found := k.GetReferrerStats(ctx, creator)
	require.True(t, found)
	require.Equal(t, int64(0), stats.TotalReferrals)
	require.True(t, stats.TotalVolume.IsZero())
	require.True(t, stats.TotalEarned.IsZero())
	require.Equal(t, keeper.TierBronze, stats.Tier)
}

// ---------------------------------------------------------------------------
// TestCreateReferralCodeDuplicate
// ---------------------------------------------------------------------------

func TestCreateReferralCodeDuplicate(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	creator := testAddr()

	_, err := k.CreateReferralCode(ctx, creator, "")
	require.NoError(t, err)

	// Second attempt should fail
	_, err = k.CreateReferralCode(ctx, creator, "")
	require.ErrorIs(t, err, types.ErrReferralCodeExists)
}

// ---------------------------------------------------------------------------
// TestRegisterReferral
// ---------------------------------------------------------------------------

func TestRegisterReferral(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	referrer := testAddr()
	user := testAddr2()

	// Create referral code
	code, err := k.CreateReferralCode(ctx, referrer, "")
	require.NoError(t, err)

	// Register user with referral code
	gotReferrer, err := k.RegisterReferral(ctx, user, code)
	require.NoError(t, err)
	require.Equal(t, referrer, gotReferrer)

	// Verify registration stored
	storedReferrer, found := k.GetReferrer(ctx, user)
	require.True(t, found)
	require.Equal(t, referrer, storedReferrer)

	// Verify referrer stats updated
	stats, found := k.GetReferrerStats(ctx, referrer)
	require.True(t, found)
	require.Equal(t, int64(1), stats.TotalReferrals)
	require.Equal(t, keeper.TierBronze, stats.Tier) // 1 referral = bronze

	// Verify user stats initialized
	userStats, found := k.GetReferredUserStats(ctx, user)
	require.True(t, found)
	require.Equal(t, referrer, userStats.Referrer)
	require.True(t, userStats.TotalVolume.IsZero())

	// Verify global stats
	global := k.GetReferralGlobalStats(ctx)
	require.Equal(t, int64(1), global.TotalReferrals)
}

// ---------------------------------------------------------------------------
// TestRegisterReferralSelfReferral
// ---------------------------------------------------------------------------

func TestRegisterReferralSelfReferral(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	user := testAddr()

	code, err := k.CreateReferralCode(ctx, user, "")
	require.NoError(t, err)

	_, err = k.RegisterReferral(ctx, user, code)
	require.ErrorIs(t, err, types.ErrSelfReferral)
}

// ---------------------------------------------------------------------------
// TestRegisterReferralAlreadyReferred
// ---------------------------------------------------------------------------

func TestRegisterReferralAlreadyReferred(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	referrer1 := testAddr()
	referrer2 := testAddr3()
	user := testAddr2()

	code1, _ := k.CreateReferralCode(ctx, referrer1, "")
	code2, _ := k.CreateReferralCode(ctx, referrer2, "")

	_, err := k.RegisterReferral(ctx, user, code1)
	require.NoError(t, err)

	// Second registration should fail
	_, err = k.RegisterReferral(ctx, user, code2)
	require.ErrorIs(t, err, types.ErrAlreadyReferred)
}

// ---------------------------------------------------------------------------
// TestRegisterReferralCircular
// ---------------------------------------------------------------------------

func TestRegisterReferralCircular(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	addrA := testAddr()
	addrB := testAddr2()

	codeA, _ := k.CreateReferralCode(ctx, addrA, "")
	codeB, _ := k.CreateReferralCode(ctx, addrB, "")

	// A refers B
	_, err := k.RegisterReferral(ctx, addrB, codeA)
	require.NoError(t, err)

	// B tries to refer A (circular)
	_, err = k.RegisterReferral(ctx, addrA, codeB)
	require.ErrorIs(t, err, types.ErrCircularReferral)
}

// ---------------------------------------------------------------------------
// TestRegisterReferralInvalidCode
// ---------------------------------------------------------------------------

func TestRegisterReferralInvalidCode(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	user := testAddr()

	_, err := k.RegisterReferral(ctx, user, "NONEXISTENT")
	require.ErrorIs(t, err, types.ErrReferralCodeNotFound)
}

// ---------------------------------------------------------------------------
// TestReferralTiers
// ---------------------------------------------------------------------------

func TestReferralTiers(t *testing.T) {
	require.Equal(t, keeper.TierBronze, keeper.GetTier(0))
	require.Equal(t, keeper.TierBronze, keeper.GetTier(10))
	require.Equal(t, keeper.TierSilver, keeper.GetTier(11))
	require.Equal(t, keeper.TierSilver, keeper.GetTier(50))
	require.Equal(t, keeper.TierGold, keeper.GetTier(51))
	require.Equal(t, keeper.TierGold, keeper.GetTier(200))
	require.Equal(t, keeper.TierPlatinum, keeper.GetTier(201))
}

// ---------------------------------------------------------------------------
// TestTierFeePercentages
// ---------------------------------------------------------------------------

func TestTierFeePercentages(t *testing.T) {
	// Bronze: 5% referrer, 5% discount
	require.Equal(t, int64(500), keeper.GetTierReferrerPctBps(keeper.TierBronze))
	require.Equal(t, int64(500), keeper.GetTierDiscountPctBps(keeper.TierBronze))

	// Silver: 7% referrer, 5% discount
	require.Equal(t, int64(700), keeper.GetTierReferrerPctBps(keeper.TierSilver))
	require.Equal(t, int64(500), keeper.GetTierDiscountPctBps(keeper.TierSilver))

	// Gold: 10% referrer, 7% discount
	require.Equal(t, int64(1000), keeper.GetTierReferrerPctBps(keeper.TierGold))
	require.Equal(t, int64(700), keeper.GetTierDiscountPctBps(keeper.TierGold))

	// Platinum: 12% referrer, 10% discount
	require.Equal(t, int64(1200), keeper.GetTierReferrerPctBps(keeper.TierPlatinum))
	require.Equal(t, int64(1000), keeper.GetTierDiscountPctBps(keeper.TierPlatinum))
}

// ---------------------------------------------------------------------------
// TestReferralCodeDeterministic
// ---------------------------------------------------------------------------

func TestReferralCodeDeterministic(t *testing.T) {
	addr := testAddr()
	code1 := keeper.GenerateReferralCode(addr)
	code2 := keeper.GenerateReferralCode(addr)
	require.Equal(t, code1, code2) // Same address always gets same code
}

// ---------------------------------------------------------------------------
// TestCalculateReferralFees
// ---------------------------------------------------------------------------

func TestCalculateReferralFees(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	referrer := testAddr()
	trader := testAddr2()

	// Setup: create referral and register
	code, _ := k.CreateReferralCode(ctx, referrer, "")
	k.RegisterReferral(ctx, trader, code)

	// Calculate fees for a 1,000,000 token swap with 30 bps fee
	result := k.CalculateReferralFees(ctx, trader, math.NewInt(1_000_000), 30)
	require.True(t, result.HasReferral)

	// Total fee = 1_000_000 * 30 / 10000 = 3000
	// Bronze tier: referrer 5% of fee = 3000 * 500 / 10000 = 150
	// Bronze tier: discount 5% of fee = 3000 * 500 / 10000 = 150
	require.Equal(t, math.NewInt(150), result.ReferrerAmount)
	require.Equal(t, math.NewInt(150), result.TraderDiscount)
	require.Equal(t, referrer, result.ReferrerAddrStr)
}

// ---------------------------------------------------------------------------
// TestCalculateReferralFeesNoReferral
// ---------------------------------------------------------------------------

func TestCalculateReferralFeesNoReferral(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	trader := testAddr()

	// No referral registered
	result := k.CalculateReferralFees(ctx, trader, math.NewInt(1_000_000), 30)
	require.False(t, result.HasReferral)
}

// ---------------------------------------------------------------------------
// TestSwapWithReferral — end-to-end integration test
// ---------------------------------------------------------------------------

func TestSwapWithReferral(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()
	referrer := testAddr3()
	trader := testAddr2()

	// Create pool: 10M usyreen / 10M uusdc (1:1)
	bk.fundAccount(creator, sdk.NewCoins(
		sdk.NewInt64Coin("usyreen", 10_000_000),
		sdk.NewInt64Coin("uusdc", 10_000_000),
	))
	poolID, err := k.CreatePool(ctx, creator, "usyreen", "uusdc", math.NewInt(10_000_000), math.NewInt(10_000_000))
	require.NoError(t, err)

	// Fund module with tokens for output
	bk.balances[types.ModuleName] = bk.balances[types.ModuleName].Add(
		sdk.NewCoin("usyreen", math.NewInt(10_000_000)),
		sdk.NewCoin("uusdc", math.NewInt(10_000_000)),
	)

	// Setup referral
	code, err := k.CreateReferralCode(ctx, referrer, "")
	require.NoError(t, err)
	_, err = k.RegisterReferral(ctx, trader, code)
	require.NoError(t, err)

	// H5: Referrer balance is no longer credited directly. The reward is held in
	// a claimable balance that the referrer must withdraw via ClaimReferralRewards.
	referrerBalBefore := bk.getBalance(referrer, "usyreen")

	// Trader swaps 1M usyreen for uusdc
	bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1_000_000)))
	tokenIn := sdk.NewInt64Coin("usyreen", 1_000_000)
	tokenOut, err := k.Swap(ctx, trader, poolID, tokenIn, math.ZeroInt())
	require.NoError(t, err)
	require.True(t, tokenOut.Amount.IsPositive())

	// Referrer's spendable balance should be unchanged (rewards are claimable, not paid).
	require.Equal(t, referrerBalBefore, bk.getBalance(referrer, "usyreen"),
		"referrer balance should not change before claiming")

	// Fee = 1_000_000 * 30 / 10000 = 3000
	// Bronze: referrer = 3000 * 500 / 10000 = 150
	claimable := k.GetClaimableReferralRewards(ctx, referrer, "usyreen")
	require.Equal(t, math.NewInt(150), claimable, "referrer should have 150 claimable")

	// Now claim and verify the funds land in the referrer's account.
	paid, err := k.ClaimReferralRewards(ctx, referrer)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(150), paid.AmountOf("usyreen"))
	referrerBalAfter := bk.getBalance(referrer, "usyreen")
	require.Equal(t, math.NewInt(150), referrerBalAfter.Sub(referrerBalBefore))

	// Claimable balance should now be drained.
	require.True(t, k.GetClaimableReferralRewards(ctx, referrer, "usyreen").IsZero())

	// Verify referrer stats updated
	stats, found := k.GetReferrerStats(ctx, referrer)
	require.True(t, found)
	require.True(t, stats.TotalEarned.IsPositive())
	require.Equal(t, math.NewInt(150), stats.TotalEarned)

	// Verify user stats updated
	userStats, found := k.GetReferredUserStats(ctx, trader)
	require.True(t, found)
	require.True(t, userStats.TotalSaved.IsPositive())
	require.Equal(t, math.NewInt(150), userStats.TotalSaved)

	// Verify earnings history
	earnings := k.GetReferralEarnings(ctx, referrer)
	require.Len(t, earnings.Entries, 1)
	require.Equal(t, math.NewInt(150), earnings.Entries[0].Amount)
	require.Equal(t, trader, earnings.Entries[0].Trader)
	require.Equal(t, "usyreen", earnings.Entries[0].Denom)
	require.Equal(t, poolID, earnings.Entries[0].PoolID)

	// Verify global stats
	global := k.GetReferralGlobalStats(ctx)
	require.True(t, global.TotalFeesDistributed.IsPositive())
}

// ---------------------------------------------------------------------------
// TestSwapWithoutReferral — verify no disruption when no referral
// ---------------------------------------------------------------------------

func TestSwapWithoutReferral(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()
	trader := testAddr2()

	bk.fundAccount(creator, sdk.NewCoins(
		sdk.NewInt64Coin("usyreen", 10_000_000),
		sdk.NewInt64Coin("uusdc", 10_000_000),
	))
	poolID, err := k.CreatePool(ctx, creator, "usyreen", "uusdc", math.NewInt(10_000_000), math.NewInt(10_000_000))
	require.NoError(t, err)

	bk.balances[types.ModuleName] = bk.balances[types.ModuleName].Add(
		sdk.NewCoin("usyreen", math.NewInt(10_000_000)),
		sdk.NewCoin("uusdc", math.NewInt(10_000_000)),
	)

	bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1_000_000)))
	tokenIn := sdk.NewInt64Coin("usyreen", 1_000_000)
	tokenOut, err := k.Swap(ctx, trader, poolID, tokenIn, math.ZeroInt())
	require.NoError(t, err)
	require.True(t, tokenOut.Amount.IsPositive())
	require.True(t, tokenOut.Amount.LT(math.NewInt(1_000_000)))
}

// ---------------------------------------------------------------------------
// TestMultipleSwapsWithReferral — verify accumulation
// ---------------------------------------------------------------------------

func TestMultipleSwapsWithReferral(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()
	referrer := testAddr3()
	trader := testAddr2()

	bk.fundAccount(creator, sdk.NewCoins(
		sdk.NewInt64Coin("usyreen", 100_000_000),
		sdk.NewInt64Coin("uusdc", 100_000_000),
	))
	poolID, err := k.CreatePool(ctx, creator, "usyreen", "uusdc", math.NewInt(100_000_000), math.NewInt(100_000_000))
	require.NoError(t, err)

	bk.balances[types.ModuleName] = bk.balances[types.ModuleName].Add(
		sdk.NewCoin("usyreen", math.NewInt(100_000_000)),
		sdk.NewCoin("uusdc", math.NewInt(100_000_000)),
	)

	code, _ := k.CreateReferralCode(ctx, referrer, "")
	k.RegisterReferral(ctx, trader, code)

	// Perform 5 swaps
	for i := 0; i < 5; i++ {
		bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1_000_000)))
		_, err := k.Swap(ctx, trader, poolID, sdk.NewInt64Coin("usyreen", 1_000_000), math.ZeroInt())
		require.NoError(t, err)
	}

	// Verify accumulated earnings
	stats, found := k.GetReferrerStats(ctx, referrer)
	require.True(t, found)
	require.True(t, stats.TotalEarned.GT(math.ZeroInt()))

	// Each swap: fee = 1M * 30 / 10000 = 3000, referrer gets 150
	// 5 swaps = ~750 (approximately, fees may vary slightly as reserves change)
	require.True(t, stats.TotalEarned.GTE(math.NewInt(700)), "should have accumulated earnings from 5 swaps")

	// Verify earnings history
	earnings := k.GetReferralEarnings(ctx, referrer)
	require.Len(t, earnings.Entries, 5)
}

// ---------------------------------------------------------------------------
// TestGetAllReferredUsers
// ---------------------------------------------------------------------------

func TestGetAllReferredUsers(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	referrer := testAddr()
	user1 := testAddr2()
	user2 := testAddr3()
	user3 := testAddr4()

	code, _ := k.CreateReferralCode(ctx, referrer, "")
	k.RegisterReferral(ctx, user1, code)
	k.RegisterReferral(ctx, user2, code)
	k.RegisterReferral(ctx, user3, code)

	users := k.GetAllReferredUsers(ctx, referrer)
	require.Len(t, users, 3)

	// Verify all users are present
	addresses := make(map[string]bool)
	for _, u := range users {
		addresses[u.Address] = true
	}
	require.True(t, addresses[user1])
	require.True(t, addresses[user2])
	require.True(t, addresses[user3])
}

// ---------------------------------------------------------------------------
// TestReferralGlobalStats
// ---------------------------------------------------------------------------

func TestReferralGlobalStats(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	// Default global stats
	global := k.GetReferralGlobalStats(ctx)
	require.Equal(t, int64(0), global.TotalReferrals)
	require.True(t, global.TotalFeesDistributed.IsZero())

	// Create referrals
	referrer := testAddr()
	code, _ := k.CreateReferralCode(ctx, referrer, "")
	k.RegisterReferral(ctx, testAddr2(), code)
	k.RegisterReferral(ctx, testAddr3(), code)

	global = k.GetReferralGlobalStats(ctx)
	require.Equal(t, int64(2), global.TotalReferrals)
}

// ---------------------------------------------------------------------------
// TestReferralPoolReserveAdjustment
// ---------------------------------------------------------------------------

func TestReferralPoolReserveAdjustment(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()
	referrer := testAddr3()
	trader := testAddr2()

	bk.fundAccount(creator, sdk.NewCoins(
		sdk.NewInt64Coin("usyreen", 10_000_000),
		sdk.NewInt64Coin("uusdc", 10_000_000),
	))
	poolID, err := k.CreatePool(ctx, creator, "usyreen", "uusdc", math.NewInt(10_000_000), math.NewInt(10_000_000))
	require.NoError(t, err)

	bk.balances[types.ModuleName] = bk.balances[types.ModuleName].Add(
		sdk.NewCoin("usyreen", math.NewInt(10_000_000)),
		sdk.NewCoin("uusdc", math.NewInt(10_000_000)),
	)

	code, _ := k.CreateReferralCode(ctx, referrer, "")
	k.RegisterReferral(ctx, trader, code)

	poolBefore, _ := k.GetPool(ctx, poolID)

	bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1_000_000)))
	_, err = k.Swap(ctx, trader, poolID, sdk.NewInt64Coin("usyreen", 1_000_000), math.ZeroInt())
	require.NoError(t, err)

	poolAfter, _ := k.GetPool(ctx, poolID)

	// The pool's reserveA should be less than without referral because
	// referral fees (referrer + discount) are taken from the fee portion
	// Without referral: reserveA = 10M + 1M = 11M
	// With referral: reserveA = 11M - referralTotal (300 for bronze)
	expectedReserveA := poolBefore.ReserveA.Add(math.NewInt(1_000_000)).Sub(math.NewInt(300))
	require.Equal(t, expectedReserveA, poolAfter.ReserveA, "reserves should account for referral payouts")
}

// ---------------------------------------------------------------------------
// TestEarningsHistoryLimit
// ---------------------------------------------------------------------------

func TestEarningsHistoryLimit(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()
	referrer := testAddr3()
	trader := testAddr2()

	bk.fundAccount(creator, sdk.NewCoins(
		sdk.NewInt64Coin("usyreen", 1_000_000_000),
		sdk.NewInt64Coin("uusdc", 1_000_000_000),
	))
	poolID, err := k.CreatePool(ctx, creator, "usyreen", "uusdc", math.NewInt(1_000_000_000), math.NewInt(1_000_000_000))
	require.NoError(t, err)

	bk.balances[types.ModuleName] = bk.balances[types.ModuleName].Add(
		sdk.NewCoin("usyreen", math.NewInt(1_000_000_000)),
		sdk.NewCoin("uusdc", math.NewInt(1_000_000_000)),
	)

	code, _ := k.CreateReferralCode(ctx, referrer, "")
	k.RegisterReferral(ctx, trader, code)

	// Do 110 swaps to exceed the 100 entry limit
	for i := 0; i < 110; i++ {
		bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100_000)))
		_, err := k.Swap(ctx, trader, poolID, sdk.NewInt64Coin("usyreen", 100_000), math.ZeroInt())
		require.NoError(t, err, "swap %d", i)
	}

	earnings := k.GetReferralEarnings(ctx, referrer)
	require.Len(t, earnings.Entries, 100, "should cap at 100 entries")

	_ = poolID
}

// ---------------------------------------------------------------------------
// TestReferralFeeCalculationSilverTier
// ---------------------------------------------------------------------------

func TestReferralFeeCalculationSilverTier(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	referrer := testAddr()
	trader := testAddr2()

	code, _ := k.CreateReferralCode(ctx, referrer, "")

	// Register 11 users to get Silver tier
	for i := 0; i < 11; i++ {
		userAddr := sdk.AccAddress([]byte(string(rune('A'+i)) + "_referral_user_pad")).String()
		k.RegisterReferral(ctx, userAddr, code)
	}

	stats, _ := k.GetReferrerStats(ctx, referrer)
	require.Equal(t, keeper.TierSilver, stats.Tier)

	// Register the trader too
	k.RegisterReferral(ctx, trader, code)

	// Calculate fees: Silver = 7% referrer, 5% discount
	result := k.CalculateReferralFees(ctx, trader, math.NewInt(1_000_000), 30)
	require.True(t, result.HasReferral)

	// Total fee = 3000
	// Silver: referrer = 3000 * 700 / 10000 = 210
	// Silver: discount = 3000 * 500 / 10000 = 150
	require.Equal(t, math.NewInt(210), result.ReferrerAmount)
	require.Equal(t, math.NewInt(150), result.TraderDiscount)
}
