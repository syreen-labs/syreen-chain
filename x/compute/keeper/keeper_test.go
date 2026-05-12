package keeper_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/require"

	"syreen/x/compute/keeper"
	"syreen/x/compute/types"
)

func init() {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount("syreen", "syreenpub")
	cfg.Seal()
}

func validAddr() string {
	pk := secp256k1.GenPrivKey().PubKey()
	return sdk.AccAddress(pk.Address()).String()
}

// ---------------------------------------------------------------------------
// Mock keepers
// ---------------------------------------------------------------------------

type mockAccountKeeper struct {
	accounts map[string]sdk.AccountI
}

func newMockAccountKeeper() *mockAccountKeeper {
	return &mockAccountKeeper{accounts: make(map[string]sdk.AccountI)}
}

func (m *mockAccountKeeper) GetModuleAddress(moduleName string) sdk.AccAddress {
	return authtypes.NewModuleAddress(moduleName)
}

func (m *mockAccountKeeper) GetModuleAccount(_ context.Context, moduleName string) sdk.ModuleAccountI {
	return authtypes.NewEmptyModuleAccount(moduleName, authtypes.Minter, authtypes.Burner, authtypes.Staking)
}

func (m *mockAccountKeeper) GetAccount(_ context.Context, addr sdk.AccAddress) sdk.AccountI {
	return m.accounts[addr.String()]
}

func (m *mockAccountKeeper) NewAccountWithAddress(_ context.Context, addr sdk.AccAddress) sdk.AccountI {
	acc := authtypes.NewBaseAccountWithAddress(addr)
	return acc
}

func (m *mockAccountKeeper) SetAccount(_ context.Context, acc sdk.AccountI) {
	m.accounts[acc.GetAddress().String()] = acc
}

type mockBankKeeper struct {
	balances map[string]sdk.Coins
	failSend bool
}

func newMockBankKeeper() *mockBankKeeper {
	return &mockBankKeeper{balances: make(map[string]sdk.Coins)}
}

