package types_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/compute/types"
)

func init() {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount("syreen", "syreenpub")
	cfg.Seal()
}

func validAddr() string {
	pk := secp256k1.GenPrivKey().PubKey()
	return sdk.AccAddress(pk.Address()).String()
}

// ---------------------------------------------------------------------------
// MsgStoreCode
// ---------------------------------------------------------------------------

func TestMsgStoreCode_ValidateBasic_Valid(t *testing.T) {
	msg := &types.MsgStoreCode{
		Sender:       validAddr(),
		WASMByteCode: []byte{0x00, 0x61, 0x73, 0x6d},
	}
	require.NoError(t, msg.ValidateBasic())
}

func TestMsgStoreCode_InvalidSender(t *testing.T) {
	msg := &types.MsgStoreCode{
		Sender:       "bad_address",
		WASMByteCode: []byte{0x00, 0x61, 0x73, 0x6d},
	}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrInvalidCreator)
}

func TestMsgStoreCode_EmptyCode(t *testing.T) {
	msg := &types.MsgStoreCode{
		Sender:       validAddr(),
		WASMByteCode: nil,
	}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrEmpty)
}

func TestMsgStoreCode_CodeTooLarge(t *testing.T) {
	msg := &types.MsgStoreCode{
		Sender:       validAddr(),
		WASMByteCode: make([]byte, types.DefaultMaxWasmCodeSize+1),
	}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrCodeTooLarge)
}

func TestMsgStoreCode_ExactMaxSize(t *testing.T) {
	msg := &types.MsgStoreCode{
		Sender:       validAddr(),
		WASMByteCode: make([]byte, types.DefaultMaxWasmCodeSize),
	}
	require.NoError(t, msg.ValidateBasic())
}

// ---------------------------------------------------------------------------
// MsgInstantiateContract
// ---------------------------------------------------------------------------

func TestMsgInstantiateContract_ValidateBasic_Valid(t *testing.T) {
	msg := &types.MsgInstantiateContract{
		Sender: validAddr(),
		Admin:  validAddr(),
		CodeID: 1,
		Label:  "test-contract",
		Msg:    json.RawMessage(`{"init":{}}`),
		Funds:  sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1000)),
	}
	require.NoError(t, msg.ValidateBasic())
}

func TestMsgInstantiateContract_NoAdmin(t *testing.T) {
	msg := &types.MsgInstantiateContract{
		Sender: validAddr(),
		Admin:  "",
		CodeID: 1,
		Label:  "test",
		Msg:    json.RawMessage(`{}`),
	}
	require.NoError(t, msg.ValidateBasic())
}

func TestMsgInstantiateContract_InvalidSender(t *testing.T) {
	msg := &types.MsgInstantiateContract{
		Sender: "bad",
		CodeID: 1,
		Label:  "test",
		Msg:    json.RawMessage(`{}`),
	}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrInvalidCreator)
}

func TestMsgInstantiateContract_InvalidAdmin(t *testing.T) {
	msg := &types.MsgInstantiateContract{
		Sender: validAddr(),
		Admin:  "bad",
		CodeID: 1,
		Label:  "test",
		Msg:    json.RawMessage(`{}`),
	}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrInvalidAdmin)
}

func TestMsgInstantiateContract_ZeroCodeID(t *testing.T) {
	msg := &types.MsgInstantiateContract{
		Sender: validAddr(),
		CodeID: 0,
		Label:  "test",
		Msg:    json.RawMessage(`{}`),
	}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrNoSuchCode)
}

func TestMsgInstantiateContract_EmptyLabel(t *testing.T) {
	msg := &types.MsgInstantiateContract{
		Sender: validAddr(),
		CodeID: 1,
		Label:  "",
		Msg:    json.RawMessage(`{}`),
	}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrInvalidLabel)
}

func TestMsgInstantiateContract_InvalidJSON(t *testing.T) {
	msg := &types.MsgInstantiateContract{
		Sender: validAddr(),
		CodeID: 1,
		Label:  "test",
		Msg:    json.RawMessage(`not json`),
	}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrInvalidMsg)
}

