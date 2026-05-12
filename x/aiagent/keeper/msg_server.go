package keeper

import (
	"context"

	syreenconfig "syreen/config"

	"syreen/x/aiagent/types"
)

var _ types.MsgServer = Keeper{}

func (k Keeper) CreateAgent(ctx context.Context, msg *types.MsgCreateAgent) (*types.MsgCreateAgentResponse, error) {
	if !syreenconfig.IsModuleEnabled("aiagent") {
		return nil, syreenconfig.ErrModuleDisabled("aiagent")
	}
	id, addr, err := k.ExecuteCreateAgent(ctx, msg.Owner, msg.Name, msg.StrategyType, msg.Config, msg.InitialFunds)
	if err != nil {
		return nil, err
	}
	return &types.MsgCreateAgentResponse{AgentID: id, AgentAddress: addr}, nil
}

func (k Keeper) FundAgent(ctx context.Context, msg *types.MsgFundAgent) (*types.MsgFundAgentResponse, error) {
	if !syreenconfig.IsModuleEnabled("aiagent") {
		return nil, syreenconfig.ErrModuleDisabled("aiagent")
	}
	if err := k.ExecuteFundAgent(ctx, msg.Owner, msg.AgentID, msg.Amount); err != nil {
		return nil, err
	}
	return &types.MsgFundAgentResponse{}, nil
}

func (k Keeper) WithdrawAgentFunds(ctx context.Context, msg *types.MsgWithdrawAgentFunds) (*types.MsgWithdrawAgentFundsResponse, error) {
	if !syreenconfig.IsModuleEnabled("aiagent") {
		return nil, syreenconfig.ErrModuleDisabled("aiagent")
	}
	amount, err := k.ExecuteWithdrawAgentFunds(ctx, msg.Owner, msg.AgentID, msg.Amount)
	if err != nil {
		return nil, err
	}
	return &types.MsgWithdrawAgentFundsResponse{AmountReturned: amount}, nil
}

func (k Keeper) PauseAgent(ctx context.Context, msg *types.MsgPauseAgent) (*types.MsgPauseAgentResponse, error) {
	if !syreenconfig.IsModuleEnabled("aiagent") {
		return nil, syreenconfig.ErrModuleDisabled("aiagent")
	}
	if err := k.ExecutePauseAgent(ctx, msg.Owner, msg.AgentID); err != nil {
		return nil, err
	}
	return &types.MsgPauseAgentResponse{}, nil
}

func (k Keeper) ResumeAgent(ctx context.Context, msg *types.MsgResumeAgent) (*types.MsgResumeAgentResponse, error) {
	if !syreenconfig.IsModuleEnabled("aiagent") {
		return nil, syreenconfig.ErrModuleDisabled("aiagent")
	}
	if err := k.ExecuteResumeAgent(ctx, msg.Owner, msg.AgentID); err != nil {
		return nil, err
	}
	return &types.MsgResumeAgentResponse{}, nil
}

func (k Keeper) UpdateAgentStrategy(ctx context.Context, msg *types.MsgUpdateAgentStrategy) (*types.MsgUpdateAgentStrategyResponse, error) {
	if !syreenconfig.IsModuleEnabled("aiagent") {
		return nil, syreenconfig.ErrModuleDisabled("aiagent")
	}
	if err := k.ExecuteUpdateStrategy(ctx, msg.Owner, msg.AgentID, msg.Config); err != nil {
		return nil, err
	}
	return &types.MsgUpdateAgentStrategyResponse{}, nil
}
