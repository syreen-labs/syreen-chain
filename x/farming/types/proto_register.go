package types

import (
	"bytes"
	"compress/gzip"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

var fileDescriptorFarmingTx []byte

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
		Name:    strp("syreen/farming/tx.proto"),
		Syntax:  strp("proto3"),
		Package: strp("syreen.farming"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: strp("MsgStake"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
				{Name: strp("pool_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
				{Name: strp("amount"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: strp("amount")},
			}},
			{Name: strp("MsgStakeResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("pending_reward"), Number: int32p(1), Label: &label, Type: &typeBytes, JsonName: strp("pendingReward")},
			}},
			{Name: strp("MsgUnstake"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
				{Name: strp("pool_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
				{Name: strp("amount"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: strp("amount")},
			}},
			{Name: strp("MsgUnstakeResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("claimed_reward"), Number: int32p(1), Label: &label, Type: &typeBytes, JsonName: strp("claimedReward")},
			}},
			{Name: strp("MsgClaimReward"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
				{Name: strp("pool_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
			}},
			{Name: strp("MsgClaimRewardResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("amount"), Number: int32p(1), Label: &label, Type: &typeBytes, JsonName: strp("amount")},
			}},
			{Name: strp("MsgCreateFarm"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("authority"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("authority")},
				{Name: strp("pool_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
				{Name: strp("lp_denom"), Number: int32p(3), Label: &label, Type: &typeString, JsonName: strp("lpDenom")},
				{Name: strp("reward_per_block"), Number: int32p(4), Label: &label, Type: &typeBytes, JsonName: strp("rewardPerBlock")},
				{Name: strp("start_block"), Number: int32p(5), Label: &label, Type: &typeInt64, JsonName: strp("startBlock")},
				{Name: strp("end_block"), Number: int32p(6), Label: &label, Type: &typeInt64, JsonName: strp("endBlock")},
			}},
			{Name: strp("MsgCreateFarmResponse")},
			{Name: strp("MsgUpdateFarm"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: strp("authority"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("authority")},
				{Name: strp("pool_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
				{Name: strp("reward_per_block"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: strp("rewardPerBlock")},
				{Name: strp("active"), Number: int32p(4), Label: &label, Type: &typeBool, JsonName: strp("active")},
			}},
			{Name: strp("MsgUpdateFarmResponse")},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: strp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: strp("Stake"), InputType: strp(".syreen.farming.MsgStake"), OutputType: strp(".syreen.farming.MsgStakeResponse")},
					{Name: strp("Unstake"), InputType: strp(".syreen.farming.MsgUnstake"), OutputType: strp(".syreen.farming.MsgUnstakeResponse")},
					{Name: strp("ClaimReward"), InputType: strp(".syreen.farming.MsgClaimReward"), OutputType: strp(".syreen.farming.MsgClaimRewardResponse")},
					{Name: strp("CreateFarm"), InputType: strp(".syreen.farming.MsgCreateFarm"), OutputType: strp(".syreen.farming.MsgCreateFarmResponse")},
					{Name: strp("UpdateFarm"), InputType: strp(".syreen.farming.MsgUpdateFarm"), OutputType: strp(".syreen.farming.MsgUpdateFarmResponse")},
				},
			},
		},
	}

	rawBz, err := proto.Marshal(fd)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	w, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
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

func strp(s string) *string  { return &s }
func int32p(i int32) *int32 { return &i }
