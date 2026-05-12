package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"syreen/x/dex/types"
)

// ---------------------------------------------------------------------------
// Helper: create a pool between two denoms and fund the module account
// ---------------------------------------------------------------------------

func createPoolForRouter(t *testing.T, k interface{ CreatePool(ctx context.Context, creator, denomA, denomB string, amountA, amountB math.Int) (uint64, error) }, ctx sdk.Context, bk *mockBankKeeper, creator, denomA, denomB string, amountA, amountB math.Int) uint64 {
	t.Helper()

	bk.fundAccount(creator, sdk.NewCoins(
		sdk.NewCoin(denomA, amountA),
		sdk.NewCoin(denomB, amountB),
	))

	poolID, err := k.CreatePool(ctx, creator, denomA, denomB, amountA, amountB)
	require.NoError(t, err)

	// Fund module account with reserves so swaps can send tokens out
	moduleCoins := sdk.NewCoins(
		sdk.NewCoin(denomA, amountA),
		sdk.NewCoin(denomB, amountB),
	)
	bk.fundAccount(types.ModuleName, moduleCoins)

	return poolID
}

// ---------------------------------------------------------------------------
// TestFindAllRoutes_DirectRoute
// ---------------------------------------------------------------------------

func TestFindAllRoutes_DirectRoute(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	// Create a single pool: usyreen <-> uusdc
	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "uusdc",
		math.NewInt(10_000_000), math.NewInt(10_000_000))

	// Find routes from usyreen to uusdc
	routes := k.FindAllRoutes(ctx, "usyreen", "uusdc")
	require.Len(t, routes, 1)
	require.Len(t, routes[0], 1)
	require.Equal(t, "usyreen", routes[0][0].DenomIn)
	require.Equal(t, "uusdc", routes[0][0].DenomOut)
}

// ---------------------------------------------------------------------------
// TestFindAllRoutes_MultiHop
// ---------------------------------------------------------------------------

func TestFindAllRoutes_MultiHop(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	// Create pools: usyreen <-> ueth, ueth <-> uusdc
	// No direct usyreen <-> uusdc pool
	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "ueth",
		math.NewInt(10_000_000), math.NewInt(10_000_000))
	createPoolForRouter(t, k, ctx, bk, creator, "ueth", "uusdc",
		math.NewInt(10_000_000), math.NewInt(10_000_000))

	// Find routes from usyreen to uusdc
	routes := k.FindAllRoutes(ctx, "usyreen", "uusdc")
	require.Len(t, routes, 1) // usyreen -> ueth -> uusdc

	require.Len(t, routes[0], 2)
	require.Equal(t, "usyreen", routes[0][0].DenomIn)
	require.Equal(t, "ueth", routes[0][0].DenomOut)
	require.Equal(t, "ueth", routes[0][1].DenomIn)
	require.Equal(t, "uusdc", routes[0][1].DenomOut)
}

// ---------------------------------------------------------------------------
// TestFindAllRoutes_MultipleRoutes
// ---------------------------------------------------------------------------

func TestFindAllRoutes_MultipleRoutes(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	// Create triangle: usyreen <-> uusdc, usyreen <-> ueth, ueth <-> uusdc
	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "uusdc",
		math.NewInt(10_000_000), math.NewInt(10_000_000))
	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "ueth",
		math.NewInt(10_000_000), math.NewInt(10_000_000))
	createPoolForRouter(t, k, ctx, bk, creator, "ueth", "uusdc",
		math.NewInt(10_000_000), math.NewInt(10_000_000))

	// Find routes from usyreen to uusdc
	routes := k.FindAllRoutes(ctx, "usyreen", "uusdc")
	// Should find 2 routes: direct (usyreen->uusdc) and multi-hop (usyreen->ueth->uusdc)
	require.GreaterOrEqual(t, len(routes), 2)

	// Verify all routes start with usyreen and end with uusdc
	for _, route := range routes {
		require.Equal(t, "usyreen", route[0].DenomIn)
		require.Equal(t, "uusdc", route[len(route)-1].DenomOut)
	}
}

// ---------------------------------------------------------------------------
// TestFindAllRoutes_NoRoute
// ---------------------------------------------------------------------------

