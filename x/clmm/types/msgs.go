package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgCreateCLPool creates a new concentrated liquidity pool
type MsgCreateCLPool struct {
	Sender       string         `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	DenomA       string         `protobuf:"bytes,2,opt,name=denom_a,json=denomA,proto3" json:"denom_a"`
	DenomB       string         `protobuf:"bytes,3,opt,name=denom_b,json=denomB,proto3" json:"denom_b"`
	TickSpacing  int64          `protobuf:"varint,4,opt,name=tick_spacing,json=tickSpacing,proto3" json:"tick_spacing"`
	FeeRate      math.LegacyDec `protobuf:"bytes,5,opt,name=fee_rate,json=feeRate,proto3" json:"fee_rate"`
	InitialPrice math.LegacyDec `protobuf:"bytes,6,opt,name=initial_price,json=initialPrice,proto3" json:"initial_price"`
}

func (m *MsgCreateCLPool) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return err
	}
	if m.DenomA == m.DenomB {
		return ErrSameDenom
	}
	if m.TickSpacing <= 0 {
		return ErrInvalidTickSpacing
	}
	if m.FeeRate.IsNil() || m.FeeRate.IsNegative() || m.FeeRate.GT(math.LegacyNewDecWithPrec(1, 1)) {
		return ErrInvalidFeeRate
	} // cap at 10%
	if m.InitialPrice.IsNil() || !m.InitialPrice.IsPositive() {
		return ErrInvalidPrice
	}
	return nil
}
func (m *MsgCreateCLPool) ProtoMessage()             {}
func (m *MsgCreateCLPool) Reset()                    { *m = MsgCreateCLPool{} }
func (m *MsgCreateCLPool) String() string            { return fmt.Sprintf("MsgCreateCLPool{%s/%s}", m.DenomA, m.DenomB) }
func (m *MsgCreateCLPool) XXX_MessageName() string   { return "syreen.clmm.MsgCreateCLPool" }

// MsgCreatePosition creates a new LP position with a price range
type MsgCreatePosition struct {
	Sender         string   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	PoolID         uint64   `protobuf:"varint,2,opt,name=pool_id,json=poolId,proto3" json:"pool_id"`
	TickLower      int64    `protobuf:"varint,3,opt,name=tick_lower,json=tickLower,proto3" json:"tick_lower"`
	TickUpper      int64    `protobuf:"varint,4,opt,name=tick_upper,json=tickUpper,proto3" json:"tick_upper"`
	Amount0Desired math.Int `protobuf:"bytes,5,opt,name=amount_0_desired,json=amount0Desired,proto3" json:"amount_0_desired"`
	Amount1Desired math.Int `protobuf:"bytes,6,opt,name=amount_1_desired,json=amount1Desired,proto3" json:"amount_1_desired"`
	Amount0Min     math.Int `protobuf:"bytes,7,opt,name=amount_0_min,json=amount0Min,proto3" json:"amount_0_min"`
	Amount1Min     math.Int `protobuf:"bytes,8,opt,name=amount_1_min,json=amount1Min,proto3" json:"amount_1_min"`
}

func (m *MsgCreatePosition) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return err
	}
	if m.TickLower >= m.TickUpper {
		return ErrInvalidTickRange
	}
	if m.Amount0Desired.IsNil() || m.Amount0Desired.IsNegative() {
		return ErrInvalidAmount
	}
	if m.Amount1Desired.IsNil() || m.Amount1Desired.IsNegative() {
		return ErrInvalidAmount
	}
	return nil
}
func (m *MsgCreatePosition) ProtoMessage()             {}
func (m *MsgCreatePosition) Reset()                    { *m = MsgCreatePosition{} }
func (m *MsgCreatePosition) String() string            { return fmt.Sprintf("MsgCreatePosition{pool=%d}", m.PoolID) }
func (m *MsgCreatePosition) XXX_MessageName() string   { return "syreen.clmm.MsgCreatePosition" }

// MsgAddLiquidity adds more liquidity to an existing position
type MsgAddLiquidity struct {
	Sender         string   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	PositionID     uint64   `protobuf:"varint,2,opt,name=position_id,json=positionId,proto3" json:"position_id"`
	Amount0Desired math.Int `protobuf:"bytes,3,opt,name=amount_0_desired,json=amount0Desired,proto3" json:"amount_0_desired"`
	Amount1Desired math.Int `protobuf:"bytes,4,opt,name=amount_1_desired,json=amount1Desired,proto3" json:"amount_1_desired"`
}

func (m *MsgAddLiquidity) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return err
	}
	if m.Amount0Desired.IsNil() || m.Amount0Desired.IsNegative() {
		return ErrInvalidAmount
	}
	if m.Amount1Desired.IsNil() || m.Amount1Desired.IsNegative() {
		return ErrInvalidAmount
	}
	return nil
}
func (m *MsgAddLiquidity) ProtoMessage()             {}
func (m *MsgAddLiquidity) Reset()                    { *m = MsgAddLiquidity{} }
func (m *MsgAddLiquidity) String() string            { return fmt.Sprintf("MsgAddLiquidity{pos=%d}", m.PositionID) }
func (m *MsgAddLiquidity) XXX_MessageName() string   { return "syreen.clmm.MsgAddLiquidity" }

// MsgRemoveLiquidity removes liquidity from a position
type MsgRemoveLiquidity struct {
	Sender          string         `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	PositionID      uint64         `protobuf:"varint,2,opt,name=position_id,json=positionId,proto3" json:"position_id"`
	LiquidityAmount math.LegacyDec `protobuf:"bytes,3,opt,name=liquidity_amount,json=liquidityAmount,proto3" json:"liquidity_amount"`
}

