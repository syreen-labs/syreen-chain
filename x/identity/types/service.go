package types

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

// --- Msg Response Types ---

type MsgRegisterIdentityResponse struct{}
func (m *MsgRegisterIdentityResponse) ProtoMessage()           {}
func (m *MsgRegisterIdentityResponse) Reset()                  { *m = MsgRegisterIdentityResponse{} }
func (m *MsgRegisterIdentityResponse) String() string          { return "register_identity_response" }
func (m *MsgRegisterIdentityResponse) XXX_MessageName() string { return "syreen.identity.MsgRegisterIdentityResponse" }

type MsgVerifyIdentityResponse struct{}
func (m *MsgVerifyIdentityResponse) ProtoMessage()           {}
func (m *MsgVerifyIdentityResponse) Reset()                  { *m = MsgVerifyIdentityResponse{} }
func (m *MsgVerifyIdentityResponse) String() string          { return "verify_identity_response" }
func (m *MsgVerifyIdentityResponse) XXX_MessageName() string { return "syreen.identity.MsgVerifyIdentityResponse" }

type MsgRejectIdentityResponse struct{}
func (m *MsgRejectIdentityResponse) ProtoMessage()           {}
func (m *MsgRejectIdentityResponse) Reset()                  { *m = MsgRejectIdentityResponse{} }
func (m *MsgRejectIdentityResponse) String() string          { return "reject_identity_response" }
func (m *MsgRejectIdentityResponse) XXX_MessageName() string { return "syreen.identity.MsgRejectIdentityResponse" }

type MsgRevokeIdentityResponse struct{}
func (m *MsgRevokeIdentityResponse) ProtoMessage()           {}
func (m *MsgRevokeIdentityResponse) Reset()                  { *m = MsgRevokeIdentityResponse{} }
func (m *MsgRevokeIdentityResponse) String() string          { return "revoke_identity_response" }
func (m *MsgRevokeIdentityResponse) XXX_MessageName() string { return "syreen.identity.MsgRevokeIdentityResponse" }

type MsgUpdateIdentityResponse struct{}
func (m *MsgUpdateIdentityResponse) ProtoMessage()           {}
func (m *MsgUpdateIdentityResponse) Reset()                  { *m = MsgUpdateIdentityResponse{} }
func (m *MsgUpdateIdentityResponse) String() string          { return "update_identity_response" }
func (m *MsgUpdateIdentityResponse) XXX_MessageName() string { return "syreen.identity.MsgUpdateIdentityResponse" }

type MsgRegisterVerifierResponse struct{}
func (m *MsgRegisterVerifierResponse) ProtoMessage()           {}
func (m *MsgRegisterVerifierResponse) Reset()                  { *m = MsgRegisterVerifierResponse{} }
func (m *MsgRegisterVerifierResponse) String() string          { return "register_verifier_response" }
func (m *MsgRegisterVerifierResponse) XXX_MessageName() string { return "syreen.identity.MsgRegisterVerifierResponse" }

type MsgDeactivateVerifierResponse struct{}
func (m *MsgDeactivateVerifierResponse) ProtoMessage()           {}
func (m *MsgDeactivateVerifierResponse) Reset()                  { *m = MsgDeactivateVerifierResponse{} }
func (m *MsgDeactivateVerifierResponse) String() string          { return "deactivate_verifier_response" }
func (m *MsgDeactivateVerifierResponse) XXX_MessageName() string { return "syreen.identity.MsgDeactivateVerifierResponse" }

type MsgIncrementBookingsResponse struct{}
func (m *MsgIncrementBookingsResponse) ProtoMessage()           {}
func (m *MsgIncrementBookingsResponse) Reset()                  { *m = MsgIncrementBookingsResponse{} }
func (m *MsgIncrementBookingsResponse) String() string          { return "increment_bookings_response" }
func (m *MsgIncrementBookingsResponse) XXX_MessageName() string { return "syreen.identity.MsgIncrementBookingsResponse" }

