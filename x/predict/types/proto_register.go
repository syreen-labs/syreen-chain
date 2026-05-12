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
	registerPredictProtoFileDescriptors()
}

func registerPredictProtoFileDescriptors() {
	label := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	typeString := descriptorpb.FieldDescriptorProto_TYPE_STRING
	typeBytes := descriptorpb.FieldDescriptorProto_TYPE_BYTES
	typeUint64 := descriptorpb.FieldDescriptorProto_TYPE_UINT64
	typeInt64 := descriptorpb.FieldDescriptorProto_TYPE_INT64

	fd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/predict/tx.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.predict"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: sp("MsgCreateMarket"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("creator"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("creator")},
				{Name: sp("question"), Number: ip(2), Label: &label, Type: &typeString, JsonName: sp("question")},
				{Name: sp("resolver"), Number: ip(3), Label: &label, Type: &typeString, JsonName: sp("resolver")},
				{Name: sp("quote_denom"), Number: ip(4), Label: &label, Type: &typeString, JsonName: sp("quoteDenom")},
				{Name: sp("resolution_block"), Number: ip(5), Label: &label, Type: &typeInt64, JsonName: sp("resolutionBlock")},
				{Name: sp("initial_liquidity"), Number: ip(6), Label: &label, Type: &typeBytes, JsonName: sp("initialLiquidity")},
			}},
			{Name: sp("MsgCreateMarketResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("market_id"), Number: ip(1), Label: &label, Type: &typeUint64, JsonName: sp("marketId")},
			}},
			{Name: sp("MsgBuyShares"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("market_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("marketId")},
				{Name: sp("outcome"), Number: ip(3), Label: &label, Type: &typeString, JsonName: sp("outcome")},
				{Name: sp("amount"), Number: ip(4), Label: &label, Type: &typeBytes, JsonName: sp("amount")},
			}},
			{Name: sp("MsgBuySharesResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("shares_bought"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("sharesBought")},
				{Name: sp("avg_price"), Number: ip(2), Label: &label, Type: &typeBytes, JsonName: sp("avgPrice")},
			}},
			{Name: sp("MsgSellShares"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("market_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("marketId")},
				{Name: sp("outcome"), Number: ip(3), Label: &label, Type: &typeString, JsonName: sp("outcome")},
				{Name: sp("shares"), Number: ip(4), Label: &label, Type: &typeBytes, JsonName: sp("shares")},
			}},
			{Name: sp("MsgSellSharesResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("quote_returned"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("quoteReturned")},
			}},
			{Name: sp("MsgResolveMarket"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("resolver"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("resolver")},
				{Name: sp("market_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("marketId")},
				{Name: sp("outcome"), Number: ip(3), Label: &label, Type: &typeString, JsonName: sp("outcome")},
			}},
			{Name: sp("MsgResolveMarketResponse")},
			{Name: sp("MsgClaimWinnings"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("market_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("marketId")},
			}},
			{Name: sp("MsgClaimWinningsResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("amount"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("amount")},
			}},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: sp("CreateMarket"), InputType: sp(".syreen.predict.MsgCreateMarket"), OutputType: sp(".syreen.predict.MsgCreateMarketResponse")},
					{Name: sp("BuyShares"), InputType: sp(".syreen.predict.MsgBuyShares"), OutputType: sp(".syreen.predict.MsgBuySharesResponse")},
					{Name: sp("SellShares"), InputType: sp(".syreen.predict.MsgSellShares"), OutputType: sp(".syreen.predict.MsgSellSharesResponse")},
					{Name: sp("ResolveMarket"), InputType: sp(".syreen.predict.MsgResolveMarket"), OutputType: sp(".syreen.predict.MsgResolveMarketResponse")},
					{Name: sp("ClaimWinnings"), InputType: sp(".syreen.predict.MsgClaimWinnings"), OutputType: sp(".syreen.predict.MsgClaimWinningsResponse")},
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
	fileDescriptorPredictTx = buf.Bytes()

	file, err := protodesc.NewFile(fd, nil)
	if err != nil {
		panic("predict: failed to create proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(file); err != nil {
		// Already registered is fine
	}
}

func sp(s string) *string { return &s }
func ip(n int32) *int32   { return &n }
