package types

import (
	"context"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected account keeper interface
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	GetModuleAddress(moduleName string) sdk.AccAddress
}

// BankKeeper defines the expected bank keeper interface
type BankKeeper interface {
	SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

// DexKeeper defines the expected DEX keeper interface for trading intent execution.
// GetSpotPrice returns the spot price of denomIn in terms of denomOut for the given pool.
// Swap executes a swap through the DEX on behalf of sender.
type DexKeeper interface {
	GetSpotPrice(ctx context.Context, poolID uint64, denomIn, denomOut string) (math.LegacyDec, error)
	Swap(ctx context.Context, sender string, poolID uint64, tokenIn sdk.Coin, minTokenOut math.Int) (sdk.Coin, error)

	// AI signal/sentiment getters — used by AI-reactive strategy templates.
	// GetSignalForPool returns (signal label, composite score 0-100, RSI).
	GetSignalForPool(ctx context.Context, poolID uint64) (signal string, compositeScore int64, rsi math.LegacyDec)
	// GetSignalVolatility returns the current volatility score for a pool.
	GetSignalVolatility(ctx context.Context, poolID uint64) math.LegacyDec
	// GetFearGreedIndex returns (index 0-100, label, ok).
	GetFearGreedIndex(ctx context.Context, poolID uint64) (index int64, label string, ok bool)
}

// TransferKeeper defines the expected IBC transfer keeper interface for cross-chain intents.
type TransferKeeper interface {
	Transfer(ctx context.Context, msg *IBCTransferMsg) (*IBCTransferResponse, error)
}

// IBCTransferMsg is a minimal representation of ibc MsgTransfer to avoid importing ibc-go directly.
type IBCTransferMsg struct {
	SourcePort       string   `json:"source_port"`
	SourceChannel    string   `json:"source_channel"`
	Token            sdk.Coin `json:"token"`
	Sender           string   `json:"sender"`
	Receiver         string   `json:"receiver"`
	TimeoutHeight    uint64   `json:"timeout_height"`
	TimeoutTimestamp uint64   `json:"timeout_timestamp"`
}

type IBCTransferResponse struct {
	Sequence uint64 `json:"sequence"`
}

// InterfaceRegistry is used for unpacking execution messages from Any format
type InterfaceRegistry interface {
	types.InterfaceRegistry
}
