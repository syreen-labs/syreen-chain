package keeper_test

import (
	"encoding/json"
	"testing"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"syreen/x/intent/keeper"
	"syreen/x/intent/types"
)

// ---------------------------------------------------------------------------
// Helper
// ---------------------------------------------------------------------------

func makeStrategyMsg(template string, risk types.RiskLevel) *types.MsgCreateStrategy {
	msg := &types.MsgCreateStrategy{
		Creator:      creatorAddr,
		TemplateName: template,
		PoolID:       1,
		InputDenom:   "usyreen",
		OutputDenom:  "uusdc",
		TotalBudget:  math.NewInt(10000000), // 10 SYR
		RiskLevel:    risk,
		ExpiryBlocks: 1000,
	}
	// Template-specific defaults
	if template == types.StrategySafeAccumulate || template == types.StrategySmartDCA {
		msg.NumExecutions = 5
		msg.IntervalBlocks = 10
	}
	if template == types.StrategyGridTrading {
		msg.GridLevels = 3
	}
	return msg
}

// ---------------------------------------------------------------------------
// Template listing
// ---------------------------------------------------------------------------

func TestGetAvailableTemplates(t *testing.T) {
	templates := types.GetAvailableTemplates()
	require.Len(t, templates, 12)

	names := make(map[string]bool)
	for _, tmpl := range templates {
		names[tmpl.Name] = true
		require.NotEmpty(t, tmpl.Description)
		require.True(t, tmpl.MinSteps >= 2)
	}

	require.True(t, names[types.StrategySafeAccumulate])
	require.True(t, names[types.StrategyMomentumRide])
	require.True(t, names[types.StrategyGridTrading])
	require.True(t, names[types.StrategySwingTrade])
	require.True(t, names[types.StrategyProtectiveSell])
	require.True(t, names[types.StrategyScalp])

	// Tier 1 AI-reactive templates
	require.True(t, names[types.StrategySentimentContrarian])
	require.True(t, names[types.StrategySmartDCA])
	require.True(t, names[types.StrategySignalComposite])
	require.True(t, names[types.StrategyRSIMeanReversion])
	require.True(t, names[types.StrategyVolatilityBreakout])
	require.True(t, names[types.StrategyTrendConfluence])
}

// ---------------------------------------------------------------------------
// Risk params
// ---------------------------------------------------------------------------

func TestGetRiskParams(t *testing.T) {
	sl, tp, gs := types.GetRiskParams(types.RiskConservative)
	require.Equal(t, "0.050000000000000000", sl.String())
	require.Equal(t, "0.100000000000000000", tp.String())
	require.Equal(t, "0.020000000000000000", gs.String())

	sl, tp, gs = types.GetRiskParams(types.RiskModerate)
	require.Equal(t, "0.100000000000000000", sl.String())
	require.Equal(t, "0.200000000000000000", tp.String())
	require.Equal(t, "0.030000000000000000", gs.String())

	sl, tp, gs = types.GetRiskParams(types.RiskAggressive)
	require.Equal(t, "0.150000000000000000", sl.String())
	require.Equal(t, "0.300000000000000000", tp.String())
	require.Equal(t, "0.050000000000000000", gs.String())
}

// ---------------------------------------------------------------------------
// Validation
// ---------------------------------------------------------------------------

