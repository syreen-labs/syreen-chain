package types

import (
	"fmt"
	"io"
	"math/bits"

	proto "github.com/cosmos/gogoproto/proto"
)

// Reference imports to suppress errors
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = io.ErrUnexpectedEOF

var fileDescriptorFlashloanTx []byte

// Descriptor methods
func (*MsgFlashLoan) Descriptor() ([]byte, []int)                  { return fileDescriptorFlashloanTx, []int{0} }
func (*MsgFlashLoanResponse) Descriptor() ([]byte, []int)          { return fileDescriptorFlashloanTx, []int{1} }
func (*MsgCreateFlashPool) Descriptor() ([]byte, []int)            { return fileDescriptorFlashloanTx, []int{2} }
func (*MsgCreateFlashPoolResponse) Descriptor() ([]byte, []int)    { return fileDescriptorFlashloanTx, []int{3} }
func (*MsgFundFlashPool) Descriptor() ([]byte, []int)              { return fileDescriptorFlashloanTx, []int{4} }
func (*MsgFundFlashPoolResponse) Descriptor() ([]byte, []int)      { return fileDescriptorFlashloanTx, []int{5} }
func (*MsgWithdrawFlashPool) Descriptor() ([]byte, []int)          { return fileDescriptorFlashloanTx, []int{6} }
func (*MsgWithdrawFlashPoolResponse) Descriptor() ([]byte, []int)  { return fileDescriptorFlashloanTx, []int{7} }

func init() {
	proto.RegisterType((*MsgFlashLoan)(nil), "syreen.flashloan.MsgFlashLoan")
	proto.RegisterType((*MsgFlashLoanResponse)(nil), "syreen.flashloan.MsgFlashLoanResponse")
	proto.RegisterType((*MsgCreateFlashPool)(nil), "syreen.flashloan.MsgCreateFlashPool")
	proto.RegisterType((*MsgCreateFlashPoolResponse)(nil), "syreen.flashloan.MsgCreateFlashPoolResponse")
	proto.RegisterType((*MsgFundFlashPool)(nil), "syreen.flashloan.MsgFundFlashPool")
	proto.RegisterType((*MsgFundFlashPoolResponse)(nil), "syreen.flashloan.MsgFundFlashPoolResponse")
	proto.RegisterType((*MsgWithdrawFlashPool)(nil), "syreen.flashloan.MsgWithdrawFlashPool")
	proto.RegisterType((*MsgWithdrawFlashPoolResponse)(nil), "syreen.flashloan.MsgWithdrawFlashPoolResponse")
}

// ---------------------------------------------------------------------------
// MsgFlashLoan — sender(string,1) denom(string,2) amount(bytes,3)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgFlashLoan proto.InternalMessageInfo

func (m *MsgFlashLoan) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgFlashLoan) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgFlashLoan.Marshal(b, m, deterministic) }
	b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil
}
func (m *MsgFlashLoan) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgFlashLoan.Merge(m, src) }
func (m *MsgFlashLoan) XXX_Size() int                { return m.Size() }
func (m *MsgFlashLoan) XXX_DiscardUnknown()          { xxx_messageInfo_MsgFlashLoan.DiscardUnknown(m) }

