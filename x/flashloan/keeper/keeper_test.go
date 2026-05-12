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

	"syreen/x/flashloan/keeper"
	"syreen/x/flashloan/types"

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

	return k, ctx, bankKeeper
}

func testAddr(i int) string {
	return sdk.AccAddress([]byte{byte(i), 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}).String()
}

const denomA = "usyreen"
const denomB = "uusdc"

var authority = authtypes.NewModuleAddress("gov").String()

func moduleAddr() string {
	return authtypes.NewModuleAddress(types.ModuleName).String()
}

// fundUser sets the user's balance for a denom and returns the address string.
func fundUser(bk *mockBankKeeper, userIdx int, denom string, amount math.Int) string {
	addr := testAddr(userIdx)
	bk.setBalance(addr, denom, amount)
	return addr
}

// fundModule sets the module account balance for a denom.
func fundModule(bk *mockBankKeeper, denom string, amount math.Int) {
	bk.setBalance(moduleAddr(), denom, amount)
}

// ============================================================
// Tests
// ============================================================

func TestCreateFlashPool(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	// Create pool with authority
	err := k.CreatePool(ctx, authority, denomA, types.DefaultFeeRate)
	require.NoError(t, err)

	pool, found := k.GetPool(ctx, denomA)
	require.True(t, found)
	require.Equal(t, denomA, pool.Denom)
	require.True(t, pool.Active)
	require.True(t, pool.Available.IsZero())
	require.True(t, pool.FeeRate.Equal(types.DefaultFeeRate))
}

func TestCreateFlashPoolUnauthorized(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	err := k.CreatePool(ctx, testAddr(1), denomA, types.DefaultFeeRate)
	require.ErrorIs(t, err, types.ErrUnauthorized)
}

func TestCreateFlashPoolDuplicate(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	err := k.CreatePool(ctx, authority, denomA, types.DefaultFeeRate)
	require.NoError(t, err)

	err = k.CreatePool(ctx, authority, denomA, types.DefaultFeeRate)
	require.ErrorIs(t, err, types.ErrPoolAlreadyExists)
}

func TestFundFlashPool(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	err := k.CreatePool(ctx, authority, denomA, types.DefaultFeeRate)
	require.NoError(t, err)

	user := fundUser(bk, 1, denomA, math.NewInt(1_000_000))

	err = k.FundPool(ctx, user, denomA, math.NewInt(500_000))
	require.NoError(t, err)

	pool, _ := k.GetPool(ctx, denomA)
	require.Equal(t, math.NewInt(500_000), pool.Available)

	// User should have 500_000 left
	userBal := bk.getBalance(user, denomA)
	require.Equal(t, math.NewInt(500_000), userBal)

	// Module should have 500_000
	modBal := bk.getBalance(moduleAddr(), denomA)
	require.Equal(t, math.NewInt(500_000), modBal)
}

func TestFundFlashPoolNotFound(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	user := fundUser(bk, 1, denomA, math.NewInt(1_000_000))
	err := k.FundPool(ctx, user, denomA, math.NewInt(500_000))
	require.ErrorIs(t, err, types.ErrPoolNotFound)
}

func TestFlashLoanSuccess(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	// Create and fund pool
	err := k.CreatePool(ctx, authority, denomA, types.DefaultFeeRate)
	require.NoError(t, err)

	funder := fundUser(bk, 1, denomA, math.NewInt(10_000_000))
	err = k.FundPool(ctx, funder, denomA, math.NewInt(10_000_000))
	require.NoError(t, err)

	// User who wants to flash loan needs to have principal + fee pre-funded
	// Borrow 1_000_000, fee = 1_000_000 * 0.0009 = 900
	borrowAmt := math.NewInt(1_000_000)
	fee := math.NewInt(900)
	totalNeeded := borrowAmt.Add(fee) // user needs this to repay after receiving borrow

	user := fundUser(bk, 2, denomA, totalNeeded)

	feeResult, err := k.ExecuteFlashLoan(ctx, user, denomA, borrowAmt)
	require.NoError(t, err)
	require.Equal(t, fee, feeResult)

	// After flash loan: user started with 1_000_900, got +1_000_000, paid -(1_000_000+900)
	// Net: 1_000_900 + 1_000_000 - 1_000_900 = 1_000_000
	userBal := bk.getBalance(user, denomA)
	require.Equal(t, math.NewInt(1_000_000), userBal)

	// Module should have original 10M + 900 fee
	modBal := bk.getBalance(moduleAddr(), denomA)
	require.Equal(t, math.NewInt(10_000_900), modBal)

	// Pool stats updated
	pool, _ := k.GetPool(ctx, denomA)
	require.Equal(t, math.NewInt(10_000_900), pool.Available) // original + fee
	require.Equal(t, borrowAmt, pool.TotalLoaned)
	require.Equal(t, fee, pool.TotalFees)
}

func TestFlashLoanFeeCalculation(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	// Custom fee rate: 0.1% = 0.001
	customFeeRate := math.LegacyNewDecWithPrec(1, 3)
	err := k.CreatePool(ctx, authority, denomA, customFeeRate)
	require.NoError(t, err)

	funder := fundUser(bk, 1, denomA, math.NewInt(50_000_000))
	err = k.FundPool(ctx, funder, denomA, math.NewInt(50_000_000))
	require.NoError(t, err)

	// Borrow 10_000_000, fee = 10_000_000 * 0.001 = 10_000
	borrowAmt := math.NewInt(10_000_000)
	expectedFee := math.NewInt(10_000)

	user := fundUser(bk, 2, denomA, borrowAmt.Add(expectedFee))

	feeResult, err := k.ExecuteFlashLoan(ctx, user, denomA, borrowAmt)
	require.NoError(t, err)
	require.Equal(t, expectedFee, feeResult)
}

func TestFlashLoanMinimumFee(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	err := k.CreatePool(ctx, authority, denomA, types.DefaultFeeRate)
	require.NoError(t, err)

	funder := fundUser(bk, 1, denomA, math.NewInt(1_000_000))
	err = k.FundPool(ctx, funder, denomA, math.NewInt(1_000_000))
	require.NoError(t, err)

	// Borrow 100 tokens, fee = 100 * 0.0009 = 0.09 -> truncated to 0 -> min 1
	borrowAmt := math.NewInt(100)
	expectedFee := math.NewInt(1) // minimum fee

	user := fundUser(bk, 2, denomA, borrowAmt.Add(expectedFee))

	feeResult, err := k.ExecuteFlashLoan(ctx, user, denomA, borrowAmt)
	require.NoError(t, err)
	require.Equal(t, expectedFee, feeResult)
}

func TestFlashLoanInsufficientRepayment(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	err := k.CreatePool(ctx, authority, denomA, types.DefaultFeeRate)
	require.NoError(t, err)

	funder := fundUser(bk, 1, denomA, math.NewInt(10_000_000))
	err = k.FundPool(ctx, funder, denomA, math.NewInt(10_000_000))
	require.NoError(t, err)

	// User has less than the fee amount, so after receiving borrow they can't repay principal + fee
	// Borrow 1_000_000, fee=900. User starts with 0.
	// After receiving borrow: user has 1_000_000. Must repay 1_000_900. Fails.
	borrowAmt := math.NewInt(1_000_000)
	user := fundUser(bk, 2, denomA, math.ZeroInt()) // has nothing, can't cover the fee

	_, err = k.ExecuteFlashLoan(ctx, user, denomA, borrowAmt)
	require.Error(t, err)
	require.Contains(t, err.Error(), "repayment failed")
}

func TestFlashLoanInsufficientLiquidity(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	err := k.CreatePool(ctx, authority, denomA, types.DefaultFeeRate)
	require.NoError(t, err)

	funder := fundUser(bk, 1, denomA, math.NewInt(2_000))
	err = k.FundPool(ctx, funder, denomA, math.NewInt(2_000))
	require.NoError(t, err)

	// Try to borrow more than available
	borrowAmt := math.NewInt(10_000)
	user := fundUser(bk, 2, denomA, math.NewInt(100_000))

	_, err = k.ExecuteFlashLoan(ctx, user, denomA, borrowAmt)
	require.ErrorIs(t, err, types.ErrInsufficientLiquidity)
}

func TestFlashLoanPoolNotFound(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	_, err := k.ExecuteFlashLoan(ctx, testAddr(1), "nonexistent", math.NewInt(1000))
	require.ErrorIs(t, err, types.ErrPoolNotFound)
}

