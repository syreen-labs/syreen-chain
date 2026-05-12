package types

import "fmt"

// GenesisCode contains code bytecode and metadata for genesis
type GenesisCode struct {
	CodeID   uint64   `json:"code_id"`
	CodeInfo CodeInfo `json:"code_info"`
	CodeBytes []byte  `json:"code_bytes"`
}

// GenesisContract contains contract info and state for genesis
type GenesisContract struct {
	ContractAddress string         `json:"contract_address"`
	ContractInfo    ContractInfo   `json:"contract_info"`
	ContractState   []Model        `json:"contract_state"`
}

// GenesisState defines the compute module's genesis state
type GenesisState struct {
	Params    Params            `json:"params"`
	Codes     []GenesisCode     `json:"codes"`
	Contracts []GenesisContract `json:"contracts"`
}

func (gs *GenesisState) ProtoMessage()           {}
func (gs *GenesisState) Reset()                  { *gs = GenesisState{} }
func (gs *GenesisState) String() string          { return fmt.Sprintf("compute genesis: %d codes, %d contracts", len(gs.Codes), len(gs.Contracts)) }
func (gs *GenesisState) XXX_MessageName() string { return "syreen.compute.GenesisState" }

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:    DefaultParams(),
		Codes:     []GenesisCode{},
		Contracts: []GenesisContract{},
	}
}

func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	seenCodes := map[uint64]bool{}
	for _, code := range gs.Codes {
		if seenCodes[code.CodeID] {
			return fmt.Errorf("duplicate code ID: %d", code.CodeID)
		}
		seenCodes[code.CodeID] = true
		if code.CodeID == 0 {
			return fmt.Errorf("code ID cannot be zero")
		}
	}
	seenContracts := map[string]bool{}
	for _, contract := range gs.Contracts {
		if seenContracts[contract.ContractAddress] {
			return fmt.Errorf("duplicate contract address: %s", contract.ContractAddress)
		}
		seenContracts[contract.ContractAddress] = true
		if !seenCodes[contract.ContractInfo.CodeID] {
			return fmt.Errorf("contract %s references unknown code ID %d", contract.ContractAddress, contract.ContractInfo.CodeID)
		}
	}
	return nil
}
