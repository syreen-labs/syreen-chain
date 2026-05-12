package types

import (
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

// ---------------------------------------------------------------------------
// MsgCreateInvoice (6 fields: creator, booking_id, payee, amount(Coin), description, currency)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgCreateInvoice proto.InternalMessageInfo

func (m *MsgCreateInvoice) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateInvoice) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateInvoice.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgCreateInvoice) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCreateInvoice.Merge(m, src) }
func (m *MsgCreateInvoice) XXX_Size() int               { return m.Size() }
func (m *MsgCreateInvoice) XXX_DiscardUnknown()         { xxx_messageInfo_MsgCreateInvoice.DiscardUnknown(m) }

func (m *MsgCreateInvoice) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}
func (m *MsgCreateInvoice) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgCreateInvoice) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.Currency) > 0 {
		i -= len(m.Currency)
		copy(dAtA[i:], m.Currency)
		i = encodeVarintPayments(dAtA, i, uint64(len(m.Currency)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintPayments(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x2a
	}
	{
		size, err := m.Amount.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintPayments(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if len(m.Payee) > 0 {
		i -= len(m.Payee)
		copy(dAtA[i:], m.Payee)
		i = encodeVarintPayments(dAtA, i, uint64(len(m.Payee)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.BookingID) > 0 {
		i -= len(m.BookingID)
		copy(dAtA[i:], m.BookingID)
		i = encodeVarintPayments(dAtA, i, uint64(len(m.BookingID)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintPayments(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}
func (m *MsgCreateInvoice) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovPayments(uint64(l))
	}
	l = len(m.BookingID)
	if l > 0 {
		n += 1 + l + sovPayments(uint64(l))
	}
	l = len(m.Payee)
	if l > 0 {
		n += 1 + l + sovPayments(uint64(l))
	}
	l = m.Amount.Size()
	n += 1 + l + sovPayments(uint64(l))
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovPayments(uint64(l))
	}
	l = len(m.Currency)
	if l > 0 {
		n += 1 + l + sovPayments(uint64(l))
	}
	return n
}
func (m *MsgCreateInvoice) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 { return ErrIntOverflowPayments }
			if iNdEx >= l { return io.ErrUnexpectedEOF }
			b := dAtA[iNdEx]; iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 { break }
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 { return fmt.Errorf("proto: MsgCreateInvoice: wiretype end group for non-group") }
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgCreateInvoice: illegal tag %d (wire type %d)", fieldNum, wire) }
		switch fieldNum {
		case 1:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflowPayments }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break }
			}
			intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthPayments }
			postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthPayments }; if postIndex > l { return io.ErrUnexpectedEOF }
			m.Creator = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field BookingID", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflowPayments }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break }
			}
			intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthPayments }
			postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthPayments }; if postIndex > l { return io.ErrUnexpectedEOF }
			m.BookingID = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 3:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Payee", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflowPayments }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break }
			}
			intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthPayments }
			postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthPayments }; if postIndex > l { return io.ErrUnexpectedEOF }
			m.Payee = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 4:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType) }
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflowPayments }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++; msglen |= int(b&0x7F) << shift; if b < 0x80 { break }
			}
			if msglen < 0 { return ErrInvalidLengthPayments }
			postIndex := iNdEx + msglen; if postIndex < 0 { return ErrInvalidLengthPayments }; if postIndex > l { return io.ErrUnexpectedEOF }
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }
			iNdEx = postIndex
		case 5:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflowPayments }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break }
			}
			intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthPayments }
			postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthPayments }; if postIndex > l { return io.ErrUnexpectedEOF }
			m.Description = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 6:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Currency", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflowPayments }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break }
			}
			intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthPayments }
			postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthPayments }; if postIndex > l { return io.ErrUnexpectedEOF }
			m.Currency = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		default:
			iNdEx = preIndex; skippy, err := skipPayments(dAtA[iNdEx:]); if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthPayments }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }
	return nil
}

