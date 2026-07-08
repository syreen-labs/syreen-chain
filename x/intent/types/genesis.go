package types

import "fmt"

// GenesisState defines the intent module's genesis state.
//
// In addition to Params/Intents/Solvers, it carries the full in-flight working
// set — Solutions, Chains, Strategies — plus the three monotonic ID counters so
// that an ExportGenesis→InitGenesis round-trip (state migration) is lossless and
// does NOT reset the counters (which would otherwise risk minting duplicate
// intent/chain/strategy IDs against surviving records).
type GenesisState struct {
	Params     Params        `json:"params"`
	Intents    []Intent      `json:"intents"`
	Solvers    []Solver      `json:"solvers"`
	Solutions  []Solution    `json:"solutions"`
	Chains     []IntentChain `json:"chains"`
	Strategies []Strategy    `json:"strategies"`

	// ID counters (raw stored counter values = highest ID issued so far). These
	// mirror the on-chain intent_counter / chain_counter / strategy_counter keys.
	NextIntentId   uint64 `json:"next_intent_id"`
	NextChainId    uint64 `json:"next_chain_id"`
	NextStrategyId uint64 `json:"next_strategy_id"`
}

func (gs *GenesisState) ProtoMessage()           {}
func (gs *GenesisState) Reset()                  { *gs = GenesisState{} }
func (gs *GenesisState) String() string {
	return fmt.Sprintf("intent genesis: %d intents, %d solvers, %d solutions, %d chains, %d strategies",
		len(gs.Intents), len(gs.Solvers), len(gs.Solutions), len(gs.Chains), len(gs.Strategies))
}
func (gs *GenesisState) XXX_MessageName() string { return "syreen.intent.GenesisState" }

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:     DefaultParams(),
		Intents:    []Intent{},
		Solvers:    []Solver{},
		Solutions:  []Solution{},
		Chains:     []IntentChain{},
		Strategies: []Strategy{},
	}
}

// Validate performs basic genesis state validation
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	intentIDs := make(map[string]bool)
	for _, intent := range gs.Intents {
		if intent.ID == "" {
			return fmt.Errorf("intent has empty ID")
		}
		if intentIDs[intent.ID] {
			return fmt.Errorf("duplicate intent ID: %s", intent.ID)
		}
		intentIDs[intent.ID] = true
		if intent.Creator == "" {
			return fmt.Errorf("intent %s has empty creator", intent.ID)
		}
		if intent.IntentType == "" {
			return fmt.Errorf("intent %s has empty intent type", intent.ID)
		}
	}
	solverAddrs := make(map[string]bool)
	for _, solver := range gs.Solvers {
		if solver.Address == "" {
			return fmt.Errorf("solver has empty address")
		}
		if solverAddrs[solver.Address] {
			return fmt.Errorf("duplicate solver address: %s", solver.Address)
		}
		solverAddrs[solver.Address] = true
	}
	return nil
}
