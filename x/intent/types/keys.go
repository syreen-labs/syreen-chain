package types

import "encoding/binary"

const (
	ModuleName = "intent"
	StoreKey   = ModuleName
	RouterKey  = ModuleName

	// IntentPrefix is the prefix for storing intents
	IntentPrefix = "intent/"

	// SolverPrefix is the prefix for storing solvers
	SolverPrefix    = "solver/"
	SolverPrefixKey = SolverPrefix

	// SolutionPrefix is the prefix for storing solutions
	SolutionPrefix = "solution/"

	// ResultPrefix is the prefix for storing intent results
	ResultPrefix = "result/"

	// ChainPrefix is the prefix for storing intent chains
	ChainPrefix = "chain/"

	// IntentChainMapPrefix maps intent IDs to their parent chain
	IntentChainMapPrefix = "intent_chain_map/"

	// ParamsKey is the key for storing module params
	ParamsKey = "params"
)

// IntentKey returns the store key for an intent by ID
func IntentKey(intentID string) []byte {
	return append([]byte(IntentPrefix), []byte(intentID)...)
}

// SolverKey returns the store key for a solver by address
func SolverKey(address string) []byte {
	return append([]byte(SolverPrefix), []byte(address)...)
}

// SolutionKey returns the store key for a solution by intent ID and solver address
func SolutionKey(intentID string, solverAddr string) []byte {
	key := append([]byte(SolutionPrefix), []byte(intentID)...)
	key = append(key, '/')
	return append(key, []byte(solverAddr)...)
}

// SolutionsByIntentPrefix returns the prefix for all solutions of a given intent
func SolutionsByIntentPrefix(intentID string) []byte {
	key := append([]byte(SolutionPrefix), []byte(intentID)...)
	return append(key, '/')
}

// ResultKey returns the store key for a result by intent ID
func ResultKey(intentID string) []byte {
	return append([]byte(ResultPrefix), []byte(intentID)...)
}

// ChainKey returns the store key for a chain by ID
func ChainKey(chainID string) []byte {
	return append([]byte(ChainPrefix), []byte(chainID)...)
}

// IntentChainMapKey returns the store key mapping an intent ID to its parent chain
func IntentChainMapKey(intentID string) []byte {
	return append([]byte(IntentChainMapPrefix), []byte(intentID)...)
}

// ChainCounterKey is the key for the chain ID counter
func ChainCounterKey() []byte {
	return []byte("chain_counter")
}

// IntentCounterKey is the key for the intent ID counter
func IntentCounterKey() []byte {
	return []byte("intent_counter")
}

// Uint64ToBytes converts a uint64 to big-endian bytes
func Uint64ToBytes(v uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, v)
	return bz
}

// BytesToUint64 converts big-endian bytes to uint64
func BytesToUint64(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
