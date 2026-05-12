package keeper

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sync"

	"cosmossdk.io/core/store"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"syreen/x/flashloan/types"
)

// MinimumLiquidity is the dead-share amount permanently locked on the first
// deposit to a flash pool to prevent the share-price inflation attack (C4).
const MinimumLiquidity int64 = 1000

type Keeper struct {
	cdc           codec.Codec
	storeService  store.KVStoreService
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	authority     string

	// activeFlashLoans tracks in-flight flash loans by sender address. It is
	// shared across all copies of the Keeper struct and used to reject
	// nested/recursive flash loan calls (H1).
	activeFlashLoans *sync.Map
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
		authority:        authority,
		activeFlashLoans: &sync.Map{},
	}
}

// ============================================================
// Flash Pool CRUD
// ============================================================

func (k Keeper) GetPool(ctx context.Context, denom string) (types.FlashLoanPool, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.FlashPoolKey(denom))
	if err != nil || bz == nil {
		return types.FlashLoanPool{}, false
	}
	var pool types.FlashLoanPool
	if err := json.Unmarshal(bz, &pool); err != nil {
		return types.FlashLoanPool{}, false
	}
	return pool, true
}

func (k Keeper) SetPool(ctx context.Context, pool types.FlashLoanPool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(pool)
	_ = kvStore.Set(types.FlashPoolKey(pool.Denom), bz)
}

func (k Keeper) GetAllPools(ctx context.Context) []types.FlashLoanPool {
	kvStore := k.storeService.OpenKVStore(ctx)
	iter, err := kvStore.Iterator([]byte(types.FlashPoolPrefix), append([]byte(types.FlashPoolPrefix), 0xFF))
	if err != nil {
		return nil
	}
	defer iter.Close()
	var pools []types.FlashLoanPool
	for ; iter.Valid(); iter.Next() {
		var pool types.FlashLoanPool
		if err := json.Unmarshal(iter.Value(), &pool); err == nil {
			pools = append(pools, pool)
		}
	}
	return pools
}

// ============================================================
// Flash Record CRUD
// ============================================================

func (k Keeper) GetNextRecordID(ctx context.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.NextRecordIDKey))
	if err != nil || bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetNextRecordID(ctx context.Context, id uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	_ = kvStore.Set([]byte(types.NextRecordIDKey), bz)
}

func (k Keeper) SetRecord(ctx context.Context, record types.FlashLoanRecord) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(record)
	_ = kvStore.Set(types.FlashRecordKey(record.ID), bz)
}

func (k Keeper) GetRecord(ctx context.Context, id uint64) (types.FlashLoanRecord, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.FlashRecordKey(id))
	if err != nil || bz == nil {
		return types.FlashLoanRecord{}, false
	}
	var record types.FlashLoanRecord
	if err := json.Unmarshal(bz, &record); err != nil {
		return types.FlashLoanRecord{}, false
	}
	return record, true
}

// ============================================================
// Flash Stats CRUD
// ============================================================

func (k Keeper) GetStats(ctx context.Context) types.FlashLoanStats {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.FlashStatsKey))
	if err != nil || bz == nil {
		return types.DefaultFlashLoanStats()
	}
	var stats types.FlashLoanStats
	if err := json.Unmarshal(bz, &stats); err != nil {
		return types.DefaultFlashLoanStats()
	}
	return stats
}

func (k Keeper) SetStats(ctx context.Context, stats types.FlashLoanStats) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(stats)
	_ = kvStore.Set([]byte(types.FlashStatsKey), bz)
}

// ============================================================
// FlashPoolDeposit CRUD (H-1: LP share tracking)
// ============================================================

func (k Keeper) GetFlashPoolDeposit(ctx context.Context, denom, depositor string) (types.FlashPoolDeposit, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.FlashPoolDepositKey(denom, depositor))
	if err != nil || bz == nil {
		return types.FlashPoolDeposit{}, false
	}
	var dep types.FlashPoolDeposit
	if err := json.Unmarshal(bz, &dep); err != nil {
		return types.FlashPoolDeposit{}, false
	}
	return dep, true
}

