package seq

import (
	"cmp"
	"iter"

	"golang.org/x/exp/constraints"
)

// Aggregate applies an accumulator function over a sequence.
func Aggregate[V, A any](seq iter.Seq[V], init A, f func(A, V) A) A {
	acc := init
	for v := range seq {
		acc = f(acc, v)
	}
	return acc
}

// All determines if all values of a sequence satisfy a condition.
func All[V any](seq iter.Seq[V], f func(V) bool) bool {
	for v := range seq {
		if !f(v) {
			return false
		}
	}
	return true
}

// Any determines if a sequence has any values.
func Any[V any](seq iter.Seq[V]) bool {
	for range seq {
		return true
	}
	return false
}

// AnyFunc determines if any values of a sequence satisfy a condition.
func AnyFunc[V any](seq iter.Seq[V], f func(V) bool) bool {
	for v := range seq {
		if f(v) {
			return true
		}
	}
	return false
}

// Append adds the given values to the end of a sequence.
func Append[V any](seq iter.Seq[V], vals ...V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if !yield(v) {
				return
			}
		}
		for _, v := range vals {
			if !yield(v) {
				return
			}
		}
	}
}

// Chunk splits the values of a sequence into chunks of the given size.
func Chunk[V any](seq iter.Seq[V], size int) iter.Seq[[]V] {
	return func(yield func([]V) bool) {
		chunk := make([]V, 0, size)
		for v := range seq {
			chunk = append(chunk, v)
			if len(chunk) == size {
				if !yield(chunk) {
					return
				}
				chunk = make([]V, 0, size)
			}
		}
		if len(chunk) > 0 {
			if !yield(chunk) {
				return
			}
		}
	}
}

// Concat concatenates multiples sequences into a single sequence.
func Concat[V any](seqs ...iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, seq := range seqs {
			for v := range seq {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Contains determines if the sequence contains the given value.
func Contains[V comparable](seq iter.Seq[V], val V) bool {
	for v := range seq {
		if v == val {
			return true
		}
	}
	return false
}

// Count returns the number of values in the sequence.
//
// Use [Any] instead if you only need to check whether the sequence has any values.
func Count[V any](seq iter.Seq[V]) int {
	count := 0
	for range seq {
		count++
	}
	return count
}

// CountFunc returns the number of values in the sequence that satisfy the predicate.
//
// Use [AnyFunc] instead if you only need to check whether the sequence has any values that satisfy
// the predicate.
func CountFunc[V any](seq iter.Seq[V], f func(V) bool) int {
	count := 0
	for v := range seq {
		if f(v) {
			count++
		}
	}
	return count
}

// Distinct returns distinct values from a sequence.
func Distinct[V comparable](seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		seen := make(map[V]struct{})
		for v := range seq {
			if _, ok := seen[v]; !ok {
				seen[v] = struct{}{}
				if !yield(v) {
					return
				}
			}
		}
	}
}

// DistinctFunc returns distinct values from a sequence according to the given key selector
// function.
func DistinctFunc[V any, K comparable](seq iter.Seq[V], f func(V) K) iter.Seq[V] {
	return func(yield func(V) bool) {
		seen := make(map[K]struct{})
		for v := range seq {
			k := f(v)
			if _, ok := seen[k]; !ok {
				seen[k] = struct{}{}
				if !yield(v) {
					return
				}
			}
		}
	}
}

// ElementAt returns the value at the given index in the sequence.
// A second return value indicates whether the result is valid,
// (i.e., the index was within the bounds of the sequence).
//
// This will panic if the index is negative.
func ElementAt[V any](seq iter.Seq[V], index int) (V, bool) {
	if index < 0 {
		panic("index must be non-negative")
	}

	i := 0
	for v := range seq {
		if i == index {
			return v, true
		}
		i++
	}

	var v V
	return v, false
}

// Empty returns an empty sequence.
func Empty[V any]() iter.Seq[V] {
	return func(yield func(V) bool) {}
}

// First returns the first value of a sequence.
// A second return value indicates whether the result is valid,
// (i.e., there was at least one value in the sequence).
func First[V any](seq iter.Seq[V]) (V, bool) {
	var v V
	for v = range seq {
		return v, true
	}

	return v, false
}

// Last returns the last value of a sequence.
// A second return value indicates whether the result is valid,
// (i.e., there was at least one value in the sequence).
func Last[V any](seq iter.Seq[V]) (V, bool) {
	var v V
	found := false
	for v = range seq {
		found = true
	}

	return v, found
}

// Max returns the maximum value in a sequence.
// A second return value indicates whether the result is valid,
// (i.e., there was at least one value in the sequence).
func Max[V cmp.Ordered](seq iter.Seq[V]) (V, bool) {
	var max V
	var found bool
	for v := range seq {
		if !found || v > max {
			max = v
			found = true
		}
	}
	return max, found
}

// Min returns the minimum value in a sequence.
// A second return value indicates whether the result is valid,
// (i.e., there was at least one value in the sequence).
func Min[V cmp.Ordered](seq iter.Seq[V]) (V, bool) {
	var min V
	var found bool
	for v := range seq {
		if !found || v < min {
			min = v
			found = true
		}
	}
	return min, found
}

// NewSeq returns a sequence that yields the given values.
func NewSeq[V any](vals ...V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range vals {
			if !yield(v) {
				return
			}
		}
	}
}

// Prepend adds the given values to the beginning of the sequence.
func Prepend[V any](seq iter.Seq[V], vals ...V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range vals {
			if !yield(v) {
				return
			}
		}
		for v := range seq {
			if !yield(v) {
				return
			}
		}
	}
}

