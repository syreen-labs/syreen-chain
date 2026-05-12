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

var fileDescriptorFlashloanTx []byte

// XXX methods required by gRPC decoder
func (m *MsgFlashLoan) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgFlashLoan) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgFlashLoan) XXX_Size() int                         { return m.Size() }

func (m *MsgFlashLoanResponse) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgFlashLoanResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgFlashLoanResponse) XXX_Size() int                         { return m.Size() }

func (m *MsgCreateFlashPool) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgCreateFlashPool) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCreateFlashPool) XXX_Size() int                         { return m.Size() }

func (m *MsgCreateFlashPoolResponse) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgCreateFlashPoolResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCreateFlashPoolResponse) XXX_Size() int                         { return m.Size() }

func (m *MsgFundFlashPool) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgFundFlashPool) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgFundFlashPool) XXX_Size() int                         { return m.Size() }

func (m *MsgFundFlashPoolResponse) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgFundFlashPoolResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgFundFlashPoolResponse) XXX_Size() int                         { return m.Size() }

func (m *MsgWithdrawFlashPool) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgWithdrawFlashPool) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgWithdrawFlashPool) XXX_Size() int                         { return m.Size() }

func (m *MsgWithdrawFlashPoolResponse) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgWithdrawFlashPoolResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgWithdrawFlashPoolResponse) XXX_Size() int                         { return m.Size() }

// Descriptor methods
func (*MsgFlashLoan) Descriptor() ([]byte, []int)                  { return fileDescriptorFlashloanTx, []int{0} }
func (*MsgFlashLoanResponse) Descriptor() ([]byte, []int)          { return fileDescriptorFlashloanTx, []int{1} }
func (*MsgCreateFlashPool) Descriptor() ([]byte, []int)            { return fileDescriptorFlashloanTx, []int{2} }
func (*MsgCreateFlashPoolResponse) Descriptor() ([]byte, []int)    { return fileDescriptorFlashloanTx, []int{3} }
func (*MsgFundFlashPool) Descriptor() ([]byte, []int)              { return fileDescriptorFlashloanTx, []int{4} }
func (*MsgFundFlashPoolResponse) Descriptor() ([]byte, []int)      { return fileDescriptorFlashloanTx, []int{5} }
func (*MsgWithdrawFlashPool) Descriptor() ([]byte, []int)          { return fileDescriptorFlashloanTx, []int{6} }
func (*MsgWithdrawFlashPoolResponse) Descriptor() ([]byte, []int)  { return fileDescriptorFlashloanTx, []int{7} }

func init() {
	proto.RegisterType((*MsgFlashLoan)(nil), "syreen.flashloan.MsgFlashLoan")
	proto.RegisterType((*MsgFlashLoanResponse)(nil), "syreen.flashloan.MsgFlashLoanResponse")
	proto.RegisterType((*MsgCreateFlashPool)(nil), "syreen.flashloan.MsgCreateFlashPool")
	proto.RegisterType((*MsgCreateFlashPoolResponse)(nil), "syreen.flashloan.MsgCreateFlashPoolResponse")
	proto.RegisterType((*MsgFundFlashPool)(nil), "syreen.flashloan.MsgFundFlashPool")
	proto.RegisterType((*MsgFundFlashPoolResponse)(nil), "syreen.flashloan.MsgFundFlashPoolResponse")
	proto.RegisterType((*MsgWithdrawFlashPool)(nil), "syreen.flashloan.MsgWithdrawFlashPool")
	proto.RegisterType((*MsgWithdrawFlashPoolResponse)(nil), "syreen.flashloan.MsgWithdrawFlashPoolResponse")
}

// Marshal/Unmarshal/Size for all message types using JSON
func (m *MsgFlashLoan) Marshal() ([]byte, error)                             { return json.Marshal(m) }
func (m *MsgFlashLoan) MarshalTo(dAtA []byte) (int, error)                   { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgFlashLoan) MarshalToSizedBuffer(dAtA []byte) (int, error)        { return m.MarshalTo(dAtA) }
func (m *MsgFlashLoan) Size() int                                            { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgFlashLoan) Unmarshal(bz []byte) error                            { return json.Unmarshal(bz, m) }

func (m *MsgFlashLoanResponse) Marshal() ([]byte, error)                     { return json.Marshal(m) }
func (m *MsgFlashLoanResponse) MarshalTo(dAtA []byte) (int, error)           { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgFlashLoanResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgFlashLoanResponse) Size() int                                    { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgFlashLoanResponse) Unmarshal(bz []byte) error                    { return json.Unmarshal(bz, m) }

func (m *MsgCreateFlashPool) Marshal() ([]byte, error)                       { return json.Marshal(m) }
func (m *MsgCreateFlashPool) MarshalTo(dAtA []byte) (int, error)             { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCreateFlashPool) MarshalToSizedBuffer(dAtA []byte) (int, error)  { return m.MarshalTo(dAtA) }
func (m *MsgCreateFlashPool) Size() int                                      { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCreateFlashPool) Unmarshal(bz []byte) error                      { return json.Unmarshal(bz, m) }

func (m *MsgCreateFlashPoolResponse) Marshal() ([]byte, error)               { return json.Marshal(m) }
func (m *MsgCreateFlashPoolResponse) MarshalTo(dAtA []byte) (int, error)     { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCreateFlashPoolResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCreateFlashPoolResponse) Size() int                              { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCreateFlashPoolResponse) Unmarshal(bz []byte) error              { return json.Unmarshal(bz, m) }

func (m *MsgFundFlashPool) Marshal() ([]byte, error)                         { return json.Marshal(m) }
func (m *MsgFundFlashPool) MarshalTo(dAtA []byte) (int, error)               { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgFundFlashPool) MarshalToSizedBuffer(dAtA []byte) (int, error)    { return m.MarshalTo(dAtA) }
func (m *MsgFundFlashPool) Size() int                                        { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgFundFlashPool) Unmarshal(bz []byte) error                        { return json.Unmarshal(bz, m) }

func (m *MsgFundFlashPoolResponse) Marshal() ([]byte, error)                 { return json.Marshal(m) }
func (m *MsgFundFlashPoolResponse) MarshalTo(dAtA []byte) (int, error)       { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgFundFlashPoolResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgFundFlashPoolResponse) Size() int                                { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgFundFlashPoolResponse) Unmarshal(bz []byte) error                { return json.Unmarshal(bz, m) }

func (m *MsgWithdrawFlashPool) Marshal() ([]byte, error)                     { return json.Marshal(m) }
func (m *MsgWithdrawFlashPool) MarshalTo(dAtA []byte) (int, error)           { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgWithdrawFlashPool) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgWithdrawFlashPool) Size() int                                    { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgWithdrawFlashPool) Unmarshal(bz []byte) error                   { return json.Unmarshal(bz, m) }

func (m *MsgWithdrawFlashPoolResponse) Marshal() ([]byte, error)             { return json.Marshal(m) }
func (m *MsgWithdrawFlashPoolResponse) MarshalTo(dAtA []byte) (int, error)   { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgWithdrawFlashPoolResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgWithdrawFlashPoolResponse) Size() int                            { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgWithdrawFlashPoolResponse) Unmarshal(bz []byte) error            { return json.Unmarshal(bz, m) }
