package keeper_test

import (
	"encoding/json"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"syreen/x/intent/types"
)

// ===================== C-01: Reentrancy Guard =====================

func TestFulfillIntent_ReentrantCallRejected(t *testing.T) {
	// Set up a keeper with a router that attempts to re-enter FulfillIntent
	// when handling a solution message.
	k, ctx := setupKeeper(t)

	// Register solver
	require.NoError(t, k.RegisterSolver(ctx, &types.MsgRegisterSolver{
		Address:     solverAddr,
		Moniker:     "solver",
		StakeAmount: sdk.NewCoin("usyreen", math.NewInt(2000000000)),
	}))

	// Submit intent
	intentID, err := k.SubmitIntent(ctx, &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{"test":1}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:          sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10)),
		ExpiryBlocks: 50,
	})
	require.NoError(t, err)

	// Submit solution
	require.NoError(t, k.SubmitSolution(ctx, &types.MsgSubmitSolution{
		SolverAddr:      solverAddr,
		IntentID:        intentID,
		ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{"@type":"/test"}`)},
		ExpectedOutcome: json.RawMessage(`{}`),
	}))

	// First call should succeed
	err = k.FulfillIntent(ctx, &types.MsgFulfillIntent{
		SolverAddr: solverAddr,
		IntentID:   intentID,
	})
	require.NoError(t, err)

	// Verify intent is fulfilled
	intent, found := k.GetIntent(ctx, intentID)
	require.True(t, found)
	require.Equal(t, types.StatusFulfilled, intent.Status)

	// Now test that attempting to call FulfillIntent while already executing
	// is properly rejected. We simulate this by manually setting the executing flag.
	// In production, this would happen when a solution's execution message
	// includes a MsgFulfillIntent routed back through the msg router.

	// Create a second intent to attempt re-entry
	intentID2, err := k.SubmitIntent(ctx, &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{"test":2}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:          sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10)),
		ExpiryBlocks: 50,
	})
	require.NoError(t, err)

	require.NoError(t, k.SubmitSolution(ctx, &types.MsgSubmitSolution{
		SolverAddr:      solverAddr,
		IntentID:        intentID2,
		ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{"@type":"/test"}`)},
		ExpectedOutcome: json.RawMessage(`{}`),
	}))

	// Simulate reentrancy by setting executing flag via context before calling
	reentrantCtx := k.SetExecuting(ctx, true)

	err = k.FulfillIntent(reentrantCtx, &types.MsgFulfillIntent{
		SolverAddr: solverAddr,
		IntentID:   intentID2,
	})
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrReentrant)

	// Verify the second intent was NOT fulfilled
	intent2, found := k.GetIntent(ctx, intentID2)
	require.True(t, found)
	require.NotEqual(t, types.StatusFulfilled, intent2.Status)
}

// ===================== C-11: Message Type Whitelist =====================

func TestExecuteSolution_DisallowedMsgTypeRejected(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Register solver and submit intent
	require.NoError(t, k.RegisterSolver(ctx, &types.MsgRegisterSolver{
		Address:     solverAddr,
		Moniker:     "solver",
		StakeAmount: sdk.NewCoin("usyreen", math.NewInt(2000000000)),
	}))

	intentID, err := k.SubmitIntent(ctx, &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{"test":1}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:          sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10)),
		ExpiryBlocks: 50,
	})
	require.NoError(t, err)

	// Set params with restricted whitelist (remove "/" which is the test mock type)
	params := k.GetParams(ctx)
	params.AllowedMsgTypes = []string{
		"/cosmos.bank.v1beta1.MsgSend",
		"/cosmos.bank.v1beta1.MsgMultiSend",
	}
	require.NoError(t, k.SetParams(ctx, params))

	// Submit solution with a message that resolves to type "/" (the mock type)
	require.NoError(t, k.SubmitSolution(ctx, &types.MsgSubmitSolution{
		SolverAddr:      solverAddr,
		IntentID:        intentID,
		ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{"@type":"/test"}`)},
		ExpectedOutcome: json.RawMessage(`{}`),
	}))

	// Try to fulfill — should fail because "/" is not in the whitelist
	err = k.FulfillIntent(ctx, &types.MsgFulfillIntent{
		SolverAddr: solverAddr,
		IntentID:   intentID,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "not allowed in solution execution")
}

func TestExecuteSolution_AllowedMsgTypeSucceeds(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Register solver and submit intent
	require.NoError(t, k.RegisterSolver(ctx, &types.MsgRegisterSolver{
		Address:     solverAddr,
		Moniker:     "solver",
		StakeAmount: sdk.NewCoin("usyreen", math.NewInt(2000000000)),
	}))

	intentID, err := k.SubmitIntent(ctx, &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{"test":1}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:          sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10)),
		ExpiryBlocks: 50,
	})
	require.NoError(t, err)

	// Params already have "/" in AllowedMsgTypes from setupKeeper
	require.NoError(t, k.SubmitSolution(ctx, &types.MsgSubmitSolution{
		SolverAddr:      solverAddr,
		IntentID:        intentID,
		ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{"@type":"/test"}`)},
		ExpectedOutcome: json.RawMessage(`{}`),
	}))

	// Should succeed because "/" is in the whitelist
	err = k.FulfillIntent(ctx, &types.MsgFulfillIntent{
		SolverAddr: solverAddr,
		IntentID:   intentID,
	})
	require.NoError(t, err)

	intent, found := k.GetIntent(ctx, intentID)
	require.True(t, found)
	require.Equal(t, types.StatusFulfilled, intent.Status)
}

