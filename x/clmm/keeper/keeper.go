package keeper

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"

	"cosmossdk.io/core/store"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/clmm/types"
)

type Keeper struct {
	storeService  store.KVStoreService
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	authority     string
}

func NewKeeper(
	storeService store.KVStoreService,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	authority string,
) Keeper {
	return Keeper{
		storeService:  storeService,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		authority:     authority,
	}
}

// ============================================================
// Pool CRUD
// ============================================================

func (k Keeper) GetPool(ctx context.Context, poolID uint64) (types.CLPool, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.CLPoolKey(poolID))
	if err != nil || bz == nil { return types.CLPool{}, false }
	var pool types.CLPool
	if err := json.Unmarshal(bz, &pool); err != nil { return types.CLPool{}, false }
	return pool, true
}

func (k Keeper) SetPool(ctx context.Context, pool types.CLPool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(pool)
	_ = kvStore.Set(types.CLPoolKey(pool.ID), bz)
}

func (k Keeper) GetAllPools(ctx context.Context) []types.CLPool {
	kvStore := k.storeService.OpenKVStore(ctx)
	iter, err := kvStore.Iterator([]byte(types.CLPoolPrefix), append([]byte(types.CLPoolPrefix), 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
	if err != nil { return nil }
	defer iter.Close()
	var pools []types.CLPool
	for ; iter.Valid(); iter.Next() {
		var pool types.CLPool
		if err := json.Unmarshal(iter.Value(), &pool); err == nil {
			pools = append(pools, pool)
		}
	}
	return pools
}

func (k Keeper) getNextPoolID(ctx context.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.NextCLPoolIDKey))
	if err != nil || bz == nil { return 1 }
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) setNextPoolID(ctx context.Context, id uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	_ = kvStore.Set([]byte(types.NextCLPoolIDKey), bz)
}

// ============================================================
// Position CRUD
// ============================================================

func (k Keeper) GetPosition(ctx context.Context, positionID uint64) (types.CLPosition, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.CLPositionKey(positionID))
	if err != nil || bz == nil { return types.CLPosition{}, false }
	var pos types.CLPosition
	if err := json.Unmarshal(bz, &pos); err != nil { return types.CLPosition{}, false }
	return pos, true
}

func (k Keeper) SetPosition(ctx context.Context, pos types.CLPosition) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(pos)
	_ = kvStore.Set(types.CLPositionKey(pos.ID), bz)
	_ = kvStore.Set(types.PosByOwnerKey(pos.Owner, pos.ID), []byte{1})
	_ = kvStore.Set(types.PosByPoolKey(pos.PoolID, pos.ID), []byte{1})
}

func (k Keeper) DeletePosition(ctx context.Context, pos types.CLPosition) {
	kvStore := k.storeService.OpenKVStore(ctx)
	_ = kvStore.Delete(types.CLPositionKey(pos.ID))
	_ = kvStore.Delete(types.PosByOwnerKey(pos.Owner, pos.ID))
	_ = kvStore.Delete(types.PosByPoolKey(pos.PoolID, pos.ID))
}

func (k Keeper) getNextPositionID(ctx context.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.NextCLPositionIDKey))
	if err != nil || bz == nil { return 1 }
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) setNextPositionID(ctx context.Context, id uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	_ = kvStore.Set([]byte(types.NextCLPositionIDKey), bz)
}

// ============================================================
// Tick CRUD
// ============================================================

func (k Keeper) GetTickInfo(ctx context.Context, poolID uint64, tickIndex int64) types.TickInfo {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.TickInfoKey(poolID, tickIndex))
	if err != nil || bz == nil {
		return types.TickInfo{
			LiquidityGross:    math.LegacyZeroDec(),
			LiquidityNet:      math.LegacyZeroDec(),
			FeeGrowthOutside0: math.LegacyZeroDec(),
			FeeGrowthOutside1: math.LegacyZeroDec(),
			Initialized:       false,
		}
	}
	var tick types.TickInfo
	if err := json.Unmarshal(bz, &tick); err != nil {
		return types.TickInfo{
			LiquidityGross:    math.LegacyZeroDec(),
			LiquidityNet:      math.LegacyZeroDec(),
			FeeGrowthOutside0: math.LegacyZeroDec(),
			FeeGrowthOutside1: math.LegacyZeroDec(),
			Initialized:       false,
		}
	}
	return tick
}

func (k Keeper) SetTickInfo(ctx context.Context, poolID uint64, tickIndex int64, tick types.TickInfo) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(tick)
	_ = kvStore.Set(types.TickInfoKey(poolID, tickIndex), bz)
}

