package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgRecordTrade records a trade event for portfolio tracking
type MsgRecordTrade struct {
	Sender    string   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	TradeType string   `protobuf:"bytes,2,opt,name=trade_type,json=tradeType,proto3" json:"trade_type"`
	Denom     string   `protobuf:"bytes,3,opt,name=denom,proto3" json:"denom"`
	Amount    math.Int `protobuf:"bytes,4,opt,name=amount,proto3" json:"amount"`
	Price     math.Int `protobuf:"bytes,5,opt,name=price,proto3" json:"price"`
	PnL       math.Int `protobuf:"bytes,6,opt,name=pnl,proto3" json:"pnl"`
}

func (m *MsgRecordTrade) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil { return ErrInvalidAddress }
	if !IsValidTradeType(m.TradeType) { return ErrInvalidTradeType }
	if m.Amount.IsNil() || m.Amount.IsNegative() { return ErrInvalidAmount }
	return nil
}
func (m *MsgRecordTrade) ProtoMessage() {}
func (m *MsgRecordTrade) Reset() { *m = MsgRecordTrade{} }
func (m *MsgRecordTrade) String() string { return fmt.Sprintf("MsgRecordTrade{%s,%s,%s}", m.Sender, m.TradeType, m.Denom) }
func (m *MsgRecordTrade) XXX_MessageName() string { return "syreen.portfolio.MsgRecordTrade" }

// MsgCreateCompetition creates a trading competition
type MsgCreateCompetition struct {
	Creator         string   `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator"`
	Name            string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name"`
	StartBlock      int64    `protobuf:"varint,3,opt,name=start_block,json=startBlock,proto3" json:"start_block"`
	EndBlock        int64    `protobuf:"varint,4,opt,name=end_block,json=endBlock,proto3" json:"end_block"`
	PrizeDenom      string   `protobuf:"bytes,5,opt,name=prize_denom,json=prizeDenom,proto3" json:"prize_denom"`
	PrizePool       math.Int `protobuf:"bytes,6,opt,name=prize_pool,json=prizePool,proto3" json:"prize_pool"`
	EntryFee        math.Int `protobuf:"bytes,7,opt,name=entry_fee,json=entryFee,proto3" json:"entry_fee"`
	MaxParticipants uint64   `protobuf:"varint,8,opt,name=max_participants,json=maxParticipants,proto3" json:"max_participants"`
}

func (m *MsgCreateCompetition) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil { return ErrInvalidAddress }
	if m.Name == "" { return ErrInvalidCompetition }
	if m.StartBlock >= m.EndBlock { return ErrInvalidCompetition }
	if m.PrizeDenom == "" { return ErrInvalidCompetition }
	if m.PrizePool.IsNil() || m.PrizePool.IsNegative() { return ErrInvalidAmount }
	if m.MaxParticipants == 0 { return ErrInvalidCompetition }
	return nil
}
func (m *MsgCreateCompetition) ProtoMessage() {}
func (m *MsgCreateCompetition) Reset() { *m = MsgCreateCompetition{} }
func (m *MsgCreateCompetition) String() string { return fmt.Sprintf("MsgCreateCompetition{%s,%s}", m.Creator, m.Name) }
func (m *MsgCreateCompetition) XXX_MessageName() string { return "syreen.portfolio.MsgCreateCompetition" }

// MsgJoinCompetition joins a trading competition
type MsgJoinCompetition struct {
	Sender        string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	CompetitionID uint64 `protobuf:"varint,2,opt,name=competition_id,json=competitionId,proto3" json:"competition_id"`
}

func (m *MsgJoinCompetition) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil { return ErrInvalidAddress }
	return nil
}
func (m *MsgJoinCompetition) ProtoMessage() {}
func (m *MsgJoinCompetition) Reset() { *m = MsgJoinCompetition{} }
func (m *MsgJoinCompetition) String() string { return fmt.Sprintf("MsgJoinCompetition{%s,%d}", m.Sender, m.CompetitionID) }
func (m *MsgJoinCompetition) XXX_MessageName() string { return "syreen.portfolio.MsgJoinCompetition" }

// MsgEndCompetition ends a competition and distributes prizes
type MsgEndCompetition struct {
	Authority     string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority"`
	CompetitionID uint64 `protobuf:"varint,2,opt,name=competition_id,json=competitionId,proto3" json:"competition_id"`
}

