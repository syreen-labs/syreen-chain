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
// Chain Storage
// ---------------------------------------------------------------------------

func (k Keeper) SetChain(ctx context.Context, chain types.IntentChain) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(chain)
	kvStore.Set(types.ChainKey(chain.ID), bz)
}

func (k Keeper) GetChain(ctx context.Context, chainID string) (types.IntentChain, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.ChainKey(chainID))
	if err != nil || bz == nil {
		return types.IntentChain{}, false
	}
	var chain types.IntentChain
	if err := json.Unmarshal(bz, &chain); err != nil {
		return types.IntentChain{}, false
	}
	return chain, true
}

func (k Keeper) SetIntentChainMap(ctx context.Context, intentID string, chainID string, stepIdx int) {
	kvStore := k.storeService.OpenKVStore(ctx)
	value := chainID + ":" + strconv.Itoa(stepIdx)
	kvStore.Set(types.IntentChainMapKey(intentID), []byte(value))
}

func (k Keeper) GetIntentChainMap(ctx context.Context, intentID string) (chainID string, stepIdx int, found bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.IntentChainMapKey(intentID))
	if err != nil || bz == nil {
		return "", 0, false
	}
	parts := bytes.SplitN(bz, []byte(":"), 2)
	if len(parts) != 2 {
		return "", 0, false
	}
	idx, err := strconv.Atoi(string(parts[1]))
	if err != nil {
		return "", 0, false
	}
	return string(parts[0]), idx, true
}

func (k Keeper) GetActiveChains(ctx context.Context) []types.IntentChain {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.ChainPrefix)

	iter, err := kvStore.Iterator(prefix, nil)
	if err != nil {
		return nil
	}
	defer iter.Close()

	var chains []types.IntentChain
	for ; iter.Valid(); iter.Next() {
		if !bytes.HasPrefix(iter.Key(), prefix) {
			break
		}
		var chain types.IntentChain
		if err := json.Unmarshal(iter.Value(), &chain); err != nil {
			continue
		}
		if chain.Status == types.ChainStatusActive {
			chains = append(chains, chain)
		}
	}

	sort.Slice(chains, func(i, j int) bool {
		return chains[i].ID < chains[j].ID
	})

	return chains
}

func (k Keeper) GetChainsByCreator(ctx context.Context, creator string) []types.IntentChain {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.ChainPrefix)

	iter, err := kvStore.Iterator(prefix, nil)
	if err != nil {
		return nil
	}
	defer iter.Close()

	var chains []types.IntentChain
	for ; iter.Valid(); iter.Next() {
		if !bytes.HasPrefix(iter.Key(), prefix) {
			break
		}
		var chain types.IntentChain
		if err := json.Unmarshal(iter.Value(), &chain); err != nil {
			continue
		}
		if chain.Creator == creator {
			chains = append(chains, chain)
		}
	}

	sort.Slice(chains, func(i, j int) bool {
		return chains[i].ID < chains[j].ID
	})

	return chains
}

func (k Keeper) nextChainID(ctx context.Context) string {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.ChainCounterKey())
	var counter uint64
	if err == nil && bz != nil {
		counter = types.BytesToUint64(bz)
	}
	counter++
	kvStore.Set(types.ChainCounterKey(), types.Uint64ToBytes(counter))
	return "chain-" + strconv.FormatUint(counter, 10)
}

// ---------------------------------------------------------------------------
// Chain Submission
// ---------------------------------------------------------------------------

