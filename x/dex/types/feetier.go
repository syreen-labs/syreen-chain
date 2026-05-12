package types

import (
	"cosmossdk.io/math"
)

// ---------------------------------------------------------------------------
// Fee Tier Types
// ---------------------------------------------------------------------------

// FeeTier defines a volume-based fee tier with discount percentages.
type FeeTier struct {
	Level            string         `json:"level"`
	MinVolume30d     math.Int       `json:"min_volume_30d"`     // minimum 30-day volume in usyreen
	MakerFeeDiscount math.LegacyDec `json:"maker_fee_discount"` // e.g., 0.10 = 10% discount
	TakerFeeDiscount math.LegacyDec `json:"taker_fee_discount"` // e.g., 0.10 = 10% discount
}

// DailyVolume records volume for a single block-day bucket.
type DailyVolume struct {
	Day    int64    `json:"day"`    // block height / BlocksPerDay
	Volume math.Int `json:"volume"` // cumulative volume for this day
}

// UserVolume tracks a user's rolling 30-day trade volume.
type UserVolume struct {
	Address         string        `json:"address"`
	Volume30d       math.Int      `json:"volume_30d"`
	LastUpdateBlock int64         `json:"last_update_block"`
	DailyVolumes    []DailyVolume `json:"daily_volumes"`
}

// ---------------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------------

const (
	// BlocksPerDay is the number of blocks in one day at 500ms block time.
	BlocksPerDay int64 = 172800

	// RollingWindowDays is the number of days in the rolling volume window.
	RollingWindowDays int64 = 30

	// PruneInterval is how often (in blocks) to run volume pruning in BeginBlock.
	PruneInterval int64 = 1000

	// Fee tier level names
	FeeTierStandard = "standard"
	FeeTierSilver   = "silver"
	FeeTierGold     = "gold"
	FeeTierPlatinum = "platinum"
	FeeTierDiamond  = "diamond"
)

// ---------------------------------------------------------------------------
// KV Store Prefixes
// ---------------------------------------------------------------------------

const (
	// UserVolumePrefix stores per-user volume data: user_volume/<address> -> UserVolume JSON
	UserVolumePrefix = "user_volume/"

	// UserFeeTierPrefix stores per-user computed fee tier: user_feetier/<address> -> tier level string
	UserFeeTierPrefix = "user_feetier/"
)

// ---------------------------------------------------------------------------
// Default Fee Tiers
// ---------------------------------------------------------------------------

// DefaultFeeTiers returns the default volume-based fee tiers.
func DefaultFeeTiers() []FeeTier {
	return []FeeTier{
		{
			Level:            FeeTierStandard,
			MinVolume30d:     math.ZeroInt(),
			MakerFeeDiscount: math.LegacyZeroDec(),
			TakerFeeDiscount: math.LegacyZeroDec(),
		},
		{
			Level:            FeeTierSilver,
			MinVolume30d:     math.NewInt(50_000_000_000),   // 50,000 SYR (6 decimals)
			MakerFeeDiscount: math.LegacyNewDecWithPrec(1, 1), // 0.10 = 10%
			TakerFeeDiscount: math.LegacyNewDecWithPrec(1, 1), // 0.10 = 10%
		},
		{
			Level:            FeeTierGold,
			MinVolume30d:     math.NewInt(250_000_000_000),  // 250,000 SYR (6 decimals)
			MakerFeeDiscount: math.LegacyNewDecWithPrec(2, 1), // 0.20 = 20%
			TakerFeeDiscount: math.LegacyNewDecWithPrec(2, 1), // 0.20 = 20%
		},
		{
			Level:            FeeTierPlatinum,
			MinVolume30d:     math.NewInt(1_000_000_000_000), // 1,000,000 SYR (6 decimals)
			MakerFeeDiscount: math.LegacyNewDecWithPrec(3, 1), // 0.30 = 30%
			TakerFeeDiscount: math.LegacyNewDecWithPrec(3, 1), // 0.30 = 30%
		},
		{
			Level:            FeeTierDiamond,
			MinVolume30d:     math.NewInt(5_000_000_000_000), // 5,000,000 SYR (6 decimals)
			MakerFeeDiscount: math.LegacyNewDecWithPrec(4, 1), // 0.40 = 40%
			TakerFeeDiscount: math.LegacyNewDecWithPrec(4, 1), // 0.40 = 40%
		},
	}
}
