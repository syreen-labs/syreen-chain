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

	"syreen/x/launchpad/types"
)

type Keeper struct {
	cdc                 codec.Codec
	storeService        store.KVStoreService
	accountKeeper       types.AccountKeeper
	bankKeeper          types.BankKeeper
	tokenFactoryKeeper  types.TokenFactoryKeeper
	dexKeeper           types.DexKeeper
	authority           string
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

func (k *Keeper) SetTokenFactoryKeeper(tfk types.TokenFactoryKeeper) { k.tokenFactoryKeeper = tfk }
func (k *Keeper) SetDexKeeper(dk types.DexKeeper)                     { k.dexKeeper = dk }

// ============================================================
// Launch CRUD
// ============================================================

func (k Keeper) GetLaunch(ctx context.Context, launchID uint64) (types.Launch, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.LaunchKey(launchID))
	if err != nil || bz == nil { return types.Launch{}, false }
	var launch types.Launch
	if err := json.Unmarshal(bz, &launch); err != nil { return types.Launch{}, false }
	return launch, true
}

func (k Keeper) SetLaunch(ctx context.Context, launch types.Launch) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(launch)
	_ = kvStore.Set(types.LaunchKey(launch.ID), bz)
}

func (k Keeper) GetNextLaunchID(ctx context.Context) uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.NextLaunchIDKey))
	if err != nil || bz == nil { return 1 }
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetNextLaunchID(ctx context.Context, id uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	_ = kvStore.Set([]byte(types.NextLaunchIDKey), bz)
}

// ============================================================
// Contribution CRUD
// ============================================================

func (k Keeper) GetContribution(ctx context.Context, launchID uint64, address string) (types.Contribution, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.ContributionKey(launchID, address))
	if err != nil || bz == nil { return types.Contribution{}, false }
	var contrib types.Contribution
	if err := json.Unmarshal(bz, &contrib); err != nil { return types.Contribution{}, false }
	return contrib, true
}

func (k Keeper) SetContribution(ctx context.Context, contrib types.Contribution) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(contrib)
	_ = kvStore.Set(types.ContributionKey(contrib.LaunchID, contrib.Address), bz)
	// Also set the addr->launch index
	_ = kvStore.Set(types.ContribByAddrKey(contrib.Address, contrib.LaunchID), []byte{1})
}

// ============================================================
// Vesting CRUD
// ============================================================

func (k Keeper) GetVesting(ctx context.Context, launchID uint64, address string) (types.VestingPosition, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.VestingKey(launchID, address))
	if err != nil || bz == nil { return types.VestingPosition{}, false }
	var vp types.VestingPosition
	if err := json.Unmarshal(bz, &vp); err != nil { return types.VestingPosition{}, false }
	return vp, true
}

func (k Keeper) SetVesting(ctx context.Context, vp types.VestingPosition) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(vp)
	_ = kvStore.Set(types.VestingKey(vp.LaunchID, vp.Address), bz)
}

// ============================================================
// Active Launch Index
// ============================================================

func (k Keeper) SetActiveLaunch(ctx context.Context, launchID uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	_ = kvStore.Set(types.ActiveLaunchKey(launchID), []byte{1})
}

func (k Keeper) RemoveActiveLaunch(ctx context.Context, launchID uint64) {
	kvStore := k.storeService.OpenKVStore(ctx)
	_ = kvStore.Delete(types.ActiveLaunchKey(launchID))
}

