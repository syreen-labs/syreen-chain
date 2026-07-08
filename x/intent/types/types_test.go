package types_test

import (
	"encoding/json"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"syreen/x/intent/types"
)

// Valid bech32 addresses (cosmos prefix works for ValidateBasic since AccAddressFromBech32 is prefix-agnostic in tests)
var (
	validAddr1 = sdk.AccAddress([]byte("addr1_______________")).String()
	validAddr2 = sdk.AccAddress([]byte("addr2_______________")).String()
)

// ===================== MsgSubmitIntent =====================

func validMsgSubmitIntent() *types.MsgSubmitIntent {
	return &types.MsgSubmitIntent{
		Creator:      validAddr1,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{"foo":"bar"}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:          sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10)),
		ExpiryBlocks: 50,
	}
}

func TestMsgSubmitIntent_Valid(t *testing.T) {
	msg := validMsgSubmitIntent()
	require.NoError(t, msg.ValidateBasic())
}

func TestMsgSubmitIntent_InvalidCreator(t *testing.T) {
	msg := validMsgSubmitIntent()
	msg.Creator = "bad_address"
	require.Error(t, msg.ValidateBasic())
	require.Contains(t, msg.ValidateBasic().Error(), "invalid creator address")
}

func TestMsgSubmitIntent_InvalidIntentType(t *testing.T) {
	msg := validMsgSubmitIntent()
	msg.IntentType = "unknown_type"
	err := msg.ValidateBasic()
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrIntentTypeUnsupported)
}

func TestMsgSubmitIntent_AllValidIntentTypes(t *testing.T) {
	for _, it := range []string{
		types.IntentTypeSwap, types.IntentTypeTransfer, types.IntentTypeDeFi, types.IntentTypeCustom,
		types.IntentTypeLimitBuy, types.IntentTypeLimitSell, types.IntentTypeStopLoss,
		types.IntentTypeTakeProfit, types.IntentTypeTWAP, types.IntentTypeDCA,
	} {
		msg := validMsgSubmitIntent()
		msg.IntentType = it
		require.NoError(t, msg.ValidateBasic(), "intent type %s should be valid", it)
	}
}

func TestMsgSubmitIntent_EmptyBody(t *testing.T) {
	msg := validMsgSubmitIntent()
	msg.Body = nil
	err := msg.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "body cannot be empty")
}

func TestMsgSubmitIntent_InvalidMaxFee(t *testing.T) {
	msg := validMsgSubmitIntent()
	msg.MaxFee = sdk.Coins{sdk.Coin{Denom: "usyreen", Amount: math.NewInt(-1)}}
	err := msg.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid max fee")
}

func TestMsgSubmitIntent_ZeroExpiryBlocks(t *testing.T) {
	msg := validMsgSubmitIntent()
	msg.ExpiryBlocks = 0
	err := msg.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "expiry blocks must be greater than 0")
}

// Note: message signer resolution is no longer done via a GetSigners() method
// (removed in Cosmos SDK v0.53). Signers are resolved through the protobuf
// signing context registered in app/encoding.go (DefineCustomGetSigners) and
// are exercised end-to-end. The former per-message GetSigners unit tests were
// removed here because the method they asserted no longer exists.

// ===================== MsgRegisterSolver =====================

func validMsgRegisterSolver() *types.MsgRegisterSolver {
	return &types.MsgRegisterSolver{
		Address:     validAddr1,
		Moniker:     "test-solver",
		StakeAmount: sdk.NewInt64Coin("usyreen", 2000),
	}
}

func TestMsgRegisterSolver_Valid(t *testing.T) {
	msg := validMsgRegisterSolver()
	require.NoError(t, msg.ValidateBasic())
}

func TestMsgRegisterSolver_InvalidAddress(t *testing.T) {
	msg := validMsgRegisterSolver()
	msg.Address = "notanaddress"
	err := msg.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid solver address")
}

func TestMsgRegisterSolver_EmptyMoniker(t *testing.T) {
	msg := validMsgRegisterSolver()
	msg.Moniker = ""
	err := msg.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "moniker cannot be empty")
}

