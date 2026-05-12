package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	syreenconfig "syreen/config"
	"syreen/x/evm/types"
)

// MsgServer implements the EVM message handling interface.
// H3: Provides a proper message server so MsgEthereumTx can be delivered via Cosmos SDK tx routing.
type MsgServer struct {
	keeper Keeper
}

// NewMsgServer returns a new MsgServer instance
func NewMsgServer(k Keeper) MsgServer {
	return MsgServer{keeper: k}
}

// EthereumTx handles a MsgEthereumTx by executing it against the EVM.
func (ms MsgServer) EthereumTx(goCtx context.Context, msg *types.MsgEthereumTx) (*types.MsgEthereumTxResponse, error) {
	if !syreenconfig.IsModuleEnabled("evm") {
		return nil, syreenconfig.ErrModuleDisabled("evm")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid MsgEthereumTx: %w", err)
	}

	result, err := ms.keeper.ExecuteEVMTx(ctx, msg)
	if err != nil {
		return nil, err
	}

	resp := &types.MsgEthereumTxResponse{
		GasUsed: result.GasUsed,
		VmError: result.VmError,
	}
	if result.ReturnData != nil {
		resp.ReturnData = fmt.Sprintf("%x", result.ReturnData)
	}
	if result.ContractAddr != (types.EmptyEthAddress) {
		resp.ContractAddress = result.ContractAddr.Hex()
	}

	return resp, nil
}
