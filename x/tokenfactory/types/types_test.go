package types_test

import (
	"strings"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"syreen/x/tokenfactory/types"
)

func init() {
	// Set bech32 prefix so addresses encode as syreen1...
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("syreen", "syreenpub")
	config.Seal()
}

func validAddr() string {
	addr := sdk.AccAddress([]byte("test_address_padded_")) // 20 bytes
	return addr.String()
}

func validAddr2() string {
	addr := sdk.AccAddress([]byte("second_addr_padding_")) // 20 bytes
	return addr.String()
}

// ---------------------------------------------------------------------------
// MsgCreateDenom
// ---------------------------------------------------------------------------

func TestMsgCreateDenom_Valid(t *testing.T) {
	msg := &types.MsgCreateDenom{Sender: validAddr(), Subdenom: "mytoken"}
	require.NoError(t, msg.ValidateBasic())
}

func TestMsgCreateDenom_EmptySender(t *testing.T) {
	msg := &types.MsgCreateDenom{Sender: "", Subdenom: "mytoken"}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrInvalidCreator)
}

func TestMsgCreateDenom_InvalidSender(t *testing.T) {
	msg := &types.MsgCreateDenom{Sender: "notabech32address", Subdenom: "mytoken"}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrInvalidCreator)
}

func TestMsgCreateDenom_EmptySubdenom(t *testing.T) {
	msg := &types.MsgCreateDenom{Sender: validAddr(), Subdenom: ""}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrSubdenomTooShort)
}

func TestMsgCreateDenom_TooLongSubdenom(t *testing.T) {
	long := strings.Repeat("a", types.MaxSubdenomLength+1)
	msg := &types.MsgCreateDenom{Sender: validAddr(), Subdenom: long}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrSubdenomTooLong)
}

func TestMsgCreateDenom_MaxLengthSubdenom(t *testing.T) {
	exact := strings.Repeat("a", types.MaxSubdenomLength)
	msg := &types.MsgCreateDenom{Sender: validAddr(), Subdenom: exact}
	require.NoError(t, msg.ValidateBasic())
}

func TestMsgCreateDenom_InvalidChars(t *testing.T) {
	for _, bad := range []string{"my token", "my@token", "my#token", "my/token", "my!token"} {
		msg := &types.MsgCreateDenom{Sender: validAddr(), Subdenom: bad}
		err := msg.ValidateBasic()
		require.Error(t, err, "subdenom %q should be invalid", bad)
		require.Contains(t, err.Error(), "invalid character")
	}
}

func TestMsgCreateDenom_ValidSpecialChars(t *testing.T) {
	for _, good := range []string{"my_token", "my-token", "my.token", "my_token-v2.0", "ABC123"} {
		msg := &types.MsgCreateDenom{Sender: validAddr(), Subdenom: good}
		require.NoError(t, msg.ValidateBasic(), "subdenom %q should be valid", good)
	}
}

func TestMsgCreateDenom_GetSigners(t *testing.T) {
	addr := validAddr()
	msg := &types.MsgCreateDenom{Sender: addr, Subdenom: "x"}
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, addr, signers[0].String())
}

// ---------------------------------------------------------------------------
// MsgMint
// ---------------------------------------------------------------------------

func TestMsgMint_Valid(t *testing.T) {
	denom := types.GetDenom(validAddr(), "mytoken")
	msg := &types.MsgMint{
		Sender: validAddr(),
		Amount: sdk.NewInt64Coin(denom, 1000),
		MintTo: "",
	}
	require.NoError(t, msg.ValidateBasic())
}

func TestMsgMint_InvalidSender(t *testing.T) {
	msg := &types.MsgMint{
		Sender: "bad",
		Amount: sdk.NewInt64Coin("factory/x/y", 1000),
	}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrInvalidCreator)
}

func TestMsgMint_ZeroAmount(t *testing.T) {
	msg := &types.MsgMint{
		Sender: validAddr(),
		Amount: sdk.NewInt64Coin("factory/x/y", 0),
	}
	require.Error(t, msg.ValidateBasic())
	require.Contains(t, msg.ValidateBasic().Error(), "invalid mint amount")
}

