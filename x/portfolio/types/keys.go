package types

import "encoding/binary"

const (
	ModuleName = "portfolio"
	StoreKey   = ModuleName
	RouterKey  = ModuleName

	PortfolioPrefix       = "port/"       // port/{address} → Portfolio
	PortfolioAssetPrefix  = "passet/"     // passet/{address}/{denom} → PortfolioAsset
	TradeRecordPrefix     = "trade/"      // trade/{id} → TradeRecord
	TradeByAddrPrefix     = "trade_addr/" // trade_addr/{address}/{id} → exists
	ActivityPrefix        = "activity/"   // activity/{id} → Activity
	ActivityByAddrPrefix  = "act_addr/"   // act_addr/{address}/{id} → exists
	CompetitionPrefix     = "comp/"       // comp/{id} → Competition
	CompEntryPrefix       = "comp_entry/" // comp_entry/{compID}/{address} → CompetitionEntry
	NextTradeIDKey        = "next_trade_id"
	NextActivityIDKey     = "next_activity_id"
	NextCompetitionIDKey  = "next_comp_id"
	GlobalActivityListKey = "global_activities"
	SnapshotBlockInterval = 100 // take portfolio snapshots every N blocks

	MaxGlobalActivities  = 500
	MaxPerUserActivities = 100
	WhaleThreshold       = 10_000_000_000 // 10,000 SYR in usyreen (6 decimals)
)

func PortfolioKey(address string) []byte {
	return append([]byte(PortfolioPrefix), []byte(address)...)
}

func PortfolioAssetKey(address, denom string) []byte {
	return append([]byte(PortfolioAssetPrefix), []byte(address+"/"+denom)...)
}

func PortfolioAssetPrefixKey(address string) []byte {
	return []byte(PortfolioAssetPrefix + address + "/")
}

func TradeRecordKey(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return append([]byte(TradeRecordPrefix), bz...)
}

func TradeByAddrKey(address string, id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return append([]byte(TradeByAddrPrefix+address+"/"), bz...)
}

func TradeByAddrPrefixKey(address string) []byte {
	return []byte(TradeByAddrPrefix + address + "/")
}

func ActivityKey(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return append([]byte(ActivityPrefix), bz...)
}

func ActivityByAddrKey(address string, id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return append([]byte(ActivityByAddrPrefix+address+"/"), bz...)
}

func ActivityByAddrPrefixKey(address string) []byte {
	return []byte(ActivityByAddrPrefix + address + "/")
}

func CompetitionKey(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return append([]byte(CompetitionPrefix), bz...)
}

func CompEntryKey(compID uint64, address string) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, compID)
	return append(append([]byte(CompEntryPrefix), bz...), []byte("/"+address)...)
}

func CompEntryPrefixKey(compID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, compID)
	return append([]byte(CompEntryPrefix), bz...)
}
