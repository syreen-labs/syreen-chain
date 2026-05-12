package types

import (
	sdkmath "cosmossdk.io/math"
)

// CLPool represents a concentrated liquidity pool
type CLPool struct {
	ID              uint64            `json:"id"`
	DenomA          string            `json:"denom_a"`           // token0 (sorted)
	DenomB          string            `json:"denom_b"`           // token1 (sorted)
	TickSpacing     int64             `json:"tick_spacing"`      // minimum tick distance
	FeeRate         sdkmath.LegacyDec `json:"fee_rate"`          // e.g. 0.003 = 0.3%
	CurrentTick     int64             `json:"current_tick"`      // current tick index
	SqrtPrice       sdkmath.LegacyDec `json:"sqrt_price"`        // sqrt(price) with 18 decimals precision
	TotalLiquidity  sdkmath.LegacyDec `json:"total_liquidity"`   // total in-range liquidity (L)
	FeeGrowthGlobal0 sdkmath.LegacyDec `json:"fee_growth_global_0"` // accumulated fees per unit liquidity for token0
	FeeGrowthGlobal1 sdkmath.LegacyDec `json:"fee_growth_global_1"` // accumulated fees per unit liquidity for token1
}

// CLPosition represents a liquidity provider's position in a CL pool
type CLPosition struct {
	ID                  uint64            `json:"id"`
	Owner               string            `json:"owner"`
	PoolID              uint64            `json:"pool_id"`
	TickLower           int64             `json:"tick_lower"`
	TickUpper           int64             `json:"tick_upper"`
	Liquidity           sdkmath.LegacyDec `json:"liquidity"`
	FeeGrowthInside0Last sdkmath.LegacyDec `json:"fee_growth_inside_0_last"`
	FeeGrowthInside1Last sdkmath.LegacyDec `json:"fee_growth_inside_1_last"`
	TokensOwed0         sdkmath.LegacyDec `json:"tokens_owed_0"` // uncollected fees token0
	TokensOwed1         sdkmath.LegacyDec `json:"tokens_owed_1"` // uncollected fees token1
}

// TickInfo stores liquidity and fee data per tick
type TickInfo struct {
	LiquidityGross    sdkmath.LegacyDec `json:"liquidity_gross"`     // total liquidity referencing this tick
	LiquidityNet      sdkmath.LegacyDec `json:"liquidity_net"`       // net liquidity change when crossing (positive = add, negative = remove)
	FeeGrowthOutside0 sdkmath.LegacyDec `json:"fee_growth_outside_0"`
	FeeGrowthOutside1 sdkmath.LegacyDec `json:"fee_growth_outside_1"`
	Initialized       bool              `json:"initialized"`
}

// GenesisState
type GenesisState struct {
	Pools     []CLPool     `json:"pools"`
	Positions []CLPosition `json:"positions"`
}

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Pools:     []CLPool{},
		Positions: []CLPosition{},
	}
}

func (gs GenesisState) Validate() error { return nil }

// TickToSqrtPrice converts a tick index to sqrt(1.0001^tick).
// Uses deterministic fixed-point arithmetic via repeated squaring.
// sqrtPrice = 1.0001^(tick/2)
func TickToSqrtPrice(tick int64) sdkmath.LegacyDec {
	// base = 1.0001 represented as LegacyDec
	base := sdkmath.LegacyNewDecWithPrec(10001, 4) // 1.0001

	// We compute base^|tick| via repeated squaring, then take sqrt.
	// For negative ticks, invert the result.
	absTick := tick
	negative := false
	if tick < 0 {
		absTick = -tick
		negative = true
	}

	// Compute 1.0001^absTick by repeated squaring
	result := sdkmath.LegacyOneDec()
	power := base
	remaining := absTick
	for remaining > 0 {
		if remaining%2 == 1 {
			result = result.Mul(power)
		}
		power = power.Mul(power)
		remaining /= 2
	}

	if negative {
		result = sdkmath.LegacyOneDec().Quo(result)
	}

	// Now result = 1.0001^tick, we need sqrt(result)
	return decSqrt(result)
}

// decSqrt computes the square root of a LegacyDec using Newton's method.
// Deterministic: only uses LegacyDec arithmetic.
func decSqrt(x sdkmath.LegacyDec) sdkmath.LegacyDec {
	if x.IsZero() || x.IsNegative() {
		return sdkmath.LegacyZeroDec()
	}
	if x.Equal(sdkmath.LegacyOneDec()) {
		return sdkmath.LegacyOneDec()
	}
	two := sdkmath.LegacyNewDec(2)
	// Initial guess: x/2 or 1, whichever is closer
	guess := x.Add(sdkmath.LegacyOneDec()).Quo(two)
	for i := 0; i < 100; i++ { // 100 iterations is more than enough for 18-digit precision
		next := guess.Add(x.Quo(guess)).Quo(two)
		if next.Equal(guess) {
			break
		}
		guess = next
	}
	return guess
}