func TestFindAllRoutes_NoRoute(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	// Create only usyreen <-> ueth pool
	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "ueth",
		math.NewInt(10_000_000), math.NewInt(10_000_000))

	// Try to find route from usyreen to uusdc (no such route)
	routes := k.FindAllRoutes(ctx, "usyreen", "uusdc")
	require.Len(t, routes, 0)
}

// ---------------------------------------------------------------------------
// TestFindAllRoutes_ThreeHops
// ---------------------------------------------------------------------------

func TestFindAllRoutes_ThreeHops(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	// Create chain: usyreen <-> ueth, ueth <-> ubtc, ubtc <-> uusdc
	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "ueth",
		math.NewInt(10_000_000), math.NewInt(10_000_000))
	createPoolForRouter(t, k, ctx, bk, creator, "ueth", "ubtc",
		math.NewInt(10_000_000), math.NewInt(10_000_000))
	createPoolForRouter(t, k, ctx, bk, creator, "ubtc", "uusdc",
		math.NewInt(10_000_000), math.NewInt(10_000_000))

	// Find routes from usyreen to uusdc (3 hops)
	routes := k.FindAllRoutes(ctx, "usyreen", "uusdc")
	require.GreaterOrEqual(t, len(routes), 1)

	// Should find the 3-hop route
	found3Hop := false
	for _, route := range routes {
		if len(route) == 3 {
			found3Hop = true
			require.Equal(t, "usyreen", route[0].DenomIn)
			require.Equal(t, "ueth", route[0].DenomOut)
			require.Equal(t, "ueth", route[1].DenomIn)
			require.Equal(t, "ubtc", route[1].DenomOut)
			require.Equal(t, "ubtc", route[2].DenomIn)
			require.Equal(t, "uusdc", route[2].DenomOut)
		}
	}
	require.True(t, found3Hop, "should find a 3-hop route")
}

// ---------------------------------------------------------------------------
// TestSimulateRoute_SingleHop
// ---------------------------------------------------------------------------

func TestSimulateRoute_SingleHop(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	poolID := createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "uusdc",
		math.NewInt(10_000_000), math.NewInt(10_000_000))

	hops := []struct {
		PoolID   uint64
		DenomIn  string
		DenomOut string
	}{
		{poolID, "usyreen", "uusdc"},
	}

	var routeHops []struct {
		PoolID   uint64
		DenomIn  string
		DenomOut string
	}
	routeHops = hops
	_ = routeHops

	// Use the keeper's RouteHop type
	routes := k.FindAllRoutes(ctx, "usyreen", "uusdc")
	require.Len(t, routes, 1)

	amountOut, impact, fees, err := k.SimulateRoute(ctx, routes[0], math.NewInt(1_000_000))
	require.NoError(t, err)
	require.True(t, amountOut.IsPositive())
	require.True(t, amountOut.LT(math.NewInt(1_000_000))) // less due to fee + impact
	require.True(t, impact.IsPositive())
	require.True(t, fees.IsPositive())
}

// ---------------------------------------------------------------------------
// TestSimulateRoute_MultiHop
// ---------------------------------------------------------------------------

func TestSimulateRoute_MultiHop(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	// Create two pools: usyreen<->ueth and ueth<->uusdc
	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "ueth",
		math.NewInt(10_000_000), math.NewInt(10_000_000))
	createPoolForRouter(t, k, ctx, bk, creator, "ueth", "uusdc",
		math.NewInt(10_000_000), math.NewInt(10_000_000))

	routes := k.FindAllRoutes(ctx, "usyreen", "uusdc")
	require.Len(t, routes, 1)

	amountOut, impact, fees, err := k.SimulateRoute(ctx, routes[0], math.NewInt(1_000_000))
	require.NoError(t, err)
	require.True(t, amountOut.IsPositive())
	// 2 hops means more fees, so output should be significantly less
	require.True(t, amountOut.LT(math.NewInt(1_000_000)))
	require.True(t, impact.IsPositive())
	// Total fees should be ~0.006 (0.3% * 2)
	require.True(t, fees.GT(math.LegacyNewDecWithPrec(3, 3))) // > 0.3%
}

// ---------------------------------------------------------------------------
// TestSimulateRoute_ConsistentWithGetQuote
// ---------------------------------------------------------------------------

