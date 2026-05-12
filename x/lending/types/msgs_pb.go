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

func (*MsgDeposit) Descriptor() ([]byte, []int)                   { return fileDescriptorLendingTx, []int{0} }
func (*MsgDepositResponse) Descriptor() ([]byte, []int)           { return fileDescriptorLendingTx, []int{1} }
func (*MsgWithdraw) Descriptor() ([]byte, []int)                  { return fileDescriptorLendingTx, []int{2} }
func (*MsgWithdrawResponse) Descriptor() ([]byte, []int)          { return fileDescriptorLendingTx, []int{3} }
func (*MsgBorrow) Descriptor() ([]byte, []int)                    { return fileDescriptorLendingTx, []int{4} }
func (*MsgBorrowResponse) Descriptor() ([]byte, []int)            { return fileDescriptorLendingTx, []int{5} }
func (*MsgRepay) Descriptor() ([]byte, []int)                     { return fileDescriptorLendingTx, []int{6} }
func (*MsgRepayResponse) Descriptor() ([]byte, []int)             { return fileDescriptorLendingTx, []int{7} }
func (*MsgLiquidate) Descriptor() ([]byte, []int)                 { return fileDescriptorLendingTx, []int{8} }
func (*MsgLiquidateResponse) Descriptor() ([]byte, []int)         { return fileDescriptorLendingTx, []int{9} }
func (*MsgCreateLendingPool) Descriptor() ([]byte, []int)         { return fileDescriptorLendingTx, []int{10} }
func (*MsgCreateLendingPoolResponse) Descriptor() ([]byte, []int) { return fileDescriptorLendingTx, []int{11} }

func init() {
	proto.RegisterType((*MsgDeposit)(nil), "syreen.lending.MsgDeposit")
	proto.RegisterType((*MsgDepositResponse)(nil), "syreen.lending.MsgDepositResponse")
	proto.RegisterType((*MsgWithdraw)(nil), "syreen.lending.MsgWithdraw")
	proto.RegisterType((*MsgWithdrawResponse)(nil), "syreen.lending.MsgWithdrawResponse")
	proto.RegisterType((*MsgBorrow)(nil), "syreen.lending.MsgBorrow")
	proto.RegisterType((*MsgBorrowResponse)(nil), "syreen.lending.MsgBorrowResponse")
	proto.RegisterType((*MsgRepay)(nil), "syreen.lending.MsgRepay")
	proto.RegisterType((*MsgRepayResponse)(nil), "syreen.lending.MsgRepayResponse")
	proto.RegisterType((*MsgLiquidate)(nil), "syreen.lending.MsgLiquidate")
	proto.RegisterType((*MsgLiquidateResponse)(nil), "syreen.lending.MsgLiquidateResponse")
	proto.RegisterType((*MsgCreateLendingPool)(nil), "syreen.lending.MsgCreateLendingPool")
	proto.RegisterType((*MsgCreateLendingPoolResponse)(nil), "syreen.lending.MsgCreateLendingPoolResponse")
}

var (
	ErrIntOverflow   = fmt.Errorf("proto: integer overflow")
	ErrInvalidLength = fmt.Errorf("proto: negative length found during unmarshaling")
)

