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

var fileDescriptorLaunchpadTx []byte

// XXX methods required by gRPC decoder
func (m *MsgCreateLaunch) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateLaunch) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCreateLaunch) XXX_Size() int { return m.Size() }

func (m *MsgCreateLaunchResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateLaunchResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCreateLaunchResponse) XXX_Size() int { return m.Size() }

func (m *MsgContribute) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgContribute) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgContribute) XXX_Size() int { return m.Size() }

func (m *MsgContributeResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgContributeResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgContributeResponse) XXX_Size() int { return m.Size() }

func (m *MsgClaimTokens) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgClaimTokens) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgClaimTokens) XXX_Size() int { return m.Size() }

func (m *MsgClaimTokensResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgClaimTokensResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgClaimTokensResponse) XXX_Size() int { return m.Size() }

func (m *MsgClaimRefund) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgClaimRefund) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgClaimRefund) XXX_Size() int { return m.Size() }

func (m *MsgClaimRefundResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgClaimRefundResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgClaimRefundResponse) XXX_Size() int { return m.Size() }

func (m *MsgFinalizeLaunch) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgFinalizeLaunch) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgFinalizeLaunch) XXX_Size() int { return m.Size() }

func (m *MsgFinalizeLaunchResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgFinalizeLaunchResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgFinalizeLaunchResponse) XXX_Size() int { return m.Size() }

// Descriptor methods
func (*MsgCreateLaunch) Descriptor() ([]byte, []int)             { return fileDescriptorLaunchpadTx, []int{0} }
func (*MsgCreateLaunchResponse) Descriptor() ([]byte, []int)     { return fileDescriptorLaunchpadTx, []int{1} }
func (*MsgContribute) Descriptor() ([]byte, []int)               { return fileDescriptorLaunchpadTx, []int{2} }
func (*MsgContributeResponse) Descriptor() ([]byte, []int)       { return fileDescriptorLaunchpadTx, []int{3} }
func (*MsgClaimTokens) Descriptor() ([]byte, []int)              { return fileDescriptorLaunchpadTx, []int{4} }
func (*MsgClaimTokensResponse) Descriptor() ([]byte, []int)      { return fileDescriptorLaunchpadTx, []int{5} }
func (*MsgClaimRefund) Descriptor() ([]byte, []int)              { return fileDescriptorLaunchpadTx, []int{6} }
func (*MsgClaimRefundResponse) Descriptor() ([]byte, []int)      { return fileDescriptorLaunchpadTx, []int{7} }
func (*MsgFinalizeLaunch) Descriptor() ([]byte, []int)           { return fileDescriptorLaunchpadTx, []int{8} }
func (*MsgFinalizeLaunchResponse) Descriptor() ([]byte, []int)   { return fileDescriptorLaunchpadTx, []int{9} }

func init() {
	proto.RegisterType((*MsgCreateLaunch)(nil), "syreen.launchpad.MsgCreateLaunch")
	proto.RegisterType((*MsgCreateLaunchResponse)(nil), "syreen.launchpad.MsgCreateLaunchResponse")
	proto.RegisterType((*MsgContribute)(nil), "syreen.launchpad.MsgContribute")
	proto.RegisterType((*MsgContributeResponse)(nil), "syreen.launchpad.MsgContributeResponse")
	proto.RegisterType((*MsgClaimTokens)(nil), "syreen.launchpad.MsgClaimTokens")
	proto.RegisterType((*MsgClaimTokensResponse)(nil), "syreen.launchpad.MsgClaimTokensResponse")
	proto.RegisterType((*MsgClaimRefund)(nil), "syreen.launchpad.MsgClaimRefund")
	proto.RegisterType((*MsgClaimRefundResponse)(nil), "syreen.launchpad.MsgClaimRefundResponse")
	proto.RegisterType((*MsgFinalizeLaunch)(nil), "syreen.launchpad.MsgFinalizeLaunch")
	proto.RegisterType((*MsgFinalizeLaunchResponse)(nil), "syreen.launchpad.MsgFinalizeLaunchResponse")
}

// Marshal/Unmarshal/Size for all message types using JSON
func (m *MsgCreateLaunch) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgCreateLaunch) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCreateLaunch) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCreateLaunch) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCreateLaunch) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgCreateLaunchResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgCreateLaunchResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCreateLaunchResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCreateLaunchResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCreateLaunchResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgContribute) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgContribute) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgContribute) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgContribute) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgContribute) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgContributeResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgContributeResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgContributeResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgContributeResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgContributeResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgClaimTokens) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgClaimTokens) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgClaimTokens) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgClaimTokens) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgClaimTokens) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgClaimTokensResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgClaimTokensResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgClaimTokensResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgClaimTokensResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgClaimTokensResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgClaimRefund) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgClaimRefund) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgClaimRefund) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgClaimRefund) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgClaimRefund) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgClaimRefundResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgClaimRefundResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgClaimRefundResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgClaimRefundResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgClaimRefundResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgFinalizeLaunch) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgFinalizeLaunch) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgFinalizeLaunch) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgFinalizeLaunch) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgFinalizeLaunch) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgFinalizeLaunchResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgFinalizeLaunchResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgFinalizeLaunchResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgFinalizeLaunchResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgFinalizeLaunchResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }
