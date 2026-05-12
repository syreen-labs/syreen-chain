package types

import (
	"bytes"
	"compress/gzip"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

func init() {
	registerFlashloanProtoFileDescriptors()
}

func registerFlashloanProtoFileDescriptors() {
	label := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	typeString := descriptorpb.FieldDescriptorProto_TYPE_STRING
	typeBytes := descriptorpb.FieldDescriptorProto_TYPE_BYTES

	fd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/flashloan/tx.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.flashloan"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: sp("MsgFlashLoan"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("denom"), Number: ip(2), Label: &label, Type: &typeString, JsonName: sp("denom")},
				{Name: sp("amount"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("amount")},
			}},
			{Name: sp("MsgFlashLoanResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("fee"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("fee")},
			}},
			{Name: sp("MsgCreateFlashPool"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("authority"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("authority")},
				{Name: sp("denom"), Number: ip(2), Label: &label, Type: &typeString, JsonName: sp("denom")},
				{Name: sp("fee_rate"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("feeRate")},
			}},
			{Name: sp("MsgCreateFlashPoolResponse")},
			{Name: sp("MsgFundFlashPool"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("denom"), Number: ip(2), Label: &label, Type: &typeString, JsonName: sp("denom")},
				{Name: sp("amount"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("amount")},
			}},
			{Name: sp("MsgFundFlashPoolResponse")},
			{Name: sp("MsgWithdrawFlashPool"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("denom"), Number: ip(2), Label: &label, Type: &typeString, JsonName: sp("denom")},
				{Name: sp("shares"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("shares")},
			}},
			{Name: sp("MsgWithdrawFlashPoolResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("amount_returned"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("amountReturned")},
			}},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: sp("FlashLoan"), InputType: sp(".syreen.flashloan.MsgFlashLoan"), OutputType: sp(".syreen.flashloan.MsgFlashLoanResponse")},
					{Name: sp("CreateFlashPool"), InputType: sp(".syreen.flashloan.MsgCreateFlashPool"), OutputType: sp(".syreen.flashloan.MsgCreateFlashPoolResponse")},
					{Name: sp("FundFlashPool"), InputType: sp(".syreen.flashloan.MsgFundFlashPool"), OutputType: sp(".syreen.flashloan.MsgFundFlashPoolResponse")},
					{Name: sp("WithdrawFlashPool"), InputType: sp(".syreen.flashloan.MsgWithdrawFlashPool"), OutputType: sp(".syreen.flashloan.MsgWithdrawFlashPoolResponse")},
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
	fileDescriptorFlashloanTx = buf.Bytes()

	file, err := protodesc.NewFile(fd, nil)
	if err != nil {
		panic("flashloan: failed to create proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(file); err != nil {
		// Already registered is fine
	}
}

func sp(s string) *string { return &s }
func ip(n int32) *int32   { return &n }
