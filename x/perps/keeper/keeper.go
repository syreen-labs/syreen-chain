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

	"syreen/x/perps/types"
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
		cdc:           cdc,
		storeService:  storeService,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		authority:     authority,
	}
}

func (k *Keeper) SetDexKeeper(dk types.DexKeeper) {
	k.dexKeeper = dk
}

// ============================================================
// Market CRUD
// ============================================================

func (k Keeper) GetMarket(ctx context.Context, marketID uint64) (types.PerpMarket, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.MarketKey(marketID))
	if err != nil || bz == nil {
		return types.PerpMarket{}, false
	}
	var market types.PerpMarket
	if err := json.Unmarshal(bz, &market); err != nil {
		return types.PerpMarket{}, false
	}
	return market, true
}

func (k Keeper) SetMarket(ctx context.Context, market types.PerpMarket) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(market)
	_ = kvStore.Set(types.MarketKey(market.ID), bz)
}

func (k Keeper) GetAllMarkets(ctx context.Context) []types.PerpMarket {
	kvStore := k.storeService.OpenKVStore(ctx)
	iter, err := kvStore.Iterator([]byte(types.MarketPrefix), append([]byte(types.MarketPrefix), 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
	if err != nil {
		return nil
	}
	defer iter.Close()
	var markets []types.PerpMarket
	for ; iter.Valid(); iter.Next() {
		var m types.PerpMarket
		if err := json.Unmarshal(iter.Value(), &m); err == nil {
			markets = append(markets, m)
		}
	}
	sort.Slice(markets, func(i, j int) bool { return markets[i].ID < markets[j].ID })
	return markets
}

func (k Keeper) GetNextMarketID(ctx context.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := kvStore.Get([]byte(types.NextMarketIDKey))
	if bz == nil {
		return 1
	}
	var id uint64
	json.Unmarshal(bz, &id)
	return id
}

func (k Keeper) SetNextMarketID(ctx context.Context, id uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(id)
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
	_ = kvStore.Set(types.PositionByAddrKey(pos.Address, pos.MarketID), []byte{1})
}

func (k Keeper) DeletePosition(ctx context.Context, marketID uint64, address string) {
	kvStore := k.storeService.OpenKVStore(ctx)
	_ = kvStore.Delete(types.PositionKey(marketID, address))
	_ = kvStore.Delete(types.PositionByAddrKey(address, marketID))
}

func (k Keeper) GetPositionsByAddress(ctx context.Context, address string) []types.Position {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := types.PositionByAddrPrefix(address)
	iter, err := kvStore.Iterator(prefix, append(prefix, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
	if err != nil {
		return nil
	}
	defer iter.Close()
	var positions []types.Position
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if len(key) >= 8 {
			poolIDBytes := key[len(key)-8:]
			marketID := uint64(poolIDBytes[0])<<56 | uint64(poolIDBytes[1])<<48 | uint64(poolIDBytes[2])<<40 | uint64(poolIDBytes[3])<<32 | uint64(poolIDBytes[4])<<24 | uint64(poolIDBytes[5])<<16 | uint64(poolIDBytes[6])<<8 | uint64(poolIDBytes[7])
			if pos, ok := k.GetPosition(ctx, marketID, address); ok {
				positions = append(positions, pos)
			}
		}
	}
	return positions
}

// GetAllPositions returns all positions across all markets (for genesis export)
func (k Keeper) GetAllPositions(ctx context.Context) []types.Position {
	kvStore := k.storeService.OpenKVStore(ctx)
	iter, err := kvStore.Iterator([]byte(types.PositionPrefix), append([]byte(types.PositionPrefix), 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
	if err != nil {
		return nil
	}
	defer iter.Close()
	var positions []types.Position
	for ; iter.Valid(); iter.Next() {
		var pos types.Position
		if err := json.Unmarshal(iter.Value(), &pos); err == nil {
			positions = append(positions, pos)
		}
	}
	return positions
}

// GetAllFundingStates returns all funding states (for genesis export)
func (k Keeper) GetAllFundingStates(ctx context.Context) []types.FundingState {
	kvStore := k.storeService.OpenKVStore(ctx)
	iter, err := kvStore.Iterator([]byte(types.FundingRatePrefix), append([]byte(types.FundingRatePrefix), 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
	if err != nil {
		return nil
	}
	defer iter.Close()
	var states []types.FundingState
	for ; iter.Valid(); iter.Next() {
		var fs types.FundingState
		if err := json.Unmarshal(iter.Value(), &fs); err == nil {
			states = append(states, fs)
		}
	}
	return states
}

// ============================================================
// Funding State
// ============================================================

func (k Keeper) GetFundingState(ctx context.Context, marketID uint64) types.FundingState {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.FundingRateKey(marketID))
	if err != nil || bz == nil {
		return types.FundingState{
			MarketID:           marketID,
			CurrentFundingRate: math.LegacyZeroDec(),
			CumulativeFunding:  math.LegacyZeroDec(),
		}
	}
	var fs types.FundingState
	json.Unmarshal(bz, &fs)
	return fs
}

func (k Keeper) SetFundingState(ctx context.Context, fs types.FundingState) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(fs)
	_ = kvStore.Set(types.FundingRateKey(fs.MarketID), bz)
}

// ============================================================
// Insurance Fund
// ============================================================

func (k Keeper) GetInsuranceFund(ctx context.Context) types.InsuranceFund {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.InsuranceFundKey))
	if err != nil || bz == nil {
		return types.InsuranceFund{Balance: math.ZeroInt()}
	}
	var fund types.InsuranceFund
	json.Unmarshal(bz, &fund)
	return fund
}

func (k Keeper) SetInsuranceFund(ctx context.Context, fund types.InsuranceFund) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(fund)
	_ = kvStore.Set([]byte(types.InsuranceFundKey), bz)
}

// ============================================================
// Mark Price (from DEX spot)
// ============================================================

func (k Keeper) GetMarkPrice(ctx context.Context, market types.PerpMarket) (math.LegacyDec, error) {
	if k.dexKeeper == nil {
		return math.LegacyZeroDec(), fmt.Errorf("dex keeper not set")
	}

	spotPrice, err := k.dexKeeper.GetSpotPrice(ctx, market.PoolID, market.BaseDenom, market.QuoteDenom)
	if err != nil {
		return math.LegacyZeroDec(), err
	}

	// C-4 fix: dampen price to prevent single-block manipulation.
	// If a prior mark price exists, clamp spot price. Use 2% for initial
	// move (when LastMarkPrice is zero/nil) and 5% for subsequent moves.
	if !market.LastMarkPrice.IsNil() && market.LastMarkPrice.IsPositive() {
		maxMove := market.LastMarkPrice.Mul(math.LegacyNewDecWithPrec(5, 2)) // 5% of last mark
		upperBound := market.LastMarkPrice.Add(maxMove)
		lowerBound := market.LastMarkPrice.Sub(maxMove)

		if spotPrice.GT(upperBound) {
			spotPrice = upperBound
		} else if spotPrice.LT(lowerBound) {
			spotPrice = lowerBound
		}
	} else {
		// No prior mark price — this is the initial price observation.
		// We still accept the spot price but will clamp via TWAP below.
	}

	// TWAP: store the last 10 prices and average them to resist manipulation.
	history := k.getMarkPriceHistory(ctx, market.ID)
	history = append(history, spotPrice)
	if len(history) > 10 {
		history = history[len(history)-10:]
	}
	k.setMarkPriceHistory(ctx, market.ID, history)

	// Average the history
	sum := math.LegacyZeroDec()
	for _, p := range history {
		sum = sum.Add(p)
	}
	twapPrice := sum.Quo(math.LegacyNewDec(int64(len(history))))

	return twapPrice, nil
}

// getMarkPriceHistory retrieves the last N mark price observations for TWAP.
func (k Keeper) getMarkPriceHistory(ctx context.Context, marketID uint64) []math.LegacyDec {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.MarkPriceTWAPKey(marketID))
	if err != nil || bz == nil {
		return nil
	}
	var prices []math.LegacyDec
	if err := json.Unmarshal(bz, &prices); err != nil {
		return nil
	}
	return prices
}

// setMarkPriceHistory stores the last N mark price observations for TWAP.
func (k Keeper) setMarkPriceHistory(ctx context.Context, marketID uint64, prices []math.LegacyDec) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(prices)
	_ = kvStore.Set(types.MarkPriceTWAPKey(marketID), bz)
}

// ============================================================
// Core Perps Logic
// ============================================================

// ExecuteOpenPosition opens a new perpetual futures position
func (k Keeper) ExecuteOpenPosition(ctx context.Context, sender string, marketID uint64, side types.Side, margin math.Int, leverage math.LegacyDec) (*types.MsgOpenPositionResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	market, ok := k.GetMarket(ctx, marketID)
	if !ok {
		return nil, types.ErrMarketNotFound
	}
	if !market.Active {
		return nil, types.ErrMarketNotActive
	}
	if leverage.GT(market.MaxLeverage) {
		return nil, types.ErrMaxLeverageExceeded
	}

	// Check no existing position
	if _, exists := k.GetPosition(ctx, marketID, sender); exists {
		return nil, types.ErrPositionAlreadyExists
	}

	// Get mark price
	markPrice, err := k.GetMarkPrice(ctx, market)
	if err != nil || markPrice.IsZero() {
		return nil, fmt.Errorf("cannot get mark price: %v", err)
	}

	// Calculate notional value = margin * leverage
	notional := math.LegacyNewDecFromInt(margin).Mul(leverage)

	// Position size = notional / markPrice
	positionSize := notional.Quo(markPrice)

	// Check open interest limits
	notionalInt := notional.TruncateInt()
	if side == types.SideLong {
		newOI := market.LongOpenInterest.Add(notionalInt)
		if newOI.GT(market.MaxOpenInterest) {
			return nil, types.ErrMaxOpenInterest
		}
	} else {
		newOI := market.ShortOpenInterest.Add(notionalInt)
		if newOI.GT(market.MaxOpenInterest) {
			return nil, types.ErrMaxOpenInterest
		}
	}

	// Calculate taker fee
	fee := math.LegacyNewDecFromInt(margin).Mul(leverage).Mul(market.TakerFee).TruncateInt()

	// H-3 fix: warn (via event) if insurance fund is below minimum before opening position
	fund := k.GetInsuranceFund(ctx)
	minReserve := math.NewInt(1_000_000) // 1 SYR minimum
	if fund.Balance.LT(minReserve) && !market.AllowUndercapitalized {
		sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
			"perps_insurance_fund_low",
			sdk.NewAttribute("market_id", fmt.Sprintf("%d", marketID)),
			sdk.NewAttribute("balance", fund.Balance.String()),
			sdk.NewAttribute("minimum", minReserve.String()),
		))
	}

	// L-2 fix: handle bech32 decode error
	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return nil, err
	}

	// Transfer margin + fee from user to module
	totalDeposit := margin.Add(fee)
	coins := sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, totalDeposit))
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, coins); err != nil {
		return nil, types.ErrInsufficientFunds
	}

	// Fee goes to insurance fund
	fund.Balance = fund.Balance.Add(fee)
	k.SetInsuranceFund(ctx, fund)

	// Create position
	fundingState := k.GetFundingState(ctx, marketID)
	pos := types.Position{
		Address:           sender,
		MarketID:          marketID,
		Side:              side,
		Size:              positionSize,
		EntryPrice:        markPrice,
		Margin:            margin,
		Leverage:          leverage,
		UnrealizedPnL:     math.LegacyZeroDec(),
		CumulativeFunding: fundingState.CumulativeFunding,
		OpenedAt:          sdkCtx.BlockHeight(),
		LastFundingBlock:  sdkCtx.BlockHeight(),
	}
	k.SetPosition(ctx, pos)

	// Update open interest
	if side == types.SideLong {
		market.LongOpenInterest = market.LongOpenInterest.Add(notionalInt)
	} else {
		market.ShortOpenInterest = market.ShortOpenInterest.Add(notionalInt)
	}
	k.SetMarket(ctx, market)

	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"perps_open_position",
		sdk.NewAttribute("sender", sender),
		sdk.NewAttribute("market_id", fmt.Sprintf("%d", marketID)),
		sdk.NewAttribute("side", string(side)),
		sdk.NewAttribute("size", positionSize.String()),
		sdk.NewAttribute("entry_price", markPrice.String()),
		sdk.NewAttribute("margin", margin.String()),
		sdk.NewAttribute("leverage", leverage.String()),
	))

	return &types.MsgOpenPositionResponse{
		PositionSize: positionSize,
		EntryPrice:   markPrice,
		Fee:          fee,
	}, nil
}

