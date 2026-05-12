package types

import (
	"fmt"
	"io"
	"math/bits"

	proto "github.com/cosmos/gogoproto/proto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Reference imports to suppress errors
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = io.ErrUnexpectedEOF

// Descriptor methods - return the gzipped file descriptor and the message index
func (*MsgCreateDenom) Descriptor() ([]byte, []int)          { return fileDescriptorTx, []int{0} }
func (*MsgCreateDenomResponse) Descriptor() ([]byte, []int)  { return fileDescriptorTx, []int{1} }
func (*MsgMint) Descriptor() ([]byte, []int)                 { return fileDescriptorTx, []int{2} }
func (*MsgMintResponse) Descriptor() ([]byte, []int)         { return fileDescriptorTx, []int{3} }
func (*MsgBurn) Descriptor() ([]byte, []int)                 { return fileDescriptorTx, []int{4} }
func (*MsgBurnResponse) Descriptor() ([]byte, []int)         { return fileDescriptorTx, []int{5} }
func (*MsgChangeAdmin) Descriptor() ([]byte, []int)          { return fileDescriptorTx, []int{6} }
func (*MsgChangeAdminResponse) Descriptor() ([]byte, []int)  { return fileDescriptorTx, []int{7} }

func init() {
	proto.RegisterType((*MsgCreateDenom)(nil), "syreen.tokenfactory.MsgCreateDenom")
	proto.RegisterType((*MsgCreateDenomResponse)(nil), "syreen.tokenfactory.MsgCreateDenomResponse")
	proto.RegisterType((*MsgMint)(nil), "syreen.tokenfactory.MsgMint")
	proto.RegisterType((*MsgMintResponse)(nil), "syreen.tokenfactory.MsgMintResponse")
	proto.RegisterType((*MsgBurn)(nil), "syreen.tokenfactory.MsgBurn")
	proto.RegisterType((*MsgBurnResponse)(nil), "syreen.tokenfactory.MsgBurnResponse")
	proto.RegisterType((*MsgChangeAdmin)(nil), "syreen.tokenfactory.MsgChangeAdmin")
	proto.RegisterType((*MsgChangeAdminResponse)(nil), "syreen.tokenfactory.MsgChangeAdminResponse")
}

// ---------------------------------------------------------------------------
// MsgCreateDenom
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCreateDenom proto.InternalMessageInfo

func (m *MsgCreateDenom) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateDenom) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateDenom.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgCreateDenom) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgCreateDenom.Merge(m, src) }
func (m *MsgCreateDenom) XXX_Size() int                { return m.Size() }
func (m *MsgCreateDenom) XXX_DiscardUnknown()          { xxx_messageInfo_MsgCreateDenom.DiscardUnknown(m) }

