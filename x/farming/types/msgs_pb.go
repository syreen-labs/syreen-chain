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

// Descriptor methods
func (*MsgStake) Descriptor() ([]byte, []int)              { return fileDescriptorFarmingTx, []int{0} }
func (*MsgStakeResponse) Descriptor() ([]byte, []int)      { return fileDescriptorFarmingTx, []int{1} }
func (*MsgUnstake) Descriptor() ([]byte, []int)            { return fileDescriptorFarmingTx, []int{2} }
func (*MsgUnstakeResponse) Descriptor() ([]byte, []int)    { return fileDescriptorFarmingTx, []int{3} }
func (*MsgClaimReward) Descriptor() ([]byte, []int)        { return fileDescriptorFarmingTx, []int{4} }
func (*MsgClaimRewardResponse) Descriptor() ([]byte, []int) { return fileDescriptorFarmingTx, []int{5} }
func (*MsgCreateFarm) Descriptor() ([]byte, []int)         { return fileDescriptorFarmingTx, []int{6} }
func (*MsgCreateFarmResponse) Descriptor() ([]byte, []int) { return fileDescriptorFarmingTx, []int{7} }
func (*MsgUpdateFarm) Descriptor() ([]byte, []int)         { return fileDescriptorFarmingTx, []int{8} }
func (*MsgUpdateFarmResponse) Descriptor() ([]byte, []int) { return fileDescriptorFarmingTx, []int{9} }

func init() {
	proto.RegisterType((*MsgStake)(nil), "syreen.farming.MsgStake")
	proto.RegisterType((*MsgStakeResponse)(nil), "syreen.farming.MsgStakeResponse")
	proto.RegisterType((*MsgUnstake)(nil), "syreen.farming.MsgUnstake")
	proto.RegisterType((*MsgUnstakeResponse)(nil), "syreen.farming.MsgUnstakeResponse")
	proto.RegisterType((*MsgClaimReward)(nil), "syreen.farming.MsgClaimReward")
	proto.RegisterType((*MsgClaimRewardResponse)(nil), "syreen.farming.MsgClaimRewardResponse")
	proto.RegisterType((*MsgCreateFarm)(nil), "syreen.farming.MsgCreateFarm")
	proto.RegisterType((*MsgCreateFarmResponse)(nil), "syreen.farming.MsgCreateFarmResponse")
	proto.RegisterType((*MsgUpdateFarm)(nil), "syreen.farming.MsgUpdateFarm")
	proto.RegisterType((*MsgUpdateFarmResponse)(nil), "syreen.farming.MsgUpdateFarmResponse")
}

// Helper functions
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

// Proto methods for all message types
func (m *MsgStake) ProtoMessage()     {}
func (m *MsgStake) Reset()            { *m = MsgStake{} }
func (m *MsgStake) String() string    { return fmt.Sprintf("MsgStake{%s, %d, %s}", m.Sender, m.PoolID, m.Amount) }
func (m *MsgUnstake) ProtoMessage()   {}
func (m *MsgUnstake) Reset()          { *m = MsgUnstake{} }
func (m *MsgUnstake) String() string  { return fmt.Sprintf("MsgUnstake{%s, %d, %s}", m.Sender, m.PoolID, m.Amount) }
func (m *MsgClaimReward) ProtoMessage() {}
func (m *MsgClaimReward) Reset()        { *m = MsgClaimReward{} }
func (m *MsgClaimReward) String() string { return fmt.Sprintf("MsgClaimReward{%s, %d}", m.Sender, m.PoolID) }
func (m *MsgCreateFarm) ProtoMessage()  {}
func (m *MsgCreateFarm) Reset()         { *m = MsgCreateFarm{} }
func (m *MsgCreateFarm) String() string { return fmt.Sprintf("MsgCreateFarm{%d}", m.PoolID) }
func (m *MsgUpdateFarm) ProtoMessage()  {}
func (m *MsgUpdateFarm) Reset()         { *m = MsgUpdateFarm{} }
func (m *MsgUpdateFarm) String() string { return fmt.Sprintf("MsgUpdateFarm{%d}", m.PoolID) }
func (m *MsgStakeResponse) ProtoMessage()       {}
func (m *MsgStakeResponse) Reset()              { *m = MsgStakeResponse{} }
func (m *MsgStakeResponse) String() string      { return "MsgStakeResponse" }
func (m *MsgUnstakeResponse) ProtoMessage()     {}
func (m *MsgUnstakeResponse) Reset()            { *m = MsgUnstakeResponse{} }
func (m *MsgUnstakeResponse) String() string    { return "MsgUnstakeResponse" }
func (m *MsgClaimRewardResponse) ProtoMessage()  {}
func (m *MsgClaimRewardResponse) Reset()         { *m = MsgClaimRewardResponse{} }
func (m *MsgClaimRewardResponse) String() string { return "MsgClaimRewardResponse" }
func (m *MsgCreateFarmResponse) ProtoMessage()   {}
func (m *MsgCreateFarmResponse) Reset()          { *m = MsgCreateFarmResponse{} }
func (m *MsgCreateFarmResponse) String() string  { return "MsgCreateFarmResponse" }
func (m *MsgUpdateFarmResponse) ProtoMessage()   {}
func (m *MsgUpdateFarmResponse) Reset()          { *m = MsgUpdateFarmResponse{} }
func (m *MsgUpdateFarmResponse) String() string  { return "MsgUpdateFarmResponse" }

