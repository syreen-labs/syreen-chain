package keeper

import (
	"context"
	"sort"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/dex/types"
)

// MaxRouteHops is the maximum number of intermediate pools in a route.
const MaxRouteHops = 3

// MaxSplitParts is the maximum number of split portions for order splitting.
const MaxSplitParts = 5

// RouteHop represents a single hop in a multi-pool swap route.
type RouteHop struct {
	PoolID   uint64 `json:"pool_id"`
	DenomIn  string `json:"denom_in"`
	DenomOut string `json:"denom_out"`
}

// Route represents a complete path from input denom to output denom.
type Route struct {
	Hops         []RouteHop     `json:"hops"`
	ExpectedOut  math.Int       `json:"expected_out"`
	PriceImpact  math.LegacyDec `json:"price_impact"`
	TotalFees    math.LegacyDec `json:"total_fees"`
	Score        math.LegacyDec `json:"score"`
}

// SplitRoute represents a portion of an order routed through a specific path.
type SplitRoute struct {
	Route       Route    `json:"route"`
	InputAmount math.Int `json:"input_amount"`
	OutputAmount math.Int `json:"output_amount"`
}

// OptimalRouteResult is the result of the smart order router.
type OptimalRouteResult struct {
	InputDenom   string         `json:"input_denom"`
	OutputDenom  string         `json:"output_denom"`
	InputAmount  math.Int       `json:"input_amount"`
	OutputAmount math.Int       `json:"output_amount"`
	PriceImpact  math.LegacyDec `json:"price_impact"`
	Routes       []SplitRoute   `json:"routes"`
	IsSplit      bool           `json:"is_split"`
	RouteType    string         `json:"route_type"` // "direct", "multi_hop", "split"
}

// ---------------------------------------------------------------------------
// Pool graph building
// ---------------------------------------------------------------------------

// poolEdge represents a connection between two denoms via a pool.
type poolEdge struct {
	poolID   uint64
	denomA   string
	denomB   string
	pool     types.Pool
}

// buildPoolGraph creates an adjacency map from denom -> list of edges.
// This is used for finding all possible routes between two tokens.
func (k Keeper) buildPoolGraph(ctx context.Context) map[string][]poolEdge {
	pools := k.GetAllPools(ctx)
	graph := make(map[string][]poolEdge, len(pools)*2)

	for _, pool := range pools {
		// Skip pools with zero reserves
		if pool.ReserveA.IsZero() || pool.ReserveB.IsZero() {
			continue
		}
		// Skip halted pools
		if halted, _ := k.IsPoolHalted(ctx, pool.ID); halted {
			continue
		}

		edge := poolEdge{
			poolID: pool.ID,
			denomA: pool.DenomA,
			denomB: pool.DenomB,
			pool:   pool,
		}
		graph[pool.DenomA] = append(graph[pool.DenomA], edge)
		graph[pool.DenomB] = append(graph[pool.DenomB], edge)
	}

	return graph
}

// ---------------------------------------------------------------------------
// Route finding: DFS up to MaxRouteHops
// ---------------------------------------------------------------------------

// FindAllRoutes discovers all possible routes from inputDenom to outputDenom
// up to MaxRouteHops hops. Returns unscored routes (no amounts calculated yet).
func (k Keeper) FindAllRoutes(ctx context.Context, inputDenom, outputDenom string) [][]RouteHop {
	graph := k.buildPoolGraph(ctx)
	var results [][]RouteHop

	// DFS with visited set to avoid cycles
	visited := make(map[uint64]bool)
	var currentPath []RouteHop

	var dfs func(currentDenom string, depth int)
	dfs = func(currentDenom string, depth int) {
		if depth > MaxRouteHops {
			return
		}
		if currentDenom == outputDenom && len(currentPath) > 0 {
			// Found a complete route, copy it
			route := make([]RouteHop, len(currentPath))
			copy(route, currentPath)
			results = append(results, route)
			return
		}

		for _, edge := range graph[currentDenom] {
			if visited[edge.poolID] {
				continue
			}

			// Determine the output denom for this hop
			var nextDenom string
			if currentDenom == edge.denomA {
				nextDenom = edge.denomB
			} else {
				nextDenom = edge.denomA
			}

			visited[edge.poolID] = true
			currentPath = append(currentPath, RouteHop{
				PoolID:   edge.poolID,
				DenomIn:  currentDenom,
				DenomOut: nextDenom,
			})

			dfs(nextDenom, depth+1)

			currentPath = currentPath[:len(currentPath)-1]
			visited[edge.poolID] = false
		}
	}

	dfs(inputDenom, 0)
	return results
}

