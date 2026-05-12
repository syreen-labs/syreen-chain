package types

import (
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

func init() {
	registerProtoFileDescriptors()
}

func registerProtoFileDescriptors() {
	fd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/intent/tx.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.intent"),
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: sp("MsgSubmitIntent"),
				Field: []*descriptorpb.FieldDescriptorProto{
					fdStringField(1, "creator"),
					fdStringField(2, "intent_type"),
					fdStringField(3, "body"),
					fdStringField(4, "max_fee"),
					fdStringField(5, "tip"),
					fdUint64Field(6, "expiry_blocks"),
				},
			},
			{Name: sp("MsgSubmitIntentResponse")},
			{
				Name: sp("MsgRegisterSolver"),
				Field: []*descriptorpb.FieldDescriptorProto{
					fdStringField(1, "address"),
					fdStringField(2, "moniker"),
					fdStringField(3, "stake_amount"),
				},
			},
			{Name: sp("MsgRegisterSolverResponse")},
			{
				Name: sp("MsgDeregisterSolver"),
				Field: []*descriptorpb.FieldDescriptorProto{
					fdStringField(1, "address"),
				},
			},
			{Name: sp("MsgDeregisterSolverResponse")},
			{
				Name: sp("MsgSubmitSolution"),
				Field: []*descriptorpb.FieldDescriptorProto{
					fdStringField(1, "solver_addr"),
					fdStringField(2, "intent_id"),
					fdStringField(3, "execution_msgs"),
					fdStringField(4, "expected_outcome"),
				},
			},
			{Name: sp("MsgSubmitSolutionResponse")},
			{
				Name: sp("MsgFulfillIntent"),
				Field: []*descriptorpb.FieldDescriptorProto{
					fdStringField(1, "sender"),
					fdStringField(2, "solver_addr"),
					fdStringField(3, "intent_id"),
				},
			},
			{Name: sp("MsgFulfillIntentResponse")},
			{
				Name: sp("MsgSubmitChain"),
				Field: []*descriptorpb.FieldDescriptorProto{
					fdStringField(1, "creator"),
					fdStringField(2, "steps"),
					fdStringField(3, "max_fee"),
					fdStringField(4, "tip"),
					fdUint64Field(5, "expiry_blocks"),
				},
			},
			{Name: sp("MsgSubmitChainResponse")},
			{
				Name: sp("MsgCancelChain"),
				Field: []*descriptorpb.FieldDescriptorProto{
					fdStringField(1, "creator"),
					fdStringField(2, "chain_id"),
				},
			},
			{Name: sp("MsgCancelChainResponse")},
			// Strategy messages carry full field definitions so that the
			// SDK's dynamicpb-based custom signer resolver (see
			// app/encoding.go) can locate the "creator" field by name.
			{
				Name: sp("MsgCreateStrategy"),
				Field: []*descriptorpb.FieldDescriptorProto{
					fdStringField(1, "creator"),
					fdStringField(2, "template_name"),
					fdUint64Field(3, "pool_id"),
					fdStringField(4, "input_denom"),
					fdStringField(5, "output_denom"),
					fdStringField(6, "total_budget"),
					fdStringField(7, "risk_level"),
					fdUint64Field(8, "expiry_blocks"),
					fdUint64Field(9, "num_executions"),
					fdUint64Field(10, "interval_blocks"),
					fdUint64Field(11, "grid_levels"),
				},
			},
			{
				Name: sp("MsgCreateStrategyResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					fdStringField(1, "strategy_id"),
					fdStringField(2, "chain_id"),
				},
			},
			{
				Name: sp("MsgCancelStrategy"),
				Field: []*descriptorpb.FieldDescriptorProto{
					fdStringField(1, "creator"),
					fdStringField(2, "strategy_id"),
				},
			},
			{Name: sp("MsgCancelStrategyResponse")},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: sp("SubmitIntent"), InputType: sp(".syreen.intent.MsgSubmitIntent"), OutputType: sp(".syreen.intent.MsgSubmitIntentResponse")},
					{Name: sp("RegisterSolver"), InputType: sp(".syreen.intent.MsgRegisterSolver"), OutputType: sp(".syreen.intent.MsgRegisterSolverResponse")},
					{Name: sp("DeregisterSolver"), InputType: sp(".syreen.intent.MsgDeregisterSolver"), OutputType: sp(".syreen.intent.MsgDeregisterSolverResponse")},
					{Name: sp("SubmitSolution"), InputType: sp(".syreen.intent.MsgSubmitSolution"), OutputType: sp(".syreen.intent.MsgSubmitSolutionResponse")},
					{Name: sp("FulfillIntent"), InputType: sp(".syreen.intent.MsgFulfillIntent"), OutputType: sp(".syreen.intent.MsgFulfillIntentResponse")},
					{Name: sp("SubmitChain"), InputType: sp(".syreen.intent.MsgSubmitChain"), OutputType: sp(".syreen.intent.MsgSubmitChainResponse")},
					{Name: sp("CancelChain"), InputType: sp(".syreen.intent.MsgCancelChain"), OutputType: sp(".syreen.intent.MsgCancelChainResponse")},
					{Name: sp("CreateStrategy"), InputType: sp(".syreen.intent.MsgCreateStrategy"), OutputType: sp(".syreen.intent.MsgCreateStrategyResponse")},
					{Name: sp("CancelStrategy"), InputType: sp(".syreen.intent.MsgCancelStrategy"), OutputType: sp(".syreen.intent.MsgCancelStrategyResponse")},
				},
			},
		},
	}

	file, err := protodesc.NewFile(fd, nil)
	if err != nil {
		panic("intent: failed to create proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(file); err != nil {
		// Already registered is fine
	}

	qfd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/intent/query.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.intent"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: sp("QueryParamsRequest")},
			{Name: sp("QueryParamsResponse")},
			{Name: sp("QueryIntentRequest")},
			{Name: sp("QueryIntentResponse")},
			{Name: sp("QuerySolverRequest")},
			{Name: sp("QuerySolverResponse")},
			{Name: sp("QuerySolutionsRequest")},
			{Name: sp("QuerySolutionsResponse")},
			{Name: sp("QueryIntentsRequest")},
			{Name: sp("QueryIntentsResponse")},
			{Name: sp("QueryIntentsByTypeRequest")},
			{Name: sp("QueryIntentsByTypeResponse")},
			{Name: sp("QueryStrategyTemplatesRequest")},
			{Name: sp("QueryStrategyTemplatesResponse")},
			{Name: sp("QueryStrategiesRequest")},
			{Name: sp("QueryStrategiesResponse")},
			{Name: sp("QueryStrategyRequest")},
			{Name: sp("QueryStrategyResponse")},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Query"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: sp("Params"), InputType: sp(".syreen.intent.QueryParamsRequest"), OutputType: sp(".syreen.intent.QueryParamsResponse")},
					{Name: sp("Intent"), InputType: sp(".syreen.intent.QueryIntentRequest"), OutputType: sp(".syreen.intent.QueryIntentResponse")},
					{Name: sp("Intents"), InputType: sp(".syreen.intent.QueryIntentsRequest"), OutputType: sp(".syreen.intent.QueryIntentsResponse")},
					{Name: sp("Solver"), InputType: sp(".syreen.intent.QuerySolverRequest"), OutputType: sp(".syreen.intent.QuerySolverResponse")},
					{Name: sp("Solutions"), InputType: sp(".syreen.intent.QuerySolutionsRequest"), OutputType: sp(".syreen.intent.QuerySolutionsResponse")},
					{Name: sp("IntentsByType"), InputType: sp(".syreen.intent.QueryIntentsByTypeRequest"), OutputType: sp(".syreen.intent.QueryIntentsByTypeResponse")},
					{Name: sp("StrategyTemplates"), InputType: sp(".syreen.intent.QueryStrategyTemplatesRequest"), OutputType: sp(".syreen.intent.QueryStrategyTemplatesResponse")},
					{Name: sp("Strategies"), InputType: sp(".syreen.intent.QueryStrategiesRequest"), OutputType: sp(".syreen.intent.QueryStrategiesResponse")},
					{Name: sp("Strategy"), InputType: sp(".syreen.intent.QueryStrategyRequest"), OutputType: sp(".syreen.intent.QueryStrategyResponse")},
				},
			},
		},
	}
	qfile, err := protodesc.NewFile(qfd, nil)
	if err != nil {
		panic("intent: failed to create query proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(qfile); err != nil {
		// Already registered is fine
	}
}

func sp(s string) *string { return &s }

// fdStringField / fdUint64Field build FieldDescriptorProtos for use in the
// hand-built tx.proto descriptor above. They live in this file (rather than
// strategy_proto.go) so that proto_register.go remains self-contained.
func fdStringField(num int32, name string) *descriptorpb.FieldDescriptorProto {
	t := descriptorpb.FieldDescriptorProto_TYPE_STRING
	l := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	n := num
	return &descriptorpb.FieldDescriptorProto{
		Name:     sp(name),
		Number:   &n,
		Type:     &t,
		Label:    &l,
		JsonName: sp(toLowerCamel(name)),
	}
}

func fdUint64Field(num int32, name string) *descriptorpb.FieldDescriptorProto {
	t := descriptorpb.FieldDescriptorProto_TYPE_UINT64
	l := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	n := num
	return &descriptorpb.FieldDescriptorProto{
		Name:     sp(name),
		Number:   &n,
		Type:     &t,
		Label:    &l,
		JsonName: sp(toLowerCamel(name)),
	}
}
