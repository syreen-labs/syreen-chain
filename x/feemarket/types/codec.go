package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
)

func RegisterInterfaces(_ types.InterfaceRegistry) {
	// No messages to register yet; governance messages can be added later
}

func RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {
	// No messages to register yet
}
