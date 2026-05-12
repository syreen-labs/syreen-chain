package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/identity/types"
)

// InitGenesis initializes the module's state from a provided genesis state
func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	k.SetParams(ctx, data.Params)

	for _, identity := range data.Identities {
		k.SetIdentity(ctx, identity)
	}
	for _, verifier := range data.Verifiers {
		k.SetVerifier(ctx, verifier)
	}
}

// ExportGenesis returns the module's exported genesis
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params:     k.GetParams(ctx),
		Identities: k.GetAllIdentities(ctx),
		Verifiers:  k.GetAllVerifiers(ctx),
	}
}
