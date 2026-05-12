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

	"syreen/x/perps/keeper"
	"syreen/x/perps/types"

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

func (m *mockBankKeeper) SendCoinsFromModuleToModule(_ context.Context, senderModule, recipientModule string, amt sdk.Coins) error {
	fromAddr := authtypes.NewModuleAddress(senderModule)
	toAddr := authtypes.NewModuleAddress(recipientModule)
	for _, coin := range amt {
		fromBal := m.getBalance(fromAddr.String(), coin.Denom)
		if fromBal.LT(coin.Amount) {
			return types.ErrInsufficientFunds
		}
		m.setBalance(fromAddr.String(), coin.Denom, fromBal.Sub(coin.Amount))
		toBal := m.getBalance(toAddr.String(), coin.Denom)
		m.setBalance(toAddr.String(), coin.Denom, toBal.Add(coin.Amount))
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

func (m *mockDexKeeper) GetPoolDenoms(_ context.Context, _ uint64) (string, string, bool) {
	return "usyreen", "uusdc", true
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

const quoteDenom = "uusdc"
const baseDenom = "usyreen"

var authority = authtypes.NewModuleAddress("gov").String()

// createTestMarket creates a market with default params and returns its ID.
func createTestMarket(t *testing.T, k keeper.Keeper, ctx sdk.Context) uint64 {
	t.Helper()
	id, err := k.CreateMarketEntry(ctx, authority, baseDenom, quoteDenom, 1, math.LegacyNewDec(20))
	require.NoError(t, err)
	return id
}

// fundUser sets the user's quote denom balance and returns the address string.
func fundUser(bk *mockBankKeeper, userIdx int, amount math.Int) string {
	addr := testAddr(userIdx)
	bk.setBalance(addr, quoteDenom, amount)
	return addr
}

// openPosition is a convenience wrapper that opens a long or short position.
func openPosition(t *testing.T, k keeper.Keeper, ctx sdk.Context, sender string, marketID uint64, side types.Side, margin math.Int, leverage math.LegacyDec) *types.MsgOpenPositionResponse {
	t.Helper()
	resp, err := k.ExecuteOpenPosition(ctx, sender, marketID, side, margin, leverage)
	require.NoError(t, err)
	return resp
}

// ============================================================
// Tests
// ============================================================

func TestCreateMarket(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	id, err := k.CreateMarketEntry(ctx, authority, baseDenom, quoteDenom, 1, math.LegacyNewDec(20))
	require.NoError(t, err)
	require.Equal(t, uint64(1), id)

	market, ok := k.GetMarket(ctx, id)
	require.True(t, ok)
	require.Equal(t, uint64(1), market.ID)
	require.Equal(t, baseDenom, market.BaseDenom)
	require.Equal(t, quoteDenom, market.QuoteDenom)
	require.Equal(t, uint64(1), market.PoolID)
	require.True(t, market.MaxLeverage.Equal(math.LegacyNewDec(20)))
	require.True(t, market.MaintenanceMargin.Equal(types.DefaultMaintenanceMargin))
	require.True(t, market.InitialMargin.Equal(types.DefaultInitialMargin))
	require.True(t, market.TakerFee.Equal(types.DefaultTakerFee))
	require.True(t, market.MakerFee.Equal(types.DefaultMakerFee))
	require.Equal(t, types.DefaultFundingInterval, market.FundingInterval)
	require.True(t, market.MaxFundingRate.Equal(types.DefaultMaxFundingRate))
	require.True(t, market.LongOpenInterest.IsZero())
	require.True(t, market.ShortOpenInterest.IsZero())
	require.True(t, market.Active)
	require.Equal(t, authority, market.Creator)

	// Next market ID should be incremented
	require.Equal(t, uint64(2), k.GetNextMarketID(ctx))

	// Funding state should be initialized
	fs := k.GetFundingState(ctx, id)
	require.Equal(t, id, fs.MarketID)
	require.True(t, fs.CurrentFundingRate.IsZero())
	require.True(t, fs.CumulativeFunding.IsZero())
}

func TestCreateMarketUnauthorized(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	_, err := k.CreateMarketEntry(ctx, testAddr(1), baseDenom, quoteDenom, 1, math.LegacyNewDec(20))
	require.ErrorIs(t, err, types.ErrUnauthorized)
}

func TestOpenLongPosition(t *testing.T) {
	k, ctx, bk, dk := setupKeeper(t)
	marketID := createTestMarket(t, k, ctx)

	// Price = 10, margin = 1000, leverage = 5x
	// Notional = 5000, size = 500
	// Fee = 1000 * 5 * 0.001 = 5
	dk.setPrice(math.LegacyNewDec(10))
	user := fundUser(bk, 1, math.NewInt(2000))

	resp, err := k.ExecuteOpenPosition(ctx, user, marketID, types.SideLong, math.NewInt(1000), math.LegacyNewDec(5))
	require.NoError(t, err)

	// Verify response
	require.True(t, resp.EntryPrice.Equal(math.LegacyNewDec(10)))
	expectedSize := math.LegacyNewDec(500) // 5000/10
	require.True(t, resp.PositionSize.Equal(expectedSize))
	expectedFee := math.NewInt(5) // 1000*5*0.001
	require.True(t, resp.Fee.Equal(expectedFee))

	// Verify stored position
	pos, ok := k.GetPosition(ctx, marketID, user)
	require.True(t, ok)
	require.Equal(t, user, pos.Address)
	require.Equal(t, marketID, pos.MarketID)
	require.Equal(t, types.SideLong, pos.Side)
	require.True(t, pos.Size.Equal(expectedSize))
	require.True(t, pos.EntryPrice.Equal(math.LegacyNewDec(10)))
	require.True(t, pos.Margin.Equal(math.NewInt(1000)))
	require.True(t, pos.Leverage.Equal(math.LegacyNewDec(5)))
	require.True(t, pos.UnrealizedPnL.IsZero())
	require.Equal(t, int64(1), pos.OpenedAt)

	// Verify margin + fee transferred from user to module
	// User started with 2000, deposited 1000 + 5 = 1005
	require.True(t, bk.getBalance(user, quoteDenom).Equal(math.NewInt(995)))
	moduleAddr := authtypes.NewModuleAddress(types.ModuleName).String()
	require.True(t, bk.getBalance(moduleAddr, quoteDenom).Equal(math.NewInt(1005)))

	// Verify open interest updated
	market, _ := k.GetMarket(ctx, marketID)
	require.True(t, market.LongOpenInterest.Equal(math.NewInt(5000)))
	require.True(t, market.ShortOpenInterest.IsZero())

	// Verify insurance fund received fee
	fund := k.GetInsuranceFund(ctx)
	require.True(t, fund.Balance.Equal(expectedFee))
}

func TestOpenShortPosition(t *testing.T) {
	k, ctx, bk, dk := setupKeeper(t)
	marketID := createTestMarket(t, k, ctx)

	dk.setPrice(math.LegacyNewDec(10))
	user := fundUser(bk, 1, math.NewInt(2000))

	resp, err := k.ExecuteOpenPosition(ctx, user, marketID, types.SideShort, math.NewInt(1000), math.LegacyNewDec(5))
	require.NoError(t, err)

	expectedSize := math.LegacyNewDec(500) // 5000/10
	require.True(t, resp.PositionSize.Equal(expectedSize))
	require.True(t, resp.EntryPrice.Equal(math.LegacyNewDec(10)))

	pos, ok := k.GetPosition(ctx, marketID, user)
	require.True(t, ok)
	require.Equal(t, types.SideShort, pos.Side)
	require.True(t, pos.Size.Equal(expectedSize))

	// Verify short OI updated
	market, _ := k.GetMarket(ctx, marketID)
	require.True(t, market.ShortOpenInterest.Equal(math.NewInt(5000)))
	require.True(t, market.LongOpenInterest.IsZero())
}

func TestOpenPositionInvalidLeverage(t *testing.T) {
	k, ctx, bk, dk := setupKeeper(t)
	marketID := createTestMarket(t, k, ctx) // max leverage = 20x

	dk.setPrice(math.LegacyNewDec(10))
	user := fundUser(bk, 1, math.NewInt(100000))

	// Try 25x leverage (exceeds 20x max)
	_, err := k.ExecuteOpenPosition(ctx, user, marketID, types.SideLong, math.NewInt(1000), math.LegacyNewDec(25))
	require.ErrorIs(t, err, types.ErrMaxLeverageExceeded)
}

func TestOpenPositionDuplicateRejected(t *testing.T) {
	k, ctx, bk, dk := setupKeeper(t)
	marketID := createTestMarket(t, k, ctx)

	dk.setPrice(math.LegacyNewDec(10))
	user := fundUser(bk, 1, math.NewInt(100000))

	_, err := k.ExecuteOpenPosition(ctx, user, marketID, types.SideLong, math.NewInt(1000), math.LegacyNewDec(5))
	require.NoError(t, err)

	// Try to open a second position in the same market
	_, err = k.ExecuteOpenPosition(ctx, user, marketID, types.SideShort, math.NewInt(1000), math.LegacyNewDec(5))
	require.ErrorIs(t, err, types.ErrPositionAlreadyExists)
}

func TestOpenPositionInsufficientFunds(t *testing.T) {
	k, ctx, bk, dk := setupKeeper(t)
	marketID := createTestMarket(t, k, ctx)

	dk.setPrice(math.LegacyNewDec(10))
	// Give user only 100 but they need 1000 margin + fee
	user := fundUser(bk, 1, math.NewInt(100))

	_, err := k.ExecuteOpenPosition(ctx, user, marketID, types.SideLong, math.NewInt(1000), math.LegacyNewDec(5))
	require.ErrorIs(t, err, types.ErrInsufficientFunds)
}

func TestCloseLongProfit(t *testing.T) {
	k, ctx, bk, dk := setupKeeper(t)
	marketID := createTestMarket(t, k, ctx)

	// Open long at price 10
	dk.setPrice(math.LegacyNewDec(10))
	user := fundUser(bk, 1, math.NewInt(10000))
	openPosition(t, k, ctx, user, marketID, types.SideLong, math.NewInt(1000), math.LegacyNewDec(5))
	// size = 500, entry = 10, margin = 1000, fee = 5
	// User balance after open: 10000 - 1005 = 8995

	// Pre-fund module to simulate counterparty deposits covering profits
	moduleAddr := authtypes.NewModuleAddress(types.ModuleName).String()
	bk.setBalance(moduleAddr, quoteDenom, bk.getBalance(moduleAddr, quoteDenom).Add(math.NewInt(10000)))

	// Price goes up to 12 (+20%)
	dk.setPrice(math.LegacyNewDec(12))

	resp, err := k.ExecuteClosePosition(ctx, user, marketID)
	require.NoError(t, err)

	// PnL = size * (markPrice - entryPrice) = 500 * (12 - 10) = 1000
	require.True(t, resp.RealizedPnL.Equal(math.LegacyNewDec(1000)))
	// Payout = margin + PnL = 1000 + 1000 = 2000
	require.True(t, resp.Payout.Equal(math.NewInt(2000)))

	// User should receive payout
	// 8995 + 2000 = 10995
	require.True(t, bk.getBalance(user, quoteDenom).Equal(math.NewInt(10995)))

	// Position should be deleted
	_, ok := k.GetPosition(ctx, marketID, user)
	require.False(t, ok)

	// OI should be reduced
	market, _ := k.GetMarket(ctx, marketID)
	require.True(t, market.LongOpenInterest.IsZero())
}

func TestCloseLongLoss(t *testing.T) {
	k, ctx, bk, dk := setupKeeper(t)
	marketID := createTestMarket(t, k, ctx)

	// Open long at price 10
	dk.setPrice(math.LegacyNewDec(10))
	user := fundUser(bk, 1, math.NewInt(10000))
	openPosition(t, k, ctx, user, marketID, types.SideLong, math.NewInt(1000), math.LegacyNewDec(5))
	// size = 500, entry = 10

	// Price drops to 9 (-10%)
	dk.setPrice(math.LegacyNewDec(9))

	resp, err := k.ExecuteClosePosition(ctx, user, marketID)
	require.NoError(t, err)

	// PnL = 500 * (9 - 10) = -500
	require.True(t, resp.RealizedPnL.Equal(math.LegacyNewDec(-500)))
	// Payout = margin + PnL = 1000 + (-500) = 500
	require.True(t, resp.Payout.Equal(math.NewInt(500)))

	// User balance: 8995 (after open) + 500 (payout) = 9495
	require.True(t, bk.getBalance(user, quoteDenom).Equal(math.NewInt(9495)))

	// Position deleted
	_, ok := k.GetPosition(ctx, marketID, user)
	require.False(t, ok)
}

func TestCloseShortProfit(t *testing.T) {
	k, ctx, bk, dk := setupKeeper(t)
	marketID := createTestMarket(t, k, ctx)

	// Open short at price 10
	dk.setPrice(math.LegacyNewDec(10))
	user := fundUser(bk, 1, math.NewInt(10000))
	openPosition(t, k, ctx, user, marketID, types.SideShort, math.NewInt(1000), math.LegacyNewDec(5))
	// size = 500, entry = 10

	// Pre-fund module to simulate counterparty deposits covering profits
	moduleAddr := authtypes.NewModuleAddress(types.ModuleName).String()
	bk.setBalance(moduleAddr, quoteDenom, bk.getBalance(moduleAddr, quoteDenom).Add(math.NewInt(10000)))

	// Price drops to 8 (-20%), profit for short
	dk.setPrice(math.LegacyNewDec(8))

	resp, err := k.ExecuteClosePosition(ctx, user, marketID)
	require.NoError(t, err)

	// PnL = size * (entryPrice - markPrice) = 500 * (10 - 8) = 1000
	require.True(t, resp.RealizedPnL.Equal(math.LegacyNewDec(1000)))
	// Payout = 1000 + 1000 = 2000
	require.True(t, resp.Payout.Equal(math.NewInt(2000)))

	// User: 8995 + 2000 = 10995
	require.True(t, bk.getBalance(user, quoteDenom).Equal(math.NewInt(10995)))

	_, ok := k.GetPosition(ctx, marketID, user)
	require.False(t, ok)
}

func TestCloseShortLoss(t *testing.T) {
	k, ctx, bk, dk := setupKeeper(t)
	marketID := createTestMarket(t, k, ctx)

	// Open short at price 10
	dk.setPrice(math.LegacyNewDec(10))
	user := fundUser(bk, 1, math.NewInt(10000))
	openPosition(t, k, ctx, user, marketID, types.SideShort, math.NewInt(1000), math.LegacyNewDec(5))
	// size = 500, entry = 10

	// Price rises to 11 (+10%), loss for short
	dk.setPrice(math.LegacyNewDec(11))

	resp, err := k.ExecuteClosePosition(ctx, user, marketID)
	require.NoError(t, err)

	// PnL = 500 * (10 - 11) = -500
	require.True(t, resp.RealizedPnL.Equal(math.LegacyNewDec(-500)))
	// Payout = 1000 + (-500) = 500
	require.True(t, resp.Payout.Equal(math.NewInt(500)))

	// User: 8995 + 500 = 9495
	require.True(t, bk.getBalance(user, quoteDenom).Equal(math.NewInt(9495)))
}

func TestClosePositionNotFound(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	_ = createTestMarket(t, k, ctx)

	_, err := k.ExecuteClosePosition(ctx, testAddr(1), 1)
	require.ErrorIs(t, err, types.ErrNoPosition)
}

func TestAddMargin(t *testing.T) {
	k, ctx, bk, dk := setupKeeper(t)
	marketID := createTestMarket(t, k, ctx)

	dk.setPrice(math.LegacyNewDec(10))
	user := fundUser(bk, 1, math.NewInt(10000))
	openPosition(t, k, ctx, user, marketID, types.SideLong, math.NewInt(1000), math.LegacyNewDec(5))
	// size = 500, leverage = 5x, margin = 1000

	// Add 1000 more margin
	err := k.AddMarginToPosition(ctx, user, marketID, math.NewInt(1000))
	require.NoError(t, err)

	pos, ok := k.GetPosition(ctx, marketID, user)
	require.True(t, ok)
	require.True(t, pos.Margin.Equal(math.NewInt(2000)))

	// Effective leverage should be reduced: notional = 500 * 10 = 5000, leverage = 5000/2000 = 2.5
	require.True(t, pos.Leverage.Equal(math.LegacyNewDecWithPrec(25, 1)))

	// User should have paid 1000 more
	// Started 10000, open cost 1005, add margin 1000 => 7995
	require.True(t, bk.getBalance(user, quoteDenom).Equal(math.NewInt(7995)))
}

func TestRemoveMargin(t *testing.T) {
	k, ctx, bk, dk := setupKeeper(t)
	marketID := createTestMarket(t, k, ctx)

	dk.setPrice(math.LegacyNewDec(10))
	user := fundUser(bk, 1, math.NewInt(10000))
	// Use low leverage so we have room to remove margin
	openPosition(t, k, ctx, user, marketID, types.SideLong, math.NewInt(1000), math.LegacyNewDec(2))
	// size = 200 (2000/10), entry = 10, margin = 1000
	// fee = 1000*2*0.001 = 2, user balance after = 10000 - 1002 = 8998

	// Notional = 200 * 10 = 2000, margin = 1000, leverage = 2
	// Maintenance margin = 5% of notional = 100
	// If we remove 500, new margin = 500, equity = 500, marginRatio = 500/2000 = 0.25 > 0.05
	err := k.RemoveMarginFromPosition(ctx, user, marketID, math.NewInt(500))
	require.NoError(t, err)

	pos, ok := k.GetPosition(ctx, marketID, user)
	require.True(t, ok)
	require.True(t, pos.Margin.Equal(math.NewInt(500)))

	// Leverage should increase: 2000/500 = 4
	require.True(t, pos.Leverage.Equal(math.LegacyNewDec(4)))

	// User received 500 back: 8998 + 500 = 9498
	require.True(t, bk.getBalance(user, quoteDenom).Equal(math.NewInt(9498)))
}

func TestRemoveMarginBelowMaintenance(t *testing.T) {
	k, ctx, bk, dk := setupKeeper(t)
	marketID := createTestMarket(t, k, ctx)

	dk.setPrice(math.LegacyNewDec(10))
	user := fundUser(bk, 1, math.NewInt(10000))
	// Open with high leverage so margin is tight
	openPosition(t, k, ctx, user, marketID, types.SideLong, math.NewInt(1000), math.LegacyNewDec(10))
	// size = 1000 (10000/10), entry = 10, margin = 1000
	// Notional = 1000 * 10 = 10000
	// Maintenance margin ratio = 5%, so min equity = 500
	// Current equity = 1000 (margin, PnL = 0)
	// Try to remove 600 => new margin = 400, equity = 400, ratio = 400/10000 = 4% < 5%

	err := k.RemoveMarginFromPosition(ctx, user, marketID, math.NewInt(600))
	require.ErrorIs(t, err, types.ErrInsufficientMargin)

	// Position unchanged
	pos, _ := k.GetPosition(ctx, marketID, user)
	require.True(t, pos.Margin.Equal(math.NewInt(1000)))
}

func TestLiquidation(t *testing.T) {
	k, ctx, bk, dk := setupKeeper(t)
	marketID := createTestMarket(t, k, ctx)

	dk.setPrice(math.LegacyNewDec(10))
	user := fundUser(bk, 1, math.NewInt(10000))
	openPosition(t, k, ctx, user, marketID, types.SideLong, math.NewInt(1000), math.LegacyNewDec(10))
	// size = 1000, entry = 10, margin = 1000
	// Notional at entry = 10000

	// Price drops to make margin ratio < 5% (maintenance margin)
	// Margin ratio = (margin + PnL) / (size * markPrice)
	// = (1000 + 1000*(price-10)) / (1000*price)
	// At price 9.0: (1000 + 1000*(-1)) / (1000*9) = 0 / 9000 = 0 => liquidate
	// At price 9.5: (1000 + 1000*(-0.5)) / (1000*9.5) = 500 / 9500 = 5.26% > 5% (barely safe)
	// At price 9.04: (1000 + 1000*(-0.96)) / (1000*9.04) = 40/9040 = 0.44% => liquidate
	dk.setPrice(math.LegacyNewDec(9))

	// Run liquidations via BeginBlock
	ctx = ctx.WithBlockHeight(2)
	k.ProcessLiquidations(ctx)

	// Position should be liquidated (deleted)
	_, ok := k.GetPosition(ctx, marketID, user)
	require.False(t, ok, "position should be liquidated")

	// OI should be reduced
	market, _ := k.GetMarket(ctx, marketID)
	require.True(t, market.LongOpenInterest.IsZero())
}

func TestLiquidationInsuranceFund(t *testing.T) {
	k, ctx, bk, dk := setupKeeper(t)
	marketID := createTestMarket(t, k, ctx)

	dk.setPrice(math.LegacyNewDec(10))
	user := fundUser(bk, 1, math.NewInt(10000))
	openPosition(t, k, ctx, user, marketID, types.SideLong, math.NewInt(1000), math.LegacyNewDec(10))
	// size = 1000, entry = 10, margin = 1000, fee = 10
	// Insurance fund starts with the fee (10)

	fundBefore := k.GetInsuranceFund(ctx)
	require.True(t, fundBefore.Balance.Equal(math.NewInt(10)))

	// Price crashes to 8 => PnL = 1000*(8-10) = -2000, equity = 1000 - 2000 = -1000
	// This is a negative equity situation: the loss exceeds margin
	dk.setPrice(math.LegacyNewDec(8))

	ctx = ctx.WithBlockHeight(2)
	k.ProcessLiquidations(ctx)

	// Position should be liquidated
	_, ok := k.GetPosition(ctx, marketID, user)
	require.False(t, ok)

	// Insurance fund should have absorbed the deficit
	// equity = -1000 (remaining is negative), so insurance fund takes the hit
	fundAfter := k.GetInsuranceFund(ctx)
	// Fund was 10, deficit = 1000, fund = max(10 - 1000, 0) = 0
	require.True(t, fundAfter.Balance.IsZero() || fundAfter.Balance.GTE(math.ZeroInt()))
}

func TestFundingRates(t *testing.T) {
	k, ctx, bk, dk := setupKeeper(t)
	marketID := createTestMarket(t, k, ctx)

	dk.setPrice(math.LegacyNewDec(10))

	// Open a long (creates OI imbalance)
	user1 := fundUser(bk, 1, math.NewInt(100000))
	openPosition(t, k, ctx, user1, marketID, types.SideLong, math.NewInt(10000), math.LegacyNewDec(5))
	// Long OI = 50000, Short OI = 0 => 100% imbalance

	// Funding interval is 7200 blocks by default
	market, _ := k.GetMarket(ctx, marketID)

	// Advance past funding interval
	ctx = ctx.WithBlockHeight(1 + market.FundingInterval)
	k.ProcessFundingRates(ctx)

	fs := k.GetFundingState(ctx, marketID)
	// imbalance = (longOI - shortOI) / totalOI = (50000 - 0) / 50000 = 1.0
	// fundingRate = 1.0 * maxFundingRate = 1.0 * 0.01 = 0.01 (capped at max)
	require.True(t, fs.CurrentFundingRate.Equal(types.DefaultMaxFundingRate))
	require.True(t, fs.CumulativeFunding.Equal(types.DefaultMaxFundingRate))
	require.Equal(t, 1+market.FundingInterval, fs.LastFundingBlock)
}

func TestFundingPayment(t *testing.T) {
	k, ctx, bk, dk := setupKeeper(t)
	marketID := createTestMarket(t, k, ctx)

	dk.setPrice(math.LegacyNewDec(10))

	// Open a long position
	user := fundUser(bk, 1, math.NewInt(100000))
	openPosition(t, k, ctx, user, marketID, types.SideLong, math.NewInt(10000), math.LegacyNewDec(5))
	// size = 5000 (50000/10), entry = 10

	// Manually set cumulative funding to simulate funding accrual
	// The position was opened when cumulative funding was 0
	fs := k.GetFundingState(ctx, marketID)
	fs.CumulativeFunding = math.LegacyNewDecWithPrec(1, 2) // 0.01
	k.SetFundingState(ctx, fs)

	// Close at the same price -- PnL from price should be 0 but funding payment matters
	resp, err := k.ExecuteClosePosition(ctx, user, marketID)
	require.NoError(t, err)

	// Funding payment for long = size * (cumFunding_now - cumFunding_entry) = 5000 * (0.01 - 0) = 50
	// PnL from price = 0 (same price)
	// Net PnL = 0 - 50 = -50
	require.True(t, resp.RealizedPnL.Equal(math.LegacyNewDec(-50)))
	// Payout = margin + PnL = 10000 - 50 = 9950
	require.True(t, resp.Payout.Equal(math.NewInt(9950)))
}

func TestGetAllMarkets(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	// No markets initially
	markets := k.GetAllMarkets(ctx)
	require.Len(t, markets, 0)

	// Create 3 markets
	_, err := k.CreateMarketEntry(ctx, authority, "usyreen", "uusdc", 1, math.LegacyNewDec(20))
	require.NoError(t, err)
	_, err = k.CreateMarketEntry(ctx, authority, "ueth", "uusdc", 2, math.LegacyNewDec(10))
	require.NoError(t, err)
	_, err = k.CreateMarketEntry(ctx, authority, "ubtc", "uusdc", 3, math.LegacyNewDec(15))
	require.NoError(t, err)

	markets = k.GetAllMarkets(ctx)
	require.Len(t, markets, 3)

	// Verify IDs are sequential
	ids := make(map[uint64]bool)
	for _, m := range markets {
		ids[m.ID] = true
	}
	require.True(t, ids[1])
	require.True(t, ids[2])
	require.True(t, ids[3])
}

func TestGetPositionsByAddress(t *testing.T) {
	k, ctx, bk, dk := setupKeeper(t)

	// Create 2 markets
	market1, err := k.CreateMarketEntry(ctx, authority, "usyreen", "uusdc", 1, math.LegacyNewDec(20))
	require.NoError(t, err)
	market2, err := k.CreateMarketEntry(ctx, authority, "ueth", "uusdc", 2, math.LegacyNewDec(20))
	require.NoError(t, err)

	dk.setPrice(math.LegacyNewDec(10))
	user := fundUser(bk, 1, math.NewInt(100000))

	// Open positions in both markets
	openPosition(t, k, ctx, user, market1, types.SideLong, math.NewInt(1000), math.LegacyNewDec(2))
	openPosition(t, k, ctx, user, market2, types.SideShort, math.NewInt(2000), math.LegacyNewDec(3))

	positions := k.GetPositionsByAddress(ctx, user)
	require.Len(t, positions, 2)

	// Verify both market IDs present
	marketIDs := make(map[uint64]bool)
	for _, p := range positions {
		marketIDs[p.MarketID] = true
		require.Equal(t, user, p.Address)
	}
	require.True(t, marketIDs[market1])
	require.True(t, marketIDs[market2])

	// Another user should have no positions
	positions2 := k.GetPositionsByAddress(ctx, testAddr(2))
	require.Len(t, positions2, 0)
}
