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

var fileDescriptorPerpsTx []byte

// XXX methods required by gRPC decoder
func (m *MsgOpenPosition) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgOpenPosition) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgOpenPosition) XXX_Size() int { return m.Size() }

func (m *MsgOpenPositionResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgOpenPositionResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgOpenPositionResponse) XXX_Size() int { return m.Size() }

func (m *MsgClosePosition) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgClosePosition) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgClosePosition) XXX_Size() int { return m.Size() }

func (m *MsgClosePositionResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgClosePositionResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgClosePositionResponse) XXX_Size() int { return m.Size() }

func (m *MsgAddMargin) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgAddMargin) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgAddMargin) XXX_Size() int { return m.Size() }

func (m *MsgAddMarginResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgAddMarginResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgAddMarginResponse) XXX_Size() int { return m.Size() }

func (m *MsgRemoveMargin) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRemoveMargin) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgRemoveMargin) XXX_Size() int { return m.Size() }

func (m *MsgRemoveMarginResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRemoveMarginResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgRemoveMarginResponse) XXX_Size() int { return m.Size() }

func (m *MsgCreateMarket) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateMarket) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCreateMarket) XXX_Size() int { return m.Size() }

func (m *MsgCreateMarketResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateMarketResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCreateMarketResponse) XXX_Size() int { return m.Size() }

// Descriptor methods
func (*MsgOpenPosition) Descriptor() ([]byte, []int)          { return fileDescriptorPerpsTx, []int{0} }
func (*MsgOpenPositionResponse) Descriptor() ([]byte, []int)  { return fileDescriptorPerpsTx, []int{1} }
func (*MsgClosePosition) Descriptor() ([]byte, []int)         { return fileDescriptorPerpsTx, []int{2} }
func (*MsgClosePositionResponse) Descriptor() ([]byte, []int) { return fileDescriptorPerpsTx, []int{3} }
func (*MsgAddMargin) Descriptor() ([]byte, []int)             { return fileDescriptorPerpsTx, []int{4} }
func (*MsgAddMarginResponse) Descriptor() ([]byte, []int)     { return fileDescriptorPerpsTx, []int{5} }
func (*MsgRemoveMargin) Descriptor() ([]byte, []int)          { return fileDescriptorPerpsTx, []int{6} }
func (*MsgRemoveMarginResponse) Descriptor() ([]byte, []int)  { return fileDescriptorPerpsTx, []int{7} }
func (*MsgCreateMarket) Descriptor() ([]byte, []int)          { return fileDescriptorPerpsTx, []int{8} }
func (*MsgCreateMarketResponse) Descriptor() ([]byte, []int)  { return fileDescriptorPerpsTx, []int{9} }

func init() {
	proto.RegisterType((*MsgOpenPosition)(nil), "syreen.perps.MsgOpenPosition")
	proto.RegisterType((*MsgOpenPositionResponse)(nil), "syreen.perps.MsgOpenPositionResponse")
	proto.RegisterType((*MsgClosePosition)(nil), "syreen.perps.MsgClosePosition")
	proto.RegisterType((*MsgClosePositionResponse)(nil), "syreen.perps.MsgClosePositionResponse")
	proto.RegisterType((*MsgAddMargin)(nil), "syreen.perps.MsgAddMargin")
	proto.RegisterType((*MsgAddMarginResponse)(nil), "syreen.perps.MsgAddMarginResponse")
	proto.RegisterType((*MsgRemoveMargin)(nil), "syreen.perps.MsgRemoveMargin")
	proto.RegisterType((*MsgRemoveMarginResponse)(nil), "syreen.perps.MsgRemoveMarginResponse")
	proto.RegisterType((*MsgCreateMarket)(nil), "syreen.perps.MsgCreateMarket")
	proto.RegisterType((*MsgCreateMarketResponse)(nil), "syreen.perps.MsgCreateMarketResponse")
}

// Marshal/Unmarshal/Size for all message types using JSON
func (m *MsgOpenPosition) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgOpenPosition) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgOpenPosition) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgOpenPosition) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgOpenPosition) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgOpenPositionResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgOpenPositionResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgOpenPositionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgOpenPositionResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgOpenPositionResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgClosePosition) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgClosePosition) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgClosePosition) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgClosePosition) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgClosePosition) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgClosePositionResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgClosePositionResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgClosePositionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgClosePositionResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgClosePositionResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgAddMargin) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgAddMargin) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgAddMargin) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgAddMargin) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgAddMargin) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgAddMarginResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgAddMarginResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgAddMarginResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgAddMarginResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgAddMarginResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgRemoveMargin) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgRemoveMargin) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgRemoveMargin) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgRemoveMargin) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgRemoveMargin) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgRemoveMarginResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgRemoveMarginResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgRemoveMarginResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgRemoveMarginResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgRemoveMarginResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgCreateMarket) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgCreateMarket) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCreateMarket) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCreateMarket) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCreateMarket) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgCreateMarketResponse) Marshal() ([]byte, error)    { return json.Marshal(m) }
func (m *MsgCreateMarketResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCreateMarketResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCreateMarketResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCreateMarketResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }
