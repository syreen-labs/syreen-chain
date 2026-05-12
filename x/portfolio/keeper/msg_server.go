package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/portfolio/types"
)

var _ types.MsgServer = Keeper{}

// allowedTradeRecorderModules is the set of module names whose module
// addresses are permitted to submit MsgRecordTrade. RecordTradeAndUpdatePortfolio
// affects competition leaderboards, so accepting it from arbitrary user
// senders would let anyone spoof their PnL. The dex and intent modules are
// the only legitimate origin points for executed trades on chain.
var allowedTradeRecorderModules = []string{"dex", "intent"}

func (k Keeper) RecordTrade(ctx context.Context, msg *types.MsgRecordTrade) (*types.MsgRecordTradeResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// L-3 fix: only allow the dex and intent module accounts (or chain authority)
	// to record trades — otherwise any user could spoof the leaderboard.
	if !k.isAllowedTradeRecorder(msg.Sender) {
		return nil, types.ErrUnauthorized
	}

	tradeID, err := k.RecordTradeAndUpdatePortfolio(ctx, msg.Sender, msg.TradeType, msg.Denom, msg.Amount, msg.Price, msg.PnL, sdkCtx.BlockHeight())
	if err != nil {
		return nil, err
	}
	return &types.MsgRecordTradeResponse{TradeID: tradeID}, nil
}

// isAllowedTradeRecorder reports whether sender matches the chain authority
// or one of the allowlisted module account addresses (dex, intent).
func (k Keeper) isAllowedTradeRecorder(sender string) bool {
	if sender == "" {
		return false
	}
	if sender == k.authority {
		return true
	}
	for _, moduleName := range allowedTradeRecorderModules {
		modAddr := k.accountKeeper.GetModuleAddress(moduleName)
		if modAddr != nil && modAddr.String() == sender {
			return true
		}
	}
	return false
}

func (k Keeper) CreateCompetition(ctx context.Context, msg *types.MsgCreateCompetition) (*types.MsgCreateCompetitionResponse, error) {
	compID, err := k.ExecuteCreateCompetition(ctx, msg.Creator, msg.Name, msg.StartBlock, msg.EndBlock, msg.PrizeDenom, msg.PrizePool, msg.EntryFee, msg.MaxParticipants)
	if err != nil {
		return nil, err
	}
	return &types.MsgCreateCompetitionResponse{CompetitionID: compID}, nil
}

func (k Keeper) JoinCompetition(ctx context.Context, msg *types.MsgJoinCompetition) (*types.MsgJoinCompetitionResponse, error) {
	if err := k.ExecuteJoinCompetition(ctx, msg.Sender, msg.CompetitionID); err != nil {
		return nil, err
	}
	return &types.MsgJoinCompetitionResponse{}, nil
}

func (k Keeper) EndCompetition(ctx context.Context, msg *types.MsgEndCompetition) (*types.MsgEndCompetitionResponse, error) {
	// Only the module authority can end competitions manually
	if msg.Authority != k.authority {
		return nil, fmt.Errorf("unauthorized: only the module authority can end competitions")
	}
	winners, err := k.EndCompetitionAndDistribute(ctx, msg.CompetitionID)
	if err != nil {
		return nil, err
	}
	return &types.MsgEndCompetitionResponse{Winners: winners}, nil
}

func (k Keeper) UpdatePortfolio(ctx context.Context, msg *types.MsgUpdatePortfolio) (*types.MsgUpdatePortfolioResponse, error) {
	totalValue := k.RecalculatePortfolio(ctx, msg.Sender)
	return &types.MsgUpdatePortfolioResponse{TotalValue: totalValue}, nil
}
