package types

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

// Msg response types are defined in intent_proto.go (with Descriptor/Marshal/Unmarshal).

// --- MsgServer Interface ---

type MsgServer interface {
	SubmitIntent(context.Context, *MsgSubmitIntent) (*MsgSubmitIntentResponse, error)
	RegisterSolver(context.Context, *MsgRegisterSolver) (*MsgRegisterSolverResponse, error)
	DeregisterSolver(context.Context, *MsgDeregisterSolver) (*MsgDeregisterSolverResponse, error)
	SubmitSolution(context.Context, *MsgSubmitSolution) (*MsgSubmitSolutionResponse, error)
	FulfillIntent(context.Context, *MsgFulfillIntent) (*MsgFulfillIntentResponse, error)
	SubmitChain(context.Context, *MsgSubmitChain) (*MsgSubmitChainResponse, error)
	CancelChain(context.Context, *MsgCancelChain) (*MsgCancelChainResponse, error)
	CreateStrategy(context.Context, *MsgCreateStrategy) (*MsgCreateStrategyResponse, error)
	CancelStrategy(context.Context, *MsgCancelStrategy) (*MsgCancelStrategyResponse, error)
}

// --- Query Request/Response Types ---

type QueryParamsRequest struct{}

func (m *QueryParamsRequest) ProtoMessage()           {}
func (m *QueryParamsRequest) Reset()                  { *m = QueryParamsRequest{} }
func (m *QueryParamsRequest) String() string          { return "query_params_request" }
func (m *QueryParamsRequest) XXX_MessageName() string { return "syreen.intent.QueryParamsRequest" }

type QueryParamsResponse struct {
	Params Params `json:"params"`
}

func (m *QueryParamsResponse) ProtoMessage()           {}
func (m *QueryParamsResponse) Reset()                  { *m = QueryParamsResponse{} }
func (m *QueryParamsResponse) String() string          { return fmt.Sprintf("params: %+v", m.Params) }
func (m *QueryParamsResponse) XXX_MessageName() string { return "syreen.intent.QueryParamsResponse" }

type QueryIntentRequest struct {
	IntentID string `json:"intent_id"`
}

func (m *QueryIntentRequest) ProtoMessage()           {}
func (m *QueryIntentRequest) Reset()                  { *m = QueryIntentRequest{} }
func (m *QueryIntentRequest) String() string          { return fmt.Sprintf("query_intent: id=%s", m.IntentID) }
func (m *QueryIntentRequest) XXX_MessageName() string { return "syreen.intent.QueryIntentRequest" }

type QueryIntentResponse struct {
	Intent Intent `json:"intent"`
	Found  bool   `json:"found"`
}

func (m *QueryIntentResponse) ProtoMessage()           {}
func (m *QueryIntentResponse) Reset()                  { *m = QueryIntentResponse{} }
func (m *QueryIntentResponse) String() string          { return fmt.Sprintf("intent: %+v", m.Intent) }
func (m *QueryIntentResponse) XXX_MessageName() string { return "syreen.intent.QueryIntentResponse" }

type QuerySolverRequest struct {
	Address string `json:"address"`
}

func (m *QuerySolverRequest) ProtoMessage()           {}
func (m *QuerySolverRequest) Reset()                  { *m = QuerySolverRequest{} }
func (m *QuerySolverRequest) String() string          { return fmt.Sprintf("query_solver: addr=%s", m.Address) }
func (m *QuerySolverRequest) XXX_MessageName() string { return "syreen.intent.QuerySolverRequest" }

type QuerySolverResponse struct {
	Solver Solver `json:"solver"`
	Found  bool   `json:"found"`
}

func (m *QuerySolverResponse) ProtoMessage()           {}
func (m *QuerySolverResponse) Reset()                  { *m = QuerySolverResponse{} }
func (m *QuerySolverResponse) String() string          { return fmt.Sprintf("solver: %+v", m.Solver) }
func (m *QuerySolverResponse) XXX_MessageName() string { return "syreen.intent.QuerySolverResponse" }

type QuerySolutionsRequest struct {
	IntentID string `json:"intent_id"`
}

