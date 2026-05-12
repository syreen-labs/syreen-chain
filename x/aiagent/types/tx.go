package types

import (
	"context"

	"google.golang.org/grpc"
)

type MsgServer interface {
	CreateAgent(context.Context, *MsgCreateAgent) (*MsgCreateAgentResponse, error)
	FundAgent(context.Context, *MsgFundAgent) (*MsgFundAgentResponse, error)
	WithdrawAgentFunds(context.Context, *MsgWithdrawAgentFunds) (*MsgWithdrawAgentFundsResponse, error)
	PauseAgent(context.Context, *MsgPauseAgent) (*MsgPauseAgentResponse, error)
	ResumeAgent(context.Context, *MsgResumeAgent) (*MsgResumeAgentResponse, error)
	UpdateAgentStrategy(context.Context, *MsgUpdateAgentStrategy) (*MsgUpdateAgentStrategyResponse, error)
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&_AIAgentMsg_serviceDesc, srv)
}

func _Msg_CreateAgent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateAgent)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).CreateAgent(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.aiagent.Msg/CreateAgent"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).CreateAgent(ctx, req.(*MsgCreateAgent)) }
	return interceptor(ctx, in, info, handler)
}

func _Msg_FundAgent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgFundAgent)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).FundAgent(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.aiagent.Msg/FundAgent"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).FundAgent(ctx, req.(*MsgFundAgent)) }
	return interceptor(ctx, in, info, handler)
}

func _Msg_WithdrawAgentFunds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgWithdrawAgentFunds)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).WithdrawAgentFunds(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.aiagent.Msg/WithdrawAgentFunds"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).WithdrawAgentFunds(ctx, req.(*MsgWithdrawAgentFunds)) }
	return interceptor(ctx, in, info, handler)
}

func _Msg_PauseAgent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgPauseAgent)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).PauseAgent(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.aiagent.Msg/PauseAgent"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).PauseAgent(ctx, req.(*MsgPauseAgent)) }
	return interceptor(ctx, in, info, handler)
}

func _Msg_ResumeAgent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgResumeAgent)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).ResumeAgent(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.aiagent.Msg/ResumeAgent"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).ResumeAgent(ctx, req.(*MsgResumeAgent)) }
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateAgentStrategy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateAgentStrategy)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).UpdateAgentStrategy(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.aiagent.Msg/UpdateAgentStrategy"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).UpdateAgentStrategy(ctx, req.(*MsgUpdateAgentStrategy)) }
	return interceptor(ctx, in, info, handler)
}

var _AIAgentMsg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.aiagent.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreateAgent", Handler: _Msg_CreateAgent_Handler},
		{MethodName: "FundAgent", Handler: _Msg_FundAgent_Handler},
		{MethodName: "WithdrawAgentFunds", Handler: _Msg_WithdrawAgentFunds_Handler},
		{MethodName: "PauseAgent", Handler: _Msg_PauseAgent_Handler},
		{MethodName: "ResumeAgent", Handler: _Msg_ResumeAgent_Handler},
		{MethodName: "UpdateAgentStrategy", Handler: _Msg_UpdateAgentStrategy_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/aiagent/tx.proto",
}