// Range returns a sequence of integers from start to end (inclusive) with the given step size.
func Range[V constraints.Integer](start, end, step V) iter.Seq[V] {
	if step < 1 {
		panic("step must be positive")
	}

	return func(yield func(V) bool) {
		for i := start; i <= end; i += step {
			if !yield(i) {
				return
			}
		}
	}
}

// Repeat returns a sequence that yields the given value the given number of times.
func Repeat[V any](val V, count int) iter.Seq[V] {
	return func(yield func(V) bool) {
		for i := 0; i < count; i++ {
			if !yield(val) {
				return
			}
		}
	}
}

// Select projects each value of a sequence into a new value.
func Select[V, VOut any](seq iter.Seq[V], f func(V) VOut) iter.Seq[VOut] {
	return func(yield func(VOut) bool) {
		for v := range seq {
			if !yield(f(v)) {
				return
			}
		}
	}
}

// SelectKeys projects each value in the sequence into a key-value pair using the given key selector function.
func SelectKeys[K, V any](seq iter.Seq[V], f func(V) K) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for v := range seq {
			k := f(v)
			if !yield(k, v) {
				return
			}
		}
	}
}

// SelectMany projects each value of a sequence into a sequence and then flattens the resulting sequences
// into a single sequence.
func SelectMany[V, VOut any](seq iter.Seq[V], f func(V) iter.Seq[VOut]) iter.Seq[VOut] {
	return func(yield func(VOut) bool) {
		for v := range seq {
			for out := range f(v) {
				if !yield(out) {
					return
				}
			}
		}
	}
}

// SelectSlices projects each value of a sequence into a slice and then flattens the resulting slices into
// a single sequence.
func SelectSlices[V, VOut any](seq iter.Seq[V], f func(V) []VOut) iter.Seq[VOut] {
	return func(yield func(VOut) bool) {
		for v := range seq {
			for _, out := range f(v) {
				if !yield(out) {
					return
				}
			}
		}
	}
}

// Single returns the only value in a sequence.
// A second return value indicates whether the result is valid,
// (i.e., there was exactly one value in the sequence).
func Single[V any](seq iter.Seq[V]) (V, bool) {
	var first V
	var found bool

	for v := range seq {
		if found {
			return first, false
		}
		first = v
		found = true
	}

	return first, true
}

// Skip bypasses the given number of values in a sequence and returns the remaining values.
func Skip[V any](seq iter.Seq[V], num int) iter.Seq[V] {
	return func(yield func(V) bool) {
		i := 0
		for v := range seq {
			if i >= num {
				if !yield(v) {
					return
				}
			}
			i++
		}
	}
}

// SkipWhile bypasses values in a sequence as long as the given condition is true and then
// returns the remaining values.
func SkipWhile[V any](seq iter.Seq[V], f func(int, V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		skip := true
		i := 0
		for v := range seq {
			if skip {
				if f(i, v) {
					continue
				}
				skip = false
			}
			if !yield(v) {
				return
			}
			i++
		}
	}
}

// Sum computes the sum of the values in the sequence.
func Sum[V constraints.Integer | constraints.Float](seq iter.Seq[V]) V {
	var sum V
	for v := range seq {
		sum += v
	}
	return sum
}

// Take returns the given number of values from the start of the sequence.
func Take[V any](seq iter.Seq[V], num int) iter.Seq[V] {
	return func(yield func(V) bool) {
		i := 0
		for v := range seq {
			if i >= num {
				return
			}
			if !yield(v) {
				return
			}
			i++
		}
	}
}

// TakeWhile returns values from a sequence as long as the given condition is true
// and then skips the remaining values.
func TakeWhile[V any](seq iter.Seq[V], f func(int, V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		i := 0
		for v := range seq {
			if f(i, v) {
				if !yield(v) {
					return
				}
			} else {
				return
			}
			i++
		}
	}
}

// Where filters a value sequence based on the given predicate.
func Where[V any](seq iter.Seq[V], f func(V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if f(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}
