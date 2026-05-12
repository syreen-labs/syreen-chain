package types

import (
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

var fileDescriptorAIAgentTx []byte

// Descriptors
func (*MsgCreateAgent) Descriptor() ([]byte, []int)                 { return fileDescriptorAIAgentTx, []int{0} }
func (*MsgCreateAgentResponse) Descriptor() ([]byte, []int)         { return fileDescriptorAIAgentTx, []int{1} }
func (*MsgFundAgent) Descriptor() ([]byte, []int)                   { return fileDescriptorAIAgentTx, []int{2} }
func (*MsgFundAgentResponse) Descriptor() ([]byte, []int)           { return fileDescriptorAIAgentTx, []int{3} }
func (*MsgWithdrawAgentFunds) Descriptor() ([]byte, []int)          { return fileDescriptorAIAgentTx, []int{4} }
func (*MsgWithdrawAgentFundsResponse) Descriptor() ([]byte, []int)  { return fileDescriptorAIAgentTx, []int{5} }
func (*MsgPauseAgent) Descriptor() ([]byte, []int)                  { return fileDescriptorAIAgentTx, []int{6} }
func (*MsgPauseAgentResponse) Descriptor() ([]byte, []int)          { return fileDescriptorAIAgentTx, []int{7} }
func (*MsgResumeAgent) Descriptor() ([]byte, []int)                 { return fileDescriptorAIAgentTx, []int{8} }
func (*MsgResumeAgentResponse) Descriptor() ([]byte, []int)         { return fileDescriptorAIAgentTx, []int{9} }
func (*MsgUpdateAgentStrategy) Descriptor() ([]byte, []int)         { return fileDescriptorAIAgentTx, []int{10} }
func (*MsgUpdateAgentStrategyResponse) Descriptor() ([]byte, []int) { return fileDescriptorAIAgentTx, []int{11} }

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

// ---------------------------------------------------------------------------
// Helper functions
// ---------------------------------------------------------------------------

var (
	ErrIntOverflow   = fmt.Errorf("proto: integer overflow")
	ErrInvalidLength = fmt.Errorf("proto: negative length found during unmarshaling")
)

func encodeVarint(dAtA []byte, offset int, v uint64) int {
	offset -= sov(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}

func sov(x uint64) int {
	return (bits.Len64(x|1) + 6) / 7
}

func skip(dAtA []byte) (int, error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 { return 0, ErrIntOverflow }
			if iNdEx >= l { return 0, io.ErrUnexpectedEOF }
			b := dAtA[iNdEx]; iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 { break }
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return 0, ErrIntOverflow }
				if iNdEx >= l { return 0, io.ErrUnexpectedEOF }
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 { break }
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return 0, ErrIntOverflow }
				if iNdEx >= l { return 0, io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 { break }
			}
			if length < 0 { return 0, ErrInvalidLength }
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 { return 0, fmt.Errorf("proto: unexpected end group") }
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 { return 0, ErrInvalidLength }
		if depth == 0 { return iNdEx, nil }
	}
	return 0, io.ErrUnexpectedEOF
}

func marshalBytes(dAtA []byte, i int, v interface{}) (int, error) {
	bz, err := json.Marshal(v)
	if err != nil { return 0, err }
	if len(bz) >= 2 && bz[0] == '"' && bz[len(bz)-1] == '"' {
		bz = bz[1 : len(bz)-1]
	}
	i -= len(bz)
	copy(dAtA[i:], bz)
	i = encodeVarint(dAtA, i, uint64(len(bz)))
	return i, nil
}

func bytesFieldSize(v interface{}) int {
	bz, _ := json.Marshal(v)
	if len(bz) >= 2 && bz[0] == '"' && bz[len(bz)-1] == '"' {
		bz = bz[1 : len(bz)-1]
	}
	l := len(bz)
	return 1 + l + sov(uint64(l))
}

