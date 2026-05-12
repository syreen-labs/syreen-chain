package types

import (
	"fmt"
	"io"
	"math/bits"

	proto "github.com/cosmos/gogoproto/proto"
)

// Reference imports to suppress errors
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = io.ErrUnexpectedEOF

// ---------------------------------------------------------------------------
// MsgRegisterIdentity (4 string fields)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgRegisterIdentity proto.InternalMessageInfo

func (m *MsgRegisterIdentity) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRegisterIdentity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRegisterIdentity.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgRegisterIdentity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgRegisterIdentity.Merge(m, src)
}
func (m *MsgRegisterIdentity) XXX_Size() int       { return m.Size() }
func (m *MsgRegisterIdentity) XXX_DiscardUnknown() { xxx_messageInfo_MsgRegisterIdentity.DiscardUnknown(m) }

func (m *MsgRegisterIdentity) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}
func (m *MsgRegisterIdentity) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgRegisterIdentity) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.Level) > 0 {
		i -= len(m.Level)
		copy(dAtA[i:], m.Level)
		i = encodeVarintIdentity(dAtA, i, uint64(len(m.Level)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Nationality) > 0 {
		i -= len(m.Nationality)
		copy(dAtA[i:], m.Nationality)
		i = encodeVarintIdentity(dAtA, i, uint64(len(m.Nationality)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.DocumentHash) > 0 {
		i -= len(m.DocumentHash)
		copy(dAtA[i:], m.DocumentHash)
		i = encodeVarintIdentity(dAtA, i, uint64(len(m.DocumentHash)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintIdentity(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}
func (m *MsgRegisterIdentity) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovIdentity(uint64(l))
	}
	l = len(m.DocumentHash)
	if l > 0 {
		n += 1 + l + sovIdentity(uint64(l))
	}
	l = len(m.Nationality)
	if l > 0 {
		n += 1 + l + sovIdentity(uint64(l))
	}
	l = len(m.Level)
	if l > 0 {
		n += 1 + l + sovIdentity(uint64(l))
	}
	return n
}
func (m *MsgRegisterIdentity) Unmarshal(dAtA []byte) error {
	return unmarshalStringsIdentity(dAtA, "MsgRegisterIdentity", func(fieldNum int32, val string) {
		switch fieldNum {
		case 1:
			m.Address = val
		case 2:
			m.DocumentHash = val
		case 3:
			m.Nationality = val
		case 4:
			m.Level = val
		}
	})
}

// ---------------------------------------------------------------------------
// MsgRegisterIdentityResponse (empty)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgRegisterIdentityResponse proto.InternalMessageInfo

func (m *MsgRegisterIdentityResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRegisterIdentityResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRegisterIdentityResponse.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgRegisterIdentityResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgRegisterIdentityResponse.Merge(m, src)
}
func (m *MsgRegisterIdentityResponse) XXX_Size() int { return m.Size() }
func (m *MsgRegisterIdentityResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgRegisterIdentityResponse.DiscardUnknown(m)
}
func (m *MsgRegisterIdentityResponse) Marshal() (dAtA []byte, err error)            { return nil, nil }
func (m *MsgRegisterIdentityResponse) MarshalTo(dAtA []byte) (int, error)            { return 0, nil }
func (m *MsgRegisterIdentityResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgRegisterIdentityResponse) Size() (n int)                                 { return 0 }
func (m *MsgRegisterIdentityResponse) Unmarshal(dAtA []byte) error                   { return unmarshalEmptyIdentity(dAtA, "MsgRegisterIdentityResponse") }

// ---------------------------------------------------------------------------
// MsgVerifyIdentity (3 string fields)
// ---------------------------------------------------------------------------

var xxx_messageInfo_MsgVerifyIdentity proto.InternalMessageInfo

func (m *MsgVerifyIdentity) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgVerifyIdentity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgVerifyIdentity.Marshal(b, m, deterministic)
	}
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *MsgVerifyIdentity) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgVerifyIdentity.Merge(m, src) }
func (m *MsgVerifyIdentity) XXX_Size() int               { return m.Size() }
func (m *MsgVerifyIdentity) XXX_DiscardUnknown()         { xxx_messageInfo_MsgVerifyIdentity.DiscardUnknown(m) }

func (m *MsgVerifyIdentity) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}
func (m *MsgVerifyIdentity) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}
func (m *MsgVerifyIdentity) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.Level) > 0 {
		i -= len(m.Level)
		copy(dAtA[i:], m.Level)
		i = encodeVarintIdentity(dAtA, i, uint64(len(m.Level)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintIdentity(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Verifier) > 0 {
		i -= len(m.Verifier)
		copy(dAtA[i:], m.Verifier)
		i = encodeVarintIdentity(dAtA, i, uint64(len(m.Verifier)))
		i--
		dAtA[i] = 0x0a
	}
	return len(dAtA) - i, nil
}
func (m *MsgVerifyIdentity) Size() (n int) {
	if m == nil { return 0 }
	var l int
	l = len(m.Verifier); if l > 0 { n += 1 + l + sovIdentity(uint64(l)) }
	l = len(m.Address); if l > 0 { n += 1 + l + sovIdentity(uint64(l)) }
	l = len(m.Level); if l > 0 { n += 1 + l + sovIdentity(uint64(l)) }
	return n
}
func (m *MsgVerifyIdentity) Unmarshal(dAtA []byte) error {
	return unmarshalStringsIdentity(dAtA, "MsgVerifyIdentity", func(fieldNum int32, val string) {
		switch fieldNum {
		case 1: m.Verifier = val
		case 2: m.Address = val
		case 3: m.Level = val
		}
	})
}

// MsgVerifyIdentityResponse
var xxx_messageInfo_MsgVerifyIdentityResponse proto.InternalMessageInfo
func (m *MsgVerifyIdentityResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgVerifyIdentityResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgVerifyIdentityResponse.Marshal(b, m, deterministic) }
	b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil
}
func (m *MsgVerifyIdentityResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgVerifyIdentityResponse.Merge(m, src) }
func (m *MsgVerifyIdentityResponse) XXX_Size() int { return m.Size() }
func (m *MsgVerifyIdentityResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgVerifyIdentityResponse.DiscardUnknown(m) }
func (m *MsgVerifyIdentityResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgVerifyIdentityResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgVerifyIdentityResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgVerifyIdentityResponse) Size() (n int) { return 0 }
func (m *MsgVerifyIdentityResponse) Unmarshal(dAtA []byte) error { return unmarshalEmptyIdentity(dAtA, "MsgVerifyIdentityResponse") }

// ---------------------------------------------------------------------------
// MsgRejectIdentity (3 string fields: verifier, address, reason)
// ---------------------------------------------------------------------------
var xxx_messageInfo_MsgRejectIdentity proto.InternalMessageInfo
func (m *MsgRejectIdentity) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRejectIdentity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic { return xxx_messageInfo_MsgRejectIdentity.Marshal(b, m, deterministic) }
	b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil
}
func (m *MsgRejectIdentity) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgRejectIdentity.Merge(m, src) }
func (m *MsgRejectIdentity) XXX_Size() int { return m.Size() }
func (m *MsgRejectIdentity) XXX_DiscardUnknown() { xxx_messageInfo_MsgRejectIdentity.DiscardUnknown(m) }
func (m *MsgRejectIdentity) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgRejectIdentity) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgRejectIdentity) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.Reason) > 0 { i -= len(m.Reason); copy(dAtA[i:], m.Reason); i = encodeVarintIdentity(dAtA, i, uint64(len(m.Reason))); i--; dAtA[i] = 0x1a }
	if len(m.Address) > 0 { i -= len(m.Address); copy(dAtA[i:], m.Address); i = encodeVarintIdentity(dAtA, i, uint64(len(m.Address))); i--; dAtA[i] = 0x12 }
	if len(m.Verifier) > 0 { i -= len(m.Verifier); copy(dAtA[i:], m.Verifier); i = encodeVarintIdentity(dAtA, i, uint64(len(m.Verifier))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgRejectIdentity) Size() (n int) {
	if m == nil { return 0 }; var l int
	l = len(m.Verifier); if l > 0 { n += 1 + l + sovIdentity(uint64(l)) }
	l = len(m.Address); if l > 0 { n += 1 + l + sovIdentity(uint64(l)) }
	l = len(m.Reason); if l > 0 { n += 1 + l + sovIdentity(uint64(l)) }
	return n
}
func (m *MsgRejectIdentity) Unmarshal(dAtA []byte) error {
	return unmarshalStringsIdentity(dAtA, "MsgRejectIdentity", func(fieldNum int32, val string) {
		switch fieldNum { case 1: m.Verifier = val; case 2: m.Address = val; case 3: m.Reason = val }
	})
}

// MsgRejectIdentityResponse
var xxx_messageInfo_MsgRejectIdentityResponse proto.InternalMessageInfo
func (m *MsgRejectIdentityResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRejectIdentityResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgRejectIdentityResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgRejectIdentityResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgRejectIdentityResponse.Merge(m, src) }
func (m *MsgRejectIdentityResponse) XXX_Size() int { return m.Size() }
func (m *MsgRejectIdentityResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgRejectIdentityResponse.DiscardUnknown(m) }
func (m *MsgRejectIdentityResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgRejectIdentityResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgRejectIdentityResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgRejectIdentityResponse) Size() (n int) { return 0 }
func (m *MsgRejectIdentityResponse) Unmarshal(dAtA []byte) error { return unmarshalEmptyIdentity(dAtA, "MsgRejectIdentityResponse") }

// MsgRevokeIdentity (3 string fields: authority, address, reason)
var xxx_messageInfo_MsgRevokeIdentity proto.InternalMessageInfo
func (m *MsgRevokeIdentity) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRevokeIdentity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgRevokeIdentity.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgRevokeIdentity) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgRevokeIdentity.Merge(m, src) }
func (m *MsgRevokeIdentity) XXX_Size() int { return m.Size() }
func (m *MsgRevokeIdentity) XXX_DiscardUnknown() { xxx_messageInfo_MsgRevokeIdentity.DiscardUnknown(m) }
func (m *MsgRevokeIdentity) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgRevokeIdentity) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgRevokeIdentity) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.Reason) > 0 { i -= len(m.Reason); copy(dAtA[i:], m.Reason); i = encodeVarintIdentity(dAtA, i, uint64(len(m.Reason))); i--; dAtA[i] = 0x1a }
	if len(m.Address) > 0 { i -= len(m.Address); copy(dAtA[i:], m.Address); i = encodeVarintIdentity(dAtA, i, uint64(len(m.Address))); i--; dAtA[i] = 0x12 }
	if len(m.Authority) > 0 { i -= len(m.Authority); copy(dAtA[i:], m.Authority); i = encodeVarintIdentity(dAtA, i, uint64(len(m.Authority))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgRevokeIdentity) Size() (n int) {
	if m == nil { return 0 }; var l int
	l = len(m.Authority); if l > 0 { n += 1 + l + sovIdentity(uint64(l)) }
	l = len(m.Address); if l > 0 { n += 1 + l + sovIdentity(uint64(l)) }
	l = len(m.Reason); if l > 0 { n += 1 + l + sovIdentity(uint64(l)) }
	return n
}
func (m *MsgRevokeIdentity) Unmarshal(dAtA []byte) error {
	return unmarshalStringsIdentity(dAtA, "MsgRevokeIdentity", func(fieldNum int32, val string) {
		switch fieldNum { case 1: m.Authority = val; case 2: m.Address = val; case 3: m.Reason = val }
	})
}

// MsgRevokeIdentityResponse
var xxx_messageInfo_MsgRevokeIdentityResponse proto.InternalMessageInfo
func (m *MsgRevokeIdentityResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRevokeIdentityResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgRevokeIdentityResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgRevokeIdentityResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgRevokeIdentityResponse.Merge(m, src) }
func (m *MsgRevokeIdentityResponse) XXX_Size() int { return m.Size() }
func (m *MsgRevokeIdentityResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgRevokeIdentityResponse.DiscardUnknown(m) }
func (m *MsgRevokeIdentityResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgRevokeIdentityResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgRevokeIdentityResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgRevokeIdentityResponse) Size() (n int) { return 0 }
func (m *MsgRevokeIdentityResponse) Unmarshal(dAtA []byte) error { return unmarshalEmptyIdentity(dAtA, "MsgRevokeIdentityResponse") }

// MsgUpdateIdentity (3 string fields: address, document_hash, nationality)
var xxx_messageInfo_MsgUpdateIdentity proto.InternalMessageInfo
func (m *MsgUpdateIdentity) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgUpdateIdentity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgUpdateIdentity.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgUpdateIdentity) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgUpdateIdentity.Merge(m, src) }
func (m *MsgUpdateIdentity) XXX_Size() int { return m.Size() }
func (m *MsgUpdateIdentity) XXX_DiscardUnknown() { xxx_messageInfo_MsgUpdateIdentity.DiscardUnknown(m) }
func (m *MsgUpdateIdentity) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgUpdateIdentity) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgUpdateIdentity) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.Nationality) > 0 { i -= len(m.Nationality); copy(dAtA[i:], m.Nationality); i = encodeVarintIdentity(dAtA, i, uint64(len(m.Nationality))); i--; dAtA[i] = 0x1a }
	if len(m.DocumentHash) > 0 { i -= len(m.DocumentHash); copy(dAtA[i:], m.DocumentHash); i = encodeVarintIdentity(dAtA, i, uint64(len(m.DocumentHash))); i--; dAtA[i] = 0x12 }
	if len(m.Address) > 0 { i -= len(m.Address); copy(dAtA[i:], m.Address); i = encodeVarintIdentity(dAtA, i, uint64(len(m.Address))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgUpdateIdentity) Size() (n int) {
	if m == nil { return 0 }; var l int
	l = len(m.Address); if l > 0 { n += 1 + l + sovIdentity(uint64(l)) }
	l = len(m.DocumentHash); if l > 0 { n += 1 + l + sovIdentity(uint64(l)) }
	l = len(m.Nationality); if l > 0 { n += 1 + l + sovIdentity(uint64(l)) }
	return n
}
func (m *MsgUpdateIdentity) Unmarshal(dAtA []byte) error {
	return unmarshalStringsIdentity(dAtA, "MsgUpdateIdentity", func(fieldNum int32, val string) {
		switch fieldNum { case 1: m.Address = val; case 2: m.DocumentHash = val; case 3: m.Nationality = val }
	})
}

// MsgUpdateIdentityResponse
var xxx_messageInfo_MsgUpdateIdentityResponse proto.InternalMessageInfo
func (m *MsgUpdateIdentityResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgUpdateIdentityResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgUpdateIdentityResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgUpdateIdentityResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgUpdateIdentityResponse.Merge(m, src) }
func (m *MsgUpdateIdentityResponse) XXX_Size() int { return m.Size() }
func (m *MsgUpdateIdentityResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgUpdateIdentityResponse.DiscardUnknown(m) }
func (m *MsgUpdateIdentityResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgUpdateIdentityResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgUpdateIdentityResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgUpdateIdentityResponse) Size() (n int) { return 0 }
func (m *MsgUpdateIdentityResponse) Unmarshal(dAtA []byte) error { return unmarshalEmptyIdentity(dAtA, "MsgUpdateIdentityResponse") }

// MsgRegisterVerifier (4 string fields: authority, verifier, name, max_level)
var xxx_messageInfo_MsgRegisterVerifier proto.InternalMessageInfo
func (m *MsgRegisterVerifier) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRegisterVerifier) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgRegisterVerifier.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgRegisterVerifier) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgRegisterVerifier.Merge(m, src) }
func (m *MsgRegisterVerifier) XXX_Size() int { return m.Size() }
func (m *MsgRegisterVerifier) XXX_DiscardUnknown() { xxx_messageInfo_MsgRegisterVerifier.DiscardUnknown(m) }
func (m *MsgRegisterVerifier) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgRegisterVerifier) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgRegisterVerifier) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.MaxLevel) > 0 { i -= len(m.MaxLevel); copy(dAtA[i:], m.MaxLevel); i = encodeVarintIdentity(dAtA, i, uint64(len(m.MaxLevel))); i--; dAtA[i] = 0x22 }
	if len(m.Name) > 0 { i -= len(m.Name); copy(dAtA[i:], m.Name); i = encodeVarintIdentity(dAtA, i, uint64(len(m.Name))); i--; dAtA[i] = 0x1a }
	if len(m.Verifier) > 0 { i -= len(m.Verifier); copy(dAtA[i:], m.Verifier); i = encodeVarintIdentity(dAtA, i, uint64(len(m.Verifier))); i--; dAtA[i] = 0x12 }
	if len(m.Authority) > 0 { i -= len(m.Authority); copy(dAtA[i:], m.Authority); i = encodeVarintIdentity(dAtA, i, uint64(len(m.Authority))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgRegisterVerifier) Size() (n int) {
	if m == nil { return 0 }; var l int
	l = len(m.Authority); if l > 0 { n += 1 + l + sovIdentity(uint64(l)) }
	l = len(m.Verifier); if l > 0 { n += 1 + l + sovIdentity(uint64(l)) }
	l = len(m.Name); if l > 0 { n += 1 + l + sovIdentity(uint64(l)) }
	l = len(m.MaxLevel); if l > 0 { n += 1 + l + sovIdentity(uint64(l)) }
	return n
}
func (m *MsgRegisterVerifier) Unmarshal(dAtA []byte) error {
	return unmarshalStringsIdentity(dAtA, "MsgRegisterVerifier", func(fieldNum int32, val string) {
		switch fieldNum { case 1: m.Authority = val; case 2: m.Verifier = val; case 3: m.Name = val; case 4: m.MaxLevel = val }
	})
}

// MsgRegisterVerifierResponse
var xxx_messageInfo_MsgRegisterVerifierResponse proto.InternalMessageInfo
func (m *MsgRegisterVerifierResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgRegisterVerifierResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgRegisterVerifierResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgRegisterVerifierResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgRegisterVerifierResponse.Merge(m, src) }
func (m *MsgRegisterVerifierResponse) XXX_Size() int { return m.Size() }
func (m *MsgRegisterVerifierResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgRegisterVerifierResponse.DiscardUnknown(m) }
func (m *MsgRegisterVerifierResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgRegisterVerifierResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgRegisterVerifierResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgRegisterVerifierResponse) Size() (n int) { return 0 }
func (m *MsgRegisterVerifierResponse) Unmarshal(dAtA []byte) error { return unmarshalEmptyIdentity(dAtA, "MsgRegisterVerifierResponse") }

// MsgDeactivateVerifier (2 string fields: authority, verifier)
var xxx_messageInfo_MsgDeactivateVerifier proto.InternalMessageInfo
func (m *MsgDeactivateVerifier) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgDeactivateVerifier) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgDeactivateVerifier.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgDeactivateVerifier) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgDeactivateVerifier.Merge(m, src) }
func (m *MsgDeactivateVerifier) XXX_Size() int { return m.Size() }
func (m *MsgDeactivateVerifier) XXX_DiscardUnknown() { xxx_messageInfo_MsgDeactivateVerifier.DiscardUnknown(m) }
func (m *MsgDeactivateVerifier) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgDeactivateVerifier) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgDeactivateVerifier) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.Verifier) > 0 { i -= len(m.Verifier); copy(dAtA[i:], m.Verifier); i = encodeVarintIdentity(dAtA, i, uint64(len(m.Verifier))); i--; dAtA[i] = 0x12 }
	if len(m.Authority) > 0 { i -= len(m.Authority); copy(dAtA[i:], m.Authority); i = encodeVarintIdentity(dAtA, i, uint64(len(m.Authority))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgDeactivateVerifier) Size() (n int) {
	if m == nil { return 0 }; var l int
	l = len(m.Authority); if l > 0 { n += 1 + l + sovIdentity(uint64(l)) }
	l = len(m.Verifier); if l > 0 { n += 1 + l + sovIdentity(uint64(l)) }
	return n
}
func (m *MsgDeactivateVerifier) Unmarshal(dAtA []byte) error {
	return unmarshalStringsIdentity(dAtA, "MsgDeactivateVerifier", func(fieldNum int32, val string) {
		switch fieldNum { case 1: m.Authority = val; case 2: m.Verifier = val }
	})
}

// MsgDeactivateVerifierResponse
var xxx_messageInfo_MsgDeactivateVerifierResponse proto.InternalMessageInfo
func (m *MsgDeactivateVerifierResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgDeactivateVerifierResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgDeactivateVerifierResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgDeactivateVerifierResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgDeactivateVerifierResponse.Merge(m, src) }
func (m *MsgDeactivateVerifierResponse) XXX_Size() int { return m.Size() }
func (m *MsgDeactivateVerifierResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgDeactivateVerifierResponse.DiscardUnknown(m) }
func (m *MsgDeactivateVerifierResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgDeactivateVerifierResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgDeactivateVerifierResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgDeactivateVerifierResponse) Size() (n int) { return 0 }
func (m *MsgDeactivateVerifierResponse) Unmarshal(dAtA []byte) error { return unmarshalEmptyIdentity(dAtA, "MsgDeactivateVerifierResponse") }

// MsgIncrementBookings (2 string fields: authority, address)
var xxx_messageInfo_MsgIncrementBookings proto.InternalMessageInfo
func (m *MsgIncrementBookings) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgIncrementBookings) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgIncrementBookings.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgIncrementBookings) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgIncrementBookings.Merge(m, src) }
func (m *MsgIncrementBookings) XXX_Size() int { return m.Size() }
func (m *MsgIncrementBookings) XXX_DiscardUnknown() { xxx_messageInfo_MsgIncrementBookings.DiscardUnknown(m) }
func (m *MsgIncrementBookings) Marshal() (dAtA []byte, err error) { size := m.Size(); dAtA = make([]byte, size); n, err := m.MarshalToSizedBuffer(dAtA[:size]); if err != nil { return nil, err }; return dAtA[:n], nil }
func (m *MsgIncrementBookings) MarshalTo(dAtA []byte) (int, error) { size := m.Size(); return m.MarshalToSizedBuffer(dAtA[:size]) }
func (m *MsgIncrementBookings) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if len(m.Address) > 0 { i -= len(m.Address); copy(dAtA[i:], m.Address); i = encodeVarintIdentity(dAtA, i, uint64(len(m.Address))); i--; dAtA[i] = 0x12 }
	if len(m.Authority) > 0 { i -= len(m.Authority); copy(dAtA[i:], m.Authority); i = encodeVarintIdentity(dAtA, i, uint64(len(m.Authority))); i--; dAtA[i] = 0x0a }
	return len(dAtA) - i, nil
}
func (m *MsgIncrementBookings) Size() (n int) {
	if m == nil { return 0 }; var l int
	l = len(m.Authority); if l > 0 { n += 1 + l + sovIdentity(uint64(l)) }
	l = len(m.Address); if l > 0 { n += 1 + l + sovIdentity(uint64(l)) }
	return n
}
func (m *MsgIncrementBookings) Unmarshal(dAtA []byte) error {
	return unmarshalStringsIdentity(dAtA, "MsgIncrementBookings", func(fieldNum int32, val string) {
		switch fieldNum { case 1: m.Authority = val; case 2: m.Address = val }
	})
}

