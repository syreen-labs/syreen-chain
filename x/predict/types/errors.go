package types

import "cosmossdk.io/errors"

var (
	ErrMarketNotFound     = errors.Register(ModuleName, 2, "prediction market not found")
	ErrMarketNotOpen      = errors.Register(ModuleName, 3, "market is not open for trading")
	ErrInvalidAmount      = errors.Register(ModuleName, 4, "invalid amount")
	ErrInvalidOutcome     = errors.Register(ModuleName, 5, "invalid outcome: must be yes, no, or void")
	ErrUnauthorized       = errors.Register(ModuleName, 6, "unauthorized: only resolver can resolve")
	ErrAlreadyResolved    = errors.Register(ModuleName, 7, "market already resolved")
	ErrNoPosition         = errors.Register(ModuleName, 8, "no position in this market")
	ErrInsufficientShares = errors.Register(ModuleName, 9, "insufficient shares to sell")
	ErrInsufficientFunds  = errors.Register(ModuleName, 10, "insufficient funds")
	ErrMarketVoided       = errors.Register(ModuleName, 11, "market has been voided")
	ErrNoWinnings         = errors.Register(ModuleName, 12, "no winnings to claim")
	ErrAlreadyClaimed     = errors.Register(ModuleName, 13, "winnings already claimed")
	ErrInvalidQuestion    = errors.Register(ModuleName, 14, "question must not be empty")
	ErrInvalidDenom       = errors.Register(ModuleName, 15, "invalid quote denom")
	ErrInvalidResolver    = errors.Register(ModuleName, 16, "invalid resolver address")
	ErrInvalidLiquidity   = errors.Register(ModuleName, 17, "initial liquidity must be positive")
)
