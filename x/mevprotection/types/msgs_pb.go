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

// Descriptor methods - return the gzipped file descriptor and the message index
func (*MsgCommitTx) Descriptor() ([]byte, []int)         { return fileDescriptorTx, []int{0} }
func (*MsgCommitTxResponse) Descriptor() ([]byte, []int)  { return fileDescriptorTx, []int{1} }
func (*MsgRevealTx) Descriptor() ([]byte, []int)          { return fileDescriptorTx, []int{2} }
func (*MsgRevealTxResponse) Descriptor() ([]byte, []int)  { return fileDescriptorTx, []int{3} }

func init() {
	proto.RegisterType((*MsgCommitTx)(nil), "syreen.mevprotection.MsgCommitTx")
	proto.RegisterType((*MsgCommitTxResponse)(nil), "syreen.mevprotection.MsgCommitTxResponse")
	proto.RegisterType((*MsgRevealTx)(nil), "syreen.mevprotection.MsgRevealTx")
	proto.RegisterType((*MsgRevealTxResponse)(nil), "syreen.mevprotection.MsgRevealTxResponse")
}

// Marshal / Unmarshal for MsgCommitTx
func (m *MsgCommitTx) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCommitTx) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCommitTx) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.EncryptedTx) > 0 {
		i -= len(m.EncryptedTx)
		copy(dAtA[i:], m.EncryptedTx)
		i = encodeVarint(dAtA, i, uint64(len(m.EncryptedTx)))
		i--
		dAtA[i] = 0x1a // field 3, type bytes
	}
	if len(m.TxHash) > 0 {
		i -= len(m.TxHash)
		copy(dAtA[i:], m.TxHash)
		i = encodeVarint(dAtA, i, uint64(len(m.TxHash)))
		i--
		dAtA[i] = 0x12 // field 2, type bytes
	}
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarint(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0x0a // field 1, type string
	}
	return len(dAtA) - i, nil
}

func (m *MsgCommitTx) Size() int {
	var n int
	l := len(m.Sender)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	l = len(m.TxHash)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	l = len(m.EncryptedTx)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	return n
}

func (m *MsgCommitTx) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
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
		if wireType == 2 {
			var byteLen int
			for shift := uint(0); ; shift += 7 {
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
			if byteLen < 0 || iNdEx+byteLen > l {
				return io.ErrUnexpectedEOF
			}
			switch fieldNum {
			case 1:
				m.Sender = string(dAtA[iNdEx : iNdEx+byteLen])
			case 2:
				m.TxHash = append(m.TxHash[:0], dAtA[iNdEx:iNdEx+byteLen]...)
			case 3:
				m.EncryptedTx = append(m.EncryptedTx[:0], dAtA[iNdEx:iNdEx+byteLen]...)
			}
			iNdEx += byteLen
		} else {
			if wireType < 0 || wireType > 5 {
				return fmt.Errorf("proto: MsgCommitTx: illegal tag %d (wire type %d)", fieldNum, wireType)
			}
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			iNdEx += skippy
		}
		_ = preIndex
	}
	return nil
}

// Marshal / Unmarshal for MsgCommitTxResponse
func (m *MsgCommitTxResponse) Marshal() (dAtA []byte, err error) {
	return []byte{}, nil
}
func (m *MsgCommitTxResponse) MarshalTo(dAtA []byte) (int, error)            { return 0, nil }
func (m *MsgCommitTxResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgCommitTxResponse) Size() int                                     { return 0 }
func (m *MsgCommitTxResponse) Unmarshal(dAtA []byte) error                   { return nil }

// Marshal / Unmarshal for MsgRevealTx
func (m *MsgRevealTx) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgRevealTx) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgRevealTx) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.Nonce) > 0 {
		i -= len(m.Nonce)
		copy(dAtA[i:], m.Nonce)
		i = encodeVarint(dAtA, i, uint64(len(m.Nonce)))
		i--
		dAtA[i] = 0x22 // field 4, type bytes
	}
	if len(m.TxBody) > 0 {
		i -= len(m.TxBody)
		copy(dAtA[i:], m.TxBody)
		i = encodeVarint(dAtA, i, uint64(len(m.TxBody)))
		i--
		dAtA[i] = 0x1a // field 3, type bytes
	}
	if len(m.CommitHash) > 0 {
		i -= len(m.CommitHash)
		copy(dAtA[i:], m.CommitHash)
		i = encodeVarint(dAtA, i, uint64(len(m.CommitHash)))
		i--
		dAtA[i] = 0x12 // field 2, type bytes
	}
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarint(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0x0a // field 1, type string
	}
	return len(dAtA) - i, nil
}

func (m *MsgRevealTx) Size() int {
	var n int
	l := len(m.Sender)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	l = len(m.CommitHash)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	l = len(m.TxBody)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	l = len(m.Nonce)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	return n
}

func (m *MsgRevealTx) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
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
		if wireType == 2 {
			var byteLen int
			for shift := uint(0); ; shift += 7 {
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
			if byteLen < 0 || iNdEx+byteLen > l {
				return io.ErrUnexpectedEOF
			}
			switch fieldNum {
			case 1:
				m.Sender = string(dAtA[iNdEx : iNdEx+byteLen])
			case 2:
				m.CommitHash = append(m.CommitHash[:0], dAtA[iNdEx:iNdEx+byteLen]...)
			case 3:
				m.TxBody = append(m.TxBody[:0], dAtA[iNdEx:iNdEx+byteLen]...)
			case 4:
				m.Nonce = append(m.Nonce[:0], dAtA[iNdEx:iNdEx+byteLen]...)
			}
			iNdEx += byteLen
		} else {
			if wireType < 0 || wireType > 5 {
				return fmt.Errorf("proto: MsgRevealTx: illegal tag %d (wire type %d)", fieldNum, wireType)
			}
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			iNdEx += skippy
		}
		_ = preIndex
	}
	return nil
}

// Marshal / Unmarshal for MsgRevealTxResponse
func (m *MsgRevealTxResponse) Marshal() (dAtA []byte, err error) {
	return []byte{}, nil
}
func (m *MsgRevealTxResponse) MarshalTo(dAtA []byte) (int, error)            { return 0, nil }
func (m *MsgRevealTxResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgRevealTxResponse) Size() int                                     { return 0 }
func (m *MsgRevealTxResponse) Unmarshal(dAtA []byte) error                   { return nil }

// Helper functions for protobuf encoding
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
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
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
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, fmt.Errorf("proto: negative length found during unmarshaling")
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
			return 0, fmt.Errorf("proto: negative position after skip")
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}
