package keeper

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strconv"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/dex/types"
)

// maxPoolAmount caps individual pool deposit amounts at 10^30 to prevent precision overflow.
var maxPoolAmount = math.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(30), nil))

type Keeper struct {
	cdc                 codec.Codec
	storeService        store.KVStoreService
	accountKeeper       types.AccountKeeper
	bankKeeper          types.BankKeeper
	tokenFactoryKeeper  types.TokenFactoryKeeper
	stakingKeeper       types.StakingKeeper
	authority           string
}

// SetStakingKeeper sets the staking keeper for time-weighted governance queries.
func (k *Keeper) SetStakingKeeper(sk types.StakingKeeper) {
	k.stakingKeeper = sk
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	tokenFactoryKeeper types.TokenFactoryKeeper,
	authority string,
) Keeper {
	return Keeper{
		cdc:                cdc,
		storeService:       storeService,
		accountKeeper:      accountKeeper,
		bankKeeper:         bankKeeper,
		tokenFactoryKeeper: tokenFactoryKeeper,
		authority:          authority,
	}
}

func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// ---------------------------------------------------------------------------
// Params
// ---------------------------------------------------------------------------

func (k Keeper) GetParams(ctx context.Context) types.Params {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte("params"))
	if err != nil || bz == nil {
		return types.DefaultParams()
	}
	var params types.Params
	if err := json.Unmarshal(bz, &params); err != nil {
		return types.DefaultParams()
	}
	return params
}

func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(params)
	if err != nil {
		return err
	}
	return kvStore.Set([]byte("params"), bz)
}

// SetParamsWithAuthority sets the module params after verifying the sender is the module authority.
func (k Keeper) SetParamsWithAuthority(ctx context.Context, sender string, params types.Params) error {
	if sender != k.authority {
		return fmt.Errorf("unauthorized: only %s can set params", k.authority)
	}
	return k.SetParams(ctx, params)
}

// ---------------------------------------------------------------------------
// Pool CRUD
// ---------------------------------------------------------------------------

func (k Keeper) GetPool(ctx context.Context, poolID uint64) (types.Pool, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.PoolKey(poolID))
	if err != nil || bz == nil {
		return types.Pool{}, false
	}
	var pool types.Pool
	if err := json.Unmarshal(bz, &pool); err != nil {
		return types.Pool{}, false
	}
	return pool, true
}

func (k Keeper) SetPool(ctx context.Context, pool types.Pool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(pool)
	kvStore.Set(types.PoolKey(pool.ID), bz)
}

func (k Keeper) GetAllPools(ctx context.Context) []types.Pool {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.PoolPrefix)

	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var pools []types.Pool
	for ; iter.Valid(); iter.Next() {
		var pool types.Pool
		if err := json.Unmarshal(iter.Value(), &pool); err != nil {
			continue
		}
		pools = append(pools, pool)
	}
	sort.Slice(pools, func(i, j int) bool {
		return pools[i].ID < pools[j].ID
	})
	return pools
}

func (k Keeper) GetNextPoolID(ctx context.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte("next_pool_id"))
	if err != nil || bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetNextPoolID(ctx context.Context, id uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	kvStore.Set([]byte("next_pool_id"), bz)
}

func (k Keeper) GetPoolByDenomPair(ctx context.Context, denomA, denomB string) (types.Pool, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.PoolByDenomPairKey(denomA, denomB))
	if err != nil || bz == nil {
		return types.Pool{}, false
	}
	poolID := binary.BigEndian.Uint64(bz)
	return k.GetPool(ctx, poolID)
}

func (k Keeper) SetPoolByDenomPair(ctx context.Context, denomA, denomB string, poolID uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, poolID)
	kvStore.Set(types.PoolByDenomPairKey(denomA, denomB), bz)
}

// countPoolsByCreator returns the number of pools created by the given address.
func (k Keeper) countPoolsByCreator(ctx context.Context, creator string) int {
	pools := k.GetAllPools(ctx)
	count := 0
	for _, p := range pools {
		if p.Creator == creator {
			count++
		}
	}
	return count
}

// ---------------------------------------------------------------------------
// Core AMM: CreatePool
// ---------------------------------------------------------------------------

