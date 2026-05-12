package types

import "fmt"

// GenesisState defines the mevprotection module's genesis state
type GenesisState struct {
	Params    Params     `json:"params"`
	Penalties []MEVPenalty `json:"penalties"`
}

func (gs *GenesisState) ProtoMessage()           {}
func (gs *GenesisState) Reset()                  { *gs = GenesisState{} }
func (gs *GenesisState) String() string          { return fmt.Sprintf("mevprotection genesis: %d penalties", len(gs.Penalties)) }
func (gs *GenesisState) XXX_MessageName() string { return "syreen.mevprotection.GenesisState" }

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:    DefaultParams(),
		Penalties: []MEVPenalty{},
	}
}

func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	for _, penalty := range gs.Penalties {
		if penalty.ValidatorAddr == "" {
			return fmt.Errorf("penalty has empty validator address")
		}
		if penalty.PenaltyType == "" {
			return fmt.Errorf("penalty has empty penalty type")
		}
	}
	return nil
}
