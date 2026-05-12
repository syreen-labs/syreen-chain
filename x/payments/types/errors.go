package types

import "cosmossdk.io/errors"

var (
	ErrInvoiceNotFound      = errors.Register(ModuleName, 2, "invoice not found")
	ErrInvoiceAlreadyPaid   = errors.Register(ModuleName, 3, "invoice already paid")
	ErrInvoiceExpired       = errors.Register(ModuleName, 4, "invoice has expired")
	ErrInsufficientFunds    = errors.Register(ModuleName, 5, "insufficient funds")
	ErrInvalidAmount        = errors.Register(ModuleName, 6, "invalid amount")
	ErrNotInvoicePayer      = errors.Register(ModuleName, 7, "not the invoice payer")
	ErrNotInvoicePayee      = errors.Register(ModuleName, 8, "not the invoice payee")
	ErrExchangeRateNotFound = errors.Register(ModuleName, 9, "exchange rate not found")
	ErrNoEarningsAvailable  = errors.Register(ModuleName, 10, "no earnings available")
	ErrInvoiceCancelled     = errors.Register(ModuleName, 11, "invoice has been cancelled")
	ErrPaymentsDisabled     = errors.Register(ModuleName, 12, "payments module is disabled")
	ErrNotInvoiceCreator    = errors.Register(ModuleName, 13, "not the invoice creator")
	ErrInvalidInvoiceStatus = errors.Register(ModuleName, 14, "invalid invoice status for this operation")
	ErrUnauthorized         = errors.Register(ModuleName, 15, "unauthorized")
)
