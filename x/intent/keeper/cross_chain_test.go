package keeper_test

import (
	"context"
	"encoding/json"
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
	"github.com/stretchr/testify/require"

	"syreen/x/intent/keeper"
	"syreen/x/intent/types"
)

// ---------- Mock TransferKeeper ----------

type mockTransferKeeper struct {
	calls    []transferCall
	err      error
	sequence uint64
}

type transferCall struct {
	SourcePort    string
	SourceChannel string
	Token         sdk.Coin
	Sender        string
	Receiver      string
}

func (m *mockTransferKeeper) Transfer(_ context.Context, msg *types.IBCTransferMsg) (*types.IBCTransferResponse, error) {
	m.calls = append(m.calls, transferCall{
		SourcePort:    msg.SourcePort,
		SourceChannel: msg.SourceChannel,
		Token:         msg.Token,
		Sender:        msg.Sender,
		Receiver:      msg.Receiver,
	})
	if m.err != nil {
		return nil, m.err
	}
	m.sequence++
	return &types.IBCTransferResponse{Sequence: m.sequence}, nil
}

// ---------- Setup ----------

func setupCrossChainKeeper(t *testing.T) (*keeper.Keeper, sdk.Context, *mockDexKeeper, *mockTransferKeeper) {
	storeKey := storetypes.NewKVStoreKey("intent")
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, cmtproto.Header{Height: 10}, false, log.NewNopLogger())
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	storeService := runtime.NewKVStoreService(storeKey)

	router := &mockMsgRouter{
		handler: func(_ sdk.Context, _ sdk.Msg) (*sdk.Result, error) {
			return &sdk.Result{}, nil
		},
	}

	k := keeper.NewKeeper(cdc, storeService, &mockAccountKeeper{}, &mockBankKeeper{}, router, "authority")

	dex := newMockDexKeeper()
	transfer := &mockTransferKeeper{}
	k.SetDexKeeper(dex)
	k.SetTransferKeeper(transfer)

	params := types.DefaultParams()
	params.IntentExpiryBlocks = 1000
	require.NoError(t, k.SetParams(ctx, params))

	return k, ctx, dex, transfer
}

func crossChainSwapBody(t *testing.T) (json.RawMessage, types.CrossChainSwapIntent) {
	body := types.CrossChainSwapIntent{
		InputDenom:       "usyreen",
		InputAmount:      math.NewInt(1000000),
		OutputDenom:      "ibc/ATOM",
		MinOutputAmount:  math.ZeroInt(),
		PoolID:           1,
		IBCSourcePort:    "transfer",
		IBCSourceChannel: "channel-0",
		Receiver:         "cosmos1abc123",
		TimeoutBlocks:    100,
	}
	bz, err := json.Marshal(body)
	require.NoError(t, err)
	return bz, body
}

// ---------- Tests ----------

func TestCrossChainSwap_Success(t *testing.T) {
	k, ctx, dex, transfer := setupCrossChainKeeper(t)

	dex.swapOut = sdk.NewCoin("ibc/ATOM", math.NewInt(500000))

	bodyBz, _ := crossChainSwapBody(t)

	intent := types.Intent{
		ID:         "1",
		Creator:    creatorAddr,
		IntentType: types.IntentTypeCrossChainSwap,
		Body:       bodyBz,
		MaxFee:     sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1000)),
		Tip:        sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Expiry:     100,
		Status:     types.StatusPending,
		CreatedAt:  1,
	}
	k.SetIntent(ctx, intent)

	k.ExecuteTradingIntents(ctx)

	// Verify intent fulfilled
	updated, found := k.GetIntent(ctx, "1")
	require.True(t, found)
	require.Equal(t, types.StatusFulfilled, updated.Status)

	// Verify swap was called
	require.Len(t, dex.swapCalls, 1)
	require.Equal(t, "usyreen", dex.swapCalls[0].TokenIn.Denom)
	require.Equal(t, math.NewInt(1000000), dex.swapCalls[0].TokenIn.Amount)

	// Verify IBC transfer was called
	require.Len(t, transfer.calls, 1)
	require.Equal(t, "transfer", transfer.calls[0].SourcePort)
	require.Equal(t, "channel-0", transfer.calls[0].SourceChannel)
	require.Equal(t, "ibc/ATOM", transfer.calls[0].Token.Denom)
	require.Equal(t, math.NewInt(500000), transfer.calls[0].Token.Amount)
	require.Equal(t, "cosmos1abc123", transfer.calls[0].Receiver)

	// Verify body was updated with tracking state
	var updatedBody types.CrossChainSwapIntent
	require.NoError(t, json.Unmarshal(updated.Body, &updatedBody))
	require.True(t, updatedBody.SwapDone)
	require.True(t, updatedBody.IBCSent)
}

