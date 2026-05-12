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

	"syreen/x/predict/keeper"
	"syreen/x/predict/types"

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

// govAddr returns the governance module address used as the chain authority in tests
func govAddr() string {
	return authtypes.NewModuleAddress("gov").String()
}

const quoteDenom = "uusdc"

func fundUser(bk *mockBankKeeper, userIdx int, denom string, amount math.Int) string {
	addr := testAddr(userIdx)
	bk.setBalance(addr, denom, amount)
	return addr
}

func fundModule(bk *mockBankKeeper, denom string, amount math.Int) {
	moduleAddr := authtypes.NewModuleAddress(types.ModuleName).String()
	bk.setBalance(moduleAddr, denom, amount)
}

func createTestMarket(t *testing.T, k keeper.Keeper, ctx sdk.Context, bk *mockBankKeeper, creatorIdx int, resolverIdx int, liquidity int64) (uint64, string, string) {
	t.Helper()
	creator := fundUser(bk, creatorIdx, quoteDenom, math.NewInt(liquidity*10))
	resolver := testAddr(resolverIdx)
	marketID, err := k.ExecuteCreateMarket(ctx, creator, "Will BTC hit $100K?", resolver, quoteDenom, 1000, math.NewInt(liquidity))
	require.NoError(t, err)
	return marketID, creator, resolver
}

// ============================================================
// Tests
// ============================================================

func TestCreateMarket(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	creator := fundUser(bk, 1, quoteDenom, math.NewInt(100000))
	resolver := testAddr(2)

	marketID, err := k.ExecuteCreateMarket(ctx, creator, "Will BTC hit $100K by June 2026?", resolver, quoteDenom, 1000, math.NewInt(10000))
	require.NoError(t, err)
	require.Equal(t, uint64(1), marketID)

	market, found := k.GetMarket(ctx, marketID)
	require.True(t, found)
	require.Equal(t, "Will BTC hit $100K by June 2026?", market.Question)
	require.Equal(t, creator, market.Creator)
	require.Equal(t, resolver, market.Resolver)
	require.Equal(t, types.MarketStatusOpen, market.Status)
	require.Equal(t, math.NewInt(10000), market.YesShares)
	require.Equal(t, math.NewInt(10000), market.NoShares)
	require.Equal(t, math.NewInt(10000), market.Liquidity)

	// Second market gets ID 2
	marketID2, err := k.ExecuteCreateMarket(ctx, creator, "Will ETH hit $10K?", resolver, quoteDenom, 2000, math.NewInt(5000))
	require.NoError(t, err)
	require.Equal(t, uint64(2), marketID2)
}

func TestBuyYesShares(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	marketID, _, _ := createTestMarket(t, k, ctx, bk, 1, 2, 10000)

	buyer := fundUser(bk, 3, quoteDenom, math.NewInt(50000))

	// Buy YES shares with 5000 quote tokens
	sharesBought, err := k.ExecuteBuyShares(ctx, buyer, marketID, "yes", math.NewInt(5000))
	require.NoError(t, err)
	require.True(t, sharesBought.IsPositive())

	// YES price should have gone up (YES shares decreased, NO shares increased)
	market, _ := k.GetMarket(ctx, marketID)
	total := market.YesShares.Add(market.NoShares)
	yesPriceBP := market.NoShares.Mul(math.NewInt(10000)).Quo(total)
	require.True(t, yesPriceBP.GT(math.NewInt(5000)), "YES price should be above 50%% after buying YES")

	// Position recorded
	pos, found := k.GetPosition(ctx, marketID, buyer)
	require.True(t, found)
	require.Equal(t, sharesBought, pos.YesShares)
	require.Equal(t, math.NewInt(5000), pos.TotalSpent)
}

