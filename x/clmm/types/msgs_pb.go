package types

import (
	"encoding/json"
	"fmt"
	"io"
	"math/bits"

	proto "github.com/cosmos/gogoproto/proto"
)

// Reference imports to suppress errors
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = io.ErrUnexpectedEOF

// Descriptor methods
func (*MsgCreateCLPool) Descriptor() ([]byte, []int)             { return fileDescriptorCLMMTx, []int{0} }
func (*MsgCreateCLPoolResponse) Descriptor() ([]byte, []int)     { return fileDescriptorCLMMTx, []int{1} }
func (*MsgCreatePosition) Descriptor() ([]byte, []int)           { return fileDescriptorCLMMTx, []int{2} }
func (*MsgCreatePositionResponse) Descriptor() ([]byte, []int)   { return fileDescriptorCLMMTx, []int{3} }
func (*MsgAddLiquidity) Descriptor() ([]byte, []int)             { return fileDescriptorCLMMTx, []int{4} }
func (*MsgAddLiquidityResponse) Descriptor() ([]byte, []int)     { return fileDescriptorCLMMTx, []int{5} }
func (*MsgRemoveLiquidity) Descriptor() ([]byte, []int)          { return fileDescriptorCLMMTx, []int{6} }
func (*MsgRemoveLiquidityResponse) Descriptor() ([]byte, []int)  { return fileDescriptorCLMMTx, []int{7} }
func (*MsgCollectFees) Descriptor() ([]byte, []int)              { return fileDescriptorCLMMTx, []int{8} }
func (*MsgCollectFeesResponse) Descriptor() ([]byte, []int)      { return fileDescriptorCLMMTx, []int{9} }
func (*MsgCLSwap) Descriptor() ([]byte, []int)                   { return fileDescriptorCLMMTx, []int{10} }
func (*MsgCLSwapResponse) Descriptor() ([]byte, []int)           { return fileDescriptorCLMMTx, []int{11} }

func init() {
	proto.RegisterType((*MsgCreateCLPool)(nil), "syreen.clmm.MsgCreateCLPool")
	proto.RegisterType((*MsgCreateCLPoolResponse)(nil), "syreen.clmm.MsgCreateCLPoolResponse")
	proto.RegisterType((*MsgCreatePosition)(nil), "syreen.clmm.MsgCreatePosition")
	proto.RegisterType((*MsgCreatePositionResponse)(nil), "syreen.clmm.MsgCreatePositionResponse")
	proto.RegisterType((*MsgAddLiquidity)(nil), "syreen.clmm.MsgAddLiquidity")
	proto.RegisterType((*MsgAddLiquidityResponse)(nil), "syreen.clmm.MsgAddLiquidityResponse")
	proto.RegisterType((*MsgRemoveLiquidity)(nil), "syreen.clmm.MsgRemoveLiquidity")
	proto.RegisterType((*MsgRemoveLiquidityResponse)(nil), "syreen.clmm.MsgRemoveLiquidityResponse")
	proto.RegisterType((*MsgCollectFees)(nil), "syreen.clmm.MsgCollectFees")
	proto.RegisterType((*MsgCollectFeesResponse)(nil), "syreen.clmm.MsgCollectFeesResponse")
	proto.RegisterType((*MsgCLSwap)(nil), "syreen.clmm.MsgCLSwap")
	proto.RegisterType((*MsgCLSwapResponse)(nil), "syreen.clmm.MsgCLSwapResponse")
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

// marshalBytes is a helper that JSON-encodes complex types (math.Int, LegacyDec) to bytes
func marshalBytes(v interface{}) []byte {
	bz, _ := json.Marshal(v)
	return bz
}

// ---------------------------------------------------------------------------
// All messages use JSON-over-bytes for complex fields (math.Int, LegacyDec)
// and proper protobuf for string/varint fields
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// MsgCreateCLPool
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCreateCLPool proto.InternalMessageInfo

func (m *MsgCreateCLPool) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgCreateCLPool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateCLPool.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgCreateCLPool) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgCreateCLPool.Merge(m, src) }
func (m *MsgCreateCLPool) XXX_Size() int                { return m.Size() }
func (m *MsgCreateCLPool) XXX_DiscardUnknown()          { xxx_messageInfo_MsgCreateCLPool.DiscardUnknown(m) }

