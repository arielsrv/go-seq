package seq_test

import (
	"fmt"
	"iter"
	"testing"

	"github.com/arielsrv/go-seq"
	"github.com/arielsrv/go-seq/internal/seqtest"
	"github.com/stretchr/testify/assert"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func cmpStringLen(a, b string) int      { return len(a) - len(b) }
func double[V any](v V) iter.Seq[V]     { return seq.Yield(v, v) }
func isEven(v int) bool                 { return v%2 == 0 }
func isEqual[T comparable](a, b T) bool { return a == b }
func isAbsEqual(a, b int) bool          { return abs(a) == abs(b) }
func toString[V any](v V) string        { return fmt.Sprint(v) }

func Test_Aggregate(t *testing.T) {
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
			seq:  seq.Yield[int](),
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

func Test_All(t *testing.T) {
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
			seq:  seq.Yield[int](),
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

func Test_Any(t *testing.T) {
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
			seq:  seq.Yield[int](),
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

func Test_AnyFunc(t *testing.T) {
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
			seq:  seq.Yield[int](),
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

func Test_Append(t *testing.T) {
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
			seq:  seq.Yield[int](),
			add:  []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			name: "empty values",
			seq:  seq.Yield(1, 2, 3),
			add:  []int{},
			want: []int{1, 2, 3},
		},
		{
			name: "single value to add",
			seq:  seq.Yield(1, 2, 3),
			add:  []int{4},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "multiple values to empty",
			seq:  seq.Yield[int](),
			add:  []int{1, 2, 3, 4, 5},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "empty to empty",
			seq:  seq.Yield[int](),
			add:  []int{},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Append(tt.seq, tt.add...)
			seqtest.AssertEqual(t, tt.want, got)
		})
	}
}

func Test_Average(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		want float64
	}{
		{
			name: "positive integers",
			seq:  seq.Yield(1, 2, 3, 4),
			want: 2.5,
		},
		{
			name: "single value",
			seq:  seq.Yield(5),
			want: 5.0,
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[int](),
			want: 0.0,
		},
		{
			name: "negative values",
			seq:  seq.Yield(-2, -1, 0, 1, 2),
			want: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Average(tt.seq)
			assert.InDelta(t, tt.want, got, 0.000001)
		})
	}

	// Test with float values
	floatTests := []struct {
		name string
		seq  iter.Seq[float64]
		want float64
	}{
		{
			name: "decimal values",
			seq:  seq.Yield(1.5, 2.5, 3.5),
			want: 2.5,
		},
		{
			name: "mixed integers and decimals",
			seq:  seq.Yield(1.0, 2.0, 3.5, 4.5),
			want: 2.75,
		},
	}

	for _, tt := range floatTests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Average(tt.seq)
			assert.InDelta(t, tt.want, got, 0.000001)
		})
	}
}

