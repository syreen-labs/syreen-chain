package types

import "cosmossdk.io/math"

// Portfolio represents a user's aggregated portfolio state
type Portfolio struct {
	Address       string   `json:"address"`
	TotalValue    math.Int `json:"total_value"`
	TotalPnL      math.Int `json:"total_pnl"`
	RealizedPnL   math.Int `json:"realized_pnl"`
	UnrealizedPnL math.Int `json:"unrealized_pnl"`
	CostBasis     math.Int `json:"cost_basis"`   // cumulative cost of current holdings (amountIn)
	LastUpdated   int64    `json:"last_updated"` // block height
}

// PortfolioAsset represents a single asset in a user's portfolio
type PortfolioAsset struct {
	Address        string         `json:"address"`
	Denom          string         `json:"denom"`
	Amount         math.Int       `json:"amount"`
	Value          math.Int       `json:"value"`
	PnL            math.Int       `json:"pnl"`
	PercentOfTotal math.LegacyDec `json:"percent_of_total"`
}

// TradeRecord represents a single trade event
type TradeRecord struct {
	ID        uint64   `json:"id"`
	Address   string   `json:"address"`
	TradeType string   `json:"trade_type"` // swap, perp_open, perp_close, lending_deposit, lending_withdraw
	Denom     string   `json:"denom"`
	Amount    math.Int `json:"amount"`
	Price     math.Int `json:"price"` // in usyreen equivalent
	PnL       math.Int `json:"pnl"`
	Block     int64    `json:"block"`
}

// Activity represents a social feed activity
type Activity struct {
	ID           uint64   `json:"id"`
	ActivityType string   `json:"activity_type"` // large_swap, position_opened, position_closed, liquidation, whale_transfer, new_pool, large_deposit
	Address      string   `json:"address"`
	Description  string   `json:"description"`
	Amount       math.Int `json:"amount"`
	Denom        string   `json:"denom"`
	Block        int64    `json:"block"`
}

// Competition represents a trading competition
type Competition struct {
	ID              uint64   `json:"id"`
	Name            string   `json:"name"`
	Creator         string   `json:"creator"`
	StartBlock      int64    `json:"start_block"`
	EndBlock        int64    `json:"end_block"`
	PrizeDenom      string   `json:"prize_denom"`
	PrizePool       math.Int `json:"prize_pool"`
	EntryFee        math.Int `json:"entry_fee"`
	MaxParticipants uint64   `json:"max_participants"`
	Participants    uint64   `json:"participants"`
	Status          string   `json:"status"` // upcoming, active, ended, distributed
}

// CompetitionEntry represents a user's entry in a competition
type CompetitionEntry struct {
	Address       string         `json:"address"`
	CompetitionID uint64         `json:"competition_id"`
	StartingValue math.Int       `json:"starting_value"`
	CurrentValue  math.Int       `json:"current_value"`
	PnLPercent    math.LegacyDec `json:"pnl_percent"`
	Rank          uint64         `json:"rank"`
}

// GlobalActivityList stores the last N activity IDs in order
type GlobalActivityList struct {
	ActivityIDs []uint64 `json:"activity_ids"`
}

// Valid trade types
const (
	TradeTypeSwap            = "swap"
	TradeTypePerpOpen        = "perp_open"
	TradeTypePerpClose       = "perp_close"
	TradeTypeLendingDeposit  = "lending_deposit"
	TradeTypeLendingWithdraw = "lending_withdraw"
)

// Valid activity types
const (
	ActivityLargeSwap       = "large_swap"
	ActivityPositionOpened  = "position_opened"
	ActivityPositionClosed  = "position_closed"
	ActivityLiquidation     = "liquidation"
	ActivityWhaleTransfer   = "whale_transfer"
	ActivityNewPool         = "new_pool"
	ActivityLargeDeposit    = "large_deposit"
)

// Competition statuses
const (
	CompStatusUpcoming    = "upcoming"
	CompStatusActive      = "active"
	CompStatusEnded       = "ended"
	CompStatusDistributed = "distributed"
)

// GenesisState
type GenesisState struct {
	Portfolios   []Portfolio       `json:"portfolios"`
	Assets       []PortfolioAsset  `json:"assets"`
	Trades       []TradeRecord     `json:"trades"`
	Activities   []Activity        `json:"activities"`
	Competitions []Competition     `json:"competitions"`
	Entries      []CompetitionEntry `json:"entries"`
}

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Portfolios:   []Portfolio{},
		Assets:       []PortfolioAsset{},
		Trades:       []TradeRecord{},
		Activities:   []Activity{},
		Competitions: []Competition{},
		Entries:      []CompetitionEntry{},
	}
}

func (gs GenesisState) Validate() error { return nil }

// IsValidTradeType checks if the given trade type is valid
func IsValidTradeType(t string) bool {
	switch t {
	case TradeTypeSwap, TradeTypePerpOpen, TradeTypePerpClose, TradeTypeLendingDeposit, TradeTypeLendingWithdraw:
		return true
	}
	return false
}