func TestMsgRegisterSolver_ZeroStake(t *testing.T) {
	msg := validMsgRegisterSolver()
	msg.StakeAmount = sdk.NewInt64Coin("usyreen", 0)
	err := msg.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "stake amount must be positive")
}

func TestMsgRegisterSolver_NegativeStake(t *testing.T) {
	msg := validMsgRegisterSolver()
	msg.StakeAmount = sdk.Coin{Denom: "usyreen", Amount: math.NewInt(-5)}
	err := msg.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "stake amount must be positive")
}

// ===================== MsgDeregisterSolver =====================

func TestMsgDeregisterSolver_Valid(t *testing.T) {
	msg := &types.MsgDeregisterSolver{Address: validAddr1}
	require.NoError(t, msg.ValidateBasic())
}

func TestMsgDeregisterSolver_InvalidAddress(t *testing.T) {
	msg := &types.MsgDeregisterSolver{Address: "invalid"}
	err := msg.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid solver address")
}

// ===================== MsgSubmitSolution =====================

func validMsgSubmitSolution() *types.MsgSubmitSolution {
	return &types.MsgSubmitSolution{
		SolverAddr:      validAddr1,
		IntentID:        "1",
		ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{"@type":"/test"}`)},
		ExpectedOutcome: json.RawMessage(`{"ok":true}`),
	}
}

func TestMsgSubmitSolution_Valid(t *testing.T) {
	msg := validMsgSubmitSolution()
	require.NoError(t, msg.ValidateBasic())
}

func TestMsgSubmitSolution_InvalidSolver(t *testing.T) {
	msg := validMsgSubmitSolution()
	msg.SolverAddr = "bad"
	err := msg.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid solver address")
}

func TestMsgSubmitSolution_EmptyIntentID(t *testing.T) {
	msg := validMsgSubmitSolution()
	msg.IntentID = ""
	err := msg.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "intent ID cannot be empty")
}

func TestMsgSubmitSolution_EmptyExecutionMsgs(t *testing.T) {
	msg := validMsgSubmitSolution()
	msg.ExecutionMsgs = nil
	err := msg.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "execution messages cannot be empty")
}

// ===================== MsgFulfillIntent =====================

func TestMsgFulfillIntent_Valid(t *testing.T) {
	msg := &types.MsgFulfillIntent{SolverAddr: validAddr1, IntentID: "1"}
	require.NoError(t, msg.ValidateBasic())
}

func TestMsgFulfillIntent_InvalidSolver(t *testing.T) {
	msg := &types.MsgFulfillIntent{SolverAddr: "bad", IntentID: "1"}
	err := msg.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid solver address")
}

func TestMsgFulfillIntent_EmptyIntentID(t *testing.T) {
	msg := &types.MsgFulfillIntent{SolverAddr: validAddr1, IntentID: ""}
	err := msg.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "intent ID cannot be empty")
}

// ===================== Params =====================

func TestDefaultParams(t *testing.T) {
	p := types.DefaultParams()
	require.NoError(t, p.Validate())
	require.True(t, p.MinSolverStake.IsPositive())
	require.Equal(t, uint64(100), p.SolvingWindow)
	require.Equal(t, uint64(20), p.MaxSolutionsPerIntent)
	require.True(t, p.SolverSlashFraction.IsPositive())
	require.True(t, p.SolverSlashFraction.LT(math.LegacyOneDec()))
	require.Equal(t, uint64(20000), p.IntentExpiryBlocks)
	require.True(t, p.EnableIntents)
}

func TestParams_InvalidMinSolverStake(t *testing.T) {
	p := types.DefaultParams()
	p.MinSolverStake = sdk.Coin{Denom: "usyreen", Amount: math.NewInt(0)}
	require.Error(t, p.Validate())
	require.Contains(t, p.Validate().Error(), "min solver stake must be positive")
}

func TestParams_ZeroSolvingWindow(t *testing.T) {
	p := types.DefaultParams()
	p.SolvingWindow = 0
	require.Error(t, p.Validate())
	require.Contains(t, p.Validate().Error(), "solving window must be greater than 0")
}

func TestParams_ZeroMaxSolutionsPerIntent(t *testing.T) {
	p := types.DefaultParams()
	p.MaxSolutionsPerIntent = 0
	require.Error(t, p.Validate())
	require.Contains(t, p.Validate().Error(), "max solutions per intent must be greater than 0")
}

func TestParams_NegativeSolverSlashFraction(t *testing.T) {
	p := types.DefaultParams()
	p.SolverSlashFraction = math.LegacyNewDecWithPrec(-1, 1) // -0.1
	require.Error(t, p.Validate())
	require.Contains(t, p.Validate().Error(), "solver slash fraction must be between 0 and 1")
}

func TestParams_SolverSlashFractionGreaterThanOne(t *testing.T) {
	p := types.DefaultParams()
	p.SolverSlashFraction = math.LegacyNewDecWithPrec(15, 1) // 1.5
	require.Error(t, p.Validate())
	require.Contains(t, p.Validate().Error(), "solver slash fraction must be between 0 and 1")
}

func TestParams_SolverSlashFractionBoundaryZero(t *testing.T) {
	p := types.DefaultParams()
	p.SolverSlashFraction = math.LegacyZeroDec()
	require.NoError(t, p.Validate())
}

func TestParams_SolverSlashFractionBoundaryOne(t *testing.T) {
	p := types.DefaultParams()
	p.SolverSlashFraction = math.LegacyOneDec()
	require.NoError(t, p.Validate())
}

func TestParams_ZeroIntentExpiryBlocks(t *testing.T) {
	p := types.DefaultParams()
	p.IntentExpiryBlocks = 0
	require.Error(t, p.Validate())
	require.Contains(t, p.Validate().Error(), "intent expiry blocks must be greater than 0")
}

// ===================== GenesisState =====================

func TestDefaultGenesis(t *testing.T) {
	gs := types.DefaultGenesis()
	require.NoError(t, gs.Validate())
	require.Empty(t, gs.Intents)
	require.Empty(t, gs.Solvers)
}

func TestGenesisState_DuplicateIntentIDs(t *testing.T) {
	gs := types.DefaultGenesis()
	gs.Intents = []types.Intent{
		{ID: "1", Creator: validAddr1, IntentType: "swap"},
		{ID: "1", Creator: validAddr2, IntentType: "transfer"},
	}
	err := gs.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "duplicate intent ID")
}

func TestGenesisState_DuplicateSolverAddresses(t *testing.T) {
	gs := types.DefaultGenesis()
	gs.Solvers = []types.Solver{
		{Address: validAddr1, Moniker: "s1"},
		{Address: validAddr1, Moniker: "s2"},
	}
	err := gs.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "duplicate solver address")
}

func TestGenesisState_EmptyIntentID(t *testing.T) {
	gs := types.DefaultGenesis()
	gs.Intents = []types.Intent{
		{ID: "", Creator: validAddr1, IntentType: "swap"},
	}
	err := gs.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "intent has empty ID")
}

func TestGenesisState_EmptyCreator(t *testing.T) {
	gs := types.DefaultGenesis()
	gs.Intents = []types.Intent{
		{ID: "1", Creator: "", IntentType: "swap"},
	}
	err := gs.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty creator")
}

func TestGenesisState_EmptyIntentType(t *testing.T) {
	gs := types.DefaultGenesis()
	gs.Intents = []types.Intent{
		{ID: "1", Creator: validAddr1, IntentType: ""},
	}
	err := gs.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty intent type")
}

func TestGenesisState_EmptySolverAddress(t *testing.T) {
	gs := types.DefaultGenesis()
	gs.Solvers = []types.Solver{
		{Address: "", Moniker: "s1"},
	}
	err := gs.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "solver has empty address")
}

func TestGenesisState_InvalidParams(t *testing.T) {
	gs := types.DefaultGenesis()
	gs.Params.SolvingWindow = 0 // invalid
	err := gs.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "solving window must be greater than 0")
}

func TestGenesisState_ValidWithData(t *testing.T) {
	gs := types.DefaultGenesis()
	gs.Intents = []types.Intent{
		{ID: "1", Creator: validAddr1, IntentType: "swap"},
		{ID: "2", Creator: validAddr2, IntentType: "transfer"},
	}
	gs.Solvers = []types.Solver{
		{Address: validAddr1, Moniker: "solver-a"},
		{Address: validAddr2, Moniker: "solver-b"},
	}
	require.NoError(t, gs.Validate())
}

// ===================== Trading Intent Body JSON Roundtrips =====================

func TestLimitBuyIntent_JSONRoundtrip(t *testing.T) {
	original := types.LimitBuyIntent{
		InputDenom:      "usyreen",
		InputAmount:     math.NewInt(1000000),
		OutputDenom:     "uatom",
		TargetPrice:     math.LegacyNewDecWithPrec(5, 1), // 0.5
		PoolID:          42,
		MinOutputAmount: math.NewInt(900000),
	}
	bz, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded types.LimitBuyIntent
	require.NoError(t, json.Unmarshal(bz, &decoded))
	require.Equal(t, original.InputDenom, decoded.InputDenom)
	require.True(t, original.InputAmount.Equal(decoded.InputAmount))
	require.True(t, original.TargetPrice.Equal(decoded.TargetPrice))
	require.Equal(t, original.PoolID, decoded.PoolID)
	require.True(t, original.MinOutputAmount.Equal(decoded.MinOutputAmount))
}

func TestLimitSellIntent_JSONRoundtrip(t *testing.T) {
	original := types.LimitSellIntent{
		InputDenom:      "uatom",
		InputAmount:     math.NewInt(500000),
		OutputDenom:     "usyreen",
		TargetPrice:     math.LegacyNewDec(2),
		PoolID:          7,
		MinOutputAmount: math.NewInt(950000),
	}
	bz, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded types.LimitSellIntent
	require.NoError(t, json.Unmarshal(bz, &decoded))
	require.Equal(t, original.InputDenom, decoded.InputDenom)
	require.True(t, original.InputAmount.Equal(decoded.InputAmount))
	require.True(t, original.TargetPrice.Equal(decoded.TargetPrice))
	require.Equal(t, original.PoolID, decoded.PoolID)
}

func TestStopLossIntent_JSONRoundtrip(t *testing.T) {
	original := types.StopLossIntent{
		InputDenom:      "uatom",
		InputAmount:     math.NewInt(100000),
		OutputDenom:     "usyreen",
		StopPrice:       math.LegacyNewDecWithPrec(8, 1), // 0.8
		PoolID:          3,
		MinOutputAmount: math.NewInt(70000),
	}
	bz, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded types.StopLossIntent
	require.NoError(t, json.Unmarshal(bz, &decoded))
	require.Equal(t, original.InputDenom, decoded.InputDenom)
	require.True(t, original.StopPrice.Equal(decoded.StopPrice))
	require.Equal(t, original.PoolID, decoded.PoolID)
}

func TestTakeProfitIntent_JSONRoundtrip(t *testing.T) {
	original := types.TakeProfitIntent{
		InputDenom:      "uatom",
		InputAmount:     math.NewInt(200000),
		OutputDenom:     "usyreen",
		TargetPrice:     math.LegacyNewDec(3),
		PoolID:          1,
		MinOutputAmount: math.NewInt(550000),
	}
	bz, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded types.TakeProfitIntent
	require.NoError(t, json.Unmarshal(bz, &decoded))
	require.Equal(t, original.InputDenom, decoded.InputDenom)
	require.True(t, original.TargetPrice.Equal(decoded.TargetPrice))
}

func TestDCAIntent_JSONRoundtrip(t *testing.T) {
	original := types.DCAIntent{
		InputDenom:     "usyreen",
		TotalAmount:    math.NewInt(10000000),
		OutputDenom:    "uatom",
		NumExecutions:  10,
		IntervalBlocks: 100,
		PoolID:         5,
		ExecutedCount:  3,
		NextExecBlock:  500,
	}
	bz, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded types.DCAIntent
	require.NoError(t, json.Unmarshal(bz, &decoded))
	require.Equal(t, original.InputDenom, decoded.InputDenom)
	require.True(t, original.TotalAmount.Equal(decoded.TotalAmount))
	require.Equal(t, original.NumExecutions, decoded.NumExecutions)
	require.Equal(t, original.IntervalBlocks, decoded.IntervalBlocks)
	require.Equal(t, original.ExecutedCount, decoded.ExecutedCount)
	require.Equal(t, original.NextExecBlock, decoded.NextExecBlock)
}