func (m *MsgRemoveLiquidity) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return err
	}
	if m.LiquidityAmount.IsNil() || !m.LiquidityAmount.IsPositive() {
		return ErrInvalidAmount
	}
	return nil
}
func (m *MsgRemoveLiquidity) ProtoMessage()             {}
func (m *MsgRemoveLiquidity) Reset()                    { *m = MsgRemoveLiquidity{} }
func (m *MsgRemoveLiquidity) String() string            { return fmt.Sprintf("MsgRemoveLiquidity{pos=%d}", m.PositionID) }
func (m *MsgRemoveLiquidity) XXX_MessageName() string   { return "syreen.clmm.MsgRemoveLiquidity" }

// MsgCollectFees collects accrued fees from a position
type MsgCollectFees struct {
	Sender     string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	PositionID uint64 `protobuf:"varint,2,opt,name=position_id,json=positionId,proto3" json:"position_id"`
}

func (m *MsgCollectFees) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	return err
}
func (m *MsgCollectFees) ProtoMessage()             {}
func (m *MsgCollectFees) Reset()                    { *m = MsgCollectFees{} }
func (m *MsgCollectFees) String() string            { return fmt.Sprintf("MsgCollectFees{pos=%d}", m.PositionID) }
func (m *MsgCollectFees) XXX_MessageName() string   { return "syreen.clmm.MsgCollectFees" }

// MsgCLSwap swaps tokens through a CL pool
type MsgCLSwap struct {
	Sender       string   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	PoolID       uint64   `protobuf:"varint,2,opt,name=pool_id,json=poolId,proto3" json:"pool_id"`
	DenomIn      string   `protobuf:"bytes,3,opt,name=denom_in,json=denomIn,proto3" json:"denom_in"`
	AmountIn     math.Int `protobuf:"bytes,4,opt,name=amount_in,json=amountIn,proto3" json:"amount_in"`
	MinAmountOut math.Int `protobuf:"bytes,5,opt,name=min_amount_out,json=minAmountOut,proto3" json:"min_amount_out"`
}

func (m *MsgCLSwap) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return err
	}
	if m.AmountIn.IsNil() || !m.AmountIn.IsPositive() {
		return ErrInvalidAmount
	}
	if m.MinAmountOut.IsNil() || m.MinAmountOut.IsNegative() {
		return ErrInvalidAmount
	}
	return nil
}
func (m *MsgCLSwap) ProtoMessage()             {}
func (m *MsgCLSwap) Reset()                    { *m = MsgCLSwap{} }
func (m *MsgCLSwap) String() string            { return fmt.Sprintf("MsgCLSwap{pool=%d,%s}", m.PoolID, m.DenomIn) }
func (m *MsgCLSwap) XXX_MessageName() string   { return "syreen.clmm.MsgCLSwap" }

// Response types
type MsgCreateCLPoolResponse struct {
	PoolID uint64 `protobuf:"varint,1,opt,name=pool_id,json=poolId,proto3" json:"pool_id"`
}