func (k Keeper) SetFlashPoolDeposit(ctx context.Context, dep types.FlashPoolDeposit) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(dep)
	_ = kvStore.Set(types.FlashPoolDepositKey(dep.Denom, dep.Depositor), bz)
}

func (k Keeper) DeleteFlashPoolDeposit(ctx context.Context, dep types.FlashPoolDeposit) {
	kvStore := k.storeService.OpenKVStore(ctx)
	_ = kvStore.Delete(types.FlashPoolDepositKey(dep.Denom, dep.Depositor))
}

// GetAllFlashPoolDeposits returns all flash pool deposit records (for genesis export)
func (k Keeper) GetAllFlashPoolDeposits(ctx context.Context) []types.FlashPoolDeposit {
	kvStore := k.storeService.OpenKVStore(ctx)
	iter, err := kvStore.Iterator([]byte(types.FlashPoolDepositPrefix), append([]byte(types.FlashPoolDepositPrefix), 0xFF))
	if err != nil {
		return nil
	}
	defer iter.Close()
	var deposits []types.FlashPoolDeposit
	for ; iter.Valid(); iter.Next() {
		var dep types.FlashPoolDeposit
		if err := json.Unmarshal(iter.Value(), &dep); err == nil {
			deposits = append(deposits, dep)
		}
	}
	return deposits
}

// GetAllRecords returns all flash loan records (for genesis export)
func (k Keeper) GetAllRecords(ctx context.Context) []types.FlashLoanRecord {
	kvStore := k.storeService.OpenKVStore(ctx)
	iter, err := kvStore.Iterator([]byte(types.FlashRecordPrefix), append([]byte(types.FlashRecordPrefix), 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
	if err != nil {
		return nil
	}
	defer iter.Close()
	var records []types.FlashLoanRecord
	for ; iter.Valid(); iter.Next() {
		var r types.FlashLoanRecord
		if err := json.Unmarshal(iter.Value(), &r); err == nil {
			records = append(records, r)
		}
	}
	return records
}

// ============================================================
// Core Logic
// ============================================================

// CreatePool creates a new flash loan pool for a given denom.
func (k Keeper) CreatePool(ctx context.Context, authority, denom string, feeRate math.LegacyDec) error {
	if authority != k.authority {
		return types.ErrUnauthorized
	}
	if _, found := k.GetPool(ctx, denom); found {
		return types.ErrPoolAlreadyExists
	}

	pool := types.FlashLoanPool{
		Denom:       denom,
		Available:   math.ZeroInt(),
		TotalLoaned: math.ZeroInt(),
		TotalFees:   math.ZeroInt(),
		FeeRate:     feeRate,
		TotalShares: math.ZeroInt(),
		Active:      true,
	}
	k.SetPool(ctx, pool)
	return nil
}

// FundPool adds liquidity to a flash loan pool by transferring tokens from sender to the module
// account and minting proportional LP shares for the funder (H-1 fix).
func (k Keeper) FundPool(ctx context.Context, sender, denom string, amount math.Int) error {
	if amount.IsNil() || !amount.IsPositive() {
		return types.ErrInvalidAmount
	}

	pool, found := k.GetPool(ctx, denom)
	if !found {
		return types.ErrPoolNotFound
	}
	if !pool.Active {
		return types.ErrPoolNotActive
	}

	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return err
	}

	coins := sdk.NewCoins(sdk.NewCoin(denom, amount))
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, coins); err != nil {
		return err
	}

	// H-1 / C-4: calculate LP shares for this deposit.
	// First funder: lock MinimumLiquidity dead shares in the pool to prevent the
	// share-price inflation attack — the funder receives (amount - MinimumLiquidity)
	// shares while pool.TotalShares is set to amount, leaving MinimumLiquidity
	// shares unowned and unredeemable forever.
	// Subsequent funders: shares = amount * totalShares / available.
	var sharesToMint math.Int
	firstDeposit := pool.TotalShares.IsZero() || pool.Available.IsZero()
	if firstDeposit {
		minLiq := math.NewInt(MinimumLiquidity)
		if amount.LTE(minLiq) {
			return types.ErrInitialDepositTooSmall
		}
		sharesToMint = amount.Sub(minLiq)
	} else {
		sharesToMint = amount.Mul(pool.TotalShares).Quo(pool.Available)
	}
	if !sharesToMint.IsPositive() {
		return types.ErrInvalidAmount
	}

	// Update or create the depositor's share record. Only sharesToMint is
	// credited to the funder — the MinimumLiquidity dead shares (first deposit
	// only) belong to no one and can never be redeemed.
	dep, exists := k.GetFlashPoolDeposit(ctx, denom, sender)
	if exists {
		dep.Shares = dep.Shares.Add(sharesToMint)
		dep.Amount = dep.Amount.Add(amount)
	} else {
		dep = types.FlashPoolDeposit{
			Depositor: sender,
			Denom:     denom,
			Amount:    amount,
			Shares:    sharesToMint,
		}
	}
	k.SetFlashPoolDeposit(ctx, dep)

	pool.Available = pool.Available.Add(amount)
	if firstDeposit {
		// Lock the dead shares: pool.TotalShares = amount, of which only
		// sharesToMint (= amount - MinimumLiquidity) is credited to a funder.
		pool.TotalShares = amount
	} else {
		pool.TotalShares = pool.TotalShares.Add(sharesToMint)
	}
	k.SetPool(ctx, pool)
	return nil
}

