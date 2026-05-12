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
	fd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/compute/tx.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.compute"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: sp("MsgStoreCode")},
			{Name: sp("MsgStoreCodeResponse")},
			{Name: sp("MsgInstantiateContract")},
			{Name: sp("MsgInstantiateContractResponse")},
			{Name: sp("MsgExecuteContract")},
			{Name: sp("MsgExecuteContractResponse")},
			{Name: sp("MsgMigrateContract")},
			{Name: sp("MsgMigrateContractResponse")},
			{Name: sp("MsgUpdateAdmin")},
			{Name: sp("MsgUpdateAdminResponse")},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: sp("StoreCode"), InputType: sp(".syreen.compute.MsgStoreCode"), OutputType: sp(".syreen.compute.MsgStoreCodeResponse")},
					{Name: sp("InstantiateContract"), InputType: sp(".syreen.compute.MsgInstantiateContract"), OutputType: sp(".syreen.compute.MsgInstantiateContractResponse")},
					{Name: sp("ExecuteContract"), InputType: sp(".syreen.compute.MsgExecuteContract"), OutputType: sp(".syreen.compute.MsgExecuteContractResponse")},
					{Name: sp("MigrateContract"), InputType: sp(".syreen.compute.MsgMigrateContract"), OutputType: sp(".syreen.compute.MsgMigrateContractResponse")},
					{Name: sp("UpdateAdmin"), InputType: sp(".syreen.compute.MsgUpdateAdmin"), OutputType: sp(".syreen.compute.MsgUpdateAdminResponse")},
				},
			},
		},
	}

	file, err := protodesc.NewFile(fd, nil)
	if err != nil {
		panic("compute: failed to create proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(file); err != nil {
		// Already registered is fine
	}

	qfd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/compute/query.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.compute"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: sp("QueryParamsRequest")},
			{Name: sp("QueryParamsResponse")},
			{Name: sp("QueryCodeInfoRequest")},
			{Name: sp("QueryCodeInfoResponse")},
			{Name: sp("QueryContractInfoRequest")},
			{Name: sp("QueryContractInfoResponse")},
			{Name: sp("QuerySmartContractStateRequest")},
			{Name: sp("QuerySmartContractStateResponse")},
			{Name: sp("QueryRawContractStateRequest")},
			{Name: sp("QueryRawContractStateResponse")},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Query"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: sp("Params"), InputType: sp(".syreen.compute.QueryParamsRequest"), OutputType: sp(".syreen.compute.QueryParamsResponse")},
					{Name: sp("CodeInfo"), InputType: sp(".syreen.compute.QueryCodeInfoRequest"), OutputType: sp(".syreen.compute.QueryCodeInfoResponse")},
					{Name: sp("ContractInfo"), InputType: sp(".syreen.compute.QueryContractInfoRequest"), OutputType: sp(".syreen.compute.QueryContractInfoResponse")},
					{Name: sp("SmartContractState"), InputType: sp(".syreen.compute.QuerySmartContractStateRequest"), OutputType: sp(".syreen.compute.QuerySmartContractStateResponse")},
					{Name: sp("RawContractState"), InputType: sp(".syreen.compute.QueryRawContractStateRequest"), OutputType: sp(".syreen.compute.QueryRawContractStateResponse")},
				},
			},
		},
	}
	qfile, err := protodesc.NewFile(qfd, nil)
	if err != nil {
		panic("compute: failed to create query proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(qfile); err != nil {
		// Already registered is fine
	}
}

func sp(s string) *string { return &s }
