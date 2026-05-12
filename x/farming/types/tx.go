package types

import (
	"context"

	"google.golang.org/grpc"
)

// MsgServer defines the farming message server interface
type MsgServer interface {
	Stake(context.Context, *MsgStake) (*MsgStakeResponse, error)
	Unstake(context.Context, *MsgUnstake) (*MsgUnstakeResponse, error)
	ClaimReward(context.Context, *MsgClaimReward) (*MsgClaimRewardResponse, error)
	CreateFarm(context.Context, *MsgCreateFarm) (*MsgCreateFarmResponse, error)
	UpdateFarm(context.Context, *MsgUpdateFarm) (*MsgUpdateFarmResponse, error)
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

// Handlers

func _Msg_Stake_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgStake)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Stake(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.farming.Msg/Stake"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Stake(ctx, req.(*MsgStake))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Unstake_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUnstake)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Unstake(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.farming.Msg/Unstake"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Unstake(ctx, req.(*MsgUnstake))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ClaimReward_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgClaimReward)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ClaimReward(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.farming.Msg/ClaimReward"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ClaimReward(ctx, req.(*MsgClaimReward))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CreateFarm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateFarm)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateFarm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.farming.Msg/CreateFarm"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateFarm(ctx, req.(*MsgCreateFarm))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateFarm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateFarm)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateFarm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.farming.Msg/UpdateFarm"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateFarm(ctx, req.(*MsgUpdateFarm))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.farming.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Stake", Handler: _Msg_Stake_Handler},
		{MethodName: "Unstake", Handler: _Msg_Unstake_Handler},
		{MethodName: "ClaimReward", Handler: _Msg_ClaimReward_Handler},
		{MethodName: "CreateFarm", Handler: _Msg_CreateFarm_Handler},
		{MethodName: "UpdateFarm", Handler: _Msg_UpdateFarm_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/farming/tx.proto",
}