func (m *mockBankKeeper) SendCoins(_ context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error {
	if m.failSend {
		return fmt.Errorf("insufficient funds")
	}
	from := fromAddr.String()
	to := toAddr.String()
	fromBal := m.balances[from]
	newFromBal, hasNeg := fromBal.SafeSub(amt...)
	if hasNeg {
		return fmt.Errorf("insufficient funds")
	}
	m.balances[from] = newFromBal
	m.balances[to] = m.balances[to].Add(amt...)
	return nil
}

func (m *mockBankKeeper) SendCoinsFromAccountToModule(_ context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	if m.failSend {
		return fmt.Errorf("insufficient funds")
	}
	return nil
}

func (m *mockBankKeeper) SendCoinsFromModuleToAccount(_ context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	if m.failSend {
		return fmt.Errorf("insufficient funds")
	}
	return nil
}

func (m *mockBankKeeper) SpendableCoins(_ context.Context, addr sdk.AccAddress) sdk.Coins {
	return m.balances[addr.String()]
}

// ---------------------------------------------------------------------------
// Setup
// ---------------------------------------------------------------------------

func setupKeeper(t *testing.T) (keeper.Keeper, sdk.Context, *mockAccountKeeper, *mockBankKeeper) {
	storeKey := storetypes.NewKVStoreKey("compute")
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	storeService := runtime.NewKVStoreService(storeKey)

	ak := newMockAccountKeeper()
	bk := newMockBankKeeper()
	k := keeper.NewKeeper(cdc, storeService, ak, bk, "authority")
	return k, ctx, ak, bk
}

// ---------------------------------------------------------------------------
// Create (StoreCode)
// ---------------------------------------------------------------------------

func TestCreate_Valid(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	creator := validAddr()
	code := []byte{0x00, 0x61, 0x73, 0x6d}
	codeID, err := k.Create(ctx, creator, code, nil)
	require.NoError(t, err)
	require.Equal(t, uint64(1), codeID)

	// Second upload gets ID 2
	codeID2, err := k.Create(ctx, creator, code, nil)
	require.NoError(t, err)
	require.Equal(t, uint64(2), codeID2)
}

func TestCreate_UploadDenied(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	// Set params to Nobody
	params := types.DefaultParams()
	params.CodeUploadAccess = types.AccessConfig{Permission: types.AccessTypeNobody}
	require.NoError(t, k.SetParams(ctx, params))

	_, err := k.Create(ctx, validAddr(), []byte{0x00}, nil)
	require.ErrorIs(t, err, types.ErrUploadDenied)
}

func TestCreate_OnlyAddress_Allowed(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	allowedAddr := validAddr()
	params := types.DefaultParams()
	params.CodeUploadAccess = types.AccessConfig{Permission: types.AccessTypeOnlyAddress, Address: allowedAddr}
	require.NoError(t, k.SetParams(ctx, params))

	codeID, err := k.Create(ctx, allowedAddr, []byte{0x00}, nil)
	require.NoError(t, err)
	require.Equal(t, uint64(1), codeID)
}

func TestCreate_OnlyAddress_Denied(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	allowedAddr := validAddr()
	otherAddr := validAddr()
	params := types.DefaultParams()
	params.CodeUploadAccess = types.AccessConfig{Permission: types.AccessTypeOnlyAddress, Address: allowedAddr}
	require.NoError(t, k.SetParams(ctx, params))

	_, err := k.Create(ctx, otherAddr, []byte{0x00}, nil)
	require.ErrorIs(t, err, types.ErrUploadDenied)
}

func TestCreate_CodeTooLarge(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	bigCode := make([]byte, types.DefaultMaxWasmCodeSize+1)
	_, err := k.Create(ctx, validAddr(), bigCode, nil)
	require.ErrorIs(t, err, types.ErrCodeTooLarge)
}

func TestCreate_WithCustomInstantiateAccess(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	creator := validAddr()
	access := &types.AccessConfig{Permission: types.AccessTypeOnlyAddress, Address: creator}
	codeID, err := k.Create(ctx, creator, []byte{0x00}, access)
	require.NoError(t, err)

	info, found := k.GetCodeInfo(ctx, codeID)
	require.True(t, found)
	require.Equal(t, types.AccessTypeOnlyAddress, info.InstantiatePermission.Permission)
	require.Equal(t, creator, info.InstantiatePermission.Address)
}

// ---------------------------------------------------------------------------
// GetCodeInfo / setCodeInfo roundtrip
// ---------------------------------------------------------------------------

func TestGetCodeInfo_Roundtrip(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	creator := validAddr()
	codeID, err := k.Create(ctx, creator, []byte{0x00, 0x01, 0x02}, nil)
	require.NoError(t, err)

	info, found := k.GetCodeInfo(ctx, codeID)
	require.True(t, found)
	require.Equal(t, codeID, info.CodeID)
	require.Equal(t, creator, info.Creator)
	require.NotEmpty(t, info.CodeHash)
}

func TestGetCodeInfo_NotFound(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	_, found := k.GetCodeInfo(ctx, 999)
	require.False(t, found)
}

// ---------------------------------------------------------------------------
// Instantiate
// ---------------------------------------------------------------------------

func TestInstantiate_Valid(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	creator := validAddr()
	codeID, err := k.Create(ctx, creator, []byte{0x00}, nil)
	require.NoError(t, err)

	contractAddr, err := k.Instantiate(ctx, codeID, creator, creator, "test-label", json.RawMessage(`{"init":"data"}`), nil)
	require.NoError(t, err)
	require.NotEmpty(t, contractAddr)

	// Verify contract info
	info, found := k.GetContractInfo(ctx, contractAddr)
	require.True(t, found)
	require.Equal(t, codeID, info.CodeID)
	require.Equal(t, creator, info.Creator)
	require.Equal(t, creator, info.Admin)
	require.Equal(t, "test-label", info.Label)
}

func TestInstantiate_CodeNotFound(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	_, err := k.Instantiate(ctx, 999, validAddr(), "", "label", json.RawMessage(`{}`), nil)
	require.ErrorIs(t, err, types.ErrCodeNotFound)
}

func TestInstantiate_Denied(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	creator := validAddr()
	access := &types.AccessConfig{Permission: types.AccessTypeNobody}
	codeID, err := k.Create(ctx, creator, []byte{0x00}, access)
	require.NoError(t, err)

	_, err = k.Instantiate(ctx, codeID, validAddr(), "", "label", json.RawMessage(`{}`), nil)
	require.ErrorIs(t, err, types.ErrInstantiateDenied)
}

func TestInstantiate_WithFunds(t *testing.T) {
	k, ctx, _, bk := setupKeeper(t)

	creator := validAddr()
	codeID, err := k.Create(ctx, creator, []byte{0x00}, nil)
	require.NoError(t, err)

	// Give the creator some funds
	creatorAddr, _ := sdk.AccAddressFromBech32(creator)
	bk.balances[creator] = sdk.NewCoins(sdk.NewInt64Coin("usyreen", 5000))

	funds := sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1000))
	contractAddr, err := k.Instantiate(ctx, codeID, creator, "", "funded", json.RawMessage(`{}`), funds)
	require.NoError(t, err)
	require.NotEmpty(t, contractAddr)

	// Verify funds were transferred
	require.Equal(t, sdk.NewInt64Coin("usyreen", 4000).Amount, bk.balances[creator].AmountOf("usyreen"))
	require.Equal(t, sdk.NewInt64Coin("usyreen", 1000).Amount, bk.balances[contractAddr].AmountOf("usyreen"))
	_ = creatorAddr
}

