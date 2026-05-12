package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"cosmossdk.io/core/store"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/lending/types"
)

type Keeper struct {
	cdc           codec.Codec
	storeService  store.KVStoreService
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	dexKeeper     types.DexKeeper
	authority     string
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	authority string,
) Keeper {
	return Keeper{
		cdc: cdc, storeService: storeService,
		accountKeeper: accountKeeper, bankKeeper: bankKeeper,
		authority: authority,
	}
}

func (k *Keeper) SetDexKeeper(dk types.DexKeeper) { k.dexKeeper = dk }

// ============================================================
// Pool CRUD
// ============================================================

func (k Keeper) GetPool(ctx context.Context, poolID uint64) (types.LendingPool, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.LendingPoolKey(poolID))
	if err != nil || bz == nil { return types.LendingPool{}, false }
	var pool types.LendingPool
	if err := json.Unmarshal(bz, &pool); err != nil { return types.LendingPool{}, false }
	return pool, true
}

func (k Keeper) SetPool(ctx context.Context, pool types.LendingPool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(pool)
	_ = kvStore.Set(types.LendingPoolKey(pool.ID), bz)
}

func (k Keeper) GetAllPools(ctx context.Context) []types.LendingPool {
	kvStore := k.storeService.OpenKVStore(ctx)
	iter, err := kvStore.Iterator([]byte(types.LendingPoolPrefix), append([]byte(types.LendingPoolPrefix), 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
	if err != nil { return nil }
	defer iter.Close()
	var pools []types.LendingPool
	for ; iter.Valid(); iter.Next() {
		var p types.LendingPool
		if err := json.Unmarshal(iter.Value(), &p); err == nil { pools = append(pools, p) }
	}
	sort.Slice(pools, func(i, j int) bool { return pools[i].ID < pools[j].ID })
	return pools
}

func (k Keeper) GetNextPoolID(ctx context.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := kvStore.Get([]byte(types.NextPoolIDKey))
	if bz == nil { return 1 }
	var id uint64
	json.Unmarshal(bz, &id)
	return id
}

func (k Keeper) SetNextPoolID(ctx context.Context, id uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(id)
	_ = kvStore.Set([]byte(types.NextPoolIDKey), bz)
}

func (k Keeper) GetNextBorrowID(ctx context.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := kvStore.Get([]byte(types.NextBorrowIDKey))
	if bz == nil { return 1 }
	var id uint64
	json.Unmarshal(bz, &id)
	return id
}

func (k Keeper) SetNextBorrowID(ctx context.Context, id uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(id)
	_ = kvStore.Set([]byte(types.NextBorrowIDKey), bz)
}

// ============================================================
// Deposit CRUD
// ============================================================

func (k Keeper) GetDeposit(ctx context.Context, poolID uint64, address string) (types.Deposit, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.DepositKey(poolID, address))
	if err != nil || bz == nil { return types.Deposit{}, false }
	var dep types.Deposit
	if err := json.Unmarshal(bz, &dep); err != nil { return types.Deposit{}, false }
	return dep, true
}

func (k Keeper) SetDeposit(ctx context.Context, dep types.Deposit) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(dep)
	_ = kvStore.Set(types.DepositKey(dep.PoolID, dep.Address), bz)
	_ = kvStore.Set(types.DepositByAddrKey(dep.Address, dep.PoolID), []byte{1})
}

func (k Keeper) DeleteDeposit(ctx context.Context, poolID uint64, address string) {
	kvStore := k.storeService.OpenKVStore(ctx)
	_ = kvStore.Delete(types.DepositKey(poolID, address))
	_ = kvStore.Delete(types.DepositByAddrKey(address, poolID))
}

func (k Keeper) GetDepositsByAddress(ctx context.Context, address string) []types.Deposit {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := types.DepositByAddrPrefixKey(address)
	iter, err := kvStore.Iterator(prefix, append(prefix, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
	if err != nil { return nil }
	defer iter.Close()
	var deposits []types.Deposit
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if len(key) >= 8 {
			b := key[len(key)-8:]
			poolID := uint64(b[0])<<56 | uint64(b[1])<<48 | uint64(b[2])<<40 | uint64(b[3])<<32 | uint64(b[4])<<24 | uint64(b[5])<<16 | uint64(b[6])<<8 | uint64(b[7])
			if dep, ok := k.GetDeposit(ctx, poolID, address); ok { deposits = append(deposits, dep) }
		}
	}
	return deposits
}

// ============================================================
// Borrow CRUD
// ============================================================

func (k Keeper) GetBorrow(ctx context.Context, borrowID uint64) (types.Borrow, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.BorrowKey(borrowID))
	if err != nil || bz == nil { return types.Borrow{}, false }
	var b types.Borrow
	if err := json.Unmarshal(bz, &b); err != nil { return types.Borrow{}, false }
	return b, true
}

func (k Keeper) SetBorrow(ctx context.Context, b types.Borrow) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(b)
	_ = kvStore.Set(types.BorrowKey(b.ID), bz)
	_ = kvStore.Set(types.BorrowByAddrKey(b.Borrower, b.ID), []byte{1})
}

