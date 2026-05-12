package keeper_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/require"

	"syreen/x/abstractaccount/keeper"
	"syreen/x/abstractaccount/types"
)

func init() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("syreen", "syreenpub")
	config.SetBech32PrefixForValidator("syreenvaloper", "syreenvaloperpub")
	config.SetBech32PrefixForConsensusNode("syreenvalcons", "syreenvalconspub")
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

func (m *mockAccountKeeper) GetAccount(_ context.Context, addr sdk.AccAddress) sdk.AccountI {
	acc, ok := m.accounts[addr.String()]
	if !ok {
		return nil
	}
	return acc
}

func (m *mockAccountKeeper) NewAccountWithAddress(_ context.Context, addr sdk.AccAddress) sdk.AccountI {
	acc := authtypes.NewBaseAccountWithAddress(addr)
	return acc
}

func (m *mockAccountKeeper) SetAccount(_ context.Context, acc sdk.AccountI) {
	m.accounts[acc.GetAddress().String()] = acc
}

func (m *mockAccountKeeper) GetModuleAddress(moduleName string) sdk.AccAddress {
	return sdk.AccAddress([]byte("module_" + moduleName))
}

type mockBankKeeper struct {
	sendErr error
}

func (m *mockBankKeeper) SendCoins(_ context.Context, _, _ sdk.AccAddress, _ sdk.Coins) error {
	return m.sendErr
}

func (m *mockBankKeeper) SendCoinsFromAccountToModule(_ context.Context, _ sdk.AccAddress, _ string, _ sdk.Coins) error {
	return m.sendErr
}

func (m *mockBankKeeper) GetBalance(_ context.Context, _ sdk.AccAddress, _ string) sdk.Coin {
	return sdk.NewInt64Coin("usyreen", 100_000_000) // 100M usyreen for sponsorship tests
}

func (m *mockBankKeeper) GetAllBalances(_ context.Context, _ sdk.AccAddress) sdk.Coins {
	return sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100_000_000))
}

// mockMsgRouter implements types.MsgRouter for testing.
type mockMsgRouter struct {
	handler baseapp.MsgServiceHandler
}

func (m *mockMsgRouter) Handler(_ sdk.Msg) baseapp.MsgServiceHandler {
	return m.handler
}

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

func validAddr() string {
	return sdk.AccAddress([]byte("test_address_padded_")).String()
}

func validAddr2() string {
	return sdk.AccAddress([]byte("second_addr_padding_")).String()
}

func validAddr3() string {
	return sdk.AccAddress([]byte("third__addr_padding_")).String()
}

func setupKeeper(t *testing.T) (keeper.Keeper, sdk.Context) {
	return setupKeeperWithRouter(t, nil)
}

func setupKeeperWithRouter(t *testing.T, router *mockMsgRouter) (keeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey("abstractaccount")
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, cmtproto.Header{Time: time.Now()}, false, log.NewNopLogger())
	registry := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)
	storeService := runtime.NewKVStoreService(storeKey)

	if router == nil {
		router = &mockMsgRouter{}
	}

	k := keeper.NewKeeper(cdc, storeService, newMockAccountKeeper(), &mockBankKeeper{}, router, "authority")
	return k, ctx
}

func setupKeeperWithParams(t *testing.T) (keeper.Keeper, sdk.Context) {
	k, ctx := setupKeeper(t)
	require.NoError(t, k.SetParams(ctx, types.DefaultParams()))
	return k, ctx
}

// createTestSmartAccount is a helper that creates a smart account and returns its address.
func createTestSmartAccount(t *testing.T, k keeper.Keeper, ctx sdk.Context, sender string, owners []string, threshold uint32) string {
	msg := &types.MsgCreateSmartAccount{
		Sender:      sender,
		AccountType: types.AccountTypeMultiSig,
		Owners:      owners,
		Threshold:   threshold,
	}
	addr, err := k.CreateSmartAccount(ctx, msg)
	require.NoError(t, err)
	require.NotEmpty(t, addr)
	return addr
}

// ---------------------------------------------------------------------------
// CreateSmartAccount
// ---------------------------------------------------------------------------

