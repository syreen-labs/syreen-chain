package types

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/protobuf/encoding/protowire"
)

// StrategyStatus represents the lifecycle status of a deployed strategy
type StrategyStatus string

const (
	StrategyStatusActive    StrategyStatus = "active"
	StrategyStatusCompleted StrategyStatus = "completed"
	StrategyStatusCancelled StrategyStatus = "cancelled"
	StrategyStatusFailed    StrategyStatus = "failed"
)

// StrategyTemplate names
const (
	StrategySafeAccumulate = "safe_accumulate"
	StrategyMomentumRide   = "momentum_ride"
	StrategyGridTrading    = "grid_trading"
	StrategySwingTrade     = "swing_trade"
	StrategyProtectiveSell = "protective_sell"
	StrategyScalp          = "scalp"

	// Tier 1 AI-reactive templates (consume dex signal/sentiment state)
	StrategySentimentContrarian = "sentiment_contrarian"
	StrategySmartDCA            = "smart_dca"
	StrategySignalComposite     = "signal_composite"
	StrategyRSIMeanReversion    = "rsi_mean_reversion"
	StrategyVolatilityBreakout  = "volatility_breakout"
	StrategyTrendConfluence     = "trend_confluence"
)

// RiskLevel defines the risk tolerance for a strategy
type RiskLevel string

const (
	RiskConservative RiskLevel = "conservative"
	RiskModerate     RiskLevel = "moderate"
	RiskAggressive   RiskLevel = "aggressive"
)

// StrategyPrefix is the store prefix for strategies
const StrategyPrefix = "strategy/"

// StrategyKey returns the store key for a strategy by ID
func StrategyKey(strategyID string) []byte {
	return append([]byte(StrategyPrefix), []byte(strategyID)...)
}

// StrategyCounterKey is the key for the strategy ID counter
func StrategyCounterKey() []byte {
	return []byte("strategy_counter")
}

// StrategyTemplate describes a pre-built strategy that users can deploy
type StrategyTemplate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	MinSteps    int    `json:"min_steps"`
	MaxSteps    int    `json:"max_steps"`
}

// Strategy represents a deployed auto-strategy and the chain it created
type Strategy struct {
	ID           string         `json:"id"`
	Creator      string         `json:"creator"`
	TemplateName string         `json:"template_name"`
	ChainID      string         `json:"chain_id"`
	Status       StrategyStatus `json:"status"`
	PoolID       uint64         `json:"pool_id"`
	InputDenom   string         `json:"input_denom"`
	OutputDenom  string         `json:"output_denom"`
	TotalBudget  math.Int       `json:"total_budget"`
	RiskLevel    RiskLevel      `json:"risk_level"`
	CreatedAt    int64          `json:"created_at"`
}

func (m *Strategy) ProtoMessage()           {}
func (m *Strategy) Reset()                  { *m = Strategy{} }
func (m *Strategy) String() string          { return fmt.Sprintf("strategy(%s): template=%s chain=%s status=%s", m.ID, m.TemplateName, m.ChainID, m.Status) }
func (m *Strategy) XXX_MessageName() string { return "syreen.intent.Strategy" }

// RiskParams returns stop-loss and take-profit percentages based on risk level
func GetRiskParams(risk RiskLevel) (stopLossPct, takeProfitPct, gridSpreadPct math.LegacyDec) {
	switch risk {
	case RiskConservative:
		return math.LegacyMustNewDecFromStr("0.05"),
			math.LegacyMustNewDecFromStr("0.10"),
			math.LegacyMustNewDecFromStr("0.02")
	case RiskAggressive:
		return math.LegacyMustNewDecFromStr("0.15"),
			math.LegacyMustNewDecFromStr("0.30"),
			math.LegacyMustNewDecFromStr("0.05")
	default: // moderate
		return math.LegacyMustNewDecFromStr("0.10"),
			math.LegacyMustNewDecFromStr("0.20"),
			math.LegacyMustNewDecFromStr("0.03")
	}
}

// GetAvailableTemplates returns all available strategy templates
func GetAvailableTemplates() []StrategyTemplate {
	return []StrategyTemplate{
		{
			Name:        StrategySafeAccumulate,
			Description: "DCA buy with stop-loss protection. Gradually accumulates a token over time with automatic stop-loss if price crashes.",
			MinSteps:    2,
			MaxSteps:    2,
		},
		{
			Name:        StrategyMomentumRide,
			Description: "Buy on dip, take profit on rise, stop-loss for safety. A 3-step strategy: limit buy at a discount, then take-profit or stop-loss.",
			MinSteps:    3,
			MaxSteps:    3,
		},
		{
			Name:        StrategyGridTrading,
			Description: "Place buy and sell orders at regular price intervals. Creates a grid of alternating buy/sell orders around the current price.",
			MinSteps:    2,
			MaxSteps:    10,
		},
		{
			Name:        StrategySwingTrade,
			Description: "Buy low, sell high. A 2-step strategy: limit buy below current price, then limit sell above.",
			MinSteps:    2,
			MaxSteps:    2,
		},
		{
			Name:        StrategyProtectiveSell,
			Description: "Sell with trailing stop-loss. Takes profit at a target price with a stop-loss safety net.",
			MinSteps:    2,
			MaxSteps:    2,
		},
		{
			Name:        StrategyScalp,
			Description: "Quick small-profit trades with tight stop-loss. Fast limit buy then immediate take-profit at a small margin.",
			MinSteps:    2,
			MaxSteps:    2,
		},
		{
			Name:        StrategySentimentContrarian,
			Description: "Contrarian entry on extreme fear, exit on extreme greed. Uses the live pool fear/greed index to time entry and take-profit.",
			MinSteps:    4,
			MaxSteps:    4,
		},
		{
			Name:        StrategySmartDCA,
			Description: "Multi-tranche DCA that scales in more aggressively as fear deepens. Buys at base pace, adds weight when fearful, buys hard at extreme fear.",
			MinSteps:    3,
			MaxSteps:    3,
		},
		{
			Name:        StrategySignalComposite,
			Description: "Rides the dex composite signal: buys on entry, takes profit on SELL signal, stops out on STRONG_SELL signal.",
			MinSteps:    3,
			MaxSteps:    3,
		},
		{
			Name:        StrategyRSIMeanReversion,
			Description: "Classic RSI mean-reversion: anchor entry plus a larger buy when RSI is oversold (<30), take-profit when RSI is overbought (>70).",
			MinSteps:    4,
			MaxSteps:    4,
		},
		{
			Name:        StrategyVolatilityBreakout,
			Description: "Waits for a volatility spike, then rides the breakout with a wider take-profit and a standard stop-loss.",
			MinSteps:    4,
			MaxSteps:    4,
		},
		{
			Name:        StrategyTrendConfluence,
			Description: "Only enters when sentiment turns bullish AND the composite signal confirms BUY. Cascades conditions to require both signals to fire.",
			MinSteps:    5,
			MaxSteps:    5,
		},
	}
}

// ValidateRiskLevel checks if the risk level is valid
func ValidateRiskLevel(risk RiskLevel) error {
	switch risk {
	case RiskConservative, RiskModerate, RiskAggressive:
		return nil
	default:
		return fmt.Errorf("invalid risk level: %s (must be conservative, moderate, or aggressive)", risk)
	}
}

// ValidateTemplateName checks if the template name is valid
func ValidateTemplateName(name string) error {
	for _, t := range GetAvailableTemplates() {
		if t.Name == name {
			return nil
		}
	}
	return fmt.Errorf("unknown strategy template: %s", name)
}

// ---------------------------------------------------------------------------
// Msgs
// ---------------------------------------------------------------------------

const (
	TypeMsgCreateStrategy = "create_strategy"
	TypeMsgCancelStrategy = "cancel_strategy"
)

// MsgCreateStrategy deploys a strategy from a template
type MsgCreateStrategy struct {
	Creator      string    `json:"creator"`
	TemplateName string    `json:"template_name"`
	PoolID       uint64    `json:"pool_id"`
	InputDenom   string    `json:"input_denom"`
	OutputDenom  string    `json:"output_denom"`
	TotalBudget  math.Int  `json:"total_budget"`
	RiskLevel    RiskLevel `json:"risk_level"`
	ExpiryBlocks uint64    `json:"expiry_blocks"`
	// DCA-specific params (used by safe_accumulate)
	NumExecutions  uint64 `json:"num_executions,omitempty"`
	IntervalBlocks uint64 `json:"interval_blocks,omitempty"`
	// Grid-specific params
	GridLevels uint64 `json:"grid_levels,omitempty"`
}

