package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/tokenfactory/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) CreateDenom(ctx context.Context, msg *types.MsgCreateDenom) (*types.MsgCreateDenomResponse, error) {
	newDenom, err := m.Keeper.CreateDenom(ctx, msg.Sender, msg.Subdenom)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"create_denom",
			sdk.NewAttribute("creator", msg.Sender),
			sdk.NewAttribute("new_token_denom", newDenom),
		),
	})

	return &types.MsgCreateDenomResponse{NewTokenDenom: newDenom}, nil
}

func (m msgServer) Mint(ctx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	err := m.Keeper.Mint(ctx, msg.Sender, msg.Amount, msg.MintTo)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"tf_mint",
			sdk.NewAttribute("mint_to_address", msg.MintTo),
			sdk.NewAttribute("amount", msg.Amount.String()),
		),
	})

	return &types.MsgMintResponse{}, nil
}

func (m msgServer) Burn(ctx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	err := m.Keeper.Burn(ctx, msg.Sender, msg.Amount, msg.BurnFrom)
	if err != nil {
		return nil, err
	}

	burnFrom := msg.Sender
	if msg.BurnFrom != "" {
		burnFrom = msg.BurnFrom
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"tf_burn",
			sdk.NewAttribute("sender", msg.Sender),
			sdk.NewAttribute("burn_from_address", burnFrom),
			sdk.NewAttribute("amount", msg.Amount.String()),
		),
	})

	return &types.MsgBurnResponse{}, nil
}

func (m msgServer) ChangeAdmin(ctx context.Context, msg *types.MsgChangeAdmin) (*types.MsgChangeAdminResponse, error) {
	err := m.Keeper.ChangeAdmin(ctx, msg.Sender, msg.Denom, msg.NewAdmin)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"change_admin",
			sdk.NewAttribute("denom", msg.Denom),
			sdk.NewAttribute("new_admin", msg.NewAdmin),
		),
	})

	return &types.MsgChangeAdminResponse{}, nil
}
