package keeper

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/options/types"
)

const (
	// Approximate blocks per year at ~5s block time
	blocksPerYear = int64(6_307_200)
)

type Keeper struct {
	cdc          codec.Codec
	storeService store.KVStoreService
	bankKeeper   types.BankKeeper
	dexKeeper    types.DexKeeper
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

func (k *Keeper) SetDexKeeper(dk types.DexKeeper) { k.dexKeeper = dk }

// ============================================================
// Option CRUD
// ============================================================

func (k Keeper) GetOption(ctx context.Context, id uint64) (types.Option, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.OptionKey(id))
	if err != nil || bz == nil {
		return types.Option{}, false
	}
	var opt types.Option
	if err := json.Unmarshal(bz, &opt); err != nil {
		return types.Option{}, false
	}
	return opt, true
}

func (k Keeper) SetOption(ctx context.Context, opt types.Option) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(opt)
	_ = kvStore.Set(types.OptionKey(opt.ID), bz)
}

func (k Keeper) GetAllOptions(ctx context.Context) []types.Option {
	kvStore := k.storeService.OpenKVStore(ctx)
	endKey := append([]byte(types.OptionPrefix), 0xFF)
	iter, err := kvStore.Iterator([]byte(types.OptionPrefix), endKey)
	if err != nil {
		return nil
	}
	defer iter.Close()
	var opts []types.Option
	for ; iter.Valid(); iter.Next() {
		var opt types.Option
		if err := json.Unmarshal(iter.Value(), &opt); err == nil {
			opts = append(opts, opt)
		}
	}
	return opts
}

