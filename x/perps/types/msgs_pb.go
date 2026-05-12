package types

import (
	"encoding/json"
	"fmt"
	"io"
	"math/bits"

	proto "github.com/cosmos/gogoproto/proto"
)

var _ = proto.Marshal
var _ = fmt.Errorf
var _ = io.ErrUnexpectedEOF

func (*MsgOpenPosition) Descriptor() ([]byte, []int)          { return fileDescriptorPerpsTx, []int{0} }
func (*MsgOpenPositionResponse) Descriptor() ([]byte, []int)  { return fileDescriptorPerpsTx, []int{1} }
func (*MsgClosePosition) Descriptor() ([]byte, []int)         { return fileDescriptorPerpsTx, []int{2} }
func (*MsgClosePositionResponse) Descriptor() ([]byte, []int) { return fileDescriptorPerpsTx, []int{3} }
func (*MsgAddMargin) Descriptor() ([]byte, []int)             { return fileDescriptorPerpsTx, []int{4} }
func (*MsgAddMarginResponse) Descriptor() ([]byte, []int)     { return fileDescriptorPerpsTx, []int{5} }
func (*MsgRemoveMargin) Descriptor() ([]byte, []int)          { return fileDescriptorPerpsTx, []int{6} }
func (*MsgRemoveMarginResponse) Descriptor() ([]byte, []int)  { return fileDescriptorPerpsTx, []int{7} }
func (*MsgCreateMarket) Descriptor() ([]byte, []int)          { return fileDescriptorPerpsTx, []int{8} }
func (*MsgCreateMarketResponse) Descriptor() ([]byte, []int)  { return fileDescriptorPerpsTx, []int{9} }

func init() {
	proto.RegisterType((*MsgOpenPosition)(nil), "syreen.perps.MsgOpenPosition")
	proto.RegisterType((*MsgOpenPositionResponse)(nil), "syreen.perps.MsgOpenPositionResponse")
	proto.RegisterType((*MsgClosePosition)(nil), "syreen.perps.MsgClosePosition")
	proto.RegisterType((*MsgClosePositionResponse)(nil), "syreen.perps.MsgClosePositionResponse")
	proto.RegisterType((*MsgAddMargin)(nil), "syreen.perps.MsgAddMargin")
	proto.RegisterType((*MsgAddMarginResponse)(nil), "syreen.perps.MsgAddMarginResponse")
	proto.RegisterType((*MsgRemoveMargin)(nil), "syreen.perps.MsgRemoveMargin")
	proto.RegisterType((*MsgRemoveMarginResponse)(nil), "syreen.perps.MsgRemoveMarginResponse")
	proto.RegisterType((*MsgCreateMarket)(nil), "syreen.perps.MsgCreateMarket")
	proto.RegisterType((*MsgCreateMarketResponse)(nil), "syreen.perps.MsgCreateMarketResponse")
}

var (
	ErrIntOverflow   = fmt.Errorf("proto: integer overflow")
	ErrInvalidLength = fmt.Errorf("proto: negative length found during unmarshaling")
)

func encodeVarint(dAtA []byte, offset int, v uint64) int {
	offset -= sov(v); base := offset
	for v >= 1<<7 { dAtA[offset] = uint8(v&0x7f | 0x80); v >>= 7; offset++ }
	dAtA[offset] = uint8(v); return base
}
func sov(x uint64) int { return (bits.Len64(x|1) + 6) / 7 }
func skip(dAtA []byte) (int, error) {
	l := len(dAtA); iNdEx := 0; depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return 0, ErrIntOverflow }; if iNdEx >= l { return 0, io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= (uint64(b) & 0x7F) << shift; if b < 0x80 { break } }
		wireType := int(wire & 0x7)
		switch wireType {
		case 0: for shift := uint(0); ; shift += 7 { if shift >= 64 { return 0, ErrIntOverflow }; if iNdEx >= l { return 0, io.ErrUnexpectedEOF }; iNdEx++; if dAtA[iNdEx-1] < 0x80 { break } }
		case 1: iNdEx += 8
		case 2: var length int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return 0, ErrIntOverflow }; if iNdEx >= l { return 0, io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; length |= (int(b) & 0x7F) << shift; if b < 0x80 { break } }; if length < 0 { return 0, ErrInvalidLength }; iNdEx += length
		case 3: depth++
		case 4: if depth == 0 { return 0, fmt.Errorf("proto: unexpected end group") }; depth--
		case 5: iNdEx += 4
		default: return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 { return 0, ErrInvalidLength }
		if depth == 0 { return iNdEx, nil }
	}
	return 0, io.ErrUnexpectedEOF
}
func marshalBytes(v interface{}) []byte { bz, _ := json.Marshal(v); return bz }

