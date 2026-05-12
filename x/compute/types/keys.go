package types

const (
	ModuleName = "compute"
	StoreKey   = ModuleName
	RouterKey  = ModuleName

	// CodeKeyPrefix is the prefix for storing code bytecode and metadata
	CodeKeyPrefix = "code/"

	// ContractKeyPrefix is the prefix for storing contract info
	ContractKeyPrefix = "contract/"

	// ContractStatePrefix is the prefix for storing contract state
	ContractStatePrefix = "state/"

	// NextCodeIDKey is the key for the auto-incrementing code ID
	NextCodeIDKey = "next_code_id"
)

// CodeInfoKey returns the store key for a code's metadata
func CodeInfoKey(codeID uint64) []byte {
	return append([]byte(CodeKeyPrefix+"info/"), Uint64ToBytes(codeID)...)
}

// CodeBytecodeKey returns the store key for a code's bytecode
func CodeBytecodeKey(codeID uint64) []byte {
	return append([]byte(CodeKeyPrefix+"bytecode/"), Uint64ToBytes(codeID)...)
}

// ContractInfoKey returns the store key for a contract's info
func ContractInfoKey(contractAddr string) []byte {
	return append([]byte(ContractKeyPrefix), []byte(contractAddr)...)
}

// ContractStateKey returns the store key for a contract state entry
func ContractStateKey(contractAddr string, key []byte) []byte {
	prefix := append([]byte(ContractStatePrefix), []byte(contractAddr+"/")...)
	return append(prefix, key...)
}

// ContractStateIteratorPrefix returns the prefix for iterating all state of a contract
func ContractStateIteratorPrefix(contractAddr string) []byte {
	return []byte(ContractStatePrefix + contractAddr + "/")
}

func Uint64ToBytes(v uint64) []byte {
	b := make([]byte, 8)
	b[0] = byte(v >> 56)
	b[1] = byte(v >> 48)
	b[2] = byte(v >> 40)
	b[3] = byte(v >> 32)
	b[4] = byte(v >> 24)
	b[5] = byte(v >> 16)
	b[6] = byte(v >> 8)
	b[7] = byte(v)
	return b
}

func BytesToUint64(b []byte) uint64 {
	if len(b) < 8 {
		return 0
	}
	return uint64(b[0])<<56 | uint64(b[1])<<48 | uint64(b[2])<<40 | uint64(b[3])<<32 |
		uint64(b[4])<<24 | uint64(b[5])<<16 | uint64(b[6])<<8 | uint64(b[7])
}
