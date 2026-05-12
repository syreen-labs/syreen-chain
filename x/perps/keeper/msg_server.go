package keeper

import (
	"context"

	syreenconfig "syreen/config"

	"cosmossdk.io/math"
	"syreen/x/perps/types"
)

var _ types.MsgServer = Keeper{}

func (k Keeper) OpenPosition(ctx context.Context, msg *types.MsgOpenPosition) (*types.MsgOpenPositionResponse, error) {
	if !syreenconfig.IsModuleEnabled("perps") {
		return nil, syreenconfig.ErrModuleDisabled("perps")
	}
	return k.ExecuteOpenPosition(ctx, msg.Sender, msg.MarketID, msg.Side, msg.Margin, msg.Leverage)
}

func (k Keeper) ClosePosition(ctx context.Context, msg *types.MsgClosePosition) (*types.MsgClosePositionResponse, error) {
	if !syreenconfig.IsModuleEnabled("perps") {
		return nil, syreenconfig.ErrModuleDisabled("perps")
	}
	return k.ExecuteClosePosition(ctx, msg.Sender, msg.MarketID)
}

func (k Keeper) AddMargin(ctx context.Context, msg *types.MsgAddMargin) (*types.MsgAddMarginResponse, error) {
	if !syreenconfig.IsModuleEnabled("perps") {
		return nil, syreenconfig.ErrModuleDisabled("perps")
	}
	if err := k.AddMarginToPosition(ctx, msg.Sender, msg.MarketID, msg.Amount); err != nil {
		return nil, err
	}
	return &types.MsgAddMarginResponse{}, nil
}

func (k Keeper) RemoveMargin(ctx context.Context, msg *types.MsgRemoveMargin) (*types.MsgRemoveMarginResponse, error) {
	if !syreenconfig.IsModuleEnabled("perps") {
		return nil, syreenconfig.ErrModuleDisabled("perps")
	}
	if err := k.RemoveMarginFromPosition(ctx, msg.Sender, msg.MarketID, msg.Amount); err != nil {
		return nil, err
	}
	return &types.MsgRemoveMarginResponse{}, nil
}

func (k Keeper) CreateMarket(ctx context.Context, msg *types.MsgCreateMarket) (*types.MsgCreateMarketResponse, error) {
	if !syreenconfig.IsModuleEnabled("perps") {
		return nil, syreenconfig.ErrModuleDisabled("perps")
	}
	maxLev := msg.MaxLeverage
	if maxLev.IsNil() || maxLev.IsZero() {
		maxLev = math.LegacyNewDec(20)
	}
	marketID, err := k.CreateMarketEntry(ctx, msg.Authority, msg.BaseDenom, msg.QuoteDenom, msg.PoolID, maxLev)
	if err != nil {
		return nil, err
	}
	return &types.MsgCreateMarketResponse{MarketID: marketID}, nil
}
