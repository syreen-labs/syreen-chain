package types

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/bits"

	proto "github.com/cosmos/gogoproto/proto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Reference imports to suppress errors
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = io.ErrUnexpectedEOF

// Descriptor methods - return the gzipped file descriptor and the message index
func (*MsgCreatePool) Descriptor() ([]byte, []int)             { return fileDescriptorTx, []int{0} }
func (*MsgCreatePoolResponse) Descriptor() ([]byte, []int)     { return fileDescriptorTx, []int{1} }
func (*MsgAddLiquidity) Descriptor() ([]byte, []int)           { return fileDescriptorTx, []int{2} }
func (*MsgAddLiquidityResponse) Descriptor() ([]byte, []int)   { return fileDescriptorTx, []int{3} }
func (*MsgRemoveLiquidity) Descriptor() ([]byte, []int)        { return fileDescriptorTx, []int{4} }
func (*MsgRemoveLiquidityResponse) Descriptor() ([]byte, []int) { return fileDescriptorTx, []int{5} }
func (*MsgSwap) Descriptor() ([]byte, []int)                   { return fileDescriptorTx, []int{6} }
func (*MsgSwapResponse) Descriptor() ([]byte, []int)           { return fileDescriptorTx, []int{7} }
func (*MsgCreateReferralCode) Descriptor() ([]byte, []int)     { return fileDescriptorTx, []int{8} }
func (*MsgCreateReferralCodeResponse) Descriptor() ([]byte, []int) { return fileDescriptorTx, []int{9} }
func (*MsgRegisterReferral) Descriptor() ([]byte, []int)       { return fileDescriptorTx, []int{10} }
func (*MsgRegisterReferralResponse) Descriptor() ([]byte, []int) { return fileDescriptorTx, []int{11} }
func (*MsgPlaceOrder) Descriptor() ([]byte, []int)              { return fileDescriptorTx, []int{12} }
func (*MsgPlaceOrderResponse) Descriptor() ([]byte, []int)      { return fileDescriptorTx, []int{13} }
func (*MsgCancelOrder) Descriptor() ([]byte, []int)             { return fileDescriptorTx, []int{14} }
func (*MsgCancelOrderResponse) Descriptor() ([]byte, []int)     { return fileDescriptorTx, []int{15} }
func (*MsgModifyOrder) Descriptor() ([]byte, []int)             { return fileDescriptorTx, []int{16} }
func (*MsgModifyOrderResponse) Descriptor() ([]byte, []int)     { return fileDescriptorTx, []int{17} }
func (*MsgFollowTrader) Descriptor() ([]byte, []int)            { return fileDescriptorTx, []int{18} }
func (*MsgFollowTraderResponse) Descriptor() ([]byte, []int)    { return fileDescriptorTx, []int{19} }
func (*MsgUnfollowTrader) Descriptor() ([]byte, []int)          { return fileDescriptorTx, []int{20} }
func (*MsgUnfollowTraderResponse) Descriptor() ([]byte, []int)  { return fileDescriptorTx, []int{21} }
func (*MsgUpdateCopySettings) Descriptor() ([]byte, []int)      { return fileDescriptorTx, []int{22} }
func (*MsgUpdateCopySettingsResponse) Descriptor() ([]byte, []int) { return fileDescriptorTx, []int{23} }
func (*MsgClaimReferralRewards) Descriptor() ([]byte, []int)         { return fileDescriptorTx, []int{24} }
func (*MsgClaimReferralRewardsResponse) Descriptor() ([]byte, []int) { return fileDescriptorTx, []int{25} }
func (*MsgSetPoolFeeConfig) Descriptor() ([]byte, []int)            { return fileDescriptorTx, []int{28} }
func (*MsgSetPoolFeeConfigResponse) Descriptor() ([]byte, []int)    { return fileDescriptorTx, []int{29} }
func (*MsgMultiHopSwap) Descriptor() ([]byte, []int)                { return fileDescriptorTx, []int{26} }
func (*MsgMultiHopSwapResponse) Descriptor() ([]byte, []int)        { return fileDescriptorTx, []int{27} }

