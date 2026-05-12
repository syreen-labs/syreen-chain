package types

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math/bits"

	proto "github.com/cosmos/gogoproto/proto"
)

var _ = proto.Marshal
var _ = fmt.Errorf
var _ = io.ErrUnexpectedEOF
var _ = binary.BigEndian
var _ = bits.Len64

var fileDescriptorAIAgentTx []byte

func (m *MsgCreateAgent) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgCreateAgent) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCreateAgent) XXX_Size() int                         { return m.Size() }

func (m *MsgCreateAgentResponse) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgCreateAgentResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCreateAgentResponse) XXX_Size() int                         { return m.Size() }

func (m *MsgFundAgent) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgFundAgent) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgFundAgent) XXX_Size() int                         { return m.Size() }

func (m *MsgFundAgentResponse) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgFundAgentResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgFundAgentResponse) XXX_Size() int                         { return m.Size() }

func (m *MsgWithdrawAgentFunds) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgWithdrawAgentFunds) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgWithdrawAgentFunds) XXX_Size() int                         { return m.Size() }

func (m *MsgWithdrawAgentFundsResponse) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgWithdrawAgentFundsResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgWithdrawAgentFundsResponse) XXX_Size() int                         { return m.Size() }

func (m *MsgPauseAgent) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgPauseAgent) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgPauseAgent) XXX_Size() int                         { return m.Size() }

func (m *MsgPauseAgentResponse) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgPauseAgentResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgPauseAgentResponse) XXX_Size() int                         { return m.Size() }

func (m *MsgResumeAgent) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgResumeAgent) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgResumeAgent) XXX_Size() int                         { return m.Size() }

func (m *MsgResumeAgentResponse) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgResumeAgentResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgResumeAgentResponse) XXX_Size() int                         { return m.Size() }

func (m *MsgUpdateAgentStrategy) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgUpdateAgentStrategy) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgUpdateAgentStrategy) XXX_Size() int                         { return m.Size() }

func (m *MsgUpdateAgentStrategyResponse) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgUpdateAgentStrategyResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgUpdateAgentStrategyResponse) XXX_Size() int                         { return m.Size() }

// Descriptors
func (*MsgCreateAgent) Descriptor() ([]byte, []int)                  { return fileDescriptorAIAgentTx, []int{0} }
func (*MsgCreateAgentResponse) Descriptor() ([]byte, []int)          { return fileDescriptorAIAgentTx, []int{1} }
func (*MsgFundAgent) Descriptor() ([]byte, []int)                    { return fileDescriptorAIAgentTx, []int{2} }
func (*MsgFundAgentResponse) Descriptor() ([]byte, []int)            { return fileDescriptorAIAgentTx, []int{3} }
func (*MsgWithdrawAgentFunds) Descriptor() ([]byte, []int)           { return fileDescriptorAIAgentTx, []int{4} }
func (*MsgWithdrawAgentFundsResponse) Descriptor() ([]byte, []int)   { return fileDescriptorAIAgentTx, []int{5} }
func (*MsgPauseAgent) Descriptor() ([]byte, []int)                   { return fileDescriptorAIAgentTx, []int{6} }
func (*MsgPauseAgentResponse) Descriptor() ([]byte, []int)           { return fileDescriptorAIAgentTx, []int{7} }
func (*MsgResumeAgent) Descriptor() ([]byte, []int)                  { return fileDescriptorAIAgentTx, []int{8} }
func (*MsgResumeAgentResponse) Descriptor() ([]byte, []int)          { return fileDescriptorAIAgentTx, []int{9} }
func (*MsgUpdateAgentStrategy) Descriptor() ([]byte, []int)          { return fileDescriptorAIAgentTx, []int{10} }
func (*MsgUpdateAgentStrategyResponse) Descriptor() ([]byte, []int)  { return fileDescriptorAIAgentTx, []int{11} }

func init() {
	proto.RegisterType((*MsgCreateAgent)(nil), "syreen.aiagent.MsgCreateAgent")
	proto.RegisterType((*MsgCreateAgentResponse)(nil), "syreen.aiagent.MsgCreateAgentResponse")
	proto.RegisterType((*MsgFundAgent)(nil), "syreen.aiagent.MsgFundAgent")
	proto.RegisterType((*MsgFundAgentResponse)(nil), "syreen.aiagent.MsgFundAgentResponse")
	proto.RegisterType((*MsgWithdrawAgentFunds)(nil), "syreen.aiagent.MsgWithdrawAgentFunds")
	proto.RegisterType((*MsgWithdrawAgentFundsResponse)(nil), "syreen.aiagent.MsgWithdrawAgentFundsResponse")
	proto.RegisterType((*MsgPauseAgent)(nil), "syreen.aiagent.MsgPauseAgent")
	proto.RegisterType((*MsgPauseAgentResponse)(nil), "syreen.aiagent.MsgPauseAgentResponse")
	proto.RegisterType((*MsgResumeAgent)(nil), "syreen.aiagent.MsgResumeAgent")
	proto.RegisterType((*MsgResumeAgentResponse)(nil), "syreen.aiagent.MsgResumeAgentResponse")
	proto.RegisterType((*MsgUpdateAgentStrategy)(nil), "syreen.aiagent.MsgUpdateAgentStrategy")
	proto.RegisterType((*MsgUpdateAgentStrategyResponse)(nil), "syreen.aiagent.MsgUpdateAgentStrategyResponse")
}

// Marshal/Unmarshal/Size
func (m *MsgCreateAgent) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgCreateAgent) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCreateAgent) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCreateAgent) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCreateAgent) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgCreateAgentResponse) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgCreateAgentResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCreateAgentResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCreateAgentResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCreateAgentResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgFundAgent) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgFundAgent) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgFundAgent) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgFundAgent) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgFundAgent) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgFundAgentResponse) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgFundAgentResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgFundAgentResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgFundAgentResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgFundAgentResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgWithdrawAgentFunds) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgWithdrawAgentFunds) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgWithdrawAgentFunds) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgWithdrawAgentFunds) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgWithdrawAgentFunds) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgWithdrawAgentFundsResponse) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgWithdrawAgentFundsResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgWithdrawAgentFundsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgWithdrawAgentFundsResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgWithdrawAgentFundsResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgPauseAgent) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgPauseAgent) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgPauseAgent) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgPauseAgent) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgPauseAgent) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgPauseAgentResponse) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgPauseAgentResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgPauseAgentResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgPauseAgentResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgPauseAgentResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgResumeAgent) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgResumeAgent) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgResumeAgent) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgResumeAgent) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgResumeAgent) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgResumeAgentResponse) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgResumeAgentResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgResumeAgentResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgResumeAgentResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgResumeAgentResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgUpdateAgentStrategy) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgUpdateAgentStrategy) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgUpdateAgentStrategy) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgUpdateAgentStrategy) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgUpdateAgentStrategy) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgUpdateAgentStrategyResponse) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgUpdateAgentStrategyResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgUpdateAgentStrategyResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgUpdateAgentStrategyResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgUpdateAgentStrategyResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }
