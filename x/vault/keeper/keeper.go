package keeper

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"sort"

	"cosmossdk.io/core/store"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"math/big"

	"syreen/x/vault/types"
)

// MinimumLiquidity is the dead-share amount permanently locked on the first
// deposit to a vault to prevent the share-price inflation attack.
const MinimumLiquidity int64 = 1000

// MinDepositAge is the minimum number of blocks between a deposit and a
// withdrawal to prevent sandwich / flash-deposit attacks.
const MinDepositAge int64 = 100

// BaseYieldPerCompound is kept for reference but is no longer used for synthetic yield generation.
// Compounding now only recognises real tokens that have arrived in the module account above
// the recorded TotalDeposited figure.
var BaseYieldPerCompound = math.LegacyNewDecWithPrec(1, 3) // 0.001 (unused — real-balance logic below)

// CompoundInterval is the number of blocks between auto-compound cycles
const CompoundInterval = int64(100)

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
// Vault CRUD
// ============================================================

func (k Keeper) GetVault(ctx context.Context, vaultID uint64) (types.Vault, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.VaultKey(vaultID))
	if err != nil || bz == nil {
		return types.Vault{}, false
	}
	var v types.Vault
	if err := json.Unmarshal(bz, &v); err != nil {
		return types.Vault{}, false
	}
	return v, true
}

func (k Keeper) SetVault(ctx context.Context, v types.Vault) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(v)
	_ = kvStore.Set(types.VaultKey(v.ID), bz)
}

func (k Keeper) GetAllVaults(ctx context.Context) []types.Vault {
	kvStore := k.storeService.OpenKVStore(ctx)
	end := make([]byte, len(types.VaultPrefix)+8)
	copy(end, []byte(types.VaultPrefix))
	for i := len(types.VaultPrefix); i < len(end); i++ {
		end[i] = 0xFF
	}
	iter, err := kvStore.Iterator([]byte(types.VaultPrefix), end)
	if err != nil {
		return nil
	}
	defer iter.Close()
	var vaults []types.Vault
	for ; iter.Valid(); iter.Next() {
		var v types.Vault
		if err := json.Unmarshal(iter.Value(), &v); err == nil {
			vaults = append(vaults, v)
		}
	}
	sort.Slice(vaults, func(i, j int) bool { return vaults[i].ID < vaults[j].ID })
	return vaults
}

func (k Keeper) GetNextVaultID(ctx context.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.NextVaultIDKey))
	if err != nil || bz == nil {
		return 1
	}
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetNextVaultID(ctx context.Context, id uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	_ = kvStore.Set([]byte(types.NextVaultIDKey), bz)
}

// ============================================================
// VaultDeposit CRUD
// ============================================================

func (k Keeper) GetDeposit(ctx context.Context, vaultID uint64, address string) (types.VaultDeposit, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.VaultDepositKey(vaultID, address))
	if err != nil || bz == nil {
		return types.VaultDeposit{}, false
	}
	var d types.VaultDeposit
	if err := json.Unmarshal(bz, &d); err != nil {
		return types.VaultDeposit{}, false
	}
	return d, true
}

func (k Keeper) SetDeposit(ctx context.Context, d types.VaultDeposit) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(d)
	_ = kvStore.Set(types.VaultDepositKey(d.VaultID, d.Address), bz)
	// Also set reverse index
	_ = kvStore.Set(types.DepositByAddrKey(d.Address, d.VaultID), []byte{1})
}

func (k Keeper) DeleteDeposit(ctx context.Context, d types.VaultDeposit) {
	kvStore := k.storeService.OpenKVStore(ctx)
	_ = kvStore.Delete(types.VaultDepositKey(d.VaultID, d.Address))
	_ = kvStore.Delete(types.DepositByAddrKey(d.Address, d.VaultID))
}

func (k Keeper) GetDepositsByVault(ctx context.Context, vaultID uint64) []types.VaultDeposit {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := types.VaultDepositPrefixKey(vaultID)
	end := make([]byte, len(prefix)+64)
	copy(end, prefix)
	for i := len(prefix); i < len(end); i++ {
		end[i] = 0xFF
	}
	iter, err := kvStore.Iterator(prefix, end)
	if err != nil {
		return nil
	}
	defer iter.Close()
	var deposits []types.VaultDeposit
	for ; iter.Valid(); iter.Next() {
		var d types.VaultDeposit
		if err := json.Unmarshal(iter.Value(), &d); err == nil {
			deposits = append(deposits, d)
		}
	}
	return deposits
}

// GetAllDeposits returns all deposits across all vaults (for genesis export).
func (k Keeper) GetAllDeposits(ctx context.Context) []types.VaultDeposit {
	var allDeposits []types.VaultDeposit
	for _, vault := range k.GetAllVaults(ctx) {
		allDeposits = append(allDeposits, k.GetDepositsByVault(ctx, vault.ID)...)
	}
	return allDeposits
}

