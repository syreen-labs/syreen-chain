package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InvoiceStatus represents the state of an invoice
type InvoiceStatus string

const (
	InvoiceStatusPending   InvoiceStatus = "pending"
	InvoiceStatusPaid      InvoiceStatus = "paid"
	InvoiceStatusRefunded  InvoiceStatus = "refunded"
	InvoiceStatusExpired   InvoiceStatus = "expired"
	InvoiceStatusCancelled InvoiceStatus = "cancelled"
)

// PaymentMethod represents accepted payment methods
type PaymentMethod string

const (
	PaymentMethodCrypto PaymentMethod = "crypto"
	PaymentMethodUSDC   PaymentMethod = "usdc"
	PaymentMethodSYR    PaymentMethod = "syr"
)

// Invoice represents a payment invoice for a booking
type Invoice struct {
	ID             string        `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	BookingID      string        `protobuf:"bytes,2,opt,name=booking_id,json=bookingId,proto3" json:"booking_id"`
	Payer          string        `protobuf:"bytes,3,opt,name=payer,proto3" json:"payer"`
	Payee          string        `protobuf:"bytes,4,opt,name=payee,proto3" json:"payee"`
	Amount         sdk.Coin      `protobuf:"bytes,5,opt,name=amount,proto3" json:"amount"`
	PlatformFee    sdk.Coin      `protobuf:"bytes,6,opt,name=platform_fee,json=platformFee,proto3" json:"platform_fee"`
	InsuranceFee   sdk.Coin      `protobuf:"bytes,7,opt,name=insurance_fee,json=insuranceFee,proto3" json:"insurance_fee"`
	NetAmount      sdk.Coin      `protobuf:"bytes,8,opt,name=net_amount,json=netAmount,proto3" json:"net_amount"`
	Status         InvoiceStatus `protobuf:"bytes,9,opt,name=status,proto3,casttype=InvoiceStatus" json:"status"`
	PaymentMethod  string        `protobuf:"bytes,10,opt,name=payment_method,json=paymentMethod,proto3" json:"payment_method"`
	PaymentTxHash  string        `protobuf:"bytes,11,opt,name=payment_tx_hash,json=paymentTxHash,proto3" json:"payment_tx_hash"`
	CreatedAtBlock int64         `protobuf:"varint,12,opt,name=created_at_block,json=createdAtBlock,proto3" json:"created_at_block"`
	PaidAtBlock    int64         `protobuf:"varint,13,opt,name=paid_at_block,json=paidAtBlock,proto3" json:"paid_at_block"`
	ExpiresAtBlock int64         `protobuf:"varint,14,opt,name=expires_at_block,json=expiresAtBlock,proto3" json:"expires_at_block"`
	Description    string        `protobuf:"bytes,15,opt,name=description,proto3" json:"description"`
	Currency       string        `protobuf:"bytes,16,opt,name=currency,proto3" json:"currency"`
	Creator        string        `protobuf:"bytes,17,opt,name=creator,proto3" json:"creator"`
}

func (m *Invoice) ProtoMessage()           {}
func (m *Invoice) Reset()                  { *m = Invoice{} }
func (m *Invoice) String() string          { return fmt.Sprintf("invoice(%s): booking=%s status=%s amount=%s", m.ID, m.BookingID, m.Status, m.Amount) }
func (m *Invoice) XXX_MessageName() string { return "syreen.payments.Invoice" }

// ExchangeRate represents a rate between two denominations
type ExchangeRate struct {
	FromDenom      string         `protobuf:"bytes,1,opt,name=from_denom,json=fromDenom,proto3" json:"from_denom"`
	ToDenom        string         `protobuf:"bytes,2,opt,name=to_denom,json=toDenom,proto3" json:"to_denom"`
	Rate           math.LegacyDec `protobuf:"bytes,3,opt,name=rate,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"rate"`
	UpdatedAtBlock int64          `protobuf:"varint,4,opt,name=updated_at_block,json=updatedAtBlock,proto3" json:"updated_at_block"`
	UpdatedBy      string         `protobuf:"bytes,5,opt,name=updated_by,json=updatedBy,proto3" json:"updated_by"`
}

func (m *ExchangeRate) ProtoMessage()           {}
func (m *ExchangeRate) Reset()                  { *m = ExchangeRate{} }
func (m *ExchangeRate) String() string          { return fmt.Sprintf("rate(%s/%s): %s", m.FromDenom, m.ToDenom, m.Rate) }
func (m *ExchangeRate) XXX_MessageName() string { return "syreen.payments.ExchangeRate" }

// Earnings represents accumulated earnings for an address
type Earnings struct {
	Address        string    `protobuf:"bytes,1,opt,name=address,proto3" json:"address"`
	Available      sdk.Coins `protobuf:"bytes,2,rep,name=available,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"available"`
	TotalEarned    sdk.Coins `protobuf:"bytes,3,rep,name=total_earned,json=totalEarned,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"total_earned"`
	TotalWithdrawn sdk.Coins `protobuf:"bytes,4,rep,name=total_withdrawn,json=totalWithdrawn,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"total_withdrawn"`
}

func (m *Earnings) ProtoMessage()           {}
func (m *Earnings) Reset()                  { *m = Earnings{} }
func (m *Earnings) String() string          { return fmt.Sprintf("earnings(%s): available=%s", m.Address, m.Available) }
func (m *Earnings) XXX_MessageName() string { return "syreen.payments.Earnings" }
