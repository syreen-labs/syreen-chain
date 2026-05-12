package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/require"

	"syreen/x/tokenfactory/keeper"
	"syreen/x/tokenfactory/types"
)

func init() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("syreen", "syreenpub")
	config.Seal()
}

// ---------------------------------------------------------------------------
// Mock keepers
// ---------------------------------------------------------------------------

type mockAccountKeeper struct{}

func (m *mockAccountKeeper) GetModuleAddress(moduleName string) sdk.AccAddress {
	return sdk.AccAddress([]byte("module_addr_padded__")) // 20 bytes
}

func (m *mockAccountKeeper) GetModuleAccount(_ context.Context, _ string) sdk.ModuleAccountI {
	return nil
}

type mockBankKeeper struct {
	minted   sdk.Coins
	burned   sdk.Coins
	metadata map[string]banktypes.Metadata
	// track balances per address for send simulation
	balances map[string]sdk.Coins
}

func newMockBankKeeper() *mockBankKeeper {
	return &mockBankKeeper{
		metadata: make(map[string]banktypes.Metadata),
		balances: make(map[string]sdk.Coins),
	}
}

func (m *mockBankKeeper) SendCoins(_ context.Context, _, _ sdk.AccAddress, _ sdk.Coins) error {
	return nil
}

func (m *mockBankKeeper) SendCoinsFromAccountToModule(_ context.Context, _ sdk.AccAddress, _ string, amt sdk.Coins) error {
	return nil
}

func (m *mockBankKeeper) SendCoinsFromModuleToAccount(_ context.Context, _ string, _ sdk.AccAddress, amt sdk.Coins) error {
	return nil
}

func (m *mockBankKeeper) MintCoins(_ context.Context, _ string, amt sdk.Coins) error {
	m.minted = m.minted.Add(amt...)
	return nil
}

func (m *mockBankKeeper) BurnCoins(_ context.Context, _ string, amt sdk.Coins) error {
	m.burned = m.burned.Add(amt...)
	return nil
}

func (m *mockBankKeeper) SetDenomMetaData(_ context.Context, denomMetaData banktypes.Metadata) {
	m.metadata[denomMetaData.Base] = denomMetaData
}

func (m *mockBankKeeper) GetDenomMetaData(_ context.Context, denom string) (banktypes.Metadata, bool) {
	md, ok := m.metadata[denom]
	return md, ok
}

type mockDistrKeeper struct {
	funded sdk.Coins
}

func (m *mockDistrKeeper) FundCommunityPool(_ context.Context, amount sdk.Coins, _ sdk.AccAddress) error {
	m.funded = m.funded.Add(amount...)
	return nil
}

// ---------------------------------------------------------------------------
// Setup helper
// ---------------------------------------------------------------------------

func setupKeeper(t *testing.T) (keeper.Keeper, sdk.Context, *mockBankKeeper, *mockDistrKeeper) {
	t.Helper()

	storeKey := storetypes.NewKVStoreKey("tokenfactory")
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	storeService := runtime.NewKVStoreService(storeKey)

	bk := newMockBankKeeper()
	dk := &mockDistrKeeper{}
	ak := &mockAccountKeeper{}

	k := keeper.NewKeeper(cdc, storeService, ak, bk, dk, "authority")
	return k, ctx, bk, dk
}

func testAddr() string {
	return sdk.AccAddress([]byte("test_address_padded_")).String()
}

func testAddr2() string {
	return sdk.AccAddress([]byte("second_addr_padding_")).String()
}

// ---------------------------------------------------------------------------
// Params
// ---------------------------------------------------------------------------

func TestSetGetParams(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	// Default params returned when nothing set
	p := k.GetParams(ctx)
	require.Equal(t, types.DefaultParams(), p)

	// Custom params roundtrip
	custom := types.Params{DenomCreationFee: sdk.NewCoins(sdk.NewInt64Coin("usyreen", 999))}
	require.NoError(t, k.SetParams(ctx, custom))
	got := k.GetParams(ctx)
	require.Equal(t, custom, got)
}

// ---------------------------------------------------------------------------
// DenomAuthorityMetadata
// ---------------------------------------------------------------------------

func TestSetGetDenomAuthorityMetadata(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	denom := "factory/syreen1abc/mytoken"
	meta := types.DenomAuthorityMetadata{Admin: "syreen1abc"}

	k.SetDenomAuthorityMetadata(ctx, denom, meta)
	got, found := k.GetDenomAuthorityMetadata(ctx, denom)
	require.True(t, found)
	require.Equal(t, meta, got)
}