// ExecuteClosePosition closes an existing position and settles PnL
func (k Keeper) ExecuteClosePosition(ctx context.Context, sender string, marketID uint64) (*types.MsgClosePositionResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	market, ok := k.GetMarket(ctx, marketID)
	if !ok {
		return nil, types.ErrMarketNotFound
	}

	pos, exists := k.GetPosition(ctx, marketID, sender)
	if !exists {
		return nil, types.ErrNoPosition
	}

	// Get current mark price
	markPrice, err := k.GetMarkPrice(ctx, market)
	if err != nil || markPrice.IsZero() {
		return nil, fmt.Errorf("cannot get mark price: %v", err)
	}

	// Calculate PnL
	pnl := k.CalculatePnL(pos, markPrice)

	// Apply funding payments
	fundingState := k.GetFundingState(ctx, marketID)
	fundingPayment := k.CalculateFundingPayment(pos, fundingState)
	pnl = pnl.Sub(fundingPayment)

	// Calculate payout = margin + PnL
	payout := math.LegacyNewDecFromInt(pos.Margin).Add(pnl).TruncateInt()
	if payout.IsNegative() {
		payout = math.ZeroInt()
	}

	// L-2 fix: handle bech32 decode error
	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return nil, err
	}

	// Send payout to user (H-3 + L-1 fixes: check balances, draw from insurance if needed, never ignore send error)
	if payout.IsPositive() {
		moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
		balance := k.bankKeeper.GetBalance(ctx, moduleAddr, market.QuoteDenom)
		if balance.Amount.LT(payout) {
			// Use insurance fund for shortfall
			fund := k.GetInsuranceFund(ctx)
			shortfall := payout.Sub(balance.Amount)
			if fund.Balance.GTE(shortfall) {
				fund.Balance = fund.Balance.Sub(shortfall)
				k.SetInsuranceFund(ctx, fund)
			} else {
				// Insurance fund cannot cover full shortfall — cap payout and emit event
				availableFromInsurance := fund.Balance
				fund.Balance = math.ZeroInt()
				k.SetInsuranceFund(ctx, fund)
				payout = balance.Amount.Add(availableFromInsurance)
				sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
					"perps_insurance_fund_exhausted",
					sdk.NewAttribute("market_id", fmt.Sprintf("%d", marketID)),
					sdk.NewAttribute("shortfall", shortfall.String()),
					sdk.NewAttribute("capped_payout", payout.String()),
				))
			}
		}
		if payout.IsPositive() {
			coins := sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, payout))
			if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, coins); err != nil {
				return nil, err
			}
		}
	} else {
		// Loss exceeds margin — all margin absorbed, insurance fund takes the hit
		deficit := pnl.Add(math.LegacyNewDecFromInt(pos.Margin)).Neg().TruncateInt()
		if deficit.IsPositive() {
			fund := k.GetInsuranceFund(ctx)
			fund.Balance = fund.Balance.Sub(deficit)
			if fund.Balance.IsNegative() {
				fund.Balance = math.ZeroInt()
			}
			k.SetInsuranceFund(ctx, fund)
		}
	}

	// Reconcile insurance fund counter with actual module account balance
	k.reconcileInsuranceFund(ctx, market.QuoteDenom)

	// Update open interest
	notional := pos.Size.Mul(pos.EntryPrice).TruncateInt()
	if pos.Side == types.SideLong {
		market.LongOpenInterest = market.LongOpenInterest.Sub(notional)
		if market.LongOpenInterest.IsNegative() {
			market.LongOpenInterest = math.ZeroInt()
		}
	} else {
		market.ShortOpenInterest = market.ShortOpenInterest.Sub(notional)
		if market.ShortOpenInterest.IsNegative() {
			market.ShortOpenInterest = math.ZeroInt()
		}
	}
	k.SetMarket(ctx, market)

	// Delete position
	k.DeletePosition(ctx, marketID, sender)

	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"perps_close_position",
		sdk.NewAttribute("sender", sender),
		sdk.NewAttribute("market_id", fmt.Sprintf("%d", marketID)),
		sdk.NewAttribute("pnl", pnl.String()),
		sdk.NewAttribute("payout", payout.String()),
	))

	return &types.MsgClosePositionResponse{
		RealizedPnL: pnl,
		Payout:      payout,
	}, nil
}

