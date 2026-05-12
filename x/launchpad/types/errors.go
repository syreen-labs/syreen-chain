package types

import "cosmossdk.io/errors"

var (
	ErrLaunchNotFound       = errors.Register(ModuleName, 2, "launch not found")
	ErrInvalidAmount        = errors.Register(ModuleName, 3, "invalid amount")
	ErrLaunchNotActive      = errors.Register(ModuleName, 4, "launch is not active")
	ErrLaunchNotStarted     = errors.Register(ModuleName, 5, "launch has not started yet")
	ErrLaunchEnded          = errors.Register(ModuleName, 6, "launch has already ended")
	ErrHardCapExceeded      = errors.Register(ModuleName, 7, "contribution would exceed hard cap")
	ErrMaxPerWalletExceeded = errors.Register(ModuleName, 8, "contribution would exceed max per wallet")
	ErrLaunchNotSuccessful  = errors.Register(ModuleName, 9, "launch was not successful")
	ErrLaunchNotFailed      = errors.Register(ModuleName, 10, "launch did not fail")
	ErrAlreadyClaimed       = errors.Register(ModuleName, 11, "already claimed")
	ErrNoContribution       = errors.Register(ModuleName, 12, "no contribution found")
	ErrUnauthorized         = errors.Register(ModuleName, 13, "unauthorized")
	ErrInsufficientFunds    = errors.Register(ModuleName, 14, "insufficient funds")
	ErrInvalidLaunchParams  = errors.Register(ModuleName, 15, "invalid launch parameters")
	ErrLaunchAlreadyFinalized = errors.Register(ModuleName, 16, "launch already finalized")
	ErrVestingNotStarted    = errors.Register(ModuleName, 17, "vesting has not started")
	ErrNothingToClaim       = errors.Register(ModuleName, 18, "nothing to claim")
	ErrLaunchNotEnded       = errors.Register(ModuleName, 19, "launch has not ended yet")
)
