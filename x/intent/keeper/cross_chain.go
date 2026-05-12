package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/intent/types"
)

// tryCrossChainSwap handles cross-chain swap intents: swap locally, then IBC transfer the output.
func (k *Keeper) tryCrossChainSwap(ctx context.Context, intent types.Intent) error {
	if k.dexKeeper == nil {
		return fmt.Errorf("dex keeper not set")
	}
	if k.transferKeeper == nil {
		return fmt.Errorf("transfer keeper not set")
	}

	var body types.CrossChainSwapIntent
	if err := json.Unmarshal(intent.Body, &body); err != nil {
		return fmt.Errorf("invalid cross_chain_swap body: %w", err)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Phase 1: Execute local DEX swap (if not already done)
	if !body.SwapDone {
		moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
		tokenIn := sdk.NewCoin(body.InputDenom, body.InputAmount)

		tokenOut, err := k.dexKeeper.Swap(ctx, moduleAddr.String(), body.PoolID, tokenIn, body.MinOutputAmount)
		if err != nil {
			intent.Status = types.StatusFailed
			k.SetIntent(ctx, intent)
			k.refundIntentFees(ctx, intent)
			k.refundTradingTokens(ctx, intent)
			k.FailChainIfLinked(ctx, intent)
			return fmt.Errorf("cross-chain swap failed: %w", err)
		}

		body.SwapDone = true

		// Phase 2: IBC transfer the output to the destination chain.
		// Use timestamp-based timeout instead of block height, because
		// block heights are meaningless across chains with different block times.
		// Convert TimeoutBlocks to a duration (500ms per Syreen block).
		timeoutDuration := time.Duration(body.TimeoutBlocks) * 500 * time.Millisecond
		if timeoutDuration < 10*time.Minute {
			timeoutDuration = 10 * time.Minute // minimum 10 minute timeout
		}
		timeoutTimestamp := uint64(sdkCtx.BlockTime().UnixNano()) + uint64(timeoutDuration)

		ibcMsg := &types.IBCTransferMsg{
			SourcePort:       body.IBCSourcePort,
			SourceChannel:    body.IBCSourceChannel,
			Token:            tokenOut,
			Sender:           moduleAddr.String(),
			Receiver:         body.Receiver,
			TimeoutHeight:    0, // use timestamp-based timeout, not height
			TimeoutTimestamp: timeoutTimestamp,
		}

		_, err = k.transferKeeper.Transfer(ctx, ibcMsg)
		if err != nil {
			// Swap succeeded but IBC transfer failed — send swap output to creator instead
			creatorAddr, addrErr := sdk.AccAddressFromBech32(intent.Creator)
			if addrErr == nil {
				_ = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, sdk.NewCoins(tokenOut))
			}

			intent.Status = types.StatusFailed
			k.SetIntent(ctx, intent)
			k.refundIntentFees(ctx, intent)
			k.FailChainIfLinked(ctx, intent)
			return fmt.Errorf("IBC transfer failed: %w", err)
		}

		body.IBCSent = true

		// Update intent body with tracking state
		updatedBody, _ := json.Marshal(body)
		intent.Body = updatedBody

		// Mark as fulfilled — IBC transport handles delivery/timeout
		intent.Status = types.StatusFulfilled
		k.SetIntent(ctx, intent)
		k.refundIntentFees(ctx, intent)

		// For chain-linked intents, advance the chain with the IBC-sent amount
		// (the tokens are in flight, but logically the step is complete)
		k.AdvanceChainIfLinked(ctx, intent, sdk.NewCoins(tokenOut))

		k.Logger(ctx).Info("cross-chain swap executed",
			"intent_id", intent.ID,
			"creator", intent.Creator,
			"token_in", tokenIn,
			"token_out", tokenOut,
			"channel", body.IBCSourceChannel,
			"receiver", body.Receiver,
		)

		sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
			"cross_chain_swap_executed",
			sdk.NewAttribute("intent_id", intent.ID),
			sdk.NewAttribute("creator", intent.Creator),
			sdk.NewAttribute("token_in", tokenIn.String()),
			sdk.NewAttribute("token_out", tokenOut.String()),
			sdk.NewAttribute("ibc_channel", body.IBCSourceChannel),
			sdk.NewAttribute("receiver", body.Receiver),
		))
	}

	return nil
}

// extractCrossChainInputCoins returns the input coins for a cross-chain swap intent.
func extractCrossChainInputCoins(body json.RawMessage) sdk.Coins {
	var b types.CrossChainSwapIntent
	if err := json.Unmarshal(body, &b); err != nil {
		return nil
	}
	if b.InputAmount.IsNil() || !b.InputAmount.IsPositive() {
		return nil
	}
	return sdk.NewCoins(sdk.NewCoin(b.InputDenom, b.InputAmount))
}

// Note: Cross-chain swap intents execute immediately when submitted (no price trigger).
// They're effectively "swap + IBC send" in one atomic operation.
// For price-triggered cross-chain swaps, users should use a conditional intent chain:
//   Step 1: limit_buy with price condition
//   Step 2: cross_chain_swap with immediate condition (uses step 1's output)

// Unused import guard
var _ = math.ZeroInt
