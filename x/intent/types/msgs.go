package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgSubmitIntent    = "submit_intent"
	TypeMsgRegisterSolver  = "register_solver"
	TypeMsgDeregisterSolver = "deregister_solver"
	TypeMsgSubmitSolution  = "submit_solution"
	TypeMsgFulfillIntent   = "fulfill_intent"
	TypeMsgCancelIntent    = "cancel_intent"
	TypeMsgSubmitChain     = "submit_chain"
	TypeMsgCancelChain     = "cancel_chain"
)

// Validation bounds enforced in ValidateBasic.
const (
	// MaxExecutionMsgs caps the number of execution messages a single solution
	// may contain. Prevents resource exhaustion / DoS via unbounded arrays.
	MaxExecutionMsgs = 100

	// MaxExpiryBlocks caps how many blocks in the future an intent may expire.
	// At ~1 block/sec this is roughly 27 hours.
	MaxExpiryBlocks = 100_000

	// MaxFeeOrTipAmount is the per-coin sanity cap for MaxFee and Tip amounts
	// (10 trillion usyreen = 10M SYR).
	MaxFeeOrTipAmount = 10_000_000_000_000
)

// --- MsgSubmitIntent ---

// MsgSubmitIntent submits a new intent declaration
type MsgSubmitIntent struct {
	Creator string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator"`
	IntentType string `protobuf:"bytes,2,opt,name=intent_type,proto3" json:"intent_type"`
	Body json.RawMessage `protobuf:"bytes,3,opt,name=body,proto3" json:"body"`
	MaxFee sdk.Coins `protobuf:"bytes,4,opt,name=max_fee,proto3" json:"max_fee"`
	Tip sdk.Coins `protobuf:"bytes,5,opt,name=tip,proto3" json:"tip"`
	ExpiryBlocks uint64 `protobuf:"bytes,6,opt,name=expiry_blocks,proto3" json:"expiry_blocks"`
}

func (m *MsgSubmitIntent) ProtoMessage()           {}
func (m *MsgSubmitIntent) Reset()                  { *m = MsgSubmitIntent{} }
func (m *MsgSubmitIntent) String() string          { return fmt.Sprintf("submit_intent: creator=%s type=%s", m.Creator, m.IntentType) }
func (m *MsgSubmitIntent) XXX_MessageName() string { return "syreen.intent.MsgSubmitIntent" }

func (m *MsgSubmitIntent) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}
	switch m.IntentType {
	case IntentTypeSwap, IntentTypeTransfer, IntentTypeDeFi, IntentTypeCustom,
		IntentTypeLimitBuy, IntentTypeLimitSell, IntentTypeStopLoss, IntentTypeTakeProfit,
		IntentTypeTWAP, IntentTypeDCA, IntentTypeCrossChainSwap:
	default:
		return ErrIntentTypeUnsupported
	}
	if len(m.Body) == 0 {
		return fmt.Errorf("intent body cannot be empty")
	}
	if !m.MaxFee.IsValid() {
		return fmt.Errorf("invalid max fee: %s", m.MaxFee)
	}
	for _, c := range m.MaxFee {
		if c.Amount.IsInt64() && c.Amount.Int64() > MaxFeeOrTipAmount {
			return fmt.Errorf("max fee amount %s exceeds sanity cap %d", c.Amount, MaxFeeOrTipAmount)
		}
		if !c.Amount.IsInt64() {
			return fmt.Errorf("max fee amount %s exceeds sanity cap %d", c.Amount, MaxFeeOrTipAmount)
		}
	}
	if !m.Tip.IsValid() {
		return fmt.Errorf("invalid tip: %s", m.Tip)
	}
	for _, c := range m.Tip {
		if c.Amount.IsInt64() && c.Amount.Int64() > MaxFeeOrTipAmount {
			return fmt.Errorf("tip amount %s exceeds sanity cap %d", c.Amount, MaxFeeOrTipAmount)
		}
		if !c.Amount.IsInt64() {
			return fmt.Errorf("tip amount %s exceeds sanity cap %d", c.Amount, MaxFeeOrTipAmount)
		}
	}
	if m.ExpiryBlocks == 0 {
		return fmt.Errorf("expiry blocks must be greater than 0")
	}
	if m.ExpiryBlocks > MaxExpiryBlocks {
		return fmt.Errorf("expiry blocks %d exceeds maximum %d", m.ExpiryBlocks, MaxExpiryBlocks)
	}
	return nil
}

