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

var fileDescriptorPortfolioTx []byte

// Descriptor methods
func (*MsgRecordTrade) Descriptor() ([]byte, []int)               { return fileDescriptorPortfolioTx, []int{0} }
func (*MsgRecordTradeResponse) Descriptor() ([]byte, []int)       { return fileDescriptorPortfolioTx, []int{1} }
func (*MsgCreateCompetition) Descriptor() ([]byte, []int)         { return fileDescriptorPortfolioTx, []int{2} }
func (*MsgCreateCompetitionResponse) Descriptor() ([]byte, []int) { return fileDescriptorPortfolioTx, []int{3} }
func (*MsgJoinCompetition) Descriptor() ([]byte, []int)           { return fileDescriptorPortfolioTx, []int{4} }
func (*MsgJoinCompetitionResponse) Descriptor() ([]byte, []int)   { return fileDescriptorPortfolioTx, []int{5} }
func (*MsgEndCompetition) Descriptor() ([]byte, []int)            { return fileDescriptorPortfolioTx, []int{6} }
func (*MsgEndCompetitionResponse) Descriptor() ([]byte, []int)    { return fileDescriptorPortfolioTx, []int{7} }
func (*MsgUpdatePortfolio) Descriptor() ([]byte, []int)           { return fileDescriptorPortfolioTx, []int{8} }
func (*MsgUpdatePortfolioResponse) Descriptor() ([]byte, []int)   { return fileDescriptorPortfolioTx, []int{9} }

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

// ---------------------------------------------------------------------------
// Helper functions
// ---------------------------------------------------------------------------

var (
	ErrIntOverflow  = fmt.Errorf("proto: integer overflow")
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
			if shift >= 64 {
				return 0, ErrIntOverflow
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflow
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflow
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLength
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, fmt.Errorf("proto: unexpected end group")
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLength
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

// marshalMathInt marshals a math.Int as bytes (json encoding for complex types)
func marshalMathInt(dAtA []byte, i int, v interface{ MarshalJSON() ([]byte, error) }) (int, error) {
	bz, err := json.Marshal(v)
	if err != nil {
		return 0, err
	}
	// Strip quotes from JSON string representation
	if len(bz) >= 2 && bz[0] == '"' && bz[len(bz)-1] == '"' {
		bz = bz[1 : len(bz)-1]
	}
	i -= len(bz)
	copy(dAtA[i:], bz)
	i = encodeVarint(dAtA, i, uint64(len(bz)))
	return i, nil
}

func mathIntSize(v interface{ MarshalJSON() ([]byte, error) }) int {
	bz, _ := json.Marshal(v)
	if len(bz) >= 2 && bz[0] == '"' && bz[len(bz)-1] == '"' {
		bz = bz[1 : len(bz)-1]
	}
	l := len(bz)
	return 1 + l + sov(uint64(l))
}

// ---------------------------------------------------------------------------
// MsgRecordTrade
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgRecordTrade proto.InternalMessageInfo

func (m *MsgRecordTrade) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRecordTrade) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRecordTrade.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgRecordTrade) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgRecordTrade.Merge(m, src) }
func (m *MsgRecordTrade) XXX_Size() int               { return m.Size() }
func (m *MsgRecordTrade) XXX_DiscardUnknown()          { xxx_messageInfo_MsgRecordTrade.DiscardUnknown(m) }

func (m *MsgRecordTrade) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgRecordTrade) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgRecordTrade) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	// field 6: pnl (bytes)
	if !m.PnL.IsNil() {
		ii, err := marshalMathInt(dAtA, i, &m.PnL)
		if err != nil {
			return 0, err
		}
		i = ii
		i--
		dAtA[i] = 0x32
	}
	// field 5: price (bytes)
	if !m.Price.IsNil() {
		ii, err := marshalMathInt(dAtA, i, &m.Price)
		if err != nil {
			return 0, err
		}
		i = ii
		i--
		dAtA[i] = 0x2a
	}
	// field 4: amount (bytes)
	if !m.Amount.IsNil() {
		ii, err := marshalMathInt(dAtA, i, &m.Amount)
		if err != nil {
			return 0, err
		}
		i = ii
		i--
		dAtA[i] = 0x22
	}
	// field 3: denom (string)
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarint(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0x1a
	}
	// field 2: trade_type (string)
	if len(m.TradeType) > 0 {
		i -= len(m.TradeType)
		copy(dAtA[i:], m.TradeType)
		i = encodeVarint(dAtA, i, uint64(len(m.TradeType)))
		i--
		dAtA[i] = 0x12
	}
	// field 1: sender (string)
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarint(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}