func (m *MsgCreateStrategy) ProtoMessage()           {}
func (m *MsgCreateStrategy) Reset()                  { *m = MsgCreateStrategy{} }
func (m *MsgCreateStrategy) String() string          { return fmt.Sprintf("create_strategy: template=%s creator=%s", m.TemplateName, m.Creator) }
func (m *MsgCreateStrategy) XXX_MessageName() string { return "syreen.intent.MsgCreateStrategy" }

func (m *MsgCreateStrategy) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}
	if err := ValidateTemplateName(m.TemplateName); err != nil {
		return err
	}
	if err := ValidateRiskLevel(m.RiskLevel); err != nil {
		return err
	}
	if m.PoolID == 0 {
		return fmt.Errorf("pool_id must be greater than 0")
	}
	if m.InputDenom == "" {
		return fmt.Errorf("input_denom cannot be empty")
	}
	if m.OutputDenom == "" {
		return fmt.Errorf("output_denom cannot be empty")
	}
	if !m.TotalBudget.IsPositive() {
		return fmt.Errorf("total_budget must be positive")
	}
	if m.ExpiryBlocks == 0 {
		return fmt.Errorf("expiry_blocks must be greater than 0")
	}
	// Template-specific validation
	if m.TemplateName == StrategySafeAccumulate {
		if m.NumExecutions == 0 {
			return fmt.Errorf("num_executions required for safe_accumulate strategy")
		}
		if m.IntervalBlocks == 0 {
			return fmt.Errorf("interval_blocks required for safe_accumulate strategy")
		}
	}
	if m.TemplateName == StrategyGridTrading {
		if m.GridLevels < 2 {
			return fmt.Errorf("grid_levels must be at least 2 for grid_trading strategy")
		}
		if m.GridLevels > 5 {
			return fmt.Errorf("grid_levels must be at most 5 for grid_trading strategy (max 10 chain steps)")
		}
	}
	return nil
}

// MsgCancelStrategy cancels an active strategy
type MsgCancelStrategy struct {
	Creator    string `json:"creator"`
	StrategyID string `json:"strategy_id"`
}

func (m *MsgCancelStrategy) ProtoMessage()           {}
func (m *MsgCancelStrategy) Reset()                  { *m = MsgCancelStrategy{} }
func (m *MsgCancelStrategy) String() string          { return fmt.Sprintf("cancel_strategy: id=%s", m.StrategyID) }
func (m *MsgCancelStrategy) XXX_MessageName() string { return "syreen.intent.MsgCancelStrategy" }

func (m *MsgCancelStrategy) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}
	if m.StrategyID == "" {
		return fmt.Errorf("strategy_id cannot be empty")
	}
	return nil
}

// Response types

type MsgCreateStrategyResponse struct {
	StrategyID string `json:"strategy_id"`
	ChainID    string `json:"chain_id"`
}

func (m *MsgCreateStrategyResponse) ProtoMessage()           {}
func (m *MsgCreateStrategyResponse) Reset()                  { *m = MsgCreateStrategyResponse{} }
func (m *MsgCreateStrategyResponse) String() string          { return fmt.Sprintf("strategy_id=%s chain_id=%s", m.StrategyID, m.ChainID) }
func (m *MsgCreateStrategyResponse) XXX_MessageName() string { return "syreen.intent.MsgCreateStrategyResponse" }

type MsgCancelStrategyResponse struct{}

func (m *MsgCancelStrategyResponse) ProtoMessage()           {}
func (m *MsgCancelStrategyResponse) Reset()                  { *m = MsgCancelStrategyResponse{} }
func (m *MsgCancelStrategyResponse) String() string          { return "cancel_strategy_response" }
func (m *MsgCancelStrategyResponse) XXX_MessageName() string { return "syreen.intent.MsgCancelStrategyResponse" }

// ---------------------------------------------------------------------------
// Query types
// ---------------------------------------------------------------------------

type QueryStrategyTemplatesRequest struct{}

func (m *QueryStrategyTemplatesRequest) ProtoMessage()           {}
func (m *QueryStrategyTemplatesRequest) Reset()                  { *m = QueryStrategyTemplatesRequest{} }
func (m *QueryStrategyTemplatesRequest) String() string          { return "query_strategy_templates" }
func (m *QueryStrategyTemplatesRequest) XXX_MessageName() string { return "syreen.intent.QueryStrategyTemplatesRequest" }

