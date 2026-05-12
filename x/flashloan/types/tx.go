package types

import (
	"context"

	"google.golang.org/grpc"
)

// MsgServer defines the flashloan message server interface
type MsgServer interface {
	FlashLoan(context.Context, *MsgFlashLoan) (*MsgFlashLoanResponse, error)
	CreateFlashPool(context.Context, *MsgCreateFlashPool) (*MsgCreateFlashPoolResponse, error)
	FundFlashPool(context.Context, *MsgFundFlashPool) (*MsgFundFlashPoolResponse, error)
	WithdrawFlashPool(context.Context, *MsgWithdrawFlashPool) (*MsgWithdrawFlashPoolResponse, error)
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

// Handlers

func _Msg_FlashLoan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgFlashLoan)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).FlashLoan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.flashloan.Msg/FlashLoan"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).FlashLoan(ctx, req.(*MsgFlashLoan))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CreateFlashPool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateFlashPool)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateFlashPool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.flashloan.Msg/CreateFlashPool"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateFlashPool(ctx, req.(*MsgCreateFlashPool))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_FundFlashPool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgFundFlashPool)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).FundFlashPool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.flashloan.Msg/FundFlashPool"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).FundFlashPool(ctx, req.(*MsgFundFlashPool))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_WithdrawFlashPool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgWithdrawFlashPool)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).WithdrawFlashPool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.flashloan.Msg/WithdrawFlashPool"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).WithdrawFlashPool(ctx, req.(*MsgWithdrawFlashPool))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.flashloan.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "FlashLoan", Handler: _Msg_FlashLoan_Handler},
		{MethodName: "CreateFlashPool", Handler: _Msg_CreateFlashPool_Handler},
		{MethodName: "FundFlashPool", Handler: _Msg_FundFlashPool_Handler},
		{MethodName: "WithdrawFlashPool", Handler: _Msg_WithdrawFlashPool_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/flashloan/tx.proto",
}
