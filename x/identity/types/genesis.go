package types

// GenesisState defines the identity module's genesis state.
type GenesisState struct {
	Params     Params     `json:"params"`
	Identities []Identity `json:"identities"`
	Verifiers  []Verifier `json:"verifiers"`
}

// DefaultGenesis returns the default genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:     DefaultParams(),
		Identities: []Identity{},
		Verifiers:  []Verifier{},
	}
}

// Validate performs basic genesis state validation.
func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}