func (k Keeper) GetNextOptionID(ctx context.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := kvStore.Get([]byte(types.NextOptionIDKey))
	if bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetNextOptionID(ctx context.Context, id uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	_ = kvStore.Set([]byte(types.NextOptionIDKey), bz)
}

func (k Keeper) setExpiryIndex(ctx context.Context, opt types.Option) {
	kvStore := k.storeService.OpenKVStore(ctx)
	_ = kvStore.Set(types.ExpiryIndexKey(opt.ExpiryBlock, opt.ID), []byte{1})
}

func (k Keeper) deleteExpiryIndex(ctx context.Context, opt types.Option) {
	kvStore := k.storeService.OpenKVStore(ctx)
	_ = kvStore.Delete(types.ExpiryIndexKey(opt.ExpiryBlock, opt.ID))
}

// ============================================================
// Black-Scholes Premium Calculation
// ============================================================

// CalculatePremium uses a simplified Black-Scholes model safe for on-chain integer math.
// S = current spot price, K = strike, T = blocks to expiry, sigma = annualized volatility (0-1).
// Returns premium in quote denom base units.
func (k Keeper) CalculatePremium(
	spotPrice, strikePrice math.LegacyDec,
	blocksToExpiry int64,
	sigma math.LegacyDec,
	optionType types.OptionType,
	amount math.Int,
) math.Int {
	if blocksToExpiry <= 0 || spotPrice.IsZero() || strikePrice.IsZero() {
		return math.ZeroInt()
	}
	if sigma.IsNil() || sigma.IsZero() {
		sigma = math.LegacyNewDecWithPrec(30, 2) // default 30% volatility
	}

	S := spotPrice
	K := strikePrice

	// T in years
	T := math.LegacyNewDec(blocksToExpiry).Quo(math.LegacyNewDec(blocksPerYear))

	// intrinsic value = max(S-K, 0) for call, max(K-S, 0) for put
	var intrinsic math.LegacyDec
	if optionType == types.OptionTypeCall {
		intrinsic = S.Sub(K)
		if intrinsic.IsNegative() {
			intrinsic = math.LegacyZeroDec()
		}
	} else {
		intrinsic = K.Sub(S)
		if intrinsic.IsNegative() {
			intrinsic = math.LegacyZeroDec()
		}
	}

	// time value ≈ S * sigma * sqrt(T) * 0.4
	// This is the at-the-money approximation from simplified B-S
	sqrtT := sqrtDec(T)
	timeValue := S.Mul(sigma).Mul(sqrtT).Mul(math.LegacyNewDecWithPrec(4, 1))

	premiumPerUnit := intrinsic.Add(timeValue)
	if premiumPerUnit.IsNegative() {
		premiumPerUnit = math.LegacyZeroDec()
	}

	// scale by amount
	total := premiumPerUnit.MulInt(amount)
	return total.TruncateInt()
}

// sqrtDec computes square root of a LegacyDec using Newton's method (10 iterations).
func sqrtDec(x math.LegacyDec) math.LegacyDec {
	if x.IsZero() || x.IsNegative() {
		return math.LegacyZeroDec()
	}
	// initial guess
	guess := x.Quo(math.LegacyNewDec(2))
	if guess.IsZero() {
		guess = math.LegacyNewDecWithPrec(1, 6)
	}
	two := math.LegacyNewDec(2)
	for i := 0; i < 10; i++ {
		guess = guess.Add(x.Quo(guess)).Quo(two)
	}
	return guess
}

// ============================================================
// Core Logic
// ============================================================

// ExecuteWriteOption creates an option, locks collateral from writer, returns option ID + premium.
func (k Keeper) ExecuteWriteOption(
	ctx context.Context,
	writer string,
	poolID uint64,
	optionType types.OptionType,
	strikePrice math.LegacyDec,
	amount math.Int,
	expiryBlock int64,
	customPremium math.Int,
) (uint64, math.Int, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// get pool denoms
	denomA, denomB, found := k.dexKeeper.GetPoolDenoms(ctx, poolID)
	if !found {
		return 0, math.ZeroInt(), fmt.Errorf("pool %d not found", poolID)
	}
	underlyingDenom := denomA
	quoteDenom := denomB

	if expiryBlock <= sdkCtx.BlockHeight() {
		return 0, math.ZeroInt(), types.ErrInvalidExpiry
	}

	// get spot price
	spotPrice, err := k.dexKeeper.GetSpotPrice(ctx, poolID, underlyingDenom, quoteDenom)
	if err != nil || spotPrice.IsZero() {
		return 0, math.ZeroInt(), fmt.Errorf("could not get spot price: %w", err)
	}

	// calculate premium
	sigma := k.dexKeeper.GetSignalVolatility(ctx, poolID)
	blocksToExpiry := expiryBlock - sdkCtx.BlockHeight()
	calculatedPremium := k.CalculatePremium(spotPrice, strikePrice, blocksToExpiry, sigma, optionType, amount)
	if calculatedPremium.IsZero() {
		calculatedPremium = math.NewInt(1) // minimum 1 unit
	}

	var premium math.Int
	if !customPremium.IsNil() && customPremium.IsPositive() {
		// Custom premium must be at least 50% of the calculated premium to
		// prevent drastically underpriced options used for wash trading.
		minAllowed := calculatedPremium.Quo(math.NewInt(2))
		if customPremium.LT(minAllowed) {
			return 0, math.ZeroInt(), fmt.Errorf("custom premium %s is below minimum allowed %s (50%% of calculated)", customPremium, minAllowed)
		}
		premium = customPremium
	} else {
		premium = calculatedPremium
	}

	// lock collateral from writer
	// For a call: writer locks the underlying tokens (amount of underlyingDenom)
	// For a put: writer locks the quote tokens (strikePrice * amount)
	writerAddr, err := sdk.AccAddressFromBech32(writer)
	if err != nil {
		return 0, math.ZeroInt(), err
	}

	var collateral sdk.Coins
	if optionType == types.OptionTypeCall {
		collateral = sdk.NewCoins(sdk.NewCoin(underlyingDenom, amount))
	} else {
		// put: lock strike * amount in quote denom
		collateralAmt := strikePrice.MulInt(amount).TruncateInt()
		collateral = sdk.NewCoins(sdk.NewCoin(quoteDenom, collateralAmt))
	}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, writerAddr, types.ModuleName, collateral); err != nil {
		return 0, math.ZeroInt(), types.ErrInsufficientFunds
	}

	// store option
	id := k.GetNextOptionID(ctx)
	opt := types.Option{
		ID:              id,
		Writer:          writer,
		PoolID:          poolID,
		UnderlyingDenom: underlyingDenom,
		QuoteDenom:      quoteDenom,
		OptionType:      optionType,
		StrikePrice:     strikePrice,
		Premium:         premium,
		Amount:          amount,
		ExpiryBlock:     expiryBlock,
		Status:          types.OptionStatusOpen,
		CreatedAt:       sdkCtx.BlockHeight(),
	}
	k.SetOption(ctx, opt)
	k.setExpiryIndex(ctx, opt)
	k.SetNextOptionID(ctx, id+1)

	return id, premium, nil
}

