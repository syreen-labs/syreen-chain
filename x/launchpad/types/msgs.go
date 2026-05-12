package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgCreateLaunch creates a new token launch
type MsgCreateLaunch struct {
	Creator       string   `json:"creator"`
	TokenDenom    string   `json:"token_denom"`
	TokenSupply   math.Int `json:"token_supply"`
	PricePerToken math.Int `json:"price_per_token"`
	QuoteDenom    string   `json:"quote_denom"`
	SoftCap       math.Int `json:"soft_cap"`
	HardCap       math.Int `json:"hard_cap"`
	MaxPerWallet  math.Int `json:"max_per_wallet"`
	StartBlock    int64    `json:"start_block"`
	EndBlock      int64    `json:"end_block"`
	VestingBlocks int64    `json:"vesting_blocks"`
	TGEPercent    uint64   `json:"tge_percent"`
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
	Sender   string   `json:"sender"`
	LaunchID uint64   `json:"launch_id"`
	Amount   math.Int `json:"amount"`
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
	Sender   string `json:"sender"`
	LaunchID uint64 `json:"launch_id"`
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
	Sender   string `json:"sender"`
	LaunchID uint64 `json:"launch_id"`
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
	Authority string `json:"authority"`
	LaunchID  uint64 `json:"launch_id"`
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
type MsgCreateLaunchResponse struct{ LaunchID uint64 `json:"launch_id"` }
func (m *MsgCreateLaunchResponse) ProtoMessage() {}; func (m *MsgCreateLaunchResponse) Reset() { *m = MsgCreateLaunchResponse{} }; func (m *MsgCreateLaunchResponse) String() string { return "MsgCreateLaunchResponse" }; func (m *MsgCreateLaunchResponse) XXX_MessageName() string { return "syreen.launchpad.MsgCreateLaunchResponse" }

type MsgContributeResponse struct{}
func (m *MsgContributeResponse) ProtoMessage() {}; func (m *MsgContributeResponse) Reset() { *m = MsgContributeResponse{} }; func (m *MsgContributeResponse) String() string { return "MsgContributeResponse" }; func (m *MsgContributeResponse) XXX_MessageName() string { return "syreen.launchpad.MsgContributeResponse" }

type MsgClaimTokensResponse struct{ Amount math.Int `json:"amount"` }
func (m *MsgClaimTokensResponse) ProtoMessage() {}; func (m *MsgClaimTokensResponse) Reset() { *m = MsgClaimTokensResponse{} }; func (m *MsgClaimTokensResponse) String() string { return "MsgClaimTokensResponse" }; func (m *MsgClaimTokensResponse) XXX_MessageName() string { return "syreen.launchpad.MsgClaimTokensResponse" }

type MsgClaimRefundResponse struct{ Amount math.Int `json:"amount"` }
func (m *MsgClaimRefundResponse) ProtoMessage() {}; func (m *MsgClaimRefundResponse) Reset() { *m = MsgClaimRefundResponse{} }; func (m *MsgClaimRefundResponse) String() string { return "MsgClaimRefundResponse" }; func (m *MsgClaimRefundResponse) XXX_MessageName() string { return "syreen.launchpad.MsgClaimRefundResponse" }

type MsgFinalizeLaunchResponse struct{ Status string `json:"status"` }
func (m *MsgFinalizeLaunchResponse) ProtoMessage() {}; func (m *MsgFinalizeLaunchResponse) Reset() { *m = MsgFinalizeLaunchResponse{} }; func (m *MsgFinalizeLaunchResponse) String() string { return "MsgFinalizeLaunchResponse" }; func (m *MsgFinalizeLaunchResponse) XXX_MessageName() string { return "syreen.launchpad.MsgFinalizeLaunchResponse" }