func (m *MsgRecordTrade) Size() (n int) {
	if m == nil {
		return 0
	}
	l := len(m.Sender)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	l = len(m.TradeType)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	if !m.Amount.IsNil() {
		n += mathIntSize(&m.Amount)
	}
	if !m.Price.IsNil() {
		n += mathIntSize(&m.Price)
	}
	if !m.PnL.IsNil() {
		n += mathIntSize(&m.PnL)
	}
	return n
}

func (m *MsgRecordTrade) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflow
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgRecordTrade: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgRecordTrade: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1: // sender
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 { break }
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 { return ErrInvalidLength }
			postIndex := iNdEx + intStringLen
			if postIndex < 0 { return ErrInvalidLength }
			if postIndex > l { return io.ErrUnexpectedEOF }
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2: // trade_type
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TradeType", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 { break }
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 { return ErrInvalidLength }
			postIndex := iNdEx + intStringLen
			if postIndex < 0 { return ErrInvalidLength }
			if postIndex > l { return io.ErrUnexpectedEOF }
			m.TradeType = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3: // denom
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 { break }
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 { return ErrInvalidLength }
			postIndex := iNdEx + intStringLen
			if postIndex < 0 { return ErrInvalidLength }
			if postIndex > l { return io.ErrUnexpectedEOF }
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4: // amount (bytes -> math.Int)
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 { break }
			}
			if byteLen < 0 { return ErrInvalidLength }
			postIndex := iNdEx + byteLen
			if postIndex < 0 { return ErrInvalidLength }
			if postIndex > l { return io.ErrUnexpectedEOF }
			// math.Int from string repr
			quoted := fmt.Sprintf(`"%s"`, string(dAtA[iNdEx:postIndex]))
			if err := json.Unmarshal([]byte(quoted), &m.Amount); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5: // price (bytes -> math.Int)
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Price", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 { break }
			}
			if byteLen < 0 { return ErrInvalidLength }
			postIndex := iNdEx + byteLen
			if postIndex < 0 { return ErrInvalidLength }
			if postIndex > l { return io.ErrUnexpectedEOF }
			quoted := fmt.Sprintf(`"%s"`, string(dAtA[iNdEx:postIndex]))
			if err := json.Unmarshal([]byte(quoted), &m.Price); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6: // pnl (bytes -> math.Int)
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PnL", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 { break }
			}
			if byteLen < 0 { return ErrInvalidLength }
			postIndex := iNdEx + byteLen
			if postIndex < 0 { return ErrInvalidLength }
			if postIndex > l { return io.ErrUnexpectedEOF }
			quoted := fmt.Sprintf(`"%s"`, string(dAtA[iNdEx:postIndex]))
			if err := json.Unmarshal([]byte(quoted), &m.PnL); err != nil {
				return err
			}
			iNdEx = postIndex
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
// MsgRecordTradeResponse
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgRecordTradeResponse proto.InternalMessageInfo

func (m *MsgRecordTradeResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRecordTradeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRecordTradeResponse.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgRecordTradeResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgRecordTradeResponse.Merge(m, src) }
func (m *MsgRecordTradeResponse) XXX_Size() int               { return m.Size() }
func (m *MsgRecordTradeResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgRecordTradeResponse.DiscardUnknown(m) }

func (m *MsgRecordTradeResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgRecordTradeResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgRecordTradeResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.TradeID != 0 {
		i = encodeVarint(dAtA, i, m.TradeID)
		i--
		dAtA[i] = 0x08
	}
	return len(dAtA) - i, nil
}

func (m *MsgRecordTradeResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	if m.TradeID != 0 {
		n += 1 + sov(m.TradeID)
	}
	return n
}

func (m *MsgRecordTradeResponse) Unmarshal(dAtA []byte) error {
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
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgRecordTradeResponse: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: // trade_id
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TradeID", wireType)
			}
			m.TradeID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				m.TradeID |= uint64(b&0x7F) << shift
				if b < 0x80 { break }
			}
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
// MsgCreateCompetition
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCreateCompetition proto.InternalMessageInfo