func (m *MsgEndCompetition) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil { return ErrInvalidAddress }
	return nil
}
func (m *MsgEndCompetition) ProtoMessage() {}
func (m *MsgEndCompetition) Reset() { *m = MsgEndCompetition{} }
func (m *MsgEndCompetition) String() string { return fmt.Sprintf("MsgEndCompetition{%s,%d}", m.Authority, m.CompetitionID) }
func (m *MsgEndCompetition) XXX_MessageName() string { return "syreen.portfolio.MsgEndCompetition" }

// MsgUpdatePortfolio triggers a portfolio recalculation
type MsgUpdatePortfolio struct {
	Sender string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
}

func (m *MsgUpdatePortfolio) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil { return ErrInvalidAddress }
	return nil
}
func (m *MsgUpdatePortfolio) ProtoMessage() {}
func (m *MsgUpdatePortfolio) Reset() { *m = MsgUpdatePortfolio{} }
func (m *MsgUpdatePortfolio) String() string { return fmt.Sprintf("MsgUpdatePortfolio{%s}", m.Sender) }
func (m *MsgUpdatePortfolio) XXX_MessageName() string { return "syreen.portfolio.MsgUpdatePortfolio" }

// Response types
type MsgRecordTradeResponse struct {
	TradeID uint64 `protobuf:"varint,1,opt,name=trade_id,json=tradeId,proto3" json:"trade_id"`
}
func (m *MsgRecordTradeResponse) ProtoMessage() {}
func (m *MsgRecordTradeResponse) Reset() { *m = MsgRecordTradeResponse{} }
func (m *MsgRecordTradeResponse) String() string { return "MsgRecordTradeResponse" }
func (m *MsgRecordTradeResponse) XXX_MessageName() string { return "syreen.portfolio.MsgRecordTradeResponse" }

type MsgCreateCompetitionResponse struct {
	CompetitionID uint64 `protobuf:"varint,1,opt,name=competition_id,json=competitionId,proto3" json:"competition_id"`
}
func (m *MsgCreateCompetitionResponse) ProtoMessage() {}
func (m *MsgCreateCompetitionResponse) Reset() { *m = MsgCreateCompetitionResponse{} }
func (m *MsgCreateCompetitionResponse) String() string { return "MsgCreateCompetitionResponse" }
func (m *MsgCreateCompetitionResponse) XXX_MessageName() string { return "syreen.portfolio.MsgCreateCompetitionResponse" }

type MsgJoinCompetitionResponse struct{}
func (m *MsgJoinCompetitionResponse) ProtoMessage() {}
func (m *MsgJoinCompetitionResponse) Reset() { *m = MsgJoinCompetitionResponse{} }
func (m *MsgJoinCompetitionResponse) String() string { return "MsgJoinCompetitionResponse" }
func (m *MsgJoinCompetitionResponse) XXX_MessageName() string { return "syreen.portfolio.MsgJoinCompetitionResponse" }

type MsgEndCompetitionResponse struct {
	Winners []string `protobuf:"bytes,1,rep,name=winners,proto3" json:"winners"`
}
func (m *MsgEndCompetitionResponse) ProtoMessage() {}
func (m *MsgEndCompetitionResponse) Reset() { *m = MsgEndCompetitionResponse{} }
func (m *MsgEndCompetitionResponse) String() string { return "MsgEndCompetitionResponse" }
func (m *MsgEndCompetitionResponse) XXX_MessageName() string { return "syreen.portfolio.MsgEndCompetitionResponse" }

type MsgUpdatePortfolioResponse struct {
	TotalValue math.Int `protobuf:"bytes,1,opt,name=total_value,json=totalValue,proto3" json:"total_value"`
}
func (m *MsgUpdatePortfolioResponse) ProtoMessage() {}
func (m *MsgUpdatePortfolioResponse) Reset() { *m = MsgUpdatePortfolioResponse{} }
func (m *MsgUpdatePortfolioResponse) String() string { return "MsgUpdatePortfolioResponse" }
func (m *MsgUpdatePortfolioResponse) XXX_MessageName() string { return "syreen.portfolio.MsgUpdatePortfolioResponse" }
