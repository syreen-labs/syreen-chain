package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	"syreen/x/clmm/keeper"
	"syreen/x/clmm/types"

	dbm "github.com/cosmos/cosmos-db"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
)

// ============================================================
// Mock Keepers
// ============================================================

type mockAccountKeeper struct {
	moduleAddrs map[string]sdk.AccAddress
}

func (m mockAccountKeeper) GetModuleAddress(name string) sdk.AccAddress {
	if addr, ok := m.moduleAddrs[name]; ok {
		return addr
	}
	return authtypes.NewModuleAddress(name)
}

func (m mockAccountKeeper) GetModuleAccount(_ context.Context, name string) sdk.ModuleAccountI {
	return nil
}

type mockBankKeeper struct {
	balances map[string]map[string]math.Int // addr -> denom -> amount
}

func newMockBankKeeper() *mockBankKeeper {
	return &mockBankKeeper{balances: make(map[string]map[string]math.Int)}
}

func (m *mockBankKeeper) setBalance(addr, denom string, amount math.Int) {
	if m.balances[addr] == nil {
		m.balances[addr] = make(map[string]math.Int)
	}
	m.balances[addr][denom] = amount
}

func (m *mockBankKeeper) getBalance(addr, denom string) math.Int {
	if m.balances[addr] == nil {
		return math.ZeroInt()
	}
	if amt, ok := m.balances[addr][denom]; ok {
		return amt
	}
	return math.ZeroInt()
}

func (m *mockBankKeeper) GetBalance(_ context.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	amt := m.getBalance(addr.String(), denom)
	return sdk.NewCoin(denom, amt)
}

func (m *mockBankKeeper) SendCoins(_ context.Context, from, to sdk.AccAddress, amt sdk.Coins) error {
	for _, coin := range amt {
		fromBal := m.getBalance(from.String(), coin.Denom)
		if fromBal.LT(coin.Amount) {
			return types.ErrInsufficientFunds
		}
		m.setBalance(from.String(), coin.Denom, fromBal.Sub(coin.Amount))
		toBal := m.getBalance(to.String(), coin.Denom)
		m.setBalance(to.String(), coin.Denom, toBal.Add(coin.Amount))
	}
	return nil
}

func (m *mockBankKeeper) SendCoinsFromAccountToModule(_ context.Context, sender sdk.AccAddress, module string, amt sdk.Coins) error {
	moduleAddr := authtypes.NewModuleAddress(module)
	for _, coin := range amt {
		fromBal := m.getBalance(sender.String(), coin.Denom)
		if fromBal.LT(coin.Amount) {
			return types.ErrInsufficientFunds
		}
		m.setBalance(sender.String(), coin.Denom, fromBal.Sub(coin.Amount))
		toBal := m.getBalance(moduleAddr.String(), coin.Denom)
		m.setBalance(moduleAddr.String(), coin.Denom, toBal.Add(coin.Amount))
	}
	return nil
}

func (m *mockBankKeeper) SendCoinsFromModuleToAccount(_ context.Context, module string, recipient sdk.AccAddress, amt sdk.Coins) error {
	moduleAddr := authtypes.NewModuleAddress(module)
	for _, coin := range amt {
		fromBal := m.getBalance(moduleAddr.String(), coin.Denom)
		if fromBal.LT(coin.Amount) {
			return types.ErrInsufficientFunds
		}
		m.setBalance(moduleAddr.String(), coin.Denom, fromBal.Sub(coin.Amount))
		toBal := m.getBalance(recipient.String(), coin.Denom)
		m.setBalance(recipient.String(), coin.Denom, toBal.Add(coin.Amount))
	}
	return nil
}

// ============================================================
// Test Setup
// ============================================================

func setupKeeper(t *testing.T) (keeper.Keeper, sdk.Context, *mockBankKeeper) {
	t.Helper()
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, cmtproto.Header{Height: 1}, false, log.NewNopLogger())

	_ = codectypes.NewInterfaceRegistry()

	accountKeeper := mockAccountKeeper{
		moduleAddrs: map[string]sdk.AccAddress{
			types.ModuleName: authtypes.NewModuleAddress(types.ModuleName),
		},
	}
	bankKeeper := newMockBankKeeper()

	k := keeper.NewKeeper(
		runtime.NewKVStoreService(storeKey),
		accountKeeper,
		bankKeeper,
		authtypes.NewModuleAddress("gov").String(),
	)

	return k, ctx, bankKeeper
}

