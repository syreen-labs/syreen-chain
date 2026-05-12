package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"syreen/x/abstractaccount/types"
)

type queryServer struct {
	Keeper
}

// NewQueryServerImpl returns an implementation of the QueryServer interface
func NewQueryServerImpl(keeper Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
}

var _ types.QueryServer = queryServer{}

func (q queryServer) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	params := q.Keeper.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}

func (q queryServer) SmartAccount(ctx context.Context, req *types.QuerySmartAccountRequest) (*types.QuerySmartAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
	}

	account, found := q.Keeper.GetSmartAccount(ctx, req.Address)
	if !found {
		return nil, status.Errorf(codes.NotFound, "smart account %s not found", req.Address)
	}

	return &types.QuerySmartAccountResponse{Account: account}, nil
}

func (q queryServer) SessionKeys(ctx context.Context, req *types.QuerySessionKeysRequest) (*types.QuerySessionKeysResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.Granter == "" || req.Grantee == "" {
		return nil, status.Error(codes.InvalidArgument, "granter and grantee cannot be empty")
	}

	sessionKey, found := q.Keeper.GetSessionKey(ctx, req.Granter, req.Grantee)
	if !found {
		return nil, status.Errorf(codes.NotFound, "session key not found for granter %s and grantee %s", req.Granter, req.Grantee)
	}

	return &types.QuerySessionKeysResponse{SessionKey: sessionKey}, nil
}
