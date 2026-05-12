package ante

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	dextypes "syreen/x/dex/types"
	feemarketante "syreen/x/feemarket/ante"
)

// GaslessDexDecorator checks if a transaction contains ONLY dex trading messages.
// If so, it sets the gas price to zero so the user pays no gas.
// The trading fee (0.3%) is still collected by the swap logic.
//
// Anti-spam protection:
// - Sender must have a minimum balance (MinGaslessBalance) to prevent zero-cost spam
// - Gas is capped at MaxGaslessGas per tx
//
// Supported gasless message types:
// - MsgSwap
// - MsgPlaceOrder
// - MsgCancelOrder
// - MsgModifyOrder
// - MsgMultiHopSwap
//
// Pool creation, liquidity add/remove, and referral operations still require gas.

const (
	// MaxGaslessGas is the maximum gas allowed for a gasless tx to prevent abuse.
	MaxGaslessGas uint64 = 1_000_000

	// MinGaslessBalanceUsyreen is the minimum account balance (in usyreen) required
	// to qualify for gasless transactions. This prevents zero-cost DoS spam.
	// 1 SYR = 1_000_000 usyreen — so attacker needs at least 1 SYR to use gasless.
	MinGaslessBalanceUsyreen int64 = 1_000_000
)

// GaslessDexDecorator implements sdk.AnteDecorator.
type GaslessDexDecorator struct {
	bankKeeper BankKeeper
}

// BankKeeper is the minimal interface needed by the gasless decorator.
type BankKeeper interface {
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

// NewGaslessDexDecorator creates a new GaslessDexDecorator.
func NewGaslessDexDecorator(bk BankKeeper) GaslessDexDecorator {
	return GaslessDexDecorator{bankKeeper: bk}
}

// AnteHandle checks if all messages in the tx are gasless-eligible DEX messages.
// If so, it replaces the gas meter with a free (but capped) meter so the user
// pays no gas fees. Otherwise, it passes through to the next decorator.
func (gd GaslessDexDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	msgs := tx.GetMsgs()
	if len(msgs) == 0 {
		return next(ctx, tx, simulate)
	}

	// Check if ALL messages are gasless-eligible
	allGasless := true
	for _, msg := range msgs {
		if !isGaslessEligible(msg) {
			allGasless = false
			break
		}
	}

	if !allGasless {
		return next(ctx, tx, simulate)
	}

	// Anti-spam: ALL signers must hold a minimum balance to qualify for gasless.
	// This prevents zero-cost DoS attacks while keeping trading free.
	// Checking only the first signer would let additional signers bypass the
	// balance requirement, so we iterate over every message's signer address.
	if gd.bankKeeper != nil {
		seen := make(map[string]bool)
		for _, msg := range msgs {
			signerAddr := extractGaslessMsgSigner(msg)
			if signerAddr == "" {
				return next(ctx, tx, simulate)
			}
			if seen[signerAddr] {
				continue
			}
			seen[signerAddr] = true
			signer, err := sdk.AccAddressFromBech32(signerAddr)
			if err != nil {
				return next(ctx, tx, simulate)
			}
			bal := gd.bankKeeper.GetBalance(ctx, signer, "usyreen")
			if bal.Amount.LT(math.NewInt(MinGaslessBalanceUsyreen)) {
				// Signer doesn't have minimum balance — charge normal gas
				return next(ctx, tx, simulate)
			}
		}
	}

	// All messages are gasless DEX trading messages and sender has min balance.
	// Replace the gas meter with a capped free meter and set the gasless flag
	// so the feemarket decorator skips the base fee check.
	gasMeter := storetypes.NewGasMeter(MaxGaslessGas)
	ctx = ctx.WithGasMeter(gasMeter).
		WithMinGasPrices(sdk.DecCoins{}).
		WithValue(feemarketante.GaslessTxKey{}, true)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"gasless_tx",
			sdk.NewAttribute("sponsored", "true"),
			sdk.NewAttribute("msg_count", fmt.Sprintf("%d", len(msgs))),
			sdk.NewAttribute("gas_cap", fmt.Sprintf("%d", MaxGaslessGas)),
		),
	)

	return next(ctx, tx, simulate)
}

// extractGaslessMsgSigner returns the signer address string from a gasless-eligible
// DEX message. Returns empty string if the message type is unrecognized.
func extractGaslessMsgSigner(msg sdk.Msg) string {
	switch m := msg.(type) {
	case *dextypes.MsgSwap:
		return m.Sender
	case *dextypes.MsgPlaceOrder:
		return m.Creator
	case *dextypes.MsgCancelOrder:
		return m.Creator
	case *dextypes.MsgModifyOrder:
		return m.Creator
	case *dextypes.MsgMultiHopSwap:
		return m.Sender
	default:
		return ""
	}
}

// isGaslessEligible returns true if the message is a DEX trading message
// that qualifies for gasless execution.
func isGaslessEligible(msg sdk.Msg) bool {
	switch msg.(type) {
	case *dextypes.MsgSwap:
		return true
	case *dextypes.MsgPlaceOrder:
		return true
	case *dextypes.MsgCancelOrder:
		return true
	case *dextypes.MsgModifyOrder:
		return true
	case *dextypes.MsgMultiHopSwap:
		return true
	default:
		return false
	}
}
