package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgCreateLaunch creates a new token launch
type MsgCreateLaunch struct {
	Creator       string   `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator"`
	TokenDenom    string   `protobuf:"bytes,2,opt,name=token_denom,json=tokenDenom,proto3" json:"token_denom"`
	TokenSupply   math.Int `protobuf:"bytes,3,opt,name=token_supply,json=tokenSupply,proto3" json:"token_supply"`
	PricePerToken math.Int `protobuf:"bytes,4,opt,name=price_per_token,json=pricePerToken,proto3" json:"price_per_token"`
	QuoteDenom    string   `protobuf:"bytes,5,opt,name=quote_denom,json=quoteDenom,proto3" json:"quote_denom"`
	SoftCap       math.Int `protobuf:"bytes,6,opt,name=soft_cap,json=softCap,proto3" json:"soft_cap"`
	HardCap       math.Int `protobuf:"bytes,7,opt,name=hard_cap,json=hardCap,proto3" json:"hard_cap"`
	MaxPerWallet  math.Int `protobuf:"bytes,8,opt,name=max_per_wallet,json=maxPerWallet,proto3" json:"max_per_wallet"`
	StartBlock    int64    `protobuf:"varint,9,opt,name=start_block,json=startBlock,proto3" json:"start_block"`
	EndBlock      int64    `protobuf:"varint,10,opt,name=end_block,json=endBlock,proto3" json:"end_block"`
	VestingBlocks int64    `protobuf:"varint,11,opt,name=vesting_blocks,json=vestingBlocks,proto3" json:"vesting_blocks"`
	TGEPercent    uint64   `protobuf:"varint,12,opt,name=tge_percent,json=tgePercent,proto3" json:"tge_percent"`
}

func (m *MsgCreateLaunch) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil { return err }
	if m.TokenDenom == "" { return ErrInvalidLaunchParams }
	if m.TokenSupply.IsNil() || !m.TokenSupply.IsPositive() { return ErrInvalidAmount }
	if m.PricePerToken.IsNil() || !m.PricePerToken.IsPositive() { return ErrInvalidAmount }
	if m.QuoteDenom == "" { return ErrInvalidLaunchParams }
	if m.SoftCap.IsNil() || !m.SoftCap.IsPositive() { return ErrInvalidAmount }
	if m.HardCap.IsNil() || m.HardCap.LT(m.SoftCap) { return ErrInvalidLaunchParams }
	if m.MaxPerWallet.IsNil() || !m.MaxPerWallet.IsPositive() { return ErrInvalidAmount }
	if m.StartBlock <= 0 || m.EndBlock <= m.StartBlock { return ErrInvalidLaunchParams }
	if m.TGEPercent > 100 { return ErrInvalidLaunchParams }
	if m.VestingBlocks < 0 { return ErrInvalidLaunchParams }
	return nil
}
func (m *MsgCreateLaunch) ProtoMessage() {}
func (m *MsgCreateLaunch) Reset() { *m = MsgCreateLaunch{} }
func (m *MsgCreateLaunch) String() string { return fmt.Sprintf("MsgCreateLaunch{%s,%s}", m.Creator, m.TokenDenom) }
func (m *MsgCreateLaunch) XXX_MessageName() string { return "syreen.launchpad.MsgCreateLaunch" }

// MsgContribute contributes quote tokens to a launch
type MsgContribute struct {
	Sender   string   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	LaunchID uint64   `protobuf:"varint,2,opt,name=launch_id,json=launchId,proto3" json:"launch_id"`
	Amount   math.Int `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount"`
}

func (m *MsgContribute) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil { return err }
	if m.Amount.IsNil() || !m.Amount.IsPositive() { return ErrInvalidAmount }
	return nil
}
func (m *MsgContribute) ProtoMessage() {}
func (m *MsgContribute) Reset() { *m = MsgContribute{} }
func (m *MsgContribute) String() string { return fmt.Sprintf("MsgContribute{%s,%d,%s}", m.Sender, m.LaunchID, m.Amount) }
func (m *MsgContribute) XXX_MessageName() string { return "syreen.launchpad.MsgContribute" }

// MsgClaimTokens claims tokens after a successful launch
type MsgClaimTokens struct {
	Sender   string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	LaunchID uint64 `protobuf:"varint,2,opt,name=launch_id,json=launchId,proto3" json:"launch_id"`
}

