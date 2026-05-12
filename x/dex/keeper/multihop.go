package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/dex/types"
)

// MaxMultiHopSwapHops is the maximum number of pools in a multi-hop swap route.
const MaxMultiHopSwapHops = 4

// MultiHopSwap executes an atomic multi-hop swap through a sequence of pools.
// The swap is atomic: if any hop fails, the entire transaction reverts
// (Cosmos SDK handles this via tx-level cache-wrapping).
//
// Parameters:
//   - sender: the address initiating the swap
//   - route: ordered list of pool IDs to swap through
//   - tokenIn: the input token (denom + amount)
//   - minTokenOut: minimum acceptable output amount (slippage protection)
//
// Returns the final output coin or an error.
func (k Keeper) MultiHopSwap(ctx context.Context, sender string, route []uint64, tokenIn sdk.Coin, minTokenOut math.Int) (sdk.Coin, error) {
	// Validate route length
	if len(route) == 0 {
		return sdk.Coin{}, fmt.Errorf("route must contain at least one pool")
	}
	if len(route) > MaxMultiHopSwapHops {
		return sdk.Coin{}, types.ErrTooManyHops
	}

	// Reject duplicate pool IDs to prevent fee extraction loops
	seen := make(map[uint64]bool, len(route))
	for _, pid := range route {
		if seen[pid] {
			return sdk.Coin{}, fmt.Errorf("duplicate pool ID %d in route", pid)
		}
		seen[pid] = true
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentToken := tokenIn

	// Pre-validate that all pools exist and route denoms are connected
	currentDenom := tokenIn.Denom
	for i, poolID := range route {
		if poolID == 0 {
			return sdk.Coin{}, types.ErrPoolNotFound
		}
		pool, found := k.GetPool(ctx, poolID)
		if !found {
			return sdk.Coin{}, fmt.Errorf("pool %d not found at hop %d", poolID, i+1)
		}
		// Verify the current denom is one of the pool's denoms
		if currentDenom != pool.DenomA && currentDenom != pool.DenomB {
			return sdk.Coin{}, fmt.Errorf("denom %s not in pool %d (has %s/%s) at hop %d", currentDenom, poolID, pool.DenomA, pool.DenomB, i+1)
		}
		// Next hop's input denom is the other side of this pool
		if currentDenom == pool.DenomA {
			currentDenom = pool.DenomB
		} else {
			currentDenom = pool.DenomA
		}
	}

	// Execute swaps sequentially through each pool
	for i, poolID := range route {
		if poolID == 0 {
			return sdk.Coin{}, types.ErrPoolNotFound
		}

		// For intermediate hops, use zero min output (final slippage check below)
		hopMinOut := math.ZeroInt()

		tokenOut, err := k.Swap(ctx, sender, poolID, currentToken, hopMinOut)
		if err != nil {
			return sdk.Coin{}, fmt.Errorf("multi-hop swap failed at hop %d (pool %d): %w", i+1, poolID, err)
		}

		k.Logger(ctx).Info("multi-hop swap hop completed",
			"hop", i+1,
			"pool_id", poolID,
			"token_in", currentToken.String(),
			"token_out", tokenOut.String(),
		)

		// The output of this hop becomes the input of the next hop
		currentToken = tokenOut
	}

	// Final slippage check on the overall output
	if !minTokenOut.IsNil() && !minTokenOut.IsZero() && currentToken.Amount.LT(minTokenOut) {
		return sdk.Coin{}, types.ErrSlippageExceeded
	}

	// Emit multi-hop swap event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"multi_hop_swap",
			sdk.NewAttribute("sender", sender),
			sdk.NewAttribute("hops", fmt.Sprintf("%d", len(route))),
			sdk.NewAttribute("token_in", tokenIn.String()),
			sdk.NewAttribute("token_out", currentToken.String()),
		),
	)

	return currentToken, nil
}

// GetBestRoute finds the optimal path through available pools from denomIn to denomOut.
// Uses BFS up to MaxMultiHopSwapHops hops, then scores routes by output amount.
// Returns the best route as a list of pool IDs and the expected output amount.
func (k Keeper) GetBestRoute(ctx context.Context, denomIn, denomOut string, amount math.Int) ([]uint64, math.Int, error) {
	if denomIn == denomOut {
		return nil, math.Int{}, types.ErrSameDenom
	}
	if !amount.IsPositive() {
		return nil, math.Int{}, types.ErrInvalidAmount
	}

	// Use the existing router infrastructure to find all routes
	allRoutePaths := k.FindAllRoutes(ctx, denomIn, denomOut)
	if len(allRoutePaths) == 0 {
		return nil, math.Int{}, types.ErrNoRouteFound
	}

	// Filter to routes within MaxMultiHopSwapHops
	var bestPoolIDs []uint64
	bestOutput := math.ZeroInt()

	for _, hops := range allRoutePaths {
		if len(hops) > MaxMultiHopSwapHops {
			continue
		}

		// Simulate the route to get expected output
		amountOut, _, _, err := k.SimulateRoute(ctx, hops, amount)
		if err != nil {
			continue
		}
		if amountOut.IsZero() {
			continue
		}

		if amountOut.GT(bestOutput) {
			bestOutput = amountOut
			bestPoolIDs = make([]uint64, len(hops))
			for i, hop := range hops {
				bestPoolIDs[i] = hop.PoolID
			}
		}
	}

	if len(bestPoolIDs) == 0 {
		return nil, math.Int{}, types.ErrNoRouteFound
	}

	return bestPoolIDs, bestOutput, nil
}
