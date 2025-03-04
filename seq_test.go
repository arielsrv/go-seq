package seq_test

import (
	"fmt"
	"iter"
	"testing"

	"github.com/sectrean/go-seq"
	"github.com/sectrean/go-seq/internal/seqtest"
	"github.com/stretchr/testify/assert"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Append(tt.seq, tt.add...)
			seqtest.AssertEqual(t, tt.want, got)
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

func Test_Distinct(t *testing.T) {
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
			seq:  seq.Yield[int](),
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

func Test_Empty(t *testing.T) {
	got := seq.Empty[int]()
	seqtest.AssertEqual(t, nil, got)
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
