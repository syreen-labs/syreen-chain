package types

import (
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

func init() {
	registerProtoFileDescriptors()
}

func registerProtoFileDescriptors() {
	optionalLabel := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	stringType := descriptorpb.FieldDescriptorProto_TYPE_STRING
	uint64Type := descriptorpb.FieldDescriptorProto_TYPE_UINT64

	fd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/evm/tx.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.evm"),
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: sp("MsgEthereumTx"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("from"), Number: i32p(1), Type: &stringType, Label: &optionalLabel},
					{Name: sp("to"), Number: i32p(2), Type: &stringType, Label: &optionalLabel},
					{Name: sp("value"), Number: i32p(3), Type: &stringType, Label: &optionalLabel},
					{Name: sp("gas_limit"), Number: i32p(4), Type: &uint64Type, Label: &optionalLabel},
					{Name: sp("data"), Number: i32p(5), Type: &stringType, Label: &optionalLabel},
					{Name: sp("nonce"), Number: i32p(6), Type: &uint64Type, Label: &optionalLabel},
					{Name: sp("raw_tx_hex"), Number: i32p(7), Type: &stringType, Label: &optionalLabel},
				},
			},
			{
				Name: sp("MsgEthereumTxResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("gas_used"), Number: i32p(1), Type: &uint64Type, Label: &optionalLabel},
					{Name: sp("vm_error"), Number: i32p(2), Type: &stringType, Label: &optionalLabel},
					{Name: sp("return_data"), Number: i32p(3), Type: &stringType, Label: &optionalLabel},
					{Name: sp("contract_address"), Number: i32p(4), Type: &stringType, Label: &optionalLabel},
				},
			},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{
						Name:       sp("EthereumTx"),
						InputType:  sp(".syreen.evm.MsgEthereumTx"),
						OutputType: sp(".syreen.evm.MsgEthereumTxResponse"),
					},
				},
			},
		},
	}

	file, err := protodesc.NewFile(fd, nil)
	if err != nil {
		panic("evm: failed to create proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(file); err != nil {
		// File may already be registered from a previous init
		return
	}
}

func sp(s string) *string  { return &s }
func i32p(n int32) *int32  { return &n }