func testAddr(i int) string {
	return sdk.AccAddress([]byte{byte(i), 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}).String()
}

const denomA = "tokenA"
const denomB = "tokenB"

func fundUser(bk *mockBankKeeper, userIdx int, denom string, amount math.Int) string {
	addr := testAddr(userIdx)
	bk.setBalance(addr, denom, amount)
	return addr
}

func fundModule(bk *mockBankKeeper, denom string, amount math.Int) {
	moduleAddr := authtypes.NewModuleAddress(types.ModuleName).String()
	bk.setBalance(moduleAddr, denom, amount)
}

// createTestPool creates a CL pool with price=1.0, tick spacing=10, fee=0.3%
func createTestPool(t *testing.T, k keeper.Keeper, ctx sdk.Context) uint64 {
	t.Helper()
	poolID, err := k.CreatePool(ctx, testAddr(99), denomA, denomB, 10, math.LegacyNewDecWithPrec(3, 3), math.LegacyOneDec())
	require.NoError(t, err)
	return poolID
}

// ============================================================
// Tests
// ============================================================

func TestCreatePool(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	poolID, err := k.CreatePool(ctx, testAddr(1), denomA, denomB, 10, math.LegacyNewDecWithPrec(3, 3), math.LegacyOneDec())
	require.NoError(t, err)
	require.Equal(t, uint64(1), poolID)

	pool, found := k.GetPool(ctx, poolID)
	require.True(t, found)
	require.Equal(t, denomA, pool.DenomA)
	require.Equal(t, denomB, pool.DenomB)
	require.Equal(t, int64(10), pool.TickSpacing)
	require.True(t, pool.FeeRate.Equal(math.LegacyNewDecWithPrec(3, 3)))
	require.True(t, pool.TotalLiquidity.IsZero())
}

func TestCreatePoolSortsDenoms(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	// Pass denoms in reverse order
	poolID, err := k.CreatePool(ctx, testAddr(1), "zzz", "aaa", 10, math.LegacyNewDecWithPrec(3, 3), math.LegacyOneDec())
	require.NoError(t, err)

	pool, found := k.GetPool(ctx, poolID)
	require.True(t, found)
	require.Equal(t, "aaa", pool.DenomA)
	require.Equal(t, "zzz", pool.DenomB)
}

func TestCreatePositionInRange(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	poolID := createTestPool(t, k, ctx)

	user := fundUser(bk, 1, denomA, math.NewInt(1_000_000))
	fundUser(bk, 1, denomB, math.NewInt(1_000_000))

	// Create position at tick range [-100, 100] around price=1.0 (tick=0)
	posID, amt0, amt1, liq, err := k.ExecuteCreatePosition(
		ctx, user, poolID,
		-100, 100,
		math.NewInt(500_000), math.NewInt(500_000),
		math.ZeroInt(), math.ZeroInt(),
	)
	require.NoError(t, err)
	require.Equal(t, uint64(1), posID)
	require.True(t, amt0.IsPositive(), "should deposit some token0")
	require.True(t, amt1.IsPositive(), "should deposit some token1")
	require.True(t, liq.IsPositive(), "should have positive liquidity")

	// Pool should now have liquidity
	pool, found := k.GetPool(ctx, poolID)
	require.True(t, found)
	require.True(t, pool.TotalLiquidity.IsPositive())
}

func TestCreatePositionBelowRange(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	poolID := createTestPool(t, k, ctx)

	user := fundUser(bk, 1, denomA, math.NewInt(1_000_000))
	fundUser(bk, 1, denomB, math.NewInt(1_000_000))

	// Create position entirely below current price (tick=0): range [-200, -100]
	// Current price is ABOVE this range, so only token1 is deposited (position fully in token1)
	posID, amt0, amt1, liq, err := k.ExecuteCreatePosition(
		ctx, user, poolID,
		-200, -100,
		math.NewInt(500_000), math.NewInt(500_000),
		math.ZeroInt(), math.ZeroInt(),
	)
	require.NoError(t, err)
	require.Equal(t, uint64(1), posID)
	require.True(t, amt0.IsZero(), "should NOT deposit token0 for below-range position")
	require.True(t, amt1.IsPositive(), "should deposit token1 for below-range position")
	require.True(t, liq.IsPositive())

	// Pool should NOT have active liquidity (position is out of range)
	pool, found := k.GetPool(ctx, poolID)
	require.True(t, found)
	require.True(t, pool.TotalLiquidity.IsZero(), "below-range position should not add to active liquidity")
}

