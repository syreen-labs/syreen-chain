package types

import (
	"context"

	"google.golang.org/grpc"
)

// MsgServer defines the portfolio message server interface
type MsgServer interface {
	RecordTrade(context.Context, *MsgRecordTrade) (*MsgRecordTradeResponse, error)
	CreateCompetition(context.Context, *MsgCreateCompetition) (*MsgCreateCompetitionResponse, error)
	JoinCompetition(context.Context, *MsgJoinCompetition) (*MsgJoinCompetitionResponse, error)
	EndCompetition(context.Context, *MsgEndCompetition) (*MsgEndCompetitionResponse, error)
	UpdatePortfolio(context.Context, *MsgUpdatePortfolio) (*MsgUpdatePortfolioResponse, error)
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

// Handlers

func _Msg_RecordTrade_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRecordTrade)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RecordTrade(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.portfolio.Msg/RecordTrade"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RecordTrade(ctx, req.(*MsgRecordTrade))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CreateCompetition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateCompetition)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateCompetition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.portfolio.Msg/CreateCompetition"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateCompetition(ctx, req.(*MsgCreateCompetition))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_JoinCompetition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgJoinCompetition)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).JoinCompetition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.portfolio.Msg/JoinCompetition"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).JoinCompetition(ctx, req.(*MsgJoinCompetition))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_EndCompetition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgEndCompetition)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).EndCompetition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.portfolio.Msg/EndCompetition"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).EndCompetition(ctx, req.(*MsgEndCompetition))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdatePortfolio_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdatePortfolio)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdatePortfolio(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.portfolio.Msg/UpdatePortfolio"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdatePortfolio(ctx, req.(*MsgUpdatePortfolio))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.portfolio.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "RecordTrade", Handler: _Msg_RecordTrade_Handler},
		{MethodName: "CreateCompetition", Handler: _Msg_CreateCompetition_Handler},
		{MethodName: "JoinCompetition", Handler: _Msg_JoinCompetition_Handler},
		{MethodName: "EndCompetition", Handler: _Msg_EndCompetition_Handler},
		{MethodName: "UpdatePortfolio", Handler: _Msg_UpdatePortfolio_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/portfolio/tx.proto",
}
