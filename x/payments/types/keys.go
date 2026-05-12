package types

const (
	ModuleName   = "payments"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	// Key prefixes for store
	InvoiceKeyPrefix      = "invoice/"
	PaymentKeyPrefix      = "payment/"
	ExchangeRateKeyPrefix = "exchange_rate/"
	EarningsKeyPrefix     = "earnings/"
	InvoiceCounterKey     = "invoice_counter"
	ParamsKey             = "params"

	// Index prefixes
	InvoiceByBookingPrefix = "invoice_booking/"
	InvoiceByPayerPrefix   = "invoice_payer/"
)

func InvoiceStoreKey(id string) []byte {
	return []byte(InvoiceKeyPrefix + id)
}

func ExchangeRateStoreKey(fromDenom, toDenom string) []byte {
	return []byte(ExchangeRateKeyPrefix + fromDenom + "/" + toDenom)
}

func EarningsStoreKey(address string) []byte {
	return []byte(EarningsKeyPrefix + address)
}

func InvoiceByBookingKey(bookingID string) []byte {
	return []byte(InvoiceByBookingPrefix + bookingID)
}

func InvoiceByPayerKey(payer string) []byte {
	return []byte(InvoiceByPayerPrefix + payer)
}
