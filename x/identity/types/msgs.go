package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgRegisterIdentity  = "register_identity"
	TypeMsgVerifyIdentity    = "verify_identity"
	TypeMsgRejectIdentity    = "reject_identity"
	TypeMsgRevokeIdentity    = "revoke_identity"
	TypeMsgUpdateIdentity    = "update_identity"
	TypeMsgRegisterVerifier  = "register_verifier"
	TypeMsgDeactivateVerifier = "deactivate_verifier"
	TypeMsgIncrementBookings = "increment_bookings"
)

// --- MsgRegisterIdentity --- User registers their identity for verification
type MsgRegisterIdentity struct {
	Address      string `protobuf:"bytes,1,opt,name=address,proto3" json:"address"`
	DocumentHash string `protobuf:"bytes,2,opt,name=document_hash,json=documentHash,proto3" json:"document_hash"`
	Nationality  string `protobuf:"bytes,3,opt,name=nationality,proto3" json:"nationality"`
	Level        string `protobuf:"bytes,4,opt,name=level,proto3" json:"level"`
}

func (m *MsgRegisterIdentity) ProtoMessage()           {}
func (m *MsgRegisterIdentity) Reset()                  { *m = MsgRegisterIdentity{} }
func (m *MsgRegisterIdentity) String() string          { return fmt.Sprintf("register_identity: %s level=%s", m.Address, m.Level) }
func (m *MsgRegisterIdentity) XXX_MessageName() string { return "syreen.identity.MsgRegisterIdentity" }

func (m *MsgRegisterIdentity) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}
	switch VerificationLevel(m.Level) {
	case VerificationBasic, VerificationStandard, VerificationEnhanced:
	default:
		return fmt.Errorf("invalid verification level: %s (must be basic, standard, or enhanced)", m.Level)
	}
	if m.Level != string(VerificationBasic) && m.DocumentHash == "" {
		return fmt.Errorf("document hash required for standard and enhanced verification")
	}
	if len(m.DocumentHash) > 128 {
		return fmt.Errorf("document hash too long (max 128 chars)")
	}
	if len(m.Nationality) > 64 {
		return fmt.Errorf("nationality too long (max 64 chars)")
	}
	return nil
}

// --- MsgVerifyIdentity --- Verifier approves an identity
type MsgVerifyIdentity struct {
	Verifier string `protobuf:"bytes,1,opt,name=verifier,proto3" json:"verifier"`
	Address  string `protobuf:"bytes,2,opt,name=address,proto3" json:"address"`
	Level    string `protobuf:"bytes,3,opt,name=level,proto3" json:"level"`
}

func (m *MsgVerifyIdentity) ProtoMessage()           {}
func (m *MsgVerifyIdentity) Reset()                  { *m = MsgVerifyIdentity{} }
func (m *MsgVerifyIdentity) String() string          { return fmt.Sprintf("verify_identity: verifier=%s address=%s", m.Verifier, m.Address) }
func (m *MsgVerifyIdentity) XXX_MessageName() string { return "syreen.identity.MsgVerifyIdentity" }

func (m *MsgVerifyIdentity) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Verifier); err != nil {
		return fmt.Errorf("invalid verifier address: %w", err)
	}
	if _, err := sdk.AccAddressFromBech32(m.Address); err != nil {
		return fmt.Errorf("invalid identity address: %w", err)
	}
	switch VerificationLevel(m.Level) {
	case VerificationBasic, VerificationStandard, VerificationEnhanced:
	default:
		return fmt.Errorf("invalid verification level: %s", m.Level)
	}
	return nil
}

// --- MsgRejectIdentity --- Verifier rejects an identity
type MsgRejectIdentity struct {
	Verifier string `protobuf:"bytes,1,opt,name=verifier,proto3" json:"verifier"`
	Address  string `protobuf:"bytes,2,opt,name=address,proto3" json:"address"`
	Reason   string `protobuf:"bytes,3,opt,name=reason,proto3" json:"reason"`
}

func (m *MsgRejectIdentity) ProtoMessage()           {}
func (m *MsgRejectIdentity) Reset()                  { *m = MsgRejectIdentity{} }
func (m *MsgRejectIdentity) String() string          { return fmt.Sprintf("reject_identity: %s", m.Address) }
func (m *MsgRejectIdentity) XXX_MessageName() string { return "syreen.identity.MsgRejectIdentity" }

func (m *MsgRejectIdentity) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Verifier); err != nil {
		return fmt.Errorf("invalid verifier address: %w", err)
	}
	if _, err := sdk.AccAddressFromBech32(m.Address); err != nil {
		return fmt.Errorf("invalid identity address: %w", err)
	}
	if m.Reason == "" {
		return fmt.Errorf("rejection reason cannot be empty")
	}
	if len(m.Reason) > 512 {
		return fmt.Errorf("rejection reason too long (max 512 chars)")
	}
	return nil
}

// --- MsgRevokeIdentity --- Verifier or governance revokes an identity
type MsgRevokeIdentity struct {
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority"`
	Address   string `protobuf:"bytes,2,opt,name=address,proto3" json:"address"`
	Reason    string `protobuf:"bytes,3,opt,name=reason,proto3" json:"reason"`
}

