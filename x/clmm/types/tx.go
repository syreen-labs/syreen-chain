package types

import (
	"context"

	"google.golang.org/grpc"
)

// MsgServer defines the clmm message server interface
type MsgServer interface {
	CreateCLPool(context.Context, *MsgCreateCLPool) (*MsgCreateCLPoolResponse, error)
	CreatePosition(context.Context, *MsgCreatePosition) (*MsgCreatePositionResponse, error)
	AddLiquidity(context.Context, *MsgAddLiquidity) (*MsgAddLiquidityResponse, error)
	RemoveLiquidity(context.Context, *MsgRemoveLiquidity) (*MsgRemoveLiquidityResponse, error)
	CollectFees(context.Context, *MsgCollectFees) (*MsgCollectFeesResponse, error)
	CLSwap(context.Context, *MsgCLSwap) (*MsgCLSwapResponse, error)
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_CreateCLPool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateCLPool)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).CreateCLPool(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.clmm.Msg/CreateCLPool"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).CreateCLPool(ctx, req.(*MsgCreateCLPool)) }
	return interceptor(ctx, in, info, handler)
}

func _Msg_CreatePosition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreatePosition)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).CreatePosition(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.clmm.Msg/CreatePosition"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).CreatePosition(ctx, req.(*MsgCreatePosition)) }
	return interceptor(ctx, in, info, handler)
}

func _Msg_AddLiquidity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgAddLiquidity)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).AddLiquidity(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.clmm.Msg/AddLiquidity"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).AddLiquidity(ctx, req.(*MsgAddLiquidity)) }
	return interceptor(ctx, in, info, handler)
}

func _Msg_RemoveLiquidity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRemoveLiquidity)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).RemoveLiquidity(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.clmm.Msg/RemoveLiquidity"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).RemoveLiquidity(ctx, req.(*MsgRemoveLiquidity)) }
	return interceptor(ctx, in, info, handler)
}

func _Msg_CollectFees_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCollectFees)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).CollectFees(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.clmm.Msg/CollectFees"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).CollectFees(ctx, req.(*MsgCollectFees)) }
	return interceptor(ctx, in, info, handler)
}

func _Msg_CLSwap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCLSwap)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).CLSwap(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.clmm.Msg/CLSwap"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).CLSwap(ctx, req.(*MsgCLSwap)) }
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.clmm.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreateCLPool", Handler: _Msg_CreateCLPool_Handler},
		{MethodName: "CreatePosition", Handler: _Msg_CreatePosition_Handler},
		{MethodName: "AddLiquidity", Handler: _Msg_AddLiquidity_Handler},
		{MethodName: "RemoveLiquidity", Handler: _Msg_RemoveLiquidity_Handler},
		{MethodName: "CollectFees", Handler: _Msg_CollectFees_Handler},
		{MethodName: "CLSwap", Handler: _Msg_CLSwap_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/clmm/tx.proto",
}