func TestMsgInstantiateContract_ValidFunds(t *testing.T) {
	msg := &types.MsgInstantiateContract{
		Sender: validAddr(),
		CodeID: 1,
		Label:  "test",
		Msg:    json.RawMessage(`{}`),
		Funds:  sdk.NewCoins(sdk.NewInt64Coin("usyreen", 500)),
	}
	require.NoError(t, msg.ValidateBasic())
}

// ---------------------------------------------------------------------------
// MsgExecuteContract
// ---------------------------------------------------------------------------

func TestMsgExecuteContract_ValidateBasic_Valid(t *testing.T) {
	msg := &types.MsgExecuteContract{
		Sender:   validAddr(),
		Contract: validAddr(),
		Msg:      json.RawMessage(`{"do_something":{}}`),
	}
	require.NoError(t, msg.ValidateBasic())
}

func TestMsgExecuteContract_InvalidSender(t *testing.T) {
	msg := &types.MsgExecuteContract{
		Sender:   "bad",
		Contract: validAddr(),
		Msg:      json.RawMessage(`{}`),
	}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrInvalidCreator)
}

func TestMsgExecuteContract_InvalidContract(t *testing.T) {
	msg := &types.MsgExecuteContract{
		Sender:   validAddr(),
		Contract: "bad",
		Msg:      json.RawMessage(`{}`),
	}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrContractNotFound)
}

func TestMsgExecuteContract_InvalidJSON(t *testing.T) {
	msg := &types.MsgExecuteContract{
		Sender:   validAddr(),
		Contract: validAddr(),
		Msg:      json.RawMessage(`{broken`),
	}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrInvalidMsg)
}

// ---------------------------------------------------------------------------
// MsgMigrateContract
// ---------------------------------------------------------------------------

func TestMsgMigrateContract_ValidateBasic_Valid(t *testing.T) {
	msg := &types.MsgMigrateContract{
		Sender:   validAddr(),
		Contract: validAddr(),
		CodeID:   2,
		Msg:      json.RawMessage(`{}`),
	}
	require.NoError(t, msg.ValidateBasic())
}

func TestMsgMigrateContract_InvalidSender(t *testing.T) {
	msg := &types.MsgMigrateContract{
		Sender:   "bad",
		Contract: validAddr(),
		CodeID:   2,
		Msg:      json.RawMessage(`{}`),
	}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrInvalidCreator)
}

func TestMsgMigrateContract_InvalidContract(t *testing.T) {
	msg := &types.MsgMigrateContract{
		Sender:   validAddr(),
		Contract: "bad",
		CodeID:   2,
		Msg:      json.RawMessage(`{}`),
	}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrContractNotFound)
}

func TestMsgMigrateContract_ZeroCodeID(t *testing.T) {
	msg := &types.MsgMigrateContract{
		Sender:   validAddr(),
		Contract: validAddr(),
		CodeID:   0,
		Msg:      json.RawMessage(`{}`),
	}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrNoSuchCode)
}

func TestMsgMigrateContract_InvalidJSON(t *testing.T) {
	msg := &types.MsgMigrateContract{
		Sender:   validAddr(),
		Contract: validAddr(),
		CodeID:   2,
		Msg:      json.RawMessage(`{bad`),
	}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrInvalidMsg)
}

// ---------------------------------------------------------------------------
// MsgUpdateAdmin
// ---------------------------------------------------------------------------

func TestMsgUpdateAdmin_ValidateBasic_Valid(t *testing.T) {
	msg := &types.MsgUpdateAdmin{
		Sender:   validAddr(),
		NewAdmin: validAddr(),
		Contract: validAddr(),
	}
	require.NoError(t, msg.ValidateBasic())
}

func TestMsgUpdateAdmin_InvalidSender(t *testing.T) {
	msg := &types.MsgUpdateAdmin{
		Sender:   "bad",
		NewAdmin: validAddr(),
		Contract: validAddr(),
	}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrInvalidCreator)
}

func TestMsgUpdateAdmin_InvalidNewAdmin(t *testing.T) {
	msg := &types.MsgUpdateAdmin{
		Sender:   validAddr(),
		NewAdmin: "bad",
		Contract: validAddr(),
	}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrInvalidAdmin)
}

