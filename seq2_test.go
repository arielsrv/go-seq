package seq_test

import (
	"iter"
	"testing"

	"github.com/arielsrv/go-seq"
	"github.com/arielsrv/go-seq/internal/seqtest"
	"github.com/stretchr/testify/assert"
)

func Test_Concat2(t *testing.T) {
	tests := []struct {
		name string
		seqs []iter.Seq2[string, int]
		want []seqtest.KeyValuePair[string, int]
	}{
		{
			name: "two seqs",
			seqs: []iter.Seq2[string, int]{
				seq.Zip(seq.Yield("a", "b"), seq.Yield(1, 2)),
				seq.Zip(seq.Yield("c", "d"), seq.Yield(3, 4)),
			},
			want: []seqtest.KeyValuePair[string, int]{
				{Key: "a", Value: 1},
				{Key: "b", Value: 2},
				{Key: "c", Value: 3},
				{Key: "d", Value: 4},
			},
		},
		{
			name: "one empty",
			seqs: []iter.Seq2[string, int]{
				seq.Zip(seq.Yield("a"), seq.Yield(1)),
				seq.Empty2[string, int](),
				seq.Zip(seq.Yield("c"), seq.Yield(3)),
			},
			want: []seqtest.KeyValuePair[string, int]{
				{Key: "a", Value: 1},
				{Key: "c", Value: 3},
			},
		},
		{
			name: "all empty",
			seqs: []iter.Seq2[string, int]{
				seq.Empty2[string, int](),
				seq.Empty2[string, int](),
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Concat2(tt.seqs...)
			seqtest.AssertEqual2(t, tt.want, got)
		})
	}
}

func Test_ContainsKey(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq2[string, int]
		key  string
		want bool
	}{
		{
			name: "found",
			seq:  seq.Zip(seq.Yield("a", "b", "c"), seq.Yield(1, 2, 3)),
			key:  "b",
			want: true,
		},
		{
			name: "not found",
			seq:  seq.Zip(seq.Yield("a", "b"), seq.Yield(1, 2)),
			key:  "c",
			want: false,
		},
		{
			name: "empty",
			seq:  seq.Empty2[string, int](),
			key:  "a",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.ContainsKey(tt.seq, tt.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Empty2(t *testing.T) {
	got := seq.Empty2[string, int]()
	seqtest.AssertEqual2(t, nil, got)
}

func Test_Keys(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq2[string, int]
		want []string
	}{
		{
			name: "multiple",
			seq:  seq.Zip(seq.Yield("a", "b", "c"), seq.Yield(1, 2, 3)),
			want: []string{"a", "b", "c"},
		},
		{
			name: "single",
			seq:  seq.Zip(seq.Yield("x"), seq.Yield(42)),
			want: []string{"x"},
		},
		{
			name: "empty",
			seq:  seq.Empty2[string, int](),
			want: nil,
		},
		{
			name: "large sequence",
			seq:  seq.Zip(seq.Yield("a", "b", "c", "d", "e"), seq.Yield(1, 2, 3, 4, 5)),
			want: []string{"a", "b", "c", "d", "e"},
		},
		{
			name: "duplicate keys",
			seq:  seq.Zip(seq.Yield("a", "a", "b"), seq.Yield(1, 2, 3)),
			want: []string{"a", "a", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Keys(tt.seq)
			seqtest.AssertEqual(t, tt.want, got)
		})
	}
}

func Test_Keys_EdgeCases(t *testing.T) {
	t.Run("keys with large sequence", func(t *testing.T) {
		seq2 := seq.Zip(seq.Yield("a", "b", "c", "d", "e"), seq.Yield(1, 2, 3, 4, 5))
		got := seq.Keys(seq2)
		seqtest.AssertEqual(t, []string{"a", "b", "c", "d", "e"}, got)
	})
	t.Run("keys with duplicate keys", func(t *testing.T) {
		seq2 := seq.Zip(seq.Yield("a", "a", "b"), seq.Yield(1, 2, 3))
		got := seq.Keys(seq2)
		seqtest.AssertEqual(t, []string{"a", "a", "b"}, got)
	})
}

func Test_Select2(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq2[string, int]
		f    func(string, int) (string, string)
		want []seqtest.KeyValuePair[string, string]
	}{
		{
			name: "transform both",
			seq:  seq.Zip(seq.Yield("a", "b"), seq.Yield(1, 2)),
			f: func(k string, v int) (string, string) {
				return k + k, toString(v) + toString(v)
			},
			want: []seqtest.KeyValuePair[string, string]{
				{Key: "aa", Value: "11"},
				{Key: "bb", Value: "22"},
			},
		},
		{
			name: "empty",
			seq:  seq.Empty2[string, int](),
			f: func(k string, v int) (string, string) {
				return k + k, toString(v)
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Select2(tt.seq, tt.f)
			seqtest.AssertEqual2(t, tt.want, got)
		})
	}
}

func Test_SelectValues(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq2[string, int]
		f    func(string, int) string
		want []string
	}{
		{
			name: "combine key and value",
			seq:  seq.Zip(seq.Yield("a", "b"), seq.Yield(1, 2)),
			f: func(k string, v int) string {
				return k + toString(v)
			},
			want: []string{"a1", "b2"},
		},
		{
			name: "empty",
			seq:  seq.Empty2[string, int](),
			f: func(k string, v int) string {
				return k + toString(v)
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.SelectValues(tt.seq, tt.f)
			seqtest.AssertEqual(t, tt.want, got)
		})
	}
}

