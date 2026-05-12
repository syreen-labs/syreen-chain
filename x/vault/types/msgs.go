package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgCreateVault creates a new yield vault
type MsgCreateVault struct {
	Creator        string         `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator"`
	Name           string         `protobuf:"bytes,2,opt,name=name,proto3" json:"name"`
	DepositDenom   string         `protobuf:"bytes,3,opt,name=deposit_denom,json=depositDenom,proto3" json:"deposit_denom"`
	StrategyType   string         `protobuf:"bytes,4,opt,name=strategy_type,json=strategyType,proto3" json:"strategy_type"`
	TargetPoolIDs  []uint64       `protobuf:"bytes,5,opt,name=target_pool_ids,json=targetPoolIds,proto3" json:"target_pool_ids"`
	PerformanceFee math.LegacyDec `protobuf:"bytes,6,opt,name=performance_fee,json=performanceFee,proto3" json:"performance_fee"`
}

func (m *MsgCreateVault) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil { return err }
	if m.Name == "" { return ErrInvalidName }
	if m.DepositDenom == "" { return ErrInvalidDenom }
	if !ValidStrategy(m.StrategyType) { return ErrInvalidStrategy }
	if m.PerformanceFee.IsNegative() || m.PerformanceFee.GT(math.LegacyNewDecWithPrec(20, 2)) { return ErrInvalidFee }
	return nil
}
func (m *MsgCreateVault) ProtoMessage() {}
func (m *MsgCreateVault) Reset() { *m = MsgCreateVault{} }
func (m *MsgCreateVault) String() string { return fmt.Sprintf("MsgCreateVault{%s,%s}", m.Creator, m.Name) }
func (m *MsgCreateVault) XXX_MessageName() string { return "syreen.vault.MsgCreateVault" }

// MsgDepositVault deposits tokens into a vault
type MsgDepositVault struct {
	Sender  string   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	VaultID uint64   `protobuf:"varint,2,opt,name=vault_id,json=vaultId,proto3" json:"vault_id"`
	Amount  math.Int `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount"`
}

func (m *MsgDepositVault) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil { return err }
	if m.Amount.IsNil() || !m.Amount.IsPositive() { return ErrInvalidAmount }
	return nil
}
func (m *MsgDepositVault) ProtoMessage() {}
func (m *MsgDepositVault) Reset() { *m = MsgDepositVault{} }
func (m *MsgDepositVault) String() string { return fmt.Sprintf("MsgDepositVault{%s,%d,%s}", m.Sender, m.VaultID, m.Amount) }
func (m *MsgDepositVault) XXX_MessageName() string { return "syreen.vault.MsgDepositVault" }

// MsgWithdrawVault withdraws from a vault by burning shares
type MsgWithdrawVault struct {
	Sender  string   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	VaultID uint64   `protobuf:"varint,2,opt,name=vault_id,json=vaultId,proto3" json:"vault_id"`
	Shares  math.Int `protobuf:"bytes,3,opt,name=shares,proto3" json:"shares"`
}

func (m *MsgWithdrawVault) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil { return err }
	if m.Shares.IsNil() || !m.Shares.IsPositive() { return ErrInvalidAmount }
	return nil
}
func (m *MsgWithdrawVault) ProtoMessage() {}
func (m *MsgWithdrawVault) Reset() { *m = MsgWithdrawVault{} }
func (m *MsgWithdrawVault) String() string { return fmt.Sprintf("MsgWithdrawVault{%s,%d,%s}", m.Sender, m.VaultID, m.Shares) }
func (m *MsgWithdrawVault) XXX_MessageName() string { return "syreen.vault.MsgWithdrawVault" }

// MsgCompoundVault triggers manual compounding
type MsgCompoundVault struct {
	Sender  string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	VaultID uint64 `protobuf:"varint,2,opt,name=vault_id,json=vaultId,proto3" json:"vault_id"`
}

func (m *MsgCompoundVault) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	return err
}
func (m *MsgCompoundVault) ProtoMessage() {}
func (m *MsgCompoundVault) Reset() { *m = MsgCompoundVault{} }
func (m *MsgCompoundVault) String() string { return fmt.Sprintf("MsgCompoundVault{%s,%d}", m.Sender, m.VaultID) }
func (m *MsgCompoundVault) XXX_MessageName() string { return "syreen.vault.MsgCompoundVault" }

