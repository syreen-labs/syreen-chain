package types

import "cosmossdk.io/errors"

var (
	ErrIntentNotFound         = errors.Register(ModuleName, 2, "intent not found")
	ErrIntentExpired          = errors.Register(ModuleName, 3, "intent has expired")
	ErrIntentAlreadyFulfilled = errors.Register(ModuleName, 4, "intent already fulfilled")
	ErrSolverNotRegistered    = errors.Register(ModuleName, 5, "solver is not registered")
	ErrSolverAlreadyRegistered = errors.Register(ModuleName, 6, "solver is already registered")
	ErrInsufficientStake      = errors.Register(ModuleName, 7, "insufficient solver stake")
	ErrInvalidSolution        = errors.Register(ModuleName, 8, "invalid solution")
	ErrSolvingWindowClosed    = errors.Register(ModuleName, 9, "solving window has closed")
	ErrIntentTypeUnsupported  = errors.Register(ModuleName, 10, "unsupported intent type")
	ErrMaxSolutionsReached    = errors.Register(ModuleName, 11, "maximum solutions reached for intent")
	ErrNoSolutions            = errors.Register(ModuleName, 12, "no solutions submitted for intent")
	ErrReentrant              = errors.Register(ModuleName, 13, "reentrant call to FulfillIntent is not allowed")
	ErrDisallowedMsgType      = errors.Register(ModuleName, 14, "message type is not allowed in solution execution")
	ErrChainNotFound          = errors.Register(ModuleName, 15, "intent chain not found")
	ErrChainAlreadyComplete   = errors.Register(ModuleName, 16, "intent chain already completed")
	ErrChainTooFewSteps       = errors.Register(ModuleName, 17, "intent chain requires at least 2 steps")
	ErrChainTooManySteps      = errors.Register(ModuleName, 18, "intent chain exceeds maximum steps")
	ErrChainInvalidCondition  = errors.Register(ModuleName, 19, "invalid chain step condition")
	ErrChainNotCreator        = errors.Register(ModuleName, 20, "only chain creator can cancel")
	ErrCrossChainTransferFail = errors.Register(ModuleName, 21, "IBC transfer failed for cross-chain intent")
	ErrIntentNotCreator       = errors.Register(ModuleName, 22, "only the intent creator can cancel it")
	ErrIntentNotCancellable   = errors.Register(ModuleName, 23, "intent cannot be cancelled in its current state")
)
