package keeper_test

import (
	"context"
	"encoding/json"
	"fmt"
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
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank"
	staking "github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/stretchr/testify/require"

	"syreen/x/intent/keeper"
	"syreen/x/intent/types"
)

// newTestCodec builds a codec whose interface registry has the std, bank, and
// staking interfaces/messages registered, and an address codec configured, so
// that real sdk.Msgs (e.g. bank MsgSend) can be decoded from protojson via
// UnmarshalInterfaceJSON and have their signers resolved via GetMsgV1Signers.
// A bare codec.NewProtoCodec(codectypes.NewInterfaceRegistry()) has nothing
// registered and a failing address codec, so executeSolutionMsgs can neither
// decode nor authorize any real execution message with it.
func newTestCodec() codec.Codec {
	return moduletestutil.MakeTestEncodingConfig(bank.AppModuleBasic{}, staking.AppModuleBasic{}).Codec
}

// bankSendExecMsg returns a real, registered, signer-resolvable execution
// message: a bank MsgSend whose signer (from_address) is the intent module
// account, which executeSolutionMsgs authorizes. The mock msg router stubs the
// handler, so the send itself is a no-op (no funding required); the message only
// needs to decode and resolve an authorized signer.
func bankSendExecMsg() json.RawMessage {
	moduleAddr := authtypes.NewModuleAddress(types.ModuleName).String()
	return json.RawMessage(fmt.Sprintf(
		`{"@type":"/cosmos.bank.v1beta1.MsgSend","from_address":"%s","to_address":"%s","amount":[{"denom":"usyreen","amount":"1"}]}`,
		moduleAddr, creatorAddr))
}

// swapIntentBody returns a well-formed SwapIntent body. Swap intents run through
// verifySwapOutcome on fulfillment, which requires min_output_amount to be a real
// (non-nil) value; the old placeholder body `{"test":1}` left it nil and caused a
// nil-pointer panic in the outcome check. MinOutputAmount is zero here so the
// mock bank's flat balance satisfies the "delivered >= required" check.
func swapIntentBody() json.RawMessage {
	return json.RawMessage(`{"input_denom":"usyreen","input_amount":"1000","output_denom":"uusdc","min_output_amount":"0","max_slippage":"0"}`)
}

// stakingDelegateExecMsg returns a real, registered staking MsgDelegate. It is
// used to exercise the whitelist rejection path: the message decodes fine (so we
// are past the decode step) but is NOT in the restricted whitelist, so it must be
// rejected with ErrDisallowedMsgType.
func stakingDelegateExecMsg() json.RawMessage {
	return json.RawMessage(fmt.Sprintf(
		`{"@type":"/cosmos.staking.v1beta1.MsgDelegate","delegator_address":"%s","validator_address":"%s","amount":{"denom":"usyreen","amount":"1"}}`,
		creatorAddr, solverAddr))
}

// ---------- Valid test addresses ----------

var (
	creatorAddr = sdk.AccAddress([]byte("creator_____________")).String()
	solverAddr  = sdk.AccAddress([]byte("solver______________")).String()
	solverAddr2 = sdk.AccAddress([]byte("solver2_____________")).String()
)

// ---------- Mock AccountKeeper ----------

type mockAccountKeeper struct{}

func (m *mockAccountKeeper) GetAccount(_ context.Context, _ sdk.AccAddress) sdk.AccountI {
	return nil
}

func (m *mockAccountKeeper) GetModuleAddress(moduleName string) sdk.AccAddress {
	return sdk.AccAddress([]byte(moduleName))
}

func (m *mockAccountKeeper) GetModuleAccount(_ context.Context, moduleName string) sdk.ModuleAccountI {
	return authtypes.NewEmptyModuleAccount(moduleName)
}

// ---------- Mock BankKeeper ----------

type mockBankKeeper struct {
	// Track balances to allow failure simulation
	sendErr error
}

func (m *mockBankKeeper) SendCoins(_ context.Context, _, _ sdk.AccAddress, _ sdk.Coins) error {
	return m.sendErr
}

func (m *mockBankKeeper) SendCoinsFromAccountToModule(_ context.Context, _ sdk.AccAddress, _ string, _ sdk.Coins) error {
	return m.sendErr
}

func (m *mockBankKeeper) SendCoinsFromModuleToAccount(_ context.Context, _ string, _ sdk.AccAddress, _ sdk.Coins) error {
	return m.sendErr
}

func (m *mockBankKeeper) SendCoinsFromModuleToModule(_ context.Context, _, _ string, _ sdk.Coins) error {
	return m.sendErr
}

func (m *mockBankKeeper) GetBalance(_ context.Context, _ sdk.AccAddress, _ string) sdk.Coin {
	return sdk.NewInt64Coin("usyreen", 1000000)
}
func (m *mockBankKeeper) GetAllBalances(_ context.Context, _ sdk.AccAddress) sdk.Coins {
	return sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1000000))
}

// ---------- Mock MsgRouter ----------

type mockMsgRouter struct {
	handler baseapp.MsgServiceHandler
}

func (m *mockMsgRouter) Handler(_ sdk.Msg) baseapp.MsgServiceHandler {
	return m.handler
}

// ---------- Tracking Bank Keeper ----------

type moduleToAccountCall struct {
	module    string
	recipient sdk.AccAddress
	amount    sdk.Coins
}

type trackingBankKeeper struct {
	mockBankKeeper
	moduleToAccountCalls []moduleToAccountCall
}

func (m *trackingBankKeeper) SendCoinsFromModuleToAccount(_ context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	m.moduleToAccountCalls = append(m.moduleToAccountCalls, moduleToAccountCall{
		module:    senderModule,
		recipient: recipientAddr,
		amount:    amt,
	})
	return m.sendErr
}

// ---------- Setup ----------

func setupKeeper(t *testing.T) (*keeper.Keeper, sdk.Context) {
	return setupKeeperWithBankErr(t, nil)
}