func Test_Where2(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq2[string, int]
		f    func(string, int) bool
		want []seqtest.KeyValuePair[string, int]
	}{
		{
			name: "some matches",
			seq:  seq.Zip(seq.Yield("a", "b", "c"), seq.Yield(1, 2, 3)),
			f: func(k string, v int) bool {
				return v%2 == 0
			},
			want: []seqtest.KeyValuePair[string, int]{
				{Key: "b", Value: 2},
			},
		},
		{
			name: "all matches",
			seq:  seq.Zip(seq.Yield("a", "b"), seq.Yield(2, 4)),
			f: func(k string, v int) bool {
				return v%2 == 0
			},
			want: []seqtest.KeyValuePair[string, int]{
				{Key: "a", Value: 2},
				{Key: "b", Value: 4},
			},
		},
		{
			name: "no matches",
			seq:  seq.Zip(seq.Yield("a", "b"), seq.Yield(1, 3)),
			f: func(k string, v int) bool {
				return v%2 == 0
			},
			want: nil,
		},
		{
			name: "empty",
			seq:  seq.Empty2[string, int](),
			f: func(k string, v int) bool {
				return true
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Where2(tt.seq, tt.f)
			seqtest.AssertEqual2(t, tt.want, got)
		})
	}
}

func Test_WithIndex(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[string]
		want []seqtest.KeyValuePair[int, string]
	}{
		{
			name: "multiple values",
			seq:  seq.Yield("a", "b", "c"),
			want: []seqtest.KeyValuePair[int, string]{
				{Key: 0, Value: "a"},
				{Key: 1, Value: "b"},
				{Key: 2, Value: "c"},
			},
		},
		{
			name: "single value",
			seq:  seq.Yield("x"),
			want: []seqtest.KeyValuePair[int, string]{
				{Key: 0, Value: "x"},
			},
		},
		{
			name: "empty",
			seq:  seq.Empty[string](),
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.WithIndex(tt.seq)
			seqtest.AssertEqual2(t, tt.want, got)
		})
	}
}

func Test_Values(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq2[string, int]
		want []int
	}{
		{
			name: "multiple",
			seq:  seq.Zip(seq.Yield("a", "b", "c"), seq.Yield(1, 2, 3)),
			want: []int{1, 2, 3},
		},
		{
			name: "single",
			seq:  seq.Zip(seq.Yield("x"), seq.Yield(42)),
			want: []int{42},
		},
		{
			name: "empty",
			seq:  seq.Empty2[string, int](),
			want: nil,
		},
		{
			name: "large sequence",
			seq:  seq.Zip(seq.Yield("a", "b", "c", "d", "e"), seq.Yield(1, 2, 3, 4, 5)),
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "duplicate keys",
			seq:  seq.Zip(seq.Yield("a", "a", "b"), seq.Yield(1, 2, 3)),
			want: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Values(tt.seq)
			seqtest.AssertEqual(t, tt.want, got)
		})
	}
}

