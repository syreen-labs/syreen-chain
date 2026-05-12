package keeper

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ---------------------------------------------------------------------------
// Time-Weighted Governance
//
// Voting power is multiplied by how long tokens have been staked.
// This prevents governance attacks from flash-stakers.
//
// Time multiplier = min(stakeDurationDays / 365, 3.0) — caps at 3x after 1 year.
// Effective voting power = tokenAmount * timeMultiplier
// ---------------------------------------------------------------------------

const (
	// KV store prefixes
	stakeStartPrefix      = "stake_start/"       // delegator -> start_height (int64)
	stakeLastTotalPrefix  = "stake_last_total/"   // delegator -> last known total delegation (int64 in uSYR)

	// Blocks per day estimate (6.5s block time)
	blocksPerDay = int64(13292)

	// Maximum time multiplier (3x after 365 days)
	maxTimeMultiplier = 3
)

// ---------------------------------------------------------------------------
// Record / Remove Delegation
// ---------------------------------------------------------------------------

// RecordDelegation records or updates the stake start height for a delegator.
//
// When new tokens are delegated by a delegator who already has a recorded start
// height, the start height is updated using a weighted average:
//
//	newStart = (oldStake * oldStart + newStake * currentHeight) / (oldStake + newStake)
//
// This prevents gaming where an old staker delegates additional tokens and
// immediately gets the full 3x time multiplier on the new tokens.
// The previous total delegation is tracked in a separate KV entry so we can
// compute the exact new delegation amount.
func (k Keeper) RecordDelegation(ctx context.Context, delegator string) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentHeight := sdkCtx.BlockHeight()
	kvStore := k.storeService.OpenKVStore(ctx)

	oldStart, exists := k.GetStakeStartHeight(ctx, delegator)

	// Get current total delegation from staking keeper
	currentTotal := math.ZeroInt()
	if k.stakingKeeper != nil {
		delegatorAddr, err := sdk.AccAddressFromBech32(delegator)
		if err == nil {
			delegations, err := k.stakingKeeper.GetDelegatorDelegations(ctx, delegatorAddr, 100)
			if err == nil {
				totalDec := math.LegacyZeroDec()
				for _, del := range delegations {
					totalDec = totalDec.Add(del.GetShares())
				}
				currentTotal = totalDec.TruncateInt()
			}
		}
	}

	if !exists {
		// First delegation — record current height and total
		bz := make([]byte, 8)
		binary.BigEndian.PutUint64(bz, uint64(currentHeight))
		kvStore.Set(stakeStartKey(delegator), bz)

		// Store current total for next weighted average calculation
		k.setLastDelegationTotal(ctx, delegator, currentTotal)

		k.Logger(ctx).Info("time-weight: recorded delegation start",
			"delegator", delegator,
			"start_height", currentHeight,
			"total_staked", currentTotal,
		)
		return
	}

	// Already has a start height — compute weighted average
	if oldStart >= currentHeight {
		// Delegating in same block, just update the total
		k.setLastDelegationTotal(ctx, delegator, currentTotal)
		return
	}

	// Get the last known total to compute newStake = currentTotal - lastTotal
	lastTotal := k.getLastDelegationTotal(ctx, delegator)
	newStake := currentTotal.Sub(lastTotal)

	if newStake.IsPositive() && lastTotal.IsPositive() {
		// Weighted average: newStart = (oldStake * oldStart + newStake * currentHeight) / totalStake
		oldComponent := lastTotal.MulRaw(oldStart)
		newComponent := newStake.MulRaw(currentHeight)
		newStart := oldComponent.Add(newComponent).Quo(currentTotal).Int64()

		bz := make([]byte, 8)
		binary.BigEndian.PutUint64(bz, uint64(newStart))
		kvStore.Set(stakeStartKey(delegator), bz)

		k.Logger(ctx).Info("time-weight: updated delegation start (weighted average)",
			"delegator", delegator,
			"old_start", oldStart,
			"new_start", newStart,
			"old_stake", lastTotal,
			"new_stake", newStake,
			"current_height", currentHeight,
		)
	}

	// Update stored total
	k.setLastDelegationTotal(ctx, delegator, currentTotal)
}

