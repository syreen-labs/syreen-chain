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

var fileDescriptorOptionsTx []byte

// Descriptor methods
func (*MsgWriteOption) Descriptor() ([]byte, []int)           { return fileDescriptorOptionsTx, []int{0} }
func (*MsgWriteOptionResponse) Descriptor() ([]byte, []int)   { return fileDescriptorOptionsTx, []int{1} }
func (*MsgBuyOption) Descriptor() ([]byte, []int)             { return fileDescriptorOptionsTx, []int{2} }
func (*MsgBuyOptionResponse) Descriptor() ([]byte, []int)     { return fileDescriptorOptionsTx, []int{3} }
func (*MsgExerciseOption) Descriptor() ([]byte, []int)        { return fileDescriptorOptionsTx, []int{4} }
func (*MsgExerciseOptionResponse) Descriptor() ([]byte, []int) { return fileDescriptorOptionsTx, []int{5} }
func (*MsgCancelOption) Descriptor() ([]byte, []int)          { return fileDescriptorOptionsTx, []int{6} }
func (*MsgCancelOptionResponse) Descriptor() ([]byte, []int)  { return fileDescriptorOptionsTx, []int{7} }

func init() {
	proto.RegisterType((*MsgWriteOption)(nil), "syreen.options.MsgWriteOption")
	proto.RegisterType((*MsgWriteOptionResponse)(nil), "syreen.options.MsgWriteOptionResponse")
	proto.RegisterType((*MsgBuyOption)(nil), "syreen.options.MsgBuyOption")
	proto.RegisterType((*MsgBuyOptionResponse)(nil), "syreen.options.MsgBuyOptionResponse")
	proto.RegisterType((*MsgExerciseOption)(nil), "syreen.options.MsgExerciseOption")
	proto.RegisterType((*MsgExerciseOptionResponse)(nil), "syreen.options.MsgExerciseOptionResponse")
	proto.RegisterType((*MsgCancelOption)(nil), "syreen.options.MsgCancelOption")
	proto.RegisterType((*MsgCancelOptionResponse)(nil), "syreen.options.MsgCancelOptionResponse")
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

func marshalMathInt(dAtA []byte, i int, v interface{ MarshalJSON() ([]byte, error) }) (int, error) {
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

func mathIntSize(v interface{ MarshalJSON() ([]byte, error) }) int {
	bz, _ := json.Marshal(v)
	if len(bz) >= 2 && bz[0] == '"' && bz[len(bz)-1] == '"' {
		bz = bz[1 : len(bz)-1]
	}
	l := len(bz)
	return 1 + l + sov(uint64(l))
}

// ---------------------------------------------------------------------------
// MsgWriteOption
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgWriteOption proto.InternalMessageInfo

func (m *MsgWriteOption) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgWriteOption) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgWriteOption.Marshal(b, m, deterministic) }
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil { return nil, err }
	return b[:n], nil
}
func (m *MsgWriteOption) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgWriteOption.Merge(m, src) }
func (m *MsgWriteOption) XXX_Size() int               { return m.Size() }
func (m *MsgWriteOption) XXX_DiscardUnknown()          { xxx_messageInfo_MsgWriteOption.DiscardUnknown(m) }