func (m *MsgCreateCompetition) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateCompetition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateCompetition.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgCreateCompetition) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCreateCompetition.Merge(m, src) }
func (m *MsgCreateCompetition) XXX_Size() int               { return m.Size() }
func (m *MsgCreateCompetition) XXX_DiscardUnknown()          { xxx_messageInfo_MsgCreateCompetition.DiscardUnknown(m) }

func (m *MsgCreateCompetition) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCreateCompetition) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCreateCompetition) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	// field 8: max_participants (varint)
	if m.MaxParticipants != 0 {
		i = encodeVarint(dAtA, i, m.MaxParticipants)
		i--
		dAtA[i] = 0x40
	}
	// field 7: entry_fee (bytes)
	if !m.EntryFee.IsNil() {
		ii, err := marshalMathInt(dAtA, i, &m.EntryFee)
		if err != nil {
			return 0, err
		}
		i = ii
		i--
		dAtA[i] = 0x3a
	}
	// field 6: prize_pool (bytes)
	if !m.PrizePool.IsNil() {
		ii, err := marshalMathInt(dAtA, i, &m.PrizePool)
		if err != nil {
			return 0, err
		}
		i = ii
		i--
		dAtA[i] = 0x32
	}
	// field 5: prize_denom (string)
	if len(m.PrizeDenom) > 0 {
		i -= len(m.PrizeDenom)
		copy(dAtA[i:], m.PrizeDenom)
		i = encodeVarint(dAtA, i, uint64(len(m.PrizeDenom)))
		i--
		dAtA[i] = 0x2a
	}
	// field 4: end_block (varint)
	if m.EndBlock != 0 {
		i = encodeVarint(dAtA, i, uint64(m.EndBlock))
		i--
		dAtA[i] = 0x20
	}
	// field 3: start_block (varint)
	if m.StartBlock != 0 {
		i = encodeVarint(dAtA, i, uint64(m.StartBlock))
		i--
		dAtA[i] = 0x18
	}
	// field 2: name (string)
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarint(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x12
	}
	// field 1: creator (string)
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarint(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}

func (m *MsgCreateCompetition) Size() (n int) {
	if m == nil {
		return 0
	}
	l := len(m.Creator)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	if m.StartBlock != 0 {
		n += 1 + sov(uint64(m.StartBlock))
	}
	if m.EndBlock != 0 {
		n += 1 + sov(uint64(m.EndBlock))
	}
	l = len(m.PrizeDenom)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	if !m.PrizePool.IsNil() {
		n += mathIntSize(&m.PrizePool)
	}
	if !m.EntryFee.IsNil() {
		n += mathIntSize(&m.EntryFee)
	}
	if m.MaxParticipants != 0 {
		n += 1 + sov(m.MaxParticipants)
	}
	return n
}

