package types

import (
	"bytes"
	"compress/gzip"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

func init() { registerOptionsProtoFileDescriptors() }

func registerOptionsProtoFileDescriptors() {
	label := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	typeString := descriptorpb.FieldDescriptorProto_TYPE_STRING
	typeBytes := descriptorpb.FieldDescriptorProto_TYPE_BYTES
	typeUint64 := descriptorpb.FieldDescriptorProto_TYPE_UINT64
	typeInt64 := descriptorpb.FieldDescriptorProto_TYPE_INT64

	fd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/options/tx.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.options"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: sp("MsgWriteOption"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("writer"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("writer")},
				{Name: sp("pool_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("poolId")},
				{Name: sp("option_type"), Number: ip(3), Label: &label, Type: &typeString, JsonName: sp("optionType")},
				{Name: sp("strike_price"), Number: ip(4), Label: &label, Type: &typeBytes, JsonName: sp("strikePrice")},
				{Name: sp("amount"), Number: ip(5), Label: &label, Type: &typeBytes, JsonName: sp("amount")},
				{Name: sp("expiry_block"), Number: ip(6), Label: &label, Type: &typeInt64, JsonName: sp("expiryBlock")},
				{Name: sp("custom_premium"), Number: ip(7), Label: &label, Type: &typeBytes, JsonName: sp("customPremium")},
			}},
			{Name: sp("MsgWriteOptionResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("option_id"), Number: ip(1), Label: &label, Type: &typeUint64, JsonName: sp("optionId")},
				{Name: sp("premium"), Number: ip(2), Label: &label, Type: &typeBytes, JsonName: sp("premium")},
			}},
			{Name: sp("MsgBuyOption"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("buyer"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("buyer")},
				{Name: sp("option_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("optionId")},
			}},
			{Name: sp("MsgBuyOptionResponse")},
			{Name: sp("MsgExerciseOption"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("buyer"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("buyer")},
				{Name: sp("option_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("optionId")},
			}},
			{Name: sp("MsgExerciseOptionResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("payout"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("payout")},
			}},
			{Name: sp("MsgCancelOption"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("writer"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("writer")},
				{Name: sp("option_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("optionId")},
			}},
			{Name: sp("MsgCancelOptionResponse")},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: sp("WriteOption"), InputType: sp(".syreen.options.MsgWriteOption"), OutputType: sp(".syreen.options.MsgWriteOptionResponse")},
					{Name: sp("BuyOption"), InputType: sp(".syreen.options.MsgBuyOption"), OutputType: sp(".syreen.options.MsgBuyOptionResponse")},
					{Name: sp("ExerciseOption"), InputType: sp(".syreen.options.MsgExerciseOption"), OutputType: sp(".syreen.options.MsgExerciseOptionResponse")},
					{Name: sp("CancelOption"), InputType: sp(".syreen.options.MsgCancelOption"), OutputType: sp(".syreen.options.MsgCancelOptionResponse")},
				},
			},
		},
	}

	rawBz, err := proto.Marshal(fd)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	_, _ = w.Write(rawBz)
	_ = w.Close()
	fileDescriptorOptionsTx = buf.Bytes()

	file, err := protodesc.NewFile(fd, nil)
	if err != nil {
		panic("options: failed to create proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(file); err != nil {
		// already registered is fine
	}
}

func sp(s string) *string { return &s }
func ip(n int32) *int32   { return &n }
