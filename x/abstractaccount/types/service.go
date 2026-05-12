package types

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

// ---------------------------------------------------------------------------
// Msg Response Types
// ---------------------------------------------------------------------------

type MsgCreateSmartAccountResponse struct {
	Address string `json:"address"`
}

func (m *MsgCreateSmartAccountResponse) ProtoMessage()  {}
func (m *MsgCreateSmartAccountResponse) Reset()         { *m = MsgCreateSmartAccountResponse{} }
func (m *MsgCreateSmartAccountResponse) String() string { return fmt.Sprintf("address=%s", m.Address) }
func (m *MsgCreateSmartAccountResponse) XXX_MessageName() string {
	return "syreen.abstractaccount.MsgCreateSmartAccountResponse"
}

type MsgCreateSessionKeyResponse struct{}

func (m *MsgCreateSessionKeyResponse) ProtoMessage()           {}
func (m *MsgCreateSessionKeyResponse) Reset()                  { *m = MsgCreateSessionKeyResponse{} }
func (m *MsgCreateSessionKeyResponse) String() string          { return "MsgCreateSessionKeyResponse" }
func (m *MsgCreateSessionKeyResponse) XXX_MessageName() string {
	return "syreen.abstractaccount.MsgCreateSessionKeyResponse"
}

type MsgRevokeSessionKeyResponse struct{}

func (m *MsgRevokeSessionKeyResponse) ProtoMessage()           {}
func (m *MsgRevokeSessionKeyResponse) Reset()                  { *m = MsgRevokeSessionKeyResponse{} }
func (m *MsgRevokeSessionKeyResponse) String() string          { return "MsgRevokeSessionKeyResponse" }
func (m *MsgRevokeSessionKeyResponse) XXX_MessageName() string {
	return "syreen.abstractaccount.MsgRevokeSessionKeyResponse"
}

type MsgInitiateRecoveryResponse struct{}

func (m *MsgInitiateRecoveryResponse) ProtoMessage()           {}
func (m *MsgInitiateRecoveryResponse) Reset()                  { *m = MsgInitiateRecoveryResponse{} }
func (m *MsgInitiateRecoveryResponse) String() string          { return "MsgInitiateRecoveryResponse" }
func (m *MsgInitiateRecoveryResponse) XXX_MessageName() string {
	return "syreen.abstractaccount.MsgInitiateRecoveryResponse"
}

type MsgApproveRecoveryResponse struct{}

func (m *MsgApproveRecoveryResponse) ProtoMessage()           {}
func (m *MsgApproveRecoveryResponse) Reset()                  { *m = MsgApproveRecoveryResponse{} }
func (m *MsgApproveRecoveryResponse) String() string          { return "MsgApproveRecoveryResponse" }
func (m *MsgApproveRecoveryResponse) XXX_MessageName() string {
	return "syreen.abstractaccount.MsgApproveRecoveryResponse"
}

type MsgExecuteRecoveryResponse struct{}

func (m *MsgExecuteRecoveryResponse) ProtoMessage()           {}
func (m *MsgExecuteRecoveryResponse) Reset()                  { *m = MsgExecuteRecoveryResponse{} }
func (m *MsgExecuteRecoveryResponse) String() string          { return "MsgExecuteRecoveryResponse" }
func (m *MsgExecuteRecoveryResponse) XXX_MessageName() string {
	return "syreen.abstractaccount.MsgExecuteRecoveryResponse"
}

type MsgSponsorGasResponse struct{}

func (m *MsgSponsorGasResponse) ProtoMessage()           {}
func (m *MsgSponsorGasResponse) Reset()                  { *m = MsgSponsorGasResponse{} }
func (m *MsgSponsorGasResponse) String() string          { return "MsgSponsorGasResponse" }
func (m *MsgSponsorGasResponse) XXX_MessageName() string {
	return "syreen.abstractaccount.MsgSponsorGasResponse"
}

type MsgBatchExecuteResponse struct{}

func (m *MsgBatchExecuteResponse) ProtoMessage()           {}
func (m *MsgBatchExecuteResponse) Reset()                  { *m = MsgBatchExecuteResponse{} }
func (m *MsgBatchExecuteResponse) String() string          { return "MsgBatchExecuteResponse" }
func (m *MsgBatchExecuteResponse) XXX_MessageName() string {
	return "syreen.abstractaccount.MsgBatchExecuteResponse"
}

