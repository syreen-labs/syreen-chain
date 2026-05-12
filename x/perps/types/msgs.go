package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgOpenPosition opens a new perpetual position
type MsgOpenPosition struct {
	Sender   string         `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	MarketID uint64         `protobuf:"varint,2,opt,name=market_id,json=marketId,proto3" json:"market_id"`
	Side     Side           `protobuf:"bytes,3,opt,name=side,proto3" json:"side"`
	Margin   math.Int       `protobuf:"bytes,4,opt,name=margin,proto3" json:"margin"`
	Leverage math.LegacyDec `protobuf:"bytes,5,opt,name=leverage,proto3" json:"leverage"`
}

func (m *MsgOpenPosition) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return err
	}
	if m.Margin.IsNil() || !m.Margin.IsPositive() {
		return ErrInvalidAmount
	}
	if m.Leverage.IsNil() || m.Leverage.LT(math.LegacyOneDec()) {
		return ErrInvalidLeverage
	}
	if m.Side != SideLong && m.Side != SideShort {
		return fmt.Errorf("invalid side: %s", m.Side)
	}
	return nil
}
func (m *MsgOpenPosition) ProtoMessage()             {}
func (m *MsgOpenPosition) Reset()                    { *m = MsgOpenPosition{} }
func (m *MsgOpenPosition) String() string            { return fmt.Sprintf("MsgOpenPosition{%s}", m.Sender) }
func (m *MsgOpenPosition) XXX_MessageName() string   { return "syreen.perps.MsgOpenPosition" }

// MsgClosePosition closes an existing position
type MsgClosePosition struct {
	Sender   string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	MarketID uint64 `protobuf:"varint,2,opt,name=market_id,json=marketId,proto3" json:"market_id"`
}

func (m *MsgClosePosition) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	return err
}
func (m *MsgClosePosition) ProtoMessage()             {}
func (m *MsgClosePosition) Reset()                    { *m = MsgClosePosition{} }
func (m *MsgClosePosition) String() string            { return fmt.Sprintf("MsgClosePosition{%s}", m.Sender) }
func (m *MsgClosePosition) XXX_MessageName() string   { return "syreen.perps.MsgClosePosition" }

// MsgAddMargin adds collateral to an existing position
type MsgAddMargin struct {
	Sender   string   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	MarketID uint64   `protobuf:"varint,2,opt,name=market_id,json=marketId,proto3" json:"market_id"`
	Amount   math.Int `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount"`
}

func (m *MsgAddMargin) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return err
	}
	if m.Amount.IsNil() || !m.Amount.IsPositive() {
		return ErrInvalidAmount
	}
	return nil
}
func (m *MsgAddMargin) ProtoMessage()             {}
func (m *MsgAddMargin) Reset()                    { *m = MsgAddMargin{} }
func (m *MsgAddMargin) String() string            { return fmt.Sprintf("MsgAddMargin{%s}", m.Sender) }
func (m *MsgAddMargin) XXX_MessageName() string   { return "syreen.perps.MsgAddMargin" }

// MsgRemoveMargin removes excess collateral
type MsgRemoveMargin struct {
	Sender   string   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	MarketID uint64   `protobuf:"varint,2,opt,name=market_id,json=marketId,proto3" json:"market_id"`
	Amount   math.Int `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount"`
}

func (m *MsgRemoveMargin) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return err
	}
	if m.Amount.IsNil() || !m.Amount.IsPositive() {
		return ErrInvalidAmount
	}
	return nil
}
func (m *MsgRemoveMargin) ProtoMessage()             {}
func (m *MsgRemoveMargin) Reset()                    { *m = MsgRemoveMargin{} }
func (m *MsgRemoveMargin) String() string            { return fmt.Sprintf("MsgRemoveMargin{%s}", m.Sender) }
func (m *MsgRemoveMargin) XXX_MessageName() string   { return "syreen.perps.MsgRemoveMargin" }