// MsgCreateInvoiceResponse
var xxx_messageInfo_MsgCreateInvoiceResponse proto.InternalMessageInfo
func (m *MsgCreateInvoiceResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCreateInvoiceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgCreateInvoiceResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgCreateInvoiceResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCreateInvoiceResponse.Merge(m, src) }
func (m *MsgCreateInvoiceResponse) XXX_Size() int { return m.Size() }
func (m *MsgCreateInvoiceResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgCreateInvoiceResponse.DiscardUnknown(m) }
func (m *MsgCreateInvoiceResponse) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgCreateInvoiceResponse) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgCreateInvoiceResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.InvoiceID) > 0 { i -= len(m.InvoiceID); copy(dAtA[i:], m.InvoiceID); i = encodeVarintPayments(dAtA, i, uint64(len(m.InvoiceID))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgCreateInvoiceResponse) Size() (n int) { if m == nil { return 0 }; l := len(m.InvoiceID); if l > 0 { n += 1 + l + sovPayments(uint64(l)) }; return n }
func (m *MsgCreateInvoiceResponse) Unmarshal(dAtA []byte) error {
	return unmarshalStringsPayments(dAtA, "MsgCreateInvoiceResponse", func(fieldNum int32, val string) {
		if fieldNum == 1 { m.InvoiceID = val }
	})
}

// MsgPayInvoice (3 string fields)
var xxx_messageInfo_MsgPayInvoice proto.InternalMessageInfo
func (m *MsgPayInvoice) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgPayInvoice) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgPayInvoice.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgPayInvoice) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgPayInvoice.Merge(m, src) }
func (m *MsgPayInvoice) XXX_Size() int { return m.Size() }
func (m *MsgPayInvoice) XXX_DiscardUnknown() { xxx_messageInfo_MsgPayInvoice.DiscardUnknown(m) }
func (m *MsgPayInvoice) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgPayInvoice) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgPayInvoice) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.PaymentMethod) > 0 { i -= len(m.PaymentMethod); copy(dAtA[i:], m.PaymentMethod); i = encodeVarintPayments(dAtA, i, uint64(len(m.PaymentMethod))); i--; dAtA[i] = 0x1a }
	if len(m.InvoiceID) > 0 { i -= len(m.InvoiceID); copy(dAtA[i:], m.InvoiceID); i = encodeVarintPayments(dAtA, i, uint64(len(m.InvoiceID))); i--; dAtA[i] = 0x12 }
	if len(m.Payer) > 0 { i -= len(m.Payer); copy(dAtA[i:], m.Payer); i = encodeVarintPayments(dAtA, i, uint64(len(m.Payer))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgPayInvoice) Size() (n int) {
	if m == nil { return 0 }; var l int
	l = len(m.Payer); if l > 0 { n += 1 + l + sovPayments(uint64(l)) }
	l = len(m.InvoiceID); if l > 0 { n += 1 + l + sovPayments(uint64(l)) }
	l = len(m.PaymentMethod); if l > 0 { n += 1 + l + sovPayments(uint64(l)) }
	return n
}
func (m *MsgPayInvoice) Unmarshal(dAtA []byte) error {
	return unmarshalStringsPayments(dAtA, "MsgPayInvoice", func(fieldNum int32, val string) {
		switch fieldNum { case 1: m.Payer = val; case 2: m.InvoiceID = val; case 3: m.PaymentMethod = val }
	})
}

// MsgPayInvoiceResponse
var xxx_messageInfo_MsgPayInvoiceResponse proto.InternalMessageInfo
func (m *MsgPayInvoiceResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgPayInvoiceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgPayInvoiceResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgPayInvoiceResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgPayInvoiceResponse.Merge(m, src) }
func (m *MsgPayInvoiceResponse) XXX_Size() int { return m.Size() }
func (m *MsgPayInvoiceResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgPayInvoiceResponse.DiscardUnknown(m) }
func (m *MsgPayInvoiceResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgPayInvoiceResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgPayInvoiceResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgPayInvoiceResponse) Size() (n int) { return 0 }
func (m *MsgPayInvoiceResponse) Unmarshal(dAtA []byte) error { return unmarshalEmptyPayments(dAtA, "MsgPayInvoiceResponse") }

// MsgRefundPayment (3 string fields)
var xxx_messageInfo_MsgRefundPayment proto.InternalMessageInfo
func (m *MsgRefundPayment) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRefundPayment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgRefundPayment.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgRefundPayment) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgRefundPayment.Merge(m, src) }
func (m *MsgRefundPayment) XXX_Size() int { return m.Size() }
func (m *MsgRefundPayment) XXX_DiscardUnknown() { xxx_messageInfo_MsgRefundPayment.DiscardUnknown(m) }
func (m *MsgRefundPayment) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgRefundPayment) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgRefundPayment) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.Reason) > 0 { i -= len(m.Reason); copy(dAtA[i:], m.Reason); i = encodeVarintPayments(dAtA, i, uint64(len(m.Reason))); i--; dAtA[i] = 0x1a }
	if len(m.InvoiceID) > 0 { i -= len(m.InvoiceID); copy(dAtA[i:], m.InvoiceID); i = encodeVarintPayments(dAtA, i, uint64(len(m.InvoiceID))); i--; dAtA[i] = 0x12 }
	if len(m.Authority) > 0 { i -= len(m.Authority); copy(dAtA[i:], m.Authority); i = encodeVarintPayments(dAtA, i, uint64(len(m.Authority))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgRefundPayment) Size() (n int) {
	if m == nil { return 0 }; var l int
	l = len(m.Authority); if l > 0 { n += 1 + l + sovPayments(uint64(l)) }
	l = len(m.InvoiceID); if l > 0 { n += 1 + l + sovPayments(uint64(l)) }
	l = len(m.Reason); if l > 0 { n += 1 + l + sovPayments(uint64(l)) }
	return n
}
func (m *MsgRefundPayment) Unmarshal(dAtA []byte) error {
	return unmarshalStringsPayments(dAtA, "MsgRefundPayment", func(fieldNum int32, val string) {
		switch fieldNum { case 1: m.Authority = val; case 2: m.InvoiceID = val; case 3: m.Reason = val }
	})
}

// MsgRefundPaymentResponse
var xxx_messageInfo_MsgRefundPaymentResponse proto.InternalMessageInfo
func (m *MsgRefundPaymentResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRefundPaymentResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgRefundPaymentResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgRefundPaymentResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgRefundPaymentResponse.Merge(m, src) }
func (m *MsgRefundPaymentResponse) XXX_Size() int { return m.Size() }
func (m *MsgRefundPaymentResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgRefundPaymentResponse.DiscardUnknown(m) }
func (m *MsgRefundPaymentResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgRefundPaymentResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgRefundPaymentResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgRefundPaymentResponse) Size() (n int) { return 0 }
func (m *MsgRefundPaymentResponse) Unmarshal(dAtA []byte) error { return unmarshalEmptyPayments(dAtA, "MsgRefundPaymentResponse") }

// MsgSetExchangeRate (4 string fields)
var xxx_messageInfo_MsgSetExchangeRate proto.InternalMessageInfo
func (m *MsgSetExchangeRate) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgSetExchangeRate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgSetExchangeRate.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgSetExchangeRate) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgSetExchangeRate.Merge(m, src) }
func (m *MsgSetExchangeRate) XXX_Size() int { return m.Size() }
func (m *MsgSetExchangeRate) XXX_DiscardUnknown() { xxx_messageInfo_MsgSetExchangeRate.DiscardUnknown(m) }
func (m *MsgSetExchangeRate) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgSetExchangeRate) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgSetExchangeRate) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.Rate) > 0 { i -= len(m.Rate); copy(dAtA[i:], m.Rate); i = encodeVarintPayments(dAtA, i, uint64(len(m.Rate))); i--; dAtA[i] = 0x22 }
	if len(m.ToDenom) > 0 { i -= len(m.ToDenom); copy(dAtA[i:], m.ToDenom); i = encodeVarintPayments(dAtA, i, uint64(len(m.ToDenom))); i--; dAtA[i] = 0x1a }
	if len(m.FromDenom) > 0 { i -= len(m.FromDenom); copy(dAtA[i:], m.FromDenom); i = encodeVarintPayments(dAtA, i, uint64(len(m.FromDenom))); i--; dAtA[i] = 0x12 }
	if len(m.Authority) > 0 { i -= len(m.Authority); copy(dAtA[i:], m.Authority); i = encodeVarintPayments(dAtA, i, uint64(len(m.Authority))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgSetExchangeRate) Size() (n int) {
	if m == nil { return 0 }; var l int
	l = len(m.Authority); if l > 0 { n += 1 + l + sovPayments(uint64(l)) }
	l = len(m.FromDenom); if l > 0 { n += 1 + l + sovPayments(uint64(l)) }
	l = len(m.ToDenom); if l > 0 { n += 1 + l + sovPayments(uint64(l)) }
	l = len(m.Rate); if l > 0 { n += 1 + l + sovPayments(uint64(l)) }
	return n
}
func (m *MsgSetExchangeRate) Unmarshal(dAtA []byte) error {
	return unmarshalStringsPayments(dAtA, "MsgSetExchangeRate", func(fieldNum int32, val string) {
		switch fieldNum { case 1: m.Authority = val; case 2: m.FromDenom = val; case 3: m.ToDenom = val; case 4: m.Rate = val }
	})
}

// MsgSetExchangeRateResponse
var xxx_messageInfo_MsgSetExchangeRateResponse proto.InternalMessageInfo
func (m *MsgSetExchangeRateResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgSetExchangeRateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgSetExchangeRateResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgSetExchangeRateResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgSetExchangeRateResponse.Merge(m, src) }
func (m *MsgSetExchangeRateResponse) XXX_Size() int { return m.Size() }
func (m *MsgSetExchangeRateResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgSetExchangeRateResponse.DiscardUnknown(m) }
func (m *MsgSetExchangeRateResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgSetExchangeRateResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgSetExchangeRateResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgSetExchangeRateResponse) Size() (n int) { return 0 }
func (m *MsgSetExchangeRateResponse) Unmarshal(dAtA []byte) error { return unmarshalEmptyPayments(dAtA, "MsgSetExchangeRateResponse") }

// MsgWithdrawEarnings (address + amount Coin)
var xxx_messageInfo_MsgWithdrawEarnings proto.InternalMessageInfo
func (m *MsgWithdrawEarnings) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgWithdrawEarnings) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgWithdrawEarnings.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgWithdrawEarnings) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgWithdrawEarnings.Merge(m, src) }
func (m *MsgWithdrawEarnings) XXX_Size() int { return m.Size() }
func (m *MsgWithdrawEarnings) XXX_DiscardUnknown() { xxx_messageInfo_MsgWithdrawEarnings.DiscardUnknown(m) }
func (m *MsgWithdrawEarnings) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgWithdrawEarnings) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgWithdrawEarnings) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	{
		size, err := m.Amount.MarshalToSizedBuffer(dAtA[:i])
		if err != nil { return 0, err }
		i -= size
		i = encodeVarintPayments(dAtA, i, uint64(size))
	}
	i--; dAtA[i] = 0x12
	if len(m.Address) > 0 { i -= len(m.Address); copy(dAtA[i:], m.Address); i = encodeVarintPayments(dAtA, i, uint64(len(m.Address))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgWithdrawEarnings) Size() (n int) {
	if m == nil { return 0 }; var l int
	l = len(m.Address); if l > 0 { n += 1 + l + sovPayments(uint64(l)) }
	l = m.Amount.Size(); n += 1 + l + sovPayments(uint64(l))
	return n
}
func (m *MsgWithdrawEarnings) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 { return ErrIntOverflowPayments }
			if iNdEx >= l { return io.ErrUnexpectedEOF }
			b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break }
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 { return fmt.Errorf("proto: MsgWithdrawEarnings: wiretype end group for non-group") }
		if fieldNum <= 0 { return fmt.Errorf("proto: MsgWithdrawEarnings: illegal tag %d (wire type %d)", fieldNum, wire) }
		switch fieldNum {
		case 1:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType) }
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflowPayments }; if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break }
			}
			intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthPayments }
			postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthPayments }; if postIndex > l { return io.ErrUnexpectedEOF }
			m.Address = string(dAtA[iNdEx:postIndex]); iNdEx = postIndex
		case 2:
			if wireType != 2 { return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType) }
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflowPayments }; if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++; msglen |= int(b&0x7F) << shift; if b < 0x80 { break }
			}
			if msglen < 0 { return ErrInvalidLengthPayments }
			postIndex := iNdEx + msglen; if postIndex < 0 { return ErrInvalidLengthPayments }; if postIndex > l { return io.ErrUnexpectedEOF }
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil { return err }
			iNdEx = postIndex
		default:
			iNdEx = preIndex; skippy, err := skipPayments(dAtA[iNdEx:]); if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthPayments }; if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }; iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }
	return nil
}

// MsgWithdrawEarningsResponse
var xxx_messageInfo_MsgWithdrawEarningsResponse proto.InternalMessageInfo
func (m *MsgWithdrawEarningsResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgWithdrawEarningsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgWithdrawEarningsResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgWithdrawEarningsResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgWithdrawEarningsResponse.Merge(m, src) }
func (m *MsgWithdrawEarningsResponse) XXX_Size() int { return m.Size() }
func (m *MsgWithdrawEarningsResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgWithdrawEarningsResponse.DiscardUnknown(m) }
func (m *MsgWithdrawEarningsResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgWithdrawEarningsResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgWithdrawEarningsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgWithdrawEarningsResponse) Size() (n int) { return 0 }
func (m *MsgWithdrawEarningsResponse) Unmarshal(dAtA []byte) error { return unmarshalEmptyPayments(dAtA, "MsgWithdrawEarningsResponse") }

// MsgCancelInvoice (2 string fields)
var xxx_messageInfo_MsgCancelInvoice proto.InternalMessageInfo
func (m *MsgCancelInvoice) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCancelInvoice) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgCancelInvoice.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgCancelInvoice) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCancelInvoice.Merge(m, src) }
func (m *MsgCancelInvoice) XXX_Size() int { return m.Size() }
func (m *MsgCancelInvoice) XXX_DiscardUnknown() { xxx_messageInfo_MsgCancelInvoice.DiscardUnknown(m) }
func (m *MsgCancelInvoice) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgCancelInvoice) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgCancelInvoice) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.InvoiceID) > 0 { i -= len(m.InvoiceID); copy(dAtA[i:], m.InvoiceID); i = encodeVarintPayments(dAtA, i, uint64(len(m.InvoiceID))); i--; dAtA[i] = 0x12 }
	if len(m.Creator) > 0 { i -= len(m.Creator); copy(dAtA[i:], m.Creator); i = encodeVarintPayments(dAtA, i, uint64(len(m.Creator))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgCancelInvoice) Size() (n int) {
	if m == nil { return 0 }; var l int
	l = len(m.Creator); if l > 0 { n += 1 + l + sovPayments(uint64(l)) }
	l = len(m.InvoiceID); if l > 0 { n += 1 + l + sovPayments(uint64(l)) }
	return n
}
func (m *MsgCancelInvoice) Unmarshal(dAtA []byte) error {
	return unmarshalStringsPayments(dAtA, "MsgCancelInvoice", func(fieldNum int32, val string) {
		switch fieldNum { case 1: m.Creator = val; case 2: m.InvoiceID = val }
	})
}

// MsgCancelInvoiceResponse
var xxx_messageInfo_MsgCancelInvoiceResponse proto.InternalMessageInfo
func (m *MsgCancelInvoiceResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgCancelInvoiceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgCancelInvoiceResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgCancelInvoiceResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgCancelInvoiceResponse.Merge(m, src) }
func (m *MsgCancelInvoiceResponse) XXX_Size() int { return m.Size() }
func (m *MsgCancelInvoiceResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgCancelInvoiceResponse.DiscardUnknown(m) }
func (m *MsgCancelInvoiceResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgCancelInvoiceResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgCancelInvoiceResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgCancelInvoiceResponse) Size() (n int) { return 0 }
func (m *MsgCancelInvoiceResponse) Unmarshal(dAtA []byte) error { return unmarshalEmptyPayments(dAtA, "MsgCancelInvoiceResponse") }

// ---------------------------------------------------------------------------
// Helper functions (payments-specific names to avoid conflicts)
// ---------------------------------------------------------------------------

var (
	ErrIntOverflowPayments  = fmt.Errorf("proto: integer overflow")
	ErrInvalidLengthPayments = fmt.Errorf("proto: negative length found during unmarshaling")
)

func encodeVarintPayments(dAtA []byte, offset int, v uint64) int {
	offset -= sovPayments(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}

func sovPayments(x uint64) int {
	return (bits.Len64(x|1) + 6) / 7
}

func skipPayments(dAtA []byte) (int, error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 { return 0, ErrIntOverflowPayments }
			if iNdEx >= l { return 0, io.ErrUnexpectedEOF }
			b := dAtA[iNdEx]; iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 { break }
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return 0, ErrIntOverflowPayments }
				if iNdEx >= l { return 0, io.ErrUnexpectedEOF }
				iNdEx++; if dAtA[iNdEx-1] < 0x80 { break }
			}
		case 1: iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return 0, ErrIntOverflowPayments }
				if iNdEx >= l { return 0, io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 { break }
			}
			if length < 0 { return 0, ErrInvalidLengthPayments }
			iNdEx += length
		case 3: depth++
		case 4:
			if depth == 0 { return 0, fmt.Errorf("proto: wiretype end group for non-group") }
			depth--
		case 5: iNdEx += 4
		default: return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 { return 0, ErrInvalidLengthPayments }
		if depth == 0 { return iNdEx, nil }
	}
	return 0, io.ErrUnexpectedEOF
}

