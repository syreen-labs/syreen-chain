package types

import (
	"encoding/json"
	"encoding/binary"
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

var fileDescriptorFarmingTx []byte

// Proto methods for all message types
func (m *MsgStake) ProtoMessage()             {}
func (m *MsgStake) Reset()                    { *m = MsgStake{} }
func (m *MsgStake) String() string            { return fmt.Sprintf("MsgStake{%s, %d, %s}", m.Sender, m.PoolID, m.Amount) }

func (m *MsgUnstake) ProtoMessage()           {}
func (m *MsgUnstake) Reset()                  { *m = MsgUnstake{} }
func (m *MsgUnstake) String() string          { return fmt.Sprintf("MsgUnstake{%s, %d, %s}", m.Sender, m.PoolID, m.Amount) }

func (m *MsgClaimReward) ProtoMessage()       {}
func (m *MsgClaimReward) Reset()              { *m = MsgClaimReward{} }
func (m *MsgClaimReward) String() string      { return fmt.Sprintf("MsgClaimReward{%s, %d}", m.Sender, m.PoolID) }

func (m *MsgCreateFarm) ProtoMessage()        {}
func (m *MsgCreateFarm) Reset()               { *m = MsgCreateFarm{} }
func (m *MsgCreateFarm) String() string       { return fmt.Sprintf("MsgCreateFarm{%d}", m.PoolID) }

func (m *MsgUpdateFarm) ProtoMessage()        {}
func (m *MsgUpdateFarm) Reset()               { *m = MsgUpdateFarm{} }
func (m *MsgUpdateFarm) String() string       { return fmt.Sprintf("MsgUpdateFarm{%d}", m.PoolID) }

func (m *MsgStakeResponse) ProtoMessage()        {}
func (m *MsgStakeResponse) Reset()               { *m = MsgStakeResponse{} }
func (m *MsgStakeResponse) String() string       { return "MsgStakeResponse" }

func (m *MsgUnstakeResponse) ProtoMessage()      {}
func (m *MsgUnstakeResponse) Reset()             { *m = MsgUnstakeResponse{} }
func (m *MsgUnstakeResponse) String() string     { return "MsgUnstakeResponse" }

func (m *MsgClaimRewardResponse) ProtoMessage()  {}
func (m *MsgClaimRewardResponse) Reset()         { *m = MsgClaimRewardResponse{} }
func (m *MsgClaimRewardResponse) String() string { return "MsgClaimRewardResponse" }

func (m *MsgCreateFarmResponse) ProtoMessage()   {}
func (m *MsgCreateFarmResponse) Reset()          { *m = MsgCreateFarmResponse{} }
func (m *MsgCreateFarmResponse) String() string  { return "MsgCreateFarmResponse" }

func (m *MsgUpdateFarmResponse) ProtoMessage()   {}
func (m *MsgUpdateFarmResponse) Reset()          { *m = MsgUpdateFarmResponse{} }
func (m *MsgUpdateFarmResponse) String() string  { return "MsgUpdateFarmResponse" }

// XXX_MessageName for type URL resolution
func (m *MsgStake) XXX_MessageName() string              { return "syreen.farming.MsgStake" }
func (m *MsgStakeResponse) XXX_MessageName() string      { return "syreen.farming.MsgStakeResponse" }
func (m *MsgUnstake) XXX_MessageName() string            { return "syreen.farming.MsgUnstake" }
func (m *MsgUnstakeResponse) XXX_MessageName() string    { return "syreen.farming.MsgUnstakeResponse" }
func (m *MsgClaimReward) XXX_MessageName() string        { return "syreen.farming.MsgClaimReward" }
func (m *MsgClaimRewardResponse) XXX_MessageName() string { return "syreen.farming.MsgClaimRewardResponse" }
func (m *MsgCreateFarm) XXX_MessageName() string         { return "syreen.farming.MsgCreateFarm" }
func (m *MsgCreateFarmResponse) XXX_MessageName() string { return "syreen.farming.MsgCreateFarmResponse" }
func (m *MsgUpdateFarm) XXX_MessageName() string         { return "syreen.farming.MsgUpdateFarm" }
func (m *MsgUpdateFarmResponse) XXX_MessageName() string { return "syreen.farming.MsgUpdateFarmResponse" }

// XXX methods required by gRPC decoder
func (m *MsgStake) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgStake) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgStake) XXX_Size() int { return m.Size() }

func (m *MsgStakeResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgStakeResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgStakeResponse) XXX_Size() int { return m.Size() }

func (m *MsgUnstake) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgUnstake) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgUnstake) XXX_Size() int { return m.Size() }

