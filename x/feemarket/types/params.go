package types

import (
	"fmt"

	"cosmossdk.io/math"
)

// Default parameter values
var (
	DefaultBaseFee                = math.LegacyNewDecWithPrec(1, 2) // 0.01 usyreen per gas
	DefaultMinBaseFee             = math.LegacyNewDecWithPrec(1, 4) // 0.0001 usyreen per gas
	DefaultMaxBaseFee             = math.LegacyNewDec(1000)         // 1000 usyreen per gas
	DefaultBaseFeeChangeDenom     = uint64(8)                        // EIP-1559: max 12.5% change per block
	DefaultElasticityMultiplier   = uint64(2)                        // EIP-1559: target = max / 2
	DefaultBurnRatio              = math.LegacyNewDecWithPrec(8, 1) // 0.8 = 80% burned
	DefaultEnableFeeBurn          = true
)

// Params defines the parameters for the feemarket module
type Params struct {
	DefaultBaseFee             math.LegacyDec `json:"default_base_fee"`
	MinBaseFee                 math.LegacyDec `json:"min_base_fee"`
	MaxBaseFee                 math.LegacyDec `json:"max_base_fee"`
	BaseFeeChangeDenominator   uint64         `json:"base_fee_change_denominator"`
	ElasticityMultiplier       uint64         `json:"elasticity_multiplier"`
	BurnRatio                  math.LegacyDec `json:"burn_ratio"`
	EnableFeeBurn              bool           `json:"enable_fee_burn"`
	FeeLanes                   []FeeLane      `json:"fee_lanes"`
}

// DefaultParams returns the default feemarket parameters
func DefaultParams() Params {
	return Params{
		DefaultBaseFee:           DefaultBaseFee,
		MinBaseFee:               DefaultMinBaseFee,
		MaxBaseFee:               DefaultMaxBaseFee,
		BaseFeeChangeDenominator: DefaultBaseFeeChangeDenom,
		ElasticityMultiplier:     DefaultElasticityMultiplier,
		BurnRatio:                DefaultBurnRatio,
		EnableFeeBurn:            DefaultEnableFeeBurn,
		FeeLanes:                 DefaultFeeLanes(),
	}
}

// DefaultFeeLanes returns the default fee lane configurations
func DefaultFeeLanes() []FeeLane {
	return []FeeLane{
		{
			Name:               "default",
			BaseFeeMultiplier:  math.LegacyNewDec(1),
			MaxBlockGas:        50_000_000,
			CurrentUtilization: math.LegacyZeroDec(),
		},
		{
			Name:               "defi",
			BaseFeeMultiplier:  math.LegacyNewDecWithPrec(15, 1), // 1.5x multiplier
			MaxBlockGas:        30_000_000,
			CurrentUtilization: math.LegacyZeroDec(),
		},
		{
			Name:               "ibc",
			BaseFeeMultiplier:  math.LegacyNewDecWithPrec(12, 1), // 1.2x multiplier
			MaxBlockGas:        20_000_000,
			CurrentUtilization: math.LegacyZeroDec(),
		},
		{
			Name:               "governance",
			BaseFeeMultiplier:  math.LegacyNewDecWithPrec(8, 1), // 0.8x multiplier (discounted)
			MaxBlockGas:        10_000_000,
			CurrentUtilization: math.LegacyZeroDec(),
		},
	}
}

// Validate validates the feemarket parameters
func (p Params) Validate() error {
	if p.DefaultBaseFee.IsNegative() {
		return fmt.Errorf("default base fee cannot be negative: %s", p.DefaultBaseFee)
	}
	if p.MinBaseFee.IsNegative() {
		return fmt.Errorf("min base fee cannot be negative: %s", p.MinBaseFee)
	}
	if p.MaxBaseFee.IsNegative() || p.MaxBaseFee.IsZero() {
		return fmt.Errorf("max base fee must be positive: %s", p.MaxBaseFee)
	}
	if p.MinBaseFee.GT(p.MaxBaseFee) {
		return fmt.Errorf("min base fee (%s) cannot be greater than max base fee (%s)", p.MinBaseFee, p.MaxBaseFee)
	}
	if p.BaseFeeChangeDenominator == 0 {
		return fmt.Errorf("base fee change denominator cannot be zero")
	}
	if p.ElasticityMultiplier == 0 {
		return fmt.Errorf("elasticity multiplier cannot be zero")
	}
	if p.BurnRatio.IsNegative() || p.BurnRatio.GT(math.LegacyOneDec()) {
		return fmt.Errorf("burn ratio must be between 0 and 1: %s", p.BurnRatio)
	}
	// H-11: Enforce minimum burn ratio of 50% to prevent governance from
	// disabling the deflationary mechanism entirely.
	if p.BurnRatio.LT(math.LegacyNewDecWithPrec(5, 1)) {
		return fmt.Errorf("burn_ratio must be at least 0.5 (50%%), got %s", p.BurnRatio)
	}
	for _, lane := range p.FeeLanes {
		if lane.Name == "" {
			return fmt.Errorf("fee lane name cannot be empty")
		}
		if lane.BaseFeeMultiplier.IsNegative() || lane.BaseFeeMultiplier.IsZero() {
			return fmt.Errorf("fee lane %s: base fee multiplier must be positive: %s", lane.Name, lane.BaseFeeMultiplier)
		}
		if lane.MaxBlockGas <= 0 {
			return fmt.Errorf("fee lane %s: max block gas must be positive: %d", lane.Name, lane.MaxBlockGas)
		}
	}
	return nil
}