// setLastDelegationTotal stores the last known total delegation for weighted average calculations.
func (k Keeper) setLastDelegationTotal(ctx context.Context, delegator string, total math.Int) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, total.Uint64())
	kvStore.Set([]byte(stakeLastTotalPrefix+delegator), bz)
}

// getLastDelegationTotal retrieves the last known total delegation.
func (k Keeper) getLastDelegationTotal(ctx context.Context, delegator string) math.Int {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(stakeLastTotalPrefix + delegator))
	if err != nil || bz == nil || len(bz) < 8 {
		return math.ZeroInt()
	}
	return math.NewIntFromUint64(binary.BigEndian.Uint64(bz))
}

// RemoveDelegation removes the stake start record for a delegator.
func (k Keeper) RemoveDelegation(ctx context.Context, delegator string) {
	kvStore := k.storeService.OpenKVStore(ctx)
	kvStore.Delete(stakeStartKey(delegator))

	k.Logger(ctx).Info("time-weight: removed delegation record",
		"delegator", delegator,
	)
}

// ---------------------------------------------------------------------------
// Queries
// ---------------------------------------------------------------------------

// GetStakeStartHeight returns the block height when the delegator first staked.
func (k Keeper) GetStakeStartHeight(ctx context.Context, delegator string) (int64, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(stakeStartKey(delegator))
	if err != nil || bz == nil {
		return 0, false
	}
	if len(bz) < 8 {
		return 0, false
	}
	height := int64(binary.BigEndian.Uint64(bz))
	return height, true
}

// GetTimeMultiplier calculates the time multiplier for a delegator.
// Formula: min(stakeDurationDays / 365, 3.0)
// Returns 0 if the delegator has no recorded stake start.
func (k Keeper) GetTimeMultiplier(ctx context.Context, delegator string) math.LegacyDec {
	startHeight, exists := k.GetStakeStartHeight(ctx, delegator)
	if !exists {
		return math.LegacyZeroDec()
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentHeight := sdkCtx.BlockHeight()

	// Calculate duration in days
	blocksDelta := currentHeight - startHeight
	if blocksDelta < 0 {
		blocksDelta = 0
	}

	// durationDays = blocksDelta / blocksPerDay
	durationDays := math.LegacyNewDec(blocksDelta).Quo(math.LegacyNewDec(blocksPerDay))

	// multiplier = durationDays / 365
	multiplier := durationDays.Quo(math.LegacyNewDec(365))

	// Cap at 3.0
	maxMul := math.LegacyNewDec(maxTimeMultiplier)
	if multiplier.GT(maxMul) {
		multiplier = maxMul
	}

	return multiplier
}

// GetTimeWeightedVotingPower returns the effective voting power for a delegator,
// which is their total staked tokens multiplied by the time multiplier.
// Requires the staking keeper to be set via SetStakingKeeper.
func (k Keeper) GetTimeWeightedVotingPower(ctx context.Context, delegator string) (math.Int, error) {
	if k.stakingKeeper == nil {
		return math.ZeroInt(), fmt.Errorf("staking keeper not set")
	}

	delegatorAddr, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		return math.ZeroInt(), fmt.Errorf("invalid delegator address: %w", err)
	}

	// Get all delegations for this delegator
	delegations, err := k.stakingKeeper.GetDelegatorDelegations(ctx, delegatorAddr, 100)
	if err != nil {
		return math.ZeroInt(), fmt.Errorf("failed to get delegations: %w", err)
	}

	// Sum total staked tokens
	totalStaked := math.LegacyZeroDec()
	for _, del := range delegations {
		totalStaked = totalStaked.Add(del.GetShares())
	}

	if totalStaked.IsZero() {
		return math.ZeroInt(), nil
	}

	// Get time multiplier
	multiplier := k.GetTimeMultiplier(ctx, delegator)
	if multiplier.IsZero() {
		// No recorded start but has delegations — use base multiplier of 0
		return math.ZeroInt(), nil
	}

	// Effective power = totalStaked * multiplier
	effectivePower := totalStaked.Mul(multiplier).TruncateInt()

	return effectivePower, nil
}