func init() {
	proto.RegisterType((*MsgCreatePool)(nil), "syreen.dex.MsgCreatePool")
	proto.RegisterType((*MsgCreatePoolResponse)(nil), "syreen.dex.MsgCreatePoolResponse")
	proto.RegisterType((*MsgAddLiquidity)(nil), "syreen.dex.MsgAddLiquidity")
	proto.RegisterType((*MsgAddLiquidityResponse)(nil), "syreen.dex.MsgAddLiquidityResponse")
	proto.RegisterType((*MsgRemoveLiquidity)(nil), "syreen.dex.MsgRemoveLiquidity")
	proto.RegisterType((*MsgRemoveLiquidityResponse)(nil), "syreen.dex.MsgRemoveLiquidityResponse")
	proto.RegisterType((*MsgSwap)(nil), "syreen.dex.MsgSwap")
	proto.RegisterType((*MsgSwapResponse)(nil), "syreen.dex.MsgSwapResponse")
	proto.RegisterType((*MsgCreateReferralCode)(nil), "syreen.dex.MsgCreateReferralCode")
	proto.RegisterType((*MsgCreateReferralCodeResponse)(nil), "syreen.dex.MsgCreateReferralCodeResponse")
	proto.RegisterType((*MsgRegisterReferral)(nil), "syreen.dex.MsgRegisterReferral")
	proto.RegisterType((*MsgRegisterReferralResponse)(nil), "syreen.dex.MsgRegisterReferralResponse")
	proto.RegisterType((*MsgPlaceOrder)(nil), "syreen.dex.MsgPlaceOrder")
	proto.RegisterType((*MsgPlaceOrderResponse)(nil), "syreen.dex.MsgPlaceOrderResponse")
	proto.RegisterType((*MsgCancelOrder)(nil), "syreen.dex.MsgCancelOrder")
	proto.RegisterType((*MsgCancelOrderResponse)(nil), "syreen.dex.MsgCancelOrderResponse")
	proto.RegisterType((*MsgModifyOrder)(nil), "syreen.dex.MsgModifyOrder")
	proto.RegisterType((*MsgModifyOrderResponse)(nil), "syreen.dex.MsgModifyOrderResponse")
	proto.RegisterType((*MsgFollowTrader)(nil), "syreen.dex.MsgFollowTrader")
	proto.RegisterType((*MsgFollowTraderResponse)(nil), "syreen.dex.MsgFollowTraderResponse")
	proto.RegisterType((*MsgUnfollowTrader)(nil), "syreen.dex.MsgUnfollowTrader")
	proto.RegisterType((*MsgUnfollowTraderResponse)(nil), "syreen.dex.MsgUnfollowTraderResponse")
	proto.RegisterType((*MsgUpdateCopySettings)(nil), "syreen.dex.MsgUpdateCopySettings")
	proto.RegisterType((*MsgUpdateCopySettingsResponse)(nil), "syreen.dex.MsgUpdateCopySettingsResponse")
	proto.RegisterType((*MsgClaimReferralRewards)(nil), "syreen.dex.MsgClaimReferralRewards")
	proto.RegisterType((*MsgClaimReferralRewardsResponse)(nil), "syreen.dex.MsgClaimReferralRewardsResponse")
	proto.RegisterType((*MsgSetPoolFeeConfig)(nil), "syreen.dex.MsgSetPoolFeeConfig")
	proto.RegisterType((*MsgMultiHopSwap)(nil), "syreen.dex.MsgMultiHopSwap")
	proto.RegisterType((*MsgMultiHopSwapResponse)(nil), "syreen.dex.MsgMultiHopSwapResponse")
	proto.RegisterType((*MsgSetPoolFeeConfigResponse)(nil), "syreen.dex.MsgSetPoolFeeConfigResponse")
}

// ---------------------------------------------------------------------------
// MsgCreatePool
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCreatePool proto.InternalMessageInfo

func (m *MsgCreatePool) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreatePool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreatePool.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgCreatePool) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgCreatePool.Merge(m, src) }
func (m *MsgCreatePool) XXX_Size() int                { return m.Size() }
func (m *MsgCreatePool) XXX_DiscardUnknown()          { xxx_messageInfo_MsgCreatePool.DiscardUnknown(m) }

