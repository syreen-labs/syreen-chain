package types

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

// --- Msg Response Types ---

type MsgCommitTxResponse struct{}

func (m *MsgCommitTxResponse) ProtoMessage()           {}
func (m *MsgCommitTxResponse) Reset()                  { *m = MsgCommitTxResponse{} }
func (m *MsgCommitTxResponse) String() string          { return "commit_tx_response" }
func (m *MsgCommitTxResponse) XXX_MessageName() string { return "syreen.mevprotection.MsgCommitTxResponse" }

type MsgRevealTxResponse struct{}

func (m *MsgRevealTxResponse) ProtoMessage()           {}
func (m *MsgRevealTxResponse) Reset()                  { *m = MsgRevealTxResponse{} }
func (m *MsgRevealTxResponse) String() string          { return "reveal_tx_response" }
func (m *MsgRevealTxResponse) XXX_MessageName() string { return "syreen.mevprotection.MsgRevealTxResponse" }

// --- MsgServer Interface ---

type MsgServer interface {
	CommitTx(context.Context, *MsgCommitTx) (*MsgCommitTxResponse, error)
	RevealTx(context.Context, *MsgRevealTx) (*MsgRevealTxResponse, error)
}

// --- Query Request/Response Types ---

type QueryParamsRequest struct{}

func (m *QueryParamsRequest) ProtoMessage()           {}
func (m *QueryParamsRequest) Reset()                  { *m = QueryParamsRequest{} }
func (m *QueryParamsRequest) String() string          { return "query_params_request" }
func (m *QueryParamsRequest) XXX_MessageName() string { return "syreen.mevprotection.QueryParamsRequest" }

type QueryParamsResponse struct {
	Params Params `json:"params"`
}

func (m *QueryParamsResponse) ProtoMessage()           {}
func (m *QueryParamsResponse) Reset()                  { *m = QueryParamsResponse{} }
func (m *QueryParamsResponse) String() string          { return fmt.Sprintf("params: %+v", m.Params) }
func (m *QueryParamsResponse) XXX_MessageName() string { return "syreen.mevprotection.QueryParamsResponse" }

type QueryCommittedTxRequest struct {
	TxHash []byte `json:"tx_hash"`
}

func (m *QueryCommittedTxRequest) ProtoMessage()           {}
func (m *QueryCommittedTxRequest) Reset()                  { *m = QueryCommittedTxRequest{} }
func (m *QueryCommittedTxRequest) String() string          { return fmt.Sprintf("query_committed_tx: hash=%x", m.TxHash) }
func (m *QueryCommittedTxRequest) XXX_MessageName() string { return "syreen.mevprotection.QueryCommittedTxRequest" }

type QueryCommittedTxResponse struct {
	CommittedTx CommittedTx `json:"committed_tx"`
	Found       bool        `json:"found"`
}

func (m *QueryCommittedTxResponse) ProtoMessage()           {}
func (m *QueryCommittedTxResponse) Reset()                  { *m = QueryCommittedTxResponse{} }
func (m *QueryCommittedTxResponse) String() string          { return fmt.Sprintf("committed_tx: %+v", m.CommittedTx) }
func (m *QueryCommittedTxResponse) XXX_MessageName() string { return "syreen.mevprotection.QueryCommittedTxResponse" }

type QueryPendingRevealsRequest struct{}

func (m *QueryPendingRevealsRequest) ProtoMessage()           {}
func (m *QueryPendingRevealsRequest) Reset()                  { *m = QueryPendingRevealsRequest{} }
func (m *QueryPendingRevealsRequest) String() string          { return "query_pending_reveals_request" }
func (m *QueryPendingRevealsRequest) XXX_MessageName() string { return "syreen.mevprotection.QueryPendingRevealsRequest" }

type QueryPendingRevealsResponse struct {
	PendingReveals []CommittedTx `json:"pending_reveals"`
}

