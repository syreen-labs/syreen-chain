package types

import "cosmossdk.io/errors"

var (
	ErrAgentNotFound    = errors.Register(ModuleName, 2, "agent not found")
	ErrUnauthorized     = errors.Register(ModuleName, 3, "unauthorized")
	ErrAgentNotActive   = errors.Register(ModuleName, 4, "agent is not active")
	ErrAgentNotPaused   = errors.Register(ModuleName, 5, "agent is not paused")
	ErrInvalidAmount    = errors.Register(ModuleName, 6, "invalid amount")
	ErrInsufficientFunds = errors.Register(ModuleName, 7, "insufficient agent funds")
	ErrInvalidStrategy  = errors.Register(ModuleName, 8, "invalid strategy type")
	ErrMaxTradesReached = errors.Register(ModuleName, 9, "max trades per day reached")
)