// MsgIncrementBookingsResponse
var xxx_messageInfo_MsgIncrementBookingsResponse proto.InternalMessageInfo
func (m *MsgIncrementBookingsResponse) XXX_Unmarshal(b []byte) error { return m.Unmarshal(b) }
func (m *MsgIncrementBookingsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) { if deterministic { return xxx_messageInfo_MsgIncrementBookingsResponse.Marshal(b, m, deterministic) }; b = b[:cap(b)]; n, err := m.MarshalToSizedBuffer(b); if err != nil { return nil, err }; return b[:n], nil }
func (m *MsgIncrementBookingsResponse) XXX_Merge(src proto.Message) { xxx_messageInfo_MsgIncrementBookingsResponse.Merge(m, src) }
func (m *MsgIncrementBookingsResponse) XXX_Size() int { return m.Size() }
func (m *MsgIncrementBookingsResponse) XXX_DiscardUnknown() { xxx_messageInfo_MsgIncrementBookingsResponse.DiscardUnknown(m) }
func (m *MsgIncrementBookingsResponse) Marshal() (dAtA []byte, err error) { return nil, nil }
func (m *MsgIncrementBookingsResponse) MarshalTo(dAtA []byte) (int, error) { return 0, nil }
func (m *MsgIncrementBookingsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) { return len(dAtA), nil }
func (m *MsgIncrementBookingsResponse) Size() (n int) { return 0 }
func (m *MsgIncrementBookingsResponse) Unmarshal(dAtA []byte) error { return unmarshalEmptyIdentity(dAtA, "MsgIncrementBookingsResponse") }

