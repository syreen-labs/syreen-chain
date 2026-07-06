package keeper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/intent/types"
)

// ---------------------------------------------------------------------------
// Strategy Storage
// ---------------------------------------------------------------------------

func (k Keeper) SetStrategy(ctx context.Context, strategy types.Strategy) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(strategy)
	kvStore.Set(types.StrategyKey(strategy.ID), bz)
}

func (k Keeper) GetStrategy(ctx context.Context, strategyID string) (types.Strategy, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.StrategyKey(strategyID))
	if err != nil || bz == nil {
		return types.Strategy{}, false
	}
	var strategy types.Strategy
	if err := json.Unmarshal(bz, &strategy); err != nil {
		return types.Strategy{}, false
	}
	return strategy, true
}

func (k Keeper) GetStrategiesByCreator(ctx context.Context, creator string) []types.Strategy {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.StrategyPrefix)

	iter, err := kvStore.Iterator(prefix, nil)
	if err != nil {
		return nil
	}
	defer iter.Close()

	var strategies []types.Strategy
	for ; iter.Valid(); iter.Next() {
		if !bytes.HasPrefix(iter.Key(), prefix) {
			break
		}
		var strategy types.Strategy
		if err := json.Unmarshal(iter.Value(), &strategy); err != nil {
			continue
		}
		if strategy.Creator == creator {
			strategies = append(strategies, strategy)
		}
	}

	sort.Slice(strategies, func(i, j int) bool {
		return strategies[i].ID < strategies[j].ID
	})

	return strategies
}

func (k Keeper) nextStrategyID(ctx context.Context) string {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.StrategyCounterKey())
	var counter uint64
	if err == nil && bz != nil {
		counter = types.BytesToUint64(bz)
	}
	counter++
	kvStore.Set(types.StrategyCounterKey(), types.Uint64ToBytes(counter))
	return "strategy-" + strconv.FormatUint(counter, 10)
}

// ---------------------------------------------------------------------------
// Strategy Creation
// ---------------------------------------------------------------------------

// Sensible defaults for DCA-based templates when the caller omits NumExecutions
// or IntervalBlocks. Match the original live mainnet behavior (10 executions,
// 600 blocks ~ 50 min apart).
const (
	defaultNumExecutions  uint64 = 10
	defaultIntervalBlocks uint64 = 600
	// Safety buffer added on top of the computed min expiry so that the chain
	// still has time to complete the final step after the last DCA tranche fires.
	dcaExpirySafetyBuffer uint64 = 1200
)

// CreateStrategy generates an intent chain from a strategy template and deploys it.
func (k *Keeper) CreateStrategy(ctx context.Context, msg *types.MsgCreateStrategy) (string, string, error) {
	if k.dexKeeper == nil {
		return "", "", fmt.Errorf("dex keeper not available")
	}

	// Apply sensible defaults for DCA-based templates. Previously, a frontend
	// that omitted these fields (or a user CLI call) would cause divide-by-zero
	// panics in tryDCA or would create strategies that literally cannot execute.
	if msg.NumExecutions == 0 {
		msg.NumExecutions = defaultNumExecutions
	}
	if msg.IntervalBlocks == 0 {
		msg.IntervalBlocks = defaultIntervalBlocks
	}

	// Enforce minimum expiry: every DCA-based template needs enough blocks to
	// fire all tranches plus a safety buffer for the final step. Without this
	// check, a user (or frontend) can submit a strategy that expires before
	// the DCA schedule can complete, causing only tranche 1 to fire and the
	// remaining budget to refund early. Auto-correct upward rather than
	// rejecting — the UX goal is "it just works".
	minExpiry := msg.NumExecutions*msg.IntervalBlocks + dcaExpirySafetyBuffer
	if msg.ExpiryBlocks < minExpiry {
		k.Logger(ctx).Info("strategy expiry too short for DCA schedule; auto-extending",
			"template", msg.TemplateName,
			"requested_expiry", msg.ExpiryBlocks,
			"min_expiry", minExpiry,
			"num_executions", msg.NumExecutions,
			"interval_blocks", msg.IntervalBlocks,
		)
		msg.ExpiryBlocks = minExpiry
	}

	// Get current spot price for the pool
	currentPrice, err := k.dexKeeper.GetSpotPrice(ctx, msg.PoolID, msg.InputDenom, msg.OutputDenom)
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", types.ErrStrategyPriceFetch, err)
	}

	// Generate chain steps based on template
	steps, err := k.generateSteps(msg, currentPrice)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate strategy steps: %w", err)
	}

	// Submit as an intent chain
	chainMsg := &types.MsgSubmitChain{
		Creator:      msg.Creator,
		Steps:        steps,
		MaxFee:       sdk.NewCoins(),
		Tip:          sdk.NewCoins(),
		ExpiryBlocks: msg.ExpiryBlocks,
	}

	chainID, err := k.SubmitChain(ctx, chainMsg)
	if err != nil {
		return "", "", fmt.Errorf("failed to submit strategy chain: %w", err)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	strategyID := k.nextStrategyID(ctx)

	strategy := types.Strategy{
		ID:           strategyID,
		Creator:      msg.Creator,
		TemplateName: msg.TemplateName,
		ChainID:      chainID,
		Status:       types.StrategyStatusActive,
		PoolID:       msg.PoolID,
		InputDenom:   msg.InputDenom,
		OutputDenom:  msg.OutputDenom,
		TotalBudget:  msg.TotalBudget,
		RiskLevel:    msg.RiskLevel,
		CreatedAt:    sdkCtx.BlockHeight(),
	}

	k.SetStrategy(ctx, strategy)

	k.Logger(ctx).Info("strategy created",
		"strategy_id", strategyID,
		"chain_id", chainID,
		"template", msg.TemplateName,
		"creator", msg.Creator,
		"risk", msg.RiskLevel,
		"budget", msg.TotalBudget,
	)

	return strategyID, chainID, nil
}

