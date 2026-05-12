package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/feemarket/types"
)

// InitGenesis initializes the feemarket module's state from a provided genesis state
func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	k.SetParams(ctx, data.Params)
	k.SetFeeState(ctx, data.FeeState)

	// Initialize fee lanes from params
	for _, lane := range data.Params.FeeLanes {
		k.SetFeeLane(ctx, lane)
	}
}

// ExportGenesis returns the feemarket module's exported genesis state
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params:   k.GetParams(ctx),
		FeeState: k.GetFeeState(ctx),
	}
}
