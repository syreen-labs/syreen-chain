package keeper

import (
	"context"
	"encoding/json"
	"sort"

	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/farming/types"
)

type Keeper struct {
	cdc          codec.Codec
	storeService store.KVStoreService
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
		cdc:          cdc,
		storeService: storeService,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		authority:     authority,
	}
}

func (k Keeper) Logger(ctx context.Context) interface{} {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", types.ModuleName)
}

// ============================================================
// Farm CRUD
// ============================================================

func (k Keeper) GetFarm(ctx context.Context, poolID uint64) (types.FarmPool, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.FarmKey(poolID))
	if err != nil || bz == nil {
		return types.FarmPool{}, false
	}
	var farm types.FarmPool
	if err := json.Unmarshal(bz, &farm); err != nil {
		return types.FarmPool{}, false
	}
	return farm, true
}

func (k Keeper) SetFarm(ctx context.Context, farm types.FarmPool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(farm)
	_ = kvStore.Set(types.FarmKey(farm.PoolID), bz)
}

func (k Keeper) GetAllFarms(ctx context.Context) []types.FarmPool {
	kvStore := k.storeService.OpenKVStore(ctx)
	iter, err := kvStore.Iterator([]byte(types.FarmPrefix), append([]byte(types.FarmPrefix), 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var farms []types.FarmPool
	for ; iter.Valid(); iter.Next() {
		var farm types.FarmPool
		if err := json.Unmarshal(iter.Value(), &farm); err == nil {
			farms = append(farms, farm)
		}
	}
	sort.Slice(farms, func(i, j int) bool {
		return farms[i].PoolID < farms[j].PoolID
	})
	return farms
}

// ============================================================
// Position CRUD
// ============================================================

func (k Keeper) GetPosition(ctx context.Context, poolID uint64, address string) (types.Position, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.PositionKey(poolID, address))
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
	_ = kvStore.Set(types.PositionKey(pos.PoolID, pos.Address), bz)
	// Also set address index
	_ = kvStore.Set(types.PositionByAddrKey(pos.Address, pos.PoolID), []byte{1})
}

func (k Keeper) DeletePosition(ctx context.Context, poolID uint64, address string) {
	kvStore := k.storeService.OpenKVStore(ctx)
	_ = kvStore.Delete(types.PositionKey(poolID, address))
	_ = kvStore.Delete(types.PositionByAddrKey(address, poolID))
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
		// Extract poolID from last 8 bytes of key
		if len(key) >= 8 {
			poolIDBytes := key[len(key)-8:]
			poolID := uint64(poolIDBytes[0])<<56 | uint64(poolIDBytes[1])<<48 | uint64(poolIDBytes[2])<<40 | uint64(poolIDBytes[3])<<32 | uint64(poolIDBytes[4])<<24 | uint64(poolIDBytes[5])<<16 | uint64(poolIDBytes[6])<<8 | uint64(poolIDBytes[7])
			if pos, ok := k.GetPosition(ctx, poolID, address); ok {
				positions = append(positions, pos)
			}
		}
	}
	return positions
}

func (k Keeper) GetAllPositionsForPool(ctx context.Context, poolID uint64) []types.Position {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := types.PositionPoolPrefix(poolID)
	iter, err := kvStore.Iterator(prefix, append(prefix, 0xFF))
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
	sort.Slice(positions, func(i, j int) bool {
		return positions[i].Address < positions[j].Address
	})
	return positions
}

// GetAllPositions returns all farming positions across all pools (for genesis export).
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
	sort.Slice(positions, func(i, j int) bool {
		if positions[i].PoolID != positions[j].PoolID {
			return positions[i].PoolID < positions[j].PoolID
		}
		return positions[i].Address < positions[j].Address
	})
	return positions
}

// ============================================================
// Core Farming Logic
// ============================================================

