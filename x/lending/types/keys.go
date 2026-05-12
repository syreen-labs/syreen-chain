package types

import "encoding/binary"

const (
	ModuleName = "lending"
	StoreKey   = ModuleName
	RouterKey  = ModuleName

	LendingPoolPrefix   = "lpool/"      // lpool/{poolID} → LendingPool
	DepositPrefix       = "deposit/"    // deposit/{poolID}/{address} → Deposit
	DepositByAddrPrefix = "dep_addr/"   // dep_addr/{address}/{poolID} → exists
	BorrowPrefix        = "borrow/"     // borrow/{borrowID} → Borrow
	BorrowByAddrPrefix  = "brw_addr/"   // brw_addr/{address}/{borrowID} → exists
	NextPoolIDKey       = "next_lpool_id"
	NextBorrowIDKey     = "next_borrow_id"
	GlobalStatsKey      = "lending_stats"
	TWAPSeriesPrefix    = "twap/"       // twap/{poolID} → TWAPSeries
)

// TWAPSeriesKey returns the storage key for the TWAP samples of a lending pool.
func TWAPSeriesKey(poolID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	return append([]byte(TWAPSeriesPrefix), bz...)
}

func LendingPoolKey(poolID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	return append([]byte(LendingPoolPrefix), bz...)
}

func DepositKey(poolID uint64, address string) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	return append(append([]byte(DepositPrefix), bz...), []byte("/"+address)...)
}

func DepositByAddrKey(address string, poolID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	return append(append([]byte(DepositByAddrPrefix), []byte(address+"/")...), bz...)
}

func DepositByAddrPrefixKey(address string) []byte {
	return []byte(DepositByAddrPrefix + address + "/")
}

func DepositPoolPrefix(poolID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	return append([]byte(DepositPrefix), bz...)
}

func BorrowKey(borrowID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, borrowID)
	return append([]byte(BorrowPrefix), bz...)
}

func BorrowByAddrKey(address string, borrowID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, borrowID)
	return append(append([]byte(BorrowByAddrPrefix), []byte(address+"/")...), bz...)
}

func BorrowByAddrPrefixKey(address string) []byte {
	return []byte(BorrowByAddrPrefix + address + "/")
}
