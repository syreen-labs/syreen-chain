package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/dex/types"
)

// PoolRiskState tracks real-time risk metrics for a pool
type PoolRiskState struct {
	PoolID          uint64         `json:"pool_id"`
	PrevReserveA    math.Int       `json:"prev_reserve_a"`
	PrevReserveB    math.Int       `json:"prev_reserve_b"`
	PrevPrice       math.LegacyDec `json:"prev_price"`       // price at last snapshot
	SnapshotHeight  int64          `json:"snapshot_height"`   // block of last snapshot
	HaltedUntil     int64          `json:"halted_until"`      // block height when halt expires (0 = not halted)
	HaltReason      string         `json:"halt_reason"`
	PriceChangePct  math.LegacyDec `json:"price_change_pct"`  // cumulative price change since snapshot
	VolumeThisBlock math.Int       `json:"volume_this_block"` // swap volume in current block
}

const (
	riskStatePrefix    = "risk_state/"
	riskSnapshotWindow = int64(100) // blocks between price snapshots
)

// Risk thresholds (governable via params in production — hardcoded for now)
var (
	maxPriceDropPct     = math.LegacyNewDecWithPrec(30, 2)  // 30% price drop triggers halt
	maxPriceSpikePct    = math.LegacyNewDecWithPrec(50, 2)  // 50% price spike triggers halt
	maxReserveDrainPct  = math.LegacyNewDecWithPrec(80, 2)  // 80% reserve drain triggers halt
	haltDurationBlocks  = int64(100)                         // halt lasts 100 blocks (~50 seconds)
	maxVolumePerBlock   = math.NewInt(500000000000)          // 500K token volume per block per pool
)

func riskStateKey(poolID uint64) []byte {
	return append([]byte(riskStatePrefix), types.PoolKey(poolID)...)
}

func (k Keeper) GetRiskState(ctx context.Context, poolID uint64) (PoolRiskState, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(riskStateKey(poolID))
	if err != nil || bz == nil {
		return PoolRiskState{}, false
	}
	var state PoolRiskState
	if err := json.Unmarshal(bz, &state); err != nil {
		return PoolRiskState{}, false
	}
	return state, true
}

func (k Keeper) SetRiskState(ctx context.Context, state PoolRiskState) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(state)
	kvStore.Set(riskStateKey(state.PoolID), bz)
}

// IsPoolHalted returns true if the pool is currently halted by the risk engine.
func (k Keeper) IsPoolHalted(ctx context.Context, poolID uint64) (bool, string) {
	state, found := k.GetRiskState(ctx, poolID)
	if !found {
		return false, ""
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if state.HaltedUntil > 0 && sdkCtx.BlockHeight() < state.HaltedUntil {
		return true, state.HaltReason
	}
	return false, ""
}

// CheckPoolRisk evaluates risk for a single pool and halts if thresholds are breached.
// Called in BeginBlock for every active pool.
func (k Keeper) CheckPoolRisk(ctx context.Context, pool types.Pool) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	state, found := k.GetRiskState(ctx, pool.ID)
	if !found {
		// Initialize risk state for this pool
		state = PoolRiskState{
			PoolID:         pool.ID,
			PrevReserveA:   pool.ReserveA,
			PrevReserveB:   pool.ReserveB,
			SnapshotHeight: sdkCtx.BlockHeight(),
			VolumeThisBlock: math.ZeroInt(),
		}
		// Calculate initial price
		if !pool.ReserveA.IsZero() {
			state.PrevPrice = math.LegacyNewDecFromInt(pool.ReserveB).Quo(math.LegacyNewDecFromInt(pool.ReserveA))
		}
		k.SetRiskState(ctx, state)
		return
	}

	// If currently halted, check if halt has expired
	if state.HaltedUntil > 0 && sdkCtx.BlockHeight() >= state.HaltedUntil {
		// Halt expired — reset and take new snapshot
		state.HaltedUntil = 0
		state.HaltReason = ""
		state.PrevReserveA = pool.ReserveA
		state.PrevReserveB = pool.ReserveB
		state.SnapshotHeight = sdkCtx.BlockHeight()
		state.PriceChangePct = math.LegacyZeroDec()
		state.VolumeThisBlock = math.ZeroInt()
		if !pool.ReserveA.IsZero() {
			state.PrevPrice = math.LegacyNewDecFromInt(pool.ReserveB).Quo(math.LegacyNewDecFromInt(pool.ReserveA))
		}
		k.SetRiskState(ctx, state)
		k.Logger(ctx).Info("pool halt expired, trading resumed", "pool_id", pool.ID)
		return
	}

	// Skip checks if already halted
	if state.HaltedUntil > 0 {
		return
	}

	// Check 1: Reserve drain — if either reserve dropped by more than threshold
	// D-14: Skip drain check if reserves grew (drain is zero or negative = no drain happening).
	if !state.PrevReserveA.IsZero() && pool.ReserveA.LT(state.PrevReserveA) {
		drainA := math.LegacyNewDecFromInt(state.PrevReserveA.Sub(pool.ReserveA)).Quo(math.LegacyNewDecFromInt(state.PrevReserveA))
		if drainA.GT(maxReserveDrainPct) {
			k.haltPool(ctx, &state, fmt.Sprintf("reserve_a drained %s%% since block %d", drainA.MulInt64(100).TruncateInt(), state.SnapshotHeight))
			k.SetRiskState(ctx, state)
			return
		}
	}
	if !state.PrevReserveB.IsZero() && pool.ReserveB.LT(state.PrevReserveB) {
		drainB := math.LegacyNewDecFromInt(state.PrevReserveB.Sub(pool.ReserveB)).Quo(math.LegacyNewDecFromInt(state.PrevReserveB))
		if drainB.GT(maxReserveDrainPct) {
			k.haltPool(ctx, &state, fmt.Sprintf("reserve_b drained %s%% since block %d", drainB.MulInt64(100).TruncateInt(), state.SnapshotHeight))
			k.SetRiskState(ctx, state)
			return
		}
	}

	// Check 2: Price crash / spike
	if !pool.ReserveA.IsZero() && !state.PrevPrice.IsZero() {
		currentPrice := math.LegacyNewDecFromInt(pool.ReserveB).Quo(math.LegacyNewDecFromInt(pool.ReserveA))
		priceChange := currentPrice.Sub(state.PrevPrice).Quo(state.PrevPrice)

		// Price dropped more than threshold (negative change)
		if priceChange.IsNegative() && priceChange.Abs().GT(maxPriceDropPct) {
			k.haltPool(ctx, &state, fmt.Sprintf("price dropped %s%% (from %s to %s)", priceChange.Abs().MulInt64(100).TruncateInt(), state.PrevPrice, currentPrice))
			k.SetRiskState(ctx, state)
			return
		}

		// Price spiked more than threshold (positive change)
		if priceChange.IsPositive() && priceChange.GT(maxPriceSpikePct) {
			k.haltPool(ctx, &state, fmt.Sprintf("price spiked %s%% (from %s to %s)", priceChange.MulInt64(100).TruncateInt(), state.PrevPrice, currentPrice))
			k.SetRiskState(ctx, state)
			return
		}

		state.PriceChangePct = priceChange
	}

	// Update snapshot if window has elapsed
	if sdkCtx.BlockHeight()-state.SnapshotHeight >= riskSnapshotWindow {
		state.PrevReserveA = pool.ReserveA
		state.PrevReserveB = pool.ReserveB
		state.SnapshotHeight = sdkCtx.BlockHeight()
		if !pool.ReserveA.IsZero() {
			state.PrevPrice = math.LegacyNewDecFromInt(pool.ReserveB).Quo(math.LegacyNewDecFromInt(pool.ReserveA))
		}
		state.PriceChangePct = math.LegacyZeroDec()
	}

	// Reset per-block volume
	state.VolumeThisBlock = math.ZeroInt()
	k.SetRiskState(ctx, state)
}

func (k Keeper) haltPool(ctx context.Context, state *PoolRiskState, reason string) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	state.HaltedUntil = sdkCtx.BlockHeight() + haltDurationBlocks
	state.HaltReason = reason

	k.Logger(ctx).Error("RISK ENGINE: Pool halted",
		"pool_id", state.PoolID,
		"reason", reason,
		"halted_until", state.HaltedUntil,
	)

	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"pool_halted",
		sdk.NewAttribute("pool_id", fmt.Sprintf("%d", state.PoolID)),
		sdk.NewAttribute("reason", reason),
		sdk.NewAttribute("halted_until", fmt.Sprintf("%d", state.HaltedUntil)),
	))
}

// CheckSwapRisk is called before every swap to enforce per-block volume limits.
// Returns an error if the swap should be rejected.
//
// NOTE: This function does NOT persist the volume increment. The caller must
// call RecordSwapVolume after the swap succeeds to actually record it. This
// prevents failed swaps from inflating the per-block volume counter.
func (k Keeper) CheckSwapRisk(ctx context.Context, poolID uint64, swapAmount math.Int) error {
	// Check if pool is halted
	halted, reason := k.IsPoolHalted(ctx, poolID)
	if halted {
		return fmt.Errorf("pool %d is halted by risk engine: %s", poolID, reason)
	}

	// Check per-block volume limit (read-only — do not persist yet)
	state, found := k.GetRiskState(ctx, poolID)
	if found {
		newVolume := state.VolumeThisBlock.Add(swapAmount)
		if newVolume.GT(maxVolumePerBlock) {
			return fmt.Errorf("pool %d: per-block volume limit exceeded (%s > %s)", poolID, newVolume, maxVolumePerBlock)
		}
	}

	return nil
}

// RecordSwapVolume persists the swap volume to the risk state AFTER a swap
// succeeds. Must be called after every successful swap to keep the per-block
// volume counter accurate.
func (k Keeper) RecordSwapVolume(ctx context.Context, poolID uint64, swapAmount math.Int) {
	state, found := k.GetRiskState(ctx, poolID)
	if !found {
		return
	}
	state.VolumeThisBlock = state.VolumeThisBlock.Add(swapAmount)
	k.SetRiskState(ctx, state)
}

// RunRiskChecks runs risk analysis for all pools. Called in BeginBlock.
func (k Keeper) RunRiskChecks(ctx context.Context) {
	pools := k.GetAllPools(ctx)
	for _, pool := range pools {
		k.CheckPoolRisk(ctx, pool)
	}
}
