package parallel

import (
	"sort"
	"sync"
)

// MVMemory implements a multi-version memory structure for Block-STM.
// Each write is tagged with the transaction index that produced it.
// Reads validate against the latest write from a lower-indexed tx.
type MVMemory struct {
	mu   sync.RWMutex
	data map[string][]MVEntry // key -> sorted entries by txIndex
}

// MVEntry represents a single write to a key by a specific transaction.
type MVEntry struct {
	TxIndex int
	Value   []byte
	Deleted bool
}

// NewMVMemory creates a new multi-version memory store.
func NewMVMemory() *MVMemory {
	return &MVMemory{
		data: make(map[string][]MVEntry),
	}
}

// Write records a write by txIndex for the given key.
func (mv *MVMemory) Write(txIndex int, key string, value []byte) {
	mv.mu.Lock()
	defer mv.mu.Unlock()

	entries := mv.data[key]

	// Check if an entry for this txIndex already exists and update it
	for i, e := range entries {
		if e.TxIndex == txIndex {
			entries[i] = MVEntry{TxIndex: txIndex, Value: copyBytes(value), Deleted: false}
			return
		}
	}

	// Insert new entry maintaining sorted order by txIndex
	entry := MVEntry{TxIndex: txIndex, Value: copyBytes(value), Deleted: false}
	idx := sort.Search(len(entries), func(i int) bool {
		return entries[i].TxIndex >= txIndex
	})
	entries = append(entries, MVEntry{})
	copy(entries[idx+1:], entries[idx:])
	entries[idx] = entry
	mv.data[key] = entries
}

// MarkDeleted records a deletion by txIndex for the given key.
func (mv *MVMemory) MarkDeleted(txIndex int, key string) {
	mv.mu.Lock()
	defer mv.mu.Unlock()

	entries := mv.data[key]

	// Check if an entry for this txIndex already exists and update it
	for i, e := range entries {
		if e.TxIndex == txIndex {
			entries[i] = MVEntry{TxIndex: txIndex, Value: nil, Deleted: true}
			return
		}
	}

	// Insert new entry maintaining sorted order by txIndex
	entry := MVEntry{TxIndex: txIndex, Value: nil, Deleted: true}
	idx := sort.Search(len(entries), func(i int) bool {
		return entries[i].TxIndex >= txIndex
	})
	entries = append(entries, MVEntry{})
	copy(entries[idx+1:], entries[idx:])
	entries[idx] = entry
	mv.data[key] = entries
}

// Read returns the value written by the highest txIndex < readerTxIndex.
// Returns (value, writerTxIndex, found).
// If the found entry is a deletion, returns (nil, writerTxIndex, true).
// If no lower-indexed tx has written this key, returns (nil, -1, false),
// meaning "read from base state".
func (mv *MVMemory) Read(readerTxIndex int, key string) ([]byte, int, bool) {
	mv.mu.RLock()
	defer mv.mu.RUnlock()

	entries, ok := mv.data[key]
	if !ok || len(entries) == 0 {
		return nil, -1, false
	}

	// Find the highest txIndex < readerTxIndex using binary search.
	// We search for the insertion point of readerTxIndex, then go one step back.
	idx := sort.Search(len(entries), func(i int) bool {
		return entries[i].TxIndex >= readerTxIndex
	})

	// idx is the first entry with TxIndex >= readerTxIndex.
	// We want the entry just before that.
	if idx == 0 {
		// No entry with TxIndex < readerTxIndex
		return nil, -1, false
	}

	entry := entries[idx-1]
	if entry.Deleted {
		return nil, entry.TxIndex, true
	}
	return copyBytes(entry.Value), entry.TxIndex, true
}

// ValidateReadSet checks if all reads in a tx's read set are still valid.
// The readSet maps key -> writerTxIndex that was seen at read time.
// A read is valid if the current latest writer below txIndex is still the same
// writer that was recorded in the read set.
func (mv *MVMemory) ValidateReadSet(txIndex int, readSet map[string]int) bool {
	mv.mu.RLock()
	defer mv.mu.RUnlock()

	for key, expectedWriter := range readSet {
		entries, ok := mv.data[key]
		if !ok || len(entries) == 0 {
			// No writes exist for this key. Valid only if we expected to read
			// from base state (writerTxIndex == -1).
			if expectedWriter != -1 {
				return false
			}
			continue
		}

		// Find current latest writer below txIndex
		idx := sort.Search(len(entries), func(i int) bool {
			return entries[i].TxIndex >= txIndex
		})

		if idx == 0 {
			// No entry with TxIndex < txIndex exists now
			if expectedWriter != -1 {
				return false
			}
			continue
		}

		currentWriter := entries[idx-1].TxIndex
		if currentWriter != expectedWriter {
			return false
		}
	}
	return true
}

// ClearWritesForTx removes all writes by the given txIndex.
// This is called before re-executing a transaction to clear stale entries.
func (mv *MVMemory) ClearWritesForTx(txIndex int) {
	mv.mu.Lock()
	defer mv.mu.Unlock()

	for key, entries := range mv.data {
		filtered := entries[:0]
		for _, e := range entries {
			if e.TxIndex != txIndex {
				filtered = append(filtered, e)
			}
		}
		if len(filtered) == 0 {
			delete(mv.data, key)
		} else {
			mv.data[key] = filtered
		}
	}
}

// copyBytes returns a copy of the byte slice to avoid aliasing issues.
func copyBytes(b []byte) []byte {
	if b == nil {
		return nil
	}
	cp := make([]byte, len(b))
	copy(cp, b)
	return cp
}
