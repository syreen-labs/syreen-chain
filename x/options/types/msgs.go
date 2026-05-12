package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgWriteOption — writer locks collateral and creates an option for sale
type MsgWriteOption struct {
	Writer        string         `protobuf:"bytes,1,opt,name=writer,proto3" json:"writer"`
	PoolID        uint64         `protobuf:"varint,2,opt,name=pool_id,json=poolId,proto3" json:"pool_id"`
	OptionType    OptionType     `protobuf:"bytes,3,opt,name=option_type,json=optionType,proto3" json:"option_type"`
	StrikePrice   math.LegacyDec `protobuf:"bytes,4,opt,name=strike_price,json=strikePrice,proto3" json:"strike_price"`
	Amount        math.Int       `protobuf:"bytes,5,opt,name=amount,proto3" json:"amount"`
	ExpiryBlock   int64          `protobuf:"varint,6,opt,name=expiry_block,json=expiryBlock,proto3" json:"expiry_block"`
	CustomPremium math.Int       `protobuf:"bytes,7,opt,name=custom_premium,json=customPremium,proto3" json:"custom_premium"`
}

func (m *MsgWriteOption) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Writer); err != nil {
		return err
	}
	if m.Amount.IsNil() || !m.Amount.IsPositive() {
		return ErrInvalidAmount
	}
	if m.StrikePrice.IsNil() || !m.StrikePrice.IsPositive() {
		return ErrInvalidStrike
	}
	if m.ExpiryBlock <= 0 {
		return ErrInvalidExpiry
	}
	if m.OptionType != OptionTypeCall && m.OptionType != OptionTypePut {
		return fmt.Errorf("invalid option type: %s", m.OptionType)
	}
	return nil
}
func (m *MsgWriteOption) ProtoMessage()           {}
func (m *MsgWriteOption) Reset()                  { *m = MsgWriteOption{} }
func (m *MsgWriteOption) String() string          { return fmt.Sprintf("MsgWriteOption{%s}", m.Writer) }
func (m *MsgWriteOption) XXX_MessageName() string { return "syreen.options.MsgWriteOption" }

type MsgWriteOptionResponse struct {
	OptionID uint64   `protobuf:"varint,1,opt,name=option_id,json=optionId,proto3" json:"option_id"`
	Premium  math.Int `protobuf:"bytes,2,opt,name=premium,proto3" json:"premium"`
}

func (m *MsgWriteOptionResponse) ProtoMessage()           {}
func (m *MsgWriteOptionResponse) Reset()                  { *m = MsgWriteOptionResponse{} }
func (m *MsgWriteOptionResponse) String() string          { return "MsgWriteOptionResponse" }
func (m *MsgWriteOptionResponse) XXX_MessageName() string { return "syreen.options.MsgWriteOptionResponse" }

// MsgBuyOption — buyer pays the premium to acquire the option
type MsgBuyOption struct {
	Buyer    string `protobuf:"bytes,1,opt,name=buyer,proto3" json:"buyer"`
	OptionID uint64 `protobuf:"varint,2,opt,name=option_id,json=optionId,proto3" json:"option_id"`
}

func (m *MsgBuyOption) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Buyer)
	return err
}
func (m *MsgBuyOption) ProtoMessage()           {}
func (m *MsgBuyOption) Reset()                  { *m = MsgBuyOption{} }
func (m *MsgBuyOption) String() string          { return fmt.Sprintf("MsgBuyOption{%s,%d}", m.Buyer, m.OptionID) }
func (m *MsgBuyOption) XXX_MessageName() string { return "syreen.options.MsgBuyOption" }

type MsgBuyOptionResponse struct{}

func (m *MsgBuyOptionResponse) ProtoMessage()           {}
func (m *MsgBuyOptionResponse) Reset()                  { *m = MsgBuyOptionResponse{} }
func (m *MsgBuyOptionResponse) String() string          { return "MsgBuyOptionResponse" }
func (m *MsgBuyOptionResponse) XXX_MessageName() string { return "syreen.options.MsgBuyOptionResponse" }

// MsgExerciseOption — buyer exercises an active option
type MsgExerciseOption struct {
	Buyer    string `protobuf:"bytes,1,opt,name=buyer,proto3" json:"buyer"`
	OptionID uint64 `protobuf:"varint,2,opt,name=option_id,json=optionId,proto3" json:"option_id"`
}

func (m *MsgExerciseOption) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Buyer)
	return err
}
func (m *MsgExerciseOption) ProtoMessage()           {}
func (m *MsgExerciseOption) Reset()                  { *m = MsgExerciseOption{} }
func (m *MsgExerciseOption) String() string          { return fmt.Sprintf("MsgExerciseOption{%s,%d}", m.Buyer, m.OptionID) }
func (m *MsgExerciseOption) XXX_MessageName() string { return "syreen.options.MsgExerciseOption" }

type MsgExerciseOptionResponse struct {
	Payout math.Int `protobuf:"bytes,1,opt,name=payout,proto3" json:"payout"`
}

func (m *MsgExerciseOptionResponse) ProtoMessage()           {}
func (m *MsgExerciseOptionResponse) Reset()                  { *m = MsgExerciseOptionResponse{} }
func (m *MsgExerciseOptionResponse) String() string          { return "MsgExerciseOptionResponse" }
func (m *MsgExerciseOptionResponse) XXX_MessageName() string { return "syreen.options.MsgExerciseOptionResponse" }

// MsgCancelOption — writer cancels an unsold option and reclaims collateral
type MsgCancelOption struct {
	Writer   string `protobuf:"bytes,1,opt,name=writer,proto3" json:"writer"`
	OptionID uint64 `protobuf:"varint,2,opt,name=option_id,json=optionId,proto3" json:"option_id"`
}

func (m *MsgCancelOption) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Writer)
	return err
}
func (m *MsgCancelOption) ProtoMessage()           {}
func (m *MsgCancelOption) Reset()                  { *m = MsgCancelOption{} }
func (m *MsgCancelOption) String() string          { return fmt.Sprintf("MsgCancelOption{%s,%d}", m.Writer, m.OptionID) }
func (m *MsgCancelOption) XXX_MessageName() string { return "syreen.options.MsgCancelOption" }

type MsgCancelOptionResponse struct{}

func (m *MsgCancelOptionResponse) ProtoMessage()           {}
func (m *MsgCancelOptionResponse) Reset()                  { *m = MsgCancelOptionResponse{} }
func (m *MsgCancelOptionResponse) String() string          { return "MsgCancelOptionResponse" }
func (m *MsgCancelOptionResponse) XXX_MessageName() string { return "syreen.options.MsgCancelOptionResponse" }
