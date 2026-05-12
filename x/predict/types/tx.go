package types

import (
	"context"

	"google.golang.org/grpc"
)

// MsgServer defines the predict message server interface
type MsgServer interface {
	CreateMarket(context.Context, *MsgCreateMarket) (*MsgCreateMarketResponse, error)
	BuyShares(context.Context, *MsgBuyShares) (*MsgBuySharesResponse, error)
	SellShares(context.Context, *MsgSellShares) (*MsgSellSharesResponse, error)
	ResolveMarket(context.Context, *MsgResolveMarket) (*MsgResolveMarketResponse, error)
	ClaimWinnings(context.Context, *MsgClaimWinnings) (*MsgClaimWinningsResponse, error)
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

// Handlers

func _Msg_CreateMarket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateMarket)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateMarket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.predict.Msg/CreateMarket"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateMarket(ctx, req.(*MsgCreateMarket))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_BuyShares_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgBuyShares)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).BuyShares(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.predict.Msg/BuyShares"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).BuyShares(ctx, req.(*MsgBuyShares))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_SellShares_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSellShares)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SellShares(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.predict.Msg/SellShares"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SellShares(ctx, req.(*MsgSellShares))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ResolveMarket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgResolveMarket)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ResolveMarket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.predict.Msg/ResolveMarket"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ResolveMarket(ctx, req.(*MsgResolveMarket))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ClaimWinnings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgClaimWinnings)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ClaimWinnings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.predict.Msg/ClaimWinnings"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ClaimWinnings(ctx, req.(*MsgClaimWinnings))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.predict.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreateMarket", Handler: _Msg_CreateMarket_Handler},
		{MethodName: "BuyShares", Handler: _Msg_BuyShares_Handler},
		{MethodName: "SellShares", Handler: _Msg_SellShares_Handler},
		{MethodName: "ResolveMarket", Handler: _Msg_ResolveMarket_Handler},
		{MethodName: "ClaimWinnings", Handler: _Msg_ClaimWinnings_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/predict/tx.proto",
}