func (m *MsgCreateCompetition) Unmarshal(dAtA []byte) error {
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
		if wireType == 4 { return fmt.Errorf("proto: MsgCreateCompetition: wiretype end group for non-group") }
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgCreateCompetition: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: // creator
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 { break }
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 { return ErrInvalidLength }
			postIndex := iNdEx + intStringLen
			if postIndex < 0 { return ErrInvalidLength }
			if postIndex > l { return io.ErrUnexpectedEOF }
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2: // name
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 { break }
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 { return ErrInvalidLength }
			postIndex := iNdEx + intStringLen
			if postIndex < 0 { return ErrInvalidLength }
			if postIndex > l { return io.ErrUnexpectedEOF }
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3: // start_block
			if wireType != 0 { return fmt.Errorf("proto: wrong wireType = %d for field StartBlock", wireType) }
			m.StartBlock = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				m.StartBlock |= int64(b&0x7F) << shift
				if b < 0x80 { break }
			}
		case 4: // end_block
			if wireType != 0 { return fmt.Errorf("proto: wrong wireType = %d for field EndBlock", wireType) }
			m.EndBlock = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				m.EndBlock |= int64(b&0x7F) << shift
				if b < 0x80 { break }
			}
		case 5: // prize_denom
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field PrizeDenom", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 { break }
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 { return ErrInvalidLength }
			postIndex := iNdEx + intStringLen
			if postIndex < 0 { return ErrInvalidLength }
			if postIndex > l { return io.ErrUnexpectedEOF }
			m.PrizeDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6: // prize_pool (bytes -> math.Int)
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field PrizePool", wireType) }
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 { break }
			}
			if byteLen < 0 { return ErrInvalidLength }
			postIndex := iNdEx + byteLen
			if postIndex < 0 { return ErrInvalidLength }
			if postIndex > l { return io.ErrUnexpectedEOF }
			quoted := fmt.Sprintf(`"%s"`, string(dAtA[iNdEx:postIndex]))
			if err := json.Unmarshal([]byte(quoted), &m.PrizePool); err != nil { return err }
			iNdEx = postIndex
		case 7: // entry_fee (bytes -> math.Int)
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field EntryFee", wireType) }
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 { break }
			}
			if byteLen < 0 { return ErrInvalidLength }
			postIndex := iNdEx + byteLen
			if postIndex < 0 { return ErrInvalidLength }
			if postIndex > l { return io.ErrUnexpectedEOF }
			quoted := fmt.Sprintf(`"%s"`, string(dAtA[iNdEx:postIndex]))
			if err := json.Unmarshal([]byte(quoted), &m.EntryFee); err != nil { return err }
			iNdEx = postIndex
		case 8: // max_participants
			if wireType != 0 { return fmt.Errorf("proto: wrong wireType = %d for field MaxParticipants", wireType) }
			m.MaxParticipants = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				m.MaxParticipants |= uint64(b&0x7F) << shift
				if b < 0x80 { break }
			}
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
// MsgCreateCompetitionResponse
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCreateCompetitionResponse proto.InternalMessageInfo

func (m *MsgCreateCompetitionResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateCompetitionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateCompetitionResponse.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgCreateCompetitionResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCreateCompetitionResponse.Merge(m, src) }
func (m *MsgCreateCompetitionResponse) XXX_Size() int               { return m.Size() }
func (m *MsgCreateCompetitionResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgCreateCompetitionResponse.DiscardUnknown(m) }

func (m *MsgCreateCompetitionResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCreateCompetitionResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCreateCompetitionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.CompetitionID != 0 {
		i = encodeVarint(dAtA, i, m.CompetitionID)
		i--
		dAtA[i] = 0x08
	}
	return len(dAtA) - i, nil
}

func (m *MsgCreateCompetitionResponse) Size() (n int) {
	if m == nil { return 0 }
	if m.CompetitionID != 0 {
		n += 1 + sov(m.CompetitionID)
	}
	return n
}

func (m *MsgCreateCompetitionResponse) Unmarshal(dAtA []byte) error {
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
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgCreateCompetitionResponse: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if int(wire&0x7) != 0 { return fmt.Errorf("proto: wrong wireType for field CompetitionID") }
			m.CompetitionID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				m.CompetitionID |= uint64(b&0x7F) << shift
				if b < 0x80 { break }
			}
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
// MsgJoinCompetition
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgJoinCompetition proto.InternalMessageInfo

func (m *MsgJoinCompetition) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgJoinCompetition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgJoinCompetition.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgJoinCompetition) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgJoinCompetition.Merge(m, src) }
func (m *MsgJoinCompetition) XXX_Size() int               { return m.Size() }
func (m *MsgJoinCompetition) XXX_DiscardUnknown()          { xxx_messageInfo_MsgJoinCompetition.DiscardUnknown(m) }

func (m *MsgJoinCompetition) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgJoinCompetition) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgJoinCompetition) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.CompetitionID != 0 {
		i = encodeVarint(dAtA, i, m.CompetitionID)
		i--
		dAtA[i] = 0x10
	}
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarint(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}

func (m *MsgJoinCompetition) Size() (n int) {
	if m == nil { return 0 }
	l := len(m.Sender)
	if l > 0 { n += 1 + l + sov(uint64(l)) }
	if m.CompetitionID != 0 { n += 1 + sov(m.CompetitionID) }
	return n
}

