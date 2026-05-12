package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	"syreen/x/lending/keeper"
	"syreen/x/lending/types"

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

// mockDexKeeper returns a configurable spot price
type mockDexKeeper struct {
	spotPrice math.LegacyDec
}

func newMockDexKeeper(price math.LegacyDec) *mockDexKeeper {
	return &mockDexKeeper{spotPrice: price}
}

func (m *mockDexKeeper) setPrice(price math.LegacyDec) {
	m.spotPrice = price
}

func (m *mockDexKeeper) GetSpotPrice(_ context.Context, _ uint64, _, _ string) (math.LegacyDec, error) {
	return m.spotPrice, nil
}

// ============================================================
// Test Setup
// ============================================================

func setupKeeper(t *testing.T) (keeper.Keeper, sdk.Context, *mockBankKeeper, *mockDexKeeper) {
	t.Helper()
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, cmtproto.Header{Height: 1}, false, log.NewNopLogger())

	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())

	accountKeeper := mockAccountKeeper{
		moduleAddrs: map[string]sdk.AccAddress{
			types.ModuleName: authtypes.NewModuleAddress(types.ModuleName),
		},
	}
	bankKeeper := newMockBankKeeper()

	k := keeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		accountKeeper,
		bankKeeper,
		authtypes.NewModuleAddress("gov").String(),
	)

	dexKeeper := newMockDexKeeper(math.LegacyNewDec(10)) // default price $10
	k.SetDexKeeper(dexKeeper)

	return k, ctx, bankKeeper, dexKeeper
}

func testAddr(i int) string {
	return sdk.AccAddress([]byte{byte(i), 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}).String()
}

const denomA = "usyreen"
const denomB = "uusdc"

var authority = authtypes.NewModuleAddress("gov").String()

// createTestPool creates a lending pool with default params and returns its ID.
func createTestPool(t *testing.T, k keeper.Keeper, ctx sdk.Context, denom string) uint64 {
	t.Helper()
	id, err := k.CreateNewLendingPool(ctx, authority, denom, math.LegacyNewDecWithPrec(75, 2), 1, "uusdc")
	require.NoError(t, err)
	return id
}

// fundUser sets the user's balance for a denom and returns the address string.
func fundUser(bk *mockBankKeeper, userIdx int, denom string, amount math.Int) string {
	addr := testAddr(userIdx)
	bk.setBalance(addr, denom, amount)
	return addr
}

// fundModule sets the module account balance for a denom.
func fundModule(bk *mockBankKeeper, denom string, amount math.Int) {
	moduleAddr := authtypes.NewModuleAddress(types.ModuleName).String()
	bk.setBalance(moduleAddr, denom, amount)
}

// ============================================================
// Tests
// ============================================================

func TestCreateLendingPool(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	id, err := k.CreateNewLendingPool(ctx, authority, denomA, math.LegacyNewDecWithPrec(75, 2), 1, "uusdc")
	require.NoError(t, err)
	require.Equal(t, uint64(1), id)

	pool, ok := k.GetPool(ctx, id)
	require.True(t, ok)
	require.Equal(t, uint64(1), pool.ID)
	require.Equal(t, denomA, pool.Denom)
	require.True(t, pool.TotalDeposited.IsZero())
	require.True(t, pool.TotalBorrowed.IsZero())
	require.True(t, pool.CollateralFactor.Equal(math.LegacyNewDecWithPrec(75, 2)))
	require.True(t, pool.LiquidationBonus.Equal(types.DefaultLiquidationBonus))
	require.True(t, pool.LiquidationThreshold.Equal(math.LegacyNewDecWithPrec(80, 2)))
	require.True(t, pool.ReserveRatio.Equal(math.LegacyNewDecWithPrec(10, 2)))
	require.True(t, pool.UtilizationRate.IsZero())
	require.True(t, pool.AccInterestPerShare.IsZero())
	require.True(t, pool.AccBorrowIndex.Equal(math.LegacyOneDec()))
	require.True(t, pool.Active)
	require.Equal(t, uint64(1), pool.DexPoolID)
	require.Equal(t, "uusdc", pool.PriceDenom)
	require.Equal(t, int64(1), pool.LastUpdateBlock)

	// Next pool ID should be incremented
	require.Equal(t, uint64(2), k.GetNextPoolID(ctx))
}

