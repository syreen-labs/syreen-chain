package parallel

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMVMemory_WriteAndRead(t *testing.T) {
	mv := NewMVMemory()

	// Write a value at tx 0
	mv.Write(0, "key1", []byte("value_from_tx0"))

	// Read from tx 1 should see tx 0's write
	val, writerIdx, found := mv.Read(1, "key1")
	require.True(t, found)
	require.Equal(t, 0, writerIdx)
	require.Equal(t, []byte("value_from_tx0"), val)

	// Read from tx 0 should NOT see its own write (reads from lower-indexed only)
	val, writerIdx, found = mv.Read(0, "key1")
	require.False(t, found)
	require.Equal(t, -1, writerIdx)
	require.Nil(t, val)
}

func TestMVMemory_ReadFromLowerIndex(t *testing.T) {
	mv := NewMVMemory()

	// Tx 0 writes val_0, tx 2 writes val_2, tx 5 writes val_5
	mv.Write(0, "key", []byte("val_0"))
	mv.Write(2, "key", []byte("val_2"))
	mv.Write(5, "key", []byte("val_5"))

	// Tx 1 reads -> sees tx 0's write
	val, writerIdx, found := mv.Read(1, "key")
	require.True(t, found)
	require.Equal(t, 0, writerIdx)
	require.Equal(t, []byte("val_0"), val)

	// Tx 3 reads -> sees tx 2's write (highest below 3)
	val, writerIdx, found = mv.Read(3, "key")
	require.True(t, found)
	require.Equal(t, 2, writerIdx)
	require.Equal(t, []byte("val_2"), val)

	// Tx 6 reads -> sees tx 5's write (highest below 6)
	val, writerIdx, found = mv.Read(6, "key")
	require.True(t, found)
	require.Equal(t, 5, writerIdx)
	require.Equal(t, []byte("val_5"), val)

	// Tx 0 reads -> no lower-indexed writer
	val, writerIdx, found = mv.Read(0, "key")
	require.False(t, found)
	require.Equal(t, -1, writerIdx)
}

func TestMVMemory_MarkDeleted(t *testing.T) {
	mv := NewMVMemory()

	mv.Write(0, "key", []byte("exists"))
	mv.MarkDeleted(1, "key")

	// Tx 2 reads -> sees the deletion by tx 1
	val, writerIdx, found := mv.Read(2, "key")
	require.True(t, found) // found=true because there IS an entry
	require.Equal(t, 1, writerIdx)
	require.Nil(t, val) // nil value means deleted

	// Tx 1 reads -> sees tx 0's write (not its own deletion)
	val, writerIdx, found = mv.Read(1, "key")
	require.True(t, found)
	require.Equal(t, 0, writerIdx)
	require.Equal(t, []byte("exists"), val)
}

func TestMVMemory_ValidateReadSet_Valid(t *testing.T) {
	mv := NewMVMemory()

	// Tx 0 writes key_a
	mv.Write(0, "key_a", []byte("v0"))
	// Tx 1 writes key_b
	mv.Write(1, "key_b", []byte("v1"))

	// Tx 2's read set: read key_a from tx 0, read key_b from tx 1
	readSet := map[string]int{
		"key_a": 0,
		"key_b": 1,
	}

	valid := mv.ValidateReadSet(2, readSet)
	require.True(t, valid)
}

func TestMVMemory_ValidateReadSet_Invalid(t *testing.T) {
	mv := NewMVMemory()

	// Tx 0 writes key_a
	mv.Write(0, "key_a", []byte("v0"))

	// Tx 2's read set recorded that it read key_a from tx 0
	readSet := map[string]int{
		"key_a": 0,
	}

	// Initially valid
	valid := mv.ValidateReadSet(2, readSet)
	require.True(t, valid)

	// Now tx 1 also writes key_a -> tx 2's read is invalidated
	// because the latest writer below tx 2 is now tx 1, not tx 0
	mv.Write(1, "key_a", []byte("v1"))

	valid = mv.ValidateReadSet(2, readSet)
	require.False(t, valid)
}

func TestMVMemory_ValidateReadSet_BaseStateRead(t *testing.T) {
	mv := NewMVMemory()

	// Tx 0 read key_x from base state (no writer, writerIdx = -1)
	readSet := map[string]int{
		"key_x": -1,
	}

	// Valid: no one has written key_x
	valid := mv.ValidateReadSet(0, readSet)
	require.True(t, valid)

	// Now some tx writes key_x -> tx 0 should see invalidation
	// But wait, tx 0's readSet is checking for txIndex=0, which means
	// we look for writers below tx 0 -- there are none. Still valid.
	mv.Write(5, "key_x", []byte("written"))
	valid = mv.ValidateReadSet(0, readSet)
	require.True(t, valid) // writer 5 is NOT below tx 0

	// But for tx 6, reading key_x from base state is now invalid
	readSet6 := map[string]int{
		"key_x": -1,
	}
	valid = mv.ValidateReadSet(6, readSet6)
	require.False(t, valid) // latest writer below 6 is tx 5, not -1
}