func TestMsgUpdateAdmin_InvalidContract(t *testing.T) {
	msg := &types.MsgUpdateAdmin{
		Sender:   validAddr(),
		NewAdmin: validAddr(),
		Contract: "bad",
	}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrContractNotFound)
}

// ---------------------------------------------------------------------------
// AccessConfig
// ---------------------------------------------------------------------------

func TestAccessConfig_Everybody(t *testing.T) {
	ac := types.AccessConfig{Permission: types.AccessTypeEverybody}
	require.Equal(t, types.AccessTypeEverybody, ac.Permission)
	require.Empty(t, ac.Address)
}

func TestAccessConfig_OnlyAddress(t *testing.T) {
	addr := validAddr()
	ac := types.AccessConfig{Permission: types.AccessTypeOnlyAddress, Address: addr}
	require.Equal(t, types.AccessTypeOnlyAddress, ac.Permission)
	require.Equal(t, addr, ac.Address)
}

func TestAccessConfig_Nobody(t *testing.T) {
	ac := types.AccessConfig{Permission: types.AccessTypeNobody}
	require.Equal(t, types.AccessTypeNobody, ac.Permission)
}

// ---------------------------------------------------------------------------
// Params
// ---------------------------------------------------------------------------

func TestDefaultParams(t *testing.T) {
	p := types.DefaultParams()
	require.Equal(t, types.AccessTypeEverybody, p.CodeUploadAccess.Permission)
	require.Equal(t, types.AccessTypeEverybody, p.InstantiateDefaultPermission)
	require.Equal(t, types.DefaultMaxWasmCodeSize, p.MaxWasmCodeSize)
	require.Equal(t, types.DefaultMaxContractGas, p.MaxContractGas)
	require.NoError(t, p.Validate())
}

func TestParams_Validate_ZeroCodeSize(t *testing.T) {
	p := types.DefaultParams()
	p.MaxWasmCodeSize = 0
	require.Error(t, p.Validate())
}

func TestParams_Validate_ZeroGas(t *testing.T) {
	p := types.DefaultParams()
	p.MaxContractGas = 0
	require.Error(t, p.Validate())
}

func TestParams_Validate_InvalidInstantiatePermission(t *testing.T) {
	p := types.DefaultParams()
	p.InstantiateDefaultPermission = "InvalidPerm"
	require.Error(t, p.Validate())
}

func TestParams_Validate_InvalidUploadPermission(t *testing.T) {
	p := types.DefaultParams()
	p.CodeUploadAccess.Permission = "InvalidPerm"
	require.Error(t, p.Validate())
}

func TestParams_Validate_OnlyAddress(t *testing.T) {
	p := types.DefaultParams()
	p.CodeUploadAccess = types.AccessConfig{
		Permission: types.AccessTypeOnlyAddress,
		Address:    validAddr(),
	}
	p.InstantiateDefaultPermission = types.AccessTypeNobody
	require.NoError(t, p.Validate())
}

// ---------------------------------------------------------------------------
// GenesisState
// ---------------------------------------------------------------------------

func TestDefaultGenesis(t *testing.T) {
	gs := types.DefaultGenesis()
	require.NotNil(t, gs)
	require.NoError(t, gs.Validate())
	require.Empty(t, gs.Codes)
	require.Empty(t, gs.Contracts)
}

func TestGenesisState_Validate_Valid(t *testing.T) {
	gs := &types.GenesisState{
		Params: types.DefaultParams(),
		Codes: []types.GenesisCode{
			{
				CodeID:    1,
				CodeInfo:  types.CodeInfo{CodeID: 1, Creator: validAddr()},
				CodeBytes: []byte{0x00},
			},
		},
		Contracts: []types.GenesisContract{
			{
				ContractAddress: validAddr(),
				ContractInfo:    types.ContractInfo{CodeID: 1, Creator: validAddr()},
				ContractState:   []types.Model{{Key: []byte("k"), Value: []byte("v")}},
			},
		},
	}
	require.NoError(t, gs.Validate())
}

