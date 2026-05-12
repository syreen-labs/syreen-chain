package types

import "cosmossdk.io/errors"

var (
	ErrVaultNotFound        = errors.Register(ModuleName, 2, "vault not found")
	ErrVaultNotActive       = errors.Register(ModuleName, 3, "vault not active")
	ErrInvalidAmount        = errors.Register(ModuleName, 4, "invalid amount")
	ErrInsufficientShares   = errors.Register(ModuleName, 5, "insufficient shares")
	ErrNoDeposit            = errors.Register(ModuleName, 6, "no deposit found")
	ErrUnauthorized         = errors.Register(ModuleName, 7, "unauthorized")
	ErrInsufficientFunds    = errors.Register(ModuleName, 8, "insufficient funds")
	ErrInvalidStrategy      = errors.Register(ModuleName, 9, "invalid strategy type")
	ErrInvalidName          = errors.Register(ModuleName, 10, "invalid vault name")
	ErrInvalidDenom         = errors.Register(ModuleName, 11, "invalid deposit denom")
	ErrInvalidFee           = errors.Register(ModuleName, 12, "invalid performance fee")
	ErrNothingToCompound    = errors.Register(ModuleName, 13, "nothing to compound")
	ErrInitialDepositTooSmall = errors.Register(ModuleName, 14, "initial vault deposit must exceed minimum liquidity (1000)")
	ErrWithdrawTooEarly     = errors.Register(ModuleName, 15, "must wait at least 100 blocks after last deposit before withdrawing")
)
