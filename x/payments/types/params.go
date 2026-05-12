package types

import (
	"fmt"

	"cosmossdk.io/math"
)

var (
	DefaultPlatformFeePercent  = math.LegacyNewDecWithPrec(25, 1) // 2.5%
	DefaultInsuranceFeePercent = math.LegacyNewDecWithPrec(5, 1)  // 0.5%
	DefaultInvoiceExpiryBlocks uint64 = 345600                     // ~48h at ~500ms blocks
	DefaultEnablePayments      bool   = true
	DefaultPlatformAddress     string = ""
	DefaultInsurancePoolAddress string = ""
)

type Params struct {
	PlatformFeePercent   math.LegacyDec `protobuf:"bytes,1,opt,name=platform_fee_percent,json=platformFeePercent,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"platform_fee_percent"`
	InsuranceFeePercent  math.LegacyDec `protobuf:"bytes,2,opt,name=insurance_fee_percent,json=insuranceFeePercent,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"insurance_fee_percent"`
	InvoiceExpiryBlocks  uint64         `protobuf:"varint,3,opt,name=invoice_expiry_blocks,json=invoiceExpiryBlocks,proto3" json:"invoice_expiry_blocks"`
	EnablePayments       bool           `protobuf:"varint,4,opt,name=enable_payments,json=enablePayments,proto3" json:"enable_payments"`
	PlatformAddress      string         `protobuf:"bytes,5,opt,name=platform_address,json=platformAddress,proto3" json:"platform_address"`
	InsurancePoolAddress string         `protobuf:"bytes,6,opt,name=insurance_pool_address,json=insurancePoolAddress,proto3" json:"insurance_pool_address"`
}

func DefaultParams() Params {
	return Params{
		PlatformFeePercent:   DefaultPlatformFeePercent,
		InsuranceFeePercent:  DefaultInsuranceFeePercent,
		InvoiceExpiryBlocks:  DefaultInvoiceExpiryBlocks,
		EnablePayments:       DefaultEnablePayments,
		PlatformAddress:      DefaultPlatformAddress,
		InsurancePoolAddress: DefaultInsurancePoolAddress,
	}
}

func (m *Params) ProtoMessage()           {}
func (m *Params) Reset()                  { *m = Params{} }
func (m *Params) String() string {
	return fmt.Sprintf("params: platform_fee=%s%% insurance_fee=%s%% expiry=%d enabled=%t",
		m.PlatformFeePercent, m.InsuranceFeePercent, m.InvoiceExpiryBlocks, m.EnablePayments)
}
func (m *Params) XXX_MessageName() string { return "syreen.payments.Params" }

func (p Params) Validate() error {
	if p.PlatformFeePercent.IsNegative() {
		return fmt.Errorf("platform fee percent cannot be negative")
	}
	if p.PlatformFeePercent.GT(math.LegacyNewDec(100)) {
		return fmt.Errorf("platform fee percent cannot exceed 100")
	}
	if p.InsuranceFeePercent.IsNegative() {
		return fmt.Errorf("insurance fee percent cannot be negative")
	}
	if p.InsuranceFeePercent.GT(math.LegacyNewDec(100)) {
		return fmt.Errorf("insurance fee percent cannot exceed 100")
	}
	totalFee := p.PlatformFeePercent.Add(p.InsuranceFeePercent)
	if totalFee.GT(math.LegacyNewDec(100)) {
		return fmt.Errorf("combined fees cannot exceed 100%%")
	}
	if p.InvoiceExpiryBlocks == 0 {
		return fmt.Errorf("invoice expiry blocks must be greater than 0")
	}
	return nil
}