func (m *MsgWriteOption) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil { return nil, err }
	return dAtA[:n], nil
}
func (m *MsgWriteOption) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgWriteOption) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	// field 7: custom_premium (bytes)
	if !m.CustomPremium.IsNil() {
		ii, err := marshalMathInt(dAtA, i, &m.CustomPremium)
		if err != nil { return 0, err }
		i = ii
		i--
		dAtA[i] = 0x3a
	}
	// field 6: expiry_block (varint)
	if m.ExpiryBlock != 0 {
		i = encodeVarint(dAtA, i, uint64(m.ExpiryBlock))
		i--
		dAtA[i] = 0x30
	}
	// field 5: amount (bytes)
	if !m.Amount.IsNil() {
		ii, err := marshalMathInt(dAtA, i, &m.Amount)
		if err != nil { return 0, err }
		i = ii
		i--
		dAtA[i] = 0x2a
	}
	// field 4: strike_price (bytes)
	if !m.StrikePrice.IsNil() {
		ii, err := marshalMathInt(dAtA, i, &m.StrikePrice)
		if err != nil { return 0, err }
		i = ii
		i--
		dAtA[i] = 0x22
	}
	// field 3: option_type (string)
	if len(m.OptionType) > 0 {
		s := string(m.OptionType)
		i -= len(s)
		copy(dAtA[i:], s)
		i = encodeVarint(dAtA, i, uint64(len(s)))
		i--
		dAtA[i] = 0x1a
	}
	// field 2: pool_id (varint)
	if m.PoolID != 0 {
		i = encodeVarint(dAtA, i, m.PoolID)
		i--
		dAtA[i] = 0x10
	}
	// field 1: writer (string)
	if len(m.Writer) > 0 {
		i -= len(m.Writer)
		copy(dAtA[i:], m.Writer)
		i = encodeVarint(dAtA, i, uint64(len(m.Writer)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}

func (m *MsgWriteOption) Size() (n int) {
	if m == nil { return 0 }
	l := len(m.Writer)
	if l > 0 { n += 1 + l + sov(uint64(l)) }
	if m.PoolID != 0 { n += 1 + sov(m.PoolID) }
	l = len(string(m.OptionType))
	if l > 0 { n += 1 + l + sov(uint64(l)) }
	if !m.StrikePrice.IsNil() { n += mathIntSize(&m.StrikePrice) }
	if !m.Amount.IsNil() { n += mathIntSize(&m.Amount) }
	if m.ExpiryBlock != 0 { n += 1 + sov(uint64(m.ExpiryBlock)) }
	if !m.CustomPremium.IsNil() { n += mathIntSize(&m.CustomPremium) }
	return n
}

func (m *MsgWriteOption) Unmarshal(dAtA []byte) error {
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
		if wireType == 4 { return fmt.Errorf("proto: MsgWriteOption: wiretype end group for non-group") }
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgWriteOption: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: // writer
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Writer", wireType) }
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
			m.Writer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2: // pool_id
			if wireType != 0 { return fmt.Errorf("proto: wrong wireType = %d for field PoolID", wireType) }
			m.PoolID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				m.PoolID |= uint64(b&0x7F) << shift
				if b < 0x80 { break }
			}
		case 3: // option_type
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field OptionType", wireType) }
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
			m.OptionType = OptionType(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4: // strike_price (bytes -> LegacyDec)
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field StrikePrice", wireType) }
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
			if err := json.Unmarshal([]byte(quoted), &m.StrikePrice); err != nil { return err }
			iNdEx = postIndex
		case 5: // amount (bytes -> math.Int)
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType) }
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
			if err := json.Unmarshal([]byte(quoted), &m.Amount); err != nil { return err }
			iNdEx = postIndex
		case 6: // expiry_block
			if wireType != 0 { return fmt.Errorf("proto: wrong wireType = %d for field ExpiryBlock", wireType) }
			m.ExpiryBlock = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				m.ExpiryBlock |= int64(b&0x7F) << shift
				if b < 0x80 { break }
			}
		case 7: // custom_premium (bytes -> math.Int)
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field CustomPremium", wireType) }
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
			if err := json.Unmarshal([]byte(quoted), &m.CustomPremium); err != nil { return err }
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
// MsgWriteOptionResponse
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgWriteOptionResponse proto.InternalMessageInfo

func (m *MsgWriteOptionResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgWriteOptionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgWriteOptionResponse.Marshal(b, m, deterministic) }
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil { return nil, err }
	return b[:n], nil
}
func (m *MsgWriteOptionResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgWriteOptionResponse.Merge(m, src) }
func (m *MsgWriteOptionResponse) XXX_Size() int               { return m.Size() }
func (m *MsgWriteOptionResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgWriteOptionResponse.DiscardUnknown(m) }

func (m *MsgWriteOptionResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil { return nil, err }
	return dAtA[:n], nil
}
func (m *MsgWriteOptionResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgWriteOptionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	// field 2: premium (bytes)
	if !m.Premium.IsNil() {
		ii, err := marshalMathInt(dAtA, i, &m.Premium)
		if err != nil { return 0, err }
		i = ii
		i--
		dAtA[i] = 0x12
	}
	// field 1: option_id (varint)
	if m.OptionID != 0 {
		i = encodeVarint(dAtA, i, m.OptionID)
		i--
		dAtA[i] = 0x08
	}
	return len(dAtA) - i, nil
}
func (m *MsgWriteOptionResponse) Size() (n int) {
	if m == nil { return 0 }
	if m.OptionID != 0 { n += 1 + sov(m.OptionID) }
	if !m.Premium.IsNil() { n += mathIntSize(&m.Premium) }
	return n
}
func (m *MsgWriteOptionResponse) Unmarshal(dAtA []byte) error {
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
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgWriteOptionResponse: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if wireType != 0 { return fmt.Errorf("proto: wrong wireType for OptionID") }
			m.OptionID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				m.OptionID |= uint64(b&0x7F) << shift
				if b < 0x80 { break }
			}
		case 2:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for Premium") }
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
			if err := json.Unmarshal([]byte(quoted), &m.Premium); err != nil { return err }
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
// MsgBuyOption
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgBuyOption proto.InternalMessageInfo

