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
	registerFarmingProtoFileDescriptors()
}

func registerFarmingProtoFileDescriptors() {
	label := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	typeString := descriptorpb.FieldDescriptorProto_TYPE_STRING
	typeBytes := descriptorpb.FieldDescriptorProto_TYPE_BYTES
	typeUint64 := descriptorpb.FieldDescriptorProto_TYPE_UINT64
	typeInt64 := descriptorpb.FieldDescriptorProto_TYPE_INT64
	typeBool := descriptorpb.FieldDescriptorProto_TYPE_BOOL

	fd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/farming/tx.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.farming"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: sp("MsgStake"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("pool_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("poolId")},
				{Name: sp("amount"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("amount")},
			}},
			{Name: sp("MsgStakeResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("pending_reward"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("pendingReward")},
			}},
			{Name: sp("MsgUnstake"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("pool_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("poolId")},
				{Name: sp("amount"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("amount")},
			}},
			{Name: sp("MsgUnstakeResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("claimed_reward"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("claimedReward")},
			}},
			{Name: sp("MsgClaimReward"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("pool_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("poolId")},
			}},
			{Name: sp("MsgClaimRewardResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("amount"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("amount")},
			}},
			{Name: sp("MsgCreateFarm"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("authority"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("authority")},
				{Name: sp("pool_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("poolId")},
				{Name: sp("reward_per_block"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("rewardPerBlock")},
				{Name: sp("start_block"), Number: ip(4), Label: &label, Type: &typeInt64, JsonName: sp("startBlock")},
				{Name: sp("end_block"), Number: ip(5), Label: &label, Type: &typeInt64, JsonName: sp("endBlock")},
			}},
			{Name: sp("MsgCreateFarmResponse")},
			{Name: sp("MsgUpdateFarm"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("authority"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("authority")},
				{Name: sp("pool_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("poolId")},
				{Name: sp("reward_per_block"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("rewardPerBlock")},
				{Name: sp("active"), Number: ip(4), Label: &label, Type: &typeBool, JsonName: sp("active")},
			}},
			{Name: sp("MsgUpdateFarmResponse")},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: sp("Stake"), InputType: sp(".syreen.farming.MsgStake"), OutputType: sp(".syreen.farming.MsgStakeResponse")},
					{Name: sp("Unstake"), InputType: sp(".syreen.farming.MsgUnstake"), OutputType: sp(".syreen.farming.MsgUnstakeResponse")},
					{Name: sp("ClaimReward"), InputType: sp(".syreen.farming.MsgClaimReward"), OutputType: sp(".syreen.farming.MsgClaimRewardResponse")},
					{Name: sp("CreateFarm"), InputType: sp(".syreen.farming.MsgCreateFarm"), OutputType: sp(".syreen.farming.MsgCreateFarmResponse")},
					{Name: sp("UpdateFarm"), InputType: sp(".syreen.farming.MsgUpdateFarm"), OutputType: sp(".syreen.farming.MsgUpdateFarmResponse")},
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
	fileDescriptorFarmingTx = buf.Bytes()

	file, err := protodesc.NewFile(fd, nil)
	if err != nil {
		panic("farming: failed to create proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(file); err != nil {
		// Already registered is fine
	}
}

func sp(s string) *string { return &s }
func ip(n int32) *int32   { return &n }
