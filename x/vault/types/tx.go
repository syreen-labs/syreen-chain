package types

import (
	"context"

	"google.golang.org/grpc"
)

// MsgServer defines the vault message server interface
type MsgServer interface {
	CreateVault(context.Context, *MsgCreateVault) (*MsgCreateVaultResponse, error)
	DepositVault(context.Context, *MsgDepositVault) (*MsgDepositVaultResponse, error)
	WithdrawVault(context.Context, *MsgWithdrawVault) (*MsgWithdrawVaultResponse, error)
	CompoundVault(context.Context, *MsgCompoundVault) (*MsgCompoundVaultResponse, error)
	UpdateVaultStrategy(context.Context, *MsgUpdateVaultStrategy) (*MsgUpdateVaultStrategyResponse, error)
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

// Handlers

func _Msg_CreateVault_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateVault)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateVault(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.vault.Msg/CreateVault"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateVault(ctx, req.(*MsgCreateVault))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_DepositVault_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgDepositVault)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).DepositVault(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.vault.Msg/DepositVault"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).DepositVault(ctx, req.(*MsgDepositVault))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_WithdrawVault_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgWithdrawVault)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).WithdrawVault(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.vault.Msg/WithdrawVault"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).WithdrawVault(ctx, req.(*MsgWithdrawVault))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CompoundVault_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCompoundVault)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CompoundVault(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.vault.Msg/CompoundVault"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CompoundVault(ctx, req.(*MsgCompoundVault))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateVaultStrategy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateVaultStrategy)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateVaultStrategy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.vault.Msg/UpdateVaultStrategy"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateVaultStrategy(ctx, req.(*MsgUpdateVaultStrategy))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.vault.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreateVault", Handler: _Msg_CreateVault_Handler},
		{MethodName: "DepositVault", Handler: _Msg_DepositVault_Handler},
		{MethodName: "WithdrawVault", Handler: _Msg_WithdrawVault_Handler},
		{MethodName: "CompoundVault", Handler: _Msg_CompoundVault_Handler},
		{MethodName: "UpdateVaultStrategy", Handler: _Msg_UpdateVaultStrategy_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/vault/tx.proto",
}