// --- MsgRegisterSolver ---

// MsgRegisterSolver registers a new solver
type MsgRegisterSolver struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address"`
	Moniker string `protobuf:"bytes,2,opt,name=moniker,proto3" json:"moniker"`
	StakeAmount sdk.Coin `protobuf:"bytes,3,opt,name=stake_amount,proto3" json:"stake_amount"`
}

func (m *MsgRegisterSolver) ProtoMessage()           {}
func (m *MsgRegisterSolver) Reset()                  { *m = MsgRegisterSolver{} }
func (m *MsgRegisterSolver) String() string          { return fmt.Sprintf("register_solver: addr=%s moniker=%s", m.Address, m.Moniker) }
func (m *MsgRegisterSolver) XXX_MessageName() string { return "syreen.intent.MsgRegisterSolver" }

func (m *MsgRegisterSolver) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		return fmt.Errorf("invalid solver address: %w", err)
	}
	if m.Moniker == "" {
		return fmt.Errorf("solver moniker cannot be empty")
	}
	if !m.StakeAmount.IsPositive() {
		return fmt.Errorf("stake amount must be positive: %s", m.StakeAmount)
	}
	return nil
}

// --- MsgDeregisterSolver ---

// MsgDeregisterSolver removes a solver from the registry
type MsgDeregisterSolver struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address"`
}

func (m *MsgDeregisterSolver) ProtoMessage()           {}
func (m *MsgDeregisterSolver) Reset()                  { *m = MsgDeregisterSolver{} }
func (m *MsgDeregisterSolver) String() string          { return fmt.Sprintf("deregister_solver: addr=%s", m.Address) }
func (m *MsgDeregisterSolver) XXX_MessageName() string { return "syreen.intent.MsgDeregisterSolver" }

func (m *MsgDeregisterSolver) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		return fmt.Errorf("invalid solver address: %w", err)
	}
	return nil
}

// --- MsgSubmitSolution ---

// MsgSubmitSolution submits a solution for an intent
type MsgSubmitSolution struct {
	SolverAddr string `protobuf:"bytes,1,opt,name=solver_addr,proto3" json:"solver_addr"`
	IntentID string `protobuf:"bytes,2,opt,name=intent_id,proto3" json:"intent_id"`
	ExecutionMsgs []json.RawMessage `protobuf:"bytes,3,opt,name=execution_msgs,proto3" json:"execution_msgs"`
	ExpectedOutcome json.RawMessage `protobuf:"bytes,4,opt,name=expected_outcome,proto3" json:"expected_outcome"`
}

func (m *MsgSubmitSolution) ProtoMessage()           {}
func (m *MsgSubmitSolution) Reset()                  { *m = MsgSubmitSolution{} }
func (m *MsgSubmitSolution) String() string          { return fmt.Sprintf("submit_solution: solver=%s intent=%s", m.SolverAddr, m.IntentID) }
func (m *MsgSubmitSolution) XXX_MessageName() string { return "syreen.intent.MsgSubmitSolution" }

func (m *MsgSubmitSolution) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.SolverAddr)
	if err != nil {
		return fmt.Errorf("invalid solver address: %w", err)
	}
	if m.IntentID == "" {
		return fmt.Errorf("intent ID cannot be empty")
	}
	if len(m.ExecutionMsgs) == 0 {
		return fmt.Errorf("execution messages cannot be empty")
	}
	if len(m.ExecutionMsgs) > MaxExecutionMsgs {
		return fmt.Errorf("execution messages count %d exceeds maximum %d", len(m.ExecutionMsgs), MaxExecutionMsgs)
	}
	if len(m.ExpectedOutcome) > 0 && !json.Valid(m.ExpectedOutcome) {
		return fmt.Errorf("expected_outcome is not valid JSON")
	}
	return nil
}

// --- MsgFulfillIntent ---