func (m *MsgBuyOption) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgBuyOption) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgBuyOption.Marshal(b, m, deterministic) }
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil { return nil, err }
	return b[:n], nil
}
func (m *MsgBuyOption) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgBuyOption.Merge(m, src) }
func (m *MsgBuyOption) XXX_Size() int               { return m.Size() }
func (m *MsgBuyOption) XXX_DiscardUnknown()          { xxx_messageInfo_MsgBuyOption.DiscardUnknown(m) }

func (m *MsgBuyOption) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil { return nil, err }
	return dAtA[:n], nil
}
func (m *MsgBuyOption) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgBuyOption) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.OptionID != 0 {
		i = encodeVarint(dAtA, i, m.OptionID)
		i--
		dAtA[i] = 0x10
	}
	if len(m.Buyer) > 0 {
		i -= len(m.Buyer)
		copy(dAtA[i:], m.Buyer)
		i = encodeVarint(dAtA, i, uint64(len(m.Buyer)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}
func (m *MsgBuyOption) Size() (n int) {
	if m == nil { return 0 }
	l := len(m.Buyer)
	if l > 0 { n += 1 + l + sov(uint64(l)) }
	if m.OptionID != 0 { n += 1 + sov(m.OptionID) }
	return n
}
func (m *MsgBuyOption) Unmarshal(dAtA []byte) error {
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
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgBuyOption: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for Buyer") }
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
			m.Buyer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 { return fmt.Errorf("proto: wrong wireType for OptionID") }
			m.OptionID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				m.OptionID |= uint64(b&0x7F) << shift
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
// MsgBuyOptionResponse (empty)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgBuyOptionResponse proto.InternalMessageInfo

func (m *MsgBuyOptionResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgBuyOptionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgBuyOptionResponse.Marshal(b, m, deterministic) }
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil { return nil, err }
	return b[:n], nil
}
func (m *MsgBuyOptionResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgBuyOptionResponse.Merge(m, src) }
func (m *MsgBuyOptionResponse) XXX_Size() int               { return m.Size() }
func (m *MsgBuyOptionResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgBuyOptionResponse.DiscardUnknown(m) }

func (m *MsgBuyOptionResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgBuyOptionResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgBuyOptionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgBuyOptionResponse) Size() (n int) { return 0 }
func (m *MsgBuyOptionResponse) Unmarshal(dAtA []byte) error {
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
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgBuyOptionResponse: illegal tag %d", fieldNum) }
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
// MsgExerciseOption (same layout as MsgBuyOption)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgExerciseOption proto.InternalMessageInfo

func (m *MsgExerciseOption) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgExerciseOption) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgExerciseOption.Marshal(b, m, deterministic) }
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil { return nil, err }
	return b[:n], nil
}
func (m *MsgExerciseOption) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgExerciseOption.Merge(m, src) }
func (m *MsgExerciseOption) XXX_Size() int               { return m.Size() }
func (m *MsgExerciseOption) XXX_DiscardUnknown()          { xxx_messageInfo_MsgExerciseOption.DiscardUnknown(m) }