func (m *MsgCreatePool) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCreatePool) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCreatePool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	// field 5: amount_b (bytes - math.Int)
	{
		bz, err := m.AmountB.Marshal()
		if err != nil {
			return 0, err
		}
		if len(bz) > 0 {
			i -= len(bz)
			copy(dAtA[i:], bz)
			i = encodeVarint(dAtA, i, uint64(len(bz)))
			i--
			dAtA[i] = 0x2a
		}
	}
	// field 4: amount_a (bytes - math.Int)
	{
		bz, err := m.AmountA.Marshal()
		if err != nil {
			return 0, err
		}
		if len(bz) > 0 {
			i -= len(bz)
			copy(dAtA[i:], bz)
			i = encodeVarint(dAtA, i, uint64(len(bz)))
			i--
			dAtA[i] = 0x22
		}
	}
	// field 3: denom_b
	if len(m.DenomB) > 0 {
		i -= len(m.DenomB)
		copy(dAtA[i:], m.DenomB)
		i = encodeVarint(dAtA, i, uint64(len(m.DenomB)))
		i--
		dAtA[i] = 0x1a
	}
	// field 2: denom_a
	if len(m.DenomA) > 0 {
		i -= len(m.DenomA)
		copy(dAtA[i:], m.DenomA)
		i = encodeVarint(dAtA, i, uint64(len(m.DenomA)))
		i--
		dAtA[i] = 0x12
	}
	// field 1: sender
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarint(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}

func (m *MsgCreatePool) Size() (n int) {
	if m == nil {
		return 0
	}
	l := len(m.Sender)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	l = len(m.DenomA)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	l = len(m.DenomB)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	bz, _ := m.AmountA.Marshal()
	l = len(bz)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	bz, _ = m.AmountB.Marshal()
	l = len(bz)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	return n
}

func (m *MsgCreatePool) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflow
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgCreatePool: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCreatePool: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1: // sender
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2: // denom_a
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DenomA", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DenomA = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3: // denom_b
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DenomB", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DenomB = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4: // amount_a
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AmountA", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.AmountA.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5: // amount_b
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AmountB", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.AmountB.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgCreatePoolResponse
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCreatePoolResponse proto.InternalMessageInfo

func (m *MsgCreatePoolResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreatePoolResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreatePoolResponse.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgCreatePoolResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCreatePoolResponse.Merge(m, src)
}
func (m *MsgCreatePoolResponse) XXX_Size() int       { return m.Size() }
func (m *MsgCreatePoolResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgCreatePoolResponse.DiscardUnknown(m) }

func (m *MsgCreatePoolResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCreatePoolResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCreatePoolResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.PoolID != 0 {
		i = encodeVarint(dAtA, i, uint64(m.PoolID))
		i--
		dAtA[i] = 0x08
	}
	return len(dAtA) - i, nil
}

func (m *MsgCreatePoolResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	if m.PoolID != 0 {
		n += 1 + sov(uint64(m.PoolID))
	}
	return n
}

func (m *MsgCreatePoolResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflow
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgCreatePoolResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCreatePoolResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1: // pool_id
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolID", wireType)
			}
			m.PoolID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PoolID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgAddLiquidity
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgAddLiquidity proto.InternalMessageInfo

func (m *MsgAddLiquidity) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgAddLiquidity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgAddLiquidity.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgAddLiquidity) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgAddLiquidity.Merge(m, src) }
func (m *MsgAddLiquidity) XXX_Size() int                { return m.Size() }
func (m *MsgAddLiquidity) XXX_DiscardUnknown()          { xxx_messageInfo_MsgAddLiquidity.DiscardUnknown(m) }