func (k *Keeper) SubmitChain(ctx context.Context, msg *types.MsgSubmitChain) (string, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(ctx)

	if !params.EnableChains {
		return "", fmt.Errorf("intent chains are currently disabled")
	}

	if uint64(len(msg.Steps)) > params.MaxChainSteps {
		return "", types.ErrChainTooManySteps
	}

	chainID := k.nextChainID(ctx)

	// Cap the expiry to prevent uint64->int64 overflow.
	// maxExpiry ~5.8 days at 500ms blocks.
	const maxExpiry int64 = 1_000_000
	maxPerStep := params.IntentExpiryBlocks * uint64(len(msg.Steps))
	expiryBlocks := msg.ExpiryBlocks
	if expiryBlocks > maxPerStep {
		expiryBlocks = maxPerStep
	}
	if expiryBlocks > uint64(maxExpiry) {
		expiryBlocks = uint64(maxExpiry)
	}
	expiry := sdkCtx.BlockHeight() + int64(expiryBlocks)

	chain := types.IntentChain{
		ID:         chainID,
		Creator:    msg.Creator,
		Steps:      msg.Steps,
		CurrentIdx: 0,
		Status:     types.ChainStatusActive,
		CreatedAt:  sdkCtx.BlockHeight(),
		Expiry:     expiry,
		MaxFee:     msg.MaxFee,
		Tip:        msg.Tip,
	}

	// Lock fees + tip
	creatorAddr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return "", fmt.Errorf("invalid creator: %w", err)
	}

	totalLock := msg.MaxFee.Add(msg.Tip...)
	if totalLock.IsAllPositive() {
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, creatorAddr, types.ModuleName, totalLock); err != nil {
			return "", fmt.Errorf("failed to lock chain fees: %w", err)
		}
	}

	// Lock step 0's input tokens
	step0Coins := k.extractTradingInputCoins(msg.Steps[0].IntentType, msg.Steps[0].Body)
	if step0Coins.IsAllPositive() {
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, creatorAddr, types.ModuleName, step0Coins); err != nil {
			return "", fmt.Errorf("failed to lock step 0 input tokens: %w", err)
		}
	}

	// Lock fixed-input tokens for later steps
	for i := 1; i < len(msg.Steps); i++ {
		if msg.Steps[i].InputSource.Type == types.InputFromFixed {
			fixedCoins := sdk.NewCoins(sdk.NewCoin(msg.Steps[i].InputSource.Denom, msg.Steps[i].InputSource.Amount))
			if fixedCoins.IsAllPositive() {
				if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, creatorAddr, types.ModuleName, fixedCoins); err != nil {
					return "", fmt.Errorf("failed to lock step %d fixed input: %w", i, err)
				}
			}
		}
	}

	k.SetChain(ctx, chain)

	// Activate step 0
	if err := k.activateChainStep(ctx, &chain, 0); err != nil {
		return "", fmt.Errorf("failed to activate first step: %w", err)
	}

	k.SetChain(ctx, chain)

	k.Logger(ctx).Info("intent chain submitted",
		"chain_id", chainID, "creator", msg.Creator, "steps", len(msg.Steps), "expiry", expiry)

	return chainID, nil
}

// activateChainStep creates an intent for the given step and links it to the chain.
func (k *Keeper) activateChainStep(ctx context.Context, chain *types.IntentChain, stepIdx int) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	step := &chain.Steps[stepIdx]

	// Rewrite body if using previous step's output
	if stepIdx > 0 {
		prevStep := chain.Steps[stepIdx-1]
		if err := k.rewriteStepBody(step, prevStep); err != nil {
			return fmt.Errorf("failed to rewrite step %d body: %w", stepIdx, err)
		}
	}

	// Calculate per-step expiry
	remainingBlocks := chain.Expiry - sdkCtx.BlockHeight()
	if remainingBlocks <= 0 {
		return fmt.Errorf("chain has expired")
	}

	// Submit the intent internally (tokens already locked at chain level)
	intentID, err := k.submitChainStepIntent(ctx, chain.Creator, step.IntentType, step.Body, uint64(remainingBlocks))
	if err != nil {
		return err
	}

	step.IntentID = intentID
	step.Status = types.StatusPending
	chain.CurrentIdx = stepIdx

	// Link intent to chain
	k.SetIntentChainMap(ctx, intentID, chain.ID, stepIdx)

	k.Logger(ctx).Info("chain step activated",
		"chain_id", chain.ID, "step", stepIdx, "intent_id", intentID, "type", step.IntentType)

	return nil
}