// getFeeGrowthInside calculates the fee growth inside a position's tick range
func (k Keeper) getFeeGrowthInside(ctx context.Context, pool types.CLPool, tickLower, tickUpper int64) (math.LegacyDec, math.LegacyDec) {
	lower := k.GetTickInfo(ctx, pool.ID, tickLower)
	upper := k.GetTickInfo(ctx, pool.ID, tickUpper)

	// Fee growth below lower tick
	var feeGrowthBelow0, feeGrowthBelow1 math.LegacyDec
	if pool.CurrentTick >= tickLower {
		feeGrowthBelow0 = lower.FeeGrowthOutside0
		feeGrowthBelow1 = lower.FeeGrowthOutside1
	} else {
		feeGrowthBelow0 = pool.FeeGrowthGlobal0.Sub(lower.FeeGrowthOutside0)
		feeGrowthBelow1 = pool.FeeGrowthGlobal1.Sub(lower.FeeGrowthOutside1)
	}

	// Fee growth above upper tick
	var feeGrowthAbove0, feeGrowthAbove1 math.LegacyDec
	if pool.CurrentTick < tickUpper {
		feeGrowthAbove0 = upper.FeeGrowthOutside0
		feeGrowthAbove1 = upper.FeeGrowthOutside1
	} else {
		feeGrowthAbove0 = pool.FeeGrowthGlobal0.Sub(upper.FeeGrowthOutside0)
		feeGrowthAbove1 = pool.FeeGrowthGlobal1.Sub(upper.FeeGrowthOutside1)
	}

	feeGrowthInside0 := pool.FeeGrowthGlobal0.Sub(feeGrowthBelow0).Sub(feeGrowthAbove0)
	feeGrowthInside1 := pool.FeeGrowthGlobal1.Sub(feeGrowthBelow1).Sub(feeGrowthAbove1)

	return feeGrowthInside0, feeGrowthInside1
}

// updatePositionFees updates uncollected fees for a position
func (k Keeper) updatePositionFees(ctx context.Context, pool types.CLPool, pos *types.CLPosition) {
	feeGrowthInside0, feeGrowthInside1 := k.getFeeGrowthInside(ctx, pool, pos.TickLower, pos.TickUpper)

	fees0 := pos.Liquidity.Mul(feeGrowthInside0.Sub(pos.FeeGrowthInside0Last))
	fees1 := pos.Liquidity.Mul(feeGrowthInside1.Sub(pos.FeeGrowthInside1Last))

	if fees0.IsPositive() {
		pos.TokensOwed0 = pos.TokensOwed0.Add(fees0)
	}
	if fees1.IsPositive() {
		pos.TokensOwed1 = pos.TokensOwed1.Add(fees1)
	}

	pos.FeeGrowthInside0Last = feeGrowthInside0
	pos.FeeGrowthInside1Last = feeGrowthInside1
}

// ============================================================
// Core Operations
// ============================================================

// CreatePool creates a new CL pool. Only the chain authority may create pools.
func (k Keeper) CreatePool(ctx context.Context, sender, denomA, denomB string, tickSpacing int64, feeRate, initialPrice math.LegacyDec) (uint64, error) {
	// Access control: only chain authority can create CL pools (governance-gated by design).
	if sender != k.authority {
		return 0, types.ErrUnauthorizedPoolCreate
	}

	// Validate initial price: must be positive and non-zero
	if initialPrice.IsNil() || initialPrice.IsZero() || initialPrice.IsNegative() {
		return 0, types.ErrInvalidPrice
	}

	// Validate tick spacing: must be positive
	if tickSpacing <= 0 {
		return 0, types.ErrInvalidTickSpacing
	}

	// Validate fee rate: must be non-negative and capped at 10%
	maxFeeRate := math.LegacyNewDecWithPrec(1, 1) // 0.1 = 10%
	if feeRate.IsNil() || feeRate.IsNegative() || feeRate.GT(maxFeeRate) {
		return 0, types.ErrInvalidFeeRate
	}

	// Validate denoms are different
	if denomA == denomB {
		return 0, types.ErrSameDenom
	}

	// Sort denoms
	if denomA > denomB {
		denomA, denomB = denomB, denomA
		// Invert price since we swapped
		if initialPrice.IsPositive() {
			initialPrice = math.LegacyOneDec().Quo(initialPrice)
		}
	}

	sqrtPrice, err := initialPrice.ApproxSqrt()
	if err != nil {
		return 0, err
	}
	currentTick := types.SqrtPriceToTick(sqrtPrice)

	poolID := k.getNextPoolID(ctx)
	pool := types.CLPool{
		ID:               poolID,
		DenomA:           denomA,
		DenomB:           denomB,
		TickSpacing:      tickSpacing,
		FeeRate:          feeRate,
		CurrentTick:      currentTick,
		SqrtPrice:        sqrtPrice,
		TotalLiquidity:   math.LegacyZeroDec(),
		FeeGrowthGlobal0: math.LegacyZeroDec(),
		FeeGrowthGlobal1: math.LegacyZeroDec(),
	}

	k.SetPool(ctx, pool)
	k.setNextPoolID(ctx, poolID+1)
	return poolID, nil
}