func TestGetDenomAuthorityMetadata_NotFound(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	_, found := k.GetDenomAuthorityMetadata(ctx, "factory/nobody/nothing")
	require.False(t, found)
}

// ---------------------------------------------------------------------------
// CreateDenom
// ---------------------------------------------------------------------------

func TestCreateDenom(t *testing.T) {
	k, ctx, bk, dk := setupKeeper(t)

	creator := testAddr()
	subdenom := "mytoken"

	denom, err := k.CreateDenom(ctx, creator, subdenom)
	require.NoError(t, err)
	require.Equal(t, types.GetDenom(creator, subdenom), denom)

	// Authority metadata stored
	meta, found := k.GetDenomAuthorityMetadata(ctx, denom)
	require.True(t, found)
	require.Equal(t, creator, meta.Admin)

	// Bank metadata set
	md, ok := bk.GetDenomMetaData(ctx, denom)
	require.True(t, ok)
	require.Equal(t, denom, md.Base)
	require.Equal(t, subdenom, md.Symbol)

	// Creation fee no longer charged (free denom creation)
	_ = dk

	// Creator index
	denoms := k.GetDenomsFromCreator(ctx, creator)
	require.Contains(t, denoms, denom)
}

func TestCreateDenom_Duplicate(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	creator := testAddr()

	_, err := k.CreateDenom(ctx, creator, "dup")
	require.NoError(t, err)

	_, err = k.CreateDenom(ctx, creator, "dup")
	require.ErrorIs(t, err, types.ErrDenomExists)
}

// ---------------------------------------------------------------------------
// Mint
// ---------------------------------------------------------------------------

func TestMint_Authorized(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	denom, err := k.CreateDenom(ctx, creator, "mintable")
	require.NoError(t, err)

	coin := sdk.NewInt64Coin(denom, 5000)
	err = k.Mint(ctx, creator, coin, "")
	require.NoError(t, err)
	require.True(t, bk.minted.AmountOf(denom).GT(math.ZeroInt()))
}

func TestMint_Unauthorized(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	creator := testAddr()
	other := testAddr2()

	denom, err := k.CreateDenom(ctx, creator, "restricted")
	require.NoError(t, err)

	coin := sdk.NewInt64Coin(denom, 100)
	err = k.Mint(ctx, other, coin, "")
	require.ErrorIs(t, err, types.ErrUnauthorized)
}

func TestMint_DenomDoesNotExist(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	coin := sdk.NewInt64Coin("factory/x/nonexistent", 100)
	err := k.Mint(ctx, testAddr(), coin, "")
	require.ErrorIs(t, err, types.ErrDenomDoesNotExist)
}

func TestMint_WithMintTo(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()
	recipient := testAddr2()

	denom, err := k.CreateDenom(ctx, creator, "gift")
	require.NoError(t, err)

	coin := sdk.NewInt64Coin(denom, 1000)
	err = k.Mint(ctx, creator, coin, recipient)
	require.NoError(t, err)
	require.True(t, bk.minted.AmountOf(denom).Equal(math.NewInt(1000)))
}

// ---------------------------------------------------------------------------
// Burn
// ---------------------------------------------------------------------------

func TestBurn_Authorized(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	denom, err := k.CreateDenom(ctx, creator, "burnable")
	require.NoError(t, err)

	coin := sdk.NewInt64Coin(denom, 300)
	err = k.Burn(ctx, creator, coin, "")
	require.NoError(t, err)
	require.True(t, bk.burned.AmountOf(denom).Equal(math.NewInt(300)))
}

func TestBurn_Unauthorized(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	creator := testAddr()
	other := testAddr2()

	denom, err := k.CreateDenom(ctx, creator, "noburn")
	require.NoError(t, err)

	coin := sdk.NewInt64Coin(denom, 100)
	err = k.Burn(ctx, other, coin, "")
	require.ErrorIs(t, err, types.ErrUnauthorized)
}

func TestBurn_DenomDoesNotExist(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	coin := sdk.NewInt64Coin("factory/x/nonexistent", 100)
	err := k.Burn(ctx, testAddr(), coin, "")
	require.ErrorIs(t, err, types.ErrDenomDoesNotExist)
}

// I-07: Admin can no longer burn from another user's address.
func TestBurn_WithBurnFrom_Rejected(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	creator := testAddr()
	target := testAddr2()

	denom, err := k.CreateDenom(ctx, creator, "burnfrom")
	require.NoError(t, err)

	coin := sdk.NewInt64Coin(denom, 200)
	err = k.Burn(ctx, creator, coin, target)
	require.Error(t, err)
	require.Contains(t, err.Error(), "can only burn from your own balance")
}

