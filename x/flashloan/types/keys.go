package types

import "encoding/binary"

const (
	ModuleName = "flashloan"
	StoreKey   = ModuleName
	RouterKey  = ModuleName

	FlashPoolPrefix        = "fpool/"        // fpool/{denom} -> FlashLoanPool
	FlashRecordPrefix      = "frecord/"      // frecord/{id} -> FlashLoanRecord
	FlashStatsKey          = "flash_stats"   // -> FlashLoanStats
	NextRecordIDKey        = "next_frecord_id"
	FlashPoolDepositPrefix = "fpdep/"        // fpdep/{denom}/{depositor} -> FlashPoolDeposit (H-1)
)

func FlashPoolKey(denom string) []byte {
	return append([]byte(FlashPoolPrefix), []byte(denom)...)
}

func FlashRecordKey(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return append([]byte(FlashRecordPrefix), bz...)
}

// FlashPoolDepositKey returns the store key for a depositor's position in a flash pool.
func FlashPoolDepositKey(denom, depositor string) []byte {
	return append([]byte(FlashPoolDepositPrefix+denom+"/"), []byte(depositor)...)
}

// FlashPoolDepositPrefixKey returns the prefix for iterating all depositors in a pool.
func FlashPoolDepositPrefixKey(denom string) []byte {
	return []byte(FlashPoolDepositPrefix + denom + "/")
}
