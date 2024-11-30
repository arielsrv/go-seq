package seq

import (
	"cmp"
	"iter"
	"slices"
)

// Collect collects values from a sequence into a new slice.
func Collect[V any](seq iter.Seq[V]) []V {
	return slices.Collect(seq)
}

// Reversed collects values from a sequence into a new slice and then reverses it.
func Reversed[V any](seq iter.Seq[V]) []V {
	s := slices.Collect(seq)
	slices.Reverse(s)
	return s
}

// Sorted collects values from a sequence into a new slice and then sorts it.
func Sorted[V cmp.Ordered](seq iter.Seq[V]) []V {
	return slices.Sorted(seq)
}

// SortedFunc collects values from a sequence into a new slice and then sorts it using the
// given comparison function.
func SortedFunc[V any](seq iter.Seq[V], cmp func(V, V) int) []V {
	return slices.SortedFunc(seq, cmp)
}

// SortedStable collects values from a sequence into a new slice and then sorts it using the
// given comparison function maintaining the order of equal values.
func SortedStableFunc[V any](seq iter.Seq[V], cmp func(V, V) int) []V {
	return slices.SortedStableFunc(seq, cmp)
}