// TimeWeightedPowerInfo is the response structure for time-weighted power queries.
type TimeWeightedPowerInfo struct {
	Address        string         `json:"address"`
	StartHeight    int64          `json:"start_height"`
	CurrentHeight  int64          `json:"current_height"`
	StakeDuration  int64          `json:"stake_duration_blocks"`
	DurationDays   math.LegacyDec `json:"duration_days"`
	TimeMultiplier math.LegacyDec `json:"time_multiplier"`
	BaseVotingPower math.Int      `json:"base_voting_power"`
	EffectivePower math.Int       `json:"effective_voting_power"`
}

// GetTimeWeightedPowerInfo returns detailed time-weighted power information.
func (k Keeper) GetTimeWeightedPowerInfo(ctx context.Context, delegator string) (*TimeWeightedPowerInfo, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentHeight := sdkCtx.BlockHeight()

	info := &TimeWeightedPowerInfo{
		Address:       delegator,
		CurrentHeight: currentHeight,
	}

	startHeight, exists := k.GetStakeStartHeight(ctx, delegator)
	if exists {
		info.StartHeight = startHeight
		blocksDelta := currentHeight - startHeight
		if blocksDelta < 0 {
			blocksDelta = 0
		}
		info.StakeDuration = blocksDelta
		info.DurationDays = math.LegacyNewDec(blocksDelta).Quo(math.LegacyNewDec(blocksPerDay))
	} else {
		info.DurationDays = math.LegacyZeroDec()
	}

	info.TimeMultiplier = k.GetTimeMultiplier(ctx, delegator)

	// Get base voting power from staking if available
	if k.stakingKeeper != nil {
		delegatorAddr, err := sdk.AccAddressFromBech32(delegator)
		if err == nil {
			delegations, err := k.stakingKeeper.GetDelegatorDelegations(ctx, delegatorAddr, 100)
			if err == nil {
				totalStaked := math.LegacyZeroDec()
				for _, del := range delegations {
					totalStaked = totalStaked.Add(del.GetShares())
				}
				info.BaseVotingPower = totalStaked.TruncateInt()
			} else {
				info.BaseVotingPower = math.ZeroInt()
			}
		} else {
			info.BaseVotingPower = math.ZeroInt()
		}
	} else {
		info.BaseVotingPower = math.ZeroInt()
	}

	// Effective power
	effectivePower, err := k.GetTimeWeightedVotingPower(ctx, delegator)
	if err != nil {
		info.EffectivePower = math.ZeroInt()
	} else {
		info.EffectivePower = effectivePower
	}

	return info, nil
}

// GetAllStakeStarts returns all recorded delegation start heights (for genesis export).
func (k Keeper) GetAllStakeStarts(ctx context.Context) map[string]int64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(stakeStartPrefix)

	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	result := make(map[string]int64)
	for ; iter.Valid(); iter.Next() {
		key := string(iter.Key())
		delegator := key[len(stakeStartPrefix):]
		if len(iter.Value()) >= 8 {
			height := int64(binary.BigEndian.Uint64(iter.Value()))
			result[delegator] = height
		}
	}
	return result
}

// ExportStakeStarts exports all stake start records as JSON.
func (k Keeper) ExportStakeStarts(ctx context.Context) json.RawMessage {
	starts := k.GetAllStakeStarts(ctx)
	if starts == nil {
		starts = make(map[string]int64)
	}
	bz, _ := json.Marshal(starts)
	return bz
}

// ---------------------------------------------------------------------------
// KV Store Key
// ---------------------------------------------------------------------------

func stakeStartKey(delegator string) []byte {
	return []byte(stakeStartPrefix + delegator)
}