func (m *MsgCreateCLPool) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}
func (m *MsgCreateCLPool) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgCreateCLPool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	// field 6: initial_price (bytes)
	if !m.InitialPrice.IsNil() {
		bz := marshalBytes(m.InitialPrice)
		if len(bz) > 0 {
			i -= len(bz)
			copy(dAtA[i:], bz)
			i = encodeVarint(dAtA, i, uint64(len(bz)))
			i--
			dAtA[i] = 0x32
		}
	}
	// field 5: fee_rate (bytes)
	if !m.FeeRate.IsNil() {
		bz := marshalBytes(m.FeeRate)
		if len(bz) > 0 {
			i -= len(bz)
			copy(dAtA[i:], bz)
			i = encodeVarint(dAtA, i, uint64(len(bz)))
			i--
			dAtA[i] = 0x2a
		}
	}
	// field 4: tick_spacing (varint)
	if m.TickSpacing != 0 {
		i = encodeVarint(dAtA, i, uint64(m.TickSpacing))
		i--
		dAtA[i] = 0x20
	}
	// field 3: denom_b (string)
	if len(m.DenomB) > 0 {
		i -= len(m.DenomB)
		copy(dAtA[i:], m.DenomB)
		i = encodeVarint(dAtA, i, uint64(len(m.DenomB)))
		i--
		dAtA[i] = 0x1a
	}
	// field 2: denom_a (string)
	if len(m.DenomA) > 0 {
		i -= len(m.DenomA)
		copy(dAtA[i:], m.DenomA)
		i = encodeVarint(dAtA, i, uint64(len(m.DenomA)))
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
func (m *MsgCreateCLPool) Size() (n int) {
	if m == nil { return 0 }
	var l int
	l = len(m.Sender)
	if l > 0 { n += 1 + l + sov(uint64(l)) }
	l = len(m.DenomA)
	if l > 0 { n += 1 + l + sov(uint64(l)) }
	l = len(m.DenomB)
	if l > 0 { n += 1 + l + sov(uint64(l)) }
	if m.TickSpacing != 0 { n += 1 + sov(uint64(m.TickSpacing)) }
	if !m.FeeRate.IsNil() { bz := marshalBytes(m.FeeRate); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	if !m.InitialPrice.IsNil() { bz := marshalBytes(m.InitialPrice); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	return n
}
func (m *MsgCreateCLPool) Unmarshal(dAtA []byte) error {
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
		if wireType == 4 { return fmt.Errorf("proto: MsgCreateCLPool: wiretype end group for non-group") }
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgCreateCLPool: illegal tag %d (wire type %d)", fieldNum, wire) }
		switch fieldNum {
		case 1: // sender
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
			intStringLen := int(stringLen)
			if intStringLen < 0 { return ErrInvalidLength }
			postIndex := iNdEx + intStringLen
			if postIndex < 0 { return ErrInvalidLength }
			if postIndex > l { return io.ErrUnexpectedEOF }
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2: // denom_a
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field DenomA", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
			intStringLen := int(stringLen)
			if intStringLen < 0 { return ErrInvalidLength }
			postIndex := iNdEx + intStringLen
			if postIndex < 0 { return ErrInvalidLength }
			if postIndex > l { return io.ErrUnexpectedEOF }
			m.DenomA = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3: // denom_b
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field DenomB", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
			intStringLen := int(stringLen)
			if intStringLen < 0 { return ErrInvalidLength }
			postIndex := iNdEx + intStringLen
			if postIndex < 0 { return ErrInvalidLength }
			if postIndex > l { return io.ErrUnexpectedEOF }
			m.DenomB = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4: // tick_spacing
			if wireType != 0 { return fmt.Errorf("proto: wrong wireType = %d for field TickSpacing", wireType) }
			m.TickSpacing = 0
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.TickSpacing |= int64(b&0x7F) << shift; if b < 0x80 { break } }
		case 5: // fee_rate (bytes -> json)
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field FeeRate", wireType) }
			var byteLen int
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }
			if byteLen < 0 { return ErrInvalidLength }
			postIndex := iNdEx + byteLen
			if postIndex < 0 { return ErrInvalidLength }
			if postIndex > l { return io.ErrUnexpectedEOF }
			if err := json.Unmarshal(dAtA[iNdEx:postIndex], &m.FeeRate); err != nil { return err }
			iNdEx = postIndex
		case 6: // initial_price (bytes -> json)
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field InitialPrice", wireType) }
			var byteLen int
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }
			if byteLen < 0 { return ErrInvalidLength }
			postIndex := iNdEx + byteLen
			if postIndex < 0 { return ErrInvalidLength }
			if postIndex > l { return io.ErrUnexpectedEOF }
			if err := json.Unmarshal(dAtA[iNdEx:postIndex], &m.InitialPrice); err != nil { return err }
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
// MsgCreateCLPoolResponse
// ---------------------------------------------------------------------------
var xxx_messageInfo_MsgCreateCLPoolResponse proto.InternalMessageInfo
func (m *MsgCreateCLPoolResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateCLPoolResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgCreateCLPoolResponse.Marshal(b, m, deterministic) }
	b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil
}
func (m *MsgCreateCLPoolResponse) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgCreateCLPoolResponse.Merge(m, src) }
func (m *MsgCreateCLPoolResponse) XXX_Size() int                { return m.Size() }
func (m *MsgCreateCLPoolResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgCreateCLPoolResponse.DiscardUnknown(m) }

