package types

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

// --- Msg Response Types ---

type MsgCreateInvoiceResponse struct {
	InvoiceID string `protobuf:"bytes,1,opt,name=invoice_id,json=invoiceId,proto3" json:"invoice_id"`
}
func (m *MsgCreateInvoiceResponse) ProtoMessage()           {}
func (m *MsgCreateInvoiceResponse) Reset()                  { *m = MsgCreateInvoiceResponse{} }
func (m *MsgCreateInvoiceResponse) String() string          { return fmt.Sprintf("create_invoice_response: %s", m.InvoiceID) }
func (m *MsgCreateInvoiceResponse) XXX_MessageName() string { return "syreen.payments.MsgCreateInvoiceResponse" }

type MsgPayInvoiceResponse struct{}
func (m *MsgPayInvoiceResponse) ProtoMessage()           {}
func (m *MsgPayInvoiceResponse) Reset()                  { *m = MsgPayInvoiceResponse{} }
func (m *MsgPayInvoiceResponse) String() string          { return "pay_invoice_response" }
func (m *MsgPayInvoiceResponse) XXX_MessageName() string { return "syreen.payments.MsgPayInvoiceResponse" }

type MsgRefundPaymentResponse struct{}
func (m *MsgRefundPaymentResponse) ProtoMessage()           {}
func (m *MsgRefundPaymentResponse) Reset()                  { *m = MsgRefundPaymentResponse{} }
func (m *MsgRefundPaymentResponse) String() string          { return "refund_payment_response" }
func (m *MsgRefundPaymentResponse) XXX_MessageName() string { return "syreen.payments.MsgRefundPaymentResponse" }

type MsgSetExchangeRateResponse struct{}
func (m *MsgSetExchangeRateResponse) ProtoMessage()           {}
func (m *MsgSetExchangeRateResponse) Reset()                  { *m = MsgSetExchangeRateResponse{} }
func (m *MsgSetExchangeRateResponse) String() string          { return "set_exchange_rate_response" }
func (m *MsgSetExchangeRateResponse) XXX_MessageName() string { return "syreen.payments.MsgSetExchangeRateResponse" }

type MsgWithdrawEarningsResponse struct{}
func (m *MsgWithdrawEarningsResponse) ProtoMessage()           {}
func (m *MsgWithdrawEarningsResponse) Reset()                  { *m = MsgWithdrawEarningsResponse{} }
func (m *MsgWithdrawEarningsResponse) String() string          { return "withdraw_earnings_response" }
func (m *MsgWithdrawEarningsResponse) XXX_MessageName() string { return "syreen.payments.MsgWithdrawEarningsResponse" }

type MsgCancelInvoiceResponse struct{}
func (m *MsgCancelInvoiceResponse) ProtoMessage()           {}
func (m *MsgCancelInvoiceResponse) Reset()                  { *m = MsgCancelInvoiceResponse{} }
func (m *MsgCancelInvoiceResponse) String() string          { return "cancel_invoice_response" }
func (m *MsgCancelInvoiceResponse) XXX_MessageName() string { return "syreen.payments.MsgCancelInvoiceResponse" }

// --- Query Types ---

type QueryParamsRequest struct{}
func (m *QueryParamsRequest) ProtoMessage()           {}
func (m *QueryParamsRequest) Reset()                  { *m = QueryParamsRequest{} }
func (m *QueryParamsRequest) String() string          { return "query_params_request" }
func (m *QueryParamsRequest) XXX_MessageName() string { return "syreen.payments.QueryParamsRequest" }

type QueryParamsResponse struct {
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}
func (m *QueryParamsResponse) ProtoMessage()           {}
func (m *QueryParamsResponse) Reset()                  { *m = QueryParamsResponse{} }
func (m *QueryParamsResponse) String() string          { return fmt.Sprintf("params: %+v", m.Params) }
func (m *QueryParamsResponse) XXX_MessageName() string { return "syreen.payments.QueryParamsResponse" }

