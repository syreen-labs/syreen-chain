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

var fileDescriptorVaultTx []byte

// XXX methods required by gRPC decoder
func (m *MsgCreateVault) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateVault) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCreateVault) XXX_Size() int { return m.Size() }

func (m *MsgCreateVaultResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateVaultResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCreateVaultResponse) XXX_Size() int { return m.Size() }

func (m *MsgDepositVault) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgDepositVault) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgDepositVault) XXX_Size() int { return m.Size() }

func (m *MsgDepositVaultResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgDepositVaultResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgDepositVaultResponse) XXX_Size() int { return m.Size() }

func (m *MsgWithdrawVault) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgWithdrawVault) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgWithdrawVault) XXX_Size() int { return m.Size() }

func (m *MsgWithdrawVaultResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgWithdrawVaultResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgWithdrawVaultResponse) XXX_Size() int { return m.Size() }

func (m *MsgCompoundVault) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCompoundVault) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCompoundVault) XXX_Size() int { return m.Size() }

func (m *MsgCompoundVaultResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCompoundVaultResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCompoundVaultResponse) XXX_Size() int { return m.Size() }

func (m *MsgUpdateVaultStrategy) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgUpdateVaultStrategy) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgUpdateVaultStrategy) XXX_Size() int { return m.Size() }

func (m *MsgUpdateVaultStrategyResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgUpdateVaultStrategyResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgUpdateVaultStrategyResponse) XXX_Size() int { return m.Size() }

// Descriptor methods
func (*MsgCreateVault) Descriptor() ([]byte, []int)                    { return fileDescriptorVaultTx, []int{0} }
func (*MsgCreateVaultResponse) Descriptor() ([]byte, []int)            { return fileDescriptorVaultTx, []int{1} }
func (*MsgDepositVault) Descriptor() ([]byte, []int)                   { return fileDescriptorVaultTx, []int{2} }
func (*MsgDepositVaultResponse) Descriptor() ([]byte, []int)           { return fileDescriptorVaultTx, []int{3} }
func (*MsgWithdrawVault) Descriptor() ([]byte, []int)                  { return fileDescriptorVaultTx, []int{4} }
func (*MsgWithdrawVaultResponse) Descriptor() ([]byte, []int)          { return fileDescriptorVaultTx, []int{5} }
func (*MsgCompoundVault) Descriptor() ([]byte, []int)                  { return fileDescriptorVaultTx, []int{6} }
func (*MsgCompoundVaultResponse) Descriptor() ([]byte, []int)          { return fileDescriptorVaultTx, []int{7} }
func (*MsgUpdateVaultStrategy) Descriptor() ([]byte, []int)            { return fileDescriptorVaultTx, []int{8} }
func (*MsgUpdateVaultStrategyResponse) Descriptor() ([]byte, []int)    { return fileDescriptorVaultTx, []int{9} }

func init() {
	proto.RegisterType((*MsgCreateVault)(nil), "syreen.vault.MsgCreateVault")
	proto.RegisterType((*MsgCreateVaultResponse)(nil), "syreen.vault.MsgCreateVaultResponse")
	proto.RegisterType((*MsgDepositVault)(nil), "syreen.vault.MsgDepositVault")
	proto.RegisterType((*MsgDepositVaultResponse)(nil), "syreen.vault.MsgDepositVaultResponse")
	proto.RegisterType((*MsgWithdrawVault)(nil), "syreen.vault.MsgWithdrawVault")
	proto.RegisterType((*MsgWithdrawVaultResponse)(nil), "syreen.vault.MsgWithdrawVaultResponse")
	proto.RegisterType((*MsgCompoundVault)(nil), "syreen.vault.MsgCompoundVault")
	proto.RegisterType((*MsgCompoundVaultResponse)(nil), "syreen.vault.MsgCompoundVaultResponse")
	proto.RegisterType((*MsgUpdateVaultStrategy)(nil), "syreen.vault.MsgUpdateVaultStrategy")
	proto.RegisterType((*MsgUpdateVaultStrategyResponse)(nil), "syreen.vault.MsgUpdateVaultStrategyResponse")
}

// Marshal/Unmarshal/Size for all message types using JSON
func (m *MsgCreateVault) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgCreateVault) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCreateVault) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCreateVault) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCreateVault) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgCreateVaultResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgCreateVaultResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCreateVaultResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCreateVaultResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCreateVaultResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgDepositVault) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgDepositVault) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgDepositVault) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgDepositVault) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgDepositVault) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgDepositVaultResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgDepositVaultResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgDepositVaultResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgDepositVaultResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgDepositVaultResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgWithdrawVault) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgWithdrawVault) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgWithdrawVault) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgWithdrawVault) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgWithdrawVault) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgWithdrawVaultResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgWithdrawVaultResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgWithdrawVaultResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgWithdrawVaultResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgWithdrawVaultResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgCompoundVault) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgCompoundVault) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCompoundVault) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCompoundVault) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCompoundVault) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgCompoundVaultResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgCompoundVaultResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCompoundVaultResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCompoundVaultResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCompoundVaultResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgUpdateVaultStrategy) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgUpdateVaultStrategy) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgUpdateVaultStrategy) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgUpdateVaultStrategy) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgUpdateVaultStrategy) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgUpdateVaultStrategyResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgUpdateVaultStrategyResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgUpdateVaultStrategyResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgUpdateVaultStrategyResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgUpdateVaultStrategyResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }
