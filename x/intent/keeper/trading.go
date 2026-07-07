package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/intent/types"
)

// DefaultDCASlippageBps is the conservative default maximum slippage (in basis
// points, 1 bp = 0.01%) applied to each DCA/TWAP tranche. The DCA intent body
// carries no user-configured slippage field, so per-tranche swaps derive their
// minimum acceptable output from the current pool spot price reduced by this
// tolerance. This prevents tranche swaps from executing with zero slippage
// protection (minAmountOut = 0), which would leave them open to sandwiching.
// LIMITATION: because DCAIntent has no slippage field, this default cannot be
// overridden per-intent; if a user needs tighter/looser bounds, the intent body
// schema must be extended with a MaxSlippage field.
const DefaultDCASlippageBps = 100 // 1.00%

// MaxIntentsPerBlock caps how many trading intents BeginBlock will process
// in a single block. Unprocessed intents remain pending and will be considered
// in subsequent blocks. Prevents BeginBlock from blowing past block-time budgets
// when the pending-intent set grows large.
const MaxIntentsPerBlock = 1000

// ExecuteTradingIntents checks all pending trading intents against current DEX prices
// and executes swaps when price conditions are met. Called in BeginBlock.
func (k *Keeper) ExecuteTradingIntents(ctx context.Context) {
	if k.dexKeeper == nil {
		return
	}

	tradingTypes := []string{
		types.IntentTypeLimitBuy,
		types.IntentTypeLimitSell,
		types.IntentTypeStopLoss,
		types.IntentTypeTakeProfit,
		types.IntentTypeTWAP,
		types.IntentTypeDCA,
		types.IntentTypeCrossChainSwap,
	}

	processed := 0
outer:
	for _, intentType := range tradingTypes {
		intents := k.GetIntentsByType(ctx, intentType)
		for _, intent := range intents {
			if processed >= MaxIntentsPerBlock {
				k.Logger(ctx).Info("trading intent per-block cap reached; remaining intents deferred",
					"cap", MaxIntentsPerBlock)
				break outer
			}
			if err := k.tryExecuteTradingIntent(ctx, intent); err != nil {
				k.Logger(ctx).Debug("trading intent not triggered",
					"intent_id", intent.ID, "type", intent.IntentType, "reason", err)
			}
			processed++
		}
	}
}

// tryExecuteTradingIntent checks if a trading intent's price condition is met and executes it.
func (k *Keeper) tryExecuteTradingIntent(ctx context.Context, intent types.Intent) error {
	switch intent.IntentType {
	case types.IntentTypeLimitBuy:
		return k.tryLimitBuy(ctx, intent)
	case types.IntentTypeLimitSell:
		return k.tryLimitSell(ctx, intent)
	case types.IntentTypeStopLoss:
		return k.tryStopLoss(ctx, intent)
	case types.IntentTypeTakeProfit:
		return k.tryTakeProfit(ctx, intent)
	case types.IntentTypeDCA:
		return k.tryDCA(ctx, intent)
	case types.IntentTypeTWAP:
		return k.tryTWAP(ctx, intent)
	case types.IntentTypeCrossChainSwap:
		return k.tryCrossChainSwap(ctx, intent)
	default:
		return fmt.Errorf("unknown trading intent type: %s", intent.IntentType)
	}
}

// tryLimitBuy executes a limit buy when the spot price drops to or below TargetPrice.
// TargetPrice is denominated as InputDenom per OutputDenom (cost per unit of output).
func (k *Keeper) tryLimitBuy(ctx context.Context, intent types.Intent) error {
	var body types.LimitBuyIntent
	if err := json.Unmarshal(intent.Body, &body); err != nil {
		return fmt.Errorf("invalid limit_buy body: %w", err)
	}

	// Spot price: how much InputDenom per 1 OutputDenom
	spotPrice, err := k.dexKeeper.GetSpotPrice(ctx, body.PoolID, body.OutputDenom, body.InputDenom)
	if err != nil {
		return fmt.Errorf("price fetch failed: %w", err)
	}

	// Buy when price <= target (cheaper than or equal to what user wants to pay)
	if spotPrice.GT(body.TargetPrice) {
		return fmt.Errorf("price %s > target %s", spotPrice, body.TargetPrice)
	}

	return k.executeTradingSwap(ctx, intent, body.PoolID, body.InputDenom, body.InputAmount, body.MinOutputAmount)
}