// marshalJSONBytes marshals a complex struct (AgentConfig) as JSON bytes for the bytes field
func marshalJSONBytes(dAtA []byte, i int, v interface{}) (int, error) {
	bz, err := json.Marshal(v)
	if err != nil { return 0, err }
	i -= len(bz)
	copy(dAtA[i:], bz)
	i = encodeVarint(dAtA, i, uint64(len(bz)))
	return i, nil
}

func jsonBytesFieldSize(v interface{}) int {
	bz, _ := json.Marshal(v)
	l := len(bz)
	return 1 + l + sov(uint64(l))
}

// unmarshalStringField is a helper for reading a string field
func unmarshalStringField(dAtA []byte, iNdEx int, l int) (string, int, error) {
	var stringLen uint64
	for shift := uint(0); ; shift += 7 {
		if shift >= 64 { return "", 0, ErrIntOverflow }
		if iNdEx >= l { return "", 0, io.ErrUnexpectedEOF }
		b := dAtA[iNdEx]; iNdEx++
		stringLen |= uint64(b&0x7F) << shift
		if b < 0x80 { break }
	}
	intStringLen := int(stringLen)
	if intStringLen < 0 { return "", 0, ErrInvalidLength }
	postIndex := iNdEx + intStringLen
	if postIndex < 0 { return "", 0, ErrInvalidLength }
	if postIndex > l { return "", 0, io.ErrUnexpectedEOF }
	return string(dAtA[iNdEx:postIndex]), postIndex, nil
}

// unmarshalVarintField reads a varint
func unmarshalVarintField(dAtA []byte, iNdEx int, l int) (uint64, int, error) {
	var v uint64
	for shift := uint(0); ; shift += 7 {
		if shift >= 64 { return 0, 0, ErrIntOverflow }
		if iNdEx >= l { return 0, 0, io.ErrUnexpectedEOF }
		b := dAtA[iNdEx]; iNdEx++
		v |= uint64(b&0x7F) << shift
		if b < 0x80 { break }
	}
	return v, iNdEx, nil
}

// unmarshalBytesField reads a length-delimited bytes field
func unmarshalBytesField(dAtA []byte, iNdEx int, l int) ([]byte, int, error) {
	var byteLen int
	for shift := uint(0); ; shift += 7 {
		if shift >= 64 { return nil, 0, ErrIntOverflow }
		if iNdEx >= l { return nil, 0, io.ErrUnexpectedEOF }
		b := dAtA[iNdEx]; iNdEx++
		byteLen |= int(b&0x7F) << shift
		if b < 0x80 { break }
	}
	if byteLen < 0 { return nil, 0, ErrInvalidLength }
	postIndex := iNdEx + byteLen
	if postIndex < 0 { return nil, 0, ErrInvalidLength }
	if postIndex > l { return nil, 0, io.ErrUnexpectedEOF }
	return dAtA[iNdEx:postIndex], postIndex, nil
}

// ---------------------------------------------------------------------------
// MsgCreateAgent
// fields: owner(1,string), name(2,string), strategy_type(3,string),
//         config(4,bytes/json), initial_funds(5,bytes/mathInt)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCreateAgent proto.InternalMessageInfo

func (m *MsgCreateAgent) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateAgent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgCreateAgent.Marshal(b, m, deterministic) }
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil { return nil, err }
	return b[:n], nil
}
func (m *MsgCreateAgent) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCreateAgent.Merge(m, src) }
func (m *MsgCreateAgent) XXX_Size() int               { return m.Size() }
func (m *MsgCreateAgent) XXX_DiscardUnknown()          { xxx_messageInfo_MsgCreateAgent.DiscardUnknown(m) }