func TestMsgMint_InvalidMintTo(t *testing.T) {
	msg := &types.MsgMint{
		Sender: validAddr(),
		Amount: sdk.NewInt64Coin("factory/x/y", 100),
		MintTo: "not_a_valid_address",
	}
	err := msg.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid mint_to address")
}

func TestMsgMint_EmptyMintToIsOk(t *testing.T) {
	msg := &types.MsgMint{
		Sender: validAddr(),
		Amount: sdk.NewInt64Coin("factory/x/y", 100),
		MintTo: "",
	}
	require.NoError(t, msg.ValidateBasic())
}

func TestMsgMint_GetSigners(t *testing.T) {
	addr := validAddr()
	msg := &types.MsgMint{Sender: addr, Amount: sdk.NewInt64Coin("factory/x/y", 1)}
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, addr, signers[0].String())
}

// ---------------------------------------------------------------------------
// MsgBurn
// ---------------------------------------------------------------------------

func TestMsgBurn_Valid(t *testing.T) {
	msg := &types.MsgBurn{
		Sender:   validAddr(),
		Amount:   sdk.NewInt64Coin("factory/x/y", 500),
		BurnFrom: "",
	}
	require.NoError(t, msg.ValidateBasic())
}

func TestMsgBurn_InvalidSender(t *testing.T) {
	msg := &types.MsgBurn{
		Sender: "bad",
		Amount: sdk.NewInt64Coin("factory/x/y", 100),
	}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrInvalidCreator)
}

func TestMsgBurn_ZeroAmount(t *testing.T) {
	msg := &types.MsgBurn{
		Sender: validAddr(),
		Amount: sdk.NewInt64Coin("factory/x/y", 0),
	}
	err := msg.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid burn amount")
}

func TestMsgBurn_InvalidBurnFrom(t *testing.T) {
	msg := &types.MsgBurn{
		Sender:   validAddr(),
		Amount:   sdk.NewInt64Coin("factory/x/y", 100),
		BurnFrom: "invalid_address",
	}
	err := msg.ValidateBasic()
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid burn_from address")
}

func TestMsgBurn_EmptyBurnFromIsOk(t *testing.T) {
	msg := &types.MsgBurn{
		Sender:   validAddr(),
		Amount:   sdk.NewInt64Coin("factory/x/y", 100),
		BurnFrom: "",
	}
	require.NoError(t, msg.ValidateBasic())
}

func TestMsgBurn_GetSigners(t *testing.T) {
	addr := validAddr()
	msg := &types.MsgBurn{Sender: addr, Amount: sdk.NewInt64Coin("factory/x/y", 1)}
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, addr, signers[0].String())
}

// ---------------------------------------------------------------------------
// MsgChangeAdmin
// ---------------------------------------------------------------------------

func TestMsgChangeAdmin_Valid(t *testing.T) {
	msg := &types.MsgChangeAdmin{
		Sender:   validAddr(),
		Denom:    "factory/x/y",
		NewAdmin: validAddr2(),
	}
	require.NoError(t, msg.ValidateBasic())
}

func TestMsgChangeAdmin_InvalidSender(t *testing.T) {
	msg := &types.MsgChangeAdmin{
		Sender:   "bad",
		Denom:    "factory/x/y",
		NewAdmin: validAddr(),
	}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrInvalidCreator)
}

func TestMsgChangeAdmin_InvalidNewAdmin(t *testing.T) {
	msg := &types.MsgChangeAdmin{
		Sender:   validAddr(),
		Denom:    "factory/x/y",
		NewAdmin: "invalid",
	}
	require.ErrorIs(t, msg.ValidateBasic(), types.ErrInvalidAdmin)
}

func TestMsgChangeAdmin_GetSigners(t *testing.T) {
	addr := validAddr()
	msg := &types.MsgChangeAdmin{Sender: addr, Denom: "x", NewAdmin: validAddr2()}
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, addr, signers[0].String())
}

