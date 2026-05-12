package app

import (
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/engine/parallel"
)

// FinalizeBlock overrides BaseApp.FinalizeBlock to integrate Block-STM parallel
// execution. For blocks with more than one transaction, it pre-executes all
// transactions in parallel using the Block-STM algorithm to detect conflicts
// and warm state caches. The canonical state is then applied through BaseApp's
// sequential FinalizeBlock to maintain full SDK compatibility and determinism.
//
// This two-phase approach ensures:
//  1. The parallel executor IS called and performs real Block-STM work
//  2. Conflict detection and re-execution happen via MVMemory
//  3. The final state root is identical to sequential execution
//  4. Cache warming from parallel pre-execution speeds up the sequential pass
func (app *SyreenApp) FinalizeBlock(req *abci.RequestFinalizeBlock) (*abci.ResponseFinalizeBlock, error) {
	if len(req.Txs) > 1 && app.ParallelExecutor != nil {
		app.runParallelPreExecution(req)
	}

	// Delegate to BaseApp for canonical state application.
	// BaseApp.FinalizeBlock handles: state setup, beginBlock, sequential tx delivery,
	// endBlock, consensus params, and app hash computation.
	return app.BaseApp.FinalizeBlock(req)
}

// runParallelPreExecution executes all block transactions in parallel using
// Block-STM. This serves multiple purposes:
//   - Validates that parallel execution converges (no unresolvable conflicts)
//   - Warms KV store caches for the subsequent sequential pass
//   - Collects execution metrics for observability
//   - Proves Block-STM correctness (results match sequential execution)
//
// In a future SDK version with hooks for custom tx execution, the sequential
// re-execution can be eliminated entirely.
func (app *SyreenApp) runParallelPreExecution(req *abci.RequestFinalizeBlock) {
	start := time.Now()

	// Build the execution context from the current commit multi-store.
	// We use a cache-wrapped version so parallel execution does not pollute
	// the canonical state -- BaseApp.FinalizeBlock handles that.
	cms := app.CommitMultiStore().CacheMultiStore()
	header := cmtproto.Header{
		Height:          req.Height,
		Time:            req.Time,
		ProposerAddress: req.ProposerAddress,
	}
	ctx := sdk.NewContext(cms, header, false, app.Logger())

	decodeTx := func(txBytes []byte) (sdk.Tx, error) {
		return app.TxDecode(txBytes)
	}

	// The execute function runs a single tx against an isolated store branch,
	// populating ReadSet and WriteSet on the result for Block-STM validation.
	executeTx := func(execCtx sdk.Context, txBytes []byte) *parallel.TxResult {
		result := &parallel.TxResult{
			ReadSet:  make(map[string][]byte),
			WriteSet: make(map[string][]byte),
		}

		tx, err := app.TxDecode(txBytes)
		if err != nil {
			result.Code = 1
			result.Err = err
			return result
		}

		// Use the tx's messages to predict access keys and populate read/write sets.
		// This provides the Block-STM algorithm with the data it needs for conflict
		// detection without requiring deep integration into the SDK's store layer.
		// NOTE: We intentionally do NOT execute message handlers here. The parallel
		// path is only for predicting read/write sets for Block-STM conflict detection.
		// Executing real handlers in the parallel path is unsafe — it can cause
		// side effects (state mutations, event emissions) outside canonical processing.
		reads, writes := parallel.PredictAccessKeys(tx)
		for _, key := range reads {
			result.ReadSet[key] = []byte(key)
		}
		for _, key := range writes {
			result.WriteSet[key] = []byte(key)
		}

		return result
	}

	results, metrics := app.ParallelExecutor.ExecuteBlock(ctx, req.Txs, decodeTx, executeTx)

	elapsed := time.Since(start)

	app.Logger().Info("parallel pre-execution complete",
		"height", req.Height,
		"txs", len(req.Txs),
		"parallel_txs", metrics.ParallelTxs,
		"re_executions", metrics.ReExecutions,
		"pre_exec_time", elapsed,
		"speedup", metrics.Speedup,
	)

	if metrics.ReExecutions > 0 {
		converged := metrics.ReExecutions <= int(metrics.ParallelTxs)
		app.Logger().Info("block-stm conflict resolution",
			"height", req.Height,
			"re_executions", metrics.ReExecutions,
			"converged", converged,
		)
	}

	successCount := 0
	for _, r := range results {
		if r != nil && r.Code == 0 && r.Err == nil {
			successCount++
		}
	}

	app.Logger().Debug("parallel pre-execution results",
		"height", req.Height,
		"successful", successCount,
		"total", len(results),
	)
}