// ExecuteCreatePosition creates a new LP position in a price range
func (k Keeper) ExecuteCreatePosition(ctx context.Context, sender string, poolID uint64, tickLower, tickUpper int64, amount0Desired, amount1Desired, amount0Min, amount1Min math.Int) (uint64, math.Int, math.Int, math.LegacyDec, error) {
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return 0, math.Int{}, math.Int{}, math.LegacyDec{}, types.ErrPoolNotFound
	}

	// Validate tick alignment
	if tickLower%pool.TickSpacing != 0 || tickUpper%pool.TickSpacing != 0 {
		return 0, math.Int{}, math.Int{}, math.LegacyDec{}, types.ErrInvalidTickSpacing
	}
	if tickLower >= tickUpper {
		return 0, math.Int{}, math.Int{}, math.LegacyDec{}, types.ErrInvalidTickRange
	}
	if tickLower < types.MinTick || tickUpper > types.MaxTick {
		return 0, math.Int{}, math.Int{}, math.LegacyDec{}, types.ErrInvalidTickRange
	}

	sqrtPriceLower := types.TickToSqrtPrice(tickLower)
	sqrtPriceUpper := types.TickToSqrtPrice(tickUpper)

	// Validate desired amounts are positive (at least one must be > 0)
	if (amount0Desired.IsZero() || amount0Desired.IsNegative()) && (amount1Desired.IsZero() || amount1Desired.IsNegative()) {
		return 0, math.Int{}, math.Int{}, math.LegacyDec{}, types.ErrInvalidAmount
	}

	// Calculate liquidity from desired amounts
	liquidity := types.LiquidityFromAmounts(pool.SqrtPrice, sqrtPriceLower, sqrtPriceUpper, amount0Desired, amount1Desired)
	if liquidity.IsZero() || liquidity.IsNegative() {
		return 0, math.Int{}, math.Int{}, math.LegacyDec{}, types.ErrZeroLiquidity
	}

	// Guard against liquidity overflow: cap at 10^36 to prevent Dec math overflow
	maxLiquidity := math.LegacyNewDec(1).MulInt(math.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(36), nil)))
	if liquidity.GT(maxLiquidity) {
		return 0, math.Int{}, math.Int{}, math.LegacyDec{}, types.ErrInvalidAmount
	}

	// Calculate actual token amounts needed
	var amount0, amount1 math.LegacyDec
	if pool.SqrtPrice.LTE(sqrtPriceLower) {
		// Below range: only token0
		amount0 = types.CalcAmount0Delta(liquidity, sqrtPriceLower, sqrtPriceUpper)
		amount1 = math.LegacyZeroDec()
	} else if pool.SqrtPrice.GTE(sqrtPriceUpper) {
		// Above range: only token1
		amount0 = math.LegacyZeroDec()
		amount1 = types.CalcAmount1Delta(liquidity, sqrtPriceLower, sqrtPriceUpper)
	} else {
		// In range: both tokens
		amount0 = types.CalcAmount0Delta(liquidity, pool.SqrtPrice, sqrtPriceUpper)
		amount1 = types.CalcAmount1Delta(liquidity, sqrtPriceLower, pool.SqrtPrice)
	}

	amt0Int := amount0.Ceil().TruncateInt()
	amt1Int := amount1.Ceil().TruncateInt()

	// Slippage check
	if amt0Int.LT(amount0Min) {
		return 0, math.Int{}, math.Int{}, math.LegacyDec{}, types.ErrAmountBelowMin
	}
	if amt1Int.LT(amount1Min) {
		return 0, math.Int{}, math.Int{}, math.LegacyDec{}, types.ErrAmountBelowMin
	}

	// Transfer tokens from sender to module
	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return 0, math.Int{}, math.Int{}, math.LegacyDec{}, err
	}

	var coins sdk.Coins
	if amt0Int.IsPositive() {
		coins = coins.Add(sdk.NewCoin(pool.DenomA, amt0Int))
	}
	if amt1Int.IsPositive() {
		coins = coins.Add(sdk.NewCoin(pool.DenomB, amt1Int))
	}
	if coins.IsAllPositive() {
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, coins); err != nil {
			return 0, math.Int{}, math.Int{}, math.LegacyDec{}, err
		}
	}

	// Update tick info
	k.updateTickForLiquidityChange(ctx, pool, tickLower, tickUpper, liquidity, false)

	// Update pool total liquidity if position is in range
	if pool.CurrentTick >= tickLower && pool.CurrentTick < tickUpper {
		pool.TotalLiquidity = pool.TotalLiquidity.Add(liquidity)
	}
	k.SetPool(ctx, pool)

	// Get fee growth inside for the position
	feeGrowthInside0, feeGrowthInside1 := k.getFeeGrowthInside(ctx, pool, tickLower, tickUpper)

	// Create position
	posID := k.getNextPositionID(ctx)
	pos := types.CLPosition{
		ID:                   posID,
		Owner:                sender,
		PoolID:               poolID,
		TickLower:            tickLower,
		TickUpper:            tickUpper,
		Liquidity:            liquidity,
		FeeGrowthInside0Last: feeGrowthInside0,
		FeeGrowthInside1Last: feeGrowthInside1,
		TokensOwed0:          math.LegacyZeroDec(),
		TokensOwed1:          math.LegacyZeroDec(),
	}
	k.SetPosition(ctx, pos)
	k.setNextPositionID(ctx, posID+1)

	return posID, amt0Int, amt1Int, liquidity, nil
}

