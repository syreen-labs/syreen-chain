package types

import (
	"fmt"
	"io"
	"math/bits"

	proto "github.com/cosmos/gogoproto/proto"
)

var _ = proto.Marshal
var _ = fmt.Errorf
var _ = io.ErrUnexpectedEOF

var fileDescriptorPredictTx []byte

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

// ---------------------------------------------------------------------------
// MsgCreateMarket — creator(s,1) question(s,2) resolver(s,3) quote_denom(s,4) resolution_block(varint,5) initial_liquidity(bytes,6)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCreateMarket proto.InternalMessageInfo
func (m *MsgCreateMarket) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateMarket) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgCreateMarket.Marshal(b, m, deterministic) }
	b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil
}
func (m *MsgCreateMarket) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCreateMarket.Merge(m, src) }
func (m *MsgCreateMarket) XXX_Size() int               { return m.Size() }
func (m *MsgCreateMarket) XXX_DiscardUnknown()         { xxx_messageInfo_MsgCreateMarket.DiscardUnknown(m) }

func (m *MsgCreateMarket) Marshal() (dAtA []byte, err error) {
	size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil
}
func (m *MsgCreateMarket) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgCreateMarket) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	// f6: initial_liquidity (bytes)
	{ bz, err := m.InitialLiquidity.Marshal(); if err != nil { return 0, err }; if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintPredict(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x32 } }
	// f5: resolution_block (varint)
	if m.ResolutionBlock != 0 { i = encodeVarintPredict(dAtA, i, uint64(m.ResolutionBlock)); i--; dAtA[i] = 0x28 }
	// f4: quote_denom (string)
	if len(m.QuoteDenom) > 0 { i -= len(m.QuoteDenom); copy(dAtA[i:], m.QuoteDenom); i = encodeVarintPredict(dAtA, i, uint64(len(m.QuoteDenom))); i--; dAtA[i] = 0x22 }
	// f3: resolver (string)
	if len(m.Resolver) > 0 { i -= len(m.Resolver); copy(dAtA[i:], m.Resolver); i = encodeVarintPredict(dAtA, i, uint64(len(m.Resolver))); i--; dAtA[i] = 0x1a }
	// f2: question (string)
	if len(m.Question) > 0 { i -= len(m.Question); copy(dAtA[i:], m.Question); i = encodeVarintPredict(dAtA, i, uint64(len(m.Question))); i--; dAtA[i] = 0x12 }
	// f1: creator (string)
	if len(m.Creator) > 0 { i -= len(m.Creator); copy(dAtA[i:], m.Creator); i = encodeVarintPredict(dAtA, i, uint64(len(m.Creator))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgCreateMarket) Size() (n int) {
	if m == nil { return 0 }
	var l int
	l = len(m.Creator); if l > 0 { n += 1 + l + sovPredict(uint64(l)) }
	l = len(m.Question); if l > 0 { n += 1 + l + sovPredict(uint64(l)) }
	l = len(m.Resolver); if l > 0 { n += 1 + l + sovPredict(uint64(l)) }
	l = len(m.QuoteDenom); if l > 0 { n += 1 + l + sovPredict(uint64(l)) }
	if m.ResolutionBlock != 0 { n += 1 + sovPredict(uint64(m.ResolutionBlock)) }
	bz, _ := m.InitialLiquidity.Marshal(); l = len(bz); if l > 0 { n += 1 + l + sovPredict(uint64(l)) }
	return n
}
func (m *MsgCreateMarket) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7)
		if wireType == 4 { return fmt.Errorf("proto: MsgCreateMarket: wiretype end group") }
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgCreateMarket: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: // creator
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
			intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthPredict }
			postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthPredict }; if postIndex > l { return io.ErrUnexpectedEOF }
			m.Creator = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2: // question
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Question", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
			intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthPredict }
			postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthPredict }; if postIndex > l { return io.ErrUnexpectedEOF }
			m.Question = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 3: // resolver
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Resolver", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
			intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthPredict }
			postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthPredict }; if postIndex > l { return io.ErrUnexpectedEOF }
			m.Resolver = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 4: // quote_denom
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field QuoteDenom", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
			intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthPredict }
			postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthPredict }; if postIndex > l { return io.ErrUnexpectedEOF }
			m.QuoteDenom = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 5: // resolution_block (varint)
			if wireType != 0 { return fmt.Errorf("proto: wrong wireType = %d for field ResolutionBlock", wireType) }
			m.ResolutionBlock = 0
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.ResolutionBlock |= int64(b&0x7F) << shift; if b < 0x80 { break } }
		case 6: // initial_liquidity (bytes)
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field InitialLiquidity", wireType) }
			var byteLen int
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }
			if byteLen < 0 { return ErrInvalidLengthPredict }
			postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthPredict }; if postIndex > l { return io.ErrUnexpectedEOF }
			if err := m.InitialLiquidity.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex
		default:
			iNdEx = preIndex; skippy, err := skipPredict(dAtA[iNdEx:]); if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthPredict }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// ---------------------------------------------------------------------------
