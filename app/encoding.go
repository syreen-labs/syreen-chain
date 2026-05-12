package app

import (
	"fmt"

	"cosmossdk.io/x/evidence"
	feegrantmodule "cosmossdk.io/x/feegrant/module"
	"cosmossdk.io/x/upgrade"

	"cosmossdk.io/x/tx/signing"
	"github.com/cosmos/cosmos-sdk/client"
	"google.golang.org/protobuf/reflect/protoreflect"
	googleproto "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/dynamicpb"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/gogoproto/proto"

	ibc "github.com/cosmos/ibc-go/v10/modules/core"
	ibctm "github.com/cosmos/ibc-go/v10/modules/light-clients/07-tendermint"
	ibctransfer "github.com/cosmos/ibc-go/v10/modules/apps/transfer"

	"syreen/x/abstractaccount"
	"syreen/x/compute"
	"syreen/x/dex"
	"syreen/x/feemarket"
	"syreen/x/intent"
	"syreen/x/mevprotection"
	"syreen/x/tokenfactory"
	"syreen/x/evm"
	// Archived module imports removed (aiagent, streampay, travelescrow, loyalty, travelinsurance, reputation)
	"syreen/x/identity"
	"syreen/x/payments"
	"syreen/x/farming"
	"syreen/x/perps"
	"syreen/x/lending"
	"syreen/x/predict"
	"syreen/x/launchpad"
	"syreen/x/clmm"
	"syreen/x/portfolio"
	"syreen/x/flashloan"
	"syreen/x/vault"
	"syreen/x/options"
	"syreen/x/aiagent"
)

// EncodingConfig specifies the concrete encoding types to use for a given app.
type EncodingConfig struct {
	InterfaceRegistry types.InterfaceRegistry
	Codec             codec.Codec
	TxConfig          client.TxConfig
	Amino             *codec.LegacyAmino
}

