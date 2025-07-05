package seq_test

import (
	"iter"
	"testing"

	"github.com/arielsrv/go-seq"
	"github.com/stretchr/testify/assert"
)

func Test_CollectLast(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		n    int
		want []int
	}{
		{
			name: "n < len",
			seq:  seq.Yield(1, 2, 3, 4, 5),
			n:    3,
			want: []int{3, 4, 5},
		},
		{
			name: "n*2 < len",
			seq:  seq.Yield(1, 2, 3, 4, 5, 6, 7, 8),
			n:    3,
			want: []int{6, 7, 8},
		},
		{
			name: "n = 0",
			seq:  seq.Yield(1, 2, 3, 4, 5),
			n:    0,
			want: nil,
		},
		{
			name: "n = 1",
			seq:  seq.Yield(1, 2, 3, 4, 5),
			n:    1,
			want: []int{5},
		},
		{
			name: "n > len",
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
			got := seq.CollectLast(tt.seq, tt.n)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Collect(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		want []int
	}{
		{
			name: "normal sequence",
			seq:  seq.Yield(1, 2, 3, 4, 5),
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[int](),
			want: nil,
		},
		{
			name: "single value",
			seq:  seq.Yield(42),
			want: []int{42},
		},
		{
			name: "duplicate values",
			seq:  seq.Yield(1, 1, 2, 2, 3),
			want: []int{1, 1, 2, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Collect(tt.seq)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Reversed(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		want []int
	}{
		{
			name: "normal sequence",
			seq:  seq.Yield(1, 2, 3, 4, 5),
			want: []int{5, 4, 3, 2, 1},
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[int](),
			want: nil,
		},
		{
			name: "single value",
			seq:  seq.Yield(42),
			want: []int{42},
		},
		{
			name: "two values",
			seq:  seq.Yield(1, 2),
			want: []int{2, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Reversed(tt.seq)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Sorted(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		want []int
	}{
		{
			name: "unsorted sequence",
			seq:  seq.Yield(3, 1, 4, 1, 5, 9, 2, 6),
			want: []int{1, 1, 2, 3, 4, 5, 6, 9},
		},
		{
			name: "already sorted",
			seq:  seq.Yield(1, 2, 3, 4, 5),
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "reverse sorted",
			seq:  seq.Yield(5, 4, 3, 2, 1),
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[int](),
			want: nil,
		},
		{
			name: "single value",
			seq:  seq.Yield(42),
			want: []int{42},
		},
		{
			name: "duplicate values",
			seq:  seq.Yield(3, 1, 3, 1, 2),
			want: []int{1, 1, 2, 3, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Sorted(tt.seq)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_SortedBy(t *testing.T) {
	tests := []struct {
		name             string
		seq              iter.Seq[string]
		f                func(string) int
		want             []string
		useElementsMatch bool
	}{
		{
			name:             "sort by length",
			seq:              seq.Yield("cat", "dog", "elephant", "ant"),
			f:                func(s string) int { return len(s) },
			want:             []string{"ant", "cat", "dog", "elephant"},
			useElementsMatch: true,
		},
		{
			name: "sort by first character",
			seq:  seq.Yield("zebra", "apple", "banana", "cherry"),
			f:    func(s string) int { return int(s[0]) },
			want: []string{"apple", "banana", "cherry", "zebra"},
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[string](),
			f:    func(s string) int { return len(s) },
			want: nil,
		},
		{
			name: "single value",
			seq:  seq.Yield("hello"),
			f:    func(s string) int { return len(s) },
			want: []string{"hello"},
		},
		{
			name: "different lengths",
			seq:  seq.Yield("cat", "elephant", "dog"),
			f:    func(s string) int { return len(s) },
			want: []string{"cat", "dog", "elephant"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.SortedBy(tt.seq, tt.f)
			if tt.useElementsMatch {
				assert.ElementsMatch(t, tt.want, got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_SortedStableBy(t *testing.T) {
	tests := []struct {
		name             string
		seq              iter.Seq[string]
		f                func(string) int
		want             []string
		useElementsMatch bool
	}{
		{
			name:             "stable sort by length",
			seq:              seq.Yield("cat", "dog", "ant", "bat"),
			f:                func(s string) int { return len(s) },
			want:             []string{"ant", "bat", "cat", "dog"},
			useElementsMatch: true,
		},
		{
			name:             "stable sort with equal keys",
			seq:              seq.Yield("apple", "banana", "cherry", "date"),
			f:                func(s string) int { return len(s) },
			want:             []string{"date", "apple", "banana", "cherry"},
			useElementsMatch: true,
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[string](),
			f:    func(s string) int { return len(s) },
			want: nil,
		},
		{
			name: "single value",
			seq:  seq.Yield("hello"),
			f:    func(s string) int { return len(s) },
			want: []string{"hello"},
		},
		{
			name: "different lengths",
			seq:  seq.Yield("cat", "elephant", "dog"),
			f:    func(s string) int { return len(s) },
			want: []string{"cat", "dog", "elephant"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.SortedStableBy(tt.seq, tt.f)
			if tt.useElementsMatch {
				assert.ElementsMatch(t, tt.want, got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_SortedFunc(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		f    func(int, int) int
		want []int
	}{
		{
			name: "sort in ascending order",
			seq:  seq.Yield(3, 1, 4, 1, 5),
			f:    func(a, b int) int { return a - b },
			want: []int{1, 1, 3, 4, 5},
		},
		{
			name: "sort in descending order",
			seq:  seq.Yield(3, 1, 4, 1, 5),
			f:    func(a, b int) int { return b - a },
			want: []int{5, 4, 3, 1, 1},
		},
		{
			name: "sort by absolute value",
			seq:  seq.Yield(-3, 1, -4, 1, 5),
			f:    func(a, b int) int { return abs(a) - abs(b) },
			want: []int{1, 1, -3, -4, 5},
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[int](),
			f:    func(a, b int) int { return a - b },
			want: nil,
		},
		{
			name: "single value",
			seq:  seq.Yield(42),
			f:    func(a, b int) int { return a - b },
			want: []int{42},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.SortedFunc(tt.seq, tt.f)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_SortedStableFunc(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		f    func(int, int) int
		want []int
	}{
		{
			name: "stable sort in ascending order",
			seq:  seq.Yield(3, 1, 4, 1, 5),
			f:    func(a, b int) int { return a - b },
			want: []int{1, 1, 3, 4, 5},
		},
		{
			name: "stable sort in descending order",
			seq:  seq.Yield(3, 1, 4, 1, 5),
			f:    func(a, b int) int { return b - a },
			want: []int{5, 4, 3, 1, 1},
		},
		{
			name: "stable sort by absolute value",
			seq:  seq.Yield(-3, 1, -4, 1, 5),
			f:    func(a, b int) int { return abs(a) - abs(b) },
			want: []int{1, 1, -3, -4, 5},
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[int](),
			f:    func(a, b int) int { return a - b },
			want: nil,
		},
		{
			name: "single value",
			seq:  seq.Yield(42),
			f:    func(a, b int) int { return a - b },
			want: []int{42},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.SortedStableFunc(tt.seq, tt.f)
			assert.Equal(t, tt.want, got)
		})
	}
}