func Test_Values_EdgeCases(t *testing.T) {
	t.Run("values with large sequence", func(t *testing.T) {
		seq2 := seq.Zip(seq.Yield("a", "b", "c", "d", "e"), seq.Yield(1, 2, 3, 4, 5))
		got := seq.Values(seq2)
		seqtest.AssertEqual(t, []int{1, 2, 3, 4, 5}, got)
	})
	t.Run("values with duplicate keys", func(t *testing.T) {
		seq2 := seq.Zip(seq.Yield("a", "a", "b"), seq.Yield(1, 2, 3))
		got := seq.Values(seq2)
		seqtest.AssertEqual(t, []int{1, 2, 3}, got)
	})
}

func Test_YieldKeyValues(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]int
		want []seqtest.KeyValuePair[string, int]
	}{
		{
			name: "multiple entries",
			m:    map[string]int{"a": 1, "b": 2, "c": 3},
			want: []seqtest.KeyValuePair[string, int]{
				{Key: "a", Value: 1},
				{Key: "b", Value: 2},
				{Key: "c", Value: 3},
			},
		},
		{
			name: "single entry",
			m:    map[string]int{"x": 42},
			want: []seqtest.KeyValuePair[string, int]{
				{Key: "x", Value: 42},
			},
		},
		{
			name: "empty map",
			m:    map[string]int{},
			want: nil,
		},
		{
			name: "large map",
			m:    map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5},
			want: []seqtest.KeyValuePair[string, int]{
				{Key: "a", Value: 1},
				{Key: "b", Value: 2},
				{Key: "c", Value: 3},
				{Key: "d", Value: 4},
				{Key: "e", Value: 5},
			},
		},
		{
			name: "nil map",
			m:    nil,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.YieldKeyValues(tt.m)
			seqtest.AssertElementsMatch2(t, tt.want, got)
		})
	}
}

func Test_YieldKeyValues_EdgeCases(t *testing.T) {
	t.Run("yield key values with large map", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}
		got := seq.YieldKeyValues(m)
		seqtest.AssertElementsMatch2(t, []seqtest.KeyValuePair[string, int]{
			{Key: "a", Value: 1},
			{Key: "b", Value: 2},
			{Key: "c", Value: 3},
			{Key: "d", Value: 4},
			{Key: "e", Value: 5},
		}, got)
	})
	t.Run("yield key values with nil map", func(t *testing.T) {
		var m map[string]int
		got := seq.YieldKeyValues(m)
		seqtest.AssertEqual2(t, nil, got)
	})
}

func Test_Zip(t *testing.T) {
	tests := []struct {
		name string
		keys iter.Seq[string]
		vals iter.Seq[int]
		want []seqtest.KeyValuePair[string, int]
	}{
		{
			name: "equal length",
			keys: seq.Yield("a", "b", "c"),
			vals: seq.Yield(1, 2, 3),
			want: []seqtest.KeyValuePair[string, int]{
				{Key: "a", Value: 1},
				{Key: "b", Value: 2},
				{Key: "c", Value: 3},
			},
		},
		{
			name: "more keys than values",
			keys: seq.Yield("a", "b", "c"),
			vals: seq.Yield(1, 2),
			want: []seqtest.KeyValuePair[string, int]{
				{Key: "a", Value: 1},
				{Key: "b", Value: 2},
			},
		},
		{
			name: "more values than keys",
			keys: seq.Yield("a", "b"),
			vals: seq.Yield(1, 2, 3),
			want: []seqtest.KeyValuePair[string, int]{
				{Key: "a", Value: 1},
				{Key: "b", Value: 2},
			},
		},
		{
			name: "empty keys",
			keys: seq.Empty[string](),
			vals: seq.Yield(1, 2, 3),
			want: nil,
		},
		{
			name: "empty values",
			keys: seq.Yield("a", "b", "c"),
			vals: seq.Empty[int](),
			want: nil,
		},
		{
			name: "both empty",
			keys: seq.Empty[string](),
			vals: seq.Empty[int](),
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Zip(tt.keys, tt.vals)
			seqtest.AssertEqual2(t, tt.want, got)
		})
	}
}