func (m *MsgFlashLoan) Marshal() (dAtA []byte, err error) {
	size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil
}
func (m *MsgFlashLoan) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgFlashLoan) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	// field 3: amount (bytes — math.Int)
	{
		bz, err := m.Amount.Marshal()
		if err != nil { return 0, err }
		if len(bz) > 0 {
			i -= len(bz); copy(dAtA[i:], bz)
			i = encodeVarintFlashloan(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x1a
		}
	}
	// field 2: denom (string)
	if len(m.Denom) > 0 { i -= len(m.Denom); copy(dAtA[i:], m.Denom); i = encodeVarintFlashloan(dAtA, i, uint64(len(m.Denom))); i--; dAtA[i] = 0x12 }
	// field 1: sender (string)
	if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarintFlashloan(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgFlashLoan) Size() (n int) {
	if m == nil { return 0 }
	l := len(m.Sender); if l > 0 { n += 1 + l + sovFlashloan(uint64(l)) }
	l = len(m.Denom); if l > 0 { n += 1 + l + sovFlashloan(uint64(l)) }
	bz, _ := m.Amount.Marshal(); l = len(bz); if l > 0 { n += 1 + l + sovFlashloan(uint64(l)) }
	return n
}
func (m *MsgFlashLoan) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowFlashloan }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7)
		if wireType == 4 { return fmt.Errorf("proto: MsgFlashLoan: wiretype end group for non-group") }
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgFlashLoan: illegal tag %d (wire type %d)", fieldNum, wire) }
		switch fieldNum {
		case 1: // sender
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowFlashloan }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
			intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthFlashloan }
			postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthFlashloan }; if postIndex > l { return io.ErrUnexpectedEOF }
			m.Sender = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2: // denom
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowFlashloan }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
			intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthFlashloan }
			postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthFlashloan }; if postIndex > l { return io.ErrUnexpectedEOF }
			m.Denom = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 3: // amount
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType) }
			var byteLen int
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowFlashloan }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }
			if byteLen < 0 { return ErrInvalidLengthFlashloan }
			postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthFlashloan }; if postIndex > l { return io.ErrUnexpectedEOF }
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }
			iNdEx = postIndex
		default:
			iNdEx = preIndex; skippy, err := skipFlashloan(dAtA[iNdEx:]); if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthFlashloan }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// ---------------------------------------------------------------------------
// MsgFlashLoanResponse — fee(bytes,1)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgFlashLoanResponse proto.InternalMessageInfo

func (m *MsgFlashLoanResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgFlashLoanResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgFlashLoanResponse.Marshal(b, m, deterministic) }
	b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil
}
func (m *MsgFlashLoanResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgFlashLoanResponse.Merge(m, src) }
func (m *MsgFlashLoanResponse) XXX_Size() int               { return m.Size() }
func (m *MsgFlashLoanResponse) XXX_DiscardUnknown()         { xxx_messageInfo_MsgFlashLoanResponse.DiscardUnknown(m) }

func (m *MsgFlashLoanResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil
}
func (m *MsgFlashLoanResponse) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgFlashLoanResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	bz, err := m.Fee.Marshal(); if err != nil { return 0, err }
	if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintFlashloan(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgFlashLoanResponse) Size() (n int) {
	if m == nil { return 0 }
	bz, _ := m.Fee.Marshal(); l := len(bz); if l > 0 { n += 1 + l + sovFlashloan(uint64(l)) }
	return n
}
func (m *MsgFlashLoanResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowFlashloan }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3)
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgFlashLoanResponse: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType for field Fee") }
			var byteLen int
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowFlashloan }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }
			if byteLen < 0 { return ErrInvalidLengthFlashloan }
			postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthFlashloan }; if postIndex > l { return io.ErrUnexpectedEOF }
			if err := m.Fee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex
		default:
			iNdEx = preIndex; skippy, err := skipFlashloan(dAtA[iNdEx:]); if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthFlashloan }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// ---------------------------------------------------------------------------
// MsgCreateFlashPool — authority(string,1) denom(string,2) fee_rate(bytes,3)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCreateFlashPool proto.InternalMessageInfo

func (m *MsgCreateFlashPool) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateFlashPool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgCreateFlashPool.Marshal(b, m, deterministic) }
	b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil
}
func (m *MsgCreateFlashPool) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCreateFlashPool.Merge(m, src) }
func (m *MsgCreateFlashPool) XXX_Size() int               { return m.Size() }
func (m *MsgCreateFlashPool) XXX_DiscardUnknown()         { xxx_messageInfo_MsgCreateFlashPool.DiscardUnknown(m) }

