package types

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// IntentStatus represents the lifecycle status of an intent
type IntentStatus string

const (
	StatusPending   IntentStatus = "pending"
	StatusSolving   IntentStatus = "solving"
	StatusFulfilled IntentStatus = "fulfilled"
	StatusExpired   IntentStatus = "expired"
	StatusFailed    IntentStatus = "failed"
)

// Intent type constants
const (
	IntentTypeSwap     = "swap"
	IntentTypeTransfer = "transfer"
	IntentTypeDeFi     = "defi"
	IntentTypeCustom   = "custom"

	// Trading intent types
	IntentTypeLimitBuy   = "limit_buy"
	IntentTypeLimitSell  = "limit_sell"
	IntentTypeStopLoss   = "stop_loss"
	IntentTypeTakeProfit = "take_profit"
	IntentTypeTWAP           = "twap"
	IntentTypeDCA            = "dca"
	IntentTypeCrossChainSwap = "cross_chain_swap"
)

// Intent represents a user's declared intent (outcome they want)
type Intent struct {
	ID         string          `json:"id"`
	Creator    string          `json:"creator"`
	IntentType string          `json:"intent_type"`
	Body       json.RawMessage `json:"body"`
	MaxFee     sdk.Coins       `json:"max_fee"`
	Tip        sdk.Coins       `json:"tip"`
	Expiry     int64           `json:"expiry"`
	Status     IntentStatus    `json:"status"`
	CreatedAt  int64           `json:"created_at"`
	SolverAddr string          `json:"solver_addr"`
}

func (m *Intent) ProtoMessage()           {}
func (m *Intent) Reset()                  { *m = Intent{} }
func (m *Intent) String() string          { return fmt.Sprintf("intent(%s): type=%s creator=%s status=%s", m.ID, m.IntentType, m.Creator, m.Status) }
func (m *Intent) XXX_MessageName() string { return "syreen.intent.Intent" }

// SwapIntent describes a token swap intent
type SwapIntent struct {
	InputDenom      string        `json:"input_denom"`
	InputAmount     math.Int      `json:"input_amount"`
	OutputDenom     string        `json:"output_denom"`
	MinOutputAmount math.Int      `json:"min_output_amount"`
	MaxSlippage     math.LegacyDec `json:"max_slippage"`
}

func (m *SwapIntent) ProtoMessage()           {}
func (m *SwapIntent) Reset()                  { *m = SwapIntent{} }
func (m *SwapIntent) String() string          { return fmt.Sprintf("swap: %s %s -> %s (min %s)", m.InputAmount, m.InputDenom, m.OutputDenom, m.MinOutputAmount) }
func (m *SwapIntent) XXX_MessageName() string { return "syreen.intent.SwapIntent" }

// TransferIntent describes a token transfer intent (including cross-chain)
type TransferIntent struct {
	Recipient string    `json:"recipient"`
	Amount    sdk.Coins `json:"amount"`
	ChainID   string    `json:"chain_id"`
}

func (m *TransferIntent) ProtoMessage()           {}
func (m *TransferIntent) Reset()                  { *m = TransferIntent{} }
func (m *TransferIntent) String() string          { return fmt.Sprintf("transfer: %s to %s chain=%s", m.Amount, m.Recipient, m.ChainID) }
func (m *TransferIntent) XXX_MessageName() string { return "syreen.intent.TransferIntent" }

// LimitBuyIntent — buy OutputDenom when price drops to or below TargetPrice
type LimitBuyIntent struct {
	InputDenom      string         `json:"input_denom"`
	InputAmount     math.Int       `json:"input_amount"`
	OutputDenom     string         `json:"output_denom"`
	TargetPrice     math.LegacyDec `json:"target_price"`
	PoolID          uint64         `json:"pool_id"`
	MinOutputAmount math.Int       `json:"min_output_amount"`
}

func (m *LimitBuyIntent) ProtoMessage()           {}
func (m *LimitBuyIntent) Reset()                  { *m = LimitBuyIntent{} }
func (m *LimitBuyIntent) String() string {
	return fmt.Sprintf("limit_buy: %s %s -> %s at price %s (min %s, pool %d)", m.InputAmount, m.InputDenom, m.OutputDenom, m.TargetPrice, m.MinOutputAmount, m.PoolID)
}
func (m *LimitBuyIntent) XXX_MessageName() string { return "syreen.intent.LimitBuyIntent" }

