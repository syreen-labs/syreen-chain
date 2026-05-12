package types

import (
	"encoding/json"
	"fmt"
	"io"
	"math/bits"
	"time"

	proto "github.com/cosmos/gogoproto/proto"
)

// Reference imports to suppress errors
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = io.ErrUnexpectedEOF

// Descriptor methods - return the gzipped file descriptor and the message index
func (*MsgCreateSmartAccount) Descriptor() ([]byte, []int)          { return fileDescriptorTx, []int{0} }
func (*MsgCreateSmartAccountResponse) Descriptor() ([]byte, []int)  { return fileDescriptorTx, []int{1} }
func (*MsgCreateSessionKey) Descriptor() ([]byte, []int)            { return fileDescriptorTx, []int{2} }
func (*MsgCreateSessionKeyResponse) Descriptor() ([]byte, []int)    { return fileDescriptorTx, []int{3} }
func (*MsgRevokeSessionKey) Descriptor() ([]byte, []int)            { return fileDescriptorTx, []int{4} }
func (*MsgRevokeSessionKeyResponse) Descriptor() ([]byte, []int)    { return fileDescriptorTx, []int{5} }
func (*MsgInitiateRecovery) Descriptor() ([]byte, []int)            { return fileDescriptorTx, []int{6} }
func (*MsgInitiateRecoveryResponse) Descriptor() ([]byte, []int)    { return fileDescriptorTx, []int{7} }
func (*MsgApproveRecovery) Descriptor() ([]byte, []int)             { return fileDescriptorTx, []int{8} }
func (*MsgApproveRecoveryResponse) Descriptor() ([]byte, []int)     { return fileDescriptorTx, []int{9} }
func (*MsgExecuteRecovery) Descriptor() ([]byte, []int)             { return fileDescriptorTx, []int{10} }
func (*MsgExecuteRecoveryResponse) Descriptor() ([]byte, []int)     { return fileDescriptorTx, []int{11} }
func (*MsgSponsorGas) Descriptor() ([]byte, []int)                  { return fileDescriptorTx, []int{12} }
func (*MsgSponsorGasResponse) Descriptor() ([]byte, []int)          { return fileDescriptorTx, []int{13} }
func (*MsgBatchExecute) Descriptor() ([]byte, []int)                { return fileDescriptorTx, []int{14} }
func (*MsgBatchExecuteResponse) Descriptor() ([]byte, []int)        { return fileDescriptorTx, []int{15} }

func init() {
	proto.RegisterType((*MsgCreateSmartAccount)(nil), "syreen.abstractaccount.MsgCreateSmartAccount")
	proto.RegisterType((*MsgCreateSmartAccountResponse)(nil), "syreen.abstractaccount.MsgCreateSmartAccountResponse")
	proto.RegisterType((*MsgCreateSessionKey)(nil), "syreen.abstractaccount.MsgCreateSessionKey")
	proto.RegisterType((*MsgCreateSessionKeyResponse)(nil), "syreen.abstractaccount.MsgCreateSessionKeyResponse")
	proto.RegisterType((*MsgRevokeSessionKey)(nil), "syreen.abstractaccount.MsgRevokeSessionKey")
	proto.RegisterType((*MsgRevokeSessionKeyResponse)(nil), "syreen.abstractaccount.MsgRevokeSessionKeyResponse")
	proto.RegisterType((*MsgInitiateRecovery)(nil), "syreen.abstractaccount.MsgInitiateRecovery")
	proto.RegisterType((*MsgInitiateRecoveryResponse)(nil), "syreen.abstractaccount.MsgInitiateRecoveryResponse")
	proto.RegisterType((*MsgApproveRecovery)(nil), "syreen.abstractaccount.MsgApproveRecovery")
	proto.RegisterType((*MsgApproveRecoveryResponse)(nil), "syreen.abstractaccount.MsgApproveRecoveryResponse")
	proto.RegisterType((*MsgExecuteRecovery)(nil), "syreen.abstractaccount.MsgExecuteRecovery")
	proto.RegisterType((*MsgExecuteRecoveryResponse)(nil), "syreen.abstractaccount.MsgExecuteRecoveryResponse")
	proto.RegisterType((*MsgSponsorGas)(nil), "syreen.abstractaccount.MsgSponsorGas")
	proto.RegisterType((*MsgSponsorGasResponse)(nil), "syreen.abstractaccount.MsgSponsorGasResponse")
	proto.RegisterType((*MsgBatchExecute)(nil), "syreen.abstractaccount.MsgBatchExecute")
	proto.RegisterType((*MsgBatchExecuteResponse)(nil), "syreen.abstractaccount.MsgBatchExecuteResponse")
}