// Helper: string+uint64+bytes msg pattern
func marshalSenderUint64Bytes3(dAtA []byte, sender string, id uint64, b3 []byte) (int, error) {
	i := len(dAtA)
	if len(b3) > 0 { i -= len(b3); copy(dAtA[i:], b3); i = encodeVarint(dAtA, i, uint64(len(b3))); i--; dAtA[i] = 0x1a }
	if id != 0 { i = encodeVarint(dAtA, i, id); i--; dAtA[i] = 0x10 }
	if len(sender) > 0 { i -= len(sender); copy(dAtA[i:], sender); i = encodeVarint(dAtA, i, uint64(len(sender))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}

// MsgOpenPosition
var xxx_messageInfo_MsgOpenPosition proto.InternalMessageInfo
func (m *MsgOpenPosition) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgOpenPosition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgOpenPosition.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgOpenPosition) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgOpenPosition.Merge(m, src) }
func (m *MsgOpenPosition) XXX_Size() int { return m.Size() }
func (m *MsgOpenPosition) XXX_DiscardUnknown() { xxx_messageInfo_MsgOpenPosition.DiscardUnknown(m) }
func (m *MsgOpenPosition) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgOpenPosition) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgOpenPosition) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.Leverage.IsNil() { bz := marshalBytes(m.Leverage); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x2a } }
	if !m.Margin.IsNil() { bz := marshalBytes(m.Margin); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x22 } }
	if len(m.Side) > 0 { s := string(m.Side); i -= len(s); copy(dAtA[i:], s); i = encodeVarint(dAtA, i, uint64(len(s))); i--; dAtA[i] = 0x1a }
	if m.MarketID != 0 { i = encodeVarint(dAtA, i, uint64(m.MarketID)); i--; dAtA[i] = 0x10 }
	if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarint(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgOpenPosition) Size() (n int) {
	if m == nil { return 0 }; var l int
	l = len(m.Sender); if l > 0 { n += 1 + l + sov(uint64(l)) }
	if m.MarketID != 0 { n += 1 + sov(uint64(m.MarketID)) }
	l = len(string(m.Side)); if l > 0 { n += 1 + l + sov(uint64(l)) }
	if !m.Margin.IsNil() { bz := marshalBytes(m.Margin); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	if !m.Leverage.IsNil() { bz := marshalBytes(m.Leverage); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	return n
}
func (m *MsgOpenPosition) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; postIndex := iNdEx + int(stringLen); if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; m.Sender = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.MarketID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.MarketID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		case 3: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; postIndex := iNdEx + int(stringLen); if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; m.Side = Side(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 4: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLength }; postIndex := iNdEx + byteLen; if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.Margin); iNdEx = postIndex
		case 5: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLength }; postIndex := iNdEx + byteLen; if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.Leverage); iNdEx = postIndex
		default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgOpenPositionResponse
var xxx_messageInfo_MsgOpenPositionResponse proto.InternalMessageInfo
func (m *MsgOpenPositionResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgOpenPositionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgOpenPositionResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgOpenPositionResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgOpenPositionResponse.Merge(m, src) }
func (m *MsgOpenPositionResponse) XXX_Size() int { return m.Size() }
func (m *MsgOpenPositionResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgOpenPositionResponse.DiscardUnknown(m) }
func (m *MsgOpenPositionResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgOpenPositionResponse) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgOpenPositionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.Fee.IsNil() { bz := marshalBytes(m.Fee); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x1a } }
	if !m.EntryPrice.IsNil() { bz := marshalBytes(m.EntryPrice); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x12 } }
	if !m.PositionSize.IsNil() { bz := marshalBytes(m.PositionSize); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x0a } }
	return len(dAtA) - i, nil
}
func (m *MsgOpenPositionResponse) Size() (n int) {
	if m == nil { return 0 }; var l int
	if !m.PositionSize.IsNil() { bz := marshalBytes(m.PositionSize); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	if !m.EntryPrice.IsNil() { bz := marshalBytes(m.EntryPrice); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	if !m.Fee.IsNil() { bz := marshalBytes(m.Fee); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	return n
}
func (m *MsgOpenPositionResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1, 2, 3: if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLength }; postIndex := iNdEx + byteLen; if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; switch fieldNum { case 1: _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.PositionSize); case 2: _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.EntryPrice); case 3: _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.Fee) }; iNdEx = postIndex
		default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgClosePosition (sender:string, market_id:uint64)
var xxx_messageInfo_MsgClosePosition proto.InternalMessageInfo
func (m *MsgClosePosition) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgClosePosition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgClosePosition.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgClosePosition) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgClosePosition.Merge(m, src) }
func (m *MsgClosePosition) XXX_Size() int { return m.Size() }
func (m *MsgClosePosition) XXX_DiscardUnknown() { xxx_messageInfo_MsgClosePosition.DiscardUnknown(m) }
func (m *MsgClosePosition) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgClosePosition) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgClosePosition) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.MarketID != 0 { i = encodeVarint(dAtA, i, uint64(m.MarketID)); i--; dAtA[i] = 0x10 }
	if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarint(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgClosePosition) Size() (n int) { if m == nil { return 0 }; var l int; l = len(m.Sender); if l > 0 { n += 1 + l + sov(uint64(l)) }; if m.MarketID != 0 { n += 1 + sov(uint64(m.MarketID)) }; return n }
func (m *MsgClosePosition) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; postIndex := iNdEx + int(stringLen); if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; m.Sender = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.MarketID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.MarketID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgClosePositionResponse (realized_pnl:bytes, payout:bytes)
var xxx_messageInfo_MsgClosePositionResponse proto.InternalMessageInfo
func (m *MsgClosePositionResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgClosePositionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgClosePositionResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgClosePositionResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgClosePositionResponse.Merge(m, src) }
func (m *MsgClosePositionResponse) XXX_Size() int { return m.Size() }
func (m *MsgClosePositionResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgClosePositionResponse.DiscardUnknown(m) }
func (m *MsgClosePositionResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgClosePositionResponse) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgClosePositionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.Payout.IsNil() { bz := marshalBytes(m.Payout); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x12 } }
	if !m.RealizedPnL.IsNil() { bz := marshalBytes(m.RealizedPnL); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x0a } }
	return len(dAtA) - i, nil
}
func (m *MsgClosePositionResponse) Size() (n int) { if m == nil { return 0 }; var l int; if !m.RealizedPnL.IsNil() { bz := marshalBytes(m.RealizedPnL); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }; if !m.Payout.IsNil() { bz := marshalBytes(m.Payout); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }; return n }
func (m *MsgClosePositionResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1, 2: if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLength }; postIndex := iNdEx + byteLen; if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; switch fieldNum { case 1: _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.RealizedPnL); case 2: _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.Payout) }; iNdEx = postIndex
		default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgAddMargin (sender, market_id, amount)