// ---------------------------------------------------------------------------
// Helper functions (identity-specific names to avoid conflicts)
// ---------------------------------------------------------------------------

var (
	ErrIntOverflowIdentity  = fmt.Errorf("proto: integer overflow")
	ErrInvalidLengthIdentity = fmt.Errorf("proto: negative length found during unmarshaling")
)

func encodeVarintIdentity(dAtA []byte, offset int, v uint64) int {
	offset -= sovIdentity(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}

func sovIdentity(x uint64) int {
	return (bits.Len64(x|1) + 6) / 7
}

func skipIdentity(dAtA []byte) (int, error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 { return 0, ErrIntOverflowIdentity }
			if iNdEx >= l { return 0, io.ErrUnexpectedEOF }
			b := dAtA[iNdEx]; iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 { break }
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return 0, ErrIntOverflowIdentity }
				if iNdEx >= l { return 0, io.ErrUnexpectedEOF }
				iNdEx++; if dAtA[iNdEx-1] < 0x80 { break }
			}
		case 1: iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return 0, ErrIntOverflowIdentity }
				if iNdEx >= l { return 0, io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 { break }
			}
			if length < 0 { return 0, ErrInvalidLengthIdentity }
			iNdEx += length
		case 3: depth++
		case 4:
			if depth == 0 { return 0, fmt.Errorf("proto: wiretype end group for non-group") }
			depth--
		case 5: iNdEx += 4
		default: return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 { return 0, ErrInvalidLengthIdentity }
		if depth == 0 { return iNdEx, nil }
	}
	return 0, io.ErrUnexpectedEOF
}

