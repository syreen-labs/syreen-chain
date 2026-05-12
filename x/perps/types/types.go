package types

import (
	"cosmossdk.io/math"
)

// Side represents long or short
type Side string

const (
	SideLong  Side = "long"
	SideShort Side = "short"
)

// PerpMarket represents a perpetual futures market
type PerpMarket struct {
	ID               uint64         `json:"id"`
	BaseDenom        string         `json:"base_denom"`        // e.g. usyreen (the asset being traded)
	QuoteDenom       string         `json:"quote_denom"`       // e.g. uusdc (margin/settlement)
	PoolID           uint64         `json:"pool_id"`           // DEX pool for mark price
	MaxLeverage      math.LegacyDec `json:"max_leverage"`      // e.g. 20x
	MaintenanceMargin math.LegacyDec `json:"maintenance_margin"` // e.g. 0.05 (5%)
	InitialMargin    math.LegacyDec `json:"initial_margin"`    // e.g. 0.10 (10% for 10x)
	TakerFee         math.LegacyDec `json:"taker_fee"`         // e.g. 0.001 (0.1%)
	MakerFee         math.LegacyDec `json:"maker_fee"`         // e.g. 0.0005 (0.05%)
	FundingInterval  int64          `json:"funding_interval"`  // blocks between funding (e.g. 7200 = ~1hr)
	MaxFundingRate   math.LegacyDec `json:"max_funding_rate"`  // max per interval (e.g. 0.01 = 1%)
	MaxOpenInterest  math.Int       `json:"max_open_interest"` // max total notional
	LongOpenInterest math.Int       `json:"long_open_interest"`
	ShortOpenInterest math.Int      `json:"short_open_interest"`
	Active           bool           `json:"active"`
	Creator          string         `json:"creator"`
	// Price dampening fields (C-4 fix): updated each block in ProcessFundingRates
	LastMarkPrice  math.LegacyDec `json:"last_mark_price"`  // dampened mark price from prior block
	LastMarkHeight int64          `json:"last_mark_height"` // block height of last mark price update
	// AllowUndercapitalized allows opening positions even if insurance fund is below minimum (H-3)
	AllowUndercapitalized bool `json:"allow_undercapitalized"`
}

// Position represents a trader's perpetual position
type Position struct {
	Address          string         `json:"address"`
	MarketID         uint64         `json:"market_id"`
	Side             Side           `json:"side"`               // long or short
	Size             math.LegacyDec `json:"size"`               // position size in base units
	EntryPrice       math.LegacyDec `json:"entry_price"`        // average entry price
	Margin           math.Int       `json:"margin"`             // collateral deposited (in quote denom)
	Leverage         math.LegacyDec `json:"leverage"`           // effective leverage
	UnrealizedPnL    math.LegacyDec `json:"unrealized_pnl"`     // current unrealized PnL
	CumulativeFunding math.LegacyDec `json:"cumulative_funding"` // funding payments accumulated
	OpenedAt         int64          `json:"opened_at"`          // block height
	LastFundingBlock int64          `json:"last_funding_block"` // last block funding was applied
}

// FundingState tracks the funding rate for a market
type FundingState struct {
	MarketID            uint64         `json:"market_id"`
	CurrentFundingRate  math.LegacyDec `json:"current_funding_rate"`  // positive = longs pay shorts
	CumulativeFunding   math.LegacyDec `json:"cumulative_funding"`    // running total
	LastFundingBlock    int64          `json:"last_funding_block"`
	LastFundingTime     int64          `json:"last_funding_time"`
}

// InsuranceFund holds funds to cover losses from liquidations
type InsuranceFund struct {
	Balance math.Int `json:"balance"` // in quote denom
}

// GenesisState
type GenesisState struct {
	Markets      []PerpMarket   `json:"markets"`
	Positions    []Position     `json:"positions"`
	Fundings     []FundingState `json:"fundings"`
	Insurance    InsuranceFund  `json:"insurance"`
	NextMarketID uint64         `json:"next_market_id"`
}

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Markets:   []PerpMarket{},
		Positions: []Position{},
		Fundings:  []FundingState{},
		Insurance: InsuranceFund{Balance: math.ZeroInt()},
	}
}

func (gs GenesisState) Validate() error { return nil }

// Default market parameters
var (
	DefaultMaxLeverage      = math.LegacyNewDec(20)       // 20x
	DefaultMaintenanceMargin = math.LegacyNewDecWithPrec(5, 2)  // 5%
	DefaultInitialMargin    = math.LegacyNewDecWithPrec(10, 2) // 10%
	DefaultTakerFee         = math.LegacyNewDecWithPrec(1, 3)  // 0.1%
	DefaultMakerFee         = math.LegacyNewDecWithPrec(5, 4)  // 0.05%
	DefaultFundingInterval  = int64(7200)                       // ~1 hour at 500ms blocks
	DefaultMaxFundingRate   = math.LegacyNewDecWithPrec(1, 2)  // 1% per interval
)
