package types

import (
	"bytes"
	"compress/gzip"
	"sync"

	gogoproto "github.com/cosmos/gogoproto/proto"
	"google.golang.org/protobuf/proto"
)

// This file registers all payments types with gogoproto and provides
// Descriptor() methods so the Cosmos SDK can properly encode/decode transactions.

var (
	fileDescOnce sync.Once
	fileDescGzip []byte
)

func txFileDescGzip() []byte {
	fileDescOnce.Do(func() {
		if TxFileDescProto == nil {
			return
		}
		raw, err := proto.Marshal(TxFileDescProto)
		if err != nil {
			return
		}
		var buf bytes.Buffer
		w := gzip.NewWriter(&buf)
		w.Write(raw)
		w.Close()
		fileDescGzip = buf.Bytes()
	})
	return fileDescGzip
}

func init() {
	gogoproto.RegisterType((*MsgCreateInvoice)(nil), "syreen.payments.MsgCreateInvoice")
	gogoproto.RegisterType((*MsgCreateInvoiceResponse)(nil), "syreen.payments.MsgCreateInvoiceResponse")
	gogoproto.RegisterType((*MsgPayInvoice)(nil), "syreen.payments.MsgPayInvoice")
	gogoproto.RegisterType((*MsgPayInvoiceResponse)(nil), "syreen.payments.MsgPayInvoiceResponse")
	gogoproto.RegisterType((*MsgRefundPayment)(nil), "syreen.payments.MsgRefundPayment")
	gogoproto.RegisterType((*MsgRefundPaymentResponse)(nil), "syreen.payments.MsgRefundPaymentResponse")
	gogoproto.RegisterType((*MsgSetExchangeRate)(nil), "syreen.payments.MsgSetExchangeRate")
	gogoproto.RegisterType((*MsgSetExchangeRateResponse)(nil), "syreen.payments.MsgSetExchangeRateResponse")
	gogoproto.RegisterType((*MsgWithdrawEarnings)(nil), "syreen.payments.MsgWithdrawEarnings")
	gogoproto.RegisterType((*MsgWithdrawEarningsResponse)(nil), "syreen.payments.MsgWithdrawEarningsResponse")
	gogoproto.RegisterType((*MsgCancelInvoice)(nil), "syreen.payments.MsgCancelInvoice")
	gogoproto.RegisterType((*MsgCancelInvoiceResponse)(nil), "syreen.payments.MsgCancelInvoiceResponse")
	gogoproto.RegisterType((*Invoice)(nil), "syreen.payments.Invoice")
	gogoproto.RegisterType((*ExchangeRate)(nil), "syreen.payments.ExchangeRate")
	gogoproto.RegisterType((*Earnings)(nil), "syreen.payments.Earnings")
	gogoproto.RegisterType((*Params)(nil), "syreen.payments.Params")
}

// Descriptor returns the gzipped proto file descriptor and path to this message.
// Required by gogoproto for proper Any encoding in Cosmos SDK transactions.
// Indices match the order in proto_register.go MessageType array:
// 0=Coin, 1=MsgCreateInvoice, 2=MsgCreateInvoiceResponse, 3=MsgPayInvoice, ...

func (*MsgCreateInvoice) Descriptor() ([]byte, []int)          { return txFileDescGzip(), []int{1} }
func (*MsgCreateInvoiceResponse) Descriptor() ([]byte, []int)  { return txFileDescGzip(), []int{2} }
func (*MsgPayInvoice) Descriptor() ([]byte, []int)             { return txFileDescGzip(), []int{3} }
func (*MsgPayInvoiceResponse) Descriptor() ([]byte, []int)     { return txFileDescGzip(), []int{4} }
func (*MsgRefundPayment) Descriptor() ([]byte, []int)          { return txFileDescGzip(), []int{5} }
func (*MsgRefundPaymentResponse) Descriptor() ([]byte, []int)  { return txFileDescGzip(), []int{6} }
func (*MsgSetExchangeRate) Descriptor() ([]byte, []int)        { return txFileDescGzip(), []int{7} }
func (*MsgSetExchangeRateResponse) Descriptor() ([]byte, []int) { return txFileDescGzip(), []int{8} }
func (*MsgWithdrawEarnings) Descriptor() ([]byte, []int)       { return txFileDescGzip(), []int{9} }
func (*MsgWithdrawEarningsResponse) Descriptor() ([]byte, []int) { return txFileDescGzip(), []int{10} }
func (*MsgCancelInvoice) Descriptor() ([]byte, []int)          { return txFileDescGzip(), []int{11} }
func (*MsgCancelInvoiceResponse) Descriptor() ([]byte, []int)  { return txFileDescGzip(), []int{12} }
func (*Invoice) Descriptor() ([]byte, []int)                   { return txFileDescGzip(), []int{13} }
func (*ExchangeRate) Descriptor() ([]byte, []int)              { return txFileDescGzip(), []int{14} }
func (*Earnings) Descriptor() ([]byte, []int)                  { return txFileDescGzip(), []int{15} }
func (*Params) Descriptor() ([]byte, []int)                    { return txFileDescGzip(), []int{16} }
