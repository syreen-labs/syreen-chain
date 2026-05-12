package types

import (
	"cosmossdk.io/math"
)

// Pool represents an AMM liquidity pool with constant-product invariant (x*y=k).
type Pool struct {
	ID          uint64         `json:"id"`
	DenomA      string         `json:"denom_a"`
	DenomB      string         `json:"denom_b"`
	ReserveA    math.Int       `json:"reserve_a"`
	ReserveB    math.Int       `json:"reserve_b"`
	TotalShares math.Int       `json:"total_shares"`
	SwapFee     math.LegacyDec `json:"swap_fee"`
	Creator     string         `json:"creator"`
	CreatedAt   int64          `json:"created_at"`

	// Dynamic fee fields (programmable fee tiers)
	VolatilityFeeEnabled bool           `json:"volatility_fee_enabled"`
	BaseFee              math.LegacyDec `json:"base_fee,omitempty"`
	MaxFee               math.LegacyDec `json:"max_fee,omitempty"`
	VolatilityMultiplier math.LegacyDec `json:"volatility_multiplier,omitempty"`
}

// SortDenoms returns two denominations in alphabetical order.
func SortDenoms(a, b string) (string, string) {
	if a < b {
		return a, b
	}
	return b, a
}