func (m *MsgCreateCLPoolResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil
}
func (m *MsgCreateCLPoolResponse) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgCreateCLPoolResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.PoolID != 0 { i = encodeVarint(dAtA, i, uint64(m.PoolID)); i--; dAtA[i] = 0x08 }
	return len(dAtA) - i, nil
}
func (m *MsgCreateCLPoolResponse) Size() (n int) {
	if m == nil { return 0 }
	if m.PoolID != 0 { n += 1 + sov(uint64(m.PoolID)) }
	return n
}
func (m *MsgCreateCLPoolResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3)
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgCreateCLPoolResponse: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if int(wire&0x7) != 0 { return fmt.Errorf("proto: wrong wireType for field PoolID") }
			m.PoolID = 0
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.PoolID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		default:
			iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }
			if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }
	return nil
}

// ---------------------------------------------------------------------------
// Generic bytes-field message patterns for remaining types
// Using JSON encoding for math.Int/LegacyDec fields over protobuf bytes wire
// ---------------------------------------------------------------------------

// MsgCreatePosition
var xxx_messageInfo_MsgCreatePosition proto.InternalMessageInfo
func (m *MsgCreatePosition) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreatePosition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgCreatePosition.Marshal(b, m, deterministic) }
	b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil
}
func (m *MsgCreatePosition) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgCreatePosition.Merge(m, src) }
func (m *MsgCreatePosition) XXX_Size() int                { return m.Size() }
func (m *MsgCreatePosition) XXX_DiscardUnknown()          { xxx_messageInfo_MsgCreatePosition.DiscardUnknown(m) }

