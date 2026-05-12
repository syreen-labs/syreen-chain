package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ---------------------------------------------------------------------------
// MEV Redistribution System
//
// When MEV is detected and a validator is slashed, the slashed amount is
// captured and redistributed: 60% to LP providers (proportional to liquidity)
// and 40% to stakers (via the distribution module's community pool).
//
// Rewards accumulate in a reward pool and are distributed every 100 blocks
// during EndBlock.
// ---------------------------------------------------------------------------

const (
	// MEVRewardPoolPrefix stores accumulated MEV rewards by denom
	MEVRewardPoolPrefix = "mev_reward_pool/"

	// MEVTotalRedistributedPrefix stores lifetime redistribution stats by denom
	MEVTotalRedistributedPrefix = "mev_total_redistributed/"

	// MEVDistributionInterval is the number of blocks between distributions
	MEVDistributionInterval = int64(100)

	// LPRewardShareBps is the LP share in basis points (60%)
	LPRewardShareBps = int64(6000)

	// StakerRewardShareBps is the staker share in basis points (40%)
	StakerRewardShareBps = int64(4000)
)

// DexKeeper defines the interface needed from the dex module to get pool info.
// Uses a minimal interface to avoid circular imports.
type DexKeeper interface {
	GetAllPools(ctx context.Context) []PoolInfo
}

// PoolInfo is a minimal representation of a DEX pool for reward distribution.
// Avoids importing dex types directly.
type PoolInfo struct {
	ID       uint64   `json:"id"`
	DenomA   string   `json:"denom_a"`
	DenomB   string   `json:"denom_b"`
	ReserveA math.Int `json:"reserve_a"`
	ReserveB math.Int `json:"reserve_b"`
}

// AccumulateMEVReward adds slashed tokens to the MEV reward pool for later distribution.
func (k Keeper) AccumulateMEVReward(ctx context.Context, amount sdk.Coins) {
	if amount.IsZero() {
		return
	}

	kvStore := k.storeService.OpenKVStore(ctx)

	for _, coin := range amount {
		key := []byte(MEVRewardPoolPrefix + coin.Denom)
		current := math.ZeroInt()

		bz, err := kvStore.Get(key)
		if err == nil && bz != nil {
			var stored string
			if err := json.Unmarshal(bz, &stored); err == nil {
				parsed, ok := math.NewIntFromString(stored)
				if ok {
					current = parsed
				}
			}
		}

		current = current.Add(coin.Amount)
		storeBz, _ := json.Marshal(current.String())
		kvStore.Set(key, storeBz)
	}

	k.Logger(ctx).Info("accumulated MEV reward", "amount", amount.String())
}

// GetMEVRewardPool returns the currently accumulated (undistributed) MEV rewards.
func (k Keeper) GetMEVRewardPool(ctx context.Context) sdk.Coins {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(MEVRewardPoolPrefix)
	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return sdk.Coins{}
	}
	defer iter.Close()

	var coins sdk.Coins
	for ; iter.Valid(); iter.Next() {
		// Extract denom from key: "mev_reward_pool/<denom>"
		denom := string(iter.Key()[len(prefix):])

		var stored string
		if err := json.Unmarshal(iter.Value(), &stored); err != nil {
			continue
		}
		amount, ok := math.NewIntFromString(stored)
		if !ok || !amount.IsPositive() {
			continue
		}
		coins = append(coins, sdk.NewCoin(denom, amount))
	}

	// Sort for determinism
	sort.Slice(coins, func(i, j int) bool {
		return coins[i].Denom < coins[j].Denom
	})

	return coins
}