// MsgCreateMarket creates a new perps market (governance only)
type MsgCreateMarket struct {
	Authority   string         `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority"`
	BaseDenom   string         `protobuf:"bytes,2,opt,name=base_denom,json=baseDenom,proto3" json:"base_denom"`
	QuoteDenom  string         `protobuf:"bytes,3,opt,name=quote_denom,json=quoteDenom,proto3" json:"quote_denom"`
	PoolID      uint64         `protobuf:"varint,4,opt,name=pool_id,json=poolId,proto3" json:"pool_id"`
	MaxLeverage math.LegacyDec `protobuf:"bytes,5,opt,name=max_leverage,json=maxLeverage,proto3" json:"max_leverage"`
}

func (m *MsgCreateMarket) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Authority)
	return err
}
func (m *MsgCreateMarket) ProtoMessage()             {}
func (m *MsgCreateMarket) Reset()                    { *m = MsgCreateMarket{} }
func (m *MsgCreateMarket) String() string            { return fmt.Sprintf("MsgCreateMarket{%s}", m.BaseDenom) }
func (m *MsgCreateMarket) XXX_MessageName() string   { return "syreen.perps.MsgCreateMarket" }

// Response types
type MsgOpenPositionResponse struct {
	PositionSize math.LegacyDec `protobuf:"bytes,1,opt,name=position_size,json=positionSize,proto3" json:"position_size"`
	EntryPrice   math.LegacyDec `protobuf:"bytes,2,opt,name=entry_price,json=entryPrice,proto3" json:"entry_price"`
	Fee          math.Int       `protobuf:"bytes,3,opt,name=fee,proto3" json:"fee"`
}

func (m *MsgOpenPositionResponse) ProtoMessage()           {}
func (m *MsgOpenPositionResponse) Reset()                  { *m = MsgOpenPositionResponse{} }
func (m *MsgOpenPositionResponse) String() string          { return "MsgOpenPositionResponse" }
func (m *MsgOpenPositionResponse) XXX_MessageName() string { return "syreen.perps.MsgOpenPositionResponse" }

type MsgClosePositionResponse struct {
	RealizedPnL math.LegacyDec `protobuf:"bytes,1,opt,name=realized_pnl,json=realizedPnl,proto3" json:"realized_pnl"`
	Payout      math.Int       `protobuf:"bytes,2,opt,name=payout,proto3" json:"payout"`
}

func (m *MsgClosePositionResponse) ProtoMessage()           {}
func (m *MsgClosePositionResponse) Reset()                  { *m = MsgClosePositionResponse{} }
func (m *MsgClosePositionResponse) String() string          { return "MsgClosePositionResponse" }
func (m *MsgClosePositionResponse) XXX_MessageName() string { return "syreen.perps.MsgClosePositionResponse" }

type MsgAddMarginResponse struct{}

func (m *MsgAddMarginResponse) ProtoMessage()           {}
func (m *MsgAddMarginResponse) Reset()                  { *m = MsgAddMarginResponse{} }
func (m *MsgAddMarginResponse) String() string          { return "MsgAddMarginResponse" }
func (m *MsgAddMarginResponse) XXX_MessageName() string { return "syreen.perps.MsgAddMarginResponse" }

type MsgRemoveMarginResponse struct{}

func (m *MsgRemoveMarginResponse) ProtoMessage()           {}
func (m *MsgRemoveMarginResponse) Reset()                  { *m = MsgRemoveMarginResponse{} }
func (m *MsgRemoveMarginResponse) String() string          { return "MsgRemoveMarginResponse" }
func (m *MsgRemoveMarginResponse) XXX_MessageName() string { return "syreen.perps.MsgRemoveMarginResponse" }

type MsgCreateMarketResponse struct {
	MarketID uint64 `protobuf:"varint,1,opt,name=market_id,json=marketId,proto3" json:"market_id"`
}

func (m *MsgCreateMarketResponse) ProtoMessage()           {}
func (m *MsgCreateMarketResponse) Reset()                  { *m = MsgCreateMarketResponse{} }
func (m *MsgCreateMarketResponse) String() string          { return "MsgCreateMarketResponse" }
func (m *MsgCreateMarketResponse) XXX_MessageName() string { return "syreen.perps.MsgCreateMarketResponse" }
