package types

import "fmt"

// DenomAuthorityMetadata specifies metadata about the admin of a denom
type DenomAuthorityMetadata struct {
	Admin    string `json:"admin"`
	MaxSupply string `json:"max_supply,omitempty"` // max mintable supply; "0" or "" means unlimited
}

// GenesisDenom defines a genesis denom with its authority metadata
type GenesisDenom struct {
	Denom             string                 `json:"denom"`
	AuthorityMetadata DenomAuthorityMetadata `json:"authority_metadata"`
}

// GenesisState defines the tokenfactory module's genesis state
type GenesisState struct {
	Params        Params         `json:"params"`
	FactoryDenoms []GenesisDenom `json:"factory_denoms"`
}

func (gs *GenesisState) ProtoMessage()             {}
func (gs *GenesisState) Reset()                    { *gs = GenesisState{} }
func (gs *GenesisState) String() string            { return fmt.Sprintf("tokenfactory genesis: %d denoms", len(gs.FactoryDenoms)) }
func (gs *GenesisState) XXX_MessageName() string   { return "syreen.tokenfactory.GenesisState" }

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:       DefaultParams(),
		FactoryDenoms: []GenesisDenom{},
	}
}

func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	seenDenoms := map[string]bool{}
	for _, denom := range gs.FactoryDenoms {
		if seenDenoms[denom.Denom] {
			return fmt.Errorf("duplicate denom: %s", denom.Denom)
		}
		seenDenoms[denom.Denom] = true
	}
	return nil
}
