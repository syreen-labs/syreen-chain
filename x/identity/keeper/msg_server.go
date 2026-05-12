package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/identity/types"
)

type msgServer struct {
	*Keeper
}

func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = &msgServer{}

func (m msgServer) RegisterIdentity(ctx context.Context, msg *types.MsgRegisterIdentity) (*types.MsgRegisterIdentityResponse, error) {
	if err := m.Keeper.CreateIdentity(ctx, msg); err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	// Nationality is PII — do not expose in events. Emit only non-sensitive fields.
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"register_identity",
		sdk.NewAttribute("address", msg.Address),
		sdk.NewAttribute("level", msg.Level),
	))
	return &types.MsgRegisterIdentityResponse{}, nil
}

func (m msgServer) VerifyIdentity(ctx context.Context, msg *types.MsgVerifyIdentity) (*types.MsgVerifyIdentityResponse, error) {
	if err := m.Keeper.VerifyIdentity(ctx, msg); err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"verify_identity",
		sdk.NewAttribute("verifier", msg.Verifier),
		sdk.NewAttribute("address", msg.Address),
		sdk.NewAttribute("level", msg.Level),
	))
	return &types.MsgVerifyIdentityResponse{}, nil
}

func (m msgServer) RejectIdentity(ctx context.Context, msg *types.MsgRejectIdentity) (*types.MsgRejectIdentityResponse, error) {
	if err := m.Keeper.RejectIdentity(ctx, msg); err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"reject_identity",
		sdk.NewAttribute("verifier", msg.Verifier),
		sdk.NewAttribute("address", msg.Address),
		sdk.NewAttribute("reason", msg.Reason),
	))
	return &types.MsgRejectIdentityResponse{}, nil
}

func (m msgServer) RevokeIdentity(ctx context.Context, msg *types.MsgRevokeIdentity) (*types.MsgRevokeIdentityResponse, error) {
	if err := m.Keeper.RevokeIdentity(ctx, msg); err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"revoke_identity",
		sdk.NewAttribute("authority", msg.Authority),
		sdk.NewAttribute("address", msg.Address),
		sdk.NewAttribute("reason", msg.Reason),
	))
	return &types.MsgRevokeIdentityResponse{}, nil
}

func (m msgServer) UpdateIdentity(ctx context.Context, msg *types.MsgUpdateIdentity) (*types.MsgUpdateIdentityResponse, error) {
	if err := m.Keeper.UpdateIdentity(ctx, msg); err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"update_identity",
		sdk.NewAttribute("address", msg.Address),
	))
	return &types.MsgUpdateIdentityResponse{}, nil
}

func (m msgServer) RegisterVerifier(ctx context.Context, msg *types.MsgRegisterVerifier) (*types.MsgRegisterVerifierResponse, error) {
	if err := m.Keeper.RegisterVerifier(ctx, msg); err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"register_verifier",
		sdk.NewAttribute("authority", msg.Authority),
		sdk.NewAttribute("verifier", msg.Verifier),
		sdk.NewAttribute("name", msg.Name),
		sdk.NewAttribute("max_level", msg.MaxLevel),
	))
	return &types.MsgRegisterVerifierResponse{}, nil
}

func (m msgServer) DeactivateVerifier(ctx context.Context, msg *types.MsgDeactivateVerifier) (*types.MsgDeactivateVerifierResponse, error) {
	if err := m.Keeper.DeactivateVerifier(ctx, msg); err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"deactivate_verifier",
		sdk.NewAttribute("authority", msg.Authority),
		sdk.NewAttribute("verifier", msg.Verifier),
	))
	return &types.MsgDeactivateVerifierResponse{}, nil
}

func (m msgServer) IncrementBookings(ctx context.Context, msg *types.MsgIncrementBookings) (*types.MsgIncrementBookingsResponse, error) {
	if err := m.Keeper.IncrementBookings(ctx, msg); err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"increment_bookings",
		sdk.NewAttribute("authority", msg.Authority),
		sdk.NewAttribute("address", msg.Address),
	))
	return &types.MsgIncrementBookingsResponse{}, nil
}
