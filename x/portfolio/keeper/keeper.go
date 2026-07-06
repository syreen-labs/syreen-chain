package keeper

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sort"

	"cosmossdk.io/core/store"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/portfolio/types"
)

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
// Portfolio CRUD
// ============================================================

func (k Keeper) GetPortfolio(ctx context.Context, address string) (types.Portfolio, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.PortfolioKey(address))
	if err != nil || bz == nil {
		return types.Portfolio{}, false
	}
	var p types.Portfolio
	if err := json.Unmarshal(bz, &p); err != nil {
		return types.Portfolio{}, false
	}
	return p, true
}

func (k Keeper) SetPortfolio(ctx context.Context, p types.Portfolio) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(p)
	_ = kvStore.Set(types.PortfolioKey(p.Address), bz)
}

// GetOrCreatePortfolio returns existing or creates a new empty portfolio
func (k Keeper) GetOrCreatePortfolio(ctx context.Context, address string) types.Portfolio {
	p, found := k.GetPortfolio(ctx, address)
	if found {
		return p
	}
	return types.Portfolio{
		Address:       address,
		TotalValue:    math.ZeroInt(),
		TotalPnL:      math.ZeroInt(),
		RealizedPnL:   math.ZeroInt(),
		UnrealizedPnL: math.ZeroInt(),
		CostBasis:     math.ZeroInt(),
		LastUpdated:   0,
	}
}

// ============================================================
// Portfolio Asset CRUD
// ============================================================

func (k Keeper) GetPortfolioAsset(ctx context.Context, address, denom string) (types.PortfolioAsset, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.PortfolioAssetKey(address, denom))
	if err != nil || bz == nil {
		return types.PortfolioAsset{}, false
	}
	var pa types.PortfolioAsset
	if err := json.Unmarshal(bz, &pa); err != nil {
		return types.PortfolioAsset{}, false
	}
	return pa, true
}

func (k Keeper) SetPortfolioAsset(ctx context.Context, pa types.PortfolioAsset) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(pa)
	_ = kvStore.Set(types.PortfolioAssetKey(pa.Address, pa.Denom), bz)
}

func (k Keeper) GetAllPortfolioAssets(ctx context.Context, address string) []types.PortfolioAsset {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := types.PortfolioAssetPrefixKey(address)
	iter, err := kvStore.Iterator(prefix, append(prefix, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var assets []types.PortfolioAsset
	for ; iter.Valid(); iter.Next() {
		var pa types.PortfolioAsset
		if err := json.Unmarshal(iter.Value(), &pa); err == nil {
			assets = append(assets, pa)
		}
	}
	return assets
}

// ============================================================
// Trade Record CRUD
// ============================================================

func (k Keeper) GetNextTradeID(ctx context.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.NextTradeIDKey))
	if err != nil || bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetNextTradeID(ctx context.Context, id uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	_ = kvStore.Set([]byte(types.NextTradeIDKey), bz)
}

func (k Keeper) GetTradeRecord(ctx context.Context, id uint64) (types.TradeRecord, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.TradeRecordKey(id))
	if err != nil || bz == nil {
		return types.TradeRecord{}, false
	}
	var tr types.TradeRecord
	if err := json.Unmarshal(bz, &tr); err != nil {
		return types.TradeRecord{}, false
	}
	return tr, true
}

func (k Keeper) SetTradeRecord(ctx context.Context, tr types.TradeRecord) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(tr)
	_ = kvStore.Set(types.TradeRecordKey(tr.ID), bz)
	// Index by address
	_ = kvStore.Set(types.TradeByAddrKey(tr.Address, tr.ID), []byte{1})
}