// ExecuteBuyOption transfers premium from buyer to writer and activates the option.
func (k Keeper) ExecuteBuyOption(ctx context.Context, buyer string, optionID uint64) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	opt, found := k.GetOption(ctx, optionID)
	if !found {
		return types.ErrOptionNotFound
	}
	if opt.Status != types.OptionStatusOpen {
		return types.ErrOptionNotOpen
	}
	if sdkCtx.BlockHeight() >= opt.ExpiryBlock {
		return types.ErrOptionExpired
	}

	buyerAddr, err := sdk.AccAddressFromBech32(buyer)
	if err != nil {
		return err
	}
	writerAddr, err := sdk.AccAddressFromBech32(opt.Writer)
	if err != nil {
		return err
	}

	// transfer premium from buyer to writer
	premiumCoins := sdk.NewCoins(sdk.NewCoin(opt.QuoteDenom, opt.Premium))
	if err := k.bankKeeper.SendCoins(ctx, buyerAddr, writerAddr, premiumCoins); err != nil {
		return types.ErrInsufficientFunds
	}

	opt.Buyer = buyer
	opt.Status = types.OptionStatusActive
	k.SetOption(ctx, opt)

	return nil
}

// ExecuteExerciseOption settles an active option and pays out the buyer.
func (k Keeper) ExecuteExerciseOption(ctx context.Context, buyer string, optionID uint64) (math.Int, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	opt, found := k.GetOption(ctx, optionID)
	if !found {
		return math.ZeroInt(), types.ErrOptionNotFound
	}
	if opt.Status != types.OptionStatusActive {
		return math.ZeroInt(), types.ErrOptionNotActive
	}
	if opt.Buyer != buyer {
		return math.ZeroInt(), types.ErrUnauthorized
	}
	if sdkCtx.BlockHeight() > opt.ExpiryBlock {
		return math.ZeroInt(), types.ErrOptionExpired
	}

	// get current price
	spotPrice, err := k.dexKeeper.GetSpotPrice(ctx, opt.PoolID, opt.UnderlyingDenom, opt.QuoteDenom)
	if err != nil {
		return math.ZeroInt(), fmt.Errorf("could not get spot price: %w", err)
	}

	buyerAddr, err := sdk.AccAddressFromBech32(buyer)
	if err != nil {
		return math.ZeroInt(), err
	}
	writerAddr, err := sdk.AccAddressFromBech32(opt.Writer)
	if err != nil {
		return math.ZeroInt(), err
	}

	var payout math.Int
	if opt.OptionType == types.OptionTypeCall {
		// call is profitable if spot > strike
		if spotPrice.LTE(opt.StrikePrice) {
			return math.ZeroInt(), types.ErrNotProfitable
		}

		// Physical settlement: buyer pays strike * amount in quote denom to writer
		strikePayment := opt.StrikePrice.MulInt(opt.Amount).TruncateInt()
		strikeCoins := sdk.NewCoins(sdk.NewCoin(opt.QuoteDenom, strikePayment))
		if err := k.bankKeeper.SendCoins(ctx, buyerAddr, writerAddr, strikeCoins); err != nil {
			return math.ZeroInt(), types.ErrInsufficientFunds
		}

		// Release underlying collateral from module to buyer
		underlyingCoins := sdk.NewCoins(sdk.NewCoin(opt.UnderlyingDenom, opt.Amount))
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, buyerAddr, underlyingCoins); err != nil {
			return math.ZeroInt(), err
		}
		payout = opt.Amount // return in underlying units

	} else {
		// put is profitable if spot < strike
		if spotPrice.GTE(opt.StrikePrice) {
			return math.ZeroInt(), types.ErrNotProfitable
		}
		// payout = (strike - spot) * amount in quote denom
		payout = opt.StrikePrice.Sub(spotPrice).MulInt(opt.Amount).TruncateInt()
		// module locked strikePrice * amount in quote denom
		totalLocked := opt.StrikePrice.MulInt(opt.Amount).TruncateInt()
		payoutCoins := sdk.NewCoins(sdk.NewCoin(opt.QuoteDenom, payout))
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, buyerAddr, payoutCoins); err != nil {
			return math.ZeroInt(), err
		}
		// return remainder to writer
		remainder := totalLocked.Sub(payout)
		if remainder.IsPositive() {
			writerAddr2, err := sdk.AccAddressFromBech32(opt.Writer)
			if err != nil {
				return math.ZeroInt(), err
			}
			remainderCoins := sdk.NewCoins(sdk.NewCoin(opt.QuoteDenom, remainder))
			if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, writerAddr2, remainderCoins); err != nil {
				return math.ZeroInt(), err
			}
		}
	}

	opt.Status = types.OptionStatusExercised
	opt.ExercisedAt = sdkCtx.BlockHeight()
	k.SetOption(ctx, opt)
	k.deleteExpiryIndex(ctx, opt)

	return payout, nil
}

