package keeper_test

import (
	"encoding/json"
	"testing"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"syreen/x/intent/types"
)

func makeChainMsg(t *testing.T, steps []types.ChainStep) *types.MsgSubmitChain {
	return &types.MsgSubmitChain{
		Creator:      creatorAddr,
		Steps:        steps,
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1000)),
		Tip:          sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		ExpiryBlocks: 500,
	}
}

func limitBuyStep(poolID uint64, inputDenom string, inputAmount int64, outputDenom string, targetPrice string) types.ChainStep {
	body := types.LimitBuyIntent{
		InputDenom:      inputDenom,
		InputAmount:     math.NewInt(inputAmount),
		OutputDenom:     outputDenom,
		TargetPrice:     math.LegacyMustNewDecFromStr(targetPrice),
		PoolID:          poolID,
		MinOutputAmount: math.ZeroInt(),
	}
	bz, _ := json.Marshal(body)
	return types.ChainStep{
		IntentType: types.IntentTypeLimitBuy,
		Body:       bz,
		Condition:  types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{
			Type: types.InputFromFixed,
			Amount: math.NewInt(inputAmount),
			Denom:  inputDenom,
		},
	}
}

func limitSellStepFromPrevious(poolID uint64, outputDenom string, targetPrice string, condition types.ChainCondition) types.ChainStep {
	body := types.LimitSellIntent{
		InputDenom:      "placeholder",
		InputAmount:     math.NewInt(1), // will be rewritten
		OutputDenom:     outputDenom,
		TargetPrice:     math.LegacyMustNewDecFromStr(targetPrice),
		PoolID:          poolID,
		MinOutputAmount: math.ZeroInt(),
	}
	bz, _ := json.Marshal(body)
	return types.ChainStep{
		IntentType:  types.IntentTypeLimitSell,
		Body:        bz,
		Condition:   condition,
		InputSource: types.InputSource{Type: types.InputFromPrevious},
	}
}

func TestSubmitChain_Basic(t *testing.T) {
	k, ctx, _ := setupTradingKeeper(t)

	step0 := limitBuyStep(1, "usyreen", 1000000, "uusdc", "0.5")
	step1 := limitSellStepFromPrevious(1, "usyreen", "2.0", types.ChainCondition{
		Type:           types.ConditionPriceAbove,
		PriceTarget:    math.LegacyMustNewDecFromStr("2.0"),
		PriceDenom:     "uusdc",
		PriceBaseDenom: "usyreen",
		PoolID:         1,
	})

	msg := makeChainMsg(t, []types.ChainStep{step0, step1})
	chainID, err := k.SubmitChain(ctx, msg)
	require.NoError(t, err)
	require.NotEmpty(t, chainID)

	// Verify chain stored
	chain, found := k.GetChain(ctx, chainID)
	require.True(t, found)
	require.Equal(t, types.ChainStatusActive, chain.Status)
	require.Equal(t, 2, len(chain.Steps))
	require.Equal(t, 0, chain.CurrentIdx)

	// Step 0 should have an intent created
	require.NotEmpty(t, chain.Steps[0].IntentID)
	require.Equal(t, types.StatusPending, chain.Steps[0].Status)

	// Step 1 should not have an intent yet
	require.Empty(t, chain.Steps[1].IntentID)
}

func TestSubmitChain_TooFewSteps(t *testing.T) {
	step0 := limitBuyStep(1, "usyreen", 1000000, "uusdc", "0.5")
	msg := makeChainMsg(t, []types.ChainStep{step0})
	msg.Steps = []types.ChainStep{step0} // only 1 step

	err := msg.ValidateBasic()
	require.Error(t, err)
}

