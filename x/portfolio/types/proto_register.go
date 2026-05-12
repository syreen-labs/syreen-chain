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
	registerPortfolioProtoFileDescriptors()
}

func registerPortfolioProtoFileDescriptors() {
	label := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	typeString := descriptorpb.FieldDescriptorProto_TYPE_STRING
	typeBytes := descriptorpb.FieldDescriptorProto_TYPE_BYTES
	typeUint64 := descriptorpb.FieldDescriptorProto_TYPE_UINT64
	typeInt64 := descriptorpb.FieldDescriptorProto_TYPE_INT64

	fd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/portfolio/tx.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.portfolio"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: sp("MsgRecordTrade"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("trade_type"), Number: ip(2), Label: &label, Type: &typeString, JsonName: sp("tradeType")},
				{Name: sp("denom"), Number: ip(3), Label: &label, Type: &typeString, JsonName: sp("denom")},
				{Name: sp("amount"), Number: ip(4), Label: &label, Type: &typeBytes, JsonName: sp("amount")},
				{Name: sp("price"), Number: ip(5), Label: &label, Type: &typeBytes, JsonName: sp("price")},
				{Name: sp("pnl"), Number: ip(6), Label: &label, Type: &typeBytes, JsonName: sp("pnl")},
			}},
			{Name: sp("MsgRecordTradeResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("trade_id"), Number: ip(1), Label: &label, Type: &typeUint64, JsonName: sp("tradeId")},
			}},
			{Name: sp("MsgCreateCompetition"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("creator"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("creator")},
				{Name: sp("name"), Number: ip(2), Label: &label, Type: &typeString, JsonName: sp("name")},
				{Name: sp("start_block"), Number: ip(3), Label: &label, Type: &typeInt64, JsonName: sp("startBlock")},
				{Name: sp("end_block"), Number: ip(4), Label: &label, Type: &typeInt64, JsonName: sp("endBlock")},
				{Name: sp("prize_denom"), Number: ip(5), Label: &label, Type: &typeString, JsonName: sp("prizeDenom")},
				{Name: sp("prize_pool"), Number: ip(6), Label: &label, Type: &typeBytes, JsonName: sp("prizePool")},
				{Name: sp("entry_fee"), Number: ip(7), Label: &label, Type: &typeBytes, JsonName: sp("entryFee")},
				{Name: sp("max_participants"), Number: ip(8), Label: &label, Type: &typeUint64, JsonName: sp("maxParticipants")},
			}},
			{Name: sp("MsgCreateCompetitionResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("competition_id"), Number: ip(1), Label: &label, Type: &typeUint64, JsonName: sp("competitionId")},
			}},
			{Name: sp("MsgJoinCompetition"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("competition_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("competitionId")},
			}},
			{Name: sp("MsgJoinCompetitionResponse")},
			{Name: sp("MsgEndCompetition"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("authority"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("authority")},
				{Name: sp("competition_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("competitionId")},
			}},
			{Name: sp("MsgEndCompetitionResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("winners"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("winners")},
			}},
			{Name: sp("MsgUpdatePortfolio"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
			}},
			{Name: sp("MsgUpdatePortfolioResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("total_value"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("totalValue")},
			}},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: sp("RecordTrade"), InputType: sp(".syreen.portfolio.MsgRecordTrade"), OutputType: sp(".syreen.portfolio.MsgRecordTradeResponse")},
					{Name: sp("CreateCompetition"), InputType: sp(".syreen.portfolio.MsgCreateCompetition"), OutputType: sp(".syreen.portfolio.MsgCreateCompetitionResponse")},
					{Name: sp("JoinCompetition"), InputType: sp(".syreen.portfolio.MsgJoinCompetition"), OutputType: sp(".syreen.portfolio.MsgJoinCompetitionResponse")},
					{Name: sp("EndCompetition"), InputType: sp(".syreen.portfolio.MsgEndCompetition"), OutputType: sp(".syreen.portfolio.MsgEndCompetitionResponse")},
					{Name: sp("UpdatePortfolio"), InputType: sp(".syreen.portfolio.MsgUpdatePortfolio"), OutputType: sp(".syreen.portfolio.MsgUpdatePortfolioResponse")},
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
	fileDescriptorPortfolioTx = buf.Bytes()

	file, err := protodesc.NewFile(fd, nil)
	if err != nil {
		panic("portfolio: failed to create proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(file); err != nil {
		// Already registered is fine
	}
}

func sp(s string) *string { return &s }
func ip(n int32) *int32   { return &n }
