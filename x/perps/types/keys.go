package types

import "encoding/binary"

const (
	ModuleName = "perps"
	StoreKey   = ModuleName
	RouterKey  = ModuleName

	// KV prefixes
	MarketPrefix        = "market/"       // market/{marketID} → PerpMarket
	PositionPrefix      = "position/"     // position/{marketID}/{address} → Position
	PositionByAddr      = "pos_addr/"     // pos_addr/{address}/{marketID} → exists
	InsuranceFundKey    = "insurance"     // insurance → InsuranceFund
	FundingRatePrefix   = "funding/"      // funding/{marketID} → FundingState
	NextMarketIDKey     = "next_market_id"
	LiquidationPrefix   = "liq/"         // liq/{marketID}/{address} → pending liquidation
	MarkPriceTWAPPrefix = "twap/"        // twap/{marketID} → []LegacyDec (last 10 prices)

	// Oracle-free TWAP price samples: perps_twap/<poolID> → []PriceSample
	PerpsTWAPPrefix = "perps_twap/"
)

func MarketKey(marketID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, marketID)
	return append([]byte(MarketPrefix), bz...)
}

func PositionKey(marketID uint64, address string) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, marketID)
	return append(append([]byte(PositionPrefix), bz...), []byte("/"+address)...)
}

func PositionByAddrKey(address string, marketID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, marketID)
	return append(append([]byte(PositionByAddr), []byte(address+"/")...), bz...)
}

func PositionByAddrPrefix(address string) []byte {
	return []byte(PositionByAddr + address + "/")
}

func PositionMarketPrefix(marketID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, marketID)
	return append([]byte(PositionPrefix), bz...)
}

func FundingRateKey(marketID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, marketID)
	return append([]byte(FundingRatePrefix), bz...)
}

func MarkPriceTWAPKey(marketID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, marketID)
	return append([]byte(MarkPriceTWAPPrefix), bz...)
}

// PerpsTWAPKey returns the store key for oracle-free TWAP price samples.
// Format: perps_twap/<poolID_8bytes>
func PerpsTWAPKey(poolID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	return append([]byte(PerpsTWAPPrefix), bz...)
}