func setupKeeperWithBankErr(t *testing.T, bankErr error) (*keeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey("intent")
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, cmtproto.Header{Height: 1}, false, log.NewNopLogger())
	cdc := newTestCodec()
	storeService := runtime.NewKVStoreService(storeKey)

	// Mock msg router that succeeds by default
	router := &mockMsgRouter{
		handler: func(_ sdk.Context, _ sdk.Msg) (*sdk.Result, error) {
			return &sdk.Result{}, nil
		},
	}

	k := keeper.NewKeeper(cdc, storeService, &mockAccountKeeper{}, &mockBankKeeper{sendErr: bankErr}, router, "authority")

	// Set default params so intents are enabled, with a test type in the whitelist
	params := types.DefaultParams()
	params.AllowedMsgTypes = append(params.AllowedMsgTypes, "/")
	require.NoError(t, k.SetParams(ctx, params))

	return k, ctx
}

func setupKeeperWithTrackingBank(t *testing.T) (*keeper.Keeper, sdk.Context, *trackingBankKeeper) {
	storeKey := storetypes.NewKVStoreKey("intent")
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, cmtproto.Header{Height: 1}, false, log.NewNopLogger())
	cdc := newTestCodec()
	storeService := runtime.NewKVStoreService(storeKey)

	router := &mockMsgRouter{
		handler: func(_ sdk.Context, _ sdk.Msg) (*sdk.Result, error) {
			return nil, fmt.Errorf("execution failed")
		},
	}

	bank := &trackingBankKeeper{}

	k := keeper.NewKeeper(cdc, storeService, &mockAccountKeeper{}, bank, router, "authority")
	params := types.DefaultParams()
	params.AllowedMsgTypes = append(params.AllowedMsgTypes, "/")
	require.NoError(t, k.SetParams(ctx, params))

	return k, ctx, bank
}

// ===================== SetParams / GetParams =====================

func TestSetGetParams(t *testing.T) {
	k, ctx := setupKeeper(t)

	p := types.DefaultParams()
	p.SolvingWindow = 42
	require.NoError(t, k.SetParams(ctx, p))

	got := k.GetParams(ctx)
	require.Equal(t, uint64(42), got.SolvingWindow)
}

func TestGetParams_DefaultWhenUnset(t *testing.T) {
	// Create a keeper without setting params
	storeKey := storetypes.NewKVStoreKey("intent")
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, cmtproto.Header{Height: 1}, false, log.NewNopLogger())
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	storeService := runtime.NewKVStoreService(storeKey)
	router := &mockMsgRouter{}
	k := keeper.NewKeeper(cdc, storeService, &mockAccountKeeper{}, &mockBankKeeper{}, router, "authority")

	got := k.GetParams(ctx)
	require.Equal(t, types.DefaultParams().SolvingWindow, got.SolvingWindow)
}

// ===================== SetIntent / GetIntent =====================

func TestSetGetIntent(t *testing.T) {
	k, ctx := setupKeeper(t)

	intent := types.Intent{
		ID:         "42",
		Creator:    creatorAddr,
		IntentType: types.IntentTypeSwap,
		Body:       json.RawMessage(`{"test":true}`),
		MaxFee:     sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:        sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10)),
		Expiry:     50,
		Status:     types.StatusPending,
		CreatedAt:  1,
	}

	k.SetIntent(ctx, intent)
	got, found := k.GetIntent(ctx, "42")
	require.True(t, found)
	require.Equal(t, "42", got.ID)
	require.Equal(t, creatorAddr, got.Creator)
	require.Equal(t, types.StatusPending, got.Status)
}

func TestGetIntent_NotFound(t *testing.T) {
	k, ctx := setupKeeper(t)
	_, found := k.GetIntent(ctx, "nonexistent")
	require.False(t, found)
}

// ===================== SetSolver / GetSolver =====================

func TestSetGetSolver(t *testing.T) {
	k, ctx := setupKeeper(t)

	solver := types.Solver{
		Address:         solverAddr,
		Moniker:         "test-solver",
		StakedAmount:    sdk.NewInt64Coin("usyreen", 5000),
		ReputationScore: 100,
		Active:          true,
		JoinedAt:        1,
	}

	k.SetSolver(ctx, solver)
	got, found := k.GetSolver(ctx, solverAddr)
	require.True(t, found)
	require.Equal(t, "test-solver", got.Moniker)
	require.Equal(t, uint64(100), got.ReputationScore)
	require.True(t, got.Active)
}

func TestGetSolver_NotFound(t *testing.T) {
	k, ctx := setupKeeper(t)
	_, found := k.GetSolver(ctx, "nonexistent")
	require.False(t, found)
}

// ===================== SubmitIntent =====================

func TestSubmitIntent_Success(t *testing.T) {
	k, ctx := setupKeeper(t)

	msg := &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{"swap":"data"}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:          sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10)),
		ExpiryBlocks: 50,
	}

	id, err := k.SubmitIntent(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, "1", id)

	intent, found := k.GetIntent(ctx, id)
	require.True(t, found)
	require.Equal(t, creatorAddr, intent.Creator)
	require.Equal(t, types.StatusPending, intent.Status)
	require.Equal(t, int64(51), intent.Expiry) // block 1 + 50
}

