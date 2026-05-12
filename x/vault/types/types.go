package types

import "cosmossdk.io/math"

// Strategy types
const (
	StrategyAutoCompound = "auto_compound"
	StrategyBalanced     = "balanced"
	StrategyAggressive   = "aggressive"
	StrategyConservative = "conservative"
)

// ValidStrategy returns true if the strategy type is valid
func ValidStrategy(s string) bool {
	switch s {
	case StrategyAutoCompound, StrategyBalanced, StrategyAggressive, StrategyConservative:
		return true
	}
	return false
}

// StrategyYieldMultiplier returns the simulated yield multiplier per compound cycle
// auto_compound=100%, balanced=75%, aggressive=150%, conservative=50%
func StrategyYieldMultiplier(s string) math.LegacyDec {
	switch s {
	case StrategyAutoCompound:
		return math.LegacyNewDecWithPrec(100, 2) // 1.00
	case StrategyBalanced:
		return math.LegacyNewDecWithPrec(75, 2) // 0.75
	case StrategyAggressive:
		return math.LegacyNewDecWithPrec(150, 2) // 1.50
	case StrategyConservative:
		return math.LegacyNewDecWithPrec(50, 2) // 0.50
	}
	return math.LegacyOneDec()
}

// Vault represents a yield vault
type Vault struct {
	ID                uint64         `json:"id"`
	Name              string         `json:"name"`
	Creator           string         `json:"creator"`
	DepositDenom      string         `json:"deposit_denom"`
	TotalDeposited    math.Int       `json:"total_deposited"`
	TotalShares       math.Int       `json:"total_shares"`
	StrategyType      string         `json:"strategy_type"`
	TargetPoolIDs     []uint64       `json:"target_pool_ids"`
	PerformanceFee    math.LegacyDec `json:"performance_fee"`    // e.g. 0.10 = 10% to creator
	ProtocolFee       math.LegacyDec `json:"protocol_fee"`       // e.g. 0.02 = 2% to protocol
	LastCompoundBlock int64          `json:"last_compound_block"`
	AccYieldPerShare  math.LegacyDec `json:"acc_yield_per_share"`
	Active            bool           `json:"active"`
}

// VaultDeposit represents a user's deposit in a vault
type VaultDeposit struct {
	Address         string         `json:"address"`
	VaultID         uint64         `json:"vault_id"`
	Shares          math.Int       `json:"shares"`
	DepositedAmount math.Int       `json:"deposited_amount"`
	YieldIndex      math.LegacyDec `json:"yield_index"` // snapshot of AccYieldPerShare at deposit
	DepositedAt     int64          `json:"deposited_at"`
}

// GenesisState
type GenesisState struct {
	Vaults   []Vault        `json:"vaults"`
	Deposits []VaultDeposit `json:"deposits"`
}

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Vaults:   []Vault{},
		Deposits: []VaultDeposit{},
	}
}

func (gs GenesisState) Validate() error {
	return nil
}
