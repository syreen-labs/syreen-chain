package types

import (
	"encoding/json"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgCreateSmartAccount = "create_smart_account"
	TypeMsgCreateSessionKey   = "create_session_key"
	TypeMsgRevokeSessionKey   = "revoke_session_key"
	TypeMsgInitiateRecovery   = "initiate_recovery"
	TypeMsgApproveRecovery    = "approve_recovery"
	TypeMsgExecuteRecovery    = "execute_recovery"
	TypeMsgSponsorGas         = "sponsor_gas"
	TypeMsgBatchExecute       = "batch_execute"
)

// --- MsgCreateSmartAccount ---

type MsgCreateSmartAccount struct {
	Sender      string      `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	AccountType AccountType `protobuf:"bytes,2,opt,name=account_type,json=accountType,proto3,casttype=AccountType" json:"account_type"`
	Owners      []string    `protobuf:"bytes,3,rep,name=owners,proto3" json:"owners"`
	Threshold   uint32      `protobuf:"varint,4,opt,name=threshold,proto3" json:"threshold"`
}

func (m *MsgCreateSmartAccount) ProtoMessage()           {}
func (m *MsgCreateSmartAccount) Reset()                  { *m = MsgCreateSmartAccount{} }
func (m *MsgCreateSmartAccount) String() string          { return fmt.Sprintf("create_smart_account: sender=%s type=%s", m.Sender, m.AccountType) }
func (m *MsgCreateSmartAccount) XXX_MessageName() string { return "syreen.abstractaccount.MsgCreateSmartAccount" }

func (m *MsgCreateSmartAccount) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return ErrInvalidAddress
	}
	if m.AccountType != AccountTypeMultiSig && m.AccountType != AccountTypeSocial && m.AccountType != AccountTypeSession {
		return ErrInvalidAccountType
	}
	if len(m.Owners) == 0 {
		return fmt.Errorf("at least one owner is required")
	}
	for _, owner := range m.Owners {
		if _, err := sdk.AccAddressFromBech32(owner); err != nil {
			return ErrInvalidAddress
		}
	}
	if m.Threshold == 0 || m.Threshold > uint32(len(m.Owners)) {
		return ErrInvalidThreshold
	}
	return nil
}

// --- MsgCreateSessionKey ---

type MsgCreateSessionKey struct {
	Granter     string       `protobuf:"bytes,1,opt,name=granter,proto3" json:"granter"`
	Grantee     string       `protobuf:"bytes,2,opt,name=grantee,proto3" json:"grantee"`
	Permissions []Permission `protobuf:"bytes,3,rep,name=permissions,proto3" json:"permissions"`
	Duration    time.Duration `protobuf:"varint,4,opt,name=duration,proto3,casttype=time.Duration" json:"duration"`
}

func (m *MsgCreateSessionKey) ProtoMessage()           {}
func (m *MsgCreateSessionKey) Reset()                  { *m = MsgCreateSessionKey{} }
func (m *MsgCreateSessionKey) String() string          { return fmt.Sprintf("create_session_key: granter=%s grantee=%s", m.Granter, m.Grantee) }
func (m *MsgCreateSessionKey) XXX_MessageName() string { return "syreen.abstractaccount.MsgCreateSessionKey" }

func (m *MsgCreateSessionKey) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Granter); err != nil {
		return ErrInvalidAddress
	}
	if _, err := sdk.AccAddressFromBech32(m.Grantee); err != nil {
		return ErrInvalidAddress
	}
	if m.Duration <= 0 {
		return fmt.Errorf("duration must be positive")
	}
	if len(m.Permissions) == 0 {
		return fmt.Errorf("at least one permission is required")
	}
	return nil
}

// --- MsgRevokeSessionKey ---

type MsgRevokeSessionKey struct {
	Granter        string `protobuf:"bytes,1,opt,name=granter,proto3" json:"granter"`
	SessionKeyAddr string `protobuf:"bytes,2,opt,name=session_key_addr,json=sessionKeyAddr,proto3" json:"session_key_addr"`
}

func (m *MsgRevokeSessionKey) ProtoMessage()           {}
func (m *MsgRevokeSessionKey) Reset()                  { *m = MsgRevokeSessionKey{} }
func (m *MsgRevokeSessionKey) String() string          { return fmt.Sprintf("revoke_session_key: granter=%s key=%s", m.Granter, m.SessionKeyAddr) }
func (m *MsgRevokeSessionKey) XXX_MessageName() string { return "syreen.abstractaccount.MsgRevokeSessionKey" }

func (m *MsgRevokeSessionKey) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Granter); err != nil {
		return ErrInvalidAddress
	}
	if _, err := sdk.AccAddressFromBech32(m.SessionKeyAddr); err != nil {
		return ErrInvalidAddress
	}
	return nil
}

// --- MsgInitiateRecovery ---

type MsgInitiateRecovery struct {
	Guardian  string   `protobuf:"bytes,1,opt,name=guardian,proto3" json:"guardian"`
	Account   string   `protobuf:"bytes,2,opt,name=account,proto3" json:"account"`
	NewOwners []string `protobuf:"bytes,3,rep,name=new_owners,json=newOwners,proto3" json:"new_owners"`
}

func (m *MsgInitiateRecovery) ProtoMessage()           {}
func (m *MsgInitiateRecovery) Reset()                  { *m = MsgInitiateRecovery{} }
func (m *MsgInitiateRecovery) String() string          { return fmt.Sprintf("initiate_recovery: guardian=%s account=%s", m.Guardian, m.Account) }
func (m *MsgInitiateRecovery) XXX_MessageName() string { return "syreen.abstractaccount.MsgInitiateRecovery" }

func (m *MsgInitiateRecovery) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Guardian); err != nil {
		return ErrInvalidAddress
	}
	if _, err := sdk.AccAddressFromBech32(m.Account); err != nil {
		return ErrInvalidAddress
	}
	if len(m.NewOwners) == 0 {
		return fmt.Errorf("at least one new owner is required")
	}
	for _, owner := range m.NewOwners {
		if _, err := sdk.AccAddressFromBech32(owner); err != nil {
			return ErrInvalidAddress
		}
	}
	return nil
}

// --- MsgApproveRecovery ---

type MsgApproveRecovery struct {
	Guardian string `protobuf:"bytes,1,opt,name=guardian,proto3" json:"guardian"`
	Account  string `protobuf:"bytes,2,opt,name=account,proto3" json:"account"`
}

func (m *MsgApproveRecovery) ProtoMessage()           {}
func (m *MsgApproveRecovery) Reset()                  { *m = MsgApproveRecovery{} }
func (m *MsgApproveRecovery) String() string          { return fmt.Sprintf("approve_recovery: guardian=%s account=%s", m.Guardian, m.Account) }
func (m *MsgApproveRecovery) XXX_MessageName() string { return "syreen.abstractaccount.MsgApproveRecovery" }

func (m *MsgApproveRecovery) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Guardian); err != nil {
		return ErrInvalidAddress
	}
	if _, err := sdk.AccAddressFromBech32(m.Account); err != nil {
		return ErrInvalidAddress
	}
	return nil
}

// --- MsgExecuteRecovery ---

type MsgExecuteRecovery struct {
	Sender  string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`  // Must be one of the approving guardians
	Account string `protobuf:"bytes,2,opt,name=account,proto3" json:"account"` // The smart account being recovered
}

