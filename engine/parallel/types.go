// Package parallel implements an optimistic parallel transaction execution engine
// using a Block-STM inspired approach. Transactions are speculatively executed
// in parallel against isolated store branches, with conflict detection and
// sequential re-execution for conflicting transactions.
package parallel

import (
	"sync"
	"time"

	storetypes "cosmossdk.io/store/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TxTask represents a single transaction to be executed
type TxTask struct {
	Index    int      // Original index in the block
	TxBytes  []byte   // Raw transaction bytes
	Tx       sdk.Tx   // Decoded transaction (may be nil if decode failed)
	DecodeErr error   // Error from decoding, if any
}

// TxResult holds the result of executing a single transaction
type TxResult struct {
	Index       int
	Code        uint32
	Data        []byte
	Log         string
	GasWanted   int64
	GasUsed     int64
	Events      []abci.Event
	Err         error

	// Internal tracking for conflict detection
	ReadSet     map[string][]byte  // store key -> value hash at read time
	WriteSet    map[string][]byte  // store key -> written value

	// CacheMS holds the branched multi-store for deferred commit.
	// State is only committed after Block-STM validation passes.
	CacheMS     storetypes.CacheMultiStore
}

// AccessSet tracks all store keys read and written by a transaction
type AccessSet struct {
	mu       sync.RWMutex
	Reads    map[string]uint64  // key -> version (tx index that wrote the value)
	Writes   map[string]struct{}
}

func NewAccessSet() *AccessSet {
	return &AccessSet{
		Reads:  make(map[string]uint64),
		Writes: make(map[string]struct{}),
	}
}

func (a *AccessSet) RecordRead(key string, version uint64) {
	a.mu.Lock()
	a.Reads[key] = version
	a.mu.Unlock()
}

func (a *AccessSet) RecordWrite(key string) {
	a.mu.Lock()
	a.Writes[key] = struct{}{}
	a.mu.Unlock()
}

func (a *AccessSet) GetReads() map[string]uint64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	cp := make(map[string]uint64, len(a.Reads))
	for k, v := range a.Reads {
		cp[k] = v
	}
	return cp
}

func (a *AccessSet) GetWrites() map[string]struct{} {
	a.mu.RLock()
	defer a.mu.RUnlock()
	cp := make(map[string]struct{}, len(a.Writes))
	for k := range a.Writes {
		cp[k] = struct{}{}
	}
	return cp
}

// ExecutionGroup represents a set of transactions that can execute in parallel
type ExecutionGroup struct {
	Tasks []*TxTask
}

// BlockMetrics tracks performance metrics for parallel execution
type BlockMetrics struct {
	TotalTxs        int
	ParallelTxs     int
	SequentialTxs   int
	ReExecutions    int
	Groups          int
	ParallelTime    time.Duration
	SequentialTime  time.Duration
	TotalTime       time.Duration
	Speedup         float64
}

// TxExecuteFn is the function signature for executing a single transaction
// against a given store. It must return the result and the cache write to apply.
type TxExecuteFn func(ctx sdk.Context, txBytes []byte) *TxResult

// TrackingMultiStore wraps a CacheMultiStore to track read/write access
type TrackingMultiStore struct {
	storetypes.CacheMultiStore
	accessSet *AccessSet
	storeKey  string
}