type QueryInvoiceRequest struct {
	InvoiceID string `protobuf:"bytes,1,opt,name=invoice_id,json=invoiceId,proto3" json:"invoice_id"`
}
func (m *QueryInvoiceRequest) ProtoMessage()           {}
func (m *QueryInvoiceRequest) Reset()                  { *m = QueryInvoiceRequest{} }
func (m *QueryInvoiceRequest) String() string          { return fmt.Sprintf("query_invoice: %s", m.InvoiceID) }
func (m *QueryInvoiceRequest) XXX_MessageName() string { return "syreen.payments.QueryInvoiceRequest" }

type QueryInvoiceResponse struct {
	Invoice Invoice `protobuf:"bytes,1,opt,name=invoice,proto3" json:"invoice"`
	Found   bool    `protobuf:"varint,2,opt,name=found,proto3" json:"found"`
}
func (m *QueryInvoiceResponse) ProtoMessage()           {}
func (m *QueryInvoiceResponse) Reset()                  { *m = QueryInvoiceResponse{} }
func (m *QueryInvoiceResponse) String() string          { return fmt.Sprintf("invoice: %+v", m.Invoice) }
func (m *QueryInvoiceResponse) XXX_MessageName() string { return "syreen.payments.QueryInvoiceResponse" }

type QueryInvoicesByBookingRequest struct {
	BookingID string `protobuf:"bytes,1,opt,name=booking_id,json=bookingId,proto3" json:"booking_id"`
}
func (m *QueryInvoicesByBookingRequest) ProtoMessage()           {}
func (m *QueryInvoicesByBookingRequest) Reset()                  { *m = QueryInvoicesByBookingRequest{} }
func (m *QueryInvoicesByBookingRequest) String() string          { return fmt.Sprintf("query_invoices_by_booking: %s", m.BookingID) }
func (m *QueryInvoicesByBookingRequest) XXX_MessageName() string { return "syreen.payments.QueryInvoicesByBookingRequest" }

type QueryInvoicesByBookingResponse struct {
	Invoices []Invoice `protobuf:"bytes,1,rep,name=invoices,proto3" json:"invoices"`
}
func (m *QueryInvoicesByBookingResponse) ProtoMessage()           {}
func (m *QueryInvoicesByBookingResponse) Reset()                  { *m = QueryInvoicesByBookingResponse{} }
func (m *QueryInvoicesByBookingResponse) String() string          { return fmt.Sprintf("invoices: %d", len(m.Invoices)) }
func (m *QueryInvoicesByBookingResponse) XXX_MessageName() string { return "syreen.payments.QueryInvoicesByBookingResponse" }

type QueryInvoicesByPayerRequest struct {
	Payer string `protobuf:"bytes,1,opt,name=payer,proto3" json:"payer"`
}
func (m *QueryInvoicesByPayerRequest) ProtoMessage()           {}
func (m *QueryInvoicesByPayerRequest) Reset()                  { *m = QueryInvoicesByPayerRequest{} }
func (m *QueryInvoicesByPayerRequest) String() string          { return fmt.Sprintf("query_invoices_by_payer: %s", m.Payer) }
func (m *QueryInvoicesByPayerRequest) XXX_MessageName() string { return "syreen.payments.QueryInvoicesByPayerRequest" }

type QueryInvoicesByPayerResponse struct {
	Invoices []Invoice `protobuf:"bytes,1,rep,name=invoices,proto3" json:"invoices"`
}
func (m *QueryInvoicesByPayerResponse) ProtoMessage()           {}
func (m *QueryInvoicesByPayerResponse) Reset()                  { *m = QueryInvoicesByPayerResponse{} }
func (m *QueryInvoicesByPayerResponse) String() string          { return fmt.Sprintf("invoices: %d", len(m.Invoices)) }
func (m *QueryInvoicesByPayerResponse) XXX_MessageName() string { return "syreen.payments.QueryInvoicesByPayerResponse" }

