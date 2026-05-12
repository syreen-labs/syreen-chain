package types

const (
	ModuleName = "mevprotection"
	StoreKey   = ModuleName
	RouterKey  = ModuleName

	// CommittedTxPrefixKey is the prefix for storing committed (encrypted) transactions
	CommittedTxPrefixKey = "committed"

	// RevealedTxPrefixKey is the prefix for storing revealed transactions
	RevealedTxPrefixKey = "revealed"

	// PenaltyPrefixKey is the prefix for storing MEV penalty records
	PenaltyPrefixKey = "penalty"
)

// CommittedTxKey returns the store key for a committed transaction by its hash
func CommittedTxKey(txHash []byte) []byte {
	return append([]byte(CommittedTxPrefixKey), txHash...)
}

// RevealedTxKey returns the store key for a revealed transaction by its commit hash
func RevealedTxKey(commitHash []byte) []byte {
	return append([]byte(RevealedTxPrefixKey), commitHash...)
}

// PenaltyKey returns the store key for a penalty record by validator address and block height
func PenaltyKey(validatorAddr string, blockHeight int64) []byte {
	key := append([]byte(PenaltyPrefixKey), []byte(validatorAddr)...)
	// Append block height as 8-byte big-endian
	h := uint64(blockHeight)
	heightBytes := []byte{
		byte(h >> 56), byte(h >> 48), byte(h >> 40), byte(h >> 32),
		byte(h >> 24), byte(h >> 16), byte(h >> 8), byte(h),
	}
	return append(key, heightBytes...)
}