func (k Keeper) GetAllActiveLaunchIDs(ctx context.Context) []uint64 {
	kvStore := k.storeService.OpenKVStore(ctx)
	iter, err := kvStore.Iterator([]byte(types.ActiveLaunchPrefix), append([]byte(types.ActiveLaunchPrefix), 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
	if err != nil { return nil }
	defer iter.Close()

	var ids []uint64
	prefix := []byte(types.ActiveLaunchPrefix)
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if len(key) >= len(prefix)+8 {
			id := binary.BigEndian.Uint64(key[len(prefix):])
			ids = append(ids, id)
		}
	}
	return ids
}

// ============================================================
// Genesis Helpers
// ============================================================

func (k Keeper) GetAllLaunches(ctx context.Context) []types.Launch {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.LaunchPrefix)
	iter, err := kvStore.Iterator(prefix, append(prefix, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
	if err != nil { return nil }
	defer iter.Close()

	var launches []types.Launch
	for ; iter.Valid(); iter.Next() {
		var l types.Launch
		if err := json.Unmarshal(iter.Value(), &l); err == nil {
			launches = append(launches, l)
		}
	}
	return launches
}

func (k Keeper) GetAllContributions(ctx context.Context) []types.Contribution {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.ContributionPrefix)
	iter, err := kvStore.Iterator(prefix, append(prefix, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
	if err != nil { return nil }
	defer iter.Close()
	var items []types.Contribution
	for ; iter.Valid(); iter.Next() {
		var c types.Contribution
		if err := json.Unmarshal(iter.Value(), &c); err == nil {
			items = append(items, c)
		}
	}
	return items
}

func (k Keeper) GetAllVestings(ctx context.Context) []types.VestingPosition {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.VestingPrefix)
	iter, err := kvStore.Iterator(prefix, append(prefix, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
	if err != nil { return nil }
	defer iter.Close()
	var items []types.VestingPosition
	for ; iter.Valid(); iter.Next() {
		var v types.VestingPosition
		if err := json.Unmarshal(iter.Value(), &v); err == nil {
			items = append(items, v)
		}
	}
	return items
}

// ============================================================
// Core Logic
// ============================================================

func (k Keeper) ExecuteCreateLaunch(ctx context.Context, msg *types.MsgCreateLaunch) (uint64, error) {
	// C-6: validate that creator is a valid bech32 address
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return 0, fmt.Errorf("invalid creator address: %w", err)
	}

	// Validate token supply and price are positive
	if msg.TokenSupply.IsNil() || !msg.TokenSupply.IsPositive() {
		return 0, fmt.Errorf("token supply must be positive")
	}
	if msg.PricePerToken.IsNil() || !msg.PricePerToken.IsPositive() {
		return 0, fmt.Errorf("price per token must be positive")
	}

	// Validate caps
	if msg.SoftCap.IsNil() || !msg.SoftCap.IsPositive() {
		return 0, fmt.Errorf("soft cap must be positive")
	}
	if msg.HardCap.IsNil() || msg.HardCap.LT(msg.SoftCap) {
		return 0, fmt.Errorf("hard cap must be >= soft cap")
	}
	if msg.MaxPerWallet.IsNil() || !msg.MaxPerWallet.IsPositive() {
		return 0, fmt.Errorf("max per wallet must be positive")
	}

	// Validate block range
	if msg.StartBlock <= 0 || msg.EndBlock <= msg.StartBlock {
		return 0, fmt.Errorf("invalid block range: start=%d end=%d", msg.StartBlock, msg.EndBlock)
	}

	// Validate TGE percent: must be 0-100
	if msg.TGEPercent > 100 {
		return 0, fmt.Errorf("TGE percent must be 0-100, got %d", msg.TGEPercent)
	}

	// Validate vesting: if TGE < 100%, vesting blocks must be positive
	if msg.TGEPercent < 100 && msg.VestingBlocks <= 0 {
		return 0, fmt.Errorf("vesting blocks must be positive when TGE percent < 100")
	}

	// Validate vesting blocks is non-negative
	if msg.VestingBlocks < 0 {
		return 0, fmt.Errorf("vesting blocks must be non-negative")
	}

	// Validate launch economics: HardCap must not exceed TokenSupply * PricePerToken.
	// This prevents creators from collecting more funds than the tokens are worth.
	maxRaise := msg.TokenSupply.Mul(msg.PricePerToken)
	if msg.HardCap.GT(maxRaise) {
		return 0, fmt.Errorf("hard cap %s exceeds max raise %s (token_supply * price_per_token)", msg.HardCap, maxRaise)
	}

	// Validate token denom is not empty
	if msg.TokenDenom == "" {
		return 0, fmt.Errorf("token denom must not be empty")
	}
	if msg.QuoteDenom == "" {
		return 0, fmt.Errorf("quote denom must not be empty")
	}

	launchID := k.GetNextLaunchID(ctx)
	k.SetNextLaunchID(ctx, launchID+1)

	// Default TGE if no vesting
	tgePercent := msg.TGEPercent
	if msg.VestingBlocks == 0 {
		tgePercent = 100
	}

	launch := types.Launch{
		ID:            launchID,
		Creator:       msg.Creator,
		TokenDenom:    msg.TokenDenom,
		TokenSupply:   msg.TokenSupply,
		PricePerToken: msg.PricePerToken,
		QuoteDenom:    msg.QuoteDenom,
		SoftCap:       msg.SoftCap,
		HardCap:       msg.HardCap,
		MaxPerWallet:  msg.MaxPerWallet,
		StartBlock:    msg.StartBlock,
		EndBlock:      msg.EndBlock,
		Status:        types.LaunchStatusPending,
		TotalRaised:   math.ZeroInt(),
		Contributors:  0,
		VestingBlocks: msg.VestingBlocks,
		TGEPercent:    tgePercent,
	}

	k.SetLaunch(ctx, launch)
	k.SetActiveLaunch(ctx, launchID)

	return launchID, nil
}

func (k Keeper) ExecuteContribute(ctx context.Context, sender string, launchID uint64, amount math.Int) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentBlock := sdkCtx.BlockHeight()

	launch, found := k.GetLaunch(ctx, launchID)
	if !found { return types.ErrLaunchNotFound }

	// Activate launch if we've reached startBlock and it's still pending
	if launch.Status == types.LaunchStatusPending && currentBlock >= launch.StartBlock {
		launch.Status = types.LaunchStatusActive
	}

	if launch.Status != types.LaunchStatusActive { return types.ErrLaunchNotActive }
	if currentBlock < launch.StartBlock { return types.ErrLaunchNotStarted }
	if currentBlock > launch.EndBlock { return types.ErrLaunchEnded }

	// Check hard cap
	newTotal := launch.TotalRaised.Add(amount)
	if newTotal.GT(launch.HardCap) { return types.ErrHardCapExceeded }

	// Check max per wallet
	contrib, _ := k.GetContribution(ctx, launchID, sender)
	var newContrib math.Int
	if !contrib.Amount.IsNil() && contrib.Amount.IsPositive() {
		newContrib = contrib.Amount.Add(amount)
	} else {
		newContrib = amount
	}
	if newContrib.GT(launch.MaxPerWallet) { return types.ErrMaxPerWalletExceeded }

	// Transfer quote tokens from sender to module
	senderAddr, _ := sdk.AccAddressFromBech32(sender)
	coins := sdk.NewCoins(sdk.NewCoin(launch.QuoteDenom, amount))
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, types.ModuleName, coins); err != nil {
		return types.ErrInsufficientFunds
	}

	// Update contribution
	isNew := contrib.Amount.IsNil() || contrib.Amount.IsZero()
	contrib.Address = sender
	contrib.LaunchID = launchID
	contrib.Amount = newContrib
	k.SetContribution(ctx, contrib)

	// Update launch
	launch.TotalRaised = newTotal
	if isNew {
		launch.Contributors++
	}
	k.SetLaunch(ctx, launch)

	return nil
}

func (k Keeper) ExecuteFinalizeLaunch(ctx context.Context, authority string, launchID uint64) (string, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentBlock := sdkCtx.BlockHeight()

	launch, found := k.GetLaunch(ctx, launchID)
	if !found { return "", types.ErrLaunchNotFound }

	// M-8: only chain governance authority or the launch creator may finalize
	if authority != k.authority && authority != launch.Creator {
		return "", types.ErrUnauthorized
	}

	// M-8: launch must have ended before it can be finalized
	if currentBlock < launch.EndBlock {
		return "", types.ErrLaunchNotEnded
	}

	if launch.Status == types.LaunchStatusFinalized || launch.Status == types.LaunchStatusSuccessful || launch.Status == types.LaunchStatusFailed {
		return "", types.ErrLaunchAlreadyFinalized
	}

	// Check if soft cap was met
	if launch.TotalRaised.GTE(launch.SoftCap) {
		launch.Status = types.LaunchStatusSuccessful
		launch.FinalizedBlock = currentBlock

		// Token creation and DEX pool creation are handled by external keepers
		// In tests we mock these; in production the tokenfactory and dex keepers handle it
		if k.tokenFactoryKeeper != nil {
			moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
			// Create denom via tokenfactory
			fullDenom, err := k.tokenFactoryKeeper.CreateDenom(ctx, moduleAddr.String(), launch.TokenDenom)
			if err != nil {
				return "", err
			}
			launch.FullTokenDenom = fullDenom

			// Mint the full token supply to the module account
			mintCoin := sdk.NewCoin(fullDenom, launch.TokenSupply)
			if err := k.tokenFactoryKeeper.Mint(ctx, moduleAddr.String(), mintCoin, moduleAddr.String()); err != nil {
				return "", err
			}

			// Create DEX pool with a portion of raised funds and tokens
			// Use 20% of raised funds and proportional tokens for initial liquidity
			if k.dexKeeper != nil {
				liquidityQuote := launch.TotalRaised.Quo(math.NewInt(5)) // 20%
				liquidityTokens := launch.TokenSupply.Quo(math.NewInt(5))
				if liquidityQuote.IsPositive() && liquidityTokens.IsPositive() {
					// Send quote tokens from module to module address for pool creation
					poolID, err := k.dexKeeper.CreatePool(ctx, moduleAddr.String(), fullDenom, launch.QuoteDenom, liquidityTokens, liquidityQuote)
					if err == nil {
						launch.DexPoolID = poolID
					}
				}
			}

			// BUG-6 fix: send the remaining 80% of raised funds to the launch creator.
			// Without this, the funds stay stuck in the module account forever.
			creatorRaised := launch.TotalRaised.Sub(launch.TotalRaised.Quo(math.NewInt(5))) // 80%
			if creatorRaised.IsPositive() {
				creatorAddr, creatorErr := sdk.AccAddressFromBech32(launch.Creator)
				if creatorErr == nil {
					creatorCoins := sdk.NewCoins(sdk.NewCoin(launch.QuoteDenom, creatorRaised))
					if sendErr := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, creatorCoins); sendErr != nil {
						return "", fmt.Errorf("failed to send raised funds to creator: %w", sendErr)
					}
				}
			}
		}
	} else {
		launch.Status = types.LaunchStatusFailed
		launch.FinalizedBlock = currentBlock
	}

	k.SetLaunch(ctx, launch)
	k.RemoveActiveLaunch(ctx, launchID)

	return launch.Status, nil
}