// tryLimitSell executes a limit sell when the spot price rises to or above TargetPrice.
// TargetPrice is denominated as OutputDenom per InputDenom (revenue per unit sold).
func (k *Keeper) tryLimitSell(ctx context.Context, intent types.Intent) error {
	var body types.LimitSellIntent
	if err := json.Unmarshal(intent.Body, &body); err != nil {
		return fmt.Errorf("invalid limit_sell body: %w", err)
	}

	// Spot price: how much OutputDenom per 1 InputDenom
	spotPrice, err := k.dexKeeper.GetSpotPrice(ctx, body.PoolID, body.InputDenom, body.OutputDenom)
	if err != nil {
		return fmt.Errorf("price fetch failed: %w", err)
	}

	// Sell when price >= target (getting at least as much as user wants)
	if spotPrice.LT(body.TargetPrice) {
		return fmt.Errorf("price %s < target %s", spotPrice, body.TargetPrice)
	}

	return k.executeTradingSwap(ctx, intent, body.PoolID, body.InputDenom, body.InputAmount, body.MinOutputAmount)
}

// tryStopLoss executes a stop-loss sell when the spot price drops to or below StopPrice.
// StopPrice is denominated as OutputDenom per InputDenom.
func (k *Keeper) tryStopLoss(ctx context.Context, intent types.Intent) error {
	var body types.StopLossIntent
	if err := json.Unmarshal(intent.Body, &body); err != nil {
		return fmt.Errorf("invalid stop_loss body: %w", err)
	}

	// Spot price: how much OutputDenom per 1 InputDenom
	spotPrice, err := k.dexKeeper.GetSpotPrice(ctx, body.PoolID, body.InputDenom, body.OutputDenom)
	if err != nil {
		return fmt.Errorf("price fetch failed: %w", err)
	}

	// Trigger stop-loss when price <= stop price (asset losing value)
	if spotPrice.GT(body.StopPrice) {
		return fmt.Errorf("price %s > stop %s", spotPrice, body.StopPrice)
	}

	return k.executeTradingSwap(ctx, intent, body.PoolID, body.InputDenom, body.InputAmount, body.MinOutputAmount)
}

// tryTakeProfit executes a take-profit sell when the spot price rises to or above TargetPrice.
// TargetPrice is denominated as OutputDenom per InputDenom.
func (k *Keeper) tryTakeProfit(ctx context.Context, intent types.Intent) error {
	var body types.TakeProfitIntent
	if err := json.Unmarshal(intent.Body, &body); err != nil {
		return fmt.Errorf("invalid take_profit body: %w", err)
	}

	// Spot price: how much OutputDenom per 1 InputDenom
	spotPrice, err := k.dexKeeper.GetSpotPrice(ctx, body.PoolID, body.InputDenom, body.OutputDenom)
	if err != nil {
		return fmt.Errorf("price fetch failed: %w", err)
	}

	// Trigger take-profit when price >= target (asset reached profit goal)
	if spotPrice.LT(body.TargetPrice) {
		return fmt.Errorf("price %s < target %s", spotPrice, body.TargetPrice)
	}

	return k.executeTradingSwap(ctx, intent, body.PoolID, body.InputDenom, body.InputAmount, body.MinOutputAmount)
}

