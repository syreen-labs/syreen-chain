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

var fileDescriptorPredictTx []byte

// XXX methods required by gRPC decoder — MsgCreateMarket
func (m *MsgCreateMarket) XXX_Unmarshal(b []byte) error                { return m.Unmarshal(b) }
func (m *MsgCreateMarket) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCreateMarket) XXX_Size() int                               { return m.Size() }

func (m *MsgCreateMarketResponse) XXX_Unmarshal(b []byte) error                { return m.Unmarshal(b) }
func (m *MsgCreateMarketResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgCreateMarketResponse) XXX_Size() int                               { return m.Size() }

// MsgBuyShares
func (m *MsgBuyShares) XXX_Unmarshal(b []byte) error                { return m.Unmarshal(b) }
func (m *MsgBuyShares) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgBuyShares) XXX_Size() int                               { return m.Size() }

func (m *MsgBuySharesResponse) XXX_Unmarshal(b []byte) error                { return m.Unmarshal(b) }
func (m *MsgBuySharesResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgBuySharesResponse) XXX_Size() int                               { return m.Size() }

// MsgSellShares
func (m *MsgSellShares) XXX_Unmarshal(b []byte) error                { return m.Unmarshal(b) }
func (m *MsgSellShares) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgSellShares) XXX_Size() int                               { return m.Size() }

func (m *MsgSellSharesResponse) XXX_Unmarshal(b []byte) error                { return m.Unmarshal(b) }
func (m *MsgSellSharesResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgSellSharesResponse) XXX_Size() int                               { return m.Size() }

// MsgResolveMarket
func (m *MsgResolveMarket) XXX_Unmarshal(b []byte) error                { return m.Unmarshal(b) }
func (m *MsgResolveMarket) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgResolveMarket) XXX_Size() int                               { return m.Size() }

func (m *MsgResolveMarketResponse) XXX_Unmarshal(b []byte) error                { return m.Unmarshal(b) }
func (m *MsgResolveMarketResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgResolveMarketResponse) XXX_Size() int                               { return m.Size() }

// MsgClaimWinnings
func (m *MsgClaimWinnings) XXX_Unmarshal(b []byte) error                { return m.Unmarshal(b) }
func (m *MsgClaimWinnings) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgClaimWinnings) XXX_Size() int                               { return m.Size() }

func (m *MsgClaimWinningsResponse) XXX_Unmarshal(b []byte) error                { return m.Unmarshal(b) }
func (m *MsgClaimWinningsResponse) XXX_Marshal(b []byte, _ bool) ([]byte, error) { return m.Marshal() }
func (m *MsgClaimWinningsResponse) XXX_Size() int                               { return m.Size() }

// Descriptor methods
func (*MsgCreateMarket) Descriptor() ([]byte, []int)           { return fileDescriptorPredictTx, []int{0} }
func (*MsgCreateMarketResponse) Descriptor() ([]byte, []int)   { return fileDescriptorPredictTx, []int{1} }
func (*MsgBuyShares) Descriptor() ([]byte, []int)              { return fileDescriptorPredictTx, []int{2} }
func (*MsgBuySharesResponse) Descriptor() ([]byte, []int)      { return fileDescriptorPredictTx, []int{3} }
func (*MsgSellShares) Descriptor() ([]byte, []int)             { return fileDescriptorPredictTx, []int{4} }
func (*MsgSellSharesResponse) Descriptor() ([]byte, []int)     { return fileDescriptorPredictTx, []int{5} }
func (*MsgResolveMarket) Descriptor() ([]byte, []int)          { return fileDescriptorPredictTx, []int{6} }
func (*MsgResolveMarketResponse) Descriptor() ([]byte, []int)  { return fileDescriptorPredictTx, []int{7} }
func (*MsgClaimWinnings) Descriptor() ([]byte, []int)          { return fileDescriptorPredictTx, []int{8} }
func (*MsgClaimWinningsResponse) Descriptor() ([]byte, []int)  { return fileDescriptorPredictTx, []int{9} }

func init() {
	proto.RegisterType((*MsgCreateMarket)(nil), "syreen.predict.MsgCreateMarket")
	proto.RegisterType((*MsgCreateMarketResponse)(nil), "syreen.predict.MsgCreateMarketResponse")
	proto.RegisterType((*MsgBuyShares)(nil), "syreen.predict.MsgBuyShares")
	proto.RegisterType((*MsgBuySharesResponse)(nil), "syreen.predict.MsgBuySharesResponse")
	proto.RegisterType((*MsgSellShares)(nil), "syreen.predict.MsgSellShares")
	proto.RegisterType((*MsgSellSharesResponse)(nil), "syreen.predict.MsgSellSharesResponse")
	proto.RegisterType((*MsgResolveMarket)(nil), "syreen.predict.MsgResolveMarket")
	proto.RegisterType((*MsgResolveMarketResponse)(nil), "syreen.predict.MsgResolveMarketResponse")
	proto.RegisterType((*MsgClaimWinnings)(nil), "syreen.predict.MsgClaimWinnings")
	proto.RegisterType((*MsgClaimWinningsResponse)(nil), "syreen.predict.MsgClaimWinningsResponse")
}

