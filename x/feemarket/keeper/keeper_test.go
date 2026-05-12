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

	"syreen/x/feemarket/keeper"
	"syreen/x/feemarket/types"
)

// ---------------------------------------------------------------------------
// Mock BankKeeper
// ---------------------------------------------------------------------------

type mockBankKeeper struct {
	burnedCoins sdk.Coins
}

func (m *mockBankKeeper) BurnCoins(_ context.Context, _ string, amt sdk.Coins) error {
	m.burnedCoins = m.burnedCoins.Add(amt...)
	return nil
}

func (m *mockBankKeeper) SendCoinsFromModuleToModule(_ context.Context, _, _ string, _ sdk.Coins) error {
	return nil
}

func (m *mockBankKeeper) SendCoinsFromAccountToModule(_ context.Context, _ sdk.AccAddress, _ string, _ sdk.Coins) error {
	return nil
}

func (m *mockBankKeeper) SendCoinsFromModuleToAccount(_ context.Context, _ string, _ sdk.AccAddress, _ sdk.Coins) error {
	return nil
}

func (m *mockBankKeeper) GetBalance(_ context.Context, _ sdk.AccAddress, _ string) sdk.Coin {
	return sdk.NewCoin("usyreen", math.ZeroInt())
}

// ---------------------------------------------------------------------------
// Setup
// ---------------------------------------------------------------------------

func setupKeeper(t *testing.T) (keeper.Keeper, sdk.Context, *mockBankKeeper) {
	t.Helper()
	storeKey := storetypes.NewKVStoreKey("feemarket")
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	storeService := runtime.NewKVStoreService(storeKey)

	bk := &mockBankKeeper{burnedCoins: sdk.NewCoins()}
	k := keeper.NewKeeper(cdc, storeService, bk, "authority")
	return k, ctx, bk
}

// ---------------------------------------------------------------------------
// FeeState Tests
// ---------------------------------------------------------------------------

func TestGetSetFeeState_Roundtrip(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	state := types.FeeState{
		BaseFee:         math.LegacyNewDec(42),
		MinBaseFee:      math.LegacyNewDec(1),
		MaxBaseFee:      math.LegacyNewDec(9999),
		BlockGasTarget:  75_000_000,
		AdjustmentSpeed: math.LegacyNewDecWithPrec(25, 2),
	}

	require.NoError(t, k.SetFeeState(ctx, state))
	got := k.GetFeeState(ctx)

	require.True(t, state.BaseFee.Equal(got.BaseFee))
	require.True(t, state.MinBaseFee.Equal(got.MinBaseFee))
	require.True(t, state.MaxBaseFee.Equal(got.MaxBaseFee))
	require.Equal(t, state.BlockGasTarget, got.BlockGasTarget)
	require.True(t, state.AdjustmentSpeed.Equal(got.AdjustmentSpeed))
}

func TestGetFeeState_DefaultWhenEmpty(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	got := k.GetFeeState(ctx)
	def := types.DefaultFeeState()
	require.True(t, def.BaseFee.Equal(got.BaseFee))
}

// ---------------------------------------------------------------------------
// BaseFee Tests
// ---------------------------------------------------------------------------

func TestGetSetBaseFee(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	newFee := math.LegacyNewDec(123)
	require.NoError(t, k.SetBaseFee(ctx, newFee))
	got := k.GetBaseFee(ctx)
	require.True(t, newFee.Equal(got))
}

// ---------------------------------------------------------------------------
// PrevBlockGasUsed Tests
// ---------------------------------------------------------------------------

func TestGetSetPrevBlockGasUsed_Roundtrip(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	k.SetPrevBlockGasUsed(ctx, 12345678)
	got := k.GetPrevBlockGasUsed(ctx)
	require.Equal(t, int64(12345678), got)
}

func TestGetPrevBlockGasUsed_DefaultZero(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	got := k.GetPrevBlockGasUsed(ctx)
	require.Equal(t, int64(0), got)
}

