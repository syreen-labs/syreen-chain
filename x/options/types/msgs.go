package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgWriteOption — writer locks collateral and creates an option for sale
type MsgWriteOption struct {
	Writer          string         `json:"writer"`
	PoolID          uint64         `json:"pool_id"`
	OptionType      OptionType     `json:"option_type"`   // call or put
	StrikePrice     math.LegacyDec `json:"strike_price"`  // quote per underlying
	Amount          math.Int       `json:"amount"`        // underlying amount
	ExpiryBlock     int64          `json:"expiry_block"`
	CustomPremium   math.Int       `json:"custom_premium"` // 0 = auto-calculate via Black-Scholes
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
	OptionID uint64   `json:"option_id"`
	Premium  math.Int `json:"premium"`
}

func (m *MsgWriteOptionResponse) ProtoMessage()           {}
func (m *MsgWriteOptionResponse) Reset()                  { *m = MsgWriteOptionResponse{} }
func (m *MsgWriteOptionResponse) String() string          { return "MsgWriteOptionResponse" }
func (m *MsgWriteOptionResponse) XXX_MessageName() string { return "syreen.options.MsgWriteOptionResponse" }

// MsgBuyOption — buyer pays the premium to acquire the option
type MsgBuyOption struct {
	Buyer    string `json:"buyer"`
	OptionID uint64 `json:"option_id"`
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
	Buyer    string `json:"buyer"`
	OptionID uint64 `json:"option_id"`
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
	Payout math.Int `json:"payout"`
}

func (m *MsgExerciseOptionResponse) ProtoMessage()           {}
func (m *MsgExerciseOptionResponse) Reset()                  { *m = MsgExerciseOptionResponse{} }
func (m *MsgExerciseOptionResponse) String() string          { return "MsgExerciseOptionResponse" }
func (m *MsgExerciseOptionResponse) XXX_MessageName() string { return "syreen.options.MsgExerciseOptionResponse" }

// MsgCancelOption — writer cancels an unsold option and reclaims collateral
type MsgCancelOption struct {
	Writer   string `json:"writer"`
	OptionID uint64 `json:"option_id"`
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
