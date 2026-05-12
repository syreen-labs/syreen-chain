package keeper

import (
	"context"
	"fmt"
	"strconv"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/dex/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) CreatePool(ctx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	poolID, err := m.Keeper.CreatePool(ctx, msg.Sender, msg.DenomA, msg.DenomB, msg.AmountA, msg.AmountB)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"create_pool",
			sdk.NewAttribute("pool_id", strconv.FormatUint(poolID, 10)),
			sdk.NewAttribute("creator", msg.Sender),
			sdk.NewAttribute("denom_a", msg.DenomA),
			sdk.NewAttribute("denom_b", msg.DenomB),
			sdk.NewAttribute("amount_a", msg.AmountA.String()),
			sdk.NewAttribute("amount_b", msg.AmountB.String()),
		),
	})

	return &types.MsgCreatePoolResponse{PoolID: poolID}, nil
}

func (m msgServer) AddLiquidity(ctx context.Context, msg *types.MsgAddLiquidity) (*types.MsgAddLiquidityResponse, error) {
	shares, err := m.Keeper.AddLiquidity(ctx, msg.Sender, msg.PoolID, msg.AmountA, msg.AmountB, msg.MinSharesOut)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"add_liquidity",
			sdk.NewAttribute("pool_id", strconv.FormatUint(msg.PoolID, 10)),
			sdk.NewAttribute("sender", msg.Sender),
			sdk.NewAttribute("shares_minted", shares.String()),
		),
	})

	return &types.MsgAddLiquidityResponse{SharesMinted: shares}, nil
}

func (m msgServer) RemoveLiquidity(ctx context.Context, msg *types.MsgRemoveLiquidity) (*types.MsgRemoveLiquidityResponse, error) {
	amountA, amountB, err := m.Keeper.RemoveLiquidity(ctx, msg.Sender, msg.PoolID, msg.SharesIn, msg.MinAmountAOut, msg.MinAmountBOut)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"remove_liquidity",
			sdk.NewAttribute("pool_id", strconv.FormatUint(msg.PoolID, 10)),
			sdk.NewAttribute("sender", msg.Sender),
			sdk.NewAttribute("amount_a", amountA.String()),
			sdk.NewAttribute("amount_b", amountB.String()),
			sdk.NewAttribute("shares_burned", msg.SharesIn.String()),
		),
	})

	return &types.MsgRemoveLiquidityResponse{AmountA: amountA, AmountB: amountB}, nil
}

func (m msgServer) Swap(ctx context.Context, msg *types.MsgSwap) (*types.MsgSwapResponse, error) {
	tokenOut, err := m.Keeper.Swap(ctx, msg.Sender, msg.PoolID, msg.TokenIn, msg.MinTokenOut)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"swap",
			sdk.NewAttribute("pool_id", strconv.FormatUint(msg.PoolID, 10)),
			sdk.NewAttribute("sender", msg.Sender),
			sdk.NewAttribute("token_in", msg.TokenIn.String()),
			sdk.NewAttribute("token_out", tokenOut.String()),
		),
	})

	return &types.MsgSwapResponse{TokenOut: tokenOut}, nil
}

func (m msgServer) PlaceOrder(ctx context.Context, msg *types.MsgPlaceOrder) (*types.MsgPlaceOrderResponse, error) {
	orderID, status, err := m.Keeper.PlaceOrder(ctx, msg)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"place_order",
			sdk.NewAttribute("order_id", strconv.FormatUint(orderID, 10)),
			sdk.NewAttribute("pool_id", strconv.FormatUint(msg.PoolID, 10)),
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("side", msg.Side),
			sdk.NewAttribute("order_type", msg.OrderType),
			sdk.NewAttribute("price", msg.Price),
			sdk.NewAttribute("quantity", msg.Quantity),
			sdk.NewAttribute("status", status),
		),
	})

	return &types.MsgPlaceOrderResponse{OrderID: orderID, Status: status}, nil
}

