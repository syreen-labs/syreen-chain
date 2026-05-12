package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	DefaultSwapFee          = math.LegacyNewDecWithPrec(3, 3)    // 0.003 = 0.3%
	DefaultProtocolFeeShare = math.LegacyNewDecWithPrec(333, 3)  // 0.333 = 1/3 of swap fee
	DefaultMinInitLiquidity = math.NewInt(1_000_000)              // 1 token minimum (in base units)
	DefaultPoolCreationFee  = sdk.NewCoins()                      // free
)

// Params defines the parameters for the dex module.
type Params struct {
	DefaultSwapFee      math.LegacyDec `json:"default_swap_fee"`
	ProtocolFeeShare    math.LegacyDec `json:"protocol_fee_share"`
	MinInitialLiquidity math.Int       `json:"min_initial_liquidity"`
	PoolCreationFee     sdk.Coins      `json:"pool_creation_fee"`
}

// DefaultParams returns the default dex module parameters.
func DefaultParams() Params {
	return Params{
		DefaultSwapFee:      DefaultSwapFee,
		ProtocolFeeShare:    DefaultProtocolFeeShare,
		MinInitialLiquidity: DefaultMinInitLiquidity,
		PoolCreationFee:     DefaultPoolCreationFee,
	}
}

// Validate performs basic validation of dex module parameters.
func (p Params) Validate() error {
	if p.DefaultSwapFee.IsNegative() {
		return fmt.Errorf("default swap fee must not be negative: %s", p.DefaultSwapFee)
	}
	if p.DefaultSwapFee.GTE(math.LegacyOneDec()) {
		return fmt.Errorf("default swap fee must be less than 1: %s", p.DefaultSwapFee)
	}
	if p.ProtocolFeeShare.IsNegative() {
		return fmt.Errorf("protocol fee share must not be negative: %s", p.ProtocolFeeShare)
	}
	if p.ProtocolFeeShare.GT(math.LegacyOneDec()) {
		return fmt.Errorf("protocol fee share must not exceed 1: %s", p.ProtocolFeeShare)
	}
	if p.MinInitialLiquidity.IsNegative() {
		return fmt.Errorf("min initial liquidity must not be negative: %s", p.MinInitialLiquidity)
	}
	if !p.PoolCreationFee.IsValid() {
		return fmt.Errorf("invalid pool creation fee: %s", p.PoolCreationFee)
	}
	return nil
}
