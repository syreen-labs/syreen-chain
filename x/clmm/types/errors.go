package types

import "cosmossdk.io/errors"

var (
	ErrPoolNotFound         = errors.Register(ModuleName, 2, "CL pool not found")
	ErrPoolAlreadyExists    = errors.Register(ModuleName, 3, "CL pool already exists for this pair")
	ErrPositionNotFound     = errors.Register(ModuleName, 4, "CL position not found")
	ErrInvalidAmount        = errors.Register(ModuleName, 5, "invalid amount")
	ErrInvalidTickRange     = errors.Register(ModuleName, 6, "invalid tick range")
	ErrInvalidTickSpacing   = errors.Register(ModuleName, 7, "tick not aligned with spacing")
	ErrInsufficientFunds    = errors.Register(ModuleName, 8, "insufficient funds")
	ErrUnauthorized         = errors.Register(ModuleName, 9, "unauthorized: not position owner")
	ErrInsufficientLiquidity = errors.Register(ModuleName, 10, "insufficient liquidity for swap")
	ErrSlippageExceeded     = errors.Register(ModuleName, 11, "swap output below minimum")
	ErrInvalidDenom         = errors.Register(ModuleName, 12, "invalid denom for this pool")
	ErrZeroLiquidity        = errors.Register(ModuleName, 13, "zero liquidity")
	ErrInvalidFeeRate       = errors.Register(ModuleName, 14, "invalid fee rate")
	ErrInvalidPrice         = errors.Register(ModuleName, 15, "invalid initial price")
	ErrSameDenom            = errors.Register(ModuleName, 16, "denomA and denomB must be different")
	ErrAmountBelowMin       = errors.Register(ModuleName, 17, "token amount below minimum")
	// ErrUnauthorizedPoolCreate is returned when a non-authority account attempts
	// to create a CL pool. Pool creation is intentionally governance-gated; this is
	// a distinct, clear error (the generic ErrUnauthorized "not position owner"
	// message was misleading on the pool-creation path).
	ErrUnauthorizedPoolCreate = errors.Register(ModuleName, 18, "only chain authority (governance) may create CL pools")
)
