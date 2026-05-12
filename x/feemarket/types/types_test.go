package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"syreen/x/feemarket/types"
)

// ---------------------------------------------------------------------------
// Params Tests
// ---------------------------------------------------------------------------

func TestDefaultParams(t *testing.T) {
	p := types.DefaultParams()
	require.NoError(t, p.Validate())
	require.Equal(t, math.LegacyNewDecWithPrec(1, 2), p.DefaultBaseFee)
	require.Equal(t, math.LegacyNewDecWithPrec(1, 4), p.MinBaseFee)
	require.Equal(t, math.LegacyNewDec(1000), p.MaxBaseFee)
	require.Equal(t, uint64(8), p.BaseFeeChangeDenominator)
	require.Equal(t, uint64(2), p.ElasticityMultiplier)
	require.Equal(t, math.LegacyNewDecWithPrec(8, 1), p.BurnRatio)
	require.True(t, p.EnableFeeBurn)
	require.Len(t, p.FeeLanes, 4)
}

func TestParamsValidate_Valid(t *testing.T) {
	p := types.DefaultParams()
	require.NoError(t, p.Validate())
}

func TestParamsValidate_NegativeDefaultBaseFee(t *testing.T) {
	p := types.DefaultParams()
	p.DefaultBaseFee = math.LegacyNewDec(-1)
	require.Error(t, p.Validate())
	require.Contains(t, p.Validate().Error(), "default base fee cannot be negative")
}

func TestParamsValidate_NegativeMinBaseFee(t *testing.T) {
	p := types.DefaultParams()
	p.MinBaseFee = math.LegacyNewDec(-1)
	require.Error(t, p.Validate())
	require.Contains(t, p.Validate().Error(), "min base fee cannot be negative")
}

func TestParamsValidate_NegativeMaxBaseFee(t *testing.T) {
	p := types.DefaultParams()
	p.MaxBaseFee = math.LegacyNewDec(-1)
	require.Error(t, p.Validate())
	require.Contains(t, p.Validate().Error(), "max base fee must be positive")
}

func TestParamsValidate_ZeroMaxBaseFee(t *testing.T) {
	p := types.DefaultParams()
	p.MaxBaseFee = math.LegacyZeroDec()
	require.Error(t, p.Validate())
	require.Contains(t, p.Validate().Error(), "max base fee must be positive")
}

func TestParamsValidate_MinGreaterThanMax(t *testing.T) {
	p := types.DefaultParams()
	p.MinBaseFee = math.LegacyNewDec(2000)
	p.MaxBaseFee = math.LegacyNewDec(1000)
	require.Error(t, p.Validate())
	require.Contains(t, p.Validate().Error(), "min base fee")
	require.Contains(t, p.Validate().Error(), "cannot be greater than max base fee")
}

func TestParamsValidate_ZeroDenominator(t *testing.T) {
	p := types.DefaultParams()
	p.BaseFeeChangeDenominator = 0
	require.Error(t, p.Validate())
	require.Contains(t, p.Validate().Error(), "base fee change denominator cannot be zero")
}

func TestParamsValidate_ZeroElasticityMultiplier(t *testing.T) {
	p := types.DefaultParams()
	p.ElasticityMultiplier = 0
	require.Error(t, p.Validate())
	require.Contains(t, p.Validate().Error(), "elasticity multiplier cannot be zero")
}

func TestParamsValidate_BurnRatioNegative(t *testing.T) {
	p := types.DefaultParams()
	p.BurnRatio = math.LegacyNewDec(-1)
	require.Error(t, p.Validate())
	require.Contains(t, p.Validate().Error(), "burn ratio must be between 0 and 1")
}

func TestParamsValidate_BurnRatioAboveOne(t *testing.T) {
	p := types.DefaultParams()
	p.BurnRatio = math.LegacyNewDecWithPrec(15, 1) // 1.5
	require.Error(t, p.Validate())
	require.Contains(t, p.Validate().Error(), "burn ratio must be between 0 and 1")
}

func TestParamsValidate_BurnRatioExactlyOne(t *testing.T) {
	p := types.DefaultParams()
	p.BurnRatio = math.LegacyOneDec()
	require.NoError(t, p.Validate())
}

