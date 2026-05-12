package types_test

import (
	"bytes"
	"crypto/sha256"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"syreen/x/mevprotection/types"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func validAddr() string {
	addr := sdk.AccAddress([]byte("test_address_padded_"))
	return addr.String()
}

func validHash() []byte {
	h := sha256.Sum256([]byte("test-body-nonce"))
	return h[:]
}

// ---------------------------------------------------------------------------
// MsgCommitTx
// ---------------------------------------------------------------------------

func TestMsgCommitTx_Valid(t *testing.T) {
	msg := &types.MsgCommitTx{
		Sender:      validAddr(),
		TxHash:      validHash(),
		EncryptedTx: []byte("encrypted-payload"),
	}
	require.NoError(t, msg.ValidateBasic())
}

func TestMsgCommitTx_EmptySender(t *testing.T) {
	msg := &types.MsgCommitTx{
		Sender:      "",
		TxHash:      validHash(),
		EncryptedTx: []byte("encrypted-payload"),
	}
	require.Error(t, msg.ValidateBasic())
	require.Contains(t, msg.ValidateBasic().Error(), "invalid sender")
}

func TestMsgCommitTx_InvalidSender(t *testing.T) {
	msg := &types.MsgCommitTx{
		Sender:      "not-a-bech32-address",
		TxHash:      validHash(),
		EncryptedTx: []byte("encrypted-payload"),
	}
	require.Error(t, msg.ValidateBasic())
	require.Contains(t, msg.ValidateBasic().Error(), "invalid sender")
}

func TestMsgCommitTx_WrongHashLength(t *testing.T) {
	msg := &types.MsgCommitTx{
		Sender:      validAddr(),
		TxHash:      []byte("short"),
		EncryptedTx: []byte("encrypted-payload"),
	}
	err := msg.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "32 bytes")
}

func TestMsgCommitTx_EmptyEncryptedTx(t *testing.T) {
	msg := &types.MsgCommitTx{
		Sender:      validAddr(),
		TxHash:      validHash(),
		EncryptedTx: nil,
	}
	err := msg.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty")
}

func TestMsgCommitTx_TooLargeEncryptedTx(t *testing.T) {
	msg := &types.MsgCommitTx{
		Sender:      validAddr(),
		TxHash:      validHash(),
		EncryptedTx: make([]byte, 262145), // 256*1024 + 1
	}
	err := msg.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "too large")
}

func TestMsgCommitTx_GetSigners(t *testing.T) {
	addr := sdk.AccAddress([]byte("test_address_padded_"))
	msg := &types.MsgCommitTx{
		Sender:      addr.String(),
		TxHash:      validHash(),
		EncryptedTx: []byte("payload"),
	}
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.True(t, bytes.Equal(signers[0], addr))
}

// ---------------------------------------------------------------------------
// MsgRevealTx
// ---------------------------------------------------------------------------

func TestMsgRevealTx_Valid(t *testing.T) {
	body := []byte("tx-body")
	nonce := []byte("nonce-value")
	h := sha256.Sum256(append(body, nonce...))

	msg := &types.MsgRevealTx{
		Sender:     validAddr(),
		CommitHash: h[:],
		TxBody:     body,
		Nonce:      nonce,
	}
	require.NoError(t, msg.ValidateBasic())
}

func TestMsgRevealTx_InvalidSender(t *testing.T) {
	msg := &types.MsgRevealTx{
		Sender:     "bad",
		CommitHash: validHash(),
		TxBody:     []byte("body"),
		Nonce:      []byte("nonce"),
	}
	require.Error(t, msg.ValidateBasic())
	require.Contains(t, msg.ValidateBasic().Error(), "invalid sender")
}

func TestMsgRevealTx_WrongCommitHashLength(t *testing.T) {
	msg := &types.MsgRevealTx{
		Sender:     validAddr(),
		CommitHash: []byte("short-hash"),
		TxBody:     []byte("body"),
		Nonce:      []byte("nonce"),
	}
	err := msg.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "32 bytes")
}

func TestMsgRevealTx_EmptyBody(t *testing.T) {
	msg := &types.MsgRevealTx{
		Sender:     validAddr(),
		CommitHash: validHash(),
		TxBody:     nil,
		Nonce:      []byte("nonce"),
	}
	err := msg.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty")
}

func TestMsgRevealTx_EmptyNonce(t *testing.T) {
	msg := &types.MsgRevealTx{
		Sender:     validAddr(),
		CommitHash: validHash(),
		TxBody:     []byte("body"),
		Nonce:      nil,
	}
	err := msg.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty")
}

func TestMsgRevealTx_NonceTooLarge(t *testing.T) {
	msg := &types.MsgRevealTx{
		Sender:     validAddr(),
		CommitHash: validHash(),
		TxBody:     []byte("body"),
		Nonce:      make([]byte, 65),
	}
	err := msg.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "too large")
}

func TestMsgRevealTx_HashVerification(t *testing.T) {
	body := []byte("my-transaction-body")
	nonce := []byte("random-nonce-1234")
	h := sha256.Sum256(append(body, nonce...))

	// The hash produced by SHA256(body||nonce) must be exactly 32 bytes
	require.Len(t, h[:], 32)

	// ValidateBasic should pass when commit hash matches the expected length
	msg := &types.MsgRevealTx{
		Sender:     validAddr(),
		CommitHash: h[:],
		TxBody:     body,
		Nonce:      nonce,
	}
	require.NoError(t, msg.ValidateBasic())

	// Wrong hash (different body) should still pass ValidateBasic (hash matching
	// is done by the keeper, not ValidateBasic), but the hash differs
	wrongHash := sha256.Sum256(append([]byte("different-body"), nonce...))
	require.False(t, bytes.Equal(h[:], wrongHash[:]))
}