// ExecuteWithdrawFlashPool redeems LP shares and returns proportional tokens to the depositor (H-1 fix).
func (k Keeper) ExecuteWithdrawFlashPool(ctx context.Context, sender, denom string, shares math.Int) (math.Int, error) {
	if shares.IsNil() || !shares.IsPositive() {
		return math.ZeroInt(), types.ErrInvalidAmount
	}

	pool, found := k.GetPool(ctx, denom)
	if !found {
		return math.ZeroInt(), types.ErrPoolNotFound
	}
	if !pool.Active {
		return math.ZeroInt(), types.ErrPoolNotActive
	}

	dep, found := k.GetFlashPoolDeposit(ctx, denom, sender)
	if !found {
		return math.ZeroInt(), types.ErrNoDeposit
	}
	if dep.Shares.LT(shares) {
		return math.ZeroInt(), types.ErrInsufficientShares
	}

	if pool.TotalShares.IsZero() {
		return math.ZeroInt(), types.ErrInsufficientLiquidity
	}

	// Calculate the token amount proportional to the shares being redeemed.
	amountToReturn := shares.Mul(pool.Available).Quo(pool.TotalShares)
	if amountToReturn.IsZero() {
		return math.ZeroInt(), types.ErrInvalidAmount
	}
	if pool.Available.LT(amountToReturn) {
		return math.ZeroInt(), types.ErrInsufficientLiquidity
	}

	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return math.ZeroInt(), err
	}

	returnCoins := sdk.NewCoins(sdk.NewCoin(denom, amountToReturn))
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, returnCoins); err != nil {
		return math.ZeroInt(), err
	}

	// Burn the redeemed shares and update pool state.
	pool.Available = pool.Available.Sub(amountToReturn)
	pool.TotalShares = pool.TotalShares.Sub(shares)
	k.SetPool(ctx, pool)

	originalShares := dep.Shares
	dep.Shares = dep.Shares.Sub(shares)
	if dep.Shares.IsZero() {
		k.DeleteFlashPoolDeposit(ctx, dep)
	} else {
		// Reduce the recorded cost-basis amount proportionally to the shares redeemed.
		depAmtReduce := shares.Mul(dep.Amount).Quo(originalShares)
		dep.Amount = dep.Amount.Sub(depAmtReduce)
		k.SetFlashPoolDeposit(ctx, dep)
	}

	return amountToReturn, nil
}

