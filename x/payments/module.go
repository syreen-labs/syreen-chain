package payments

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

	"syreen/x/payments/client/cli"
	"syreen/x/payments/keeper"
	"syreen/x/payments/types"
)

var (
	_ module.HasGenesis  = AppModule{}
	_ module.AppModule   = AppModule{}
)

type AppModule struct {
	keeper *keeper.Keeper
}

func NewAppModule(keeper *keeper.Keeper) AppModule {
	return AppModule{keeper: keeper}
}

func (AppModule) Name() string                                                           { return types.ModuleName }
func (AppModule) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino)                        { types.RegisterLegacyAminoCodec(cdc) }
func (AppModule) RegisterInterfaces(registry cdctypes.InterfaceRegistry)                  { types.RegisterInterfaces(registry) }
func (AppModule) DefaultGenesis(_ codec.JSONCodec) json.RawMessage                        { bz, _ := json.Marshal(types.DefaultGenesis()); return bz }
func (AppModule) RegisterGRPCGatewayRoutes(_ client.Context, _ *runtime.ServeMux)         {}
func (AppModule) GetTxCmd() *cobra.Command                                                { return cli.GetTxCmd() }
func (AppModule) GetQueryCmd() *cobra.Command                                             { return cli.GetQueryCmd() }
func (AppModule) ConsensusVersion() uint64                                                { return 1 }

func (AppModule) ValidateGenesis(_ codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	if err := json.Unmarshal(bz, &data); err != nil {
		panic(fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err))
	}
	return data.Validate()
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(*am.keeper))
}

func (am AppModule) InitGenesis(ctx sdk.Context, _ codec.JSONCodec, data json.RawMessage) {
	var gs types.GenesisState
	if err := json.Unmarshal(data, &gs); err != nil {
		panic(fmt.Errorf("failed to unmarshal %s genesis: %w", types.ModuleName, err))
	}
	am.keeper.InitGenesis(ctx, gs)
	return
}

func (am AppModule) ExportGenesis(ctx sdk.Context, _ codec.JSONCodec) json.RawMessage {
	gs := am.keeper.ExportGenesis(ctx)
	bz, _ := json.Marshal(gs)
	return bz
}

func (am AppModule) BeginBlock(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	am.keeper.ExpireInvoices(ctx)
	return nil
}

func (am AppModule) IsOnePerModuleType() {}
func (am AppModule) IsAppModule()        {}