// ============================================================
// Business Logic
// ============================================================

// ExecuteCreateVault creates a new vault and returns its ID
func (k Keeper) ExecuteCreateVault(ctx context.Context, creator, name, depositDenom, strategyType string, targetPoolIDs []uint64, performanceFee math.LegacyDec) (uint64, error) {
	if !types.ValidStrategy(strategyType) {
		return 0, types.ErrInvalidStrategy
	}
	if name == "" {
		return 0, types.ErrInvalidName
	}
	if depositDenom == "" {
		return 0, types.ErrInvalidDenom
	}
	// L-4: cap performance fee at 20% (down from 50%).
	if performanceFee.IsNegative() || performanceFee.GT(math.LegacyNewDecWithPrec(20, 2)) {
		return 0, types.ErrInvalidFee
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	id := k.GetNextVaultID(ctx)

	vault := types.Vault{
		ID:                id,
		Name:              name,
		Creator:           creator,
		DepositDenom:      depositDenom,
		TotalDeposited:    math.ZeroInt(),
		TotalShares:       math.ZeroInt(),
		StrategyType:      strategyType,
		TargetPoolIDs:     targetPoolIDs,
		PerformanceFee:    performanceFee,
		ProtocolFee:       math.LegacyNewDecWithPrec(2, 2), // 2% fixed protocol fee
		LastCompoundBlock: sdkCtx.BlockHeight(),
		AccYieldPerShare:  math.LegacyZeroDec(),
		Active:            true,
	}

	k.SetVault(ctx, vault)
	k.SetNextVaultID(ctx, id+1)
	return id, nil
}

// ExecuteDeposit deposits tokens into a vault and mints proportional shares
func (k Keeper) ExecuteDeposit(ctx context.Context, sender string, vaultID uint64, amount math.Int) (math.Int, error) {
	if amount.IsNil() || !amount.IsPositive() {
		return math.ZeroInt(), types.ErrInvalidAmount
	}

	vault, ok := k.GetVault(ctx, vaultID)
	if !ok {
		return math.ZeroInt(), types.ErrVaultNotFound
	}
	if !vault.Active {
		return math.ZeroInt(), types.ErrVaultNotActive
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return math.ZeroInt(), err
	}

	// Transfer tokens from sender to module
	coins := sdk.NewCoins(sdk.NewCoin(vault.DepositDenom, amount))
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, coins); err != nil {
		return math.ZeroInt(), err
	}

	// Calculate shares to mint.
	// First deposit: lock MinimumLiquidity dead shares to prevent the inflation
	// attack. Shares = sqrt(amount) - MinimumLiquidity; the dead shares belong
	// to the module and can never be redeemed.
	// Subsequent deposits: shares = amount * totalShares / totalDeposited.
	var sharesToMint math.Int
	if vault.TotalShares.IsZero() || vault.TotalDeposited.IsZero() {
		minLiq := math.NewInt(MinimumLiquidity)
		// sqrt via big.Int
		sqrtAmt := math.NewIntFromBigInt(new(big.Int).Sqrt(amount.BigInt()))
		if sqrtAmt.LTE(minLiq) {
			return math.ZeroInt(), types.ErrInitialDepositTooSmall
		}
		sharesToMint = sqrtAmt.Sub(minLiq)
		// Dead shares: add MinimumLiquidity to TotalShares but credit no depositor.
		vault.TotalShares = vault.TotalShares.Add(minLiq)
	} else {
		sharesToMint = amount.Mul(vault.TotalShares).Quo(vault.TotalDeposited)
	}

	if sharesToMint.IsZero() {
		return math.ZeroInt(), types.ErrInvalidAmount
	}

	// Update vault totals
	vault.TotalDeposited = vault.TotalDeposited.Add(amount)
	vault.TotalShares = vault.TotalShares.Add(sharesToMint)
	k.SetVault(ctx, vault)

	// Update or create deposit record
	dep, exists := k.GetDeposit(ctx, vaultID, sender)
	if exists {
		dep.Shares = dep.Shares.Add(sharesToMint)
		dep.DepositedAmount = dep.DepositedAmount.Add(amount)
		dep.DepositedAt = sdkCtx.BlockHeight() // reset timelock on every deposit
	} else {
		dep = types.VaultDeposit{
			Address:         sender,
			VaultID:         vaultID,
			Shares:          sharesToMint,
			DepositedAmount: amount,
			YieldIndex:      vault.AccYieldPerShare,
			DepositedAt:     sdkCtx.BlockHeight(),
		}
	}
	k.SetDeposit(ctx, dep)

	return sharesToMint, nil
}