func (m *MsgClaimTokens) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	return err
}
func (m *MsgClaimTokens) ProtoMessage() {}
func (m *MsgClaimTokens) Reset() { *m = MsgClaimTokens{} }
func (m *MsgClaimTokens) String() string { return fmt.Sprintf("MsgClaimTokens{%s,%d}", m.Sender, m.LaunchID) }
func (m *MsgClaimTokens) XXX_MessageName() string { return "syreen.launchpad.MsgClaimTokens" }

// MsgClaimRefund claims a refund after a failed launch
type MsgClaimRefund struct {
	Sender   string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	LaunchID uint64 `protobuf:"varint,2,opt,name=launch_id,json=launchId,proto3" json:"launch_id"`
}

func (m *MsgClaimRefund) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	return err
}
func (m *MsgClaimRefund) ProtoMessage() {}
func (m *MsgClaimRefund) Reset() { *m = MsgClaimRefund{} }
func (m *MsgClaimRefund) String() string { return fmt.Sprintf("MsgClaimRefund{%s,%d}", m.Sender, m.LaunchID) }
func (m *MsgClaimRefund) XXX_MessageName() string { return "syreen.launchpad.MsgClaimRefund" }

// MsgFinalizeLaunch finalizes a launch after it ends
type MsgFinalizeLaunch struct {
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority"`
	LaunchID  uint64 `protobuf:"varint,2,opt,name=launch_id,json=launchId,proto3" json:"launch_id"`
}

func (m *MsgFinalizeLaunch) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Authority)
	return err
}
func (m *MsgFinalizeLaunch) ProtoMessage() {}
func (m *MsgFinalizeLaunch) Reset() { *m = MsgFinalizeLaunch{} }
func (m *MsgFinalizeLaunch) String() string { return fmt.Sprintf("MsgFinalizeLaunch{%s,%d}", m.Authority, m.LaunchID) }
func (m *MsgFinalizeLaunch) XXX_MessageName() string { return "syreen.launchpad.MsgFinalizeLaunch" }

// Response types
type MsgCreateLaunchResponse struct {
	LaunchID uint64 `protobuf:"varint,1,opt,name=launch_id,json=launchId,proto3" json:"launch_id"`
}
func (m *MsgCreateLaunchResponse) ProtoMessage() {}
func (m *MsgCreateLaunchResponse) Reset() { *m = MsgCreateLaunchResponse{} }
func (m *MsgCreateLaunchResponse) String() string { return "MsgCreateLaunchResponse" }
func (m *MsgCreateLaunchResponse) XXX_MessageName() string { return "syreen.launchpad.MsgCreateLaunchResponse" }

type MsgContributeResponse struct{}
func (m *MsgContributeResponse) ProtoMessage() {}
func (m *MsgContributeResponse) Reset() { *m = MsgContributeResponse{} }
func (m *MsgContributeResponse) String() string { return "MsgContributeResponse" }
func (m *MsgContributeResponse) XXX_MessageName() string { return "syreen.launchpad.MsgContributeResponse" }

type MsgClaimTokensResponse struct {
	Amount math.Int `protobuf:"bytes,1,opt,name=amount,proto3" json:"amount"`
}
func (m *MsgClaimTokensResponse) ProtoMessage() {}
func (m *MsgClaimTokensResponse) Reset() { *m = MsgClaimTokensResponse{} }
func (m *MsgClaimTokensResponse) String() string { return "MsgClaimTokensResponse" }
func (m *MsgClaimTokensResponse) XXX_MessageName() string { return "syreen.launchpad.MsgClaimTokensResponse" }

type MsgClaimRefundResponse struct {
	Amount math.Int `protobuf:"bytes,1,opt,name=amount,proto3" json:"amount"`
}
func (m *MsgClaimRefundResponse) ProtoMessage() {}
func (m *MsgClaimRefundResponse) Reset() { *m = MsgClaimRefundResponse{} }
func (m *MsgClaimRefundResponse) String() string { return "MsgClaimRefundResponse" }
func (m *MsgClaimRefundResponse) XXX_MessageName() string { return "syreen.launchpad.MsgClaimRefundResponse" }

type MsgFinalizeLaunchResponse struct {
	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status"`
}
func (m *MsgFinalizeLaunchResponse) ProtoMessage() {}
func (m *MsgFinalizeLaunchResponse) Reset() { *m = MsgFinalizeLaunchResponse{} }
func (m *MsgFinalizeLaunchResponse) String() string { return "MsgFinalizeLaunchResponse" }
func (m *MsgFinalizeLaunchResponse) XXX_MessageName() string { return "syreen.launchpad.MsgFinalizeLaunchResponse" }