func TestFlashLoanStatsTracking(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	err := k.CreatePool(ctx, authority, denomA, types.DefaultFeeRate)
	require.NoError(t, err)

	funder := fundUser(bk, 1, denomA, math.NewInt(100_000_000))
	err = k.FundPool(ctx, funder, denomA, math.NewInt(100_000_000))
	require.NoError(t, err)

	// Execute first flash loan
	borrow1 := math.NewInt(1_000_000)
	fee1 := math.NewInt(900)
	user1 := fundUser(bk, 2, denomA, borrow1.Add(fee1))
	_, err = k.ExecuteFlashLoan(ctx, user1, denomA, borrow1)
	require.NoError(t, err)

	stats := k.GetStats(ctx)
	require.Equal(t, uint64(1), stats.TotalExecuted)
	require.Equal(t, borrow1, stats.TotalVolume)
	require.Equal(t, fee1, stats.TotalFeesEarned)

	// Execute second flash loan
	borrow2 := math.NewInt(5_000_000)
	fee2 := math.NewInt(4500)
	user2 := fundUser(bk, 3, denomA, borrow2.Add(fee2))
	_, err = k.ExecuteFlashLoan(ctx, user2, denomA, borrow2)
	require.NoError(t, err)

	stats = k.GetStats(ctx)
	require.Equal(t, uint64(2), stats.TotalExecuted)
	require.Equal(t, borrow1.Add(borrow2), stats.TotalVolume)
	require.Equal(t, fee1.Add(fee2), stats.TotalFeesEarned)
}

func TestFlashLoanRecordTracking(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	err := k.CreatePool(ctx, authority, denomA, types.DefaultFeeRate)
	require.NoError(t, err)

	funder := fundUser(bk, 1, denomA, math.NewInt(100_000_000))
	err = k.FundPool(ctx, funder, denomA, math.NewInt(100_000_000))
	require.NoError(t, err)

	borrowAmt := math.NewInt(2_000_000)
	fee := math.NewInt(1800)
	user := fundUser(bk, 2, denomA, borrowAmt.Add(fee))

	_, err = k.ExecuteFlashLoan(ctx, user, denomA, borrowAmt)
	require.NoError(t, err)

	record, found := k.GetRecord(ctx, 1)
	require.True(t, found)
	require.Equal(t, uint64(1), record.ID)
	require.Equal(t, user, record.Borrower)
	require.Equal(t, denomA, record.Denom)
	require.Equal(t, borrowAmt, record.Amount)
	require.Equal(t, fee, record.Fee)
	require.Equal(t, int64(1), record.Block)
	require.True(t, record.Success)
}

func TestMultipleFlashLoans(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	err := k.CreatePool(ctx, authority, denomA, types.DefaultFeeRate)
	require.NoError(t, err)

	funder := fundUser(bk, 1, denomA, math.NewInt(100_000_000))
	err = k.FundPool(ctx, funder, denomA, math.NewInt(100_000_000))
	require.NoError(t, err)

	// Execute 3 flash loans
	for i := 0; i < 3; i++ {
		borrowAmt := math.NewInt(1_000_000)
		fee := math.NewInt(900)
		user := fundUser(bk, 10+i, denomA, borrowAmt.Add(fee))
		_, err = k.ExecuteFlashLoan(ctx, user, denomA, borrowAmt)
		require.NoError(t, err)
	}

	stats := k.GetStats(ctx)
	require.Equal(t, uint64(3), stats.TotalExecuted)
	require.Equal(t, math.NewInt(3_000_000), stats.TotalVolume)
	require.Equal(t, math.NewInt(2700), stats.TotalFeesEarned)

	pool, _ := k.GetPool(ctx, denomA)
	require.Equal(t, math.NewInt(100_002_700), pool.Available)
}

func TestMultipleDenomPools(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	// Create pools for two denoms
	err := k.CreatePool(ctx, authority, denomA, types.DefaultFeeRate)
	require.NoError(t, err)
	err = k.CreatePool(ctx, authority, denomB, types.DefaultFeeRate)
	require.NoError(t, err)

	// Fund both pools
	funderA := fundUser(bk, 1, denomA, math.NewInt(10_000_000))
	err = k.FundPool(ctx, funderA, denomA, math.NewInt(10_000_000))
	require.NoError(t, err)

	funderB := fundUser(bk, 2, denomB, math.NewInt(10_000_000))
	err = k.FundPool(ctx, funderB, denomB, math.NewInt(10_000_000))
	require.NoError(t, err)

	// Flash loan from denomA
	borrowAmt := math.NewInt(1_000_000)
	fee := math.NewInt(900)
	user := fundUser(bk, 3, denomA, borrowAmt.Add(fee))
	_, err = k.ExecuteFlashLoan(ctx, user, denomA, borrowAmt)
	require.NoError(t, err)

	// Flash loan from denomB
	user2 := fundUser(bk, 4, denomB, borrowAmt.Add(fee))
	_, err = k.ExecuteFlashLoan(ctx, user2, denomB, borrowAmt)
	require.NoError(t, err)

	// Both pools should have their fees
	poolA, _ := k.GetPool(ctx, denomA)
	require.Equal(t, math.NewInt(10_000_900), poolA.Available)

	poolB, _ := k.GetPool(ctx, denomB)
	require.Equal(t, math.NewInt(10_000_900), poolB.Available)

	// Global stats should aggregate
	stats := k.GetStats(ctx)
	require.Equal(t, uint64(2), stats.TotalExecuted)
	require.Equal(t, math.NewInt(2_000_000), stats.TotalVolume)
	require.Equal(t, math.NewInt(1800), stats.TotalFeesEarned)
}

