package types_test

import (
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	"syreen/x/intent/types"
)

// TestMsgUpdateParams_RoundTrip proves the hand-rolled wire codec for
// MsgUpdateParams preserves BOTH the authority (field 1, proto3 string) and the
// full Params (field 2, JSON bytes). This is critical: MsgUpdateParams is the
// gov param-update path, and a wrong Marshal/Unmarshal would silently wipe the
// module params when a proposal executes. The Params carries a non-default
// SolvingWindow so a dropped/zeroed field would be caught by the assertion.
func TestMsgUpdateParams_RoundTrip(t *testing.T) {
	authority := authtypes.NewModuleAddress("gov").String()

	params := types.DefaultParams()
	params.SolvingWindow = 100

	orig := &types.MsgUpdateParams{
		Authority: authority,
		Params:    params,
	}

	bz, err := orig.Marshal()
	require.NoError(t, err)
	require.NotEmpty(t, bz)

	var got types.MsgUpdateParams
	require.NoError(t, got.Unmarshal(bz))

	require.Equal(t, orig.Authority, got.Authority)
	require.Equal(t, uint64(100), got.Params.SolvingWindow)
	// The whole Params struct must survive, not just SolvingWindow.
	require.Equal(t, orig.Params.MaxSolutionsPerIntent, got.Params.MaxSolutionsPerIntent)
	require.Equal(t, orig.Params.MinSolverStake.String(), got.Params.MinSolverStake.String())
	require.Equal(t, orig.Params.AllowedMsgTypes, got.Params.AllowedMsgTypes)
}

// TestMsgUpdateParams_ValidateBasic checks authority bech32 validation and that
// invalid Params are rejected.
func TestMsgUpdateParams_ValidateBasic(t *testing.T) {
	authority := authtypes.NewModuleAddress("gov").String()

	// Valid.
	valid := &types.MsgUpdateParams{Authority: authority, Params: types.DefaultParams()}
	require.NoError(t, valid.ValidateBasic())

	// Bad authority address.
	badAuth := &types.MsgUpdateParams{Authority: "not-an-address", Params: types.DefaultParams()}
	require.Error(t, badAuth.ValidateBasic())

	// Invalid params (SolvingWindow == 0).
	badParams := types.DefaultParams()
	badParams.SolvingWindow = 0
	require.Error(t, (&types.MsgUpdateParams{Authority: authority, Params: badParams}).ValidateBasic())
}