func (k Keeper) CreatePool(ctx context.Context, creator, denomA, denomB string, amountA, amountB math.Int) (uint64, error) {
	// Sort denoms alphabetically
	sortedA, sortedB := types.SortDenoms(denomA, denomB)
	// If the user passed them in reverse order, swap amounts too
	if sortedA != denomA {
		amountA, amountB = amountB, amountA
	}

	// Check pool doesn't already exist for this pair
	_, exists := k.GetPoolByDenomPair(ctx, sortedA, sortedB)
	if exists {
		return 0, types.ErrPoolAlreadyExists
	}

	// Check maximum amounts to prevent precision overflow
	if amountA.GT(maxPoolAmount) || amountB.GT(maxPoolAmount) {
		return 0, fmt.Errorf("amount exceeds maximum allowed")
	}

	// Check minimum initial liquidity
	params := k.GetParams(ctx)
	if amountA.LT(params.MinInitialLiquidity) || amountB.LT(params.MinInitialLiquidity) {
		return 0, types.ErrMinLiquidityNotMet
	}

	creatorAddr, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return 0, types.ErrInvalidSender
	}

	// Rate-limit pool creation per creator (DoS prevention)
	const MaxPoolsPerCreator = 50
	if k.countPoolsByCreator(ctx, creator) >= MaxPoolsPerCreator {
		return 0, fmt.Errorf("creator has reached maximum pool limit (%d)", MaxPoolsPerCreator)
	}

	// Transfer both token amounts from creator to module account
	coins := sdk.NewCoins(
		sdk.NewCoin(sortedA, amountA),
		sdk.NewCoin(sortedB, amountB),
	)
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, creatorAddr, types.ModuleName, coins); err != nil {
		return 0, err
	}

	// Get next pool ID
	poolID := k.GetNextPoolID(ctx)
	k.SetNextPoolID(ctx, poolID+1)

	// Calculate initial LP shares = sqrt(amountA * amountB)
	initialShares := isqrt(amountA.Mul(amountB))
	if initialShares.IsZero() {
		return 0, types.ErrZeroLiquidity
	}

	// Create LP token denom via tokenfactory
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	subdenom := "pool-" + strconv.FormatUint(poolID, 10)
	lpDenom, err := k.tokenFactoryKeeper.CreateDenom(ctx, moduleAddr.String(), subdenom)
	if err != nil {
		return 0, fmt.Errorf("failed to create LP denom: %w", err)
	}

	// Mint LP tokens to creator
	lpCoin := sdk.NewCoin(lpDenom, initialShares)
	if err := k.tokenFactoryKeeper.Mint(ctx, moduleAddr.String(), lpCoin, creator); err != nil {
		return 0, fmt.Errorf("failed to mint LP tokens: %w", err)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Create and store pool
	pool := types.Pool{
		ID:          poolID,
		DenomA:      sortedA,
		DenomB:      sortedB,
		ReserveA:    amountA,
		ReserveB:    amountB,
		TotalShares: initialShares,
		SwapFee:     params.DefaultSwapFee,
		Creator:     creator,
		CreatedAt:   sdkCtx.BlockHeight(),
	}
	k.SetPool(ctx, pool)
	k.SetPoolByDenomPair(ctx, sortedA, sortedB, poolID)

	k.Logger(ctx).Info("created pool",
		"pool_id", poolID,
		"denom_a", sortedA,
		"denom_b", sortedB,
		"reserve_a", amountA,
		"reserve_b", amountB,
		"lp_shares", initialShares,
		"creator", creator,
	)

	return poolID, nil
}

// ---------------------------------------------------------------------------
// Core AMM: AddLiquidity
// ---------------------------------------------------------------------------

