package ante

// GaslessTxKey is a context key set by the GaslessDexDecorator
// to signal that this transaction is gasless and should bypass
// the fee market minimum fee check.
type GaslessTxKey struct{}
