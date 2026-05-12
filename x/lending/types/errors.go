package types

import "cosmossdk.io/errors"

var (
	ErrPoolNotFound        = errors.Register(ModuleName, 2, "lending pool not found")
	ErrPoolNotActive       = errors.Register(ModuleName, 3, "lending pool not active")
	ErrPoolAlreadyExists   = errors.Register(ModuleName, 4, "lending pool already exists for this denom")
	ErrInvalidAmount       = errors.Register(ModuleName, 5, "invalid amount")
	ErrInsufficientDeposit = errors.Register(ModuleName, 6, "insufficient deposit balance")
	ErrInsufficientLiquidity = errors.Register(ModuleName, 7, "insufficient pool liquidity for borrow")
	ErrBorrowNotFound      = errors.Register(ModuleName, 8, "borrow position not found")
	ErrUnauthorized        = errors.Register(ModuleName, 9, "unauthorized")
	ErrInsufficientFunds   = errors.Register(ModuleName, 10, "insufficient funds")
	ErrOverCollateral      = errors.Register(ModuleName, 11, "borrow exceeds collateral limit")
	ErrNoDeposit           = errors.Register(ModuleName, 12, "no deposit found")
	ErrBorrowNotOwner      = errors.Register(ModuleName, 13, "not the borrow owner")
	ErrHealthFactorTooLow  = errors.Register(ModuleName, 14, "health factor too low for this action")
	ErrLiquidationNotEligible = errors.Register(ModuleName, 15, "position not eligible for liquidation")
	ErrInsufficientCollateral = errors.Register(ModuleName, 16, "collateral value insufficient for borrow amount")
	ErrNoPriceOracle          = errors.Register(ModuleName, 17, "cannot price collateral: no oracle data available")
	ErrInvalidParam           = errors.Register(ModuleName, 18, "invalid parameter value")
	ErrPriceDeviation         = errors.Register(ModuleName, 19, "spot price deviates from TWAP beyond allowed bound")
	ErrLiquidationDelay       = errors.Register(ModuleName, 20, "position has not been unhealthy long enough to liquidate")
)