func encodeVarint(dAtA []byte, offset int, v uint64) int {
	offset -= sov(v); base := offset; for v >= 1<<7 { dAtA[offset] = uint8(v&0x7f | 0x80); v >>= 7; offset++ }; dAtA[offset] = uint8(v); return base
}
func sov(x uint64) int { return (bits.Len64(x|1) + 6) / 7 }
func skip(dAtA []byte) (int, error) {
	l := len(dAtA); iNdEx := 0; depth := 0
	for iNdEx < l {
		var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return 0, ErrIntOverflow }; if iNdEx >= l { return 0, io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= (uint64(b) & 0x7F) << shift; if b < 0x80 { break } }
		switch int(wire & 0x7) {
		case 0: for shift := uint(0); ; shift += 7 { if shift >= 64 { return 0, ErrIntOverflow }; if iNdEx >= l { return 0, io.ErrUnexpectedEOF }; iNdEx++; if dAtA[iNdEx-1] < 0x80 { break } }
		case 1: iNdEx += 8
		case 2: var length int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return 0, ErrIntOverflow }; if iNdEx >= l { return 0, io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; length |= (int(b) & 0x7F) << shift; if b < 0x80 { break } }; if length < 0 { return 0, ErrInvalidLength }; iNdEx += length
		case 3: depth++
		case 4: if depth == 0 { return 0, fmt.Errorf("proto: unexpected end group") }; depth--
		case 5: iNdEx += 4
		default: return 0, fmt.Errorf("proto: illegal wireType %d", int(wire&0x7))
		}; if iNdEx < 0 { return 0, ErrInvalidLength }; if depth == 0 { return iNdEx, nil }
	}; return 0, io.ErrUnexpectedEOF
}
func marshalBytes(v interface{}) []byte { bz, _ := json.Marshal(v); return bz }

// Helper: standard string+uint64+bytes pattern marshal/unmarshal
func marshalSUB(dAtA []byte, s string, id uint64, bz []byte) (int, error) {
	i := len(dAtA)
	if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x1a }
	if id != 0 { i = encodeVarint(dAtA, i, id); i--; dAtA[i] = 0x10 }
	if len(s) > 0 { i -= len(s); copy(dAtA[i:], s); i = encodeVarint(dAtA, i, uint64(len(s))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}

// MsgDeposit (sender:string/f1, pool_id:uint64/f2, amount:bytes/f3)
var xxx_messageInfo_MsgDeposit proto.InternalMessageInfo
func (m *MsgDeposit) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgDeposit) XXX_Marshal(b []byte, det bool) ([]byte, error) { if det { return xxx_messageInfo_MsgDeposit.Marshal(b, m, det) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgDeposit) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgDeposit.Merge(m, src) }
func (m *MsgDeposit) XXX_Size() int { return m.Size() }
func (m *MsgDeposit) XXX_DiscardUnknown() { xxx_messageInfo_MsgDeposit.DiscardUnknown(m) }
func (m *MsgDeposit) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgDeposit) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgDeposit) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.Amount.IsNil() { bz := marshalBytes(m.Amount); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x1a } }
	if m.PoolID != 0 { i = encodeVarint(dAtA, i, uint64(m.PoolID)); i--; dAtA[i] = 0x10 }
	if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarint(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgDeposit) Size() (n int) { if m == nil { return 0 }; var l int; l = len(m.Sender); if l > 0 { n += 1 + l + sov(uint64(l)) }; if m.PoolID != 0 { n += 1 + sov(uint64(m.PoolID)) }; if !m.Amount.IsNil() { bz := marshalBytes(m.Amount); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }; return n }
func (m *MsgDeposit) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var sLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; sLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; p := iNdEx + int(sLen); if p < 0 || p > l { return io.ErrUnexpectedEOF }; m.Sender = string(dAtA[iNdEx:p]); iNdEx = p
		case 2: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.PoolID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.PoolID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		case 3: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var bLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; bLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if bLen < 0 { return ErrInvalidLength }; p := iNdEx + bLen; if p < 0 || p > l { return io.ErrUnexpectedEOF }; _ = json.Unmarshal(dAtA[iNdEx:p], &m.Amount); iNdEx = p
		default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// Single bytes-field response helper
func marshalOneBytesField(dAtA []byte, v interface{}, isNil bool) (int, error) {
	i := len(dAtA); if !isNil { bz := marshalBytes(v); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x0a } }; return len(dAtA) - i, nil
}
func sizeOneBytesField(v interface{}, isNil bool) int { if isNil { return 0 }; bz := marshalBytes(v); l := len(bz); if l > 0 { return 1 + l + sov(uint64(l)) }; return 0 }
func unmarshalOneBytesField(dAtA []byte, target interface{}) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType") }; var bLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; bLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if bLen < 0 { return ErrInvalidLength }; p := iNdEx + bLen; if p < 0 || p > l { return io.ErrUnexpectedEOF }; _ = json.Unmarshal(dAtA[iNdEx:p], target); iNdEx = p
		default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgDepositResponse
var xxx_messageInfo_MsgDepositResponse proto.InternalMessageInfo
func (m *MsgDepositResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgDepositResponse) XXX_Marshal(b []byte, det bool) ([]byte, error) { if det { return xxx_messageInfo_MsgDepositResponse.Marshal(b, m, det) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgDepositResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgDepositResponse.Merge(m, src) }
func (m *MsgDepositResponse) XXX_Size() int { return m.Size() }
func (m *MsgDepositResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgDepositResponse.DiscardUnknown(m) }
func (m *MsgDepositResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgDepositResponse) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgDepositResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return marshalOneBytesField(dAtA, m.InterestEarned, m.InterestEarned.IsNil()) }
func (m *MsgDepositResponse) Size() (n int) { if m == nil { return 0 }; return sizeOneBytesField(m.InterestEarned, m.InterestEarned.IsNil()) }
func (m *MsgDepositResponse) Unmarshal(dAtA []byte) error { return unmarshalOneBytesField(dAtA, &m.InterestEarned) }

// MsgWithdraw = same as MsgDeposit
var xxx_messageInfo_MsgWithdraw proto.InternalMessageInfo
func (m *MsgWithdraw) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgWithdraw) XXX_Marshal(b []byte, det bool) ([]byte, error) { if det { return xxx_messageInfo_MsgWithdraw.Marshal(b, m, det) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgWithdraw) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgWithdraw.Merge(m, src) }
func (m *MsgWithdraw) XXX_Size() int { return m.Size() }
func (m *MsgWithdraw) XXX_DiscardUnknown() { xxx_messageInfo_MsgWithdraw.DiscardUnknown(m) }
func (m *MsgWithdraw) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgWithdraw) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgWithdraw) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.Amount.IsNil() { bz := marshalBytes(m.Amount); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x1a } }
	if m.PoolID != 0 { i = encodeVarint(dAtA, i, uint64(m.PoolID)); i--; dAtA[i] = 0x10 }
	if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarint(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgWithdraw) Size() (n int) { if m == nil { return 0 }; var l int; l = len(m.Sender); if l > 0 { n += 1 + l + sov(uint64(l)) }; if m.PoolID != 0 { n += 1 + sov(uint64(m.PoolID)) }; if !m.Amount.IsNil() { bz := marshalBytes(m.Amount); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }; return n }
func (m *MsgWithdraw) Unmarshal(dAtA []byte) error { tmp := MsgDeposit{}; if err := tmp.Unmarshal(dAtA); err != nil { return err }; m.Sender = tmp.Sender; m.PoolID = tmp.PoolID; m.Amount = tmp.Amount; return nil }

var xxx_messageInfo_MsgWithdrawResponse proto.InternalMessageInfo
func (m *MsgWithdrawResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgWithdrawResponse) XXX_Marshal(b []byte, det bool) ([]byte, error) { if det { return xxx_messageInfo_MsgWithdrawResponse.Marshal(b, m, det) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgWithdrawResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgWithdrawResponse.Merge(m, src) }
func (m *MsgWithdrawResponse) XXX_Size() int { return m.Size() }
func (m *MsgWithdrawResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgWithdrawResponse.DiscardUnknown(m) }
func (m *MsgWithdrawResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgWithdrawResponse) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgWithdrawResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return marshalOneBytesField(dAtA, m.InterestEarned, m.InterestEarned.IsNil()) }
func (m *MsgWithdrawResponse) Size() (n int) { if m == nil { return 0 }; return sizeOneBytesField(m.InterestEarned, m.InterestEarned.IsNil()) }
func (m *MsgWithdrawResponse) Unmarshal(dAtA []byte) error { return unmarshalOneBytesField(dAtA, &m.InterestEarned) }

