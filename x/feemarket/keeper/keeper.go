package keeper

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	stdmath "math"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/feemarket/types"
)

type Keeper struct {
	cdc          codec.Codec
	storeService store.KVStoreService
	bankKeeper   types.BankKeeper
	authority    string
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	bankKeeper types.BankKeeper,
	authority string,
) Keeper {
	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		bankKeeper:   bankKeeper,
		authority:    authority,
	}
}

func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// ---------------------------------------------------------------------------
// Base Fee State
// ---------------------------------------------------------------------------

// GetFeeState returns the current fee state from the store
func (k Keeper) GetFeeState(ctx context.Context) types.FeeState {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.BaseFeeKey))
	if err != nil || bz == nil {
		return types.DefaultFeeState()
	}
	var state types.FeeState
	if err := json.Unmarshal(bz, &state); err != nil {
		// Log corruption instead of silently returning defaults
		k.Logger(ctx).Error("STORE CORRUPTION: failed to unmarshal fee state, using defaults", "error", err)
		return types.DefaultFeeState()
	}
	return state
}

// SetFeeState stores the fee state
func (k Keeper) SetFeeState(ctx context.Context, state types.FeeState) error {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal fee state: %w", err)
	}
	return kvStore.Set([]byte(types.BaseFeeKey), bz)
}

// GetBaseFee returns the current base fee
func (k Keeper) GetBaseFee(ctx context.Context) math.LegacyDec {
	return k.GetFeeState(ctx).BaseFee
}

// SetBaseFee updates only the base fee in the current fee state
func (k Keeper) SetBaseFee(ctx context.Context, fee math.LegacyDec) error {
	state := k.GetFeeState(ctx)
	state.BaseFee = fee
	return k.SetFeeState(ctx, state)
}

// ---------------------------------------------------------------------------
// Previous Block Gas Tracking
// ---------------------------------------------------------------------------

// SetPrevBlockGasUsed stores the gas consumed in this block for next block's fee adjustment
func (k Keeper) SetPrevBlockGasUsed(ctx context.Context, gasUsed uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, gasUsed)
	kvStore.Set([]byte(types.PrevBlockGasKey), bz)
}

// GetPrevBlockGasUsed returns the gas consumed in the previous block
func (k Keeper) GetPrevBlockGasUsed(ctx context.Context) int64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.PrevBlockGasKey))
	if err != nil || bz == nil || len(bz) != 8 {
		return 0
	}
	gasUsed := binary.BigEndian.Uint64(bz)
	// Safe cast: cap at MaxInt64 to prevent overflow
	if gasUsed > uint64(stdmath.MaxInt64) {
		return stdmath.MaxInt64
	}
	return int64(gasUsed)
}

// ---------------------------------------------------------------------------
// EIP-1559 Base Fee Adjustment
// ---------------------------------------------------------------------------

// AdjustBaseFee implements the EIP-1559 style base fee adjustment algorithm.
// If block gas used > target, base fee increases; if below target, it decreases.
// The maximum change per block is baseFee / BaseFeeChangeDenominator (12.5%).
func (k Keeper) AdjustBaseFee(ctx context.Context, blockGasUsed int64) math.LegacyDec {
	state := k.GetFeeState(ctx)
	params := k.GetParams(ctx)

	gasTarget := state.BlockGasTarget
	baseFee := state.BaseFee

	if gasTarget == 0 {
		return baseFee
	}

	changeDenom := math.LegacyNewDec(int64(params.BaseFeeChangeDenominator))

	if blockGasUsed == gasTarget {
		// At target: no change
		return baseFee
	}

	if blockGasUsed > gasTarget {
		// Above target: increase base fee
		// delta = baseFee * (gasUsed - gasTarget) / gasTarget / changeDenom
		gasUsedDec := math.LegacyNewDec(blockGasUsed)
		gasTargetDec := math.LegacyNewDec(gasTarget)
		delta := baseFee.Mul(gasUsedDec.Sub(gasTargetDec)).Quo(gasTargetDec).Quo(changeDenom)

		// Ensure minimum increase of 1 unit to prevent stalling
		if delta.LT(math.LegacyNewDecWithPrec(1, 6)) {
			delta = math.LegacyNewDecWithPrec(1, 6)
		}

		baseFee = baseFee.Add(delta)
	} else {
		// Below target: decrease base fee
		// delta = baseFee * (gasTarget - gasUsed) / gasTarget / changeDenom
		gasUsedDec := math.LegacyNewDec(blockGasUsed)
		gasTargetDec := math.LegacyNewDec(gasTarget)
		delta := baseFee.Mul(gasTargetDec.Sub(gasUsedDec)).Quo(gasTargetDec).Quo(changeDenom)

		baseFee = baseFee.Sub(delta)
	}

	// Clamp to [min, max]
	if baseFee.LT(state.MinBaseFee) {
		baseFee = state.MinBaseFee
	}
	if baseFee.GT(state.MaxBaseFee) {
		baseFee = state.MaxBaseFee
	}

	state.BaseFee = baseFee
	if err := k.SetFeeState(ctx, state); err != nil {
		k.Logger(ctx).Error("failed to persist fee state", "error", err)
	}

	return baseFee
}

// ---------------------------------------------------------------------------
// Fee Calculation
// ---------------------------------------------------------------------------

