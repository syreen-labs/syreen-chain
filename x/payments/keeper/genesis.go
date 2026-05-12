package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/payments/types"
)

// InitGenesis initializes the module's state from a provided genesis state
func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	k.SetParams(ctx, data.Params)
	k.SetInvoiceCount(ctx, data.InvoiceCount)

	for _, invoice := range data.Invoices {
		k.SetInvoice(ctx, invoice)
	}
	for _, rate := range data.ExchangeRates {
		k.SetExchangeRate(ctx, rate)
	}
	for _, earnings := range data.Earnings {
		k.SetEarnings(ctx, earnings)
	}
}

// ExportGenesis returns the module's exported genesis
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params:        k.GetParams(ctx),
		Invoices:      k.GetAllInvoices(ctx),
		ExchangeRates: k.GetAllExchangeRates(ctx),
		Earnings:      k.GetAllEarnings(ctx),
		InvoiceCount:  k.GetInvoiceCount(ctx),
	}
}