var xxx_messageInfo_MsgAddMargin proto.InternalMessageInfo
func (m *MsgAddMargin) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgAddMargin) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgAddMargin.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgAddMargin) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgAddMargin.Merge(m, src) }
func (m *MsgAddMargin) XXX_Size() int { return m.Size() }
func (m *MsgAddMargin) XXX_DiscardUnknown() { xxx_messageInfo_MsgAddMargin.DiscardUnknown(m) }
func (m *MsgAddMargin) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgAddMargin) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgAddMargin) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.Amount.IsNil() { bz := marshalBytes(m.Amount); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x1a } }
	if m.MarketID != 0 { i = encodeVarint(dAtA, i, uint64(m.MarketID)); i--; dAtA[i] = 0x10 }
	if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarint(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgAddMargin) Size() (n int) { if m == nil { return 0 }; var l int; l = len(m.Sender); if l > 0 { n += 1 + l + sov(uint64(l)) }; if m.MarketID != 0 { n += 1 + sov(uint64(m.MarketID)) }; if !m.Amount.IsNil() { bz := marshalBytes(m.Amount); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }; return n }
func (m *MsgAddMargin) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; postIndex := iNdEx + int(stringLen); if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; m.Sender = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.MarketID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.MarketID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		case 3: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLength }; postIndex := iNdEx + byteLen; if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.Amount); iNdEx = postIndex
		default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// Empty responses
var xxx_messageInfo_MsgAddMarginResponse proto.InternalMessageInfo
func (m *MsgAddMarginResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgAddMarginResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return nil, nil }
func (m *MsgAddMarginResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgAddMarginResponse.Merge(m, src) }
func (m *MsgAddMarginResponse) XXX_Size() int { return 0 }
func (m *MsgAddMarginResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgAddMarginResponse.DiscardUnknown(m) }
func (m *MsgAddMarginResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgAddMarginResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgAddMarginResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgAddMarginResponse) Size() (n int) { return 0 }
func (m *MsgAddMarginResponse) Unmarshal(dAtA []byte) error { l := len(dAtA); iNdEx := 0; for iNdEx < l { preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }; iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy }; return nil }

// MsgRemoveMargin = same as MsgAddMargin
var xxx_messageInfo_MsgRemoveMargin proto.InternalMessageInfo
func (m *MsgRemoveMargin) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRemoveMargin) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgRemoveMargin.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgRemoveMargin) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgRemoveMargin.Merge(m, src) }
func (m *MsgRemoveMargin) XXX_Size() int { return m.Size() }
func (m *MsgRemoveMargin) XXX_DiscardUnknown() { xxx_messageInfo_MsgRemoveMargin.DiscardUnknown(m) }
func (m *MsgRemoveMargin) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgRemoveMargin) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgRemoveMargin) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.Amount.IsNil() { bz := marshalBytes(m.Amount); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x1a } }
	if m.MarketID != 0 { i = encodeVarint(dAtA, i, uint64(m.MarketID)); i--; dAtA[i] = 0x10 }
	if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarint(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgRemoveMargin) Size() (n int) { if m == nil { return 0 }; var l int; l = len(m.Sender); if l > 0 { n += 1 + l + sov(uint64(l)) }; if m.MarketID != 0 { n += 1 + sov(uint64(m.MarketID)) }; if !m.Amount.IsNil() { bz := marshalBytes(m.Amount); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }; return n }
func (m *MsgRemoveMargin) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; postIndex := iNdEx + int(stringLen); if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; m.Sender = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.MarketID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.MarketID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		case 3: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLength }; postIndex := iNdEx + byteLen; if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.Amount); iNdEx = postIndex
		default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

var xxx_messageInfo_MsgRemoveMarginResponse proto.InternalMessageInfo
func (m *MsgRemoveMarginResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRemoveMarginResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return nil, nil }
func (m *MsgRemoveMarginResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgRemoveMarginResponse.Merge(m, src) }
func (m *MsgRemoveMarginResponse) XXX_Size() int { return 0 }
func (m *MsgRemoveMarginResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgRemoveMarginResponse.DiscardUnknown(m) }
func (m *MsgRemoveMarginResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgRemoveMarginResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgRemoveMarginResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgRemoveMarginResponse) Size() (n int) { return 0 }
func (m *MsgRemoveMarginResponse) Unmarshal(dAtA []byte) error { l := len(dAtA); iNdEx := 0; for iNdEx < l { preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }; iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy }; return nil }