func (m *MsgAddLiquidity) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgAddLiquidity) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgAddLiquidity) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	// field 5: min_shares_out
	{
		bz, err := m.MinSharesOut.Marshal()
		if err != nil {
			return 0, err
		}
		if len(bz) > 0 {
			i -= len(bz)
			copy(dAtA[i:], bz)
			i = encodeVarint(dAtA, i, uint64(len(bz)))
			i--
			dAtA[i] = 0x2a
		}
	}
	// field 4: amount_b
	{
		bz, err := m.AmountB.Marshal()
		if err != nil {
			return 0, err
		}
		if len(bz) > 0 {
			i -= len(bz)
			copy(dAtA[i:], bz)
			i = encodeVarint(dAtA, i, uint64(len(bz)))
			i--
			dAtA[i] = 0x22
		}
	}
	// field 3: amount_a
	{
		bz, err := m.AmountA.Marshal()
		if err != nil {
			return 0, err
		}
		if len(bz) > 0 {
			i -= len(bz)
			copy(dAtA[i:], bz)
			i = encodeVarint(dAtA, i, uint64(len(bz)))
			i--
			dAtA[i] = 0x1a
		}
	}
	// field 2: pool_id (varint)
	if m.PoolID != 0 {
		i = encodeVarint(dAtA, i, uint64(m.PoolID))
		i--
		dAtA[i] = 0x10
	}
	// field 1: sender
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarint(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}

func (m *MsgAddLiquidity) Size() (n int) {
	if m == nil {
		return 0
	}
	l := len(m.Sender)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	if m.PoolID != 0 {
		n += 1 + sov(uint64(m.PoolID))
	}
	bz, _ := m.AmountA.Marshal()
	l = len(bz)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	bz, _ = m.AmountB.Marshal()
	l = len(bz)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	bz, _ = m.MinSharesOut.Marshal()
	l = len(bz)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	return n
}

func (m *MsgAddLiquidity) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflow
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgAddLiquidity: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgAddLiquidity: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1: // sender
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2: // pool_id
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolID", wireType)
			}
			m.PoolID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PoolID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3: // amount_a
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AmountA", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.AmountA.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4: // amount_b
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AmountB", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.AmountB.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5: // min_shares_out
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinSharesOut", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MinSharesOut.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgAddLiquidityResponse
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgAddLiquidityResponse proto.InternalMessageInfo

func (m *MsgAddLiquidityResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgAddLiquidityResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgAddLiquidityResponse.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgAddLiquidityResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgAddLiquidityResponse.Merge(m, src)
}
func (m *MsgAddLiquidityResponse) XXX_Size() int       { return m.Size() }
func (m *MsgAddLiquidityResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgAddLiquidityResponse.DiscardUnknown(m) }

func (m *MsgAddLiquidityResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgAddLiquidityResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgAddLiquidityResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	{
		bz, err := m.SharesMinted.Marshal()
		if err != nil {
			return 0, err
		}
		if len(bz) > 0 {
			i -= len(bz)
			copy(dAtA[i:], bz)
			i = encodeVarint(dAtA, i, uint64(len(bz)))
			i--
			dAtA[i] = 0x0a
		}
	}
	return len(dAtA) - i, nil
}

func (m *MsgAddLiquidityResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	bz, _ := m.SharesMinted.Marshal()
	l := len(bz)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	return n
}

func (m *MsgAddLiquidityResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflow
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgAddLiquidityResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgAddLiquidityResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1: // shares_minted
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SharesMinted", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.SharesMinted.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgRemoveLiquidity
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgRemoveLiquidity proto.InternalMessageInfo

func (m *MsgRemoveLiquidity) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRemoveLiquidity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRemoveLiquidity.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgRemoveLiquidity) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgRemoveLiquidity.Merge(m, src) }
func (m *MsgRemoveLiquidity) XXX_Size() int                { return m.Size() }
func (m *MsgRemoveLiquidity) XXX_DiscardUnknown()          { xxx_messageInfo_MsgRemoveLiquidity.DiscardUnknown(m) }