// MsgCreateMarketResponse — market_id(varint,1)
// ---------------------------------------------------------------------------
var xxx_messageInfo_MsgCreateMarketResponse proto.InternalMessageInfo
func (m *MsgCreateMarketResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateMarketResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgCreateMarketResponse.Marshal(b, m, deterministic) }
	b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil
}
func (m *MsgCreateMarketResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCreateMarketResponse.Merge(m, src) }
func (m *MsgCreateMarketResponse) XXX_Size() int               { return m.Size() }
func (m *MsgCreateMarketResponse) XXX_DiscardUnknown()         { xxx_messageInfo_MsgCreateMarketResponse.DiscardUnknown(m) }

func (m *MsgCreateMarketResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil
}
func (m *MsgCreateMarketResponse) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgCreateMarketResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.MarketID != 0 { i = encodeVarintPredict(dAtA, i, uint64(m.MarketID)); i--; dAtA[i] = 0x08 }
	return len(dAtA) - i, nil
}
func (m *MsgCreateMarketResponse) Size() (n int) {
	if m == nil { return 0 }; if m.MarketID != 0 { n += 1 + sovPredict(uint64(m.MarketID)) }; return n
}
func (m *MsgCreateMarketResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3)
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgCreateMarketResponse: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1:
			if int(wire&0x7) != 0 { return fmt.Errorf("proto: wrong wireType for MarketID") }
			m.MarketID = 0
			for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.MarketID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		default:
			iNdEx = preIndex; skippy, err := skipPredict(dAtA[iNdEx:]); if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthPredict }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// ---------------------------------------------------------------------------
// MsgBuyShares — sender(s,1) market_id(varint,2) outcome(s,3) amount(bytes,4)
// ---------------------------------------------------------------------------
var xxx_messageInfo_MsgBuyShares proto.InternalMessageInfo
func (m *MsgBuyShares) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgBuyShares) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgBuyShares.Marshal(b, m, deterministic) }
	b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil
}
func (m *MsgBuyShares) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgBuyShares.Merge(m, src) }
func (m *MsgBuyShares) XXX_Size() int               { return m.Size() }
func (m *MsgBuyShares) XXX_DiscardUnknown()         { xxx_messageInfo_MsgBuyShares.DiscardUnknown(m) }

