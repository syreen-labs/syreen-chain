package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	syreenconfig "syreen/config"
	"syreen/x/compute/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) StoreCode(ctx context.Context, msg *types.MsgStoreCode) (*types.MsgStoreCodeResponse, error) {
	if !syreenconfig.IsModuleEnabled("compute") {
		return nil, syreenconfig.ErrModuleDisabled("compute")
	}
	codeID, err := m.Keeper.Create(ctx, msg.Sender, msg.WASMByteCode, msg.InstantiatePermission)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"store_code",
		sdk.NewAttribute("sender", msg.Sender),
		sdk.NewAttribute("code_id", fmt.Sprintf("%d", codeID)),
	))

	return &types.MsgStoreCodeResponse{CodeID: codeID}, nil
}

func (m msgServer) InstantiateContract(ctx context.Context, msg *types.MsgInstantiateContract) (*types.MsgInstantiateContractResponse, error) {
	if !syreenconfig.IsModuleEnabled("compute") {
		return nil, syreenconfig.ErrModuleDisabled("compute")
	}
	contractAddr, err := m.Keeper.Instantiate(ctx, msg.CodeID, msg.Sender, msg.Admin, msg.Label, msg.Msg, msg.Funds)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"instantiate_contract",
		sdk.NewAttribute("sender", msg.Sender),
		sdk.NewAttribute("code_id", fmt.Sprintf("%d", msg.CodeID)),
		sdk.NewAttribute("contract_address", contractAddr),
	))

	return &types.MsgInstantiateContractResponse{Address: contractAddr}, nil
}

func (m msgServer) ExecuteContract(ctx context.Context, msg *types.MsgExecuteContract) (*types.MsgExecuteContractResponse, error) {
	if !syreenconfig.IsModuleEnabled("compute") {
		return nil, syreenconfig.ErrModuleDisabled("compute")
	}
	data, err := m.Keeper.Execute(ctx, msg.Contract, msg.Sender, msg.Msg, msg.Funds)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"execute_contract",
		sdk.NewAttribute("sender", msg.Sender),
		sdk.NewAttribute("contract", msg.Contract),
	))

	return &types.MsgExecuteContractResponse{Data: data}, nil
}

func (m msgServer) MigrateContract(ctx context.Context, msg *types.MsgMigrateContract) (*types.MsgMigrateContractResponse, error) {
	if !syreenconfig.IsModuleEnabled("compute") {
		return nil, syreenconfig.ErrModuleDisabled("compute")
	}
	data, err := m.Keeper.Migrate(ctx, msg.Contract, msg.Sender, msg.CodeID, msg.Msg)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"migrate_contract",
		sdk.NewAttribute("sender", msg.Sender),
		sdk.NewAttribute("contract", msg.Contract),
		sdk.NewAttribute("code_id", fmt.Sprintf("%d", msg.CodeID)),
	))

	return &types.MsgMigrateContractResponse{Data: data}, nil
}

func (m msgServer) UpdateAdmin(ctx context.Context, msg *types.MsgUpdateAdmin) (*types.MsgUpdateAdminResponse, error) {
	if !syreenconfig.IsModuleEnabled("compute") {
		return nil, syreenconfig.ErrModuleDisabled("compute")
	}
	err := m.Keeper.UpdateAdmin(ctx, msg.Contract, msg.Sender, msg.NewAdmin)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"update_admin",
		sdk.NewAttribute("sender", msg.Sender),
		sdk.NewAttribute("contract", msg.Contract),
		sdk.NewAttribute("new_admin", msg.NewAdmin),
	))

	return &types.MsgUpdateAdminResponse{}, nil
}
