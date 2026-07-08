package keeper

import (
	"bytes"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/intent/types"
)

// InitGenesis initializes the module's state from a provided genesis state
func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	// Force-create the module account so its address resolves as a real
	// ModuleAccount from block 0. The intent account is intentionally left out of
	// the bank blocklist (it must receive DEX swap outputs), so without this a user
	// could bank-send to the address before the module first uses it, creating a
	// BaseAccount and corrupting the account (later escrow panics with
	// "account is not a module account").
	k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)

	k.SetParams(ctx, data.Params)

	for _, intent := range data.Intents {
		k.SetIntent(ctx, intent)
	}
	for _, solver := range data.Solvers {
		k.SetSolver(ctx, solver)
	}
	for _, solution := range data.Solutions {
		k.SetSolution(ctx, solution)
	}
	for _, chain := range data.Chains {
		k.SetChain(ctx, chain)
	}
	for _, strategy := range data.Strategies {
		k.SetStrategy(ctx, strategy)
	}

	// Restore the ID counters so post-migration ID allocation continues from
	// where the exported chain left off (rather than colliding from 0).
	k.setIntentCounter(ctx, data.NextIntentId)
	k.setChainCounter(ctx, data.NextChainId)
	k.setStrategyCounter(ctx, data.NextStrategyId)
}

// ExportGenesis returns the module's exported genesis
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params:         k.GetParams(ctx),
		Intents:        k.getAllIntents(ctx),
		Solvers:        k.getAllSolvers(ctx),
		Solutions:      k.getAllSolutions(ctx),
		Chains:         k.getAllChains(ctx),
		Strategies:     k.getAllStrategies(ctx),
		NextIntentId:   k.getIntentCounter(ctx),
		NextChainId:    k.getChainCounter(ctx),
		NextStrategyId: k.getStrategyCounter(ctx),
	}
}

// getAllIntents returns all intents from the store
func (k Keeper) getAllIntents(ctx sdk.Context) []types.Intent {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.IntentPrefix)

	iter, err := kvStore.Iterator(prefix, nil)
	if err != nil {
		return nil
	}
	defer iter.Close()

	var intents []types.Intent
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if !bytes.HasPrefix(key, prefix) {
			break
		}
		var intent types.Intent
		if err := json.Unmarshal(iter.Value(), &intent); err != nil {
			continue
		}
		intents = append(intents, intent)
	}

	return intents
}

// getAllSolvers returns all solvers from the store
func (k Keeper) getAllSolvers(ctx sdk.Context) []types.Solver {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.SolverPrefix)

	iter, err := kvStore.Iterator(prefix, nil)
	if err != nil {
		return nil
	}
	defer iter.Close()

	var solvers []types.Solver
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if !bytes.HasPrefix(key, prefix) {
			break
		}
		var solver types.Solver
		if err := json.Unmarshal(iter.Value(), &solver); err != nil {
			continue
		}
		solvers = append(solvers, solver)
	}

	return solvers
}

// getAllSolutions returns all solutions from the store
func (k Keeper) getAllSolutions(ctx sdk.Context) []types.Solution {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.SolutionPrefix)

	iter, err := kvStore.Iterator(prefix, nil)
	if err != nil {
		return nil
	}
	defer iter.Close()

	var solutions []types.Solution
	for ; iter.Valid(); iter.Next() {
		if !bytes.HasPrefix(iter.Key(), prefix) {
			break
		}
		var solution types.Solution
		if err := json.Unmarshal(iter.Value(), &solution); err != nil {
			continue
		}
		solutions = append(solutions, solution)
	}

	return solutions
}

// getAllChains returns all intent chains from the store
func (k Keeper) getAllChains(ctx sdk.Context) []types.IntentChain {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.ChainPrefix)

	iter, err := kvStore.Iterator(prefix, nil)
	if err != nil {
		return nil
	}
	defer iter.Close()

	var chains []types.IntentChain
	for ; iter.Valid(); iter.Next() {
		if !bytes.HasPrefix(iter.Key(), prefix) {
			break
		}
		var chain types.IntentChain
		if err := json.Unmarshal(iter.Value(), &chain); err != nil {
			continue
		}
		chains = append(chains, chain)
	}

	return chains
}

// getAllStrategies returns all strategies from the store
func (k Keeper) getAllStrategies(ctx sdk.Context) []types.Strategy {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.StrategyPrefix)

	iter, err := kvStore.Iterator(prefix, nil)
	if err != nil {
		return nil
	}
	defer iter.Close()

	var strategies []types.Strategy
	for ; iter.Valid(); iter.Next() {
		if !bytes.HasPrefix(iter.Key(), prefix) {
			break
		}
		var strategy types.Strategy
		if err := json.Unmarshal(iter.Value(), &strategy); err != nil {
			continue
		}
		strategies = append(strategies, strategy)
	}

	return strategies
}

// ---------------------------------------------------------------------------
// Raw ID-counter accessors (peek/set without incrementing)
//
// These read/write the same keys used by nextIntentID / nextChainID /
// nextStrategyID (see keeper.go, chain.go, strategy.go). ExportGenesis snapshots
// the current counter value; InitGenesis restores it so ID allocation resumes.
// ---------------------------------------------------------------------------

func (k Keeper) getIntentCounter(ctx sdk.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.IntentCounterKey())
	if err != nil || bz == nil {
		return 0
	}
	return types.BytesToUint64(bz)
}

func (k Keeper) setIntentCounter(ctx sdk.Context, v uint64) {
	if v == 0 {
		return
	}
	kvStore := k.storeService.OpenKVStore(ctx)
	kvStore.Set(types.IntentCounterKey(), types.Uint64ToBytes(v))
}

func (k Keeper) getChainCounter(ctx sdk.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.ChainCounterKey())
	if err != nil || bz == nil {
		return 0
	}
	return types.BytesToUint64(bz)
}

func (k Keeper) setChainCounter(ctx sdk.Context, v uint64) {
	if v == 0 {
		return
	}
	kvStore := k.storeService.OpenKVStore(ctx)
	kvStore.Set(types.ChainCounterKey(), types.Uint64ToBytes(v))
}

func (k Keeper) getStrategyCounter(ctx sdk.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.StrategyCounterKey())
	if err != nil || bz == nil {
		return 0
	}
	return types.BytesToUint64(bz)
}

func (k Keeper) setStrategyCounter(ctx sdk.Context, v uint64) {
	if v == 0 {
		return
	}
	kvStore := k.storeService.OpenKVStore(ctx)
	kvStore.Set(types.StrategyCounterKey(), types.Uint64ToBytes(v))
}