// MsgFulfillIntent triggers execution of the winning solution for an intent.
// When SolverAddr is empty (auto-selection mode), Sender must be the intent
// creator to prevent unauthorized third parties from triggering fulfillment.
type MsgFulfillIntent struct {
	Sender string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	SolverAddr string `protobuf:"bytes,2,opt,name=solver_addr,proto3" json:"solver_addr"`
	IntentID string `protobuf:"bytes,3,opt,name=intent_id,proto3" json:"intent_id"`
}

func (m *MsgFulfillIntent) ProtoMessage()           {}
func (m *MsgFulfillIntent) Reset()                  { *m = MsgFulfillIntent{} }
func (m *MsgFulfillIntent) String() string          { return fmt.Sprintf("fulfill_intent: solver=%s intent=%s", m.SolverAddr, m.IntentID) }
func (m *MsgFulfillIntent) XXX_MessageName() string { return "syreen.intent.MsgFulfillIntent" }

func (m *MsgFulfillIntent) ValidateBasic() error {
	// SolverAddr may be empty to trigger auto-selection via auction
	if m.SolverAddr != "" {
		_, err := sdk.AccAddressFromBech32(m.SolverAddr)
		if err != nil {
			return fmt.Errorf("invalid solver address: %w", err)
		}
	}
	// H-14: When SolverAddr is empty (auto-selection), require a valid Sender
	// so that the keeper can verify the sender is the intent creator.
	if m.SolverAddr == "" {
		if m.Sender == "" {
			return fmt.Errorf("sender is required when solver_addr is empty (auto-selection)")
		}
		if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
			return fmt.Errorf("invalid sender address: %w", err)
		}
	}
	if m.IntentID == "" {
		return fmt.Errorf("intent ID cannot be empty")
	}
	return nil
}

// --- MsgCancelIntent ---

// MsgCancelIntent lets an intent's creator cancel it before hard expiry and
// reclaim any locked funds (MaxFee + Tip and any locked trading input). Only the
// creator may cancel, and only while the intent is still Pending or Solving.
type MsgCancelIntent struct {
	Creator  string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator"`
	IntentID string `protobuf:"bytes,2,opt,name=intent_id,proto3" json:"intent_id"`
}

func (m *MsgCancelIntent) ProtoMessage()           {}
func (m *MsgCancelIntent) Reset()                  { *m = MsgCancelIntent{} }
func (m *MsgCancelIntent) String() string          { return fmt.Sprintf("cancel_intent: creator=%s intent=%s", m.Creator, m.IntentID) }
func (m *MsgCancelIntent) XXX_MessageName() string { return "syreen.intent.MsgCancelIntent" }

func (m *MsgCancelIntent) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}
	if m.IntentID == "" {
		return fmt.Errorf("intent ID cannot be empty")
	}
	return nil
}

type MsgCancelIntentResponse struct{}

func (m *MsgCancelIntentResponse) ProtoMessage()           {}
func (m *MsgCancelIntentResponse) Reset()                  { *m = MsgCancelIntentResponse{} }
func (m *MsgCancelIntentResponse) String() string          { return "cancel_intent_response" }
func (m *MsgCancelIntentResponse) XXX_MessageName() string { return "syreen.intent.MsgCancelIntentResponse" }

// --- MsgUpdateParams ---

// MsgUpdateParams is the standard Cosmos gov param-update message for the intent
// module. Only the gov module authority may execute it. The Params are carried
// on the wire as JSON bytes (see intent_proto.go Marshal/Unmarshal) because
// Params is a plain-JSON struct in this module, mirroring how the query
// responses hand-roll complex fields as JSON bytes.
type MsgUpdateParams struct {
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority"`
	Params    Params `protobuf:"bytes,2,opt,name=params,proto3" json:"params"`
}

func (m *MsgUpdateParams) ProtoMessage()           {}
func (m *MsgUpdateParams) Reset()                  { *m = MsgUpdateParams{} }
func (m *MsgUpdateParams) String() string          { return fmt.Sprintf("update_params: authority=%s", m.Authority) }
func (m *MsgUpdateParams) XXX_MessageName() string { return "syreen.intent.MsgUpdateParams" }

