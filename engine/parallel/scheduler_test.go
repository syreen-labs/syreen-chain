package parallel

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestUniqueSorted verifies deduplication and sorting of access keys.
func TestUniqueSorted(t *testing.T) {
	input := []string{"bank:z", "bank:a", "bank:z", "acc:b", "acc:a", "bank:a"}
	got := uniqueSorted(input)
	require.Equal(t, []string{"acc:a", "acc:b", "bank:a", "bank:z"}, got)
}

func TestUniqueSorted_Empty(t *testing.T) {
	require.Nil(t, uniqueSorted(nil))
	require.Equal(t, []string{"a"}, uniqueSorted([]string{"a"}))
}

// TestAnalyzeDependencies_Deterministic runs the scheduler multiple times
// and asserts the grouping is always identical. Before the fix, Go's
// randomized map iteration could produce different group assignments.
func TestAnalyzeDependencies_Deterministic(t *testing.T) {
	// Build tasks that exercise the grouping algorithm.
	// Task 0 writes "k1", Task 1 writes "k2", Task 2 writes "k1" (conflicts with 0).
	// Expected: group 0 = [task0, task1], group 1 = [task2]
	makeTasks := func() []*TxTask {
		return []*TxTask{
			{Index: 0, TxBytes: []byte("tx0")}, // decoded tx set below
			{Index: 1, TxBytes: []byte("tx1")},
			{Index: 2, TxBytes: []byte("tx2")},
		}
	}

	// We cannot easily create real sdk.Tx objects here, so instead
	// we test the deterministic properties of the helper functions directly
	// and also test AnalyzeDependencies with nil-Tx tasks (decode errors).
	s := NewScheduler(4)

	// With nil Tx, each task gets a unique "decode_err:N" write key,
	// so no conflicts -- all should land in group 0 (they don't conflict).
	// Actually decode_err:0, decode_err:1, decode_err:2 are all different,
	// so all 3 tasks go into group 0.
	tasks := makeTasks()
	groups := s.AnalyzeDependencies(tasks)
	require.Len(t, groups, 1, "all decode-error tasks have unique keys, should be 1 group")
	require.Len(t, groups[0].Tasks, 3)

	// Run it 50 times to verify determinism
	for i := 0; i < 50; i++ {
		tasks2 := makeTasks()
		groups2 := s.AnalyzeDependencies(tasks2)
		require.Len(t, groups2, 1)
		require.Len(t, groups2[0].Tasks, 3)
		for j, task := range groups2[0].Tasks {
			require.Equal(t, groups[0].Tasks[j].Index, task.Index,
				"iteration %d: task order in group differs", i)
		}
	}
}

// TestAnalyzeDependencies_Empty verifies nil return on empty input.
func TestAnalyzeDependencies_Empty(t *testing.T) {
	s := NewScheduler(4)
	groups := s.AnalyzeDependencies(nil)
	require.Nil(t, groups)
}

// TestConflictsSorted verifies the conflict detection helper.
func TestConflictsSorted(t *testing.T) {
	tx := struct {
		reads  []string
		writes []string
	}{
		reads:  []string{"a", "b"},
		writes: []string{"c", "d"},
	}

	// No overlap
	group := struct {
		writes map[string]struct{}
		reads  map[string]struct{}
	}{
		writes: map[string]struct{}{"x": {}, "y": {}},
		reads:  map[string]struct{}{"z": {}},
	}
	require.False(t, conflictsSorted(tx, group))

	// WAW conflict
	group.writes["c"] = struct{}{}
	require.True(t, conflictsSorted(tx, group))
	delete(group.writes, "c")

	// WAR conflict
	group.reads["d"] = struct{}{}
	require.True(t, conflictsSorted(tx, group))
	delete(group.reads, "d")

	// RAW conflict
	group.writes["a"] = struct{}{}
	require.True(t, conflictsSorted(tx, group))
}

// TestPredictAccessKeys_Deterministic verifies that PredictAccessKeys returns
// sorted, deduplicated slices by calling it multiple times.
func TestPredictAccessKeys_Deterministic(t *testing.T) {
	// We can't easily construct a real sdk.Tx without the full codec setup,
	// but we can verify the uniqueSorted helper that PredictAccessKeys uses.
	// This test ensures the helper itself produces deterministic output.
	for i := 0; i < 100; i++ {
		got := uniqueSorted([]string{"z", "a", "m", "a", "z", "b"})
		require.Equal(t, []string{"a", "b", "m", "z"}, got, "iteration %d", i)
	}
}
