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
}

// ExportGenesis returns the module's exported genesis
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params:  k.GetParams(ctx),
		Intents: k.getAllIntents(ctx),
		Solvers: k.getAllSolvers(ctx),
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
