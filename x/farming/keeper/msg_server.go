package keeper

import (
	"context"

	syreenconfig "syreen/config"
	"syreen/x/farming/types"
)

var _ types.MsgServer = Keeper{}

func (k Keeper) Stake(ctx context.Context, msg *types.MsgStake) (*types.MsgStakeResponse, error) {
	if !syreenconfig.IsModuleEnabled("farming") {
		return nil, syreenconfig.ErrModuleDisabled("farming")
	}
	pendingReward, err := k.StakeLP(ctx, msg.Sender, msg.PoolID, msg.Amount)
	if err != nil {
		return nil, err
	}
	return &types.MsgStakeResponse{PendingReward: pendingReward}, nil
}

func (k Keeper) Unstake(ctx context.Context, msg *types.MsgUnstake) (*types.MsgUnstakeResponse, error) {
	if !syreenconfig.IsModuleEnabled("farming") {
		return nil, syreenconfig.ErrModuleDisabled("farming")
	}
	claimedReward, err := k.UnstakeLP(ctx, msg.Sender, msg.PoolID, msg.Amount)
	if err != nil {
		return nil, err
	}
	return &types.MsgUnstakeResponse{ClaimedReward: claimedReward}, nil
}

func (k Keeper) ClaimReward(ctx context.Context, msg *types.MsgClaimReward) (*types.MsgClaimRewardResponse, error) {
	if !syreenconfig.IsModuleEnabled("farming") {
		return nil, syreenconfig.ErrModuleDisabled("farming")
	}
	amount, err := k.ClaimFarmReward(ctx, msg.Sender, msg.PoolID)
	if err != nil {
		return nil, err
	}
	return &types.MsgClaimRewardResponse{Amount: amount}, nil
}

func (k Keeper) CreateFarm(ctx context.Context, msg *types.MsgCreateFarm) (*types.MsgCreateFarmResponse, error) {
	if !syreenconfig.IsModuleEnabled("farming") {
		return nil, syreenconfig.ErrModuleDisabled("farming")
	}
	if err := k.CreateFarmPool(ctx, msg.Authority, msg.PoolID, msg.LPDenom, msg.RewardPerBlock, msg.StartBlock, msg.EndBlock); err != nil {
		return nil, err
	}
	return &types.MsgCreateFarmResponse{}, nil
}

func (k Keeper) UpdateFarm(ctx context.Context, msg *types.MsgUpdateFarm) (*types.MsgUpdateFarmResponse, error) {
	if !syreenconfig.IsModuleEnabled("farming") {
		return nil, syreenconfig.ErrModuleDisabled("farming")
	}
	if err := k.UpdateFarmConfig(ctx, msg.Authority, msg.PoolID, msg.RewardPerBlock, msg.Active); err != nil {
		return nil, err
	}
	return &types.MsgUpdateFarmResponse{}, nil
}