func (m *MsgCreateAgent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil { return nil, err }
	return dAtA[:n], nil
}
func (m *MsgCreateAgent) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgCreateAgent) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	// field 5: initial_funds (bytes)
	if !m.InitialFunds.IsNil() {
		ii, err := marshalBytes(dAtA, i, &m.InitialFunds)
		if err != nil { return 0, err }
		i = ii
		i--
		dAtA[i] = 0x2a
	}
	// field 4: config (bytes/json)
	{
		ii, err := marshalJSONBytes(dAtA, i, &m.Config)
		if err != nil { return 0, err }
		i = ii
		i--
		dAtA[i] = 0x22
	}
	// field 3: strategy_type (string)
	if len(m.StrategyType) > 0 {
		s := string(m.StrategyType)
		i -= len(s)
		copy(dAtA[i:], s)
		i = encodeVarint(dAtA, i, uint64(len(s)))
		i--
		dAtA[i] = 0x1a
	}
	// field 2: name (string)
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarint(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x12
	}
	// field 1: owner (string)
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarint(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}
func (m *MsgCreateAgent) Size() (n int) {
	if m == nil { return 0 }
	l := len(m.Owner)
	if l > 0 { n += 1 + l + sov(uint64(l)) }
	l = len(m.Name)
	if l > 0 { n += 1 + l + sov(uint64(l)) }
	l = len(string(m.StrategyType))
	if l > 0 { n += 1 + l + sov(uint64(l)) }
	n += jsonBytesFieldSize(&m.Config)
	if !m.InitialFunds.IsNil() { n += bytesFieldSize(&m.InitialFunds) }
	return n
}
func (m *MsgCreateAgent) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 { return ErrIntOverflow }
			if iNdEx >= l { return io.ErrUnexpectedEOF }
			b := dAtA[iNdEx]; iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 { break }
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgCreateAgent: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: // owner
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for Owner") }
			s, next, err := unmarshalStringField(dAtA, iNdEx, l)
			if err != nil { return err }
			m.Owner = s; iNdEx = next
		case 2: // name
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for Name") }
			s, next, err := unmarshalStringField(dAtA, iNdEx, l)
			if err != nil { return err }
			m.Name = s; iNdEx = next
		case 3: // strategy_type
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for StrategyType") }
			s, next, err := unmarshalStringField(dAtA, iNdEx, l)
			if err != nil { return err }
			m.StrategyType = StrategyType(s); iNdEx = next
		case 4: // config (json bytes)
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for Config") }
			bz, next, err := unmarshalBytesField(dAtA, iNdEx, l)
			if err != nil { return err }
			if err := json.Unmarshal(bz, &m.Config); err != nil { return err }
			iNdEx = next
		case 5: // initial_funds (math.Int as bytes)
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for InitialFunds") }
			bz, next, err := unmarshalBytesField(dAtA, iNdEx, l)
			if err != nil { return err }
			quoted := fmt.Sprintf(`"%s"`, string(bz))
			if err := json.Unmarshal([]byte(quoted), &m.InitialFunds); err != nil { return err }
			iNdEx = next
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }
			if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }
			iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }
	return nil
}

// ---------------------------------------------------------------------------
// MsgCreateAgentResponse
// fields: agent_id(1,varint), agent_address(2,string)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCreateAgentResponse proto.InternalMessageInfo

func (m *MsgCreateAgentResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateAgentResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgCreateAgentResponse.Marshal(b, m, deterministic) }
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil { return nil, err }
	return b[:n], nil
}
func (m *MsgCreateAgentResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCreateAgentResponse.Merge(m, src) }
func (m *MsgCreateAgentResponse) XXX_Size() int               { return m.Size() }
func (m *MsgCreateAgentResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgCreateAgentResponse.DiscardUnknown(m) }