func TestMsgCreateStrategy_ValidateBasic(t *testing.T) {
	// Valid
	msg := makeStrategyMsg(types.StrategySafeAccumulate, types.RiskModerate)
	require.NoError(t, msg.ValidateBasic())

	// Invalid creator
	bad := *msg
	bad.Creator = "bad"
	require.Error(t, bad.ValidateBasic())

	// Invalid template
	bad2 := *msg
	bad2.TemplateName = "nonexistent"
	require.Error(t, bad2.ValidateBasic())

	// Invalid risk
	bad3 := *msg
	bad3.RiskLevel = "yolo"
	require.Error(t, bad3.ValidateBasic())

	// Zero pool
	bad4 := *msg
	bad4.PoolID = 0
	require.Error(t, bad4.ValidateBasic())

	// Zero budget
	bad5 := *msg
	bad5.TotalBudget = math.ZeroInt()
	require.Error(t, bad5.ValidateBasic())

	// Safe accumulate without num_executions
	bad6 := *msg
	bad6.NumExecutions = 0
	require.Error(t, bad6.ValidateBasic())

	// Grid trading needs grid_levels >= 2
	gridMsg := makeStrategyMsg(types.StrategyGridTrading, types.RiskModerate)
	gridMsg.GridLevels = 1
	require.Error(t, gridMsg.ValidateBasic())

	// Grid trading max 5
	gridMsg2 := makeStrategyMsg(types.StrategyGridTrading, types.RiskModerate)
	gridMsg2.GridLevels = 6
	require.Error(t, gridMsg2.ValidateBasic())
}

func TestMsgCancelStrategy_ValidateBasic(t *testing.T) {
	msg := &types.MsgCancelStrategy{
		Creator:    creatorAddr,
		StrategyID: "strategy-1",
	}
	require.NoError(t, msg.ValidateBasic())

	bad := *msg
	bad.StrategyID = ""
	require.Error(t, bad.ValidateBasic())

	bad2 := *msg
	bad2.Creator = "bad"
	require.Error(t, bad2.ValidateBasic())
}

// ---------------------------------------------------------------------------
// Strategy creation — safe_accumulate
// ---------------------------------------------------------------------------

func TestCreateStrategy_SafeAccumulate(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	// Set spot price: 1 usyreen = 2 uusdc (or 2 usyreen per 1 uusdc)
	dex.SetSpotPrice(1, "usyreen", "uusdc", math.LegacyMustNewDecFromStr("2.0"))

	msg := makeStrategyMsg(types.StrategySafeAccumulate, types.RiskModerate)
	strategyID, chainID, err := k.CreateStrategy(ctx, msg)
	require.NoError(t, err)
	require.NotEmpty(t, strategyID)
	require.NotEmpty(t, chainID)

	// Verify strategy stored
	strategy, found := k.GetStrategy(ctx, strategyID)
	require.True(t, found)
	require.Equal(t, types.StrategyStatusActive, strategy.Status)
	require.Equal(t, types.StrategySafeAccumulate, strategy.TemplateName)
	require.Equal(t, chainID, strategy.ChainID)
	require.Equal(t, creatorAddr, strategy.Creator)

	// Verify underlying chain
	chain, found := k.GetChain(ctx, chainID)
	require.True(t, found)
	require.Equal(t, types.ChainStatusActive, chain.Status)
	require.Equal(t, 2, len(chain.Steps))

	// Step 0: DCA
	require.Equal(t, types.IntentTypeDCA, chain.Steps[0].IntentType)
	var dcaBody types.DCAIntent
	require.NoError(t, json.Unmarshal(chain.Steps[0].Body, &dcaBody))
	require.Equal(t, "usyreen", dcaBody.InputDenom)
	require.Equal(t, "uusdc", dcaBody.OutputDenom)
	require.Equal(t, math.NewInt(10000000), dcaBody.TotalAmount)
	require.Equal(t, uint64(5), dcaBody.NumExecutions)
	require.Equal(t, uint64(10), dcaBody.IntervalBlocks)

	// Step 1: stop-loss
	require.Equal(t, types.IntentTypeStopLoss, chain.Steps[1].IntentType)
	var slBody types.StopLossIntent
	require.NoError(t, json.Unmarshal(chain.Steps[1].Body, &slBody))
	require.Equal(t, "uusdc", slBody.InputDenom)
	require.Equal(t, "usyreen", slBody.OutputDenom)
	// Stop price should be 0.5 * (1 - 0.10) = 0.45
	require.Equal(t, "0.450000000000000000", slBody.StopPrice.String())
}

// ---------------------------------------------------------------------------
// Strategy creation — momentum_ride
// ---------------------------------------------------------------------------