func (m *MsgRemoveLiquidity) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgRemoveLiquidity) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgRemoveLiquidity) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	// field 5: min_amount_b_out
	{
		bz, err := m.MinAmountBOut.Marshal()
		if err != nil {
			return 0, err
		}
		if len(bz) > 0 {
			i -= len(bz)
			copy(dAtA[i:], bz)
			i = encodeVarint(dAtA, i, uint64(len(bz)))
			i--
			dAtA[i] = 0x2a
		}
	}
	// field 4: min_amount_a_out
	{
		bz, err := m.MinAmountAOut.Marshal()
		if err != nil {
			return 0, err
		}
		if len(bz) > 0 {
			i -= len(bz)
			copy(dAtA[i:], bz)
			i = encodeVarint(dAtA, i, uint64(len(bz)))
			i--
			dAtA[i] = 0x22
		}
	}
	// field 3: shares_in
	{
		bz, err := m.SharesIn.Marshal()
		if err != nil {
			return 0, err
		}
		if len(bz) > 0 {
			i -= len(bz)
			copy(dAtA[i:], bz)
			i = encodeVarint(dAtA, i, uint64(len(bz)))
			i--
			dAtA[i] = 0x1a
		}
	}
	// field 2: pool_id
	if m.PoolID != 0 {
		i = encodeVarint(dAtA, i, uint64(m.PoolID))
		i--
		dAtA[i] = 0x10
	}
	// field 1: sender
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarint(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}

func (m *MsgRemoveLiquidity) Size() (n int) {
	if m == nil {
		return 0
	}
	l := len(m.Sender)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	if m.PoolID != 0 {
		n += 1 + sov(uint64(m.PoolID))
	}
	bz, _ := m.SharesIn.Marshal()
	l = len(bz)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	bz, _ = m.MinAmountAOut.Marshal()
	l = len(bz)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	bz, _ = m.MinAmountBOut.Marshal()
	l = len(bz)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	return n
}

func (m *MsgRemoveLiquidity) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflow
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgRemoveLiquidity: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgRemoveLiquidity: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1: // sender
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2: // pool_id
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolID", wireType)
			}
			m.PoolID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PoolID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3: // shares_in
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SharesIn", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.SharesIn.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4: // min_amount_a_out
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinAmountAOut", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MinAmountAOut.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5: // min_amount_b_out
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinAmountBOut", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MinAmountBOut.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgRemoveLiquidityResponse
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgRemoveLiquidityResponse proto.InternalMessageInfo

func (m *MsgRemoveLiquidityResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRemoveLiquidityResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRemoveLiquidityResponse.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgRemoveLiquidityResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgRemoveLiquidityResponse.Merge(m, src)
}
func (m *MsgRemoveLiquidityResponse) XXX_Size() int       { return m.Size() }
func (m *MsgRemoveLiquidityResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgRemoveLiquidityResponse.DiscardUnknown(m) }

func (m *MsgRemoveLiquidityResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgRemoveLiquidityResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgRemoveLiquidityResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	// field 2: amount_b
	{
		bz, err := m.AmountB.Marshal()
		if err != nil {
			return 0, err
		}
		if len(bz) > 0 {
			i -= len(bz)
			copy(dAtA[i:], bz)
			i = encodeVarint(dAtA, i, uint64(len(bz)))
			i--
			dAtA[i] = 0x12
		}
	}
	// field 1: amount_a
	{
		bz, err := m.AmountA.Marshal()
		if err != nil {
			return 0, err
		}
		if len(bz) > 0 {
			i -= len(bz)
			copy(dAtA[i:], bz)
			i = encodeVarint(dAtA, i, uint64(len(bz)))
			i--
			dAtA[i] = 0x0a
		}
	}
	return len(dAtA) - i, nil
}

func (m *MsgRemoveLiquidityResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	bz, _ := m.AmountA.Marshal()
	l := len(bz)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	bz, _ = m.AmountB.Marshal()
	l = len(bz)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	return n
}

func (m *MsgRemoveLiquidityResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflow
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgRemoveLiquidityResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgRemoveLiquidityResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1: // amount_a
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AmountA", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.AmountA.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2: // amount_b
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AmountB", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.AmountB.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgSwap
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgSwap proto.InternalMessageInfo

func (m *MsgSwap) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgSwap) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSwap.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgSwap) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgSwap.Merge(m, src) }
func (m *MsgSwap) XXX_Size() int                { return m.Size() }
func (m *MsgSwap) XXX_DiscardUnknown()          { xxx_messageInfo_MsgSwap.DiscardUnknown(m) }

