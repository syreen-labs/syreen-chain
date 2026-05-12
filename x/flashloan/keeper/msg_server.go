package keeper

import (
	"context"

	syreenconfig "syreen/config"
	"syreen/x/flashloan/types"
)

var _ types.MsgServer = Keeper{}

func (k Keeper) FlashLoan(ctx context.Context, msg *types.MsgFlashLoan) (*types.MsgFlashLoanResponse, error) {
	if !syreenconfig.IsModuleEnabled("flashloan") {
		return nil, syreenconfig.ErrModuleDisabled("flashloan")
	}
	fee, err := k.ExecuteFlashLoan(ctx, msg.Sender, msg.Denom, msg.Amount)
	if err != nil {
		return nil, err
	}
	return &types.MsgFlashLoanResponse{Fee: fee}, nil
}

func (k Keeper) CreateFlashPool(ctx context.Context, msg *types.MsgCreateFlashPool) (*types.MsgCreateFlashPoolResponse, error) {
	if !syreenconfig.IsModuleEnabled("flashloan") {
		return nil, syreenconfig.ErrModuleDisabled("flashloan")
	}
	if err := k.CreatePool(ctx, msg.Authority, msg.Denom, msg.FeeRate); err != nil {
		return nil, err
	}
	return &types.MsgCreateFlashPoolResponse{}, nil
}

func (k Keeper) FundFlashPool(ctx context.Context, msg *types.MsgFundFlashPool) (*types.MsgFundFlashPoolResponse, error) {
	if !syreenconfig.IsModuleEnabled("flashloan") {
		return nil, syreenconfig.ErrModuleDisabled("flashloan")
	}
	if err := k.FundPool(ctx, msg.Sender, msg.Denom, msg.Amount); err != nil {
		return nil, err
	}
	return &types.MsgFundFlashPoolResponse{}, nil
}

func (k Keeper) WithdrawFlashPool(ctx context.Context, msg *types.MsgWithdrawFlashPool) (*types.MsgWithdrawFlashPoolResponse, error) {
	if !syreenconfig.IsModuleEnabled("flashloan") {
		return nil, syreenconfig.ErrModuleDisabled("flashloan")
	}
	amount, err := k.ExecuteWithdrawFlashPool(ctx, msg.Sender, msg.Denom, msg.Shares)
	if err != nil {
		return nil, err
	}
	return &types.MsgWithdrawFlashPoolResponse{AmountReturned: amount}, nil
}
