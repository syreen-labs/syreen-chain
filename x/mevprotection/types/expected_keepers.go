package types

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type StakingKeeper interface {
	GetValidator(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error)
	GetValidatorByConsAddr(ctx context.Context, consAddr sdk.ConsAddress) (stakingtypes.Validator, error)
	Slash(ctx context.Context, consAddr sdk.ConsAddress, infractionHeight int64, power int64, slashFactor math.LegacyDec) (math.Int, error)
	// GetBondedValidatorsByPower returns the current set of bonded (active)
	// validators. Used to guard slashing/jailing against halting a chain that
	// has too few validators.
	GetBondedValidatorsByPower(ctx context.Context) ([]stakingtypes.Validator, error)
}

type SlashingKeeper interface {
	Slash(ctx context.Context, consAddr sdk.ConsAddress, fraction math.LegacyDec, power int64, height int64) error
	Jail(ctx context.Context, consAddr sdk.ConsAddress) error
}

type BankKeeper interface {
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}