func TestSubmitIntent_SequentialIDs(t *testing.T) {
	k, ctx := setupKeeper(t)

	msg := &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeTransfer,
		Body:         json.RawMessage(`{}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1)),
		Tip:          sdk.Coins{},
		ExpiryBlocks: 10,
	}

	id1, err := k.SubmitIntent(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, "1", id1)

	id2, err := k.SubmitIntent(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, "2", id2)

	id3, err := k.SubmitIntent(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, "3", id3)
}

func TestSubmitIntent_CapsExpiryToParam(t *testing.T) {
	k, ctx := setupKeeper(t)

	msg := &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1)),
		Tip:          sdk.Coins{},
		ExpiryBlocks: 99999, // much larger than default
	}

	id, err := k.SubmitIntent(ctx, msg)
	require.NoError(t, err)

	intent, found := k.GetIntent(ctx, id)
	require.True(t, found)
	// Should be capped to DefaultIntentExpiryBlocks (20000) + block height (1)
	require.Equal(t, int64(int64(types.DefaultIntentExpiryBlocks)+1), intent.Expiry)
}

func TestSubmitIntent_DisabledIntents(t *testing.T) {
	k, ctx := setupKeeper(t)

	p := types.DefaultParams()
	p.EnableIntents = false
	require.NoError(t, k.SetParams(ctx, p))

	msg := &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1)),
		Tip:          sdk.Coins{},
		ExpiryBlocks: 10,
	}

	_, err := k.SubmitIntent(ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "intents are currently disabled")
}

func TestSubmitIntent_LockFundsFailure(t *testing.T) {
	k, ctx := setupKeeperWithBankErr(t, fmt.Errorf("insufficient funds"))

	msg := &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:          sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10)),
		ExpiryBlocks: 10,
	}

	_, err := k.SubmitIntent(ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to lock fees")
}

// ===================== RegisterSolver =====================

func TestRegisterSolver_Success(t *testing.T) {
	k, ctx := setupKeeper(t)

	msg := &types.MsgRegisterSolver{
		Address:     solverAddr,
		Moniker:     "solver-one",
		StakeAmount: sdk.NewCoin("usyreen", math.NewInt(2000000000)),
	}

	err := k.RegisterSolver(ctx, msg)
	require.NoError(t, err)

	solver, found := k.GetSolver(ctx, solverAddr)
	require.True(t, found)
	require.Equal(t, uint64(100), solver.ReputationScore)
	require.True(t, solver.Active)
	require.Equal(t, "solver-one", solver.Moniker)
	require.Equal(t, int64(1), solver.JoinedAt)
}

func TestRegisterSolver_Duplicate(t *testing.T) {
	k, ctx := setupKeeper(t)

	msg := &types.MsgRegisterSolver{
		Address:     solverAddr,
		Moniker:     "solver-one",
		StakeAmount: sdk.NewCoin("usyreen", math.NewInt(2000000000)),
	}

	require.NoError(t, k.RegisterSolver(ctx, msg))
	err := k.RegisterSolver(ctx, msg)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrSolverAlreadyRegistered)
}

func TestRegisterSolver_InsufficientStake(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Default MinSolverStake = 1000000000usyreen (1000 SYR)
	msg := &types.MsgRegisterSolver{
		Address:     solverAddr,
		Moniker:     "solver",
		StakeAmount: sdk.NewInt64Coin("usyreen", 500), // below min
	}

	err := k.RegisterSolver(ctx, msg)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrInsufficientStake)
}

// ===================== DeregisterSolver =====================

func TestDeregisterSolver_Success(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Register first
	regMsg := &types.MsgRegisterSolver{
		Address:     solverAddr,
		Moniker:     "solver",
		StakeAmount: sdk.NewCoin("usyreen", math.NewInt(2000000000)),
	}
	require.NoError(t, k.RegisterSolver(ctx, regMsg))

	// Deregister
	deregMsg := &types.MsgDeregisterSolver{Address: solverAddr}
	err := k.DeregisterSolver(ctx, deregMsg)
	require.NoError(t, err)

	solver, found := k.GetSolver(ctx, solverAddr)
	require.True(t, found)
	require.False(t, solver.Active)
	// UnbondingHeight = current block (1) + DefaultSolverUnbondingBlocks (100) = 101
	require.Equal(t, int64(101), solver.UnbondingHeight)
}

func TestDeregisterSolver_NotRegistered(t *testing.T) {
	k, ctx := setupKeeper(t)

	msg := &types.MsgDeregisterSolver{Address: solverAddr}
	err := k.DeregisterSolver(ctx, msg)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrSolverNotRegistered)
}

func TestDeregisterSolver_AlreadyDeregistered(t *testing.T) {
	k, ctx := setupKeeper(t)

	regMsg := &types.MsgRegisterSolver{
		Address:     solverAddr,
		Moniker:     "solver",
		StakeAmount: sdk.NewCoin("usyreen", math.NewInt(2000000000)),
	}
	require.NoError(t, k.RegisterSolver(ctx, regMsg))

	deregMsg := &types.MsgDeregisterSolver{Address: solverAddr}
	require.NoError(t, k.DeregisterSolver(ctx, deregMsg))

	// Trying again should fail
	err := k.DeregisterSolver(ctx, deregMsg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "already deregistered")
}

// ===================== SubmitSolution =====================

func registerSolverAndIntent(t *testing.T, k *keeper.Keeper, ctx sdk.Context) string {
	t.Helper()
	// Register solver
	require.NoError(t, k.RegisterSolver(ctx, &types.MsgRegisterSolver{
		Address:     solverAddr,
		Moniker:     "solver",
		StakeAmount: sdk.NewCoin("usyreen", math.NewInt(2000000000)),
	}))

	// Submit intent from a different creator
	id, err := k.SubmitIntent(ctx, &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         swapIntentBody(),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:          sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10)),
		ExpiryBlocks: 50,
	})
	require.NoError(t, err)
	return id
}

func TestSubmitSolution_Success(t *testing.T) {
	k, ctx := setupKeeper(t)
	intentID := registerSolverAndIntent(t, k, ctx)

	msg := &types.MsgSubmitSolution{
		SolverAddr:      solverAddr,
		IntentID:        intentID,
		ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{"@type":"/test"}`)},
		ExpectedOutcome: json.RawMessage(`{"ok":true}`),
	}

	err := k.SubmitSolution(ctx, msg)
	require.NoError(t, err)

	// Intent should move to solving
	intent, found := k.GetIntent(ctx, intentID)
	require.True(t, found)
	require.Equal(t, types.StatusSolving, intent.Status)

	// Solution should be stored
	solutions := k.GetSolutionsForIntent(ctx, intentID)
	require.Len(t, solutions, 1)
	require.Equal(t, solverAddr, solutions[0].SolverAddr)
}

func TestSubmitSolution_SolverNotRegistered(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Submit intent but don't register solver
	id, err := k.SubmitIntent(ctx, &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1)),
		Tip:          sdk.Coins{},
		ExpiryBlocks: 10,
	})
	require.NoError(t, err)

	msg := &types.MsgSubmitSolution{
		SolverAddr:      solverAddr,
		IntentID:        id,
		ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{}`)},
		ExpectedOutcome: json.RawMessage(`{}`),
	}

	err = k.SubmitSolution(ctx, msg)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrSolverNotRegistered)
}

func TestSubmitSolution_IntentNotFound(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Register solver only, no intent
	require.NoError(t, k.RegisterSolver(ctx, &types.MsgRegisterSolver{
		Address:     solverAddr,
		Moniker:     "solver",
		StakeAmount: sdk.NewCoin("usyreen", math.NewInt(2000000000)),
	}))

	msg := &types.MsgSubmitSolution{
		SolverAddr:      solverAddr,
		IntentID:        "nonexistent",
		ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{}`)},
		ExpectedOutcome: json.RawMessage(`{}`),
	}

	err := k.SubmitSolution(ctx, msg)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrIntentNotFound)
}