// XXX_MessageName
func (m *MsgStake) XXX_MessageName() string              { return "syreen.farming.MsgStake" }
func (m *MsgStakeResponse) XXX_MessageName() string      { return "syreen.farming.MsgStakeResponse" }
func (m *MsgUnstake) XXX_MessageName() string            { return "syreen.farming.MsgUnstake" }
func (m *MsgUnstakeResponse) XXX_MessageName() string    { return "syreen.farming.MsgUnstakeResponse" }
func (m *MsgClaimReward) XXX_MessageName() string        { return "syreen.farming.MsgClaimReward" }
func (m *MsgClaimRewardResponse) XXX_MessageName() string { return "syreen.farming.MsgClaimRewardResponse" }
func (m *MsgCreateFarm) XXX_MessageName() string         { return "syreen.farming.MsgCreateFarm" }
func (m *MsgCreateFarmResponse) XXX_MessageName() string { return "syreen.farming.MsgCreateFarmResponse" }
func (m *MsgUpdateFarm) XXX_MessageName() string         { return "syreen.farming.MsgUpdateFarm" }
func (m *MsgUpdateFarmResponse) XXX_MessageName() string { return "syreen.farming.MsgUpdateFarmResponse" }

// ---------------------------------------------------------------------------
// MsgStake (sender:string/f1, pool_id:uint64/f2, amount:bytes/f3)
// ---------------------------------------------------------------------------
var xxx_messageInfo_MsgStake proto.InternalMessageInfo
func (m *MsgStake) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgStake) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgStake.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgStake) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgStake.Merge(m, src) }
func (m *MsgStake) XXX_Size() int { return m.Size() }
func (m *MsgStake) XXX_DiscardUnknown() { xxx_messageInfo_MsgStake.DiscardUnknown(m) }
func (m *MsgStake) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgStake) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgStake) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.Amount.IsNil() { bz := marshalBytes(m.Amount); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x1a } }
	if m.PoolID != 0 { i = encodeVarint(dAtA, i, uint64(m.PoolID)); i--; dAtA[i] = 0x10 }
	if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarint(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgStake) Size() (n int) {
	if m == nil { return 0 }; var l int
	l = len(m.Sender); if l > 0 { n += 1 + l + sov(uint64(l)) }
	if m.PoolID != 0 { n += 1 + sov(uint64(m.PoolID)) }
	if !m.Amount.IsNil() { bz := marshalBytes(m.Amount); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	return n
}
func (m *MsgStake) Unmarshal(dAtA []byte) error {
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
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType for Amount") }
			var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }
			if byteLen < 0 { return ErrInvalidLength }; postIndex := iNdEx + byteLen; if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }
			_ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.Amount); iNdEx = postIndex
		default:
			iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgStakeResponse (pending_reward:bytes/f1)