func (m *MsgCreatePosition) Marshal() (dAtA []byte, err error) {
	size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil
}
func (m *MsgCreatePosition) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgCreatePosition) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	// f8: amount_1_min
	if !m.Amount1Min.IsNil() { bz := marshalBytes(m.Amount1Min); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x42 } }
	// f7: amount_0_min
	if !m.Amount0Min.IsNil() { bz := marshalBytes(m.Amount0Min); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x3a } }
	// f6: amount_1_desired
	if !m.Amount1Desired.IsNil() { bz := marshalBytes(m.Amount1Desired); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x32 } }
	// f5: amount_0_desired
	if !m.Amount0Desired.IsNil() { bz := marshalBytes(m.Amount0Desired); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x2a } }
	// f4: tick_upper
	if m.TickUpper != 0 { i = encodeVarint(dAtA, i, uint64(m.TickUpper)); i--; dAtA[i] = 0x20 }
	// f3: tick_lower
	if m.TickLower != 0 { i = encodeVarint(dAtA, i, uint64(m.TickLower)); i--; dAtA[i] = 0x18 }
	// f2: pool_id
	if m.PoolID != 0 { i = encodeVarint(dAtA, i, uint64(m.PoolID)); i--; dAtA[i] = 0x10 }
	// f1: sender
	if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarint(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgCreatePosition) Size() (n int) {
	if m == nil { return 0 }
	var l int
	l = len(m.Sender); if l > 0 { n += 1 + l + sov(uint64(l)) }
	if m.PoolID != 0 { n += 1 + sov(uint64(m.PoolID)) }
	if m.TickLower != 0 { n += 1 + sov(uint64(m.TickLower)) }
	if m.TickUpper != 0 { n += 1 + sov(uint64(m.TickUpper)) }
	if !m.Amount0Desired.IsNil() { bz := marshalBytes(m.Amount0Desired); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	if !m.Amount1Desired.IsNil() { bz := marshalBytes(m.Amount1Desired); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	if !m.Amount0Min.IsNil() { bz := marshalBytes(m.Amount0Min); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	if !m.Amount1Min.IsNil() { bz := marshalBytes(m.Amount1Min); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	return n
}
func (m *MsgCreatePosition) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7)
		if wireType == 4 { return fmt.Errorf("proto: MsgCreatePosition: wiretype end group") }
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgCreatePosition: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for Sender") }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
			postIndex := iNdEx + int(stringLen); if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }
			m.Sender = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2:
			if wireType != 0 { return fmt.Errorf("proto: wrong wireType for PoolID") }
			m.PoolID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.PoolID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		case 3:
			if wireType != 0 { return fmt.Errorf("proto: wrong wireType for TickLower") }
			m.TickLower = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.TickLower |= int64(b&0x7F) << shift; if b < 0x80 { break } }
		case 4:
			if wireType != 0 { return fmt.Errorf("proto: wrong wireType for TickUpper") }
			m.TickUpper = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.TickUpper |= int64(b&0x7F) << shift; if b < 0x80 { break } }
		case 5, 6, 7, 8: // bytes fields for math.Int
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for bytes field %d", fieldNum) }
			var byteLen int
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }
			if byteLen < 0 { return ErrInvalidLength }
			postIndex := iNdEx + byteLen; if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }
			switch fieldNum {
			case 5: if err := json.Unmarshal(dAtA[iNdEx:postIndex], &m.Amount0Desired); err != nil { return err }
			case 6: if err := json.Unmarshal(dAtA[iNdEx:postIndex], &m.Amount1Desired); err != nil { return err }
			case 7: if err := json.Unmarshal(dAtA[iNdEx:postIndex], &m.Amount0Min); err != nil { return err }
			case 8: if err := json.Unmarshal(dAtA[iNdEx:postIndex], &m.Amount1Min); err != nil { return err }
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }
			if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }
	return nil
}