func TestCreatePositionAboveRange(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	poolID := createTestPool(t, k, ctx)

	user := fundUser(bk, 1, denomA, math.NewInt(1_000_000))
	fundUser(bk, 1, denomB, math.NewInt(1_000_000))

	// Create position entirely above current price (tick=0): range [100, 200]
	// Current price is BELOW this range, so only token0 is deposited (position fully in token0)
	posID, amt0, amt1, liq, err := k.ExecuteCreatePosition(
		ctx, user, poolID,
		100, 200,
		math.NewInt(500_000), math.NewInt(500_000),
		math.ZeroInt(), math.ZeroInt(),
	)
	require.NoError(t, err)
	require.Equal(t, uint64(1), posID)
	require.True(t, amt0.IsPositive(), "should deposit token0 for above-range position")
	require.True(t, amt1.IsZero(), "should NOT deposit token1 for above-range position")
	require.True(t, liq.IsPositive())

	// Pool should NOT have active liquidity
	pool, found := k.GetPool(ctx, poolID)
	require.True(t, found)
	require.True(t, pool.TotalLiquidity.IsZero())
}

func TestCreatePositionInvalidTicks(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	poolID := createTestPool(t, k, ctx)
	user := fundUser(bk, 1, denomA, math.NewInt(1_000_000))
	fundUser(bk, 1, denomB, math.NewInt(1_000_000))

	// tickLower >= tickUpper
	_, _, _, _, err := k.ExecuteCreatePosition(ctx, user, poolID, 100, 100, math.NewInt(100), math.NewInt(100), math.ZeroInt(), math.ZeroInt())
	require.Error(t, err)

	// Ticks not aligned with spacing=10
	_, _, _, _, err = k.ExecuteCreatePosition(ctx, user, poolID, -5, 105, math.NewInt(100), math.NewInt(100), math.ZeroInt(), math.ZeroInt())
	require.Error(t, err)
}

func TestCreatePositionPoolNotFound(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	user := fundUser(bk, 1, denomA, math.NewInt(1_000_000))
	fundUser(bk, 1, denomB, math.NewInt(1_000_000))

	_, _, _, _, err := k.ExecuteCreatePosition(ctx, user, 999, -100, 100, math.NewInt(100), math.NewInt(100), math.ZeroInt(), math.ZeroInt())
	require.ErrorIs(t, err, types.ErrPoolNotFound)
}

func TestAddLiquidity(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	poolID := createTestPool(t, k, ctx)
	user := fundUser(bk, 1, denomA, math.NewInt(2_000_000))
	fundUser(bk, 1, denomB, math.NewInt(2_000_000))

	posID, _, _, liq1, err := k.ExecuteCreatePosition(ctx, user, poolID, -100, 100, math.NewInt(500_000), math.NewInt(500_000), math.ZeroInt(), math.ZeroInt())
	require.NoError(t, err)

	// Add more liquidity
	_, _, liq2, err := k.AddLiquidityToPosition(ctx, user, posID, math.NewInt(500_000), math.NewInt(500_000))
	require.NoError(t, err)
	require.True(t, liq2.IsPositive())

	pos, found := k.GetPosition(ctx, posID)
	require.True(t, found)
	require.True(t, pos.Liquidity.GT(liq1), "total liquidity should increase")
}

func TestAddLiquidityUnauthorized(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	poolID := createTestPool(t, k, ctx)
	user1 := fundUser(bk, 1, denomA, math.NewInt(1_000_000))
	fundUser(bk, 1, denomB, math.NewInt(1_000_000))
	user2 := fundUser(bk, 2, denomA, math.NewInt(1_000_000))
	fundUser(bk, 2, denomB, math.NewInt(1_000_000))

	posID, _, _, _, err := k.ExecuteCreatePosition(ctx, user1, poolID, -100, 100, math.NewInt(500_000), math.NewInt(500_000), math.ZeroInt(), math.ZeroInt())
	require.NoError(t, err)

	// user2 trying to add to user1's position
	_, _, _, err = k.AddLiquidityToPosition(ctx, user2, posID, math.NewInt(100), math.NewInt(100))
	require.ErrorIs(t, err, types.ErrUnauthorized)
}