func TestCreateLendingPoolUnauthorized(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	_, err := k.CreateNewLendingPool(ctx, testAddr(1), denomA, math.LegacyNewDecWithPrec(75, 2), 1, "uusdc")
	require.ErrorIs(t, err, types.ErrUnauthorized)
}

func TestDeposit(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	poolID := createTestPool(t, k, ctx, denomA)

	user := fundUser(bk, 1, denomA, math.NewInt(10000))
	depositAmt := math.NewInt(5000)

	_, err := k.ExecuteDeposit(ctx, user, poolID, depositAmt)
	require.NoError(t, err)

	// Verify user balance decreased
	require.True(t, bk.getBalance(user, denomA).Equal(math.NewInt(5000)))

	// Verify module received the tokens
	moduleAddr := authtypes.NewModuleAddress(types.ModuleName).String()
	require.True(t, bk.getBalance(moduleAddr, denomA).Equal(math.NewInt(5000)))

	// Verify pool total deposited updated
	pool, ok := k.GetPool(ctx, poolID)
	require.True(t, ok)
	require.True(t, pool.TotalDeposited.Equal(math.NewInt(5000)))

	// Verify deposit record created
	dep, ok := k.GetDeposit(ctx, poolID, user)
	require.True(t, ok)
	require.Equal(t, user, dep.Address)
	require.Equal(t, poolID, dep.PoolID)
	require.True(t, dep.Amount.Equal(math.NewInt(5000)))
	require.Equal(t, int64(1), dep.DepositedAt)
}

func TestDepositMultiple(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	poolID := createTestPool(t, k, ctx, denomA)

	user := fundUser(bk, 1, denomA, math.NewInt(10000))

	// First deposit
	_, err := k.ExecuteDeposit(ctx, user, poolID, math.NewInt(3000))
	require.NoError(t, err)

	// Second deposit
	_, err = k.ExecuteDeposit(ctx, user, poolID, math.NewInt(2000))
	require.NoError(t, err)

	// Verify user balance
	require.True(t, bk.getBalance(user, denomA).Equal(math.NewInt(5000)))

	// Verify deposit record accumulated
	dep, ok := k.GetDeposit(ctx, poolID, user)
	require.True(t, ok)
	require.True(t, dep.Amount.Equal(math.NewInt(5000)))

	// Verify pool total
	pool, ok := k.GetPool(ctx, poolID)
	require.True(t, ok)
	require.True(t, pool.TotalDeposited.Equal(math.NewInt(5000)))
}

func TestWithdraw(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	poolID := createTestPool(t, k, ctx, denomA)

	user := fundUser(bk, 1, denomA, math.NewInt(10000))

	// Deposit first
	_, err := k.ExecuteDeposit(ctx, user, poolID, math.NewInt(5000))
	require.NoError(t, err)

	// Withdraw
	_, err = k.ExecuteWithdraw(ctx, user, poolID, math.NewInt(3000))
	require.NoError(t, err)

	// Verify user got tokens back
	require.True(t, bk.getBalance(user, denomA).Equal(math.NewInt(8000)))

	// Verify deposit reduced
	dep, ok := k.GetDeposit(ctx, poolID, user)
	require.True(t, ok)
	require.True(t, dep.Amount.Equal(math.NewInt(2000)))

	// Verify pool total reduced
	pool, ok := k.GetPool(ctx, poolID)
	require.True(t, ok)
	require.True(t, pool.TotalDeposited.Equal(math.NewInt(2000)))
}

func TestWithdrawInsufficient(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	poolID := createTestPool(t, k, ctx, denomA)

	user := fundUser(bk, 1, denomA, math.NewInt(10000))

	// Deposit 5000
	_, err := k.ExecuteDeposit(ctx, user, poolID, math.NewInt(5000))
	require.NoError(t, err)

	// Try to withdraw more than deposited
	_, err = k.ExecuteWithdraw(ctx, user, poolID, math.NewInt(6000))
	require.ErrorIs(t, err, types.ErrInsufficientDeposit)
}

