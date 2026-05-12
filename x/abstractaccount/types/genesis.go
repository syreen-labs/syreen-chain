package types

import "fmt"

// GenesisState defines the abstractaccount module's genesis state
type GenesisState struct {
	Params Params `json:"params"`
}

func (gs *GenesisState) ProtoMessage()           {}
func (gs *GenesisState) Reset()                  { *gs = GenesisState{} }
func (gs *GenesisState) String() string          { return fmt.Sprintf("abstractaccount genesis: params=%+v", gs.Params) }
func (gs *GenesisState) XXX_MessageName() string { return "syreen.abstractaccount.GenesisState" }

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}
