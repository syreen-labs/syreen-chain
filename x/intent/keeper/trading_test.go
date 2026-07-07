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

// ---------- Mock DexKeeper ----------

type mockDexKeeper struct {
	spotPrices map[string]math.LegacyDec // key: "poolID:denomIn:denomOut"
	swapOut    sdk.Coin
	swapErr    error
	swapCalls  []swapCall

	// AI signal/sentiment mock state — keyed by poolID.
	signalLabels    map[uint64]string
	signalScores    map[uint64]int64
	rsiValues       map[uint64]math.LegacyDec
	volatility      map[uint64]math.LegacyDec
	fearGreed       map[uint64]int64
	fearGreedLabels map[uint64]string
	fearGreedOK     map[uint64]bool
}

type swapCall struct {
	Sender   string
	PoolID   uint64
	TokenIn  sdk.Coin
	MinOut   math.Int
}

func newMockDexKeeper() *mockDexKeeper {
	return &mockDexKeeper{
		spotPrices:      make(map[string]math.LegacyDec),
		signalLabels:    make(map[uint64]string),
		signalScores:    make(map[uint64]int64),
		rsiValues:       make(map[uint64]math.LegacyDec),
		volatility:      make(map[uint64]math.LegacyDec),
		fearGreed:       make(map[uint64]int64),
		fearGreedLabels: make(map[uint64]string),
		fearGreedOK:     make(map[uint64]bool),
	}
}

// SetSignal configures the mock signal for a pool.
func (m *mockDexKeeper) SetSignal(poolID uint64, label string, score int64, rsi math.LegacyDec) {
	m.signalLabels[poolID] = label
	m.signalScores[poolID] = score
	m.rsiValues[poolID] = rsi
}

// SetVolatility sets the mock volatility score for a pool.
func (m *mockDexKeeper) SetVolatility(poolID uint64, vol math.LegacyDec) {
	m.volatility[poolID] = vol
}

// SetFearGreed sets the mock fear/greed index for a pool.
func (m *mockDexKeeper) SetFearGreed(poolID uint64, index int64, label string) {
	m.fearGreed[poolID] = index
	m.fearGreedLabels[poolID] = label
	m.fearGreedOK[poolID] = true
}

func (m *mockDexKeeper) GetSignalForPool(_ context.Context, poolID uint64) (string, int64, math.LegacyDec) {
	label, ok := m.signalLabels[poolID]
	if !ok {
		return "NEUTRAL", 50, math.LegacyNewDec(50)
	}
	return label, m.signalScores[poolID], m.rsiValues[poolID]
}

func (m *mockDexKeeper) GetSignalVolatility(_ context.Context, poolID uint64) math.LegacyDec {
	v, ok := m.volatility[poolID]
	if !ok {
		return math.LegacyZeroDec()
	}
	return v
}

func (m *mockDexKeeper) GetFearGreedIndex(_ context.Context, poolID uint64) (int64, string, bool) {
	if !m.fearGreedOK[poolID] {
		return 0, "", false
	}
	return m.fearGreed[poolID], m.fearGreedLabels[poolID], true
}

func (m *mockDexKeeper) priceKey(poolID uint64, denomIn, denomOut string) string {
	return string(rune(poolID)) + ":" + denomIn + ":" + denomOut
}

func (m *mockDexKeeper) SetSpotPrice(poolID uint64, denomIn, denomOut string, price math.LegacyDec) {
	key := m.priceKey(poolID, denomIn, denomOut)
	m.spotPrices[key] = price
}

func (m *mockDexKeeper) GetSpotPrice(_ context.Context, poolID uint64, denomIn, denomOut string) (math.LegacyDec, error) {
	key := m.priceKey(poolID, denomIn, denomOut)
	price, ok := m.spotPrices[key]
	if !ok {
		return math.LegacyDec{}, types.ErrIntentNotFound // reuse error
	}
	return price, nil
}

