package types

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

// ---------------------------------------------------------------------------
// Proto-compatible message methods for MsgEthereumTx
// ---------------------------------------------------------------------------

func (m *MsgEthereumTx) ProtoMessage()           {}
func (m *MsgEthereumTx) Reset()                  { *m = MsgEthereumTx{} }
func (m *MsgEthereumTx) String() string          { return fmt.Sprintf("from=%s to=%s", m.From, m.To) }
func (m *MsgEthereumTx) XXX_MessageName() string { return "syreen.evm.MsgEthereumTx" }

// ---------------------------------------------------------------------------
// Proto-compatible message methods for MsgEthereumTxResponse
// ---------------------------------------------------------------------------

func (m *MsgEthereumTxResponse) ProtoMessage()           {}
func (m *MsgEthereumTxResponse) Reset()                  { *m = MsgEthereumTxResponse{} }
func (m *MsgEthereumTxResponse) String() string          { return fmt.Sprintf("gas_used=%d", m.GasUsed) }
func (m *MsgEthereumTxResponse) XXX_MessageName() string { return "syreen.evm.MsgEthereumTxResponse" }

// ---------------------------------------------------------------------------
// MsgServer Interface (H3)
// ---------------------------------------------------------------------------

type MsgServer interface {
	EthereumTx(context.Context, *MsgEthereumTx) (*MsgEthereumTxResponse, error)
}

// ---------------------------------------------------------------------------
// gRPC Msg Handler
// ---------------------------------------------------------------------------

func _Msg_EthereumTx_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgEthereumTx)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).EthereumTx(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syreen.evm.Msg/EthereumTx",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).EthereumTx(ctx, req.(*MsgEthereumTx))
	}
	return interceptor(ctx, in, info, handler)
}

// ---------------------------------------------------------------------------
// gRPC Service Descriptor
// ---------------------------------------------------------------------------

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.evm.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "EthereumTx",
			Handler:    _Msg_EthereumTx_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/evm/tx.proto",
}

// ---------------------------------------------------------------------------
// Registration Functions (H3)
// ---------------------------------------------------------------------------

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}