func TestWithdrawNoDeposit(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	poolID := createTestPool(t, k, ctx, denomA)

	user := testAddr(1)

	// Withdraw without depositing
	_, err := k.ExecuteWithdraw(ctx, user, poolID, math.NewInt(1000))
	require.ErrorIs(t, err, types.ErrNoDeposit)
}

func TestBorrow(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)

	// Create two pools: denomA for collateral, denomB for borrowing
	collateralPoolID := createTestPool(t, k, ctx, denomA)
	borrowPoolID := createTestPool(t, k, ctx, denomB)

	// User has denomA for collateral
	user := fundUser(bk, 1, denomA, math.NewInt(10000))

	// Lender deposits denomB into the borrow pool so there is liquidity
	lender := fundUser(bk, 2, denomB, math.NewInt(50000))
	_, err := k.ExecuteDeposit(ctx, lender, borrowPoolID, math.NewInt(50000))
	require.NoError(t, err)

	// User borrows 5000 denomB with 10000 denomA collateral
	// Collateral factor = 0.75, so max borrow = 10000 * 0.75 = 7500
	borrowAmt := math.NewInt(5000)
	collateralAmt := math.NewInt(10000)

	borrowID, err := k.ExecuteBorrow(ctx, user, borrowPoolID, borrowAmt, collateralPoolID, collateralAmt)
	require.NoError(t, err)
	require.Equal(t, uint64(1), borrowID)

	// Verify collateral was taken from user
	require.True(t, bk.getBalance(user, denomA).IsZero())

	// Verify user received borrowed tokens
	require.True(t, bk.getBalance(user, denomB).Equal(math.NewInt(5000)))

	// Verify borrow record
	borrow, ok := k.GetBorrow(ctx, borrowID)
	require.True(t, ok)
	require.Equal(t, user, borrow.Borrower)
	require.Equal(t, borrowPoolID, borrow.BorrowPoolID)
	require.True(t, borrow.BorrowAmount.Equal(math.NewInt(5000)))
	require.Equal(t, collateralPoolID, borrow.CollateralPoolID)
	require.True(t, borrow.CollateralAmount.Equal(math.NewInt(10000)))
	require.Equal(t, int64(1), borrow.BorrowedAt)

	// Verify pool total borrowed updated
	pool, ok := k.GetPool(ctx, borrowPoolID)
	require.True(t, ok)
	require.True(t, pool.TotalBorrowed.Equal(math.NewInt(5000)))
}

func TestBorrowOverCollateral(t *testing.T) {
	k, ctx, bk, dk := setupKeeper(t)

	// Set price to 1:1 so collateral value equals token amount
	// collateral = 10000 denomA, price = 1 -> value = 10000 denomB
	// collateral factor = 0.75 -> max borrow = 7500 denomB
	dk.setPrice(math.LegacyOneDec())

	collateralPoolID := createTestPool(t, k, ctx, denomA)
	borrowPoolID := createTestPool(t, k, ctx, denomB)

	user := fundUser(bk, 1, denomA, math.NewInt(10000))

	// Lender deposits into borrow pool
	lender := fundUser(bk, 2, denomB, math.NewInt(50000))
	_, err := k.ExecuteDeposit(ctx, lender, borrowPoolID, math.NewInt(50000))
	require.NoError(t, err)

	// Collateral factor = 0.75, collateral = 10000, spot price = 1:1, max borrow = 7500
	// Try to borrow 8000 — should fail
	_, err = k.ExecuteBorrow(ctx, user, borrowPoolID, math.NewInt(8000), collateralPoolID, math.NewInt(10000))
	require.ErrorIs(t, err, types.ErrOverCollateral)
}

