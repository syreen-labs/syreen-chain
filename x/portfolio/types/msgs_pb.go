package types

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math/bits"

	proto "github.com/cosmos/gogoproto/proto"
)

// Reference imports
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = io.ErrUnexpectedEOF
var _ = binary.BigEndian
var _ = bits.Len64

var fileDescriptorPortfolioTx []byte

// XXX methods required by gRPC decoder
func (m *MsgRecordTrade) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRecordTrade) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgRecordTrade) XXX_Size() int { return m.Size() }

func (m *MsgRecordTradeResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRecordTradeResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgRecordTradeResponse) XXX_Size() int { return m.Size() }

func (m *MsgCreateCompetition) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateCompetition) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCreateCompetition) XXX_Size() int { return m.Size() }

func (m *MsgCreateCompetitionResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateCompetitionResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCreateCompetitionResponse) XXX_Size() int { return m.Size() }

func (m *MsgJoinCompetition) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgJoinCompetition) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgJoinCompetition) XXX_Size() int { return m.Size() }

func (m *MsgJoinCompetitionResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgJoinCompetitionResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgJoinCompetitionResponse) XXX_Size() int { return m.Size() }

func (m *MsgEndCompetition) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgEndCompetition) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgEndCompetition) XXX_Size() int { return m.Size() }

func (m *MsgEndCompetitionResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgEndCompetitionResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgEndCompetitionResponse) XXX_Size() int { return m.Size() }

func (m *MsgUpdatePortfolio) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgUpdatePortfolio) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgUpdatePortfolio) XXX_Size() int { return m.Size() }

func (m *MsgUpdatePortfolioResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgUpdatePortfolioResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgUpdatePortfolioResponse) XXX_Size() int { return m.Size() }

// Descriptor methods
func (*MsgRecordTrade) Descriptor() ([]byte, []int)                  { return fileDescriptorPortfolioTx, []int{0} }
func (*MsgRecordTradeResponse) Descriptor() ([]byte, []int)          { return fileDescriptorPortfolioTx, []int{1} }
func (*MsgCreateCompetition) Descriptor() ([]byte, []int)            { return fileDescriptorPortfolioTx, []int{2} }
func (*MsgCreateCompetitionResponse) Descriptor() ([]byte, []int)    { return fileDescriptorPortfolioTx, []int{3} }
func (*MsgJoinCompetition) Descriptor() ([]byte, []int)              { return fileDescriptorPortfolioTx, []int{4} }
func (*MsgJoinCompetitionResponse) Descriptor() ([]byte, []int)      { return fileDescriptorPortfolioTx, []int{5} }
func (*MsgEndCompetition) Descriptor() ([]byte, []int)               { return fileDescriptorPortfolioTx, []int{6} }
func (*MsgEndCompetitionResponse) Descriptor() ([]byte, []int)       { return fileDescriptorPortfolioTx, []int{7} }
func (*MsgUpdatePortfolio) Descriptor() ([]byte, []int)              { return fileDescriptorPortfolioTx, []int{8} }
func (*MsgUpdatePortfolioResponse) Descriptor() ([]byte, []int)      { return fileDescriptorPortfolioTx, []int{9} }

func init() {
	proto.RegisterType((*MsgRecordTrade)(nil), "syreen.portfolio.MsgRecordTrade")
	proto.RegisterType((*MsgRecordTradeResponse)(nil), "syreen.portfolio.MsgRecordTradeResponse")
	proto.RegisterType((*MsgCreateCompetition)(nil), "syreen.portfolio.MsgCreateCompetition")
	proto.RegisterType((*MsgCreateCompetitionResponse)(nil), "syreen.portfolio.MsgCreateCompetitionResponse")
	proto.RegisterType((*MsgJoinCompetition)(nil), "syreen.portfolio.MsgJoinCompetition")
	proto.RegisterType((*MsgJoinCompetitionResponse)(nil), "syreen.portfolio.MsgJoinCompetitionResponse")
	proto.RegisterType((*MsgEndCompetition)(nil), "syreen.portfolio.MsgEndCompetition")
	proto.RegisterType((*MsgEndCompetitionResponse)(nil), "syreen.portfolio.MsgEndCompetitionResponse")
	proto.RegisterType((*MsgUpdatePortfolio)(nil), "syreen.portfolio.MsgUpdatePortfolio")
	proto.RegisterType((*MsgUpdatePortfolioResponse)(nil), "syreen.portfolio.MsgUpdatePortfolioResponse")
}

// Marshal/Unmarshal/Size for all message types using JSON
func (m *MsgRecordTrade) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgRecordTrade) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgRecordTrade) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgRecordTrade) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgRecordTrade) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgRecordTradeResponse) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgRecordTradeResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgRecordTradeResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgRecordTradeResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgRecordTradeResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgCreateCompetition) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgCreateCompetition) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCreateCompetition) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCreateCompetition) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCreateCompetition) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgCreateCompetitionResponse) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgCreateCompetitionResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCreateCompetitionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCreateCompetitionResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCreateCompetitionResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgJoinCompetition) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgJoinCompetition) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgJoinCompetition) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgJoinCompetition) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgJoinCompetition) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgJoinCompetitionResponse) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgJoinCompetitionResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgJoinCompetitionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgJoinCompetitionResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgJoinCompetitionResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgEndCompetition) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgEndCompetition) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgEndCompetition) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgEndCompetition) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgEndCompetition) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgEndCompetitionResponse) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgEndCompetitionResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgEndCompetitionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgEndCompetitionResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgEndCompetitionResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgUpdatePortfolio) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgUpdatePortfolio) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgUpdatePortfolio) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgUpdatePortfolio) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgUpdatePortfolio) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgUpdatePortfolioResponse) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgUpdatePortfolioResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgUpdatePortfolioResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgUpdatePortfolioResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgUpdatePortfolioResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }
