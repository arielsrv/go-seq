package seq

import (
	"cmp"
	"iter"
	"slices"
)

// func FromSlice[Slice ~[]E, E any](s Slice) iter.Seq2[int, E] {
// 	return slices.All(s)
// }

// func FromSliceValues[Slice ~[]E, E any](s Slice) iter.Seq[E] {
// 	return slices.Values(s)
// }

// func FromSliceBackward[Slice ~[]E, E any](s Slice) iter.Seq2[int, E] {
// 	return slices.Backward(s)
// }

func Collect[V any](seq iter.Seq[V]) []V {
	return slices.Collect(seq)
}

func Reversed[V any](seq iter.Seq[V]) []V {
	s := slices.Collect(seq)
	slices.Reverse(s)
	return s
}

func Sorted[V cmp.Ordered](seq iter.Seq[V]) []V {
	return slices.Sorted(seq)
}

func SortedFunc[V any](seq iter.Seq[V], cmp func(V, V) int) []V {
	return slices.SortedFunc(seq, cmp)
}

func SortedStableFunc[V any](seq iter.Seq[V], cmp func(V, V) int) []V {
	return slices.SortedStableFunc(seq, cmp)
}