// UpdateFarmRewards updates the accumulated reward per share for a farm
// Called before any stake/unstake/claim to ensure accurate calculations
func (k Keeper) UpdateFarmRewards(ctx context.Context, farm *types.FarmPool) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentBlock := sdkCtx.BlockHeight()

	if currentBlock <= farm.LastRewardBlock {
		return
	}

	if farm.TotalStaked.IsZero() || !farm.Active {
		farm.LastRewardBlock = currentBlock
		return
	}

	// Don't distribute rewards before start or after end
	fromBlock := farm.LastRewardBlock
	if fromBlock < farm.StartBlock {
		fromBlock = farm.StartBlock
	}
	toBlock := currentBlock
	if farm.EndBlock > 0 && toBlock > farm.EndBlock {
		toBlock = farm.EndBlock
	}

	if toBlock <= fromBlock {
		farm.LastRewardBlock = currentBlock
		return
	}

	// Calculate rewards for the period
	blocks := toBlock - fromBlock
	totalReward := farm.RewardPerBlock.Mul(math.NewInt(blocks))

	// Check if farming module has enough funds to distribute
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	balance := k.bankKeeper.GetBalance(ctx, moduleAddr, farm.RewardDenom)
	if balance.Amount.LT(totalReward) {
		// M-9 fix: shortfall — emit an event and a log warning so operators
		// can detect under-funded farms instead of silently capping.
		requested := totalReward
		actualPaid := balance.Amount
		sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
			"reward_shortfall",
			sdk.NewAttribute("pool_id", fmt.Sprintf("%d", farm.PoolID)),
			sdk.NewAttribute("requested", requested.String()),
			sdk.NewAttribute("actual_paid", actualPaid.String()),
			sdk.NewAttribute("reward_denom", farm.RewardDenom),
		))
		sdkCtx.Logger().Error(
			"farming reward shortfall",
			"module", types.ModuleName,
			"pool_id", farm.PoolID,
			"requested", requested.String(),
			"actual_paid", actualPaid.String(),
			"reward_denom", farm.RewardDenom,
		)

		// If the module account has zero of this denom, do NOT update the
		// per-share accumulator at all — corrupting AccRewardPerShare with
		// a zero would permanently lock those reward blocks. Instead skip
		// this update so the rewards can be paid once funds arrive.
		if balance.Amount.IsZero() {
			// Intentionally do not advance LastRewardBlock either, so the
			// missed blocks will be retried next time funds are present.
			return
		}

		// Distribute whatever is available
		totalReward = balance.Amount
	}

	if totalReward.IsPositive() {
		// accRewardPerShare += totalReward * 1e12 / totalStaked
		rewardScaled := math.LegacyNewDecFromInt(totalReward).Mul(types.RewardPrecision)
		perShare := rewardScaled.Quo(math.LegacyNewDecFromInt(farm.TotalStaked))
		farm.AccRewardPerShare = farm.AccRewardPerShare.Add(perShare)
	}

	farm.LastRewardBlock = currentBlock
}

// CalculatePendingReward calculates the pending reward for a position
func (k Keeper) CalculatePendingReward(farm types.FarmPool, pos types.Position) math.Int {
	if pos.Amount.IsZero() {
		return pos.PendingReward
	}

	// pending = (amount * accRewardPerShare / 1e12) - rewardDebt + pendingReward
	accReward := math.LegacyNewDecFromInt(pos.Amount).Mul(farm.AccRewardPerShare).Quo(types.RewardPrecision)
	pending := accReward.Sub(pos.RewardDebt).TruncateInt()

	if pending.IsNegative() {
		pending = math.ZeroInt()
	}

	return pending.Add(pos.PendingReward)
}

// StakeLP adds LP tokens to a farming pool
func (k Keeper) StakeLP(ctx context.Context, sender string, poolID uint64, amount math.Int) (math.Int, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	farm, ok := k.GetFarm(ctx, poolID)
	if !ok {
		return math.Int{}, types.ErrFarmNotFound
	}
	if !farm.Active {
		return math.Int{}, types.ErrFarmNotActive
	}
	if farm.StartBlock > 0 && sdkCtx.BlockHeight() < farm.StartBlock {
		return math.Int{}, types.ErrFarmNotStarted
	}
	if farm.EndBlock > 0 && sdkCtx.BlockHeight() > farm.EndBlock {
		return math.Int{}, types.ErrFarmEnded
	}

	// Validate LP denom is configured
	if farm.LPDenom == "" {
		return math.ZeroInt(), fmt.Errorf("farm LP denom not configured")
	}

	// Update farm rewards before any changes
	k.UpdateFarmRewards(ctx, &farm)

	// Get or create position
	pos, exists := k.GetPosition(ctx, poolID, sender)
	pendingReward := math.ZeroInt()

	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return math.Int{}, err
	}

	if exists && pos.Amount.IsPositive() {
		// Calculate pending rewards and send them immediately (M-6 fix)
		pendingReward = k.CalculatePendingReward(farm, pos)
		if pendingReward.IsPositive() {
			rewardCoins := sdk.NewCoins(sdk.NewCoin(farm.RewardDenom, pendingReward))
			moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
			balance := k.bankKeeper.GetBalance(ctx, moduleAddr, farm.RewardDenom)
			if balance.Amount.GTE(pendingReward) {
				if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, rewardCoins); err == nil {
					pos.PendingReward = math.ZeroInt() // cleared since sent
				} else {
					pos.PendingReward = pendingReward // keep for later claim
				}
			} else {
				pos.PendingReward = pendingReward // keep for later claim
			}
		}
	} else {
		pos = types.Position{
			Address:       sender,
			PoolID:        poolID,
			Amount:        math.ZeroInt(),
			RewardDebt:    math.LegacyZeroDec(),
			PendingReward: math.ZeroInt(),
			StakedAt:      sdkCtx.BlockHeight(),
		}
	}

	// Transfer LP tokens from user to farming module
	lpCoins := sdk.NewCoins(sdk.NewCoin(farm.LPDenom, amount))
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, lpCoins); err != nil {
		return math.Int{}, err
	}

	// Update position
	pos.Amount = pos.Amount.Add(amount)
	// rewardDebt = amount * accRewardPerShare / 1e12
	pos.RewardDebt = math.LegacyNewDecFromInt(pos.Amount).Mul(farm.AccRewardPerShare).Quo(types.RewardPrecision)
	k.SetPosition(ctx, pos)

	// Update farm total
	farm.TotalStaked = farm.TotalStaked.Add(amount)
	k.SetFarm(ctx, farm)

	// Emit event
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"farming_stake",
		sdk.NewAttribute("sender", sender),
		sdk.NewAttribute("pool_id", fmt.Sprintf("%d", poolID)),
		sdk.NewAttribute("amount", amount.String()),
	))

	return pendingReward, nil
}