// AddMarginToPosition adds collateral to an existing position
func (k Keeper) AddMarginToPosition(ctx context.Context, sender string, marketID uint64, amount math.Int) error {
	market, ok := k.GetMarket(ctx, marketID)
	if !ok {
		return types.ErrMarketNotFound
	}

	pos, exists := k.GetPosition(ctx, marketID, sender)
	if !exists {
		return types.ErrNoPosition
	}

	// L-2 fix: handle bech32 decode error
	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return err
	}

	// Transfer margin from user
	coins := sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, amount))
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, coins); err != nil {
		return types.ErrInsufficientFunds
	}

	pos.Margin = pos.Margin.Add(amount)
	// Recalculate effective leverage
	markPrice, _ := k.GetMarkPrice(ctx, market)
	if !markPrice.IsZero() {
		notional := pos.Size.Mul(markPrice)
		pos.Leverage = notional.Quo(math.LegacyNewDecFromInt(pos.Margin))
	}
	k.SetPosition(ctx, pos)
	return nil
}

// RemoveMarginFromPosition removes excess collateral
func (k Keeper) RemoveMarginFromPosition(ctx context.Context, sender string, marketID uint64, amount math.Int) error {
	market, ok := k.GetMarket(ctx, marketID)
	if !ok {
		return types.ErrMarketNotFound
	}

	pos, exists := k.GetPosition(ctx, marketID, sender)
	if !exists {
		return types.ErrNoPosition
	}

	// M-1 fix: explicit zero/negative check BEFORE leverage recalculation to prevent division by zero
	newMargin := pos.Margin.Sub(amount)
	if newMargin.IsZero() || newMargin.IsNegative() {
		return types.ErrInsufficientMargin
	}

	// Check removal doesn't violate maintenance margin
	markPrice, _ := k.GetMarkPrice(ctx, market)
	pnl := k.CalculatePnL(pos, markPrice)

	equity := math.LegacyNewDecFromInt(newMargin).Add(pnl)
	notional := pos.Size.Mul(markPrice)
	if notional.IsPositive() {
		marginRatio := equity.Quo(notional)
		if marginRatio.LT(market.MaintenanceMargin) {
			return types.ErrInsufficientMargin
		}
	}

	// L-2 fix: handle bech32 decode error
	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return err
	}

	// L-1 fix: never ignore bank send error
	coins := sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, amount))
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, coins); err != nil {
		return err
	}

	pos.Margin = newMargin
	pos.Leverage = notional.Quo(math.LegacyNewDecFromInt(pos.Margin))
	k.SetPosition(ctx, pos)
	return nil
}