func TestBuyNoShares(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	marketID, _, _ := createTestMarket(t, k, ctx, bk, 1, 2, 10000)

	buyer := fundUser(bk, 3, quoteDenom, math.NewInt(50000))

	sharesBought, err := k.ExecuteBuyShares(ctx, buyer, marketID, "no", math.NewInt(3000))
	require.NoError(t, err)
	require.True(t, sharesBought.IsPositive())

	// NO price should have gone up (NO shares decreased)
	market, _ := k.GetMarket(ctx, marketID)
	total := market.YesShares.Add(market.NoShares)
	noPriceBP := market.YesShares.Mul(math.NewInt(10000)).Quo(total)
	require.True(t, noPriceBP.GT(math.NewInt(5000)), "NO price should be above 50%% after buying NO")

	pos, found := k.GetPosition(ctx, marketID, buyer)
	require.True(t, found)
	require.Equal(t, sharesBought, pos.NoShares)
}

func TestPricesSumToOne(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	marketID, _, _ := createTestMarket(t, k, ctx, bk, 1, 2, 10000)

	// Initial prices should be 50/50
	yesPrice, ok := k.GetYesPrice(ctx, marketID)
	require.True(t, ok)
	noPrice, ok := k.GetNoPrice(ctx, marketID)
	require.True(t, ok)
	require.Equal(t, math.NewInt(10000), yesPrice.Add(noPrice))

	// Buy some YES shares
	buyer := fundUser(bk, 3, quoteDenom, math.NewInt(50000))
	_, err := k.ExecuteBuyShares(ctx, buyer, marketID, "yes", math.NewInt(3000))
	require.NoError(t, err)

	// Prices should still sum to ~10000 (within rounding)
	yesPrice, _ = k.GetYesPrice(ctx, marketID)
	noPrice, _ = k.GetNoPrice(ctx, marketID)
	sum := yesPrice.Add(noPrice)
	require.True(t, sum.GTE(math.NewInt(9999)) && sum.LTE(math.NewInt(10001)),
		"prices should sum to ~10000 basis points, got %s", sum)
}

func TestSellShares(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	marketID, _, _ := createTestMarket(t, k, ctx, bk, 1, 2, 10000)

	buyer := fundUser(bk, 3, quoteDenom, math.NewInt(50000))

	sharesBought, err := k.ExecuteBuyShares(ctx, buyer, marketID, "yes", math.NewInt(5000))
	require.NoError(t, err)

	// Sell half the shares back
	halfShares := sharesBought.Quo(math.NewInt(2))
	quoteReturned, err := k.ExecuteSellShares(ctx, buyer, marketID, "yes", halfShares)
	require.NoError(t, err)
	require.True(t, quoteReturned.IsPositive())

	// Position should be reduced
	pos, found := k.GetPosition(ctx, marketID, buyer)
	require.True(t, found)
	require.Equal(t, sharesBought.Sub(halfShares), pos.YesShares)
}

func TestSellInsufficientShares(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	marketID, _, _ := createTestMarket(t, k, ctx, bk, 1, 2, 10000)

	buyer := fundUser(bk, 3, quoteDenom, math.NewInt(50000))

	sharesBought, err := k.ExecuteBuyShares(ctx, buyer, marketID, "yes", math.NewInt(5000))
	require.NoError(t, err)

	// Try to sell more than owned
	_, err = k.ExecuteSellShares(ctx, buyer, marketID, "yes", sharesBought.Add(math.NewInt(1)))
	require.ErrorIs(t, err, types.ErrInsufficientShares)
}

func TestResolveMarketYes(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	marketID, _, _ := createTestMarket(t, k, ctx, bk, 1, 2, 10000)

	// Buy YES and NO shares
	yesBuyer := fundUser(bk, 3, quoteDenom, math.NewInt(50000))
	noBuyer := fundUser(bk, 4, quoteDenom, math.NewInt(50000))

	yesShares, err := k.ExecuteBuyShares(ctx, yesBuyer, marketID, "yes", math.NewInt(5000))
	require.NoError(t, err)
	_, err = k.ExecuteBuyShares(ctx, noBuyer, marketID, "no", math.NewInt(3000))
	require.NoError(t, err)

	// Resolve YES — only chain governance authority may resolve
	err = k.ExecuteResolveMarket(ctx, govAddr(), marketID, "yes")
	require.NoError(t, err)

	market, _ := k.GetMarket(ctx, marketID)
	require.Equal(t, types.MarketStatusResolvedYes, market.Status)

	res, found := k.GetResolution(ctx, marketID)
	require.True(t, found)
	require.Equal(t, "yes", res.Outcome)

	// YES holder claims winnings = their YES shares
	fundModule(bk, quoteDenom, math.NewInt(100000)) // ensure module has funds
	payout, err := k.ExecuteClaimWinnings(ctx, yesBuyer, marketID)
	require.NoError(t, err)
	require.Equal(t, yesShares, payout)
}

func TestResolveMarketNo(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	marketID, _, _ := createTestMarket(t, k, ctx, bk, 1, 2, 10000)

	yesBuyer := fundUser(bk, 3, quoteDenom, math.NewInt(50000))
	noBuyer := fundUser(bk, 4, quoteDenom, math.NewInt(50000))

	_, err := k.ExecuteBuyShares(ctx, yesBuyer, marketID, "yes", math.NewInt(5000))
	require.NoError(t, err)
	noShares, err := k.ExecuteBuyShares(ctx, noBuyer, marketID, "no", math.NewInt(3000))
	require.NoError(t, err)

	// Resolve NO — only chain governance authority may resolve
	err = k.ExecuteResolveMarket(ctx, govAddr(), marketID, "no")
	require.NoError(t, err)

	market, _ := k.GetMarket(ctx, marketID)
	require.Equal(t, types.MarketStatusResolvedNo, market.Status)

	// NO holder claims winnings
	fundModule(bk, quoteDenom, math.NewInt(100000))
	payout, err := k.ExecuteClaimWinnings(ctx, noBuyer, marketID)
	require.NoError(t, err)
	require.Equal(t, noShares, payout)

	// YES holder should get nothing (no YES winnings when resolved NO)
	_, err = k.ExecuteClaimWinnings(ctx, yesBuyer, marketID)
	require.ErrorIs(t, err, types.ErrNoWinnings)
}

func TestVoidMarket(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	marketID, _, _ := createTestMarket(t, k, ctx, bk, 1, 2, 10000)

	buyer := fundUser(bk, 3, quoteDenom, math.NewInt(50000))
	_, err := k.ExecuteBuyShares(ctx, buyer, marketID, "yes", math.NewInt(5000))
	require.NoError(t, err)

	// Void the market — only chain governance authority may resolve
	err = k.ExecuteResolveMarket(ctx, govAddr(), marketID, "void")
	require.NoError(t, err)

	market, _ := k.GetMarket(ctx, marketID)
	require.Equal(t, types.MarketStatusVoided, market.Status)

	// Buyer gets full refund of TotalSpent
	fundModule(bk, quoteDenom, math.NewInt(100000))
	payout, err := k.ExecuteClaimWinnings(ctx, buyer, marketID)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(5000), payout)
}

func TestCannotTradeAfterClosed(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	marketID, _, _ := createTestMarket(t, k, ctx, bk, 1, 2, 10000)

	// Resolve the market to close it — only chain governance authority may resolve
	err := k.ExecuteResolveMarket(ctx, govAddr(), marketID, "yes")
	require.NoError(t, err)

	buyer := fundUser(bk, 3, quoteDenom, math.NewInt(50000))
	_, err = k.ExecuteBuyShares(ctx, buyer, marketID, "yes", math.NewInt(1000))
	require.ErrorIs(t, err, types.ErrMarketNotOpen)

	_, err = k.ExecuteSellShares(ctx, buyer, marketID, "yes", math.NewInt(1000))
	require.ErrorIs(t, err, types.ErrMarketNotOpen)
}

