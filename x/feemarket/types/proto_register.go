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
	qfd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/feemarket/query.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.feemarket"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: sp("QueryParamsRequest")},
			{Name: sp("QueryParamsResponse")},
			{Name: sp("QueryBaseFeeRequest")},
			{Name: sp("QueryBaseFeeResponse")},
			{Name: sp("QueryFeeStateRequest")},
			{Name: sp("QueryFeeStateResponse")},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Query"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: sp("Params"), InputType: sp(".syreen.feemarket.QueryParamsRequest"), OutputType: sp(".syreen.feemarket.QueryParamsResponse")},
					{Name: sp("BaseFee"), InputType: sp(".syreen.feemarket.QueryBaseFeeRequest"), OutputType: sp(".syreen.feemarket.QueryBaseFeeResponse")},
					{Name: sp("FeeState"), InputType: sp(".syreen.feemarket.QueryFeeStateRequest"), OutputType: sp(".syreen.feemarket.QueryFeeStateResponse")},
				},
			},
		},
	}
	qfile, err := protodesc.NewFile(qfd, nil)
	if err != nil {
		panic("feemarket: failed to create query proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(qfile); err != nil {
		// Already registered is fine
	}
}

func sp(s string) *string { return &s }