func TestCreateSmartAccount(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()
	addr2 := validAddr2()

	msg := &types.MsgCreateSmartAccount{
		Sender:      addr1,
		AccountType: types.AccountTypeMultiSig,
		Owners:      []string{addr1, addr2},
		Threshold:   2,
	}

	smartAddr, err := k.CreateSmartAccount(ctx, msg)
	require.NoError(t, err)
	require.NotEmpty(t, smartAddr)

	// Verify the account is stored
	account, found := k.GetSmartAccount(ctx, smartAddr)
	require.True(t, found)
	require.Equal(t, smartAddr, account.Address)
	require.Equal(t, types.AccountTypeMultiSig, account.AccountType)
	require.Equal(t, []string{addr1, addr2}, account.Owners)
	require.Equal(t, uint32(2), account.Threshold)
}

func TestCreateSmartAccount_DifferentTypes(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()

	// Social type should also create a recovery config
	msg := &types.MsgCreateSmartAccount{
		Sender:      addr1,
		AccountType: types.AccountTypeSocial,
		Owners:      []string{addr1},
		Threshold:   1,
	}

	smartAddr, err := k.CreateSmartAccount(ctx, msg)
	require.NoError(t, err)
	require.NotEmpty(t, smartAddr)

	// Verify recovery config was auto-created for social accounts
	config, found := k.GetRecoveryConfig(ctx, smartAddr)
	require.True(t, found)
	require.Equal(t, smartAddr, config.Account)
	require.Equal(t, []string{addr1}, config.Guardians)
	require.Equal(t, uint32(1), config.Threshold)
}

func TestCreateSmartAccount_UniqueAddresses(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()

	msg := &types.MsgCreateSmartAccount{
		Sender:      addr1,
		AccountType: types.AccountTypeMultiSig,
		Owners:      []string{addr1},
		Threshold:   1,
	}

	addr1Result, err := k.CreateSmartAccount(ctx, msg)
	require.NoError(t, err)

	// Creating another account (nonce increments so address differs)
	addr2Result, err := k.CreateSmartAccount(ctx, msg)
	require.NoError(t, err)
	require.NotEqual(t, addr1Result, addr2Result)
}

func TestCreateSmartAccount_InvalidSender(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)

	msg := &types.MsgCreateSmartAccount{
		Sender:      "invalid",
		AccountType: types.AccountTypeMultiSig,
		Owners:      []string{validAddr()},
		Threshold:   1,
	}

	_, err := k.CreateSmartAccount(ctx, msg)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrInvalidAddress)
}

// ---------------------------------------------------------------------------
// GetSmartAccount / SetSmartAccount roundtrip
// ---------------------------------------------------------------------------

func TestSmartAccountRoundtrip(t *testing.T) {
	k, ctx := setupKeeper(t)

	account := types.SmartAccount{
		Address:     validAddr(),
		AccountType: types.AccountTypeMultiSig,
		Owners:      []string{validAddr(), validAddr2()},
		Threshold:   2,
	}

	k.SetSmartAccount(ctx, account)

	got, found := k.GetSmartAccount(ctx, account.Address)
	require.True(t, found)
	require.Equal(t, account.Address, got.Address)
	require.Equal(t, account.AccountType, got.AccountType)
	require.Equal(t, account.Owners, got.Owners)
	require.Equal(t, account.Threshold, got.Threshold)
}

func TestGetSmartAccount_NotFound(t *testing.T) {
	k, ctx := setupKeeper(t)

	_, found := k.GetSmartAccount(ctx, validAddr())
	require.False(t, found)
}

// ---------------------------------------------------------------------------
// CreateSessionKey
// ---------------------------------------------------------------------------

func TestCreateSessionKey_Valid(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()
	addr2 := validAddr2()

	// The granter must be an owner of a smart account whose address is the granter.
	// So we need a smart account whose address matches the granter.
	// Since CreateSmartAccount derives a new address, we manually set one.
	account := types.SmartAccount{
		Address:     addr1,
		AccountType: types.AccountTypeMultiSig,
		Owners:      []string{addr1},
		Threshold:   1,
	}
	k.SetSmartAccount(ctx, account)

	perms := []types.Permission{
		{MsgType: "/cosmos.bank.v1beta1.MsgSend"},
	}
	err := k.CreateSessionKey(ctx, addr1, addr2, perms, time.Hour)
	require.NoError(t, err)

	// Verify the session key is stored
	sk, found := k.GetSessionKey(ctx, addr1, addr2)
	require.True(t, found)
	require.Equal(t, addr2, sk.Key)
	require.Equal(t, addr1, sk.Granter)
	require.Len(t, sk.Permissions, 1)
}

func TestCreateSessionKey_GranterNotOwner(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()
	addr2 := validAddr2()
	addr3 := validAddr3()

	// Smart account exists but addr2 is not an owner
	account := types.SmartAccount{
		Address:     addr2,
		AccountType: types.AccountTypeMultiSig,
		Owners:      []string{addr1}, // Only addr1 is owner
		Threshold:   1,
	}
	k.SetSmartAccount(ctx, account)

	perms := []types.Permission{
		{MsgType: "/cosmos.bank.v1beta1.MsgSend"},
	}
	err := k.CreateSessionKey(ctx, addr2, addr3, perms, time.Hour)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not an owner")
}

func TestCreateSessionKey_AccountNotFound(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()
	addr2 := validAddr2()

	perms := []types.Permission{
		{MsgType: "/cosmos.bank.v1beta1.MsgSend"},
	}
	err := k.CreateSessionKey(ctx, addr1, addr2, perms, time.Hour)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrAccountNotFound)
}

func TestCreateSessionKey_DurationTooLong(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()
	addr2 := validAddr2()

	account := types.SmartAccount{
		Address:     addr1,
		AccountType: types.AccountTypeMultiSig,
		Owners:      []string{addr1},
		Threshold:   1,
	}
	k.SetSmartAccount(ctx, account)

	perms := []types.Permission{
		{MsgType: "/cosmos.bank.v1beta1.MsgSend"},
	}
	// Default max is 24h; exceed it
	err := k.CreateSessionKey(ctx, addr1, addr2, perms, 48*time.Hour)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrSessionKeyDuration)
}

// ---------------------------------------------------------------------------
// ValidateSessionKey
// ---------------------------------------------------------------------------

func TestValidateSessionKey_Valid(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()
	addr2 := validAddr2()

	account := types.SmartAccount{
		Address:     addr1,
		AccountType: types.AccountTypeMultiSig,
		Owners:      []string{addr1},
		Threshold:   1,
	}
	k.SetSmartAccount(ctx, account)

	perms := []types.Permission{
		{MsgType: "/cosmos.bank.v1beta1.MsgSend"},
	}
	require.NoError(t, k.CreateSessionKey(ctx, addr1, addr2, perms, time.Hour))

	err := k.ValidateSessionKey(ctx, addr2, "/cosmos.bank.v1beta1.MsgSend", nil)
	require.NoError(t, err)
}

func TestValidateSessionKey_Expired(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()
	addr2 := validAddr2()

	account := types.SmartAccount{
		Address:     addr1,
		AccountType: types.AccountTypeMultiSig,
		Owners:      []string{addr1},
		Threshold:   1,
	}
	k.SetSmartAccount(ctx, account)

	perms := []types.Permission{
		{MsgType: "/cosmos.bank.v1beta1.MsgSend"},
	}
	require.NoError(t, k.CreateSessionKey(ctx, addr1, addr2, perms, time.Hour))

	// Advance time past expiry
	futureCtx := ctx.WithBlockTime(ctx.BlockTime().Add(2 * time.Hour))
	err := k.ValidateSessionKey(futureCtx, addr2, "/cosmos.bank.v1beta1.MsgSend", nil)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrSessionExpired)
}

func TestValidateSessionKey_PermissionDenied(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()
	addr2 := validAddr2()

	account := types.SmartAccount{
		Address:     addr1,
		AccountType: types.AccountTypeMultiSig,
		Owners:      []string{addr1},
		Threshold:   1,
	}
	k.SetSmartAccount(ctx, account)

	perms := []types.Permission{
		{MsgType: "/cosmos.bank.v1beta1.MsgSend"},
	}
	require.NoError(t, k.CreateSessionKey(ctx, addr1, addr2, perms, time.Hour))

	// Try a different message type
	err := k.ValidateSessionKey(ctx, addr2, "/cosmos.staking.v1beta1.MsgDelegate", nil)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrPermissionDenied)
}

