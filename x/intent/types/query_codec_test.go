package types_test

import (
	"encoding/json"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"syreen/x/intent/types"
)

// These tests prove the hand-rolled proto codecs on the intent query
// request/response types actually round-trip their fields. Before the codecs
// were added, the reflection-based proto codec produced EMPTY output for these
// tag-less structs, so every field was silently dropped over gRPC/REST and
// clients saw blank responses. A round-trip through Marshal -> Unmarshal that
// preserves the data is exactly what the reflection path failed to do.

func sampleIntent(id string) types.Intent {
	return types.Intent{
		ID:         id,
		Creator:    "syreen1creator",
		IntentType: types.IntentTypeSwap,
		Body:       json.RawMessage(`{"input_denom":"usyreen","output_denom":"uatom"}`),
		MaxFee:     sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Status:     types.StatusPending,
		Expiry:     500,
		CreatedAt:  10,
	}
}

// TestMsgCancelIntent_RoundTrip proves the hand-rolled wire codec for the new
// MsgCancelIntent preserves both fields. A wrong Marshal here would silently drop
// the signer (creator) or intent_id on broadcast.
func TestMsgCancelIntent_RoundTrip(t *testing.T) {
	orig := &types.MsgCancelIntent{Creator: "syreen1creator", IntentID: "intent-99"}
	bz, err := orig.Marshal()
	require.NoError(t, err)
	require.NotEmpty(t, bz)

	var got types.MsgCancelIntent
	require.NoError(t, got.Unmarshal(bz))
	require.Equal(t, orig.Creator, got.Creator)
	require.Equal(t, orig.IntentID, got.IntentID)
	require.Equal(t, orig.Size(), len(bz))
}

func TestQueryIntentRequest_RoundTrip(t *testing.T) {
	orig := &types.QueryIntentRequest{IntentID: "intent-42"}
	bz, err := orig.Marshal()
	require.NoError(t, err)
	require.NotEmpty(t, bz)

	var got types.QueryIntentRequest
	require.NoError(t, got.Unmarshal(bz))
	require.Equal(t, orig.IntentID, got.IntentID)
}

func TestQueryIntentResponse_RoundTrip(t *testing.T) {
	orig := &types.QueryIntentResponse{Intent: sampleIntent("intent-1"), Found: true}
	bz, err := orig.Marshal()
	require.NoError(t, err)

	var got types.QueryIntentResponse
	require.NoError(t, got.Unmarshal(bz))
	require.True(t, got.Found)
	require.Equal(t, "intent-1", got.Intent.ID)
	require.Equal(t, orig.Intent.IntentType, got.Intent.IntentType)
	require.JSONEq(t, string(orig.Intent.Body), string(got.Intent.Body))
}

func TestQueryIntentResponse_NotFound(t *testing.T) {
	orig := &types.QueryIntentResponse{Found: false}
	bz, err := orig.Marshal()
	require.NoError(t, err)

	var got types.QueryIntentResponse
	require.NoError(t, got.Unmarshal(bz))
	require.False(t, got.Found)
}

func TestQueryIntentsResponse_RoundTrip(t *testing.T) {
	orig := &types.QueryIntentsResponse{Intents: []types.Intent{sampleIntent("a"), sampleIntent("b")}}
	bz, err := orig.Marshal()
	require.NoError(t, err)

	var got types.QueryIntentsResponse
	require.NoError(t, got.Unmarshal(bz))
	require.Len(t, got.Intents, 2)
	require.Equal(t, "a", got.Intents[0].ID)
	require.Equal(t, "b", got.Intents[1].ID)
}

func TestQueryIntentsByType_RoundTrip(t *testing.T) {
	req := &types.QueryIntentsByTypeRequest{IntentType: types.IntentTypeLimitBuy}
	bz, err := req.Marshal()
	require.NoError(t, err)
	var gotReq types.QueryIntentsByTypeRequest
	require.NoError(t, gotReq.Unmarshal(bz))
	require.Equal(t, types.IntentTypeLimitBuy, gotReq.IntentType)

	resp := &types.QueryIntentsByTypeResponse{Intents: []types.Intent{sampleIntent("x")}}
	rbz, err := resp.Marshal()
	require.NoError(t, err)
	var gotResp types.QueryIntentsByTypeResponse
	require.NoError(t, gotResp.Unmarshal(rbz))
	require.Len(t, gotResp.Intents, 1)
	require.Equal(t, "x", gotResp.Intents[0].ID)
}

