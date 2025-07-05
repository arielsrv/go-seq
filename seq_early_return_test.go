package seq_test

import (
	"iter"
	"strconv"
	"testing"

	"github.com/arielsrv/go-seq"
	"github.com/stretchr/testify/assert"
)

// limitedCollector collects values from a sequence but stops after a specified limit.
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

// limitedCollector2 collects key-value pairs from a sequence but stops after a specified limit.
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
		assert.Len(t, result, 2)
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
		assert.Len(t, result, 2)
	})
}

func Test_YieldKeyValues_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
		keyValueSeq := seq.YieldKeyValues(m)

		// Only collect the first 2 key-value pairs
		result := limitedCollector2(keyValueSeq, 2)

		// Should only get 2 key-value pairs
		assert.Len(t, result, 2)
	})
}

func Test_SetValues_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		set := seq.NewSet(1, 2, 3, 4, 5)
		valuesSeq := set.Values()

		// Only collect the first 3 values
		result := limitedCollector(valuesSeq, 3)

		// Should only get 3 values
		assert.Len(t, result, 3)
	})
}

func Test_Concat_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		seq1 := seq.Yield(1, 2, 3)
		seq2 := seq.Yield(4, 5, 6)
		seq3 := seq.Yield(7, 8, 9)

		concatSeq := seq.Concat(seq1, seq2, seq3)

		// Only collect the first 4 items
		result := limitedCollector(concatSeq, 4)

		// Should get all items from seq1 and first item from seq2
		assert.Equal(t, []int{1, 2, 3, 4}, result)
	})
}

func Test_OfType_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		// Create a slice with mixed types
		mixedSlice := []any{1, "string", 2, 3.14, 4}

		// Create a sequence from the slice
		mixedSeq := seq.Yield(mixedSlice...)

		// Get only int values
		intSeq := seq.OfType[any, int](mixedSeq)

		// Only collect the first 2 items
		result := limitedCollector(intSeq, 2)

		// Should only get the first 2 int values
		assert.Equal(t, []int{1, 2}, result)
	})
}

func Test_Select_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		numbers := seq.Yield(1, 2, 3, 4, 5)
		doubled := seq.Select(numbers, func(n int) int {
			return n * 2
		})

		// Only collect the first 3 items
		result := limitedCollector(doubled, 3)

		// Should only get the first 3 doubled values
		assert.Equal(t, []int{2, 4, 6}, result)
	})
}

func Test_SelectKeys_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		numbers := seq.Yield(1, 2, 3, 4, 5)
		keyValueSeq := seq.SelectKeys(numbers, func(n int) string {
			return "key" + strconv.Itoa(n)
		})

		// Only collect the first 2 key-value pairs
		result := limitedCollector2(keyValueSeq, 2)

		// Should only get the first 2 key-value pairs
		assert.Len(t, result, 2)
		assert.Equal(t, "key1", result[0].Key)
		assert.Equal(t, 1, result[0].Value)
		assert.Equal(t, "key2", result[1].Key)
		assert.Equal(t, 2, result[1].Value)
	})
}

func Test_SelectMany_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		numbers := seq.Yield(1, 2, 3)
		expanded := seq.SelectMany(numbers, func(n int) iter.Seq[int] {
			return seq.Yield(n, n*10, n*100)
		})

		// Only collect the first 5 items
		result := limitedCollector(expanded, 5)

		// Should only get the first 5 expanded values
		assert.Equal(t, []int{1, 10, 100, 2, 20}, result)
	})
}

func Test_Where_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		numbers := seq.Yield(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		evens := seq.Where(numbers, func(n int) bool {
			return n%2 == 0
		})

		// Only collect the first 2 items
		result := limitedCollector(evens, 2)

		// Should only get the first 2 even values
		assert.Equal(t, []int{2, 4}, result)
	})
}

func Test_Concat2_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		seq1 := seq.YieldKeyValues(map[string]int{"a": 1, "b": 2})
		seq2 := seq.YieldKeyValues(map[string]int{"c": 3, "d": 4})

		concatSeq := seq.Concat2(seq1, seq2)

		// Only collect the first 3 key-value pairs
		result := limitedCollector2(concatSeq, 3)

		// Should only get 3 key-value pairs
		assert.Len(t, result, 3)
	})
}

