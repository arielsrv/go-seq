package seq_test

import (
	"iter"
	"testing"

	"github.com/sectrean/go-seq"
	"github.com/stretchr/testify/assert"
)

func Test_AggregateGrouped(t *testing.T) {
	tests := []struct {
		name     string
		seq      iter.Seq2[string, int]
		initFunc func(string) int
		f        func(int, int) int
		expected map[string]int
	}{
		{
			name:     "single group",
			seq:      seq.Zip(seq.Yield("a", "a", "a"), seq.Yield(1, 2, 3)),
			initFunc: func(k string) int { return 0 },
			f:        func(acc int, v int) int { return acc + v },
			expected: map[string]int{"a": 6},
		},
		{
			name:     "multiple groups",
			seq:      seq.Zip(seq.Yield("a", "b", "a"), seq.Yield(1, 2, 3)),
			initFunc: func(k string) int { return 0 },
			f:        func(acc int, v int) int { return acc + v },
			expected: map[string]int{"a": 4, "b": 2},
		},
		{
			name:     "empty sequence",
			seq:      seq.Empty2[string, int](),
			initFunc: func(k string) int { return 0 },
			f:        func(acc int, v int) int { return acc + v },
			expected: map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := seq.AggregateGrouped(tt.seq, tt.initFunc, tt.f)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_CountGrouped(t *testing.T) {
	tests := []struct {
		name     string
		seq      iter.Seq2[string, int]
		expected map[string]int
	}{
		{
			name:     "single group",
			seq:      seq.Zip(seq.Yield("a", "a", "a"), seq.Yield(1, 2, 3)),
			expected: map[string]int{"a": 3},
		},
		{
			name:     "multiple groups",
			seq:      seq.Zip(seq.Yield("a", "b", "a"), seq.Yield(1, 2, 3)),
			expected: map[string]int{"a": 2, "b": 1},
		},
		{
			name:     "empty sequence",
			seq:      seq.Empty2[string, int](),
			expected: map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := seq.CountGrouped(tt.seq)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_CountFuncGrouped(t *testing.T) {
	tests := []struct {
		name      string
		seq       iter.Seq2[string, int]
		predicate func(string, int) bool
		expected  map[string]int
	}{
		{
			name: "count even values",
			seq:  seq.Zip(seq.Yield("a", "b", "a", "a"), seq.Yield(1, 2, 3, 4)),
			predicate: func(k string, v int) bool {
				return v%2 == 0
			},
			expected: map[string]int{"a": 1, "b": 1},
		},
		{
			name: "count all values",
			seq:  seq.Zip(seq.Yield("a", "b", "a", "b"), seq.Yield(1, 2, 3, 4)),
			predicate: func(k string, v int) bool {
				return true
			},
			expected: map[string]int{"a": 2, "b": 2},
		},
		{
			name: "empty sequence",
			seq:  seq.Empty2[string, int](),
			predicate: func(k string, v int) bool {
				return false
			},
			expected: map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := seq.CountFuncGrouped(tt.seq, tt.predicate)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_Grouped(t *testing.T) {
	tests := []struct {
		name     string
		seq      iter.Seq2[string, int]
		expected map[string][]int
	}{
		{
			name:     "single group",
			seq:      seq.Zip(seq.Yield("a", "a", "a"), seq.Yield(1, 2, 3)),
			expected: map[string][]int{"a": {1, 2, 3}},
		},
		{
			name:     "multiple groups",
			seq:      seq.Zip(seq.Yield("a", "b", "a"), seq.Yield(1, 2, 3)),
			expected: map[string][]int{"a": {1, 3}, "b": {2}},
		},
		{
			name:     "empty sequence",
			seq:      seq.Empty2[string, int](),
			expected: map[string][]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := seq.Grouped(tt.seq)
			assert.Equal(t, tt.expected, result)
		})
	}
}