func (m *MsgCreateCLPoolResponse) ProtoMessage()             {}
func (m *MsgCreateCLPoolResponse) Reset()                    { *m = MsgCreateCLPoolResponse{} }
func (m *MsgCreateCLPoolResponse) String() string            { return "MsgCreateCLPoolResponse" }
func (m *MsgCreateCLPoolResponse) XXX_MessageName() string   { return "syreen.clmm.MsgCreateCLPoolResponse" }

type MsgCreatePositionResponse struct {
	PositionID uint64         `protobuf:"varint,1,opt,name=position_id,json=positionId,proto3" json:"position_id"`
	Amount0    math.Int       `protobuf:"bytes,2,opt,name=amount_0,json=amount0,proto3" json:"amount_0"`
	Amount1    math.Int       `protobuf:"bytes,3,opt,name=amount_1,json=amount1,proto3" json:"amount_1"`
	Liquidity  math.LegacyDec `protobuf:"bytes,4,opt,name=liquidity,proto3" json:"liquidity"`
}

func (m *MsgCreatePositionResponse) ProtoMessage()             {}
func (m *MsgCreatePositionResponse) Reset()                    { *m = MsgCreatePositionResponse{} }
func (m *MsgCreatePositionResponse) String() string            { return "MsgCreatePositionResponse" }
func (m *MsgCreatePositionResponse) XXX_MessageName() string   { return "syreen.clmm.MsgCreatePositionResponse" }

type MsgAddLiquidityResponse struct {
	Amount0   math.Int       `protobuf:"bytes,1,opt,name=amount_0,json=amount0,proto3" json:"amount_0"`
	Amount1   math.Int       `protobuf:"bytes,2,opt,name=amount_1,json=amount1,proto3" json:"amount_1"`
	Liquidity math.LegacyDec `protobuf:"bytes,3,opt,name=liquidity,proto3" json:"liquidity"`
}

func (m *MsgAddLiquidityResponse) ProtoMessage()             {}
func (m *MsgAddLiquidityResponse) Reset()                    { *m = MsgAddLiquidityResponse{} }
func (m *MsgAddLiquidityResponse) String() string            { return "MsgAddLiquidityResponse" }
func (m *MsgAddLiquidityResponse) XXX_MessageName() string   { return "syreen.clmm.MsgAddLiquidityResponse" }

type MsgRemoveLiquidityResponse struct {
	Amount0 math.Int `protobuf:"bytes,1,opt,name=amount_0,json=amount0,proto3" json:"amount_0"`
	Amount1 math.Int `protobuf:"bytes,2,opt,name=amount_1,json=amount1,proto3" json:"amount_1"`
}

func (m *MsgRemoveLiquidityResponse) ProtoMessage()             {}
func (m *MsgRemoveLiquidityResponse) Reset()                    { *m = MsgRemoveLiquidityResponse{} }
func (m *MsgRemoveLiquidityResponse) String() string            { return "MsgRemoveLiquidityResponse" }
func (m *MsgRemoveLiquidityResponse) XXX_MessageName() string   { return "syreen.clmm.MsgRemoveLiquidityResponse" }

type MsgCollectFeesResponse struct {
	Amount0 math.Int `protobuf:"bytes,1,opt,name=amount_0,json=amount0,proto3" json:"amount_0"`
	Amount1 math.Int `protobuf:"bytes,2,opt,name=amount_1,json=amount1,proto3" json:"amount_1"`
}

func (m *MsgCollectFeesResponse) ProtoMessage()             {}
func (m *MsgCollectFeesResponse) Reset()                    { *m = MsgCollectFeesResponse{} }
func (m *MsgCollectFeesResponse) String() string            { return "MsgCollectFeesResponse" }
func (m *MsgCollectFeesResponse) XXX_MessageName() string   { return "syreen.clmm.MsgCollectFeesResponse" }

type MsgCLSwapResponse struct {
	AmountOut math.Int `protobuf:"bytes,1,opt,name=amount_out,json=amountOut,proto3" json:"amount_out"`
}

func (m *MsgCLSwapResponse) ProtoMessage()             {}
func (m *MsgCLSwapResponse) Reset()                    { *m = MsgCLSwapResponse{} }
func (m *MsgCLSwapResponse) String() string            { return "MsgCLSwapResponse" }
func (m *MsgCLSwapResponse) XXX_MessageName() string   { return "syreen.clmm.MsgCLSwapResponse" }