// MsgCreateMarket (authority, base_denom, quote_denom, pool_id, max_leverage)
var xxx_messageInfo_MsgCreateMarket proto.InternalMessageInfo
func (m *MsgCreateMarket) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateMarket) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgCreateMarket.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgCreateMarket) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCreateMarket.Merge(m, src) }
func (m *MsgCreateMarket) XXX_Size() int { return m.Size() }
func (m *MsgCreateMarket) XXX_DiscardUnknown() { xxx_messageInfo_MsgCreateMarket.DiscardUnknown(m) }
func (m *MsgCreateMarket) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgCreateMarket) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgCreateMarket) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.MaxLeverage.IsNil() { bz := marshalBytes(m.MaxLeverage); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x2a } }
	if m.PoolID != 0 { i = encodeVarint(dAtA, i, uint64(m.PoolID)); i--; dAtA[i] = 0x20 }
	if len(m.QuoteDenom) > 0 { i -= len(m.QuoteDenom); copy(dAtA[i:], m.QuoteDenom); i = encodeVarint(dAtA, i, uint64(len(m.QuoteDenom))); i--; dAtA[i] = 0x1a }
	if len(m.BaseDenom) > 0 { i -= len(m.BaseDenom); copy(dAtA[i:], m.BaseDenom); i = encodeVarint(dAtA, i, uint64(len(m.BaseDenom))); i--; dAtA[i] = 0x12 }
	if len(m.Authority) > 0 { i -= len(m.Authority); copy(dAtA[i:], m.Authority); i = encodeVarint(dAtA, i, uint64(len(m.Authority))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgCreateMarket) Size() (n int) {
	if m == nil { return 0 }; var l int
	l = len(m.Authority); if l > 0 { n += 1 + l + sov(uint64(l)) }
	l = len(m.BaseDenom); if l > 0 { n += 1 + l + sov(uint64(l)) }
	l = len(m.QuoteDenom); if l > 0 { n += 1 + l + sov(uint64(l)) }
	if m.PoolID != 0 { n += 1 + sov(uint64(m.PoolID)) }
	if !m.MaxLeverage.IsNil() { bz := marshalBytes(m.MaxLeverage); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	return n
}
func (m *MsgCreateMarket) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; postIndex := iNdEx + int(stringLen); if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; m.Authority = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; postIndex := iNdEx + int(stringLen); if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; m.BaseDenom = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 3: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; postIndex := iNdEx + int(stringLen); if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; m.QuoteDenom = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 4: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.PoolID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.PoolID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		case 5: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLength }; postIndex := iNdEx + byteLen; if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.MaxLeverage); iNdEx = postIndex
		default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgCreateMarketResponse (market_id:uint64)
var xxx_messageInfo_MsgCreateMarketResponse proto.InternalMessageInfo
func (m *MsgCreateMarketResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateMarketResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgCreateMarketResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgCreateMarketResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCreateMarketResponse.Merge(m, src) }
func (m *MsgCreateMarketResponse) XXX_Size() int { return m.Size() }
func (m *MsgCreateMarketResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgCreateMarketResponse.DiscardUnknown(m) }
func (m *MsgCreateMarketResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgCreateMarketResponse) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgCreateMarketResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA); if m.MarketID != 0 { i = encodeVarint(dAtA, i, uint64(m.MarketID)); i--; dAtA[i] = 0x08 }; return len(dAtA) - i, nil
}
func (m *MsgCreateMarketResponse) Size() (n int) { if m == nil { return 0 }; if m.MarketID != 0 { n += 1 + sov(uint64(m.MarketID)) }; return n }
func (m *MsgCreateMarketResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if int(wire&0x7) != 0 { return fmt.Errorf("proto: wrong wireType") }; m.MarketID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.MarketID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}