func (m *MsgBuyShares) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgBuyShares) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgBuyShares) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	{ bz, err := m.Amount.Marshal(); if err != nil { return 0, err }; if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintPredict(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x22 } }
	if len(m.Outcome) > 0 { i -= len(m.Outcome); copy(dAtA[i:], m.Outcome); i = encodeVarintPredict(dAtA, i, uint64(len(m.Outcome))); i--; dAtA[i] = 0x1a }
	if m.MarketID != 0 { i = encodeVarintPredict(dAtA, i, uint64(m.MarketID)); i--; dAtA[i] = 0x10 }
	if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarintPredict(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgBuyShares) Size() (n int) {
	if m == nil { return 0 }
	l := len(m.Sender); if l > 0 { n += 1 + l + sovPredict(uint64(l)) }
	if m.MarketID != 0 { n += 1 + sovPredict(uint64(m.MarketID)) }
	l = len(m.Outcome); if l > 0 { n += 1 + l + sovPredict(uint64(l)) }
	bz, _ := m.Amount.Marshal(); l = len(bz); if l > 0 { n += 1 + l + sovPredict(uint64(l)) }
	return n
}
func (m *MsgBuyShares) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7)
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgBuyShares: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthPredict }; postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthPredict }; if postIndex > l { return io.ErrUnexpectedEOF }; m.Sender = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.MarketID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.MarketID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		case 3: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthPredict }; postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthPredict }; if postIndex > l { return io.ErrUnexpectedEOF }; m.Outcome = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 4: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLengthPredict }; postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthPredict }; if postIndex > l { return io.ErrUnexpectedEOF }; if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex
		default: iNdEx = preIndex; skippy, err := skipPredict(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthPredict }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// ---------------------------------------------------------------------------
// MsgBuySharesResponse — shares_bought(bytes,1) avg_price(bytes,2)
// ---------------------------------------------------------------------------
var xxx_messageInfo_MsgBuySharesResponse proto.InternalMessageInfo
func (m *MsgBuySharesResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgBuySharesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgBuySharesResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgBuySharesResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgBuySharesResponse.Merge(m, src) }
func (m *MsgBuySharesResponse) XXX_Size() int { return m.Size() }
func (m *MsgBuySharesResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgBuySharesResponse.DiscardUnknown(m) }
func (m *MsgBuySharesResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgBuySharesResponse) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgBuySharesResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	{ bz, err := m.AvgPrice.Marshal(); if err != nil { return 0, err }; if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintPredict(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x12 } }
	{ bz, err := m.SharesBought.Marshal(); if err != nil { return 0, err }; if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintPredict(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x0a } }
	return len(dAtA) - i, nil
}
func (m *MsgBuySharesResponse) Size() (n int) {
	if m == nil { return 0 }
	bz, _ := m.SharesBought.Marshal(); l := len(bz); if l > 0 { n += 1 + l + sovPredict(uint64(l)) }
	bz2, _ := m.AvgPrice.Marshal(); l = len(bz2); if l > 0 { n += 1 + l + sovPredict(uint64(l)) }
	return n
}
func (m *MsgBuySharesResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3)
		if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLengthPredict }; postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthPredict }; if postIndex > l { return io.ErrUnexpectedEOF }; if err := m.SharesBought.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex
		case 2: if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLengthPredict }; postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthPredict }; if postIndex > l { return io.ErrUnexpectedEOF }; if err := m.AvgPrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex
		default: iNdEx = preIndex; skippy, err := skipPredict(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthPredict }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// ---------------------------------------------------------------------------