// UnstakeLP removes LP tokens from a farming pool and claims rewards
func (k Keeper) UnstakeLP(ctx context.Context, sender string, poolID uint64, amount math.Int) (math.Int, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	farm, ok := k.GetFarm(ctx, poolID)
	if !ok {
		return math.Int{}, types.ErrFarmNotFound
	}

	pos, exists := k.GetPosition(ctx, poolID, sender)
	if !exists {
		return math.Int{}, types.ErrNoPosition
	}
	if pos.Amount.LT(amount) {
		return math.Int{}, types.ErrInsufficientStake
	}

	// Update farm rewards
	k.UpdateFarmRewards(ctx, &farm)

	// Calculate pending rewards
	claimedReward := k.CalculatePendingReward(farm, pos)

	// Send rewards to user
	if claimedReward.IsPositive() {
		senderAddr, _ := sdk.AccAddressFromBech32(sender)
		rewardCoins := sdk.NewCoins(sdk.NewCoin(farm.RewardDenom, claimedReward))
		moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
		balance := k.bankKeeper.GetBalance(ctx, moduleAddr, farm.RewardDenom)
		if balance.Amount.GTE(claimedReward) {
			if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, rewardCoins); err != nil {
				claimedReward = math.ZeroInt()
			}
		} else {
			claimedReward = math.ZeroInt()
		}
	}

	// Return LP tokens to user
	senderAddr, _ := sdk.AccAddressFromBech32(sender)
	lpCoins := sdk.NewCoins(sdk.NewCoin(farm.LPDenom, amount))
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, lpCoins); err != nil {
		return math.Int{}, err
	}

	// Update position — preserve unclaimed reward if send failed
	pos.Amount = pos.Amount.Sub(amount)
	if claimedReward.IsZero() {
		// Reward send failed or insufficient balance — keep the pending reward
		// so the user can try claiming later.
		pos.PendingReward = k.CalculatePendingReward(farm, pos)
	} else {
		pos.PendingReward = math.ZeroInt()
	}

	if pos.Amount.IsZero() && pos.PendingReward.IsZero() {
		k.DeletePosition(ctx, poolID, sender)
	} else {
		pos.RewardDebt = math.LegacyNewDecFromInt(pos.Amount).Mul(farm.AccRewardPerShare).Quo(types.RewardPrecision)
		k.SetPosition(ctx, pos)
	}

	// Update farm total
	farm.TotalStaked = farm.TotalStaked.Sub(amount)
	if farm.TotalStaked.IsNegative() {
		farm.TotalStaked = math.ZeroInt()
	}
	k.SetFarm(ctx, farm)

	// Emit event
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"farming_unstake",
		sdk.NewAttribute("sender", sender),
		sdk.NewAttribute("pool_id", fmt.Sprintf("%d", poolID)),
		sdk.NewAttribute("amount", amount.String()),
		sdk.NewAttribute("reward", claimedReward.String()),
	))

	return claimedReward, nil
}

