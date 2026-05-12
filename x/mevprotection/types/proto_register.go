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
	typeBytes := descriptorpb.FieldDescriptorProto_TYPE_BYTES

	fd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/mevprotection/tx.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.mevprotection"),
		MessageType: []*descriptorpb.DescriptorProto{
			{ // 0: MsgCommitTx
				Name: sp("MsgCommitTx"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
					{Name: sp("tx_hash"), Number: int32p(2), Label: &label, Type: &typeBytes, JsonName: sp("txHash")},
					{Name: sp("encrypted_tx"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: sp("encryptedTx")},
				},
			},
			{ // 1: MsgCommitTxResponse
				Name: sp("MsgCommitTxResponse"),
			},
			{ // 2: MsgRevealTx
				Name: sp("MsgRevealTx"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
					{Name: sp("commit_hash"), Number: int32p(2), Label: &label, Type: &typeBytes, JsonName: sp("commitHash")},
					{Name: sp("tx_body"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: sp("txBody")},
					{Name: sp("nonce"), Number: int32p(4), Label: &label, Type: &typeBytes, JsonName: sp("nonce")},
				},
			},
			{ // 3: MsgRevealTxResponse
				Name: sp("MsgRevealTxResponse"),
			},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: sp("CommitTx"), InputType: sp(".syreen.mevprotection.MsgCommitTx"), OutputType: sp(".syreen.mevprotection.MsgCommitTxResponse")},
					{Name: sp("RevealTx"), InputType: sp(".syreen.mevprotection.MsgRevealTx"), OutputType: sp(".syreen.mevprotection.MsgRevealTxResponse")},
				},
			},
		},
	}

	// Serialize and gzip the file descriptor for Descriptor() methods
	rawDesc, err := proto.Marshal(fd)
	if err != nil {
		panic("mevprotection: failed to marshal file descriptor: " + err.Error())
	}
	var buf bytes.Buffer
	gz, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	gz.Write(rawDesc)
	gz.Close()
	fileDescriptorTx = buf.Bytes()

	file, err := protodesc.NewFile(fd, nil)
	if err != nil {
		panic("mevprotection: failed to create proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(file); err != nil {
		// Already registered is fine
	}

	qfd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/mevprotection/query.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.mevprotection"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: sp("QueryParamsRequest")},
			{Name: sp("QueryParamsResponse")},
			{Name: sp("QueryCommittedTxRequest")},
			{Name: sp("QueryCommittedTxResponse")},
			{Name: sp("QueryPendingRevealsRequest")},
			{Name: sp("QueryPendingRevealsResponse")},
			{Name: sp("QueryMEVRedistributionRequest")},
			{Name: sp("QueryMEVRedistributionResponse")},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Query"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: sp("Params"), InputType: sp(".syreen.mevprotection.QueryParamsRequest"), OutputType: sp(".syreen.mevprotection.QueryParamsResponse")},
					{Name: sp("CommittedTx"), InputType: sp(".syreen.mevprotection.QueryCommittedTxRequest"), OutputType: sp(".syreen.mevprotection.QueryCommittedTxResponse")},
					{Name: sp("PendingReveals"), InputType: sp(".syreen.mevprotection.QueryPendingRevealsRequest"), OutputType: sp(".syreen.mevprotection.QueryPendingRevealsResponse")},
					{Name: sp("MEVRedistribution"), InputType: sp(".syreen.mevprotection.QueryMEVRedistributionRequest"), OutputType: sp(".syreen.mevprotection.QueryMEVRedistributionResponse")},
				},
			},
		},
	}
	qfile, err := protodesc.NewFile(qfd, nil)
	if err != nil {
		panic("mevprotection: failed to create query proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(qfile); err != nil {
		// Already registered is fine
	}
}

func sp(s string) *string    { return &s }
func int32p(i int32) *int32 { return &i }
