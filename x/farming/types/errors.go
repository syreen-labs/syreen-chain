package types

import "cosmossdk.io/errors"

var (
	ErrFarmNotFound      = errors.Register(ModuleName, 2, "farming pool not found")
	ErrFarmNotActive     = errors.Register(ModuleName, 3, "farming pool not active")
	ErrFarmAlreadyExists = errors.Register(ModuleName, 4, "farming pool already exists")
	ErrInvalidAmount     = errors.Register(ModuleName, 5, "invalid amount")
	ErrInsufficientStake = errors.Register(ModuleName, 6, "insufficient staked amount")
	ErrNoPosition        = errors.Register(ModuleName, 7, "no farming position found")
	ErrNoPendingReward   = errors.Register(ModuleName, 8, "no pending rewards to claim")
	ErrUnauthorized      = errors.Register(ModuleName, 9, "unauthorized")
	ErrInsufficientFunds = errors.Register(ModuleName, 10, "insufficient reward funds")
	ErrFarmNotStarted    = errors.Register(ModuleName, 11, "farming pool has not started yet")
	ErrFarmEnded         = errors.Register(ModuleName, 12, "farming pool has ended")
)