func TestSimulateRoute_ConsistentWithGetQuote(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	poolID := createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "uusdc",
		math.NewInt(10_000_000), math.NewInt(10_000_000))

	// GetQuote for direct pool
	tokenIn := sdk.NewInt64Coin("usyreen", 1_000_000)
	quoteOut, _, err := k.GetQuote(ctx, poolID, tokenIn)
	require.NoError(t, err)

	// SimulateRoute should give same result for a single-hop route
	routes := k.FindAllRoutes(ctx, "usyreen", "uusdc")
	require.Len(t, routes, 1)
	routeOut, _, _, routeErr := k.SimulateRoute(ctx, routes[0], math.NewInt(1_000_000))
	require.NoError(t, routeErr)

	require.Equal(t, quoteOut.Amount, routeOut, "SimulateRoute should match GetQuote for single-hop")
}

// ---------------------------------------------------------------------------
// TestFindOptimalRoute_DirectPool
// ---------------------------------------------------------------------------

func TestFindOptimalRoute_DirectPool(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "uusdc",
		math.NewInt(10_000_000), math.NewInt(10_000_000))

	result, err := k.FindOptimalRoute(ctx, "usyreen", "uusdc", math.NewInt(1_000_000))
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "direct", result.RouteType)
	require.False(t, result.IsSplit)
	require.True(t, result.OutputAmount.IsPositive())
	require.Equal(t, "usyreen", result.InputDenom)
	require.Equal(t, "uusdc", result.OutputDenom)
	require.Len(t, result.Routes, 1)
	require.Len(t, result.Routes[0].Route.Hops, 1)
}

// ---------------------------------------------------------------------------
// TestFindOptimalRoute_MultiHopBetterThanDirect
// ---------------------------------------------------------------------------

func TestFindOptimalRoute_ChoosesBestRoute(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	// Create direct pool with low liquidity (high price impact)
	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "uusdc",
		math.NewInt(2_000_000), math.NewInt(2_000_000))

	// Create indirect route with high liquidity (low price impact)
	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "ueth",
		math.NewInt(100_000_000), math.NewInt(100_000_000))
	createPoolForRouter(t, k, ctx, bk, creator, "ueth", "uusdc",
		math.NewInt(100_000_000), math.NewInt(100_000_000))

	// For a large swap, the multi-hop through deep pools should yield more
	result, err := k.FindOptimalRoute(ctx, "usyreen", "uusdc", math.NewInt(1_000_000))
	require.NoError(t, err)
	require.NotNil(t, result)
	require.True(t, result.OutputAmount.IsPositive())
	require.Len(t, result.Routes, 1)

	// The router should pick the route with the best output
	// For a 1M swap into a 2M pool, the direct route has massive impact
	// The multi-hop through 100M pools should be better
	require.Equal(t, "multi_hop", result.RouteType)
}

// ---------------------------------------------------------------------------
// TestFindOptimalRoute_PrefersDirect_WhenBetter
// ---------------------------------------------------------------------------

func TestFindOptimalRoute_PrefersDirect_WhenBetter(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	// Create direct pool with high liquidity
	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "uusdc",
		math.NewInt(100_000_000), math.NewInt(100_000_000))

	// Create indirect route with low liquidity
	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "ueth",
		math.NewInt(2_000_000), math.NewInt(2_000_000))
	createPoolForRouter(t, k, ctx, bk, creator, "ueth", "uusdc",
		math.NewInt(2_000_000), math.NewInt(2_000_000))

	// For a moderate swap, direct should be better
	result, err := k.FindOptimalRoute(ctx, "usyreen", "uusdc", math.NewInt(100_000))
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "direct", result.RouteType)
}

// ---------------------------------------------------------------------------
// TestFindOptimalRoute_NoRoute
// ---------------------------------------------------------------------------

func TestFindOptimalRoute_NoRoute(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "ueth",
		math.NewInt(10_000_000), math.NewInt(10_000_000))

	// No route from usyreen to uusdc
	_, err := k.FindOptimalRoute(ctx, "usyreen", "uusdc", math.NewInt(1_000_000))
	require.ErrorIs(t, err, types.ErrNoRouteFound)
}

// ---------------------------------------------------------------------------
// TestFindOptimalRoute_SameDenom
// ---------------------------------------------------------------------------

func TestFindOptimalRoute_SameDenom(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	_, err := k.FindOptimalRoute(ctx, "usyreen", "usyreen", math.NewInt(1_000_000))
	require.ErrorIs(t, err, types.ErrSameDenom)
}

// ---------------------------------------------------------------------------
// TestFindOptimalRoute_InvalidAmount
// ---------------------------------------------------------------------------

