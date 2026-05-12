package types

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

// ---------------------------------------------------------------------------
// Msg Response Types
// ---------------------------------------------------------------------------

type MsgStoreCodeResponse struct {
	CodeID uint64 `protobuf:"varint,1,opt,name=code_id,json=codeId,proto3" json:"code_id"`
}

func (m *MsgStoreCodeResponse) ProtoMessage()           {}
func (m *MsgStoreCodeResponse) Reset()                  { *m = MsgStoreCodeResponse{} }
func (m *MsgStoreCodeResponse) String() string          { return fmt.Sprintf("code_id=%d", m.CodeID) }
func (m *MsgStoreCodeResponse) XXX_MessageName() string {
	return "syreen.compute.MsgStoreCodeResponse"
}

type MsgInstantiateContractResponse struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address"`
}

func (m *MsgInstantiateContractResponse) ProtoMessage()  {}
func (m *MsgInstantiateContractResponse) Reset()         { *m = MsgInstantiateContractResponse{} }
func (m *MsgInstantiateContractResponse) String() string { return fmt.Sprintf("address=%s", m.Address) }
func (m *MsgInstantiateContractResponse) XXX_MessageName() string {
	return "syreen.compute.MsgInstantiateContractResponse"
}

type MsgExecuteContractResponse struct {
	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *MsgExecuteContractResponse) ProtoMessage()           {}
func (m *MsgExecuteContractResponse) Reset()                  { *m = MsgExecuteContractResponse{} }
func (m *MsgExecuteContractResponse) String() string          { return fmt.Sprintf("data=%x", m.Data) }
func (m *MsgExecuteContractResponse) XXX_MessageName() string {
	return "syreen.compute.MsgExecuteContractResponse"
}

type MsgMigrateContractResponse struct {
	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *MsgMigrateContractResponse) ProtoMessage()           {}
func (m *MsgMigrateContractResponse) Reset()                  { *m = MsgMigrateContractResponse{} }
func (m *MsgMigrateContractResponse) String() string          { return fmt.Sprintf("data=%x", m.Data) }
func (m *MsgMigrateContractResponse) XXX_MessageName() string {
	return "syreen.compute.MsgMigrateContractResponse"
}

type MsgUpdateAdminResponse struct{}

func (m *MsgUpdateAdminResponse) ProtoMessage()           {}
func (m *MsgUpdateAdminResponse) Reset()                  { *m = MsgUpdateAdminResponse{} }
func (m *MsgUpdateAdminResponse) String() string          { return "MsgUpdateAdminResponse" }
func (m *MsgUpdateAdminResponse) XXX_MessageName() string {
	return "syreen.compute.MsgUpdateAdminResponse"
}

// ---------------------------------------------------------------------------
// MsgServer Interface
// ---------------------------------------------------------------------------

type MsgServer interface {
	StoreCode(context.Context, *MsgStoreCode) (*MsgStoreCodeResponse, error)
	InstantiateContract(context.Context, *MsgInstantiateContract) (*MsgInstantiateContractResponse, error)
	ExecuteContract(context.Context, *MsgExecuteContract) (*MsgExecuteContractResponse, error)
	MigrateContract(context.Context, *MsgMigrateContract) (*MsgMigrateContractResponse, error)
	UpdateAdmin(context.Context, *MsgUpdateAdmin) (*MsgUpdateAdminResponse, error)
}

// ---------------------------------------------------------------------------
// Query Request/Response Types
// ---------------------------------------------------------------------------

type QueryParamsRequest struct{}

func (m *QueryParamsRequest) ProtoMessage()           {}
func (m *QueryParamsRequest) Reset()                  { *m = QueryParamsRequest{} }
func (m *QueryParamsRequest) String() string          { return "QueryParamsRequest" }
func (m *QueryParamsRequest) XXX_MessageName() string { return "syreen.compute.QueryParamsRequest" }

type QueryParamsResponse struct {
	Params Params `json:"params"`
}

