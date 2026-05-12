package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgFlashLoan borrows tokens from the flash pool, sends to borrower,
// then checks repayment (principal + fee) at end of handler.
type MsgFlashLoan struct {
	Sender string   `json:"sender"`
	Denom  string   `json:"denom"`
	Amount math.Int `json:"amount"`
}

func (m *MsgFlashLoan) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return err
	}
	if m.Denom == "" {
		return ErrInvalidDenom
	}
	if m.Amount.IsNil() || !m.Amount.IsPositive() {
		return ErrInvalidAmount
	}
	return nil
}
func (m *MsgFlashLoan) ProtoMessage()            {}
func (m *MsgFlashLoan) Reset()                   { *m = MsgFlashLoan{} }
func (m *MsgFlashLoan) String() string           { return fmt.Sprintf("MsgFlashLoan{%s,%s,%s}", m.Sender, m.Denom, m.Amount) }
func (m *MsgFlashLoan) XXX_MessageName() string  { return "syreen.flashloan.MsgFlashLoan" }

// MsgCreateFlashPool creates a new flash loan pool (authority only)
type MsgCreateFlashPool struct {
	Authority string         `json:"authority"`
	Denom     string         `json:"denom"`
	FeeRate   math.LegacyDec `json:"fee_rate"`
}

func (m *MsgCreateFlashPool) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return err
	}
	if m.Denom == "" {
		return ErrInvalidDenom
	}
	if m.FeeRate.IsNil() || m.FeeRate.IsNegative() {
		return ErrInvalidFeeRate
	}
	// Enforce minimum fee rate of 0.01% to prevent free flash loans.
	minFeeRate := math.LegacyNewDecWithPrec(1, 4) // 0.0001
	if m.FeeRate.LT(minFeeRate) {
		return fmt.Errorf("fee rate %s is below minimum %s", m.FeeRate, minFeeRate)
	}
	return nil
}
func (m *MsgCreateFlashPool) ProtoMessage()            {}
func (m *MsgCreateFlashPool) Reset()                   { *m = MsgCreateFlashPool{} }
func (m *MsgCreateFlashPool) String() string           { return fmt.Sprintf("MsgCreateFlashPool{%s,%s}", m.Authority, m.Denom) }
func (m *MsgCreateFlashPool) XXX_MessageName() string  { return "syreen.flashloan.MsgCreateFlashPool" }

// MsgFundFlashPool adds liquidity to a flash loan pool
type MsgFundFlashPool struct {
	Sender string   `json:"sender"`
	Denom  string   `json:"denom"`
	Amount math.Int `json:"amount"`
}

func (m *MsgFundFlashPool) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return err
	}
	if m.Denom == "" {
		return ErrInvalidDenom
	}
	if m.Amount.IsNil() || !m.Amount.IsPositive() {
		return ErrInvalidAmount
	}
	return nil
}
func (m *MsgFundFlashPool) ProtoMessage()            {}
func (m *MsgFundFlashPool) Reset()                   { *m = MsgFundFlashPool{} }
func (m *MsgFundFlashPool) String() string           { return fmt.Sprintf("MsgFundFlashPool{%s,%s,%s}", m.Sender, m.Denom, m.Amount) }
func (m *MsgFundFlashPool) XXX_MessageName() string  { return "syreen.flashloan.MsgFundFlashPool" }

// MsgWithdrawFlashPool redeems LP shares from a flash loan pool (H-1 fix).
type MsgWithdrawFlashPool struct {
	Sender string   `json:"sender"`
	Denom  string   `json:"denom"`
	Shares math.Int `json:"shares"`
}

func (m *MsgWithdrawFlashPool) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return err
	}
	if m.Denom == "" {
		return ErrInvalidDenom
	}
	if m.Shares.IsNil() || !m.Shares.IsPositive() {
		return ErrInvalidAmount
	}
	return nil
}
func (m *MsgWithdrawFlashPool) ProtoMessage()           {}
func (m *MsgWithdrawFlashPool) Reset()                  { *m = MsgWithdrawFlashPool{} }
func (m *MsgWithdrawFlashPool) String() string          { return fmt.Sprintf("MsgWithdrawFlashPool{%s,%s,%s}", m.Sender, m.Denom, m.Shares) }
func (m *MsgWithdrawFlashPool) XXX_MessageName() string { return "syreen.flashloan.MsgWithdrawFlashPool" }

// Response types
type MsgFlashLoanResponse struct {
	Fee math.Int `json:"fee"`
}

func (m *MsgFlashLoanResponse) ProtoMessage()            {}
func (m *MsgFlashLoanResponse) Reset()                   { *m = MsgFlashLoanResponse{} }
func (m *MsgFlashLoanResponse) String() string           { return "MsgFlashLoanResponse" }
func (m *MsgFlashLoanResponse) XXX_MessageName() string  { return "syreen.flashloan.MsgFlashLoanResponse" }

type MsgCreateFlashPoolResponse struct{}

func (m *MsgCreateFlashPoolResponse) ProtoMessage()            {}
func (m *MsgCreateFlashPoolResponse) Reset()                   { *m = MsgCreateFlashPoolResponse{} }
func (m *MsgCreateFlashPoolResponse) String() string           { return "MsgCreateFlashPoolResponse" }
func (m *MsgCreateFlashPoolResponse) XXX_MessageName() string  { return "syreen.flashloan.MsgCreateFlashPoolResponse" }

type MsgFundFlashPoolResponse struct{}

func (m *MsgFundFlashPoolResponse) ProtoMessage()           {}
func (m *MsgFundFlashPoolResponse) Reset()                  { *m = MsgFundFlashPoolResponse{} }
func (m *MsgFundFlashPoolResponse) String() string          { return "MsgFundFlashPoolResponse" }
func (m *MsgFundFlashPoolResponse) XXX_MessageName() string { return "syreen.flashloan.MsgFundFlashPoolResponse" }

type MsgWithdrawFlashPoolResponse struct {
	AmountReturned math.Int `json:"amount_returned"`
}

func (m *MsgWithdrawFlashPoolResponse) ProtoMessage()           {}
func (m *MsgWithdrawFlashPoolResponse) Reset()                  { *m = MsgWithdrawFlashPoolResponse{} }
func (m *MsgWithdrawFlashPoolResponse) String() string          { return "MsgWithdrawFlashPoolResponse" }
func (m *MsgWithdrawFlashPoolResponse) XXX_MessageName() string { return "syreen.flashloan.MsgWithdrawFlashPoolResponse" }
