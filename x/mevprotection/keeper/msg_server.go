package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/mevprotection/types"
)

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) CommitTx(ctx context.Context, msg *types.MsgCommitTx) (*types.MsgCommitTxResponse, error) {
	if err := m.Keeper.CommitTx(ctx, msg.Sender, msg.TxHash, msg.EncryptedTx); err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"commit_tx",
		sdk.NewAttribute("sender", msg.Sender),
	))
	return &types.MsgCommitTxResponse{}, nil
}

func (m msgServer) RevealTx(ctx context.Context, msg *types.MsgRevealTx) (*types.MsgRevealTxResponse, error) {
	if err := m.Keeper.RevealTx(ctx, msg.Sender, msg.CommitHash, msg.TxBody, msg.Nonce); err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"reveal_tx",
		sdk.NewAttribute("sender", msg.Sender),
	))
	return &types.MsgRevealTxResponse{}, nil
}