func (m msgServer) CancelOrder(ctx context.Context, msg *types.MsgCancelOrder) (*types.MsgCancelOrderResponse, error) {
	err := m.Keeper.CancelOrder(ctx, msg.Creator, msg.OrderID)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"cancel_order",
			sdk.NewAttribute("order_id", strconv.FormatUint(msg.OrderID, 10)),
			sdk.NewAttribute("creator", msg.Creator),
		),
	})

	return &types.MsgCancelOrderResponse{}, nil
}

func (m msgServer) ModifyOrder(ctx context.Context, msg *types.MsgModifyOrder) (*types.MsgModifyOrderResponse, error) {
	newOrderID, err := m.Keeper.ModifyOrder(ctx, msg.Creator, msg.OrderID, msg.NewPrice, msg.NewQuantity)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"modify_order",
			sdk.NewAttribute("old_order_id", strconv.FormatUint(msg.OrderID, 10)),
			sdk.NewAttribute("new_order_id", strconv.FormatUint(newOrderID, 10)),
			sdk.NewAttribute("creator", msg.Creator),
		),
	})

	return &types.MsgModifyOrderResponse{OrderID: newOrderID}, nil
}

func (m msgServer) CreateReferralCode(ctx context.Context, msg *types.MsgCreateReferralCode) (*types.MsgCreateReferralCodeResponse, error) {
	code, err := m.Keeper.CreateReferralCode(ctx, msg.Creator, msg.CustomCode)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"create_referral_code",
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("referral_code", code),
		),
	})

	return &types.MsgCreateReferralCodeResponse{ReferralCode: code}, nil
}

func (m msgServer) RegisterReferral(ctx context.Context, msg *types.MsgRegisterReferral) (*types.MsgRegisterReferralResponse, error) {
	referrer, err := m.Keeper.RegisterReferral(ctx, msg.User, msg.ReferralCode)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"register_referral",
			sdk.NewAttribute("user", msg.User),
			sdk.NewAttribute("referrer", referrer),
			sdk.NewAttribute("referral_code", msg.ReferralCode),
		),
	})

	return &types.MsgRegisterReferralResponse{Referrer: referrer}, nil
}

func (m msgServer) FollowTrader(ctx context.Context, msg *types.MsgFollowTrader) (*types.MsgFollowTraderResponse, error) {
	maxPerTrade, ok := math.NewIntFromString(msg.MaxPerTrade)
	if !ok {
		return nil, fmt.Errorf("invalid max_per_trade: %s", msg.MaxPerTrade)
	}
	totalBudget, ok := math.NewIntFromString(msg.TotalBudget)
	if !ok {
		return nil, fmt.Errorf("invalid total_budget: %s", msg.TotalBudget)
	}
	err := m.Keeper.FollowTrader(ctx, msg.Follower, msg.Trader, maxPerTrade, totalBudget, msg.CopyRatio)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"follow_trader",
			sdk.NewAttribute("follower", msg.Follower),
			sdk.NewAttribute("trader", msg.Trader),
			sdk.NewAttribute("copy_ratio", strconv.FormatInt(msg.CopyRatio, 10)),
			sdk.NewAttribute("max_per_trade", msg.MaxPerTrade),
			sdk.NewAttribute("total_budget", msg.TotalBudget),
		),
	})

	return &types.MsgFollowTraderResponse{}, nil
}

func (m msgServer) UnfollowTrader(ctx context.Context, msg *types.MsgUnfollowTrader) (*types.MsgUnfollowTraderResponse, error) {
	err := m.Keeper.UnfollowTrader(ctx, msg.Follower, msg.Trader)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"unfollow_trader",
			sdk.NewAttribute("follower", msg.Follower),
			sdk.NewAttribute("trader", msg.Trader),
		),
	})

	return &types.MsgUnfollowTraderResponse{}, nil
}

func (m msgServer) UpdateCopySettings(ctx context.Context, msg *types.MsgUpdateCopySettings) (*types.MsgUpdateCopySettingsResponse, error) {
	ucMaxPerTrade, ok := math.NewIntFromString(msg.MaxPerTrade)
	if !ok {
		return nil, fmt.Errorf("invalid max_per_trade: %s", msg.MaxPerTrade)
	}
	ucTotalBudget, ok := math.NewIntFromString(msg.TotalBudget)
	if !ok {
		return nil, fmt.Errorf("invalid total_budget: %s", msg.TotalBudget)
	}
	err := m.Keeper.UpdateCopySettings(ctx, msg.Follower, msg.Trader, ucMaxPerTrade, ucTotalBudget, msg.CopyRatio)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"update_copy_settings",
			sdk.NewAttribute("follower", msg.Follower),
			sdk.NewAttribute("trader", msg.Trader),
		),
	})

	return &types.MsgUpdateCopySettingsResponse{}, nil
}

