package keeper

import (
	"context"

	syreenconfig "syreen/config"

	"cosmossdk.io/math"
	"syreen/x/predict/types"
)

var _ types.MsgServer = Keeper{}

func (k Keeper) CreateMarket(ctx context.Context, msg *types.MsgCreateMarket) (*types.MsgCreateMarketResponse, error) {
	if !syreenconfig.IsModuleEnabled("predict") {
		return nil, syreenconfig.ErrModuleDisabled("predict")
	}
	marketID, err := k.ExecuteCreateMarket(ctx, msg.Creator, msg.Question, msg.Resolver, msg.QuoteDenom, msg.ResolutionBlock, msg.InitialLiquidity)
	if err != nil {
		return nil, err
	}
	return &types.MsgCreateMarketResponse{MarketID: marketID}, nil
}

func (k Keeper) BuyShares(ctx context.Context, msg *types.MsgBuyShares) (*types.MsgBuySharesResponse, error) {
	if !syreenconfig.IsModuleEnabled("predict") {
		return nil, syreenconfig.ErrModuleDisabled("predict")
	}
	sharesBought, err := k.ExecuteBuyShares(ctx, msg.Sender, msg.MarketID, msg.Outcome, msg.Amount)
	if err != nil {
		return nil, err
	}
	avgPrice := math.ZeroInt()
	if !sharesBought.IsZero() {
		avgPrice = msg.Amount.Mul(math.NewInt(10000)).Quo(sharesBought)
	}
	return &types.MsgBuySharesResponse{SharesBought: sharesBought, AvgPrice: avgPrice}, nil
}

func (k Keeper) SellShares(ctx context.Context, msg *types.MsgSellShares) (*types.MsgSellSharesResponse, error) {
	if !syreenconfig.IsModuleEnabled("predict") {
		return nil, syreenconfig.ErrModuleDisabled("predict")
	}
	quoteReturned, err := k.ExecuteSellShares(ctx, msg.Sender, msg.MarketID, msg.Outcome, msg.Shares)
	if err != nil {
		return nil, err
	}
	return &types.MsgSellSharesResponse{QuoteReturned: quoteReturned}, nil
}

func (k Keeper) ResolveMarket(ctx context.Context, msg *types.MsgResolveMarket) (*types.MsgResolveMarketResponse, error) {
	if !syreenconfig.IsModuleEnabled("predict") {
		return nil, syreenconfig.ErrModuleDisabled("predict")
	}
	if err := k.ExecuteResolveMarket(ctx, msg.Resolver, msg.MarketID, msg.Outcome); err != nil {
		return nil, err
	}
	return &types.MsgResolveMarketResponse{}, nil
}

func (k Keeper) ClaimWinnings(ctx context.Context, msg *types.MsgClaimWinnings) (*types.MsgClaimWinningsResponse, error) {
	if !syreenconfig.IsModuleEnabled("predict") {
		return nil, syreenconfig.ErrModuleDisabled("predict")
	}
	amount, err := k.ExecuteClaimWinnings(ctx, msg.Sender, msg.MarketID)
	if err != nil {
		return nil, err
	}
	return &types.MsgClaimWinningsResponse{Amount: amount}, nil
}
