package aiagent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"syreen/config"
	"syreen/x/aiagent/client/cli"
	"syreen/x/aiagent/keeper"
	"syreen/x/aiagent/types"
)

var (
	_ module.HasGenesis         = AppModule{}
	_ module.AppModule          = AppModule{}
	_ appmodule.HasBeginBlocker = AppModule{}
)

type AppModule struct {
	keeper *keeper.Keeper
}

func NewAppModule(k *keeper.Keeper) AppModule {
	return AppModule{keeper: k}
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
func (AppModule) GetTxCmd() *cobra.Command                                        { return cli.GetTxCmd() }
func (AppModule) GetQueryCmd() *cobra.Command                                     { return cli.GetQueryCmd() }

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), *am.keeper)
}

func (am AppModule) InitGenesis(ctx sdk.Context, _ codec.JSONCodec, data json.RawMessage) {
	var gs types.GenesisState
	if err := json.Unmarshal(data, &gs); err != nil {
		panic(fmt.Errorf("failed to unmarshal %s genesis: %w", types.ModuleName, err))
	}
	for _, agent := range gs.Agents {
		am.keeper.SetAgent(ctx, agent)
	}
	if gs.NextID > 0 {
		am.keeper.SetNextAgentID(ctx, gs.NextID)
	}
	return
}

func (am AppModule) ExportGenesis(ctx sdk.Context, _ codec.JSONCodec) json.RawMessage {
	gs := types.GenesisState{
		Agents: am.keeper.GetAllAgents(ctx),
		NextID: am.keeper.GetNextAgentID(ctx),
	}
	bz, _ := json.Marshal(gs)
	return bz
}

func (AppModule) ConsensusVersion() uint64 { return 1 }

func (am AppModule) BeginBlock(ctx context.Context) error {
	if !config.IsModuleEnabled("aiagent") {
		return nil
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	am.keeper.BeginBlockHandler(sdkCtx)
	return nil
}

func (am AppModule) EndBlock(_ context.Context) error { return nil }

func (am AppModule) IsOnePerModuleType() {}
func (am AppModule) IsAppModule()        {}