func (k Keeper) AddLiquidity(ctx context.Context, sender string, poolID uint64, amountA, amountB, minSharesOut math.Int) (math.Int, error) {
	// Check maximum amounts to prevent precision overflow
	if amountA.GT(maxPoolAmount) || amountB.GT(maxPoolAmount) {
		return math.Int{}, fmt.Errorf("amount exceeds maximum allowed")
	}

	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return math.Int{}, types.ErrPoolNotFound
	}

	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return math.Int{}, types.ErrInvalidSender
	}

	// Calculate proportional deposit using the smaller ratio
	// ratioA = amountA / reserveA, ratioB = amountB / reserveB
	// shares = min(ratioA, ratioB) * totalShares
	// actual amounts = ratio * reserves

	// Use cross multiplication to avoid division:
	// ratioA < ratioB  iff  amountA * reserveB < amountB * reserveA
	crossA := amountA.Mul(pool.ReserveB) // amountA * reserveB
	crossB := amountB.Mul(pool.ReserveA) // amountB * reserveA

	var actualA, actualB, shares math.Int

	// Prevent re-adding liquidity to a fully drained pool — must recreate instead.
	if !pool.ReserveA.IsPositive() || !pool.ReserveB.IsPositive() {
		return math.Int{}, types.ErrPoolDrained
	}

	// Division-by-zero guard: reject if either reserve is zero
	if pool.ReserveA.IsZero() || pool.ReserveB.IsZero() {
		return math.Int{}, fmt.Errorf("pool %d has zero reserves, cannot add liquidity", poolID)
	}

	if crossA.LTE(crossB) {
		// ratioA <= ratioB: use ratioA
		actualA = amountA
		// actualB = amountA * reserveB / reserveA
		actualB = amountA.Mul(pool.ReserveB).Quo(pool.ReserveA)
		// shares = amountA * totalShares / reserveA
		shares = amountA.Mul(pool.TotalShares).Quo(pool.ReserveA)
	} else {
		// ratioB < ratioA: use ratioB
		actualB = amountB
		// actualA = amountB * reserveA / reserveB
		actualA = amountB.Mul(pool.ReserveA).Quo(pool.ReserveB)
		// shares = amountB * totalShares / reserveB
		shares = amountB.Mul(pool.TotalShares).Quo(pool.ReserveB)
	}

	if shares.IsZero() {
		return math.Int{}, types.ErrZeroLiquidity
	}

	// Check min shares out
	if shares.LT(minSharesOut) {
		return math.Int{}, types.ErrSlippageExceeded
	}

	// Transfer tokens from sender to module
	coins := sdk.NewCoins(
		sdk.NewCoin(pool.DenomA, actualA),
		sdk.NewCoin(pool.DenomB, actualB),
	)
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, coins); err != nil {
		return math.Int{}, err
	}

	// Mint LP tokens to sender
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	lpDenom := fmt.Sprintf("factory/%s/pool-%d", moduleAddr.String(), poolID)
	lpCoin := sdk.NewCoin(lpDenom, shares)
	if err := k.tokenFactoryKeeper.Mint(ctx, moduleAddr.String(), lpCoin, sender); err != nil {
		return math.Int{}, fmt.Errorf("failed to mint LP tokens: %w", err)
	}

	// Update pool
	pool.ReserveA = pool.ReserveA.Add(actualA)
	pool.ReserveB = pool.ReserveB.Add(actualB)
	pool.TotalShares = pool.TotalShares.Add(shares)
	k.SetPool(ctx, pool)

	k.Logger(ctx).Info("added liquidity",
		"pool_id", poolID,
		"sender", sender,
		"amount_a", actualA,
		"amount_b", actualB,
		"shares_minted", shares,
	)

	// Record liquidity add for sentiment oracle (total value = actualA + actualB)
	k.RecordLiquidityForSentiment(ctx, poolID, actualA.Add(actualB), true)

	return shares, nil
}

// ---------------------------------------------------------------------------
// Core AMM: RemoveLiquidity
// ---------------------------------------------------------------------------