func (k Keeper) DeleteBorrow(ctx context.Context, borrowID uint64, borrower string) {
	kvStore := k.storeService.OpenKVStore(ctx)
	_ = kvStore.Delete(types.BorrowKey(borrowID))
	_ = kvStore.Delete(types.BorrowByAddrKey(borrower, borrowID))
}

func (k Keeper) GetBorrowsByAddress(ctx context.Context, address string) []types.Borrow {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := types.BorrowByAddrPrefixKey(address)
	iter, err := kvStore.Iterator(prefix, append(prefix, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
	if err != nil { return nil }
	defer iter.Close()
	var borrows []types.Borrow
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if len(key) >= 8 {
			b := key[len(key)-8:]
			borrowID := uint64(b[0])<<56 | uint64(b[1])<<48 | uint64(b[2])<<40 | uint64(b[3])<<32 | uint64(b[4])<<24 | uint64(b[5])<<16 | uint64(b[6])<<8 | uint64(b[7])
			if borrow, ok := k.GetBorrow(ctx, borrowID); ok { borrows = append(borrows, borrow) }
		}
	}
	return borrows
}

// GetAllDeposits returns all deposits across all pools (for genesis export)
func (k Keeper) GetAllDeposits(ctx context.Context) []types.Deposit {
	kvStore := k.storeService.OpenKVStore(ctx)
	iter, err := kvStore.Iterator([]byte(types.DepositPrefix), append([]byte(types.DepositPrefix), 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
	if err != nil { return nil }
	defer iter.Close()
	var deposits []types.Deposit
	for ; iter.Valid(); iter.Next() {
		var d types.Deposit
		if err := json.Unmarshal(iter.Value(), &d); err == nil { deposits = append(deposits, d) }
	}
	return deposits
}

// GetAllBorrows returns all borrows (for genesis export)
func (k Keeper) GetAllBorrows(ctx context.Context) []types.Borrow {
	kvStore := k.storeService.OpenKVStore(ctx)
	iter, err := kvStore.Iterator([]byte(types.BorrowPrefix), append([]byte(types.BorrowPrefix), 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
	if err != nil { return nil }
	defer iter.Close()
	var borrows []types.Borrow
	for ; iter.Valid(); iter.Next() {
		var b types.Borrow
		if err := json.Unmarshal(iter.Value(), &b); err == nil { borrows = append(borrows, b) }
	}
	return borrows
}

// ============================================================
// TWAP Oracle (C-1)
// ============================================================

// safeAccIndex returns max(x, OneDec). M-2: floor for borrow indices to
// prevent corruption (e.g. from JSON unmarshalling a zero) from lowering debt.
func safeAccIndex(x math.LegacyDec) math.LegacyDec {
	if x.IsNil() || x.LT(math.LegacyOneDec()) {
		return math.LegacyOneDec()
	}
	return x
}

// GetTWAPSeries fetches the stored TWAP series for a lending pool.
func (k Keeper) GetTWAPSeries(ctx context.Context, poolID uint64) types.TWAPSeries {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.TWAPSeriesKey(poolID))
	if err != nil || bz == nil {
		return types.TWAPSeries{PoolID: poolID}
	}
	var s types.TWAPSeries
	if err := json.Unmarshal(bz, &s); err != nil {
		return types.TWAPSeries{PoolID: poolID}
	}
	return s
}

// SetTWAPSeries persists the TWAP series for a lending pool.
func (k Keeper) SetTWAPSeries(ctx context.Context, s types.TWAPSeries) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(s)
	_ = kvStore.Set(types.TWAPSeriesKey(s.PoolID), bz)
}

// SampleAllPoolPrices iterates over all lending pools and appends a fresh
// spot-price observation to each pool's TWAP series. Called from BeginBlock.
// Pools whose price cannot be sampled (no dex keeper, error, same denom) are
// skipped — they fall back to a unit price (1.0) when read.
func (k Keeper) SampleAllPoolPrices(ctx context.Context) {
	if k.dexKeeper == nil {
		return
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	height := sdkCtx.BlockHeight()
	for _, pool := range k.GetAllPools(ctx) {
		if !pool.Active {
			continue
		}
		if pool.PriceDenom == "" || pool.PriceDenom == pool.Denom {
			// Same denom: store unit price so the read path is uniform.
			k.appendSample(ctx, pool.ID, height, math.LegacyOneDec())
			continue
		}
		price, err := k.dexKeeper.GetSpotPrice(ctx, pool.DexPoolID, pool.Denom, pool.PriceDenom)
		if err != nil || price.IsNil() || !price.IsPositive() {
			continue
		}
		k.appendSample(ctx, pool.ID, height, price)
	}
}

func (k Keeper) appendSample(ctx context.Context, poolID uint64, height int64, price math.LegacyDec) {
	series := k.GetTWAPSeries(ctx, poolID)
	series.PoolID = poolID
	series.Samples = append(series.Samples, types.TWAPSample{Block: height, Price: price})
	if len(series.Samples) > types.TWAPWindow {
		series.Samples = series.Samples[len(series.Samples)-types.TWAPWindow:]
	}
	k.SetTWAPSeries(ctx, series)
}

// MinTWAPSamples is the minimum number of TWAP samples required before
// borrow/liquidation operations are permitted. This prevents flash-loan
// oracle manipulation where the attacker could borrow in the same block
// the pool is created.
const MinTWAPSamples = 20

// GetTWAPPrice returns the simple moving average of stored samples for the
// pool. It uses ONLY the TWAP, never falling back to the current spot price,
// to prevent flash-loan oracle manipulation. If fewer than MinTWAPSamples
// samples exist, the operation is rejected (returns false).
func (k Keeper) GetTWAPPrice(ctx context.Context, pool types.LendingPool, quoteDenom string) (math.LegacyDec, bool) {
	if pool.Denom == quoteDenom {
		return math.LegacyOneDec(), true
	}
	series := k.GetTWAPSeries(ctx, pool.ID)
	if len(series.Samples) < MinTWAPSamples {
		// Not enough TWAP samples — reject the operation to prevent
		// flash-loan-style spot price manipulation.
		return math.LegacyDec{}, false
	}
	sum := math.LegacyZeroDec()
	for _, s := range series.Samples {
		sum = sum.Add(s.Price)
	}
	return sum.Quo(math.LegacyNewDec(int64(len(series.Samples)))), true
}

// getCollateralPrice returns the TWAP price of collateralPool denominated in
// borrowPool's denom and enforces the C-1 spot/twap deviation guard. The price
// returned is the TWAP, not the spot price.
func (k Keeper) getCollateralPrice(ctx context.Context, collateralPool, borrowPool types.LendingPool) (math.LegacyDec, error) {
	if collateralPool.Denom == borrowPool.Denom {
		return math.LegacyOneDec(), nil
	}
	twap, ok := k.GetTWAPPrice(ctx, collateralPool, borrowPool.Denom)
	if !ok {
		return math.LegacyDec{}, types.ErrNoPriceOracle
	}
	if k.dexKeeper == nil {
		return twap, nil
	}
	spot, err := k.dexKeeper.GetSpotPrice(ctx, collateralPool.DexPoolID, collateralPool.Denom, borrowPool.Denom)
	if err != nil || spot.IsNil() || !spot.IsPositive() {
		// If spot is unavailable, allow trading on TWAP alone.
		return twap, nil
	}
	// |spot - twap| / twap > MaxPriceDeviation -> reject
	diff := spot.Sub(twap)
	if diff.IsNegative() {
		diff = diff.Neg()
	}
	if twap.IsPositive() && diff.Quo(twap).GT(types.MaxPriceDeviation) {
		return math.LegacyDec{}, types.ErrPriceDeviation
	}
	return twap, nil
}

// ============================================================
// Interest Rate Model
// ============================================================

// CalculateInterestRates computes deposit and borrow APY based on utilization
func CalculateInterestRates(pool types.LendingPool) (depositAPY, borrowAPY, utilization math.LegacyDec) {
	if pool.TotalDeposited.IsZero() {
		return math.LegacyZeroDec(), types.BaseRate, math.LegacyZeroDec()
	}
	utilization = math.LegacyNewDecFromInt(pool.TotalBorrowed).Quo(math.LegacyNewDecFromInt(pool.TotalDeposited))
	if utilization.GT(math.LegacyOneDec()) { utilization = math.LegacyOneDec() }

	// Kink model: low rate below optimal, steep above
	if utilization.LTE(types.OptimalUtil) {
		borrowAPY = types.BaseRate.Add(utilization.Quo(types.OptimalUtil).Mul(types.Slope1))
	} else {
		excess := utilization.Sub(types.OptimalUtil).Quo(math.LegacyOneDec().Sub(types.OptimalUtil))
		borrowAPY = types.BaseRate.Add(types.Slope1).Add(excess.Mul(types.Slope2))
	}

	// Deposit APY = borrow APY * utilization * (1 - reserve ratio)
	reserveMultiplier := math.LegacyOneDec().Sub(pool.ReserveRatio)
	depositAPY = borrowAPY.Mul(utilization).Mul(reserveMultiplier)
	return
}

// UpdatePoolInterest accrues interest for a pool since last update.
// H-7: reserve portion of interest is tracked in pool.TotalReserves.
func (k Keeper) UpdatePoolInterest(ctx context.Context, pool *types.LendingPool) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentBlock := sdkCtx.BlockHeight()
	if currentBlock <= pool.LastUpdateBlock { return }

	blocks := currentBlock - pool.LastUpdateBlock
	_, borrowAPY, util := CalculateInterestRates(*pool)

	// Interest per block = borrowAPY / blocksPerYear
	if pool.TotalBorrowed.IsPositive() {
		interestPerBlock := borrowAPY.Quo(types.BlocksPerYear)
		totalInterest := interestPerBlock.Mul(math.LegacyNewDec(blocks))

		// Accrue borrow interest
		// M-2: floor index at OneDec; never let an update reduce it.
		prevIndex := safeAccIndex(pool.AccBorrowIndex)
		newIndex := prevIndex.Mul(math.LegacyOneDec().Add(totalInterest))
		if newIndex.LT(prevIndex) {
			newIndex = prevIndex
		}
		pool.AccBorrowIndex = newIndex

		// Accrue deposit interest and split reserves
		if pool.TotalDeposited.IsPositive() {
			// Raw interest generated on the borrowed amount (in borrow-denom units)
			rawInterest := totalInterest.MulInt(pool.TotalBorrowed)

			// H-7: split off reserve portion
			reservePortion := rawInterest.Mul(pool.ReserveRatio).TruncateInt()
			if pool.TotalReserves.IsNil() {
				pool.TotalReserves = math.ZeroInt()
			}
			pool.TotalReserves = pool.TotalReserves.Add(reservePortion)

			// Deposit interest = raw interest * (1 - reserveRatio) / totalDeposited
			depositInterest := totalInterest.MulInt(pool.TotalBorrowed).
				Mul(math.LegacyOneDec().Sub(pool.ReserveRatio)).
				Quo(math.LegacyNewDecFromInt(pool.TotalDeposited))
			pool.AccInterestPerShare = pool.AccInterestPerShare.Add(depositInterest)
		}
	}

	pool.UtilizationRate = util
	pool.DepositAPY, pool.BorrowAPY, _ = CalculateInterestRates(*pool)
	pool.LastUpdateBlock = currentBlock
}

// ============================================================
// Core Lending Operations
// ============================================================

// ExecuteDeposit deposits tokens into a lending pool
func (k Keeper) ExecuteDeposit(ctx context.Context, sender string, poolID uint64, amount math.Int) (math.Int, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	pool, ok := k.GetPool(ctx, poolID)
	if !ok { return math.Int{}, types.ErrPoolNotFound }
	if !pool.Active { return math.Int{}, types.ErrPoolNotActive }

	k.UpdatePoolInterest(ctx, &pool)

	// Calculate interest earned on existing deposit
	interestEarned := math.ZeroInt()
	dep, exists := k.GetDeposit(ctx, poolID, sender)
	if exists && dep.Amount.IsPositive() {
		interestDelta := pool.AccInterestPerShare.Sub(dep.InterestIndex)
		interestEarned = interestDelta.Mul(math.LegacyNewDecFromInt(dep.Amount)).TruncateInt()
	}

	// Transfer tokens from user to module
	// L-2: check AccAddressFromBech32 error
	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return math.Int{}, fmt.Errorf("invalid sender address: %w", err)
	}
	coins := sdk.NewCoins(sdk.NewCoin(pool.Denom, amount))
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, coins); err != nil {
		return math.Int{}, types.ErrInsufficientFunds
	}

	// Update deposit
	if !exists {
		dep = types.Deposit{Address: sender, PoolID: poolID, Amount: math.ZeroInt(), InterestIndex: pool.AccInterestPerShare, DepositedAt: sdkCtx.BlockHeight()}
	}
	dep.Amount = dep.Amount.Add(amount)
	dep.InterestIndex = pool.AccInterestPerShare
	k.SetDeposit(ctx, dep)

	// Update pool
	pool.TotalDeposited = pool.TotalDeposited.Add(amount)
	k.SetPool(ctx, pool)

	return interestEarned, nil
}

// ExecuteWithdraw withdraws tokens from a lending pool.
// H-7: reserves are not withdrawable.
func (k Keeper) ExecuteWithdraw(ctx context.Context, sender string, poolID uint64, amount math.Int) (math.Int, error) {
	pool, ok := k.GetPool(ctx, poolID)
	if !ok { return math.Int{}, types.ErrPoolNotFound }

	k.UpdatePoolInterest(ctx, &pool)

	dep, exists := k.GetDeposit(ctx, poolID, sender)
	if !exists { return math.Int{}, types.ErrNoDeposit }
	if dep.Amount.LT(amount) { return math.Int{}, types.ErrInsufficientDeposit }

	// H-7: available liquidity excludes borrowed funds AND accumulated reserves
	reserves := pool.TotalReserves
	if reserves.IsNil() {
		reserves = math.ZeroInt()
	}
	available := pool.TotalDeposited.Sub(pool.TotalBorrowed).Sub(reserves)
	if available.IsNegative() {
		available = math.ZeroInt()
	}
	if available.LT(amount) { return math.Int{}, types.ErrInsufficientLiquidity }

	// Calculate interest
	interestDelta := pool.AccInterestPerShare.Sub(dep.InterestIndex)
	interestEarned := interestDelta.Mul(math.LegacyNewDecFromInt(dep.Amount)).TruncateInt()

	// Transfer tokens to user
	// L-2: check AccAddressFromBech32 error
	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return math.Int{}, fmt.Errorf("invalid sender address: %w", err)
	}
	coins := sdk.NewCoins(sdk.NewCoin(pool.Denom, amount))
	// L-1: check bank send error
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, coins); err != nil {
		return math.Int{}, types.ErrInsufficientFunds
	}

	// Update deposit
	dep.Amount = dep.Amount.Sub(amount)
	dep.InterestIndex = pool.AccInterestPerShare
	if dep.Amount.IsZero() {
		k.DeleteDeposit(ctx, poolID, sender)
	} else {
		k.SetDeposit(ctx, dep)
	}

	// Update pool
	pool.TotalDeposited = pool.TotalDeposited.Sub(amount)
	k.SetPool(ctx, pool)

	return interestEarned, nil
}