// LimitSellIntent — sell InputDenom when price rises to or above TargetPrice
type LimitSellIntent struct {
	InputDenom      string         `json:"input_denom"`
	InputAmount     math.Int       `json:"input_amount"`
	OutputDenom     string         `json:"output_denom"`
	TargetPrice     math.LegacyDec `json:"target_price"`
	PoolID          uint64         `json:"pool_id"`
	MinOutputAmount math.Int       `json:"min_output_amount"`
}

func (m *LimitSellIntent) ProtoMessage()           {}
func (m *LimitSellIntent) Reset()                  { *m = LimitSellIntent{} }
func (m *LimitSellIntent) String() string {
	return fmt.Sprintf("limit_sell: %s %s -> %s at price %s (min %s, pool %d)", m.InputAmount, m.InputDenom, m.OutputDenom, m.TargetPrice, m.MinOutputAmount, m.PoolID)
}
func (m *LimitSellIntent) XXX_MessageName() string { return "syreen.intent.LimitSellIntent" }

// StopLossIntent — sell InputDenom if price drops to or below StopPrice (protective sell)
type StopLossIntent struct {
	InputDenom      string         `json:"input_denom"`
	InputAmount     math.Int       `json:"input_amount"`
	OutputDenom     string         `json:"output_denom"`
	StopPrice       math.LegacyDec `json:"stop_price"`
	PoolID          uint64         `json:"pool_id"`
	MinOutputAmount math.Int       `json:"min_output_amount"`
}

func (m *StopLossIntent) ProtoMessage()           {}
func (m *StopLossIntent) Reset()                  { *m = StopLossIntent{} }
func (m *StopLossIntent) String() string {
	return fmt.Sprintf("stop_loss: %s %s -> %s at stop %s (min %s, pool %d)", m.InputAmount, m.InputDenom, m.OutputDenom, m.StopPrice, m.MinOutputAmount, m.PoolID)
}
func (m *StopLossIntent) XXX_MessageName() string { return "syreen.intent.StopLossIntent" }

// TakeProfitIntent — sell InputDenom if price rises to or above TargetPrice (profit taking)
type TakeProfitIntent struct {
	InputDenom      string         `json:"input_denom"`
	InputAmount     math.Int       `json:"input_amount"`
	OutputDenom     string         `json:"output_denom"`
	TargetPrice     math.LegacyDec `json:"target_price"`
	PoolID          uint64         `json:"pool_id"`
	MinOutputAmount math.Int       `json:"min_output_amount"`
}

func (m *TakeProfitIntent) ProtoMessage()           {}
func (m *TakeProfitIntent) Reset()                  { *m = TakeProfitIntent{} }
func (m *TakeProfitIntent) String() string {
	return fmt.Sprintf("take_profit: %s %s -> %s at price %s (min %s, pool %d)", m.InputAmount, m.InputDenom, m.OutputDenom, m.TargetPrice, m.MinOutputAmount, m.PoolID)
}
func (m *TakeProfitIntent) XXX_MessageName() string { return "syreen.intent.TakeProfitIntent" }

// DCAIntent — dollar cost average: split TotalAmount into NumExecutions equal swaps over interval blocks
type DCAIntent struct {
	InputDenom     string   `json:"input_denom"`
	TotalAmount    math.Int `json:"total_amount"`
	OutputDenom    string   `json:"output_denom"`
	NumExecutions  uint64   `json:"num_executions"`
	IntervalBlocks uint64   `json:"interval_blocks"`
	PoolID         uint64   `json:"pool_id"`
	ExecutedCount  uint64   `json:"executed_count"`
	NextExecBlock  int64    `json:"next_exec_block"`
}

func (m *DCAIntent) ProtoMessage()           {}
func (m *DCAIntent) Reset()                  { *m = DCAIntent{} }
func (m *DCAIntent) String() string {
	return fmt.Sprintf("dca: %s %s -> %s (%d/%d executions, interval %d, pool %d)", m.TotalAmount, m.InputDenom, m.OutputDenom, m.ExecutedCount, m.NumExecutions, m.IntervalBlocks, m.PoolID)
}
func (m *DCAIntent) XXX_MessageName() string { return "syreen.intent.DCAIntent" }