// ---------------------------------------------------------------------------
// MsgServer Interface
// ---------------------------------------------------------------------------

type MsgServer interface {
	CreateSmartAccount(context.Context, *MsgCreateSmartAccount) (*MsgCreateSmartAccountResponse, error)
	CreateSessionKey(context.Context, *MsgCreateSessionKey) (*MsgCreateSessionKeyResponse, error)
	RevokeSessionKey(context.Context, *MsgRevokeSessionKey) (*MsgRevokeSessionKeyResponse, error)
	InitiateRecovery(context.Context, *MsgInitiateRecovery) (*MsgInitiateRecoveryResponse, error)
	ApproveRecovery(context.Context, *MsgApproveRecovery) (*MsgApproveRecoveryResponse, error)
	ExecuteRecovery(context.Context, *MsgExecuteRecovery) (*MsgExecuteRecoveryResponse, error)
	SponsorGas(context.Context, *MsgSponsorGas) (*MsgSponsorGasResponse, error)
	BatchExecute(context.Context, *MsgBatchExecute) (*MsgBatchExecuteResponse, error)
}

// ---------------------------------------------------------------------------
// Query Request/Response Types
// ---------------------------------------------------------------------------

type QueryParamsRequest struct{}

func (m *QueryParamsRequest) ProtoMessage()           {}
func (m *QueryParamsRequest) Reset()                  { *m = QueryParamsRequest{} }
func (m *QueryParamsRequest) String() string          { return "QueryParamsRequest" }
func (m *QueryParamsRequest) XXX_MessageName() string {
	return "syreen.abstractaccount.QueryParamsRequest"
}

type QueryParamsResponse struct {
	Params Params `json:"params"`
}

func (m *QueryParamsResponse) ProtoMessage()           {}
func (m *QueryParamsResponse) Reset()                  { *m = QueryParamsResponse{} }
func (m *QueryParamsResponse) String() string          { return fmt.Sprintf("params=%v", m.Params) }
func (m *QueryParamsResponse) XXX_MessageName() string {
	return "syreen.abstractaccount.QueryParamsResponse"
}

type QuerySmartAccountRequest struct {
	Address string `json:"address"`
}

func (m *QuerySmartAccountRequest) ProtoMessage()           {}
func (m *QuerySmartAccountRequest) Reset()                  { *m = QuerySmartAccountRequest{} }
func (m *QuerySmartAccountRequest) String() string          { return fmt.Sprintf("address=%s", m.Address) }
func (m *QuerySmartAccountRequest) XXX_MessageName() string {
	return "syreen.abstractaccount.QuerySmartAccountRequest"
}

type QuerySmartAccountResponse struct {
	Account SmartAccount `json:"account"`
}

func (m *QuerySmartAccountResponse) ProtoMessage()  {}
func (m *QuerySmartAccountResponse) Reset()         { *m = QuerySmartAccountResponse{} }
func (m *QuerySmartAccountResponse) String() string { return fmt.Sprintf("account=%v", m.Account) }
func (m *QuerySmartAccountResponse) XXX_MessageName() string {
	return "syreen.abstractaccount.QuerySmartAccountResponse"
}

type QuerySessionKeysRequest struct {
	Granter string `json:"granter"`
	Grantee string `json:"grantee"`
}

func (m *QuerySessionKeysRequest) ProtoMessage()  {}
func (m *QuerySessionKeysRequest) Reset()         { *m = QuerySessionKeysRequest{} }
func (m *QuerySessionKeysRequest) String() string { return fmt.Sprintf("granter=%s grantee=%s", m.Granter, m.Grantee) }
func (m *QuerySessionKeysRequest) XXX_MessageName() string {
	return "syreen.abstractaccount.QuerySessionKeysRequest"
}

type QuerySessionKeysResponse struct {
	SessionKey SessionKey `json:"session_key"`
}

func (m *QuerySessionKeysResponse) ProtoMessage()  {}
func (m *QuerySessionKeysResponse) Reset()         { *m = QuerySessionKeysResponse{} }
func (m *QuerySessionKeysResponse) String() string { return fmt.Sprintf("session_key=%v", m.SessionKey) }
func (m *QuerySessionKeysResponse) XXX_MessageName() string {
	return "syreen.abstractaccount.QuerySessionKeysResponse"
}

// ---------------------------------------------------------------------------
// QueryServer Interface
// ---------------------------------------------------------------------------

type QueryServer interface {
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	SmartAccount(context.Context, *QuerySmartAccountRequest) (*QuerySmartAccountResponse, error)
	SessionKeys(context.Context, *QuerySessionKeysRequest) (*QuerySessionKeysResponse, error)
}

// ---------------------------------------------------------------------------
// Msg Handlers
// ---------------------------------------------------------------------------

func _Msg_CreateSmartAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateSmartAccount)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateSmartAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syreen.abstractaccount.Msg/CreateSmartAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateSmartAccount(ctx, req.(*MsgCreateSmartAccount))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CreateSessionKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateSessionKey)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateSessionKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syreen.abstractaccount.Msg/CreateSessionKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateSessionKey(ctx, req.(*MsgCreateSessionKey))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_RevokeSessionKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRevokeSessionKey)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RevokeSessionKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syreen.abstractaccount.Msg/RevokeSessionKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RevokeSessionKey(ctx, req.(*MsgRevokeSessionKey))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_InitiateRecovery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgInitiateRecovery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).InitiateRecovery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syreen.abstractaccount.Msg/InitiateRecovery",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).InitiateRecovery(ctx, req.(*MsgInitiateRecovery))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ApproveRecovery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgApproveRecovery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ApproveRecovery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syreen.abstractaccount.Msg/ApproveRecovery",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ApproveRecovery(ctx, req.(*MsgApproveRecovery))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ExecuteRecovery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgExecuteRecovery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ExecuteRecovery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syreen.abstractaccount.Msg/ExecuteRecovery",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ExecuteRecovery(ctx, req.(*MsgExecuteRecovery))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_SponsorGas_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSponsorGas)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SponsorGas(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syreen.abstractaccount.Msg/SponsorGas",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SponsorGas(ctx, req.(*MsgSponsorGas))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_BatchExecute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgBatchExecute)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).BatchExecute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syreen.abstractaccount.Msg/BatchExecute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).BatchExecute(ctx, req.(*MsgBatchExecute))
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
		FullMethod: "/syreen.abstractaccount.Query/Params",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_SmartAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuerySmartAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).SmartAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syreen.abstractaccount.Query/SmartAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).SmartAccount(ctx, req.(*QuerySmartAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_SessionKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuerySessionKeysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).SessionKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syreen.abstractaccount.Query/SessionKeys",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).SessionKeys(ctx, req.(*QuerySessionKeysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ---------------------------------------------------------------------------
// Service Descriptors
// ---------------------------------------------------------------------------

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.abstractaccount.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateSmartAccount",
			Handler:    _Msg_CreateSmartAccount_Handler,
		},
		{
			MethodName: "CreateSessionKey",
			Handler:    _Msg_CreateSessionKey_Handler,
		},
		{
			MethodName: "RevokeSessionKey",
			Handler:    _Msg_RevokeSessionKey_Handler,
		},
		{
			MethodName: "InitiateRecovery",
			Handler:    _Msg_InitiateRecovery_Handler,
		},
		{
			MethodName: "ApproveRecovery",
			Handler:    _Msg_ApproveRecovery_Handler,
		},
		{
			MethodName: "ExecuteRecovery",
			Handler:    _Msg_ExecuteRecovery_Handler,
		},
		{
			MethodName: "SponsorGas",
			Handler:    _Msg_SponsorGas_Handler,
		},
		{
			MethodName: "BatchExecute",
			Handler:    _Msg_BatchExecute_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/abstractaccount/tx.proto",
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.abstractaccount.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "SmartAccount",
			Handler:    _Query_SmartAccount_Handler,
		},
		{
			MethodName: "SessionKeys",
			Handler:    _Query_SessionKeys_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/abstractaccount/query.proto",
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
