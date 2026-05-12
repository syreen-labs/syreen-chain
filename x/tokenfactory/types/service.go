package types

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

// ---------------------------------------------------------------------------
// Query request/response types
// ---------------------------------------------------------------------------

type QueryParamsRequest struct{}

func (m *QueryParamsRequest) ProtoMessage()           {}
func (m *QueryParamsRequest) Reset()                  { *m = QueryParamsRequest{} }
func (m *QueryParamsRequest) String() string          { return "query_params_request" }
func (m *QueryParamsRequest) XXX_MessageName() string { return "syreen.tokenfactory.QueryParamsRequest" }

type QueryParamsResponse struct {
	Params Params `json:"params"`
}

func (m *QueryParamsResponse) ProtoMessage()           {}
func (m *QueryParamsResponse) Reset()                  { *m = QueryParamsResponse{} }
func (m *QueryParamsResponse) String() string          { return fmt.Sprintf("params: %+v", m.Params) }
func (m *QueryParamsResponse) XXX_MessageName() string { return "syreen.tokenfactory.QueryParamsResponse" }

type QueryDenomAuthorityRequest struct {
	Denom string `json:"denom"`
}

func (m *QueryDenomAuthorityRequest) ProtoMessage()           {}
func (m *QueryDenomAuthorityRequest) Reset()                  { *m = QueryDenomAuthorityRequest{} }
func (m *QueryDenomAuthorityRequest) String() string          { return fmt.Sprintf("query_denom_authority: %s", m.Denom) }
func (m *QueryDenomAuthorityRequest) XXX_MessageName() string { return "syreen.tokenfactory.QueryDenomAuthorityRequest" }

type QueryDenomAuthorityResponse struct {
	Admin string `json:"admin"`
}

func (m *QueryDenomAuthorityResponse) ProtoMessage()           {}
func (m *QueryDenomAuthorityResponse) Reset()                  { *m = QueryDenomAuthorityResponse{} }
func (m *QueryDenomAuthorityResponse) String() string          { return fmt.Sprintf("admin: %s", m.Admin) }
func (m *QueryDenomAuthorityResponse) XXX_MessageName() string { return "syreen.tokenfactory.QueryDenomAuthorityResponse" }

// ---------------------------------------------------------------------------
// QueryServer interface
// ---------------------------------------------------------------------------

type QueryServer interface {
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	DenomAuthority(context.Context, *QueryDenomAuthorityRequest) (*QueryDenomAuthorityResponse, error)
}

// ---------------------------------------------------------------------------
// RegisterMsgServer
// ---------------------------------------------------------------------------

// RegisterMsgServer registers the MsgServer with the gRPC service registrar.
func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

// ---------------------------------------------------------------------------
// RegisterQueryServer
// ---------------------------------------------------------------------------

// RegisterQueryServer registers the QueryServer with the gRPC service registrar.
func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

// ---------------------------------------------------------------------------
// Msg service handlers
// ---------------------------------------------------------------------------

func _Msg_CreateDenom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateDenom)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateDenom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.tokenfactory.Msg/CreateDenom"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateDenom(ctx, req.(*MsgCreateDenom))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Mint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgMint)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Mint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.tokenfactory.Msg/Mint"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Mint(ctx, req.(*MsgMint))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Burn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgBurn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Burn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.tokenfactory.Msg/Burn"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Burn(ctx, req.(*MsgBurn))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ChangeAdmin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgChangeAdmin)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ChangeAdmin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.tokenfactory.Msg/ChangeAdmin"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ChangeAdmin(ctx, req.(*MsgChangeAdmin))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.tokenfactory.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreateDenom", Handler: _Msg_CreateDenom_Handler},
		{MethodName: "Mint", Handler: _Msg_Mint_Handler},
		{MethodName: "Burn", Handler: _Msg_Burn_Handler},
		{MethodName: "ChangeAdmin", Handler: _Msg_ChangeAdmin_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/tokenfactory/tx.proto",
}

// ---------------------------------------------------------------------------
// Query service handlers
// ---------------------------------------------------------------------------

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.tokenfactory.Query/Params"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_DenomAuthority_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryDenomAuthorityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).DenomAuthority(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.tokenfactory.Query/DenomAuthority"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).DenomAuthority(ctx, req.(*QueryDenomAuthorityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.tokenfactory.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Params", Handler: _Query_Params_Handler},
		{MethodName: "DenomAuthority", Handler: _Query_DenomAuthority_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/tokenfactory/query.proto",
}
