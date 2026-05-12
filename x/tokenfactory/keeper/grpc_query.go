package keeper

import (
	"context"

	"syreen/x/tokenfactory/types"
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

func (q queryServer) DenomAuthority(ctx context.Context, req *types.QueryDenomAuthorityRequest) (*types.QueryDenomAuthorityResponse, error) {
	if req.Denom == "" {
		return nil, types.ErrInvalidDenom
	}

	authority, found := q.Keeper.GetDenomAuthorityMetadata(ctx, req.Denom)
	if !found {
		return &types.QueryDenomAuthorityResponse{Admin: ""}, nil
	}

	return &types.QueryDenomAuthorityResponse{Admin: authority.Admin}, nil
}