// MakeEncodingConfig creates an EncodingConfig for the app.
func MakeEncodingConfig() EncodingConfig {
	amino := codec.NewLegacyAmino()
	signingOpts := signing.Options{
		AddressCodec:          address.NewBech32Codec(Bech32MainPrefix),
		ValidatorAddressCodec: address.NewBech32Codec(Bech32MainPrefix + "valoper"),
	}
	// Register custom GetSigners for MsgEthereumTx — extracts signer from the "from" field
	signingOpts.DefineCustomGetSigners(
		protoreflect.FullName("syreen.evm.MsgEthereumTx"),
		func(msg googleproto.Message) ([][]byte, error) {
			// The message comes as a dynamic protobuf message — extract "from" field
			dynMsg, ok := msg.(*dynamicpb.Message)
			if !ok {
				return nil, fmt.Errorf("unexpected message type for MsgEthereumTx signer extraction")
			}
			fromField := dynMsg.Descriptor().Fields().ByName("from")
			if fromField == nil {
				return nil, fmt.Errorf("MsgEthereumTx missing 'from' field")
			}
			fromStr := dynMsg.Get(fromField).String()
			if fromStr == "" {
				return nil, nil
			}
			addrCodec := address.NewBech32Codec(Bech32MainPrefix)
			bz, err := addrCodec.StringToBytes(fromStr)
			if err != nil {
				return nil, err
			}
			return [][]byte{bz}, nil
		},
	)
	// Register custom GetSigners for tokenfactory messages
	for _, msgName := range []protoreflect.FullName{
		"syreen.tokenfactory.MsgCreateDenom",
		"syreen.tokenfactory.MsgMint",
		"syreen.tokenfactory.MsgBurn",
		"syreen.tokenfactory.MsgChangeAdmin",
	} {
		name := msgName
		signingOpts.DefineCustomGetSigners(
			name,
			func(msg googleproto.Message) ([][]byte, error) {
				dynMsg, ok := msg.(*dynamicpb.Message)
				if !ok {
					return nil, fmt.Errorf("unexpected message type for %s signer extraction", name)
				}
				senderField := dynMsg.Descriptor().Fields().ByName("sender")
				if senderField == nil {
					return nil, fmt.Errorf("%s missing 'sender' field", name)
				}
				senderStr := dynMsg.Get(senderField).String()
				if senderStr == "" {
					return nil, nil
				}
				addrCodec := address.NewBech32Codec(Bech32MainPrefix)
				bz, err := addrCodec.StringToBytes(senderStr)
				if err != nil {
					return nil, err
				}
				return [][]byte{bz}, nil
			},
		)
	}

	// Register custom GetSigners for DEX messages (sender field)
	for _, msgName := range []protoreflect.FullName{
		"syreen.dex.MsgCreatePool",
		"syreen.dex.MsgAddLiquidity",
		"syreen.dex.MsgRemoveLiquidity",
		"syreen.dex.MsgSwap",
		"syreen.dex.MsgMultiHopSwap",
		"syreen.dex.MsgSetPoolFeeConfig",
	} {
		name := msgName
		signingOpts.DefineCustomGetSigners(
			name,
			func(msg googleproto.Message) ([][]byte, error) {
				dynMsg, ok := msg.(*dynamicpb.Message)
				if !ok {
					return nil, fmt.Errorf("unexpected message type for %s signer extraction", name)
				}
				senderField := dynMsg.Descriptor().Fields().ByName("sender")
				if senderField == nil {
					return nil, fmt.Errorf("%s missing 'sender' field", name)
				}
				senderStr := dynMsg.Get(senderField).String()
				if senderStr == "" {
					return nil, nil
				}
				addrCodec := address.NewBech32Codec(Bech32MainPrefix)
				bz, err := addrCodec.StringToBytes(senderStr)
				if err != nil {
					return nil, err
				}
				return [][]byte{bz}, nil
			},
		)
	}

	// Register custom GetSigners for DEX order/referral/copytrade messages
	for _, pair := range []struct {
		name  protoreflect.FullName
		field string
	}{
		{"syreen.dex.MsgPlaceOrder", "creator"},
		{"syreen.dex.MsgCancelOrder", "creator"},
		{"syreen.dex.MsgModifyOrder", "creator"},
		{"syreen.dex.MsgCreateReferralCode", "creator"},
		{"syreen.dex.MsgRegisterReferral", "user"},
		{"syreen.dex.MsgClaimReferralRewards", "referrer"},
		{"syreen.dex.MsgFollowTrader", "follower"},
		{"syreen.dex.MsgUnfollowTrader", "follower"},
		{"syreen.dex.MsgUpdateCopySettings", "follower"},
	} {
		p := pair
		signingOpts.DefineCustomGetSigners(
			p.name,
			func(msg googleproto.Message) ([][]byte, error) {
				dynMsg, ok := msg.(*dynamicpb.Message)
				if !ok {
					return nil, fmt.Errorf("unexpected message type for %s signer extraction", p.name)
				}
				f := dynMsg.Descriptor().Fields().ByName(protoreflect.Name(p.field))
				if f == nil {
					return nil, fmt.Errorf("%s missing '%s' field", p.name, p.field)
				}
				s := dynMsg.Get(f).String()
				if s == "" {
					return nil, nil
				}
				addrCodec := address.NewBech32Codec(Bech32MainPrefix)
				bz, err := addrCodec.StringToBytes(s)
				if err != nil {
					return nil, err
				}
				return [][]byte{bz}, nil
			},
		)
	}

	// Register custom GetSigners for intent messages
	for _, pair := range []struct {
		name  protoreflect.FullName
		field string
	}{
		{"syreen.intent.MsgSubmitIntent", "creator"},
		{"syreen.intent.MsgRegisterSolver", "address"},
		{"syreen.intent.MsgDeregisterSolver", "address"},
		{"syreen.intent.MsgSubmitSolution", "solver_addr"},
		{"syreen.intent.MsgFulfillIntent", "sender"},
		{"syreen.intent.MsgSubmitChain", "creator"},
		{"syreen.intent.MsgCancelChain", "creator"},
		{"syreen.intent.MsgCreateStrategy", "creator"},
		{"syreen.intent.MsgCancelStrategy", "creator"},
	} {
		p := pair
		signingOpts.DefineCustomGetSigners(
			p.name,
			func(msg googleproto.Message) ([][]byte, error) {
				dynMsg, ok := msg.(*dynamicpb.Message)
				if !ok {
					return nil, fmt.Errorf("unexpected message type for %s signer extraction", p.name)
				}
				f := dynMsg.Descriptor().Fields().ByName(protoreflect.Name(p.field))
				if f == nil {
					return nil, fmt.Errorf("%s missing '%s' field", p.name, p.field)
				}
				s := dynMsg.Get(f).String()
				if s == "" {
					return nil, nil
				}
				addrCodec := address.NewBech32Codec(Bech32MainPrefix)
				bz, err := addrCodec.StringToBytes(s)
				if err != nil {
					return nil, err
				}
				return [][]byte{bz}, nil
			},
		)
	}

	// Register custom GetSigners for farming + perps messages
	for _, pair := range []struct {
		name  protoreflect.FullName
		field string
	}{
		{"syreen.farming.MsgStake", "sender"},
		{"syreen.farming.MsgUnstake", "sender"},
		{"syreen.farming.MsgClaimReward", "sender"},
		{"syreen.farming.MsgCreateFarm", "authority"},
		{"syreen.farming.MsgUpdateFarm", "authority"},
		{"syreen.perps.MsgOpenPosition", "sender"},
		{"syreen.perps.MsgClosePosition", "sender"},
		{"syreen.perps.MsgAddMargin", "sender"},
		{"syreen.perps.MsgRemoveMargin", "sender"},
		{"syreen.perps.MsgCreateMarket", "authority"},
		{"syreen.lending.MsgDeposit", "sender"},
		{"syreen.lending.MsgWithdraw", "sender"},
		{"syreen.lending.MsgBorrow", "sender"},
		{"syreen.lending.MsgRepay", "sender"},
		{"syreen.lending.MsgLiquidate", "liquidator"},
		{"syreen.lending.MsgCreateLendingPool", "authority"},
		{"syreen.predict.MsgCreateMarket", "creator"},
		{"syreen.predict.MsgBuyShares", "sender"},
		{"syreen.predict.MsgSellShares", "sender"},
		{"syreen.predict.MsgResolveMarket", "resolver"},
		{"syreen.predict.MsgClaimWinnings", "sender"},
		{"syreen.launchpad.MsgCreateLaunch", "creator"},
		{"syreen.launchpad.MsgContribute", "sender"},
		{"syreen.launchpad.MsgClaimTokens", "sender"},
		{"syreen.launchpad.MsgClaimRefund", "sender"},
		{"syreen.launchpad.MsgFinalizeLaunch", "authority"},
		{"syreen.clmm.MsgCreateCLPool", "sender"},
		{"syreen.clmm.MsgCreatePosition", "sender"},
		{"syreen.clmm.MsgAddLiquidity", "sender"},
		{"syreen.clmm.MsgRemoveLiquidity", "sender"},
		{"syreen.clmm.MsgCollectFees", "sender"},
		{"syreen.clmm.MsgCLSwap", "sender"},
		{"syreen.flashloan.MsgFlashLoan", "sender"},
		{"syreen.flashloan.MsgCreateFlashPool", "authority"},
		{"syreen.flashloan.MsgFundFlashPool", "sender"},
		{"syreen.flashloan.MsgWithdrawFlashPool", "sender"},
		{"syreen.vault.MsgCreateVault", "creator"},
		{"syreen.vault.MsgDepositVault", "sender"},
		{"syreen.vault.MsgWithdrawVault", "sender"},
		{"syreen.vault.MsgCompoundVault", "sender"},
		{"syreen.vault.MsgUpdateVaultStrategy", "creator"},
		{"syreen.options.MsgWriteOption", "writer"},
		{"syreen.options.MsgBuyOption", "buyer"},
		{"syreen.options.MsgExerciseOption", "buyer"},
		{"syreen.options.MsgCancelOption", "writer"},
		{"syreen.aiagent.MsgCreateAgent", "owner"},
		{"syreen.aiagent.MsgFundAgent", "owner"},
		{"syreen.aiagent.MsgWithdrawAgentFunds", "owner"},
		{"syreen.aiagent.MsgPauseAgent", "owner"},
		{"syreen.aiagent.MsgResumeAgent", "owner"},
		{"syreen.aiagent.MsgUpdateAgentStrategy", "owner"},
		{"syreen.portfolio.MsgRecordTrade", "sender"},
		{"syreen.portfolio.MsgCreateCompetition", "creator"},
		{"syreen.portfolio.MsgJoinCompetition", "sender"},
		{"syreen.portfolio.MsgEndCompetition", "authority"},
		{"syreen.portfolio.MsgUpdatePortfolio", "sender"},
		// mevprotection
		{"syreen.mevprotection.MsgCommitTx", "sender"},
		{"syreen.mevprotection.MsgRevealTx", "sender"},
		// abstractaccount
		{"syreen.abstractaccount.MsgCreateSmartAccount", "sender"},
		{"syreen.abstractaccount.MsgCreateSessionKey", "granter"},
		{"syreen.abstractaccount.MsgRevokeSessionKey", "granter"},
		{"syreen.abstractaccount.MsgInitiateRecovery", "guardian"},
		{"syreen.abstractaccount.MsgApproveRecovery", "guardian"},
		{"syreen.abstractaccount.MsgExecuteRecovery", "sender"},
		{"syreen.abstractaccount.MsgSponsorGas", "sponsor"},
		{"syreen.abstractaccount.MsgBatchExecute", "sender"},
		// compute
		{"syreen.compute.MsgStoreCode", "sender"},
		{"syreen.compute.MsgInstantiateContract", "sender"},
		{"syreen.compute.MsgExecuteContract", "sender"},
		{"syreen.compute.MsgMigrateContract", "sender"},
		{"syreen.compute.MsgUpdateAdmin", "sender"},
		// identity
		{"syreen.identity.MsgRegisterIdentity", "address"},
		{"syreen.identity.MsgVerifyIdentity", "verifier"},
		{"syreen.identity.MsgRejectIdentity", "verifier"},
		{"syreen.identity.MsgRevokeIdentity", "authority"},
		{"syreen.identity.MsgUpdateIdentity", "address"},
		{"syreen.identity.MsgRegisterVerifier", "authority"},
		{"syreen.identity.MsgDeactivateVerifier", "authority"},
		{"syreen.identity.MsgIncrementBookings", "authority"},
		// payments
		{"syreen.payments.MsgCreateInvoice", "creator"},
		{"syreen.payments.MsgPayInvoice", "payer"},
		{"syreen.payments.MsgRefundPayment", "authority"},
		{"syreen.payments.MsgSetExchangeRate", "authority"},
		{"syreen.payments.MsgWithdrawEarnings", "address"},
		{"syreen.payments.MsgCancelInvoice", "creator"},
	} {
		p := pair
		signingOpts.DefineCustomGetSigners(
			p.name,
			func(msg googleproto.Message) ([][]byte, error) {
				dynMsg, ok := msg.(*dynamicpb.Message)
				if !ok {
					return nil, fmt.Errorf("unexpected message type for %s signer extraction", p.name)
				}
				f := dynMsg.Descriptor().Fields().ByName(protoreflect.Name(p.field))
				if f == nil {
					return nil, fmt.Errorf("%s missing '%s' field", p.name, p.field)
				}
				s := dynMsg.Get(f).String()
				if s == "" {
					return nil, nil
				}
				addrCodec := address.NewBech32Codec(Bech32MainPrefix)
				bz, err := addrCodec.StringToBytes(s)
				if err != nil {
					return nil, err
				}
				return [][]byte{bz}, nil
			},
		)
	}

	interfaceRegistry, err := types.NewInterfaceRegistryWithOptions(types.InterfaceRegistryOptions{
		ProtoFiles:     proto.HybridResolver,
		SigningOptions: signingOpts,
	})
	if err != nil {
		panic(err)
	}
	cdc := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(cdc, tx.DefaultSignModes)

	std.RegisterLegacyAminoCodec(amino)
	std.RegisterInterfaces(interfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(amino)
	ModuleBasics.RegisterInterfaces(interfaceRegistry)

	return EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             cdc,
		TxConfig:          txCfg,
		Amino:             amino,
	}
}

