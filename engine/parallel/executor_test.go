package parallel

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSchedulerNoDeps(t *testing.T) {
	s := NewScheduler(4)

	// Tasks with nil Tx each get unique decode_err keys -- they don't conflict
	tasks := []*TxTask{
		{Index: 0},
		{Index: 1},
		{Index: 2},
	}

	groups := s.AnalyzeDependencies(tasks)
	require.NotEmpty(t, groups)
	totalTasks := 0
	for _, g := range groups {
		totalTasks += len(g.Tasks)
	}
	require.Equal(t, 3, totalTasks)
}

func TestSchedulerEmpty(t *testing.T) {
	s := NewScheduler(4)
	groups := s.AnalyzeDependencies(nil)
	require.Empty(t, groups)
}

func TestAccessSet(t *testing.T) {
	as := NewAccessSet()
	as.RecordRead("key1", 0)
	as.RecordWrite("key2")

	reads := as.GetReads()
	writes := as.GetWrites()

	require.Equal(t, uint64(0), reads["key1"])
	require.Contains(t, writes, "key2")
}

func TestBlockMetricsString(t *testing.T) {
	m := BlockMetrics{
		TotalTxs:      10,
		ParallelTxs:   8,
		SequentialTxs: 2,
		Groups:        3,
	}
	s := m.String()
	require.Contains(t, s, "txs=10")
	require.Contains(t, s, "parallel=8")
}

func TestTxFingerprint(t *testing.T) {
	fp := TxFingerprint([]byte("test tx"))
	require.Len(t, fp, 8) // 4 bytes = 8 hex chars
}

func TestMVMemoryIntegrationWithExecutor(t *testing.T) {
	// Verify that MVMemory correctly handles the Block-STM conflict resolution
	// pattern used by the executor: execute -> validate -> clear -> re-execute.
	mv := NewMVMemory()

	// Simulate initial parallel execution
	// Tx 0 writes key_a
	mv.Write(0, "key_a", []byte("tx0_val"))

	// Tx 1 writes key_a (conflict!) and key_b
	mv.Write(1, "key_a", []byte("tx1_val"))
	mv.Write(1, "key_b", []byte("tx1_b"))

	// Tx 1's read set: it read key_a from base state before tx 0 wrote
	readSetTx1 := map[string]int{"key_a": -1}

	// Validation should fail for tx 1
	require.False(t, mv.ValidateReadSet(1, readSetTx1))

	// Clear tx 1's writes (preparation for re-execution)
	mv.ClearWritesForTx(1)

	// Re-execute tx 1 with updated reads
	mv.Write(1, "key_a", []byte("tx1_val_v2"))
	mv.Write(1, "key_b", []byte("tx1_b_v2"))

	// New read set reflects reading from tx 0
	readSetTx1Fixed := map[string]int{"key_a": 0}
	require.True(t, mv.ValidateReadSet(1, readSetTx1Fixed))
}

func TestSchedulerConflictDetection(t *testing.T) {
	s := NewScheduler(4)

	// Two tasks with nil Tx (decode errors) should each get unique keys
	// and be placed in the same group since they don't conflict.
	tasks := []*TxTask{
		{Index: 0},
		{Index: 1},
	}

	groups := s.AnalyzeDependencies(tasks)
	require.Equal(t, 1, len(groups))
	require.Equal(t, 2, len(groups[0].Tasks))
}