func TestRemoveLiquidity(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	poolID := createTestPool(t, k, ctx)
	user := fundUser(bk, 1, denomA, math.NewInt(1_000_000))
	fundUser(bk, 1, denomB, math.NewInt(1_000_000))

	posID, _, _, liq, err := k.ExecuteCreatePosition(ctx, user, poolID, -100, 100, math.NewInt(500_000), math.NewInt(500_000), math.ZeroInt(), math.ZeroInt())
	require.NoError(t, err)

	// Remove half the liquidity
	halfLiq := liq.Quo(math.LegacyNewDec(2))
	amt0, amt1, err := k.RemoveLiquidityFromPosition(ctx, user, posID, halfLiq)
	require.NoError(t, err)
	require.True(t, amt0.IsPositive() || amt1.IsPositive(), "should return some tokens")

	pos, found := k.GetPosition(ctx, posID)
	require.True(t, found)
	require.True(t, pos.Liquidity.GT(math.LegacyZeroDec()), "should still have remaining liquidity")
}

func TestRemoveAllLiquidity(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	poolID := createTestPool(t, k, ctx)
	user := fundUser(bk, 1, denomA, math.NewInt(1_000_000))
	fundUser(bk, 1, denomB, math.NewInt(1_000_000))

	posID, _, _, liq, err := k.ExecuteCreatePosition(ctx, user, poolID, -100, 100, math.NewInt(500_000), math.NewInt(500_000), math.ZeroInt(), math.ZeroInt())
	require.NoError(t, err)

	// Remove all liquidity
	_, _, err = k.RemoveLiquidityFromPosition(ctx, user, posID, liq)
	require.NoError(t, err)

	// Position should be deleted (no fees owed)
	_, found := k.GetPosition(ctx, posID)
	require.False(t, found, "position with zero liquidity and no fees should be deleted")
}

func TestRemoveLiquidityUnauthorized(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	poolID := createTestPool(t, k, ctx)
	user1 := fundUser(bk, 1, denomA, math.NewInt(1_000_000))
	fundUser(bk, 1, denomB, math.NewInt(1_000_000))
	user2 := testAddr(2)

	posID, _, _, _, err := k.ExecuteCreatePosition(ctx, user1, poolID, -100, 100, math.NewInt(500_000), math.NewInt(500_000), math.ZeroInt(), math.ZeroInt())
	require.NoError(t, err)

	_, _, err = k.RemoveLiquidityFromPosition(ctx, user2, posID, math.LegacyOneDec())
	require.ErrorIs(t, err, types.ErrUnauthorized)
}

func TestSwapWithinRange(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	poolID := createTestPool(t, k, ctx)
	lp := fundUser(bk, 1, denomA, math.NewInt(10_000_000))
	fundUser(bk, 1, denomB, math.NewInt(10_000_000))

	// Create a wide position with lots of liquidity
	_, _, _, _, err := k.ExecuteCreatePosition(ctx, lp, poolID, -1000, 1000, math.NewInt(5_000_000), math.NewInt(5_000_000), math.ZeroInt(), math.ZeroInt())
	require.NoError(t, err)

	// Swap: sell token0 for token1
	trader := fundUser(bk, 2, denomA, math.NewInt(100_000))
	amountOut, err := k.ExecuteSwap(ctx, trader, poolID, denomA, math.NewInt(10_000), math.ZeroInt())
	require.NoError(t, err)
	require.True(t, amountOut.IsPositive(), "swap should produce output")

	// Verify trader's balances
	traderAddr, _ := sdk.AccAddressFromBech32(trader)
	balA := bk.GetBalance(ctx, traderAddr, denomA)
	require.True(t, balA.Amount.LT(math.NewInt(100_000)), "trader should have less tokenA")
	balB := bk.GetBalance(ctx, traderAddr, denomB)
	require.True(t, balB.Amount.IsPositive(), "trader should receive tokenB")
}