// CancelStrategy cancels a strategy by cancelling its underlying intent chain.
func (k *Keeper) CancelStrategy(ctx context.Context, msg *types.MsgCancelStrategy) error {
	strategy, found := k.GetStrategy(ctx, msg.StrategyID)
	if !found {
		return types.ErrStrategyNotFound
	}
	if strategy.Creator != msg.Creator {
		return types.ErrStrategyNotCreator
	}
	if strategy.Status != types.StrategyStatusActive {
		return types.ErrStrategyNotActive
	}

	// Cancel the underlying chain
	err := k.CancelChain(ctx, &types.MsgCancelChain{
		Creator: msg.Creator,
		ChainID: strategy.ChainID,
	})
	if err != nil {
		return fmt.Errorf("failed to cancel strategy chain: %w", err)
	}

	strategy.Status = types.StrategyStatusCancelled
	k.SetStrategy(ctx, strategy)

	k.Logger(ctx).Info("strategy cancelled",
		"strategy_id", msg.StrategyID,
		"creator", msg.Creator,
	)

	return nil
}

// SyncStrategyStatus updates a strategy's status based on its underlying chain status.
// Called during BeginBlock.
func (k *Keeper) SyncStrategyStatuses(ctx context.Context) {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.StrategyPrefix)

	iter, err := kvStore.Iterator(prefix, nil)
	if err != nil {
		return
	}

	// Collect the strategies whose status changed while iterating, then persist
	// them after the iterator is closed. Calling SetStrategy (a store write) on a
	// key under the live iterator is unsafe.
	var toUpdate []types.Strategy
	for ; iter.Valid(); iter.Next() {
		if !bytes.HasPrefix(iter.Key(), prefix) {
			break
		}
		var strategy types.Strategy
		if err := json.Unmarshal(iter.Value(), &strategy); err != nil {
			continue
		}
		if strategy.Status != types.StrategyStatusActive {
			continue
		}

		chain, found := k.GetChain(ctx, strategy.ChainID)
		if !found {
			continue
		}

		changed := true
		switch chain.Status {
		case types.ChainStatusCompleted:
			strategy.Status = types.StrategyStatusCompleted
		case types.ChainStatusFailed:
			strategy.Status = types.StrategyStatusFailed
		case types.ChainStatusExpired:
			strategy.Status = types.StrategyStatusFailed
		case types.ChainStatusCancelled:
			strategy.Status = types.StrategyStatusCancelled
		default:
			changed = false
		}
		if changed {
			toUpdate = append(toUpdate, strategy)
		}
	}
	iter.Close()

	for _, strategy := range toUpdate {
		k.SetStrategy(ctx, strategy)
	}
}

// ---------------------------------------------------------------------------
// Step Generation — builds ChainStep slices for each template
// ---------------------------------------------------------------------------

func (k *Keeper) generateSteps(msg *types.MsgCreateStrategy, currentPrice math.LegacyDec) ([]types.ChainStep, error) {
	switch msg.TemplateName {
	case types.StrategySafeAccumulate:
		return k.generateSafeAccumulate(msg, currentPrice)
	case types.StrategyMomentumRide:
		return k.generateMomentumRide(msg, currentPrice)
	case types.StrategyGridTrading:
		return k.generateGridTrading(msg, currentPrice)
	case types.StrategySwingTrade:
		return k.generateSwingTrade(msg, currentPrice)
	case types.StrategyProtectiveSell:
		return k.generateProtectiveSell(msg, currentPrice)
	case types.StrategyScalp:
		return k.generateScalp(msg, currentPrice)
	case types.StrategySentimentContrarian:
		return k.generateSentimentContrarian(msg, currentPrice)
	case types.StrategySmartDCA:
		return k.generateSmartDCA(msg, currentPrice)
	case types.StrategySignalComposite:
		return k.generateSignalComposite(msg, currentPrice)
	case types.StrategyRSIMeanReversion:
		return k.generateRSIMeanReversion(msg, currentPrice)
	case types.StrategyVolatilityBreakout:
		return k.generateVolatilityBreakout(msg, currentPrice)
	case types.StrategyTrendConfluence:
		return k.generateTrendConfluence(msg, currentPrice)
	default:
		return nil, fmt.Errorf("unknown strategy template: %s", msg.TemplateName)
	}
}