func (m *MsgJoinCompetition) Unmarshal(dAtA []byte) error {
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
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgJoinCompetition: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 { break }
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 { return ErrInvalidLength }
			postIndex := iNdEx + intStringLen
			if postIndex < 0 { return ErrInvalidLength }
			if postIndex > l { return io.ErrUnexpectedEOF }
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 { return fmt.Errorf("proto: wrong wireType = %d for field CompetitionID", wireType) }
			m.CompetitionID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				m.CompetitionID |= uint64(b&0x7F) << shift
				if b < 0x80 { break }
			}
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
// MsgJoinCompetitionResponse (empty)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgJoinCompetitionResponse proto.InternalMessageInfo

func (m *MsgJoinCompetitionResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgJoinCompetitionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgJoinCompetitionResponse.Marshal(b, m, deterministic) }
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil { return nil, err }
	return b[:n], nil
}
func (m *MsgJoinCompetitionResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgJoinCompetitionResponse.Merge(m, src) }
func (m *MsgJoinCompetitionResponse) XXX_Size() int               { return m.Size() }
func (m *MsgJoinCompetitionResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgJoinCompetitionResponse.DiscardUnknown(m) }

func (m *MsgJoinCompetitionResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgJoinCompetitionResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgJoinCompetitionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgJoinCompetitionResponse) Size() (n int) { return 0 }
func (m *MsgJoinCompetitionResponse) Unmarshal(dAtA []byte) error {
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
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgJoinCompetitionResponse: illegal tag %d", fieldNum) }
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
// MsgEndCompetition
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgEndCompetition proto.InternalMessageInfo

func (m *MsgEndCompetition) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgEndCompetition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgEndCompetition.Marshal(b, m, deterministic) }
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil { return nil, err }
	return b[:n], nil
}
func (m *MsgEndCompetition) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgEndCompetition.Merge(m, src) }
func (m *MsgEndCompetition) XXX_Size() int               { return m.Size() }
func (m *MsgEndCompetition) XXX_DiscardUnknown()          { xxx_messageInfo_MsgEndCompetition.DiscardUnknown(m) }

func (m *MsgEndCompetition) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil { return nil, err }
	return dAtA[:n], nil
}
func (m *MsgEndCompetition) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgEndCompetition) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.CompetitionID != 0 {
		i = encodeVarint(dAtA, i, m.CompetitionID)
		i--
		dAtA[i] = 0x10
	}
	if len(m.Authority) > 0 {
		i -= len(m.Authority)
		copy(dAtA[i:], m.Authority)
		i = encodeVarint(dAtA, i, uint64(len(m.Authority)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}
func (m *MsgEndCompetition) Size() (n int) {
	if m == nil { return 0 }
	l := len(m.Authority)
	if l > 0 { n += 1 + l + sov(uint64(l)) }
	if m.CompetitionID != 0 { n += 1 + sov(m.CompetitionID) }
	return n
}
func (m *MsgEndCompetition) Unmarshal(dAtA []byte) error {
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
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgEndCompetition: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Authority", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 { break }
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 { return ErrInvalidLength }
			postIndex := iNdEx + intStringLen
			if postIndex < 0 { return ErrInvalidLength }
			if postIndex > l { return io.ErrUnexpectedEOF }
			m.Authority = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 { return fmt.Errorf("proto: wrong wireType = %d for field CompetitionID", wireType) }
			m.CompetitionID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				m.CompetitionID |= uint64(b&0x7F) << shift
				if b < 0x80 { break }
			}
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
// MsgEndCompetitionResponse
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgEndCompetitionResponse proto.InternalMessageInfo

func (m *MsgEndCompetitionResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgEndCompetitionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgEndCompetitionResponse.Marshal(b, m, deterministic) }
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil { return nil, err }
	return b[:n], nil
}
func (m *MsgEndCompetitionResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgEndCompetitionResponse.Merge(m, src) }
func (m *MsgEndCompetitionResponse) XXX_Size() int               { return m.Size() }
func (m *MsgEndCompetitionResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgEndCompetitionResponse.DiscardUnknown(m) }

func (m *MsgEndCompetitionResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil { return nil, err }
	return dAtA[:n], nil
}
func (m *MsgEndCompetitionResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgEndCompetitionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	// winners is repeated string, field 1
	if len(m.Winners) > 0 {
		for iNdEx := len(m.Winners) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Winners[iNdEx])
			copy(dAtA[i:], m.Winners[iNdEx])
			i = encodeVarint(dAtA, i, uint64(len(m.Winners[iNdEx])))
			i--
			dAtA[i] = 0x0a
		}
	}
	return len(dAtA) - i, nil
}
func (m *MsgEndCompetitionResponse) Size() (n int) {
	if m == nil { return 0 }
	if len(m.Winners) > 0 {
		for _, s := range m.Winners {
			l := len(s)
			n += 1 + l + sov(uint64(l))
		}
	}
	return n
}
func (m *MsgEndCompetitionResponse) Unmarshal(dAtA []byte) error {
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
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgEndCompetitionResponse: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Winners", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 { break }
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 { return ErrInvalidLength }
			postIndex := iNdEx + intStringLen
			if postIndex < 0 { return ErrInvalidLength }
			if postIndex > l { return io.ErrUnexpectedEOF }
			m.Winners = append(m.Winners, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
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
// MsgUpdatePortfolio
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgUpdatePortfolio proto.InternalMessageInfo

func (m *MsgUpdatePortfolio) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgUpdatePortfolio) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgUpdatePortfolio.Marshal(b, m, deterministic) }
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil { return nil, err }
	return b[:n], nil
}
func (m *MsgUpdatePortfolio) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgUpdatePortfolio.Merge(m, src) }
func (m *MsgUpdatePortfolio) XXX_Size() int               { return m.Size() }
func (m *MsgUpdatePortfolio) XXX_DiscardUnknown()          { xxx_messageInfo_MsgUpdatePortfolio.DiscardUnknown(m) }