func TestSubmitSolution_SelfSolvingPrevented(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Register solver with same address as creator
	require.NoError(t, k.RegisterSolver(ctx, &types.MsgRegisterSolver{
		Address:     creatorAddr, // same as intent creator
		Moniker:     "creator-solver",
		StakeAmount: sdk.NewCoin("usyreen", math.NewInt(2000000000)),
	}))

	id, err := k.SubmitIntent(ctx, &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1)),
		Tip:          sdk.Coins{},
		ExpiryBlocks: 10,
	})
	require.NoError(t, err)

	msg := &types.MsgSubmitSolution{
		SolverAddr:      creatorAddr,
		IntentID:        id,
		ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{}`)},
		ExpectedOutcome: json.RawMessage(`{}`),
	}

	err = k.SubmitSolution(ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "creator cannot solve their own intent")
}

func TestSubmitSolution_MaxSolutionsExceeded(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Set max solutions to 1
	p := types.DefaultParams()
	p.MaxSolutionsPerIntent = 1
	require.NoError(t, k.SetParams(ctx, p))

	intentID := registerSolverAndIntent(t, k, ctx)

	// First solution succeeds
	require.NoError(t, k.SubmitSolution(ctx, &types.MsgSubmitSolution{
		SolverAddr:      solverAddr,
		IntentID:        intentID,
		ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{}`)},
		ExpectedOutcome: json.RawMessage(`{}`),
	}))

	// Register second solver
	require.NoError(t, k.RegisterSolver(ctx, &types.MsgRegisterSolver{
		Address:     solverAddr2,
		Moniker:     "solver2",
		StakeAmount: sdk.NewCoin("usyreen", math.NewInt(2000000000)),
	}))

	// Second solution should fail
	err := k.SubmitSolution(ctx, &types.MsgSubmitSolution{
		SolverAddr:      solverAddr2,
		IntentID:        intentID,
		ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{}`)},
		ExpectedOutcome: json.RawMessage(`{}`),
	})
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrMaxSolutionsReached)
}

func TestSubmitSolution_SolvingWindowClosed(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Pin SolvingWindow=10 for this test (the module default is larger) so the
	// window closes at block 11 while the intent (expiry 51) is still live.
	swParams := k.GetParams(ctx)
	swParams.SolvingWindow = 10
	require.NoError(t, k.SetParams(ctx, swParams))

	// Use long expiry so the solving window check triggers before the expiry check.
	// SolvingWindow=10, so intent created at block 1 has window closing at block 11.
	// We need an intent that hasn't expired but whose solving window has closed.
	// Register solver
	require.NoError(t, k.RegisterSolver(ctx, &types.MsgRegisterSolver{
		Address:     solverAddr,
		Moniker:     "solver",
		StakeAmount: sdk.NewCoin("usyreen", math.NewInt(2000000000)),
	}))

	// Submit intent with long expiry (50 blocks), created at block 1, expiry = 51
	id, err := k.SubmitIntent(ctx, &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{"test":1}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:          sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10)),
		ExpiryBlocks: 50,
	})
	require.NoError(t, err)

	// Advance to block 12: past solving window (1+10=11) but before expiry (51)
	ctx = ctx.WithBlockHeight(12)

	msg := &types.MsgSubmitSolution{
		SolverAddr:      solverAddr,
		IntentID:        id,
		ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{}`)},
		ExpectedOutcome: json.RawMessage(`{}`),
	}

	err = k.SubmitSolution(ctx, msg)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrSolvingWindowClosed)
}

// ===================== FulfillIntent =====================

// A swap intent whose body OMITS min_output_amount leaves MinOutputAmount nil.
// verifySwapOutcome runs inside AutoFulfillIntents (BeginBlock), so a nil-Int
// comparison there (big.Int.Cmp(nil)) would panic and HALT THE CHAIN. This
// guards that the nil floor is normalized and fulfillment succeeds without panic.
func TestFulfillIntent_SwapBodyMissingMinOutput_NoPanic(t *testing.T) {
	k, ctx := setupKeeper(t)
	require.NoError(t, k.RegisterSolver(ctx, &types.MsgRegisterSolver{
		Address:     solverAddr,
		Moniker:     "solver",
		StakeAmount: sdk.NewCoin("usyreen", math.NewInt(2000000000)),
	}))
	// Body has denoms but NO min_output_amount → MinOutputAmount stays nil.
	id, err := k.SubmitIntent(ctx, &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{"input_denom":"usyreen","output_denom":"uusdc"}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:          sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10)),
		ExpiryBlocks: 50,
	})
	require.NoError(t, err, "nil min_output_amount is allowed (treated as zero floor)")

	require.NoError(t, k.SubmitSolution(ctx, &types.MsgSubmitSolution{
		SolverAddr:      solverAddr,
		IntentID:        id,
		ExecutionMsgs:   []json.RawMessage{bankSendExecMsg()},
		ExpectedOutcome: json.RawMessage(`{}`),
	}))

	// Must NOT panic despite the nil floor.
	require.NotPanics(t, func() {
		err = k.FulfillIntent(ctx, &types.MsgFulfillIntent{SolverAddr: solverAddr, IntentID: id})
	})
	require.NoError(t, err)
	intent, found := k.GetIntent(ctx, id)
	require.True(t, found)
	require.Equal(t, types.StatusFulfilled, intent.Status)
}

// A swap intent that explicitly declares a NEGATIVE floor is rejected at submit.
func TestSubmitIntent_SwapNegativeMinOutput_Rejected(t *testing.T) {
	k, ctx := setupKeeper(t)
	_, err := k.SubmitIntent(ctx, &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{"input_denom":"usyreen","output_denom":"uusdc","min_output_amount":"-5"}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		ExpiryBlocks: 50,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "must not be negative")
}