func TestCreateStrategy_MomentumRide(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	dex.SetSpotPrice(1, "usyreen", "uusdc", math.LegacyMustNewDecFromStr("2.0"))

	msg := makeStrategyMsg(types.StrategyMomentumRide, types.RiskAggressive)
	strategyID, chainID, err := k.CreateStrategy(ctx, msg)
	require.NoError(t, err)
	require.NotEmpty(t, strategyID)

	chain, found := k.GetChain(ctx, chainID)
	require.True(t, found)
	require.Equal(t, 3, len(chain.Steps))

	// Step 0: limit buy
	require.Equal(t, types.IntentTypeLimitBuy, chain.Steps[0].IntentType)
	var buyBody types.LimitBuyIntent
	require.NoError(t, json.Unmarshal(chain.Steps[0].Body, &buyBody))
	require.Equal(t, "usyreen", buyBody.InputDenom)
	require.Equal(t, "uusdc", buyBody.OutputDenom)
	// Buy price = 2.0 * (1 - 0.15) = 1.7
	require.Equal(t, "1.700000000000000000", buyBody.TargetPrice.String())

	// Step 1: take-profit (50% portion)
	require.Equal(t, types.IntentTypeTakeProfit, chain.Steps[1].IntentType)
	require.Equal(t, types.InputFromPortion, chain.Steps[1].InputSource.Type)
	require.Equal(t, "0.500000000000000000", chain.Steps[1].InputSource.PortionPercent.String())

	// Step 2: stop-loss (previous = remaining)
	require.Equal(t, types.IntentTypeStopLoss, chain.Steps[2].IntentType)
	require.Equal(t, types.InputFromPrevious, chain.Steps[2].InputSource.Type)
}

// ---------------------------------------------------------------------------
// Strategy creation — grid_trading
// ---------------------------------------------------------------------------

func TestCreateStrategy_GridTrading(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	dex.SetSpotPrice(1, "usyreen", "uusdc", math.LegacyMustNewDecFromStr("1.0"))

	msg := makeStrategyMsg(types.StrategyGridTrading, types.RiskModerate)
	msg.GridLevels = 3
	strategyID, chainID, err := k.CreateStrategy(ctx, msg)
	require.NoError(t, err)
	require.NotEmpty(t, strategyID)

	chain, found := k.GetChain(ctx, chainID)
	require.True(t, found)
	// 3 levels * 2 steps = 6 steps
	require.Equal(t, 6, len(chain.Steps))

	// Alternating buy/sell
	for i := 0; i < len(chain.Steps); i += 2 {
		require.Equal(t, types.IntentTypeLimitBuy, chain.Steps[i].IntentType)
		require.Equal(t, types.IntentTypeLimitSell, chain.Steps[i+1].IntentType)
		require.Equal(t, types.InputFromPrevious, chain.Steps[i+1].InputSource.Type)
	}

	// Budget split: 10M / 3 = 3333333 per level
	var buyBody types.LimitBuyIntent
	require.NoError(t, json.Unmarshal(chain.Steps[0].Body, &buyBody))
	require.Equal(t, math.NewInt(3333333), buyBody.InputAmount)
}

// ---------------------------------------------------------------------------
// Strategy creation — swing_trade
// ---------------------------------------------------------------------------

