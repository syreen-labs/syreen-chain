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
	registerPerpsProtoFileDescriptors()
}

func registerPerpsProtoFileDescriptors() {
	label := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	typeString := descriptorpb.FieldDescriptorProto_TYPE_STRING
	typeBytes := descriptorpb.FieldDescriptorProto_TYPE_BYTES
	typeUint64 := descriptorpb.FieldDescriptorProto_TYPE_UINT64

	fd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/perps/tx.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.perps"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: sp("MsgOpenPosition"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("market_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("marketId")},
				{Name: sp("side"), Number: ip(3), Label: &label, Type: &typeString, JsonName: sp("side")},
				{Name: sp("margin"), Number: ip(4), Label: &label, Type: &typeBytes, JsonName: sp("margin")},
				{Name: sp("leverage"), Number: ip(5), Label: &label, Type: &typeBytes, JsonName: sp("leverage")},
			}},
			{Name: sp("MsgOpenPositionResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("position_size"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("positionSize")},
				{Name: sp("entry_price"), Number: ip(2), Label: &label, Type: &typeBytes, JsonName: sp("entryPrice")},
				{Name: sp("fee"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("fee")},
			}},
			{Name: sp("MsgClosePosition"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("market_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("marketId")},
			}},
			{Name: sp("MsgClosePositionResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("realized_pnl"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("realizedPnl")},
				{Name: sp("payout"), Number: ip(2), Label: &label, Type: &typeBytes, JsonName: sp("payout")},
			}},
			{Name: sp("MsgAddMargin"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("market_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("marketId")},
				{Name: sp("amount"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("amount")},
			}},
			{Name: sp("MsgAddMarginResponse")},
			{Name: sp("MsgRemoveMargin"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("market_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("marketId")},
				{Name: sp("amount"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("amount")},
			}},
			{Name: sp("MsgRemoveMarginResponse")},
			{Name: sp("MsgCreateMarket"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("authority"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("authority")},
				{Name: sp("base_denom"), Number: ip(2), Label: &label, Type: &typeString, JsonName: sp("baseDenom")},
				{Name: sp("quote_denom"), Number: ip(3), Label: &label, Type: &typeString, JsonName: sp("quoteDenom")},
				{Name: sp("pool_id"), Number: ip(4), Label: &label, Type: &typeUint64, JsonName: sp("poolId")},
				{Name: sp("max_leverage"), Number: ip(5), Label: &label, Type: &typeBytes, JsonName: sp("maxLeverage")},
			}},
			{Name: sp("MsgCreateMarketResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("market_id"), Number: ip(1), Label: &label, Type: &typeUint64, JsonName: sp("marketId")},
			}},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: sp("OpenPosition"), InputType: sp(".syreen.perps.MsgOpenPosition"), OutputType: sp(".syreen.perps.MsgOpenPositionResponse")},
					{Name: sp("ClosePosition"), InputType: sp(".syreen.perps.MsgClosePosition"), OutputType: sp(".syreen.perps.MsgClosePositionResponse")},
					{Name: sp("AddMargin"), InputType: sp(".syreen.perps.MsgAddMargin"), OutputType: sp(".syreen.perps.MsgAddMarginResponse")},
					{Name: sp("RemoveMargin"), InputType: sp(".syreen.perps.MsgRemoveMargin"), OutputType: sp(".syreen.perps.MsgRemoveMarginResponse")},
					{Name: sp("CreateMarket"), InputType: sp(".syreen.perps.MsgCreateMarket"), OutputType: sp(".syreen.perps.MsgCreateMarketResponse")},
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
	fileDescriptorPerpsTx = buf.Bytes()

	file, err := protodesc.NewFile(fd, nil)
	if err != nil {
		panic("perps: failed to create proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(file); err != nil {
		// Already registered is fine
	}
}

func sp(s string) *string { return &s }
func ip(n int32) *int32   { return &n }