func unmarshalEmptyPayments(dAtA []byte, msgName string) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 { return ErrIntOverflowPayments }
			if iNdEx >= l { return io.ErrUnexpectedEOF }
			b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break }
		}
		fieldNum := int32(wire >> 3)
		if fieldNum <= 0 { return fmt.Errorf("proto: %s: illegal tag %d", msgName, fieldNum) }
		iNdEx = preIndex
		skippy, err := skipPayments(dAtA[iNdEx:])
		if err != nil { return err }
		if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthPayments }
		if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }
		iNdEx += skippy
	}
	return nil
}

func unmarshalStringsPayments(dAtA []byte, msgName string, setter func(fieldNum int32, val string)) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 { return ErrIntOverflowPayments }
			if iNdEx >= l { return io.ErrUnexpectedEOF }
			b := dAtA[iNdEx]; iNdEx++; wire |= uint64(b&0x7F) << shift; if b < 0x80 { break }
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 { return fmt.Errorf("proto: %s: wiretype end group for non-group", msgName) }
		if fieldNum <= 0 { return fmt.Errorf("proto: %s: illegal tag %d (wire type %d)", msgName, fieldNum, wire) }
		if wireType == 2 {
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflowPayments }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++; stringLen |= uint64(b&0x7F) << shift; if b < 0x80 { break }
			}
			intStringLen := int(stringLen); if intStringLen < 0 { return ErrInvalidLengthPayments }
			postIndex := iNdEx + intStringLen; if postIndex < 0 { return ErrInvalidLengthPayments }; if postIndex > l { return io.ErrUnexpectedEOF }
			setter(fieldNum, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		} else {
			iNdEx = preIndex
			skippy, err := skipPayments(dAtA[iNdEx:])
			if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthPayments }
			if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }
			iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }
	return nil
}

// suppress unused import
var _ sdk.Coin