// MsgCreatePositionResponse
var xxx_messageInfo_MsgCreatePositionResponse proto.InternalMessageInfo
func (m *MsgCreatePositionResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreatePositionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgCreatePositionResponse.Marshal(b, m, deterministic) }
	b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil
}
func (m *MsgCreatePositionResponse) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgCreatePositionResponse.Merge(m, src) }
func (m *MsgCreatePositionResponse) XXX_Size() int                { return m.Size() }
func (m *MsgCreatePositionResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgCreatePositionResponse.DiscardUnknown(m) }
func (m *MsgCreatePositionResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgCreatePositionResponse) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgCreatePositionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.Liquidity.IsNil() { bz := marshalBytes(m.Liquidity); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x22 } }
	if !m.Amount1.IsNil() { bz := marshalBytes(m.Amount1); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x1a } }
	if !m.Amount0.IsNil() { bz := marshalBytes(m.Amount0); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x12 } }
	if m.PositionID != 0 { i = encodeVarint(dAtA, i, uint64(m.PositionID)); i--; dAtA[i] = 0x08 }
	return len(dAtA) - i, nil
}
func (m *MsgCreatePositionResponse) Size() (n int) {
	if m == nil { return 0 }; var l int
	if m.PositionID != 0 { n += 1 + sov(uint64(m.PositionID)) }
	if !m.Amount0.IsNil() { bz := marshalBytes(m.Amount0); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	if !m.Amount1.IsNil() { bz := marshalBytes(m.Amount1); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	if !m.Liquidity.IsNil() { bz := marshalBytes(m.Liquidity); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	return n
}
func (m *MsgCreatePositionResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3)
		if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if int(wire&0x7) != 0 { return fmt.Errorf("proto: wrong wireType for PositionID") }
			m.PositionID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.PositionID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		case 2, 3, 4:
			if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType for bytes field") }
			var byteLen int
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }
			if byteLen < 0 { return ErrInvalidLength }; postIndex := iNdEx + byteLen; if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }
			switch fieldNum {
			case 2: _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.Amount0)
			case 3: _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.Amount1)
			case 4: _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.Liquidity)
			}; iNdEx = postIndex
		default:
			iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// ---------------------------------------------------------------------------
// MsgAddLiquidity (sender:string, position_id:uint64, amount_0_desired:bytes, amount_1_desired:bytes)
// ---------------------------------------------------------------------------
var xxx_messageInfo_MsgAddLiquidity proto.InternalMessageInfo
func (m *MsgAddLiquidity) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgAddLiquidity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgAddLiquidity.Marshal(b, m, deterministic) }
	b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil
}
func (m *MsgAddLiquidity) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgAddLiquidity.Merge(m, src) }
func (m *MsgAddLiquidity) XXX_Size() int                { return m.Size() }
func (m *MsgAddLiquidity) XXX_DiscardUnknown()          { xxx_messageInfo_MsgAddLiquidity.DiscardUnknown(m) }
func (m *MsgAddLiquidity) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgAddLiquidity) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgAddLiquidity) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.Amount1Desired.IsNil() { bz := marshalBytes(m.Amount1Desired); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x22 } }
	if !m.Amount0Desired.IsNil() { bz := marshalBytes(m.Amount0Desired); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x1a } }
	if m.PositionID != 0 { i = encodeVarint(dAtA, i, uint64(m.PositionID)); i--; dAtA[i] = 0x10 }
	if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarint(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgAddLiquidity) Size() (n int) {
	if m == nil { return 0 }; var l int
	l = len(m.Sender); if l > 0 { n += 1 + l + sov(uint64(l)) }
	if m.PositionID != 0 { n += 1 + sov(uint64(m.PositionID)) }
	if !m.Amount0Desired.IsNil() { bz := marshalBytes(m.Amount0Desired); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	if !m.Amount1Desired.IsNil() { bz := marshalBytes(m.Amount1Desired); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	return n
}
func (m *MsgAddLiquidity) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7)
		if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for Sender") }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
			postIndex := iNdEx + int(stringLen); if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; m.Sender = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2:
			if wireType != 0 { return fmt.Errorf("proto: wrong wireType for PositionID") }
			m.PositionID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.PositionID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		case 3, 4:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for bytes field") }
			var byteLen int
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }
			if byteLen < 0 { return ErrInvalidLength }; postIndex := iNdEx + byteLen; if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }
			switch fieldNum {
			case 3: _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.Amount0Desired)
			case 4: _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.Amount1Desired)
			}; iNdEx = postIndex
		default:
			iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgAddLiquidityResponse
var xxx_messageInfo_MsgAddLiquidityResponse proto.InternalMessageInfo
func (m *MsgAddLiquidityResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgAddLiquidityResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgAddLiquidityResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgAddLiquidityResponse) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgAddLiquidityResponse.Merge(m, src) }
func (m *MsgAddLiquidityResponse) XXX_Size() int                { return m.Size() }
func (m *MsgAddLiquidityResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgAddLiquidityResponse.DiscardUnknown(m) }
func (m *MsgAddLiquidityResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgAddLiquidityResponse) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgAddLiquidityResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.Liquidity.IsNil() { bz := marshalBytes(m.Liquidity); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x1a } }
	if !m.Amount1.IsNil() { bz := marshalBytes(m.Amount1); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x12 } }
	if !m.Amount0.IsNil() { bz := marshalBytes(m.Amount0); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x0a } }
	return len(dAtA) - i, nil
}
func (m *MsgAddLiquidityResponse) Size() (n int) {
	if m == nil { return 0 }; var l int
	if !m.Amount0.IsNil() { bz := marshalBytes(m.Amount0); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	if !m.Amount1.IsNil() { bz := marshalBytes(m.Amount1); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	if !m.Liquidity.IsNil() { bz := marshalBytes(m.Liquidity); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	return n
}
func (m *MsgAddLiquidityResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3)
		if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1, 2, 3:
			if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType for bytes field") }
			var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }
			if byteLen < 0 { return ErrInvalidLength }; postIndex := iNdEx + byteLen; if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }
			switch fieldNum {
			case 1: _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.Amount0)
			case 2: _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.Amount1)
			case 3: _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.Liquidity)
			}; iNdEx = postIndex
		default:
			iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgRemoveLiquidity (sender:string, position_id:uint64, liquidity_amount:bytes)
var xxx_messageInfo_MsgRemoveLiquidity proto.InternalMessageInfo
func (m *MsgRemoveLiquidity) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRemoveLiquidity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgRemoveLiquidity.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgRemoveLiquidity) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgRemoveLiquidity.Merge(m, src) }
func (m *MsgRemoveLiquidity) XXX_Size() int                { return m.Size() }
func (m *MsgRemoveLiquidity) XXX_DiscardUnknown()          { xxx_messageInfo_MsgRemoveLiquidity.DiscardUnknown(m) }
func (m *MsgRemoveLiquidity) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgRemoveLiquidity) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgRemoveLiquidity) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.LiquidityAmount.IsNil() { bz := marshalBytes(m.LiquidityAmount); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x1a } }
	if m.PositionID != 0 { i = encodeVarint(dAtA, i, uint64(m.PositionID)); i--; dAtA[i] = 0x10 }
	if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarint(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgRemoveLiquidity) Size() (n int) {
	if m == nil { return 0 }; var l int
	l = len(m.Sender); if l > 0 { n += 1 + l + sov(uint64(l)) }
	if m.PositionID != 0 { n += 1 + sov(uint64(m.PositionID)) }
	if !m.LiquidityAmount.IsNil() { bz := marshalBytes(m.LiquidityAmount); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	return n
}
func (m *MsgRemoveLiquidity) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7)
		if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for Sender") }
			var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
			postIndex := iNdEx + int(stringLen); if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; m.Sender = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2:
			if wireType != 0 { return fmt.Errorf("proto: wrong wireType for PositionID") }
			m.PositionID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.PositionID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		case 3:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for LiquidityAmount") }
			var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }
			if byteLen < 0 { return ErrInvalidLength }; postIndex := iNdEx + byteLen; if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }
			_ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.LiquidityAmount); iNdEx = postIndex
		default:
			iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgRemoveLiquidityResponse (amount_0:bytes, amount_1:bytes)