func TestSwapReversDirection(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	poolID := createTestPool(t, k, ctx)
	lp := fundUser(bk, 1, denomA, math.NewInt(10_000_000))
	fundUser(bk, 1, denomB, math.NewInt(10_000_000))

	_, _, _, _, err := k.ExecuteCreatePosition(ctx, lp, poolID, -1000, 1000, math.NewInt(5_000_000), math.NewInt(5_000_000), math.ZeroInt(), math.ZeroInt())
	require.NoError(t, err)

	// Swap: sell token1 for token0
	trader := fundUser(bk, 2, denomB, math.NewInt(100_000))
	amountOut, err := k.ExecuteSwap(ctx, trader, poolID, denomB, math.NewInt(10_000), math.ZeroInt())
	require.NoError(t, err)
	require.True(t, amountOut.IsPositive())
}

func TestSwapInvalidDenom(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	poolID := createTestPool(t, k, ctx)
	lp := fundUser(bk, 1, denomA, math.NewInt(10_000_000))
	fundUser(bk, 1, denomB, math.NewInt(10_000_000))
	_, _, _, _, err := k.ExecuteCreatePosition(ctx, lp, poolID, -1000, 1000, math.NewInt(5_000_000), math.NewInt(5_000_000), math.ZeroInt(), math.ZeroInt())
	require.NoError(t, err)

	trader := fundUser(bk, 2, "tokenX", math.NewInt(100_000))
	_, err = k.ExecuteSwap(ctx, trader, poolID, "tokenX", math.NewInt(10_000), math.ZeroInt())
	require.ErrorIs(t, err, types.ErrInvalidDenom)
}

func TestSwapSlippageProtection(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	poolID := createTestPool(t, k, ctx)
	lp := fundUser(bk, 1, denomA, math.NewInt(10_000_000))
	fundUser(bk, 1, denomB, math.NewInt(10_000_000))
	_, _, _, _, err := k.ExecuteCreatePosition(ctx, lp, poolID, -1000, 1000, math.NewInt(5_000_000), math.NewInt(5_000_000), math.ZeroInt(), math.ZeroInt())
	require.NoError(t, err)

	// Set unrealistically high minAmountOut
	trader := fundUser(bk, 2, denomA, math.NewInt(100_000))
	_, err = k.ExecuteSwap(ctx, trader, poolID, denomA, math.NewInt(100), math.NewInt(999_999_999))
	require.ErrorIs(t, err, types.ErrSlippageExceeded)
}

func TestSwapAcrossTicks(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	poolID := createTestPool(t, k, ctx)

	// Create two narrow positions at different ranges
	lp := fundUser(bk, 1, denomA, math.NewInt(50_000_000))
	fundUser(bk, 1, denomB, math.NewInt(50_000_000))

	// Position 1: tick [-100, 100]
	_, _, _, _, err := k.ExecuteCreatePosition(ctx, lp, poolID, -100, 100, math.NewInt(2_000_000), math.NewInt(2_000_000), math.ZeroInt(), math.ZeroInt())
	require.NoError(t, err)

	// Position 2: tick [-200, -100]  (below current price, has only tokenA)
	_, _, _, _, err = k.ExecuteCreatePosition(ctx, lp, poolID, -200, -100, math.NewInt(2_000_000), math.NewInt(2_000_000), math.ZeroInt(), math.ZeroInt())
	require.NoError(t, err)

	// Large swap selling token0 should push price down and cross into position 2
	trader := fundUser(bk, 2, denomA, math.NewInt(10_000_000))
	amountOut, err := k.ExecuteSwap(ctx, trader, poolID, denomA, math.NewInt(3_000_000), math.ZeroInt())
	require.NoError(t, err)
	require.True(t, amountOut.IsPositive(), "should get output even when crossing ticks")
}

