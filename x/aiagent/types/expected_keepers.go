package types

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
	HasAccount(ctx context.Context, addr sdk.AccAddress) bool
	NewAccountWithAddress(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	SetAccount(ctx context.Context, acc sdk.AccountI)
}

type BankKeeper interface {
	SendCoins(ctx context.Context, from, to sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, sender sdk.AccAddress, module string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, module string, recipient sdk.AccAddress, amt sdk.Coins) error
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

type DexKeeper interface {
	GetSpotPrice(ctx context.Context, poolID uint64, denomIn, denomOut string) (math.LegacyDec, error)
	SmartSwap(ctx context.Context, sender string, inputDenom, outputDenom string, inputAmount, minOutputAmount math.Int) (sdk.Coin, error)
	GetSignalForPool(ctx context.Context, poolID uint64) (string, int64, math.LegacyDec) // signal, score, RSI
	GetPoolDenoms(ctx context.Context, poolID uint64) (string, string, bool)
}
