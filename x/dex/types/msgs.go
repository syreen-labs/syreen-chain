package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgCreatePool         = "create_pool"
	TypeMsgAddLiquidity       = "add_liquidity"
	TypeMsgRemoveLiquidity    = "remove_liquidity"
	TypeMsgSwap               = "swap"
	TypeMsgCreateReferralCode  = "create_referral_code"
	TypeMsgRegisterReferral    = "register_referral"
	TypeMsgClaimReferralRewards = "claim_referral_rewards"
	TypeMsgPlaceOrder         = "place_order"
	TypeMsgCancelOrder        = "cancel_order"
	TypeMsgModifyOrder        = "modify_order"
	TypeMsgMultiHopSwap       = "multi_hop_swap"
)

// ---------------------------------------------------------------------------
// MsgCreatePool
// ---------------------------------------------------------------------------

type MsgCreatePool struct {
	Sender  string   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	DenomA  string   `protobuf:"bytes,2,opt,name=denom_a,json=denomA,proto3" json:"denom_a"`
	DenomB  string   `protobuf:"bytes,3,opt,name=denom_b,json=denomB,proto3" json:"denom_b"`
	AmountA math.Int `protobuf:"bytes,4,opt,name=amount_a,json=amountA,proto3" json:"amount_a"`
	AmountB math.Int `protobuf:"bytes,5,opt,name=amount_b,json=amountB,proto3" json:"amount_b"`
}

func (m *MsgCreatePool) ProtoMessage()          {}
func (m *MsgCreatePool) Reset()                 { *m = MsgCreatePool{} }
func (m *MsgCreatePool) String() string         { return fmt.Sprintf("create_pool: %s/%s by %s", m.DenomA, m.DenomB, m.Sender) }
func (m *MsgCreatePool) XXX_MessageName() string { return "syreen.dex.MsgCreatePool" }

