package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreatePool{},
		&MsgAddLiquidity{},
		&MsgRemoveLiquidity{},
		&MsgSwap{},
		&MsgCreateReferralCode{},
		&MsgRegisterReferral{},
		&MsgPlaceOrder{},
		&MsgCancelOrder{},
		&MsgModifyOrder{},
		&MsgFollowTrader{},
		&MsgUnfollowTrader{},
		&MsgUpdateCopySettings{},
		&MsgClaimReferralRewards{},
		&MsgMultiHopSwap{},
		&MsgSetPoolFeeConfig{},
	)
}

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreatePool{}, "syreen/dex/MsgCreatePool", nil)
	cdc.RegisterConcrete(&MsgAddLiquidity{}, "syreen/dex/MsgAddLiquidity", nil)
	cdc.RegisterConcrete(&MsgRemoveLiquidity{}, "syreen/dex/MsgRemoveLiquidity", nil)
	cdc.RegisterConcrete(&MsgSwap{}, "syreen/dex/MsgSwap", nil)
	cdc.RegisterConcrete(&MsgPlaceOrder{}, "syreen/dex/MsgPlaceOrder", nil)
	cdc.RegisterConcrete(&MsgCancelOrder{}, "syreen/dex/MsgCancelOrder", nil)
	cdc.RegisterConcrete(&MsgModifyOrder{}, "syreen/dex/MsgModifyOrder", nil)
	cdc.RegisterConcrete(&MsgCreateReferralCode{}, "syreen/dex/MsgCreateReferralCode", nil)
	cdc.RegisterConcrete(&MsgRegisterReferral{}, "syreen/dex/MsgRegisterReferral", nil)
	cdc.RegisterConcrete(&MsgClaimReferralRewards{}, "syreen/dex/MsgClaimReferralRewards", nil)
	cdc.RegisterConcrete(&MsgFollowTrader{}, "syreen/dex/MsgFollowTrader", nil)
	cdc.RegisterConcrete(&MsgUnfollowTrader{}, "syreen/dex/MsgUnfollowTrader", nil)
	cdc.RegisterConcrete(&MsgUpdateCopySettings{}, "syreen/dex/MsgUpdateCopySettings", nil)
	cdc.RegisterConcrete(&MsgMultiHopSwap{}, "syreen/dex/MsgMultiHopSwap", nil)
	cdc.RegisterConcrete(&MsgSetPoolFeeConfig{}, "syreen/dex/MsgSetPoolFeeConfig", nil)
}
