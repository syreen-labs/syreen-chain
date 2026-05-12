package options

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"syreen/config"
	"syreen/x/options/keeper"
	"syreen/x/options/types"
)

var (
	_ module.HasGenesis = AppModule{}
	_ module.AppModule  = AppModule{}
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
func (AppModule) GetTxCmd() *cobra.Command    { return nil }
func (AppModule) GetQueryCmd() *cobra.Command  { return nil }
func (AppModule) ConsensusVersion() uint64 { return 1 }

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), *am.keeper)
}

func (am AppModule) InitGenesis(ctx sdk.Context, _ codec.JSONCodec, data json.RawMessage) {
	var gs types.GenesisState
	if err := json.Unmarshal(data, &gs); err != nil {
		panic(fmt.Errorf("failed to unmarshal %s genesis: %w", types.ModuleName, err))
	}
	for _, opt := range gs.Options {
		am.keeper.SetOption(ctx, opt)
	}
	if gs.NextID > 0 {
		am.keeper.SetNextOptionID(ctx, gs.NextID)
	}
	return
}

func (am AppModule) ExportGenesis(ctx sdk.Context, _ codec.JSONCodec) json.RawMessage {
	gs := types.GenesisState{
		Options: am.keeper.GetAllOptions(ctx),
		NextID:  am.keeper.GetNextOptionID(ctx),
	}
	bz, _ := json.Marshal(gs)
	return bz
}

func (am AppModule) BeginBlock(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if !config.IsModuleEnabled("options") {
		return nil
	}
	am.keeper.BeginBlockHandler(ctx)
	return nil
}

func (am AppModule) EndBlock(_ context.Context) error { return nil }

func (am AppModule) IsOnePerModuleType() {}
func (am AppModule) IsAppModule()        {}