func TestFulfillIntent_Success(t *testing.T) {
	k, ctx := setupKeeper(t)
	intentID := registerSolverAndIntent(t, k, ctx)

	// Submit solution
	require.NoError(t, k.SubmitSolution(ctx, &types.MsgSubmitSolution{
		SolverAddr:      solverAddr,
		IntentID:        intentID,
		ExecutionMsgs:   []json.RawMessage{bankSendExecMsg()},
		ExpectedOutcome: json.RawMessage(`{}`),
	}))

	// Fulfill
	err := k.FulfillIntent(ctx, &types.MsgFulfillIntent{
		SolverAddr: solverAddr,
		IntentID:   intentID,
	})
	require.NoError(t, err)

	// Check intent is fulfilled
	intent, found := k.GetIntent(ctx, intentID)
	require.True(t, found)
	require.Equal(t, types.StatusFulfilled, intent.Status)
	require.Equal(t, solverAddr, intent.SolverAddr)

	// Check solver reputation increased
	solver, found := k.GetSolver(ctx, solverAddr)
	require.True(t, found)
	require.Equal(t, uint64(110), solver.ReputationScore) // 100 + 10
	require.Equal(t, uint64(1), solver.TotalSolved)
}

func TestFulfillIntent_NotFound(t *testing.T) {
	k, ctx := setupKeeper(t)

	err := k.FulfillIntent(ctx, &types.MsgFulfillIntent{
		SolverAddr: solverAddr,
		IntentID:   "nonexistent",
	})
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrIntentNotFound)
}

func TestFulfillIntent_AlreadyFulfilled(t *testing.T) {
	k, ctx := setupKeeper(t)
	intentID := registerSolverAndIntent(t, k, ctx)

	require.NoError(t, k.SubmitSolution(ctx, &types.MsgSubmitSolution{
		SolverAddr:      solverAddr,
		IntentID:        intentID,
		ExecutionMsgs:   []json.RawMessage{bankSendExecMsg()},
		ExpectedOutcome: json.RawMessage(`{}`),
	}))

	fulfillMsg := &types.MsgFulfillIntent{
		SolverAddr: solverAddr,
		IntentID:   intentID,
	}

	require.NoError(t, k.FulfillIntent(ctx, fulfillMsg))

	// Second attempt should fail
	err := k.FulfillIntent(ctx, fulfillMsg)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrIntentAlreadyFulfilled)
}

func TestFulfillIntent_Expired(t *testing.T) {
	k, ctx := setupKeeper(t)
	intentID := registerSolverAndIntent(t, k, ctx)

	require.NoError(t, k.SubmitSolution(ctx, &types.MsgSubmitSolution{
		SolverAddr:      solverAddr,
		IntentID:        intentID,
		ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{"@type":"/test"}`)},
		ExpectedOutcome: json.RawMessage(`{}`),
	}))

	// Advance past expiry (intent expires at block 51)
	ctx = ctx.WithBlockHeight(200)

	err := k.FulfillIntent(ctx, &types.MsgFulfillIntent{
		SolverAddr: solverAddr,
		IntentID:   intentID,
	})
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrIntentExpired)
}

func TestFulfillIntent_NoSolution(t *testing.T) {
	k, ctx := setupKeeper(t)
	intentID := registerSolverAndIntent(t, k, ctx)

	// Try to fulfill without submitting a solution
	err := k.FulfillIntent(ctx, &types.MsgFulfillIntent{
		SolverAddr: solverAddr,
		IntentID:   intentID,
	})
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrInvalidSolution)
}

// ===================== ExpireIntents =====================

func TestExpireIntents_RefundsCreator(t *testing.T) {
	k, ctx := setupKeeper(t)

	id, err := k.SubmitIntent(ctx, &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:          sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10)),
		ExpiryBlocks: 10,
	})
	require.NoError(t, err)

	// Verify intent is pending
	intent, found := k.GetIntent(ctx, id)
	require.True(t, found)
	require.Equal(t, types.StatusPending, intent.Status)

	// Advance past expiry (created at block 1, expiry = 1 + 10 = 11)
	ctx = ctx.WithBlockHeight(12)
	k.ExpireIntents(ctx)

	// Intent should now be expired
	intent, found = k.GetIntent(ctx, id)
	require.True(t, found)
	require.Equal(t, types.StatusExpired, intent.Status)
}

func TestExpireIntents_DoesNotExpireActiveIntents(t *testing.T) {
	k, ctx := setupKeeper(t)

	id, err := k.SubmitIntent(ctx, &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:          sdk.Coins{},
		ExpiryBlocks: 50,
	})
	require.NoError(t, err)

	// Still within expiry window
	ctx = ctx.WithBlockHeight(5)
	k.ExpireIntents(ctx)

	intent, found := k.GetIntent(ctx, id)
	require.True(t, found)
	require.Equal(t, types.StatusPending, intent.Status)
}

// ===================== CompleteSolverUnbonding =====================

func TestCompleteSolverUnbonding_ReturnsStake(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Register and deregister
	require.NoError(t, k.RegisterSolver(ctx, &types.MsgRegisterSolver{
		Address:     solverAddr,
		Moniker:     "solver",
		StakeAmount: sdk.NewCoin("usyreen", math.NewInt(2000000000)),
	}))
	require.NoError(t, k.DeregisterSolver(ctx, &types.MsgDeregisterSolver{Address: solverAddr}))

	// Before unbonding period ends, solver still exists
	ctx = ctx.WithBlockHeight(10)
	k.CompleteSolverUnbonding(ctx)
	_, found := k.GetSolver(ctx, solverAddr)
	require.True(t, found, "solver should still exist before unbonding completes")

	// After unbonding period (unbonding_height = 1 + 100 = 101)
	ctx = ctx.WithBlockHeight(102)
	k.CompleteSolverUnbonding(ctx)

	// Solver should be removed from store
	_, found = k.GetSolver(ctx, solverAddr)
	require.False(t, found, "solver should be removed after unbonding completes")
}

// ===================== GetPendingIntents =====================

func TestGetPendingIntents(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Create multiple intents
	for i := 0; i < 3; i++ {
		_, err := k.SubmitIntent(ctx, &types.MsgSubmitIntent{
			Creator:      creatorAddr,
			IntentType:   types.IntentTypeSwap,
			Body:         json.RawMessage(`{}`),
			MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1)),
			Tip:          sdk.Coins{},
			ExpiryBlocks: 50,
		})
		require.NoError(t, err)
	}

	// Manually set one intent to fulfilled
	intent, found := k.GetIntent(ctx, "2")
	require.True(t, found)
	intent.Status = types.StatusFulfilled
	k.SetIntent(ctx, intent)

	pending := k.GetPendingIntents(ctx)
	require.Len(t, pending, 2)
	for _, p := range pending {
		require.Equal(t, types.StatusPending, p.Status)
	}
}

// ===================== Genesis InitGenesis / ExportGenesis =====================

func TestGenesis_InitExportRoundtrip(t *testing.T) {
	k, ctx := setupKeeper(t)

	params := types.DefaultParams()
	params.SolvingWindow = 42

	genesisIn := types.GenesisState{
		Params: params,
		Intents: []types.Intent{
			{
				ID:         "10",
				Creator:    creatorAddr,
				IntentType: types.IntentTypeSwap,
				Body:       json.RawMessage(`{"a":1}`),
				MaxFee:     sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
				Tip:        sdk.NewCoins(sdk.NewInt64Coin("usyreen", 5)),
				Expiry:     200,
				Status:     types.StatusPending,
				CreatedAt:  1,
			},
			{
				ID:         "11",
				Creator:    creatorAddr,
				IntentType: types.IntentTypeTransfer,
				Body:       json.RawMessage(`{"b":2}`),
				MaxFee:     sdk.NewCoins(sdk.NewInt64Coin("usyreen", 50)),
				Tip:        sdk.Coins{},
				Expiry:     300,
				Status:     types.StatusFulfilled,
				CreatedAt:  5,
				SolverAddr: solverAddr,
			},
		},
		Solvers: []types.Solver{
			{
				Address:         solverAddr,
				Moniker:         "gen-solver",
				StakedAmount:    sdk.NewInt64Coin("usyreen", 5000),
				ReputationScore: 120,
				TotalSolved:     3,
				Active:          true,
				JoinedAt:        1,
			},
		},
	}

	k.InitGenesis(ctx, genesisIn)

	// Verify state
	got := k.GetParams(ctx)
	require.Equal(t, uint64(42), got.SolvingWindow)

	intent, found := k.GetIntent(ctx, "10")
	require.True(t, found)
	require.Equal(t, types.StatusPending, intent.Status)

	intent11, found := k.GetIntent(ctx, "11")
	require.True(t, found)
	require.Equal(t, types.StatusFulfilled, intent11.Status)

	solver, found := k.GetSolver(ctx, solverAddr)
	require.True(t, found)
	require.Equal(t, uint64(120), solver.ReputationScore)

	// Export and verify roundtrip
	exported := k.ExportGenesis(ctx)
	require.Equal(t, uint64(42), exported.Params.SolvingWindow)
	require.Len(t, exported.Intents, 2)
	require.Len(t, exported.Solvers, 1)

	// Verify exported intents have correct IDs
	exportedIDs := map[string]bool{}
	for _, i := range exported.Intents {
		exportedIDs[i.ID] = true
	}
	require.True(t, exportedIDs["10"])
	require.True(t, exportedIDs["11"])
}

func TestGenesis_EmptyState(t *testing.T) {
	k, ctx := setupKeeper(t)

	gs := *types.DefaultGenesis()
	k.InitGenesis(ctx, gs)

	exported := k.ExportGenesis(ctx)
	require.NoError(t, exported.Validate())
	// May be nil (no intents/solvers written), which is fine
	require.Equal(t, types.DefaultParams().SolvingWindow, exported.Params.SolvingWindow)
}

// ===================== Edge cases =====================

func TestSubmitSolution_IntentAlreadyFulfilled(t *testing.T) {
	k, ctx := setupKeeper(t)
	intentID := registerSolverAndIntent(t, k, ctx)

	// Manually set intent to fulfilled
	intent, _ := k.GetIntent(ctx, intentID)
	intent.Status = types.StatusFulfilled
	k.SetIntent(ctx, intent)

	err := k.SubmitSolution(ctx, &types.MsgSubmitSolution{
		SolverAddr:      solverAddr,
		IntentID:        intentID,
		ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{}`)},
		ExpectedOutcome: json.RawMessage(`{}`),
	})
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrIntentAlreadyFulfilled)
}