func (m *mockDexKeeper) Swap(_ context.Context, sender string, poolID uint64, tokenIn sdk.Coin, minTokenOut math.Int) (sdk.Coin, error) {
	m.swapCalls = append(m.swapCalls, swapCall{sender, poolID, tokenIn, minTokenOut})
	if m.swapErr != nil {
		return sdk.Coin{}, m.swapErr
	}
	return m.swapOut, nil
}

// ---------- Setup with DexKeeper ----------

func setupTradingKeeper(t *testing.T) (*keeper.Keeper, sdk.Context, *mockDexKeeper) {
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
	k.SetDexKeeper(dex)

	params := types.DefaultParams()
	params.IntentExpiryBlocks = 1000
	require.NoError(t, k.SetParams(ctx, params))

	return k, ctx, dex
}

func mustMarshal(t *testing.T, v interface{}) json.RawMessage {
	bz, err := json.Marshal(v)
	require.NoError(t, err)
	return bz
}

// ---------- Tests ----------

func TestLimitBuy_TriggersWhenPriceBelowTarget(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	body := types.LimitBuyIntent{
		InputDenom:      "usyreen",
		InputAmount:     math.NewInt(1000000),
		OutputDenom:     "uusdc",
		TargetPrice:     math.LegacyNewDecWithPrec(5, 1), // 0.5 usyreen per uusdc
		PoolID:          1,
		MinOutputAmount: math.NewInt(900000),
	}

	// Set spot price below target: 0.4 usyreen per uusdc (cheaper)
	dex.SetSpotPrice(1, "uusdc", "usyreen", math.LegacyNewDecWithPrec(4, 1))
	dex.swapOut = sdk.NewCoin("uusdc", math.NewInt(950000))

	// Create intent directly in store
	intent := types.Intent{
		ID:         "1",
		Creator:    creatorAddr,
		IntentType: types.IntentTypeLimitBuy,
		Body:       mustMarshal(t, body),
		MaxFee:     sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1000)),
		Tip:        sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Expiry:     100,
		Status:     types.StatusPending,
		CreatedAt:  1,
	}
	k.SetIntent(ctx, intent)

	// Execute
	k.ExecuteTradingIntents(ctx)

	// Verify intent was fulfilled
	updated, found := k.GetIntent(ctx, "1")
	require.True(t, found)
	require.Equal(t, types.StatusFulfilled, updated.Status)

	// Verify swap was called
	require.Len(t, dex.swapCalls, 1)
	require.Equal(t, uint64(1), dex.swapCalls[0].PoolID)
	require.Equal(t, "usyreen", dex.swapCalls[0].TokenIn.Denom)
	require.Equal(t, math.NewInt(1000000), dex.swapCalls[0].TokenIn.Amount)
}

func TestLimitBuy_DoesNotTriggerWhenPriceAboveTarget(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	body := types.LimitBuyIntent{
		InputDenom:      "usyreen",
		InputAmount:     math.NewInt(1000000),
		OutputDenom:     "uusdc",
		TargetPrice:     math.LegacyNewDecWithPrec(5, 1), // 0.5
		PoolID:          1,
		MinOutputAmount: math.NewInt(900000),
	}

	// Set spot price above target: 0.6 usyreen per uusdc (too expensive)
	dex.SetSpotPrice(1, "uusdc", "usyreen", math.LegacyNewDecWithPrec(6, 1))

	intent := types.Intent{
		ID:         "1",
		Creator:    creatorAddr,
		IntentType: types.IntentTypeLimitBuy,
		Body:       mustMarshal(t, body),
		MaxFee:     sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1000)),
		Tip:        sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Expiry:     100,
		Status:     types.StatusPending,
		CreatedAt:  1,
	}
	k.SetIntent(ctx, intent)

	k.ExecuteTradingIntents(ctx)

	// Should still be pending
	updated, found := k.GetIntent(ctx, "1")
	require.True(t, found)
	require.Equal(t, types.StatusPending, updated.Status)
	require.Empty(t, dex.swapCalls)
}