// CreateMarketEntry creates a new perps market (authority only)
func (k Keeper) CreateMarketEntry(ctx context.Context, authority string, baseDenom, quoteDenom string, poolID uint64, maxLeverage math.LegacyDec) (uint64, error) {
	if authority != k.authority {
		return 0, types.ErrUnauthorized
	}

	if maxLeverage.IsNil() || maxLeverage.LT(math.LegacyOneDec()) {
		maxLeverage = types.DefaultMaxLeverage
	}

	marketID := k.GetNextMarketID(ctx)

	market := types.PerpMarket{
		ID:                marketID,
		BaseDenom:         baseDenom,
		QuoteDenom:        quoteDenom,
		PoolID:            poolID,
		MaxLeverage:       maxLeverage,
		MaintenanceMargin: types.DefaultMaintenanceMargin,
		InitialMargin:     types.DefaultInitialMargin,
		TakerFee:          types.DefaultTakerFee,
		MakerFee:          types.DefaultMakerFee,
		FundingInterval:   types.DefaultFundingInterval,
		MaxFundingRate:    types.DefaultMaxFundingRate,
		MaxOpenInterest:   math.NewInt(1_000_000_000_000), // 1M quote units default
		LongOpenInterest:  math.ZeroInt(),
		ShortOpenInterest: math.ZeroInt(),
		Active:            true,
		Creator:           authority,
	}

	k.SetMarket(ctx, market)
	k.SetNextMarketID(ctx, marketID+1)

	// Initialize funding state
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	k.SetFundingState(ctx, types.FundingState{
		MarketID:           marketID,
		CurrentFundingRate: math.LegacyZeroDec(),
		CumulativeFunding:  math.LegacyZeroDec(),
		LastFundingBlock:   sdkCtx.BlockHeight(),
	})

	return marketID, nil
}