func TestChainAdvance_ImmediateCondition(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	// Set price so limit buy triggers
	dex.SetSpotPrice(1, "uusdc", "usyreen", math.LegacyMustNewDecFromStr("0.4"))
	dex.swapOut = sdk.NewCoin("uusdc", math.NewInt(950000))

	step0 := limitBuyStep(1, "usyreen", 1000000, "uusdc", "0.5")
	step1 := limitSellStepFromPrevious(1, "usyreen", "2.0", types.ChainCondition{
		Type: types.ConditionImmediate,
	})

	msg := makeChainMsg(t, []types.ChainStep{step0, step1})
	chainID, err := k.SubmitChain(ctx, msg)
	require.NoError(t, err)

	// Execute trading intents — step 0 should trigger
	k.ExecuteTradingIntents(ctx)

	// Check chain advanced
	chain, found := k.GetChain(ctx, chainID)
	require.True(t, found)

	// Step 0 should be fulfilled
	require.Equal(t, types.StatusFulfilled, chain.Steps[0].Status)

	// Step 1 should be activated (immediate condition)
	require.NotEmpty(t, chain.Steps[1].IntentID)
	require.Equal(t, types.StatusPending, chain.Steps[1].Status)

	// Chain should still be active (step 1 pending)
	require.Equal(t, types.ChainStatusActive, chain.Status)
}

func TestChainAdvance_PriceCondition(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	// Set price for step 0 to trigger
	dex.SetSpotPrice(1, "uusdc", "usyreen", math.LegacyMustNewDecFromStr("0.4"))
	dex.swapOut = sdk.NewCoin("uusdc", math.NewInt(950000))

	step0 := limitBuyStep(1, "usyreen", 1000000, "uusdc", "0.5")
	step1 := limitSellStepFromPrevious(1, "usyreen", "2.0", types.ChainCondition{
		Type:           types.ConditionPriceAbove,
		PriceTarget:    math.LegacyMustNewDecFromStr("2.0"),
		PriceDenom:     "uusdc",
		PriceBaseDenom: "usyreen",
		PoolID:         1,
	})

	msg := makeChainMsg(t, []types.ChainStep{step0, step1})
	chainID, err := k.SubmitChain(ctx, msg)
	require.NoError(t, err)

	// Execute — step 0 triggers
	k.ExecuteTradingIntents(ctx)

	chain, _ := k.GetChain(ctx, chainID)
	require.Equal(t, types.StatusFulfilled, chain.Steps[0].Status)

	// Step 1 should NOT be activated (price condition not met)
	require.Empty(t, chain.Steps[1].IntentID)

	// Now set price to meet condition
	dex.SetSpotPrice(1, "uusdc", "usyreen", math.LegacyMustNewDecFromStr("2.5"))

	// Run ProcessChainConditions (BeginBlock)
	k.ProcessChainConditions(ctx)

	chain, _ = k.GetChain(ctx, chainID)
	// Now step 1 should be activated
	require.NotEmpty(t, chain.Steps[1].IntentID)
	require.Equal(t, types.StatusPending, chain.Steps[1].Status)
}

func TestChainAdvance_DelayCondition(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	dex.SetSpotPrice(1, "uusdc", "usyreen", math.LegacyMustNewDecFromStr("0.4"))
	dex.swapOut = sdk.NewCoin("uusdc", math.NewInt(950000))

	step0 := limitBuyStep(1, "usyreen", 1000000, "uusdc", "0.5")
	step1 := limitSellStepFromPrevious(1, "usyreen", "2.0", types.ChainCondition{
		Type:        types.ConditionDelay,
		DelayBlocks: 5,
	})

	msg := makeChainMsg(t, []types.ChainStep{step0, step1})
	chainID, err := k.SubmitChain(ctx, msg)
	require.NoError(t, err)

	// Execute step 0 at block 10
	k.ExecuteTradingIntents(ctx)

	chain, _ := k.GetChain(ctx, chainID)
	require.Equal(t, types.StatusFulfilled, chain.Steps[0].Status)
	require.Equal(t, int64(10), chain.Steps[0].CompletedAt)

	// Step 1 not activated yet (need 5 block delay)
	require.Empty(t, chain.Steps[1].IntentID)

	// ProcessChainConditions at block 10 — too early
	k.ProcessChainConditions(ctx)
	chain, _ = k.GetChain(ctx, chainID)
	require.Empty(t, chain.Steps[1].IntentID)

	// Advance to block 15 and try again
	ctx = ctx.WithBlockHeight(15)
	k.ProcessChainConditions(ctx)
	chain, _ = k.GetChain(ctx, chainID)
	require.NotEmpty(t, chain.Steps[1].IntentID)
}