type QueryExchangeRateRequest struct {
	FromDenom string `protobuf:"bytes,1,opt,name=from_denom,json=fromDenom,proto3" json:"from_denom"`
	ToDenom   string `protobuf:"bytes,2,opt,name=to_denom,json=toDenom,proto3" json:"to_denom"`
}
func (m *QueryExchangeRateRequest) ProtoMessage()           {}
func (m *QueryExchangeRateRequest) Reset()                  { *m = QueryExchangeRateRequest{} }
func (m *QueryExchangeRateRequest) String() string          { return fmt.Sprintf("query_exchange_rate: %s/%s", m.FromDenom, m.ToDenom) }
func (m *QueryExchangeRateRequest) XXX_MessageName() string { return "syreen.payments.QueryExchangeRateRequest" }

type QueryExchangeRateResponse struct {
	ExchangeRate ExchangeRate `protobuf:"bytes,1,opt,name=exchange_rate,json=exchangeRate,proto3" json:"exchange_rate"`
	Found        bool         `protobuf:"varint,2,opt,name=found,proto3" json:"found"`
}
func (m *QueryExchangeRateResponse) ProtoMessage()           {}
func (m *QueryExchangeRateResponse) Reset()                  { *m = QueryExchangeRateResponse{} }
func (m *QueryExchangeRateResponse) String() string          { return fmt.Sprintf("exchange_rate: %+v", m.ExchangeRate) }
func (m *QueryExchangeRateResponse) XXX_MessageName() string { return "syreen.payments.QueryExchangeRateResponse" }

type QueryEarningsRequest struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address"`
}
func (m *QueryEarningsRequest) ProtoMessage()           {}
func (m *QueryEarningsRequest) Reset()                  { *m = QueryEarningsRequest{} }
func (m *QueryEarningsRequest) String() string          { return fmt.Sprintf("query_earnings: %s", m.Address) }
func (m *QueryEarningsRequest) XXX_MessageName() string { return "syreen.payments.QueryEarningsRequest" }

type QueryEarningsResponse struct {
	Earnings Earnings `protobuf:"bytes,1,opt,name=earnings,proto3" json:"earnings"`
	Found    bool     `protobuf:"varint,2,opt,name=found,proto3" json:"found"`
}
func (m *QueryEarningsResponse) ProtoMessage()           {}
func (m *QueryEarningsResponse) Reset()                  { *m = QueryEarningsResponse{} }
func (m *QueryEarningsResponse) String() string          { return fmt.Sprintf("earnings: %+v", m.Earnings) }
func (m *QueryEarningsResponse) XXX_MessageName() string { return "syreen.payments.QueryEarningsResponse" }

// --- Server Interfaces ---

type MsgServer interface {
	CreateInvoice(context.Context, *MsgCreateInvoice) (*MsgCreateInvoiceResponse, error)
	PayInvoice(context.Context, *MsgPayInvoice) (*MsgPayInvoiceResponse, error)
	RefundPayment(context.Context, *MsgRefundPayment) (*MsgRefundPaymentResponse, error)
	SetExchangeRate(context.Context, *MsgSetExchangeRate) (*MsgSetExchangeRateResponse, error)
	WithdrawEarnings(context.Context, *MsgWithdrawEarnings) (*MsgWithdrawEarningsResponse, error)
	CancelInvoice(context.Context, *MsgCancelInvoice) (*MsgCancelInvoiceResponse, error)
}

type QueryServer interface {
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	Invoice(context.Context, *QueryInvoiceRequest) (*QueryInvoiceResponse, error)
	InvoicesByBooking(context.Context, *QueryInvoicesByBookingRequest) (*QueryInvoicesByBookingResponse, error)
	InvoicesByPayer(context.Context, *QueryInvoicesByPayerRequest) (*QueryInvoicesByPayerResponse, error)
	ExchangeRate(context.Context, *QueryExchangeRateRequest) (*QueryExchangeRateResponse, error)
	Earnings(context.Context, *QueryEarningsRequest) (*QueryEarningsResponse, error)
}