// MsgSellShares — same as BuyShares but field 4 is shares(bytes)
// ---------------------------------------------------------------------------
var xxx_messageInfo_MsgSellShares proto.InternalMessageInfo
func (m *MsgSellShares) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgSellShares) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgSellShares.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgSellShares) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgSellShares.Merge(m, src) }
func (m *MsgSellShares) XXX_Size() int { return m.Size() }
func (m *MsgSellShares) XXX_DiscardUnknown() { xxx_messageInfo_MsgSellShares.DiscardUnknown(m) }
func (m *MsgSellShares) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgSellShares) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgSellShares) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	{ bz, err := m.Shares.Marshal(); if err != nil { return 0, err }; if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintPredict(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x22 } }
	if len(m.Outcome) > 0 { i -= len(m.Outcome); copy(dAtA[i:], m.Outcome); i = encodeVarintPredict(dAtA, i, uint64(len(m.Outcome))); i--; dAtA[i] = 0x1a }
	if m.MarketID != 0 { i = encodeVarintPredict(dAtA, i, uint64(m.MarketID)); i--; dAtA[i] = 0x10 }
	if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarintPredict(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgSellShares) Size() (n int) {
	if m == nil { return 0 }
	l := len(m.Sender); if l > 0 { n += 1 + l + sovPredict(uint64(l)) }
	if m.MarketID != 0 { n += 1 + sovPredict(uint64(m.MarketID)) }
	l = len(m.Outcome); if l > 0 { n += 1 + l + sovPredict(uint64(l)) }
	bz, _ := m.Shares.Marshal(); l = len(bz); if l > 0 { n += 1 + l + sovPredict(uint64(l)) }
	return n
}
func (m *MsgSellShares) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7)
		if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthPredict }; postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthPredict }; if postIndex > l { return io.ErrUnexpectedEOF }; m.Sender = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.MarketID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.MarketID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		case 3: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthPredict }; postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthPredict }; if postIndex > l { return io.ErrUnexpectedEOF }; m.Outcome = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 4: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLengthPredict }; postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthPredict }; if postIndex > l { return io.ErrUnexpectedEOF }; if err := m.Shares.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex
		default: iNdEx = preIndex; skippy, err := skipPredict(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthPredict }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgSellSharesResponse — quote_returned(bytes,1)
var xxx_messageInfo_MsgSellSharesResponse proto.InternalMessageInfo
func (m *MsgSellSharesResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgSellSharesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgSellSharesResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgSellSharesResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgSellSharesResponse.Merge(m, src) }
func (m *MsgSellSharesResponse) XXX_Size() int { return m.Size() }
func (m *MsgSellSharesResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgSellSharesResponse.DiscardUnknown(m) }
func (m *MsgSellSharesResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgSellSharesResponse) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgSellSharesResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	bz, err := m.QuoteReturned.Marshal(); if err != nil { return 0, err }; if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintPredict(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgSellSharesResponse) Size() (n int) { if m == nil { return 0 }; bz, _ := m.QuoteReturned.Marshal(); l := len(bz); if l > 0 { n += 1 + l + sovPredict(uint64(l)) }; return n }
func (m *MsgSellSharesResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3)
		if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLengthPredict }; postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthPredict }; if postIndex > l { return io.ErrUnexpectedEOF }; if err := m.QuoteReturned.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex
		default: iNdEx = preIndex; skippy, err := skipPredict(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthPredict }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgResolveMarket — resolver(s,1) market_id(varint,2) outcome(s,3)
var xxx_messageInfo_MsgResolveMarket proto.InternalMessageInfo
func (m *MsgResolveMarket) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgResolveMarket) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgResolveMarket.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgResolveMarket) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgResolveMarket.Merge(m, src) }
func (m *MsgResolveMarket) XXX_Size() int { return m.Size() }
func (m *MsgResolveMarket) XXX_DiscardUnknown() { xxx_messageInfo_MsgResolveMarket.DiscardUnknown(m) }
func (m *MsgResolveMarket) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgResolveMarket) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgResolveMarket) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.Outcome) > 0 { i -= len(m.Outcome); copy(dAtA[i:], m.Outcome); i = encodeVarintPredict(dAtA, i, uint64(len(m.Outcome))); i--; dAtA[i] = 0x1a }
	if m.MarketID != 0 { i = encodeVarintPredict(dAtA, i, uint64(m.MarketID)); i--; dAtA[i] = 0x10 }
	if len(m.Resolver) > 0 { i -= len(m.Resolver); copy(dAtA[i:], m.Resolver); i = encodeVarintPredict(dAtA, i, uint64(len(m.Resolver))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgResolveMarket) Size() (n int) {
	if m == nil { return 0 }
	l := len(m.Resolver); if l > 0 { n += 1 + l + sovPredict(uint64(l)) }
	if m.MarketID != 0 { n += 1 + sovPredict(uint64(m.MarketID)) }
	l = len(m.Outcome); if l > 0 { n += 1 + l + sovPredict(uint64(l)) }
	return n
}
func (m *MsgResolveMarket) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7)
		if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthPredict }; postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthPredict }; if postIndex > l { return io.ErrUnexpectedEOF }; m.Resolver = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.MarketID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.MarketID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		case 3: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthPredict }; postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthPredict }; if postIndex > l { return io.ErrUnexpectedEOF }; m.Outcome = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		default: iNdEx = preIndex; skippy, err := skipPredict(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthPredict }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgResolveMarketResponse (empty)