// ---------------------------------------------------------------------------
// MsgCreateSmartAccount
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCreateSmartAccount proto.InternalMessageInfo

func (m *MsgCreateSmartAccount) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateSmartAccount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateSmartAccount.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgCreateSmartAccount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCreateSmartAccount.Merge(m, src)
}
func (m *MsgCreateSmartAccount) XXX_Size() int       { return m.Size() }
func (m *MsgCreateSmartAccount) XXX_DiscardUnknown() { xxx_messageInfo_MsgCreateSmartAccount.DiscardUnknown(m) }

func (m *MsgCreateSmartAccount) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCreateSmartAccount) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCreateSmartAccount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.Threshold != 0 {
		i = encodeVarint(dAtA, i, uint64(m.Threshold))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Owners) > 0 {
		for iNdEx := len(m.Owners) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Owners[iNdEx])
			copy(dAtA[i:], m.Owners[iNdEx])
			i = encodeVarint(dAtA, i, uint64(len(m.Owners[iNdEx])))
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.AccountType) > 0 {
		i -= len(m.AccountType)
		copy(dAtA[i:], m.AccountType)
		i = encodeVarint(dAtA, i, uint64(len(m.AccountType)))
		i--
		dAtA[i] = 0x12
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

func (m *MsgCreateSmartAccount) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	l = len(m.AccountType)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	if len(m.Owners) > 0 {
		for _, s := range m.Owners {
			l = len(s)
			n += 1 + l + sov(uint64(l))
		}
	}
	if m.Threshold != 0 {
		n += 1 + sov(uint64(m.Threshold))
	}
	return n
}

func (m *MsgCreateSmartAccount) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgCreateSmartAccount: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCreateSmartAccount: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1: // sender
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2: // account_type
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AccountType", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AccountType = AccountType(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3: // owners
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owners", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owners = append(m.Owners, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 4: // threshold
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Threshold", wireType)
			}
			m.Threshold = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Threshold |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgCreateSmartAccountResponse
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCreateSmartAccountResponse proto.InternalMessageInfo

func (m *MsgCreateSmartAccountResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateSmartAccountResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateSmartAccountResponse.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgCreateSmartAccountResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCreateSmartAccountResponse.Merge(m, src)
}
func (m *MsgCreateSmartAccountResponse) XXX_Size() int { return m.Size() }
func (m *MsgCreateSmartAccountResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgCreateSmartAccountResponse.DiscardUnknown(m)
}

func (m *MsgCreateSmartAccountResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}
func (m *MsgCreateSmartAccountResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgCreateSmartAccountResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarint(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}
func (m *MsgCreateSmartAccountResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	l := len(m.Address)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	return n
}
func (m *MsgCreateSmartAccountResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgCreateSmartAccountResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCreateSmartAccountResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgCreateSessionKey
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCreateSessionKey proto.InternalMessageInfo

func (m *MsgCreateSessionKey) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateSessionKey) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateSessionKey.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgCreateSessionKey) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCreateSessionKey.Merge(m, src)
}
func (m *MsgCreateSessionKey) XXX_Size() int       { return m.Size() }
func (m *MsgCreateSessionKey) XXX_DiscardUnknown() { xxx_messageInfo_MsgCreateSessionKey.DiscardUnknown(m) }