type QueryStrategyTemplatesResponse struct {
	Templates []StrategyTemplate `json:"templates"`
}

func (m *QueryStrategyTemplatesResponse) ProtoMessage()           {}
func (m *QueryStrategyTemplatesResponse) Reset()                  { *m = QueryStrategyTemplatesResponse{} }
func (m *QueryStrategyTemplatesResponse) String() string          { return fmt.Sprintf("templates: %d", len(m.Templates)) }
func (m *QueryStrategyTemplatesResponse) XXX_MessageName() string { return "syreen.intent.QueryStrategyTemplatesResponse" }

// Hand-rolled proto codec: without it the reflection codec drops Templates over
// gRPC/REST (see service.go for the full rationale). Encoded as repeated JSON
// bytes in field #1.
func (m *QueryStrategyTemplatesResponse) Marshal() ([]byte, error) {
	var b []byte
	var err error
	for i := range m.Templates {
		if b, err = appendJSON(b, 1, &m.Templates[i]); err != nil {
			return nil, err
		}
	}
	return b, nil
}
func (m *QueryStrategyTemplatesResponse) MarshalToSizedBuffer(d []byte) (int, error) {
	return sizedBuffer(d, m.Marshal)
}
func (m *QueryStrategyTemplatesResponse) Size() int { b, _ := m.Marshal(); return len(b) }
func (m *QueryStrategyTemplatesResponse) Unmarshal(data []byte) error {
	m.Templates = nil
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return protowire.ParseError(n)
		}
		data = data[n:]
		if num == 1 && typ == protowire.BytesType {
			v, vn := protowire.ConsumeBytes(data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
			var t StrategyTemplate
			if err := json.Unmarshal(v, &t); err != nil {
				return err
			}
			m.Templates = append(m.Templates, t)
		} else {
			vn := protowire.ConsumeFieldValue(num, typ, data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
		}
	}
	return nil
}

type QueryStrategiesRequest struct {
	Address string `json:"address"`
}

func (m *QueryStrategiesRequest) ProtoMessage()           {}
func (m *QueryStrategiesRequest) Reset()                  { *m = QueryStrategiesRequest{} }
func (m *QueryStrategiesRequest) String() string          { return fmt.Sprintf("query_strategies: addr=%s", m.Address) }
func (m *QueryStrategiesRequest) XXX_MessageName() string { return "syreen.intent.QueryStrategiesRequest" }

