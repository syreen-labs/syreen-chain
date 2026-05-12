package keeper

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/compute/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	k.SetParams(ctx, data.Params)

	// Restore codes
	var maxCodeID uint64
	for _, code := range data.Codes {
		k.setCodeInfo(ctx, code.CodeID, code.CodeInfo)
		k.setCodeBytecode(ctx, code.CodeID, code.CodeBytes)
		if code.CodeID > maxCodeID {
			maxCodeID = code.CodeID
		}
	}
	k.setNextCodeID(ctx, maxCodeID+1)

	// Restore contracts and their state
	for _, contract := range data.Contracts {
		k.setContractInfo(ctx, contract.ContractAddress, contract.ContractInfo)
		for _, model := range contract.ContractState {
			k.SetContractState(ctx, contract.ContractAddress, model.Key, model.Value)
		}
	}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	gs := &types.GenesisState{
		Params: k.GetParams(ctx),
	}

	// Export all codes
	nextCodeID := k.getNextCodeID(ctx)
	for id := uint64(1); id < nextCodeID; id++ {
		info, found := k.GetCodeInfo(ctx, id)
		if !found {
			continue
		}
		bytecode := k.getCodeBytecode(ctx, id)
		gs.Codes = append(gs.Codes, types.GenesisCode{
			CodeID:    id,
			CodeInfo:  info,
			CodeBytes: bytecode,
		})
	}

	// Export all contracts by iterating the contract prefix
	kvStore := k.storeService.OpenKVStore(ctx)
	contractPrefix := []byte(types.ContractKeyPrefix)
	iter, err := kvStore.Iterator(contractPrefix, append(contractPrefix, 0xFF))
	if err == nil {
		defer iter.Close()
		for ; iter.Valid(); iter.Next() {
			contractAddr := string(iter.Key()[len(contractPrefix):])
			var info types.ContractInfo
			if err := json.Unmarshal(iter.Value(), &info); err != nil {
				continue
			}

			// Export contract state
			var state []types.Model
			statePrefix := types.ContractStateIteratorPrefix(contractAddr)
			stateIter, sErr := kvStore.Iterator(statePrefix, append(statePrefix, 0xFF))
			if sErr == nil {
				for ; stateIter.Valid(); stateIter.Next() {
					key := make([]byte, len(stateIter.Key())-len(statePrefix))
					copy(key, stateIter.Key()[len(statePrefix):])
					state = append(state, types.Model{
						Key:   key,
						Value: stateIter.Value(),
					})
				}
				stateIter.Close()
			}

			gs.Contracts = append(gs.Contracts, types.GenesisContract{
				ContractAddress: contractAddr,
				ContractInfo:    info,
				ContractState:   state,
			})
		}
	}

	return gs
}