var xxx_messageInfo_MsgResolveMarketResponse proto.InternalMessageInfo
func (m *MsgResolveMarketResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgResolveMarketResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgResolveMarketResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgResolveMarketResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgResolveMarketResponse.Merge(m, src) }
func (m *MsgResolveMarketResponse) XXX_Size() int { return m.Size() }
func (m *MsgResolveMarketResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgResolveMarketResponse.DiscardUnknown(m) }
func (m *MsgResolveMarketResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgResolveMarketResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgResolveMarketResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgResolveMarketResponse) Size() (n int) { return 0 }
func (m *MsgResolveMarketResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		iNdEx = preIndex; skippy, err := skipPredict(dAtA[iNdEx:]); if err != nil { return err }
		if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthPredict }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
	}
	return nil
}

// MsgClaimWinnings — sender(s,1) market_id(varint,2)
var xxx_messageInfo_MsgClaimWinnings proto.InternalMessageInfo
func (m *MsgClaimWinnings) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgClaimWinnings) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgClaimWinnings.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgClaimWinnings) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgClaimWinnings.Merge(m, src) }
func (m *MsgClaimWinnings) XXX_Size() int { return m.Size() }
func (m *MsgClaimWinnings) XXX_DiscardUnknown() { xxx_messageInfo_MsgClaimWinnings.DiscardUnknown(m) }
func (m *MsgClaimWinnings) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgClaimWinnings) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgClaimWinnings) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.MarketID != 0 { i = encodeVarintPredict(dAtA, i, uint64(m.MarketID)); i--; dAtA[i] = 0x10 }
	if len(m.Sender) > 0 { i -= len(m.Sender); copy(dAtA[i:], m.Sender); i = encodeVarintPredict(dAtA, i, uint64(len(m.Sender))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgClaimWinnings) Size() (n int) {
	if m == nil { return 0 }
	l := len(m.Sender); if l > 0 { n += 1 + l + sovPredict(uint64(l)) }
	if m.MarketID != 0 { n += 1 + sovPredict(uint64(m.MarketID)) }
	return n
}
func (m *MsgClaimWinnings) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3); wireType := int(wire & 0x7)
		if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if wireType != 2 { return fmt.Errorf("proto: wrong wireType") }; var stringLen uint64; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break } }; intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthPredict }; postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthPredict }; if postIndex > l { return io.ErrUnexpectedEOF }; m.Sender = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2: if wireType != 0 { return fmt.Errorf("proto: wrong wireType") }; m.MarketID = 0; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; m.MarketID |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		default: iNdEx = preIndex; skippy, err := skipPredict(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthPredict }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// MsgClaimWinningsResponse — amount(bytes,1)
