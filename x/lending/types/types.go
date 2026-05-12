package types

import "cosmossdk.io/math"

// LendingPool represents a pool where users deposit tokens to earn interest
type LendingPool struct {
	ID                uint64         `json:"id"`
	Denom             string         `json:"denom"`              // token denom for this pool
	TotalDeposited    math.Int       `json:"total_deposited"`    // total tokens deposited
	TotalBorrowed     math.Int       `json:"total_borrowed"`     // total tokens borrowed
	DepositAPY        math.LegacyDec `json:"deposit_apy"`        // current deposit interest rate
	BorrowAPY         math.LegacyDec `json:"borrow_apy"`         // current borrow interest rate
	CollateralFactor  math.LegacyDec `json:"collateral_factor"`  // e.g. 0.75 = can borrow 75% of collateral value
	LiquidationBonus  math.LegacyDec `json:"liquidation_bonus"`  // e.g. 0.05 = 5% bonus for liquidators
	LiquidationThreshold math.LegacyDec `json:"liquidation_threshold"` // e.g. 0.80 = liquidate below 80% health
	ReserveRatio      math.LegacyDec `json:"reserve_ratio"`      // portion of interest kept as reserves
	TotalReserves     math.Int       `json:"total_reserves"`     // accumulated reserve tokens (from interest split)
	BadDebt           math.Int       `json:"bad_debt"`           // uncovered debt from underwater liquidations
	UtilizationRate   math.LegacyDec `json:"utilization_rate"`   // borrowed/deposited ratio
	AccInterestPerShare math.LegacyDec `json:"acc_interest_per_share"` // accumulated interest per deposited token
	AccBorrowIndex    math.LegacyDec `json:"acc_borrow_index"`   // accumulated borrow interest index
	LastUpdateBlock   int64          `json:"last_update_block"`
	Active            bool           `json:"active"`
	DexPoolID         uint64         `json:"dex_pool_id"`        // for price oracle
	PriceDenom        string         `json:"price_denom"`        // quote denom for pricing (e.g. uusdc)
}

// Deposit represents a user's deposit in a lending pool
type Deposit struct {
	Address       string         `json:"address"`
	PoolID        uint64         `json:"pool_id"`
	Amount        math.Int       `json:"amount"`         // tokens deposited
	InterestIndex math.LegacyDec `json:"interest_index"` // snapshot of AccInterestPerShare at deposit
	DepositedAt   int64          `json:"deposited_at"`
}

// Borrow represents an active borrow position
type Borrow struct {
	ID              uint64         `json:"id"`
	Borrower        string         `json:"borrower"`
	BorrowPoolID    uint64         `json:"borrow_pool_id"`    // pool being borrowed from
	BorrowAmount    math.Int       `json:"borrow_amount"`     // original borrow amount
	BorrowIndex     math.LegacyDec `json:"borrow_index"`      // snapshot of AccBorrowIndex at borrow time
	CollateralPoolID uint64        `json:"collateral_pool_id"` // pool where collateral is deposited
	CollateralAmount math.Int      `json:"collateral_amount"` // collateral locked
	BorrowedAt      int64          `json:"borrowed_at"`
	// M-1: block height at which the position first became unhealthy.
	// 0 means the position is currently considered healthy.
	UnhealthySinceBlock int64       `json:"unhealthy_since_block"`
}

// GenesisState
type GenesisState struct {
	Pools        []LendingPool `json:"pools"`
	Deposits     []Deposit     `json:"deposits"`
	Borrows      []Borrow      `json:"borrows"`
	NextPoolID   uint64        `json:"next_pool_id"`
	NextBorrowID uint64        `json:"next_borrow_id"`
}

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Pools:    []LendingPool{},
		Deposits: []Deposit{},
		Borrows:  []Borrow{},
	}
}

func (gs GenesisState) Validate() error { return nil }

// Interest rate model constants (similar to Aave/Compound)
var (
	BaseRate         = math.LegacyNewDecWithPrec(2, 2)  // 2% base rate
	Slope1           = math.LegacyNewDecWithPrec(4, 2)  // 4% slope below optimal
	Slope2           = math.LegacyNewDec(3)             // 300% slope above optimal (steep)
	OptimalUtil      = math.LegacyNewDecWithPrec(80, 2) // 80% optimal utilization
	InterestPrecision = math.LegacyNewDec(1_000_000_000_000)
	BlocksPerYear    = math.LegacyNewDec(63_072_000)    // ~500ms blocks

	// C-1 / M-1 oracle and liquidation safety constants.
	// MaxPriceDeviation: maximum allowed |spot - twap| / twap before borrow/liquidate
	// is rejected (5%).
	MaxPriceDeviation = math.LegacyNewDecWithPrec(5, 2)
	// DefaultLiquidationBonus: default liquidator bonus on new pools (2%, M-1).
	DefaultLiquidationBonus = math.LegacyNewDecWithPrec(2, 2)
)

const (
	// TWAPWindow: number of recent samples kept per pool for TWAP averaging.
	TWAPWindow = 30
	// LiquidationDelayBlocks: required number of blocks an account must remain
	// unhealthy before a liquidation is permitted (M-1).
	LiquidationDelayBlocks int64 = 10
)

// TWAPSample is a single oracle observation stored in lending state.
type TWAPSample struct {
	Block int64          `json:"block"`
	Price math.LegacyDec `json:"price"`
}

// TWAPSeries holds the most recent samples for a (lending pool) collateral.
type TWAPSeries struct {
	PoolID  uint64       `json:"pool_id"`
	Samples []TWAPSample `json:"samples"`
}