// tryDCA executes a DCA (dollar-cost averaging) swap if the next execution block has arrived.
func (k *Keeper) tryDCA(ctx context.Context, intent types.Intent) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	var body types.DCAIntent
	if err := json.Unmarshal(intent.Body, &body); err != nil {
		return fmt.Errorf("invalid dca body: %w", err)
	}

	// Initialize NextExecBlock on first check
	if body.NextExecBlock == 0 {
		body.NextExecBlock = intent.CreatedAt
	}

	// Not time yet
	if sdkCtx.BlockHeight() < body.NextExecBlock {
		return fmt.Errorf("next exec at block %d, current %d", body.NextExecBlock, sdkCtx.BlockHeight())
	}

	// All executions done
	if body.ExecutedCount >= body.NumExecutions {
		return fmt.Errorf("all %d DCA executions completed", body.NumExecutions)
	}

	// Calculate per-execution amount
	remaining := body.NumExecutions - body.ExecutedCount
	perExecAmount := body.TotalAmount.QuoRaw(int64(remaining))
	if perExecAmount.IsZero() {
		return fmt.Errorf("per-execution amount is zero")
	}

	// Execute the swap using the intent module account (funds were locked on submit).
	// Wrap the entire swap+delivery sequence in a cache context so that if anything
	// fails mid-sequence (e.g. output delivery rejected, slippage, insufficient liquidity),
	// the input transfer is rolled back and the locked funds stay in the intent module.
	// Without this, a partial failure drains funds block-by-block into the DEX module.
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	tokenIn := sdk.NewCoin(body.InputDenom, perExecAmount)

	// Derive a per-tranche minimum output from the current pool spot price so the
	// swap is protected against sandwiching. GetSpotPrice returns OutputDenom per
	// 1 InputDenom, so the quoted output for this tranche is spotPrice * perExecAmount.
	// We require at least (1 - DefaultDCASlippageBps/10000) of that quote.
	spotPrice, err := k.dexKeeper.GetSpotPrice(ctx, body.PoolID, body.InputDenom, body.OutputDenom)
	if err != nil {
		return fmt.Errorf("DCA price fetch failed (needed for slippage protection): %w", err)
	}
	quotedOut := spotPrice.MulInt(perExecAmount)
	slippageFactor := math.LegacyNewDec(10000 - DefaultDCASlippageBps).QuoInt64(10000)
	minOut := quotedOut.Mul(slippageFactor).TruncateInt()

	cacheCtx, write := sdkCtx.CacheContext()

	tokenOut, err := k.dexKeeper.Swap(cacheCtx, moduleAddr.String(), body.PoolID, tokenIn, minOut)
	if err != nil {
		return fmt.Errorf("DCA swap failed: %w", err)
	}

	// Send output tokens to intent creator
	creatorAddr, err := sdk.AccAddressFromBech32(intent.Creator)
	if err != nil {
		return fmt.Errorf("invalid creator: %w", err)
	}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(cacheCtx, types.ModuleName, creatorAddr, sdk.NewCoins(tokenOut)); err != nil {
		return fmt.Errorf("failed to send DCA output to creator: %w", err)
	}

	// Commit the atomic swap + delivery to parent context
	write()

	// Update DCA state
	body.ExecutedCount++
	body.TotalAmount = body.TotalAmount.Sub(perExecAmount)
	body.NextExecBlock = sdkCtx.BlockHeight() + int64(body.IntervalBlocks)

	updatedBody, _ := json.Marshal(body)
	intent.Body = updatedBody

	if body.ExecutedCount >= body.NumExecutions {
		// All done — mark fulfilled, refund remaining locked fees
		intent.Status = types.StatusFulfilled
		k.refundIntentFees(ctx, intent)
	}

	k.SetIntent(ctx, intent)

	k.Logger(ctx).Info("DCA execution",
		"intent_id", intent.ID,
		"execution", body.ExecutedCount,
		"of", body.NumExecutions,
		"swapped", tokenIn,
		"received", tokenOut,
	)

	return nil
}

// tryTWAP executes a TWAP (time-weighted average price) swap — same as DCA but named differently.
// TWAP splits a large order across blocks to minimize price impact.
func (k *Keeper) tryTWAP(ctx context.Context, intent types.Intent) error {
	// TWAP reuses DCA mechanics — split amount over N blocks at regular intervals
	return k.tryDCA(ctx, intent)
}

// executeTradingSwap performs the actual swap for a trading intent that has been triggered.
// It swaps from the intent module account (where funds were locked) and sends output to the creator.
// The swap + delivery are wrapped in a cache context so a partial failure (e.g. blocked
// recipient, slippage, insufficient liquidity) rolls the input transfer back atomically.
func (k *Keeper) executeTradingSwap(ctx context.Context, intent types.Intent, poolID uint64, inputDenom string, inputAmount, minOutput math.Int) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)

	// This runs inside ExecuteTradingIntents in BeginBlock. A trading intent body
	// with a nil/non-positive input_amount would make sdk.NewCoin below PANIC
	// (nil amount fails Coin.Validate) — halting the chain. Fail the intent
	// instead of trusting the stored body.
	if inputAmount.IsNil() || !inputAmount.IsPositive() {
		intent.Status = types.StatusFailed
		k.SetIntent(ctx, intent)
		k.refundIntentFees(ctx, intent)
		k.FailChainIfLinked(ctx, intent)
		return fmt.Errorf("trading intent has invalid input amount")
	}

	tokenIn := sdk.NewCoin(inputDenom, inputAmount)

	// Atomic swap + delivery via cache context
	cacheCtx, write := sdkCtx.CacheContext()

	// Execute swap from module account
	tokenOut, err := k.dexKeeper.Swap(cacheCtx, moduleAddr.String(), poolID, tokenIn, minOutput)
	if err != nil {
		// Nothing was committed, but still mark intent as failed and refund fees
		intent.Status = types.StatusFailed
		k.SetIntent(ctx, intent)
		k.refundIntentFees(ctx, intent)
		k.FailChainIfLinked(ctx, intent)
		return fmt.Errorf("swap execution failed: %w", err)
	}

	outputCoins := sdk.NewCoins(tokenOut)

	// Check if this intent is an intermediate chain step — if so, keep output in module account
	isIntermediate := k.IsIntermediateChainStep(ctx, intent.ID)

	if !isIntermediate {
		// Send output tokens to intent creator
		creatorAddr, err := sdk.AccAddressFromBech32(intent.Creator)
		if err != nil {
			return fmt.Errorf("invalid creator: %w", err)
		}

		if err := k.bankKeeper.SendCoinsFromModuleToAccount(cacheCtx, types.ModuleName, creatorAddr, outputCoins); err != nil {
			return fmt.Errorf("failed to send swap output: %w", err)
		}
	}
	// If intermediate, output stays in module account for the next step

	// Commit the atomic swap + delivery to parent context
	write()

	// Mark intent as fulfilled and refund fees (chain-linked intents have zero fees)
	intent.Status = types.StatusFulfilled
	k.SetIntent(ctx, intent)
	k.refundIntentFees(ctx, intent)

	k.Logger(ctx).Info("trading intent executed",
		"intent_id", intent.ID,
		"type", intent.IntentType,
		"creator", intent.Creator,
		"token_in", tokenIn,
		"token_out", tokenOut,
		"chain_linked", isIntermediate,
	)

	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"trading_intent_executed",
		sdk.NewAttribute("intent_id", intent.ID),
		sdk.NewAttribute("intent_type", intent.IntentType),
		sdk.NewAttribute("creator", intent.Creator),
		sdk.NewAttribute("token_in", tokenIn.String()),
		sdk.NewAttribute("token_out", tokenOut.String()),
	))

	// Advance the chain if this intent is part of one
	k.AdvanceChainIfLinked(ctx, intent, outputCoins)

	return nil
}