// AddLiquidityToPosition adds more liquidity to an existing position
func (k Keeper) AddLiquidityToPosition(ctx context.Context, sender string, positionID uint64, amount0Desired, amount1Desired math.Int) (math.Int, math.Int, math.LegacyDec, error) {
	pos, found := k.GetPosition(ctx, positionID)
	if !found {
		return math.Int{}, math.Int{}, math.LegacyDec{}, types.ErrPositionNotFound
	}
	if pos.Owner != sender {
		return math.Int{}, math.Int{}, math.LegacyDec{}, types.ErrUnauthorized
	}

	pool, found := k.GetPool(ctx, pos.PoolID)
	if !found {
		return math.Int{}, math.Int{}, math.LegacyDec{}, types.ErrPoolNotFound
	}

	// Update fees first
	k.updatePositionFees(ctx, pool, &pos)

	sqrtPriceLower := types.TickToSqrtPrice(pos.TickLower)
	sqrtPriceUpper := types.TickToSqrtPrice(pos.TickUpper)

	liquidity := types.LiquidityFromAmounts(pool.SqrtPrice, sqrtPriceLower, sqrtPriceUpper, amount0Desired, amount1Desired)
	if liquidity.IsZero() || liquidity.IsNegative() {
		return math.Int{}, math.Int{}, math.LegacyDec{}, types.ErrZeroLiquidity
	}

	// Guard against liquidity overflow: cap at 10^36
	maxLiquidity := math.LegacyNewDec(1).MulInt(math.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(36), nil)))
	newTotalLiquidity := pos.Liquidity.Add(liquidity)
	if newTotalLiquidity.GT(maxLiquidity) {
		return math.Int{}, math.Int{}, math.LegacyDec{}, types.ErrInvalidAmount
	}

	// Calculate actual token amounts
	var amount0, amount1 math.LegacyDec
	if pool.SqrtPrice.LTE(sqrtPriceLower) {
		amount0 = types.CalcAmount0Delta(liquidity, sqrtPriceLower, sqrtPriceUpper)
		amount1 = math.LegacyZeroDec()
	} else if pool.SqrtPrice.GTE(sqrtPriceUpper) {
		amount0 = math.LegacyZeroDec()
		amount1 = types.CalcAmount1Delta(liquidity, sqrtPriceLower, sqrtPriceUpper)
	} else {
		amount0 = types.CalcAmount0Delta(liquidity, pool.SqrtPrice, sqrtPriceUpper)
		amount1 = types.CalcAmount1Delta(liquidity, sqrtPriceLower, pool.SqrtPrice)
	}

	amt0Int := amount0.Ceil().TruncateInt()
	amt1Int := amount1.Ceil().TruncateInt()

	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return math.Int{}, math.Int{}, math.LegacyDec{}, err
	}

	var coins sdk.Coins
	if amt0Int.IsPositive() {
		coins = coins.Add(sdk.NewCoin(pool.DenomA, amt0Int))
	}
	if amt1Int.IsPositive() {
		coins = coins.Add(sdk.NewCoin(pool.DenomB, amt1Int))
	}
	if coins.IsAllPositive() {
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, coins); err != nil {
			return math.Int{}, math.Int{}, math.LegacyDec{}, err
		}
	}

	// Update ticks
	k.updateTickForLiquidityChange(ctx, pool, pos.TickLower, pos.TickUpper, liquidity, false)

	// Update pool liquidity if in range
	if pool.CurrentTick >= pos.TickLower && pool.CurrentTick < pos.TickUpper {
		pool.TotalLiquidity = pool.TotalLiquidity.Add(liquidity)
	}
	k.SetPool(ctx, pool)

	pos.Liquidity = pos.Liquidity.Add(liquidity)
	k.SetPosition(ctx, pos)

	return amt0Int, amt1Int, liquidity, nil
}

