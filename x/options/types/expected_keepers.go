package types

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BankKeeper interface {
	SendCoins(ctx context.Context, from, to sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, sender sdk.AccAddress, module string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, module string, recipient sdk.AccAddress, amt sdk.Coins) error
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

type DexKeeper interface {
	GetSpotPrice(ctx context.Context, poolID uint64, denomIn, denomOut string) (math.LegacyDec, error)
	GetPoolDenoms(ctx context.Context, poolID uint64) (string, string, bool)
	GetSignalVolatility(ctx context.Context, poolID uint64) math.LegacyDec
}
