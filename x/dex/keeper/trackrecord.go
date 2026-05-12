package keeper

import (
	"context"
	"encoding/json"
	"sort"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ---------------------------------------------------------------------------
// On-Chain Track Records for Copy Trading
//
// Tracks verifiable on-chain performance for every trader, including:
// - Win/loss record, PnL, volume, best/worst trade
// - Copy trader follower count and aggregate copy PnL
// - Leaderboard sorted by multiple criteria
// ---------------------------------------------------------------------------

const (
	// KV store prefix
	traderRecordPrefix = "trader_record/" // address -> TraderRecord JSON
)

// TraderRecord stores comprehensive per-trader performance stats.
type TraderRecord struct {
	Address         string         `json:"address"`
	TotalTrades     uint64         `json:"total_trades"`
	WinningTrades   uint64         `json:"winning_trades"`
	LosingTrades    uint64         `json:"losing_trades"`
	TotalVolume     math.Int       `json:"total_volume"`
	TotalPnL        math.LegacyDec `json:"total_pnl"`
	BestTrade       math.LegacyDec `json:"best_trade"`
	WorstTrade      math.LegacyDec `json:"worst_trade"`
	WinRate         math.LegacyDec `json:"win_rate"`
	AvgTradeSize    math.Int       `json:"avg_trade_size"`
	ActiveSince     int64          `json:"active_since"`
	LastTradeHeight int64          `json:"last_trade_height"`
	Followers       uint64         `json:"followers"`
	CopyPnL         math.LegacyDec `json:"copy_pnl"`
}

// newEmptyTraderRecord returns a zero-valued TraderRecord for the given address.
func newEmptyTraderRecord(address string) TraderRecord {
	return TraderRecord{
		Address:      address,
		TotalVolume:  math.ZeroInt(),
		TotalPnL:     math.LegacyZeroDec(),
		BestTrade:    math.LegacyZeroDec(),
		WorstTrade:   math.LegacyZeroDec(),
		WinRate:      math.LegacyZeroDec(),
		AvgTradeSize: math.ZeroInt(),
		CopyPnL:      math.LegacyZeroDec(),
	}
}

// ---------------------------------------------------------------------------
// Update / Get
// ---------------------------------------------------------------------------

// UpdateTraderRecord updates the on-chain track record after each swap.
// pnlPercent is the PnL percentage for this trade (positive = profit, negative = loss).
// tokenInAmount is the trade input amount (used for volume tracking).
func (k Keeper) UpdateTraderRecord(ctx context.Context, trader string, poolID uint64, tokenInAmount, tokenOutAmount math.Int, pnlPercent math.LegacyDec) {
	// Anti-dust: skip track record updates for trades below minimum threshold.
	// This prevents leaderboard manipulation via thousands of tiny profitable trades.
	minTradeSize := math.NewInt(1_000_000) // 1 token minimum (6 decimals)
	if tokenInAmount.LT(minTradeSize) {
		return
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	height := sdkCtx.BlockHeight()

	record := k.GetTraderRecord(ctx, trader)

	// Update trade count
	record.TotalTrades++

	// Update volume
	record.TotalVolume = record.TotalVolume.Add(tokenInAmount)

	// Update PnL
	record.TotalPnL = record.TotalPnL.Add(pnlPercent)

	// Win/loss tracking
	if pnlPercent.IsPositive() {
		record.WinningTrades++
	} else if pnlPercent.IsNegative() {
		record.LosingTrades++
	}

	// Best/worst trade
	if pnlPercent.GT(record.BestTrade) {
		record.BestTrade = pnlPercent
	}
	if pnlPercent.LT(record.WorstTrade) {
		record.WorstTrade = pnlPercent
	}

	// Win rate
	if record.TotalTrades > 0 {
		record.WinRate = math.LegacyNewDec(int64(record.WinningTrades)).Quo(math.LegacyNewDec(int64(record.TotalTrades)))
	}

	// Average trade size
	if record.TotalTrades > 0 {
		record.AvgTradeSize = record.TotalVolume.Quo(math.NewInt(int64(record.TotalTrades)))
	}

	// Active since (first trade)
	if record.ActiveSince == 0 {
		record.ActiveSince = height
	}

	// Last trade
	record.LastTradeHeight = height

	// Follower count from copy trading system
	followers := k.GetFollowers(ctx, trader)
	record.Followers = uint64(len(followers))

	k.SetTraderRecord(ctx, record)
}

// GetTraderRecord returns the track record for a trader.
func (k Keeper) GetTraderRecord(ctx context.Context, address string) TraderRecord {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(traderRecordKey(address))
	if err != nil || bz == nil {
		return newEmptyTraderRecord(address)
	}
	var record TraderRecord
	if err := json.Unmarshal(bz, &record); err != nil {
		return newEmptyTraderRecord(address)
	}
	return record
}

// SetTraderRecord stores a trader's track record.
func (k Keeper) SetTraderRecord(ctx context.Context, record TraderRecord) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(record)
	kvStore.Set(traderRecordKey(record.Address), bz)
}

// ---------------------------------------------------------------------------
// Leaderboards
// ---------------------------------------------------------------------------

// GetTopTraders returns the top traders sorted by cumulative PnL.
func (k Keeper) GetTopTraders(ctx context.Context, limit int) []TraderRecord {
	return k.GetTraderLeaderboard(ctx, "pnl", limit)
}

// GetTraderLeaderboard returns traders sorted by the specified criterion.
// sortBy can be: "pnl", "volume", "win_rate", "trades", "followers"
func (k Keeper) GetTraderLeaderboard(ctx context.Context, sortBy string, limit int) []TraderRecord {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(traderRecordPrefix)

	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var allRecords []TraderRecord
	for ; iter.Valid(); iter.Next() {
		var record TraderRecord
		if err := json.Unmarshal(iter.Value(), &record); err != nil {
			continue
		}
		if record.TotalTrades > 0 {
			allRecords = append(allRecords, record)
		}
	}

	// Sort based on criterion with deterministic tiebreaker by address
	switch sortBy {
	case "volume":
		sort.SliceStable(allRecords, func(i, j int) bool {
			if allRecords[i].TotalVolume.Equal(allRecords[j].TotalVolume) {
				return allRecords[i].Address < allRecords[j].Address
			}
			return allRecords[i].TotalVolume.GT(allRecords[j].TotalVolume)
		})
	case "win_rate":
		sort.SliceStable(allRecords, func(i, j int) bool {
			if allRecords[i].WinRate.Equal(allRecords[j].WinRate) {
				return allRecords[i].Address < allRecords[j].Address
			}
			return allRecords[i].WinRate.GT(allRecords[j].WinRate)
		})
	case "trades":
		sort.SliceStable(allRecords, func(i, j int) bool {
			if allRecords[i].TotalTrades == allRecords[j].TotalTrades {
				return allRecords[i].Address < allRecords[j].Address
			}
			return allRecords[i].TotalTrades > allRecords[j].TotalTrades
		})
	case "followers":
		sort.SliceStable(allRecords, func(i, j int) bool {
			if allRecords[i].Followers == allRecords[j].Followers {
				return allRecords[i].Address < allRecords[j].Address
			}
			return allRecords[i].Followers > allRecords[j].Followers
		})
	default: // "pnl"
		sort.SliceStable(allRecords, func(i, j int) bool {
			if allRecords[i].TotalPnL.Equal(allRecords[j].TotalPnL) {
				return allRecords[i].Address < allRecords[j].Address
			}
			return allRecords[i].TotalPnL.GT(allRecords[j].TotalPnL)
		})
	}

	if limit > 0 && len(allRecords) > limit {
		allRecords = allRecords[:limit]
	}

	return allRecords
}

// UpdateCopyPnL updates the aggregate PnL of copy traders following a given trader.
func (k Keeper) UpdateCopyPnL(ctx context.Context, trader string, copyPnlDelta math.LegacyDec) {
	record := k.GetTraderRecord(ctx, trader)
	record.CopyPnL = record.CopyPnL.Add(copyPnlDelta)
	k.SetTraderRecord(ctx, record)
}

// ---------------------------------------------------------------------------
// KV Store Key
// ---------------------------------------------------------------------------

func traderRecordKey(address string) []byte {
	return []byte(traderRecordPrefix + address)
}