func TestCannotResolveIfNotAuthority(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	marketID, _, _ := createTestMarket(t, k, ctx, bk, 1, 2, 10000)

	// Any address that is not the chain governance authority must be rejected
	nonAuthority := testAddr(99)
	err := k.ExecuteResolveMarket(ctx, nonAuthority, marketID, "yes")
	require.ErrorIs(t, err, types.ErrUnauthorized)
}

func TestCannotResolveAlreadyResolved(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	marketID, _, _ := createTestMarket(t, k, ctx, bk, 1, 2, 10000)

	err := k.ExecuteResolveMarket(ctx, govAddr(), marketID, "yes")
	require.NoError(t, err)

	err = k.ExecuteResolveMarket(ctx, govAddr(), marketID, "no")
	require.ErrorIs(t, err, types.ErrAlreadyResolved)
}

func TestAutoCloseAtResolutionBlock(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	creator := fundUser(bk, 1, quoteDenom, math.NewInt(100000))
	resolver := testAddr(2)

	// Market with resolution at block 50
	marketID, err := k.ExecuteCreateMarket(ctx, creator, "test?", resolver, quoteDenom, 50, math.NewInt(10000))
	require.NoError(t, err)

	// At block 1, market is open
	market, _ := k.GetMarket(ctx, marketID)
	require.Equal(t, types.MarketStatusOpen, market.Status)

	// Simulate block 50
	ctx = ctx.WithBlockHeight(50)
	k.AutoCloseExpiredMarkets(ctx)

	market, _ = k.GetMarket(ctx, marketID)
	require.Equal(t, types.MarketStatusClosed, market.Status)

	// Cannot trade on closed market
	buyer := fundUser(bk, 3, quoteDenom, math.NewInt(50000))
	_, err = k.ExecuteBuyShares(ctx, buyer, marketID, "yes", math.NewInt(1000))
	require.ErrorIs(t, err, types.ErrMarketNotOpen)

	// Governance authority can still resolve a closed market
	err = k.ExecuteResolveMarket(ctx, govAddr(), marketID, "yes")
	require.NoError(t, err)
}

func TestDoubleClaim(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	marketID, _, _ := createTestMarket(t, k, ctx, bk, 1, 2, 10000)

	buyer := fundUser(bk, 3, quoteDenom, math.NewInt(50000))
	_, err := k.ExecuteBuyShares(ctx, buyer, marketID, "yes", math.NewInt(5000))
	require.NoError(t, err)

	err = k.ExecuteResolveMarket(ctx, govAddr(), marketID, "yes")
	require.NoError(t, err)

	fundModule(bk, quoteDenom, math.NewInt(100000))
	_, err = k.ExecuteClaimWinnings(ctx, buyer, marketID)
	require.NoError(t, err)

	// Second claim should fail
	_, err = k.ExecuteClaimWinnings(ctx, buyer, marketID)
	require.ErrorIs(t, err, types.ErrAlreadyClaimed)
}

func TestBuyIncreasesVolume(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	marketID, _, _ := createTestMarket(t, k, ctx, bk, 1, 2, 10000)

	buyer := fundUser(bk, 3, quoteDenom, math.NewInt(50000))
	_, err := k.ExecuteBuyShares(ctx, buyer, marketID, "yes", math.NewInt(3000))
	require.NoError(t, err)
	_, err = k.ExecuteBuyShares(ctx, buyer, marketID, "no", math.NewInt(2000))
	require.NoError(t, err)

	market, _ := k.GetMarket(ctx, marketID)
	require.Equal(t, math.NewInt(5000), market.TotalVolume)
}

