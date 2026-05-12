package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	syreenconfig "syreen/config"
	"syreen/x/abstractaccount/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) CreateSmartAccount(ctx context.Context, msg *types.MsgCreateSmartAccount) (*types.MsgCreateSmartAccountResponse, error) {
	if !syreenconfig.IsModuleEnabled("abstractaccount") {
		return nil, syreenconfig.ErrModuleDisabled("abstractaccount")
	}
	addr, err := m.Keeper.CreateSmartAccount(ctx, msg)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"create_smart_account",
		sdk.NewAttribute("sender", msg.Sender),
		sdk.NewAttribute("address", addr),
		sdk.NewAttribute("account_type", string(msg.AccountType)),
	))

	return &types.MsgCreateSmartAccountResponse{Address: addr}, nil
}

func (m msgServer) CreateSessionKey(ctx context.Context, msg *types.MsgCreateSessionKey) (*types.MsgCreateSessionKeyResponse, error) {
	if !syreenconfig.IsModuleEnabled("abstractaccount") {
		return nil, syreenconfig.ErrModuleDisabled("abstractaccount")
	}
	err := m.Keeper.CreateSessionKey(ctx, msg.Granter, msg.Grantee, msg.Permissions, msg.Duration)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"create_session_key",
		sdk.NewAttribute("granter", msg.Granter),
		sdk.NewAttribute("grantee", msg.Grantee),
	))

	return &types.MsgCreateSessionKeyResponse{}, nil
}

func (m msgServer) RevokeSessionKey(ctx context.Context, msg *types.MsgRevokeSessionKey) (*types.MsgRevokeSessionKeyResponse, error) {
	if !syreenconfig.IsModuleEnabled("abstractaccount") {
		return nil, syreenconfig.ErrModuleDisabled("abstractaccount")
	}
	err := m.Keeper.RevokeSessionKey(ctx, msg.Granter, msg.SessionKeyAddr)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"revoke_session_key",
		sdk.NewAttribute("granter", msg.Granter),
		sdk.NewAttribute("session_key_addr", msg.SessionKeyAddr),
	))

	return &types.MsgRevokeSessionKeyResponse{}, nil
}

func (m msgServer) InitiateRecovery(ctx context.Context, msg *types.MsgInitiateRecovery) (*types.MsgInitiateRecoveryResponse, error) {
	if !syreenconfig.IsModuleEnabled("abstractaccount") {
		return nil, syreenconfig.ErrModuleDisabled("abstractaccount")
	}
	err := m.Keeper.InitiateRecovery(ctx, msg.Guardian, msg.Account, msg.NewOwners)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"initiate_recovery",
		sdk.NewAttribute("guardian", msg.Guardian),
		sdk.NewAttribute("account", msg.Account),
	))

	return &types.MsgInitiateRecoveryResponse{}, nil
}

func (m msgServer) ApproveRecovery(ctx context.Context, msg *types.MsgApproveRecovery) (*types.MsgApproveRecoveryResponse, error) {
	if !syreenconfig.IsModuleEnabled("abstractaccount") {
		return nil, syreenconfig.ErrModuleDisabled("abstractaccount")
	}
	err := m.Keeper.ApproveRecovery(ctx, msg.Guardian, msg.Account)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"approve_recovery",
		sdk.NewAttribute("guardian", msg.Guardian),
		sdk.NewAttribute("account", msg.Account),
	))

	return &types.MsgApproveRecoveryResponse{}, nil
}

func (m msgServer) ExecuteRecovery(ctx context.Context, msg *types.MsgExecuteRecovery) (*types.MsgExecuteRecoveryResponse, error) {
	if !syreenconfig.IsModuleEnabled("abstractaccount") {
		return nil, syreenconfig.ErrModuleDisabled("abstractaccount")
	}
	err := m.Keeper.ExecuteRecovery(ctx, msg.Sender, msg.Account)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"execute_recovery",
		sdk.NewAttribute("sender", msg.Sender),
		sdk.NewAttribute("account", msg.Account),
	))

	return &types.MsgExecuteRecoveryResponse{}, nil
}

func (m msgServer) SponsorGas(ctx context.Context, msg *types.MsgSponsorGas) (*types.MsgSponsorGasResponse, error) {
	if !syreenconfig.IsModuleEnabled("abstractaccount") {
		return nil, syreenconfig.ErrModuleDisabled("abstractaccount")
	}
	err := m.Keeper.SponsorGas(ctx, msg.Sponsor, msg.Sponsored, msg.GasLimit, msg.Duration)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"sponsor_gas",
		sdk.NewAttribute("sponsor", msg.Sponsor),
		sdk.NewAttribute("sponsored", msg.Sponsored),
	))

	return &types.MsgSponsorGasResponse{}, nil
}

func (m msgServer) BatchExecute(ctx context.Context, msg *types.MsgBatchExecute) (*types.MsgBatchExecuteResponse, error) {
	if !syreenconfig.IsModuleEnabled("abstractaccount") {
		return nil, syreenconfig.ErrModuleDisabled("abstractaccount")
	}
	_, err := m.Keeper.ExecuteBatch(ctx, msg.Sender, msg.Messages)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"batch_execute",
		sdk.NewAttribute("sender", msg.Sender),
	))

	return &types.MsgBatchExecuteResponse{}, nil
}