func TestCreateStrategy_SwingTrade(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	dex.SetSpotPrice(1, "usyreen", "uusdc", math.LegacyMustNewDecFromStr("1.0"))

	msg := makeStrategyMsg(types.StrategySwingTrade, types.RiskConservative)
	strategyID, chainID, err := k.CreateStrategy(ctx, msg)
	require.NoError(t, err)
	require.NotEmpty(t, strategyID)

	chain, found := k.GetChain(ctx, chainID)
	require.True(t, found)
	require.Equal(t, 2, len(chain.Steps))

	// Step 0: limit buy at 5% discount (conservative)
	require.Equal(t, types.IntentTypeLimitBuy, chain.Steps[0].IntentType)
	var buyBody types.LimitBuyIntent
	require.NoError(t, json.Unmarshal(chain.Steps[0].Body, &buyBody))
	require.Equal(t, "0.950000000000000000", buyBody.TargetPrice.String())

	// Step 1: limit sell at 10% premium
	require.Equal(t, types.IntentTypeLimitSell, chain.Steps[1].IntentType)
	var sellBody types.LimitSellIntent
	require.NoError(t, json.Unmarshal(chain.Steps[1].Body, &sellBody))
	// 1/1.0 * 1.1 = 1.1
	require.Equal(t, "1.100000000000000000", sellBody.TargetPrice.String())
}

// ---------------------------------------------------------------------------
// Strategy creation — protective_sell
// ---------------------------------------------------------------------------

func TestCreateStrategy_ProtectiveSell(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	dex.SetSpotPrice(1, "usyreen", "uusdc", math.LegacyMustNewDecFromStr("2.0"))

	msg := makeStrategyMsg(types.StrategyProtectiveSell, types.RiskModerate)
	strategyID, chainID, err := k.CreateStrategy(ctx, msg)
	require.NoError(t, err)
	require.NotEmpty(t, strategyID)

	chain, found := k.GetChain(ctx, chainID)
	require.True(t, found)
	require.Equal(t, 2, len(chain.Steps))

	// Step 0: take-profit (half budget)
	require.Equal(t, types.IntentTypeTakeProfit, chain.Steps[0].IntentType)
	var tpBody types.TakeProfitIntent
	require.NoError(t, json.Unmarshal(chain.Steps[0].Body, &tpBody))
	require.Equal(t, "usyreen", tpBody.InputDenom)
	require.Equal(t, math.NewInt(5000000), tpBody.InputAmount) // half of 10M

	// Step 1: stop-loss (other half)
	require.Equal(t, types.IntentTypeStopLoss, chain.Steps[1].IntentType)
	var slBody types.StopLossIntent
	require.NoError(t, json.Unmarshal(chain.Steps[1].Body, &slBody))
	require.Equal(t, "usyreen", slBody.InputDenom)
	require.Equal(t, math.NewInt(5000000), slBody.InputAmount) // remaining half
}

// ---------------------------------------------------------------------------
// Strategy creation — scalp
// ---------------------------------------------------------------------------

func TestCreateStrategy_Scalp(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	dex.SetSpotPrice(1, "usyreen", "uusdc", math.LegacyMustNewDecFromStr("1.0"))

	msg := makeStrategyMsg(types.StrategyScalp, types.RiskModerate)
	strategyID, chainID, err := k.CreateStrategy(ctx, msg)
	require.NoError(t, err)
	require.NotEmpty(t, strategyID)

	chain, found := k.GetChain(ctx, chainID)
	require.True(t, found)
	require.Equal(t, 2, len(chain.Steps))

	// Step 0: limit buy at tiny discount (half of 10% = 5%)
	require.Equal(t, types.IntentTypeLimitBuy, chain.Steps[0].IntentType)
	var buyBody types.LimitBuyIntent
	require.NoError(t, json.Unmarshal(chain.Steps[0].Body, &buyBody))
	// 1.0 * (1 - 0.05) = 0.95
	require.Equal(t, "0.950000000000000000", buyBody.TargetPrice.String())

	// Step 1: take-profit at tiny premium
	require.Equal(t, types.IntentTypeTakeProfit, chain.Steps[1].IntentType)
	var tpBody types.TakeProfitIntent
	require.NoError(t, json.Unmarshal(chain.Steps[1].Body, &tpBody))
	// 1/1.0 * (1 + 0.05) = 1.05
	require.Equal(t, "1.050000000000000000", tpBody.TargetPrice.String())
}

// ---------------------------------------------------------------------------
// Strategy storage and retrieval
// ---------------------------------------------------------------------------