// RemoveLiquidityFromPosition removes liquidity and returns tokens
func (k Keeper) RemoveLiquidityFromPosition(ctx context.Context, sender string, positionID uint64, liquidityAmount math.LegacyDec) (math.Int, math.Int, error) {
	pos, found := k.GetPosition(ctx, positionID)
	if !found {
		return math.Int{}, math.Int{}, types.ErrPositionNotFound
	}
	if pos.Owner != sender {
		return math.Int{}, math.Int{}, types.ErrUnauthorized
	}

	pool, found := k.GetPool(ctx, pos.PoolID)
	if !found {
		return math.Int{}, math.Int{}, types.ErrPoolNotFound
	}

	if liquidityAmount.IsNil() || liquidityAmount.IsZero() || liquidityAmount.IsNegative() {
		return math.Int{}, math.Int{}, types.ErrZeroLiquidity
	}

	if liquidityAmount.GT(pos.Liquidity) {
		liquidityAmount = pos.Liquidity
	}

	// Update fees first
	k.updatePositionFees(ctx, pool, &pos)

	sqrtPriceLower := types.TickToSqrtPrice(pos.TickLower)
	sqrtPriceUpper := types.TickToSqrtPrice(pos.TickUpper)

	// Calculate tokens to return
	var amount0, amount1 math.LegacyDec
	if pool.SqrtPrice.LTE(sqrtPriceLower) {
		amount0 = types.CalcAmount0Delta(liquidityAmount, sqrtPriceLower, sqrtPriceUpper)
		amount1 = math.LegacyZeroDec()
	} else if pool.SqrtPrice.GTE(sqrtPriceUpper) {
		amount0 = math.LegacyZeroDec()
		amount1 = types.CalcAmount1Delta(liquidityAmount, sqrtPriceLower, sqrtPriceUpper)
	} else {
		amount0 = types.CalcAmount0Delta(liquidityAmount, pool.SqrtPrice, sqrtPriceUpper)
		amount1 = types.CalcAmount1Delta(liquidityAmount, sqrtPriceLower, pool.SqrtPrice)
	}

	amt0Int := amount0.TruncateInt()
	amt1Int := amount1.TruncateInt()

	// Transfer tokens from module to sender
	recipientAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return math.Int{}, math.Int{}, err
	}

	var coins sdk.Coins
	if amt0Int.IsPositive() {
		coins = coins.Add(sdk.NewCoin(pool.DenomA, amt0Int))
	}
	if amt1Int.IsPositive() {
		coins = coins.Add(sdk.NewCoin(pool.DenomB, amt1Int))
	}
	if coins.IsAllPositive() {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipientAddr, coins); err != nil {
			return math.Int{}, math.Int{}, err
		}
	}

	// Update ticks (remove liquidity = negate)
	k.updateTickForLiquidityChange(ctx, pool, pos.TickLower, pos.TickUpper, liquidityAmount, true)

	// Update pool liquidity if in range
	if pool.CurrentTick >= pos.TickLower && pool.CurrentTick < pos.TickUpper {
		pool.TotalLiquidity = pool.TotalLiquidity.Sub(liquidityAmount)
		if pool.TotalLiquidity.IsNegative() {
			pool.TotalLiquidity = math.LegacyZeroDec()
		}
	}
	k.SetPool(ctx, pool)

	pos.Liquidity = pos.Liquidity.Sub(liquidityAmount)
	if pos.Liquidity.IsZero() || pos.Liquidity.IsNegative() {
		// Delete position if no liquidity left and no uncollected fees
		if pos.TokensOwed0.IsZero() && pos.TokensOwed1.IsZero() {
			k.DeletePosition(ctx, pos)
		} else {
			pos.Liquidity = math.LegacyZeroDec()
			k.SetPosition(ctx, pos)
		}
	} else {
		k.SetPosition(ctx, pos)
	}

	return amt0Int, amt1Int, nil
}

// CollectPositionFees collects accrued fees from a position
func (k Keeper) CollectPositionFees(ctx context.Context, sender string, positionID uint64) (math.Int, math.Int, error) {
	pos, found := k.GetPosition(ctx, positionID)
	if !found {
		return math.Int{}, math.Int{}, types.ErrPositionNotFound
	}
	if pos.Owner != sender {
		return math.Int{}, math.Int{}, types.ErrUnauthorized
	}

	pool, found := k.GetPool(ctx, pos.PoolID)
	if !found {
		return math.Int{}, math.Int{}, types.ErrPoolNotFound
	}

	// Update fees
	k.updatePositionFees(ctx, pool, &pos)

	amt0Int := pos.TokensOwed0.TruncateInt()
	amt1Int := pos.TokensOwed1.TruncateInt()

	recipientAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return math.Int{}, math.Int{}, err
	}

	var coins sdk.Coins
	if amt0Int.IsPositive() {
		coins = coins.Add(sdk.NewCoin(pool.DenomA, amt0Int))
	}
	if amt1Int.IsPositive() {
		coins = coins.Add(sdk.NewCoin(pool.DenomB, amt1Int))
	}
	if coins.IsAllPositive() {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipientAddr, coins); err != nil {
			return math.Int{}, math.Int{}, err
		}
	}

	pos.TokensOwed0 = math.LegacyZeroDec()
	pos.TokensOwed1 = math.LegacyZeroDec()
	k.SetPosition(ctx, pos)

	return amt0Int, amt1Int, nil
}

