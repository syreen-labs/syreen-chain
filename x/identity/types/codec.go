package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegisterIdentity{},
		&MsgVerifyIdentity{},
		&MsgRejectIdentity{},
		&MsgRevokeIdentity{},
		&MsgUpdateIdentity{},
		&MsgRegisterVerifier{},
		&MsgDeactivateVerifier{},
		&MsgIncrementBookings{},
	)
}

func RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {}