// RecordTradeAndUpdatePortfolio creates a trade record and updates the portfolio
func (k Keeper) RecordTradeAndUpdatePortfolio(ctx context.Context, address, tradeType, denom string, amount, price, pnl math.Int, blockHeight int64) (uint64, error) {
	tradeID := k.GetNextTradeID(ctx)
	tr := types.TradeRecord{
		ID:        tradeID,
		Address:   address,
		TradeType: tradeType,
		Denom:     denom,
		Amount:    amount,
		Price:     price,
		PnL:       pnl,
		Block:     blockHeight,
	}
	k.SetTradeRecord(ctx, tr)
	k.SetNextTradeID(ctx, tradeID+1)

	// Update portfolio
	portfolio := k.GetOrCreatePortfolio(ctx, address)
	portfolio.RealizedPnL = portfolio.RealizedPnL.Add(pnl)
	// Accumulate cost basis: add amount paid in (amountIn) for buys/deposits
	switch tradeType {
	case types.TradeTypeSwap, types.TradeTypePerpOpen, types.TradeTypeLendingDeposit:
		if portfolio.CostBasis.IsNil() {
			portfolio.CostBasis = math.ZeroInt()
		}
		portfolio.CostBasis = portfolio.CostBasis.Add(amount)
	}
	portfolio.TotalPnL = portfolio.RealizedPnL.Add(portfolio.UnrealizedPnL)
	portfolio.LastUpdated = blockHeight
	k.SetPortfolio(ctx, portfolio)

	// Update asset entry
	asset, found := k.GetPortfolioAsset(ctx, address, denom)
	if !found {
		asset = types.PortfolioAsset{
			Address:        address,
			Denom:          denom,
			Amount:         math.ZeroInt(),
			Value:          math.ZeroInt(),
			PnL:            math.ZeroInt(),
			PercentOfTotal: math.LegacyZeroDec(),
		}
	}

	// For buys/deposits, add amount; for sells/withdraws, subtract
	switch tradeType {
	case types.TradeTypeSwap, types.TradeTypePerpOpen, types.TradeTypeLendingDeposit:
		asset.Amount = asset.Amount.Add(amount)
	case types.TradeTypePerpClose, types.TradeTypeLendingWithdraw:
		asset.Amount = asset.Amount.Sub(amount)
		if asset.Amount.IsNegative() {
			asset.Amount = math.ZeroInt()
		}
	}
	asset.Value = asset.Amount.Mul(price)
	asset.PnL = asset.PnL.Add(pnl)
	k.SetPortfolioAsset(ctx, asset)

	// Check if whale activity
	valueInSyr := amount.Mul(price)
	if valueInSyr.GTE(math.NewInt(types.WhaleThreshold)) {
		activityType := types.ActivityLargeSwap
		switch tradeType {
		case types.TradeTypePerpOpen:
			activityType = types.ActivityPositionOpened
		case types.TradeTypePerpClose:
			activityType = types.ActivityPositionClosed
		case types.TradeTypeLendingDeposit:
			activityType = types.ActivityLargeDeposit
		}
		desc := fmt.Sprintf("%s %s %s (value: %s)", tradeType, amount.String(), denom, valueInSyr.String())
		k.RecordActivity(ctx, activityType, address, desc, amount, denom, blockHeight)
	}

	return tradeID, nil
}

// ============================================================
// Activity Feed CRUD
// ============================================================

func (k Keeper) GetNextActivityID(ctx context.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.NextActivityIDKey))
	if err != nil || bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetNextActivityID(ctx context.Context, id uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	_ = kvStore.Set([]byte(types.NextActivityIDKey), bz)
}

func (k Keeper) GetActivity(ctx context.Context, id uint64) (types.Activity, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.ActivityKey(id))
	if err != nil || bz == nil {
		return types.Activity{}, false
	}
	var a types.Activity
	if err := json.Unmarshal(bz, &a); err != nil {
		return types.Activity{}, false
	}
	return a, true
}

func (k Keeper) SetActivity(ctx context.Context, a types.Activity) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(a)
	_ = kvStore.Set(types.ActivityKey(a.ID), bz)
	// Index by address
	_ = kvStore.Set(types.ActivityByAddrKey(a.Address, a.ID), []byte{1})
}

func (k Keeper) GetGlobalActivityList(ctx context.Context) types.GlobalActivityList {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.GlobalActivityListKey))
	if err != nil || bz == nil {
		return types.GlobalActivityList{ActivityIDs: []uint64{}}
	}
	var gal types.GlobalActivityList
	if err := json.Unmarshal(bz, &gal); err != nil {
		return types.GlobalActivityList{ActivityIDs: []uint64{}}
	}
	return gal
}

func (k Keeper) SetGlobalActivityList(ctx context.Context, gal types.GlobalActivityList) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(gal)
	_ = kvStore.Set([]byte(types.GlobalActivityListKey), bz)
}