// MsgUpdateVaultStrategy updates vault strategy parameters
type MsgUpdateVaultStrategy struct {
	Creator       string   `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator"`
	VaultID       uint64   `protobuf:"varint,2,opt,name=vault_id,json=vaultId,proto3" json:"vault_id"`
	StrategyType  string   `protobuf:"bytes,3,opt,name=strategy_type,json=strategyType,proto3" json:"strategy_type"`
	TargetPoolIDs []uint64 `protobuf:"bytes,4,opt,name=target_pool_ids,json=targetPoolIds,proto3" json:"target_pool_ids"`
}

func (m *MsgUpdateVaultStrategy) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil { return err }
	if !ValidStrategy(m.StrategyType) { return ErrInvalidStrategy }
	return nil
}
func (m *MsgUpdateVaultStrategy) ProtoMessage() {}
func (m *MsgUpdateVaultStrategy) Reset() { *m = MsgUpdateVaultStrategy{} }
func (m *MsgUpdateVaultStrategy) String() string { return fmt.Sprintf("MsgUpdateVaultStrategy{%s,%d,%s}", m.Creator, m.VaultID, m.StrategyType) }
func (m *MsgUpdateVaultStrategy) XXX_MessageName() string { return "syreen.vault.MsgUpdateVaultStrategy" }

// Response types
type MsgCreateVaultResponse struct {
	VaultID uint64 `protobuf:"varint,1,opt,name=vault_id,json=vaultId,proto3" json:"vault_id"`
}
func (m *MsgCreateVaultResponse) ProtoMessage() {}
func (m *MsgCreateVaultResponse) Reset() { *m = MsgCreateVaultResponse{} }
func (m *MsgCreateVaultResponse) String() string { return "MsgCreateVaultResponse" }
func (m *MsgCreateVaultResponse) XXX_MessageName() string { return "syreen.vault.MsgCreateVaultResponse" }

type MsgDepositVaultResponse struct {
	SharesMinted math.Int `protobuf:"bytes,1,opt,name=shares_minted,json=sharesMinted,proto3" json:"shares_minted"`
}
func (m *MsgDepositVaultResponse) ProtoMessage() {}
func (m *MsgDepositVaultResponse) Reset() { *m = MsgDepositVaultResponse{} }
func (m *MsgDepositVaultResponse) String() string { return "MsgDepositVaultResponse" }
func (m *MsgDepositVaultResponse) XXX_MessageName() string { return "syreen.vault.MsgDepositVaultResponse" }

type MsgWithdrawVaultResponse struct {
	AmountReturned math.Int `protobuf:"bytes,1,opt,name=amount_returned,json=amountReturned,proto3" json:"amount_returned"`
}
func (m *MsgWithdrawVaultResponse) ProtoMessage() {}
func (m *MsgWithdrawVaultResponse) Reset() { *m = MsgWithdrawVaultResponse{} }
func (m *MsgWithdrawVaultResponse) String() string { return "MsgWithdrawVaultResponse" }
func (m *MsgWithdrawVaultResponse) XXX_MessageName() string { return "syreen.vault.MsgWithdrawVaultResponse" }

type MsgCompoundVaultResponse struct {
	YieldGenerated math.Int `protobuf:"bytes,1,opt,name=yield_generated,json=yieldGenerated,proto3" json:"yield_generated"`
}
func (m *MsgCompoundVaultResponse) ProtoMessage() {}
func (m *MsgCompoundVaultResponse) Reset() { *m = MsgCompoundVaultResponse{} }
func (m *MsgCompoundVaultResponse) String() string { return "MsgCompoundVaultResponse" }
func (m *MsgCompoundVaultResponse) XXX_MessageName() string { return "syreen.vault.MsgCompoundVaultResponse" }

type MsgUpdateVaultStrategyResponse struct{}
func (m *MsgUpdateVaultStrategyResponse) ProtoMessage() {}
func (m *MsgUpdateVaultStrategyResponse) Reset() { *m = MsgUpdateVaultStrategyResponse{} }
func (m *MsgUpdateVaultStrategyResponse) String() string { return "MsgUpdateVaultStrategyResponse" }
func (m *MsgUpdateVaultStrategyResponse) XXX_MessageName() string { return "syreen.vault.MsgUpdateVaultStrategyResponse" }