func (m *MsgCreateSessionKey) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}
func (m *MsgCreateSessionKey) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgCreateSessionKey) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.Duration != 0 {
		i = encodeVarint(dAtA, i, uint64(m.Duration))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Permissions) > 0 {
		for iNdEx := len(m.Permissions) - 1; iNdEx >= 0; iNdEx-- {
			bz, err := json.Marshal(m.Permissions[iNdEx])
			if err != nil {
				return 0, err
			}
			i -= len(bz)
			copy(dAtA[i:], bz)
			i = encodeVarint(dAtA, i, uint64(len(bz)))
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Grantee) > 0 {
		i -= len(m.Grantee)
		copy(dAtA[i:], m.Grantee)
		i = encodeVarint(dAtA, i, uint64(len(m.Grantee)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Granter) > 0 {
		i -= len(m.Granter)
		copy(dAtA[i:], m.Granter)
		i = encodeVarint(dAtA, i, uint64(len(m.Granter)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}
func (m *MsgCreateSessionKey) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	l = len(m.Granter)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	l = len(m.Grantee)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	if len(m.Permissions) > 0 {
		for _, p := range m.Permissions {
			bz, _ := json.Marshal(p)
			l = len(bz)
			n += 1 + l + sov(uint64(l))
		}
	}
	if m.Duration != 0 {
		n += 1 + sov(uint64(m.Duration))
	}
	return n
}
func (m *MsgCreateSessionKey) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgCreateSessionKey: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCreateSessionKey: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Granter", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Granter = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Grantee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Grantee = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Permissions", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var perm Permission
			if err := json.Unmarshal(dAtA[iNdEx:postIndex], &perm); err != nil {
				return err
			}
			m.Permissions = append(m.Permissions, perm)
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Duration", wireType)
			}
			m.Duration = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Duration |= time.Duration(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgCreateSessionKeyResponse (empty)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCreateSessionKeyResponse proto.InternalMessageInfo

func (m *MsgCreateSessionKeyResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateSessionKeyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateSessionKeyResponse.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgCreateSessionKeyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCreateSessionKeyResponse.Merge(m, src)
}
func (m *MsgCreateSessionKeyResponse) XXX_Size() int { return m.Size() }
func (m *MsgCreateSessionKeyResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgCreateSessionKeyResponse.DiscardUnknown(m)
}
func (m *MsgCreateSessionKeyResponse) Marshal() (dAtA []byte, err error)              { return nil, nil }
func (m *MsgCreateSessionKeyResponse) MarshalTo(dAtA []byte) (int, error)              { return 0, nil }
func (m *MsgCreateSessionKeyResponse) MarshalToSizedBuffer(dAtA []byte) (int, error)   { return len(dAtA), nil }
func (m *MsgCreateSessionKeyResponse) Size() (n int)                                   { return 0 }
func (m *MsgCreateSessionKeyResponse) Unmarshal(dAtA []byte) error                     { return unmarshalEmpty(dAtA, "MsgCreateSessionKeyResponse") }

// ---------------------------------------------------------------------------
// MsgRevokeSessionKey (2 string fields)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgRevokeSessionKey proto.InternalMessageInfo

func (m *MsgRevokeSessionKey) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRevokeSessionKey) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRevokeSessionKey.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgRevokeSessionKey) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgRevokeSessionKey.Merge(m, src)
}
func (m *MsgRevokeSessionKey) XXX_Size() int       { return m.Size() }
func (m *MsgRevokeSessionKey) XXX_DiscardUnknown() { xxx_messageInfo_MsgRevokeSessionKey.DiscardUnknown(m) }

func (m *MsgRevokeSessionKey) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}
func (m *MsgRevokeSessionKey) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgRevokeSessionKey) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.SessionKeyAddr) > 0 {
		i -= len(m.SessionKeyAddr)
		copy(dAtA[i:], m.SessionKeyAddr)
		i = encodeVarint(dAtA, i, uint64(len(m.SessionKeyAddr)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Granter) > 0 {
		i -= len(m.Granter)
		copy(dAtA[i:], m.Granter)
		i = encodeVarint(dAtA, i, uint64(len(m.Granter)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}
func (m *MsgRevokeSessionKey) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	l = len(m.Granter)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	l = len(m.SessionKeyAddr)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	return n
}
func (m *MsgRevokeSessionKey) Unmarshal(dAtA []byte) error {
	return unmarshalTwoStrings(dAtA, &m.Granter, &m.SessionKeyAddr, "MsgRevokeSessionKey")
}

// ---------------------------------------------------------------------------
// MsgRevokeSessionKeyResponse (empty)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgRevokeSessionKeyResponse proto.InternalMessageInfo

func (m *MsgRevokeSessionKeyResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRevokeSessionKeyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRevokeSessionKeyResponse.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgRevokeSessionKeyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgRevokeSessionKeyResponse.Merge(m, src)
}
func (m *MsgRevokeSessionKeyResponse) XXX_Size() int { return m.Size() }
func (m *MsgRevokeSessionKeyResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgRevokeSessionKeyResponse.DiscardUnknown(m)
}
func (m *MsgRevokeSessionKeyResponse) Marshal() (dAtA []byte, err error)              { return nil, nil }
func (m *MsgRevokeSessionKeyResponse) MarshalTo(dAtA []byte) (int, error)              { return 0, nil }
func (m *MsgRevokeSessionKeyResponse) MarshalToSizedBuffer(dAtA []byte) (int, error)   { return len(dAtA), nil }
func (m *MsgRevokeSessionKeyResponse) Size() (n int)                                   { return 0 }
func (m *MsgRevokeSessionKeyResponse) Unmarshal(dAtA []byte) error                     { return unmarshalEmpty(dAtA, "MsgRevokeSessionKeyResponse") }

// ---------------------------------------------------------------------------
// MsgInitiateRecovery (guardian, account, new_owners)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgInitiateRecovery proto.InternalMessageInfo

func (m *MsgInitiateRecovery) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgInitiateRecovery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgInitiateRecovery.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgInitiateRecovery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgInitiateRecovery.Merge(m, src)
}
func (m *MsgInitiateRecovery) XXX_Size() int       { return m.Size() }
func (m *MsgInitiateRecovery) XXX_DiscardUnknown() { xxx_messageInfo_MsgInitiateRecovery.DiscardUnknown(m) }

func (m *MsgInitiateRecovery) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}
func (m *MsgInitiateRecovery) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgInitiateRecovery) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.NewOwners) > 0 {
		for iNdEx := len(m.NewOwners) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.NewOwners[iNdEx])
			copy(dAtA[i:], m.NewOwners[iNdEx])
			i = encodeVarint(dAtA, i, uint64(len(m.NewOwners[iNdEx])))
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Account) > 0 {
		i -= len(m.Account)
		copy(dAtA[i:], m.Account)
		i = encodeVarint(dAtA, i, uint64(len(m.Account)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Guardian) > 0 {
		i -= len(m.Guardian)
		copy(dAtA[i:], m.Guardian)
		i = encodeVarint(dAtA, i, uint64(len(m.Guardian)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}
func (m *MsgInitiateRecovery) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	l = len(m.Guardian)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	l = len(m.Account)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	if len(m.NewOwners) > 0 {
		for _, s := range m.NewOwners {
			l = len(s)
			n += 1 + l + sov(uint64(l))
		}
	}
	return n
}
func (m *MsgInitiateRecovery) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgInitiateRecovery: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgInitiateRecovery: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Guardian", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Guardian = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Account", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Account = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NewOwners", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NewOwners = append(m.NewOwners, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgInitiateRecoveryResponse (empty)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgInitiateRecoveryResponse proto.InternalMessageInfo

func (m *MsgInitiateRecoveryResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgInitiateRecoveryResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgInitiateRecoveryResponse.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgInitiateRecoveryResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgInitiateRecoveryResponse.Merge(m, src)
}
func (m *MsgInitiateRecoveryResponse) XXX_Size() int { return m.Size() }
func (m *MsgInitiateRecoveryResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgInitiateRecoveryResponse.DiscardUnknown(m)
}
func (m *MsgInitiateRecoveryResponse) Marshal() (dAtA []byte, err error)              { return nil, nil }
func (m *MsgInitiateRecoveryResponse) MarshalTo(dAtA []byte) (int, error)              { return 0, nil }
func (m *MsgInitiateRecoveryResponse) MarshalToSizedBuffer(dAtA []byte) (int, error)   { return len(dAtA), nil }
func (m *MsgInitiateRecoveryResponse) Size() (n int)                                   { return 0 }
func (m *MsgInitiateRecoveryResponse) Unmarshal(dAtA []byte) error                     { return unmarshalEmpty(dAtA, "MsgInitiateRecoveryResponse") }

// ---------------------------------------------------------------------------
// MsgApproveRecovery (2 string fields)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgApproveRecovery proto.InternalMessageInfo

func (m *MsgApproveRecovery) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgApproveRecovery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgApproveRecovery.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgApproveRecovery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgApproveRecovery.Merge(m, src)
}
func (m *MsgApproveRecovery) XXX_Size() int       { return m.Size() }
func (m *MsgApproveRecovery) XXX_DiscardUnknown() { xxx_messageInfo_MsgApproveRecovery.DiscardUnknown(m) }

func (m *MsgApproveRecovery) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}
func (m *MsgApproveRecovery) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgApproveRecovery) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.Account) > 0 {
		i -= len(m.Account)
		copy(dAtA[i:], m.Account)
		i = encodeVarint(dAtA, i, uint64(len(m.Account)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Guardian) > 0 {
		i -= len(m.Guardian)
		copy(dAtA[i:], m.Guardian)
		i = encodeVarint(dAtA, i, uint64(len(m.Guardian)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}
func (m *MsgApproveRecovery) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	l = len(m.Guardian)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	l = len(m.Account)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	return n
}
func (m *MsgApproveRecovery) Unmarshal(dAtA []byte) error {
	return unmarshalTwoStrings(dAtA, &m.Guardian, &m.Account, "MsgApproveRecovery")
}

// ---------------------------------------------------------------------------
// MsgApproveRecoveryResponse (empty)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgApproveRecoveryResponse proto.InternalMessageInfo

func (m *MsgApproveRecoveryResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgApproveRecoveryResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgApproveRecoveryResponse.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgApproveRecoveryResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgApproveRecoveryResponse.Merge(m, src)
}
func (m *MsgApproveRecoveryResponse) XXX_Size() int { return m.Size() }
func (m *MsgApproveRecoveryResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgApproveRecoveryResponse.DiscardUnknown(m)
}
func (m *MsgApproveRecoveryResponse) Marshal() (dAtA []byte, err error)              { return nil, nil }
func (m *MsgApproveRecoveryResponse) MarshalTo(dAtA []byte) (int, error)              { return 0, nil }
func (m *MsgApproveRecoveryResponse) MarshalToSizedBuffer(dAtA []byte) (int, error)   { return len(dAtA), nil }
func (m *MsgApproveRecoveryResponse) Size() (n int)                                   { return 0 }
func (m *MsgApproveRecoveryResponse) Unmarshal(dAtA []byte) error                     { return unmarshalEmpty(dAtA, "MsgApproveRecoveryResponse") }

// ---------------------------------------------------------------------------
// MsgExecuteRecovery (2 string fields)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgExecuteRecovery proto.InternalMessageInfo

func (m *MsgExecuteRecovery) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgExecuteRecovery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgExecuteRecovery.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgExecuteRecovery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgExecuteRecovery.Merge(m, src)
}
func (m *MsgExecuteRecovery) XXX_Size() int       { return m.Size() }
func (m *MsgExecuteRecovery) XXX_DiscardUnknown() { xxx_messageInfo_MsgExecuteRecovery.DiscardUnknown(m) }

func (m *MsgExecuteRecovery) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}
func (m *MsgExecuteRecovery) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgExecuteRecovery) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.Account) > 0 {
		i -= len(m.Account)
		copy(dAtA[i:], m.Account)
		i = encodeVarint(dAtA, i, uint64(len(m.Account)))
		i--
		dAtA[i] = 0x12
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
func (m *MsgExecuteRecovery) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	l = len(m.Account)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	return n
}
func (m *MsgExecuteRecovery) Unmarshal(dAtA []byte) error {
	return unmarshalTwoStrings(dAtA, &m.Sender, &m.Account, "MsgExecuteRecovery")
}

