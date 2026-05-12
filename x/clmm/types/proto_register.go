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
	registerCLMMProtoFileDescriptors()
}

func registerCLMMProtoFileDescriptors() {
	label := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	typeString := descriptorpb.FieldDescriptorProto_TYPE_STRING
	typeBytes := descriptorpb.FieldDescriptorProto_TYPE_BYTES
	typeUint64 := descriptorpb.FieldDescriptorProto_TYPE_UINT64
	typeInt64 := descriptorpb.FieldDescriptorProto_TYPE_INT64

	fd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/clmm/tx.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.clmm"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: sp("MsgCreateCLPool"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("denom_a"), Number: ip(2), Label: &label, Type: &typeString, JsonName: sp("denomA")},
				{Name: sp("denom_b"), Number: ip(3), Label: &label, Type: &typeString, JsonName: sp("denomB")},
				{Name: sp("tick_spacing"), Number: ip(4), Label: &label, Type: &typeInt64, JsonName: sp("tickSpacing")},
				{Name: sp("fee_rate"), Number: ip(5), Label: &label, Type: &typeBytes, JsonName: sp("feeRate")},
				{Name: sp("initial_price"), Number: ip(6), Label: &label, Type: &typeBytes, JsonName: sp("initialPrice")},
			}},
			{Name: sp("MsgCreateCLPoolResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("pool_id"), Number: ip(1), Label: &label, Type: &typeUint64, JsonName: sp("poolId")},
			}},
			{Name: sp("MsgCreatePosition"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("pool_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("poolId")},
				{Name: sp("tick_lower"), Number: ip(3), Label: &label, Type: &typeInt64, JsonName: sp("tickLower")},
				{Name: sp("tick_upper"), Number: ip(4), Label: &label, Type: &typeInt64, JsonName: sp("tickUpper")},
				{Name: sp("amount_0_desired"), Number: ip(5), Label: &label, Type: &typeBytes, JsonName: sp("amount0Desired")},
				{Name: sp("amount_1_desired"), Number: ip(6), Label: &label, Type: &typeBytes, JsonName: sp("amount1Desired")},
				{Name: sp("amount_0_min"), Number: ip(7), Label: &label, Type: &typeBytes, JsonName: sp("amount0Min")},
				{Name: sp("amount_1_min"), Number: ip(8), Label: &label, Type: &typeBytes, JsonName: sp("amount1Min")},
			}},
			{Name: sp("MsgCreatePositionResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("position_id"), Number: ip(1), Label: &label, Type: &typeUint64, JsonName: sp("positionId")},
				{Name: sp("amount_0"), Number: ip(2), Label: &label, Type: &typeBytes, JsonName: sp("amount0")},
				{Name: sp("amount_1"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("amount1")},
				{Name: sp("liquidity"), Number: ip(4), Label: &label, Type: &typeBytes, JsonName: sp("liquidity")},
			}},
			{Name: sp("MsgAddLiquidity"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("position_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("positionId")},
				{Name: sp("amount_0_desired"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("amount0Desired")},
				{Name: sp("amount_1_desired"), Number: ip(4), Label: &label, Type: &typeBytes, JsonName: sp("amount1Desired")},
			}},
			{Name: sp("MsgAddLiquidityResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("amount_0"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("amount0")},
				{Name: sp("amount_1"), Number: ip(2), Label: &label, Type: &typeBytes, JsonName: sp("amount1")},
				{Name: sp("liquidity"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("liquidity")},
			}},
			{Name: sp("MsgRemoveLiquidity"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("position_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("positionId")},
				{Name: sp("liquidity_amount"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("liquidityAmount")},
			}},
			{Name: sp("MsgRemoveLiquidityResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("amount_0"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("amount0")},
				{Name: sp("amount_1"), Number: ip(2), Label: &label, Type: &typeBytes, JsonName: sp("amount1")},
			}},
			{Name: sp("MsgCollectFees"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("position_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("positionId")},
			}},
			{Name: sp("MsgCollectFeesResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("amount_0"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("amount0")},
				{Name: sp("amount_1"), Number: ip(2), Label: &label, Type: &typeBytes, JsonName: sp("amount1")},
			}},
			{Name: sp("MsgCLSwap"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("pool_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("poolId")},
				{Name: sp("denom_in"), Number: ip(3), Label: &label, Type: &typeString, JsonName: sp("denomIn")},
				{Name: sp("amount_in"), Number: ip(4), Label: &label, Type: &typeBytes, JsonName: sp("amountIn")},
				{Name: sp("min_amount_out"), Number: ip(5), Label: &label, Type: &typeBytes, JsonName: sp("minAmountOut")},
			}},
			{Name: sp("MsgCLSwapResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("amount_out"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("amountOut")},
			}},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: sp("CreateCLPool"), InputType: sp(".syreen.clmm.MsgCreateCLPool"), OutputType: sp(".syreen.clmm.MsgCreateCLPoolResponse")},
					{Name: sp("CreatePosition"), InputType: sp(".syreen.clmm.MsgCreatePosition"), OutputType: sp(".syreen.clmm.MsgCreatePositionResponse")},
					{Name: sp("AddLiquidity"), InputType: sp(".syreen.clmm.MsgAddLiquidity"), OutputType: sp(".syreen.clmm.MsgAddLiquidityResponse")},
					{Name: sp("RemoveLiquidity"), InputType: sp(".syreen.clmm.MsgRemoveLiquidity"), OutputType: sp(".syreen.clmm.MsgRemoveLiquidityResponse")},
					{Name: sp("CollectFees"), InputType: sp(".syreen.clmm.MsgCollectFees"), OutputType: sp(".syreen.clmm.MsgCollectFeesResponse")},
					{Name: sp("CLSwap"), InputType: sp(".syreen.clmm.MsgCLSwap"), OutputType: sp(".syreen.clmm.MsgCLSwapResponse")},
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
	fileDescriptorCLMMTx = buf.Bytes()

	file, err := protodesc.NewFile(fd, nil)
	if err != nil {
		panic("clmm: failed to create proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(file); err != nil {
		// Already registered is fine
	}
}

func sp(s string) *string { return &s }
func ip(n int32) *int32   { return &n }