func (m *MsgUnstakeResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgUnstakeResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgUnstakeResponse) XXX_Size() int { return m.Size() }

func (m *MsgClaimReward) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgClaimReward) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgClaimReward) XXX_Size() int { return m.Size() }

func (m *MsgClaimRewardResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgClaimRewardResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgClaimRewardResponse) XXX_Size() int { return m.Size() }

func (m *MsgCreateFarm) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateFarm) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCreateFarm) XXX_Size() int { return m.Size() }

func (m *MsgCreateFarmResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateFarmResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCreateFarmResponse) XXX_Size() int { return m.Size() }

func (m *MsgUpdateFarm) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgUpdateFarm) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgUpdateFarm) XXX_Size() int { return m.Size() }

func (m *MsgUpdateFarmResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgUpdateFarmResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgUpdateFarmResponse) XXX_Size() int { return m.Size() }

// Descriptor methods
func (*MsgStake) Descriptor() ([]byte, []int)              { return fileDescriptorFarmingTx, []int{0} }
func (*MsgStakeResponse) Descriptor() ([]byte, []int)      { return fileDescriptorFarmingTx, []int{1} }
func (*MsgUnstake) Descriptor() ([]byte, []int)            { return fileDescriptorFarmingTx, []int{2} }
func (*MsgUnstakeResponse) Descriptor() ([]byte, []int)    { return fileDescriptorFarmingTx, []int{3} }
func (*MsgClaimReward) Descriptor() ([]byte, []int)        { return fileDescriptorFarmingTx, []int{4} }
func (*MsgClaimRewardResponse) Descriptor() ([]byte, []int){ return fileDescriptorFarmingTx, []int{5} }
func (*MsgCreateFarm) Descriptor() ([]byte, []int)         { return fileDescriptorFarmingTx, []int{6} }
func (*MsgCreateFarmResponse) Descriptor() ([]byte, []int) { return fileDescriptorFarmingTx, []int{7} }
func (*MsgUpdateFarm) Descriptor() ([]byte, []int)         { return fileDescriptorFarmingTx, []int{8} }
func (*MsgUpdateFarmResponse) Descriptor() ([]byte, []int) { return fileDescriptorFarmingTx, []int{9} }

func init() {
	proto.RegisterType((*MsgStake)(nil), "syreen.farming.MsgStake")
	proto.RegisterType((*MsgStakeResponse)(nil), "syreen.farming.MsgStakeResponse")
	proto.RegisterType((*MsgUnstake)(nil), "syreen.farming.MsgUnstake")
	proto.RegisterType((*MsgUnstakeResponse)(nil), "syreen.farming.MsgUnstakeResponse")
	proto.RegisterType((*MsgClaimReward)(nil), "syreen.farming.MsgClaimReward")
	proto.RegisterType((*MsgClaimRewardResponse)(nil), "syreen.farming.MsgClaimRewardResponse")
	proto.RegisterType((*MsgCreateFarm)(nil), "syreen.farming.MsgCreateFarm")
	proto.RegisterType((*MsgCreateFarmResponse)(nil), "syreen.farming.MsgCreateFarmResponse")
	proto.RegisterType((*MsgUpdateFarm)(nil), "syreen.farming.MsgUpdateFarm")
	proto.RegisterType((*MsgUpdateFarmResponse)(nil), "syreen.farming.MsgUpdateFarmResponse")
}

// Marshal/Unmarshal/Size for all message types using JSON
func (m *MsgStake) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgStake) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgStake) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgStake) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgStake) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgStakeResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgStakeResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgStakeResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgStakeResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgStakeResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgUnstake) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgUnstake) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgUnstake) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgUnstake) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgUnstake) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgUnstakeResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgUnstakeResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgUnstakeResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgUnstakeResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgUnstakeResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgClaimReward) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgClaimReward) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgClaimReward) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgClaimReward) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgClaimReward) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgClaimRewardResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgClaimRewardResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgClaimRewardResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgClaimRewardResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgClaimRewardResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgCreateFarm) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgCreateFarm) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCreateFarm) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCreateFarm) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCreateFarm) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgCreateFarmResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgCreateFarmResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCreateFarmResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCreateFarmResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCreateFarmResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgUpdateFarm) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgUpdateFarm) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgUpdateFarm) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgUpdateFarm) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgUpdateFarm) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgUpdateFarmResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgUpdateFarmResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgUpdateFarmResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgUpdateFarmResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgUpdateFarmResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }
