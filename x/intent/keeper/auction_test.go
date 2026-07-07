package keeper_test

import (
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
	"github.com/stretchr/testify/require"

	"syreen/x/intent/keeper"
	"syreen/x/intent/types"
)

// ---------- Tests ----------

func TestSelectWinningSolver_SingleSolver(t *testing.T) {
	k, ctx := setupKeeper(t)
	intentID := registerSolverAndIntent(t, k, ctx)

	// Submit one solution
	require.NoError(t, k.SubmitSolution(ctx, &types.MsgSubmitSolution{
		SolverAddr:      solverAddr,
		IntentID:        intentID,
		ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{"@type":"/test"}`)},
		ExpectedOutcome: json.RawMessage(`{}`),
	}))

	winner, err := k.SelectWinningSolver(ctx, intentID)
	require.NoError(t, err)
	require.NotNil(t, winner)
	require.Equal(t, solverAddr, winner.SolverAddr)
}

func TestSelectWinningSolver_HigherReputationWins(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Register solver1 (default reputation 100)
	require.NoError(t, k.RegisterSolver(ctx, &types.MsgRegisterSolver{
		Address:     solverAddr,
		Moniker:     "solver-low-rep",
		StakeAmount: sdk.NewCoin("usyreen", math.NewInt(2000000000)),
	}))

	// Register solver2 (default reputation 100, we'll bump it)
	require.NoError(t, k.RegisterSolver(ctx, &types.MsgRegisterSolver{
		Address:     solverAddr2,
		Moniker:     "solver-high-rep",
		StakeAmount: sdk.NewCoin("usyreen", math.NewInt(2000000000)),
	}))

	// Boost solver2's reputation to 500
	s2, found := k.GetSolver(ctx, solverAddr2)
	require.True(t, found)
	s2.ReputationScore = 500
	k.SetSolver(ctx, s2)

	// Submit intent from creator
	intentID, err := k.SubmitIntent(ctx, &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{"test":1}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:          sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10)),
		ExpiryBlocks: 50,
	})
	require.NoError(t, err)

	// Both solvers submit solutions with equal gas estimates
	require.NoError(t, k.SubmitSolution(ctx, &types.MsgSubmitSolution{
		SolverAddr:      solverAddr,
		IntentID:        intentID,
		ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{"@type":"/test"}`)},
		ExpectedOutcome: json.RawMessage(`{}`),
	}))
	require.NoError(t, k.SubmitSolution(ctx, &types.MsgSubmitSolution{
		SolverAddr:      solverAddr2,
		IntentID:        intentID,
		ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{"@type":"/test"}`)},
		ExpectedOutcome: json.RawMessage(`{}`),
	}))

	winner, err := k.SelectWinningSolver(ctx, intentID)
	require.NoError(t, err)
	require.NotNil(t, winner)
	// solver2 has 500 reputation vs 100 => solver2 should win due to 40% weight on reputation
	require.Equal(t, solverAddr2, winner.SolverAddr)
}

func TestSelectWinningSolver_NoSolutions(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Create intent without any solutions
	intentID, err := k.SubmitIntent(ctx, &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1)),
		Tip:          sdk.Coins{},
		ExpiryBlocks: 50,
	})
	require.NoError(t, err)

	_, err = k.SelectWinningSolver(ctx, intentID)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrNoSolutions)
}

func TestFulfillIntent_AutoSelect_EmptySolverAddr(t *testing.T) {
	k, ctx := setupKeeper(t)
	intentID := registerSolverAndIntent(t, k, ctx)

	// Submit solution from solver
	require.NoError(t, k.SubmitSolution(ctx, &types.MsgSubmitSolution{
		SolverAddr:      solverAddr,
		IntentID:        intentID,
		ExecutionMsgs:   []json.RawMessage{bankSendExecMsg()},
		ExpectedOutcome: json.RawMessage(`{}`),
	}))

	// Fulfill with empty SolverAddr => auto-select
	err := k.FulfillIntent(ctx, &types.MsgFulfillIntent{
		SolverAddr: "", // triggers auction
		IntentID:   intentID,
	})
	require.NoError(t, err)

	// Verify intent is fulfilled with the right solver
	intent, found := k.GetIntent(ctx, intentID)
	require.True(t, found)
	require.Equal(t, types.StatusFulfilled, intent.Status)
	require.Equal(t, solverAddr, intent.SolverAddr)
}

func TestSlashing_ReducesStakeBy10Percent(t *testing.T) {
	kFail, ctxFail := setupKeeperWithFailingRouter(t)

	stakeAmt := math.NewInt(2000000000) // 2000 SYR
	solver := types.Solver{
		Address:         solverAddr,
		Moniker:         "solver-to-slash",
		StakedAmount:    sdk.NewCoin("usyreen", stakeAmt),
		ReputationScore: 100,
		Active:          true,
		JoinedAt:        1,
	}
	kFail.SetSolver(ctxFail, solver)

	intentID := "1"
	intent := types.Intent{
		ID:         intentID,
		Creator:    creatorAddr,
		IntentType: types.IntentTypeSwap,
		Body:       json.RawMessage(`{"test":1}`),
		MaxFee:     sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:        sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10)),
		Expiry:     51,
		Status:     types.StatusSolving,
		CreatedAt:  1,
	}
	kFail.SetIntent(ctxFail, intent)

	sol := types.Solution{
		IntentID:        intentID,
		SolverAddr:      solverAddr,
		ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{"@type":"/test"}`)},
		ExpectedOutcome: json.RawMessage(`{}`),
		SubmittedAt:     1,
	}
	kFail.SetSolution(ctxFail, sol)

	// Attempt to fulfill -- should fail because of the failing router
	err := kFail.FulfillIntent(ctxFail, &types.MsgFulfillIntent{
		SolverAddr: solverAddr,
		IntentID:   intentID,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "solution execution failed")

	// Check that solver was slashed
	slashedSolver, found := kFail.GetSolver(ctxFail, solverAddr)
	require.True(t, found)

	// 10% of 2000000000 = 200000000
	expectedRemaining := stakeAmt.Sub(math.NewInt(200000000))
	require.Equal(t, expectedRemaining, slashedSolver.StakedAmount.Amount,
		"stake should be reduced by 10%%: expected %s, got %s", expectedRemaining, slashedSolver.StakedAmount.Amount)

	// Reputation should have decreased by 5
	require.Equal(t, uint64(95), slashedSolver.ReputationScore)
	require.Equal(t, uint64(1), slashedSolver.TotalFailed)

	// Solver should still be active (1800000000 > MinSolverStake of 1000000000)
	require.True(t, slashedSolver.Active)
}

