package types

import (
	"fmt"
	"io"
	"math/bits"

	proto "github.com/cosmos/gogoproto/proto"
)

var _ = proto.Marshal
var _ = fmt.Errorf
var _ = io.ErrUnexpectedEOF

var fileDescriptorLaunchpadTx []byte

// Descriptor methods
func (*MsgCreateLaunch) Descriptor() ([]byte, []int)             { return fileDescriptorLaunchpadTx, []int{0} }
func (*MsgCreateLaunchResponse) Descriptor() ([]byte, []int)     { return fileDescriptorLaunchpadTx, []int{1} }
func (*MsgContribute) Descriptor() ([]byte, []int)               { return fileDescriptorLaunchpadTx, []int{2} }
func (*MsgContributeResponse) Descriptor() ([]byte, []int)       { return fileDescriptorLaunchpadTx, []int{3} }
func (*MsgClaimTokens) Descriptor() ([]byte, []int)              { return fileDescriptorLaunchpadTx, []int{4} }
func (*MsgClaimTokensResponse) Descriptor() ([]byte, []int)      { return fileDescriptorLaunchpadTx, []int{5} }
func (*MsgClaimRefund) Descriptor() ([]byte, []int)              { return fileDescriptorLaunchpadTx, []int{6} }
func (*MsgClaimRefundResponse) Descriptor() ([]byte, []int)      { return fileDescriptorLaunchpadTx, []int{7} }
func (*MsgFinalizeLaunch) Descriptor() ([]byte, []int)           { return fileDescriptorLaunchpadTx, []int{8} }
func (*MsgFinalizeLaunchResponse) Descriptor() ([]byte, []int)   { return fileDescriptorLaunchpadTx, []int{9} }

func init() {
	proto.RegisterType((*MsgCreateLaunch)(nil), "syreen.launchpad.MsgCreateLaunch")
	proto.RegisterType((*MsgCreateLaunchResponse)(nil), "syreen.launchpad.MsgCreateLaunchResponse")
	proto.RegisterType((*MsgContribute)(nil), "syreen.launchpad.MsgContribute")
	proto.RegisterType((*MsgContributeResponse)(nil), "syreen.launchpad.MsgContributeResponse")
	proto.RegisterType((*MsgClaimTokens)(nil), "syreen.launchpad.MsgClaimTokens")
	proto.RegisterType((*MsgClaimTokensResponse)(nil), "syreen.launchpad.MsgClaimTokensResponse")
	proto.RegisterType((*MsgClaimRefund)(nil), "syreen.launchpad.MsgClaimRefund")
	proto.RegisterType((*MsgClaimRefundResponse)(nil), "syreen.launchpad.MsgClaimRefundResponse")
	proto.RegisterType((*MsgFinalizeLaunch)(nil), "syreen.launchpad.MsgFinalizeLaunch")
	proto.RegisterType((*MsgFinalizeLaunchResponse)(nil), "syreen.launchpad.MsgFinalizeLaunchResponse")
}

// ---------------------------------------------------------------------------
// MsgCreateLaunch — 12 fields
// ---------------------------------------------------------------------------
var xxx_messageInfo_MsgCreateLaunch proto.InternalMessageInfo
func (m *MsgCreateLaunch) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateLaunch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgCreateLaunch.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgCreateLaunch) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCreateLaunch.Merge(m, src) }
func (m *MsgCreateLaunch) XXX_Size() int { return m.Size() }
func (m *MsgCreateLaunch) XXX_DiscardUnknown() { xxx_messageInfo_MsgCreateLaunch.DiscardUnknown(m) }

func (m *MsgCreateLaunch) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgCreateLaunch) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgCreateLaunch) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	// f12: tge_percent (varint)
	if m.TGEPercent != 0 { i = encodeVarintLaunchpad(dAtA, i, uint64(m.TGEPercent)); i--; dAtA[i] = 0x60 }
	// f11: vesting_blocks (varint)
	if m.VestingBlocks != 0 { i = encodeVarintLaunchpad(dAtA, i, uint64(m.VestingBlocks)); i--; dAtA[i] = 0x58 }
	// f10: end_block (varint)
	if m.EndBlock != 0 { i = encodeVarintLaunchpad(dAtA, i, uint64(m.EndBlock)); i--; dAtA[i] = 0x50 }
	// f9: start_block (varint)
	if m.StartBlock != 0 { i = encodeVarintLaunchpad(dAtA, i, uint64(m.StartBlock)); i--; dAtA[i] = 0x48 }
	// f8: max_per_wallet (bytes)
	{ bz, err := m.MaxPerWallet.Marshal(); if err != nil { return 0, err }; if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintLaunchpad(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x42 } }
	// f7: hard_cap (bytes)
	{ bz, err := m.HardCap.Marshal(); if err != nil { return 0, err }; if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintLaunchpad(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x3a } }
	// f6: soft_cap (bytes)
	{ bz, err := m.SoftCap.Marshal(); if err != nil { return 0, err }; if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintLaunchpad(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x32 } }
	// f5: quote_denom (string)
	if len(m.QuoteDenom) > 0 { i -= len(m.QuoteDenom); copy(dAtA[i:], m.QuoteDenom); i = encodeVarintLaunchpad(dAtA, i, uint64(len(m.QuoteDenom))); i--; dAtA[i] = 0x2a }
	// f4: price_per_token (bytes)
	{ bz, err := m.PricePerToken.Marshal(); if err != nil { return 0, err }; if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintLaunchpad(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x22 } }
	// f3: token_supply (bytes)
	{ bz, err := m.TokenSupply.Marshal(); if err != nil { return 0, err }; if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintLaunchpad(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x1a } }
	// f2: token_denom (string)
	if len(m.TokenDenom) > 0 { i -= len(m.TokenDenom); copy(dAtA[i:], m.TokenDenom); i = encodeVarintLaunchpad(dAtA, i, uint64(len(m.TokenDenom))); i--; dAtA[i] = 0x12 }
	// f1: creator (string)
	if len(m.Creator) > 0 { i -= len(m.Creator); copy(dAtA[i:], m.Creator); i = encodeVarintLaunchpad(dAtA, i, uint64(len(m.Creator))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgCreateLaunch) Size() (n int) {
	if m == nil { return 0 }; var l int
	l = len(m.Creator); if l > 0 { n += 1 + l + sovLaunchpad(uint64(l)) }
	l = len(m.TokenDenom); if l > 0 { n += 1 + l + sovLaunchpad(uint64(l)) }
	bz3, _ := m.TokenSupply.Marshal(); l = len(bz3); if l > 0 { n += 1 + l + sovLaunchpad(uint64(l)) }
	bz4, _ := m.PricePerToken.Marshal(); l = len(bz4); if l > 0 { n += 1 + l + sovLaunchpad(uint64(l)) }
	l = len(m.QuoteDenom); if l > 0 { n += 1 + l + sovLaunchpad(uint64(l)) }
	bz6, _ := m.SoftCap.Marshal(); l = len(bz6); if l > 0 { n += 1 + l + sovLaunchpad(uint64(l)) }
	bz7, _ := m.HardCap.Marshal(); l = len(bz7); if l > 0 { n += 1 + l + sovLaunchpad(uint64(l)) }
	bz8, _ := m.MaxPerWallet.Marshal(); l = len(bz8); if l > 0 { n += 1 + l + sovLaunchpad(uint64(l)) }
	if m.StartBlock != 0 { n += 1 + sovLaunchpad(uint64(m.StartBlock)) }
	if m.EndBlock != 0 { n += 1 + sovLaunchpad(uint64(m.EndBlock)) }
	if m.VestingBlocks != 0 { n += 1 + sovLaunchpad(uint64(m.VestingBlocks)) }
	if m.TGEPercent != 0 { n += 1 + sovLaunchpad(uint64(m.TGEPercent)) }
	return n
}
func (m *MsgCreateLaunch) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7)
		if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthLaunchpad }; postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthLaunchpad }; if postIndex > l { return io.ErrUnexpectedEOF }; m.Creator = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthLaunchpad }; postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthLaunchpad }; if postIndex > l { return io.ErrUnexpectedEOF }; m.TokenDenom = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 3: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLengthLaunchpad }; postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthLaunchpad }; if postIndex > l { return io.ErrUnexpectedEOF }; if err := m.TokenSupply.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex
		case 4: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLengthLaunchpad }; postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthLaunchpad }; if postIndex > l { return io.ErrUnexpectedEOF }; if err := m.PricePerToken.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex
		case 5: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthLaunchpad }; postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthLaunchpad }; if postIndex > l { return io.ErrUnexpectedEOF }; m.QuoteDenom = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 6: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLengthLaunchpad }; postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthLaunchpad }; if postIndex > l { return io.ErrUnexpectedEOF }; if err := m.SoftCap.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex
		case 7: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLengthLaunchpad }; postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthLaunchpad }; if postIndex > l { return io.ErrUnexpectedEOF }; if err := m.HardCap.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex
		case 8: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLengthLaunchpad }; postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthLaunchpad }; if postIndex > l { return io.ErrUnexpectedEOF }; if err := m.MaxPerWallet.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex
		case 9: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.StartBlock = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.StartBlock |= int64(b&0x7F) << shift; if b < 0x80 { break } }
		case 10: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.EndBlock = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.EndBlock |= int64(b&0x7F) << shift; if b < 0x80 { break } }
		case 11: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.VestingBlocks = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.VestingBlocks |= int64(b&0x7F) << shift; if b < 0x80 { break } }
		case 12: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.TGEPercent = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.TGEPercent |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		default: iNdEx = preIndex; skippy, err := skipLaunchpad(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthLaunchpad }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgCreateLaunchResponse — launch_id(varint,1)
