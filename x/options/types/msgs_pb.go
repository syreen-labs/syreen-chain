package types

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math/bits"

	proto "github.com/cosmos/gogoproto/proto"
)

var _ = proto.Marshal
var _ = fmt.Errorf
var _ = io.ErrUnexpectedEOF
var _ = binary.BigEndian
var _ = bits.Len64

var fileDescriptorOptionsTx []byte

// XXX methods
func (m *MsgWriteOption) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgWriteOption) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgWriteOption) XXX_Size() int                         { return m.Size() }

func (m *MsgWriteOptionResponse) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgWriteOptionResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgWriteOptionResponse) XXX_Size() int                         { return m.Size() }

func (m *MsgBuyOption) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgBuyOption) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgBuyOption) XXX_Size() int                         { return m.Size() }

func (m *MsgBuyOptionResponse) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgBuyOptionResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgBuyOptionResponse) XXX_Size() int                         { return m.Size() }

func (m *MsgExerciseOption) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgExerciseOption) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgExerciseOption) XXX_Size() int                         { return m.Size() }

func (m *MsgExerciseOptionResponse) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgExerciseOptionResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgExerciseOptionResponse) XXX_Size() int                         { return m.Size() }

func (m *MsgCancelOption) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgCancelOption) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCancelOption) XXX_Size() int                         { return m.Size() }

func (m *MsgCancelOptionResponse) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgCancelOptionResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCancelOptionResponse) XXX_Size() int                         { return m.Size() }

// Descriptor methods
func (*MsgWriteOption) Descriptor() ([]byte, []int)           { return fileDescriptorOptionsTx, []int{0} }
func (*MsgWriteOptionResponse) Descriptor() ([]byte, []int)   { return fileDescriptorOptionsTx, []int{1} }
func (*MsgBuyOption) Descriptor() ([]byte, []int)             { return fileDescriptorOptionsTx, []int{2} }
func (*MsgBuyOptionResponse) Descriptor() ([]byte, []int)     { return fileDescriptorOptionsTx, []int{3} }
func (*MsgExerciseOption) Descriptor() ([]byte, []int)        { return fileDescriptorOptionsTx, []int{4} }
func (*MsgExerciseOptionResponse) Descriptor() ([]byte, []int){ return fileDescriptorOptionsTx, []int{5} }
func (*MsgCancelOption) Descriptor() ([]byte, []int)          { return fileDescriptorOptionsTx, []int{6} }
func (*MsgCancelOptionResponse) Descriptor() ([]byte, []int)  { return fileDescriptorOptionsTx, []int{7} }

func init() {
	proto.RegisterType((*MsgWriteOption)(nil), "syreen.options.MsgWriteOption")
	proto.RegisterType((*MsgWriteOptionResponse)(nil), "syreen.options.MsgWriteOptionResponse")
	proto.RegisterType((*MsgBuyOption)(nil), "syreen.options.MsgBuyOption")
	proto.RegisterType((*MsgBuyOptionResponse)(nil), "syreen.options.MsgBuyOptionResponse")
	proto.RegisterType((*MsgExerciseOption)(nil), "syreen.options.MsgExerciseOption")
	proto.RegisterType((*MsgExerciseOptionResponse)(nil), "syreen.options.MsgExerciseOptionResponse")
	proto.RegisterType((*MsgCancelOption)(nil), "syreen.options.MsgCancelOption")
	proto.RegisterType((*MsgCancelOptionResponse)(nil), "syreen.options.MsgCancelOptionResponse")
}

// Marshal/Unmarshal/Size
func (m *MsgWriteOption) Marshal() ([]byte, error)                  { return json.Marshal(m) }
func (m *MsgWriteOption) MarshalTo(dAtA []byte) (int, error)        { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgWriteOption) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgWriteOption) Size() int                                 { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgWriteOption) Unmarshal(bz []byte) error                 { return json.Unmarshal(bz, m) }

func (m *MsgWriteOptionResponse) Marshal() ([]byte, error)                  { return json.Marshal(m) }
func (m *MsgWriteOptionResponse) MarshalTo(dAtA []byte) (int, error)        { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgWriteOptionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgWriteOptionResponse) Size() int                                 { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgWriteOptionResponse) Unmarshal(bz []byte) error                 { return json.Unmarshal(bz, m) }

func (m *MsgBuyOption) Marshal() ([]byte, error)                  { return json.Marshal(m) }
func (m *MsgBuyOption) MarshalTo(dAtA []byte) (int, error)        { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgBuyOption) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgBuyOption) Size() int                                 { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgBuyOption) Unmarshal(bz []byte) error                 { return json.Unmarshal(bz, m) }

func (m *MsgBuyOptionResponse) Marshal() ([]byte, error)                  { return json.Marshal(m) }
func (m *MsgBuyOptionResponse) MarshalTo(dAtA []byte) (int, error)        { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgBuyOptionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgBuyOptionResponse) Size() int                                 { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgBuyOptionResponse) Unmarshal(bz []byte) error                 { return json.Unmarshal(bz, m) }

func (m *MsgExerciseOption) Marshal() ([]byte, error)                  { return json.Marshal(m) }
func (m *MsgExerciseOption) MarshalTo(dAtA []byte) (int, error)        { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgExerciseOption) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgExerciseOption) Size() int                                 { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgExerciseOption) Unmarshal(bz []byte) error                 { return json.Unmarshal(bz, m) }

func (m *MsgExerciseOptionResponse) Marshal() ([]byte, error)                  { return json.Marshal(m) }
func (m *MsgExerciseOptionResponse) MarshalTo(dAtA []byte) (int, error)        { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgExerciseOptionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgExerciseOptionResponse) Size() int                                 { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgExerciseOptionResponse) Unmarshal(bz []byte) error                 { return json.Unmarshal(bz, m) }

func (m *MsgCancelOption) Marshal() ([]byte, error)                  { return json.Marshal(m) }
func (m *MsgCancelOption) MarshalTo(dAtA []byte) (int, error)        { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCancelOption) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCancelOption) Size() int                                 { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCancelOption) Unmarshal(bz []byte) error                 { return json.Unmarshal(bz, m) }

func (m *MsgCancelOptionResponse) Marshal() ([]byte, error)                  { return json.Marshal(m) }
func (m *MsgCancelOptionResponse) MarshalTo(dAtA []byte) (int, error)        { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCancelOptionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCancelOptionResponse) Size() int                                 { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCancelOptionResponse) Unmarshal(bz []byte) error                 { return json.Unmarshal(bz, m) }
