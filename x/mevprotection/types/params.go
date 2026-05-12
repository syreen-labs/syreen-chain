package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
)

// MinCommitWindow is the minimum allowed commit window to prevent MEV via too-short windows.
const MinCommitWindow = uint64(3)

var DefaultFairOrderConfig = FairOrderConfig{
	EnableCommitReveal: true,
	CommitWindow:       3, // 3 blocks (minimum)
	RevealWindow:       1, // 1 block
	MaxTxDelay:         3, // 3 blocks
}

// DefaultSlashFraction is 5% — the fraction of stake slashed when MEV is detected.
const DefaultSlashFraction = "0.050000000000000000"

// Params defines the parameters for the mevprotection module
type Params struct {
	FairOrderConfig FairOrderConfig `json:"fair_order_config"`
	SlashFraction   string          `json:"slash_fraction"` // decimal string, e.g. "0.050000000000000000"
}

func DefaultParams() Params {
	return Params{
		FairOrderConfig: DefaultFairOrderConfig,
		SlashFraction:   DefaultSlashFraction,
	}
}

func (p Params) Validate() error {
	if p.FairOrderConfig.CommitWindow < MinCommitWindow {
		return fmt.Errorf("commit window must be at least %d blocks, got %d", MinCommitWindow, p.FairOrderConfig.CommitWindow)
	}
	if p.FairOrderConfig.RevealWindow == 0 {
		return fmt.Errorf("reveal window must be greater than 0")
	}
	if p.FairOrderConfig.MaxTxDelay == 0 {
		return fmt.Errorf("max tx delay must be greater than 0")
	}
	if p.SlashFraction != "" {
		// Validate that the slash fraction is a parseable decimal between 0 and 1
		// We allow empty for backwards compat (defaults applied at use site)
		parsed, err := parseDecimal(p.SlashFraction)
		if err != nil {
			return fmt.Errorf("invalid slash_fraction %q: %w", p.SlashFraction, err)
		}
		if parsed < 0 || parsed > 1 {
			return fmt.Errorf("slash_fraction must be between 0 and 1, got %s", p.SlashFraction)
		}
		// H-12: Ensure slash fraction is at least 1% to prevent governance
		// from neutralizing MEV penalties entirely.
		if parsed < 0.01 {
			return fmt.Errorf("slash_fraction must be at least 0.01 (1%%), got %s", p.SlashFraction)
		}
	}
	return nil
}

// parseDecimal validates decimal strings using deterministic sdk.Dec arithmetic.
func parseDecimal(s string) (float64, error) {
	d, err := sdkmath.LegacyNewDecFromStr(s)
	if err != nil {
		return 0, err
	}
	// Return as float64 for the comparison callers; the important thing is
	// that validation is deterministic via LegacyDec parsing.
	f, _ := d.Float64()
	return f, nil
}