func (m *MsgSwap) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSwap) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSwap) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	// field 4: min_token_out (bytes - math.Int)
	{
		bz, err := m.MinTokenOut.Marshal()
		if err != nil {
			return 0, err
		}
		if len(bz) > 0 {
			i -= len(bz)
			copy(dAtA[i:], bz)
			i = encodeVarint(dAtA, i, uint64(len(bz)))
			i--
			dAtA[i] = 0x2a // field 5 tag would be 0x2a, but we use field 4 = 0x22... wait
		}
	}
	// field 3: token_in (bytes - sdk.Coin, embedded message)
	{
		size, err := m.TokenIn.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	// field 2: pool_id (varint)
	if m.PoolID != 0 {
		i = encodeVarint(dAtA, i, uint64(m.PoolID))
		i--
		dAtA[i] = 0x10
	}
	// field 1: sender
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarint(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}

func (m *MsgSwap) Size() (n int) {
	if m == nil {
		return 0
	}
	l := len(m.Sender)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	if m.PoolID != 0 {
		n += 1 + sov(uint64(m.PoolID))
	}
	l = m.TokenIn.Size()
	n += 1 + l + sov(uint64(l))
	bz, _ := m.MinTokenOut.Marshal()
	l = len(bz)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	return n
}

func (m *MsgSwap) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflow
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgSwap: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSwap: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1: // sender
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2: // pool_id
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolID", wireType)
			}
			m.PoolID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PoolID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3: // token_in (embedded message)
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenIn", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TokenIn.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4: // min_token_out
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinTokenOut", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MinTokenOut.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgSwapResponse
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgSwapResponse proto.InternalMessageInfo

func (m *MsgSwapResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgSwapResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSwapResponse.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgSwapResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSwapResponse.Merge(m, src)
}
func (m *MsgSwapResponse) XXX_Size() int       { return m.Size() }
func (m *MsgSwapResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgSwapResponse.DiscardUnknown(m) }

func (m *MsgSwapResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSwapResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSwapResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	// field 1: token_out (embedded message - sdk.Coin)
	{
		size, err := m.TokenOut.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x0a
	return len(dAtA) - i, nil
}

func (m *MsgSwapResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	l := m.TokenOut.Size()
	n += 1 + l + sov(uint64(l))
	return n
}

func (m *MsgSwapResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflow
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgSwapResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSwapResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1: // token_out
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenOut", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TokenOut.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgCreateReferralCode protobuf methods
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCreateReferralCode proto.InternalMessageInfo

func (m *MsgCreateReferralCode) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateReferralCode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateReferralCode.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgCreateReferralCode) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgCreateReferralCode.Merge(m, src) }
func (m *MsgCreateReferralCode) XXX_Size() int                { return m.Size() }
func (m *MsgCreateReferralCode) XXX_DiscardUnknown()          { xxx_messageInfo_MsgCreateReferralCode.DiscardUnknown(m) }

func (m *MsgCreateReferralCode) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCreateReferralCode) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCreateReferralCode) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarint(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}

func (m *MsgCreateReferralCode) Size() (n int) {
	if m == nil {
		return 0
	}
	l := len(m.Creator)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	return n
}

func (m *MsgCreateReferralCode) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflow
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgCreateReferralCode: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCreateReferralCode: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgCreateReferralCodeResponse protobuf methods
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCreateReferralCodeResponse proto.InternalMessageInfo

func (m *MsgCreateReferralCodeResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateReferralCodeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateReferralCodeResponse.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgCreateReferralCodeResponse) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgCreateReferralCodeResponse.Merge(m, src) }
func (m *MsgCreateReferralCodeResponse) XXX_Size() int                { return m.Size() }
func (m *MsgCreateReferralCodeResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgCreateReferralCodeResponse.DiscardUnknown(m) }

func (m *MsgCreateReferralCodeResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCreateReferralCodeResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCreateReferralCodeResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.ReferralCode) > 0 {
		i -= len(m.ReferralCode)
		copy(dAtA[i:], m.ReferralCode)
		i = encodeVarint(dAtA, i, uint64(len(m.ReferralCode)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}

func (m *MsgCreateReferralCodeResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	l := len(m.ReferralCode)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	return n
}

func (m *MsgCreateReferralCodeResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflow
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgCreateReferralCodeResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCreateReferralCodeResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReferralCode", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ReferralCode = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgRegisterReferral protobuf methods
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgRegisterReferral proto.InternalMessageInfo

func (m *MsgRegisterReferral) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRegisterReferral) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRegisterReferral.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgRegisterReferral) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgRegisterReferral.Merge(m, src) }
func (m *MsgRegisterReferral) XXX_Size() int                { return m.Size() }
func (m *MsgRegisterReferral) XXX_DiscardUnknown()          { xxx_messageInfo_MsgRegisterReferral.DiscardUnknown(m) }

func (m *MsgRegisterReferral) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgRegisterReferral) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgRegisterReferral) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.ReferralCode) > 0 {
		i -= len(m.ReferralCode)
		copy(dAtA[i:], m.ReferralCode)
		i = encodeVarint(dAtA, i, uint64(len(m.ReferralCode)))
		i--
		dAtA[i] = 0x12 // field 2, wire type 2
	}
	if len(m.User) > 0 {
		i -= len(m.User)
		copy(dAtA[i:], m.User)
		i = encodeVarint(dAtA, i, uint64(len(m.User)))
		i--
		dAtA[i] = 0x0a // field 1, wire type 2
	}
	return len(dAtA) - i, nil
}

func (m *MsgRegisterReferral) Size() (n int) {
	if m == nil {
		return 0
	}
	l := len(m.User)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	l = len(m.ReferralCode)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	return n
}

func (m *MsgRegisterReferral) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflow
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgRegisterReferral: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgRegisterReferral: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field User", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.User = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReferralCode", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ReferralCode = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// ---------------------------------------------------------------------------
// MsgRegisterReferralResponse protobuf methods
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgRegisterReferralResponse proto.InternalMessageInfo

func (m *MsgRegisterReferralResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRegisterReferralResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRegisterReferralResponse.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgRegisterReferralResponse) XXX_Merge(src proto.Message)  { xxx_messageInfo_MsgRegisterReferralResponse.Merge(m, src) }
func (m *MsgRegisterReferralResponse) XXX_Size() int                { return m.Size() }
func (m *MsgRegisterReferralResponse) XXX_DiscardUnknown()          { xxx_messageInfo_MsgRegisterReferralResponse.DiscardUnknown(m) }

func (m *MsgRegisterReferralResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgRegisterReferralResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgRegisterReferralResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.Referrer) > 0 {
		i -= len(m.Referrer)
		copy(dAtA[i:], m.Referrer)
		i = encodeVarint(dAtA, i, uint64(len(m.Referrer)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}

func (m *MsgRegisterReferralResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	l := len(m.Referrer)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	return n
}

func (m *MsgRegisterReferralResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflow
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgRegisterReferralResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgRegisterReferralResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Referrer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Referrer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}
	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

// ---------------------------------------------------------------------------
// Helper functions (standard protobuf encoding helpers)
// ---------------------------------------------------------------------------

var (
	ErrIntOverflow   = fmt.Errorf("proto: integer overflow")
	ErrInvalidLength = fmt.Errorf("proto: negative length found during unmarshaling")
)

func encodeVarint(dAtA []byte, offset int, v uint64) int {
	offset -= sov(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}

func sov(x uint64) int {
	return (bits.Len64(x|1) + 6) / 7
}

func skip(dAtA []byte) (int, error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflow
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflow
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflow
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLength
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, fmt.Errorf("proto: wiretype end group for non-group")
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLength
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

// suppress unused imports
var _ sdk.Coin
var _ = binary.BigEndian
