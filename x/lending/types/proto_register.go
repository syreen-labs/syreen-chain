package types

import (
	"bytes"
	"compress/gzip"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

var fileDescriptorLendingTx []byte

func init() {
	registerLendingProtoFileDescriptors()
}

func registerLendingProtoFileDescriptors() {
	label := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	typeString := descriptorpb.FieldDescriptorProto_TYPE_STRING
	typeBytes := descriptorpb.FieldDescriptorProto_TYPE_BYTES
	typeUint64 := descriptorpb.FieldDescriptorProto_TYPE_UINT64

	fd := &descriptorpb.FileDescriptorProto{
		Name:    strp("syreen/lending/tx.proto"),
		Syntax:  strp("proto3"),
		Package: strp("syreen.lending"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: strp("MsgDeposit"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
				{Name: strp("pool_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
				{Name: strp("amount"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: strp("amount")},
			}},
			{Name: strp("MsgDepositResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("interest_earned"), Number: int32p(1), Label: &label, Type: &typeBytes, JsonName: strp("interestEarned")},
			}},
			{Name: strp("MsgWithdraw"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
				{Name: strp("pool_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
				{Name: strp("amount"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: strp("amount")},
			}},
			{Name: strp("MsgWithdrawResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("interest_earned"), Number: int32p(1), Label: &label, Type: &typeBytes, JsonName: strp("interestEarned")},
			}},
			{Name: strp("MsgBorrow"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
				{Name: strp("borrow_pool_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("borrowPoolId")},
				{Name: strp("amount"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: strp("amount")},
				{Name: strp("collateral_pool_id"), Number: int32p(4), Label: &label, Type: &typeUint64, JsonName: strp("collateralPoolId")},
				{Name: strp("collateral_amount"), Number: int32p(5), Label: &label, Type: &typeBytes, JsonName: strp("collateralAmount")},
			}},
			{Name: strp("MsgBorrowResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("borrow_id"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("borrowId")},
			}},
			{Name: strp("MsgRepay"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
				{Name: strp("borrow_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("borrowId")},
				{Name: strp("amount"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: strp("amount")},
			}},
			{Name: strp("MsgRepayResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("interest_paid"), Number: int32p(1), Label: &label, Type: &typeBytes, JsonName: strp("interestPaid")},
				{Name: strp("collateral_returned"), Number: int32p(2), Label: &label, Type: &typeBytes, JsonName: strp("collateralReturned")},
			}},
			{Name: strp("MsgLiquidate"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("liquidator"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("liquidator")},
				{Name: strp("borrow_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("borrowId")},
			}},
			{Name: strp("MsgLiquidateResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("collateral_seized"), Number: int32p(1), Label: &label, Type: &typeBytes, JsonName: strp("collateralSeized")},
				{Name: strp("debt_repaid"), Number: int32p(2), Label: &label, Type: &typeBytes, JsonName: strp("debtRepaid")},
			}},
			{Name: strp("MsgCreateLendingPool"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("authority"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("authority")},
				{Name: strp("denom"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: strp("denom")},
				{Name: strp("collateral_factor"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: strp("collateralFactor")},
				{Name: strp("dex_pool_id"), Number: int32p(4), Label: &label, Type: &typeUint64, JsonName: strp("dexPoolId")},
				{Name: strp("price_denom"), Number: int32p(5), Label: &label, Type: &typeString, JsonName: strp("priceDenom")},
			}},
			{Name: strp("MsgCreateLendingPoolResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("pool_id"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
			}},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: strp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: strp("Deposit"), InputType: strp(".syreen.lending.MsgDeposit"), OutputType: strp(".syreen.lending.MsgDepositResponse")},
					{Name: strp("Withdraw"), InputType: strp(".syreen.lending.MsgWithdraw"), OutputType: strp(".syreen.lending.MsgWithdrawResponse")},
					{Name: strp("Borrow"), InputType: strp(".syreen.lending.MsgBorrow"), OutputType: strp(".syreen.lending.MsgBorrowResponse")},
					{Name: strp("Repay"), InputType: strp(".syreen.lending.MsgRepay"), OutputType: strp(".syreen.lending.MsgRepayResponse")},
					{Name: strp("Liquidate"), InputType: strp(".syreen.lending.MsgLiquidate"), OutputType: strp(".syreen.lending.MsgLiquidateResponse")},
					{Name: strp("CreateLendingPool"), InputType: strp(".syreen.lending.MsgCreateLendingPool"), OutputType: strp(".syreen.lending.MsgCreateLendingPoolResponse")},
				},
			},
		},
	}

	rawBz, err := proto.Marshal(fd)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	w, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	_, _ = w.Write(rawBz)
	_ = w.Close()
	fileDescriptorLendingTx = buf.Bytes()

	file, err := protodesc.NewFile(fd, nil)
	if err != nil {
		panic("lending: failed to create proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(file); err != nil {
		// Already registered is fine
	}
}

func strp(s string) *string  { return &s }
func int32p(i int32) *int32 { return &i }
