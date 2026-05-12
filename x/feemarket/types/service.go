package types

import (
	"context"
	"fmt"

	"cosmossdk.io/math"

	"google.golang.org/grpc"
)

// ---------------------------------------------------------------------------
// Query request/response types
// ---------------------------------------------------------------------------

type QueryParamsRequest struct{}

func (m *QueryParamsRequest) ProtoMessage()           {}
func (m *QueryParamsRequest) Reset()                  { *m = QueryParamsRequest{} }
func (m *QueryParamsRequest) String() string          { return "query_params_request" }
func (m *QueryParamsRequest) XXX_MessageName() string { return "syreen.feemarket.QueryParamsRequest" }

type QueryParamsResponse struct {
	Params Params `json:"params"`
}

func (m *QueryParamsResponse) ProtoMessage()           {}
func (m *QueryParamsResponse) Reset()                  { *m = QueryParamsResponse{} }
func (m *QueryParamsResponse) String() string          { return fmt.Sprintf("params: %+v", m.Params) }
func (m *QueryParamsResponse) XXX_MessageName() string { return "syreen.feemarket.QueryParamsResponse" }

type QueryBaseFeeRequest struct{}

func (m *QueryBaseFeeRequest) ProtoMessage()           {}
func (m *QueryBaseFeeRequest) Reset()                  { *m = QueryBaseFeeRequest{} }
func (m *QueryBaseFeeRequest) String() string          { return "query_base_fee_request" }
func (m *QueryBaseFeeRequest) XXX_MessageName() string { return "syreen.feemarket.QueryBaseFeeRequest" }

type QueryBaseFeeResponse struct {
	BaseFee math.LegacyDec `json:"base_fee"`
}

func (m *QueryBaseFeeResponse) ProtoMessage()           {}
func (m *QueryBaseFeeResponse) Reset()                  { *m = QueryBaseFeeResponse{} }
func (m *QueryBaseFeeResponse) String() string          { return fmt.Sprintf("base_fee: %s", m.BaseFee) }
func (m *QueryBaseFeeResponse) XXX_MessageName() string { return "syreen.feemarket.QueryBaseFeeResponse" }

type QueryFeeStateRequest struct{}

func (m *QueryFeeStateRequest) ProtoMessage()           {}
func (m *QueryFeeStateRequest) Reset()                  { *m = QueryFeeStateRequest{} }
func (m *QueryFeeStateRequest) String() string          { return "query_fee_state_request" }
func (m *QueryFeeStateRequest) XXX_MessageName() string { return "syreen.feemarket.QueryFeeStateRequest" }

type QueryFeeStateResponse struct {
	FeeState FeeState `json:"fee_state"`
}

func (m *QueryFeeStateResponse) ProtoMessage()           {}
func (m *QueryFeeStateResponse) Reset()                  { *m = QueryFeeStateResponse{} }
func (m *QueryFeeStateResponse) String() string          { return fmt.Sprintf("fee_state: %+v", m.FeeState) }
func (m *QueryFeeStateResponse) XXX_MessageName() string { return "syreen.feemarket.QueryFeeStateResponse" }

// ---------------------------------------------------------------------------
// QueryServer interface
// ---------------------------------------------------------------------------

type QueryServer interface {
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	BaseFee(context.Context, *QueryBaseFeeRequest) (*QueryBaseFeeResponse, error)
	FeeState(context.Context, *QueryFeeStateRequest) (*QueryFeeStateResponse, error)
}

// ---------------------------------------------------------------------------
// RegisterQueryServer
// ---------------------------------------------------------------------------

// RegisterQueryServer registers the QueryServer with the gRPC service registrar.
func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
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
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.feemarket.Query/Params"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_BaseFee_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryBaseFeeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).BaseFee(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.feemarket.Query/BaseFee"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).BaseFee(ctx, req.(*QueryBaseFeeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_FeeState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryFeeStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).FeeState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.feemarket.Query/FeeState"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).FeeState(ctx, req.(*QueryFeeStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.feemarket.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Params", Handler: _Query_Params_Handler},
		{MethodName: "BaseFee", Handler: _Query_BaseFee_Handler},
		{MethodName: "FeeState", Handler: _Query_FeeState_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/feemarket/query.proto",
}
