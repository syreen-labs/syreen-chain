package types

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// IntentChainStatus represents the lifecycle status of an intent chain
type IntentChainStatus string

const (
	ChainStatusActive    IntentChainStatus = "active"
	ChainStatusCompleted IntentChainStatus = "completed"
	ChainStatusFailed    IntentChainStatus = "failed"
	ChainStatusExpired   IntentChainStatus = "expired"
	ChainStatusCancelled IntentChainStatus = "cancelled"
)

// ConditionType defines when a chain step should activate
type ConditionType string

const (
	ConditionImmediate  ConditionType = "immediate"    // activate as soon as previous step completes
	ConditionPriceAbove ConditionType = "price_above"  // activate when price rises above target
	ConditionPriceBelow ConditionType = "price_below"  // activate when price drops below target
	ConditionDelay      ConditionType = "delay"        // activate after DelayBlocks blocks

	// AI-reactive conditions — evaluated each block against live dex signal/sentiment state.
	ConditionRSIAbove        ConditionType = "rsi_above"          // fires when pool RSI >= RSIThreshold
	ConditionRSIBelow        ConditionType = "rsi_below"          // fires when pool RSI <= RSIThreshold
	ConditionFearGreedAbove  ConditionType = "fear_greed_above"   // fires when F&G index >= FearGreedThreshold
	ConditionFearGreedBelow  ConditionType = "fear_greed_below"   // fires when F&G index <= FearGreedThreshold
	ConditionSignalEquals    ConditionType = "signal_equals"      // fires when signal label equals SignalLabel
	ConditionVolatilityAbove ConditionType = "volatility_above"   // fires when volatility >= VolatilityThreshold
)

// Standard AI thresholds — used by strategy templates. Declared as named constants
// so they can be tuned in one place as the signal model evolves.
var (
	RSIOversoldThreshold   = math.LegacyNewDec(30) // RSI < 30 → oversold (buy)
	RSIOverboughtThreshold = math.LegacyNewDec(70) // RSI > 70 → overbought (sell)
)

const (
	FearGreedExtremeFear  int64 = 25 // index <= 25 → extreme fear (buy the dip)
	FearGreedFear         int64 = 40 // index <= 40 → fear
	FearGreedNeutralUpper int64 = 50 // index >= 50 → neutral or bullish
	FearGreedExtremeGreed int64 = 75 // index >= 75 → extreme greed (take profit)
)

// VolatilityBreakoutThreshold is the minimum volatility score required for a
// volatility-breakout template to fire (0.05 = 5% range).
var VolatilityBreakoutThreshold = math.LegacyMustNewDecFromStr("0.05")

// InputSourceType controls where a step gets its input tokens
type InputSourceType string

const (
	InputFromPrevious InputSourceType = "previous" // use full output of previous step
	InputFromPortion  InputSourceType = "portion"  // use PortionPercent of previous step output
	InputFromFixed    InputSourceType = "fixed"    // use a specific locked amount
)

// IntentChain represents a multi-step conditional trading strategy
type IntentChain struct {
	ID          string            `json:"id"`
	Creator     string            `json:"creator"`
	Steps       []ChainStep       `json:"steps"`
	CurrentIdx  int               `json:"current_idx"`
	Status      IntentChainStatus `json:"status"`
	CreatedAt   int64             `json:"created_at"`
	Expiry      int64             `json:"expiry"`
	MaxFee      sdk.Coins         `json:"max_fee"`
	Tip         sdk.Coins         `json:"tip"`
}

func (m *IntentChain) ProtoMessage()           {}
func (m *IntentChain) Reset()                  { *m = IntentChain{} }
func (m *IntentChain) String() string          { return fmt.Sprintf("chain(%s): steps=%d current=%d status=%s", m.ID, len(m.Steps), m.CurrentIdx, m.Status) }
func (m *IntentChain) XXX_MessageName() string { return "syreen.intent.IntentChain" }

// ChainStep represents one step in a conditional intent chain
type ChainStep struct {
	IntentType    string          `json:"intent_type"`
	Body          json.RawMessage `json:"body"`
	Condition     ChainCondition  `json:"condition"`
	InputSource   InputSource     `json:"input_source"`
	IntentID      string          `json:"intent_id"`
	Status        IntentStatus    `json:"status"`
	OutputCoins   sdk.Coins       `json:"output_coins"`
	CompletedAt   int64           `json:"completed_at"`
}

// ChainCondition defines when a chain step should activate
type ChainCondition struct {
	Type           ConditionType  `json:"type"`
	PriceTarget    math.LegacyDec `json:"price_target,omitempty"`
	PriceDenom     string         `json:"price_denom,omitempty"`
	PriceBaseDenom string         `json:"price_base_denom,omitempty"`
	PoolID         uint64         `json:"pool_id,omitempty"`
	DelayBlocks    uint64         `json:"delay_blocks,omitempty"`

	// AI-reactive condition parameters
	RSIThreshold        math.LegacyDec `json:"rsi_threshold,omitempty"`         // 0-100
	FearGreedThreshold  int64          `json:"fear_greed_threshold,omitempty"`  // 0-100
	SignalLabel         string         `json:"signal_label,omitempty"`          // e.g. "STRONG_BUY"
	VolatilityThreshold math.LegacyDec `json:"volatility_threshold,omitempty"`  // positive dec
}

