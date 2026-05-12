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

	"syreen/x/portfolio/keeper"
	"syreen/x/portfolio/types"

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

func (m *mockBankKeeper) GetAllBalances(_ context.Context, addr sdk.AccAddress) sdk.Coins {
	addrStr := addr.String()
	if m.balances[addrStr] == nil {
		return sdk.Coins{}
	}
	var coins sdk.Coins
	for denom, amt := range m.balances[addrStr] {
		if amt.IsPositive() {
			coins = append(coins, sdk.NewCoin(denom, amt))
		}
	}
	return coins.Sort()
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

func fundUser(bk *mockBankKeeper, userIdx int, denom string, amount math.Int) string {
	addr := testAddr(userIdx)
	bk.setBalance(addr, denom, amount)
	return addr
}

func fundModule(bk *mockBankKeeper, denom string, amount math.Int) {
	moduleAddr := authtypes.NewModuleAddress(types.ModuleName).String()
	bk.setBalance(moduleAddr, denom, amount)
}

const denomA = "usyreen"
const denomB = "uusdc"

// ============================================================
// Portfolio Tests
// ============================================================

func TestRecordTrade_UpdatesPortfolio(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	addr := testAddr(1)

	tradeID, err := k.RecordTradeAndUpdatePortfolio(ctx, addr, types.TradeTypeSwap, denomA, math.NewInt(1000), math.NewInt(1), math.NewInt(50), 1)
	require.NoError(t, err)
	require.Equal(t, uint64(1), tradeID)

	// Check trade record
	tr, found := k.GetTradeRecord(ctx, 1)
	require.True(t, found)
	require.Equal(t, addr, tr.Address)
	require.Equal(t, types.TradeTypeSwap, tr.TradeType)
	require.Equal(t, denomA, tr.Denom)
	require.Equal(t, math.NewInt(1000), tr.Amount)

	// Check portfolio updated
	p, found := k.GetPortfolio(ctx, addr)
	require.True(t, found)
	require.Equal(t, math.NewInt(50), p.RealizedPnL)
	require.Equal(t, int64(1), p.LastUpdated)
}

func TestRecordMultipleTrades_CumulativePnL(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	addr := testAddr(1)

	_, err := k.RecordTradeAndUpdatePortfolio(ctx, addr, types.TradeTypeSwap, denomA, math.NewInt(1000), math.NewInt(1), math.NewInt(50), 1)
	require.NoError(t, err)

	_, err = k.RecordTradeAndUpdatePortfolio(ctx, addr, types.TradeTypeSwap, denomB, math.NewInt(500), math.NewInt(2), math.NewInt(-20), 2)
	require.NoError(t, err)

	p, found := k.GetPortfolio(ctx, addr)
	require.True(t, found)
	require.Equal(t, math.NewInt(30), p.RealizedPnL) // 50 + (-20) = 30
}

func TestPortfolioAsset_TracksAmount(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	addr := testAddr(1)

	// Buy 1000
	_, err := k.RecordTradeAndUpdatePortfolio(ctx, addr, types.TradeTypeSwap, denomA, math.NewInt(1000), math.NewInt(1), math.ZeroInt(), 1)
	require.NoError(t, err)

	asset, found := k.GetPortfolioAsset(ctx, addr, denomA)
	require.True(t, found)
	require.Equal(t, math.NewInt(1000), asset.Amount)

	// Sell 400
	_, err = k.RecordTradeAndUpdatePortfolio(ctx, addr, types.TradeTypePerpClose, denomA, math.NewInt(400), math.NewInt(1), math.ZeroInt(), 2)
	require.NoError(t, err)

	asset, found = k.GetPortfolioAsset(ctx, addr, denomA)
	require.True(t, found)
	require.Equal(t, math.NewInt(600), asset.Amount)
}

func TestPortfolioAsset_AllAssets(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	addr := testAddr(1)

	_, _ = k.RecordTradeAndUpdatePortfolio(ctx, addr, types.TradeTypeSwap, denomA, math.NewInt(1000), math.NewInt(1), math.ZeroInt(), 1)
	_, _ = k.RecordTradeAndUpdatePortfolio(ctx, addr, types.TradeTypeSwap, denomB, math.NewInt(500), math.NewInt(2), math.ZeroInt(), 1)

	assets := k.GetAllPortfolioAssets(ctx, addr)
	require.Len(t, assets, 2)
}

func TestRecalculatePortfolio(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	addr := fundUser(bk, 1, denomA, math.NewInt(5000))
	bk.setBalance(addr, denomB, math.NewInt(3000))

	totalValue := k.RecalculatePortfolio(ctx, addr)
	require.Equal(t, math.NewInt(8000), totalValue)

	p, found := k.GetPortfolio(ctx, addr)
	require.True(t, found)
	require.Equal(t, math.NewInt(8000), p.TotalValue)
}

// ============================================================
// Activity Feed Tests
// ============================================================

func TestRecordActivity(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	addr := testAddr(1)

	id := k.RecordActivity(ctx, types.ActivityLargeSwap, addr, "Swapped 50000 usyreen", math.NewInt(50000), denomA, 10)
	require.Equal(t, uint64(1), id)

	activity, found := k.GetActivity(ctx, 1)
	require.True(t, found)
	require.Equal(t, types.ActivityLargeSwap, activity.ActivityType)
	require.Equal(t, addr, activity.Address)
	require.Equal(t, int64(10), activity.Block)
}

func TestGlobalActivityList_TrimsToMax(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	addr := testAddr(1)

	// Record more than MaxGlobalActivities
	for i := 0; i < types.MaxGlobalActivities+10; i++ {
		k.RecordActivity(ctx, types.ActivityLargeSwap, addr, "test", math.NewInt(100), denomA, int64(i))
	}

	gal := k.GetGlobalActivityList(ctx)
	require.Len(t, gal.ActivityIDs, types.MaxGlobalActivities)
	// First ID should be 11 (trimmed first 10)
	require.Equal(t, uint64(11), gal.ActivityIDs[0])
}

func TestWhaleDetection_LargeSwap(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	addr := testAddr(1)

	// Whale threshold is 10,000,000,000 usyreen. Amount * price >= threshold triggers whale activity.
	// Amount=10_000_000_000 * price=1 = 10_000_000_000 -> whale
	_, err := k.RecordTradeAndUpdatePortfolio(ctx, addr, types.TradeTypeSwap, denomA, math.NewInt(10_000_000_000), math.NewInt(1), math.ZeroInt(), 1)
	require.NoError(t, err)

	gal := k.GetGlobalActivityList(ctx)
	require.Len(t, gal.ActivityIDs, 1)

	activity, found := k.GetActivity(ctx, gal.ActivityIDs[0])
	require.True(t, found)
	require.Equal(t, types.ActivityLargeSwap, activity.ActivityType)
}

func TestWhaleDetection_SmallSwapNoActivity(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	addr := testAddr(1)

	// Small trade: 100 * 1 = 100 < whale threshold
	_, err := k.RecordTradeAndUpdatePortfolio(ctx, addr, types.TradeTypeSwap, denomA, math.NewInt(100), math.NewInt(1), math.ZeroInt(), 1)
	require.NoError(t, err)

	gal := k.GetGlobalActivityList(ctx)
	require.Len(t, gal.ActivityIDs, 0)
}

func TestUserActivities(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	addr1 := testAddr(1)
	addr2 := testAddr(2)

	k.RecordActivity(ctx, types.ActivityLargeSwap, addr1, "swap1", math.NewInt(100), denomA, 1)
	k.RecordActivity(ctx, types.ActivityLargeDeposit, addr2, "deposit1", math.NewInt(200), denomA, 2)
	k.RecordActivity(ctx, types.ActivityPositionOpened, addr1, "open1", math.NewInt(300), denomA, 3)

	acts := k.GetUserActivities(ctx, addr1)
	require.Len(t, acts, 2)
	require.Equal(t, types.ActivityLargeSwap, acts[0].ActivityType)
	require.Equal(t, types.ActivityPositionOpened, acts[1].ActivityType)
}

// ============================================================
// Competition Tests
// ============================================================

func TestCreateCompetition(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	creator := fundUser(bk, 1, denomA, math.NewInt(1_000_000))

	compID, err := k.ExecuteCreateCompetition(ctx, creator, "Test Comp", 5, 100, denomA, math.NewInt(100_000), math.NewInt(1000), 50)
	require.NoError(t, err)
	require.Equal(t, uint64(1), compID)

	comp, found := k.GetCompetition(ctx, 1)
	require.True(t, found)
	require.Equal(t, "Test Comp", comp.Name)
	require.Equal(t, creator, comp.Creator)
	require.Equal(t, math.NewInt(100_000), comp.PrizePool)
	require.Equal(t, types.CompStatusUpcoming, comp.Status)
}

func TestCreateCompetition_InsufficientFunds(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	creator := fundUser(bk, 1, denomA, math.NewInt(100)) // not enough

	_, err := k.ExecuteCreateCompetition(ctx, creator, "Test", 5, 100, denomA, math.NewInt(100_000), math.NewInt(1000), 50)
	require.Error(t, err)
}

func TestJoinCompetition(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	creator := fundUser(bk, 1, denomA, math.NewInt(1_000_000))
	joiner := fundUser(bk, 2, denomA, math.NewInt(50_000))

	// Create with entry fee of 1000
	compID, err := k.ExecuteCreateCompetition(ctx, creator, "Test", 1, 100, denomA, math.NewInt(10_000), math.NewInt(1000), 10)
	require.NoError(t, err)

	err = k.ExecuteJoinCompetition(ctx, joiner, compID)
	require.NoError(t, err)

	// Check entry
	entry, found := k.GetCompetitionEntry(ctx, compID, joiner)
	require.True(t, found)
	require.Equal(t, joiner, entry.Address)

	// Check prize pool increased by entry fee
	comp, _ := k.GetCompetition(ctx, compID)
	require.Equal(t, math.NewInt(11_000), comp.PrizePool) // 10000 + 1000
	require.Equal(t, uint64(1), comp.Participants)
}

func TestJoinCompetition_AlreadyJoined(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	creator := fundUser(bk, 1, denomA, math.NewInt(1_000_000))
	joiner := fundUser(bk, 2, denomA, math.NewInt(50_000))

	compID, _ := k.ExecuteCreateCompetition(ctx, creator, "Test", 1, 100, denomA, math.ZeroInt(), math.ZeroInt(), 10)
	_ = k.ExecuteJoinCompetition(ctx, joiner, compID)
	err := k.ExecuteJoinCompetition(ctx, joiner, compID)
	require.ErrorIs(t, err, types.ErrAlreadyJoined)
}

func TestJoinCompetition_Full(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	creator := fundUser(bk, 1, denomA, math.NewInt(1_000_000))

	compID, _ := k.ExecuteCreateCompetition(ctx, creator, "Test", 1, 100, denomA, math.ZeroInt(), math.ZeroInt(), 1)
	joiner1 := fundUser(bk, 2, denomA, math.NewInt(50_000))
	joiner2 := fundUser(bk, 3, denomA, math.NewInt(50_000))

	err := k.ExecuteJoinCompetition(ctx, joiner1, compID)
	require.NoError(t, err)

	err = k.ExecuteJoinCompetition(ctx, joiner2, compID)
	require.ErrorIs(t, err, types.ErrCompetitionFull)
}

func TestEndCompetition_DistributePrizes(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	creator := fundUser(bk, 1, denomA, math.NewInt(1_000_000))

	compID, _ := k.ExecuteCreateCompetition(ctx, creator, "Test", 1, 50, denomA, math.NewInt(100_000), math.ZeroInt(), 10)

	// Join 3 users
	user1 := fundUser(bk, 10, denomA, math.NewInt(50_000))
	user2 := fundUser(bk, 11, denomA, math.NewInt(50_000))
	user3 := fundUser(bk, 12, denomA, math.NewInt(50_000))

	_ = k.ExecuteJoinCompetition(ctx, user1, compID)
	_ = k.ExecuteJoinCompetition(ctx, user2, compID)
	_ = k.ExecuteJoinCompetition(ctx, user3, compID)

	// Give user1 the best portfolio value
	p1 := types.Portfolio{Address: user1, TotalValue: math.NewInt(200_000), TotalPnL: math.NewInt(100), RealizedPnL: math.NewInt(100), UnrealizedPnL: math.ZeroInt(), LastUpdated: 10}
	k.SetPortfolio(ctx, p1)
	// user2 second best
	p2 := types.Portfolio{Address: user2, TotalValue: math.NewInt(150_000), TotalPnL: math.NewInt(50), RealizedPnL: math.NewInt(50), UnrealizedPnL: math.ZeroInt(), LastUpdated: 10}
	k.SetPortfolio(ctx, p2)
	// user3 worst
	p3 := types.Portfolio{Address: user3, TotalValue: math.NewInt(80_000), TotalPnL: math.NewInt(-10), RealizedPnL: math.NewInt(-10), UnrealizedPnL: math.ZeroInt(), LastUpdated: 10}
	k.SetPortfolio(ctx, p3)

	winners, err := k.EndCompetitionAndDistribute(ctx, compID)
	require.NoError(t, err)
	require.Len(t, winners, 3)

	// Verify competition is distributed
	comp, _ := k.GetCompetition(ctx, compID)
	require.Equal(t, types.CompStatusDistributed, comp.Status)

	// Winner 1 gets 50% = 50000
	user1Bal := bk.getBalance(user1, denomA)
	require.Equal(t, math.NewInt(100_000), user1Bal) // 50000 original + 50000 prize

	// Winner 2 gets 30% = 30000
	user2Bal := bk.getBalance(user2, denomA)
	require.Equal(t, math.NewInt(80_000), user2Bal) // 50000 + 30000

	// Winner 3 gets 20% = 20000
	user3Bal := bk.getBalance(user3, denomA)
	require.Equal(t, math.NewInt(70_000), user3Bal) // 50000 + 20000
}

func TestEndCompetition_AlreadyDistributed(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	creator := fundUser(bk, 1, denomA, math.NewInt(1_000_000))

	compID, _ := k.ExecuteCreateCompetition(ctx, creator, "Test", 1, 50, denomA, math.ZeroInt(), math.ZeroInt(), 10)

	_, err := k.EndCompetitionAndDistribute(ctx, compID)
	require.NoError(t, err) // first time: ok (no entries, just marks distributed)

	_, err = k.EndCompetitionAndDistribute(ctx, compID)
	require.ErrorIs(t, err, types.ErrPrizesAlreadyDistributed)
}

func TestEndCompetition_NotFound(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	_, err := k.EndCompetitionAndDistribute(ctx, 999)
	require.ErrorIs(t, err, types.ErrCompetitionNotFound)
}

func TestBeginBlock_AutoEndCompetition(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	creator := fundUser(bk, 1, denomA, math.NewInt(1_000_000))

	// Create competition ending at block 5
	compID, _ := k.ExecuteCreateCompetition(ctx, creator, "Auto End", 1, 5, denomA, math.ZeroInt(), math.ZeroInt(), 10)

	// Advance to block 5
	ctx = ctx.WithBlockHeight(5)
	k.ProcessBeginBlock(ctx)

	comp, _ := k.GetCompetition(ctx, compID)
	require.Equal(t, types.CompStatusDistributed, comp.Status)
}

// ============================================================
// MsgServer Tests
// ============================================================

func TestMsgRecordTrade(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	// L-3 fix: MsgRecordTrade is now restricted to the dex and intent
	// module accounts (and the chain authority). A trade submitted by
	// the dex module address must succeed.
	dexAddr := authtypes.NewModuleAddress("dex").String()

	resp, err := k.RecordTrade(ctx, &types.MsgRecordTrade{
		Sender:    dexAddr,
		TradeType: types.TradeTypeSwap,
		Denom:     denomA,
		Amount:    math.NewInt(500),
		Price:     math.NewInt(2),
		PnL:       math.NewInt(10),
	})
	require.NoError(t, err)
	require.Equal(t, uint64(1), resp.TradeID)
}

// TestMsgRecordTrade_UnauthorizedUser verifies that an arbitrary user cannot
// submit MsgRecordTrade and thereby spoof the leaderboard (L-3).
func TestMsgRecordTrade_UnauthorizedUser(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	addr := testAddr(1)

	_, err := k.RecordTrade(ctx, &types.MsgRecordTrade{
		Sender:    addr,
		TradeType: types.TradeTypeSwap,
		Denom:     denomA,
		Amount:    math.NewInt(500),
		Price:     math.NewInt(2),
		PnL:       math.NewInt(10),
	})
	require.ErrorIs(t, err, types.ErrUnauthorized)
}

func TestMsgUpdatePortfolio(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	addr := fundUser(bk, 1, denomA, math.NewInt(10_000))

	resp, err := k.UpdatePortfolio(ctx, &types.MsgUpdatePortfolio{Sender: addr})
	require.NoError(t, err)
	require.Equal(t, math.NewInt(10_000), resp.TotalValue)
}

func TestPnLCalculation_RealisedUpdates(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	addr := testAddr(1)

	// Trade 1: +100 PnL
	_, _ = k.RecordTradeAndUpdatePortfolio(ctx, addr, types.TradeTypeSwap, denomA, math.NewInt(100), math.NewInt(1), math.NewInt(100), 1)
	// Trade 2: -30 PnL
	_, _ = k.RecordTradeAndUpdatePortfolio(ctx, addr, types.TradeTypeSwap, denomA, math.NewInt(200), math.NewInt(1), math.NewInt(-30), 2)
	// Trade 3: +50 PnL
	_, _ = k.RecordTradeAndUpdatePortfolio(ctx, addr, types.TradeTypePerpClose, denomA, math.NewInt(50), math.NewInt(1), math.NewInt(50), 3)

	p, found := k.GetPortfolio(ctx, addr)
	require.True(t, found)
	require.Equal(t, math.NewInt(120), p.RealizedPnL) // 100 - 30 + 50
}

func TestPortfolioPercentOfTotal(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	addr := fundUser(bk, 1, denomA, math.NewInt(7000))
	bk.setBalance(addr, denomB, math.NewInt(3000))

	k.RecalculatePortfolio(ctx, addr)

	assetA, found := k.GetPortfolioAsset(ctx, addr, denomA)
	require.True(t, found)
	// 7000/10000 * 100 = 70%
	require.True(t, assetA.PercentOfTotal.GT(math.LegacyNewDec(69)))
	require.True(t, assetA.PercentOfTotal.LT(math.LegacyNewDec(71)))

	assetB, found := k.GetPortfolioAsset(ctx, addr, denomB)
	require.True(t, found)
	// 3000/10000 * 100 = 30%
	require.True(t, assetB.PercentOfTotal.GT(math.LegacyNewDec(29)))
	require.True(t, assetB.PercentOfTotal.LT(math.LegacyNewDec(31)))
}

func TestCompetitionSingleWinner(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	creator := fundUser(bk, 1, denomA, math.NewInt(1_000_000))

	compID, _ := k.ExecuteCreateCompetition(ctx, creator, "Solo", 1, 50, denomA, math.NewInt(10_000), math.ZeroInt(), 10)

	user1 := fundUser(bk, 10, denomA, math.NewInt(50_000))
	_ = k.ExecuteJoinCompetition(ctx, user1, compID)

	p1 := types.Portfolio{Address: user1, TotalValue: math.NewInt(100_000), TotalPnL: math.NewInt(50), RealizedPnL: math.NewInt(50), UnrealizedPnL: math.ZeroInt(), LastUpdated: 10}
	k.SetPortfolio(ctx, p1)

	winners, err := k.EndCompetitionAndDistribute(ctx, compID)
	require.NoError(t, err)
	require.Len(t, winners, 1)

	// Single winner gets 100% = 10000
	user1Bal := bk.getBalance(user1, denomA)
	require.Equal(t, math.NewInt(60_000), user1Bal) // 50000 + 10000
}
