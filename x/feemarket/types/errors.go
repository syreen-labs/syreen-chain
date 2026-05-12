package types

import "cosmossdk.io/errors"

var (
	ErrFeeTooLow        = errors.Register(ModuleName, 2, "fee is below the required minimum base fee")
	ErrInvalidBaseFee   = errors.Register(ModuleName, 3, "invalid base fee")
	ErrLaneNotFound     = errors.Register(ModuleName, 4, "fee lane not found")
	ErrBlockGasExceeded = errors.Register(ModuleName, 5, "block gas limit exceeded for lane")
)
