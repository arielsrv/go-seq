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

// CollectLast collects values from a sequence into a new slice,
// keeping only the last n values.
func CollectLast[V any](seq iter.Seq[V], n int) []V {
	if n < 0 {
		panic("seq.CollectLast: n must be non-negative")
	}

	if n == 0 {
		return nil
	}

	buf := make([]V, n)
	i := 0

	// fill a circular buffer
	for v := range seq {
		buf[i%n] = v
		i++
	}

	switch {
	case i == 0:
		buf = nil

	case i < n:
		// return only the values that were filled
		buf = buf[:i]

	case i > n:
		// rotate the buffer to the start of the final n values
		i %= n
		copy(buf, append(buf[i:], buf[:i]...))
	}

	return buf
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

// SortedBy collects values from a sequence into a new slice and then uses a
// function to select a key for sorting.
func SortedBy[V any, K cmp.Ordered](seq iter.Seq[V], f func(V) K) []V {
	s := slices.Collect(seq)

	slices.SortFunc(s, func(a, b V) int {
		return cmp.Compare(f(a), f(b))
	})

	return s
}

// SortedStableBy collects values from a sequence into a new slice and then uses a
// function to select a key for sorting.
func SortedStableBy[V any, K cmp.Ordered](seq iter.Seq[V], f func(V) K) []V {
	s := slices.Collect(seq)

	slices.SortStableFunc(s, func(a, b V) int {
		return cmp.Compare(f(a), f(b))
	})

	return s
}

// SortedFunc collects values from a sequence into a new slice and then sorts it using the
// given comparison function.
func SortedFunc[V any](seq iter.Seq[V], f func(V, V) int) []V {
	return slices.SortedFunc(seq, f)
}

// SortedStableFunc collects values from a sequence into a new slice and then sorts it using the
// given comparison function maintaining the order of equal values.
func SortedStableFunc[V any](seq iter.Seq[V], f func(V, V) int) []V {
	return slices.SortedStableFunc(seq, f)
}