// --- Query Types ---

type QueryParamsRequest struct{}
func (m *QueryParamsRequest) ProtoMessage()           {}
func (m *QueryParamsRequest) Reset()                  { *m = QueryParamsRequest{} }
func (m *QueryParamsRequest) String() string          { return "query_params_request" }
func (m *QueryParamsRequest) XXX_MessageName() string { return "syreen.identity.QueryParamsRequest" }

type QueryParamsResponse struct {
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}
func (m *QueryParamsResponse) ProtoMessage()           {}
func (m *QueryParamsResponse) Reset()                  { *m = QueryParamsResponse{} }
func (m *QueryParamsResponse) String() string          { return fmt.Sprintf("params: %+v", m.Params) }
func (m *QueryParamsResponse) XXX_MessageName() string { return "syreen.identity.QueryParamsResponse" }

type QueryIdentityRequest struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address"`
}
func (m *QueryIdentityRequest) ProtoMessage()           {}
func (m *QueryIdentityRequest) Reset()                  { *m = QueryIdentityRequest{} }
func (m *QueryIdentityRequest) String() string          { return fmt.Sprintf("query_identity: %s", m.Address) }
func (m *QueryIdentityRequest) XXX_MessageName() string { return "syreen.identity.QueryIdentityRequest" }

type QueryIdentityResponse struct {
	Identity Identity `protobuf:"bytes,1,opt,name=identity,proto3" json:"identity"`
	Found    bool     `protobuf:"varint,2,opt,name=found,proto3" json:"found"`
}
func (m *QueryIdentityResponse) ProtoMessage()           {}
func (m *QueryIdentityResponse) Reset()                  { *m = QueryIdentityResponse{} }
func (m *QueryIdentityResponse) String() string          { return fmt.Sprintf("identity: %+v", m.Identity) }
func (m *QueryIdentityResponse) XXX_MessageName() string { return "syreen.identity.QueryIdentityResponse" }

type QueryIdentitiesByVerifierRequest struct {
	Verifier string `protobuf:"bytes,1,opt,name=verifier,proto3" json:"verifier"`
}
func (m *QueryIdentitiesByVerifierRequest) ProtoMessage()           {}
func (m *QueryIdentitiesByVerifierRequest) Reset()                  { *m = QueryIdentitiesByVerifierRequest{} }
func (m *QueryIdentitiesByVerifierRequest) String() string          { return fmt.Sprintf("query_by_verifier: %s", m.Verifier) }
func (m *QueryIdentitiesByVerifierRequest) XXX_MessageName() string { return "syreen.identity.QueryIdentitiesByVerifierRequest" }

type QueryIdentitiesByVerifierResponse struct {
	Identities []Identity `protobuf:"bytes,1,rep,name=identities,proto3" json:"identities"`
}
func (m *QueryIdentitiesByVerifierResponse) ProtoMessage()           {}
func (m *QueryIdentitiesByVerifierResponse) Reset()                  { *m = QueryIdentitiesByVerifierResponse{} }
func (m *QueryIdentitiesByVerifierResponse) String() string          { return fmt.Sprintf("identities: %d", len(m.Identities)) }
func (m *QueryIdentitiesByVerifierResponse) XXX_MessageName() string { return "syreen.identity.QueryIdentitiesByVerifierResponse" }

type QueryVerifierRequest struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address"`
}
func (m *QueryVerifierRequest) ProtoMessage()           {}
func (m *QueryVerifierRequest) Reset()                  { *m = QueryVerifierRequest{} }
func (m *QueryVerifierRequest) String() string          { return fmt.Sprintf("query_verifier: %s", m.Address) }
func (m *QueryVerifierRequest) XXX_MessageName() string { return "syreen.identity.QueryVerifierRequest" }