var xxx_messageInfo_MsgStakeResponse proto.InternalMessageInfo
func (m *MsgStakeResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgStakeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgStakeResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgStakeResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgStakeResponse.Merge(m, src) }
func (m *MsgStakeResponse) XXX_Size() int { return m.Size() }
func (m *MsgStakeResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgStakeResponse.DiscardUnknown(m) }
func (m *MsgStakeResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgStakeResponse) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgStakeResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.PendingReward.IsNil() { bz := marshalBytes(m.PendingReward); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x0a } }
	return len(dAtA) - i, nil
}
func (m *MsgStakeResponse) Size() (n int) { if m == nil { return 0 }; var l int; if !m.PendingReward.IsNil() { bz := marshalBytes(m.PendingReward); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }; return n }
func (m *MsgStakeResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLength }; postIndex := iNdEx + byteLen; if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.PendingReward); iNdEx = postIndex
		default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgUnstake = same shape as MsgStake
var xxx_messageInfo_MsgUnstake proto.InternalMessageInfo
func (m *MsgUnstake) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgUnstake) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgUnstake.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgUnstake) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgUnstake.Merge(m, src) }
func (m *MsgUnstake) XXX_Size() int { return m.Size() }
func (m *MsgUnstake) XXX_DiscardUnknown() { xxx_messageInfo_MsgUnstake.DiscardUnknown(m) }
func (m *MsgUnstake) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgUnstake) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgUnstake) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.Amount.IsNil() { bz := marshalBytes(m.Amount); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x1a } }
	if m.PoolID != 0 { i = encodeVarint(dAtA, i, uint64(m.PoolID)); i--; dAtA[i] = 0x10 }
	if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarint(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgUnstake) Size() (n int) {
	if m == nil { return 0 }; var l int; l = len(m.Sender); if l > 0 { n += 1 + l + sov(uint64(l)) }; if m.PoolID != 0 { n += 1 + sov(uint64(m.PoolID)) }; if !m.Amount.IsNil() { bz := marshalBytes(m.Amount); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }; return n
}
func (m *MsgUnstake) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; postIndex := iNdEx + int(stringLen); if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; m.Sender = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.PoolID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.PoolID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		case 3: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLength }; postIndex := iNdEx + byteLen; if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.Amount); iNdEx = postIndex
		default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgUnstakeResponse
var xxx_messageInfo_MsgUnstakeResponse proto.InternalMessageInfo
func (m *MsgUnstakeResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgUnstakeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgUnstakeResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgUnstakeResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgUnstakeResponse.Merge(m, src) }
func (m *MsgUnstakeResponse) XXX_Size() int { return m.Size() }
func (m *MsgUnstakeResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgUnstakeResponse.DiscardUnknown(m) }
func (m *MsgUnstakeResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgUnstakeResponse) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgUnstakeResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.ClaimedReward.IsNil() { bz := marshalBytes(m.ClaimedReward); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x0a } }
	return len(dAtA) - i, nil
}
func (m *MsgUnstakeResponse) Size() (n int) { if m == nil { return 0 }; var l int; if !m.ClaimedReward.IsNil() { bz := marshalBytes(m.ClaimedReward); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }; return n }
func (m *MsgUnstakeResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLength }; postIndex := iNdEx + byteLen; if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.ClaimedReward); iNdEx = postIndex
		default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgClaimReward (sender:string/f1, pool_id:uint64/f2)
var xxx_messageInfo_MsgClaimReward proto.InternalMessageInfo
func (m *MsgClaimReward) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgClaimReward) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgClaimReward.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgClaimReward) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgClaimReward.Merge(m, src) }
func (m *MsgClaimReward) XXX_Size() int { return m.Size() }
func (m *MsgClaimReward) XXX_DiscardUnknown() { xxx_messageInfo_MsgClaimReward.DiscardUnknown(m) }
func (m *MsgClaimReward) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgClaimReward) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgClaimReward) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.PoolID != 0 { i = encodeVarint(dAtA, i, uint64(m.PoolID)); i--; dAtA[i] = 0x10 }
	if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarint(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgClaimReward) Size() (n int) { if m == nil { return 0 }; var l int; l = len(m.Sender); if l > 0 { n += 1 + l + sov(uint64(l)) }; if m.PoolID != 0 { n += 1 + sov(uint64(m.PoolID)) }; return n }