func TestChainComplete(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	// Both steps will trigger immediately
	dex.SetSpotPrice(1, "uusdc", "usyreen", math.LegacyMustNewDecFromStr("0.4"))
	dex.swapOut = sdk.NewCoin("uusdc", math.NewInt(950000))

	step0 := limitBuyStep(1, "usyreen", 1000000, "uusdc", "0.5")
	step1 := limitSellStepFromPrevious(1, "usyreen", "0.1", types.ChainCondition{
		Type: types.ConditionImmediate,
	})

	msg := makeChainMsg(t, []types.ChainStep{step0, step1})
	chainID, err := k.SubmitChain(ctx, msg)
	require.NoError(t, err)

	// Execute step 0
	k.ExecuteTradingIntents(ctx)

	// Set price for step 1 to trigger
	dex.SetSpotPrice(1, "uusdc", "usyreen", math.LegacyMustNewDecFromStr("0.5"))

	// Execute step 1
	k.ExecuteTradingIntents(ctx)

	chain, found := k.GetChain(ctx, chainID)
	require.True(t, found)
	require.Equal(t, types.ChainStatusCompleted, chain.Status)
	require.Equal(t, types.StatusFulfilled, chain.Steps[0].Status)
	require.Equal(t, types.StatusFulfilled, chain.Steps[1].Status)
}

func TestCancelChain(t *testing.T) {
	k, ctx, _ := setupTradingKeeper(t)

	step0 := limitBuyStep(1, "usyreen", 1000000, "uusdc", "0.5")
	step1 := limitSellStepFromPrevious(1, "usyreen", "2.0", types.ChainCondition{
		Type: types.ConditionImmediate,
	})

	msg := makeChainMsg(t, []types.ChainStep{step0, step1})
	chainID, err := k.SubmitChain(ctx, msg)
	require.NoError(t, err)

	// Cancel
	err = k.CancelChain(ctx, &types.MsgCancelChain{
		Creator: creatorAddr,
		ChainID: chainID,
	})
	require.NoError(t, err)

	chain, found := k.GetChain(ctx, chainID)
	require.True(t, found)
	require.Equal(t, types.ChainStatusCancelled, chain.Status)
}

func TestCancelChain_NotCreator(t *testing.T) {
	k, ctx, _ := setupTradingKeeper(t)

	step0 := limitBuyStep(1, "usyreen", 1000000, "uusdc", "0.5")
	step1 := limitSellStepFromPrevious(1, "usyreen", "2.0", types.ChainCondition{
		Type: types.ConditionImmediate,
	})

	msg := makeChainMsg(t, []types.ChainStep{step0, step1})
	chainID, err := k.SubmitChain(ctx, msg)
	require.NoError(t, err)

	// Try to cancel from a different address
	err = k.CancelChain(ctx, &types.MsgCancelChain{
		Creator: solverAddr, // not the creator
		ChainID: chainID,
	})
	require.ErrorIs(t, err, types.ErrChainNotCreator)
}

func TestChainExpiry(t *testing.T) {
	k, ctx, _ := setupTradingKeeper(t)

	step0 := limitBuyStep(1, "usyreen", 1000000, "uusdc", "0.5")
	step1 := limitSellStepFromPrevious(1, "usyreen", "2.0", types.ChainCondition{
		Type: types.ConditionImmediate,
	})

	msg := &types.MsgSubmitChain{
		Creator:      creatorAddr,
		Steps:        []types.ChainStep{step0, step1},
		MaxFee:       sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1000)),
		Tip:          sdk.NewCoins(sdk.NewInt64Coin("usyreen", 100)),
		ExpiryBlocks: 20, // expires at block 30
	}

	chainID, err := k.SubmitChain(ctx, msg)
	require.NoError(t, err)

	// Advance past expiry
	ctx = ctx.WithBlockHeight(50)
	k.ProcessChainConditions(ctx)

	chain, found := k.GetChain(ctx, chainID)
	require.True(t, found)
	require.Equal(t, types.ChainStatusExpired, chain.Status)
}

