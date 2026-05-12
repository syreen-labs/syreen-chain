package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"syreen/x/options/keeper"
	"syreen/x/options/types"
)

// ---- mock bank ----

type mockBank struct{ balances map[string]sdk.Coins }

func newMockBank() *mockBank { return &mockBank{balances: map[string]sdk.Coins{}} }

func (m *mockBank) credit(addr sdk.AccAddress, coins sdk.Coins) {
	m.balances[addr.String()] = m.balances[addr.String()].Add(coins...)
}

func (m *mockBank) GetBalance(_ context.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	for _, c := range m.balances[addr.String()] {
		if c.Denom == denom {
			return c
		}
	}
	return sdk.NewCoin(denom, math.ZeroInt())
}

func (m *mockBank) SendCoins(_ context.Context, from, to sdk.AccAddress, amt sdk.Coins) error {
	m.balances[from.String()] = m.balances[from.String()].Sub(amt...)
	m.balances[to.String()] = m.balances[to.String()].Add(amt...)
	return nil
}

func (m *mockBank) SendCoinsFromAccountToModule(ctx context.Context, sender sdk.AccAddress, module string, amt sdk.Coins) error {
	return m.SendCoins(ctx, sender, sdk.AccAddress([]byte(module)), amt)
}

func (m *mockBank) SendCoinsFromModuleToAccount(ctx context.Context, module string, recipient sdk.AccAddress, amt sdk.Coins) error {
	return m.SendCoins(ctx, sdk.AccAddress([]byte(module)), recipient, amt)
}

// ---- mock dex ----

type mockDex struct {
	spotPrice  math.LegacyDec
	volatility math.LegacyDec
	denomA     string
	denomB     string
}

func (m *mockDex) GetSpotPrice(_ context.Context, _ uint64, _, _ string) (math.LegacyDec, error) {
	return m.spotPrice, nil
}
func (m *mockDex) GetPoolDenoms(_ context.Context, _ uint64) (string, string, bool) {
	return m.denomA, m.denomB, true
}
func (m *mockDex) GetSignalVolatility(_ context.Context, _ uint64) math.LegacyDec {
	return m.volatility
}

// ---- setup ----

func setupKeeper(t *testing.T) (keeper.Keeper, sdk.Context, *mockBank, *mockDex) {
	t.Helper()
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, cmtproto.Header{Height: 100}, false, log.NewNopLogger())

	ir := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(ir)

	bank := newMockBank()
	dex := &mockDex{
		spotPrice:  math.LegacyNewDec(100),
		volatility: math.LegacyNewDecWithPrec(30, 2),
		denomA:     "usyreen",
		denomB:     "uusdc",
	}

	k := keeper.NewKeeper(cdc, runtime.NewKVStoreService(storeKey), bank, "authority")
	k.SetDexKeeper(dex)
	return k, ctx, bank, dex
}

func testAddr(i int) string {
	return sdk.AccAddress([]byte{byte(i), 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}).String()
}

var (
	writerAddr = testAddr(1)
	buyerAddr  = testAddr(2)
)

// ---- tests ----

func TestCalculatePremium_Call(t *testing.T) {
	k, _, _, _ := setupKeeper(t)
	premium := k.CalculatePremium(math.LegacyNewDec(100), math.LegacyNewDec(110), 630720, math.LegacyNewDecWithPrec(30, 2), types.OptionTypeCall, math.NewInt(1_000_000))
	require.True(t, premium.IsPositive())
}

func TestCalculatePremium_Put(t *testing.T) {
	k, _, _, _ := setupKeeper(t)
	premium := k.CalculatePremium(math.LegacyNewDec(100), math.LegacyNewDec(90), 630720, math.LegacyNewDecWithPrec(30, 2), types.OptionTypePut, math.NewInt(1_000_000))
	require.True(t, premium.IsPositive())
}

func TestCalculatePremium_ITM_HigherThanOTM(t *testing.T) {
	k, _, _, _ := setupKeeper(t)
	amount := math.NewInt(1_000_000)
	sigma := math.LegacyNewDecWithPrec(30, 2)
	blocks := int64(630720)
	spot := math.LegacyNewDec(110)
	itm := k.CalculatePremium(spot, math.LegacyNewDec(100), blocks, sigma, types.OptionTypeCall, amount)
	otm := k.CalculatePremium(spot, math.LegacyNewDec(120), blocks, sigma, types.OptionTypeCall, amount)
	require.True(t, itm.GT(otm))
}

func TestCalculatePremium_ZeroExpiry(t *testing.T) {
	k, _, _, _ := setupKeeper(t)
	premium := k.CalculatePremium(math.LegacyNewDec(100), math.LegacyNewDec(100), 0, math.LegacyNewDecWithPrec(30, 2), types.OptionTypeCall, math.NewInt(1_000_000))
	require.True(t, premium.IsZero())
}