// safe_accumulate: Step 0 = DCA buy, Step 1 = stop-loss on accumulated tokens
func (k *Keeper) generateSafeAccumulate(msg *types.MsgCreateStrategy, currentPrice math.LegacyDec) ([]types.ChainStep, error) {
	stopLossPct, _, _ := types.GetRiskParams(msg.RiskLevel)

	// Step 0: DCA buy — gradually accumulate OutputDenom using InputDenom
	dcaBody := types.DCAIntent{
		InputDenom:     msg.InputDenom,
		TotalAmount:    msg.TotalBudget,
		OutputDenom:    msg.OutputDenom,
		NumExecutions:  msg.NumExecutions,
		IntervalBlocks: msg.IntervalBlocks,
		PoolID:         msg.PoolID,
	}
	dcaBz, _ := json.Marshal(dcaBody)
	step0 := types.ChainStep{
		IntentType: types.IntentTypeDCA,
		Body:       dcaBz,
		Condition:  types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{
			Type:   types.InputFromFixed,
			Amount: msg.TotalBudget,
			Denom:  msg.InputDenom,
		},
	}

	// Step 1: stop-loss on accumulated output — sell if price drops below threshold
	// Stop price = currentPrice * (1 - stopLossPct)
	// The price here is OutputDenom priced in InputDenom
	// We need the inverse: price of OutputDenom in terms of InputDenom
	if currentPrice.IsZero() {
		return nil, fmt.Errorf("current price is zero, cannot calculate stop-loss")
	}
	outputPrice := math.LegacyOneDec().Quo(currentPrice) // price of 1 OutputDenom in InputDenom
	stopPrice := outputPrice.Mul(math.LegacyOneDec().Sub(stopLossPct))

	stopBody := types.StopLossIntent{
		InputDenom:      msg.OutputDenom, // selling what we accumulated
		InputAmount:     math.OneInt(),    // placeholder, rewritten from previous
		OutputDenom:     msg.InputDenom,
		StopPrice:       stopPrice,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	stopBz, _ := json.Marshal(stopBody)
	step1 := types.ChainStep{
		IntentType:  types.IntentTypeStopLoss,
		Body:        stopBz,
		Condition:   types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{Type: types.InputFromPrevious},
	}

	return []types.ChainStep{step0, step1}, nil
}

// momentum_ride: Step 0 = limit buy on dip, Step 1 = take-profit (portion), Step 2 = stop-loss (rest)
func (k *Keeper) generateMomentumRide(msg *types.MsgCreateStrategy, currentPrice math.LegacyDec) ([]types.ChainStep, error) {
	stopLossPct, takeProfitPct, _ := types.GetRiskParams(msg.RiskLevel)

	// Buy target: below current price by stopLossPct (buy the dip)
	buyPrice := currentPrice.Mul(math.LegacyOneDec().Sub(stopLossPct))

	// Step 0: limit buy on dip
	buyBody := types.LimitBuyIntent{
		InputDenom:      msg.InputDenom,
		InputAmount:     msg.TotalBudget,
		OutputDenom:     msg.OutputDenom,
		TargetPrice:     buyPrice,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	buyBz, _ := json.Marshal(buyBody)
	step0 := types.ChainStep{
		IntentType: types.IntentTypeLimitBuy,
		Body:       buyBz,
		Condition:  types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{
			Type:   types.InputFromFixed,
			Amount: msg.TotalBudget,
			Denom:  msg.InputDenom,
		},
	}

	// Sell target: above buy price by takeProfitPct
	if currentPrice.IsZero() {
		return nil, fmt.Errorf("current price is zero")
	}
	outputPrice := math.LegacyOneDec().Quo(currentPrice)
	takeProfitPrice := outputPrice.Mul(math.LegacyOneDec().Add(takeProfitPct))
	stopLossPrice := outputPrice.Mul(math.LegacyOneDec().Sub(stopLossPct))

	// Step 1: take-profit sell 50% of output
	tpBody := types.TakeProfitIntent{
		InputDenom:      msg.OutputDenom,
		InputAmount:     math.OneInt(), // rewritten
		OutputDenom:     msg.InputDenom,
		TargetPrice:     takeProfitPrice,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	tpBz, _ := json.Marshal(tpBody)
	step1 := types.ChainStep{
		IntentType: types.IntentTypeTakeProfit,
		Body:       tpBz,
		Condition:  types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{
			Type:           types.InputFromPortion,
			PortionPercent: math.LegacyMustNewDecFromStr("0.5"),
		},
	}

	// Step 2: stop-loss on remaining 50%
	slBody := types.StopLossIntent{
		InputDenom:      msg.OutputDenom,
		InputAmount:     math.OneInt(), // rewritten
		OutputDenom:     msg.InputDenom,
		StopPrice:       stopLossPrice,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	slBz, _ := json.Marshal(slBody)
	step2 := types.ChainStep{
		IntentType:  types.IntentTypeStopLoss,
		Body:        slBz,
		Condition:   types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{Type: types.InputFromPrevious},
	}

	return []types.ChainStep{step0, step1, step2}, nil
}

// grid_trading: alternating buy/sell orders at price intervals around current price
func (k *Keeper) generateGridTrading(msg *types.MsgCreateStrategy, currentPrice math.LegacyDec) ([]types.ChainStep, error) {
	_, _, gridSpreadPct := types.GetRiskParams(msg.RiskLevel)
	levels := msg.GridLevels
	if levels < 2 {
		levels = 2
	}
	if levels > 5 {
		levels = 5
	}

	// Budget per level pair
	budgetPerPair := msg.TotalBudget.QuoRaw(int64(levels))
	if budgetPerPair.IsZero() {
		return nil, fmt.Errorf("budget too small for %d grid levels", levels)
	}

	var steps []types.ChainStep

	for i := uint64(0); i < levels; i++ {
		// Buy price = currentPrice * (1 - gridSpreadPct * (i+1))
		discount := gridSpreadPct.MulInt64(int64(i + 1))
		buyPrice := currentPrice.Mul(math.LegacyOneDec().Sub(discount))
		if buyPrice.IsNegative() || buyPrice.IsZero() {
			break
		}

		// Step: limit buy at discount
		buyBody := types.LimitBuyIntent{
			InputDenom:      msg.InputDenom,
			InputAmount:     budgetPerPair,
			OutputDenom:     msg.OutputDenom,
			TargetPrice:     buyPrice,
			PoolID:          msg.PoolID,
			MinOutputAmount: math.ZeroInt(),
		}
		buyBz, _ := json.Marshal(buyBody)

		condition := types.ChainCondition{Type: types.ConditionImmediate}
		inputSource := types.InputSource{
			Type:   types.InputFromFixed,
			Amount: budgetPerPair,
			Denom:  msg.InputDenom,
		}

		// For step 0, condition must be immediate; for subsequent, also immediate
		// since each grid level operates independently with fixed input
		step := types.ChainStep{
			IntentType:  types.IntentTypeLimitBuy,
			Body:        buyBz,
			Condition:   condition,
			InputSource: inputSource,
		}

		// Only the first step can have immediate + fixed input without prefix constraint
		if len(steps) == 0 {
			steps = append(steps, step)
		} else {
			// Subsequent buy steps use fixed input with immediate condition
			steps = append(steps, step)
		}

		// Corresponding sell: sell the output from the buy at a premium
		premium := gridSpreadPct.MulInt64(int64(i + 1))
		if currentPrice.IsZero() {
			return nil, fmt.Errorf("current price is zero")
		}
		outputPriceForSell := math.LegacyOneDec().Quo(currentPrice)
		sellPrice := outputPriceForSell.Mul(math.LegacyOneDec().Add(premium))

		sellBody := types.LimitSellIntent{
			InputDenom:      msg.OutputDenom,
			InputAmount:     math.OneInt(), // rewritten from previous
			OutputDenom:     msg.InputDenom,
			TargetPrice:     sellPrice,
			PoolID:          msg.PoolID,
			MinOutputAmount: math.ZeroInt(),
		}
		sellBz, _ := json.Marshal(sellBody)

		sellStep := types.ChainStep{
			IntentType:  types.IntentTypeLimitSell,
			Body:        sellBz,
			Condition:   types.ChainCondition{Type: types.ConditionImmediate},
			InputSource: types.InputSource{Type: types.InputFromPrevious},
		}
		steps = append(steps, sellStep)

		// Enforce chain max steps
		if len(steps) >= 10 {
			break
		}
	}

	if len(steps) < 2 {
		return nil, fmt.Errorf("grid trading requires at least 2 steps")
	}

	return steps, nil
}

// swing_trade: Step 0 = limit buy below current price, Step 1 = limit sell above current price
func (k *Keeper) generateSwingTrade(msg *types.MsgCreateStrategy, currentPrice math.LegacyDec) ([]types.ChainStep, error) {
	stopLossPct, takeProfitPct, _ := types.GetRiskParams(msg.RiskLevel)

	// Buy at a discount
	buyPrice := currentPrice.Mul(math.LegacyOneDec().Sub(stopLossPct))

	buyBody := types.LimitBuyIntent{
		InputDenom:      msg.InputDenom,
		InputAmount:     msg.TotalBudget,
		OutputDenom:     msg.OutputDenom,
		TargetPrice:     buyPrice,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	buyBz, _ := json.Marshal(buyBody)
	step0 := types.ChainStep{
		IntentType: types.IntentTypeLimitBuy,
		Body:       buyBz,
		Condition:  types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{
			Type:   types.InputFromFixed,
			Amount: msg.TotalBudget,
			Denom:  msg.InputDenom,
		},
	}

	// Sell at a premium
	if currentPrice.IsZero() {
		return nil, fmt.Errorf("current price is zero")
	}
	outputPrice := math.LegacyOneDec().Quo(currentPrice)
	sellPrice := outputPrice.Mul(math.LegacyOneDec().Add(takeProfitPct))

	sellBody := types.LimitSellIntent{
		InputDenom:      msg.OutputDenom,
		InputAmount:     math.OneInt(), // rewritten
		OutputDenom:     msg.InputDenom,
		TargetPrice:     sellPrice,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	sellBz, _ := json.Marshal(sellBody)
	step1 := types.ChainStep{
		IntentType:  types.IntentTypeLimitSell,
		Body:        sellBz,
		Condition:   types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{Type: types.InputFromPrevious},
	}

	return []types.ChainStep{step0, step1}, nil
}

// protective_sell: Step 0 = take-profit sell, Step 1 = stop-loss sell (safety net)
// User already holds OutputDenom and wants to sell it — InputDenom here is the token they hold
func (k *Keeper) generateProtectiveSell(msg *types.MsgCreateStrategy, currentPrice math.LegacyDec) ([]types.ChainStep, error) {
	stopLossPct, takeProfitPct, _ := types.GetRiskParams(msg.RiskLevel)

	// The user holds InputDenom and wants to sell to OutputDenom
	// currentPrice = InputDenom per OutputDenom
	// Price of InputDenom in OutputDenom = 1/currentPrice
	if currentPrice.IsZero() {
		return nil, fmt.Errorf("current price is zero")
	}
	sellPrice := math.LegacyOneDec().Quo(currentPrice)
	takeProfitTarget := sellPrice.Mul(math.LegacyOneDec().Add(takeProfitPct))
	stopLossTarget := sellPrice.Mul(math.LegacyOneDec().Sub(stopLossPct))

	halfBudget := msg.TotalBudget.QuoRaw(2)
	remainderBudget := msg.TotalBudget.Sub(halfBudget)

	// Step 0: take-profit sell half
	tpBody := types.TakeProfitIntent{
		InputDenom:      msg.InputDenom,
		InputAmount:     halfBudget,
		OutputDenom:     msg.OutputDenom,
		TargetPrice:     takeProfitTarget,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	tpBz, _ := json.Marshal(tpBody)
	step0 := types.ChainStep{
		IntentType: types.IntentTypeTakeProfit,
		Body:       tpBz,
		Condition:  types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{
			Type:   types.InputFromFixed,
			Amount: halfBudget,
			Denom:  msg.InputDenom,
		},
	}

	// Step 1: stop-loss sell the other half
	slBody := types.StopLossIntent{
		InputDenom:      msg.InputDenom,
		InputAmount:     remainderBudget,
		OutputDenom:     msg.OutputDenom,
		StopPrice:       stopLossTarget,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	slBz, _ := json.Marshal(slBody)
	step1 := types.ChainStep{
		IntentType: types.IntentTypeStopLoss,
		Body:       slBz,
		Condition:  types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{
			Type:   types.InputFromFixed,
			Amount: remainderBudget,
			Denom:  msg.InputDenom,
		},
	}

	return []types.ChainStep{step0, step1}, nil
}

// scalp: Step 0 = limit buy at small discount, Step 1 = take-profit at small premium (tight)
func (k *Keeper) generateScalp(msg *types.MsgCreateStrategy, currentPrice math.LegacyDec) ([]types.ChainStep, error) {
	// Scalp uses very tight margins: half of the risk params
	stopLossPct, _, _ := types.GetRiskParams(msg.RiskLevel)
	scalpMargin := stopLossPct.Quo(math.LegacyNewDec(2)) // half the stop-loss pct as profit target

	// Buy at a tiny discount
	buyPrice := currentPrice.Mul(math.LegacyOneDec().Sub(scalpMargin))

	buyBody := types.LimitBuyIntent{
		InputDenom:      msg.InputDenom,
		InputAmount:     msg.TotalBudget,
		OutputDenom:     msg.OutputDenom,
		TargetPrice:     buyPrice,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	buyBz, _ := json.Marshal(buyBody)
	step0 := types.ChainStep{
		IntentType: types.IntentTypeLimitBuy,
		Body:       buyBz,
		Condition:  types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{
			Type:   types.InputFromFixed,
			Amount: msg.TotalBudget,
			Denom:  msg.InputDenom,
		},
	}

	// Take profit at a tiny premium
	if currentPrice.IsZero() {
		return nil, fmt.Errorf("current price is zero")
	}
	outputPrice := math.LegacyOneDec().Quo(currentPrice)
	takeProfitPrice := outputPrice.Mul(math.LegacyOneDec().Add(scalpMargin))

	tpBody := types.TakeProfitIntent{
		InputDenom:      msg.OutputDenom,
		InputAmount:     math.OneInt(), // rewritten
		OutputDenom:     msg.InputDenom,
		TargetPrice:     takeProfitPrice,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	tpBz, _ := json.Marshal(tpBody)
	step1 := types.ChainStep{
		IntentType:  types.IntentTypeTakeProfit,
		Body:        tpBz,
		Condition:   types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{Type: types.InputFromPrevious},
	}

	return []types.ChainStep{step0, step1}, nil
}

// ---------------------------------------------------------------------------
// Tier 1 AI-reactive templates
//
// Each of these generators produces a chain whose first step is an Immediate
// anchor action (required by ValidateChainStep) and whose subsequent steps
// depend on live dex signal/sentiment data — evaluated block-by-block in
// ProcessChainConditions. This is what makes them "AI-reactive" rather than
// frozen-at-create-time.
// ---------------------------------------------------------------------------

// splitBudget10_90 returns (10% anchor, 90% main) rounded so they sum to total.
func splitBudget10_90(total math.Int) (anchor, main math.Int) {
	anchor = total.QuoRaw(10)
	if anchor.IsZero() {
		anchor = math.OneInt()
	}
	main = total.Sub(anchor)
	return
}

// sentiment_contrarian:
//
//	Step 0: small anchor limit_buy (10% budget), immediate
//	Step 1: main limit_buy (90% budget), when F&G <= ExtremeFear (25)
//	Step 2: take_profit on held position, when F&G >= ExtremeGreed (75)
//	Step 3: stop_loss safety net, immediate
func (k *Keeper) generateSentimentContrarian(msg *types.MsgCreateStrategy, currentPrice math.LegacyDec) ([]types.ChainStep, error) {
	if currentPrice.IsZero() {
		return nil, fmt.Errorf("current price is zero")
	}
	stopLossPct, takeProfitPct, _ := types.GetRiskParams(msg.RiskLevel)

	anchorBudget, mainBudget := splitBudget10_90(msg.TotalBudget)

	buyPrice := currentPrice.Mul(math.LegacyOneDec().Add(math.LegacyMustNewDecFromStr("0.005")))

	anchorBody := types.LimitBuyIntent{
		InputDenom:      msg.InputDenom,
		InputAmount:     anchorBudget,
		OutputDenom:     msg.OutputDenom,
		TargetPrice:     buyPrice,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	anchorBz, _ := json.Marshal(anchorBody)
	step0 := types.ChainStep{
		IntentType:  types.IntentTypeLimitBuy,
		Body:        anchorBz,
		Condition:   types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{Type: types.InputFromFixed, Amount: anchorBudget, Denom: msg.InputDenom},
	}

	mainBody := types.LimitBuyIntent{
		InputDenom:      msg.InputDenom,
		InputAmount:     mainBudget,
		OutputDenom:     msg.OutputDenom,
		TargetPrice:     buyPrice,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	mainBz, _ := json.Marshal(mainBody)
	step1 := types.ChainStep{
		IntentType: types.IntentTypeLimitBuy,
		Body:       mainBz,
		Condition: types.ChainCondition{
			Type:               types.ConditionFearGreedBelow,
			PoolID:             msg.PoolID,
			FearGreedThreshold: types.FearGreedExtremeFear,
		},
		InputSource: types.InputSource{Type: types.InputFromFixed, Amount: mainBudget, Denom: msg.InputDenom},
	}

	outputPrice := math.LegacyOneDec().Quo(currentPrice)
	takeProfitPrice := outputPrice.Mul(math.LegacyOneDec().Add(takeProfitPct))
	stopPrice := outputPrice.Mul(math.LegacyOneDec().Sub(stopLossPct))

	tpBody := types.TakeProfitIntent{
		InputDenom:      msg.OutputDenom,
		InputAmount:     math.OneInt(),
		OutputDenom:     msg.InputDenom,
		TargetPrice:     takeProfitPrice,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	tpBz, _ := json.Marshal(tpBody)
	step2 := types.ChainStep{
		IntentType: types.IntentTypeTakeProfit,
		Body:       tpBz,
		Condition: types.ChainCondition{
			Type:               types.ConditionFearGreedAbove,
			PoolID:             msg.PoolID,
			FearGreedThreshold: types.FearGreedExtremeGreed,
		},
		InputSource: types.InputSource{
			Type:           types.InputFromPortion,
			PortionPercent: math.LegacyMustNewDecFromStr("0.5"),
		},
	}

	slBody := types.StopLossIntent{
		InputDenom:      msg.OutputDenom,
		InputAmount:     math.OneInt(),
		OutputDenom:     msg.InputDenom,
		StopPrice:       stopPrice,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	slBz, _ := json.Marshal(slBody)
	step3 := types.ChainStep{
		IntentType:  types.IntentTypeStopLoss,
		Body:        slBz,
		Condition:   types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{Type: types.InputFromPrevious},
	}

	return []types.ChainStep{step0, step1, step2, step3}, nil
}

// smart_dca:
//
//	Step 0: DCA 50% at base pace (immediate)
//	Step 1: DCA 30% when F&G <= 40 (fear)
//	Step 2: DCA 20% when F&G <= 25 (extreme fear)
func (k *Keeper) generateSmartDCA(msg *types.MsgCreateStrategy, _ math.LegacyDec) ([]types.ChainStep, error) {
	if msg.NumExecutions == 0 {
		return nil, fmt.Errorf("num_executions required for smart_dca strategy")
	}
	if msg.IntervalBlocks == 0 {
		return nil, fmt.Errorf("interval_blocks required for smart_dca strategy")
	}

	baseBudget := msg.TotalBudget.MulRaw(50).QuoRaw(100)
	fearBudget := msg.TotalBudget.MulRaw(30).QuoRaw(100)
	extremeBudget := msg.TotalBudget.Sub(baseBudget).Sub(fearBudget)
	if baseBudget.IsZero() || fearBudget.IsZero() || extremeBudget.IsZero() {
		return nil, fmt.Errorf("budget too small for smart_dca tranches")
	}

	base := types.DCAIntent{
		InputDenom:     msg.InputDenom,
		TotalAmount:    baseBudget,
		OutputDenom:    msg.OutputDenom,
		NumExecutions:  msg.NumExecutions,
		IntervalBlocks: msg.IntervalBlocks,
		PoolID:         msg.PoolID,
	}
	baseBz, _ := json.Marshal(base)
	step0 := types.ChainStep{
		IntentType:  types.IntentTypeDCA,
		Body:        baseBz,
		Condition:   types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{Type: types.InputFromFixed, Amount: baseBudget, Denom: msg.InputDenom},
	}

	fearExecs := msg.NumExecutions / 2
	if fearExecs == 0 {
		fearExecs = 1
	}
	fear := types.DCAIntent{
		InputDenom:     msg.InputDenom,
		TotalAmount:    fearBudget,
		OutputDenom:    msg.OutputDenom,
		NumExecutions:  fearExecs,
		IntervalBlocks: msg.IntervalBlocks,
		PoolID:         msg.PoolID,
	}
	fearBz, _ := json.Marshal(fear)
	step1 := types.ChainStep{
		IntentType: types.IntentTypeDCA,
		Body:       fearBz,
		Condition: types.ChainCondition{
			Type:               types.ConditionFearGreedBelow,
			PoolID:             msg.PoolID,
			FearGreedThreshold: types.FearGreedFear,
		},
		InputSource: types.InputSource{Type: types.InputFromFixed, Amount: fearBudget, Denom: msg.InputDenom},
	}

	extremeExecs := msg.NumExecutions / 4
	if extremeExecs == 0 {
		extremeExecs = 1
	}
	extreme := types.DCAIntent{
		InputDenom:     msg.InputDenom,
		TotalAmount:    extremeBudget,
		OutputDenom:    msg.OutputDenom,
		NumExecutions:  extremeExecs,
		IntervalBlocks: msg.IntervalBlocks,
		PoolID:         msg.PoolID,
	}
	extremeBz, _ := json.Marshal(extreme)
	step2 := types.ChainStep{
		IntentType: types.IntentTypeDCA,
		Body:       extremeBz,
		Condition: types.ChainCondition{
			Type:               types.ConditionFearGreedBelow,
			PoolID:             msg.PoolID,
			FearGreedThreshold: types.FearGreedExtremeFear,
		},
		InputSource: types.InputSource{Type: types.InputFromFixed, Amount: extremeBudget, Denom: msg.InputDenom},
	}

	return []types.ChainStep{step0, step1, step2}, nil
}

// signal_composite:
//
//	Step 0: entry limit_buy (immediate)
//	Step 1: take_profit portion on SELL signal
//	Step 2: stop_loss on STRONG_SELL signal
func (k *Keeper) generateSignalComposite(msg *types.MsgCreateStrategy, currentPrice math.LegacyDec) ([]types.ChainStep, error) {
	if currentPrice.IsZero() {
		return nil, fmt.Errorf("current price is zero")
	}
	stopLossPct, takeProfitPct, _ := types.GetRiskParams(msg.RiskLevel)

	entryPrice := currentPrice.Mul(math.LegacyOneDec().Add(math.LegacyMustNewDecFromStr("0.005")))

	buyBody := types.LimitBuyIntent{
		InputDenom:      msg.InputDenom,
		InputAmount:     msg.TotalBudget,
		OutputDenom:     msg.OutputDenom,
		TargetPrice:     entryPrice,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	buyBz, _ := json.Marshal(buyBody)
	step0 := types.ChainStep{
		IntentType:  types.IntentTypeLimitBuy,
		Body:        buyBz,
		Condition:   types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{Type: types.InputFromFixed, Amount: msg.TotalBudget, Denom: msg.InputDenom},
	}

	outputPrice := math.LegacyOneDec().Quo(currentPrice)
	tpTarget := outputPrice.Mul(math.LegacyOneDec().Add(takeProfitPct))
	slTarget := outputPrice.Mul(math.LegacyOneDec().Sub(stopLossPct))

	tpBody := types.TakeProfitIntent{
		InputDenom:      msg.OutputDenom,
		InputAmount:     math.OneInt(),
		OutputDenom:     msg.InputDenom,
		TargetPrice:     tpTarget,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	tpBz, _ := json.Marshal(tpBody)
	step1 := types.ChainStep{
		IntentType: types.IntentTypeTakeProfit,
		Body:       tpBz,
		Condition: types.ChainCondition{
			Type:        types.ConditionSignalEquals,
			PoolID:      msg.PoolID,
			SignalLabel: "SELL",
		},
		InputSource: types.InputSource{
			Type:           types.InputFromPortion,
			PortionPercent: math.LegacyMustNewDecFromStr("0.5"),
		},
	}

	slBody := types.StopLossIntent{
		InputDenom:      msg.OutputDenom,
		InputAmount:     math.OneInt(),
		OutputDenom:     msg.InputDenom,
		StopPrice:       slTarget,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	slBz, _ := json.Marshal(slBody)
	step2 := types.ChainStep{
		IntentType: types.IntentTypeStopLoss,
		Body:       slBz,
		Condition: types.ChainCondition{
			Type:        types.ConditionSignalEquals,
			PoolID:      msg.PoolID,
			SignalLabel: "STRONG_SELL",
		},
		InputSource: types.InputSource{Type: types.InputFromPrevious},
	}

	return []types.ChainStep{step0, step1, step2}, nil
}

// rsi_mean_reversion:
//
//	Step 0: anchor limit_buy (10%), immediate
//	Step 1: main limit_buy (90%) when RSI <= 30
//	Step 2: take_profit when RSI >= 70
//	Step 3: stop_loss safety net, immediate
func (k *Keeper) generateRSIMeanReversion(msg *types.MsgCreateStrategy, currentPrice math.LegacyDec) ([]types.ChainStep, error) {
	if currentPrice.IsZero() {
		return nil, fmt.Errorf("current price is zero")
	}
	stopLossPct, takeProfitPct, _ := types.GetRiskParams(msg.RiskLevel)

	anchorBudget, mainBudget := splitBudget10_90(msg.TotalBudget)

	buyPrice := currentPrice.Mul(math.LegacyOneDec().Add(math.LegacyMustNewDecFromStr("0.005")))

	anchorBody := types.LimitBuyIntent{
		InputDenom:      msg.InputDenom,
		InputAmount:     anchorBudget,
		OutputDenom:     msg.OutputDenom,
		TargetPrice:     buyPrice,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	anchorBz, _ := json.Marshal(anchorBody)
	step0 := types.ChainStep{
		IntentType:  types.IntentTypeLimitBuy,
		Body:        anchorBz,
		Condition:   types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{Type: types.InputFromFixed, Amount: anchorBudget, Denom: msg.InputDenom},
	}

	mainBody := types.LimitBuyIntent{
		InputDenom:      msg.InputDenom,
		InputAmount:     mainBudget,
		OutputDenom:     msg.OutputDenom,
		TargetPrice:     buyPrice,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	mainBz, _ := json.Marshal(mainBody)
	step1 := types.ChainStep{
		IntentType: types.IntentTypeLimitBuy,
		Body:       mainBz,
		Condition: types.ChainCondition{
			Type:         types.ConditionRSIBelow,
			PoolID:       msg.PoolID,
			RSIThreshold: types.RSIOversoldThreshold,
		},
		InputSource: types.InputSource{Type: types.InputFromFixed, Amount: mainBudget, Denom: msg.InputDenom},
	}

	outputPrice := math.LegacyOneDec().Quo(currentPrice)
	takeProfitPrice := outputPrice.Mul(math.LegacyOneDec().Add(takeProfitPct))
	stopPrice := outputPrice.Mul(math.LegacyOneDec().Sub(stopLossPct))

	tpBody := types.TakeProfitIntent{
		InputDenom:      msg.OutputDenom,
		InputAmount:     math.OneInt(),
		OutputDenom:     msg.InputDenom,
		TargetPrice:     takeProfitPrice,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	tpBz, _ := json.Marshal(tpBody)
	step2 := types.ChainStep{
		IntentType: types.IntentTypeTakeProfit,
		Body:       tpBz,
		Condition: types.ChainCondition{
			Type:         types.ConditionRSIAbove,
			PoolID:       msg.PoolID,
			RSIThreshold: types.RSIOverboughtThreshold,
		},
		InputSource: types.InputSource{
			Type:           types.InputFromPortion,
			PortionPercent: math.LegacyMustNewDecFromStr("0.5"),
		},
	}

	slBody := types.StopLossIntent{
		InputDenom:      msg.OutputDenom,
		InputAmount:     math.OneInt(),
		OutputDenom:     msg.InputDenom,
		StopPrice:       stopPrice,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	slBz, _ := json.Marshal(slBody)
	step3 := types.ChainStep{
		IntentType:  types.IntentTypeStopLoss,
		Body:        slBz,
		Condition:   types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{Type: types.InputFromPrevious},
	}

	return []types.ChainStep{step0, step1, step2, step3}, nil
}

// volatility_breakout:
//
//	Step 0: anchor limit_buy (10%), immediate
//	Step 1: main limit_buy (90%) when volatility >= breakout threshold
//	Step 2: wider take_profit (2x takeProfitPct) (portion)
//	Step 3: stop_loss safety net (immediate)
func (k *Keeper) generateVolatilityBreakout(msg *types.MsgCreateStrategy, currentPrice math.LegacyDec) ([]types.ChainStep, error) {
	if currentPrice.IsZero() {
		return nil, fmt.Errorf("current price is zero")
	}
	stopLossPct, takeProfitPct, _ := types.GetRiskParams(msg.RiskLevel)

	anchorBudget, mainBudget := splitBudget10_90(msg.TotalBudget)

	anchorPrice := currentPrice.Mul(math.LegacyOneDec().Add(math.LegacyMustNewDecFromStr("0.01")))
	mainBuyPrice := currentPrice.Mul(math.LegacyOneDec().Add(math.LegacyMustNewDecFromStr("0.005")))

	anchorBody := types.LimitBuyIntent{
		InputDenom:      msg.InputDenom,
		InputAmount:     anchorBudget,
		OutputDenom:     msg.OutputDenom,
		TargetPrice:     anchorPrice,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	anchorBz, _ := json.Marshal(anchorBody)
	step0 := types.ChainStep{
		IntentType:  types.IntentTypeLimitBuy,
		Body:        anchorBz,
		Condition:   types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{Type: types.InputFromFixed, Amount: anchorBudget, Denom: msg.InputDenom},
	}

	mainBody := types.LimitBuyIntent{
		InputDenom:      msg.InputDenom,
		InputAmount:     mainBudget,
		OutputDenom:     msg.OutputDenom,
		TargetPrice:     mainBuyPrice,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	mainBz, _ := json.Marshal(mainBody)
	step1 := types.ChainStep{
		IntentType: types.IntentTypeLimitBuy,
		Body:       mainBz,
		Condition: types.ChainCondition{
			Type:                types.ConditionVolatilityAbove,
			PoolID:              msg.PoolID,
			VolatilityThreshold: types.VolatilityBreakoutThreshold,
		},
		InputSource: types.InputSource{Type: types.InputFromFixed, Amount: mainBudget, Denom: msg.InputDenom},
	}

	outputPrice := math.LegacyOneDec().Quo(currentPrice)
	wideTP := outputPrice.Mul(math.LegacyOneDec().Add(takeProfitPct.Mul(math.LegacyNewDec(2))))
	slPrice := outputPrice.Mul(math.LegacyOneDec().Sub(stopLossPct))

	tpBody := types.TakeProfitIntent{
		InputDenom:      msg.OutputDenom,
		InputAmount:     math.OneInt(),
		OutputDenom:     msg.InputDenom,
		TargetPrice:     wideTP,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	tpBz, _ := json.Marshal(tpBody)
	step2 := types.ChainStep{
		IntentType: types.IntentTypeTakeProfit,
		Body:       tpBz,
		Condition:  types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{
			Type:           types.InputFromPortion,
			PortionPercent: math.LegacyMustNewDecFromStr("0.5"),
		},
	}

	slBody := types.StopLossIntent{
		InputDenom:      msg.OutputDenom,
		InputAmount:     math.OneInt(),
		OutputDenom:     msg.InputDenom,
		StopPrice:       slPrice,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	slBz, _ := json.Marshal(slBody)
	step3 := types.ChainStep{
		IntentType:  types.IntentTypeStopLoss,
		Body:        slBz,
		Condition:   types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{Type: types.InputFromPrevious},
	}

	return []types.ChainStep{step0, step1, step2, step3}, nil
}

// trend_confluence:
//
// Requires BOTH bullish sentiment AND a BUY composite signal to fully enter.
// Since chain conditions are a single predicate per step, we cascade: step 1
// waits for F&G >= NeutralUpper, then step 2 waits for the composite signal
// to fire "BUY" before the take-profit path activates.
//
//	Step 0: anchor limit_buy (10%), immediate
//	Step 1: main limit_buy (90%) when F&G >= NeutralUpper (50)
//	Step 2: take_profit (portion 50%) gated on SignalEquals "BUY" — trend confirmation
//	Step 3: take_profit wider (remaining) immediate
//	Step 4: stop_loss safety net immediate
func (k *Keeper) generateTrendConfluence(msg *types.MsgCreateStrategy, currentPrice math.LegacyDec) ([]types.ChainStep, error) {
	if currentPrice.IsZero() {
		return nil, fmt.Errorf("current price is zero")
	}
	stopLossPct, takeProfitPct, _ := types.GetRiskParams(msg.RiskLevel)

	anchorBudget, mainBudget := splitBudget10_90(msg.TotalBudget)

	buyPrice := currentPrice.Mul(math.LegacyOneDec().Add(math.LegacyMustNewDecFromStr("0.005")))

	anchorBody := types.LimitBuyIntent{
		InputDenom:      msg.InputDenom,
		InputAmount:     anchorBudget,
		OutputDenom:     msg.OutputDenom,
		TargetPrice:     buyPrice,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	anchorBz, _ := json.Marshal(anchorBody)
	step0 := types.ChainStep{
		IntentType:  types.IntentTypeLimitBuy,
		Body:        anchorBz,
		Condition:   types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{Type: types.InputFromFixed, Amount: anchorBudget, Denom: msg.InputDenom},
	}

	mainBody := types.LimitBuyIntent{
		InputDenom:      msg.InputDenom,
		InputAmount:     mainBudget,
		OutputDenom:     msg.OutputDenom,
		TargetPrice:     buyPrice,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	mainBz, _ := json.Marshal(mainBody)
	step1 := types.ChainStep{
		IntentType: types.IntentTypeLimitBuy,
		Body:       mainBz,
		Condition: types.ChainCondition{
			Type:               types.ConditionFearGreedAbove,
			PoolID:             msg.PoolID,
			FearGreedThreshold: types.FearGreedNeutralUpper,
		},
		InputSource: types.InputSource{Type: types.InputFromFixed, Amount: mainBudget, Denom: msg.InputDenom},
	}

	outputPrice := math.LegacyOneDec().Quo(currentPrice)
	tpTarget := outputPrice.Mul(math.LegacyOneDec().Add(takeProfitPct))
	slTarget := outputPrice.Mul(math.LegacyOneDec().Sub(stopLossPct))

	tpBody := types.TakeProfitIntent{
		InputDenom:      msg.OutputDenom,
		InputAmount:     math.OneInt(),
		OutputDenom:     msg.InputDenom,
		TargetPrice:     tpTarget,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	tpBz, _ := json.Marshal(tpBody)
	step2 := types.ChainStep{
		IntentType: types.IntentTypeTakeProfit,
		Body:       tpBz,
		Condition: types.ChainCondition{
			Type:        types.ConditionSignalEquals,
			PoolID:      msg.PoolID,
			SignalLabel: "BUY",
		},
		InputSource: types.InputSource{
			Type:           types.InputFromPortion,
			PortionPercent: math.LegacyMustNewDecFromStr("0.5"),
		},
	}

	wideTP := outputPrice.Mul(math.LegacyOneDec().Add(takeProfitPct.Mul(math.LegacyMustNewDecFromStr("1.5"))))
	tp2Body := types.TakeProfitIntent{
		InputDenom:      msg.OutputDenom,
		InputAmount:     math.OneInt(),
		OutputDenom:     msg.InputDenom,
		TargetPrice:     wideTP,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	tp2Bz, _ := json.Marshal(tp2Body)
	step3 := types.ChainStep{
		IntentType:  types.IntentTypeTakeProfit,
		Body:        tp2Bz,
		Condition:   types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{Type: types.InputFromPrevious},
	}

	slBody := types.StopLossIntent{
		InputDenom:      msg.OutputDenom,
		InputAmount:     math.OneInt(),
		OutputDenom:     msg.InputDenom,
		StopPrice:       slTarget,
		PoolID:          msg.PoolID,
		MinOutputAmount: math.ZeroInt(),
	}
	slBz, _ := json.Marshal(slBody)
	step4 := types.ChainStep{
		IntentType:  types.IntentTypeStopLoss,
		Body:        slBz,
		Condition:   types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{Type: types.InputFromPrevious},
	}

	return []types.ChainStep{step0, step1, step2, step3, step4}, nil
}