func TestSubmitSolution_IntentExpired(t *testing.T) {
	k, ctx := setupKeeper(t)
	intentID := registerSolverAndIntent(t, k, ctx)

	// Advance past intent expiry but within solving window conceptually
	// Intent expires at block 51, so go past that
	ctx = ctx.WithBlockHeight(52)

	err := k.SubmitSolution(ctx, &types.MsgSubmitSolution{
		SolverAddr:      solverAddr,
		IntentID:        intentID,
		ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{}`)},
		ExpectedOutcome: json.RawMessage(`{}`),
	})
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrIntentExpired)
}

func TestFulfillIntent_TerminalState(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Create intent and set to failed
	id, err := k.SubmitIntent(ctx, &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1)),
		Tip:          sdk.Coins{},
		ExpiryBlocks: 50,
	})
	require.NoError(t, err)

	intent, _ := k.GetIntent(ctx, id)
	intent.Status = types.StatusFailed
	k.SetIntent(ctx, intent)

	err = k.FulfillIntent(ctx, &types.MsgFulfillIntent{
		SolverAddr: solverAddr,
		IntentID:   id,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "terminal state")
}

func TestIntentCounter_Persistence(t *testing.T) {
	k, ctx := setupKeeper(t)

	msg := &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeCustom,
		Body:         json.RawMessage(`{}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1)),
		Tip:          sdk.Coins{},
		ExpiryBlocks: 10,
	}

	// Submit 5 intents
	for i := 0; i < 5; i++ {
		_, err := k.SubmitIntent(ctx, msg)
		require.NoError(t, err)
	}

	// Next should be 6
	id, err := k.SubmitIntent(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, "6", id)
}

func TestSolverStake_DenomMismatch(t *testing.T) {
	k, ctx := setupKeeper(t)

	// sdk.Coin.IsLT panics when denoms differ. Verify the keeper propagates that panic.
	msg := &types.MsgRegisterSolver{
		Address:     solverAddr,
		Moniker:     "solver",
		StakeAmount: sdk.NewInt64Coin("uatom", 999999),
	}

	require.Panics(t, func() {
		_ = k.RegisterSolver(ctx, msg)
	})
}

