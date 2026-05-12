package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgOpenPosition opens a new perpetual position
type MsgOpenPosition struct {
	Sender   string         `json:"sender"`
	MarketID uint64         `json:"market_id"`
	Side     Side           `json:"side"`      // long or short
	Margin   math.Int       `json:"margin"`    // collateral amount in quote denom
	Leverage math.LegacyDec `json:"leverage"`  // desired leverage (2-20x)
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
	Sender   string `json:"sender"`
	MarketID uint64 `json:"market_id"`
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
	Sender   string   `json:"sender"`
	MarketID uint64   `json:"market_id"`
	Amount   math.Int `json:"amount"`
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
	Sender   string   `json:"sender"`
	MarketID uint64   `json:"market_id"`
	Amount   math.Int `json:"amount"`
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
	Authority  string         `json:"authority"`
	BaseDenom  string         `json:"base_denom"`
	QuoteDenom string         `json:"quote_denom"`
	PoolID     uint64         `json:"pool_id"`
	MaxLeverage math.LegacyDec `json:"max_leverage"`
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
	PositionSize math.LegacyDec `json:"position_size"`
	EntryPrice   math.LegacyDec `json:"entry_price"`
	Fee          math.Int       `json:"fee"`
}
func (m *MsgOpenPositionResponse) ProtoMessage()           {}
func (m *MsgOpenPositionResponse) Reset()                  { *m = MsgOpenPositionResponse{} }
func (m *MsgOpenPositionResponse) String() string          { return "MsgOpenPositionResponse" }
func (m *MsgOpenPositionResponse) XXX_MessageName() string { return "syreen.perps.MsgOpenPositionResponse" }

type MsgClosePositionResponse struct {
	RealizedPnL math.LegacyDec `json:"realized_pnl"`
	Payout      math.Int       `json:"payout"`
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
	MarketID uint64 `json:"market_id"`
}
func (m *MsgCreateMarketResponse) ProtoMessage()           {}
func (m *MsgCreateMarketResponse) Reset()                  { *m = MsgCreateMarketResponse{} }
func (m *MsgCreateMarketResponse) String() string          { return "MsgCreateMarketResponse" }
func (m *MsgCreateMarketResponse) XXX_MessageName() string { return "syreen.perps.MsgCreateMarketResponse" }