func TestSlashing_DeactivatesWhenBelowMinStake(t *testing.T) {
	kFail, ctxFail := setupKeeperWithFailingRouter(t)

	// Set solver with stake just above min (MinSolverStake = 1000000000)
	// 10% slash of 1050000000 = 105000000 => remaining = 945000000 < 1000000000 => deactivated
	stakeAmt := math.NewInt(1050000000)
	solver := types.Solver{
		Address:         solverAddr,
		Moniker:         "barely-staked",
		StakedAmount:    sdk.NewCoin("usyreen", stakeAmt),
		ReputationScore: 100,
		Active:          true,
		JoinedAt:        1,
	}
	kFail.SetSolver(ctxFail, solver)

	intentID := "deactivate-1"
	intent := types.Intent{
		ID:         intentID,
		Creator:    creatorAddr,
		IntentType: types.IntentTypeSwap,
		Body:       json.RawMessage(`{}`),
		MaxFee:     sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:        sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10)),
		Expiry:     51,
		Status:     types.StatusSolving,
		CreatedAt:  1,
	}
	kFail.SetIntent(ctxFail, intent)

	sol := types.Solution{
		IntentID:        intentID,
		SolverAddr:      solverAddr,
		ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{"@type":"/test"}`)},
		ExpectedOutcome: json.RawMessage(`{}`),
		SubmittedAt:     1,
	}
	kFail.SetSolution(ctxFail, sol)

	err := kFail.FulfillIntent(ctxFail, &types.MsgFulfillIntent{
		SolverAddr: solverAddr,
		IntentID:   intentID,
	})
	require.Error(t, err)

	slashedSolver, found := kFail.GetSolver(ctxFail, solverAddr)
	require.True(t, found)
	require.False(t, slashedSolver.Active, "solver should be deactivated when stake drops below minimum")
	expectedRemaining := stakeAmt.Sub(math.NewInt(105000000))
	require.Equal(t, expectedRemaining, slashedSolver.StakedAmount.Amount)
}

func TestAutoFulfillIntents_FulfillsAfterSolvingWindow(t *testing.T) {
	k, ctx := setupKeeper(t)
	intentID := registerSolverAndIntent(t, k, ctx)

	// Submit solution
	require.NoError(t, k.SubmitSolution(ctx, &types.MsgSubmitSolution{
		SolverAddr:      solverAddr,
		IntentID:        intentID,
		ExecutionMsgs:   []json.RawMessage{bankSendExecMsg()},
		ExpectedOutcome: json.RawMessage(`{}`),
	}))

	// Verify intent is in "solving" state
	intent, found := k.GetIntent(ctx, intentID)
	require.True(t, found)
	require.Equal(t, types.StatusSolving, intent.Status)

	// Advance past solving window (CreatedAt=1, SolvingWindow=10 => deadline = 11)
	ctx = ctx.WithBlockHeight(12)

	// Run auto-fulfill
	k.AutoFulfillIntents(ctx)

	// Intent should now be fulfilled
	intent, found = k.GetIntent(ctx, intentID)
	require.True(t, found)
	require.Equal(t, types.StatusFulfilled, intent.Status)
	require.Equal(t, solverAddr, intent.SolverAddr)
}

func TestAutoFulfillIntents_ExpiresWhenNoSolutions(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Submit intent
	intentID, err := k.SubmitIntent(ctx, &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:          sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10)),
		ExpiryBlocks: 50,
	})
	require.NoError(t, err)

	// Manually set to solving
	intent, found := k.GetIntent(ctx, intentID)
	require.True(t, found)
	intent.Status = types.StatusSolving
	k.SetIntent(ctx, intent)

	// Advance past solving window
	ctx = ctx.WithBlockHeight(12)
	k.AutoFulfillIntents(ctx)

	// Intent should be expired (no solutions to fulfill)
	intent, found = k.GetIntent(ctx, intentID)
	require.True(t, found)
	require.Equal(t, types.StatusExpired, intent.Status)
}

func TestAutoFulfillIntents_SkipsWithinWindow(t *testing.T) {
	k, ctx := setupKeeper(t)
	intentID := registerSolverAndIntent(t, k, ctx)

	require.NoError(t, k.SubmitSolution(ctx, &types.MsgSubmitSolution{
		SolverAddr:      solverAddr,
		IntentID:        intentID,
		ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{"@type":"/test"}`)},
		ExpectedOutcome: json.RawMessage(`{}`),
	}))

	// Still within solving window (block 5, deadline = 1+10 = 11)
	ctx = ctx.WithBlockHeight(5)
	k.AutoFulfillIntents(ctx)

	// Intent should still be solving, not fulfilled
	intent, found := k.GetIntent(ctx, intentID)
	require.True(t, found)
	require.Equal(t, types.StatusSolving, intent.Status)
}

// ---------- setupKeeperWithFailingRouter ----------

func setupKeeperWithFailingRouter(t *testing.T) (*keeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey("intent")
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, cmtproto.Header{Height: 1}, false, log.NewNopLogger())
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	storeService := runtime.NewKVStoreService(storeKey)

	router := &mockMsgRouter{
		handler: func(_ sdk.Context, _ sdk.Msg) (*sdk.Result, error) {
			return nil, fmt.Errorf("execution failed")
		},
	}

	k := keeper.NewKeeper(cdc, storeService, &mockAccountKeeper{}, &mockBankKeeper{}, router, "authority")
	params := types.DefaultParams()
	params.AllowedMsgTypes = append(params.AllowedMsgTypes, "/")
	require.NoError(t, k.SetParams(ctx, params))

	return k, ctx
}
