package types

import "cosmossdk.io/errors"

var (
	ErrInvalidCommitment = errors.Register(ModuleName, 2, "invalid commitment: tx hash is empty or malformed")
	ErrRevealMismatch    = errors.Register(ModuleName, 3, "reveal does not match prior commitment")
	ErrRevealExpired     = errors.Register(ModuleName, 4, "reveal window has expired for this commitment")
	ErrDuplicateCommit   = errors.Register(ModuleName, 5, "duplicate commitment: tx hash already committed")
	ErrMEVDetected       = errors.Register(ModuleName, 6, "MEV extraction detected: suspicious transaction reordering")
	ErrRevealTooEarly    = errors.Register(ModuleName, 7, "reveal submitted too early: must wait at least half the commit window")
)
