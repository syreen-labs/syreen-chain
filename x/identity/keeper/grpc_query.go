package keeper

import (
	"context"
	"fmt"

	"syreen/x/identity/types"
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

func (q queryServer) Identity(ctx context.Context, req *types.QueryIdentityRequest) (*types.QueryIdentityResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request: request cannot be nil")
	}
	if req.Address == "" {
		return nil, fmt.Errorf("invalid request: address cannot be empty")
	}
	identity, found := q.Keeper.GetIdentity(ctx, req.Address)
	return &types.QueryIdentityResponse{
		Identity: identity,
		Found:    found,
	}, nil
}

func (q queryServer) IdentitiesByVerifier(ctx context.Context, req *types.QueryIdentitiesByVerifierRequest) (*types.QueryIdentitiesByVerifierResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request: request cannot be nil")
	}
	if req.Verifier == "" {
		return nil, fmt.Errorf("invalid request: verifier cannot be empty")
	}
	identities := q.Keeper.GetIdentitiesByVerifier(ctx, req.Verifier)
	return &types.QueryIdentitiesByVerifierResponse{Identities: identities}, nil
}

func (q queryServer) Verifier(ctx context.Context, req *types.QueryVerifierRequest) (*types.QueryVerifierResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request: request cannot be nil")
	}
	if req.Address == "" {
		return nil, fmt.Errorf("invalid request: address cannot be empty")
	}
	verifier, found := q.Keeper.GetVerifier(ctx, req.Address)
	return &types.QueryVerifierResponse{
		Verifier: verifier,
		Found:    found,
	}, nil
}

func (q queryServer) VerificationStatus(ctx context.Context, req *types.QueryVerificationStatusRequest) (*types.QueryVerificationStatusResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request: request cannot be nil")
	}
	if req.Address == "" {
		return nil, fmt.Errorf("invalid request: address cannot be empty")
	}

	identity, found := q.Keeper.GetIdentity(ctx, req.Address)
	if !found {
		return &types.QueryVerificationStatusResponse{
			Address:    req.Address,
			Verified:   false,
			Level:      types.VerificationNone,
			Status:     types.StatusPending,
			TrustScore: 0,
		}, nil
	}

	verified := q.Keeper.IsVerified(ctx, req.Address)

	return &types.QueryVerificationStatusResponse{
		Address:    req.Address,
		Verified:   verified,
		Level:      identity.Level,
		Status:     identity.Status,
		TrustScore: identity.TrustScore,
	}, nil
}