func (m *QuerySolutionsRequest) ProtoMessage()           {}
func (m *QuerySolutionsRequest) Reset()                  { *m = QuerySolutionsRequest{} }
func (m *QuerySolutionsRequest) String() string          { return fmt.Sprintf("query_solutions: intent=%s", m.IntentID) }
func (m *QuerySolutionsRequest) XXX_MessageName() string { return "syreen.intent.QuerySolutionsRequest" }

type QuerySolutionsResponse struct {
	Solutions []Solution `json:"solutions"`
}

func (m *QuerySolutionsResponse) ProtoMessage()           {}
func (m *QuerySolutionsResponse) Reset()                  { *m = QuerySolutionsResponse{} }
func (m *QuerySolutionsResponse) String() string          { return fmt.Sprintf("solutions: %d", len(m.Solutions)) }
func (m *QuerySolutionsResponse) XXX_MessageName() string { return "syreen.intent.QuerySolutionsResponse" }

// --- Intents (list all) Query ---

type QueryIntentsRequest struct{}

func (m *QueryIntentsRequest) ProtoMessage()           {}
func (m *QueryIntentsRequest) Reset()                  { *m = QueryIntentsRequest{} }
func (m *QueryIntentsRequest) String() string          { return "query_intents_request" }
func (m *QueryIntentsRequest) XXX_MessageName() string { return "syreen.intent.QueryIntentsRequest" }

type QueryIntentsResponse struct {
	Intents []Intent `json:"intents"`
}

func (m *QueryIntentsResponse) ProtoMessage()           {}
func (m *QueryIntentsResponse) Reset()                  { *m = QueryIntentsResponse{} }
func (m *QueryIntentsResponse) String() string          { return fmt.Sprintf("intents: %d", len(m.Intents)) }
func (m *QueryIntentsResponse) XXX_MessageName() string { return "syreen.intent.QueryIntentsResponse" }

// --- IntentsByType Query ---

type QueryIntentsByTypeRequest struct {
	IntentType string `json:"intent_type"`
}

func (m *QueryIntentsByTypeRequest) ProtoMessage()           {}
func (m *QueryIntentsByTypeRequest) Reset()                  { *m = QueryIntentsByTypeRequest{} }
func (m *QueryIntentsByTypeRequest) String() string          { return fmt.Sprintf("query_intents_by_type: type=%s", m.IntentType) }
func (m *QueryIntentsByTypeRequest) XXX_MessageName() string { return "syreen.intent.QueryIntentsByTypeRequest" }

type QueryIntentsByTypeResponse struct {
	Intents []Intent `json:"intents"`
}

func (m *QueryIntentsByTypeResponse) ProtoMessage()           {}
func (m *QueryIntentsByTypeResponse) Reset()                  { *m = QueryIntentsByTypeResponse{} }
func (m *QueryIntentsByTypeResponse) String() string          { return fmt.Sprintf("intents_by_type: %d", len(m.Intents)) }
func (m *QueryIntentsByTypeResponse) XXX_MessageName() string { return "syreen.intent.QueryIntentsByTypeResponse" }

// --- QueryServer Interface ---

type QueryServer interface {
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	Intent(context.Context, *QueryIntentRequest) (*QueryIntentResponse, error)
	Intents(context.Context, *QueryIntentsRequest) (*QueryIntentsResponse, error)
	Solver(context.Context, *QuerySolverRequest) (*QuerySolverResponse, error)
	Solutions(context.Context, *QuerySolutionsRequest) (*QuerySolutionsResponse, error)
	IntentsByType(context.Context, *QueryIntentsByTypeRequest) (*QueryIntentsByTypeResponse, error)
	StrategyTemplates(context.Context, *QueryStrategyTemplatesRequest) (*QueryStrategyTemplatesResponse, error)
	Strategies(context.Context, *QueryStrategiesRequest) (*QueryStrategiesResponse, error)
	Strategy(context.Context, *QueryStrategyRequest) (*QueryStrategyResponse, error)
}

// --- Msg Handler Functions ---

