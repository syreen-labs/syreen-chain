package ante

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"syreen/x/abstractaccount/types"
)

const (
	// sponsorFeeDenom is the native denom a gas sponsor may cover.
	sponsorFeeDenom = "usyreen"
	// maxSponsoredGasPrice caps the fee a sponsor can be charged per unit of
	// authorized gas (usyreen/gas). The chain min gas price is 0.001 and the
	// feemarket base fee is ~0.01 usyreen/gas, so this leaves ~1000x headroom for
	// legitimate sponsored txs while bounding total sponsor exposure to
	// GasLimit * maxSponsoredGasPrice. A fixed constant keeps the check
	// deterministic across all nodes (node-local MinGasPrices must not be used here).
	maxSponsoredGasPrice = 10
)

// contextKey is an unexported type for context keys in this package.
type contextKey string

// GasSponsoredKey is the context key indicating fees have already been paid by a sponsor.
const GasSponsoredKey contextKey = "aa_gas_sponsored"

// AbstractAccountKeeper defines the keeper methods needed by the gas sponsorship decorator.
type AbstractAccountKeeper interface {
	CheckGasSponsor(ctx context.Context, sponsored string, gasWanted uint64) (*types.GasSponsor, bool)
	DeductGasSponsorship(ctx context.Context, sponsor *types.GasSponsor, gasUsed uint64)
}

// GasSponsorBankKeeper defines the bank methods needed to transfer fees from the sponsor.
type GasSponsorBankKeeper interface {
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

// GasSponsorshipDecorator checks if the tx sender has an active gas sponsor
// and if so, deducts fees from the sponsor instead of the sender.
type GasSponsorshipDecorator struct {
	aaKeeper   AbstractAccountKeeper
	bankKeeper GasSponsorBankKeeper
}

// NewGasSponsorshipDecorator returns a new GasSponsorshipDecorator.
func NewGasSponsorshipDecorator(aaKeeper AbstractAccountKeeper, bankKeeper GasSponsorBankKeeper) GasSponsorshipDecorator {
	return GasSponsorshipDecorator{
		aaKeeper:   aaKeeper,
		bankKeeper: bankKeeper,
	}
}

func (gsd GasSponsorshipDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.ErrTxDecode.Wrap("tx must implement FeeTx")
	}

	fees := feeTx.GetFee()
	if fees.IsZero() {
		// No fees to sponsor; pass through.
		return next(ctx, tx, simulate)
	}

	feePayer := feeTx.FeePayer()
	feePayerStr := sdk.AccAddress(feePayer).String()

	gasWanted := feeTx.GetGas()

	sponsor, found := gsd.aaKeeper.CheckGasSponsor(ctx, feePayerStr, gasWanted)
	if !found {
		// No active sponsor for this account; pass through to normal fee deduction.
		return next(ctx, tx, simulate)
	}

	// SECURITY: The sponsor's budget is denominated in GAS (GasLimit/TotalGasUsed),
	// but `fees` is an attacker-controlled coin amount. Transferring the full `fees`
	// let a sponsored account drain its sponsor by attaching an arbitrarily large fee
	// to a tiny-gas tx. Bound the sponsorable fee to the gas the sponsor authorized,
	// priced at a fixed maximum gas price. Anything above that (or paid in a denom
	// other than usyreen) is NOT sponsored — the tx falls through to normal fee
	// deduction so the sponsored account pays it themselves.
	feeAmount := fees.AmountOf(sponsorFeeDenom)
	maxSponsoredAmt := sdkmath.NewIntFromUint64(gasWanted).Mul(sdkmath.NewIntFromUint64(maxSponsoredGasPrice))
	if !fees.Equal(sdk.NewCoins(sdk.NewCoin(sponsorFeeDenom, feeAmount))) || feeAmount.GT(maxSponsoredAmt) {
		// Fee is in a non-native denom or exceeds the sponsored allowance.
		return next(ctx, tx, simulate)
	}

	// Sponsor found and fee within the authorized allowance: deduct fees from the
	// sponsor's account to fee_collector.
	sponsorAddr, err := sdk.AccAddressFromBech32(sponsor.Sponsor)
	if err != nil {
		return ctx, sdkerrors.ErrInvalidAddress.Wrapf("invalid sponsor address: %s", sponsor.Sponsor)
	}

	err = gsd.bankKeeper.SendCoinsFromAccountToModule(ctx, sponsorAddr, authtypes.FeeCollectorName, fees)
	if err != nil {
		// If sponsor cannot pay, fall through to normal fee deduction.
		return next(ctx, tx, simulate)
	}

	// Record gas usage against the sponsor's budget. CheckGasSponsor already verified
	// TotalGasUsed+gasWanted <= GasLimit, so deduct the authorized gas (units match the
	// budget). This bounds total sponsor exposure to GasLimit * maxSponsoredGasPrice.
	gsd.aaKeeper.DeductGasSponsorship(ctx, sponsor, gasWanted)

	// Mark the context so downstream DeductFeeDecorator knows fees are paid.
	ctx = ctx.WithValue(GasSponsoredKey, true)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		"gas_sponsored",
		sdk.NewAttribute("sponsor", sponsor.Sponsor),
		sdk.NewAttribute("sponsored", feePayerStr),
		sdk.NewAttribute("fees", fees.String()),
	))

	return next(ctx, tx, simulate)
}