// ExecuteCancelOption returns collateral to writer for an unsold (open) option.
func (k Keeper) ExecuteCancelOption(ctx context.Context, writer string, optionID uint64) error {
	opt, found := k.GetOption(ctx, optionID)
	if !found {
		return types.ErrOptionNotFound
	}
	if opt.Writer != writer {
		return types.ErrUnauthorized
	}
	if opt.Status != types.OptionStatusOpen {
		return fmt.Errorf("can only cancel open options, status: %s", opt.Status)
	}

	writerAddr, err := sdk.AccAddressFromBech32(writer)
	if err != nil {
		return err
	}

	// return collateral
	var collateral sdk.Coins
	if opt.OptionType == types.OptionTypeCall {
		collateral = sdk.NewCoins(sdk.NewCoin(opt.UnderlyingDenom, opt.Amount))
	} else {
		collateralAmt := opt.StrikePrice.MulInt(opt.Amount).TruncateInt()
		collateral = sdk.NewCoins(sdk.NewCoin(opt.QuoteDenom, collateralAmt))
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, writerAddr, collateral); err != nil {
		return err
	}

	opt.Status = types.OptionStatusCancelled
	k.SetOption(ctx, opt)
	k.deleteExpiryIndex(ctx, opt)

	return nil
}

// BeginBlockHandler expires options that have passed their expiry block.
// Returns collateral to writers for options that were never exercised.
func (k Keeper) BeginBlockHandler(ctx context.Context) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	height := sdkCtx.BlockHeight()

	// iterate expiry index for this block
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := types.ExpiryIndexPrefixForBlock(height)
	endKey := types.ExpiryIndexPrefixForBlock(height + 1)
	iter, err := kvStore.Iterator(prefix, endKey)
	if err != nil {
		return
	}
	defer iter.Close()

	var toExpire []uint64
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		// key = prefix/block(8bytes)/id(8bytes) - check still in range
		if len(key) < 8 {
			break
		}
		// extract id from last 8 bytes
		idBytes := key[len(key)-8:]
		id := binary.BigEndian.Uint64(idBytes)
		toExpire = append(toExpire, id)
	}
	iter.Close()

	for _, id := range toExpire {
		opt, found := k.GetOption(ctx, id)
		if !found || (opt.Status != types.OptionStatusOpen && opt.Status != types.OptionStatusActive) {
			continue
		}

		// return collateral to writer
		writerAddr, err := sdk.AccAddressFromBech32(opt.Writer)
		if err != nil {
			continue
		}

		var collateral sdk.Coins
		if opt.OptionType == types.OptionTypeCall {
			collateral = sdk.NewCoins(sdk.NewCoin(opt.UnderlyingDenom, opt.Amount))
		} else {
			collateralAmt := opt.StrikePrice.MulInt(opt.Amount).TruncateInt()
			if collateralAmt.IsPositive() {
				collateral = sdk.NewCoins(sdk.NewCoin(opt.QuoteDenom, collateralAmt))
			}
		}

		if len(collateral) > 0 {
			if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, writerAddr, collateral); err != nil {
				// collateral return failed; leave option for retry next block
				continue
			}
		}

		opt.Status = types.OptionStatusExpired
		k.SetOption(ctx, opt)
		k.deleteExpiryIndex(ctx, opt)
	}
}
