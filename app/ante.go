package app

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"

	ibcante "github.com/cosmos/ibc-go/v10/modules/core/ante"
	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"

	aaante "syreen/x/abstractaccount/ante"
	abstractaccountkeeper "syreen/x/abstractaccount/keeper"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	dexante "syreen/x/dex/ante"
	evmante "syreen/x/evm/ante"
	evmkeeper "syreen/x/evm/keeper"
	feemarketante "syreen/x/feemarket/ante"
	feemarketkeeper "syreen/x/feemarket/keeper"
	feemarkettypes "syreen/x/feemarket/types"
)

// HandlerOptions extend the SDK's AnteHandler options
type HandlerOptions struct {
	ante.HandlerOptions
	IBCKeeper                 *ibckeeper.Keeper
	FeeMarketKeeper           feemarketkeeper.Keeper
	FeeMarketBankKeeper       feemarkettypes.BankKeeper
	AbstractAccountKeeper     abstractaccountkeeper.Keeper
	AbstractAccountBankKeeper aaante.GasSponsorBankKeeper
	EVMKeeper                 evmkeeper.Keeper
	FullAccountKeeper         authkeeper.AccountKeeper
}

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {
	if options.AccountKeeper == nil {
		return nil, fmt.Errorf("account keeper is required for ante builder")
	}
	if options.BankKeeper == nil {
		return nil, fmt.Errorf("bank keeper is required for ante builder")
	}
	if options.SignModeHandler == nil {
		return nil, fmt.Errorf("sign mode handler is required for ante builder")
	}

	anteDecorators := []sdk.AnteDecorator{
		ante.NewSetUpContextDecorator(),
		// Gasless DEX trading: zero gas for swap, place-order, cancel-order, modify-order, multi-hop-swap
		dexante.NewGaslessDexDecorator(options.BankKeeper.(dexante.BankKeeper)),
		ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		// EVM raw tx handling: verify Ethereum signature, create account if needed
		evmante.NewEVMAccountDecorator(options.FullAccountKeeper, options.EVMKeeper),
		// Session key validation: check if signer is a session key and validate permissions
		aaante.NewSessionKeyDecorator(options.AbstractAccountKeeper),
		// Gas sponsorship: check if sender has an active sponsor, deduct fees from sponsor
		aaante.NewGasSponsorshipDecorator(options.AbstractAccountKeeper, options.AbstractAccountBankKeeper),
		aaante.NewConditionalDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
		// EIP-1559 fee market enforcement: checks base fee and burns the base fee portion
		feemarketante.NewFeeMarketDecorator(options.FeeMarketKeeper, options.FeeMarketBankKeeper),
		// Signature-related decorators: wrapped with EVM skip for raw Ethereum txs
		evmante.NewEVMSkipDecorator(ante.NewSetPubKeyDecorator(options.AccountKeeper)),
		evmante.NewEVMSkipDecorator(ante.NewValidateSigCountDecorator(options.AccountKeeper)),
		evmante.NewEVMSkipDecorator(ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer)),
		evmante.NewEVMSigVerificationDecorator(ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler)),
		// IBC redundant relay check BEFORE sequence increment to avoid wasting sequence numbers
		ibcante.NewRedundantRelayDecorator(options.IBCKeeper),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
	}

	return sdk.ChainAnteDecorators(anteDecorators...), nil
}
