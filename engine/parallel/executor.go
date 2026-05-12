package parallel

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MaxReExecutions is the safety limit on re-execution rounds to prevent
// infinite loops in pathological cases.
const MaxReExecutions = 50

// Executor performs Block-STM parallel execution of transactions in a block.
type Executor struct {
	logger    log.Logger
	scheduler *Scheduler
	workers   int
	metrics   BlockMetrics
}

// NewExecutor creates a new parallel execution engine.
func NewExecutor(logger log.Logger, workers int) *Executor {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	if workers > 32 {
		workers = 32
	}
	return &Executor{
		logger:    logger.With("module", "parallel-executor"),
		scheduler: NewScheduler(workers),
		workers:   workers,
	}
}

// ExecuteBlock runs transactions using the Block-STM algorithm:
//  1. Execute ALL txs optimistically in parallel using MVMemory
//  2. Validate each tx's read set against MVMemory
//  3. Re-execute any tx whose reads were invalidated
//  4. Repeat until all txs are validated
//  5. Return results in deterministic tx-index order
func (e *Executor) ExecuteBlock(
	ctx sdk.Context,
	txs [][]byte,
	decodeTx func([]byte) (sdk.Tx, error),
	executeTx TxExecuteFn,
) ([]*TxResult, BlockMetrics) {
	start := time.Now()
	metrics := BlockMetrics{TotalTxs: len(txs)}

	if len(txs) == 0 {
		metrics.TotalTime = time.Since(start)
		return nil, metrics
	}

	// Decode all transactions upfront
	tasks := make([]*TxTask, len(txs))
	for i, txBytes := range txs {
		tx, err := decodeTx(txBytes)
		tasks[i] = &TxTask{
			Index:     i,
			TxBytes:   txBytes,
			Tx:        tx,
			DecodeErr: err,
		}
	}

	// Single transaction: skip parallelism overhead
	if len(tasks) == 1 {
		result := executeTx(ctx, tasks[0].TxBytes)
		result.Index = 0
		metrics.SequentialTxs = 1
		metrics.Groups = 1
		metrics.TotalTime = time.Since(start)
		return []*TxResult{result}, metrics
	}

	// Run Block-STM
	results := e.blockSTM(ctx, tasks, executeTx, &metrics)

	metrics.TotalTime = time.Since(start)
	if metrics.TotalTime > 0 {
		// Estimate speedup: assume sequential would take sum of all individual execution times
		seqEstimate := metrics.SequentialTime + metrics.ParallelTime
		if seqEstimate > 0 {
			metrics.Speedup = float64(seqEstimate) / float64(metrics.TotalTime)
		}
	}

	e.logger.Info("block-stm execution complete",
		"total_txs", metrics.TotalTxs,
		"parallel_txs", metrics.ParallelTxs,
		"sequential_txs", metrics.SequentialTxs,
		"re_executions", metrics.ReExecutions,
		"groups", metrics.Groups,
		"total_time", metrics.TotalTime,
	)

	return results, metrics
}

// blockSTM implements the Block-STM algorithm.
func (e *Executor) blockSTM(
	ctx sdk.Context,
	tasks []*TxTask,
	executeTx TxExecuteFn,
	metrics *BlockMetrics,
) []*TxResult {
	n := len(tasks)
	mvMem := NewMVMemory()

	// Per-tx state
	results := make([]*TxResult, n)
	readSets := make([]map[string]int, n) // key -> writerTxIndex seen at read time
	validated := make([]bool, n)
	var reExecCount int64

	// Phase 1: Execute all txs in parallel, recording reads/writes into MVMemory
	e.executeAllParallel(ctx, tasks, executeTx, mvMem, results, readSets, metrics)
	metrics.ParallelTxs = n

	// Phases 2-3: Validate and re-execute until convergence
	for round := 0; round < MaxReExecutions; round++ {
		// Reset validation flags
		for i := range validated {
			validated[i] = false
		}

		allValid := true

		// Validate in tx-index order (important for Block-STM correctness)
		for i := 0; i < n; i++ {
			if results[i] == nil {
				// Failed to execute initially; skip validation
				validated[i] = true
				continue
			}

			if readSets[i] == nil {
				// No reads recorded (e.g., decode error); consider valid
				validated[i] = true
				continue
			}

			if mvMem.ValidateReadSet(i, readSets[i]) {
				validated[i] = true
			} else {
				// Read set invalidated: a lower-indexed tx has written
				// a value that this tx read, and the writer changed.
				allValid = false
				validated[i] = false

				// Clear stale MVMemory entries for this tx and re-execute
				mvMem.ClearWritesForTx(i)
				e.reExecuteTx(ctx, tasks[i], executeTx, mvMem, results, readSets)
				atomic.AddInt64(&reExecCount, 1)
			}
		}

		if allValid {
			break
		}
	}

	metrics.ReExecutions = int(atomic.LoadInt64(&reExecCount))
	metrics.Groups = 1 // Block-STM treats the block as one group

	// Phase 4: Commit results in tx-index order by writing CacheMultiStores.
	// Only successful transactions have their state committed.
	for i := 0; i < n; i++ {
		r := results[i]
		if r != nil && r.CacheMS != nil && r.Err == nil && r.Code == 0 {
			r.CacheMS.Write()
		}
		// Clear the CacheMS reference after commit to allow GC
		if r != nil {
			r.CacheMS = nil
		}
	}

	return results
}

// executeAllParallel runs all transactions concurrently, each recording
// their read/write sets into MVMemory.
func (e *Executor) executeAllParallel(
	ctx sdk.Context,
	tasks []*TxTask,
	executeTx TxExecuteFn,
	mvMem *MVMemory,
	results []*TxResult,
	readSets []map[string]int,
	metrics *BlockMetrics,
) {
	n := len(tasks)
	var wg sync.WaitGroup
	wg.Add(n)
	sem := make(chan struct{}, e.workers)
	parStart := time.Now()

	for i := 0; i < n; i++ {
		sem <- struct{}{}
		go func(idx int) {
			defer func() {
				if r := recover(); r != nil {
					results[idx] = &TxResult{
						Index: tasks[idx].Index,
						Err:   fmt.Errorf("panic during parallel tx execution: %v", r),
						Code:  1,
					}
				}
				<-sem
				wg.Done()
			}()

			task := tasks[idx]

			// Create an isolated context with a branched store
			cacheMS := ctx.MultiStore().CacheMultiStore()
			branchCtx := ctx.WithMultiStore(cacheMS).
				WithEventManager(sdk.NewEventManager())

			// Execute the transaction
			result := executeTx(branchCtx, task.TxBytes)
			result.Index = task.Index

			// Record write set into MVMemory
			if result.WriteSet != nil {
				for key, val := range result.WriteSet {
					if val == nil {
						mvMem.MarkDeleted(idx, key)
					} else {
						mvMem.Write(idx, key, val)
					}
				}
			}

			// Record read set (key -> writer tx index)
			rs := make(map[string]int)
			if result.ReadSet != nil {
				for key := range result.ReadSet {
					_, writerIdx, found := mvMem.Read(idx, key)
					if found {
						rs[key] = writerIdx
					} else {
						rs[key] = -1 // read from base state
					}
				}
			}

			// Store the cacheMS for deferred commit after validation.
			// Do NOT call cacheMS.Write() here — Block-STM requires
			// validation to pass before committing state.
			result.CacheMS = cacheMS

			results[idx] = result
			readSets[idx] = rs
		}(i)
	}

	wg.Wait()
	metrics.ParallelTime = time.Since(parStart)
}

// reExecuteTx re-executes a single transaction after its read set was
// invalidated, recording new read/write sets into MVMemory.
func (e *Executor) reExecuteTx(
	ctx sdk.Context,
	task *TxTask,
	executeTx TxExecuteFn,
	mvMem *MVMemory,
	results []*TxResult,
	readSets []map[string]int,
) {
	idx := task.Index

	defer func() {
		if r := recover(); r != nil {
			results[idx] = &TxResult{
				Index: idx,
				Err:   fmt.Errorf("panic during re-execution: %v", r),
				Code:  1,
			}
		}
	}()

	cacheMS := ctx.MultiStore().CacheMultiStore()
	branchCtx := ctx.WithMultiStore(cacheMS).
		WithEventManager(sdk.NewEventManager())

	result := executeTx(branchCtx, task.TxBytes)
	result.Index = idx

	// Record write set into MVMemory
	if result.WriteSet != nil {
		for key, val := range result.WriteSet {
			if val == nil {
				mvMem.MarkDeleted(idx, key)
			} else {
				mvMem.Write(idx, key, val)
			}
		}
	}

	// Record read set
	rs := make(map[string]int)
	if result.ReadSet != nil {
		for key := range result.ReadSet {
			_, writerIdx, found := mvMem.Read(idx, key)
			if found {
				rs[key] = writerIdx
			} else {
				rs[key] = -1
			}
		}
	}

	// Store fresh cacheMS for deferred commit; old one is discarded.
	result.CacheMS = cacheMS

	results[idx] = result
	readSets[idx] = rs
}

// String returns a human-readable summary of block metrics.
func (m BlockMetrics) String() string {
	return fmt.Sprintf(
		"txs=%d parallel=%d sequential=%d groups=%d re_exec=%d speedup=%.2fx time=%v",
		m.TotalTxs, m.ParallelTxs, m.SequentialTxs, m.Groups, m.ReExecutions, m.Speedup, m.TotalTime,
	)
}
