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

var fileDescriptorCLMMTx []byte

// XXX methods required by gRPC decoder — MsgCreateCLPool
func (m *MsgCreateCLPool) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgCreateCLPool) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCreateCLPool) XXX_Size() int                         { return m.Size() }

func (m *MsgCreateCLPoolResponse) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgCreateCLPoolResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCreateCLPoolResponse) XXX_Size() int                         { return m.Size() }

// MsgCreatePosition
func (m *MsgCreatePosition) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgCreatePosition) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCreatePosition) XXX_Size() int                         { return m.Size() }

func (m *MsgCreatePositionResponse) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgCreatePositionResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCreatePositionResponse) XXX_Size() int                         { return m.Size() }

// MsgAddLiquidity
func (m *MsgAddLiquidity) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgAddLiquidity) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgAddLiquidity) XXX_Size() int                         { return m.Size() }

func (m *MsgAddLiquidityResponse) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgAddLiquidityResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgAddLiquidityResponse) XXX_Size() int                         { return m.Size() }

// MsgRemoveLiquidity
func (m *MsgRemoveLiquidity) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgRemoveLiquidity) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgRemoveLiquidity) XXX_Size() int                         { return m.Size() }

func (m *MsgRemoveLiquidityResponse) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgRemoveLiquidityResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgRemoveLiquidityResponse) XXX_Size() int                         { return m.Size() }

// MsgCollectFees
func (m *MsgCollectFees) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgCollectFees) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCollectFees) XXX_Size() int                         { return m.Size() }

func (m *MsgCollectFeesResponse) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgCollectFeesResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCollectFeesResponse) XXX_Size() int                         { return m.Size() }

// MsgCLSwap
func (m *MsgCLSwap) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgCLSwap) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCLSwap) XXX_Size() int                         { return m.Size() }

func (m *MsgCLSwapResponse) XXX_Unmarshal(b []byte) error          { return m.Unmarshal(b) }
func (m *MsgCLSwapResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCLSwapResponse) XXX_Size() int                         { return m.Size() }

// Descriptor methods
func (*MsgCreateCLPool) Descriptor() ([]byte, []int)             { return fileDescriptorCLMMTx, []int{0} }
func (*MsgCreateCLPoolResponse) Descriptor() ([]byte, []int)     { return fileDescriptorCLMMTx, []int{1} }
func (*MsgCreatePosition) Descriptor() ([]byte, []int)           { return fileDescriptorCLMMTx, []int{2} }
func (*MsgCreatePositionResponse) Descriptor() ([]byte, []int)   { return fileDescriptorCLMMTx, []int{3} }
func (*MsgAddLiquidity) Descriptor() ([]byte, []int)             { return fileDescriptorCLMMTx, []int{4} }
func (*MsgAddLiquidityResponse) Descriptor() ([]byte, []int)     { return fileDescriptorCLMMTx, []int{5} }
func (*MsgRemoveLiquidity) Descriptor() ([]byte, []int)          { return fileDescriptorCLMMTx, []int{6} }
func (*MsgRemoveLiquidityResponse) Descriptor() ([]byte, []int)  { return fileDescriptorCLMMTx, []int{7} }
func (*MsgCollectFees) Descriptor() ([]byte, []int)              { return fileDescriptorCLMMTx, []int{8} }
func (*MsgCollectFeesResponse) Descriptor() ([]byte, []int)      { return fileDescriptorCLMMTx, []int{9} }
func (*MsgCLSwap) Descriptor() ([]byte, []int)                   { return fileDescriptorCLMMTx, []int{10} }
func (*MsgCLSwapResponse) Descriptor() ([]byte, []int)           { return fileDescriptorCLMMTx, []int{11} }

func init() {
	proto.RegisterType((*MsgCreateCLPool)(nil), "syreen.clmm.MsgCreateCLPool")
	proto.RegisterType((*MsgCreateCLPoolResponse)(nil), "syreen.clmm.MsgCreateCLPoolResponse")
	proto.RegisterType((*MsgCreatePosition)(nil), "syreen.clmm.MsgCreatePosition")
	proto.RegisterType((*MsgCreatePositionResponse)(nil), "syreen.clmm.MsgCreatePositionResponse")
	proto.RegisterType((*MsgAddLiquidity)(nil), "syreen.clmm.MsgAddLiquidity")
	proto.RegisterType((*MsgAddLiquidityResponse)(nil), "syreen.clmm.MsgAddLiquidityResponse")
	proto.RegisterType((*MsgRemoveLiquidity)(nil), "syreen.clmm.MsgRemoveLiquidity")
	proto.RegisterType((*MsgRemoveLiquidityResponse)(nil), "syreen.clmm.MsgRemoveLiquidityResponse")
	proto.RegisterType((*MsgCollectFees)(nil), "syreen.clmm.MsgCollectFees")
	proto.RegisterType((*MsgCollectFeesResponse)(nil), "syreen.clmm.MsgCollectFeesResponse")
	proto.RegisterType((*MsgCLSwap)(nil), "syreen.clmm.MsgCLSwap")
	proto.RegisterType((*MsgCLSwapResponse)(nil), "syreen.clmm.MsgCLSwapResponse")
}

// Marshal/Unmarshal/Size for all message types using JSON

func (m *MsgCreateCLPool) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgCreateCLPool) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCreateCLPool) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCreateCLPool) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCreateCLPool) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgCreateCLPoolResponse) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgCreateCLPoolResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCreateCLPoolResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCreateCLPoolResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCreateCLPoolResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgCreatePosition) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgCreatePosition) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCreatePosition) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCreatePosition) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCreatePosition) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgCreatePositionResponse) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgCreatePositionResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCreatePositionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCreatePositionResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCreatePositionResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgAddLiquidity) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgAddLiquidity) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgAddLiquidity) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgAddLiquidity) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgAddLiquidity) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgAddLiquidityResponse) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgAddLiquidityResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgAddLiquidityResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgAddLiquidityResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgAddLiquidityResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgRemoveLiquidity) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgRemoveLiquidity) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgRemoveLiquidity) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgRemoveLiquidity) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgRemoveLiquidity) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgRemoveLiquidityResponse) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgRemoveLiquidityResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgRemoveLiquidityResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgRemoveLiquidityResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgRemoveLiquidityResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgCollectFees) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgCollectFees) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCollectFees) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCollectFees) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCollectFees) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgCollectFeesResponse) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgCollectFeesResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCollectFeesResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCollectFeesResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCollectFeesResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgCLSwap) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgCLSwap) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCLSwap) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCLSwap) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCLSwap) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }

func (m *MsgCLSwapResponse) Marshal() ([]byte, error) { return json.Marshal(m) }
func (m *MsgCLSwapResponse) MarshalTo(dAtA []byte) (int, error) { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCLSwapResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCLSwapResponse) Size() int { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCLSwapResponse) Unmarshal(bz []byte) error { return json.Unmarshal(bz, m) }
