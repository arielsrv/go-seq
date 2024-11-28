package seq_test

import (
	"fmt"
	"iter"
	"testing"

	"github.com/sectrean/go-seq"
	"github.com/sectrean/go-seq/internal/seqtest"
	"github.com/stretchr/testify/assert"
)

func isEven(v int) bool              { return v%2 == 0 }
func repeat2[V any](v V) iter.Seq[V] { return seq.Yield(v, v) }
func valString[V any](v V) string    { return fmt.Sprint(v) }

func TestAggregate(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		init int
		f    func(int, int) int
		want int
	}{
		{
			name: "sum",
			seq:  seq.Yield(1, 2, 3, 4),
			init: 0,
			f:    func(acc, v int) int { return acc + v },
			want: 10,
		},
		{
			name: "product",
			seq:  seq.Yield(1, 2, 3, 4),
			init: 1,
			f:    func(acc, v int) int { return acc * v },
			want: 24,
		},
		{
			name: "max",
			seq:  seq.Yield(1, 2, 3, 4),
			init: 0,
			f: func(acc, v int) int {
				if v > acc {
					return v
				}
				return acc
			},
			want: 4,
		},
		{
			name: "empty",
			seq:  seq.Empty[int](),
			init: 0,
			f:    func(acc, v int) int { return acc + v },
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Aggregate(tt.seq, tt.init, tt.f)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAll(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		f    func(int) bool
		want bool
	}{
		{
			name: "all",
			seq:  seq.Yield(2, 4, 6, 8),
			f:    isEven,
			want: true,
		},
		{
			name: "not all",
			seq:  seq.Yield(2, 3, 6, 8),
			f:    isEven,
			want: false,
		},
		{
			name: "empty",
			seq:  seq.Empty[int](),
			f:    isEven,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.All(tt.seq, tt.f)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAny(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		want bool
	}{
		{
			name: "non-empty",
			seq:  seq.Yield(1, 2, 3),
			want: true,
		},
		{
			name: "empty",
			seq:  seq.Empty[int](),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Any(tt.seq)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAnyFunc(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		f    func(int) bool
		want bool
	}{
		{
			name: "found",
			seq:  seq.Yield(1, 2, 3),
			f:    isEven,
			want: true,
		},
		{
			name: "none",
			seq:  seq.Yield(1, 3, 5),
			f:    isEven,
			want: false,
		},
		{
			name: "empty",
			seq:  seq.Empty[int](),
			f:    isEven,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.AnyFunc(tt.seq, tt.f)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAppend(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		add  []int
		want []int
	}{
		{
			name: "non-empty",
			seq:  seq.Yield(1, 2, 3),
			add:  []int{4, 5},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "empty",
			seq:  seq.Empty[int](),
			add:  []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			name: "empty values",
			seq:  seq.Yield(1, 2, 3),
			add:  []int{},
			want: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Append(tt.seq, tt.add...)
			seqtest.AssertEqual(t, tt.want, got)
		})
	}
}

func TestChunk(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		size int
		want [][]int
	}{
		{
			name: "exact chunks",
			seq:  seq.Yield(1, 2, 3, 4, 5, 6),
			size: 2,
			want: [][]int{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			name: "last chunk smaller",
			seq:  seq.Yield(1, 2, 3, 4, 5),
			size: 2,
			want: [][]int{{1, 2}, {3, 4}, {5}},
		},
		{
			name: "single chunk",
			seq:  seq.Yield(1, 2, 3),
			size: 3,
			want: [][]int{{1, 2, 3}},
		},
		{
			name: "chunk size larger than sequence",
			seq:  seq.Yield(1, 2),
			size: 5,
			want: [][]int{{1, 2}},
		},
		{
			name: "empty sequence",
			seq:  seq.Empty[int](),
			size: 3,
			want: nil,
		},
		{
			name: "chunk size 1",
			seq:  seq.Yield(1, 2, 3),
			size: 1,
			want: [][]int{{1}, {2}, {3}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Chunk(tt.seq, tt.size)
			seqtest.AssertEqual(t, tt.want, got)
		})
	}
}

func TestConcat(t *testing.T) {
	tests := []struct {
		name string
		seqs []iter.Seq[int]
		want []int
	}{
		{
			name: "two seqs",
			seqs: []iter.Seq[int]{seq.Yield(1, 2), seq.Yield(3, 4)},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "multiple seqs",
			seqs: []iter.Seq[int]{seq.Yield(1), seq.Yield(2, 2), seq.Yield(3, 3, 3)},
			want: []int{1, 2, 2, 3, 3, 3},
		},
		{
			name: "one empty",
			seqs: []iter.Seq[int]{seq.Yield(1, 2), seq.Empty[int](), seq.Yield(3, 4)},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "all empty",
			seqs: []iter.Seq[int]{seq.Empty[int](), seq.Empty[int]()},
			want: nil,
		},
		{
			name: "single sequence",
			seqs: []iter.Seq[int]{seq.Yield(1, 2, 3)},
			want: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Concat(tt.seqs...)
			seqtest.AssertEqual(t, tt.want, got)
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		val  int
		want bool
	}{
		{
			name: "found",
			seq:  seq.Yield(1, 2, 3, 4),
			val:  3,
			want: true,
		},
		{
			name: "not found",
			seq:  seq.Yield(1, 2, 3, 4),
			val:  5,
			want: false,
		},
		{
			name: "empty",
			seq:  seq.Empty[int](),
			val:  1,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Contains(tt.seq, tt.val)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCount(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		want int
	}{
		{
			name: "non-empty",
			seq:  seq.Yield(1, 2, 3, 4),
			want: 4,
		},
		{
			name: "empty",
			seq:  seq.Empty[int](),
			want: 0,
		},
		{
			name: "single",
			seq:  seq.Yield(1),
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Count(tt.seq)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCountFunc(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		f    func(int) bool
		want int
	}{
		{
			name: "some",
			seq:  seq.Yield(1, 2, 3, 4, 5, 6),
			f:    isEven,
			want: 3,
		},
		{
			name: "none",
			seq:  seq.Yield(1, 3, 5),
			f:    isEven,
			want: 0,
		},
		{
			name: "empty",
			seq:  seq.Empty[int](),
			f:    isEven,
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.CountFunc(tt.seq, tt.f)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDistinct(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		want []int
	}{
		{
			name: "distinct values",
			seq:  seq.Yield(1, 2, 2, 3, 4, 4, 5),
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "all unique",
			seq:  seq.Yield(1, 2, 3, 4, 5),
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "all duplicates",
			seq:  seq.Yield(1, 1, 1, 1),
			want: []int{1},
		},
		{
			name: "empty",
			seq:  seq.Empty[int](),
			want: nil,
		},
		{
			name: "only",
			seq:  seq.Yield(1),
			want: []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Distinct(tt.seq)
			seqtest.AssertEqual(t, tt.want, got)
		})
	}
}

func TestDistinctFunc(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		f    func(int) string
		want []int
	}{
		{
			name: "some",
			seq:  seq.Yield(1, 2, 2, 3, 4, 4, 5),
			f:    valString[int],
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "all unique",
			seq:  seq.Yield(1, 2, 3, 4, 5),
			f:    valString[int],
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "empty",
			seq:  seq.Empty[int](),
			f:    valString[int],
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.DistinctFunc(tt.seq, tt.f)
			seqtest.AssertEqual(t, tt.want, got)
		})
	}
}

func TestValueAt(t *testing.T) {
	tests := []struct {
		name  string
		seq   iter.Seq[int]
		index int
		want  int
		ok    bool
	}{
		{
			name:  "found",
			seq:   seq.Yield(1, 2, 3, 4, 5),
			index: 2,
			want:  3,
			ok:    true,
		},
		{
			name:  "first",
			seq:   seq.Yield(1, 2, 3),
			index: 0,
			want:  1,
			ok:    true,
		},
		{
			name:  "last",
			seq:   seq.Yield(1, 2, 3),
			index: 2,
			want:  3,
			ok:    true,
		},
		{
			name:  "index out of bounds",
			seq:   seq.Yield(1, 2, 3),
			index: 5,
			want:  0,
			ok:    false,
		},
		{
			name:  "empty sequence",
			seq:   seq.Empty[int](),
			index: 2,
			want:  0,
			ok:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := seq.ValueAt(tt.seq, tt.index)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.ok, ok)
		})
	}

	t.Run("panic on negative index", func(t *testing.T) {
		assert.Panics(t, func() {
			seq.ValueAt(seq.Yield(1, 2, 3), -1)
		})
	})
}

func TestEmpty(t *testing.T) {
	got := seq.Empty[int]()
	seqtest.AssertEqual(t, nil, got)
}

func TestFirst(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		want int
		ok   bool
	}{
		{
			name: "multiple",
			seq:  seq.Yield(1, 2, 3),
			want: 1,
			ok:   true,
		},
		{
			name: "single",
			seq:  seq.Yield(5),
			want: 5,
			ok:   true,
		},
		{
			name: "empty",
			seq:  seq.Empty[int](),
			want: 0,
			ok:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := seq.First(tt.seq)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func TestLast(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		want int
		ok   bool
	}{
		{
			name: "multiple",
			seq:  seq.Yield(1, 2, 3),
			want: 3,
			ok:   true,
		},
		{
			name: "single",
			seq:  seq.Yield(5),
			want: 5,
			ok:   true,
		},
		{
			name: "empty",
			seq:  seq.Empty[int](),
			want: 0,
			ok:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := seq.Last(tt.seq)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		want int
		ok   bool
	}{
		{
			name: "multiple",
			seq:  seq.Yield(1, 5, 3, 4, 2),
			want: 5,
			ok:   true,
		},
		{
			name: "single",
			seq:  seq.Yield(42),
			want: 42,
			ok:   true,
		},
		{
			name: "empty sequence",
			seq:  seq.Empty[int](),
			want: 0,
			ok:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := seq.Max(tt.seq)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		want int
		ok   bool
	}{
		{
			name: "multiple",
			seq:  seq.Yield(1, 5, 3, 4, 2),
			want: 1,
			ok:   true,
		},
		{
			name: "single",
			seq:  seq.Yield(42),
			want: 42,
			ok:   true,
		},
		{
			name: "empty",
			seq:  seq.Empty[int](),
			want: 0,
			ok:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := seq.Min(tt.seq)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func TestPrepend(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		vals []int
		want []int
	}{
		{
			name: "non-empty",
			seq:  seq.Yield(4, 5, 6),
			vals: []int{1, 2, 3},
			want: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name: "empty",
			seq:  seq.Empty[int](),
			vals: []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			name: "empty values",
			seq:  seq.Yield(1, 2, 3),
			vals: []int{},
			want: []int{1, 2, 3},
		},
		{
			name: "single",
			seq:  seq.Yield(2, 3),
			vals: []int{1},
			want: []int{1, 2, 3},
		},
		{
			name: "both empty",
			seq:  seq.Empty[int](),
			vals: []int{},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Prepend(tt.seq, tt.vals...)
			seqtest.AssertEqual(t, tt.want, got)
		})
	}
}

func TestRange(t *testing.T) {
	tests := []struct {
		name  string
		start int
		end   int
		step  int
		want  []int
	}{
		{
			name:  "step 1",
			start: 1,
			end:   5,
			step:  1,
			want:  []int{1, 2, 3, 4, 5},
		},
		{
			name:  "step 2",
			start: 1,
			end:   5,
			step:  2,
			want:  []int{1, 3, 5},
		},
		{
			name:  "start equals end",
			start: 3,
			end:   3,
			step:  1,
			want:  []int{3},
		},
		{
			name:  "start greater than end",
			start: 5,
			end:   1,
			step:  1,
			want:  nil,
		},
		{
			name:  "larger step",
			start: 0,
			end:   10,
			step:  3,
			want:  []int{0, 3, 6, 9},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Range(tt.start, tt.end, tt.step)
			seqtest.AssertEqual(t, tt.want, got)
		})
	}

	t.Run("panic on non-positive step", func(t *testing.T) {
		assert.Panics(t, func() {
			seq.Range(1, 5, 0)
		})
		assert.Panics(t, func() {
			seq.Range(1, 5, -1)
		})
	})
}

func TestRepeat(t *testing.T) {
	tests := []struct {
		name  string
		val   int
		count int
		want  []int
	}{
		{
			name:  "multiple",
			val:   42,
			count: 3,
			want:  []int{42, 42, 42},
		},
		{
			name:  "single",
			val:   1,
			count: 1,
			want:  []int{1},
		},
		{
			name:  "zero count",
			val:   5,
			count: 0,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Repeat(tt.val, tt.count)
			seqtest.AssertEqual(t, tt.want, got)
		})
	}

	t.Run("panic on negative count", func(t *testing.T) {
		assert.Panics(t, func() {
			seq.Repeat(1, -1)
		})
	})
}

func TestSelect(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		f    func(int) string
		want []string
	}{
		{
			name: "multiple",
			seq:  seq.Yield(1, 2, 3),
			f:    valString[int],
			want: []string{"1", "2", "3"},
		},
		{
			name: "single",
			seq:  seq.Yield(42),
			f:    valString[int],
			want: []string{"42"},
		},
		{
			name: "empty",
			seq:  seq.Empty[int](),
			f:    valString[int],
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Select(tt.seq, tt.f)
			seqtest.AssertEqual(t, tt.want, got)
		})
	}
}

func TestSelectKeys(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		f    func(int) string
		want []seqtest.KeyValuePair[string, int]
	}{
		{
			name: "multiple",
			seq:  seq.Yield(1, 2, 3),
			f:    valString[int],
			want: []seqtest.KeyValuePair[string, int]{
				{Key: "1", Value: 1},
				{Key: "2", Value: 2},
				{Key: "3", Value: 3},
			},
		},
		{
			name: "single",
			seq:  seq.Yield(42),
			f:    valString[int],
			want: []seqtest.KeyValuePair[string, int]{
				{Key: "42", Value: 42},
			},
		},
		{
			name: "empty",
			seq:  seq.Empty[int](),
			f:    valString[int],
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.SelectKeys(tt.seq, tt.f)
			seqtest.AssertEqual2(t, tt.want, got)
		})
	}
}

func TestSelectMany(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		f    func(int) iter.Seq[int]
		want []int
	}{
		{
			name: "multiple",
			seq:  seq.Yield(1, 2, 3),
			f:    repeat2[int],
			want: []int{1, 1, 2, 2, 3, 3},
		},
		{
			name: "single",
			seq:  seq.Yield(5),
			f:    repeat2[int],
			want: []int{5, 5},
		},
		{
			name: "empty outer",
			seq:  seq.Empty[int](),
			f:    repeat2[int],
			want: nil,
		},
		{
			name: "empty inner",
			seq:  seq.Yield(1, 2, 3),
			f: func(int) iter.Seq[int] {
				return seq.Empty[int]()
			},
			want: nil,
		},
		{
			name: "varying inner lengths",
			seq:  seq.Yield(1, 2, 3),
			f: func(v int) iter.Seq[int] {
				return seq.Range(1, v, 1)
			},
			want: []int{1, 1, 2, 1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.SelectMany(tt.seq, tt.f)
			seqtest.AssertEqual(t, tt.want, got)
		})
	}
}