func TestBorrowInsufficientLiquidity(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)

	collateralPoolID := createTestPool(t, k, ctx, denomA)
	borrowPoolID := createTestPool(t, k, ctx, denomB)

	user := fundUser(bk, 1, denomA, math.NewInt(10000))

	// No deposits in borrow pool — zero liquidity
	// Try to borrow — should fail
	_, err := k.ExecuteBorrow(ctx, user, borrowPoolID, math.NewInt(1000), collateralPoolID, math.NewInt(10000))
	require.ErrorIs(t, err, types.ErrInsufficientLiquidity)
}

func TestRepayFull(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)

	collateralPoolID := createTestPool(t, k, ctx, denomA)
	borrowPoolID := createTestPool(t, k, ctx, denomB)

	user := fundUser(bk, 1, denomA, math.NewInt(10000))

	// Lender deposits
	lender := fundUser(bk, 2, denomB, math.NewInt(50000))
	_, err := k.ExecuteDeposit(ctx, lender, borrowPoolID, math.NewInt(50000))
	require.NoError(t, err)

	// Borrow 5000
	borrowID, err := k.ExecuteBorrow(ctx, user, borrowPoolID, math.NewInt(5000), collateralPoolID, math.NewInt(10000))
	require.NoError(t, err)

	// User now has 5000 denomB. Repay full amount (pass 0 to repay all)
	_, collateralReturned, err := k.ExecuteRepay(ctx, user, borrowID, math.ZeroInt())
	require.NoError(t, err)

	// Collateral should be returned
	require.True(t, collateralReturned.Equal(math.NewInt(10000)))
	require.True(t, bk.getBalance(user, denomA).Equal(math.NewInt(10000)))

	// Borrow should be deleted
	_, ok := k.GetBorrow(ctx, borrowID)
	require.False(t, ok)

	// Pool total borrowed should be zero
	pool, ok := k.GetPool(ctx, borrowPoolID)
	require.True(t, ok)
	require.True(t, pool.TotalBorrowed.IsZero())
}

func TestRepayPartial(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)

	collateralPoolID := createTestPool(t, k, ctx, denomA)
	borrowPoolID := createTestPool(t, k, ctx, denomB)

	user := fundUser(bk, 1, denomA, math.NewInt(10000))

	// Lender deposits
	lender := fundUser(bk, 2, denomB, math.NewInt(50000))
	_, err := k.ExecuteDeposit(ctx, lender, borrowPoolID, math.NewInt(50000))
	require.NoError(t, err)

	// Borrow 5000
	borrowID, err := k.ExecuteBorrow(ctx, user, borrowPoolID, math.NewInt(5000), collateralPoolID, math.NewInt(10000))
	require.NoError(t, err)

	// Partial repay 2000 out of 5000
	_, collateralReturned, err := k.ExecuteRepay(ctx, user, borrowID, math.NewInt(2000))
	require.NoError(t, err)

	// Partial repay should not return collateral
	require.True(t, collateralReturned.IsZero())

	// Borrow should still exist but be reduced
	borrow, ok := k.GetBorrow(ctx, borrowID)
	require.True(t, ok)
	require.True(t, borrow.BorrowAmount.LT(math.NewInt(5000)))

	// Pool total borrowed should be reduced
	pool, ok := k.GetPool(ctx, borrowPoolID)
	require.True(t, ok)
	require.True(t, pool.TotalBorrowed.LT(math.NewInt(5000)))
}

func TestRepayNotOwner(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)

	collateralPoolID := createTestPool(t, k, ctx, denomA)
	borrowPoolID := createTestPool(t, k, ctx, denomB)

	user := fundUser(bk, 1, denomA, math.NewInt(10000))

	// Lender deposits
	lender := fundUser(bk, 2, denomB, math.NewInt(50000))
	_, err := k.ExecuteDeposit(ctx, lender, borrowPoolID, math.NewInt(50000))
	require.NoError(t, err)

	// User borrows
	borrowID, err := k.ExecuteBorrow(ctx, user, borrowPoolID, math.NewInt(5000), collateralPoolID, math.NewInt(10000))
	require.NoError(t, err)

	// Another user tries to repay
	otherUser := fundUser(bk, 3, denomB, math.NewInt(10000))
	_, _, err = k.ExecuteRepay(ctx, otherUser, borrowID, math.NewInt(1000))
	require.ErrorIs(t, err, types.ErrBorrowNotOwner)
}

