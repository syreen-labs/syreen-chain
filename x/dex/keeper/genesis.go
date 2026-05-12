package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/dex/types"
)

// InitGenesis initializes the dex module genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	k.SetParams(ctx, data.Params)
	k.SetNextPoolID(ctx, data.NextPoolID)

	for _, pool := range data.Pools {
		k.SetPool(ctx, pool)
		k.SetPoolByDenomPair(ctx, pool.DenomA, pool.DenomB, pool.ID)
	}
}

// ExportGenesis exports the dex module genesis state.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params:     k.GetParams(ctx),
		Pools:      k.GetAllPools(ctx),
		NextPoolID: k.GetNextPoolID(ctx),
	}
}
