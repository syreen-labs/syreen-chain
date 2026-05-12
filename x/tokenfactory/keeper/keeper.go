package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"syreen/x/tokenfactory/types"
)

type Keeper struct {
	cdc          codec.Codec
	storeService store.KVStoreService
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	distrKeeper   types.DistrKeeper
	authority     string
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	distrKeeper types.DistrKeeper,
	authority string,
) Keeper {
	return Keeper{
		cdc:           cdc,
		storeService:  storeService,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		distrKeeper:   distrKeeper,
		authority:     authority,
	}
}

func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// CreateDenom creates a new factory denom
func (k Keeper) CreateDenom(ctx context.Context, creator, subdenom string) (string, error) {
	// L-05: Blacklist native denom names to prevent confusion
	blacklisted := map[string]bool{
		"usyreen": true, "syreen": true, "syr": true, "usyr": true,
		"uatom": true, "atom": true, "stake": true,
	}
	if blacklisted[strings.ToLower(subdenom)] {
		return "", fmt.Errorf("subdenom %q is a reserved native denom name", subdenom)
	}

	denom := types.GetDenom(creator, subdenom)

	// Check if denom already exists
	_, found := k.GetDenomAuthorityMetadata(ctx, denom)
	if found {
		return "", types.ErrDenomExists
	}

	// Validate creator address
	if _, err := sdk.AccAddressFromBech32(creator); err != nil {
		return "", types.ErrInvalidCreator
	}

	// Migrate stored params: remove old creation fee
	params := k.GetParams(ctx)
	if params.DenomCreationFee.IsAllPositive() {
		params.DenomCreationFee = sdk.NewCoins()
		k.SetParams(ctx, params)
	}

	// Set denom metadata in bank module
	k.bankKeeper.SetDenomMetaData(ctx, banktypes.Metadata{
		Description: fmt.Sprintf("Token factory denom created by %s", creator),
		DenomUnits: []*banktypes.DenomUnit{
			{Denom: denom, Exponent: 0, Aliases: nil},
		},
		Base:    denom,
		Display: denom,
		Name:    subdenom,
		Symbol:  subdenom,
	})

	// Store authority metadata
	k.SetDenomAuthorityMetadata(ctx, denom, types.DenomAuthorityMetadata{Admin: creator})

	// Index by creator
	k.addDenomByCreator(ctx, creator, denom)

	k.Logger(ctx).Info("created new factory denom", "denom", denom, "creator", creator)
	return denom, nil
}

// Mint mints tokens of a factory denom
func (k Keeper) Mint(ctx context.Context, sender string, amount sdk.Coin, mintTo string) error {
	authority, found := k.GetDenomAuthorityMetadata(ctx, amount.Denom)
	if !found {
		return types.ErrDenomDoesNotExist
	}
	if authority.Admin != sender {
		return types.ErrUnauthorized
	}

	// Fix #1: enforce MaxSupply cap when set (non-zero)
	if authority.MaxSupply != "" && authority.MaxSupply != "0" {
		maxSupply, ok := math.NewIntFromString(authority.MaxSupply)
		if !ok {
			return fmt.Errorf("invalid max supply value: %s", authority.MaxSupply)
		}
		{
			currentSupply := k.bankKeeper.GetSupply(ctx, amount.Denom)
			if currentSupply.Amount.Add(amount.Amount).GT(maxSupply) {
				return fmt.Errorf("minting %s would exceed max supply of %s (current supply: %s)", amount.Amount, authority.MaxSupply, currentSupply.Amount)
			}
		}
	}

	// Mint to module account
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(amount)); err != nil {
		return err
	}

	// Send to recipient
	recipientAddr := sender
	if mintTo != "" {
		recipientAddr = mintTo
	}
	addr, err := sdk.AccAddressFromBech32(recipientAddr)
	if err != nil {
		return err
	}
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.NewCoins(amount))
}

// Burn burns tokens of a factory denom from the sender's own balance.
// I-07: If burnFrom is specified and differs from sender, the operation is
// rejected to prevent admin from burning tokens from arbitrary addresses.
func (k Keeper) Burn(ctx context.Context, sender string, amount sdk.Coin, burnFrom string) error {
	authority, found := k.GetDenomAuthorityMetadata(ctx, amount.Denom)
	if !found {
		return types.ErrDenomDoesNotExist
	}
	if authority.Admin != sender {
		return types.ErrUnauthorized
	}

	// I-07: Reject burning from another user's address
	if burnFrom != "" && burnFrom != sender {
		return fmt.Errorf("can only burn from your own balance or the module account")
	}

	// Determine the address to burn from
	burnAddr := sender
	if burnFrom != "" {
		burnAddr = burnFrom
	}

	addr, err := sdk.AccAddressFromBech32(burnAddr)
	if err != nil {
		return err
	}

	// Send from burn address to module
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.NewCoins(amount)); err != nil {
		return err
	}

	return k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(amount))
}

// ChangeAdmin changes the admin of a factory denom
func (k Keeper) ChangeAdmin(ctx context.Context, sender, denom, newAdmin string) error {
	authority, found := k.GetDenomAuthorityMetadata(ctx, denom)
	if !found {
		return types.ErrDenomDoesNotExist
	}
	if authority.Admin != sender {
		return types.ErrUnauthorized
	}
	if newAdmin == "" {
		return fmt.Errorf("new admin address cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(newAdmin); err != nil {
		return fmt.Errorf("invalid new admin address: %w", err)
	}

	authority.Admin = newAdmin
	k.SetDenomAuthorityMetadata(ctx, denom, authority)

	k.Logger(ctx).Info("changed denom admin", "denom", denom, "new_admin", newAdmin)
	return nil
}

// GetDenomAuthorityMetadata returns the authority metadata for a denom
func (k Keeper) GetDenomAuthorityMetadata(ctx context.Context, denom string) (types.DenomAuthorityMetadata, bool) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.DenomAuthorityKey(denom))
	if err != nil || bz == nil {
		return types.DenomAuthorityMetadata{}, false
	}
	var metadata types.DenomAuthorityMetadata
	if err := json.Unmarshal(bz, &metadata); err != nil {
		return types.DenomAuthorityMetadata{}, false
	}
	return metadata, true
}

// SetDenomAuthorityMetadata stores the authority metadata for a denom
func (k Keeper) SetDenomAuthorityMetadata(ctx context.Context, denom string, metadata types.DenomAuthorityMetadata) {
	store := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(metadata)
	store.Set(types.DenomAuthorityKey(denom), bz)
}

// GetParams returns the module parameters
func (k Keeper) GetParams(ctx context.Context) types.Params {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get([]byte("params"))
	if err != nil || bz == nil {
		return types.DefaultParams()
	}
	var params types.Params
	if err := json.Unmarshal(bz, &params); err != nil {
		return types.DefaultParams()
	}
	return params
}

// SetParams stores the module parameters
func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(params)
	if err != nil {
		return err
	}
	return store.Set([]byte("params"), bz)
}

// GetDenomsFromCreator returns all denoms created by a specific creator
func (k Keeper) GetDenomsFromCreator(ctx context.Context, creator string) []string {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.CreatorDenomsKey(creator))
	if err != nil || bz == nil {
		return nil
	}
	var denoms []string
	json.Unmarshal(bz, &denoms)
	return denoms
}

func (k Keeper) addDenomByCreator(ctx context.Context, creator, denom string) {
	denoms := k.GetDenomsFromCreator(ctx, creator)
	denoms = append(denoms, denom)
	store := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(denoms)
	store.Set(types.CreatorDenomsKey(creator), bz)
}