// --- gRPC Handlers ---

func _Msg_CreateInvoice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateInvoice)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).CreateInvoice(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.payments.Msg/CreateInvoice"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).CreateInvoice(ctx, req.(*MsgCreateInvoice)) })
}

func _Msg_PayInvoice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgPayInvoice)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).PayInvoice(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.payments.Msg/PayInvoice"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).PayInvoice(ctx, req.(*MsgPayInvoice)) })
}

func _Msg_RefundPayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRefundPayment)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).RefundPayment(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.payments.Msg/RefundPayment"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).RefundPayment(ctx, req.(*MsgRefundPayment)) })
}

func _Msg_SetExchangeRate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSetExchangeRate)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).SetExchangeRate(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.payments.Msg/SetExchangeRate"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).SetExchangeRate(ctx, req.(*MsgSetExchangeRate)) })
}

func _Msg_WithdrawEarnings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgWithdrawEarnings)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).WithdrawEarnings(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.payments.Msg/WithdrawEarnings"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).WithdrawEarnings(ctx, req.(*MsgWithdrawEarnings)) })
}

func _Msg_CancelInvoice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCancelInvoice)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MsgServer).CancelInvoice(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.payments.Msg/CancelInvoice"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MsgServer).CancelInvoice(ctx, req.(*MsgCancelInvoice)) })
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(QueryServer).Params(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.payments.Query/Params"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest)) })
}

func _Query_Invoice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryInvoiceRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(QueryServer).Invoice(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.payments.Query/Invoice"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(QueryServer).Invoice(ctx, req.(*QueryInvoiceRequest)) })
}

func _Query_InvoicesByBooking_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryInvoicesByBookingRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(QueryServer).InvoicesByBooking(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.payments.Query/InvoicesByBooking"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(QueryServer).InvoicesByBooking(ctx, req.(*QueryInvoicesByBookingRequest)) })
}

func _Query_InvoicesByPayer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryInvoicesByPayerRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(QueryServer).InvoicesByPayer(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.payments.Query/InvoicesByPayer"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(QueryServer).InvoicesByPayer(ctx, req.(*QueryInvoicesByPayerRequest)) })
}

func _Query_ExchangeRate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryExchangeRateRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(QueryServer).ExchangeRate(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.payments.Query/ExchangeRate"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(QueryServer).ExchangeRate(ctx, req.(*QueryExchangeRateRequest)) })
}

func _Query_Earnings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryEarningsRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(QueryServer).Earnings(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/syreen.payments.Query/Earnings"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(QueryServer).Earnings(ctx, req.(*QueryEarningsRequest)) })
}

// --- Service Descriptors ---

var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.payments.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreateInvoice", Handler: _Msg_CreateInvoice_Handler},
		{MethodName: "PayInvoice", Handler: _Msg_PayInvoice_Handler},
		{MethodName: "RefundPayment", Handler: _Msg_RefundPayment_Handler},
		{MethodName: "SetExchangeRate", Handler: _Msg_SetExchangeRate_Handler},
		{MethodName: "WithdrawEarnings", Handler: _Msg_WithdrawEarnings_Handler},
		{MethodName: "CancelInvoice", Handler: _Msg_CancelInvoice_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/payments/tx.proto",
}

var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "syreen.payments.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Params", Handler: _Query_Params_Handler},
		{MethodName: "Invoice", Handler: _Query_Invoice_Handler},
		{MethodName: "InvoicesByBooking", Handler: _Query_InvoicesByBooking_Handler},
		{MethodName: "InvoicesByPayer", Handler: _Query_InvoicesByPayer_Handler},
		{MethodName: "ExchangeRate", Handler: _Query_ExchangeRate_Handler},
		{MethodName: "Earnings", Handler: _Query_Earnings_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "syreen/payments/query.proto",
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
}