var xxx_messageInfo_MsgClaimWinningsResponse proto.InternalMessageInfo
func (m *MsgClaimWinningsResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgClaimWinningsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgClaimWinningsResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgClaimWinningsResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgClaimWinningsResponse.Merge(m, src) }
func (m *MsgClaimWinningsResponse) XXX_Size() int { return m.Size() }
func (m *MsgClaimWinningsResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgClaimWinningsResponse.DiscardUnknown(m) }
func (m *MsgClaimWinningsResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgClaimWinningsResponse) MarshalTo(dAtA []byte) (int, error) { return m.MarshalToSizedBuffer(dAtA[:m.Size()]) }
func (m *MsgClaimWinningsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	bz, err := m.Amount.Marshal(); if err != nil { return 0, err }; if len(bz) > 0 { i -= len(bz); copy(dAtA[i:], bz); i = encodeVarintPredict(dAtA, i, uint64(len(bz))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgClaimWinningsResponse) Size() (n int) { if m == nil { return 0 }; bz, _ := m.Amount.Marshal(); l := len(bz); if l > 0 { n += 1 + l + sovPredict(uint64(l)) }; return n }
func (m *MsgClaimWinningsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA); iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx; var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break } }
		fieldNum := int32(wire >> 3)
		if fieldNum <= 0 { return fmt.Errorf("proto: illegal tag %d", fieldNum) }
		switch fieldNum {
		case 1: if int(wire&0x7) != 2 { return fmt.Errorf("proto: wrong wireType") }; var byteLen int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return ErrIntOverflowPredict }; if iNdEx >= l { return io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; byteLen |= int(b&0x7F) << shift; if b < 0x80 { break } }; if byteLen < 0 { return ErrInvalidLengthPredict }; postIndex := iNdEx + byteLen; if postIndex < 0 { return ErrInvalidLengthPredict }; if postIndex > l { return io.ErrUnexpectedEOF }; if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }; iNdEx = postIndex
		default: iNdEx = preIndex; skippy, err := skipPredict(dAtA[iNdEx:]); if err != nil { return err }; if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthPredict }; if (iNdEx+skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }; return nil
}

// Helpers
var (
	ErrIntOverflowPredict  = fmt.Errorf("proto: integer overflow")
	ErrInvalidLengthPredict = fmt.Errorf("proto: negative length found during unmarshaling")
)

func encodeVarintPredict(dAtA []byte, offset int, v uint64) int {
	offset -= sovPredict(v); base := offset
	for v >= 1<<7 { dAtA[offset] = uint8(v&0x7f | 0x80); v >>= 7; offset++ }; dAtA[offset] = uint8(v); return base
}
func sovPredict(x uint64) int { return (bits.Len64(x|1) + 6) / 7 }
func skipPredict(dAtA []byte) (int, error) {
	l := len(dAtA); iNdEx := 0; depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 { if shift >= 64 { return 0, ErrIntOverflowPredict }; if iNdEx >= l { return 0, io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; wire |= (uint64(b) & 0x7F) << shift; if b < 0x80 { break } }
		wireType := int(wire & 0x7)
		switch wireType {
		case 0: for shift := uint(0); ; shift += 7 { if shift >= 64 { return 0, ErrIntOverflowPredict }; if iNdEx >= l { return 0, io.ErrUnexpectedEOF }; iNdEx++; if dAtA[iNdEx-1] < 0x80 { break } }
		case 1: iNdEx += 8
		case 2: var length int; for shift := uint(0); ; shift += 7 { if shift >= 64 { return 0, ErrIntOverflowPredict }; if iNdEx >= l { return 0, io.ErrUnexpectedEOF }; b := dAtA[iNdEx]; iNdEx++; length |= (int(b) & 0x7F) << shift; if b < 0x80 { break } }; if length < 0 { return 0, ErrInvalidLengthPredict }; iNdEx += length
		case 3: depth++
		case 4: if depth == 0 { return 0, fmt.Errorf("proto: wiretype end group for non-group") }; depth--
		case 5: iNdEx += 4
		default: return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 { return 0, ErrInvalidLengthPredict }
		if depth == 0 { return iNdEx, nil }
	}
	return 0, io.ErrUnexpectedEOF
}