// I-07: Burning from your own address (burnFrom == sender) still works.
func TestBurn_WithBurnFrom_Self(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	denom, err := k.CreateDenom(ctx, creator, "burnself")
	require.NoError(t, err)

	coin := sdk.NewInt64Coin(denom, 200)
	err = k.Burn(ctx, creator, coin, creator)
	require.NoError(t, err)
	require.True(t, bk.burned.AmountOf(denom).Equal(math.NewInt(200)))
}

// ---------------------------------------------------------------------------
// ChangeAdmin
// ---------------------------------------------------------------------------

func TestChangeAdmin(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	creator := testAddr()
	newAdmin := testAddr2()

	denom, err := k.CreateDenom(ctx, creator, "admintest")
	require.NoError(t, err)

	// Change admin
	err = k.ChangeAdmin(ctx, creator, denom, newAdmin)
	require.NoError(t, err)

	meta, found := k.GetDenomAuthorityMetadata(ctx, denom)
	require.True(t, found)
	require.Equal(t, newAdmin, meta.Admin)

	// Old admin can no longer mint
	coin := sdk.NewInt64Coin(denom, 100)
	err = k.Mint(ctx, creator, coin, "")
	require.ErrorIs(t, err, types.ErrUnauthorized)

	// New admin can mint
	err = k.Mint(ctx, newAdmin, coin, "")
	require.NoError(t, err)
}

func TestChangeAdmin_Unauthorized(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	creator := testAddr()
	other := testAddr2()

	denom, err := k.CreateDenom(ctx, creator, "noadmin")
	require.NoError(t, err)

	err = k.ChangeAdmin(ctx, other, denom, other)
	require.ErrorIs(t, err, types.ErrUnauthorized)
}

func TestChangeAdmin_DenomNotFound(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	err := k.ChangeAdmin(ctx, testAddr(), "factory/x/missing", testAddr2())
	require.ErrorIs(t, err, types.ErrDenomDoesNotExist)
}

// ---------------------------------------------------------------------------
// GetDenomsFromCreator
// ---------------------------------------------------------------------------

func TestGetDenomsFromCreator(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	creator := testAddr()

	// No denoms initially
	denoms := k.GetDenomsFromCreator(ctx, creator)
	require.Empty(t, denoms)

	// Create multiple denoms
	d1, err := k.CreateDenom(ctx, creator, "alpha")
	require.NoError(t, err)
	d2, err := k.CreateDenom(ctx, creator, "beta")
	require.NoError(t, err)
	d3, err := k.CreateDenom(ctx, creator, "gamma")
	require.NoError(t, err)

	denoms = k.GetDenomsFromCreator(ctx, creator)
	require.Len(t, denoms, 3)
	require.Contains(t, denoms, d1)
	require.Contains(t, denoms, d2)
	require.Contains(t, denoms, d3)

	// Different creator has none
	other := testAddr2()
	require.Empty(t, k.GetDenomsFromCreator(ctx, other))
}

// ---------------------------------------------------------------------------
// Genesis InitGenesis / ExportGenesis
// ---------------------------------------------------------------------------

func TestGenesisRoundtrip(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)
	creator := testAddr()

	// Create some state
	_, err := k.CreateDenom(ctx, creator, "one")
	require.NoError(t, err)
	_, err = k.CreateDenom(ctx, creator, "two")
	require.NoError(t, err)

	// Export
	exported := k.ExportGenesis(ctx)
	require.NotNil(t, exported)
	require.Len(t, exported.FactoryDenoms, 2)
	require.Equal(t, k.GetParams(ctx), exported.Params)

	// Re-init into a fresh keeper
	k2, ctx2, _, _ := setupKeeper(t)
	k2.InitGenesis(ctx2, *exported)

	// Verify state matches
	exported2 := k2.ExportGenesis(ctx2)
	require.Equal(t, len(exported.FactoryDenoms), len(exported2.FactoryDenoms))

	for _, gd := range exported.FactoryDenoms {
		meta, found := k2.GetDenomAuthorityMetadata(ctx2, gd.Denom)
		require.True(t, found, "denom %s should exist after InitGenesis", gd.Denom)
		require.Equal(t, gd.AuthorityMetadata.Admin, meta.Admin)
	}

	// Creator index rebuilt
	denoms := k2.GetDenomsFromCreator(ctx2, creator)
	require.Len(t, denoms, 2)
}

func TestInitGenesis_DefaultState(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	defaultGS := types.DefaultGenesis()
	k.InitGenesis(ctx, *defaultGS)

	params := k.GetParams(ctx)
	require.Equal(t, types.DefaultParams(), params)

	exported := k.ExportGenesis(ctx)
	require.Empty(t, exported.FactoryDenoms)
}
