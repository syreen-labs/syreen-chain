package statedb

import "errors"

var (
	// ErrNegativeAmount is returned when a negative amount is passed to AddBalance/SubBalance
	ErrNegativeAmount = errors.New("evm statedb: negative amount")

	// ErrCodeTooLarge is returned when contract code exceeds EIP-170 limit (24KB)
	ErrCodeTooLarge = errors.New("evm statedb: code size exceeds 24KB limit (EIP-170)")

	// ErrRefundUnderflow is returned when SubRefund tries to subtract more than the current refund
	ErrRefundUnderflow = errors.New("evm statedb: refund counter underflow")
)
