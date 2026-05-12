package types

import (
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

// combinedResolver tries a local file first, then falls back to global registry.
type combinedResolver struct {
	local protoreflect.FileDescriptor
}

func (r combinedResolver) FindFileByPath(path string) (protoreflect.FileDescriptor, error) {
	if r.local != nil && r.local.Path() == path {
		return r.local, nil
	}
	return protoregistry.GlobalFiles.FindFileByPath(path)
}

func (r combinedResolver) FindDescriptorByName(name protoreflect.FullName) (protoreflect.Descriptor, error) {
	if r.local != nil {
		if d := r.local.Messages().ByName(name.Name()); d != nil {
			return d, nil
		}
	}
	return protoregistry.GlobalFiles.FindDescriptorByName(name)
}

// TxFileDescProto holds the raw FileDescriptorProto for codec_proto.go to use.
var TxFileDescProto *descriptorpb.FileDescriptorProto

func init() {
	registerProtoFileDescriptors()
}

func registerProtoFileDescriptors() {
	strType := descriptorpb.FieldDescriptorProto_TYPE_STRING
	uint64Type := descriptorpb.FieldDescriptorProto_TYPE_UINT64
	int64Type := descriptorpb.FieldDescriptorProto_TYPE_INT64
	boolType := descriptorpb.FieldDescriptorProto_TYPE_BOOL
	msgType := descriptorpb.FieldDescriptorProto_TYPE_MESSAGE
	labelOpt := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL

	fd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/identity/tx.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.identity"),
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: sp("MsgRegisterIdentity"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("address"), Number: ip(1), Type: &strType, Label: &labelOpt},
					{Name: sp("document_hash"), Number: ip(2), Type: &strType, Label: &labelOpt},
					{Name: sp("nationality"), Number: ip(3), Type: &strType, Label: &labelOpt},
					{Name: sp("level"), Number: ip(4), Type: &strType, Label: &labelOpt},
				},
			},
			{Name: sp("MsgRegisterIdentityResponse")},
			{
				Name: sp("MsgVerifyIdentity"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("verifier"), Number: ip(1), Type: &strType, Label: &labelOpt},
					{Name: sp("address"), Number: ip(2), Type: &strType, Label: &labelOpt},
					{Name: sp("level"), Number: ip(3), Type: &strType, Label: &labelOpt},
				},
			},
			{Name: sp("MsgVerifyIdentityResponse")},
			{
				Name: sp("MsgRejectIdentity"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("verifier"), Number: ip(1), Type: &strType, Label: &labelOpt},
					{Name: sp("address"), Number: ip(2), Type: &strType, Label: &labelOpt},
					{Name: sp("reason"), Number: ip(3), Type: &strType, Label: &labelOpt},
				},
			},
			{Name: sp("MsgRejectIdentityResponse")},
			{
				Name: sp("MsgRevokeIdentity"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("authority"), Number: ip(1), Type: &strType, Label: &labelOpt},
					{Name: sp("address"), Number: ip(2), Type: &strType, Label: &labelOpt},
					{Name: sp("reason"), Number: ip(3), Type: &strType, Label: &labelOpt},
				},
			},
			{Name: sp("MsgRevokeIdentityResponse")},
			{
				Name: sp("MsgUpdateIdentity"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("address"), Number: ip(1), Type: &strType, Label: &labelOpt},
					{Name: sp("document_hash"), Number: ip(2), Type: &strType, Label: &labelOpt},
					{Name: sp("nationality"), Number: ip(3), Type: &strType, Label: &labelOpt},
				},
			},
			{Name: sp("MsgUpdateIdentityResponse")},
			{
				Name: sp("MsgRegisterVerifier"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("authority"), Number: ip(1), Type: &strType, Label: &labelOpt},
					{Name: sp("verifier"), Number: ip(2), Type: &strType, Label: &labelOpt},
					{Name: sp("name"), Number: ip(3), Type: &strType, Label: &labelOpt},
					{Name: sp("max_level"), Number: ip(4), Type: &strType, Label: &labelOpt},
				},
			},
			{Name: sp("MsgRegisterVerifierResponse")},
			{
				Name: sp("MsgDeactivateVerifier"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("authority"), Number: ip(1), Type: &strType, Label: &labelOpt},
					{Name: sp("verifier"), Number: ip(2), Type: &strType, Label: &labelOpt},
				},
			},
			{Name: sp("MsgDeactivateVerifierResponse")},
			{
				Name: sp("MsgIncrementBookings"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("authority"), Number: ip(1), Type: &strType, Label: &labelOpt},
					{Name: sp("address"), Number: ip(2), Type: &strType, Label: &labelOpt},
				},
			},
			{Name: sp("MsgIncrementBookingsResponse")},
			{
				Name: sp("Identity"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("address"), Number: ip(1), Type: &strType, Label: &labelOpt},
					{Name: sp("level"), Number: ip(2), Type: &strType, Label: &labelOpt},
					{Name: sp("status"), Number: ip(3), Type: &strType, Label: &labelOpt},
					{Name: sp("verifier"), Number: ip(4), Type: &strType, Label: &labelOpt},
					{Name: sp("document_hash"), Number: ip(5), Type: &strType, Label: &labelOpt},
					{Name: sp("nationality"), Number: ip(6), Type: &strType, Label: &labelOpt},
					{Name: sp("verified_at_block"), Number: ip(7), Type: &int64Type, Label: &labelOpt},
					{Name: sp("expires_at_block"), Number: ip(8), Type: &int64Type, Label: &labelOpt},
					{Name: sp("created_at_block"), Number: ip(9), Type: &int64Type, Label: &labelOpt},
					{Name: sp("rejection_reason"), Number: ip(10), Type: &strType, Label: &labelOpt},
					{Name: sp("total_bookings"), Number: ip(11), Type: &uint64Type, Label: &labelOpt},
					{Name: sp("trust_score"), Number: ip(12), Type: &uint64Type, Label: &labelOpt},
				},
			},
			{
				Name: sp("Verifier"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("address"), Number: ip(1), Type: &strType, Label: &labelOpt},
					{Name: sp("name"), Number: ip(2), Type: &strType, Label: &labelOpt},
					{Name: sp("max_level"), Number: ip(3), Type: &strType, Label: &labelOpt},
					{Name: sp("active"), Number: ip(4), Type: &boolType, Label: &labelOpt},
					{Name: sp("total_verified"), Number: ip(5), Type: &uint64Type, Label: &labelOpt},
					{Name: sp("registered_at"), Number: ip(6), Type: &int64Type, Label: &labelOpt},
				},
			},
			{
				Name: sp("Params"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("verification_expiry_blocks"), Number: ip(1), Type: &uint64Type, Label: &labelOpt},
					{Name: sp("min_trust_score"), Number: ip(2), Type: &uint64Type, Label: &labelOpt},
					{Name: sp("enable_identity"), Number: ip(3), Type: &boolType, Label: &labelOpt},
					{Name: sp("basic_verification_free"), Number: ip(4), Type: &boolType, Label: &labelOpt},
				},
			},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: sp("RegisterIdentity"), InputType: sp(".syreen.identity.MsgRegisterIdentity"), OutputType: sp(".syreen.identity.MsgRegisterIdentityResponse")},
					{Name: sp("VerifyIdentity"), InputType: sp(".syreen.identity.MsgVerifyIdentity"), OutputType: sp(".syreen.identity.MsgVerifyIdentityResponse")},
					{Name: sp("RejectIdentity"), InputType: sp(".syreen.identity.MsgRejectIdentity"), OutputType: sp(".syreen.identity.MsgRejectIdentityResponse")},
					{Name: sp("RevokeIdentity"), InputType: sp(".syreen.identity.MsgRevokeIdentity"), OutputType: sp(".syreen.identity.MsgRevokeIdentityResponse")},
					{Name: sp("UpdateIdentity"), InputType: sp(".syreen.identity.MsgUpdateIdentity"), OutputType: sp(".syreen.identity.MsgUpdateIdentityResponse")},
					{Name: sp("RegisterVerifier"), InputType: sp(".syreen.identity.MsgRegisterVerifier"), OutputType: sp(".syreen.identity.MsgRegisterVerifierResponse")},
					{Name: sp("DeactivateVerifier"), InputType: sp(".syreen.identity.MsgDeactivateVerifier"), OutputType: sp(".syreen.identity.MsgDeactivateVerifierResponse")},
					{Name: sp("IncrementBookings"), InputType: sp(".syreen.identity.MsgIncrementBookings"), OutputType: sp(".syreen.identity.MsgIncrementBookingsResponse")},
				},
			},
		},
	}

	// Query file descriptor
	qfd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/identity/query.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.identity"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: sp("QueryParamsRequest")},
			{
				Name: sp("QueryParamsResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("params"), Number: ip(1), Type: &msgType, TypeName: sp(".syreen.identity.Params"), Label: &labelOpt},
				},
			},
			{
				Name: sp("QueryIdentityRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("address"), Number: ip(1), Type: &strType, Label: &labelOpt},
				},
			},
			{
				Name: sp("QueryIdentityResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("identity"), Number: ip(1), Type: &msgType, TypeName: sp(".syreen.identity.Identity"), Label: &labelOpt},
					{Name: sp("found"), Number: ip(2), Type: &boolType, Label: &labelOpt},
				},
			},
			{
				Name: sp("QueryIdentitiesByVerifierRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("verifier"), Number: ip(1), Type: &strType, Label: &labelOpt},
				},
			},
			{
				Name: sp("QueryIdentitiesByVerifierResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("identities"), Number: ip(1), Type: &msgType, TypeName: sp(".syreen.identity.Identity"), Label: &labelOpt},
				},
			},
			{
				Name: sp("QueryVerifierRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("address"), Number: ip(1), Type: &strType, Label: &labelOpt},
				},
			},
			{
				Name: sp("QueryVerifierResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("verifier"), Number: ip(1), Type: &msgType, TypeName: sp(".syreen.identity.Verifier"), Label: &labelOpt},
					{Name: sp("found"), Number: ip(2), Type: &boolType, Label: &labelOpt},
				},
			},
			{
				Name: sp("QueryVerificationStatusRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("address"), Number: ip(1), Type: &strType, Label: &labelOpt},
				},
			},
			{
				Name: sp("QueryVerificationStatusResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("address"), Number: ip(1), Type: &strType, Label: &labelOpt},
					{Name: sp("verified"), Number: ip(2), Type: &boolType, Label: &labelOpt},
					{Name: sp("level"), Number: ip(3), Type: &strType, Label: &labelOpt},
					{Name: sp("status"), Number: ip(4), Type: &strType, Label: &labelOpt},
					{Name: sp("trust_score"), Number: ip(5), Type: &uint64Type, Label: &labelOpt},
				},
			},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Query"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: sp("Params"), InputType: sp(".syreen.identity.QueryParamsRequest"), OutputType: sp(".syreen.identity.QueryParamsResponse")},
					{Name: sp("Identity"), InputType: sp(".syreen.identity.QueryIdentityRequest"), OutputType: sp(".syreen.identity.QueryIdentityResponse")},
					{Name: sp("IdentitiesByVerifier"), InputType: sp(".syreen.identity.QueryIdentitiesByVerifierRequest"), OutputType: sp(".syreen.identity.QueryIdentitiesByVerifierResponse")},
					{Name: sp("Verifier"), InputType: sp(".syreen.identity.QueryVerifierRequest"), OutputType: sp(".syreen.identity.QueryVerifierResponse")},
					{Name: sp("VerificationStatus"), InputType: sp(".syreen.identity.QueryVerificationStatusRequest"), OutputType: sp(".syreen.identity.QueryVerificationStatusResponse")},
				},
			},
		},
	}

	TxFileDescProto = fd

	file, err := protodesc.NewFile(fd, protoregistry.GlobalFiles)
	if err != nil {
		// Message type refs not yet registered - convert to bytes as fallback
		for _, mt := range fd.MessageType {
			for _, f := range mt.Field {
				if f.TypeName != nil && *f.Type == msgType {
					bytesT := descriptorpb.FieldDescriptorProto_TYPE_BYTES
					f.Type = &bytesT
					f.TypeName = nil
				}
			}
		}
		file, err = protodesc.NewFile(fd, nil)
		if err != nil {
			panic("identity: failed to create proto file descriptor: " + err.Error())
		}
	}
	if err := protoregistry.GlobalFiles.RegisterFile(file); err != nil {
		// Already registered is fine
	}

	// Same fallback for query file
	for _, mt := range qfd.MessageType {
		for _, f := range mt.Field {
			if f.TypeName != nil && *f.Type == msgType {
				tn := *f.TypeName
				if tn != ".syreen.identity.Identity" &&
					tn != ".syreen.identity.Verifier" &&
					tn != ".syreen.identity.Params" {
					bytesT := descriptorpb.FieldDescriptorProto_TYPE_BYTES
					f.Type = &bytesT
					f.TypeName = nil
				}
			}
		}
	}

	qfile, err := protodesc.NewFile(qfd, combinedResolver{file})
	if err != nil {
		// Last resort: strip all message type refs
		for _, mt := range qfd.MessageType {
			for _, f := range mt.Field {
				if f.TypeName != nil && *f.Type == msgType {
					bytesT := descriptorpb.FieldDescriptorProto_TYPE_BYTES
					f.Type = &bytesT
					f.TypeName = nil
				}
			}
		}
		qfile, err = protodesc.NewFile(qfd, nil)
		if err != nil {
			panic("identity: failed to create query proto file descriptor: " + err.Error())
		}
	}
	if err := protoregistry.GlobalFiles.RegisterFile(qfile); err != nil {
		// Already registered is fine
	}
}

func sp(s string) *string { return &s }
func ip(i int32) *int32   { return &i }
