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

var fileDescriptorVaultTx []byte

// Descriptor methods
func (*MsgCreateVault) Descriptor() ([]byte, []int)                    { return fileDescriptorVaultTx, []int{0} }
func (*MsgCreateVaultResponse) Descriptor() ([]byte, []int)            { return fileDescriptorVaultTx, []int{1} }
func (*MsgDepositVault) Descriptor() ([]byte, []int)                   { return fileDescriptorVaultTx, []int{2} }
func (*MsgDepositVaultResponse) Descriptor() ([]byte, []int)           { return fileDescriptorVaultTx, []int{3} }
func (*MsgWithdrawVault) Descriptor() ([]byte, []int)                  { return fileDescriptorVaultTx, []int{4} }
func (*MsgWithdrawVaultResponse) Descriptor() ([]byte, []int)          { return fileDescriptorVaultTx, []int{5} }
func (*MsgCompoundVault) Descriptor() ([]byte, []int)                  { return fileDescriptorVaultTx, []int{6} }
func (*MsgCompoundVaultResponse) Descriptor() ([]byte, []int)          { return fileDescriptorVaultTx, []int{7} }
func (*MsgUpdateVaultStrategy) Descriptor() ([]byte, []int)            { return fileDescriptorVaultTx, []int{8} }
func (*MsgUpdateVaultStrategyResponse) Descriptor() ([]byte, []int)    { return fileDescriptorVaultTx, []int{9} }

func init() {
	proto.RegisterType((*MsgCreateVault)(nil), "syreen.vault.MsgCreateVault")
	proto.RegisterType((*MsgCreateVaultResponse)(nil), "syreen.vault.MsgCreateVaultResponse")
	proto.RegisterType((*MsgDepositVault)(nil), "syreen.vault.MsgDepositVault")
	proto.RegisterType((*MsgDepositVaultResponse)(nil), "syreen.vault.MsgDepositVaultResponse")
	proto.RegisterType((*MsgWithdrawVault)(nil), "syreen.vault.MsgWithdrawVault")
	proto.RegisterType((*MsgWithdrawVaultResponse)(nil), "syreen.vault.MsgWithdrawVaultResponse")
	proto.RegisterType((*MsgCompoundVault)(nil), "syreen.vault.MsgCompoundVault")
	proto.RegisterType((*MsgCompoundVaultResponse)(nil), "syreen.vault.MsgCompoundVaultResponse")
	proto.RegisterType((*MsgUpdateVaultStrategy)(nil), "syreen.vault.MsgUpdateVaultStrategy")
	proto.RegisterType((*MsgUpdateVaultStrategyResponse)(nil), "syreen.vault.MsgUpdateVaultStrategyResponse")
}

// ---------------------------------------------------------------------------
// MsgCreateVault — creator(s,1) name(s,2) deposit_denom(s,3) strategy_type(s,4) target_pool_ids(bytes,5) performance_fee(bytes,6)
// target_pool_ids is []uint64, encoded as bytes via JSON for simplicity (matches descriptor)
// ---------------------------------------------------------------------------
var xxx_messageInfo_MsgCreateVault proto.InternalMessageInfo
func (m *MsgCreateVault) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateVault) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgCreateVault.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgCreateVault) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCreateVault.Merge(m, src) }
func (m *MsgCreateVault) XXX_Size() int { return m.Size() }
func (m *MsgCreateVault) XXX_DiscardUnknown() { xxx_messageInfo_MsgCreateVault.DiscardUnknown(m) }

