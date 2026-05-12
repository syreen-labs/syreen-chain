package keeper

import (
	"context"
	"encoding/json"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/dex/types"
)

// ---------------------------------------------------------------------------
// UpdateCandles — called after each trade execution
// ---------------------------------------------------------------------------

// UpdateCandles updates the current candle for every interval when a trade occurs.
func (k Keeper) UpdateCandles(ctx context.Context, poolID uint64, tradePrice math.LegacyDec, tradeVolume math.Int) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	blockHeight := sdkCtx.BlockHeight()

	for _, interval := range types.AllCandleIntervals {
		k.updateCandleForInterval(ctx, poolID, interval, tradePrice, tradeVolume, blockHeight)
	}
}

func (k Keeper) updateCandleForInterval(ctx context.Context, poolID uint64, interval string, tradePrice math.LegacyDec, tradeVolume math.Int, blockHeight int64) {
	intervalBlocks := types.CandleIntervalBlocks[interval]

	// Determine which candle period this block belongs to
	candleOpenBlock := (blockHeight / intervalBlocks) * intervalBlocks

	// Get current candle
	candle, found := k.GetCurrentCandle(ctx, poolID, interval)

	if !found || candle.OpenBlock != candleOpenBlock {
		// If we had an old candle from a different period, finalize it
		if found && candle.OpenBlock != candleOpenBlock {
			k.finalizeCandle(ctx, candle)
		}

		// Start a new candle
		candle = types.Candle{
			PoolID:     poolID,
			Interval:   interval,
			OpenBlock:  candleOpenBlock,
			Open:       tradePrice,
			High:       tradePrice,
			Low:        tradePrice,
			Close:      tradePrice,
			Volume:     tradeVolume,
			TradeCount: 1,
			Timestamp:  blockHeight,
		}
	} else {
		// Update existing candle
		if tradePrice.GT(candle.High) {
			candle.High = tradePrice
		}
		if tradePrice.LT(candle.Low) {
			candle.Low = tradePrice
		}
		candle.Close = tradePrice
		candle.Volume = candle.Volume.Add(tradeVolume)
		candle.TradeCount++
		candle.Timestamp = blockHeight
	}

	k.SetCurrentCandle(ctx, candle)
}

// ---------------------------------------------------------------------------
// FinalizeCandles — called in EndBlock
// ---------------------------------------------------------------------------

// FinalizeCandles checks all pools and intervals for candle periods that have ended.
func (k Keeper) FinalizeCandles(ctx sdk.Context) {
	blockHeight := ctx.BlockHeight()
	pools := k.GetAllPools(ctx)

	for _, pool := range pools {
		for _, interval := range types.AllCandleIntervals {
			intervalBlocks := types.CandleIntervalBlocks[interval]
			currentPeriod := (blockHeight / intervalBlocks) * intervalBlocks

			candle, found := k.GetCurrentCandle(ctx, pool.ID, interval)
			if !found {
				continue
			}

			// If the current candle belongs to a previous period, finalize it
			if candle.OpenBlock < currentPeriod {
				if k.finalizeCandle(ctx, candle) {
					k.DeleteCurrentCandle(ctx, pool.ID, interval)
				}
			}
		}
	}
}

// finalizeCandle persists a completed candle to the history store and prunes old ones.
// Returns true if the candle was successfully persisted.
func (k Keeper) finalizeCandle(ctx context.Context, candle types.Candle) bool {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(candle)
	if err != nil {
		return false
	}
	kvStore.Set(types.CandleKey(candle.PoolID, candle.Interval, candle.OpenBlock), bz)

	// Prune if over max
	k.pruneCandles(ctx, candle.PoolID, candle.Interval)
	return true
}

// pruneCandles removes the oldest candles if we exceed MaxCandlesPerInterval.
func (k Keeper) pruneCandles(ctx context.Context, poolID uint64, interval string) {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := types.CandlePoolIntervalPrefix(poolID, interval)

	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return
	}
	defer iter.Close()

	// Collect all keys
	var keys [][]byte
	for ; iter.Valid(); iter.Next() {
		keys = append(keys, append([]byte{}, iter.Key()...))
	}

	// Delete oldest if over max
	excess := len(keys) - types.MaxCandlesPerInterval
	for i := 0; i < excess; i++ {
		kvStore.Delete(keys[i])
	}
}

// ---------------------------------------------------------------------------
// Current Candle CRUD
// ---------------------------------------------------------------------------

func (k Keeper) GetCurrentCandle(ctx context.Context, poolID uint64, interval string) (types.Candle, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.CandleCurrentKey(poolID, interval))
	if err != nil || bz == nil {
		return types.Candle{}, false
	}
	var candle types.Candle
	if err := json.Unmarshal(bz, &candle); err != nil {
		return types.Candle{}, false
	}
	return candle, true
}

func (k Keeper) SetCurrentCandle(ctx context.Context, candle types.Candle) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(candle)
	if err != nil {
		return
	}
	kvStore.Set(types.CandleCurrentKey(candle.PoolID, candle.Interval), bz)
}

func (k Keeper) DeleteCurrentCandle(ctx context.Context, poolID uint64, interval string) {
	kvStore := k.storeService.OpenKVStore(ctx)
	kvStore.Delete(types.CandleCurrentKey(poolID, interval))
}

// ---------------------------------------------------------------------------
// Query helpers
// ---------------------------------------------------------------------------

// GetCandle returns a single candle for a pool at a specific interval and block height.
func (k Keeper) GetCandle(ctx context.Context, poolID uint64, interval string, blockHeight int64) *types.Candle {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.CandleKey(poolID, interval, blockHeight))
	if err != nil || bz == nil {
		return nil
	}
	var candle types.Candle
	if err := json.Unmarshal(bz, &candle); err != nil {
		return nil
	}
	return &candle
}

// GetCandles returns finalized candles for a pool within a block range, oldest first.
func (k Keeper) GetCandles(ctx context.Context, poolID uint64, interval string, fromBlock, toBlock int64) []types.Candle {
	kvStore := k.storeService.OpenKVStore(ctx)

	startKey := types.CandleKey(poolID, interval, fromBlock)
	// endKey is exclusive — use toBlock+1 so we include candles starting at toBlock
	endKey := types.CandleKey(poolID, interval, toBlock+1)

	iter, err := kvStore.Iterator(startKey, endKey)
	if err != nil {
		return nil
	}
	defer iter.Close()

	var candles []types.Candle
	for ; iter.Valid(); iter.Next() {
		var c types.Candle
		if err := json.Unmarshal(iter.Value(), &c); err != nil {
			continue
		}
		candles = append(candles, c)
	}
	return candles
}