// ---------------------------------------------------------------------------
// Simulate swap through a route (read-only, no state changes)
// ---------------------------------------------------------------------------

// simulateHop calculates the output amount for a single hop without modifying state.
func (k Keeper) simulateHop(ctx context.Context, hop RouteHop, amountIn math.Int) (amountOut math.Int, priceImpact math.LegacyDec, err error) {
	pool, found := k.GetPool(ctx, hop.PoolID)
	if !found {
		return math.Int{}, math.LegacyZeroDec(), types.ErrPoolNotFound
	}

	var reserveIn, reserveOut math.Int

	if hop.DenomIn == pool.DenomA {
		reserveIn = pool.ReserveA
		reserveOut = pool.ReserveB
	} else if hop.DenomIn == pool.DenomB {
		reserveIn = pool.ReserveB
		reserveOut = pool.ReserveA
	} else {
		return math.Int{}, math.LegacyZeroDec(), types.ErrInvalidDenom
	}

	if reserveIn.IsZero() || reserveOut.IsZero() {
		return math.Int{}, math.LegacyZeroDec(), types.ErrInsufficientLiquidity
	}

	// Calculate output using constant-product formula with fee
	feeBps := pool.SwapFee.MulInt64(10000).TruncateInt64()
	amountInAfterFee := amountIn.MulRaw(10000 - feeBps).QuoRaw(10000)
	outputAmount := reserveOut.Mul(amountInAfterFee).Quo(reserveIn.Add(amountInAfterFee))

	if outputAmount.IsZero() {
		return math.Int{}, math.LegacyZeroDec(), types.ErrInsufficientLiquidity
	}

	// Calculate price impact for this hop
	spotNumerator := math.LegacyNewDecFromInt(outputAmount).Mul(math.LegacyNewDecFromInt(reserveIn))
	spotDenominator := math.LegacyNewDecFromInt(amountIn).Mul(math.LegacyNewDecFromInt(reserveOut))

	impact := math.LegacyZeroDec()
	if spotDenominator.IsPositive() {
		ratio := spotNumerator.Quo(spotDenominator)
		impact = math.LegacyOneDec().Sub(ratio)
		if impact.IsNegative() {
			impact = math.LegacyZeroDec()
		}
	}

	return outputAmount, impact, nil
}

// SimulateRoute calculates the expected output for an entire multi-hop route.
func (k Keeper) SimulateRoute(ctx context.Context, hops []RouteHop, amountIn math.Int) (amountOut math.Int, totalImpact math.LegacyDec, totalFee math.LegacyDec, err error) {
	currentAmount := amountIn
	totalImpact = math.LegacyZeroDec()
	totalFee = math.LegacyZeroDec()

	for _, hop := range hops {
		pool, found := k.GetPool(ctx, hop.PoolID)
		if !found {
			return math.Int{}, math.LegacyZeroDec(), math.LegacyZeroDec(), types.ErrPoolNotFound
		}

		hopOut, hopImpact, hopErr := k.simulateHop(ctx, hop, currentAmount)
		if hopErr != nil {
			return math.Int{}, math.LegacyZeroDec(), math.LegacyZeroDec(), hopErr
		}

		// Accumulate price impact (compound: 1 - (1-i1)*(1-i2)*...)
		// Simplified for lightweight computation: sum of impacts
		totalImpact = totalImpact.Add(hopImpact)

		// Accumulate fees
		totalFee = totalFee.Add(pool.SwapFee)

		currentAmount = hopOut
	}

	return currentAmount, totalImpact, totalFee, nil
}

// ---------------------------------------------------------------------------
// Route scoring
// ---------------------------------------------------------------------------