func TestWriteOption(t *testing.T) {
	k, ctx, bank, _ := setupKeeper(t)
	writer, _ := sdk.AccAddressFromBech32(writerAddr)
	bank.credit(writer, sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(10_000_000))))

	id, premium, err := k.ExecuteWriteOption(ctx, writerAddr, 1, types.OptionTypeCall, math.LegacyNewDec(110), math.NewInt(1_000_000), 200, math.ZeroInt())
	require.NoError(t, err)
	require.Equal(t, uint64(1), id)
	require.True(t, premium.IsPositive())

	opt, found := k.GetOption(ctx, id)
	require.True(t, found)
	require.Equal(t, types.OptionStatusOpen, opt.Status)
}

func TestBuyOption(t *testing.T) {
	k, ctx, bank, _ := setupKeeper(t)
	writer, _ := sdk.AccAddressFromBech32(writerAddr)
	buyer, _ := sdk.AccAddressFromBech32(buyerAddr)
	bank.credit(writer, sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(10_000_000))))

	id, premium, err := k.ExecuteWriteOption(ctx, writerAddr, 1, types.OptionTypeCall, math.LegacyNewDec(110), math.NewInt(1_000_000), 200, math.ZeroInt())
	require.NoError(t, err)
	bank.credit(buyer, sdk.NewCoins(sdk.NewCoin("uusdc", premium)))

	require.NoError(t, k.ExecuteBuyOption(ctx, buyerAddr, id))
	opt, _ := k.GetOption(ctx, id)
	require.Equal(t, types.OptionStatusActive, opt.Status)
	require.Equal(t, buyerAddr, opt.Buyer)
}

func TestExerciseCall_InTheMoney(t *testing.T) {
	k, ctx, bank, dex := setupKeeper(t)
	writer, _ := sdk.AccAddressFromBech32(writerAddr)
	buyer, _ := sdk.AccAddressFromBech32(buyerAddr)
	bank.credit(writer, sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(10_000_000))))

	dex.spotPrice = math.LegacyNewDec(100)
	id, premium, err := k.ExecuteWriteOption(ctx, writerAddr, 1, types.OptionTypeCall, math.LegacyNewDec(90), math.NewInt(1_000_000), 200, math.ZeroInt())
	require.NoError(t, err)
	bank.credit(buyer, sdk.NewCoins(sdk.NewCoin("uusdc", premium)))
	require.NoError(t, k.ExecuteBuyOption(ctx, buyerAddr, id))

	payout, err := k.ExecuteExerciseOption(ctx, buyerAddr, id)
	require.NoError(t, err)
	require.True(t, payout.IsPositive())

	opt, _ := k.GetOption(ctx, id)
	require.Equal(t, types.OptionStatusExercised, opt.Status)
}

func TestExerciseCall_OutOfMoney(t *testing.T) {
	k, ctx, bank, dex := setupKeeper(t)
	writer, _ := sdk.AccAddressFromBech32(writerAddr)
	buyer, _ := sdk.AccAddressFromBech32(buyerAddr)
	bank.credit(writer, sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(10_000_000))))

	dex.spotPrice = math.LegacyNewDec(100)
	id, premium, err := k.ExecuteWriteOption(ctx, writerAddr, 1, types.OptionTypeCall, math.LegacyNewDec(120), math.NewInt(1_000_000), 200, math.ZeroInt())
	require.NoError(t, err)
	bank.credit(buyer, sdk.NewCoins(sdk.NewCoin("uusdc", premium)))
	require.NoError(t, k.ExecuteBuyOption(ctx, buyerAddr, id))

	_, err = k.ExecuteExerciseOption(ctx, buyerAddr, id)
	require.Error(t, err)
}

func TestExercisePut_InTheMoney(t *testing.T) {
	k, ctx, bank, dex := setupKeeper(t)
	writer, _ := sdk.AccAddressFromBech32(writerAddr)
	buyer, _ := sdk.AccAddressFromBech32(buyerAddr)
	bank.credit(writer, sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(200_000_000))))

	dex.spotPrice = math.LegacyNewDec(100)
	id, premium, err := k.ExecuteWriteOption(ctx, writerAddr, 1, types.OptionTypePut, math.LegacyNewDec(110), math.NewInt(1_000_000), 200, math.ZeroInt())
	require.NoError(t, err)
	bank.credit(buyer, sdk.NewCoins(sdk.NewCoin("uusdc", premium)))
	require.NoError(t, k.ExecuteBuyOption(ctx, buyerAddr, id))

	payout, err := k.ExecuteExerciseOption(ctx, buyerAddr, id)
	require.NoError(t, err)
	require.True(t, payout.IsPositive())
}

