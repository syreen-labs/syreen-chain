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
	registerVaultProtoFileDescriptors()
}

func registerVaultProtoFileDescriptors() {
	label := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	typeString := descriptorpb.FieldDescriptorProto_TYPE_STRING
	typeBytes := descriptorpb.FieldDescriptorProto_TYPE_BYTES
	typeUint64 := descriptorpb.FieldDescriptorProto_TYPE_UINT64

	fd := &descriptorpb.FileDescriptorProto{
		Name:    sp("syreen/vault/tx.proto"),
		Syntax:  sp("proto3"),
		Package: sp("syreen.vault"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: sp("MsgCreateVault"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("creator"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("creator")},
				{Name: sp("name"), Number: ip(2), Label: &label, Type: &typeString, JsonName: sp("name")},
				{Name: sp("deposit_denom"), Number: ip(3), Label: &label, Type: &typeString, JsonName: sp("depositDenom")},
				{Name: sp("strategy_type"), Number: ip(4), Label: &label, Type: &typeString, JsonName: sp("strategyType")},
				{Name: sp("target_pool_ids"), Number: ip(5), Label: &label, Type: &typeBytes, JsonName: sp("targetPoolIds")},
				{Name: sp("performance_fee"), Number: ip(6), Label: &label, Type: &typeBytes, JsonName: sp("performanceFee")},
			}},
			{Name: sp("MsgCreateVaultResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("vault_id"), Number: ip(1), Label: &label, Type: &typeUint64, JsonName: sp("vaultId")},
			}},
			{Name: sp("MsgDepositVault"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("vault_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("vaultId")},
				{Name: sp("amount"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("amount")},
			}},
			{Name: sp("MsgDepositVaultResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("shares_minted"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("sharesMinted")},
			}},
			{Name: sp("MsgWithdrawVault"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("vault_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("vaultId")},
				{Name: sp("shares"), Number: ip(3), Label: &label, Type: &typeBytes, JsonName: sp("shares")},
			}},
			{Name: sp("MsgWithdrawVaultResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("amount_returned"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("amountReturned")},
			}},
			{Name: sp("MsgCompoundVault"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("sender"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("sender")},
				{Name: sp("vault_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("vaultId")},
			}},
			{Name: sp("MsgCompoundVaultResponse"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("yield_generated"), Number: ip(1), Label: &label, Type: &typeBytes, JsonName: sp("yieldGenerated")},
			}},
			{Name: sp("MsgUpdateVaultStrategy"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: sp("creator"), Number: ip(1), Label: &label, Type: &typeString, JsonName: sp("creator")},
				{Name: sp("vault_id"), Number: ip(2), Label: &label, Type: &typeUint64, JsonName: sp("vaultId")},
				{Name: sp("strategy_type"), Number: ip(3), Label: &label, Type: &typeString, JsonName: sp("strategyType")},
				{Name: sp("target_pool_ids"), Number: ip(4), Label: &label, Type: &typeBytes, JsonName: sp("targetPoolIds")},
			}},
			{Name: sp("MsgUpdateVaultStrategyResponse")},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: sp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: sp("CreateVault"), InputType: sp(".syreen.vault.MsgCreateVault"), OutputType: sp(".syreen.vault.MsgCreateVaultResponse")},
					{Name: sp("DepositVault"), InputType: sp(".syreen.vault.MsgDepositVault"), OutputType: sp(".syreen.vault.MsgDepositVaultResponse")},
					{Name: sp("WithdrawVault"), InputType: sp(".syreen.vault.MsgWithdrawVault"), OutputType: sp(".syreen.vault.MsgWithdrawVaultResponse")},
					{Name: sp("CompoundVault"), InputType: sp(".syreen.vault.MsgCompoundVault"), OutputType: sp(".syreen.vault.MsgCompoundVaultResponse")},
					{Name: sp("UpdateVaultStrategy"), InputType: sp(".syreen.vault.MsgUpdateVaultStrategy"), OutputType: sp(".syreen.vault.MsgUpdateVaultStrategyResponse")},
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
	fileDescriptorVaultTx = buf.Bytes()

	file, err := protodesc.NewFile(fd, nil)
	if err != nil {
		panic("vault: failed to create proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(file); err != nil {
		// Already registered is fine
	}
}

func sp(s string) *string { return &s }
func ip(n int32) *int32   { return &n }
