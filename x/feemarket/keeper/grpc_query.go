package keeper

import (
	"context"

	"syreen/x/feemarket/types"
)

type queryServer struct {
	Keeper
}

// NewQueryServerImpl returns an implementation of the QueryServer interface.
func NewQueryServerImpl(keeper Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
}

var _ types.QueryServer = queryServer{}

func (q queryServer) Params(ctx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params := q.Keeper.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}

func (q queryServer) BaseFee(ctx context.Context, _ *types.QueryBaseFeeRequest) (*types.QueryBaseFeeResponse, error) {
	baseFee := q.Keeper.GetBaseFee(ctx)
	return &types.QueryBaseFeeResponse{BaseFee: baseFee}, nil
}

func (q queryServer) FeeState(ctx context.Context, _ *types.QueryFeeStateRequest) (*types.QueryFeeStateResponse, error) {
	feeState := q.Keeper.GetFeeState(ctx)
	return &types.QueryFeeStateResponse{FeeState: feeState}, nil
}