func TestExpireIntents_SolvingStatusAlsoExpires(t *testing.T) {
	k, ctx := setupKeeper(t)

	id, err := k.SubmitIntent(ctx, &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:          sdk.Coins{},
		ExpiryBlocks: 10,
	})
	require.NoError(t, err)

	// Set status to solving
	intent, _ := k.GetIntent(ctx, id)
	intent.Status = types.StatusSolving
	k.SetIntent(ctx, intent)

	// Advance past expiry
	ctx = ctx.WithBlockHeight(12)
	k.ExpireIntents(ctx)

	intent, found := k.GetIntent(ctx, id)
	require.True(t, found)
	require.Equal(t, types.StatusExpired, intent.Status)
}

func TestExpireIntents_DoesNotExpireFulfilledOrFailed(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Create two intents and set them to fulfilled and failed
	for _, status := range []types.IntentStatus{types.StatusFulfilled, types.StatusFailed} {
		id, err := k.SubmitIntent(ctx, &types.MsgSubmitIntent{
			Creator:      creatorAddr,
			IntentType:   types.IntentTypeSwap,
			Body:         json.RawMessage(`{}`),
			MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1)),
			Tip:          sdk.Coins{},
			ExpiryBlocks: 5,
		})
		require.NoError(t, err)
		intent, _ := k.GetIntent(ctx, id)
		intent.Status = status
		k.SetIntent(ctx, intent)
	}

	ctx = ctx.WithBlockHeight(100)
	k.ExpireIntents(ctx)

	// Both should retain their original status
	intent1, _ := k.GetIntent(ctx, "1")
	require.Equal(t, types.StatusFulfilled, intent1.Status)
	intent2, _ := k.GetIntent(ctx, "2")
	require.Equal(t, types.StatusFailed, intent2.Status)
}

func TestMathHelpers(t *testing.T) {
	val := uint64(1234567890)
	bz := types.Uint64ToBytes(val)
	require.Equal(t, val, types.BytesToUint64(bz))
}

func TestStoreKeys(t *testing.T) {
	require.Equal(t, []byte("intent/abc"), types.IntentKey("abc"))
	require.Equal(t, []byte("solver/xyz"), types.SolverKey("xyz"))
	require.Equal(t, []byte("solution/i1/s1"), types.SolutionKey("i1", "s1"))
	require.Equal(t, []byte("solution/i1/"), types.SolutionsByIntentPrefix("i1"))
	require.Equal(t, []byte("intent_counter"), types.IntentCounterKey())
}

func TestRegisterSolver_StakeLockFailure(t *testing.T) {
	k, ctx := setupKeeperWithBankErr(t, fmt.Errorf("insufficient funds"))

	msg := &types.MsgRegisterSolver{
		Address:     solverAddr,
		Moniker:     "solver",
		StakeAmount: sdk.NewCoin("usyreen", math.NewInt(2000000000)),
	}

	err := k.RegisterSolver(ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to lock solver stake")
}

func TestSubmitSolution_InactiveSolver(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Register and deregister solver
	require.NoError(t, k.RegisterSolver(ctx, &types.MsgRegisterSolver{
		Address:     solverAddr,
		Moniker:     "solver",
		StakeAmount: sdk.NewCoin("usyreen", math.NewInt(2000000000)),
	}))
	require.NoError(t, k.DeregisterSolver(ctx, &types.MsgDeregisterSolver{Address: solverAddr}))

	// Create intent
	id, err := k.SubmitIntent(ctx, &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1)),
		Tip:          sdk.Coins{},
		ExpiryBlocks: 50,
	})
	require.NoError(t, err)

	// Inactive solver tries to submit solution
	err = k.SubmitSolution(ctx, &types.MsgSubmitSolution{
		SolverAddr:      solverAddr,
		IntentID:        id,
		ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{}`)},
		ExpectedOutcome: json.RawMessage(`{}`),
	})
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrSolverNotRegistered)
}

// Ensure coins with Amount=0 don't cause lock issues (should skip locking)
func TestSubmitIntent_ZeroFees(t *testing.T) {
	k, ctx := setupKeeper(t)

	msg := &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{}`),
		MaxFee:       sdk.Coins{},
		Tip:          sdk.Coins{},
		ExpiryBlocks: 10,
	}

	id, err := k.SubmitIntent(ctx, msg)
	require.NoError(t, err)
	require.NotEmpty(t, id)
}

// Verify that sdk.Coin.IsLT with same denom works for min stake check.
func TestRegisterSolver_ExactMinStake(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Exact minimum (1000000000 usyreen = 1000 SYR)
	msg := &types.MsgRegisterSolver{
		Address:     solverAddr,
		Moniker:     "solver",
		StakeAmount: sdk.NewCoin("usyreen", math.NewInt(1000000000)),
	}

	err := k.RegisterSolver(ctx, msg)
	require.NoError(t, err)
}

// ===================== Trading Intent Types =====================

func TestSubmitIntent_LimitBuy(t *testing.T) {
	k, ctx := setupKeeper(t)

	body := types.LimitBuyIntent{
		InputDenom:      "usyreen",
		InputAmount:     math.NewInt(1000000),
		OutputDenom:     "uatom",
		TargetPrice:     math.LegacyNewDecWithPrec(5, 1),
		PoolID:          1,
		MinOutputAmount: math.NewInt(900000),
	}
	bodyBz, err := json.Marshal(body)
	require.NoError(t, err)

	msg := &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeLimitBuy,
		Body:         bodyBz,
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:          sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10)),
		ExpiryBlocks: 50,
	}
	require.NoError(t, msg.ValidateBasic())

	id, err := k.SubmitIntent(ctx, msg)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	intent, found := k.GetIntent(ctx, id)
	require.True(t, found)
	require.Equal(t, types.IntentTypeLimitBuy, intent.IntentType)
}

func TestSubmitIntent_StopLoss(t *testing.T) {
	k, ctx := setupKeeper(t)

	body := types.StopLossIntent{
		InputDenom:      "uatom",
		InputAmount:     math.NewInt(500000),
		OutputDenom:     "usyreen",
		StopPrice:       math.LegacyNewDecWithPrec(8, 1),
		PoolID:          2,
		MinOutputAmount: math.NewInt(350000),
	}
	bodyBz, err := json.Marshal(body)
	require.NoError(t, err)

	msg := &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeStopLoss,
		Body:         bodyBz,
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:          sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10)),
		ExpiryBlocks: 50,
	}
	require.NoError(t, msg.ValidateBasic())

	id, err := k.SubmitIntent(ctx, msg)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	intent, found := k.GetIntent(ctx, id)
	require.True(t, found)
	require.Equal(t, types.IntentTypeStopLoss, intent.IntentType)
}

