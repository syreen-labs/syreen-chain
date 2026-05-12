package keeper

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"sort"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/dex/types"
)

// ---------------------------------------------------------------------------
// Dynamic Fee System (Programmable Fee Tiers)
//
// Pool creators can enable volatility-based dynamic fees. The system:
// 1. Tracks price samples per pool over the last 100 blocks
// 2. Calculates price volatility as standard deviation of price changes
// 3. Dynamic fee = BaseFee + (Volatility * VolatilityMultiplier), capped at MaxFee
// 4. Falls back to regular SwapFee when VolatilityFeeEnabled is false
// ---------------------------------------------------------------------------

const (
	// PriceSamplePrefix stores price samples: price_sample/<poolID>/<height>
	PriceSamplePrefix = "price_sample/"

	// VolatilityWindowBlocks is the number of blocks used to calculate volatility
	VolatilityWindowBlocks = int64(100)
)

// PriceSample records a price observation at a specific block height.
type PriceSample struct {
	Price  math.LegacyDec `json:"price"`
	Height int64          `json:"height"`
}

// priceSampleKey returns the store key for a price sample.
// Format: price_sample/<poolID_8bytes><height_8bytes>
func priceSampleKey(poolID uint64, height int64) []byte {
	key := make([]byte, 16)
	binary.BigEndian.PutUint64(key[:8], poolID)
	binary.BigEndian.PutUint64(key[8:], uint64(height))
	return append([]byte(PriceSamplePrefix), key...)
}

// priceSamplePoolPrefix returns the prefix for all price samples of a pool.
func priceSamplePoolPrefix(poolID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	return append([]byte(PriceSamplePrefix), bz...)
}

// RecordPriceSample stores the current price for a pool at the given height.
// Called during swaps to maintain the price history for volatility calculation.
func (k Keeper) RecordPriceSample(ctx context.Context, poolID uint64, price math.LegacyDec) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	height := sdkCtx.BlockHeight()

	sample := PriceSample{
		Price:  price,
		Height: height,
	}

	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(sample)
	if err != nil {
		return
	}
	kvStore.Set(priceSampleKey(poolID, height), bz)

	// Prune old samples beyond the volatility window
	k.pruneOldPriceSamples(ctx, poolID, height-VolatilityWindowBlocks*2)
}

// pruneOldPriceSamples removes price samples older than the cutoff height.
func (k Keeper) pruneOldPriceSamples(ctx context.Context, poolID uint64, cutoffHeight int64) {
	if cutoffHeight <= 0 {
		return
	}

	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := priceSamplePoolPrefix(poolID)
	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return
	}
	defer iter.Close()

	var keysToDelete [][]byte
	for ; iter.Valid(); iter.Next() {
		// Extract height from key (last 8 bytes after prefix + poolID)
		keyBytes := iter.Key()
		if len(keyBytes) < len(prefix)+8 {
			continue
		}
		heightBytes := keyBytes[len(keyBytes)-8:]
		sampleHeight := int64(binary.BigEndian.Uint64(heightBytes))

		if sampleHeight < cutoffHeight {
			keysToDelete = append(keysToDelete, append([]byte(nil), keyBytes...))
		}
	}

	for _, key := range keysToDelete {
		kvStore.Delete(key)
	}
}

// GetPriceSamples returns recent price samples for a pool within the volatility window.
func (k Keeper) GetPriceSamples(ctx context.Context, poolID uint64) []PriceSample {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	minHeight := sdkCtx.BlockHeight() - VolatilityWindowBlocks
	if minHeight < 0 {
		minHeight = 0
	}

	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := priceSamplePoolPrefix(poolID)
	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var samples []PriceSample
	for ; iter.Valid(); iter.Next() {
		var sample PriceSample
		if err := json.Unmarshal(iter.Value(), &sample); err != nil {
			continue
		}
		if sample.Height >= minHeight {
			samples = append(samples, sample)
		}
	}

	// Sort by height for determinism
	sort.Slice(samples, func(i, j int) bool {
		return samples[i].Height < samples[j].Height
	})

	return samples
}

// CalculateVolatility computes the price volatility for a pool as the standard
// deviation of block-to-block price changes over the last 100 blocks.
// Returns zero if there are fewer than 2 price samples.
func (k Keeper) CalculateVolatility(ctx context.Context, poolID uint64) math.LegacyDec {
	samples := k.GetPriceSamples(ctx, poolID)
	if len(samples) < 2 {
		return math.LegacyZeroDec()
	}

	// Calculate price changes (returns)
	var returns []math.LegacyDec
	for i := 1; i < len(samples); i++ {
		if samples[i-1].Price.IsZero() {
			continue
		}
		// Percentage change: (price[i] - price[i-1]) / price[i-1]
		change := samples[i].Price.Sub(samples[i-1].Price).Quo(samples[i-1].Price)
		returns = append(returns, change)
	}

	if len(returns) < 2 {
		return math.LegacyZeroDec()
	}

	// Calculate mean
	sum := math.LegacyZeroDec()
	for _, r := range returns {
		sum = sum.Add(r)
	}
	mean := sum.Quo(math.LegacyNewDec(int64(len(returns))))

	// Calculate variance = sum((r - mean)^2) / n
	varianceSum := math.LegacyZeroDec()
	for _, r := range returns {
		diff := r.Sub(mean)
		varianceSum = varianceSum.Add(diff.Mul(diff))
	}
	variance := varianceSum.Quo(math.LegacyNewDec(int64(len(returns))))

	// Standard deviation = sqrt(variance)
	// Use Newton's method for square root of a Dec
	stdDev := sqrtDec(variance)

	return stdDev
}