func (k Keeper) RecordActivity(ctx context.Context, activityType, address, description string, amount math.Int, denom string, blockHeight int64) uint64 {
	actID := k.GetNextActivityID(ctx)
	activity := types.Activity{
		ID:           actID,
		ActivityType: activityType,
		Address:      address,
		Description:  description,
		Amount:       amount,
		Denom:        denom,
		Block:        blockHeight,
	}
	k.SetActivity(ctx, activity)
	k.SetNextActivityID(ctx, actID+1)

	// Add to global list, trim to max
	gal := k.GetGlobalActivityList(ctx)
	gal.ActivityIDs = append(gal.ActivityIDs, actID)
	if len(gal.ActivityIDs) > types.MaxGlobalActivities {
		gal.ActivityIDs = gal.ActivityIDs[len(gal.ActivityIDs)-types.MaxGlobalActivities:]
	}
	k.SetGlobalActivityList(ctx, gal)

	return actID
}

// GetUserActivities returns the last N activities for a user (up to MaxPerUserActivities)
func (k Keeper) GetUserActivities(ctx context.Context, address string) []types.Activity {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := types.ActivityByAddrPrefixKey(address)
	iter, err := kvStore.Iterator(prefix, append(prefix, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var ids []uint64
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if len(key) >= 8 {
			id := binary.BigEndian.Uint64(key[len(key)-8:])
			ids = append(ids, id)
		}
	}

	// Take last MaxPerUserActivities
	if len(ids) > types.MaxPerUserActivities {
		ids = ids[len(ids)-types.MaxPerUserActivities:]
	}

	var activities []types.Activity
	for _, id := range ids {
		if a, found := k.GetActivity(ctx, id); found {
			activities = append(activities, a)
		}
	}
	return activities
}

// ============================================================
// Competition CRUD
// ============================================================

func (k Keeper) GetNextCompetitionID(ctx context.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.NextCompetitionIDKey))
	if err != nil || bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetNextCompetitionID(ctx context.Context, id uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	_ = kvStore.Set([]byte(types.NextCompetitionIDKey), bz)
}

func (k Keeper) GetCompetition(ctx context.Context, id uint64) (types.Competition, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.CompetitionKey(id))
	if err != nil || bz == nil {
		return types.Competition{}, false
	}
	var c types.Competition
	if err := json.Unmarshal(bz, &c); err != nil {
		return types.Competition{}, false
	}
	return c, true
}

func (k Keeper) SetCompetition(ctx context.Context, c types.Competition) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(c)
	_ = kvStore.Set(types.CompetitionKey(c.ID), bz)
}

func (k Keeper) GetCompetitionEntry(ctx context.Context, compID uint64, address string) (types.CompetitionEntry, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.CompEntryKey(compID, address))
	if err != nil || bz == nil {
		return types.CompetitionEntry{}, false
	}
	var e types.CompetitionEntry
	if err := json.Unmarshal(bz, &e); err != nil {
		return types.CompetitionEntry{}, false
	}
	return e, true
}

func (k Keeper) SetCompetitionEntry(ctx context.Context, e types.CompetitionEntry) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(e)
	_ = kvStore.Set(types.CompEntryKey(e.CompetitionID, e.Address), bz)
}