func TestValidateSessionKey_SpendLimitExceeded(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()
	addr2 := validAddr2()

	account := types.SmartAccount{
		Address:     addr1,
		AccountType: types.AccountTypeMultiSig,
		Owners:      []string{addr1},
		Threshold:   1,
	}
	k.SetSmartAccount(ctx, account)

	perms := []types.Permission{
		{
			MsgType:   "/cosmos.bank.v1beta1.MsgSend",
			MaxAmount: sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		},
	}
	require.NoError(t, k.CreateSessionKey(ctx, addr1, addr2, perms, time.Hour))

	// First spend within limit (read-only check)
	err := k.ValidateSessionKey(ctx, addr2, "/cosmos.bank.v1beta1.MsgSend",
		sdk.NewCoins(sdk.NewInt64Coin("usyreen", 50)))
	require.NoError(t, err)

	// Record the usage (simulates PostHandler after successful tx)
	err = k.RecordSessionKeyUsage(ctx, addr2, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 50)))
	require.NoError(t, err)

	// Second spend exceeds limit (50 used + 60 = 110 > 100)
	err = k.ValidateSessionKey(ctx, addr2, "/cosmos.bank.v1beta1.MsgSend",
		sdk.NewCoins(sdk.NewInt64Coin("usyreen", 60)))
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrSessionLimitExceeded)
}

func TestValidateSessionKey_NotFound(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)

	err := k.ValidateSessionKey(ctx, validAddr(), "/cosmos.bank.v1beta1.MsgSend", nil)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrSessionKeyNotFound)
}

// ---------------------------------------------------------------------------
// RevokeSessionKey
// ---------------------------------------------------------------------------

func TestRevokeSessionKey_Valid(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()
	addr2 := validAddr2()

	account := types.SmartAccount{
		Address:     addr1,
		AccountType: types.AccountTypeMultiSig,
		Owners:      []string{addr1},
		Threshold:   1,
	}
	k.SetSmartAccount(ctx, account)

	perms := []types.Permission{
		{MsgType: "/cosmos.bank.v1beta1.MsgSend"},
	}
	require.NoError(t, k.CreateSessionKey(ctx, addr1, addr2, perms, time.Hour))

	// Revoke
	err := k.RevokeSessionKey(ctx, addr1, addr2)
	require.NoError(t, err)

	// Verify it's gone
	_, found := k.GetSessionKey(ctx, addr1, addr2)
	require.False(t, found)
}

func TestRevokeSessionKey_NotFound(t *testing.T) {
	k, ctx := setupKeeper(t)

	err := k.RevokeSessionKey(ctx, validAddr(), validAddr2())
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrSessionKeyNotFound)
}

// ---------------------------------------------------------------------------
// InitiateRecovery
// ---------------------------------------------------------------------------

func TestInitiateRecovery_Valid(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()
	addr2 := validAddr2()
	addr3 := validAddr3()

	// Set up a smart account and recovery config
	account := types.SmartAccount{
		Address:     addr1,
		AccountType: types.AccountTypeSocial,
		Owners:      []string{addr1},
		Threshold:   1,
	}
	k.SetSmartAccount(ctx, account)

	config := types.RecoveryConfig{
		Account:     addr1,
		Guardians:   []string{addr2},
		Threshold:   1,
		DelayPeriod: 48 * time.Hour,
	}
	k.SetRecoveryConfig(ctx, config)

	err := k.InitiateRecovery(ctx, addr2, addr1, []string{addr3})
	require.NoError(t, err)

	// Verify request exists
	request, found := k.GetRecoveryRequest(ctx, addr1)
	require.True(t, found)
	require.Equal(t, []string{addr3}, request.NewOwners)
	require.Equal(t, []string{addr2}, request.Approvals)
}

