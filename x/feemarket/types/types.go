package types

import (
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// FeeLane represents a fee tier for a specific transaction type
type FeeLane struct {
	Name               string       `json:"name"`
	BaseFeeMultiplier  math.LegacyDec `json:"base_fee_multiplier"`
	MaxBlockGas        int64        `json:"max_block_gas"`
	CurrentUtilization math.LegacyDec `json:"current_utilization"`
}

func (fl *FeeLane) ProtoMessage()           {}
func (fl *FeeLane) Reset()                  { *fl = FeeLane{} }
func (fl *FeeLane) String() string          { return fmt.Sprintf("FeeLane{%s, multiplier=%s}", fl.Name, fl.BaseFeeMultiplier) }
func (fl *FeeLane) XXX_MessageName() string { return "syreen.feemarket.FeeLane" }

// FeeState tracks the current dynamic fee state for the EIP-1559 algorithm
type FeeState struct {
	BaseFee         math.LegacyDec `json:"base_fee"`
	MinBaseFee      math.LegacyDec `json:"min_base_fee"`
	MaxBaseFee      math.LegacyDec `json:"max_base_fee"`
	BlockGasTarget  int64          `json:"block_gas_target"`
	AdjustmentSpeed math.LegacyDec `json:"adjustment_speed"`
}

func (fs *FeeState) ProtoMessage()           {}
func (fs *FeeState) Reset()                  { *fs = FeeState{} }
func (fs *FeeState) String() string          { return fmt.Sprintf("FeeState{baseFee=%s}", fs.BaseFee) }
func (fs *FeeState) XXX_MessageName() string { return "syreen.feemarket.FeeState" }

// BurnRecord tracks fees burned at a specific block height
type BurnRecord struct {
	Height           int64     `json:"height"`
	Amount           sdk.Coins `json:"amount"`
	CumulativeBurned sdk.Coins `json:"cumulative_burned"`
}

func (br *BurnRecord) ProtoMessage()           {}
func (br *BurnRecord) Reset()                  { *br = BurnRecord{} }
func (br *BurnRecord) String() string          { return fmt.Sprintf("BurnRecord{height=%d, amount=%s}", br.Height, br.Amount) }
func (br *BurnRecord) XXX_MessageName() string { return "syreen.feemarket.BurnRecord" }

// DefaultFeeState returns the default fee state for genesis
func DefaultFeeState() FeeState {
	return FeeState{
		BaseFee:         math.LegacyNewDecWithPrec(1, 2), // 0.01 usyreen per gas
		MinBaseFee:      math.LegacyNewDecWithPrec(1, 4), // 0.0001 usyreen per gas
		MaxBaseFee:      math.LegacyNewDec(1000),         // 1000 usyreen per gas
		BlockGasTarget:  50_000_000,                       // 50M gas target (half of 100M max)
		AdjustmentSpeed: math.LegacyNewDecWithPrec(125, 3), // 0.125 = 1/8 (EIP-1559 style)
	}
}
