package lending

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
	"syreen/x/lending/keeper"
	"syreen/x/lending/types"
)

var (
	_ module.HasGenesis  = AppModule{}
	_ module.AppModule   = AppModule{}
)

type AppModule struct {
	keeper *keeper.Keeper
}

func NewAppModule(k *keeper.Keeper) AppModule { return AppModule{keeper: k} }

func (AppModule) Name() string { return types.ModuleName }
func (AppModule) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) { types.RegisterLegacyAminoCodec(cdc) }
func (AppModule) RegisterInterfaces(registry cdctypes.InterfaceRegistry) { types.RegisterInterfaces(registry) }
func (AppModule) DefaultGenesis(_ codec.JSONCodec) json.RawMessage { bz, _ := json.Marshal(types.DefaultGenesis()); return bz }
func (AppModule) ValidateGenesis(_ codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	if err := json.Unmarshal(bz, &data); err != nil { panic(fmt.Errorf("failed to unmarshal %s genesis: %w", types.ModuleName, err)) }
	return data.Validate()
}
func (AppModule) RegisterGRPCGatewayRoutes(_ client.Context, _ *runtime.ServeMux) {}
func (AppModule) GetTxCmd() *cobra.Command    { return nil }
func (AppModule) GetQueryCmd() *cobra.Command  { return nil }
func (AppModule) ConsensusVersion() uint64 { return 1 }

func (am AppModule) RegisterServices(cfg module.Configurator) { types.RegisterMsgServer(cfg.MsgServer(), *am.keeper) }

func (am AppModule) InitGenesis(ctx sdk.Context, _ codec.JSONCodec, data json.RawMessage) {
	var gs types.GenesisState
	if err := json.Unmarshal(data, &gs); err != nil { panic(fmt.Errorf("failed to unmarshal %s genesis: %w", types.ModuleName, err)) }
	for _, p := range gs.Pools { am.keeper.SetPool(ctx, p) }
	for _, d := range gs.Deposits { am.keeper.SetDeposit(ctx, d) }
	for _, b := range gs.Borrows { am.keeper.SetBorrow(ctx, b) }
	if gs.NextPoolID > 0 {
		am.keeper.SetNextPoolID(ctx, gs.NextPoolID)
	}
	if gs.NextBorrowID > 0 {
		am.keeper.SetNextBorrowID(ctx, gs.NextBorrowID)
	}
	return
}

func (am AppModule) ExportGenesis(ctx sdk.Context, _ codec.JSONCodec) json.RawMessage {
	gs := types.GenesisState{
		Pools:        am.keeper.GetAllPools(ctx),
		Deposits:     am.keeper.GetAllDeposits(ctx),
		Borrows:      am.keeper.GetAllBorrows(ctx),
		NextPoolID:   am.keeper.GetNextPoolID(ctx),
		NextBorrowID: am.keeper.GetNextBorrowID(ctx),
	}
	bz, _ := json.Marshal(gs)
	return bz
}

func (am AppModule) BeginBlock(goCtx context.Context) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if !config.IsModuleEnabled("lending") {
		return nil
	}
	// C-1: sample collateral spot prices every block to feed the in-lending TWAP.
	am.keeper.SampleAllPoolPrices(ctx)
	// Accrue interest every 10 blocks to reduce overhead
	if ctx.BlockHeight()%10 == 0 {
		am.keeper.AccrueAllInterest(ctx)
	}
	return nil
}

func (am AppModule) EndBlock(_ context.Context) error { return nil }
func (am AppModule) IsOnePerModuleType() {}
func (am AppModule) IsAppModule()        {}
