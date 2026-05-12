package keeper

import (
	"context"

	"syreen/x/vault/types"
)

// msgServer wraps the Keeper to implement the MsgServer interface
type msgServer struct {
	k Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServer returns an implementation of the MsgServer interface
func NewMsgServer(k Keeper) types.MsgServer {
	return msgServer{k: k}
}

func (m msgServer) CreateVault(ctx context.Context, msg *types.MsgCreateVault) (*types.MsgCreateVaultResponse, error) {
	id, err := m.k.ExecuteCreateVault(ctx, msg.Creator, msg.Name, msg.DepositDenom, msg.StrategyType, msg.TargetPoolIDs, msg.PerformanceFee)
	if err != nil {
		return nil, err
	}
	return &types.MsgCreateVaultResponse{VaultID: id}, nil
}

func (m msgServer) DepositVault(ctx context.Context, msg *types.MsgDepositVault) (*types.MsgDepositVaultResponse, error) {
	shares, err := m.k.ExecuteDeposit(ctx, msg.Sender, msg.VaultID, msg.Amount)
	if err != nil {
		return nil, err
	}
	return &types.MsgDepositVaultResponse{SharesMinted: shares}, nil
}

func (m msgServer) WithdrawVault(ctx context.Context, msg *types.MsgWithdrawVault) (*types.MsgWithdrawVaultResponse, error) {
	amount, err := m.k.ExecuteWithdraw(ctx, msg.Sender, msg.VaultID, msg.Shares)
	if err != nil {
		return nil, err
	}
	return &types.MsgWithdrawVaultResponse{AmountReturned: amount}, nil
}

func (m msgServer) CompoundVault(ctx context.Context, msg *types.MsgCompoundVault) (*types.MsgCompoundVaultResponse, error) {
	yield, err := m.k.ExecuteCompound(ctx, msg.Sender, msg.VaultID)
	if err != nil {
		return nil, err
	}
	return &types.MsgCompoundVaultResponse{YieldGenerated: yield}, nil
}

func (m msgServer) UpdateVaultStrategy(ctx context.Context, msg *types.MsgUpdateVaultStrategy) (*types.MsgUpdateVaultStrategyResponse, error) {
	err := m.k.ExecuteUpdateStrategy(ctx, msg.Creator, msg.VaultID, msg.StrategyType, msg.TargetPoolIDs)
	if err != nil {
		return nil, err
	}
	return &types.MsgUpdateVaultStrategyResponse{}, nil
}