func unmarshalEmptyIdentity(dAtA []byte, msgName string) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 { return ErrIntOverflowIdentity }
			if iNdEx >= l { return io.ErrUnexpectedEOF }
			b := dAtA[iNdEx]; iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 { break }
		}
		fieldNum := int32(wire >> 3)
		if fieldNum <= 0 { return fmt.Errorf("proto: %s: illegal tag %d", msgName, fieldNum) }
		iNdEx = preIndex
		skippy, err := skipIdentity(dAtA[iNdEx:])
		if err != nil { return err }
		if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthIdentity }
		if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }
		iNdEx += skippy
	}
	return nil
}

// unmarshalStringsIdentity is a generic unmarshaler for messages with only string fields
func unmarshalStringsIdentity(dAtA []byte, msgName string, setter func(fieldNum int32, val string)) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 { return ErrIntOverflowIdentity }
			if iNdEx >= l { return io.ErrUnexpectedEOF }
			b := dAtA[iNdEx]; iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 { break }
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 { return fmt.Errorf("proto: %s: wiretype end group for non-group", msgName) }
		if fieldNum <= 0 { return fmt.Errorf("proto: %s: illegal tag %d (wire type %d)", msgName, fieldNum, wire) }
		if wireType == 2 {
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 { return ErrIntOverflowIdentity }
				if iNdEx >= l { return io.ErrUnexpectedEOF }
				b := dAtA[iNdEx]; iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 { break }
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 { return ErrInvalidLengthIdentity }
			postIndex := iNdEx + intStringLen
			if postIndex < 0 { return ErrInvalidLengthIdentity }
			if postIndex > l { return io.ErrUnexpectedEOF }
			setter(fieldNum, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		} else {
			iNdEx = preIndex
			skippy, err := skipIdentity(dAtA[iNdEx:])
			if err != nil { return err }
			if (skippy < 0) || (iNdEx+skippy) < 0 { return ErrInvalidLengthIdentity }
			if (iNdEx + skippy) > l { return io.ErrUnexpectedEOF }
			iNdEx += skippy
		}
	}
	if iNdEx > l { return io.ErrUnexpectedEOF }
	return nil
}
