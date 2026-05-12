package types

import "cosmossdk.io/errors"

var (
	ErrIdentityNotFound     = errors.Register(ModuleName, 2, "identity not found")
	ErrIdentityExists       = errors.Register(ModuleName, 3, "identity already registered")
	ErrNotVerifier          = errors.Register(ModuleName, 4, "address is not an authorized verifier")
	ErrVerifierNotActive    = errors.Register(ModuleName, 5, "verifier is not active")
	ErrInvalidLevel         = errors.Register(ModuleName, 6, "invalid verification level")
	ErrLevelTooHigh         = errors.Register(ModuleName, 7, "verification level exceeds verifier's max level")
	ErrAlreadyVerified      = errors.Register(ModuleName, 8, "identity already verified at this level or higher")
	ErrIdentityRevoked      = errors.Register(ModuleName, 9, "identity has been revoked")
	ErrNotIdentityOwner     = errors.Register(ModuleName, 10, "not the identity owner")
	ErrVerifierExists       = errors.Register(ModuleName, 11, "verifier already registered")
	ErrIdentityExpired      = errors.Register(ModuleName, 12, "identity verification has expired")
	ErrCannotDowngrade      = errors.Register(ModuleName, 13, "cannot downgrade verification level")
	ErrInvalidDocumentHash  = errors.Register(ModuleName, 14, "invalid document hash")
)