func (m *QueryStrategiesRequest) Marshal() ([]byte, error)              { return appendString(nil, 1, m.Address), nil }
func (m *QueryStrategiesRequest) MarshalToSizedBuffer(d []byte) (int, error) { return sizedBuffer(d, m.Marshal) }
func (m *QueryStrategiesRequest) Size() int                            { b, _ := m.Marshal(); return len(b) }
func (m *QueryStrategiesRequest) Unmarshal(data []byte) error {
	*m = QueryStrategiesRequest{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return protowire.ParseError(n)
		}
		data = data[n:]
		if num == 1 && typ == protowire.BytesType {
			v, vn := protowire.ConsumeString(data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
			m.Address = v
		} else {
			vn := protowire.ConsumeFieldValue(num, typ, data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
		}
	}
	return nil
}

type QueryStrategiesResponse struct {
	Strategies []Strategy `json:"strategies"`
}

func (m *QueryStrategiesResponse) ProtoMessage()           {}
func (m *QueryStrategiesResponse) Reset()                  { *m = QueryStrategiesResponse{} }
func (m *QueryStrategiesResponse) String() string          { return fmt.Sprintf("strategies: %d", len(m.Strategies)) }
func (m *QueryStrategiesResponse) XXX_MessageName() string { return "syreen.intent.QueryStrategiesResponse" }

func (m *QueryStrategiesResponse) Marshal() ([]byte, error) {
	var b []byte
	var err error
	for i := range m.Strategies {
		if b, err = appendJSON(b, 1, &m.Strategies[i]); err != nil {
			return nil, err
		}
	}
	return b, nil
}
func (m *QueryStrategiesResponse) MarshalToSizedBuffer(d []byte) (int, error) { return sizedBuffer(d, m.Marshal) }
func (m *QueryStrategiesResponse) Size() int                            { b, _ := m.Marshal(); return len(b) }
func (m *QueryStrategiesResponse) Unmarshal(data []byte) error {
	m.Strategies = nil
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return protowire.ParseError(n)
		}
		data = data[n:]
		if num == 1 && typ == protowire.BytesType {
			v, vn := protowire.ConsumeBytes(data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
			var s Strategy
			if err := json.Unmarshal(v, &s); err != nil {
				return err
			}
			m.Strategies = append(m.Strategies, s)
		} else {
			vn := protowire.ConsumeFieldValue(num, typ, data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
		}
	}
	return nil
}

type QueryStrategyRequest struct {
	StrategyID string `json:"strategy_id"`
}

func (m *QueryStrategyRequest) ProtoMessage()           {}
func (m *QueryStrategyRequest) Reset()                  { *m = QueryStrategyRequest{} }
func (m *QueryStrategyRequest) String() string          { return fmt.Sprintf("query_strategy: id=%s", m.StrategyID) }
func (m *QueryStrategyRequest) XXX_MessageName() string { return "syreen.intent.QueryStrategyRequest" }

func (m *QueryStrategyRequest) Marshal() ([]byte, error)              { return appendString(nil, 1, m.StrategyID), nil }
func (m *QueryStrategyRequest) MarshalToSizedBuffer(d []byte) (int, error) { return sizedBuffer(d, m.Marshal) }
func (m *QueryStrategyRequest) Size() int                            { b, _ := m.Marshal(); return len(b) }
func (m *QueryStrategyRequest) Unmarshal(data []byte) error {
	*m = QueryStrategyRequest{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return protowire.ParseError(n)
		}
		data = data[n:]
		if num == 1 && typ == protowire.BytesType {
			v, vn := protowire.ConsumeString(data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
			m.StrategyID = v
		} else {
			vn := protowire.ConsumeFieldValue(num, typ, data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
		}
	}
	return nil
}

type QueryStrategyResponse struct {
	Strategy Strategy         `json:"strategy"`
	Chain    IntentChain      `json:"chain"`
	Found    bool             `json:"found"`
}

func (m *QueryStrategyResponse) ProtoMessage()           {}
func (m *QueryStrategyResponse) Reset()                  { *m = QueryStrategyResponse{} }
func (m *QueryStrategyResponse) String() string          { return fmt.Sprintf("strategy: %+v", m.Strategy) }
func (m *QueryStrategyResponse) XXX_MessageName() string { return "syreen.intent.QueryStrategyResponse" }

func (m *QueryStrategyResponse) Marshal() ([]byte, error) {
	b, err := appendJSON(nil, 1, &m.Strategy)
	if err != nil {
		return nil, err
	}
	if b, err = appendJSON(b, 2, &m.Chain); err != nil {
		return nil, err
	}
	return appendBool(b, 3, m.Found), nil
}
func (m *QueryStrategyResponse) MarshalToSizedBuffer(d []byte) (int, error) { return sizedBuffer(d, m.Marshal) }
func (m *QueryStrategyResponse) Size() int                            { b, _ := m.Marshal(); return len(b) }
func (m *QueryStrategyResponse) Unmarshal(data []byte) error {
	*m = QueryStrategyResponse{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return protowire.ParseError(n)
		}
		data = data[n:]
		switch {
		case num == 1 && typ == protowire.BytesType:
			v, vn := protowire.ConsumeBytes(data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
			if err := json.Unmarshal(v, &m.Strategy); err != nil {
				return err
			}
		case num == 2 && typ == protowire.BytesType:
			v, vn := protowire.ConsumeBytes(data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
			if err := json.Unmarshal(v, &m.Chain); err != nil {
				return err
			}
		case num == 3 && typ == protowire.VarintType:
			v, vn := protowire.ConsumeVarint(data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
			m.Found = v != 0
		default:
			vn := protowire.ConsumeFieldValue(num, typ, data)
			if vn < 0 {
				return protowire.ParseError(vn)
			}
			data = data[vn:]
		}
	}
	return nil
}

// Error codes for strategies
var (
	ErrStrategyNotFound    = fmt.Errorf("strategy not found")
	ErrStrategyNotCreator  = fmt.Errorf("only strategy creator can cancel")
	ErrStrategyNotActive   = fmt.Errorf("strategy is not active")
	ErrStrategyPriceFetch  = fmt.Errorf("failed to fetch current price for strategy")
)