func TestMVMemory_ClearWritesForTx(t *testing.T) {
	mv := NewMVMemory()

	mv.Write(0, "key_a", []byte("v0"))
	mv.Write(1, "key_a", []byte("v1"))
	mv.Write(1, "key_b", []byte("v1b"))

	// Clear tx 1's writes
	mv.ClearWritesForTx(1)

	// Tx 2 reading key_a should now see tx 0 (not tx 1)
	val, writerIdx, found := mv.Read(2, "key_a")
	require.True(t, found)
	require.Equal(t, 0, writerIdx)
	require.Equal(t, []byte("v0"), val)

	// Tx 2 reading key_b should see nothing (tx 1 was the only writer)
	_, writerIdx, found = mv.Read(2, "key_b")
	require.False(t, found)
	require.Equal(t, -1, writerIdx)
}

func TestMVMemory_OverwriteSameTx(t *testing.T) {
	mv := NewMVMemory()

	// Tx 0 writes twice to the same key
	mv.Write(0, "key", []byte("first"))
	mv.Write(0, "key", []byte("second"))

	// Reader should see the latest value
	val, writerIdx, found := mv.Read(1, "key")
	require.True(t, found)
	require.Equal(t, 0, writerIdx)
	require.Equal(t, []byte("second"), val)
}

func TestMVMemory_ConcurrentAccess(t *testing.T) {
	mv := NewMVMemory()
	var wg sync.WaitGroup

	// Concurrent writes from different tx indices
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(txIdx int) {
			defer wg.Done()
			key := fmt.Sprintf("key_%d", txIdx%10)
			mv.Write(txIdx, key, []byte(fmt.Sprintf("val_%d", txIdx)))
		}(i)
	}
	wg.Wait()

	// Concurrent reads
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(txIdx int) {
			defer wg.Done()
			key := fmt.Sprintf("key_%d", txIdx%10)
			mv.Read(txIdx, key)
		}(i)
	}
	wg.Wait()
}

func TestBlockSTM_ConflictDetection(t *testing.T) {
	// Simulate two txs writing the same key, causing tx 1 to need re-execution.
	mv := NewMVMemory()

	// Tx 0 writes key_shared
	mv.Write(0, "key_shared", []byte("val_tx0"))

	// Tx 1 read key_shared from base state (before tx 0 wrote it, in a parallel scenario)
	readSetTx1 := map[string]int{
		"key_shared": -1, // tx 1 saw base state (no writer)
	}

	// Tx 1's read set should now be invalid because tx 0 (lower-indexed) wrote key_shared
	valid := mv.ValidateReadSet(1, readSetTx1)
	require.False(t, valid, "tx 1's read set should be invalid after tx 0 wrote the key")

	// After re-execution, tx 1 would read from tx 0
	readSetTx1Fixed := map[string]int{
		"key_shared": 0, // now sees tx 0's write
	}
	valid = mv.ValidateReadSet(1, readSetTx1Fixed)
	require.True(t, valid, "tx 1's re-executed read set should be valid")
}

func TestBlockSTM_DeterministicResults(t *testing.T) {
	// Run the same set of operations multiple times and verify deterministic outcomes.
	for trial := 0; trial < 10; trial++ {
		mv := NewMVMemory()

		// 5 txs, each writing a unique key and reading a shared key
		for i := 0; i < 5; i++ {
			mv.Write(i, fmt.Sprintf("unique_%d", i), []byte(fmt.Sprintf("val_%d", i)))
			mv.Write(i, "shared", []byte(fmt.Sprintf("shared_val_%d", i)))
		}

		// Validate that reads are deterministic
		for i := 0; i < 5; i++ {
			// Each tx reads "shared" and should see the write from the highest
			// lower-indexed tx
			val, writerIdx, found := mv.Read(i, "shared")
			if i == 0 {
				require.False(t, found, "tx 0 should not see any lower writer")
			} else {
				require.True(t, found)
				require.Equal(t, i-1, writerIdx)
				require.Equal(t, []byte(fmt.Sprintf("shared_val_%d", i-1)), val)
			}
		}

		// Validate read sets
		for i := 1; i < 5; i++ {
			readSet := map[string]int{
				"shared": i - 1,
			}
			valid := mv.ValidateReadSet(i, readSet)
			require.True(t, valid, "trial %d, tx %d should have valid read set", trial, i)
		}
	}
}