func (m *MsgExerciseOption) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil { return nil, err }
	return dAtA[:n], nil
}
func (m *MsgExerciseOption) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgExerciseOption) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.OptionID != 0 {
		i = encodeVarint(dAtA, i, m.OptionID)
		i--
		dAtA[i] = 0x10
	}
	if len(m.Buyer) > 0 {
		i -= len(m.Buyer)
		copy(dAtA[i:], m.Buyer)
		i = encodeVarint(dAtA, i, uint64(len(m.Buyer)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}
func (m *MsgExerciseOption) Size() (n int) {
	if m == nil { return 0 }
	l := len(m.Buyer)
	if l > 0 { n += 1 + l + sov(uint64(l)) }
	if m.OptionID != 0 { n += 1 + sov(m.OptionID) }
	return n
}
func (m *MsgExerciseOption) Unmarshal(dAtA []byte) error {
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
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgExerciseOption: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for Buyer") }
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
			m.Buyer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 { return fmt.Errorf("proto: wrong wireType for OptionID") }
			m.OptionID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				m.OptionID |= uint64(b&0x7F) << shift
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
// MsgExerciseOptionResponse
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgExerciseOptionResponse proto.InternalMessageInfo

func (m *MsgExerciseOptionResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgExerciseOptionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgExerciseOptionResponse.Marshal(b, m, deterministic) }
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil { return nil, err }
	return b[:n], nil
}
func (m *MsgExerciseOptionResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgExerciseOptionResponse.Merge(m, src) }
func (m *MsgExerciseOptionResponse) XXX_Size() int               { return m.Size() }
func (m *MsgExerciseOptionResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgExerciseOptionResponse.DiscardUnknown(m) }

func (m *MsgExerciseOptionResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil { return nil, err }
	return dAtA[:n], nil
}
func (m *MsgExerciseOptionResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgExerciseOptionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.Payout.IsNil() {
		ii, err := marshalMathInt(dAtA, i, &m.Payout)
		if err != nil { return 0, err }
		i = ii
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}
func (m *MsgExerciseOptionResponse) Size() (n int) {
	if m == nil { return 0 }
	if !m.Payout.IsNil() { n += mathIntSize(&m.Payout) }
	return n
}
func (m *MsgExerciseOptionResponse) Unmarshal(dAtA []byte) error {
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
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgExerciseOptionResponse: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for Payout") }
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
			if err := json.Unmarshal([]byte(quoted), &m.Payout); err != nil { return err }
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
// MsgCancelOption (writer + option_id)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCancelOption proto.InternalMessageInfo

func (m *MsgCancelOption) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCancelOption) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgCancelOption.Marshal(b, m, deterministic) }
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil { return nil, err }
	return b[:n], nil
}
func (m *MsgCancelOption) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCancelOption.Merge(m, src) }
func (m *MsgCancelOption) XXX_Size() int               { return m.Size() }
func (m *MsgCancelOption) XXX_DiscardUnknown()          { xxx_messageInfo_MsgCancelOption.DiscardUnknown(m) }

func (m *MsgCancelOption) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil { return nil, err }
	return dAtA[:n], nil
}
func (m *MsgCancelOption) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgCancelOption) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.OptionID != 0 {
		i = encodeVarint(dAtA, i, m.OptionID)
		i--
		dAtA[i] = 0x10
	}
	if len(m.Writer) > 0 {
		i -= len(m.Writer)
		copy(dAtA[i:], m.Writer)
		i = encodeVarint(dAtA, i, uint64(len(m.Writer)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}
func (m *MsgCancelOption) Size() (n int) {
	if m == nil { return 0 }
	l := len(m.Writer)
	if l > 0 { n += 1 + l + sov(uint64(l)) }
	if m.OptionID != 0 { n += 1 + sov(m.OptionID) }
	return n
}
func (m *MsgCancelOption) Unmarshal(dAtA []byte) error {
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
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgCancelOption: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for Writer") }
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
			m.Writer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 { return fmt.Errorf("proto: wrong wireType for OptionID") }
			m.OptionID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflow }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				m.OptionID |= uint64(b&0x7F) << shift
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
// MsgCancelOptionResponse (empty)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCancelOptionResponse proto.InternalMessageInfo

func (m *MsgCancelOptionResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCancelOptionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgCancelOptionResponse.Marshal(b, m, deterministic) }
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil { return nil, err }
	return b[:n], nil
}
func (m *MsgCancelOptionResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCancelOptionResponse.Merge(m, src) }
func (m *MsgCancelOptionResponse) XXX_Size() int               { return m.Size() }
func (m *MsgCancelOptionResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgCancelOptionResponse.DiscardUnknown(m) }

func (m *MsgCancelOptionResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgCancelOptionResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgCancelOptionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgCancelOptionResponse) Size() (n int) { return 0 }
func (m *MsgCancelOptionResponse) Unmarshal(dAtA []byte) error {
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
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgCancelOptionResponse: illegal tag %d", fieldNum) }
		iNdEx = preIndex
		skippy, err := skip(dAtA[iNdEx:])
		if err != nil { return err }
		if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }
		if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }
		iNdEx += skippy
	}
	return nil
}
