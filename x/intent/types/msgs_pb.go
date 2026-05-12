package types

// This file adds the XXX_* gogoproto methods and MarshalToSizedBuffer to all
// intent Msg types and their responses. Without these, the SDK's amino/proto
// codec may silently fall back to JSON or panic when deterministic encoding is
// requested.

import (
	proto "github.com/cosmos/gogoproto/proto"
)

// ---------------------------------------------------------------------------
// MsgSubmitIntent
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgSubmitIntent proto.InternalMessageInfo

func (m *MsgSubmitIntent) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgSubmitIntent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSubmitIntent.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgSubmitIntent) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgSubmitIntent.Merge(m, src) }
func (m *MsgSubmitIntent) XXX_Size() int                { return m.Size() }
func (m *MsgSubmitIntent) XXX_DiscardUnknown()          { xxx_messageInfo_MsgSubmitIntent.DiscardUnknown(m) }

func (m *MsgSubmitIntent) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil {
		return 0, err
	}
	n := len(bz)
	copy(dAtA[len(dAtA)-n:], bz)
	return n, nil
}

// ---------------------------------------------------------------------------
// MsgSubmitIntentResponse
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgSubmitIntentResponse proto.InternalMessageInfo

func (m *MsgSubmitIntentResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgSubmitIntentResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSubmitIntentResponse.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgSubmitIntentResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSubmitIntentResponse.Merge(m, src)
}
func (m *MsgSubmitIntentResponse) XXX_Size() int       { return m.Size() }
func (m *MsgSubmitIntentResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgSubmitIntentResponse.DiscardUnknown(m) }

func (m *MsgSubmitIntentResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil {
		return 0, err
	}
	n := len(bz)
	copy(dAtA[len(dAtA)-n:], bz)
	return n, nil
}

// ---------------------------------------------------------------------------
// MsgRegisterSolver
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgRegisterSolver proto.InternalMessageInfo

func (m *MsgRegisterSolver) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRegisterSolver) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRegisterSolver.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgRegisterSolver) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgRegisterSolver.Merge(m, src) }
func (m *MsgRegisterSolver) XXX_Size() int                { return m.Size() }
func (m *MsgRegisterSolver) XXX_DiscardUnknown()          { xxx_messageInfo_MsgRegisterSolver.DiscardUnknown(m) }

func (m *MsgRegisterSolver) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil {
		return 0, err
	}
	n := len(bz)
	copy(dAtA[len(dAtA)-n:], bz)
	return n, nil
}

// ---------------------------------------------------------------------------
// MsgRegisterSolverResponse
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgRegisterSolverResponse proto.InternalMessageInfo

func (m *MsgRegisterSolverResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRegisterSolverResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRegisterSolverResponse.Marshal(b, m, deterministic)
	}
	return nil, nil
}
func (m *MsgRegisterSolverResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgRegisterSolverResponse.Merge(m, src)
}
func (m *MsgRegisterSolverResponse) XXX_Size() int       { return m.Size() }
func (m *MsgRegisterSolverResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgRegisterSolverResponse.DiscardUnknown(m) }

func (m *MsgRegisterSolverResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	return len(dAtA), nil
}

// ---------------------------------------------------------------------------
// MsgDeregisterSolver
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgDeregisterSolver proto.InternalMessageInfo

func (m *MsgDeregisterSolver) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgDeregisterSolver) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgDeregisterSolver.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgDeregisterSolver) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgDeregisterSolver.Merge(m, src) }
func (m *MsgDeregisterSolver) XXX_Size() int                { return m.Size() }
func (m *MsgDeregisterSolver) XXX_DiscardUnknown()          { xxx_messageInfo_MsgDeregisterSolver.DiscardUnknown(m) }

func (m *MsgDeregisterSolver) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil {
		return 0, err
	}
	n := len(bz)
	copy(dAtA[len(dAtA)-n:], bz)
	return n, nil
}

// ---------------------------------------------------------------------------
// MsgDeregisterSolverResponse
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgDeregisterSolverResponse proto.InternalMessageInfo

func (m *MsgDeregisterSolverResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgDeregisterSolverResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgDeregisterSolverResponse.Marshal(b, m, deterministic)
	}
	return nil, nil
}
func (m *MsgDeregisterSolverResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgDeregisterSolverResponse.Merge(m, src)
}
func (m *MsgDeregisterSolverResponse) XXX_Size() int       { return m.Size() }
func (m *MsgDeregisterSolverResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgDeregisterSolverResponse.DiscardUnknown(m) }

func (m *MsgDeregisterSolverResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	return len(dAtA), nil
}