func (m *QueryParamsResponse) ProtoMessage()           {}
func (m *QueryParamsResponse) Reset()                  { *m = QueryParamsResponse{} }
func (m *QueryParamsResponse) String() string          { return fmt.Sprintf("params=%v", m.Params) }
func (m *QueryParamsResponse) XXX_MessageName() string { return "syreen.compute.QueryParamsResponse" }

type QueryCodeInfoRequest struct {
	CodeID uint64 `json:"code_id"`
}

func (m *QueryCodeInfoRequest) ProtoMessage()           {}
func (m *QueryCodeInfoRequest) Reset()                  { *m = QueryCodeInfoRequest{} }
func (m *QueryCodeInfoRequest) String() string          { return fmt.Sprintf("code_id=%d", m.CodeID) }
func (m *QueryCodeInfoRequest) XXX_MessageName() string { return "syreen.compute.QueryCodeInfoRequest" }

type QueryCodeInfoResponse struct {
	CodeInfo CodeInfo `json:"code_info"`
}

func (m *QueryCodeInfoResponse) ProtoMessage()  {}
func (m *QueryCodeInfoResponse) Reset()         { *m = QueryCodeInfoResponse{} }
func (m *QueryCodeInfoResponse) String() string { return fmt.Sprintf("code_info=%v", m.CodeInfo) }
func (m *QueryCodeInfoResponse) XXX_MessageName() string {
	return "syreen.compute.QueryCodeInfoResponse"
}

type QueryContractInfoRequest struct {
	Address string `json:"address"`
}

func (m *QueryContractInfoRequest) ProtoMessage()  {}
func (m *QueryContractInfoRequest) Reset()         { *m = QueryContractInfoRequest{} }
func (m *QueryContractInfoRequest) String() string { return fmt.Sprintf("address=%s", m.Address) }
func (m *QueryContractInfoRequest) XXX_MessageName() string {
	return "syreen.compute.QueryContractInfoRequest"
}

type QueryContractInfoResponse struct {
	ContractInfo ContractInfo `json:"contract_info"`
}

func (m *QueryContractInfoResponse) ProtoMessage()  {}
func (m *QueryContractInfoResponse) Reset()         { *m = QueryContractInfoResponse{} }
func (m *QueryContractInfoResponse) String() string { return fmt.Sprintf("contract_info=%v", m.ContractInfo) }
func (m *QueryContractInfoResponse) XXX_MessageName() string {
	return "syreen.compute.QueryContractInfoResponse"
}

type QuerySmartContractStateRequest struct {
	Address  string `json:"address"`
	QueryMsg []byte `json:"query_msg"`
}

func (m *QuerySmartContractStateRequest) ProtoMessage()  {}
func (m *QuerySmartContractStateRequest) Reset()         { *m = QuerySmartContractStateRequest{} }
func (m *QuerySmartContractStateRequest) String() string { return fmt.Sprintf("address=%s", m.Address) }
func (m *QuerySmartContractStateRequest) XXX_MessageName() string {
	return "syreen.compute.QuerySmartContractStateRequest"
}

type QuerySmartContractStateResponse struct {
	Data []byte `json:"data"`
}

func (m *QuerySmartContractStateResponse) ProtoMessage()  {}
func (m *QuerySmartContractStateResponse) Reset()         { *m = QuerySmartContractStateResponse{} }
func (m *QuerySmartContractStateResponse) String() string { return fmt.Sprintf("data=%x", m.Data) }
func (m *QuerySmartContractStateResponse) XXX_MessageName() string {
	return "syreen.compute.QuerySmartContractStateResponse"
}

type QueryRawContractStateRequest struct {
	Address string `json:"address"`
	Key     []byte `json:"key"`
}

func (m *QueryRawContractStateRequest) ProtoMessage()  {}
func (m *QueryRawContractStateRequest) Reset()         { *m = QueryRawContractStateRequest{} }
func (m *QueryRawContractStateRequest) String() string { return fmt.Sprintf("address=%s", m.Address) }
func (m *QueryRawContractStateRequest) XXX_MessageName() string {
	return "syreen.compute.QueryRawContractStateRequest"
}

type QueryRawContractStateResponse struct {
	Data []byte `json:"data"`
}

