package types

import "encoding/binary"

const (
	ModuleName = "predict"
	StoreKey   = ModuleName
	RouterKey  = ModuleName

	MarketPrefix      = "market/"      // market/{marketID} -> Market
	PositionPrefix    = "position/"    // position/{marketID}/{address} -> Position
	PosByAddrPrefix   = "pos_addr/"    // pos_addr/{address}/{marketID} -> exists
	ResolutionPrefix  = "resolution/"  // resolution/{marketID} -> Resolution
	NextMarketIDKey   = "next_market_id"
)

func MarketKey(marketID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, marketID)
	return append([]byte(MarketPrefix), bz...)
}

func MarketPrefixBytes() []byte {
	return []byte(MarketPrefix)
}

func PositionKey(marketID uint64, address string) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, marketID)
	return append(append([]byte(PositionPrefix), bz...), []byte("/"+address)...)
}

func PositionMarketPrefix(marketID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, marketID)
	return append([]byte(PositionPrefix), bz...)
}

func PosByAddrKey(address string, marketID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, marketID)
	return append(append([]byte(PosByAddrPrefix), []byte(address+"/")...), bz...)
}

func ResolutionKey(marketID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, marketID)
	return append([]byte(ResolutionPrefix), bz...)
}