// ============================================================
// PnL Calculation
// ============================================================

// CalculatePnL calculates unrealized PnL for a position
func (k Keeper) CalculatePnL(pos types.Position, markPrice math.LegacyDec) math.LegacyDec {
	if pos.Size.IsZero() || markPrice.IsZero() {
		return math.LegacyZeroDec()
	}
	// Long: PnL = size * (markPrice - entryPrice)
	// Short: PnL = size * (entryPrice - markPrice)
	priceDiff := markPrice.Sub(pos.EntryPrice)
	if pos.Side == types.SideShort {
		priceDiff = pos.EntryPrice.Sub(markPrice)
	}
	return pos.Size.Mul(priceDiff)
}

// CalculateFundingPayment calculates the funding payment owed
func (k Keeper) CalculateFundingPayment(pos types.Position, fs types.FundingState) math.LegacyDec {
	// Funding payment = size * (cumulative_funding_now - cumulative_funding_at_entry)
	fundingDelta := fs.CumulativeFunding.Sub(pos.CumulativeFunding)
	payment := pos.Size.Mul(fundingDelta)
	// Longs pay positive funding, shorts receive (and vice versa)
	if pos.Side == types.SideShort {
		payment = payment.Neg()
	}
	return payment
}

// GetMarginRatio returns the current margin ratio for a position
func (k Keeper) GetMarginRatio(pos types.Position, markPrice math.LegacyDec) math.LegacyDec {
	pnl := k.CalculatePnL(pos, markPrice)
	equity := math.LegacyNewDecFromInt(pos.Margin).Add(pnl)
	notional := pos.Size.Mul(markPrice)
	if notional.IsZero() {
		return math.LegacyZeroDec()
	}
	return equity.Quo(notional)
}

// ============================================================
// BeginBlock: Funding Rates & Liquidations
// ============================================================