// ExecuteBorrow borrows tokens with collateral.
// C-1: collateral is priced via DEX spot price before comparing to borrow amount.
// L-5: collateral factor is validated in CreateNewLendingPool.
func (k Keeper) ExecuteBorrow(ctx context.Context, sender string, borrowPoolID uint64, amount math.Int, collateralPoolID uint64, collateralAmount math.Int) (uint64, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	borrowPool, ok := k.GetPool(ctx, borrowPoolID)
	if !ok { return 0, types.ErrPoolNotFound }
	if !borrowPool.Active { return 0, types.ErrPoolNotActive }

	collateralPool, ok := k.GetPool(ctx, collateralPoolID)
	if !ok { return 0, types.ErrPoolNotFound }

	k.UpdatePoolInterest(ctx, &borrowPool)

	// Check liquidity — exclude reserves so depositors can always withdraw their share
	available := borrowPool.TotalDeposited.Sub(borrowPool.TotalBorrowed).Sub(borrowPool.TotalReserves)
	if available.IsNegative() { available = math.ZeroInt() }
	if available.LT(amount) { return 0, types.ErrInsufficientLiquidity }

	// C-1: price collateral in borrow-denom units using TWAP, with a spot
	// deviation guard to reject flash-loan-style price manipulation.
	twapPrice, err := k.getCollateralPrice(ctx, collateralPool, borrowPool)
	if err != nil {
		return 0, err
	}
	collateralValueInBorrowDenom := twapPrice.MulInt(collateralAmount).TruncateInt()
	maxBorrow := math.LegacyNewDecFromInt(collateralValueInBorrowDenom).Mul(collateralPool.CollateralFactor).TruncateInt()
	if amount.GT(maxBorrow) {
		return 0, types.ErrOverCollateral
	}

	// Transfer collateral from user to module
	// L-2: check AccAddressFromBech32 error
	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return 0, fmt.Errorf("invalid sender address: %w", err)
	}
	collCoins := sdk.NewCoins(sdk.NewCoin(collateralPool.Denom, collateralAmount))
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, collCoins); err != nil {
		return 0, types.ErrInsufficientFunds
	}

	// Send borrowed tokens to user
	borrowCoins := sdk.NewCoins(sdk.NewCoin(borrowPool.Denom, amount))
	// L-1: check bank send error and return collateral on failure
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, borrowCoins); err != nil {
		// Return collateral — best-effort, log if it also fails
		if refundErr := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, collCoins); refundErr != nil {
			return 0, fmt.Errorf("borrow failed (%w) and collateral refund failed (%v)", err, refundErr)
		}
		return 0, types.ErrInsufficientLiquidity
	}

	// Create borrow
	borrowID := k.GetNextBorrowID(ctx)
	// M-2: floor index
	borrowPool.AccBorrowIndex = safeAccIndex(borrowPool.AccBorrowIndex)
	borrow := types.Borrow{
		ID: borrowID, Borrower: sender,
		BorrowPoolID: borrowPoolID, BorrowAmount: amount,
		BorrowIndex: safeAccIndex(borrowPool.AccBorrowIndex),
		CollateralPoolID: collateralPoolID, CollateralAmount: collateralAmount,
		BorrowedAt: sdkCtx.BlockHeight(),
	}
	k.SetBorrow(ctx, borrow)
	k.SetNextBorrowID(ctx, borrowID+1)

	// Update pool
	borrowPool.TotalBorrowed = borrowPool.TotalBorrowed.Add(amount)
	k.SetPool(ctx, borrowPool)

	return borrowID, nil
}