func (m *MsgCreateFlashPool) Marshal() (dAtA []byte, err error) {
	size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil
}
func (m *MsgCreateFlashPool) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgCreateFlashPool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	// field 3: fee_rate (bytes — LegacyDec)
	{
		bz, err := m.FeeRate.Marshal()
		if err != nil { return 0, err }
		if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintFlashloan(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x1a }
	}
	if len(m.Denom) > 0 { i -= len(m.Denom); copy(dAtA[i:], m.Denom); i = encodeVarintFlashloan(dAtA, i, uint64(len(m.Denom))); i--; dAtA[i] = 0x12 }
	if len(m.Authority) > 0 { i -= len(m.Authority); copy(dAtA[i:], m.Authority); i = encodeVarintFlashloan(dAtA, i, uint64(len(m.Authority))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgCreateFlashPool) Size() (n int) {
	if m == nil { return 0 }
	l := len(m.Authority); if l > 0 { n += 1 + l + sovFlashloan(uint64(l)) }
	l = len(m.Denom); if l > 0 { n += 1 + l + sovFlashloan(uint64(l)) }
	bz, _ := m.FeeRate.Marshal(); l = len(bz); if l > 0 { n += 1 + l + sovFlashloan(uint64(l)) }
	return n
}
func (m *MsgCreateFlashPool) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowFlashloan }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7)
		if wireType == 4 { return fmt.Errorf("proto: MsgCreateFlashPool: wiretype end group for non-group") }
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgCreateFlashPool: illegal tag %d", fieldNum, wire) }
		switch fieldNum {
		case 1:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Authority", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowFlashloan }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
			intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthFlashloan }
			postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthFlashloan }; if postIndex > l { return io.ErrUnexpectedEOF }
			m.Authority = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowFlashloan }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
			intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthFlashloan }
			postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthFlashloan }; if postIndex > l { return io.ErrUnexpectedEOF }
			m.Denom = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 3:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field FeeRate", wireType) }
			var byteLen int
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowFlashloan }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }
			if byteLen < 0 { return ErrInvalidLengthFlashloan }
			postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthFlashloan }; if postIndex > l { return io.ErrUnexpectedEOF }
			if err := m.FeeRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex
		default:
			iNdEx = preIndex; skippy, err := skipFlashloan(dAtA[iNdEx:]); if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthFlashloan }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// ---------------------------------------------------------------------------
// MsgCreateFlashPoolResponse (empty)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCreateFlashPoolResponse proto.InternalMessageInfo

func (m *MsgCreateFlashPoolResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateFlashPoolResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgCreateFlashPoolResponse.Marshal(b, m, deterministic) }
	b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil
}
func (m *MsgCreateFlashPoolResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCreateFlashPoolResponse.Merge(m, src) }
func (m *MsgCreateFlashPoolResponse) XXX_Size() int               { return m.Size() }
func (m *MsgCreateFlashPoolResponse) XXX_DiscardUnknown()         { xxx_messageInfo_MsgCreateFlashPoolResponse.DiscardUnknown(m) }

func (m *MsgCreateFlashPoolResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgCreateFlashPoolResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgCreateFlashPoolResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgCreateFlashPoolResponse) Size() (n int) { return 0 }
func (m *MsgCreateFlashPoolResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowFlashloan }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: MsgCreateFlashPoolResponse: illegal tag %d", fieldNum) }
		iNdEx = preIndex; skippy, err := skipFlashloan(dAtA[iNdEx:]); if err != nil { return err }
		if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthFlashloan }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgFundFlashPool — sender(string,1) denom(string,2) amount(bytes,3)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgFundFlashPool proto.InternalMessageInfo

func (m *MsgFundFlashPool) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgFundFlashPool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgFundFlashPool.Marshal(b, m, deterministic) }
	b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil
}
func (m *MsgFundFlashPool) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgFundFlashPool.Merge(m, src) }
func (m *MsgFundFlashPool) XXX_Size() int               { return m.Size() }
func (m *MsgFundFlashPool) XXX_DiscardUnknown()         { xxx_messageInfo_MsgFundFlashPool.DiscardUnknown(m) }

func (m *MsgFundFlashPool) Marshal() (dAtA []byte, err error) {
	size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil
}
func (m *MsgFundFlashPool) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgFundFlashPool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	{ bz, err := m.Amount.Marshal(); if err != nil { return 0, err }; if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintFlashloan(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x1a } }
	if len(m.Denom) > 0 { i -= len(m.Denom); copy(dAtA[i:], m.Denom); i = encodeVarintFlashloan(dAtA, i, uint64(len(m.Denom))); i--; dAtA[i] = 0x12 }
	if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarintFlashloan(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgFundFlashPool) Size() (n int) {
	if m == nil { return 0 }
	l := len(m.Sender); if l > 0 { n += 1 + l + sovFlashloan(uint64(l)) }
	l = len(m.Denom); if l > 0 { n += 1 + l + sovFlashloan(uint64(l)) }
	bz, _ := m.Amount.Marshal(); l = len(bz); if l > 0 { n += 1 + l + sovFlashloan(uint64(l)) }
	return n
}
func (m *MsgFundFlashPool) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowFlashloan }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7)
		if wireType == 4 { return fmt.Errorf("proto: MsgFundFlashPool: wiretype end group") }
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgFundFlashPool: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowFlashloan }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
			intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthFlashloan }
			postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthFlashloan }; if postIndex > l { return io.ErrUnexpectedEOF }
			m.Sender = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowFlashloan }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
			intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthFlashloan }
			postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthFlashloan }; if postIndex > l { return io.ErrUnexpectedEOF }
			m.Denom = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 3:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType) }
			var byteLen int
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowFlashloan }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }
			if byteLen < 0 { return ErrInvalidLengthFlashloan }
			postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthFlashloan }; if postIndex > l { return io.ErrUnexpectedEOF }
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex
		default:
			iNdEx = preIndex; skippy, err := skipFlashloan(dAtA[iNdEx:]); if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthFlashloan }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// ---------------------------------------------------------------------------
// MsgFundFlashPoolResponse (empty)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgFundFlashPoolResponse proto.InternalMessageInfo

func (m *MsgFundFlashPoolResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgFundFlashPoolResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgFundFlashPoolResponse.Marshal(b, m, deterministic) }
	b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil
}
func (m *MsgFundFlashPoolResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgFundFlashPoolResponse.Merge(m, src) }
func (m *MsgFundFlashPoolResponse) XXX_Size() int               { return m.Size() }
func (m *MsgFundFlashPoolResponse) XXX_DiscardUnknown()         { xxx_messageInfo_MsgFundFlashPoolResponse.DiscardUnknown(m) }

func (m *MsgFundFlashPoolResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgFundFlashPoolResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgFundFlashPoolResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgFundFlashPoolResponse) Size() (n int) { return 0 }
func (m *MsgFundFlashPoolResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowFlashloan }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: MsgFundFlashPoolResponse: illegal tag %d", fieldNum) }
		iNdEx = preIndex; skippy, err := skipFlashloan(dAtA[iNdEx:]); if err != nil { return err }
		if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthFlashloan }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgWithdrawFlashPool — sender(string,1) denom(string,2) shares(bytes,3)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgWithdrawFlashPool proto.InternalMessageInfo

func (m *MsgWithdrawFlashPool) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgWithdrawFlashPool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgWithdrawFlashPool.Marshal(b, m, deterministic) }
	b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil
}
func (m *MsgWithdrawFlashPool) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgWithdrawFlashPool.Merge(m, src) }
func (m *MsgWithdrawFlashPool) XXX_Size() int               { return m.Size() }
func (m *MsgWithdrawFlashPool) XXX_DiscardUnknown()         { xxx_messageInfo_MsgWithdrawFlashPool.DiscardUnknown(m) }

func (m *MsgWithdrawFlashPool) Marshal() (dAtA []byte, err error) {
	size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil
}
func (m *MsgWithdrawFlashPool) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgWithdrawFlashPool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	{ bz, err := m.Shares.Marshal(); if err != nil { return 0, err }; if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintFlashloan(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x1a } }
	if len(m.Denom) > 0 { i -= len(m.Denom); copy(dAtA[i:], m.Denom); i = encodeVarintFlashloan(dAtA, i, uint64(len(m.Denom))); i--; dAtA[i] = 0x12 }
	if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarintFlashloan(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgWithdrawFlashPool) Size() (n int) {
	if m == nil { return 0 }
	l := len(m.Sender); if l > 0 { n += 1 + l + sovFlashloan(uint64(l)) }
	l = len(m.Denom); if l > 0 { n += 1 + l + sovFlashloan(uint64(l)) }
	bz, _ := m.Shares.Marshal(); l = len(bz); if l > 0 { n += 1 + l + sovFlashloan(uint64(l)) }
	return n
}
func (m *MsgWithdrawFlashPool) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowFlashloan }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7)
		if wireType == 4 { return fmt.Errorf("proto: MsgWithdrawFlashPool: wiretype end group") }
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgWithdrawFlashPool: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowFlashloan }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
			intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthFlashloan }
			postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthFlashloan }; if postIndex > l { return io.ErrUnexpectedEOF }
			m.Sender = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowFlashloan }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
			intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthFlashloan }
			postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthFlashloan }; if postIndex > l { return io.ErrUnexpectedEOF }
			m.Denom = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 3:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Shares", wireType) }
			var byteLen int
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowFlashloan }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }
			if byteLen < 0 { return ErrInvalidLengthFlashloan }
			postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthFlashloan }; if postIndex > l { return io.ErrUnexpectedEOF }
			if err := m.Shares.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex
		default:
			iNdEx = preIndex; skippy, err := skipFlashloan(dAtA[iNdEx:]); if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthFlashloan }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// ---------------------------------------------------------------------------