// InputSource controls where a step gets its input tokens
type InputSource struct {
	Type           InputSourceType `json:"type"`
	PortionPercent math.LegacyDec `json:"portion_percent,omitempty"` // 0-1, e.g., 0.5 = 50%
	Amount         math.Int        `json:"amount,omitempty"`
	Denom          string          `json:"denom,omitempty"`
}

// ValidateChainStep validates a single chain step
func ValidateChainStep(step ChainStep, idx int, maxSteps int) error {
	// Validate intent type is a trading type
	switch step.IntentType {
	case IntentTypeLimitBuy, IntentTypeLimitSell, IntentTypeStopLoss,
		IntentTypeTakeProfit, IntentTypeDCA, IntentTypeTWAP, IntentTypeSwap,
		IntentTypeCrossChainSwap:
		// valid
	default:
		return fmt.Errorf("step %d: invalid intent type %s for chain", idx, step.IntentType)
	}

	if len(step.Body) == 0 || !json.Valid(step.Body) {
		return fmt.Errorf("step %d: invalid body JSON", idx)
	}

	// First step must be immediate
	if idx == 0 && step.Condition.Type != ConditionImmediate {
		return fmt.Errorf("step 0 must have condition type 'immediate'")
	}

	// Validate condition
	switch step.Condition.Type {
	case ConditionImmediate:
		// no extra validation
	case ConditionPriceAbove, ConditionPriceBelow:
		if step.Condition.PriceTarget.IsNil() || !step.Condition.PriceTarget.IsPositive() {
			return fmt.Errorf("step %d: price condition requires positive price_target", idx)
		}
		if step.Condition.PoolID == 0 {
			return fmt.Errorf("step %d: price condition requires pool_id", idx)
		}
	case ConditionDelay:
		if step.Condition.DelayBlocks == 0 {
			return fmt.Errorf("step %d: delay condition requires delay_blocks > 0", idx)
		}
	case ConditionRSIAbove, ConditionRSIBelow:
		if step.Condition.PoolID == 0 {
			return fmt.Errorf("step %d: rsi condition requires pool_id", idx)
		}
		if step.Condition.RSIThreshold.IsNil() ||
			step.Condition.RSIThreshold.IsNegative() ||
			step.Condition.RSIThreshold.GT(math.LegacyNewDec(100)) {
			return fmt.Errorf("step %d: rsi_threshold must be between 0 and 100", idx)
		}
	case ConditionFearGreedAbove, ConditionFearGreedBelow:
		if step.Condition.PoolID == 0 {
			return fmt.Errorf("step %d: fear/greed condition requires pool_id", idx)
		}
		if step.Condition.FearGreedThreshold < 0 || step.Condition.FearGreedThreshold > 100 {
			return fmt.Errorf("step %d: fear_greed_threshold must be between 0 and 100", idx)
		}
	case ConditionSignalEquals:
		if step.Condition.PoolID == 0 {
			return fmt.Errorf("step %d: signal condition requires pool_id", idx)
		}
		if step.Condition.SignalLabel == "" {
			return fmt.Errorf("step %d: signal_label cannot be empty", idx)
		}
	case ConditionVolatilityAbove:
		if step.Condition.PoolID == 0 {
			return fmt.Errorf("step %d: volatility condition requires pool_id", idx)
		}
		if step.Condition.VolatilityThreshold.IsNil() || !step.Condition.VolatilityThreshold.IsPositive() {
			return fmt.Errorf("step %d: volatility_threshold must be positive", idx)
		}
	default:
		return fmt.Errorf("step %d: unknown condition type %s", idx, step.Condition.Type)
	}

	// Validate input source
	switch step.InputSource.Type {
	case InputFromPrevious:
		if idx == 0 {
			return fmt.Errorf("step 0 cannot use 'previous' input source")
		}
	case InputFromPortion:
		if idx == 0 {
			return fmt.Errorf("step 0 cannot use 'portion' input source")
		}
		if step.InputSource.PortionPercent.IsNil() || !step.InputSource.PortionPercent.IsPositive() || step.InputSource.PortionPercent.GT(math.LegacyOneDec()) {
			return fmt.Errorf("step %d: portion_percent must be between 0 (exclusive) and 1 (inclusive)", idx)
		}
	case InputFromFixed:
		if step.InputSource.Amount.IsNil() || !step.InputSource.Amount.IsPositive() {
			return fmt.Errorf("step %d: fixed input requires positive amount", idx)
		}
		if step.InputSource.Denom == "" {
			return fmt.Errorf("step %d: fixed input requires denom", idx)
		}
	default:
		if idx == 0 {
			// Step 0 uses its body's input directly, input_source can be empty
		} else {
			return fmt.Errorf("step %d: unknown input source type %s", idx, step.InputSource.Type)
		}
	}

	return nil
}