func (k Keeper) RemoveLiquidity(ctx context.Context, sender string, poolID uint64, sharesIn, minAmountAOut, minAmountBOut math.Int) (math.Int, math.Int, error) {
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return math.Int{}, math.Int{}, types.ErrPoolNotFound
	}

	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return math.Int{}, math.Int{}, types.ErrInvalidSender
	}

	if sharesIn.GT(pool.TotalShares) {
		return math.Int{}, math.Int{}, types.ErrInsufficientLiquidity
	}

	// Calculate proportional withdrawal
	// amountA = sharesIn * reserveA / totalShares
	// amountB = sharesIn * reserveB / totalShares
	amountA := sharesIn.Mul(pool.ReserveA).Quo(pool.TotalShares)
	amountB := sharesIn.Mul(pool.ReserveB).Quo(pool.TotalShares)

	// Check min amounts
	if amountA.LT(minAmountAOut) || amountB.LT(minAmountBOut) {
		return math.Int{}, math.Int{}, types.ErrSlippageExceeded
	}

	// Prevent draining pool below minimum reserves (unless removing ALL liquidity)
	if !sharesIn.Equal(pool.TotalShares) {
		remainingA := pool.ReserveA.Sub(amountA)
		remainingB := pool.ReserveB.Sub(amountB)
		minReserve := math.NewInt(1000)
		if remainingA.LT(minReserve) || remainingB.LT(minReserve) {
			return math.Int{}, math.Int{}, fmt.Errorf("cannot drain pool below minimum reserve of %s", minReserve)
		}
	}

	// Send LP tokens from sender to module account, then burn
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	lpDenom := fmt.Sprintf("factory/%s/pool-%d", moduleAddr.String(), poolID)
	lpCoin := sdk.NewCoin(lpDenom, sharesIn)

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, sdk.NewCoins(lpCoin)); err != nil {
		return math.Int{}, math.Int{}, fmt.Errorf("failed to transfer LP tokens: %w", err)
	}
	if err := k.tokenFactoryKeeper.Burn(ctx, moduleAddr.String(), lpCoin, ""); err != nil {
		return math.Int{}, math.Int{}, fmt.Errorf("failed to burn LP tokens: %w", err)
	}

	// Send underlying tokens from module to sender
	coins := sdk.NewCoins(
		sdk.NewCoin(pool.DenomA, amountA),
		sdk.NewCoin(pool.DenomB, amountB),
	)
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, coins); err != nil {
		return math.Int{}, math.Int{}, err
	}

	// Update pool
	pool.ReserveA = pool.ReserveA.Sub(amountA)
	pool.ReserveB = pool.ReserveB.Sub(amountB)
	pool.TotalShares = pool.TotalShares.Sub(sharesIn)
	k.SetPool(ctx, pool)

	k.Logger(ctx).Info("removed liquidity",
		"pool_id", poolID,
		"sender", sender,
		"amount_a", amountA,
		"amount_b", amountB,
		"shares_burned", sharesIn,
	)

	// Record liquidity remove for sentiment oracle (total value = amountA + amountB)
	k.RecordLiquidityForSentiment(ctx, poolID, amountA.Add(amountB), false)

	return amountA, amountB, nil
}

// ---------------------------------------------------------------------------
// Core AMM: Swap
// ---------------------------------------------------------------------------

func (k Keeper) Swap(ctx context.Context, sender string, poolID uint64, tokenIn sdk.Coin, minTokenOut math.Int) (sdk.Coin, error) {
	// Risk engine: check if pool is halted or volume limit exceeded
	if err := k.CheckSwapRisk(ctx, poolID, tokenIn.Amount); err != nil {
		return sdk.Coin{}, err
	}

	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return sdk.Coin{}, types.ErrPoolNotFound
	}

	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return sdk.Coin{}, types.ErrInvalidSender
	}

	// Determine input/output reserves
	var reserveIn, reserveOut *math.Int
	var denomOut string

	if tokenIn.Denom == pool.DenomA {
		reserveIn = &pool.ReserveA
		reserveOut = &pool.ReserveB
		denomOut = pool.DenomB
	} else if tokenIn.Denom == pool.DenomB {
		reserveIn = &pool.ReserveB
		reserveOut = &pool.ReserveA
		denomOut = pool.DenomA
	} else {
		return sdk.Coin{}, types.ErrInvalidDenom
	}

	// Reject swaps that exceed 30% of the pool's input reserve
	maxSwapIn := reserveIn.MulRaw(30).QuoRaw(100)
	if tokenIn.Amount.GT(maxSwapIn) {
		return sdk.Coin{}, types.ErrSwapTooLarge
	}

	// Calculate output using constant-product formula with fee
	// Use dynamic fee if enabled, otherwise regular SwapFee
	effectiveFee := k.GetEffectiveFee(ctx, poolID)
	baseFeeBps := effectiveFee.MulInt64(10000).TruncateInt64()

	// Apply volume-based fee tier discount (multiplicative).
	// E.g., if baseFeeBps=30 and volume discount=10%, effectiveFeeBps=27.
	volumeDiscount := k.GetFeeDiscount(ctx, sender)
	if volumeDiscount.GT(math.LegacyOneDec()) {
		volumeDiscount = math.LegacyOneDec()
	}
	effectiveFeeBps := baseFeeBps
	if volumeDiscount.IsPositive() {
		// effectiveFeeBps = baseFeeBps * (1 - volumeDiscount)
		discountedFee := math.LegacyNewDec(baseFeeBps).Mul(math.LegacyOneDec().Sub(volumeDiscount))
		effectiveFeeBps = discountedFee.TruncateInt64()
	}
	// Clamp fee to valid range [0, 10000] to prevent underflow/overflow
	if effectiveFeeBps < 0 {
		effectiveFeeBps = 0
	}
	if effectiveFeeBps > 10000 {
		effectiveFeeBps = 10000
	}

	// amountInAfterFee = amountIn * (10000 - effectiveFeeBps) / 10000
	// outputAmount = reserveOut * amountInAfterFee / (reserveIn + amountInAfterFee)
	amountInAfterFee := tokenIn.Amount.MulRaw(10000 - effectiveFeeBps).QuoRaw(10000)

	// Guard against dust amounts: if the swap input rounds to zero after fees,
	// reject early to prevent zero-output swaps.
	if amountInAfterFee.IsZero() {
		return sdk.Coin{}, fmt.Errorf("swap amount too small after fees")
	}

	denominator := reserveIn.Add(amountInAfterFee)
	if denominator.IsZero() {
		return sdk.Coin{}, fmt.Errorf("swap denominator is zero")
	}
	outputAmount := reserveOut.Mul(amountInAfterFee).Quo(denominator)

	if outputAmount.IsZero() {
		return sdk.Coin{}, fmt.Errorf("output amount is zero, increase swap amount")
	}

	if outputAmount.GT(*reserveOut) {
		return sdk.Coin{}, types.ErrInsufficientLiquidity
	}

	// Check slippage
	if !minTokenOut.IsNil() && outputAmount.LT(minTokenOut) {
		return sdk.Coin{}, types.ErrSlippageExceeded
	}

	// Calculate referral fee split before transferring.
	// The referral discount is applied on the EFFECTIVE fee (after volume discount),
	// so both discounts stack multiplicatively.
	feeResult := k.CalculateReferralFees(ctx, sender, tokenIn.Amount, effectiveFeeBps)

	// FIX: Subtract BOTH TraderDiscount AND ReferrerAmount from reserves.
	// TraderDiscount leaves module immediately (sent to trader).
	// ReferrerAmount stays as claimable balance but is NOT pool liquidity —
	// it must be excluded from reserves to prevent insolvency when claimed.
	feeAmount := tokenIn.Amount.Sub(amountInAfterFee)
	referralTotal := math.ZeroInt()
	if feeResult.HasReferral {
		referralTotal = feeResult.TraderDiscount.Add(feeResult.ReferrerAmount)
		if referralTotal.GT(feeAmount) {
			referralTotal = feeAmount
		}
	}

	// Transfer tokenIn from sender to module
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, sdk.NewCoins(tokenIn)); err != nil {
		return sdk.Coin{}, err
	}

	// Transfer tokenOut from module to sender
	tokenOutCoin := sdk.NewCoin(denomOut, outputAmount)
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, sdk.NewCoins(tokenOutCoin)); err != nil {
		return sdk.Coin{}, err
	}

	// Apply referral fee sharing BEFORE updating pool state so that all operations
	// succeed atomically. If referral transfer fails, the swap is aborted entirely.
	if feeResult.HasReferral {
		if err := k.ApplyReferralFees(ctx, sender, poolID, tokenIn.Denom, feeResult, tokenIn.Amount); err != nil {
			return sdk.Coin{}, fmt.Errorf("referral fee application failed: %w", err)
		}
	}

	// Update reserves: full input amount goes to pool, minus referral payouts.
	// Referral fees are deducted from the fee portion so reserves stay consistent.
	reserveIncrease := tokenIn.Amount.Sub(referralTotal)
	if tokenIn.Denom == pool.DenomA {
		pool.ReserveA = pool.ReserveA.Add(reserveIncrease)
		pool.ReserveB = pool.ReserveB.Sub(outputAmount)
	} else {
		pool.ReserveB = pool.ReserveB.Add(reserveIncrease)
		pool.ReserveA = pool.ReserveA.Sub(outputAmount)
	}
	k.SetPool(ctx, pool)

	// Record swap volume in risk engine AFTER successful swap
	k.RecordSwapVolume(ctx, poolID, tokenIn.Amount)

	k.Logger(ctx).Info("swap executed",
		"pool_id", poolID,
		"sender", sender,
		"token_in", tokenIn,
		"token_out", tokenOutCoin,
	)

	// Record swap for sentiment oracle
	// "buy" means buying DenomA (i.e., input is DenomB)
	isBuy := tokenIn.Denom == pool.DenomB
	k.RecordSwapForSentiment(ctx, poolID, tokenIn.Amount, isBuy, pool)

	// Check for whale alert (swap > 5% of reserve)
	k.CheckAndRecordWhaleAlert(ctx, sender, poolID, tokenIn, pool)

	// Record trade for copy trading system. Pass reserveIncrease (L1) so the
	// PnL reconstruction sees the same delta the pool actually applied.
	k.RecordTradeForCopyTrading(ctx, sender, poolID, tokenIn, tokenOutCoin, reserveIncrease)

	// Record trade in persistent trade history for chart/candle queries.
	// Price convention: DenomA per DenomB (matches /pool/{id} price_ab).
	// Side is from the taker's perspective: "buy" means buying DenomB with DenomA.
	var tradePrice math.LegacyDec
	var tradeQty math.Int
	var tradeSide string
	if tokenIn.Denom == pool.DenomA {
		// Buying DenomB: spent DenomA, received DenomB
		if !outputAmount.IsZero() {
			tradePrice = math.LegacyNewDecFromInt(tokenIn.Amount).Quo(math.LegacyNewDecFromInt(outputAmount))
		} else {
			tradePrice = math.LegacyZeroDec()
		}
		tradeQty = outputAmount
		tradeSide = "buy"
	} else {
		// Selling DenomB: spent DenomB, received DenomA
		if !tokenIn.Amount.IsZero() {
			tradePrice = math.LegacyNewDecFromInt(outputAmount).Quo(math.LegacyNewDecFromInt(tokenIn.Amount))
		} else {
			tradePrice = math.LegacyZeroDec()
		}
		tradeQty = tokenIn.Amount
		tradeSide = "sell"
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	k.recordTrade(ctx, poolID, tradePrice, tradeQty, "", sender, 0, 0, tradeSide, "amm", sdkCtx.BlockHeight(), sdkCtx.BlockTime())

	// Record trade volume for fee tier tracking (uses tokenIn amount as volume)
	k.RecordTradeVolume(sdkCtx, sender, tokenIn.Amount)

	// Record price sample for dynamic fee volatility tracking
	if tradePrice.IsPositive() {
		k.RecordPriceSample(ctx, poolID, tradePrice)
	}

	return tokenOutCoin, nil
}