// ---------------------------------------------------------------------------
// MsgSubmitSolution
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgSubmitSolution proto.InternalMessageInfo

func (m *MsgSubmitSolution) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgSubmitSolution) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSubmitSolution.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgSubmitSolution) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgSubmitSolution.Merge(m, src) }
func (m *MsgSubmitSolution) XXX_Size() int                { return m.Size() }
func (m *MsgSubmitSolution) XXX_DiscardUnknown()          { xxx_messageInfo_MsgSubmitSolution.DiscardUnknown(m) }

func (m *MsgSubmitSolution) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil {
		return 0, err
	}
	n := len(bz)
	copy(dAtA[len(dAtA)-n:], bz)
	return n, nil
}

// ---------------------------------------------------------------------------
// MsgSubmitSolutionResponse
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgSubmitSolutionResponse proto.InternalMessageInfo

func (m *MsgSubmitSolutionResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgSubmitSolutionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSubmitSolutionResponse.Marshal(b, m, deterministic)
	}
	return nil, nil
}
func (m *MsgSubmitSolutionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSubmitSolutionResponse.Merge(m, src)
}
func (m *MsgSubmitSolutionResponse) XXX_Size() int       { return m.Size() }
func (m *MsgSubmitSolutionResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgSubmitSolutionResponse.DiscardUnknown(m) }

func (m *MsgSubmitSolutionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	return len(dAtA), nil
}

// ---------------------------------------------------------------------------
// MsgFulfillIntent
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgFulfillIntent proto.InternalMessageInfo

func (m *MsgFulfillIntent) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgFulfillIntent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgFulfillIntent.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgFulfillIntent) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgFulfillIntent.Merge(m, src) }
func (m *MsgFulfillIntent) XXX_Size() int                { return m.Size() }
func (m *MsgFulfillIntent) XXX_DiscardUnknown()          { xxx_messageInfo_MsgFulfillIntent.DiscardUnknown(m) }

func (m *MsgFulfillIntent) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil {
		return 0, err
	}
	n := len(bz)
	copy(dAtA[len(dAtA)-n:], bz)
	return n, nil
}

// ---------------------------------------------------------------------------
// MsgFulfillIntentResponse
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgFulfillIntentResponse proto.InternalMessageInfo

func (m *MsgFulfillIntentResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgFulfillIntentResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgFulfillIntentResponse.Marshal(b, m, deterministic)
	}
	return nil, nil
}
func (m *MsgFulfillIntentResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgFulfillIntentResponse.Merge(m, src)
}
func (m *MsgFulfillIntentResponse) XXX_Size() int       { return m.Size() }
func (m *MsgFulfillIntentResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgFulfillIntentResponse.DiscardUnknown(m) }

func (m *MsgFulfillIntentResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	return len(dAtA), nil
}

// ---------------------------------------------------------------------------
// MsgSubmitChain
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgSubmitChain proto.InternalMessageInfo

func (m *MsgSubmitChain) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgSubmitChain) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSubmitChain.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgSubmitChain) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgSubmitChain.Merge(m, src) }
func (m *MsgSubmitChain) XXX_Size() int                { return m.Size() }
func (m *MsgSubmitChain) XXX_DiscardUnknown()          { xxx_messageInfo_MsgSubmitChain.DiscardUnknown(m) }

func (m *MsgSubmitChain) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil {
		return 0, err
	}
	n := len(bz)
	copy(dAtA[len(dAtA)-n:], bz)
	return n, nil
}

// ---------------------------------------------------------------------------
// MsgSubmitChainResponse
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgSubmitChainResponse proto.InternalMessageInfo

func (m *MsgSubmitChainResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgSubmitChainResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSubmitChainResponse.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgSubmitChainResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSubmitChainResponse.Merge(m, src)
}
func (m *MsgSubmitChainResponse) XXX_Size() int       { return m.Size() }
func (m *MsgSubmitChainResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgSubmitChainResponse.DiscardUnknown(m) }

func (m *MsgSubmitChainResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil {
		return 0, err
	}
	n := len(bz)
	copy(dAtA[len(dAtA)-n:], bz)
	return n, nil
}

// ---------------------------------------------------------------------------
// MsgCancelChain
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCancelChain proto.InternalMessageInfo

func (m *MsgCancelChain) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCancelChain) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCancelChain.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgCancelChain) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgCancelChain.Merge(m, src) }
func (m *MsgCancelChain) XXX_Size() int                { return m.Size() }
func (m *MsgCancelChain) XXX_DiscardUnknown()          { xxx_messageInfo_MsgCancelChain.DiscardUnknown(m) }

