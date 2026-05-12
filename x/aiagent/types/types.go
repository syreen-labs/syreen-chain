package types

import "cosmossdk.io/math"

// AgentStatus tracks lifecycle
type AgentStatus string

const (
	AgentStatusActive  AgentStatus = "active"
	AgentStatusPaused  AgentStatus = "paused"
	AgentStatusDrained AgentStatus = "drained" // budget exhausted
)

// StrategyType determines the AI agent's trading logic
type StrategyType string

const (
	StrategyMomentum      StrategyType = "momentum"       // follow trading signals
	StrategyMeanReversion StrategyType = "mean_reversion"  // RSI-based contrarian
	StrategyGrid          StrategyType = "grid"            // buy/sell at fixed intervals
	StrategyDCA           StrategyType = "dca"             // dollar cost average
	StrategyArbitrage     StrategyType = "arbitrage"       // exploit cross-pool price gaps
)

// AgentConfig holds strategy-specific parameters
type AgentConfig struct {
	// Common
	PoolID          uint64         `json:"pool_id"`           // primary pool to trade
	PoolIDB         uint64         `json:"pool_id_b"`         // secondary pool (arbitrage only)
	InputDenom      string         `json:"input_denom"`       // denom agent trades FROM
	OutputDenom     string         `json:"output_denom"`      // denom agent trades TO
	TradePercent    math.LegacyDec `json:"trade_percent"`     // % of balance per trade (e.g. 0.10 = 10%)
	MaxSlippage     math.LegacyDec `json:"max_slippage"`      // max allowed slippage per trade
	IntervalBlocks  int64          `json:"interval_blocks"`   // min blocks between trades
	MaxTradesPerDay int64          `json:"max_trades_per_day"` // safety cap

	// Momentum
	MinSignalStrength string `json:"min_signal_strength"` // "BUY", "STRONG_BUY" etc.

	// Mean reversion
	RSIBuyThreshold  int64 `json:"rsi_buy_threshold"`  // buy when RSI < this (e.g. 30)
	RSISellThreshold int64 `json:"rsi_sell_threshold"` // sell when RSI > this (e.g. 70)

	// Grid
	GridLower  math.LegacyDec `json:"grid_lower"`  // lower bound price
	GridUpper  math.LegacyDec `json:"grid_upper"`  // upper bound price
	GridLevels int64          `json:"grid_levels"` // number of grid levels

	// Arbitrage
	MinProfitPercent math.LegacyDec `json:"min_profit_percent"` // min profit to trigger arb
}

// AgentStats tracks performance
type AgentStats struct {
	TotalTrades    int64          `json:"total_trades"`
	WinningTrades  int64          `json:"winning_trades"`
	TotalPnL       math.LegacyDec `json:"total_pnl"`       // approximate, in input denom
	TotalVolumeIn  math.Int       `json:"total_volume_in"`
	TotalVolumeOut math.Int       `json:"total_volume_out"`
	LastTradeBlock int64          `json:"last_trade_block"`
	TradesThisDay  int64          `json:"trades_this_day"`
	DayStartBlock  int64          `json:"day_start_block"`
}

// Agent is the core AI trading agent struct
type Agent struct {
	ID           uint64       `json:"id"`
	Owner        string       `json:"owner"`
	Name         string       `json:"name"`
	StrategyType StrategyType `json:"strategy_type"`
	Status       AgentStatus  `json:"status"`
	Config       AgentConfig  `json:"config"`
	Stats        AgentStats   `json:"stats"`
	CreatedAt    int64        `json:"created_at"`
	// AgentAddress is derived deterministically from ID
	AgentAddress string `json:"agent_address"`
}

// GenesisState
type GenesisState struct {
	Agents []Agent `json:"agents"`
	NextID uint64  `json:"next_id"`
}

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Agents: []Agent{},
		NextID: 1,
	}
}

func (gs GenesisState) Validate() error { return nil }
