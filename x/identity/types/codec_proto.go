package types

import (
	"bytes"
	"compress/gzip"
	"sync"

	gogoproto "github.com/cosmos/gogoproto/proto"
	"google.golang.org/protobuf/proto"
)

// This file registers all identity types with gogoproto and provides
// Descriptor() methods so the Cosmos SDK can properly encode/decode transactions.

var (
	fileDescOnce sync.Once
	fileDescGzip []byte
)

func txFileDescGzip() []byte {
	fileDescOnce.Do(func() {
		if TxFileDescProto == nil {
			return
		}
		raw, err := proto.Marshal(TxFileDescProto)
		if err != nil {
			return
		}
		var buf bytes.Buffer
		w := gzip.NewWriter(&buf)
		w.Write(raw)
		w.Close()
		fileDescGzip = buf.Bytes()
	})
	return fileDescGzip
}

func init() {
	gogoproto.RegisterType((*MsgRegisterIdentity)(nil), "syreen.identity.MsgRegisterIdentity")
	gogoproto.RegisterType((*MsgRegisterIdentityResponse)(nil), "syreen.identity.MsgRegisterIdentityResponse")
	gogoproto.RegisterType((*MsgVerifyIdentity)(nil), "syreen.identity.MsgVerifyIdentity")
	gogoproto.RegisterType((*MsgVerifyIdentityResponse)(nil), "syreen.identity.MsgVerifyIdentityResponse")
	gogoproto.RegisterType((*MsgRejectIdentity)(nil), "syreen.identity.MsgRejectIdentity")
	gogoproto.RegisterType((*MsgRejectIdentityResponse)(nil), "syreen.identity.MsgRejectIdentityResponse")
	gogoproto.RegisterType((*MsgRevokeIdentity)(nil), "syreen.identity.MsgRevokeIdentity")
	gogoproto.RegisterType((*MsgRevokeIdentityResponse)(nil), "syreen.identity.MsgRevokeIdentityResponse")
	gogoproto.RegisterType((*MsgUpdateIdentity)(nil), "syreen.identity.MsgUpdateIdentity")
	gogoproto.RegisterType((*MsgUpdateIdentityResponse)(nil), "syreen.identity.MsgUpdateIdentityResponse")
	gogoproto.RegisterType((*MsgRegisterVerifier)(nil), "syreen.identity.MsgRegisterVerifier")
	gogoproto.RegisterType((*MsgRegisterVerifierResponse)(nil), "syreen.identity.MsgRegisterVerifierResponse")
	gogoproto.RegisterType((*MsgDeactivateVerifier)(nil), "syreen.identity.MsgDeactivateVerifier")
	gogoproto.RegisterType((*MsgDeactivateVerifierResponse)(nil), "syreen.identity.MsgDeactivateVerifierResponse")
	gogoproto.RegisterType((*MsgIncrementBookings)(nil), "syreen.identity.MsgIncrementBookings")
	gogoproto.RegisterType((*MsgIncrementBookingsResponse)(nil), "syreen.identity.MsgIncrementBookingsResponse")
	gogoproto.RegisterType((*Identity)(nil), "syreen.identity.Identity")
	gogoproto.RegisterType((*Verifier)(nil), "syreen.identity.Verifier")
	gogoproto.RegisterType((*Params)(nil), "syreen.identity.Params")
}

// Descriptor returns the gzipped proto file descriptor and path to this message.
// Required by gogoproto for proper Any encoding in Cosmos SDK transactions.

func (*MsgRegisterIdentity) Descriptor() ([]byte, []int)          { return txFileDescGzip(), []int{0} }
func (*MsgRegisterIdentityResponse) Descriptor() ([]byte, []int)  { return txFileDescGzip(), []int{1} }
func (*MsgVerifyIdentity) Descriptor() ([]byte, []int)            { return txFileDescGzip(), []int{2} }
func (*MsgVerifyIdentityResponse) Descriptor() ([]byte, []int)    { return txFileDescGzip(), []int{3} }
func (*MsgRejectIdentity) Descriptor() ([]byte, []int)            { return txFileDescGzip(), []int{4} }
func (*MsgRejectIdentityResponse) Descriptor() ([]byte, []int)    { return txFileDescGzip(), []int{5} }
func (*MsgRevokeIdentity) Descriptor() ([]byte, []int)            { return txFileDescGzip(), []int{6} }
func (*MsgRevokeIdentityResponse) Descriptor() ([]byte, []int)    { return txFileDescGzip(), []int{7} }
func (*MsgUpdateIdentity) Descriptor() ([]byte, []int)            { return txFileDescGzip(), []int{8} }
func (*MsgUpdateIdentityResponse) Descriptor() ([]byte, []int)    { return txFileDescGzip(), []int{9} }
func (*MsgRegisterVerifier) Descriptor() ([]byte, []int)          { return txFileDescGzip(), []int{10} }
func (*MsgRegisterVerifierResponse) Descriptor() ([]byte, []int)  { return txFileDescGzip(), []int{11} }
func (*MsgDeactivateVerifier) Descriptor() ([]byte, []int)        { return txFileDescGzip(), []int{12} }
func (*MsgDeactivateVerifierResponse) Descriptor() ([]byte, []int) { return txFileDescGzip(), []int{13} }
func (*MsgIncrementBookings) Descriptor() ([]byte, []int)         { return txFileDescGzip(), []int{14} }
func (*MsgIncrementBookingsResponse) Descriptor() ([]byte, []int) { return txFileDescGzip(), []int{15} }
func (*Identity) Descriptor() ([]byte, []int)                     { return txFileDescGzip(), []int{16} }
func (*Verifier) Descriptor() ([]byte, []int)                     { return txFileDescGzip(), []int{17} }
func (*Params) Descriptor() ([]byte, []int)                       { return txFileDescGzip(), []int{18} }
