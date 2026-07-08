package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"syreen/x/mevprotection/types"
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
	// MEVRewardPoolPrefix stores accumulated MEV rewards by denom (lifetime stats)
	MEVRewardPoolPrefix = "mev_reward_pool/"

	// MEVTotalRedistributedPrefix stores lifetime redistribution stats by denom
	MEVTotalRedistributedPrefix = "mev_total_redistributed/"

	// RebateWeightPrefix stores each eligible user's rebate weight (a trade-size
	// proxy) for the current distribution window: rebate_weight/<addr> -> Int.
	RebateWeightPrefix = "rebate_weight/"

	// MEVDistributionInterval is the number of blocks between distributions
	MEVDistributionInterval = int64(100)

	// Fairness Engine · Return: the harmed/trading USER is the first-priority
	// recipient, ahead of LPs and stakers. Shares are basis points and sum to 10000.
	UserRebateShareBps   = int64(5000) // 50% -> users (rebated first)
	LPRewardShareBps     = int64(3000) // 30% -> LPs (added to DEX reserves)
	StakerRewardShareBps = int64(2000) // 20% -> stakers (via fee collector)
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

// CreditFairnessPool moves REAL coins from another module account into the
// fairness pool (this module's account) and, if a beneficiary is given, records
// them as eligible for the user-first rebate weighted by the credited amount.
// This is the single funding entry point (e.g. intent routes slashed solver
// stake here instead of orphaning it in the distribution account).
func (k Keeper) CreditFairnessPool(ctx context.Context, fromModule string, coins sdk.Coins, beneficiary string) error {
	if coins.IsZero() {
		return nil
	}
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, fromModule, types.ModuleName, coins); err != nil {
		return err
	}
	// Mirror into the lifetime accumulator (drives the MEVRedistribution query).
	k.AccumulateMEVReward(ctx, coins)
	if beneficiary != "" {
		k.RecordRebateBeneficiary(ctx, beneficiary, coins.AmountOf(sdk.DefaultBondDenom))
	}
	return nil
}

// RecordRebateBeneficiary adds rebate weight (a trade-size proxy, e.g. the
// intent's tip/fee) for a user so the next distribution rebates them first.
// Callers on trade fulfillment mark the trading user here.
func (k Keeper) RecordRebateBeneficiary(ctx context.Context, addr string, weight math.Int) {
	if addr == "" || weight.IsNil() || !weight.IsPositive() {
		return
	}
	if _, err := sdk.AccAddressFromBech32(addr); err != nil {
		return
	}
	kv := k.storeService.OpenKVStore(ctx)
	key := []byte(RebateWeightPrefix + addr)
	existing := math.ZeroInt()
	if bz, err := kv.Get(key); err == nil && bz != nil {
		_ = existing.Unmarshal(bz)
	}
	total := existing.Add(weight)
	if bz, err := total.Marshal(); err == nil {
		_ = kv.Set(key, bz)
	}
}

// rebateBeneficiaries returns the eligible users and their weights, sorted by
// address for deterministic distribution.
func (k Keeper) rebateBeneficiaries(ctx context.Context) (addrs []string, weights []math.Int, total math.Int) {
	kv := k.storeService.OpenKVStore(ctx)
	prefix := []byte(RebateWeightPrefix)
	iter, err := kv.Iterator(prefix, prefixEndBytes(prefix))
	total = math.ZeroInt()
	if err != nil {
		return nil, nil, total
	}
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		addr := string(iter.Key()[len(prefix):])
		w := math.ZeroInt()
		if err := w.Unmarshal(iter.Value()); err != nil || !w.IsPositive() {
			continue
		}
		addrs = append(addrs, addr)
		weights = append(weights, w)
		total = total.Add(w)
	}
	sort.Slice(addrs, func(i, j int) bool { return addrs[i] < addrs[j] })
	// re-derive weights in the sorted order
	wmap := make(map[string]math.Int, len(addrs))
	for i := range addrs {
		wmap[addrs[i]] = weights[i]
	}
	sortedW := make([]math.Int, len(addrs))
	for i, a := range addrs {
		sortedW[i] = wmap[a]
	}
	return addrs, sortedW, total
}