func Test_Select2_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		keyValueSeq := seq.YieldKeyValues(map[string]int{"a": 1, "b": 2, "c": 3})
		transformedSeq := seq.Select2(keyValueSeq, func(k string, v int) (string, string) {
			return k + "_key", "value_" + strconv.Itoa(v)
		})

		// Only collect the first 2 key-value pairs
		result := limitedCollector2(transformedSeq, 2)

		// Should only get 2 transformed key-value pairs
		assert.Len(t, result, 2)
	})
}

func Test_SelectValues_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		keyValueSeq := seq.YieldKeyValues(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4})
		valuesSeq := seq.SelectValues(keyValueSeq, func(k string, v int) string {
			return k + "_" + strconv.Itoa(v)
		})

		// Only collect the first 2 values
		result := limitedCollector(valuesSeq, 2)

		// Should only get 2 transformed values
		assert.Len(t, result, 2)
	})
}

func Test_Where2_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		keyValueSeq := seq.YieldKeyValues(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5})
		filteredSeq := seq.Where2(keyValueSeq, func(k string, v int) bool {
			return v%2 == 0
		})

		// Only collect the first 1 key-value pair
		result := limitedCollector2(filteredSeq, 1)

		// Should only get 1 filtered key-value pair
		assert.Len(t, result, 1)
	})
}

func Test_Chunk_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		numbers := seq.Yield(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		chunks := seq.Chunk(numbers, 3)

		// Only collect the first 2 chunks
		result := limitedCollector(chunks, 2)

		// Should only get the first 2 chunks
		assert.Len(t, result, 2)
		assert.Equal(t, []int{1, 2, 3}, result[0])
		assert.Equal(t, []int{4, 5, 6}, result[1])
	})
}

func Test_Range_EarlyReturn(t *testing.T) {
	t.Run("ascending range early return", func(t *testing.T) {
		rangeSeq, _ := seq.Range(1, 10, 1)

		// Only collect the first 3 items
		result := limitedCollector(rangeSeq, 3)

		// Should only get the first 3 items
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("descending range early return", func(t *testing.T) {
		rangeSeq, _ := seq.Range(10, 1, -1)

		// Only collect the first 3 items
		result := limitedCollector(rangeSeq, 3)

		// Should only get the first 3 items
		assert.Equal(t, []int{10, 9, 8}, result)
	})
}

func Test_Skip_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		numbers := seq.Yield(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		skipped := seq.Skip(numbers, 3)

		// Only collect the first 2 items
		result := limitedCollector(skipped, 2)

		// Should only get the first 2 items after skipping 3
		assert.Equal(t, []int{4, 5}, result)
	})
}

func Test_SkipWhile_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		numbers := seq.Yield(2, 4, 6, 7, 8, 9, 10)
		skipped := seq.SkipWhile(numbers, func(i, v int) bool {
			return v%2 == 0
		})

		// Only collect the first 2 items
		result := limitedCollector(skipped, 2)

		// Should only get the first 2 items after skipping even numbers
		assert.Equal(t, []int{7, 8}, result)
	})
}

func Test_Take_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		numbers := seq.Yield(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		taken := seq.Take(numbers, 5)

		// Only collect the first 3 items
		result := limitedCollector(taken, 3)

		// Should only get the first 3 items
		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

func Test_TakeWhile_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		numbers := seq.Yield(2, 4, 6, 7, 8, 10)
		taken := seq.TakeWhile(numbers, func(i, v int) bool {
			return v%2 == 0
		})

		// Only collect the first 2 items
		result := limitedCollector(taken, 2)

		// Should only get the first 2 even numbers
		assert.Equal(t, []int{2, 4}, result)
	})
}

func Test_WithIndex_EarlyReturn(t *testing.T) {
	t.Run("early return", func(t *testing.T) {
		values := seq.Yield("a", "b", "c", "d", "e")
		indexed := seq.WithIndex(values)

		// Only collect the first 2 key-value pairs
		result := limitedCollector2(indexed, 2)

		// Should only get the first 2 indexed values
		assert.Len(t, result, 2)
		assert.Equal(t, 0, result[0].Key)
		assert.Equal(t, "a", result[0].Value)
		assert.Equal(t, 1, result[1].Key)
		assert.Equal(t, "b", result[1].Value)
	})
}
