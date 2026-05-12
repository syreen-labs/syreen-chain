package keeper

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"sort"

	"cosmossdk.io/core/store"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/predict/types"
)

// prefixEnd returns the end key for a prefix range scan
func prefixEnd(prefix []byte) []byte {
	return append(prefix, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF)
}

type Keeper struct {
	cdc           codec.Codec
	storeService  store.KVStoreService
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
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

// ============================================================
// Market CRUD
// ============================================================

func (k Keeper) GetMarket(ctx context.Context, marketID uint64) (types.Market, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.MarketKey(marketID))
	if err != nil || bz == nil {
		return types.Market{}, false
	}
	var market types.Market
	if err := json.Unmarshal(bz, &market); err != nil {
		return types.Market{}, false
	}
	return market, true
}

func (k Keeper) SetMarket(ctx context.Context, market types.Market) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(market)
	_ = kvStore.Set(types.MarketKey(market.ID), bz)
}

func (k Keeper) GetAllMarkets(ctx context.Context) []types.Market {
	kvStore := k.storeService.OpenKVStore(ctx)
	iter, err := kvStore.Iterator([]byte(types.MarketPrefix), prefixEnd([]byte(types.MarketPrefix)))
	if err != nil {
		return nil
	}
	defer iter.Close()
	var markets []types.Market
	for ; iter.Valid(); iter.Next() {
		var m types.Market
		if err := json.Unmarshal(iter.Value(), &m); err == nil {
			markets = append(markets, m)
		}
	}
	sort.Slice(markets, func(i, j int) bool { return markets[i].ID < markets[j].ID })
	return markets
}

func (k Keeper) GetNextMarketID(ctx context.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.NextMarketIDKey))
	if err != nil || bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetNextMarketID(ctx context.Context, id uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	_ = kvStore.Set([]byte(types.NextMarketIDKey), bz)
}

// ============================================================
// Position CRUD
// ============================================================

func (k Keeper) GetPosition(ctx context.Context, marketID uint64, address string) (types.Position, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.PositionKey(marketID, address))
	if err != nil || bz == nil {
		return types.Position{}, false
	}
	var pos types.Position
	if err := json.Unmarshal(bz, &pos); err != nil {
		return types.Position{}, false
	}
	return pos, true
}

func (k Keeper) SetPosition(ctx context.Context, pos types.Position) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(pos)
	_ = kvStore.Set(types.PositionKey(pos.MarketID, pos.Address), bz)
	_ = kvStore.Set(types.PosByAddrKey(pos.Address, pos.MarketID), []byte{1})
}

func (k Keeper) GetPositionsForMarket(ctx context.Context, marketID uint64) []types.Position {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := types.PositionMarketPrefix(marketID)
	iter, err := kvStore.Iterator(prefix, prefixEnd(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()
	var positions []types.Position
	for ; iter.Valid(); iter.Next() {
		var p types.Position
		if err := json.Unmarshal(iter.Value(), &p); err == nil {
			positions = append(positions, p)
		}
	}
	return positions
}

// ============================================================
// Resolution CRUD
// ============================================================

func (k Keeper) GetResolution(ctx context.Context, marketID uint64) (types.Resolution, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.ResolutionKey(marketID))
	if err != nil || bz == nil {
		return types.Resolution{}, false
	}
	var res types.Resolution
	if err := json.Unmarshal(bz, &res); err != nil {
		return types.Resolution{}, false
	}
	return res, true
}

func (k Keeper) SetResolution(ctx context.Context, res types.Resolution) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(res)
	_ = kvStore.Set(types.ResolutionKey(res.MarketID), bz)
}

// ============================================================
// Core Business Logic
// ============================================================

// ExecuteCreateMarket creates a new prediction market with equal YES/NO share pools
func (k Keeper) ExecuteCreateMarket(ctx context.Context, creator, question, resolver, quoteDenom string, resolutionBlock int64, initialLiquidity math.Int) (uint64, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	creatorAddr, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return 0, err
	}

	coins := sdk.NewCoins(sdk.NewCoin(quoteDenom, initialLiquidity))
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, creatorAddr, types.ModuleName, coins); err != nil {
		return 0, err
	}

	marketID := k.GetNextMarketID(ctx)
	k.SetNextMarketID(ctx, marketID+1)

	// Start with equal YES and NO shares = initialLiquidity each
	// This gives 50/50 pricing: yesPrice = noShares/(yesShares+noShares) = 0.5
	market := types.Market{
		ID:              marketID,
		Question:        question,
		Creator:         creator,
		Resolver:        resolver,
		QuoteDenom:      quoteDenom,
		ResolutionBlock: resolutionBlock,
		Status:          types.MarketStatusOpen,
		YesShares:       initialLiquidity,
		NoShares:        initialLiquidity,
		Liquidity:       initialLiquidity,
		TotalVolume:     math.ZeroInt(),
		CreatedAt:       sdkCtx.BlockHeight(),
	}

	k.SetMarket(ctx, market)
	return marketID, nil
}