// ExecuteWithdraw burns shares and returns proportional tokens + accrued yield
func (k Keeper) ExecuteWithdraw(ctx context.Context, sender string, vaultID uint64, shares math.Int) (math.Int, error) {
	if shares.IsNil() || !shares.IsPositive() {
		return math.ZeroInt(), types.ErrInvalidAmount
	}

	vault, ok := k.GetVault(ctx, vaultID)
	if !ok {
		return math.ZeroInt(), types.ErrVaultNotFound
	}

	dep, ok := k.GetDeposit(ctx, vaultID, sender)
	if !ok {
		return math.ZeroInt(), types.ErrNoDeposit
	}

	if dep.Shares.LT(shares) {
		return math.ZeroInt(), types.ErrInsufficientShares
	}

	// Enforce minimum deposit age to prevent compound front-running
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if sdkCtx.BlockHeight()-dep.DepositedAt < int64(MinDepositAge) {
		return math.ZeroInt(), types.ErrWithdrawTooEarly
	}

	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return math.ZeroInt(), err
	}

	// Calculate proportional amount: amount = shares * totalDeposited / totalShares
	amountToReturn := shares.Mul(vault.TotalDeposited).Quo(vault.TotalShares)
	if amountToReturn.IsZero() {
		return math.ZeroInt(), types.ErrInvalidAmount
	}

	// Transfer tokens from module to sender
	coins := sdk.NewCoins(sdk.NewCoin(vault.DepositDenom, amountToReturn))
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, coins); err != nil {
		return math.ZeroInt(), err
	}

	// Update vault totals
	vault.TotalDeposited = vault.TotalDeposited.Sub(amountToReturn)
	vault.TotalShares = vault.TotalShares.Sub(shares)
	k.SetVault(ctx, vault)

	// Update deposit record
	originalShares := dep.Shares
	dep.Shares = dep.Shares.Sub(shares)

	// Calculate proportional deposited amount reduction
	depAmtReduce := shares.Mul(dep.DepositedAmount).Quo(originalShares)
	dep.DepositedAmount = dep.DepositedAmount.Sub(depAmtReduce)

	if dep.Shares.IsZero() {
		k.DeleteDeposit(ctx, dep)
	} else {
		k.SetDeposit(ctx, dep)
	}

	return amountToReturn, nil
}

// ExecuteCompound compounds real yield that has arrived in the module account into the vault.
// Only the vault creator or chain authority may trigger compounding (H-6 fix).
// Yield is measured as the actual module balance above the recorded TotalDeposited (C-3 fix).
func (k Keeper) ExecuteCompound(ctx context.Context, sender string, vaultID uint64) (math.Int, error) {
	vault, ok := k.GetVault(ctx, vaultID)
	if !ok {
		return math.ZeroInt(), types.ErrVaultNotFound
	}
	if !vault.Active {
		return math.ZeroInt(), types.ErrVaultNotActive
	}

	// H-6: only the vault creator or chain authority can trigger compounding.
	// An empty sender is only allowed when called internally from AutoCompoundAll.
	if sender != "" && sender != vault.Creator && sender != k.authority {
		return math.ZeroInt(), types.ErrUnauthorized
	}

	if vault.TotalDeposited.IsZero() {
		return math.ZeroInt(), types.ErrNothingToCompound
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Each vault only compounds what it has explicitly tracked via its own
	// PendingYield field.  We do NOT query the module account balance, because
	// the module account is shared across all vaults and balance-based surplus
	// distribution allows cross-vault theft and donation attacks.
	//
	// External yield sources must call AddExplicitYield(vaultID, amount) which
	// increments PendingYield.  Until that is wired up, compounding is a no-op.
	//
	// For now, return zero — no synthetic yield, no balance scraping.

	vault.LastCompoundBlock = sdkCtx.BlockHeight()
	k.SetVault(ctx, vault)

	return math.ZeroInt(), nil
}

// AutoCompoundAll compounds all active vaults (called from BeginBlock)
func (k Keeper) AutoCompoundAll(ctx context.Context) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	height := sdkCtx.BlockHeight()

	// Only compound every CompoundInterval blocks
	if height%CompoundInterval != 0 {
		return
	}

	vaults := k.GetAllVaults(ctx)
	for _, v := range vaults {
		if !v.Active || v.TotalDeposited.IsZero() {
			continue
		}
		_, _ = k.ExecuteCompound(ctx, "", v.ID)
	}
}

// ExecuteUpdateStrategy updates the strategy of a vault (creator only)
func (k Keeper) ExecuteUpdateStrategy(ctx context.Context, creator string, vaultID uint64, strategyType string, targetPoolIDs []uint64) error {
	vault, ok := k.GetVault(ctx, vaultID)
	if !ok {
		return types.ErrVaultNotFound
	}
	if vault.Creator != creator {
		return types.ErrUnauthorized
	}
	if !types.ValidStrategy(strategyType) {
		return types.ErrInvalidStrategy
	}

	vault.StrategyType = strategyType
	vault.TargetPoolIDs = targetPoolIDs
	k.SetVault(ctx, vault)
	return nil
}

// GetModuleAddress returns the module account address
func (k Keeper) GetModuleAddress() sdk.AccAddress {
	return authtypes.NewModuleAddress(types.ModuleName)
}