var xxx_messageInfo_MsgRemoveLiquidityResponse proto.InternalMessageInfo
func (m *MsgRemoveLiquidityResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRemoveLiquidityResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgRemoveLiquidityResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgRemoveLiquidityResponse) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgRemoveLiquidityResponse.Merge(m, src) }
func (m *MsgRemoveLiquidityResponse) XXX_Size() int                { return m.Size() }
func (m *MsgRemoveLiquidityResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgRemoveLiquidityResponse.DiscardUnknown(m) }
func (m *MsgRemoveLiquidityResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgRemoveLiquidityResponse) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgRemoveLiquidityResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.Amount1.IsNil() { bz := marshalBytes(m.Amount1); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x12 } }
	if !m.Amount0.IsNil() { bz := marshalBytes(m.Amount0); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x0a } }
	return len(dAtA) - i, nil
}
func (m *MsgRemoveLiquidityResponse) Size() (n int) {
	if m == nil { return 0 }; var l int
	if !m.Amount0.IsNil() { bz := marshalBytes(m.Amount0); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	if !m.Amount1.IsNil() { bz := marshalBytes(m.Amount1); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	return n
}
func (m *MsgRemoveLiquidityResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3)
		if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1, 2:
			if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType for bytes field") }
			var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }
			if byteLen < 0 { return ErrInvalidLength }; postIndex := iNdEx + byteLen; if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }
			switch fieldNum {
			case 1: _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.Amount0)
			case 2: _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.Amount1)
			}; iNdEx = postIndex
		default:
			iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgCollectFees (sender:string, position_id:uint64)
var xxx_messageInfo_MsgCollectFees proto.InternalMessageInfo
func (m *MsgCollectFees) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCollectFees) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgCollectFees.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgCollectFees) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgCollectFees.Merge(m, src) }
func (m *MsgCollectFees) XXX_Size() int                { return m.Size() }
func (m *MsgCollectFees) XXX_DiscardUnknown()          { xxx_messageInfo_MsgCollectFees.DiscardUnknown(m) }
func (m *MsgCollectFees) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgCollectFees) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgCollectFees) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.PositionID != 0 { i = encodeVarint(dAtA, i, uint64(m.PositionID)); i--; dAtA[i] = 0x10 }
	if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarint(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgCollectFees) Size() (n int) {
	if m == nil { return 0 }; var l int
	l = len(m.Sender); if l > 0 { n += 1 + l + sov(uint64(l)) }
	if m.PositionID != 0 { n += 1 + sov(uint64(m.PositionID)) }
	return n
}
func (m *MsgCollectFees) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7)
		if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for Sender") }
			var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
			postIndex := iNdEx + int(stringLen); if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; m.Sender = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2:
			if wireType != 0 { return fmt.Errorf("proto: wrong wireType for PositionID") }
			m.PositionID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.PositionID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		default:
			iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgCollectFeesResponse = MsgRemoveLiquidityResponse shape (amount_0, amount_1)
var xxx_messageInfo_MsgCollectFeesResponse proto.InternalMessageInfo
func (m *MsgCollectFeesResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCollectFeesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgCollectFeesResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgCollectFeesResponse) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgCollectFeesResponse.Merge(m, src) }
func (m *MsgCollectFeesResponse) XXX_Size() int                { return m.Size() }
func (m *MsgCollectFeesResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgCollectFeesResponse.DiscardUnknown(m) }
func (m *MsgCollectFeesResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgCollectFeesResponse) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgCollectFeesResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.Amount1.IsNil() { bz := marshalBytes(m.Amount1); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x12 } }
	if !m.Amount0.IsNil() { bz := marshalBytes(m.Amount0); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x0a } }
	return len(dAtA) - i, nil
}
func (m *MsgCollectFeesResponse) Size() (n int) {
	if m == nil { return 0 }; var l int
	if !m.Amount0.IsNil() { bz := marshalBytes(m.Amount0); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	if !m.Amount1.IsNil() { bz := marshalBytes(m.Amount1); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	return n
}
func (m *MsgCollectFeesResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3)
		if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1, 2:
			if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType for bytes field") }
			var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }
			if byteLen < 0 { return ErrInvalidLength }; postIndex := iNdEx + byteLen; if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }
			switch fieldNum {
			case 1: _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.Amount0)
			case 2: _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.Amount1)
			}; iNdEx = postIndex
		default:
			iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgCLSwap (sender:string, pool_id:uint64, denom_in:string, amount_in:bytes, min_amount_out:bytes)
