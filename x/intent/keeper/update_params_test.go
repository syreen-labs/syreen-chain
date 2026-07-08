package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"syreen/x/intent/types"
)

// TestMsgServer_UpdateParams_WrongAuthority proves the gov guard: a caller whose
// authority does not match the keeper's authority is rejected and params are
// left untouched.
func TestMsgServer_UpdateParams_WrongAuthority(t *testing.T) {
	k, ctx := setupKeeper(t)
	ms := newMsgServer(k)

	before := k.GetParams(ctx)

	newParams := types.DefaultParams()
	newParams.SolvingWindow = 250

	_, err := ms.UpdateParams(ctx, &types.MsgUpdateParams{
		Authority: "not-the-authority",
		Params:    newParams,
	})
	require.Error(t, err)

	// Params unchanged.
	after := k.GetParams(ctx)
	require.Equal(t, before.SolvingWindow, after.SolvingWindow)
}

// TestMsgServer_UpdateParams_CorrectAuthority proves the happy path: the keeper's
// authority ("authority" in the test setup) can update params, and GetParams
// reflects the new SolvingWindow.
func TestMsgServer_UpdateParams_CorrectAuthority(t *testing.T) {
	k, ctx := setupKeeper(t)
	ms := newMsgServer(k)

	require.Equal(t, "authority", k.GetAuthority())

	newParams := types.DefaultParams()
	newParams.SolvingWindow = 250

	resp, err := ms.UpdateParams(ctx, &types.MsgUpdateParams{
		Authority: k.GetAuthority(),
		Params:    newParams,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	got := k.GetParams(ctx)
	require.Equal(t, uint64(250), got.SolvingWindow)
}