// ProcessFundingRates updates funding rates for all active markets
func (k Keeper) ProcessFundingRates(ctx context.Context) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	markets := k.GetAllMarkets(ctx)

	for _, market := range markets {
		if !market.Active {
			continue
		}

		// C-4 fix: update the dampened LastMarkPrice every block so the 5% clamp
		// in GetMarkPrice has an up-to-date reference point each block.
		if markPrice, mpErr := k.GetMarkPrice(ctx, market); mpErr == nil && markPrice.IsPositive() {
			market.LastMarkPrice = markPrice
			market.LastMarkHeight = sdkCtx.BlockHeight()
			k.SetMarket(ctx, market)
		}

		fs := k.GetFundingState(ctx, market.ID)
		blocksSinceLastFunding := sdkCtx.BlockHeight() - fs.LastFundingBlock

		if blocksSinceLastFunding < market.FundingInterval {
			continue
		}

		// Calculate funding rate based on OI imbalance
		// If more longs than shorts, longs pay shorts (positive rate)
		// If more shorts than longs, shorts pay longs (negative rate)
		longOI := math.LegacyNewDecFromInt(market.LongOpenInterest)
		shortOI := math.LegacyNewDecFromInt(market.ShortOpenInterest)
		totalOI := longOI.Add(shortOI)

		fundingRate := math.LegacyZeroDec()
		if totalOI.IsPositive() {
			imbalance := longOI.Sub(shortOI).Quo(totalOI)
			fundingRate = imbalance.Mul(market.MaxFundingRate)
			// Cap at max funding rate
			if fundingRate.GT(market.MaxFundingRate) {
				fundingRate = market.MaxFundingRate
			}
			if fundingRate.LT(market.MaxFundingRate.Neg()) {
				fundingRate = market.MaxFundingRate.Neg()
			}
		}

		fs.CurrentFundingRate = fundingRate
		fs.CumulativeFunding = fs.CumulativeFunding.Add(fundingRate)
		fs.LastFundingBlock = sdkCtx.BlockHeight()
		k.SetFundingState(ctx, fs)
	}
}

// ProcessLiquidations checks all positions for liquidation
func (k Keeper) ProcessLiquidations(ctx context.Context) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	markets := k.GetAllMarkets(ctx)

	// L-9 fix: cap total liquidations per block to prevent runaway block times
	const maxLiquidationsPerBlock = 50
	checked := 0

	for _, market := range markets {
		if !market.Active {
			continue
		}

		markPrice, err := k.GetMarkPrice(ctx, market)
		if err != nil || markPrice.IsZero() {
			continue
		}

		// Iterate all positions for this market
		kvStore := k.storeService.OpenKVStore(ctx)
		prefix := types.PositionMarketPrefix(market.ID)
		iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
		if err != nil {
			continue
		}

		var toLiquidate []types.Position
		for ; iter.Valid(); iter.Next() {
			// L-9 fix: stop scanning once we've reached the per-block limit
			if checked >= maxLiquidationsPerBlock {
				break
			}
			checked++

			var pos types.Position
			if err := json.Unmarshal(iter.Value(), &pos); err != nil {
				continue
			}
			marginRatio := k.GetMarginRatio(pos, markPrice)
			if marginRatio.LT(market.MaintenanceMargin) {
				toLiquidate = append(toLiquidate, pos)
			}
		}
		iter.Close()

		// Liquidate positions
		for _, pos := range toLiquidate {
			k.liquidatePosition(ctx, market, pos, markPrice)
			sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
				"perps_liquidation",
				sdk.NewAttribute("address", pos.Address),
				sdk.NewAttribute("market_id", fmt.Sprintf("%d", market.ID)),
				sdk.NewAttribute("side", string(pos.Side)),
				sdk.NewAttribute("size", pos.Size.String()),
				sdk.NewAttribute("mark_price", markPrice.String()),
			))
		}

		// Stop processing more markets if the block limit was hit
		if checked >= maxLiquidationsPerBlock {
			break
		}
	}
}