// scoreRoute assigns a composite score to a route. Higher score = better.
// Score factors: output amount (most important), price impact (penalty), fees (penalty).
func (k Keeper) scoreRoute(route Route, inputAmount math.Int) math.LegacyDec {
	if route.ExpectedOut.IsZero() || inputAmount.IsZero() {
		return math.LegacyZeroDec()
	}

	// Base score: effective output per input (higher is better)
	effectiveRate := math.LegacyNewDecFromInt(route.ExpectedOut).Quo(math.LegacyNewDecFromInt(inputAmount))

	// Penalty for price impact (subtract impact from score)
	impactPenalty := route.PriceImpact

	// Penalty for fees (subtract total fee percentage from score)
	feePenalty := route.TotalFees

	// Penalty for number of hops (more hops = more execution risk)
	hopPenalty := math.LegacyNewDecWithPrec(int64(len(route.Hops)-1), 2) // 0.01 per extra hop

	score := effectiveRate.Sub(impactPenalty).Sub(feePenalty).Sub(hopPenalty)
	return score
}

// ---------------------------------------------------------------------------
// Order splitting
// ---------------------------------------------------------------------------

// simulateRouteSplit simulates splitting an order across a route to find
// the output for a partial amount. Returns the output amount.
func (k Keeper) simulateRouteSplit(ctx context.Context, hops []RouteHop, amount math.Int) (math.Int, error) {
	out, _, _, err := k.SimulateRoute(ctx, hops, amount)
	return out, err
}

// findOptimalSplit finds the best split ratio between two routes to maximize output.
// Uses a simple grid search with MaxSplitParts divisions for lightweight computation.
func (k Keeper) findOptimalSplit(ctx context.Context, route1Hops, route2Hops []RouteHop, totalAmount math.Int) (amount1, amount2, out1, out2 math.Int, err error) {
	bestTotal := math.ZeroInt()

	for i := 0; i <= MaxSplitParts; i++ {
		// split ratio: i/MaxSplitParts goes to route1, rest to route2
		split1 := totalAmount.MulRaw(int64(i)).QuoRaw(int64(MaxSplitParts))
		split2 := totalAmount.Sub(split1)

		var o1, o2 math.Int

		if split1.IsPositive() {
			o1, err = k.simulateRouteSplit(ctx, route1Hops, split1)
			if err != nil {
				o1 = math.ZeroInt()
			}
		} else {
			o1 = math.ZeroInt()
		}

		if split2.IsPositive() {
			o2, err = k.simulateRouteSplit(ctx, route2Hops, split2)
			if err != nil {
				o2 = math.ZeroInt()
			}
		} else {
			o2 = math.ZeroInt()
		}

		total := o1.Add(o2)
		if total.GT(bestTotal) {
			bestTotal = total
			amount1 = split1
			amount2 = split2
			out1 = o1
			out2 = o2
		}
	}

	return amount1, amount2, out1, out2, nil
}

// ---------------------------------------------------------------------------
// Main routing entry point: FindOptimalRoute
// ---------------------------------------------------------------------------