func TestStrategyStorage(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	dex.SetSpotPrice(1, "usyreen", "uusdc", math.LegacyMustNewDecFromStr("1.0"))

	// Create two strategies
	msg1 := makeStrategyMsg(types.StrategySwingTrade, types.RiskModerate)
	sid1, _, err := k.CreateStrategy(ctx, msg1)
	require.NoError(t, err)

	msg2 := makeStrategyMsg(types.StrategyScalp, types.RiskAggressive)
	sid2, _, err := k.CreateStrategy(ctx, msg2)
	require.NoError(t, err)

	// Get by creator
	strategies := k.GetStrategiesByCreator(ctx, creatorAddr)
	require.Len(t, strategies, 2)

	// IDs are sequential
	require.Equal(t, sid1, strategies[0].ID)
	require.Equal(t, sid2, strategies[1].ID)

	// Not found for different creator
	strategies2 := k.GetStrategiesByCreator(ctx, solverAddr)
	require.Len(t, strategies2, 0)
}

// ---------------------------------------------------------------------------
// Strategy cancellation
// ---------------------------------------------------------------------------

func TestCancelStrategy(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	dex.SetSpotPrice(1, "usyreen", "uusdc", math.LegacyMustNewDecFromStr("1.0"))

	msg := makeStrategyMsg(types.StrategySwingTrade, types.RiskModerate)
	strategyID, chainID, err := k.CreateStrategy(ctx, msg)
	require.NoError(t, err)

	// Cancel
	err = k.CancelStrategy(ctx, &types.MsgCancelStrategy{
		Creator:    creatorAddr,
		StrategyID: strategyID,
	})
	require.NoError(t, err)

	// Verify strategy cancelled
	strategy, found := k.GetStrategy(ctx, strategyID)
	require.True(t, found)
	require.Equal(t, types.StrategyStatusCancelled, strategy.Status)

	// Verify underlying chain cancelled
	chain, found := k.GetChain(ctx, chainID)
	require.True(t, found)
	require.Equal(t, types.ChainStatusCancelled, chain.Status)
}

func TestCancelStrategy_NotCreator(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	dex.SetSpotPrice(1, "usyreen", "uusdc", math.LegacyMustNewDecFromStr("1.0"))

	msg := makeStrategyMsg(types.StrategySwingTrade, types.RiskModerate)
	strategyID, _, err := k.CreateStrategy(ctx, msg)
	require.NoError(t, err)

	// Try to cancel from different user
	err = k.CancelStrategy(ctx, &types.MsgCancelStrategy{
		Creator:    solverAddr,
		StrategyID: strategyID,
	})
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrStrategyNotCreator)
}

func TestCancelStrategy_NotFound(t *testing.T) {
	k, ctx, _ := setupTradingKeeper(t)

	err := k.CancelStrategy(ctx, &types.MsgCancelStrategy{
		Creator:    creatorAddr,
		StrategyID: "nonexistent",
	})
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrStrategyNotFound)
}

func TestCancelStrategy_AlreadyCancelled(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	dex.SetSpotPrice(1, "usyreen", "uusdc", math.LegacyMustNewDecFromStr("1.0"))

	msg := makeStrategyMsg(types.StrategySwingTrade, types.RiskModerate)
	strategyID, _, err := k.CreateStrategy(ctx, msg)
	require.NoError(t, err)

	// Cancel once
	err = k.CancelStrategy(ctx, &types.MsgCancelStrategy{
		Creator:    creatorAddr,
		StrategyID: strategyID,
	})
	require.NoError(t, err)

	// Cancel again — should fail
	err = k.CancelStrategy(ctx, &types.MsgCancelStrategy{
		Creator:    creatorAddr,
		StrategyID: strategyID,
	})
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrStrategyNotActive)
}

// ---------------------------------------------------------------------------
// Strategy status sync
// ---------------------------------------------------------------------------