// ClaimReferralRewards (H5) lets a referrer withdraw their accumulated
// claimable referral fees from the dex module account.
func (m msgServer) ClaimReferralRewards(ctx context.Context, msg *types.MsgClaimReferralRewards) (*types.MsgClaimReferralRewardsResponse, error) {
	paid, err := m.Keeper.ClaimReferralRewards(ctx, msg.Referrer)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"claim_referral_rewards",
			sdk.NewAttribute("referrer", msg.Referrer),
			sdk.NewAttribute("paid", paid.String()),
		),
	})

	return &types.MsgClaimReferralRewardsResponse{Paid: paid}, nil
}

func (m msgServer) MultiHopSwap(ctx context.Context, msg *types.MsgMultiHopSwap) (*types.MsgMultiHopSwapResponse, error) {
	tokenInAmount, ok := math.NewIntFromString(msg.TokenInAmount)
	if !ok {
		return nil, fmt.Errorf("invalid token_in_amount: %s", msg.TokenInAmount)
	}
	tokenIn := sdk.NewCoin(msg.TokenInDenom, tokenInAmount)
	minTokenOut, ok := math.NewIntFromString(msg.MinTokenOut)
	if !ok {
		return nil, fmt.Errorf("invalid min_token_out: %s", msg.MinTokenOut)
	}

	tokenOut, err := m.Keeper.MultiHopSwap(ctx, msg.Sender, msg.Route, tokenIn, minTokenOut)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"multi_hop_swap",
			sdk.NewAttribute("sender", msg.Sender),
			sdk.NewAttribute("route", fmt.Sprintf("%v", msg.Route)),
			sdk.NewAttribute("token_in", tokenIn.String()),
			sdk.NewAttribute("token_out", tokenOut.String()),
			sdk.NewAttribute("hops", strconv.Itoa(len(msg.Route))),
		),
	})

	return &types.MsgMultiHopSwapResponse{TokenOut: tokenOut}, nil
}

func (m msgServer) SetPoolFeeConfig(ctx context.Context, msg *types.MsgSetPoolFeeConfig) (*types.MsgSetPoolFeeConfigResponse, error) {
	baseFee, err := math.LegacyNewDecFromStr(msg.BaseFee)
	if err != nil {
		return nil, fmt.Errorf("invalid base_fee: %s", msg.BaseFee)
	}
	maxFee, err := math.LegacyNewDecFromStr(msg.MaxFee)
	if err != nil {
		return nil, fmt.Errorf("invalid max_fee: %s", msg.MaxFee)
	}
	multiplier, err := math.LegacyNewDecFromStr(msg.VolatilityMultiplier)
	if err != nil {
		return nil, fmt.Errorf("invalid volatility_multiplier: %s", msg.VolatilityMultiplier)
	}

	err = m.Keeper.SetPoolFeeConfig(ctx, msg.Sender, msg.PoolID, msg.VolatilityFeeEnabled, baseFee, maxFee, multiplier)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"set_pool_fee_config",
			sdk.NewAttribute("pool_id", strconv.FormatUint(msg.PoolID, 10)),
			sdk.NewAttribute("sender", msg.Sender),
			sdk.NewAttribute("volatility_fee_enabled", strconv.FormatBool(msg.VolatilityFeeEnabled)),
			sdk.NewAttribute("base_fee", msg.BaseFee),
			sdk.NewAttribute("max_fee", msg.MaxFee),
			sdk.NewAttribute("volatility_multiplier", msg.VolatilityMultiplier),
		),
	})

	return &types.MsgSetPoolFeeConfigResponse{}, nil
}

// suppress unused import
var _ = fmt.Sprintf