func TestLiquidate(t *testing.T) {
	k, ctx, bk, dk := setupKeeper(t)

	// Use a 1:1 spot price so collateral value is straightforward to reason about.
	dk.setPrice(math.LegacyOneDec())

	collateralPoolID := createTestPool(t, k, ctx, denomA)
	borrowPoolID := createTestPool(t, k, ctx, denomB)

	user := fundUser(bk, 1, denomA, math.NewInt(10000))

	// Lender deposits
	lender := fundUser(bk, 2, denomB, math.NewInt(50000))
	_, err := k.ExecuteDeposit(ctx, lender, borrowPoolID, math.NewInt(50000))
	require.NoError(t, err)

	// Seed the TWAP series so the borrow check passes against TWAP at price=1.
	k.SampleAllPoolPrices(ctx)

	// Borrow near max: collateral = 10000, factor = 0.75, borrow = 7500
	borrowID, err := k.ExecuteBorrow(ctx, user, borrowPoolID, math.NewInt(7500), collateralPoolID, math.NewInt(10000))
	require.NoError(t, err)

	// Make position unhealthy by inflating the borrow index — debt doubles to 15000.
	pool, _ := k.GetPool(ctx, borrowPoolID)
	pool.AccBorrowIndex = math.LegacyNewDec(2)
	k.SetPool(ctx, pool)

	// Liquidator has funds to pay debt
	liquidator := fundUser(bk, 3, denomB, math.NewInt(50000))

	// Fund module with collateral so it can send it to liquidator
	fundModule(bk, denomA, math.NewInt(10000))

	// First call: marks position unhealthy and returns ErrLiquidationDelay (M-1).
	_, _, err = k.ExecuteLiquidate(ctx, liquidator, borrowID)
	require.ErrorIs(t, err, types.ErrLiquidationDelay)

	// Advance past LiquidationDelayBlocks and retry.
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + types.LiquidationDelayBlocks + 1)

	collateralSeized, debtRepaid, err := k.ExecuteLiquidate(ctx, liquidator, borrowID)
	require.NoError(t, err)
	require.True(t, collateralSeized.IsPositive())
	require.True(t, debtRepaid.IsPositive())

	// Borrow should be deleted after liquidation
	_, ok := k.GetBorrow(ctx, borrowID)
	require.False(t, ok)
}

func TestLiquidateHealthy(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)

	collateralPoolID := createTestPool(t, k, ctx, denomA)
	borrowPoolID := createTestPool(t, k, ctx, denomB)

	user := fundUser(bk, 1, denomA, math.NewInt(10000))

	// Lender deposits
	lender := fundUser(bk, 2, denomB, math.NewInt(50000))
	_, err := k.ExecuteDeposit(ctx, lender, borrowPoolID, math.NewInt(50000))
	require.NoError(t, err)

	// Borrow conservatively: collateral = 10000, borrow only 1000
	// Health factor = 10000 * 0.80 / 1000 = 8.0 — very healthy
	borrowID, err := k.ExecuteBorrow(ctx, user, borrowPoolID, math.NewInt(1000), collateralPoolID, math.NewInt(10000))
	require.NoError(t, err)

	liquidator := fundUser(bk, 3, denomB, math.NewInt(50000))

	_, _, err = k.ExecuteLiquidate(ctx, liquidator, borrowID)
	require.ErrorIs(t, err, types.ErrLiquidationNotEligible)
}