func TestFindOptimalRoute_InvalidAmount(t *testing.T) {
	k, ctx, _, _ := setupKeeper(t)

	_, err := k.FindOptimalRoute(ctx, "usyreen", "uusdc", math.ZeroInt())
	require.ErrorIs(t, err, types.ErrInvalidAmount)
}

// ---------------------------------------------------------------------------
// TestSmartSwap_DirectSwap
// ---------------------------------------------------------------------------

func TestSmartSwap_DirectSwap(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()
	trader := testAddr2()

	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "uusdc",
		math.NewInt(10_000_000), math.NewInt(10_000_000))

	// Fund trader
	bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1_000_000)))

	tokenOut, err := k.SmartSwap(ctx, trader, "usyreen", "uusdc", math.NewInt(1_000_000), math.ZeroInt())
	require.NoError(t, err)
	require.Equal(t, "uusdc", tokenOut.Denom)
	require.True(t, tokenOut.Amount.IsPositive())
	require.True(t, tokenOut.Amount.LT(math.NewInt(1_000_000)))
}

// ---------------------------------------------------------------------------
// TestSmartSwap_MultiHop
// ---------------------------------------------------------------------------

func TestSmartSwap_MultiHop(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()
	trader := testAddr2()

	// No direct pool, only indirect route
	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "ueth",
		math.NewInt(10_000_000), math.NewInt(10_000_000))
	createPoolForRouter(t, k, ctx, bk, creator, "ueth", "uusdc",
		math.NewInt(10_000_000), math.NewInt(10_000_000))

	// Fund trader
	bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1_000_000)))

	tokenOut, err := k.SmartSwap(ctx, trader, "usyreen", "uusdc", math.NewInt(1_000_000), math.ZeroInt())
	require.NoError(t, err)
	require.Equal(t, "uusdc", tokenOut.Denom)
	require.True(t, tokenOut.Amount.IsPositive())
}

// ---------------------------------------------------------------------------
// TestSmartSwap_SlippageCheck
// ---------------------------------------------------------------------------

func TestSmartSwap_SlippageCheck(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()
	trader := testAddr2()

	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "uusdc",
		math.NewInt(10_000_000), math.NewInt(10_000_000))

	bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1_000_000)))

	// Set unrealistically high min output
	_, err := k.SmartSwap(ctx, trader, "usyreen", "uusdc", math.NewInt(1_000_000), math.NewInt(999_999))
	require.ErrorIs(t, err, types.ErrSlippageExceeded)
}

// ---------------------------------------------------------------------------
// TestSmartSwap_NoRoute
// ---------------------------------------------------------------------------

func TestSmartSwap_NoRoute(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()
	trader := testAddr2()

	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "ueth",
		math.NewInt(10_000_000), math.NewInt(10_000_000))

	bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 1_000_000)))

	_, err := k.SmartSwap(ctx, trader, "usyreen", "uusdc", math.NewInt(1_000_000), math.ZeroInt())
	require.ErrorIs(t, err, types.ErrNoRouteFound)
}

// ---------------------------------------------------------------------------
// TestGetAllRoutesWithScores
// ---------------------------------------------------------------------------

func TestGetAllRoutesWithScores(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	// Create triangle
	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "uusdc",
		math.NewInt(10_000_000), math.NewInt(10_000_000))
	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "ueth",
		math.NewInt(10_000_000), math.NewInt(10_000_000))
	createPoolForRouter(t, k, ctx, bk, creator, "ueth", "uusdc",
		math.NewInt(10_000_000), math.NewInt(10_000_000))

	routes, err := k.GetAllRoutesWithScores(ctx, "usyreen", "uusdc", math.NewInt(1_000_000))
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(routes), 2)

	// Routes should be sorted by score (descending)
	for i := 0; i < len(routes)-1; i++ {
		require.True(t, routes[i].Score.GTE(routes[i+1].Score),
			"routes should be sorted by score descending")
	}

	// All routes should have positive output
	for _, r := range routes {
		require.True(t, r.ExpectedOut.IsPositive())
	}
}

// ---------------------------------------------------------------------------
// TestGetAllRoutesWithScores_DirectBetterScore
// ---------------------------------------------------------------------------

