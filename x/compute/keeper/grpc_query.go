package keeper

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"syreen/x/compute/types"
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

func (q queryServer) CodeInfo(ctx context.Context, req *types.QueryCodeInfoRequest) (*types.QueryCodeInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.CodeID == 0 {
		return nil, status.Error(codes.InvalidArgument, "code_id must be positive")
	}

	codeInfo, found := q.Keeper.GetCodeInfo(ctx, req.CodeID)
	if !found {
		return nil, status.Errorf(codes.NotFound, "code %d not found", req.CodeID)
	}

	return &types.QueryCodeInfoResponse{CodeInfo: codeInfo}, nil
}

func (q queryServer) ContractInfo(ctx context.Context, req *types.QueryContractInfoRequest) (*types.QueryContractInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
	}

	contractInfo, found := q.Keeper.GetContractInfo(ctx, req.Address)
	if !found {
		return nil, status.Errorf(codes.NotFound, "contract %s not found", req.Address)
	}

	return &types.QueryContractInfoResponse{ContractInfo: contractInfo}, nil
}

func (q queryServer) SmartContractState(ctx context.Context, req *types.QuerySmartContractStateRequest) (*types.QuerySmartContractStateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
	}
	if len(req.QueryMsg) == 0 {
		return nil, status.Error(codes.InvalidArgument, "query message cannot be empty")
	}

	data, err := q.Keeper.QuerySmart(ctx, req.Address, json.RawMessage(req.QueryMsg))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "smart query failed: %s", err.Error())
	}

	return &types.QuerySmartContractStateResponse{Data: data}, nil
}

func (q queryServer) RawContractState(ctx context.Context, req *types.QueryRawContractStateRequest) (*types.QueryRawContractStateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
	}
	if len(req.Key) == 0 {
		return nil, status.Error(codes.InvalidArgument, "key cannot be empty")
	}

	data := q.Keeper.QueryRaw(ctx, req.Address, req.Key)
	return &types.QueryRawContractStateResponse{Data: data}, nil
}
