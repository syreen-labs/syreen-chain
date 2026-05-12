package types

import (
	"encoding/binary"
)

const (
	ModuleName = "dex"
	StoreKey   = ModuleName
	RouterKey  = ModuleName

	PoolPrefix        = "pool/"
	PoolByDenomPrefix = "pool_by_denom/"

	// Sentiment oracle store prefixes
	SentimentBlockDataPrefix = "sentiment_block/"  // per-pool per-block trading data
	SentimentStatePrefix     = "sentiment_state/"  // current sentiment state per pool
	SentimentAlertPrefix     = "sentiment_alert/"  // active sentiment alerts per pool

	// Referral system store prefixes
	ReferralCodePrefix       = "referral_code/"    // referral code -> referrer address
	ReferralByAddrPrefix     = "referral_byaddr/"  // address -> referral code
	ReferralRegistrationPrefix = "referral_reg/"   // referred user -> referrer address
	ReferralStatsPrefix      = "referral_stats/"   // referrer address -> stats
	ReferralUserStatsPrefix  = "referral_ustats/"  // referred user -> user stats
	ReferralGlobalPrefix     = "referral_global"   // global referral stats (single key)
	ReferralEarningsPrefix   = "referral_earn/"    // referrer address -> earnings entries
	ReferralClaimablePrefix  = "referral_claim/"   // <referrer>/<denom> -> claimable amount

	// Risk score round-robin cursor (H6)
	RiskScoreCursorKey = "risk_score_cursor"

	// Sentiment oracle primitive store prefixes
	OracleCompositePrefix = "oracle_composite/" // oracle_composite/<poolID> -> OracleComposite JSON
	WhaleAlertPrefix      = "whale_alert/"      // whale_alert/<poolID>/<height>/<index> -> WhaleAlert JSON
)

// SentimentBlockDataKey returns the store key for block-level trading data.
// Format: sentiment_block/<poolID_8bytes><height_8bytes>
func SentimentBlockDataKey(poolID uint64, height int64) []byte {
	key := make([]byte, 16)
	binary.BigEndian.PutUint64(key[:8], poolID)
	binary.BigEndian.PutUint64(key[8:], uint64(height))
	return append([]byte(SentimentBlockDataPrefix), key...)
}

// SentimentBlockDataPoolPrefix returns the prefix for all block data for a pool.
func SentimentBlockDataPoolPrefix(poolID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	return append([]byte(SentimentBlockDataPrefix), bz...)
}

// SentimentStateKey returns the store key for a pool's current sentiment state.
func SentimentStateKey(poolID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	return append([]byte(SentimentStatePrefix), bz...)
}

// SentimentAlertKey returns the store key for a pool's sentiment alerts.
func SentimentAlertKey(poolID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	return append([]byte(SentimentAlertPrefix), bz...)
}

// OracleCompositeKey returns the store key for a pool's oracle composite state.
func OracleCompositeKey(poolID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	return append([]byte(OracleCompositePrefix), bz...)
}

// WhaleAlertKey returns the store key for a whale alert.
// Format: whale_alert/<poolID_8bytes><height_8bytes><index_8bytes>
func WhaleAlertKey(poolID uint64, height int64, index uint64) []byte {
	key := make([]byte, 24)
	binary.BigEndian.PutUint64(key[:8], poolID)
	binary.BigEndian.PutUint64(key[8:16], uint64(height))
	binary.BigEndian.PutUint64(key[16:], index)
	return append([]byte(WhaleAlertPrefix), key...)
}

// WhaleAlertPoolPrefix returns the prefix for all whale alerts for a pool.
func WhaleAlertPoolPrefix(poolID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	return append([]byte(WhaleAlertPrefix), bz...)
}

// PoolKey returns the store key for a pool by its ID.
func PoolKey(poolID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	return append([]byte(PoolPrefix), bz...)
}

// PoolByDenomPairKey returns the index key for a denom pair (always sorted alphabetically).
func PoolByDenomPairKey(denomA, denomB string) []byte {
	a, b := SortDenoms(denomA, denomB)
	return append([]byte(PoolByDenomPrefix), []byte(a+"/"+b)...)
}