func (k Keeper) ExecuteClaimTokens(ctx context.Context, sender string, launchID uint64) (math.Int, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentBlock := sdkCtx.BlockHeight()

	launch, found := k.GetLaunch(ctx, launchID)
	if !found { return math.Int{}, types.ErrLaunchNotFound }
	if launch.Status != types.LaunchStatusSuccessful && launch.Status != types.LaunchStatusFinalized {
		return math.Int{}, types.ErrLaunchNotSuccessful
	}

	contrib, found := k.GetContribution(ctx, launchID, sender)
	if !found || contrib.Amount.IsNil() || contrib.Amount.IsZero() {
		return math.Int{}, types.ErrNoContribution
	}

	// Calculate tokens owed: (contribution / totalRaised) * tokenSupply
	// But cap distribution at tokenSupply (minus liquidity portion if DEX pool was created)
	distributableTokens := launch.TokenSupply
	if launch.DexPoolID > 0 {
		// 20% went to DEX pool
		distributableTokens = launch.TokenSupply.Mul(math.NewInt(4)).Quo(math.NewInt(5))
	}

	// Guard: totalRaised must be positive to avoid division by zero
	if launch.TotalRaised.IsNil() || launch.TotalRaised.IsZero() || !launch.TotalRaised.IsPositive() {
		return math.Int{}, types.ErrNothingToClaim
	}

	// Proportional allocation: tokens = (userContrib * distributableTokens) / totalRaised
	totalTokensOwed := contrib.Amount.Mul(distributableTokens).Quo(launch.TotalRaised)

	// Safety cap: never exceed distributable tokens
	if totalTokensOwed.GT(distributableTokens) {
		totalTokensOwed = distributableTokens
	}

	if totalTokensOwed.IsZero() {
		return math.Int{}, types.ErrNothingToClaim
	}

	// Handle vesting
	if launch.VestingBlocks > 0 && launch.TGEPercent < 100 {
		// Check if vesting position exists
		vp, vpFound := k.GetVesting(ctx, launchID, sender)
		if !vpFound {
			// First claim - create vesting position and send TGE amount
			tgeAmount := totalTokensOwed.Mul(math.NewInt(int64(launch.TGEPercent))).Quo(math.NewInt(100))
			vestingTokens := totalTokensOwed.Sub(tgeAmount)
			vestingPerBlock := math.ZeroInt()
			if launch.VestingBlocks > 0 {
				vestingPerBlock = vestingTokens.Quo(math.NewInt(launch.VestingBlocks))
			}

			vp = types.VestingPosition{
				Address:        sender,
				LaunchID:       launchID,
				TotalTokens:    totalTokensOwed,
				Claimed:        tgeAmount,
				TGEAmount:      tgeAmount,
				VestingPerBlock: vestingPerBlock,
				VestStart:      launch.FinalizedBlock,
			}
			k.SetVesting(ctx, vp)

			// Send TGE tokens
			if tgeAmount.IsPositive() {
				senderAddr, _ := sdk.AccAddressFromBech32(sender)
				tokenDenom := launch.FullTokenDenom
				if tokenDenom == "" {
					tokenDenom = launch.TokenDenom
				}
				coins := sdk.NewCoins(sdk.NewCoin(tokenDenom, tgeAmount))
				if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, coins); err != nil {
					return math.Int{}, err
				}
			}

			contrib.Claimed = true
			contrib.ClaimedAt = currentBlock
			k.SetContribution(ctx, contrib)

			return tgeAmount, nil
		}

		// Subsequent claim - calculate vested amount
		blocksSinceStart := currentBlock - vp.VestStart
		if blocksSinceStart <= 0 { return math.Int{}, types.ErrVestingNotStarted }

		totalVested := vp.TGEAmount.Add(vp.VestingPerBlock.Mul(math.NewInt(blocksSinceStart)))
		if totalVested.GT(vp.TotalTokens) {
			totalVested = vp.TotalTokens
		}

		claimable := totalVested.Sub(vp.Claimed)
		if claimable.IsZero() || !claimable.IsPositive() {
			return math.Int{}, types.ErrNothingToClaim
		}

		// Send vested tokens
		senderAddr, _ := sdk.AccAddressFromBech32(sender)
		tokenDenom := launch.FullTokenDenom
		if tokenDenom == "" {
			tokenDenom = launch.TokenDenom
		}
		coins := sdk.NewCoins(sdk.NewCoin(tokenDenom, claimable))
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, coins); err != nil {
			return math.Int{}, err
		}

		vp.Claimed = vp.Claimed.Add(claimable)
		k.SetVesting(ctx, vp)

		return claimable, nil
	}

	// No vesting - send all tokens at once
	if contrib.Claimed { return math.Int{}, types.ErrAlreadyClaimed }

	senderAddr, _ := sdk.AccAddressFromBech32(sender)
	tokenDenom := launch.FullTokenDenom
	if tokenDenom == "" {
		tokenDenom = launch.TokenDenom
	}
	coins := sdk.NewCoins(sdk.NewCoin(tokenDenom, totalTokensOwed))
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, coins); err != nil {
		return math.Int{}, err
	}

	contrib.Claimed = true
	contrib.ClaimedAt = currentBlock
	k.SetContribution(ctx, contrib)

	return totalTokensOwed, nil
}