func TestGetSetPrevBlockGasUsed_Zero(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	k.SetPrevBlockGasUsed(ctx, 0)
	got := k.GetPrevBlockGasUsed(ctx)
	require.Equal(t, int64(0), got)
}

func TestGetSetPrevBlockGasUsed_LargeValue(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	// Use a large value (just under MaxInt64)
	largeVal := uint64(1<<63 - 1)
	k.SetPrevBlockGasUsed(ctx, largeVal)
	got := k.GetPrevBlockGasUsed(ctx)
	require.Equal(t, int64(largeVal), got)
}

func TestGetPrevBlockGasUsed_OverflowCap(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	// Value exceeding MaxInt64 should be capped
	overflowVal := uint64(1<<63 + 100)
	k.SetPrevBlockGasUsed(ctx, overflowVal)
	got := k.GetPrevBlockGasUsed(ctx)
	require.Equal(t, int64(1<<63-1), got) // MaxInt64
}

// ---------------------------------------------------------------------------
// AdjustBaseFee Tests
// ---------------------------------------------------------------------------

func TestAdjustBaseFee_AtTarget_NoChange(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	require.NoError(t, k.SetFeeState(ctx, types.DefaultFeeState()))
	require.NoError(t, k.SetParams(ctx, types.DefaultParams()))

	state := types.DefaultFeeState()
	result := k.AdjustBaseFee(ctx, state.BlockGasTarget)
	require.True(t, state.BaseFee.Equal(result), "base fee should not change when at target")
}

func TestAdjustBaseFee_AboveTarget_Increase(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	require.NoError(t, k.SetFeeState(ctx, types.DefaultFeeState()))
	require.NoError(t, k.SetParams(ctx, types.DefaultParams()))

	state := types.DefaultFeeState()
	// Double the gas target
	result := k.AdjustBaseFee(ctx, state.BlockGasTarget*2)
	require.True(t, result.GT(state.BaseFee), "base fee should increase when gas > target")
}

func TestAdjustBaseFee_BelowTarget_Decrease(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	require.NoError(t, k.SetFeeState(ctx, types.DefaultFeeState()))
	require.NoError(t, k.SetParams(ctx, types.DefaultParams()))

	state := types.DefaultFeeState()
	// Half the gas target
	result := k.AdjustBaseFee(ctx, state.BlockGasTarget/2)
	require.True(t, result.LT(state.BaseFee), "base fee should decrease when gas < target")
}

func TestAdjustBaseFee_ClampToMin(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	state := types.DefaultFeeState()
	// Set base fee very close to min so a decrease would push it below
	state.BaseFee = state.MinBaseFee.Add(math.LegacyNewDecWithPrec(1, 6))
	require.NoError(t, k.SetFeeState(ctx, state))
	require.NoError(t, k.SetParams(ctx, types.DefaultParams()))

	// Zero gas used => max decrease
	result := k.AdjustBaseFee(ctx, 0)
	require.True(t, result.GTE(state.MinBaseFee), "base fee should not go below min")
}

func TestAdjustBaseFee_ClampToMax(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	state := types.DefaultFeeState()
	// Set base fee very close to max
	state.BaseFee = state.MaxBaseFee.Sub(math.LegacyNewDecWithPrec(1, 6))
	require.NoError(t, k.SetFeeState(ctx, state))
	require.NoError(t, k.SetParams(ctx, types.DefaultParams()))

	// Very high gas used => max increase
	result := k.AdjustBaseFee(ctx, state.BlockGasTarget*100)
	require.True(t, result.LTE(state.MaxBaseFee), "base fee should not exceed max")
}

func TestAdjustBaseFee_GasTargetZero(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	state := types.DefaultFeeState()
	state.BlockGasTarget = 0
	require.NoError(t, k.SetFeeState(ctx, state))
	require.NoError(t, k.SetParams(ctx, types.DefaultParams()))

	result := k.AdjustBaseFee(ctx, 1000)
	require.True(t, state.BaseFee.Equal(result), "should return unchanged base fee when gasTarget=0")
}

// ---------------------------------------------------------------------------
// CalculateFee Tests
// ---------------------------------------------------------------------------

func TestCalculateFee_DefaultLane(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	require.NoError(t, k.SetFeeState(ctx, types.DefaultFeeState()))

	lane := types.FeeLane{
		Name:               "default",
		BaseFeeMultiplier:  math.LegacyNewDec(1),
		MaxBlockGas:        50_000_000,
		CurrentUtilization: math.LegacyZeroDec(),
	}
	require.NoError(t, k.SetFeeLane(ctx, lane))

	fee, err := k.CalculateFee(ctx, "default", 100_000)
	require.NoError(t, err)

	baseFee := types.DefaultFeeState().BaseFee // 0.01
	expected := baseFee.Mul(math.LegacyNewDec(1)).Mul(math.LegacyNewDec(100_000))
	require.True(t, expected.Equal(fee))
}

func TestCalculateFee_HighMultiplierLane(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	require.NoError(t, k.SetFeeState(ctx, types.DefaultFeeState()))

	lane := types.FeeLane{
		Name:               "defi",
		BaseFeeMultiplier:  math.LegacyNewDecWithPrec(15, 1), // 1.5x
		MaxBlockGas:        30_000_000,
		CurrentUtilization: math.LegacyZeroDec(),
	}
	require.NoError(t, k.SetFeeLane(ctx, lane))

	fee, err := k.CalculateFee(ctx, "defi", 100_000)
	require.NoError(t, err)

	baseFee := types.DefaultFeeState().BaseFee
	expected := baseFee.Mul(math.LegacyNewDecWithPrec(15, 1)).Mul(math.LegacyNewDec(100_000))
	require.True(t, expected.Equal(fee))
}

func TestCalculateFee_LaneNotFound(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	require.NoError(t, k.SetFeeState(ctx, types.DefaultFeeState()))

	_, err := k.CalculateFee(ctx, "nonexistent", 100_000)
	require.Error(t, err)
	require.Contains(t, err.Error(), "fee lane not found")
}

// ---------------------------------------------------------------------------
// BurnBaseFee Tests
// ---------------------------------------------------------------------------

func TestBurnBaseFee_Enabled(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	params := types.DefaultParams()
	params.EnableFeeBurn = true
	params.BurnRatio = math.LegacyNewDecWithPrec(5, 1) // 50%
	require.NoError(t, k.SetParams(ctx, params))

	fees := sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(1000)))
	err := k.BurnBaseFee(ctx, fees)
	require.NoError(t, err)

	// 50% of 1000 = 500 burned
	require.True(t, bk.burnedCoins.AmountOf("usyreen").Equal(math.NewInt(500)))
}

func TestBurnBaseFee_Disabled(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	params := types.DefaultParams()
	params.EnableFeeBurn = false
	require.NoError(t, k.SetParams(ctx, params))

	fees := sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(1000)))
	err := k.BurnBaseFee(ctx, fees)
	require.NoError(t, err)

	// Nothing should be burned
	require.True(t, bk.burnedCoins.IsZero())
}

func TestBurnBaseFee_ZeroFees(t *testing.T) {
	k, ctx, bk := setupKeeper(t)

	params := types.DefaultParams()
	params.EnableFeeBurn = true
	require.NoError(t, k.SetParams(ctx, params))

	err := k.BurnBaseFee(ctx, sdk.NewCoins())
	require.NoError(t, err)
	require.True(t, bk.burnedCoins.IsZero())
}

// ---------------------------------------------------------------------------
// GetBurnStats Tests
// ---------------------------------------------------------------------------

func TestGetBurnStats_AfterBurns(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	params := types.DefaultParams()
	params.EnableFeeBurn = true
	params.BurnRatio = math.LegacyOneDec() // 100% burn
	require.NoError(t, k.SetParams(ctx, params))

	fees := sdk.NewCoins(sdk.NewCoin("usyreen", math.NewInt(500)))
	require.NoError(t, k.BurnBaseFee(ctx, fees))

	totalBurned, _, lastBurnAmount := k.GetBurnStats(ctx)
	require.True(t, totalBurned.AmountOf("usyreen").Equal(math.NewInt(500)))
	require.True(t, lastBurnAmount.AmountOf("usyreen").Equal(math.NewInt(500)))
}