func TestInstantiate_InsufficientFunds(t *testing.T) {
	k, ctx, _, bk := setupKeeper(t)

	creator := validAddr()
	codeID, err := k.Create(ctx, creator, []byte{0x00}, nil)
	require.NoError(t, err)

	// No balance
	bk.balances[creator] = sdk.NewCoins()

	funds := sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1000))
	_, err = k.Instantiate(ctx, codeID, creator, "", "funded", json.RawMessage(`{}`), funds)
	require.Error(t, err)
}

// ---------------------------------------------------------------------------
// GetContractInfo roundtrip
// ---------------------------------------------------------------------------

func TestGetContractInfo_NotFound(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	_, found := k.GetContractInfo(ctx, "syreen1nonexistent")
	require.False(t, found)
}

// ---------------------------------------------------------------------------
// Execute
// ---------------------------------------------------------------------------

func TestExecute_Legacy_Valid(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	creator := validAddr()
	codeID, err := k.Create(ctx, creator, []byte{0x00}, nil)
	require.NoError(t, err)

	contractAddr, err := k.Instantiate(ctx, codeID, creator, creator, "exec-test", json.RawMessage(`{}`), nil)
	require.NoError(t, err)

	resp, err := k.Execute(ctx, contractAddr, creator, json.RawMessage(`{"action":"do_something"}`), nil)
	require.NoError(t, err)
	require.NotNil(t, resp)

	var result map[string]string
	require.NoError(t, json.Unmarshal(resp, &result))
	require.Equal(t, "ok", result["status"])
}

func TestExecute_ContractNotFound(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	_, err := k.Execute(ctx, validAddr(), validAddr(), json.RawMessage(`{}`), nil)
	require.ErrorIs(t, err, types.ErrContractNotFound)
}