func TestQuerySolverRoundTrip(t *testing.T) {
	req := &types.QuerySolverRequest{Address: "syreen1solver"}
	bz, err := req.Marshal()
	require.NoError(t, err)
	var gotReq types.QuerySolverRequest
	require.NoError(t, gotReq.Unmarshal(bz))
	require.Equal(t, "syreen1solver", gotReq.Address)

	resp := &types.QuerySolverResponse{
		Solver: types.Solver{Address: "syreen1solver", Moniker: "alice", ReputationScore: 100, Active: true},
		Found:  true,
	}
	rbz, err := resp.Marshal()
	require.NoError(t, err)
	var gotResp types.QuerySolverResponse
	require.NoError(t, gotResp.Unmarshal(rbz))
	require.True(t, gotResp.Found)
	require.Equal(t, "alice", gotResp.Solver.Moniker)
	require.Equal(t, uint64(100), gotResp.Solver.ReputationScore)
}

func TestQueryParamsResponse_RoundTrip(t *testing.T) {
	orig := &types.QueryParamsResponse{Params: types.DefaultParams()}
	bz, err := orig.Marshal()
	require.NoError(t, err)

	var got types.QueryParamsResponse
	require.NoError(t, got.Unmarshal(bz))
	require.Equal(t, orig.Params, got.Params)
}

func TestQueryStrategyRoundTrip(t *testing.T) {
	req := &types.QueryStrategyRequest{StrategyID: "strat-9"}
	bz, err := req.Marshal()
	require.NoError(t, err)
	var gotReq types.QueryStrategyRequest
	require.NoError(t, gotReq.Unmarshal(bz))
	require.Equal(t, "strat-9", gotReq.StrategyID)

	resp := &types.QueryStrategyResponse{
		Strategy: types.Strategy{ID: "strat-9", Creator: "syreen1c", TemplateName: "dca", TotalBudget: math.NewInt(1000)},
		Chain:    types.IntentChain{ID: "chain-9", Creator: "syreen1c"},
		Found:    true,
	}
	rbz, err := resp.Marshal()
	require.NoError(t, err)
	var gotResp types.QueryStrategyResponse
	require.NoError(t, gotResp.Unmarshal(rbz))
	require.True(t, gotResp.Found)
	require.Equal(t, "strat-9", gotResp.Strategy.ID)
	require.Equal(t, "dca", gotResp.Strategy.TemplateName)
	require.Equal(t, math.NewInt(1000), gotResp.Strategy.TotalBudget)
	require.Equal(t, "chain-9", gotResp.Chain.ID)
}

func TestQueryStrategiesResponse_RoundTrip(t *testing.T) {
	orig := &types.QueryStrategiesResponse{Strategies: []types.Strategy{
		{ID: "s1", TotalBudget: math.NewInt(1)},
		{ID: "s2", TotalBudget: math.NewInt(2)},
	}}
	bz, err := orig.Marshal()
	require.NoError(t, err)

	var got types.QueryStrategiesResponse
	require.NoError(t, got.Unmarshal(bz))
	require.Len(t, got.Strategies, 2)
	require.Equal(t, "s1", got.Strategies[0].ID)
	require.Equal(t, "s2", got.Strategies[1].ID)
}

func TestQueryStrategyTemplatesResponse_RoundTrip(t *testing.T) {
	orig := &types.QueryStrategyTemplatesResponse{Templates: []types.StrategyTemplate{
		{Name: "dca", Description: "dollar cost average", MinSteps: 1, MaxSteps: 10},
	}}
	bz, err := orig.Marshal()
	require.NoError(t, err)

	var got types.QueryStrategyTemplatesResponse
	require.NoError(t, got.Unmarshal(bz))
	require.Len(t, got.Templates, 1)
	require.Equal(t, "dca", got.Templates[0].Name)
	require.Equal(t, 10, got.Templates[0].MaxSteps)
}