// MsgBorrow (sender:s/1, borrow_pool_id:u64/2, amount:b/3, collateral_pool_id:u64/4, collateral_amount:b/5)
var xxx_messageInfo_MsgBorrow proto.InternalMessageInfo
func (m *MsgBorrow) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgBorrow) XXX_Marshal(b []byte, det bool) ([]byte, error) { if det { return xxx_messageInfo_MsgBorrow.Marshal(b, m, det) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgBorrow) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgBorrow.Merge(m, src) }
func (m *MsgBorrow) XXX_Size() int { return m.Size() }
func (m *MsgBorrow) XXX_DiscardUnknown() { xxx_messageInfo_MsgBorrow.DiscardUnknown(m) }
func (m *MsgBorrow) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgBorrow) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgBorrow) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.CollateralAmount.IsNil() { bz := marshalBytes(m.CollateralAmount); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x2a } }
	if m.CollateralPoolID != 0 { i = encodeVarint(dAtA, i, uint64(m.CollateralPoolID)); i--; dAtA[i] = 0x20 }
	if !m.Amount.IsNil() { bz := marshalBytes(m.Amount); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x1a } }
	if m.BorrowPoolID != 0 { i = encodeVarint(dAtA, i, uint64(m.BorrowPoolID)); i--; dAtA[i] = 0x10 }
	if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarint(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgBorrow) Size() (n int) {
	if m == nil { return 0 }; var l int; l = len(m.Sender); if l > 0 { n += 1 + l + sov(uint64(l)) }; if m.BorrowPoolID != 0 { n += 1 + sov(uint64(m.BorrowPoolID)) }; if !m.Amount.IsNil() { bz := marshalBytes(m.Amount); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }; if m.CollateralPoolID != 0 { n += 1 + sov(uint64(m.CollateralPoolID)) }; if !m.CollateralAmount.IsNil() { bz := marshalBytes(m.CollateralAmount); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }; return n
}
func (m *MsgBorrow) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var sLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; sLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; p := iNdEx + int(sLen); if p < 0 || p > l { return io.ErrUnexpectedEOF }; m.Sender = string(dAtA[iNdEx:p]); iNdEx = p
		case 2: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.BorrowPoolID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.BorrowPoolID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		case 3: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var bLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; bLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if bLen < 0 { return ErrInvalidLength }; p := iNdEx + bLen; if p < 0 || p > l { return io.ErrUnexpectedEOF }; _ = json.Unmarshal(dAtA[iNdEx:p], &m.Amount); iNdEx = p
		case 4: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.CollateralPoolID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.CollateralPoolID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		case 5: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var bLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; bLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if bLen < 0 { return ErrInvalidLength }; p := iNdEx + bLen; if p < 0 || p > l { return io.ErrUnexpectedEOF }; _ = json.Unmarshal(dAtA[iNdEx:p], &m.CollateralAmount); iNdEx = p
		default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgBorrowResponse (borrow_id:uint64/f1)