func (m *MsgCreateAgentResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil { return nil, err }
	return dAtA[:n], nil
}
func (m *MsgCreateAgentResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgCreateAgentResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.AgentAddress) > 0 {
		i -= len(m.AgentAddress)
		copy(dAtA[i:], m.AgentAddress)
		i = encodeVarint(dAtA, i, uint64(len(m.AgentAddress)))
		i--
		dAtA[i] = 0x12
	}
	if m.AgentID != 0 {
		i = encodeVarint(dAtA, i, m.AgentID)
		i--
		dAtA[i] = 0x08
	}
	return len(dAtA) - i, nil
}
func (m *MsgCreateAgentResponse) Size() (n int) {
	if m == nil { return 0 }
	if m.AgentID != 0 { n += 1 + sov(m.AgentID) }
	l := len(m.AgentAddress)
	if l > 0 { n += 1 + l + sov(uint64(l)) }
	return n
}
func (m *MsgCreateAgentResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 { return ErrIntOverflow }
			if iNdEx >= l { return io.ErrUnexpectedEOF }
			b := dAtA[iNdEx]; iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 { break }
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgCreateAgentResponse: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if wireType != 0 { return fmt.Errorf("proto: wrong wireType for AgentID") }
			v, next, err := unmarshalVarintField(dAtA, iNdEx, l)
			if err != nil { return err }
			m.AgentID = v; iNdEx = next
		case 2:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for AgentAddress") }
			s, next, err := unmarshalStringField(dAtA, iNdEx, l)
			if err != nil { return err }
			m.AgentAddress = s; iNdEx = next
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }
			if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }
			iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }
	return nil
}

// ---------------------------------------------------------------------------
// MsgFundAgent: owner(1,string), agent_id(2,varint), amount(3,bytes)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgFundAgent proto.InternalMessageInfo

func (m *MsgFundAgent) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgFundAgent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgFundAgent.Marshal(b, m, deterministic) }
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil { return nil, err }
	return b[:n], nil
}
func (m *MsgFundAgent) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgFundAgent.Merge(m, src) }
func (m *MsgFundAgent) XXX_Size() int               { return m.Size() }
func (m *MsgFundAgent) XXX_DiscardUnknown()          { xxx_messageInfo_MsgFundAgent.DiscardUnknown(m) }

func (m *MsgFundAgent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil { return nil, err }
	return dAtA[:n], nil
}
func (m *MsgFundAgent) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgFundAgent) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.Amount.IsNil() {
		ii, err := marshalBytes(dAtA, i, &m.Amount)
		if err != nil { return 0, err }
		i = ii
		i--
		dAtA[i] = 0x1a
	}
	if m.AgentID != 0 {
		i = encodeVarint(dAtA, i, m.AgentID)
		i--
		dAtA[i] = 0x10
	}
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarint(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}
func (m *MsgFundAgent) Size() (n int) {
	if m == nil { return 0 }
	l := len(m.Owner)
	if l > 0 { n += 1 + l + sov(uint64(l)) }
	if m.AgentID != 0 { n += 1 + sov(m.AgentID) }
	if !m.Amount.IsNil() { n += bytesFieldSize(&m.Amount) }
	return n
}
func (m *MsgFundAgent) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 { return ErrIntOverflow }
			if iNdEx >= l { return io.ErrUnexpectedEOF }
			b := dAtA[iNdEx]; iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 { break }
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgFundAgent: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for Owner") }
			s, next, err := unmarshalStringField(dAtA, iNdEx, l)
			if err != nil { return err }
			m.Owner = s; iNdEx = next
		case 2:
			if wireType != 0 { return fmt.Errorf("proto: wrong wireType for AgentID") }
			v, next, err := unmarshalVarintField(dAtA, iNdEx, l)
			if err != nil { return err }
			m.AgentID = v; iNdEx = next
		case 3:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for Amount") }
			bz, next, err := unmarshalBytesField(dAtA, iNdEx, l)
			if err != nil { return err }
			quoted := fmt.Sprintf(`"%s"`, string(bz))
			if err := json.Unmarshal([]byte(quoted), &m.Amount); err != nil { return err }
			iNdEx = next
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }
			if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }
			iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }
	return nil
}

// ---------------------------------------------------------------------------
// MsgFundAgentResponse (empty)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgFundAgentResponse proto.InternalMessageInfo

func (m *MsgFundAgentResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgFundAgentResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgFundAgentResponse.Marshal(b, m, deterministic) }
	return nil, nil
}
func (m *MsgFundAgentResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgFundAgentResponse.Merge(m, src) }
func (m *MsgFundAgentResponse) XXX_Size() int               { return 0 }
func (m *MsgFundAgentResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgFundAgentResponse.DiscardUnknown(m) }