func TestInterestAccrual(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)

	collateralPoolID := createTestPool(t, k, ctx, denomA)
	borrowPoolID := createTestPool(t, k, ctx, denomB)

	user := fundUser(bk, 1, denomA, math.NewInt(100000))

	// Lender deposits 100000 denomB
	lender := fundUser(bk, 2, denomB, math.NewInt(100000))
	_, err := k.ExecuteDeposit(ctx, lender, borrowPoolID, math.NewInt(100000))
	require.NoError(t, err)

	// Borrow 50000 (50% utilization)
	borrowID, err := k.ExecuteBorrow(ctx, user, borrowPoolID, math.NewInt(50000), collateralPoolID, math.NewInt(100000))
	require.NoError(t, err)

	// Record initial borrow index
	poolBefore, _ := k.GetPool(ctx, borrowPoolID)
	indexBefore := poolBefore.AccBorrowIndex

	// Advance blocks significantly
	ctx = ctx.WithBlockHeight(1_000_000)

	// Trigger interest accrual by calling UpdatePoolInterest
	poolAfter, _ := k.GetPool(ctx, borrowPoolID)
	k.UpdatePoolInterest(ctx, &poolAfter)
	k.SetPool(ctx, poolAfter)

	// Borrow index should have increased
	require.True(t, poolAfter.AccBorrowIndex.GT(indexBefore))

	// AccInterestPerShare should have increased (depositors earn interest)
	require.True(t, poolAfter.AccInterestPerShare.GT(math.LegacyZeroDec()))

	// Borrower's debt should be more than original 50000 when calculated via index
	borrow, ok := k.GetBorrow(ctx, borrowID)
	require.True(t, ok)
	indexRatio := poolAfter.AccBorrowIndex.Quo(borrow.BorrowIndex)
	totalDebt := math.LegacyNewDecFromInt(borrow.BorrowAmount).Mul(indexRatio).TruncateInt()
	require.True(t, totalDebt.GT(math.NewInt(50000)))
}

func TestInterestRateModel(t *testing.T) {
	// Test that borrow rates increase with utilization
	lowUtilPool := types.LendingPool{
		TotalDeposited: math.NewInt(100000),
		TotalBorrowed:  math.NewInt(10000), // 10% utilization
		ReserveRatio:   math.LegacyNewDecWithPrec(10, 2),
	}
	_, borrowAPYLow, utilLow := keeper.CalculateInterestRates(lowUtilPool)
	require.True(t, utilLow.Equal(math.LegacyNewDecWithPrec(10, 2)))

	medUtilPool := types.LendingPool{
		TotalDeposited: math.NewInt(100000),
		TotalBorrowed:  math.NewInt(50000), // 50% utilization
		ReserveRatio:   math.LegacyNewDecWithPrec(10, 2),
	}
	_, borrowAPYMed, utilMed := keeper.CalculateInterestRates(medUtilPool)
	require.True(t, utilMed.Equal(math.LegacyNewDecWithPrec(50, 2)))

	highUtilPool := types.LendingPool{
		TotalDeposited: math.NewInt(100000),
		TotalBorrowed:  math.NewInt(90000), // 90% utilization — above optimal
		ReserveRatio:   math.LegacyNewDecWithPrec(10, 2),
	}
	_, borrowAPYHigh, utilHigh := keeper.CalculateInterestRates(highUtilPool)
	require.True(t, utilHigh.Equal(math.LegacyNewDecWithPrec(90, 2)))

	// Rates should increase with utilization
	require.True(t, borrowAPYMed.GT(borrowAPYLow), "medium util rate should exceed low util rate")
	require.True(t, borrowAPYHigh.GT(borrowAPYMed), "high util rate should exceed medium util rate")

	// Above optimal utilization (80%), rate should be steep
	// 90% utilization = above kink, should be significantly higher
	require.True(t, borrowAPYHigh.GT(math.LegacyNewDecWithPrec(10, 2)), "high util rate should be > 10%")

	// Zero deposits = zero rates
	emptyPool := types.LendingPool{
		TotalDeposited: math.ZeroInt(),
		TotalBorrowed:  math.ZeroInt(),
		ReserveRatio:   math.LegacyNewDecWithPrec(10, 2),
	}
	depAPY, borrowAPY, util := keeper.CalculateInterestRates(emptyPool)
	require.True(t, depAPY.IsZero())
	require.True(t, borrowAPY.Equal(types.BaseRate))
	require.True(t, util.IsZero())
}

func TestGetAllPools(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	// No pools initially
	pools := k.GetAllPools(ctx)
	require.Len(t, pools, 0)

	// Create 3 pools
	createTestPool(t, k, ctx, denomA)
	createTestPool(t, k, ctx, denomB)
	createTestPool(t, k, ctx, "uatom")

	pools = k.GetAllPools(ctx)
	require.Len(t, pools, 3)

	// Verify they have sequential IDs
	require.Equal(t, uint64(1), pools[0].ID)
	require.Equal(t, uint64(2), pools[1].ID)
	require.Equal(t, uint64(3), pools[2].ID)

	// Verify denoms
	require.Equal(t, denomA, pools[0].Denom)
	require.Equal(t, denomB, pools[1].Denom)
	require.Equal(t, "uatom", pools[2].Denom)
}

func TestGetDepositsByAddress(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)

	pool1 := createTestPool(t, k, ctx, denomA)
	pool2 := createTestPool(t, k, ctx, denomB)

	user := testAddr(1)
	bk.setBalance(user, denomA, math.NewInt(10000))
	bk.setBalance(user, denomB, math.NewInt(10000))

	// Deposit into two pools
	_, err := k.ExecuteDeposit(ctx, user, pool1, math.NewInt(3000))
	require.NoError(t, err)
	_, err = k.ExecuteDeposit(ctx, user, pool2, math.NewInt(5000))
	require.NoError(t, err)

	// Another user deposits into pool1 only
	otherUser := fundUser(bk, 2, denomA, math.NewInt(10000))
	_, err = k.ExecuteDeposit(ctx, otherUser, pool1, math.NewInt(2000))
	require.NoError(t, err)

	// Query deposits by user
	deposits := k.GetDepositsByAddress(ctx, user)
	require.Len(t, deposits, 2)

	// Verify amounts (order may vary, so check both)
	amounts := map[uint64]math.Int{}
	for _, d := range deposits {
		amounts[d.PoolID] = d.Amount
	}
	require.True(t, amounts[pool1].Equal(math.NewInt(3000)))
	require.True(t, amounts[pool2].Equal(math.NewInt(5000)))

	// Query other user — should have 1 deposit
	otherDeposits := k.GetDepositsByAddress(ctx, otherUser)
	require.Len(t, otherDeposits, 1)
	require.True(t, otherDeposits[0].Amount.Equal(math.NewInt(2000)))
}

func TestGetBorrowsByAddress(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)

	collateralPoolID := createTestPool(t, k, ctx, denomA)
	borrowPoolID := createTestPool(t, k, ctx, denomB)

	// Lender deposits liquidity
	lender := fundUser(bk, 2, denomB, math.NewInt(100000))
	_, err := k.ExecuteDeposit(ctx, lender, borrowPoolID, math.NewInt(100000))
	require.NoError(t, err)

	// User 1 borrows twice
	user := fundUser(bk, 1, denomA, math.NewInt(100000))
	borrowID1, err := k.ExecuteBorrow(ctx, user, borrowPoolID, math.NewInt(5000), collateralPoolID, math.NewInt(10000))
	require.NoError(t, err)
	borrowID2, err := k.ExecuteBorrow(ctx, user, borrowPoolID, math.NewInt(3000), collateralPoolID, math.NewInt(8000))
	require.NoError(t, err)

	// User 2 borrows once
	user2 := fundUser(bk, 3, denomA, math.NewInt(50000))
	_, err = k.ExecuteBorrow(ctx, user2, borrowPoolID, math.NewInt(2000), collateralPoolID, math.NewInt(5000))
	require.NoError(t, err)

	// Query borrows for user 1
	borrows := k.GetBorrowsByAddress(ctx, user)
	require.Len(t, borrows, 2)

	borrowIDs := map[uint64]bool{}
	for _, b := range borrows {
		borrowIDs[b.ID] = true
	}
	require.True(t, borrowIDs[borrowID1])
	require.True(t, borrowIDs[borrowID2])

	// Query borrows for user 2 — should have 1
	borrows2 := k.GetBorrowsByAddress(ctx, user2)
	require.Len(t, borrows2, 1)
	require.True(t, borrows2[0].BorrowAmount.Equal(math.NewInt(2000)))
}
