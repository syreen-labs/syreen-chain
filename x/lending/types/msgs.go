package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgDeposit deposits tokens into a lending pool
type MsgDeposit struct {
	Sender string   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	PoolID uint64   `protobuf:"varint,2,opt,name=pool_id,json=poolId,proto3" json:"pool_id"`
	Amount math.Int `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount"`
}

func (m *MsgDeposit) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return err
	}
	if m.Amount.IsNil() || !m.Amount.IsPositive() {
		return ErrInvalidAmount
	}
	return nil
}
func (m *MsgDeposit) ProtoMessage()             {}
func (m *MsgDeposit) Reset()                    { *m = MsgDeposit{} }
func (m *MsgDeposit) String() string            { return fmt.Sprintf("MsgDeposit{%s,%d,%s}", m.Sender, m.PoolID, m.Amount) }
func (m *MsgDeposit) XXX_MessageName() string   { return "syreen.lending.MsgDeposit" }

// MsgWithdraw withdraws tokens from a lending pool
type MsgWithdraw struct {
	Sender string   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	PoolID uint64   `protobuf:"varint,2,opt,name=pool_id,json=poolId,proto3" json:"pool_id"`
	Amount math.Int `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount"`
}

func (m *MsgWithdraw) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return err
	}
	if m.Amount.IsNil() || !m.Amount.IsPositive() {
		return ErrInvalidAmount
	}
	return nil
}
func (m *MsgWithdraw) ProtoMessage()             {}
func (m *MsgWithdraw) Reset()                    { *m = MsgWithdraw{} }
func (m *MsgWithdraw) String() string            { return fmt.Sprintf("MsgWithdraw{%s,%d,%s}", m.Sender, m.PoolID, m.Amount) }
func (m *MsgWithdraw) XXX_MessageName() string   { return "syreen.lending.MsgWithdraw" }

// MsgBorrow borrows tokens using deposited collateral
type MsgBorrow struct {
	Sender           string   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	BorrowPoolID     uint64   `protobuf:"varint,2,opt,name=borrow_pool_id,json=borrowPoolId,proto3" json:"borrow_pool_id"`
	Amount           math.Int `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount"`
	CollateralPoolID uint64   `protobuf:"varint,4,opt,name=collateral_pool_id,json=collateralPoolId,proto3" json:"collateral_pool_id"`
	CollateralAmount math.Int `protobuf:"bytes,5,opt,name=collateral_amount,json=collateralAmount,proto3" json:"collateral_amount"`
}

func (m *MsgBorrow) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return err
	}
	if m.Amount.IsNil() || !m.Amount.IsPositive() {
		return ErrInvalidAmount
	}
	if m.CollateralAmount.IsNil() || !m.CollateralAmount.IsPositive() {
		return ErrInvalidAmount
	}
	return nil
}
func (m *MsgBorrow) ProtoMessage()             {}
func (m *MsgBorrow) Reset()                    { *m = MsgBorrow{} }
func (m *MsgBorrow) String() string            { return fmt.Sprintf("MsgBorrow{%s,%d,%s}", m.Sender, m.BorrowPoolID, m.Amount) }
func (m *MsgBorrow) XXX_MessageName() string   { return "syreen.lending.MsgBorrow" }

// MsgRepay repays a borrow position
type MsgRepay struct {
	Sender   string   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	BorrowID uint64   `protobuf:"varint,2,opt,name=borrow_id,json=borrowId,proto3" json:"borrow_id"`
	Amount   math.Int `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount"`
}

func (m *MsgRepay) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return err
	}
	return nil
}
func (m *MsgRepay) ProtoMessage()             {}
func (m *MsgRepay) Reset()                    { *m = MsgRepay{} }
func (m *MsgRepay) String() string            { return fmt.Sprintf("MsgRepay{%s,%d}", m.Sender, m.BorrowID) }
func (m *MsgRepay) XXX_MessageName() string   { return "syreen.lending.MsgRepay" }

// MsgLiquidate liquidates an unhealthy borrow position
type MsgLiquidate struct {
	Liquidator string `protobuf:"bytes,1,opt,name=liquidator,proto3" json:"liquidator"`
	BorrowID   uint64 `protobuf:"varint,2,opt,name=borrow_id,json=borrowId,proto3" json:"borrow_id"`
}

func (m *MsgLiquidate) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Liquidator)
	return err
}
func (m *MsgLiquidate) ProtoMessage()             {}
func (m *MsgLiquidate) Reset()                    { *m = MsgLiquidate{} }
func (m *MsgLiquidate) String() string            { return fmt.Sprintf("MsgLiquidate{%s,%d}", m.Liquidator, m.BorrowID) }
func (m *MsgLiquidate) XXX_MessageName() string   { return "syreen.lending.MsgLiquidate" }