func TestInitiateRecovery_NotGuardian(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()
	addr2 := validAddr2()
	addr3 := validAddr3()

	config := types.RecoveryConfig{
		Account:     addr1,
		Guardians:   []string{addr2},
		Threshold:   1,
		DelayPeriod: 48 * time.Hour,
	}
	k.SetRecoveryConfig(ctx, config)

	// addr3 is not a guardian
	err := k.InitiateRecovery(ctx, addr3, addr1, []string{addr2})
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrNotGuardian)
}

func TestInitiateRecovery_AlreadyInProgress(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()
	addr2 := validAddr2()
	addr3 := validAddr3()

	config := types.RecoveryConfig{
		Account:     addr1,
		Guardians:   []string{addr2},
		Threshold:   1,
		DelayPeriod: 48 * time.Hour,
	}
	k.SetRecoveryConfig(ctx, config)

	// First initiation succeeds
	require.NoError(t, k.InitiateRecovery(ctx, addr2, addr1, []string{addr3}))

	// Second initiation fails
	err := k.InitiateRecovery(ctx, addr2, addr1, []string{addr3})
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrRecoveryInProgress)
}

func TestInitiateRecovery_NoConfig(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)

	err := k.InitiateRecovery(ctx, validAddr(), validAddr2(), []string{validAddr3()})
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrRecoveryConfigNotFound)
}

// ---------------------------------------------------------------------------
// ApproveRecovery
// ---------------------------------------------------------------------------

func TestApproveRecovery_Valid(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()
	addr2 := validAddr2()
	addr3 := validAddr3()

	config := types.RecoveryConfig{
		Account:     addr1,
		Guardians:   []string{addr2, addr3},
		Threshold:   2,
		DelayPeriod: 48 * time.Hour,
	}
	k.SetRecoveryConfig(ctx, config)

	// Initiate with guardian addr2
	require.NoError(t, k.InitiateRecovery(ctx, addr2, addr1, []string{addr3}))

	// Approve with guardian addr3
	err := k.ApproveRecovery(ctx, addr3, addr1)
	require.NoError(t, err)

	// Verify approval was added
	request, found := k.GetRecoveryRequest(ctx, addr1)
	require.True(t, found)
	require.Len(t, request.Approvals, 2)
	require.Contains(t, request.Approvals, addr2)
	require.Contains(t, request.Approvals, addr3)
}

func TestApproveRecovery_DuplicateApproval(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()
	addr2 := validAddr2()
	addr3 := validAddr3()

	config := types.RecoveryConfig{
		Account:     addr1,
		Guardians:   []string{addr2, addr3},
		Threshold:   2,
		DelayPeriod: 48 * time.Hour,
	}
	k.SetRecoveryConfig(ctx, config)

	require.NoError(t, k.InitiateRecovery(ctx, addr2, addr1, []string{addr3}))

	// addr2 already approved during initiation; trying again should fail
	err := k.ApproveRecovery(ctx, addr2, addr1)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrDuplicateApproval)
}

func TestApproveRecovery_NotGuardian(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()
	addr2 := validAddr2()
	addr3 := validAddr3()

	config := types.RecoveryConfig{
		Account:     addr1,
		Guardians:   []string{addr2},
		Threshold:   1,
		DelayPeriod: 48 * time.Hour,
	}
	k.SetRecoveryConfig(ctx, config)

	require.NoError(t, k.InitiateRecovery(ctx, addr2, addr1, []string{addr3}))

	// addr3 is not a guardian
	err := k.ApproveRecovery(ctx, addr3, addr1)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrNotGuardian)
}

// ---------------------------------------------------------------------------
// ExecuteRecovery
// ---------------------------------------------------------------------------

func TestExecuteRecovery_Valid(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()
	addr2 := validAddr2()
	addr3 := validAddr3()

	// Set up smart account
	account := types.SmartAccount{
		Address:     addr1,
		AccountType: types.AccountTypeSocial,
		Owners:      []string{addr1},
		Threshold:   1,
	}
	k.SetSmartAccount(ctx, account)

	config := types.RecoveryConfig{
		Account:     addr1,
		Guardians:   []string{addr2},
		Threshold:   1,
		DelayPeriod: time.Hour, // 1 hour delay
	}
	k.SetRecoveryConfig(ctx, config)

	require.NoError(t, k.InitiateRecovery(ctx, addr2, addr1, []string{addr3}))

	// Advance time past the delay period
	futureCtx := ctx.WithBlockTime(ctx.BlockTime().Add(2 * time.Hour))

	err := k.ExecuteRecovery(futureCtx, addr2, addr1)
	require.NoError(t, err)

	// Verify owners were swapped
	updatedAccount, found := k.GetSmartAccount(futureCtx, addr1)
	require.True(t, found)
	require.Equal(t, []string{addr3}, updatedAccount.Owners)

	// Verify recovery request was cleaned up
	_, found = k.GetRecoveryRequest(futureCtx, addr1)
	require.False(t, found)
}

func TestExecuteRecovery_InsufficientGuardians(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()
	addr2 := validAddr2()
	addr3 := validAddr3()

	account := types.SmartAccount{
		Address:     addr1,
		AccountType: types.AccountTypeSocial,
		Owners:      []string{addr1},
		Threshold:   1,
	}
	k.SetSmartAccount(ctx, account)

	config := types.RecoveryConfig{
		Account:     addr1,
		Guardians:   []string{addr2, addr3},
		Threshold:   2, // Need 2 approvals
		DelayPeriod: time.Hour,
	}
	k.SetRecoveryConfig(ctx, config)

	// Only one guardian initiates (which counts as 1 approval)
	require.NoError(t, k.InitiateRecovery(ctx, addr2, addr1, []string{addr3}))

	futureCtx := ctx.WithBlockTime(ctx.BlockTime().Add(2 * time.Hour))
	err := k.ExecuteRecovery(futureCtx, addr2, addr1)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrInsufficientGuardians)
}

func TestExecuteRecovery_DelayNotMet(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()
	addr2 := validAddr2()
	addr3 := validAddr3()

	account := types.SmartAccount{
		Address:     addr1,
		AccountType: types.AccountTypeSocial,
		Owners:      []string{addr1},
		Threshold:   1,
	}
	k.SetSmartAccount(ctx, account)

	config := types.RecoveryConfig{
		Account:     addr1,
		Guardians:   []string{addr2},
		Threshold:   1,
		DelayPeriod: 48 * time.Hour,
	}
	k.SetRecoveryConfig(ctx, config)

	require.NoError(t, k.InitiateRecovery(ctx, addr2, addr1, []string{addr3}))

	// Try to execute immediately (delay not met)
	err := k.ExecuteRecovery(ctx, addr2, addr1)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrRecoveryDelayNotMet)
}

func TestExecuteRecovery_NoRequest(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()
	addr2 := validAddr2()

	config := types.RecoveryConfig{
		Account:     addr1,
		Guardians:   []string{addr2},
		Threshold:   1,
		DelayPeriod: time.Hour,
	}
	k.SetRecoveryConfig(ctx, config)

	err := k.ExecuteRecovery(ctx, addr2, addr1)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrRecoveryNotFound)
}

// ---------------------------------------------------------------------------
// SponsorGas
// ---------------------------------------------------------------------------

func TestSponsorGas_Valid(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()
	addr2 := validAddr2()

	err := k.SponsorGas(ctx, addr1, addr2, 1000000, time.Hour)
	require.NoError(t, err)

	// Verify the sponsorship record
	sponsor, found := k.GetGasSponsor(ctx, addr1, addr2)
	require.True(t, found)
	require.Equal(t, addr1, sponsor.Sponsor)
	require.Equal(t, addr2, sponsor.Sponsored)
	require.Equal(t, uint64(1000000), sponsor.GasLimit)
}

func TestSponsorGas_Disabled(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)

	// Disable gas sponsorship
	params := types.DefaultParams()
	params.EnableGasSponsorship = false
	require.NoError(t, k.SetParams(ctx, params))

	err := k.SponsorGas(ctx, validAddr(), validAddr2(), 1000000, time.Hour)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrGasSponsorshipDisabled)
}

// ---------------------------------------------------------------------------
// CheckGasSponsor
// ---------------------------------------------------------------------------

