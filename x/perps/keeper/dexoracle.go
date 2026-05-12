package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/perps/types"
)

// ---------------------------------------------------------------------------
// Oracle-Free Perps Pricing — uses native DEX pool prices instead of
// external oracle feeds. Applies TWAP over the last 100 blocks to prevent
// single-block manipulation.
// ---------------------------------------------------------------------------

const (
	// TWAPWindowBlocks is the number of blocks over which the TWAP is computed.
	TWAPWindowBlocks = 100

	// MaxPriceSamples limits stored price samples to prevent unbounded growth.
	MaxPriceSamples = 200

	// FundingRateClampBps is the maximum funding rate deviation in basis points.
	FundingRateClampBps = 100 // 1%
)

// ---------------------------------------------------------------------------
// Price Sample Storage
// ---------------------------------------------------------------------------

// PriceSample represents a single price observation at a specific block height.
type PriceSample struct {
	Height int64          `json:"height"`
	Price  math.LegacyDec `json:"price"`
}

// RecordPriceSample stores a price observation for a pool at the current block height.
// Called in BeginBlock for each active perps market.
func (k Keeper) RecordPriceSample(ctx context.Context, poolID uint64, price math.LegacyDec) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	height := sdkCtx.BlockHeight()

	samples := k.getPriceSamples(ctx, poolID)

	// Append new sample
	samples = append(samples, PriceSample{
		Height: height,
		Price:  price,
	})

	// Prune old samples beyond MaxPriceSamples
	if len(samples) > MaxPriceSamples {
		samples = samples[len(samples)-MaxPriceSamples:]
	}

	k.setPriceSamples(ctx, poolID, samples)
}

// getPriceSamples retrieves all stored price samples for a pool.
func (k Keeper) getPriceSamples(ctx context.Context, poolID uint64) []PriceSample {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.PerpsTWAPKey(poolID))
	if err != nil || bz == nil {
		return nil
	}
	var samples []PriceSample
	if err := json.Unmarshal(bz, &samples); err != nil {
		return nil
	}
	return samples
}

// setPriceSamples stores price samples for a pool.
func (k Keeper) setPriceSamples(ctx context.Context, poolID uint64, samples []PriceSample) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(samples)
	_ = kvStore.Set(types.PerpsTWAPKey(poolID), bz)
}

// ---------------------------------------------------------------------------
// TWAP Calculation
// ---------------------------------------------------------------------------

// CalculateTWAP computes the Time-Weighted Average Price over the last windowBlocks.
// Returns zero if no samples are available.
func (k Keeper) CalculateTWAP(ctx context.Context, poolID uint64, windowBlocks int64) math.LegacyDec {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentHeight := sdkCtx.BlockHeight()
	cutoff := currentHeight - windowBlocks

	samples := k.getPriceSamples(ctx, poolID)
	if len(samples) == 0 {
		return math.LegacyZeroDec()
	}

	// Filter samples within the window
	sum := math.LegacyZeroDec()
	count := int64(0)
	for _, s := range samples {
		if s.Height >= cutoff {
			sum = sum.Add(s.Price)
			count++
		}
	}

	if count == 0 {
		// No samples in window — use the most recent sample available
		return samples[len(samples)-1].Price
	}

	return sum.Quo(math.LegacyNewDec(count))
}

// ---------------------------------------------------------------------------
// DEX Oracle Price
// ---------------------------------------------------------------------------