// MsgCreateLendingPool creates a new lending pool (governance only)
type MsgCreateLendingPool struct {
	Authority        string         `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority"`
	Denom            string         `protobuf:"bytes,2,opt,name=denom,proto3" json:"denom"`
	CollateralFactor math.LegacyDec `protobuf:"bytes,3,opt,name=collateral_factor,json=collateralFactor,proto3" json:"collateral_factor"`
	DexPoolID        uint64         `protobuf:"varint,4,opt,name=dex_pool_id,json=dexPoolId,proto3" json:"dex_pool_id"`
	PriceDenom       string         `protobuf:"bytes,5,opt,name=price_denom,json=priceDenom,proto3" json:"price_denom"`
}

func (m *MsgCreateLendingPool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Authority)
	return err
}
func (m *MsgCreateLendingPool) ProtoMessage()             {}
func (m *MsgCreateLendingPool) Reset()                    { *m = MsgCreateLendingPool{} }
func (m *MsgCreateLendingPool) String() string            { return fmt.Sprintf("MsgCreateLendingPool{%s}", m.Denom) }
func (m *MsgCreateLendingPool) XXX_MessageName() string   { return "syreen.lending.MsgCreateLendingPool" }

// Response types
type MsgDepositResponse struct {
	InterestEarned math.Int `protobuf:"bytes,1,opt,name=interest_earned,json=interestEarned,proto3" json:"interest_earned"`
}

func (m *MsgDepositResponse) ProtoMessage()             {}
func (m *MsgDepositResponse) Reset()                    { *m = MsgDepositResponse{} }
func (m *MsgDepositResponse) String() string            { return "MsgDepositResponse" }
func (m *MsgDepositResponse) XXX_MessageName() string   { return "syreen.lending.MsgDepositResponse" }

type MsgWithdrawResponse struct {
	InterestEarned math.Int `protobuf:"bytes,1,opt,name=interest_earned,json=interestEarned,proto3" json:"interest_earned"`
}

func (m *MsgWithdrawResponse) ProtoMessage()             {}
func (m *MsgWithdrawResponse) Reset()                    { *m = MsgWithdrawResponse{} }
func (m *MsgWithdrawResponse) String() string            { return "MsgWithdrawResponse" }
func (m *MsgWithdrawResponse) XXX_MessageName() string   { return "syreen.lending.MsgWithdrawResponse" }

type MsgBorrowResponse struct {
	BorrowID uint64 `protobuf:"varint,1,opt,name=borrow_id,json=borrowId,proto3" json:"borrow_id"`
}

func (m *MsgBorrowResponse) ProtoMessage()             {}
func (m *MsgBorrowResponse) Reset()                    { *m = MsgBorrowResponse{} }
func (m *MsgBorrowResponse) String() string            { return "MsgBorrowResponse" }
func (m *MsgBorrowResponse) XXX_MessageName() string   { return "syreen.lending.MsgBorrowResponse" }

type MsgRepayResponse struct {
	InterestPaid       math.Int `protobuf:"bytes,1,opt,name=interest_paid,json=interestPaid,proto3" json:"interest_paid"`
	CollateralReturned math.Int `protobuf:"bytes,2,opt,name=collateral_returned,json=collateralReturned,proto3" json:"collateral_returned"`
}

func (m *MsgRepayResponse) ProtoMessage()             {}
func (m *MsgRepayResponse) Reset()                    { *m = MsgRepayResponse{} }
func (m *MsgRepayResponse) String() string            { return "MsgRepayResponse" }
func (m *MsgRepayResponse) XXX_MessageName() string   { return "syreen.lending.MsgRepayResponse" }

type MsgLiquidateResponse struct {
	CollateralSeized math.Int `protobuf:"bytes,1,opt,name=collateral_seized,json=collateralSeized,proto3" json:"collateral_seized"`
	DebtRepaid       math.Int `protobuf:"bytes,2,opt,name=debt_repaid,json=debtRepaid,proto3" json:"debt_repaid"`
}

func (m *MsgLiquidateResponse) ProtoMessage()             {}
func (m *MsgLiquidateResponse) Reset()                    { *m = MsgLiquidateResponse{} }
func (m *MsgLiquidateResponse) String() string            { return "MsgLiquidateResponse" }
func (m *MsgLiquidateResponse) XXX_MessageName() string   { return "syreen.lending.MsgLiquidateResponse" }

type MsgCreateLendingPoolResponse struct {
	PoolID uint64 `protobuf:"varint,1,opt,name=pool_id,json=poolId,proto3" json:"pool_id"`
}

func (m *MsgCreateLendingPoolResponse) ProtoMessage()             {}
func (m *MsgCreateLendingPoolResponse) Reset()                    { *m = MsgCreateLendingPoolResponse{} }
func (m *MsgCreateLendingPoolResponse) String() string            { return "MsgCreateLendingPoolResponse" }
func (m *MsgCreateLendingPoolResponse) XXX_MessageName() string   { return "syreen.lending.MsgCreateLendingPoolResponse" }