// MsgWithdrawFlashPoolResponse — amount_returned(bytes,1)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgWithdrawFlashPoolResponse proto.InternalMessageInfo

func (m *MsgWithdrawFlashPoolResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgWithdrawFlashPoolResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgWithdrawFlashPoolResponse.Marshal(b, m, deterministic) }
	b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil
}
func (m *MsgWithdrawFlashPoolResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgWithdrawFlashPoolResponse.Merge(m, src) }
func (m *MsgWithdrawFlashPoolResponse) XXX_Size() int               { return m.Size() }
func (m *MsgWithdrawFlashPoolResponse) XXX_DiscardUnknown()         { xxx_messageInfo_MsgWithdrawFlashPoolResponse.DiscardUnknown(m) }

func (m *MsgWithdrawFlashPoolResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil
}
func (m *MsgWithdrawFlashPoolResponse) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgWithdrawFlashPoolResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	bz, err := m.AmountReturned.Marshal(); if err != nil { return 0, err }
	if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintFlashloan(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgWithdrawFlashPoolResponse) Size() (n int) {
	if m == nil { return 0 }
	bz, _ := m.AmountReturned.Marshal(); l := len(bz); if l > 0 { n += 1 + l + sovFlashloan(uint64(l)) }
	return n
}
func (m *MsgWithdrawFlashPoolResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowFlashloan }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3)
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgWithdrawFlashPoolResponse: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType for field AmountReturned") }
			var byteLen int
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowFlashloan }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }
			if byteLen < 0 { return ErrInvalidLengthFlashloan }
			postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthFlashloan }; if postIndex > l { return io.ErrUnexpectedEOF }
			if err := m.AmountReturned.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex
		default:
			iNdEx = preIndex; skippy, err := skipFlashloan(dAtA[iNdEx:]); if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthFlashloan }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// ---------------------------------------------------------------------------
// Helper functions
// ---------------------------------------------------------------------------

var (
	ErrIntOverflowFlashloan  = fmt.Errorf("proto: integer overflow")
	ErrInvalidLengthFlashloan = fmt.Errorf("proto: negative length found during unmarshaling")
)

func encodeVarintFlashloan(dAtA []byte, offset int, v uint64) int {
	offset -= sovFlashloan(v)
	base := offset
	for v >= 1<<7 { dAtA[offset] = uint8(v&0x7f | 0x80); v >>= 7; offset++ }
	dAtA[offset] = uint8(v)
	return base
}

func sovFlashloan(x uint64) int {
	return (bits.Len64(x|1) + 6) / 7
}

func skipFlashloan(dAtA []byte) (int, error) {
	l := len(dAtA); iNdEx := 0; depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return 0, ErrIntOverflowFlashloan }; if iNdEx >= l { return 0, io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= (uint64(b) & 0x7F) << shift; if b < 0x80 { break } }
		wireType := int(wire & 0x7)
		switch wireType {
		case 0: for shift := uint(0); ; shift += 7 { if shift >= 64 { return 0, ErrIntOverflowFlashloan }; if iNdEx >= l { return 0, io.ErrUnexpectedEOF }; iNdEx++; if dAtA[iNdEx-1] < 0x80 { break } }
		case 1: iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return 0, ErrIntOverflowFlashloan }; if iNdEx >= l { return 0, io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; length |= (int(b) & 0x7F) << shift; if b < 0x80 { break } }
			if length < 0 { return 0, ErrInvalidLengthFlashloan }; iNdEx += length
		case 3: depth++
		case 4: if depth == 0 { return 0, fmt.Errorf("proto: wiretype end group for non-group") }; depth--
		case 5: iNdEx += 4
		default: return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 { return 0, ErrInvalidLengthFlashloan }
		if depth == 0 { return iNdEx, nil }
	}
	return 0, io.ErrUnexpectedEOF
}