func TestMsgRevealTx_GetSigners(t *testing.T) {
	addr := sdk.AccAddress([]byte("test_address_padded_"))
	msg := &types.MsgRevealTx{
		Sender:     addr.String(),
		CommitHash: validHash(),
		TxBody:     []byte("body"),
		Nonce:      []byte("nonce"),
	}
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.True(t, bytes.Equal(signers[0], addr))
}

// ---------------------------------------------------------------------------
// Params
// ---------------------------------------------------------------------------

func TestDefaultParams(t *testing.T) {
	p := types.DefaultParams()
	require.True(t, p.FairOrderConfig.EnableCommitReveal)
	require.Equal(t, uint64(3), p.FairOrderConfig.CommitWindow)
	require.Equal(t, uint64(1), p.FairOrderConfig.RevealWindow)
	require.Equal(t, uint64(3), p.FairOrderConfig.MaxTxDelay)
	require.NoError(t, p.Validate())
}

func TestParams_Validate_ZeroCommitWindow(t *testing.T) {
	p := types.DefaultParams()
	p.FairOrderConfig.CommitWindow = 0
	err := p.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "commit window")
}

func TestParams_Validate_ZeroRevealWindow(t *testing.T) {
	p := types.DefaultParams()
	p.FairOrderConfig.RevealWindow = 0
	err := p.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "reveal window")
}

func TestParams_Validate_ZeroMaxTxDelay(t *testing.T) {
	p := types.DefaultParams()
	p.FairOrderConfig.MaxTxDelay = 0
	err := p.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "max tx delay")
}

// ---------------------------------------------------------------------------
// GenesisState
// ---------------------------------------------------------------------------

func TestDefaultGenesis(t *testing.T) {
	gs := types.DefaultGenesis()
	require.NotNil(t, gs)
	require.NoError(t, gs.Validate())
	require.Empty(t, gs.Penalties)
	require.Equal(t, types.DefaultParams(), gs.Params)
}

func TestGenesisState_Validate_EmptyPenalties(t *testing.T) {
	gs := types.DefaultGenesis()
	gs.Penalties = []types.MEVPenalty{}
	require.NoError(t, gs.Validate())
}

func TestGenesisState_Validate_PenaltyEmptyValidatorAddr(t *testing.T) {
	gs := types.DefaultGenesis()
	gs.Penalties = []types.MEVPenalty{
		{
			ValidatorAddr: "",
			PenaltyType:   "reordering",
			Amount:        1,
			BlockHeight:   10,
		},
	}
	err := gs.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty validator address")
}

func TestGenesisState_Validate_PenaltyEmptyType(t *testing.T) {
	gs := types.DefaultGenesis()
	gs.Penalties = []types.MEVPenalty{
		{
			ValidatorAddr: "syreen1validator",
			PenaltyType:   "",
			Amount:        1,
			BlockHeight:   10,
		},
	}
	err := gs.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty penalty type")
}

func TestGenesisState_Validate_BadParams(t *testing.T) {
	gs := types.DefaultGenesis()
	gs.Params.FairOrderConfig.CommitWindow = 0
	require.Error(t, gs.Validate())
}

func TestGenesisState_Validate_ValidPenalties(t *testing.T) {
	gs := types.DefaultGenesis()
	gs.Penalties = []types.MEVPenalty{
		{
			ValidatorAddr: "syreen1abc",
			PenaltyType:   "reordering",
			Amount:        5,
			Evidence:      []byte("proof"),
			BlockHeight:   100,
		},
	}
	require.NoError(t, gs.Validate())
}

// ---------------------------------------------------------------------------
// FairOrderConfig defaults
// ---------------------------------------------------------------------------

func TestDefaultFairOrderConfig(t *testing.T) {
	cfg := types.DefaultFairOrderConfig
	require.True(t, cfg.EnableCommitReveal)
	require.Equal(t, uint64(3), cfg.CommitWindow)
	require.Equal(t, uint64(1), cfg.RevealWindow)
	require.Equal(t, uint64(3), cfg.MaxTxDelay)
}

// ---------------------------------------------------------------------------
// Proto/String methods basic smoke tests
// ---------------------------------------------------------------------------

func TestCommittedTx_String(t *testing.T) {
	ct := types.CommittedTx{Sender: "alice", BlockHeight: 42}
	s := ct.String()
	require.True(t, strings.Contains(s, "alice"))
	require.True(t, strings.Contains(s, "42"))
}

func TestRevealedTx_String(t *testing.T) {
	rt := types.RevealedTx{CommittedHash: []byte{0xab, 0xcd}}
	s := rt.String()
	require.True(t, strings.Contains(s, "abcd"))
}

func TestMEVPenalty_String(t *testing.T) {
	p := types.MEVPenalty{ValidatorAddr: "val1", PenaltyType: "front-run", BlockHeight: 99}
	s := p.String()
	require.True(t, strings.Contains(s, "val1"))
	require.True(t, strings.Contains(s, "99"))
}

func TestFairOrderConfig_String(t *testing.T) {
	cfg := types.DefaultFairOrderConfig
	s := cfg.String()
	require.True(t, strings.Contains(s, "true"))
}