func TestOutOfRangePositionNoFees(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	poolID := createTestPool(t, k, ctx)
	lp := fundUser(bk, 1, denomA, math.NewInt(20_000_000))
	fundUser(bk, 1, denomB, math.NewInt(20_000_000))

	// Position 1: in range [-100, 100] (will earn fees)
	inRangePosID, _, _, _, err := k.ExecuteCreatePosition(ctx, lp, poolID, -100, 100, math.NewInt(5_000_000), math.NewInt(5_000_000), math.ZeroInt(), math.ZeroInt())
	require.NoError(t, err)

	// Position 2: out of range [500, 600] (should NOT earn fees from swaps near tick 0)
	outRangePosID, _, _, _, err := k.ExecuteCreatePosition(ctx, lp, poolID, 500, 600, math.NewInt(5_000_000), math.NewInt(5_000_000), math.ZeroInt(), math.ZeroInt())
	require.NoError(t, err)

	// Execute a swap within range [-100, 100] to generate fees
	trader := fundUser(bk, 2, denomA, math.NewInt(1_000_000))
	_, err = k.ExecuteSwap(ctx, trader, poolID, denomA, math.NewInt(100_000), math.ZeroInt())
	require.NoError(t, err)

	// Collect fees for in-range position
	fees0In, fees1In, err := k.CollectPositionFees(ctx, lp, inRangePosID)
	require.NoError(t, err)
	totalFeesIn := fees0In.Add(fees1In)

	// Collect fees for out-of-range position
	fees0Out, fees1Out, err := k.CollectPositionFees(ctx, lp, outRangePosID)
	require.NoError(t, err)
	totalFeesOut := fees0Out.Add(fees1Out)

	require.True(t, totalFeesIn.IsPositive(), "in-range position should earn fees")
	require.True(t, totalFeesOut.IsZero(), "out-of-range position should NOT earn fees")
}

func TestCollectFeesPositionNotFound(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	_, _, err := k.CollectPositionFees(ctx, testAddr(1), 999)
	require.ErrorIs(t, err, types.ErrPositionNotFound)
}

func TestCollectFeesUnauthorized(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	poolID := createTestPool(t, k, ctx)
	user1 := fundUser(bk, 1, denomA, math.NewInt(1_000_000))
	fundUser(bk, 1, denomB, math.NewInt(1_000_000))

	posID, _, _, _, err := k.ExecuteCreatePosition(ctx, user1, poolID, -100, 100, math.NewInt(500_000), math.NewInt(500_000), math.ZeroInt(), math.ZeroInt())
	require.NoError(t, err)

	_, _, err = k.CollectPositionFees(ctx, testAddr(2), posID)
	require.ErrorIs(t, err, types.ErrUnauthorized)
}

func TestMultiplePoolsIndependent(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	// Create two pools with different pairs
	poolID1, err := k.CreatePool(ctx, testAddr(99), "aaa", "bbb", 10, math.LegacyNewDecWithPrec(3, 3), math.LegacyOneDec())
	require.NoError(t, err)
	poolID2, err := k.CreatePool(ctx, testAddr(99), "ccc", "ddd", 20, math.LegacyNewDecWithPrec(5, 3), math.LegacyNewDec(2))
	require.NoError(t, err)
	require.NotEqual(t, poolID1, poolID2)

	pool1, found := k.GetPool(ctx, poolID1)
	require.True(t, found)
	pool2, found := k.GetPool(ctx, poolID2)
	require.True(t, found)

	require.Equal(t, "aaa", pool1.DenomA)
	require.Equal(t, "ccc", pool2.DenomA)
	require.Equal(t, int64(10), pool1.TickSpacing)
	require.Equal(t, int64(20), pool2.TickSpacing)

	_ = bk // not needed for this test
}

func TestTickToSqrtPriceAndBack(t *testing.T) {
	// Test that tick conversion is roughly invertible
	for _, tick := range []int64{0, 100, -100, 1000, -1000, 10000} {
		sqrtPrice := types.TickToSqrtPrice(tick)
		require.True(t, sqrtPrice.IsPositive(), "sqrt price should be positive for tick %d", tick)

		recovered := types.SqrtPriceToTick(sqrtPrice)
		// Allow +-1 tick tolerance due to floating point
		diff := tick - recovered
		if diff < 0 { diff = -diff }
		require.True(t, diff <= 1, "tick conversion should be roughly invertible: tick=%d, recovered=%d", tick, recovered)
	}
}

func TestSwapPoolNotFound(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	trader := fundUser(bk, 1, denomA, math.NewInt(100_000))
	_, err := k.ExecuteSwap(ctx, trader, 999, denomA, math.NewInt(100), math.ZeroInt())
	require.ErrorIs(t, err, types.ErrPoolNotFound)
}
