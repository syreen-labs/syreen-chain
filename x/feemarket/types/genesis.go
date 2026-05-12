package types

import "fmt"

// GenesisState defines the feemarket module's genesis state
type GenesisState struct {
	Params   Params   `json:"params"`
	FeeState FeeState `json:"fee_state"`
}

func (gs *GenesisState) ProtoMessage()           {}
func (gs *GenesisState) Reset()                  { *gs = GenesisState{} }
func (gs *GenesisState) String() string          { return fmt.Sprintf("feemarket genesis: baseFee=%s", gs.FeeState.BaseFee) }
func (gs *GenesisState) XXX_MessageName() string { return "syreen.feemarket.GenesisState" }

// DefaultGenesis returns the default genesis state for the feemarket module
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:   DefaultParams(),
		FeeState: DefaultFeeState(),
	}
}

// Validate validates the genesis state
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}
	if gs.FeeState.BaseFee.IsNegative() {
		return fmt.Errorf("base fee cannot be negative: %s", gs.FeeState.BaseFee)
	}
	if gs.FeeState.MinBaseFee.IsNegative() {
		return fmt.Errorf("min base fee cannot be negative: %s", gs.FeeState.MinBaseFee)
	}
	if gs.FeeState.MaxBaseFee.IsNegative() || gs.FeeState.MaxBaseFee.IsZero() {
		return fmt.Errorf("max base fee must be positive: %s", gs.FeeState.MaxBaseFee)
	}
	if gs.FeeState.BlockGasTarget <= 0 {
		return fmt.Errorf("block gas target must be positive: %d", gs.FeeState.BlockGasTarget)
	}
	return nil
}