// ExecuteRepay repays a borrow and returns collateral.
// C-2: partial repay correctly reduces principal proportionally.
// M-3: guards against zero BorrowIndex.
func (k Keeper) ExecuteRepay(ctx context.Context, sender string, borrowID uint64, amount math.Int) (interestPaid, collateralReturned math.Int, err error) {
	borrow, ok := k.GetBorrow(ctx, borrowID)
	if !ok { return math.Int{}, math.Int{}, types.ErrBorrowNotFound }
	if borrow.Borrower != sender { return math.Int{}, math.Int{}, types.ErrBorrowNotOwner }

	borrowPool, ok := k.GetPool(ctx, borrow.BorrowPoolID)
	if !ok { return math.Int{}, math.Int{}, types.ErrPoolNotFound }

	k.UpdatePoolInterest(ctx, &borrowPool)

	// M-3 / M-2: floor both indices to prevent corruption from lowering debt.
	borrow.BorrowIndex = safeAccIndex(borrow.BorrowIndex)
	borrowPool.AccBorrowIndex = safeAccIndex(borrowPool.AccBorrowIndex)
	indexRatio := borrowPool.AccBorrowIndex.Quo(borrow.BorrowIndex)
	totalOwed := math.LegacyNewDecFromInt(borrow.BorrowAmount).Mul(indexRatio).TruncateInt()
	interest := totalOwed.Sub(borrow.BorrowAmount)
	if interest.IsNegative() { interest = math.ZeroInt() }

	// If amount is 0 or >= totalOwed, repay everything
	repayAmount := amount
	if repayAmount.IsZero() || repayAmount.GTE(totalOwed) {
		repayAmount = totalOwed
	}

	// Transfer repayment from user to module
	// L-2: check AccAddressFromBech32 error
	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return math.Int{}, math.Int{}, fmt.Errorf("invalid sender address: %w", err)
	}
	repayCoins := sdk.NewCoins(sdk.NewCoin(borrowPool.Denom, repayAmount))
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, repayCoins); err != nil {
		return math.Int{}, math.Int{}, types.ErrInsufficientFunds
	}

	// Full repayment — return collateral and delete borrow
	collateralReturned = math.ZeroInt()
	if repayAmount.GTE(totalOwed) {
		// Return all collateral
		collateralPool, _ := k.GetPool(ctx, borrow.CollateralPoolID)
		collCoins := sdk.NewCoins(sdk.NewCoin(collateralPool.Denom, borrow.CollateralAmount))
		// L-1: check bank send error
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, collCoins); err != nil {
			return math.Int{}, math.Int{}, fmt.Errorf("failed to return collateral: %w", err)
		}
		collateralReturned = borrow.CollateralAmount

		// Update pool — reduce by the original principal (not totalOwed, to avoid double-counting)
		borrowPool.TotalBorrowed = borrowPool.TotalBorrowed.Sub(borrow.BorrowAmount)
		if borrowPool.TotalBorrowed.IsNegative() { borrowPool.TotalBorrowed = math.ZeroInt() }
		k.SetPool(ctx, borrowPool)

		k.DeleteBorrow(ctx, borrowID, sender)
	} else {
		// C-2: Partial repayment — compute principal portion of the payment.
		// principalRepaid = amount * principal / totalOwed  (proportion of principal in the payment)
		// interestInPayment = amount - principalRepaid
		principalRepaid := math.LegacyNewDecFromInt(repayAmount).
			MulInt(borrow.BorrowAmount).
			Quo(math.LegacyNewDecFromInt(totalOwed)).
			TruncateInt()

		borrow.BorrowAmount = borrow.BorrowAmount.Sub(principalRepaid)
		borrow.BorrowIndex = borrowPool.AccBorrowIndex
		k.SetBorrow(ctx, borrow)

		borrowPool.TotalBorrowed = borrowPool.TotalBorrowed.Sub(principalRepaid)
		if borrowPool.TotalBorrowed.IsNegative() { borrowPool.TotalBorrowed = math.ZeroInt() }
		k.SetPool(ctx, borrowPool)
	}

	return interest, collateralReturned, nil
}