// ExecuteSwap performs a swap through the concentrated liquidity pool
func (k Keeper) ExecuteSwap(ctx context.Context, sender string, poolID uint64, denomIn string, amountIn, minAmountOut math.Int) (math.Int, error) {
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return math.Int{}, types.ErrPoolNotFound
	}

	// Validate swap input amount
	if amountIn.IsNil() || amountIn.IsZero() || !amountIn.IsPositive() {
		return math.Int{}, types.ErrInvalidAmount
	}

	// Determine swap direction: zeroForOne = selling token0 for token1
	var zeroForOne bool
	if denomIn == pool.DenomA {
		zeroForOne = true
	} else if denomIn == pool.DenomB {
		zeroForOne = false
	} else {
		return math.Int{}, types.ErrInvalidDenom
	}

	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return math.Int{}, err
	}

	// Send input tokens to module
	inCoins := sdk.NewCoins(sdk.NewCoin(denomIn, amountIn))
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, inCoins); err != nil {
		return math.Int{}, err
	}

	amountRemaining := math.LegacyNewDecFromInt(amountIn)
	totalAmountOut := math.LegacyZeroDec()

	// Iterative swap across ticks
	maxIterations := 1000
	for i := 0; i < maxIterations && amountRemaining.IsPositive(); i++ {
		// Find the next initialized tick
		nextTick := k.findNextInitializedTick(ctx, pool, zeroForOne)
		nextSqrtPrice := types.TickToSqrtPrice(nextTick)

		// Determine the target sqrt price for this step
		targetSqrtPrice := nextSqrtPrice
		if zeroForOne {
			// Price decreases when selling token0
			if targetSqrtPrice.LT(pool.SqrtPrice) {
				// nextTick is our target
			} else {
				// No more ticks in this direction
				break
			}
		} else {
			// Price increases when selling token1
			if targetSqrtPrice.GT(pool.SqrtPrice) {
				// nextTick is our target
			} else {
				break
			}
		}

		if pool.TotalLiquidity.IsZero() || pool.TotalLiquidity.IsNegative() {
			// No liquidity — try to cross to next tick
			pool.SqrtPrice = targetSqrtPrice
			pool.CurrentTick = nextTick
			k.crossTick(ctx, &pool, nextTick, zeroForOne)
			continue
		}

		// Calculate how much input is needed to reach the target price
		amountInForStep, amountOutForStep := k.computeSwapStep(pool.SqrtPrice, targetSqrtPrice, pool.TotalLiquidity, amountRemaining, pool.FeeRate, zeroForOne)

		if amountInForStep.IsZero() && amountOutForStep.IsZero() {
			// Can't make progress, cross the tick
			pool.SqrtPrice = targetSqrtPrice
			pool.CurrentTick = nextTick
			k.crossTick(ctx, &pool, nextTick, zeroForOne)
			continue
		}

		// Calculate fee
		fee := amountInForStep.Mul(pool.FeeRate)
		amountInWithFee := amountInForStep.Add(fee)

		if amountInWithFee.GT(amountRemaining) {
			// Partial fill: recalculate with available amount
			amountInForStep = amountRemaining.Quo(math.LegacyOneDec().Add(pool.FeeRate))
			fee = amountRemaining.Sub(amountInForStep)

			// Recalculate output and new sqrt price
			amountOutForStep, pool.SqrtPrice = k.computePartialStep(pool.SqrtPrice, pool.TotalLiquidity, amountInForStep, zeroForOne)
			pool.CurrentTick = types.SqrtPriceToTick(pool.SqrtPrice)

			// Update fee growth
			k.updateFeeGrowth(&pool, fee, zeroForOne)

			amountRemaining = math.LegacyZeroDec()
			totalAmountOut = totalAmountOut.Add(amountOutForStep)
		} else {
			amountRemaining = amountRemaining.Sub(amountInWithFee)
			totalAmountOut = totalAmountOut.Add(amountOutForStep)

			// Update fee growth
			k.updateFeeGrowth(&pool, fee, zeroForOne)

			// Move to the target price and cross the tick
			pool.SqrtPrice = targetSqrtPrice
			pool.CurrentTick = nextTick
			k.crossTick(ctx, &pool, nextTick, zeroForOne)
		}
	}

	amountOutInt := totalAmountOut.TruncateInt()

	if amountOutInt.LT(minAmountOut) {
		// Refund input — revert the swap
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, inCoins); err != nil {
			return math.Int{}, err
		}
		return math.Int{}, types.ErrSlippageExceeded
	}

	// Send output tokens to sender
	var denomOut string
	if zeroForOne {
		denomOut = pool.DenomB
	} else {
		denomOut = pool.DenomA
	}
	if amountOutInt.IsPositive() {
		outCoins := sdk.NewCoins(sdk.NewCoin(denomOut, amountOutInt))
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, outCoins); err != nil {
			return math.Int{}, err
		}
	}

	// Refund unused input
	amountUsed := math.LegacyNewDecFromInt(amountIn).Sub(amountRemaining).Ceil().TruncateInt()
	refund := amountIn.Sub(amountUsed)
	if refund.IsPositive() {
		refundCoins := sdk.NewCoins(sdk.NewCoin(denomIn, refund))
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, refundCoins); err != nil {
			return math.Int{}, err
		}
	}

	k.SetPool(ctx, pool)

	return amountOutInt, nil
}