func TestLimitSell_TriggersWhenPriceAboveTarget(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	body := types.LimitSellIntent{
		InputDenom:      "usyreen",
		InputAmount:     math.NewInt(500000),
		OutputDenom:     "uusdc",
		TargetPrice:     math.LegacyNewDec(2), // 2 uusdc per usyreen
		PoolID:          1,
		MinOutputAmount: math.NewInt(900000),
	}

	// Set spot price above target: 2.5 uusdc per usyreen
	dex.SetSpotPrice(1, "usyreen", "uusdc", math.LegacyNewDecWithPrec(25, 1))
	dex.swapOut = sdk.NewCoin("uusdc", math.NewInt(1200000))

	intent := types.Intent{
		ID:         "1",
		Creator:    creatorAddr,
		IntentType: types.IntentTypeLimitSell,
		Body:       mustMarshal(t, body),
		MaxFee:     sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1000)),
		Tip:        sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Expiry:     100,
		Status:     types.StatusPending,
		CreatedAt:  1,
	}
	k.SetIntent(ctx, intent)

	k.ExecuteTradingIntents(ctx)

	updated, found := k.GetIntent(ctx, "1")
	require.True(t, found)
	require.Equal(t, types.StatusFulfilled, updated.Status)
	require.Len(t, dex.swapCalls, 1)
}

func TestStopLoss_TriggersWhenPriceBelowStop(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	body := types.StopLossIntent{
		InputDenom:      "usyreen",
		InputAmount:     math.NewInt(1000000),
		OutputDenom:     "uusdc",
		StopPrice:       math.LegacyNewDec(1), // 1 uusdc per usyreen stop
		PoolID:          1,
		MinOutputAmount: math.NewInt(800000),
	}

	// Set price below stop: 0.8 uusdc per usyreen (price crashed)
	dex.SetSpotPrice(1, "usyreen", "uusdc", math.LegacyNewDecWithPrec(8, 1))
	dex.swapOut = sdk.NewCoin("uusdc", math.NewInt(800000))

	intent := types.Intent{
		ID:         "1",
		Creator:    creatorAddr,
		IntentType: types.IntentTypeStopLoss,
		Body:       mustMarshal(t, body),
		MaxFee:     sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1000)),
		Tip:        sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Expiry:     100,
		Status:     types.StatusPending,
		CreatedAt:  1,
	}
	k.SetIntent(ctx, intent)

	k.ExecuteTradingIntents(ctx)

	updated, found := k.GetIntent(ctx, "1")
	require.True(t, found)
	require.Equal(t, types.StatusFulfilled, updated.Status)
}

func TestTakeProfit_TriggersWhenPriceAboveTarget(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	body := types.TakeProfitIntent{
		InputDenom:      "usyreen",
		InputAmount:     math.NewInt(500000),
		OutputDenom:     "uusdc",
		TargetPrice:     math.LegacyNewDec(5), // 5 uusdc per usyreen take-profit
		PoolID:          1,
		MinOutputAmount: math.NewInt(2000000),
	}

	// Set price above target: 6 uusdc per usyreen (moon!)
	dex.SetSpotPrice(1, "usyreen", "uusdc", math.LegacyNewDec(6))
	dex.swapOut = sdk.NewCoin("uusdc", math.NewInt(2800000))

	intent := types.Intent{
		ID:         "1",
		Creator:    creatorAddr,
		IntentType: types.IntentTypeTakeProfit,
		Body:       mustMarshal(t, body),
		MaxFee:     sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1000)),
		Tip:        sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Expiry:     100,
		Status:     types.StatusPending,
		CreatedAt:  1,
	}
	k.SetIntent(ctx, intent)

	k.ExecuteTradingIntents(ctx)

	updated, found := k.GetIntent(ctx, "1")
	require.True(t, found)
	require.Equal(t, types.StatusFulfilled, updated.Status)
}