func (k Keeper) ExecuteClaimRefund(ctx context.Context, sender string, launchID uint64) (math.Int, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	launch, found := k.GetLaunch(ctx, launchID)
	if !found { return math.Int{}, types.ErrLaunchNotFound }
	if launch.Status != types.LaunchStatusFailed {
		return math.Int{}, types.ErrLaunchNotFailed
	}

	contrib, found := k.GetContribution(ctx, launchID, sender)
	if !found || contrib.Amount.IsNil() || contrib.Amount.IsZero() {
		return math.Int{}, types.ErrNoContribution
	}
	if contrib.Claimed { return math.Int{}, types.ErrAlreadyClaimed }

	// Send refund from module to user
	senderAddr, _ := sdk.AccAddressFromBech32(sender)
	coins := sdk.NewCoins(sdk.NewCoin(launch.QuoteDenom, contrib.Amount))
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, coins); err != nil {
		return math.Int{}, err
	}

	contrib.Claimed = true
	contrib.ClaimedAt = sdkCtx.BlockHeight()
	k.SetContribution(ctx, contrib)

	return contrib.Amount, nil
}

// CheckAndFinalizeExpiredLaunches is called from BeginBlock to auto-finalize expired launches
func (k Keeper) CheckAndFinalizeExpiredLaunches(ctx context.Context) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentBlock := sdkCtx.BlockHeight()

	activeIDs := k.GetAllActiveLaunchIDs(ctx)
	for _, launchID := range activeIDs {
		launch, found := k.GetLaunch(ctx, launchID)
		if !found { continue }

		// Activate pending launches
		if launch.Status == types.LaunchStatusPending && currentBlock >= launch.StartBlock {
			launch.Status = types.LaunchStatusActive
			k.SetLaunch(ctx, launch)
		}

		// Finalize expired launches
		if currentBlock > launch.EndBlock && (launch.Status == types.LaunchStatusActive || launch.Status == types.LaunchStatusPending) {
			k.ExecuteFinalizeLaunch(ctx, k.authority, launchID)
		}
	}
}