func TestCheckGasSponsor_Valid(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()
	addr2 := validAddr2()

	require.NoError(t, k.SponsorGas(ctx, addr1, addr2, 1000000, time.Hour))

	sponsor, found := k.CheckGasSponsor(ctx, addr2, 500)
	require.True(t, found)
	require.NotNil(t, sponsor)
	require.Equal(t, addr1, sponsor.Sponsor)
}

func TestCheckGasSponsor_Expired(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()
	addr2 := validAddr2()

	require.NoError(t, k.SponsorGas(ctx, addr1, addr2, 1000000, time.Hour))

	futureCtx := ctx.WithBlockTime(ctx.BlockTime().Add(2 * time.Hour))
	_, found := k.CheckGasSponsor(futureCtx, addr2, 500)
	require.False(t, found)
}

func TestCheckGasSponsor_LimitExceeded(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()
	addr2 := validAddr2()

	require.NoError(t, k.SponsorGas(ctx, addr1, addr2, 100, time.Hour))

	// Request more gas than the limit
	_, found := k.CheckGasSponsor(ctx, addr2, 200)
	require.False(t, found)
}

// ---------------------------------------------------------------------------
// ExecuteBatch
// ---------------------------------------------------------------------------

func TestExecuteBatch_TooLarge(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)
	addr1 := validAddr()

	account := types.SmartAccount{
		Address:     addr1,
		AccountType: types.AccountTypeMultiSig,
		Owners:      []string{addr1},
		Threshold:   1,
	}
	k.SetSmartAccount(ctx, account)

	// Default max batch size is 10; create 11 messages
	msgs := make([]json.RawMessage, 11)
	for i := range msgs {
		msgs[i] = json.RawMessage(`{"type":"send"}`)
	}

	_, err := k.ExecuteBatch(ctx, addr1, msgs)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrBatchTooLarge)
}

func TestExecuteBatch_AccountNotFound(t *testing.T) {
	k, ctx := setupKeeperWithParams(t)

	msgs := []json.RawMessage{json.RawMessage(`{"type":"send"}`)}
	_, err := k.ExecuteBatch(ctx, validAddr(), msgs)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrAccountNotFound)
}

// TestExecuteBatch_TwoSuccessful tests that a batch of 2 valid messages both execute.
func TestExecuteBatch_TwoSuccessful(t *testing.T) {
	callCount := 0
	router := &mockMsgRouter{
		handler: func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
			callCount++
			return &sdk.Result{}, nil
		},
	}
	k, ctx := setupKeeperWithRouter(t, router)
	require.NoError(t, k.SetParams(ctx, types.DefaultParams()))

	addr1 := validAddr()
	account := types.SmartAccount{
		Address:     addr1,
		AccountType: types.AccountTypeMultiSig,
		Owners:      []string{addr1},
		Threshold:   1,
	}
	k.SetSmartAccount(ctx, account)

	// Build two valid batch messages. The @type must resolve via the codec registry.
	msg1 := json.RawMessage(fmt.Sprintf(`{"@type":"/syreen.abstractaccount.MsgSponsorGas","sponsor":"%s","sponsored":"%s","gas_limit":1000,"duration":3600000000000}`, addr1, validAddr2()))
	msg2 := json.RawMessage(fmt.Sprintf(`{"@type":"/syreen.abstractaccount.MsgSponsorGas","sponsor":"%s","sponsored":"%s","gas_limit":2000,"duration":3600000000000}`, addr1, validAddr3()))

	results, err := k.ExecuteBatch(ctx, addr1, []json.RawMessage{msg1, msg2})
	require.NoError(t, err)
	require.Len(t, results, 2)
	require.Equal(t, 2, callCount)
}