// submitChainStepIntent creates an intent for a chain step without locking additional tokens.
func (k *Keeper) submitChainStepIntent(ctx context.Context, creator string, intentType string, body json.RawMessage, expiryBlocks uint64) (string, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(ctx)

	if expiryBlocks > params.IntentExpiryBlocks {
		expiryBlocks = params.IntentExpiryBlocks
	}
	expiry := sdkCtx.BlockHeight() + int64(expiryBlocks)

	intentID := k.nextIntentID(ctx)

	intent := types.Intent{
		ID:         intentID,
		Creator:    creator,
		IntentType: intentType,
		Body:       body,
		MaxFee:     sdk.NewCoins(), // fees locked at chain level
		Tip:        sdk.NewCoins(),
		Expiry:     expiry,
		Status:     types.StatusPending,
		CreatedAt:  sdkCtx.BlockHeight(),
	}

	// Tokens are already locked at the chain level — don't lock again
	k.SetIntent(ctx, intent)

	return intentID, nil
}

// rewriteStepBody updates the step's body with the output from the previous step.
func (k *Keeper) rewriteStepBody(step *types.ChainStep, prevStep types.ChainStep) error {
	if len(prevStep.OutputCoins) == 0 {
		return fmt.Errorf("previous step has no output coins")
	}

	// Get the output coin from previous step
	outputCoin := prevStep.OutputCoins[0]

	var inputAmount math.Int
	inputDenom := outputCoin.Denom

	switch step.InputSource.Type {
	case types.InputFromPrevious:
		inputAmount = outputCoin.Amount
	case types.InputFromPortion:
		inputAmount = step.InputSource.PortionPercent.MulInt(outputCoin.Amount).TruncateInt()
		if inputAmount.IsZero() {
			return fmt.Errorf("portion amount is zero")
		}
	case types.InputFromFixed:
		// Fixed input — don't rewrite, body already has the right amount
		return nil
	default:
		return fmt.Errorf("unknown input source type: %s", step.InputSource.Type)
	}

	// Rewrite the body's input_denom and input_amount
	var bodyMap map[string]json.RawMessage
	if err := json.Unmarshal(step.Body, &bodyMap); err != nil {
		return fmt.Errorf("failed to parse step body: %w", err)
	}

	amountBz, _ := json.Marshal(inputAmount.String())
	denomBz, _ := json.Marshal(inputDenom)
	bodyMap["input_denom"] = denomBz
	bodyMap["input_amount"] = amountBz

	// For DCA/TWAP, rewrite total_amount instead of input_amount
	if step.IntentType == types.IntentTypeDCA || step.IntentType == types.IntentTypeTWAP {
		bodyMap["total_amount"] = amountBz
		bodyMap["input_denom"] = denomBz
	}

	newBody, err := json.Marshal(bodyMap)
	if err != nil {
		return fmt.Errorf("failed to serialize rewritten body: %w", err)
	}

	step.Body = newBody
	return nil
}

// ---------------------------------------------------------------------------
// Chain Advancement (called when an intent is fulfilled)
// ---------------------------------------------------------------------------