// ExecuteFlashLoan executes a flash loan:
// 1. Sends borrowed amount from module to borrower
// 2. Checks that borrower has sent back principal + fee to the module account
// 3. If not repaid, returns error (tx reverts)
func (k Keeper) ExecuteFlashLoan(ctx context.Context, sender, denom string, amount math.Int) (math.Int, error) {
	if amount.IsNil() || !amount.IsPositive() {
		return math.Int{}, types.ErrInvalidAmount
	}

	// H1: prevent nested/recursive flash loans for the same sender. The guard is
	// keyed by sender and is held only for the duration of this single keeper
	// call (disburse + repay are sequential within this function, so the loan is
	// fully atomic). Any attempt to invoke ExecuteFlashLoan again for the same
	// sender — directly or via a nested keeper hook — will be rejected.
	guardKey := "fl/" + sender + "/" + denom
	if _, busy := k.activeFlashLoans.LoadOrStore(guardKey, struct{}{}); busy {
		return math.Int{}, types.ErrFlashLoanInProgress
	}
	defer k.activeFlashLoans.Delete(guardKey)

	pool, found := k.GetPool(ctx, denom)
	if !found {
		return math.Int{}, types.ErrPoolNotFound
	}
	if !pool.Active {
		return math.Int{}, types.ErrPoolNotActive
	}
	if pool.Available.LT(amount) {
		return math.Int{}, types.ErrInsufficientLiquidity
	}

	// Calculate fee: 0.09% of borrowed amount
	fee := math.LegacyNewDecFromInt(amount).Mul(pool.FeeRate).TruncateInt()
	if fee.IsZero() {
		fee = math.OneInt() // minimum 1 unit fee
	}

	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return math.Int{}, err
	}

	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)

	// Record the module balance BEFORE the loan (used for final safety check)
	preFlashBalance := k.bankKeeper.GetBalance(ctx, moduleAddr, denom).Amount

	// Send borrowed funds from module to borrower
	coins := sdk.NewCoins(sdk.NewCoin(denom, amount))
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, coins); err != nil {
		return math.Int{}, fmt.Errorf("failed to send borrowed funds: %w", err)
	}

	// The borrower must have already pre-funded the repayment (principal + fee)
	// back to the module account. In Cosmos SDK, since there's no internal callback,
	// the user must include the repayment in the same transaction (via multi-msg tx)
	// or pre-fund the module account before this message executes.
	//
	// Check: module balance after lending should be >= balance_before - amount + amount + fee
	// i.e., module balance should be >= balance_before + fee
	// This means borrower must send (amount + fee) back to the module.

	// The borrower needs to send repayment back to module.
	// We attempt to pull principal + fee from the borrower.
	repayCoins := sdk.NewCoins(sdk.NewCoin(denom, amount.Add(fee)))
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, repayCoins); err != nil {
		return math.Int{}, types.ErrRepaymentFailed.Wrapf("borrower must have principal + fee (%s + %s = %s %s): %s",
			amount.String(), fee.String(), amount.Add(fee).String(), denom, err.Error())
	}

	// Verify module balance is at least what it was before + fee
	balanceAfter := k.bankKeeper.GetBalance(ctx, moduleAddr, denom)
	expectedMin := preFlashBalance.Add(fee)
	if balanceAfter.Amount.LT(expectedMin) {
		return math.Int{}, types.ErrRepaymentFailed.Wrapf(
			"module balance after repayment %s < expected %s",
			balanceAfter.Amount.String(), expectedMin.String(),
		)
	}

	// Final safety check: verify module received full repayment
	finalBalance := k.bankKeeper.GetBalance(ctx, moduleAddr, denom)
	if finalBalance.Amount.LT(preFlashBalance.Add(fee)) {
		return math.Int{}, types.ErrRepaymentInsufficient
	}

	// Success — update pool state
	pool.Available = pool.Available.Add(fee) // net gain is the fee
	pool.TotalLoaned = pool.TotalLoaned.Add(amount)
	pool.TotalFees = pool.TotalFees.Add(fee)
	k.SetPool(ctx, pool)

	// Record the flash loan
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	recordID := k.GetNextRecordID(ctx)
	record := types.FlashLoanRecord{
		ID:       recordID,
		Borrower: sender,
		Denom:    denom,
		Amount:   amount,
		Fee:      fee,
		Block:    sdkCtx.BlockHeight(),
		Success:  true,
	}
	k.SetRecord(ctx, record)
	k.SetNextRecordID(ctx, recordID+1)

	// Update global stats
	stats := k.GetStats(ctx)
	stats.TotalExecuted++
	stats.TotalVolume = stats.TotalVolume.Add(amount)
	stats.TotalFeesEarned = stats.TotalFeesEarned.Add(fee)
	k.SetStats(ctx, stats)

	return fee, nil
}