var xxx_messageInfo_MsgCLSwap proto.InternalMessageInfo
func (m *MsgCLSwap) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCLSwap) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgCLSwap.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgCLSwap) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgCLSwap.Merge(m, src) }
func (m *MsgCLSwap) XXX_Size() int                { return m.Size() }
func (m *MsgCLSwap) XXX_DiscardUnknown()          { xxx_messageInfo_MsgCLSwap.DiscardUnknown(m) }
func (m *MsgCLSwap) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgCLSwap) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgCLSwap) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.MinAmountOut.IsNil() { bz := marshalBytes(m.MinAmountOut); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x2a } }
	if !m.AmountIn.IsNil() { bz := marshalBytes(m.AmountIn); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x22 } }
	if len(m.DenomIn) > 0 { i -= len(m.DenomIn); copy(dAtA[i:], m.DenomIn); i = encodeVarint(dAtA, i, uint64(len(m.DenomIn))); i--; dAtA[i] = 0x1a }
	if m.PoolID != 0 { i = encodeVarint(dAtA, i, uint64(m.PoolID)); i--; dAtA[i] = 0x10 }
	if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarint(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgCLSwap) Size() (n int) {
	if m == nil { return 0 }; var l int
	l = len(m.Sender); if l > 0 { n += 1 + l + sov(uint64(l)) }
	if m.PoolID != 0 { n += 1 + sov(uint64(m.PoolID)) }
	l = len(m.DenomIn); if l > 0 { n += 1 + l + sov(uint64(l)) }
	if !m.AmountIn.IsNil() { bz := marshalBytes(m.AmountIn); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	if !m.MinAmountOut.IsNil() { bz := marshalBytes(m.MinAmountOut); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	return n
}
func (m *MsgCLSwap) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7)
		if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for Sender") }
			var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
			postIndex := iNdEx + int(stringLen); if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; m.Sender = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2:
			if wireType != 0 { return fmt.Errorf("proto: wrong wireType for PoolID") }
			m.PoolID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.PoolID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		case 3:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for DenomIn") }
			var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
			postIndex := iNdEx + int(stringLen); if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; m.DenomIn = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 4, 5:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for bytes field") }
			var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }
			if byteLen < 0 { return ErrInvalidLength }; postIndex := iNdEx + byteLen; if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }
			switch fieldNum {
			case 4: _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.AmountIn)
			case 5: _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.MinAmountOut)
			}; iNdEx = postIndex
		default:
			iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgCLSwapResponse (amount_out:bytes)
var xxx_messageInfo_MsgCLSwapResponse proto.InternalMessageInfo
func (m *MsgCLSwapResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCLSwapResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgCLSwapResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgCLSwapResponse) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgCLSwapResponse.Merge(m, src) }
func (m *MsgCLSwapResponse) XXX_Size() int                { return m.Size() }
func (m *MsgCLSwapResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgCLSwapResponse.DiscardUnknown(m) }
func (m *MsgCLSwapResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgCLSwapResponse) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgCLSwapResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.AmountOut.IsNil() { bz := marshalBytes(m.AmountOut); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x0a } }
	return len(dAtA) - i, nil
}
func (m *MsgCLSwapResponse) Size() (n int) {
	if m == nil { return 0 }; var l int
	if !m.AmountOut.IsNil() { bz := marshalBytes(m.AmountOut); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	return n
}
func (m *MsgCLSwapResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3)
		if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType for AmountOut") }
			var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }
			if byteLen < 0 { return ErrInvalidLength }; postIndex := iNdEx + byteLen; if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }
			_ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.AmountOut); iNdEx = postIndex
		default:
			iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}