// sqrtDec computes the square root of a LegacyDec using Newton's method.
func sqrtDec(x math.LegacyDec) math.LegacyDec {
	if x.IsZero() || x.IsNegative() {
		return math.LegacyZeroDec()
	}

	// Initial guess
	z := x
	half := math.LegacyNewDecWithPrec(5, 1) // 0.5

	// Newton iterations: z = (z + x/z) / 2
	for i := 0; i < 50; i++ {
		if z.IsZero() {
			return math.LegacyZeroDec()
		}
		next := z.Add(x.Quo(z)).Mul(half)
		// Converged when diff < 1e-18
		diff := next.Sub(z)
		if diff.IsNegative() {
			diff = diff.Neg()
		}
		z = next
		if diff.LT(math.LegacyNewDecWithPrec(1, 18)) {
			break
		}
	}

	return z
}

// GetEffectiveFee returns the current fee for a pool, accounting for dynamic fees.
// If VolatilityFeeEnabled is false, returns the regular SwapFee (backward compatible).
// Otherwise: fee = BaseFee + (Volatility * VolatilityMultiplier), capped at MaxFee.
func (k Keeper) GetEffectiveFee(ctx context.Context, poolID uint64) math.LegacyDec {
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return math.LegacyZeroDec()
	}

	// Backward compatible: if dynamic fees not enabled, use regular SwapFee
	if !pool.VolatilityFeeEnabled {
		return pool.SwapFee
	}

	// Calculate dynamic fee
	baseFee := pool.BaseFee
	if baseFee.IsNil() || baseFee.IsNegative() {
		baseFee = pool.SwapFee // fallback to SwapFee
	}

	maxFee := pool.MaxFee
	if maxFee.IsNil() || maxFee.IsNegative() {
		maxFee = math.LegacyNewDecWithPrec(5, 2) // 5% default max
	}

	multiplier := pool.VolatilityMultiplier
	if multiplier.IsNil() || multiplier.IsNegative() {
		multiplier = math.LegacyOneDec()
	}

	volatility := k.CalculateVolatility(ctx, poolID)

	// Dynamic fee = BaseFee + (Volatility * Multiplier)
	dynamicFee := baseFee.Add(volatility.Mul(multiplier))

	// Cap at MaxFee
	if dynamicFee.GT(maxFee) {
		dynamicFee = maxFee
	}

	// Floor at BaseFee
	if dynamicFee.LT(baseFee) {
		dynamicFee = baseFee
	}

	return dynamicFee
}

// SetPoolFeeConfig allows the pool creator to configure dynamic fee settings.
func (k Keeper) SetPoolFeeConfig(ctx context.Context, sender string, poolID uint64, enabled bool, baseFee, maxFee, multiplier math.LegacyDec) error {
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return types.ErrPoolNotFound
	}

	// Only the pool creator can configure fees
	if pool.Creator != sender {
		return types.ErrNotPoolCreator
	}

	// Validate fee parameters
	if !baseFee.IsNil() && baseFee.IsNegative() {
		return types.ErrInvalidAmount
	}
	if !maxFee.IsNil() && maxFee.IsNegative() {
		return types.ErrInvalidAmount
	}
	if !multiplier.IsNil() && multiplier.IsNegative() {
		return types.ErrInvalidAmount
	}
	// Cap multiplier at 100 to prevent unbounded fee scaling during volatility.
	// Even at 100x, the MaxFee cap (10%) still applies as the ultimate ceiling.
	maxAllowedMultiplier := math.LegacyNewDec(100)
	if !multiplier.IsNil() && multiplier.GT(maxAllowedMultiplier) {
		return types.ErrInvalidFeeConfig
	}

	// Cap MaxFee at 10% (0.1) to prevent pool creators from stealing swap value
	maxAllowedFee := math.LegacyNewDecWithPrec(1, 1) // 0.1 = 10%
	if !maxFee.IsNil() && maxFee.GT(maxAllowedFee) {
		return types.ErrInvalidFeeConfig
	}
	// Cap BaseFee at 5% (0.05)
	maxAllowedBase := math.LegacyNewDecWithPrec(5, 2) // 0.05 = 5%
	if !baseFee.IsNil() && baseFee.GT(maxAllowedBase) {
		return types.ErrInvalidFeeConfig
	}

	// Ensure maxFee >= baseFee if both set
	if !baseFee.IsNil() && !maxFee.IsNil() && maxFee.LT(baseFee) {
		return types.ErrInvalidFeeConfig
	}

	pool.VolatilityFeeEnabled = enabled
	if !baseFee.IsNil() {
		pool.BaseFee = baseFee
	}
	if !maxFee.IsNil() {
		pool.MaxFee = maxFee
	}
	if !multiplier.IsNil() {
		pool.VolatilityMultiplier = multiplier
	}

	k.SetPool(ctx, pool)

	k.Logger(ctx).Info("pool fee config updated",
		"pool_id", poolID,
		"enabled", enabled,
		"base_fee", baseFee.String(),
		"max_fee", maxFee.String(),
		"multiplier", multiplier.String(),
		"sender", sender,
	)

	return nil
}
