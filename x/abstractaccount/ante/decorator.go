package ante

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"syreen/x/abstractaccount/types"
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

	// Sponsor found: deduct fees from the sponsor's account to fee_collector.
	sponsorAddr, err := sdk.AccAddressFromBech32(sponsor.Sponsor)
	if err != nil {
		return ctx, sdkerrors.ErrInvalidAddress.Wrapf("invalid sponsor address: %s", sponsor.Sponsor)
	}

	err = gsd.bankKeeper.SendCoinsFromAccountToModule(ctx, sponsorAddr, authtypes.FeeCollectorName, fees)
	if err != nil {
		// If sponsor cannot pay, fall through to normal fee deduction.
		return next(ctx, tx, simulate)
	}

	// Record gas usage against the sponsor's budget.
	// Use a conservative estimate: min(gasWanted, fee amount) to avoid overcharging.
	// NOTE: A PostHandler should be added for exact accounting using actual gas consumed.
	budgetDeduction := gasWanted
	feeAmount := fees.AmountOf("usyreen")
	if !feeAmount.IsZero() && feeAmount.Uint64() < budgetDeduction {
		budgetDeduction = feeAmount.Uint64()
	}
	gsd.aaKeeper.DeductGasSponsorship(ctx, sponsor, budgetDeduction)

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