// CrossChainSwapIntent — swap locally then IBC transfer the output to another chain
type CrossChainSwapIntent struct {
	InputDenom      string         `json:"input_denom"`
	InputAmount     math.Int       `json:"input_amount"`
	OutputDenom     string         `json:"output_denom"`
	MinOutputAmount math.Int       `json:"min_output_amount"`
	PoolID          uint64         `json:"pool_id"`

	// IBC transfer parameters (post-swap)
	IBCSourcePort    string `json:"ibc_source_port"`
	IBCSourceChannel string `json:"ibc_source_channel"`
	Receiver         string `json:"receiver"`
	TimeoutBlocks    uint64 `json:"timeout_blocks"`

	// Tracking
	SwapDone bool `json:"swap_done"`
	IBCSent  bool `json:"ibc_sent"`
}

func (m *CrossChainSwapIntent) ProtoMessage()           {}
func (m *CrossChainSwapIntent) Reset()                  { *m = CrossChainSwapIntent{} }
func (m *CrossChainSwapIntent) String() string {
	return fmt.Sprintf("cross_chain_swap: %s %s -> %s via %s/%s to %s", m.InputAmount, m.InputDenom, m.OutputDenom, m.IBCSourcePort, m.IBCSourceChannel, m.Receiver)
}
func (m *CrossChainSwapIntent) XXX_MessageName() string { return "syreen.intent.CrossChainSwapIntent" }

// Solver represents a registered intent solver
type Solver struct {
	Address         string   `json:"address"`
	Moniker         string   `json:"moniker"`
	StakedAmount    sdk.Coin `json:"staked_amount"`
	ReputationScore uint64   `json:"reputation_score"`
	TotalSolved     uint64   `json:"total_solved"`
	TotalFailed     uint64   `json:"total_failed"`
	Active          bool     `json:"active"`
	JoinedAt        int64    `json:"joined_at"`
	UnbondingHeight int64    `json:"unbonding_height,omitempty"` // Block height when unbonding completes
}

func (m *Solver) ProtoMessage()           {}
func (m *Solver) Reset()                  { *m = Solver{} }
func (m *Solver) String() string          { return fmt.Sprintf("solver(%s): moniker=%s active=%v reputation=%d", m.Address, m.Moniker, m.Active, m.ReputationScore) }
func (m *Solver) XXX_MessageName() string { return "syreen.intent.Solver" }

// Solution represents a solver's proposed solution for an intent
type Solution struct {
	IntentID        string            `json:"intent_id"`
	SolverAddr      string            `json:"solver_addr"`
	ExecutionMsgs   []json.RawMessage `json:"execution_msgs"`
	ExpectedOutcome json.RawMessage   `json:"expected_outcome"`
	GasEstimate     uint64            `json:"gas_estimate"`
	Tip             sdk.Coins         `json:"tip,omitempty"`
	SubmittedAt     int64             `json:"submitted_at"`
}

// SolutionOutcome is the solver's declared execution result, carried in
// Solution.ExpectedOutcome. For swap intents, OutputAmount is the amount of the
// intent's OutputDenom the solver commits to deliver to the creator.
//
// Fairness Engine (best execution): the auction ranks solvers by this declared
// output — highest wins — and the winner is HELD to it at settlement
// (verifySwapOutcome requires delivered >= declared). So over-declaring to win
// the auction fails and slashes the solver; competition drives the surplus to
// the user instead of the solver.
type SolutionOutcome struct {
	OutputAmount math.Int `json:"output_amount"`
}

// DeclaredOutput returns the output amount the solver committed to in
// ExpectedOutcome, and whether a valid, positive amount was declared.
func (m *Solution) DeclaredOutput() (math.Int, bool) {
	if len(m.ExpectedOutcome) == 0 {
		return math.ZeroInt(), false
	}
	var o SolutionOutcome
	if err := json.Unmarshal(m.ExpectedOutcome, &o); err != nil || o.OutputAmount.IsNil() || !o.OutputAmount.IsPositive() {
		return math.ZeroInt(), false
	}
	return o.OutputAmount, true
}

func (m *Solution) ProtoMessage()           {}
func (m *Solution) Reset()                  { *m = Solution{} }
func (m *Solution) String() string          { return fmt.Sprintf("solution: intent=%s solver=%s gas=%d", m.IntentID, m.SolverAddr, m.GasEstimate) }
func (m *Solution) XXX_MessageName() string { return "syreen.intent.Solution" }
