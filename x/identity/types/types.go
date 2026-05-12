package types

import "fmt"

// VerificationLevel represents the level of identity verification
type VerificationLevel string

const (
	VerificationNone     VerificationLevel = "none"
	VerificationBasic    VerificationLevel = "basic"    // Email verified
	VerificationStandard VerificationLevel = "standard" // ID document verified
	VerificationEnhanced VerificationLevel = "enhanced" // Biometric/full KYC
)

// VerificationStatus represents the state of a verification request
type VerificationStatus string

const (
	StatusPending  VerificationStatus = "pending"
	StatusVerified VerificationStatus = "verified"
	StatusRejected VerificationStatus = "rejected"
	StatusRevoked  VerificationStatus = "revoked"
	StatusExpired  VerificationStatus = "expired"
)

// Identity represents an on-chain identity record.
// No personal data is stored — only verification status and metadata hashes.
type Identity struct {
	Address           string             `protobuf:"bytes,1,opt,name=address,proto3" json:"address"`
	Level             VerificationLevel  `protobuf:"bytes,2,opt,name=level,proto3,casttype=VerificationLevel" json:"level"`
	Status            VerificationStatus `protobuf:"bytes,3,opt,name=status,proto3,casttype=VerificationStatus" json:"status"`
	Verifier          string             `protobuf:"bytes,4,opt,name=verifier,proto3" json:"verifier"`
	DocumentHash      string             `protobuf:"bytes,5,opt,name=document_hash,json=documentHash,proto3" json:"document_hash"`
	Nationality       string             `protobuf:"bytes,6,opt,name=nationality,proto3" json:"nationality"`
	VerifiedAtBlock   int64              `protobuf:"varint,7,opt,name=verified_at_block,json=verifiedAtBlock,proto3" json:"verified_at_block"`
	ExpiresAtBlock    int64              `protobuf:"varint,8,opt,name=expires_at_block,json=expiresAtBlock,proto3" json:"expires_at_block"`
	CreatedAtBlock    int64              `protobuf:"varint,9,opt,name=created_at_block,json=createdAtBlock,proto3" json:"created_at_block"`
	RejectionReason   string             `protobuf:"bytes,10,opt,name=rejection_reason,json=rejectionReason,proto3" json:"rejection_reason"`
	TotalBookings     uint64             `protobuf:"varint,11,opt,name=total_bookings,json=totalBookings,proto3" json:"total_bookings"`
	TrustScore        uint64             `protobuf:"varint,12,opt,name=trust_score,json=trustScore,proto3" json:"trust_score"`
}

func (m *Identity) ProtoMessage()           {}
func (m *Identity) Reset()                  { *m = Identity{} }
func (m *Identity) String() string          { return fmt.Sprintf("identity(%s): level=%s status=%s", m.Address, m.Level, m.Status) }
func (m *Identity) XXX_MessageName() string { return "syreen.identity.Identity" }

// Verifier represents an authorized identity verifier
type Verifier struct {
	Address       string `protobuf:"bytes,1,opt,name=address,proto3" json:"address"`
	Name          string `protobuf:"bytes,2,opt,name=name,proto3" json:"name"`
	MaxLevel      VerificationLevel `protobuf:"bytes,3,opt,name=max_level,json=maxLevel,proto3,casttype=VerificationLevel" json:"max_level"`
	Active        bool   `protobuf:"varint,4,opt,name=active,proto3" json:"active"`
	TotalVerified uint64 `protobuf:"varint,5,opt,name=total_verified,json=totalVerified,proto3" json:"total_verified"`
	RegisteredAt  int64  `protobuf:"varint,6,opt,name=registered_at,json=registeredAt,proto3" json:"registered_at"`
}

func (m *Verifier) ProtoMessage()           {}
func (m *Verifier) Reset()                  { *m = Verifier{} }
func (m *Verifier) String() string          { return fmt.Sprintf("verifier(%s): %s level=%s", m.Address, m.Name, m.MaxLevel) }
func (m *Verifier) XXX_MessageName() string { return "syreen.identity.Verifier" }
