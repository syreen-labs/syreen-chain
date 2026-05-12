package types

import (
	"bytes"
	"compress/gzip"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

var fileDescriptorCLMMTx []byte

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
		Name:    strp("syreen/clmm/tx.proto"),
		Syntax:  strp("proto3"),
		Package: strp("syreen.clmm"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: strp("MsgCreateCLPool"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
				{Name: strp("denom_a"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: strp("denomA")},
				{Name: strp("denom_b"), Number: int32p(3), Label: &label, Type: &typeString, JsonName: strp("denomB")},
				{Name: strp("tick_spacing"), Number: int32p(4), Label: &label, Type: &typeInt64, JsonName: strp("tickSpacing")},
				{Name: strp("fee_rate"), Number: int32p(5), Label: &label, Type: &typeBytes, JsonName: strp("feeRate")},
				{Name: strp("initial_price"), Number: int32p(6), Label: &label, Type: &typeBytes, JsonName: strp("initialPrice")},
			}},
			{Name: strp("MsgCreateCLPoolResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("pool_id"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
			}},
			{Name: strp("MsgCreatePosition"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
				{Name: strp("pool_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
				{Name: strp("tick_lower"), Number: int32p(3), Label: &label, Type: &typeInt64, JsonName: strp("tickLower")},
				{Name: strp("tick_upper"), Number: int32p(4), Label: &label, Type: &typeInt64, JsonName: strp("tickUpper")},
				{Name: strp("amount_0_desired"), Number: int32p(5), Label: &label, Type: &typeBytes, JsonName: strp("amount0Desired")},
				{Name: strp("amount_1_desired"), Number: int32p(6), Label: &label, Type: &typeBytes, JsonName: strp("amount1Desired")},
				{Name: strp("amount_0_min"), Number: int32p(7), Label: &label, Type: &typeBytes, JsonName: strp("amount0Min")},
				{Name: strp("amount_1_min"), Number: int32p(8), Label: &label, Type: &typeBytes, JsonName: strp("amount1Min")},
			}},
			{Name: strp("MsgCreatePositionResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("position_id"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("positionId")},
				{Name: strp("amount_0"), Number: int32p(2), Label: &label, Type: &typeBytes, JsonName: strp("amount0")},
				{Name: strp("amount_1"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: strp("amount1")},
				{Name: strp("liquidity"), Number: int32p(4), Label: &label, Type: &typeBytes, JsonName: strp("liquidity")},
			}},
			{Name: strp("MsgAddLiquidity"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
				{Name: strp("position_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("positionId")},
				{Name: strp("amount_0_desired"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: strp("amount0Desired")},
				{Name: strp("amount_1_desired"), Number: int32p(4), Label: &label, Type: &typeBytes, JsonName: strp("amount1Desired")},
			}},
			{Name: strp("MsgAddLiquidityResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("amount_0"), Number: int32p(1), Label: &label, Type: &typeBytes, JsonName: strp("amount0")},
				{Name: strp("amount_1"), Number: int32p(2), Label: &label, Type: &typeBytes, JsonName: strp("amount1")},
				{Name: strp("liquidity"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: strp("liquidity")},
			}},
			{Name: strp("MsgRemoveLiquidity"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
				{Name: strp("position_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("positionId")},
				{Name: strp("liquidity_amount"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: strp("liquidityAmount")},
			}},
			{Name: strp("MsgRemoveLiquidityResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("amount_0"), Number: int32p(1), Label: &label, Type: &typeBytes, JsonName: strp("amount0")},
				{Name: strp("amount_1"), Number: int32p(2), Label: &label, Type: &typeBytes, JsonName: strp("amount1")},
			}},
			{Name: strp("MsgCollectFees"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
				{Name: strp("position_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("positionId")},
			}},
			{Name: strp("MsgCollectFeesResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("amount_0"), Number: int32p(1), Label: &label, Type: &typeBytes, JsonName: strp("amount0")},
				{Name: strp("amount_1"), Number: int32p(2), Label: &label, Type: &typeBytes, JsonName: strp("amount1")},
			}},
			{Name: strp("MsgCLSwap"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
				{Name: strp("pool_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
				{Name: strp("denom_in"), Number: int32p(3), Label: &label, Type: &typeString, JsonName: strp("denomIn")},
				{Name: strp("amount_in"), Number: int32p(4), Label: &label, Type: &typeBytes, JsonName: strp("amountIn")},
				{Name: strp("min_amount_out"), Number: int32p(5), Label: &label, Type: &typeBytes, JsonName: strp("minAmountOut")},
			}},
			{Name: strp("MsgCLSwapResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("amount_out"), Number: int32p(1), Label: &label, Type: &typeBytes, JsonName: strp("amountOut")},
			}},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: strp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: strp("CreateCLPool"), InputType: strp(".syreen.clmm.MsgCreateCLPool"), OutputType: strp(".syreen.clmm.MsgCreateCLPoolResponse")},
					{Name: strp("CreatePosition"), InputType: strp(".syreen.clmm.MsgCreatePosition"), OutputType: strp(".syreen.clmm.MsgCreatePositionResponse")},
					{Name: strp("AddLiquidity"), InputType: strp(".syreen.clmm.MsgAddLiquidity"), OutputType: strp(".syreen.clmm.MsgAddLiquidityResponse")},
					{Name: strp("RemoveLiquidity"), InputType: strp(".syreen.clmm.MsgRemoveLiquidity"), OutputType: strp(".syreen.clmm.MsgRemoveLiquidityResponse")},
					{Name: strp("CollectFees"), InputType: strp(".syreen.clmm.MsgCollectFees"), OutputType: strp(".syreen.clmm.MsgCollectFeesResponse")},
					{Name: strp("CLSwap"), InputType: strp(".syreen.clmm.MsgCLSwap"), OutputType: strp(".syreen.clmm.MsgCLSwapResponse")},
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
	fileDescriptorCLMMTx = buf.Bytes()

	file, err := protodesc.NewFile(fd, nil)
	if err != nil {
		panic("clmm: failed to create proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(file); err != nil {
		// Already registered is fine
	}
}

func strp(s string) *string  { return &s }
func int32p(i int32) *int32 { return &i }
