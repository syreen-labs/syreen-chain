package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	syreenconfig "syreen/config"
	"syreen/x/intent/types"
)

type msgServer struct {
	*Keeper
}

func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = &msgServer{}

func (m msgServer) SubmitIntent(ctx context.Context, msg *types.MsgSubmitIntent) (*types.MsgSubmitIntentResponse, error) {
	if !syreenconfig.IsModuleEnabled("intent") {
		return nil, syreenconfig.ErrModuleDisabled("intent")
	}
	intentID, err := m.Keeper.SubmitIntent(ctx, msg)
	if err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"submit_intent",
		sdk.NewAttribute("creator", msg.Creator),
		sdk.NewAttribute("intent_id", intentID),
		sdk.NewAttribute("intent_type", msg.IntentType),
	))
	return &types.MsgSubmitIntentResponse{IntentID: intentID}, nil
}

func (m msgServer) RegisterSolver(ctx context.Context, msg *types.MsgRegisterSolver) (*types.MsgRegisterSolverResponse, error) {
	if !syreenconfig.IsModuleEnabled("intent") {
		return nil, syreenconfig.ErrModuleDisabled("intent")
	}
	if err := m.Keeper.RegisterSolver(ctx, msg); err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"register_solver",
		sdk.NewAttribute("address", msg.Address),
		sdk.NewAttribute("moniker", msg.Moniker),
	))
	return &types.MsgRegisterSolverResponse{}, nil
}

func (m msgServer) DeregisterSolver(ctx context.Context, msg *types.MsgDeregisterSolver) (*types.MsgDeregisterSolverResponse, error) {
	if !syreenconfig.IsModuleEnabled("intent") {
		return nil, syreenconfig.ErrModuleDisabled("intent")
	}
	if err := m.Keeper.DeregisterSolver(ctx, msg); err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"deregister_solver",
		sdk.NewAttribute("address", msg.Address),
	))
	return &types.MsgDeregisterSolverResponse{}, nil
}

func (m msgServer) SubmitSolution(ctx context.Context, msg *types.MsgSubmitSolution) (*types.MsgSubmitSolutionResponse, error) {
	if !syreenconfig.IsModuleEnabled("intent") {
		return nil, syreenconfig.ErrModuleDisabled("intent")
	}
	if err := m.Keeper.SubmitSolution(ctx, msg); err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"submit_solution",
		sdk.NewAttribute("solver_addr", msg.SolverAddr),
		sdk.NewAttribute("intent_id", msg.IntentID),
	))
	return &types.MsgSubmitSolutionResponse{}, nil
}

func (m msgServer) FulfillIntent(ctx context.Context, msg *types.MsgFulfillIntent) (*types.MsgFulfillIntentResponse, error) {
	if !syreenconfig.IsModuleEnabled("intent") {
		return nil, syreenconfig.ErrModuleDisabled("intent")
	}
	if err := m.Keeper.FulfillIntent(ctx, msg); err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"fulfill_intent",
		sdk.NewAttribute("solver_addr", msg.SolverAddr),
		sdk.NewAttribute("intent_id", msg.IntentID),
	))
	return &types.MsgFulfillIntentResponse{}, nil
}

func (m msgServer) SubmitChain(ctx context.Context, msg *types.MsgSubmitChain) (*types.MsgSubmitChainResponse, error) {
	if !syreenconfig.IsModuleEnabled("intent") {
		return nil, syreenconfig.ErrModuleDisabled("intent")
	}
	chainID, err := m.Keeper.SubmitChain(ctx, msg)
	if err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"submit_chain",
		sdk.NewAttribute("creator", msg.Creator),
		sdk.NewAttribute("chain_id", chainID),
		sdk.NewAttribute("steps", fmt.Sprintf("%d", len(msg.Steps))),
	))
	return &types.MsgSubmitChainResponse{ChainID: chainID}, nil
}

func (m msgServer) CancelChain(ctx context.Context, msg *types.MsgCancelChain) (*types.MsgCancelChainResponse, error) {
	if !syreenconfig.IsModuleEnabled("intent") {
		return nil, syreenconfig.ErrModuleDisabled("intent")
	}
	if err := m.Keeper.CancelChain(ctx, msg); err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"cancel_chain",
		sdk.NewAttribute("creator", msg.Creator),
		sdk.NewAttribute("chain_id", msg.ChainID),
	))
	return &types.MsgCancelChainResponse{}, nil
}

func (m msgServer) CreateStrategy(ctx context.Context, msg *types.MsgCreateStrategy) (*types.MsgCreateStrategyResponse, error) {
	if !syreenconfig.IsModuleEnabled("intent") {
		return nil, syreenconfig.ErrModuleDisabled("intent")
	}
	strategyID, chainID, err := m.Keeper.CreateStrategy(ctx, msg)
	if err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"create_strategy",
		sdk.NewAttribute("creator", msg.Creator),
		sdk.NewAttribute("strategy_id", strategyID),
		sdk.NewAttribute("chain_id", chainID),
		sdk.NewAttribute("template", msg.TemplateName),
	))
	return &types.MsgCreateStrategyResponse{StrategyID: strategyID, ChainID: chainID}, nil
}

func (m msgServer) CancelStrategy(ctx context.Context, msg *types.MsgCancelStrategy) (*types.MsgCancelStrategyResponse, error) {
	if !syreenconfig.IsModuleEnabled("intent") {
		return nil, syreenconfig.ErrModuleDisabled("intent")
	}
	if err := m.Keeper.CancelStrategy(ctx, msg); err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"cancel_strategy",
		sdk.NewAttribute("creator", msg.Creator),
		sdk.NewAttribute("strategy_id", msg.StrategyID),
	))
	return &types.MsgCancelStrategyResponse{}, nil
}