var xxx_messageInfo_MsgBorrowResponse proto.InternalMessageInfo
func (m *MsgBorrowResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgBorrowResponse) XXX_Marshal(b []byte, det bool) ([]byte, error) { if det { return xxx_messageInfo_MsgBorrowResponse.Marshal(b, m, det) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgBorrowResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgBorrowResponse.Merge(m, src) }
func (m *MsgBorrowResponse) XXX_Size() int { return m.Size() }
func (m *MsgBorrowResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgBorrowResponse.DiscardUnknown(m) }
func (m *MsgBorrowResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgBorrowResponse) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgBorrowResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { i := len(dAtA); if m.BorrowID != 0 { i = encodeVarint(dAtA, i, uint64(m.BorrowID)); i--; dAtA[i] = 0x08 }; return len(dAtA) - i, nil }
func (m *MsgBorrowResponse) Size() (n int) { if m == nil { return 0 }; if m.BorrowID != 0 { n += 1 + sov(uint64(m.BorrowID)) }; return n }
func (m *MsgBorrowResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l { preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }; switch fieldNum { case 1: if int(wire&0x7) != 0 { return fmt.Errorf("proto: wrong wireType") }; m.BorrowID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.BorrowID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy } }; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgRepay = same shape as MsgDeposit (sender, borrow_id, amount)
var xxx_messageInfo_MsgRepay proto.InternalMessageInfo
func (m *MsgRepay) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRepay) XXX_Marshal(b []byte, det bool) ([]byte, error) { if det { return xxx_messageInfo_MsgRepay.Marshal(b, m, det) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgRepay) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgRepay.Merge(m, src) }
func (m *MsgRepay) XXX_Size() int { return m.Size() }
func (m *MsgRepay) XXX_DiscardUnknown() { xxx_messageInfo_MsgRepay.DiscardUnknown(m) }
func (m *MsgRepay) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgRepay) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgRepay) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.Amount.IsNil() { bz := marshalBytes(m.Amount); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x1a } }
	if m.BorrowID != 0 { i = encodeVarint(dAtA, i, uint64(m.BorrowID)); i--; dAtA[i] = 0x10 }
	if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarint(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgRepay) Size() (n int) { if m == nil { return 0 }; var l int; l = len(m.Sender); if l > 0 { n += 1 + l + sov(uint64(l)) }; if m.BorrowID != 0 { n += 1 + sov(uint64(m.BorrowID)) }; if !m.Amount.IsNil() { bz := marshalBytes(m.Amount); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }; return n }
func (m *MsgRepay) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var sLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; sLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; p := iNdEx + int(sLen); if p < 0 || p > l { return io.ErrUnexpectedEOF }; m.Sender = string(dAtA[iNdEx:p]); iNdEx = p
		case 2: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.BorrowID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.BorrowID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		case 3: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var bLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; bLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if bLen < 0 { return ErrInvalidLength }; p := iNdEx + bLen; if p < 0 || p > l { return io.ErrUnexpectedEOF }; _ = json.Unmarshal(dAtA[iNdEx:p], &m.Amount); iNdEx = p
		default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgRepayResponse (interest_paid:b/1, collateral_returned:b/2)
var xxx_messageInfo_MsgRepayResponse proto.InternalMessageInfo
func (m *MsgRepayResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRepayResponse) XXX_Marshal(b []byte, det bool) ([]byte, error) { if det { return xxx_messageInfo_MsgRepayResponse.Marshal(b, m, det) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgRepayResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgRepayResponse.Merge(m, src) }
func (m *MsgRepayResponse) XXX_Size() int { return m.Size() }
func (m *MsgRepayResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgRepayResponse.DiscardUnknown(m) }
func (m *MsgRepayResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgRepayResponse) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgRepayResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.CollateralReturned.IsNil() { bz := marshalBytes(m.CollateralReturned); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x12 } }
	if !m.InterestPaid.IsNil() { bz := marshalBytes(m.InterestPaid); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x0a } }
	return len(dAtA) - i, nil
}
func (m *MsgRepayResponse) Size() (n int) { if m == nil { return 0 }; var l int; if !m.InterestPaid.IsNil() { bz := marshalBytes(m.InterestPaid); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }; if !m.CollateralReturned.IsNil() { bz := marshalBytes(m.CollateralReturned); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }; return n }
func (m *MsgRepayResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1, 2: if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType") }; var bLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; bLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if bLen < 0 { return ErrInvalidLength }; p := iNdEx + bLen; if p < 0 || p > l { return io.ErrUnexpectedEOF }; switch fieldNum { case 1: _ = json.Unmarshal(dAtA[iNdEx:p], &m.InterestPaid); case 2: _ = json.Unmarshal(dAtA[iNdEx:p], &m.CollateralReturned) }; iNdEx = p
		default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgLiquidate (liquidator:string/1, borrow_id:uint64/2)
var xxx_messageInfo_MsgLiquidate proto.InternalMessageInfo
func (m *MsgLiquidate) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgLiquidate) XXX_Marshal(b []byte, det bool) ([]byte, error) { if det { return xxx_messageInfo_MsgLiquidate.Marshal(b, m, det) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgLiquidate) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgLiquidate.Merge(m, src) }
func (m *MsgLiquidate) XXX_Size() int { return m.Size() }
func (m *MsgLiquidate) XXX_DiscardUnknown() { xxx_messageInfo_MsgLiquidate.DiscardUnknown(m) }
func (m *MsgLiquidate) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgLiquidate) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgLiquidate) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA); if m.BorrowID != 0 { i = encodeVarint(dAtA, i, uint64(m.BorrowID)); i--; dAtA[i] = 0x10 }; if len(m.Liquidator) > 0 { i -= len(m.Liquidator); copy(dAtA[i:], m.Liquidator); i = encodeVarint(dAtA, i, uint64(len(m.Liquidator))); i--; dAtA[i] = 0x0a }; return len(dAtA) - i, nil
}
func (m *MsgLiquidate) Size() (n int) { if m == nil { return 0 }; var l int; l = len(m.Liquidator); if l > 0 { n += 1 + l + sov(uint64(l)) }; if m.BorrowID != 0 { n += 1 + sov(uint64(m.BorrowID)) }; return n }
func (m *MsgLiquidate) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var sLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; sLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; p := iNdEx + int(sLen); if p < 0 || p > l { return io.ErrUnexpectedEOF }; m.Liquidator = string(dAtA[iNdEx:p]); iNdEx = p
		case 2: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.BorrowID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.BorrowID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgLiquidateResponse = same shape as MsgRepayResponse
