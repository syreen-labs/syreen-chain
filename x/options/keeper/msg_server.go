package keeper

import (
	"context"

	syreenconfig "syreen/config"

	"syreen/x/options/types"
)

var _ types.MsgServer = Keeper{}

func (k Keeper) WriteOption(ctx context.Context, msg *types.MsgWriteOption) (*types.MsgWriteOptionResponse, error) {
	if !syreenconfig.IsModuleEnabled("options") {
		return nil, syreenconfig.ErrModuleDisabled("options")
	}
	id, premium, err := k.ExecuteWriteOption(ctx, msg.Writer, msg.PoolID, msg.OptionType, msg.StrikePrice, msg.Amount, msg.ExpiryBlock, msg.CustomPremium)
	if err != nil {
		return nil, err
	}
	return &types.MsgWriteOptionResponse{OptionID: id, Premium: premium}, nil
}

func (k Keeper) BuyOption(ctx context.Context, msg *types.MsgBuyOption) (*types.MsgBuyOptionResponse, error) {
	if !syreenconfig.IsModuleEnabled("options") {
		return nil, syreenconfig.ErrModuleDisabled("options")
	}
	if err := k.ExecuteBuyOption(ctx, msg.Buyer, msg.OptionID); err != nil {
		return nil, err
	}
	return &types.MsgBuyOptionResponse{}, nil
}

func (k Keeper) ExerciseOption(ctx context.Context, msg *types.MsgExerciseOption) (*types.MsgExerciseOptionResponse, error) {
	if !syreenconfig.IsModuleEnabled("options") {
		return nil, syreenconfig.ErrModuleDisabled("options")
	}
	payout, err := k.ExecuteExerciseOption(ctx, msg.Buyer, msg.OptionID)
	if err != nil {
		return nil, err
	}
	return &types.MsgExerciseOptionResponse{Payout: payout}, nil
}

func (k Keeper) CancelOption(ctx context.Context, msg *types.MsgCancelOption) (*types.MsgCancelOptionResponse, error) {
	if !syreenconfig.IsModuleEnabled("options") {
		return nil, syreenconfig.ErrModuleDisabled("options")
	}
	if err := k.ExecuteCancelOption(ctx, msg.Writer, msg.OptionID); err != nil {
		return nil, err
	}
	return &types.MsgCancelOptionResponse{}, nil
}