func TestGetAllRoutesWithScores_DirectBetterScore(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	// Create deep direct pool and shallow indirect route
	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "uusdc",
		math.NewInt(100_000_000), math.NewInt(100_000_000))
	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "ueth",
		math.NewInt(2_000_000), math.NewInt(2_000_000))
	createPoolForRouter(t, k, ctx, bk, creator, "ueth", "uusdc",
		math.NewInt(2_000_000), math.NewInt(2_000_000))

	routes, err := k.GetAllRoutesWithScores(ctx, "usyreen", "uusdc", math.NewInt(100_000))
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(routes), 2)

	// Direct route (1 hop) should be ranked first
	require.Len(t, routes[0].Hops, 1, "best route should be the direct one")
}

// ---------------------------------------------------------------------------
// TestRouteScoring_PriceImpactMatters
// ---------------------------------------------------------------------------

func TestRouteScoring_PriceImpactMatters(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	// Small pool (high impact)
	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "uusdc",
		math.NewInt(2_000_000), math.NewInt(2_000_000))

	// Large swap to cause significant impact
	routes, err := k.GetAllRoutesWithScores(ctx, "usyreen", "uusdc", math.NewInt(1_000_000))
	require.NoError(t, err)
	require.Len(t, routes, 1)

	// Price impact should be significant for 50% of pool liquidity
	require.True(t, routes[0].PriceImpact.GT(math.LegacyNewDecWithPrec(10, 2)),
		"price impact should be >10%% for a 50%% pool swap, got %s", routes[0].PriceImpact)
}

// ---------------------------------------------------------------------------
// TestFindOptimalRoute_SplitOrder
// ---------------------------------------------------------------------------

func TestFindOptimalRoute_SplitOrder(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	// Create two routes with equal liquidity
	// Direct route
	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "uusdc",
		math.NewInt(5_000_000), math.NewInt(5_000_000))
	// Indirect route via ueth
	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "ueth",
		math.NewInt(5_000_000), math.NewInt(5_000_000))
	createPoolForRouter(t, k, ctx, bk, creator, "ueth", "uusdc",
		math.NewInt(5_000_000), math.NewInt(5_000_000))

	// Large swap that would have high impact on any single route
	result, err := k.FindOptimalRoute(ctx, "usyreen", "uusdc", math.NewInt(3_000_000))
	require.NoError(t, err)
	require.NotNil(t, result)
	require.True(t, result.OutputAmount.IsPositive())

	// Whether split or not, the result should be valid
	if result.IsSplit {
		require.Equal(t, "split", result.RouteType)
		require.Len(t, result.Routes, 2)
		// Total of split amounts should equal input
		totalInput := math.ZeroInt()
		for _, sr := range result.Routes {
			totalInput = totalInput.Add(sr.InputAmount)
		}
		require.Equal(t, math.NewInt(3_000_000), totalInput)
	}
}

// ---------------------------------------------------------------------------
// TestSmartSwap_ConsistentWithDirectSwap
// ---------------------------------------------------------------------------

func TestSmartSwap_ConsistentWithDirectSwap(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()
	_ = testAddr2()

	// Single pool — SmartSwap should give same result as direct Swap
	poolID := createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "uusdc",
		math.NewInt(10_000_000), math.NewInt(10_000_000))

	// Direct swap quote
	tokenIn := sdk.NewInt64Coin("usyreen", 100_000)
	directOut, _, err := k.GetQuote(ctx, poolID, tokenIn)
	require.NoError(t, err)

	// Smart route quote (just simulation, don't execute)
	result, err := k.FindOptimalRoute(ctx, "usyreen", "uusdc", math.NewInt(100_000))
	require.NoError(t, err)

	// Should give the same output for a single-pool scenario
	require.Equal(t, directOut.Amount, result.OutputAmount,
		"SmartSwap output should match direct swap for single pool")
}

// ---------------------------------------------------------------------------
// TestFindAllRoutes_NoCycles
// ---------------------------------------------------------------------------

