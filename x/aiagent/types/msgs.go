package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgCreateAgent deploys a new AI agent with initial funding
type MsgCreateAgent struct {
	Owner        string       `json:"owner"`
	Name         string       `json:"name"`
	StrategyType StrategyType `json:"strategy_type"`
	Config       AgentConfig  `json:"config"`
	InitialFunds math.Int     `json:"initial_funds"` // amount in config.InputDenom
}

func (m *MsgCreateAgent) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return err
	}
	if m.Name == "" {
		return fmt.Errorf("agent name cannot be empty")
	}
	if m.InitialFunds.IsNil() || !m.InitialFunds.IsPositive() {
		return ErrInvalidAmount
	}
	// Validate agent config to prevent misconfigured agents from wasting gas every block.
	cfg := m.Config
	if cfg.InputDenom == "" || cfg.OutputDenom == "" {
		return fmt.Errorf("input_denom and output_denom must be set")
	}
	if cfg.PoolID == 0 {
		return fmt.Errorf("pool_id must be positive")
	}
	if cfg.IntervalBlocks <= 0 {
		return fmt.Errorf("interval_blocks must be positive")
	}
	if !cfg.TradePercent.IsNil() && (cfg.TradePercent.IsNegative() || cfg.TradePercent.GT(math.LegacyOneDec())) {
		return fmt.Errorf("trade_percent must be between 0 and 1")
	}
	if m.StrategyType == StrategyGrid && cfg.GridLevels <= 0 {
		return fmt.Errorf("grid_levels must be positive for grid strategy")
	}
	return nil
}
func (m *MsgCreateAgent) ProtoMessage()           {}
func (m *MsgCreateAgent) Reset()                  { *m = MsgCreateAgent{} }
func (m *MsgCreateAgent) String() string          { return fmt.Sprintf("MsgCreateAgent{%s}", m.Owner) }
func (m *MsgCreateAgent) XXX_MessageName() string { return "syreen.aiagent.MsgCreateAgent" }

type MsgCreateAgentResponse struct {
	AgentID      uint64 `json:"agent_id"`
	AgentAddress string `json:"agent_address"`
}

func (m *MsgCreateAgentResponse) ProtoMessage()           {}
func (m *MsgCreateAgentResponse) Reset()                  { *m = MsgCreateAgentResponse{} }
func (m *MsgCreateAgentResponse) String() string          { return "MsgCreateAgentResponse" }
func (m *MsgCreateAgentResponse) XXX_MessageName() string { return "syreen.aiagent.MsgCreateAgentResponse" }

// MsgFundAgent adds funds to an existing agent
type MsgFundAgent struct {
	Owner   string   `json:"owner"`
	AgentID uint64   `json:"agent_id"`
	Amount  math.Int `json:"amount"`
}

func (m *MsgFundAgent) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return err
	}
	if m.Amount.IsNil() || !m.Amount.IsPositive() {
		return ErrInvalidAmount
	}
	return nil
}
func (m *MsgFundAgent) ProtoMessage()           {}
func (m *MsgFundAgent) Reset()                  { *m = MsgFundAgent{} }
func (m *MsgFundAgent) String() string          { return fmt.Sprintf("MsgFundAgent{%s,%d}", m.Owner, m.AgentID) }
func (m *MsgFundAgent) XXX_MessageName() string { return "syreen.aiagent.MsgFundAgent" }

type MsgFundAgentResponse struct{}

func (m *MsgFundAgentResponse) ProtoMessage()           {}
func (m *MsgFundAgentResponse) Reset()                  { *m = MsgFundAgentResponse{} }
func (m *MsgFundAgentResponse) String() string          { return "MsgFundAgentResponse" }
func (m *MsgFundAgentResponse) XXX_MessageName() string { return "syreen.aiagent.MsgFundAgentResponse" }

// MsgWithdrawAgentFunds withdraws funds from agent back to owner
type MsgWithdrawAgentFunds struct {
	Owner   string   `json:"owner"`
	AgentID uint64   `json:"agent_id"`
	Amount  math.Int `json:"amount"`
}

func (m *MsgWithdrawAgentFunds) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return err
	}
	if m.Amount.IsNil() || !m.Amount.IsPositive() {
		return ErrInvalidAmount
	}
	return nil
}
func (m *MsgWithdrawAgentFunds) ProtoMessage()           {}
func (m *MsgWithdrawAgentFunds) Reset()                  { *m = MsgWithdrawAgentFunds{} }
func (m *MsgWithdrawAgentFunds) String() string          { return fmt.Sprintf("MsgWithdrawAgentFunds{%s,%d}", m.Owner, m.AgentID) }
func (m *MsgWithdrawAgentFunds) XXX_MessageName() string { return "syreen.aiagent.MsgWithdrawAgentFunds" }

type MsgWithdrawAgentFundsResponse struct {
	AmountReturned math.Int `json:"amount_returned"`
}

