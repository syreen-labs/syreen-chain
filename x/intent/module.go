package intent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"cosmossdk.io/core/appmodule"
	"syreen/config"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"syreen/x/intent/client/cli"
	"syreen/x/intent/keeper"
	"syreen/x/intent/types"
)

var (
	_ module.HasGenesis         = AppModule{}
	_ module.AppModule          = AppModule{}
	_ appmodule.HasBeginBlocker = AppModule{}
)

// AppModule implements the AppModule interface
type AppModule struct {
	keeper *keeper.Keeper
}

func NewAppModule(keeper *keeper.Keeper) AppModule {
	return AppModule{
		keeper: keeper,
	}
}

func (AppModule) Name() string { return types.ModuleName }

func (AppModule) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

func (AppModule) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

func (AppModule) DefaultGenesis(_ codec.JSONCodec) json.RawMessage {
	bz, _ := json.Marshal(types.DefaultGenesis())
	return bz
}

func (AppModule) ValidateGenesis(_ codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	if err := json.Unmarshal(bz, &data); err != nil {
		panic(fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err))
	}
	return data.Validate()
}

func (AppModule) RegisterGRPCGatewayRoutes(_ client.Context, _ *runtime.ServeMux) {}

func (AppModule) GetTxCmd() *cobra.Command { return cli.GetTxCmd() }

func (AppModule) GetQueryCmd() *cobra.Command { return cli.GetQueryCmd() }

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(*am.keeper))
}

func (am AppModule) InitGenesis(ctx sdk.Context, _ codec.JSONCodec, data json.RawMessage) {
	var genesisState types.GenesisState
	if err := json.Unmarshal(data, &genesisState); err != nil {
		panic(fmt.Errorf("failed to unmarshal %s genesis: %w", types.ModuleName, err))
	}
	am.keeper.InitGenesis(ctx, genesisState)
	return
}

func (am AppModule) ExportGenesis(ctx sdk.Context, _ codec.JSONCodec) json.RawMessage {
	gs := am.keeper.ExportGenesis(ctx)
	bz, _ := json.Marshal(gs)
	return bz
}

func (AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock expires intents, completes solver unbonding, auto-fulfills intents
// whose solving window has elapsed, and executes trading intents when price conditions are met.
//
// NOTE: signature uses context.Context (not sdk.Context) so that this method
// actually satisfies appmodule.HasBeginBlocker — historically this file used
// sdk.Context, which meant the method was a dead no-op and none of the
// intent BeginBlock work (strategy syncing, trading-intent execution, param
// migrations, ...) ever ran.
func (am AppModule) BeginBlock(ctx context.Context) error {
	if !config.IsModuleEnabled("intent") {
		return nil
	}
	// One-shot param migration: early genesis versions of this module
	// persisted intent params without the EnableChains / MaxChainSteps
	// fields, so they decoded as zero values and permanently disabled the
	// strategy engine. If we detect that exact stale signature, rewrite
	// params to the current defaults. This is idempotent — once the fields
	// are populated, the condition never matches again, and the normal
	// param path (including any future gov-controlled updates) takes over.
	p := am.keeper.GetParams(ctx)
	if !p.EnableChains && p.MaxChainSteps == 0 {
		p.EnableChains = types.DefaultEnableChains
		p.MaxChainSteps = types.DefaultMaxChainSteps
		_ = am.keeper.SetParams(ctx, p)
	}

	// Second one-shot migration: old genesis persisted IntentExpiryBlocks=200
	// which caps chain expiry at 200 * numSteps. For a 2-step safe_accumulate
	// that's 400 blocks — far less than a single 600-block DCA interval,
	// so multi-tranche strategies could never run to completion. Bump to
	// the current default (20000) if we detect the stale value. Idempotent:
	// once raised, subsequent gov updates still take precedence.
	if p.IntentExpiryBlocks <= 200 {
		p.IntentExpiryBlocks = types.DefaultIntentExpiryBlocks
		_ = am.keeper.SetParams(ctx, p)
	}

	am.keeper.ExpireIntents(ctx)
	am.keeper.CompleteSolverUnbonding(ctx)
	am.keeper.AutoFulfillIntents(ctx)
	am.keeper.ExecuteTradingIntents(ctx)
	am.keeper.ProcessChainConditions(ctx)
	am.keeper.SyncStrategyStatuses(ctx)
	return nil
}

func (am AppModule) IsOnePerModuleType() {}
func (am AppModule) IsAppModule()        {}