func TestCancelOption(t *testing.T) {
	k, ctx, bank, _ := setupKeeper(t)
	writer, _ := sdk.AccAddressFromBech32(writerAddr)
	bank.credit(writer, sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(10_000_000))))
	balBefore := bank.GetBalance(ctx, writer, "usyreen")

	id, _, err := k.ExecuteWriteOption(ctx, writerAddr, 1, types.OptionTypeCall, math.LegacyNewDec(110), math.NewInt(1_000_000), 200, math.ZeroInt())
	require.NoError(t, err)

	require.NoError(t, k.ExecuteCancelOption(ctx, writerAddr, id))

	balAfter := bank.GetBalance(ctx, writer, "usyreen")
	require.Equal(t, balBefore, balAfter)

	opt, _ := k.GetOption(ctx, id)
	require.Equal(t, types.OptionStatusCancelled, opt.Status)
}

func TestCancelOption_Unauthorized(t *testing.T) {
	k, ctx, bank, _ := setupKeeper(t)
	writer, _ := sdk.AccAddressFromBech32(writerAddr)
	bank.credit(writer, sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(10_000_000))))

	id, _, err := k.ExecuteWriteOption(ctx, writerAddr, 1, types.OptionTypeCall, math.LegacyNewDec(110), math.NewInt(1_000_000), 200, math.ZeroInt())
	require.NoError(t, err)

	err = k.ExecuteCancelOption(ctx, buyerAddr, id)
	require.Error(t, err)
}

func TestExpireOptions_BeginBlock(t *testing.T) {
	k, ctx, bank, _ := setupKeeper(t)
	writer, _ := sdk.AccAddressFromBech32(writerAddr)
	bank.credit(writer, sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(10_000_000))))
	balBefore := bank.GetBalance(ctx, writer, "usyreen")

	// write option expiring at block 150 (current ctx is block 100)
	id, _, err := k.ExecuteWriteOption(ctx, writerAddr, 1, types.OptionTypeCall, math.LegacyNewDec(110), math.NewInt(1_000_000), 150, math.ZeroInt())
	require.NoError(t, err)

	// advance context to expiry block
	ctx = sdk.NewContext(ctx.MultiStore(), cmtproto.Header{Height: 150}, false, log.NewNopLogger())
	k.BeginBlockHandler(ctx)

	opt, _ := k.GetOption(ctx, id)
	require.Equal(t, types.OptionStatusExpired, opt.Status)
	balAfter := bank.GetBalance(ctx, writer, "usyreen")
	require.Equal(t, balBefore, balAfter)
}

func TestGetAllOptions(t *testing.T) {
	k, ctx, bank, _ := setupKeeper(t)
	writer, _ := sdk.AccAddressFromBech32(writerAddr)
	bank.credit(writer, sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(100_000_000))))

	for i := 0; i < 5; i++ {
		_, _, err := k.ExecuteWriteOption(ctx, writerAddr, 1, types.OptionTypeCall, math.LegacyNewDec(110), math.NewInt(1_000_000), 200, math.ZeroInt())
		require.NoError(t, err)
	}
	require.Len(t, k.GetAllOptions(ctx), 5)
}

func TestWritePutOption(t *testing.T) {
	k, ctx, bank, _ := setupKeeper(t)
	writer, _ := sdk.AccAddressFromBech32(writerAddr)
	bank.credit(writer, sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(200_000_000))))

	id, premium, err := k.ExecuteWriteOption(ctx, writerAddr, 1, types.OptionTypePut, math.LegacyNewDec(90), math.NewInt(1_000_000), 200, math.ZeroInt())
	require.NoError(t, err)
	require.Equal(t, uint64(1), id)
	require.True(t, premium.IsPositive())
}

func TestOptionNotFound(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	_, found := k.GetOption(ctx, 999)
	require.False(t, found)
}

func TestNextOptionID(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	require.Equal(t, uint64(1), k.GetNextOptionID(ctx))
	k.SetNextOptionID(ctx, 5)
	require.Equal(t, uint64(5), k.GetNextOptionID(ctx))
}

func TestCustomPremium(t *testing.T) {
	k, ctx, bank, _ := setupKeeper(t)
	writer, _ := sdk.AccAddressFromBech32(writerAddr)
	bank.credit(writer, sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(10_000_000))))

	custom := math.NewInt(500_000)
	_, premium, err := k.ExecuteWriteOption(ctx, writerAddr, 1, types.OptionTypeCall, math.LegacyNewDec(110), math.NewInt(1_000_000), 200, custom)
	require.NoError(t, err)
	require.Equal(t, custom, premium)
}