func (m *MsgClaimReward) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; postIndex := iNdEx + int(stringLen); if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; m.Sender = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.PoolID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.PoolID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgClaimRewardResponse (amount:bytes/f1) — same pattern as MsgStakeResponse
var xxx_messageInfo_MsgClaimRewardResponse proto.InternalMessageInfo
func (m *MsgClaimRewardResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgClaimRewardResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgClaimRewardResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgClaimRewardResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgClaimRewardResponse.Merge(m, src) }
func (m *MsgClaimRewardResponse) XXX_Size() int { return m.Size() }
func (m *MsgClaimRewardResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgClaimRewardResponse.DiscardUnknown(m) }
func (m *MsgClaimRewardResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgClaimRewardResponse) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgClaimRewardResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if !m.Amount.IsNil() { bz := marshalBytes(m.Amount); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x0a } }
	return len(dAtA) - i, nil
}
func (m *MsgClaimRewardResponse) Size() (n int) { if m == nil { return 0 }; var l int; if !m.Amount.IsNil() { bz := marshalBytes(m.Amount); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }; return n }
func (m *MsgClaimRewardResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLength }; postIndex := iNdEx + byteLen; if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.Amount); iNdEx = postIndex
		default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgCreateFarm (authority:string/f1, pool_id:uint64/f2, lp_denom:string/f3, reward_per_block:bytes/f4, start_block:int64/f5, end_block:int64/f6)
var xxx_messageInfo_MsgCreateFarm proto.InternalMessageInfo
func (m *MsgCreateFarm) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateFarm) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgCreateFarm.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgCreateFarm) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCreateFarm.Merge(m, src) }
func (m *MsgCreateFarm) XXX_Size() int { return m.Size() }
func (m *MsgCreateFarm) XXX_DiscardUnknown() { xxx_messageInfo_MsgCreateFarm.DiscardUnknown(m) }
func (m *MsgCreateFarm) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgCreateFarm) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgCreateFarm) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.EndBlock != 0 { i = encodeVarint(dAtA, i, uint64(m.EndBlock)); i--; dAtA[i] = 0x30 }
	if m.StartBlock != 0 { i = encodeVarint(dAtA, i, uint64(m.StartBlock)); i--; dAtA[i] = 0x28 }
	if !m.RewardPerBlock.IsNil() { bz := marshalBytes(m.RewardPerBlock); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x22 } }
	if len(m.LPDenom) > 0 { i -= len(m.LPDenom); copy(dAtA[i:], m.LPDenom); i = encodeVarint(dAtA, i, uint64(len(m.LPDenom))); i--; dAtA[i] = 0x1a }
	if m.PoolID != 0 { i = encodeVarint(dAtA, i, uint64(m.PoolID)); i--; dAtA[i] = 0x10 }
	if len(m.Authority) > 0 { i -= len(m.Authority); copy(dAtA[i:], m.Authority); i = encodeVarint(dAtA, i, uint64(len(m.Authority))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgCreateFarm) Size() (n int) {
	if m == nil { return 0 }; var l int
	l = len(m.Authority); if l > 0 { n += 1 + l + sov(uint64(l)) }
	if m.PoolID != 0 { n += 1 + sov(uint64(m.PoolID)) }
	l = len(m.LPDenom); if l > 0 { n += 1 + l + sov(uint64(l)) }
	if !m.RewardPerBlock.IsNil() { bz := marshalBytes(m.RewardPerBlock); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	if m.StartBlock != 0 { n += 1 + sov(uint64(m.StartBlock)) }
	if m.EndBlock != 0 { n += 1 + sov(uint64(m.EndBlock)) }
	return n
}
func (m *MsgCreateFarm) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; postIndex := iNdEx + int(stringLen); if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; m.Authority = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.PoolID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.PoolID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		case 3: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; postIndex := iNdEx + int(stringLen); if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; m.LPDenom = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 4: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLength }; postIndex := iNdEx + byteLen; if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.RewardPerBlock); iNdEx = postIndex
		case 5: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.StartBlock = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.StartBlock |= int64(b&0x7F) << shift; if b < 0x80 { break } }
		case 6: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.EndBlock = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.EndBlock |= int64(b&0x7F) << shift; if b < 0x80 { break } }
		default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgCreateFarmResponse (empty)
