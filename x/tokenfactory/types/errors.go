package types

import "cosmossdk.io/errors"

var (
	ErrDenomExists              = errors.Register(ModuleName, 2, "denom already exists")
	ErrUnauthorized             = errors.Register(ModuleName, 3, "unauthorized: only the admin can perform this action")
	ErrInvalidDenom             = errors.Register(ModuleName, 4, "invalid denom")
	ErrInvalidCreator           = errors.Register(ModuleName, 5, "invalid creator address")
	ErrDenomDoesNotExist        = errors.Register(ModuleName, 6, "denom does not exist")
	ErrSubdenomTooLong          = errors.Register(ModuleName, 7, "subdenom too long, max 44 characters")
	ErrSubdenomTooShort         = errors.Register(ModuleName, 8, "subdenom too short, min 1 character")
	ErrInvalidAdmin             = errors.Register(ModuleName, 9, "invalid admin address")
	ErrCreationFeeNotEnough     = errors.Register(ModuleName, 10, "not enough funds to pay denom creation fee")
)
