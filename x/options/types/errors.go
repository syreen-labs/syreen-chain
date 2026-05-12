package types

import "cosmossdk.io/errors"

var (
	ErrOptionNotFound    = errors.Register(ModuleName, 2, "option not found")
	ErrOptionNotOpen     = errors.Register(ModuleName, 3, "option is not open for purchase")
	ErrOptionNotActive   = errors.Register(ModuleName, 4, "option is not active")
	ErrOptionExpired     = errors.Register(ModuleName, 5, "option has expired")
	ErrOptionNotExpired  = errors.Register(ModuleName, 6, "option has not expired yet")
	ErrUnauthorized      = errors.Register(ModuleName, 7, "unauthorized")
	ErrInvalidAmount     = errors.Register(ModuleName, 8, "invalid amount")
	ErrInvalidStrike     = errors.Register(ModuleName, 9, "invalid strike price")
	ErrInvalidExpiry     = errors.Register(ModuleName, 10, "invalid expiry block")
	ErrInsufficientFunds = errors.Register(ModuleName, 11, "insufficient funds")
	ErrNotProfitable     = errors.Register(ModuleName, 12, "option is not in the money")
)