func TestMultipleBuyersPriceMovement(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	marketID, _, _ := createTestMarket(t, k, ctx, bk, 1, 2, 100000)

	// Initial price 50/50
	yesPrice1, _ := k.GetYesPrice(ctx, marketID)
	require.Equal(t, math.NewInt(5000), yesPrice1)

	// Buyer1 buys YES — price goes up
	buyer1 := fundUser(bk, 3, quoteDenom, math.NewInt(500000))
	_, err := k.ExecuteBuyShares(ctx, buyer1, marketID, "yes", math.NewInt(50000))
	require.NoError(t, err)
	yesPrice2, _ := k.GetYesPrice(ctx, marketID)
	require.True(t, yesPrice2.GT(yesPrice1))

	// Buyer2 also buys YES — price goes up even more
	buyer2 := fundUser(bk, 4, quoteDenom, math.NewInt(500000))
	_, err = k.ExecuteBuyShares(ctx, buyer2, marketID, "yes", math.NewInt(50000))
	require.NoError(t, err)
	yesPrice3, _ := k.GetYesPrice(ctx, marketID)
	require.True(t, yesPrice3.GT(yesPrice2))

	// Buyer3 buys NO — pushes YES price down
	buyer3 := fundUser(bk, 5, quoteDenom, math.NewInt(500000))
	_, err = k.ExecuteBuyShares(ctx, buyer3, marketID, "no", math.NewInt(100000))
	require.NoError(t, err)
	yesPrice4, _ := k.GetYesPrice(ctx, marketID)
	require.True(t, yesPrice4.LT(yesPrice3))
}

func TestMarketNotFound(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	buyer := fundUser(bk, 3, quoteDenom, math.NewInt(50000))

	_, err := k.ExecuteBuyShares(ctx, buyer, 999, "yes", math.NewInt(1000))
	require.ErrorIs(t, err, types.ErrMarketNotFound)
}

func TestInsufficientFundsOnBuy(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	marketID, _, _ := createTestMarket(t, k, ctx, bk, 1, 2, 10000)

	// Fund with less than they try to spend
	buyer := fundUser(bk, 3, quoteDenom, math.NewInt(100))
	_, err := k.ExecuteBuyShares(ctx, buyer, marketID, "yes", math.NewInt(5000))
	require.Error(t, err)
}

func TestNoPositionClaim(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	marketID, _, _ := createTestMarket(t, k, ctx, bk, 1, 2, 10000)

	err := k.ExecuteResolveMarket(ctx, govAddr(), marketID, "yes")
	require.NoError(t, err)

	// User with no position tries to claim
	stranger := testAddr(99)
	_, err = k.ExecuteClaimWinnings(ctx, stranger, marketID)
	require.ErrorIs(t, err, types.ErrNoPosition)
}

func TestMsgServerCreateMarket(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	creator := fundUser(bk, 1, quoteDenom, math.NewInt(100000))
	resolver := testAddr(2)

	resp, err := k.CreateMarket(ctx, &types.MsgCreateMarket{
		Creator:          creator,
		Question:         "Will SOL hit $500?",
		Resolver:         resolver,
		QuoteDenom:       quoteDenom,
		ResolutionBlock:  500,
		InitialLiquidity: math.NewInt(10000),
	})
	require.NoError(t, err)
	require.Equal(t, uint64(1), resp.MarketID)
}

func TestMsgServerBuyAndSell(t *testing.T) {
	k, ctx, bk := setupKeeper(t)
	marketID, _, _ := createTestMarket(t, k, ctx, bk, 1, 2, 10000)

	buyer := fundUser(bk, 3, quoteDenom, math.NewInt(50000))

	buyResp, err := k.BuyShares(ctx, &types.MsgBuyShares{
		Sender:   buyer,
		MarketID: marketID,
		Outcome:  "yes",
		Amount:   math.NewInt(5000),
	})
	require.NoError(t, err)
	require.True(t, buyResp.SharesBought.IsPositive())
	require.True(t, buyResp.AvgPrice.IsPositive())

	sellResp, err := k.SellShares(ctx, &types.MsgSellShares{
		Sender:   buyer,
		MarketID: marketID,
		Outcome:  "yes",
		Shares:   buyResp.SharesBought.Quo(math.NewInt(2)),
	})
	require.NoError(t, err)
	require.True(t, sellResp.QuoteReturned.IsPositive())
}