// ---------------------------------------------------------------------------
// MsgExecuteRecoveryResponse (empty)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgExecuteRecoveryResponse proto.InternalMessageInfo

func (m *MsgExecuteRecoveryResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgExecuteRecoveryResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgExecuteRecoveryResponse.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgExecuteRecoveryResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgExecuteRecoveryResponse.Merge(m, src)
}
func (m *MsgExecuteRecoveryResponse) XXX_Size() int { return m.Size() }
func (m *MsgExecuteRecoveryResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgExecuteRecoveryResponse.DiscardUnknown(m)
}
func (m *MsgExecuteRecoveryResponse) Marshal() (dAtA []byte, err error)              { return nil, nil }
func (m *MsgExecuteRecoveryResponse) MarshalTo(dAtA []byte) (int, error)              { return 0, nil }
func (m *MsgExecuteRecoveryResponse) MarshalToSizedBuffer(dAtA []byte) (int, error)   { return len(dAtA), nil }
func (m *MsgExecuteRecoveryResponse) Size() (n int)                                   { return 0 }
func (m *MsgExecuteRecoveryResponse) Unmarshal(dAtA []byte) error                     { return unmarshalEmpty(dAtA, "MsgExecuteRecoveryResponse") }

// ---------------------------------------------------------------------------
// MsgSponsorGas (2 strings, 1 uint64, 1 int64/Duration)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgSponsorGas proto.InternalMessageInfo