// ---------------------------------------------------------------------------
// GetQuote: read-only swap calculation
// ---------------------------------------------------------------------------

func (k Keeper) GetQuote(ctx context.Context, poolID uint64, tokenIn sdk.Coin) (sdk.Coin, math.LegacyDec, error) {
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return sdk.Coin{}, math.LegacyDec{}, types.ErrPoolNotFound
	}

	var reserveIn, reserveOut math.Int
	var denomOut string

	if tokenIn.Denom == pool.DenomA {
		reserveIn = pool.ReserveA
		reserveOut = pool.ReserveB
		denomOut = pool.DenomB
	} else if tokenIn.Denom == pool.DenomB {
		reserveIn = pool.ReserveB
		reserveOut = pool.ReserveA
		denomOut = pool.DenomA
	} else {
		return sdk.Coin{}, math.LegacyDec{}, types.ErrInvalidDenom
	}

	// Calculate output with fee
	feeBps := pool.SwapFee.MulInt64(10000).TruncateInt64()
	amountInAfterFee := tokenIn.Amount.MulRaw(10000 - feeBps).QuoRaw(10000)
	outputAmount := reserveOut.Mul(amountInAfterFee).Quo(reserveIn.Add(amountInAfterFee))

	// Calculate price impact
	// spotPrice = reserveOut / reserveIn
	// effectivePrice = outputAmount / amountIn
	// priceImpact = 1 - (effectivePrice / spotPrice)
	// Simplified: priceImpact = 1 - (outputAmount * reserveIn) / (amountIn * reserveOut)
	// Using sdk.Dec for precision
	spotNumerator := math.LegacyNewDecFromInt(outputAmount).Mul(math.LegacyNewDecFromInt(reserveIn))
	spotDenominator := math.LegacyNewDecFromInt(tokenIn.Amount).Mul(math.LegacyNewDecFromInt(reserveOut))

	priceImpact := math.LegacyZeroDec()
	if spotDenominator.IsPositive() {
		ratio := spotNumerator.Quo(spotDenominator)
		priceImpact = math.LegacyOneDec().Sub(ratio)
		if priceImpact.IsNegative() {
			priceImpact = math.LegacyZeroDec()
		}
	}

	tokenOut := sdk.NewCoin(denomOut, outputAmount)
	return tokenOut, priceImpact, nil
}