func (m *MsgExecuteRecovery) ProtoMessage()           {}
func (m *MsgExecuteRecovery) Reset()                  { *m = MsgExecuteRecovery{} }
func (m *MsgExecuteRecovery) String() string          { return fmt.Sprintf("execute_recovery: sender=%s account=%s", m.Sender, m.Account) }
func (m *MsgExecuteRecovery) XXX_MessageName() string { return "syreen.abstractaccount.MsgExecuteRecovery" }

func (m *MsgExecuteRecovery) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return fmt.Errorf("invalid sender address: %w", err)
	}
	if _, err := sdk.AccAddressFromBech32(m.Account); err != nil {
		return ErrInvalidAddress
	}
	return nil
}

// --- MsgSponsorGas ---

type MsgSponsorGas struct {
	Sponsor   string        `protobuf:"bytes,1,opt,name=sponsor,proto3" json:"sponsor"`
	Sponsored string        `protobuf:"bytes,2,opt,name=sponsored,proto3" json:"sponsored"`
	GasLimit  uint64        `protobuf:"varint,3,opt,name=gas_limit,json=gasLimit,proto3" json:"gas_limit"`
	Duration  time.Duration `protobuf:"varint,4,opt,name=duration,proto3,casttype=time.Duration" json:"duration"`
}

func (m *MsgSponsorGas) ProtoMessage()           {}
func (m *MsgSponsorGas) Reset()                  { *m = MsgSponsorGas{} }
func (m *MsgSponsorGas) String() string          { return fmt.Sprintf("sponsor_gas: %s -> %s limit=%d", m.Sponsor, m.Sponsored, m.GasLimit) }
func (m *MsgSponsorGas) XXX_MessageName() string { return "syreen.abstractaccount.MsgSponsorGas" }

func (m *MsgSponsorGas) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sponsor); err != nil {
		return ErrInvalidAddress
	}
	if _, err := sdk.AccAddressFromBech32(m.Sponsored); err != nil {
		return ErrInvalidAddress
	}
	if m.GasLimit == 0 {
		return fmt.Errorf("gas limit must be positive")
	}
	if m.Duration <= 0 {
		return fmt.Errorf("duration must be positive")
	}
	return nil
}

// --- MsgBatchExecute ---

type MsgBatchExecute struct {
	Sender   string             `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	Messages []json.RawMessage  `protobuf:"bytes,2,rep,name=messages,proto3" json:"messages"`
}

func (m *MsgBatchExecute) ProtoMessage()           {}
func (m *MsgBatchExecute) Reset()                  { *m = MsgBatchExecute{} }
func (m *MsgBatchExecute) String() string          { return fmt.Sprintf("batch_execute: sender=%s msgs=%d", m.Sender, len(m.Messages)) }
func (m *MsgBatchExecute) XXX_MessageName() string { return "syreen.abstractaccount.MsgBatchExecute" }

func (m *MsgBatchExecute) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return ErrInvalidAddress
	}
	if len(m.Messages) == 0 {
		return fmt.Errorf("at least one message is required")
	}
	return nil
}