func (m *MsgUpdatePortfolio) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil { return nil, err }
	return dAtA[:n], nil
}
func (m *MsgUpdatePortfolio) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgUpdatePortfolio) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarint(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}
func (m *MsgUpdatePortfolio) Size() (n int) {
	if m == nil { return 0 }
	l := len(m.Sender)
	if l > 0 { n += 1 + l + sov(uint64(l)) }
	return n
}
func (m *MsgUpdatePortfolio) Unmarshal(dAtA []byte) error {
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
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgUpdatePortfolio: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 { break }
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 { return ErrInvalidLength }
			postIndex := iNdEx + intStringLen
			if postIndex < 0 { return ErrInvalidLength }
			if postIndex > l { return io.ErrUnexpectedEOF }
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
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
// MsgUpdatePortfolioResponse
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgUpdatePortfolioResponse proto.InternalMessageInfo

func (m *MsgUpdatePortfolioResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgUpdatePortfolioResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgUpdatePortfolioResponse.Marshal(b, m, deterministic) }
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil { return nil, err }
	return b[:n], nil
}
func (m *MsgUpdatePortfolioResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgUpdatePortfolioResponse.Merge(m, src) }
func (m *MsgUpdatePortfolioResponse) XXX_Size() int               { return m.Size() }
func (m *MsgUpdatePortfolioResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgUpdatePortfolioResponse.DiscardUnknown(m) }

func (m *MsgUpdatePortfolioResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil { return nil, err }
	return dAtA[:n], nil
}
func (m *MsgUpdatePortfolioResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgUpdatePortfolioResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.TotalValue.IsNil() {
		ii, err := marshalMathInt(dAtA, i, &m.TotalValue)
		if err != nil { return 0, err }
		i = ii
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}
func (m *MsgUpdatePortfolioResponse) Size() (n int) {
	if m == nil { return 0 }
	if !m.TotalValue.IsNil() {
		n += mathIntSize(&m.TotalValue)
	}
	return n
}
func (m *MsgUpdatePortfolioResponse) Unmarshal(dAtA []byte) error {
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
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgUpdatePortfolioResponse: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field TotalValue", wireType) }
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 { break }
			}
			if byteLen < 0 { return ErrInvalidLength }
			postIndex := iNdEx + byteLen
			if postIndex < 0 { return ErrInvalidLength }
			if postIndex > l { return io.ErrUnexpectedEOF }
			quoted := fmt.Sprintf(`"%s"`, string(dAtA[iNdEx:postIndex]))
			if err := json.Unmarshal([]byte(quoted), &m.TotalValue); err != nil { return err }
			iNdEx = postIndex
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