// ExecuteBuyShares executes a share purchase using constant-product AMM math.
// Buying YES shares: user adds quote tokens which increase the NO pool,
// and receives YES shares removed from the YES pool.
// The invariant is: yesShares * noShares = k (constant product).
//
// For buying YES:
//   - newNoShares = noShares + amount
//   - newYesShares = k / newNoShares
//   - sharesBought = yesShares - newYesShares
func (k Keeper) ExecuteBuyShares(ctx context.Context, sender string, marketID uint64, outcome string, amount math.Int) (math.Int, error) {
	if amount.IsNil() || !amount.IsPositive() {
		return math.Int{}, types.ErrInvalidAmount
	}

	market, found := k.GetMarket(ctx, marketID)
	if !found {
		return math.Int{}, types.ErrMarketNotFound
	}
	if market.Status != types.MarketStatusOpen {
		return math.Int{}, types.ErrMarketNotOpen
	}

	// Reject trades at or after the resolution block to prevent insider trading.
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if market.ResolutionBlock > 0 && sdkCtx.BlockHeight() >= market.ResolutionBlock {
		return math.Int{}, types.ErrMarketNotOpen
	}

	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return math.Int{}, err
	}

	coins := sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, amount))
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, coins); err != nil {
		return math.Int{}, err
	}

	kProduct := market.YesShares.Mul(market.NoShares)

	var sharesBought math.Int
	if outcome == "yes" {
		newNoShares := market.NoShares.Add(amount)
		newYesShares := kProduct.Quo(newNoShares)
		sharesBought = market.YesShares.Sub(newYesShares)
		if sharesBought.IsZero() || sharesBought.IsNegative() {
			return math.Int{}, types.ErrInvalidAmount
		}
		market.YesShares = newYesShares
		market.NoShares = newNoShares
	} else {
		newYesShares := market.YesShares.Add(amount)
		newNoShares := kProduct.Quo(newYesShares)
		sharesBought = market.NoShares.Sub(newNoShares)
		if sharesBought.IsZero() || sharesBought.IsNegative() {
			return math.Int{}, types.ErrInvalidAmount
		}
		market.YesShares = newYesShares
		market.NoShares = newNoShares
	}

	market.TotalVolume = market.TotalVolume.Add(amount)
	k.SetMarket(ctx, market)

	// Update position
	pos, found := k.GetPosition(ctx, marketID, sender)
	if !found {
		pos = types.Position{
			Address:    sender,
			MarketID:   marketID,
			YesShares:  math.ZeroInt(),
			NoShares:   math.ZeroInt(),
			TotalSpent: math.ZeroInt(),
		}
	}
	if outcome == "yes" {
		pos.YesShares = pos.YesShares.Add(sharesBought)
	} else {
		pos.NoShares = pos.NoShares.Add(sharesBought)
	}
	pos.TotalSpent = pos.TotalSpent.Add(amount)
	k.SetPosition(ctx, pos)

	return sharesBought, nil
}