func (k Keeper) GetAllCompetitionEntries(ctx context.Context, compID uint64) []types.CompetitionEntry {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := types.CompEntryPrefixKey(compID)
	iter, err := kvStore.Iterator(prefix, append(prefix, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var entries []types.CompetitionEntry
	for ; iter.Valid(); iter.Next() {
		var e types.CompetitionEntry
		if err := json.Unmarshal(iter.Value(), &e); err == nil {
			entries = append(entries, e)
		}
	}
	return entries
}

// ExecuteCreateCompetition creates a new trading competition
func (k Keeper) ExecuteCreateCompetition(ctx context.Context, creator, name string, startBlock, endBlock int64, prizeDenom string, prizePool, entryFee math.Int, maxParticipants uint64) (uint64, error) {
	compID := k.GetNextCompetitionID(ctx)

	// If prize pool > 0, transfer from creator to module
	if prizePool.IsPositive() {
		creatorAddr, err := sdk.AccAddressFromBech32(creator)
		if err != nil {
			return 0, types.ErrInvalidAddress
		}
		coins := sdk.NewCoins(sdk.NewCoin(prizeDenom, prizePool))
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, creatorAddr, types.ModuleName, coins); err != nil {
			return 0, types.ErrInsufficientFunds
		}
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	status := types.CompStatusUpcoming
	if sdkCtx.BlockHeight() >= startBlock {
		status = types.CompStatusActive
	}

	comp := types.Competition{
		ID:              compID,
		Name:            name,
		Creator:         creator,
		StartBlock:      startBlock,
		EndBlock:        endBlock,
		PrizeDenom:      prizeDenom,
		PrizePool:       prizePool,
		EntryFee:        entryFee,
		MaxParticipants: maxParticipants,
		Participants:    0,
		Status:          status,
	}
	k.SetCompetition(ctx, comp)
	k.SetNextCompetitionID(ctx, compID+1)

	return compID, nil
}

// ExecuteJoinCompetition adds a user to a competition
func (k Keeper) ExecuteJoinCompetition(ctx context.Context, sender string, compID uint64) error {
	comp, found := k.GetCompetition(ctx, compID)
	if !found {
		return types.ErrCompetitionNotFound
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Update status if needed
	if comp.Status == types.CompStatusUpcoming && sdkCtx.BlockHeight() >= comp.StartBlock {
		comp.Status = types.CompStatusActive
	}

	if comp.Status != types.CompStatusActive && comp.Status != types.CompStatusUpcoming {
		return types.ErrCompetitionNotActive
	}

	if comp.Participants >= comp.MaxParticipants {
		return types.ErrCompetitionFull
	}

	// Check if already joined
	if _, joined := k.GetCompetitionEntry(ctx, compID, sender); joined {
		return types.ErrAlreadyJoined
	}

	// Collect entry fee
	if comp.EntryFee.IsPositive() {
		senderAddr, err := sdk.AccAddressFromBech32(sender)
		if err != nil {
			return types.ErrInvalidAddress
		}
		coins := sdk.NewCoins(sdk.NewCoin(comp.PrizeDenom, comp.EntryFee))
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, coins); err != nil {
			return types.ErrInsufficientFunds
		}
		comp.PrizePool = comp.PrizePool.Add(comp.EntryFee)
	}

	// Get current portfolio value as starting value
	portfolio := k.GetOrCreatePortfolio(ctx, sender)
	startingValue := portfolio.TotalValue
	if startingValue.IsZero() {
		// Use a baseline of 1 to avoid division by zero in PnL%
		startingValue = math.NewInt(1)
	}

	entry := types.CompetitionEntry{
		Address:       sender,
		CompetitionID: compID,
		StartingValue: startingValue,
		CurrentValue:  startingValue,
		PnLPercent:    math.LegacyZeroDec(),
		Rank:          0,
	}
	k.SetCompetitionEntry(ctx, entry)

	comp.Participants++
	k.SetCompetition(ctx, comp)

	return nil
}

// EndCompetitionAndDistribute ends a competition and distributes prizes to top 3
func (k Keeper) EndCompetitionAndDistribute(ctx context.Context, compID uint64) ([]string, error) {
	comp, found := k.GetCompetition(ctx, compID)
	if !found {
		return nil, types.ErrCompetitionNotFound
	}

	if comp.Status == types.CompStatusDistributed {
		return nil, types.ErrPrizesAlreadyDistributed
	}

	// Get all entries and rank them
	entries := k.GetAllCompetitionEntries(ctx, compID)
	if len(entries) == 0 {
		comp.Status = types.CompStatusDistributed
		k.SetCompetition(ctx, comp)
		return []string{}, nil
	}

	// Update current values from portfolios
	for i := range entries {
		portfolio := k.GetOrCreatePortfolio(ctx, entries[i].Address)
		entries[i].CurrentValue = portfolio.TotalValue
		if entries[i].StartingValue.IsPositive() {
			diff := entries[i].CurrentValue.Sub(entries[i].StartingValue)
			entries[i].PnLPercent = math.LegacyNewDecFromInt(diff).Quo(math.LegacyNewDecFromInt(entries[i].StartingValue)).MulInt64(100)
		}
	}

	// Sort by RealizedPnL descending (M-5 fix: harder to game than TotalValue which
	// uses 1:1 token valuation that can be exploited with low-value tokens).
	sort.Slice(entries, func(i, j int) bool {
		pi, _ := k.GetPortfolio(ctx, entries[i].Address)
		pj, _ := k.GetPortfolio(ctx, entries[j].Address)
		if pi.RealizedPnL.Equal(pj.RealizedPnL) {
			return entries[i].Address < entries[j].Address
		}
		return pi.RealizedPnL.GT(pj.RealizedPnL)
	})

	// Assign ranks and save
	for i := range entries {
		entries[i].Rank = uint64(i + 1)
		k.SetCompetitionEntry(ctx, entries[i])
	}

	// Distribute prizes: top 3 get 50%, 30%, 20%
	var winners []string
	if comp.PrizePool.IsPositive() {
		shares := []math.LegacyDec{
			math.LegacyNewDecWithPrec(50, 2), // 50%
			math.LegacyNewDecWithPrec(30, 2), // 30%
			math.LegacyNewDecWithPrec(20, 2), // 20%
		}

		maxWinners := len(entries)
		if maxWinners > 3 {
			maxWinners = 3
		}

		// Redistribute shares if fewer than 3 participants
		if maxWinners == 1 {
			shares = []math.LegacyDec{math.LegacyOneDec()}
		} else if maxWinners == 2 {
			shares = []math.LegacyDec{
				math.LegacyNewDecWithPrec(60, 2),
				math.LegacyNewDecWithPrec(40, 2),
			}
		}

		for i := 0; i < maxWinners; i++ {
			prizeAmount := shares[i].MulInt(comp.PrizePool).TruncateInt()
			if prizeAmount.IsPositive() {
				winnerAddr, err := sdk.AccAddressFromBech32(entries[i].Address)
				if err != nil {
					continue
				}
				coins := sdk.NewCoins(sdk.NewCoin(comp.PrizeDenom, prizeAmount))
				if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, winnerAddr, coins); err != nil {
					continue
				}
				winners = append(winners, entries[i].Address)
			}
		}
	}

	comp.Status = types.CompStatusDistributed
	k.SetCompetition(ctx, comp)

	return winners, nil
}

// RecalculatePortfolio recalculates a user's portfolio from their bank balances
func (k Keeper) RecalculatePortfolio(ctx context.Context, address string) math.Int {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return math.ZeroInt()
	}

	balances := k.bankKeeper.GetAllBalances(ctx, addr)
	totalValue := math.ZeroInt()

	for _, coin := range balances {
		asset := types.PortfolioAsset{
			Address:        address,
			Denom:          coin.Denom,
			Amount:         coin.Amount,
			Value:          coin.Amount, // 1:1 valuation (since everything is in usyreen base)
			PnL:            math.ZeroInt(),
			PercentOfTotal: math.LegacyZeroDec(),
		}

		// Preserve existing PnL
		if existing, found := k.GetPortfolioAsset(ctx, address, coin.Denom); found {
			asset.PnL = existing.PnL
		}

		totalValue = totalValue.Add(coin.Amount)
		k.SetPortfolioAsset(ctx, asset)
	}

	// Calculate percentages
	if totalValue.IsPositive() {
		for _, coin := range balances {
			if a, found := k.GetPortfolioAsset(ctx, address, coin.Denom); found {
				a.PercentOfTotal = math.LegacyNewDecFromInt(a.Value).Quo(math.LegacyNewDecFromInt(totalValue)).MulInt64(100)
				k.SetPortfolioAsset(ctx, a)
			}
		}
	}

	portfolio := k.GetOrCreatePortfolio(ctx, address)
	portfolio.TotalValue = totalValue
	// L-7 fix: UnrealizedPnL = currentHoldingValue - costBasis (not minus realizedPnL,
	// which incorrectly double-counts realized profits).
	if portfolio.CostBasis.IsNil() {
		portfolio.CostBasis = math.ZeroInt()
	}
	portfolio.UnrealizedPnL = totalValue.Sub(portfolio.CostBasis)
	portfolio.TotalPnL = portfolio.RealizedPnL.Add(portfolio.UnrealizedPnL)
	portfolio.LastUpdated = sdkCtx.BlockHeight()
	k.SetPortfolio(ctx, portfolio)

	return totalValue
}

