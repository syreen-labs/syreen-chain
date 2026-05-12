package keeper

import (
	"context"
	"encoding/json"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/dex/types"
)

// ---------------------------------------------------------------------------
// Volume Tracking: KV Store operations
// ---------------------------------------------------------------------------

// getUserVolume loads a user's volume record from the KV store.
func (k Keeper) getUserVolume(ctx context.Context, trader string) types.UserVolume {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.UserVolumePrefix + trader))
	if err != nil || bz == nil {
		return types.UserVolume{
			Address:         trader,
			Volume30d:       math.ZeroInt(),
			LastUpdateBlock: 0,
			DailyVolumes:    []types.DailyVolume{},
		}
	}
	var uv types.UserVolume
	if err := json.Unmarshal(bz, &uv); err != nil {
		return types.UserVolume{
			Address:         trader,
			Volume30d:       math.ZeroInt(),
			LastUpdateBlock: 0,
			DailyVolumes:    []types.DailyVolume{},
		}
	}
	return uv
}

// setUserVolume stores a user's volume record.
func (k Keeper) setUserVolume(ctx context.Context, uv types.UserVolume) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(uv)
	if err != nil {
		return
	}
	kvStore.Set([]byte(types.UserVolumePrefix+uv.Address), bz)
}

// ---------------------------------------------------------------------------
// RecordTradeVolume — called after each trade to update 30-day rolling volume
// ---------------------------------------------------------------------------

// RecordTradeVolume records trade volume for a user, accumulating into the
// current block-day bucket. The volumeUsyreen parameter is the trade volume
// in usyreen (base units).
func (k Keeper) RecordTradeVolume(ctx sdk.Context, trader string, volumeUsyreen math.Int) {
	if volumeUsyreen.IsNil() || !volumeUsyreen.IsPositive() {
		return
	}

	uv := k.getUserVolume(ctx, trader)
	height := ctx.BlockHeight()
	currentDay := height / types.BlocksPerDay

	// Find or create the bucket for the current day
	found := false
	for i := range uv.DailyVolumes {
		if uv.DailyVolumes[i].Day == currentDay {
			uv.DailyVolumes[i].Volume = uv.DailyVolumes[i].Volume.Add(volumeUsyreen)
			found = true
			break
		}
	}
	if !found {
		uv.DailyVolumes = append(uv.DailyVolumes, types.DailyVolume{
			Day:    currentDay,
			Volume: volumeUsyreen,
		})
	}

	// Prune expired buckets (older than 30 days)
	cutoffDay := currentDay - types.RollingWindowDays + 1
	pruned := make([]types.DailyVolume, 0, len(uv.DailyVolumes))
	for _, dv := range uv.DailyVolumes {
		if dv.Day >= cutoffDay {
			pruned = append(pruned, dv)
		}
	}
	uv.DailyVolumes = pruned

	// Recompute 30-day total
	total := math.ZeroInt()
	for _, dv := range uv.DailyVolumes {
		total = total.Add(dv.Volume)
	}
	uv.Volume30d = total
	uv.LastUpdateBlock = height

	k.setUserVolume(ctx, uv)
}

// ---------------------------------------------------------------------------
// GetUserVolume30d — returns 30-day rolling volume for a user
// ---------------------------------------------------------------------------

func (k Keeper) GetUserVolume30d(ctx context.Context, trader string) math.Int {
	uv := k.getUserVolume(ctx, trader)
	if uv.Volume30d.IsNil() {
		return math.ZeroInt()
	}
	return uv.Volume30d
}

// ---------------------------------------------------------------------------
// GetUserFeeTier — returns the fee tier based on 30-day volume
// ---------------------------------------------------------------------------

func (k Keeper) GetUserFeeTier(ctx context.Context, trader string) types.FeeTier {
	volume := k.GetUserVolume30d(ctx, trader)
	tiers := types.DefaultFeeTiers()

	// Walk tiers from highest to lowest, return the first one the user qualifies for
	for i := len(tiers) - 1; i >= 0; i-- {
		if volume.GTE(tiers[i].MinVolume30d) {
			return tiers[i]
		}
	}

	// Fallback: standard tier
	return tiers[0]
}

// ---------------------------------------------------------------------------
// GetFeeDiscount — returns the taker fee discount percentage for a user
// ---------------------------------------------------------------------------

// GetFeeDiscount returns the taker fee discount as a math.LegacyDec (e.g., 0.10 for 10%).
func (k Keeper) GetFeeDiscount(ctx context.Context, trader string) math.LegacyDec {
	tier := k.GetUserFeeTier(ctx, trader)
	return tier.TakerFeeDiscount
}

// ---------------------------------------------------------------------------
// PruneExpiredVolumes — called in BeginBlock to clean up stale volume entries
// ---------------------------------------------------------------------------

// PruneExpiredVolumes iterates over all user volume records and removes
// daily volume buckets older than 30 days. This keeps KV store size bounded.
func (k Keeper) PruneExpiredVolumes(ctx sdk.Context) {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.UserVolumePrefix)

	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return
	}

	height := ctx.BlockHeight()
	currentDay := height / types.BlocksPerDay
	cutoffDay := currentDay - types.RollingWindowDays + 1

	// Collect entries to modify — do NOT mutate KVStore during iteration
	type pruneAction struct {
		key    []byte
		uv     types.UserVolume
		delete bool
	}
	var actions []pruneAction

	for ; iter.Valid(); iter.Next() {
		var uv types.UserVolume
		if err := json.Unmarshal(iter.Value(), &uv); err != nil {
			continue
		}

		needsPrune := false
		for _, dv := range uv.DailyVolumes {
			if dv.Day < cutoffDay {
				needsPrune = true
				break
			}
		}
		if !needsPrune {
			continue
		}

		pruned := make([]types.DailyVolume, 0, len(uv.DailyVolumes))
		for _, dv := range uv.DailyVolumes {
			if dv.Day >= cutoffDay {
				pruned = append(pruned, dv)
			}
		}
		uv.DailyVolumes = pruned

		total := math.ZeroInt()
		for _, dv := range uv.DailyVolumes {
			total = total.Add(dv.Volume)
		}
		uv.Volume30d = total
		uv.LastUpdateBlock = height

		actions = append(actions, pruneAction{
			key:    append([]byte{}, iter.Key()...),
			uv:     uv,
			delete: len(uv.DailyVolumes) == 0,
		})
	}
	iter.Close()

	// Apply modifications after iterator is closed
	for _, a := range actions {
		if a.delete {
			kvStore.Delete(a.key)
		} else {
			bz, err := json.Marshal(a.uv)
			if err != nil {
				continue
			}
			kvStore.Set(a.key, bz)
		}
	}
}