// FindOptimalRoute discovers the best execution path for a swap.
// It compares direct swap, multi-hop routes, and split orders, then returns
// the strategy that yields the best output amount.
func (k Keeper) FindOptimalRoute(ctx context.Context, inputDenom, outputDenom string, inputAmount math.Int) (*OptimalRouteResult, error) {
	if inputDenom == outputDenom {
		return nil, types.ErrSameDenom
	}
	if !inputAmount.IsPositive() {
		return nil, types.ErrInvalidAmount
	}

	// Find all possible routes
	allRoutePaths := k.FindAllRoutes(ctx, inputDenom, outputDenom)

	if len(allRoutePaths) == 0 {
		return nil, types.ErrNoRouteFound
	}

	// Score all routes
	var scoredRoutes []Route
	for _, hops := range allRoutePaths {
		amountOut, impact, fees, err := k.SimulateRoute(ctx, hops, inputAmount)
		if err != nil {
			continue
		}
		if amountOut.IsZero() {
			continue
		}

		route := Route{
			Hops:        hops,
			ExpectedOut: amountOut,
			PriceImpact: impact,
			TotalFees:   fees,
		}
		route.Score = k.scoreRoute(route, inputAmount)
		scoredRoutes = append(scoredRoutes, route)
	}

	if len(scoredRoutes) == 0 {
		return nil, types.ErrNoRouteFound
	}

	// Sort by output amount descending (best output first)
	sort.Slice(scoredRoutes, func(i, j int) bool {
		if scoredRoutes[i].ExpectedOut.Equal(scoredRoutes[j].ExpectedOut) {
			// Tiebreaker: prefer shorter routes, then lower pool IDs
			if len(scoredRoutes[i].Hops) != len(scoredRoutes[j].Hops) {
				return len(scoredRoutes[i].Hops) < len(scoredRoutes[j].Hops)
			}
			for h := 0; h < len(scoredRoutes[i].Hops) && h < len(scoredRoutes[j].Hops); h++ {
				if scoredRoutes[i].Hops[h].PoolID != scoredRoutes[j].Hops[h].PoolID {
					return scoredRoutes[i].Hops[h].PoolID < scoredRoutes[j].Hops[h].PoolID
				}
			}
			return false
		}
		return scoredRoutes[i].ExpectedOut.GT(scoredRoutes[j].ExpectedOut)
	})

	bestRoute := scoredRoutes[0]

	// Try order splitting if there are at least 2 routes
	// Split is worthwhile when a large order causes high price impact
	result := &OptimalRouteResult{
		InputDenom:  inputDenom,
		OutputDenom: outputDenom,
		InputAmount: inputAmount,
	}

	if len(scoredRoutes) >= 2 {
		// Try splitting between top 2 routes
		a1, a2, o1, o2, err := k.findOptimalSplit(ctx, scoredRoutes[0].Hops, scoredRoutes[1].Hops, inputAmount)
		if err == nil {
			splitTotal := o1.Add(o2)
			if splitTotal.GT(bestRoute.ExpectedOut) {
				// Split is better
				result.OutputAmount = splitTotal
				result.IsSplit = true
				result.RouteType = "split"

				// Calculate combined price impact by comparing split output vs best single-route output
				if !bestRoute.ExpectedOut.IsZero() {
					result.PriceImpact = math.LegacyOneDec().Sub(math.LegacyNewDecFromInt(splitTotal).Quo(math.LegacyNewDecFromInt(bestRoute.ExpectedOut)))
					if result.PriceImpact.IsNegative() {
						result.PriceImpact = math.LegacyZeroDec()
					}
				}

				if o1.IsPositive() && a1.IsPositive() {
					_, impact1, fees1, _ := k.SimulateRoute(ctx, scoredRoutes[0].Hops, a1)
					result.Routes = append(result.Routes, SplitRoute{
						Route: Route{
							Hops:        scoredRoutes[0].Hops,
							ExpectedOut: o1,
							PriceImpact: impact1,
							TotalFees:   fees1,
						},
						InputAmount:  a1,
						OutputAmount: o1,
					})
				}
				if o2.IsPositive() && a2.IsPositive() {
					_, impact2, fees2, _ := k.SimulateRoute(ctx, scoredRoutes[1].Hops, a2)
					result.Routes = append(result.Routes, SplitRoute{
						Route: Route{
							Hops:        scoredRoutes[1].Hops,
							ExpectedOut: o2,
							PriceImpact: impact2,
							TotalFees:   fees2,
						},
						InputAmount:  a2,
						OutputAmount: o2,
					})
				}

				return result, nil
			}
		}
	}

	// Single best route (no split)
	result.OutputAmount = bestRoute.ExpectedOut
	result.PriceImpact = bestRoute.PriceImpact
	result.IsSplit = false

	if len(bestRoute.Hops) == 1 {
		result.RouteType = "direct"
	} else {
		result.RouteType = "multi_hop"
	}

	result.Routes = []SplitRoute{
		{
			Route:        bestRoute,
			InputAmount:  inputAmount,
			OutputAmount: bestRoute.ExpectedOut,
		},
	}

	return result, nil
}

// ---------------------------------------------------------------------------
// Smart Swap Execution: execute a swap using the optimal route
// ---------------------------------------------------------------------------

