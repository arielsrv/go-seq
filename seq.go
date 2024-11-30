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
//
// This will return true if the sequence was empty.
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
//
// This will return false if the sequence was empty.
func AnyFunc[V any](seq iter.Seq[V], f func(V) bool) bool {
	for v := range seq {
		if f(v) {
			return true
		}
	}
	return false
}

// Append adds values to the end of a sequence.
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

// Chunk splits the values of a sequence into chunks of a fixed size.
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

// Contains determines if a sequence contains a value.
func Contains[V comparable](seq iter.Seq[V], val V) bool {
	for v := range seq {
		if v == val {
			return true
		}
	}
	return false
}

// ContainsFunc determines if a sequence contains a value using a function to select a
// comparable value.
//
// Example:
//
//	var itemID int
//	var items iter.Seq[*Item]
//	found := seq.ContainsFunc(items, itemID, func(i *Item) string {
//		return i.ItemID
//	})
func ContainsFunc[V any, C comparable](seq iter.Seq[V], val C, f func(V) C) bool {
	for v := range seq {
		if f(v) == val {
			return true
		}
	}
	return false
}

// Count returns the number of values in a sequence.
//
// This will iterate over the entire sequence to count the values.
// Use [Any] instead if you only need to check whether the sequence has any values.
// Use [Single] instead if you only need to check whether the sequence has exactly one value.
func Count[V any](seq iter.Seq[V]) int {
	count := 0
	for range seq {
		count++
	}
	return count
}

// CountFunc returns the number of values in a sequence that satisfy a predicate.
//
// This will iterate over the entire sequence to count the values.
// Use [AnyFunc] instead if you only need to check whether the sequence has any values that satisfy
// a predicate.
// Use [SingleFunc] instead if you only need to check whether the sequence has exactly one value that
// satisfies a predicate.
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
//
// The first occurrence is yielded, and any subsequent occurrences are ignored.
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

// DistinctFunc returns distinct values from a sequence using a function to select a comparable value.
//
// The first occurrence is yielded, and any subsequent occurrences are ignored.
func DistinctFunc[V any, C comparable](seq iter.Seq[V], f func(V) C) iter.Seq[V] {
	return func(yield func(V) bool) {
		seen := make(map[C]struct{})
		for v := range seq {
			c := f(v)
			if _, ok := seen[c]; !ok {
				seen[c] = struct{}{}
				if !yield(v) {
					return
				}
			}
		}
	}
}

// First returns the first value of a sequence.
//
// A second return value indicates whether the sequence contained any values.
func First[V any](seq iter.Seq[V]) (V, bool) {
	var v V
	for v = range seq {
		return v, true
	}

	return v, false
}

// Last returns the last value of a sequence.
//
// A second return value indicates whether the sequence contained any values.
func Last[V any](seq iter.Seq[V]) (V, bool) {
	var v V
	found := false
	for v = range seq {
		found = true
	}

	return v, found
}

// Max returns the maximum value in a sequence.
//
// A second return value indicates whether the sequence contained any values.
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
//
// A second return value indicates whether the sequence contained any values.
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

// Prepend adds values to the beginning of a sequence.
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
func Repeat[V any](val V, n int) iter.Seq[V] {
	if n < 0 {
		panic("count must be non-negative")
	}

	return func(yield func(V) bool) {
		for i := 0; i < n; i++ {
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
			out := f(v)
			if !yield(out) {
				return
			}
		}
	}
}

// SelectKeys projects each value of a sequence into a key-value pair using a function to
// select a key.
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
//
// A second return value indicates whether the sequence contained exactly one value.
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

// SingleFunc returns the only value in a sequence that satisfies a predicate.
//
// A second return value indicates whether the sequence contained exactly one value
// that satisfied the predicate.
func SingleFunc[V any](seq iter.Seq[V], f func(V) bool) (V, bool) {
	var first V
	var found bool

	for v := range seq {
		if f(v) {
			if found {
				var zero V
				return zero, false
			}
			first = v
			found = true
		}
	}

	return first, found
}

// Skip bypasses a given number of values in a sequence and returns the remaining values.
func Skip[V any](seq iter.Seq[V], n int) iter.Seq[V] {
	return func(yield func(V) bool) {
		i := 0
		for v := range seq {
			if i >= n {
				if !yield(v) {
					return
				}
			}
			i++
		}
	}
}

// SkipWhile bypasses values in a sequence as long as a condition is true and then
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

// Sum computes the sum of the values in a sequence.
func Sum[V constraints.Integer | constraints.Float](seq iter.Seq[V]) V {
	var sum V
	for v := range seq {
		sum += v
	}
	return sum
}

// Take returns a given number of values from the start of a sequence.
func Take[V any](seq iter.Seq[V], n int) iter.Seq[V] {
	return func(yield func(V) bool) {
		i := 0
		for v := range seq {
			if i >= n {
				return
			}
			if !yield(v) {
				return
			}
			i++
		}
	}
}

// TakeWhile returns values from a sequence as long as a given condition is true
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

// Where filters a sequence based on a predicate.
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

// ValueAt returns the value at a given index in a sequence.
//
// A second return value indicates whether the given index was within the bounds of the sequence.
// This will panic if the given index is negative.
func ValueAt[V any](seq iter.Seq[V], index int) (V, bool) {
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

	var zero V
	return zero, false
}

// Yield returns a sequence of values.
//
// This is useful for creating a sequence from a slice or variadic arguments.
// Providing no arguments creates an empty sequence.
//
// Examples:
//
//	// yield each element of the slice
//	// yields ("a"), ("b"), ("c")
//	vals := seq.Yield([]string{"a", "b", "c"}...)
//
//	// yields (1), (2), (3)
//	vals := seq.Yield(1, 2, 3)
//
//	// empty sequence yields no values
//	empty := seq.Yield[int]()
func Yield[V any](vals ...V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range vals {
			if !yield(v) {
				return
			}
		}
	}
}

// YieldBackward returns a sequence of values in reverse order.
//
// See Yield for more information.
func YieldBackward[V any](vals ...V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for i := len(vals) - 1; i >= 0; i-- {
			if !yield(vals[i]) {
				return
			}
		}
	}
}
