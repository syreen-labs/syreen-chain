package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgCreateInvoice    = "create_invoice"
	TypeMsgPayInvoice       = "pay_invoice"
	TypeMsgRefundPayment    = "refund_payment"
	TypeMsgSetExchangeRate  = "set_exchange_rate"
	TypeMsgWithdrawEarnings = "withdraw_earnings"
	TypeMsgCancelInvoice    = "cancel_invoice"
)

// --- MsgCreateInvoice --- Creates a new payment invoice
type MsgCreateInvoice struct {
	Creator     string   `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator"`
	BookingID   string   `protobuf:"bytes,2,opt,name=booking_id,json=bookingId,proto3" json:"booking_id"`
	Payee       string   `protobuf:"bytes,3,opt,name=payee,proto3" json:"payee"`
	Amount      sdk.Coin `protobuf:"bytes,4,opt,name=amount,proto3" json:"amount"`
	Description string   `protobuf:"bytes,5,opt,name=description,proto3" json:"description"`
	Currency    string   `protobuf:"bytes,6,opt,name=currency,proto3" json:"currency"`
}

func (m *MsgCreateInvoice) ProtoMessage()           {}
func (m *MsgCreateInvoice) Reset()                  { *m = MsgCreateInvoice{} }
func (m *MsgCreateInvoice) String() string          { return fmt.Sprintf("create_invoice: booking=%s amount=%s", m.BookingID, m.Amount) }
func (m *MsgCreateInvoice) XXX_MessageName() string { return "syreen.payments.MsgCreateInvoice" }

func (m *MsgCreateInvoice) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}
	if _, err := sdk.AccAddressFromBech32(m.Payee); err != nil {
		return fmt.Errorf("invalid payee address: %w", err)
	}
	if m.BookingID == "" {
		return fmt.Errorf("booking ID cannot be empty")
	}
	if !m.Amount.IsValid() || m.Amount.IsZero() {
		return ErrInvalidAmount
	}
	if len(m.Description) > 512 {
		return fmt.Errorf("description too long (max 512 chars)")
	}
	return nil
}

// --- MsgPayInvoice --- Pay an existing invoice
type MsgPayInvoice struct {
	Payer         string `protobuf:"bytes,1,opt,name=payer,proto3" json:"payer"`
	InvoiceID     string `protobuf:"bytes,2,opt,name=invoice_id,json=invoiceId,proto3" json:"invoice_id"`
	PaymentMethod string `protobuf:"bytes,3,opt,name=payment_method,json=paymentMethod,proto3" json:"payment_method"`
}

func (m *MsgPayInvoice) ProtoMessage()           {}
func (m *MsgPayInvoice) Reset()                  { *m = MsgPayInvoice{} }
func (m *MsgPayInvoice) String() string          { return fmt.Sprintf("pay_invoice: %s method=%s", m.InvoiceID, m.PaymentMethod) }
func (m *MsgPayInvoice) XXX_MessageName() string { return "syreen.payments.MsgPayInvoice" }

func (m *MsgPayInvoice) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Payer); err != nil {
		return fmt.Errorf("invalid payer address: %w", err)
	}
	if m.InvoiceID == "" {
		return fmt.Errorf("invoice ID cannot be empty")
	}
	switch PaymentMethod(m.PaymentMethod) {
	case PaymentMethodCrypto, PaymentMethodUSDC, PaymentMethodSYR:
	default:
		return fmt.Errorf("invalid payment method: %s (must be crypto, usdc, or syr)", m.PaymentMethod)
	}
	return nil
}

// --- MsgRefundPayment --- Refund a paid invoice (authority only)
type MsgRefundPayment struct {
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority"`
	InvoiceID string `protobuf:"bytes,2,opt,name=invoice_id,json=invoiceId,proto3" json:"invoice_id"`
	Reason    string `protobuf:"bytes,3,opt,name=reason,proto3" json:"reason"`
}

func (m *MsgRefundPayment) ProtoMessage()           {}
func (m *MsgRefundPayment) Reset()                  { *m = MsgRefundPayment{} }
func (m *MsgRefundPayment) String() string          { return fmt.Sprintf("refund_payment: %s", m.InvoiceID) }
func (m *MsgRefundPayment) XXX_MessageName() string { return "syreen.payments.MsgRefundPayment" }

