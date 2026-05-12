package types_test

import (
	"encoding/json"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"syreen/x/abstractaccount/types"
)

func init() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("syreen", "syreenpub")
	config.SetBech32PrefixForValidator("syreenvaloper", "syreenvaloperpub")
	config.SetBech32PrefixForConsensusNode("syreenvalcons", "syreenvalconspub")
}

func validAddr() string {
	return sdk.AccAddress([]byte("test_address_padded_")).String() // 20 bytes
}

func validAddr2() string {
	return sdk.AccAddress([]byte("second_addr_padding_")).String() // 20 bytes
}

func validAddr3() string {
	return sdk.AccAddress([]byte("third__addr_padding_")).String() // 20 bytes
}

// ---------------------------------------------------------------------------
// MsgCreateSmartAccount
// ---------------------------------------------------------------------------

func TestMsgCreateSmartAccount_ValidateBasic(t *testing.T) {
	addr1 := validAddr()
	addr2 := validAddr2()

	tests := []struct {
		name    string
		msg     types.MsgCreateSmartAccount
		wantErr bool
	}{
		{
			name: "valid multisig",
			msg: types.MsgCreateSmartAccount{
				Sender:      addr1,
				AccountType: types.AccountTypeMultiSig,
				Owners:      []string{addr1, addr2},
				Threshold:   2,
			},
			wantErr: false,
		},
		{
			name: "valid social single owner",
			msg: types.MsgCreateSmartAccount{
				Sender:      addr1,
				AccountType: types.AccountTypeSocial,
				Owners:      []string{addr1},
				Threshold:   1,
			},
			wantErr: false,
		},
		{
			name: "valid session type",
			msg: types.MsgCreateSmartAccount{
				Sender:      addr1,
				AccountType: types.AccountTypeSession,
				Owners:      []string{addr1},
				Threshold:   1,
			},
			wantErr: false,
		},
		{
			name: "invalid sender",
			msg: types.MsgCreateSmartAccount{
				Sender:      "bad",
				AccountType: types.AccountTypeMultiSig,
				Owners:      []string{addr1},
				Threshold:   1,
			},
			wantErr: true,
		},
		{
			name: "invalid owner address",
			msg: types.MsgCreateSmartAccount{
				Sender:      addr1,
				AccountType: types.AccountTypeMultiSig,
				Owners:      []string{"invalid_owner"},
				Threshold:   1,
			},
			wantErr: true,
		},
		{
			name: "invalid account type",
			msg: types.MsgCreateSmartAccount{
				Sender:      addr1,
				AccountType: "banana",
				Owners:      []string{addr1},
				Threshold:   1,
			},
			wantErr: true,
		},
		{
			name: "threshold 0",
			msg: types.MsgCreateSmartAccount{
				Sender:      addr1,
				AccountType: types.AccountTypeMultiSig,
				Owners:      []string{addr1},
				Threshold:   0,
			},
			wantErr: true,
		},
		{
			name: "threshold > owners",
			msg: types.MsgCreateSmartAccount{
				Sender:      addr1,
				AccountType: types.AccountTypeMultiSig,
				Owners:      []string{addr1},
				Threshold:   5,
			},
			wantErr: true,
		},
		{
			name: "no owners",
			msg: types.MsgCreateSmartAccount{
				Sender:      addr1,
				AccountType: types.AccountTypeMultiSig,
				Owners:      []string{},
				Threshold:   1,
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// MsgCreateSessionKey
// ---------------------------------------------------------------------------

func TestMsgCreateSessionKey_ValidateBasic(t *testing.T) {
	addr1 := validAddr()
	addr2 := validAddr2()

	tests := []struct {
		name    string
		msg     types.MsgCreateSessionKey
		wantErr bool
	}{
		{
			name: "valid",
			msg: types.MsgCreateSessionKey{
				Granter:     addr1,
				Grantee:     addr2,
				Permissions: []types.Permission{{MsgType: "/cosmos.bank.v1beta1.MsgSend"}},
				Duration:    time.Hour,
			},
			wantErr: false,
		},
		{
			name: "invalid granter",
			msg: types.MsgCreateSessionKey{
				Granter:     "bad",
				Grantee:     addr2,
				Permissions: []types.Permission{{MsgType: "/cosmos.bank.v1beta1.MsgSend"}},
				Duration:    time.Hour,
			},
			wantErr: true,
		},
		{
			name: "invalid grantee",
			msg: types.MsgCreateSessionKey{
				Granter:     addr1,
				Grantee:     "bad",
				Permissions: []types.Permission{{MsgType: "/cosmos.bank.v1beta1.MsgSend"}},
				Duration:    time.Hour,
			},
			wantErr: true,
		},
		{
			name: "zero duration",
			msg: types.MsgCreateSessionKey{
				Granter:     addr1,
				Grantee:     addr2,
				Permissions: []types.Permission{{MsgType: "/cosmos.bank.v1beta1.MsgSend"}},
				Duration:    0,
			},
			wantErr: true,
		},
		{
			name: "negative duration",
			msg: types.MsgCreateSessionKey{
				Granter:     addr1,
				Grantee:     addr2,
				Permissions: []types.Permission{{MsgType: "/cosmos.bank.v1beta1.MsgSend"}},
				Duration:    -time.Hour,
			},
			wantErr: true,
		},
		{
			name: "empty permissions",
			msg: types.MsgCreateSessionKey{
				Granter:     addr1,
				Grantee:     addr2,
				Permissions: []types.Permission{},
				Duration:    time.Hour,
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// MsgRevokeSessionKey
// ---------------------------------------------------------------------------

func TestMsgRevokeSessionKey_ValidateBasic(t *testing.T) {
	addr1 := validAddr()
	addr2 := validAddr2()

	tests := []struct {
		name    string
		msg     types.MsgRevokeSessionKey
		wantErr bool
	}{
		{
			name:    "valid",
			msg:     types.MsgRevokeSessionKey{Granter: addr1, SessionKeyAddr: addr2},
			wantErr: false,
		},
		{
			name:    "invalid granter",
			msg:     types.MsgRevokeSessionKey{Granter: "bad", SessionKeyAddr: addr2},
			wantErr: true,
		},
		{
			name:    "invalid session key addr",
			msg:     types.MsgRevokeSessionKey{Granter: addr1, SessionKeyAddr: "bad"},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// MsgInitiateRecovery
// ---------------------------------------------------------------------------

func TestMsgInitiateRecovery_ValidateBasic(t *testing.T) {
	addr1 := validAddr()
	addr2 := validAddr2()
	addr3 := validAddr3()

	tests := []struct {
		name    string
		msg     types.MsgInitiateRecovery
		wantErr bool
	}{
		{
			name:    "valid",
			msg:     types.MsgInitiateRecovery{Guardian: addr1, Account: addr2, NewOwners: []string{addr3}},
			wantErr: false,
		},
		{
			name:    "invalid guardian",
			msg:     types.MsgInitiateRecovery{Guardian: "bad", Account: addr2, NewOwners: []string{addr3}},
			wantErr: true,
		},
		{
			name:    "invalid account",
			msg:     types.MsgInitiateRecovery{Guardian: addr1, Account: "bad", NewOwners: []string{addr3}},
			wantErr: true,
		},
		{
			name:    "empty new owners",
			msg:     types.MsgInitiateRecovery{Guardian: addr1, Account: addr2, NewOwners: []string{}},
			wantErr: true,
		},
		{
			name:    "invalid new owner address",
			msg:     types.MsgInitiateRecovery{Guardian: addr1, Account: addr2, NewOwners: []string{"invalid"}},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// MsgApproveRecovery
// ---------------------------------------------------------------------------

func TestMsgApproveRecovery_ValidateBasic(t *testing.T) {
	addr1 := validAddr()
	addr2 := validAddr2()

	tests := []struct {
		name    string
		msg     types.MsgApproveRecovery
		wantErr bool
	}{
		{
			name:    "valid",
			msg:     types.MsgApproveRecovery{Guardian: addr1, Account: addr2},
			wantErr: false,
		},
		{
			name:    "invalid guardian",
			msg:     types.MsgApproveRecovery{Guardian: "bad", Account: addr2},
			wantErr: true,
		},
		{
			name:    "invalid account",
			msg:     types.MsgApproveRecovery{Guardian: addr1, Account: "bad"},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// MsgExecuteRecovery
// ---------------------------------------------------------------------------

func TestMsgExecuteRecovery_ValidateBasic(t *testing.T) {
	addr1 := validAddr()
	addr2 := validAddr2()

	tests := []struct {
		name    string
		msg     types.MsgExecuteRecovery
		wantErr bool
	}{
		{
			name:    "valid",
			msg:     types.MsgExecuteRecovery{Sender: addr1, Account: addr2},
			wantErr: false,
		},
		{
			name:    "invalid sender",
			msg:     types.MsgExecuteRecovery{Sender: "bad", Account: addr2},
			wantErr: true,
		},
		{
			name:    "invalid account",
			msg:     types.MsgExecuteRecovery{Sender: addr1, Account: "bad"},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// MsgSponsorGas
// ---------------------------------------------------------------------------

func TestMsgSponsorGas_ValidateBasic(t *testing.T) {
	addr1 := validAddr()
	addr2 := validAddr2()

	tests := []struct {
		name    string
		msg     types.MsgSponsorGas
		wantErr bool
	}{
		{
			name:    "valid",
			msg:     types.MsgSponsorGas{Sponsor: addr1, Sponsored: addr2, GasLimit: 1000000, Duration: time.Hour},
			wantErr: false,
		},
		{
			name:    "invalid sponsor",
			msg:     types.MsgSponsorGas{Sponsor: "bad", Sponsored: addr2, GasLimit: 1000000, Duration: time.Hour},
			wantErr: true,
		},
		{
			name:    "invalid sponsored",
			msg:     types.MsgSponsorGas{Sponsor: addr1, Sponsored: "bad", GasLimit: 1000000, Duration: time.Hour},
			wantErr: true,
		},
		{
			name:    "zero gas limit",
			msg:     types.MsgSponsorGas{Sponsor: addr1, Sponsored: addr2, GasLimit: 0, Duration: time.Hour},
			wantErr: true,
		},
		{
			name:    "zero duration",
			msg:     types.MsgSponsorGas{Sponsor: addr1, Sponsored: addr2, GasLimit: 1000000, Duration: 0},
			wantErr: true,
		},
		{
			name:    "negative duration",
			msg:     types.MsgSponsorGas{Sponsor: addr1, Sponsored: addr2, GasLimit: 1000000, Duration: -time.Hour},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// MsgBatchExecute
// ---------------------------------------------------------------------------

func TestMsgBatchExecute_ValidateBasic(t *testing.T) {
	addr1 := validAddr()

	tests := []struct {
		name    string
		msg     types.MsgBatchExecute
		wantErr bool
	}{
		{
			name: "valid",
			msg: types.MsgBatchExecute{
				Sender:   addr1,
				Messages: []json.RawMessage{json.RawMessage(`{"type":"send"}`)},
			},
			wantErr: false,
		},
		{
			name: "invalid sender",
			msg: types.MsgBatchExecute{
				Sender:   "bad",
				Messages: []json.RawMessage{json.RawMessage(`{"type":"send"}`)},
			},
			wantErr: true,
		},
		{
			name: "empty messages",
			msg: types.MsgBatchExecute{
				Sender:   addr1,
				Messages: []json.RawMessage{},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Params
// ---------------------------------------------------------------------------

func TestDefaultParams(t *testing.T) {
	p := types.DefaultParams()
	require.Equal(t, 24*time.Hour, p.MaxSessionKeyDuration)
	require.Equal(t, uint32(7), p.MaxGuardians)
	require.Equal(t, 48*time.Hour, p.RecoveryDelayPeriod)
	require.Equal(t, uint32(10), p.MaxBatchSize)
	require.True(t, p.EnableGasSponsorship)
	require.NoError(t, p.Validate())
}

func TestParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  types.Params
		wantErr bool
	}{
		{
			name:    "default params valid",
			params:  types.DefaultParams(),
			wantErr: false,
		},
		{
			name: "zero session key duration",
			params: types.Params{
				MaxSessionKeyDuration: 0,
				MaxGuardians:          7,
				RecoveryDelayPeriod:   48 * time.Hour,
				MaxBatchSize:          10,
				EnableGasSponsorship:  true,
			},
			wantErr: true,
		},
		{
			name: "zero max guardians",
			params: types.Params{
				MaxSessionKeyDuration: 24 * time.Hour,
				MaxGuardians:          0,
				RecoveryDelayPeriod:   48 * time.Hour,
				MaxBatchSize:          10,
				EnableGasSponsorship:  true,
			},
			wantErr: true,
		},
		{
			name: "zero recovery delay",
			params: types.Params{
				MaxSessionKeyDuration: 24 * time.Hour,
				MaxGuardians:          7,
				RecoveryDelayPeriod:   0,
				MaxBatchSize:          10,
				EnableGasSponsorship:  true,
			},
			wantErr: true,
		},
		{
			name: "zero max batch size",
			params: types.Params{
				MaxSessionKeyDuration: 24 * time.Hour,
				MaxGuardians:          7,
				RecoveryDelayPeriod:   48 * time.Hour,
				MaxBatchSize:          0,
				EnableGasSponsorship:  true,
			},
			wantErr: true,
		},
		{
			name: "negative session key duration",
			params: types.Params{
				MaxSessionKeyDuration: -time.Hour,
				MaxGuardians:          7,
				RecoveryDelayPeriod:   48 * time.Hour,
				MaxBatchSize:          10,
				EnableGasSponsorship:  true,
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.params.Validate()
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Genesis
// ---------------------------------------------------------------------------

func TestDefaultGenesis(t *testing.T) {
	gs := types.DefaultGenesis()
	require.NotNil(t, gs)
	require.NoError(t, gs.Validate())
	require.Equal(t, types.DefaultParams(), gs.Params)
}

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		name    string
		gs      types.GenesisState
		wantErr bool
	}{
		{
			name:    "default is valid",
			gs:      *types.DefaultGenesis(),
			wantErr: false,
		},
		{
			name: "invalid params - zero batch size",
			gs: types.GenesisState{
				Params: types.Params{
					MaxSessionKeyDuration: 24 * time.Hour,
					MaxGuardians:          7,
					RecoveryDelayPeriod:   48 * time.Hour,
					MaxBatchSize:          0,
					EnableGasSponsorship:  true,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid params - zero guardians",
			gs: types.GenesisState{
				Params: types.Params{
					MaxSessionKeyDuration: 24 * time.Hour,
					MaxGuardians:          0,
					RecoveryDelayPeriod:   48 * time.Hour,
					MaxBatchSize:          10,
					EnableGasSponsorship:  true,
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.gs.Validate()
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