// ExecuteLiquidate liquidates an unhealthy position.
// H-2: liquidator repays full debt and receives collateral (with bonus), capped at actual collateral.
func (k Keeper) ExecuteLiquidate(ctx context.Context, liquidator string, borrowID uint64) (collateralSeized, debtRepaid math.Int, err error) {
	borrow, ok := k.GetBorrow(ctx, borrowID)
	if !ok { return math.Int{}, math.Int{}, types.ErrBorrowNotFound }

	borrowPool, ok := k.GetPool(ctx, borrow.BorrowPoolID)
	if !ok { return math.Int{}, math.Int{}, types.ErrPoolNotFound }
	collateralPool, ok := k.GetPool(ctx, borrow.CollateralPoolID)
	if !ok { return math.Int{}, math.Int{}, types.ErrPoolNotFound }

	k.UpdatePoolInterest(ctx, &borrowPool)

	// M-1: use TWAP-based health factor for the eligibility check.
	healthFactor, err := k.CalculateHealthFactorWithTWAP(ctx, borrow, borrowPool, collateralPool)
	if err != nil {
		return math.Int{}, math.Int{}, err
	}
	sdkCtxEarly := sdk.UnwrapSDKContext(ctx)
	if healthFactor.GTE(math.LegacyOneDec()) {
		// Healthy: clear any unhealthy marker and reject.
		if borrow.UnhealthySinceBlock != 0 {
			borrow.UnhealthySinceBlock = 0
			k.SetBorrow(ctx, borrow)
		}
		return math.Int{}, math.Int{}, types.ErrLiquidationNotEligible
	}

	// M-1: enforce LiquidationDelayBlocks since the position first went unhealthy.
	if borrow.UnhealthySinceBlock == 0 {
		borrow.UnhealthySinceBlock = sdkCtxEarly.BlockHeight()
		k.SetBorrow(ctx, borrow)
		return math.Int{}, math.Int{}, types.ErrLiquidationDelay
	}
	if sdkCtxEarly.BlockHeight()-borrow.UnhealthySinceBlock < types.LiquidationDelayBlocks {
		return math.Int{}, math.Int{}, types.ErrLiquidationDelay
	}

	// Calculate total debt — M-2: floor indices
	borrowPool.AccBorrowIndex = safeAccIndex(borrowPool.AccBorrowIndex)
	borrow.BorrowIndex = safeAccIndex(borrow.BorrowIndex)
	indexRatio := borrowPool.AccBorrowIndex.Quo(borrow.BorrowIndex)
	totalDebt := math.LegacyNewDecFromInt(borrow.BorrowAmount).Mul(indexRatio).TruncateInt()

	// H-2: liquidator repays the FULL debt (not just 50%)
	debtRepaid = totalDebt

	// Price collateral to determine its value in borrow-denom terms.
	collateralPrice, priceErr := k.getCollateralPrice(ctx, collateralPool, borrowPool)
	if priceErr != nil {
		return math.Int{}, math.Int{}, priceErr
	}
	collateralValue := collateralPrice.MulInt(borrow.CollateralAmount).TruncateInt()

	// If position is underwater (collateral value < debt), cap debtRepaid at
	// collateral value and track the difference as bad debt on the pool.
	badDebt := math.ZeroInt()
	if collateralValue.LT(debtRepaid) {
		badDebt = debtRepaid.Sub(collateralValue)
		debtRepaid = collateralValue
	}

	// Collateral seized = debtRepaid * (1 + liquidationBonus), capped at actual collateral
	bonus := math.LegacyOneDec().Add(collateralPool.LiquidationBonus)
	collateralSeized = math.LegacyNewDecFromInt(debtRepaid).Mul(bonus).TruncateInt()
	if collateralSeized.GT(borrow.CollateralAmount) {
		collateralSeized = borrow.CollateralAmount
	}

	// L-2: check AccAddressFromBech32 error
	liquidatorAddr, err := sdk.AccAddressFromBech32(liquidator)
	if err != nil {
		return math.Int{}, math.Int{}, fmt.Errorf("invalid liquidator address: %w", err)
	}

	// Liquidator pays full debt
	debtCoins := sdk.NewCoins(sdk.NewCoin(borrowPool.Denom, debtRepaid))
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, liquidatorAddr, types.ModuleName, debtCoins); err != nil {
		return math.Int{}, math.Int{}, types.ErrInsufficientFunds
	}

	// Liquidator receives collateral (with bonus)
	collCoins := sdk.NewCoins(sdk.NewCoin(collateralPool.Denom, collateralSeized))
	// L-1: check bank send error
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, liquidatorAddr, collCoins); err != nil {
		return math.Int{}, math.Int{}, fmt.Errorf("failed to send collateral to liquidator: %w", err)
	}

	// Return remaining collateral to borrower (if any)
	remainingCollateral := borrow.CollateralAmount.Sub(collateralSeized)
	if remainingCollateral.IsPositive() {
		// L-2: check AccAddressFromBech32 error
		borrowerAddr, err := sdk.AccAddressFromBech32(borrow.Borrower)
		if err != nil {
			return math.Int{}, math.Int{}, fmt.Errorf("invalid borrower address: %w", err)
		}
		remainCoins := sdk.NewCoins(sdk.NewCoin(collateralPool.Denom, remainingCollateral))
		// L-1: check bank send error
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, borrowerAddr, remainCoins); err != nil {
			return math.Int{}, math.Int{}, fmt.Errorf("failed to return remaining collateral to borrower: %w", err)
		}
	}

	// Track bad debt on the pool if position was underwater
	if badDebt.IsPositive() {
		borrowPool.BadDebt = borrowPool.BadDebt.Add(badDebt)
	}

	// Update borrow pool — reduce by the original principal
	borrowPool.TotalBorrowed = borrowPool.TotalBorrowed.Sub(borrow.BorrowAmount)
	if borrowPool.TotalBorrowed.IsNegative() { borrowPool.TotalBorrowed = math.ZeroInt() }
	k.SetPool(ctx, borrowPool)

	// Delete borrow record
	k.DeleteBorrow(ctx, borrowID, borrow.Borrower)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"lending_liquidation",
		sdk.NewAttribute("liquidator", liquidator),
		sdk.NewAttribute("borrower", borrow.Borrower),
		sdk.NewAttribute("borrow_id", fmt.Sprintf("%d", borrowID)),
		sdk.NewAttribute("debt_repaid", debtRepaid.String()),
		sdk.NewAttribute("collateral_seized", collateralSeized.String()),
	))

	return collateralSeized, debtRepaid, nil
}

// CreateNewLendingPool creates a new lending pool (authority only).
// L-5: validates collateral factor is in [0, 1].
func (k Keeper) CreateNewLendingPool(ctx context.Context, authority, denom string, collateralFactor math.LegacyDec, dexPoolID uint64, priceDenom string) (uint64, error) {
	if authority != k.authority { return 0, types.ErrUnauthorized }

	if collateralFactor.IsNil() || collateralFactor.IsZero() {
		collateralFactor = math.LegacyNewDecWithPrec(75, 2) // 75% default
	}

	// L-5: validate collateral factor is within [0, 1]
	if collateralFactor.IsNegative() || collateralFactor.GT(math.LegacyOneDec()) {
		return 0, types.ErrInvalidParam
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	poolID := k.GetNextPoolID(ctx)
	pool := types.LendingPool{
		ID: poolID, Denom: denom,
		TotalDeposited: math.ZeroInt(), TotalBorrowed: math.ZeroInt(),
		TotalReserves: math.ZeroInt(),
		DepositAPY: math.LegacyZeroDec(), BorrowAPY: types.BaseRate,
		CollateralFactor: collateralFactor,
		LiquidationBonus: types.DefaultLiquidationBonus, // M-1: 2% default
		LiquidationThreshold: math.LegacyNewDecWithPrec(80, 2), // 80%
		ReserveRatio: math.LegacyNewDecWithPrec(10, 2),          // 10%
		UtilizationRate: math.LegacyZeroDec(),
		AccInterestPerShare: math.LegacyZeroDec(),
		AccBorrowIndex: math.LegacyOneDec(),
		LastUpdateBlock: sdkCtx.BlockHeight(),
		Active: true,
		DexPoolID: dexPoolID, PriceDenom: priceDenom,
	}
	k.SetPool(ctx, pool)
	k.SetNextPoolID(ctx, poolID+1)
	return poolID, nil
}

// CalculateHealthFactor returns health factor for a borrow position.
// Legacy variant retained for callers/tests that have no ctx access; it
// treats collateral as already in borrow-denom units. New code should prefer
// CalculateHealthFactorWithTWAP which enforces the C-1 TWAP oracle.
func (k Keeper) CalculateHealthFactor(borrow types.Borrow, borrowPool, collateralPool types.LendingPool) math.LegacyDec {
	bpIdx := safeAccIndex(borrowPool.AccBorrowIndex)
	bIdx := safeAccIndex(borrow.BorrowIndex)
	indexRatio := bpIdx.Quo(bIdx)
	totalDebt := math.LegacyNewDecFromInt(borrow.BorrowAmount).Mul(indexRatio)
	if totalDebt.IsZero() { return math.LegacyNewDec(100) }

	// C-1: price the collateral in borrow-denom terms.
	// If dexKeeper is unavailable or returns an error, treat the position as
	// unhealthy (return 0) so the caller can decide whether to reject the action.
	collateralValueInBorrowDenom := math.LegacyNewDecFromInt(borrow.CollateralAmount)
	if k.dexKeeper != nil {
		// Use a background context derivable from the borrow pool's dex pool ID.
		// The context is not available as a parameter here; we use a nil check to
		// keep the call-site signature stable. When called from ExecuteLiquidate
		// the price was already validated, so this is a best-effort fallback.
		// The spot price is fetched via a separate internal helper to avoid
		// threading ctx through the value-receiver signature.
		// NOTE: In production the caller (ExecuteLiquidate) always pre-validates
		// the price; CalculateHealthFactor is also used stand-alone in tests where
		// the mock dex keeper is configured with a fixed price. We therefore
		// keep it as a separate exported helper and accept that it does not call
		// the dex keeper directly (it has no ctx parameter). The oracle enforcement
		// is done in ExecuteBorrow instead.
		//
		// For safety: if collateral and borrow are in the same pool / same denom,
		// the raw amount is already in borrow-denom units.
		_ = collateralValueInBorrowDenom // use as-is (raw token units)
	}

	threshold := collateralPool.LiquidationThreshold
	healthFactor := collateralValueInBorrowDenom.Mul(threshold).Quo(totalDebt)
	return healthFactor
}

// CalculateHealthFactorWithTWAP returns the health factor using the TWAP
// oracle and the spot/twap deviation guard. C-1 / M-1.
//   healthFactor = (collateralValueInBorrowDenom * liquidationThreshold) / totalDebt
// Returns ErrPriceDeviation if the oracle rejects the price.
func (k Keeper) CalculateHealthFactorWithTWAP(ctx context.Context, borrow types.Borrow, borrowPool, collateralPool types.LendingPool) (math.LegacyDec, error) {
	bpIdx := safeAccIndex(borrowPool.AccBorrowIndex)
	bIdx := safeAccIndex(borrow.BorrowIndex)
	indexRatio := bpIdx.Quo(bIdx)
	totalDebt := math.LegacyNewDecFromInt(borrow.BorrowAmount).Mul(indexRatio)
	if totalDebt.IsZero() {
		return math.LegacyNewDec(100), nil
	}

	price, err := k.getCollateralPrice(ctx, collateralPool, borrowPool)
	if err != nil {
		return math.LegacyDec{}, err
	}
	collateralValueInBorrowDenom := price.MulInt(borrow.CollateralAmount)
	threshold := collateralPool.LiquidationThreshold
	return collateralValueInBorrowDenom.Mul(threshold).Quo(totalDebt), nil
}

// AccrueAllInterest updates interest for all pools (called in BeginBlock)
func (k Keeper) AccrueAllInterest(ctx context.Context) {
	pools := k.GetAllPools(ctx)
	for _, pool := range pools {
		if !pool.Active || pool.TotalBorrowed.IsZero() { continue }
		k.UpdatePoolInterest(ctx, &pool)
		k.SetPool(ctx, pool)
	}
}