func (m *MsgSponsorGas) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgSponsorGas) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSponsorGas.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgSponsorGas) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSponsorGas.Merge(m, src)
}
func (m *MsgSponsorGas) XXX_Size() int       { return m.Size() }
func (m *MsgSponsorGas) XXX_DiscardUnknown() { xxx_messageInfo_MsgSponsorGas.DiscardUnknown(m) }

func (m *MsgSponsorGas) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}
func (m *MsgSponsorGas) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgSponsorGas) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.Duration != 0 {
		i = encodeVarint(dAtA, i, uint64(m.Duration))
		i--
		dAtA[i] = 0x20
	}
	if m.GasLimit != 0 {
		i = encodeVarint(dAtA, i, m.GasLimit)
		i--
		dAtA[i] = 0x18
	}
	if len(m.Sponsored) > 0 {
		i -= len(m.Sponsored)
		copy(dAtA[i:], m.Sponsored)
		i = encodeVarint(dAtA, i, uint64(len(m.Sponsored)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Sponsor) > 0 {
		i -= len(m.Sponsor)
		copy(dAtA[i:], m.Sponsor)
		i = encodeVarint(dAtA, i, uint64(len(m.Sponsor)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}
func (m *MsgSponsorGas) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	l = len(m.Sponsor)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	l = len(m.Sponsored)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	if m.GasLimit != 0 {
		n += 1 + sov(m.GasLimit)
	}
	if m.Duration != 0 {
		n += 1 + sov(uint64(m.Duration))
	}
	return n
}
func (m *MsgSponsorGas) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgSponsorGas: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSponsorGas: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sponsor", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sponsor = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sponsored", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sponsored = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasLimit", wireType)
			}
			m.GasLimit = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GasLimit |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Duration", wireType)
			}
			m.Duration = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Duration |= time.Duration(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgSponsorGasResponse (empty)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgSponsorGasResponse proto.InternalMessageInfo

func (m *MsgSponsorGasResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgSponsorGasResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSponsorGasResponse.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgSponsorGasResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSponsorGasResponse.Merge(m, src)
}
func (m *MsgSponsorGasResponse) XXX_Size() int { return m.Size() }
func (m *MsgSponsorGasResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSponsorGasResponse.DiscardUnknown(m)
}
func (m *MsgSponsorGasResponse) Marshal() (dAtA []byte, err error)              { return nil, nil }
func (m *MsgSponsorGasResponse) MarshalTo(dAtA []byte) (int, error)              { return 0, nil }
func (m *MsgSponsorGasResponse) MarshalToSizedBuffer(dAtA []byte) (int, error)   { return len(dAtA), nil }
func (m *MsgSponsorGasResponse) Size() (n int)                                   { return 0 }
func (m *MsgSponsorGasResponse) Unmarshal(dAtA []byte) error                     { return unmarshalEmpty(dAtA, "MsgSponsorGasResponse") }

// ---------------------------------------------------------------------------
// MsgBatchExecute (sender + messages as bytes)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgBatchExecute proto.InternalMessageInfo

func (m *MsgBatchExecute) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgBatchExecute) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgBatchExecute.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgBatchExecute) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgBatchExecute.Merge(m, src)
}
func (m *MsgBatchExecute) XXX_Size() int       { return m.Size() }
func (m *MsgBatchExecute) XXX_DiscardUnknown() { xxx_messageInfo_MsgBatchExecute.DiscardUnknown(m) }

