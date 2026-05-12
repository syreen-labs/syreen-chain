package types

import (
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

// combinedResolver tries a local file first, then falls back to global registry.
type combinedResolver struct {
	local protoreflect.FileDescriptor
}

func (r combinedResolver) FindFileByPath(path string) (protoreflect.FileDescriptor, error) {
	if r.local != nil && r.local.Path() == path {
		return r.local, nil
	}
	return protoregistry.GlobalFiles.FindFileByPath(path)
}

func (r combinedResolver) FindDescriptorByName(name protoreflect.FullName) (protoreflect.Descriptor, error) {
	if r.local != nil {
		if d := r.local.Messages().ByName(name.Name()); d != nil {
			return d, nil
		}
	}
	return protoregistry.GlobalFiles.FindDescriptorByName(name)
}

// TxFileDescProto holds the raw FileDescriptorProto for codec_proto.go to use.
var TxFileDescProto *descriptorpb.FileDescriptorProto

func init() {
	registerProtoFileDescriptors()
}

func registerProtoFileDescriptors() {
	strType := descriptorpb.FieldDescriptorProto_TYPE_STRING
	uint64Type := descriptorpb.FieldDescriptorProto_TYPE_UINT64
	int64Type := descriptorpb.FieldDescriptorProto_TYPE_INT64
	boolType := descriptorpb.FieldDescriptorProto_TYPE_BOOL
	msgType := descriptorpb.FieldDescriptorProto_TYPE_MESSAGE
	labelOpt := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	labelRep := descriptorpb.FieldDescriptorProto_LABEL_REPEATED

	// Coin message descriptor (inline, since cosmos.base.v1beta1.Coin may not be registered yet)
	coinDesc := &descriptorpb.DescriptorProto{
		Name: sp("Coin"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{Name: sp("denom"), Number: ip(1), Type: &strType, Label: &labelOpt},
			{Name: sp("amount"), Number: ip(2), Type: &strType, Label: &labelOpt},
		},
	}

	fd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/payments/tx.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.payments"),
		MessageType: []*descriptorpb.DescriptorProto{
			coinDesc,
			{
				Name: sp("MsgCreateInvoice"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("creator"), Number: ip(1), Type: &strType, Label: &labelOpt},
					{Name: sp("booking_id"), Number: ip(2), Type: &strType, Label: &labelOpt},
					{Name: sp("payee"), Number: ip(3), Type: &strType, Label: &labelOpt},
					{Name: sp("amount"), Number: ip(4), Type: &msgType, TypeName: sp(".syreen.payments.Coin"), Label: &labelOpt},
					{Name: sp("description"), Number: ip(5), Type: &strType, Label: &labelOpt},
					{Name: sp("currency"), Number: ip(6), Type: &strType, Label: &labelOpt},
				},
			},
			{
				Name: sp("MsgCreateInvoiceResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("invoice_id"), Number: ip(1), Type: &strType, Label: &labelOpt},
				},
			},
			{
				Name: sp("MsgPayInvoice"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("payer"), Number: ip(1), Type: &strType, Label: &labelOpt},
					{Name: sp("invoice_id"), Number: ip(2), Type: &strType, Label: &labelOpt},
					{Name: sp("payment_method"), Number: ip(3), Type: &strType, Label: &labelOpt},
				},
			},
			{Name: sp("MsgPayInvoiceResponse")},
			{
				Name: sp("MsgRefundPayment"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("authority"), Number: ip(1), Type: &strType, Label: &labelOpt},
					{Name: sp("invoice_id"), Number: ip(2), Type: &strType, Label: &labelOpt},
					{Name: sp("reason"), Number: ip(3), Type: &strType, Label: &labelOpt},
				},
			},
			{Name: sp("MsgRefundPaymentResponse")},
			{
				Name: sp("MsgSetExchangeRate"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("authority"), Number: ip(1), Type: &strType, Label: &labelOpt},
					{Name: sp("from_denom"), Number: ip(2), Type: &strType, Label: &labelOpt},
					{Name: sp("to_denom"), Number: ip(3), Type: &strType, Label: &labelOpt},
					{Name: sp("rate"), Number: ip(4), Type: &strType, Label: &labelOpt},
				},
			},
			{Name: sp("MsgSetExchangeRateResponse")},
			{
				Name: sp("MsgWithdrawEarnings"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("address"), Number: ip(1), Type: &strType, Label: &labelOpt},
					{Name: sp("amount"), Number: ip(2), Type: &msgType, TypeName: sp(".syreen.payments.Coin"), Label: &labelOpt},
				},
			},
			{Name: sp("MsgWithdrawEarningsResponse")},
			{
				Name: sp("MsgCancelInvoice"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("creator"), Number: ip(1), Type: &strType, Label: &labelOpt},
					{Name: sp("invoice_id"), Number: ip(2), Type: &strType, Label: &labelOpt},
				},
			},
			{Name: sp("MsgCancelInvoiceResponse")},
			{
				Name: sp("Invoice"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("id"), Number: ip(1), Type: &strType, Label: &labelOpt},
					{Name: sp("booking_id"), Number: ip(2), Type: &strType, Label: &labelOpt},
					{Name: sp("payer"), Number: ip(3), Type: &strType, Label: &labelOpt},
					{Name: sp("payee"), Number: ip(4), Type: &strType, Label: &labelOpt},
					{Name: sp("amount"), Number: ip(5), Type: &msgType, TypeName: sp(".syreen.payments.Coin"), Label: &labelOpt},
					{Name: sp("platform_fee"), Number: ip(6), Type: &msgType, TypeName: sp(".syreen.payments.Coin"), Label: &labelOpt},
					{Name: sp("insurance_fee"), Number: ip(7), Type: &msgType, TypeName: sp(".syreen.payments.Coin"), Label: &labelOpt},
					{Name: sp("net_amount"), Number: ip(8), Type: &msgType, TypeName: sp(".syreen.payments.Coin"), Label: &labelOpt},
					{Name: sp("status"), Number: ip(9), Type: &strType, Label: &labelOpt},
					{Name: sp("payment_method"), Number: ip(10), Type: &strType, Label: &labelOpt},
					{Name: sp("payment_tx_hash"), Number: ip(11), Type: &strType, Label: &labelOpt},
					{Name: sp("created_at_block"), Number: ip(12), Type: &int64Type, Label: &labelOpt},
					{Name: sp("paid_at_block"), Number: ip(13), Type: &int64Type, Label: &labelOpt},
					{Name: sp("expires_at_block"), Number: ip(14), Type: &int64Type, Label: &labelOpt},
					{Name: sp("description"), Number: ip(15), Type: &strType, Label: &labelOpt},
					{Name: sp("currency"), Number: ip(16), Type: &strType, Label: &labelOpt},
					{Name: sp("creator"), Number: ip(17), Type: &strType, Label: &labelOpt},
				},
			},
			{
				Name: sp("ExchangeRate"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("from_denom"), Number: ip(1), Type: &strType, Label: &labelOpt},
					{Name: sp("to_denom"), Number: ip(2), Type: &strType, Label: &labelOpt},
					{Name: sp("rate"), Number: ip(3), Type: &strType, Label: &labelOpt},
					{Name: sp("updated_at_block"), Number: ip(4), Type: &int64Type, Label: &labelOpt},
					{Name: sp("updated_by"), Number: ip(5), Type: &strType, Label: &labelOpt},
				},
			},
			{
				Name: sp("Earnings"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("address"), Number: ip(1), Type: &strType, Label: &labelOpt},
					{Name: sp("available"), Number: ip(2), Type: &msgType, TypeName: sp(".syreen.payments.Coin"), Label: &labelRep},
					{Name: sp("total_earned"), Number: ip(3), Type: &msgType, TypeName: sp(".syreen.payments.Coin"), Label: &labelRep},
					{Name: sp("total_withdrawn"), Number: ip(4), Type: &msgType, TypeName: sp(".syreen.payments.Coin"), Label: &labelRep},
				},
			},
			{
				Name: sp("Params"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("platform_fee_percent"), Number: ip(1), Type: &strType, Label: &labelOpt},
					{Name: sp("insurance_fee_percent"), Number: ip(2), Type: &strType, Label: &labelOpt},
					{Name: sp("invoice_expiry_blocks"), Number: ip(3), Type: &uint64Type, Label: &labelOpt},
					{Name: sp("enable_payments"), Number: ip(4), Type: &boolType, Label: &labelOpt},
					{Name: sp("platform_address"), Number: ip(5), Type: &strType, Label: &labelOpt},
					{Name: sp("insurance_pool_address"), Number: ip(6), Type: &strType, Label: &labelOpt},
				},
			},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: sp("CreateInvoice"), InputType: sp(".syreen.payments.MsgCreateInvoice"), OutputType: sp(".syreen.payments.MsgCreateInvoiceResponse")},
					{Name: sp("PayInvoice"), InputType: sp(".syreen.payments.MsgPayInvoice"), OutputType: sp(".syreen.payments.MsgPayInvoiceResponse")},
					{Name: sp("RefundPayment"), InputType: sp(".syreen.payments.MsgRefundPayment"), OutputType: sp(".syreen.payments.MsgRefundPaymentResponse")},
					{Name: sp("SetExchangeRate"), InputType: sp(".syreen.payments.MsgSetExchangeRate"), OutputType: sp(".syreen.payments.MsgSetExchangeRateResponse")},
					{Name: sp("WithdrawEarnings"), InputType: sp(".syreen.payments.MsgWithdrawEarnings"), OutputType: sp(".syreen.payments.MsgWithdrawEarningsResponse")},
					{Name: sp("CancelInvoice"), InputType: sp(".syreen.payments.MsgCancelInvoice"), OutputType: sp(".syreen.payments.MsgCancelInvoiceResponse")},
				},
			},
		},
	}

	// Query file descriptor
	qfd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/payments/query.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.payments"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: sp("QueryParamsRequest")},
			{
				Name: sp("QueryParamsResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("params"), Number: ip(1), Type: &msgType, TypeName: sp(".syreen.payments.Params"), Label: &labelOpt},
				},
			},
			{
				Name: sp("QueryInvoiceRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("invoice_id"), Number: ip(1), Type: &strType, Label: &labelOpt},
				},
			},
			{
				Name: sp("QueryInvoiceResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("invoice"), Number: ip(1), Type: &msgType, TypeName: sp(".syreen.payments.Invoice"), Label: &labelOpt},
					{Name: sp("found"), Number: ip(2), Type: &boolType, Label: &labelOpt},
				},
			},
			{
				Name: sp("QueryInvoicesByBookingRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("booking_id"), Number: ip(1), Type: &strType, Label: &labelOpt},
				},
			},
			{
				Name: sp("QueryInvoicesByBookingResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("invoices"), Number: ip(1), Type: &msgType, TypeName: sp(".syreen.payments.Invoice"), Label: &labelRep},
				},
			},
			{
				Name: sp("QueryInvoicesByPayerRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("payer"), Number: ip(1), Type: &strType, Label: &labelOpt},
				},
			},
			{
				Name: sp("QueryInvoicesByPayerResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("invoices"), Number: ip(1), Type: &msgType, TypeName: sp(".syreen.payments.Invoice"), Label: &labelRep},
				},
			},
			{
				Name: sp("QueryExchangeRateRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("from_denom"), Number: ip(1), Type: &strType, Label: &labelOpt},
					{Name: sp("to_denom"), Number: ip(2), Type: &strType, Label: &labelOpt},
				},
			},
			{
				Name: sp("QueryExchangeRateResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("exchange_rate"), Number: ip(1), Type: &msgType, TypeName: sp(".syreen.payments.ExchangeRate"), Label: &labelOpt},
					{Name: sp("found"), Number: ip(2), Type: &boolType, Label: &labelOpt},
				},
			},
			{
				Name: sp("QueryEarningsRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("address"), Number: ip(1), Type: &strType, Label: &labelOpt},
				},
			},
			{
				Name: sp("QueryEarningsResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("earnings"), Number: ip(1), Type: &msgType, TypeName: sp(".syreen.payments.Earnings"), Label: &labelOpt},
					{Name: sp("found"), Number: ip(2), Type: &boolType, Label: &labelOpt},
				},
			},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Query"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: sp("Params"), InputType: sp(".syreen.payments.QueryParamsRequest"), OutputType: sp(".syreen.payments.QueryParamsResponse")},
					{Name: sp("Invoice"), InputType: sp(".syreen.payments.QueryInvoiceRequest"), OutputType: sp(".syreen.payments.QueryInvoiceResponse")},
					{Name: sp("InvoicesByBooking"), InputType: sp(".syreen.payments.QueryInvoicesByBookingRequest"), OutputType: sp(".syreen.payments.QueryInvoicesByBookingResponse")},
					{Name: sp("InvoicesByPayer"), InputType: sp(".syreen.payments.QueryInvoicesByPayerRequest"), OutputType: sp(".syreen.payments.QueryInvoicesByPayerResponse")},
					{Name: sp("ExchangeRate"), InputType: sp(".syreen.payments.QueryExchangeRateRequest"), OutputType: sp(".syreen.payments.QueryExchangeRateResponse")},
					{Name: sp("Earnings"), InputType: sp(".syreen.payments.QueryEarningsRequest"), OutputType: sp(".syreen.payments.QueryEarningsResponse")},
				},
			},
		},
	}

	TxFileDescProto = fd

	file, err := protodesc.NewFile(fd, protoregistry.GlobalFiles)
	if err != nil {
		// Message type refs not yet registered - convert to bytes as fallback
		for _, mt := range fd.MessageType {
			for _, f := range mt.Field {
				if f.TypeName != nil && *f.Type == msgType {
					bytesT := descriptorpb.FieldDescriptorProto_TYPE_BYTES
					f.Type = &bytesT
					f.TypeName = nil
				}
			}
		}
		file, err = protodesc.NewFile(fd, nil)
		if err != nil {
			panic("payments: failed to create proto file descriptor: " + err.Error())
		}
	}
	if err := protoregistry.GlobalFiles.RegisterFile(file); err != nil {
		// Already registered is fine
	}

	// Same fallback for query file
	for _, mt := range qfd.MessageType {
		for _, f := range mt.Field {
			if f.TypeName != nil && *f.Type == msgType {
				tn := *f.TypeName
				if tn != ".syreen.payments.Invoice" &&
					tn != ".syreen.payments.ExchangeRate" &&
					tn != ".syreen.payments.Earnings" &&
					tn != ".syreen.payments.Params" {
					bytesT := descriptorpb.FieldDescriptorProto_TYPE_BYTES
					f.Type = &bytesT
					f.TypeName = nil
				}
			}
		}
	}

	qfile, err := protodesc.NewFile(qfd, combinedResolver{file})
	if err != nil {
		// Last resort: strip all message type refs
		for _, mt := range qfd.MessageType {
			for _, f := range mt.Field {
				if f.TypeName != nil && *f.Type == msgType {
					bytesT := descriptorpb.FieldDescriptorProto_TYPE_BYTES
					f.Type = &bytesT
					f.TypeName = nil
				}
			}
		}
		qfile, err = protodesc.NewFile(qfd, nil)
		if err != nil {
			panic("payments: failed to create query proto file descriptor: " + err.Error())
		}
	}
	if err := protoregistry.GlobalFiles.RegisterFile(qfile); err != nil {
		// Already registered is fine
	}
}

func sp(s string) *string { return &s }
func ip(i int32) *int32   { return &i }