func (m *MsgCancelChain) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil {
		return 0, err
	}
	n := len(bz)
	copy(dAtA[len(dAtA)-n:], bz)
	return n, nil
}

// ---------------------------------------------------------------------------
// MsgCancelChainResponse
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCancelChainResponse proto.InternalMessageInfo

func (m *MsgCancelChainResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCancelChainResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCancelChainResponse.Marshal(b, m, deterministic)
	}
	return nil, nil
}
func (m *MsgCancelChainResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCancelChainResponse.Merge(m, src)
}
func (m *MsgCancelChainResponse) XXX_Size() int       { return m.Size() }
func (m *MsgCancelChainResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgCancelChainResponse.DiscardUnknown(m) }

func (m *MsgCancelChainResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	return len(dAtA), nil
}

// ---------------------------------------------------------------------------
// MsgCreateStrategy
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCreateStrategy proto.InternalMessageInfo

func (m *MsgCreateStrategy) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateStrategy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateStrategy.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgCreateStrategy) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgCreateStrategy.Merge(m, src) }
func (m *MsgCreateStrategy) XXX_Size() int                { return m.Size() }
func (m *MsgCreateStrategy) XXX_DiscardUnknown()          { xxx_messageInfo_MsgCreateStrategy.DiscardUnknown(m) }

func (m *MsgCreateStrategy) MarshalTo(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil {
		return 0, err
	}
	copy(dAtA, bz)
	return len(bz), nil
}

func (m *MsgCreateStrategy) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil {
		return 0, err
	}
	n := len(bz)
	copy(dAtA[len(dAtA)-n:], bz)
	return n, nil
}

func (m *MsgCreateStrategy) Size() int {
	bz, _ := m.Marshal()
	return len(bz)
}

// ---------------------------------------------------------------------------
// MsgCreateStrategyResponse
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCreateStrategyResponse proto.InternalMessageInfo

func (m *MsgCreateStrategyResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateStrategyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateStrategyResponse.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgCreateStrategyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCreateStrategyResponse.Merge(m, src)
}
func (m *MsgCreateStrategyResponse) XXX_Size() int       { return m.Size() }
func (m *MsgCreateStrategyResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgCreateStrategyResponse.DiscardUnknown(m) }

func (m *MsgCreateStrategyResponse) MarshalTo(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil {
		return 0, err
	}
	copy(dAtA, bz)
	return len(bz), nil
}

func (m *MsgCreateStrategyResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil {
		return 0, err
	}
	n := len(bz)
	copy(dAtA[len(dAtA)-n:], bz)
	return n, nil
}

func (m *MsgCreateStrategyResponse) Size() int {
	bz, _ := m.Marshal()
	return len(bz)
}

// ---------------------------------------------------------------------------
// MsgCancelStrategy
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCancelStrategy proto.InternalMessageInfo

func (m *MsgCancelStrategy) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCancelStrategy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCancelStrategy.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgCancelStrategy) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgCancelStrategy.Merge(m, src) }
func (m *MsgCancelStrategy) XXX_Size() int                { return m.Size() }
func (m *MsgCancelStrategy) XXX_DiscardUnknown()          { xxx_messageInfo_MsgCancelStrategy.DiscardUnknown(m) }

func (m *MsgCancelStrategy) MarshalTo(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil {
		return 0, err
	}
	copy(dAtA, bz)
	return len(bz), nil
}

func (m *MsgCancelStrategy) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil {
		return 0, err
	}
	n := len(bz)
	copy(dAtA[len(dAtA)-n:], bz)
	return n, nil
}

func (m *MsgCancelStrategy) Size() int {
	bz, _ := m.Marshal()
	return len(bz)
}

// ---------------------------------------------------------------------------
// MsgCancelStrategyResponse
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCancelStrategyResponse proto.InternalMessageInfo

func (m *MsgCancelStrategyResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCancelStrategyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCancelStrategyResponse.Marshal(b, m, deterministic)
	}
	return nil, nil
}
func (m *MsgCancelStrategyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCancelStrategyResponse.Merge(m, src)
}
func (m *MsgCancelStrategyResponse) XXX_Size() int       { return m.Size() }
func (m *MsgCancelStrategyResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgCancelStrategyResponse.DiscardUnknown(m) }

func (m *MsgCancelStrategyResponse) MarshalTo(dAtA []byte) (int, error) {
	return 0, nil
}

func (m *MsgCancelStrategyResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	return len(dAtA), nil
}

func (m *MsgCancelStrategyResponse) Size() int {
	return 0
}