var xxx_messageInfo_MsgCreateLaunchResponse proto.InternalMessageInfo
func (m *MsgCreateLaunchResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateLaunchResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgCreateLaunchResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgCreateLaunchResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCreateLaunchResponse.Merge(m, src) }
func (m *MsgCreateLaunchResponse) XXX_Size() int { return m.Size() }
func (m *MsgCreateLaunchResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgCreateLaunchResponse.DiscardUnknown(m) }
func (m *MsgCreateLaunchResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgCreateLaunchResponse) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgCreateLaunchResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { i := len(dAtA); if m.LaunchID != 0 { i = encodeVarintLaunchpad(dAtA, i, uint64(m.LaunchID)); i--; dAtA[i] = 0x08 }; return len(dAtA) - i, nil }
func (m *MsgCreateLaunchResponse) Size() (n int) { if m == nil { return 0 }; if m.LaunchID != 0 { n += 1 + sovLaunchpad(uint64(m.LaunchID)) }; return n }
func (m *MsgCreateLaunchResponse) Unmarshal(dAtA []byte) error { l := len(dAtA); iNdEx := 0; for iNdEx < l { preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }; switch fieldNum { case 1: if int(wire&0x7) != 0 { return fmt.Errorf("proto: wrong wireType") }; m.LaunchID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.LaunchID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; default: iNdEx = preIndex; skippy, err := skipLaunchpad(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthLaunchpad }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy } }; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil }

// MsgContribute — sender(s,1) launch_id(varint,2) amount(bytes,3)
var xxx_messageInfo_MsgContribute proto.InternalMessageInfo
func (m *MsgContribute) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgContribute) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgContribute.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgContribute) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgContribute.Merge(m, src) }
func (m *MsgContribute) XXX_Size() int { return m.Size() }
func (m *MsgContribute) XXX_DiscardUnknown() { xxx_messageInfo_MsgContribute.DiscardUnknown(m) }
func (m *MsgContribute) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgContribute) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgContribute) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	{ bz, err := m.Amount.Marshal(); if err != nil { return 0, err }; if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintLaunchpad(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x1a } }
	if m.LaunchID != 0 { i = encodeVarintLaunchpad(dAtA, i, uint64(m.LaunchID)); i--; dAtA[i] = 0x10 }
	if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarintLaunchpad(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgContribute) Size() (n int) { if m == nil { return 0 }; l := len(m.Sender); if l > 0 { n += 1 + l + sovLaunchpad(uint64(l)) }; if m.LaunchID != 0 { n += 1 + sovLaunchpad(uint64(m.LaunchID)) }; bz, _ := m.Amount.Marshal(); l = len(bz); if l > 0 { n += 1 + l + sovLaunchpad(uint64(l)) }; return n }
func (m *MsgContribute) Unmarshal(dAtA []byte) error { l := len(dAtA); iNdEx := 0; for iNdEx < l { preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; fieldNum := int32(wire >> 3); wireType := int(wire & 0x7); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }; switch fieldNum { case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthLaunchpad }; postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthLaunchpad }; if postIndex > l { return io.ErrUnexpectedEOF }; m.Sender = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex; case 2: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.LaunchID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.LaunchID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; case 3: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLengthLaunchpad }; postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthLaunchpad }; if postIndex > l { return io.ErrUnexpectedEOF }; if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex; default: iNdEx = preIndex; skippy, err := skipLaunchpad(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthLaunchpad }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy } }; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil }

// MsgContributeResponse (empty)
var xxx_messageInfo_MsgContributeResponse proto.InternalMessageInfo
func (m *MsgContributeResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgContributeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgContributeResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgContributeResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgContributeResponse.Merge(m, src) }
func (m *MsgContributeResponse) XXX_Size() int { return m.Size() }
func (m *MsgContributeResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgContributeResponse.DiscardUnknown(m) }
func (m *MsgContributeResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgContributeResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgContributeResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgContributeResponse) Size() (n int) { return 0 }
func (m *MsgContributeResponse) Unmarshal(dAtA []byte) error { l := len(dAtA); iNdEx := 0; for iNdEx < l { preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }; iNdEx = preIndex; skippy, err := skipLaunchpad(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthLaunchpad }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy }; return nil }

// MsgClaimTokens — sender(s,1) launch_id(varint,2)
var xxx_messageInfo_MsgClaimTokens proto.InternalMessageInfo
func (m *MsgClaimTokens) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgClaimTokens) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgClaimTokens.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgClaimTokens) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgClaimTokens.Merge(m, src) }
func (m *MsgClaimTokens) XXX_Size() int { return m.Size() }
func (m *MsgClaimTokens) XXX_DiscardUnknown() { xxx_messageInfo_MsgClaimTokens.DiscardUnknown(m) }
func (m *MsgClaimTokens) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgClaimTokens) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgClaimTokens) MarshalToSizedBuffer(dAtA []byte) (int, error) { i := len(dAtA); if m.LaunchID != 0 { i = encodeVarintLaunchpad(dAtA, i, uint64(m.LaunchID)); i--; dAtA[i] = 0x10 }; if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarintLaunchpad(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }; return len(dAtA) - i, nil }
func (m *MsgClaimTokens) Size() (n int) { if m == nil { return 0 }; l := len(m.Sender); if l > 0 { n += 1 + l + sovLaunchpad(uint64(l)) }; if m.LaunchID != 0 { n += 1 + sovLaunchpad(uint64(m.LaunchID)) }; return n }
func (m *MsgClaimTokens) Unmarshal(dAtA []byte) error { l := len(dAtA); iNdEx := 0; for iNdEx < l { preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; fieldNum := int32(wire >> 3); wireType := int(wire & 0x7); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }; switch fieldNum { case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthLaunchpad }; postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthLaunchpad }; if postIndex > l { return io.ErrUnexpectedEOF }; m.Sender = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex; case 2: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.LaunchID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.LaunchID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; default: iNdEx = preIndex; skippy, err := skipLaunchpad(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthLaunchpad }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy } }; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil }

// MsgClaimTokensResponse — amount(bytes,1)
var xxx_messageInfo_MsgClaimTokensResponse proto.InternalMessageInfo
func (m *MsgClaimTokensResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgClaimTokensResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgClaimTokensResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgClaimTokensResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgClaimTokensResponse.Merge(m, src) }
func (m *MsgClaimTokensResponse) XXX_Size() int { return m.Size() }
func (m *MsgClaimTokensResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgClaimTokensResponse.DiscardUnknown(m) }
func (m *MsgClaimTokensResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgClaimTokensResponse) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgClaimTokensResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { i := len(dAtA); bz, err := m.Amount.Marshal(); if err != nil { return 0, err }; if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintLaunchpad(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x0a }; return len(dAtA) - i, nil }
func (m *MsgClaimTokensResponse) Size() (n int) { if m == nil { return 0 }; bz, _ := m.Amount.Marshal(); l := len(bz); if l > 0 { n += 1 + l + sovLaunchpad(uint64(l)) }; return n }
func (m *MsgClaimTokensResponse) Unmarshal(dAtA []byte) error { l := len(dAtA); iNdEx := 0; for iNdEx < l { preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }; switch fieldNum { case 1: if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLengthLaunchpad }; postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthLaunchpad }; if postIndex > l { return io.ErrUnexpectedEOF }; if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex; default: iNdEx = preIndex; skippy, err := skipLaunchpad(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthLaunchpad }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy } }; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil }

