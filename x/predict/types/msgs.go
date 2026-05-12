package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgCreateMarket creates a new binary outcome prediction market
type MsgCreateMarket struct {
	Creator          string   `json:"creator"`
	Question         string   `json:"question"`
	Resolver         string   `json:"resolver"`
	QuoteDenom       string   `json:"quote_denom"`
	ResolutionBlock  int64    `json:"resolution_block"`
	InitialLiquidity math.Int `json:"initial_liquidity"`
}

func (m *MsgCreateMarket) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil { return err }
	if _, err := sdk.AccAddressFromBech32(m.Resolver); err != nil { return ErrInvalidResolver }
	if m.Question == "" { return ErrInvalidQuestion }
	if m.QuoteDenom == "" { return ErrInvalidDenom }
	if m.InitialLiquidity.IsNil() || !m.InitialLiquidity.IsPositive() { return ErrInvalidLiquidity }
	if m.ResolutionBlock <= 0 { return fmt.Errorf("resolution_block must be positive") }
	return nil
}
func (m *MsgCreateMarket) ProtoMessage()                {}
func (m *MsgCreateMarket) Reset()                       { *m = MsgCreateMarket{} }
func (m *MsgCreateMarket) String() string               { return fmt.Sprintf("MsgCreateMarket{%s}", m.Question) }
func (m *MsgCreateMarket) XXX_MessageName() string      { return "syreen.predict.MsgCreateMarket" }

// MsgBuyShares buys YES or NO shares in a market
type MsgBuyShares struct {
	Sender   string   `json:"sender"`
	MarketID uint64   `json:"market_id"`
	Outcome  string   `json:"outcome"` // "yes" or "no"
	Amount   math.Int `json:"amount"`  // quote tokens to spend
}

func (m *MsgBuyShares) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil { return err }
	if m.Outcome != "yes" && m.Outcome != "no" { return ErrInvalidOutcome }
	if m.Amount.IsNil() || !m.Amount.IsPositive() { return ErrInvalidAmount }
	return nil
}
func (m *MsgBuyShares) ProtoMessage()                {}
func (m *MsgBuyShares) Reset()                       { *m = MsgBuyShares{} }
func (m *MsgBuyShares) String() string               { return fmt.Sprintf("MsgBuyShares{%s,%d,%s}", m.Sender, m.MarketID, m.Outcome) }
func (m *MsgBuyShares) XXX_MessageName() string      { return "syreen.predict.MsgBuyShares" }

// MsgSellShares sells YES or NO shares back to the market
type MsgSellShares struct {
	Sender   string   `json:"sender"`
	MarketID uint64   `json:"market_id"`
	Outcome  string   `json:"outcome"` // "yes" or "no"
	Shares   math.Int `json:"shares"`  // number of shares to sell
}

func (m *MsgSellShares) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil { return err }
	if m.Outcome != "yes" && m.Outcome != "no" { return ErrInvalidOutcome }
	if m.Shares.IsNil() || !m.Shares.IsPositive() { return ErrInvalidAmount }
	return nil
}
func (m *MsgSellShares) ProtoMessage()                {}
func (m *MsgSellShares) Reset()                       { *m = MsgSellShares{} }
func (m *MsgSellShares) String() string               { return fmt.Sprintf("MsgSellShares{%s,%d,%s}", m.Sender, m.MarketID, m.Outcome) }
func (m *MsgSellShares) XXX_MessageName() string      { return "syreen.predict.MsgSellShares" }

// MsgResolveMarket resolves a market to YES, NO, or VOID
type MsgResolveMarket struct {
	Resolver string `json:"resolver"`
	MarketID uint64 `json:"market_id"`
	Outcome  string `json:"outcome"` // "yes", "no", or "void"
}

func (m *MsgResolveMarket) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Resolver); err != nil { return err }
	if m.Outcome != "yes" && m.Outcome != "no" && m.Outcome != "void" { return ErrInvalidOutcome }
	return nil
}
func (m *MsgResolveMarket) ProtoMessage()                {}
func (m *MsgResolveMarket) Reset()                       { *m = MsgResolveMarket{} }
func (m *MsgResolveMarket) String() string               { return fmt.Sprintf("MsgResolveMarket{%d,%s}", m.MarketID, m.Outcome) }
func (m *MsgResolveMarket) XXX_MessageName() string      { return "syreen.predict.MsgResolveMarket" }

// MsgClaimWinnings claims payout from a resolved market
type MsgClaimWinnings struct {
	Sender   string `json:"sender"`
	MarketID uint64 `json:"market_id"`
}

func (m *MsgClaimWinnings) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	return err
}
func (m *MsgClaimWinnings) ProtoMessage()                {}
func (m *MsgClaimWinnings) Reset()                       { *m = MsgClaimWinnings{} }
func (m *MsgClaimWinnings) String() string               { return fmt.Sprintf("MsgClaimWinnings{%s,%d}", m.Sender, m.MarketID) }
func (m *MsgClaimWinnings) XXX_MessageName() string      { return "syreen.predict.MsgClaimWinnings" }

// Response types

type MsgCreateMarketResponse struct {
	MarketID uint64 `json:"market_id"`
}
func (m *MsgCreateMarketResponse) ProtoMessage()           {}
func (m *MsgCreateMarketResponse) Reset()                  { *m = MsgCreateMarketResponse{} }
func (m *MsgCreateMarketResponse) String() string          { return "MsgCreateMarketResponse" }
func (m *MsgCreateMarketResponse) XXX_MessageName() string { return "syreen.predict.MsgCreateMarketResponse" }

type MsgBuySharesResponse struct {
	SharesBought math.Int `json:"shares_bought"`
	AvgPrice     math.Int `json:"avg_price"` // in basis points (e.g. 5000 = 0.50)
}
func (m *MsgBuySharesResponse) ProtoMessage()           {}
func (m *MsgBuySharesResponse) Reset()                  { *m = MsgBuySharesResponse{} }
func (m *MsgBuySharesResponse) String() string          { return "MsgBuySharesResponse" }
func (m *MsgBuySharesResponse) XXX_MessageName() string { return "syreen.predict.MsgBuySharesResponse" }

type MsgSellSharesResponse struct {
	QuoteReturned math.Int `json:"quote_returned"`
}
func (m *MsgSellSharesResponse) ProtoMessage()           {}
func (m *MsgSellSharesResponse) Reset()                  { *m = MsgSellSharesResponse{} }
func (m *MsgSellSharesResponse) String() string          { return "MsgSellSharesResponse" }
func (m *MsgSellSharesResponse) XXX_MessageName() string { return "syreen.predict.MsgSellSharesResponse" }

type MsgResolveMarketResponse struct{}
func (m *MsgResolveMarketResponse) ProtoMessage()           {}
func (m *MsgResolveMarketResponse) Reset()                  { *m = MsgResolveMarketResponse{} }
func (m *MsgResolveMarketResponse) String() string          { return "MsgResolveMarketResponse" }
func (m *MsgResolveMarketResponse) XXX_MessageName() string { return "syreen.predict.MsgResolveMarketResponse" }

type MsgClaimWinningsResponse struct {
	Amount math.Int `json:"amount"`
}
func (m *MsgClaimWinningsResponse) ProtoMessage()           {}
func (m *MsgClaimWinningsResponse) Reset()                  { *m = MsgClaimWinningsResponse{} }
func (m *MsgClaimWinningsResponse) String() string          { return "MsgClaimWinningsResponse" }
func (m *MsgClaimWinningsResponse) XXX_MessageName() string { return "syreen.predict.MsgClaimWinningsResponse" }