// ---------------------------------------------------------------------------
// GetDenom / DeconstructDenom
// ---------------------------------------------------------------------------

func TestGetDenom(t *testing.T) {
	denom := types.GetDenom("syreen1abc", "mytoken")
	require.Equal(t, "factory/syreen1abc/mytoken", denom)
}

func TestDeconstructDenom_Valid(t *testing.T) {
	creator, subdenom, err := types.DeconstructDenom("factory/syreen1abc/mytoken")
	require.NoError(t, err)
	require.Equal(t, "syreen1abc", creator)
	require.Equal(t, "mytoken", subdenom)
}

func TestDeconstructDenom_NotFactory(t *testing.T) {
	_, _, err := types.DeconstructDenom("usyreen")
	require.Error(t, err)
	require.Contains(t, err.Error(), "not a factory denom")
}

func TestDeconstructDenom_MissingParts(t *testing.T) {
	_, _, err := types.DeconstructDenom("factory/onlyonepart")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid factory denom")
}

func TestDeconstructDenom_Roundtrip(t *testing.T) {
	original := types.GetDenom("syreen1creator", "subdenom")
	creator, subdenom, err := types.DeconstructDenom(original)
	require.NoError(t, err)
	require.Equal(t, "syreen1creator", creator)
	require.Equal(t, "subdenom", subdenom)
	require.Equal(t, original, types.GetDenom(creator, subdenom))
}

// ---------------------------------------------------------------------------
// Params
// ---------------------------------------------------------------------------

func TestDefaultParams(t *testing.T) {
	p := types.DefaultParams()
	require.True(t, p.DenomCreationFee.IsValid())
	require.NoError(t, p.Validate())
}

func TestParams_Validate_InvalidFee(t *testing.T) {
	p := types.Params{
		DenomCreationFee: sdk.Coins{sdk.Coin{Denom: "", Amount: math.NewInt(0)}},
	}
	require.Error(t, p.Validate())
}

func TestParams_Validate_ZeroFee(t *testing.T) {
	// Empty coins (no fee) is still valid per Coins.IsValid()
	p := types.Params{DenomCreationFee: sdk.Coins{}}
	require.NoError(t, p.Validate())
}

// ---------------------------------------------------------------------------
// GenesisState
// ---------------------------------------------------------------------------

func TestDefaultGenesisState(t *testing.T) {
	gs := types.DefaultGenesis()
	require.NotNil(t, gs)
	require.Equal(t, types.DefaultParams(), gs.Params)
	require.Empty(t, gs.FactoryDenoms)
	require.NoError(t, gs.Validate())
}

func TestGenesisState_Validate_DuplicateDenom(t *testing.T) {
	gs := types.GenesisState{
		Params: types.DefaultParams(),
		FactoryDenoms: []types.GenesisDenom{
			{Denom: "factory/creator/a", AuthorityMetadata: types.DenomAuthorityMetadata{Admin: "creator"}},
			{Denom: "factory/creator/a", AuthorityMetadata: types.DenomAuthorityMetadata{Admin: "creator"}},
		},
	}
	err := gs.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "duplicate denom")
}

func TestGenesisState_Validate_UniqueDenoms(t *testing.T) {
	gs := types.GenesisState{
		Params: types.DefaultParams(),
		FactoryDenoms: []types.GenesisDenom{
			{Denom: "factory/creator/a", AuthorityMetadata: types.DenomAuthorityMetadata{Admin: "creator"}},
			{Denom: "factory/creator/b", AuthorityMetadata: types.DenomAuthorityMetadata{Admin: "creator"}},
		},
	}
	require.NoError(t, gs.Validate())
}

func TestGenesisState_Validate_InvalidParams(t *testing.T) {
	gs := types.GenesisState{
		Params: types.Params{
			DenomCreationFee: sdk.Coins{sdk.Coin{Denom: "", Amount: math.NewInt(0)}},
		},
		FactoryDenoms: nil,
	}
	require.Error(t, gs.Validate())
}
