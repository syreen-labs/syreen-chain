package mevprotection

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

	"syreen/x/mevprotection/client/cli"
	"syreen/x/mevprotection/keeper"
	"syreen/x/mevprotection/types"
)

var (
	_ module.HasGenesis = AppModule{}
	_ module.AppModule  = AppModule{}
)

// AppModule implements the AppModule interface
type AppModule struct {
	keeper keeper.Keeper
}

func NewAppModule(keeper keeper.Keeper) AppModule {
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
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(am.keeper))
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

// BeginBlock prunes expired committed transactions and runs MEV detection.
// The primary MEV enforcement is in ProcessProposal (which rejects bad blocks).
// BeginBlock provides a secondary check: if a block was somehow committed with
// ordering violations, DetectMEV records penalties and slashes the proposer.
func (am AppModule) BeginBlock(ctx context.Context) error {
	am.keeper.PruneExpiredCommits(ctx)
	am.keeper.BeginBlockMEVCheck(ctx)
	return nil
}

// EndBlock distributes accumulated MEV rewards to LPs and stakers every 100 blocks.
func (am AppModule) EndBlock(ctx context.Context) error {
	am.keeper.DistributeMEVRewards(ctx)
	return nil
}

func (am AppModule) IsOnePerModuleType() {}
func (am AppModule) IsAppModule()        {}
