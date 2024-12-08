package seq_test

import (
	"iter"
	"testing"

	"github.com/sectrean/go-seq"
	"github.com/stretchr/testify/assert"
)

func TestCollectLast(t *testing.T) {
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