func (m *MsgBatchExecute) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}
func (m *MsgBatchExecute) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgBatchExecute) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.Messages) > 0 {
		for iNdEx := len(m.Messages) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Messages[iNdEx])
			copy(dAtA[i:], m.Messages[iNdEx])
			i = encodeVarint(dAtA, i, uint64(len(m.Messages[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
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
func (m *MsgBatchExecute) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	if len(m.Messages) > 0 {
		for _, msg := range m.Messages {
			l = len(msg)
			n += 1 + l + sov(uint64(l))
		}
	}
	return n
}
func (m *MsgBatchExecute) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgBatchExecute: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgBatchExecute: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Messages", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Messages = append(m.Messages, json.RawMessage(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgBatchExecuteResponse (empty)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgBatchExecuteResponse proto.InternalMessageInfo

func (m *MsgBatchExecuteResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgBatchExecuteResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgBatchExecuteResponse.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgBatchExecuteResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgBatchExecuteResponse.Merge(m, src)
}
func (m *MsgBatchExecuteResponse) XXX_Size() int { return m.Size() }
func (m *MsgBatchExecuteResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgBatchExecuteResponse.DiscardUnknown(m)
}
func (m *MsgBatchExecuteResponse) Marshal() (dAtA []byte, err error)              { return nil, nil }
func (m *MsgBatchExecuteResponse) MarshalTo(dAtA []byte) (int, error)              { return 0, nil }
func (m *MsgBatchExecuteResponse) MarshalToSizedBuffer(dAtA []byte) (int, error)   { return len(dAtA), nil }
func (m *MsgBatchExecuteResponse) Size() (n int)                                   { return 0 }
func (m *MsgBatchExecuteResponse) Unmarshal(dAtA []byte) error                     { return unmarshalEmpty(dAtA, "MsgBatchExecuteResponse") }

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
				return 0, fmt.Errorf("proto: wiretype end group for non-group")
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

// unmarshalEmpty skips all fields in an empty response message
func unmarshalEmpty(dAtA []byte, msgName string) error {
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
		if fieldNum <= 0 {
			return fmt.Errorf("proto: %s: illegal tag %d", msgName, fieldNum)
		}
		iNdEx = preIndex
		skippy, err := skip(dAtA[iNdEx:])
		if err != nil {
			return err
		}
		if (skippy < 0) || (iNdEx+skippy) < 0 {
			return ErrInvalidLength
		}
		if (iNdEx + skippy) > l {
			return io.ErrUnexpectedEOF
		}
		iNdEx += skippy
	}
	return nil
}

// unmarshalTwoStrings is a helper for messages with exactly 2 string fields
func unmarshalTwoStrings(dAtA []byte, field1 *string, field2 *string, msgName string) error {
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
			return fmt.Errorf("proto: %s: wiretype end group for non-group", msgName)
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: %s: illegal tag %d (wire type %d)", msgName, fieldNum, wire)
		}
		switch fieldNum {
		case 1, 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field %d", wireType, fieldNum)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if fieldNum == 1 {
				*field1 = string(dAtA[iNdEx:postIndex])
			} else {
				*field2 = string(dAtA[iNdEx:postIndex])
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// suppress unused import
var _ json.RawMessage
var _ = time.Second
