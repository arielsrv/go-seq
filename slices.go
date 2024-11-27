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

// Collect collects the values from the input sequence into a slice.
func Collect[V any](seq iter.Seq[V]) []V {
	return slices.Collect(seq)
}

// Reversed collects the values from the input sequence into a slice and then reverses it.
func Reversed[V any](seq iter.Seq[V]) []V {
	s := slices.Collect(seq)
	slices.Reverse(s)
	return s
}

// Sorted collects the values from the input sequence into a slice and then sorts it.
func Sorted[V cmp.Ordered](seq iter.Seq[V]) []V {
	return slices.Sorted(seq)
}

// SortedFunc collects the values from the input sequence into a slice and then sorts it using the
// given comparison function.
func SortedFunc[V any](seq iter.Seq[V], cmp func(V, V) int) []V {
	return slices.SortedFunc(seq, cmp)
}

// SortedStable collects the values from the input sequence into a slice and then sorts it using the
// given comparison function maintaining the order of equal elements.
func SortedStableFunc[V any](seq iter.Seq[V], cmp func(V, V) int) []V {
	return slices.SortedStableFunc(seq, cmp)
}