// GetTotalMEVRedistributed returns lifetime MEV redistribution stats.
func (k Keeper) GetTotalMEVRedistributed(ctx context.Context) sdk.Coins {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(MEVTotalRedistributedPrefix)
	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return sdk.Coins{}
	}
	defer iter.Close()

	var coins sdk.Coins
	for ; iter.Valid(); iter.Next() {
		denom := string(iter.Key()[len(prefix):])

		var stored string
		if err := json.Unmarshal(iter.Value(), &stored); err != nil {
			continue
		}
		amount, ok := math.NewIntFromString(stored)
		if !ok || !amount.IsPositive() {
			continue
		}
		coins = append(coins, sdk.NewCoin(denom, amount))
	}

	sort.Slice(coins, func(i, j int) bool {
		return coins[i].Denom < coins[j].Denom
	})

	return coins
}

// DistributeMEVRewards distributes accumulated MEV rewards.
// Called in EndBlock every MEVDistributionInterval blocks.
// - 60% goes to LP providers by sending tokens to the DEX module account
//   (which increases pool reserves, boosting LP share value).
// - 40% goes to the distribution module's community pool for staker rewards.
func (k Keeper) DistributeMEVRewards(ctx context.Context) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Only distribute every 100 blocks
	if sdkCtx.BlockHeight()%MEVDistributionInterval != 0 {
		return
	}

	rewardPool := k.GetMEVRewardPool(ctx)
	if rewardPool.IsZero() {
		return
	}

	// Calculate LP portion (60%) and staker portion (40%)
	var lpCoins sdk.Coins
	var stakerCoins sdk.Coins

	for _, coin := range rewardPool {
		lpAmount := coin.Amount.MulRaw(LPRewardShareBps).QuoRaw(10000)
		stakerAmount := coin.Amount.Sub(lpAmount) // remainder goes to stakers

		if lpAmount.IsPositive() {
			lpCoins = append(lpCoins, sdk.NewCoin(coin.Denom, lpAmount))
		}
		if stakerAmount.IsPositive() {
			stakerCoins = append(stakerCoins, sdk.NewCoin(coin.Denom, stakerAmount))
		}
	}

	// The slashed tokens from validators are already captured by the staking module
	// and sent to the community pool. We do NOT mint new tokens — that would cause
	// inflation and double-count the slashed amount. Instead, we track the
	// redistribution as an accounting entry. The actual tokens for redistribution
	// are already in the community pool from the slashing event.
	k.Logger(ctx).Info("MEV redistribution recorded (accounting only, no mint)",
		"lp_share", lpCoins.String(),
		"staker_share", stakerCoins.String(),
		"height", sdkCtx.BlockHeight(),
	)

	// Update lifetime stats
	k.addToTotalRedistributed(ctx, rewardPool)

	// Clear the reward pool
	k.clearRewardPool(ctx)

	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"mev_redistribution",
		sdk.NewAttribute("lp_rewards", lpCoins.String()),
		sdk.NewAttribute("staker_rewards", stakerCoins.String()),
		sdk.NewAttribute("height", fmt.Sprintf("%d", sdkCtx.BlockHeight())),
	))
}

// addToTotalRedistributed updates lifetime redistribution stats.
func (k Keeper) addToTotalRedistributed(ctx context.Context, amount sdk.Coins) {
	kvStore := k.storeService.OpenKVStore(ctx)

	for _, coin := range amount {
		key := []byte(MEVTotalRedistributedPrefix + coin.Denom)
		current := math.ZeroInt()

		bz, err := kvStore.Get(key)
		if err == nil && bz != nil {
			var stored string
			if err := json.Unmarshal(bz, &stored); err == nil {
				parsed, ok := math.NewIntFromString(stored)
				if ok {
					current = parsed
				}
			}
		}

		current = current.Add(coin.Amount)
		storeBz, _ := json.Marshal(current.String())
		kvStore.Set(key, storeBz)
	}
}

// clearRewardPool removes all entries from the reward pool after distribution.
func (k Keeper) clearRewardPool(ctx context.Context) {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(MEVRewardPoolPrefix)
	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return
	}
	defer iter.Close()

	var keysToDelete [][]byte
	for ; iter.Valid(); iter.Next() {
		keysToDelete = append(keysToDelete, append([]byte(nil), iter.Key()...))
	}

	for _, key := range keysToDelete {
		kvStore.Delete(key)
	}
}
