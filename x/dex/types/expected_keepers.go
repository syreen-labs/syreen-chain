package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type BankKeeper interface {
	SendCoins(ctx context.Context, from, to sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
}

type TokenFactoryKeeper interface {
	CreateDenom(ctx context.Context, creator, subdenom string) (string, error)
	Mint(ctx context.Context, sender string, amount sdk.Coin, mintTo string) error
	Burn(ctx context.Context, sender string, amount sdk.Coin, burnFrom string) error
}

// StakingKeeper defines the expected staking module keeper interface.
// Used by the time-weighted governance system to query delegator staking info.
type StakingKeeper interface {
	GetDelegatorDelegations(ctx context.Context, delegator sdk.AccAddress, maxRetrieve uint16) ([]stakingtypes.Delegation, error)
}
