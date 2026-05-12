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
		Name:    sp("syreen/abstractaccount/tx.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.abstractaccount"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: sp("MsgCreateSmartAccount")},
			{Name: sp("MsgCreateSmartAccountResponse")},
			{Name: sp("MsgCreateSessionKey")},
			{Name: sp("MsgCreateSessionKeyResponse")},
			{Name: sp("MsgRevokeSessionKey")},
			{Name: sp("MsgRevokeSessionKeyResponse")},
			{Name: sp("MsgInitiateRecovery")},
			{Name: sp("MsgInitiateRecoveryResponse")},
			{Name: sp("MsgApproveRecovery")},
			{Name: sp("MsgApproveRecoveryResponse")},
			{Name: sp("MsgExecuteRecovery")},
			{Name: sp("MsgExecuteRecoveryResponse")},
			{Name: sp("MsgSponsorGas")},
			{Name: sp("MsgSponsorGasResponse")},
			{Name: sp("MsgBatchExecute")},
			{Name: sp("MsgBatchExecuteResponse")},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: sp("CreateSmartAccount"), InputType: sp(".syreen.abstractaccount.MsgCreateSmartAccount"), OutputType: sp(".syreen.abstractaccount.MsgCreateSmartAccountResponse")},
					{Name: sp("CreateSessionKey"), InputType: sp(".syreen.abstractaccount.MsgCreateSessionKey"), OutputType: sp(".syreen.abstractaccount.MsgCreateSessionKeyResponse")},
					{Name: sp("RevokeSessionKey"), InputType: sp(".syreen.abstractaccount.MsgRevokeSessionKey"), OutputType: sp(".syreen.abstractaccount.MsgRevokeSessionKeyResponse")},
					{Name: sp("InitiateRecovery"), InputType: sp(".syreen.abstractaccount.MsgInitiateRecovery"), OutputType: sp(".syreen.abstractaccount.MsgInitiateRecoveryResponse")},
					{Name: sp("ApproveRecovery"), InputType: sp(".syreen.abstractaccount.MsgApproveRecovery"), OutputType: sp(".syreen.abstractaccount.MsgApproveRecoveryResponse")},
					{Name: sp("ExecuteRecovery"), InputType: sp(".syreen.abstractaccount.MsgExecuteRecovery"), OutputType: sp(".syreen.abstractaccount.MsgExecuteRecoveryResponse")},
					{Name: sp("SponsorGas"), InputType: sp(".syreen.abstractaccount.MsgSponsorGas"), OutputType: sp(".syreen.abstractaccount.MsgSponsorGasResponse")},
					{Name: sp("BatchExecute"), InputType: sp(".syreen.abstractaccount.MsgBatchExecute"), OutputType: sp(".syreen.abstractaccount.MsgBatchExecuteResponse")},
				},
			},
		},
	}

	file, err := protodesc.NewFile(fd, nil)
	if err != nil {
		panic("abstractaccount: failed to create proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(file); err != nil {
		// Already registered is fine
	}

	qfd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/abstractaccount/query.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.abstractaccount"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: sp("QueryParamsRequest")},
			{Name: sp("QueryParamsResponse")},
			{Name: sp("QuerySmartAccountRequest")},
			{Name: sp("QuerySmartAccountResponse")},
			{Name: sp("QuerySessionKeysRequest")},
			{Name: sp("QuerySessionKeysResponse")},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Query"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: sp("Params"), InputType: sp(".syreen.abstractaccount.QueryParamsRequest"), OutputType: sp(".syreen.abstractaccount.QueryParamsResponse")},
					{Name: sp("SmartAccount"), InputType: sp(".syreen.abstractaccount.QuerySmartAccountRequest"), OutputType: sp(".syreen.abstractaccount.QuerySmartAccountResponse")},
					{Name: sp("SessionKeys"), InputType: sp(".syreen.abstractaccount.QuerySessionKeysRequest"), OutputType: sp(".syreen.abstractaccount.QuerySessionKeysResponse")},
				},
			},
		},
	}
	qfile, err := protodesc.NewFile(qfd, nil)
	if err != nil {
		panic("abstractaccount: failed to create query proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(qfile); err != nil {
		// Already registered is fine
	}
}

func sp(s string) *string { return &s }
