package dex

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

	"syreen/x/dex/client/cli"
	"syreen/x/dex/keeper"
	"syreen/x/dex/types"
)

var (
	_ module.HasGenesis = AppModule{}
	_ module.AppModule  = AppModule{}
)

// AppModule implements the AppModule interface.
type AppModule struct {
	keeper *keeper.Keeper
}

func NewAppModule(k *keeper.Keeper) AppModule {
	return AppModule{
		keeper: k,
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
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(*am.keeper))
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

// BeginBlock runs DEX operations only when pools exist.
func (am AppModule) BeginBlock(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	pools := am.keeper.GetAllPools(ctx)
	if len(pools) == 0 {
		return nil
	}

	am.keeper.RunRiskChecks(ctx)
	am.keeper.ProcessConditionalOrders(ctx)
	am.keeper.RunOrderBookMatching(ctx)
	am.keeper.UpdateAllPoolRiskScores(ctx)

	if ctx.BlockHeight()%5 == 0 {
		am.keeper.UpdateAllSignals(ctx)
		am.keeper.UpdateSentiment(ctx)
		am.keeper.UpdateAllOracleSignals(ctx)
	}

	if ctx.BlockHeight()%1000 == 0 {
		am.keeper.PruneExpiredVolumes(ctx)
	}
	return nil
}

// EndBlock processes pending trades only when pools exist.
func (am AppModule) EndBlock(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	pools := am.keeper.GetAllPools(ctx)
	if len(pools) == 0 {
		return nil
	}

	am.keeper.ProcessCopyTrades(ctx)
	am.keeper.FinalizeCandles(ctx)
	am.keeper.ProcessIBCReturns(ctx)
	return nil
}

func (am AppModule) IsOnePerModuleType() {}
func (am AppModule) IsAppModule()        {}