// ============================================================
// Internal helpers
// ============================================================

func (k Keeper) updateTickForLiquidityChange(ctx context.Context, pool types.CLPool, tickLower, tickUpper int64, liquidity math.LegacyDec, isRemove bool) {
	lowerTick := k.GetTickInfo(ctx, pool.ID, tickLower)
	upperTick := k.GetTickInfo(ctx, pool.ID, tickUpper)

	if isRemove {
		lowerTick.LiquidityGross = lowerTick.LiquidityGross.Sub(liquidity)
		lowerTick.LiquidityNet = lowerTick.LiquidityNet.Sub(liquidity)
		upperTick.LiquidityGross = upperTick.LiquidityGross.Sub(liquidity)
		upperTick.LiquidityNet = upperTick.LiquidityNet.Add(liquidity) // upper tick has negative net on add
	} else {
		// Initialize fee growth outside for new ticks
		if !lowerTick.Initialized {
			if pool.CurrentTick >= tickLower {
				lowerTick.FeeGrowthOutside0 = pool.FeeGrowthGlobal0
				lowerTick.FeeGrowthOutside1 = pool.FeeGrowthGlobal1
			}
			lowerTick.Initialized = true
		}
		if !upperTick.Initialized {
			if pool.CurrentTick >= tickUpper {
				upperTick.FeeGrowthOutside0 = pool.FeeGrowthGlobal0
				upperTick.FeeGrowthOutside1 = pool.FeeGrowthGlobal1
			}
			upperTick.Initialized = true
		}

		lowerTick.LiquidityGross = lowerTick.LiquidityGross.Add(liquidity)
		lowerTick.LiquidityNet = lowerTick.LiquidityNet.Add(liquidity)
		upperTick.LiquidityGross = upperTick.LiquidityGross.Add(liquidity)
		upperTick.LiquidityNet = upperTick.LiquidityNet.Sub(liquidity) // crossing upper tick removes liquidity
	}

	k.SetTickInfo(ctx, pool.ID, tickLower, lowerTick)
	k.SetTickInfo(ctx, pool.ID, tickUpper, upperTick)
}

// maxTickSearchIterations caps the O(N) tick scan to prevent DoS via sparse liquidity.
const maxTickSearchIterations = 10000

func (k Keeper) findNextInitializedTick(ctx context.Context, pool types.CLPool, zeroForOne bool) int64 {
	// Search for the next initialized tick in the direction of the swap.
	// For simplicity, iterate through ticks at tick spacing.
	currentTick := pool.CurrentTick
	spacing := pool.TickSpacing
	if spacing <= 0 {
		// Safety: invalid tick spacing, return boundary to halt swap progress
		if zeroForOne {
			return types.MinTick
		}
		return types.MaxTick
	}

	if zeroForOne {
		// Search downward (price decreasing)
		// Snap current tick down to spacing using integer floor division
		// (Go's % truncates toward zero, which rounds negative ticks UP).
		q := currentTick / spacing
		if currentTick%spacing != 0 && currentTick < 0 {
			q--
		}
		searchTick := q * spacing
		if searchTick == currentTick {
			searchTick -= spacing
		}
		iterations := 0
		for t := searchTick; t >= types.MinTick; t -= spacing {
			iterations++
			if iterations > maxTickSearchIterations {
				// Return the boundary — swap will use MinTick as the limit
				return types.MinTick
			}
			tick := k.GetTickInfo(ctx, pool.ID, t)
			if tick.Initialized {
				return t
			}
		}
		return types.MinTick
	}

	// Search upward (price increasing)
	// Snap current tick down to spacing using integer floor division
	// (Go's % truncates toward zero, which rounds negative ticks UP), then step up.
	q := currentTick / spacing
	if currentTick%spacing != 0 && currentTick < 0 {
		q--
	}
	searchTick := q*spacing + spacing
	iterations := 0
	for t := searchTick; t <= types.MaxTick; t += spacing {
		iterations++
		if iterations > maxTickSearchIterations {
			// Return the boundary — swap will use MaxTick as the limit
			return types.MaxTick
		}
		tick := k.GetTickInfo(ctx, pool.ID, t)
		if tick.Initialized {
			return t
		}
	}
	return types.MaxTick
}

