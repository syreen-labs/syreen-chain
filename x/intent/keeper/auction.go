package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	sdkmath "cosmossdk.io/math"

	"syreen/x/intent/types"
)

// SelectWinningSolver ranks all solutions for an intent and returns the best one.
// Ranking criteria (weighted score):
//   - Solver reputation (40% weight)
//   - Solution gas efficiency - lower gas estimate = better (30% weight)
//   - Solver stake amount - higher = more skin in game (20% weight)
//   - Solution tip offered to intent creator (10% weight)
func (k Keeper) SelectWinningSolver(ctx context.Context, intentID string) (*types.Solution, error) {
	solutions := k.GetSolutionsForIntent(ctx, intentID)
	if len(solutions) == 0 {
		return nil, types.ErrNoSolutions
	}

	// Single solution: return immediately (if solver is still active)
	if len(solutions) == 1 {
		sol := solutions[0]
		solver, found := k.GetSolver(ctx, sol.SolverAddr)
		if !found || !solver.Active {
			return nil, fmt.Errorf("only solver %s is no longer active", sol.SolverAddr)
		}
		return &sol, nil
	}

	// Collect raw metrics for all solutions from active solvers.
	// All arithmetic uses deterministic sdkmath.LegacyDec (18-digit fixed-point).
	type rawMetrics struct {
		solution   *types.Solution
		reputation sdkmath.LegacyDec
		gas        sdkmath.LegacyDec
		stake      sdkmath.LegacyDec
		tip        sdkmath.LegacyDec
	}

	var metrics []rawMetrics
	for i := range solutions {
		sol := &solutions[i]
		solver, found := k.GetSolver(ctx, sol.SolverAddr)
		if !found || !solver.Active {
			continue // Skip inactive/missing solvers
		}

		tipAmount := sdkmath.LegacyZeroDec()
		if sol.Tip != nil && sol.Tip.IsAllPositive() {
			for _, c := range sol.Tip {
				tipAmount = tipAmount.Add(sdkmath.LegacyNewDecFromBigInt(c.Amount.BigInt()))
			}
		}

		metrics = append(metrics, rawMetrics{
			solution:   sol,
			reputation: sdkmath.LegacyNewDec(int64(solver.ReputationScore)),
			gas:        sdkmath.LegacyNewDec(int64(sol.GasEstimate)),
			stake:      sdkmath.LegacyNewDecFromBigInt(solver.StakedAmount.Amount.BigInt()),
			tip:        tipAmount,
		})
	}

	if len(metrics) == 0 {
		return nil, fmt.Errorf("no active solvers among submitted solutions")
	}

	// Single active solver after filtering
	if len(metrics) == 1 {
		return metrics[0].solution, nil
	}

	// ── Fairness Engine · best execution ────────────────────────────────────
	// For SWAP intents, the winner is the solver who delivers the most output to
	// the USER, not the one with the best reputation/tip. Solvers declare their
	// output in ExpectedOutcome and are bound to it at settlement (a winner who
	// under-delivers fails verification and is slashed), so the auction can't be
	// gamed by over-declaring. Reputation and stake are deterministic tiebreakers.
	// Non-declaring (or below-floor) solvers are treated as offering exactly the
	// intent's MinOutputAmount. Non-swap intents fall through to the multi-factor
	// score below.
	if intent, found := k.GetIntent(ctx, intentID); found && intent.IntentType == types.IntentTypeSwap {
		var swapBody types.SwapIntent
		if err := json.Unmarshal(intent.Body, &swapBody); err == nil && !swapBody.MinOutputAmount.IsNil() {
			floor := sdkmath.LegacyNewDecFromBigInt(swapBody.MinOutputAmount.BigInt())
			type bestExec struct {
				m      *rawMetrics
				output sdkmath.LegacyDec
			}
			ranked := make([]bestExec, 0, len(metrics))
			for i := range metrics {
				out := floor
				if declared, ok := metrics[i].solution.DeclaredOutput(); ok {
					d := sdkmath.LegacyNewDecFromBigInt(declared.BigInt())
					if d.GT(out) {
						out = d
					}
				}
				ranked = append(ranked, bestExec{m: &metrics[i], output: out})
			}
			sort.SliceStable(ranked, func(i, j int) bool {
				if !ranked[i].output.Equal(ranked[j].output) {
					return ranked[i].output.GT(ranked[j].output) // most output to the user wins
				}
				if !ranked[i].m.reputation.Equal(ranked[j].m.reputation) {
					return ranked[i].m.reputation.GT(ranked[j].m.reputation)
				}
				if !ranked[i].m.stake.Equal(ranked[j].m.stake) {
					return ranked[i].m.stake.GT(ranked[j].m.stake)
				}
				return ranked[i].m.solution.SolverAddr < ranked[j].m.solution.SolverAddr
			})
			return ranked[0].m.solution, nil
		}
	}

	// Find min/max for each metric to normalize to 0-100 scale
	minRep, maxRep := metrics[0].reputation, metrics[0].reputation
	minGas, maxGas := metrics[0].gas, metrics[0].gas
	minStk, maxStk := metrics[0].stake, metrics[0].stake
	minTip, maxTip := metrics[0].tip, metrics[0].tip

	for _, m := range metrics[1:] {
		if m.reputation.LT(minRep) {
			minRep = m.reputation
		}
		if m.reputation.GT(maxRep) {
			maxRep = m.reputation
		}
		if m.gas.LT(minGas) {
			minGas = m.gas
		}
		if m.gas.GT(maxGas) {
			maxGas = m.gas
		}
		if m.stake.LT(minStk) {
			minStk = m.stake
		}
		if m.stake.GT(maxStk) {
			maxStk = m.stake
		}
		if m.tip.LT(minTip) {
			minTip = m.tip
		}
		if m.tip.GT(maxTip) {
			maxTip = m.tip
		}
	}

	hundred := sdkmath.LegacyNewDec(100)

	// normalize returns a 0-100 score. If all values are the same, returns 100 for everyone.
	normalize := func(val, min, max sdkmath.LegacyDec) sdkmath.LegacyDec {
		if max.Equal(min) {
			return hundred
		}
		return val.Sub(min).Quo(max.Sub(min)).Mul(hundred)
	}

	// normalizeInverse returns a 0-100 score where LOWER values are better (for gas).
	normalizeInverse := func(val, min, max sdkmath.LegacyDec) sdkmath.LegacyDec {
		if max.Equal(min) {
			return hundred
		}
		return max.Sub(val).Quo(max.Sub(min)).Mul(hundred)
	}

	// Weights as LegacyDec
	wReputation := sdkmath.LegacyNewDecWithPrec(40, 2) // 0.40
	wGas := sdkmath.LegacyNewDecWithPrec(30, 2)        // 0.30
	wStake := sdkmath.LegacyNewDecWithPrec(20, 2)      // 0.20
	wTip := sdkmath.LegacyNewDecWithPrec(10, 2)        // 0.10

	// Score each solution
	type scoredSolution struct {
		solution *types.Solution
		score    sdkmath.LegacyDec
	}
	scored := make([]scoredSolution, 0, len(metrics))

	for _, m := range metrics {
		repScore := normalize(m.reputation, minRep, maxRep)
		gasScore := normalizeInverse(m.gas, minGas, maxGas)
		stkScore := normalize(m.stake, minStk, maxStk)
		tipScore := normalize(m.tip, minTip, maxTip)

		totalScore := repScore.Mul(wReputation).Add(gasScore.Mul(wGas)).Add(stkScore.Mul(wStake)).Add(tipScore.Mul(wTip))

		scored = append(scored, scoredSolution{solution: m.solution, score: totalScore})
	}

	// Deterministic sort: by score descending, then by solver address ascending as tiebreaker
	sort.SliceStable(scored, func(i, j int) bool {
		if !scored[i].score.Equal(scored[j].score) {
			return scored[i].score.GT(scored[j].score)
		}
		// Deterministic tiebreaker: lowest address wins
		return scored[i].solution.SolverAddr < scored[j].solution.SolverAddr
	})

	if len(scored) == 0 {
		return nil, fmt.Errorf("auction scoring produced no winner")
	}

	return scored[0].solution, nil
}