// SmartSwap executes a swap using the smart order router.
// It finds the optimal route and executes each hop sequentially.
func (k Keeper) SmartSwap(ctx context.Context, sender string, inputDenom, outputDenom string, inputAmount, minOutputAmount math.Int) (sdk.Coin, error) {
	// Find optimal route
	optimalRoute, err := k.FindOptimalRoute(ctx, inputDenom, outputDenom, inputAmount)
	if err != nil {
		return sdk.Coin{}, err
	}

	// Check slippage against total expected output
	if !minOutputAmount.IsNil() && !minOutputAmount.IsZero() && optimalRoute.OutputAmount.LT(minOutputAmount) {
		return sdk.Coin{}, types.ErrSlippageExceeded
	}

	totalOutput := math.ZeroInt()

	// Execute each split route
	for _, splitRoute := range optimalRoute.Routes {
		if !splitRoute.InputAmount.IsPositive() {
			continue
		}

		// D-06: Simulate each hop to get per-hop expected outputs for sandwich protection
		hopExpectedOutputs := make([]math.Int, len(splitRoute.Route.Hops))
		simAmount := splitRoute.InputAmount
		for i, hop := range splitRoute.Route.Hops {
			simOut, _, simErr := k.simulateHop(ctx, hop, simAmount)
			if simErr != nil {
				// Fallback: use 90% of input as minimum to prevent catastrophic loss
				hopExpectedOutputs[i] = simAmount.MulRaw(90).QuoRaw(100)
			} else {
				hopExpectedOutputs[i] = simOut
			}
			simAmount = hopExpectedOutputs[i]
		}

		currentDenom := inputDenom
		currentAmount := splitRoute.InputAmount

		// Execute each hop in this route
		for i, hop := range splitRoute.Route.Hops {
			tokenIn := sdk.NewCoin(currentDenom, currentAmount)

			var minOut math.Int
			if i < len(splitRoute.Route.Hops)-1 {
				// D-06: Intermediate hops use 2% per-hop slippage tolerance
				// to prevent sandwich attacks between hops
				minOut = hopExpectedOutputs[i].MulRaw(98).QuoRaw(100)
			} else {
				// Final hop: the overall minOutputAmount check below protects the user
				minOut = math.ZeroInt()
			}

			tokenOut, swapErr := k.Swap(ctx, sender, hop.PoolID, tokenIn, minOut)
			if swapErr != nil {
				return sdk.Coin{}, swapErr
			}

			currentDenom = hop.DenomOut
			currentAmount = tokenOut.Amount

			// Verify denom matches what we expected
			if i < len(splitRoute.Route.Hops)-1 {
				// Intermediate hop - output denom should match next hop input
				if tokenOut.Denom != splitRoute.Route.Hops[i+1].DenomIn {
					return sdk.Coin{}, types.ErrInvalidDenom
				}
			}
		}

		totalOutput = totalOutput.Add(currentAmount)
	}

	// Final slippage check on actual execution
	if !minOutputAmount.IsNil() && !minOutputAmount.IsZero() && totalOutput.LT(minOutputAmount) {
		// This shouldn't happen since we checked the simulation, but safety first
		return sdk.Coin{}, types.ErrSlippageExceeded
	}

	return sdk.NewCoin(outputDenom, totalOutput), nil
}

// ---------------------------------------------------------------------------
// GetAllRoutes returns all discoverable routes between two denoms (query only)
// ---------------------------------------------------------------------------

// GetAllRoutesWithScores returns all routes between two denoms with scoring.
func (k Keeper) GetAllRoutesWithScores(ctx context.Context, inputDenom, outputDenom string, amount math.Int) ([]Route, error) {
	if inputDenom == outputDenom {
		return nil, types.ErrSameDenom
	}

	allRoutePaths := k.FindAllRoutes(ctx, inputDenom, outputDenom)

	var routes []Route
	for _, hops := range allRoutePaths {
		amountOut, impact, fees, err := k.SimulateRoute(ctx, hops, amount)
		if err != nil {
			continue
		}
		if amountOut.IsZero() {
			continue
		}

		route := Route{
			Hops:        hops,
			ExpectedOut: amountOut,
			PriceImpact: impact,
			TotalFees:   fees,
		}
		route.Score = k.scoreRoute(route, amount)
		routes = append(routes, route)
	}

	// Sort by score descending
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Score.Equal(routes[j].Score) {
			// Tiebreaker: first hop pool ID
			if len(routes[i].Hops) > 0 && len(routes[j].Hops) > 0 {
				return routes[i].Hops[0].PoolID < routes[j].Hops[0].PoolID
			}
			return false
		}
		return routes[i].Score.GT(routes[j].Score)
	})

	return routes, nil
}