// SqrtPriceToTick converts a sqrt price to the nearest tick below it.
// Uses deterministic binary search instead of floating-point logarithms.
func SqrtPriceToTick(sqrtPrice sdkmath.LegacyDec) int64 {
	if sqrtPrice.LTE(sdkmath.LegacyZeroDec()) {
		return MinTick
	}

	// Binary search for the tick such that TickToSqrtPrice(tick) <= sqrtPrice < TickToSqrtPrice(tick+1)
	lo := MinTick
	hi := MaxTick

	for lo < hi {
		mid := lo + (hi-lo+1)/2
		midPrice := TickToSqrtPrice(mid)
		if midPrice.LTE(sqrtPrice) {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return lo
}

// MinTick and MaxTick define the boundaries
const (
	MinTick int64 = -887272
	MaxTick int64 = 887272
)

// CalcAmount0Delta calculates the amount of token0 needed for a liquidity change between two sqrt prices.
// amount0 = L * (1/sqrtPriceLower - 1/sqrtPriceUpper)
func CalcAmount0Delta(liquidity, sqrtPriceLower, sqrtPriceUpper sdkmath.LegacyDec) sdkmath.LegacyDec {
	if sqrtPriceLower.GT(sqrtPriceUpper) {
		sqrtPriceLower, sqrtPriceUpper = sqrtPriceUpper, sqrtPriceLower
	}
	if sqrtPriceLower.IsZero() {
		return sdkmath.LegacyZeroDec()
	}
	// L * (sqrtPriceUpper - sqrtPriceLower) / (sqrtPriceLower * sqrtPriceUpper)
	diff := sqrtPriceUpper.Sub(sqrtPriceLower)
	product := sqrtPriceLower.Mul(sqrtPriceUpper)
	if product.IsZero() {
		return sdkmath.LegacyZeroDec()
	}
	return liquidity.Mul(diff).Quo(product)
}

// CalcAmount1Delta calculates the amount of token1 needed for a liquidity change between two sqrt prices.
// amount1 = L * (sqrtPriceUpper - sqrtPriceLower)
func CalcAmount1Delta(liquidity, sqrtPriceLower, sqrtPriceUpper sdkmath.LegacyDec) sdkmath.LegacyDec {
	if sqrtPriceLower.GT(sqrtPriceUpper) {
		sqrtPriceLower, sqrtPriceUpper = sqrtPriceUpper, sqrtPriceLower
	}
	return liquidity.Mul(sqrtPriceUpper.Sub(sqrtPriceLower))
}

// LiquidityFromAmounts calculates the max liquidity that can be provided given amounts and price range.
func LiquidityFromAmounts(sqrtPrice, sqrtPriceLower, sqrtPriceUpper sdkmath.LegacyDec, amount0, amount1 sdkmath.Int) sdkmath.LegacyDec {
	if sqrtPriceLower.GT(sqrtPriceUpper) {
		sqrtPriceLower, sqrtPriceUpper = sqrtPriceUpper, sqrtPriceLower
	}

	amt0Dec := sdkmath.LegacyNewDecFromInt(amount0)
	amt1Dec := sdkmath.LegacyNewDecFromInt(amount1)

	if sqrtPrice.LTE(sqrtPriceLower) {
		// Current price below range — only token0 needed
		// L = amount0 * sqrtPriceLower * sqrtPriceUpper / (sqrtPriceUpper - sqrtPriceLower)
		diff := sqrtPriceUpper.Sub(sqrtPriceLower)
		if diff.IsZero() {
			return sdkmath.LegacyZeroDec()
		}
		return amt0Dec.Mul(sqrtPriceLower).Mul(sqrtPriceUpper).Quo(diff)
	} else if sqrtPrice.GTE(sqrtPriceUpper) {
		// Current price above range — only token1 needed
		// L = amount1 / (sqrtPriceUpper - sqrtPriceLower)
		diff := sqrtPriceUpper.Sub(sqrtPriceLower)
		if diff.IsZero() {
			return sdkmath.LegacyZeroDec()
		}
		return amt1Dec.Quo(diff)
	} else {
		// Current price in range — use minimum of both liquidity calculations
		diff0 := sqrtPriceUpper.Sub(sqrtPrice)
		diff1 := sqrtPrice.Sub(sqrtPriceLower)
		if diff0.IsZero() || diff1.IsZero() {
			return sdkmath.LegacyZeroDec()
		}
		liq0 := amt0Dec.Mul(sqrtPrice).Mul(sqrtPriceUpper).Quo(diff0)
		liq1 := amt1Dec.Quo(diff1)
		if liq0.LT(liq1) {
			return liq0
		}
		return liq1
	}
}
