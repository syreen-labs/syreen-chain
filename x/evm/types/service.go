package types

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"

	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// fileDescriptor is a compressed proto file descriptor for the EVM tx types.
// Used by Descriptor() methods for gogoproto compatibility.
var fileDescriptor_evm_tx []byte

func init() {
	optionalLabel := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	stringType := descriptorpb.FieldDescriptorProto_TYPE_STRING
	uint64Type := descriptorpb.FieldDescriptorProto_TYPE_UINT64

	fd := &descriptorpb.FileDescriptorProto{
		Name:    strp("syreen/evm/tx.proto"),
		Syntax:  strp("proto3"),
		Package: strp("syreen.evm"),
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: strp("MsgEthereumTx"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("from"), Number: intp(1), Type: &stringType, Label: &optionalLabel},
					{Name: strp("to"), Number: intp(2), Type: &stringType, Label: &optionalLabel},
					{Name: strp("value"), Number: intp(3), Type: &stringType, Label: &optionalLabel},
					{Name: strp("gas_limit"), Number: intp(4), Type: &uint64Type, Label: &optionalLabel},
					{Name: strp("data"), Number: intp(5), Type: &stringType, Label: &optionalLabel},
					{Name: strp("nonce"), Number: intp(6), Type: &uint64Type, Label: &optionalLabel},
					{Name: strp("raw_tx_hex"), Number: intp(7), Type: &stringType, Label: &optionalLabel},
				},
			},
			{
				Name: strp("MsgEthereumTxResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("gas_used"), Number: intp(1), Type: &uint64Type, Label: &optionalLabel},
					{Name: strp("vm_error"), Number: intp(2), Type: &stringType, Label: &optionalLabel},
					{Name: strp("return_data"), Number: intp(3), Type: &stringType, Label: &optionalLabel},
					{Name: strp("contract_address"), Number: intp(4), Type: &stringType, Label: &optionalLabel},
				},
			},
		},
	}
	raw, _ := proto.Marshal(fd)
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	w.Write(raw)
	w.Close()
	fileDescriptor_evm_tx = buf.Bytes()
}

func strp(s string) *string { return &s }
func intp(n int32) *int32   { return &n }

// ---------------------------------------------------------------------------
// Proto wire format encoding helpers
// ---------------------------------------------------------------------------

func appendVarint(buf []byte, v uint64) []byte {
	for v >= 0x80 {
		buf = append(buf, byte(v)|0x80)
		v >>= 7
	}
	return append(buf, byte(v))
}

func appendString(buf []byte, fieldNum int, s string) []byte {
	if s == "" {
		return buf
	}
	tag := uint64(fieldNum<<3 | 2) // wire type 2 = length-delimited
	buf = appendVarint(buf, tag)
	buf = appendVarint(buf, uint64(len(s)))
	buf = append(buf, s...)
	return buf
}

func appendUint64(buf []byte, fieldNum int, v uint64) []byte {
	if v == 0 {
		return buf
	}
	tag := uint64(fieldNum<<3 | 0) // wire type 0 = varint
	buf = appendVarint(buf, tag)
	buf = appendVarint(buf, v)
	return buf
}

func readVarint(data []byte, offset int) (uint64, int, error) {
	val := uint64(0)
	shift := uint(0)
	for i := offset; i < len(data); i++ {
		b := data[i]
		val |= uint64(b&0x7F) << shift
		if b < 0x80 {
			return val, i + 1, nil
		}
		shift += 7
		if shift >= 64 {
			return 0, 0, fmt.Errorf("varint overflow")
		}
	}
	return 0, 0, io.ErrUnexpectedEOF
}

// ---------------------------------------------------------------------------
// Proto-compatible message methods for MsgEthereumTx
// ---------------------------------------------------------------------------

func (m *MsgEthereumTx) ProtoMessage()           {}
func (m *MsgEthereumTx) Reset()                  { *m = MsgEthereumTx{} }
func (m *MsgEthereumTx) String() string          { return fmt.Sprintf("from=%s to=%s", m.From, m.To) }
func (m *MsgEthereumTx) XXX_MessageName() string { return "syreen.evm.MsgEthereumTx" }

// Descriptor returns the gzipped file descriptor bytes and the path to this message.
func (*MsgEthereumTx) Descriptor() ([]byte, []int) {
	return fileDescriptor_evm_tx, []int{0}
}

// Marshal encodes MsgEthereumTx in proto wire format.
func (m *MsgEthereumTx) Marshal() ([]byte, error) {
	buf := make([]byte, 0, m.Size())
	buf = appendString(buf, 1, m.From)
	buf = appendString(buf, 2, m.To)
	buf = appendString(buf, 3, m.Value)
	buf = appendUint64(buf, 4, m.GasLimit)
	buf = appendString(buf, 5, m.Data)
	buf = appendUint64(buf, 6, m.Nonce)
	buf = appendString(buf, 7, m.RawTxHex)
	return buf, nil
}

func (m *MsgEthereumTx) MarshalTo(data []byte) (int, error) {
	bz, err := m.Marshal()
	if err != nil {
		return 0, err
	}
	copy(data, bz)
	return len(bz), nil
}

// Unmarshal decodes MsgEthereumTx from proto wire format.
func (m *MsgEthereumTx) Unmarshal(data []byte) error {
	offset := 0
	for offset < len(data) {
		tagVal, newOff, err := readVarint(data, offset)
		if err != nil {
			return fmt.Errorf("read tag: %w", err)
		}
		offset = newOff
		fieldNum := int(tagVal >> 3)
		wireType := tagVal & 0x7

		switch wireType {
		case 0: // varint
			val, newOff, err := readVarint(data, offset)
			if err != nil {
				return fmt.Errorf("read varint field %d: %w", fieldNum, err)
			}
			offset = newOff
			switch fieldNum {
			case 4:
				m.GasLimit = val
			case 6:
				m.Nonce = val
			}
		case 2: // length-delimited
			length, newOff, err := readVarint(data, offset)
			if err != nil {
				return fmt.Errorf("read length field %d: %w", fieldNum, err)
			}
			offset = newOff
			if offset+int(length) > len(data) {
				return io.ErrUnexpectedEOF
			}
			s := string(data[offset : offset+int(length)])
			offset += int(length)
			switch fieldNum {
			case 1:
				m.From = s
			case 2:
				m.To = s
			case 3:
				m.Value = s
			case 5:
				m.Data = s
			case 7:
				m.RawTxHex = s
			}
		default:
			return fmt.Errorf("unsupported wire type %d for field %d", wireType, fieldNum)
		}
	}
	return nil
}

// Size returns the proto wire format size.
func (m *MsgEthereumTx) Size() int {
	n := 0
	if m.From != "" {
		n += 1 + sovSize(uint64(len(m.From))) + len(m.From)
	}
	if m.To != "" {
		n += 1 + sovSize(uint64(len(m.To))) + len(m.To)
	}
	if m.Value != "" {
		n += 1 + sovSize(uint64(len(m.Value))) + len(m.Value)
	}
	if m.GasLimit != 0 {
		n += 1 + sovSize(m.GasLimit)
	}
	if m.Data != "" {
		n += 1 + sovSize(uint64(len(m.Data))) + len(m.Data)
	}
	if m.Nonce != 0 {
		n += 1 + sovSize(m.Nonce)
	}
	if m.RawTxHex != "" {
		n += 1 + sovSize(uint64(len(m.RawTxHex))) + len(m.RawTxHex)
	}
	return n
}

func sovSize(x uint64) int {
	n := 0
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}

// Ensure binary is imported
var _ = binary.MaxVarintLen64

// ---------------------------------------------------------------------------
// Proto-compatible message methods for MsgEthereumTxResponse
// ---------------------------------------------------------------------------

func (m *MsgEthereumTxResponse) ProtoMessage()           {}
func (m *MsgEthereumTxResponse) Reset()                  { *m = MsgEthereumTxResponse{} }
func (m *MsgEthereumTxResponse) String() string          { return fmt.Sprintf("gas_used=%d", m.GasUsed) }
func (m *MsgEthereumTxResponse) XXX_MessageName() string { return "syreen.evm.MsgEthereumTxResponse" }

// Descriptor returns the gzipped file descriptor bytes and the path to this message.
func (*MsgEthereumTxResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_evm_tx, []int{1}
}

// ---------------------------------------------------------------------------
// MsgServer Interface (H3)
// ---------------------------------------------------------------------------

type MsgServer interface {
	EthereumTx(context.Context, *MsgEthereumTx) (*MsgEthereumTxResponse, error)
}

// ---------------------------------------------------------------------------
// gRPC Msg Handler
// ---------------------------------------------------------------------------

func _Msg_EthereumTx_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgEthereumTx)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).EthereumTx(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syreen.evm.Msg/EthereumTx",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).EthereumTx(ctx, req.(*MsgEthereumTx))
	}
	return interceptor(ctx, in, info, handler)
}

// ---------------------------------------------------------------------------
// gRPC Service Descriptor
// ---------------------------------------------------------------------------

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.evm.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "EthereumTx",
			Handler:    _Msg_EthereumTx_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/evm/tx.proto",
}

// ---------------------------------------------------------------------------
// Registration Functions (H3)
// ---------------------------------------------------------------------------

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}