func TestExecuteSolution_WhitelistUpdatableViaParams(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Verify default params include the standard whitelist entries
	params := k.GetParams(ctx)
	require.Contains(t, params.AllowedMsgTypes, "/cosmos.bank.v1beta1.MsgSend")
	require.Contains(t, params.AllowedMsgTypes, "/cosmos.bank.v1beta1.MsgMultiSend")
	require.Contains(t, params.AllowedMsgTypes, "/ibc.applications.transfer.v1.MsgTransfer")
	require.Contains(t, params.AllowedMsgTypes, "/syreen.tokenfactory.MsgMint")
	require.Contains(t, params.AllowedMsgTypes, "/syreen.tokenfactory.MsgBurn")
	require.Contains(t, params.AllowedMsgTypes, "/syreen.compute.MsgExecuteContract")

	// Update whitelist via params (simulating governance)
	params.AllowedMsgTypes = []string{"/cosmos.bank.v1beta1.MsgSend"}
	require.NoError(t, k.SetParams(ctx, params))

	updated := k.GetParams(ctx)
	require.Len(t, updated.AllowedMsgTypes, 1)
	require.Equal(t, "/cosmos.bank.v1beta1.MsgSend", updated.AllowedMsgTypes[0])
}

// ===================== C-09: Refund on Failed Solution =====================

func TestFulfillIntent_FailedSolutionRefundsCreator(t *testing.T) {
	// Use the failing router so solution execution fails
	kFail, ctxFail := setupKeeperWithFailingRouter(t)

	stakeAmt := math.NewInt(2000000000)
	solver := types.Solver{
		Address:         solverAddr,
		Moniker:         "solver-to-slash",
		StakedAmount:    sdk.NewCoin("usyreen", stakeAmt),
		ReputationScore: 100,
		Active:          true,
		JoinedAt:        1,
	}
	kFail.SetSolver(ctxFail, solver)

	intentID := "refund-test-1"
	maxFee := sdk.NewCoins(sdk.NewInt64Coin("usyreen", 500))
	tip := sdk.NewCoins(sdk.NewInt64Coin("usyreen", 50))

	intent := types.Intent{
		ID:         intentID,
		Creator:    creatorAddr,
		IntentType: types.IntentTypeSwap,
		Body:       json.RawMessage(`{"test":1}`),
		MaxFee:     maxFee,
		Tip:        tip,
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

	// Track bank calls to verify refund happens
	// The mock bank keeper will succeed (no error), so the refund call will be made
	err := kFail.FulfillIntent(ctxFail, &types.MsgFulfillIntent{
		SolverAddr: solverAddr,
		IntentID:   intentID,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "solution execution failed")

	// Verify intent is marked as failed
	failedIntent, found := kFail.GetIntent(ctxFail, intentID)
	require.True(t, found)
	require.Equal(t, types.StatusFailed, failedIntent.Status)

	// The refund was attempted (mock bank keeper succeeds, so no error logged).
	// In a real scenario, the bank keeper would transfer MaxFee + Tip back to creator.
	// We verify this by using a tracking bank keeper.
}

func TestFulfillIntent_FailedSolutionRefundsCreator_WithTracking(t *testing.T) {
	// Use a custom setup with a tracking bank keeper to verify refund calls
	kFail, ctxFail, bankKeeper := setupKeeperWithTrackingBank(t)

	stakeAmt := math.NewInt(2000000000)
	solver := types.Solver{
		Address:         solverAddr,
		Moniker:         "solver-refund",
		StakedAmount:    sdk.NewCoin("usyreen", stakeAmt),
		ReputationScore: 100,
		Active:          true,
		JoinedAt:        1,
	}
	kFail.SetSolver(ctxFail, solver)

	intentID := "refund-track-1"
	maxFee := sdk.NewCoins(sdk.NewInt64Coin("usyreen", 500))
	tip := sdk.NewCoins(sdk.NewInt64Coin("usyreen", 50))

	intent := types.Intent{
		ID:         intentID,
		Creator:    creatorAddr,
		IntentType: types.IntentTypeSwap,
		Body:       json.RawMessage(`{"test":1}`),
		MaxFee:     maxFee,
		Tip:        tip,
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
	require.Contains(t, err.Error(), "solution execution failed")

	// Verify the bank keeper received a refund call (SendCoinsFromModuleToAccount)
	// The refund should be MaxFee + Tip = 550usyreen
	creatorAccAddr, _ := sdk.AccAddressFromBech32(creatorAddr)
	expectedRefund := maxFee.Add(tip...)

	found := false
	for _, call := range bankKeeper.moduleToAccountCalls {
		if call.module == types.ModuleName &&
			call.recipient.Equals(creatorAccAddr) &&
			call.amount.Equal(expectedRefund) {
			found = true
			break
		}
	}
	require.True(t, found, "expected refund of %s to creator %s but got calls: %v",
		expectedRefund, creatorAddr, bankKeeper.moduleToAccountCalls)
}