// ClaimFarmReward claims pending rewards without unstaking
func (k Keeper) ClaimFarmReward(ctx context.Context, sender string, poolID uint64) (math.Int, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	farm, ok := k.GetFarm(ctx, poolID)
	if !ok {
		return math.Int{}, types.ErrFarmNotFound
	}

	pos, exists := k.GetPosition(ctx, poolID, sender)
	if !exists {
		return math.Int{}, types.ErrNoPosition
	}

	// Update farm rewards
	k.UpdateFarmRewards(ctx, &farm)

	// Calculate pending
	pending := k.CalculatePendingReward(farm, pos)
	if pending.IsZero() {
		return math.Int{}, types.ErrNoPendingReward
	}

	// Send rewards
	senderAddr, _ := sdk.AccAddressFromBech32(sender)
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	balance := k.bankKeeper.GetBalance(ctx, moduleAddr, farm.RewardDenom)
	if balance.Amount.LT(pending) {
		pending = balance.Amount
	}
	if pending.IsPositive() {
		rewardCoins := sdk.NewCoins(sdk.NewCoin(farm.RewardDenom, pending))
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, rewardCoins); err != nil {
			return math.Int{}, err
		}
	}

	// Update position debt
	pos.PendingReward = math.ZeroInt()
	pos.RewardDebt = math.LegacyNewDecFromInt(pos.Amount).Mul(farm.AccRewardPerShare).Quo(types.RewardPrecision)
	k.SetPosition(ctx, pos)
	k.SetFarm(ctx, farm)

	// Emit event
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"farming_claim",
		sdk.NewAttribute("sender", sender),
		sdk.NewAttribute("pool_id", fmt.Sprintf("%d", poolID)),
		sdk.NewAttribute("amount", pending.String()),
	))

	return pending, nil
}

// CreateFarmPool creates a new farming pool (authority only).
// lpDenom must be the exact LP token denom produced by the DEX module for poolID.
func (k Keeper) CreateFarmPool(ctx context.Context, authority string, poolID uint64, lpDenom string, rewardPerBlock math.Int, startBlock, endBlock int64) error {
	if authority != k.authority {
		return types.ErrUnauthorized
	}

	if lpDenom == "" {
		return fmt.Errorf("LP denom must not be empty")
	}

	if _, exists := k.GetFarm(ctx, poolID); exists {
		return types.ErrFarmAlreadyExists
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if startBlock == 0 {
		startBlock = sdkCtx.BlockHeight()
	}

	farm := types.FarmPool{
		PoolID:            poolID,
		LPDenom:           lpDenom,
		RewardDenom:       "usyreen",
		RewardPerBlock:    rewardPerBlock,
		TotalStaked:       math.ZeroInt(),
		AccRewardPerShare: math.LegacyZeroDec(),
		LastRewardBlock:   startBlock,
		StartBlock:        startBlock,
		EndBlock:          endBlock,
		Active:            true,
		Creator:           authority,
	}

	k.SetFarm(ctx, farm)
	return nil
}

// UpdateFarmConfig updates an existing farm's settings
func (k Keeper) UpdateFarmConfig(ctx context.Context, authority string, poolID uint64, rewardPerBlock math.Int, active bool) error {
	if authority != k.authority {
		return types.ErrUnauthorized
	}

	farm, ok := k.GetFarm(ctx, poolID)
	if !ok {
		return types.ErrFarmNotFound
	}

	// Update rewards before changing rate
	k.UpdateFarmRewards(ctx, &farm)

	if !rewardPerBlock.IsNil() && rewardPerBlock.IsPositive() {
		farm.RewardPerBlock = rewardPerBlock
	}
	farm.Active = active
	k.SetFarm(ctx, farm)
	return nil
}

// DistributeRewards is called in BeginBlock to update all active farms
func (k Keeper) DistributeRewards(ctx context.Context) {
	farms := k.GetAllFarms(ctx)
	for _, farm := range farms {
		if !farm.Active || farm.TotalStaked.IsZero() {
			continue
		}
		k.UpdateFarmRewards(ctx, &farm)
		k.SetFarm(ctx, farm)
	}
}

// GetFarmAPR calculates the approximate APR for a farm
func (k Keeper) GetFarmAPR(ctx context.Context, farm types.FarmPool) string {
	if farm.TotalStaked.IsZero() {
		return "∞"
	}
	// APR = (rewardPerBlock * blocks_per_year) / totalStaked * 100
	// ~6,307,200 blocks/year at 5s blocks, but we use 500ms so ~63,072,000
	blocksPerYear := math.NewInt(63_072_000)
	yearlyReward := farm.RewardPerBlock.Mul(blocksPerYear)
	apr := math.LegacyNewDecFromInt(yearlyReward).Quo(math.LegacyNewDecFromInt(farm.TotalStaked)).Mul(math.LegacyNewDec(100))
	return apr.TruncateInt().String() + "%"
}

// getLPDenom is intentionally removed — LP denom must be explicitly provided
// when creating a farm via CreateFarmPool's lpDenom parameter.

// SetLPDenom allows setting the LP denom for a farm after creation
func (k Keeper) setLPDenom(ctx context.Context, poolID uint64, lpDenom string) error {
	farm, ok := k.GetFarm(ctx, poolID)
	if !ok {
		return types.ErrFarmNotFound
	}
	farm.LPDenom = lpDenom
	k.SetFarm(ctx, farm)
	return nil
}