func TestCrossChainSwap_SwapFails(t *testing.T) {
	k, ctx, dex, transfer := setupCrossChainKeeper(t)

	dex.swapErr = types.ErrIntentNotFound // simulate swap failure

	bodyBz, _ := crossChainSwapBody(t)

	intent := types.Intent{
		ID:         "1",
		Creator:    creatorAddr,
		IntentType: types.IntentTypeCrossChainSwap,
		Body:       bodyBz,
		MaxFee:     sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1000)),
		Tip:        sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Expiry:     100,
		Status:     types.StatusPending,
		CreatedAt:  1,
	}
	k.SetIntent(ctx, intent)

	k.ExecuteTradingIntents(ctx)

	// Intent should be failed
	updated, found := k.GetIntent(ctx, "1")
	require.True(t, found)
	require.Equal(t, types.StatusFailed, updated.Status)

	// No IBC transfer should have been attempted
	require.Empty(t, transfer.calls)
}

func TestCrossChainSwap_IBCTransferFails(t *testing.T) {
	k, ctx, dex, transfer := setupCrossChainKeeper(t)

	dex.swapOut = sdk.NewCoin("ibc/ATOM", math.NewInt(500000))
	transfer.err = types.ErrCrossChainTransferFail

	bodyBz, _ := crossChainSwapBody(t)

	intent := types.Intent{
		ID:         "1",
		Creator:    creatorAddr,
		IntentType: types.IntentTypeCrossChainSwap,
		Body:       bodyBz,
		MaxFee:     sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1000)),
		Tip:        sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Expiry:     100,
		Status:     types.StatusPending,
		CreatedAt:  1,
	}
	k.SetIntent(ctx, intent)

	k.ExecuteTradingIntents(ctx)

	// Intent should be failed (swap succeeded but IBC failed)
	updated, found := k.GetIntent(ctx, "1")
	require.True(t, found)
	require.Equal(t, types.StatusFailed, updated.Status)

	// Swap was called
	require.Len(t, dex.swapCalls, 1)

	// IBC transfer was attempted but failed
	require.Len(t, transfer.calls, 1)
}

func TestCrossChainSwap_NoTransferKeeper(t *testing.T) {
	k, ctx, _ := setupTradingKeeper(t) // no transfer keeper set

	bodyBz, _ := crossChainSwapBody(t)

	intent := types.Intent{
		ID:         "1",
		Creator:    creatorAddr,
		IntentType: types.IntentTypeCrossChainSwap,
		Body:       bodyBz,
		MaxFee:     sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1000)),
		Tip:        sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Expiry:     100,
		Status:     types.StatusPending,
		CreatedAt:  1,
	}
	k.SetIntent(ctx, intent)

	// Should not panic, just skip
	k.ExecuteTradingIntents(ctx)

	updated, found := k.GetIntent(ctx, "1")
	require.True(t, found)
	require.Equal(t, types.StatusPending, updated.Status) // stays pending
}

func TestCrossChainSwap_ExtractInputCoins(t *testing.T) {
	k, ctx, _, _ := setupCrossChainKeeper(t)

	bodyBz, _ := crossChainSwapBody(t)

	msg := &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeCrossChainSwap,
		Body:         bodyBz,
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1000)),
		Tip:          sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		ExpiryBlocks: 100,
	}

	intentID, err := k.SubmitIntent(ctx, msg)
	require.NoError(t, err)
	require.NotEmpty(t, intentID)
}

func TestCrossChainSwap_InChain(t *testing.T) {
	k, ctx, dex, transfer := setupCrossChainKeeper(t)

	// Step 0: limit buy SYR
	dex.SetSpotPrice(1, "uusdc", "usyreen", math.LegacyMustNewDecFromStr("0.4"))
	dex.swapOut = sdk.NewCoin("ibc/ATOM", math.NewInt(500000))

	step0 := limitBuyStep(1, "usyreen", 1000000, "uusdc", "0.5")

	// Step 1: cross-chain swap to send output via IBC
	ccBody := types.CrossChainSwapIntent{
		InputDenom:       "placeholder",
		InputAmount:      math.NewInt(1),
		OutputDenom:      "ibc/ATOM",
		MinOutputAmount:  math.ZeroInt(),
		PoolID:           1,
		IBCSourcePort:    "transfer",
		IBCSourceChannel: "channel-0",
		Receiver:         "cosmos1xyz",
		TimeoutBlocks:    100,
	}
	ccBz, _ := json.Marshal(ccBody)
	step1 := types.ChainStep{
		IntentType:  types.IntentTypeCrossChainSwap,
		Body:        ccBz,
		Condition:   types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{Type: types.InputFromPrevious},
	}

	msg := makeChainMsg(t, []types.ChainStep{step0, step1})
	chainID, err := k.SubmitChain(ctx, msg)
	require.NoError(t, err)

	// Execute step 0
	k.ExecuteTradingIntents(ctx)

	chain, found := k.GetChain(ctx, chainID)
	require.True(t, found)
	require.Equal(t, types.StatusFulfilled, chain.Steps[0].Status)

	// Step 1 should be activated
	require.NotEmpty(t, chain.Steps[1].IntentID)

	// Execute step 1 (cross-chain swap)
	dex.swapCalls = nil
	dex.swapOut = sdk.NewCoin("ibc/ATOM", math.NewInt(400000))
	k.ExecuteTradingIntents(ctx)

	chain, _ = k.GetChain(ctx, chainID)
	require.Equal(t, types.ChainStatusCompleted, chain.Status)
	require.Len(t, transfer.calls, 1)
	require.Equal(t, "cosmos1xyz", transfer.calls[0].Receiver)
}