// AdvanceChainIfLinked checks if a fulfilled intent belongs to a chain and advances it.
func (k *Keeper) AdvanceChainIfLinked(ctx context.Context, intent types.Intent, outputCoins sdk.Coins) {
	chainID, stepIdx, found := k.GetIntentChainMap(ctx, intent.ID)
	if !found {
		return
	}

	chain, found := k.GetChain(ctx, chainID)
	if !found || chain.Status != types.ChainStatusActive {
		return
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Record output and mark step complete
	chain.Steps[stepIdx].OutputCoins = outputCoins
	chain.Steps[stepIdx].Status = types.StatusFulfilled
	chain.Steps[stepIdx].CompletedAt = sdkCtx.BlockHeight()

	// Send portion remainder to creator if applicable
	if stepIdx+1 < len(chain.Steps) {
		nextStep := chain.Steps[stepIdx+1]
		if nextStep.InputSource.Type == types.InputFromPortion && len(outputCoins) > 0 {
			outputCoin := outputCoins[0]
			portionAmount := nextStep.InputSource.PortionPercent.MulInt(outputCoin.Amount).TruncateInt()
			remainder := outputCoin.Amount.Sub(portionAmount)
			if remainder.IsPositive() {
				creatorAddr, err := sdk.AccAddressFromBech32(chain.Creator)
				if err == nil {
					remainderCoins := sdk.NewCoins(sdk.NewCoin(outputCoin.Denom, remainder))
					if sendErr := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, remainderCoins); sendErr != nil {
						k.Logger(ctx).Error("AdvanceChainIfLinked: remainder send failed",
							"chain_id", chainID, "step", stepIdx, "amount", remainderCoins.String(), "error", sendErr)
					}
				}
			}
		}
	}

	// Check if this was the last step
	if stepIdx+1 >= len(chain.Steps) {
		chain.Status = types.ChainStatusCompleted
		k.refundChainFees(ctx, chain)
		// Send final output to creator
		if len(outputCoins) > 0 {
			creatorAddr, err := sdk.AccAddressFromBech32(chain.Creator)
			if err == nil {
				if sendErr := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, outputCoins); sendErr != nil {
					k.Logger(ctx).Error("AdvanceChainIfLinked: final output send failed",
						"chain_id", chainID, "amount", outputCoins.String(), "error", sendErr)
				}
			}
		}
		k.SetChain(ctx, chain)
		k.Logger(ctx).Info("intent chain completed", "chain_id", chainID, "steps", len(chain.Steps))
		return
	}

	// Try to activate next step if condition is met
	nextIdx := stepIdx + 1
	chain.CurrentIdx = nextIdx // advance index so ProcessChainConditions knows which step to check
	if k.checkChainCondition(ctx, chain, nextIdx) {
		if err := k.activateChainStep(ctx, &chain, nextIdx); err != nil {
			chain.Status = types.ChainStatusFailed
			k.refundChainTokens(ctx, chain, nextIdx)
			k.Logger(ctx).Error("failed to activate next chain step", "chain_id", chainID, "step", nextIdx, "error", err)
		}
	}
	// If condition not met, chain stays active — ProcessChainConditions will check in BeginBlock

	k.SetChain(ctx, chain)
}

// FailChainIfLinked marks a chain as failed when one of its intents fails.
func (k *Keeper) FailChainIfLinked(ctx context.Context, intent types.Intent) {
	chainID, _, found := k.GetIntentChainMap(ctx, intent.ID)
	if !found {
		return
	}

	chain, found := k.GetChain(ctx, chainID)
	if !found || chain.Status != types.ChainStatusActive {
		return
	}

	chain.Status = types.ChainStatusFailed
	k.refundChainTokens(ctx, chain, chain.CurrentIdx+1)
	k.refundChainFees(ctx, chain)
	k.SetChain(ctx, chain)

	k.Logger(ctx).Info("intent chain failed", "chain_id", chainID, "failed_step", chain.CurrentIdx)
}

// ---------------------------------------------------------------------------
// BeginBlock: Process Chain Conditions
// ---------------------------------------------------------------------------