func TestDCA_ExecutesOnSchedule(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	body := types.DCAIntent{
		InputDenom:     "usyreen",
		TotalAmount:    math.NewInt(3000000),
		OutputDenom:    "uusdc",
		NumExecutions:  3,
		IntervalBlocks: 10,
		PoolID:         1,
		ExecutedCount:  0,
		NextExecBlock:  10, // execute at block 10
	}

	// tryDCA fetches the spot price for per-tranche slippage protection; without a
	// configured price GetSpotPrice errors and the swap is skipped.
	dex.SetSpotPrice(1, "usyreen", "uusdc", math.LegacyNewDec(1))
	dex.swapOut = sdk.NewCoin("uusdc", math.NewInt(900000))

	intent := types.Intent{
		ID:         "1",
		Creator:    creatorAddr,
		IntentType: types.IntentTypeDCA,
		Body:       mustMarshal(t, body),
		MaxFee:     sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1000)),
		Tip:        sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Expiry:     100,
		Status:     types.StatusPending,
		CreatedAt:  1,
	}
	k.SetIntent(ctx, intent)

	// Execute at block 10 — should trigger first DCA swap
	k.ExecuteTradingIntents(ctx)

	updated, found := k.GetIntent(ctx, "1")
	require.True(t, found)
	// After first execution, still pending (2 more to go)
	require.Equal(t, types.StatusPending, updated.Status)

	// Verify swap was called with 1/3 of total amount
	require.Len(t, dex.swapCalls, 1)
	require.Equal(t, math.NewInt(1000000), dex.swapCalls[0].TokenIn.Amount) // 3M / 3 = 1M

	// Verify body was updated
	var updatedBody types.DCAIntent
	require.NoError(t, json.Unmarshal(updated.Body, &updatedBody))
	require.Equal(t, uint64(1), updatedBody.ExecutedCount)
	require.Equal(t, int64(20), updatedBody.NextExecBlock) // 10 + 10
	require.Equal(t, math.NewInt(2000000), updatedBody.TotalAmount)
}

func TestDCA_CompletesAndFulfills(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	// Last execution remaining
	body := types.DCAIntent{
		InputDenom:     "usyreen",
		TotalAmount:    math.NewInt(1000000),
		OutputDenom:    "uusdc",
		NumExecutions:  3,
		IntervalBlocks: 10,
		PoolID:         1,
		ExecutedCount:  2,
		NextExecBlock:  10,
	}

	// Configure spot price so tryDCA's slippage-protection price fetch succeeds.
	dex.SetSpotPrice(1, "usyreen", "uusdc", math.LegacyNewDec(1))
	dex.swapOut = sdk.NewCoin("uusdc", math.NewInt(900000))

	intent := types.Intent{
		ID:         "1",
		Creator:    creatorAddr,
		IntentType: types.IntentTypeDCA,
		Body:       mustMarshal(t, body),
		MaxFee:     sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1000)),
		Tip:        sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Expiry:     100,
		Status:     types.StatusPending,
		CreatedAt:  1,
	}
	k.SetIntent(ctx, intent)

	k.ExecuteTradingIntents(ctx)

	updated, found := k.GetIntent(ctx, "1")
	require.True(t, found)
	require.Equal(t, types.StatusFulfilled, updated.Status)
}

func TestDCA_DoesNotExecuteBeforeNextBlock(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	body := types.DCAIntent{
		InputDenom:     "usyreen",
		TotalAmount:    math.NewInt(2000000),
		OutputDenom:    "uusdc",
		NumExecutions:  2,
		IntervalBlocks: 10,
		PoolID:         1,
		ExecutedCount:  0,
		NextExecBlock:  20, // not ready at block 10
	}

	intent := types.Intent{
		ID:         "1",
		Creator:    creatorAddr,
		IntentType: types.IntentTypeDCA,
		Body:       mustMarshal(t, body),
		MaxFee:     sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1000)),
		Tip:        sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Expiry:     100,
		Status:     types.StatusPending,
		CreatedAt:  1,
	}
	k.SetIntent(ctx, intent)

	k.ExecuteTradingIntents(ctx)

	// Should not have swapped
	require.Empty(t, dex.swapCalls)

	updated, found := k.GetIntent(ctx, "1")
	require.True(t, found)
	require.Equal(t, types.StatusPending, updated.Status)
}