func Test_Chunk(t *testing.T) {
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
			seq:  seq.Yield[int](),
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

func Test_Concat(t *testing.T) {
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
			seqs: []iter.Seq[int]{seq.Yield(1, 2), seq.Yield[int](), seq.Yield(3, 4)},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "all empty",
			seqs: []iter.Seq[int]{seq.Yield[int](), seq.Yield[int]()},
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

func Test_Contains(t *testing.T) {
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
			seq:  seq.Yield[int](),
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

func Test_Count(t *testing.T) {
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
			seq:  seq.Yield[int](),
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

func Test_CountFunc(t *testing.T) {
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
			seq:  seq.Yield[int](),
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

func Test_Empty(t *testing.T) {
	got := seq.Empty[int]()
	seqtest.AssertEqual(t, nil, got)
}

func Test_Equal(t *testing.T) {
	tests := []struct {
		name string
		seq1 iter.Seq[int]
		seq2 iter.Seq[int]
		want bool
	}{
		{
			name: "equal sequences",
			seq1: seq.Yield(1, 2, 3),
			seq2: seq.Yield(1, 2, 3),
			want: true,
		},
		{
			name: "different lengths",
			seq1: seq.Yield(1, 2, 3),
			seq2: seq.Yield(1, 2),
			want: false,
		},
		{
			name: "different values",
			seq1: seq.Yield(1, 2, 3),
			seq2: seq.Yield(1, 2, 4),
			want: false,
		},
		{
			name: "both empty",
			seq1: seq.Yield[int](),
			seq2: seq.Yield[int](),
			want: true,
		},
		{
			name: "one empty",
			seq1: seq.Yield(1, 2),
			seq2: seq.Yield[int](),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Equal(tt.seq1, tt.seq2)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_EqualFunc(t *testing.T) {
	tests := []struct {
		name string
		seq1 iter.Seq[int]
		seq2 iter.Seq[int]
		f    func(int, int) bool
		want bool
	}{
		{
			name: "equal sequences",
			seq1: seq.Yield(1, 2, 3),
			seq2: seq.Yield(1, 2, 3),
			f:    isEqual[int],
			want: true,
		},
		{
			name: "equal with custom comparison",
			seq1: seq.Yield(1, -2, 3),
			seq2: seq.Yield(-1, 2, -3),
			f:    isAbsEqual,
			want: true,
		},
		{
			name: "different lengths",
			seq1: seq.Yield(1, 2, 3),
			seq2: seq.Yield(1, 2),
			f:    isEqual[int],
			want: false,
		},
		{
			name: "different values",
			seq1: seq.Yield(1, 2, 3),
			seq2: seq.Yield(1, 2, 4),
			f:    isEqual[int],
			want: false,
		},
		{
			name: "both empty",
			seq1: seq.Yield[int](),
			seq2: seq.Yield[int](),
			f:    isEqual[int],
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.EqualFunc(tt.seq1, tt.seq2, tt.f)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_First(t *testing.T) {
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
			seq:  seq.Yield[int](),
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

func Test_FirstFunc(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		f    func(int) bool
		want int
		ok   bool
	}{
		{
			name: "found first",
			seq:  seq.Yield(1, 2, 3, 4),
			f:    isEven,
			want: 2,
			ok:   true,
		},
		{
			name: "not found",
			seq:  seq.Yield(1, 3, 5),
			f:    isEven,
			want: 0,
			ok:   false,
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[int](),
			f:    isEven,
			want: 0,
			ok:   false,
		},
		{
			name: "first element matches",
			seq:  seq.Yield(2, 4, 6),
			f:    isEven,
			want: 2,
			ok:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := seq.FirstFunc(tt.seq, tt.f)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func Test_Last(t *testing.T) {
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
			seq:  seq.Yield[int](),
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

func Test_LastFunc(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		f    func(int) bool
		want int
		ok   bool
	}{
		{
			name: "found last",
			seq:  seq.Yield(1, 2, 3, 4, 5, 6),
			f:    isEven,
			want: 6,
			ok:   true,
		},
		{
			name: "found in middle",
			seq:  seq.Yield(2, 4, 5, 7, 9),
			f:    isEven,
			want: 4,
			ok:   true,
		},
		{
			name: "not found",
			seq:  seq.Yield(1, 3, 5),
			f:    isEven,
			want: 0,
			ok:   false,
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[int](),
			f:    isEven,
			want: 0,
			ok:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := seq.LastFunc(tt.seq, tt.f)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func Test_Max(t *testing.T) {
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
			name: "single",
			seq:  seq.Yield(-10),
			want: -10,
			ok:   true,
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[int](),
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

func Test_MaxBy(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[string]
		f    func(string) int
		want string
		ok   bool
	}{
		{
			name: "by length",
			seq:  seq.Yield("a", "bb", "ccc", "dd"),
			f:    func(s string) int { return len(s) },
			want: "ccc",
			ok:   true,
		},
		{
			name: "single value",
			seq:  seq.Yield("test"),
			f:    func(s string) int { return len(s) },
			want: "test",
			ok:   true,
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[string](),
			f:    func(s string) int { return len(s) },
			want: "",
			ok:   false,
		},
		{
			name: "same key values",
			seq:  seq.Yield("aa", "bb", "cc"),
			f:    func(s string) int { return len(s) },
			want: "aa", // first occurrence is kept
			ok:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := seq.MaxBy(tt.seq, tt.f)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func Test_MaxFunc(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[string]
		f    func(string, string) int
		want string
		ok   bool
	}{
		{
			name: "by length comparison",
			seq:  seq.Yield("a", "bb", "ccc", "dd"),
			f:    cmpStringLen,
			want: "ccc",
			ok:   true,
		},
		{
			name: "single value",
			seq:  seq.Yield("test"),
			f:    cmpStringLen,
			want: "test",
			ok:   true,
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[string](),
			f:    cmpStringLen,
			want: "",
			ok:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := seq.MaxFunc(tt.seq, tt.f)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func Test_Min(t *testing.T) {
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
			seq:  seq.Yield[int](),
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

func Test_MinBy(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[string]
		f    func(string) int
		want string
		ok   bool
	}{
		{
			name: "by length",
			seq:  seq.Yield("a", "bb", "ccc", "dd"),
			f:    func(s string) int { return len(s) },
			want: "a",
			ok:   true,
		},
		{
			name: "single value",
			seq:  seq.Yield("test"),
			f:    func(s string) int { return len(s) },
			want: "test",
			ok:   true,
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[string](),
			f:    func(s string) int { return len(s) },
			want: "",
			ok:   false,
		},
		{
			name: "same key values",
			seq:  seq.Yield("aa", "bb", "cc"),
			f:    func(s string) int { return len(s) },
			want: "aa", // first occurrence is kept
			ok:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := seq.MinBy(tt.seq, tt.f)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func Test_MinFunc(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[string]
		f    func(string, string) int
		want string
		ok   bool
	}{
		{
			name: "by length comparison",
			seq:  seq.Yield("a", "bb", "ccc", "dd"),
			f:    cmpStringLen,
			want: "a",
			ok:   true,
		},
		{
			name: "single value",
			seq:  seq.Yield("test"),
			f:    cmpStringLen,
			want: "test",
			ok:   true,
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[string](),
			f:    cmpStringLen,
			want: "",
			ok:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := seq.MinFunc(tt.seq, tt.f)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func Test_OfType(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[any]
		want []int
	}{
		{
			name: "mixed types",
			seq:  seq.Yield[any](1, "two", 3, true, 4),
			want: []int{1, 3, 4},
		},
		{
			name: "all matching",
			seq:  seq.Yield[any](1, 2, 3),
			want: []int{1, 2, 3},
		},
		{
			name: "none matching",
			seq:  seq.Yield[any]("one", "two", "three"),
			want: nil,
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[any](),
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.OfType[any, int](tt.seq)
			seqtest.AssertEqual(t, tt.want, got)
		})
	}
}

func Test_Prepend(t *testing.T) {
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
			seq:  seq.Yield[int](),
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
			seq:  seq.Yield[int](),
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

func Test_Range(t *testing.T) {
	tests := []struct {
		name    string
		start   int
		end     int
		step    int
		want    []int
		wantErr string
	}{
		{
			name:  "ascending positive step",
			start: 1,
			end:   5,
			step:  1,
			want:  []int{1, 2, 3, 4, 5},
		},
		{
			name:    "ascending negative step",
			start:   1,
			end:     5,
			step:    -1,
			wantErr: "seq.Range: step must be positive for ascending ranges",
		},
		{
			name:  "descending negative step",
			start: 5,
			end:   1,
			step:  -1,
			want:  []int{5, 4, 3, 2, 1},
		},
		{
			name:    "descending positive step",
			start:   5,
			end:     1,
			step:    1,
			wantErr: "seq.Range: step must be negative for descending ranges",
		},
		{
			name:  "start equals end",
			start: 5,
			end:   5,
			step:  1,
			want:  []int{5},
		},
		{
			name:    "zero step",
			start:   1,
			end:     5,
			step:    0,
			wantErr: "seq.Range: step must be non-zero",
		},
		{
			name:  "step past end ascending",
			start: 1,
			end:   5,
			step:  3,
			want:  []int{1, 4},
		},
		{
			name:  "step past end descending",
			start: 5,
			end:   1,
			step:  -3,
			want:  []int{5, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := seq.Range(tt.start, tt.end, tt.step)
			seqtest.AssertEqual(t, tt.want, got)

			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_Range_Floats(t *testing.T) {
	tests := []struct {
		name    string
		start   float64
		end     float64
		step    float64
		want    []float64
		wantErr string
	}{
		{
			name:  "ascending positive step",
			start: 0.5,
			end:   2.5,
			step:  0.5,
			want:  []float64{0.5, 1.0, 1.5, 2.0, 2.5},
		},
		{
			name:    "ascending negative step",
			start:   0.5,
			end:     2.5,
			step:    -1,
			wantErr: "seq.Range: step must be positive for ascending ranges",
		},
		{
			name:  "descending negative step",
			start: 2.5,
			end:   0.5,
			step:  -0.5,
			want:  []float64{2.5, 2.0, 1.5, 1.0, 0.5},
		},
		{
			name:    "descending positive step",
			start:   2.5,
			end:     0.5,
			step:    1,
			wantErr: "seq.Range: step must be negative for descending ranges",
		},
		{
			name:  "step past end ascending",
			start: 0.5,
			end:   2.5,
			step:  1.5,
			want:  []float64{0.5, 2.0},
		},
		{
			name:  "step past end descending",
			start: 2.5,
			end:   0.5,
			step:  -1.5,
			want:  []float64{2.5, 1.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := seq.Range(tt.start, tt.end, tt.step)
			seqtest.AssertEqual(t, tt.want, got)

			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_Repeat(t *testing.T) {
	tests := []struct {
		name string
		val  int
		n    int
		want []int
	}{
		{
			name: "positive count",
			val:  42,
			n:    3,
			want: []int{42, 42, 42},
		},
		{
			name: "zero count",
			val:  42,
			n:    0,
			want: nil,
		},
		{
			name: "negative count",
			val:  42,
			n:    -1,
			want: nil,
		},
		{
			name: "single value",
			val:  42,
			n:    1,
			want: []int{42},
		},
		{
			name: "large count",
			val:  7,
			n:    5,
			want: []int{7, 7, 7, 7, 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Repeat(tt.val, tt.n)
			seqtest.AssertEqual(t, tt.want, got)
		})
	}
}

func Test_Select(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		f    func(int) string
		want []string
	}{
		{
			name: "multiple",
			seq:  seq.Yield(1, 2, 3),
			f:    toString[int],
			want: []string{"1", "2", "3"},
		},
		{
			name: "single",
			seq:  seq.Yield(42),
			f:    toString[int],
			want: []string{"42"},
		},
		{
			name: "empty",
			seq:  seq.Yield[int](),
			f:    toString[int],
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

func Test_SelectKeys(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		f    func(int) string
		want []seqtest.KeyValuePair[string, int]
	}{
		{
			name: "multiple",
			seq:  seq.Yield(1, 2, 3),
			f:    toString[int],
			want: []seqtest.KeyValuePair[string, int]{
				{Key: "1", Value: 1},
				{Key: "2", Value: 2},
				{Key: "3", Value: 3},
			},
		},
		{
			name: "single",
			seq:  seq.Yield(42),
			f:    toString[int],
			want: []seqtest.KeyValuePair[string, int]{
				{Key: "42", Value: 42},
			},
		},
		{
			name: "empty",
			seq:  seq.Yield[int](),
			f:    toString[int],
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

func Test_SelectMany(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		f    func(int) iter.Seq[int]
		want []int
	}{
		{
			name: "multiple",
			seq:  seq.Yield(1, 2, 3),
			f:    double[int],
			want: []int{1, 1, 2, 2, 3, 3},
		},
		{
			name: "single",
			seq:  seq.Yield(5),
			f:    double[int],
			want: []int{5, 5},
		},
		{
			name: "empty outer",
			seq:  seq.Yield[int](),
			f:    double[int],
			want: nil,
		},
		{
			name: "empty inner",
			seq:  seq.Yield(1, 2, 3),
			f: func(int) iter.Seq[int] {
				return seq.Yield[int]()
			},
			want: nil,
		},
		{
			name: "varying inner lengths",
			seq:  seq.Yield(1, 2, 3),
			f: func(v int) iter.Seq[int] {
				r, _ := seq.Range(1, v, 1)
				return r
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

func Test_Single(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		want int
		ok   bool
	}{
		{
			name: "single value",
			seq:  seq.Yield(42),
			want: 42,
			ok:   true,
		},
		{
			name: "multiple values",
			seq:  seq.Yield(1, 2, 3),
			want: 0,
			ok:   false,
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[int](),
			want: 0,
			ok:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := seq.Single(tt.seq)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func Test_SingleFunc(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		f    func(int) bool
		want int
		ok   bool
	}{
		{
			name: "single match",
			seq:  seq.Yield(1, 2, 3),
			f:    func(x int) bool { return x == 2 },
			want: 2,
			ok:   true,
		},
		{
			name: "no matches",
			seq:  seq.Yield(1, 3, 5),
			f:    isEven,
			want: 0,
			ok:   false,
		},
		{
			name: "multiple matches",
			seq:  seq.Yield(2, 4, 6),
			f:    isEven,
			want: 0,
			ok:   false,
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[int](),
			f:    isEven,
			want: 0,
			ok:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := seq.SingleFunc(tt.seq, tt.f)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func Test_Skip(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		n    int
		want []int
	}{
		{
			name: "skip some",
			seq:  seq.Yield(1, 2, 3, 4, 5),
			n:    2,
			want: []int{3, 4, 5},
		},
		{
			name: "skip none",
			seq:  seq.Yield(1, 2, 3),
			n:    0,
			want: []int{1, 2, 3},
		},
		{
			name: "skip all",
			seq:  seq.Yield(1, 2, 3),
			n:    3,
			want: nil,
		},
		{
			name: "skip more than length",
			seq:  seq.Yield(1, 2, 3),
			n:    5,
			want: nil,
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[int](),
			n:    2,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Skip(tt.seq, tt.n)
			seqtest.AssertEqual(t, tt.want, got)
		})
	}
}

func Test_SkipWhile(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		f    func(int, int) bool
		want []int
	}{
		{
			name: "skip even numbers",
			seq:  seq.Yield(2, 4, 3, 6, 8),
			f: func(i, v int) bool {
				return v%2 == 0
			},
			want: []int{3, 6, 8},
		},
		{
			name: "skip all numbers",
			seq:  seq.Yield(2, 4, 6, 8),
			f: func(i, v int) bool {
				return v%2 == 0
			},
			want: nil,
		},
		{
			name: "skip none",
			seq:  seq.Yield(1, 3, 5, 7),
			f: func(i, v int) bool {
				return v%2 == 0
			},
			want: []int{1, 3, 5, 7},
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[int](),
			f: func(i, v int) bool {
				return v%2 == 0
			},
			want: nil,
		},
		{
			name: "skip based on index",
			seq:  seq.Yield(1, 2, 3, 4, 5),
			f: func(i, v int) bool {
				return i < 2
			},
			want: []int{3, 4, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.SkipWhile(tt.seq, tt.f)
			seqtest.AssertEqual(t, tt.want, got)
		})
	}
}

func Test_Sum(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		want int
	}{
		{
			name: "positive integers",
			seq:  seq.Yield(1, 2, 3, 4),
			want: 10,
		},
		{
			name: "mixed signs",
			seq:  seq.Yield(-2, -1, 0, 1, 2),
			want: 0,
		},
		{
			name: "single value",
			seq:  seq.Yield(42),
			want: 42,
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[int](),
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Sum(tt.seq)
			assert.Equal(t, tt.want, got)
		})
	}

	// Test with float values
	floatTests := []struct {
		name string
		seq  iter.Seq[float64]
		want float64
	}{
		{
			name: "decimal values",
			seq:  seq.Yield(1.5, 2.5, 3.5),
			want: 7.5,
		},
		{
			name: "mixed integers and decimals",
			seq:  seq.Yield(1.0, 2.0, 3.5, 4.5),
			want: 11.0,
		},
	}

	for _, tt := range floatTests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Sum(tt.seq)
			assert.InDelta(t, tt.want, got, 0.000001)
		})
	}
}

func Test_Take(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		n    int
		want []int
	}{
		{
			name: "take some",
			seq:  seq.Yield(1, 2, 3, 4, 5),
			n:    3,
			want: []int{1, 2, 3},
		},
		{
			name: "take none",
			seq:  seq.Yield(1, 2, 3),
			n:    0,
			want: nil,
		},
		{
			name: "take all",
			seq:  seq.Yield(1, 2, 3),
			n:    3,
			want: []int{1, 2, 3},
		},
		{
			name: "take more than length",
			seq:  seq.Yield(1, 2, 3),
			n:    5,
			want: []int{1, 2, 3},
		},
		{
			name: "empty",
			seq:  seq.Yield[int](),
			n:    2,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Take(tt.seq, tt.n)
			seqtest.AssertEqual(t, tt.want, got)
		})
	}
}

func Test_TakeWhile(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		f    func(int, int) bool
		want []int
	}{
		{
			name: "take even numbers",
			seq:  seq.Yield(2, 4, 3, 6, 8),
			f: func(i, v int) bool {
				return v%2 == 0
			},
			want: []int{2, 4},
		},
		{
			name: "take all numbers",
			seq:  seq.Yield(2, 4, 6, 8),
			f: func(i, v int) bool {
				return v%2 == 0
			},
			want: []int{2, 4, 6, 8},
		},
		{
			name: "take none",
			seq:  seq.Yield(1, 3, 5, 7),
			f: func(i, v int) bool {
				return v%2 == 0
			},
			want: nil,
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[int](),
			f: func(i, v int) bool {
				return v%2 == 0
			},
			want: nil,
		},
		{
			name: "take based on index",
			seq:  seq.Yield(1, 2, 3, 4, 5),
			f: func(i, v int) bool {
				return i < 2
			},
			want: []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.TakeWhile(tt.seq, tt.f)
			seqtest.AssertEqual(t, tt.want, got)
		})
	}
}

func Test_ValueAt(t *testing.T) {
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
			seq:   seq.Yield[int](),
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

func Test_Where(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		f    func(int) bool
		want []int
	}{
		{
			name: "some matches",
			seq:  seq.Yield(1, 2, 3, 4, 5, 6),
			f:    isEven,
			want: []int{2, 4, 6},
		},
		{
			name: "all matches",
			seq:  seq.Yield(2, 4, 6, 8),
			f:    isEven,
			want: []int{2, 4, 6, 8},
		},
		{
			name: "no matches",
			seq:  seq.Yield(1, 3, 5, 7),
			f:    isEven,
			want: nil,
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[int](),
			f:    isEven,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Where(tt.seq, tt.f)
			seqtest.AssertEqual(t, tt.want, got)
		})
	}
}

func Test_Yield(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		want []int
	}{
		{
			name: "multiple values",
			seq:  seq.Yield(1, 2, 3, 4),
			want: []int{1, 2, 3, 4},
		},
		{
			name: "single value",
			seq:  seq.Yield(42),
			want: []int{42},
		},
		{
			name: "empty",
			seq:  seq.Yield[int](),
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seqtest.AssertEqual(t, tt.want, tt.seq)
		})
	}
}

func Test_YieldBackwards(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		want []int
	}{
		{
			name: "multiple values",
			seq:  seq.YieldBackwards(1, 2, 3, 4),
			want: []int{4, 3, 2, 1},
		},
		{
			name: "single value",
			seq:  seq.YieldBackwards(42),
			want: []int{42},
		},
		{
			name: "empty",
			seq:  seq.YieldBackwards[int](),
			want: nil,
		},
		{
			name: "two values",
			seq:  seq.YieldBackwards(1, 2),
			want: []int{2, 1},
		},
		{
			name: "large sequence",
			seq:  seq.YieldBackwards(1, 2, 3, 4, 5, 6, 7, 8, 9, 10),
			want: []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seqtest.AssertEqual(t, tt.want, tt.seq)
		})
	}
}

func Test_YieldChan(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		want []int
	}{
		{
			name: "multiple values",
			seq: func() iter.Seq[int] {
				ch := make(chan int)
				go func() {
					defer close(ch)
					ch <- 1
					ch <- 2
					ch <- 3
					ch <- 4
				}()
				return seq.YieldChan(ch)
			}(),
			want: []int{1, 2, 3, 4},
		},
		{
			name: "single value",
			seq: func() iter.Seq[int] {
				ch := make(chan int)
				go func() {
					defer close(ch)
					ch <- 42
				}()
				return seq.YieldChan(ch)
			}(),
			want: []int{42},
		},
		{
			name: "empty",
			seq: func() iter.Seq[int] {
				ch := make(chan int)
				close(ch)
				return seq.YieldChan(ch)
			}(),
			want: nil,
		},
		{
			name: "large sequence",
			seq: func() iter.Seq[int] {
				ch := make(chan int)
				go func() {
					defer close(ch)
					for i := 1; i <= 10; i++ {
						ch <- i
					}
				}()
				return seq.YieldChan(ch)
			}(),
			want: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name: "zero values",
			seq: func() iter.Seq[int] {
				ch := make(chan int, 1)
				close(ch)
				return seq.YieldChan(ch)
			}(),
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seqtest.AssertEqual(t, tt.want, tt.seq)
		})
	}
}
