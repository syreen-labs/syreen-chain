package types

import (
	"bytes"
	"compress/gzip"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

var fileDescriptorPerpsTx []byte

func init() {
	registerPerpsProtoFileDescriptors()
}

func registerPerpsProtoFileDescriptors() {
	label := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	typeString := descriptorpb.FieldDescriptorProto_TYPE_STRING
	typeBytes := descriptorpb.FieldDescriptorProto_TYPE_BYTES
	typeUint64 := descriptorpb.FieldDescriptorProto_TYPE_UINT64

	fd := &descriptorpb.FileDescriptorProto{
		Name:    strp("syreen/perps/tx.proto"),
		Syntax:  strp("proto3"),
		Package: strp("syreen.perps"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: strp("MsgOpenPosition"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
				{Name: strp("market_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("marketId")},
				{Name: strp("side"), Number: int32p(3), Label: &label, Type: &typeString, JsonName: strp("side")},
				{Name: strp("margin"), Number: int32p(4), Label: &label, Type: &typeBytes, JsonName: strp("margin")},
				{Name: strp("leverage"), Number: int32p(5), Label: &label, Type: &typeBytes, JsonName: strp("leverage")},
			}},
			{Name: strp("MsgOpenPositionResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("position_size"), Number: int32p(1), Label: &label, Type: &typeBytes, JsonName: strp("positionSize")},
				{Name: strp("entry_price"), Number: int32p(2), Label: &label, Type: &typeBytes, JsonName: strp("entryPrice")},
				{Name: strp("fee"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: strp("fee")},
			}},
			{Name: strp("MsgClosePosition"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
				{Name: strp("market_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("marketId")},
			}},
			{Name: strp("MsgClosePositionResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("realized_pnl"), Number: int32p(1), Label: &label, Type: &typeBytes, JsonName: strp("realizedPnl")},
				{Name: strp("payout"), Number: int32p(2), Label: &label, Type: &typeBytes, JsonName: strp("payout")},
			}},
			{Name: strp("MsgAddMargin"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
				{Name: strp("market_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("marketId")},
				{Name: strp("amount"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: strp("amount")},
			}},
			{Name: strp("MsgAddMarginResponse")},
			{Name: strp("MsgRemoveMargin"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
				{Name: strp("market_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("marketId")},
				{Name: strp("amount"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: strp("amount")},
			}},
			{Name: strp("MsgRemoveMarginResponse")},
			{Name: strp("MsgCreateMarket"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("authority"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("authority")},
				{Name: strp("base_denom"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: strp("baseDenom")},
				{Name: strp("quote_denom"), Number: int32p(3), Label: &label, Type: &typeString, JsonName: strp("quoteDenom")},
				{Name: strp("pool_id"), Number: int32p(4), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
				{Name: strp("max_leverage"), Number: int32p(5), Label: &label, Type: &typeBytes, JsonName: strp("maxLeverage")},
			}},
			{Name: strp("MsgCreateMarketResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("market_id"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("marketId")},
			}},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: strp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: strp("OpenPosition"), InputType: strp(".syreen.perps.MsgOpenPosition"), OutputType: strp(".syreen.perps.MsgOpenPositionResponse")},
					{Name: strp("ClosePosition"), InputType: strp(".syreen.perps.MsgClosePosition"), OutputType: strp(".syreen.perps.MsgClosePositionResponse")},
					{Name: strp("AddMargin"), InputType: strp(".syreen.perps.MsgAddMargin"), OutputType: strp(".syreen.perps.MsgAddMarginResponse")},
					{Name: strp("RemoveMargin"), InputType: strp(".syreen.perps.MsgRemoveMargin"), OutputType: strp(".syreen.perps.MsgRemoveMarginResponse")},
					{Name: strp("CreateMarket"), InputType: strp(".syreen.perps.MsgCreateMarket"), OutputType: strp(".syreen.perps.MsgCreateMarketResponse")},
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
	fileDescriptorPerpsTx = buf.Bytes()

	file, err := protodesc.NewFile(fd, nil)
	if err != nil {
		panic("perps: failed to create proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(file); err != nil {
		// Already registered is fine
	}
}

func strp(s string) *string  { return &s }
func int32p(i int32) *int32 { return &i }