func (m *MsgCreateDenom) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCreateDenom) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCreateDenom) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.Subdenom) > 0 {
		i -= len(m.Subdenom)
		copy(dAtA[i:], m.Subdenom)
		i = encodeVarint(dAtA, i, uint64(len(m.Subdenom)))
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

func (m *MsgCreateDenom) Size() (n int) {
	if m == nil {
		return 0
	}
	l := len(m.Sender)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	l = len(m.Subdenom)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	return n
}

func (m *MsgCreateDenom) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgCreateDenom: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCreateDenom: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Subdenom", wireType)
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
			m.Subdenom = string(dAtA[iNdEx:postIndex])
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
// MsgCreateDenomResponse
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCreateDenomResponse proto.InternalMessageInfo

func (m *MsgCreateDenomResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateDenomResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateDenomResponse.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgCreateDenomResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCreateDenomResponse.Merge(m, src)
}
func (m *MsgCreateDenomResponse) XXX_Size() int       { return m.Size() }
func (m *MsgCreateDenomResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgCreateDenomResponse.DiscardUnknown(m) }

func (m *MsgCreateDenomResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCreateDenomResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCreateDenomResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.NewTokenDenom) > 0 {
		i -= len(m.NewTokenDenom)
		copy(dAtA[i:], m.NewTokenDenom)
		i = encodeVarint(dAtA, i, uint64(len(m.NewTokenDenom)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}

func (m *MsgCreateDenomResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	l := len(m.NewTokenDenom)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	return n
}

func (m *MsgCreateDenomResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgCreateDenomResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCreateDenomResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NewTokenDenom", wireType)
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
			m.NewTokenDenom = string(dAtA[iNdEx:postIndex])
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
// MsgMint
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgMint proto.InternalMessageInfo

func (m *MsgMint) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgMint) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgMint.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgMint) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgMint.Merge(m, src) }
func (m *MsgMint) XXX_Size() int                { return m.Size() }
func (m *MsgMint) XXX_DiscardUnknown()          { xxx_messageInfo_MsgMint.DiscardUnknown(m) }

func (m *MsgMint) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgMint) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgMint) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.MintTo) > 0 {
		i -= len(m.MintTo)
		copy(dAtA[i:], m.MintTo)
		i = encodeVarint(dAtA, i, uint64(len(m.MintTo)))
		i--
		dAtA[i] = 0x1a
	}
	{
		size, err := m.Amount.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarint(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}

func (m *MsgMint) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	l = m.Amount.Size()
	n += 1 + l + sov(uint64(l))
	l = len(m.MintTo)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	return n
}

func (m *MsgMint) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgMint: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgMint: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MintTo", wireType)
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
			m.MintTo = string(dAtA[iNdEx:postIndex])
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
// MsgMintResponse
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgMintResponse proto.InternalMessageInfo

func (m *MsgMintResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgMintResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgMintResponse.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgMintResponse) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgMintResponse.Merge(m, src) }
func (m *MsgMintResponse) XXX_Size() int                { return m.Size() }
func (m *MsgMintResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgMintResponse.DiscardUnknown(m) }
func (m *MsgMintResponse) ProtoMessage()                {}
func (m *MsgMintResponse) Reset()                       { *m = MsgMintResponse{} }
func (m *MsgMintResponse) String() string               { return "mint_response" }
func (m *MsgMintResponse) XXX_MessageName() string      { return "syreen.tokenfactory.MsgMintResponse" }

func (m *MsgMintResponse) Marshal() (dAtA []byte, err error) {
	return nil, nil
}
func (m *MsgMintResponse) MarshalTo(dAtA []byte) (int, error) {
	return 0, nil
}
func (m *MsgMintResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	return len(dAtA), nil
}
func (m *MsgMintResponse) Size() (n int) { return 0 }
func (m *MsgMintResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgMintResponse: illegal tag %d", fieldNum)
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

// ---------------------------------------------------------------------------
// MsgBurn
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgBurn proto.InternalMessageInfo

func (m *MsgBurn) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgBurn) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgBurn.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgBurn) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgBurn.Merge(m, src) }
func (m *MsgBurn) XXX_Size() int                { return m.Size() }
func (m *MsgBurn) XXX_DiscardUnknown()          { xxx_messageInfo_MsgBurn.DiscardUnknown(m) }

func (m *MsgBurn) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgBurn) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgBurn) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.BurnFrom) > 0 {
		i -= len(m.BurnFrom)
		copy(dAtA[i:], m.BurnFrom)
		i = encodeVarint(dAtA, i, uint64(len(m.BurnFrom)))
		i--
		dAtA[i] = 0x1a
	}
	{
		size, err := m.Amount.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarint(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}

func (m *MsgBurn) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	l = m.Amount.Size()
	n += 1 + l + sov(uint64(l))
	l = len(m.BurnFrom)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	return n
}

func (m *MsgBurn) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgBurn: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgBurn: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BurnFrom", wireType)
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
			m.BurnFrom = string(dAtA[iNdEx:postIndex])
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
// MsgBurnResponse
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgBurnResponse proto.InternalMessageInfo