// refundIntentFees refunds the locked MaxFee and Tip back to the intent creator.
func (k *Keeper) refundIntentFees(ctx context.Context, intent types.Intent) {
	creatorAddr, err := sdk.AccAddressFromBech32(intent.Creator)
	if err != nil {
		return
	}

	totalRefund := intent.MaxFee.Add(intent.Tip...)
	if totalRefund.IsAllPositive() {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, totalRefund); err != nil {
			k.Logger(ctx).Error("failed to refund trading intent fees",
				"intent_id", intent.ID, "error", err)
		}
	}
}

// refundTradingTokens refunds locked input tokens for a trading intent that was not executed.
// Called when a trading intent expires or fails before execution.
func (k *Keeper) refundTradingTokens(ctx context.Context, intent types.Intent) {
	coins := k.extractTradingInputCoins(intent.IntentType, intent.Body)
	if !coins.IsAllPositive() {
		return
	}

	creatorAddr, err := sdk.AccAddressFromBech32(intent.Creator)
	if err != nil {
		return
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, coins); err != nil {
		k.Logger(ctx).Error("failed to refund trading input tokens",
			"intent_id", intent.ID, "error", err)
	}
}

// extractTradingInputCoins parses the intent body and returns the input coins that need to be locked.
func (k *Keeper) extractTradingInputCoins(intentType string, body json.RawMessage) sdk.Coins {
	switch intentType {
	case types.IntentTypeLimitBuy:
		var b types.LimitBuyIntent
		if err := json.Unmarshal(body, &b); err != nil {
			return nil
		}
		return sdk.NewCoins(sdk.NewCoin(b.InputDenom, b.InputAmount))

	case types.IntentTypeLimitSell:
		var b types.LimitSellIntent
		if err := json.Unmarshal(body, &b); err != nil {
			return nil
		}
		return sdk.NewCoins(sdk.NewCoin(b.InputDenom, b.InputAmount))

	case types.IntentTypeStopLoss:
		var b types.StopLossIntent
		if err := json.Unmarshal(body, &b); err != nil {
			return nil
		}
		return sdk.NewCoins(sdk.NewCoin(b.InputDenom, b.InputAmount))

	case types.IntentTypeTakeProfit:
		var b types.TakeProfitIntent
		if err := json.Unmarshal(body, &b); err != nil {
			return nil
		}
		return sdk.NewCoins(sdk.NewCoin(b.InputDenom, b.InputAmount))

	case types.IntentTypeDCA:
		var b types.DCAIntent
		if err := json.Unmarshal(body, &b); err != nil {
			return nil
		}
		return sdk.NewCoins(sdk.NewCoin(b.InputDenom, b.TotalAmount))

	case types.IntentTypeTWAP:
		var b types.DCAIntent // TWAP uses same body as DCA
		if err := json.Unmarshal(body, &b); err != nil {
			return nil
		}
		return sdk.NewCoins(sdk.NewCoin(b.InputDenom, b.TotalAmount))

	case types.IntentTypeCrossChainSwap:
		return extractCrossChainInputCoins(body)

	default:
		return nil
	}
}
