package types

import "cosmossdk.io/math"

// MarketStatus represents the state of a prediction market
type MarketStatus string

const (
	MarketStatusOpen       MarketStatus = "open"
	MarketStatusClosed     MarketStatus = "closed"
	MarketStatusResolvedYes MarketStatus = "resolved_yes"
	MarketStatusResolvedNo  MarketStatus = "resolved_no"
	MarketStatusVoided     MarketStatus = "voided"
)

// Market represents a binary outcome prediction market
type Market struct {
	ID              uint64           `json:"id"`
	Question        string           `json:"question"`
	Creator         string           `json:"creator"`
	Resolver        string           `json:"resolver"` // informational only; resolution authority is always chain governance
	QuoteDenom      string           `json:"quote_denom"`
	ResolutionBlock int64            `json:"resolution_block"`
	Status          MarketStatus     `json:"status"`
	Outcome         string           `json:"outcome"`      // set at resolution: "yes", "no", or "void"
	PayoutRatio     math.LegacyDec   `json:"payout_ratio"` // ratio applied to winning shares at payout (1.0 = full, <1.0 = proportional)
	YesShares       math.Int         `json:"yes_shares"`   // total YES shares in the pool
	NoShares        math.Int         `json:"no_shares"`    // total NO shares in the pool
	Liquidity       math.Int         `json:"liquidity"`    // initial liquidity parameter
	TotalVolume     math.Int         `json:"total_volume"` // cumulative trading volume
	CreatedAt       int64            `json:"created_at"`
}

// Position represents a user's share holdings in a market
type Position struct {
	Address    string   `json:"address"`
	MarketID   uint64   `json:"market_id"`
	YesShares  math.Int `json:"yes_shares"`
	NoShares   math.Int `json:"no_shares"`
	TotalSpent math.Int `json:"total_spent"`
	Claimed    bool     `json:"claimed"`
}

// Resolution records the final outcome of a market
type Resolution struct {
	MarketID   uint64 `json:"market_id"`
	Outcome    string `json:"outcome"` // "yes", "no", "void"
	ResolvedAt int64  `json:"resolved_at"`
	ResolvedBy string `json:"resolved_by"`
}

// GenesisState defines the predict module genesis state
type GenesisState struct {
	Markets     []Market     `json:"markets"`
	Positions   []Position   `json:"positions"`
	Resolutions []Resolution `json:"resolutions"`
}

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Markets:     []Market{},
		Positions:   []Position{},
		Resolutions: []Resolution{},
	}
}

func (gs GenesisState) Validate() error {
	return nil
}