// GetDEXOraclePrice returns the TWAP from the DEX pool, using a 100-block window.
// This is the primary oracle-free price source for perps mark price.
func (k Keeper) GetDEXOraclePrice(ctx context.Context, poolID uint64) (math.LegacyDec, error) {
	if k.dexKeeper == nil {
		return math.LegacyZeroDec(), fmt.Errorf("dex keeper not set")
	}

	// First, get the current spot price and record it as a sample
	market, found := k.findMarketByPoolID(ctx, poolID)
	if !found {
		return math.LegacyZeroDec(), fmt.Errorf("no perps market found for pool %d", poolID)
	}

	spotPrice, err := k.dexKeeper.GetSpotPrice(ctx, poolID, market.BaseDenom, market.QuoteDenom)
	if err != nil {
		return math.LegacyZeroDec(), fmt.Errorf("failed to get spot price: %w", err)
	}

	// Note: samples are recorded in BeginBlock via RecordAllPriceSamples.
	// Do NOT record again here to avoid double-counting.

	// Calculate TWAP over the standard window
	twap := k.CalculateTWAP(ctx, poolID, TWAPWindowBlocks)
	if twap.IsZero() {
		// If TWAP is zero (shouldn't happen after recording), fall back to spot
		return spotPrice, nil
	}

	return twap, nil
}

// ---------------------------------------------------------------------------
// Funding Rate Based on TWAP
// ---------------------------------------------------------------------------

// GetFundingRate calculates the funding rate based on the difference between
// the current mark price and the TWAP. A positive rate means mark > TWAP
// (longs pay shorts), negative means shorts pay longs.
func (k Keeper) GetFundingRate(ctx context.Context, poolID uint64) math.LegacyDec {
	if k.dexKeeper == nil {
		return math.LegacyZeroDec()
	}

	market, found := k.findMarketByPoolID(ctx, poolID)
	if !found {
		return math.LegacyZeroDec()
	}

	// Get current spot price (mark price proxy)
	spotPrice, err := k.dexKeeper.GetSpotPrice(ctx, poolID, market.BaseDenom, market.QuoteDenom)
	if err != nil || spotPrice.IsZero() {
		return math.LegacyZeroDec()
	}

	// Get TWAP
	twap := k.CalculateTWAP(ctx, poolID, TWAPWindowBlocks)
	if twap.IsZero() {
		return math.LegacyZeroDec()
	}

	// Funding rate = (markPrice - twap) / twap
	// Positive = longs pay shorts (mark is above fair value)
	// Negative = shorts pay longs (mark is below fair value)
	deviation := spotPrice.Sub(twap)
	fundingRate := deviation.Quo(twap)

	// Clamp to max funding rate
	maxRate := math.LegacyNewDecWithPrec(FundingRateClampBps, 4) // 1%
	if fundingRate.GT(maxRate) {
		fundingRate = maxRate
	}
	if fundingRate.LT(maxRate.Neg()) {
		fundingRate = maxRate.Neg()
	}

	return fundingRate
}

// ---------------------------------------------------------------------------
// BeginBlock Integration — Record Price Samples
// ---------------------------------------------------------------------------

// RecordAllPriceSamples records a price sample for every active perps market.
// Should be called in BeginBlock.
func (k Keeper) RecordAllPriceSamples(ctx context.Context) {
	if k.dexKeeper == nil {
		return
	}

	markets := k.GetAllMarkets(ctx)
	for _, market := range markets {
		if !market.Active {
			continue
		}

		spotPrice, err := k.dexKeeper.GetSpotPrice(ctx, market.PoolID, market.BaseDenom, market.QuoteDenom)
		if err != nil || spotPrice.IsZero() {
			continue
		}

		k.RecordPriceSample(ctx, market.PoolID, spotPrice)
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// findMarketByPoolID finds the perps market associated with a DEX pool ID.
func (k Keeper) findMarketByPoolID(ctx context.Context, poolID uint64) (types.PerpMarket, bool) {
	markets := k.GetAllMarkets(ctx)
	for _, m := range markets {
		if m.PoolID == poolID {
			return m, true
		}
	}
	return types.PerpMarket{}, false
}

// GetTWAPHistory returns the stored price samples for a pool (for debugging/queries).
func (k Keeper) GetTWAPHistory(ctx context.Context, poolID uint64) []PriceSample {
	samples := k.getPriceSamples(ctx, poolID)
	if samples == nil {
		return []PriceSample{}
	}
	return samples
}
