package types

const (
	ModuleName = "feemarket"
	StoreKey   = ModuleName
	RouterKey  = ModuleName

	// BaseFeeKey is the store key for the current base fee state
	BaseFeeKey = "basefee"

	// LaneConfigPrefix is the prefix for storing fee lane configurations
	LaneConfigPrefix = "lane/"

	// BurnTotalKey is the store key for cumulative burn totals
	BurnTotalKey = "burntotal"

	// BurnRecordPrefix is the prefix for per-block burn records
	BurnRecordPrefix = "burnrecord/"

	// ParamsKey is the store key for module parameters
	ParamsKey = "params"

	// PrevBlockGasKey stores the gas used in the previous block for fee adjustment
	PrevBlockGasKey = "prevblockgas"
)

// LaneConfigKey returns the store key for a specific fee lane
func LaneConfigKey(laneName string) []byte {
	return append([]byte(LaneConfigPrefix), []byte(laneName)...)
}

// BurnRecordKey returns the store key for a burn record at a specific height
func BurnRecordKey(height int64) []byte {
	return append([]byte(BurnRecordPrefix), Int64ToBytes(height)...)
}

// Int64ToBytes converts an int64 to a big-endian byte slice
func Int64ToBytes(i int64) []byte {
	b := make([]byte, 8)
	b[0] = byte(i >> 56)
	b[1] = byte(i >> 48)
	b[2] = byte(i >> 40)
	b[3] = byte(i >> 32)
	b[4] = byte(i >> 24)
	b[5] = byte(i >> 16)
	b[6] = byte(i >> 8)
	b[7] = byte(i)
	return b
}