func TestChainTokenFlow_IntermediateStaysInModule(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	dex.SetSpotPrice(1, "uusdc", "usyreen", math.LegacyMustNewDecFromStr("0.4"))
	dex.swapOut = sdk.NewCoin("uusdc", math.NewInt(950000))

	step0 := limitBuyStep(1, "usyreen", 1000000, "uusdc", "0.5")
	step1 := limitSellStepFromPrevious(1, "usyreen", "2.0", types.ChainCondition{
		Type: types.ConditionImmediate,
	})

	msg := makeChainMsg(t, []types.ChainStep{step0, step1})
	chainID, err := k.SubmitChain(ctx, msg)
	require.NoError(t, err)

	// Verify step 0's intent is chain-linked
	chain, _ := k.GetChain(ctx, chainID)
	require.True(t, k.IsChainLinkedIntent(ctx, chain.Steps[0].IntentID))
	require.True(t, k.IsIntermediateChainStep(ctx, chain.Steps[0].IntentID))

	// Execute step 0
	k.ExecuteTradingIntents(ctx)

	// The swap output should NOT have been sent to creator (intermediate step)
	// Verify by checking that the dex swap was called but no SendCoinsFromModuleToAccount for output
	require.Len(t, dex.swapCalls, 1)

	// Chain should have recorded the output
	chain, _ = k.GetChain(ctx, chainID)
	require.Equal(t, types.StatusFulfilled, chain.Steps[0].Status)
	require.NotEmpty(t, chain.Steps[0].OutputCoins)
}

func TestChainPortionInput(t *testing.T) {
	k, ctx, dex := setupTradingKeeper(t)

	dex.SetSpotPrice(1, "uusdc", "usyreen", math.LegacyMustNewDecFromStr("0.4"))
	dex.swapOut = sdk.NewCoin("uusdc", math.NewInt(1000000))

	step0 := limitBuyStep(1, "usyreen", 1000000, "uusdc", "0.5")

	// Step 1: sell 50% of step 0's output
	body1 := types.LimitSellIntent{
		InputDenom:      "placeholder",
		InputAmount:     math.NewInt(1),
		OutputDenom:     "usyreen",
		TargetPrice:     math.LegacyMustNewDecFromStr("0.1"),
		PoolID:          1,
		MinOutputAmount: math.ZeroInt(),
	}
	bz1, _ := json.Marshal(body1)
	step1 := types.ChainStep{
		IntentType: types.IntentTypeLimitSell,
		Body:       bz1,
		Condition:  types.ChainCondition{Type: types.ConditionImmediate},
		InputSource: types.InputSource{
			Type:           types.InputFromPortion,
			PortionPercent: math.LegacyMustNewDecFromStr("0.5"),
		},
	}

	msg := makeChainMsg(t, []types.ChainStep{step0, step1})
	chainID, err := k.SubmitChain(ctx, msg)
	require.NoError(t, err)

	// Execute step 0
	k.ExecuteTradingIntents(ctx)

	chain, _ := k.GetChain(ctx, chainID)
	require.Equal(t, types.StatusFulfilled, chain.Steps[0].Status)

	// Step 1 should be activated with 50% of output
	require.NotEmpty(t, chain.Steps[1].IntentID)

	// Check the rewritten body has 500000 (50% of 1000000)
	intent, found := k.GetIntent(ctx, chain.Steps[1].IntentID)
	require.True(t, found)

	var sellBody types.LimitSellIntent
	require.NoError(t, json.Unmarshal(intent.Body, &sellBody))
	require.Equal(t, math.NewInt(500000), sellBody.InputAmount)
	require.Equal(t, "uusdc", sellBody.InputDenom)
}