// ProcessChainConditions checks dormant chain steps for condition triggers.
func (k *Keeper) ProcessChainConditions(ctx context.Context) {
	if k.dexKeeper == nil {
		return
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	chains := k.GetActiveChains(ctx)

	for _, chain := range chains {
		// Check overall expiry
		if sdkCtx.BlockHeight() > chain.Expiry {
			chain.Status = types.ChainStatusExpired
			k.refundChainTokens(ctx, chain, chain.CurrentIdx)
			k.refundChainFees(ctx, chain)
			k.SetChain(ctx, chain)
			k.Logger(ctx).Info("intent chain expired", "chain_id", chain.ID)
			continue
		}

		currentStep := chain.Steps[chain.CurrentIdx]

		// If current step already has an intent, let the trading engine handle it
		if currentStep.IntentID != "" {
			continue
		}

		// Check if the dormant step's condition is now met
		if k.checkChainCondition(ctx, chain, chain.CurrentIdx) {
			if err := k.activateChainStep(ctx, &chain, chain.CurrentIdx); err != nil {
				chain.Status = types.ChainStatusFailed
				k.refundChainTokens(ctx, chain, chain.CurrentIdx)
				k.refundChainFees(ctx, chain)
				k.Logger(ctx).Error("failed to activate dormant chain step", "chain_id", chain.ID, "step", chain.CurrentIdx, "error", err)
			}
			k.SetChain(ctx, chain)
		}
	}
}

// ---------------------------------------------------------------------------
// Chain Cancellation
// ---------------------------------------------------------------------------

func (k *Keeper) CancelChain(ctx context.Context, msg *types.MsgCancelChain) error {
	chain, found := k.GetChain(ctx, msg.ChainID)
	if !found {
		return types.ErrChainNotFound
	}

	if chain.Creator != msg.Creator {
		return types.ErrChainNotCreator
	}

	if chain.Status != types.ChainStatusActive {
		return types.ErrChainAlreadyComplete
	}

	// Cancel active intent if one exists
	currentStep := chain.Steps[chain.CurrentIdx]
	if currentStep.IntentID != "" {
		intent, found := k.GetIntent(ctx, currentStep.IntentID)
		if found && (intent.Status == types.StatusPending || intent.Status == types.StatusSolving) {
			intent.Status = types.StatusExpired
			k.SetIntent(ctx, intent)
			// Refund intent's locked trading tokens
			k.refundTradingTokens(ctx, intent)
		}
	}

	// Refund remaining fixed-input tokens
	k.refundChainTokens(ctx, chain, chain.CurrentIdx+1)
	k.refundChainFees(ctx, chain)

	chain.Status = types.ChainStatusCancelled
	k.SetChain(ctx, chain)

	k.Logger(ctx).Info("intent chain cancelled", "chain_id", msg.ChainID, "creator", msg.Creator)
	return nil
}

// ---------------------------------------------------------------------------
// Condition Checking
// ---------------------------------------------------------------------------

func (k *Keeper) checkChainCondition(ctx context.Context, chain types.IntentChain, stepIdx int) bool {
	if stepIdx >= len(chain.Steps) {
		return false
	}

	step := chain.Steps[stepIdx]
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	switch step.Condition.Type {
	case types.ConditionImmediate:
		return true

	case types.ConditionPriceAbove:
		if k.dexKeeper == nil {
			return false
		}
		price, err := k.dexKeeper.GetSpotPrice(ctx, step.Condition.PoolID, step.Condition.PriceDenom, step.Condition.PriceBaseDenom)
		if err != nil {
			return false
		}
		return price.GTE(step.Condition.PriceTarget)

	case types.ConditionPriceBelow:
		if k.dexKeeper == nil {
			return false
		}
		price, err := k.dexKeeper.GetSpotPrice(ctx, step.Condition.PoolID, step.Condition.PriceDenom, step.Condition.PriceBaseDenom)
		if err != nil {
			return false
		}
		return price.LTE(step.Condition.PriceTarget)

	case types.ConditionDelay:
		if stepIdx == 0 {
			return true
		}
		prevStep := chain.Steps[stepIdx-1]
		if prevStep.CompletedAt == 0 {
			return false
		}
		return sdkCtx.BlockHeight() >= prevStep.CompletedAt+int64(step.Condition.DelayBlocks)

	case types.ConditionRSIAbove:
		if k.dexKeeper == nil {
			return false
		}
		_, _, rsi := k.dexKeeper.GetSignalForPool(ctx, step.Condition.PoolID)
		if rsi.IsNil() || !rsi.IsPositive() {
			return false
		}
		return rsi.GTE(step.Condition.RSIThreshold)

	case types.ConditionRSIBelow:
		if k.dexKeeper == nil {
			return false
		}
		_, _, rsi := k.dexKeeper.GetSignalForPool(ctx, step.Condition.PoolID)
		if rsi.IsNil() || !rsi.IsPositive() {
			return false
		}
		return rsi.LTE(step.Condition.RSIThreshold)

	case types.ConditionFearGreedAbove:
		if k.dexKeeper == nil {
			return false
		}
		idx, _, ok := k.dexKeeper.GetFearGreedIndex(ctx, step.Condition.PoolID)
		if !ok {
			return false
		}
		return idx >= step.Condition.FearGreedThreshold

	case types.ConditionFearGreedBelow:
		if k.dexKeeper == nil {
			return false
		}
		idx, _, ok := k.dexKeeper.GetFearGreedIndex(ctx, step.Condition.PoolID)
		if !ok {
			return false
		}
		return idx <= step.Condition.FearGreedThreshold

	case types.ConditionSignalEquals:
		if k.dexKeeper == nil {
			return false
		}
		sig, _, _ := k.dexKeeper.GetSignalForPool(ctx, step.Condition.PoolID)
		return sig == step.Condition.SignalLabel

	case types.ConditionVolatilityAbove:
		if k.dexKeeper == nil {
			return false
		}
		vol := k.dexKeeper.GetSignalVolatility(ctx, step.Condition.PoolID)
		if vol.IsNil() {
			return false
		}
		return vol.GTE(step.Condition.VolatilityThreshold)

	default:
		return false
	}
}

// ---------------------------------------------------------------------------
// Refund Helpers
// ---------------------------------------------------------------------------

func (k *Keeper) refundChainFees(ctx context.Context, chain types.IntentChain) {
	creatorAddr, err := sdk.AccAddressFromBech32(chain.Creator)
	if err != nil {
		k.Logger(ctx).Error("refundChainFees: invalid creator", "chain_id", chain.ID, "error", err)
		return
	}
	totalRefund := chain.MaxFee.Add(chain.Tip...)
	if totalRefund.IsAllPositive() {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, totalRefund); err != nil {
			k.Logger(ctx).Error("refundChainFees: bank transfer failed",
				"chain_id", chain.ID, "creator", chain.Creator, "amount", totalRefund.String(), "error", err)
		}
	}
}

// refundChainTokens refunds locked tokens for steps from startIdx onward.
func (k *Keeper) refundChainTokens(ctx context.Context, chain types.IntentChain, startIdx int) {
	creatorAddr, err := sdk.AccAddressFromBech32(chain.Creator)
	if err != nil {
		k.Logger(ctx).Error("refundChainTokens: invalid creator", "chain_id", chain.ID, "error", err)
		return
	}

	for i := startIdx; i < len(chain.Steps); i++ {
		step := chain.Steps[i]
		if step.InputSource.Type == types.InputFromFixed {
			coins := sdk.NewCoins(sdk.NewCoin(step.InputSource.Denom, step.InputSource.Amount))
			if coins.IsAllPositive() {
				if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, coins); err != nil {
					k.Logger(ctx).Error("refundChainTokens: bank transfer failed",
						"chain_id", chain.ID, "step", i, "creator", chain.Creator, "amount", coins.String(), "error", err)
				}
			}
		}
	}
}

// IsChainLinkedIntent returns true if the given intent ID belongs to a chain.
func (k *Keeper) IsChainLinkedIntent(ctx context.Context, intentID string) bool {
	_, _, found := k.GetIntentChainMap(ctx, intentID)
	return found
}

// IsIntermediateChainStep returns true if the intent is a chain step that is NOT the final step.
func (k *Keeper) IsIntermediateChainStep(ctx context.Context, intentID string) bool {
	chainID, stepIdx, found := k.GetIntentChainMap(ctx, intentID)
	if !found {
		return false
	}
	chain, found := k.GetChain(ctx, chainID)
	if !found {
		return false
	}
	return stepIdx < len(chain.Steps)-1
}