// ============================================================
// BeginBlock: auto-end expired competitions, detect whale activities
// ============================================================

// maxCompetitionsPerBlock bounds how many competitions are settled
// (EndCompetitionAndDistribute) in a single block. Settlement scans every
// participant entry, so an attacker able to grow the competition keyspace could
// otherwise force unbounded per-block work and stall the chain. Competitions
// beyond this cap settle in subsequent blocks (settlement a few blocks late is
// acceptable). Deterministic: comps are iterated in sorted ID order below, so
// all validators settle the identical bounded subset.
const maxCompetitionsPerBlock = 50

func (k Keeper) ProcessBeginBlock(ctx context.Context) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	height := sdkCtx.BlockHeight()

	// Collect competitions first, then mutate (never mutate during iteration)
	kvStore := k.storeService.OpenKVStore(ctx)
	iter, err := kvStore.Iterator([]byte(types.CompetitionPrefix), append([]byte(types.CompetitionPrefix), 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
	if err != nil {
		return
	}

	var comps []types.Competition
	for ; iter.Valid(); iter.Next() {
		var comp types.Competition
		if err := json.Unmarshal(iter.Value(), &comp); err == nil {
			comps = append(comps, comp)
		}
	}
	iter.Close()

	sort.Slice(comps, func(i, j int) bool { return comps[i].ID < comps[j].ID })

	settled := 0
	for _, comp := range comps {
		// Activate upcoming competitions
		if comp.Status == types.CompStatusUpcoming && height >= comp.StartBlock {
			comp.Status = types.CompStatusActive
			k.SetCompetition(ctx, comp)
		}

		// Auto-end expired active competitions
		if comp.Status == types.CompStatusActive && height >= comp.EndBlock {
			// Cap settlement work per block; remaining competitions settle in
			// subsequent blocks. Iteration is in sorted ID order (deterministic).
			if settled >= maxCompetitionsPerBlock {
				continue
			}
			comp.Status = types.CompStatusEnded
			k.SetCompetition(ctx, comp)
			settled++

			// Per-competition panic isolation: a single malformed competition
			// must never panic BeginBlocker and halt the chain. Settlement state
			// writes already happened above; wrap only the distribution call.
			compID := comp.ID
			func() {
				defer func() {
					if r := recover(); r != nil {
						sdkCtx.Logger().Error("portfolio: recovered from panic settling competition",
							"competition_id", compID, "recover", r)
					}
				}()
				k.EndCompetitionAndDistribute(ctx, compID)
			}()
		}
	}
}

// GetAllCompetitions returns all competitions (for genesis export)
func (k Keeper) GetAllCompetitions(ctx context.Context) []types.Competition {
	kvStore := k.storeService.OpenKVStore(ctx)
	iter, err := kvStore.Iterator([]byte(types.CompetitionPrefix), append([]byte(types.CompetitionPrefix), 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var comps []types.Competition
	for ; iter.Valid(); iter.Next() {
		var c types.Competition
		if err := json.Unmarshal(iter.Value(), &c); err == nil {
			comps = append(comps, c)
		}
	}
	return comps
}

func (k Keeper) getAllByPrefix(ctx context.Context, prefix string, unmarshal func([]byte) bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	iter, err := kvStore.Iterator([]byte(prefix), append([]byte(prefix), 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
	if err != nil { return }
	defer iter.Close()
	for ; iter.Valid(); iter.Next() { unmarshal(iter.Value()) }
}

func (k Keeper) GetAllPortfolios(ctx context.Context) []types.Portfolio {
	var items []types.Portfolio
	k.getAllByPrefix(ctx, types.PortfolioPrefix, func(bz []byte) bool {
		var v types.Portfolio; if json.Unmarshal(bz, &v) == nil { items = append(items, v) }; return true
	})
	return items
}

func (k Keeper) GetAllAssets(ctx context.Context) []types.PortfolioAsset {
	var items []types.PortfolioAsset
	k.getAllByPrefix(ctx, types.PortfolioAssetPrefix, func(bz []byte) bool {
		var v types.PortfolioAsset; if json.Unmarshal(bz, &v) == nil { items = append(items, v) }; return true
	})
	return items
}

func (k Keeper) GetAllTrades(ctx context.Context) []types.TradeRecord {
	var items []types.TradeRecord
	k.getAllByPrefix(ctx, types.TradeRecordPrefix, func(bz []byte) bool {
		var v types.TradeRecord; if json.Unmarshal(bz, &v) == nil { items = append(items, v) }; return true
	})
	return items
}

func (k Keeper) GetAllActivities(ctx context.Context) []types.Activity {
	var items []types.Activity
	k.getAllByPrefix(ctx, types.ActivityPrefix, func(bz []byte) bool {
		var v types.Activity; if json.Unmarshal(bz, &v) == nil { items = append(items, v) }; return true
	})
	return items
}

func (k Keeper) GetAllEntries(ctx context.Context) []types.CompetitionEntry {
	var items []types.CompetitionEntry
	k.getAllByPrefix(ctx, types.CompEntryPrefix, func(bz []byte) bool {
		var v types.CompetitionEntry; if json.Unmarshal(bz, &v) == nil { items = append(items, v) }; return true
	})
	return items
}