var xxx_messageInfo_MsgCreateFarmResponse proto.InternalMessageInfo
func (m *MsgCreateFarmResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateFarmResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgCreateFarmResponse.Marshal(b, m, deterministic) }; return nil, nil }
func (m *MsgCreateFarmResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCreateFarmResponse.Merge(m, src) }
func (m *MsgCreateFarmResponse) XXX_Size() int { return 0 }
func (m *MsgCreateFarmResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgCreateFarmResponse.DiscardUnknown(m) }
func (m *MsgCreateFarmResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgCreateFarmResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgCreateFarmResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgCreateFarmResponse) Size() (n int) { return 0 }
func (m *MsgCreateFarmResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l { preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }; iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy }; return nil
}

// MsgUpdateFarm (authority:string/f1, pool_id:uint64/f2, reward_per_block:bytes/f3, active:bool/f4)
var xxx_messageInfo_MsgUpdateFarm proto.InternalMessageInfo
func (m *MsgUpdateFarm) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgUpdateFarm) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgUpdateFarm.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgUpdateFarm) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgUpdateFarm.Merge(m, src) }
func (m *MsgUpdateFarm) XXX_Size() int { return m.Size() }
func (m *MsgUpdateFarm) XXX_DiscardUnknown() { xxx_messageInfo_MsgUpdateFarm.DiscardUnknown(m) }
func (m *MsgUpdateFarm) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgUpdateFarm) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgUpdateFarm) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.Active { i--; if m.Active { dAtA[i] = 1 } else { dAtA[i] = 0 }; i--; dAtA[i] = 0x20 }
	if !m.RewardPerBlock.IsNil() { bz := marshalBytes(m.RewardPerBlock); if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarint(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x1a } }
	if m.PoolID != 0 { i = encodeVarint(dAtA, i, uint64(m.PoolID)); i--; dAtA[i] = 0x10 }
	if len(m.Authority) > 0 { i -= len(m.Authority); copy(dAtA[i:], m.Authority); i = encodeVarint(dAtA, i, uint64(len(m.Authority))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgUpdateFarm) Size() (n int) {
	if m == nil { return 0 }; var l int
	l = len(m.Authority); if l > 0 { n += 1 + l + sov(uint64(l)) }
	if m.PoolID != 0 { n += 1 + sov(uint64(m.PoolID)) }
	if !m.RewardPerBlock.IsNil() { bz := marshalBytes(m.RewardPerBlock); l = len(bz); if l > 0 { n += 1 + l + sov(uint64(l)) } }
	if m.Active { n += 1 + 1 }
	return n
}
func (m *MsgUpdateFarm) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; postIndex := iNdEx + int(stringLen); if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; m.Authority = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.PoolID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.PoolID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		case 3: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLength }; postIndex := iNdEx + byteLen; if postIndex < 0 || postIndex > l { return io.ErrUnexpectedEOF }; _ = json.Unmarshal(dAtA[iNdEx:postIndex], &m.RewardPerBlock); iNdEx = postIndex
		case 4: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; var v int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; v |= int(b&0x7F) << shift; if b < 0x80 { break } }; m.Active = v != 0
		default: iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}; if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgUpdateFarmResponse (empty) — same as MsgCreateFarmResponse
var xxx_messageInfo_MsgUpdateFarmResponse proto.InternalMessageInfo
func (m *MsgUpdateFarmResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgUpdateFarmResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgUpdateFarmResponse.Marshal(b, m, deterministic) }; return nil, nil }
func (m *MsgUpdateFarmResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgUpdateFarmResponse.Merge(m, src) }
func (m *MsgUpdateFarmResponse) XXX_Size() int { return 0 }
func (m *MsgUpdateFarmResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgUpdateFarmResponse.DiscardUnknown(m) }
func (m *MsgUpdateFarmResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgUpdateFarmResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgUpdateFarmResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgUpdateFarmResponse) Size() (n int) { return 0 }
func (m *MsgUpdateFarmResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l { preIndex := iNdEx; var wire uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflow }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }; iNdEx = preIndex; skippy, err := skip(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLength }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy }; return nil
}