// ExecuteSellShares sells shares back to the market using reverse CPMM math.
func (k Keeper) ExecuteSellShares(ctx context.Context, sender string, marketID uint64, outcome string, shares math.Int) (math.Int, error) {
	if shares.IsNil() || !shares.IsPositive() {
		return math.Int{}, types.ErrInvalidAmount
	}

	market, found := k.GetMarket(ctx, marketID)
	if !found {
		return math.Int{}, types.ErrMarketNotFound
	}
	if market.Status != types.MarketStatusOpen {
		return math.Int{}, types.ErrMarketNotOpen
	}

	// Reject trades at or after the resolution block
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if market.ResolutionBlock > 0 && sdkCtx.BlockHeight() >= market.ResolutionBlock {
		return math.Int{}, types.ErrMarketNotOpen
	}

	pos, found := k.GetPosition(ctx, marketID, sender)
	if !found {
		return math.Int{}, types.ErrNoPosition
	}

	if outcome == "yes" {
		if pos.YesShares.LT(shares) {
			return math.Int{}, types.ErrInsufficientShares
		}
	} else {
		if pos.NoShares.LT(shares) {
			return math.Int{}, types.ErrInsufficientShares
		}
	}

	kProduct := market.YesShares.Mul(market.NoShares)

	var quoteReturned math.Int
	if outcome == "yes" {
		newYesShares := market.YesShares.Add(shares)
		newNoShares := kProduct.Quo(newYesShares)
		quoteReturned = market.NoShares.Sub(newNoShares)
		if quoteReturned.IsZero() || quoteReturned.IsNegative() {
			return math.Int{}, types.ErrInvalidAmount
		}
		market.YesShares = newYesShares
		market.NoShares = newNoShares
		pos.YesShares = pos.YesShares.Sub(shares)
	} else {
		newNoShares := market.NoShares.Add(shares)
		newYesShares := kProduct.Quo(newNoShares)
		quoteReturned = market.YesShares.Sub(newYesShares)
		if quoteReturned.IsZero() || quoteReturned.IsNegative() {
			return math.Int{}, types.ErrInvalidAmount
		}
		market.YesShares = newYesShares
		market.NoShares = newNoShares
		pos.NoShares = pos.NoShares.Sub(shares)
	}

	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return math.Int{}, err
	}
	coins := sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, quoteReturned))
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, coins); err != nil {
		return math.Int{}, err
	}

	k.SetMarket(ctx, market)
	k.SetPosition(ctx, pos)

	return quoteReturned, nil
}

// ExecuteResolveMarket resolves a market to yes, no, or void.
// Only the chain governance authority (k.authority) may resolve markets.
// The Resolver field stored on the market is informational only.
func (k Keeper) ExecuteResolveMarket(ctx context.Context, resolver string, marketID uint64, outcome string) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// C-5: only chain governance authority can resolve — prevents creator self-resolution
	if resolver != k.authority {
		return types.ErrUnauthorized
	}

	market, found := k.GetMarket(ctx, marketID)
	if !found {
		return types.ErrMarketNotFound
	}

	if market.Status == types.MarketStatusResolvedYes || market.Status == types.MarketStatusResolvedNo || market.Status == types.MarketStatusVoided {
		return types.ErrAlreadyResolved
	}

	switch outcome {
	case "yes":
		market.Status = types.MarketStatusResolvedYes
	case "no":
		market.Status = types.MarketStatusResolvedNo
	case "void":
		market.Status = types.MarketStatusVoided
	default:
		return types.ErrInvalidOutcome
	}

	// Store the resolved outcome on the market for use during claims
	market.Outcome = outcome

	// M-2: verify module solvency and compute payout ratio
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	balance := k.bankKeeper.GetBalance(ctx, moduleAddr, market.QuoteDenom).Amount

	// BUG-5 fix: market.YesShares / market.NoShares are AMM pool reserves, NOT
	// the total shares held by users.  We must iterate all positions to sum the
	// winning side's shares for an accurate payout ratio.
	var winnersTotal math.Int
	if outcome == "yes" || outcome == "no" {
		winnersTotal = math.ZeroInt()
		positions := k.GetPositionsForMarket(ctx, marketID)
		for _, pos := range positions {
			if outcome == "yes" {
				winnersTotal = winnersTotal.Add(pos.YesShares)
			} else {
				winnersTotal = winnersTotal.Add(pos.NoShares)
			}
		}
	}

	if outcome == "void" {
		// For void, calculate total refunds needed and apply insolvency ratio
		totalRefunds := math.ZeroInt()
		voidPositions := k.GetPositionsForMarket(ctx, marketID)
		for _, pos := range voidPositions {
			totalRefunds = totalRefunds.Add(pos.TotalSpent)
		}
		if totalRefunds.IsPositive() && balance.LT(totalRefunds) {
			market.PayoutRatio = math.LegacyNewDecFromInt(balance).Quo(math.LegacyNewDecFromInt(totalRefunds))
		} else {
			market.PayoutRatio = math.LegacyOneDec()
		}
	} else if winnersTotal.IsZero() {
		market.PayoutRatio = math.LegacyOneDec()
	} else if balance.LT(winnersTotal) {
		// Module is insolvent: cap payouts proportionally
		market.PayoutRatio = math.LegacyNewDecFromInt(balance).Quo(math.LegacyNewDecFromInt(winnersTotal))
	} else {
		market.PayoutRatio = math.LegacyOneDec()
	}

	k.SetMarket(ctx, market)
	k.SetResolution(ctx, types.Resolution{
		MarketID:   marketID,
		Outcome:    outcome,
		ResolvedAt: sdkCtx.BlockHeight(),
		ResolvedBy: resolver,
	})

	return nil
}