// NewDefaultGenesisState generates the default state for the application
// with usyreen as the bond denom. Delegates to DefaultGenesisState in genesis.go.
func NewDefaultGenesisState(cdc codec.JSONCodec) GenesisState {
	return DefaultGenesisState(cdc)
}

// ModuleBasics defines the module BasicManager is in charge of setting up basic,
// non-dependant module elements, such as codec registration
// and genesis verification.
var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
	bank.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distr.AppModuleBasic{},
	gov.NewAppModuleBasic(
		[]govclient.ProposalHandler{
			paramsclient.ProposalHandler,
		},
	),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	vesting.AppModuleBasic{},
	consensus.AppModuleBasic{},
	evidence.AppModuleBasic{},
	feegrantmodule.AppModuleBasic{},
	authzmodule.AppModuleBasic{},
	ibc.AppModuleBasic{},
	ibctransfer.AppModuleBasic{},
	ibctm.AppModuleBasic{},
	tokenfactory.AppModule{},
	feemarket.AppModule{},
	mevprotection.AppModule{},
	abstractaccount.AppModule{},
	intent.AppModule{},
	compute.AppModule{},
	evm.AppModule{},
	dex.AppModule{},
	identity.AppModule{},
	payments.AppModule{},
	farming.AppModule{},
	perps.AppModule{},
	lending.AppModule{},
	predict.AppModule{},
	launchpad.AppModule{},
	clmm.AppModule{},
	flashloan.AppModule{},
	vault.AppModule{},
	portfolio.AppModule{},
	options.AppModule{},
	aiagent.AppModule{},
)
