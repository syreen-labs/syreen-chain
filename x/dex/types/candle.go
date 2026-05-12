package types

import (
	"encoding/binary"

	"cosmossdk.io/math"
)

// ---------------------------------------------------------------------------
// OHLCV Candle Types
// ---------------------------------------------------------------------------

// Candle stores OHLCV data for a pool at a specific interval.
type Candle struct {
	PoolID     uint64         `json:"pool_id"`
	Interval   string         `json:"interval"`    // e.g. "1m", "5m", "15m", "1h", "4h", "1d"
	OpenBlock  int64          `json:"open_block"`   // block height when this candle started
	Open       math.LegacyDec `json:"open"`
	High       math.LegacyDec `json:"high"`
	Low        math.LegacyDec `json:"low"`
	Close      math.LegacyDec `json:"close"`
	Volume     math.Int       `json:"volume"`       // total base-denom volume
	TradeCount uint64         `json:"trade_count"`
	Timestamp  int64          `json:"timestamp"`    // block height of last update
}

// ---------------------------------------------------------------------------
// Interval definitions (block-based, ~500ms per block)
// ---------------------------------------------------------------------------

const (
	CandleInterval1m  = "1m"   // ~120 blocks (~1 minute)
	CandleInterval5m  = "5m"   // ~600 blocks (~5 minutes)
	CandleInterval15m = "15m"  // ~1800 blocks (~15 minutes)
	CandleInterval1h  = "1h"   // ~7200 blocks (~1 hour)
	CandleInterval4h  = "4h"   // ~28800 blocks (~4 hours)
	CandleInterval1d  = "1d"   // ~172800 blocks (~1 day)
)

// CandleIntervalBlocks maps interval names to their block count.
var CandleIntervalBlocks = map[string]int64{
	CandleInterval1m:  120,
	CandleInterval5m:  600,
	CandleInterval15m: 1800,
	CandleInterval1h:  7200,
	CandleInterval4h:  28800,
	CandleInterval1d:  172800,
}

// AllCandleIntervals is the ordered list of intervals for deterministic iteration.
var AllCandleIntervals = []string{
	CandleInterval1m,
	CandleInterval5m,
	CandleInterval15m,
	CandleInterval1h,
	CandleInterval4h,
	CandleInterval1d,
}

// ---------------------------------------------------------------------------
// KV Store keys
// ---------------------------------------------------------------------------

const (
	CandlePrefix        = "candle/"         // candle/<poolID>/<interval>/<openBlock>
	CandleCurrentPrefix = "candle_current/" // candle_current/<poolID>/<interval> -> current open candle
	MaxCandlesPerInterval = 1000
)

// CandleKey returns the store key for a finalized candle.
// Format: candle/<poolID_8bytes><interval_padded_4bytes><openBlock_8bytes>
func CandleKey(poolID uint64, interval string, openBlock int64) []byte {
	poolBz := make([]byte, 8)
	binary.BigEndian.PutUint64(poolBz, poolID)

	intervalBz := encodeInterval(interval)

	blockBz := make([]byte, 8)
	binary.BigEndian.PutUint64(blockBz, uint64(openBlock))

	key := append([]byte(CandlePrefix), poolBz...)
	key = append(key, intervalBz...)
	key = append(key, blockBz...)
	return key
}

// CandlePoolIntervalPrefix returns the prefix for all candles of a pool+interval.
func CandlePoolIntervalPrefix(poolID uint64, interval string) []byte {
	poolBz := make([]byte, 8)
	binary.BigEndian.PutUint64(poolBz, poolID)

	intervalBz := encodeInterval(interval)

	key := append([]byte(CandlePrefix), poolBz...)
	key = append(key, intervalBz...)
	return key
}

// CandleCurrentKey returns the store key for the current (in-progress) candle.
// Format: candle_current/<poolID_8bytes><interval_padded_4bytes>
func CandleCurrentKey(poolID uint64, interval string) []byte {
	poolBz := make([]byte, 8)
	binary.BigEndian.PutUint64(poolBz, poolID)

	intervalBz := encodeInterval(interval)

	key := append([]byte(CandleCurrentPrefix), poolBz...)
	key = append(key, intervalBz...)
	return key
}

// encodeInterval encodes an interval string as a fixed 4-byte sortable value.
func encodeInterval(interval string) []byte {
	bz := make([]byte, 4)
	// Use block count for deterministic sorting
	blocks, ok := CandleIntervalBlocks[interval]
	if !ok {
		blocks = 0
	}
	binary.BigEndian.PutUint32(bz, uint32(blocks))
	return bz
}