func (m *MsgCreateVault) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgCreateVault) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgCreateVault) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	// f6: performance_fee (bytes — LegacyDec)
	{ bz, err := m.PerformanceFee.Marshal(); if err != nil { return 0, err }; if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintVault(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x32 } }
	// f5: target_pool_ids (bytes — JSON-encoded []uint64)
	if len(m.TargetPoolIDs) > 0 { bz, _ := json.Marshal(m.TargetPoolIDs); i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintVault(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x2a }
	// f4: strategy_type (string)
	if len(m.StrategyType) > 0 { i -= len(m.StrategyType); copy(dAtA[i:], m.StrategyType); i = encodeVarintVault(dAtA, i, uint64(len(m.StrategyType))); i--; dAtA[i] = 0x22 }
	// f3: deposit_denom (string)
	if len(m.DepositDenom) > 0 { i -= len(m.DepositDenom); copy(dAtA[i:], m.DepositDenom); i = encodeVarintVault(dAtA, i, uint64(len(m.DepositDenom))); i--; dAtA[i] = 0x1a }
	// f2: name (string)
	if len(m.Name) > 0 { i -= len(m.Name); copy(dAtA[i:], m.Name); i = encodeVarintVault(dAtA, i, uint64(len(m.Name))); i--; dAtA[i] = 0x12 }
	// f1: creator (string)
	if len(m.Creator) > 0 { i -= len(m.Creator); copy(dAtA[i:], m.Creator); i = encodeVarintVault(dAtA, i, uint64(len(m.Creator))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgCreateVault) Size() (n int) {
	if m == nil { return 0 }; var l int
	l = len(m.Creator); if l > 0 { n += 1 + l + sovVault(uint64(l)) }
	l = len(m.Name); if l > 0 { n += 1 + l + sovVault(uint64(l)) }
	l = len(m.DepositDenom); if l > 0 { n += 1 + l + sovVault(uint64(l)) }
	l = len(m.StrategyType); if l > 0 { n += 1 + l + sovVault(uint64(l)) }
	if len(m.TargetPoolIDs) > 0 { bz, _ := json.Marshal(m.TargetPoolIDs); l = len(bz); n += 1 + l + sovVault(uint64(l)) }
	bz, _ := m.PerformanceFee.Marshal(); l = len(bz); if l > 0 { n += 1 + l + sovVault(uint64(l)) }
	return n
}
func (m *MsgCreateVault) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7)
		if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthVault }; postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthVault }; if postIndex > l { return io.ErrUnexpectedEOF }; m.Creator = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthVault }; postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthVault }; if postIndex > l { return io.ErrUnexpectedEOF }; m.Name = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 3: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthVault }; postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthVault }; if postIndex > l { return io.ErrUnexpectedEOF }; m.DepositDenom = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 4: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthVault }; postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthVault }; if postIndex > l { return io.ErrUnexpectedEOF }; m.StrategyType = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 5: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLengthVault }; postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthVault }; if postIndex > l { return io.ErrUnexpectedEOF }; if err := json.Unmarshal(dAtA[iNdEx:postIndex], &m.TargetPoolIDs); err != nil { return err }; iNdEx = postIndex
		case 6: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLengthVault }; postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthVault }; if postIndex > l { return io.ErrUnexpectedEOF }; if err := m.PerformanceFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex
		default: iNdEx = preIndex; skippy, err := skipVault(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthVault }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgCreateVaultResponse — vault_id(varint,1)
var xxx_messageInfo_MsgCreateVaultResponse proto.InternalMessageInfo
func (m *MsgCreateVaultResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateVaultResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgCreateVaultResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgCreateVaultResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCreateVaultResponse.Merge(m, src) }
func (m *MsgCreateVaultResponse) XXX_Size() int { return m.Size() }
func (m *MsgCreateVaultResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgCreateVaultResponse.DiscardUnknown(m) }
func (m *MsgCreateVaultResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgCreateVaultResponse) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgCreateVaultResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { i := len(dAtA); if m.VaultID != 0 { i = encodeVarintVault(dAtA, i, uint64(m.VaultID)); i--; dAtA[i] = 0x08 }; return len(dAtA) - i, nil }
func (m *MsgCreateVaultResponse) Size() (n int) { if m == nil { return 0 }; if m.VaultID != 0 { n += 1 + sovVault(uint64(m.VaultID)) }; return n }
func (m *MsgCreateVaultResponse) Unmarshal(dAtA []byte) error { l := len(dAtA); iNdEx := 0; for iNdEx < l { preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }; switch fieldNum { case 1: if int(wire&0x7) != 0 { return fmt.Errorf("proto: wrong wireType") }; m.VaultID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.VaultID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; default: iNdEx = preIndex; skippy, err := skipVault(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthVault }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy } }; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil }

