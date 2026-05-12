package portfolio

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

	"syreen/x/portfolio/keeper"
	"syreen/x/portfolio/types"
)

var (
	_ module.HasGenesis  = AppModule{}
	_ module.AppModule   = AppModule{}
	_ module.HasServices = AppModule{}
)

type AppModule struct {
	keeper *keeper.Keeper
}

func NewAppModule(k *keeper.Keeper) AppModule {
	return AppModule{keeper: k}
}

func (AppModule) Name() string                                                     { return types.ModuleName }
func (AppModule) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino)                 { types.RegisterLegacyAminoCodec(cdc) }
func (AppModule) RegisterInterfaces(registry cdctypes.InterfaceRegistry)           { types.RegisterInterfaces(registry) }
func (AppModule) DefaultGenesis(_ codec.JSONCodec) json.RawMessage                 { bz, _ := json.Marshal(types.DefaultGenesis()); return bz }
func (AppModule) ValidateGenesis(_ codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	if err := json.Unmarshal(bz, &data); err != nil { panic(fmt.Errorf("failed to unmarshal %s genesis: %w", types.ModuleName, err)) }
	return data.Validate()
}
func (AppModule) RegisterGRPCGatewayRoutes(_ client.Context, _ *runtime.ServeMux) {}
func (AppModule) GetTxCmd() *cobra.Command                                        { return nil }
func (AppModule) GetQueryCmd() *cobra.Command                                     { return nil }

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), *am.keeper)
}

func (am AppModule) InitGenesis(ctx sdk.Context, _ codec.JSONCodec, data json.RawMessage) {
	var gs types.GenesisState
	if err := json.Unmarshal(data, &gs); err != nil {
		ctx.Logger().Error("FATAL: malformed portfolio genesis JSON — halting startup", "module", types.ModuleName, "error", err)
		panic(fmt.Errorf("failed to unmarshal %s genesis: %w", types.ModuleName, err))
	}
	for _, p := range gs.Portfolios { am.keeper.SetPortfolio(ctx, p) }
	for _, a := range gs.Assets { am.keeper.SetPortfolioAsset(ctx, a) }
	for _, t := range gs.Trades { am.keeper.SetTradeRecord(ctx, t) }
	for _, act := range gs.Activities { am.keeper.SetActivity(ctx, act) }
	for _, c := range gs.Competitions { am.keeper.SetCompetition(ctx, c) }
	for _, e := range gs.Entries { am.keeper.SetCompetitionEntry(ctx, e) }
	return
}

func (am AppModule) ExportGenesis(ctx sdk.Context, _ codec.JSONCodec) json.RawMessage {
	gs := types.GenesisState{
		Portfolios:   am.keeper.GetAllPortfolios(ctx),
		Assets:       am.keeper.GetAllAssets(ctx),
		Trades:       am.keeper.GetAllTrades(ctx),
		Activities:   am.keeper.GetAllActivities(ctx),
		Competitions: am.keeper.GetAllCompetitions(ctx),
		Entries:      am.keeper.GetAllEntries(ctx),
	}
	bz, _ := json.Marshal(gs)
	return bz
}

func (AppModule) ConsensusVersion() uint64 { return 1 }

func (am AppModule) BeginBlock(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	am.keeper.ProcessBeginBlock(ctx)
	return nil
}

func (am AppModule) EndBlock(_ context.Context) error { return nil }
func (am AppModule) IsOnePerModuleType()              {}
func (am AppModule) IsAppModule()                     {}
