package types

import (
	"context"

	"google.golang.org/grpc"
)

// MsgServer defines the lending message server interface
type MsgServer interface {
	Deposit(context.Context, *MsgDeposit) (*MsgDepositResponse, error)
	Withdraw(context.Context, *MsgWithdraw) (*MsgWithdrawResponse, error)
	Borrow(context.Context, *MsgBorrow) (*MsgBorrowResponse, error)
	Repay(context.Context, *MsgRepay) (*MsgRepayResponse, error)
	Liquidate(context.Context, *MsgLiquidate) (*MsgLiquidateResponse, error)
	CreateLendingPool(context.Context, *MsgCreateLendingPool) (*MsgCreateLendingPoolResponse, error)
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

// Handlers

func _Msg_Deposit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgDeposit)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Deposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.lending.Msg/Deposit"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Deposit(ctx, req.(*MsgDeposit))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Withdraw_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgWithdraw)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Withdraw(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.lending.Msg/Withdraw"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Withdraw(ctx, req.(*MsgWithdraw))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Borrow_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgBorrow)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Borrow(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.lending.Msg/Borrow"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Borrow(ctx, req.(*MsgBorrow))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Repay_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRepay)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Repay(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.lending.Msg/Repay"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Repay(ctx, req.(*MsgRepay))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Liquidate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgLiquidate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Liquidate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.lending.Msg/Liquidate"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Liquidate(ctx, req.(*MsgLiquidate))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CreateLendingPool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateLendingPool)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateLendingPool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.lending.Msg/CreateLendingPool"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateLendingPool(ctx, req.(*MsgCreateLendingPool))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.lending.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Deposit", Handler: _Msg_Deposit_Handler},
		{MethodName: "Withdraw", Handler: _Msg_Withdraw_Handler},
		{MethodName: "Borrow", Handler: _Msg_Borrow_Handler},
		{MethodName: "Repay", Handler: _Msg_Repay_Handler},
		{MethodName: "Liquidate", Handler: _Msg_Liquidate_Handler},
		{MethodName: "CreateLendingPool", Handler: _Msg_CreateLendingPool_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/lending/tx.proto",
}