// MsgDepositVault — sender(s,1) vault_id(varint,2) amount(bytes,3)
var xxx_messageInfo_MsgDepositVault proto.InternalMessageInfo
func (m *MsgDepositVault) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgDepositVault) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgDepositVault.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgDepositVault) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgDepositVault.Merge(m, src) }
func (m *MsgDepositVault) XXX_Size() int { return m.Size() }
func (m *MsgDepositVault) XXX_DiscardUnknown() { xxx_messageInfo_MsgDepositVault.DiscardUnknown(m) }
func (m *MsgDepositVault) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgDepositVault) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgDepositVault) MarshalToSizedBuffer(dAtA []byte) (int, error) { i := len(dAtA); { bz, err := m.Amount.Marshal(); if err != nil { return 0, err }; if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintVault(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x1a } }; if m.VaultID != 0 { i = encodeVarintVault(dAtA, i, uint64(m.VaultID)); i--; dAtA[i] = 0x10 }; if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarintVault(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }; return len(dAtA) - i, nil }
func (m *MsgDepositVault) Size() (n int) { if m == nil { return 0 }; l := len(m.Sender); if l > 0 { n += 1 + l + sovVault(uint64(l)) }; if m.VaultID != 0 { n += 1 + sovVault(uint64(m.VaultID)) }; bz, _ := m.Amount.Marshal(); l = len(bz); if l > 0 { n += 1 + l + sovVault(uint64(l)) }; return n }
func (m *MsgDepositVault) Unmarshal(dAtA []byte) error { l := len(dAtA); iNdEx := 0; for iNdEx < l { preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; fieldNum := int32(wire >> 3); wireType := int(wire & 0x7); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }; switch fieldNum { case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthVault }; postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthVault }; if postIndex > l { return io.ErrUnexpectedEOF }; m.Sender = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex; case 2: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.VaultID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.VaultID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; case 3: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLengthVault }; postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthVault }; if postIndex > l { return io.ErrUnexpectedEOF }; if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex; default: iNdEx = preIndex; skippy, err := skipVault(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthVault }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy } }; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil }

// MsgDepositVaultResponse — shares_minted(bytes,1)
var xxx_messageInfo_MsgDepositVaultResponse proto.InternalMessageInfo
func (m *MsgDepositVaultResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgDepositVaultResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgDepositVaultResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgDepositVaultResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgDepositVaultResponse.Merge(m, src) }
func (m *MsgDepositVaultResponse) XXX_Size() int { return m.Size() }
func (m *MsgDepositVaultResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgDepositVaultResponse.DiscardUnknown(m) }
func (m *MsgDepositVaultResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgDepositVaultResponse) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgDepositVaultResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { i := len(dAtA); bz, err := m.SharesMinted.Marshal(); if err != nil { return 0, err }; if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintVault(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x0a }; return len(dAtA) - i, nil }
func (m *MsgDepositVaultResponse) Size() (n int) { if m == nil { return 0 }; bz, _ := m.SharesMinted.Marshal(); l := len(bz); if l > 0 { n += 1 + l + sovVault(uint64(l)) }; return n }
func (m *MsgDepositVaultResponse) Unmarshal(dAtA []byte) error { l := len(dAtA); iNdEx := 0; for iNdEx < l { preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }; switch fieldNum { case 1: if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLengthVault }; postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthVault }; if postIndex > l { return io.ErrUnexpectedEOF }; if err := m.SharesMinted.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex; default: iNdEx = preIndex; skippy, err := skipVault(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthVault }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy } }; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil }

// MsgWithdrawVault — sender(s,1) vault_id(varint,2) shares(bytes,3) — same shape as DepositVault
var xxx_messageInfo_MsgWithdrawVault proto.InternalMessageInfo
func (m *MsgWithdrawVault) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgWithdrawVault) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgWithdrawVault.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgWithdrawVault) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgWithdrawVault.Merge(m, src) }
func (m *MsgWithdrawVault) XXX_Size() int { return m.Size() }
func (m *MsgWithdrawVault) XXX_DiscardUnknown() { xxx_messageInfo_MsgWithdrawVault.DiscardUnknown(m) }
func (m *MsgWithdrawVault) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgWithdrawVault) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgWithdrawVault) MarshalToSizedBuffer(dAtA []byte) (int, error) { i := len(dAtA); { bz, err := m.Shares.Marshal(); if err != nil { return 0, err }; if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintVault(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x1a } }; if m.VaultID != 0 { i = encodeVarintVault(dAtA, i, uint64(m.VaultID)); i--; dAtA[i] = 0x10 }; if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarintVault(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }; return len(dAtA) - i, nil }
func (m *MsgWithdrawVault) Size() (n int) { if m == nil { return 0 }; l := len(m.Sender); if l > 0 { n += 1 + l + sovVault(uint64(l)) }; if m.VaultID != 0 { n += 1 + sovVault(uint64(m.VaultID)) }; bz, _ := m.Shares.Marshal(); l = len(bz); if l > 0 { n += 1 + l + sovVault(uint64(l)) }; return n }
func (m *MsgWithdrawVault) Unmarshal(dAtA []byte) error { l := len(dAtA); iNdEx := 0; for iNdEx < l { preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; fieldNum := int32(wire >> 3); wireType := int(wire & 0x7); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }; switch fieldNum { case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthVault }; postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthVault }; if postIndex > l { return io.ErrUnexpectedEOF }; m.Sender = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex; case 2: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.VaultID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.VaultID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; case 3: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLengthVault }; postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthVault }; if postIndex > l { return io.ErrUnexpectedEOF }; if err := m.Shares.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex; default: iNdEx = preIndex; skippy, err := skipVault(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthVault }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy } }; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil }

// MsgWithdrawVaultResponse — amount_returned(bytes,1)
var xxx_messageInfo_MsgWithdrawVaultResponse proto.InternalMessageInfo
func (m *MsgWithdrawVaultResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgWithdrawVaultResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgWithdrawVaultResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgWithdrawVaultResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgWithdrawVaultResponse.Merge(m, src) }
func (m *MsgWithdrawVaultResponse) XXX_Size() int { return m.Size() }
func (m *MsgWithdrawVaultResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgWithdrawVaultResponse.DiscardUnknown(m) }
func (m *MsgWithdrawVaultResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgWithdrawVaultResponse) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgWithdrawVaultResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { i := len(dAtA); bz, err := m.AmountReturned.Marshal(); if err != nil { return 0, err }; if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintVault(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x0a }; return len(dAtA) - i, nil }
func (m *MsgWithdrawVaultResponse) Size() (n int) { if m == nil { return 0 }; bz, _ := m.AmountReturned.Marshal(); l := len(bz); if l > 0 { n += 1 + l + sovVault(uint64(l)) }; return n }
func (m *MsgWithdrawVaultResponse) Unmarshal(dAtA []byte) error { l := len(dAtA); iNdEx := 0; for iNdEx < l { preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }; switch fieldNum { case 1: if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLengthVault }; postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthVault }; if postIndex > l { return io.ErrUnexpectedEOF }; if err := m.AmountReturned.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex; default: iNdEx = preIndex; skippy, err := skipVault(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthVault }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy } }; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil }

// MsgCompoundVault — sender(s,1) vault_id(varint,2)
var xxx_messageInfo_MsgCompoundVault proto.InternalMessageInfo
func (m *MsgCompoundVault) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCompoundVault) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgCompoundVault.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgCompoundVault) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCompoundVault.Merge(m, src) }
func (m *MsgCompoundVault) XXX_Size() int { return m.Size() }
func (m *MsgCompoundVault) XXX_DiscardUnknown() { xxx_messageInfo_MsgCompoundVault.DiscardUnknown(m) }
func (m *MsgCompoundVault) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgCompoundVault) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgCompoundVault) MarshalToSizedBuffer(dAtA []byte) (int, error) { i := len(dAtA); if m.VaultID != 0 { i = encodeVarintVault(dAtA, i, uint64(m.VaultID)); i--; dAtA[i] = 0x10 }; if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarintVault(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }; return len(dAtA) - i, nil }
func (m *MsgCompoundVault) Size() (n int) { if m == nil { return 0 }; l := len(m.Sender); if l > 0 { n += 1 + l + sovVault(uint64(l)) }; if m.VaultID != 0 { n += 1 + sovVault(uint64(m.VaultID)) }; return n }
func (m *MsgCompoundVault) Unmarshal(dAtA []byte) error { l := len(dAtA); iNdEx := 0; for iNdEx < l { preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; fieldNum := int32(wire >> 3); wireType := int(wire & 0x7); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }; switch fieldNum { case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthVault }; postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthVault }; if postIndex > l { return io.ErrUnexpectedEOF }; m.Sender = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex; case 2: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.VaultID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.VaultID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; default: iNdEx = preIndex; skippy, err := skipVault(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthVault }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy } }; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil }

// MsgCompoundVaultResponse — yield_generated(bytes,1)
var xxx_messageInfo_MsgCompoundVaultResponse proto.InternalMessageInfo
func (m *MsgCompoundVaultResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCompoundVaultResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgCompoundVaultResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgCompoundVaultResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCompoundVaultResponse.Merge(m, src) }
func (m *MsgCompoundVaultResponse) XXX_Size() int { return m.Size() }
func (m *MsgCompoundVaultResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgCompoundVaultResponse.DiscardUnknown(m) }
func (m *MsgCompoundVaultResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgCompoundVaultResponse) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgCompoundVaultResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { i := len(dAtA); bz, err := m.YieldGenerated.Marshal(); if err != nil { return 0, err }; if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintVault(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x0a }; return len(dAtA) - i, nil }
func (m *MsgCompoundVaultResponse) Size() (n int) { if m == nil { return 0 }; bz, _ := m.YieldGenerated.Marshal(); l := len(bz); if l > 0 { n += 1 + l + sovVault(uint64(l)) }; return n }
func (m *MsgCompoundVaultResponse) Unmarshal(dAtA []byte) error { l := len(dAtA); iNdEx := 0; for iNdEx < l { preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }; switch fieldNum { case 1: if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLengthVault }; postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthVault }; if postIndex > l { return io.ErrUnexpectedEOF }; if err := m.YieldGenerated.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex; default: iNdEx = preIndex; skippy, err := skipVault(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthVault }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy } }; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil }

// MsgUpdateVaultStrategy — creator(s,1) vault_id(varint,2) strategy_type(s,3) target_pool_ids(bytes,4)
var xxx_messageInfo_MsgUpdateVaultStrategy proto.InternalMessageInfo
func (m *MsgUpdateVaultStrategy) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgUpdateVaultStrategy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgUpdateVaultStrategy.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgUpdateVaultStrategy) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgUpdateVaultStrategy.Merge(m, src) }
func (m *MsgUpdateVaultStrategy) XXX_Size() int { return m.Size() }
func (m *MsgUpdateVaultStrategy) XXX_DiscardUnknown() { xxx_messageInfo_MsgUpdateVaultStrategy.DiscardUnknown(m) }
func (m *MsgUpdateVaultStrategy) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgUpdateVaultStrategy) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgUpdateVaultStrategy) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.TargetPoolIDs) > 0 { bz, _ := json.Marshal(m.TargetPoolIDs); i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintVault(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x22 }
	if len(m.StrategyType) > 0 { i -= len(m.StrategyType); copy(dAtA[i:], m.StrategyType); i = encodeVarintVault(dAtA, i, uint64(len(m.StrategyType))); i--; dAtA[i] = 0x1a }
	if m.VaultID != 0 { i = encodeVarintVault(dAtA, i, uint64(m.VaultID)); i--; dAtA[i] = 0x10 }
	if len(m.Creator) > 0 { i -= len(m.Creator); copy(dAtA[i:], m.Creator); i = encodeVarintVault(dAtA, i, uint64(len(m.Creator))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgUpdateVaultStrategy) Size() (n int) { if m == nil { return 0 }; var l int; l = len(m.Creator); if l > 0 { n += 1 + l + sovVault(uint64(l)) }; if m.VaultID != 0 { n += 1 + sovVault(uint64(m.VaultID)) }; l = len(m.StrategyType); if l > 0 { n += 1 + l + sovVault(uint64(l)) }; if len(m.TargetPoolIDs) > 0 { bz, _ := json.Marshal(m.TargetPoolIDs); l = len(bz); n += 1 + l + sovVault(uint64(l)) }; return n }
func (m *MsgUpdateVaultStrategy) Unmarshal(dAtA []byte) error { l := len(dAtA); iNdEx := 0; for iNdEx < l { preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; fieldNum := int32(wire >> 3); wireType := int(wire & 0x7); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }; switch fieldNum { case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthVault }; postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthVault }; if postIndex > l { return io.ErrUnexpectedEOF }; m.Creator = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex; case 2: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.VaultID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.VaultID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; case 3: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthVault }; postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthVault }; if postIndex > l { return io.ErrUnexpectedEOF }; m.StrategyType = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex; case 4: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLengthVault }; postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthVault }; if postIndex > l { return io.ErrUnexpectedEOF }; if err := json.Unmarshal(dAtA[iNdEx:postIndex], &m.TargetPoolIDs); err != nil { return err }; iNdEx = postIndex; default: iNdEx = preIndex; skippy, err := skipVault(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthVault }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy } }; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil }

// MsgUpdateVaultStrategyResponse (empty)
var xxx_messageInfo_MsgUpdateVaultStrategyResponse proto.InternalMessageInfo
func (m *MsgUpdateVaultStrategyResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgUpdateVaultStrategyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgUpdateVaultStrategyResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgUpdateVaultStrategyResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgUpdateVaultStrategyResponse.Merge(m, src) }
func (m *MsgUpdateVaultStrategyResponse) XXX_Size() int { return m.Size() }
func (m *MsgUpdateVaultStrategyResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgUpdateVaultStrategyResponse.DiscardUnknown(m) }
func (m *MsgUpdateVaultStrategyResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgUpdateVaultStrategyResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgUpdateVaultStrategyResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgUpdateVaultStrategyResponse) Size() (n int) { return 0 }
func (m *MsgUpdateVaultStrategyResponse) Unmarshal(dAtA []byte) error { l := len(dAtA); iNdEx := 0; for iNdEx < l { preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowVault }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }; iNdEx = preIndex; skippy, err := skipVault(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthVault }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy }; return nil }

// Helpers
var (
	ErrIntOverflowVault  = fmt.Errorf("proto: integer overflow")
	ErrInvalidLengthVault = fmt.Errorf("proto: negative length found during unmarshaling")
)
func encodeVarintVault(dAtA []byte, offset int, v uint64) int { offset -= sovVault(v); base := offset; for v >= 1<<7 { dAtA[offset] = uint8(v&0x7f | 0x80); v >>= 7; offset++ }; dAtA[offset] = uint8(v); return base }
func sovVault(x uint64) int { return (bits.Len64(x|1) + 6) / 7 }
func skipVault(dAtA []byte) (int, error) { l := len(dAtA); iNdEx := 0; depth := 0; for iNdEx < l { var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return 0, ErrIntOverflowVault }; if iNdEx >= l { return 0, io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= (uint64(b) & 0x7F) << shift; if b < 0x80 { break } }; wireType := int(wire & 0x7); switch wireType { case 0: for shift := uint(0); ; shift += 7 { if shift >= 64 { return 0, ErrIntOverflowVault }; if iNdEx >= l { return 0, io.ErrUnexpectedEOF }; iNdEx++; if dAtA[iNdEx-1] < 0x80 { break } }; case 1: iNdEx += 8; case 2: var length int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return 0, ErrIntOverflowVault }; if iNdEx >= l { return 0, io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; length |= (int(b) & 0x7F) << shift; if b < 0x80 { break } }; if length < 0 { return 0, ErrInvalidLengthVault }; iNdEx += length; case 3: depth++; case 4: if depth == 0 { return 0, fmt.Errorf("proto: wiretype end group for non-group") }; depth--; case 5: iNdEx += 4; default: return 0, fmt.Errorf("proto: illegal wireType %d", wireType) }; if iNdEx < 0 { return 0, ErrInvalidLengthVault }; if depth == 0 { return iNdEx, nil } }; return 0, io.ErrUnexpectedEOF }