func TestParamsValidate_BurnRatioZero(t *testing.T) {
	// H-11: BurnRatio of 0 should now fail validation (minimum is 50%)
	p := types.DefaultParams()
	p.BurnRatio = math.LegacyZeroDec()
	require.Error(t, p.Validate())
	require.Contains(t, p.Validate().Error(), "burn_ratio must be at least 0.5")
}

func TestParamsValidate_EmptyLaneName(t *testing.T) {
	p := types.DefaultParams()
	p.FeeLanes = []types.FeeLane{
		{
			Name:               "",
			BaseFeeMultiplier:  math.LegacyNewDec(1),
			MaxBlockGas:        50_000_000,
			CurrentUtilization: math.LegacyZeroDec(),
		},
	}
	require.Error(t, p.Validate())
	require.Contains(t, p.Validate().Error(), "fee lane name cannot be empty")
}

func TestParamsValidate_NegativeLaneMultiplier(t *testing.T) {
	p := types.DefaultParams()
	p.FeeLanes = []types.FeeLane{
		{
			Name:               "bad",
			BaseFeeMultiplier:  math.LegacyNewDec(-1),
			MaxBlockGas:        50_000_000,
			CurrentUtilization: math.LegacyZeroDec(),
		},
	}
	require.Error(t, p.Validate())
	require.Contains(t, p.Validate().Error(), "base fee multiplier must be positive")
}

func TestParamsValidate_ZeroLaneMultiplier(t *testing.T) {
	p := types.DefaultParams()
	p.FeeLanes = []types.FeeLane{
		{
			Name:               "zero",
			BaseFeeMultiplier:  math.LegacyZeroDec(),
			MaxBlockGas:        50_000_000,
			CurrentUtilization: math.LegacyZeroDec(),
		},
	}
	require.Error(t, p.Validate())
	require.Contains(t, p.Validate().Error(), "base fee multiplier must be positive")
}

func TestParamsValidate_NegativeMaxBlockGas(t *testing.T) {
	p := types.DefaultParams()
	p.FeeLanes = []types.FeeLane{
		{
			Name:               "bad",
			BaseFeeMultiplier:  math.LegacyNewDec(1),
			MaxBlockGas:        -1,
			CurrentUtilization: math.LegacyZeroDec(),
		},
	}
	require.Error(t, p.Validate())
	require.Contains(t, p.Validate().Error(), "max block gas must be positive")
}

// ---------------------------------------------------------------------------
// FeeState Tests
// ---------------------------------------------------------------------------

func TestDefaultFeeState(t *testing.T) {
	fs := types.DefaultFeeState()
	require.Equal(t, math.LegacyNewDecWithPrec(1, 2), fs.BaseFee)
	require.Equal(t, math.LegacyNewDecWithPrec(1, 4), fs.MinBaseFee)
	require.Equal(t, math.LegacyNewDec(1000), fs.MaxBaseFee)
	require.Equal(t, int64(50_000_000), fs.BlockGasTarget)
	require.Equal(t, math.LegacyNewDecWithPrec(125, 3), fs.AdjustmentSpeed)
}

func TestFeeState_ProtoMessage(t *testing.T) {
	fs := types.DefaultFeeState()
	fs.ProtoMessage()
	require.Equal(t, "syreen.feemarket.FeeState", fs.XXX_MessageName())
	require.Contains(t, fs.String(), "FeeState{baseFee=")
}

func TestFeeState_Reset(t *testing.T) {
	fs := types.DefaultFeeState()
	fs.Reset()
	require.True(t, fs.BaseFee.IsNil())
}

// ---------------------------------------------------------------------------
// FeeLane Tests
// ---------------------------------------------------------------------------

func TestDefaultFeeLanes(t *testing.T) {
	lanes := types.DefaultFeeLanes()
	require.Len(t, lanes, 4)

	expectedNames := []string{"default", "defi", "ibc", "governance"}
	for i, name := range expectedNames {
		require.Equal(t, name, lanes[i].Name)
	}

	// Check multipliers
	require.True(t, lanes[0].BaseFeeMultiplier.Equal(math.LegacyNewDec(1)))              // default: 1.0x
	require.True(t, lanes[1].BaseFeeMultiplier.Equal(math.LegacyNewDecWithPrec(15, 1)))  // defi: 1.5x
	require.True(t, lanes[2].BaseFeeMultiplier.Equal(math.LegacyNewDecWithPrec(12, 1)))  // ibc: 1.2x
	require.True(t, lanes[3].BaseFeeMultiplier.Equal(math.LegacyNewDecWithPrec(8, 1)))   // governance: 0.8x
}