// Marshal/Unmarshal/Size for all message types using JSON

func (m *MsgCreateMarket) Marshal() ([]byte, error)                            { return json.Marshal(m) }
func (m *MsgCreateMarket) MarshalTo(dAtA []byte) (int, error)                  { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCreateMarket) MarshalToSizedBuffer(dAtA []byte) (int, error)       { return m.MarshalTo(dAtA) }
func (m *MsgCreateMarket) Size() int                                            { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCreateMarket) Unmarshal(bz []byte) error                           { return json.Unmarshal(bz, m) }

func (m *MsgCreateMarketResponse) Marshal() ([]byte, error)                    { return json.Marshal(m) }
func (m *MsgCreateMarketResponse) MarshalTo(dAtA []byte) (int, error)          { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgCreateMarketResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgCreateMarketResponse) Size() int                                    { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgCreateMarketResponse) Unmarshal(bz []byte) error                   { return json.Unmarshal(bz, m) }

func (m *MsgBuyShares) Marshal() ([]byte, error)                               { return json.Marshal(m) }
func (m *MsgBuyShares) MarshalTo(dAtA []byte) (int, error)                     { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgBuyShares) MarshalToSizedBuffer(dAtA []byte) (int, error)          { return m.MarshalTo(dAtA) }
func (m *MsgBuyShares) Size() int                                               { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgBuyShares) Unmarshal(bz []byte) error                              { return json.Unmarshal(bz, m) }

func (m *MsgBuySharesResponse) Marshal() ([]byte, error)                       { return json.Marshal(m) }
func (m *MsgBuySharesResponse) MarshalTo(dAtA []byte) (int, error)             { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgBuySharesResponse) MarshalToSizedBuffer(dAtA []byte) (int, error)  { return m.MarshalTo(dAtA) }
func (m *MsgBuySharesResponse) Size() int                                       { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgBuySharesResponse) Unmarshal(bz []byte) error                      { return json.Unmarshal(bz, m) }

func (m *MsgSellShares) Marshal() ([]byte, error)                              { return json.Marshal(m) }
func (m *MsgSellShares) MarshalTo(dAtA []byte) (int, error)                    { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgSellShares) MarshalToSizedBuffer(dAtA []byte) (int, error)         { return m.MarshalTo(dAtA) }
func (m *MsgSellShares) Size() int                                              { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgSellShares) Unmarshal(bz []byte) error                             { return json.Unmarshal(bz, m) }

func (m *MsgSellSharesResponse) Marshal() ([]byte, error)                      { return json.Marshal(m) }
func (m *MsgSellSharesResponse) MarshalTo(dAtA []byte) (int, error)            { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgSellSharesResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgSellSharesResponse) Size() int                                      { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgSellSharesResponse) Unmarshal(bz []byte) error                     { return json.Unmarshal(bz, m) }

func (m *MsgResolveMarket) Marshal() ([]byte, error)                           { return json.Marshal(m) }
func (m *MsgResolveMarket) MarshalTo(dAtA []byte) (int, error)                 { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgResolveMarket) MarshalToSizedBuffer(dAtA []byte) (int, error)      { return m.MarshalTo(dAtA) }
func (m *MsgResolveMarket) Size() int                                           { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgResolveMarket) Unmarshal(bz []byte) error                          { return json.Unmarshal(bz, m) }

func (m *MsgResolveMarketResponse) Marshal() ([]byte, error)                   { return json.Marshal(m) }
func (m *MsgResolveMarketResponse) MarshalTo(dAtA []byte) (int, error)         { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgResolveMarketResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgResolveMarketResponse) Size() int                                   { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgResolveMarketResponse) Unmarshal(bz []byte) error                  { return json.Unmarshal(bz, m) }

func (m *MsgClaimWinnings) Marshal() ([]byte, error)                           { return json.Marshal(m) }
func (m *MsgClaimWinnings) MarshalTo(dAtA []byte) (int, error)                 { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgClaimWinnings) MarshalToSizedBuffer(dAtA []byte) (int, error)      { return m.MarshalTo(dAtA) }
func (m *MsgClaimWinnings) Size() int                                           { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgClaimWinnings) Unmarshal(bz []byte) error                          { return json.Unmarshal(bz, m) }

func (m *MsgClaimWinningsResponse) Marshal() ([]byte, error)                   { return json.Marshal(m) }
func (m *MsgClaimWinningsResponse) MarshalTo(dAtA []byte) (int, error)         { bz, err := json.Marshal(m); if err != nil { return 0, err }; copy(dAtA, bz); return len(bz), nil }
func (m *MsgClaimWinningsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return m.MarshalTo(dAtA) }
func (m *MsgClaimWinningsResponse) Size() int                                   { bz, _ := json.Marshal(m); return len(bz) }
func (m *MsgClaimWinningsResponse) Unmarshal(bz []byte) error                  { return json.Unmarshal(bz, m) }