// ---------------------------------------------------------------------------
// GetSpotPrice returns the spot price of denomIn in terms of denomOut.
// Price = reserveOut / reserveIn (how much denomOut you get per 1 denomIn).
// ---------------------------------------------------------------------------

func (k Keeper) GetSpotPrice(ctx context.Context, poolID uint64, denomIn, denomOut string) (math.LegacyDec, error) {
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return math.LegacyDec{}, types.ErrPoolNotFound
	}

	var reserveIn, reserveOut math.Int

	if denomIn == pool.DenomA && denomOut == pool.DenomB {
		reserveIn = pool.ReserveA
		reserveOut = pool.ReserveB
	} else if denomIn == pool.DenomB && denomOut == pool.DenomA {
		reserveIn = pool.ReserveB
		reserveOut = pool.ReserveA
	} else {
		return math.LegacyDec{}, types.ErrInvalidDenom
	}

	if reserveIn.IsZero() {
		return math.LegacyDec{}, fmt.Errorf("zero reserve for %s", denomIn)
	}

	price := math.LegacyNewDecFromInt(reserveOut).Quo(math.LegacyNewDecFromInt(reserveIn))
	return price, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// isqrt computes the integer square root of x using Newton's method.
func isqrt(x math.Int) math.Int {
	if x.IsZero() || x.IsNegative() {
		return math.ZeroInt()
	}
	if x.Equal(math.OneInt()) {
		return math.OneInt()
	}

	// Newton's method: z = (z + x/z) / 2, starting with z = x
	z := x
	y := z.Add(math.OneInt()).QuoRaw(2) // (z+1)/2

	for y.LT(z) {
		z = y
		y = z.Add(x.Quo(z)).QuoRaw(2)
	}

	return z
}

// prefixEndBytes returns the end key for a prefix scan.
func prefixEndBytes(prefix []byte) []byte {
	if len(prefix) == 0 {
		return nil
	}
	end := make([]byte, len(prefix))
	copy(end, prefix)
	for i := len(end) - 1; i >= 0; i-- {
		end[i]++
		if end[i] != 0 {
			return end
		}
	}
	return nil
}

// GetPoolDenoms returns the two denoms of a pool (for cross-module use).
func (k Keeper) GetPoolDenoms(ctx context.Context, poolID uint64) (string, string, bool) {
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return "", "", false
	}
	return pool.DenomA, pool.DenomB, true
}