func (k Keeper) clearRebateLedger(ctx context.Context) {
	kv := k.storeService.OpenKVStore(ctx)
	prefix := []byte(RebateWeightPrefix)
	iter, err := kv.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return
	}
	var keys [][]byte
	for ; iter.Valid(); iter.Next() {
		keys = append(keys, append([]byte{}, iter.Key()...))
	}
	iter.Close()
	for _, key := range keys {
		_ = kv.Delete(key)
	}
}

// DistributeMEVRewards pays out the fairness pool's REAL balance every
// MEVDistributionInterval blocks, user-first: 50% rebated to trading users
// (proportional to weight), 30% to LPs (into the DEX module account, boosting
// reserves), 20% to stakers (via the fee collector, distributed next block).
// Unlike the previous accounting-only stub, this moves actual coins.
func (k Keeper) DistributeMEVRewards(ctx context.Context) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Only distribute every 100 blocks
	if sdkCtx.BlockHeight()%MEVDistributionInterval != 0 {
		return
	}

	poolAddr := authtypes.NewModuleAddress(types.ModuleName)
	pool := k.bankKeeper.GetAllBalances(ctx, poolAddr)
	if pool.IsZero() {
		return
	}

	users, weights, totalWeight := k.rebateBeneficiaries(ctx)
	hasUsers := len(users) > 0 && totalWeight.IsPositive()

	var distributed sdk.Coins
	for _, coin := range pool {
		userAmt := coin.Amount.MulRaw(UserRebateShareBps).QuoRaw(10000)
		lpAmt := coin.Amount.MulRaw(LPRewardShareBps).QuoRaw(10000)
		// Stakers get the remainder so nothing is lost to rounding.
		stakerAmt := coin.Amount.Sub(userAmt).Sub(lpAmt)

		// USER tranche (first). If there are no eligible users, fold it into the
		// staker tranche rather than stranding it.
		if hasUsers && userAmt.IsPositive() {
			paid := math.ZeroInt()
			for i, addr := range users {
				share := userAmt.Mul(weights[i]).Quo(totalWeight)
				if i == len(users)-1 {
					share = userAmt.Sub(paid) // last user absorbs rounding dust
				}
				if !share.IsPositive() {
					continue
				}
				acc, err := sdk.AccAddressFromBech32(addr)
				if err != nil {
					stakerAmt = stakerAmt.Add(share)
					paid = paid.Add(share)
					continue
				}
				if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, acc, sdk.NewCoins(sdk.NewCoin(coin.Denom, share))); err != nil {
					stakerAmt = stakerAmt.Add(share)
				}
				paid = paid.Add(share)
			}
		} else {
			stakerAmt = stakerAmt.Add(userAmt)
		}

		// LP tranche -> DEX module account (increases reserves).
		if lpAmt.IsPositive() {
			if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, "dex", sdk.NewCoins(sdk.NewCoin(coin.Denom, lpAmt))); err != nil {
				stakerAmt = stakerAmt.Add(lpAmt) // fall back to stakers if DEX send fails
			}
		}

		// STAKER tranche -> fee collector (distribution pays stakers next block).
		if stakerAmt.IsPositive() {
			if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, authtypes.FeeCollectorName, sdk.NewCoins(sdk.NewCoin(coin.Denom, stakerAmt))); err != nil {
				k.Logger(ctx).Error("fairness: staker tranche send failed", "denom", coin.Denom, "error", err)
				continue
			}
		}
		distributed = distributed.Add(coin)
	}

	k.Logger(ctx).Info("fairness pool distributed (real coins)",
		"total", distributed.String(), "users", len(users), "height", sdkCtx.BlockHeight())

	// Update lifetime stats and clear the window ledgers.
	k.addToTotalRedistributed(ctx, distributed)
	k.clearRebateLedger(ctx)
	// The reward-pool counter mirrors the pool; clear it now that we paid out.
	k.clearRewardPool(ctx)

	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"mev_redistribution",
		sdk.NewAttribute("distributed", distributed.String()),
		sdk.NewAttribute("users_rebated", fmt.Sprintf("%d", len(users))),
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
