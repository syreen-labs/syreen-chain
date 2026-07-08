package keeper_test

import (
	"encoding/json"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"syreen/x/intent/types"
)

// TestExportImportGenesisRoundTrip proves that ExportGenesis followed by
// InitGenesis into a fresh keeper is lossless: intents, solvers, solutions,
// chains, strategies AND the three monotonic ID counters all survive. This is
// the state-migration path; before the fix it silently dropped solutions/chains/
// strategies and reset the counters (risking duplicate IDs).
func TestExportImportGenesisRoundTrip(t *testing.T) {
	k1, ctx1 := setupKeeper(t)

	// --- Seed a rich working set ---
	intent := types.Intent{
		ID:         "7",
		Creator:    creatorAddr,
		IntentType: types.IntentTypeSwap,
		Body:       swapIntentBody(),
		MaxFee:     sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		Tip:        sdk.NewCoins(sdk.NewInt64Coin("usyreen", 10)),
		Expiry:     50,
		Status:     types.StatusSolving,
		CreatedAt:  1,
	}
	k1.SetIntent(ctx1, intent)

	solver := types.Solver{
		Address:      solverAddr,
		Moniker:      "sv",
		StakedAmount: sdk.NewInt64Coin("usyreen", 5000),
		Active:       true,
		JoinedAt:     1,
	}
	k1.SetSolver(ctx1, solver)

	solution := types.Solution{
		IntentID:      "7",
		SolverAddr:    solverAddr,
		ExecutionMsgs: []json.RawMessage{bankSendExecMsg()},
		GasEstimate:   1234,
		SubmittedAt:   2,
	}
	k1.SetSolution(ctx1, solution)

	chain := types.IntentChain{
		ID:         "chain-3",
		Creator:    creatorAddr,
		Status:     types.ChainStatusActive,
		CurrentIdx: 1,
		CreatedAt:  1,
		Expiry:     500,
		MaxFee:     sdk.NewCoins(),
		Tip:        sdk.NewCoins(),
	}
	k1.SetChain(ctx1, chain)

	strategy := types.Strategy{
		ID:           "strategy-2",
		Creator:      creatorAddr,
		TemplateName: types.StrategySafeAccumulate,
		ChainID:      "chain-3",
		Status:       types.StrategyStatusActive,
		TotalBudget:  math.NewInt(1000),
		RiskLevel:    types.RiskModerate,
		CreatedAt:    1,
	}
	k1.SetStrategy(ctx1, strategy)

	// Advance the counters to non-trivial values by writing them directly via
	// the exported genesis restore path is not available, so seed through a
	// full export/import: first stamp the counters using InitGenesis below.
	gsSeed := k1.ExportGenesis(ctx1)
	gsSeed.NextIntentId = 7
	gsSeed.NextChainId = 3
	gsSeed.NextStrategyId = 2

	// --- Round-trip into a fresh keeper ---
	k2, ctx2 := setupKeeper(t)
	k2.InitGenesis(ctx2, *gsSeed)
	gsOut := k2.ExportGenesis(ctx2)

	// Counters preserved
	require.Equal(t, uint64(7), gsOut.NextIntentId, "intent counter lost")
	require.Equal(t, uint64(3), gsOut.NextChainId, "chain counter lost")
	require.Equal(t, uint64(2), gsOut.NextStrategyId, "strategy counter lost")

	// Intent preserved
	gotIntent, found := k2.GetIntent(ctx2, "7")
	require.True(t, found)
	require.Equal(t, types.StatusSolving, gotIntent.Status)

	// Solver preserved
	gotSolver, found := k2.GetSolver(ctx2, solverAddr)
	require.True(t, found)
	require.Equal(t, "sv", gotSolver.Moniker)

	// Solution preserved (was silently dropped before the fix)
	sols := k2.GetSolutionsForIntent(ctx2, "7")
	require.Len(t, sols, 1)
	require.Equal(t, uint64(1234), sols[0].GasEstimate)

	// Chain preserved
	gotChain, found := k2.GetChain(ctx2, "chain-3")
	require.True(t, found)
	require.Equal(t, types.ChainStatusActive, gotChain.Status)

	// Strategy preserved
	gotStrategy, found := k2.GetStrategy(ctx2, "strategy-2")
	require.True(t, found)
	require.Equal(t, "chain-3", gotStrategy.ChainID)

	// Whole-genesis equivalence for the slices/counters
	require.Len(t, gsOut.Intents, 1)
	require.Len(t, gsOut.Solvers, 1)
	require.Len(t, gsOut.Solutions, 1)
	require.Len(t, gsOut.Chains, 1)
	require.Len(t, gsOut.Strategies, 1)
}
