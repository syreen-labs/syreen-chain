package keeper

import (
	"context"

	syreenconfig "syreen/config"

	"cosmossdk.io/math"
	"syreen/x/lending/types"
)

var _ types.MsgServer = Keeper{}

func (k Keeper) Deposit(ctx context.Context, msg *types.MsgDeposit) (*types.MsgDepositResponse, error) {
	if !syreenconfig.IsModuleEnabled("lending") {
		return nil, syreenconfig.ErrModuleDisabled("lending")
	}
	earned, err := k.ExecuteDeposit(ctx, msg.Sender, msg.PoolID, msg.Amount)
	if err != nil { return nil, err }
	return &types.MsgDepositResponse{InterestEarned: earned}, nil
}

func (k Keeper) Withdraw(ctx context.Context, msg *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	if !syreenconfig.IsModuleEnabled("lending") {
		return nil, syreenconfig.ErrModuleDisabled("lending")
	}
	earned, err := k.ExecuteWithdraw(ctx, msg.Sender, msg.PoolID, msg.Amount)
	if err != nil { return nil, err }
	return &types.MsgWithdrawResponse{InterestEarned: earned}, nil
}

func (k Keeper) Borrow(ctx context.Context, msg *types.MsgBorrow) (*types.MsgBorrowResponse, error) {
	if !syreenconfig.IsModuleEnabled("lending") {
		return nil, syreenconfig.ErrModuleDisabled("lending")
	}
	borrowID, err := k.ExecuteBorrow(ctx, msg.Sender, msg.BorrowPoolID, msg.Amount, msg.CollateralPoolID, msg.CollateralAmount)
	if err != nil { return nil, err }
	return &types.MsgBorrowResponse{BorrowID: borrowID}, nil
}

func (k Keeper) Repay(ctx context.Context, msg *types.MsgRepay) (*types.MsgRepayResponse, error) {
	if !syreenconfig.IsModuleEnabled("lending") {
		return nil, syreenconfig.ErrModuleDisabled("lending")
	}
	amount := msg.Amount
	if amount.IsNil() { amount = math.ZeroInt() }
	interestPaid, collateralReturned, err := k.ExecuteRepay(ctx, msg.Sender, msg.BorrowID, amount)
	if err != nil { return nil, err }
	return &types.MsgRepayResponse{InterestPaid: interestPaid, CollateralReturned: collateralReturned}, nil
}

func (k Keeper) Liquidate(ctx context.Context, msg *types.MsgLiquidate) (*types.MsgLiquidateResponse, error) {
	if !syreenconfig.IsModuleEnabled("lending") {
		return nil, syreenconfig.ErrModuleDisabled("lending")
	}
	collateralSeized, debtRepaid, err := k.ExecuteLiquidate(ctx, msg.Liquidator, msg.BorrowID)
	if err != nil { return nil, err }
	return &types.MsgLiquidateResponse{CollateralSeized: collateralSeized, DebtRepaid: debtRepaid}, nil
}

func (k Keeper) CreateLendingPool(ctx context.Context, msg *types.MsgCreateLendingPool) (*types.MsgCreateLendingPoolResponse, error) {
	if !syreenconfig.IsModuleEnabled("lending") {
		return nil, syreenconfig.ErrModuleDisabled("lending")
	}
	poolID, err := k.CreateNewLendingPool(ctx, msg.Authority, msg.Denom, msg.CollateralFactor, msg.DexPoolID, msg.PriceDenom)
	if err != nil { return nil, err }
	return &types.MsgCreateLendingPoolResponse{PoolID: poolID}, nil
}
