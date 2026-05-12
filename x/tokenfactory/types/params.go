package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var DefaultDenomCreationFee = sdk.NewCoins() // Free denom creation

// Params defines the parameters for the tokenfactory module
type Params struct {
	DenomCreationFee sdk.Coins `json:"denom_creation_fee"`
}

func DefaultParams() Params {
	return Params{
		DenomCreationFee: DefaultDenomCreationFee,
	}
}

func (p Params) Validate() error {
	if !p.DenomCreationFee.IsValid() {
		return ErrCreationFeeNotEnough
	}
	return nil
}
