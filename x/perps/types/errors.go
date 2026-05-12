package types

import "cosmossdk.io/errors"

var (
	ErrMarketNotFound       = errors.Register(ModuleName, 2, "perpetual market not found")
	ErrMarketNotActive      = errors.Register(ModuleName, 3, "perpetual market not active")
	ErrMarketAlreadyExists  = errors.Register(ModuleName, 4, "perpetual market already exists")
	ErrInvalidAmount        = errors.Register(ModuleName, 5, "invalid amount")
	ErrInvalidLeverage      = errors.Register(ModuleName, 6, "invalid leverage")
	ErrInsufficientMargin   = errors.Register(ModuleName, 7, "insufficient margin")
	ErrNoPosition           = errors.Register(ModuleName, 8, "no position found")
	ErrPositionAlreadyExists = errors.Register(ModuleName, 9, "position already exists for this market")
	ErrUnauthorized         = errors.Register(ModuleName, 10, "unauthorized")
	ErrMaxLeverageExceeded  = errors.Register(ModuleName, 11, "max leverage exceeded")
	ErrBelowMinMargin       = errors.Register(ModuleName, 12, "below minimum margin")
	ErrMaxOpenInterest      = errors.Register(ModuleName, 13, "max open interest exceeded")
	ErrInsufficientFunds    = errors.Register(ModuleName, 14, "insufficient funds")
	ErrPositionTooSmall     = errors.Register(ModuleName, 15, "position size too small")
	ErrLiquidationFailed    = errors.Register(ModuleName, 16, "liquidation failed")
)