func (m *MsgFundAgentResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgFundAgentResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgFundAgentResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgFundAgentResponse) Size() (n int) { return 0 }
func (m *MsgFundAgentResponse) Unmarshal(dAtA []byte) error { return skipAll(dAtA, "MsgFundAgentResponse") }

// ---------------------------------------------------------------------------
// MsgWithdrawAgentFunds: owner(1,string), agent_id(2,varint), amount(3,bytes)
// Same layout as MsgFundAgent
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgWithdrawAgentFunds proto.InternalMessageInfo

func (m *MsgWithdrawAgentFunds) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgWithdrawAgentFunds) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgWithdrawAgentFunds.Marshal(b, m, deterministic) }
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil { return nil, err }
	return b[:n], nil
}
func (m *MsgWithdrawAgentFunds) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgWithdrawAgentFunds.Merge(m, src) }
func (m *MsgWithdrawAgentFunds) XXX_Size() int               { return m.Size() }
func (m *MsgWithdrawAgentFunds) XXX_DiscardUnknown()          { xxx_messageInfo_MsgWithdrawAgentFunds.DiscardUnknown(m) }

func (m *MsgWithdrawAgentFunds) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil { return nil, err }
	return dAtA[:n], nil
}
func (m *MsgWithdrawAgentFunds) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgWithdrawAgentFunds) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.Amount.IsNil() {
		ii, err := marshalBytes(dAtA, i, &m.Amount)
		if err != nil { return 0, err }
		i = ii
		i--
		dAtA[i] = 0x1a
	}
	if m.AgentID != 0 {
		i = encodeVarint(dAtA, i, m.AgentID)
		i--
		dAtA[i] = 0x10
	}
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarint(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}
func (m *MsgWithdrawAgentFunds) Size() (n int) {
	if m == nil { return 0 }
	l := len(m.Owner)
	if l > 0 { n += 1 + l + sov(uint64(l)) }
	if m.AgentID != 0 { n += 1 + sov(m.AgentID) }
	if !m.Amount.IsNil() { n += bytesFieldSize(&m.Amount) }
	return n
}
func (m *MsgWithdrawAgentFunds) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 { return ErrIntOverflow }
			if iNdEx >= l { return io.ErrUnexpectedEOF }
			b := dAtA[iNdEx]; iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 { break }
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgWithdrawAgentFunds: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for Owner") }
			s, next, err := unmarshalStringField(dAtA, iNdEx, l)
			if err != nil { return err }
			m.Owner = s; iNdEx = next
		case 2:
			if wireType != 0 { return fmt.Errorf("proto: wrong wireType for AgentID") }
			v, next, err := unmarshalVarintField(dAtA, iNdEx, l)
			if err != nil { return err }
			m.AgentID = v; iNdEx = next
		case 3:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for Amount") }
			bz, next, err := unmarshalBytesField(dAtA, iNdEx, l)
			if err != nil { return err }
			quoted := fmt.Sprintf(`"%s"`, string(bz))
			if err := json.Unmarshal([]byte(quoted), &m.Amount); err != nil { return err }
			iNdEx = next
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }
			if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }
			iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }
	return nil
}

// ---------------------------------------------------------------------------
// MsgWithdrawAgentFundsResponse: amount_returned(1,bytes)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgWithdrawAgentFundsResponse proto.InternalMessageInfo

func (m *MsgWithdrawAgentFundsResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgWithdrawAgentFundsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgWithdrawAgentFundsResponse.Marshal(b, m, deterministic) }
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil { return nil, err }
	return b[:n], nil
}
func (m *MsgWithdrawAgentFundsResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgWithdrawAgentFundsResponse.Merge(m, src) }
func (m *MsgWithdrawAgentFundsResponse) XXX_Size() int               { return m.Size() }
func (m *MsgWithdrawAgentFundsResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgWithdrawAgentFundsResponse.DiscardUnknown(m) }

func (m *MsgWithdrawAgentFundsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil { return nil, err }
	return dAtA[:n], nil
}
func (m *MsgWithdrawAgentFundsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgWithdrawAgentFundsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.AmountReturned.IsNil() {
		ii, err := marshalBytes(dAtA, i, &m.AmountReturned)
		if err != nil { return 0, err }
		i = ii
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}
func (m *MsgWithdrawAgentFundsResponse) Size() (n int) {
	if m == nil { return 0 }
	if !m.AmountReturned.IsNil() { n += bytesFieldSize(&m.AmountReturned) }
	return n
}
func (m *MsgWithdrawAgentFundsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 { return ErrIntOverflow }
			if iNdEx >= l { return io.ErrUnexpectedEOF }
			b := dAtA[iNdEx]; iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 { break }
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgWithdrawAgentFundsResponse: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for AmountReturned") }
			bz, next, err := unmarshalBytesField(dAtA, iNdEx, l)
			if err != nil { return err }
			quoted := fmt.Sprintf(`"%s"`, string(bz))
			if err := json.Unmarshal([]byte(quoted), &m.AmountReturned); err != nil { return err }
			iNdEx = next
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }
			if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }
			iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }
	return nil
}

// ---------------------------------------------------------------------------
// Helper for owner+agent_id pattern (MsgPauseAgent, MsgResumeAgent)
// ---------------------------------------------------------------------------

func marshalOwnerAgentID(dAtA []byte, owner string, agentID uint64) (int, error) {
	i := len(dAtA)
	if agentID != 0 {
		i = encodeVarint(dAtA, i, agentID)
		i--
		dAtA[i] = 0x10
	}
	if len(owner) > 0 {
		i -= len(owner)
		copy(dAtA[i:], owner)
		i = encodeVarint(dAtA, i, uint64(len(owner)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}

func sizeOwnerAgentID(owner string, agentID uint64) int {
	n := 0
	l := len(owner)
	if l > 0 { n += 1 + l + sov(uint64(l)) }
	if agentID != 0 { n += 1 + sov(agentID) }
	return n
}

func unmarshalOwnerAgentID(dAtA []byte, msgName string) (string, uint64, error) {
	l := len(dAtA)
	iNdEx := 0
	var owner string
	var agentID uint64
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 { return "", 0, ErrIntOverflow }
			if iNdEx >= l { return "", 0, io.ErrUnexpectedEOF }
			b := dAtA[iNdEx]; iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 { break }
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if fieldNum <= 0 { return "", 0, fmt.Errorf("proto: %s: illegal tag %d", msgName, fieldNum) }
		switch fieldNum {
		case 1:
			if wireType != 2 { return "", 0, fmt.Errorf("proto: wrong wireType for Owner") }
			s, next, err := unmarshalStringField(dAtA, iNdEx, l)
			if err != nil { return "", 0, err }
			owner = s; iNdEx = next
		case 2:
			if wireType != 0 { return "", 0, fmt.Errorf("proto: wrong wireType for AgentID") }
			v, next, err := unmarshalVarintField(dAtA, iNdEx, l)
			if err != nil { return "", 0, err }
			agentID = v; iNdEx = next
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil { return "", 0, err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return "", 0, ErrInvalidLength }
			if (iNdEx + skippy) > l { return "", 0, io.ErrUnexpectedEOF }
			iNdEx += skippy
		}
	}
	return owner, agentID, nil
}

// skipAll skips all fields in an empty response
func skipAll(dAtA []byte, msgName string) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 { return ErrIntOverflow }
			if iNdEx >= l { return io.ErrUnexpectedEOF }
			b := dAtA[iNdEx]; iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 { break }
		}
		fieldNum := int32(wire >> 3)
		if fieldNum <= 0 { return fmt.Errorf("proto: %s: illegal tag %d", msgName, fieldNum) }
		iNdEx = preIndex
		skippy, err := skip(dAtA[iNdEx:])
		if err != nil { return err }
		if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }
		if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }
		iNdEx += skippy
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgPauseAgent: owner(1,string), agent_id(2,varint)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgPauseAgent proto.InternalMessageInfo