func (m *MsgBurnResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgBurnResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgBurnResponse.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgBurnResponse) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgBurnResponse.Merge(m, src) }
func (m *MsgBurnResponse) XXX_Size() int                { return m.Size() }
func (m *MsgBurnResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgBurnResponse.DiscardUnknown(m) }
func (m *MsgBurnResponse) ProtoMessage()                {}
func (m *MsgBurnResponse) Reset()                       { *m = MsgBurnResponse{} }
func (m *MsgBurnResponse) String() string               { return "burn_response" }
func (m *MsgBurnResponse) XXX_MessageName() string      { return "syreen.tokenfactory.MsgBurnResponse" }

func (m *MsgBurnResponse) Marshal() (dAtA []byte, err error) {
	return nil, nil
}
func (m *MsgBurnResponse) MarshalTo(dAtA []byte) (int, error) {
	return 0, nil
}
func (m *MsgBurnResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	return len(dAtA), nil
}
func (m *MsgBurnResponse) Size() (n int) { return 0 }
func (m *MsgBurnResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgBurnResponse: illegal tag %d", fieldNum)
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

// ---------------------------------------------------------------------------
// MsgChangeAdmin
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgChangeAdmin proto.InternalMessageInfo

func (m *MsgChangeAdmin) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgChangeAdmin) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgChangeAdmin.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgChangeAdmin) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgChangeAdmin.Merge(m, src) }
func (m *MsgChangeAdmin) XXX_Size() int                { return m.Size() }
func (m *MsgChangeAdmin) XXX_DiscardUnknown()          { xxx_messageInfo_MsgChangeAdmin.DiscardUnknown(m) }

func (m *MsgChangeAdmin) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgChangeAdmin) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgChangeAdmin) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.NewAdmin) > 0 {
		i -= len(m.NewAdmin)
		copy(dAtA[i:], m.NewAdmin)
		i = encodeVarint(dAtA, i, uint64(len(m.NewAdmin)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarint(dAtA, i, uint64(len(m.Denom)))
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

func (m *MsgChangeAdmin) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	l = len(m.NewAdmin)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	return n
}

func (m *MsgChangeAdmin) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgChangeAdmin: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgChangeAdmin: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
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
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NewAdmin", wireType)
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
			m.NewAdmin = string(dAtA[iNdEx:postIndex])
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
// MsgChangeAdminResponse
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgChangeAdminResponse proto.InternalMessageInfo

func (m *MsgChangeAdminResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgChangeAdminResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgChangeAdminResponse.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgChangeAdminResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgChangeAdminResponse.Merge(m, src)
}
func (m *MsgChangeAdminResponse) XXX_Size() int       { return m.Size() }
func (m *MsgChangeAdminResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgChangeAdminResponse.DiscardUnknown(m) }
func (m *MsgChangeAdminResponse) ProtoMessage()       {}
func (m *MsgChangeAdminResponse) Reset()              { *m = MsgChangeAdminResponse{} }
func (m *MsgChangeAdminResponse) String() string      { return "change_admin_response" }
func (m *MsgChangeAdminResponse) XXX_MessageName() string {
	return "syreen.tokenfactory.MsgChangeAdminResponse"
}

func (m *MsgChangeAdminResponse) Marshal() (dAtA []byte, err error) {
	return nil, nil
}
func (m *MsgChangeAdminResponse) MarshalTo(dAtA []byte) (int, error) {
	return 0, nil
}
func (m *MsgChangeAdminResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	return len(dAtA), nil
}
func (m *MsgChangeAdminResponse) Size() (n int) { return 0 }
func (m *MsgChangeAdminResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgChangeAdminResponse: illegal tag %d", fieldNum)
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

// ---------------------------------------------------------------------------
// Helper functions (standard protobuf encoding helpers)
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
				return 0, fmt.Errorf("proto: MsgCreateDenom: wiretype end group for non-group")
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

// suppress unused import
var _ sdk.Coin