func TestNoDexKeeper_DoesNotPanic(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Intent exists but no dex keeper set — should be a no-op
	body := types.LimitBuyIntent{
		InputDenom:  "usyreen",
		InputAmount: math.NewInt(1000000),
		OutputDenom: "uusdc",
		TargetPrice: math.LegacyNewDec(1),
		PoolID:      1,
	}

	intent := types.Intent{
		ID:         "1",
		Creator:    creatorAddr,
		IntentType: types.IntentTypeLimitBuy,
		Body:       mustMarshal(t, body),
		MaxFee:     sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1000)),
		Tip:        sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Expiry:     100,
		Status:     types.StatusPending,
		CreatedAt:  1,
	}
	k.SetIntent(ctx, intent)

	// Should not panic
	k.ExecuteTradingIntents(ctx)

	updated, found := k.GetIntent(ctx, "1")
	require.True(t, found)
	require.Equal(t, types.StatusPending, updated.Status)
}

func TestExtractTradingInputCoins(t *testing.T) {
	k, ctx, _ := setupTradingKeeper(t)

	// Submit a limit_buy intent and verify input tokens are parsed
	body := types.LimitBuyIntent{
		InputDenom:  "usyreen",
		InputAmount: math.NewInt(5000000),
		OutputDenom: "uusdc",
		TargetPrice: math.LegacyNewDec(1),
		PoolID:      1,
	}

	msg := &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeLimitBuy,
		Body:         mustMarshal(t, body),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1000)),
		Tip:          sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		ExpiryBlocks: 100,
	}

	// This will try to lock both fee+tip AND input tokens
	// With mock bank keeper it always succeeds
	intentID, err := k.SubmitIntent(ctx, msg)
	require.NoError(t, err)
	require.NotEmpty(t, intentID)
}

func TestStopLoss_DoesNotTriggerAboveStopPrice(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	body := types.StopLossIntent{
		InputDenom:      "usyreen",
		InputAmount:     math.NewInt(1000000),
		OutputDenom:     "uusdc",
		StopPrice:       math.LegacyNewDec(1),
		PoolID:          1,
		MinOutputAmount: math.NewInt(800000),
	}

	// Price above stop: 1.5 uusdc per usyreen (still healthy)
	dex.SetSpotPrice(1, "usyreen", "uusdc", math.LegacyNewDecWithPrec(15, 1))

	intent := types.Intent{
		ID:         "1",
		Creator:    creatorAddr,
		IntentType: types.IntentTypeStopLoss,
		Body:       mustMarshal(t, body),
		MaxFee:     sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1000)),
		Tip:        sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Expiry:     100,
		Status:     types.StatusPending,
		CreatedAt:  1,
	}
	k.SetIntent(ctx, intent)

	k.ExecuteTradingIntents(ctx)

	updated, found := k.GetIntent(ctx, "1")
	require.True(t, found)
	require.Equal(t, types.StatusPending, updated.Status)
	require.Empty(t, dex.swapCalls)
}

func TestTWAP_UsesDAMechanics(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	body := types.DCAIntent{
		InputDenom:     "usyreen",
		TotalAmount:    math.NewInt(2000000),
		OutputDenom:    "uusdc",
		NumExecutions:  2,
		IntervalBlocks: 5,
		PoolID:         1,
		ExecutedCount:  0,
		NextExecBlock:  10,
	}

	// Configure spot price so tryDCA's slippage-protection price fetch succeeds.
	dex.SetSpotPrice(1, "usyreen", "uusdc", math.LegacyNewDec(1))
	dex.swapOut = sdk.NewCoin("uusdc", math.NewInt(900000))

	intent := types.Intent{
		ID:         "1",
		Creator:    creatorAddr,
		IntentType: types.IntentTypeTWAP,
		Body:       mustMarshal(t, body),
		MaxFee:     sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1000)),
		Tip:        sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Expiry:     100,
		Status:     types.StatusPending,
		CreatedAt:  1,
	}
	k.SetIntent(ctx, intent)

	k.ExecuteTradingIntents(ctx)

	require.Len(t, dex.swapCalls, 1)
	require.Equal(t, math.NewInt(1000000), dex.swapCalls[0].TokenIn.Amount)
}
