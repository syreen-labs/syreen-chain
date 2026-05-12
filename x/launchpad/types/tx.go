package types

import (
	"context"

	"google.golang.org/grpc"
)

// MsgServer defines the launchpad message server interface
type MsgServer interface {
	CreateLaunch(context.Context, *MsgCreateLaunch) (*MsgCreateLaunchResponse, error)
	Contribute(context.Context, *MsgContribute) (*MsgContributeResponse, error)
	ClaimTokens(context.Context, *MsgClaimTokens) (*MsgClaimTokensResponse, error)
	ClaimRefund(context.Context, *MsgClaimRefund) (*MsgClaimRefundResponse, error)
	FinalizeLaunch(context.Context, *MsgFinalizeLaunch) (*MsgFinalizeLaunchResponse, error)
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

// Handlers

func _Msg_CreateLaunch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateLaunch)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateLaunch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.launchpad.Msg/CreateLaunch"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateLaunch(ctx, req.(*MsgCreateLaunch))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Contribute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgContribute)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Contribute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.launchpad.Msg/Contribute"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Contribute(ctx, req.(*MsgContribute))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ClaimTokens_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgClaimTokens)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ClaimTokens(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.launchpad.Msg/ClaimTokens"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ClaimTokens(ctx, req.(*MsgClaimTokens))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ClaimRefund_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgClaimRefund)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ClaimRefund(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.launchpad.Msg/ClaimRefund"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ClaimRefund(ctx, req.(*MsgClaimRefund))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_FinalizeLaunch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgFinalizeLaunch)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).FinalizeLaunch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.launchpad.Msg/FinalizeLaunch"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).FinalizeLaunch(ctx, req.(*MsgFinalizeLaunch))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.launchpad.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreateLaunch", Handler: _Msg_CreateLaunch_Handler},
		{MethodName: "Contribute", Handler: _Msg_Contribute_Handler},
		{MethodName: "ClaimTokens", Handler: _Msg_ClaimTokens_Handler},
		{MethodName: "ClaimRefund", Handler: _Msg_ClaimRefund_Handler},
		{MethodName: "FinalizeLaunch", Handler: _Msg_FinalizeLaunch_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/launchpad/tx.proto",
}
