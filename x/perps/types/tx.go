package types

import (
	"context"

	"google.golang.org/grpc"
)

// MsgServer defines the perps message server interface
type MsgServer interface {
	OpenPosition(context.Context, *MsgOpenPosition) (*MsgOpenPositionResponse, error)
	ClosePosition(context.Context, *MsgClosePosition) (*MsgClosePositionResponse, error)
	AddMargin(context.Context, *MsgAddMargin) (*MsgAddMarginResponse, error)
	RemoveMargin(context.Context, *MsgRemoveMargin) (*MsgRemoveMarginResponse, error)
	CreateMarket(context.Context, *MsgCreateMarket) (*MsgCreateMarketResponse, error)
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

// Handlers

func _Msg_OpenPosition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgOpenPosition)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).OpenPosition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.perps.Msg/OpenPosition"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).OpenPosition(ctx, req.(*MsgOpenPosition))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ClosePosition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgClosePosition)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ClosePosition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.perps.Msg/ClosePosition"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ClosePosition(ctx, req.(*MsgClosePosition))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_AddMargin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgAddMargin)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).AddMargin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.perps.Msg/AddMargin"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).AddMargin(ctx, req.(*MsgAddMargin))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_RemoveMargin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRemoveMargin)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RemoveMargin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.perps.Msg/RemoveMargin"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RemoveMargin(ctx, req.(*MsgRemoveMargin))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CreateMarket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateMarket)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateMarket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.perps.Msg/CreateMarket"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateMarket(ctx, req.(*MsgCreateMarket))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.perps.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "OpenPosition", Handler: _Msg_OpenPosition_Handler},
		{MethodName: "ClosePosition", Handler: _Msg_ClosePosition_Handler},
		{MethodName: "AddMargin", Handler: _Msg_AddMargin_Handler},
		{MethodName: "RemoveMargin", Handler: _Msg_RemoveMargin_Handler},
		{MethodName: "CreateMarket", Handler: _Msg_CreateMarket_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/perps/tx.proto",
}