func TestMsgFlashLoanViaServer(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	err := k.CreatePool(ctx, authority, denomA, types.DefaultFeeRate)
	require.NoError(t, err)

	funder := fundUser(bk, 1, denomA, math.NewInt(10_000_000))
	err = k.FundPool(ctx, funder, denomA, math.NewInt(10_000_000))
	require.NoError(t, err)

	borrowAmt := math.NewInt(1_000_000)
	fee := math.NewInt(900)
	user := fundUser(bk, 2, denomA, borrowAmt.Add(fee))

	msg := &types.MsgFlashLoan{
		Sender: user,
		Denom:  denomA,
		Amount: borrowAmt,
	}

	resp, err := k.FlashLoan(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, fee, resp.Fee)
}

func TestMsgCreateFlashPoolViaServer(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	msg := &types.MsgCreateFlashPool{
		Authority: authority,
		Denom:     denomA,
		FeeRate:   types.DefaultFeeRate,
	}

	_, err := k.CreateFlashPool(ctx, msg)
	require.NoError(t, err)

	pool, found := k.GetPool(ctx, denomA)
	require.True(t, found)
	require.Equal(t, denomA, pool.Denom)
}

func TestMsgFundFlashPoolViaServer(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	err := k.CreatePool(ctx, authority, denomA, types.DefaultFeeRate)
	require.NoError(t, err)

	user := fundUser(bk, 1, denomA, math.NewInt(5_000_000))

	msg := &types.MsgFundFlashPool{
		Sender: user,
		Denom:  denomA,
		Amount: math.NewInt(3_000_000),
	}

	_, err = k.FundFlashPool(ctx, msg)
	require.NoError(t, err)

	pool, _ := k.GetPool(ctx, denomA)
	require.Equal(t, math.NewInt(3_000_000), pool.Available)
}

func TestFlashLoanZeroAmount(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	err := k.CreatePool(ctx, authority, denomA, types.DefaultFeeRate)
	require.NoError(t, err)

	msg := &types.MsgFlashLoan{
		Sender: testAddr(1),
		Denom:  denomA,
		Amount: math.ZeroInt(),
	}

	err = msg.ValidateBasic()
	require.Error(t, err)
}

func TestFundInsufficientBalance(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	err := k.CreatePool(ctx, authority, denomA, types.DefaultFeeRate)
	require.NoError(t, err)

	// User only has 100 but tries to fund 1000
	user := fundUser(bk, 1, denomA, math.NewInt(100))
	err = k.FundPool(ctx, user, denomA, math.NewInt(1000))
	require.Error(t, err)
}

func TestGetAllPools(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	err := k.CreatePool(ctx, authority, denomA, types.DefaultFeeRate)
	require.NoError(t, err)
	err = k.CreatePool(ctx, authority, denomB, types.DefaultFeeRate)
	require.NoError(t, err)

	pools := k.GetAllPools(ctx)
	require.Len(t, pools, 2)
}

func TestDefaultStats(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	stats := k.GetStats(ctx)
	require.Equal(t, uint64(0), stats.TotalExecuted)
	require.True(t, stats.TotalVolume.IsZero())
	require.True(t, stats.TotalFeesEarned.IsZero())
}

func TestFlashLoanInactivePool(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	err := k.CreatePool(ctx, authority, denomA, types.DefaultFeeRate)
	require.NoError(t, err)

	// Deactivate pool manually
	pool, _ := k.GetPool(ctx, denomA)
	pool.Active = false
	k.SetPool(ctx, pool)

	user := fundUser(bk, 1, denomA, math.NewInt(10_000_000))
	_, err = k.ExecuteFlashLoan(ctx, user, denomA, math.NewInt(1_000_000))
	require.ErrorIs(t, err, types.ErrPoolNotActive)
}

func TestFundInactivePool(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	err := k.CreatePool(ctx, authority, denomA, types.DefaultFeeRate)
	require.NoError(t, err)

	pool, _ := k.GetPool(ctx, denomA)
	pool.Active = false
	k.SetPool(ctx, pool)

	user := fundUser(bk, 1, denomA, math.NewInt(1_000_000))
	err = k.FundPool(ctx, user, denomA, math.NewInt(500_000))
	require.ErrorIs(t, err, types.ErrPoolNotActive)
}