type QueryVerifierResponse struct {
	Verifier Verifier `protobuf:"bytes,1,opt,name=verifier,proto3" json:"verifier"`
	Found    bool     `protobuf:"varint,2,opt,name=found,proto3" json:"found"`
}
func (m *QueryVerifierResponse) ProtoMessage()           {}
func (m *QueryVerifierResponse) Reset()                  { *m = QueryVerifierResponse{} }
func (m *QueryVerifierResponse) String() string          { return fmt.Sprintf("verifier: %+v", m.Verifier) }
func (m *QueryVerifierResponse) XXX_MessageName() string { return "syreen.identity.QueryVerifierResponse" }

type QueryVerificationStatusRequest struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address"`
}
func (m *QueryVerificationStatusRequest) ProtoMessage()           {}
func (m *QueryVerificationStatusRequest) Reset()                  { *m = QueryVerificationStatusRequest{} }
func (m *QueryVerificationStatusRequest) String() string          { return fmt.Sprintf("query_status: %s", m.Address) }
func (m *QueryVerificationStatusRequest) XXX_MessageName() string { return "syreen.identity.QueryVerificationStatusRequest" }

type QueryVerificationStatusResponse struct {
	Address  string             `protobuf:"bytes,1,opt,name=address,proto3" json:"address"`
	Verified bool               `protobuf:"varint,2,opt,name=verified,proto3" json:"verified"`
	Level    VerificationLevel  `protobuf:"bytes,3,opt,name=level,proto3,casttype=VerificationLevel" json:"level"`
	Status   VerificationStatus `protobuf:"bytes,4,opt,name=status,proto3,casttype=VerificationStatus" json:"status"`
	TrustScore uint64           `protobuf:"varint,5,opt,name=trust_score,json=trustScore,proto3" json:"trust_score"`
}
func (m *QueryVerificationStatusResponse) ProtoMessage()           {}
func (m *QueryVerificationStatusResponse) Reset()                  { *m = QueryVerificationStatusResponse{} }
func (m *QueryVerificationStatusResponse) String() string          { return fmt.Sprintf("status: %s %s %s", m.Address, m.Level, m.Status) }
func (m *QueryVerificationStatusResponse) XXX_MessageName() string { return "syreen.identity.QueryVerificationStatusResponse" }

// --- Server Interfaces ---

type MsgServer interface {
	RegisterIdentity(context.Context, *MsgRegisterIdentity) (*MsgRegisterIdentityResponse, error)
	VerifyIdentity(context.Context, *MsgVerifyIdentity) (*MsgVerifyIdentityResponse, error)
	RejectIdentity(context.Context, *MsgRejectIdentity) (*MsgRejectIdentityResponse, error)
	RevokeIdentity(context.Context, *MsgRevokeIdentity) (*MsgRevokeIdentityResponse, error)
	UpdateIdentity(context.Context, *MsgUpdateIdentity) (*MsgUpdateIdentityResponse, error)
	RegisterVerifier(context.Context, *MsgRegisterVerifier) (*MsgRegisterVerifierResponse, error)
	DeactivateVerifier(context.Context, *MsgDeactivateVerifier) (*MsgDeactivateVerifierResponse, error)
	IncrementBookings(context.Context, *MsgIncrementBookings) (*MsgIncrementBookingsResponse, error)
}

type QueryServer interface {
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	Identity(context.Context, *QueryIdentityRequest) (*QueryIdentityResponse, error)
	IdentitiesByVerifier(context.Context, *QueryIdentitiesByVerifierRequest) (*QueryIdentitiesByVerifierResponse, error)
	Verifier(context.Context, *QueryVerifierRequest) (*QueryVerifierResponse, error)
	VerificationStatus(context.Context, *QueryVerificationStatusRequest) (*QueryVerificationStatusResponse, error)
}

// --- gRPC Handlers ---

func _Msg_RegisterIdentity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRegisterIdentity)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).RegisterIdentity(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.identity.Msg/RegisterIdentity"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).RegisterIdentity(ctx, req.(*MsgRegisterIdentity)) })
}

