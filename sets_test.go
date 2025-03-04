package seq_test

import (
	"iter"
	"testing"

	"github.com/sectrean/go-seq"
	"github.com/sectrean/go-seq/internal/seqtest"
)

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

func Test_DistinctKeys(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq2[string, int]
		want []seqtest.KeyValuePair[string, int]
	}{
		{
			name: "distinct keys",
			seq:  seq.SelectKeys(seq.Yield(1, 2, 3), toString[int]),
			want: []seqtest.KeyValuePair[string, int]{
				{Key: "1", Value: 1},
				{Key: "2", Value: 2},
				{Key: "3", Value: 3},
			},
		},
		{
			name: "duplicate keys",
			seq:  seq.SelectKeys(seq.Yield(1, 2, 2, 3), toString[int]),
			want: []seqtest.KeyValuePair[string, int]{
				{Key: "1", Value: 1},
				{Key: "2", Value: 2},
				{Key: "3", Value: 3},
			},
		},
		{
			name: "empty sequence",
			seq:  seq.SelectKeys(seq.Yield[int](), toString[int]),
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.DistinctKeys(tt.seq)
			seqtest.AssertEqual2(t, tt.want, got)
		})
	}
}

func Test_Except(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		set  seq.Set[int]
		want []int
	}{
		{
			name: "some excluded",
			seq:  seq.Yield(1, 2, 3, 4, 5),
			set:  seq.NewSet(2, 4),
			want: []int{1, 3, 5},
		},
		{
			name: "none excluded",
			seq:  seq.Yield(1, 2, 3),
			set:  seq.NewSet[int](),
			want: []int{1, 2, 3},
		},
		{
			name: "all excluded",
			seq:  seq.Yield(1, 2, 3),
			set:  seq.NewSet(1, 2, 3),
			want: nil,
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[int](),
			set:  seq.NewSet(1, 2, 3),
			want: nil,
		},
		{
			name: "duplicates in sequence",
			seq:  seq.Yield(1, 2, 2, 3, 3, 3),
			set:  seq.NewSet(2),
			want: []int{1, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Except(tt.seq, tt.set)
			seqtest.AssertEqual(t, tt.want, got)
		})
	}
}

func Test_ExceptKeys(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq2[string, int]
		set  seq.Set[string]
		want []seqtest.KeyValuePair[string, int]
	}{
		{
			name: "some excluded",
			seq:  seq.SelectKeys(seq.Yield(1, 2, 3), toString[int]),
			set:  seq.NewSet("2"),
			want: []seqtest.KeyValuePair[string, int]{
				{Key: "1", Value: 1},
				{Key: "3", Value: 3},
			},
		},
		{
			name: "none excluded",
			seq:  seq.SelectKeys(seq.Yield(1, 2, 3), toString[int]),
			set:  seq.NewSet[string](),
			want: []seqtest.KeyValuePair[string, int]{
				{Key: "1", Value: 1},
				{Key: "2", Value: 2},
				{Key: "3", Value: 3},
			},
		},
		{
			name: "all excluded",
			seq:  seq.SelectKeys(seq.Yield(1, 2, 3), toString[int]),
			set:  seq.NewSet("1", "2", "3"),
			want: nil,
		},
		{
			name: "empty sequence",
			seq:  seq.SelectKeys(seq.Yield[int](), toString[int]),
			set:  seq.NewSet("1", "2", "3"),
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.ExceptKeys(tt.seq, tt.set)
			seqtest.AssertEqual2(t, tt.want, got)
		})
	}
}

func Test_Intersect(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq[int]
		set  seq.Set[int]
		want []int
	}{
		{
			name: "some intersect",
			seq:  seq.Yield(1, 2, 3, 4, 5),
			set:  seq.NewSet(2, 4, 6),
			want: []int{2, 4},
		},
		{
			name: "all intersect",
			seq:  seq.Yield(1, 2, 3),
			set:  seq.NewSet(1, 2, 3),
			want: []int{1, 2, 3},
		},
		{
			name: "none intersect",
			seq:  seq.Yield(1, 2, 3),
			set:  seq.NewSet(4, 5, 6),
			want: nil,
		},
		{
			name: "empty sequence",
			seq:  seq.Yield[int](),
			set:  seq.NewSet(1, 2, 3),
			want: nil,
		},
		{
			name: "duplicates in sequence",
			seq:  seq.Yield(1, 2, 2, 3, 3, 3),
			set:  seq.NewSet(2, 3),
			want: []int{2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Intersect(tt.seq, tt.set)
			seqtest.AssertEqual(t, tt.want, got)
		})
	}
}

