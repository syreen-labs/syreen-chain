package types

import "cosmossdk.io/errors"

var (
	ErrPortfolioNotFound      = errors.Register(ModuleName, 2, "portfolio not found")
	ErrInvalidAmount          = errors.Register(ModuleName, 3, "invalid amount")
	ErrInvalidAddress         = errors.Register(ModuleName, 4, "invalid address")
	ErrCompetitionNotFound    = errors.Register(ModuleName, 5, "competition not found")
	ErrCompetitionNotActive   = errors.Register(ModuleName, 6, "competition is not active")
	ErrCompetitionFull        = errors.Register(ModuleName, 7, "competition is full")
	ErrAlreadyJoined          = errors.Register(ModuleName, 8, "already joined competition")
	ErrCompetitionNotEnded    = errors.Register(ModuleName, 9, "competition has not ended")
	ErrInsufficientFunds      = errors.Register(ModuleName, 10, "insufficient funds for entry fee")
	ErrUnauthorized           = errors.Register(ModuleName, 11, "unauthorized")
	ErrInvalidCompetition     = errors.Register(ModuleName, 12, "invalid competition parameters")
	ErrCompetitionAlreadyEnded = errors.Register(ModuleName, 13, "competition already ended")
	ErrPrizesAlreadyDistributed = errors.Register(ModuleName, 14, "prizes already distributed")
	ErrNotJoined              = errors.Register(ModuleName, 15, "not joined in competition")
	ErrInvalidTradeType       = errors.Register(ModuleName, 16, "invalid trade type")
)