func (m *MsgPauseAgent) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgPauseAgent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgPauseAgent.Marshal(b, m, deterministic) }
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil { return nil, err }
	return b[:n], nil
}
func (m *MsgPauseAgent) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgPauseAgent.Merge(m, src) }
func (m *MsgPauseAgent) XXX_Size() int               { return m.Size() }
func (m *MsgPauseAgent) XXX_DiscardUnknown()          { xxx_messageInfo_MsgPauseAgent.DiscardUnknown(m) }

func (m *MsgPauseAgent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil { return nil, err }
	return dAtA[:n], nil
}
func (m *MsgPauseAgent) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgPauseAgent) MarshalToSizedBuffer(dAtA []byte) (int, error) { return marshalOwnerAgentID(dAtA, m.Owner, m.AgentID) }
func (m *MsgPauseAgent) Size() (n int) { if m == nil { return 0 }; return sizeOwnerAgentID(m.Owner, m.AgentID) }
func (m *MsgPauseAgent) Unmarshal(dAtA []byte) error {
	owner, agentID, err := unmarshalOwnerAgentID(dAtA, "MsgPauseAgent")
	if err != nil { return err }
	m.Owner = owner; m.AgentID = agentID
	return nil
}

// ---------------------------------------------------------------------------
// MsgPauseAgentResponse (empty)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgPauseAgentResponse proto.InternalMessageInfo

func (m *MsgPauseAgentResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgPauseAgentResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgPauseAgentResponse.Marshal(b, m, deterministic) }
	return nil, nil
}
func (m *MsgPauseAgentResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgPauseAgentResponse.Merge(m, src) }
func (m *MsgPauseAgentResponse) XXX_Size() int               { return 0 }
func (m *MsgPauseAgentResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgPauseAgentResponse.DiscardUnknown(m) }

func (m *MsgPauseAgentResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgPauseAgentResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgPauseAgentResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgPauseAgentResponse) Size() (n int) { return 0 }
func (m *MsgPauseAgentResponse) Unmarshal(dAtA []byte) error { return skipAll(dAtA, "MsgPauseAgentResponse") }

// ---------------------------------------------------------------------------
// MsgResumeAgent: owner(1,string), agent_id(2,varint)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgResumeAgent proto.InternalMessageInfo

func (m *MsgResumeAgent) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgResumeAgent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgResumeAgent.Marshal(b, m, deterministic) }
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil { return nil, err }
	return b[:n], nil
}
func (m *MsgResumeAgent) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgResumeAgent.Merge(m, src) }
func (m *MsgResumeAgent) XXX_Size() int               { return m.Size() }
func (m *MsgResumeAgent) XXX_DiscardUnknown()          { xxx_messageInfo_MsgResumeAgent.DiscardUnknown(m) }

func (m *MsgResumeAgent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil { return nil, err }
	return dAtA[:n], nil
}
func (m *MsgResumeAgent) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgResumeAgent) MarshalToSizedBuffer(dAtA []byte) (int, error) { return marshalOwnerAgentID(dAtA, m.Owner, m.AgentID) }
func (m *MsgResumeAgent) Size() (n int) { if m == nil { return 0 }; return sizeOwnerAgentID(m.Owner, m.AgentID) }
func (m *MsgResumeAgent) Unmarshal(dAtA []byte) error {
	owner, agentID, err := unmarshalOwnerAgentID(dAtA, "MsgResumeAgent")
	if err != nil { return err }
	m.Owner = owner; m.AgentID = agentID
	return nil
}

// ---------------------------------------------------------------------------
// MsgResumeAgentResponse (empty)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgResumeAgentResponse proto.InternalMessageInfo

func (m *MsgResumeAgentResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgResumeAgentResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgResumeAgentResponse.Marshal(b, m, deterministic) }
	return nil, nil
}
func (m *MsgResumeAgentResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgResumeAgentResponse.Merge(m, src) }
func (m *MsgResumeAgentResponse) XXX_Size() int               { return 0 }
func (m *MsgResumeAgentResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgResumeAgentResponse.DiscardUnknown(m) }

