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
	labelRep := descriptorpb.FieldDescriptorProto_LABEL_REPEATED
	typeString := descriptorpb.FieldDescriptorProto_TYPE_STRING
	typeBytes := descriptorpb.FieldDescriptorProto_TYPE_BYTES
	typeUint32 := descriptorpb.FieldDescriptorProto_TYPE_UINT32
	typeUint64 := descriptorpb.FieldDescriptorProto_TYPE_UINT64
	typeInt64 := descriptorpb.FieldDescriptorProto_TYPE_INT64

	fd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/abstractaccount/tx.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.abstractaccount"),
		MessageType: []*descriptorpb.DescriptorProto{
			{ // 0: MsgCreateSmartAccount
				Name: sp("MsgCreateSmartAccount"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
					{Name: sp("account_type"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: sp("accountType")},
					{Name: sp("owners"), Number: int32p(3), Label: &labelRep, Type: &typeString, JsonName: sp("owners")},
					{Name: sp("threshold"), Number: int32p(4), Label: &label, Type: &typeUint32, JsonName: sp("threshold")},
				},
			},
			{ // 1: MsgCreateSmartAccountResponse
				Name: sp("MsgCreateSmartAccountResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("address"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: sp("address")},
				},
			},
			{ // 2: MsgCreateSessionKey
				Name: sp("MsgCreateSessionKey"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("granter"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: sp("granter")},
					{Name: sp("grantee"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: sp("grantee")},
					{Name: sp("permissions"), Number: int32p(3), Label: &labelRep, Type: &typeBytes, JsonName: sp("permissions")},
					{Name: sp("duration"), Number: int32p(4), Label: &label, Type: &typeInt64, JsonName: sp("duration")},
				},
			},
			{ // 3: MsgCreateSessionKeyResponse
				Name: sp("MsgCreateSessionKeyResponse"),
			},
			{ // 4: MsgRevokeSessionKey
				Name: sp("MsgRevokeSessionKey"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("granter"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: sp("granter")},
					{Name: sp("session_key_addr"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: sp("sessionKeyAddr")},
				},
			},
			{ // 5: MsgRevokeSessionKeyResponse
				Name: sp("MsgRevokeSessionKeyResponse"),
			},
			{ // 6: MsgInitiateRecovery
				Name: sp("MsgInitiateRecovery"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("guardian"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: sp("guardian")},
					{Name: sp("account"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: sp("account")},
					{Name: sp("new_owners"), Number: int32p(3), Label: &labelRep, Type: &typeString, JsonName: sp("newOwners")},
				},
			},
			{ // 7: MsgInitiateRecoveryResponse
				Name: sp("MsgInitiateRecoveryResponse"),
			},
			{ // 8: MsgApproveRecovery
				Name: sp("MsgApproveRecovery"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("guardian"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: sp("guardian")},
					{Name: sp("account"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: sp("account")},
				},
			},
			{ // 9: MsgApproveRecoveryResponse
				Name: sp("MsgApproveRecoveryResponse"),
			},
			{ // 10: MsgExecuteRecovery
				Name: sp("MsgExecuteRecovery"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
					{Name: sp("account"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: sp("account")},
				},
			},
			{ // 11: MsgExecuteRecoveryResponse
				Name: sp("MsgExecuteRecoveryResponse"),
			},
			{ // 12: MsgSponsorGas
				Name: sp("MsgSponsorGas"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("sponsor"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: sp("sponsor")},
					{Name: sp("sponsored"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: sp("sponsored")},
					{Name: sp("gas_limit"), Number: int32p(3), Label: &label, Type: &typeUint64, JsonName: sp("gasLimit")},
					{Name: sp("duration"), Number: int32p(4), Label: &label, Type: &typeInt64, JsonName: sp("duration")},
				},
			},
			{ // 13: MsgSponsorGasResponse
				Name: sp("MsgSponsorGasResponse"),
			},
			{ // 14: MsgBatchExecute
				Name: sp("MsgBatchExecute"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
					{Name: sp("messages"), Number: int32p(2), Label: &labelRep, Type: &typeBytes, JsonName: sp("messages")},
				},
			},
			{ // 15: MsgBatchExecuteResponse
				Name: sp("MsgBatchExecuteResponse"),
			},
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

	// Serialize and gzip the file descriptor for Descriptor() methods
	rawDesc, err := proto.Marshal(fd)
	if err != nil {
		panic("abstractaccount: failed to marshal file descriptor: " + err.Error())
	}
	var buf bytes.Buffer
	gz, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	gz.Write(rawDesc)
	gz.Close()
	fileDescriptorTx = buf.Bytes()

	// Register with nil resolver (no dependencies needed)
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

func sp(s string) *string  { return &s }
func int32p(i int32) *int32 { return &i }
