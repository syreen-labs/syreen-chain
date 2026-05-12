package keeper

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/tokenfactory/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	k.SetParams(ctx, data.Params)
	for _, denom := range data.FactoryDenoms {
		k.SetDenomAuthorityMetadata(ctx, denom.Denom, denom.AuthorityMetadata)
		// Re-index by creator
		creator, _, err := types.DeconstructDenom(denom.Denom)
		if err == nil {
			k.addDenomByCreator(ctx, creator, denom.Denom)
		}
	}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	gs := &types.GenesisState{
		Params:        k.GetParams(ctx),
		FactoryDenoms: k.getAllFactoryDenoms(ctx),
	}
	return gs
}

// getAllFactoryDenoms iterates the denom authority store and returns all factory denoms.
func (k Keeper) getAllFactoryDenoms(ctx sdk.Context) []types.GenesisDenom {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.DenomsPrefixKey)

	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var denoms []types.GenesisDenom
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		// Strip the prefix to get the denom name
		denom := string(key[len(prefix):])

		var metadata types.DenomAuthorityMetadata
		if err := json.Unmarshal(iter.Value(), &metadata); err != nil {
			continue
		}

		denoms = append(denoms, types.GenesisDenom{
			Denom:             denom,
			AuthorityMetadata: metadata,
		})
	}

	return denoms
}

// prefixEndBytes returns the end key for a prefix scan (increment last byte).
func prefixEndBytes(prefix []byte) []byte {
	if len(prefix) == 0 {
		return nil
	}
	end := make([]byte, len(prefix))
	copy(end, prefix)
	for i := len(end) - 1; i >= 0; i-- {
		end[i]++
		if end[i] != 0 {
			return end
		}
	}
	return nil // overflow: prefix is all 0xFF
}
