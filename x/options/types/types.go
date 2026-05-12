package types

import "cosmossdk.io/math"

// OptionType is call or put
type OptionType string

const (
	OptionTypeCall OptionType = "call"
	OptionTypePut  OptionType = "put"
)

// OptionStatus tracks lifecycle
type OptionStatus string

const (
	OptionStatusOpen      OptionStatus = "open"      // written, not yet bought
	OptionStatusActive    OptionStatus = "active"     // bought, not yet exercised/expired
	OptionStatusExercised OptionStatus = "exercised"
	OptionStatusExpired   OptionStatus = "expired"
	OptionStatusCancelled OptionStatus = "cancelled"
)

// Option represents an on-chain option contract
type Option struct {
	ID             uint64         `json:"id"`
	Writer         string         `json:"writer"`          // who wrote/sold the option
	Buyer          string         `json:"buyer"`           // who bought it (empty until bought)
	PoolID         uint64         `json:"pool_id"`         // DEX pool for price reference
	UnderlyingDenom string        `json:"underlying_denom"`
	QuoteDenom     string         `json:"quote_denom"`
	OptionType     OptionType     `json:"option_type"`
	StrikePrice    math.LegacyDec `json:"strike_price"`    // in quote per underlying
	Premium        math.Int       `json:"premium"`         // cost to buy, in quote denom base units
	Amount         math.Int       `json:"amount"`          // underlying amount controlled
	ExpiryBlock    int64          `json:"expiry_block"`
	Status         OptionStatus   `json:"status"`
	CreatedAt      int64          `json:"created_at"`
	ExercisedAt    int64          `json:"exercised_at"`
}

// GenesisState
type GenesisState struct {
	Options    []Option `json:"options"`
	NextID     uint64   `json:"next_id"`
}

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Options: []Option{},
		NextID:  1,
	}
}

func (gs GenesisState) Validate() error { return nil }
