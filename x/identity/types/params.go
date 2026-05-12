package types

import "fmt"

var (
	DefaultVerificationExpiryBlocks uint64 = 63115200 // ~1 year at 500ms blocks
	DefaultMinTrustScore            uint64 = 0
	DefaultEnableIdentity           bool   = true
	DefaultBasicVerificationFree    bool   = true
)

type Params struct {
	VerificationExpiryBlocks uint64 `protobuf:"varint,1,opt,name=verification_expiry_blocks,json=verificationExpiryBlocks,proto3" json:"verification_expiry_blocks"`
	MinTrustScore            uint64 `protobuf:"varint,2,opt,name=min_trust_score,json=minTrustScore,proto3" json:"min_trust_score"`
	EnableIdentity           bool   `protobuf:"varint,3,opt,name=enable_identity,json=enableIdentity,proto3" json:"enable_identity"`
	BasicVerificationFree    bool   `protobuf:"varint,4,opt,name=basic_verification_free,json=basicVerificationFree,proto3" json:"basic_verification_free"`
}

func DefaultParams() Params {
	return Params{
		VerificationExpiryBlocks: DefaultVerificationExpiryBlocks,
		MinTrustScore:            DefaultMinTrustScore,
		EnableIdentity:           DefaultEnableIdentity,
		BasicVerificationFree:    DefaultBasicVerificationFree,
	}
}

func (m *Params) ProtoMessage()           {}
func (m *Params) Reset()                  { *m = Params{} }
func (m *Params) String() string          { return fmt.Sprintf("params: expiry=%d trust=%d enabled=%t", m.VerificationExpiryBlocks, m.MinTrustScore, m.EnableIdentity) }
func (m *Params) XXX_MessageName() string { return "syreen.identity.Params" }

func (p Params) Validate() error {
	if p.VerificationExpiryBlocks == 0 {
		return fmt.Errorf("verification expiry blocks must be greater than 0")
	}
	return nil
}
