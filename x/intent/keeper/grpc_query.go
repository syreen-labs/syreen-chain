package keeper

import (
	"context"
	"fmt"

	"syreen/x/intent/types"
)

type queryServer struct {
	Keeper
}

func NewQueryServerImpl(keeper Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
}

var _ types.QueryServer = queryServer{}

func (q queryServer) Params(ctx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params := q.Keeper.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}

func (q queryServer) Intent(ctx context.Context, req *types.QueryIntentRequest) (*types.QueryIntentResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request: request cannot be nil")
	}
	if req.IntentID == "" {
		return nil, fmt.Errorf("invalid request: intent_id cannot be empty")
	}
	intent, found := q.Keeper.GetIntent(ctx, req.IntentID)
	return &types.QueryIntentResponse{
		Intent: intent,
		Found:  found,
	}, nil
}

func (q queryServer) Intents(ctx context.Context, _ *types.QueryIntentsRequest) (*types.QueryIntentsResponse, error) {
	intents := q.Keeper.GetAllIntents(ctx)
	if intents == nil {
		intents = []types.Intent{}
	}
	return &types.QueryIntentsResponse{Intents: intents}, nil
}

func (q queryServer) Solver(ctx context.Context, req *types.QuerySolverRequest) (*types.QuerySolverResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request: request cannot be nil")
	}
	if req.Address == "" {
		return nil, fmt.Errorf("invalid request: address cannot be empty")
	}
	solver, found := q.Keeper.GetSolver(ctx, req.Address)
	return &types.QuerySolverResponse{
		Solver: solver,
		Found:  found,
	}, nil
}

func (q queryServer) Solutions(ctx context.Context, req *types.QuerySolutionsRequest) (*types.QuerySolutionsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request: request cannot be nil")
	}
	if req.IntentID == "" {
		return nil, fmt.Errorf("invalid request: intent_id cannot be empty")
	}
	solutions := q.Keeper.GetSolutionsForIntent(ctx, req.IntentID)
	return &types.QuerySolutionsResponse{Solutions: solutions}, nil
}

func (q queryServer) IntentsByType(ctx context.Context, req *types.QueryIntentsByTypeRequest) (*types.QueryIntentsByTypeResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request: request cannot be nil")
	}
	if req.IntentType == "" {
		return nil, fmt.Errorf("invalid request: intent_type cannot be empty")
	}
	intents := q.Keeper.GetIntentsByType(ctx, req.IntentType)
	return &types.QueryIntentsByTypeResponse{Intents: intents}, nil
}

func (q queryServer) StrategyTemplates(_ context.Context, _ *types.QueryStrategyTemplatesRequest) (*types.QueryStrategyTemplatesResponse, error) {
	return &types.QueryStrategyTemplatesResponse{Templates: types.GetAvailableTemplates()}, nil
}

func (q queryServer) Strategies(ctx context.Context, req *types.QueryStrategiesRequest) (*types.QueryStrategiesResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request: request cannot be nil")
	}
	if req.Address == "" {
		return nil, fmt.Errorf("invalid request: address cannot be empty")
	}
	strategies := q.Keeper.GetStrategiesByCreator(ctx, req.Address)
	if strategies == nil {
		strategies = []types.Strategy{}
	}
	return &types.QueryStrategiesResponse{Strategies: strategies}, nil
}

func (q queryServer) Strategy(ctx context.Context, req *types.QueryStrategyRequest) (*types.QueryStrategyResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request: request cannot be nil")
	}
	if req.StrategyID == "" {
		return nil, fmt.Errorf("invalid request: strategy_id cannot be empty")
	}
	strategy, found := q.Keeper.GetStrategy(ctx, req.StrategyID)
	if !found {
		return &types.QueryStrategyResponse{Found: false}, nil
	}
	chain, _ := q.Keeper.GetChain(ctx, strategy.ChainID)
	return &types.QueryStrategyResponse{
		Strategy: strategy,
		Chain:    chain,
		Found:    true,
	}, nil
}
