package types

import "cosmossdk.io/errors"

var (
	ErrAccountExists          = errors.Register(ModuleName, 2, "smart account already exists")
	ErrUnauthorized           = errors.Register(ModuleName, 3, "unauthorized: caller is not authorized to perform this action")
	ErrSessionExpired         = errors.Register(ModuleName, 4, "session key has expired")
	ErrSessionLimitExceeded   = errors.Register(ModuleName, 5, "session key spend limit exceeded")
	ErrRecoveryInProgress     = errors.Register(ModuleName, 6, "recovery is already in progress for this account")
	ErrInsufficientGuardians  = errors.Register(ModuleName, 7, "insufficient guardian approvals to execute recovery")
	ErrBatchTooLarge          = errors.Register(ModuleName, 8, "batch size exceeds maximum allowed")
	ErrSponsorLimitExceeded   = errors.Register(ModuleName, 9, "gas sponsorship limit exceeded")
	ErrInvalidAddress         = errors.Register(ModuleName, 10, "invalid address")
	ErrAccountNotFound        = errors.Register(ModuleName, 11, "smart account not found")
	ErrSessionKeyNotFound     = errors.Register(ModuleName, 12, "session key not found")
	ErrRecoveryNotFound       = errors.Register(ModuleName, 13, "no active recovery request found")
	ErrRecoveryDelayNotMet    = errors.Register(ModuleName, 14, "recovery delay period has not elapsed")
	ErrInvalidAccountType     = errors.Register(ModuleName, 15, "invalid account type")
	ErrInvalidThreshold       = errors.Register(ModuleName, 16, "threshold must be > 0 and <= number of owners")
	ErrTooManyGuardians       = errors.Register(ModuleName, 17, "number of guardians exceeds maximum allowed")
	ErrSponsorNotFound        = errors.Register(ModuleName, 18, "gas sponsor not found")
	ErrGasSponsorshipDisabled = errors.Register(ModuleName, 19, "gas sponsorship is disabled")
	ErrDuplicateApproval      = errors.Register(ModuleName, 20, "guardian has already approved this recovery")
	ErrNotGuardian            = errors.Register(ModuleName, 21, "address is not a guardian for this account")
	ErrRecoveryConfigNotFound = errors.Register(ModuleName, 22, "recovery config not found for this account")
	ErrSessionKeyDuration     = errors.Register(ModuleName, 23, "session key duration exceeds maximum allowed")
	ErrPermissionDenied       = errors.Register(ModuleName, 24, "session key does not have permission for this message type")
)