func TestGenesisState_Validate_DuplicateCodeIDs(t *testing.T) {
	gs := &types.GenesisState{
		Params: types.DefaultParams(),
		Codes: []types.GenesisCode{
			{CodeID: 1, CodeInfo: types.CodeInfo{CodeID: 1}},
			{CodeID: 1, CodeInfo: types.CodeInfo{CodeID: 1}},
		},
	}
	err := gs.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "duplicate code ID")
}

func TestGenesisState_Validate_ZeroCodeID(t *testing.T) {
	gs := &types.GenesisState{
		Params: types.DefaultParams(),
		Codes: []types.GenesisCode{
			{CodeID: 0, CodeInfo: types.CodeInfo{CodeID: 0}},
		},
	}
	err := gs.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "code ID cannot be zero")
}

func TestGenesisState_Validate_DuplicateContractAddresses(t *testing.T) {
	addr := validAddr()
	gs := &types.GenesisState{
		Params: types.DefaultParams(),
		Codes: []types.GenesisCode{
			{CodeID: 1, CodeInfo: types.CodeInfo{CodeID: 1}},
		},
		Contracts: []types.GenesisContract{
			{ContractAddress: addr, ContractInfo: types.ContractInfo{CodeID: 1}},
			{ContractAddress: addr, ContractInfo: types.ContractInfo{CodeID: 1}},
		},
	}
	err := gs.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "duplicate contract address")
}

func TestGenesisState_Validate_UnknownCodeRef(t *testing.T) {
	gs := &types.GenesisState{
		Params: types.DefaultParams(),
		Codes: []types.GenesisCode{
			{CodeID: 1, CodeInfo: types.CodeInfo{CodeID: 1}},
		},
		Contracts: []types.GenesisContract{
			{ContractAddress: validAddr(), ContractInfo: types.ContractInfo{CodeID: 99}},
		},
	}
	err := gs.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown code ID")
}

func TestGenesisState_Validate_InvalidParams(t *testing.T) {
	gs := &types.GenesisState{
		Params: types.Params{MaxWasmCodeSize: 0, MaxContractGas: 1, InstantiateDefaultPermission: types.AccessTypeEverybody, CodeUploadAccess: types.AccessConfig{Permission: types.AccessTypeEverybody}},
	}
	require.Error(t, gs.Validate())
}

// ---------------------------------------------------------------------------
// ProtoMessage and String helpers
// ---------------------------------------------------------------------------

func TestCodeInfo_StringAndProto(t *testing.T) {
	ci := &types.CodeInfo{CodeID: 5, Creator: "someone"}
	ci.ProtoMessage()
	require.Contains(t, ci.String(), "5")
	require.Equal(t, "syreen.compute.CodeInfo", ci.XXX_MessageName())
	ci.Reset()
	require.Equal(t, uint64(0), ci.CodeID)
}

func TestContractInfo_StringAndProto(t *testing.T) {
	ci := &types.ContractInfo{Address: "addr", CodeID: 3}
	ci.ProtoMessage()
	require.Contains(t, ci.String(), "addr")
	require.Equal(t, "syreen.compute.ContractInfo", ci.XXX_MessageName())
	ci.Reset()
	require.Equal(t, uint64(0), ci.CodeID)
}

func TestContractState_StringAndProto(t *testing.T) {
	cs := &types.ContractState{Key: []byte("k"), Value: []byte("v")}
	cs.ProtoMessage()
	require.Contains(t, cs.String(), "state key=")
	require.Equal(t, "syreen.compute.ContractState", cs.XXX_MessageName())
	cs.Reset()
	require.Nil(t, cs.Key)
}

func TestGenesisState_StringAndProto(t *testing.T) {
	gs := types.DefaultGenesis()
	gs.ProtoMessage()
	require.Contains(t, gs.String(), "compute genesis")
	require.Equal(t, "syreen.compute.GenesisState", gs.XXX_MessageName())
}

func TestMsgStoreCode_StringAndProto(t *testing.T) {
	msg := &types.MsgStoreCode{Sender: "alice", WASMByteCode: []byte{1, 2}}
	msg.ProtoMessage()
	require.Contains(t, msg.String(), "store_code")
	require.Equal(t, "syreen.compute.MsgStoreCode", msg.XXX_MessageName())
	msg.Reset()
	require.Equal(t, "", msg.Sender)
}