// CalculateFee returns the required fee for a transaction in a given lane.
// requiredFee = baseFee * laneMultiplier * gasWanted
func (k Keeper) CalculateFee(ctx context.Context, laneName string, gasWanted int64) (math.LegacyDec, error) {
	baseFee := k.GetBaseFee(ctx)
	lane, found := k.GetFeeLane(ctx, laneName)
	if !found {
		return math.LegacyDec{}, types.ErrLaneNotFound.Wrapf("lane: %s", laneName)
	}

	effectiveFee := baseFee.Mul(lane.BaseFeeMultiplier)
	totalFee := effectiveFee.Mul(math.LegacyNewDec(gasWanted))
	return totalFee, nil
}

// ---------------------------------------------------------------------------
// Fee Burning
// ---------------------------------------------------------------------------

// BurnBaseFee burns the base fee portion and returns the tip (remainder).
// burnAmount = fees * burnRatio
// tip = fees - burnAmount (sent to the proposer via fee collector)
func (k Keeper) BurnBaseFee(ctx context.Context, fees sdk.Coins) error {
	params := k.GetParams(ctx)
	if !params.EnableFeeBurn || fees.IsZero() {
		return nil
	}

	burnCoins := sdk.NewCoins()
	for _, coin := range fees {
		burnAmt := params.BurnRatio.MulInt(coin.Amount).TruncateInt()
		if burnAmt.IsPositive() {
			burnCoins = burnCoins.Add(sdk.NewCoin(coin.Denom, burnAmt))
		}
	}

	if !burnCoins.IsZero() {
		if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, burnCoins); err != nil {
			return fmt.Errorf("failed to burn base fee: %w", err)
		}

		if err := k.recordBurn(ctx, burnCoins); err != nil {
			return fmt.Errorf("failed to record burn: %w", err)
		}
	}

	return nil
}

// recordBurn updates the cumulative burn record
func (k Keeper) recordBurn(ctx context.Context, amount sdk.Coins) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	height := sdkCtx.BlockHeight()

	record := k.getBurnRecord(ctx)
	record.Height = height
	record.Amount = amount
	record.CumulativeBurned = record.CumulativeBurned.Add(amount...)

	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal burn record: %w", err)
	}
	return kvStore.Set([]byte(types.BurnTotalKey), bz)
}

func (k Keeper) getBurnRecord(ctx context.Context) types.BurnRecord {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.BurnTotalKey))
	if err != nil || bz == nil {
		return types.BurnRecord{
			CumulativeBurned: sdk.NewCoins(),
		}
	}
	var record types.BurnRecord
	if err := json.Unmarshal(bz, &record); err != nil {
		k.Logger(ctx).Error("STORE CORRUPTION: failed to unmarshal burn record", "error", err)
		return types.BurnRecord{
			CumulativeBurned: sdk.NewCoins(),
		}
	}
	return record
}

// GetBurnStats returns the cumulative burn stats
func (k Keeper) GetBurnStats(ctx context.Context) (totalBurned sdk.Coins, lastBurnHeight int64, lastBurnAmount sdk.Coins) {
	record := k.getBurnRecord(ctx)
	return record.CumulativeBurned, record.Height, record.Amount
}

// ---------------------------------------------------------------------------
// Fee Lanes
// ---------------------------------------------------------------------------

// GetFeeLane returns the fee lane configuration by name
func (k Keeper) GetFeeLane(ctx context.Context, name string) (types.FeeLane, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.LaneConfigKey(name))
	if err != nil || bz == nil {
		return types.FeeLane{}, false
	}
	var lane types.FeeLane
	if err := json.Unmarshal(bz, &lane); err != nil {
		return types.FeeLane{}, false
	}
	return lane, true
}

// SetFeeLane stores a fee lane configuration
func (k Keeper) SetFeeLane(ctx context.Context, lane types.FeeLane) error {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(lane)
	if err != nil {
		return fmt.Errorf("failed to marshal fee lane: %w", err)
	}
	return kvStore.Set(types.LaneConfigKey(lane.Name), bz)
}

// ---------------------------------------------------------------------------
// Params
// ---------------------------------------------------------------------------

// GetParams returns the module parameters
func (k Keeper) GetParams(ctx context.Context) types.Params {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.ParamsKey))
	if err != nil || bz == nil {
		return types.DefaultParams()
	}
	var params types.Params
	if err := json.Unmarshal(bz, &params); err != nil {
		k.Logger(ctx).Error("STORE CORRUPTION: failed to unmarshal feemarket params", "error", err)
		return types.DefaultParams()
	}
	return params
}

// SetParams stores the module parameters
func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(params)
	if err != nil {
		return err
	}
	return kvStore.Set([]byte(types.ParamsKey), bz)
}

// ---------------------------------------------------------------------------
// ABCI Hooks
// ---------------------------------------------------------------------------

// BeginBlocker is a no-op. Base fee stays at minimum floor.
// Lane utilization resets are unnecessary when no dynamic fee routing is active.
// Previously this caused state writes every block that contributed to AppHash divergence.
func (k Keeper) BeginBlocker(ctx context.Context) error {
	return nil
}

// EndBlocker is a no-op. Gas tracking has been removed from state because
// BlockGasMeter().GasConsumed() returns different values across binary
// compilations (Go compiler optimizations change BeginBlock/EndBlock gas),
// causing app hash divergence. The base fee is kept at the minimum floor
// and adjusts only based on transaction gas (tracked via ante handler).
func (k Keeper) EndBlocker(ctx context.Context) error {
	return nil
}