func (m *MsgRevokeIdentity) ProtoMessage()           {}
func (m *MsgRevokeIdentity) Reset()                  { *m = MsgRevokeIdentity{} }
func (m *MsgRevokeIdentity) String() string          { return fmt.Sprintf("revoke_identity: %s", m.Address) }
func (m *MsgRevokeIdentity) XXX_MessageName() string { return "syreen.identity.MsgRevokeIdentity" }

func (m *MsgRevokeIdentity) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return fmt.Errorf("invalid authority address: %w", err)
	}
	if _, err := sdk.AccAddressFromBech32(m.Address); err != nil {
		return fmt.Errorf("invalid identity address: %w", err)
	}
	if m.Reason == "" {
		return fmt.Errorf("revocation reason cannot be empty")
	}
	return nil
}

// --- MsgUpdateIdentity --- User updates their identity (e.g. new document)
type MsgUpdateIdentity struct {
	Address      string `protobuf:"bytes,1,opt,name=address,proto3" json:"address"`
	DocumentHash string `protobuf:"bytes,2,opt,name=document_hash,json=documentHash,proto3" json:"document_hash"`
	Nationality  string `protobuf:"bytes,3,opt,name=nationality,proto3" json:"nationality"`
}

func (m *MsgUpdateIdentity) ProtoMessage()           {}
func (m *MsgUpdateIdentity) Reset()                  { *m = MsgUpdateIdentity{} }
func (m *MsgUpdateIdentity) String() string          { return fmt.Sprintf("update_identity: %s", m.Address) }
func (m *MsgUpdateIdentity) XXX_MessageName() string { return "syreen.identity.MsgUpdateIdentity" }

func (m *MsgUpdateIdentity) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Address); err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}
	if len(m.DocumentHash) > 128 {
		return fmt.Errorf("document hash too long (max 128 chars)")
	}
	return nil
}

// --- MsgRegisterVerifier --- Governance registers a new verifier
type MsgRegisterVerifier struct {
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority"`
	Verifier  string `protobuf:"bytes,2,opt,name=verifier,proto3" json:"verifier"`
	Name      string `protobuf:"bytes,3,opt,name=name,proto3" json:"name"`
	MaxLevel  string `protobuf:"bytes,4,opt,name=max_level,json=maxLevel,proto3" json:"max_level"`
}

func (m *MsgRegisterVerifier) ProtoMessage()           {}
func (m *MsgRegisterVerifier) Reset()                  { *m = MsgRegisterVerifier{} }
func (m *MsgRegisterVerifier) String() string          { return fmt.Sprintf("register_verifier: %s", m.Verifier) }
func (m *MsgRegisterVerifier) XXX_MessageName() string { return "syreen.identity.MsgRegisterVerifier" }

func (m *MsgRegisterVerifier) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return fmt.Errorf("invalid authority address: %w", err)
	}
	if _, err := sdk.AccAddressFromBech32(m.Verifier); err != nil {
		return fmt.Errorf("invalid verifier address: %w", err)
	}
	if m.Name == "" {
		return fmt.Errorf("verifier name cannot be empty")
	}
	if len(m.Name) > 128 {
		return fmt.Errorf("verifier name too long (max 128 chars)")
	}
	switch VerificationLevel(m.MaxLevel) {
	case VerificationBasic, VerificationStandard, VerificationEnhanced:
	default:
		return fmt.Errorf("invalid max level: %s", m.MaxLevel)
	}
	return nil
}

// --- MsgDeactivateVerifier --- Governance deactivates a verifier
type MsgDeactivateVerifier struct {
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority"`
	Verifier  string `protobuf:"bytes,2,opt,name=verifier,proto3" json:"verifier"`
}

func (m *MsgDeactivateVerifier) ProtoMessage()           {}
func (m *MsgDeactivateVerifier) Reset()                  { *m = MsgDeactivateVerifier{} }
func (m *MsgDeactivateVerifier) String() string          { return fmt.Sprintf("deactivate_verifier: %s", m.Verifier) }
func (m *MsgDeactivateVerifier) XXX_MessageName() string { return "syreen.identity.MsgDeactivateVerifier" }

func (m *MsgDeactivateVerifier) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return fmt.Errorf("invalid authority address: %w", err)
	}
	if _, err := sdk.AccAddressFromBech32(m.Verifier); err != nil {
		return fmt.Errorf("invalid verifier address: %w", err)
	}
	return nil
}

// --- MsgIncrementBookings --- Called by travelescrow to track booking count
type MsgIncrementBookings struct {
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority"`
	Address   string `protobuf:"bytes,2,opt,name=address,proto3" json:"address"`
}

func (m *MsgIncrementBookings) ProtoMessage()           {}
func (m *MsgIncrementBookings) Reset()                  { *m = MsgIncrementBookings{} }
func (m *MsgIncrementBookings) String() string          { return fmt.Sprintf("increment_bookings: %s", m.Address) }
func (m *MsgIncrementBookings) XXX_MessageName() string { return "syreen.identity.MsgIncrementBookings" }

func (m *MsgIncrementBookings) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return fmt.Errorf("invalid authority address: %w", err)
	}
	if _, err := sdk.AccAddressFromBech32(m.Address); err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}
	return nil
}