func TestMsgInstantiateContract_String(t *testing.T) {
	msg := &types.MsgInstantiateContract{Sender: "alice", CodeID: 7}
	require.Contains(t, msg.String(), "7")
	require.Equal(t, "syreen.compute.MsgInstantiate", msg.XXX_MessageName())
}

func TestMsgExecuteContract_String(t *testing.T) {
	msg := &types.MsgExecuteContract{Sender: "alice", Contract: "contract1"}
	require.Contains(t, msg.String(), "contract1")
	require.Equal(t, "syreen.compute.MsgExecute", msg.XXX_MessageName())
}

func TestMsgMigrateContract_String(t *testing.T) {
	msg := &types.MsgMigrateContract{Contract: "c1", CodeID: 3}
	require.Contains(t, msg.String(), "3")
	require.Equal(t, "syreen.compute.MsgMigrate", msg.XXX_MessageName())
}

func TestMsgUpdateAdmin_String(t *testing.T) {
	msg := &types.MsgUpdateAdmin{Contract: "c1", NewAdmin: "new1"}
	require.Contains(t, msg.String(), "new1")
	require.Equal(t, "syreen.compute.MsgUpdateAdmin", msg.XXX_MessageName())
}

// ---------------------------------------------------------------------------
// GetSigners
// ---------------------------------------------------------------------------

func TestMsgStoreCode_GetSigners(t *testing.T) {
	addr := validAddr()
	msg := &types.MsgStoreCode{Sender: addr, WASMByteCode: []byte{1}}
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, addr, signers[0].String())
}

func TestMsgInstantiateContract_GetSigners(t *testing.T) {
	addr := validAddr()
	msg := &types.MsgInstantiateContract{Sender: addr, CodeID: 1, Label: "l", Msg: json.RawMessage(`{}`)}
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, addr, signers[0].String())
}

func TestMsgExecuteContract_GetSigners(t *testing.T) {
	addr := validAddr()
	msg := &types.MsgExecuteContract{Sender: addr, Contract: validAddr(), Msg: json.RawMessage(`{}`)}
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, addr, signers[0].String())
}

func TestMsgMigrateContract_GetSigners(t *testing.T) {
	addr := validAddr()
	msg := &types.MsgMigrateContract{Sender: addr, Contract: validAddr(), CodeID: 1, Msg: json.RawMessage(`{}`)}
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, addr, signers[0].String())
}

func TestMsgUpdateAdmin_GetSigners(t *testing.T) {
	addr := validAddr()
	msg := &types.MsgUpdateAdmin{Sender: addr, NewAdmin: validAddr(), Contract: validAddr()}
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, addr, signers[0].String())
}

// ---------------------------------------------------------------------------
// Keys helpers
// ---------------------------------------------------------------------------

func TestUint64BytesRoundtrip(t *testing.T) {
	for _, val := range []uint64{0, 1, 42, 1<<32 - 1, 1<<64 - 1} {
		b := types.Uint64ToBytes(val)
		require.Len(t, b, 8)
		require.Equal(t, val, types.BytesToUint64(b))
	}
}

func TestBytesToUint64_Short(t *testing.T) {
	require.Equal(t, uint64(0), types.BytesToUint64([]byte{1, 2, 3}))
}

func TestCodeInfoKey(t *testing.T) {
	key := types.CodeInfoKey(1)
	require.True(t, strings.HasPrefix(string(key), "code/info/"))
}

func TestCodeBytecodeKey(t *testing.T) {
	key := types.CodeBytecodeKey(1)
	require.True(t, strings.HasPrefix(string(key), "code/bytecode/"))
}

func TestContractInfoKey(t *testing.T) {
	key := types.ContractInfoKey("syreen1abc")
	require.Equal(t, "contract/syreen1abc", string(key))
}

func TestContractStateKey(t *testing.T) {
	key := types.ContractStateKey("syreen1abc", []byte("mykey"))
	require.Equal(t, "state/syreen1abc/mykey", string(key))
}

func TestContractStateIteratorPrefix(t *testing.T) {
	prefix := types.ContractStateIteratorPrefix("syreen1abc")
	require.Equal(t, "state/syreen1abc/", string(prefix))
}
