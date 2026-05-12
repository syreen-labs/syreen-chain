package types

import (
	"bytes"
	"compress/gzip"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

func init() {
	registerLaunchpadProtoFileDescriptors()
}

func registerLaunchpadProtoFileDescriptors() {
	label := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	typeString := descriptorpb.FieldDescriptorProto_TYPE_STRING
	typeBytes := descriptorpb.FieldDescriptorProto_TYPE_BYTES
	typeUint64 := descriptorpb.FieldDescriptorProto_TYPE_UINT64
	typeInt64 := descriptorpb.FieldDescriptorProto_TYPE_INT64

	fd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/launchpad/tx.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.launchpad"),
		MessageType: []*descriptorpb.DescriptorProto{
			{ // 0: MsgCreateLaunch
				Name: sp("MsgCreateLaunch"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("creator"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("creator")},
					{Name: sp("token_denom"), Number: ip(2), Label: &label, Type: &typeString, JsonName: sp("tokenDenom")},
					{Name: sp("token_supply"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("tokenSupply")},
					{Name: sp("price_per_token"), Number: ip(4), Label: &label, Type: &typeBytes, JsonName: sp("pricePerToken")},
					{Name: sp("quote_denom"), Number: ip(5), Label: &label, Type: &typeString, JsonName: sp("quoteDenom")},
					{Name: sp("soft_cap"), Number: ip(6), Label: &label, Type: &typeBytes, JsonName: sp("softCap")},
					{Name: sp("hard_cap"), Number: ip(7), Label: &label, Type: &typeBytes, JsonName: sp("hardCap")},
					{Name: sp("max_per_wallet"), Number: ip(8), Label: &label, Type: &typeBytes, JsonName: sp("maxPerWallet")},
					{Name: sp("start_block"), Number: ip(9), Label: &label, Type: &typeInt64, JsonName: sp("startBlock")},
					{Name: sp("end_block"), Number: ip(10), Label: &label, Type: &typeInt64, JsonName: sp("endBlock")},
					{Name: sp("vesting_blocks"), Number: ip(11), Label: &label, Type: &typeInt64, JsonName: sp("vestingBlocks")},
					{Name: sp("tge_percent"), Number: ip(12), Label: &label, Type: &typeUint64, JsonName: sp("tgePercent")},
				},
			},
			{ // 1: MsgCreateLaunchResponse
				Name: sp("MsgCreateLaunchResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("launch_id"), Number: ip(1), Label: &label, Type: &typeUint64, JsonName: sp("launchId")},
				},
			},
			{ // 2: MsgContribute
				Name: sp("MsgContribute"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
					{Name: sp("launch_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("launchId")},
					{Name: sp("amount"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("amount")},
				},
			},
			{ // 3: MsgContributeResponse
				Name: sp("MsgContributeResponse"),
			},
			{ // 4: MsgClaimTokens
				Name: sp("MsgClaimTokens"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
					{Name: sp("launch_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("launchId")},
				},
			},
			{ // 5: MsgClaimTokensResponse
				Name: sp("MsgClaimTokensResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("amount"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("amount")},
				},
			},
			{ // 6: MsgClaimRefund
				Name: sp("MsgClaimRefund"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
					{Name: sp("launch_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("launchId")},
				},
			},
			{ // 7: MsgClaimRefundResponse
				Name: sp("MsgClaimRefundResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("amount"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("amount")},
				},
			},
			{ // 8: MsgFinalizeLaunch
				Name: sp("MsgFinalizeLaunch"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("authority"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("authority")},
					{Name: sp("launch_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("launchId")},
				},
			},
			{ // 9: MsgFinalizeLaunchResponse
				Name: sp("MsgFinalizeLaunchResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("status"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("status")},
				},
			},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: sp("CreateLaunch"), InputType: sp(".syreen.launchpad.MsgCreateLaunch"), OutputType: sp(".syreen.launchpad.MsgCreateLaunchResponse")},
					{Name: sp("Contribute"), InputType: sp(".syreen.launchpad.MsgContribute"), OutputType: sp(".syreen.launchpad.MsgContributeResponse")},
					{Name: sp("ClaimTokens"), InputType: sp(".syreen.launchpad.MsgClaimTokens"), OutputType: sp(".syreen.launchpad.MsgClaimTokensResponse")},
					{Name: sp("ClaimRefund"), InputType: sp(".syreen.launchpad.MsgClaimRefund"), OutputType: sp(".syreen.launchpad.MsgClaimRefundResponse")},
					{Name: sp("FinalizeLaunch"), InputType: sp(".syreen.launchpad.MsgFinalizeLaunch"), OutputType: sp(".syreen.launchpad.MsgFinalizeLaunchResponse")},
				},
			},
		},
	}

	rawBz, err := proto.Marshal(fd)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	_, _ = w.Write(rawBz)
	_ = w.Close()
	fileDescriptorLaunchpadTx = buf.Bytes()

	file, err := protodesc.NewFile(fd, nil)
	if err != nil {
		panic("launchpad: failed to create proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(file); err != nil {
		// Already registered is fine
	}
}

func sp(s string) *string { return &s }
func ip(n int32) *int32   { return &n }
