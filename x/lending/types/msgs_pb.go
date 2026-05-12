package types

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math/bits"

	proto "github.com/cosmos/gogoproto/proto"
)

// Reference imports
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = io.ErrUnexpectedEOF
var _ = binary.BigEndian
var _ = bits.Len64

var fileDescriptorLendingTx []byte

// XXX methods required by gRPC decoder
func (m *MsgDeposit) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgDeposit) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgDeposit) XXX_Size() int { return m.Size() }

func (m *MsgDepositResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgDepositResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgDepositResponse) XXX_Size() int { return m.Size() }

func (m *MsgWithdraw) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgWithdraw) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgWithdraw) XXX_Size() int { return m.Size() }

func (m *MsgWithdrawResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgWithdrawResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgWithdrawResponse) XXX_Size() int { return m.Size() }

func (m *MsgBorrow) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgBorrow) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgBorrow) XXX_Size() int { return m.Size() }

func (m *MsgBorrowResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgBorrowResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgBorrowResponse) XXX_Size() int { return m.Size() }

func (m *MsgRepay) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRepay) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgRepay) XXX_Size() int { return m.Size() }

func (m *MsgRepayResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRepayResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgRepayResponse) XXX_Size() int { return m.Size() }

func (m *MsgLiquidate) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgLiquidate) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgLiquidate) XXX_Size() int { return m.Size() }

func (m *MsgLiquidateResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgLiquidateResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgLiquidateResponse) XXX_Size() int { return m.Size() }

func (m *MsgCreateLendingPool) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateLendingPool) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCreateLendingPool) XXX_Size() int { return m.Size() }

func (m *MsgCreateLendingPoolResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateLendingPoolResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCreateLendingPoolResponse) XXX_Size() int { return m.Size() }

// Descriptor methods
func (*MsgDeposit) Descriptor() ([]byte, []int)                  { return fileDescriptorLendingTx, []int{0} }
func (*MsgDepositResponse) Descriptor() ([]byte, []int)          { return fileDescriptorLendingTx, []int{1} }
func (*MsgWithdraw) Descriptor() ([]byte, []int)                 { return fileDescriptorLendingTx, []int{2} }
func (*MsgWithdrawResponse) Descriptor() ([]byte, []int)         { return fileDescriptorLendingTx, []int{3} }
func (*MsgBorrow) Descriptor() ([]byte, []int)                   { return fileDescriptorLendingTx, []int{4} }
func (*MsgBorrowResponse) Descriptor() ([]byte, []int)           { return fileDescriptorLendingTx, []int{5} }
func (*MsgRepay) Descriptor() ([]byte, []int)                    { return fileDescriptorLendingTx, []int{6} }
func (*MsgRepayResponse) Descriptor() ([]byte, []int)            { return fileDescriptorLendingTx, []int{7} }
func (*MsgLiquidate) Descriptor() ([]byte, []int)                { return fileDescriptorLendingTx, []int{8} }
func (*MsgLiquidateResponse) Descriptor() ([]byte, []int)        { return fileDescriptorLendingTx, []int{9} }
func (*MsgCreateLendingPool) Descriptor() ([]byte, []int)        { return fileDescriptorLendingTx, []int{10} }
func (*MsgCreateLendingPoolResponse) Descriptor() ([]byte, []int) { return fileDescriptorLendingTx, []int{11} }

func init() {
	proto.RegisterType((*MsgDeposit)(nil), "syreen.lending.MsgDeposit")
	proto.RegisterType((*MsgDepositResponse)(nil), "syreen.lending.MsgDepositResponse")
	proto.RegisterType((*MsgWithdraw)(nil), "syreen.lending.MsgWithdraw")
	proto.RegisterType((*MsgWithdrawResponse)(nil), "syreen.lending.MsgWithdrawResponse")
	proto.RegisterType((*MsgBorrow)(nil), "syreen.lending.MsgBorrow")
	proto.RegisterType((*MsgBorrowResponse)(nil), "syreen.lending.MsgBorrowResponse")
	proto.RegisterType((*MsgRepay)(nil), "syreen.lending.MsgRepay")
	proto.RegisterType((*MsgRepayResponse)(nil), "syreen.lending.MsgRepayResponse")
	proto.RegisterType((*MsgLiquidate)(nil), "syreen.lending.MsgLiquidate")
	proto.RegisterType((*MsgLiquidateResponse)(nil), "syreen.lending.MsgLiquidateResponse")
	proto.RegisterType((*MsgCreateLendingPool)(nil), "syreen.lending.MsgCreateLendingPool")
	proto.RegisterType((*MsgCreateLendingPoolResponse)(nil), "syreen.lending.MsgCreateLendingPoolResponse")
}

// Marshal/Unmarshal/Size for all message types using JSON
func (m *MsgDeposit) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgDeposit) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgDeposit) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgDeposit) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgDeposit) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgDepositResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgDepositResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgDepositResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgDepositResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgDepositResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgWithdraw) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgWithdraw) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgWithdraw) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgWithdraw) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgWithdraw) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgWithdrawResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgWithdrawResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgWithdrawResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgWithdrawResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgWithdrawResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgBorrow) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgBorrow) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgBorrow) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgBorrow) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgBorrow) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgBorrowResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgBorrowResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgBorrowResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgBorrowResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgBorrowResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgRepay) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgRepay) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgRepay) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgRepay) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgRepay) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgRepayResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgRepayResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgRepayResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgRepayResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgRepayResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgLiquidate) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgLiquidate) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgLiquidate) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgLiquidate) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgLiquidate) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgLiquidateResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgLiquidateResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgLiquidateResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgLiquidateResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgLiquidateResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgCreateLendingPool) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgCreateLendingPool) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCreateLendingPool) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCreateLendingPool) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCreateLendingPool) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgCreateLendingPoolResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgCreateLendingPoolResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCreateLendingPoolResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCreateLendingPoolResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCreateLendingPoolResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }
