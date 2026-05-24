package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountType defines the type of smart account
type AccountType string

const (
	AccountTypeMultiSig AccountType = "multisig"
	AccountTypeSocial   AccountType = "social"
	AccountTypeSession  AccountType = "session"
)

// SmartAccount represents a programmable smart wallet account
type SmartAccount struct {
	Address     string      `json:"address"`
	AccountType AccountType `json:"account_type"`
	Owners      []string    `json:"owners"`
	Threshold   uint32      `json:"threshold"`
	Metadata    []byte      `json:"metadata,omitempty"`
}

func (m *SmartAccount) ProtoMessage()           {}
func (m *SmartAccount) Reset()                  { *m = SmartAccount{} }
func (m *SmartAccount) String() string          { return fmt.Sprintf("SmartAccount{%s, type=%s, owners=%d}", m.Address, m.AccountType, len(m.Owners)) }
func (m *SmartAccount) XXX_MessageName() string { return "syreen.abstractaccount.SmartAccount" }

// Permission defines scoped permissions for a session key
type Permission struct {
	MsgType       string    `json:"msg_type"`
	MaxAmount     sdk.Coins `json:"max_amount,omitempty"`
	ContractAddr  string    `json:"contract_addr,omitempty"`
	AllowedDenoms []string  `json:"allowed_denoms,omitempty"`
}

func (m *Permission) ProtoMessage()           {}
func (m *Permission) Reset()                  { *m = Permission{} }
func (m *Permission) String() string          { return fmt.Sprintf("Permission{%s, max=%s}", m.MsgType, m.MaxAmount) }
func (m *Permission) XXX_MessageName() string { return "syreen.abstractaccount.Permission" }

// SessionKey represents a temporary key with scoped permissions
type SessionKey struct {
	Key         string       `json:"key"`
	Granter     string       `json:"granter"`
	Permissions []Permission `json:"permissions"`
	Expiry      time.Time    `json:"expiry"`
	ExpiresAt   int64        `json:"expires_at,omitempty"` // Block height expiry (0 = no block-height limit)
	SpendLimit  sdk.Coins    `json:"spend_limit,omitempty"`
	Used        sdk.Coins    `json:"used,omitempty"`
}

func (m *SessionKey) ProtoMessage()           {}
func (m *SessionKey) Reset()                  { *m = SessionKey{} }
func (m *SessionKey) String() string          { return fmt.Sprintf("SessionKey{%s, granter=%s, expiry=%s}", m.Key, m.Granter, m.Expiry) }
func (m *SessionKey) XXX_MessageName() string { return "syreen.abstractaccount.SessionKey" }

// RecoveryConfig defines the social recovery configuration for an account
type RecoveryConfig struct {
	Account     string        `json:"account"`
	Guardians   []string      `json:"guardians"`
	Threshold   uint32        `json:"threshold"`
	DelayPeriod time.Duration `json:"delay_period"`
}

func (m *RecoveryConfig) ProtoMessage()           {}
func (m *RecoveryConfig) Reset()                  { *m = RecoveryConfig{} }
func (m *RecoveryConfig) String() string          { return fmt.Sprintf("RecoveryConfig{%s, guardians=%d, threshold=%d}", m.Account, len(m.Guardians), m.Threshold) }
func (m *RecoveryConfig) XXX_MessageName() string { return "syreen.abstractaccount.RecoveryConfig" }

// RecoveryRequest represents an active recovery attempt
type RecoveryRequest struct {
	Account     string    `json:"account"`
	NewOwners   []string  `json:"new_owners"`
	Approvals   []string  `json:"approvals"`
	InitiatedAt time.Time `json:"initiated_at"`
}

func (m *RecoveryRequest) ProtoMessage()           {}
func (m *RecoveryRequest) Reset()                  { *m = RecoveryRequest{} }
func (m *RecoveryRequest) String() string          { return fmt.Sprintf("RecoveryRequest{%s, approvals=%d}", m.Account, len(m.Approvals)) }
func (m *RecoveryRequest) XXX_MessageName() string { return "syreen.abstractaccount.RecoveryRequest" }

// GasSponsor represents a gas sponsorship record
type GasSponsor struct {
	Sponsor        string    `json:"sponsor"`
	Sponsored      string    `json:"sponsored"`
	GasLimit       uint64    `json:"gas_limit"`
	Expiry         time.Time `json:"expiry"`
	TotalSponsored sdk.Coins `json:"total_sponsored,omitempty"` // Deprecated: kept for JSON compatibility
	TotalGasUsed   uint64    `json:"total_gas_used"`            // I-14: uint64 gas counter replaces fake "gas" denom
}

func (m *GasSponsor) ProtoMessage()           {}
func (m *GasSponsor) Reset()                  { *m = GasSponsor{} }
func (m *GasSponsor) String() string          { return fmt.Sprintf("GasSponsor{%s -> %s, limit=%d}", m.Sponsor, m.Sponsored, m.GasLimit) }
func (m *GasSponsor) XXX_MessageName() string { return "syreen.abstractaccount.GasSponsor" }
