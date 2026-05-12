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
		Name:    strp("syreen/tokenfactory/tx.proto"),
		Syntax:  strp("proto3"),
		Package: strp("syreen.tokenfactory"),
		MessageType: []*descriptorpb.DescriptorProto{
			{ // 0: MsgCreateDenom
				Name: strp("MsgCreateDenom"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
					{Name: strp("subdenom"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: strp("subdenom")},
				},
			},
			{ // 1: MsgCreateDenomResponse
				Name: strp("MsgCreateDenomResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("new_token_denom"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("newTokenDenom")},
				},
			},
			{ // 2: MsgMint
				Name: strp("MsgMint"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
					{Name: strp("amount"), Number: int32p(2), Label: &label, Type: &typeBytes, JsonName: strp("amount")},
					{Name: strp("mint_to_address"), Number: int32p(3), Label: &label, Type: &typeString, JsonName: strp("mintToAddress")},
				},
			},
			{ // 3: MsgMintResponse
				Name: strp("MsgMintResponse"),
			},
			{ // 4: MsgBurn
				Name: strp("MsgBurn"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
					{Name: strp("amount"), Number: int32p(2), Label: &label, Type: &typeBytes, JsonName: strp("amount")},
					{Name: strp("burn_from_address"), Number: int32p(3), Label: &label, Type: &typeString, JsonName: strp("burnFromAddress")},
				},
			},
			{ // 5: MsgBurnResponse
				Name: strp("MsgBurnResponse"),
			},
			{ // 6: MsgChangeAdmin
				Name: strp("MsgChangeAdmin"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
					{Name: strp("denom"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: strp("denom")},
					{Name: strp("new_admin"), Number: int32p(3), Label: &label, Type: &typeString, JsonName: strp("newAdmin")},
				},
			},
			{ // 7: MsgChangeAdminResponse
				Name: strp("MsgChangeAdminResponse"),
			},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: strp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: strp("CreateDenom"), InputType: strp(".syreen.tokenfactory.MsgCreateDenom"), OutputType: strp(".syreen.tokenfactory.MsgCreateDenomResponse")},
					{Name: strp("Mint"), InputType: strp(".syreen.tokenfactory.MsgMint"), OutputType: strp(".syreen.tokenfactory.MsgMintResponse")},
					{Name: strp("Burn"), InputType: strp(".syreen.tokenfactory.MsgBurn"), OutputType: strp(".syreen.tokenfactory.MsgBurnResponse")},
					{Name: strp("ChangeAdmin"), InputType: strp(".syreen.tokenfactory.MsgChangeAdmin"), OutputType: strp(".syreen.tokenfactory.MsgChangeAdminResponse")},
				},
			},
		},
	}

	// Serialize and gzip the file descriptor for Descriptor() methods
	rawDesc, err := proto.Marshal(fd)
	if err != nil {
		panic("tokenfactory: failed to marshal file descriptor: " + err.Error())
	}
	var buf bytes.Buffer
	gz, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	gz.Write(rawDesc)
	gz.Close()
	fileDescriptorTx = buf.Bytes()

	// Register with nil resolver (no dependencies needed, Coin encoded as bytes)
	file, err := protodesc.NewFile(fd, nil)
	if err != nil {
		panic("tokenfactory: failed to create proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(file); err != nil {
		// Already registered is fine
	}

	// Register query proto
	qfd := &descriptorpb.FileDescriptorProto{
		Name:    strp("syreen/tokenfactory/query.proto"),
		Syntax:  strp("proto3"),
		Package: strp("syreen.tokenfactory"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: strp("QueryParamsRequest")},
			{Name: strp("QueryParamsResponse")},
			{Name: strp("QueryDenomAuthorityRequest")},
			{Name: strp("QueryDenomAuthorityResponse")},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: strp("Query"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: strp("Params"), InputType: strp(".syreen.tokenfactory.QueryParamsRequest"), OutputType: strp(".syreen.tokenfactory.QueryParamsResponse")},
					{Name: strp("DenomAuthority"), InputType: strp(".syreen.tokenfactory.QueryDenomAuthorityRequest"), OutputType: strp(".syreen.tokenfactory.QueryDenomAuthorityResponse")},
				},
			},
		},
	}
	qfile, err := protodesc.NewFile(qfd, nil)
	if err != nil {
		panic("tokenfactory: failed to create query proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(qfile); err != nil {
		// Already registered is fine
	}
}

func strp(s string) *string  { return &s }
func int32p(i int32) *int32 { return &i }