// TestExecuteBatch_SecondFails_FirstRolledBack tests atomicity:
// if the second message fails, the first should be rolled back.
func TestExecuteBatch_SecondFails_FirstRolledBack(t *testing.T) {
	callCount := 0
	router := &mockMsgRouter{
		handler: func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
			callCount++
			if callCount == 2 {
				return nil, fmt.Errorf("simulated failure on second message")
			}
			return &sdk.Result{}, nil
		},
	}
	k, ctx := setupKeeperWithRouter(t, router)
	require.NoError(t, k.SetParams(ctx, types.DefaultParams()))

	addr1 := validAddr()
	account := types.SmartAccount{
		Address:     addr1,
		AccountType: types.AccountTypeMultiSig,
		Owners:      []string{addr1},
		Threshold:   1,
	}
	k.SetSmartAccount(ctx, account)

	msg1 := json.RawMessage(fmt.Sprintf(`{"@type":"/syreen.abstractaccount.MsgSponsorGas","sponsor":"%s","sponsored":"%s","gas_limit":1000,"duration":3600000000000}`, addr1, validAddr2()))
	msg2 := json.RawMessage(fmt.Sprintf(`{"@type":"/syreen.abstractaccount.MsgSponsorGas","sponsor":"%s","sponsored":"%s","gas_limit":2000,"duration":3600000000000}`, addr1, validAddr3()))

	_, err := k.ExecuteBatch(ctx, addr1, []json.RawMessage{msg1, msg2})
	require.Error(t, err)
	require.Contains(t, err.Error(), "simulated failure on second message")

	// Verify the first message's side effects were rolled back:
	// callCount is 2 (both were attempted), but since we used CacheContext
	// and didn't commit, state changes from the first message are not persisted.
	require.Equal(t, 2, callCount)
}

// ---------------------------------------------------------------------------
// SetParams / GetParams roundtrip
// ---------------------------------------------------------------------------

func TestParamsRoundtrip(t *testing.T) {
	k, ctx := setupKeeper(t)

	params := types.Params{
		MaxSessionKeyDuration: 12 * time.Hour,
		MaxGuardians:          5,
		RecoveryDelayPeriod:   72 * time.Hour,
		MaxBatchSize:          20,
		EnableGasSponsorship:  false,
	}

	require.NoError(t, k.SetParams(ctx, params))

	got := k.GetParams(ctx)
	require.Equal(t, params.MaxSessionKeyDuration, got.MaxSessionKeyDuration)
	require.Equal(t, params.MaxGuardians, got.MaxGuardians)
	require.Equal(t, params.RecoveryDelayPeriod, got.RecoveryDelayPeriod)
	require.Equal(t, params.MaxBatchSize, got.MaxBatchSize)
	require.Equal(t, params.EnableGasSponsorship, got.EnableGasSponsorship)
}

func TestGetParams_DefaultsWhenUnset(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Without setting params, should return defaults
	params := k.GetParams(ctx)
	require.Equal(t, types.DefaultParams(), params)
}

// ---------------------------------------------------------------------------
// Genesis InitGenesis / ExportGenesis
// ---------------------------------------------------------------------------

func TestGenesis_InitAndExport(t *testing.T) {
	k, ctx := setupKeeper(t)

	genesis := types.GenesisState{
		Params: types.Params{
			MaxSessionKeyDuration: 6 * time.Hour,
			MaxGuardians:          3,
			RecoveryDelayPeriod:   24 * time.Hour,
			MaxBatchSize:          5,
			EnableGasSponsorship:  false,
		},
	}

	k.InitGenesis(ctx, genesis)

	exported := k.ExportGenesis(ctx)
	require.NotNil(t, exported)
	require.Equal(t, genesis.Params.MaxSessionKeyDuration, exported.Params.MaxSessionKeyDuration)
	require.Equal(t, genesis.Params.MaxGuardians, exported.Params.MaxGuardians)
	require.Equal(t, genesis.Params.RecoveryDelayPeriod, exported.Params.RecoveryDelayPeriod)
	require.Equal(t, genesis.Params.MaxBatchSize, exported.Params.MaxBatchSize)
	require.Equal(t, genesis.Params.EnableGasSponsorship, exported.Params.EnableGasSponsorship)
}

func TestGenesis_DefaultRoundtrip(t *testing.T) {
	k, ctx := setupKeeper(t)

	defaultGenesis := *types.DefaultGenesis()
	k.InitGenesis(ctx, defaultGenesis)

	exported := k.ExportGenesis(ctx)
	require.Equal(t, defaultGenesis.Params, exported.Params)
}
