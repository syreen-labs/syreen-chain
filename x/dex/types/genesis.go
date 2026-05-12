package types

import "fmt"

// GenesisState defines the dex module's genesis state.
type GenesisState struct {
	Params     Params `json:"params"`
	Pools      []Pool `json:"pools"`
	NextPoolID uint64 `json:"next_pool_id"`
}

func (gs *GenesisState) ProtoMessage()          {}
func (gs *GenesisState) Reset()                 { *gs = GenesisState{} }
func (gs *GenesisState) String() string         { return fmt.Sprintf("dex genesis: %d pools", len(gs.Pools)) }
func (gs *GenesisState) XXX_MessageName() string { return "syreen.dex.GenesisState" }

// DefaultGenesis returns the default genesis state for the dex module.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:     DefaultParams(),
		Pools:      []Pool{},
		NextPoolID: 1,
	}
}

// Validate performs basic genesis state validation.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	seenIDs := map[uint64]bool{}
	seenPairs := map[string]bool{}
	for _, pool := range gs.Pools {
		if seenIDs[pool.ID] {
			return fmt.Errorf("duplicate pool ID: %d", pool.ID)
		}
		seenIDs[pool.ID] = true

		a, b := SortDenoms(pool.DenomA, pool.DenomB)
		pairKey := a + "/" + b
		if seenPairs[pairKey] {
			return fmt.Errorf("duplicate pool for denom pair: %s", pairKey)
		}
		seenPairs[pairKey] = true

		if pool.DenomA == pool.DenomB {
			return fmt.Errorf("pool %d has same denom on both sides: %s", pool.ID, pool.DenomA)
		}
		if pool.ReserveA.IsNegative() || pool.ReserveB.IsNegative() {
			return fmt.Errorf("pool %d has negative reserves", pool.ID)
		}
		if pool.TotalShares.IsNegative() {
			return fmt.Errorf("pool %d has negative total shares", pool.ID)
		}
	}
	if gs.NextPoolID == 0 {
		return fmt.Errorf("next pool ID must be positive")
	}
	return nil
}
