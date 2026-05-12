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
	registerLendingProtoFileDescriptors()
}

func registerLendingProtoFileDescriptors() {
	label := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	typeString := descriptorpb.FieldDescriptorProto_TYPE_STRING
	typeBytes := descriptorpb.FieldDescriptorProto_TYPE_BYTES
	typeUint64 := descriptorpb.FieldDescriptorProto_TYPE_UINT64

	fd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/lending/tx.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.lending"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: sp("MsgDeposit"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("pool_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("poolId")},
				{Name: sp("amount"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("amount")},
			}},
			{Name: sp("MsgDepositResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("interest_earned"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("interestEarned")},
			}},
			{Name: sp("MsgWithdraw"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("pool_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("poolId")},
				{Name: sp("amount"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("amount")},
			}},
			{Name: sp("MsgWithdrawResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("interest_earned"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("interestEarned")},
			}},
			{Name: sp("MsgBorrow"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("borrow_pool_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("borrowPoolId")},
				{Name: sp("amount"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("amount")},
				{Name: sp("collateral_pool_id"), Number: ip(4), Label: &label, Type: &typeUint64, JsonName: sp("collateralPoolId")},
				{Name: sp("collateral_amount"), Number: ip(5), Label: &label, Type: &typeBytes, JsonName: sp("collateralAmount")},
			}},
			{Name: sp("MsgBorrowResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("borrow_id"), Number: ip(1), Label: &label, Type: &typeUint64, JsonName: sp("borrowId")},
			}},
			{Name: sp("MsgRepay"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("borrow_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("borrowId")},
				{Name: sp("amount"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("amount")},
			}},
			{Name: sp("MsgRepayResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("interest_paid"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("interestPaid")},
				{Name: sp("collateral_returned"), Number: ip(2), Label: &label, Type: &typeBytes, JsonName: sp("collateralReturned")},
			}},
			{Name: sp("MsgLiquidate"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("liquidator"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("liquidator")},
				{Name: sp("borrow_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("borrowId")},
			}},
			{Name: sp("MsgLiquidateResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("collateral_seized"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("collateralSeized")},
				{Name: sp("debt_repaid"), Number: ip(2), Label: &label, Type: &typeBytes, JsonName: sp("debtRepaid")},
			}},
			{Name: sp("MsgCreateLendingPool"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("authority"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("authority")},
				{Name: sp("denom"), Number: ip(2), Label: &label, Type: &typeString, JsonName: sp("denom")},
				{Name: sp("collateral_factor"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("collateralFactor")},
				{Name: sp("dex_pool_id"), Number: ip(4), Label: &label, Type: &typeUint64, JsonName: sp("dexPoolId")},
				{Name: sp("price_denom"), Number: ip(5), Label: &label, Type: &typeString, JsonName: sp("priceDenom")},
			}},
			{Name: sp("MsgCreateLendingPoolResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("pool_id"), Number: ip(1), Label: &label, Type: &typeUint64, JsonName: sp("poolId")},
			}},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: sp("Deposit"), InputType: sp(".syreen.lending.MsgDeposit"), OutputType: sp(".syreen.lending.MsgDepositResponse")},
					{Name: sp("Withdraw"), InputType: sp(".syreen.lending.MsgWithdraw"), OutputType: sp(".syreen.lending.MsgWithdrawResponse")},
					{Name: sp("Borrow"), InputType: sp(".syreen.lending.MsgBorrow"), OutputType: sp(".syreen.lending.MsgBorrowResponse")},
					{Name: sp("Repay"), InputType: sp(".syreen.lending.MsgRepay"), OutputType: sp(".syreen.lending.MsgRepayResponse")},
					{Name: sp("Liquidate"), InputType: sp(".syreen.lending.MsgLiquidate"), OutputType: sp(".syreen.lending.MsgLiquidateResponse")},
					{Name: sp("CreateLendingPool"), InputType: sp(".syreen.lending.MsgCreateLendingPool"), OutputType: sp(".syreen.lending.MsgCreateLendingPoolResponse")},
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
	fileDescriptorLendingTx = buf.Bytes()

	file, err := protodesc.NewFile(fd, nil)
	if err != nil {
		panic("lending: failed to create proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(file); err != nil {
		// Already registered is fine
	}
}

func sp(s string) *string { return &s }
func ip(n int32) *int32   { return &n }
