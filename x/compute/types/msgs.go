package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgStoreCode          = "store_code"
	TypeMsgInstantiateContract = "instantiate"
	TypeMsgExecuteContract    = "execute"
	TypeMsgMigrateContract    = "migrate"
	TypeMsgUpdateAdmin        = "update_admin"
)

// MsgStoreCode uploads Wasm bytecode to the chain
type MsgStoreCode struct {
	Sender                string       `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	WASMByteCode          []byte       `protobuf:"bytes,2,opt,name=wasm_byte_code,json=wasmByteCode,proto3" json:"wasm_byte_code"`
	InstantiatePermission *AccessConfig `protobuf:"bytes,3,opt,name=instantiate_permission,json=instantiatePermission,proto3" json:"instantiate_permission,omitempty"`
}

func (m *MsgStoreCode) ProtoMessage()           {}
func (m *MsgStoreCode) Reset()                  { *m = MsgStoreCode{} }
func (m *MsgStoreCode) String() string          { return fmt.Sprintf("store_code from %s (%d bytes)", m.Sender, len(m.WASMByteCode)) }
func (m *MsgStoreCode) XXX_MessageName() string { return "syreen.compute.MsgStoreCode" }

func (m *MsgStoreCode) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return ErrInvalidCreator
	}
	if len(m.WASMByteCode) == 0 {
		return ErrEmpty
	}
	if uint64(len(m.WASMByteCode)) > DefaultMaxWasmCodeSize {
		return ErrCodeTooLarge
	}
	return nil
}

// MsgInstantiateContract creates a new smart contract instance
type MsgInstantiateContract struct {
	Sender string          `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	Admin  string          `protobuf:"bytes,2,opt,name=admin,proto3" json:"admin,omitempty"`
	CodeID uint64          `protobuf:"varint,3,opt,name=code_id,json=codeId,proto3" json:"code_id"`
	Label  string          `protobuf:"bytes,4,opt,name=label,proto3" json:"label"`
	Msg    json.RawMessage `protobuf:"bytes,5,opt,name=msg,proto3,casttype=encoding/json.RawMessage" json:"msg"`
	Funds  sdk.Coins       `protobuf:"bytes,6,rep,name=funds,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"funds,omitempty"`
}

func (m *MsgInstantiateContract) ProtoMessage()           {}
func (m *MsgInstantiateContract) Reset()                  { *m = MsgInstantiateContract{} }
func (m *MsgInstantiateContract) String() string          { return fmt.Sprintf("instantiate code %d from %s", m.CodeID, m.Sender) }
func (m *MsgInstantiateContract) XXX_MessageName() string { return "syreen.compute.MsgInstantiateContract" }

func (m *MsgInstantiateContract) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return ErrInvalidCreator
	}
	if m.Admin != "" {
		_, err := sdk.AccAddressFromBech32(m.Admin)
		if err != nil {
			return ErrInvalidAdmin
		}
	}
	if m.CodeID == 0 {
		return ErrNoSuchCode
	}
	if len(m.Label) == 0 {
		return ErrInvalidLabel
	}
	if len(m.Label) > 128 {
		return fmt.Errorf("label exceeds maximum length of 128 bytes")
	}
	if !json.Valid(m.Msg) {
		return ErrInvalidMsg
	}
	if !m.Funds.IsValid() {
		return ErrInvalidMsg
	}
	return nil
}

// MsgExecuteContract calls a smart contract
type MsgExecuteContract struct {
	Sender   string          `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	Contract string          `protobuf:"bytes,2,opt,name=contract,proto3" json:"contract"`
	Msg      json.RawMessage `protobuf:"bytes,3,opt,name=msg,proto3,casttype=encoding/json.RawMessage" json:"msg"`
	Funds    sdk.Coins       `protobuf:"bytes,4,rep,name=funds,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"funds,omitempty"`
}

func (m *MsgExecuteContract) ProtoMessage()           {}
func (m *MsgExecuteContract) Reset()                  { *m = MsgExecuteContract{} }
func (m *MsgExecuteContract) String() string          { return fmt.Sprintf("execute %s from %s", m.Contract, m.Sender) }
func (m *MsgExecuteContract) XXX_MessageName() string { return "syreen.compute.MsgExecuteContract" }

func (m *MsgExecuteContract) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return ErrInvalidCreator
	}
	_, err = sdk.AccAddressFromBech32(m.Contract)
	if err != nil {
		return ErrContractNotFound
	}
	if !json.Valid(m.Msg) {
		return ErrInvalidMsg
	}
	if !m.Funds.IsValid() {
		return ErrInvalidMsg
	}
	return nil
}

// MsgMigrateContract runs a code upgrade for a smart contract
type MsgMigrateContract struct {
	Sender   string          `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	Contract string          `protobuf:"bytes,2,opt,name=contract,proto3" json:"contract"`
	CodeID   uint64          `protobuf:"varint,3,opt,name=code_id,json=codeId,proto3" json:"code_id"`
	Msg      json.RawMessage `protobuf:"bytes,4,opt,name=msg,proto3,casttype=encoding/json.RawMessage" json:"msg"`
}

func (m *MsgMigrateContract) ProtoMessage()           {}
func (m *MsgMigrateContract) Reset()                  { *m = MsgMigrateContract{} }
func (m *MsgMigrateContract) String() string          { return fmt.Sprintf("migrate %s to code %d", m.Contract, m.CodeID) }
func (m *MsgMigrateContract) XXX_MessageName() string { return "syreen.compute.MsgMigrateContract" }

func (m *MsgMigrateContract) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return ErrInvalidCreator
	}
	_, err = sdk.AccAddressFromBech32(m.Contract)
	if err != nil {
		return ErrContractNotFound
	}
	if m.CodeID == 0 {
		return ErrNoSuchCode
	}
	if !json.Valid(m.Msg) {
		return ErrInvalidMsg
	}
	return nil
}

// MsgUpdateAdmin sets a new admin for a contract
type MsgUpdateAdmin struct {
	Sender   string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender"`
	NewAdmin string `protobuf:"bytes,2,opt,name=new_admin,json=newAdmin,proto3" json:"new_admin"`
	Contract string `protobuf:"bytes,3,opt,name=contract,proto3" json:"contract"`
}

func (m *MsgUpdateAdmin) ProtoMessage()           {}
func (m *MsgUpdateAdmin) Reset()                  { *m = MsgUpdateAdmin{} }
func (m *MsgUpdateAdmin) String() string          { return fmt.Sprintf("update_admin %s to %s", m.Contract, m.NewAdmin) }
func (m *MsgUpdateAdmin) XXX_MessageName() string { return "syreen.compute.MsgUpdateAdmin" }

func (m *MsgUpdateAdmin) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return ErrInvalidCreator
	}
	_, err = sdk.AccAddressFromBech32(m.NewAdmin)
	if err != nil {
		return ErrInvalidAdmin
	}
	_, err = sdk.AccAddressFromBech32(m.Contract)
	if err != nil {
		return ErrContractNotFound
	}
	return nil
}
