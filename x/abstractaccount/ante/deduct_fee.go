package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	feemarketante "syreen/x/feemarket/ante"
)

// ConditionalDeductFeeDecorator wraps the SDK's DeductFeeDecorator
// and skips fee deduction if gas was already sponsored via GasSponsorshipDecorator.
// This prevents the double-fee-deduction bug (H-02) where a sponsored tx would
// have fees taken from both the sponsor AND the sender.
type ConditionalDeductFeeDecorator struct {
	wrapped sdkante.DeductFeeDecorator
}

// NewConditionalDeductFeeDecorator creates a ConditionalDeductFeeDecorator that
// delegates to the standard SDK DeductFeeDecorator only when fees have not
// already been paid by a gas sponsor.
func NewConditionalDeductFeeDecorator(
	ak sdkante.AccountKeeper,
	bk authtypes.BankKeeper,
	fk sdkante.FeegrantKeeper,
	tfc sdkante.TxFeeChecker,
) ConditionalDeductFeeDecorator {
	return ConditionalDeductFeeDecorator{
		wrapped: sdkante.NewDeductFeeDecorator(ak, bk, fk, tfc),
	}
}

func (cdd ConditionalDeductFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// If gas was already sponsored, skip standard fee deduction to avoid
	// charging fees twice (once from sponsor, once from sender).
	if ctx.Value(GasSponsoredKey) != nil {
		return next(ctx, tx, simulate)
	}
	// If this is a gasless DEX transaction, skip fee deduction entirely.
	if ctx.Value(feemarketante.GaslessTxKey{}) != nil {
		return next(ctx, tx, simulate)
	}
	// Normal fee deduction path.
	return cdd.wrapped.AnteHandle(ctx, tx, simulate, next)
}