func (m *MsgWithdrawAgentFundsResponse) ProtoMessage()           {}
func (m *MsgWithdrawAgentFundsResponse) Reset()                  { *m = MsgWithdrawAgentFundsResponse{} }
func (m *MsgWithdrawAgentFundsResponse) String() string          { return "MsgWithdrawAgentFundsResponse" }
func (m *MsgWithdrawAgentFundsResponse) XXX_MessageName() string { return "syreen.aiagent.MsgWithdrawAgentFundsResponse" }

// MsgPauseAgent pauses an active agent
type MsgPauseAgent struct {
	Owner   string `json:"owner"`
	AgentID uint64 `json:"agent_id"`
}

func (m *MsgPauseAgent) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Owner)
	return err
}
func (m *MsgPauseAgent) ProtoMessage()           {}
func (m *MsgPauseAgent) Reset()                  { *m = MsgPauseAgent{} }
func (m *MsgPauseAgent) String() string          { return fmt.Sprintf("MsgPauseAgent{%s,%d}", m.Owner, m.AgentID) }
func (m *MsgPauseAgent) XXX_MessageName() string { return "syreen.aiagent.MsgPauseAgent" }

type MsgPauseAgentResponse struct{}

func (m *MsgPauseAgentResponse) ProtoMessage()           {}
func (m *MsgPauseAgentResponse) Reset()                  { *m = MsgPauseAgentResponse{} }
func (m *MsgPauseAgentResponse) String() string          { return "MsgPauseAgentResponse" }
func (m *MsgPauseAgentResponse) XXX_MessageName() string { return "syreen.aiagent.MsgPauseAgentResponse" }

// MsgResumeAgent resumes a paused agent
type MsgResumeAgent struct {
	Owner   string `json:"owner"`
	AgentID uint64 `json:"agent_id"`
}

func (m *MsgResumeAgent) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Owner)
	return err
}
func (m *MsgResumeAgent) ProtoMessage()           {}
func (m *MsgResumeAgent) Reset()                  { *m = MsgResumeAgent{} }
func (m *MsgResumeAgent) String() string          { return fmt.Sprintf("MsgResumeAgent{%s,%d}", m.Owner, m.AgentID) }
func (m *MsgResumeAgent) XXX_MessageName() string { return "syreen.aiagent.MsgResumeAgent" }

type MsgResumeAgentResponse struct{}

func (m *MsgResumeAgentResponse) ProtoMessage()           {}
func (m *MsgResumeAgentResponse) Reset()                  { *m = MsgResumeAgentResponse{} }
func (m *MsgResumeAgentResponse) String() string          { return "MsgResumeAgentResponse" }
func (m *MsgResumeAgentResponse) XXX_MessageName() string { return "syreen.aiagent.MsgResumeAgentResponse" }

// MsgUpdateAgentStrategy updates agent config without redeploying
type MsgUpdateAgentStrategy struct {
	Owner   string      `json:"owner"`
	AgentID uint64      `json:"agent_id"`
	Config  AgentConfig `json:"config"`
}

func (m *MsgUpdateAgentStrategy) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return err
	}
	cfg := m.Config
	if cfg.InputDenom == "" || cfg.OutputDenom == "" {
		return fmt.Errorf("input_denom and output_denom must be set")
	}
	if cfg.PoolID == 0 {
		return fmt.Errorf("pool_id must be positive")
	}
	if cfg.IntervalBlocks <= 0 {
		return fmt.Errorf("interval_blocks must be positive")
	}
	if !cfg.TradePercent.IsNil() && (cfg.TradePercent.IsNegative() || cfg.TradePercent.GT(math.LegacyOneDec())) {
		return fmt.Errorf("trade_percent must be between 0 and 1")
	}
	return nil
}
func (m *MsgUpdateAgentStrategy) ProtoMessage()           {}
func (m *MsgUpdateAgentStrategy) Reset()                  { *m = MsgUpdateAgentStrategy{} }
func (m *MsgUpdateAgentStrategy) String() string          { return fmt.Sprintf("MsgUpdateAgentStrategy{%s,%d}", m.Owner, m.AgentID) }
func (m *MsgUpdateAgentStrategy) XXX_MessageName() string { return "syreen.aiagent.MsgUpdateAgentStrategy" }

type MsgUpdateAgentStrategyResponse struct{}

func (m *MsgUpdateAgentStrategyResponse) ProtoMessage()           {}
func (m *MsgUpdateAgentStrategyResponse) Reset()                  { *m = MsgUpdateAgentStrategyResponse{} }
func (m *MsgUpdateAgentStrategyResponse) String() string          { return "MsgUpdateAgentStrategyResponse" }
func (m *MsgUpdateAgentStrategyResponse) XXX_MessageName() string { return "syreen.aiagent.MsgUpdateAgentStrategyResponse" }
