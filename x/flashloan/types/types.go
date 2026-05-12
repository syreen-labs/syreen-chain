package types

import "cosmossdk.io/math"

// FlashLoanPool represents a pool from which flash loans can be borrowed
type FlashLoanPool struct {
	Denom       string         `json:"denom"`
	Available   math.Int       `json:"available"`     // current available liquidity
	TotalLoaned math.Int       `json:"total_loaned"`  // cumulative amount loaned
	TotalFees   math.Int       `json:"total_fees"`    // cumulative fees earned
	FeeRate     math.LegacyDec `json:"fee_rate"`      // e.g. 0.0009 = 0.09%
	TotalShares math.Int       `json:"total_shares"`  // H-1: total LP shares outstanding
	Active      bool           `json:"active"`
}

// FlashPoolDeposit records an individual funder's LP share position in a flash pool (H-1 fix).
type FlashPoolDeposit struct {
	Depositor string   `json:"depositor"`
	Denom     string   `json:"denom"`   // pool identifier
	Amount    math.Int `json:"amount"`  // tokens originally deposited
	Shares    math.Int `json:"shares"`  // proportional shares of the pool
}

// FlashLoanRecord tracks an individual flash loan execution
type FlashLoanRecord struct {
	ID       uint64   `json:"id"`
	Borrower string   `json:"borrower"`
	Denom    string   `json:"denom"`
	Amount   math.Int `json:"amount"`
	Fee      math.Int `json:"fee"`
	Block    int64    `json:"block"`
	Success  bool     `json:"success"`
}

// FlashLoanStats tracks aggregate flash loan statistics
type FlashLoanStats struct {
	TotalExecuted  uint64   `json:"total_executed"`
	TotalVolume    math.Int `json:"total_volume"`     // cumulative borrowed amount across all denoms
	TotalFeesEarned math.Int `json:"total_fees_earned"` // cumulative fees across all denoms
}

// DefaultFlashLoanStats returns zero-value stats
func DefaultFlashLoanStats() FlashLoanStats {
	return FlashLoanStats{
		TotalExecuted:   0,
		TotalVolume:     math.ZeroInt(),
		TotalFeesEarned: math.ZeroInt(),
	}
}

// DefaultFeeRate is 0.09% (like Aave)
var DefaultFeeRate = math.LegacyNewDecWithPrec(9, 4) // 0.0009

// GenesisState
type GenesisState struct {
	Pools        []FlashLoanPool    `json:"pools"`
	Records      []FlashLoanRecord  `json:"records"`
	Deposits     []FlashPoolDeposit `json:"deposits"`
	Stats        FlashLoanStats     `json:"stats"`
	NextRecordID uint64             `json:"next_record_id"`
}

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Pools:    []FlashLoanPool{},
		Records:  []FlashLoanRecord{},
		Deposits: []FlashPoolDeposit{},
		Stats:    DefaultFlashLoanStats(),
	}
}

func (gs GenesisState) Validate() error { return nil }