func TestGetBurnStats_EmptyWhenNoBurns(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	totalBurned, height, lastBurnAmount := k.GetBurnStats(ctx)
	require.True(t, totalBurned.IsZero())
	require.Equal(t, int64(0), height)
	require.True(t, lastBurnAmount.IsZero())
}

// ---------------------------------------------------------------------------
// FeeLane Tests
// ---------------------------------------------------------------------------

func TestGetSetFeeLane_Roundtrip(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	lane := types.FeeLane{
		Name:               "testlane",
		BaseFeeMultiplier:  math.LegacyNewDecWithPrec(25, 1), // 2.5x
		MaxBlockGas:        10_000_000,
		CurrentUtilization: math.LegacyNewDecWithPrec(3, 1), // 0.3
	}

	require.NoError(t, k.SetFeeLane(ctx, lane))

	got, found := k.GetFeeLane(ctx, "testlane")
	require.True(t, found)
	require.Equal(t, lane.Name, got.Name)
	require.True(t, lane.BaseFeeMultiplier.Equal(got.BaseFeeMultiplier))
	require.Equal(t, lane.MaxBlockGas, got.MaxBlockGas)
	require.True(t, lane.CurrentUtilization.Equal(got.CurrentUtilization))
}

func TestGetFeeLane_NotFound(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	_, found := k.GetFeeLane(ctx, "nonexistent")
	require.False(t, found)
}

// ---------------------------------------------------------------------------
// Params Tests
// ---------------------------------------------------------------------------

func TestGetSetParams_Roundtrip(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	params := types.DefaultParams()
	params.BurnRatio = math.LegacyNewDecWithPrec(75, 2) // 0.75
	params.BaseFeeChangeDenominator = 16
	params.EnableFeeBurn = false

	require.NoError(t, k.SetParams(ctx, params))
	got := k.GetParams(ctx)

	require.True(t, params.BurnRatio.Equal(got.BurnRatio))
	require.Equal(t, params.BaseFeeChangeDenominator, got.BaseFeeChangeDenominator)
	require.Equal(t, params.EnableFeeBurn, got.EnableFeeBurn)
}

func TestGetParams_DefaultWhenEmpty(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	got := k.GetParams(ctx)
	def := types.DefaultParams()
	require.True(t, def.BurnRatio.Equal(got.BurnRatio))
	require.Equal(t, def.BaseFeeChangeDenominator, got.BaseFeeChangeDenominator)
}

// ---------------------------------------------------------------------------
// Genesis Tests
// ---------------------------------------------------------------------------

func TestInitExportGenesis_Roundtrip(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	gs := types.DefaultGenesis()
	gs.FeeState.BaseFee = math.LegacyNewDec(77)
	gs.Params.BurnRatio = math.LegacyNewDecWithPrec(9, 1) // 0.9

	k.InitGenesis(ctx, *gs)

	exported := k.ExportGenesis(ctx)
	require.True(t, gs.FeeState.BaseFee.Equal(exported.FeeState.BaseFee))
	require.True(t, gs.Params.BurnRatio.Equal(exported.Params.BurnRatio))

	// Verify lanes were stored
	for _, lane := range gs.Params.FeeLanes {
		got, found := k.GetFeeLane(ctx, lane.Name)
		require.True(t, found, "lane %s should exist after InitGenesis", lane.Name)
		require.True(t, lane.BaseFeeMultiplier.Equal(got.BaseFeeMultiplier))
	}
}

func TestInitGenesis_OverwritesExisting(t *testing.T) {
	k, ctx, _ := setupKeeper(t)

	// Init with defaults first
	k.InitGenesis(ctx, *types.DefaultGenesis())

	// Init again with custom values
	gs := types.DefaultGenesis()
	gs.FeeState.BaseFee = math.LegacyNewDec(999)
	k.InitGenesis(ctx, *gs)

	got := k.GetBaseFee(ctx)
	require.True(t, math.LegacyNewDec(999).Equal(got))
}