func (m *MsgRefundPayment) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return fmt.Errorf("invalid authority address: %w", err)
	}
	if m.InvoiceID == "" {
		return fmt.Errorf("invoice ID cannot be empty")
	}
	if m.Reason == "" {
		return fmt.Errorf("refund reason cannot be empty")
	}
	if len(m.Reason) > 512 {
		return fmt.Errorf("reason too long (max 512 chars)")
	}
	return nil
}

// --- MsgSetExchangeRate --- Set exchange rate between two denoms (authority only)
type MsgSetExchangeRate struct {
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority"`
	FromDenom string `protobuf:"bytes,2,opt,name=from_denom,json=fromDenom,proto3" json:"from_denom"`
	ToDenom   string `protobuf:"bytes,3,opt,name=to_denom,json=toDenom,proto3" json:"to_denom"`
	Rate      string `protobuf:"bytes,4,opt,name=rate,proto3" json:"rate"`
}

func (m *MsgSetExchangeRate) ProtoMessage()           {}
func (m *MsgSetExchangeRate) Reset()                  { *m = MsgSetExchangeRate{} }
func (m *MsgSetExchangeRate) String() string          { return fmt.Sprintf("set_exchange_rate: %s/%s = %s", m.FromDenom, m.ToDenom, m.Rate) }
func (m *MsgSetExchangeRate) XXX_MessageName() string { return "syreen.payments.MsgSetExchangeRate" }

func (m *MsgSetExchangeRate) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return fmt.Errorf("invalid authority address: %w", err)
	}
	if m.FromDenom == "" {
		return fmt.Errorf("from_denom cannot be empty")
	}
	if m.ToDenom == "" {
		return fmt.Errorf("to_denom cannot be empty")
	}
	if m.FromDenom == m.ToDenom {
		return fmt.Errorf("from_denom and to_denom cannot be the same")
	}
	rate, err := math.LegacyNewDecFromStr(m.Rate)
	if err != nil {
		return fmt.Errorf("invalid rate: %w", err)
	}
	if !rate.IsPositive() {
		return fmt.Errorf("rate must be positive")
	}
	return nil
}

// --- MsgWithdrawEarnings --- Withdraw accumulated earnings
type MsgWithdrawEarnings struct {
	Address string   `protobuf:"bytes,1,opt,name=address,proto3" json:"address"`
	Amount  sdk.Coin `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount"`
}

func (m *MsgWithdrawEarnings) ProtoMessage()           {}
func (m *MsgWithdrawEarnings) Reset()                  { *m = MsgWithdrawEarnings{} }
func (m *MsgWithdrawEarnings) String() string          { return fmt.Sprintf("withdraw_earnings: %s amount=%s", m.Address, m.Amount) }
func (m *MsgWithdrawEarnings) XXX_MessageName() string { return "syreen.payments.MsgWithdrawEarnings" }

func (m *MsgWithdrawEarnings) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Address); err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}
	if !m.Amount.IsValid() || m.Amount.IsZero() {
		return ErrInvalidAmount
	}
	return nil
}

// --- MsgCancelInvoice --- Cancel a pending invoice (creator only)
type MsgCancelInvoice struct {
	Creator   string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator"`
	InvoiceID string `protobuf:"bytes,2,opt,name=invoice_id,json=invoiceId,proto3" json:"invoice_id"`
}

func (m *MsgCancelInvoice) ProtoMessage()           {}
func (m *MsgCancelInvoice) Reset()                  { *m = MsgCancelInvoice{} }
func (m *MsgCancelInvoice) String() string          { return fmt.Sprintf("cancel_invoice: %s", m.InvoiceID) }
func (m *MsgCancelInvoice) XXX_MessageName() string { return "syreen.payments.MsgCancelInvoice" }

func (m *MsgCancelInvoice) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}
	if m.InvoiceID == "" {
		return fmt.Errorf("invoice ID cannot be empty")
	}
	return nil
}