func (m *QueryRawContractStateResponse) ProtoMessage()  {}
func (m *QueryRawContractStateResponse) Reset()         { *m = QueryRawContractStateResponse{} }
func (m *QueryRawContractStateResponse) String() string { return fmt.Sprintf("data=%x", m.Data) }
func (m *QueryRawContractStateResponse) XXX_MessageName() string {
	return "syreen.compute.QueryRawContractStateResponse"
}

// ---------------------------------------------------------------------------
// QueryServer Interface
// ---------------------------------------------------------------------------

type QueryServer interface {
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	CodeInfo(context.Context, *QueryCodeInfoRequest) (*QueryCodeInfoResponse, error)
	ContractInfo(context.Context, *QueryContractInfoRequest) (*QueryContractInfoResponse, error)
	SmartContractState(context.Context, *QuerySmartContractStateRequest) (*QuerySmartContractStateResponse, error)
	RawContractState(context.Context, *QueryRawContractStateRequest) (*QueryRawContractStateResponse, error)
}

// ---------------------------------------------------------------------------
// Msg Handlers
// ---------------------------------------------------------------------------

func _Msg_StoreCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgStoreCode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).StoreCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syreen.compute.Msg/StoreCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).StoreCode(ctx, req.(*MsgStoreCode))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_InstantiateContract_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgInstantiateContract)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).InstantiateContract(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syreen.compute.Msg/InstantiateContract",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).InstantiateContract(ctx, req.(*MsgInstantiateContract))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ExecuteContract_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgExecuteContract)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ExecuteContract(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syreen.compute.Msg/ExecuteContract",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ExecuteContract(ctx, req.(*MsgExecuteContract))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_MigrateContract_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgMigrateContract)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).MigrateContract(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syreen.compute.Msg/MigrateContract",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).MigrateContract(ctx, req.(*MsgMigrateContract))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateAdmin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateAdmin)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateAdmin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syreen.compute.Msg/UpdateAdmin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateAdmin(ctx, req.(*MsgUpdateAdmin))
	}
	return interceptor(ctx, in, info, handler)
}

// ---------------------------------------------------------------------------
// Query Handlers
// ---------------------------------------------------------------------------

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syreen.compute.Query/Params",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_CodeInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryCodeInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).CodeInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syreen.compute.Query/CodeInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).CodeInfo(ctx, req.(*QueryCodeInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ContractInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryContractInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ContractInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syreen.compute.Query/ContractInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ContractInfo(ctx, req.(*QueryContractInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_SmartContractState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuerySmartContractStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).SmartContractState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syreen.compute.Query/SmartContractState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).SmartContractState(ctx, req.(*QuerySmartContractStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_RawContractState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRawContractStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).RawContractState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syreen.compute.Query/RawContractState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).RawContractState(ctx, req.(*QueryRawContractStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ---------------------------------------------------------------------------
// Service Descriptors
// ---------------------------------------------------------------------------

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.compute.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StoreCode",
			Handler:    _Msg_StoreCode_Handler,
		},
		{
			MethodName: "InstantiateContract",
			Handler:    _Msg_InstantiateContract_Handler,
		},
		{
			MethodName: "ExecuteContract",
			Handler:    _Msg_ExecuteContract_Handler,
		},
		{
			MethodName: "MigrateContract",
			Handler:    _Msg_MigrateContract_Handler,
		},
		{
			MethodName: "UpdateAdmin",
			Handler:    _Msg_UpdateAdmin_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/compute/tx.proto",
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.compute.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "CodeInfo",
			Handler:    _Query_CodeInfo_Handler,
		},
		{
			MethodName: "ContractInfo",
			Handler:    _Query_ContractInfo_Handler,
		},
		{
			MethodName: "SmartContractState",
			Handler:    _Query_SmartContractState_Handler,
		},
		{
			MethodName: "RawContractState",
			Handler:    _Query_RawContractState_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/compute/query.proto",
}

// ---------------------------------------------------------------------------
// Registration Functions
// ---------------------------------------------------------------------------

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}
