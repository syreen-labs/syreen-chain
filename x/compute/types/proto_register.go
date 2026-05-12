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
	strType := descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum()
	bytesType := descriptorpb.FieldDescriptorProto_TYPE_BYTES.Enum()
	uint64Type := descriptorpb.FieldDescriptorProto_TYPE_UINT64.Enum()
	optLabel := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum()
	repLabel := descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum()

	fd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/compute/tx.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.compute"),
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: sp("MsgStoreCode"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("sender"), Number: ip(1), Type: strType, Label: optLabel},
					{Name: sp("wasm_byte_code"), Number: ip(2), Type: bytesType, Label: optLabel},
					{Name: sp("instantiate_permission"), Number: ip(3), Type: bytesType, Label: optLabel},
				},
			},
			{
				Name: sp("MsgStoreCodeResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("code_id"), Number: ip(1), Type: uint64Type, Label: optLabel},
				},
			},
			{
				Name: sp("MsgInstantiateContract"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("sender"), Number: ip(1), Type: strType, Label: optLabel},
					{Name: sp("admin"), Number: ip(2), Type: strType, Label: optLabel},
					{Name: sp("code_id"), Number: ip(3), Type: uint64Type, Label: optLabel},
					{Name: sp("label"), Number: ip(4), Type: strType, Label: optLabel},
					{Name: sp("msg"), Number: ip(5), Type: bytesType, Label: optLabel},
					{Name: sp("funds"), Number: ip(6), Type: bytesType, Label: repLabel},
				},
			},
			{
				Name: sp("MsgInstantiateContractResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("address"), Number: ip(1), Type: strType, Label: optLabel},
				},
			},
			{
				Name: sp("MsgExecuteContract"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("sender"), Number: ip(1), Type: strType, Label: optLabel},
					{Name: sp("contract"), Number: ip(2), Type: strType, Label: optLabel},
					{Name: sp("msg"), Number: ip(3), Type: bytesType, Label: optLabel},
					{Name: sp("funds"), Number: ip(4), Type: bytesType, Label: repLabel},
				},
			},
			{
				Name: sp("MsgExecuteContractResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("data"), Number: ip(1), Type: bytesType, Label: optLabel},
				},
			},
			{
				Name: sp("MsgMigrateContract"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("sender"), Number: ip(1), Type: strType, Label: optLabel},
					{Name: sp("contract"), Number: ip(2), Type: strType, Label: optLabel},
					{Name: sp("code_id"), Number: ip(3), Type: uint64Type, Label: optLabel},
					{Name: sp("msg"), Number: ip(4), Type: bytesType, Label: optLabel},
				},
			},
			{
				Name: sp("MsgMigrateContractResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("data"), Number: ip(1), Type: bytesType, Label: optLabel},
				},
			},
			{
				Name: sp("MsgUpdateAdmin"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("sender"), Number: ip(1), Type: strType, Label: optLabel},
					{Name: sp("new_admin"), Number: ip(2), Type: strType, Label: optLabel},
					{Name: sp("contract"), Number: ip(3), Type: strType, Label: optLabel},
				},
			},
			{
				Name: sp("MsgUpdateAdminResponse"),
			},
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

	// Serialize and gzip the file descriptor for Descriptor() methods
	rawDesc, err := proto.Marshal(fd)
	if err != nil {
		panic("compute: failed to marshal file descriptor: " + err.Error())
	}
	var buf bytes.Buffer
	gz, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	gz.Write(rawDesc)
	gz.Close()
	fileDescriptorTx = buf.Bytes()

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
func ip(i int32) *int32   { return &i }