func TestExecute_Legacy_AnyoneCanExecute(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	creator := validAddr()
	other := validAddr()

	codeID, err := k.Create(ctx, creator, []byte{0x00}, nil)
	require.NoError(t, err)

	contractAddr, err := k.Instantiate(ctx, codeID, creator, creator, "exec-test", json.RawMessage(`{}`), nil)
	require.NoError(t, err)

	// With VM-based execution, authorization is handled by the contract itself.
	// Legacy contracts (non-VM) allow any sender.
	resp, err := k.Execute(ctx, contractAddr, other, json.RawMessage(`{"key":"value"}`), nil)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestExecute_Legacy_StoresState(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	creator := validAddr()
	codeID, err := k.Create(ctx, creator, []byte{0x00}, nil)
	require.NoError(t, err)

	contractAddr, err := k.Instantiate(ctx, codeID, creator, creator, "state-test", json.RawMessage(`{}`), nil)
	require.NoError(t, err)

	_, err = k.Execute(ctx, contractAddr, creator, json.RawMessage(`{"counter":"42"}`), nil)
	require.NoError(t, err)

	// Verify state was written with sender namespace
	stateKey := fmt.Sprintf("exec/%s/counter", creator)
	state := k.GetContractState(ctx, contractAddr, []byte(stateKey))
	require.NotNil(t, state)
	require.Equal(t, `"42"`, string(state))
}

// ---------------------------------------------------------------------------
// Migrate
// ---------------------------------------------------------------------------

func TestMigrate_Valid(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	admin := validAddr()
	codeID1, err := k.Create(ctx, admin, []byte{0x00}, nil)
	require.NoError(t, err)
	codeID2, err := k.Create(ctx, admin, []byte{0x01}, nil)
	require.NoError(t, err)

	contractAddr, err := k.Instantiate(ctx, codeID1, admin, admin, "migrate-test", json.RawMessage(`{}`), nil)
	require.NoError(t, err)

	resp, err := k.Migrate(ctx, contractAddr, admin, codeID2, json.RawMessage(`{"migrate":"data"}`))
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify code ID was updated
	info, found := k.GetContractInfo(ctx, contractAddr)
	require.True(t, found)
	require.Equal(t, codeID2, info.CodeID)
}

func TestMigrate_NotAdmin(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	admin := validAddr()
	other := validAddr()
	codeID, err := k.Create(ctx, admin, []byte{0x00}, nil)
	require.NoError(t, err)
	codeID2, err := k.Create(ctx, admin, []byte{0x01}, nil)
	require.NoError(t, err)

	contractAddr, err := k.Instantiate(ctx, codeID, admin, admin, "migrate-test", json.RawMessage(`{}`), nil)
	require.NoError(t, err)

	_, err = k.Migrate(ctx, contractAddr, other, codeID2, json.RawMessage(`{}`))
	require.ErrorIs(t, err, types.ErrUnauthorized)
}

func TestMigrate_NoAdmin(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	creator := validAddr()
	codeID, err := k.Create(ctx, creator, []byte{0x00}, nil)
	require.NoError(t, err)
	codeID2, err := k.Create(ctx, creator, []byte{0x01}, nil)
	require.NoError(t, err)

	// No admin
	contractAddr, err := k.Instantiate(ctx, codeID, creator, "", "migrate-test", json.RawMessage(`{}`), nil)
	require.NoError(t, err)

	_, err = k.Migrate(ctx, contractAddr, creator, codeID2, json.RawMessage(`{}`))
	require.ErrorIs(t, err, types.ErrUnauthorized)
}

func TestMigrate_CodeNotFound(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	admin := validAddr()
	codeID, err := k.Create(ctx, admin, []byte{0x00}, nil)
	require.NoError(t, err)

	contractAddr, err := k.Instantiate(ctx, codeID, admin, admin, "migrate-test", json.RawMessage(`{}`), nil)
	require.NoError(t, err)

	_, err = k.Migrate(ctx, contractAddr, admin, 999, json.RawMessage(`{}`))
	require.ErrorIs(t, err, types.ErrCodeNotFound)
}

func TestMigrate_ContractNotFound(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	_, err := k.Migrate(ctx, validAddr(), validAddr(), 1, json.RawMessage(`{}`))
	require.ErrorIs(t, err, types.ErrContractNotFound)
}

// ---------------------------------------------------------------------------
// UpdateAdmin
// ---------------------------------------------------------------------------

func TestUpdateAdmin_Valid(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	admin := validAddr()
	newAdmin := validAddr()
	codeID, err := k.Create(ctx, admin, []byte{0x00}, nil)
	require.NoError(t, err)

	contractAddr, err := k.Instantiate(ctx, codeID, admin, admin, "admin-test", json.RawMessage(`{}`), nil)
	require.NoError(t, err)

	err = k.UpdateAdmin(ctx, contractAddr, admin, newAdmin)
	require.NoError(t, err)

	info, found := k.GetContractInfo(ctx, contractAddr)
	require.True(t, found)
	require.Equal(t, newAdmin, info.Admin)
}

func TestUpdateAdmin_NotAdmin(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	admin := validAddr()
	other := validAddr()
	codeID, err := k.Create(ctx, admin, []byte{0x00}, nil)
	require.NoError(t, err)

	contractAddr, err := k.Instantiate(ctx, codeID, admin, admin, "admin-test", json.RawMessage(`{}`), nil)
	require.NoError(t, err)

	err = k.UpdateAdmin(ctx, contractAddr, other, validAddr())
	require.ErrorIs(t, err, types.ErrUnauthorized)
}

func TestUpdateAdmin_ContractNotFound(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	err := k.UpdateAdmin(ctx, validAddr(), validAddr(), validAddr())
	require.ErrorIs(t, err, types.ErrContractNotFound)
}

// ---------------------------------------------------------------------------
// ContractState
// ---------------------------------------------------------------------------

func TestContractState_Roundtrip(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	contractAddr := "syreen1testcontract"
	key := []byte("my_key")
	value := []byte("my_value")

	k.SetContractState(ctx, contractAddr, key, value)
	got := k.GetContractState(ctx, contractAddr, key)
	require.Equal(t, value, got)
}

func TestContractState_NotFound(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	got := k.GetContractState(ctx, "syreen1missing", []byte("key"))
	require.Nil(t, got)
}

func TestContractState_MultipleKeys(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	addr := "syreen1multi"
	k.SetContractState(ctx, addr, []byte("k1"), []byte("v1"))
	k.SetContractState(ctx, addr, []byte("k2"), []byte("v2"))

	require.Equal(t, []byte("v1"), k.GetContractState(ctx, addr, []byte("k1")))
	require.Equal(t, []byte("v2"), k.GetContractState(ctx, addr, []byte("k2")))
}

// ---------------------------------------------------------------------------
// QueryRaw
// ---------------------------------------------------------------------------

func TestQueryRaw_ReturnsValue(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	creator := validAddr()
	codeID, err := k.Create(ctx, creator, []byte{0x00}, nil)
	require.NoError(t, err)

	contractAddr, err := k.Instantiate(ctx, codeID, creator, creator, "qr-test", json.RawMessage(`{}`), nil)
	require.NoError(t, err)

	// Set some state
	k.SetContractState(ctx, contractAddr, []byte("raw_key"), []byte("raw_value"))

	result := k.QueryRaw(ctx, contractAddr, []byte("raw_key"))
	require.Equal(t, []byte("raw_value"), result)
}

func TestQueryRaw_ContractNotFound(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	result := k.QueryRaw(ctx, validAddr(), []byte("key"))
	require.Nil(t, result)
}

func TestQueryRaw_KeyNotFound(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	creator := validAddr()
	codeID, err := k.Create(ctx, creator, []byte{0x00}, nil)
	require.NoError(t, err)

	contractAddr, err := k.Instantiate(ctx, codeID, creator, creator, "qr-test", json.RawMessage(`{}`), nil)
	require.NoError(t, err)

	result := k.QueryRaw(ctx, contractAddr, []byte("nonexistent"))
	require.Nil(t, result)
}

// ---------------------------------------------------------------------------
// Params
// ---------------------------------------------------------------------------

func TestParams_SetGet_Roundtrip(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	custom := types.Params{
		CodeUploadAccess:             types.AccessConfig{Permission: types.AccessTypeNobody},
		InstantiateDefaultPermission: types.AccessTypeOnlyAddress,
		MaxWasmCodeSize:              500_000,
		MaxContractGas:               5_000_000_000,
	}
	require.NoError(t, k.SetParams(ctx, custom))

	got := k.GetParams(ctx)
	require.Equal(t, custom.CodeUploadAccess.Permission, got.CodeUploadAccess.Permission)
	require.Equal(t, custom.InstantiateDefaultPermission, got.InstantiateDefaultPermission)
	require.Equal(t, custom.MaxWasmCodeSize, got.MaxWasmCodeSize)
	require.Equal(t, custom.MaxContractGas, got.MaxContractGas)
}

func TestParams_DefaultWhenNotSet(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	got := k.GetParams(ctx)
	def := types.DefaultParams()
	require.Equal(t, def.MaxWasmCodeSize, got.MaxWasmCodeSize)
	require.Equal(t, def.MaxContractGas, got.MaxContractGas)
}

// ---------------------------------------------------------------------------
// Genesis InitGenesis / ExportGenesis roundtrip
// ---------------------------------------------------------------------------

func TestGenesis_InitExport_Roundtrip(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	creator := validAddr()
	gs := types.GenesisState{
		Params: types.Params{
			CodeUploadAccess:             types.AccessConfig{Permission: types.AccessTypeEverybody},
			InstantiateDefaultPermission: types.AccessTypeEverybody,
			MaxWasmCodeSize:              600_000,
			MaxContractGas:               8_000_000_000,
		},
		Codes: []types.GenesisCode{
			{
				CodeID: 1,
				CodeInfo: types.CodeInfo{
					CodeID:                1,
					Creator:               creator,
					CodeHash:              []byte{0xab, 0xcd},
					InstantiatePermission: types.AccessConfig{Permission: types.AccessTypeEverybody},
				},
				CodeBytes: []byte{0x00, 0x61, 0x73, 0x6d},
			},
			{
				CodeID: 2,
				CodeInfo: types.CodeInfo{
					CodeID:                2,
					Creator:               creator,
					CodeHash:              []byte{0xef},
					InstantiatePermission: types.AccessConfig{Permission: types.AccessTypeEverybody},
				},
				CodeBytes: []byte{0x01},
			},
		},
		Contracts: []types.GenesisContract{
			{
				ContractAddress: "syreen1contract1",
				ContractInfo: types.ContractInfo{
					Address: "syreen1contract1",
					CodeID:  1,
					Creator: creator,
					Admin:   creator,
					Label:   "contract-1",
				},
				ContractState: []types.Model{
					{Key: []byte("state_key"), Value: []byte("state_value")},
				},
			},
		},
	}

	k.InitGenesis(ctx, gs)

	// Verify params were set
	params := k.GetParams(ctx)
	require.Equal(t, uint64(600_000), params.MaxWasmCodeSize)

	// Verify codes were restored
	info1, found := k.GetCodeInfo(ctx, 1)
	require.True(t, found)
	require.Equal(t, creator, info1.Creator)

	info2, found := k.GetCodeInfo(ctx, 2)
	require.True(t, found)
	require.Equal(t, uint64(2), info2.CodeID)

	// Verify contract was restored
	ci, found := k.GetContractInfo(ctx, "syreen1contract1")
	require.True(t, found)
	require.Equal(t, "contract-1", ci.Label)

	// Verify contract state was restored
	state := k.GetContractState(ctx, "syreen1contract1", []byte("state_key"))
	require.Equal(t, []byte("state_value"), state)

	// Export and verify params come back
	exported := k.ExportGenesis(ctx)
	require.NotNil(t, exported)
	require.Equal(t, uint64(600_000), exported.Params.MaxWasmCodeSize)
}

func TestGenesis_InitSetsNextCodeID(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	creator := validAddr()
	gs := types.GenesisState{
		Params: types.DefaultParams(),
		Codes: []types.GenesisCode{
			{CodeID: 5, CodeInfo: types.CodeInfo{CodeID: 5, Creator: creator}, CodeBytes: []byte{0x00}},
			{CodeID: 10, CodeInfo: types.CodeInfo{CodeID: 10, Creator: creator}, CodeBytes: []byte{0x01}},
		},
	}

	k.InitGenesis(ctx, gs)

	// Next Create should get code ID 11
	codeID, err := k.Create(ctx, creator, []byte{0x02}, nil)
	require.NoError(t, err)
	require.Equal(t, uint64(11), codeID)
}

func TestGenesis_EmptyInit(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	gs := *types.DefaultGenesis()
	k.InitGenesis(ctx, gs)

	exported := k.ExportGenesis(ctx)
	require.NotNil(t, exported)
	require.Equal(t, types.DefaultParams().MaxWasmCodeSize, exported.Params.MaxWasmCodeSize)
}

// ---------------------------------------------------------------------------
// QuerySmart
// ---------------------------------------------------------------------------

func TestQuerySmart_Valid(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	creator := validAddr()
	codeID, err := k.Create(ctx, creator, []byte{0x00}, nil)
	require.NoError(t, err)

	contractAddr, err := k.Instantiate(ctx, codeID, creator, creator, "smart-test", json.RawMessage(`{}`), nil)
	require.NoError(t, err)

	// Execute to store some state
	_, err = k.Execute(ctx, contractAddr, creator, json.RawMessage(`{"counter":"99"}`), nil)
	require.NoError(t, err)

	// Query for a key that uses the exec/ prefix format
	// QuerySmart looks up "exec/<key>" in state - it doesn't include sender namespace
	resp, err := k.QuerySmart(ctx, contractAddr, json.RawMessage(`{"something":{}}`))
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestQuerySmart_ContractNotFound(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	_, err := k.QuerySmart(ctx, validAddr(), json.RawMessage(`{}`))
	require.ErrorIs(t, err, types.ErrContractNotFound)
}

func TestQuerySmart_InvalidJSON(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	creator := validAddr()
	codeID, err := k.Create(ctx, creator, []byte{0x00}, nil)
	require.NoError(t, err)

	contractAddr, err := k.Instantiate(ctx, codeID, creator, creator, "smart-test", json.RawMessage(`{}`), nil)
	require.NoError(t, err)

	_, err = k.QuerySmart(ctx, contractAddr, json.RawMessage(`{bad`))
	require.ErrorIs(t, err, types.ErrInvalidMsg)
}

// ---------------------------------------------------------------------------
// VM-based contract tests
// ---------------------------------------------------------------------------

func counterContractCode() []byte {
	return []byte(`{
		"version": "1.0",
		"state": {
			"count": "0",
			"owner": ""
		},
		"instantiate": {
			"actions": [
				{"op": "set", "key": "owner", "value": "$msg.owner"},
				{"op": "set", "key": "count", "value": "$msg.initial_count"},
				{"op": "response", "data": {"status": "initialized"}}
			]
		},
		"execute": {
			"increment": {
				"actions": [
					{"op": "get", "key": "count", "into": "$count"},
					{"op": "int_add", "left": "$count", "right": "1", "into": "$count"},
					{"op": "set", "key": "count", "value": "$count"},
					{"op": "response", "data": {"new_count": "$count"}}
				]
			},
			"transfer_ownership": {
				"actions": [
					{"op": "get", "key": "owner", "into": "$current_owner"},
					{"op": "require_eq", "left": "$sender", "right": "$current_owner", "error": "unauthorized"},
					{"op": "set", "key": "owner", "value": "$msg.new_owner"},
					{"op": "response", "data": {"new_owner": "$msg.new_owner"}}
				]
			},
			"send_tokens": {
				"actions": [
					{"op": "get", "key": "owner", "into": "$current_owner"},
					{"op": "require_eq", "left": "$sender", "right": "$current_owner", "error": "unauthorized"},
					{"op": "bank_send", "from": "$contract", "to": "$msg.recipient", "amount": "$msg.amount", "denom": "$msg.denom"},
					{"op": "response", "data": {"status": "sent"}}
				]
			}
		},
		"query": {
			"get_count": {
				"actions": [
					{"op": "get", "key": "count", "into": "$count"},
					{"op": "response", "data": {"count": "$count"}}
				]
			},
			"get_owner": {
				"actions": [
					{"op": "get", "key": "owner", "into": "$owner"},
					{"op": "response", "data": {"owner": "$owner"}}
				]
			}
		}
	}`)
}

func TestVM_InstantiateAndQuery(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	creator := validAddr()
	codeID, err := k.Create(ctx, creator, counterContractCode(), nil)
	require.NoError(t, err)

	initMsg := json.RawMessage(fmt.Sprintf(`{"owner":"%s","initial_count":"10"}`, creator))
	contractAddr, err := k.Instantiate(ctx, codeID, creator, creator, "vm-counter", initMsg, nil)
	require.NoError(t, err)
	require.NotEmpty(t, contractAddr)

	// Query the count
	resp, err := k.QuerySmart(ctx, contractAddr, json.RawMessage(`{"get_count":{}}`))
	require.NoError(t, err)

	var result map[string]string
	require.NoError(t, json.Unmarshal(resp, &result))
	require.Equal(t, "10", result["count"])
}

func TestVM_ExecuteIncrement(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	creator := validAddr()
	codeID, err := k.Create(ctx, creator, counterContractCode(), nil)
	require.NoError(t, err)

	initMsg := json.RawMessage(fmt.Sprintf(`{"owner":"%s","initial_count":"0"}`, creator))
	contractAddr, err := k.Instantiate(ctx, codeID, creator, creator, "vm-counter", initMsg, nil)
	require.NoError(t, err)

	// Increment twice
	_, err = k.Execute(ctx, contractAddr, creator, json.RawMessage(`{"increment":{}}`), nil)
	require.NoError(t, err)

	_, err = k.Execute(ctx, contractAddr, creator, json.RawMessage(`{"increment":{}}`), nil)
	require.NoError(t, err)

	// Query count should be 2
	resp, err := k.QuerySmart(ctx, contractAddr, json.RawMessage(`{"get_count":{}}`))
	require.NoError(t, err)

	var result map[string]string
	require.NoError(t, json.Unmarshal(resp, &result))
	require.Equal(t, "2", result["count"])
}

func TestVM_ExecuteUnauthorized(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	owner := validAddr()
	attacker := validAddr()

	codeID, err := k.Create(ctx, owner, counterContractCode(), nil)
	require.NoError(t, err)

	initMsg := json.RawMessage(fmt.Sprintf(`{"owner":"%s","initial_count":"0"}`, owner))
	contractAddr, err := k.Instantiate(ctx, codeID, owner, owner, "vm-counter", initMsg, nil)
	require.NoError(t, err)

	// Attacker tries to transfer ownership
	transferMsg := json.RawMessage(fmt.Sprintf(`{"transfer_ownership":{"new_owner":"%s"}}`, attacker))
	_, err = k.Execute(ctx, contractAddr, attacker, transferMsg, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unauthorized")

	// Owner should be unchanged
	resp, err := k.QuerySmart(ctx, contractAddr, json.RawMessage(`{"get_owner":{}}`))
	require.NoError(t, err)

	var result map[string]string
	require.NoError(t, json.Unmarshal(resp, &result))
	require.Equal(t, owner, result["owner"])
}

func TestVM_ExecuteTransferOwnership(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	owner := validAddr()
	newOwner := validAddr()

	codeID, err := k.Create(ctx, owner, counterContractCode(), nil)
	require.NoError(t, err)

	initMsg := json.RawMessage(fmt.Sprintf(`{"owner":"%s","initial_count":"0"}`, owner))
	contractAddr, err := k.Instantiate(ctx, codeID, owner, owner, "vm-counter", initMsg, nil)
	require.NoError(t, err)

	// Owner transfers ownership
	transferMsg := json.RawMessage(fmt.Sprintf(`{"transfer_ownership":{"new_owner":"%s"}}`, newOwner))
	resp, err := k.Execute(ctx, contractAddr, owner, transferMsg, nil)
	require.NoError(t, err)

	var result map[string]string
	require.NoError(t, json.Unmarshal(resp, &result))
	require.Equal(t, newOwner, result["new_owner"])
}

func TestVM_BankSend(t *testing.T) {
	k, ctx, _, bk := setupKeeper(t)

	owner := validAddr()
	recipient := validAddr()

	codeID, err := k.Create(ctx, owner, counterContractCode(), nil)
	require.NoError(t, err)

	initMsg := json.RawMessage(fmt.Sprintf(`{"owner":"%s","initial_count":"0"}`, owner))
	contractAddr, err := k.Instantiate(ctx, codeID, owner, owner, "vm-counter", initMsg, nil)
	require.NoError(t, err)

	// Give the contract some funds
	contractAccAddr, _ := sdk.AccAddressFromBech32(contractAddr)
	bk.balances[contractAddr] = sdk.NewCoins(sdk.NewInt64Coin("usyreen", 5000))
	_ = contractAccAddr

	// Send tokens from the contract
	sendMsg := json.RawMessage(fmt.Sprintf(`{"send_tokens":{"recipient":"%s","amount":"1000","denom":"usyreen"}}`, recipient))
	resp, err := k.Execute(ctx, contractAddr, owner, sendMsg, nil)
	require.NoError(t, err)

	var result map[string]string
	require.NoError(t, json.Unmarshal(resp, &result))
	require.Equal(t, "sent", result["status"])

	// Verify funds were transferred
	require.Equal(t, sdk.NewInt64Coin("usyreen", 4000).Amount, bk.balances[contractAddr].AmountOf("usyreen"))
	require.Equal(t, sdk.NewInt64Coin("usyreen", 1000).Amount, bk.balances[recipient].AmountOf("usyreen"))
}