// ExecuteClaimWinnings claims payout from a resolved market.
// For resolved_yes: payout = position.YesShares (1 token per winning share)
// For resolved_no:  payout = position.NoShares
// For voided:       payout = position.TotalSpent (full refund)
func (k Keeper) ExecuteClaimWinnings(ctx context.Context, sender string, marketID uint64) (math.Int, error) {
	market, found := k.GetMarket(ctx, marketID)
	if !found {
		return math.Int{}, types.ErrMarketNotFound
	}

	if market.Status != types.MarketStatusResolvedYes && market.Status != types.MarketStatusResolvedNo && market.Status != types.MarketStatusVoided {
		return math.Int{}, types.ErrMarketNotOpen
	}

	pos, found := k.GetPosition(ctx, marketID, sender)
	if !found {
		return math.Int{}, types.ErrNoPosition
	}

	if pos.Claimed {
		return math.Int{}, types.ErrAlreadyClaimed
	}

	var payout math.Int
	switch market.Status {
	case types.MarketStatusResolvedYes:
		// M-2: apply payout ratio to cover potential insolvency
		rawShares := pos.YesShares
		if !market.PayoutRatio.IsNil() && !market.PayoutRatio.IsZero() {
			payout = market.PayoutRatio.MulInt(rawShares).TruncateInt()
		} else {
			payout = rawShares
		}
	case types.MarketStatusResolvedNo:
		// M-2: apply payout ratio to cover potential insolvency
		rawShares := pos.NoShares
		if !market.PayoutRatio.IsNil() && !market.PayoutRatio.IsZero() {
			payout = market.PayoutRatio.MulInt(rawShares).TruncateInt()
		} else {
			payout = rawShares
		}
	case types.MarketStatusVoided:
		// Apply payout ratio for void too — module may be insolvent
		rawRefund := pos.TotalSpent
		if !market.PayoutRatio.IsNil() && !market.PayoutRatio.IsZero() {
			payout = market.PayoutRatio.MulInt(rawRefund).TruncateInt()
		} else {
			payout = rawRefund
		}
	}

	if payout.IsZero() {
		return math.Int{}, types.ErrNoWinnings
	}

	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return math.Int{}, err
	}

	coins := sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, payout))
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, coins); err != nil {
		return math.Int{}, err
	}

	pos.Claimed = true
	k.SetPosition(ctx, pos)

	return payout, nil
}

// AutoCloseExpiredMarkets closes markets past their resolution block
func (k Keeper) AutoCloseExpiredMarkets(ctx context.Context) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	height := sdkCtx.BlockHeight()

	markets := k.GetAllMarkets(ctx)
	for _, market := range markets {
		if market.Status == types.MarketStatusOpen && market.ResolutionBlock > 0 && height >= market.ResolutionBlock {
			market.Status = types.MarketStatusClosed
			k.SetMarket(ctx, market)
		}
	}
}

// GetYesPrice returns the current YES price in basis points (0-10000)
// yesPrice = noShares / (yesShares + noShares)
func (k Keeper) GetYesPrice(ctx context.Context, marketID uint64) (math.Int, bool) {
	market, found := k.GetMarket(ctx, marketID)
	if !found {
		return math.Int{}, false
	}
	total := market.YesShares.Add(market.NoShares)
	if total.IsZero() {
		return math.ZeroInt(), true
	}
	price := market.NoShares.Mul(math.NewInt(10000)).Quo(total)
	return price, true
}

// GetNoPrice returns the current NO price in basis points (0-10000)
// noPrice = yesShares / (yesShares + noShares)
func (k Keeper) GetNoPrice(ctx context.Context, marketID uint64) (math.Int, bool) {
	market, found := k.GetMarket(ctx, marketID)
	if !found {
		return math.Int{}, false
	}
	total := market.YesShares.Add(market.NoShares)
	if total.IsZero() {
		return math.ZeroInt(), true
	}
	price := market.YesShares.Mul(math.NewInt(10000)).Quo(total)
	return price, true
}
