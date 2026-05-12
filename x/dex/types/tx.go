package types

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgServer defines the dex msg server interface.
type MsgServer interface {
	CreatePool(context.Context, *MsgCreatePool) (*MsgCreatePoolResponse, error)
	AddLiquidity(context.Context, *MsgAddLiquidity) (*MsgAddLiquidityResponse, error)
	RemoveLiquidity(context.Context, *MsgRemoveLiquidity) (*MsgRemoveLiquidityResponse, error)
	Swap(context.Context, *MsgSwap) (*MsgSwapResponse, error)
	CreateReferralCode(context.Context, *MsgCreateReferralCode) (*MsgCreateReferralCodeResponse, error)
	RegisterReferral(context.Context, *MsgRegisterReferral) (*MsgRegisterReferralResponse, error)
	PlaceOrder(context.Context, *MsgPlaceOrder) (*MsgPlaceOrderResponse, error)
	CancelOrder(context.Context, *MsgCancelOrder) (*MsgCancelOrderResponse, error)
	ModifyOrder(context.Context, *MsgModifyOrder) (*MsgModifyOrderResponse, error)
	FollowTrader(context.Context, *MsgFollowTrader) (*MsgFollowTraderResponse, error)
	UnfollowTrader(context.Context, *MsgUnfollowTrader) (*MsgUnfollowTraderResponse, error)
	UpdateCopySettings(context.Context, *MsgUpdateCopySettings) (*MsgUpdateCopySettingsResponse, error)
	ClaimReferralRewards(context.Context, *MsgClaimReferralRewards) (*MsgClaimReferralRewardsResponse, error)
	MultiHopSwap(context.Context, *MsgMultiHopSwap) (*MsgMultiHopSwapResponse, error)
	SetPoolFeeConfig(context.Context, *MsgSetPoolFeeConfig) (*MsgSetPoolFeeConfigResponse, error)
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

type MsgCreatePoolResponse struct {
	PoolID uint64 `protobuf:"varint,1,opt,name=pool_id,json=poolId,proto3" json:"pool_id"`
}

func (m *MsgCreatePoolResponse) ProtoMessage()          {}
func (m *MsgCreatePoolResponse) Reset()                 { *m = MsgCreatePoolResponse{} }
func (m *MsgCreatePoolResponse) String() string         { return "create_pool_response" }
func (m *MsgCreatePoolResponse) XXX_MessageName() string { return "syreen.dex.MsgCreatePoolResponse" }

type MsgAddLiquidityResponse struct {
	SharesMinted math.Int `protobuf:"bytes,1,opt,name=shares_minted,json=sharesMinted,proto3" json:"shares_minted"`
}

func (m *MsgAddLiquidityResponse) ProtoMessage()          {}
func (m *MsgAddLiquidityResponse) Reset()                 { *m = MsgAddLiquidityResponse{} }
func (m *MsgAddLiquidityResponse) String() string         { return "add_liquidity_response" }
func (m *MsgAddLiquidityResponse) XXX_MessageName() string { return "syreen.dex.MsgAddLiquidityResponse" }

type MsgRemoveLiquidityResponse struct {
	AmountA math.Int `protobuf:"bytes,1,opt,name=amount_a,json=amountA,proto3" json:"amount_a"`
	AmountB math.Int `protobuf:"bytes,2,opt,name=amount_b,json=amountB,proto3" json:"amount_b"`
}

func (m *MsgRemoveLiquidityResponse) ProtoMessage()          {}
func (m *MsgRemoveLiquidityResponse) Reset()                 { *m = MsgRemoveLiquidityResponse{} }
func (m *MsgRemoveLiquidityResponse) String() string         { return "remove_liquidity_response" }
func (m *MsgRemoveLiquidityResponse) XXX_MessageName() string { return "syreen.dex.MsgRemoveLiquidityResponse" }

type MsgSwapResponse struct {
	TokenOut sdk.Coin `protobuf:"bytes,1,opt,name=token_out,json=tokenOut,proto3" json:"token_out"`
}

func (m *MsgSwapResponse) ProtoMessage()          {}
func (m *MsgSwapResponse) Reset()                 { *m = MsgSwapResponse{} }
func (m *MsgSwapResponse) String() string         { return "swap_response" }
func (m *MsgSwapResponse) XXX_MessageName() string { return "syreen.dex.MsgSwapResponse" }

// ---------------------------------------------------------------------------
// Referral response types
// ---------------------------------------------------------------------------

type MsgCreateReferralCodeResponse struct {
	ReferralCode string `protobuf:"bytes,1,opt,name=referral_code,json=referralCode,proto3" json:"referral_code"`
}

func (m *MsgCreateReferralCodeResponse) ProtoMessage()          {}
func (m *MsgCreateReferralCodeResponse) Reset()                 { *m = MsgCreateReferralCodeResponse{} }
func (m *MsgCreateReferralCodeResponse) String() string         { return "create_referral_code_response" }
func (m *MsgCreateReferralCodeResponse) XXX_MessageName() string { return "syreen.dex.MsgCreateReferralCodeResponse" }

type MsgRegisterReferralResponse struct {
	Referrer string `protobuf:"bytes,1,opt,name=referrer,proto3" json:"referrer"`
}

func (m *MsgRegisterReferralResponse) ProtoMessage()          {}
func (m *MsgRegisterReferralResponse) Reset()                 { *m = MsgRegisterReferralResponse{} }
func (m *MsgRegisterReferralResponse) String() string         { return "register_referral_response" }
func (m *MsgRegisterReferralResponse) XXX_MessageName() string { return "syreen.dex.MsgRegisterReferralResponse" }

// ---------------------------------------------------------------------------
// Order Book response types
// ---------------------------------------------------------------------------

type MsgPlaceOrderResponse struct {
	OrderID uint64 `protobuf:"varint,1,opt,name=order_id,json=orderId,proto3" json:"order_id"`
	Status  string `protobuf:"bytes,2,opt,name=status,proto3" json:"status"`
}

func (m *MsgPlaceOrderResponse) ProtoMessage()          {}
func (m *MsgPlaceOrderResponse) Reset()                 { *m = MsgPlaceOrderResponse{} }
func (m *MsgPlaceOrderResponse) String() string         { return "place_order_response" }
func (m *MsgPlaceOrderResponse) XXX_MessageName() string { return "syreen.dex.MsgPlaceOrderResponse" }

type MsgCancelOrderResponse struct{}

func (m *MsgCancelOrderResponse) ProtoMessage()          {}
func (m *MsgCancelOrderResponse) Reset()                 { *m = MsgCancelOrderResponse{} }
func (m *MsgCancelOrderResponse) String() string         { return "cancel_order_response" }
func (m *MsgCancelOrderResponse) XXX_MessageName() string { return "syreen.dex.MsgCancelOrderResponse" }

type MsgModifyOrderResponse struct {
	OrderID uint64 `protobuf:"varint,1,opt,name=order_id,json=orderId,proto3" json:"order_id"`
}

func (m *MsgModifyOrderResponse) ProtoMessage()          {}
func (m *MsgModifyOrderResponse) Reset()                 { *m = MsgModifyOrderResponse{} }
func (m *MsgModifyOrderResponse) String() string         { return "modify_order_response" }
func (m *MsgModifyOrderResponse) XXX_MessageName() string { return "syreen.dex.MsgModifyOrderResponse" }

// ---------------------------------------------------------------------------
// Copy Trading response types
// ---------------------------------------------------------------------------

type MsgFollowTraderResponse struct{}

func (m *MsgFollowTraderResponse) ProtoMessage()           {}
func (m *MsgFollowTraderResponse) Reset()                  { *m = MsgFollowTraderResponse{} }
func (m *MsgFollowTraderResponse) String() string          { return "follow_trader_response" }
func (m *MsgFollowTraderResponse) XXX_MessageName() string { return "syreen.dex.MsgFollowTraderResponse" }

type MsgUnfollowTraderResponse struct{}

func (m *MsgUnfollowTraderResponse) ProtoMessage()           {}
func (m *MsgUnfollowTraderResponse) Reset()                  { *m = MsgUnfollowTraderResponse{} }
func (m *MsgUnfollowTraderResponse) String() string          { return "unfollow_trader_response" }
func (m *MsgUnfollowTraderResponse) XXX_MessageName() string { return "syreen.dex.MsgUnfollowTraderResponse" }

type MsgUpdateCopySettingsResponse struct{}

func (m *MsgUpdateCopySettingsResponse) ProtoMessage()           {}
func (m *MsgUpdateCopySettingsResponse) Reset()                  { *m = MsgUpdateCopySettingsResponse{} }
func (m *MsgUpdateCopySettingsResponse) String() string          { return "update_copy_settings_response" }
func (m *MsgUpdateCopySettingsResponse) XXX_MessageName() string { return "syreen.dex.MsgUpdateCopySettingsResponse" }

// ---------------------------------------------------------------------------
// MsgClaimReferralRewardsResponse (H5)
// ---------------------------------------------------------------------------

type MsgClaimReferralRewardsResponse struct {
	Paid sdk.Coins `protobuf:"bytes,1,rep,name=paid,proto3" json:"paid"`
}

func (m *MsgClaimReferralRewardsResponse) ProtoMessage()           {}
func (m *MsgClaimReferralRewardsResponse) Reset()                  { *m = MsgClaimReferralRewardsResponse{} }
func (m *MsgClaimReferralRewardsResponse) String() string          { return "claim_referral_rewards_response" }
func (m *MsgClaimReferralRewardsResponse) XXX_MessageName() string { return "syreen.dex.MsgClaimReferralRewardsResponse" }

// ---------------------------------------------------------------------------
// Multi-Hop Swap response type
// ---------------------------------------------------------------------------

type MsgMultiHopSwapResponse struct {
	TokenOut sdk.Coin `protobuf:"bytes,1,opt,name=token_out,json=tokenOut,proto3" json:"token_out"`
}

func (m *MsgMultiHopSwapResponse) ProtoMessage()          {}
func (m *MsgMultiHopSwapResponse) Reset()                 { *m = MsgMultiHopSwapResponse{} }
func (m *MsgMultiHopSwapResponse) String() string         { return "multi_hop_swap_response" }
func (m *MsgMultiHopSwapResponse) XXX_MessageName() string { return "syreen.dex.MsgMultiHopSwapResponse" }

// ---------------------------------------------------------------------------
// Set Pool Fee Config response type
// ---------------------------------------------------------------------------

type MsgSetPoolFeeConfigResponse struct{}

func (m *MsgSetPoolFeeConfigResponse) ProtoMessage()           {}
func (m *MsgSetPoolFeeConfigResponse) Reset()                  { *m = MsgSetPoolFeeConfigResponse{} }
func (m *MsgSetPoolFeeConfigResponse) String() string          { return "set_pool_fee_config_response" }
func (m *MsgSetPoolFeeConfigResponse) XXX_MessageName() string { return "syreen.dex.MsgSetPoolFeeConfigResponse" }
