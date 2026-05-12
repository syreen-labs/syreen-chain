package keeper

import (
	"context"

	syreenconfig "syreen/config"
	"syreen/x/clmm/types"
)

var _ types.MsgServer = Keeper{}

func (k Keeper) CreateCLPool(ctx context.Context, msg *types.MsgCreateCLPool) (*types.MsgCreateCLPoolResponse, error) {
	if !syreenconfig.IsModuleEnabled("clmm") {
		return nil, syreenconfig.ErrModuleDisabled("clmm")
	}
	poolID, err := k.CreatePool(ctx, msg.Sender, msg.DenomA, msg.DenomB, msg.TickSpacing, msg.FeeRate, msg.InitialPrice)
	if err != nil {
		return nil, err
	}
	return &types.MsgCreateCLPoolResponse{PoolID: poolID}, nil
}

func (k Keeper) CreatePosition(ctx context.Context, msg *types.MsgCreatePosition) (*types.MsgCreatePositionResponse, error) {
	if !syreenconfig.IsModuleEnabled("clmm") {
		return nil, syreenconfig.ErrModuleDisabled("clmm")
	}
	posID, amt0, amt1, liq, err := k.ExecuteCreatePosition(ctx, msg.Sender, msg.PoolID, msg.TickLower, msg.TickUpper, msg.Amount0Desired, msg.Amount1Desired, msg.Amount0Min, msg.Amount1Min)
	if err != nil {
		return nil, err
	}
	return &types.MsgCreatePositionResponse{PositionID: posID, Amount0: amt0, Amount1: amt1, Liquidity: liq}, nil
}

func (k Keeper) AddLiquidity(ctx context.Context, msg *types.MsgAddLiquidity) (*types.MsgAddLiquidityResponse, error) {
	if !syreenconfig.IsModuleEnabled("clmm") {
		return nil, syreenconfig.ErrModuleDisabled("clmm")
	}
	amt0, amt1, liq, err := k.AddLiquidityToPosition(ctx, msg.Sender, msg.PositionID, msg.Amount0Desired, msg.Amount1Desired)
	if err != nil {
		return nil, err
	}
	return &types.MsgAddLiquidityResponse{Amount0: amt0, Amount1: amt1, Liquidity: liq}, nil
}

func (k Keeper) RemoveLiquidity(ctx context.Context, msg *types.MsgRemoveLiquidity) (*types.MsgRemoveLiquidityResponse, error) {
	if !syreenconfig.IsModuleEnabled("clmm") {
		return nil, syreenconfig.ErrModuleDisabled("clmm")
	}
	amt0, amt1, err := k.RemoveLiquidityFromPosition(ctx, msg.Sender, msg.PositionID, msg.LiquidityAmount)
	if err != nil {
		return nil, err
	}
	return &types.MsgRemoveLiquidityResponse{Amount0: amt0, Amount1: amt1}, nil
}

func (k Keeper) CollectFees(ctx context.Context, msg *types.MsgCollectFees) (*types.MsgCollectFeesResponse, error) {
	if !syreenconfig.IsModuleEnabled("clmm") {
		return nil, syreenconfig.ErrModuleDisabled("clmm")
	}
	amt0, amt1, err := k.CollectPositionFees(ctx, msg.Sender, msg.PositionID)
	if err != nil {
		return nil, err
	}
	return &types.MsgCollectFeesResponse{Amount0: amt0, Amount1: amt1}, nil
}

func (k Keeper) CLSwap(ctx context.Context, msg *types.MsgCLSwap) (*types.MsgCLSwapResponse, error) {
	if !syreenconfig.IsModuleEnabled("clmm") {
		return nil, syreenconfig.ErrModuleDisabled("clmm")
	}
	amountOut, err := k.ExecuteSwap(ctx, msg.Sender, msg.PoolID, msg.DenomIn, msg.AmountIn, msg.MinAmountOut)
	if err != nil {
		return nil, err
	}
	return &types.MsgCLSwapResponse{AmountOut: amountOut}, nil
}
