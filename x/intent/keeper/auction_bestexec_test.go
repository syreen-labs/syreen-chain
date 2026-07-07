package keeper_test

import (
	"encoding/json"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"syreen/x/intent/types"
)

// Fairness Engine · Ticket 1 (best execution).
// The winner of a swap auction must be the solver who delivers the MOST output
// to the user — not the one with the best reputation/stake. This is the keystone
// change: it flips the auction from "reward the solver" to "reward the user".
func TestSelectWinningSolver_BestExecution_HighestOutputWins(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Two solvers, equal stake. solverAddr gets a HUGE reputation edge so that
	// under the OLD reputation-weighted auction it would win outright.
	require.NoError(t, k.RegisterSolver(ctx, &types.MsgRegisterSolver{
		Address: solverAddr, Moniker: "high-rep-low-price",
		StakeAmount: sdk.NewCoin("usyreen", math.NewInt(2000000000)),
	}))
	require.NoError(t, k.RegisterSolver(ctx, &types.MsgRegisterSolver{
		Address: solverAddr2, Moniker: "low-rep-best-price",
		StakeAmount: sdk.NewCoin("usyreen", math.NewInt(2000000000)),
	}))
	s1, found := k.GetSolver(ctx, solverAddr)
	require.True(t, found)
	s1.ReputationScore = 1000 // would dominate the old 40%-reputation score
	k.SetSolver(ctx, s1)

	// A real swap intent: user wants uusdc, floor of 100.
	intentID, err := k.SubmitIntent(ctx, &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{"output_denom":"uusdc","min_output_amount":"100"}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:          sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10)),
		ExpiryBlocks: 50,
	})
	require.NoError(t, err)

	// High-reputation solver promises only 150 to the user...
	require.NoError(t, k.SubmitSolution(ctx, &types.MsgSubmitSolution{
		SolverAddr:      solverAddr,
		IntentID:        intentID,
		ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{"@type":"/test"}`)},
		ExpectedOutcome: json.RawMessage(`{"output_amount":"150"}`),
	}))
	// ...low-reputation solver promises 200 — a better deal for the user.
	require.NoError(t, k.SubmitSolution(ctx, &types.MsgSubmitSolution{
		SolverAddr:      solverAddr2,
		IntentID:        intentID,
		ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{"@type":"/test"}`)},
		ExpectedOutcome: json.RawMessage(`{"output_amount":"200"}`),
	}))

	winner, err := k.SelectWinningSolver(ctx, intentID)
	require.NoError(t, err)
	require.NotNil(t, winner)
	require.Equal(t, solverAddr2, winner.SolverAddr,
		"best execution: the solver delivering the most output to the user must win, despite lower reputation")
}

// Equal declared output falls back to deterministic tiebreakers (reputation).
func TestSelectWinningSolver_BestExecution_TieBreaksByReputation(t *testing.T) {
	k, ctx := setupKeeper(t)
	require.NoError(t, k.RegisterSolver(ctx, &types.MsgRegisterSolver{
		Address: solverAddr, Moniker: "a", StakeAmount: sdk.NewCoin("usyreen", math.NewInt(2000000000)),
	}))
	require.NoError(t, k.RegisterSolver(ctx, &types.MsgRegisterSolver{
		Address: solverAddr2, Moniker: "b", StakeAmount: sdk.NewCoin("usyreen", math.NewInt(2000000000)),
	}))
	s2, _ := k.GetSolver(ctx, solverAddr2)
	s2.ReputationScore = 500
	k.SetSolver(ctx, s2)

	intentID, err := k.SubmitIntent(ctx, &types.MsgSubmitIntent{
		Creator:      creatorAddr,
		IntentType:   types.IntentTypeSwap,
		Body:         json.RawMessage(`{"output_denom":"uusdc","min_output_amount":"100"}`),
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		ExpiryBlocks: 50,
	})
	require.NoError(t, err)
	// Both promise the SAME 180 → tiebreak should pick the higher-reputation solver.
	for _, addr := range []string{solverAddr, solverAddr2} {
		require.NoError(t, k.SubmitSolution(ctx, &types.MsgSubmitSolution{
			SolverAddr:      addr,
			IntentID:        intentID,
			ExecutionMsgs:   []json.RawMessage{json.RawMessage(`{"@type":"/test"}`)},
			ExpectedOutcome: json.RawMessage(`{"output_amount":"180"}`),
		}))
	}
	winner, err := k.SelectWinningSolver(ctx, intentID)
	require.NoError(t, err)
	require.Equal(t, solverAddr2, winner.SolverAddr, "equal output ties break to higher reputation")
}

// The declared-output parser: valid/positive amounts parse; missing, empty,
// zero, and malformed declarations are treated as "no declaration".
func TestSolution_DeclaredOutput(t *testing.T) {
	cases := []struct {
		raw  string
		want string
		ok   bool
	}{
		{`{"output_amount":"200"}`, "200", true},
		{`{"output_amount":"1"}`, "1", true},
		{`{}`, "", false},
		{``, "", false},
		{`{"output_amount":"0"}`, "", false},  // zero is not a positive offer
		{`{"output_amount":"-5"}`, "", false}, // negative rejected
		{`not-json`, "", false},
	}
	for _, c := range cases {
		s := types.Solution{ExpectedOutcome: json.RawMessage(c.raw)}
		got, ok := s.DeclaredOutput()
		require.Equal(t, c.ok, ok, "raw=%q", c.raw)
		if c.ok {
			require.Equal(t, c.want, got.String(), "raw=%q", c.raw)
		}
	}
}