func (m *MsgCreatePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return ErrInvalidSender
	}
	if m.DenomA == "" || m.DenomB == "" {
		return ErrInvalidDenom
	}
	if err := sdk.ValidateDenom(m.DenomA); err != nil {
		return ErrInvalidDenom
	}
	if err := sdk.ValidateDenom(m.DenomB); err != nil {
		return ErrInvalidDenom
	}
	if m.DenomA == m.DenomB {
		return ErrSameDenom
	}
	if !m.AmountA.IsPositive() || !m.AmountB.IsPositive() {
		return ErrInvalidAmount
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgAddLiquidity
// ---------------------------------------------------------------------------

type MsgAddLiquidity struct {
	Sender      string   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	PoolID      uint64   `protobuf:"varint,2,opt,name=pool_id,json=poolId,proto3" json:"pool_id"`
	AmountA     math.Int `protobuf:"bytes,3,opt,name=amount_a,json=amountA,proto3" json:"amount_a"`
	AmountB     math.Int `protobuf:"bytes,4,opt,name=amount_b,json=amountB,proto3" json:"amount_b"`
	MinSharesOut math.Int `protobuf:"bytes,5,opt,name=min_shares_out,json=minSharesOut,proto3" json:"min_shares_out"`
}

func (m *MsgAddLiquidity) ProtoMessage()          {}
func (m *MsgAddLiquidity) Reset()                 { *m = MsgAddLiquidity{} }
func (m *MsgAddLiquidity) String() string         { return fmt.Sprintf("add_liquidity: pool %d by %s", m.PoolID, m.Sender) }
func (m *MsgAddLiquidity) XXX_MessageName() string { return "syreen.dex.MsgAddLiquidity" }

func (m *MsgAddLiquidity) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return ErrInvalidSender
	}
	if m.PoolID == 0 {
		return ErrPoolNotFound
	}
	if m.AmountA.IsNil() || !m.AmountA.IsPositive() || m.AmountB.IsNil() || !m.AmountB.IsPositive() {
		return ErrInvalidAmount
	}
	if !m.MinSharesOut.IsNil() && m.MinSharesOut.IsNegative() {
		return ErrInvalidAmount
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgRemoveLiquidity
// ---------------------------------------------------------------------------

type MsgRemoveLiquidity struct {
	Sender       string   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	PoolID       uint64   `protobuf:"varint,2,opt,name=pool_id,json=poolId,proto3" json:"pool_id"`
	SharesIn     math.Int `protobuf:"bytes,3,opt,name=shares_in,json=sharesIn,proto3" json:"shares_in"`
	MinAmountAOut math.Int `protobuf:"bytes,4,opt,name=min_amount_a_out,json=minAmountAOut,proto3" json:"min_amount_a_out"`
	MinAmountBOut math.Int `protobuf:"bytes,5,opt,name=min_amount_b_out,json=minAmountBOut,proto3" json:"min_amount_b_out"`
}

func (m *MsgRemoveLiquidity) ProtoMessage()          {}
func (m *MsgRemoveLiquidity) Reset()                 { *m = MsgRemoveLiquidity{} }
func (m *MsgRemoveLiquidity) String() string         { return fmt.Sprintf("remove_liquidity: pool %d by %s", m.PoolID, m.Sender) }
func (m *MsgRemoveLiquidity) XXX_MessageName() string { return "syreen.dex.MsgRemoveLiquidity" }

func (m *MsgRemoveLiquidity) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return ErrInvalidSender
	}
	if m.PoolID == 0 {
		return ErrPoolNotFound
	}
	if m.SharesIn.IsNil() || !m.SharesIn.IsPositive() {
		return ErrInvalidAmount
	}
	if (!m.MinAmountAOut.IsNil() && m.MinAmountAOut.IsNegative()) || (!m.MinAmountBOut.IsNil() && m.MinAmountBOut.IsNegative()) {
		return ErrInvalidAmount
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgSwap
// ---------------------------------------------------------------------------

type MsgSwap struct {
	Sender      string   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	PoolID      uint64   `protobuf:"varint,2,opt,name=pool_id,json=poolId,proto3" json:"pool_id"`
	TokenIn     sdk.Coin `protobuf:"bytes,3,opt,name=token_in,json=tokenIn,proto3" json:"token_in"`
	MinTokenOut math.Int `protobuf:"bytes,4,opt,name=min_token_out,json=minTokenOut,proto3" json:"min_token_out"`
}

func (m *MsgSwap) ProtoMessage()          {}
func (m *MsgSwap) Reset()                 { *m = MsgSwap{} }
func (m *MsgSwap) String() string         { return fmt.Sprintf("swap: pool %d %s by %s", m.PoolID, m.TokenIn, m.Sender) }
func (m *MsgSwap) XXX_MessageName() string { return "syreen.dex.MsgSwap" }

func (m *MsgSwap) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return ErrInvalidSender
	}
	if m.PoolID == 0 {
		return ErrPoolNotFound
	}
	if !m.TokenIn.IsValid() || m.TokenIn.IsZero() {
		return ErrInvalidAmount
	}
	if !m.MinTokenOut.IsNil() && m.MinTokenOut.IsNegative() {
		return ErrInvalidAmount
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgCreateReferralCode
// ---------------------------------------------------------------------------

type MsgCreateReferralCode struct {
	Creator    string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator"`
	CustomCode string `protobuf:"bytes,2,opt,name=custom_code,proto3" json:"custom_code,omitempty"`
}

func (m *MsgCreateReferralCode) ProtoMessage()          {}
func (m *MsgCreateReferralCode) Reset()                 { *m = MsgCreateReferralCode{} }
func (m *MsgCreateReferralCode) String() string         { return fmt.Sprintf("create_referral_code: by %s", m.Creator) }
func (m *MsgCreateReferralCode) XXX_MessageName() string { return "syreen.dex.MsgCreateReferralCode" }

func (m *MsgCreateReferralCode) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		return ErrInvalidSender
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgRegisterReferral
// ---------------------------------------------------------------------------

type MsgRegisterReferral struct {
	User         string `protobuf:"bytes,1,opt,name=user,proto3" json:"user"`
	ReferralCode string `protobuf:"bytes,2,opt,name=referral_code,json=referralCode,proto3" json:"referral_code"`
}

func (m *MsgRegisterReferral) ProtoMessage()          {}
func (m *MsgRegisterReferral) Reset()                 { *m = MsgRegisterReferral{} }
func (m *MsgRegisterReferral) String() string         { return fmt.Sprintf("register_referral: %s with code %s", m.User, m.ReferralCode) }
func (m *MsgRegisterReferral) XXX_MessageName() string { return "syreen.dex.MsgRegisterReferral" }

func (m *MsgRegisterReferral) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.User)
	if err != nil {
		return ErrInvalidSender
	}
	if m.ReferralCode == "" {
		return ErrReferralCodeNotFound
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgClaimReferralRewards (H5)
// ---------------------------------------------------------------------------

// MsgClaimReferralRewards lets a referrer withdraw all of their accumulated
// referral fee rewards (across every denom) from the dex module account into
// their own account.
type MsgClaimReferralRewards struct {
	Referrer string `protobuf:"bytes,1,opt,name=referrer,proto3" json:"referrer"`
}

func (m *MsgClaimReferralRewards) ProtoMessage()           {}
func (m *MsgClaimReferralRewards) Reset()                  { *m = MsgClaimReferralRewards{} }
func (m *MsgClaimReferralRewards) String() string          { return fmt.Sprintf("claim_referral_rewards: by %s", m.Referrer) }
func (m *MsgClaimReferralRewards) XXX_MessageName() string { return "syreen.dex.MsgClaimReferralRewards" }

func (m *MsgClaimReferralRewards) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Referrer)
	if err != nil {
		return ErrInvalidSender
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgPlaceOrder
// ---------------------------------------------------------------------------

type MsgPlaceOrder struct {
	Creator      string         `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator"`
	PoolID       uint64         `protobuf:"varint,2,opt,name=pool_id,json=poolId,proto3" json:"pool_id"`
	Side         string         `protobuf:"bytes,3,opt,name=side,proto3" json:"side"`
	OrderType    string         `protobuf:"bytes,4,opt,name=order_type,json=orderType,proto3" json:"order_type"`
	Price        string   `protobuf:"bytes,5,opt,name=price,proto3" json:"price"`
	Quantity     string   `protobuf:"bytes,6,opt,name=quantity,proto3" json:"quantity"`
	TimeInForce  string   `protobuf:"bytes,7,opt,name=time_in_force,json=timeInForce,proto3" json:"time_in_force"`
	TriggerPrice string   `protobuf:"bytes,8,opt,name=trigger_price,json=triggerPrice,proto3" json:"trigger_price"`
}

func (m *MsgPlaceOrder) ProtoMessage()          {}
func (m *MsgPlaceOrder) Reset()                 { *m = MsgPlaceOrder{} }
func (m *MsgPlaceOrder) String() string         { return fmt.Sprintf("place_order: pool %d %s %s by %s", m.PoolID, m.Side, m.OrderType, m.Creator) }
func (m *MsgPlaceOrder) XXX_MessageName() string { return "syreen.dex.MsgPlaceOrder" }

func (m *MsgPlaceOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		return ErrInvalidSender
	}
	if m.PoolID == 0 {
		return ErrPoolNotFound
	}
	if m.Side != OrderSideBuy && m.Side != OrderSideSell {
		return ErrInvalidOrderSide
	}

	isConditional := IsConditionalOrderType(m.OrderType)

	if m.OrderType != OrderTypeLimit && m.OrderType != OrderTypeMarket && !isConditional {
		return ErrInvalidOrderType
	}

	// Parse and validate Price
	var priceDec math.LegacyDec
	if m.Price != "" && m.Price != "0" {
		var err2 error
		priceDec, err2 = math.LegacyNewDecFromStr(m.Price)
		if err2 != nil {
			return ErrInvalidOrderPrice
		}
		if !priceDec.IsPositive() {
			return ErrInvalidOrderPrice
		}
	}

	// Price validation: required for limit orders and conditional-limit orders
	if m.OrderType == OrderTypeLimit && (m.Price == "" || m.Price == "0") {
		return ErrInvalidOrderPrice
	}
	if (m.OrderType == OrderTypeStopLossLimit || m.OrderType == OrderTypeTakeProfitLimit) && (m.Price == "" || m.Price == "0") {
		return ErrInvalidOrderPrice
	}

	// Parse and validate TriggerPrice
	var triggerPriceDec math.LegacyDec
	if m.TriggerPrice != "" && m.TriggerPrice != "0" {
		var err3 error
		triggerPriceDec, err3 = math.LegacyNewDecFromStr(m.TriggerPrice)
		if err3 != nil {
			return ErrInvalidTriggerPrice
		}
		if !triggerPriceDec.IsPositive() {
			return ErrInvalidTriggerPrice
		}
	}

	// TriggerPrice validation
	if isConditional {
		if m.TriggerPrice == "" || m.TriggerPrice == "0" {
			return ErrInvalidTriggerPrice
		}
	} else {
		// Regular limit/market orders must NOT have a trigger price
		if m.TriggerPrice != "" && m.TriggerPrice != "0" && !triggerPriceDec.IsZero() {
			return ErrTriggerPriceNotAllowed
		}
	}
	_ = priceDec

	if m.Quantity == "" || m.Quantity == "0" {
		return ErrInvalidAmount
	}
	if qty, ok := math.NewIntFromString(m.Quantity); !ok || !qty.IsPositive() {
		return ErrInvalidAmount
	}
	if m.TimeInForce != TimeInForceGTC && m.TimeInForce != TimeInForceIOC && m.TimeInForce != TimeInForceFOK {
		return ErrInvalidTimeInForce
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgCancelOrder
// ---------------------------------------------------------------------------

type MsgCancelOrder struct {
	Creator string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator"`
	OrderID uint64 `protobuf:"varint,2,opt,name=order_id,json=orderId,proto3" json:"order_id"`
}

func (m *MsgCancelOrder) ProtoMessage()          {}
func (m *MsgCancelOrder) Reset()                 { *m = MsgCancelOrder{} }
func (m *MsgCancelOrder) String() string         { return fmt.Sprintf("cancel_order: %d by %s", m.OrderID, m.Creator) }
func (m *MsgCancelOrder) XXX_MessageName() string { return "syreen.dex.MsgCancelOrder" }

func (m *MsgCancelOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		return ErrInvalidSender
	}
	if m.OrderID == 0 {
		return ErrOrderNotFound
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgModifyOrder
// ---------------------------------------------------------------------------

type MsgModifyOrder struct {
	Creator     string         `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator"`
	OrderID     uint64         `protobuf:"varint,2,opt,name=order_id,json=orderId,proto3" json:"order_id"`
	NewPrice    string         `protobuf:"bytes,3,opt,name=new_price,json=newPrice,proto3" json:"new_price"`
	NewQuantity math.Int       `protobuf:"bytes,4,opt,name=new_quantity,json=newQuantity,proto3,customtype=cosmossdk.io/math.Int" json:"new_quantity"`
}

func (m *MsgModifyOrder) ProtoMessage()          {}
func (m *MsgModifyOrder) Reset()                 { *m = MsgModifyOrder{} }
func (m *MsgModifyOrder) String() string         { return fmt.Sprintf("modify_order: %d by %s", m.OrderID, m.Creator) }
func (m *MsgModifyOrder) XXX_MessageName() string { return "syreen.dex.MsgModifyOrder" }

func (m *MsgModifyOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		return ErrInvalidSender
	}
	if m.OrderID == 0 {
		return ErrOrderNotFound
	}
	if m.NewPrice != "" && m.NewPrice != "0" {
		p, err2 := math.LegacyNewDecFromStr(m.NewPrice)
		if err2 != nil {
			return ErrInvalidOrderPrice
		}
		if !p.IsPositive() {
			return ErrInvalidOrderPrice
		}
	}
	if !m.NewQuantity.IsNil() && !m.NewQuantity.IsPositive() {
		return ErrInvalidAmount
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgSetPoolFeeConfig
// ---------------------------------------------------------------------------

const TypeMsgSetPoolFeeConfig = "set_pool_fee_config"

type MsgSetPoolFeeConfig struct {
	Sender               string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	PoolID               uint64 `protobuf:"varint,2,opt,name=pool_id,json=poolId,proto3" json:"pool_id"`
	VolatilityFeeEnabled bool   `protobuf:"varint,3,opt,name=volatility_fee_enabled,json=volatilityFeeEnabled,proto3" json:"volatility_fee_enabled"`
	BaseFee              string `protobuf:"bytes,4,opt,name=base_fee,json=baseFee,proto3" json:"base_fee"`
	MaxFee               string `protobuf:"bytes,5,opt,name=max_fee,json=maxFee,proto3" json:"max_fee"`
	VolatilityMultiplier string `protobuf:"bytes,6,opt,name=volatility_multiplier,json=volatilityMultiplier,proto3" json:"volatility_multiplier"`
}

func (m *MsgSetPoolFeeConfig) ProtoMessage()           {}
func (m *MsgSetPoolFeeConfig) Reset()                  { *m = MsgSetPoolFeeConfig{} }
func (m *MsgSetPoolFeeConfig) String() string          { return fmt.Sprintf("set_pool_fee_config: pool %d by %s", m.PoolID, m.Sender) }
func (m *MsgSetPoolFeeConfig) XXX_MessageName() string { return "syreen.dex.MsgSetPoolFeeConfig" }

func (m *MsgSetPoolFeeConfig) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return ErrInvalidSender
	}
	if m.PoolID == 0 {
		return ErrPoolNotFound
	}
	if m.BaseFee != "" {
		bf, err := math.LegacyNewDecFromStr(m.BaseFee)
		if err != nil {
			return fmt.Errorf("invalid base_fee: %w", err)
		}
		if bf.IsNegative() {
			return ErrInvalidAmount
		}
	}
	if m.MaxFee != "" {
		mf, err := math.LegacyNewDecFromStr(m.MaxFee)
		if err != nil {
			return fmt.Errorf("invalid max_fee: %w", err)
		}
		if mf.IsNegative() {
			return ErrInvalidAmount
		}
	}
	if m.VolatilityMultiplier != "" {
		vm, err := math.LegacyNewDecFromStr(m.VolatilityMultiplier)
		if err != nil {
			return fmt.Errorf("invalid volatility_multiplier: %w", err)
		}
		if vm.IsNegative() {
			return ErrInvalidAmount
		}
	}
	if m.BaseFee != "" && m.MaxFee != "" {
		bf, _ := math.LegacyNewDecFromStr(m.BaseFee)
		mf, _ := math.LegacyNewDecFromStr(m.MaxFee)
		if mf.LT(bf) {
			return ErrInvalidFeeConfig
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgFollowTrader
// ---------------------------------------------------------------------------

type MsgFollowTrader struct {
	Follower    string `protobuf:"bytes,1,opt,name=follower,proto3" json:"follower"`
	Trader      string `protobuf:"bytes,2,opt,name=trader,proto3" json:"trader"`
	MaxPerTrade string `protobuf:"bytes,3,opt,name=max_per_trade,json=maxPerTrade,proto3" json:"max_per_trade"`
	TotalBudget string `protobuf:"bytes,4,opt,name=total_budget,json=totalBudget,proto3" json:"total_budget"`
	CopyRatio   int64  `protobuf:"varint,5,opt,name=copy_ratio,json=copyRatio,proto3" json:"copy_ratio"`
}

func (m *MsgFollowTrader) ProtoMessage()           {}
func (m *MsgFollowTrader) Reset()                  { *m = MsgFollowTrader{} }
func (m *MsgFollowTrader) String() string          { return fmt.Sprintf("follow_trader: %s follows %s", m.Follower, m.Trader) }
func (m *MsgFollowTrader) XXX_MessageName() string { return "syreen.dex.MsgFollowTrader" }

func (m *MsgFollowTrader) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Follower)
	if err != nil {
		return ErrInvalidSender
	}
	_, err = sdk.AccAddressFromBech32(m.Trader)
	if err != nil {
		return fmt.Errorf("invalid trader address")
	}
	if m.Follower == m.Trader {
		return fmt.Errorf("cannot follow yourself")
	}
	maxPT, _ := math.NewIntFromString(m.MaxPerTrade); if maxPT.IsNil() || !maxPT.IsPositive() {
		return ErrInvalidAmount
	}
	totalB, _ := math.NewIntFromString(m.TotalBudget); if totalB.IsNil() || !totalB.IsPositive() {
		return ErrInvalidAmount
	}
	if m.CopyRatio <= 0 || m.CopyRatio > 10000 {
		return fmt.Errorf("copy_ratio must be between 1 and 10000 basis points")
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgUnfollowTrader
// ---------------------------------------------------------------------------

type MsgUnfollowTrader struct {
	Follower string `protobuf:"bytes,1,opt,name=follower,proto3" json:"follower"`
	Trader   string `protobuf:"bytes,2,opt,name=trader,proto3" json:"trader"`
}

func (m *MsgUnfollowTrader) ProtoMessage()           {}
func (m *MsgUnfollowTrader) Reset()                  { *m = MsgUnfollowTrader{} }
func (m *MsgUnfollowTrader) String() string          { return fmt.Sprintf("unfollow_trader: %s unfollows %s", m.Follower, m.Trader) }
func (m *MsgUnfollowTrader) XXX_MessageName() string { return "syreen.dex.MsgUnfollowTrader" }

func (m *MsgUnfollowTrader) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Follower)
	if err != nil {
		return ErrInvalidSender
	}
	_, err = sdk.AccAddressFromBech32(m.Trader)
	if err != nil {
		return fmt.Errorf("invalid trader address")
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgMultiHopSwap
// ---------------------------------------------------------------------------

type MsgMultiHopSwap struct {
	Sender         string   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	Route          []uint64 `protobuf:"varint,2,rep,packed,name=route,proto3" json:"route"`
	TokenInDenom   string   `protobuf:"bytes,3,opt,name=token_in_denom,json=tokenInDenom,proto3" json:"token_in_denom"`
	TokenInAmount  string   `protobuf:"bytes,4,opt,name=token_in_amount,json=tokenInAmount,proto3" json:"token_in_amount"`
	MinTokenOut    string   `protobuf:"bytes,5,opt,name=min_token_out,json=minTokenOut,proto3" json:"min_token_out"`
}

func (m *MsgMultiHopSwap) ProtoMessage()          {}
func (m *MsgMultiHopSwap) Reset()                 { *m = MsgMultiHopSwap{} }
func (m *MsgMultiHopSwap) String() string         { return fmt.Sprintf("multi_hop_swap: %d hops %s by %s", len(m.Route), m.TokenInDenom, m.Sender) }
func (m *MsgMultiHopSwap) XXX_MessageName() string { return "syreen.dex.MsgMultiHopSwap" }

func (m *MsgMultiHopSwap) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return ErrInvalidSender
	}
	if len(m.Route) == 0 {
		return fmt.Errorf("route must contain at least one pool")
	}
	if len(m.Route) > 4 {
		return ErrTooManyHops
	}
	for _, poolID := range m.Route {
		if poolID == 0 {
			return ErrPoolNotFound
		}
	}
	if m.TokenInDenom == "" {
		return ErrInvalidDenom
	}
	amt, ok := math.NewIntFromString(m.TokenInAmount)
	if !ok || !amt.IsPositive() {
		return ErrInvalidAmount
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgUpdateCopySettings
// ---------------------------------------------------------------------------

type MsgUpdateCopySettings struct {
	Follower    string `protobuf:"bytes,1,opt,name=follower,proto3" json:"follower"`
	Trader      string `protobuf:"bytes,2,opt,name=trader,proto3" json:"trader"`
	MaxPerTrade string `protobuf:"bytes,3,opt,name=max_per_trade,json=maxPerTrade,proto3" json:"max_per_trade"`
	TotalBudget string `protobuf:"bytes,4,opt,name=total_budget,json=totalBudget,proto3" json:"total_budget"`
	CopyRatio   int64  `protobuf:"varint,5,opt,name=copy_ratio,json=copyRatio,proto3" json:"copy_ratio"`
}

func (m *MsgUpdateCopySettings) ProtoMessage()           {}
func (m *MsgUpdateCopySettings) Reset()                  { *m = MsgUpdateCopySettings{} }
func (m *MsgUpdateCopySettings) String() string          { return fmt.Sprintf("update_copy_settings: %s -> %s", m.Follower, m.Trader) }
func (m *MsgUpdateCopySettings) XXX_MessageName() string { return "syreen.dex.MsgUpdateCopySettings" }

func (m *MsgUpdateCopySettings) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Follower)
	if err != nil {
		return ErrInvalidSender
	}
	_, err = sdk.AccAddressFromBech32(m.Trader)
	if err != nil {
		return fmt.Errorf("invalid trader address")
	}
	return nil
}
