package types

// GenesisState defines the payments module's genesis state.
type GenesisState struct {
	Params        Params         `json:"params"`
	Invoices      []Invoice      `json:"invoices"`
	ExchangeRates []ExchangeRate `json:"exchange_rates"`
	Earnings      []Earnings     `json:"earnings"`
	InvoiceCount  uint64         `json:"invoice_count"`
}

// DefaultGenesis returns the default genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:        DefaultParams(),
		Invoices:      []Invoice{},
		ExchangeRates: []ExchangeRate{},
		Earnings:      []Earnings{},
		InvoiceCount:  0,
	}
}

// Validate performs basic genesis state validation.
func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}