func (k Keeper) crossTick(ctx context.Context, pool *types.CLPool, tickIndex int64, zeroForOne bool) {
	tick := k.GetTickInfo(ctx, pool.ID, tickIndex)
	if !tick.Initialized {
		return
	}

	// Flip fee growth outside
	tick.FeeGrowthOutside0 = pool.FeeGrowthGlobal0.Sub(tick.FeeGrowthOutside0)
	tick.FeeGrowthOutside1 = pool.FeeGrowthGlobal1.Sub(tick.FeeGrowthOutside1)

	// Update liquidity
	if zeroForOne {
		// Moving left: subtract net liquidity (because net is positive at lower boundary)
		pool.TotalLiquidity = pool.TotalLiquidity.Sub(tick.LiquidityNet)
		if pool.TotalLiquidity.IsNegative() {
			pool.TotalLiquidity = math.LegacyZeroDec()
		}
		pool.CurrentTick = tickIndex - 1
	} else {
		// Moving right: add net liquidity
		pool.TotalLiquidity = pool.TotalLiquidity.Add(tick.LiquidityNet)
		if pool.TotalLiquidity.IsNegative() {
			pool.TotalLiquidity = math.LegacyZeroDec()
		}
		pool.CurrentTick = tickIndex
	}

	k.SetTickInfo(ctx, pool.ID, tickIndex, tick)
}

func (k Keeper) updateFeeGrowth(pool *types.CLPool, feeAmount math.LegacyDec, zeroForOne bool) {
	if pool.TotalLiquidity.IsZero() || pool.TotalLiquidity.IsNegative() {
		return
	}
	feePerLiquidity := feeAmount.Quo(pool.TotalLiquidity)
	if zeroForOne {
		pool.FeeGrowthGlobal0 = pool.FeeGrowthGlobal0.Add(feePerLiquidity)
	} else {
		pool.FeeGrowthGlobal1 = pool.FeeGrowthGlobal1.Add(feePerLiquidity)
	}
}

// computeSwapStep calculates the input/output amounts for a swap step to reach the target price
func (k Keeper) computeSwapStep(sqrtPriceCurrent, sqrtPriceTarget, liquidity, amountRemaining, feeRate math.LegacyDec, zeroForOne bool) (amountIn, amountOut math.LegacyDec) {
	if zeroForOne {
		// Selling token0: price decreases, get token1
		// amountIn (token0) = L * (1/sqrtTarget - 1/sqrtCurrent) = L * (sqrtCurrent - sqrtTarget) / (sqrtTarget * sqrtCurrent)
		amountIn = types.CalcAmount0Delta(liquidity, sqrtPriceTarget, sqrtPriceCurrent)
		amountOut = types.CalcAmount1Delta(liquidity, sqrtPriceTarget, sqrtPriceCurrent)
	} else {
		// Selling token1: price increases, get token0
		// amountIn (token1) = L * (sqrtTarget - sqrtCurrent)
		amountIn = types.CalcAmount1Delta(liquidity, sqrtPriceCurrent, sqrtPriceTarget)
		amountOut = types.CalcAmount0Delta(liquidity, sqrtPriceCurrent, sqrtPriceTarget)
	}
	return amountIn, amountOut
}

// computePartialStep calculates the output for a partial swap step (when input is limited)
func (k Keeper) computePartialStep(sqrtPriceCurrent, liquidity, amountIn math.LegacyDec, zeroForOne bool) (amountOut math.LegacyDec, newSqrtPrice math.LegacyDec) {
	if liquidity.IsZero() {
		return math.LegacyZeroDec(), sqrtPriceCurrent
	}

	if zeroForOne {
		// token0 in: new_sqrt_price = L * sqrtPrice / (L + amountIn * sqrtPrice)
		denominator := liquidity.Add(amountIn.Mul(sqrtPriceCurrent))
		if denominator.IsZero() {
			return math.LegacyZeroDec(), sqrtPriceCurrent
		}
		newSqrtPrice = liquidity.Mul(sqrtPriceCurrent).Quo(denominator)
		amountOut = types.CalcAmount1Delta(liquidity, newSqrtPrice, sqrtPriceCurrent)
	} else {
		// token1 in: new_sqrt_price = sqrtPrice + amountIn / L
		newSqrtPrice = sqrtPriceCurrent.Add(amountIn.Quo(liquidity))
		amountOut = types.CalcAmount0Delta(liquidity, sqrtPriceCurrent, newSqrtPrice)
	}
	return amountOut, newSqrtPrice
}

// Ensure Keeper implements MsgServer (checked in msg_server.go)
var _ = fmt.Sprintf // suppress unused import

// GetAllPositions returns all positions (for genesis export)
func (k Keeper) GetAllPositions(ctx context.Context) []types.CLPosition {
	kvStore := k.storeService.OpenKVStore(ctx)
	iter, err := kvStore.Iterator([]byte(types.CLPositionPrefix), append([]byte(types.CLPositionPrefix), 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
	if err != nil { return nil }
	defer iter.Close()
	var positions []types.CLPosition
	for ; iter.Valid(); iter.Next() {
		var pos types.CLPosition
		if err := json.Unmarshal(iter.Value(), &pos); err == nil {
			positions = append(positions, pos)
		}
	}
	return positions
}
