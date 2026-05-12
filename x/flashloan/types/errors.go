package types

import "cosmossdk.io/errors"

var (
	ErrPoolNotFound          = errors.Register(ModuleName, 2, "flash loan pool not found")
	ErrPoolNotActive         = errors.Register(ModuleName, 3, "flash loan pool not active")
	ErrPoolAlreadyExists     = errors.Register(ModuleName, 4, "flash loan pool already exists for this denom")
	ErrInvalidAmount         = errors.Register(ModuleName, 5, "invalid amount")
	ErrInsufficientLiquidity = errors.Register(ModuleName, 6, "insufficient flash pool liquidity")
	ErrRepaymentFailed       = errors.Register(ModuleName, 7, "flash loan repayment failed: module account balance insufficient")
	ErrUnauthorized          = errors.Register(ModuleName, 8, "unauthorized")
	ErrInsufficientFunds     = errors.Register(ModuleName, 9, "insufficient funds")
	ErrInvalidDenom          = errors.Register(ModuleName, 10, "invalid denom")
	ErrInvalidFeeRate        = errors.Register(ModuleName, 11, "invalid fee rate")
	ErrNoDeposit             = errors.Register(ModuleName, 12, "no flash pool deposit found for this address")
	ErrInsufficientShares    = errors.Register(ModuleName, 13, "insufficient flash pool shares")
	ErrInitialDepositTooSmall = errors.Register(ModuleName, 14, "initial flash pool deposit must exceed minimum liquidity")
	ErrFlashLoanInProgress    = errors.Register(ModuleName, 15, "a flash loan is already in progress for this sender")
	ErrRepaymentInsufficient  = errors.Register(ModuleName, 16, "flash loan repayment insufficient: module balance below required minimum")
)
