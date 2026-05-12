package types

import (
	"bytes"
	"compress/gzip"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

func init() { registerAIAgentProtoFileDescriptors() }

func registerAIAgentProtoFileDescriptors() {
	label := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	typeString := descriptorpb.FieldDescriptorProto_TYPE_STRING
	typeBytes := descriptorpb.FieldDescriptorProto_TYPE_BYTES
	typeUint64 := descriptorpb.FieldDescriptorProto_TYPE_UINT64

	fd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/aiagent/tx.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.aiagent"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: sp("MsgCreateAgent"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("owner"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("owner")},
				{Name: sp("name"), Number: ip(2), Label: &label, Type: &typeString, JsonName: sp("name")},
				{Name: sp("strategy_type"), Number: ip(3), Label: &label, Type: &typeString, JsonName: sp("strategyType")},
				{Name: sp("config"), Number: ip(4), Label: &label, Type: &typeBytes, JsonName: sp("config")},
				{Name: sp("initial_funds"), Number: ip(5), Label: &label, Type: &typeBytes, JsonName: sp("initialFunds")},
			}},
			{Name: sp("MsgCreateAgentResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("agent_id"), Number: ip(1), Label: &label, Type: &typeUint64, JsonName: sp("agentId")},
				{Name: sp("agent_address"), Number: ip(2), Label: &label, Type: &typeString, JsonName: sp("agentAddress")},
			}},
			{Name: sp("MsgFundAgent"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("owner"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("owner")},
				{Name: sp("agent_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("agentId")},
				{Name: sp("amount"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("amount")},
			}},
			{Name: sp("MsgFundAgentResponse")},
			{Name: sp("MsgWithdrawAgentFunds"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("owner"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("owner")},
				{Name: sp("agent_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("agentId")},
				{Name: sp("amount"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("amount")},
			}},
			{Name: sp("MsgWithdrawAgentFundsResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("amount_returned"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("amountReturned")},
			}},
			{Name: sp("MsgPauseAgent"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("owner"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("owner")},
				{Name: sp("agent_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("agentId")},
			}},
			{Name: sp("MsgPauseAgentResponse")},
			{Name: sp("MsgResumeAgent"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("owner"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("owner")},
				{Name: sp("agent_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("agentId")},
			}},
			{Name: sp("MsgResumeAgentResponse")},
			{Name: sp("MsgUpdateAgentStrategy"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("owner"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("owner")},
				{Name: sp("agent_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("agentId")},
				{Name: sp("config"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("config")},
			}},
			{Name: sp("MsgUpdateAgentStrategyResponse")},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: sp("CreateAgent"), InputType: sp(".syreen.aiagent.MsgCreateAgent"), OutputType: sp(".syreen.aiagent.MsgCreateAgentResponse")},
					{Name: sp("FundAgent"), InputType: sp(".syreen.aiagent.MsgFundAgent"), OutputType: sp(".syreen.aiagent.MsgFundAgentResponse")},
					{Name: sp("WithdrawAgentFunds"), InputType: sp(".syreen.aiagent.MsgWithdrawAgentFunds"), OutputType: sp(".syreen.aiagent.MsgWithdrawAgentFundsResponse")},
					{Name: sp("PauseAgent"), InputType: sp(".syreen.aiagent.MsgPauseAgent"), OutputType: sp(".syreen.aiagent.MsgPauseAgentResponse")},
					{Name: sp("ResumeAgent"), InputType: sp(".syreen.aiagent.MsgResumeAgent"), OutputType: sp(".syreen.aiagent.MsgResumeAgentResponse")},
					{Name: sp("UpdateAgentStrategy"), InputType: sp(".syreen.aiagent.MsgUpdateAgentStrategy"), OutputType: sp(".syreen.aiagent.MsgUpdateAgentStrategyResponse")},
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
	fileDescriptorAIAgentTx = buf.Bytes()

	file, err := protodesc.NewFile(fd, nil)
	if err != nil {
		panic("aiagent: failed to create proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(file); err != nil {
		// already registered is fine
	}
}

func sp(s string) *string { return &s }
func ip(n int32) *int32   { return &n }