func _Msg_VerifyIdentity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgVerifyIdentity)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).VerifyIdentity(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.identity.Msg/VerifyIdentity"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).VerifyIdentity(ctx, req.(*MsgVerifyIdentity)) })
}

func _Msg_RejectIdentity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRejectIdentity)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).RejectIdentity(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.identity.Msg/RejectIdentity"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).RejectIdentity(ctx, req.(*MsgRejectIdentity)) })
}

func _Msg_RevokeIdentity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRevokeIdentity)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).RevokeIdentity(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.identity.Msg/RevokeIdentity"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).RevokeIdentity(ctx, req.(*MsgRevokeIdentity)) })
}

func _Msg_UpdateIdentity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateIdentity)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).UpdateIdentity(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.identity.Msg/UpdateIdentity"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).UpdateIdentity(ctx, req.(*MsgUpdateIdentity)) })
}

func _Msg_RegisterVerifier_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRegisterVerifier)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).RegisterVerifier(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.identity.Msg/RegisterVerifier"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).RegisterVerifier(ctx, req.(*MsgRegisterVerifier)) })
}

func _Msg_DeactivateVerifier_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgDeactivateVerifier)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).DeactivateVerifier(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.identity.Msg/DeactivateVerifier"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).DeactivateVerifier(ctx, req.(*MsgDeactivateVerifier)) })
}

func _Msg_IncrementBookings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgIncrementBookings)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).IncrementBookings(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.identity.Msg/IncrementBookings"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).IncrementBookings(ctx, req.(*MsgIncrementBookings)) })
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(QueryServer).Params(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.identity.Query/Params"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest)) })
}

func _Query_Identity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryIdentityRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(QueryServer).Identity(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.identity.Query/Identity"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(QueryServer).Identity(ctx, req.(*QueryIdentityRequest)) })
}

func _Query_IdentitiesByVerifier_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryIdentitiesByVerifierRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(QueryServer).IdentitiesByVerifier(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.identity.Query/IdentitiesByVerifier"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(QueryServer).IdentitiesByVerifier(ctx, req.(*QueryIdentitiesByVerifierRequest)) })
}

func _Query_Verifier_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryVerifierRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(QueryServer).Verifier(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.identity.Query/Verifier"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(QueryServer).Verifier(ctx, req.(*QueryVerifierRequest)) })
}

func _Query_VerificationStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryVerificationStatusRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(QueryServer).VerificationStatus(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.identity.Query/VerificationStatus"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(QueryServer).VerificationStatus(ctx, req.(*QueryVerificationStatusRequest)) })
}

// --- Service Descriptors ---

var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.identity.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "RegisterIdentity", Handler: _Msg_RegisterIdentity_Handler},
		{MethodName: "VerifyIdentity", Handler: _Msg_VerifyIdentity_Handler},
		{MethodName: "RejectIdentity", Handler: _Msg_RejectIdentity_Handler},
		{MethodName: "RevokeIdentity", Handler: _Msg_RevokeIdentity_Handler},
		{MethodName: "UpdateIdentity", Handler: _Msg_UpdateIdentity_Handler},
		{MethodName: "RegisterVerifier", Handler: _Msg_RegisterVerifier_Handler},
		{MethodName: "DeactivateVerifier", Handler: _Msg_DeactivateVerifier_Handler},
		{MethodName: "IncrementBookings", Handler: _Msg_IncrementBookings_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/identity/tx.proto",
}

var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.identity.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Params", Handler: _Query_Params_Handler},
		{MethodName: "Identity", Handler: _Query_Identity_Handler},
		{MethodName: "IdentitiesByVerifier", Handler: _Query_IdentitiesByVerifier_Handler},
		{MethodName: "Verifier", Handler: _Query_Verifier_Handler},
		{MethodName: "VerificationStatus", Handler: _Query_VerificationStatus_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/identity/query.proto",
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
}
