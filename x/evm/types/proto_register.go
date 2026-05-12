package types

import (
	"bytes"
	"compress/gzip"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

var fileDescriptorTx []byte

func init() {
	registerProtoFileDescriptors()
}

func registerProtoFileDescriptors() {
	label := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	typeString := descriptorpb.FieldDescriptorProto_TYPE_STRING
	typeUint64 := descriptorpb.FieldDescriptorProto_TYPE_UINT64

	fd := &descriptorpb.FileDescriptorProto{
		Name:    strp("syreen/evm/tx.proto"),
		Syntax:  strp("proto3"),
		Package: strp("syreen.evm"),
		MessageType: []*descriptorpb.DescriptorProto{
			{ // 0: MsgEthereumTx
				Name: strp("MsgEthereumTx"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("from"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("from")},
					{Name: strp("to"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: strp("to")},
					{Name: strp("value"), Number: int32p(3), Label: &label, Type: &typeString, JsonName: strp("value")},
					{Name: strp("gas_limit"), Number: int32p(4), Label: &label, Type: &typeUint64, JsonName: strp("gasLimit")},
					{Name: strp("data"), Number: int32p(5), Label: &label, Type: &typeString, JsonName: strp("data")},
					{Name: strp("nonce"), Number: int32p(6), Label: &label, Type: &typeUint64, JsonName: strp("nonce")},
					{Name: strp("raw_tx_hex"), Number: int32p(7), Label: &label, Type: &typeString, JsonName: strp("rawTxHex")},
				},
			},
			{ // 1: MsgEthereumTxResponse
				Name: strp("MsgEthereumTxResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("gas_used"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("gasUsed")},
					{Name: strp("vm_error"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: strp("vmError")},
					{Name: strp("return_data"), Number: int32p(3), Label: &label, Type: &typeString, JsonName: strp("returnData")},
					{Name: strp("contract_address"), Number: int32p(4), Label: &label, Type: &typeString, JsonName: strp("contractAddress")},
				},
			},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: strp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: strp("EthereumTx"), InputType: strp(".syreen.evm.MsgEthereumTx"), OutputType: strp(".syreen.evm.MsgEthereumTxResponse")},
				},
			},
		},
	}

	// Serialize and gzip the file descriptor for Descriptor() methods
	rawDesc, err := proto.Marshal(fd)
	if err != nil {
		panic("evm: failed to marshal file descriptor: " + err.Error())
	}
	var buf bytes.Buffer
	gz, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	gz.Write(rawDesc)
	gz.Close()
	fileDescriptorTx = buf.Bytes()

	// Register with nil resolver (no dependencies needed)
	file, err := protodesc.NewFile(fd, nil)
	if err != nil {
		panic("evm: failed to create proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(file); err != nil {
		// Already registered is fine
	}
}

func strp(s string) *string  { return &s }
func int32p(i int32) *int32 { return &i }