func _Msg_SubmitIntent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSubmitIntent)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SubmitIntent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Msg/SubmitIntent"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SubmitIntent(ctx, req.(*MsgSubmitIntent))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_RegisterSolver_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRegisterSolver)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RegisterSolver(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Msg/RegisterSolver"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RegisterSolver(ctx, req.(*MsgRegisterSolver))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_DeregisterSolver_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgDeregisterSolver)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).DeregisterSolver(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Msg/DeregisterSolver"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).DeregisterSolver(ctx, req.(*MsgDeregisterSolver))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_SubmitSolution_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSubmitSolution)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SubmitSolution(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Msg/SubmitSolution"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SubmitSolution(ctx, req.(*MsgSubmitSolution))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_FulfillIntent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgFulfillIntent)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).FulfillIntent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Msg/FulfillIntent"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).FulfillIntent(ctx, req.(*MsgFulfillIntent))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_SubmitChain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSubmitChain)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SubmitChain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Msg/SubmitChain"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SubmitChain(ctx, req.(*MsgSubmitChain))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CancelChain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCancelChain)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CancelChain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Msg/CancelChain"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CancelChain(ctx, req.(*MsgCancelChain))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CreateStrategy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateStrategy)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateStrategy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Msg/CreateStrategy"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateStrategy(ctx, req.(*MsgCreateStrategy))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CancelStrategy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCancelStrategy)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CancelStrategy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Msg/CancelStrategy"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CancelStrategy(ctx, req.(*MsgCancelStrategy))
	}
	return interceptor(ctx, in, info, handler)
}

// --- Query Handler Functions ---

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Query/Params"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Intent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryIntentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Intent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Query/Intent"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Intent(ctx, req.(*QueryIntentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Solver_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuerySolverRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Solver(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Query/Solver"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Solver(ctx, req.(*QuerySolverRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Solutions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuerySolutionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Solutions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Query/Solutions"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Solutions(ctx, req.(*QuerySolutionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Intents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryIntentsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Intents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Query/Intents"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Intents(ctx, req.(*QueryIntentsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_IntentsByType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryIntentsByTypeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).IntentsByType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Query/IntentsByType"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).IntentsByType(ctx, req.(*QueryIntentsByTypeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_StrategyTemplates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryStrategyTemplatesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).StrategyTemplates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Query/StrategyTemplates"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).StrategyTemplates(ctx, req.(*QueryStrategyTemplatesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Strategies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryStrategiesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Strategies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Query/Strategies"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Strategies(ctx, req.(*QueryStrategiesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Strategy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryStrategyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Strategy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.intent.Query/Strategy"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Strategy(ctx, req.(*QueryStrategyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// --- Service Descriptors ---

var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.intent.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubmitIntent",
			Handler:    _Msg_SubmitIntent_Handler,
		},
		{
			MethodName: "RegisterSolver",
			Handler:    _Msg_RegisterSolver_Handler,
		},
		{
			MethodName: "DeregisterSolver",
			Handler:    _Msg_DeregisterSolver_Handler,
		},
		{
			MethodName: "SubmitSolution",
			Handler:    _Msg_SubmitSolution_Handler,
		},
		{
			MethodName: "FulfillIntent",
			Handler:    _Msg_FulfillIntent_Handler,
		},
		{
			MethodName: "SubmitChain",
			Handler:    _Msg_SubmitChain_Handler,
		},
		{
			MethodName: "CancelChain",
			Handler:    _Msg_CancelChain_Handler,
		},
		{
			MethodName: "CreateStrategy",
			Handler:    _Msg_CreateStrategy_Handler,
		},
		{
			MethodName: "CancelStrategy",
			Handler:    _Msg_CancelStrategy_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/intent/tx.proto",
}

var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.intent.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "Intent",
			Handler:    _Query_Intent_Handler,
		},
		{
			MethodName: "Intents",
			Handler:    _Query_Intents_Handler,
		},
		{
			MethodName: "Solver",
			Handler:    _Query_Solver_Handler,
		},
		{
			MethodName: "Solutions",
			Handler:    _Query_Solutions_Handler,
		},
		{
			MethodName: "IntentsByType",
			Handler:    _Query_IntentsByType_Handler,
		},
		{
			MethodName: "StrategyTemplates",
			Handler:    _Query_StrategyTemplates_Handler,
		},
		{
			MethodName: "Strategies",
			Handler:    _Query_Strategies_Handler,
		},
		{
			MethodName: "Strategy",
			Handler:    _Query_Strategy_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/intent/query.proto",
}

// RegisterMsgServer registers the MsgServer implementation with the gRPC server
func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&Msg_ServiceDesc, srv)
}

// RegisterQueryServer registers the QueryServer implementation with the gRPC server
func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
}
