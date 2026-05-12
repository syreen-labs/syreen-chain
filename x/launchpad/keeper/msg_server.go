package keeper

import (
	"context"

	syreenconfig "syreen/config"
	"syreen/x/launchpad/types"
)

var _ types.MsgServer = Keeper{}

func (k Keeper) CreateLaunch(ctx context.Context, msg *types.MsgCreateLaunch) (*types.MsgCreateLaunchResponse, error) {
	if !syreenconfig.IsModuleEnabled("launchpad") {
		return nil, syreenconfig.ErrModuleDisabled("launchpad")
	}
	launchID, err := k.ExecuteCreateLaunch(ctx, msg)
	if err != nil { return nil, err }
	return &types.MsgCreateLaunchResponse{LaunchID: launchID}, nil
}

func (k Keeper) Contribute(ctx context.Context, msg *types.MsgContribute) (*types.MsgContributeResponse, error) {
	if !syreenconfig.IsModuleEnabled("launchpad") {
		return nil, syreenconfig.ErrModuleDisabled("launchpad")
	}
	if err := k.ExecuteContribute(ctx, msg.Sender, msg.LaunchID, msg.Amount); err != nil {
		return nil, err
	}
	return &types.MsgContributeResponse{}, nil
}

func (k Keeper) ClaimTokens(ctx context.Context, msg *types.MsgClaimTokens) (*types.MsgClaimTokensResponse, error) {
	if !syreenconfig.IsModuleEnabled("launchpad") {
		return nil, syreenconfig.ErrModuleDisabled("launchpad")
	}
	amount, err := k.ExecuteClaimTokens(ctx, msg.Sender, msg.LaunchID)
	if err != nil { return nil, err }
	return &types.MsgClaimTokensResponse{Amount: amount}, nil
}

func (k Keeper) ClaimRefund(ctx context.Context, msg *types.MsgClaimRefund) (*types.MsgClaimRefundResponse, error) {
	if !syreenconfig.IsModuleEnabled("launchpad") {
		return nil, syreenconfig.ErrModuleDisabled("launchpad")
	}
	amount, err := k.ExecuteClaimRefund(ctx, msg.Sender, msg.LaunchID)
	if err != nil { return nil, err }
	return &types.MsgClaimRefundResponse{Amount: amount}, nil
}

func (k Keeper) FinalizeLaunch(ctx context.Context, msg *types.MsgFinalizeLaunch) (*types.MsgFinalizeLaunchResponse, error) {
	if !syreenconfig.IsModuleEnabled("launchpad") {
		return nil, syreenconfig.ErrModuleDisabled("launchpad")
	}
	status, err := k.ExecuteFinalizeLaunch(ctx, msg.Authority, msg.LaunchID)
	if err != nil { return nil, err }
	return &types.MsgFinalizeLaunchResponse{Status: status}, nil
}
