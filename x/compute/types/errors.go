package types

import "cosmossdk.io/errors"

var (
	ErrCodeNotFound     = errors.Register(ModuleName, 2, "code not found")
	ErrContractNotFound = errors.Register(ModuleName, 3, "contract not found")
	ErrUnauthorized     = errors.Register(ModuleName, 4, "unauthorized")
	ErrInvalidMsg       = errors.Register(ModuleName, 5, "invalid message")
	ErrOutOfGas         = errors.Register(ModuleName, 6, "out of gas for contract execution")
	ErrCodeTooLarge     = errors.Register(ModuleName, 7, "code size exceeds limit")
	ErrInvalidCreator   = errors.Register(ModuleName, 8, "invalid creator address")
	ErrInvalidAdmin     = errors.Register(ModuleName, 9, "invalid admin address")
	ErrInvalidLabel     = errors.Register(ModuleName, 10, "invalid label")
	ErrDuplicateCode    = errors.Register(ModuleName, 11, "duplicate code")
	ErrMigrationFailed  = errors.Register(ModuleName, 12, "contract migration failed")
	ErrNoSuchCode       = errors.Register(ModuleName, 13, "no such code")
	ErrEmpty            = errors.Register(ModuleName, 14, "empty")
	ErrUploadDenied     = errors.Register(ModuleName, 15, "code upload not allowed")
	ErrInstantiateDenied = errors.Register(ModuleName, 16, "instantiate not allowed")
)