func TestSyncStrategyStatuses_Completed(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	dex.SetSpotPrice(1, "usyreen", "uusdc", math.LegacyMustNewDecFromStr("1.0"))
	dex.SetSpotPrice(1, "uusdc", "usyreen", math.LegacyMustNewDecFromStr("1.0"))
	dex.swapOut = sdk.NewCoin("uusdc", math.NewInt(9500000))

	msg := makeStrategyMsg(types.StrategySwingTrade, types.RiskConservative)
	strategyID, chainID, err := k.CreateStrategy(ctx, msg)
	require.NoError(t, err)

	// Execute step 0 (limit buy triggers if price <= target)
	// We need price to match the trigger condition
	// Buy price = 1.0 * (1 - 0.05) = 0.95 — limit buy triggers when spot price of uusdc (in usyreen) <= 0.95
	dex.SetSpotPrice(1, "uusdc", "usyreen", math.LegacyMustNewDecFromStr("0.90"))
	k.ExecuteTradingIntents(ctx)

	// Now step 1 should be pending. Set price for limit sell to trigger.
	// Sell price for step 1 = 1/1.0 * 1.1 = 1.1 — triggers when price of uusdc (in usyreen) >= 1.1
	dex.SetSpotPrice(1, "uusdc", "usyreen", math.LegacyMustNewDecFromStr("1.2"))
	dex.swapOut = sdk.NewCoin("usyreen", math.NewInt(10000000))
	k.ExecuteTradingIntents(ctx)

	// Chain should be completed
	chain, _ := k.GetChain(ctx, chainID)
	require.Equal(t, types.ChainStatusCompleted, chain.Status)

	// Sync strategies
	k.SyncStrategyStatuses(ctx)

	strategy, found := k.GetStrategy(ctx, strategyID)
	require.True(t, found)
	require.Equal(t, types.StrategyStatusCompleted, strategy.Status)
}

// ---------------------------------------------------------------------------
// No dex keeper
// ---------------------------------------------------------------------------

func TestCreateStrategy_NoDexKeeper(t *testing.T) {
	k, ctx := setupKeeper(t)

	msg := makeStrategyMsg(types.StrategySwingTrade, types.RiskModerate)
	_, _, err := k.CreateStrategy(ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "dex keeper not available")
}

// ---------------------------------------------------------------------------
// Price fetch failure
// ---------------------------------------------------------------------------

