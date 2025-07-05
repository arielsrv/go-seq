package seq_test

import (
	"iter"
	"testing"

	"github.com/arielsrv/go-seq"
	"github.com/stretchr/testify/assert"
)

// limitedCollector collects values from a sequence but stops after a specified limit
func limitedCollector[T any](seq iter.Seq[T], limit int) []T {
	result := make([]T, 0)
	count := 0

	for v := range seq {
		result = append(result, v)
		count++
		if count >= limit {
			break
		}
	}

	return result
}

// limitedCollector2 collects key-value pairs from a sequence but stops after a specified limit
func limitedCollector2[K, V any](seq iter.Seq2[K, V], limit int) []struct {
	Key   K
	Value V
} {
	result := make([]struct {
		Key   K
		Value V
	}, 0)
	count := 0

	for k, v := range seq {
		result = append(result, struct {
			Key   K
			Value V
		}{Key: k, Value: v})
		count++
		if count >= limit {
			break
		}
	}

	return result
}

func Test_Append_EarlyReturn(t *testing.T) {
	// Test early return in the first loop (original sequence)
	t.Run("early return in original sequence", func(t *testing.T) {
		original := seq.Yield(1, 2, 3, 4, 5)
		appended := seq.Append(original, 6, 7, 8)

		// Only collect the first 2 items
		result := limitedCollector(appended, 2)

		// Should only get the first 2 items from the original sequence
		assert.Equal(t, []int{1, 2}, result)
	})

	// Test early return in the second loop (appended values)
	t.Run("early return in appended values", func(t *testing.T) {
		original := seq.Yield(1, 2)
		appended := seq.Append(original, 3, 4, 5, 6, 7)

		// Only collect the first 3 items
		result := limitedCollector(appended, 3)

		// Should get all items from original sequence and first item from appended values
		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

func Test_Prepend_EarlyReturn(t *testing.T) {
	// Test early return in the first loop (prepended values)
	t.Run("early return in prepended values", func(t *testing.T) {
		original := seq.Yield(3, 4, 5)
		prepended := seq.Prepend(original, 1, 2)

		// Only collect the first 1 item
		result := limitedCollector(prepended, 1)

		// Should only get the first item from the prepended values
		assert.Equal(t, []int{1}, result)
	})

	// Test early return in the second loop (original sequence)
	t.Run("early return in original sequence", func(t *testing.T) {
		original := seq.Yield(3, 4, 5, 6, 7)
		prepended := seq.Prepend(original, 1, 2)

		// Only collect the first 3 items
		result := limitedCollector(prepended, 3)

		// Should get all prepended values and first item from original sequence
		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

func Test_Repeat_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		repeated := seq.Repeat("test", 10)

		// Only collect the first 3 items
		result := limitedCollector(repeated, 3)

		// Should only get the first 3 repetitions
		assert.Equal(t, []string{"test", "test", "test"}, result)
	})
}

func Test_YieldBackwards_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		backwards := seq.YieldBackwards(1, 2, 3, 4, 5)

		// Only collect the first 2 items
		result := limitedCollector(backwards, 2)

		// Should only get the first 2 items in reverse order
		assert.Equal(t, []int{5, 4}, result)
	})
}

func Test_YieldChan_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		ch := make(chan int, 5)
		ch <- 1
		ch <- 2
		ch <- 3
		ch <- 4
		ch <- 5
		close(ch)

		chanSeq := seq.YieldChan(ch)

		// Only collect the first 3 items
		result := limitedCollector(chanSeq, 3)

		// Should only get the first 3 items from the channel
		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

func Test_Keys_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
		keyValueSeq := seq.YieldKeyValues(m)
		keysSeq := seq.Keys(keyValueSeq)

		// Only collect the first 2 keys
		result := limitedCollector(keysSeq, 2)

		// Should only get 2 keys
		assert.Equal(t, 2, len(result))
	})
}

func Test_Values_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
		keyValueSeq := seq.YieldKeyValues(m)
		valuesSeq := seq.Values(keyValueSeq)

		// Only collect the first 2 values
		result := limitedCollector(valuesSeq, 2)

		// Should only get 2 values
		assert.Equal(t, 2, len(result))
	})
}

func Test_YieldKeyValues_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
		keyValueSeq := seq.YieldKeyValues(m)

		// Only collect the first 2 key-value pairs
		result := limitedCollector2(keyValueSeq, 2)

		// Should only get 2 key-value pairs
		assert.Equal(t, 2, len(result))
	})
}

func Test_SetValues_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		set := seq.NewSet(1, 2, 3, 4, 5)
		valuesSeq := set.Values()

		// Only collect the first 3 values
		result := limitedCollector(valuesSeq, 3)

		// Should only get 3 values
		assert.Equal(t, 3, len(result))
	})
}
