package keeper

import (
	"context"
	"fmt"

	"syreen/x/mevprotection/types"
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

func (q queryServer) CommittedTx(ctx context.Context, req *types.QueryCommittedTxRequest) (*types.QueryCommittedTxResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request: request cannot be nil")
	}
	if len(req.TxHash) == 0 {
		return nil, fmt.Errorf("invalid request: tx_hash cannot be empty")
	}
	committed, found := q.Keeper.GetCommittedTx(ctx, req.TxHash)
	return &types.QueryCommittedTxResponse{
		CommittedTx: committed,
		Found:       found,
	}, nil
}

func (q queryServer) PendingReveals(ctx context.Context, _ *types.QueryPendingRevealsRequest) (*types.QueryPendingRevealsResponse, error) {
	pending := q.Keeper.GetPendingReveals(ctx)
	return &types.QueryPendingRevealsResponse{PendingReveals: pending}, nil
}

func (q queryServer) MEVRedistribution(ctx context.Context, _ *types.QueryMEVRedistributionRequest) (*types.QueryMEVRedistributionResponse, error) {
	rewardPool := q.Keeper.GetMEVRewardPool(ctx)
	totalRedistributed := q.Keeper.GetTotalMEVRedistributed(ctx)

	return &types.QueryMEVRedistributionResponse{
		RewardPool:           rewardPool.String(),
		TotalRedistributed:   totalRedistributed.String(),
		LPSharePct:           "60",
		StakerSharePct:       "40",
		DistributionInterval: MEVDistributionInterval,
	}, nil
}