func TestCreateStrategy_PriceFetchFail(t *testing.T) {
	k, ctx, _ := setupTradingKeeper(t)
	// No price set for pool 1 — should fail

	msg := makeStrategyMsg(types.StrategySwingTrade, types.RiskModerate)
	_, _, err := k.CreateStrategy(ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to fetch current price")
}

// ---------------------------------------------------------------------------
// Query endpoints
// ---------------------------------------------------------------------------

func TestQueryStrategyTemplates(t *testing.T) {
	k, ctx, _ := setupTradingKeeper(t)

	qs := newQueryServer(k)
	resp, err := qs.StrategyTemplates(ctx, &types.QueryStrategyTemplatesRequest{})
	require.NoError(t, err)
	require.Len(t, resp.Templates, 12)
}

func TestQueryStrategies(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	dex.SetSpotPrice(1, "usyreen", "uusdc", math.LegacyMustNewDecFromStr("1.0"))

	msg := makeStrategyMsg(types.StrategySwingTrade, types.RiskModerate)
	_, _, err := k.CreateStrategy(ctx, msg)
	require.NoError(t, err)

	qs := newQueryServer(k)
	resp, err := qs.Strategies(ctx, &types.QueryStrategiesRequest{Address: creatorAddr})
	require.NoError(t, err)
	require.Len(t, resp.Strategies, 1)

	// Empty for other user
	resp2, err := qs.Strategies(ctx, &types.QueryStrategiesRequest{Address: solverAddr})
	require.NoError(t, err)
	require.Len(t, resp2.Strategies, 0)
}

func TestQueryStrategy(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	dex.SetSpotPrice(1, "usyreen", "uusdc", math.LegacyMustNewDecFromStr("1.0"))

	msg := makeStrategyMsg(types.StrategySwingTrade, types.RiskModerate)
	strategyID, _, err := k.CreateStrategy(ctx, msg)
	require.NoError(t, err)

	qs := newQueryServer(k)
	resp, err := qs.Strategy(ctx, &types.QueryStrategyRequest{StrategyID: strategyID})
	require.NoError(t, err)
	require.True(t, resp.Found)
	require.Equal(t, strategyID, resp.Strategy.ID)
	require.Equal(t, types.ChainStatusActive, resp.Chain.Status)

	// Not found
	resp2, err := qs.Strategy(ctx, &types.QueryStrategyRequest{StrategyID: "nonexistent"})
	require.NoError(t, err)
	require.False(t, resp2.Found)
}

// ---------------------------------------------------------------------------
// All risk levels for each template
// ---------------------------------------------------------------------------

func TestCreateStrategy_AllRiskLevels(t *testing.T) {
	templates := []string{
		types.StrategySafeAccumulate,
		types.StrategyMomentumRide,
		types.StrategyGridTrading,
		types.StrategySwingTrade,
		types.StrategyProtectiveSell,
		types.StrategyScalp,
		types.StrategySentimentContrarian,
		types.StrategySmartDCA,
		types.StrategySignalComposite,
		types.StrategyRSIMeanReversion,
		types.StrategyVolatilityBreakout,
		types.StrategyTrendConfluence,
	}

	risks := []types.RiskLevel{
		types.RiskConservative,
		types.RiskModerate,
		types.RiskAggressive,
	}

	for _, tmpl := range templates {
		for _, risk := range risks {
			t.Run(tmpl+"_"+string(risk), func(t *testing.T) {
				k, ctx, dex := setupTradingKeeper(t)
				dex.SetSpotPrice(1, "usyreen", "uusdc", math.LegacyMustNewDecFromStr("1.0"))

				msg := makeStrategyMsg(tmpl, risk)
				strategyID, chainID, err := k.CreateStrategy(ctx, msg)
				require.NoError(t, err)
				require.NotEmpty(t, strategyID)
				require.NotEmpty(t, chainID)

				// Verify chain was created with valid steps
				chain, found := k.GetChain(ctx, chainID)
				require.True(t, found)
				require.True(t, len(chain.Steps) >= 2)
				require.Equal(t, types.ChainStatusActive, chain.Status)
			})
		}
	}
}

// ---------------------------------------------------------------------------
// Msg server integration
// ---------------------------------------------------------------------------

func TestMsgServer_CreateStrategy(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)
	dex.SetSpotPrice(1, "usyreen", "uusdc", math.LegacyMustNewDecFromStr("1.0"))

	ms := newMsgServer(k)
	msg := makeStrategyMsg(types.StrategySwingTrade, types.RiskModerate)
	resp, err := ms.CreateStrategy(ctx, msg)
	require.NoError(t, err)
	require.NotEmpty(t, resp.StrategyID)
	require.NotEmpty(t, resp.ChainID)
}

func TestMsgServer_CancelStrategy(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)
	dex.SetSpotPrice(1, "usyreen", "uusdc", math.LegacyMustNewDecFromStr("1.0"))

	ms := newMsgServer(k)
	createResp, err := ms.CreateStrategy(ctx, makeStrategyMsg(types.StrategySwingTrade, types.RiskModerate))
	require.NoError(t, err)

	cancelResp, err := ms.CancelStrategy(ctx, &types.MsgCancelStrategy{
		Creator:    creatorAddr,
		StrategyID: createResp.StrategyID,
	})
	require.NoError(t, err)
	require.NotNil(t, cancelResp)
}

// ---------------------------------------------------------------------------
// Helpers for msg/query server
// ---------------------------------------------------------------------------

func newQueryServer(k *keeper.Keeper) types.QueryServer {
	return keeper.NewQueryServerImpl(*k)
}

func newMsgServer(k *keeper.Keeper) types.MsgServer {
	return keeper.NewMsgServerImpl(k)
}