func (m *QueryPendingRevealsResponse) ProtoMessage()           {}
func (m *QueryPendingRevealsResponse) Reset()                  { *m = QueryPendingRevealsResponse{} }
func (m *QueryPendingRevealsResponse) String() string          { return fmt.Sprintf("pending_reveals: %d", len(m.PendingReveals)) }
func (m *QueryPendingRevealsResponse) XXX_MessageName() string { return "syreen.mevprotection.QueryPendingRevealsResponse" }

// --- MEV Redistribution Query Types ---

type QueryMEVRedistributionRequest struct{}

func (m *QueryMEVRedistributionRequest) ProtoMessage()           {}
func (m *QueryMEVRedistributionRequest) Reset()                  { *m = QueryMEVRedistributionRequest{} }
func (m *QueryMEVRedistributionRequest) String() string          { return "query_mev_redistribution_request" }
func (m *QueryMEVRedistributionRequest) XXX_MessageName() string { return "syreen.mevprotection.QueryMEVRedistributionRequest" }

type QueryMEVRedistributionResponse struct {
	RewardPool         string `json:"reward_pool"`
	TotalRedistributed string `json:"total_redistributed"`
	LPSharePct         string `json:"lp_share_pct"`
	StakerSharePct     string `json:"staker_share_pct"`
	DistributionInterval int64 `json:"distribution_interval"`
}

func (m *QueryMEVRedistributionResponse) ProtoMessage()           {}
func (m *QueryMEVRedistributionResponse) Reset()                  { *m = QueryMEVRedistributionResponse{} }
func (m *QueryMEVRedistributionResponse) String() string          { return fmt.Sprintf("mev_redistribution: pool=%s total=%s", m.RewardPool, m.TotalRedistributed) }
func (m *QueryMEVRedistributionResponse) XXX_MessageName() string { return "syreen.mevprotection.QueryMEVRedistributionResponse" }

// --- QueryServer Interface ---

type QueryServer interface {
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	CommittedTx(context.Context, *QueryCommittedTxRequest) (*QueryCommittedTxResponse, error)
	PendingReveals(context.Context, *QueryPendingRevealsRequest) (*QueryPendingRevealsResponse, error)
	MEVRedistribution(context.Context, *QueryMEVRedistributionRequest) (*QueryMEVRedistributionResponse, error)
}

// --- Msg Handler Functions ---

func _Msg_CommitTx_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCommitTx)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CommitTx(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.mevprotection.Msg/CommitTx"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CommitTx(ctx, req.(*MsgCommitTx))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_RevealTx_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRevealTx)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RevealTx(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.mevprotection.Msg/RevealTx"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RevealTx(ctx, req.(*MsgRevealTx))
	}
	return interceptor(ctx, in, info, handler)
}

// --- Query Handler Functions ---

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.mevprotection.Query/Params"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_CommittedTx_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryCommittedTxRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).CommittedTx(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.mevprotection.Query/CommittedTx"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).CommittedTx(ctx, req.(*QueryCommittedTxRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_PendingReveals_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryPendingRevealsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).PendingReveals(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.mevprotection.Query/PendingReveals"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).PendingReveals(ctx, req.(*QueryPendingRevealsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// --- Service Descriptors ---

var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.mevprotection.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CommitTx",
			Handler:    _Msg_CommitTx_Handler,
		},
		{
			MethodName: "RevealTx",
			Handler:    _Msg_RevealTx_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/mevprotection/tx.proto",
}

func _Query_MEVRedistribution_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryMEVRedistributionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).MEVRedistribution(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.mevprotection.Query/MEVRedistribution"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).MEVRedistribution(ctx, req.(*QueryMEVRedistributionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.mevprotection.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "CommittedTx",
			Handler:    _Query_CommittedTx_Handler,
		},
		{
			MethodName: "PendingReveals",
			Handler:    _Query_PendingReveals_Handler,
		},
		{
			MethodName: "MEVRedistribution",
			Handler:    _Query_MEVRedistribution_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/mevprotection/query.proto",
}

// RegisterMsgServer registers the MsgServer implementation with the gRPC server
func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&Msg_ServiceDesc, srv)
}

// RegisterQueryServer registers the QueryServer implementation with the gRPC server
func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
}
