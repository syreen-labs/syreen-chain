package types

const (
	ModuleName = "evm"
	StoreKey   = ModuleName
	TStoreKey  = "transient_" + ModuleName
	RouterKey  = ModuleName

	// Prefixes for KV store
	PrefixCode     = 0x01 // contract bytecode: prefix + address -> code
	PrefixStorage  = 0x02 // contract storage: prefix + address + slot -> value
	PrefixNonce    = 0x03 // EVM nonces: prefix + address -> nonce
	PrefixCodeHash = 0x04 // code hashes: prefix + address -> hash
	PrefixAccount  = 0x05 // EVM accounts (tracks existence): prefix + address -> 0x01
	PrefixLogs     = 0x06 // tx logs: prefix + txHash -> logs
	PrefixBloom    = 0x07 // block bloom: prefix + height -> bloom
	PrefixParams   = 0x08 // module params: prefix -> JSON params
)

// KeyCode returns the store key for contract code
func KeyCode(addr []byte) []byte {
	return append([]byte{PrefixCode}, addr...)
}

// KeyStorage returns the store key for a storage slot.
// addr must be exactly 20 bytes (Ethereum address), slot is 0 or 32 bytes.
func KeyStorage(addr []byte, slot []byte) []byte {
	if len(addr) != 0 && len(addr) != 20 {
		panic("KeyStorage: addr must be 20 bytes")
	}
	key := append([]byte{PrefixStorage}, addr...)
	return append(key, slot...)
}

// KeyNonce returns the store key for an EVM nonce
func KeyNonce(addr []byte) []byte {
	return append([]byte{PrefixNonce}, addr...)
}

// KeyCodeHash returns the store key for a code hash
func KeyCodeHash(addr []byte) []byte {
	return append([]byte{PrefixCodeHash}, addr...)
}

// KeyAccount returns the store key for EVM account existence
func KeyAccount(addr []byte) []byte {
	return append([]byte{PrefixAccount}, addr...)
}

// KeyParams returns the store key for EVM module params (M6)
func KeyParams() []byte {
	return []byte{PrefixParams}
}
