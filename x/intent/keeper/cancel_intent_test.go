package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"syreen/x/intent/types"
)

// TestCancelIntent_Creator proves the happy path: the creator cancels a pending
// intent, locked MaxFee+Tip are refunded to the creator, submitted solutions are
// cleaned up, and the intent moves to a terminal state.
func TestCancelIntent_Creator(t *testing.T) {
	k, ctx, bank := setupKeeperWithTrackingBank(t)

	intent := types.Intent{
		ID:         "1",
		Creator:    creatorAddr,
		IntentType: types.IntentTypeSwap,
		Body:       swapIntentBody(),
		MaxFee:     sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:        sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10)),
		Expiry:     100,
		Status:     types.StatusPending,
		CreatedAt:  1,
	}
	k.SetIntent(ctx, intent)

	// Seed a solution to verify it gets cleaned up.
	k.SetSolution(ctx, types.Solution{IntentID: "1", SolverAddr: solverAddr})
	require.Len(t, k.GetSolutionsForIntent(ctx, "1"), 1)

	err := k.CancelIntent(ctx, &types.MsgCancelIntent{Creator: creatorAddr, IntentID: "1"})
	require.NoError(t, err)

	// Status is terminal.
	got, found := k.GetIntent(ctx, "1")
	require.True(t, found)
	require.True(t, got.IsTerminal())
	require.Equal(t, types.StatusExpired, got.Status)

	// Solutions cleaned up.
	require.Empty(t, k.GetSolutionsForIntent(ctx, "1"))

	// MaxFee+Tip (110usyreen) refunded to the creator.
	creatorAcc, err := sdk.AccAddressFromBech32(creatorAddr)
	require.NoError(t, err)
	var refunded sdk.Coins
	for _, c := range bank.moduleToAccountCalls {
		if c.recipient.Equals(creatorAcc) {
			refunded = refunded.Add(c.amount...)
		}
	}
	require.True(t, refunded.AmountOf("usyreen").GTE(sdk.NewInt64Coin("usyreen", 110).Amount),
		"expected at least 110usyreen refunded, got %s", refunded)
}

// TestCancelIntent_NonCreator: only the creator may cancel.
func TestCancelIntent_NonCreator(t *testing.T) {
	k, ctx := setupKeeper(t)

	k.SetIntent(ctx, types.Intent{
		ID:         "2",
		Creator:    creatorAddr,
		IntentType: types.IntentTypeSwap,
		Body:       swapIntentBody(),
		MaxFee:     sdk.NewCoins(),
		Tip:        sdk.NewCoins(),
		Status:     types.StatusPending,
		CreatedAt:  1,
	})

	err := k.CancelIntent(ctx, &types.MsgCancelIntent{Creator: solverAddr, IntentID: "2"})
	require.ErrorIs(t, err, types.ErrIntentNotCreator)

	// Untouched.
	got, _ := k.GetIntent(ctx, "2")
	require.Equal(t, types.StatusPending, got.Status)
}

// TestCancelIntent_Fulfilled: a fulfilled (terminal) intent cannot be cancelled.
func TestCancelIntent_Fulfilled(t *testing.T) {
	k, ctx := setupKeeper(t)

	k.SetIntent(ctx, types.Intent{
		ID:         "3",
		Creator:    creatorAddr,
		IntentType: types.IntentTypeSwap,
		Body:       swapIntentBody(),
		MaxFee:     sdk.NewCoins(),
		Tip:        sdk.NewCoins(),
		Status:     types.StatusFulfilled,
		CreatedAt:  1,
	})

	err := k.CancelIntent(ctx, &types.MsgCancelIntent{Creator: creatorAddr, IntentID: "3"})
	require.ErrorIs(t, err, types.ErrIntentNotCancellable)
}

// TestCancelIntent_NotFound: cancelling a missing intent errors cleanly.
func TestCancelIntent_NotFound(t *testing.T) {
	k, ctx := setupKeeper(t)
	err := k.CancelIntent(ctx, &types.MsgCancelIntent{Creator: creatorAddr, IntentID: "nope"})
	require.ErrorIs(t, err, types.ErrIntentNotFound)
}