func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return fmt.Errorf("invalid authority address: %w", err)
	}
	return m.Params.Validate()
}

type MsgUpdateParamsResponse struct{}

func (m *MsgUpdateParamsResponse) ProtoMessage()           {}
func (m *MsgUpdateParamsResponse) Reset()                  { *m = MsgUpdateParamsResponse{} }
func (m *MsgUpdateParamsResponse) String() string          { return "update_params_response" }
func (m *MsgUpdateParamsResponse) XXX_MessageName() string { return "syreen.intent.MsgUpdateParamsResponse" }

// --- MsgSubmitChain ---

// MsgSubmitChain submits a multi-step conditional trading strategy
type MsgSubmitChain struct {
	Creator string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator"`
	Steps []ChainStep `protobuf:"bytes,2,opt,name=steps,proto3" json:"steps"`
	MaxFee sdk.Coins `protobuf:"bytes,3,opt,name=max_fee,proto3" json:"max_fee"`
	Tip sdk.Coins `protobuf:"bytes,4,opt,name=tip,proto3" json:"tip"`
	ExpiryBlocks uint64 `protobuf:"bytes,5,opt,name=expiry_blocks,proto3" json:"expiry_blocks"`
}

func (m *MsgSubmitChain) ProtoMessage()           {}
func (m *MsgSubmitChain) Reset()                  { *m = MsgSubmitChain{} }
func (m *MsgSubmitChain) String() string          { return fmt.Sprintf("submit_chain: creator=%s steps=%d", m.Creator, len(m.Steps)) }
func (m *MsgSubmitChain) XXX_MessageName() string { return "syreen.intent.MsgSubmitChain" }

func (m *MsgSubmitChain) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}
	if len(m.Steps) < 2 {
		return ErrChainTooFewSteps
	}
	if len(m.Steps) > 10 {
		return ErrChainTooManySteps
	}
	for i, step := range m.Steps {
		if err := ValidateChainStep(step, i, len(m.Steps)); err != nil {
			return err
		}
	}
	if !m.MaxFee.IsValid() {
		return fmt.Errorf("invalid max fee")
	}
	if !m.Tip.IsValid() {
		return fmt.Errorf("invalid tip")
	}
	if m.ExpiryBlocks == 0 {
		return fmt.Errorf("expiry blocks must be greater than 0")
	}
	return nil
}

// --- MsgCancelChain ---

// MsgCancelChain cancels an active intent chain and refunds locked tokens
type MsgCancelChain struct {
	Creator string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator"`
	ChainID string `protobuf:"bytes,2,opt,name=chain_id,proto3" json:"chain_id"`
}

func (m *MsgCancelChain) ProtoMessage()           {}
func (m *MsgCancelChain) Reset()                  { *m = MsgCancelChain{} }
func (m *MsgCancelChain) String() string          { return fmt.Sprintf("cancel_chain: chain=%s", m.ChainID) }
func (m *MsgCancelChain) XXX_MessageName() string { return "syreen.intent.MsgCancelChain" }

func (m *MsgCancelChain) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}
	if m.ChainID == "" {
		return fmt.Errorf("chain_id cannot be empty")
	}
	return nil
}

// --- Response types ---

type MsgSubmitChainResponse struct {
	ChainID string `protobuf:"bytes,1,opt,name=chain_id,proto3" json:"chain_id"`
}

func (m *MsgSubmitChainResponse) ProtoMessage()           {}
func (m *MsgSubmitChainResponse) Reset()                  { *m = MsgSubmitChainResponse{} }
func (m *MsgSubmitChainResponse) String() string          { return m.ChainID }
func (m *MsgSubmitChainResponse) XXX_MessageName() string { return "syreen.intent.MsgSubmitChainResponse" }

type MsgCancelChainResponse struct{}

func (m *MsgCancelChainResponse) ProtoMessage()           {}
func (m *MsgCancelChainResponse) Reset()                  { *m = MsgCancelChainResponse{} }
func (m *MsgCancelChainResponse) String() string          { return "cancel_chain_response" }
func (m *MsgCancelChainResponse) XXX_MessageName() string { return "syreen.intent.MsgCancelChainResponse" }