func (m *MsgResumeAgentResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgResumeAgentResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgResumeAgentResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgResumeAgentResponse) Size() (n int) { return 0 }
func (m *MsgResumeAgentResponse) Unmarshal(dAtA []byte) error { return skipAll(dAtA, "MsgResumeAgentResponse") }

// ---------------------------------------------------------------------------
// MsgUpdateAgentStrategy: owner(1,string), agent_id(2,varint), config(3,bytes)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgUpdateAgentStrategy proto.InternalMessageInfo

func (m *MsgUpdateAgentStrategy) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgUpdateAgentStrategy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgUpdateAgentStrategy.Marshal(b, m, deterministic) }
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil { return nil, err }
	return b[:n], nil
}
func (m *MsgUpdateAgentStrategy) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgUpdateAgentStrategy.Merge(m, src) }
func (m *MsgUpdateAgentStrategy) XXX_Size() int               { return m.Size() }
func (m *MsgUpdateAgentStrategy) XXX_DiscardUnknown()          { xxx_messageInfo_MsgUpdateAgentStrategy.DiscardUnknown(m) }

func (m *MsgUpdateAgentStrategy) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil { return nil, err }
	return dAtA[:n], nil
}
func (m *MsgUpdateAgentStrategy) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgUpdateAgentStrategy) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	// field 3: config (bytes/json)
	{
		ii, err := marshalJSONBytes(dAtA, i, &m.Config)
		if err != nil { return 0, err }
		i = ii
		i--
		dAtA[i] = 0x1a
	}
	if m.AgentID != 0 {
		i = encodeVarint(dAtA, i, m.AgentID)
		i--
		dAtA[i] = 0x10
	}
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarint(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}
func (m *MsgUpdateAgentStrategy) Size() (n int) {
	if m == nil { return 0 }
	l := len(m.Owner)
	if l > 0 { n += 1 + l + sov(uint64(l)) }
	if m.AgentID != 0 { n += 1 + sov(m.AgentID) }
	n += jsonBytesFieldSize(&m.Config)
	return n
}
func (m *MsgUpdateAgentStrategy) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 { return ErrIntOverflow }
			if iNdEx >= l { return io.ErrUnexpectedEOF }
			b := dAtA[iNdEx]; iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 { break }
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgUpdateAgentStrategy: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for Owner") }
			s, next, err := unmarshalStringField(dAtA, iNdEx, l)
			if err != nil { return err }
			m.Owner = s; iNdEx = next
		case 2:
			if wireType != 0 { return fmt.Errorf("proto: wrong wireType for AgentID") }
			v, next, err := unmarshalVarintField(dAtA, iNdEx, l)
			if err != nil { return err }
			m.AgentID = v; iNdEx = next
		case 3:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for Config") }
			bz, next, err := unmarshalBytesField(dAtA, iNdEx, l)
			if err != nil { return err }
			if err := json.Unmarshal(bz, &m.Config); err != nil { return err }
			iNdEx = next
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }
			if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }
			iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }
	return nil
}

// ---------------------------------------------------------------------------
// MsgUpdateAgentStrategyResponse (empty)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgUpdateAgentStrategyResponse proto.InternalMessageInfo

func (m *MsgUpdateAgentStrategyResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgUpdateAgentStrategyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgUpdateAgentStrategyResponse.Marshal(b, m, deterministic) }
	return nil, nil
}
func (m *MsgUpdateAgentStrategyResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgUpdateAgentStrategyResponse.Merge(m, src) }
func (m *MsgUpdateAgentStrategyResponse) XXX_Size() int               { return 0 }
func (m *MsgUpdateAgentStrategyResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgUpdateAgentStrategyResponse.DiscardUnknown(m) }

func (m *MsgUpdateAgentStrategyResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgUpdateAgentStrategyResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgUpdateAgentStrategyResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgUpdateAgentStrategyResponse) Size() (n int) { return 0 }
func (m *MsgUpdateAgentStrategyResponse) Unmarshal(dAtA []byte) error { return skipAll(dAtA, "MsgUpdateAgentStrategyResponse") }