func TestFindAllRoutes_NoCycles(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	// Create a graph that could have cycles: A-B, B-C, C-A, A-D
	createPoolForRouter(t, k, ctx, bk, creator, "tokena", "tokenb",
		math.NewInt(10_000_000), math.NewInt(10_000_000))
	createPoolForRouter(t, k, ctx, bk, creator, "tokenb", "tokenc",
		math.NewInt(10_000_000), math.NewInt(10_000_000))
	createPoolForRouter(t, k, ctx, bk, creator, "tokena", "tokenc",
		math.NewInt(10_000_000), math.NewInt(10_000_000))
	createPoolForRouter(t, k, ctx, bk, creator, "tokena", "tokend",
		math.NewInt(10_000_000), math.NewInt(10_000_000))

	routes := k.FindAllRoutes(ctx, "tokena", "tokend")
	require.GreaterOrEqual(t, len(routes), 1)

	// Verify no route visits the same pool twice
	for _, route := range routes {
		seen := make(map[uint64]bool)
		for _, hop := range route {
			require.False(t, seen[hop.PoolID], "route should not visit same pool twice")
			seen[hop.PoolID] = true
		}
	}
}

// ---------------------------------------------------------------------------
// TestFindAllRoutes_MaxHopsRespected
// ---------------------------------------------------------------------------

func TestFindAllRoutes_MaxHopsRespected(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	// Create a long chain: A->B->C->D->E (4 hops)
	createPoolForRouter(t, k, ctx, bk, creator, "tokena", "tokenb",
		math.NewInt(10_000_000), math.NewInt(10_000_000))
	createPoolForRouter(t, k, ctx, bk, creator, "tokenb", "tokenc",
		math.NewInt(10_000_000), math.NewInt(10_000_000))
	createPoolForRouter(t, k, ctx, bk, creator, "tokenc", "tokend",
		math.NewInt(10_000_000), math.NewInt(10_000_000))
	createPoolForRouter(t, k, ctx, bk, creator, "tokend", "tokene",
		math.NewInt(10_000_000), math.NewInt(10_000_000))

	// Routes from A to E would need 4 hops, but max is 3
	routes := k.FindAllRoutes(ctx, "tokena", "tokene")

	// Should not find any routes (4 hops exceeds MaxRouteHops=3)
	for _, route := range routes {
		require.LessOrEqual(t, len(route), 3, "route should not exceed MaxRouteHops")
	}
}

// ---------------------------------------------------------------------------
// TestSmartSwap_ActualExecution_MultiHop
// ---------------------------------------------------------------------------

func TestSmartSwap_ActualExecution_MultiHop(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()
	trader := testAddr2()

	// Only indirect route available
	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "ueth",
		math.NewInt(10_000_000), math.NewInt(10_000_000))
	createPoolForRouter(t, k, ctx, bk, creator, "ueth", "uusdc",
		math.NewInt(10_000_000), math.NewInt(10_000_000))

	bk.fundAccount(trader, sdk.NewCoins(sdk.NewInt64Coin("usyreen", 500_000)))

	tokenOut, err := k.SmartSwap(ctx, trader, "usyreen", "uusdc", math.NewInt(500_000), math.ZeroInt())
	require.NoError(t, err)
	require.Equal(t, "uusdc", tokenOut.Denom)
	require.True(t, tokenOut.Amount.IsPositive())

	// Verify trader received the tokens
	traderUusdc := bk.getBalance(trader, "uusdc")
	require.Equal(t, tokenOut.Amount, traderUusdc)

	// Verify trader's usyreen was deducted
	traderUsyreen := bk.getBalance(trader, "usyreen")
	require.True(t, traderUsyreen.IsZero(), "trader should have spent all usyreen")
}

// ---------------------------------------------------------------------------
// TestFindOptimalRoute_ReverseDirection
// ---------------------------------------------------------------------------

func TestFindOptimalRoute_ReverseDirection(t *testing.T) {
	k, ctx, bk, _ := setupKeeper(t)
	creator := testAddr()

	createPoolForRouter(t, k, ctx, bk, creator, "usyreen", "uusdc",
		math.NewInt(10_000_000), math.NewInt(10_000_000))

	// Forward: usyreen -> uusdc
	fwdResult, err := k.FindOptimalRoute(ctx, "usyreen", "uusdc", math.NewInt(1_000_000))
	require.NoError(t, err)

	// Reverse: uusdc -> usyreen
	revResult, err := k.FindOptimalRoute(ctx, "uusdc", "usyreen", math.NewInt(1_000_000))
	require.NoError(t, err)

	// Both should work (symmetric pool)
	require.True(t, fwdResult.OutputAmount.IsPositive())
	require.True(t, revResult.OutputAmount.IsPositive())

	// For a 1:1 pool, outputs should be equal
	require.Equal(t, fwdResult.OutputAmount, revResult.OutputAmount)
}
