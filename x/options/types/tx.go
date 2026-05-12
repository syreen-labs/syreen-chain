package types

import (
	"context"

	"google.golang.org/grpc"
)

type MsgServer interface {
	WriteOption(context.Context, *MsgWriteOption) (*MsgWriteOptionResponse, error)
	BuyOption(context.Context, *MsgBuyOption) (*MsgBuyOptionResponse, error)
	ExerciseOption(context.Context, *MsgExerciseOption) (*MsgExerciseOptionResponse, error)
	CancelOption(context.Context, *MsgCancelOption) (*MsgCancelOptionResponse, error)
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_WriteOption_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgWriteOption)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).WriteOption(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.options.Msg/WriteOption"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).WriteOption(ctx, req.(*MsgWriteOption))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_BuyOption_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgBuyOption)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).BuyOption(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.options.Msg/BuyOption"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).BuyOption(ctx, req.(*MsgBuyOption))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ExerciseOption_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgExerciseOption)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ExerciseOption(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.options.Msg/ExerciseOption"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ExerciseOption(ctx, req.(*MsgExerciseOption))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CancelOption_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCancelOption)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CancelOption(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.options.Msg/CancelOption"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CancelOption(ctx, req.(*MsgCancelOption))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.options.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "WriteOption", Handler: _Msg_WriteOption_Handler},
		{MethodName: "BuyOption", Handler: _Msg_BuyOption_Handler},
		{MethodName: "ExerciseOption", Handler: _Msg_ExerciseOption_Handler},
		{MethodName: "CancelOption", Handler: _Msg_CancelOption_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/options/tx.proto",
}
