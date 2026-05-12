package keeper

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/mevprotection/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	k.SetParams(ctx, data.Params)
	for _, penalty := range data.Penalties {
		kvStore := k.storeService.OpenKVStore(ctx)
		bz, _ := json.Marshal(penalty)
		kvStore.Set(types.PenaltyKey(penalty.ValidatorAddr, penalty.BlockHeight), bz)
	}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params:    k.GetParams(ctx),
		Penalties: []types.MEVPenalty{}, // simplified: iterate store for full export
	}
}