var xxx_messageInfo_MsgLiquidateResponse proto.InternalMessageInfo
func (m *MsgLiquidateResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgLiquidateResponse) XXX_Marshal(b []byte, det bool) ([]byte, error) { if det { return xxx_messageInfo_MsgLiquidateResponse.Marshal(b, m, det) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgLiquidateResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgLiquidateResponse.Merge(m, src) }
func (m *MsgLiquidateResponse) XXX_Size() int { return m.Size() }
func (m *MsgLiquidateResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgLiquidateResponse.DiscardUnknown(m) }
func (m *MsgLiquidateResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgLiquidateResponse) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgLiquidateResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.DebtRepaid.IsNil() { bz := marshalBytes(m.DebtRepaid); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x12 } }
	if !m.CollateralSeized.IsNil() { bz := marshalBytes(m.CollateralSeized); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x0a } }
	return len(dAtA) - i, nil
}
func (m *MsgLiquidateResponse) Size() (n int) { if m == nil { return 0 }; var l int; if !m.CollateralSeized.IsNil() { bz := marshalBytes(m.CollateralSeized); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }; if !m.DebtRepaid.IsNil() { bz := marshalBytes(m.DebtRepaid); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }; return n }
func (m *MsgLiquidateResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1, 2: if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType") }; var bLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; bLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if bLen < 0 { return ErrInvalidLength }; p := iNdEx + bLen; if p < 0 || p > l { return io.ErrUnexpectedEOF }; switch fieldNum { case 1: _ = json.Unmarshal(dAtA[iNdEx:p], &m.CollateralSeized); case 2: _ = json.Unmarshal(dAtA[iNdEx:p], &m.DebtRepaid) }; iNdEx = p
		default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgCreateLendingPool (authority:s/1, denom:s/2, collateral_factor:b/3, dex_pool_id:u64/4, price_denom:s/5)
var xxx_messageInfo_MsgCreateLendingPool proto.InternalMessageInfo
func (m *MsgCreateLendingPool) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateLendingPool) XXX_Marshal(b []byte, det bool) ([]byte, error) { if det { return xxx_messageInfo_MsgCreateLendingPool.Marshal(b, m, det) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgCreateLendingPool) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCreateLendingPool.Merge(m, src) }
func (m *MsgCreateLendingPool) XXX_Size() int { return m.Size() }
func (m *MsgCreateLendingPool) XXX_DiscardUnknown() { xxx_messageInfo_MsgCreateLendingPool.DiscardUnknown(m) }
func (m *MsgCreateLendingPool) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgCreateLendingPool) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgCreateLendingPool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.PriceDenom) > 0 { i -= len(m.PriceDenom); copy(dAtA[i:], m.PriceDenom); i = encodeVarint(dAtA, i, uint64(len(m.PriceDenom))); i--; dAtA[i] = 0x2a }
	if m.DexPoolID != 0 { i = encodeVarint(dAtA, i, uint64(m.DexPoolID)); i--; dAtA[i] = 0x20 }
	if !m.CollateralFactor.IsNil() { bz := marshalBytes(m.CollateralFactor); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x1a } }
	if len(m.Denom) > 0 { i -= len(m.Denom); copy(dAtA[i:], m.Denom); i = encodeVarint(dAtA, i, uint64(len(m.Denom))); i--; dAtA[i] = 0x12 }
	if len(m.Authority) > 0 { i -= len(m.Authority); copy(dAtA[i:], m.Authority); i = encodeVarint(dAtA, i, uint64(len(m.Authority))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgCreateLendingPool) Size() (n int) {
	if m == nil { return 0 }; var l int; l = len(m.Authority); if l > 0 { n += 1 + l + sov(uint64(l)) }; l = len(m.Denom); if l > 0 { n += 1 + l + sov(uint64(l)) }; if !m.CollateralFactor.IsNil() { bz := marshalBytes(m.CollateralFactor); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }; if m.DexPoolID != 0 { n += 1 + sov(uint64(m.DexPoolID)) }; l = len(m.PriceDenom); if l > 0 { n += 1 + l + sov(uint64(l)) }; return n
}
func (m *MsgCreateLendingPool) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var sLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; sLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; p := iNdEx + int(sLen); if p < 0 || p > l { return io.ErrUnexpectedEOF }; m.Authority = string(dAtA[iNdEx:p]); iNdEx = p
		case 2: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var sLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; sLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; p := iNdEx + int(sLen); if p < 0 || p > l { return io.ErrUnexpectedEOF }; m.Denom = string(dAtA[iNdEx:p]); iNdEx = p
		case 3: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var bLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; bLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if bLen < 0 { return ErrInvalidLength }; p := iNdEx + bLen; if p < 0 || p > l { return io.ErrUnexpectedEOF }; _ = json.Unmarshal(dAtA[iNdEx:p], &m.CollateralFactor); iNdEx = p
		case 4: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.DexPoolID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.DexPoolID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		case 5: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var sLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; sLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; p := iNdEx + int(sLen); if p < 0 || p > l { return io.ErrUnexpectedEOF }; m.PriceDenom = string(dAtA[iNdEx:p]); iNdEx = p
		default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgCreateLendingPoolResponse (pool_id:uint64/f1) = same as MsgBorrowResponse
var xxx_messageInfo_MsgCreateLendingPoolResponse proto.InternalMessageInfo
func (m *MsgCreateLendingPoolResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateLendingPoolResponse) XXX_Marshal(b []byte, det bool) ([]byte, error) { if det { return xxx_messageInfo_MsgCreateLendingPoolResponse.Marshal(b, m, det) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgCreateLendingPoolResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCreateLendingPoolResponse.Merge(m, src) }
func (m *MsgCreateLendingPoolResponse) XXX_Size() int { return m.Size() }
func (m *MsgCreateLendingPoolResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgCreateLendingPoolResponse.DiscardUnknown(m) }
func (m *MsgCreateLendingPoolResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgCreateLendingPoolResponse) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgCreateLendingPoolResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { i := len(dAtA); if m.PoolID != 0 { i = encodeVarint(dAtA, i, uint64(m.PoolID)); i--; dAtA[i] = 0x08 }; return len(dAtA) - i, nil }
func (m *MsgCreateLendingPoolResponse) Size() (n int) { if m == nil { return 0 }; if m.PoolID != 0 { n += 1 + sov(uint64(m.PoolID)) }; return n }
func (m *MsgCreateLendingPoolResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l { preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }; switch fieldNum { case 1: if int(wire&0x7) != 0 { return fmt.Errorf("proto: wrong wireType") }; m.PoolID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.PoolID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy } }; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}