func Test_IntersectKeys(t *testing.T) {
	tests := []struct {
		name string
		seq  iter.Seq2[string, int]
		set  seq.Set[string]
		want []seqtest.KeyValuePair[string, int]
	}{
		{
			name: "some intersect",
			seq:  seq.SelectKeys(seq.Yield(1, 2, 3), toString[int]),
			set:  seq.NewSet("2", "3", "4"),
			want: []seqtest.KeyValuePair[string, int]{
				{Key: "2", Value: 2},
				{Key: "3", Value: 3},
			},
		},
		{
			name: "all intersect",
			seq:  seq.SelectKeys(seq.Yield(1, 2, 3), toString[int]),
			set:  seq.NewSet("1", "2", "3"),
			want: []seqtest.KeyValuePair[string, int]{
				{Key: "1", Value: 1},
				{Key: "2", Value: 2},
				{Key: "3", Value: 3},
			},
		},
		{
			name: "none intersect",
			seq:  seq.SelectKeys(seq.Yield(1, 2, 3), toString[int]),
			set:  seq.NewSet("4", "5", "6"),
			want: nil,
		},
		{
			name: "empty sequence",
			seq:  seq.SelectKeys(seq.Yield[int](), toString[int]),
			set:  seq.NewSet("1", "2", "3"),
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.IntersectKeys(tt.seq, tt.set)
			seqtest.AssertEqual2(t, tt.want, got)
		})
	}
}

func Test_Union(t *testing.T) {
	tests := []struct {
		name string
		seqs []iter.Seq[int]
		want []int
	}{
		{
			name: "disjoint sets",
			seqs: []iter.Seq[int]{
				seq.Yield(1, 2),
				seq.Yield(3, 4),
				seq.Yield(5, 6),
			},
			want: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name: "overlapping sets",
			seqs: []iter.Seq[int]{
				seq.Yield(1, 2, 3),
				seq.Yield(2, 3, 4),
				seq.Yield(3, 4, 5),
			},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "some empty sets",
			seqs: []iter.Seq[int]{
				seq.Yield(1, 2),
				seq.Yield[int](),
				seq.Yield(3, 4),
			},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "all empty sets",
			seqs: []iter.Seq[int]{
				seq.Yield[int](),
				seq.Yield[int](),
			},
			want: nil,
		},
		{
			name: "duplicates within sets",
			seqs: []iter.Seq[int]{
				seq.Yield(1, 1, 2),
				seq.Yield(2, 2, 3),
			},
			want: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Union(tt.seqs...)
			seqtest.AssertEqual(t, tt.want, got)
		})
	}
}

func Test_UnionKeys(t *testing.T) {
	tests := []struct {
		name string
		seqs []iter.Seq2[string, int]
		want []seqtest.KeyValuePair[string, int]
	}{
		{
			name: "disjoint sets",
			seqs: []iter.Seq2[string, int]{
				seq.SelectKeys(seq.Yield(1, 2), toString[int]),
				seq.SelectKeys(seq.Yield(3, 4), toString[int]),
			},
			want: []seqtest.KeyValuePair[string, int]{
				{Key: "1", Value: 1},
				{Key: "2", Value: 2},
				{Key: "3", Value: 3},
				{Key: "4", Value: 4},
			},
		},
		{
			name: "overlapping sets",
			seqs: []iter.Seq2[string, int]{
				seq.SelectKeys(seq.Yield(1, 2), toString[int]),
				seq.SelectKeys(seq.Yield(2, 3), toString[int]),
			},
			want: []seqtest.KeyValuePair[string, int]{
				{Key: "1", Value: 1},
				{Key: "2", Value: 2},
				{Key: "3", Value: 3},
			},
		},
		{
			name: "some empty sets",
			seqs: []iter.Seq2[string, int]{
				seq.SelectKeys(seq.Yield(1, 2), toString[int]),
				seq.SelectKeys(seq.Yield[int](), toString[int]),
				seq.SelectKeys(seq.Yield(3, 4), toString[int]),
			},
			want: []seqtest.KeyValuePair[string, int]{
				{Key: "1", Value: 1},
				{Key: "2", Value: 2},
				{Key: "3", Value: 3},
				{Key: "4", Value: 4},
			},
		},
		{
			name: "all empty sets",
			seqs: []iter.Seq2[string, int]{
				seq.SelectKeys(seq.Yield[int](), toString[int]),
				seq.SelectKeys(seq.Yield[int](), toString[int]),
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.UnionKeys(tt.seqs...)
			seqtest.AssertEqual2(t, tt.want, got)
		})
	}
}