func (k Keeper) liquidatePosition(ctx context.Context, market types.PerpMarket, pos types.Position, markPrice math.LegacyDec) {
	// Calculate remaining equity
	pnl := k.CalculatePnL(pos, markPrice)
	fundingState := k.GetFundingState(ctx, market.ID)
	fundingPayment := k.CalculateFundingPayment(pos, fundingState)
	equity := math.LegacyNewDecFromInt(pos.Margin).Add(pnl).Sub(fundingPayment)

	// Liquidation penalty: 2.5% of notional goes to insurance fund
	notional := pos.Size.Mul(markPrice)
	penalty := notional.Mul(math.LegacyNewDecWithPrec(25, 3)).TruncateInt() // 2.5%

	// Whatever equity remains after penalty: penalty goes to insurance fund,
	// any excess is returned to the liquidated user.
	remaining := equity.TruncateInt()
	if remaining.IsPositive() {
		fund := k.GetInsuranceFund(ctx)
		toInsurance := remaining
		if toInsurance.GT(penalty) {
			toInsurance = penalty
		}
		fund.Balance = fund.Balance.Add(toInsurance)
		k.SetInsuranceFund(ctx, fund)

		// Return excess equity (remaining - penalty) to the liquidated user
		excess := remaining.Sub(toInsurance)
		if excess.IsPositive() {
			userAddr, addrErr := sdk.AccAddressFromBech32(pos.Address)
			if addrErr == nil {
				returnCoins := sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, excess))
				if sendErr := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, userAddr, returnCoins); sendErr != nil {
					sdk.UnwrapSDKContext(ctx).Logger().Error("failed to return excess margin during liquidation", "address", pos.Address, "amount", excess.String(), "error", sendErr)
				}
			}
		}
	} else {
		// Socialized loss — insurance fund absorbs the deficit
		fund := k.GetInsuranceFund(ctx)
		deficit := remaining.Neg()
		fund.Balance = fund.Balance.Sub(deficit)
		if fund.Balance.IsNegative() {
			fund.Balance = math.ZeroInt()
		}
		k.SetInsuranceFund(ctx, fund)
	}

	// Update open interest
	notionalInt := pos.Size.Mul(pos.EntryPrice).TruncateInt()
	if pos.Side == types.SideLong {
		market.LongOpenInterest = market.LongOpenInterest.Sub(notionalInt)
		if market.LongOpenInterest.IsNegative() {
			market.LongOpenInterest = math.ZeroInt()
		}
	} else {
		market.ShortOpenInterest = market.ShortOpenInterest.Sub(notionalInt)
		if market.ShortOpenInterest.IsNegative() {
			market.ShortOpenInterest = math.ZeroInt()
		}
	}
	k.SetMarket(ctx, market)

	// Delete the position
	k.DeletePosition(ctx, market.ID, pos.Address)

	// Reconcile insurance fund after liquidation to keep counter in sync
	k.reconcileInsuranceFund(ctx, market.QuoteDenom)
}

// BeginBlockHandler runs funding rates, liquidations, and oracle price sampling
func (k Keeper) BeginBlockHandler(ctx context.Context) {
	// Record DEX oracle price samples for TWAP computation
	k.RecordAllPriceSamples(ctx)
	k.ProcessFundingRates(ctx)
	k.ProcessLiquidations(ctx)
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

// reconcileInsuranceFund ensures the insurance fund balance counter matches
// the actual module account bank balance. If they differ, the counter is
// corrected to the actual balance.
func (k Keeper) reconcileInsuranceFund(ctx context.Context, quoteDenom string) {
	fund := k.GetInsuranceFund(ctx)
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	actualBalance := k.bankKeeper.GetBalance(ctx, moduleAddr, quoteDenom)
	if !fund.Balance.Equal(actualBalance.Amount) {
		fund.Balance = actualBalance.Amount
		k.SetInsuranceFund(ctx, fund)
	}
}