// ===================== Body validation (defense in depth) =====================

// A trading body with a nil input_amount must be rejected cleanly at submit,
// NOT reach extractTradingInputCoins where sdk.NewCoin(nil) would panic.
func TestSubmitIntent_TradingNilInputAmount_Rejected(t *testing.T) {
	k, ctx := setupKeeper(t)
	// total_amount / input_amount omitted -> nil math.Int after unmarshal.
	msg := &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeLimitBuy,
		Body:         json.RawMessage(`{"input_denom":"usyreen","output_denom":"uatom","target_price":"0.5","pool_id":1}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		ExpiryBlocks: 50,
	}
	var err error
	require.NotPanics(t, func() { _, err = k.SubmitIntent(ctx, msg) })
	require.Error(t, err)
	require.Contains(t, err.Error(), "input_amount")
}

// A trading body with an empty denom must be rejected at submit (sdk.NewCoin
// would otherwise panic on an invalid denom).
func TestSubmitIntent_TradingEmptyDenom_Rejected(t *testing.T) {
	k, ctx := setupKeeper(t)
	msg := &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeLimitSell,
		Body:         json.RawMessage(`{"input_denom":"","output_denom":"uatom","input_amount":"1000","target_price":"0.5","pool_id":1}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		ExpiryBlocks: 50,
	}
	var err error
	require.NotPanics(t, func() { _, err = k.SubmitIntent(ctx, msg) })
	require.Error(t, err)
	require.Contains(t, err.Error(), "input_denom")
}

// A DCA body with a nil total_amount must be rejected at submit — otherwise it
// would panic in extractTradingInputCoins on submit and (if it somehow reached
// state) in tryDCA's QuoRaw during BeginBlock.
func TestSubmitIntent_DCANilTotalAmount_Rejected(t *testing.T) {
	k, ctx := setupKeeper(t)
	msg := &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeDCA,
		Body:         json.RawMessage(`{"input_denom":"usyreen","output_denom":"uatom","num_executions":5,"interval_blocks":10,"pool_id":1}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		ExpiryBlocks: 500,
	}
	var err error
	require.NotPanics(t, func() { _, err = k.SubmitIntent(ctx, msg) })
	require.Error(t, err)
	require.Contains(t, err.Error(), "total_amount")
}

// A DCA body with zero num_executions is rejected (would strand locked funds and
// risks a divide path in execution).
func TestSubmitIntent_DCAZeroNumExecutions_Rejected(t *testing.T) {
	k, ctx := setupKeeper(t)
	msg := &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeDCA,
		Body:         json.RawMessage(`{"input_denom":"usyreen","output_denom":"uatom","total_amount":"1000000","num_executions":0,"interval_blocks":10,"pool_id":1}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		ExpiryBlocks: 500,
	}
	_, err := k.SubmitIntent(ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "num_executions")
}

func TestGetIntentsByType(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Create a limit_buy intent
	buyBody, _ := json.Marshal(types.LimitBuyIntent{
		InputDenom:      "usyreen",
		InputAmount:     math.NewInt(1000000),
		OutputDenom:     "uatom",
		TargetPrice:     math.LegacyNewDec(1),
		PoolID:          1,
		MinOutputAmount: math.NewInt(900000),
	})
	id1, err := k.SubmitIntent(ctx, &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeLimitBuy,
		Body:         buyBody,
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:          sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10)),
		ExpiryBlocks: 50,
	})
	require.NoError(t, err)

	// Create a stop_loss intent
	slBody, _ := json.Marshal(types.StopLossIntent{
		InputDenom:      "uatom",
		InputAmount:     math.NewInt(500000),
		OutputDenom:     "usyreen",
		StopPrice:       math.LegacyNewDecWithPrec(5, 1),
		PoolID:          2,
		MinOutputAmount: math.NewInt(200000),
	})
	_, err = k.SubmitIntent(ctx, &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeStopLoss,
		Body:         slBody,
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:          sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10)),
		ExpiryBlocks: 50,
	})
	require.NoError(t, err)

	// Create another limit_buy intent
	id3, err := k.SubmitIntent(ctx, &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeLimitBuy,
		Body:         buyBody,
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:          sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10)),
		ExpiryBlocks: 50,
	})
	require.NoError(t, err)

	// Query by limit_buy — should return 2 intents
	limitBuys := k.GetIntentsByType(ctx, types.IntentTypeLimitBuy)
	require.Len(t, limitBuys, 2)
	require.Equal(t, id1, limitBuys[0].ID)
	require.Equal(t, id3, limitBuys[1].ID)

	// Query by stop_loss — should return 1 intent
	stopLosses := k.GetIntentsByType(ctx, types.IntentTypeStopLoss)
	require.Len(t, stopLosses, 1)

	// Query by dca — should return 0 intents
	dcas := k.GetIntentsByType(ctx, types.IntentTypeDCA)
	require.Len(t, dcas, 0)
}

// Terminal intents older than the retention window are pruned from state so the
// BeginBlock scan cost stays bounded (liveness DoS defense).
func TestExpireIntents_PrunesStaleTerminalIntents(t *testing.T) {
	k, ctx := setupKeeper(t)
	// A fulfilled intent created long ago (older than retention).
	old := types.Intent{
		ID: "old", Creator: creatorAddr, IntentType: types.IntentTypeSwap,
		Body: swapIntentBody(), Status: types.StatusFulfilled,
		CreatedAt: 1, Expiry: 10,
	}
	k.SetIntent(ctx, old)
	// A recent fulfilled intent (within retention).
	recent := types.Intent{
		ID: "recent", Creator: creatorAddr, IntentType: types.IntentTypeSwap,
		Body: swapIntentBody(), Status: types.StatusFulfilled,
		CreatedAt: 100, Expiry: 110,
	}
	k.SetIntent(ctx, recent)

	ctx = ctx.WithBlockHeight(1 + types.IntentPruneRetentionBlocks + 5)
	k.ExpireIntents(ctx)

	_, foundOld := k.GetIntent(ctx, "old")
	require.False(t, foundOld, "stale terminal intent should be pruned")
	_, foundRecent := k.GetIntent(ctx, "recent")
	require.True(t, foundRecent, "recent terminal intent should be retained")
}