// MsgClaimRefund — same as MsgClaimTokens
var xxx_messageInfo_MsgClaimRefund proto.InternalMessageInfo
func (m *MsgClaimRefund) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgClaimRefund) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgClaimRefund.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgClaimRefund) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgClaimRefund.Merge(m, src) }
func (m *MsgClaimRefund) XXX_Size() int { return m.Size() }
func (m *MsgClaimRefund) XXX_DiscardUnknown() { xxx_messageInfo_MsgClaimRefund.DiscardUnknown(m) }
func (m *MsgClaimRefund) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgClaimRefund) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgClaimRefund) MarshalToSizedBuffer(dAtA []byte) (int, error) { i := len(dAtA); if m.LaunchID != 0 { i = encodeVarintLaunchpad(dAtA, i, uint64(m.LaunchID)); i--; dAtA[i] = 0x10 }; if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarintLaunchpad(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }; return len(dAtA) - i, nil }
func (m *MsgClaimRefund) Size() (n int) { if m == nil { return 0 }; l := len(m.Sender); if l > 0 { n += 1 + l + sovLaunchpad(uint64(l)) }; if m.LaunchID != 0 { n += 1 + sovLaunchpad(uint64(m.LaunchID)) }; return n }
func (m *MsgClaimRefund) Unmarshal(dAtA []byte) error { l := len(dAtA); iNdEx := 0; for iNdEx < l { preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; fieldNum := int32(wire >> 3); wireType := int(wire & 0x7); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }; switch fieldNum { case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthLaunchpad }; postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthLaunchpad }; if postIndex > l { return io.ErrUnexpectedEOF }; m.Sender = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex; case 2: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.LaunchID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.LaunchID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; default: iNdEx = preIndex; skippy, err := skipLaunchpad(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthLaunchpad }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy } }; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil }

// MsgClaimRefundResponse — amount(bytes,1), same as MsgClaimTokensResponse
var xxx_messageInfo_MsgClaimRefundResponse proto.InternalMessageInfo
func (m *MsgClaimRefundResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgClaimRefundResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgClaimRefundResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgClaimRefundResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgClaimRefundResponse.Merge(m, src) }
func (m *MsgClaimRefundResponse) XXX_Size() int { return m.Size() }
func (m *MsgClaimRefundResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgClaimRefundResponse.DiscardUnknown(m) }
func (m *MsgClaimRefundResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgClaimRefundResponse) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgClaimRefundResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { i := len(dAtA); bz, err := m.Amount.Marshal(); if err != nil { return 0, err }; if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintLaunchpad(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x0a }; return len(dAtA) - i, nil }
func (m *MsgClaimRefundResponse) Size() (n int) { if m == nil { return 0 }; bz, _ := m.Amount.Marshal(); l := len(bz); if l > 0 { n += 1 + l + sovLaunchpad(uint64(l)) }; return n }
func (m *MsgClaimRefundResponse) Unmarshal(dAtA []byte) error { l := len(dAtA); iNdEx := 0; for iNdEx < l { preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }; switch fieldNum { case 1: if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLengthLaunchpad }; postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthLaunchpad }; if postIndex > l { return io.ErrUnexpectedEOF }; if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex; default: iNdEx = preIndex; skippy, err := skipLaunchpad(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthLaunchpad }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy } }; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil }

// MsgFinalizeLaunch — authority(s,1) launch_id(varint,2)
var xxx_messageInfo_MsgFinalizeLaunch proto.InternalMessageInfo
func (m *MsgFinalizeLaunch) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgFinalizeLaunch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgFinalizeLaunch.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgFinalizeLaunch) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgFinalizeLaunch.Merge(m, src) }
func (m *MsgFinalizeLaunch) XXX_Size() int { return m.Size() }
func (m *MsgFinalizeLaunch) XXX_DiscardUnknown() { xxx_messageInfo_MsgFinalizeLaunch.DiscardUnknown(m) }
func (m *MsgFinalizeLaunch) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgFinalizeLaunch) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgFinalizeLaunch) MarshalToSizedBuffer(dAtA []byte) (int, error) { i := len(dAtA); if m.LaunchID != 0 { i = encodeVarintLaunchpad(dAtA, i, uint64(m.LaunchID)); i--; dAtA[i] = 0x10 }; if len(m.Authority) > 0 { i -= len(m.Authority); copy(dAtA[i:], m.Authority); i = encodeVarintLaunchpad(dAtA, i, uint64(len(m.Authority))); i--; dAtA[i] = 0x0a }; return len(dAtA) - i, nil }
func (m *MsgFinalizeLaunch) Size() (n int) { if m == nil { return 0 }; l := len(m.Authority); if l > 0 { n += 1 + l + sovLaunchpad(uint64(l)) }; if m.LaunchID != 0 { n += 1 + sovLaunchpad(uint64(m.LaunchID)) }; return n }
func (m *MsgFinalizeLaunch) Unmarshal(dAtA []byte) error { l := len(dAtA); iNdEx := 0; for iNdEx < l { preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; fieldNum := int32(wire >> 3); wireType := int(wire & 0x7); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }; switch fieldNum { case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthLaunchpad }; postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthLaunchpad }; if postIndex > l { return io.ErrUnexpectedEOF }; m.Authority = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex; case 2: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.LaunchID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.LaunchID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; default: iNdEx = preIndex; skippy, err := skipLaunchpad(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthLaunchpad }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy } }; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil }

// MsgFinalizeLaunchResponse — status(s,1)
var xxx_messageInfo_MsgFinalizeLaunchResponse proto.InternalMessageInfo
func (m *MsgFinalizeLaunchResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgFinalizeLaunchResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgFinalizeLaunchResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgFinalizeLaunchResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgFinalizeLaunchResponse.Merge(m, src) }
func (m *MsgFinalizeLaunchResponse) XXX_Size() int { return m.Size() }
func (m *MsgFinalizeLaunchResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgFinalizeLaunchResponse.DiscardUnknown(m) }
func (m *MsgFinalizeLaunchResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgFinalizeLaunchResponse) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgFinalizeLaunchResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { i := len(dAtA); if len(m.Status) > 0 { i -= len(m.Status); copy(dAtA[i:], m.Status); i = encodeVarintLaunchpad(dAtA, i, uint64(len(m.Status))); i--; dAtA[i] = 0x0a }; return len(dAtA) - i, nil }
func (m *MsgFinalizeLaunchResponse) Size() (n int) { if m == nil { return 0 }; l := len(m.Status); if l > 0 { n += 1 + l + sovLaunchpad(uint64(l)) }; return n }
func (m *MsgFinalizeLaunchResponse) Unmarshal(dAtA []byte) error { l := len(dAtA); iNdEx := 0; for iNdEx < l { preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }; switch fieldNum { case 1: if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowLaunchpad }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthLaunchpad }; postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthLaunchpad }; if postIndex > l { return io.ErrUnexpectedEOF }; m.Status = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex; default: iNdEx = preIndex; skippy, err := skipLaunchpad(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthLaunchpad }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy } }; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil }

// Helpers
var (
	ErrIntOverflowLaunchpad  = fmt.Errorf("proto: integer overflow")
	ErrInvalidLengthLaunchpad = fmt.Errorf("proto: negative length found during unmarshaling")
)
func encodeVarintLaunchpad(dAtA []byte, offset int, v uint64) int { offset -= sovLaunchpad(v); base := offset; for v >= 1<<7 { dAtA[offset] = uint8(v&0x7f | 0x80); v >>= 7; offset++ }; dAtA[offset] = uint8(v); return base }
func sovLaunchpad(x uint64) int { return (bits.Len64(x|1) + 6) / 7 }
func skipLaunchpad(dAtA []byte) (int, error) { l := len(dAtA); iNdEx := 0; depth := 0; for iNdEx < l { var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return 0, ErrIntOverflowLaunchpad }; if iNdEx >= l { return 0, io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= (uint64(b) & 0x7F) << shift; if b < 0x80 { break } }; wireType := int(wire & 0x7); switch wireType { case 0: for shift := uint(0); ; shift += 7 { if shift >= 64 { return 0, ErrIntOverflowLaunchpad }; if iNdEx >= l { return 0, io.ErrUnexpectedEOF }; iNdEx++; if dAtA[iNdEx-1] < 0x80 { break } }; case 1: iNdEx += 8; case 2: var length int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return 0, ErrIntOverflowLaunchpad }; if iNdEx >= l { return 0, io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; length |= (int(b) & 0x7F) << shift; if b < 0x80 { break } }; if length < 0 { return 0, ErrInvalidLengthLaunchpad }; iNdEx += length; case 3: depth++; case 4: if depth == 0 { return 0, fmt.Errorf("proto: wiretype end group for non-group") }; depth--; case 5: iNdEx += 4; default: return 0, fmt.Errorf("proto: illegal wireType %d", wireType) }; if iNdEx < 0 { return 0, ErrInvalidLengthLaunchpad }; if depth == 0 { return iNdEx, nil } }; return 0, io.ErrUnexpectedEOF }