func TestFeeLane_ProtoMessage(t *testing.T) {
	fl := &types.FeeLane{
		Name:              "test",
		BaseFeeMultiplier: math.LegacyNewDec(1),
	}
	fl.ProtoMessage()
	require.Equal(t, "syreen.feemarket.FeeLane", fl.XXX_MessageName())
	require.Contains(t, fl.String(), "FeeLane{test")
}

// ---------------------------------------------------------------------------
// Genesis Tests
// ---------------------------------------------------------------------------

func TestDefaultGenesis(t *testing.T) {
	gs := types.DefaultGenesis()
	require.NotNil(t, gs)
	require.NoError(t, gs.Validate())
	require.Equal(t, types.DefaultParams(), gs.Params)
	require.Equal(t, types.DefaultFeeState(), gs.FeeState)
}

func TestGenesisState_Validate_Valid(t *testing.T) {
	gs := types.DefaultGenesis()
	require.NoError(t, gs.Validate())
}

func TestGenesisState_Validate_InvalidParams(t *testing.T) {
	gs := types.DefaultGenesis()
	gs.Params.BaseFeeChangeDenominator = 0
	err := gs.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid params")
}

func TestGenesisState_Validate_NegativeBaseFee(t *testing.T) {
	gs := types.DefaultGenesis()
	gs.FeeState.BaseFee = math.LegacyNewDec(-1)
	err := gs.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "base fee cannot be negative")
}

func TestGenesisState_Validate_NegativeMinBaseFee(t *testing.T) {
	gs := types.DefaultGenesis()
	gs.FeeState.MinBaseFee = math.LegacyNewDec(-1)
	err := gs.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "min base fee cannot be negative")
}

func TestGenesisState_Validate_ZeroMaxBaseFee(t *testing.T) {
	gs := types.DefaultGenesis()
	gs.FeeState.MaxBaseFee = math.LegacyZeroDec()
	err := gs.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "max base fee must be positive")
}

func TestGenesisState_Validate_ZeroGasTarget(t *testing.T) {
	gs := types.DefaultGenesis()
	gs.FeeState.BlockGasTarget = 0
	err := gs.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "block gas target must be positive")
}

func TestGenesisState_Validate_NegativeGasTarget(t *testing.T) {
	gs := types.DefaultGenesis()
	gs.FeeState.BlockGasTarget = -1
	err := gs.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "block gas target must be positive")
}

func TestGenesisState_ProtoMessage(t *testing.T) {
	gs := types.DefaultGenesis()
	gs.ProtoMessage()
	require.Equal(t, "syreen.feemarket.GenesisState", gs.XXX_MessageName())
	require.Contains(t, gs.String(), "feemarket genesis")
}

// ---------------------------------------------------------------------------
// BurnRecord Tests
// ---------------------------------------------------------------------------

func TestBurnRecord_ProtoMessage(t *testing.T) {
	br := &types.BurnRecord{Height: 10}
	br.ProtoMessage()
	require.Equal(t, "syreen.feemarket.BurnRecord", br.XXX_MessageName())
	require.Contains(t, br.String(), "BurnRecord{height=10")
}

func TestBurnRecord_Reset(t *testing.T) {
	br := &types.BurnRecord{Height: 10}
	br.Reset()
	require.Equal(t, int64(0), br.Height)
}

// ---------------------------------------------------------------------------
// Keys Tests
// ---------------------------------------------------------------------------

func TestLaneConfigKey(t *testing.T) {
	key := types.LaneConfigKey("default")
	require.Equal(t, []byte("lane/default"), key)
}

func TestBurnRecordKey(t *testing.T) {
	key := types.BurnRecordKey(100)
	require.Equal(t, len("burnrecord/")+8, len(key))
}

func TestInt64ToBytes(t *testing.T) {
	bz := types.Int64ToBytes(1)
	require.Len(t, bz, 8)
	require.Equal(t, byte(1), bz[7])
}
